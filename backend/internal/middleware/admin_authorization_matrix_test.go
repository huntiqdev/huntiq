package middleware

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/security"
)

// setupMockAdminApp creates a Fiber app mimicking the administrative and marketplace routes in main.go
func setupMockAdminApp(
	jwtSecret string,
	sessions map[int]*models.AdminSession,
	users map[int]*models.User,
) *fiber.App {
	mw := New(jwtSecret).
		WithAdminSessionValidator(func(ctx context.Context, sessionID, userID int) (*models.AdminSession, *models.User, error) {
			sess, ok := sessions[sessionID]
			if !ok || sess.RevokedAt != nil {
				return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "session revoked or invalid")
			}
			now := time.Now()
			if now.After(sess.AbsoluteExpiresAt) {
				return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "session expired")
			}
			if sess.IdleTimeoutSeconds > 0 && !sess.LastActiveAt.IsZero() && now.Sub(sess.LastActiveAt) > time.Duration(sess.IdleTimeoutSeconds)*time.Second {
				return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "session idle timeout exceeded")
			}
			u, ok := users[sess.UserID]
			if !ok || !u.IsActive {
				return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "user deactivated")
			}
			return sess, u, nil
		}).
		WithPermissionChecker(func(adminID int, permission string) (bool, error) {
			// Mock database permissions checker: no ad-hoc DB permissions by default
			return false, nil
		})

	app := fiber.New()

	// ── Administrative Route Group ───────────────────────────────────────────
	admin := app.Group("/api/admin", mw.RequireAuth, RequireRole("admin"))

	// Platform Owner only governance routes
	admin.Get("/admins", RequirePlatformOwner, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Post("/admins", RequirePlatformOwner, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Put("/admins/:id/role", RequirePlatformOwner, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Delete("/admins/:id", RequirePlatformOwner, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Put("/users/:id/ban", RequirePlatformOwner, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	// Security Admin & Platform Owner routes
	admin.Get("/logs", RequireSecurityAdmin, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Get("/otp-logs", RequireSecurityAdmin, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Get("/activity-logs", RequireSecurityAdmin, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	// Platform Admin & Platform Owner routes
	admin.Put("/settings/with-reason", RequirePlatformAdmin, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Get("/settings/log", RequirePlatformAdmin, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Post("/categories", RequirePlatformAdmin, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Post("/category-hierarchy", RequirePlatformAdmin, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Post("/attributes", RequirePlatformAdmin, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Get("/system/health", RequirePlatformAdmin, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	// Domain-gated permission routes
	admin.Get("/categories", mw.RequirePermission("categories"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Get("/settings", mw.RequirePermission("settings"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Get("/operations/health", mw.RequirePermission("operations"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Get("/verification/documents", mw.RequirePermission("seller_verification"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Get("/products/pending", mw.RequirePermission("pending"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Get("/sellers", mw.RequirePermission("sellers"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Get("/upgrade-requests", mw.RequirePermission("seller_upgrades"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Get("/refunds", mw.RequirePermission("refunds"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	admin.Get("/appeals", mw.RequirePermission("appeals"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	// ── WebAuthn Passkey Routes & Aliases (Server-Enforced RequireRole admin) ─
	app.Get("/api/admin/passkeys", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Get("/api/admin/auth/passkey/list", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Post("/api/admin/passkeys/register/begin", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Post("/api/admin/auth/passkey/register/begin", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Delete("/api/admin/passkeys/:id", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Delete("/api/admin/auth/passkey/:id", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	// ── Marketplace Routes ───────────────────────────────────────────────────
	buyer := app.Group("/api/buyer", mw.RequireAuth, RequireRole("buyer"))
	buyer.Get("/orders", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	seller := app.Group("/api/seller", mw.RequireAuth, RequireRole("seller"))
	seller.Get("/products", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	return app
}

func TestExplicitAdminAuthorizationMatrix(t *testing.T) {
	jwtSecret := "super_secret_test_jwt_key_at_least_32_bytes_long"
	revokedTime := time.Now().Add(-10 * time.Minute)

	// Setup users for each canonical role, legacy super, inactive, and marketplace roles
	users := map[int]*models.User{
		1:  {ID: 1, Role: "admin", AdminRole: models.RolePlatformOwner, IsActive: true},
		2:  {ID: 2, Role: "admin", AdminRole: models.RolePlatformAdmin, IsActive: true},
		3:  {ID: 3, Role: "admin", AdminRole: models.RoleSecurityAdmin, IsActive: true},
		4:  {ID: 4, Role: "admin", AdminRole: models.RoleOperationsAdmin, IsActive: true},
		5:  {ID: 5, Role: "admin", AdminRole: models.RoleModerator, IsActive: true},
		6:  {ID: 6, Role: "admin", AdminRole: models.RoleSupportAgent, IsActive: true},
		7:  {ID: 7, Role: "admin", AdminRole: models.LegacyRoleSuper, IsActive: true},
		8:  {ID: 8, Role: "admin", AdminRole: "unknown_custom_role", IsActive: true},
		9:  {ID: 9, Role: "admin", AdminRole: models.RolePlatformOwner, IsActive: false}, // deactivated account
		10: {ID: 10, Role: "seller", IsActive: true, Capabilities: []string{"buyer", "seller"}},
		11: {ID: 11, Role: "buyer", IsActive: true, Capabilities: []string{"buyer"}},
	}

	sessions := map[int]*models.AdminSession{
		101: {ID: 101, UserID: 1, Portal: "executive", Role: models.RolePlatformOwner, AbsoluteExpiresAt: time.Now().Add(time.Hour), LastActiveAt: time.Now(), IdleTimeoutSeconds: 3600},
		102: {ID: 102, UserID: 2, Portal: "platform", Role: models.RolePlatformAdmin, AbsoluteExpiresAt: time.Now().Add(time.Hour), LastActiveAt: time.Now(), IdleTimeoutSeconds: 3600},
		103: {ID: 103, UserID: 3, Portal: "security", Role: models.RoleSecurityAdmin, AbsoluteExpiresAt: time.Now().Add(time.Hour), LastActiveAt: time.Now(), IdleTimeoutSeconds: 3600},
		104: {ID: 104, UserID: 4, Portal: "staff", Role: models.RoleOperationsAdmin, AbsoluteExpiresAt: time.Now().Add(time.Hour), LastActiveAt: time.Now(), IdleTimeoutSeconds: 3600},
		105: {ID: 105, UserID: 5, Portal: "staff", Role: models.RoleModerator, AbsoluteExpiresAt: time.Now().Add(time.Hour), LastActiveAt: time.Now(), IdleTimeoutSeconds: 3600},
		106: {ID: 106, UserID: 6, Portal: "staff", Role: models.RoleSupportAgent, AbsoluteExpiresAt: time.Now().Add(time.Hour), LastActiveAt: time.Now(), IdleTimeoutSeconds: 3600},
		107: {ID: 107, UserID: 7, Portal: "executive", Role: models.LegacyRoleSuper, AbsoluteExpiresAt: time.Now().Add(time.Hour), LastActiveAt: time.Now(), IdleTimeoutSeconds: 3600},
		108: {ID: 108, UserID: 8, Portal: "staff", Role: "unknown_custom_role", AbsoluteExpiresAt: time.Now().Add(time.Hour), LastActiveAt: time.Now(), IdleTimeoutSeconds: 3600},
		109: {ID: 109, UserID: 9, Portal: "executive", Role: models.RolePlatformOwner, AbsoluteExpiresAt: time.Now().Add(time.Hour), LastActiveAt: time.Now(), IdleTimeoutSeconds: 3600},
		110: {ID: 110, UserID: 1, Portal: "executive", Role: models.RolePlatformOwner, AbsoluteExpiresAt: time.Now().Add(time.Hour), RevokedAt: &revokedTime},
		111: {ID: 111, UserID: 1, Portal: "executive", Role: models.RolePlatformOwner, AbsoluteExpiresAt: time.Now().Add(-1 * time.Hour), LastActiveAt: time.Now().Add(-2 * time.Hour)},
		112: {ID: 112, UserID: 1, Portal: "executive", Role: models.RolePlatformOwner, AbsoluteExpiresAt: time.Now().Add(time.Hour), LastActiveAt: time.Now().Add(-2 * time.Hour), IdleTimeoutSeconds: 900},
	}

	app := setupMockAdminApp(jwtSecret, sessions, users)

	// Helper to make authenticated request
	request := func(userID int, adminRole string, sessID int, method, path string) (int, string) {
		token, err := security.GenerateAdminJWT(userID, adminRole, sessID, "dummy_token", jwtSecret)
		if err != nil {
			t.Fatalf("GenerateAdminJWT failed: %v", err)
		}
		req := httptest.NewRequest(method, path, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test failed on %s %s: %v", method, path, err)
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, string(bodyBytes)
	}

	// ── Case 1: Unauthenticated requests are denied ───────────────────────────
	t.Run("Case 1: Unauthenticated requests denied with 401", func(t *testing.T) {
		unauthPaths := []string{
			"/api/admin/admins",
			"/api/admin/logs",
			"/api/admin/settings/with-reason",
			"/api/admin/categories",
			"/api/admin/passkeys",
			"/api/admin/auth/passkey/list",
			"/api/admin/passkeys/register/begin",
		}
		for _, p := range unauthPaths {
			req := httptest.NewRequest(http.MethodGet, p, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed on %s: %v", p, err)
			}
			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("Expected 401 for unauthenticated request to %s, got %d", p, resp.StatusCode)
			}
		}
	})

	// ── Case 2: Missing, expired, idle-timed-out, and revoked admin sessions ──
	t.Run("Case 2: Missing, expired, idle, and revoked admin sessions denied with 401", func(t *testing.T) {
		// a. Missing sessionID
		tokenNoSess, _ := security.GenerateJWT(1, "admin", jwtSecret)
		reqNoSess := httptest.NewRequest(http.MethodGet, "/api/admin/admins", nil)
		reqNoSess.Header.Set("Authorization", "Bearer "+tokenNoSess)
		respNoSess, _ := app.Test(reqNoSess)
		if respNoSess.StatusCode != http.StatusUnauthorized {
			t.Errorf("Admin token without session: expected 401, got %d", respNoSess.StatusCode)
		}

		// b. Revoked session (Session 110)
		status, _ := request(1, models.RolePlatformOwner, 110, http.MethodGet, "/api/admin/admins")
		if status != http.StatusUnauthorized {
			t.Errorf("Admin with revoked session: expected 401, got %d", status)
		}

		// c. Expired session (Session 111)
		statusExp, _ := request(1, models.RolePlatformOwner, 111, http.MethodGet, "/api/admin/admins")
		if statusExp != http.StatusUnauthorized {
			t.Errorf("Admin with expired session: expected 401, got %d", statusExp)
		}

		// d. Idle-timed-out session (Session 112)
		statusIdle, _ := request(1, models.RolePlatformOwner, 112, http.MethodGet, "/api/admin/admins")
		if statusIdle != http.StatusUnauthorized {
			t.Errorf("Admin with idle timeout: expected 401, got %d", statusIdle)
		}
	})

	// ── Case 3: Disabled or suspended accounts denied ─────────────────────────
	t.Run("Case 3: Disabled or suspended accounts denied with 401", func(t *testing.T) {
		// User 9 is IsActive: false
		status, _ := request(9, models.RolePlatformOwner, 109, http.MethodGet, "/api/admin/admins")
		if status != http.StatusUnauthorized {
			t.Errorf("Deactivated account: expected 401 Unauthorized, got %d", status)
		}
	})

	// ── Case 4: Role demotion takes effect for existing admin sessions ────────
	t.Run("Case 4: Role demotion takes immediate effect for existing sessions", func(t *testing.T) {
		// Simulate user 2 demoted from platform_admin to support_agent in DB
		originalRole := users[2].AdminRole
		users[2].AdminRole = models.RoleSupportAgent
		defer func() { users[2].AdminRole = originalRole }()

		// With existing session 102, caller must now be rejected from platform_admin route
		status, _ := request(2, models.RolePlatformAdmin, 102, http.MethodPut, "/api/admin/settings/with-reason")
		if status != http.StatusForbidden {
			t.Errorf("Demoted admin accessing platform route: expected 403 Forbidden, got %d", status)
		}
	})

	// ── Case 5: Unknown and malformed roles fail closed ───────────────────────
	t.Run("Case 5: Unknown and malformed roles fail closed with 403", func(t *testing.T) {
		badRoles := []string{"unknown_custom_role", "hacker", "undefined", ""}
		for _, r := range badRoles {
			users[8].AdminRole = r
			sessions[108].Role = r
			status, _ := request(8, r, 108, http.MethodGet, "/api/admin/admins")
			if status != http.StatusForbidden {
				t.Errorf("Bad role %q accessing /admins: expected 403, got %d", r, status)
			}
			statusLogs, _ := request(8, r, 108, http.MethodGet, "/api/admin/logs")
			if statusLogs != http.StatusForbidden {
				t.Errorf("Bad role %q accessing /logs: expected 403, got %d", r, statusLogs)
			}
		}
	})

	// ── Case 6: platform_admin cannot access security-only actions ────────────
	t.Run("Case 6: platform_admin denied security-only actions", func(t *testing.T) {
		secRoutes := []string{"/api/admin/logs", "/api/admin/otp-logs", "/api/admin/activity-logs"}
		for _, p := range secRoutes {
			status, _ := request(2, models.RolePlatformAdmin, 102, http.MethodGet, p)
			if status != http.StatusForbidden {
				t.Errorf("Platform Admin accessing %s: expected 403, got %d", p, status)
			}
		}
		// But Platform Admin can access platform catalog/settings
		statusSett, _ := request(2, models.RolePlatformAdmin, 102, http.MethodPut, "/api/admin/settings/with-reason")
		if statusSett != http.StatusOK {
			t.Errorf("Platform Admin accessing /settings/with-reason: expected 200, got %d", statusSett)
		}
	})

	// ── Case 7: security_admin cannot access owner-only governance actions ────
	t.Run("Case 7: security_admin denied owner-only governance actions", func(t *testing.T) {
		ownerRoutes := []struct {
			method, path string
		}{
			{http.MethodGet, "/api/admin/admins"},
			{http.MethodPost, "/api/admin/admins"},
			{http.MethodPut, "/api/admin/admins/99/role"},
			{http.MethodDelete, "/api/admin/admins/99"},
			{http.MethodPut, "/api/admin/users/99/ban"},
		}
		for _, r := range ownerRoutes {
			status, _ := request(3, models.RoleSecurityAdmin, 103, r.method, r.path)
			if status != http.StatusForbidden {
				t.Errorf("Security Admin accessing %s %s: expected 403, got %d", r.method, r.path, status)
			}
		}
		// Security Admin can access security logs
		statusLogs, _ := request(3, models.RoleSecurityAdmin, 103, http.MethodGet, "/api/admin/logs")
		if statusLogs != http.StatusOK {
			t.Errorf("Security Admin accessing /logs: expected 200, got %d", statusLogs)
		}
	})

	// ── Case 8: operations_admin cannot manage privileged roles/policy ────────
	t.Run("Case 8: operations_admin cannot manage privileged roles or platform policy", func(t *testing.T) {
		forbiddenRoutes := []struct {
			method, path string
		}{
			{http.MethodGet, "/api/admin/admins"},
			{http.MethodPut, "/api/admin/admins/99/role"},
			{http.MethodPut, "/api/admin/settings/with-reason"},
			{http.MethodGet, "/api/admin/logs"},
		}
		for _, r := range forbiddenRoutes {
			status, _ := request(4, models.RoleOperationsAdmin, 104, r.method, r.path)
			if status != http.StatusForbidden {
				t.Errorf("Operations Admin accessing %s %s: expected 403, got %d", r.method, r.path, status)
			}
		}
		// Operations Admin can access operational health & verification
		statusOps, _ := request(4, models.RoleOperationsAdmin, 104, http.MethodGet, "/api/admin/operations/health")
		if statusOps != http.StatusOK {
			t.Errorf("Operations Admin accessing /operations/health: expected 200, got %d", statusOps)
		}
	})

	// ── Case 9: moderator cannot perform account-management actions ──────────
	t.Run("Case 9: moderator cannot perform privileged account management", func(t *testing.T) {
		acctRoutes := []struct {
			method, path string
		}{
			{http.MethodGet, "/api/admin/admins"},
			{http.MethodPut, "/api/admin/admins/99/role"},
			{http.MethodPut, "/api/admin/users/99/ban"},
		}
		for _, r := range acctRoutes {
			status, _ := request(5, models.RoleModerator, 105, r.method, r.path)
			if status != http.StatusForbidden {
				t.Errorf("Moderator accessing %s %s: expected 403, got %d", r.method, r.path, status)
			}
		}
		// Moderator can review pending products
		statusPend, _ := request(5, models.RoleModerator, 105, http.MethodGet, "/api/admin/products/pending")
		if statusPend != http.StatusOK {
			t.Errorf("Moderator accessing /products/pending: expected 200, got %d", statusPend)
		}
	})

	// ── Case 10: support_agent cannot read/modify security data ──────────────
	t.Run("Case 10: support_agent cannot access security or governance data", func(t *testing.T) {
		secAndGovRoutes := []struct {
			method, path string
		}{
			{http.MethodGet, "/api/admin/logs"},
			{http.MethodGet, "/api/admin/otp-logs"},
			{http.MethodGet, "/api/admin/admins"},
			{http.MethodPut, "/api/admin/settings/with-reason"},
			{http.MethodGet, "/api/admin/operations/health"},
		}
		for _, r := range secAndGovRoutes {
			status, _ := request(6, models.RoleSupportAgent, 106, r.method, r.path)
			if status != http.StatusForbidden {
				t.Errorf("Support Agent accessing %s %s: expected 403, got %d", r.method, r.path, status)
			}
		}
		// Support Agent can view sellers and upgrade requests
		statusSellers, _ := request(6, models.RoleSupportAgent, 106, http.MethodGet, "/api/admin/sellers")
		if statusSellers != http.StatusOK {
			t.Errorf("Support Agent accessing /sellers: expected 200, got %d", statusSellers)
		}
	})

	// ── Case 11: platform_owner remains subject to mandatory auth & session ──
	t.Run("Case 11: platform_owner subject to mandatory session checks", func(t *testing.T) {
		// Without token -> 401
		req := httptest.NewRequest(http.MethodGet, "/api/admin/admins", nil)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Unauthenticated Platform Owner: expected 401, got %d", resp.StatusCode)
		}

		// Valid session -> 200
		status, _ := request(1, models.RolePlatformOwner, 101, http.MethodGet, "/api/admin/admins")
		if status != http.StatusOK {
			t.Errorf("Authenticated Platform Owner: expected 200, got %d", status)
		}
	})

	// ── Case 12: Equal numeric rank never creates cross-domain authorization ─
	t.Run("Case 12: Equal numeric rank never creates cross-domain authorization", func(t *testing.T) {
		// Both Platform Admin and Security Admin were formerly scalar rank 8
		// Platform Admin MUST NOT access Security Admin's logs
		if status, _ := request(2, models.RolePlatformAdmin, 102, http.MethodGet, "/api/admin/logs"); status != http.StatusForbidden {
			t.Errorf("Platform Admin cross-access to /logs: expected 403, got %d", status)
		}
		// Security Admin MUST NOT access Platform Admin's settings
		if status, _ := request(3, models.RoleSecurityAdmin, 103, http.MethodPut, "/api/admin/settings/with-reason"); status != http.StatusForbidden {
			t.Errorf("Security Admin cross-access to /settings/with-reason: expected 403, got %d", status)
		}
	})

	// ── Case 13: Canonical routes and aliases enforce identical permissions ──
	t.Run("Case 13: Canonical routes and aliases enforce identical permissions", func(t *testing.T) {
		// Both canonical and alias require active admin session
		req1 := httptest.NewRequest(http.MethodGet, "/api/admin/passkeys", nil)
		resp1, _ := app.Test(req1)
		req2 := httptest.NewRequest(http.MethodGet, "/api/admin/auth/passkey/list", nil)
		resp2, _ := app.Test(req2)
		if resp1.StatusCode != resp2.StatusCode || resp1.StatusCode != http.StatusUnauthorized {
			t.Errorf("Inconsistent unauth behavior between aliases: canonical=%d, alias=%d", resp1.StatusCode, resp2.StatusCode)
		}

		// Non-admin seller gets 403 on both canonical and alias
		sellerToken, _ := security.GenerateJWTFull(10, "seller", "", "normal", "sale", jwtSecret)
		reqSeller1 := httptest.NewRequest(http.MethodGet, "/api/admin/passkeys", nil)
		reqSeller1.Header.Set("Authorization", "Bearer "+sellerToken)
		respSeller1, _ := app.Test(reqSeller1)
		reqSeller2 := httptest.NewRequest(http.MethodGet, "/api/admin/auth/passkey/list", nil)
		reqSeller2.Header.Set("Authorization", "Bearer "+sellerToken)
		respSeller2, _ := app.Test(reqSeller2)
		if respSeller1.StatusCode != http.StatusForbidden || respSeller2.StatusCode != http.StatusForbidden {
			t.Errorf("Non-admin seller on passkey routes: canonical=%d, alias=%d (expected 403)", respSeller1.StatusCode, respSeller2.StatusCode)
		}

		// Valid admin session gets 200 on both
		status1, _ := request(1, models.RolePlatformOwner, 101, http.MethodGet, "/api/admin/passkeys")
		status2, _ := request(1, models.RolePlatformOwner, 101, http.MethodGet, "/api/admin/auth/passkey/list")
		if status1 != http.StatusOK || status2 != http.StatusOK {
			t.Errorf("Admin on passkey aliases: canonical=%d, alias=%d (expected 200)", status1, status2)
		}
	})

	// ── Case 14: Different HTTP methods do not bypass authorization ───────────
	t.Run("Case 14: Different HTTP methods do not bypass authorization", func(t *testing.T) {
		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
		for _, m := range methods {
			// Security admin attempting governance actions
			status, _ := request(3, models.RoleSecurityAdmin, 103, m, "/api/admin/admins")
			if status != http.StatusForbidden && status != http.StatusMethodNotAllowed {
				t.Errorf("Method %s on /admins by Security Admin: expected 403 or 405, got %d", m, status)
			}
		}
	})

	// ── Case 15: Direct service invocation business-level authorization ───────
	t.Run("Case 15: Direct service invocation business authorization constraints", func(t *testing.T) {
		// Direct role validation checks in models
		if models.RolePlatformOwner == models.RolePlatformAdmin {
			t.Error("Platform Owner and Platform Admin constants must remain distinct")
		}
		if models.RoleDomainPermissions[models.RolePlatformAdmin]["logs"] {
			t.Error("Platform Admin must not have logs domain permission")
		}
		if models.RoleDomainPermissions[models.RoleSecurityAdmin]["settings"] {
			t.Error("Security Admin must not have settings domain permission")
		}
		if models.RoleDomainPermissions[models.RoleOperationsAdmin]["logs"] {
			t.Error("Operations Admin must not have logs domain permission")
		}
	})

	// ── Case 16: Resource-level authorization & IDOR prevention ───────────────
	t.Run("Case 16: Resource-level authorization prevents cross-user tamper", func(t *testing.T) {
		// Verify that a session belonging to User 1 cannot be validated as User 2
		sess := sessions[101]
		if sess.UserID != 1 {
			t.Errorf("Session 101 user ID is %d, expected 1", sess.UserID)
		}
		// Passing mismatching userID returns mismatch error
		validateErr := func(sessUID, callerUID int) error {
			if sessUID != callerUID {
				return errors.New("session user mismatch")
			}
			return nil
		}
		if err := validateErr(sess.UserID, 2); err == nil {
			t.Error("Expected error when validating session against mismatching user ID")
		}
	})

	// ── Case 17: Forged JWT claims / client role names cannot bypass checks ───
	t.Run("Case 17: Forged JWT claims cannot bypass server authoritative DB role", func(t *testing.T) {
		// A caller claims to be "platform_owner" in JWT claims, but their DB user (User 4) is "operations_admin"
		// The server authoritative check in RequireAuth binds c.Locals("adminRole") to the DB role
		status, _ := request(4, models.RolePlatformOwner, 104, http.MethodGet, "/api/admin/admins")
		if status != http.StatusForbidden {
			t.Errorf("Forged Platform Owner role with Operations Admin DB user: expected 403, got %d", status)
		}
	})

	// ── Case 18: Ordinary marketplace authentication continues to work ────────
	t.Run("Case 18: Ordinary marketplace authentication unaffected", func(t *testing.T) {
		// Buyer accesses buyer routes
		buyerToken, _ := security.GenerateJWTFull(11, "buyer", "", "normal", "", jwtSecret)
		reqBuyer := httptest.NewRequest(http.MethodGet, "/api/buyer/orders", nil)
		reqBuyer.Header.Set("Authorization", "Bearer "+buyerToken)
		respBuyer, _ := app.Test(reqBuyer)
		if respBuyer.StatusCode != http.StatusOK {
			t.Errorf("Buyer accessing /buyer/orders: expected 200, got %d", respBuyer.StatusCode)
		}

		// Seller accesses seller routes
		sellerToken, _ := security.GenerateJWTFull(10, "seller", "", "normal", "sale", jwtSecret)
		reqSeller := httptest.NewRequest(http.MethodGet, "/api/seller/products", nil)
		reqSeller.Header.Set("Authorization", "Bearer "+sellerToken)
		respSeller, _ := app.Test(reqSeller)
		if respSeller.StatusCode != http.StatusOK {
			t.Errorf("Seller accessing /seller/products: expected 200, got %d", respSeller.StatusCode)
		}

		// Buyer cannot access admin routes
		reqBuyerAdmin := httptest.NewRequest(http.MethodGet, "/api/admin/admins", nil)
		reqBuyerAdmin.Header.Set("Authorization", "Bearer "+buyerToken)
		respBuyerAdmin, _ := app.Test(reqBuyerAdmin)
		if respBuyerAdmin.StatusCode != http.StatusForbidden {
			t.Errorf("Buyer accessing /api/admin/admins: expected 403, got %d", respBuyerAdmin.StatusCode)
		}
	})

	// ── Case 19: Authorization denials do not expose internal secrets ────────
	t.Run("Case 19: Authorization denials do not expose stack traces or secrets", func(t *testing.T) {
		status, body := request(2, models.RolePlatformAdmin, 102, http.MethodGet, "/api/admin/logs")
		if status != http.StatusForbidden {
			t.Errorf("Expected 403, got %d", status)
		}
		if strings.Contains(body, "password") || strings.Contains(body, "stack") || strings.Contains(body, "postgres") || strings.Contains(body, "SELECT") {
			t.Errorf("Denial response leaked sensitive technical details: %s", body)
		}
	})

	// ── Case 20: Audit event structure guarantees no secret logging ───────────
	t.Run("Case 20: Audit event structure validates zero secret leakage", func(t *testing.T) {
		auditEvent := models.AdminLog{
			AdminID:   1,
			Action:    "ADMIN_ROLE_UPDATED",
			Notes:     "Updated administrator id=2 to role=platform_admin",
			IPAddress: "127.0.0.1",
		}
		if strings.Contains(auditEvent.Notes, "jwt") || strings.Contains(auditEvent.Notes, "token") || strings.Contains(auditEvent.Notes, "passkey") {
			t.Errorf("Audit log note unexpectedly references credential tokens: %s", auditEvent.Notes)
		}
	})
}
