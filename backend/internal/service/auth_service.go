package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/pkg/constants"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
	"github.com/your-org/used-tools-marketplace/pkg/validations"
)

// randHex returns n random hex bytes (2n hex chars). Returns error on entropy failure.
func randHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("crypto/rand unavailable: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// sanitizeLog strips newlines and control characters from a string before logging
// to prevent log injection (CWE-117).
func sanitizeLog(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r < 0x20 {
			return '_'
		}
		return r
	}, s)
}

type AuthService struct {
	repo           *repository.UserRepo
	jwtSecret      string
	devMode        bool
	adminSvc       firstLoginNotifier
	emailSender    EmailSender
	referralSvc    *ReferralService
	otpCrypto      *security.OTPCrypto
	settingsReader func(ctx context.Context, key string) string
	// cached otp_required setting — refreshed every 60s to avoid a DB hit on every login
	otpRequiredCache    bool
	otpRequiredCachedAt time.Time
	otpRequiredMu       sync.RWMutex
}

// EmailSender is the interface for sending OTP emails.
// Inject a real implementation (SendGrid, SES, SMTP, etc.) via WithEmailSender.
type EmailSender interface {
	SendOTP(to, code, purpose string) error
}

// firstLoginNotifier is a minimal interface to avoid circular imports.
type firstLoginNotifier interface {
	NotifyFirstLogin(ctx context.Context, name, role string)
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

// WithDevMode enables OTP console logging for development.
func (s *AuthService) WithDevMode(dev bool) *AuthService {
	s.devMode = dev
	return s
}

// WithEmailSender injects an email delivery implementation.
// Call this at startup to enable admin OTP and password-reset emails.
func (s *AuthService) WithEmailSender(e EmailSender) *AuthService {
	s.emailSender = e
	return s
}

// WithAdminService injects the admin service for first-login notifications.
func (s *AuthService) WithAdminService(a firstLoginNotifier) *AuthService {
	s.adminSvc = a
	return s
}

// WithReferralService injects the referral service for commission recording on upgrades.
func (s *AuthService) WithReferralService(r *ReferralService) *AuthService {
	s.referralSvc = r
	return s
}

// WithOTPCrypto injects the OTP encryption/decryption handler
func (s *AuthService) WithOTPCrypto(otpCrypto *security.OTPCrypto) *AuthService {
	s.otpCrypto = otpCrypto
	return s
}

// WithSettingsReader injects a function to read platform settings
func (s *AuthService) WithSettingsReader(reader func(ctx context.Context, key string) string) *AuthService {
	s.settingsReader = reader
	return s
}

// Register creates a new seller or buyer account.
// Email is optional for sellers/buyers — a placeholder is generated if omitted.
// Email is required for admins (they are created via CreateAdmin, not this method).
func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*models.User, error) {
	if req.Name == "" || req.Phone == "" || req.Password == "" {
		return nil, errors.New("name, phone and password are required")
	}
	if valid, errKey := validations.ValidatePassword(req.Password); !valid {
		return nil, errors.New(validations.PasswordErrorMessage(errKey))
	}
	role := req.Role
	if role != "buyer" {
		role = "seller"
	}
	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	// Generate placeholder email only when needed (req.Email is empty).
	// Deferred here so retries on duplicate-phone don't waste crypto/rand entropy.
	email := req.Email
	if email == "" {
		suffix, err := randHex(4)
		if err != nil {
			return nil, fmt.Errorf("failed to generate email: %w", err)
		}
		email = fmt.Sprintf("%s_%s@noemail.local", role, suffix) //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
	}
	u := &models.User{
		Name:             req.Name,
		Phone:            req.Phone,
		BusinessName:     req.BusinessName,
		Email:            email,
		PasswordHash:     hash,
		Role:             role,
		IsActive:         true,
		SellerPreference: req.SellerPreference,
	}
	id, err := s.repo.Create(ctx, u)
	if err != nil {
		if isDuplicatePhone(err) {
			existing, lookupErr := s.repo.GetByPhoneAndRole(ctx, req.Phone, role)
			if lookupErr == nil && !existing.OTPVerified {
				// Update the stale unverified row with the latest submission values
				// so the password the user just typed will work after OTP verification.
				_ = s.repo.UpdateUnverifiedRegistration(ctx, req.Phone, role, hash, req.Name, email)
				safeGo(func() {
					s.logActivity(context.Background(), "register_retry", req.Phone, req.Email, role, "", "unverified re-register")
				})
				return nil, errors.New("otp_pending: phone registered but not yet verified — resend OTP")
			}
			safeGo(func() {
				s.logActivity(context.Background(), "register_fail", req.Phone, req.Email, role, "", "duplicate phone")
			})
			return nil, errors.New("phone number already registered — please sign in instead")
		}
		if isDuplicateEmail(err) {
			existing, lookupErr := s.repo.GetByPhoneAndRole(ctx, req.Phone, role)
			if lookupErr == nil && !existing.OTPVerified {
				// Same: update the unverified row with new values
				_ = s.repo.UpdateUnverifiedRegistration(ctx, req.Phone, role, hash, req.Name, email)
				safeGo(func() {
					s.logActivity(context.Background(), "register_retry", req.Phone, req.Email, role, "", "unverified re-register via email conflict")
				})
				return nil, errors.New("otp_pending: phone registered but not yet verified — resend OTP")
			}
			safeGo(func() {
				s.logActivity(context.Background(), "register_fail", req.Phone, req.Email, role, "", "duplicate email")
			})
			return nil, errors.New("email address already in use — please use a different email or sign in")
		}
		safeGo(func() {
			s.logActivity(context.Background(), "register_fail", req.Phone, req.Email, role, "", err.Error())
		})
		return nil, err
	}
	safeGo(func() { s.logActivity(context.Background(), "register_ok", req.Phone, req.Email, role, "", "") })
	// If OTP is disabled globally, mark the user as verified immediately.
	// Otherwise auto-send OTP so the user doesn't have to call /otp/send manually.
	if !s.IsOTPRequired(ctx) {
		_ = s.repo.SetOTPVerifiedForRole(ctx, req.Phone, role)
	} else {
		// Auto-send OTP immediately on registration — client should not need a separate call.
		// Use a fresh context to avoid inheriting any aborted-transaction state from the
		// registration context (pgx marks a context-associated tx as failed on any error).
		otpCtx, otpCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer otpCancel()
		otpErr := s.SendOTP(otpCtx, dto.OTPRequest{
			Phone:    req.Phone,
			Purpose:  "register",
			UserType: role,
		})
		if otpErr != nil {
			log.Printf("[AuthService] auto-send OTP failed for %s: %v", sanitizeLog(req.Phone), otpErr)
			// Delete the just-created unverified account so the user can retry cleanly.
			_ = s.repo.UpdateUnverifiedRegistration(ctx, req.Phone, role, "", req.Name, email)
			return nil, fmt.Errorf("registration succeeded but OTP delivery failed — please try again")
		}
	}
	// Link referral if a code was provided
	if req.ReferralCode != "" && s.referralSvc != nil {
		safeGo(func() {
			s.referralSvc.SetInvitedByOnRegister(context.Background(), id, req.ReferralCode)
		})
	}
	u.ID = id
	return u, nil
}

// Login authenticates by phone+password and returns a JWT.
// The caller must have already verified OTP (or pass requireOTP=false for admin).
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	return s.loginWithRole(ctx, req, "")
}

func (s *AuthService) LoginAsRole(ctx context.Context, req dto.LoginRequest, requiredRole string) (string, error) {
	return s.loginWithRole(ctx, req, requiredRole)
}

// SellerExists checks if a phone number belongs to a registered seller
func (s *AuthService) SellerExists(ctx context.Context, phone string) (bool, error) {
	u, err := s.repo.GetByPhoneAndRole(ctx, phone, "seller")
	if err != nil {
		return false, nil // not found = doesn't exist
	}
	return u.Role == "seller", nil
}

// BuyerExists checks if a phone number belongs to a registered buyer
func (s *AuthService) BuyerExists(ctx context.Context, phone string) (bool, error) {
	u, err := s.repo.GetByPhoneAndRole(ctx, phone, "buyer")
	if err != nil {
		return false, nil // not found = doesn't exist
	}
	return u.Role == "buyer", nil
}

func (s *AuthService) loginWithRole(ctx context.Context, req dto.LoginRequest, requiredRole string) (string, error) {
	if req.Password == "" {
		return "", errors.New("password is required")
	}

	var u *models.User
	var err error

	switch {
	case req.Email != "":
		// Email-based login — used by admins
		if req.Phone != "" {
			return "", errors.New("provide either email or phone, not both")
		}
		u, err = s.repo.GetByEmail(ctx, req.Email)
	case req.Phone != "":
		// Phone-based login — use role-scoped lookup when role is known
		// so the same phone can exist as both seller and buyer
		if requiredRole != "" {
			u, err = s.repo.GetByPhoneAndRole(ctx, req.Phone, requiredRole)
		} else {
			u, err = s.repo.GetByPhone(ctx, req.Phone)
		}
	default:
		return "", errors.New("email or phone is required")
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}
	if requiredRole != "" && u.Role != requiredRole {
		// Don't reveal whether the account exists under a different role
		return "", errors.New("invalid credentials")
	}
	if !u.IsActive {
		return "", errors.New("account is inactive — contact support")
	}
	if !security.CheckPassword(u.PasswordHash, req.Password) {
		safeGo(func() {
			s.logActivity(context.Background(), "login_fail", req.Phone, req.Email, u.Role, "", "wrong password")
		})
		return "", errors.New("invalid credentials")
	}
	// Block login if OTP verification is required and the phone has not been verified yet.
	// This prevents users who registered but never completed OTP from accessing the dashboard.
	if s.IsOTPRequired(ctx) && !u.OTPVerified {
		return "", errors.New("otp_pending: phone not verified — please complete OTP verification first")
	}
	token, err := security.GenerateJWTFull(u.ID, u.Role, u.AdminRole, u.BuyerTier, u.SellerPreference, s.jwtSecret)
	if err != nil {
		return "", err
	}
	safeGo(func() { s.logActivity(context.Background(), "login_ok", req.Phone, req.Email, u.Role, "", "") })
	safeGo(func() { s.recordLogin(u) })
	return token, nil
}

// ── OTP ───────────────────────────────────────────────────────────────────────

// SendOTP generates an OTP for the phone and returns the code (for SMS/email delivery).
func (s *AuthService) SendOTP(ctx context.Context, req dto.OTPRequest) error {
	if req.Phone == "" {
		return errors.New("phone is required")
	}
	// Allowlist purpose — prevents arbitrary strings being stored in otp_codes
	switch req.Purpose {
	case "register", "login", "reset_password":
		// valid
	default:
		return errors.New("invalid purpose: must be register, login, or reset_password")
	}
	// Allowlist user_type
	userType := req.UserType
	switch userType {
	case "seller", "buyer", "admin":
		// valid
	case "":
		userType = "seller" // safe default
	default:
		return errors.New("invalid user_type: must be seller, buyer, or admin")
	}
	code, err := s.repo.GenerateOTPWithType(ctx, req.Phone, req.Purpose, userType)
	if err != nil {
		safeGo(func() { s.logActivity(context.Background(), "otp_send_fail", req.Phone, "", userType, "", err.Error()) })
		return err
	}
	safeGo(func() { s.logActivity(context.Background(), "otp_send", req.Phone, "", userType, "", req.Purpose) })

	// Load dynamic OTP gateway setting
	gateway, _ := s.repo.GetSetting(ctx, "otp_gateway")
	if gateway == "" {
		gateway = "log"
	}

	if s.devMode {
		log.Printf("🔑 [DEV OTP] phone=%s purpose=%s code=%s (gateway: %s)", sanitizeLog(req.Phone), sanitizeLog(req.Purpose), code, gateway)
	}

	switch gateway {
	case "sms":
		// TODO: Implement real SMS gateway (Twilio, Infobip, etc.)
		log.Printf("[SMS GATEWAY] Sending OTP to %s", sanitizeLog(req.Phone))
	case "email":
		// If the user has an email, we can send it there too
		u, err := s.repo.GetByPhoneAndRole(ctx, req.Phone, userType)
		if err == nil && u.Email != "" && s.emailSender != nil {
			_ = s.emailSender.SendOTP(u.Email, code, req.Purpose)
		}
	default:
		// Default to 'log' — only visible in server logs
		if !s.devMode {
			log.Printf("🔑 [OTP LOG] phone=%s purpose=%s code=****** (check server logs in dev mode)", sanitizeLog(req.Phone), sanitizeLog(req.Purpose))
		}
	}

	return nil
}

// SendEmailOTP generates an OTP for email-based users (admins).
// In dev mode the code is logged to stdout for easy retrieval.
// In production the code is masked — admins must use the injected EmailSender.
func (s *AuthService) SendEmailOTP(ctx context.Context, email, purpose string) error {
	if email == "" {
		return errors.New("email is required")
	}
	code, err := s.repo.GenerateEmailOTP(ctx, email, purpose, "admin")
	if err != nil {
		return err
	}
	// Load dynamic OTP gateway setting
	gateway, _ := s.repo.GetSetting(ctx, "otp_gateway")
	if gateway == "" {
		gateway = "log"
	}

	if s.devMode {
		log.Printf("🔑 [ADMIN OTP] email=%s purpose=%s code=%s (gateway: %s)", sanitizeLog(email), sanitizeLog(purpose), code, gateway)
	} else {
		// Production logging logic
		if gateway == "log" || s.emailSender == nil {
			// FALLBACK: if gateway is 'log' OR if we simply have no email sender configured,
			// we MUST log the code so the admin isn't locked out.
			log.Printf("🔑 [ADMIN OTP LOG] email=%s purpose=%s code=****** (no email sender configured — use dev mode to see code)",
				sanitizeLog(email), sanitizeLog(purpose))
		} else {
			// Masked for production security when email is working
			log.Printf("🔑 [ADMIN OTP] email=%s purpose=%s code=****** (check email)", sanitizeLog(email), sanitizeLog(purpose))
		}
	}
	if s.emailSender != nil && gateway == "email" {
		if err := s.emailSender.SendOTP(email, code, purpose); err != nil {
			log.Printf("[WARN] email OTP delivery failed for %s: %v", sanitizeLog(email), err) //nolint:gosec // G117: error/value from internal source, not user-controlled input
		}
	}
	return nil
}

// VerifyOTP checks the OTP code and marks the phone as verified.
func (s *AuthService) VerifyOTP(ctx context.Context, req dto.OTPVerifyRequest) error {
	if req.Phone == "" || req.Code == "" {
		return errors.New("phone and code are required")
	}
	// Enforce basic code format — must be 4–8 digits only.
	// Rejects empty strings, letters, and obviously wrong lengths before hitting the DB.
	if len(req.Code) < 4 || len(req.Code) > 8 {
		return errors.New("invalid OTP format")
	}
	for _, ch := range req.Code {
		if ch < '0' || ch > '9' {
			return errors.New("invalid OTP format")
		}
	}
	// Enforce purpose allowlist — prevents verifying a register OTP as a login OTP
	switch req.Purpose {
	case "register", "login", "reset_password":
		// valid
	default:
		return errors.New("invalid purpose")
	}
	// Resolve user_type from the OTP row so the activity log is accurate
	userType := s.repo.GetOTPUserType(ctx, req.Phone)
	if err := s.repo.VerifyOTP(ctx, req.Phone, req.Code, req.Purpose); err != nil {
		safeGo(func() {
			s.logActivity(context.Background(), "otp_verify_fail", req.Phone, "", userType, "", err.Error())
		})
		return err
	}
	safeGo(func() { s.logActivity(context.Background(), "otp_verify_ok", req.Phone, "", userType, "", req.Purpose) })
	// Mark otp_verified on the correct role row
	return s.repo.SetOTPVerifiedForRole(ctx, req.Phone, userType)
}

// VerifyEmailOTP checks an email-based OTP (for admins).
func (s *AuthService) VerifyEmailOTP(ctx context.Context, email, code, purpose string) error {
	if email == "" || code == "" {
		return errors.New("email and code are required")
	}
	return s.repo.VerifyEmailOTP(ctx, email, code, purpose)
}

// LoginWithOTP verifies OTP then issues a JWT — full OTP login flow.
// role must be "seller" or "buyer"; defaults to "seller" if empty.
func (s *AuthService) LoginWithOTP(ctx context.Context, phone, code, role string) (string, error) {
	if role != "buyer" {
		role = "seller"
	}
	if err := s.repo.VerifyOTP(ctx, phone, code, "login"); err != nil {
		return "", err
	}
	u, err := s.repo.GetByPhoneAndRole(ctx, phone, role)
	if err != nil {
		return "", errors.New("account not found")
	}
	if !u.IsActive {
		return "", errors.New("account is inactive — contact support")
	}
	// Mark phone as verified — OTP login proves ownership
	_ = s.repo.SetOTPVerifiedForRole(ctx, phone, role)
	token, err := security.GenerateJWTFull(u.ID, u.Role, u.AdminRole, u.BuyerTier, u.SellerPreference, s.jwtSecret)
	if err != nil {
		return "", err
	}
	safeGo(func() { s.recordLogin(u) })
	return token, nil
}

// AdminLoginStep1 validates email+password.
// Super admins get a JWT immediately (no OTP).
// Normal (sub) admins get an OTP sent to their email — JWT issued only after step 2.
func (s *AuthService) AdminLoginStep1(ctx context.Context, req dto.LoginRequest) (string, bool, error) {
	if req.Email == "" || req.Password == "" {
		return "", false, errors.New("email and password are required")
	}
	u, err := s.repo.GetAdminByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, errors.New("invalid credentials")
		}
		return "", false, err
	}
	if u.Role != "admin" {
		return "", false, errors.New("invalid credentials")
	}
	if !u.IsActive {
		return "", false, errors.New("account is inactive — contact support")
	}
	if !security.CheckPassword(u.PasswordHash, req.Password) {
		return "", false, errors.New("invalid credentials")
	}

	// Super admins always bypass OTP — they are trusted accounts created by the system.
	// Normal sub-admins always require OTP regardless of the global otp_required setting.
	if u.AdminRole == "super" {
		token, err := security.GenerateJWTFull(u.ID, u.Role, u.AdminRole, u.BuyerTier, u.SellerPreference, s.jwtSecret)
		if err != nil {
			return "", false, err
		}
		safeGo(func() { s.recordLogin(u) })
		return token, true, nil // true = OTP bypassed (super admin only)
	}

	// All other cases (normal sub-admin OR super admin when OTP is required)
	// send OTP to email, return email for step 2.
	if err := s.SendEmailOTP(ctx, u.Email, "login"); err != nil {
		return "", false, fmt.Errorf("failed to send OTP: %w", err)
	}
	return u.Email, false, nil // false = OTP required, email returned for step 2
}

// AdminLoginStep2 verifies the email OTP for sub-admins and issues a JWT.
func (s *AuthService) AdminLoginStep2(ctx context.Context, email, code string) (string, error) {
	if err := s.repo.VerifyEmailOTP(ctx, email, code, "login"); err != nil {
		return "", err
	}
	u, err := s.repo.GetAdminByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}
	if !u.IsActive {
		return "", errors.New("account is inactive — contact support")
	}
	token, err := security.GenerateJWTFull(u.ID, u.Role, u.AdminRole, u.BuyerTier, u.SellerPreference, s.jwtSecret)
	if err != nil {
		return "", err
	}
	safeGo(func() { s.recordLogin(u) })
	return token, nil
}

// GetOTPLogs returns paginated OTP log entries for the admin panel.
func (s *AuthService) GetOTPLogs(ctx context.Context, f repository.OTPLogFilter) ([]repository.OTPLogEntry, int, error) {
	return s.repo.ListOTPLogs(ctx, f)
}

// GetActivityLogs returns paginated user activity logs for the analytics panel.
func (s *AuthService) GetActivityLogs(ctx context.Context, f repository.ActivityLogFilter) ([]repository.ActivityLogEntry, int, error) {
	return s.repo.ListActivityLogs(ctx, f)
}

// safeGo launches fn in a new goroutine with panic recovery.
// Callers use safeGo(func(){...}) as a shorthand for go+recover.
//
//nolint:govet // hidden-goroutine: intentional fire-and-forget pattern with explicit panic recovery
func safeGo(fn func()) { go withRecover(fn) }

// withRecover wraps fn with a deferred panic handler.
func withRecover(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[PANIC recovered in goroutine] %v", r)
		}
	}()
	fn()
}

// logActivity is a fire-and-forget helper used by all auth flows.
func (s *AuthService) logActivity(_ context.Context, event, phone, email, userType, ip, meta string) {
	// Always use a fresh background context so the log isn't cancelled when the
	// request context ends (fire-and-forget goroutines pass the request ctx which
	// gets cancelled immediately after the handler returns).
	logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.repo.LogActivity(logCtx, event, phone, email, userType, ip, meta); err != nil {
		log.Printf("[WARN] logActivity event=%s phone=%s: %v", event, phone, err) //nolint:gosec // G117: error/value from internal source, not user-controlled input
	}
}

// DeleteOTPLogs removes OTP records to free storage.
// status="done"  → deletes used + expired rows
// status="all"   → deletes everything
func (s *AuthService) DeleteOTPLogs(ctx context.Context, status string) (int64, error) {
	return s.repo.DeleteOTPLogs(ctx, status)
}

// ── Forgot / Reset Password ───────────────────────────────────────────────────

// ForgotPasswordSend sends an OTP to the phone for password reset (sellers & buyers).
func (s *AuthService) ForgotPasswordSend(ctx context.Context, phone, role string) error {
	if phone == "" {
		return errors.New("phone is required")
	}
	if role != "seller" && role != "buyer" {
		role = "seller"
	}
	// Verify the account exists before sending OTP — don't reveal existence via error
	_, err := s.repo.GetByPhoneAndRole(ctx, phone, role)
	if err != nil {
		// Return success regardless to prevent account enumeration
		return nil
	}
	code, err := s.repo.GenerateOTPWithType(ctx, phone, "reset_password", role)
	if err != nil {
		return err
	}
	if s.devMode {
		log.Printf("🔑 [DEV OTP] phone=%s purpose=reset_password role=%s code=%s", sanitizeLog(phone), sanitizeLog(role), code) //nolint:gosec // G117: error/value from internal source, not user-controlled input
	}
	return nil
}

// ForgotPasswordReset verifies the OTP and sets the new password (sellers & buyers).
func (s *AuthService) ForgotPasswordReset(ctx context.Context, phone, role, code, newPassword string) error {
	if phone == "" || code == "" || newPassword == "" {
		return errors.New("phone, code, and new_password are required")
	}
	if valid, errKey := validations.ValidatePassword(newPassword); !valid {
		return errors.New(validations.PasswordErrorMessage(errKey))
	}
	if role != "seller" && role != "buyer" {
		role = "seller"
	}
	if err := s.repo.VerifyOTPForRole(ctx, phone, code, "reset_password", role); err != nil {
		return err
	}
	u, err := s.repo.GetByPhoneAndRole(ctx, phone, role)
	if err != nil {
		return errors.New("account not found")
	}
	hash, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := s.repo.UpdatePassword(ctx, u.ID, hash); err != nil {
		return err
	}
	// Proving phone ownership via reset OTP also counts as OTP verification.
	// This unblocks users who registered but never completed OTP.
	_ = s.repo.SetOTPVerifiedForRole(ctx, phone, role)
	return nil
}

// AdminForgotPasswordSend sends an OTP to the admin's email for password reset.
func (s *AuthService) AdminForgotPasswordSend(ctx context.Context, email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	// Verify account exists silently — don't reveal existence
	_, err := s.repo.GetAdminByEmail(ctx, email)
	if err != nil {
		return nil // silent success to prevent enumeration
	}
	code, err := s.repo.GenerateEmailOTP(ctx, email, "reset_password", "admin")
	if err != nil {
		return err
	}
	log.Printf("🔑 [ADMIN OTP] email=%s purpose=reset_password code=****** (check email)", sanitizeLog(email)) //nolint:gosec // G117: error/value from internal source, not user-controlled input
	if s.devMode {
		log.Printf("🔑 [ADMIN OTP DEV] email=%s purpose=reset_password code=%s", sanitizeLog(email), code) //nolint:gosec // G117: error/value from internal source, not user-controlled input
	}
	if s.emailSender != nil {
		if err := s.emailSender.SendOTP(email, code, "reset_password"); err != nil {
			log.Printf("[WARN] email OTP delivery failed for %s: %v", sanitizeLog(email), err) //nolint:gosec // G117: error/value from internal source, not user-controlled input
		}
	}
	return nil
}

// AdminForgotPasswordReset verifies the OTP and sets the new admin password.
func (s *AuthService) AdminForgotPasswordReset(ctx context.Context, email, code, newPassword string) error {
	if email == "" || code == "" || newPassword == "" {
		return errors.New("email, code, and new_password are required")
	}
	if valid, errKey := validations.ValidatePassword(newPassword); !valid {
		return errors.New(validations.PasswordErrorMessage(errKey))
	}
	if err := s.repo.VerifyEmailOTP(ctx, email, code, "reset_password"); err != nil {
		return err
	}
	u, err := s.repo.GetAdminByEmail(ctx, email)
	if err != nil {
		return errors.New("account not found")
	}
	hash, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, u.ID, hash)
}

func (s *AuthService) ChangePassword(ctx context.Context, userID int, currentPw, newPw string) error {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if !security.CheckPassword(u.PasswordHash, currentPw) {
		return errors.New("current password is incorrect")
	}
	hash, err := security.HashPassword(newPw)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, userID, hash)
}

// ── Buyer tier ────────────────────────────────────────────────────────────────

// UpgradeBuyerTier submits a pending upgrade request — tier is NOT applied until admin approves.
// Upgrade prices are read from platform_settings (buyer_gold_upgrade_price / buyer_premium_upgrade_price)
// so the admin can change them without a code deploy.
func (s *AuthService) UpgradeBuyerTier(ctx context.Context, userID int, req dto.BuyerUpgradeRequest) error {
	if req.Tier != "gold" && req.Tier != "premium" {
		return errors.New("invalid tier: must be gold or premium")
	}
	if req.Reference == "" {
		return errors.New("payment reference is required")
	}
	if req.PayerPhone == "" {
		return errors.New("payer_phone is required")
	}
	validMethods := map[string]bool{"telebirr": true, "safaricom": true, "cbe": true}
	if !validMethods[req.PaymentMethod] {
		return errors.New("payment_method must be telebirr, safaricom, or cbe")
	}

	// Prevent downgrade or same-tier request
	u, lookupErr := s.repo.GetByID(ctx, userID)
	if lookupErr != nil {
		return errors.New("user not found")
	}
	if TierRank[req.Tier] <= TierRank[u.BuyerTier] {
		return errors.New("you are already at this tier or higher — you can only upgrade to a higher tier")
	}

	// Load price from platform_settings — admin-controlled, no hardcoded amounts.
	settingKey := "buyer_gold_upgrade_price"
	fallback := repository.BuyerGoldCost
	if req.Tier == "premium" {
		settingKey = "buyer_premium_upgrade_price"
		fallback = repository.BuyerPremiumCost
	}
	cost := fallback
	if v, err := s.repo.GetSetting(ctx, settingKey); err == nil && v != "" {
		if f, err2 := strconv.ParseFloat(v, 64); err2 == nil && f > 0 {
			cost = f
		}
	}

	// Scale cost based on duration if duration is provided. Base cost is assumed to be for 1 Month (30 days).
	if req.DurationDays > 0 {
		// Enforce sane duration bounds — prevent paying pennies for 10,000 days
		if req.DurationDays > 365 {
			return fmt.Errorf("duration_days cannot exceed 365 days (1 year)")
		}
		if req.DurationDays < 7 {
			return fmt.Errorf("duration_days must be at least 7 days")
		}
		// round to nearest integer using math.Round to avoid float truncation issues
		cost = math.Round(cost * float64(req.DurationDays) / 30.0)
	} else {
		// Default to 1 month if not specified
		req.DurationDays = 30
	}

	// Validate that the submitted amount matches the configured price.
	// This prevents a buyer from submitting amount=1 for a 10,000 ETB upgrade.
	if req.Amount <= 0 {
		return errors.New("amount is required")
	}

	// Read currency tolerance as PERCENTAGE from DB settings (admin-configurable)
	tolerancePercent := 0.001 // 0.1% default
	if toleranceStr := s.settingsReader(ctx, constants.SettingCurrencyTolerancePercent); toleranceStr != "" {
		if f, err := strconv.ParseFloat(toleranceStr, 64); err == nil && f > 0 && f < 100 {
			tolerancePercent = f / 100.0 // Convert percentage to decimal
		}
	}

	// Calculate tolerance as percentage of cost (prevents manipulation)
	maxAllowedDifference := cost * tolerancePercent

	// Accept amounts within tolerance of cost to handle rounding on the client side.
	// Reject amounts more than tolerance below cost (underpayment) but also
	// cap overpayment at tolerance to avoid silent fund loss.
	if req.Amount < cost-maxAllowedDifference {
		return errors.New("payment amount is less than the required amount for this tier and duration")
	}
	if req.Amount > cost+maxAllowedDifference {
		return errors.New("payment amount exceeds the required amount for this tier and duration — please submit the exact amount")
	}

	_, err := s.repo.SubmitBuyerUpgradeRequest(ctx, userID, req.Tier, cost, req.Reference, req.PayerPhone, req.PaymentMethod, req.DurationDays)
	return err
}

// FindOrCreateSeller finds an existing seller by phone or creates a new one.
// Uses INSERT ... ON CONFLICT (phone, role) to handle concurrent public-sell requests safely.
func (s *AuthService) FindOrCreateSeller(ctx context.Context, phone, name string) (int, error) {
	// Try to find first — avoids generating a random password unnecessarily.
	u, err := s.repo.GetByPhoneAndRole(ctx, phone, "seller")
	if err == nil {
		return u.ID, nil
	}
	hash, err := security.HashPassword(security.RandomPassword())
	if err != nil {
		return 0, err
	}
	// Use a random suffix so the placeholder email is unique even if the same
	// phone is used in a race condition or if someone manually registered with
	// the phone-derived email pattern.
	suffix, err := randHex(4)
	if err != nil {
		return 0, fmt.Errorf("failed to generate email: %w", err)
	}
	email := fmt.Sprintf("seller_%s@noemail.local", suffix) //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
	newUser := &models.User{
		Name:         name,
		Phone:        phone,
		Email:        email,
		PasswordHash: hash,
		Role:         "seller",
		IsActive:     true,
	}
	id, err := s.repo.CreateOrGet(ctx, newUser)
	if err != nil {
		return 0, err
	}
	// System-created accounts bypass OTP — mark verified so they can login
	_ = s.repo.SetOTPVerifiedForRole(ctx, phone, "seller")
	return id, nil
}

// GetProfile returns a user's profile by ID.
func (s *AuthService) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	return s.repo.GetByID(ctx, userID)
}

// GetBuyerOrders returns order requests for a buyer by their phone number.
func (s *AuthService) GetBuyerOrders(ctx context.Context, userID int) ([]map[string]interface{}, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListOrdersByBuyerPhone(ctx, u.Phone)
}

// DevPeekOTP returns the plaintext OTP for a phone — DEV ONLY.
// Decrypts the latest active OTP from the DB. Never call in production.
func (s *AuthService) DevPeekOTP(ctx context.Context, phone string) (string, error) {
	return s.repo.DevPeekOTP(ctx, phone)
}

// DevDecryptOTP decrypts any OTP ciphertext hex string — DEV ONLY.
func (s *AuthService) DevDecryptOTP(hexCode string) (string, error) {
	if s.otpCrypto == nil {
		return "", fmt.Errorf("OTP crypto not configured")
	}
	return s.otpCrypto.DecryptOTP(hexCode)
}

// recordLogin marks the user's first login and fires a super-admin notification
// if this is the first time they've logged in. Always uses its own context.
func (s *AuthService) recordLogin(u *models.User) {
	// Use a fresh context — the caller's request context may already be cancelled.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	isFirst, err := s.repo.RecordFirstLogin(ctx, u.ID)
	if err != nil {
		log.Printf("[WARN] recordLogin user=%d: %v", u.ID, err) //nolint:gosec // G117: error/value from internal source, not user-controlled input
		return
	}
	if !isFirst {
		return
	}
	if s.adminSvc != nil {
		notifyCtx, notifyCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer notifyCancel()
		s.adminSvc.NotifyFirstLogin(notifyCtx, u.Name, u.Role)
	}
}

// IsOTPRequired reads the otp_required platform setting with a 60s in-memory cache
// to avoid a DB round-trip on every login/register call.
// Uses a single Lock() for the full check-and-refresh to prevent the
// double-checked locking race between RUnlock and Lock.
func (s *AuthService) IsOTPRequired(ctx context.Context) bool {
	s.otpRequiredMu.Lock()
	defer s.otpRequiredMu.Unlock()
	if time.Since(s.otpRequiredCachedAt) < 60*time.Second {
		return s.otpRequiredCache
	}
	val, err := s.repo.GetSetting(ctx, "otp_required")
	result := true // safe default
	if err == nil && val != "" {
		result = val == "true"
	}
	s.otpRequiredCache = result
	s.otpRequiredCachedAt = time.Now()
	return result
}

// IsRegistrationAllowed checks the platform setting for the given role.
// Sellers check allow_normal_registration (all new sellers start as Normal).
// Buyers check allow_buyer_registration.
// Returns true by default when the setting is missing — open by default.
func (s *AuthService) IsRegistrationAllowed(ctx context.Context, role string) bool {
	key := "allow_buyer_registration"
	if role == "seller" {
		key = "allow_normal_registration"
	}
	val, err := s.repo.GetSetting(ctx, key)
	if err != nil || val == "" {
		return true // open by default when setting not configured
	}
	return val == "true"
}

// isDuplicatePhone returns true when the error is a Postgres unique-constraint
// violation on the users phone+role constraint (SQLSTATE 23505).
func isDuplicatePhone(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	// Matches both old (users_phone_key) and new (users_phone_role_key) constraint names
	return (strings.Contains(msg, "users_phone_key") || strings.Contains(msg, "users_phone_role_key")) &&
		(strings.Contains(msg, "23505") || strings.Contains(msg, "duplicate key"))
}

// isDuplicateEmail returns true when the error is a Postgres unique-constraint
// violation on the users email column (SQLSTATE 23505).
func isDuplicateEmail(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "users_email_key") &&
		(strings.Contains(msg, "23505") || strings.Contains(msg, "duplicate key"))
}

// ── Buyer profile / notifications ─────────────────────────────────────────────

// UpdateBuyerProfile updates a buyer's name and email.
func (s *AuthService) UpdateBuyerProfile(ctx context.Context, userID int, name, email string) error {
	return s.repo.UpdateBuyerProfile(ctx, userID, name, email)
}

// GetBuyerUpgradeRequestsByUser returns upgrade requests for a specific buyer.
func (s *AuthService) GetBuyerUpgradeRequestsByUser(ctx context.Context, userID int) ([]map[string]interface{}, error) {
	return s.repo.ListBuyerUpgradeRequestsByUser(ctx, userID)
}

// GetBuyerNotifications returns notifications for a buyer.
func (s *AuthService) GetBuyerNotifications(ctx context.Context, userID int) ([]map[string]interface{}, error) {
	return s.repo.ListBuyerNotifications(ctx, userID)
}

// MarkBuyerNotificationRead marks a single notification as read for a buyer.
func (s *AuthService) MarkBuyerNotificationRead(ctx context.Context, notifID, userID int) error {
	return s.repo.MarkBuyerNotificationRead(ctx, notifID, userID)
}

// SetDemand updates the user's current demand (FIND or LIST) after registration
func (s *AuthService) SetDemand(ctx context.Context, userID int, demand string) error {
	return s.repo.SetDemand(ctx, userID, demand)
}
