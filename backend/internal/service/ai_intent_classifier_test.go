package service_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/service"
)

// 1. Social Greetings & Conversational Etiquette
func TestConversational_Greetings(t *testing.T) {
	greetings := []struct {
		query    string
		locale   string
		expected string
	}{
		{"hello", "en", "HUNTIQ"},
		{"hi there", "en", "HUNTIQ"},
		{"ሰላም", "am", "ሰላም"},
		{"እንደምን አለህ", "am", "እንደምን"},
		{"akkam", "om", "Akkam"},
		{"nagaa", "om", "Nagaa"},
		{"ሰላም", "ti", "ሰላም"},
	}

	for _, tc := range greetings {
		plan := service.ClassifyMessage(tc.query, tc.locale, nil)
		if plan.RequiresDatabaseSearch {
			t.Errorf("query '%s' should NOT require database search, but plan marked it true", tc.query)
		}
		if len(plan.Searches) > 0 {
			t.Errorf("query '%s' should NOT have extracted searches, got %d", tc.query, len(plan.Searches))
		}
		if plan.SocialGreeting == "" {
			t.Errorf("query '%s' expected social greeting, got empty", tc.query)
		}
	}
}

// 2. Identity & System Nature Inquiries
func TestConversational_IdentityInquiries(t *testing.T) {
	identityQueries := []struct {
		query   string
		locale  string
		mustSub string
	}{
		{"what is your name?", "en", "HUNTIQ"},
		{"who are you?", "en", "assistant"},
		{"who made you?", "en", "HUNTIQ"},
		{"what is huntiq?", "en", "marketplace"},
		{"ስምህ ማን ነው?", "am", "HUNTIQ"},
		{"ማን ነህ?", "am", "ረዳት"},
		{"maqaan kee eenyu?", "om", "HUNTIQ"},
		{"ati eenyu?", "om", "gargaaraa"},
		{"መን ኢኻ?", "ti", "ሓጋዚ"},
	}

	for _, tc := range identityQueries {
		plan := service.ClassifyMessage(tc.query, tc.locale, nil)
		if plan.RequiresDatabaseSearch {
			t.Errorf("identity query '%s' should NOT trigger DB search", tc.query)
		}
		if len(plan.Searches) > 0 {
			t.Errorf("identity query '%s' should NOT have search terms, got: %v", tc.query, plan.Searches)
		}
		if !strings.Contains(strings.ToLower(plan.IdentityResponse), strings.ToLower(tc.mustSub)) {
			t.Errorf("expected identity response for '%s' to contain '%s', got: %s", tc.query, tc.mustSub, plan.IdentityResponse)
		}
	}
}

// 3. Ambiguous Sub-4 Character Token Clarification
func TestConversational_AmbiguousSubTokens(t *testing.T) {
	subTokens := []struct {
		query    string
		locale   string
		expected string
	}{
		{"pho", "en", "phone"},
		{"lap", "en", "laptop"},
		{"apa", "en", "apartment"},
		{"car", "en", ""}, // car is complete (3 chars but legitimate)
	}

	for _, tc := range subTokens {
		plan := service.ClassifyMessage(tc.query, tc.locale, nil)
		if tc.expected != "" {
			if !plan.RequiresClarification {
				t.Errorf("expected '%s' to require clarification", tc.query)
			}
			if !strings.Contains(strings.ToLower(plan.ClarificationQuestion), tc.expected) {
				t.Errorf("expected clarification for '%s' to suggest '%s', got: %s", tc.query, tc.expected, plan.ClarificationQuestion)
			}
		} else {
			// legitimate short word like "car"
			if len(plan.Searches) == 0 {
				t.Errorf("expected '%s' to be recognized as valid search", tc.query)
			}
		}
	}
}

// 4. Combined Greeting and Search
func TestConversational_GreetingAndSearch(t *testing.T) {
	queries := []struct {
		query        string
		locale       string
		expectedCat  string
		expectedTerm string
	}{
		{"hello, I want a phone", "en", "Mobile Phones", "phone"},
		{"hi there, find me a laptop under 40000 birr", "en", "Computers & Laptops", "laptop"},
		{"ሰላም ስልክ እፈልጋለሁ", "am", "Mobile Phones", "ስልክ"},
		{"akkam, bilbila barbaada", "om", "Mobile Phones", "bilbila"},
	}

	for _, tc := range queries {
		plan := service.ClassifyMessage(tc.query, tc.locale, nil)
		if plan.SocialGreeting == "" {
			t.Errorf("query '%s' expected social greeting prefix, got empty", tc.query)
		}
		if !plan.RequiresDatabaseSearch {
			t.Errorf("query '%s' should require database search", tc.query)
		}
		if len(plan.Searches) == 0 {
			t.Fatalf("query '%s' expected extracted searches, got none", tc.query)
		}
		s0 := plan.Searches[0]
		if s0.CategoryName != tc.expectedCat {
			t.Errorf("expected category '%s', got '%s'", tc.expectedCat, s0.CategoryName)
		}
	}
}

// 5. Multi-Category Extraction
func TestConversational_MultiCategoryExtraction(t *testing.T) {
	query := "I want a phone and a PC and a house"
	plan := service.ClassifyMessage(query, "en", nil)

	if len(plan.Searches) != 3 {
		t.Fatalf("expected 3 distinct search requests, got %d: %+v", len(plan.Searches), plan.Searches)
	}

	expectedCategories := []string{"Mobile Phones", "Computers & Laptops", "Residential Properties & Houses"}
	for i, exp := range expectedCategories {
		if plan.Searches[i].CategoryName != exp {
			t.Errorf("index %d: expected category '%s', got '%s'", i, exp, plan.Searches[i].CategoryName)
		}
	}
}

// 6. Context Follow-Up and Pronoun Resolution
func TestConversational_ContextFollowUp(t *testing.T) {
	sm := service.NewSessionManager(10*time.Minute, 100)
	sessionID := "test_session_123"

	// Turn 1: User asks for iPhone
	sess := sm.GetOrCreate(sessionID, "en")
	plan1 := service.ClassifyMessage("show me iphone", "en", sess)
	if len(plan1.Searches) == 0 {
		t.Fatalf("Turn 1 failed: expected search for iphone")
	}
	sm.RecordTurn(sessionID, "show me iphone", "Here are iPhones", plan1, plan1.Searches)

	// Turn 2: User says "show me cheaper ones"
	sess = sm.GetOrCreate(sessionID, "en")
	plan2 := service.ClassifyMessage("show me cheaper ones", "en", sess)

	if len(plan2.Searches) == 0 {
		t.Fatalf("Turn 2 follow-up failed: expected inherited search for cheaper phone")
	}

	inherited := plan2.Searches[0]
	if inherited.CategoryName != "Mobile Phones" {
		t.Errorf("expected inherited category 'Mobile Phones', got '%s'", inherited.CategoryName)
	}

	// Turn 3: User says "only samsung"
	sess = sm.GetOrCreate(sessionID, "en")
	plan3 := service.ClassifyMessage("only samsung", "en", sess)

	if len(plan3.Searches) == 0 {
		t.Fatalf("Turn 3 brand follow-up failed: expected search for samsung")
	}
	if !strings.Contains(strings.ToLower(plan3.Searches[0].SearchTerm), "samsung") {
		t.Errorf("expected search term to contain samsung, got: %s", plan3.Searches[0].SearchTerm)
	}
}

// 7. Bounded Session Manager Limits and Sliding Expiry
func TestSessionManager_BoundedContextAndTTL(t *testing.T) {
	sm := service.NewSessionManager(50*time.Millisecond, 2) // Max 2 sessions, 50ms TTL

	// Create 2 sessions
	s1 := sm.GetOrCreate("sess_1", "en")
	s2 := sm.GetOrCreate("sess_2", "en")
	if s1 == nil || s2 == nil {
		t.Fatalf("failed to create sessions")
	}

	// Record 8 turns in sess_1 -> should cap at MaxTurnsPerSession (6)
	for i := 1; i <= 8; i++ {
		sm.RecordTurn("sess_1", "msg", "reply", nil, nil)
	}

	s1After, ok := sm.Get("sess_1")
	if !ok {
		t.Fatalf("expected sess_1 to be present")
	}
	if len(s1After.Turns) != service.MaxTurnsPerSession {
		t.Errorf("expected turns capped at %d, got %d", service.MaxTurnsPerSession, len(s1After.Turns))
	}

	// Add 3rd session -> should evict oldest session (sess_2 since sess_1 was updated with 8 turns)
	_ = sm.GetOrCreate("sess_3", "en")
	if sm.Len() > 2 {
		t.Errorf("expected max 2 sessions, got %d", sm.Len())
	}

	// Wait for TTL expiration
	time.Sleep(70 * time.Millisecond)
	_, stillExists := sm.Get("sess_1")
	if stillExists {
		t.Errorf("expected sess_1 to be expired after 70ms TTL")
	}
}

// 8. Integration test for ProcessGuestChat with Conversational Response
func TestProcessGuestChat_ConversationalPurity(t *testing.T) {
	svc := service.NewAIChatService(nil, nil, nil)
	ctx := context.Background()

	// "What is your name?"
	resp, err := svc.ProcessGuestChat(ctx, service.GuestChatRequestDTO{
		Message: "what is your name?",
		Locale:  "en",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Answer, "HUNTIQ") {
		t.Errorf("expected name to be HUNTIQ, got: %s", resp.Answer)
	}
	if len(resp.Listings) != 0 {
		t.Errorf("expected 0 listings for identity question, got %d", len(resp.Listings))
	}

	// "Hello"
	respGreeting, err := svc.ProcessGuestChat(ctx, service.GuestChatRequestDTO{
		Message: "Hello",
		Locale:  "en",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if respGreeting.GreetingPrefix == "" && !strings.Contains(respGreeting.Answer, "HUNTIQ") {
		t.Errorf("expected greeting, got: %s", respGreeting.Answer)
	}
	if len(respGreeting.Listings) != 0 {
		t.Errorf("expected 0 listings for pure hello, got %d", len(respGreeting.Listings))
	}

	// "pho" -> clarification
	respClarify, err := svc.ProcessGuestChat(ctx, service.GuestChatRequestDTO{
		Message: "pho",
		Locale:  "en",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if respClarify.Clarification == nil || !strings.Contains(*respClarify.Clarification, "phone") {
		t.Errorf("expected clarification for 'pho', got: %+v", respClarify)
	}
}
