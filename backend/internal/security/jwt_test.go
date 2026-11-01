package security_test

import (
	"testing"
	"github.com/your-org/used-tools-marketplace/internal/security"
)

func TestTierHierarchyAndJWTClaims(t *testing.T) {
	secret := "test-secret-at-least-32-chars-long-security"

	// 1. Generate JWT for Gold Buyer
	token, err := security.GenerateJWTWithTier(37, "buyer", "gold", secret)
	if err != nil {
		t.Fatalf("Failed to generate Gold JWT: %v", err)
	}

	claims, err := security.ParseJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to parse Gold JWT: %v", err)
	}

	if claims.UserID != 37 {
		t.Errorf("expected UserID 37, got %d", claims.UserID)
	}
	if claims.BuyerTier != "gold" {
		t.Errorf("expected BuyerTier 'gold', got %q", claims.BuyerTier)
	}

	// 2. Generate JWT for Normal Registered Buyer
	normToken, err := security.GenerateJWTWithTier(41, "buyer", "normal", secret)
	if err != nil {
		t.Fatalf("Failed to generate Normal JWT: %v", err)
	}

	normClaims, err := security.ParseJWT(normToken, secret)
	if err != nil {
		t.Fatalf("Failed to parse Normal JWT: %v", err)
	}

	if normClaims.BuyerTier != "normal" {
		t.Errorf("expected BuyerTier 'normal', got %q", normClaims.BuyerTier)
	}

	// 3. Generate JWT with Profile (Name, Phone, Capabilities)
	profileToken, err := security.GenerateJWTWithProfile(55, "Abebe Bikila", "+251911223344", "seller", "", "normal", "sale", []string{"buyer", "seller"}, secret)
	if err != nil {
		t.Fatalf("Failed to generate Profile JWT: %v", err)
	}

	profileClaims, err := security.ParseJWT(profileToken, secret)
	if err != nil {
		t.Fatalf("Failed to parse Profile JWT: %v", err)
	}

	if profileClaims.Name != "Abebe Bikila" {
		t.Errorf("expected Name 'Abebe Bikila', got %q", profileClaims.Name)
	}
	if profileClaims.Phone != "+251911223344" {
		t.Errorf("expected Phone '+251911223344', got %q", profileClaims.Phone)
	}
	if !profileClaims.HasCapability("seller") {
		t.Errorf("expected seller capability to be true")
	}
}
