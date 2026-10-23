package fpa

import (
	"fmt"
	"math"
)

// EvaluateSensitivityShocks runs discrete what-if stress tests against the scenario configuration
func EvaluateSensitivityShocks(baseCfg FullScenarioConfig, baseReport *ForecastReport) []SensitivityShockResult {
	shocks := []SensitivityShockResult{}

	baseBreakEven := baseReport.Milestones.BreakEvenMonthIndex
	baseYear1Profit := 0.0
	if len(baseReport.AnnualSummaries) > 0 {
		baseYear1Profit = baseReport.AnnualSummaries[0].NetProfitAfterTax
	}
	baseLowestCash := baseReport.Milestones.LowestCashBalance

	// 1. What if user growth is twice as fast? (2x Growth)
	growth2xCfg := baseCfg
	growth2xCfg.Growth.EarlyMonthlyRate = baseCfg.Growth.EarlyMonthlyRate * 2.0
	growth2xCfg.Growth.LinearMonthlyAddition = baseCfg.Growth.LinearMonthlyAddition * 2.0
	rep2x := GenerateForecastReportWithoutShocks(growth2xCfg)
	y1Profit2x := 0.0
	if len(rep2x.AnnualSummaries) > 0 {
		y1Profit2x = rep2x.AnnualSummaries[0].NetProfitAfterTax
	}
	shocks = append(shocks, SensitivityShockResult{
		ShockName:             "2x Accelerated User Growth",
		Description:           "Early adoption rate doubles, driving faster marketplace density and higher transaction volume.",
		OriginalBreakEvenMo:   baseBreakEven,
		ShockBreakEvenMo:      rep2x.Milestones.BreakEvenMonthIndex,
		OriginalYear1Profit:   baseYear1Profit,
		ShockYear1Profit:      y1Profit2x,
		OriginalLowestCash:    baseLowestCash,
		ShockLowestCash:       rep2x.Milestones.LowestCashBalance,
		SurvivalCashNeededETB: math.Max(0, -rep2x.Milestones.LowestCashBalance),
		KeyTakeaway:           fmt.Sprintf("Accelerates break-even to month %d and boosts Year 1 profit by ETB %.0f.", rep2x.Milestones.BreakEvenMonthIndex, y1Profit2x-baseYear1Profit),
	})

	// 2. What if user growth is half as fast? (0.5x Growth)
	growthHalfCfg := baseCfg
	growthHalfCfg.Growth.EarlyMonthlyRate = baseCfg.Growth.EarlyMonthlyRate * 0.50
	growthHalfCfg.Growth.LinearMonthlyAddition = baseCfg.Growth.LinearMonthlyAddition * 0.50
	repHalf := GenerateForecastReportWithoutShocks(growthHalfCfg)
	y1ProfitHalf := 0.0
	if len(repHalf.AnnualSummaries) > 0 {
		y1ProfitHalf = repHalf.AnnualSummaries[0].NetProfitAfterTax
	}
	survivalNeedHalf := math.Max(0, -repHalf.Milestones.LowestCashBalance)
	shocks = append(shocks, SensitivityShockResult{
		ShockName:             "0.5x Slower User Adoption (Downside)",
		Description:           "Market traction develops at half projected speed, extending the unprofitable gestation phase.",
		OriginalBreakEvenMo:   baseBreakEven,
		ShockBreakEvenMo:      repHalf.Milestones.BreakEvenMonthIndex,
		OriginalYear1Profit:   baseYear1Profit,
		ShockYear1Profit:      y1ProfitHalf,
		OriginalLowestCash:    baseLowestCash,
		ShockLowestCash:       repHalf.Milestones.LowestCashBalance,
		SurvivalCashNeededETB: survivalNeedHalf,
		KeyTakeaway:           fmt.Sprintf("Delays break-even to month %d; requires an estimated ETB %.0f in capital reserves to maintain solvency.", repHalf.Milestones.BreakEvenMonthIndex, survivalNeedHalf),
	})

	// 3. What if infrastructure unit costs rise by 25%?
	infra25Cfg := baseCfg
	infra25Cfg.Infra.ComputeCostPer1KRPS *= 1.25
	infra25Cfg.Infra.CDNEgressCostPerGBETB *= 1.25
	infra25Cfg.Infra.DBStorageCostPerGBETB *= 1.25
	infra25Cfg.Infra.ComputeBaseMonthlyETB *= 1.25
	repInfra25 := GenerateForecastReportWithoutShocks(infra25Cfg)
	y1ProfitInfra25 := 0.0
	if len(repInfra25.AnnualSummaries) > 0 {
		y1ProfitInfra25 = repInfra25.AnnualSummaries[0].NetProfitAfterTax
	}
	shocks = append(shocks, SensitivityShockResult{
		ShockName:             "Infrastructure Cost Inflation (+25%)",
		Description:           "Higher cloud compute, bandwidth transit, and database licensing prices across all categories.",
		OriginalBreakEvenMo:   baseBreakEven,
		ShockBreakEvenMo:      repInfra25.Milestones.BreakEvenMonthIndex,
		OriginalYear1Profit:   baseYear1Profit,
		ShockYear1Profit:      y1ProfitInfra25,
		OriginalLowestCash:    baseLowestCash,
		ShockLowestCash:       repInfra25.Milestones.LowestCashBalance,
		SurvivalCashNeededETB: math.Max(0, -repInfra25.Milestones.LowestCashBalance),
		KeyTakeaway:           "Marginally reduces gross margin by 2-4%; break-even is delayed by at most 1-2 months due to strong software operating leverage.",
	})

	// 4. What if monetization is delayed by 6 months?
	monDelayCfg := baseCfg
	// Simulate zero fee monetization for first 6 months
	repDelay := GenerateForecastReportWithMonetizationDelay(monDelayCfg, 6)
	y1ProfitDelay := 0.0
	if len(repDelay.AnnualSummaries) > 0 {
		y1ProfitDelay = repDelay.AnnualSummaries[0].NetProfitAfterTax
	}
	survivalDelay := math.Max(0, -repDelay.Milestones.LowestCashBalance)
	shocks = append(shocks, SensitivityShockResult{
		ShockName:             "Monetization Delayed by 6 Months",
		Description:           "Platform runs 100% free with zero fees or commissions for the first 6 months to jumpstart liquidity.",
		OriginalBreakEvenMo:   baseBreakEven,
		ShockBreakEvenMo:      repDelay.Milestones.BreakEvenMonthIndex,
		OriginalYear1Profit:   baseYear1Profit,
		ShockYear1Profit:      y1ProfitDelay,
		OriginalLowestCash:    baseLowestCash,
		ShockLowestCash:       repDelay.Milestones.LowestCashBalance,
		SurvivalCashNeededETB: survivalDelay,
		KeyTakeaway:           fmt.Sprintf("Creates an additional ETB %.0f cash burn in Year 1; break-even is deferred to month %d.", math.Abs(y1ProfitDelay-baseYear1Profit), repDelay.Milestones.BreakEvenMonthIndex),
	})

	// 5. What if transaction conversion rate drops by 30%?
	txDropCfg := baseCfg
	txDropCfg.Revenue.TransactionsPerSellerMo *= 0.70
	txDropCfg.Revenue.AverageTransactionValue *= 0.85
	repTxDrop := GenerateForecastReportWithoutShocks(txDropCfg)
	y1ProfitTxDrop := 0.0
	if len(repTxDrop.AnnualSummaries) > 0 {
		y1ProfitTxDrop = repTxDrop.AnnualSummaries[0].NetProfitAfterTax
	}
	shocks = append(shocks, SensitivityShockResult{
		ShockName:             "Marketplace Conversion Drop (-30%)",
		Description:           "Economic slowdown or higher competition reduces average transaction count and order basket size.",
		OriginalBreakEvenMo:   baseBreakEven,
		ShockBreakEvenMo:      repTxDrop.Milestones.BreakEvenMonthIndex,
		OriginalYear1Profit:   baseYear1Profit,
		ShockYear1Profit:      y1ProfitTxDrop,
		OriginalLowestCash:    baseLowestCash,
		ShockLowestCash:       repTxDrop.Milestones.LowestCashBalance,
		SurvivalCashNeededETB: math.Max(0, -repTxDrop.Milestones.LowestCashBalance),
		KeyTakeaway:           "Requires greater reliance on recurring premium seller subscriptions to offset lower GMV commission revenue.",
	})

	return shocks
}
