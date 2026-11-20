package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// PlatformSettingsGetter defines reading platform settings dynamically.
type PlatformSettingsGetter func(ctx context.Context, key string) string

// AIService orchestrates HUNTIQ AI interactions, domain restrictions, and fallback workflows.
type AIService struct {
	client         ProviderClient
	repo           AIRepository
	tools          *ToolRegistry
	classifier     *DomainClassifier
	settingsGetter PlatformSettingsGetter
	enabled        bool
	rateLimitMu    sync.Mutex
	rateLimitCache map[string][]time.Time
	rateLimitMax   int
}

// NewAIService constructs an AIService instance.
func NewAIService(client ProviderClient, repo AIRepository, enabled bool, rateLimitMax int) *AIService {
	if rateLimitMax <= 0 {
		rateLimitMax = 30
	}
	return &AIService{
		client:         client,
		repo:           repo,
		classifier:     NewDomainClassifier(),
		enabled:        enabled,
		rateLimitCache: make(map[string][]time.Time),
		rateLimitMax:   rateLimitMax,
	}
}

// WithToolRegistry connects the authorized marketplace tool registry.
func (s *AIService) WithToolRegistry(tools *ToolRegistry) *AIService {
	s.tools = tools
	return s
}

// WithSettingsGetter attaches dynamic platform settings lookup.
func (s *AIService) WithSettingsGetter(fn PlatformSettingsGetter) *AIService {
	s.settingsGetter = fn
	return s
}

// IsAIEnabled checks both static config and dynamic platform settings.
func (s *AIService) IsAIEnabled(ctx context.Context) bool {
	if !s.enabled {
		return false
	}
	if s.settingsGetter != nil {
		if s.settingsGetter(ctx, "ai_enabled") == "false" {
			return false
		}
		if s.settingsGetter(ctx, "ai_maintenance_mode") == "true" {
			return false
		}
	}
	return true
}

// CheckRateLimit enforces sliding-window rate limits per session or IP.
func (s *AIService) CheckRateLimit(key string) bool {
	s.rateLimitMu.Lock()
	defer s.rateLimitMu.Unlock()

	now := time.Now()
	cutoff := now.Add(-1 * time.Minute)

	timestamps := s.rateLimitCache[key]
	valid := make([]time.Time, 0, len(timestamps))
	for _, t := range timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= s.rateLimitMax {
		s.rateLimitCache[key] = valid
		return false
	}

	valid = append(valid, now)
	s.rateLimitCache[key] = valid
	return true
}

// DetectLanguage resolves the language from input text or requested parameter.
func (s *AIService) DetectLanguage(text, requestedLang string) string {
	if IsSupportedLanguage(requestedLang) {
		return requestedLang
	}

	// Ge'ez Unicode script range check (Amharic / Tigrinya)
	hasEthiopic := false
	for _, r := range text {
		if r >= 0x1200 && r <= 0x137F { // Ethiopic Unicode block
			hasEthiopic = true
			break
		}
	}

	if hasEthiopic {
		// Differentiate Tigrinya vs Amharic indicators
		lower := strings.ToLower(text)
		if strings.Contains(lower, "ከመይ") || strings.Contains(lower, "ደሊየ") || strings.Contains(lower, "ኣሎ") {
			return LangTI
		}
		return LangAM
	}

	// Afaan Oromo indicators
	lower := strings.ToLower(text)
	if strings.Contains(lower, "akkam") || strings.Contains(lower, "barbaada") || strings.Contains(lower, "biyya") || strings.Contains(lower, "maaloo") {
		return LangOR
	}

	return LangEN
}

// PreScreenDomain classifies whether a prompt is strictly in the HUNTIQ marketplace domain.
func (s *AIService) PreScreenDomain(prompt string) (DomainIntent, bool) {
	lower := strings.ToLower(strings.TrimSpace(prompt))

	// Check Adversarial jailbreak attempts
	adversarialSignals := []string{
		"ignore previous instructions",
		"ignore all previous instructions",
		"disregard previous instructions",
		"pretend you are",
		"jailbreak",
		"bypass security",
		"system prompt",
		"dan mode",
		"reveal secret",
		"drop table",
		"select * from",
	}
	for _, sig := range adversarialSignals {
		if strings.Contains(lower, sig) {
			return IntentAdversarialBypass, false
		}
	}

	// Explicit out-of-domain patterns (Physics, Mathematics, General Academic, Generic Code)
	outOfDomainSignals := []string{
		"e=mc^2",
		"e=mc²",
		"quantum mechanics",
		"general physics",
		"derivative of",
		"integral of",
		"solve for x",
		"write a python script to",
		"write code for",
		"write an essay about",
		"capital of france",
		"who won the world cup",
		"tell me a joke",
	}
	for _, sig := range outOfDomainSignals {
		if strings.Contains(lower, sig) {
			return IntentOutOfDomain, false
		}
	}

	// In-domain marketplace signals
	marketplaceSignals := []string{
		"buy", "sell", "rent", "price", "cost", "house", "apartment", "car", "vehicle", "truck",
		"phone", "laptop", "tools", "equipment", "job", "hire", "service", "electrician",
		"plumber", "tender", "gig", "bole", "addis ababa", "hawassa", "adama", "dire dawa",
		"search", "find", "listing", "product", "huntiq", "account", "register", "contact",
		"መግዛት", "መሸጥ", "ኪራይ", "ዋጋ", "ቤት", "መኪና", "ስልክ", "ላፕቶፕ", "ስራ", "አዲስ አበባ", // Amharic
		"ቢዛ", "ቢቱ", "ቤቃ", "ኮርሲ", "ቢታ", "ሆቴል", // Oromo
		"ምግዛእ", "ምሻጥ", "ክራይ", "ዋጋ", "ገዛ", "መኪና", // Tigrinya
	}
	for _, sig := range marketplaceSignals {
		if strings.Contains(lower, sig) {
			return IntentMarketplaceDiscovery, true
		}
	}

	// If short generic message like greeting or ambiguous query
	if len(strings.Fields(lower)) <= 3 {
		return IntentAmbiguous, true
	}

	return IntentMarketplaceDiscovery, true
}

// BuildSystemInstruction compiles the authoritative system prompt enforcing domain restrictions.
func (s *AIService) BuildSystemInstruction(lang string) string {
	base := `You are HUNTIQ AI, Made by HUNTIQ. You are the official domain-restricted marketplace assistant for HUNTIQ (huntiq.com).

STRICT SCOPE AND RULES:
1. You assist ONLY with HUNTIQ marketplace items: physical products, properties/housing, vehicle transport, services, jobs, gigs, tenders, events, and platform guidance.
2. OUT-OF-DOMAIN QUESTIONS: You MUST politely decline any unrelated questions (e.g., mathematics, physics, general academic essays, general programming tutorials, unrelated trivia).
3. FACTUAL INTEGRITY: Never invent listing IDs, prices, fake phone numbers, fake reviews, or nonexistent listings.
4. IDENTITY: You are HUNTIQ AI, Made by HUNTIQ. Never claim to be standard Google Gemini or an unrestricted assistant.
5. CONCISE & POLITE: Be helpful, direct, and concise.`

	switch lang {
	case LangAM:
		return base + "\nRespond primarily in clear, natural Amharic (አማርኛ). If listing names or technical terms are in English, preserve them accurately."
	case LangTI:
		return base + "\nRespond primarily in clear, natural Tigrinya (ትግርኛ). If listing names or technical terms are in English, preserve them accurately."
	case LangOR:
		return base + "\nRespond primarily in clear, natural Afaan Oromo. If listing names or technical terms are in English, preserve them accurately."
	default:
		return base + "\nRespond in English."
	}
}

// GenerateOutOfDomainReply crafts a polite rejection redirecting the user back to HUNTIQ.
func (s *AIService) GenerateOutOfDomainReply(lang string) string {
	switch lang {
	case LangAM:
		return "እኔ HUNTIQ AI ነኝ፤ የHUNTIQ የገበያ ቦታ ረዳትዎ። በHUNTIQ ላይ ያሉ እቃዎችን፣ ቤቶችን፣ መኪኖችን፣ አገልግሎቶችን፣ ስራዎችን እና የጨረታ እድሎችን ለማግኘት ብቻ ልረዳዎ እችላለሁ። ዛሬ ምን መፈለግ ይፈልጋሉ?"
	case LangTI:
		return "ኣነ HUNTIQ AI እየ፤ ናይ HUNTIQ ዕዳጋ ረዳኢኹም። ኣብ HUNTIQ ዘለዉ ንብረታት፣ ኣባይቲ፣ መካይን፣ ኣገልግሎታት፣ ስራሕን ጨረታታትን ንምርካብ ጥራይ ክሕግዘኩም እኽእል። እንታይ ክትደልዩ ትደልዩ?"
	case LangOR:
		return "Ani HUNTIQ AI dha; gargaaraa gabaa HUNTIQ keessan. Meeshaalee, manneen, konkolaattota, tajaajiloota, hojiilee fi carraawwan caalbaasii HUNTIQ irratti argaman qofa argachuuf isin gargaaruu danda'a. Maal barbaaduu barbaaddu?"
	default:
		return "I am HUNTIQ AI, Made by HUNTIQ — your dedicated marketplace assistant. I can help you discover products, properties, vehicles, services, jobs, and opportunities available on HUNTIQ. What would you like to find today?"
	}
}

// GenerateFallbackReply provides a polite notice when the AI provider is unavailable.
func (s *AIService) GenerateFallbackReply(lang string) string {
	switch lang {
	case LangAM:
		return "HUNTIQ AI በአሁኑ ጊዜ በጊዜያዊ ጥገና ላይ ነው። እባክዎ በቀጥታ በዋናው ገጽ ላይ ባለው የፍለጋ ሳጥን በመጠቀም እቃዎችን እና አገልግሎቶችን መፈለግዎን ይቀጥሉ።"
	case LangTI:
		return "HUNTIQ AI ኣብዚ እዋን ብግዝያዊ ምዕራይ ኣሎ። በጃኹም ብቐጥታ ኣብቲ ዋና ገጽ ዘሎ ናይ ምድላይ ሳንዱቕ ተጠቒምኩም ንብረታትን ኣገልግሎታትን ምድላይኩም ቀጽሉ።"
	case LangOR:
		return "HUNTIQ AI yeroo ammaa suphaa yerootiif qophaa'aa jira. Maaloo kallattiin saanduqa barbaacha fuula duraa fayyadamuun meeshaalee fi tajaajiloota barbaaduu itti fufaa."
	default:
		return "HUNTIQ AI is temporarily undergoing service maintenance. You can continue discovering all products, properties, and services using the manual search and filters on HUNTIQ."
	}
}

// ProcessChat handles the complete conversational request lifecycle.
func (s *AIService) ProcessChat(ctx context.Context, req ChatRequest, userID *int, clientIP string) (*ChatResponse, error) {
	start := time.Now()

	// 1. Language and Domain Classification via DomainClassifier
	classification := s.classifier.Classify(req.Message, req.Language)
	lang := classification.Language

	// 2. Unsupported Language Check (Strict 4-Language Policy)
	if !classification.IsSupportedLang {
		_ = s.repo.RecordAuditEvent(ctx, &AuditEvent{
			SessionID: req.SessionID,
			UserID:    userID,
			EventType: AuditEventUnsupportedLang,
			Details:   map[string]interface{}{"message": req.Message, "detected": classification.Language},
			IPAddress: clientIP,
		})

		return &ChatResponse{
			SessionID:     req.SessionID,
			Language:      LangEN,
			Reply:         GenerateUnsupportedLanguageReply(),
			Intent:        IntentOutOfDomain,
			IsOutOfDomain: true,
			LatencyMs:     time.Since(start).Milliseconds(),
		}, nil
	}

	// 3. AI Enablement Check
	if !s.IsAIEnabled(ctx) {
		return &ChatResponse{
			SessionID:     req.SessionID,
			Language:      lang,
			Reply:         s.GenerateFallbackReply(lang),
			Intent:        IntentMarketplaceDiscovery,
			FallbackUsed:  true,
			IsOutOfDomain: false,
			LatencyMs:     time.Since(start).Milliseconds(),
		}, nil
	}

	// 4. Rate Limit Check
	rateKey := req.SessionID
	if rateKey == "" {
		rateKey = clientIP
	}
	if !s.CheckRateLimit(rateKey) {
		_ = s.repo.RecordAuditEvent(ctx, &AuditEvent{
			SessionID: req.SessionID,
			UserID:    userID,
			EventType: AuditEventRateLimited,
			Details:   map[string]interface{}{"ip": clientIP},
			IPAddress: clientIP,
		})
		return &ChatResponse{
			SessionID:     req.SessionID,
			Language:      lang,
			Reply:         "You have reached the temporary request limit for HUNTIQ AI. Please wait a minute and try again.",
			Intent:        IntentMarketplaceDiscovery,
			FallbackUsed:  false,
			IsOutOfDomain: false,
			LatencyMs:     time.Since(start).Milliseconds(),
		}, nil
	}

	// 5. Conversation lookup or creation
	conv, err := s.repo.GetOrCreateConversation(ctx, req.SessionID, userID, lang)
	if err != nil {
		conv = &Conversation{ID: 0, SessionID: req.SessionID, LanguageCode: lang}
	}

	// 6. Pre-Screening Domain Restrictions & Prompt Injection
	if !classification.IsAllowed {
		eventType := AuditEventOutOfDomain
		if classification.Intent == IntentAdversarialBypass {
			eventType = AuditEventAdversarialBlocked
		}

		_ = s.repo.RecordAuditEvent(ctx, &AuditEvent{
			SessionID: req.SessionID,
			UserID:    userID,
			EventType: eventType,
			Details:   map[string]interface{}{"message": req.Message, "intent": classification.Intent, "reason": classification.Reason},
			IPAddress: clientIP,
		})

		reply := s.GenerateOutOfDomainReply(lang)

		// Record message in conversation history
		if conv.ID > 0 {
			_ = s.repo.SaveMessage(ctx, &Message{
				ConversationID: conv.ID,
				Role:           "user",
				Content:        req.Message,
				Intent:         classification.Intent,
			})
			_ = s.repo.SaveMessage(ctx, &Message{
				ConversationID: conv.ID,
				Role:           "assistant",
				Content:        reply,
				Intent:         classification.Intent,
				LatencyMs:      int(time.Since(start).Milliseconds()),
			})
		}

		return &ChatResponse{
			SessionID:     req.SessionID,
			Language:      lang,
			Reply:         reply,
			Intent:        classification.Intent,
			IsOutOfDomain: true,
			FallbackUsed:  false,
			LatencyMs:     time.Since(start).Milliseconds(),
			SuggestedQuery: []string{
				"Find houses for rent in Bole",
				"Laptops under 50,000 ETB",
				"Plumbers in Addis Ababa",
			},
		}, nil
	}

	// 7. Retrieve Real Authorized Marketplace Listings (Phases Four, Five, Six)
	var matchedCards []ListingCardDTO
	if s.tools != nil {
		searchTerm, minP, maxP := ExtractSearchFilters(req.Message)
		if searchTerm != "" || minP != nil || maxP != nil {
			cards, searchErr := s.tools.SearchPublicListings(ctx, searchTerm, minP, maxP, 6)
			if searchErr == nil && len(cards) > 0 {
				matchedCards = cards
			}
		}
	}

	// 8. Retrieve recent conversation history (multi-turn context)
	var geminiHistory []GeminiContent
	if conv.ID > 0 {
		recentMsgs, _ := s.repo.GetRecentMessages(ctx, conv.ID, 6)
		for _, m := range recentMsgs {
			role := "user"
			if m.Role == "assistant" {
				role = "model"
			}
			geminiHistory = append(geminiHistory, GeminiContent{
				Role:  role,
				Parts: []GeminiPart{{Text: m.Content}},
			})
		}
	}

	// 9. Augment System Instruction with Verified DB Data
	systemInstruction := s.BuildSystemInstruction(lang)
	augmentedPrompt := req.Message
	if len(matchedCards) > 0 {
		augmentedPrompt += "\n\n[AUTHORIZED DATABASE CONTEXT: Real matching listings found in HUNTIQ database:\n"
		for _, c := range matchedCards {
			augmentedPrompt += fmt.Sprintf("- Listing #%d: %s | Price: %.2f ETB | Category: %s | Location: %s\n",
				c.ID, c.Title, c.Price, c.CategoryName, c.LocationName)
		}
		augmentedPrompt += "Summarize or recommend these real verified listings to the user directly without inventing details.]"
	}

	// 10. Gemini Generation Call
	replyText, tokensUsed, genErr := s.client.Generate(ctx, systemInstruction, geminiHistory, augmentedPrompt)
	latency := time.Since(start).Milliseconds()

	if genErr != nil {
		_ = s.repo.RecordAuditEvent(ctx, &AuditEvent{
			SessionID: req.SessionID,
			UserID:    userID,
			EventType: AuditEventProviderError,
			Details:   map[string]interface{}{"error": genErr.Error()},
			IPAddress: clientIP,
		})

		// If real listings were found, we can still present them gracefully even if LLM had a provider timeout!
		fallbackReply := s.GenerateFallbackReply(lang)
		if len(matchedCards) > 0 {
			fallbackReply = fmt.Sprintf("I retrieved %d matching listings from HUNTIQ for you below:", len(matchedCards))
		}

		return &ChatResponse{
			SessionID:     req.SessionID,
			Language:      lang,
			Reply:         fallbackReply,
			Intent:        classification.Intent,
			Listings:      matchedCards,
			FallbackUsed:  true,
			IsOutOfDomain: false,
			LatencyMs:     latency,
		}, nil
	}

	// 11. Persist Turn in Database
	if conv.ID > 0 {
		_ = s.repo.SaveMessage(ctx, &Message{
			ConversationID: conv.ID,
			Role:           "user",
			Content:        req.Message,
			Intent:         classification.Intent,
		})
		_ = s.repo.SaveMessage(ctx, &Message{
			ConversationID: conv.ID,
			Role:           "assistant",
			Content:        replyText,
			Intent:         classification.Intent,
			TokenCount:     tokensUsed,
			LatencyMs:      int(latency),
		})
	}

	_ = s.repo.RecordAuditEvent(ctx, &AuditEvent{
		SessionID: req.SessionID,
		UserID:    userID,
		EventType: AuditEventQuery,
		Details:   map[string]interface{}{"tokens": tokensUsed, "latency_ms": latency, "matches": len(matchedCards)},
		IPAddress: clientIP,
	})

	return &ChatResponse{
		SessionID:      req.SessionID,
		Language:       lang,
		Reply:          replyText,
		Intent:         classification.Intent,
		Listings:       matchedCards,
		IsOutOfDomain:  false,
		FallbackUsed:   false,
		TokensUsed:     tokensUsed,
		LatencyMs:      latency,
		SuggestedQuery: []string{
			"Show pictures of these listings",
			"Filter by nearest location",
			"Sort by lowest price",
		},
	}, nil
}

// HealthCheck provides the status of the AI provider and system.
func (s *AIService) HealthCheck(ctx context.Context) (*AIHealthStatus, error) {
	if !s.IsAIEnabled(ctx) {
		return &AIHealthStatus{
			Status:             "disabled",
			Provider:           "Google Gemini",
			Model:              s.client.GetModel(),
			Enabled:            false,
			LastChecked:        time.Now(),
			SupportedLanguages: []string{LangEN, LangAM, LangTI, LangOR},
		}, nil
	}
	return s.client.HealthCheck(ctx)
}

// GetAdminMetrics returns usage, performance, and cost statistics.
func (s *AIService) GetAdminMetrics(ctx context.Context) (*AIMetricsReport, error) {
	return s.repo.GetAdminMetrics(ctx)
}
