package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/security"
)

func TestRequirePortalMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		userRole   string
		adminRole  string
		portal     string
		wantStatus int
	}{
		{"Non-admin blocked", "buyer", "", "executive", fiber.StatusForbidden},
		{"Platform Owner allowed executive", "admin", models.RolePlatformOwner, "executive", fiber.StatusOK},
		{"Legacy Super allowed executive", "admin", models.LegacyRoleSuper, "executive", fiber.StatusOK},
		{"Platform Admin blocked executive", "admin", models.RolePlatformAdmin, "executive", fiber.StatusForbidden},
		{"Platform Admin allowed platform console", "admin", models.RolePlatformAdmin, "platform", fiber.StatusOK},
		{"Operations Admin blocked platform console", "admin", models.RoleOperationsAdmin, "platform", fiber.StatusForbidden},
		{"Security Admin allowed security console", "admin", models.RoleSecurityAdmin, "security", fiber.StatusOK},
		{"Platform Admin blocked security console", "admin", models.RolePlatformAdmin, "security", fiber.StatusForbidden},
		{"Operations Admin allowed staff", "admin", models.RoleOperationsAdmin, "staff", fiber.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/test", func(c *fiber.Ctx) error {
				c.Locals("role", tt.userRole)
				c.Locals("adminRole", tt.adminRole)
				return c.Next()
			}, RequirePortal(tt.portal), func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Test request failed: %v", err)
			}
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("RequirePortal status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestRequirePermissionOrSuper(t *testing.T) {
	mw := New("test_secret").WithPermissionChecker(func(adminID int, permission string) (bool, error) {
		return permission == "pending", nil
	})

	tests := []struct {
		name       string
		adminRole  string
		adminID    int
		perm       string
		wantStatus int
	}{
		{"Platform Owner bypasses permission check", models.RolePlatformOwner, 1, "settings", fiber.StatusOK},
		{"Legacy Super bypasses permission check", models.LegacyRoleSuper, 1, "settings", fiber.StatusOK},
		{"Operations Admin with granted permission", models.RoleOperationsAdmin, 2, "pending", fiber.StatusOK},
		{"Operations Admin without permission", models.RoleOperationsAdmin, 2, "settings", fiber.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/test", func(c *fiber.Ctx) error {
				c.Locals("role", "admin")
				c.Locals("adminRole", tt.adminRole)
				c.Locals("userID", tt.adminID)
				return c.Next()
			}, mw.RequirePermission(tt.perm), func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Test request failed: %v", err)
			}
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("RequirePermission status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestRequireAuth_AdminJWTSessionBypass(t *testing.T) {
	jwtSecret := "test_super_secret_key_1234567890"

	// Mock database session store
	type sessionEntry struct {
		session *models.AdminSession
		user    *models.User
		err     error
	}

	sessions := map[int]sessionEntry{
		100: {
			session: &models.AdminSession{ID: 100, UserID: 1, Portal: "executive", Role: models.RolePlatformOwner},
			user:    &models.User{ID: 1, Role: "admin", AdminRole: models.RolePlatformOwner, IsActive: true},
			err:     nil,
		},
		200: {
			session: nil,
			user:    nil,
			err:     errors.New("session is revoked"),
		},
		300: {
			session: nil,
			user:    nil,
			err:     errors.New("session timed out due to inactivity"),
		},
		400: {
			session: &models.AdminSession{ID: 400, UserID: 4, Portal: "staff", Role: models.RolePlatformAdmin},
			user:    &models.User{ID: 4, Role: "admin", AdminRole: models.RolePlatformAdmin, IsActive: false}, // Disabled
			err:     nil,
		},
		500: {
			session: &models.AdminSession{ID: 500, UserID: 5, Portal: "staff", Role: models.RolePlatformOwner},
			user:    &models.User{ID: 5, Role: "admin", AdminRole: models.RoleModerator, IsActive: true}, // Demoted in DB to moderator!
			err:     nil,
		},
	}

	mw := New(jwtSecret).
		WithAdminSessionValidator(func(ctx context.Context, sessionID int, userID int) (*models.AdminSession, *models.User, error) {
			entry, ok := sessions[sessionID]
			if !ok || entry.err != nil {
				if entry.err != nil {
					return nil, nil, entry.err
				}
				return nil, nil, errors.New("session not found")
			}
			return entry.session, entry.user, nil
		}).
		WithSessionValidator(func(ctx context.Context, token string) (*models.AdminSession, *models.User, error) {
			if token == "11112222333344445555666677778888aaaabbbbccccddddeeeeffff00001111" {
				return &models.AdminSession{ID: 100, UserID: 1, Portal: "executive", Role: models.RolePlatformOwner},
					&models.User{ID: 1, Role: "admin", AdminRole: models.RolePlatformOwner, IsActive: true}, nil
			}
			return nil, nil, errors.New("revoked or invalid session")
		})

	t.Run("1. Admin JWT without server session fails with 401", func(t *testing.T) {
		app := fiber.New()
		app.Get("/admin/test", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		// Generate raw admin JWT without session (SessionID = 0)
		token, err := security.GenerateJWTFull(99, "admin", models.RolePlatformOwner, "", "", jwtSecret)
		if err != nil {
			t.Fatalf("GenerateJWTFull failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected status 401 Unauthorized for admin JWT without session, got %d", resp.StatusCode)
		}
	})

	t.Run("2. Admin JWT with revoked session fails with 401", func(t *testing.T) {
		app := fiber.New()
		app.Get("/admin/test", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		// JWT points to session 200 which is revoked
		token, err := security.GenerateAdminJWT(2, models.RolePlatformOwner, 200, "token_200", jwtSecret)
		if err != nil {
			t.Fatalf("GenerateAdminJWT failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected status 401 Unauthorized for revoked session, got %d", resp.StatusCode)
		}
	})

	t.Run("3. Admin JWT with expired/timed-out session fails with 401", func(t *testing.T) {
		app := fiber.New()
		app.Get("/admin/test", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		// JWT points to session 300 which timed out
		token, err := security.GenerateAdminJWT(3, models.RolePlatformOwner, 300, "token_300", jwtSecret)
		if err != nil {
			t.Fatalf("GenerateAdminJWT failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected status 401 Unauthorized for timed out session, got %d", resp.StatusCode)
		}
	})

	t.Run("4. Admin JWT with disabled account fails with 401", func(t *testing.T) {
		app := fiber.New()
		app.Get("/admin/test", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		// JWT points to session 400 where user.IsActive == false
		token, err := security.GenerateAdminJWT(4, models.RolePlatformAdmin, 400, "token_400", jwtSecret)
		if err != nil {
			t.Fatalf("GenerateAdminJWT failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected status 401 Unauthorized for disabled account, got %d", resp.StatusCode)
		}
	})

	t.Run("5. Demoted admin JWT reflects live DB role and is blocked from privileged route with 403", func(t *testing.T) {
		app := fiber.New()
		// Privileged route requiring platform_owner
		app.Get("/admin/owner-only", mw.RequireAuth, RequireRole("admin"), RequireSuperAdmin, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		// JWT claims platform_owner, but session 500 has live DB role moderator
		token, err := security.GenerateAdminJWT(5, models.RolePlatformOwner, 500, "token_500", jwtSecret)
		if err != nil {
			t.Fatalf("GenerateAdminJWT failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/admin/owner-only", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if resp.StatusCode != fiber.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden for demoted admin accessing owner route, got %d", resp.StatusCode)
		}
	})

	t.Run("6. Admin JWT with active valid session succeeds with 200", func(t *testing.T) {
		app := fiber.New()
		app.Get("/admin/test", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		// Valid session 100
		token, err := security.GenerateAdminJWT(1, models.RolePlatformOwner, 100, "token_100", jwtSecret)
		if err != nil {
			t.Fatalf("GenerateAdminJWT failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status 200 OK for valid admin session, got %d", resp.StatusCode)
		}
	})

	t.Run("7. 64-char raw session token validates correctly", func(t *testing.T) {
		app := fiber.New()
		app.Get("/admin/test", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		// Active 64-char token
		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		req.Header.Set("Authorization", "Bearer 11112222333344445555666677778888aaaabbbbccccddddeeeeffff00001111")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status 200 OK for active 64-char token, got %d", resp.StatusCode)
		}

		// Revoked 64-char token
		reqBad := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		reqBad.Header.Set("Authorization", "Bearer 99998888777766665555444433332222aaaabbbbccccddddeeeeffff00001111")
		respBad, err := app.Test(reqBad)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if respBad.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected status 401 Unauthorized for revoked 64-char token, got %d", respBad.StatusCode)
		}
	})

	t.Run("8. Marketplace users with regular JWT are unaffected", func(t *testing.T) {
		app := fiber.New()
		app.Get("/seller/dashboard", mw.RequireAuth, RequireRole("seller"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})
		app.Get("/admin/restricted", mw.RequireAuth, RequireRole("admin"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		sellerToken, err := security.GenerateJWTFull(10, "seller", "", "normal", "sale", jwtSecret)
		if err != nil {
			t.Fatalf("GenerateJWTFull failed: %v", err)
		}

		// Seller route succeeds
		reqSeller := httptest.NewRequest(http.MethodGet, "/seller/dashboard", nil)
		reqSeller.Header.Set("Authorization", "Bearer "+sellerToken)
		respSeller, err := app.Test(reqSeller)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if respSeller.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status 200 OK for seller on seller route, got %d", respSeller.StatusCode)
		}

		// Seller accessing admin route is forbidden
		reqAdmin := httptest.NewRequest(http.MethodGet, "/admin/restricted", nil)
		reqAdmin.Header.Set("Authorization", "Bearer "+sellerToken)
		respAdmin, err := app.Test(reqAdmin)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if respAdmin.StatusCode != fiber.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden for seller on admin route, got %d", respAdmin.StatusCode)
		}
	})
}

