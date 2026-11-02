package service

import (
	"context"
	"crypto/rand"
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
	Message   string             `json:"message"`
	Locale    string             `json:"locale"`
	SessionID string             `json:"session_id,omitempty"`
	History   []GuestChatMessage `json:"history,omitempty"`
}

// GuestChatResponseDTO represents the context-aware output returned to the visitor.
type GuestChatResponseDTO struct {
	SessionID      string                         `json:"session_id,omitempty"`
	GreetingPrefix string                         `json:"greeting_prefix,omitempty"`
	Answer         string                         `json:"answer"`
	CategoryGroups []MultiCategorySearchResultDTO `json:"category_groups,omitempty"`
	Listings       []ListingSummaryDTO            `json:"listings,omitempty"`
	Clarification  *string                        `json:"clarification,omitempty"`
	LatencyMS      float64                        `json:"latency_ms"`
	Model          string                         `json:"model"`
}

// AIChatService coordinates conversation state, user message persistence,
// LLM worker communication, and authorized tool dispatching.
type AIChatService struct {
	conversationRepo *repository.AIConversationRepo
	aiClient         *aiclient.Client
	toolDispatcher   *AIToolDispatcher
	sessionManager   *SessionManager
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
		sessionManager:   NewSessionManager(DefaultSessionTTL, DefaultMaxSessions),
	}
}

// SetSessionManager allows injecting a custom session manager (e.g. for testing).
func (s *AIChatService) SetSessionManager(sm *SessionManager) {
	s.sessionManager = sm
}

// SessionManager returns the underlying in-memory session manager.
func (s *AIChatService) SessionManager() *SessionManager {
	return s.sessionManager
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

// ProcessGuestChat handles completely anonymous, context-aware AI assistant requests.
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

	// 1. Session resolution
	sessionID := strings.TrimSpace(req.SessionID)
	if sessionID == "" {
		sessionID = generateSessionID()
	}

	var sessionCtx *ConversationContext
	if s.sessionManager != nil {
		sessionCtx = s.sessionManager.GetOrCreate(sessionID, locale)
		// Hydrate session turns if fresh and client provided history
		if sessionCtx != nil && len(sessionCtx.Turns) == 0 && len(req.History) > 0 {
			for _, h := range req.History {
				sessionCtx.Turns = append(sessionCtx.Turns, ConversationTurn{
					Role:      h.Role,
					Content:   h.Content,
					Timestamp: time.Now(),
				})
			}
		}
	}

	// 2. Multidimensional intent classification and conversational planning
	plan := ClassifyMessage(cleanMsg, locale, sessionCtx)
	effectiveLocale := plan.DetectedLanguage
	if effectiveLocale == "" {
		effectiveLocale = locale
	}

	// 3. Out-of-Domain Guard (Coding tutorials, physics/math equations, biology, trivia)
	if containsIntent(plan.Intents, IntentOutOfDomain) || IsOutOfDomain(cleanMsg) {
		ans := GetOutOfDomainRedirection(effectiveLocale)
		if s.sessionManager != nil {
			s.sessionManager.RecordTurn(sessionID, cleanMsg, ans, plan, nil)
		}
		return &GuestChatResponseDTO{
			SessionID: sessionID,
			Answer:    ans,
			LatencyMS: float64(time.Since(t0).Milliseconds()),
			Model:     "huntiq-guard:domain-policy-v1",
		}, nil
	}

	// 4. Platform FAQ Check (Posting process, registration rules, membership tiers)
	if isFAQ, faqAnswer := CheckPlatformFAQ(cleanMsg, effectiveLocale); isFAQ {
		if s.sessionManager != nil {
			s.sessionManager.RecordTurn(sessionID, cleanMsg, faqAnswer, plan, nil)
		}
		return &GuestChatResponseDTO{
			SessionID: sessionID,
			Answer:    faqAnswer,
			LatencyMS: float64(time.Since(t0).Milliseconds()),
			Model:     "huntiq-kb:platform-policy-v1",
		}, nil
	}

	// 5. Incomplete / Ambiguous Word Clarification Needed
	if plan.RequiresClarification && len(plan.Searches) == 0 {
		clarification := plan.ClarificationQuestion
		if s.sessionManager != nil {
			s.sessionManager.RecordTurn(sessionID, cleanMsg, clarification, plan, nil)
		}
		return &GuestChatResponseDTO{
			SessionID:     sessionID,
			Answer:        clarification,
			Clarification: &clarification,
			LatencyMS:     float64(time.Since(t0).Milliseconds()),
			Model:         "huntiq-assistant:clarification-v1",
		}, nil
	}

	// 6. Conversational / Social / Identity Inquiries (Zero DB search)
	if !plan.RequiresDatabaseSearch && len(plan.Searches) == 0 {
		var ansParts []string
		if plan.SocialGreeting != "" {
			ansParts = append(ansParts, plan.SocialGreeting)
		}
		if plan.IdentityResponse != "" {
			ansParts = append(ansParts, plan.IdentityResponse)
		} else if plan.SocialFarewell != "" {
			ansParts = append(ansParts, plan.SocialFarewell)
		} else if plan.GeneralResponse != "" {
			ansParts = append(ansParts, plan.GeneralResponse)
		}

		ans := strings.Join(ansParts, " ")
		if ans == "" {
			ans = buildDefaultConversationalResponse(effectiveLocale)
		}

		if s.sessionManager != nil {
			s.sessionManager.RecordTurn(sessionID, cleanMsg, ans, plan, nil)
		}
		return &GuestChatResponseDTO{
			SessionID:      sessionID,
			GreetingPrefix: plan.SocialGreeting,
			Answer:         ans,
			LatencyMS:      float64(time.Since(t0).Milliseconds()),
			Model:          "huntiq-assistant:conversational-v2",
		}, nil
	}

	// 7. Marketplace Search & Multi-Category Execution
	var categoryGroups []MultiCategorySearchResultDTO
	var flatListings []ListingSummaryDTO

	if s.toolDispatcher != nil {
		categoryGroups = s.toolDispatcher.ExecuteMultiSearch(ctx, plan.Searches)
		for _, grp := range categoryGroups {
			flatListings = append(flatListings, grp.Listings...)
		}
	}

	// 8. Grounded Summary Synthesis
	answer := s.synthesizeMultiCategoryGroundedResponse(ctx, cleanMsg, plan.SocialGreeting, categoryGroups, effectiveLocale)

	if s.sessionManager != nil {
		s.sessionManager.RecordTurn(sessionID, cleanMsg, answer, plan, plan.Searches)
	}

	return &GuestChatResponseDTO{
		SessionID:      sessionID,
		GreetingPrefix: plan.SocialGreeting,
		Answer:         answer,
		CategoryGroups: categoryGroups,
		Listings:       flatListings,
		LatencyMS:      float64(time.Since(t0).Milliseconds()),
		Model:          "huntiq-assistant:conversational-v2",
	}, nil
}

func (s *AIChatService) synthesizeMultiCategoryGroundedResponse(
	ctx context.Context,
	userQuery string,
	greetingPrefix string,
	groups []MultiCategorySearchResultDTO,
	locale string,
) string {
	totalListings := 0
	for _, g := range groups {
		totalListings += len(g.Listings)
	}

	if totalListings == 0 {
		var noResultMsg string
		if len(groups) == 1 {
			displayTerm := groups[0].SearchTerm
			if displayTerm == "" {
				displayTerm = groups[0].CategoryName
			}
			localizedTerm := aiclient.GetLocalizedEntityTerm(displayTerm, locale)
			if localizedTerm != "" {
				displayTerm = localizedTerm
			}

			switch locale {
			case "am":
				noResultMsg = fmt.Sprintf("በHUNTIQ ዳታቤዝ ውስጥ ለ'%s' በአሁኑ ሰዓት የሚስማማ የተረጋገጠ ዝርዝር አልተገኘም። እባክዎ የፍለጋ ቃላትን በማስተካከል ወይም የዋጋ ገደቡን በማስፋት ይሞክሩ።", displayTerm)
			case "om":
				noResultMsg = fmt.Sprintf("Kuusaa gabaa HUNTIQ keessatti '%s' tiif yeroo ammaa tarreeffamni walsimu hin argamne. Maaloo jechoota biroo fayyadamuun ykn daangaa gatii jijjiiruun yaalaa.", displayTerm)
			case "ti":
				noResultMsg = fmt.Sprintf("ኣብ ዳታቤዝ HUNTIQ ን'%s' ኣብዚ እዋን ዝሰማማዕ ዝተረጋገፀ ዝርዝር ኣይተረኽበን። በጃኹም ናይ ፍለጋ ቃላት ወይ ደረት ዋጋ ብምስትኻኻል ፈትኑ።", displayTerm)
			default:
				noResultMsg = fmt.Sprintf("I searched our live marketplace database for '%s', but found no approved listings matching those exact criteria right now. You might try searching for broader keywords, adjusting your budget, or checking a different location.", displayTerm)
			}
		} else if len(groups) > 1 {
			var catNames []string
			for _, g := range groups {
				catNames = append(catNames, g.CategoryName)
			}
			joined := strings.Join(catNames, ", ")
			switch locale {
			case "am":
				noResultMsg = fmt.Sprintf("በHUNTIQ ዳታቤዝ ውስጥ ለ%s በአሁኑ ሰዓት የሚስማማ የተረጋገጠ ዝርዝር አልተገኘም። እባክዎ የፍለጋ ቃላትን በማስተካከል ይሞክሩ።", joined)
			case "om":
				noResultMsg = fmt.Sprintf("Kuusaa gabaa HUNTIQ keessatti %s tiif yeroo ammaa tarreeffamni walsimu hin argamne. Maaloo jechoota biroo fayyadamuun yaalaa.", joined)
			case "ti":
				noResultMsg = fmt.Sprintf("ኣብ ዳታቤዝ HUNTIQ ን%s ኣብዚ እዋን ዝሰማማዕ ዝተረጋገፀ ዝርዝር ኣይተረኽበን። በጃኹም ናይ ፍለጋ ቃላት ብምስትኻኻል ፈትኑ።", joined)
			default:
				noResultMsg = fmt.Sprintf("I searched our live marketplace database across %s, but found no approved listings matching those categories right now. Try broader search terms or check back soon as new verified items are posted.", joined)
			}
		} else {
			switch locale {
			case "am":
				noResultMsg = "በHUNTIQ ዳታቤዝ ውስጥ በአሁኑ ሰዓት የሚስማማ የተረጋገጠ ዝርዝር አልተገኘም። እባክዎ የፍለጋ ቃላትን በማስተካከል ይሞክሩ።"
			case "om":
				noResultMsg = "Kuusaa gabaa HUNTIQ keessatti yeroo ammaa tarreeffamni walsimu hin argamne. Maaloo jechoota biroo fayyadamuun yaalaa."
			case "ti":
				noResultMsg = "ኣብ ዳታቤዝ HUNTIQ ኣብዚ እዋን ዝሰማማዕ ዝተረጋገፀ ዝርዝር ኣይተረኽበን። በጃኹም ናይ ፍለጋ ቃላት ብምስትኻኻል ፈትኑ።"
			default:
				noResultMsg = "I searched our live marketplace database, but found no approved listings matching your search right now. Try searching with different keywords or exploring popular categories."
			}
		}

		if greetingPrefix != "" {
			return greetingPrefix + "\n\n" + noResultMsg
		}
		return noResultMsg
	}

	// Format results by category group
	var b strings.Builder
	if greetingPrefix != "" {
		b.WriteString(greetingPrefix)
		b.WriteString("\n\n")
	}

	if len(groups) == 1 {
		grp := groups[0]
		var lines []string
		for _, l := range grp.Listings {
			typeStr := "For Sale"
			if l.ListingType == "rent" {
				typeStr = "For Rent"
			}
			lines = append(lines, fmt.Sprintf("• **%s** — %.0f %s (%s, %s)", l.Title, l.Price, l.Currency, l.City, typeStr))
		}

		switch locale {
		case "am":
			b.WriteString(fmt.Sprintf("በHUNTIQ ዳታቤዝ ውስጥ ለ**%s** የተገኙ **%d የተረጋገጡ ዝርዝሮች**:\n\n%s\n\nሙሉ መረጃ ለማየትና ከሻጩ/አከራዩ ጋር ለመገናኘት ከስር ያሉትን ካርዶች መጫን ይችላሉ።", grp.CategoryName, len(grp.Listings), strings.Join(lines, "\n")))
		case "om":
			b.WriteString(fmt.Sprintf("Kuusaa gabaa HUNTIQ keessatti **%s** tiif **tarreeffamoota mirkanaa'an %d** argamaniiru:\n\n%s\n\nOdeeffannoo dabalataaf tarreeffama armaan gadii cuqaasaa.", grp.CategoryName, len(grp.Listings), strings.Join(lines, "\n")))
		case "ti":
			b.WriteString(fmt.Sprintf("ኣብ ዳታቤዝ HUNTIQ ን**%s** **%d ዝተረጋገፁ ዝርዝራት** ተረኺቦም ኣለዉ:\n\n%s\n\nዝርዝር ሓበሬታ ንምርኣይ ካብዞም ኣብ ታሕቲ ዘለዉ ካርዶች ምረጹ።", grp.CategoryName, len(grp.Listings), strings.Join(lines, "\n")))
		default:
			b.WriteString(fmt.Sprintf("Here are the **%d verified live marketplace listings** for **%s** in the HUNTIQ database:\n\n%s\n\nYou can click on any listing card below to view full specifications, photos, and contact the provider directly.", len(grp.Listings), grp.CategoryName, strings.Join(lines, "\n")))
		}
	} else {
		// Multi-category breakdown
		switch locale {
		case "am":
			b.WriteString(fmt.Sprintf("በHUNTIQ ዳታቤዝ ውስጥ በ**%d ክፍሎች** የተገኙ ዝርዝሮች:\n\n", len(groups)))
		case "om":
			b.WriteString(fmt.Sprintf("Gabaa HUNTIQ keessatti gareewwan **%d** keessaa tarreeffamoota argaman:\n\n", len(groups)))
		case "ti":
			b.WriteString(fmt.Sprintf("ኣብ HUNTIQ ብ**%d ምድባት** ዝተረኸቡ ዝርዝራት:\n\n", len(groups)))
		default:
			b.WriteString(fmt.Sprintf("I found verified marketplace listings across **%d categories** in our live database:\n\n", len(groups)))
		}

		for _, grp := range groups {
			b.WriteString(fmt.Sprintf("### %s (%d found)\n", grp.CategoryName, len(grp.Listings)))
			if len(grp.Listings) == 0 {
				switch locale {
				case "am":
					b.WriteString("በአሁኑ ሰዓት በዚህ ክፍል ውስጥ የተመዘገበ ዝርዝር የለም።\n\n")
				case "om":
					b.WriteString("Yeroo ammaa garee kana keessatti tarreeffamni hin jiru.\n\n")
				case "ti":
					b.WriteString("ኣብዚ እዋን ኣብዚ ምድብ ዝተመዝገበ ዝርዝር የለን።\n\n")
				default:
					b.WriteString("No active listings currently found in this category.\n\n")
				}
				continue
			}

			for _, l := range grp.Listings {
				typeStr := "For Sale"
				if l.ListingType == "rent" {
					typeStr = "For Rent"
				}
				b.WriteString(fmt.Sprintf("• **%s** — %.0f %s (%s, %s)\n", l.Title, l.Price, l.Currency, l.City, typeStr))
			}
			b.WriteString("\n")
		}

		switch locale {
		case "am":
			b.WriteString("ዝርዝር መረጃ ለማየት ከስር ያሉትን ካርዶች ይመልከቱ።")
		case "om":
			b.WriteString("Odeeffannoo dabalataaf tarreeffama armaan gadii ilaalaa.")
		case "ti":
			b.WriteString("ተወሳኺ ሓበሬታ ንምርኣይ ካብዞም ኣብ ታሕቲ ዘለዉ ካርዶች ተወከሱ።")
		default:
			b.WriteString("Explore the listing cards below to review full photos, prices, and seller contacts.")
		}
	}

	return b.String()
}

// synthesizeGroundedResponse provides backward compatibility for single-search queries.
func (s *AIChatService) synthesizeGroundedResponse(
	ctx context.Context,
	userQuery string,
	listings []ListingSummaryDTO,
	searchArgs map[string]interface{},
	locale string,
) string {
	searchTerm, _ := searchArgs["search_term"].(string)
	grp := MultiCategorySearchResultDTO{
		CategoryName: searchTerm,
		SearchTerm:   searchTerm,
		TotalCount:   len(listings),
		Listings:     listings,
	}
	return s.synthesizeMultiCategoryGroundedResponse(ctx, userQuery, "", []MultiCategorySearchResultDTO{grp}, locale)
}

func containsIntent(intents []IntentType, target IntentType) bool {
	for _, it := range intents {
		if it == target {
			return true
		}
	}
	return false
}

func generateSessionID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("guest_%x_%d", b, time.Now().Unix())
}

func buildDefaultConversationalResponse(locale string) string {
	switch locale {
	case "am":
		return "እንዴት ልርዳዎት እችላለሁ? በHUNTIQ ላይ ስልኮችን፣ ላፕቶፖችን፣ መኪናዎችን ወይም ቤቶችን መፈለግ ይችላሉ።"
	case "om":
		return "Akkamitti isin gargaaruu danda'a? HUNTIQ irratti bilbila, laaptooppii, konkolaataa ykn mana barbaaduu dandeessu."
	case "ti":
		return "ብኸመይ ክሕግዘኩም ይኽእል? ኣብ HUNTIQ ስልኪ፣ ላፕቶፕ፣ መኪና ወይ ገዛ ክትደልዩ ትኽእሉ ኢኹም።"
	default:
		return "How can I help you today? You can ask me to search for phones, laptops, cars, apartments for rent, or explain how HUNTIQ works."
	}
}

