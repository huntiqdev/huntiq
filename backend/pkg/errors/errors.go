// Package errors provides structured error handling with error codes and multilingual support.
package errors

import (
	"fmt"
)

// ErrorCode represents a unique error identifier that can be mapped to translations
type ErrorCode string

// Error categories
const (
	// Authentication errors (AUTH_xxx)
	ErrAuthInvalidCredentials     ErrorCode = "AUTH_INVALID_CREDENTIALS"
	ErrAuthWrongPassword          ErrorCode = "AUTH_WRONG_PASSWORD"
	ErrAuthAccountInactive        ErrorCode = "AUTH_ACCOUNT_INACTIVE"
	ErrAuthPhoneNotVerified       ErrorCode = "AUTH_PHONE_NOT_VERIFIED"
	ErrAuthSessionExpired         ErrorCode = "AUTH_SESSION_EXPIRED"
	ErrAuthTokenInvalid           ErrorCode = "AUTH_TOKEN_INVALID"
	ErrAuthRegistrationClosed     ErrorCode = "AUTH_REGISTRATION_CLOSED"
	ErrAuthBuyerRegistrationClosed ErrorCode = "AUTH_BUYER_REGISTRATION_CLOSED"
	ErrAuthPhoneAlreadyRegistered ErrorCode = "AUTH_PHONE_ALREADY_REGISTERED"
	ErrAuthEmailAlreadyInUse      ErrorCode = "AUTH_EMAIL_ALREADY_IN_USE"
	ErrAuthAccountNotFound        ErrorCode = "AUTH_ACCOUNT_NOT_FOUND"

	// OTP errors (OTP_xxx)
	ErrOTPLocked          ErrorCode = "OTP_LOCKED"
	ErrOTPExpired         ErrorCode = "OTP_EXPIRED"
	ErrOTPRequired        ErrorCode = "OTP_REQUIRED"
	ErrOTPSendFailed      ErrorCode = "OTP_SEND_FAILED"
	ErrOTPInvalid         ErrorCode = "OTP_INVALID"
	ErrOTPNoActive        ErrorCode = "OTP_NO_ACTIVE"
	ErrOTPRateLimitOTP    ErrorCode = "OTP_RATE_LIMIT"
	ErrOTPRateLimitGeneral ErrorCode = "RATE_LIMIT_GENERAL"

	// Permission errors (PERM_xxx)
	ErrPermSuperAdminRequired     ErrorCode = "PERM_SUPER_ADMIN_REQUIRED"
	ErrPermAdminRequired          ErrorCode = "PERM_ADMIN_REQUIRED"
	ErrPermDenied                 ErrorCode = "PERM_DENIED"
	ErrPermSellerTierRequired     ErrorCode = "PERM_SELLER_TIER_REQUIRED"
	ErrPermCannotDeleteSelf       ErrorCode = "PERM_CANNOT_DELETE_SELF"
	ErrPermCannotChangeOwnRole    ErrorCode = "PERM_CANNOT_CHANGE_OWN_ROLE"
	ErrPermLastSuperAdmin         ErrorCode = "PERM_LAST_SUPER_ADMIN"
	ErrPermSuperAdminDelete       ErrorCode = "PERM_SUPER_ADMIN_DELETE"
	ErrPermReferralBlocked        ErrorCode = "PERM_REFERRAL_BLOCKED"
	ErrPermReferralSuspended      ErrorCode = "PERM_REFERRAL_SUSPENDED"
	ErrPermOwnReferralCode        ErrorCode = "PERM_OWN_REFERRAL_CODE"
	ErrPermInvitationAlreadyAccepted ErrorCode = "PERM_INVITATION_ALREADY_ACCEPTED"

	// Upload errors (UPLOAD_xxx)
	ErrUploadCloudFailed      ErrorCode = "UPLOAD_CLOUD_FAILED"
	ErrUploadFailed           ErrorCode = "UPLOAD_FAILED"
	ErrUploadFileTooLarge     ErrorCode = "UPLOAD_FILE_TOO_LARGE"
	ErrUploadWrongType        ErrorCode = "UPLOAD_WRONG_TYPE"
	ErrUploadURLNotLoadable   ErrorCode = "UPLOAD_URL_NOT_LOADABLE"
	ErrUploadLimitReached     ErrorCode = "UPLOAD_LIMIT_REACHED"

	// Payment errors (PAY_xxx)
	ErrPayPendingExists           ErrorCode = "PAY_PENDING_EXISTS"
	ErrPayAlreadyPaid             ErrorCode = "PAY_ALREADY_PAID"
	ErrPaySoldProduct             ErrorCode = "PAY_SOLD_PRODUCT"
	ErrPayUpgradeClosed           ErrorCode = "PAY_UPGRADE_CLOSED"
	ErrPayInsufficientBalance     ErrorCode = "PAY_INSUFFICIENT_BALANCE"
	ErrPayPendingWithdrawal       ErrorCode = "PAY_PENDING_WITHDRAWAL"
	ErrPayMinWithdrawal           ErrorCode = "PAY_MIN_WITHDRAWAL"
	ErrPayMaxWithdrawal           ErrorCode = "PAY_MAX_WITHDRAWAL"
	ErrPayInvalidMethod           ErrorCode = "PAY_INVALID_METHOD"
	ErrPayRejectionReasonRequired ErrorCode = "PAY_REJECTION_REASON_REQUIRED"
	ErrPayReferenceRequired       ErrorCode = "PAY_REFERENCE_REQUIRED"

	// Product errors (PROD_xxx)
	ErrProdNotAvailable    ErrorCode = "PROD_NOT_AVAILABLE"
	ErrProdNotFound        ErrorCode = "PROD_NOT_FOUND"
	ErrProdListingExpired  ErrorCode = "PROD_LISTING_EXPIRED"
	ErrProdPriceMismatch   ErrorCode = "PROD_PRICE_MISMATCH"
	ErrProdCannotEdit      ErrorCode = "PROD_CANNOT_EDIT"
	ErrProdCannotDelete    ErrorCode = "PROD_CANNOT_DELETE"
	ErrProdMinImages       ErrorCode = "PROD_MIN_IMAGES"

	// Validation errors (VAL_xxx)
	ErrValPhoneRequired        ErrorCode = "VAL_PHONE_REQUIRED"
	ErrValPhoneInvalid         ErrorCode = "VAL_PHONE_INVALID"
	ErrValPasswordMin          ErrorCode = "VAL_PASSWORD_MIN"
	ErrValNameRequired         ErrorCode = "VAL_NAME_REQUIRED"
	ErrValTitleRequired        ErrorCode = "VAL_TITLE_REQUIRED"
	ErrValTitleTooLong         ErrorCode = "VAL_TITLE_TOO_LONG"
	ErrValDescriptionTooLong   ErrorCode = "VAL_DESCRIPTION_TOO_LONG"
	ErrValLocationRequired     ErrorCode = "VAL_LOCATION_REQUIRED"
	ErrValLocationInvalid      ErrorCode = "VAL_LOCATION_INVALID"
	ErrValTierInvalid          ErrorCode = "VAL_TIER_INVALID"
	ErrValAmountPositive       ErrorCode = "VAL_AMOUNT_POSITIVE"
	ErrValPriceRequired        ErrorCode = "VAL_PRICE_REQUIRED"
	ErrValRequired             ErrorCode = "VAL_REQUIRED"
	ErrValSettingKeyUnknown    ErrorCode = "VAL_SETTING_KEY_UNKNOWN"
	ErrValSettingNumeric       ErrorCode = "VAL_SETTING_NUMERIC"
	ErrValSettingBoolean       ErrorCode = "VAL_SETTING_BOOLEAN"
	ErrValDateFormat           ErrorCode = "VAL_DATE_FORMAT"
	ErrValPayerPhoneRequired   ErrorCode = "VAL_PAYER_PHONE_REQUIRED"
	ErrValMaxImages            ErrorCode = "VAL_MAX_IMAGES"
	ErrValInvalidImageURL      ErrorCode = "VAL_INVALID_IMAGE_URL"
	ErrValPasswordsMismatch    ErrorCode = "VAL_PASSWORDS_MISMATCH"

	// Not found errors (NF_xxx)
	ErrNFUser              ErrorCode = "NF_USER"
	ErrNFReferralCode      ErrorCode = "NF_REFERRAL_CODE"
	ErrNFUpgradeRequest    ErrorCode = "NF_UPGRADE_REQUEST"
	ErrNFWithdrawal        ErrorCode = "NF_WITHDRAWAL"
	ErrNFTransaction       ErrorCode = "NF_TRANSACTION"
	ErrNFNotification      ErrorCode = "NF_NOTIFICATION"
	ErrNFCategory          ErrorCode = "NF_CATEGORY"
	ErrNFAdmin             ErrorCode = "NF_ADMIN"

	// Server errors (SRV_xxx)
	ErrSrvInternal ErrorCode = "SRV_INTERNAL"
	ErrSrvDatabase ErrorCode = "SRV_DATABASE"

	// Network errors (NET_xxx)
	ErrNetOffline ErrorCode = "NET_OFFLINE"
	ErrNetTimeout ErrorCode = "NET_TIMEOUT"
)

// AppError represents a structured application error with code and context
type AppError struct {
	Code       ErrorCode              // Machine-readable error code
	Message    string                 // Human-readable English message (fallback)
	HTTPStatus int                    // HTTP status code
	Context    map[string]interface{} // Additional context for interpolation
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return string(e.Code)
}

// New creates a new AppError with the given code and message
func New(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Context:    make(map[string]interface{}),
	}
}

// WithContext adds context data for error message interpolation
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	e.Context[key] = value
	return e
}

// Common error constructors for convenience
func AuthInvalidCredentials() *AppError {
	return New(ErrAuthInvalidCredentials, "invalid credentials", 401)
}

func AuthWrongPassword() *AppError {
	return New(ErrAuthWrongPassword, "current password is incorrect", 401)
}

func AuthPhoneNotVerified() *AppError {
	return New(ErrAuthPhoneNotVerified, "phone not verified — complete OTP verification first", 403)
}

func AuthSessionExpired() *AppError {
	return New(ErrAuthSessionExpired, "session expired — please log in again", 401)
}

func AuthTokenInvalid() *AppError {
	return New(ErrAuthTokenInvalid, "invalid or missing authorization token", 401)
}

func AuthRegistrationClosed() *AppError {
	return New(ErrAuthRegistrationClosed, "seller registration is currently closed — contact admin", 403)
}

func AuthPhoneAlreadyRegistered() *AppError {
	return New(ErrAuthPhoneAlreadyRegistered, "phone number already registered", 400)
}

func OTPLocked(minutes int) *AppError {
	return New(ErrOTPLocked, fmt.Sprintf("too many failed attempts — locked for %d minutes", minutes), 429).
		WithContext("minutes", minutes)
}

func OTPExpired() *AppError {
	return New(ErrOTPExpired, "invalid or expired OTP code", 400)
}

func OTPInvalid(attemptsRemaining int) *AppError {
	return New(ErrOTPInvalid, fmt.Sprintf("invalid OTP — %d attempt(s) remaining", attemptsRemaining), 400).
		WithContext("attempts", attemptsRemaining)
}

func OTPRequired() *AppError {
	return New(ErrOTPRequired, "OTP code is required", 400)
}

func OTPSendFailed() *AppError {
	return New(ErrOTPSendFailed, "failed to send OTP — please try again", 500)
}

func PermDenied() *AppError {
	return New(ErrPermDenied, "you do not have permission to perform this action", 403)
}

func PermSuperAdminRequired() *AppError {
	return New(ErrPermSuperAdminRequired, "super admin access required", 403)
}

func PermSellerTierRequired(tier string) *AppError {
	return New(ErrPermSellerTierRequired, fmt.Sprintf("this feature requires %s tier or higher", tier), 403).
		WithContext("tier", tier)
}

func UploadFileTooLarge(maxMB int) *AppError {
	return New(ErrUploadFileTooLarge, fmt.Sprintf("image must be smaller than %d MB", maxMB), 400).
		WithContext("max_mb", maxMB)
}

func UploadWrongType() *AppError {
	return New(ErrUploadWrongType, "only JPEG, PNG and WebP images are allowed", 400)
}

func PayInsufficientBalance(required string) *AppError {
	return New(ErrPayInsufficientBalance, fmt.Sprintf("insufficient balance — ETB %s required", required), 400).
		WithContext("amount", required)
}

func ProdNotFound() *AppError {
	return New(ErrProdNotFound, "product not found", 404)
}

func ProdPriceMismatch() *AppError {
	return New(ErrProdPriceMismatch, "transaction amount does not match the listed price", 400)
}

func ValPhoneRequired() *AppError {
	return New(ErrValPhoneRequired, "phone number is required", 400)
}

func ValPhoneInvalid() *AppError {
	return New(ErrValPhoneInvalid, "invalid phone number format", 400)
}

func ValPasswordMin() *AppError {
	return New(ErrValPasswordMin, "password must be at least 8 characters", 400)
}

func ValRequired(field string) *AppError {
	return New(ErrValRequired, fmt.Sprintf("%s is required", field), 400).
		WithContext("field", field)
}

func NotFound(resource string) *AppError {
	return New(ErrNFUser, fmt.Sprintf("%s not found", resource), 404).
		WithContext("resource", resource)
}

func InternalError() *AppError {
	return New(ErrSrvInternal, "an internal error occurred", 500)
}

func DatabaseError() *AppError {
	return New(ErrSrvDatabase, "database error occurred", 500)
}
