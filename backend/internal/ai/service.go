package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/repository"
)

type AIService struct {
	supervisor  *PythonSupervisor
	productTool *ProductSearchTool
	trendTool   *TrendTool
}

func NewAIService(supervisor *PythonSupervisor, productRepo *repository.ProductRepo, searchRepo *repository.SearchRepo) *AIService {
	return &AIService{
		supervisor:  supervisor,
		productTool: NewProductSearchTool(productRepo),
		trendTool:   NewTrendTool(productRepo, searchRepo),
	}
}

func (s *AIService) ProcessQuestion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	trimmed := strings.TrimSpace(req.Message)
	if trimmed == "" {
		return nil, fmt.Errorf("question cannot be empty")
	}

	t0 := time.Now()

	// 1. Single Lightweight Routing & Pre-Gemini Domain Validation via Python Engine
	classifyRes, err := s.supervisor.Classify(ctx, ClassifyRequest{
		Message:  trimmed,
		Language: req.Language,
		History:  req.History,
	})
	if err != nil {
		return nil, fmt.Errorf("classification error: %w", err)
	}

	// 2. Reject Out-of-Domain Inquiries Locally (0 Gemini Calls, 0 DB Queries)
	if !classifyRes.IsAllowed {
		return &ChatResponse{
			Reply:            classifyRes.InstantReply,
			Language:         classifyRes.Language,
			Intent:           classifyRes.Intent,
			Path:             classifyRes.Path,
			IsOutOfDomain:    true,
			FallbackUsed:     false,
			TokensUsed:       0,
			LatencyMs:        time.Since(t0).Milliseconds(),
			SuggestedQueries: classifyRes.SuggestedQueries,
		}, nil
	}

	// 3. Fast Deterministic Greeting (< 1ms, 0 Gemini Calls)
	if classifyRes.Intent == "GREETING" && classifyRes.InstantReply != "" {
		return &ChatResponse{
			Reply:            classifyRes.InstantReply,
			Language:         classifyRes.Language,
			Intent:           "GREETING",
			Path:             "PATH_ADVISORY",
			IsOutOfDomain:    false,
			FallbackUsed:     false,
			TokensUsed:       0,
			LatencyMs:        time.Since(t0).Milliseconds(),
			SuggestedQueries: classifyRes.SuggestedQueries,
		}, nil
	}

	// 4. Product Database Authority: Retrieve real database listings when path is search or filters present
	var cards []ListingCard
	if classifyRes.Path == "PATH_SEARCH" || classifyRes.SearchTerm != "" || classifyRes.MinPrice != nil || classifyRes.MaxPrice != nil {
		cards, _ = s.productTool.SearchListings(ctx, classifyRes.SearchTerm, classifyRes.MinPrice, classifyRes.MaxPrice, 6)
	}

	// 5. Daily Marketplace Intelligence (Top categories & searches)
	trendSummary := s.trendTool.GetTrendSummary(ctx)

	// 6. Python Gemini Orchestration: Send prompt augmented with authoritative DB listings
	pyChatReq := PythonChatRequest{
		Message:      trimmed,
		Language:     classifyRes.Language,
		History:      req.History,
		Listings:     cards,
		TrendSummary: trendSummary,
	}

	chatRes, err := s.supervisor.Chat(ctx, pyChatReq)
	if err != nil {
		return nil, fmt.Errorf("chat processing error: %w", err)
	}

	// 7. Attach Real Database Listing Cards to Response for UI Card Rendering
	chatRes.Listings = cards
	chatRes.LatencyMs = time.Since(t0).Milliseconds()

	return chatRes, nil
}

func (s *AIService) Health(ctx context.Context) error {
	return s.supervisor.healthCheck(ctx)
}
