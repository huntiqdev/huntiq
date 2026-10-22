package fees_test

import (
	"testing"

	"github.com/your-org/used-tools-marketplace/pkg/fees"
)

func TestCalculateListingFee(t *testing.T) {
	cfg := fees.FeeConfig{
		BrokerFeeRate:    5.0,  // 5% broker fee
		StartupFee:       2.0,  // 2% startup fee
		WeeklyIncrement:  1.0,  // 1% per billable week
		ListingFreeWeeks: 2,    // 2 weeks free
	}

	tests := []struct {
		name                 string
		price                float64
		durationDays         int
		expectedFeePercent   float64
		expectedBrokerAmount float64
		expectedPayout       float64
	}{
		{
			name:                 "14 days (2 weeks free)",
			price:                1000.0,
			durationDays:         14,
			expectedFeePercent:   2.0, // startup only, 0 billable extra weeks
			expectedBrokerAmount: 50.0,
			expectedPayout:       950.0,
		},
		{
			name:                 "21 days (3 weeks = 1 billable week)",
			price:                2000.0,
			durationDays:         21,
			expectedFeePercent:   3.0, // 2% startup + 1% billable week
			expectedBrokerAmount: 100.0,
			expectedPayout:       1900.0,
		},
		{
			name:                 "0 days",
			price:                500.0,
			durationDays:         0,
			expectedFeePercent:   2.0,
			expectedBrokerAmount: 25.0,
			expectedPayout:       475.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := fees.CalculateListingFee(tt.price, tt.durationDays, cfg)
			if res.TotalFeePercent != tt.expectedFeePercent {
				t.Errorf("expected fee percent %f, got %f", tt.expectedFeePercent, res.TotalFeePercent)
			}
			if res.BrokerFeeAmount != tt.expectedBrokerAmount {
				t.Errorf("expected broker amount %f, got %f", tt.expectedBrokerAmount, res.BrokerFeeAmount)
			}
			if res.SellerPayout != tt.expectedPayout {
				t.Errorf("expected seller payout %f, got %f", tt.expectedPayout, res.SellerPayout)
			}
		})
	}
}

func TestAmountsMatch(t *testing.T) {
	if !fees.AmountsMatch(1000.0, 1000.5) {
		t.Errorf("expected 1000.0 and 1000.5 to match within 0.1%% tolerance")
	}
	if fees.AmountsMatch(1000.0, 1005.0) {
		t.Errorf("expected 1000.0 and 1005.0 not to match outside 0.1%% tolerance")
	}
}
