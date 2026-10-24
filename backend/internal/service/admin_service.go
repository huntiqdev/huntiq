package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/pkg/validations"
)

type AdminService struct {
	repo           *repository.AdminRepo
	referralSvc    *ReferralService
	sellerSvc      *SellerService
	sessionRevoker func(ctx context.Context, userID int, reason string) error
}

func NewAdminService(repo repository.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

// WithReferralService injects the referral service for commission recording on upgrade approvals.
func (s *AdminService) WithReferralService(r *ReferralService) *AdminService {
	s.referralSvc = r
	return s
}

// WithSellerService injects the seller service for sending seller notifications.
func (s *AdminService) WithSellerService(ss *SellerService) *AdminService {
	s.sellerSvc = ss
	return s
}

// WithSessionRevoker injects a callback to revoke server-side sessions when roles or access change.
func (s *AdminService) WithSessionRevoker(revoker func(ctx context.Context, userID int, reason string) error) *AdminService {
	s.sessionRevoker = revoker
	return s
}

func (s *AdminService) ListPending(ctx context.Context) ([]models.Product, error) {
	return s.repo.ListPending(ctx)
}

func (s *AdminService) ListAllProducts(ctx context.Context, page, limit int) ([]models.Product, int, error) {
	return s.repo.ListAllProductsPaged(ctx, page, limit)
}

func (s *AdminService) ListProductsByStatus(ctx context.Context, status string, page, limit int) ([]models.Product, int, error) {
	return s.repo.ListProductsByStatusPaged(ctx, status, page, limit)
}

func (s *AdminService) ReviewProduct(ctx context.Context, productID, adminID int, status, notes string) error {
	if status != "approved" && status != "rejected" {
		return errors.New("status must be approved or rejected")
	}
	if err := s.repo.ReviewProduct(ctx, productID, adminID, status, notes); err != nil {
		return err
	}
	// Notify the seller — fire-and-forget, never blocks the admin action.
	// Use context.Background() — the request ctx is cancelled once the response is sent.
	if s.sellerSvc != nil {
		safeGo(func() {
			sellerID, err := s.repo.GetProductSellerID(context.Background(), productID)
			if err != nil || sellerID == 0 {
				return
			}
			var title, msg, ntype string
			if status == "approved" {
				title = "Product Approved"
				msg = "Your product listing has been approved and is now live on the marketplace."
				ntype = "success"
			} else {
				title = "Product Rejected"
				msg = "Your product listing was rejected."
				if notes != "" {
					msg += " Admin note: " + notes
				}
				ntype = "warning"
			}
			_ = s.sellerSvc.CreateNotification(context.Background(), sellerID, title, msg, ntype)
		})
	}
	return nil
}

func (s *AdminService) DeleteProduct(ctx context.Context, productID int) error {
	return s.repo.DeleteProduct(ctx, productID)
}

func (s *AdminService) ListSellers(ctx context.Context, page, limit int) ([]models.User, int, error) {
	return s.repo.ListSellersPaged(ctx, page, limit)
}

func (s *AdminService) ListSellersFiltered(ctx context.Context, page, limit int, search, status, tier string) ([]models.User, int, error) {
	return s.repo.ListSellersFilteredPaged(ctx, page, limit, search, status, tier)
}

func (s *AdminService) SetSellerActive(ctx context.Context, sellerID int, active bool, reason string) error {
	if err := s.repo.SetSellerActive(ctx, sellerID, active, reason); err != nil {
		return err
	}
	if !active && s.sessionRevoker != nil {
		_ = s.sessionRevoker(ctx, sellerID, fmt.Sprintf("seller_deactivated: %s", reason))
	}
	return nil
}

func (s *AdminService) ListAllTransactions(ctx context.Context, page, limit int) ([]models.Transaction, int, error) {
	return s.repo.ListAllTransactionsPaged(ctx, page, limit)
}

func (s *AdminService) ListAllTransactionsFiltered(ctx context.Context, page, limit int, search, status, paymentMethod, txType, dateFrom, dateTo string) ([]models.Transaction, int, error) {
	return s.repo.ListAllTransactionsFilteredPaged(ctx, page, limit, search, status, paymentMethod, txType, dateFrom, dateTo)
}

func (s *AdminService) AddLocation(ctx context.Context, l models.Location) (int, error) {
	if l.Region == "" {
		return 0, errors.New("region is required")
	}
	return s.repo.AddLocation(ctx, l)
}

func (s *AdminService) UpdateLocation(ctx context.Context, l models.Location) error {
	if l.ID == 0 || l.Region == "" {
		return errors.New("id and region are required")
	}
	return s.repo.UpdateLocation(ctx, l)
}

func (s *AdminService) DeleteLocation(ctx context.Context, id int) error {
	return s.repo.DeleteLocation(ctx, id)
}

func (s *AdminService) ListLocations(ctx context.Context) ([]models.Location, error) {
	return s.repo.ListLocations(ctx)
}

// ── Analytics ─────────────────────────────────────────────────────────────────

func (s *AdminService) AnalyticsOverview(ctx context.Context) (map[string]interface{}, error) {
	return s.repo.AnalyticsOverview(ctx)
}

func (s *AdminService) AnalyticsWoredaHeatmap(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.AnalyticsWoredaHeatmap(ctx)
}

func (s *AdminService) AnalyticsMonthly(ctx context.Context, months int) ([]map[string]interface{}, error) {
	if months < 1 || months > 24 {
		months = 6
	}
	return s.repo.AnalyticsMonthly(ctx, months)
}

func (s *AdminService) AnalyticsGantt(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	if limit < 1 || limit > 100 {
		limit = 30
	}
	return s.repo.AnalyticsGantt(ctx, limit)
}

func (s *AdminService) AnalyticsFeeRevenue(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.AnalyticsFeeRevenue(ctx)
}

// ── Settings ──────────────────────────────────────────────────────────────────

func (s *AdminService) GetAllSettings(ctx context.Context) (map[string]string, error) {
	return s.repo.GetAllSettings(ctx)
}

// GetSetting returns a single platform setting value, or "" if not found.
// Uses a targeted single-row query instead of fetching all settings.
func (s *AdminService) GetSetting(ctx context.Context, key string) string {
	val, _ := s.repo.GetSetting(ctx, key)
	return val
}

func (s *AdminService) UpdateSetting(ctx context.Context, key, value string, adminID int) error {
	if key == "" {
		return errors.New("key is required")
	}
	return s.repo.UpdateSetting(ctx, key, value, adminID)
}

// settingNumericKeys is the set of platform_settings keys whose values must be non-negative numbers.
var settingNumericKeys = map[string]bool{
	"broker_fee_rate": true, "normal_startup_fee": true, "gold_startup_fee": true,
	"premium_startup_fee": true, "normal_weekly_increment": true,
	"gold_weekly_increment": true, "premium_weekly_increment": true,
	"buyer_gold_cost": true, "buyer_premium_cost": true,
	"gold_threshold": true, "premium_threshold": true,
	"buyer_gold_upgrade_price": true, "buyer_premium_upgrade_price": true,
	"gold_upgrade_price": true, "premium_upgrade_price": true,
	"listing_free_weeks": true, "otp_expiry_minutes": true,
	"normal_max_images": true, "gold_max_images": true, "premium_max_images": true,
	"min_images_per_listing": true,
	// Content limits
	"max_title_length": true, "max_description_length": true,
	// Location limits
	"normal_max_locations": true, "gold_max_locations": true, "premium_max_locations": true,
	// Referral
	"referral_commission_rate": true, "min_referral_withdrawal": true, "max_referral_withdrawal": true,
	// Banner limits
	"max_banners_buyer": true, "max_banners_seller": true,
	// File upload limits
	"max_image_size_mb": true,
	// Currency & payment
	"currency_tolerance": true,
}

// settingBooleanKeys is the set of platform_settings keys whose values must be "true" or "false".
var settingBooleanKeys = map[string]bool{
	"otp_required": true, "maintenance_mode": true,
	"allow_normal_registration": true, "allow_buyer_registration": true,
	"allow_seller_upgrade_to_gold": true, "allow_seller_upgrade_to_premium": true,
	"allow_buyer_tier_upgrade": true, "telebirr_payment_enabled": true,
	"buyer_tier_upgrade_enabled": true, "ai_recommendations_enabled": true,
	"hero_banners_enabled": true, "gold_products_enabled": true,
	"premium_products_enabled": true, "featured_products_enabled": true,
	"trending_categories_enabled": true,
	// Promotion
	"promotion_active": true,
}

// validateSettings checks all setting values against type rules.
// Returns an error on the first invalid value, nil if all are valid.
func validateSettings(settings map[string]string) error {
	for k, v := range settings {
		if settingNumericKeys[k] {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil || f < 0 {
				return fmt.Errorf("setting %q must be a non-negative number, got %q", k, v)
			}
		}
		if settingBooleanKeys[k] && v != "true" && v != "false" {
			return fmt.Errorf("setting %q must be \"true\" or \"false\", got %q", k, v)
		}
	}
	return nil
}

func (s *AdminService) UpdateSettings(ctx context.Context, settings map[string]string, adminID int) error {
	if err := validateSettings(settings); err != nil {
		return err
	}
	for k, v := range settings {
		if err := s.repo.UpdateSetting(ctx, k, v, adminID); err != nil {
			return err
		}
	}
	return nil
}

// ── Notifications ─────────────────────────────────────────────────────────────

func (s *AdminService) ListNotifications(ctx context.Context, adminID int, adminRole string) ([]map[string]interface{}, error) {
	return s.repo.ListAdminNotifications(ctx, adminID, adminRole)
}

func (s *AdminService) MarkNotificationRead(ctx context.Context, notifID, adminID int) error {
	return s.repo.MarkNotificationRead(ctx, notifID, adminID)
}

func (s *AdminService) CreateNotification(ctx context.Context, title, message, ntype, targetRole string) error {
	return s.repo.CreateAdminNotification(ctx, title, message, ntype, targetRole)
}

// NotifyFirstLogin fires a super-admin notification when a sub-admin or seller
// logs in for the very first time. Fire-and-forget — never blocks the login response.
func (s *AdminService) NotifyFirstLogin(ctx context.Context, name, role string) {
	var title, message string
	switch role {
	case "admin":
		title = "Sub-Admin First Login"
		message = fmt.Sprintf("Sub-admin \"%s\" has logged in for the first time.", name) //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
	case "seller":
		title = "Seller First Login"
		message = fmt.Sprintf("Seller \"%s\" has logged in for the first time.", name) //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
	default:
		return // buyers and super admins don't trigger this
	}
	// Target only super admins
	_ = s.repo.CreateAdminNotification(ctx, title, message, "info", "super")
}

// ── Listing Fee Payment Verification ─────────────────────────────────────────

func (s *AdminService) ListPendingFeePayments(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.ListPendingFeePayments(ctx)
}

func (s *AdminService) VerifyListingFeePayment(ctx context.Context, productID, adminID int) error {
	return s.repo.VerifyListingFeePayment(ctx, productID, adminID)
}

func (s *AdminService) RejectListingFeePayment(ctx context.Context, productID, adminID int, reason string) error {
	return s.repo.RejectListingFeePayment(ctx, productID, adminID, reason)
}

// ── Profit (Super Admin Only) ─────────────────────────────────────────────────

func (s *AdminService) ProfitSummary(ctx context.Context) (map[string]interface{}, error) {	return s.repo.ProfitSummary(ctx)
}

func (s *AdminService) ProfitByTier(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.ProfitByTier(ctx)
}

func (s *AdminService) MonthlyProfitTrend(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.MonthlyProfitTrend(ctx)
}

func (s *AdminService) WeeklyProfitTrend(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.WeeklyProfitTrend(ctx)
}

func (s *AdminService) AnnualProfitTrend(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.AnnualProfitTrend(ctx)
}

func (s *AdminService) GetSettingsChangeLog(ctx context.Context, key string) ([]map[string]interface{}, error) {
	return s.repo.GetSettingsChangeLog(ctx, key, 100)
}

func (s *AdminService) UpdateSettingsWithReason(ctx context.Context, settings map[string]string, reason string, adminID int) error {
	if err := validateSettings(settings); err != nil {
		return err
	}
	for k, v := range settings {
		if err := s.repo.UpdateSettingWithReason(ctx, k, v, reason, adminID); err != nil {
			return err
		}
	}
	return nil
}

// ── Sub-Admin Management (Super Admin Only) ───────────────────────────────────

func (s *AdminService) ListAdmins(ctx context.Context) ([]models.User, error) {
	return s.repo.ListAdmins(ctx)
}

func (s *AdminService) CreateAdmin(ctx context.Context, name, phone, email, password, adminRole string) (int, error) {
	if name == "" || phone == "" || password == "" {
		return 0, errors.New("name, phone, and password are required")
	}
	if valid, errKey := validations.ValidatePassword(password); !valid {
		if errKey == "password_min" {
			return 0, errors.New("password must be at least 8 characters")
		}
		return 0, errors.New("invalid password")
	}
	if adminRole == "" {
		adminRole = models.RoleOperationsAdmin
	}
	// Platform Owner accounts cannot be provisioned via web API; they require direct bootstrap
	if adminRole == models.RolePlatformOwner || adminRole == models.LegacyRoleSuper {
		return 0, errors.New("Platform Owner accounts cannot be provisioned via the web API")
	}

	validRoles := map[string]bool{
		models.RolePlatformAdmin:   true,
		models.RoleSecurityAdmin:   true,
		models.RoleOperationsAdmin: true,
		models.RoleModerator:       true,
		models.RoleSupportAgent:    true,
		models.LegacyRoleNormal:    true,
	}
	if !validRoles[adminRole] {
		return 0, errors.New("invalid admin_role: must be platform_admin, security_admin, operations_admin, moderator, or support_agent")
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}
	id, err := s.repo.CreateAdmin(ctx, name, phone, email, hash, adminRole)
	if err != nil {
		return 0, err
	}
	// Grant baseline permissions based on domain role, or default permissions for legacy normal
	var defaultPerms []string
	if domainPerms, ok := models.RoleDomainPermissions[adminRole]; ok {
		for p := range domainPerms {
			defaultPerms = append(defaultPerms, p)
		}
	} else {
		defaultPerms = models.DefaultAdminPermissions
	}
	_ = s.repo.SetAdminPermissions(ctx, id, defaultPerms)
	return id, nil
}

func (s *AdminService) UpdateAdminRole(ctx context.Context, adminID int, adminRole string) error {
	validRoles := map[string]bool{
		models.RolePlatformOwner:   true,
		models.RolePlatformAdmin:   true,
		models.RoleSecurityAdmin:   true,
		models.RoleOperationsAdmin: true,
		models.RoleModerator:       true,
		models.RoleSupportAgent:    true,
		models.LegacyRoleSuper:     true,
		models.LegacyRoleNormal:    true,
	}
	if !validRoles[adminRole] {
		return errors.New("admin_role must be a valid enterprise role: platform_owner, platform_admin, security_admin, operations_admin, moderator, or support_agent")
	}
	// Prevent demoting the last platform owner / super admin — always keep at least one.
	if adminRole != models.RolePlatformOwner && adminRole != models.LegacyRoleSuper {
		isSuper, _ := s.repo.IsAdminSuper(ctx, adminID)
		if isSuper {
			count, err := s.repo.CountSuperAdmins(ctx)
			if err != nil {
				return err
			}
			if count <= 1 {
				return errors.New("cannot demote the last Platform Owner")
			}
		}
	}
	if err := s.repo.UpdateAdminRole(ctx, adminID, adminRole); err != nil {
		return err
	}

	// Reset domain baseline permissions for the newly assigned role
	var newPerms []string
	if perms, ok := models.RoleDomainPermissions[adminRole]; ok {
		for p := range perms {
			newPerms = append(newPerms, p)
		}
	}
	_ = s.repo.SetAdminPermissions(ctx, adminID, newPerms)

	// Immediately invalidate all active server-side sessions so stale permissions cannot be used
	if s.sessionRevoker != nil {
		_ = s.sessionRevoker(ctx, adminID, "admin_role_updated")
	}

	return nil
}

func (s *AdminService) DeleteAdmin(ctx context.Context, adminID int) error {
	// Prevent deleting a super admin — super admins can only be removed via direct DB access.
	isSuper, err := s.repo.IsAdminSuper(ctx, adminID)
	if err != nil {
		return err
	}
	if isSuper {
		return errors.New("super admins cannot be deleted via the API")
	}
	if err := s.repo.DeleteAdmin(ctx, adminID); err != nil {
		return err
	}
	if s.sessionRevoker != nil {
		_ = s.sessionRevoker(ctx, adminID, "admin_deleted")
	}
	return nil
}

// ── Granular Permissions ──────────────────────────────────────────────────────

func (s *AdminService) GetAdminPermissions(ctx context.Context, adminID int) ([]string, error) {
	// Super admins have all permissions — querying their permission rows is meaningless.
	isSuper, err := s.repo.IsAdminSuper(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if isSuper {
		return nil, errors.New("super admins have full access — permissions are not applicable")
	}
	perms, err := s.repo.GetAdminPermissions(ctx, adminID)
	if err != nil {
		return nil, err
	}
	// Filter to only currently-valid keys — drops any stale keys from old schema
	valid := map[string]bool{}
	for _, p := range models.AllAdminPermissions {
		valid[p] = true
	}
	filtered := make([]string, 0, len(perms))
	for _, p := range perms {
		if valid[p] {
			filtered = append(filtered, p)
		}
	}
	return filtered, nil
}

func (s *AdminService) SetAdminPermissions(ctx context.Context, adminID int, permissions []string) error {
	// Super admins have all permissions by definition — setting them is a no-op that should be blocked.
	isSuper, err := s.repo.IsAdminSuper(ctx, adminID)
	if err != nil {
		return err
	}
	if isSuper {
		return errors.New("cannot set permissions on a super admin — they have full access by default")
	}
	// Build valid set
	valid := map[string]bool{}
	for _, p := range models.AllAdminPermissions {
		valid[p] = true
	}
	// Silently drop unknown keys instead of rejecting — handles stale data from old permission sets
	filtered := make([]string, 0, len(permissions))
	for _, p := range permissions {
		if valid[p] {
			filtered = append(filtered, p)
		}
	}
	return s.repo.SetAdminPermissions(ctx, adminID, filtered)
}

func (s *AdminService) HasPermission(ctx context.Context, adminID int, permission string) (bool, error) {
	return s.repo.HasPermission(ctx, adminID, permission)
}

// ── Admin Activity Logs ───────────────────────────────────────────────────────

func (s *AdminService) InsertAdminLog(ctx context.Context, l models.AdminLog) {
	// Fire-and-forget — logging failures must never break the main action.
	// Use a background context so the log isn't dropped when the request context ends.
	logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.repo.InsertAdminLog(logCtx, l); err != nil {
		log.Printf("[WARN] InsertAdminLog action=%s adminID=%d: %v", l.Action, l.AdminID, err) //nolint:gosec // G117: error/value from internal source, not user-controlled input
	}
}

func (s *AdminService) ListAdminLogs(ctx context.Context, f models.AdminLogFilter) ([]models.AdminLog, int, error) {
	return s.repo.ListAdminLogs(ctx, f)
}

// DeleteAdminLogs removes admin activity log records to free storage.
// before="" → deletes all logs; before="YYYY-MM-DD" → deletes logs older than that date.
func (s *AdminService) DeleteAdminLogs(ctx context.Context, before string) (int64, error) {
	return s.repo.DeleteAdminLogs(ctx, before)
}

// ── Tier Management (Super Admin) ─────────────────────────────────────────────

func (s *AdminService) SetSellerTier(ctx context.Context, userID int, tier string, changedBy int, reason string) error {
	valid := map[string]bool{"normal": true, "gold": true, "premium": true}
	if !valid[tier] {
		return errors.New("tier must be normal, gold, or premium")
	}
	return s.repo.SetSellerTier(ctx, userID, tier, changedBy, reason)
}

func (s *AdminService) SetBuyerTier(ctx context.Context, userID int, tier string, changedBy int, reason string) error {
	valid := map[string]bool{"normal": true, "gold": true, "premium": true}
	if !valid[tier] {
		return errors.New("tier must be normal, gold, or premium")
	}
	return s.repo.SetBuyerTier(ctx, userID, tier, changedBy, reason)
}

func (s *AdminService) ListTierHistory(ctx context.Context, userID int) ([]models.TierHistory, error) {
	return s.repo.ListTierHistory(ctx, userID)
}

func (s *AdminService) ListPremiumUsers(ctx context.Context, role, tierType, tier string) ([]models.User, error) {
	return s.repo.ListPremiumUsers(ctx, role, tierType, tier)
}

// ── Seller Upgrade Requests ───────────────────────────────────────────────────

func (s *AdminService) ListUpgradeRequests(ctx context.Context, status string) ([]map[string]interface{}, error) {
	return s.repo.ListUpgradeRequests(ctx, status)
}

func (s *AdminService) ApproveUpgradeRequest(ctx context.Context, requestID, adminID int) error {
	var userID int
	var role, requestedTier string
	var amount float64
	if s.referralSvc != nil || s.sellerSvc != nil {
		userID, role, requestedTier, amount, _ = s.repo.GetUpgradeRequestDetails(ctx, requestID)
	}
	if err := s.repo.ApproveUpgradeRequest(ctx, requestID, adminID); err != nil {
		return err
	}
	// Notify the seller
	if s.sellerSvc != nil && userID > 0 {
		safeGo(func() {
			msg := fmt.Sprintf("Your upgrade request to %s tier has been approved. Your new tier is now active.", requestedTier) //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
			_ = s.sellerSvc.CreateNotification(context.Background(), userID, "Tier Upgrade Approved", msg, "success")
		})
	}
	// Fire commission after successful approval (with 30-day chargeback protection).
	if s.referralSvc != nil && userID > 0 && amount > 0 {
		var txType string
		if role == "seller" {
			txType = requestedTier + "_upgrade"
		} else {
			txType = "buyer_" + requestedTier + "_upgrade"
		}
		srcID := requestID
		safeGo(func() {
			s.referralSvc.RecordPendingCommission(
				context.Background(), userID, amount, txType,
				role, requestedTier,
				&srcID,
				30*24*time.Hour, // 30 days pending period for chargeback protection
			)
		})
	}
	return nil
}

func (s *AdminService) RejectUpgradeRequest(ctx context.Context, requestID, adminID int, reason string) error {
	var userID int
	var requestedTier string
	if s.sellerSvc != nil {
		userID, _, requestedTier, _, _ = s.repo.GetUpgradeRequestDetails(ctx, requestID)
	}
	if err := s.repo.RejectUpgradeRequest(ctx, requestID, adminID); err != nil {
		return err
	}
	if s.sellerSvc != nil && userID > 0 {
		safeGo(func() {
			msg := fmt.Sprintf("Your upgrade request to %s tier was not approved.", requestedTier) //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
			if reason != "" {
				msg += " Reason: " + reason
			} else {
				msg += " Please contact support for details."
			}
			_ = s.sellerSvc.CreateNotification(context.Background(), userID, "Tier Upgrade Rejected", msg, "warning")
		})
	}
	return nil
}

// ── Buyer Upgrade Requests ────────────────────────────────────────────────────

func (s *AdminService) ListBuyerUpgradeRequests(ctx context.Context, status string, page, limit int) ([]map[string]interface{}, int, error) {
	return s.repo.ListBuyerUpgradeRequests(ctx, status, page, limit)
}

func (s *AdminService) ApproveBuyerUpgradeRequest(ctx context.Context, requestID, adminID int) error {
	var userID int
	var requestedTier string
	var amount float64
	if s.referralSvc != nil || s.sellerSvc != nil {
		userID, _, requestedTier, amount, _ = s.repo.GetUpgradeRequestDetails(ctx, requestID)
	}
	if err := s.repo.ApproveBuyerUpgradeRequest(ctx, requestID, adminID); err != nil {
		return err
	}
	// Buyer upgrading to gold/premium qualifies for commission (they paid) — with 30-day chargeback protection.
	if s.referralSvc != nil && userID > 0 && amount > 0 {
		txType := "buyer_" + requestedTier + "_upgrade"
		srcID := requestID
		safeGo(func() {
			s.referralSvc.RecordPendingCommission(
				context.Background(), userID, amount, txType,
				"buyer", requestedTier,
				&srcID,
				30*24*time.Hour, // 30 days pending period for chargeback protection
			)
		})
	}
	return nil
}

func (s *AdminService) RejectBuyerUpgradeRequest(ctx context.Context, requestID, adminID int, reason string) error {
	return s.repo.RejectBuyerUpgradeRequest(ctx, requestID, adminID, reason)
}

func (s *AdminService) IssueWarning(ctx context.Context, userID, adminID int, reason string) error {
	if err := s.repo.IssueWarning(ctx, userID, adminID, reason); err != nil {
		return err
	}
	// Notify the user
	if s.sellerSvc != nil {
		safeGo(func() {
			_ = s.sellerSvc.CreateNotification(context.Background(), userID, "Official Warning", reason, "warning")
		})
	}
	return nil
}

func (s *AdminService) SuspendUser(ctx context.Context, userID, adminID int, reason string) error {
	if err := s.repo.SuspendUser(ctx, userID, adminID, reason); err != nil {
		return err
	}
	if s.sessionRevoker != nil {
		_ = s.sessionRevoker(ctx, userID, "user_suspended")
	}
	return nil
}

func (s *AdminService) BanUser(ctx context.Context, userID, adminID int, reason string) error {
	if err := s.repo.BanUser(ctx, userID, adminID, reason); err != nil {
		return err
	}
	if s.sessionRevoker != nil {
		_ = s.sessionRevoker(ctx, userID, "user_banned")
	}
	return nil
}

func (s *AdminService) CheckSystemHealth(ctx context.Context) (map[string]interface{}, error) {
	return s.repo.CheckSystemHealth(ctx)
}
func (s *AdminService) RepairMissingAdminEmails(ctx context.Context, adminID int, emailMap map[int]string) error {
	return s.repo.RepairMissingAdminEmails(ctx, adminID, emailMap)
}

// ── Operations Monitoring (Role 02) ───────────────────────────────────────────

func (s *AdminService) GetOperationsHealth(ctx context.Context) (*models.OperationsHealth, error) {
	return s.repo.GetOperationsHealth(ctx)
}

func (s *AdminService) ListJobStatusSummary(ctx context.Context) ([]models.JobStatusSummary, error) {
	return s.repo.ListJobStatusSummary(ctx)
}

func (s *AdminService) ListRecentFailedJobs(ctx context.Context, limit int) ([]models.JobExecutionLog, error) {
	return s.repo.ListRecentFailedJobs(ctx, limit)
}

// ── Seller Verification (Role 08) ─────────────────────────────────────────────

func (s *AdminService) ListVerificationDocuments(ctx context.Context, status string, page, limit int) ([]models.SellerVerificationDocument, int, error) {
	return s.repo.ListVerificationDocuments(ctx, status, page, limit)
}

func (s *AdminService) ReviewVerificationDocument(ctx context.Context, docID, adminID int, approved bool, reason, notes string) error {
	status := "approved"
	if !approved {
		status = "rejected"
	}
	return s.repo.ReviewVerificationDocument(ctx, docID, adminID, status, reason, notes)
}

// ── Trust & Safety (Role 21) ──────────────────────────────────────────────────

func (s *AdminService) ListAbuseReports(ctx context.Context, f models.AbuseReportFilter) ([]models.AbuseReport, int, error) {
	return s.repo.ListAbuseReports(ctx, f)
}

func (s *AdminService) ReviewAbuseReport(ctx context.Context, reportID, adminID int, status, actionTaken, notes string) error {
	reportedUserID, err := s.repo.ReviewAbuseReport(ctx, reportID, adminID, status, actionTaken, notes)
	if err != nil {
		return err
	}
	if reportedUserID != nil && *reportedUserID > 0 {
		if actionTaken == "banned" && s.sessionRevoker != nil {
			_ = s.sessionRevoker(ctx, *reportedUserID, "user_banned")
		} else if actionTaken == "suspended" && s.sessionRevoker != nil {
			_ = s.sessionRevoker(ctx, *reportedUserID, "user_suspended")
		}
	}
	return nil
}

func (s *AdminService) UpdateAbuseReportPriority(ctx context.Context, reportID int, priority string) error {
	return s.repo.UpdateAbuseReportPriority(ctx, reportID, priority)
}

// ── Refund Management (Role 20) ───────────────────────────────────────────────

func (s *AdminService) ListRefundRequests(ctx context.Context, f models.RefundRequestFilter) ([]models.RefundRequest, int, error) {
	return s.repo.ListRefundRequests(ctx, f)
}

func (s *AdminService) ReviewRefundRequest(ctx context.Context, refundID, adminID int, approved bool, reason string) error {
	return s.repo.ReviewRefundRequest(ctx, refundID, adminID, approved, reason)
}

func (s *AdminService) MarkRefundProcessed(ctx context.Context, refundID int, reference string) error {
	return s.repo.MarkRefundProcessed(ctx, refundID, reference)
}

// ── Account Appeals (Role 30) ─────────────────────────────────────────────────

func (s *AdminService) ListAccountAppeals(ctx context.Context, f models.AccountAppealFilter) ([]models.AccountAppeal, int, error) {
	return s.repo.ListAccountAppeals(ctx, f)
}

func (s *AdminService) ReviewAccountAppeal(ctx context.Context, appealID, adminID int, approved bool, decision, actionTaken string) error {
	return s.repo.ReviewAccountAppeal(ctx, appealID, adminID, approved, decision, actionTaken)
}

func (s *AdminService) UpdateAppealStatus(ctx context.Context, appealID int, status string) error {
	return s.repo.UpdateAppealStatus(ctx, appealID, status)
}

func (s *AdminService) CreateVerificationDocument(ctx context.Context, doc *models.SellerVerificationDocument) (int, error) {
	return s.repo.CreateVerificationDocument(ctx, doc)
}

func (s *AdminService) ListSellerVerificationDocuments(ctx context.Context, sellerID int) ([]models.SellerVerificationDocument, error) {
	return s.repo.ListSellerVerificationDocuments(ctx, sellerID)
}

func (s *AdminService) CreateAbuseReport(ctx context.Context, report *models.AbuseReport) (int, error) {
	return s.repo.CreateAbuseReport(ctx, report)
}

func (s *AdminService) CreateRefundRequest(ctx context.Context, req *models.RefundRequest) (int, error) {
	return s.repo.CreateRefundRequest(ctx, req)
}

func (s *AdminService) CreateAccountAppeal(ctx context.Context, appeal *models.AccountAppeal) (int, error) {
	return s.repo.CreateAccountAppeal(ctx, appeal)
}
