package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	handlers "github.com/your-org/used-tools-marketplace/internal/handler"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

func setupTestAppWithGuestChat() *fiber.App {
	app := fiber.New()
	chatSvc := service.NewAIChatService(nil, nil, nil)
	aiHandler := handlers.NewAIChatHandler(chatSvc)

	app.Post("/api/ai/chat", aiHandler.GuestChat)
	return app
}

func TestAIChatHandler_GuestChat_OutOfDomain(t *testing.T) {
	app := setupTestAppWithGuestChat()

	body, _ := json.Marshal(map[string]any{
		"message": "write a python function to reverse a string",
		"locale":  "en",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/ai/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var jsonResp map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&jsonResp)

	data, ok := jsonResp["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object in response, got %v", jsonResp)
	}

	answer, _ := data["answer"].(string)
	if !strings.Contains(answer, "marketplace assistant") {
		t.Errorf("expected marketplace redirection, got: %s", answer)
	}
}

func TestAIChatHandler_GuestChat_PlatformFAQ(t *testing.T) {
	app := setupTestAppWithGuestChat()

	body, _ := json.Marshal(map[string]any{
		"message": "how do I post an apartment?",
		"locale":  "en",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/ai/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var jsonResp map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&jsonResp)

	data, _ := jsonResp["data"].(map[string]any)
	answer, _ := data["answer"].(string)
	if !strings.Contains(answer, "Post Listing") {
		t.Errorf("expected posting guide in answer, got: %s", answer)
	}
}

func TestAIChatHandler_GuestChat_Validation(t *testing.T) {
	app := setupTestAppWithGuestChat()

	// Empty message
	body, _ := json.Marshal(map[string]any{"message": ""})
	req := httptest.NewRequest(http.MethodPost, "/api/ai/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for empty message, got %d", resp.StatusCode)
	}

	// Oversized message
	longMsg := strings.Repeat("M", 1005)
	bodyLong, _ := json.Marshal(map[string]any{"message": longMsg})
	reqLong := httptest.NewRequest(http.MethodPost, "/api/ai/chat", bytes.NewReader(bodyLong))
	reqLong.Header.Set("Content-Type", "application/json")
	respLong, _ := app.Test(reqLong)
	if respLong.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for oversized message, got %d", respLong.StatusCode)
	}
}
