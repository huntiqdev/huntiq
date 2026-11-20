package fees

import (
	"github.com/your-org/used-tools-marketplace/pkg/constants"
	"math"
)

// FeeConfig holds fee calculation parameters
type FeeConfig struct {
	BrokerFeeRate    float64
	StartupFee       float64
	WeeklyIncrement  float64
	ListingFreeWeeks int
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

// CalculateListingFee calculates all fees for a product listing based on days
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

// CalculateListingFeeHourly calculates demand-based fees for hourly durations.
func CalculateListingFeeHourly(price float64, durationHours int, config FeeConfig) FeeResult {
	if durationHours <= 0 {
		durationHours = 24
	}

	hourlyPercentile := config.WeeklyIncrement / 168.0 // 168 hours in a week
	demandFactor := 1.0
	if durationHours <= 4 {
		demandFactor = 0.6 // Discounted flash rate
	} else if durationHours <= 24 {
		demandFactor = 0.85 // Urgent daily rate
	}

	billableHours := float64(durationHours - 24)
	if billableHours < 0 {
		billableHours = 0
	}

	totalFeePercent := config.StartupFee + billableHours*hourlyPercentile*demandFactor
	if totalFeePercent < 0.1 {
		totalFeePercent = 0.1
	}
	totalFeePercent = math.Round(totalFeePercent*100) / 100
	totalFeeAmount := math.Round(price*(totalFeePercent/100.0)*100) / 100

	brokerFeePercent := config.BrokerFeeRate
	brokerFeeAmount := math.Round(price*(brokerFeePercent/100.0)*100) / 100
	sellerPayout := math.Round((price-brokerFeeAmount)*100) / 100
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
