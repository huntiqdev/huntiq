package models

import "time"

// ReferralToken represents an exact, high-entropy attribution token.
type ReferralToken struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Token           string    `json:"token"`
	ReferralCode    string    `json:"referral_code"`
	ClickCount      int       `json:"click_count"`
	ConversionCount int       `json:"conversion_count"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ReferralAttribution tracks an attribution event originating from a referral token.
type ReferralAttribution struct {
	ID               int64     `json:"id"`
	TokenID          int       `json:"token_id"`
	InviterID        int       `json:"inviter_id"`
	IPHash           string    `json:"ip_hash"`
	UserAgentHash    *string   `json:"user_agent_hash,omitempty"`
	LandingURL       string    `json:"landing_url"`
	Status           string    `json:"status"` // clicked, registered, converted, expired
	RegisteredUserID *int      `json:"registered_user_id,omitempty"`
	AttributedAt     time.Time `json:"attributed_at"`
	ExpiresAt        time.Time `json:"expires_at"`
}

// ReferralRewardPolicy defines versioned business rules for referral rewards.
type ReferralRewardPolicy struct {
	ID                        int        `json:"id"`
	PolicyName                string     `json:"policy_name"`
	RegistrationRewardAmount  float64    `json:"registration_reward_amount"`
	TransactionRewardPercent  float64    `json:"transaction_reward_percent"`
	MinQualifyingTxAmount     float64    `json:"min_qualifying_tx_amount"`
	MaxRewardPerTx            float64    `json:"max_reward_per_tx"`
	MaxLifetimePerReferred    float64    `json:"max_lifetime_per_referred"`
	MaxRewardsPerInviter      float64    `json:"max_rewards_per_inviter"`
	RewardHoldingDays         int        `json:"reward_holding_days"`
	IsActive                  bool       `json:"is_active"`
	Version                   int        `json:"version"`
	EffectiveFrom             time.Time  `json:"effective_from"`
	EffectiveTo               *time.Time `json:"effective_to,omitempty"`
	UpdatedBy                 *int       `json:"updated_by,omitempty"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
}

// ReferralRewardPolicyVersion preserves immutable history of financial policy changes.
type ReferralRewardPolicyVersion struct {
	ID                       int       `json:"id"`
	PolicyID                 int       `json:"policy_id"`
	Version                  int       `json:"version"`
	RegistrationRewardAmount float64   `json:"registration_reward_amount"`
	TransactionRewardPercent float64   `json:"transaction_reward_percent"`
	MinQualifyingTxAmount    float64   `json:"min_qualifying_tx_amount"`
	MaxRewardPerTx           float64   `json:"max_reward_per_tx"`
	RewardHoldingDays        int       `json:"reward_holding_days"`
	ChangedBy                *int      `json:"changed_by,omitempty"`
	ChangeReason             string    `json:"change_reason"`
	CreatedAt                time.Time `json:"created_at"`
}
