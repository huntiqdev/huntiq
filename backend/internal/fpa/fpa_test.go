package fpa

import (
	"math"
	"strings"
	"testing"
)

func TestUserGrowthModels(t *testing.T) {
	// 1. Geometric Growth
	geoCfg := GrowthConfig{
		ModelType:             GrowthGeometric,
		InitialUsers:          1000,
		EarlyMonthlyRate:      0.10,
		ForecastHorizonMonths: 12,
		MarketCeiling:         100000,
	}
	geo := CalculateUserGrowth(geoCfg)
	if len(geo) != 12 {
		t.Fatalf("expected 12 months, got %d", len(geo))
	}
	if geo[0] != 1000 {
		t.Errorf("expected month 1 = 1000, got %f", geo[0])
	}
	if geo[1] <= geo[0] {
		t.Errorf("expected growth in month 2, got %f", geo[1])
	}

	// 2. Linear Growth
	linCfg := GrowthConfig{
		ModelType:             GrowthLinear,
		InitialUsers:          500,
		LinearMonthlyAddition: 250,
		ForecastHorizonMonths: 6,
		MarketCeiling:         50000,
	}
	lin := CalculateUserGrowth(linCfg)
	if lin[1] != 750 {
		t.Errorf("expected month 2 = 750, got %f", lin[1])
	}
	if lin[5] != 1750 {
		t.Errorf("expected month 6 = 1750, got %f", lin[5])
	}

	// 3. Logistic Saturation (clamping at ceiling)
	logCfg := GrowthConfig{
		ModelType:             GrowthLogistic,
		InitialUsers:          10000,
		EarlyMonthlyRate:      0.50,
		MarketCeiling:         15000,
		ForecastHorizonMonths: 24,
	}
	logSeries := CalculateUserGrowth(logCfg)
	lastVal := logSeries[len(logSeries)-1]
	if lastVal > 15000.0 {
		t.Errorf("logistic growth exceeded market ceiling: got %f, max 15000", lastVal)
	}

	// 4. Stable Transition Model
	transCfg := GrowthConfig{
		ModelType:             GrowthStableTransition,
		InitialUsers:          1000,
		EarlyMonthlyRate:      0.20,
		EarlyDurationMonths:   6,
		MatureMonthlyRate:     0.02,
		MarketCeiling:         1000000,
		ForecastHorizonMonths: 12,
	}
	trans := CalculateUserGrowth(transCfg)
	if trans[6] <= trans[5] {
		t.Errorf("expected continuous growth, got m6=%f, m7=%f", trans[5], trans[6])
	}

	// 5. Quarterly Smooth Growth Model
	qSmoothCfg := GrowthConfig{
		ModelType:             GrowthQuarterlySmooth,
		InitialUsers:          1000,
		EarlyQuarterlyRate:    0.20, // 20% quarterly
		EarlyDurationQuarters: 4,    // 4 quarters = 12 months
		MatureMonthlyRate:     0.02,
		MarketCeiling:         1000000,
		ForecastHorizonMonths: 13,
	}
	qSmooth := CalculateUserGrowth(qSmoothCfg)
	if len(qSmooth) != 13 {
		t.Fatalf("expected 13 months, got %d", len(qSmooth))
	}
	// At month 4 (after 3 months = 1 quarter), 1000 * (1.20)^1 should be ~1200
	if math.Abs(qSmooth[3]-1200.0) > 0.01 {
		t.Errorf("expected Quarter 1 end (month 4) = 1200.0, got %f", qSmooth[3])
	}
	// At month 13 (after 12 months = 4 quarters), 1000 * (1.20)^4 should be ~2073.6
	if math.Abs(qSmooth[12]-2073.6) > 0.1 {
		t.Errorf("expected Quarter 4 end (month 13) = 2073.6, got %f", qSmooth[12])
	}

	// 6. Quarterly Step Growth Model
	qStepCfg := GrowthConfig{
		ModelType:             GrowthQuarterlyStep,
		InitialUsers:          1000,
		EarlyQuarterlyRate:    0.20,
		EarlyDurationQuarters: 4,
		ForecastHorizonMonths: 12,
		MarketCeiling:         1000000,
	}
	qStep := CalculateUserGrowth(qStepCfg)
	// Months 1, 2, 3 should stay at 1000
	if qStep[0] != 1000 || qStep[1] != 1000 || qStep[2] != 1000 {
		t.Errorf("expected months 1-3 to remain 1000, got %f, %f, %f", qStep[0], qStep[1], qStep[2])
	}
	// Month 4 should step to 1200
	if qStep[3] != 1200 {
		t.Errorf("expected month 4 to step to 1200, got %f", qStep[3])
	}
}

func TestWorkloadCalculations(t *testing.T) {
	cfg := WorkloadConfig{
		DAUtoMAURatio:            0.25,
		SessionsPerUserPerDay:    2.0,
		RequestsPerSession:       10.0,
		PeakToAverageRatio:       3.0,
		AvgImageSizeKB:           300.0,
		DBReadsPerRequest:        3.0,
		DBWritesPerRequest:       0.5,
		BackgroundJobsPerUserDay: 0.1,
	}

	dau := 10000.0
	mau := 40000.0
	metrics := CalculateWorkload(dau, mau, 100000.0, 12, cfg)

	expectedDailyReqs := dau * 2.0 * 10.0 // 200,000 req/day
	if metrics.DailyRequests != expectedDailyReqs {
		t.Errorf("expected daily reqs %f, got %f", expectedDailyReqs, metrics.DailyRequests)
	}
	expectedAvgRPS := expectedDailyReqs / 86400.0
	if metrics.AverageRPS != expectedAvgRPS {
		t.Errorf("expected avg RPS %f, got %f", expectedAvgRPS, metrics.AverageRPS)
	}
	expectedPeakRPS := expectedAvgRPS * 3.0
	if metrics.PeakRPS != expectedPeakRPS {
		t.Errorf("expected peak RPS %f, got %f", expectedPeakRPS, metrics.PeakRPS)
	}
	if metrics.DBReadIOPS <= 0 || metrics.DBWriteIOPS <= 0 {
		t.Errorf("expected positive DB IOPS, got read=%f, write=%f", metrics.DBReadIOPS, metrics.DBWriteIOPS)
	}
}

func TestInfraCostEstimator(t *testing.T) {
	workload := WorkloadMetrics{
		PeakRPS:              2500.0,
		DBWriteIOPS:          400.0,
		AccumulatedStorageTB: 5.0,
		MonthlyBandwidthTB:   12.0,
	}
	cfg := InfraCostConfig{
		Profile:                ProfileGrowingProduction,
		ComputeBaseMonthlyETB:  5000.0,
		ComputeCostPer1KRPS:    250.0,
		DBBaseMonthlyETB:       7000.0,
		DBStorageCostPerGBETB:  18.0,
		ObjectStorageCostPerGB: 3.50,
		CDNEgressCostPerGBETB:  8.50,
		RedundancyBufferPct:    0.20,
	}

	breakdown := CalculateInfraCost(workload, 50000.0, 12, cfg)
	if breakdown.ComputeAndNetwork <= 0 {
		t.Errorf("expected positive compute cost, got %f", breakdown.ComputeAndNetwork)
	}
	if breakdown.Database <= 0 {
		t.Errorf("expected positive database cost, got %f", breakdown.Database)
	}
	if breakdown.StorageAndCDN <= 0 {
		t.Errorf("expected positive storage & CDN cost, got %f", breakdown.StorageAndCDN)
	}
	if breakdown.TotalMonthlyCost <= (breakdown.ComputeAndNetwork + breakdown.Database) {
		t.Errorf("total monthly cost must include all categories + buffer, got %f", breakdown.TotalMonthlyCost)
	}
}

func TestRevenueAndEscrowSeparation(t *testing.T) {
	cfg := RevenueConfig{
		SellersPctOfUsers:              0.10,
		TransactionsPerSellerMo:        2.0,
		AverageTransactionValue:        5000.0, // 5k ETB
		CommissionRatePct:              0.02,   // 2%
		RefundCancellationRate:         0.0,
		CommissionDiscountRate:         0.0,
		ListingsPerSellerMo:            2.0,
		PaidListingFeePerItem:          50.0,
		PaidListingAdoptionPct:         0.20,
		GoldTierConversionPct:          0.01,
		GoldTierMonthlyPrice:           500.0,
		PremiumTierConversionPct:       0.0,
		FeaturedListingMonthlyAdoption: 0.0,
		VerifiedBusinessSellersPct:     0.0,
	}

	mau := 10000.0 // 1,000 sellers, 2,000 transactions
	rev := CalculateRevenue(mau, cfg)

	expectedTransactions := 2000
	if rev.CompletedTransactions != expectedTransactions {
		t.Errorf("expected %d transactions, got %d", expectedTransactions, rev.CompletedTransactions)
	}

	expectedGMV := float64(expectedTransactions) * 5000.0 // 10,000,000 ETB
	if rev.GrossMerchandiseVolume != expectedGMV {
		t.Errorf("expected GMV %f, got %f", expectedGMV, rev.GrossMerchandiseVolume)
	}

	expectedCommission := expectedGMV * 0.02 // 200,000 ETB
	if rev.CommissionRevenue != expectedCommission {
		t.Errorf("expected commission %f, got %f", expectedCommission, rev.CommissionRevenue)
	}

	// Verify HUNTIQ revenue is strictly commission + fees, NOT the 10M GMV!
	if rev.NetRecognizedRevenue >= rev.GrossMerchandiseVolume {
		t.Errorf("CRITICAL VIOLATION: Platform revenue (%f) must never conflate with seller escrow GMV (%f)", rev.NetRecognizedRevenue, rev.GrossMerchandiseVolume)
	}
}

func TestEthiopianTaxCalculations(t *testing.T) {
	cfg := TaxConfig{
		CorporateIncomeTaxRate:   0.30,
		StandardVATRate:          0.15,
		VATRegistrationThreshold: 1000000.0,
		InputVATRecoverablePct:   0.85,
		DisallowedExpenseRate:    0.05,
		LossCarryforwardYears:    5,
	}

	// Case 1: Early-stage loss
	resLoss := CalculateEthiopiaTax(100000.0, 300000.0, 250000.0, 50000.0, 20000.0, 1200000.0, 0.0, cfg)
	if resLoss.EstimatedCorporateIncomeTax != 0 {
		t.Errorf("CIT must be zero when in a tax loss, got %f", resLoss.EstimatedCorporateIncomeTax)
	}
	if resLoss.EndingLossPool <= 0 {
		t.Errorf("Loss carryforward pool should accumulate, got %f", resLoss.EndingLossPool)
	}
	if resLoss.OutputVAT != 15000.0 { // 100k * 15%
		t.Errorf("expected output VAT 15,000, got %f", resLoss.OutputVAT)
	}

	// Case 2: Profitable month utilizing previous loss relief
	resProfit := CalculateEthiopiaTax(500000.0, 200000.0, 150000.0, 50000.0, 30000.0, 6000000.0, 100000.0, cfg)
	if resProfit.LossReliefApplied != 100000.0 {
		t.Errorf("expected 100,000 loss relief applied, got %f", resProfit.LossReliefApplied)
	}
	if resProfit.EndingLossPool != 0.0 {
		t.Errorf("loss pool should be fully consumed, got %f", resProfit.EndingLossPool)
	}
	if resProfit.EstimatedCorporateIncomeTax <= 0 {
		t.Errorf("expected positive CIT after utilizing loss, got %f", resProfit.EstimatedCorporateIncomeTax)
	}
}

func TestFullForecastEngineAndBreakEven(t *testing.T) {
	cfg := GetBaselineScenarioConfig("base")
	cfg.Growth.ForecastHorizonMonths = 60

	report := GenerateForecastReport(cfg)
	if report == nil {
		t.Fatal("expected non-nil forecast report")
	}

	if len(report.MonthlySeries) != 60 {
		t.Fatalf("expected 60 monthly series, got %d", len(report.MonthlySeries))
	}
	if len(report.AnnualSummaries) != 5 {
		t.Fatalf("expected 5 annual summaries, got %d", len(report.AnnualSummaries))
	}

	// Milestone check
	if report.Milestones.BreakEvenMonthIndex <= 0 || report.Milestones.BreakEvenMonthIndex > 60 {
		t.Logf("Break-even month: %d", report.Milestones.BreakEvenMonthIndex)
	}

	// Capacity 20M check
	cap20M := report.Capacity20M
	if cap20M.TargetDAU != 20000000.0 {
		t.Errorf("expected 20M target DAU, got %f", cap20M.TargetDAU)
	}
	if len(cap20M.TechnicalBottlenecks) == 0 {
		t.Errorf("expected documented technical bottlenecks for 20M DAU")
	}

	// Sensitivity check
	if len(report.SensitivityShocks) < 5 {
		t.Errorf("expected at least 5 sensitivity shock results, got %d", len(report.SensitivityShocks))
	}

	// CSV Export test
	csv := GenerateCSVExport(report)
	if !strings.Contains(csv, "GrossMerchandiseVolume") || !strings.Contains(csv, "NetProfitAfterTax") {
		t.Errorf("invalid CSV export structure")
	}
}
