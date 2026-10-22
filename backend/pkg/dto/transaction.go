package dto

type CreateTransactionRequest struct {
	ProductID  int     `json:"product_id"`
	BuyerName  string  `json:"buyer_name"`
	BuyerPhone string  `json:"buyer_phone"`
	Amount     float64 `json:"amount"`
	Notes      string  `json:"notes"`
}

// ListingFeePaymentRequest is submitted by a seller to record their listing fee payment.
// This is separate from buyer transactions — it does NOT mark the product as sold.
type ListingFeePaymentRequest struct {
	ProductID     int     `json:"product_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"` // telebirr | safaricom | cbe
	PayerPhone    string  `json:"payer_phone"`    // normalized +251XXXXXXXXX
	Reference     string  `json:"reference"`      // transaction ID from SMS
	Notes         string  `json:"notes"`
}
