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
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

func setupTrustSafetyApp(t *testing.T) (*fiber.App, *AdminHandler, *service.AdminAuthService, *pgxpool.Pool, string, func()) {
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "cRbQBaeoQakbnOO2FNC7BNzQf8UkubRxkhMxpIqKYEtXN4Y9o9g7pJwKQcHCUnmi")
	}
	if os.Getenv("OTP_ENCRYPT_KEY") == "" {
		os.Setenv("OTP_ENCRYPT_KEY", "yEr0PYPyolvDVLDJ8JaNtx6yln3cI+54fkRSulMCVHA=")
	}
	if os.Getenv("ALLOW_ORIGIN") == "" {
		os.Setenv("ALLOW_ORIGIN", "http://localhost:3000")
	}
	if os.Getenv("FRONTEND_URL") == "" {
		os.Setenv("FRONTEND_URL", "http://localhost:3000")
	}
	if os.Getenv("APP_BASE_URL") == "" {
		os.Setenv("APP_BASE_URL", "http://localhost:8080")
	}
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_USER", "postgres")
		os.Setenv("DB_PASSWORD", "gedam")
		os.Setenv("DB_NAME", "hulubedej_db")
		os.Setenv("DB_SSLMODE", "disable")
	}

	cfg := config.Load()
	db, err := database.NewPool(cfg)
	if err != nil {
		t.Skipf("Skipping database-dependent test: DB connection failed: %v", err)
		return nil, nil, nil, nil, "", nil
	}

	adminRepo := repository.NewAdminRepo(db)
	userRepo := repository.NewUserRepo(db)
	adminPasskeyRepo := repository.NewAdminPasskeyRepo(db)

	waManager, _ := security.NewWebAuthnManager("localhost", []string{"http://localhost:3000", "http://localhost:8080"})
	adminAuthSvc := service.NewAdminAuthService(adminPasskeyRepo, userRepo, adminRepo, waManager, cfg.JWTSecret)
	adminSvc := service.NewAdminService(adminRepo)

	// Active cache simulating main.go
	var (
		isActiveMu    sync.Mutex
		isActiveCache = map[int]struct {
			active  bool
			expires time.Time
		}{}
	)

	adminSvc.WithSessionRevoker(func(ctx context.Context, userID int, reason string) error {
		isActiveMu.Lock()
		delete(isActiveCache, userID)
		isActiveMu.Unlock()
		return adminAuthSvc.RevokeAllSessions(ctx, userID, reason)
	})

	mw := middleware.New(cfg.JWTSecret).
		WithSessionValidator(adminAuthSvc.ValidateSession).
		WithAdminSessionValidator(adminAuthSvc.ValidateSessionByID).
		WithIsActiveChecker(func(userID int) (bool, error) {
			isActiveMu.Lock()
			if e, ok := isActiveCache[userID]; ok && time.Now().Before(e.expires) {
				isActiveMu.Unlock()
				return e.active, nil
			}
			isActiveMu.Unlock()

			u, err := userRepo.GetByID(context.Background(), userID)
			if err != nil {
				return false, err
			}

			isActiveMu.Lock()
			isActiveCache[userID] = struct {
				active  bool
				expires time.Time
			}{u.IsActive, time.Now().Add(10 * time.Second)}
			isActiveMu.Unlock()

			return u.IsActive, nil
		}).
		WithPermissionChecker(func(adminID int, permission string) (bool, error) {
			// Super/platform owners bypass, or return true for trust_safety
			return true, nil
		})

	adminHandler := NewAdminHandler(adminSvc)

	app := fiber.New()
	admin := app.Group("/api/admin", mw.RequireAuth, middleware.RequireRole("admin"))
	admin.Put("/trust/reports/:id/review", mw.RequirePermission("trust_safety"), adminHandler.ReviewAbuseReport)

	// Protected seller route for testing access revocation
	seller := app.Group("/api/seller", mw.RequireAuth, middleware.RequireRole("seller"))
	seller.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	cleanup := func() {
		db.Close()
	}

	return app, adminHandler, adminAuthSvc, db, cfg.JWTSecret, cleanup
}

func TestReviewAbuseReport_HandlerSuspensionFlow(t *testing.T) {
	app, _, adminAuthSvc, db, jwtSecret, cleanup := setupTrustSafetyApp(t)
	if app == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()
	ts := time.Now().UnixNano()
	adminPhone := fmt.Sprintf("+25191%07d", ts%10000000)
	targetPhone := fmt.Sprintf("+25192%07d", (ts+5)%10000000)

	var adminID, targetUserID int
	_ = db.QueryRow(ctx,
		`INSERT INTO users (name, phone, role, admin_role, is_active, is_banned, referral_code, created_at, updated_at)
		 VALUES ('Admin Inc', $1, 'admin', 'super', TRUE, FALSE, $2, NOW(), NOW()) RETURNING id`,
		adminPhone, "REF"+adminPhone[len(adminPhone)-6:]).Scan(&adminID)

	_ = db.QueryRow(ctx,
		`INSERT INTO users (name, phone, role, is_active, is_banned, referral_code, created_at, updated_at)
		 VALUES ('Bad Seller', $1, 'seller', TRUE, FALSE, $2, NOW(), NOW()) RETURNING id`,
		targetPhone, "REF"+targetPhone[len(targetPhone)-6:]).Scan(&targetUserID)

	defer func() {
		_, _ = db.Exec(ctx, `DELETE FROM abuse_reports WHERE reported_user_id=$1 OR reporter_id=$1`, targetUserID)
		_, _ = db.Exec(ctx, `DELETE FROM users WHERE id IN ($1, $2)`, adminID, targetUserID)
	}()

	var reportID int
	_ = db.QueryRow(ctx,
		`INSERT INTO abuse_reports (reporter_id, reported_user_id, reason, description, status, priority)
		 VALUES ($1, $2, 'fraud', 'Counterfeit tools sold', 'pending', 'critical') RETURNING id`,
		adminID, targetUserID).Scan(&reportID)

	// Generate JWTs
	rawToken, sess, err := adminAuthSvc.CreateAdminSession(ctx, adminID, "staff", "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("failed to create admin session: %v", err)
	}
	adminJWT, err := security.GenerateAdminJWT(adminID, "super", sess.ID, rawToken, jwtSecret)
	if err != nil {
		t.Fatalf("failed to generate admin JWT: %v", err)
	}

	sellerJWT, err := security.GenerateJWTFull(targetUserID, "seller", "", "normal", "sale", jwtSecret)
	if err != nil {
		t.Fatalf("failed to generate seller JWT: %v", err)
	}

	// 1. Seller ping before suspension should succeed
	reqPing := httptest.NewRequest(http.MethodGet, "/api/seller/ping", nil)
	reqPing.Header.Set("Authorization", "Bearer "+sellerJWT)
	respPing, err := app.Test(reqPing)
	if err != nil || respPing.StatusCode != http.StatusOK {
		t.Fatalf("Expected seller ping before suspension to return 200, got %v (%d)", err, respPing.StatusCode)
	}

	// 2. Admin reviews abuse report sending { "status": "resolved", "action": "suspended", "notes": "Confirmed counterfeit" }
	// (frontend client sends 'action')
	reviewBody, _ := json.Marshal(map[string]interface{}{
		"status": "resolved",
		"action": "suspended",
		"notes":  "Confirmed counterfeit seller",
	})
	reqReview := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/admin/trust/reports/%d/review", reportID), bytes.NewReader(reviewBody))
	reqReview.Header.Set("Authorization", "Bearer "+adminJWT)
	reqReview.Header.Set("Content-Type", "application/json")
	respReview, err := app.Test(reqReview)
	if err != nil || respReview.StatusCode != http.StatusOK {
		t.Fatalf("Expected ReviewAbuseReport to return 200, got %v (%d)", err, respReview.StatusCode)
	}

	// 3. Verify user in DB is inactive and not banned
	var isActive, isBanned bool
	_ = db.QueryRow(ctx, "SELECT is_active, is_banned FROM users WHERE id=$1", targetUserID).Scan(&isActive, &isBanned)
	if isActive || isBanned {
		t.Errorf("Expected is_active=false, is_banned=false, got is_active=%v, is_banned=%v", isActive, isBanned)
	}

	// 4. Seller request with previously issued token MUST NOW BE DENIED WITH 401
	reqBlocked := httptest.NewRequest(http.MethodGet, "/api/seller/ping", nil)
	reqBlocked.Header.Set("Authorization", "Bearer "+sellerJWT)
	respBlocked, err := app.Test(reqBlocked)
	if err != nil {
		t.Fatalf("Test request failed: %v", err)
	}
	if respBlocked.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected previously issued seller token to be rejected with 401 Unauthorized after suspension, got %d", respBlocked.StatusCode)
	}

	// 5. Verify action_taken key variant also works
	targetUser2Phone := fmt.Sprintf("+25193%07d", (ts+15)%10000000)
	var targetUser2ID int
	_ = db.QueryRow(ctx,
		`INSERT INTO users (name, phone, role, is_active, is_banned, referral_code, created_at, updated_at)
		 VALUES ('Bad Seller 2', $1, 'seller', TRUE, FALSE, $2, NOW(), NOW()) RETURNING id`,
		targetUser2Phone, "REF"+targetUser2Phone[len(targetUser2Phone)-6:]).Scan(&targetUser2ID)
	defer func() {
		_, _ = db.Exec(ctx, `DELETE FROM abuse_reports WHERE reported_user_id=$1 OR reporter_id=$1`, targetUser2ID)
		_, _ = db.Exec(ctx, `DELETE FROM users WHERE id=$1`, targetUser2ID)
	}()

	var report2ID int
	_ = db.QueryRow(ctx,
		`INSERT INTO abuse_reports (reporter_id, reported_user_id, reason, description, status, priority)
		 VALUES ($1, $2, 'spam', 'Spam report', 'pending', 'normal') RETURNING id`,
		adminID, targetUser2ID).Scan(&report2ID)

	review2Body, _ := json.Marshal(map[string]interface{}{
		"status":       "resolved",
		"action_taken": "suspended",
		"notes":        "Suspension via action_taken key",
	})
	reqReview2 := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/admin/trust/reports/%d/review", report2ID), bytes.NewReader(review2Body))
	reqReview2.Header.Set("Authorization", "Bearer "+adminJWT)
	reqReview2.Header.Set("Content-Type", "application/json")
	respReview2, err := app.Test(reqReview2)
	if err != nil || respReview2.StatusCode != http.StatusOK {
		t.Fatalf("Expected ReviewAbuseReport with action_taken key to return 200, got %v (%d)", err, respReview2.StatusCode)
	}

	var isUser2Active bool
	_ = db.QueryRow(ctx, "SELECT is_active FROM users WHERE id=$1", targetUser2ID).Scan(&isUser2Active)
	if isUser2Active {
		t.Errorf("Expected user 2 is_active=false, got true")
	}
}
