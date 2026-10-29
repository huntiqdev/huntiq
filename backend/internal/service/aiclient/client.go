package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ExtractedPriceRange holds extracted monetary bounds and currency.
type ExtractedPriceRange struct {
	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`
	Currency string   `json:"currency"`
}

// ExtractedEntities holds parsed brand, model, price, condition, location, and attributes.
type ExtractedEntities struct {
	Brand      *string             `json:"brand"`
	Model      *string             `json:"model"`
	PriceRange ExtractedPriceRange `json:"price_range"`
	Condition  *string             `json:"condition"`
	Location   *string             `json:"location"`
	Attributes map[string]any      `json:"attributes"`
}

// QueryUnderstandingResponse represents the validated contract from Python AI worker.
type QueryUnderstandingResponse struct {
	NormalizedQuery      string            `json:"normalized_query"`
	CleanedSearchTerm    string            `json:"cleaned_search_term"`
	SortBy               *string           `json:"sort_by"`
	Mode                 string            `json:"mode"`
	PrimaryIntent        string            `json:"primary_intent"`
	Confidence           float64           `json:"confidence"`
	SuggestedCategoryIDs []int             `json:"suggested_category_ids"`
	ExtractedEntities    ExtractedEntities `json:"extracted_entities"`
	DeepParameters       map[string]any    `json:"deep_parameters"`
	LocalizedAnswer      string            `json:"localized_answer"`
	PIIDetected          bool              `json:"pii_detected"`
	DetectedLocale       string            `json:"detected_locale"`
	ExecutionTimeMS      float64                    `json:"execution_time_ms"`
	IsSafe               bool                       `json:"is_safe"`
	SafetyFlags          []string                   `json:"safety_flags"`
	ReasoningChain       *MarketplaceReasoningChain `json:"reasoning_chain,omitempty"`
}

// ReasoningStep represents a stage in the multi-stage Chain-of-Thought reasoning.
type ReasoningStep struct {
	Step   int    `json:"step"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// MarketplaceReasoningChain holds deep logical reasoning breakdown for complex marketplace queries.
type MarketplaceReasoningChain struct {
	IntentBreakdown    string          `json:"intent_breakdown"`
	MarketFeasibility  string          `json:"market_feasibility"`
	TradeOffAnalysis   string          `json:"trade_off_analysis"`
	InspectionProtocol string          `json:"inspection_protocol"`
	Steps              []ReasoningStep `json:"steps"`
}

// UnderstandQueryParams contains configurable parameters for query inference.
type UnderstandQueryParams struct {
	Query             string `json:"query"`
	Locale            string `json:"locale,omitempty"`
	Mode              string `json:"mode,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
}

// Client defines the AI client communicating with local Ollama, Python worker, or deterministic fallback.
type Client struct {
	baseURL      string
	secretKey    string
	providerType string // "ollama", "python", "mock"
	ollama       *OllamaProvider
	httpClient   *http.Client
	chatClient   *http.Client
	enabled      bool
}

// NewClient creates a new thread-safe AI service client.
func NewClient(baseURL, secretKey string) *Client {
	enabled := (baseURL != "" && secretKey != "")
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression: false,
	}
	return &Client{
		baseURL:      baseURL,
		secretKey:    secretKey,
		providerType: "ollama",
		enabled:      enabled,
		httpClient: &http.Client{
			Timeout:   500 * time.Millisecond,
			Transport: transport,
		},
		chatClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}
}

// WithOllama attaches an Ollama provider pointing to the local daemon.
func (c *Client) WithOllama(baseURL, model string) *Client {
	if c == nil {
		return c
	}
	c.ollama = NewOllamaProvider(baseURL, model)
	c.enabled = true
	return c
}

// WithProvider overrides the active provider strategy ("ollama", "python", "mock").
func (c *Client) WithProvider(provider string) *Client {
	if c == nil {
		return c
	}
	c.providerType = provider
	return c
}

// GetProviderType returns the active provider identifier.
func (c *Client) GetProviderType() string {
	if c == nil {
		return "disabled"
	}
	if c.providerType != "" {
		return c.providerType
	}
	if c.ollama != nil {
		return "ollama"
	}
	if c.baseURL != "" {
		return "python"
	}
	return "mock"
}

// IsEnabled returns true if any AI provider is configured.
func (c *Client) IsEnabled() bool {
	return c != nil && (c.enabled || c.ollama != nil || c.providerType == "mock")
}

// UnderstandQuery sends search query to the internal AI worker with a 400ms context timeout.
// Returns an error if the AI worker is unreachable, fails authentication, or times out,
// signaling the caller to execute immediate deterministic fallback.
func (c *Client) UnderstandQuery(ctx context.Context, query, locale string) (*QueryUnderstandingResponse, error) {
	return c.UnderstandQueryWithParams(ctx, UnderstandQueryParams{
		Query:  query,
		Locale: locale,
		Mode:   "ai",
	})
}

// fallbackUnderstandQuery extracts search parameters deterministically in Go when remote AI worker is offline.
func (c *Client) fallbackUnderstandQuery(params UnderstandQueryParams) *QueryUnderstandingResponse {
	hasCommercial, args := ParseFallbackSearchArgs(nil, params.Query)
	cleanTerm := cleanConversationalSearchTerm(params.Query)
	if hasCommercial && args != nil {
		if st, ok := args["search_term"].(string); ok && st != "" {
			cleanTerm = st
		}
	}

	primaryIntent := "buy_product"
	lower := strings.ToLower(params.Query)
	if strings.Contains(lower, "rent") || strings.Contains(lower, "ኪራይ") || strings.Contains(lower, "kiraa") || strings.Contains(lower, "ክራይ") {
		primaryIntent = "rent_property"
	}

	deepParams := make(map[string]any)
	extracted := ExtractedEntities{}
	if args != nil {
		deepParams["extracted_entities"] = args
		if loc, ok := args["location"].(string); ok && loc != "" {
			extracted.Location = &loc
		}
		if cond, ok := args["condition"].(string); ok && cond != "" {
			extracted.Condition = &cond
		}
		if maxP, ok := args["max_price"].(float64); ok && maxP > 0 {
			extracted.PriceRange.MaxPrice = &maxP
		}
		if minP, ok := args["min_price"].(float64); ok && minP > 0 {
			extracted.PriceRange.MinPrice = &minP
		}
	}

	return &QueryUnderstandingResponse{
		NormalizedQuery:   strings.TrimSpace(params.Query),
		CleanedSearchTerm: cleanTerm,
		PrimaryIntent:     primaryIntent,
		Confidence:        0.90,
		DetectedLocale:    params.Locale,
		ExtractedEntities: extracted,
		DeepParameters:    deepParams,
		IsSafe:            true,
	}
}

// UnderstandQueryWithParams sends parameterized search query options to the internal AI worker.
func (c *Client) UnderstandQueryWithParams(ctx context.Context, params UnderstandQueryParams) (*QueryUnderstandingResponse, error) {
	if c == nil || !c.enabled || params.Query == "" {
		return nil, fmt.Errorf("ai client disabled or empty query")
	}

	reqCtx, cancel := context.WithTimeout(ctx, 400*time.Millisecond)
	defer cancel()

	payload, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ai request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/v1/query/understand", c.baseURL)
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create ai request: %w", err)
	}

	req.Header.Set("X-HUNTIQ-AI-SECRET", c.secretKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ai service unreachable or timed out: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ai service returned non-200 status: %d", resp.StatusCode)
	}

	var result QueryUnderstandingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode ai response: %w", err)
	}

	return &result, nil
}

// ProductCandidateSummary represents a real database product matching the search criteria.
type ProductCandidateSummary struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	SellerTier  string  `json:"seller_tier"`
	Condition   string  `json:"condition"`
	Location    *string `json:"location,omitempty"`
	ListingType string  `json:"listing_type"`
}

// TopProductPick holds the highlighted best Gold or Premium listing.
type TopProductPick struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	Price      float64 `json:"price"`
	SellerTier string  `json:"seller_tier"`
	Badge      string  `json:"badge"`
	Reason     string  `json:"reason"`
}

// ProductDealScore holds an individual product's deal score and badge.
type ProductDealScore struct {
	ProductID    int     `json:"product_id"`
	DealScore    float64 `json:"deal_score"`
	Badge        string  `json:"badge"`
	ValueVerdict string  `json:"value_verdict"`
}

// PriceStats holds market pricing aggregates across actual database listings.
type PriceStats struct {
	MinPrice float64 `json:"min_price"`
	AvgPrice float64 `json:"avg_price"`
	MaxPrice float64 `json:"max_price"`
	Currency string  `json:"currency"`
}

// ProductAnalysisRequest payload sent to Python AI worker.
type ProductAnalysisRequest struct {
	Query             string                    `json:"query"`
	Locale            string                    `json:"locale"`
	PreferredLanguage string                    `json:"preferred_language,omitempty"`
	TotalMatching     int                       `json:"total_matching,omitempty"`
	Products          []ProductCandidateSummary `json:"products"`
}

// ProductAnalysisResponse received from Python AI worker or computed via deterministic fallback.
type ProductAnalysisResponse struct {
	TotalAnalyzed   int                `json:"total_analyzed"`
	TotalMatching   int                `json:"total_matching"`
	PremiumCount    int                `json:"premium_count"`
	GoldCount       int                `json:"gold_count"`
	NormalCount     int                `json:"normal_count"`
	PriceStats      PriceStats         `json:"price_stats"`
	TopPremiumPick  *TopProductPick    `json:"top_premium_pick"`
	TopGoldPick     *TopProductPick    `json:"top_gold_pick"`
	MarketVerdict   string             `json:"market_verdict"`
	DealRatings     []ProductDealScore `json:"deal_ratings"`
	ExecutionTimeMS float64            `json:"execution_time_ms"`
}

// ComputeDeterministicAnalysis calculates real-time product intelligence in pure Go
// if the Python AI service is unreachable, ensuring 100% platform fault tolerance.
func ComputeDeterministicAnalysis(req ProductAnalysisRequest) *ProductAnalysisResponse {
	total := len(req.Products)
	if total == 0 {
		return &ProductAnalysisResponse{
			TotalAnalyzed: 0,
			MarketVerdict: "No matching listings currently found in the marketplace database.",
			PriceStats:    PriceStats{Currency: "ETB"},
		}
	}

	var premCount, goldCount, normCount int
	var minP, maxP, sumP float64
	var bestPrem, bestGold *ProductCandidateSummary

	for i := range req.Products {
		p := &req.Products[i]
		switch p.SellerTier {
		case "premium":
			premCount++
			if bestPrem == nil || (p.Price > 0 && p.Price < bestPrem.Price) {
				bestPrem = p
			}
		case "gold":
			goldCount++
			if bestGold == nil || (p.Price > 0 && p.Price < bestGold.Price) {
				bestGold = p
			}
		default:
			normCount++
		}

		if p.Price > 0 {
			if minP == 0 || p.Price < minP {
				minP = p.Price
			}
			if p.Price > maxP {
				maxP = p.Price
			}
			sumP += p.Price
		}
	}

	avgP := 0.0
	if total > 0 && sumP > 0 {
		avgP = sumP / float64(total)
	}

	var topPrem *TopProductPick
	if bestPrem != nil {
		titleSnip := bestPrem.Title
		if len(titleSnip) > 40 {
			titleSnip = titleSnip[:40]
		}
		topPrem = &TopProductPick{
			ID:         bestPrem.ID,
			Title:      bestPrem.Title,
			Price:      bestPrem.Price,
			SellerTier: "premium",
			Badge:      "👑 Best Premium Deal",
			Reason:     fmt.Sprintf("Top verified listing starting at %.0f ETB.", bestPrem.Price),
		}
	}

	var topG *TopProductPick
	if bestGold != nil {
		topG = &TopProductPick{
			ID:         bestGold.ID,
			Title:      bestGold.Title,
			Price:      bestGold.Price,
			SellerTier: "gold",
			Badge:      "⭐ Best Gold Deal",
			Reason:     fmt.Sprintf("High-trust Gold listing starting at %.0f ETB.", bestGold.Price),
		}
	}

	dealRatings := make([]ProductDealScore, 0, total)
	for _, p := range req.Products {
		score := 65.0
		if avgP > 0 {
			score += ((avgP - p.Price) / avgP) * 25.0
		}
		badge := "Standard Listing"
		verdict := "Fair market listing."
		if p.SellerTier == "premium" {
			score += 15.0
			badge = "👑 Premium Verified"
			verdict = "Top tier verified seller with platform buyer protection."
		} else if p.SellerTier == "gold" {
			score += 10.0
			badge = "⭐ Gold Certified"
			verdict = "Verified Gold seller with high trust standing."
		} else if avgP > 0 && p.Price <= avgP*0.85 {
			badge = "🔥 Great Value"
			verdict = "Competitively priced below average market rates."
		}
		if score < 5.0 {
			score = 5.0
		} else if score > 99.0 {
			score = 99.0
		}
		dealRatings = append(dealRatings, ProductDealScore{
			ProductID:    p.ID,
			DealScore:    score,
			Badge:        badge,
			ValueVerdict: verdict,
		})
	}

	lang := strings.ToLower(strings.TrimSpace(req.PreferredLanguage))
	if lang == "" {
		lang = strings.ToLower(strings.TrimSpace(req.Locale))
	}

	totalMatching := req.TotalMatching
	if totalMatching == 0 {
		totalMatching = total
	}

	var verdictText string
	switch lang {
	case "am":
		scopeText := fmt.Sprintf("በዳታቤዝ ውስጥ %d ንቁ ምርቶች በቅጽበት ተተንትነዋል: ", total)
		if totalMatching > total {
			scopeText = fmt.Sprintf("በዳታቤዝ ውስጥ ካሉ %d አጠቃላይ ተዛማጆች መካከል %d ንቁ ምርቶች በቅጽበት ተተንትነዋል: ", totalMatching, total)
		}
		verdictText = scopeText + fmt.Sprintf("%d ፕሪሚየም እና %d የወርቅ ደረጃ ሻጮች ተገኝተዋል። የገበያ ዋጋ ከ%.0f እስከ %.0f ብር (አማካይ %.0f ብር) ነው።",
			premCount, goldCount, minP, maxP, avgP)
		if topPrem != nil {
			titleSnip := topPrem.Title
			if len(titleSnip) > 40 {
				titleSnip = titleSnip[:40]
			}
			verdictText += fmt.Sprintf(" 💡 የተመረጠ ፕሪሚየም: '%s' በ%.0f ብር።", titleSnip, topPrem.Price)
		}
	case "om":
		scopeText := fmt.Sprintf("Kuusaa odeeffannoo keessatti oomishaalee %d battalumatti xiinxalamaniiru: ", total)
		if totalMatching > total {
			scopeText = fmt.Sprintf("Kuusaa odeeffannoo keessatti oomishaalee %d keessaa %d battalumatti xiinxalamaniiru: ", totalMatching, total)
		}
		verdictText = scopeText + fmt.Sprintf("gurgurtoonni mirkanaa'an Premium %d fi Gold %d jiru. Gatiin gabaa %.0f hanga %.0f ETB (giddu-galeessi %.0f ETB) dha.",
			premCount, goldCount, minP, maxP, avgP)
		if topPrem != nil {
			titleSnip := topPrem.Title
			if len(titleSnip) > 40 {
				titleSnip = titleSnip[:40]
			}
			verdictText += fmt.Sprintf(" 💡 Filatamaa Premium: '%s' gatii %.0f ETB tiin.", titleSnip, topPrem.Price)
		}
	case "ti":
		scopeText := fmt.Sprintf("ኣብ ዳታቤዝ %d ንጡፋት ፍርያት ብቕጽበት ተተንቲኖም: ", total)
		if totalMatching > total {
			scopeText = fmt.Sprintf("ኣብ ዳታቤዝ ካብ ዘለዉ %d ሓፈሻዊ ፍርያት %d ንጡፋት ፍርያት ብቕጽበት ተተንቲኖም: ", totalMatching, total)
		}
		verdictText = scopeText + fmt.Sprintf("%d ፕሪሚየምን %d ወርቂን ደረጃ ዘለዎም ነጋዶ ተረኺቦም። ናይ ዕዳጋ ዋጋ ካብ %.0f ክሳብ %.0f ቅርሺ (ማእከላይ %.0f ቅርሺ) እዩ።",
			premCount, goldCount, minP, maxP, avgP)
		if topPrem != nil {
			titleSnip := topPrem.Title
			if len(titleSnip) > 40 {
				titleSnip = titleSnip[:40]
			}
			verdictText += fmt.Sprintf(" 💡 ፍሉይ ፕሪሚየም: '%s' ብ%.0f ቅርሺ።", titleSnip, topPrem.Price)
		}
	default:
		scopeText := fmt.Sprintf("Real-time market analysis of %d active database listings: ", total)
		if totalMatching > total {
			scopeText = fmt.Sprintf("Real-time market analysis of %d listings (out of %d matching in database): ", total, totalMatching)
		}
		verdictText = scopeText + fmt.Sprintf("%d Premium & %d Gold verified sellers found. Market prices range from %.0f to %.0f ETB (average %.0f ETB).",
			premCount, goldCount, minP, maxP, avgP)
		if topPrem != nil {
			titleSnip := topPrem.Title
			if len(titleSnip) > 40 {
				titleSnip = titleSnip[:40]
			}
			verdictText += fmt.Sprintf(" 💡 Top Verified Pick: '%s' at %.0f ETB.", titleSnip, topPrem.Price)
		}
	}

	return &ProductAnalysisResponse{
		TotalAnalyzed: total,
		TotalMatching: totalMatching,
		PremiumCount:  premCount,
		GoldCount:     goldCount,
		NormalCount:   normCount,
		PriceStats: PriceStats{
			MinPrice: minP,
			AvgPrice: avgP,
			MaxPrice: maxP,
			Currency: "ETB",
		},
		TopPremiumPick: topPrem,
		TopGoldPick:    topG,
		MarketVerdict:  verdictText,
		DealRatings:    dealRatings,
	}
}

// AnalyzeProducts invokes the Python AI worker to perform deep statistical analysis,
// tier classification, and localized market verdict synthesis on real database products.
// If the AI service is disabled, unreachable, or times out (>300ms), it immediately
// falls back to a deterministic Go local calculation so the response never fails.
func (c *Client) AnalyzeProducts(ctx context.Context, req ProductAnalysisRequest) (*ProductAnalysisResponse, error) {
	if len(req.Products) == 0 {
		return ComputeDeterministicAnalysis(req), nil
	}
	if c == nil || !c.enabled {
		return ComputeDeterministicAnalysis(req), nil
	}

	reqCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	payload, err := json.Marshal(req)
	if err != nil {
		return ComputeDeterministicAnalysis(req), nil
	}

	endpoint := fmt.Sprintf("%s/api/v1/query/analyze-products", c.baseURL)
	httpReq, err := http.NewRequestWithContext(reqCtx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return ComputeDeterministicAnalysis(req), nil
	}

	httpReq.Header.Set("X-HUNTIQ-AI-SECRET", c.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		return ComputeDeterministicAnalysis(req), nil
	}
	defer resp.Body.Close()

	var result ProductAnalysisResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ComputeDeterministicAnalysis(req), nil
	}

	return &result, nil
}

// ChatFunctionCall represents the function call invocation from LLM.
type ChatFunctionCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ChatToolCall represents a structured tool call item from LLM.
type ChatToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function ChatFunctionCall `json:"function"`
}

// ChatLLMMessage represents a conversational turn sent to the AI service.
type ChatLLMMessage struct {
	Role       string         `json:"role"`
	Content    string         `json:"content"`
	Name       *string        `json:"name,omitempty"`
	ToolCallID *string        `json:"tool_call_id,omitempty"`
	ToolCalls  []ChatToolCall `json:"tool_calls,omitempty"`
}

// ChatMessageRequest represents the payload sent to POST /api/v1/chat/message.
type ChatMessageRequest struct {
	ConversationID *string          `json:"conversation_id,omitempty"`
	Messages       []ChatLLMMessage `json:"messages"`
	Locale         string           `json:"locale"`
	ToolsEnabled   bool             `json:"tools_enabled"`
}

// ChatMessageResponse represents the returned completion from POST /api/v1/chat/message.
type ChatMessageResponse struct {
	Content      *string        `json:"content"`
	ToolCalls    []ChatToolCall `json:"tool_calls"`
	FinishReason string         `json:"finish_reason"`
	Model        string         `json:"model"`
	LatencyMS    float64        `json:"latency_ms"`
}

// SendChatMessage communicates with the configured AI provider (Ollama, Python, or fallback).
func (c *Client) SendChatMessage(ctx context.Context, req ChatMessageRequest) (*ChatMessageResponse, error) {
	if c == nil || !c.IsEnabled() {
		return nil, fmt.Errorf("ai client is not enabled or configured")
	}

	provider := c.GetProviderType()

	// Extract latest user text
	lastUserText := ""
	for i := len(req.Messages) - 1; i >= 0; i-- {
		if req.Messages[i].Role == "user" {
			lastUserText = req.Messages[i].Content
			break
		}
	}

	// 1. Ollama Provider (Primary Local Model for English/Latin)
	// Ge'ez script (Amharic & Tigrinya) and Ethiopian languages (om, am, ti) are routed to sovereign authority
	// because small 1.5B LLMs lack Ethiopian low-resource tokenizers and falsely trigger safety guardrails.
	isEthiopianLang := req.Locale == "am" || req.Locale == "om" || req.Locale == "ti" || ContainsGeez(lastUserText) || isOromoText(lastUserText)
	if provider == "ollama" && c.ollama != nil && !isEthiopianLang {
		resp, err := c.ollama.SendChat(ctx, req.Messages, defaultOllamaTools())
		if err == nil && resp != nil {
			return resp, nil
		}
		// When Ollama is offline or model is still downloading, fall back cleanly
	}

	// 2. Python AI Worker (if configured)
	if provider == "python" && c.baseURL != "" && c.secretKey != "" {
		resp, err := c.sendPythonChat(ctx, req)
		if err == nil && resp != nil {
			return resp, nil
		}
	}

	// 3. Truthful Deterministic Fallback
	return c.fallbackChat(ctx, req)
}

func (c *Client) sendPythonChat(ctx context.Context, req ChatMessageRequest) (*ChatMessageResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chat request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/v1/chat/message", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create chat request: %w", err)
	}

	httpReq.Header.Set("X-HUNTIQ-AI-SECRET", c.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	client := c.chatClient
	if client == nil {
		client = c.httpClient
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ai service chat request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ai service chat returned status: %d", resp.StatusCode)
	}

	var result ChatMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode chat response: %w", err)
	}
	return &result, nil
}

func (c *Client) fallbackChat(ctx context.Context, req ChatMessageRequest) (*ChatMessageResponse, error) {
	t0 := time.Now()
	if len(req.Messages) == 0 {
		msg := "Hello! How can I assist you today?"
		return &ChatMessageResponse{
			Content:      &msg,
			FinishReason: "stop",
			Model:        "fallback:deterministic-v1",
			LatencyMS:    float64(time.Since(t0).Milliseconds()),
		}, nil
	}

	lastMsg := req.Messages[len(req.Messages)-1]
	userText := strings.TrimSpace(lastMsg.Content)
	lowerQ := strings.ToLower(userText)

	// Check if previous message was a tool result
	for i := len(req.Messages) - 1; i >= 0; i-- {
		if req.Messages[i].Role == "tool" {
			var toolData map[string]any
			_ = json.Unmarshal([]byte(req.Messages[i].Content), &toolData)

			if total, ok := toolData["total_matching"].(float64); ok {
				if total == 0 {
					msg := "I searched our live marketplace database, but found no approved listings matching those exact criteria right now. You might consider broadening your price range or location."
					if req.Locale == "am" || ContainsGeez(userText) {
						msg = "በHUNTIQ ዳታቤዝ ውስጥ በአሁኑ ሰዓት የሚስማማ የተረጋገጠ ዝርዝር አልተገኘም። እባክዎ የፍለጋ ቃላትን ወይም የዋጋ ገደቡን በማስተካከል ይሞክሩ።"
					} else if req.Locale == "om" {
						msg = "Kuusaa gabaa HUNTIQ keessatti yeroo ammaa tarreeffamni walsimu hin argamne. Maaloo daangaa gatii ykn bakka jijjiiruun yaalaa."
					} else if req.Locale == "ti" {
						msg = "ኣብ ዳታቤዝ HUNTIQ ኣብዚ እዋን ዝሰማማዕ ዝተረጋገፀ ዝርዝር ኣይተረኽበን። በጃኹም ናይ ፍለጋ ቃላት ወይ ደረት ዋጋ ብምስትኻኻል ፈትኑ።"
					}
					return &ChatMessageResponse{
						Content:      &msg,
						FinishReason: "stop",
						Model:        "fallback:deterministic-v1",
						LatencyMS:    float64(time.Since(t0).Milliseconds()),
					}, nil
				}
				listings, _ := toolData["listings"].([]any)
				var lines []string
				for _, item := range listings {
					if m, ok := item.(map[string]any); ok {
						title, _ := m["title"].(string)
						price, _ := m["price"].(float64)
						city, _ := m["city"].(string)
						if city == "" {
							city = "Addis Ababa"
						}
						lines = append(lines, fmt.Sprintf("• **%s** — %.0f ETB (%s)", title, price, city))
					}
				}
				msg := fmt.Sprintf("Found **%.0f verified listings** in the marketplace database:\n\n%s\n\nWould you like me to compare any of these?", total, strings.Join(lines, "\n"))
				if req.Locale == "am" || ContainsGeez(userText) {
					msg = fmt.Sprintf("በHUNTIQ ዳታቤዝ ውስጥ **%.0f የተረጋገጡ ዝርዝሮች** ተገኝተዋል:\n\n%s\n\nከእነዚህ ውስጥ አንዱን በዝርዝር ማየት ይፈልጋሉ?", total, strings.Join(lines, "\n"))
				} else if req.Locale == "om" {
					msg = fmt.Sprintf("Kuusaa gabaa HUNTIQ keessatti **tarreeffamoota mirkanaa'an %.0f** argamaniiru:\n\n%s\n\nIsaan kana keessaa kam caalaatti argu barbaaddu?", total, strings.Join(lines, "\n"))
				} else if req.Locale == "ti" {
					msg = fmt.Sprintf("ኣብ ዳታቤዝ HUNTIQ **%.0f ዝተረጋገፁ ዝርዝራት** ተረኺቦም ኣለዉ:\n\n%s\n\nካብዚኣቶም ሓደ ብዝርዝር ክትሪኡ ትደልዩ ዶ?", total, strings.Join(lines, "\n"))
				}
				return &ChatMessageResponse{
					Content:      &msg,
					FinishReason: "stop",
					Model:        "fallback:deterministic-v1",
					LatencyMS:    float64(time.Since(t0).Milliseconds()),
				}, nil
			}

			if avgPrice, ok := toolData["avg_price"].(float64); ok {
				minPrice, _ := toolData["min_price"].(float64)
				maxPrice, _ := toolData["max_price"].(float64)
				total, _ := toolData["total_listings"].(float64)
				term, _ := toolData["search_term"].(string)
				msg := fmt.Sprintf("Based on **%.0f active listings** for '%s' in the database:\n\n• **Average Price:** %.0f ETB\n• **Range:** %.0f ETB – %.0f ETB\n\nWould you like to browse listings within a specific budget?", total, term, avgPrice, minPrice, maxPrice)
				return &ChatMessageResponse{
					Content:      &msg,
					FinishReason: "stop",
					Model:        "fallback:deterministic-v1",
					LatencyMS:    float64(time.Since(t0).Milliseconds()),
				}, nil
			}
			break
		}
	}

	// 1. Greetings
	if lowerQ == "hi" || lowerQ == "hello" || lowerQ == "hey" || lowerQ == "selam" || lowerQ == "ሰላም" || lowerQ == "akkam" {
		msg := "Hello! I am HUNTIQ AI, your intelligent assistant for general inquiries and live marketplace discovery in Ethiopia. How can I help you today?"
		return &ChatMessageResponse{
			Content:      &msg,
			FinishReason: "stop",
			Model:        "fallback:deterministic-v1",
			LatencyMS:    float64(time.Since(t0).Milliseconds()),
		}, nil
	}

	// 2. Out-of-Domain Redirection (Coding tutorials, science equations, general trivia)
	if strings.Contains(lowerQ, "reverse a string") || strings.Contains(lowerQ, "write code") ||
		strings.Contains(lowerQ, "capital of") || strings.Contains(lowerQ, "wrote hamlet") ||
		strings.Contains(lowerQ, "physics") || strings.Contains(lowerQ, "equation") ||
		strings.Contains(lowerQ, "calculus") || strings.Contains(lowerQ, "homework") ||
		strings.Contains(lowerQ, "python") || strings.Contains(lowerQ, "connection pool") {
		msg := "I am HUNTIQ's marketplace assistant. I specialize in helping you find verified products, vehicles, property rentals, and services across Ethiopia, as well as answering questions about using HUNTIQ. For general programming, physics, or academic homework, please consult general-purpose educational resources."
		if req.Locale == "am" || ContainsGeez(userText) {
			msg = "እኔ የHUNTIQ የገበያ ረዳት ነኝ። በኢትዮጵያ ውስጥ ያሉ ምርቶችን፣ መኪናዎችን፣ የቤትና አፓርታማ ኪራዮችን እና አገልግሎቶችን እንዲያገኙ ወይም ስለ HUNTIQ አጠቃቀም መረጃ ለመስጠት ዝግጁ ነኝ። ለትምህርት፣ ለፊዚክስ ወይም ለኮዲንግ ጥያቄዎች እባክዎን አጠቃላይ መድረኮችን ይጠቀሙ።"
		} else if req.Locale == "om" {
			msg = "Ani gargaaraa gabaa HUNTIQ ti. Itoophiyaa keessatti oomishaalee, konkolaataa, kireeffannaa manaa fi tajaajiloota akka argattu ykn waa'ee HUNTIQ akka beektuuf nan gargaara. Gaaffilee barnootaa, koodii fi fiiziksiif maddoota biraa fayyadamaa."
		} else if req.Locale == "ti" {
			msg = "ኣነ ናይ HUNTIQ ናይ ዕዳጋ ሓጋዚ እየ። ኣብ ኢትዮጵያ ፍርያት፣ መካይን፣ ናይ ገዛ ክራይን ካልኦት ኣገልግሎታትን ንምርካብ ወይ ብዛዕባ HUNTIQ ንምሕታት ክሕግዘኩም እኽእል እየ። ንትምህርታዊ፣ ፊዚክስ ወይ ኮዲንግ ሕቶታት ካልእ መድረኽ ተጠቐሙ።"
		}
		return &ChatMessageResponse{
			Content:      &msg,
			FinishReason: "stop",
			Model:        "fallback:domain-policy-v1",
			LatencyMS:    float64(time.Since(t0).Milliseconds()),
		}, nil
	}

	// 3. Commercial Marketplace Search Intent -> Emit Grounded Tool Call
	if req.ToolsEnabled {
		if hasIntent, args := parseFallbackSearchArgs(req.Messages, userText); hasIntent {
			tc := ChatToolCall{
				ID:   fmt.Sprintf("call_%d", time.Now().UnixMilli()),
				Type: "function",
				Function: ChatFunctionCall{
					Name:      "search_marketplace_listings",
					Arguments: args,
				},
			}
			statusMsg := "Searching the live marketplace database for you..."
			if req.Locale == "am" || ContainsGeez(userText) {
				statusMsg = "የቀጥታ ዳታቤዝ ዝርዝሮችን በመፈለግ ላይ..."
			} else if req.Locale == "om" {
				statusMsg = "Kuusaa gabaa irraa barbaadaa jira..."
			} else if req.Locale == "ti" {
				statusMsg = "ናይ ቀጥታ ዝርዝራት ይድለ ኣሎ..."
			}
			return &ChatMessageResponse{
				Content:      &statusMsg,
				ToolCalls:    []ChatToolCall{tc},
				FinishReason: "tool_calls",
				Model:        "fallback:deterministic-v1",
				LatencyMS:    float64(time.Since(t0).Milliseconds()),
			}, nil
		}
	}

	// 4. Default Marketplace Focused Response
	respMsg := "I am HUNTIQ's marketplace assistant. I can help you search verified products, vehicles, and property rentals across Ethiopia, or answer questions about buying, selling, and posting listings on HUNTIQ."
	if req.Locale == "am" || ContainsGeez(userText) {
		respMsg = "እኔ የHUNTIQ የገበያ ረዳት ነኝ። በኢትዮጵያ ውስጥ የተረጋገጡ ምርቶችን፣ ተሽከርካሪዎችን፣ የቤት ኪራዮችን ለመፈለግ ወይም በHUNTIQ ላይ ለመግዛትና ለመሸጥ የሚረዱ መመሪያዎችን ልሰጥዎ ዝግጁ ነኝ።"
	}
	return &ChatMessageResponse{
		Content:      &respMsg,
		FinishReason: "stop",
		Model:        "fallback:deterministic-v1",
		LatencyMS:    float64(time.Since(t0).Milliseconds()),
	}, nil
}


func defaultOllamaTools() []any {
	return []any{
		map[string]any{
			"type": "function",
			"function": map[string]any{
				"name":        "search_marketplace_listings",
				"description": "Search approved products, vehicles, and real estate in the live HUNTIQ database",
				"parameters": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"search_term": map[string]any{"type": "string", "description": "Search keyword, brand, or model"},
						"min_price":   map[string]any{"type": "number", "description": "Optional minimum price in ETB. ONLY include if the user specifically asked for a price floor."},
						"max_price":   map[string]any{"type": "number", "description": "Optional maximum price in ETB. ONLY include if the user specifically requested a budget limit or maximum price."},
						"location":    map[string]any{"type": "string", "description": "Optional Ethiopian city or subcity (e.g. Bole). ONLY include if user specifically mentioned a location."},
						"condition":   map[string]any{"type": "string", "description": "Optional condition ('new' or 'used'). ONLY include if user specifically requested new or used."},
					},
					"required": []string{"search_term"},
				},
			},
		},
		map[string]any{
			"type": "function",
			"function": map[string]any{
				"name":        "get_marketplace_stats",
				"description": "Calculate average price and listing volume statistics from live database",
				"parameters": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"search_term": map[string]any{"type": "string", "description": "Product or vehicle model"},
						"location":    map[string]any{"type": "string", "description": "Location filter"},
					},
					"required": []string{"search_term"},
				},
			},
		},
	}
}

var knownEthiopianLocations = []string{
	"bole", "piassa", "merkato", "kazanchis", "megenagna", "cmc", "sarbet", "gerji",
	"jomo", "summit", "yeka", "kirkos", "arada", "gullele", "kolfe", "nifas silk",
	"akaki", "addis ababa", "hawassa", "adama", "bahir dar", "mekelle", "dire dawa",
	"gondar", "jimma", "ቦሌ", "ፒያሳ", "መርካቶ", "መገናኛ", "ሃዋሳ", "አዳማ", "ባህር ዳር", "መቀሌ", "ድሬ ዳዋ",
}

var (
	rxPriceUnder = regexp.MustCompile(`(?i)(?:under|below|less than|max(?:imum)?|up to|ከ|በታች|hanga)\s*(\d+(?:,\d+)*(?:\.\d+)?)\s*(k|thousand|ሺህ|etb|birr)?`)
	rxPriceAbove = regexp.MustCompile(`(?i)(?:above|over|more than|min(?:imum)?|starting from|ከ|በላይ)\s*(\d+(?:,\d+)*(?:\.\d+)?)\s*(k|thousand|ሺህ|etb|birr)?`)
	rxPriceK     = regexp.MustCompile(`(?i)\b(\d+(?:,\d+)*(?:\.\d+)?)\s*k\b`)
	rxPriceETB   = regexp.MustCompile(`(?i)\b(\d+(?:,\d+)*(?:\.\d+)?)\s*(?:etb|birr)\b`)
)

func parsePriceVal(numStr string, multiplierStr string) float64 {
	cleaned := strings.ReplaceAll(numStr, ",", "")
	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0
	}
	m := strings.ToLower(strings.TrimSpace(multiplierStr))
	if m == "k" || m == "thousand" || m == "ሺህ" {
		val *= 1000
	}
	return val
}

func detectCondition(text string) string {
	lower := strings.ToLower(text)
	if strings.Contains(lower, "brand new") || strings.Contains(lower, "አዲስ") || strings.Contains(lower, "haaraa") || strings.Contains(lower, "ሓድሽ") {
		return "new"
	}
	if strings.Contains(lower, "used") || strings.Contains(lower, "secondhand") || strings.Contains(lower, "second hand") || strings.Contains(lower, "ያገለገለ") || strings.Contains(lower, "duloma") {
		return "used"
	}
	return ""
}

func detectListingType(text string) string {
	lower := strings.ToLower(text)
	if strings.Contains(lower, "rent") || strings.Contains(lower, "rental") || strings.Contains(lower, "lease") || strings.Contains(lower, "ኪራይ") || strings.Contains(lower, "kiraa") || strings.Contains(lower, "ክራይ") {
		return "rent"
	}
	if strings.Contains(lower, "buy") || strings.Contains(lower, "sale") || strings.Contains(lower, "purchase") || strings.Contains(lower, "መግዛት") || strings.Contains(lower, "bituu") || strings.Contains(lower, "ክገዝእ") {
		return "sale"
	}
	return ""
}

func detectLocation(text string) string {
	lower := strings.ToLower(text)
	for _, loc := range knownEthiopianLocations {
		if strings.Contains(lower, loc) {
			if len(loc) > 0 && loc[0] >= 'a' && loc[0] <= 'z' {
				return strings.ToUpper(loc[:1]) + loc[1:]
			}
			return loc
		}
	}
	return ""
}

func isCommercialIntent(text string) bool {
	lower := strings.ToLower(text)

	// Exclude programming or trivia queries
	if strings.Contains(lower, "capital of") || strings.Contains(lower, "reverse a string") ||
		strings.Contains(lower, "tell me a joke") || strings.Contains(lower, "who are you") ||
		strings.Contains(lower, "write code") || strings.Contains(lower, "write a python") {
		return false
	}

	commercialKeywords := []string{
		"phone", "iphone", "samsung", "pixel", "android", "mobile",
		"car", "toyota", "vitz", "corolla", "yaris", "suzuki", "hyundai", "vehicle",
		"laptop", "macbook", "dell", "hp", "lenovo", "computer", "pc",
		"house", "apartment", "condo", "villa", "studio", "real estate", "land",
		"tv", "television", "camera", "watch", "shoe", "shoes", "clothes", "dress",
		"buy", "sell", "rent", "cheap", "cost", "price", "looking for", "find me",
		"search for", "show me", "product", "products", "listing", "listings",
		"available", "store", "wenna", "wanna", "want", "need",
		"ስልክ", "መኪና", "ቤት", "ኮምፒውተር", "ላፕቶፕ", "እቃ", "ምርት", "መግዛት", "ፈልጋለሁ",
		"ኪራይ", "ዋጋ", "አለ", "ሽያጭ", "ጫማ", "ልብስ", "ቦሌ", "አፓርታማ", "ኮንዶሚኒየም",
		"bilbila", "konkolaataa", "mana", "kompiitara", "bituu", "barbaada", "kiraa", "gatii", "meeshaa",
		"ሞባይል", "መኪና", "ገዛ", "ኮምፒውተር", "ክገዝእ", "ደልየ", "ክራይ", "ዋጋ", "ንብረት",
	}

	for _, kw := range commercialKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}

	if rxPriceUnder.MatchString(lower) || rxPriceK.MatchString(lower) || strings.Contains(lower, "only ones") {
		return true
	}

	return false
}

// CleanConversationalSearchTerm strips conversational filler, pronouns, prepositions,
// and location/price tokens, isolating the pure item or brand name for PostgreSQL full-text search.
func CleanConversationalSearchTerm(q string) string {
	cleaned := strings.TrimSpace(q)
	if cleaned == "" {
		return ""
	}

	// 1. Strip punctuation
	for _, ch := range []string{"?", "!", ".", ",", "\"", "'", ":", ";", "(", ")", "[", "]"} {
		cleaned = strings.ReplaceAll(cleaned, ch, " ")
	}

	lower := strings.ToLower(cleaned)

	// Remove price expressions
	lower = rxPriceUnder.ReplaceAllString(lower, " ")
	lower = rxPriceAbove.ReplaceAllString(lower, " ")
	lower = rxPriceK.ReplaceAllString(lower, " ")
	lower = rxPriceETB.ReplaceAllString(lower, " ")

	// Remove conversational prefixes and phrases in English
	fillerPhrases := []string{
		"what products are available now",
		"what products are available",
		"what are available now",
		"what are available",
		"what products do you have",
		"what is available now",
		"what is available",
		"can you find me a",
		"can you find me",
		"can you find",
		"can you show me",
		"can you show",
		"i want to buy a",
		"i want to buy",
		"i want to get",
		"i want to rent",
		"i need to buy",
		"i need to find",
		"i need to rent",
		"i want a",
		"i want an",
		"i want",
		"i wenna a",
		"i wenna an",
		"i wenna",
		"i wanna a",
		"i wanna an",
		"i wanna",
		"i need a",
		"i need an",
		"i need",
		"looking for a",
		"looking for an",
		"looking for",
		"look for a",
		"look for",
		"find me a",
		"find me an",
		"find me",
		"find a",
		"find an",
		"search for a",
		"search for an",
		"search for",
		"search a",
		"show me all",
		"show me",
		"please find me",
		"please find",
		"please show me",
		"please show",
		"where can i buy",
		"where can i find",
		"do you have any",
		"do you have a",
		"do you have",
		"are there any",
		"is there any",
		"give me a",
		"give me",
		"only ones",
		"only one",
		"near me",
		"around me",
		"for rent",
		"for sale",
		"brand new",
		"second hand",
		"secondhand",
		"used only",
	}

	sort.Slice(fillerPhrases, func(i, j int) bool {
		return len(fillerPhrases[i]) > len(fillerPhrases[j])
	})
	for _, phrase := range fillerPhrases {
		lower = strings.ReplaceAll(lower, phrase, " ")
	}

	// Remove Ethiopian conversational filler phrases
	ethiopianFillers := []string{
		"መግዛት እፈልጋለሁ", "መግዛት እፈልጋለው", "መግዛት ፈልጌ ነበር", "መግዛት ፈልጌ ነው",
		"እፈልጋለሁ", "እፈልጋለው", "ፈልጌ ነበር", "ፈልጌ ነው",
		"አለ ወይ", "አለ?", "እባክህ", "እባክሽ", "እባክዎ",
		"የሚሸጥ", "የሚከራይ", "ያገለገለ", "አዲስ",
		"bituu barbaada", "bituu barbaadu", "barbaada", "bituu",
		"ክገዝእ ደልየ", "ክገዝእ", "ደልየ", "ሓድሽ",
	}
	sort.Slice(ethiopianFillers, func(i, j int) bool {
		return len(ethiopianFillers[i]) > len(ethiopianFillers[j])
	})
	for _, filler := range ethiopianFillers {
		lower = strings.ReplaceAll(lower, filler, " ")
	}

	// Remove known location mentions
	for _, loc := range knownEthiopianLocations {
		lower = strings.ReplaceAll(lower, "in "+loc, " ")
		lower = strings.ReplaceAll(lower, "at "+loc, " ")
		lower = strings.ReplaceAll(lower, loc, " ")
	}

	// Remove stop words
	words := strings.Fields(lower)
	var finalWords []string
	stopWords := map[string]bool{
		"a": true, "an": true, "the": true, "in": true, "at": true, "on": true,
		"for": true, "to": true, "of": true, "with": true, "and": true, "or": true,
		"is": true, "are": true, "now": true, "me": true, "my": true, "please": true,
		"some": true, "any": true, "all": true, "only": true, "ones": true, "one": true,
		"used": true, "new": true, "product": true, "products": true, "item": true, "items": true,
		"listing": true, "listings": true,
	}

	for _, w := range words {
		if !stopWords[w] && len(w) > 0 {
			finalWords = append(finalWords, w)
		}
	}

	return strings.TrimSpace(strings.Join(finalWords, " "))
}

func cleanConversationalSearchTerm(q string) string {
	return CleanConversationalSearchTerm(q)
}

// ParseFallbackSearchArgs extracts structured search filters from user intent and dialogue history.
func ParseFallbackSearchArgs(messages []ChatLLMMessage, userText string) (bool, map[string]any) {
	if !isCommercialIntent(userText) {
		hasCommercialHistory := false
		for i := len(messages) - 1; i >= 0; i-- {
			if isCommercialIntent(messages[i].Content) {
				hasCommercialHistory = true
				break
			}
		}
		if !hasCommercialHistory {
			return false, nil
		}
	}

	args := make(map[string]any)

	// Extract attributes from current turn
	cleanTerm := cleanConversationalSearchTerm(userText)

	// Multi-turn inheritance if current cleanTerm is empty
	if cleanTerm == "" {
		for i := len(messages) - 1; i >= 0; i-- {
			prev := messages[i]
			if prev.Role == "user" && prev.Content != userText {
				prevClean := cleanConversationalSearchTerm(prev.Content)
				if prevClean != "" {
					cleanTerm = prevClean
					break
				}
			} else if prev.Role == "assistant" && len(prev.ToolCalls) > 0 {
				for _, tc := range prev.ToolCalls {
					if st, ok := tc.Function.Arguments["search_term"].(string); ok && st != "" {
						cleanTerm = st
						break
					}
				}
				if cleanTerm != "" {
					break
				}
			}
		}
	}

	catalogTerm := cleanTerm
	if mapped, ok := ethiopianKeywordMap[cleanTerm]; ok {
		catalogTerm = mapped
	} else {
		for k, v := range ethiopianKeywordMap {
			if strings.Contains(cleanTerm, k) {
				catalogTerm = v
				break
			}
		}
	}

	args["search_term"] = catalogTerm
	if catalogTerm != cleanTerm {
		args["localized_term"] = cleanTerm
	}

	// Extract condition
	cond := detectCondition(userText)
	if cond != "" {
		args["condition"] = cond
	}

	// Extract listing type
	lt := detectListingType(userText)
	if lt != "" {
		args["listing_type"] = lt
	}

	// Extract location
	loc := detectLocation(userText)
	if loc != "" {
		args["location"] = loc
	}

	// Extract max_price
	if m := rxPriceUnder.FindStringSubmatch(userText); len(m) > 1 {
		multiplier := ""
		if len(m) > 2 {
			multiplier = m[2]
		}
		val := parsePriceVal(m[1], multiplier)
		if val > 0 {
			args["max_price"] = val
		}
	} else if m := rxPriceK.FindStringSubmatch(userText); len(m) > 1 && strings.Contains(strings.ToLower(userText), "under") {
		val := parsePriceVal(m[1], "k")
		if val > 0 {
			args["max_price"] = val
		}
	}

	// Extract min_price
	if m := rxPriceAbove.FindStringSubmatch(userText); len(m) > 1 {
		multiplier := ""
		if len(m) > 2 {
			multiplier = m[2]
		}
		val := parsePriceVal(m[1], multiplier)
		if val > 0 {
			args["min_price"] = val
		}
	}

	// Inherit historical filters if omitted in current turn
	if len(messages) > 0 {
		for i := len(messages) - 1; i >= 0; i-- {
			prev := messages[i]
			if prev.Role == "assistant" && len(prev.ToolCalls) > 0 {
				for _, tc := range prev.ToolCalls {
					if _, hasLoc := args["location"]; !hasLoc {
						if prevLoc, ok := tc.Function.Arguments["location"].(string); ok && prevLoc != "" {
							args["location"] = prevLoc
						}
					}
					if _, hasCond := args["condition"]; !hasCond {
						if prevCond, ok := tc.Function.Arguments["condition"].(string); ok && prevCond != "" {
							args["condition"] = prevCond
						}
					}
					if _, hasLT := args["listing_type"]; !hasLT {
						if prevLT, ok := tc.Function.Arguments["listing_type"].(string); ok && prevLT != "" {
							args["listing_type"] = prevLT
						}
					}
				}
			}
		}
	}

	return true, args
}

func parseFallbackSearchArgs(messages []ChatLLMMessage, userText string) (bool, map[string]any) {
	return ParseFallbackSearchArgs(messages, userText)
}

// ContainsGeez returns true if the string contains any Ethiopic / Ge'ez unicode code points.
func ContainsGeez(s string) bool {
	for _, r := range s {
		if r >= 0x1200 && r <= 0x137F {
			return true
		}
	}
	return false
}
func isOromoText(text string) bool {
	lower := strings.ToLower(text)
	oromoWords := []string{"bituu", "barbaada", "barbaadu", "bilbila", "konkolaataa", "mana", "gatii", "kiraa", "meeshaa", "akkam", "jirtu"}
	for _, w := range oromoWords {
		if strings.Contains(lower, w) {
			return true
		}
	}
	return false
}
var ethiopianKeywordMap = map[string]string{
	"ስልክ":         "phone",
	"ሞባይል":       "phone",
	"bilbila":     "phone",
	"መኪና":        "car",
	"konkolaataa": "car",
	"ቤት":         "house",
	"mana":        "house",
	"ገዛ":          "house",
	"ኮምፒውተር":     "laptop",
	"ላፕቶፕ":       "laptop",
	"kompiitara":  "laptop",
	"ጫማ":          "shoes",
	"ልብስ":        "clothes",
	"አፓርታማ":      "apartment",
	"condo":       "apartment",
}



