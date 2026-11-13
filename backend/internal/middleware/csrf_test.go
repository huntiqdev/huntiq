package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestCSRFProtect(t *testing.T) {
	isAllowedOrigin := func(origin string) bool {
		return origin == "https://huntiq.et" || origin == "http://localhost:3000"
	}

	cfg := CSRFConfig{
		CookieName:      "huntiq_csrf",
		HeaderName:      "X-CSRF-Token",
		IsOriginAllowed: isAllowedOrigin,
		ExemptPaths: []string{
			"/api/chapa/webhook",
			"/api/telebirr/webhook",
		},
	}

	setupApp := func() *fiber.App {
		app := fiber.New()
		app.Use(CSRFProtect(cfg))
		app.Get("/api/items", func(c *fiber.Ctx) error {
			return c.SendString("items list")
		})
		app.Post("/api/items", func(c *fiber.Ctx) error {
			return c.SendString("item created")
		})
		app.Put("/api/items/1", func(c *fiber.Ctx) error {
			return c.SendString("item updated")
		})
		app.Delete("/api/items/1", func(c *fiber.Ctx) error {
			return c.SendString("item deleted")
		})
		app.Post("/api/auth/login", func(c *fiber.Ctx) error {
			return c.SendString("login ok")
		})
		app.Post("/api/chapa/webhook", func(c *fiber.Ctx) error {
			return c.SendString("webhook ok")
		})
		return app
	}

	t.Run("Safe method GET passes without CSRF", func(t *testing.T) {
		app := setupApp()
		req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Mutating POST without CSRF cookie or header returns 403", func(t *testing.T) {
		app := setupApp()
		req := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		// Add a session cookie to indicate a browser session
		req.AddCookie(&http.Cookie{Name: "huntiq_session", Value: "jwt-token"})
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusForbidden {
			t.Errorf("expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("Mutating POST with CSRF cookie but missing header returns 403", func(t *testing.T) {
		app := setupApp()
		req := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "huntiq_session", Value: "jwt-token"})
		req.AddCookie(&http.Cookie{Name: "huntiq_csrf", Value: "valid-csrf-token"})
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusForbidden {
			t.Errorf("expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("Mutating POST with mismatched CSRF header returns 403", func(t *testing.T) {
		app := setupApp()
		req := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-CSRF-Token", "wrong-token")
		req.AddCookie(&http.Cookie{Name: "huntiq_session", Value: "jwt-token"})
		req.AddCookie(&http.Cookie{Name: "huntiq_csrf", Value: "valid-csrf-token"})
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusForbidden {
			t.Errorf("expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("Mutating POST with matching CSRF cookie and header succeeds", func(t *testing.T) {
		app := setupApp()
		csrfToken := "super-secure-anti-csrf-token-12345678"
		req := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-CSRF-Token", csrfToken)
		req.AddCookie(&http.Cookie{Name: "huntiq_session", Value: "jwt-token"})
		req.AddCookie(&http.Cookie{Name: "huntiq_csrf", Value: csrfToken})
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Mutating POST with untrusted Origin returns 403", func(t *testing.T) {
		app := setupApp()
		csrfToken := "valid-csrf-token"
		req := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://attacker.com")
		req.Header.Set("X-CSRF-Token", csrfToken)
		req.AddCookie(&http.Cookie{Name: "huntiq_session", Value: "jwt-token"})
		req.AddCookie(&http.Cookie{Name: "huntiq_csrf", Value: csrfToken})
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusForbidden {
			t.Errorf("expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("Mutating POST with trusted Origin succeeds", func(t *testing.T) {
		app := setupApp()
		csrfToken := "valid-csrf-token"
		req := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://huntiq.et")
		req.Header.Set("X-CSRF-Token", csrfToken)
		req.AddCookie(&http.Cookie{Name: "huntiq_session", Value: "jwt-token"})
		req.AddCookie(&http.Cookie{Name: "huntiq_csrf", Value: csrfToken})
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Public auth POST /api/auth/login passes without CSRF", func(t *testing.T) {
		app := setupApp()
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"phone":"0911000000","password":"password"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Exempt webhook path passes without CSRF", func(t *testing.T) {
		app := setupApp()
		req := httptest.NewRequest(http.MethodPost, "/api/chapa/webhook", strings.NewReader(`{"event":"charge.success"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Pure M2M Bearer token request without cookies passes without CSRF", func(t *testing.T) {
		app := setupApp()
		req := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-m2m-jwt-token")
		// NO cookies attached
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}
