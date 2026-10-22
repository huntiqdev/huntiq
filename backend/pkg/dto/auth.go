package dto

type RegisterRequest struct {
	Name             string `json:"name"`
	Phone            string `json:"phone"`
	BusinessName     string `json:"business_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Role             string `json:"role"`              // seller | buyer
	ReferralCode     string `json:"referral_code"`     // optional ?ref= code from registration link
	SellerPreference string `json:"seller_preference"` // sale | rent | both (for sellers only)
}

type LoginRequest struct {
	// Admins log in with email; sellers/buyers log in with phone.
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// OTPRequest is used to request an OTP be sent to a phone number.
type OTPRequest struct {
	Phone    string `json:"phone"`
	Purpose  string `json:"purpose"`   // login | register
	UserType string `json:"user_type"` // seller | buyer | admin (optional, defaults to seller)
}

// OTPVerifyRequest verifies the OTP and completes login/registration.
type OTPVerifyRequest struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`   // used for admin OTP
	Code    string `json:"code"`
	Purpose string `json:"purpose"`
	Role    string `json:"role"`    // seller | buyer (used by OTP login endpoint)
}

// AdminOTPRequest is sent after step-1 admin login to verify the email OTP.
type AdminOTPRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// AdminLoginStep1Response is returned after credentials are verified.
// If is_super=true, token is set and no OTP step is needed.
// If is_super=false, email is set and the client must proceed to step 2 with OTP.
type AdminLoginStep1Response struct {
	Email   string `json:"email,omitempty"` // set when OTP required (sub-admin)
	Token   string `json:"token,omitempty"` // set when no OTP needed (super admin)
	IsSuper bool   `json:"is_super"`
	Message string `json:"message"`
}

// BuyerUpgradeRequest requests a buyer tier upgrade.
type BuyerUpgradeRequest struct {
	Tier          string  `json:"tier"`           // gold | premium
	Amount        float64 `json:"amount"`         // amount the user claims to have paid — validated server-side
	DurationDays  int     `json:"duration_days"`  // duration of the tier upgrade
	Reference     string  `json:"reference"`      // payment transaction reference
	PayerPhone    string  `json:"payer_phone"`    // phone used to make the payment
	PaymentMethod string  `json:"payment_method"` // telebirr | safaricom | cbe
}

// DemandSelectionRequest is sent after successful registration to indicate user's current demand
type DemandSelectionRequest struct {
	Demand string `json:"demand"` // FIND | LIST
}

// DemandSelectionResponse returns the next route based on demand
type DemandSelectionResponse struct {
	Demand   string `json:"demand"`
	NextPath string `json:"next_path"`
	Message  string `json:"message"`
}
