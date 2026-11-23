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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":    false,
			"error_code": "ERR_INVALID_PAYLOAD",
			"error":      "invalid JSON payload",
		})
	}

	trimmed := strings.TrimSpace(req.Message)
	if trimmed == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":    false,
			"error_code": "ERR_EMPTY_MESSAGE",
			"error":      "message cannot be empty",
		})
	}

	// Truncate overly long queries (safety cap at 2000 chars)
	req.Message = truncate(trimmed, 2000)

	resp, err := h.svc.ProcessQuestion(c.Context(), req)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
			return c.Status(fiber.StatusGatewayTimeout).JSON(fiber.Map{
				"success":    false,
				"error_code": "ERR_TIMEOUT",
				"error":      "AI service request timed out",
			})
		}
		if strings.Contains(errStr, "python ai engine is not running") || strings.Contains(errStr, "error contacting python ai engine") {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success":    false,
				"error_code": "ERR_AI_SERVICE_UNAVAILABLE",
				"error":      "HUNTIQ AI standalone engine is currently unavailable",
			})
		}
		if strings.Contains(errStr, "Gemini") || strings.Contains(errStr, "model provider") {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"success":    false,
				"error_code": "ERR_PROVIDER_UNAVAILABLE",
				"error":      "External AI model provider is temporarily unavailable",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success":    false,
			"error_code": "ERR_INTERNAL",
			"error":      errStr,
		})
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
