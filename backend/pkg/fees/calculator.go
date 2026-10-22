package fees

import (
	"math"
	"github.com/your-org/used-tools-marketplace/pkg/constants"
)

// FeeConfig holds fee calculation parameters
type FeeConfig struct {
	BrokerFeeRate     float64
	StartupFee        float64
	WeeklyIncrement   float64
	ListingFreeWeeks  int
}

// FeeResult contains calculated fee breakdown
type FeeResult struct {
	TotalFeePercent  float64 // Total listing fee as percentage
	TotalFeeAmount   float64 // Listing fee in ETB
	BrokerFeePercent float64 // Broker fee as percentage
	BrokerFeeAmount  float64 // Broker fee in ETB
	SellerPayout     float64 // Amount seller receives
	TotalCost        float64 // Buyer pays (price + fees)
}

// CalculateListingFee calculates all fees for a product listing
func CalculateListingFee(price float64, durationDays int, config FeeConfig) FeeResult {
	// Convert days to weeks (round up)
	totalWeeks := int(math.Ceil(float64(durationDays) / 7.0))
	
	// Calculate billable weeks (subtract free weeks)
	billableWeeks := totalWeeks - config.ListingFreeWeeks
	if billableWeeks < 0 {
		billableWeeks = 0
	}

	// Calculate total listing fee percentage
	totalFeePercent := config.StartupFee + float64(billableWeeks)*config.WeeklyIncrement
	totalFeeAmount := math.Round(price*(totalFeePercent/100.0)*100) / 100

	// Calculate broker fee
	brokerFeePercent := config.BrokerFeeRate
	brokerFeeAmount := math.Round(price*(brokerFeePercent/100.0)*100) / 100

	// Calculate seller payout
	sellerPayout := math.Round((price-brokerFeeAmount)*100) / 100

	// Total cost for buyer (price + listing fee)
	totalCost := math.Round((price+totalFeeAmount)*100) / 100

	return FeeResult{
		TotalFeePercent:  totalFeePercent,
		TotalFeeAmount:   totalFeeAmount,
		BrokerFeePercent: brokerFeePercent,
		BrokerFeeAmount:  brokerFeeAmount,
		SellerPayout:     sellerPayout,
		TotalCost:        totalCost,
	}
}

// AmountsMatch checks if two currency amounts match within percentage-based tolerance.
// Uses constants.DefaultCurrencyTolerancePercent (0.1%) expressed as a fraction.
func AmountsMatch(amount1, amount2 float64) bool {
	tolerancePercent := constants.DefaultCurrencyTolerancePercent / 100.0 // 0.1% → 0.001
	maxDiff := math.Max(amount1, amount2) * tolerancePercent
	diff := math.Abs(amount1 - amount2)
	return diff <= maxDiff
}
