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

// 9. Housing Ambiguity & Multi-Turn Slot Resolution
func TestConversational_HousingAmbiguityAndSlotResolution(t *testing.T) {
	// Turn 1: "I want a house" -> should ask whether user wants to buy or rent
	plan1 := service.ClassifyMessage("I want a house", "en", nil)
	if !plan1.RequiresClarification {
		t.Errorf("expected 'I want a house' to require clarification (buy vs rent)")
	}
	if plan1.UnresolvedSlot != "listing_type" {
		t.Errorf("expected unresolved slot 'listing_type', got '%s'", plan1.UnresolvedSlot)
	}
	if !strings.Contains(plan1.ClarificationQuestion, "buy or rent") {
		t.Errorf("expected buy or rent clarification question, got: %s", plan1.ClarificationQuestion)
	}

	// Prepare session context with unresolved slot
	sessionCtx := &service.ConversationContext{
		UnresolvedSlot: plan1.UnresolvedSlot,
		PendingSearch:  plan1.PendingSearch,
	}

	// Turn 2: User answers "Rent" -> slot resolved, search executes!
	plan2 := service.ClassifyMessage("Rent", "en", sessionCtx)
	if plan2.RequiresClarification {
		t.Errorf("expected turn 2 'Rent' to resolve clarification, but required clarification again")
	}
	if len(plan2.Searches) != 1 {
		t.Fatalf("expected 1 search in turn 2, got %d", len(plan2.Searches))
	}
	if plan2.Searches[0].ListingType == nil || *plan2.Searches[0].ListingType != "rent" {
		t.Errorf("expected listing_type 'rent', got %+v", plan2.Searches[0].ListingType)
	}
}

// 10. Relative Location Clarification
func TestConversational_RelativeLocationClarification(t *testing.T) {
	plan := service.ClassifyMessage("Find a house near my workplace", "en", nil)
	if !plan.RequiresClarification {
		t.Errorf("expected 'near my workplace' to trigger location clarification")
	}
	if plan.UnresolvedSlot != "location" {
		t.Errorf("expected unresolved slot 'location', got '%s'", plan.UnresolvedSlot)
	}
	if !strings.Contains(plan.ClarificationQuestion, "workplace") && !strings.Contains(plan.ClarificationQuestion, "city") {
		t.Errorf("expected workplace location question, got: %s", plan.ClarificationQuestion)
	}
}

// 11. Conversational Intent Correction: Buy -> Rent
func TestConversational_IntentCorrection(t *testing.T) {
	saleType := "sale"
	sessionCtx := &service.ConversationContext{
		RecentSearches: []service.ExtractedSearchRequest{
			{
				CategoryName: "Residential Properties & Houses",
				SearchTerm:   "house",
				ListingType:  &saleType,
			},
		},
	}

	plan := service.ClassifyMessage("I don't want to buy, I want to rent", "en", sessionCtx)
	if len(plan.Searches) != 1 {
		t.Fatalf("expected 1 corrected search, got %d", len(plan.Searches))
	}
	if plan.Searches[0].ListingType == nil || *plan.Searches[0].ListingType != "rent" {
		t.Errorf("expected corrected listing_type 'rent', got %+v", plan.Searches[0].ListingType)
	}
}

// 12. Brand Filter Refinement: "Only Samsung"
func TestConversational_BrandFilterRefinement(t *testing.T) {
	sessionCtx := &service.ConversationContext{
		RecentSearches: []service.ExtractedSearchRequest{
			{
				CategoryName: "Mobile Phones",
				SearchTerm:   "phone",
			},
		},
	}

	plan := service.ClassifyMessage("Only Samsung", "en", sessionCtx)
	if len(plan.Searches) != 1 {
		t.Fatalf("expected 1 search, got %d", len(plan.Searches))
	}
	if plan.Searches[0].Brand != "samsung" {
		t.Errorf("expected brand 'samsung', got '%s'", plan.Searches[0].Brand)
	}
}

// 13. Budget and Food Extraction: "I have 4,000 birr, where can I get food?"
func TestConversational_BudgetAndFoodExtraction(t *testing.T) {
	plan := service.ClassifyMessage("I have 4,000 birr, where can I get food?", "en", nil)
	if len(plan.Searches) != 1 {
		t.Fatalf("expected 1 search for food, got %d", len(plan.Searches))
	}
	if plan.Searches[0].CategoryName != "Food & Groceries" {
		t.Errorf("expected 'Food & Groceries', got '%s'", plan.Searches[0].CategoryName)
	}
	if plan.Searches[0].MaxPrice == nil || *plan.Searches[0].MaxPrice != 4000.0 {
		t.Errorf("expected max_price 4000.0, got %+v", plan.Searches[0].MaxPrice)
	}
}

// 14. Full Constraints Extraction: "I need a two-bedroom house in Bole below 20k monthly"
func TestConversational_FullConstraintsNoUnnecessaryClarification(t *testing.T) {
	plan := service.ClassifyMessage("I need a two-bedroom house in Bole below 20k monthly", "en", nil)
	if plan.RequiresClarification {
		t.Errorf("should NOT require clarification when all constraints (house, Bole, 20k) are provided")
	}
	if len(plan.Searches) != 1 {
		t.Fatalf("expected 1 search, got %d", len(plan.Searches))
	}
	s := plan.Searches[0]
	if s.Location != "Bole" {
		t.Errorf("expected location 'Bole', got '%s'", s.Location)
	}
	if s.MaxPrice == nil || *s.MaxPrice != 20000.0 {
		t.Errorf("expected max_price 20000.0, got %+v", s.MaxPrice)
	}
}

// 15. Expanded Platform FAQs
func TestConversational_ExpandedPlatformFAQs(t *testing.T) {
	faqs := []struct {
		query   string
		locale  string
		mustSub string
	}{
		{"Can I register as a seller?", "en", "seller"},
		{"How do I publish a listing?", "en", "Post Listing"},
		{"Is this listing verified?", "en", "Verified"},
		{"I want to sell my land", "en", "Land & Plots"},
		{"How much does it cost to list?", "en", "free"},
		{"How do I pay with Telebirr?", "en", "Telebirr"},
	}

	for _, tc := range faqs {
		isFAQ, ans := service.CheckPlatformFAQ(tc.query, tc.locale)
		if !isFAQ {
			t.Errorf("query '%s' should be recognized as platform FAQ", tc.query)
		}
		if !strings.Contains(ans, tc.mustSub) {
			t.Errorf("query '%s' answer should contain '%s', got: %s", tc.query, tc.mustSub, ans)
		}
	}
}
