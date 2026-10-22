package models

import "time"

type Transaction struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	ProductID        int       `json:"product_id"`
	SellerID         int       `json:"seller_id"`
	BuyerName        string    `json:"buyer_name"`
	BuyerPhone       string    `json:"buyer_phone"`
	Amount           float64   `json:"amount"`
	FeeAmount        float64   `json:"fee_amount"`
	BrokerFee        float64   `json:"broker_fee"`
	SellerGets       float64   `json:"seller_gets"`
	TransactionType  string    `json:"transaction_type"`
	PaymentMethod    string    `json:"payment_method,omitempty"`
	PaymentReference string    `json:"payment_reference,omitempty"`
	PayerPhone       string    `json:"payer_phone,omitempty"`
	Status           string    `json:"status"`
	Notes            string    `json:"notes,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

