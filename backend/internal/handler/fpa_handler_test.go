package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/middleware"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

func TestFPAHandlerEndpointsAndAuthorization(t *testing.T) {
	jwtSecret := "test-fpa-jwt-secret-key-32-bytes-minimum!!"
	app := fiber.New()

	// Simulated repository without live DB connection for fast, isolated testing
	fpaRepo := repository.NewFPARepo(nil)
	fpaSvc := service.NewFPAService(fpaRepo)
	fpaHandler := NewFPAHandler(fpaSvc)

	mw := middleware.New(jwtSecret).
		WithAdminSessionValidator(func(ctx context.Context, sessionID, userID int) (*models.AdminSession, *models.User, error) {
			role := "admin"
			adminRole := models.RolePlatformOwner
			if userID == 201 {
				adminRole = models.RolePlatformAdmin
			} else if userID == 202 {
				adminRole = models.RoleModerator
			}
			return &models.AdminSession{ID: sessionID, UserID: userID, Portal: "executive", LastActiveAt: time.Now(), AbsoluteExpiresAt: time.Now().Add(time.Hour)},
				&models.User{ID: userID, Role: role, AdminRole: adminRole, IsActive: true}, nil
		})

	admin := app.Group("/admin", mw.RequireAuth, middleware.RequireRole("admin"))
	financeViewAuth := mw.RequireAnyPermission("profit", "finance", "finance:view")
	financeManageAuth := mw.RequireAnyPermission("finance:manage", "finance_manage")

	admin.Get("/finance/forecast", financeViewAuth, fpaHandler.GetForecast)
	admin.Get("/finance/scenarios", financeViewAuth, fpaHandler.GetScenarios)
	admin.Get("/finance/capacity-20m", financeViewAuth, fpaHandler.GetCapacityPlanning20M)
	admin.Get("/finance/export", financeViewAuth, fpaHandler.ExportForecast)
	admin.Put("/finance/assumptions", financeManageAuth, fpaHandler.UpdateAssumptions)

	// Generate JWT tokens for test roles
	ownerToken, _ := security.GenerateAdminJWT(101, models.RolePlatformOwner, 1, "test-sess-owner", jwtSecret)
	platAdminToken, _ := security.GenerateAdminJWT(201, models.RolePlatformAdmin, 2, "test-sess-plat", jwtSecret)
	moderatorToken, _ := security.GenerateAdminJWT(202, models.RoleModerator, 3, "test-sess-mod", jwtSecret)

	// 1. Unauthenticated Request Denied
	t.Run("Unauthenticated_Denied_401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/finance/forecast", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized, got %d", resp.StatusCode)
		}
	})

	// 2. Moderator Role Denied (Separation of Duties)
	t.Run("Moderator_Denied_403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/finance/forecast", nil)
		req.Header.Set("Authorization", "Bearer "+moderatorToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden for Moderator without finance permission, got %d", resp.StatusCode)
		}
	})

	// 3. Platform Admin Authorized to View (Domain baseline permission)
	t.Run("PlatformAdmin_Authorized_To_View_200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/finance/forecast?scenario=base&horizon=12", nil)
		req.Header.Set("Authorization", "Bearer "+platAdminToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK for Platform Admin view, got %d", resp.StatusCode)
		}
	})

	// 4. Platform Admin FORBIDDEN from Mutating Assumptions (Separation of Duties enforced)
	t.Run("PlatformAdmin_Forbidden_From_Mutating_Assumptions_403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/admin/finance/assumptions", strings.NewReader(`{}`))
		req.Header.Set("Authorization", "Bearer "+platAdminToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden for Platform Admin mutating assumptions without finance:manage, got %d", resp.StatusCode)
		}
	})

	// 5. Platform Owner Full Access to View and Mutate
	t.Run("PlatformOwner_Authorized_200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/finance/forecast?scenario=base&horizon=60", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK for Platform Owner, got %d", resp.StatusCode)
		}

		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		data, ok := body["data"].(map[string]interface{})
		if !ok || data["scenario_id"] != "base" {
			t.Errorf("expected scenario_id = base in response data")
		}
	})

	// 6. Capacity 20M Scenario Endpoint
	t.Run("Capacity_20M_Endpoint_200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/finance/capacity-20m", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}
	})

	// 7. CSV Export Endpoint
	t.Run("CSV_Export_200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/finance/export?scenario=base&format=csv", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}
		contentType := resp.Header.Get("Content-Type")
		if contentType != "text/csv; charset=utf-8" {
			t.Errorf("expected text/csv; charset=utf-8, got %s", contentType)
		}
	})
}

// Ensure unused time import doesn't fail compilation
var _ = time.Now
