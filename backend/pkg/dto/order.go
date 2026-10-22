package dto

type CreateOrderRequest struct {
	ProductID  int      `json:"product_id"`
	BuyerName  string   `json:"buyer_name"`
	BuyerPhone string   `json:"buyer_phone"`
	BuyerEmail string   `json:"buyer_email"`
	Message    string   `json:"message"`
	OfferPrice *float64 `json:"offer_price"`
	// BuyerTier is injected server-side from the JWT claim — never from the request body.
	// The "-" tag prevents clients from supplying this field.
	BuyerTier string `json:"-"`
}

type UpdateOrderRequest struct {
	Status     string `json:"status"` // accepted | declined | completed | cancelled
	SellerNote string `json:"seller_note"`
}
