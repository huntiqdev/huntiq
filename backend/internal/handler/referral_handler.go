package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/referral"
	"github.com/your-org/used-tools-marketplace/pkg/response"
	"github.com/your-org/used-tools-marketplace/pkg/validations"
)

// ReferralHandler handles referral-related endpoints.
type ReferralHandler struct {
	svc           *service.ReferralService
	frontendURL   string                                                // e.g. https://yourapp.vercel.app — used to build referral links
	adminNotifier func(ctx context.Context, userID int, amount float64) // optional: fires admin notification
	adminLogger   func(ctx context.Context, adminID int, targetUserID *int, action, notes string)
}

func NewReferralHandler(svc *service.ReferralService, frontendURL string) *ReferralHandler {
	return &ReferralHandler{svc: svc, frontendURL: frontendURL}
}

// WithAdminLogger injects a function that records admin actions in the central log.
func (h *ReferralHandler) WithAdminLogger(fn func(ctx context.Context, adminID int, targetUserID *int, action, notes string)) *ReferralHandler {
	h.adminLogger = fn
	return h
}

// WithAdminNotifier injects a function that fires a super-admin notification on withdrawal request.
func (h *ReferralHandler) WithAdminNotifier(fn func(ctx context.Context, userID int, amount float64)) *ReferralHandler {
	h.adminNotifier = fn
	return h
}

// GET /api/seller/referral
func (h *ReferralHandler) GetSellerStats(c *fiber.Ctx) error {
	return h.getStats(c, "seller")
}

// GET /api/buyer/referral
func (h *ReferralHandler) GetBuyerStats(c *fiber.Ctx) error {
	return h.getStats(c, "buyer")
}

func (h *ReferralHandler) getStats(c *fiber.Ctx, role string) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	stats, err := h.svc.GetStats(c.Context(), userID)
	if err != nil {
		return response.InternalError(c, err)
	}
	// Build referral link with role hint so the register tab opens correctly
	referralLink := h.frontendURL + "/auth?ref=" + stats.ReferralCode + "&as=" + role
	return response.OK(c, fiber.Map{
		"referral_code":          stats.ReferralCode,
		"referral_link":          referralLink,
		"referral_suspended":     stats.ReferralSuspended,
		"invited_count":          stats.InvitedCount,
		"total_commission":       stats.TotalCommission,
		"inviter_code":           stats.InviterCode,
		"inviter_name":           stats.InviterName,
		"invitation_accepted_at": stats.InvitationAcceptedAt,
	})
}

// GET /api/seller/referral/commissions  or  GET /api/buyer/referral/commissions
func (h *ReferralHandler) GetCommissions(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	page, limit := parsePagination(c, 20, 100)
	commissions, total, err := h.svc.GetCommissions(c.Context(), userID, page, limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"commissions": commissions, "total": total})
}

// POST /api/seller/referral/accept  or  POST /api/buyer/referral/accept
// Body: { "referral_code": "REFxxxxxxxx" }
func (h *ReferralHandler) AcceptInvitation(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var body struct {
		ReferralCode string `json:"referral_code"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.ReferralCode == "" {
		return response.BadRequest(c, "referral_code is required")
	}
	// Use centralized referral code validation
	if !referral.ValidateReferralCode(body.ReferralCode) {
		return response.BadRequest(c, "invalid referral code format")
	}
	if err := h.svc.AcceptInvitation(c.Context(), userID, body.ReferralCode); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{
		"message": "invitation accepted — your inviter will earn 2% when you make a paid action",
	})
}

// GET /api/admin/referral/overview — super admin only
func (h *ReferralHandler) AdminOverview(c *fiber.Ctx) error {
	data, err := h.svc.AdminGetOverview(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, data)
}

// GET /api/admin/referral/list — super admin only
// Query: page, limit, search
func (h *ReferralHandler) AdminList(c *fiber.Ctx) error {
	page, limit := parsePagination(c, 50, 200)
	search := c.Query("search", "")
	rows, total, err := h.svc.AdminListReferrals(c.Context(), page, limit, search)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"referrals": rows, "total": total})
}

// PUT /api/admin/referral/users/:id/suspend — super admin only
// Body: { "suspended": true/false }
func (h *ReferralHandler) AdminSetSuspended(c *fiber.Ctx) error {
	userID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Suspended bool `json:"suspended"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if err := h.svc.AdminSetSuspended(c.Context(), userID, body.Suspended); err != nil {
		return response.BadRequest(c, err.Error())
	}
	action := "unsuspended"
	if body.Suspended {
		action = "suspended"
	}
	if h.adminLogger != nil {
		// Use c.Locals directly — mustUserID would write a second response on failure
		adminID, _ := c.Locals("userID").(int)
		h.adminLogger(c.Context(), adminID, &userID, "referral_"+action, "")
	}
	return response.OK(c, fiber.Map{"message": "referral code " + action})
}

// GET /api/seller/referral/balance  or  GET /api/buyer/referral/balance
func (h *ReferralHandler) GetBalance(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	bal, err := h.svc.GetUserBalance(c.Context(), userID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, bal)
}

// POST /api/seller/referral/withdraw  or  POST /api/buyer/referral/withdraw
// Body: { "amount": 1000, "payment_method": "telebirr", "payer_phone": "0912345678" }
func (h *ReferralHandler) RequestWithdrawal(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var body struct {
		Amount        float64 `json:"amount"`
		PaymentMethod string  `json:"payment_method"`
		PayerPhone    string  `json:"payer_phone"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.Amount <= 0 {
		return response.BadRequest(c, "amount is required")
	}
	// Use centralized phone validation
	phoneResult := validations.ValidateEthiopianPhone(body.PayerPhone)
	if !phoneResult.Valid {
		return response.BadRequest(c, "invalid payer phone number")
	}
	// Min/max limits are validated in the service layer against platform_settings
	id, err := h.svc.RequestWithdrawal(c.Context(), userID, body.Amount, body.PaymentMethod, phoneResult.Normalized)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	// Fire admin notification (non-fatal). Use context.Background() — c.Context() is
	// cancelled as soon as the HTTP response is sent, which would silently drop the DB insert.
	if h.adminNotifier != nil {
		go h.adminNotifier(context.Background(), userID, body.Amount)
	}
	return response.OK(c, fiber.Map{
		"id":      id,
		"message": "withdrawal request submitted — admin will process within 24 hours",
	})
}

// GET /api/seller/referral/withdrawals  or  GET /api/buyer/referral/withdrawals
func (h *ReferralHandler) GetMyWithdrawals(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	list, err := h.svc.GetMyWithdrawals(c.Context(), userID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, list)
}

// GET /api/admin/referral/withdrawals — super admin only
// Query: status (pending|approved|rejected|paid|all), page, limit
func (h *ReferralHandler) AdminListWithdrawals(c *fiber.Ctx) error {
	status := c.Query("status", "pending")
	switch status {
	case "pending", "approved", "rejected", "paid", "all":
	default:
		status = "pending"
	}
	page, limit := parsePagination(c, 50, 200)
	list, total, err := h.svc.AdminListWithdrawals(c.Context(), status, page, limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"withdrawals": list, "total": total})
}

// PUT /api/admin/referral/withdrawals/:id/approve — super admin only
// Body: { "note": "optional note" }
func (h *ReferralHandler) AdminApproveWithdrawal(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	wID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Note string `json:"note"`
	}
	_ = c.BodyParser(&body)
	if err := h.svc.AdminApproveWithdrawal(c.Context(), wID, adminID, body.Note); err != nil {
		return response.BadRequest(c, err.Error())
	}
	if h.adminLogger != nil {
		h.adminLogger(c.Context(), adminID, nil, "withdrawal_approve", "Withdrawal #"+strconv.Itoa(wID)+" "+body.Note)
	}
	return response.OK(c, fiber.Map{"message": "withdrawal approved — mark as paid once money is sent"})
}

// PUT /api/admin/referral/withdrawals/:id/reject — super admin only
// Body: { "reason": "..." }
func (h *ReferralHandler) AdminRejectWithdrawal(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	wID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil || body.Reason == "" {
		return response.BadRequest(c, "rejection reason is required")
	}
	if err := h.svc.AdminRejectWithdrawal(c.Context(), wID, adminID, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}
	if h.adminLogger != nil {
		h.adminLogger(c.Context(), adminID, nil, "withdrawal_reject", "Withdrawal #"+strconv.Itoa(wID)+" "+body.Reason)
	}
	return response.OK(c, fiber.Map{"message": "withdrawal rejected"})
}

// PUT /api/admin/referral/withdrawals/:id/paid — super admin only
// Body: { "reference": "TXN123456" }
func (h *ReferralHandler) AdminMarkWithdrawalPaid(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	wID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Reference string `json:"reference"`
	}
	if err := c.BodyParser(&body); err != nil || body.Reference == "" {
		return response.BadRequest(c, "payment reference is required")
	}
	if err := h.svc.AdminMarkWithdrawalPaid(c.Context(), wID, adminID, body.Reference); err != nil {
		return response.BadRequest(c, err.Error())
	}
	if h.adminLogger != nil {
		h.adminLogger(c.Context(), adminID, nil, "withdrawal_paid", "Withdrawal #"+strconv.Itoa(wID)+" Ref:"+body.Reference)
	}
	return response.OK(c, fiber.Map{"message": "withdrawal marked as paid"})
}

// PUT /api/admin/referral/users/:id/block — super admin only
// Body: { "blocked": true, "reason": "..." }
func (h *ReferralHandler) AdminSetBlocked(c *fiber.Ctx) error {
	userID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Blocked bool   `json:"blocked"`
		Reason  string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if err := h.svc.AdminSetPermanentlyBlocked(c.Context(), userID, body.Blocked, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}
	action := "unblocked"
	if body.Blocked {
		action = "permanently blocked"
	}
	if h.adminLogger != nil {
		// Use c.Locals directly — mustUserID would write a second response on failure
		adminID, _ := c.Locals("userID").(int)
		h.adminLogger(c.Context(), adminID, &userID, "referral_"+action, body.Reason)
	}
	return response.OK(c, fiber.Map{"message": "user referral account " + action})
}
