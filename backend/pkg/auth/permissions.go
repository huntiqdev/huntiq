package auth

import "github.com/your-org/used-tools-marketplace/pkg/constants"

// Permission names
const (
	PermProducts          = "products"
	PermPending           = "pending"
	PermSellers           = "sellers"
	PermBuyerUpgrades     = "buyer_upgrades"
	PermSellerUpgrades    = "seller_upgrades"
	PermTransactions      = "transactions"
	PermProfit            = "profit"
	PermAnalytics         = "analytics"
	PermSettings          = "settings"
	PermNotifications     = "notifications"
	PermCategories        = "categories"
	PermLocations         = "locations"
	PermBanners           = "banners"
	PermOTPLogs           = "otp_logs"
	PermPremiumUsers      = "premium_users"
	PermCategoryTrends    = "category_trends"
	PermWithdrawals       = "withdrawals"
)

// RequiresTier checks if a user's tier meets the requirement
func RequiresTier(userTier, requiredTier string) bool {
	return constants.MeetsRequiredTier(userTier, requiredTier)
}
