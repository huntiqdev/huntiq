package constants

// Listing duration options (in days)
var ListingDurations = []int{2, 7, 30, 180, 365}

// Fee configuration keys (these are stored in platform_settings table)
const (
	// Broker fee
	SettingBrokerFeeRate = "broker_fee_rate" // Default: 2.0%

	// Listing fees by tier
	SettingNormalStartupFee       = "normal_startup_fee"        // Default: 1.0%
	SettingGoldStartupFee         = "gold_startup_fee"          // Default: 0.5%
	SettingPremiumStartupFee      = "premium_startup_fee"       // Default: 0.25%
	SettingNormalWeeklyIncrement  = "normal_weekly_increment"   // Default: 0.5%
	SettingGoldWeeklyIncrement    = "gold_weekly_increment"     // Default: 0.25%
	SettingPremiumWeeklyIncrement = "premium_weekly_increment"  // Default: 0.1%
	SettingListingFreeWeeks       = "listing_free_weeks"        // Default: 1

	// Tier upgrade costs
	SettingGoldThreshold    = "gold_threshold"     // Default: 5900 ETB
	SettingPremiumThreshold = "premium_threshold"  // Default: 15900 ETB
	SettingBuyerGoldCost    = "buyer_gold_cost"    // Default: 5900 ETB
	SettingBuyerPremiumCost = "buyer_premium_cost" // Default: 15900 ETB

	// Image limits
	SettingNormalMaxImages  = "normal_max_images"  // Default: 2
	SettingGoldMaxImages    = "gold_max_images"    // Default: 7
	SettingPremiumMaxImages = "premium_max_images" // Default: 20
	SettingMinImagesPerListing = "min_images_per_listing" // Default: 2

	// Content limits
	SettingMaxTitleLength       = "max_title_length"       // Default: 200
	SettingMaxDescriptionLength = "max_description_length" // Default: 5000

	// Location limits per tier
	SettingNormalMaxLocations  = "normal_max_locations"  // Default: 1
	SettingGoldMaxLocations    = "gold_max_locations"    // Default: 3
	SettingPremiumMaxLocations = "premium_max_locations" // Default: 10

	// Referral
	SettingReferralCommissionRate = "referral_commission_rate" // Default: 1.2%

	// Banner limits
	SettingMaxBannersBuyer  = "max_banners_buyer"  // Default: 5
	SettingMaxBannersSeller = "max_banners_seller" // Default: 5

	// File upload limits
	SettingMaxImageSizeMB = "max_image_size_mb" // Default: 20

	// Currency & payment (percentage-based tolerance)
	SettingCurrencyTolerancePercent = "currency_tolerance_percent" // Default: 0.1 (0.1%, not 10%)

	// Referral withdrawal limits
	SettingMinReferralWithdrawal = "min_referral_withdrawal" // Default: 500
	SettingMaxReferralWithdrawal = "max_referral_withdrawal" // Default: 1000000

	// Promotion banner
	SettingPromotionMessage    = "promotion_message"
	SettingPromotionActive     = "promotion_active"
	SettingPromotionBgColor    = "promotion_bg_color"
	SettingPromotionTextColor  = "promotion_text_color"
	SettingPromotionLinkURL    = "promotion_link_url"
	SettingPromotionLinkLabel  = "promotion_link_label"
)

// Default fee parameters (fallback if DB not configured — admin can override all via platform_settings)
const (
	DefaultBrokerFeeRate          = 2.0  // 2%
	DefaultListingFreeWeeks       = 1
	DefaultNormalStartupFee       = 1.0  // 1%
	DefaultGoldStartupFee         = 0.5  // 0.5%
	DefaultPremiumStartupFee      = 0.25 // 0.25%
	DefaultNormalWeeklyIncrement  = 0.5  // 0.5%
	DefaultGoldWeeklyIncrement    = 0.25 // 0.25%
	DefaultPremiumWeeklyIncrement = 0.1  // 0.1%
	DefaultGoldThreshold          = 5900.0
	DefaultPremiumThreshold       = 15900.0
	DefaultBuyerGoldCost          = 5900.0
	DefaultBuyerPremiumCost       = 15900.0
	DefaultReferralCommissionRate = 1.2  // 1.2%
	DefaultMaxTitleLength         = 200
	DefaultMaxDescriptionLength   = 5000
	DefaultMinImages              = 2
	DefaultNormalMaxLocations     = 1
	DefaultGoldMaxLocations       = 3
	DefaultPremiumMaxLocations    = 10
	DefaultMaxBanners             = 5
	DefaultMaxImageSizeMB         = 20
	DefaultCurrencyTolerancePercent = 0.1 // 0.1% (NOT 10%!)
	DefaultMinReferralWithdrawal  = 500.0
	DefaultMaxReferralWithdrawal  = 1000000.0
)
