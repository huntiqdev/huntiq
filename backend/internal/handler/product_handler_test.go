package handlers

import (
	"testing"
)

func TestTierAuthorizationLogic(t *testing.T) {
	tests := []struct {
		name         string
		userTier     string
		requiredTier string
		isAdmin      bool
		expected     bool
	}{
		// Guest access
		{"Guest on normal product", "guest", "normal", false, true},
		{"Guest on empty required tier", "guest", "", false, true},
		{"Guest on gold product", "guest", "gold", false, false},
		{"Guest on premium product", "guest", "premium", false, false},

		// Registered normal buyer
		{"Registered buyer on normal product", "normal", "normal", false, true},
		{"Registered buyer on gold product", "normal", "gold", false, false},
		{"Registered buyer on premium product", "normal", "premium", false, false},

		// Gold member
		{"Gold member on normal product", "gold", "normal", false, true},
		{"Gold member on gold product", "gold", "gold", false, true},
		{"Gold member on premium product", "gold", "premium", false, false},

		// Premium member
		{"Premium member on normal product", "premium", "normal", false, true},
		{"Premium member on gold product", "premium", "gold", false, true},
		{"Premium member on premium product", "premium", "premium", false, true},

		// Admin overrides
		{"Admin on gold product without tier", "normal", "gold", true, true},
		{"Admin on premium product as guest", "guest", "premium", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTierAuthorized(tt.userTier, tt.requiredTier, tt.isAdmin)
			if got != tt.expected {
				t.Errorf("isTierAuthorized(%q, %q, %v) = %v; want %v",
					tt.userTier, tt.requiredTier, tt.isAdmin, got, tt.expected)
			}
		})
	}
}
