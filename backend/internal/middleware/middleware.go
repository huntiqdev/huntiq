// Package middleware provides HTTP middleware for authentication, authorization, and request handling.
package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/pkg/constants"
)

// ============================================================================
// MIDDLEWARE CONTAINER - Dependency Injection
// ============================================================================

// Middleware holds all middleware dependencies injected at startup
type Middleware struct {
	jwtSecret          string
	isActiveChecker    func(userID int) (bool, error)
	permissionChecker  func(adminID int, permission string) (bool, error)
	sellerTierChecker  func(sellerID int) (string, error)
	maintenanceChecker func() (bool, error)
	// Maintenance cache
	mCacheMu      sync.RWMutex
	mCacheEnabled bool
	mCacheExpiry  time.Time
}

// New creates a new middleware container with injected dependencies
func New(jwtSecret string) *Middleware {
	return &Middleware{
		jwtSecret: jwtSecret,
	}
}

// WithIsActiveChecker injects the user active status checker
func (m *Middleware) WithIsActiveChecker(checker func(userID int) (bool, error)) *Middleware {
	m.isActiveChecker = checker
	return m
}

// WithPermissionChecker injects the admin permission checker
func (m *Middleware) WithPermissionChecker(checker func(adminID int, permission string) (bool, error)) *Middleware {
	m.permissionChecker = checker
	return m
}

// WithSellerTierChecker injects the seller tier checker
func (m *Middleware) WithSellerTierChecker(checker func(sellerID int) (string, error)) *Middleware {
	m.sellerTierChecker = checker
	return m
}

// WithMaintenanceChecker injects the maintenance mode checker
func (m *Middleware) WithMaintenanceChecker(checker func() (bool, error)) *Middleware {
	m.maintenanceChecker = checker
	return m
}

// ============================================================================
// AUTHENTICATION MIDDLEWARE
// ============================================================================

// RequireAuth validates the Bearer token and stores claims in locals.
func (m *Middleware) RequireAuth(c *fiber.Ctx) error {
	header := c.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return fiber.NewError(fiber.StatusUnauthorized, "missing or invalid authorization header")
	}
	claims, err := security.ParseJWT(strings.TrimPrefix(header, "Bearer "), m.jwtSecret)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}
	if m.isActiveChecker != nil {
		active, err := m.isActiveChecker(claims.UserID)
		if err != nil || !active {
			return fiber.NewError(fiber.StatusUnauthorized, "account is inactive or not found")
		}
	}
	c.Locals("userID", claims.UserID)
	c.Locals("role", claims.Role)
	c.Locals("adminRole", claims.AdminRole)
	c.Locals("buyerTier", claims.BuyerTier)
	return c.Next()
}

// OptionalAuth extracts JWT claims if present, but doesn't reject unauthenticated requests.
func (m *Middleware) OptionalAuth(c *fiber.Ctx) error {
	header := c.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		if claims, err := security.ParseJWT(strings.TrimPrefix(header, "Bearer "), m.jwtSecret); err == nil {
			c.Locals("userID", claims.UserID)
			c.Locals("role", claims.Role)
			c.Locals("adminRole", claims.AdminRole)
			c.Locals("buyerTier", claims.BuyerTier)
		}
	}
	return c.Next()
}

// ============================================================================
// ROLE-BASED AUTHORIZATION MIDDLEWARE
// ============================================================================

// RequireRole checks that the authenticated user has the expected role.
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("role") != role {
			return fiber.NewError(fiber.StatusForbidden, "insufficient permissions")
		}
		return c.Next()
	}
}

// RequireSuperAdmin ensures the admin has super role.
func RequireSuperAdmin(c *fiber.Ctx) error {
	if c.Locals("adminRole") != "super" {
		return fiber.NewError(fiber.StatusForbidden, "super admin access required")
	}
	return c.Next()
}

// adminRoleRank maps admin_role to a numeric rank for hierarchy checks.
var adminRoleRank = map[string]int{"super": 2, "normal": 1, "": 0}

// RequireAdminLevel enforces a minimum admin_role level.
func RequireAdminLevel(minRole string) fiber.Handler {
	required := adminRoleRank[minRole]
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "admin" {
			return fiber.NewError(fiber.StatusForbidden, "admin access required")
		}
		adminRole, _ := c.Locals("adminRole").(string)
		if adminRoleRank[adminRole] < required {
			return fiber.NewError(fiber.StatusForbidden, "insufficient admin privileges: "+minRole+" admin required")
		}
		return c.Next()
	}
}

// ============================================================================
// PERMISSION-BASED AUTHORIZATION
// ============================================================================

// RequirePermission checks that a normal admin has the given page permission.
// Super admins always pass.
func (m *Middleware) RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		adminRole, _ := c.Locals("adminRole").(string)
		if adminRole == "super" {
			return c.Next()
		}
		adminID, _ := c.Locals("userID").(int)
		if m.permissionChecker == nil {
			return fiber.NewError(fiber.StatusInternalServerError, "permission checker not configured")
		}
		ok, err := m.permissionChecker(adminID, permission)
		if err != nil || !ok {
			return fiber.NewError(fiber.StatusForbidden, "you do not have permission to access: "+permission)
		}
		return c.Next()
	}
}

// RequirePermissionOrSuper allows super admins unconditionally, and normal admins with permission.
func (m *Middleware) RequirePermissionOrSuper(permission string) fiber.Handler {
	return m.RequirePermission(permission)
}

// ============================================================================
// TIER-BASED AUTHORIZATION
// ============================================================================

// RequireSellerTier enforces a minimum seller tier.
func (m *Middleware) RequireSellerTier(minTier string) fiber.Handler {
	required := constants.TierRank[minTier]
	return func(c *fiber.Ctx) error {
		sellerID, _ := c.Locals("userID").(int)
		if sellerID == 0 {
			return fiber.NewError(fiber.StatusUnauthorized, "authentication required")
		}
		tier := "normal"
		if m.sellerTierChecker != nil {
			if t, err := m.sellerTierChecker(sellerID); err == nil && t != "" {
				tier = t
			}
		}
		if constants.TierRank[tier] < required {
			return fiber.NewError(fiber.StatusForbidden,
				"this feature requires "+minTier+" tier or higher — upgrade your account to access it")
		}
		return c.Next()
	}
}

// ============================================================================
// MAINTENANCE MODE
// ============================================================================

// RequireNoMaintenance blocks access if the platform is in maintenance mode.
// Admin routes and auth endpoints are always allowed.
func (m *Middleware) RequireNoMaintenance(c *fiber.Ctx) error {
	path := c.Path()

	// Always allow auth and public settings
	if strings.HasPrefix(path, "/api/auth") || strings.HasPrefix(path, "/api/public") {
		return c.Next()
	}

	// Skip check for logged-in admins
	if c.Locals("role") == "admin" {
		return c.Next()
	}

	m.mCacheMu.RLock()
	if time.Now().Before(m.mCacheExpiry) {
		enabled := m.mCacheEnabled
		m.mCacheMu.RUnlock()
		if enabled {
			return fiber.NewError(fiber.StatusServiceUnavailable, "Platform is currently under maintenance. Please try again later.")
		}
		return c.Next()
	}
	m.mCacheMu.RUnlock()

	// Cache expired, fetch fresh
	if m.maintenanceChecker != nil {
		enabled, err := m.maintenanceChecker()
		if err == nil {
			m.mCacheMu.Lock()
			// Double-check: another goroutine may have refreshed the cache while we waited for the lock.
			if time.Now().After(m.mCacheExpiry) {
				m.mCacheEnabled = enabled
				m.mCacheExpiry = time.Now().Add(10 * time.Second)
			}
			m.mCacheMu.Unlock()

			if enabled {
				return fiber.NewError(fiber.StatusServiceUnavailable, "Platform is currently under maintenance. Please try again later.")
			}
		}
	}

	return c.Next()
}

// ============================================================================
// REQUEST LOGGING
// ============================================================================

// RequestLogger logs method, path, status, latency, and client IP for every request.
func RequestLogger(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	log.Printf("[%s] %s %s %s %d %s",
		c.Get("X-Request-ID"),
		c.IP(),
		c.Method(),
		c.Path(),
		c.Response().StatusCode(),
		time.Since(start),
	)
	return err
}

// CorrelationID attaches a unique X-Request-ID to every request for tracing.
// Uses the incoming header if present, otherwise generates a new one.
func CorrelationID(c *fiber.Ctx) error {
	reqID := c.Get("X-Request-ID")
	if reqID == "" {
		b := make([]byte, 8)
		_, _ = rand.Read(b)
		reqID = hex.EncodeToString(b)
	}
	c.Locals("requestID", reqID)
	c.Set("X-Request-ID", reqID)
	return c.Next()
}
