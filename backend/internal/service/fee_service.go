package service

import (
	"math"
	"github.com/your-org/used-tools-marketplace/pkg/constants"
)

// SellerTier represents the seller's tier level.
type SellerTier string

const (
	TierNormal  SellerTier = "normal"
	TierGold    SellerTier = "gold"
	TierPremium SellerTier = "premium"
)

// TierRank maps tier names to a numeric rank for hierarchy comparisons.
// DEPRECATED: Use constants.TierRank instead for single source of truth.
var TierRank = constants.TierRank

// TierAllowed returns true if buyerTier meets or exceeds requiredTier.
// Signature: TierAllowed(buyerTier, requiredTier) — buyer first, required second.
// DEPRECATED: Use constants.MeetsRequiredTier directly for clarity.
func TierAllowed(buyerTier, requiredTier string) bool {
	return constants.MeetsRequiredTier(buyerTier, requiredTier)
}

// TierThresholds — total value sold required to reach each tier.
// These are defaults; actual values come from platform_settings.
const (
	GoldThreshold    = constants.DefaultGoldThreshold
	PremiumThreshold = constants.DefaultPremiumThreshold
)

// ListingDuration options in days.
// DEPRECATED: Use constants.ListingDurations instead.
var ListingDurations = constants.ListingDurations

// FeeSettings holds all configurable fee parameters loaded from platform_settings.
type FeeSettings struct {
	NormalStartupFee       float64 // default from constants.DefaultNormalStartupFee
	GoldStartupFee         float64 // default from constants.DefaultGoldStartupFee
	PremiumStartupFee      float64 // default from constants.DefaultPremiumStartupFee
	NormalWeeklyIncrement  float64 // default from constants.DefaultNormalWeeklyIncrement
	GoldWeeklyIncrement    float64 // default from constants.DefaultGoldWeeklyIncrement
	PremiumWeeklyIncrement float64 // default from constants.DefaultPremiumWeeklyIncrement
	FreeWeeks              int     // default from constants.DefaultListingFreeWeeks
}

// DefaultFeeSettings returns the hardcoded defaults used when DB settings are unavailable.
func DefaultFeeSettings() FeeSettings {
	return FeeSettings{
		NormalStartupFee:       constants.DefaultNormalStartupFee,
		GoldStartupFee:         constants.DefaultGoldStartupFee,
		PremiumStartupFee:      constants.DefaultPremiumStartupFee,
		NormalWeeklyIncrement:  constants.DefaultNormalWeeklyIncrement,
		GoldWeeklyIncrement:    constants.DefaultGoldWeeklyIncrement,
		PremiumWeeklyIncrement: constants.DefaultPremiumWeeklyIncrement,
		FreeWeeks:              constants.DefaultListingFreeWeeks,
	}
}

// FeeBreakdown holds the full fee calculation result.
type FeeBreakdown struct {
	Tier               string  `json:"tier"`
	ProductValue       float64 `json:"product_value"`
	DurationDays       int     `json:"duration_days"`
	StartupPct         float64 `json:"startup_pct"`
	WeeklyIncrementPct float64 `json:"weekly_increment_pct"`
	FreeWeeks          int     `json:"free_weeks"`
	BillableWeeks      int     `json:"billable_weeks"`
	TotalFeePercent    float64 `json:"total_fee_percent"`
	TotalFeeAmount     float64 `json:"total_fee_amount"`
	SellerGets         float64 `json:"seller_gets"`
}

// CalculateListingFee computes the total service fee for a listing using provided settings.
//
// Formula:
//
//	totalWeeks    = ceil(durationDays / 7)
//	billableWeeks = max(0, totalWeeks - freeWeeks)   ← first N weeks are free
//	totalFee%     = startupPct + (billableWeeks * weeklyIncrementPct)
//	totalFeeAmount = productValue * totalFee% / 100
func CalculateListingFee(tier SellerTier, productValue float64, durationDays int) FeeBreakdown {
	return CalculateListingFeeWithSettings(tier, productValue, durationDays, DefaultFeeSettings())
}

// CalculateListingFeeWithSettings is the same but uses caller-supplied settings (from DB).
func CalculateListingFeeWithSettings(tier SellerTier, productValue float64, durationDays int, s FeeSettings) FeeBreakdown {
	// Guard: negative inputs produce nonsensical fee results.
	if productValue < 0 {
		productValue = 0
	}
	if durationDays < 0 {
		durationDays = 0
	}

	var startupPct, weeklyPct float64
	switch tier {
	case TierGold:
		startupPct = s.GoldStartupFee
		weeklyPct = s.GoldWeeklyIncrement
	case TierPremium:
		startupPct = s.PremiumStartupFee
		weeklyPct = s.PremiumWeeklyIncrement
	default:
		startupPct = s.NormalStartupFee
		weeklyPct = s.NormalWeeklyIncrement
	}

	// Guard: negative fee rates from misconfigured DB must not invert the fee.
	if startupPct < 0 {
		startupPct = 0
	}
	if weeklyPct < 0 {
		weeklyPct = 0
	}

	totalWeeks := int(math.Ceil(float64(durationDays) / 7.0))
	freeWeeks := s.FreeWeeks
	if freeWeeks < 0 {
		freeWeeks = 0
	}
	billableWeeks := totalWeeks - freeWeeks
	if billableWeeks < 0 {
		billableWeeks = 0
	}

	totalPct := startupPct + float64(billableWeeks)*weeklyPct
	totalPct = math.Round(totalPct*100) / 100
	feeAmount := math.Round(productValue*totalPct/100*100) / 100
	sellerGets := math.Round((productValue-feeAmount)*100) / 100

	return FeeBreakdown{
		Tier:               string(tier),
		ProductValue:       productValue,
		DurationDays:       durationDays,
		StartupPct:         startupPct,
		WeeklyIncrementPct: weeklyPct,
		FreeWeeks:          freeWeeks,
		BillableWeeks:      billableWeeks,
		TotalFeePercent:    totalPct,
		TotalFeeAmount:     feeAmount,
		SellerGets:         sellerGets,
	}
}

// ResolveTier returns the tier for a seller based on their total value sold.
// goldThreshold and premiumThreshold must be loaded from platform_settings by the caller.
func ResolveTier(totalValueSold, goldThreshold, premiumThreshold float64) SellerTier {
	if goldThreshold <= 0 {
		goldThreshold = constants.DefaultGoldThreshold
	}
	if premiumThreshold <= 0 {
		premiumThreshold = constants.DefaultPremiumThreshold
	}
	// Guard: inverted thresholds (misconfigured DB) would silently assign wrong tiers.
	// Reset to defaults if the invariant goldThreshold < premiumThreshold is violated.
	if goldThreshold >= premiumThreshold {
		goldThreshold = constants.DefaultGoldThreshold
		premiumThreshold = constants.DefaultPremiumThreshold
	}
	switch {
	case totalValueSold >= premiumThreshold:
		return TierPremium
	case totalValueSold >= goldThreshold:
		return TierGold
	default:
		return TierNormal
	}
}

// BrokerFeeRate is the fallback platform fee charged per transaction.
// DEPRECATED: Use constants.DefaultBrokerFeeRate instead.
// The live rate is always loaded from platform_settings.broker_fee_rate.
const BrokerFeeRate = constants.DefaultBrokerFeeRate

// ReferralCommissionRate is the fallback referral commission rate (1.2%).
// The live rate is always loaded from platform_settings.referral_commission_rate.
// NOTE: stored as a percentage (1.2 = 1.2%), consistent with platform_settings.
const ReferralCommissionRate = constants.DefaultReferralCommissionRate

// CalculateFee returns broker fee and seller payout for a given amount using the default rate.
// Uses the live DB rate when available — this fallback is only for non-transaction contexts.
func CalculateFee(amount float64) (brokerFee, sellerGets float64) {
	return CalculateFeeWithRate(amount, BrokerFeeRate)
}

// CalculateFeeWithRate computes broker fee and seller payout.
// rate is a percentage value (e.g. 2.0 means 2%), consistent with platform_settings.broker_fee_rate.
// Pass rate=0 for a legitimate fee-free transaction. Pass rate<0 to use the default BrokerFeeRate.
func CalculateFeeWithRate(amount, rate float64) (brokerFee, sellerGets float64) {
	if amount < 0 {
		amount = 0
	}
	// rate=0 is valid (promotional free transaction); only negative triggers the fallback.
	if rate < 0 {
		rate = BrokerFeeRate
	}
	brokerFee = math.Round(amount*(rate/100)*100) / 100
	sellerGets = math.Round((amount-brokerFee)*100) / 100
	return
}

// CalculateReferralCommission returns the commission amount for a given base amount and rate.
// rate is a percentage value (e.g. 1.2 means 1.2%), consistent with platform_settings.referral_commission_rate.
// Pass rate=0 for no commission. Pass rate<0 to use the default ReferralCommissionRate.
func CalculateReferralCommission(baseAmount, rate float64) float64 {
	if baseAmount < 0 {
		baseAmount = 0
	}
	// rate=0 is valid (no commission); only negative triggers the fallback.
	if rate < 0 {
		rate = ReferralCommissionRate
	}
	return math.Round(baseAmount*(rate/100)*100) / 100
}
