// Package pkg provides a unified import point for all shared packages.
// Import this package to access all common utilities, validations, constants, and errors.
package pkg

import (
	"github.com/your-org/used-tools-marketplace/pkg/auth"
	"github.com/your-org/used-tools-marketplace/pkg/constants"
	"github.com/your-org/used-tools-marketplace/pkg/errors"
	"github.com/your-org/used-tools-marketplace/pkg/fees"
	"github.com/your-org/used-tools-marketplace/pkg/referral"
	"github.com/your-org/used-tools-marketplace/pkg/response"
	"github.com/your-org/used-tools-marketplace/pkg/validations"
)

// Re-export commonly used types and functions for convenience
type (
	// Errors
	ErrorCode = errors.ErrorCode
	AppError  = errors.AppError

	// Validations
	PhoneResult                = validations.PhoneResult
	ProductValidationResult    = validations.ProductValidationResult
	OrderValidationResult      = validations.OrderValidationResult
	LocationValidationResult   = validations.LocationValidationResult

	// Fees
	FeeConfig = fees.FeeConfig
	FeeResult = fees.FeeResult

	// Payment
	PaymentMethod = constants.PaymentMethod
)

// Validation functions
var (
	ValidatePhone    = validations.ValidateEthiopianPhone
	ValidateEmail    = validations.ValidateEmail
	ValidatePassword = validations.ValidatePassword
	ValidateProduct  = validations.ValidateProduct
	ValidateOrder    = validations.ValidateOrderRequest
	ValidateLocation = validations.ValidateLocation
	ValidateImageURL = validations.ValidateImageURL
	IsValidPhone     = validations.IsValidPhone
)

// Constants
var (
	TierRank            = constants.TierRank
	TierNormal          = constants.TierNormal
	TierGold            = constants.TierGold
	TierPremium         = constants.TierPremium
	MeetsRequiredTier   = constants.MeetsRequiredTier
	CanUpgradeTo        = constants.CanUpgradeTo
	PaymentMethods      = constants.PaymentMethods
	IsValidPaymentMethod = constants.IsValidPaymentMethod
	ListingDurations    = constants.ListingDurations
	// MinImages and GetMaxImages removed - use platform_settings instead
)

// Fee calculation
var (
	CalculateListingFee = fees.CalculateListingFee
	AmountsMatch        = fees.AmountsMatch
)

// Referral
var (
	ValidateReferralCode = referral.ValidateReferralCode
)

// Error constructors
var (
	NewAppError              = errors.New
	AuthInvalidCredentials   = errors.AuthInvalidCredentials
	AuthWrongPassword        = errors.AuthWrongPassword
	AuthPhoneNotVerified     = errors.AuthPhoneNotVerified
	AuthSessionExpired       = errors.AuthSessionExpired
	AuthTokenInvalid         = errors.AuthTokenInvalid
	AuthRegistrationClosed   = errors.AuthRegistrationClosed
	AuthPhoneAlreadyRegistered = errors.AuthPhoneAlreadyRegistered
	OTPLocked                = errors.OTPLocked
	OTPExpired               = errors.OTPExpired
	OTPInvalid               = errors.OTPInvalid
	OTPRequired              = errors.OTPRequired
	OTPSendFailed            = errors.OTPSendFailed
	PermDenied               = errors.PermDenied
	PermSuperAdminRequired   = errors.PermSuperAdminRequired
	PermSellerTierRequired   = errors.PermSellerTierRequired
	UploadFileTooLarge       = errors.UploadFileTooLarge
	UploadWrongType          = errors.UploadWrongType
	PayInsufficientBalance   = errors.PayInsufficientBalance
	ProdNotFound             = errors.ProdNotFound
	ProdPriceMismatch        = errors.ProdPriceMismatch
	ValPhoneRequired         = errors.ValPhoneRequired
	ValPhoneInvalid          = errors.ValPhoneInvalid
	ValPasswordMin           = errors.ValPasswordMin
	ValRequired              = errors.ValRequired
	NotFound                 = errors.NotFound
	InternalError            = errors.InternalError
	DatabaseError            = errors.DatabaseError
)

// Response helpers
var (
	OK           = response.OK
	OKList       = response.OKList
	AppErrorResp = response.AppError
	BadRequest   = response.BadRequest
	Unauthorized = response.Unauthorized
	NotFoundResp = response.NotFound
)

// Permission constants
var (
	PermProducts       = auth.PermProducts
	PermPending        = auth.PermPending
	PermSellers        = auth.PermSellers
	PermBuyerUpgrades  = auth.PermBuyerUpgrades
	PermSellerUpgrades = auth.PermSellerUpgrades
	PermTransactions   = auth.PermTransactions
	PermProfit         = auth.PermProfit
	PermAnalytics      = auth.PermAnalytics
	PermSettings       = auth.PermSettings
	RequiresTier       = auth.RequiresTier
)
