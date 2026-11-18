package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/pkg/response"
	apperrors "github.com/your-org/used-tools-marketplace/pkg/errors"
)

func (h *AdminHandler) AnalyticsOverview(c *fiber.Ctx) error {
	data, err := h.svc.AnalyticsOverview(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	// Financial totals restricted to Platform Owner and Platform Admin
	if adminRole(c) != models.RolePlatformOwner && adminRole(c) != models.RolePlatformAdmin && adminRole(c) != models.LegacyRoleSuper {
		delete(data, "total_revenue")
		delete(data, "total_volume")
		delete(data, "listing_fee_revenue")
		delete(data, "net_profit")
	}
	return response.OK(c, data)
}

// GET /api/admin/analytics/heatmap
func (h *AdminHandler) AnalyticsHeatmap(c *fiber.Ctx) error {
	data, err := h.svc.AnalyticsWoredaHeatmap(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, data)
}

// GET /api/admin/analytics/monthly?months=6
func (h *AdminHandler) AnalyticsMonthly(c *fiber.Ctx) error {
	months := parseQueryInt(c, "months", 6)
	if months < 1 || months > 24 {
		months = 6
	}
	data, err := h.svc.AnalyticsMonthly(c.Context(), months)
	if err != nil {
		return response.InternalError(c, err)
	}
	// Revenue figures restricted to Platform Owner and Platform Admin
	if adminRole(c) != models.RolePlatformOwner && adminRole(c) != models.RolePlatformAdmin && adminRole(c) != models.LegacyRoleSuper {
		for _, row := range data {
			delete(row, "revenue")
			delete(row, "volume")
		}
	}
	return response.OKList(c, data)
}

// GET /api/admin/analytics/gantt?limit=30
func (h *AdminHandler) AnalyticsGantt(c *fiber.Ctx) error {
	limit := parseQueryInt(c, "limit", 30)
	if limit < 1 || limit > 200 {
		limit = 30
	}
	data, err := h.svc.AnalyticsGantt(c.Context(), limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, data)
}

// GET /api/admin/analytics/fee-revenue  — super admin only
func (h *AdminHandler) AnalyticsFeeRevenue(c *fiber.Ctx) error {
	data, err := h.svc.AnalyticsFeeRevenue(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, data)
}

// ── Settings ──────────────────────────────────────────────────────────────────

// GET /api/admin/settings
func (h *AdminHandler) GetSettings(c *fiber.Ctx) error {
	settings, err := h.svc.GetAllSettings(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	// Financial/fee settings restricted to Platform Owner and Platform Admin
	if adminRole(c) != models.RolePlatformOwner && adminRole(c) != models.RolePlatformAdmin && adminRole(c) != models.LegacyRoleSuper {
		for _, k := range []string{"broker_fee_rate", "referral_commission_rate",
			"normal_startup_fee", "gold_startup_fee",
			"premium_startup_fee", "normal_weekly_increment", "gold_weekly_increment",
			"premium_weekly_increment", "buyer_gold_cost", "buyer_premium_cost"} {
			delete(settings, k)
		}
	}
	return response.OK(c, settings)
}

// PUT /api/admin/settings  — bulk update
// validSettingKeys is the whitelist of keys that can be updated via the API.
// Any key not in this map is rejected to prevent arbitrary data injection.
var validSettingKeys = map[string]bool{
	// Fee configuration
	"broker_fee_rate":          true,
	"referral_commission_rate": true,
	"min_referral_withdrawal":  true,
	"normal_startup_fee":       true,
	"gold_startup_fee":         true,
	"premium_startup_fee":      true,
	"normal_weekly_increment":  true,
	"gold_weekly_increment":    true,
	"premium_weekly_increment": true,
	"listing_free_weeks":       true,
	// Image limits
	"normal_max_images":      true,
	"gold_max_images":        true,
	"premium_max_images":     true,
	"min_images_per_listing": true,
	// Tier thresholds & costs
	"buyer_gold_cost":    true,
	"buyer_premium_cost": true,
	"gold_threshold":     true,
	"premium_threshold":  true,
	// Upgrade pricing
	"buyer_gold_upgrade_price":    true,
	"buyer_premium_upgrade_price": true,
	"gold_upgrade_price":          true,
	"premium_upgrade_price":       true,
	// Payment merchant numbers
	"telebirr_merchant_phone":  true,
	"safaricom_merchant_phone": true,
	"cbe_account_number":       true,
	// Registration & upgrade controls
	"allow_normal_registration":       true,
	"allow_buyer_registration":        true,
	"allow_seller_upgrade_to_gold":    true,
	"allow_seller_upgrade_to_premium": true,
	"allow_buyer_tier_upgrade":        true,
	// OTP & auth
	"otp_required":       true,
	"otp_expiry_minutes": true,
	"otp_gateway":        true,
	// Feature flags
	"maintenance_mode":           true,
	"telebirr_payment_enabled":   true,
	"buyer_tier_upgrade_enabled": true,
	"ai_recommendations_enabled": true,
	// Home page controls
	"hero_banners_enabled":        true,
	"gold_products_enabled":       true,
	"premium_products_enabled":    true,
	"featured_products_enabled":   true,
	"trending_categories_enabled": true,
	// Promotion announcement box
	"promotion_message":    true,
	"promotion_active":     true,
	"promotion_bg_color":   true,
	"promotion_text_color": true,
	"promotion_link_url":   true,
	"promotion_link_label": true,
	// Banner limits
	"max_banners_buyer":  true,
	"max_banners_seller": true,
	// Location limits
	"normal_max_locations":  true,
	"gold_max_locations":    true,
	"premium_max_locations": true,
	// Content limits
	"max_title_length":       true,
	"max_description_length": true,
	// File upload limits
	"max_image_size_mb": true,
	// Currency & payment
	"currency_tolerance": true,
	// Referral withdrawal limits (already declared above on line 389)
	"max_referral_withdrawal": true,
}

func (h *AdminHandler) UpdateSettings(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	// Reject unknown keys before any permission check
	for k := range body {
		if !validSettingKeys[k] {
			return response.BadRequest(c, "unknown setting key: "+k)
		}
	}
	// Financial settings updates restricted to Platform Owner and Platform Admin
	if adminRole(c) != models.RolePlatformOwner && adminRole(c) != models.RolePlatformAdmin && adminRole(c) != models.LegacyRoleSuper {
		restricted := map[string]bool{
			"broker_fee_rate": true, "referral_commission_rate": true,
			"normal_startup_fee": true, "gold_startup_fee": true,
			"premium_startup_fee": true, "normal_weekly_increment": true,
			"gold_weekly_increment": true, "premium_weekly_increment": true,
			"buyer_gold_cost": true, "buyer_premium_cost": true,
			"gold_threshold": true, "premium_threshold": true,
			"min_referral_withdrawal": true,
		}
		for k := range body {
			if restricted[k] {
				return response.ForbiddenCode(c, "insufficient permissions to update "+k, apperrors.ErrPermDenied)
			}
		}
	}
	if err := h.svc.UpdateSettings(c.Context(), body, adminID); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"message": "settings updated"})
}

// ── Notifications ─────────────────────────────────────────────────────────────

// GET /api/admin/notifications
func (h *AdminHandler) ListNotifications(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	role := adminRole(c)
	if role == "" {
		role = "normal"
	}
	notifs, err := h.svc.ListNotifications(c.Context(), adminID, role)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, notifs)
}

// PUT /api/admin/notifications/:id/read
func (h *AdminHandler) MarkNotificationRead(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.MarkNotificationRead(c.Context(), id, adminID); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"message": "marked as read"})
}

// POST /api/admin/notifications  — super admin only
func (h *AdminHandler) CreateNotification(c *fiber.Ctx) error {
	var body struct {
		Title      string `json:"title"`
		Message    string `json:"message"`
		Type       string `json:"type"`
		TargetRole string `json:"target_role"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.Title == "" || body.Message == "" {
		return response.BadRequest(c, "title and message are required")
	}
	if body.Type == "" {
		body.Type = "info"
	}
	if err := h.svc.CreateNotification(c.Context(), body.Title, body.Message, body.Type, body.TargetRole); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"message": "notification created"})
}

// ── Listing Fee Payment Verification ─────────────────────────────────────────

// GET /api/admin/fee-payments/pending
func (h *AdminHandler) ListPendingFeePayments(c *fiber.Ctx) error {
	data, err := h.svc.ListPendingFeePayments(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, data)
}

// PUT /api/admin/fee-payments/:id/verify
func (h *AdminHandler) VerifyListingFeePayment(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	productID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.VerifyListingFeePayment(c.Context(), productID, adminID); err != nil {
		return response.BadRequest(c, err.Error())
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		ProductID: &productID,
		Action:    "verify_fee",
		IPAddress: c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "listing fee payment verified"})
}

// PUT /api/admin/fee-payments/:id/reject
func (h *AdminHandler) RejectListingFeePayment(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	productID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Reason string `json:"reason"`
	}
	_ = c.BodyParser(&body)
	if err := h.svc.RejectListingFeePayment(c.Context(), productID, adminID, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		ProductID: &productID,
		Action:    "reject_fee",
		Notes:     body.Reason,
		IPAddress: c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "listing fee payment claim rejected"})
}

// ── Profit (Super Admin Only) ─────────────────────────────────────────────────

// GET /api/admin/profit/summary  — super admin only
func (h *AdminHandler) ProfitSummary(c *fiber.Ctx) error {
	summary, err := h.svc.ProfitSummary(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, summary)
}

// GET /api/admin/profit/by-tier  — super admin only
func (h *AdminHandler) ProfitByTier(c *fiber.Ctx) error {
	data, err := h.svc.ProfitByTier(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, data)
}

// GET /api/admin/profit/monthly  — super admin only
func (h *AdminHandler) ProfitMonthly(c *fiber.Ctx) error {
	data, err := h.svc.MonthlyProfitTrend(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, data)
}

// GET /api/admin/profit/weekly  — super admin only
func (h *AdminHandler) ProfitWeekly(c *fiber.Ctx) error {
	data, err := h.svc.WeeklyProfitTrend(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, data)
}

// GET /api/admin/profit/annual  — super admin only
func (h *AdminHandler) ProfitAnnual(c *fiber.Ctx) error {
	data, err := h.svc.AnnualProfitTrend(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, data)
}

// GET /api/admin/settings/log?key=broker_fee_rate  — super admin only
func (h *AdminHandler) SettingsChangeLog(c *fiber.Ctx) error {
	key := truncate(c.Query("key"), 100)
	data, err := h.svc.GetSettingsChangeLog(c.Context(), key)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, data)
}

// PUT /api/admin/settings/with-reason  — super admin only, accepts reason field
func (h *AdminHandler) UpdateSettingsWithReason(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	if adminRole(c) != models.RolePlatformOwner && adminRole(c) != models.RolePlatformAdmin && adminRole(c) != models.LegacyRoleSuper {
		return fiber.NewError(fiber.StatusForbidden, "platform administrator access required")
	}
	var body struct {
		Settings map[string]string `json:"settings"`
		Reason   string            `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	for k := range body.Settings {
		if !validSettingKeys[k] {
			return response.BadRequest(c, "unknown setting key: "+k)
		}
	}
	if err := h.svc.UpdateSettingsWithReason(c.Context(), body.Settings, body.Reason, adminID); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"message": "settings updated"})
}

// ── Category Management (Super Admin Only) ────────────────────────────────────

// GET /api/admin/categories
func (h *AdminHandler) GetPublicPromotion(c *fiber.Ctx) error {
	all, err := h.svc.GetAllSettings(c.Context())
	if err != nil {
		return response.OK(c, fiber.Map{"active": false})
	}
	active := all["promotion_active"] == "true"
	if !active {
		return response.OK(c, fiber.Map{"active": false})
	}
	return response.OK(c, fiber.Map{
		"active":     true,
		"message":    all["promotion_message"],
		"bg_color":   all["promotion_bg_color"],
		"text_color": all["promotion_text_color"],
		"link_url":   all["promotion_link_url"],
		"link_label": all["promotion_link_label"],
	})
}

// GET /api/public/payment-settings — returns payment phone numbers for the pay page (no auth required)
func (h *AdminHandler) GetPublicPaymentSettings(c *fiber.Ctx) error {
	all, err := h.svc.GetAllSettings(c.Context())
	if err != nil {
		// Settings unavailable — return empty map, never hardcoded values
		return response.OK(c, fiber.Map{})
	}
	// Only expose the specific keys needed by the public pay page and fee calculations
	// Never expose sensitive admin-only settings
	publicKeys := []string{
		// Payment methods
		"telebirr_merchant_phone", "safaricom_merchant_phone", "cbe_account_number",
		// Tier upgrade costs
		"gold_upgrade_price", "premium_upgrade_price",
		"buyer_gold_upgrade_price", "buyer_premium_upgrade_price",
		// Tier thresholds
		"gold_threshold", "premium_threshold",
		// Fee rates (public info - users need to see fees before listing)
		"referral_commission_rate", "broker_fee_rate",
		"normal_startup_fee", "gold_startup_fee", "premium_startup_fee",
		"normal_weekly_increment", "gold_weekly_increment", "premium_weekly_increment",
		"listing_free_weeks",
		// Image limits per tier
		"normal_max_images", "gold_max_images", "premium_max_images",
		"min_images_per_listing", "max_image_size_mb",
		// Content limits
		"max_title_length", "max_description_length",
		// Location limits
		"normal_max_locations", "gold_max_locations", "premium_max_locations",
		// Referral system
		"min_referral_withdrawal", "max_referral_withdrawal",
		// Company contact
		"company_phone", "company_email", "company_location",
		// Misc
		"otp_gateway", "currency_tolerance",
	}
	out := fiber.Map{}
	for _, k := range publicKeys {
		if v, ok := all[k]; ok {
			out[k] = v
		}
	}
	return response.OK(c, out)
}

// GET /api/public/settings — alias for public platform settings and company contact info
func (h *AdminHandler) GetPublicSettings(c *fiber.Ctx) error {
	return h.GetPublicPaymentSettings(c)
}

// ── Buyer Upgrade Requests ────────────────────────────────────────────────────

// GET /api/admin/buyer-upgrade-requests?status=pending&page=1&limit=20
