package models

import "time"

type OrderRequest struct {
	ID               int        `json:"id"`
	ProductID        int        `json:"product_id"`
	SellerID         int        `json:"seller_id"`
	BuyerName        string     `json:"buyer_name"`
	BuyerPhone       string     `json:"buyer_phone"`
	BuyerEmail       string     `json:"buyer_email,omitempty"`
	Message          string     `json:"message,omitempty"`
	OfferPrice       *float64   `json:"offer_price,omitempty"`
	Status           string     `json:"status"`
	SellerNote       string     `json:"seller_note,omitempty"`
	BuyerTierAtOrder string     `json:"buyer_tier_at_order"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Joined
	ProductTitle string  `json:"product_title,omitempty"`
	ProductPrice float64 `json:"product_price,omitempty"`
}
