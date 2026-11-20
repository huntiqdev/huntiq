package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/ai"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type AIHandler struct {
	svc  *ai.AIService
	repo ai.AIRepository
}

func NewAIHandler(svc *ai.AIService, repo ai.AIRepository) *AIHandler {
	return &AIHandler{
		svc:  svc,
		repo: repo,
	}
}

// POST /api/ai/chat
func (h *AIHandler) Chat(c *fiber.Ctx) error {
	var req ai.ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		return response.BadRequest(c, "Message cannot be empty")
	}
	if len(req.Message) > 2000 {
		return response.BadRequest(c, "Message exceeds maximum allowed length of 2,000 characters")
	}

	req.SessionID = strings.TrimSpace(req.SessionID)
	if req.SessionID == "" {
		return response.BadRequest(c, "session_id is required")
	}

	var userID *int
	if uid, ok := c.Locals("userID").(int); ok && uid > 0 {
		userID = &uid
	}

	res, err := h.svc.ProcessChat(c.Context(), req, userID, c.IP())
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, res)
}

// GET /api/ai/health
func (h *AIHandler) Health(c *fiber.Ctx) error {
	status, _ := h.svc.HealthCheck(c.Context())
	return response.OK(c, status)
}

// GET /api/ai/session/:sessionId
func (h *AIHandler) GetSessionHistory(c *fiber.Ctx) error {
	sessionID := strings.TrimSpace(c.Params("sessionId"))
	if sessionID == "" {
		return response.BadRequest(c, "sessionId is required")
	}

	conv, err := h.repo.GetConversation(c.Context(), sessionID)
	if err != nil || conv == nil {
		return response.NotFound(c, "Conversation session not found")
	}

	// If session belongs to a user, ensure only that user or an admin can view it
	var currentUserID *int
	if uid, ok := c.Locals("userID").(int); ok && uid > 0 {
		currentUserID = &uid
	}
	role, _ := c.Locals("role").(string)

	if conv.UserID != nil && role != "admin" {
		if currentUserID == nil || *conv.UserID != *currentUserID {
			return response.Forbidden(c, "Unauthorized to access this session")
		}
	}

	msgs, err := h.repo.GetRecentMessages(c.Context(), conv.ID, 30)
	if err != nil {
		return response.InternalError(c, err)
	}
	if msgs == nil {
		msgs = []ai.Message{}
	}

	return response.OK(c, fiber.Map{
		"conversation": conv,
		"messages":     msgs,
	})
}

// DELETE /api/ai/session/:sessionId
func (h *AIHandler) ClearSession(c *fiber.Ctx) error {
	sessionID := strings.TrimSpace(c.Params("sessionId"))
	if sessionID == "" {
		return response.BadRequest(c, "sessionId is required")
	}

	var currentUserID *int
	if uid, ok := c.Locals("userID").(int); ok && uid > 0 {
		currentUserID = &uid
	}

	if err := h.repo.ArchiveConversation(c.Context(), sessionID, currentUserID); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"success": true, "message": "Session conversation archived successfully"})
}

// GET /api/admin/ai/metrics
func (h *AIHandler) AdminMetrics(c *fiber.Ctx) error {
	report, err := h.svc.GetAdminMetrics(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, report)
}
