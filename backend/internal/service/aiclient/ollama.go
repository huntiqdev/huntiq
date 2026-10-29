package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// OllamaTag represents a model entry returned by GET /api/tags.
type OllamaTag struct {
	Name       string `json:"name"`
	Model      string `json:"model"`
	ModifiedAt string `json:"modified_at"`
	Size       int64  `json:"size"`
	Digest     string `json:"digest"`
}

// OllamaTagsResponse represents the response from GET /api/tags.
type OllamaTagsResponse struct {
	Models []OllamaTag `json:"models"`
}

// OllamaChatMessage represents a turn in the Ollama /api/chat payload.
type OllamaChatMessage struct {
	Role      string         `json:"role"`
	Content   string         `json:"content"`
	ToolCalls []ChatToolCall `json:"tool_calls,omitempty"`
}

// OllamaChatRequest represents the request body sent to POST /api/chat.
type OllamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []OllamaChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
	Tools    []any               `json:"tools,omitempty"`
	Options  map[string]any      `json:"options,omitempty"`
}

// OllamaChatResponse represents the response returned by POST /api/chat.
type OllamaChatResponse struct {
	Model         string            `json:"model"`
	CreatedAt     string            `json:"created_at"`
	Message       OllamaChatMessage `json:"message"`
	Done          bool              `json:"done"`
	TotalDuration int64             `json:"total_duration"`
}

// OllamaProvider communicates directly with a local or hosted Ollama instance.
type OllamaProvider struct {
	baseURL     string
	targetModel string
	client      *http.Client
}

// NewOllamaProvider creates a new configured Ollama provider.
func NewOllamaProvider(baseURL, targetModel string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:11434"
	}
	if targetModel == "" {
		targetModel = "huntiq-assistant:latest"
	}
	return &OllamaProvider{
		baseURL:     strings.TrimRight(baseURL, "/"),
		targetModel: targetModel,
		client: &http.Client{
			Timeout: 90 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        50,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// ListInstalledModels queries GET /api/tags to discover models installed in the local Ollama daemon.
func (p *OllamaProvider) ListInstalledModels(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL+"/api/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama tags request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama daemon unreachable at %s: %w", p.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama /api/tags returned status %d", resp.StatusCode)
	}

	var tagsResp OllamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
		return nil, fmt.Errorf("failed to decode ollama tags: %w", err)
	}

	names := make([]string, 0, len(tagsResp.Models))
	for _, m := range tagsResp.Models {
		names = append(names, m.Name)
	}
	return names, nil
}

// ResolveActiveModel verifies if targetModel is installed; if not, selects an available model.
func (p *OllamaProvider) ResolveActiveModel(ctx context.Context) string {
	installed, err := p.ListInstalledModels(ctx)
	if err != nil || len(installed) == 0 {
		return p.targetModel
	}

	// 1. Exact match on targetModel (e.g. "huntiq-assistant:latest" or "qwen3:4b")
	for _, name := range installed {
		if strings.EqualFold(name, p.targetModel) || strings.HasPrefix(strings.ToLower(name), strings.ToLower(p.targetModel)) {
			return name
		}
	}

	// 2. Check for huntiq-assistant
	for _, name := range installed {
		if strings.HasPrefix(strings.ToLower(name), "huntiq-assistant") {
			return name
		}
	}

	// 3. Check for qwen3:4b
	for _, name := range installed {
		if strings.HasPrefix(strings.ToLower(name), "qwen3:4b") || strings.HasPrefix(strings.ToLower(name), "qwen3") {
			return name
		}
	}

	// 4. Fall back to first installed model
	return installed[0]
}

// SendChat executes a conversational turn with Ollama /api/chat.
func (p *OllamaProvider) SendChat(ctx context.Context, messages []ChatLLMMessage, tools []any) (*ChatMessageResponse, error) {
	activeModel := p.ResolveActiveModel(ctx)

	ollamaMessages := make([]OllamaChatMessage, 0, len(messages))
	for _, m := range messages {
		ollamaMessages = append(ollamaMessages, OllamaChatMessage{
			Role:      m.Role,
			Content:   m.Content,
			ToolCalls: m.ToolCalls,
		})
	}

	payload := OllamaChatRequest{
		Model:    activeModel,
		Messages: ollamaMessages,
		Stream:   false,
		Tools:    tools,
		Options: map[string]any{
			"temperature": 0.7,
			"num_ctx":     4096,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ollama chat payload: %w", err)
	}

	t0 := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/api/chat", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama chat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama chat request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama chat returned status %d", resp.StatusCode)
	}

	var chatResp OllamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode ollama chat response: %w", err)
	}

	latency := float64(time.Since(t0).Milliseconds())
	content := cleanThinkingTags(chatResp.Message.Content)

	return &ChatMessageResponse{
		Content:      &content,
		ToolCalls:    chatResp.Message.ToolCalls,
		FinishReason: "stop",
		Model:        "ollama:" + chatResp.Model,
		LatencyMS:    latency,
	}, nil
}

// cleanThinkingTags removes internal <think>...</think> reasoning blocks from the assistant's message.
func cleanThinkingTags(raw string) string {
	if idx := strings.Index(raw, "</think>"); idx != -1 {
		cleaned := strings.TrimSpace(raw[idx+len("</think>"):])
		if cleaned != "" {
			return cleaned
		}
	}
	return raw
}
