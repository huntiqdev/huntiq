package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

// AdminAuthHandler handles all privileged administrative authentication and passkey endpoints
type AdminAuthHandler struct {
	adminAuthSvc *service.AdminAuthService
	authSvc      *service.AuthService
	adminSvc     *service.AdminService
	secureCookie bool
}

// NewAdminAuthHandler creates a new AdminAuthHandler
func NewAdminAuthHandler(
	adminAuthSvc *service.AdminAuthService,
	authSvc *service.AuthService,
	adminSvc *service.AdminService,
	secureCookie bool,
) *AdminAuthHandler {
	return &AdminAuthHandler{
		adminAuthSvc: adminAuthSvc,
		authSvc:      authSvc,
		adminSvc:     adminSvc,
		secureCookie: secureCookie,
	}
}

// Helper to set the HttpOnly session cookie
func (h *AdminAuthHandler) setSessionCookie(c *fiber.Ctx, sessionToken string, maxAgeSec int) {
	c.Cookie(&fiber.Cookie{
		Name:     "hm_admin_session",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   maxAgeSec,
		HTTPOnly: true,
		Secure:   h.secureCookie,
		SameSite: "Lax",
	})
}

// Helper to clear the session cookie
func (h *AdminAuthHandler) clearSessionCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "hm_admin_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   h.secureCookie,
		SameSite: "Lax",
	})
}

// extractClientIP retrieves client IP without trusting arbitrary forwarded headers
func extractClientIP(c *fiber.Ctx) string {
	ip := c.IP()
	if ip == "" {
		ip = "127.0.0.1"
	}
	return ip
}

// ============================================================================
// WEBAUTHN PASSKEY REGISTRATION
// ============================================================================

// POST /api/admin/passkeys/register/begin
func (h *AdminAuthHandler) BeginRegistration(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok || userID == 0 {
		return response.Unauthorized(c, "authentication required")
	}

	portal := c.Query("portal", "staff")
	options, challenge, err := h.adminAuthSvc.BeginPasskeyRegistration(c.Context(), userID, portal)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"options":      options,
		"challenge":    challenge.Challenge,
		"challenge_id": challenge.Challenge,
	})
}

// POST /api/admin/passkeys/register/finish (also /admin/auth/passkey/register/verify)
func (h *AdminAuthHandler) FinishRegistration(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok || userID == 0 {
		return response.Unauthorized(c, "authentication required")
	}

	var req struct {
		Challenge    string          `json:"challenge"`
		ChallengeID  string          `json:"challenge_id"`
		Name         string          `json:"name"`
		FriendlyName string          `json:"friendly_name"`
		Response     json.RawMessage `json:"response"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid registration payload")
	}

	challengeStr := req.Challenge
	if challengeStr == "" {
		challengeStr = req.ChallengeID
	}
	name := req.Name
	if name == "" {
		name = req.FriendlyName
	}
	if name == "" {
		name = "Authenticator"
	}

	if challengeStr == "" || len(req.Response) == 0 {
		return response.BadRequest(c, "challenge and response are required")
	}

	httpReq, _ := http.NewRequestWithContext(c.Context(), "POST", "", bytes.NewReader(req.Response))
	httpReq.Header.Set("Content-Type", "application/json")
	parsedCreation, err := protocol.ParseCredentialCreationResponse(httpReq)
	if err != nil {
		return response.BadRequest(c, "failed to parse attestation response: "+err.Error())
	}

	passkey, err := h.adminAuthSvc.FinishPasskeyRegistration(c.Context(), userID, challengeStr, name, parsedCreation)
	if err != nil {
		return response.BadRequest(c, "registration failed: "+err.Error())
	}

	// Audit enrollment
	h.adminSvc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   userID,
		Action:    "PASSKEY_ENROLLED",
		Notes:     fmt.Sprintf("Registered authenticator %s (%s)", passkey.Name, passkey.Attachment),
		IPAddress: extractClientIP(c),
	})

	return response.OK(c, fiber.Map{
		"success":       true,
		"message":       "Passkey successfully enrolled",
		"passkey":       passkey,
		"credential_id": passkey.ID,
	})
}

// ============================================================================
// WEBAUTHN PASSKEY AUTHENTICATION (LOGIN)
// ============================================================================

// POST /api/admin/passkeys/login/begin (also /admin/auth/passkey/login/begin)
func (h *AdminAuthHandler) BeginLogin(c *fiber.Ctx) error {
	var req struct {
		Identifier string `json:"identifier"`
		Email      string `json:"email"`
		Portal     string `json:"portal"`
	}
	_ = c.BodyParser(&req)

	identifier := req.Identifier
	if identifier == "" {
		identifier = req.Email
	}

	if req.Portal == "" {
		req.Portal = "staff"
	}

	ip := extractClientIP(c)
	allowed, waitSec, _ := h.adminAuthSvc.CheckThrottling(c.Context(), strings.ToLower(identifier), ip)
	if !allowed {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"success": false,
			"error":   fmt.Sprintf("Too many login attempts. Please wait %d seconds.", waitSec),
		})
	}

	assertion, challenge, err := h.adminAuthSvc.BeginPasskeyLogin(c.Context(), identifier, req.Portal)
	if err != nil {
		return response.BadRequest(c, "failed to initialize passkey assertion: "+err.Error())
	}

	return response.OK(c, fiber.Map{
		"options":      assertion,
		"challenge":    challenge.Challenge,
		"challenge_id": challenge.Challenge,
	})
}

// POST /api/admin/passkeys/login/finish (also /admin/auth/passkey/login/verify)
func (h *AdminAuthHandler) FinishLogin(c *fiber.Ctx) error {
	var req struct {
		Challenge   string          `json:"challenge"`
		ChallengeID string          `json:"challenge_id"`
		Portal      string          `json:"portal"`
		Response    json.RawMessage `json:"response"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid assertion payload")
	}

	challengeStr := req.Challenge
	if challengeStr == "" {
		challengeStr = req.ChallengeID
	}
	if challengeStr == "" || len(req.Response) == 0 {
		return response.BadRequest(c, "challenge and response are required")
	}

	if req.Portal == "" {
		req.Portal = "staff"
	}

	ip := extractClientIP(c)
	userAgent := c.Get("User-Agent")

	httpReq, _ := http.NewRequestWithContext(c.Context(), "POST", "", bytes.NewReader(req.Response))
	httpReq.Header.Set("Content-Type", "application/json")
	parsedAssertion, err := protocol.ParseCredentialRequestResponse(httpReq)
	if err != nil {
		h.adminAuthSvc.RecordAuthResult(c.Context(), "passkey", ip, req.Portal, false, "assertion_parse_failure")
		return response.BadRequest(c, "failed to parse credential assertion: "+err.Error())
	}

	sessionToken, sessionID, user, err := h.adminAuthSvc.FinishPasskeyLogin(
		c.Context(),
		challengeStr,
		req.Portal,
		ip,
		userAgent,
		parsedAssertion,
	)
	if err != nil {
		h.adminAuthSvc.RecordAuthResult(c.Context(), "passkey", ip, req.Portal, false, err.Error())
		return response.Unauthorized(c, "Passkey authentication failed: "+err.Error())
	}

	// Record success for throttling
	h.adminAuthSvc.RecordAuthResult(c.Context(), user.Email, ip, req.Portal, true, "passkey_success")

	// Generate session-bound JWT
	jwtToken, _ := h.adminAuthSvc.GenerateAdminJWT(user.ID, user.AdminRole, sessionID, sessionToken)

	// Set HttpOnly cookie
	maxAge := 12 * 3600 // 12 hours
	if user.AdminRole == models.RolePlatformOwner || req.Portal == "executive" {
		maxAge = 4 * 3600 // 4 hours for executive
	}
	h.setSessionCookie(c, sessionToken, maxAge)

	// Audit login
	h.adminSvc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   user.ID,
		Action:    "AUTH_PASSKEY_SUCCESS",
		Notes:     fmt.Sprintf("Passkey login successful via portal: %s", req.Portal),
		IPAddress: ip,
	})

	return response.OK(c, fiber.Map{
		"success":       true,
		"token":         jwtToken,
		"session_token": sessionToken,
		"role":          user.AdminRole,
		"redirect_url":  "/admin/dashboard",
		"user": fiber.Map{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"admin_role": user.AdminRole,
		},
		"portal": req.Portal,
	})
}

// ============================================================================
// SESSION MANAGEMENT
// ============================================================================

// POST /api/admin/logout
func (h *AdminAuthHandler) Logout(c *fiber.Ctx) error {
	sessionToken := c.Get("Authorization")
	if strings.HasPrefix(sessionToken, "Bearer ") {
		sessionToken = strings.TrimPrefix(sessionToken, "Bearer ")
	} else {
		sessionToken = c.Cookies("hm_admin_session")
	}

	if sessionToken != "" {
		_ = h.adminAuthSvc.RevokeSession(c.Context(), sessionToken, "logout")
	}

	h.clearSessionCookie(c)
	return response.OK(c, fiber.Map{"success": true, "message": "logged out successfully"})
}

// GET /api/admin/passkeys (also /admin/auth/passkey/list)
func (h *AdminAuthHandler) ListPasskeys(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok || userID == 0 {
		return response.Unauthorized(c, "authentication required")
	}

	list, err := h.adminAuthSvc.ListPasskeys(c.Context(), userID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	if list == nil {
		list = []models.PasskeyCredential{}
	}
	return response.OK(c, fiber.Map{
		"passkeys":       list,
		"authenticators": list,
	})
}

// DELETE /api/admin/passkeys/:id (also /admin/auth/passkey/:id)
func (h *AdminAuthHandler) RevokePasskey(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok || userID == 0 {
		return response.Unauthorized(c, "authentication required")
	}

	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return response.BadRequest(c, "invalid passkey id")
	}

	if err := h.adminAuthSvc.RevokePasskey(c.Context(), id, userID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	h.adminSvc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   userID,
		Action:    "PASSKEY_REVOKED",
		Notes:     fmt.Sprintf("Revoked passkey id %d", id),
		IPAddress: extractClientIP(c),
	})

	return response.OK(c, fiber.Map{"success": true, "message": "passkey revoked"})
}

// POST /api/admin/reauth/begin
func (h *AdminAuthHandler) BeginReauth(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok || userID == 0 {
		return response.Unauthorized(c, "authentication required")
	}

	var req struct {
		ActionScope string `json:"action_scope"`
	}
	_ = c.BodyParser(&req)
	if req.ActionScope == "" {
		req.ActionScope = "sensitive_operation"
	}

	portal, _ := c.Locals("portal").(string)
	if portal == "" {
		portal = "staff"
	}

	assertion, challenge, err := h.adminAuthSvc.BeginPasskeyLogin(c.Context(), "", portal)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"options":      assertion,
		"challenge":    challenge.Challenge,
		"action_scope": req.ActionScope,
	})
}

// POST /api/admin/reauth/finish
func (h *AdminAuthHandler) FinishReauth(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok || userID == 0 {
		return response.Unauthorized(c, "authentication required")
	}

	var req struct {
		Challenge   string          `json:"challenge"`
		ActionScope string          `json:"action_scope"`
		Response    json.RawMessage `json:"response"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid reauth payload")
	}

	if req.ActionScope == "" {
		req.ActionScope = "sensitive_operation"
	}

	ip := extractClientIP(c)
	httpReq, _ := http.NewRequestWithContext(c.Context(), "POST", "", bytes.NewReader(req.Response))
	httpReq.Header.Set("Content-Type", "application/json")
	parsedAssertion, err := protocol.ParseCredentialRequestResponse(httpReq)
	if err != nil {
		return response.BadRequest(c, "parse assertion failed: "+err.Error())
	}

	// Verify assertion for user
	_, _, _, err = h.adminAuthSvc.FinishPasskeyLogin(c.Context(), req.Challenge, "staff", ip, c.Get("User-Agent"), parsedAssertion)
	if err != nil {
		return response.Unauthorized(c, "reauth assertion failed: "+err.Error())
	}

	reauthToken, err := h.adminAuthSvc.IssueReauthToken(c.Context(), userID, req.ActionScope)
	if err != nil {
		return response.BadRequest(c, "failed to issue reauth token: "+err.Error())
	}

	return response.OK(c, fiber.Map{
		"success":      true,
		"reauth_token": reauthToken,
	})
}

// POST /api/admin/auth/reauth — password-based re-authentication endpoint
func (h *AdminAuthHandler) ReauthPassword(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok || userID == 0 {
		return response.Unauthorized(c, "authentication required")
	}

	var req struct {
		Password string `json:"password"`
		Purpose  string `json:"purpose"`
		Scope    string `json:"scope"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid reauth payload")
	}

	scope := req.Scope
	if scope == "" {
		scope = req.Purpose
	}
	if scope == "" {
		scope = "sensitive_operation"
	}

	if req.Password == "" {
		return response.BadRequest(c, "password is required")
	}

	valid, err := h.adminAuthSvc.ValidateAdminPassword(c.Context(), userID, req.Password)
	if err != nil || !valid {
		return response.Unauthorized(c, "invalid administrative password")
	}

	reauthToken, err := h.adminAuthSvc.IssueReauthToken(c.Context(), userID, scope)
	if err != nil {
		return response.BadRequest(c, "failed to issue reauth token: "+err.Error())
	}

	return response.OK(c, fiber.Map{
		"success":            true,
		"reauth_token":       reauthToken,
		"expires_in_seconds": 600,
	})
}

// POST /api/admin/auth/session/revoke — session revocation endpoint
func (h *AdminAuthHandler) RevokeSessionEndpoint(c *fiber.Ctx) error {
	var req struct {
		SessionID string `json:"session_id"`
	}
	_ = c.BodyParser(&req)

	if req.SessionID != "" {
		sessID, err := strconv.Atoi(req.SessionID)
		if err == nil && sessID > 0 {
			userID, _ := c.Locals("userID").(int)
			_ = h.adminAuthSvc.RevokeSessionByID(c.Context(), sessID, userID, "revoked_via_api")
			return response.OK(c, fiber.Map{"success": true, "message": "session revoked"})
		}
	}

	// Default fallback: revoke current session
	return h.Logout(c)
}

