package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/internal/service/aiclient"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

// InternalAIToolHandler handles restricted service-to-service marketplace tool calls from Python AI.
type InternalAIToolHandler struct {
	toolDispatcher *service.AIToolDispatcher
	secretKey      string
}

// NewInternalAIToolHandler creates a new instance for internal AI tool execution.
func NewInternalAIToolHandler(toolDispatcher *service.AIToolDispatcher, secretKey string) *InternalAIToolHandler {
	return &InternalAIToolHandler{
		toolDispatcher: toolDispatcher,
		secretKey:      secretKey,
	}
}

// ToolExecuteRequest represents the untrusted incoming tool payload.
type ToolExecuteRequest struct {
	ToolName  string                 `json:"tool_name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// Execute handles POST /api/internal/ai/tools/execute
func (h *InternalAIToolHandler) Execute(c *fiber.Ctx) error {
	// 1. Authenticate service-to-service header
	providedSecret := strings.TrimSpace(c.Get("X-HUNTIQ-AI-SECRET"))
	if h.secretKey != "" && (providedSecret == "" || providedSecret != h.secretKey) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized service call: invalid or missing X-HUNTIQ-AI-SECRET",
		})
	}

	var req ToolExecuteRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	toolName := strings.TrimSpace(req.ToolName)
	if toolName == "" {
		return response.BadRequest(c, "tool_name is required")
	}

	// 2. Layer 2 Strict Tool Validation
	switch toolName {
	case "search_marketplace_listings", "get_listing_details", "compare_listings", "get_marketplace_stats":
		// Allowed tool names
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "unsupported or disallowed tool: " + toolName,
		})
	}

	if req.Arguments == nil {
		req.Arguments = make(map[string]interface{})
	}

	// Restrict result limit to prevent memory exhaustion
	if lim, ok := req.Arguments["limit"].(float64); ok {
		if lim < 1 {
			req.Arguments["limit"] = 1
		} else if lim > 10 {
			req.Arguments["limit"] = 10
		}
	} else {
		req.Arguments["limit"] = 5
	}

	// Validate min/max price bounds
	if minP, ok := req.Arguments["min_price"].(float64); ok && minP < 0 {
		return response.BadRequest(c, "min_price cannot be negative")
	}
	if maxP, ok := req.Arguments["max_price"].(float64); ok && maxP < 0 {
		return response.BadRequest(c, "max_price cannot be negative")
	}

	// 3. Execute authorized marketplace query via dispatcher
	call := aiclient.ChatToolCall{
		ID:   "internal_tool_call",
		Type: "function",
		Function: aiclient.ChatFunctionCall{
			Name:      toolName,
			Arguments: req.Arguments,
		},
	}

	rawResult, listings, err := h.toolDispatcher.ExecuteToolCall(c.Context(), call)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if listings == nil {
		listings = []service.ListingSummaryDTO{}
	}

	return c.JSON(fiber.Map{
		"success":        true,
		"tool_name":      toolName,
		"raw_result":     rawResult,
		"total_matching": len(listings),
		"listings":       listings,
	})
}
