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

	// Test D: Compliments ("You are clever") -> Fast reply (0 tokens)
	respComp, err := aiSvc.ProcessQuestion(ctx, ChatRequest{Message: "You are clever"})
	if err != nil {
		t.Fatalf("ProcessQuestion failed on compliment: %v", err)
	}
	if respComp.Intent != "COMPLIMENT" {
		t.Errorf("Expected intent COMPLIMENT, got %s", respComp.Intent)
	}
	if respComp.TokensUsed != 0 {
		t.Errorf("Expected tokens_used=0 for compliment, got %d", respComp.TokensUsed)
	}

	// Test E: Appreciation ("Thank you") -> Fast reply (0 tokens)
	respThanks, err := aiSvc.ProcessQuestion(ctx, ChatRequest{Message: "Thank you"})
	if err != nil {
		t.Fatalf("ProcessQuestion failed on appreciation: %v", err)
	}
	if respThanks.Intent != "APPRECIATION" {
		t.Errorf("Expected intent APPRECIATION, got %s", respThanks.Intent)
	}
	if respThanks.TokensUsed != 0 {
		t.Errorf("Expected tokens_used=0 for appreciation, got %d", respThanks.TokensUsed)
	}

	// Test F: Capability / Identity ("Who are you?") -> Fast reply (0 tokens)
	respWho, err := aiSvc.ProcessQuestion(ctx, ChatRequest{Message: "Who are you?"})
	if err != nil {
		t.Fatalf("ProcessQuestion failed on capability: %v", err)
	}
	if respWho.Intent != "CAPABILITY_QUESTION" {
		t.Errorf("Expected intent CAPABILITY_QUESTION, got %s", respWho.Intent)
	}
	if respWho.TokensUsed != 0 {
		t.Errorf("Expected tokens_used=0 for capability, got %d", respWho.TokensUsed)
	}

	// Test G: Privacy Guard ("Who registered yesterday?") -> Blocked / Privacy explanation (0 tokens)
	respPriv, err := aiSvc.ProcessQuestion(ctx, ChatRequest{Message: "Who registered yesterday?"})
	if err != nil {
		t.Fatalf("ProcessQuestion failed on privacy probe: %v", err)
	}
	if !respPriv.IsOutOfDomain {
		t.Errorf("Expected is_out_of_domain=true for privacy query, got false")
	}
	if respPriv.TokensUsed != 0 {
		t.Errorf("Expected tokens_used=0 for privacy query, got %d", respPriv.TokensUsed)
	}

	// Test H: Pure Math ("Solve 25 x 37") -> Blocked
	classifyMath, err := supervisor.Classify(ctx, ClassifyRequest{Message: "Solve 25 x 37"})
	if err != nil {
		t.Fatalf("Classify failed on pure math: %v", err)
	}
	if classifyMath.IsAllowed {
		t.Errorf("Expected pure math to be rejected, got is_allowed=true")
	}

	// Test I: Marketplace Math ("Calculate 10% commission on 100,000 ETB") -> Allowed
	classifyMktMath, err := supervisor.Classify(ctx, ClassifyRequest{Message: "Calculate 10% commission on 100,000 ETB"})
	if err != nil {
		t.Fatalf("Classify failed on marketplace math: %v", err)
	}
	if !classifyMktMath.IsAllowed {
		t.Errorf("Expected marketplace math to be allowed, got is_allowed=false")
	}
	if classifyMktMath.Path != "PATH_ADVISORY" {
		t.Errorf("Expected path PATH_ADVISORY for marketplace math, got %s", classifyMktMath.Path)
	}

	// Test J: Advisory Query ("Which PC is suitable for a software engineer?") -> Allowed PATH_ADVISORY
	classifyPC, err := supervisor.Classify(ctx, ClassifyRequest{Message: "Which PC is suitable for a software engineer?"})
	if err != nil {
		t.Fatalf("Classify failed on PC advisory: %v", err)
	}
	if !classifyPC.IsAllowed {
		t.Errorf("Expected PC advisory to be allowed, got is_allowed=false")
	}
	if classifyPC.Path != "PATH_ADVISORY" {
		t.Errorf("Expected path PATH_ADVISORY for PC advisory, got %s", classifyPC.Path)
	}

	// Test K: Contextual Follow-up ("What about Lenovo?") -> Allowed PATH_ADVISORY with history
	classifyFollowup, err := supervisor.Classify(ctx, ClassifyRequest{
		Message: "What about Lenovo?",
		History: []ChatMessage{
			{Role: "user", Parts: []struct {
				Text string `json:"text"`
			}{{Text: "Which PC is suitable for a software engineer?"}}},
		},
	})
	if err != nil {
		t.Fatalf("Classify failed on contextual follow-up: %v", err)
	}
	if !classifyFollowup.IsAllowed {
		t.Errorf("Expected follow-up to be allowed, got is_allowed=false")
	}
	if classifyFollowup.Path != "PATH_ADVISORY" {
		t.Errorf("Expected path PATH_ADVISORY for follow-up, got %s", classifyFollowup.Path)
	}
}
