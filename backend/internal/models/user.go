package models

import "time"

type User struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Phone              string    `json:"phone"`
	AdditionalPhones   []string  `json:"additional_phones,omitempty"`
	BusinessName       string    `json:"business_name,omitempty"`
	Email              string    `json:"email"`
	PasswordHash       string    `json:"-"`
	Rating             float64   `json:"rating"`
	Role               string    `json:"role"`
	AdminRole          string    `json:"admin_role,omitempty"`
	Permissions        []string  `json:"permissions,omitempty"`
	IsActive           bool      `json:"is_active"`
	DeactivationReason string    `json:"deactivation_reason,omitempty"`
	Tier               string    `json:"tier"`
	TotalValueSold     float64   `json:"total_value_sold"`
	NetProfit          float64   `json:"net_profit,omitempty"`
	BuyerTier          string    `json:"buyer_tier"`
	OTPVerified        bool      `json:"otp_verified"`
	SellerPreference   string    `json:"seller_preference,omitempty"` // 'sale', 'rent', or 'both'
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// MaxAdditionalPhones returns the maximum number of additional phones allowed for a tier
func (u *User) MaxAdditionalPhones() int {
	switch u.Tier {
	case "gold":
		return 3
	case "premium":
		return -1 // unlimited
	default:
		return 0 // normal tier: no additional phones
	}
}

// AllAdminPermissions is the canonical list of page/feature permission keys
// that can be assigned to normal admins. Super admins always have full access.
// Keys are grouped by feature area for the UI.
var AllAdminPermissions = []string{
	// Management
	"dashboard", "pending", "products", "sold",
	// People
	"sellers", "seller_upgrades", "buyer_upgrades", "transactions",
	// Finance (sensitive — off by default for new admins)
	"profit",
	"referrals",   // referral management + withdrawals
	"withdrawals", // alias for referral withdrawal processing
	// Users
	"premium_users",
	// System
	"analytics", "notifications", "locations", "settings", "categories", "category_trends",
	// Audit — all sensitive, off by default
	"logs",
	"otp_logs",
	// Content
	"banners",
	// Operations & Safety (Priority 1 enhancements - off by default)
	"operations",          // Platform operations dashboard, job monitoring
	"seller_verification", // Seller identity and business document verification
	"trust_safety",        // Fraud detection, abuse reports, trust & safety
	"refunds",             // Refund and dispute management
	"appeals",             // Account appeals and recovery
}

// DefaultAdminPermissions are granted automatically to new normal admins.
// Sensitive keys (profit, logs, otp_logs, transactions, premium_users, referrals, withdrawals) are excluded.
var DefaultAdminPermissions = []string{
	"dashboard", "pending", "products", "sold", "sellers",
	"analytics", "notifications", "locations", "settings", "categories",
}
