package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
	"github.com/your-org/used-tools-marketplace/internal/middleware"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

func setupTestApp(t *testing.T) (*fiber.App, *AdminAuthHandler, *service.AdminAuthService, *repository.AdminPasskeyRepo, *repository.UserRepo, *config.Config, func()) {
	if os.Getenv("DB_PASSWORD") == "" && os.Getenv("DATABASE_URL") == "" {
		t.Skip("Skipping database-dependent test: DB_PASSWORD or DATABASE_URL not set in environment")
		return nil, nil, nil, nil, nil, nil, nil
	}

	cfg := config.NewTestConfig()
	db, err := database.NewPool(cfg)
	if err != nil {
		t.Skipf("Skipping database-dependent test: DB connection failed: %v", err)
		return nil, nil, nil, nil, nil, nil, nil
	}

	adminPasskeyRepo := repository.NewAdminPasskeyRepo(db)
	userRepo := repository.NewUserRepo(db)
	adminRepo := repository.NewAdminRepo(db)

	waManager, err := security.NewWebAuthnManager("localhost", []string{"http://localhost:3000", "http://localhost:8080"})
	if err != nil {
		db.Close()
		t.Fatalf("Failed to initialize WebAuthnManager: %v", err)
	}

	adminAuthSvc := service.NewAdminAuthService(adminPasskeyRepo, userRepo, adminRepo, waManager, cfg.JWTSecret)
	adminSvc := service.NewAdminService(adminRepo)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	adminAuthHandler := NewAdminAuthHandler(adminAuthSvc, authSvc, adminSvc, false)

	mw := middleware.New(cfg.JWTSecret).
		WithSessionValidator(adminAuthSvc.ValidateSession).
		WithAdminSessionValidator(adminAuthSvc.ValidateSessionByID)

	app := fiber.New(fiber.Config{
		ProxyHeader: fiber.HeaderXForwardedFor,
	})

	// Register both canonical and alias routes as done in main.go
	api := app.Group("/api")

	// Login begin / finish
	api.Post("/admin/passkeys/login/begin", adminAuthHandler.BeginLogin)
	api.Post("/admin/auth/passkey/login/begin", adminAuthHandler.BeginLogin)
	api.Post("/admin/passkeys/login/finish", adminAuthHandler.FinishLogin)
	api.Post("/admin/auth/passkey/login/verify", adminAuthHandler.FinishLogin)

	// Registration begin / finish (requires admin role & active session)
	api.Post("/admin/passkeys/register/begin", mw.RequireAuth, middleware.RequireRole("admin"), adminAuthHandler.BeginRegistration)
	api.Post("/admin/auth/passkey/register/begin", mw.RequireAuth, middleware.RequireRole("admin"), adminAuthHandler.BeginRegistration)
	api.Post("/admin/passkeys/register/finish", mw.RequireAuth, middleware.RequireRole("admin"), adminAuthHandler.FinishRegistration)
	api.Post("/admin/auth/passkey/register/verify", mw.RequireAuth, middleware.RequireRole("admin"), adminAuthHandler.FinishRegistration)

	// Passkey management
	api.Get("/admin/passkeys", mw.RequireAuth, middleware.RequireRole("admin"), adminAuthHandler.ListPasskeys)
	api.Get("/admin/auth/passkey/list", mw.RequireAuth, middleware.RequireRole("admin"), adminAuthHandler.ListPasskeys)
	api.Delete("/admin/passkeys/:id", mw.RequireAuth, middleware.RequireRole("admin"), adminAuthHandler.RevokePasskey)
	api.Delete("/admin/auth/passkey/:id", mw.RequireAuth, middleware.RequireRole("admin"), adminAuthHandler.RevokePasskey)

	// Reauth & session
	api.Post("/admin/reauth/begin", mw.RequireAuth, middleware.RequireRole("admin"), adminAuthHandler.BeginReauth)
	api.Post("/admin/reauth/finish", mw.RequireAuth, middleware.RequireRole("admin"), adminAuthHandler.FinishReauth)
	api.Post("/admin/auth/reauth", mw.RequireAuth, middleware.RequireRole("admin"), adminAuthHandler.ReauthPassword)
	api.Post("/admin/logout", adminAuthHandler.Logout)
	api.Post("/admin/auth/session/revoke", mw.RequireAuth, middleware.RequireRole("admin"), adminAuthHandler.RevokeSessionEndpoint)

	cleanup := func() {
		db.Close()
	}

	return app, adminAuthHandler, adminAuthSvc, adminPasskeyRepo, userRepo, cfg, cleanup
}

func TestWebAuthnRoutesAndSchemasAlignment(t *testing.T) {
	app, _, adminAuthSvc, adminPasskeyRepo, userRepo, cfg, cleanup := setupTestApp(t)
	if app == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// 1. Test Login Begin on both canonical and alias routes with both {email} and {identifier}
	t.Run("Login Begin Canonical and Alias with normalized payload", func(t *testing.T) {
		testEmail := fmt.Sprintf("test_admin_%d@huntiq.local", time.Now().UnixNano())
		routes := []string{"/api/admin/passkeys/login/begin", "/api/admin/auth/passkey/login/begin"}
		payloads := []map[string]string{
			{"email": testEmail, "portal": "staff"},
			{"identifier": testEmail, "portal": "staff"},
		}

		for _, route := range routes {
			for _, p := range payloads {
				body, _ := json.Marshal(p)
				req := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Forwarded-For", fmt.Sprintf("198.51.100.%d", (time.Now().UnixNano()/1000)%250+1))

				resp, err := app.Test(req)
				if err != nil {
					t.Fatalf("Request to %s failed: %v", route, err)
				}
				if resp.StatusCode != fiber.StatusOK {
					var errBody map[string]interface{}
					_ = json.NewDecoder(resp.Body).Decode(&errBody)
					t.Fatalf("Expected status 200 on %s with payload %v, got %d (body: %v)", route, p, resp.StatusCode, errBody)
				}

				var resEnvelope map[string]interface{}
				_ = json.NewDecoder(resp.Body).Decode(&resEnvelope)

				dataMap, ok := resEnvelope["data"].(map[string]interface{})
				if !ok {
					dataMap = resEnvelope
				}

				// Verify presence of options and both challenge and challenge_id
				if dataMap["options"] == nil {
					t.Errorf("Expected options in response from %s (got %v)", route, resEnvelope)
				}
				challenge, ok1 := dataMap["challenge"].(string)
				challengeID, ok2 := dataMap["challenge_id"].(string)
				if !ok1 || challenge == "" {
					t.Errorf("Expected non-empty 'challenge' field from %s", route)
				}
				if !ok2 || challengeID == "" {
					t.Errorf("Expected non-empty 'challenge_id' field from %s", route)
				}
				if challenge != challengeID {
					t.Errorf("challenge (%s) and challenge_id (%s) must match", challenge, challengeID)
				}
			}
		}
	})

	// 2. Test Login Finish with invalid assertions, expired challenges, and replayed challenges
	t.Run("Login Finish Assertion Validation and Challenge Lifecycles", func(t *testing.T) {
		routes := []string{"/api/admin/passkeys/login/finish", "/api/admin/auth/passkey/login/verify"}

		for _, route := range routes {
			// a. Malformed/missing assertion returns 400
			badPayload, _ := json.Marshal(map[string]interface{}{
				"challenge_id": "nonexistent_challenge_12345",
				"response":     map[string]interface{}{"rawId": "invalid"},
			})
			req := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(badPayload))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			if resp.StatusCode != fiber.StatusBadRequest && resp.StatusCode != fiber.StatusUnauthorized {
				t.Errorf("Expected 400 or 401 on invalid assertion on %s, got %d", route, resp.StatusCode)
			}

			// b. Expired challenge is strictly rejected
			expiredChallenge := "expired_test_challenge_" + service.GenerateRandomToken(16)
			_ = adminPasskeyRepo.SaveChallenge(ctx, &models.AuthChallenge{
				Challenge:    expiredChallenge,
				Portal:       "staff",
				CeremonyType: "authentication",
				SessionData:  "{}",
				ExpiresAt:    time.Now().Add(-10 * time.Minute), // Already expired
			})

			expPayload, _ := json.Marshal(map[string]interface{}{
				"challenge": expiredChallenge,
				"response":  map[string]interface{}{"id": "abc"},
			})
			reqExp := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(expPayload))
			reqExp.Header.Set("Content-Type", "application/json")
			respExp, err := app.Test(reqExp)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			if respExp.StatusCode != fiber.StatusBadRequest && respExp.StatusCode != fiber.StatusUnauthorized {
				t.Errorf("Expected 400 or 401 on expired challenge on %s, got %d", route, respExp.StatusCode)
			}

			// c. Replayed / consumed challenge cannot be used a second time
			replayedChallenge := "replay_test_challenge_" + service.GenerateRandomToken(16)
			_ = adminPasskeyRepo.SaveChallenge(ctx, &models.AuthChallenge{
				Challenge:    replayedChallenge,
				Portal:       "staff",
				CeremonyType: "authentication",
				SessionData:  "{}",
				ExpiresAt:    time.Now().Add(5 * time.Minute),
			})

			// Consume it once
			_, err = adminPasskeyRepo.GetAndConsumeChallenge(ctx, replayedChallenge)
			if err != nil {
				t.Fatalf("First consumption should succeed: %v", err)
			}

			// Second consumption attempt must fail
			_, err = adminPasskeyRepo.GetAndConsumeChallenge(ctx, replayedChallenge)
			if err == nil {
				t.Errorf("Challenge replay must fail, but got nil error")
			}
		}
	})

	// 3. Test Registration Begin & Finish authentication gates and payload alignment
	t.Run("Registration Begin & Finish Gates", func(t *testing.T) {
		regRoutes := []string{"/api/admin/passkeys/register/begin", "/api/admin/auth/passkey/register/begin"}
		for _, r := range regRoutes {
			// Without auth returns 401
			req := httptest.NewRequest(http.MethodPost, r, bytes.NewReader([]byte(`{"friendly_name":"Key"}`)))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)
			if resp.StatusCode != fiber.StatusUnauthorized {
				t.Errorf("Expected 401 Unauthorized without auth on %s, got %d", r, resp.StatusCode)
			}
		}

		verifyRoutes := []string{"/api/admin/passkeys/register/finish", "/api/admin/auth/passkey/register/verify"}
		for _, r := range verifyRoutes {
			// Without auth returns 401
			req := httptest.NewRequest(http.MethodPost, r, bytes.NewReader([]byte(`{"challenge_id":"abc"}`)))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)
			if resp.StatusCode != fiber.StatusUnauthorized {
				t.Errorf("Expected 401 Unauthorized without auth on %s, got %d", r, resp.StatusCode)
			}
		}
	})

	// 4. Test Passkey List and Delete on both canonical and alias routes
	t.Run("Passkey List and Delete Route Aliases", func(t *testing.T) {
		listRoutes := []string{"/api/admin/passkeys", "/api/admin/auth/passkey/list"}
		for _, r := range listRoutes {
			req := httptest.NewRequest(http.MethodGet, r, nil)
			resp, _ := app.Test(req)
			if resp.StatusCode != fiber.StatusUnauthorized {
				t.Errorf("Expected 401 Unauthorized without auth on %s, got %d", r, resp.StatusCode)
			}
		}

		deleteRoutes := []string{"/api/admin/passkeys/9999", "/api/admin/auth/passkey/9999"}
		for _, r := range deleteRoutes {
			req := httptest.NewRequest(http.MethodDelete, r, nil)
			resp, _ := app.Test(req)
			if resp.StatusCode != fiber.StatusUnauthorized {
				t.Errorf("Expected 401 Unauthorized without auth on %s, got %d", r, resp.StatusCode)
			}
		}

		// Non-admin marketplace users (buyer / seller) MUST be rejected with 403 Forbidden
		sellerToken, _ := security.GenerateJWTFull(55, "seller", "", "normal", "sale", cfg.JWTSecret)
		for _, r := range listRoutes {
			req := httptest.NewRequest(http.MethodGet, r, nil)
			req.Header.Set("Authorization", "Bearer "+sellerToken)
			resp, _ := app.Test(req)
			if resp.StatusCode != fiber.StatusForbidden {
				t.Errorf("Expected 403 Forbidden for non-admin on %s, got %d", r, resp.StatusCode)
			}
		}
	})

	// 5. Test Reauth and Session Revoke endpoints
	t.Run("Reauth and Session Revoke Gates", func(t *testing.T) {
		reqReauth := httptest.NewRequest(http.MethodPost, "/api/admin/auth/reauth", bytes.NewReader([]byte(`{"password":"test"}`)))
		reqReauth.Header.Set("Content-Type", "application/json")
		respReauth, _ := app.Test(reqReauth)
		if respReauth.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized without auth on /api/admin/auth/reauth, got %d", respReauth.StatusCode)
		}

		reqRevoke := httptest.NewRequest(http.MethodPost, "/api/admin/auth/session/revoke", bytes.NewReader([]byte(`{"session_id":"1"}`)))
		reqRevoke.Header.Set("Content-Type", "application/json")
		respRevoke, _ := app.Test(reqRevoke)
		if respRevoke.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized without auth on /api/admin/auth/session/revoke, got %d", respRevoke.StatusCode)
		}
	})

	// 6. Test Authenticated Passkey Management Ceremonies
	t.Run("Authenticated Passkey Management Ceremonies", func(t *testing.T) {
		admin, err := userRepo.GetAdminByEmail(ctx, "admin@huntiq.local")
		if err != nil || admin == nil {
			// Find any active admin
			admin, err = userRepo.GetByID(ctx, 36)
			if err != nil || admin == nil || admin.Role != "admin" {
				t.Skipf("No test admin found in database: %v", err)
				return
			}
		}

		rawToken, sess, err := adminAuthSvc.CreateAdminSession(ctx, admin.ID, "staff", "127.0.0.1", "test-agent")
		if err != nil {
			t.Fatalf("CreateAdminSession failed: %v", err)
		}
		defer func() {
			_ = adminPasskeyRepo.RevokeSession(ctx, sess.ID, "test_complete")
		}()

		jwtToken, err := security.GenerateAdminJWT(admin.ID, admin.AdminRole, sess.ID, rawToken, cfg.JWTSecret)
		if err != nil {
			t.Fatalf("GenerateAdminJWT failed: %v", err)
		}

		// a. Registration Begin on both canonical and alias routes with Auth
		regRoutes := []string{"/api/admin/passkeys/register/begin", "/api/admin/auth/passkey/register/begin"}
		for _, r := range regRoutes {
			req := httptest.NewRequest(http.MethodPost, r, bytes.NewReader([]byte(`{"friendly_name":"Test Authenticator"}`)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+jwtToken)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			if resp.StatusCode != fiber.StatusOK {
				t.Fatalf("Expected 200 on %s with auth, got %d", r, resp.StatusCode)
			}
			var resEnv map[string]interface{}
			_ = json.NewDecoder(resp.Body).Decode(&resEnv)
			dataMap, ok := resEnv["data"].(map[string]interface{})
			if !ok {
				dataMap = resEnv
			}
			if dataMap["options"] == nil || dataMap["challenge"] == nil || dataMap["challenge_id"] == nil {
				t.Errorf("Expected options, challenge, challenge_id in %s, got %v", r, dataMap)
			}
		}

		// b. List Passkeys on both canonical and alias routes with Auth
		listRoutes := []string{"/api/admin/passkeys", "/api/admin/auth/passkey/list"}
		for _, r := range listRoutes {
			req := httptest.NewRequest(http.MethodGet, r, nil)
			req.Header.Set("Authorization", "Bearer "+jwtToken)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			if resp.StatusCode != fiber.StatusOK {
				t.Fatalf("Expected 200 on %s with auth, got %d", r, resp.StatusCode)
			}
			var resEnv map[string]interface{}
			_ = json.NewDecoder(resp.Body).Decode(&resEnv)
			dataMap, ok := resEnv["data"].(map[string]interface{})
			if !ok {
				dataMap = resEnv
			}
			if dataMap["passkeys"] == nil || dataMap["authenticators"] == nil {
				t.Errorf("Expected both passkeys and authenticators in %s, got %v", r, dataMap)
			}
		}

		// c. Revoke non-existent passkey returns 400 Bad Request
		delRoutes := []string{"/api/admin/passkeys/999999", "/api/admin/auth/passkey/999999"}
		for _, r := range delRoutes {
			req := httptest.NewRequest(http.MethodDelete, r, nil)
			req.Header.Set("Authorization", "Bearer "+jwtToken)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			if resp.StatusCode != fiber.StatusBadRequest {
				t.Fatalf("Expected 400 on %s for non-existent passkey, got %d", r, resp.StatusCode)
			}
		}

		// d. Logout endpoint clears session cookie
		reqLogout := httptest.NewRequest(http.MethodPost, "/api/admin/logout", nil)
		respLogout, err := app.Test(reqLogout)
		if err != nil {
			t.Fatalf("Logout request failed: %v", err)
		}
		if respLogout.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected 200 on /api/admin/logout, got %d", respLogout.StatusCode)
		}
	})
}
