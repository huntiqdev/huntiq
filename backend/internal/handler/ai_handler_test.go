package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/ai"
)

type mockAIClient struct{}

func (m *mockAIClient) Generate(ctx context.Context, systemInstruction string, history []ai.GeminiContent, prompt string) (string, int, error) {
	return "I can help you find products on HUNTIQ.", 20, nil
}

func (m *mockAIClient) HealthCheck(ctx context.Context) (*ai.AIHealthStatus, error) {
	return &ai.AIHealthStatus{Status: "operational", Enabled: true, Provider: "Mock"}, nil
}

func (m *mockAIClient) GetModel() string {
	return "mock-model"
}

type mockAIRepo struct {
	convs map[string]*ai.Conversation
	msgs  map[int64][]ai.Message
}

func (r *mockAIRepo) GetOrCreateConversation(ctx context.Context, sessionID string, userID *int, lang string) (*ai.Conversation, error) {
	conv := &ai.Conversation{ID: 1, SessionID: sessionID, Title: "Test Conv", LanguageCode: lang}
	r.convs[sessionID] = conv
	return conv, nil
}

func (r *mockAIRepo) GetConversation(ctx context.Context, sessionID string) (*ai.Conversation, error) {
	return r.convs[sessionID], nil
}

func (r *mockAIRepo) GetRecentMessages(ctx context.Context, convID int64, limit int) ([]ai.Message, error) {
	return r.msgs[convID], nil
}

func (r *mockAIRepo) SaveMessage(ctx context.Context, msg *ai.Message) error {
	r.msgs[msg.ConversationID] = append(r.msgs[msg.ConversationID], *msg)
	return nil
}

func (r *mockAIRepo) RecordAuditEvent(ctx context.Context, evt *ai.AuditEvent) error {
	return nil
}

func (r *mockAIRepo) ArchiveConversation(ctx context.Context, sessionID string, userID *int) error {
	if c, ok := r.convs[sessionID]; ok {
		c.IsArchived = true
	}
	return nil
}

func (r *mockAIRepo) GetAdminMetrics(ctx context.Context) (*ai.AIMetricsReport, error) {
	return &ai.AIMetricsReport{TotalQueries: 100, TotalTokens: 5000}, nil
}

func setupAITestApp() (*fiber.App, *AIHandler) {
	app := fiber.New()
	repo := &mockAIRepo{convs: make(map[string]*ai.Conversation), msgs: make(map[int64][]ai.Message)}
	svc := ai.NewAIService(&mockAIClient{}, repo, true, 30)
	handler := NewAIHandler(svc, repo)

	app.Get("/api/ai/health", handler.Health)
	app.Post("/api/ai/chat", handler.Chat)
	app.Get("/api/ai/session/:sessionId", handler.GetSessionHistory)
	app.Delete("/api/ai/session/:sessionId", handler.ClearSession)
	app.Get("/api/admin/ai/metrics", handler.AdminMetrics)

	return app, handler
}

func TestAIHandler_Health(t *testing.T) {
	app, _ := setupAITestApp()
	req := httptest.NewRequest(http.MethodGet, "/api/ai/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestAIHandler_Chat_Validation(t *testing.T) {
	app, _ := setupAITestApp()

	// Missing session_id
	body, _ := json.Marshal(ai.ChatRequest{Message: "hello"})
	req := httptest.NewRequest(http.MethodPost, "/api/ai/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", resp.StatusCode)
	}

	// Valid chat request
	body, _ = json.Marshal(ai.ChatRequest{SessionID: "sess-valid-1", Message: "Find houses for rent in Bole"})
	req = httptest.NewRequest(http.MethodPost, "/api/ai/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestAIHandler_SessionLifecycle(t *testing.T) {
	app, _ := setupAITestApp()

	// Chat to create session
	body, _ := json.Marshal(ai.ChatRequest{SessionID: "lifecycle-sess", Message: "Buy laptops"})
	req := httptest.NewRequest(http.MethodPost, "/api/ai/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	_, _ = app.Test(req)

	// Get session
	getReq := httptest.NewRequest(http.MethodGet, "/api/ai/session/lifecycle-sess", nil)
	resp, err := app.Test(getReq)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}

	// Delete session
	delReq := httptest.NewRequest(http.MethodDelete, "/api/ai/session/lifecycle-sess", nil)
	resp, err = app.Test(delReq)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestAIHandler_AdminMetrics(t *testing.T) {
	app, _ := setupAITestApp()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/ai/metrics", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}
