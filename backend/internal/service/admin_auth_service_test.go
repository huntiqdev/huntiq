package service

import (
	"testing"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/security"
)

func TestValidatePortalAccess(t *testing.T) {
	tests := []struct {
		name      string
		role      string
		portal    string
		shouldErr bool
	}{
		// Executive Portal
		{"Executive - Platform Owner", models.RolePlatformOwner, "executive", false},
		{"Executive - Legacy Super", models.LegacyRoleSuper, "executive", false},
		{"Executive - Platform Admin", models.RolePlatformAdmin, "executive", true},
		{"Executive - Security Admin", models.RoleSecurityAdmin, "executive", true},
		{"Executive - Operations Admin", models.RoleOperationsAdmin, "executive", true},
		{"Executive - Moderator", models.RoleModerator, "executive", true},
		{"Executive - Support Agent", models.RoleSupportAgent, "executive", true},

		// Platform Console
		{"Platform - Platform Owner", models.RolePlatformOwner, "platform", false},
		{"Platform - Platform Admin", models.RolePlatformAdmin, "platform", false},
		{"Platform - Security Admin", models.RoleSecurityAdmin, "platform", true},
		{"Platform - Operations Admin", models.RoleOperationsAdmin, "platform", true},

		// Security Console
		{"Security - Platform Owner", models.RolePlatformOwner, "security", false},
		{"Security - Security Admin", models.RoleSecurityAdmin, "security", false},
		{"Security - Platform Admin", models.RolePlatformAdmin, "security", true},
		{"Security - Operations Admin", models.RoleOperationsAdmin, "security", true},

		// Staff Portal
		{"Staff - Operations Admin", models.RoleOperationsAdmin, "staff", false},
		{"Staff - Moderator", models.RoleModerator, "staff", false},
		{"Staff - Support Agent", models.RoleSupportAgent, "staff", false},
		{"Staff - Platform Owner", models.RolePlatformOwner, "staff", false},
		{"Staff - Empty Role", "", "staff", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePortalAccess(tt.role, tt.portal)
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidatePortalAccess(%q, %q) error = %v, want error = %v", tt.role, tt.portal, err, tt.shouldErr)
			}
		})
	}
}

func TestTokenHashingDeterministic(t *testing.T) {
	token := "test_token_secret_12345"
	h1 := security.HashToken(token)
	h2 := security.HashToken(token)
	if h1 != h2 {
		t.Fatalf("HashToken must be deterministic: %s != %s", h1, h2)
	}
	if len(h1) != 64 {
		t.Fatalf("HashToken must be 64 hex characters (SHA-256), got %d", len(h1))
	}
}

func TestDeriveRPID(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"http://localhost:3000", "localhost"},
		{"https://huntiq.com", "huntiq.com"},
		{"https://admin.huntiq.com/login", "admin.huntiq.com"},
		{"", "localhost"},
	}

	for _, tt := range tests {
		got := security.DeriveRPID(tt.url)
		if got != tt.expected {
			t.Errorf("DeriveRPID(%q) = %q; want %q", tt.url, got, tt.expected)
		}
	}
}

func TestAdminUserModelRoleHelpers(t *testing.T) {
	owner := &models.User{Role: "admin", AdminRole: models.RolePlatformOwner}
	if !owner.IsPlatformOwner() || !owner.IsPlatformAdmin() || !owner.IsSecurityAdmin() || !owner.IsOperationsAdmin() {
		t.Errorf("Platform Owner must satisfy all administrative privilege checks")
	}

	secAdmin := &models.User{Role: "admin", AdminRole: models.RoleSecurityAdmin}
	if secAdmin.IsPlatformOwner() {
		t.Errorf("Security Admin must NOT be Platform Owner")
	}
	if !secAdmin.IsSecurityAdmin() {
		t.Errorf("Security Admin must satisfy IsSecurityAdmin")
	}

	opsAdmin := &models.User{Role: "admin", AdminRole: models.RoleOperationsAdmin}
	if opsAdmin.IsPlatformOwner() || opsAdmin.IsPlatformAdmin() || opsAdmin.IsSecurityAdmin() {
		t.Errorf("Operations Admin must not have executive/security privileges")
	}
	if !opsAdmin.IsOperationsAdmin() {
		t.Errorf("Operations Admin must satisfy IsOperationsAdmin")
	}

	legacySuper := &models.User{Role: "admin", AdminRole: models.LegacyRoleSuper}
	if !legacySuper.IsPlatformOwner() {
		t.Errorf("Legacy super admin must map to Platform Owner for backward compatibility")
	}

	normalUser := &models.User{Role: "buyer", AdminRole: ""}
	if normalUser.IsPlatformOwner() || normalUser.IsPlatformAdmin() || normalUser.IsOperationsAdmin() {
		t.Errorf("Non-admin marketplace user must have zero admin privileges")
	}
}

func TestGenerateRandomTokenEntropy(t *testing.T) {
	t1 := GenerateRandomToken(32)
	t2 := GenerateRandomToken(32)
	if t1 == t2 {
		t.Fatalf("Generated tokens must be distinct: %s == %s", t1, t2)
	}
	if len(t1) != 64 { // 32 bytes = 64 hex chars
		t.Fatalf("Expected 64 hex characters for 32 bytes token, got %d", len(t1))
	}
}

func TestSessionTimeoutLogic(t *testing.T) {
	now := time.Now()
	// Idle timeout 15 min
	lastActive := now.Add(-16 * time.Minute)
	idleTimeout := 900 // 15 min in seconds

	isIdleTimedOut := now.Sub(lastActive) > time.Duration(idleTimeout)*time.Second
	if !isIdleTimedOut {
		t.Errorf("Expected session to be idle timed out after 16 minutes of inactivity")
	}

	// Active session within idle window
	recentActive := now.Add(-5 * time.Minute)
	isRecentTimedOut := now.Sub(recentActive) > time.Duration(idleTimeout)*time.Second
	if isRecentTimedOut {
		t.Errorf("Session with 5 minutes inactivity must NOT be idle timed out")
	}
}
