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
	Capabilities       []string  `json:"capabilities,omitempty"`      // additive: 'buyer', 'seller', 'landlord', 'renter'
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// UserCapability represents an individual additive capability assigned to a user account
type UserCapability struct {
	ID          int                    `json:"id"`
	UserID      int                    `json:"user_id"`
	Capability  string                 `json:"capability"` // buyer, seller, landlord, renter
	Status      string                 `json:"status"`     // active, pending, suspended, revoked
	VerifiedAt  *time.Time             `json:"verified_at,omitempty"`
	ActivatedAt time.Time              `json:"activated_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// HasCapability checks if the user possesses an active marketplace capability
func (u *User) HasCapability(cap string) bool {
	for _, c := range u.Capabilities {
		if c == cap {
			return true
		}
	}
	// Backward compatibility fallback: seller role includes buyer and seller capabilities
	if cap == "buyer" {
		return true // all registered accounts are buyers
	}
	if cap == "seller" && u.Role == "seller" {
		return true
	}
	if cap == "landlord" && (u.SellerPreference == "rent" || u.SellerPreference == "both") {
		return true
	}
	return false
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

// IsPlatformOwner checks if user is Platform Owner or legacy super admin
func (u *User) IsPlatformOwner() bool {
	return u.Role == "admin" && (u.AdminRole == RolePlatformOwner || u.AdminRole == LegacyRoleSuper)
}

// IsPlatformAdmin checks if user is Platform Admin or Platform Owner
func (u *User) IsPlatformAdmin() bool {
	return u.Role == "admin" && (u.AdminRole == RolePlatformAdmin || u.IsPlatformOwner())
}

// IsSecurityAdmin checks if user is Security Admin or Platform Owner
func (u *User) IsSecurityAdmin() bool {
	return u.Role == "admin" && (u.AdminRole == RoleSecurityAdmin || u.IsPlatformOwner())
}

// IsOperationsAdmin checks if user is Operations Admin or Platform Owner
func (u *User) IsOperationsAdmin() bool {
	return u.Role == "admin" && (u.AdminRole == RoleOperationsAdmin || u.IsPlatformOwner())
}

// RoleDomainPermissions defines the baseline domain permissions granted to enterprise administrative roles.
// Platform Owner and Legacy Super have all permissions by default.
var RoleDomainPermissions = map[string]map[string]bool{
	RolePlatformAdmin: {
		"dashboard": true, "settings": true, "categories": true, "category_trends": true,
		"analytics": true, "locations": true, "banners": true, "notifications": true,
		"profit": true, "finance": true, "sellers": true, "transactions": true,
	},
	RoleSecurityAdmin: {
		"dashboard": true, "logs": true, "otp_logs": true, "trust_safety": true,
	},
	RoleOperationsAdmin: {
		"dashboard": true, "operations": true, "seller_verification": true, "refunds": true,
		"appeals": true, "pending": true, "products": true, "sellers": true,
		"seller_upgrades": true, "buyer_upgrades": true, "transactions": true,
	},
	RoleModerator: {
		"dashboard": true, "pending": true, "products": true, "trust_safety": true,
	},
	RoleSupportAgent: {
		"dashboard": true, "sellers": true, "seller_upgrades": true, "buyer_upgrades": true,
	},
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
	"finance",        // Financial Planning & Analysis (FP&A) view
	"finance:manage", // Financial Planning & Analysis (FP&A) manage assumptions
	"referrals",      // referral management + withdrawals
	"withdrawals",    // alias for referral withdrawal processing
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
