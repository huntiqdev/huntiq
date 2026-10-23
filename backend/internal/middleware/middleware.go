// Package middleware provides HTTP middleware for authentication, authorization, and request handling.
package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/pkg/constants"
)

// ============================================================================
// MIDDLEWARE CONTAINER - Dependency Injection
// ============================================================================

// Middleware holds all middleware dependencies injected at startup
type Middleware struct {
	jwtSecret             string
	isActiveChecker       func(userID int) (bool, error)
	permissionChecker     func(adminID int, permission string) (bool, error)
	sellerTierChecker     func(sellerID int) (string, error)
	maintenanceChecker    func() (bool, error)
	sessionValidator      func(ctx context.Context, token string) (*models.AdminSession, *models.User, error)
	adminSessionValidator func(ctx context.Context, sessionID int, userID int) (*models.AdminSession, *models.User, error)
	reauthValidator       func(ctx context.Context, token, scope string, userID int) (bool, error)
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

// WithSessionValidator injects the server-side session validator
func (m *Middleware) WithSessionValidator(v func(ctx context.Context, token string) (*models.AdminSession, *models.User, error)) *Middleware {
	m.sessionValidator = v
	return m
}

// WithAdminSessionValidator injects the validator for session-bound JWTs
func (m *Middleware) WithAdminSessionValidator(v func(ctx context.Context, sessionID int, userID int) (*models.AdminSession, *models.User, error)) *Middleware {
	m.adminSessionValidator = v
	return m
}

// WithReauthValidator injects the sensitive action reauth validator
func (m *Middleware) WithReauthValidator(v func(ctx context.Context, token, scope string, userID int) (bool, error)) *Middleware {
	m.reauthValidator = v
	return m
}

// ============================================================================
// AUTHENTICATION MIDDLEWARE
// ============================================================================

// RequireAuth validates the Bearer token or session cookie and stores claims in locals.
func (m *Middleware) RequireAuth(c *fiber.Ctx) error {
	token := ""
	header := c.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		token = strings.TrimPrefix(header, "Bearer ")
	} else if cookie := c.Cookies("hm_admin_session"); cookie != "" {
		token = cookie
	}

	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing or invalid authorization header or session cookie")
	}

	// 1. Check if token is a server-side session token (64 hex characters)
	if len(token) == 64 {
		if m.sessionValidator == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "session validator not configured")
		}
		session, user, err := m.sessionValidator(c.Context(), token)
		if err != nil || user == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid, expired, or revoked session")
		}
		c.Locals("userID", user.ID)
		c.Locals("role", user.Role)
		c.Locals("adminRole", user.AdminRole)
		c.Locals("portal", session.Portal)
		c.Locals("sessionID", session.ID)
		c.Locals("buyerTier", user.BuyerTier)
		c.Locals("capabilities", user.Capabilities)
		c.Locals("user", user)
		return c.Next()
	}

	// 2. Fallback to JWT parsing (for marketplace users and session-bound admin tokens)
	claims, err := security.ParseJWT(token, m.jwtSecret)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token or expired session")
	}
	if m.isActiveChecker != nil {
		active, err := m.isActiveChecker(claims.UserID)
		if err != nil || !active {
			return fiber.NewError(fiber.StatusUnauthorized, "account is inactive or not found")
		}
	}

	// 3. For administrative requests, verify active server-side session in admin_sessions
	if claims.Role == "admin" {
		var session *models.AdminSession
		var user *models.User
		var sErr error

		if claims.SessionID > 0 && m.adminSessionValidator != nil {
			session, user, sErr = m.adminSessionValidator(c.Context(), claims.SessionID, claims.UserID)
		} else if cookie := c.Cookies("hm_admin_session"); cookie != "" && m.sessionValidator != nil && len(cookie) == 64 {
			session, user, sErr = m.sessionValidator(c.Context(), cookie)
			if sErr == nil && (user == nil || user.ID != claims.UserID) {
				sErr = errors.New("session user mismatch")
				user = nil
			}
		} else {
			sErr = errors.New("no active session provided")
		}

		if sErr != nil || session == nil || user == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "administrative session required: session revoked, expired, or invalid")
		}

		if !user.IsActive || user.Role != "admin" {
			return fiber.NewError(fiber.StatusUnauthorized, "account is inactive, disabled, or demoted")
		}

		// Authorization claims reflect the current live database state, NOT stale JWT claims
		c.Locals("userID", user.ID)
		c.Locals("role", user.Role)
		c.Locals("adminRole", user.AdminRole) // LIVE database role!
		c.Locals("portal", session.Portal)
		c.Locals("sessionID", session.ID)
		c.Locals("buyerTier", user.BuyerTier)
		c.Locals("capabilities", user.Capabilities)
		c.Locals("user", user)
		c.Locals("claims", claims)
		return c.Next()
	}

	c.Locals("userID", claims.UserID)
	c.Locals("role", claims.Role)
	c.Locals("adminRole", claims.AdminRole)
	c.Locals("buyerTier", claims.BuyerTier)
	c.Locals("capabilities", claims.Capabilities)
	c.Locals("claims", claims)
	return c.Next()
}

// OptionalAuth extracts JWT claims or session if present, but doesn't reject unauthenticated requests.
func (m *Middleware) OptionalAuth(c *fiber.Ctx) error {
	token := ""
	header := c.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		token = strings.TrimPrefix(header, "Bearer ")
	} else if cookie := c.Cookies("hm_admin_session"); cookie != "" {
		token = cookie
	}

	if token != "" {
		if m.sessionValidator != nil && len(token) == 64 {
			if session, user, err := m.sessionValidator(c.Context(), token); err == nil && user != nil {
				c.Locals("userID", user.ID)
				c.Locals("role", user.Role)
				c.Locals("adminRole", user.AdminRole)
				c.Locals("portal", session.Portal)
				c.Locals("sessionID", session.ID)
				c.Locals("buyerTier", user.BuyerTier)
				c.Locals("capabilities", user.Capabilities)
				c.Locals("user", user)
				return c.Next()
			}
		}
		if claims, err := security.ParseJWT(token, m.jwtSecret); err == nil {
			if claims.Role == "admin" {
				var session *models.AdminSession
				var user *models.User
				if claims.SessionID > 0 && m.adminSessionValidator != nil {
					session, user, _ = m.adminSessionValidator(c.Context(), claims.SessionID, claims.UserID)
				}
				if session != nil && user != nil && user.IsActive && user.Role == "admin" {
					c.Locals("userID", user.ID)
					c.Locals("role", user.Role)
					c.Locals("adminRole", user.AdminRole)
					c.Locals("portal", session.Portal)
					c.Locals("sessionID", session.ID)
					c.Locals("buyerTier", user.BuyerTier)
					c.Locals("capabilities", user.Capabilities)
					c.Locals("user", user)
					c.Locals("claims", claims)
				}
			} else {
				c.Locals("userID", claims.UserID)
				c.Locals("role", claims.Role)
				c.Locals("adminRole", claims.AdminRole)
				c.Locals("buyerTier", claims.BuyerTier)
				c.Locals("capabilities", claims.Capabilities)
				c.Locals("claims", claims)
			}
		}
	}
	return c.Next()
}

// ============================================================================
// ROLE & CAPABILITY-BASED AUTHORIZATION MIDDLEWARE
// ============================================================================

// RequireRole checks that the authenticated user has the expected role or capability.
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Admin role is strictly isolated from marketplace capabilities
		if role == "admin" {
			if c.Locals("role") == "admin" {
				sessID, _ := c.Locals("sessionID").(int)
				if sessID <= 0 {
					return fiber.NewError(fiber.StatusUnauthorized, "active administrative session required")
				}
				return c.Next()
			}
			return fiber.NewError(fiber.StatusForbidden, "insufficient permissions: admin capability required")
		}

		// When capability claims are present, active capabilities govern authorization
		if claims, ok := c.Locals("claims").(*security.Claims); ok && claims != nil && len(claims.Capabilities) > 0 {
			if claims.HasCapability(role) {
				return c.Next()
			}
			return fiber.NewError(fiber.StatusForbidden, "insufficient permissions: "+role+" capability required")
		}

		// Fallback: Check capabilities slice in locals
		if caps, ok := c.Locals("capabilities").([]string); ok && len(caps) > 0 {
			for _, capItem := range caps {
				if capItem == role {
					return c.Next()
				}
			}
			return fiber.NewError(fiber.StatusForbidden, "insufficient permissions: "+role+" capability required")
		}

		// Legacy tokens without capabilities: fallback to direct role check
		if c.Locals("role") == role {
			return c.Next()
		}

		return fiber.NewError(fiber.StatusForbidden, "insufficient permissions: "+role+" capability required")
	}
}

// RequireCapability enforces that the authenticated user possesses an active marketplace capability.
func RequireCapability(capability string) fiber.Handler {
	return RequireRole(capability)
}

// RequireSuperAdmin ensures the admin has super or platform_owner role.
func RequireSuperAdmin(c *fiber.Ctx) error {
	adminRole, _ := c.Locals("adminRole").(string)
	if adminRole != models.RolePlatformOwner && adminRole != models.LegacyRoleSuper {
		return fiber.NewError(fiber.StatusForbidden, "platform owner access required")
	}
	return c.Next()
}

// RequirePlatformOwner enforces platform owner governance authorization.
func RequirePlatformOwner(c *fiber.Ctx) error {
	return RequireSuperAdmin(c)
}

// RequirePlatformAdmin ensures the administrator has Platform Administrator or Platform Owner access.
func RequirePlatformAdmin(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	if role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "admin access required")
	}
	adminRole, _ := c.Locals("adminRole").(string)
	if adminRole != models.RolePlatformOwner && adminRole != models.RolePlatformAdmin && adminRole != models.LegacyRoleSuper {
		return fiber.NewError(fiber.StatusForbidden, "platform administrator access required")
	}
	return c.Next()
}

// RequireSecurityAdmin ensures the administrator has Security Administrator or Platform Owner access.
// Strictly forbids platform administrators, operations admins, moderators, and support agents.
func RequireSecurityAdmin(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	if role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "admin access required")
	}
	adminRole, _ := c.Locals("adminRole").(string)
	if adminRole != models.RolePlatformOwner && adminRole != models.RoleSecurityAdmin && adminRole != models.LegacyRoleSuper {
		return fiber.NewError(fiber.StatusForbidden, "security administrator access required")
	}
	return c.Next()
}

// RequireExplicitAdminRole checks if the caller's active database role matches any of the explicitly allowed roles.
func RequireExplicitAdminRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "admin" {
			return fiber.NewError(fiber.StatusForbidden, "admin access required")
		}
		adminRole, _ := c.Locals("adminRole").(string)
		if adminRole == models.RolePlatformOwner || adminRole == models.LegacyRoleSuper {
			return c.Next()
		}
		for _, allowed := range allowedRoles {
			if adminRole == allowed {
				return c.Next()
			}
		}
		return fiber.NewError(fiber.StatusForbidden, "insufficient administrative permissions for this domain")
	}
}

// RequireAdminLevel routes legacy minRole checks to explicit domain authorization middleware.
// Deprecated: use RequirePlatformOwner, RequirePlatformAdmin, RequireSecurityAdmin, or RequireExplicitAdminRole.
func RequireAdminLevel(minRole string) fiber.Handler {
	switch minRole {
	case models.RolePlatformOwner, models.LegacyRoleSuper:
		return RequireSuperAdmin
	case models.RoleSecurityAdmin:
		return RequireSecurityAdmin
	case models.RolePlatformAdmin:
		return RequirePlatformAdmin
	default:
		return RequireExplicitAdminRole(minRole)
	}
}

// RequirePortal enforces portal-to-role access rules
func RequirePortal(portal string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "admin" {
			return fiber.NewError(fiber.StatusForbidden, "admin access required")
		}
		adminRole, _ := c.Locals("adminRole").(string)
		switch portal {
		case "executive":
			if adminRole != models.RolePlatformOwner && adminRole != models.LegacyRoleSuper {
				return fiber.NewError(fiber.StatusForbidden, "executive portal requires Platform Owner role")
			}
		case "platform":
			if adminRole != models.RolePlatformOwner && adminRole != models.RolePlatformAdmin && adminRole != models.LegacyRoleSuper {
				return fiber.NewError(fiber.StatusForbidden, "platform console requires Platform Administrator role")
			}
		case "security":
			if adminRole != models.RolePlatformOwner && adminRole != models.RoleSecurityAdmin && adminRole != models.LegacyRoleSuper {
				return fiber.NewError(fiber.StatusForbidden, "security console requires Security Administrator role")
			}
		}
		return c.Next()
	}
}

// RequireReauth enforces fresh re-authentication for high-risk administrative operations
func (m *Middleware) RequireReauth(scope string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(int)
		if !ok || userID == 0 {
			return fiber.NewError(fiber.StatusUnauthorized, "authentication required")
		}
		token := c.Get("X-Reauth-Token")
		if token == "" {
			return fiber.NewError(fiber.StatusPreconditionRequired, "re-authentication required: provide X-Reauth-Token")
		}
		if m.reauthValidator == nil {
			return fiber.NewError(fiber.StatusInternalServerError, "reauth validator not configured")
		}
		valid, err := m.reauthValidator(c.Context(), token, scope, userID)
		if err != nil || !valid {
			return fiber.NewError(fiber.StatusForbidden, "invalid or expired re-authentication token")
		}
		return c.Next()
	}
}

// ============================================================================
// PERMISSION-BASED AUTHORIZATION
// ============================================================================

// RequirePermission checks that an admin has the given page permission.
// Platform Owner and legacy super admins always pass.
// Enterprise roles automatically receive their designated domain baseline permissions.
func (m *Middleware) RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "admin" {
			return fiber.NewError(fiber.StatusForbidden, "admin access required")
		}
		adminRole, _ := c.Locals("adminRole").(string)
		if adminRole == models.RolePlatformOwner || adminRole == models.LegacyRoleSuper {
			return c.Next()
		}
		// Check enterprise role domain baseline permissions
		if domainPerms, ok := models.RoleDomainPermissions[adminRole]; ok && domainPerms[permission] {
			return c.Next()
		}
		// Otherwise verify explicitly granted permission via database checker
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

// RequirePermissionOrSuper allows super admins / platform owners unconditionally, and normal admins with permission.
func (m *Middleware) RequirePermissionOrSuper(permission string) fiber.Handler {
	return m.RequirePermission(permission)
}

// RequireAnyPermission checks that an admin has at least one of the given permissions.
// Platform Owner and Legacy Super always pass.
// Enterprise roles automatically check designated domain baseline permissions.
func (m *Middleware) RequireAnyPermission(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "admin" {
			return fiber.NewError(fiber.StatusForbidden, "admin access required")
		}
		adminRole, _ := c.Locals("adminRole").(string)
		if adminRole == models.RolePlatformOwner || adminRole == models.LegacyRoleSuper {
			return c.Next()
		}
		for _, perm := range permissions {
			if domainPerms, ok := models.RoleDomainPermissions[adminRole]; ok && domainPerms[perm] {
				return c.Next()
			}
			adminID, _ := c.Locals("userID").(int)
			if m.permissionChecker != nil {
				if ok, err := m.permissionChecker(adminID, perm); err == nil && ok {
					return c.Next()
				}
			}
		}
		return fiber.NewError(fiber.StatusForbidden, "insufficient permissions: required permission not granted")
	}
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
