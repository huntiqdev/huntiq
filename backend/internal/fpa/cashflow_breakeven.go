package fpa

import (
	"math"
)

// CashFlowAndBreakEvenResult combines liquidity and break-even metrics
type CashFlowAndBreakEvenResult struct {
	OperatingCashInflow     float64
	OperatingCashOutflow    float64
	TaxCashOutflow          float64
	NetOperatingCashFlow    float64
	CapExCashOutflow        float64
	NetMonthlyCashChange    float64
	ClosingCashBalance      float64
	MonthlyBurnRate         float64
	CashRunwayMonths        float64 // -1 = Cash flow positive / Infinite
	FixedCosts              float64
	VariableCosts           float64
	ContributionMargin      float64
	ContributionMarginRatio float64
	MonthlyBreakEvenRevenue float64
	IsBreakEven             bool
}

// CalculateCashFlowAndBreakEven models month-by-month cash movements, burn rate, runway, and break-even targets
func CalculateCashFlowAndBreakEven(
	monthIndex int,
	netRevenue float64,
	cashOpEx float64,
	infraCost float64,
	netVATPayable float64,
	annualCIT float64,
	initialCapEx float64,
	openingCash float64,
	marketingPerUserCAC float64,
	newUsers float64,
) CashFlowAndBreakEvenResult {
	// Cash Inflows: Net Revenue (assuming high payment gateway pass-through)
	operatingInflow := netRevenue

	// Cash Outflows: Cash OpEx + Hosting Infra
	operatingOutflow := cashOpEx + infraCost

	// Tax Cash Outflows:
	// VAT remitted monthly (M+1)
	taxCashOutflow := netVATPayable

	// Corporate Income Tax remitted annually in Month 4 of following fiscal cycle (Months 16, 28, 40, 52)
	if monthIndex == 16 || monthIndex == 28 || monthIndex == 40 || monthIndex == 52 {
		taxCashOutflow += annualCIT
	}

	netOperatingCash := operatingInflow - operatingOutflow - taxCashOutflow

	// CapEx Outflow: Incurred primarily at launch (Month 1)
	capexOutflow := 0.0
	if monthIndex == 1 && initialCapEx > 0 {
		capexOutflow = initialCapEx
	}

	netMonthlyCashChange := netOperatingCash - capexOutflow
	closingCash := openingCash + netMonthlyCashChange

	// Burn Rate: Monthly cash deficit
	monthlyBurn := 0.0
	if netMonthlyCashChange < 0 {
		monthlyBurn = math.Abs(netMonthlyCashChange)
	}

	// Cash Runway in Months
	cashRunway := -1.0
	if monthlyBurn > 0 {
		if closingCash > 0 {
			cashRunway = closingCash / monthlyBurn
		} else {
			cashRunway = 0.0 // Cash depleted
		}
	}

	// Break-Even Analysis (Contribution Margin Approach)
	// Variable costs: Direct user CAC, dynamic infra (storage, egress, SMS)
	variableCosts := (newUsers * marketingPerUserCAC) + (infraCost * 0.40)
	fixedCosts := cashOpEx + (infraCost * 0.60) - (newUsers * marketingPerUserCAC)
	if fixedCosts < 0 {
		fixedCosts = cashOpEx * 0.70
	}

	contributionMargin := netRevenue - variableCosts
	contributionMarginRatio := 0.0
	if netRevenue > 0 {
		contributionMarginRatio = contributionMargin / netRevenue
	}
	if contributionMarginRatio <= 0 {
		contributionMarginRatio = 0.35 // Conservative floor for marketplace platform
	}

	breakEvenRevenue := fixedCosts / contributionMarginRatio
	isBreakEven := netRevenue >= breakEvenRevenue && netRevenue > 0

	return CashFlowAndBreakEvenResult{
		OperatingCashInflow:     round2(operatingInflow),
		OperatingCashOutflow:    round2(operatingOutflow),
		TaxCashOutflow:          round2(taxCashOutflow),
		NetOperatingCashFlow:    round2(netOperatingCash),
		CapExCashOutflow:        round2(capexOutflow),
		NetMonthlyCashChange:    round2(netMonthlyCashChange),
		ClosingCashBalance:      round2(closingCash),
		MonthlyBurnRate:         round2(monthlyBurn),
		CashRunwayMonths:        round2(cashRunway),
		FixedCosts:              round2(fixedCosts),
		VariableCosts:           round2(variableCosts),
		ContributionMargin:      round2(contributionMargin),
		ContributionMarginRatio: round2(contributionMarginRatio),
		MonthlyBreakEvenRevenue: round2(breakEvenRevenue),
		IsBreakEven:             isBreakEven,
	}
}
