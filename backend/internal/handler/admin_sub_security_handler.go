package handlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

// GET /api/admin/security/telemetry/me
// Returns real-time administrative telemetry: Identity tuple, IP, Device ID,
// Trust level, Session validity, and assigned Sub-Security Zones.
func (h *AdminHandler) GetSecurityTelemetryMe(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	clientIP := c.IP()
	if xff := c.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 && strings.TrimSpace(parts[0]) != "" {
			clientIP = strings.TrimSpace(parts[0])
		}
	}

	userAgent := c.Get("User-Agent")
	clientDeviceID := c.Get("X-Device-Fingerprint")
	if clientDeviceID == "" {
		clientDeviceID = c.Query("device_id", "")
	}

	data, err := h.svc.GetAdminSecurityTelemetry(c.Context(), adminID, clientIP, userAgent, clientDeviceID)
	if err != nil {
		return response.InternalError(c, fmt.Errorf("failed to compute administrative security telemetry: %w", err))
	}

	return response.OK(c, data)
}

// GET /api/admin/security/active-sessions
// Lists all active administrative sessions across all consoles with IP and device telemetry.
func (h *AdminHandler) GetActiveAdminSessions(c *fiber.Ctx) error {
	sessions, err := h.svc.ListActiveAdminSessions(c.Context())
	if err != nil {
		return response.InternalError(c, fmt.Errorf("failed to retrieve active admin sessions: %w", err))
	}
	return response.OKList(c, sessions)
}

type RevokeSessionRequest struct {
	Reason string `json:"reason"`
}

// POST /api/admin/security/sessions/:id/revoke
// Instantly terminates an active administrative session across all clusters.
func (h *AdminHandler) RevokeAdminSession(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	sessionID, err := parseIDParam(c, "id")
	if err != nil {
		return response.BadRequest(c, "invalid session id")
	}

	var req RevokeSessionRequest
	_ = c.BodyParser(&req)
	if req.Reason == "" {
		req.Reason = "Manual revocation via Security Console"
	}

	if err := h.svc.RevokeAdminSessionByID(c.Context(), sessionID, adminID, req.Reason); err != nil {
		return response.InternalError(c, fmt.Errorf("failed to revoke admin session: %w", err))
	}

	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    "ADMIN_SESSION_REVOKED",
		Notes:     fmt.Sprintf("Revoked SessionID=%d Reason=%s", sessionID, req.Reason),
		IPAddress: c.IP(),
	})

	return response.OK(c, fiber.Map{
		"message":    "Administrative session revoked successfully",
		"session_id": sessionID,
	})
}

// GET /api/admin/security/zones
// Returns all Sub-Security Zones and their authorized administrative personnel.
func (h *AdminHandler) ListSecurityZones(c *fiber.Ctx) error {
	zones, err := h.svc.ListSecurityZones(c.Context())
	if err != nil {
		return response.InternalError(c, fmt.Errorf("failed to retrieve sub-security zones: %w", err))
	}
	return response.OKList(c, zones)
}

type AssignZoneRequest struct {
	ZoneCode      string `json:"zone_code"`
	TargetAdminID int    `json:"target_admin_id"`
	Notes         string `json:"notes"`
}

// POST /api/admin/security/zones/assign
// Delegated by Security Administrator / CISO to assign an administrator to a specific Sub-Security Zone.
func (h *AdminHandler) AssignSecurityZone(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	var req AssignZoneRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	req.ZoneCode = strings.ToUpper(strings.TrimSpace(req.ZoneCode))
	if req.ZoneCode == "" || req.TargetAdminID <= 0 {
		return response.BadRequest(c, "zone_code and valid target_admin_id are required")
	}

	if err := h.svc.AssignSecurityZone(c.Context(), req.ZoneCode, req.TargetAdminID, adminID, req.Notes); err != nil {
		return response.BadRequest(c, err.Error())
	}

	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:      adminID,
		TargetUserID: &req.TargetAdminID,
		Action:       "ASSIGN_SECURITY_ZONE",
		Notes:        fmt.Sprintf("Zone=%s AssignedToAdmin=%d Notes=%s", req.ZoneCode, req.TargetAdminID, req.Notes),
		IPAddress:    c.IP(),
	})

	return response.OK(c, fiber.Map{
		"message":         "Sub-security zone delegated successfully",
		"zone_code":       req.ZoneCode,
		"target_admin_id": req.TargetAdminID,
	})
}

type RevokeZoneRequest struct {
	ZoneCode      string `json:"zone_code"`
	TargetAdminID int    `json:"target_admin_id"`
	Reason        string `json:"reason"`
}

// POST /api/admin/security/zones/revoke
// Revokes an administrator's assignment from a Sub-Security Zone.
func (h *AdminHandler) RevokeSecurityZone(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	var req RevokeZoneRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	req.ZoneCode = strings.ToUpper(strings.TrimSpace(req.ZoneCode))
	if req.ZoneCode == "" || req.TargetAdminID <= 0 {
		return response.BadRequest(c, "zone_code and valid target_admin_id are required")
	}
	if req.Reason == "" {
		req.Reason = "Revoked by Security Officer"
	}

	if err := h.svc.RevokeSecurityZone(c.Context(), req.ZoneCode, req.TargetAdminID, adminID, req.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}

	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:      adminID,
		TargetUserID: &req.TargetAdminID,
		Action:       "REVOKE_SECURITY_ZONE",
		Notes:        fmt.Sprintf("Zone=%s RevokedFromAdmin=%d Reason=%s", req.ZoneCode, req.TargetAdminID, req.Reason),
		IPAddress:    c.IP(),
	})

	return response.OK(c, fiber.Map{
		"message":         "Sub-security zone revoked successfully",
		"zone_code":       req.ZoneCode,
		"target_admin_id": req.TargetAdminID,
	})
}
