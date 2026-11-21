package ai

import (
	"context"
	"testing"
	"time"
)

func TestPythonSupervisorLifecycleAndAIService(t *testing.T) {
	// 1. Initialize PythonSupervisor pointing to script.py and policy.json
	supervisor := NewPythonSupervisor("", "gemini-3.5-flash-lite")
	if err := supervisor.Start(); err != nil {
		t.Fatalf("Failed to start Python AI supervisor: %v", err)
	}
	defer supervisor.Stop()

	// 2. Health check verification
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := supervisor.healthCheck(ctx); err != nil {
		t.Fatalf("Health check failed on running supervisor: %v", err)
	}

	// 3. Classify Out-of-Domain Query
	classifyOut, err := supervisor.Classify(ctx, ClassifyRequest{
		Message: "What is E=mc²?",
	})
	if err != nil {
		t.Fatalf("Classify failed: %v", err)
	}
	if classifyOut.IsAllowed {
		t.Errorf("Expected 'What is E=mc²?' to be rejected, but got is_allowed=true")
	}
	if classifyOut.Intent != "OUT_OF_DOMAIN" {
		t.Errorf("Expected intent OUT_OF_DOMAIN, got %s", classifyOut.Intent)
	}
	if classifyOut.Path != "PATH_BLOCKED" {
		t.Errorf("Expected path PATH_BLOCKED, got %s", classifyOut.Path)
	}

	// 4. Classify Fast Greeting
	classifyGreet, err := supervisor.Classify(ctx, ClassifyRequest{
		Message: "Hello",
	})
	if err != nil {
		t.Fatalf("Classify failed: %v", err)
	}
	if !classifyGreet.IsAllowed {
		t.Errorf("Expected 'Hello' to be allowed, but got is_allowed=false")
	}
	if classifyGreet.Intent != "GREETING" {
		t.Errorf("Expected intent GREETING, got %s", classifyGreet.Intent)
	}
	if classifyGreet.InstantReply == "" {
		t.Errorf("Expected non-empty instant reply for greeting")
	}

	// 5. Classify Marketplace Commercial Advice Query
	classifyAdvice, err := supervisor.Classify(ctx, ClassifyRequest{
		Message: "How can I increase sales for my listings on HUNTIQ?",
	})
	if err != nil {
		t.Fatalf("Classify failed: %v", err)
	}
	if !classifyAdvice.IsAllowed {
		t.Errorf("Expected commercial advice to be allowed, got is_allowed=false")
	}
	if classifyAdvice.Path != "PATH_ADVISORY" {
		t.Errorf("Expected path PATH_ADVISORY, got %s", classifyAdvice.Path)
	}

	// 6. Classify Search Query with Price Filters
	classifySearch, err := supervisor.Classify(ctx, ClassifyRequest{
		Message: "Find laptops under 50000 ETB",
	})
	if err != nil {
		t.Fatalf("Classify failed: %v", err)
	}
	if !classifySearch.IsAllowed {
		t.Errorf("Expected laptop search to be allowed, got is_allowed=false")
	}
	if classifySearch.Path != "PATH_SEARCH" {
		t.Errorf("Expected path PATH_SEARCH, got %s", classifySearch.Path)
	}
	if classifySearch.MaxPrice == nil || *classifySearch.MaxPrice != 50000.0 {
		t.Errorf("Expected max_price=50000.0, got %v", classifySearch.MaxPrice)
	}

	// 7. Test AIService Orchestration (Mocking repositories as nil)
	aiSvc := NewAIService(supervisor, nil, nil)

	// Test A: Out-of-Domain Query should return instant reply without error
	respOut, err := aiSvc.ProcessQuestion(ctx, ChatRequest{Message: "Write a python hello world script"})
	if err != nil {
		t.Fatalf("ProcessQuestion failed on out-of-domain: %v", err)
	}
	if !respOut.IsOutOfDomain {
		t.Errorf("Expected is_out_of_domain=true for programming question")
	}
	if respOut.Reply == "" {
		t.Errorf("Expected non-empty reply for out-of-domain")
	}

	// Test B: Greeting Query should return instant reply with 0 tokens and low latency
	respGreet, err := aiSvc.ProcessQuestion(ctx, ChatRequest{Message: "ሰላም"})
	if err != nil {
		t.Fatalf("ProcessQuestion failed on greeting: %v", err)
	}
	if respGreet.Intent != "GREETING" {
		t.Errorf("Expected intent GREETING, got %s", respGreet.Intent)
	}
	if respGreet.TokensUsed != 0 {
		t.Errorf("Expected tokens_used=0 for greeting, got %d", respGreet.TokensUsed)
	}

	// Test C: Empty message should error
	_, errEmpty := aiSvc.ProcessQuestion(ctx, ChatRequest{Message: "   "})
	if errEmpty == nil {
		t.Errorf("Expected error for empty message, got nil")
	}
}
