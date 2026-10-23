# HUNTIQ FP&A — Financial Model, Tax & Forecast Reconciliation Audit

**Audit Date**: October 23, 2026  
**Auditor**: Antigravity Quality & FP&A Engineering Team  
**Repository**: `/home/poven/GO PROJECT/sellbuy`  
**Target Module**: Financial Planning & Analysis (FP&A), Infrastructure Capacity & Ethiopian Tax Engine  
**Route**: `/admin/finance`  
**Audit Mode**: Read-Only Verification Pass (Mandatory Rule: No code, DB, or migration changes)

---

## Executive Summary

An independent, line-by-line reconciliation of the HUNTIQ Financial Planning & Analysis module was conducted against the active Go backend implementation ([`backend/internal/fpa/`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/)), database migration 030 ([`backend/internal/database/migrations/030_financial_planning_and_forecasting.sql`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/database/migrations/030_financial_planning_and_forecasting.sql)), PostgreSQL schema (`hulubedej_db`), and frontend data bindings ([`frontend/src/app/admin/finance/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/finance/page.tsx)).

### Key Finding
**The Go calculation engine, database schema, and frontend UI are mathematically consistent with each other, but the previous implementation report contained extensive, severe discrepancies and fabricated forecast claims.** 

The previous report claimed that under the base scenario, HUNTIQ achieves break-even in Month 19 (at ~85,000 active users and 1.84M ETB monthly revenue) and reaches ~2,420,000 active users and +59,060,000 ETB monthly operating profit by Year 5. 

**In reality, running the actual Go engine reveals:**
1. Under the configured base scenario assumptions, HUNTIQ reaches **7,430 MAU** at Month 12, **8,535 MAU** at Month 19, and **19,222 MAU** at Month 60.
2. **Break-even is NEVER achieved in 60 months** (`BreakEvenMonthIndex = -1`).
3. Starting from 10,000,000 ETB in initial cash, the company burns through its entire cash reserves by **Month 14** (due to aggressive headcount expansion from 8 to 21 employees by Month 12 against slow 2% monthly mature user growth).
4. By Month 19, the company is operating at a **monthly net loss of -945,986 ETB** and has accumulated a **cash deficit of -5,303,802 ETB**.
5. By Month 60, the monthly operating deficit is **-2,467,071 ETB**, with a cumulative cash deficit of **-82,126,851 ETB**.
6. The 20M DAU workload figures in the previous report (69,444 Average RPS, 208,333 Peak RPS, 6.8 PB storage) were not produced by [`workload.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/workload.go); the actual Go implementation calculates **10,417 Average RPS**, **31,250 Peak RPS**, and **~470 TB total storage**.

The financial dashboard at `/admin/finance` and the CSV export endpoint correctly render the true engine outputs, not the fabricated report figures. However, before the dashboard is used for real corporate budgeting, hiring decisions, or investor presentations, the baseline assumptions must be adjusted to reflect economic viability, and administrative separation of duties must be tightened.

---

## 1. Reproduced Forecasts and Their Inputs

### Baseline Input Parameters (Seeded in [`engine.go:582-686`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/engine.go#L582-L686))
* **Growth Model**: `GrowthStableTransition`, Target: `MetricMAU`
  * Initial users: 1,000
  * Early monthly growth rate: 20.0%
  * Early duration: 12 months
  * Mature monthly growth rate: 2.0%
  * Market ceiling: 20,000,000 active users
  * Horizon: 60 months
* **Workload Multipliers**:
  * DAU-to-MAU ratio: 0.25 (25% daily activity)
  * Sessions per user per day: 2.5
  * Requests per session: 18.0 (Blended daily workload = 45.0 API requests/DAU/day)
  * Peak-to-average ratio: 3.0
  * Listing uploads: 0.05 uploads/user/mo with 4 photos @ 350 KB compressed WebP
  * DB reads/request: 3.5, DB writes/request: 0.4
* **Infrastructure Pricing Profile**: `growing` (1.0x baseline multiplier)
  * Compute: 5,000 ETB base + 250 ETB per 1k Peak RPS
  * PostgreSQL: 7,000 ETB base + 18 ETB/GB SSD
  * S3 Object Storage: 3.50 ETB/GB/mo
  * CDN Egress: 8.50 ETB/GB
  * Security & Auth: 3,000 ETB WAF base + 1.20 ETB per OTP SMS (40% active user trigger rate)
  * Buffer: 20% redundancy headroom
* **Headcount & OpEx Ramp**:
  * Headcount: 8 initial $\rightarrow$ 21 at Month 12 $\rightarrow$ 52 at Month 60
  * Roles: Go Backend, Next.js Frontend, DevOps/SRE, Security/QA, Product/Ops, Support, Marketing, Legal/Finance
  * Salaries: 26,000 – 55,000 ETB base with 8% annual escalation
  * Employer Pension: 11% statutory contribution (Proc. 715/2011 & 1268/2022)
  * Marketing: 30,000 ETB base + 15 ETB CAC per net new user
  * Rent & Admin: 45,000 ETB base escalating with team size + 16,000 ETB SaaS base + 1,200 ETB/seat tooling
  * CapEx: 500,000 ETB launch CapEx depreciated straight-line over 36 months (13,888.89 ETB/mo)
* **Marketplace Monetization**:
  * Active sellers: 15% of MAU; 3.2 completed transactions/seller/month
  * Average Transaction Value (ATV): 4,500 ETB
  * Commission: 2.0% gross (3.0% refunds, 5.0% promo discounts $\rightarrow$ 1.84% net)
  * Listing fees: 4 listings/seller, 20% paid tier @ 50 ETB
  * Subscriptions: 2% Gold @ 500 ETB/mo; 0.8% Premium @ 1,500 ETB/mo
  * Promotions: 3% featured @ 250 ETB; 1.5% verified business tools @ 800 ETB/mo
* **Ethiopian Tax Engine**:
  * Corporate Income Tax (CIT): 30% of net taxable income
  * VAT: 15% standard rate on platform revenue (1M ETB registration threshold)
  * Input VAT: 85% recoverable factor on hosting, rent, and software
  * Disallowed expenses: 5% addback on cash OpEx
  * Loss relief: Cumulative loss pool carryforward
* **Cash Flow Initial Balance**:
  * Initial capital: **10,000,000 ETB**

---

### Reproduced Month-by-Month Milestone Series

| Milestone | Active Users (MAU) | Daily Active (DAU) | Gross Merchandise Volume (GMV) | Net Recognized Revenue | Total Infrastructure | Total Cash OpEx | Total OpEx (inc. Deprec) | Operating Profit (EBIT) | Closing Cash Balance | Monthly Burn Rate | Runway (Months) |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **Month 1** | 1,000 | 250 | 2,160,000 ETB | 55,344 ETB | 24,676 ETB | 534,210 ETB | 572,775 ETB | -517,431 ETB | 8,996,458 ETB | 1,003,542 ETB | 8.96 |
| **Month 6** | 2,488 | 622 | 5,373,000 ETB | 137,707 ETB | 26,603 ETB | 785,739 ETB | 826,231 ETB | -688,524 ETB | 5,236,252 ETB | 694,778 ETB | 7.54 |
| **Month 12** | 7,430 | 1,858 | 16,047,000 ETB | 411,174 ETB | 31,474 ETB | 1,157,420 ETB | 1,202,783 ETB | -791,609 ETB | 1,160,521 ETB | 820,361 ETB | 1.41 |
| **Month 13** | 7,579 | 1,895 | 16,371,000 ETB | 419,454 ETB | 31,648 ETB | 1,219,396 ETB | 1,264,933 ETB | -845,479 ETB | 285,071 ETB | 875,450 ETB | 0.33 |
| **Month 14** | 7,730 | 1,933 | 16,699,500 ETB | 427,870 ETB | 31,827 ETB | 1,244,482 ETB | 1,290,198 ETB | -862,328 ETB | **-608,624 ETB** | 893,694 ETB | **0.00 (Depleted)** |
| **Month 19** | 8,535 | 2,134 | 18,436,500 ETB | 472,375 ETB | 32,765 ETB | 1,371,707 ETB | 1,418,361 ETB | **-945,986 ETB** | **-5,303,802 ETB** | 982,863 ETB | **0.00** |
| **Month 36** | 12,019 | 3,005 | 25,960,500 ETB | 665,116 ETB | 36,493 ETB | 2,096,440 ETB | 2,146,822 ETB | -1,481,706 ETB | -27,088,298 ETB | 1,535,420 ETB | **0.00** |
| **Month 60** | 19,222 | 4,806 | 41,521,500 ETB | 1,063,861 ETB | 45,062 ETB | 3,485,871 ETB | 3,530,932 ETB | **-2,467,071 ETB** | **-82,126,851 ETB** | 2,596,693 ETB | **0.00** |

---

### Reproduced 5-Year Annual Summaries

| Metric | Year 1 | Year 2 | Year 3 | Year 4 | Year 5 |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Ending MAU** | 7,430 | 9,423 | 11,951 | 15,157 | 19,222 |
| **Average DAU** | 754 | 2,109 | 2,674 | 3,392 | 4,302 |
| **Gross Merchandise Volume (GMV)** | 85,486,500 ETB | 219,555,000 ETB | 278,460,000 ETB | 353,137,500 ETB | 447,867,000 ETB |
| **Net Recognized Revenue** | 2,190,407 ETB | 5,625,490 ETB | 7,134,687 ETB | 9,048,193 ETB | 11,475,357 ETB |
| **Total Infrastructure Cost** | 325,523 ETB | 392,269 ETB | 423,883 ETB | 463,704 ETB | 514,019 ETB |
| **Total Staffing & OpEx** | 10,059,781 ETB | 16,700,776 ETB | 23,943,482 ETB | 30,362,869 ETB | 39,952,934 ETB |
| **Depreciation (CapEx)** | 166,667 ETB | 166,667 ETB | 166,667 ETB | 0 ETB | 0 ETB |
| **Total Operating Expenses** | 10,551,971 ETB | 17,259,711 ETB | 24,534,032 ETB | 30,826,574 ETB | 40,466,952 ETB |
| **EBITDA** | -8,194,897 ETB | -11,467,555 ETB | -17,232,679 ETB | -21,778,381 ETB | -28,991,596 ETB |
| **Operating Profit (EBIT)** | **-8,361,564 ETB** | **-11,634,221 ETB** | **-17,399,345 ETB** | **-21,778,381 ETB** | **-28,991,596 ETB** |
| **Estimated Corporate Income Tax** | 0 ETB | 0 ETB | 0 ETB | 0 ETB | 0 ETB |
| **Net VAT Payable** | 144,582 ETB | 601,087 ETB | 792,013 ETB | 1,049,654 ETB | 1,374,409 ETB |
| **Net Profit / (Loss) After Tax** | **-8,361,564 ETB** | **-11,634,221 ETB** | **-17,399,345 ETB** | **-21,778,381 ETB** | **-28,991,596 ETB** |
| **Ending Cash Balance** | **1,160,521 ETB** | **-10,908,120 ETB** | **-28,932,812 ETB** | **-51,760,846 ETB** | **-82,126,851 ETB** |

---

## 2. Discrepancy Reconciliation Table

| Metric / Scenario Point | Claimed in Previous Report | Recalculated Actual Engine Output | Variance | Root Cause Analysis |
| :--- | :--- | :--- | :--- | :--- |
| **Month 12 MAU** | ~8,916 users | **7,430.08 users** | -1,486 (-16.7%) | The previous report calculated $1000 \times (1.20)^{12} = 8,916.1$ assuming 12 compoundings. The Go engine starts Month 1 at 1,000, compounding 11 times by Month 12 ($1000 \times 1.20^{11} = 7,430.08$). |
| **Month 19 MAU** | ~85,000 users | **8,534.83 users** | -76,465 (-89.9%) | **Decimal/Zero Hallucination**: Month 19 has 7 compoundings at 2% over Month 12 ($7430.08 \times 1.02^7 = 8,534.83$). The previous report added an extra zero (8,535 became ~85,000). |
| **Month 19 Break-Even** | "Break-Even achieved in Month 19 at 1.84M ETB/mo" | **NOT ACHIEVED (`BreakEvenMonthIndex = -1`)** | Complete contradiction | Month 19 revenue is 472,375 ETB vs required break-even revenue of 1,436,339 ETB. Monthly operating loss is -945,986 ETB. |
| **Initial Capital** | ETB 5,000,000 | **ETB 10,000,000** | +5,000,000 (+100%) | `engine.go:683` sets `InitialCashBalanceETB = 10,000,000`. The report misquoted the initial balance sheet. |
| **Month 12 Cash Runway** | 14.2 months | **1.41 months** | -12.79 months (-90.1%) | Claimed runway was based on fictional low burn (173k ETB). Actual M12 cash is 1.16M ETB against an 820k ETB monthly burn, exhausting cash in Month 14. |
| **Month 60 MAU** | ~2,420,000 users | **19,222.15 users** | -2,400,778 (-99.2%) | **Fabricated Metric**: 48 months of 2% mature growth yields $7430 \times 1.02^{48} = 19,222$. Reaching 2.42M would require 14.1% monthly compounding for 5 straight years. |
| **Month 60 Net Revenue** | ~74,800,000 ETB/mo | **1,063,861 ETB/mo** | -73,736,139 (-98.6%) | Consequence of fabricated 2.42M user base. Net revenue at 19k MAU is 1.06M ETB. |
| **Month 60 Infra Cost** | ~6,840,000 ETB/mo | **45,062 ETB/mo** | -6,794,938 (-99.3%) | Consequence of fabricated 2.42M user base. Real compute/storage cost for 19k MAU is 45k ETB/mo. |
| **Month 60 Operating Profit (EBIT)** | **+59,060,000 ETB/mo** | **-2,467,071 ETB/mo** | -61,527,071 (Reversal) | The actual base model is heavily loss-making; the previous report falsely claimed massive operational profitability. |
| **Month 60 CIT Accrual** | ~17,718,000 ETB/mo | **0.00 ETB** | -17,718,000 | The company never generates taxable profit; cumulative loss pool reaches 82M ETB, resulting in zero CIT. |
| **20M DAU Average RPS** | ~69,444 RPS | **10,416.67 RPS** | -59,027 (-85.0%) | `workload.go` uses 45 req/DAU/day ($20M \times 45 / 86400 = 10,417$ RPS). 69,444 RPS requires 300 req/DAU/day. |
| **20M DAU Peak RPS** | ~208,333 RPS | **31,250.00 RPS** | -177,083 (-85.0%) | Peak RPS in `workload.go` is $10,416.67 \times 3.0 = 31,250$. The report published a 6.67x higher number. |
| **20M DAU Storage** | "Over 6.8 PB" | **469.75 TB** | -6.33 PB (-93.1%) | Image uploads (0.05/user $\times$ 4 photos $\times$ 350 KB $\times$ 60 mos $\times$ 1.5 repl) produces ~470 TB, not 6.8 PB. |
| **WHT on Services Threshold** | 10,000 ETB | **3,000 ETB** | -7,000 ETB | Proc. 979/2016 Art. 92(1)(a) mandates 3,000 ETB for services. 10,000 ETB applies to goods (Art. 92(1)(b)). |

---

## 3. Corrected Formulas and Missing Assumptions

### A. Growth Model Compounding & Phase Transition ([`growth.go:120-135`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/growth.go#L120-L135))
In the current implementation:
```go
for m := 0; m < horizon; m++ {
    if m > 0 {
        rate := earlyRate
        if m >= earlyDuration {
            rate = matureRate
        }
        ...
        current = current * (1.0 + rate)
    }
    results[m] = current
}
```
* At `m = 0` (Month 1): `results[0] = 1,000.0`.
* By Month 12 (`m = 11`), only 11 compoundings occur ($1,000 \times 1.20^{11} = 7,430.08$).
* At Month 13 (`m = 12`), `m >= earlyDuration` (`12 >= 12`), switching immediately to 2% mature growth ($7,430.08 \times 1.02 = 7,578.69$).
* **Correction**: To provide a full 12 months of early compounding, the transition condition should be `if m > earlyDuration` rather than `m >= earlyDuration`, or Month 1 must be designated as baseline launch with 12 active adoption months following.

### B. Structural Economic Unviability of the Base Scenario
The reason the base model collapses into an 82,000,000 ETB deficit is an **unaligned cost-versus-growth trajectory**:
1. **Headcount ramp is decoupled from revenue**: Headcount expands linearly from 8 to 21 employees by Month 12, and to 52 employees by Month 60, irrespective of platform revenue.
2. **Growth rate collapses too early**: Dropping from 20%/mo to 2%/mo at Month 12 caps adoption at just 19k users by Year 5. A 52-person tech firm cannot be sustained by 19,000 marketplace users at an ARPU of ~55 ETB/month.
3. **Required Adjustment**: To achieve break-even, the model requires either:
   * **Milestone-Gated Headcount**: Hiring new staff only when MAU or gross revenue crosses predetermined milestones (e.g. 5,000 MAU per support agent).
   * **Extended Growth Phase**: Sustaining 10–12% monthly growth through Month 24–36.
   * **Additional Equity Capital**: Modeling Series A / Series B equity injections to bridge the ~15,000,000 ETB funding gap in Months 14–36.

### C. `CommissionVATInclusive` Dead Code ([`types.go:129`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/types.go#L129))
* The configuration struct defines `CommissionVATInclusive bool`, but it is **never referenced** in `CalculateRevenue` or `CalculateEthiopiaTax`.
* Currently, 15% VAT is applied unconditionally on top of all recognized revenue.
* **Correction**: When `CommissionVATInclusive` is true, net platform commission must be calculated as:
  $$\text{Net Commission} = \frac{\text{Gross Commission}}{1 + \text{VAT Rate}}$$
  $$\text{Output VAT} = \text{Gross Commission} - \text{Net Commission}$$

### D. Loss Carryforward Aging ([`tax.go:41-65`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/tax.go#L41-L65))
* Article 26 of Proclamation No. 979/2016 limits business loss relief to **5 consecutive tax years**.
* The current engine uses a single scalar `lossPool` float. An unutilized loss from Year 1 carries forward indefinitely past Year 5.
* **Correction**: Implement a slice tracking annual loss tranches (`[5]float64`), expiring losses that exceed 5 fiscal cycles.

### E. Omission of Withholding Tax in Cash Flows
* Proclamation No. 979/2016 Article 92 requires 2% withholding on domestic service disbursements (> 3,000 ETB).
* While seeded in `financial_tax_rules`, withholding tax is not deducted from supplier OpEx disbursements in `CalculateStaffingAndOpEx` or `CalculateCashFlowAndBreakEven`.

---

## 4. Actual-versus-Projected Data Provenance

An audit of the database confirms strict adherence to data integrity rules:
* `SELECT count(*) FROM financial_actuals;` returns **`0`**.
* **Zero Fabricated Actuals**: No fake historical revenue, transactions, expenses, or tax payments have been entered into the system.
* **Projections Segregation**: Projections are computed on-the-fly by the Go engine. The frontend explicitly marks all forecast metrics with the advisory watermark:
  > *"ILLUSTRATIVE FINANCIAL PLANNING MODEL ONLY — NOT CERTIFIED ACCOUNTING RECORDS."*
* **Scenario Overrides**: Custom assumptions modified via `PUT /api/admin/finance/assumptions` are stored in `financial_assumptions` by `scenario_id` and `section` without touching base defaults or historical tables.

---

## 5. Ethiopian Tax-Rule Verification Status

Audit of the 7 tax rules seeded in `financial_tax_rules` via migration `030`:

| Rule Code | Description | Standard Rate / Threshold | Authoritative Source Citation | Current Status | Application in Go Engine |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `CIT_STANDARD` | Corporate Income Tax (Body Corporate) | 30.00% | Federal Income Tax Proclamation No. 979/2016, Art. 19(1) | **Statutory baseline, unverified** (`is_verified = false`) | Applied to `netTaxableIncome` in `tax.go:74` |
| `VAT_STANDARD` | Value Added Tax (Goods & Digital Supplies) | 15.00% | VAT Proclamation No. 285/2002 & Proclamation No. 1341/2024, Art. 7 | **Statutory baseline, unverified** (`is_verified = false`) | Applied to `netRevenue` in `tax.go:94` |
| `VAT_THRESHOLD` | Compulsory VAT Registration Threshold | 1,000,000.00 ETB / year | VAT Proclamation No. 285/2002, Art. 16 | **Statutory baseline, unverified** (`is_verified = false`) | Gating condition in `tax.go:82` |
| `WHT_SERVICES` | Withholding Tax on Supply of Services | 2.00% (> 3,000 ETB) | Federal Income Tax Proclamation No. 979/2016, Art. 92(1)(a) | **Statutory baseline, unverified** (`is_verified = false`) | **Not implemented in OpEx cash flows** |
| `WHT_GOODS` | Withholding Tax on Supply of Goods | 2.00% (> 10,000 ETB) | Federal Income Tax Proclamation No. 979/2016, Art. 92(1)(b) | **Statutory baseline, unverified** (`is_verified = false`) | **Not implemented in OpEx cash flows** |
| `PENSION_EMPLOYER` | Private Org Employees Pension (Employer Share) | 11.00% of basic salary | Pension Proclamations No. 715/2011 & No. 1268/2022, Art. 10 | **Statutory baseline, unverified** (`is_verified = false`) | Applied to gross salaries in `staffing_opex.go:48` |
| `LOSS_CARRYFORWARD`| Business Loss Relief Limit | 5 tax years | Federal Income Tax Proclamation No. 979/2016, Art. 26 | **Statutory baseline, unverified** (`is_verified = false`) | Applied as scalar pool; lacks 5-year vintage aging |

> [!IMPORTANT]
> All 7 rules are correctly marked with `is_verified = FALSE` in the PostgreSQL database. Each carries the reviewer note: *"Requires tax-professional confirmation."*

---

## 6. Financial-Data Authorization & Separation of Duties

### Code Audit of Middleware & Roles
1. **Route Protection** ([`backend/cmd/api/main.go:618-628`](file:///home/poven/GO%20PROJECT/sellbuy/backend/cmd/api/main.go#L618-L628)):
   ```go
   financeAuth := mw.RequireAnyPermission("profit", "finance")
   admin.Get("/finance/forecast", financeAuth, fpaHandler.GetForecast)
   admin.Put("/finance/assumptions", financeAuth, fpaHandler.UpdateAssumptions)
   admin.Post("/finance/assumptions/reset", financeAuth, fpaHandler.ResetAssumptions)
   ```
2. **Role Domain Baseline Permissions** ([`backend/internal/models/user.go:96-101`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/models/user.go#L96-L101)):
   ```go
   RolePlatformAdmin: {
       "dashboard": true, "settings": true, "categories": true, "category_trends": true,
       "analytics": true, "locations": true, "banners": true, "notifications": true,
       "profit": true, "finance": true,
   },
   ```
3. **Findings**:
   * **Zero Separation of Duties Between Read & Write**: The same middleware (`financeAuth`) gates both read-only forecast views and write mutations (`PUT /finance/assumptions`, `POST /finance/assumptions/reset`).
   * **Over-Privileged Platform Admins**: Every user assigned `RolePlatformAdmin` is granted `"finance": true` and `"profit": true` by default in code. Consequently, infrastructure/platform engineers can alter corporate tax rates, override monetization percentages, and reset financial assumptions.
   * **Moderators & Support Agents**: Properly restricted. `RoleModerator` and `RoleSupportAgent` receive HTTP 403 Forbidden on all `/admin/finance/*` routes.
   * **Recommendation**: Split permissions into `"finance:view"` (granted to Platform Admins) and `"finance:manage"` (restricted strictly to Platform Owners and designated Finance Officers).

---

## 7. Tests and Calculations Executed

1. **Go Calculation Engine Unit Tests** ([`backend/internal/fpa/fpa_test.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/fpa_test.go)):
   * `TestUserGrowthModels`: PASS (0.00s)
   * `TestWorkloadCalculations`: PASS (0.00s)
   * `TestInfraCostEstimator`: PASS (0.00s)
   * `TestRevenueAndEscrowSeparation`: PASS (0.00s)
   * `TestEthiopianTaxCalculations`: PASS (0.00s)
   * `TestFullForecastEngineAndBreakEven`: PASS (0.00s, confirmed `Break-even month: -1`)
2. **HTTP API & Permission Integration Tests** ([`backend/internal/handler/fpa_handler_test.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/handler/fpa_handler_test.go)):
   * `Unauthenticated_Denied_401`: PASS
   * `Moderator_Denied_403`: PASS
   * `PlatformAdmin_Authorized_200`: PASS
   * `PlatformOwner_Authorized_200`: PASS
   * `Capacity_20M_Endpoint_200`: PASS
   * `CSV_Export_200`: PASS
3. **Go Static Analysis**:
   * `go vet ./...`: PASS (Exit code 0, zero warnings).
4. **Frontend Production Build**:
   * `npm run build`: PASS (Next.js 16.2.4 compiled 60/60 routes with zero TypeScript or ESLint errors).
5. **Independent Python Simulation Verification**:
   * Reproduced every monthly metric in `scratch/reconcile.py` matching Go engine output to 2 decimal places.

---

## 8. Severity-Ranked Findings

| Severity | ID | Title | Summary & Risk |
| :--- | :--- | :--- | :--- |
| **CRITICAL** | `AUD-FPA-001` | **Previous Implementation Report Contains Fabricated Forecasts** | The previous report published fictional metrics (Month 19 break-even, Month 60 2.42M users / 59M profit, 20M DAU 69k RPS / 6.8 PB). Relying on these numbers for corporate budgeting or fundraising would lead to insolvency within 14 months. |
| **HIGH** | `AUD-FPA-002` | **Base Case Model Causes Corporate Insolvency at Month 14** | Headcount ramp (scaling to 52 staff) combined with low 2% mature growth creates an 82M ETB deficit over 5 years. Break-even is never achieved under current base parameters. |
| **HIGH** | `AUD-FPA-003` | **Lack of Separation of Duties on Financial Mutations** | `RolePlatformAdmin` holds default write permissions to update and reset financial planning assumptions and tax settings via `PUT /admin/finance/assumptions`. |
| **MEDIUM** | `AUD-FPA-004` | **Loss Carryforward Pool Lacks 5-Year Vintage Aging** | Tax losses are accumulated as a perpetual scalar pool in `tax.go`, violating the 5-year expiration limit under Article 26 of Proclamation No. 979/2016. |
| **MEDIUM** | `AUD-FPA-005` | **`CommissionVATInclusive` Is Unused Dead Code** | The configuration flag is ignored in calculations, applying 15% VAT additively regardless of user preference. |
| **LOW** | `AUD-FPA-006` | **Withholding Tax Omitted from OpEx Cash Outflows** | Seeded 2% service/goods withholding tax rules are catalogued in PostgreSQL but omitted from monthly cash disbursements. |

---

## 9. Minimal Recommended Fixes

*(To be implemented only after stakeholder review and approval; read-only audit mode strictly preserved).*

1. **Re-align Base Case Planning Assumptions**:
   * Gate staff additions to revenue/MAU milestones (e.g. max 12 employees at Month 12, scaling only if revenue covers OpEx).
   * Extend early growth phase to 24 months, or adjust mature monthly growth to a sustainable 6–8%.
   * Introduce a 15M ETB Series A equity funding round in Month 12 in `CashFlowConfig` to bridge the operational investment gap.
2. **Enforce Role-Based Separation of Duties**:
   * Split `finance` permission into `finance:view` (read-only forecast access) and `finance:manage` (modifying assumptions, resetting models).
   * Remove `finance:manage` from `RoleDomainPermissions[RolePlatformAdmin]`, reserving write mutations for `RolePlatformOwner` or an explicit Finance Admin identity.
3. **Implement 5-Year Loss Relief Vintage Expiration**:
   * Update `tax.go` to store `lossTranches [5]float64`, shifting array indices at each 12-month boundary and expiring losses older than 5 years.
4. **Activate `CommissionVATInclusive` Logic**:
   * In `revenue.go`, back out VAT from gross commission when the inclusive flag is active.

---

## 10. Remaining Uncertainties and Required Human Review

The following business and legal questions require executive and professional input before final configuration:

1. **Strategic Staffing Policy**:
   * Does HUNTIQ management intend to hire 52 full-time staff by Year 5 regardless of active user adoption, or should hiring be explicitly pegged to gross platform revenue?
2. **Equity Financing Expectations**:
   * Does HUNTIQ plan to operate purely on initial seed capital (requiring severe cost discipline), or should equity investment rounds (e.g. Series A) be modeled as financing cash inflows?
3. **VAT Treatment on Seller Commissions**:
   * Is HUNTIQ's 2% commission intended to be VAT-inclusive (deducted from seller payout) or VAT-exclusive (added to the fee)?
4. **Formal Ethiopian Tax Review**:
   * The catalogued tax rules must be formally reviewed and certified by an authorized Ethiopian chartered accountant or tax attorney registered with the Ministry of Revenues before using the estimates for official tax planning or filing.

---
*End of Audit Report. Standing by for stakeholder review and approval.*
