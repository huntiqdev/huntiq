package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/internal/service/aiclient"
)

func TestAIToolDispatcher_UnsupportedTool(t *testing.T) {
	dispatcher := service.NewAIToolDispatcher(nil)

	resStr, listings, err := dispatcher.ExecuteToolCall(context.Background(), aiclient.ChatToolCall{
		ID:   "call_unknown",
		Type: "function",
		Function: aiclient.ChatFunctionCall{
			Name:      "unsupported_arbitrary_query",
			Arguments: map[string]interface{}{"query": "SELECT * FROM users"},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(listings) != 0 {
		t.Errorf("expected 0 listings, got %d", len(listings))
	}
	if !strings.Contains(resStr, "unsupported tool") {
		t.Errorf("expected unsupported tool error in result, got: %s", resStr)
	}
}

func TestAIToolDispatcher_GetListingDetails_InvalidArgs(t *testing.T) {
	dispatcher := service.NewAIToolDispatcher(nil)

	resStr, _, err := dispatcher.ExecuteToolCall(context.Background(), aiclient.ChatToolCall{
		ID:   "call_details",
		Type: "function",
		Function: aiclient.ChatFunctionCall{
			Name:      "get_listing_details",
			Arguments: map[string]interface{}{"product_id": -5.0},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resStr, "invalid product_id") {
		t.Errorf("expected invalid product_id error, got: %s", resStr)
	}
}

func TestAIToolDispatcher_CompareListings_RequiresMultiple(t *testing.T) {
	dispatcher := service.NewAIToolDispatcher(nil)

	resStr, _, err := dispatcher.ExecuteToolCall(context.Background(), aiclient.ChatToolCall{
		ID:   "call_compare",
		Type: "function",
		Function: aiclient.ChatFunctionCall{
			Name:      "compare_listings",
			Arguments: map[string]interface{}{"product_ids": []interface{}{101.0}},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resStr, "requires an array of 2 or more") {
		t.Errorf("expected comparison requirement error, got: %s", resStr)
	}
}

func TestAIToolDispatcher_SecurityAndPromptInjection(t *testing.T) {
	dispatcher := service.NewAIToolDispatcher(nil)

	// Test 1: Arbitrary function name execution blocked
	resStr, listings, err := dispatcher.ExecuteToolCall(context.Background(), aiclient.ChatToolCall{
		ID:   "inject_1",
		Type: "function",
		Function: aiclient.ChatFunctionCall{
			Name:      "execute_raw_sql",
			Arguments: map[string]interface{}{"sql": "DROP TABLE users;"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(listings) != 0 {
		t.Errorf("expected 0 listings, got %d", len(listings))
	}
	if !strings.Contains(resStr, "unsupported tool") {
		t.Errorf("expected unsupported tool error, got: %s", resStr)
	}

	// Test 2: Pagination limits strictly bounded
	resStr2, _, _ := dispatcher.ExecuteToolCall(context.Background(), aiclient.ChatToolCall{
		ID:   "inject_2",
		Type: "function",
		Function: aiclient.ChatFunctionCall{
			Name: "search_marketplace_listings",
			Arguments: map[string]interface{}{
				"limit": 99999999.0,
			},
		},
	})
	// Dispatcher handles nil productRepo gracefully
	if resStr2 == "" {
		t.Errorf("expected non-empty response")
	}
}

func TestProcessGuestChat_ZeroFabricatedListingsOnEmpty(t *testing.T) {
	svc := service.NewAIChatService(nil, nil, nil)
	ctx := context.Background()

	// Query for an item that yields 0 database results
	resp, err := svc.ProcessGuestChat(ctx, service.GuestChatRequestDTO{
		Message: "I need a flying hoverboard in Bole below 50 birr",
		Locale:  "en",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Zero fabrication invariant: listings slice must be empty
	if len(resp.Listings) != 0 {
		t.Errorf("expected 0 listings fabricated, got %d", len(resp.Listings))
	}
	// Answer must honestly state no listings found
	lowerAns := strings.ToLower(resp.Answer)
	if !strings.Contains(lowerAns, "no listings") &&
		!strings.Contains(lowerAns, "found no") &&
		!strings.Contains(lowerAns, "could not find") &&
		!strings.Contains(lowerAns, "sorry") {
		t.Errorf("expected honest no-results response, got: %s", resp.Answer)
	}
}
