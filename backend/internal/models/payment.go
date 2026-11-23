package models

import "time"

// Payment statuses reflecting full financial lifecycle
const (
	PaymentStatusCreated    = "CREATED"
	PaymentStatusPending    = "PENDING"
	PaymentStatusProcessing = "PROCESSING"
	PaymentStatusSuccessful = "SUCCESSFUL"
	PaymentStatusFailed     = "FAILED"
	PaymentStatusCancelled  = "CANCELLED"
	PaymentStatusExpired    = "EXPIRED"
)

// Payment operating environments
const (
	PaymentEnvDemo       = "demo"
	PaymentEnvProduction = "production"
)

// Payment order types
const (
	OrderTypeTierUpgrade = "tier_upgrade"
	OrderTypeListingFee  = "listing_fee"
	OrderTypeMembership  = "membership"
)

// PaymentOrder represents an individual unified payment transaction
type PaymentOrder struct {
	ID                   int                    `json:"id"`
	OrderID              string                 `json:"order_id"`
	UserID               int                    `json:"user_id"`
	UserName             string                 `json:"user_name,omitempty"`
	UserPhone            string                 `json:"user_phone,omitempty"`
	Provider             string                 `json:"provider"`
	Environment          string                 `json:"environment"`
	OrderType            string                 `json:"order_type"`
	TargetTier           *string                `json:"target_tier,omitempty"`
	UserRole             string                 `json:"user_role"`
	DurationDays         int                    `json:"duration_days"`
	ProductID            *int                   `json:"product_id,omitempty"`
	ProductTitle         string                 `json:"product_title,omitempty"`
	Amount               float64                `json:"amount"`
	Currency             string                 `json:"currency"`
	Title                string                 `json:"title"`
	Status               string                 `json:"status"`
	PayerPhone           *string                `json:"payer_phone,omitempty"`
	ProviderTxID         *string                `json:"provider_tx_id,omitempty"`
	CheckoutURL          *string                `json:"checkout_url,omitempty"`
	RawCallback          map[string]interface{} `json:"raw_callback,omitempty"`
	VerificationMetadata map[string]interface{} `json:"verification_metadata,omitempty"`
	IsDemo               bool                   `json:"is_demo"`
	OTPVerified          bool                   `json:"otp_verified"`
	OTPVerifiedAt        *time.Time             `json:"otp_verified_at,omitempty"`
	VerifiedAt           *time.Time             `json:"verified_at,omitempty"`
	FailureReason        *string                `json:"failure_reason,omitempty"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

// PaymentProviderConfig represents provider setup and credentials
type PaymentProviderConfig struct {
	ID              int                    `json:"id"`
	ProviderCode    string                 `json:"provider_code"`
	ProviderName    string                 `json:"provider_name"`
	EnvironmentMode string                 `json:"environment_mode"`
	IsEnabled       bool                   `json:"is_enabled"`
	APIEndpoint     string                 `json:"api_endpoint"`
	AppID           string                 `json:"app_id"`
	AppKeyMasked    string                 `json:"app_key_masked"`
	ShortCode       string                 `json:"short_code"`
	NotifyURL       string                 `json:"notify_url"`
	ReturnURL       string                 `json:"return_url"`
	ExtraConfig     map[string]interface{} `json:"extra_config,omitempty"`
	UpdatedBy       *int                   `json:"updated_by,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}
