package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
	adminRepo      *repository.AdminRepo
}

func NewPaymentHandler(paymentService *service.PaymentService, adminRepo *repository.AdminRepo) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		adminRepo:      adminRepo,
	}
}

// GetPricing returns authoritative platform pricing and operating mode.
// GET /api/payments/pricing
func (h *PaymentHandler) GetPricing(c *fiber.Ctx) error {
	ctx := c.Context()
	prices, err := h.paymentService.GetPaymentRepo().GetTierPricing(ctx)
	if err != nil {
		return response.InternalError(c, err)
	}

	mode, _ := h.paymentService.GetPaymentRepo().GetGlobalOperatingMode(ctx)

	return response.OK(c, fiber.Map{
		"pricing":        prices,
		"operating_mode": mode,
	})
}

// Checkout initiates a unified payment order.
// POST /api/payments/checkout
func (h *PaymentHandler) Checkout(c *fiber.Ctx) error {
	callerID, ok := c.Locals("userID").(int)
	if !ok || callerID <= 0 {
		return response.Unauthorized(c, "authentication required")
	}

	userRole, _ := c.Locals("role").(string)

	var req struct {
		OrderType    string  `json:"order_type"` // tier_upgrade, listing_fee
		TargetTier   string  `json:"target_tier,omitempty"`
		DurationDays int     `json:"duration_days,omitempty"`
		ProductID    *int    `json:"product_id,omitempty"`
		Provider     string  `json:"provider"` // demo_mock, telebirr, ethio_telecom
		PayerPhone   string  `json:"payer_phone,omitempty"`
		OTPVerified  bool    `json:"otp_verified,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid checkout payload")
	}

	req.OrderType = strings.TrimSpace(req.OrderType)
	if req.OrderType == "" {
		req.OrderType = models.OrderTypeTierUpgrade
	}

	switch req.OrderType {
	case models.OrderTypeTierUpgrade, "membership":
		if req.TargetTier == "" {
			return response.BadRequest(c, "target_tier (gold or premium) is required for upgrade")
		}
		if req.DurationDays <= 0 {
			req.DurationDays = 30
		}

		order, err := h.paymentService.InitiateUpgradePayment(
			c.Context(),
			callerID,
			userRole,
			req.TargetTier,
			req.DurationDays,
			req.Provider,
			req.PayerPhone,
			req.OTPVerified,
		)
		if err != nil {
			return response.BadRequest(c, err.Error())
		}
		return response.OK(c, order)

	case models.OrderTypeListingFee:
		if req.ProductID == nil || *req.ProductID <= 0 {
			return response.BadRequest(c, "product_id is required for listing fee payment")
		}

		order, err := h.paymentService.InitiateListingFeePayment(
			c.Context(),
			callerID,
			*req.ProductID,
			req.Provider,
			req.PayerPhone,
		)
		if err != nil {
			return response.BadRequest(c, err.Error())
		}
		return response.OK(c, order)

	default:
		return response.BadRequest(c, fmt.Sprintf("unsupported order_type: %s", req.OrderType))
	}
}

// GetOrder fetches a payment order by its order_id.
// GET /api/payments/orders/:orderId
func (h *PaymentHandler) GetOrder(c *fiber.Ctx) error {
	callerID, ok := c.Locals("userID").(int)
	if !ok || callerID <= 0 {
		return response.Unauthorized(c, "authentication required")
	}

	role, _ := c.Locals("role").(string)
	orderID := c.Params("orderId")
	if orderID == "" {
		return response.BadRequest(c, "orderId parameter is required")
	}

	order, err := h.paymentService.GetPaymentRepo().GetOrderByOrderID(c.Context(), orderID)
	if err != nil || order == nil {
		return response.NotFound(c, "payment order not found")
	}

	// Security: users may only read their own orders; admins can view all
	if role != "admin" && order.UserID != callerID {
		return response.Forbidden(c, "access denied")
	}

	return response.OK(c, order)
}

// VerifyOrder triggers verification of payment with the provider.
// POST /api/payments/verify/:orderId
func (h *PaymentHandler) VerifyOrder(c *fiber.Ctx) error {
	callerID, ok := c.Locals("userID").(int)
	if !ok || callerID <= 0 {
		return response.Unauthorized(c, "authentication required")
	}

	role, _ := c.Locals("role").(string)
	orderID := c.Params("orderId")
	if orderID == "" {
		return response.BadRequest(c, "orderId parameter is required")
	}

	order, err := h.paymentService.GetPaymentRepo().GetOrderByOrderID(c.Context(), orderID)
	if err != nil || order == nil {
		return response.NotFound(c, "payment order not found")
	}

	if role != "admin" && order.UserID != callerID {
		return response.Forbidden(c, "access denied")
	}

	updated, err := h.paymentService.VerifyAndActivatePayment(c.Context(), orderID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, updated)
}

// SimulateDemoOutcome allows simulated testing of outcomes in demo mode.
// POST /api/payments/demo/simulate
func (h *PaymentHandler) SimulateDemoOutcome(c *fiber.Ctx) error {
	callerID, ok := c.Locals("userID").(int)
	if !ok || callerID <= 0 {
		return response.Unauthorized(c, "authentication required")
	}

	role, _ := c.Locals("role").(string)

	var req struct {
		OrderID string `json:"order_id"`
		Action  string `json:"action"` // success, failed, cancelled, expired
	}

	if err := c.BodyParser(&req); err != nil || req.OrderID == "" || req.Action == "" {
		return response.BadRequest(c, "order_id and action (success, failed, cancelled, expired) are required")
	}

	order, err := h.paymentService.GetPaymentRepo().GetOrderByOrderID(c.Context(), req.OrderID)
	if err != nil || order == nil {
		return response.NotFound(c, "order not found")
	}

	if role != "admin" && order.UserID != callerID {
		return response.Forbidden(c, "access denied")
	}

	updated, err := h.paymentService.SimulatePaymentAction(c.Context(), req.OrderID, req.Action)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, updated)
}

// Webhook processes provider callback notifications.
// POST /api/payments/webhook/:provider
func (h *PaymentHandler) Webhook(c *fiber.Ctx) error {
	providerCode := c.Params("provider")
	rawBody := c.Body()
	sig := c.Get("X-Signature")
	if sig == "" {
		sig = c.Get("Signature")
	}

	order, err := h.paymentService.HandleWebhook(c.Context(), providerCode, rawBody, sig)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"order_id": order.OrderID,
		"status":   order.Status,
	})
}

// SendMembershipOTP dispatches Fida OTP for membership upgrade confirmation.
// POST /api/membership/otp/send
func (h *PaymentHandler) SendMembershipOTP(c *fiber.Ctx) error {
	callerID, ok := c.Locals("userID").(int)
	if !ok || callerID <= 0 {
		return response.Unauthorized(c, "authentication required")
	}

	role, _ := c.Locals("role").(string)

	var req struct {
		Phone      string `json:"phone"`
		TargetTier string `json:"target_tier"`
	}
	if err := c.BodyParser(&req); err != nil || strings.TrimSpace(req.Phone) == "" {
		return response.BadRequest(c, "valid phone number is required")
	}
	if req.TargetTier == "" {
		req.TargetTier = "gold"
	}

	_, err := h.paymentService.SendUpgradeOTP(c.Context(), callerID, req.Phone, req.TargetTier, role)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"message": "Verification code dispatched via SMS",
	})
}

// VerifyMembershipOTP checks the OTP code for upgrade confirmation.
// POST /api/membership/otp/verify
func (h *PaymentHandler) VerifyMembershipOTP(c *fiber.Ctx) error {
	callerID, ok := c.Locals("userID").(int)
	if !ok || callerID <= 0 {
		return response.Unauthorized(c, "authentication required")
	}

	role, _ := c.Locals("role").(string)

	var req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	if err := c.BodyParser(&req); err != nil || strings.TrimSpace(req.Phone) == "" || strings.TrimSpace(req.Code) == "" {
		return response.BadRequest(c, "phone and code are required")
	}

	err := h.paymentService.VerifyUpgradeOTP(c.Context(), callerID, req.Phone, req.Code, role)
	if err != nil {
		return response.BadRequest(c, "invalid or expired verification code")
	}

	return response.OK(c, fiber.Map{
		"verified": true,
		"message":  "Phone verified for upgrade",
	})
}

// ── Admin Endpoints ──────────────────────────────────────────────────────────

// AdminListOrders lists payment orders with search and filtering.
// GET /api/admin/payments/orders
func (h *PaymentHandler) AdminListOrders(c *fiber.Ctx) error {
	filter := repository.PaymentOrderFilter{
		Provider:    c.Query("provider"),
		Status:      c.Query("status"),
		Environment: c.Query("environment"),
		Search:      c.Query("search"),
		Limit:       c.QueryInt("limit", 20),
		Offset:      c.QueryInt("offset", 0),
	}

	if uidStr := c.Query("user_id"); uidStr != "" {
		if uid, err := strconv.Atoi(uidStr); err == nil && uid > 0 {
			filter.UserID = &uid
		}
	}

	orders, total, err := h.paymentService.GetPaymentRepo().ListOrders(c.Context(), filter)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{
		"orders": orders,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// AdminGetOrderDetail fetches order detail including callback payload.
// GET /api/admin/payments/orders/:id
func (h *PaymentHandler) AdminGetOrderDetail(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// try by orderID string
		order, err2 := h.paymentService.GetPaymentRepo().GetOrderByOrderID(c.Context(), idParam)
		if err2 != nil || order == nil {
			return response.NotFound(c, "order not found")
		}
		return response.OK(c, order)
	}

	order, err := h.paymentService.GetPaymentRepo().GetOrderByID(c.Context(), id)
	if err != nil || order == nil {
		return response.NotFound(c, "order not found")
	}

	return response.OK(c, order)
}

// AdminGetProviders returns provider setups and active operating mode.
// GET /api/admin/payments/providers
func (h *PaymentHandler) AdminGetProviders(c *fiber.Ctx) error {
	configs, err := h.paymentService.GetPaymentRepo().GetProviderConfigs(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}

	mode, _ := h.paymentService.GetPaymentRepo().GetGlobalOperatingMode(c.Context())

	return response.OK(c, fiber.Map{
		"providers":      configs,
		"operating_mode": mode,
	})
}

// AdminUpdateProvider updates a provider's settings.
// PUT /api/admin/payments/providers/:code
func (h *PaymentHandler) AdminUpdateProvider(c *fiber.Ctx) error {
	adminID, _ := c.Locals("userID").(int)
	code := c.Params("code")

	var req struct {
		EnvironmentMode string `json:"environment_mode"`
		IsEnabled       bool   `json:"is_enabled"`
		APIEndpoint     string `json:"api_endpoint"`
		AppID           string `json:"app_id"`
		AppKey          string `json:"app_key"`
		ShortCode       string `json:"short_code"`
		NotifyURL       string `json:"notify_url"`
		ReturnURL       string `json:"return_url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body")
	}

	if req.EnvironmentMode == "" {
		req.EnvironmentMode = "demo"
	}

	err := h.paymentService.GetPaymentRepo().UpdateProviderConfig(
		c.Context(),
		code,
		req.EnvironmentMode,
		req.IsEnabled,
		req.APIEndpoint,
		req.AppID,
		req.AppKey,
		req.ShortCode,
		req.NotifyURL,
		req.ReturnURL,
		adminID,
	)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{
		"message": fmt.Sprintf("Provider %s configuration updated successfully", code),
	})
}

// AdminSwitchOperatingMode switches between demo and production mode with fail-closed validation.
// POST /api/admin/payments/switch-mode (Platform Owner only)
func (h *PaymentHandler) AdminSwitchOperatingMode(c *fiber.Ctx) error {
	adminID, _ := c.Locals("userID").(int)

	var req struct {
		Mode string `json:"mode"` // demo or production
	}

	if err := c.BodyParser(&req); err != nil || (req.Mode != "demo" && req.Mode != "production") {
		return response.BadRequest(c, "mode must be either 'demo' or 'production'")
	}

	if err := h.paymentService.SetOperatingMode(c.Context(), adminID, req.Mode); err != nil {
		return response.BadRequest(c, fmt.Sprintf("Mode switch failed (FAIL-CLOSED): %v", err))
	}

	// Audit log entry
	if h.adminRepo != nil {
		_ = h.adminRepo.InsertAdminLog(c.Context(), models.AdminLog{
			AdminID: adminID,
			Action:  "SWITCH_PAYMENT_OPERATING_MODE",
			Notes:   fmt.Sprintf("Payment operating mode changed to: %s", req.Mode),
		})
	}

	return response.OK(c, fiber.Map{
		"message":        fmt.Sprintf("Payment operating mode successfully set to %s", req.Mode),
		"operating_mode": req.Mode,
	})
}

// AdminTestOTP provides an authenticated endpoint for dev testing OTP codes.
// GET /api/admin/sms/test-otp (Admin only, requires otp_logs permission)
func (h *PaymentHandler) AdminTestOTP(c *fiber.Ctx) error {
	phone := c.Query("phone")
	if phone == "" {
		return response.BadRequest(c, "phone query parameter is required")
	}

	normalized, err := service.NormalizeEthiopianPhone(phone)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Query latest unused OTP code for phone
	return response.OK(c, fiber.Map{
		"phone":   normalized,
		"message": "OTP lookup available via /api/admin/otp-logs with proper decrypt permissions",
	})
}
