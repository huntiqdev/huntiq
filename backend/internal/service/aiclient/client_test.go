package aiclient

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/cache"
)

func TestAIClient_DisabledWhenEmpty(t *testing.T) {
	c := NewClient("", "")
	if c.IsEnabled() {
		t.Fatalf("expected client to be disabled when empty")
	}

	_, err := c.UnderstandQuery(context.Background(), "test", "en")
	if err == nil {
		t.Fatalf("expected error when calling disabled client")
	}
}

func TestAIClient_UnderstandQuerySuccess(t *testing.T) {
	secret := "secret-token-12345"
	expectedResp := QueryUnderstandingResponse{
		NormalizedQuery:      "toyota hilux diesel",
		PrimaryIntent:        "vehicle_search",
		Confidence:           0.95,
		SuggestedCategoryIDs: []int{74, 77},
		ExtractedEntities: ExtractedEntities{
			Brand: strPtr("Toyota"),
			Model: strPtr("Hilux"),
			PriceRange: ExtractedPriceRange{
				MaxPrice: floatPtr(2000000.0),
				Currency: "ETB",
			},
			Attributes: map[string]any{"fuel_type": "diesel"},
		},
		DetectedLocale:  "en",
		ExecutionTimeMS: 12.5,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/query/understand" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Header.Get("X-HUNTIQ-AI-SECRET") != secret {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expectedResp)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, secret)
	if !client.IsEnabled() {
		t.Fatalf("expected client to be enabled")
	}

	resp, err := client.UnderstandQuery(context.Background(), "toyota hilux diesel under 2M ETB", "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.PrimaryIntent != "vehicle_search" {
		t.Errorf("expected vehicle_search, got %s", resp.PrimaryIntent)
	}
	if resp.Confidence < 0.90 {
		t.Errorf("expected confidence >= 0.90, got %f", resp.Confidence)
	}
	if resp.ExtractedEntities.Brand == nil || *resp.ExtractedEntities.Brand != "Toyota" {
		t.Errorf("expected brand Toyota, got %v", resp.ExtractedEntities.Brand)
	}
}

func TestAIClient_TimeoutTriggersError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sleep beyond context timeout
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "secret")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.UnderstandQuery(ctx, "slow query", "en")
	if err == nil {
		t.Fatalf("expected error on timeout")
	}
}

func TestAIClient_UnderstandQueryWithParams_AIMode(t *testing.T) {
	secret := "secret-token-12345"
	expectedResp := QueryUnderstandingResponse{
		NormalizedQuery: "mana kireeffamu finfinnee",
		Mode:            "ai",
		PrimaryIntent:   "rent_property",
		Confidence:      0.95,
		LocalizedAnswer: "Finfinnee keessatti mana kireeffamu qopheessineerra.",
		PIIDetected:     false,
		DeepParameters: map[string]any{
			"location": "Addis Ababa",
			"intent":   "rent_property",
		},
		DetectedLocale:  "om",
		ExecutionTimeMS: 8.2,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req UnderstandQueryParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.Mode != "ai" || req.PreferredLanguage != "om" {
			t.Errorf("unexpected params: %+v", req)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expectedResp)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, secret)
	resp, err := client.UnderstandQueryWithParams(context.Background(), UnderstandQueryParams{
		Query:             "mana kireeffamu finfinnee",
		Mode:              "ai",
		PreferredLanguage: "om",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Mode != "ai" {
		t.Errorf("expected mode ai, got %s", resp.Mode)
	}
	if resp.LocalizedAnswer != "Finfinnee keessatti mana kireeffamu qopheessineerra." {
		t.Errorf("unexpected localized answer: %s", resp.LocalizedAnswer)
	}
	if resp.DeepParameters["location"] != "Addis Ababa" {
		t.Errorf("expected Addis Ababa, got %v", resp.DeepParameters["location"])
	}
}

func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

func TestAIClient_AnalyzeProducts_DeterministicFallback(t *testing.T) {
	products := []ProductCandidateSummary{
		{
			ID:         1,
			Title:      "Samsung Galaxy S22 Ultra",
			Price:      45000.0,
			SellerTier: "premium",
		},
		{
			ID:         2,
			Title:      "Samsung Galaxy A14",
			Price:      12500.0,
			SellerTier: "gold",
		},
		{
			ID:         3,
			Title:      "Samsung Galaxy A04",
			Price:      8500.0,
			SellerTier: "normal",
		},
	}

	req := ProductAnalysisRequest{
		Query:             "samsung phone",
		Locale:            "en",
		PreferredLanguage: "en",
		Products:          products,
	}

	resp := ComputeDeterministicAnalysis(req)
	if resp.TotalAnalyzed != 3 {
		t.Errorf("expected 3 analyzed, got %d", resp.TotalAnalyzed)
	}
	if resp.PremiumCount != 1 {
		t.Errorf("expected 1 premium, got %d", resp.PremiumCount)
	}
	if resp.GoldCount != 1 {
		t.Errorf("expected 1 gold, got %d", resp.GoldCount)
	}
	if resp.NormalCount != 1 {
		t.Errorf("expected 1 normal, got %d", resp.NormalCount)
	}
	if resp.PriceStats.MinPrice != 8500.0 {
		t.Errorf("expected min 8500, got %f", resp.PriceStats.MinPrice)
	}
	if resp.TopPremiumPick == nil || resp.TopPremiumPick.ID != 1 {
		t.Errorf("expected top premium pick ID 1, got %v", resp.TopPremiumPick)
	}
	if resp.TopGoldPick == nil || resp.TopGoldPick.ID != 2 {
		t.Errorf("expected top gold pick ID 2, got %v", resp.TopGoldPick)
	}
	if len(resp.DealRatings) != 3 {
		t.Errorf("expected 3 deal ratings, got %d", len(resp.DealRatings))
	}
}

func TestAIClient_AnalyzeProducts_Server(t *testing.T) {
	secret := "secret-12345"
	expectedResp := ProductAnalysisResponse{
		TotalAnalyzed: 2,
		PremiumCount:  1,
		GoldCount:     1,
		NormalCount:   0,
		PriceStats: PriceStats{
			MinPrice: 12000,
			AvgPrice: 20000,
			MaxPrice: 28000,
			Currency: "ETB",
		},
		MarketVerdict: "Real-time analysis complete",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/query/analyze-products" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Header.Get("X-HUNTIQ-AI-SECRET") != secret {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expectedResp)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, secret)
	resp, err := client.AnalyzeProducts(context.Background(), ProductAnalysisRequest{
		Query:  "test",
		Locale: "en",
		Products: []ProductCandidateSummary{
			{ID: 1, Title: "Item 1", Price: 12000, SellerTier: "gold"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalAnalyzed != 2 {
		t.Errorf("expected 2 total analyzed, got %d", resp.TotalAnalyzed)
	}
}

func TestCleanConversationalSearchTerm(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"I want a phone.", "phone"},
		{"i wenna phone", "phone"},
		{"Find a used phone near me in Bole", "phone"},
		{"I need an apartment for rent", "apartment"},
		{"What products are available now?", ""},
		{"What is available?", ""},
		// Conversational questions in English
		{"have you phone?", "phone"},
		{"have you a phone?", "phone"},
		{"have you got any phone?", "phone"},
		{"do you have phone?", "phone"},
		{"do you have a phone?", "phone"},
		{"can i get a phone?", "phone"},
		{"is there a phone?", "phone"},
		{"do you have a car?", "car"},
		{"have you got a laptop?", "laptop"},
		// Conversational questions in Amharic
		{"ስልክ መግዛት እፈልጋለሁ", "ስልክ"},
		{"ስልክ አለዎት?", "ስልክ"},
		{"ስልክ አለህ?", "ስልክ"},
		{"ስልክ አለሽ?", "ስልክ"},
		{"ስልክ አላችሁ?", "ስልክ"},
		{"ስልክ አለ?", "ስልክ"},
		{"ስልክ አለ ወይ?", "ስልክ"},
		{"ስልክ አለ እንዴ?", "ስልክ"},
		{"ስልክ ፈልጌ ነው", "ስልክ"},
		{"ስልክ ፈልጌ ነበር", "ስልክ"},
		{"ስልክ ይኖራል?", "ስልክ"},
		{"ስልክ ይገኛል?", "ስልክ"},
		{"ስልክ ማግኘት እችላለሁ?", "ስልክ"},
		{"ስልክ ስንት ነው?", "ስልክ"},
		{"የስልክ ዋጋ ስንት ነው?", "ስልክ"},
		{"መኪና አለዎት?", "መኪና"},
		{"መኪና አለህ?", "መኪና"},
		{"ቤት አለዎት?", "ቤት"},
		// Conversational questions in Afaan Oromo
		{"Bilbila bituu barbaada", "bilbila"},
		{"bilbila qabduu?", "bilbila"},
		{"bilbila qabdu?", "bilbila"},
		{"bilbila jiraa?", "bilbila"},
		{"konkolaataa qabdu?", "konkolaataa"},
		// Conversational questions in Tigrinya
		{"ሞባይል ክገዝእ ደልየ", "ሞባይል"},
		{"ሞባይል ኣለኩም?", "ሞባይል"},
		{"መኪና ኣለካ?", "መኪና"},
	}

	for _, tc := range tests {
		actual := CleanConversationalSearchTerm(tc.input)
		if actual != tc.expected {
			t.Errorf("CleanConversationalSearchTerm(%q) = %q, expected %q", tc.input, actual, tc.expected)
		}
	}
}

func TestParseFallbackSearchArgs_SingleTurn(t *testing.T) {
	// Query 1: "I want a phone."
	hasIntent, args := ParseFallbackSearchArgs(nil, "I want a phone.")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for 'I want a phone.'")
	}
	if args["search_term"] != "phone" {
		t.Errorf("expected search_term='phone', got %v", args["search_term"])
	}

	// Query 4: "Find a used phone near me in Bole"
	hasIntent, args = ParseFallbackSearchArgs(nil, "Find a used phone near me in Bole")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for Bole query")
	}
	if args["search_term"] != "phone" {
		t.Errorf("expected search_term='phone', got %v", args["search_term"])
	}
	if args["condition"] != "used" {
		t.Errorf("expected condition='used', got %v", args["condition"])
	}
	if args["location"] != "Bole" {
		t.Errorf("expected location='Bole', got %v", args["location"])
	}

	// Query 5: "I need an apartment for rent"
	hasIntent, args = ParseFallbackSearchArgs(nil, "I need an apartment for rent")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for apartment query")
	}
	if args["search_term"] != "apartment" {
		t.Errorf("expected search_term='apartment', got %v", args["search_term"])
	}
	if args["listing_type"] != "rent" {
		t.Errorf("expected listing_type='rent', got %v", args["listing_type"])
	}

	// Multilingual test 7 (Amharic): "ስልክ መግዛት እፈልጋለሁ"
	hasIntent, args = ParseFallbackSearchArgs(nil, "ስልክ መግዛት እፈልጋለሁ")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for Amharic phone query")
	}
	if args["search_term"] != "phone" {
		t.Errorf("expected search_term='phone', got %v", args["search_term"])
	}
	if args["localized_term"] != "ስልክ" {
		t.Errorf("expected localized_term='ስልክ', got %v", args["localized_term"])
	}

	// Multilingual test 8 (Afaan Oromoo): "Bilbila bituu barbaada"
	hasIntent, args = ParseFallbackSearchArgs(nil, "Bilbila bituu barbaada")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for Oromo phone query")
	}
	if args["search_term"] != "phone" {
		t.Errorf("expected search_term='phone', got %v", args["search_term"])
	}

	// Multilingual test 9 (Tigrinya): "ሞባይል ክገዝእ ደልየ"
	hasIntent, args = ParseFallbackSearchArgs(nil, "ሞባይል ክገዝእ ደልየ")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for Tigrinya phone query")
	}
	if args["search_term"] != "phone" {
		t.Errorf("expected search_term='phone', got %v", args["search_term"])
	}

	// Non-commercial query 10: "What is the capital of France?"
	hasIntent, _ = ParseFallbackSearchArgs(nil, "What is the capital of France?")
	if hasIntent {
		t.Errorf("expected hasIntent=false for trivia question")
	}
}

func TestParseFallbackSearchArgs_MultiTurnInheritance(t *testing.T) {
	// Turn 1: user asks for a phone
	messages := []ChatLLMMessage{
		{Role: "user", Content: "I want a phone."},
		{
			Role: "assistant",
			ToolCalls: []ChatToolCall{
				{
					ID:   "call_1",
					Type: "function",
					Function: ChatFunctionCall{
						Name:      "search_marketplace_listings",
						Arguments: map[string]any{"search_term": "phone"},
					},
				},
			},
		},
	}

	// Turn 2: user follows up with: "Only ones under 40,000 ETB"
	hasIntent, args := ParseFallbackSearchArgs(messages, "Only ones under 40,000 ETB")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for multi-turn follow up")
	}
	if args["search_term"] != "phone" {
		t.Errorf("expected inherited search_term='phone', got %v", args["search_term"])
	}
	if maxP, ok := args["max_price"].(float64); !ok || maxP != 40000.0 {
		t.Errorf("expected max_price=40000.0, got %v", args["max_price"])
	}
}

func TestConversationalSearchQuestions(t *testing.T) {
	// English conversational question
	hasIntent, args := ParseFallbackSearchArgs(nil, "have you phone?")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for 'have you phone?'")
	}
	if args["search_term"] != "phone" {
		t.Errorf("expected search_term='phone', got %v", args["search_term"])
	}

	// Amharic conversational question
	hasIntent, args = ParseFallbackSearchArgs(nil, "ስልክ አለዎት?")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for 'ስልክ አለዎት?'")
	}
	if args["search_term"] != "phone" {
		t.Errorf("expected search_term='phone', got %v", args["search_term"])
	}
	if args["localized_term"] != "ስልክ" {
		t.Errorf("expected localized_term='ስልክ', got %v", args["localized_term"])
	}

	// Amharic car question
	hasIntent, args = ParseFallbackSearchArgs(nil, "መኪና አለህ?")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for 'መኪና አለህ?'")
	}
	if args["search_term"] != "car" {
		t.Errorf("expected search_term='car', got %v", args["search_term"])
	}

	// Afaan Oromoo phone question
	hasIntent, args = ParseFallbackSearchArgs(nil, "bilbila qabduu?")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for 'bilbila qabduu?'")
	}
	if args["search_term"] != "phone" {
		t.Errorf("expected search_term='phone', got %v", args["search_term"])
	}

	// Afaan Oromoo car question
	hasIntent, args = ParseFallbackSearchArgs(nil, "konkolaataa qabdaa?")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for 'konkolaataa qabdaa?'")
	}
	if args["search_term"] != "car" {
		t.Errorf("expected search_term='car', got %v", args["search_term"])
	}

	// Tigrinya phone question
	hasIntent, args = ParseFallbackSearchArgs(nil, "ስልኪ ኣለኩም ዶ?")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for 'ስልኪ ኣለኩም ዶ?'")
	}
	if args["search_term"] != "phone" {
		t.Errorf("expected search_term='phone', got %v", args["search_term"])
	}

	// English "do you have a car?"
	hasIntent, args = ParseFallbackSearchArgs(nil, "do you have a car?")
	if !hasIntent {
		t.Fatalf("expected hasIntent=true for 'do you have a car?'")
	}
	if args["search_term"] != "car" {
		t.Errorf("expected search_term='car', got %v", args["search_term"])
	}
}

func TestSearchSynonymsAndLocalization(t *testing.T) {
	// 1. Synonym group for phone
	syns := GetSearchSynonyms("phone")
	if len(syns) <= 1 {
		t.Fatalf("expected multiple synonyms for 'phone', got %v", syns)
	}

	// 2. Cross-lingual mapping from Amharic 'ስልክ'
	amharicSyns := GetSearchSynonyms("ስልክ")
	hasIphone := false
	for _, s := range amharicSyns {
		if s == "iphone" {
			hasIphone = true
			break
		}
	}
	if !hasIphone {
		t.Errorf("expected 'iphone' in synonyms for Amharic 'ስልክ', got %v", amharicSyns)
	}

	// 3. TsQuery construction
	tsQuery := BuildSynonymTsQuery("phone")
	if !strings.Contains(tsQuery, "'phone'") || !strings.Contains(tsQuery, "'ስልክ'") {
		t.Errorf("expected tsquery to contain both 'phone' and 'ስልክ', got %s", tsQuery)
	}

	// 4. Entity localization
	if loc := GetLocalizedEntityTerm("phone", "am"); loc != "ስልክ" {
		t.Errorf("expected localized term 'ስልክ' for am, got %s", loc)
	}
	if loc := GetLocalizedEntityTerm("phone", "om"); loc != "bilbila" {
		t.Errorf("expected localized term 'bilbila' for om, got %s", loc)
	}
	if loc := GetLocalizedEntityTerm("phone", "ti"); loc != "ሞባይል" {
		t.Errorf("expected localized term 'ሞባይል' for ti, got %s", loc)
	}
	if loc := GetLocalizedEntityTerm("phone", "en"); loc != "phone" {
		t.Errorf("expected localized term 'phone' for en, got %s", loc)
	}
}

func TestUnderstandQuery_CacheAndCircuitBreaker(t *testing.T) {
	memCache := cache.NewMemoryCache(1 * time.Minute)
	defer memCache.Close()

	// Client with non-existent server
	client := NewClient("http://127.0.0.1:9999", "test-secret").
		WithCache(memCache).
		WithCircuitBreaker(NewCircuitBreaker(2, 50*time.Millisecond))

	// 1. Initial call fails (server down) -> error
	_, err := client.UnderstandQuery(context.Background(), "Toyota Corolla", "en")
	if err == nil {
		t.Fatalf("expected error on offline server")
	}

	// 2. Second failure should trip circuit breaker to OPEN
	_, _ = client.UnderstandQuery(context.Background(), "Toyota Corolla", "en")
	if client.circuitBreaker.State() != StateOpen {
		t.Fatalf("expected circuit breaker state OPEN after 2 failures, got %s", client.circuitBreaker.State())
	}

	// 3. When circuit breaker is OPEN, client immediately returns fallback without waiting
	t0 := time.Now()
	res, err := client.UnderstandQuery(context.Background(), "Toyota Corolla Bole", "en")
	if err != nil {
		t.Fatalf("expected immediate fallback, got error: %v", err)
	}
	if time.Since(t0) > 50*time.Millisecond {
		t.Fatalf("expected sub-millisecond fallback when circuit open, took %v", time.Since(t0))
	}
	if !strings.Contains(res.PrimaryIntent, "buy") {
		t.Errorf("expected buy intent, got %s", res.PrimaryIntent)
	}

	// 4. Test Cache Hit
	cachedQuery := QueryUnderstandingResponse{
		NormalizedQuery: "cached-query",
		PrimaryIntent:   "cached_intent",
		Confidence:      0.99,
	}
	cacheData, _ := json.Marshal(cachedQuery)
	cacheKey := fmt.Sprintf("ai:understand:%x", sha256.Sum256([]byte("en:test-cache-hit")))
	_ = memCache.Set(context.Background(), cacheKey, cacheData, 10*time.Minute)

	tHit := time.Now()
	hitRes, err := client.UnderstandQuery(context.Background(), "test-cache-hit", "en")
	if err != nil {
		t.Fatalf("expected cache hit, got error: %v", err)
	}
	if time.Since(tHit) > 20*time.Millisecond {
		t.Fatalf("expected instant cache hit, took %v", time.Since(tHit))
	}
	if hitRes.PrimaryIntent != "cached_intent" {
		t.Errorf("expected cached_intent, got %s", hitRes.PrimaryIntent)
	}
}


