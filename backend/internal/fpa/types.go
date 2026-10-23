package fpa

import "time"

// Currency represents supported ISO currency codes
type Currency string

const (
	CurrencyETB Currency = "ETB" // Default currency: Ethiopian Birr
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

// GrowthModelType defines the mathematical user growth projection algorithm
type GrowthModelType string

const (
	GrowthGeometric        GrowthModelType = "geometric"
	GrowthLinear           GrowthModelType = "linear"
	GrowthUserDefinedRates GrowthModelType = "user_rates"
	GrowthLogistic         GrowthModelType = "logistic_saturation"
	GrowthStableTransition GrowthModelType = "stable_transition"
	GrowthCustomMonthly    GrowthModelType = "custom_monthly"
	GrowthQuarterlySmooth  GrowthModelType = "quarterly_smooth"
	GrowthQuarterlyStep    GrowthModelType = "quarterly_step"
)

// TargetMetric defines what entity the growth model projects
type TargetMetric string

const (
	MetricMAU             TargetMetric = "mau"              // Monthly Active Users
	MetricDAU             TargetMetric = "dau"              // Daily Active Users
	MetricRegisteredUsers TargetMetric = "registered_users" // Cumulative Registered Accounts
)

// GrowthConfig encapsulates user growth modeling parameters
type GrowthConfig struct {
	ModelType             GrowthModelType `json:"model_type"`
	TargetMetric          TargetMetric    `json:"target_metric"`
	InitialUsers          float64         `json:"initial_users"`           // e.g. 1,000
	EarlyMonthlyRate      float64         `json:"early_monthly_rate"`      // e.g. 0.20 (20%)
	EarlyDurationMonths   int             `json:"early_duration_months"`   // e.g. 12 months
	MatureMonthlyRate     float64         `json:"mature_monthly_rate"`     // e.g. 0.02 (2%)
	GrowthPeriodType      string          `json:"growth_period_type"`      // "monthly" (default) or "quarterly"
	EarlyQuarterlyRate    float64         `json:"early_quarterly_rate"`    // e.g. 0.20 (20% per quarter)
	EarlyDurationQuarters int             `json:"early_duration_quarters"` // e.g. 4 quarters (= 12 months)
	MatureQuarterlyRate   float64         `json:"mature_quarterly_rate"`   // e.g. 0.0612 (~2%/mo)
	MarketCeiling         float64         `json:"market_ceiling"`          // e.g. 20,000,000
	ForecastHorizonMonths int             `json:"forecast_horizon_months"` // default 60
	LinearMonthlyAddition float64         `json:"linear_monthly_addition"` // e.g. 5,000 users/mo
	UserMonthlyRates      []float64       `json:"user_monthly_rates"`      // custom rate per month
	CustomMonthlyUsers    []float64       `json:"custom_monthly_users"`    // explicit user array
}

// WorkloadConfig defines user activity and infrastructure demand multipliers
type WorkloadConfig struct {
	DAUtoMAURatio              float64 `json:"dau_to_mau_ratio"`               // e.g. 0.25 (25% of MAU active daily)
	SessionsPerUserPerDay      float64 `json:"sessions_per_user_per_day"`      // e.g. 2.5
	RequestsPerSession         float64 `json:"requests_per_session"`           // e.g. 18.0
	PeakToAverageRatio         float64 `json:"peak_to_average_ratio"`          // e.g. 3.0
	SearchesPerUserDay         float64 `json:"searches_per_user_day"`          // e.g. 4.0
	ListingViewsPerUserDay     float64 `json:"listing_views_per_user_day"`     // e.g. 12.0
	ImageUploadsPerActiveUser  float64 `json:"image_uploads_per_active_user"`  // e.g. 0.05
	ImageViewsPerUserDay       float64 `json:"image_views_per_user_day"`       // e.g. 24.0
	AvgImageSizeKB             float64 `json:"avg_image_size_kb"`              // e.g. 350.0 KB
	DBReadsPerRequest          float64 `json:"db_reads_per_request"`           // e.g. 3.5
	DBWritesPerRequest         float64 `json:"db_writes_per_request"`          // e.g. 0.4
	DataRetainedPerUserKB      float64 `json:"data_retained_per_user_kb"`      // e.g. 50.0 KB
	DataRetainedPerListingKB   float64 `json:"data_retained_per_listing_kb"`   // e.g. 250.0 KB
	LogsRetainedDays           int     `json:"logs_retained_days"`             // e.g. 90
	BackgroundJobsPerUserDay   float64 `json:"background_jobs_per_user_day"`   // e.g. 0.20
	AnnualPriceChangeRate      float64 `json:"annual_price_change_rate"`       // e.g. -0.05 (-5% cloud efficiency)
}

// DeploymentProfile defines the target architectural scale
type DeploymentProfile string

const (
	ProfileLeanEarlyStage     DeploymentProfile = "lean"
	ProfileGrowingProduction  DeploymentProfile = "growing"
	ProfileHighAvailability   DeploymentProfile = "high_availability"
	ProfileLargeMultiRegion   DeploymentProfile = "multi_region"
)

// InfraCostConfig defines unit costs across infrastructure categories
type InfraCostConfig struct {
	Profile                DeploymentProfile `json:"profile"`
	ComputeBaseMonthlyETB  float64           `json:"compute_base_monthly_etb"`
	ComputeCostPer1KRPS    float64           `json:"compute_cost_per_1k_rps"`
	DBBaseMonthlyETB       float64           `json:"db_base_monthly_etb"`
	DBStorageCostPerGBETB  float64           `json:"db_storage_cost_per_gb_etb"`
	ObjectStorageCostPerGB float64           `json:"object_storage_cost_per_gb_etb"`
	CDNEgressCostPerGBETB  float64           `json:"cdn_egress_cost_per_gb_etb"`
	LoggingMonitoringETB   float64           `json:"logging_monitoring_base_etb"`
	SecurityWAFMonthlyETB  float64           `json:"security_waf_monthly_etb"`
	AuthEmailSMSUnitETB    float64           `json:"auth_email_sms_unit_etb"`
	BackupDisasterRecETB   float64           `json:"backup_disaster_rec_etb"`
	RedundancyBufferPct    float64           `json:"redundancy_buffer_pct"` // e.g. 0.25 (25% reserve)
}

// StaffRoleConfig defines compensation and headcount for an organizational department
type StaffRoleConfig struct {
	RoleTitle       string  `json:"role_title"`
	InitialCount    int     `json:"initial_count"`
	TargetCountM12  int     `json:"target_count_m12"`
	TargetCountM60  int     `json:"target_count_m60"`
	MonthlySalaryETB float64 `json:"monthly_salary_etb"`
}

// StaffingExpenseConfig defines organizational staffing and business operating expenses
type StaffingExpenseConfig struct {
	Roles               []StaffRoleConfig `json:"roles"`
	EmployerPensionRate float64           `json:"employer_pension_rate"` // e.g. 0.11 (11% under Proc. 715/2011 & 1268/2022)
	AnnualSalaryRaise   float64           `json:"annual_salary_raise"`   // e.g. 0.08 (8% annual inflation)
	OfficeRentMonthly   float64           `json:"office_rent_monthly"`   // ETB
	LegalComplianceMo   float64           `json:"legal_compliance_mo"`   // ETB
	MarketingBudgetMo   float64           `json:"marketing_budget_mo"`   // Base marketing ETB
	MarketingPerUserCAC float64           `json:"marketing_per_user_cac"`// Target CAC per acquired user
	SoftwareSaaSMo      float64           `json:"software_saas_mo"`      // ETB
	InsuranceMiscMo     float64           `json:"insurance_misc_mo"`     // ETB
	InitialCapExETB     float64           `json:"initial_capex_etb"`     // One-time startup equipment & servers
	CapExDeprecMonths   int               `json:"capex_deprec_months"`   // Accounting useful life (e.g. 36 months)
}

// RevenueConfig defines monetization levers across HUNTIQ revenue streams
type RevenueConfig struct {
	// 1. Marketplace Commission
	SellersPctOfUsers       float64 `json:"sellers_pct_of_users"`        // e.g. 0.15 (15% of active users are sellers)
	TransactionsPerSellerMo float64 `json:"transactions_per_seller_mo"`  // e.g. 3.2
	AverageTransactionValue float64 `json:"average_transaction_value"`   // ETB e.g. 4,500 ETB
	CommissionRatePct       float64 `json:"commission_rate_pct"`         // e.g. 0.02 (2.0% broker fee)
	RefundCancellationRate  float64 `json:"refund_cancellation_rate"`    // e.g. 0.03 (3% refund rate)
	CommissionDiscountRate  float64 `json:"commission_discount_rate"`    // e.g. 0.05 (5% promo discount)
	CommissionVATInclusive  bool    `json:"commission_vat_inclusive"`    // False = 15% VAT added on top

	// 2. Listing Fees
	ListingsPerSellerMo     float64 `json:"listings_per_seller_mo"`      // e.g. 4.0
	PaidListingFeePerItem   float64 `json:"paid_listing_fee_per_item"`   // ETB e.g. 50 ETB
	PaidListingAdoptionPct  float64 `json:"paid_listing_adoption_pct"`   // e.g. 0.20 (20% pay listing fee)

	// 3. Premium Membership Tiers (Gold & Premium)
	GoldTierConversionPct   float64 `json:"gold_tier_conversion_pct"`    // e.g. 0.02 (2% of sellers)
	GoldTierMonthlyPrice    float64 `json:"gold_tier_monthly_price"`     // ETB e.g. 500 ETB
	PremiumTierConversionPct float64 `json:"premium_tier_conversion_pct"`// e.g. 0.008 (0.8% of sellers)
	PremiumTierMonthlyPrice float64 `json:"premium_tier_monthly_price"`  // ETB e.g. 1,500 ETB
	TierMonthlyChurnRate    float64 `json:"tier_monthly_churn_rate"`     // e.g. 0.04 (4% monthly churn)

	// 4. Paid Promotion & Advertising
	FeaturedListingMonthlyAdoption float64 `json:"featured_listing_monthly_adoption"` // e.g. 0.03
	FeaturedListingPrice           float64 `json:"featured_listing_price"`            // ETB e.g. 250 ETB

	// 5. Business Tools & Subscriptions
	VerifiedBusinessSellersPct  float64 `json:"verified_business_sellers_pct"` // e.g. 0.015
	VerifiedBusinessToolFeeMo   float64 `json:"verified_business_tool_fee_mo"` // ETB e.g. 800 ETB
}

// TaxConfig defines the Ethiopian statutory tax parameters
type TaxConfig struct {
	Jurisdiction              string  `json:"jurisdiction"`                 // "ET"
	EntityLegalType           string  `json:"entity_legal_type"`            // "Private Limited Company (PLC)"
	FiscalYearConvention      string  `json:"fiscal_year_convention"`       // "Gregorian" or "Ethiopian Fiscal Year"
	CorporateIncomeTaxRate    float64 `json:"corporate_income_tax_rate"`    // Standard 30% (Proc. 979/2016 Art. 19)
	StandardVATRate           float64 `json:"standard_vat_rate"`            // Standard 15% (Proc. 285/2002 & 1341/2024 Art. 7)
	VATRegistrationThreshold  float64 `json:"vat_registration_threshold"`   // ETB 1,000,000 annual taxable turnover
	InputVATRecoverablePct    float64 `json:"input_vat_recoverable_pct"`    // e.g. 0.85 (85% of OpEx/cloud input VAT eligible)
	WithholdingTaxServicesPct float64 `json:"withholding_tax_services_pct"` // 2% on supplier services (Proc. 979/2016 Art. 92)
	DisallowedExpenseRate     float64 `json:"disallowed_expense_rate"`      // e.g. 0.05 (5% non-deductible adjustments)
	LossCarryforwardYears     int     `json:"loss_carryforward_years"`      // 5 consecutive tax years (Proc. 979/2016 Art. 26)
	RequiresAccountantReview  bool    `json:"requires_accountant_review"`   // True: advisory watermark on reports
}

// CashFlowConfig defines initial balance sheet and liquidity parameters
type CashFlowConfig struct {
	InitialCashBalanceETB float64 `json:"initial_cash_balance_etb"` // e.g. ETB 10,000,000 startup capital
	ReceivablesDays       int     `json:"receivables_days"`         // e.g. 5 days payment gateway settlement
	PayablesDays          int     `json:"payables_days"`            // e.g. 30 days vendor payment terms
	TaxRemittanceLagMonths int    `json:"tax_remittance_lag_months"`// VAT paid M+1; CIT paid annually
}

// FullScenarioConfig groups all parameters for a complete forecast scenario
type FullScenarioConfig struct {
	ScenarioID  string                `json:"scenario_id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Currency    Currency              `json:"currency"`
	Growth      GrowthConfig          `json:"growth"`
	Workload    WorkloadConfig        `json:"workload"`
	Infra       InfraCostConfig       `json:"infra"`
	Staffing    StaffingExpenseConfig `json:"staffing"`
	Revenue     RevenueConfig         `json:"revenue"`
	Tax         TaxConfig             `json:"tax"`
	CashFlow    CashFlowConfig        `json:"cash_flow"`
}

// MonthlyForecastRow contains all calculated financial, technical, and tax metrics for a single month
type MonthlyForecastRow struct {
	MonthIndex int    `json:"month_index"` // 1 to 60
	Period     string `json:"period"`      // e.g. "2026-11"

	// User metrics
	MAU             float64 `json:"mau"`
	DAU             float64 `json:"dau"`
	CumulativeUsers float64 `json:"cumulative_users"`

	// Technical demand & capacity
	DailyRequests      float64 `json:"daily_requests"`
	AverageRPS         float64 `json:"average_rps"`
	PeakRPS            float64 `json:"peak_rps"`
	MonthlyRequests    float64 `json:"monthly_requests"`
	DBReadIOPS         float64 `json:"db_read_iops"`
	DBWriteIOPS        float64 `json:"db_write_iops"`
	AccumulatedStorageTB float64 `json:"accumulated_storage_tb"`
	MonthlyBandwidthTB   float64 `json:"monthly_bandwidth_tb"`

	// Revenue & Marketplace escrow separation
	GrossMerchandiseVolume float64 `json:"gross_merchandise_volume"` // Funds belonging to sellers (GMV)
	CompletedTransactions  int     `json:"completed_transactions"`
	CommissionRevenue      float64 `json:"commission_revenue"`
	ListingFeeRevenue      float64 `json:"listing_fee_revenue"`
	SubscriptionRevenue    float64 `json:"subscription_revenue"`
	PromotedListingRevenue float64 `json:"promoted_listing_revenue"`
	BusinessToolsRevenue   float64 `json:"business_tools_revenue"`
	GrossRecognizedRevenue float64 `json:"gross_recognized_revenue"`
	DiscountsAndRefunds    float64 `json:"discounts_and_refunds"`
	NetRecognizedRevenue   float64 `json:"net_recognized_revenue"`

	// Operating Costs
	ComputeAndNetworkCost float64 `json:"compute_and_network_cost"`
	DatabaseCost          float64 `json:"database_cost"`
	StorageAndCDNCost     float64 `json:"storage_and_cdn_cost"`
	SecurityAndAuthCost   float64 `json:"security_and_auth_cost"`
	TotalInfrastructureCost float64 `json:"total_infrastructure_cost"`

	StaffSalariesCost     float64 `json:"staff_salaries_cost"`
	EmployerPensionCost   float64 `json:"employer_pension_cost"`
	MarketingExpense      float64 `json:"marketing_expense"`
	OfficeAdminLegalCost  float64 `json:"office_admin_legal_cost"`
	SoftwareSaaSCost      float64 `json:"software_saas_cost"`
	TotalCashOpEx         float64 `json:"total_cash_opex"`

	AccountingDepreciation float64 `json:"accounting_depreciation"`
	TotalOperatingExpense  float64 `json:"total_operating_expense"` // Cash OpEx + Depreciation + Infra

	// Profitability
	GrossProfit     float64 `json:"gross_profit"`     // Net Revenue - Direct Infra Cost
	GrossMarginPct  float64 `json:"gross_margin_pct"`
	EBITDA          float64 `json:"ebitda"`           // Net Revenue - Cash OpEx - Direct Infra
	EBITDAMarginPct float64 `json:"ebitda_margin_pct"`
	OperatingProfit float64 `json:"operating_profit"` // EBIT
	OperatingMarginPct float64 `json:"operating_margin_pct"`

	// Tax Calculation (Ethiopia Statutory Engine)
	AccountingProfitBeforeTax float64 `json:"accounting_profit_before_tax"`
	TaxableIncomeAdjustments  float64 `json:"taxable_income_adjustments"`
	TaxableIncome             float64 `json:"taxable_income"`
	LossReliefApplied         float64 `json:"loss_relief_applied"`
	CumulativeLossPool        float64 `json:"cumulative_loss_pool"`
	EstimatedIncomeTax        float64 `json:"estimated_income_tax"` // CIT at 30%
	OutputVAT                 float64 `json:"output_vat"`           // 15% on taxable platform revenue
	EligibleInputVAT          float64 `json:"eligible_input_vat"`    // Creditable input VAT on expenses
	NetVATPayable             float64 `json:"net_vat_payable"`       // Output - Input VAT
	TotalStatutoryTaxAccrual  float64 `json:"total_statutory_tax_accrual"`
	NetProfitAfterTax         float64 `json:"net_profit_after_tax"`

	// Cash Flow & Balance Sheet
	OperatingCashInflow   float64 `json:"operating_cash_inflow"`
	OperatingCashOutflow  float64 `json:"operating_cash_outflow"`
	TaxCashOutflow        float64 `json:"tax_cash_outflow"` // VAT paid monthly, CIT paid annually
	NetOperatingCashFlow  float64 `json:"net_operating_cash_flow"`
	CapExCashOutflow      float64 `json:"capex_cash_outflow"`
	NetMonthlyCashChange  float64 `json:"net_monthly_cash_change"`
	ClosingCashBalance    float64 `json:"closing_cash_balance"`
	MonthlyBurnRate       float64 `json:"monthly_burn_rate"`
	CashRunwayMonths      float64 `json:"cash_runway_months"` // -1 = Infinite / Cash positive

	// Unit Economics
	RevenuePerActiveUser     float64 `json:"revenue_per_active_user"` // ARPU
	CostPerActiveUser        float64 `json:"cost_per_active_user"`
	InfraCostPer1KRequests   float64 `json:"infra_cost_per_1k_requests"`
	InfraCostPerActiveUser   float64 `json:"infra_cost_per_active_user"`
	InfraCostPerTransaction  float64 `json:"infra_cost_per_transaction"`
	EstimatedCAC             float64 `json:"estimated_cac"`
	ContributionMarginRatio  float64 `json:"contribution_margin_ratio"`
	MonthlyBreakEvenRevenue  float64 `json:"monthly_break_even_revenue"`
	IsBreakEven              bool    `json:"is_break_even"`
}

// AnnualSummary aggregates monthly forecast metrics into annual planning blocks
type AnnualSummary struct {
	YearNumber             int     `json:"year_number"`
	CalendarYear           int     `json:"calendar_year"`
	EndingMAU              float64 `json:"ending_mau"`
	AverageDAU             float64 `json:"average_dau"`
	TotalGMV               float64 `json:"total_gmv"`
	GrossRevenue           float64 `json:"gross_revenue"`
	NetRevenue             float64 `json:"net_revenue"`
	TotalInfraCost         float64 `json:"total_infra_cost"`
	TotalStaffingCost      float64 `json:"total_staffing_cost"`
	TotalOpEx              float64 `json:"total_opex"`
	EBITDA                 float64 `json:"ebitda"`
	OperatingProfit        float64 `json:"operating_profit"`
	EstimatedCIT           float64 `json:"estimated_cit"`
	NetVATPayable          float64 `json:"net_vat_payable"`
	NetProfitAfterTax      float64 `json:"net_profit_after_tax"`
	EndingCashBalance      float64 `json:"ending_cash_balance"`
	TotalRequestsAnnual    float64 `json:"total_requests_annual"`
	StorageEndingTB        float64 `json:"storage_ending_tb"`
	AverageARPU            float64 `json:"average_arpu"`
	AverageCostPerUser     float64 `json:"average_cost_per_user"`
}

// MilestoneAnalysis tracks key corporate planning achievements
type MilestoneAnalysis struct {
	BreakEvenMonthIndex int     `json:"break_even_month_index"` // -1 if not achieved within horizon
	BreakEvenPeriod     string  `json:"break_even_period"`
	BreakEvenRevenue    float64 `json:"break_even_revenue"`
	PeakFundingNeedETB  float64 `json:"peak_funding_need_etb"`
	LowestCashMonth     string  `json:"lowest_cash_month"`
	LowestCashBalance   float64 `json:"lowest_cash_balance"`
	Month12MAU          float64 `json:"month_12_mau"`
	Month12NetRevenue   float64 `json:"month_12_net_revenue"`
	Month60MAU          float64 `json:"month_60_mau"`
	Month60NetRevenue   float64 `json:"month_60_net_revenue"`
}

// CapacityPlanning20M provides dedicated technical scale analysis for the 20,000,000 DAU milestone
type CapacityPlanning20M struct {
	TargetDAU               float64  `json:"target_dau"` // 20,000,000
	EstimatedMAU            float64  `json:"estimated_mau"`
	PeakRPS                 float64  `json:"peak_rps"`
	AverageRPS              float64  `json:"average_rps"`
	DailyAPIRequests        float64  `json:"daily_api_requests"`
	MonthlyAPIRequests      float64  `json:"monthly_api_requests"`
	RequiredAppCores        int      `json:"required_app_cores"`
	RequiredRAMGB           int      `json:"required_ram_gb"`
	DatabaseReadIOPS        float64  `json:"database_read_iops"`
	DatabaseWriteIOPS       float64  `json:"database_write_iops"`
	EstimatedListingImagesTB float64  `json:"estimated_listing_images_tb"`
	MonthlyEgressBandwidthTB float64 `json:"monthly_egress_bandwidth_tb"`
	EstimatedMonthlyCloudETB float64 `json:"estimated_monthly_cloud_etb"`
	ArchitectureFeasibility string   `json:"architecture_feasibility"`
	TechnicalBottlenecks    []string `json:"technical_bottlenecks"`
	RequiredScaleMeasures   []string `json:"required_scale_measures"`
}

// SensitivityShockResult evaluates the financial consequence of a specific parameter shock
type SensitivityShockResult struct {
	ShockName             string  `json:"shock_name"`
	Description           string  `json:"description"`
	OriginalBreakEvenMo   int     `json:"original_break_even_mo"`
	ShockBreakEvenMo      int     `json:"shock_break_even_mo"`
	OriginalYear1Profit   float64 `json:"original_year_1_profit"`
	ShockYear1Profit      float64 `json:"shock_year_1_profit"`
	OriginalLowestCash    float64 `json:"original_lowest_cash"`
	ShockLowestCash       float64 `json:"shock_lowest_cash"`
	SurvivalCashNeededETB float64 `json:"survival_cash_needed_etb"`
	KeyTakeaway           string  `json:"key_takeaway"`
}

// ForecastReport is the comprehensive management planning deliverable
type ForecastReport struct {
	ScenarioID            string                   `json:"scenario_id"`
	ScenarioName          string                   `json:"scenario_name"`
	Currency              Currency                 `json:"currency"`
	GeneratedAt           time.Time                `json:"generated_at"`
	HorizonMonths         int                      `json:"horizon_months"`
	DisclaimerNotice      string                   `json:"disclaimer_notice"`
	Milestones            MilestoneAnalysis        `json:"milestones"`
	Capacity20M           CapacityPlanning20M      `json:"capacity_20m"`
	AnnualSummaries       []AnnualSummary          `json:"annual_summaries"`
	MonthlySeries         []MonthlyForecastRow     `json:"monthly_series"`
	SensitivityShocks     []SensitivityShockResult `json:"sensitivity_shocks"`
}
