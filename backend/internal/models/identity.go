package models

import "time"

const (
	IdentityStatusPending   = "pending"
	IdentityStatusVerified  = "verified"
	IdentityStatusRejected  = "rejected"
	IdentityStatusEscalated = "escalated"
	IdentityStatusExpired   = "expired"
)

// IdentityVerification records official national identity (e.g. Fayda) verification audits.
type IdentityVerification struct {
	ID                   int64                  `json:"id"`
	UserID               int                    `json:"user_id"`
	ProviderCode         string                 `json:"provider_code"` // fayda, manual_review
	CredentialType       string                 `json:"credential_type"` // national_id, passport, kebele_id
	CredentialIDHash     string                 `json:"-"` // SHA-256 hash for deduplication
	IDMasked             string                 `json:"id_masked"` // Masked identifier (e.g. ******1234)
	FullLegalName        string                 `json:"full_legal_name"`
	PhoneMatched         bool                   `json:"phone_matched"`
	PhoneMatchStatus     string                 `json:"phone_match_status"` // verified, unsupported_by_provider, mismatch
	AuthReference        string                 `json:"auth_reference"`
	Status               string                 `json:"status"` // pending, verified, rejected, escalated, expired
	ConsentRecorded      bool                   `json:"consent_recorded"`
	ConsentTimestamp     time.Time              `json:"consent_timestamp"`
	VerificationMetadata map[string]interface{} `json:"verification_metadata,omitempty"`
	ExceptionReason      *string                `json:"exception_reason,omitempty"`
	ReviewedBy           *int                   `json:"reviewed_by,omitempty"`
	VerifiedAt           *time.Time             `json:"verified_at,omitempty"`
	ExpiresAt            *time.Time             `json:"expires_at,omitempty"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

