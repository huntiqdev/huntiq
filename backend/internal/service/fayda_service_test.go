package service

import (
	"context"
	"strings"
	"testing"
)

func TestMockFaydaDevProvider(t *testing.T) {
	provider := &MockFaydaDevProvider{}
	ctx := context.Background()

	t.Run("Valid 12-digit FIN succeeds with simulated official eKYC verification", func(t *testing.T) {
		req := &FaydaVerificationRequest{
			UserID:       101,
			FIN:          "123456789012",
			FullName:     "Abebe Kebede Bikila",
			Phone:        "+251911223344",
			ConsentGiven: true,
		}

		res, err := provider.VerifyIdentity(ctx, req)
		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}

		if !res.IsValid {
			t.Fatalf("expected res.IsValid to be true")
		}
		if provider.ProviderCode() != "fayda_dev_mock" {
			t.Fatalf("unexpected provider code: %s", provider.ProviderCode())
		}
		if !strings.HasPrefix(res.AuthReference, "FAYDA-AUTH-") {
			t.Fatalf("unexpected auth reference: %s", res.AuthReference)
		}
		if res.PhoneMatched {
			t.Fatalf("phone_matched must be false under privacy-preserving Level 1 eKYC")
		}
		if res.PhoneMatchStatus != "unsupported_by_provider" {
			t.Fatalf("phone_match_status must be 'unsupported_by_provider', got %s", res.PhoneMatchStatus)
		}
	})

	t.Run("Rejects FIN with non-numeric characters", func(t *testing.T) {
		req := &FaydaVerificationRequest{
			UserID:       102,
			FIN:          "1234ABCD9012",
			FullName:     "Abebe Kebede",
			ConsentGiven: true,
		}
		_, err := provider.VerifyIdentity(ctx, req)
		if err == nil {
			t.Fatalf("expected error for non-digit FIN")
		}
	})

	t.Run("Rejects FIN with invalid length", func(t *testing.T) {
		req := &FaydaVerificationRequest{
			UserID:       103,
			FIN:          "12345",
			FullName:     "Abebe Kebede",
			ConsentGiven: true,
		}
		_, err := provider.VerifyIdentity(ctx, req)
		if err == nil {
			t.Fatalf("expected error for short FIN")
		}
	})

	t.Run("Rejects empty legal name", func(t *testing.T) {
		req := &FaydaVerificationRequest{
			UserID:       104,
			FIN:          "123456789012",
			FullName:     "   ",
			ConsentGiven: true,
		}
		_, err := provider.VerifyIdentity(ctx, req)
		if err == nil {
			t.Fatalf("expected error for empty legal name")
		}
	})

	t.Run("Rejects without explicit user consent", func(t *testing.T) {
		req := &FaydaVerificationRequest{
			UserID:       105,
			FIN:          "123456789012",
			FullName:     "Abebe Kebede",
			ConsentGiven: false,
		}
		_, err := provider.VerifyIdentity(ctx, req)
		if err == nil {
			t.Fatalf("expected error for missing consent")
		}
	})
}
