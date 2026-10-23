package fpa

import (
	"math"
)

// StaffingAndOpExBreakdown holds compensation, non-headcount OpEx, and depreciation
type StaffingAndOpExBreakdown struct {
	TotalHeadcount         int
	GrossSalaries          float64
	EmployerPension        float64
	MarketingExpense       float64
	OfficeRentAndAdmin     float64
	SoftwareSaaS           float64
	LegalAndCompliance     float64
	TotalCashOpEx          float64
	AccountingDepreciation float64
	TotalOpExWithDeprec    float64
}

// CalculateStaffingAndOpEx computes organizational expenses for month m (1-indexed)
func CalculateStaffingAndOpEx(monthIndex int, newUsersThisMonth float64, cfg StaffingExpenseConfig) StaffingAndOpExBreakdown {
	yearIndex := (monthIndex - 1) / 12
	annualRaise := cfg.AnnualSalaryRaise
	if annualRaise <= 0 {
		annualRaise = 0.08 // 8% annual Ethiopian cost of living adjustment
	}
	salaryEscalationFactor := math.Pow(1.0+annualRaise, float64(yearIndex))

	totalHeadcount := 0
	grossSalaries := 0.0

	for _, role := range cfg.Roles {
		count := interpolateHeadcount(monthIndex, role.InitialCount, role.TargetCountM12, role.TargetCountM60)
		totalHeadcount += count
		monthlyBase := role.MonthlySalaryETB
		if monthlyBase <= 0 {
			monthlyBase = 35000.0 // Default baseline ETB
		}
		escalatedSalary := monthlyBase * salaryEscalationFactor
		grossSalaries += float64(count) * escalatedSalary
	}

	pensionRate := cfg.EmployerPensionRate
	if pensionRate <= 0 {
		pensionRate = 0.11 // 11% statutory employer pension
	}
	employerPension := grossSalaries * pensionRate

	// Marketing budget: fixed base + CAC per newly acquired user
	baseMarketing := cfg.MarketingBudgetMo
	if baseMarketing <= 0 {
		baseMarketing = 25000.0 // Base ETB
	}
	cac := cfg.MarketingPerUserCAC
	if cac <= 0 {
		cac = 15.0 // 15 ETB customer acquisition cost
	}
	if newUsersThisMonth < 0 {
		newUsersThisMonth = 0
	}
	marketingExpense := baseMarketing + (newUsersThisMonth * cac)

	officeRent := cfg.OfficeRentMonthly
	if officeRent <= 0 {
		officeRent = 45000.0
	}
	// Office rent escalation with team growth
	officeFactor := 1.0 + (float64(totalHeadcount) / 40.0)
	officeRent = officeRent * officeFactor

	softwareSaaS := cfg.SoftwareSaaSMo
	if softwareSaaS <= 0 {
		softwareSaaS = 15000.0
	}
	softwareSaaS += float64(totalHeadcount) * 1200.0 // per-seat tooling

	legalComp := cfg.LegalComplianceMo
	if legalComp <= 0 {
		legalComp = 12000.0
	}

	miscInsurance := cfg.InsuranceMiscMo
	if miscInsurance <= 0 {
		miscInsurance = 8000.0
	}

	totalCashOpEx := grossSalaries + employerPension + marketingExpense + officeRent + softwareSaaS + legalComp + miscInsurance

	// CapEx Depreciation (Straight-Line over useful life)
	deprecMonths := cfg.CapExDeprecMonths
	if deprecMonths <= 0 {
		deprecMonths = 36 // 3-year useful life for IT & server gear
	}
	monthlyDepreciation := 0.0
	if monthIndex <= deprecMonths && cfg.InitialCapExETB > 0 {
		monthlyDepreciation = cfg.InitialCapExETB / float64(deprecMonths)
	}

	return StaffingAndOpExBreakdown{
		TotalHeadcount:         totalHeadcount,
		GrossSalaries:          round2(grossSalaries),
		EmployerPension:        round2(employerPension),
		MarketingExpense:       round2(marketingExpense),
		OfficeRentAndAdmin:     round2(officeRent + miscInsurance),
		SoftwareSaaS:           round2(softwareSaaS),
		LegalAndCompliance:     round2(legalComp),
		TotalCashOpEx:          round2(totalCashOpEx),
		AccountingDepreciation: round2(monthlyDepreciation),
		TotalOpExWithDeprec:    round2(totalCashOpEx + monthlyDepreciation),
	}
}

func interpolateHeadcount(month, initCount, m12Count, m60Count int) int {
	if month <= 1 {
		return initCount
	}
	if month <= 12 {
		progress := float64(month-1) / 11.0
		return int(math.Round(float64(initCount) + progress*float64(m12Count-initCount)))
	}
	if month <= 60 {
		progress := float64(month-12) / 48.0
		return int(math.Round(float64(m12Count) + progress*float64(m60Count-m12Count)))
	}
	return m60Count
}
