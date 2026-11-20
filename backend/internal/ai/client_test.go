package ai

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGeminiClient_Initialization(t *testing.T) {
	client := NewGeminiClient("secret-test-key", "gemini-1.5-flash", 10, 500)
	if client.apiKey != "secret-test-key" {
		t.Errorf("expected apiKey to match, got %s", client.apiKey)
	}
	if client.model != "gemini-1.5-flash" {
		t.Errorf("expected model gemini-1.5-flash, got %s", client.model)
	}
	if client.maxTokens != 500 {
		t.Errorf("expected maxTokens 500, got %d", client.maxTokens)
	}
}

func TestGeminiClient_SanitizeError(t *testing.T) {
	client := NewGeminiClient("SUPER_SECRET_KEY_12345", "gemini-1.5-flash", 10, 500)
	err := client.sanitizeError(http.ErrHandlerTimeout)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	rawMsg := "failed connection to endpoint with key SUPER_SECRET_KEY_12345"
	testLeak := client.sanitizeError(errors.New(rawMsg))
	if strings.Contains(testLeak.Error(), "SUPER_SECRET_KEY_12345") {
		t.Errorf("leak detected: %s", testLeak.Error())
	}
	if !strings.Contains(testLeak.Error(), "[REDACTED_API_KEY]") {
		t.Errorf("expected [REDACTED_API_KEY], got %s", testLeak.Error())
	}
}

func TestGeminiClient_GenerateSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify header security
		if r.Header.Get("x-goog-api-key") != "test-api-key" {
			http.Error(w, "missing or invalid key header", http.StatusForbidden)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "invalid content type", http.StatusBadRequest)
			return
		}

		resp := GeminiGenerateResponse{
			Candidates: []GeminiCandidate{
				{
					Content: GeminiContent{
						Role: "model",
						Parts: []GeminiPart{
							{Text: "I found 3 laptops available on HUNTIQ."},
						},
					},
					FinishReason: "STOP",
				},
			},
			UsageMetadata: GeminiUsageMetadata{
				PromptTokenCount:     15,
				CandidatesTokenCount: 10,
				TotalTokenCount:      25,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewGeminiClient("test-api-key", "gemini-1.5-flash", 5, 256).
		WithBaseURL(server.URL)

	reply, tokens, err := client.Generate(context.Background(), "You are HUNTIQ AI.", nil, "Find laptops")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != "I found 3 laptops available on HUNTIQ." {
		t.Errorf("unexpected reply: %s", reply)
	}
	if tokens != 25 {
		t.Errorf("expected 25 tokens, got %d", tokens)
	}
}

func TestGeminiClient_RetryOn429(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			http.Error(w, "rate limited", http.StatusTooManyRequests)
			return
		}
		resp := GeminiGenerateResponse{
			Candidates: []GeminiCandidate{
				{
					Content: GeminiContent{
						Role: "model",
						Parts: []GeminiPart{
							{Text: "Success after retry."},
						},
					},
				},
			},
			UsageMetadata: GeminiUsageMetadata{TotalTokenCount: 10},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewGeminiClient("test-api-key", "gemini-1.5-flash", 5, 256).
		WithBaseURL(server.URL)

	reply, _, err := client.Generate(context.Background(), "System prompt", nil, "hello")
	if err != nil {
		t.Fatalf("expected success after retry, got: %v", err)
	}
	if reply != "Success after retry." {
		t.Errorf("unexpected reply: %s", reply)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestGeminiClient_HealthCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := GeminiGenerateResponse{
			Candidates: []GeminiCandidate{
				{
					Content: GeminiContent{
						Role: "model",
						Parts: []GeminiPart{{Text: "pong"}},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewGeminiClient("test-key", "gemini-1.5-flash", 5, 256).WithBaseURL(server.URL)
	status, err := client.HealthCheck(context.Background())
	if err != nil {
		t.Fatalf("unexpected health check error: %v", err)
	}
	if status.Status != "operational" {
		t.Errorf("expected status operational, got %s", status.Status)
	}
	if !status.Enabled {
		t.Errorf("expected enabled true")
	}
}
