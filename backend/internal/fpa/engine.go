package fpa

import (
	"fmt"
	"math"
	"time"
)

// GenerateForecastReport runs the end-to-end financial planning, capacity, and tax calculation model
func GenerateForecastReport(cfg FullScenarioConfig) *ForecastReport {
	return generateReportInternal(cfg, 0, true)
}

// GenerateForecastReportWithoutShocks runs the calculation model without recursive sensitivity shocks
func GenerateForecastReportWithoutShocks(cfg FullScenarioConfig) *ForecastReport {
	return generateReportInternal(cfg, 0, false)
}

// GenerateForecastReportWithMonetizationDelay generates a forecast with delayed fee monetization
func GenerateForecastReportWithMonetizationDelay(cfg FullScenarioConfig, delayMonths int) *ForecastReport {
	return generateReportInternal(cfg, delayMonths, false)
}

func generateReportInternal(cfg FullScenarioConfig, delayMonths int, includeShocks bool) *ForecastReport {
	horizon := cfg.Growth.ForecastHorizonMonths
	if horizon <= 0 {
		horizon = 60
	}
	if horizon > 120 {
		horizon = 120
	}

	userGrowth := CalculateUserGrowth(cfg.Growth)
	monthlyRows := make([]MonthlyForecastRow, horizon)

	currentCash := cfg.CashFlow.InitialCashBalanceETB
	if currentCash <= 0 {
		currentCash = 10000000.0 // Default 10M ETB initial capital
	}
	lossPool := 0.0

	// Annual accumulators
	annualCITSum := 0.0
	var annualSummaries []AnnualSummary

	// Track milestones
	firstBreakEvenMonth := -1
	firstBreakEvenPeriod := "Not Achieved"
	firstBreakEvenRevenue := 0.0
	lowestCash := currentCash
	lowestCashPeriod := "Month 1"
	peakFundingNeed := 0.0

	baseDate := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)

	prevMAU := 0.0

	for m := 1; m <= horizon; m++ {
		monthDate := baseDate.AddDate(0, m-1, 0)
		periodStr := monthDate.Format("2006-01")

		rawUsers := userGrowth[m-1]
		mau, dau, cumulativeUsers := ConvertTargetMetric(rawUsers, cfg.Growth.TargetMetric, cfg.Workload.DAUtoMAURatio)

		newUsersThisMonth := mau - prevMAU
		if newUsersThisMonth < 0 {
			newUsersThisMonth = 0
		}
		prevMAU = mau

		// 1. Workload Demand
		workload := CalculateWorkload(dau, mau, cumulativeUsers, m, cfg.Workload)

		// 2. Revenue Modeling
		rev := CalculateRevenue(mau, cfg.Revenue)
		if m <= delayMonths {
			// Zero monetization during delay phase
			rev.CommissionRevenue = 0
			rev.ListingFeeRevenue = 0
			rev.SubscriptionRevenue = 0
			rev.PromotedListingRevenue = 0
			rev.BusinessToolsRevenue = 0
			rev.GrossBilledRevenue = 0
			rev.NetRecognizedRevenue = 0
		}

		// 3. Infrastructure Cost
		infra := CalculateInfraCost(workload, mau, m, cfg.Infra)

		// 4. Staffing & OpEx
		staffOpEx := CalculateStaffingAndOpEx(m, newUsersThisMonth, cfg.Staffing)

		// 5. Profitability Metrics
		totalOperatingExpenses := infra.TotalMonthlyCost + staffOpEx.TotalOpExWithDeprec
		grossProfit := rev.NetRecognizedRevenue - infra.TotalMonthlyCost
		grossMarginPct := 0.0
		if rev.NetRecognizedRevenue > 0 {
			grossMarginPct = (grossProfit / rev.NetRecognizedRevenue) * 100.0
		}

		ebitda := rev.NetRecognizedRevenue - (infra.TotalMonthlyCost + staffOpEx.TotalCashOpEx)
		ebitdaMarginPct := 0.0
		if rev.NetRecognizedRevenue > 0 {
			ebitdaMarginPct = (ebitda / rev.NetRecognizedRevenue) * 100.0
		}

		operatingProfit := rev.NetRecognizedRevenue - totalOperatingExpenses
		operatingMarginPct := 0.0
		if rev.NetRecognizedRevenue > 0 {
			operatingMarginPct = (operatingProfit / rev.NetRecognizedRevenue) * 100.0
		}

		// 6. Ethiopia Statutory Tax Engine
		annualizedTurnover := rev.NetRecognizedRevenue * 12.0
		taxResult := CalculateEthiopiaTax(
			rev.NetRecognizedRevenue,
			totalOperatingExpenses,
			staffOpEx.TotalCashOpEx,
			infra.TotalMonthlyCost,
			staffOpEx.OfficeRentAndAdmin+staffOpEx.SoftwareSaaS,
			annualizedTurnover,
			lossPool,
			cfg.Tax,
		)
		lossPool = taxResult.EndingLossPool
		annualCITSum += taxResult.EstimatedCorporateIncomeTax

		// 7. Cash Flow, Burn & Break-Even
		annualCITDue := 0.0
		if (m-1)%12 == 3 && m > 12 { // Remitted annually in month 4 of subsequent fiscal cycle
			annualCITDue = annualCITSum
			annualCITSum = 0.0
		}

		cf := CalculateCashFlowAndBreakEven(
			m,
			rev.NetRecognizedRevenue,
			staffOpEx.TotalCashOpEx,
			infra.TotalMonthlyCost,
			taxResult.NetVATPayable,
			annualCITDue,
			cfg.Staffing.InitialCapExETB,
			currentCash,
			cfg.Staffing.MarketingPerUserCAC,
			newUsersThisMonth,
		)
		currentCash = cf.ClosingCashBalance

		if currentCash < lowestCash {
			lowestCash = currentCash
			lowestCashPeriod = periodStr
		}
		if currentCash < 0 {
			fundingDeficit := math.Abs(currentCash)
			if fundingDeficit > peakFundingNeed {
				peakFundingNeed = fundingDeficit
			}
		}

		if cf.IsBreakEven && firstBreakEvenMonth == -1 {
			firstBreakEvenMonth = m
			firstBreakEvenPeriod = periodStr
			firstBreakEvenRevenue = cf.MonthlyBreakEvenRevenue
		}

		// 8. Unit Economics
		arpu := 0.0
		costPerUser := 0.0
		infraPerUser := 0.0
		infraPer1KReqs := 0.0
		infraPerTx := 0.0
		if mau > 0 {
			arpu = rev.NetRecognizedRevenue / mau
			costPerUser = totalOperatingExpenses / mau
			infraPerUser = infra.TotalMonthlyCost / mau
		}
		if workload.MonthlyRequests > 0 {
			infraPer1KReqs = (infra.TotalMonthlyCost / workload.MonthlyRequests) * 1000.0
		}
		if rev.CompletedTransactions > 0 {
			infraPerTx = infra.TotalMonthlyCost / float64(rev.CompletedTransactions)
		}

		monthlyRows[m-1] = MonthlyForecastRow{
			MonthIndex:                m,
			Period:                    periodStr,
			MAU:                       round2(mau),
			DAU:                       round2(dau),
			CumulativeUsers:           round2(cumulativeUsers),
			DailyRequests:             round2(workload.DailyRequests),
			AverageRPS:                round2(workload.AverageRPS),
			PeakRPS:                   round2(workload.PeakRPS),
			MonthlyRequests:           round2(workload.MonthlyRequests),
			DBReadIOPS:                round2(workload.DBReadIOPS),
			DBWriteIOPS:               round2(workload.DBWriteIOPS),
			AccumulatedStorageTB:      round2(workload.AccumulatedStorageTB),
			MonthlyBandwidthTB:        round2(workload.MonthlyBandwidthTB),
			GrossMerchandiseVolume:    round2(rev.GrossMerchandiseVolume),
			CompletedTransactions:     rev.CompletedTransactions,
			CommissionRevenue:         round2(rev.CommissionRevenue),
			ListingFeeRevenue:         round2(rev.ListingFeeRevenue),
			SubscriptionRevenue:       round2(rev.SubscriptionRevenue),
			PromotedListingRevenue:    round2(rev.PromotedListingRevenue),
			BusinessToolsRevenue:      round2(rev.BusinessToolsRevenue),
			GrossRecognizedRevenue:    round2(rev.GrossBilledRevenue),
			DiscountsAndRefunds:       round2(rev.DiscountsAndRefunds),
			NetRecognizedRevenue:      round2(rev.NetRecognizedRevenue),
			ComputeAndNetworkCost:     round2(infra.ComputeAndNetwork),
			DatabaseCost:              round2(infra.Database),
			StorageAndCDNCost:         round2(infra.StorageAndCDN),
			SecurityAndAuthCost:       round2(infra.SecurityAndAuth),
			TotalInfrastructureCost:   round2(infra.TotalMonthlyCost),
			StaffSalariesCost:         round2(staffOpEx.GrossSalaries),
			EmployerPensionCost:       round2(staffOpEx.EmployerPension),
			MarketingExpense:          round2(staffOpEx.MarketingExpense),
			OfficeAdminLegalCost:      round2(staffOpEx.OfficeRentAndAdmin + staffOpEx.LegalAndCompliance),
			SoftwareSaaSCost:          round2(staffOpEx.SoftwareSaaS),
			TotalCashOpEx:             round2(staffOpEx.TotalCashOpEx),
			AccountingDepreciation:    round2(staffOpEx.AccountingDepreciation),
			TotalOperatingExpense:     round2(totalOperatingExpenses),
			GrossProfit:               round2(grossProfit),
			GrossMarginPct:            round2(grossMarginPct),
			EBITDA:                    round2(ebitda),
			EBITDAMarginPct:           round2(ebitdaMarginPct),
			OperatingProfit:           round2(operatingProfit),
			OperatingMarginPct:        round2(operatingMarginPct),
			AccountingProfitBeforeTax: round2(taxResult.AccountingProfitBeforeTax),
			TaxableIncomeAdjustments:  round2(taxResult.TaxableIncomeAdjustments),
			TaxableIncome:             round2(taxResult.NetTaxableIncome),
			LossReliefApplied:         round2(taxResult.LossReliefApplied),
			CumulativeLossPool:        round2(taxResult.EndingLossPool),
			EstimatedIncomeTax:        round2(taxResult.EstimatedCorporateIncomeTax),
			OutputVAT:                 round2(taxResult.OutputVAT),
			EligibleInputVAT:          round2(taxResult.EligibleInputVAT),
			NetVATPayable:             round2(taxResult.NetVATPayable),
			TotalStatutoryTaxAccrual:  round2(taxResult.TotalStatutoryAccrual),
			NetProfitAfterTax:         round2(taxResult.NetProfitAfterTax),
			OperatingCashInflow:       round2(cf.OperatingCashInflow),
			OperatingCashOutflow:      round2(cf.OperatingCashOutflow),
			TaxCashOutflow:            round2(cf.TaxCashOutflow),
			NetOperatingCashFlow:      round2(cf.NetOperatingCashFlow),
			CapExCashOutflow:          round2(cf.CapExCashOutflow),
			NetMonthlyCashChange:      round2(cf.NetMonthlyCashChange),
			ClosingCashBalance:        round2(cf.ClosingCashBalance),
			MonthlyBurnRate:           round2(cf.MonthlyBurnRate),
			CashRunwayMonths:          round2(cf.CashRunwayMonths),
			RevenuePerActiveUser:      round2(arpu),
			CostPerActiveUser:         round2(costPerUser),
			InfraCostPer1KRequests:    round2(infraPer1KReqs),
			InfraCostPerActiveUser:    round2(infraPerUser),
			InfraCostPerTransaction:   round2(infraPerTx),
			EstimatedCAC:              round2(cfg.Staffing.MarketingPerUserCAC),
			ContributionMarginRatio:   round2(cf.ContributionMarginRatio),
			MonthlyBreakEvenRevenue:   round2(cf.MonthlyBreakEvenRevenue),
			IsBreakEven:               cf.IsBreakEven,
		}

		// Aggregate into annual summary blocks
		if m%12 == 0 || m == horizon {
			yearNum := (m - 1) / 12 + 1
			startIdx := (yearNum - 1) * 12
			endIdx := m

			annGMV := 0.0
			annGrossRev := 0.0
			annNetRev := 0.0
			annInfra := 0.0
			annStaff := 0.0
			annOpEx := 0.0
			annEBITDA := 0.0
			annEBIT := 0.0
			annCIT := 0.0
			annVAT := 0.0
			annNetProfit := 0.0
			annReqs := 0.0
			dauSum := 0.0

			for i := startIdx; i < endIdx; i++ {
				r := monthlyRows[i]
				annGMV += r.GrossMerchandiseVolume
				annGrossRev += r.GrossRecognizedRevenue
				annNetRev += r.NetRecognizedRevenue
				annInfra += r.TotalInfrastructureCost
				annStaff += r.StaffSalariesCost + r.EmployerPensionCost
				annOpEx += r.TotalOperatingExpense
				annEBITDA += r.EBITDA
				annEBIT += r.OperatingProfit
				annCIT += r.EstimatedIncomeTax
				annVAT += r.NetVATPayable
				annNetProfit += r.NetProfitAfterTax
				annReqs += r.MonthlyRequests
				dauSum += r.DAU
			}

			monthsInYear := float64(endIdx - startIdx)
			annualSummaries = append(annualSummaries, AnnualSummary{
				YearNumber:          yearNum,
				CalendarYear:        2026 + yearNum - 1,
				EndingMAU:           monthlyRows[endIdx-1].MAU,
				AverageDAU:          round2(dauSum / monthsInYear),
				TotalGMV:            round2(annGMV),
				GrossRevenue:        round2(annGrossRev),
				NetRevenue:          round2(annNetRev),
				TotalInfraCost:      round2(annInfra),
				TotalStaffingCost:   round2(annStaff),
				TotalOpEx:           round2(annOpEx),
				EBITDA:              round2(annEBITDA),
				OperatingProfit:     round2(annEBIT),
				EstimatedCIT:        round2(annCIT),
				NetVATPayable:       round2(annVAT),
				NetProfitAfterTax:   round2(annNetProfit),
				EndingCashBalance:   monthlyRows[endIdx-1].ClosingCashBalance,
				TotalRequestsAnnual: round2(annReqs),
				StorageEndingTB:     monthlyRows[endIdx-1].AccumulatedStorageTB,
				AverageARPU:         round2(annNetRev / (monthlyRows[endIdx-1].MAU + 1.0)),
				AverageCostPerUser:  round2(annOpEx / (monthlyRows[endIdx-1].MAU + 1.0)),
			})
		}
	}

	m12MAU := 0.0
	m12NetRev := 0.0
	if horizon >= 12 {
		m12MAU = monthlyRows[11].MAU
		m12NetRev = monthlyRows[11].NetRecognizedRevenue
	}
	m60MAU := monthlyRows[horizon-1].MAU
	m60NetRev := monthlyRows[horizon-1].NetRecognizedRevenue

	milestones := MilestoneAnalysis{
		BreakEvenMonthIndex: firstBreakEvenMonth,
		BreakEvenPeriod:     firstBreakEvenPeriod,
		BreakEvenRevenue:    round2(firstBreakEvenRevenue),
		PeakFundingNeedETB:  round2(peakFundingNeed),
		LowestCashMonth:     lowestCashPeriod,
		LowestCashBalance:   round2(lowestCash),
		Month12MAU:          round2(m12MAU),
		Month12NetRevenue:   round2(m12NetRev),
		Month60MAU:          round2(m60MAU),
		Month60NetRevenue:   round2(m60NetRev),
	}

	capacity20M := Assess20MCapacityScenario(cfg.Workload, cfg.Infra)

	report := &ForecastReport{
		ScenarioID:       cfg.ScenarioID,
		ScenarioName:     cfg.Name,
		Currency:         cfg.Currency,
		GeneratedAt:      time.Now().UTC(),
		HorizonMonths:    horizon,
		DisclaimerNotice: "ILLUSTRATIVE FINANCIAL PLANNING MODEL ONLY — NOT CERTIFIED ACCOUNTING RECORDS. Statutory tax rates require validation by an authorized Ethiopian tax consultant under Proclamations 979/2016 and 285/2002 & 1341/2024.",
		Milestones:       milestones,
		Capacity20M:      capacity20M,
		AnnualSummaries:  annualSummaries,
		MonthlySeries:    monthlyRows,
	}

	if includeShocks {
		report.SensitivityShocks = EvaluateSensitivityShocks(cfg, report)
	}

	return report
}

// GetBaselineScenarioConfig returns the standard enterprise base-case configuration
func GetBaselineScenarioConfig(scenarioID string) FullScenarioConfig {
	currency := CurrencyETB

	switch scenarioID {
	case "conservative":
		return FullScenarioConfig{
			ScenarioID:  "conservative",
			Name:        "Conservative Downside Case",
			Description: "Slower adoption, lower monetization conversion, higher churn, and protective cash runway modeling.",
			Currency:    currency,
			Growth: GrowthConfig{
				ModelType:             GrowthStableTransition,
				TargetMetric:          MetricMAU,
				InitialUsers:          800.0,
				EarlyMonthlyRate:      0.10, // 10% early growth
				EarlyDurationMonths:   12,
				MatureMonthlyRate:     0.015, // 1.5% mature
				MarketCeiling:         10000000.0,
				ForecastHorizonMonths: 60,
			},
			Workload: WorkloadConfig{
				DAUtoMAURatio:             0.20,
				SessionsPerUserPerDay:     2.0,
				RequestsPerSession:        15.0,
				PeakToAverageRatio:        3.5,
				SearchesPerUserDay:        3.0,
				ListingViewsPerUserDay:    10.0,
				ImageUploadsPerActiveUser: 0.03,
				ImageViewsPerUserDay:      18.0,
				AvgImageSizeKB:            350.0,
				DBReadsPerRequest:         3.5,
				DBWritesPerRequest:        0.4,
				DataRetainedPerUserKB:     50.0,
				DataRetainedPerListingKB:  250.0,
				LogsRetainedDays:          90,
				BackgroundJobsPerUserDay:  0.15,
			},
			Infra: InfraCostConfig{
				Profile:                ProfileLeanEarlyStage,
				ComputeBaseMonthlyETB:  4000.0,
				ComputeCostPer1KRPS:    250.0,
				DBBaseMonthlyETB:       5500.0,
				DBStorageCostPerGBETB:  18.0,
				ObjectStorageCostPerGB: 3.50,
				CDNEgressCostPerGBETB:  8.50,
				LoggingMonitoringETB:   2000.0,
				SecurityWAFMonthlyETB:  2500.0,
				AuthEmailSMSUnitETB:    1.20,
				BackupDisasterRecETB:   1800.0,
				RedundancyBufferPct:    0.20,
			},
			Staffing: StaffingExpenseConfig{
				Roles: []StaffRoleConfig{
					{RoleTitle: "Backend Engineers (Go)", InitialCount: 2, TargetCountM12: 3, TargetCountM60: 6, MonthlySalaryETB: 50000.0},
					{RoleTitle: "Frontend Engineers (Next.js)", InitialCount: 1, TargetCountM12: 2, TargetCountM60: 4, MonthlySalaryETB: 45000.0},
					{RoleTitle: "DevOps & SRE", InitialCount: 1, TargetCountM12: 1, TargetCountM60: 2, MonthlySalaryETB: 55000.0},
					{RoleTitle: "Product & Operations", InitialCount: 1, TargetCountM12: 2, TargetCountM60: 4, MonthlySalaryETB: 40000.0},
					{RoleTitle: "Customer Support & Triage", InitialCount: 1, TargetCountM12: 2, TargetCountM60: 6, MonthlySalaryETB: 25000.0},
				},
				EmployerPensionRate: 0.11,
				AnnualSalaryRaise:   0.08,
				OfficeRentMonthly:   35000.0,
				LegalComplianceMo:   10000.0,
				MarketingBudgetMo:   20000.0,
				MarketingPerUserCAC: 20.0,
				SoftwareSaaSMo:      12000.0,
				InsuranceMiscMo:     7000.0,
				InitialCapExETB:     350000.0,
				CapExDeprecMonths:   36,
			},
			Revenue: RevenueConfig{
				SellersPctOfUsers:              0.12,
				TransactionsPerSellerMo:        2.5,
				AverageTransactionValue:        3800.0,
				CommissionRatePct:              0.02,
				RefundCancellationRate:         0.04,
				CommissionDiscountRate:         0.05,
				ListingsPerSellerMo:            3.0,
				PaidListingFeePerItem:          45.0,
				PaidListingAdoptionPct:         0.15,
				GoldTierConversionPct:          0.015,
				GoldTierMonthlyPrice:           500.0,
				PremiumTierConversionPct:       0.005,
				PremiumTierMonthlyPrice:        1500.0,
				TierMonthlyChurnRate:           0.06,
				FeaturedListingMonthlyAdoption: 0.02,
				FeaturedListingPrice:           200.0,
				VerifiedBusinessSellersPct:     0.01,
				VerifiedBusinessToolFeeMo:      750.0,
			},
			Tax: TaxConfig{
				Jurisdiction:             "ET",
				EntityLegalType:          "Private Limited Company (PLC)",
				FiscalYearConvention:     "Gregorian",
				CorporateIncomeTaxRate:   0.30,
				StandardVATRate:          0.15,
				VATRegistrationThreshold: 1000000.0,
				InputVATRecoverablePct:   0.85,
				DisallowedExpenseRate:    0.05,
				LossCarryforwardYears:    5,
				RequiresAccountantReview: true,
			},
			CashFlow: CashFlowConfig{
				InitialCashBalanceETB: 8000000.0,
			},
		}

	case "high_growth":
		return FullScenarioConfig{
			ScenarioID:  "high_growth",
			Name:        "High Growth Upside Case",
			Description: "Aggressive marketplace expansion, rapid active user scaling, higher tier adoption, and capacity stress testing.",
			Currency:    currency,
			Growth: GrowthConfig{
				ModelType:             GrowthStableTransition,
				TargetMetric:          MetricMAU,
				InitialUsers:          1500.0,
				EarlyMonthlyRate:      0.28, // 28% early monthly growth
				EarlyDurationMonths:   14,
				MatureMonthlyRate:     0.035, // 3.5% mature
				MarketCeiling:         20000000.0,
				ForecastHorizonMonths: 60,
			},
			Workload: WorkloadConfig{
				DAUtoMAURatio:             0.30,
				SessionsPerUserPerDay:     3.0,
				RequestsPerSession:        22.0,
				PeakToAverageRatio:        3.0,
				SearchesPerUserDay:        5.0,
				ListingViewsPerUserDay:    15.0,
				ImageUploadsPerActiveUser: 0.08,
				ImageViewsPerUserDay:      30.0,
				AvgImageSizeKB:            350.0,
				DBReadsPerRequest:         3.8,
				DBWritesPerRequest:        0.5,
				DataRetainedPerUserKB:     50.0,
				DataRetainedPerListingKB:  250.0,
				LogsRetainedDays:          90,
				BackgroundJobsPerUserDay:  0.25,
			},
			Infra: InfraCostConfig{
				Profile:                ProfileHighAvailability,
				ComputeBaseMonthlyETB:  7500.0,
				ComputeCostPer1KRPS:    250.0,
				DBBaseMonthlyETB:       9000.0,
				DBStorageCostPerGBETB:  18.0,
				ObjectStorageCostPerGB: 3.50,
				CDNEgressCostPerGBETB:  8.50,
				LoggingMonitoringETB:   3500.0,
				SecurityWAFMonthlyETB:  4500.0,
				AuthEmailSMSUnitETB:    1.20,
				BackupDisasterRecETB:   3000.0,
				RedundancyBufferPct:    0.25,
			},
			Staffing: StaffingExpenseConfig{
				Roles: []StaffRoleConfig{
					{RoleTitle: "Backend Engineers (Go)", InitialCount: 3, TargetCountM12: 6, TargetCountM60: 16, MonthlySalaryETB: 55000.0},
					{RoleTitle: "Frontend Engineers (Next.js)", InitialCount: 2, TargetCountM12: 4, TargetCountM60: 12, MonthlySalaryETB: 50000.0},
					{RoleTitle: "DevOps & SRE", InitialCount: 1, TargetCountM12: 3, TargetCountM60: 6, MonthlySalaryETB: 60000.0},
					{RoleTitle: "Security & QA", InitialCount: 1, TargetCountM12: 2, TargetCountM60: 5, MonthlySalaryETB: 50000.0},
					{RoleTitle: "Product & Operations", InitialCount: 2, TargetCountM12: 4, TargetCountM60: 10, MonthlySalaryETB: 45000.0},
					{RoleTitle: "Customer Support & Moderation", InitialCount: 2, TargetCountM12: 6, TargetCountM60: 20, MonthlySalaryETB: 28000.0},
					{RoleTitle: "Sales & Marketing", InitialCount: 2, TargetCountM12: 5, TargetCountM60: 12, MonthlySalaryETB: 40000.0},
					{RoleTitle: "Finance & Compliance", InitialCount: 1, TargetCountM12: 2, TargetCountM60: 4, MonthlySalaryETB: 50000.0},
				},
				EmployerPensionRate: 0.11,
				AnnualSalaryRaise:   0.08,
				OfficeRentMonthly:   60000.0,
				LegalComplianceMo:   18000.0,
				MarketingBudgetMo:   50000.0,
				MarketingPerUserCAC: 12.0,
				SoftwareSaaSMo:      25000.0,
				InsuranceMiscMo:     10000.0,
				InitialCapExETB:     750000.0,
				CapExDeprecMonths:   36,
			},
			Revenue: RevenueConfig{
				SellersPctOfUsers:              0.18,
				TransactionsPerSellerMo:        4.2,
				AverageTransactionValue:        5200.0,
				CommissionRatePct:              0.02,
				RefundCancellationRate:         0.025,
				CommissionDiscountRate:         0.04,
				ListingsPerSellerMo:            5.0,
				PaidListingFeePerItem:          55.0,
				PaidListingAdoptionPct:         0.25,
				GoldTierConversionPct:          0.03,
				GoldTierMonthlyPrice:           500.0,
				PremiumTierConversionPct:       0.012,
				PremiumTierMonthlyPrice:        1500.0,
				TierMonthlyChurnRate:           0.03,
				FeaturedListingMonthlyAdoption: 0.045,
				FeaturedListingPrice:           300.0,
				VerifiedBusinessSellersPct:     0.025,
				VerifiedBusinessToolFeeMo:      900.0,
			},
			Tax: TaxConfig{
				Jurisdiction:             "ET",
				EntityLegalType:          "Private Limited Company (PLC)",
				FiscalYearConvention:     "Gregorian",
				CorporateIncomeTaxRate:   0.30,
				StandardVATRate:          0.15,
				VATRegistrationThreshold: 1000000.0,
				InputVATRecoverablePct:   0.85,
				DisallowedExpenseRate:    0.05,
				LossCarryforwardYears:    5,
				RequiresAccountantReview: true,
			},
			CashFlow: CashFlowConfig{
				InitialCashBalanceETB: 15000000.0,
			},
		}

	case "base":
		fallthrough
	default:
		return FullScenarioConfig{
			ScenarioID:  "base",
			Name:        "Base Case Planning",
			Description: "Realistic adoption growth, standard infrastructure deployment, benchmark monetization rates.",
			Currency:    currency,
			Growth: GrowthConfig{
				ModelType:             GrowthStableTransition,
				TargetMetric:          MetricMAU,
				InitialUsers:          1000.0,
				EarlyMonthlyRate:      0.20, // 20% early growth
				EarlyDurationMonths:   12,
				MatureMonthlyRate:     0.02, // 2% mature growth
				MarketCeiling:         20000000.0,
				ForecastHorizonMonths: 60,
			},
			Workload: WorkloadConfig{
				DAUtoMAURatio:             0.25,
				SessionsPerUserPerDay:     2.5,
				RequestsPerSession:        18.0,
				PeakToAverageRatio:        3.0,
				SearchesPerUserDay:        4.0,
				ListingViewsPerUserDay:    12.0,
				ImageUploadsPerActiveUser: 0.05,
				ImageViewsPerUserDay:      24.0,
				AvgImageSizeKB:            350.0,
				DBReadsPerRequest:         3.5,
				DBWritesPerRequest:        0.4,
				DataRetainedPerUserKB:     50.0,
				DataRetainedPerListingKB:  250.0,
				LogsRetainedDays:          90,
				BackgroundJobsPerUserDay:  0.20,
			},
			Infra: InfraCostConfig{
				Profile:                ProfileGrowingProduction,
				ComputeBaseMonthlyETB:  5000.0,
				ComputeCostPer1KRPS:    250.0,
				DBBaseMonthlyETB:       7000.0,
				DBStorageCostPerGBETB:  18.0,
				ObjectStorageCostPerGB: 3.50,
				CDNEgressCostPerGBETB:  8.50,
				LoggingMonitoringETB:   2500.0,
				SecurityWAFMonthlyETB:  3000.0,
				AuthEmailSMSUnitETB:    1.20,
				BackupDisasterRecETB:   2200.0,
				RedundancyBufferPct:    0.20,
			},
			Staffing: StaffingExpenseConfig{
				Roles: []StaffRoleConfig{
					{RoleTitle: "Backend Engineers (Go)", InitialCount: 2, TargetCountM12: 4, TargetCountM60: 10, MonthlySalaryETB: 50000.0},
					{RoleTitle: "Frontend Engineers (Next.js)", InitialCount: 1, TargetCountM12: 3, TargetCountM60: 8, MonthlySalaryETB: 45000.0},
					{RoleTitle: "DevOps & SRE", InitialCount: 1, TargetCountM12: 2, TargetCountM60: 4, MonthlySalaryETB: 55000.0},
					{RoleTitle: "Security & QA", InitialCount: 0, TargetCountM12: 1, TargetCountM60: 3, MonthlySalaryETB: 48000.0},
					{RoleTitle: "Product & Operations", InitialCount: 1, TargetCountM12: 3, TargetCountM60: 6, MonthlySalaryETB: 42000.0},
					{RoleTitle: "Customer Support & Triage", InitialCount: 1, TargetCountM12: 4, TargetCountM60: 12, MonthlySalaryETB: 26000.0},
					{RoleTitle: "Marketing & Growth", InitialCount: 1, TargetCountM12: 3, TargetCountM60: 6, MonthlySalaryETB: 38000.0},
					{RoleTitle: "Finance & Legal", InitialCount: 1, TargetCountM12: 1, TargetCountM60: 3, MonthlySalaryETB: 45000.0},
				},
				EmployerPensionRate: 0.11,
				AnnualSalaryRaise:   0.08,
				OfficeRentMonthly:   45000.0,
				LegalComplianceMo:   12000.0,
				MarketingBudgetMo:   30000.0,
				MarketingPerUserCAC: 15.0,
				SoftwareSaaSMo:      16000.0,
				InsuranceMiscMo:     8000.0,
				InitialCapExETB:     500000.0,
				CapExDeprecMonths:   36,
			},
			Revenue: RevenueConfig{
				SellersPctOfUsers:              0.15,
				TransactionsPerSellerMo:        3.2,
				AverageTransactionValue:        4500.0,
				CommissionRatePct:              0.02,
				RefundCancellationRate:         0.03,
				CommissionDiscountRate:         0.05,
				ListingsPerSellerMo:            4.0,
				PaidListingFeePerItem:          50.0,
				PaidListingAdoptionPct:         0.20,
				GoldTierConversionPct:          0.02,
				GoldTierMonthlyPrice:           500.0,
				PremiumTierConversionPct:       0.008,
				PremiumTierMonthlyPrice:        1500.0,
				TierMonthlyChurnRate:           0.04,
				FeaturedListingMonthlyAdoption: 0.03,
				FeaturedListingPrice:           250.0,
				VerifiedBusinessSellersPct:     0.015,
				VerifiedBusinessToolFeeMo:      800.0,
			},
			Tax: TaxConfig{
				Jurisdiction:             "ET",
				EntityLegalType:          "Private Limited Company (PLC)",
				FiscalYearConvention:     "Gregorian",
				CorporateIncomeTaxRate:   0.30,
				StandardVATRate:          0.15,
				VATRegistrationThreshold: 1000000.0,
				InputVATRecoverablePct:   0.85,
				DisallowedExpenseRate:    0.05,
				LossCarryforwardYears:    5,
				RequiresAccountantReview: true,
			},
			CashFlow: CashFlowConfig{
				InitialCashBalanceETB: 10000000.0,
			},
		}
	}
}

// GenerateCSVExport generates an RFC-4180 compliant CSV export of the monthly financial forecast
func GenerateCSVExport(rep *ForecastReport) string {
	csv := "Month,Period,MAU,DAU,GrossMerchandiseVolume,NetRevenue,TotalInfraCost,TotalOpEx,EBITDA,OperatingProfit,EstimatedCIT,NetVATPayable,NetProfitAfterTax,ClosingCashBalance,MonthlyBurnRate,IsBreakEven\n"
	for _, r := range rep.MonthlySeries {
		csv += fmt.Sprintf(
			"%d,%s,%.0f,%.0f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%t\n",
			r.MonthIndex, r.Period, r.MAU, r.DAU,
			r.GrossMerchandiseVolume, r.NetRecognizedRevenue, r.TotalInfrastructureCost,
			r.TotalOperatingExpense, r.EBITDA, r.OperatingProfit, r.EstimatedIncomeTax,
			r.NetVATPayable, r.NetProfitAfterTax, r.ClosingCashBalance, r.MonthlyBurnRate, r.IsBreakEven,
		)
	}
	return csv
}
