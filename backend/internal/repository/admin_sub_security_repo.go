package repository

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type SecurityZoneMember struct {
	AdminID   int        `json:"admin_id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	AdminRole string     `json:"admin_role"`
	Notes     string     `json:"notes"`
	GrantedAt time.Time  `json:"granted_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type SecurityZone struct {
	ID                   int                  `json:"id"`
	Code                 string               `json:"code"`
	Name                 string               `json:"name"`
	Description          string               `json:"description"`
	RequiredRole         string               `json:"required_role"`
	IPWhitelist          []string             `json:"ip_whitelist"`
	RequireTrustedDevice bool                 `json:"require_trusted_device"`
	RequirePasskey       bool                 `json:"require_passkey"`
	IsActive             bool                 `json:"is_active"`
	Members              []SecurityZoneMember `json:"members"`
	CreatedAt            time.Time            `json:"created_at"`
}

type ActiveAdminSession struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	AdminRole          string    `json:"admin_role"`
	Portal             string    `json:"portal"`
	IPAddress          string    `json:"ip_address"`
	GeoLocation        string    `json:"geo_location"`
	DeviceID           string    `json:"device_id"`
	DeviceName         string    `json:"device_name"`
	ShiftCode          string    `json:"shift_code"`
	WorkstationCode    string    `json:"workstation_code"`
	FacilityLocation   string    `json:"facility_location"`
	Browser            string    `json:"browser"`
	OS                 string    `json:"os"`
	IsTrusted          bool      `json:"is_trusted"`
	LastActiveAt       time.Time `json:"last_active_at"`
	IdleTimeoutSeconds int       `json:"idle_timeout_seconds"`
	AbsoluteExpiresAt  time.Time `json:"absolute_expires_at"`
}

// GetAdminSecurityTelemetry returns complete audit telemetry for the given admin ID
func (r *AdminRepo) GetAdminSecurityTelemetry(ctx context.Context, adminID int, clientIP, userAgent, clientDeviceID string) (map[string]interface{}, error) {
	var name, email, phone, adminRole, defaultShift, assignedWorkstation string
	var isActive bool
	err := r.db.QueryRow(ctx, `
		SELECT name, COALESCE(email, ''), phone, COALESCE(admin_role::text, 'normal'),
		       COALESCE(default_shift, 'DAY_SHIFT'), COALESCE(assigned_workstation, 'WS-CORE-01'), is_active
		FROM users
		WHERE id = $1 AND role = 'admin'
	`, adminID).Scan(&name, &email, &phone, &adminRole, &defaultShift, &assignedWorkstation, &isActive)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admin profile: %w", err)
	}

	// Retrieve assigned sub-security zones
	zoneRows, err := r.db.Query(ctx, `
		SELECT z.id, z.code, z.name, z.description, za.granted_at, za.notes
		FROM admin_security_zone_assignments za
		JOIN admin_security_zones z ON za.zone_id = z.id
		WHERE za.admin_id = $1 AND z.is_active = true
		ORDER BY z.id ASC
	`, adminID)
	assignedZones := []map[string]interface{}{}
	if err == nil {
		defer zoneRows.Close()
		for zoneRows.Next() {
			var zid int
			var zcode, zname, zdesc, znotes string
			var grantedAt time.Time
			if err := zoneRows.Scan(&zid, &zcode, &zname, &zdesc, &grantedAt, &znotes); err == nil {
				assignedZones = append(assignedZones, map[string]interface{}{
					"id":          zid,
					"code":        zcode,
					"name":        zname,
					"description": zdesc,
					"granted_at":  grantedAt,
					"notes":       znotes,
				})
			}
		}
	}

	// Retrieve permissions
	permRows, err := r.db.Query(ctx, `
		SELECT permission FROM admin_permissions
		WHERE admin_id = $1 AND enabled = true
	`, adminID)
	permissions := []string{}
	if err == nil {
		defer permRows.Close()
		for permRows.Next() {
			var p string
			if err := permRows.Scan(&p); err == nil {
				permissions = append(permissions, p)
			}
		}
	}

	// Device identity check
	deviceName := "Registered Workstation"
	deviceType := "desktop"
	browser := "Chrome"
	os := "Linux/Unix"
	isTrusted := true

	if clientDeviceID != "" {
		_ = r.db.QueryRow(ctx, `
			SELECT device_name, device_type, COALESCE(browser, 'Browser'), COALESCE(os, 'OS'), is_trusted
			FROM admin_device_registrations
			WHERE device_id = $1 AND user_id = $2 AND revoked_at IS NULL
		`, clientDeviceID, adminID).Scan(&deviceName, &deviceType, &browser, &os, &isTrusted)
	} else {
		// Fallback to most recent registered device
		_ = r.db.QueryRow(ctx, `
			SELECT device_id, device_name, device_type, COALESCE(browser, 'Browser'), COALESCE(os, 'OS'), is_trusted
			FROM admin_device_registrations
			WHERE user_id = $1 AND revoked_at IS NULL
			ORDER BY last_active_at DESC LIMIT 1
		`, adminID).Scan(&clientDeviceID, &deviceName, &deviceType, &browser, &os, &isTrusted)
	}

	if clientDeviceID == "" {
		clientDeviceID = fmt.Sprintf("DEV-AUTH-%d", adminID)
	}

	// Geo-location heuristic (Ethio Telecom IP detection)
	geoLocation := "Ethiopia / Addis Ababa Central Metro"
	if strings.HasPrefix(clientIP, "197.156.") || strings.HasPrefix(clientIP, "213.55.") {
		geoLocation = "Ethio Telecom National Backbone / Addis Ababa"
	} else if clientIP == "127.0.0.1" || clientIP == "::1" || strings.HasPrefix(clientIP, "10.") || strings.HasPrefix(clientIP, "192.168.") {
		geoLocation = "Secure Private Intranet / Core Cluster"
	}

	telemetry := map[string]interface{}{
		"admin_id":       adminID,
		"name":           name,
		"email":          email,
		"phone":          phone,
		"admin_role":     adminRole,
		"is_active":      isActive,
		"ip_address":     clientIP,
		"geo_location":   geoLocation,
		"user_agent":     userAgent,
		"device_id":      clientDeviceID,
		"device_name":    deviceName,
		"device_type":    deviceType,
		"browser":        browser,
		"os":             os,
		"is_trusted":     isTrusted,
		"shift_code":     defaultShift,
		"workstation_code": assignedWorkstation,
		"session_status": "ACTIVE_AUTHENTICATED",
		"security_zones": assignedZones,
		"permissions":    permissions,
		"telemetry_time": time.Now().UTC(),
	}

	return telemetry, nil
}

// ListActiveAdminSessions returns all active, non-revoked administrator sessions
func (r *AdminRepo) ListActiveAdminSessions(ctx context.Context) ([]ActiveAdminSession, error) {
	query := `
		SELECT s.id, s.user_id, u.name, COALESCE(u.email, ''), COALESCE(u.admin_role::text, 'normal'),
		       s.portal, COALESCE(s.ip_address, '127.0.0.1'),
		       COALESCE(s.geo_location, 'Ethiopia / Addis Ababa Metro'),
		       COALESCE(s.device_id, 'UNKNOWN'),
		       COALESCE(d.device_name, 'Standard Terminal'),
		       COALESCE(s.shift_code, u.default_shift, 'DAY_SHIFT'),
		       COALESCE(s.workstation_code, u.assigned_workstation, 'WS-CORE-01'),
		       COALESCE(w.facility_location, 'Addis Ababa HQ Command Center'),
		       COALESCE(d.browser, 'Web Browser'),
		       COALESCE(d.os, 'Workstation'),
		       COALESCE(d.is_trusted, true),
		       s.last_active_at, s.idle_timeout_seconds, s.absolute_expires_at
		FROM admin_sessions s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN admin_device_registrations d ON s.device_id = d.device_id
		LEFT JOIN admin_workstations w ON COALESCE(s.workstation_code, u.assigned_workstation) = w.code
		WHERE s.revoked_at IS NULL AND s.absolute_expires_at > NOW()
		ORDER BY s.last_active_at DESC
		LIMIT 100
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []ActiveAdminSession
	for rows.Next() {
		var s ActiveAdminSession
		err := rows.Scan(
			&s.ID, &s.UserID, &s.Name, &s.Email, &s.AdminRole,
			&s.Portal, &s.IPAddress, &s.GeoLocation, &s.DeviceID,
			&s.DeviceName, &s.ShiftCode, &s.WorkstationCode, &s.FacilityLocation,
			&s.Browser, &s.OS, &s.IsTrusted,
			&s.LastActiveAt, &s.IdleTimeoutSeconds, &s.AbsoluteExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, rows.Err()
}

// RevokeAdminSessionByID terminates a specific administrative session
func (r *AdminRepo) RevokeAdminSessionByID(ctx context.Context, sessionID, revokedByAdminID int, reason string) error {
	if reason == "" {
		reason = "Security Administrator Remote Termination"
	}

	var targetUserID int
	err := r.db.QueryRow(ctx, `
		UPDATE admin_sessions
		SET revoked_at = NOW(), revocation_reason = $1
		WHERE id = $2 AND revoked_at IS NULL
		RETURNING user_id
	`, reason, sessionID).Scan(&targetUserID)
	if err != nil {
		return fmt.Errorf("session not found or already terminated: %w", err)
	}

	// Immutable audit logging
	_, _ = r.db.Exec(ctx, `
		INSERT INTO admin_logs (admin_id, target_user_id, action, target_type, target_id, notes, details, created_at)
		VALUES ($1, $2, 'REVOKE_ADMIN_SESSION', 'admin_session', $3, $4, $5, NOW())
	`, revokedByAdminID, targetUserID, sessionID, reason, fmt.Sprintf("Session #%d terminated by admin #%d", sessionID, revokedByAdminID))

	return nil
}

// ListSecurityZones returns all configured Sub-Security Zones and their member administrators
func (r *AdminRepo) ListSecurityZones(ctx context.Context) ([]SecurityZone, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, code, name, description, required_role,
		       COALESCE(ip_whitelist, '{}'), require_trusted_device, require_passkey, is_active, created_at
		FROM admin_security_zones
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zones []SecurityZone
	zoneMap := make(map[int]int) // zone.id -> index in zones

	for rows.Next() {
		var z SecurityZone
		if err := rows.Scan(
			&z.ID, &z.Code, &z.Name, &z.Description, &z.RequiredRole,
			&z.IPWhitelist, &z.RequireTrustedDevice, &z.RequirePasskey, &z.IsActive, &z.CreatedAt,
		); err != nil {
			return nil, err
		}
		z.Members = []SecurityZoneMember{}
		zoneMap[z.ID] = len(zones)
		zones = append(zones, z)
	}

	// Fetch all members
	memberRows, err := r.db.Query(ctx, `
		SELECT za.zone_id, u.id, u.name, COALESCE(u.email, ''), COALESCE(u.admin_role::text, 'normal'),
		       COALESCE(za.notes, ''), za.granted_at, za.expires_at
		FROM admin_security_zone_assignments za
		JOIN users u ON za.admin_id = u.id
		ORDER BY za.granted_at ASC
	`)
	if err == nil {
		defer memberRows.Close()
		for memberRows.Next() {
			var zoneID, adminID int
			var name, email, role, notes string
			var grantedAt time.Time
			var expiresAt *time.Time
			if err := memberRows.Scan(&zoneID, &adminID, &name, &email, &role, &notes, &grantedAt, &expiresAt); err == nil {
				if idx, ok := zoneMap[zoneID]; ok {
					zones[idx].Members = append(zones[idx].Members, SecurityZoneMember{
						AdminID:   adminID,
						Name:      name,
						Email:     email,
						AdminRole: role,
						Notes:     notes,
						GrantedAt: grantedAt,
						ExpiresAt: expiresAt,
					})
				}
			}
		}
	}

	return zones, nil
}

// AssignSecurityZone assigns an administrator to a Sub-Security Zone
func (r *AdminRepo) AssignSecurityZone(ctx context.Context, zoneCode string, targetAdminID, grantedByAdminID int, notes string) error {
	var zoneID int
	err := r.db.QueryRow(ctx, `SELECT id FROM admin_security_zones WHERE code = $1`, zoneCode).Scan(&zoneID)
	if err != nil {
		return fmt.Errorf("security zone '%s' not found: %w", zoneCode, err)
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO admin_security_zone_assignments (zone_id, admin_id, granted_by, notes, granted_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (zone_id, admin_id) DO UPDATE
		SET notes = EXCLUDED.notes, granted_by = EXCLUDED.granted_by, granted_at = NOW()
	`, zoneID, targetAdminID, grantedByAdminID, notes)
	if err != nil {
		return fmt.Errorf("failed to assign security zone: %w", err)
	}

	// Immutable audit logging
	_, _ = r.db.Exec(ctx, `
		INSERT INTO admin_logs (admin_id, target_user_id, action, target_type, target_id, notes, details, created_at)
		VALUES ($1, $2, 'ASSIGN_SECURITY_ZONE', 'security_zone', $3, $4, $5, NOW())
	`, grantedByAdminID, targetAdminID, zoneID, notes, fmt.Sprintf("Admin #%d assigned to zone %s", targetAdminID, zoneCode))

	return nil
}

// RevokeSecurityZone removes an administrator from a Sub-Security Zone
func (r *AdminRepo) RevokeSecurityZone(ctx context.Context, zoneCode string, targetAdminID, revokedByAdminID int, reason string) error {
	var zoneID int
	err := r.db.QueryRow(ctx, `SELECT id FROM admin_security_zones WHERE code = $1`, zoneCode).Scan(&zoneID)
	if err != nil {
		return fmt.Errorf("security zone '%s' not found: %w", zoneCode, err)
	}

	res, err := r.db.Exec(ctx, `
		DELETE FROM admin_security_zone_assignments
		WHERE zone_id = $1 AND admin_id = $2
	`, zoneID, targetAdminID)
	if err != nil {
		return fmt.Errorf("failed to revoke security zone assignment: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("assignment not found for admin #%d in zone %s", targetAdminID, zoneCode)
	}

	// Immutable audit logging
	_, _ = r.db.Exec(ctx, `
		INSERT INTO admin_logs (admin_id, target_user_id, action, target_type, target_id, notes, details, created_at)
		VALUES ($1, $2, 'REVOKE_SECURITY_ZONE', 'security_zone', $3, $4, $5, NOW())
	`, revokedByAdminID, targetAdminID, zoneID, reason, fmt.Sprintf("Admin #%d revoked from zone %s by admin #%d", targetAdminID, zoneCode, revokedByAdminID))

	return nil
}
