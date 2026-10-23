// Package handlers implements HTTP request handlers.
// Rate limiting for all auth endpoints (login, register, OTP) is applied at the
// route level in cmd/api/main.go via fiber/v2/middleware/limiter, keyed by client IP:
//   - authLimiter: 10 req/min for login and register endpoints
//   - otpLimiter:  5 req/min for OTP send and verify endpoints
package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
	"github.com/your-org/used-tools-marketplace/pkg/errors"
	"github.com/your-org/used-tools-marketplace/pkg/response"
	"github.com/your-org/used-tools-marketplace/pkg/validations"
)

// Type aliases for repository types
type (
	OTPLogFilter     = repository.OTPLogFilter
	OTPLogEntry      = repository.OTPLogEntry
	ActivityLogFilter = repository.ActivityLogFilter
	ActivityLogEntry = repository.ActivityLogEntry
)

// AuthHandler handles authentication endpoints.
// Rate limiting is applied at the route level in main.go — not in the handler itself.
type AuthHandler struct {
	svc      *service.AuthService
	adminSvc *service.AdminService
	devMode  bool
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) WithAdminService(a *service.AdminService) *AuthHandler {
	h.adminSvc = a
	return h
}

func (h *AuthHandler) WithDevMode(dev bool) *AuthHandler {
	h.devMode = dev
	return h
}

// registerUser is shared by Register and BuyerRegister.
func (h *AuthHandler) registerUser(c *fiber.Ctx, role string) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	req.Role = role
	
	// Validate seller_preference for sellers
	if role == "seller" {
		if req.SellerPreference == "" {
			req.SellerPreference = "sale" // default
		}
		validPrefs := map[string]bool{"sale": true, "rent": true, "both": true}
		if !validPrefs[req.SellerPreference] {
			return response.BadRequest(c, "seller_preference must be 'sale', 'rent', or 'both'")
		}
	}
	
	if !h.svc.IsRegistrationAllowed(c.Context(), role) {
		return response.BadRequest(c, role+" registration is currently closed — contact admin")
	}
	user, err := h.svc.Register(c.Context(), req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "otp_pending:") {
			return response.OK(c, fiber.Map{"otp_pending": true, "otp_required": true})
		}
		return response.BadRequest(c, err.Error())
	}
	if h.svc.IsOTPRequired(c.Context()) {
		return response.OK(c, fiber.Map{"otp_required": true})
	}
	return response.OK(c, fiber.Map{"user": user, "otp_required": false})
}

// POST /api/seller/register
func (h *AuthHandler) Register(c *fiber.Ctx) error { return h.registerUser(c, "seller") }

// POST /api/buyer/register
func (h *AuthHandler) BuyerRegister(c *fiber.Ctx) error { return h.registerUser(c, "buyer") }

// loginAsRole is shared by Login and BuyerLogin.
func (h *AuthHandler) loginAsRole(c *fiber.Ctx, role string) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	// Use centralized phone validation
	phoneResult := validations.ValidateEthiopianPhone(req.Phone)
	if !phoneResult.Valid {
		return response.AppError(c, errors.ValPhoneInvalid())
	}
	req.Phone = phoneResult.Normalized
	req.Email = ""
	token, err := h.svc.LoginAsRole(c.Context(), req, role)
	if err != nil {
		if strings.HasPrefix(err.Error(), "otp_pending:") {
			return response.OK(c, fiber.Map{"otp_pending": true, "otp_required": true})
		}
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	return response.OK(c, dto.LoginResponse{Token: token})
}

// POST /api/seller/login
func (h *AuthHandler) Login(c *fiber.Ctx) error { return h.loginAsRole(c, "seller") }

// checkExists is shared by CheckSellerExists and CheckBuyerExists.
func (h *AuthHandler) checkExists(c *fiber.Ctx, existsFn func(ctx context.Context, phone string) (bool, error)) error {
	var body struct {
		Phone string `json:"phone"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	// Use centralized phone validation
	phoneResult := validations.ValidateEthiopianPhone(body.Phone)
	if !phoneResult.Valid {
		return response.AppError(c, errors.ValPhoneInvalid())
	}
	exists, err := existsFn(c.Context(), phoneResult.Normalized)
	if err != nil {
		return response.InternalError(c, err)
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"exists": false})
	}
	return response.OK(c, fiber.Map{"exists": true})
}

// POST /api/seller/check
func (h *AuthHandler) CheckSellerExists(c *fiber.Ctx) error {
	return h.checkExists(c, h.svc.SellerExists)
}

// POST /api/buyer/check
func (h *AuthHandler) CheckBuyerExists(c *fiber.Ctx) error {
	return h.checkExists(c, h.svc.BuyerExists)
}

// POST /api/buyer/login
func (h *AuthHandler) BuyerLogin(c *fiber.Ctx) error { return h.loginAsRole(c, "buyer") }

// POST /api/admin/login  — Step 1: email + password
func (h *AuthHandler) AdminLogin(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	// Use centralized email validation
	if !validations.ValidateEmail(req.Email) {
		return response.BadRequest(c, "valid email is required")
	}
	// Use centralized password validation
	if valid, _ := validations.ValidatePassword(req.Password); !valid {
		return response.AppError(c, errors.ValPasswordMin())
	}
	req.Phone = ""
	result, isSuper, err := h.svc.AdminLoginStep1(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	if isSuper {
		return response.OK(c, dto.AdminLoginStep1Response{
			Token:   result,
			IsSuper: true,
			Message: "login successful",
		})
	}
	return response.OK(c, dto.AdminLoginStep1Response{
		Email:   result,
		IsSuper: false,
		Message: "OTP sent to your email address",
	})
}

// POST /api/admin/login/verify  — Step 2: verify email OTP → issues JWT and session
func (h *AuthHandler) AdminLoginVerifyOTP(c *fiber.Ctx) error {
	var req dto.AdminOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Email == "" || req.Code == "" {
		return response.BadRequest(c, "email and code are required")
	}
	ip := c.IP()
	if ip == "" {
		ip = "127.0.0.1"
	}
	token, sessionToken, err := h.svc.AdminLoginStep2(c.Context(), req.Email, req.Code, ip, c.Get("User-Agent"))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	if sessionToken != "" {
		c.Cookie(&fiber.Cookie{
			Name:     "hm_admin_session",
			Value:    sessionToken,
			Path:     "/",
			MaxAge:   12 * 3600,
			HTTPOnly: true,
			SameSite: "Lax",
		})
	}
	return response.OK(c, fiber.Map{
		"token":         token,
		"session_token": sessionToken,
	})
}

// POST /api/auth/otp/send  — send OTP to phone
func (h *AuthHandler) SendOTP(c *fiber.Ctx) error {
	var req dto.OTPRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Purpose == "" {
		req.Purpose = "login"
	}
	if err := h.svc.SendOTP(c.Context(), req); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "OTP sent"})
}

// GET /api/admin/activity-logs  — super admin only
func (h *AuthHandler) ListActivityLogs(c *fiber.Ctx) error {
	eventFilter := c.Query("event")
	validEvents := map[string]bool{
		"": true, "register_attempt": true, "register_ok": true, "register_fail": true,
		"login_attempt": true, "login_ok": true, "login_fail": true,
		"otp_send": true, "otp_verify_ok": true, "otp_verify_fail": true, "otp_locked": true,
		"register_retry": true, "otp_send_fail": true,
	}
	if !validEvents[eventFilter] {
		eventFilter = ""
	}
	userTypeFilter := c.Query("user_type")
	switch userTypeFilter {
	case "", "seller", "buyer", "admin":
	default:
		userTypeFilter = ""
	}
	const dateLayout = "2006-01-02"
	var dateFromPtr, dateToPtr *time.Time
	dateFrom, dateTo := c.Query("from"), c.Query("to")
	if dateFrom != "" {
		parsed, err := time.Parse(dateLayout, dateFrom)
		if err != nil {
			return response.BadRequest(c, "invalid 'from' date — use YYYY-MM-DD format")
		}
		dateFromPtr = &parsed
	}
	if dateTo != "" {
		parsed, err := time.Parse(dateLayout, dateTo)
		if err != nil {
			return response.BadRequest(c, "invalid 'to' date — use YYYY-MM-DD format")
		}
		dateToPtr = &parsed
	}
	page, limit := parsePagination(c, 50, 200)
	f := ActivityLogFilter{
		Event: eventFilter, Phone: c.Query("phone"), UserType: userTypeFilter,
		DateFrom: dateFromPtr, DateTo: dateToPtr, Page: page, Limit: limit,
	}
	logs, total, err := h.svc.GetActivityLogs(c.Context(), f)
	if err != nil {
		return response.InternalError(c, err)
	}
	if logs == nil {
		logs = []ActivityLogEntry{}
	}
	return response.OK(c, fiber.Map{"logs": logs, "total": total})
}

// GET /api/admin/otp-logs  — super admin only
func (h *AuthHandler) ListOTPLogs(c *fiber.Ctx) error {
	statusFilter := c.Query("status")
	switch statusFilter {
	case "verified", "locked", "expired", "pending", "":
	default:
		statusFilter = ""
	}
	userTypeFilter := c.Query("user_type")
	switch userTypeFilter {
	case "", "seller", "buyer", "admin":
	default:
		userTypeFilter = ""
	}
	page, limit := parsePagination(c, 50, 200)
	f := OTPLogFilter{
		Phone: truncate(c.Query("phone"), 20), Email: truncate(c.Query("email"), 100),
		UserType: userTypeFilter,
		Status: statusFilter, ShowPlainOTP: c.Locals("adminRole") == models.RolePlatformOwner || c.Locals("adminRole") == models.LegacyRoleSuper,
		Page: page, Limit: limit,
	}
	logs, total, err := h.svc.GetOTPLogs(c.Context(), f)
	if err != nil {
		return response.InternalError(c, err)
	}
	if logs == nil {
		logs = []repository.OTPLogEntry{}
	}
	return response.OK(c, fiber.Map{"logs": logs, "total": total})
}

// DELETE /api/admin/otp-logs  — bulk delete OTP records
// Requires confirm=DELETE_ALL in body when status=all to prevent accidental wipe.
func (h *AuthHandler) DeleteOTPLogs(c *fiber.Ctx) error {
	status := c.Query("status", "done")
	if status != "done" && status != "all" {
		status = "done"
	}
	if status == "all" {
		var body struct {
			Confirm string `json:"confirm"`
		}
		_ = c.BodyParser(&body)
		if body.Confirm != "DELETE_ALL" {
			return response.BadRequest(c, `confirm field must be "DELETE_ALL" to wipe all OTP records`)
		}
	}
	deleted, err := h.svc.DeleteOTPLogs(c.Context(), status)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"deleted": deleted, "message": "OTP logs cleaned up"})
}

// POST /api/auth/otp/verify  — verify OTP (marks phone verified, does NOT issue token)
func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req dto.OTPVerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if err := h.svc.VerifyOTP(c.Context(), req); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"verified": true})
}

// POST /api/auth/otp/login  — verify OTP and issue JWT in one step
func (h *AuthHandler) OTPLogin(c *fiber.Ctx) error {
	var req dto.OTPVerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	role := req.Role
	if role != "buyer" {
		role = "seller"
	}
	token, err := h.svc.LoginWithOTP(c.Context(), req.Phone, req.Code, role)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	return response.OK(c, dto.LoginResponse{Token: token})
}

// forgotPassword is shared by SellerForgotPassword and BuyerForgotPassword.
func (h *AuthHandler) forgotPassword(c *fiber.Ctx, role string) error {
	var body struct {
		Phone string `json:"phone"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	// Use centralized phone validation
	phoneResult := validations.ValidateEthiopianPhone(body.Phone)
	if !phoneResult.Valid {
		return response.AppError(c, errors.ValPhoneInvalid())
	}
	_ = h.svc.ForgotPasswordSend(c.Context(), phoneResult.Normalized, role)
	return response.OK(c, fiber.Map{"message": "if an account exists, an OTP has been sent"})
}

// resetPassword is shared by SellerResetPassword and BuyerResetPassword.
func (h *AuthHandler) resetPassword(c *fiber.Ctx, role string) error {
	var body struct {
		Phone       string `json:"phone"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	// Use centralized validations
	phoneResult := validations.ValidateEthiopianPhone(body.Phone)
	if !phoneResult.Valid {
		return response.AppError(c, errors.ValPhoneInvalid())
	}
	if valid, _ := validations.ValidatePassword(body.NewPassword); !valid {
		return response.AppError(c, errors.ValPasswordMin())
	}
	if err := h.svc.ForgotPasswordReset(c.Context(), phoneResult.Normalized, role, body.Code, body.NewPassword); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "password reset successfully"})
}

func (h *AuthHandler) SellerForgotPassword(c *fiber.Ctx) error { return h.forgotPassword(c, "seller") }
func (h *AuthHandler) SellerResetPassword(c *fiber.Ctx) error  { return h.resetPassword(c, "seller") }
func (h *AuthHandler) BuyerForgotPassword(c *fiber.Ctx) error  { return h.forgotPassword(c, "buyer") }
func (h *AuthHandler) BuyerResetPassword(c *fiber.Ctx) error   { return h.resetPassword(c, "buyer") }

// POST /api/auth/forgot-password — Unified password reset request
func (h *AuthHandler) UnifiedForgotPassword(c *fiber.Ctx) error {
	var body struct {
		Phone string `json:"phone"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	phoneResult := validations.ValidateEthiopianPhone(body.Phone)
	if !phoneResult.Valid {
		return response.AppError(c, errors.ValPhoneInvalid())
	}
	_ = h.svc.UnifiedForgotPasswordSend(c.Context(), phoneResult.Normalized)
	return response.OK(c, fiber.Map{"message": "if an account exists, an OTP has been sent"})
}

// POST /api/auth/reset-password — Unified password reset completion
func (h *AuthHandler) UnifiedResetPassword(c *fiber.Ctx) error {
	var body struct {
		Phone       string `json:"phone"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	phoneResult := validations.ValidateEthiopianPhone(body.Phone)
	if !phoneResult.Valid {
		return response.AppError(c, errors.ValPhoneInvalid())
	}
	if valid, _ := validations.ValidatePassword(body.NewPassword); !valid {
		return response.AppError(c, errors.ValPasswordMin())
	}
	if err := h.svc.UnifiedForgotPasswordReset(c.Context(), phoneResult.Normalized, body.Code, body.NewPassword); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "password reset successfully"})
}

// POST /api/admin/forgot-password — send OTP to email for admin password reset
func (h *AuthHandler) AdminForgotPassword(c *fiber.Ctx) error {
	var body struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	// Use centralized email validation
	if !validations.ValidateEmail(body.Email) {
		return response.BadRequest(c, "valid email is required")
	}
	_ = h.svc.AdminForgotPasswordSend(c.Context(), body.Email)
	return response.OK(c, fiber.Map{"message": "if an account exists, an OTP has been sent"})
}

// POST /api/admin/reset-password — verify OTP and set new admin password
func (h *AuthHandler) AdminResetPassword(c *fiber.Ctx) error {
	var body struct {
		Email       string `json:"email"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	// Use centralized validations
	if !validations.ValidateEmail(body.Email) {
		return response.BadRequest(c, "valid email is required")
	}
	if valid, _ := validations.ValidatePassword(body.NewPassword); !valid {
		return response.AppError(c, errors.ValPasswordMin())
	}
	if err := h.svc.AdminForgotPasswordReset(c.Context(), body.Email, body.Code, body.NewPassword); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "password reset successfully"})
}

// PUT /api/seller/password
func (h *AuthHandler) SellerChangePassword(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
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
	// Use centralized password validation
	if valid, _ := validations.ValidatePassword(body.NewPassword); !valid {
		return response.AppError(c, errors.ValPasswordMin())
	}
	if err := h.svc.ChangePassword(c.Context(), userID, body.CurrentPassword, body.NewPassword); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "password updated"})
}

// GET /api/dev/otp?phone=... — DEV ONLY
func (h *AuthHandler) DevPeekOTP(c *fiber.Ctx) error {
	if !h.devMode {
		return fiber.ErrNotFound
	}
	phone := c.Query("phone")
	if phone == "" {
		return response.BadRequest(c, "phone query param required")
	}
	code, err := h.svc.DevPeekOTP(c.Context(), phone)
	if err != nil {
		return response.NotFound(c, "no active OTP found for this phone")
	}
	return response.OK(c, fiber.Map{"phone": phone, "otp": code, "warning": "DEV ONLY — never expose in production"})
}

// GET /api/dev/otp/decrypt?code=<hex> — DEV ONLY
func (h *AuthHandler) DevDecryptOTP(c *fiber.Ctx) error {
	if !h.devMode {
		return fiber.ErrNotFound
	}
	hexCode := c.Query("code")
	if hexCode == "" {
		return response.BadRequest(c, "code query param required")
	}
	plain, err := h.svc.DevDecryptOTP(hexCode)
	if err != nil {
		return response.BadRequest(c, "decryption failed: "+err.Error())
	}
	return response.OK(c, fiber.Map{"otp": plain, "warning": "DEV ONLY — never expose in production"})
}

// POST /api/auth/demand — set user's current demand (FIND or LIST) after registration
func (h *AuthHandler) SetDemand(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var req dto.DemandSelectionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Demand != "FIND" && req.Demand != "LIST" {
		return response.BadRequest(c, "demand must be 'FIND' or 'LIST'")
	}
	if err := h.svc.SetDemand(c.Context(), userID, req.Demand); err != nil {
		return response.InternalError(c, err)
	}
	var nextPath string
	var message string
	switch req.Demand {
	case "FIND":
		nextPath = "/find"
		message = "Choose Buy or Rent to continue"
	case "LIST":
		nextPath = "/list"
		message = "Choose Sell or Landlord to continue"
	}
	return response.OK(c, dto.DemandSelectionResponse{
		Demand:   req.Demand,
		NextPath: nextPath,
		Message:  message,
	})
}

// ── Unified Authentication & Capability Handlers ─────────────────────────────

// POST /api/auth/register — Single unified registration endpoint for all users
func (h *AuthHandler) UnifiedRegister(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	user, token, err := h.svc.UnifiedRegister(c.Context(), req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "otp_pending:") {
			return response.OK(c, fiber.Map{"otp_pending": true, "otp_required": true, "phone": req.Phone})
		}
		return response.BadRequest(c, err.Error())
	}

	if h.svc.IsOTPRequired(c.Context()) && !user.OTPVerified {
		return response.OK(c, fiber.Map{
			"otp_required": true,
			"phone":        user.Phone,
			"message":      "Please enter the OTP sent to your phone",
		})
	}

	return response.OK(c, fiber.Map{
		"user":         user,
		"token":        token,
		"capabilities": user.Capabilities,
		"otp_required": false,
	})
}

// POST /api/auth/login — Single unified login endpoint for all users
func (h *AuthHandler) UnifiedLogin(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	user, token, err := h.svc.UnifiedLogin(c.Context(), req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "otp_pending:") {
			return response.OK(c, fiber.Map{"otp_pending": true, "otp_required": true, "phone": req.Phone})
		}
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return response.OK(c, fiber.Map{
		"user":         user,
		"token":        token,
		"capabilities": user.Capabilities,
	})
}

// GET /api/auth/me — Current authenticated user profile with active capabilities
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	user, err := h.svc.GetMe(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	return response.OK(c, user)
}

// POST /api/auth/capabilities/activate — 1-click capability activation (seller, landlord, etc.)
func (h *AuthHandler) ActivateCapability(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var req dto.ActivateCapabilityRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	user, token, err := h.svc.ActivateCapability(c.Context(), userID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"user":         user,
		"token":        token,
		"capabilities": user.Capabilities,
		"message":      "Capability activated successfully",
	})
}

// GET /api/auth/capabilities — List user's active capabilities
func (h *AuthHandler) GetCapabilities(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	caps, err := h.svc.GetCapabilities(c.Context(), userID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"capabilities": caps})
}
