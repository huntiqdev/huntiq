package fpa

import (
	"math"
)

// TaxCalculationResult details the Ethiopian statutory tax assessment for a month
type TaxCalculationResult struct {
	AccountingProfitBeforeTax float64
	TaxableIncomeAdjustments  float64
	TaxableIncomePreLoss      float64
	LossReliefApplied         float64
	EndingLossPool            float64
	NetTaxableIncome          float64
	EstimatedCorporateIncomeTax float64 // 30% standard rate
	OutputVAT                 float64 // 15% on taxable platform revenue
	EligibleInputVAT          float64 // Input credits on qualifying business expenses
	NetVATPayable             float64 // Output - Input VAT
	TotalStatutoryAccrual     float64 // CIT + VAT payable
	NetProfitAfterTax         float64 // Accounting PBT - CIT
	IsVATRegistered           bool
	RequiresAccountantReview  bool
}

// CalculateEthiopiaTax calculates statutory tax liabilities according to Proclamations 979/2016 and 285/2002 & 1341/2024
func CalculateEthiopiaTax(
	netRevenue, totalOpEx, cashOpEx, infraCost, officeSoftwareCost float64,
	annualizedTurnover float64,
	startingLossPool float64,
	cfg TaxConfig,
) TaxCalculationResult {
	pbt := netRevenue - totalOpEx

	// 1. Taxable Income Adjustments (Disallowed expenses under Proclamation 979/2016)
	disallowedRate := cfg.DisallowedExpenseRate
	if disallowedRate <= 0 {
		disallowedRate = 0.05 // 5% standard non-deductible expense assumption
	}
	disallowedAddback := cashOpEx * disallowedRate
	taxablePreLoss := pbt + disallowedAddback

	lossReliefApplied := 0.0
	endingLossPool := startingLossPool
	netTaxableIncome := 0.0

	if taxablePreLoss < 0 {
		// New tax loss: add to carryforward pool
		endingLossPool += math.Abs(taxablePreLoss)
		netTaxableIncome = 0.0
	} else {
		// Taxable profit: apply available loss carryforward relief (Art. 26 Proc. 979/2016)
		if endingLossPool > 0 {
			if taxablePreLoss >= endingLossPool {
				lossReliefApplied = endingLossPool
				netTaxableIncome = taxablePreLoss - endingLossPool
				endingLossPool = 0.0
			} else {
				lossReliefApplied = taxablePreLoss
				endingLossPool -= taxablePreLoss
				netTaxableIncome = 0.0
			}
		} else {
			netTaxableIncome = taxablePreLoss
		}
	}

	// 2. Corporate Income Tax (CIT) — 30% Standard Rate (Article 19 Proclamation 979/2016)
	citRate := cfg.CorporateIncomeTaxRate
	if citRate <= 0 {
		citRate = 0.30 // 30%
	}
	estimatedCIT := 0.0
	if netTaxableIncome > 0 {
		estimatedCIT = netTaxableIncome * citRate
	}

	// 3. Value Added Tax (VAT) — 15% Standard Rate (Article 7 Proclamations 285/2002 & 1341/2024)
	vatThreshold := cfg.VATRegistrationThreshold
	if vatThreshold <= 0 {
		vatThreshold = 1000000.0 // ETB 1M annual taxable supply threshold
	}
	isVATRegistered := annualizedTurnover >= vatThreshold

	vatRate := cfg.StandardVATRate
	if vatRate <= 0 {
		vatRate = 0.15 // 15%
	}

	outputVAT := 0.0
	eligibleInputVAT := 0.0
	netVATPayable := 0.0

	if isVATRegistered && netRevenue > 0 {
		outputVAT = netRevenue * vatRate

		// Input VAT recovery on business inputs (Cloud hosting, office rent, software SaaS)
		qualifyingInputs := infraCost + officeSoftwareCost
		recoverablePct := cfg.InputVATRecoverablePct
		if recoverablePct <= 0 {
			recoverablePct = 0.85 // 85% eligible input credit
		}
		eligibleInputVAT = qualifyingInputs * vatRate * recoverablePct
		netVATPayable = outputVAT - eligibleInputVAT
		// If input exceeds output, credit carries forward
		if netVATPayable < 0 {
			netVATPayable = 0
		}
	}

	totalStatutoryAccrual := estimatedCIT + netVATPayable
	netProfitAfterTax := pbt - estimatedCIT // Accounting Net Profit = PBT - Corporate Income Tax

	return TaxCalculationResult{
		AccountingProfitBeforeTax:   round2(pbt),
		TaxableIncomeAdjustments:    round2(disallowedAddback),
		TaxableIncomePreLoss:        round2(taxablePreLoss),
		LossReliefApplied:           round2(lossReliefApplied),
		EndingLossPool:              round2(endingLossPool),
		NetTaxableIncome:            round2(netTaxableIncome),
		EstimatedCorporateIncomeTax: round2(estimatedCIT),
		OutputVAT:                   round2(outputVAT),
		EligibleInputVAT:            round2(eligibleInputVAT),
		NetVATPayable:               round2(netVATPayable),
		TotalStatutoryAccrual:       round2(totalStatutoryAccrual),
		NetProfitAfterTax:           round2(netProfitAfterTax),
		IsVATRegistered:             isVATRegistered,
		RequiresAccountantReview:    true, // Mandatory disclaimer watermark
	}
}
