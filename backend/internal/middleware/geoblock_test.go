package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestEthiopiaGeoAccess(t *testing.T) {
	app := fiber.New()
	app.Use(EthiopiaGeoAccess(GeoPolicyConfig{
		DevMode: false,
		EnabledCheck: func() bool {
			return true
		},
		ExemptPaths: []string{"/api/telebirr/webhook", "/api/chapa/webhook"},
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	app.Post("/api/telebirr/webhook", func(c *fiber.Ctx) error {
		return c.SendString("webhook ok")
	})
	app.Get("/api/products", func(c *fiber.Ctx) error {
		return c.SendString("products")
	})

	t.Run("Exempt path /health is allowed regardless of country header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.Header.Set("CF-IPCountry", "US")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("Exempt path /api/telebirr/webhook is allowed from foreign origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/telebirr/webhook", nil)
		req.Header.Set("CF-IPCountry", "CN")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("Ethiopia CF-IPCountry header is permitted", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
		req.Header.Set("CF-IPCountry", "ET")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("Foreign CF-IPCountry header is rejected with 403 Forbidden", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
		req.Header.Set("CF-IPCountry", "US")
		// RemoteAddr non-private
		req.RemoteAddr = "8.8.8.8:12345"
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", resp.StatusCode)
		}
	})
}

