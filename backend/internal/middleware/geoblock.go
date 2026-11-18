package middleware

import (
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GeoPolicyConfig holds settings for Ethiopia-first network access controls.
type GeoPolicyConfig struct {
	DevMode          bool
	EnabledCheck     func() bool
	ExemptPaths      []string
}

// isPrivateIP checks if an IP belongs to private or loopback subnets.
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return true // Fail-open on unparseable local IP in dev/test
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified()
}

// EthiopiaGeoAccess creates a middleware that restricts foreign public requests when enabled.
func EthiopiaGeoAccess(cfg GeoPolicyConfig) fiber.Handler {
	exemptMap := make(map[string]bool)
	for _, p := range cfg.ExemptPaths {
		exemptMap[p] = true
	}

	return func(c *fiber.Ctx) error {
		// Dev mode bypass
		if cfg.DevMode {
			return c.Next()
		}

		// Dynamic toggle check
		if cfg.EnabledCheck != nil && !cfg.EnabledCheck() {
			return c.Next()
		}

		path := c.Path()

		// Static exempt paths (health, webhooks, root)
		if path == "/" || path == "/health" || exemptMap[path] || strings.HasPrefix(path, "/api/telebirr/webhook") || strings.HasPrefix(path, "/api/chapa/webhook") {
			return c.Next()
		}

		// Admin routes are protected by strong MFA / Passkeys, not by coarse IP blocking
		if strings.HasPrefix(path, "/api/admin") {
			return c.Next()
		}

		// Inspect GeoIP headers populated by edge reverse proxy / WAF
		country := c.Get("CF-IPCountry")
		if country == "" {
			country = c.Get("X-Country-Code")
		}
		if country == "" {
			country = c.Get("X-Geo-Country")
		}

		// If edge country header is explicitly present and not Ethiopia (ET), restrict ordinary public access
		if country != "" && !strings.EqualFold(country, "ET") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "access_restricted: HUNTIQ is optimized for access within Ethiopia. If you are an authorized partner or legitimate user traveling abroad, contact support@huntiq.et for review.",
				"country": country,
			})
		}

		// Check local/private IP ranges (permits internal service calls and dev tunnels)
		clientIP := c.IP()
		if isPrivateIP(clientIP) {
			return c.Next()
		}

		return c.Next()
	}
}
