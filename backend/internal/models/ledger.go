package models

import "time"

// WalletAccount represents the dedicated double-entry financial account for a user.
type WalletAccount struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	Currency         string    `json:"currency"`
	AvailableBalance float64   `json:"available_balance"`
	PendingBalance   float64   `json:"pending_balance"`
	ReservedBalance  float64   `json:"reserved_balance"`
	TotalWithdrawn   float64   `json:"total_withdrawn"`
	LifetimeEarned   float64   `json:"lifetime_earned"`
	IsLocked         bool      `json:"is_locked"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// LedgerTransaction represents an immutable balanced entry in the financial ledger.
type LedgerTransaction struct {
	ID              int64     `json:"id"`
	LedgerTxID      string    `json:"ledger_tx_id"`
	UserID          int       `json:"user_id"`
	DebitAccount    string    `json:"debit_account"`
	CreditAccount   string    `json:"credit_account"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	EntryType       string    `json:"entry_type"` // referral_reward_a, referral_commission_b, withdrawal_hold, withdrawal_settle, withdrawal_reversal, manual_compensation
	ReferenceID     *string   `json:"reference_id,omitempty"`
	SourceEvent     string    `json:"source_event"`
	Status          string    `json:"status"` // posted, pending, reversed
	IdempotencyKey  string    `json:"idempotency_key"`
	Notes           *string   `json:"notes,omitempty"`
	ActorID         *int      `json:"actor_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
