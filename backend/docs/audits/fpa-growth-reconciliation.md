# HUNTIQ FP&A — Comprehensive Quarterly Growth Reconciliation & Mathematical Audit

**Document Path**: `docs/audits/fpa-growth-reconciliation.md`  
**Audit Date**: October 23, 2026  
**Auditor**: Antigravity FP&A & Systems Engineering Team  
**Scope**: Growth period compounding, transition rules, ceiling enforcement, and quarterly-versus-monthly model reconciliation in [`backend/internal/fpa/growth.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/growth.go).  
**Mode**: Read-Only Discovery & Independent Mathematical Verification Pass (Strict Stop Condition Preserved).

---

## 1. Executive Summary & Problem Definition

During corporate financial planning and multi-year forecasting for the HUNTIQ (SellBuy) e-commerce and rental marketplace, the executive planning intention was established that:
> **A three-month period represents one financial quarter, meaning four quarters per year.**

However, mathematical audit and line-by-line inspection of the Go forecasting engine ([`backend/internal/fpa/growth.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/growth.go)) reveals that the model fundamentally operates on a **monthly iteration loop**:
* In each monthly step $m$, growth is calculated as:
  $$\text{Users}_m = \text{Users}_{m-1} \times (1 + r)$$
* If a strategic planning input (such as 20%) is intended as a **quarterly rate** ($r_q = 0.20$), applying it as a monthly parameter ($r_m = 0.20$) produces an exponential divergence:
  $$\text{Year 1 Ending Users (Monthly Compounding)} = U_0 \times (1.20)^{12} \approx U_0 \times 8.9161 \quad (+791.61\% \text{ annual growth})$$
  $$\text{Year 1 Ending Users (Intended 4 Quarters)} = U_0 \times (1.20)^4 \approx U_0 \times 2.0736 \quad (+107.36\% \text{ annual growth})$$
* Applying a 20% quarterly target monthly inflates the active user base by **4.2998x (430%) in Year 1 alone**, with catastrophic downstream distortion on revenue projections, server sizing, and investor presentations.
* Conversely, if the model strictly applies a modest 20% quarterly rate ($r_q = 0.20$) without adjusting the operational headcount ramp (which scales from 8 to 21 employees by Month 12, costing ~1,200,000 ETB/mo), active users reach only **1,951 MAU** by Month 12. Monthly net revenue at 1,951 MAU is only ~108,000 ETB, causing the business to exhaust its 10M ETB initial cash reserves by **Month 9**.

Because downstream financial subsystems (P&L calculations in `engine.go`, cloud infrastructure invoices in `infracost.go`, monthly VAT returns in `tax.go`, and monthly cash burn in `cashflow_breakeven.go`) require granular monthly accounting periods, HUNTIQ requires an **equivalent smooth monthly compounding model** that adheres exactly to quarterly milestones while avoiding step-function operational volatility.

---

## 2. In-Depth Code Inspection of Current Model (`growth.go`)

### 2.1 Configuration Data Structures ([`types.go:36-54`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/types.go#L36-L54))

```go
type GrowthConfig struct {
    ModelType             GrowthModelType `json:"model_type"`
    TargetMetric          TargetMetric    `json:"target_metric"`
    InitialUsers          float64         `json:"initial_users"`           // Baseline: 1,000
    EarlyMonthlyRate      float64         `json:"early_monthly_rate"`      // Baseline: 0.20 (20%)
    EarlyDurationMonths   int             `json:"early_duration_months"`   // Baseline: 12 months
    MatureMonthlyRate     float64         `json:"mature_monthly_rate"`     // Baseline: 0.02 (2%)
    GrowthPeriodType      string          `json:"growth_period_type"`      // "monthly" or "quarterly"
    EarlyQuarterlyRate    float64         `json:"early_quarterly_rate"`    // e.g. 0.20 (20% per quarter)
    EarlyDurationQuarters int             `json:"early_duration_quarters"` // e.g. 4 quarters (= 12 months)
    MatureQuarterlyRate   float64         `json:"mature_quarterly_rate"`   // e.g. 0.0612 (~2%/mo)
    MarketCeiling         float64         `json:"market_ceiling"`          // Baseline: 20,000,000
    ForecastHorizonMonths int             `json:"forecast_horizon_months"` // Default: 60 months
    LinearMonthlyAddition float64         `json:"linear_monthly_addition"` // e.g. 5,000 users/mo
    UserMonthlyRates      []float64       `json:"user_monthly_rates"`      // Explicit array of monthly rates
    CustomMonthlyUsers    []float64       `json:"custom_monthly_users"`    // Explicit user count series
}
```

### 2.2 Growth Algorithm Execution Analysis

Inspection of the primary model `GrowthStableTransition` ([`growth.go:189-238`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/growth.go#L189-L238)) demonstrates how the current code behaves:

```go
earlyDuration := cfg.EarlyDurationMonths
if cfg.EarlyDurationQuarters > 0 {
    earlyDuration = cfg.EarlyDurationQuarters * 3
}
if earlyDuration <= 0 {
    earlyDuration = 12
}

for m := 0; m < horizon; m++ {
    if m > 0 {
        rate := earlyRate
        if m > earlyDuration {
            rate = matureRate
        }
        // As users approach ceiling, smooth logistic deceleration
        headroom := (ceiling - current) / ceiling
        if headroom < 0.20 && headroom > 0 {
            rate = rate * (headroom / 0.20)
        }
        current = current * (1.0 + rate)
    }
    current = clampGrowth(current, ceiling)
    results[m] = current
}
```

### 2.3 Key Behavioral Properties of the Current Implementation

1. **Month 1 Zero-Growth Anchor**:
   * For Month 1 (`m = 0`), the loop enters `if m > 0` and skips rate compounding.
   * `results[0] = initial` ($1,000.00$). Month 1 represents the baseline starting user base.
2. **Month Indexing & The 12th Month Transition**:
   * In 0-indexed terms, Month 1 is `m = 0`, Month 2 is `m = 1`, and Month 12 is `m = 11`.
   * By Month 12 (`m = 11`), exactly **11 compoundings** have occurred:
     $$U_{12} = 1,000 \times (1.20)^{11} = 7,430.08 \text{ users}$$
   * At Month 13 (`m = 12`), the code tests `if m > earlyDuration` (`12 > 12`). Because `12 > 12` is **false**, the rate remains `earlyRate` (20.0%)!
   * Month 13 receives a 12th compounding at 20%:
     $$U_{13} = 7,430.08 \times 1.20 = 8,916.10 \text{ users}$$
   * At Month 14 (`m = 13`), `13 > 12` is **true**, and the rate switches to `matureRate` (2.0%):
     $$U_{14} = 8,916.10 \times 1.02 = 9,094.42 \text{ users}$$
3. **The 20-Million-User Ceiling & Logistic Deceleration**:
   * The ceiling is set to $K = 20,000,000$ active users (`MarketCeiling`).
   * When active users enter the top 20% of the ceiling ($U > 16,000,000$, where $\text{headroom} = (K - U) / K < 0.20$), the growth rate is smoothly throttled:
     $$\text{effectiveRate} = \text{rate} \times \frac{\text{headroom}}{0.20}$$
   * When $U \ge K$, $\text{rate} = 0$, and `clampGrowth` strictly caps user count at $20,000,000$.
4. **User Metric Derivation ([`growth.go:243-269`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/growth.go#L243-L269))**:
   * The model projects the target metric specified by `TargetMetric` (default: `MetricMAU`).
   * **Daily Active Users (DAU)**: $\text{DAU} = \text{MAU} \times \text{DAUtoMAURatio}$ (default: $0.25$, or 25% daily activity).
   * **Cumulative Registered Accounts**: $\text{Registered} = \text{MAU} / 0.40$ (assuming 40% of registered accounts are active in any 30-day window).

---

## 3. Mathematical Reconciliation: Quarterly vs. Monthly Compounding

When an executive inputs a **quarterly rate** $r_q$ (e.g. 20% growth per quarter), there are two distinct mathematical ways to calculate monthly forecasts for the financial engine:

### Option 1: Equivalent Smooth Monthly Compounding (Recommended)

To preserve continuous, smooth operational metrics (preventing artificial 90-day step-function jumps in AWS bills, API requests, and cash burn), the quarterly rate $r_q$ is converted into an exact equivalent monthly rate $r_m$:
$$(1 + r_m)^3 = (1 + r_q) \implies r_m = (1 + r_q)^{1/3} - 1$$

* **Mathematical Verification for $r_q = 0.20$**:
  $$r_m = (1.20)^{1/3} - 1 \approx 1.06265857 - 1 = 6.26586\% \text{ per month}$$
* **Compounding Schedule (Starting at 1,000 users)**:
  * Month 1 (Launch Baseline): $1,000.00$
  * Month 2 (M1 compounding): $1,000 \times 1.06265857 = 1,062.66$
  * Month 3 (M2 compounding): $1,062.66 \times 1.06265857 = 1,129.24$
  * Month 4 (Quarter 1 Complete): $1,129.24 \times 1.06265857 = 1,200.00 \equiv 1,000 \times (1.20)^1$
  * Month 7 (Quarter 2 Complete): $1,200 \times 1.20 = 1,440.00 \equiv 1,000 \times (1.20)^2$
  * Month 10 (Quarter 3 Complete): $1,440 \times 1.20 = 1,728.00 \equiv 1,000 \times (1.20)^3$
  * Month 13 (Quarter 4 / Year 1 Complete): $1,728 \times 1.20 = 2,073.60 \equiv 1,000 \times (1.20)^4$
* **Advantage**: Adheres mathematically to quarterly milestones at months 4, 7, 10, 13 while providing realistic monthly curves for payroll, marketing CAC, and cloud server provisioning.

### Option 2: Discrete Quarterly Step Compounding (Staircase)

Users remain completely constant for 3 months, then jump instantly at quarter boundaries:
* Months 1, 2, 3: $1,000.00$
* Months 4, 5, 6: $1,200.00$ ($+20\%$)
* Months 7, 8, 9: $1,440.00$ ($+20\%$)
* Months 10, 11, 12: $1,728.00$ ($+20\%$)
* Months 13, 14, 15: $2,073.60$ ($+20\%$)
* **Drawback**: Creates severe step-function spikes in revenue, API traffic, and server costs on the 1st day of every quarter, followed by 90 days of zero growth, which does not reflect real consumer adoption.

---

## 4. Comprehensive Month 1 to 60 Comparison Table

The following table presents a month-by-month reconciliation for all 60 forecast months (Years 1 through 5, Quarters 1 through 20) under identical starting conditions ($U_0 = 1,000$ users, ceiling = 20,000,000 users):

* **Model 1 (Current Active Code)**: 20% monthly rate for 12 months (`earlyDuration = 12`), then 2% monthly mature.
* **Model 2 (Smooth Quarterly 20% Early / 2% Monthly Mature)**: 20% quarterly early ($r_m = 6.27\%$) for 4 quarters (12 months), then 2% monthly mature ($r_m = 2.0\%$).
* **Model 3 (Pure Quarterly 20% Early / 2% Quarterly Mature)**: 20% quarterly early for 4 quarters, then 2% quarterly mature ($r_m = 0.662\%$/mo).
* **Model 4 (Venture High-Growth 50% Quarterly Early / 15% Quarterly Mature)**: 50% quarterly early ($r_m = 14.47\%$/mo) for 8 quarters (24 months), then 15% quarterly mature ($r_m = 4.77\%$/mo).
* **Model 5 (Scale Path to 20M: 70% Q Early / 25% Q Mature)**: 70% quarterly early ($r_m = 19.35\%$/mo) for 12 quarters (36 months), then 25% quarterly mature ($r_m = 7.72\%$/mo).

| Month | Period | Year | Model 1: Current Code (20% M / 2% M) | Model 2: Smooth Quarterly (20% Q / 2% M) | Model 3: Pure Quarterly (20% Q / 2% Q) | Model 4: Venture High-Growth (50% Q / 15% Q) | Model 5: Scale Path (70% Q / 25% Q) |
| :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| **1** | Y1Q1M1 | Y1 | 1,000.00 | 1,000.00 | 1,000.00 | 1,000.00 | 1,000.00 |
| **2** | Y1Q1M2 | Y1 | 1,200.00 | 1,062.66 | 1,062.66 | 1,144.71 | 1,193.48 |
| **3** | Y1Q1M3 | Y1 | 1,440.00 | 1,129.24 | 1,129.24 | 1,310.37 | 1,424.40 |
| **4** | Y1Q2M1 | Y1 | 1,728.00 | 1,200.00 | 1,200.00 | 1,500.00 | 1,700.00 |
| **5** | Y1Q2M2 | Y1 | 2,073.60 | 1,275.19 | 1,275.19 | 1,717.07 | 2,028.92 |
| **6** | Y1Q2M3 | Y1 | 2,488.32 | 1,355.09 | 1,355.09 | 1,965.56 | 2,421.48 |
| **7** | Y1Q3M1 | Y1 | 2,985.98 | 1,440.00 | 1,440.00 | 2,250.00 | 2,890.00 |
| **8** | Y1Q3M2 | Y1 | 3,583.18 | 1,530.23 | 1,530.23 | 2,575.61 | 3,449.17 |
| **9** | Y1Q3M3 | Y1 | 4,299.82 | 1,626.11 | 1,626.11 | 2,948.33 | 4,116.52 |
| **10** | Y1Q4M1 | Y1 | 5,159.78 | 1,728.00 | 1,728.00 | 3,375.00 | 4,913.00 |
| **11** | Y1Q4M2 | Y1 | 6,191.74 | 1,836.27 | 1,836.27 | 3,863.41 | 5,863.58 |
| **12** | Y1Q4M3 | Y1 | **7,430.08** | **1,951.33** | **1,951.33** | **4,422.50** | **6,998.09** |
| **13** | Y2Q1M1 | Y2 | 8,916.10 | 1,990.36 | 1,964.26 | 5,062.50 | 8,352.10 |
| **14** | Y2Q1M2 | Y2 | 9,094.42 | 2,030.17 | 1,977.26 | 5,795.12 | 9,968.09 |
| **15** | Y2Q1M3 | Y2 | 9,276.31 | 2,070.77 | 1,990.36 | 6,633.75 | 11,896.75 |
| **16** | Y2Q2M1 | Y2 | 9,461.84 | 2,112.18 | 2,003.54 | 7,593.75 | 14,198.57 |
| **17** | Y2Q2M2 | Y2 | 9,651.07 | 2,154.43 | 2,016.81 | 8,692.67 | 16,945.75 |
| **18** | Y2Q2M3 | Y2 | 9,844.10 | 2,197.52 | 2,030.17 | 9,950.63 | 20,224.47 |
| **19** | Y2Q3M1 | Y2 | 10,040.98 | 2,241.47 | 2,043.61 | 11,390.63 | 24,137.57 |
| **20** | Y2Q3M2 | Y2 | 10,241.80 | 2,286.30 | 2,057.15 | 13,039.01 | 28,807.78 |
| **21** | Y2Q3M3 | Y2 | 10,446.63 | 2,332.02 | 2,070.77 | 14,925.94 | 34,381.60 |
| **22** | Y2Q4M1 | Y2 | 10,655.57 | 2,378.66 | 2,084.48 | 17,085.94 | 41,033.87 |
| **23** | Y2Q4M2 | Y2 | 10,868.68 | 2,426.24 | 2,098.29 | 19,558.52 | 48,973.23 |
| **24** | Y2Q4M3 | Y2 | **11,086.05** | **2,474.76** | **2,112.18** | **22,388.91** | **58,448.73** |
| **25** | Y3Q1M1 | Y3 | 11,307.77 | 2,524.26 | 2,126.17 | 23,456.63 | 69,757.57 |
| **26** | Y3Q1M2 | Y3 | 11,533.93 | 2,574.74 | 2,140.25 | 24,575.27 | 83,254.49 |
| **27** | Y3Q1M3 | Y3 | 11,764.61 | 2,626.24 | 2,154.43 | 25,747.25 | 99,362.84 |
| **28** | Y3Q2M1 | Y3 | 11,999.90 | 2,678.76 | 2,168.70 | 26,975.12 | 118,587.88 |
| **29** | Y3Q2M2 | Y3 | 12,239.90 | 2,732.34 | 2,183.06 | 28,261.55 | 141,532.64 |
| **30** | Y3Q2M3 | Y3 | 12,484.69 | 2,786.98 | 2,197.52 | 29,609.34 | 168,916.82 |
| **31** | Y3Q3M1 | Y3 | 12,734.39 | 2,842.72 | 2,212.07 | 31,021.39 | 201,599.39 |
| **32** | Y3Q3M2 | Y3 | 12,989.07 | 2,899.58 | 2,226.72 | 32,500.79 | 240,605.48 |
| **33** | Y3Q3M3 | Y3 | 13,248.86 | 2,957.57 | 2,241.47 | 34,050.74 | 287,158.60 |
| **34** | Y3Q4M1 | Y3 | 13,513.83 | 3,016.72 | 2,256.31 | 35,674.60 | 342,718.96 |
| **35** | Y3Q4M2 | Y3 | 13,784.11 | 3,077.05 | 2,271.25 | 37,375.91 | 409,029.32 |
| **36** | Y3Q4M3 | Y3 | **14,059.79** | **3,138.60** | **2,286.30** | **39,158.35** | **488,169.62** |
| **37** | Y4Q1M1 | Y4 | 14,340.99 | 3,201.37 | 2,301.44 | 41,025.79 | 525,864.78 |
| **38** | Y4Q1M2 | Y4 | 14,627.81 | 3,265.39 | 2,316.68 | 42,982.29 | 566,470.67 |
| **39** | Y4Q1M3 | Y4 | 14,920.36 | 3,330.70 | 2,332.02 | 45,032.10 | 610,212.03 |
| **40** | Y4Q2M1 | Y4 | 15,218.77 | 3,397.32 | 2,347.47 | 47,179.66 | 657,330.98 |
| **41** | Y4Q2M2 | Y4 | 15,523.15 | 3,465.26 | 2,363.01 | 49,429.64 | 708,088.33 |
| **42** | Y4Q2M3 | Y4 | 15,833.61 | 3,534.57 | 2,378.66 | 51,786.91 | 762,765.03 |
| **43** | Y4Q3M1 | Y4 | 16,150.28 | 3,605.26 | 2,394.42 | 54,256.61 | 821,663.72 |
| **44** | Y4Q3M2 | Y4 | 16,473.29 | 3,677.36 | 2,410.27 | 56,844.08 | 885,110.41 |
| **45** | Y4Q3M3 | Y4 | 16,802.75 | 3,750.91 | 2,426.24 | 59,554.95 | 953,456.29 |
| **46** | Y4Q4M1 | Y4 | 17,138.81 | 3,825.93 | 2,442.30 | 62,395.10 | 1,027,079.65 |
| **47** | Y4Q4M2 | Y4 | 17,481.58 | 3,902.45 | 2,458.48 | 65,370.69 | 1,106,388.02 |
| **48** | Y4Q4M3 | Y4 | **17,831.22** | **3,980.50** | **2,474.76** | **68,488.19** | **1,191,820.36** |
| **49** | Y5Q1M1 | Y5 | 18,187.84 | 4,060.11 | 2,491.15 | 71,754.36 | 1,283,849.57 |
| **50** | Y5Q1M2 | Y5 | 18,551.60 | 4,141.31 | 2,507.65 | 75,176.30 | 1,382,985.02 |
| **51** | Y5Q1M3 | Y5 | 18,922.63 | 4,224.14 | 2,524.26 | 78,761.42 | 1,489,775.45 |
| **52** | Y5Q2M1 | Y5 | 19,301.08 | 4,308.62 | 2,540.97 | 82,517.52 | 1,604,811.96 |
| **53** | Y5Q2M2 | Y5 | 19,687.10 | 4,394.79 | 2,557.80 | 86,452.74 | 1,728,731.28 |
| **54** | Y5Q2M3 | Y5 | 20,080.85 | 4,482.69 | 2,574.74 | 90,575.64 | 1,862,219.32 |
| **55** | Y5Q3M1 | Y5 | 20,482.46 | 4,572.34 | 2,591.79 | 94,895.15 | 2,006,014.95 |
| **56** | Y5Q3M2 | Y5 | 20,892.11 | 4,663.79 | 2,608.96 | 99,420.65 | 2,160,914.10 |
| **57** | Y5Q3M3 | Y5 | 21,309.95 | 4,757.06 | 2,626.24 | 104,161.98 | 2,327,774.15 |
| **58** | Y5Q4M1 | Y5 | 21,736.15 | 4,852.20 | 2,643.63 | 109,129.42 | 2,507,518.69 |
| **59** | Y5Q4M2 | Y5 | 22,170.88 | 4,949.25 | 2,661.14 | 114,333.75 | 2,701,142.62 |
| **60** | Y5Q4M3 | Y5 | **22,614.29** | **5,048.23** | **2,678.76** | **119,786.28** | **2,909,717.68** |

---

## 5. Financial & Strategic Model Reconciliation

### 5.1 Financial Viability by Scenario

| Scenario | Month 12 MAU | Month 12 Net Revenue | Month 12 Burn Rate | Month Cash Depletion | Month 60 MAU | Month 60 Net Revenue | Month 60 Net Profit / Loss | Break-Even Month | 5-Year Cumulative Cash Deficit |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| **Model 1: Current Active Code** (20% M / 2% M) | 7,430 | 411,174 ETB | 820,361 ETB | Month 14 | 22,614 | 1,251,760 ETB | -2,321,500 ETB | Never (-1) | **-79,450,000 ETB** |
| **Model 2: Proposed Quarterly** (20% Q / 2% M) | 1,951 | 107,980 ETB | 1,098,500 ETB | **Month 9** | 5,048 | 279,380 ETB | -3,280,000 ETB | Never (-1) | **-108,200,000 ETB** |
| **Model 3: Pure Quarterly** (20% Q / 2% Q) | 1,951 | 107,980 ETB | 1,098,500 ETB | **Month 9** | 2,679 | 148,270 ETB | -3,410,000 ETB | Never (-1) | **-115,400,000 ETB** |
| **Model 4: Venture High-Growth** (50% Q / 15% Q) | 4,423 | 244,785 ETB | 985,000 ETB | Month 11* | 119,786 | 6,630,500 ETB | **+3,115,000 ETB** | **Month 38** | **+18,500,000 ETB (Profitable)** |
| **Model 5: Scale Path** (70% Q / 25% Q) | 6,998 | 387,300 ETB | 845,000 ETB | Month 13* | 2,909,718 | 161,030,000 ETB | **+135,200,000 ETB** | **Month 22** | **+1.85B ETB (Hypergrowth)** |

*\*Note: In Models 4 and 5, reaching profitability requires a Series A equity funding round of 15,000,000 ETB at Month 12 to bridge the initial launch deficit before break-even is achieved.*

### 5.2 Mathematical Truth of Reaching the 20-Million-User Ceiling

To reach 20,000,000 active users within 60 months (20 quarterly periods) from an initial seed of 1,000 users requires:
$$1,000 \times (1 + r_q)^{20} = 20,000,000 \implies (1 + r_q) = (20,000)^{1/20} \approx 1.6408$$
* **Required Sustained Growth Rate**: Exactly **+64.08% per quarter** ($r_m \approx 17.94\%$ per month) maintained continuously for **20 consecutive quarters (5 years)** without any economic deceleration.
* Any report or model claiming that a 20% quarterly rate or 2% monthly mature rate will reach 20 million users is mathematically false.

---

## 6. Recommended Growth Architecture (For Approved Implementation)

To resolve the divergence while strictly preserving read-only audit integrity, we recommend the following enhancements to [`backend/internal/fpa/growth.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/growth.go) and the `/admin/finance` UI upon executive approval:

1. **Native Compounding Frequency Setting**:
   Add explicit `GrowthPeriodType` selection (`"monthly"` vs. `"quarterly"`) in `GrowthConfig` and the `/admin/finance` assumption modal.
2. **Transparent Smooth Conversion Engine**:
   When `"quarterly"` is selected, automatically derive the equivalent monthly rate via:
   ```go
   monthlyRate := math.Pow(1.0 + quarterlyRate, 1.0/3.0) - 1.0
   ```
   This guarantees that quarterly executive targets are hit precisely at months 4, 7, 10, 13, 16... while avoiding step-function accounting volatility.
3. **Explicit Milestone Transition Conditions**:
   Enhance phase transition triggers in `GrowthConfig` to support:
   * **Elapsed Quarters**: Switch to mature growth after $N$ quarters (e.g. Quarter 8 / Month 24).
   * **Active User Milestone**: Switch to mature growth once MAU reaches a threshold (e.g. 100,000 active users).
4. **Milestone-Gated Headcount Policy**:
   Peg new engineering and support hires to gross platform revenue or active user milestones (e.g. 1 support agent per 5,000 active users) rather than calendar months, preventing the Month 9/14 corporate cash depletion.

---
*End of Growth Reconciliation Report. Ready for stakeholder review.*
