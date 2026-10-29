package aiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOllamaProvider_ListInstalledModels_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			http.NotFound(w, r)
			return
		}
		resp := OllamaTagsResponse{
			Models: []OllamaTag{
				{Name: "huntiq-assistant:latest", Size: 2500000000},
				{Name: "qwen3:4b", Size: 2500000000},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	provider := NewOllamaProvider(server.URL, "huntiq-assistant:latest")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	models, err := provider.ListInstalledModels(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(models))
	}
	if models[0] != "huntiq-assistant:latest" || models[1] != "qwen3:4b" {
		t.Errorf("unexpected models: %v", models)
	}
}

func TestOllamaProvider_ResolveActiveModel_Hierarchy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := OllamaTagsResponse{
			Models: []OllamaTag{
				{Name: "qwen3:4b"},
				{Name: "mistral:latest"},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 1. Target is huntiq-assistant, but only qwen3:4b is present -> resolves qwen3:4b
	p1 := NewOllamaProvider(server.URL, "huntiq-assistant:latest")
	resolved := p1.ResolveActiveModel(context.Background())
	if resolved != "qwen3:4b" {
		t.Errorf("expected qwen3:4b, got %s", resolved)
	}

	// 2. Exact match targetModel
	p2 := NewOllamaProvider(server.URL, "qwen3:4b")
	resolved2 := p2.ResolveActiveModel(context.Background())
	if resolved2 != "qwen3:4b" {
		t.Errorf("expected qwen3:4b, got %s", resolved2)
	}
}

func TestOllamaProvider_SendChat_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/tags" {
			resp := OllamaTagsResponse{
				Models: []OllamaTag{{Name: "huntiq-assistant:latest"}},
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		if r.URL.Path == "/api/chat" {
			var req OllamaChatRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if req.Model != "huntiq-assistant:latest" {
				t.Errorf("expected model huntiq-assistant:latest, got %s", req.Model)
			}
			chatResp := OllamaChatResponse{
				Model: req.Model,
				Message: OllamaChatMessage{
					Role:    "assistant",
					Content: "Selam! I am HUNTIQ assistant.",
				},
				Done: true,
			}
			_ = json.NewEncoder(w).Encode(chatResp)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	provider := NewOllamaProvider(server.URL, "huntiq-assistant:latest")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	messages := []ChatLLMMessage{
		{Role: "user", Content: "Hello"},
	}
	resp, err := provider.SendChat(ctx, messages, nil)
	if err != nil {
		t.Fatalf("unexpected SendChat error: %v", err)
	}
	if resp == nil || resp.Content == nil {
		t.Fatal("expected non-nil response and content")
	}
	if *resp.Content != "Selam! I am HUNTIQ assistant." {
		t.Errorf("unexpected content: %s", *resp.Content)
	}
}

func TestOllamaProvider_UnreachableGraceful(t *testing.T) {
	provider := NewOllamaProvider("http://127.0.0.1:54321", "huntiq-assistant:latest")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_, err := provider.ListInstalledModels(ctx)
	if err == nil {
		t.Fatal("expected error for unreachable address, got nil")
	}

	// ResolveActiveModel should fallback without panic
	active := provider.ResolveActiveModel(ctx)
	if active != "huntiq-assistant:latest" {
		t.Errorf("expected fallback to targetModel, got %s", active)
	}
}

func TestCleanThinkingTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "contains thinking block",
			input:    "<think>\nAnalyzing query...\nUser wants a drill.\n</think>\nYes, we have drills available.",
			expected: "Yes, we have drills available.",
		},
		{
			name:     "no thinking block",
			input:    "Hello! How can I help you?",
			expected: "Hello! How can I help you?",
		},
		{
			name:     "only thinking block and nothing else",
			input:    "<think>Thinking...</think>",
			expected: "<think>Thinking...</think>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanThinkingTags(tt.input)
			if got != tt.expected {
				t.Errorf("cleanThinkingTags(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
