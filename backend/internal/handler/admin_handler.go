package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/service"
	apperrors "github.com/your-org/used-tools-marketplace/pkg/errors"
	"github.com/your-org/used-tools-marketplace/pkg/response"
	"github.com/your-org/used-tools-marketplace/pkg/validations"
)

type AdminHandler struct {
	svc      *service.AdminService
	authSvc  *service.AuthService
	catRepo  repository.CategoryRepository
	prodRepo repository.ProductRepository
}

func NewAdminHandler(svc *service.AdminService, authSvc ...*service.AuthService) *AdminHandler {
	h := &AdminHandler{svc: svc}
	if len(authSvc) > 0 {
		h.authSvc = authSvc[0]
	}
	return h
}

func (h *AdminHandler) WithCategoryRepo(r repository.CategoryRepository) *AdminHandler {
	h.catRepo = r
	return h
}

func (h *AdminHandler) WithProductRepo(r repository.ProductRepository) *AdminHandler {
	h.prodRepo = r
	return h
}

func adminRole(c *fiber.Ctx) string {
	r, _ := c.Locals("adminRole").(string)
	return r
}

// PUT /api/admin/password
func (h *AdminHandler) ChangePassword(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.CurrentPassword == "" || body.NewPassword == "" {
		return response.BadRequest(c, "current_password and new_password are required")
	}
	if valid, errKey := validations.ValidatePassword(body.NewPassword); !valid {
		if errKey == "password_min" {
			return response.BadRequest(c, "new password must be at least 8 characters")
		}
		return response.BadRequest(c, "invalid password")
	}
	if h.authSvc == nil {
		return response.InternalError(c, fmt.Errorf("password change not configured"))
	}
	if err := h.authSvc.ChangePassword(c.Context(), adminID, body.CurrentPassword, body.NewPassword); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "password updated"})
}

// GET /api/admin/products/pending
func (h *AdminHandler) ListPending(c *fiber.Ctx) error {
	products, err := h.svc.ListPending(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, products)
}

// GET /api/admin/products
func (h *AdminHandler) ListAllProducts(c *fiber.Ctx) error {
	page, limit := parsePagination(c, 50, 200)
	var products []models.Product
	var total int
	var err error
	if status := c.Query("status"); status != "" {
		validStatuses := map[string]bool{
			"pending": true, "approved": true, "rejected": true,
			"sold": true, "expired": true,
		}
		if !validStatuses[status] {
			return response.BadRequest(c, "invalid status — must be pending, approved, rejected, sold, or expired")
		}
		products, total, err = h.svc.ListProductsByStatus(c.Context(), status, page, limit)
	} else {
		products, total, err = h.svc.ListAllProducts(c.Context(), page, limit)
	}
	if err != nil {
		return response.InternalError(c, err)
	}
	if products == nil {
		products = []models.Product{}
	}
	return response.OK(c, fiber.Map{"data": products, "total": total, "page": page, "limit": limit})
}

// PUT /api/admin/products/:id/review
func (h *AdminHandler) ReviewProduct(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	productID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if err := h.svc.ReviewProduct(c.Context(), productID, adminID, body.Status, body.Notes); err != nil {
		return response.BadRequest(c, err.Error())
	}
	action := "approve"
	if body.Status == "rejected" {
		action = "reject"
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		ProductID: &productID,
		Action:    action,
		Notes:     body.Notes,
		IPAddress: c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "product reviewed"})
}

// DELETE /api/admin/products/:id
func (h *AdminHandler) DeleteProduct(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.DeleteProduct(c.Context(), id); err != nil {
		return response.InternalError(c, err)
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		ProductID: &id,
		Action:    "delete",
		IPAddress: c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "product deleted"})
}

// GET /api/admin/sellers
func (h *AdminHandler) ListSellers(c *fiber.Ctx) error {
	page, limit := parsePagination(c, 50, 200)
	search := strings.TrimSpace(c.Query("search"))
	status := strings.TrimSpace(c.Query("status"))
	tier := strings.TrimSpace(c.Query("tier"))
	sellers, total, err := h.svc.ListSellersFiltered(c.Context(), page, limit, search, status, tier)
	if err != nil {
		return response.InternalError(c, err)
	}
	if sellers == nil {
		sellers = []models.User{}
	}
	// Hide sensitive profit data for non-platform owner
	if adminRole(c) != models.RolePlatformOwner && adminRole(c) != models.LegacyRoleSuper {
		for i := range sellers {
			sellers[i].NetProfit = 0
		}
	}
	return response.OK(c, fiber.Map{"data": sellers, "total": total, "page": page, "limit": limit})
}

// PUT /api/admin/sellers/:id/active
func (h *AdminHandler) SetSellerActive(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Active bool   `json:"active"`
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	// Reason is required when deactivating
	if !body.Active && strings.TrimSpace(body.Reason) == "" {
		return response.BadRequest(c, "reason is required when deactivating a seller")
	}
	if err := h.svc.SetSellerActive(c.Context(), id, body.Active, strings.TrimSpace(body.Reason)); err != nil {
		return response.BadRequest(c, err.Error())
	}
	action := "activate"
	if !body.Active {
		action = "deactivate"
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:      adminID,
		TargetUserID: &id,
		Action:       action,
		Notes:        strings.TrimSpace(body.Reason),
		IPAddress:    c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "seller updated"})
}

// GET /api/admin/transactions  — permission-gated
func (h *AdminHandler) ListTransactions(c *fiber.Ctx) error {
	page, limit := parsePagination(c, 50, 200)
	search := strings.TrimSpace(c.Query("search"))
	status := strings.TrimSpace(c.Query("status"))
	method := strings.TrimSpace(c.Query("method"))
	txType := strings.TrimSpace(c.Query("type"))
	fromDate := strings.TrimSpace(c.Query("from_date"))
	toDate := strings.TrimSpace(c.Query("to_date"))

	txns, total, err := h.svc.ListAllTransactionsFiltered(c.Context(), page, limit, search, status, method, txType, fromDate, toDate)
	if err != nil {
		return response.InternalError(c, err)
	}
	if txns == nil {
		txns = []models.Transaction{}
	}
	return response.OK(c, fiber.Map{"data": txns, "total": total, "page": page, "limit": limit})
}

// GET /api/admin/locations
func (h *AdminHandler) ListLocations(c *fiber.Ctx) error {
	locs, err := h.svc.ListLocations(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, locs)
}

// POST /api/admin/locations
func (h *AdminHandler) AddLocation(c *fiber.Ctx) error {
	var l models.Location
	if err := c.BodyParser(&l); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	id, err := h.svc.AddLocation(c.Context(), l)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	l.ID = id
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": l})
}

// PUT /api/admin/locations/:id
func (h *AdminHandler) UpdateLocation(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var l models.Location
	if err := c.BodyParser(&l); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	l.ID = id
	if err := h.svc.UpdateLocation(c.Context(), l); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "location updated"})
}

// DELETE /api/admin/locations/:id
func (h *AdminHandler) DeleteLocation(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.DeleteLocation(c.Context(), id); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"message": "location deleted"})
}

// ── Analytics ─────────────────────────────────────────────────────────────────

// GET /api/admin/analytics/overview
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
	"broker_fee_rate":           true,
	"referral_commission_rate":  true,
	"min_referral_withdrawal":   true,
	"normal_startup_fee":        true,
	"gold_startup_fee":          true,
	"premium_startup_fee":       true,
	"normal_weekly_increment":   true,
	"gold_weekly_increment":     true,
	"premium_weekly_increment":  true,
	"listing_free_weeks":        true,
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
func (h *AdminHandler) ListCategories(c *fiber.Ctx) error {
	if h.catRepo == nil {
		return response.InternalError(c, fiber.NewError(500, "category repo not configured"))
	}
	cats, err := h.catRepo.ListAll(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, cats)
}

// POST /api/admin/categories
func (h *AdminHandler) CreateCategory(c *fiber.Ctx) error {
	if h.catRepo == nil {
		return response.InternalError(c, fiber.NewError(500, "category repo not configured"))
	}
	var cat models.Category
	if err := c.BodyParser(&cat); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if cat.Slug == "" || cat.NameEn == "" {
		return response.BadRequest(c, "slug and name_en are required")
	}
	id, err := h.catRepo.Create(c.Context(), cat)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	cat.ID = id
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": cat})
}

// PUT /api/admin/categories/:id
func (h *AdminHandler) UpdateCategory(c *fiber.Ctx) error {
	if h.catRepo == nil {
		return response.InternalError(c, fiber.NewError(500, "category repo not configured"))
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var cat models.Category
	if err := c.BodyParser(&cat); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	cat.ID = id
	if err := h.catRepo.Update(c.Context(), cat); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "category updated"})
}

// DELETE /api/admin/categories/:id
func (h *AdminHandler) DeleteCategory(c *fiber.Ctx) error {
	if h.catRepo == nil {
		return response.InternalError(c, fiber.NewError(500, "category repo not configured"))
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.catRepo.Delete(c.Context(), id); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"message": "category deleted"})
}

// GET /api/admin/me — returns the current admin's profile
// Uses authSvc.GetProfile (single GetByID query) — avoids loading all admins.
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

	validRoles := map[string]bool{
		models.RolePlatformAdmin:   true,
		models.RoleSecurityAdmin:   true,
		models.RoleOperationsAdmin: true,
		models.RoleModerator:       true,
		models.RoleSupportAgent:    true,
		models.LegacyRoleNormal:    true,
	}
	if !validRoles[role] {
		return response.BadRequest(c, "invalid admin_role: must be platform_admin, security_admin, operations_admin, moderator, or support_agent")
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
		Action            string         `json:"action"`
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
func (h *AdminHandler) GetOperationsHealth(c *fiber.Ctx) error {
	health, err := h.svc.GetOperationsHealth(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, health)
}

// GET /api/admin/operations/jobs
func (h *AdminHandler) GetJobStatusSummary(c *fiber.Ctx) error {
	jobs, err := h.svc.ListJobStatusSummary(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, jobs)
}

// GET /api/admin/operations/failed-jobs
func (h *AdminHandler) GetFailedJobs(c *fiber.Ctx) error {
	limit := parseQueryInt(c, "limit", 50)
	jobs, err := h.svc.ListRecentFailedJobs(c.Context(), limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, jobs)
}

// ── Seller Verification (Role 08: Seller Verification Admin) ──────────────────

// GET /api/admin/verification/documents
func (h *AdminHandler) ListVerificationDocuments(c *fiber.Ctx) error {
	status := c.Query("status", "pending")
	page, limit := parsePagination(c, 50, 100)
	
	docs, total, err := h.svc.ListVerificationDocuments(c.Context(), status, page, limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"data":    docs,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// PUT /api/admin/verification/documents/:id/review
func (h *AdminHandler) ReviewVerificationDocument(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	
	docID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	
	var body struct {
		Approved bool   `json:"approved"`
		Reason   string `json:"reason"`
		Notes    string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	
	if err := h.svc.ReviewVerificationDocument(c.Context(), docID, adminID, body.Approved, body.Reason, body.Notes); err != nil {
		return response.InternalError(c, err)
	}
	
	// Log the action
	action := "approve_verification"
	if !body.Approved {
		action = "reject_verification"
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    action,
		Notes:     fmt.Sprintf("Document ID: %d, Reason: %s", docID, body.Reason),
		IPAddress: c.IP(),
	})
	
	return response.OK(c, fiber.Map{
		"message": "document reviewed successfully",
	})
}

// ── Trust & Safety (Role 21: Fraud & Trust & Safety Admin) ────────────────────

// GET /api/admin/trust/reports
func (h *AdminHandler) ListAbuseReports(c *fiber.Ctx) error {
	filter := models.AbuseReportFilter{
		Status:         c.Query("status", "pending"),
		Priority:       c.Query("priority", "all"),
		Reason:         c.Query("reason", "all"),
		ReportedUserID: parseQueryInt(c, "reported_user_id", 0),
	}
	filter.Page, filter.Limit = parsePagination(c, 50, 100)
	
	reports, total, err := h.svc.ListAbuseReports(c.Context(), filter)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"data":    reports,
		"meta": fiber.Map{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

// PUT /api/admin/trust/reports/:id/review
func (h *AdminHandler) ReviewAbuseReport(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	
	reportID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	
	var body struct {
		Status      string `json:"status"` // 'investigating', 'resolved', 'dismissed'
		ActionTaken string `json:"action_taken"`
		Action      string `json:"action"`
		Notes       string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.ActionTaken == "" && body.Action != "" {
		body.ActionTaken = body.Action
	}
	
	if body.Status == "" {
		return response.BadRequest(c, "status is required")
	}
	
	if err := h.svc.ReviewAbuseReport(c.Context(), reportID, adminID, body.Status, body.ActionTaken, body.Notes); err != nil {
		return response.BadRequest(c, err.Error())
	}
	
	// Log the action
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    "review_abuse_report",
		Notes:     fmt.Sprintf("Report ID: %d, Status: %s, Action: %s", reportID, body.Status, body.ActionTaken),
		IPAddress: c.IP(),
	})
	
	return response.OK(c, fiber.Map{
		"message": "abuse report reviewed successfully",
	})
}

// PUT /api/admin/trust/reports/:id/priority
func (h *AdminHandler) UpdateAbuseReportPriority(c *fiber.Ctx) error {
	reportID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	
	var body struct {
		Priority string `json:"priority"` // 'low', 'normal', 'high', 'critical'
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	
	if err := h.svc.UpdateAbuseReportPriority(c.Context(), reportID, body.Priority); err != nil {
		return response.BadRequest(c, err.Error())
	}
	
	return response.OK(c, fiber.Map{
		"message": "priority updated successfully",
	})
}

// ── Refund Management (Role 20: Refund & Dispute Admin) ───────────────────────

// GET /api/admin/refunds
func (h *AdminHandler) ListRefundRequests(c *fiber.Ctx) error {
	filter := models.RefundRequestFilter{
		Status:        c.Query("status", "pending"),
		RequesterRole: c.Query("requester_role", "all"),
	}
	filter.Page, filter.Limit = parsePagination(c, 50, 100)
	
	refunds, total, err := h.svc.ListRefundRequests(c.Context(), filter)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"data":    refunds,
		"meta": fiber.Map{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

// PUT /api/admin/refunds/:id/review
func (h *AdminHandler) ReviewRefundRequest(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	
	refundID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	
	var body struct {
		Approved bool   `json:"approved"`
		Reason   string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	
	if err := h.svc.ReviewRefundRequest(c.Context(), refundID, adminID, body.Approved, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}
	
	// Log the action
	action := "approve_refund"
	if !body.Approved {
		action = "reject_refund"
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    action,
		Notes:     fmt.Sprintf("Refund ID: %d, Reason: %s", refundID, body.Reason),
		IPAddress: c.IP(),
	})
	
	return response.OK(c, fiber.Map{
		"message": "refund request reviewed successfully",
	})
}

// PUT /api/admin/refunds/:id/processed
func (h *AdminHandler) MarkRefundProcessed(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	
	refundID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	
	var body struct {
		Reference string `json:"reference"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	
	if body.Reference == "" {
		return response.BadRequest(c, "refund reference is required")
	}
	
	if err := h.svc.MarkRefundProcessed(c.Context(), refundID, body.Reference); err != nil {
		return response.BadRequest(c, err.Error())
	}
	
	// Log the action
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    "process_refund",
		Notes:     fmt.Sprintf("Refund ID: %d, Reference: %s", refundID, body.Reference),
		IPAddress: c.IP(),
	})
	
	return response.OK(c, fiber.Map{
		"message": "refund marked as processed",
	})
}

// ── Account Appeals (Role 30: Appeals & Account Recovery Admin) ───────────────

// GET /api/admin/appeals
func (h *AdminHandler) ListAccountAppeals(c *fiber.Ctx) error {
	filter := models.AccountAppealFilter{
		Status:     c.Query("status", "pending"),
		AppealType: c.Query("appeal_type", "all"),
		UserID:     parseQueryInt(c, "user_id", 0),
	}
	filter.Page, filter.Limit = parsePagination(c, 50, 100)
	
	appeals, total, err := h.svc.ListAccountAppeals(c.Context(), filter)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"data":    appeals,
		"meta": fiber.Map{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

// PUT /api/admin/appeals/:id/review
func (h *AdminHandler) ReviewAccountAppeal(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	
	appealID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	
	var body struct {
		Approved    bool   `json:"approved"`
		Decision    string `json:"decision"`
		ActionTaken string `json:"action_taken"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	
	if body.Decision == "" {
		return response.BadRequest(c, "decision is required")
	}
	
	if err := h.svc.ReviewAccountAppeal(c.Context(), appealID, adminID, body.Approved, body.Decision, body.ActionTaken); err != nil {
		return response.BadRequest(c, err.Error())
	}
	
	// Log the action
	action := "approve_appeal"
	if !body.Approved {
		action = "reject_appeal"
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    action,
		Notes:     fmt.Sprintf("Appeal ID: %d, Decision: %s, Action: %s", appealID, body.Decision, body.ActionTaken),
		IPAddress: c.IP(),
	})
	
	return response.OK(c, fiber.Map{
		"message": "appeal reviewed successfully",
	})
}

// PUT /api/admin/appeals/:id/status
func (h *AdminHandler) UpdateAppealStatus(c *fiber.Ctx) error {
	appealID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	
	var body struct {
		Status string `json:"status"` // 'pending', 'under_review', 'approved', 'rejected'
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	
	if err := h.svc.UpdateAppealStatus(c.Context(), appealID, body.Status); err != nil {
		return response.BadRequest(c, err.Error())
	}
	
	return response.OK(c, fiber.Map{
		"message": "appeal status updated successfully",
	})
}
