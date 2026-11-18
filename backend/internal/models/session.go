package models

import (
	"time"

	"github.com/google/uuid"
)

// UserRefreshToken represents a persistent server-side refresh token lineage
// for cookie-based session rotation, revocation, and replay detection.
type UserRefreshToken struct {
	ID               uuid.UUID  `json:"id"`
	UserID           int        `json:"user_id"`
	FamilyID         uuid.UUID  `json:"family_id"`
	TokenHash        string     `json:"-"`
	ParentTokenHash  string     `json:"-"`
	IsRevoked        bool       `json:"is_revoked"`
	RevocationReason string     `json:"revocation_reason,omitempty"`
	RevokedAt        *time.Time `json:"revoked_at,omitempty"`
	ExpiresAt        time.Time  `json:"expires_at"`
	UserAgent        string     `json:"user_agent"`
	IPAddress        string     `json:"ip_address"`
	CreatedAt        time.Time  `json:"created_at"`
	LastUsedAt       time.Time  `json:"last_used_at"`
}
