package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Gemini API structures based on official Google Generative Language v1beta specification.
type GeminiPart struct {
	Text string `json:"text,omitempty"`
}

type GeminiContent struct {
	Role  string       `json:"role,omitempty"` // "user", "model", "system"
	Parts []GeminiPart `json:"parts"`
}

type GeminiGenerationConfig struct {
	Temperature     float64 `json:"temperature"`
	TopP            float64 `json:"topP,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

type GeminiGenerateRequest struct {
	Contents          []GeminiContent         `json:"contents"`
	SystemInstruction *GeminiContent          `json:"systemInstruction,omitempty"`
	GenerationConfig  *GeminiGenerationConfig `json:"generationConfig,omitempty"`
}

type GeminiCandidate struct {
	Content      GeminiContent `json:"content"`
	FinishReason string        `json:"finishReason"`
	Index        int           `json:"index"`
}

type GeminiUsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

type GeminiGenerateResponse struct {
	Candidates    []GeminiCandidate   `json:"candidates"`
	UsageMetadata GeminiUsageMetadata `json:"usageMetadata"`
	Error         *GeminiAPIError     `json:"error,omitempty"`
}

type GeminiAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// ProviderClient interface defines the contract for AI provider communication.
type ProviderClient interface {
	Generate(ctx context.Context, systemInstruction string, history []GeminiContent, prompt string) (string, int, error)
	HealthCheck(ctx context.Context) (*AIHealthStatus, error)
	GetModel() string
}

// GeminiClient implements ProviderClient using Google's Generative Language REST API.
type GeminiClient struct {
	apiKey     string
	model      string
	timeout    time.Duration
	maxTokens  int
	httpClient *http.Client
	baseURL    string
}

// NewGeminiClient instantiates a new Gemini client with validated parameters.
func NewGeminiClient(apiKey, model string, timeoutSeconds, maxTokens int) *GeminiClient {
	if model == "" {
		model = "gemini-1.5-flash"
	}
	if timeoutSeconds <= 0 {
		timeoutSeconds = 15
	}
	if maxTokens <= 0 {
		maxTokens = 1024
	}

	return &GeminiClient{
		apiKey:    apiKey,
		model:     model,
		timeout:   time.Duration(timeoutSeconds) * time.Second,
		maxTokens: maxTokens,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds+5) * time.Second,
		},
		baseURL: "https://generativelanguage.googleapis.com/v1beta",
	}
}

// WithBaseURL overrides the base URL for testing purposes.
func (c *GeminiClient) WithBaseURL(url string) *GeminiClient {
	c.baseURL = url
	return c
}

// WithHTTPClient overrides the HTTP client for testing purposes.
func (c *GeminiClient) WithHTTPClient(client *http.Client) *GeminiClient {
	c.httpClient = client
	return c
}

func (c *GeminiClient) GetModel() string {
	return c.model
}

// sanitizeError removes any accidental occurrences of API key from errors.
func (c *GeminiClient) sanitizeError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if c.apiKey != "" && strings.Contains(msg, c.apiKey) {
		msg = strings.ReplaceAll(msg, c.apiKey, "[REDACTED_API_KEY]")
	}
	return errors.New(msg)
}

// Generate sends a conversational turn to Google Gemini with bounded retry on transient errors.
func (c *GeminiClient) Generate(ctx context.Context, systemInstruction string, history []GeminiContent, prompt string) (string, int, error) {
	if c.apiKey == "" {
		return "", 0, errors.New("gemini api key is not configured")
	}

	contents := make([]GeminiContent, 0, len(history)+1)
	contents = append(contents, history...)
	contents = append(contents, GeminiContent{
		Role:  "user",
		Parts: []GeminiPart{{Text: prompt}},
	})

	reqBody := GeminiGenerateRequest{
		Contents: contents,
		GenerationConfig: &GeminiGenerationConfig{
			Temperature:     0.2, // Deterministic, domain-focused responses
			TopP:            0.95,
			MaxOutputTokens: c.maxTokens,
		},
	}

	if systemInstruction != "" {
		reqBody.SystemInstruction = &GeminiContent{
			Parts: []GeminiPart{{Text: systemInstruction}},
		}
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, fmt.Errorf("marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/models/%s:generateContent", c.baseURL, c.model)

	// Bounded retry policy: up to 2 retries for transient 429/503
	var lastErr error
	backoffs := []time.Duration{500 * time.Millisecond, 1500 * time.Millisecond}

	for attempt := 0; attempt <= len(backoffs); attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return "", 0, ctx.Err()
			case <-time.After(backoffs[attempt-1]):
			}
		}

		reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
		httpReq, err := http.NewRequestWithContext(reqCtx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
		if err != nil {
			cancel()
			return "", 0, c.sanitizeError(fmt.Errorf("create http request: %w", err))
		}

		httpReq.Header.Set("Content-Type", "application/json")
		// Strictly pass API key in header, NEVER as a URL query param
		httpReq.Header.Set("x-goog-api-key", c.apiKey)

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			cancel()
			lastErr = err
			continue
		}

		respBody, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		cancel()

		if readErr != nil {
			lastErr = readErr
			continue
		}

		// Handle Transient Status Codes
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			lastErr = fmt.Errorf("provider status %d: rate limit or service unavailable", resp.StatusCode)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			var apiErr GeminiGenerateResponse
			_ = json.Unmarshal(respBody, &apiErr)
			if apiErr.Error != nil {
				return "", 0, c.sanitizeError(fmt.Errorf("gemini api error (%d): %s", apiErr.Error.Code, apiErr.Error.Message))
			}
			return "", 0, c.sanitizeError(fmt.Errorf("gemini api returned status %d", resp.StatusCode))
		}

		var genResp GeminiGenerateResponse
		if err := json.Unmarshal(respBody, &genResp); err != nil {
			return "", 0, fmt.Errorf("decode gemini response: %w", err)
		}

		if len(genResp.Candidates) == 0 || len(genResp.Candidates[0].Content.Parts) == 0 {
			return "", genResp.UsageMetadata.TotalTokenCount, errors.New("empty response from gemini model")
		}

		replyText := strings.TrimSpace(genResp.Candidates[0].Content.Parts[0].Text)
		return replyText, genResp.UsageMetadata.TotalTokenCount, nil
	}

	return "", 0, c.sanitizeError(fmt.Errorf("gemini request failed after retries: %w", lastErr))
}

// HealthCheck validates connectivity to the Gemini API endpoint.
func (c *GeminiClient) HealthCheck(ctx context.Context) (*AIHealthStatus, error) {
	start := time.Now()
	status := &AIHealthStatus{
		Status:             "operational",
		Provider:           "Google Gemini",
		Model:              c.model,
		Enabled:            c.apiKey != "",
		LastChecked:        time.Now(),
		SupportedLanguages: []string{LangEN, LangAM, LangTI, LangOR},
	}

	if c.apiKey == "" {
		status.Status = "degraded"
		return status, errors.New("gemini api key is not configured")
	}

	// Ping using a minimal 1-token test prompt with 3s timeout
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, _, err := c.Generate(pingCtx, "System ping.", nil, "ping")
	status.LatencyMs = time.Since(start).Milliseconds()

	if err != nil {
		status.Status = "degraded"
		log.Printf("[HUNTIQ AI] Health check ping warning: %v", err)
		return status, err
	}

	return status, nil
}
