package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
	"github.com/your-org/used-tools-marketplace/internal/middleware"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

func setupSellerManagementApp(t *testing.T) (*fiber.App, *AdminHandler, *service.AdminAuthService, *pgxpool.Pool, string, *int, func()) {
	if os.Getenv("DB_PASSWORD") == "" && os.Getenv("DATABASE_URL") == "" {
		t.Skip("Skipping live database test: DB_PASSWORD or DATABASE_URL not set in environment")
		return nil, nil, nil, nil, "", nil, nil
	}

	cfg := config.NewTestConfig()
	db, err := database.NewPool(cfg)
	if err != nil {
		t.Skipf("Skipping live database test: DB connection failed: %v", err)
		return nil, nil, nil, nil, "", nil, nil
	}

	adminRepo := repository.NewAdminRepo(db)
	userRepo := repository.NewUserRepo(db)
	adminPasskeyRepo := repository.NewAdminPasskeyRepo(db)

	waManager, _ := security.NewWebAuthnManager("localhost", []string{"http://localhost:3000", "http://localhost:8080"})
	adminAuthSvc := service.NewAdminAuthService(adminPasskeyRepo, userRepo, adminRepo, waManager, cfg.JWTSecret)
	adminSvc := service.NewAdminService(adminRepo)

	var revokedUserID int
	var revokedMu sync.Mutex

	adminSvc.WithSessionRevoker(func(ctx context.Context, userID int, reason string) error {
		revokedMu.Lock()
		revokedUserID = userID
		revokedMu.Unlock()
		return adminAuthSvc.RevokeAllSessions(ctx, userID, reason)
	})

	mw := middleware.New(cfg.JWTSecret).
		WithSessionValidator(adminAuthSvc.ValidateSession).
		WithAdminSessionValidator(adminAuthSvc.ValidateSessionByID).
		WithIsActiveChecker(func(userID int) (bool, error) {
			u, err := userRepo.GetByID(context.Background(), userID)
			if err != nil {
				return false, err
			}
			return u.IsActive, nil
		}).
		WithPermissionChecker(func(adminID int, permission string) (bool, error) {
			u, err := userRepo.GetByID(context.Background(), adminID)
			if err != nil {
				return false, err
			}
			// Platform owner, platform admin, and operations admin have sellers permission
			if u.AdminRole == models.RolePlatformOwner || u.AdminRole == models.LegacyRoleSuper ||
				u.AdminRole == models.RolePlatformAdmin || u.AdminRole == models.RoleOperationsAdmin {
				return true, nil
			}
			if u.AdminRole == models.RoleModerator && permission != "sellers" {
				return true, nil
			}
			return false, nil
		})

	adminHandler := NewAdminHandler(adminSvc)

	app := fiber.New()
	admin := app.Group("/api/admin", mw.RequireAuth, middleware.RequireRole("admin"))
	admin.Get("/sellers", mw.RequirePermission("sellers"), adminHandler.ListSellers)
	admin.Put("/sellers/:id/active", mw.RequirePermission("sellers"), adminHandler.SetSellerActive)

	cleanup := func() {
		db.Close()
	}

	return app, adminHandler, adminAuthSvc, db, cfg.JWTSecret, &revokedUserID, cleanup
}

func TestSellerManagement_Endpoints(t *testing.T) {
	app, _, adminAuthSvc, db, jwtSecret, revokedUserID, cleanup := setupSellerManagementApp(t)
	if app == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()
	ts := time.Now().UnixNano()

	// 1. Create Test Platform Owner Admin
	ownerPhone := fmt.Sprintf("+25191%07d", ts%10000000)
	var ownerID int
	err := db.QueryRow(ctx, `
		INSERT INTO users (name, phone, role, admin_role, is_active, rating, referral_code)
		VALUES ('Owner Test', $1, 'admin', $2, true, 5.0, $3)
		RETURNING id`, ownerPhone, models.RolePlatformOwner, "REF"+ownerPhone[len(ownerPhone)-6:]).Scan(&ownerID)
	if err != nil {
		t.Fatalf("Failed to insert test admin: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM users WHERE id = $1", ownerID)

	rawOwner, sessOwner, err := adminAuthSvc.CreateAdminSession(ctx, ownerID, "executive", "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create owner admin session: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM admin_sessions WHERE user_id = $1", ownerID)

	ownerToken, _ := security.GenerateAdminJWT(ownerID, models.RolePlatformOwner, sessOwner.ID, rawOwner, jwtSecret)

	// 2. Create Test Operations Admin
	opsPhone := fmt.Sprintf("+25192%07d", (ts+1)%10000000)
	var opsID int
	err = db.QueryRow(ctx, `
		INSERT INTO users (name, phone, role, admin_role, is_active, rating, referral_code)
		VALUES ('Ops Test', $1, 'admin', $2, true, 5.0, $3)
		RETURNING id`, opsPhone, models.RoleOperationsAdmin, "REF"+opsPhone[len(opsPhone)-6:]).Scan(&opsID)
	if err != nil {
		t.Fatalf("Failed to insert ops admin: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM users WHERE id = $1", opsID)

	rawOps, sessOps, err := adminAuthSvc.CreateAdminSession(ctx, opsID, "operations", "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create ops admin session: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM admin_sessions WHERE user_id = $1", opsID)

	opsToken, _ := security.GenerateAdminJWT(opsID, models.RoleOperationsAdmin, sessOps.ID, rawOps, jwtSecret)

	// 3. Create Test Moderator Admin (Forbidden for sellers management)
	modPhone := fmt.Sprintf("+25193%07d", (ts+2)%10000000)
	var modID int
	err = db.QueryRow(ctx, `
		INSERT INTO users (name, phone, role, admin_role, is_active, rating, referral_code)
		VALUES ('Mod Test', $1, 'admin', $2, true, 5.0, $3)
		RETURNING id`, modPhone, models.RoleModerator, "REF"+modPhone[len(modPhone)-6:]).Scan(&modID)
	if err != nil {
		t.Fatalf("Failed to insert mod admin: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM users WHERE id = $1", modID)

	rawMod, sessMod, err := adminAuthSvc.CreateAdminSession(ctx, modID, "moderation", "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create mod admin session: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM admin_sessions WHERE user_id = $1", modID)

	modToken, _ := security.GenerateAdminJWT(modID, models.RoleModerator, sessMod.ID, rawMod, jwtSecret)

	// 4. Create Test Seller Account
	sellerPhone := fmt.Sprintf("+25194%07d", (ts+3)%10000000)
	var sellerID int
	err = db.QueryRow(ctx, `
		INSERT INTO users (name, phone, business_name, email, role, is_active, tier, rating, referral_code)
		VALUES ('Test Merchant Alpha', $1, 'Alpha Enterprise Tools', 'alpha@merchants.et', 'seller', true, 'gold', 4.8, $2)
		RETURNING id`, sellerPhone, "REF"+sellerPhone[len(sellerPhone)-6:]).Scan(&sellerID)
	if err != nil {
		t.Fatalf("Failed to insert test seller: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM users WHERE id = $1", sellerID)
	defer db.Exec(ctx, "DELETE FROM admin_logs WHERE target_user_id = $1", sellerID)

	// A. Authentication & Authorization Tests
	t.Run("ListSellers_Unauthenticated_Returns401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/sellers", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d", resp.StatusCode)
		}
	})

	t.Run("ListSellers_ForbiddenRole_Returns403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/sellers", nil)
		req.Header.Set("Authorization", "Bearer "+modToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected 403 Forbidden for Moderator on sellers endpoint, got %d", resp.StatusCode)
		}
	})

	t.Run("ListSellers_AuthorizedOperationsAdmin_Returns200WithPagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/sellers?limit=10&page=1", nil)
		req.Header.Set("Authorization", "Bearer "+opsToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
		}

		var body struct {
			Success bool `json:"success"`
			Data    struct {
				Data  []models.User `json:"data"`
				Total int           `json:"total"`
				Page  int           `json:"page"`
				Limit int           `json:"limit"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !body.Success {
			t.Errorf("Expected success: true")
		}
		if body.Data.Total < 1 {
			t.Errorf("Expected total >= 1, got %d", body.Data.Total)
		}
		if body.Data.Limit != 10 {
			t.Errorf("Expected limit 10, got %d", body.Data.Limit)
		}
		if body.Data.Page != 1 {
			t.Errorf("Expected page 1, got %d", body.Data.Page)
		}
	})

	// B. Search and Filter Tests
	t.Run("ListSellers_FilterSearchByName", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/sellers?search=Alpha+Enterprise", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}

		var body struct {
			Data struct {
				Data  []models.User `json:"data"`
				Total int           `json:"total"`
			} `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&body)

		found := false
		for _, s := range body.Data.Data {
			if s.ID == sellerID {
				found = true
				if s.BusinessName != "Alpha Enterprise Tools" {
					t.Errorf("Expected business name 'Alpha Enterprise Tools', got '%s'", s.BusinessName)
				}
				if s.Tier != "gold" {
					t.Errorf("Expected tier 'gold', got '%s'", s.Tier)
				}
			}
		}
		if !found {
			t.Errorf("Expected to find test seller in search results")
		}
	})

	t.Run("ListSellers_FilterByTier", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/sellers?tier=gold", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}

		var body struct {
			Data struct {
				Data []models.User `json:"data"`
			} `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&body)

		for _, s := range body.Data.Data {
			if s.Tier != "gold" {
				t.Errorf("Expected all returned sellers to have tier 'gold', got %s", s.Tier)
			}
		}
	})

	// C. Deactivation and Validation Tests
	t.Run("SetSellerActive_DeactivateWithoutReason_Returns400", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]interface{}{
			"active": false,
			"reason": "",
		})
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/admin/sellers/%d/active", sellerID), bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+opsToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400 Bad Request when deactivating without reason, got %d", resp.StatusCode)
		}
	})

	t.Run("SetSellerActive_DeactivateWithReason_SucceedsAndRevokesSessions", func(t *testing.T) {
		deactReason := "Investigation for fraudulent listing reports"
		payload, _ := json.Marshal(map[string]interface{}{
			"active": false,
			"reason": deactReason,
		})
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/admin/sellers/%d/active", sellerID), bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+opsToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
		}

		// Verify seller database record updated
		var isActive bool
		var reasonInDB string
		err = db.QueryRow(ctx, "SELECT is_active, COALESCE(deactivation_reason, '') FROM users WHERE id = $1", sellerID).Scan(&isActive, &reasonInDB)
		if err != nil {
			t.Fatalf("Failed to query user: %v", err)
		}
		if isActive {
			t.Errorf("Expected seller is_active to be false")
		}
		if reasonInDB != deactReason {
			t.Errorf("Expected deactivation reason '%s', got '%s'", deactReason, reasonInDB)
		}

		// Verify session revocation hook called
		if *revokedUserID != sellerID {
			t.Errorf("Expected revokedUserID to be %d, got %d", sellerID, *revokedUserID)
		}

		// Verify audit log record created
		var logAction, logNotes string
		err = db.QueryRow(ctx, "SELECT action, notes FROM admin_logs WHERE target_user_id = $1 ORDER BY id DESC LIMIT 1", sellerID).Scan(&logAction, &logNotes)
		if err != nil {
			t.Fatalf("Failed to query admin_logs: %v", err)
		}
		if logAction != "deactivate" {
			t.Errorf("Expected admin_log action 'deactivate', got '%s'", logAction)
		}
		if logNotes != deactReason {
			t.Errorf("Expected admin_log notes '%s', got '%s'", deactReason, logNotes)
		}
	})

	t.Run("SetSellerActive_ReactivateSeller_Succeeds", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]interface{}{
			"active": true,
			"reason": "Merchant identity verified",
		})
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/admin/sellers/%d/active", sellerID), bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
		}

		var isActive bool
		err = db.QueryRow(ctx, "SELECT is_active FROM users WHERE id = $1", sellerID).Scan(&isActive)
		if err != nil {
			t.Fatalf("Failed to query user: %v", err)
		}
		if !isActive {
			t.Errorf("Expected seller is_active to be true after reactivation")
		}
	})
}
