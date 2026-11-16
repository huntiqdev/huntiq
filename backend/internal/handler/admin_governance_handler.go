package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/pkg/response"
	"github.com/your-org/used-tools-marketplace/pkg/validations"
	"strconv"
	"strings"
	"time"
	apperrors "github.com/your-org/used-tools-marketplace/pkg/errors"
)

func (h *AdminHandler) GetMe(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	role := adminRole(c)
	if h.authSvc != nil {
		u, err := h.authSvc.GetProfile(c.Context(), adminID)
		if err == nil && u != nil {
			return response.OK(c, fiber.Map{
				"id":         u.ID,
				"name":       u.Name,
				"phone":      u.Phone,
				"email":      u.Email,
				"admin_role": u.AdminRole,
			})
		}
	}
	return response.OK(c, fiber.Map{
		"id":         adminID,
		"name":       "Admin",
		"admin_role": role,
	})
}

// GET /api/admin/me/permissions — any admin fetches their own permissions
// Normal admins use this on login to load their page access list.
// Super admins get an empty list (they have all access by default).
func (h *AdminHandler) GetMyPermissions(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	role := adminRole(c)
	if role == models.RolePlatformOwner || role == models.LegacyRoleSuper {
		// Platform Owner / Super admins have all permissions by default
		return response.OK(c, fiber.Map{
			"permissions":       []string{},
			"is_super":          true,
			"is_platform_owner": true,
			"admin_role":        role,
		})
	}
	perms, err := h.svc.GetAdminPermissions(c.Context(), adminID)
	if err != nil {
		return response.InternalError(c, err)
	}
	if perms == nil {
		perms = []string{}
	}
	return response.OK(c, fiber.Map{
		"permissions":       perms,
		"is_super":          false,
		"is_platform_owner": false,
		"admin_role":        role,
	})
}

// ── Sub-Admin Management (Super Admin Only) ───────────────────────────────────

// GET /api/admin/admins
func (h *AdminHandler) ListAdmins(c *fiber.Ctx) error {
	admins, err := h.svc.ListAdmins(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, admins)
}

// POST /api/admin/admins
func (h *AdminHandler) CreateAdmin(c *fiber.Ctx) error {
	var body struct {
		Name      string `json:"name"`
		Phone     string `json:"phone"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		AdminRole string `json:"admin_role"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.Name == "" || body.Phone == "" || body.Password == "" {
		return response.BadRequest(c, "name, phone, and password are required")
	}
	if body.Email == "" {
		return response.BadRequest(c, "email is required — it is the unique identifier for administrative staff")
	}
	if valid, errKey := validations.ValidatePassword(body.Password); !valid {
		if errKey == "password_min" {
			return response.BadRequest(c, "password must be at least 8 characters")
		}
		return response.BadRequest(c, "invalid password")
	}

	role := body.AdminRole
	if role == "" {
		role = models.RoleOperationsAdmin
	}
	// Platform Owner accounts cannot be provisioned via web API; they require direct bootstrap
	if role == models.RolePlatformOwner || role == models.LegacyRoleSuper {
		return response.BadRequest(c, "Platform Owner accounts cannot be provisioned via the web API")
	}

	if _, ok := models.RoleDomainPermissions[role]; !ok && role != models.LegacyRoleNormal {
		return response.BadRequest(c, "invalid admin_role: must be a recognized enterprise administrative role")
	}

	id, err := h.svc.CreateAdmin(c.Context(), body.Name, body.Phone, body.Email, body.Password, role)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "users_phone_role_key") ||
			strings.Contains(errMsg, "users_phone_key") ||
			(strings.Contains(errMsg, "duplicate") && strings.Contains(errMsg, "phone")) {
			return response.BadRequestCode(c, "phone number already in use by another admin", apperrors.ErrAuthPhoneAlreadyRegistered)
		}
		if strings.Contains(errMsg, "users_email_key") ||
			(strings.Contains(errMsg, "duplicate") && strings.Contains(errMsg, "email")) {
			return response.BadRequestCode(c, "email address already in use", apperrors.ErrAuthEmailAlreadyInUse)
		}
		return response.BadRequest(c, errMsg)
	}
	callerID, _ := mustUserID(c)
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:      callerID,
		TargetUserID: &id,
		Action:       "ADMIN_CREATED",
		Notes:        fmt.Sprintf("Created administrator id=%d, role=%s, email=%s", id, role, body.Email),
		IPAddress:    c.IP(),
	})
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"id": id, "name": body.Name, "phone": body.Phone, "email": body.Email, "admin_role": role},
	})
}

// PUT /api/admin/admins/:id/role
func (h *AdminHandler) UpdateAdminRole(c *fiber.Ctx) error {
	callerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	callerRole := adminRole(c)
	if callerRole != models.RolePlatformOwner && callerRole != models.LegacyRoleSuper {
		return fiber.NewError(fiber.StatusForbidden, "only Platform Owner can modify administrator roles")
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	// Prevent Platform Owner from self-modifying role
	if id == callerID {
		return response.BadRequest(c, "cannot change your own role")
	}
	var body struct {
		AdminRole string `json:"admin_role"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.AdminRole == models.RolePlatformOwner || body.AdminRole == models.LegacyRoleSuper {
		return response.BadRequest(c, "Platform Owner promotion requires direct server bootstrap authorization")
	}
	if err := h.svc.UpdateAdminRole(c.Context(), id, body.AdminRole); err != nil {
		return response.BadRequest(c, err.Error())
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:      callerID,
		TargetUserID: &id,
		Action:       "ADMIN_ROLE_UPDATED",
		Notes:        fmt.Sprintf("Updated administrator id=%d to role=%s", id, body.AdminRole),
		IPAddress:    c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "admin role updated"})
}

// DELETE /api/admin/admins/:id
func (h *AdminHandler) DeleteAdmin(c *fiber.Ctx) error {
	callerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if id == callerID {
		return response.BadRequest(c, "cannot delete yourself")
	}
	if err := h.svc.DeleteAdmin(c.Context(), id); err != nil {
		// "super admins cannot be deleted" is a client logic error, not a server fault
		return response.BadRequest(c, err.Error())
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:      callerID,
		TargetUserID: &id,
		Action:       "ADMIN_DELETED",
		Notes:        fmt.Sprintf("Deleted administrator id=%d", id),
		IPAddress:    c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "admin deleted"})
}

// GET /api/admin/admins/:id/permissions
func (h *AdminHandler) GetAdminPermissions(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	perms, err := h.svc.GetAdminPermissions(c.Context(), id)
	if err != nil {
		// Service returns an error for super admins
		return response.BadRequest(c, err.Error())
	}
	if perms == nil {
		perms = []string{}
	}
	return response.OK(c, fiber.Map{"permissions": perms})
}

// PUT /api/admin/admins/:id/permissions
func (h *AdminHandler) SetAdminPermissions(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.Permissions == nil {
		body.Permissions = []string{}
	}
	if err := h.svc.SetAdminPermissions(c.Context(), id, body.Permissions); err != nil {
		// "super admins cannot be set" is a client logic error, not a server fault
		return response.BadRequest(c, err.Error())
	}
	callerID, _ := mustUserID(c)
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:      callerID,
		TargetUserID: &id,
		Action:       "ADMIN_PERMISSIONS_UPDATED",
		Notes:        fmt.Sprintf("Updated permissions for administrator id=%d (count=%d)", id, len(body.Permissions)),
		IPAddress:    c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "permissions updated"})
}

// ── Admin Activity Logs ───────────────────────────────────────────────────────

// GET /api/admin/logs
func (h *AdminHandler) ListAdminLogs(c *fiber.Ctx) error {
	const dateLayout = "2006-01-02"
	dateFrom := c.Query("from")
	dateTo := c.Query("to")
	if dateFrom != "" {
		if _, err := time.Parse(dateLayout, dateFrom); err != nil {
			return response.BadRequest(c, "invalid 'from' date — use YYYY-MM-DD")
		}
	}
	if dateTo != "" {
		if _, err := time.Parse(dateLayout, dateTo); err != nil {
			return response.BadRequest(c, "invalid 'to' date — use YYYY-MM-DD")
		}
	}
	// Sanitise IP: strip % wildcards and cap length to prevent oversized ILIKE patterns
	ip := strings.ReplaceAll(c.Query("ip"), "%", "")
	ip = truncate(ip, 45) // max IPv6 length is 45 chars
	f := models.AdminLogFilter{
		Action:    c.Query("action"),
		IPAddress: ip,
		DateFrom:  dateFrom,
		DateTo:    dateTo,
	}
	if v := c.Query("admin_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			f.AdminID = id
		}
	}
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			f.Page = p
		}
	}
	if v := c.Query("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil {
			f.Limit = l
		}
	}
	// Cap limit to prevent memory exhaustion from oversized result sets
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	logs, total, err := h.svc.ListAdminLogs(c.Context(), f)
	if err != nil {
		return response.InternalError(c, err)
	}
	if logs == nil {
		logs = []models.AdminLog{}
	}
	return response.OK(c, fiber.Map{"logs": logs, "total": total})
}

// DELETE /api/admin/logs  — bulk delete admin activity logs to free storage
// Query param: before=YYYY-MM-DD (delete logs older than date), omit = delete all
func (h *AdminHandler) DeleteAdminLogs(c *fiber.Ctx) error {
	before := c.Query("before", "")
	if before != "" {
		// Validate date format to prevent malformed input reaching the DB
		if _, err := time.Parse("2006-01-02", before); err != nil {
			return response.BadRequest(c, "invalid 'before' date — use YYYY-MM-DD format")
		}
	}
	deleted, err := h.svc.DeleteAdminLogs(c.Context(), before)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"deleted": deleted, "message": "admin logs cleaned up"})
}

// ── Tier Management ───────────────────────────────────────────────────────────

// PUT /api/admin/users/:id/seller-tier
func (h *AdminHandler) SetSellerTier(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	userID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Tier   string `json:"tier"`
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if err := h.svc.SetSellerTier(c.Context(), userID, body.Tier, adminID, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:      adminID,
		TargetUserID: &userID,
		Action:       "tier_change",
		Details:      fmt.Sprintf("seller tier → %s", body.Tier), //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
		Notes:        body.Reason,
		IPAddress:    c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "seller tier updated"})
}

// PUT /api/admin/users/:id/buyer-tier
func (h *AdminHandler) SetBuyerTier(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	userID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Tier   string `json:"tier"`
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if err := h.svc.SetBuyerTier(c.Context(), userID, body.Tier, adminID, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:      adminID,
		TargetUserID: &userID,
		Action:       "tier_change",
		Details:      fmt.Sprintf("buyer tier → %s", body.Tier), //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
		Notes:        body.Reason,
		IPAddress:    c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "buyer tier updated"})
}

// GET /api/admin/users/:id/tier-history
func (h *AdminHandler) GetTierHistory(c *fiber.Ctx) error {
	userID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	history, err := h.svc.ListTierHistory(c.Context(), userID)
	if err != nil {
		return response.InternalError(c, err)
	}
	if history == nil {
		history = []models.TierHistory{}
	}
	return response.OKList(c, history)
}

// GET /api/admin/premium-users
func (h *AdminHandler) ListPremiumUsers(c *fiber.Ctx) error {
	// Allowlist all three params — they control DB branch logic and query values
	role := c.Query("role", "seller")
	if role != "seller" && role != "buyer" {
		role = "seller"
	}
	tierType := c.Query("tier_type", "seller")
	if tierType != "seller" && tierType != "buyer" {
		tierType = "seller"
	}
	tier := c.Query("tier", "gold")
	if tier != "gold" && tier != "premium" && tier != "all" && tier != "normal" {
		tier = "gold"
	}
	users, err := h.svc.ListPremiumUsers(c.Context(), role, tierType, tier)
	if err != nil {
		return response.InternalError(c, err)
	}
	if users == nil {
		users = []models.User{}
	}
	return response.OKList(c, users)
}

// GET /api/admin/categories/popularity
func (h *AdminHandler) CategoryPopularity(c *fiber.Ctx) error {
	days := 7
	limit := 20
	if v := c.Query("days"); v != "" {
		if d, err := strconv.Atoi(v); err == nil && d >= 1 && d <= 365 {
			days = d
		}
	}
	if v := c.Query("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l >= 1 && l <= 100 {
			limit = l
		}
	}
	if h.prodRepo == nil {
		return response.BadRequest(c, "product repo not configured")
	}
	cats, err := h.prodRepo.PopularCategories(c.Context(), days, limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	if cats == nil {
		cats = []models.PopularCategory{}
	}
	return response.OKList(c, cats)
}

// ── Seller Upgrade Requests ───────────────────────────────────────────────────

// GET /api/admin/upgrade-requests?status=pending
func (h *AdminHandler) ListUpgradeRequests(c *fiber.Ctx) error {
	status := c.Query("status", "pending")
	requests, err := h.svc.ListUpgradeRequests(c.Context(), status)
	if err != nil {
		return response.InternalError(c, err)
	}
	if requests == nil {
		requests = []map[string]interface{}{}
	}
	return response.OKList(c, requests)
}

// PUT /api/admin/upgrade-requests/:id/approve
func (h *AdminHandler) ApproveUpgradeRequest(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.ApproveUpgradeRequest(c.Context(), id, adminID); err != nil {
		return response.BadRequest(c, err.Error())
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    "approve",
		Details:   fmt.Sprintf("upgrade request #%d approved", id), //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
		IPAddress: c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "upgrade approved — seller tier updated"})
}

// PUT /api/admin/upgrade-requests/:id/reject
func (h *AdminHandler) RejectUpgradeRequest(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Reason string `json:"reason"`
	}
	_ = c.BodyParser(&body) // reason is optional
	if err := h.svc.RejectUpgradeRequest(c.Context(), id, adminID, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    "reject",
		Details:   fmt.Sprintf("upgrade request #%d rejected", id), //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
		Notes:     body.Reason,
		IPAddress: c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "upgrade request rejected"})
}

// GET /api/public/promotion — returns the active promotion banner settings (no auth required)
func (h *AdminHandler) ListBuyerUpgradeRequests(c *fiber.Ctx) error {
	status := c.Query("status") // pending | approved | rejected | "" (all)
	page, limit := parsePagination(c, 50, 200)
	requests, total, err := h.svc.ListBuyerUpgradeRequests(c.Context(), status, page, limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"requests": requests, "total": total})
}

// PUT /api/admin/buyer-upgrade-requests/:id/approve
func (h *AdminHandler) ApproveBuyerUpgradeRequest(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.ApproveBuyerUpgradeRequest(c.Context(), id, adminID); err != nil {
		return response.BadRequest(c, err.Error())
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    "approve",
		Details:   fmt.Sprintf("buyer upgrade request #%d approved", id), //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
		IPAddress: c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "buyer upgrade approved — tier updated"})
}

// PUT /api/admin/buyer-upgrade-requests/:id/reject
func (h *AdminHandler) RejectBuyerUpgradeRequest(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Reason string `json:"reason"`
	}
	_ = c.BodyParser(&body)
	if err := h.svc.RejectBuyerUpgradeRequest(c.Context(), id, adminID, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    "reject",
		Details:   fmt.Sprintf("buyer upgrade request #%d rejected: %s", id, body.Reason), //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
		IPAddress: c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "buyer upgrade request rejected"})
}

// ── User Governance (Warnings & Bans) ─────────────────────────────────────────

// POST /api/admin/users/:id/warn
func (h *AdminHandler) IssueWarning(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	userID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil || body.Reason == "" {
		return response.BadRequest(c, "reason is required")
	}
	if err := h.svc.IssueWarning(c.Context(), userID, adminID, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:      adminID,
		TargetUserID: &userID,
		Action:       "warn",
		Notes:        body.Reason,
		IPAddress:    c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "warning issued to user"})
}

// PUT /api/admin/users/:id/ban
func (h *AdminHandler) BanUser(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	// Defense-in-depth: only super admins can ban (permanent exclusion)
	if adminRole(c) != models.RolePlatformOwner && adminRole(c) != models.LegacyRoleSuper {
		return fiber.NewError(fiber.StatusForbidden, "platform owner access required to ban users")
	}
	userID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil || body.Reason == "" {
		return response.BadRequest(c, "reason is required for banning")
	}
	if err := h.svc.BanUser(c.Context(), userID, adminID, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:      adminID,
		TargetUserID: &userID,
		Action:       "ban",
		Notes:        body.Reason,
		IPAddress:    c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "user has been banned and deactivated"})
}

// GET /api/admin/system/health
func (h *AdminHandler) SystemHealth(c *fiber.Ctx) error {
	if adminRole(c) != models.RolePlatformOwner && adminRole(c) != models.RolePlatformAdmin && adminRole(c) != models.LegacyRoleSuper {
		return fiber.NewError(fiber.StatusForbidden, "platform administrator access required")
	}
	report, err := h.svc.CheckSystemHealth(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, report)
}
func (h *AdminHandler) RepairSystemHealth(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var req struct {
		Action             string         `json:"action"`
		MissingAdminEmails map[int]string `json:"missing_admin_emails"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	switch req.Action {
	case "fix_missing_admin_emails":
		if len(req.MissingAdminEmails) == 0 {
			return response.BadRequest(c, "no emails provided for repair")
		}
		if err := h.svc.RepairMissingAdminEmails(c.Context(), adminID, req.MissingAdminEmails); err != nil {
			return response.BadRequest(c, err.Error())
		}
		return response.OK(c, fiber.Map{"message": "Admin emails repaired successfully"})
	default:
		return response.BadRequest(c, "unsupported repair action")
	}
}

// ══════════════════════════════════════════════════════════════════════════════
// Priority 1 Enhancements: Operations, Verification, Trust & Safety, Refunds, Appeals
// ══════════════════════════════════════════════════════════════════════════════

// ── Operations Monitoring (Role 02: Platform Operations Admin) ─────────────────

// GET /api/admin/operations/health
