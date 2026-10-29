package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

// AIChatHandler handles HTTP routing and serialization for conversational AI endpoints.
type AIChatHandler struct {
	chatSvc *service.AIChatService
}

// NewAIChatHandler creates a new handler instance.
func NewAIChatHandler(chatSvc *service.AIChatService) *AIChatHandler {
	return &AIChatHandler{chatSvc: chatSvc}
}

// CreateConversationRequest represents payload for creating a chat session.
type CreateConversationRequest struct {
	Title string `json:"title"`
}

// CreateConversation handles POST /api/ai/conversations
func (h *AIChatHandler) CreateConversation(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	var req CreateConversationRequest
	_ = c.BodyParser(&req)

	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = "New Chat"
	}

	conv, err := h.chatSvc.CreateConversation(c.Context(), userID, title)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.Created(c, conv)
}

// ListConversations handles GET /api/ai/conversations
func (h *AIChatHandler) ListConversations(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	page, limit := parsePagination(c, 30, 50)
	offset := (page - 1) * limit

	conversations, err := h.chatSvc.ListConversations(c.Context(), userID, limit, offset)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OKList(c, conversations)
}

// GetConversation handles GET /api/ai/conversations/:id
func (h *AIChatHandler) GetConversation(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	conversationID := c.Params("id")
	if conversationID == "" {
		return response.BadRequest(c, "missing conversation id")
	}

	conv, err := h.chatSvc.GetConversation(c.Context(), conversationID, userID)
	if err != nil {
		return response.NotFound(c, "conversation not found or unauthorized")
	}

	return response.OK(c, conv)
}

// UpdateConversationRequest represents payload for renaming a chat session.
type UpdateConversationRequest struct {
	Title string `json:"title"`
}

// UpdateConversation handles PATCH /api/ai/conversations/:id
func (h *AIChatHandler) UpdateConversation(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	conversationID := c.Params("id")
	if conversationID == "" {
		return response.BadRequest(c, "missing conversation id")
	}

	var req UpdateConversationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return response.BadRequest(c, "title cannot be empty")
	}

	if err := h.chatSvc.UpdateConversationTitle(c.Context(), conversationID, userID, title); err != nil {
		return response.NotFound(c, "conversation not found or unauthorized")
	}

	return response.OK(c, fiber.Map{"updated": true, "title": title})
}

// DeleteConversation handles DELETE /api/ai/conversations/:id
func (h *AIChatHandler) DeleteConversation(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	conversationID := c.Params("id")
	if conversationID == "" {
		return response.BadRequest(c, "missing conversation id")
	}

	if err := h.chatSvc.DeleteConversation(c.Context(), conversationID, userID); err != nil {
		return response.NotFound(c, "conversation not found or unauthorized")
	}

	return response.OK(c, fiber.Map{"deleted": true})
}

// SendMessageRequest represents incoming user message payload.
type SendMessageRequest struct {
	Content string `json:"content"`
	Locale  string `json:"locale"`
}

// SendMessage handles POST /api/ai/conversations/:id/messages
func (h *AIChatHandler) SendMessage(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	conversationID := c.Params("id")
	if conversationID == "" {
		return response.BadRequest(c, "missing conversation id")
	}

	var req SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	cleanContent := strings.TrimSpace(req.Content)
	if cleanContent == "" {
		return response.BadRequest(c, "message content cannot be empty")
	}
	if len(cleanContent) > 4000 {
		return response.BadRequest(c, "message exceeds maximum length of 4000 characters")
	}

	locale := strings.TrimSpace(req.Locale)
	if locale == "" {
		locale = "en"
	}

	assistantMsg, err := h.chatSvc.SendMessage(c.Context(), userID, conversationID, cleanContent, locale)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, assistantMsg)
}

// GuestChatRequest represents incoming payload for unauthenticated public AI chat.
type GuestChatRequest struct {
	Message string                    `json:"message"`
	Locale  string                    `json:"locale"`
	History []service.GuestChatMessage `json:"history"`
}

// GuestChat handles POST /api/ai/chat (Public, unauthenticated, stateless)
func (h *AIChatHandler) GuestChat(c *fiber.Ctx) error {
	var req GuestChatRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	cleanMessage := strings.TrimSpace(req.Message)
	if cleanMessage == "" {
		return response.BadRequest(c, "message cannot be empty")
	}
	if len(cleanMessage) > 1000 {
		return response.BadRequest(c, "message exceeds maximum length of 1000 characters")
	}

	locale := strings.TrimSpace(req.Locale)
	if locale == "" {
		locale = "en"
	}

	// Bound history to at most 4 items to keep context small and fast
	history := req.History
	if len(history) > 4 {
		history = history[len(history)-4:]
	}

	res, err := h.chatSvc.ProcessGuestChat(c.Context(), service.GuestChatRequestDTO{
		Message: cleanMessage,
		Locale:  locale,
		History: history,
	})
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, res)
}

