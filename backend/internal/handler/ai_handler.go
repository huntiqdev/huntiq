package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/ai"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type AIHandler struct {
	svc *ai.AIService
}

func NewAIHandler(svc *ai.AIService) *AIHandler {
	return &AIHandler{svc: svc}
}

// POST /api/ai/chat
func (h *AIHandler) Chat(c *fiber.Ctx) error {
	var req ai.ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid JSON payload")
	}

	trimmed := strings.TrimSpace(req.Message)
	if trimmed == "" {
		return response.BadRequest(c, "message cannot be empty")
	}

	// Truncate overly long queries (safety cap at 2000 chars)
	req.Message = truncate(trimmed, 2000)

	resp, err := h.svc.ProcessQuestion(c.Context(), req)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, resp)
}

// GET /api/ai/health
func (h *AIHandler) Health(c *fiber.Ctx) error {
	if err := h.svc.Health(c.Context()); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success": false,
			"status":  "degraded",
			"error":   err.Error(),
		})
	}

	return response.OK(c, fiber.Map{
		"status":  "operational",
		"service": "HUNTIQ AI Standalone Engine",
		"version": "1.0",
	})
}
