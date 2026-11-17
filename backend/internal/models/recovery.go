package models

import "time"

const (
	RecoveryStatusPending  = "pending"
	RecoveryStatusVerified = "verified"
	RecoveryStatusConsumed = "consumed"
	RecoveryStatusExpired  = "expired"
)

// PasswordRecoveryRequest records an OTP password recovery lifecycle with 3-hour cooldown enforcement.
type PasswordRecoveryRequest struct {
	ID                  int64      `json:"id"`
	UserID              int        `json:"user_id"`
	Phone               string     `json:"phone"`
	OTPHash             string     `json:"-"`
	Attempts            int        `json:"attempts"`
	MaxAttempts         int        `json:"max_attempts"`
	Status              string     `json:"status"` // pending, verified, consumed, expired
	ResetTokenHash      *string    `json:"-"`
	ResetTokenExpiresAt *time.Time `json:"reset_token_expires_at,omitempty"`
	CooldownUntil       time.Time  `json:"cooldown_until"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

