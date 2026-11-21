package middleware

import (
	"crypto/subtle"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CSRFConfig struct {
	CookieName      string
	HeaderName      string
	IsOriginAllowed func(origin string) bool
	ExemptPaths     []string
}

// CSRFProtect creates an enterprise-grade Anti-CSRF middleware.
// It verifies that all mutating requests (POST, PUT, PATCH, DELETE) using cookie credentials
// provide a cryptographically matching Anti-CSRF token in the request header and a trusted Origin.
func CSRFProtect(cfg CSRFConfig) fiber.Handler {
	if cfg.CookieName == "" {
		cfg.CookieName = "huntiq_csrf"
	}
	if cfg.HeaderName == "" {
		cfg.HeaderName = "X-CSRF-Token"
	}

	return func(c *fiber.Ctx) error {
		method := c.Method()

		// 1. Safe HTTP methods never mutate state (RFC 7231)
		if method == fiber.MethodGet || method == fiber.MethodHead || method == fiber.MethodOptions || method == "TRACE" {
			return c.Next()
		}

		path := c.Path()

		// 2. Check path exemptions (e.g. payment provider webhooks that authenticate via HMAC signatures)
		for _, exempt := range cfg.ExemptPaths {
			if strings.HasPrefix(path, exempt) {
				return c.Next()
			}
		}

		// 3. Origin & Referer Validation
		origin := c.Get("Origin")
		if origin == "" {
			// Fallback to Referer host if Origin is omitted by older browser
			if ref := c.Get("Referer"); ref != "" {
				if parsed, err := url.Parse(ref); err == nil && parsed.Scheme != "" && parsed.Host != "" {
					origin = parsed.Scheme + "://" + parsed.Host
				}
			}
		}

		if origin != "" && cfg.IsOriginAllowed != nil {
			if !cfg.IsOriginAllowed(origin) {
				return fiber.NewError(fiber.StatusForbidden, "CSRF_ORIGIN_UNTRUSTED")
			}
		}

		// 4. Public auth endpoints (login, register, forgot/reset password, OTP)
		isPublicAuth := path == "/api/auth/login" ||
			path == "/api/auth/register" ||
			path == "/api/seller/login" ||
			path == "/api/seller/register" ||
			path == "/api/buyer/login" ||
			path == "/api/buyer/register" ||
			path == "/api/admin/login" ||
			path == "/api/admin/login/verify" ||
			path == "/api/auth/otp/send" ||
			path == "/api/auth/otp/verify" ||
			path == "/api/auth/otp/login" ||
			path == "/api/auth/forgot-password" ||
			path == "/api/auth/reset-password" ||
			path == "/api/seller/forgot-password" ||
			path == "/api/seller/reset-password" ||
			path == "/api/buyer/forgot-password" ||
			path == "/api/buyer/reset-password" ||
			path == "/api/admin/forgot-password" ||
			path == "/api/admin/reset-password" ||
			path == "/api/admin/passkeys/login/begin" ||
			path == "/api/admin/passkeys/login/finish" ||
			path == "/api/admin/auth/passkey/login/begin" ||
			path == "/api/admin/auth/passkey/login/verify"

		if isPublicAuth {
			return c.Next()
		}

		// 5. Check authentication context:
		// If the request presents an explicit Bearer Authorization header,
		// it is inherently immune to ambient browser cross-site request forgery
		// because browsers cannot attach custom Authorization headers in cross-origin requests.
		hasBearerAuth := strings.HasPrefix(c.Get("Authorization"), "Bearer ")
		if hasBearerAuth {
			// If client also supplied both CSRF cookie and header, ensure they cryptographically match
			csrfCookie := c.Cookies(cfg.CookieName)
			csrfHeader := c.Get(cfg.HeaderName)
			if csrfCookie != "" && csrfHeader != "" {
				if subtle.ConstantTimeCompare([]byte(csrfCookie), []byte(csrfHeader)) != 1 {
					return fiber.NewError(fiber.StatusForbidden, "CSRF_TOKEN_INVALID")
				}
			}
			return c.Next()
		}

		// 6. Anti-CSRF Token Verification for ambient cookie-authenticated mutating requests
		csrfCookie := c.Cookies(cfg.CookieName)
		csrfHeader := c.Get(cfg.HeaderName)

		if csrfCookie == "" || csrfHeader == "" {
			return fiber.NewError(fiber.StatusForbidden, "CSRF_TOKEN_MISSING")
		}

		// Constant-time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(csrfCookie), []byte(csrfHeader)) != 1 {
			return fiber.NewError(fiber.StatusForbidden, "CSRF_TOKEN_INVALID")
		}

		return c.Next()
	}
}
