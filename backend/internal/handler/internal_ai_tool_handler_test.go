package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

func TestInternalAIToolHandler_Authentication(t *testing.T) {
	app := fiber.New()
	secret := "test-secret-key-123"
	dispatcher := service.NewAIToolDispatcher(nil)
	handler := NewInternalAIToolHandler(dispatcher, secret)

	app.Post("/api/internal/ai/tools/execute", handler.Execute)

	// Case 1: Missing secret header -> 401
	body, _ := json.Marshal(map[string]interface{}{
		"tool_name": "search_marketplace_listings",
	})
	req := httptest.NewRequest("POST", "/api/internal/ai/tools/execute", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected status 401 for missing secret, got %d", resp.StatusCode)
	}

	// Case 2: Wrong secret header -> 401
	req2 := httptest.NewRequest("POST", "/api/internal/ai/tools/execute", bytes.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-HUNTIQ-AI-SECRET", "wrong-secret")
	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp2.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected status 401 for wrong secret, got %d", resp2.StatusCode)
	}
}

func TestInternalAIToolHandler_DisallowedTool(t *testing.T) {
	app := fiber.New()
	secret := "test-secret"
	dispatcher := service.NewAIToolDispatcher(nil)
	handler := NewInternalAIToolHandler(dispatcher, secret)

	app.Post("/api/internal/ai/tools/execute", handler.Execute)

	body, _ := json.Marshal(map[string]interface{}{
		"tool_name": "drop_database_table",
	})
	req := httptest.NewRequest("POST", "/api/internal/ai/tools/execute", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-HUNTIQ-AI-SECRET", secret)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status 400 for disallowed tool, got %d", resp.StatusCode)
	}
}

func TestInternalAIToolHandler_ValidToolExecution(t *testing.T) {
	app := fiber.New()
	secret := "test-secret"
	dispatcher := service.NewAIToolDispatcher(nil) // productRepo is nil, handled defensively
	handler := NewInternalAIToolHandler(dispatcher, secret)

	app.Post("/api/internal/ai/tools/execute", handler.Execute)

	body, _ := json.Marshal(map[string]interface{}{
		"tool_name": "search_marketplace_listings",
		"arguments": map[string]interface{}{
			"search_term": "phone",
			"limit":       5,
		},
	})
	req := httptest.NewRequest("POST", "/api/internal/ai/tools/execute", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-HUNTIQ-AI-SECRET", secret)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status 200 for valid tool execution, got %d", resp.StatusCode)
	}
}
