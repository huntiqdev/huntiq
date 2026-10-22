package dto

// UpdateSellerProfileRequest is the body for PUT /api/seller/profile.
// Only safe profile fields are accepted — role/tier fields are intentionally excluded
// to prevent injection via request body.
type UpdateSellerProfileRequest struct {
	Name             string   `json:"name"`
	BusinessName     string   `json:"business_name"`
	Email            string   `json:"email"`
	AdditionalPhones []string `json:"additional_phones"`
}

// SellerUpgradeRequest is the body for POST /api/seller/upgrade-tier.
type SellerUpgradeRequest struct {
	RequestedTier string  `json:"requested_tier"` // gold | premium
	PaymentMethod string  `json:"payment_method"` // telebirr | safaricom | cbe
	PayerPhone    string  `json:"payer_phone"`
	Reference     string  `json:"reference"`
	Amount        float64 `json:"amount"`
	Notes         string  `json:"notes"`
}
