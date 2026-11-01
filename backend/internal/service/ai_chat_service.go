package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/service/aiclient"
)

// GuestChatMessage represents a turn in temporary browser-held conversation history.
type GuestChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GuestChatRequestDTO represents incoming guest payload.
type GuestChatRequestDTO struct {
	Message string             `json:"message"`
	Locale  string             `json:"locale"`
	History []GuestChatMessage `json:"history"`
}

// GuestChatResponseDTO represents the stateless output returned to the visitor.
type GuestChatResponseDTO struct {
	Answer    string              `json:"answer"`
	Listings  []ListingSummaryDTO `json:"listings,omitempty"`
	LatencyMS float64             `json:"latency_ms"`
	Model     string              `json:"model"`
}


// AIChatService coordinates conversation state, user message persistence,
// LLM worker communication, and authorized tool dispatching.
type AIChatService struct {
	conversationRepo *repository.AIConversationRepo
	aiClient         *aiclient.Client
	toolDispatcher   *AIToolDispatcher
}

// NewAIChatService creates a new conversational service instance.
func NewAIChatService(
	conversationRepo *repository.AIConversationRepo,
	aiClient *aiclient.Client,
	toolDispatcher *AIToolDispatcher,
) *AIChatService {
	return &AIChatService{
		conversationRepo: conversationRepo,
		aiClient:         aiClient,
		toolDispatcher:   toolDispatcher,
	}
}

// CreateConversation creates a new conversation for an authenticated user.
func (s *AIChatService) CreateConversation(ctx context.Context, userID int, title string) (*models.AIConversation, error) {
	return s.conversationRepo.CreateConversation(ctx, userID, title)
}

// GetConversation retrieves conversation history after validating user ownership.
func (s *AIChatService) GetConversation(ctx context.Context, conversationID string, userID int) (*models.AIConversation, error) {
	return s.conversationRepo.GetConversationByID(ctx, conversationID, userID)
}

// ListConversations retrieves the user's recent conversations.
func (s *AIChatService) ListConversations(ctx context.Context, userID int, limit, offset int) ([]models.AIConversation, error) {
	return s.conversationRepo.ListConversationsByUserID(ctx, userID, limit, offset)
}

// UpdateConversationTitle renames a conversation verifying ownership.
func (s *AIChatService) UpdateConversationTitle(ctx context.Context, conversationID string, userID int, title string) error {
	return s.conversationRepo.UpdateConversationTitle(ctx, conversationID, userID, title)
}

// DeleteConversation deletes a conversation verifying ownership.
func (s *AIChatService) DeleteConversation(ctx context.Context, conversationID string, userID int) error {
	return s.conversationRepo.DeleteConversation(ctx, conversationID, userID)
}

// SendMessage processes a user dialogue turn, invokes the LLM worker,
// executes any requested tools via Go authority, and persists the response.
func (s *AIChatService) SendMessage(
	ctx context.Context,
	userID int,
	conversationID string,
	userContent string,
	locale string,
) (*models.AIMessage, error) {
	cleanContent := strings.TrimSpace(userContent)
	if cleanContent == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}

	locale = strings.ToLower(strings.TrimSpace(locale))
	if locale == "" {
		locale = "en"
	}
	if locale == "or" {
		locale = "om"
	}

	// 1. Authorize user ownership
	conv, err := s.conversationRepo.GetConversationByID(ctx, conversationID, userID)
	if err != nil {
		return nil, fmt.Errorf("unauthorized or invalid conversation: %w", err)
	}

	// 2. Persist user message
	userMsg := &models.AIMessage{
		ConversationID: conversationID,
		Role:           "user",
		Content:        cleanContent,
	}
	if err := s.conversationRepo.CreateMessage(ctx, userMsg); err != nil {
		return nil, fmt.Errorf("failed to persist user message: %w", err)
	}

	// Update conversation title if it is the first user message
	if (conv.Title == "New Chat" || conv.Title == "") && len(conv.Messages) == 0 {
		newTitle := cleanContent
		if len(newTitle) > 40 {
			newTitle = newTitle[:37] + "..."
		}
		_ = s.conversationRepo.UpdateConversationTitle(ctx, conversationID, userID, newTitle)
	}

	// 3. Assemble message history
	allMessages, err := s.conversationRepo.GetMessagesByConversationID(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve conversation messages: %w", err)
	}

	llmMessages := make([]aiclient.ChatLLMMessage, 0, len(allMessages))
	for _, m := range allMessages {
		var calls []aiclient.ChatToolCall
		if len(m.ToolCalls) > 0 {
			for _, tc := range m.ToolCalls {
				calls = append(calls, aiclient.ChatToolCall{
					ID:   tc.ID,
					Type: tc.Type,
					Function: aiclient.ChatFunctionCall{
						Name:      tc.Function.Name,
						Arguments: tc.Function.Arguments,
					},
				})
			}
		}

		llmMessages = append(llmMessages, aiclient.ChatLLMMessage{
			Role:       m.Role,
			Content:    m.Content,
			ToolCallID: m.ToolCallID,
			ToolCalls:  calls,
		})
	}

	// 4. Send request to AI worker
	chatReq := aiclient.ChatMessageRequest{
		ConversationID: &conversationID,
		Messages:       llmMessages,
		Locale:         locale,
		ToolsEnabled:   true,
	}

	resp, err := s.aiClient.SendChatMessage(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("conversational intelligence unavailable: %w", err)
	}

	var inlineListings []ListingSummaryDTO

	// 5. If LLM requested tool execution, execute via Go sovereign authority
	if len(resp.ToolCalls) > 0 {
		var dbToolCalls []models.AIToolCall
		for _, tc := range resp.ToolCalls {
			dbToolCalls = append(dbToolCalls, models.AIToolCall{
				ID:   tc.ID,
				Type: tc.Type,
				Function: models.AIFunctionInvocation{
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				},
			})
		}

		// Save intermediate assistant message with tool calls
		assistantToolCallMsg := &models.AIMessage{
			ConversationID: conversationID,
			Role:           "assistant",
			Content:        "",
			ToolCalls:      dbToolCalls,
		}
		_ = s.conversationRepo.CreateMessage(ctx, assistantToolCallMsg)

		// Execute each tool and append tool result messages
		for _, tc := range resp.ToolCalls {
			resultStr, listings, toolErr := s.toolDispatcher.ExecuteToolCall(ctx, tc)
			if toolErr != nil {
				resultStr = fmt.Sprintf(`{"error": "%s"}`, toolErr.Error())
			}
			if len(listings) > 0 {
				inlineListings = append(inlineListings, listings...)
			}

			toolMsg := &models.AIMessage{
				ConversationID: conversationID,
				Role:           "tool",
				Content:        resultStr,
				ToolCallID:     &tc.ID,
			}
			_ = s.conversationRepo.CreateMessage(ctx, toolMsg)

			// Append to history for the second LLM pass
			llmMessages = append(llmMessages, aiclient.ChatLLMMessage{
				Role:      "assistant",
				ToolCalls: []aiclient.ChatToolCall{tc},
			})
			llmMessages = append(llmMessages, aiclient.ChatLLMMessage{
				Role:       "tool",
				Content:    resultStr,
				ToolCallID: &tc.ID,
			})
		}

		// Re-invoke LLM for final synthesis of tool results
		chatReq.Messages = llmMessages
		finalResp, err := s.aiClient.SendChatMessage(ctx, chatReq)
		if err == nil && finalResp != nil {
			resp = finalResp
		}
	}

	// 6. Persist final assistant response
	content := ""
	if resp.Content != nil {
		content = *resp.Content
	}

	metadata := map[string]interface{}{
		"model":      resp.Model,
		"latency_ms": resp.LatencyMS,
	}
	if len(inlineListings) > 0 {
		metadata["listings"] = inlineListings
	}

	finalAssistantMsg := &models.AIMessage{
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        content,
		Metadata:       metadata,
	}

	if err := s.conversationRepo.CreateMessage(ctx, finalAssistantMsg); err != nil {
		return nil, fmt.Errorf("failed to persist assistant message: %w", err)
	}

	return finalAssistantMsg, nil
}

// ProcessGuestChat handles completely anonymous, stateless AI assistant requests.
// Zero database writes are performed: prompts, answers, and reasoning are never written to PostgreSQL.
func (s *AIChatService) ProcessGuestChat(ctx context.Context, req GuestChatRequestDTO) (*GuestChatResponseDTO, error) {
	t0 := time.Now()
	cleanMsg := strings.TrimSpace(req.Message)
	if cleanMsg == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}
	if len(cleanMsg) > 1000 {
		return nil, fmt.Errorf("message exceeds maximum length of 1000 characters")
	}

	locale := strings.ToLower(strings.TrimSpace(req.Locale))
	if locale == "" {
		locale = "en"
	}
	if locale == "or" {
		locale = "om"
	}

	// 1. Out-of-Domain Guard (Coding tutorials, physics/math equations, biology, trivia)
	if IsOutOfDomain(cleanMsg) {
		return &GuestChatResponseDTO{
			Answer:    GetOutOfDomainRedirection(locale),
			LatencyMS: float64(time.Since(t0).Milliseconds()),
			Model:     "huntiq-guard:domain-policy-v1",
		}, nil
	}

	// 2. Platform FAQ Check (Posting process, registration rules, membership tiers)
	if isFAQ, faqAnswer := CheckPlatformFAQ(cleanMsg, locale); isFAQ {
		return &GuestChatResponseDTO{
			Answer:    faqAnswer,
			LatencyMS: float64(time.Since(t0).Milliseconds()),
			Model:     "huntiq-kb:platform-policy-v1",
		}, nil
	}

	// 3. Marketplace Discovery & Grounded Search
	var searchArgs map[string]interface{}
	if IsBroadProductSearch(cleanMsg) {
		// Broad inquiry (e.g. "what products are available now"): query approved listings
		searchArgs = map[string]interface{}{
			"search_term": "",
			"limit":       float64(5),
		}
	} else {
		// Specific discovery inquiry (e.g. "i wenna phone", "apartment in bole"):
		var chatHistory []aiclient.ChatLLMMessage
		// Bound history to max 4 turns
		hist := req.History
		if len(hist) > 4 {
			hist = hist[len(hist)-4:]
		}
		for _, h := range hist {
			chatHistory = append(chatHistory, aiclient.ChatLLMMessage{
				Role:    h.Role,
				Content: h.Content,
			})
		}
		_, args := aiclient.ParseFallbackSearchArgs(chatHistory, cleanMsg)
		if args == nil {
			args = make(map[string]any)
		}
		if _, ok := args["search_term"]; !ok || args["search_term"] == "" {
			term := aiclient.CleanConversationalSearchTerm(cleanMsg)
			args["search_term"] = term
		}
		args["limit"] = float64(5)
		searchArgs = args
	}

	// Query PostgreSQL for approved, active listings (zero DB mutations)
	_, listings, err := s.toolDispatcher.executeSearchListings(ctx, searchArgs)
	if err != nil {
		listings = nil
	}

	// 4. Grounded Summary Synthesis
	answer := s.synthesizeGroundedResponse(ctx, cleanMsg, listings, searchArgs, locale)

	return &GuestChatResponseDTO{
		Answer:    answer,
		Listings:  listings,
		LatencyMS: float64(time.Since(t0).Milliseconds()),
		Model:     "huntiq-assistant:grounded-v1",
	}, nil
}

func (s *AIChatService) synthesizeGroundedResponse(
	ctx context.Context,
	userQuery string,
	listings []ListingSummaryDTO,
	searchArgs map[string]interface{},
	locale string,
) string {
	if len(listings) == 0 {
		searchTerm, _ := searchArgs["search_term"].(string)
		if searchTerm != "" {
			displayTerm := aiclient.GetLocalizedEntityTerm(searchTerm, locale)
			if displayTerm == "" {
				displayTerm = searchTerm
			}
			switch locale {
			case "am":
				return fmt.Sprintf("በHUNTIQ ዳታቤዝ ውስጥ ለ'%s' በአሁኑ ሰዓት የሚስማማ የተረጋገጠ ዝርዝር አልተገኘም። እባክዎ የፍለጋ ቃላትን በማስተካከል ወይም የዋጋ ገደቡን በማስፋት ይሞክሩ።", displayTerm)
			case "om":
				return fmt.Sprintf("Kuusaa gabaa HUNTIQ keessatti '%s' tiif yeroo ammaa tarreeffamni walsimu hin argamne. Maaloo jechoota biroo fayyadamuun ykn daangaa gatii jijjiiruun yaalaa.", displayTerm)
			case "ti":
				return fmt.Sprintf("ኣብ ዳታቤዝ HUNTIQ ን'%s' ኣብዚ እዋን ዝሰማማዕ ዝተረጋገፀ ዝርዝር ኣይተረኽበን። በጃኹም ናይ ፍለጋ ቃላት ወይ ደረት ዋጋ ብምስትኻኻል ፈትኑ።", displayTerm)
			default:
				return fmt.Sprintf("I searched our live marketplace database for '%s', but found no approved listings matching those exact criteria right now. You might try searching for broader keywords, adjusting your budget, or checking a different location.", displayTerm)
			}
		}

		switch locale {
		case "am":
			return "በHUNTIQ ዳታቤዝ ውስጥ በአሁኑ ሰዓት የሚስማማ የተረጋገጠ ዝርዝር አልተገኘም። እባክዎ የፍለጋ ቃላትን በማስተካከል ይሞክሩ።"
		case "om":
			return "Kuusaa gabaa HUNTIQ keessatti yeroo ammaa tarreeffamni walsimu hin argamne. Maaloo jechoota biroo fayyadamuun yaalaa."
		case "ti":
			return "ኣብ ዳታቤዝ HUNTIQ ኣብዚ እዋን ዝሰማማዕ ዝተረጋገፀ ዝርዝር ኣይተረኽበን። በጃኹም ናይ ፍለጋ ቃላት ብምስትኻኻል ፈትኑ።"
		default:
			return "I searched our live marketplace database, but found no approved listings matching your search right now. Try searching with different keywords or exploring popular categories."
		}
	}

	// Try Ollama inference with strict 2-second timeout if enabled
	if s.aiClient != nil && s.aiClient.IsEnabled() {
		var listingLines []string
		for i, l := range listings {
			action := "For Sale"
			if l.ListingType == "rent" {
				action = "For Rent"
			}
			listingLines = append(listingLines, fmt.Sprintf("%d. %s — %.0f ETB (%s, %s)", i+1, l.Title, l.Price, l.City, action))
		}

		promptContext := fmt.Sprintf(
			"You are HUNTIQ's marketplace assistant. Ground your answer ONLY in these verified listings from PostgreSQL:\n%s\n\nTask: Provide a concise, helpful 1-2 sentence summary in %s answering the user's inquiry: %s. Never invent products, prices, or specs not listed above.",
			strings.Join(listingLines, "\n"), locale, userQuery,
		)

		ollamaCtx, cancel := context.WithTimeout(ctx, 2500*time.Millisecond)
		defer cancel()

		req := aiclient.ChatMessageRequest{
			Messages: []aiclient.ChatLLMMessage{
				{Role: "user", Content: promptContext},
			},
			Locale:       locale,
			ToolsEnabled: false,
		}

		if resp, err := s.aiClient.SendChatMessage(ollamaCtx, req); err == nil && resp != nil && resp.Content != nil {
			cleaned := strings.TrimSpace(*resp.Content)
			if cleaned != "" && !strings.HasPrefix(resp.Model, "fallback:") && !strings.Contains(cleaned, "I am HUNTIQ's marketplace assistant") {
				return cleaned
			}
		}
	}

	// Deterministic truthful synthesis
	var lines []string
	for _, l := range listings {
		typeStr := "For Sale"
		if l.ListingType == "rent" {
			typeStr = "For Rent"
		}
		lines = append(lines, fmt.Sprintf("• **%s** — %.0f %s (%s, %s)", l.Title, l.Price, l.Currency, l.City, typeStr))
	}

	switch locale {
	case "am":
		return fmt.Sprintf("በHUNTIQ ዳታቤዝ ውስጥ በአሁኑ ሰዓት የተገኙ **%d የተረጋገጡ ዝርዝሮች**:\n\n%s\n\nሙሉ መረጃ ለማየትና ከሻጩ/አከራዩ ጋር ለመገናኘት ከስር ያሉትን ካርዶች መጫን ይችላሉ።", len(listings), strings.Join(lines, "\n"))
	case "om":
		return fmt.Sprintf("Kuusaa gabaa HUNTIQ keessatti **tarreeffamoota mirkanaa'an %d** argamaniiru:\n\n%s\n\nOdeeffannoo dabalataaf tarreeffama armaan gadii cuqaasaa.", len(listings), strings.Join(lines, "\n"))
	case "ti":
		return fmt.Sprintf("ኣብ ዳታቤዝ HUNTIQ **%d ዝተረጋገፁ ዝርዝራት** ተረኺቦም ኣለዉ:\n\n%s\n\nዝርዝር ሓበሬታ ንምርኣይ ካብዞም ኣብ ታሕቲ ዘለዉ ካርዶች ምረጹ።", len(listings), strings.Join(lines, "\n"))
	default:
		return fmt.Sprintf("Here are the **%d verified live marketplace listings** matching your search in the HUNTIQ database:\n\n%s\n\nYou can click on any listing card below to view full specifications, photos, and contact the provider directly.", len(listings), strings.Join(lines, "\n"))
	}
}

