package models

import "time"

// TelebirrOrder represents a payment order initiated through Telebirr.
type TelebirrOrder struct {
	ID            int                    `json:"id"`
	OutTradeNo    string                 `json:"out_trade_no"`
	UserID        int                    `json:"user_id"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	Title         string                 `json:"title"`
	OrderType     string                 `json:"order_type"` // listing_fee, tier_upgrade, membership
	Status        string                 `json:"status"`     // pending, completed, failed, refunded
	PrepayID      *string                `json:"prepay_id,omitempty"`
	TransactionNo *string                `json:"transaction_no,omitempty"`
	PayerPhone    *string                `json:"payer_phone,omitempty"`
	RawCallback   map[string]interface{} `json:"raw_callback,omitempty"`
	VerifiedAt    *time.Time             `json:"verified_at,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

