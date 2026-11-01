package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/your-org/used-tools-marketplace/internal/service"
)

func TestGuestAI_OutOfDomainDetection(t *testing.T) {
	outOfDomainQueries := []string{
		"solve this physics equation: E=mc^2",
		"write a python function to reverse a string",
		"write javascript code for binary search",
		"who wrote hamlet?",
		"what is the capital of france?",
		"calculate the derivative of x^3 + 2x",
		"tell me a bedtime story about dragons",
		"can you help me with my chemistry homework equation",
	}

	for _, q := range outOfDomainQueries {
		if !service.IsOutOfDomain(q) {
			t.Errorf("expected query '%s' to be detected as out of domain, but it was allowed", q)
		}
	}

	legitimateMarketplaceQueries := []string{
		"i wenna phone",
		"what products are available now",
		"I need an apartment for rent in Addis Ababa",
		"Any laptops with 16GB RAM and 512GB SSD for office use?",
		"Toyota Vitz 2015 price in Bole",
		"house for rent in CMC around 30000 birr",
		"how much is Samsung Galaxy A54?",
		"compare these two laptops for office work",
		"have you phone?",
		"have you a phone?",
		"do you have a car?",
		"ስልክ መግዛት እፈልጋለሁ",
		"ስልክ አለዎት?",
		"ስልክ አለህ?",
		"መኪና አለህ?",
		"ቤት አለዎት?",
		"bilbila bituu barbaada",
		"bilbila qabduu?",
		"ሞባይል ኣለኩም?",
	}

	for _, q := range legitimateMarketplaceQueries {
		if service.IsOutOfDomain(q) {
			t.Errorf("expected legitimate marketplace query '%s' to be allowed, but it was rejected as out of domain", q)
		}
	}
}

func TestGuestAI_PlatformFAQ(t *testing.T) {
	faqQueries := []struct {
		query   string
		mustSub string
	}{
		{"how do I post an apartment?", "Post Listing"},
		{"how to sell a car", "Post Listing"},
		{"how to register on HUNTIQ", "Register"},
		{"how does membership work?", "tier"},
	}

	for _, tc := range faqQueries {
		isFAQ, answer := service.CheckPlatformFAQ(tc.query, "en")
		if !isFAQ {
			t.Errorf("expected '%s' to be recognized as platform FAQ", tc.query)
		}
		if !strings.Contains(answer, tc.mustSub) {
			t.Errorf("expected answer for '%s' to contain '%s', got: %s", tc.query, tc.mustSub, answer)
		}
	}
}

func TestGuestAI_ProcessGuestChat_ValidationAndGuard(t *testing.T) {
	svc := service.NewAIChatService(nil, nil, nil)
	ctx := context.Background()

	// 1. Empty message error
	_, err := svc.ProcessGuestChat(ctx, service.GuestChatRequestDTO{
		Message: "   ",
		Locale:  "en",
	})
	if err == nil {
		t.Errorf("expected error for empty message, got nil")
	}

	// 2. Oversized message error (> 1000 chars)
	longMsg := strings.Repeat("A", 1005)
	_, err = svc.ProcessGuestChat(ctx, service.GuestChatRequestDTO{
		Message: longMsg,
		Locale:  "en",
	})
	if err == nil {
		t.Errorf("expected error for oversized message, got nil")
	}

	// 3. Out-of-Domain redirection (0ms guard, zero database calls)
	resp, err := svc.ProcessGuestChat(ctx, service.GuestChatRequestDTO{
		Message: "solve this physics equation: E=mc^2",
		Locale:  "en",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Answer, "marketplace assistant") {
		t.Errorf("expected domain redirection, got: %s", resp.Answer)
	}
	if len(resp.Listings) != 0 {
		t.Errorf("expected 0 listings for out of domain, got %d", len(resp.Listings))
	}
	if resp.Model != "huntiq-guard:domain-policy-v1" {
		t.Errorf("expected guard model tag, got %s", resp.Model)
	}

	// 4. Platform FAQ direct response
	faqResp, err := svc.ProcessGuestChat(ctx, service.GuestChatRequestDTO{
		Message: "how do I post an apartment?",
		Locale:  "en",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(faqResp.Answer, "Post Listing") {
		t.Errorf("expected posting guide in answer, got: %s", faqResp.Answer)
	}
	if faqResp.Model != "huntiq-kb:platform-policy-v1" {
		t.Errorf("expected kb model tag, got %s", faqResp.Model)
	}
}
