package fpa

import (
	"math"
)

// CalculateUserGrowth produces an array of target user metric projections for horizonMonths (1-indexed).
func CalculateUserGrowth(cfg GrowthConfig) []float64 {
	horizon := cfg.ForecastHorizonMonths
	if horizon <= 0 {
		horizon = 60
	}
	if horizon > 120 {
		horizon = 120
	}

	results := make([]float64, horizon)
	initial := cfg.InitialUsers
	if initial <= 0 {
		initial = 1000.0
	}

	ceiling := cfg.MarketCeiling
	if ceiling <= 0 {
		ceiling = 20000000.0 // Default 20 million market ceiling
	}

	current := initial

	switch cfg.ModelType {
	case GrowthGeometric:
		rate := cfg.EarlyMonthlyRate
		if rate < -0.99 {
			rate = -0.99
		}
		for m := 0; m < horizon; m++ {
			if m > 0 {
				current = current * (1.0 + rate)
			}
			current = clampGrowth(current, ceiling)
			results[m] = current
		}

	case GrowthLinear:
		increment := cfg.LinearMonthlyAddition
		if increment <= 0 {
			increment = 5000.0 // Default 5k new users/mo
		}
		for m := 0; m < horizon; m++ {
			if m > 0 {
				current = current + increment
			}
			current = clampGrowth(current, ceiling)
			results[m] = current
		}

	case GrowthUserDefinedRates:
		rates := cfg.UserMonthlyRates
		for m := 0; m < horizon; m++ {
			if m > 0 {
				rate := cfg.MatureMonthlyRate
				if m < len(rates) {
					rate = rates[m]
				}
				current = current * (1.0 + rate)
			}
			current = clampGrowth(current, ceiling)
			results[m] = current
		}

	case GrowthLogistic:
		// Discrete logistic equation: U_t = U_{t-1} + r * U_{t-1} * (1 - U_{t-1} / K)
		rate := cfg.EarlyMonthlyRate
		if rate <= 0 {
			rate = 0.15
		}
		for m := 0; m < horizon; m++ {
			if m > 0 {
				delta := rate * current * (1.0 - (current / ceiling))
				if delta < 0 && current >= ceiling {
					delta = 0
				}
				current = current + delta
			}
			current = clampGrowth(current, ceiling)
			results[m] = current
		}

	case GrowthCustomMonthly:
		custom := cfg.CustomMonthlyUsers
		for m := 0; m < horizon; m++ {
			val := current
			if m < len(custom) && custom[m] > 0 {
				val = custom[m]
			} else if m > 0 {
				val = results[m-1]
			}
			val = clampGrowth(val, ceiling)
			results[m] = val
			current = val
		}

	case GrowthQuarterlySmooth:
		earlyDuration := cfg.EarlyDurationMonths
		if cfg.EarlyDurationQuarters > 0 {
			earlyDuration = cfg.EarlyDurationQuarters * 3
		}
		if earlyDuration <= 0 {
			earlyDuration = 12
		}

		earlyQRate := cfg.EarlyQuarterlyRate
		if earlyQRate == 0 && cfg.EarlyMonthlyRate > 0 {
			earlyQRate = cfg.EarlyMonthlyRate
		}
		if earlyQRate == 0 {
			earlyQRate = 0.20 // 20% quarterly growth default
		}
		// Convert quarterly rate to equivalent smooth monthly rate: (1 + r_q)^(1/3) - 1
		earlyMonthlyRate := math.Pow(1.0+earlyQRate, 1.0/3.0) - 1.0

		matureMonthlyRate := cfg.MatureMonthlyRate
		if cfg.MatureQuarterlyRate > 0 {
			matureMonthlyRate = math.Pow(1.0+cfg.MatureQuarterlyRate, 1.0/3.0) - 1.0
		}
		if matureMonthlyRate == 0 {
			matureMonthlyRate = 0.02
		}

		for m := 0; m < horizon; m++ {
			if m > 0 {
				rate := earlyMonthlyRate
				if m > earlyDuration {
					rate = matureMonthlyRate
				}
				headroom := (ceiling - current) / ceiling
				if headroom < 0.20 && headroom > 0 {
					rate = rate * (headroom / 0.20)
				}
				current = current * (1.0 + rate)
			}
			current = clampGrowth(current, ceiling)
			results[m] = current
		}

	case GrowthQuarterlyStep:
		earlyDuration := cfg.EarlyDurationMonths
		if cfg.EarlyDurationQuarters > 0 {
			earlyDuration = cfg.EarlyDurationQuarters * 3
		}
		if earlyDuration <= 0 {
			earlyDuration = 12
		}

		earlyQRate := cfg.EarlyQuarterlyRate
		if earlyQRate == 0 && cfg.EarlyMonthlyRate > 0 {
			earlyQRate = cfg.EarlyMonthlyRate
		}
		if earlyQRate == 0 {
			earlyQRate = 0.20
		}

		matureQRate := cfg.MatureQuarterlyRate
		if matureQRate == 0 && cfg.MatureMonthlyRate > 0 {
			// Convert mature monthly to quarterly: (1 + r_m)^3 - 1
			matureQRate = math.Pow(1.0+cfg.MatureMonthlyRate, 3.0) - 1.0
		}
		if matureQRate == 0 {
			matureQRate = math.Pow(1.02, 3.0) - 1.0 // ~6.12% per quarter
		}

		for m := 0; m < horizon; m++ {
			// Step growth applied at start of each new quarter (m=3, m=6, m=9, ...)
			if m > 0 && m%3 == 0 {
				rate := earlyQRate
				if m > earlyDuration {
					rate = matureQRate
				}
				headroom := (ceiling - current) / ceiling
				if headroom < 0.20 && headroom > 0 {
					rate = rate * (headroom / 0.20)
				}
				current = current * (1.0 + rate)
			}
			current = clampGrowth(current, ceiling)
			results[m] = current
		}

	case GrowthStableTransition:
		fallthrough
	default:
		// Transition model: early rapid adoption for N months, then mature growth capped at ceiling
		earlyDuration := cfg.EarlyDurationMonths
		if cfg.EarlyDurationQuarters > 0 {
			earlyDuration = cfg.EarlyDurationQuarters * 3
		}
		if earlyDuration <= 0 {
			earlyDuration = 12
		}

		earlyRate := cfg.EarlyMonthlyRate
		if cfg.GrowthPeriodType == "quarterly" || (cfg.EarlyQuarterlyRate > 0 && earlyRate == 0) {
			qRate := cfg.EarlyQuarterlyRate
			if qRate == 0 {
				qRate = earlyRate
			}
			if qRate == 0 {
				qRate = 0.20
			}
			earlyRate = math.Pow(1.0+qRate, 1.0/3.0) - 1.0
		} else if earlyRate == 0 {
			earlyRate = 0.20 // 20% early monthly growth
		}

		matureRate := cfg.MatureMonthlyRate
		if cfg.GrowthPeriodType == "quarterly" && cfg.MatureQuarterlyRate > 0 {
			matureRate = math.Pow(1.0+cfg.MatureQuarterlyRate, 1.0/3.0) - 1.0
		} else if matureRate == 0 {
			matureRate = 0.02 // 2% mature monthly growth
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
	}

	return results
}

// ConvertTargetMetric translates the raw growth series into MAU, DAU, and cumulative registered users
func ConvertTargetMetric(rawUsers float64, target MetricMAUOrDAU, dauToMAURatio float64) (mau float64, dau float64, cumulative float64) {
	if dauToMAURatio <= 0 || dauToMAURatio > 1.0 {
		dauToMAURatio = 0.25 // 25% default
	}

	activeToRegisteredRatio := 0.40 // 40% of registered accounts are active in any month

	switch target {
	case MetricDAU:
		dau = rawUsers
		mau = rawUsers / dauToMAURatio
		cumulative = mau / activeToRegisteredRatio
	case MetricRegisteredUsers:
		cumulative = rawUsers
		mau = cumulative * activeToRegisteredRatio
		dau = mau * dauToMAURatio
	case MetricMAU:
		fallthrough
	default:
		mau = rawUsers
		dau = mau * dauToMAURatio
		cumulative = mau / activeToRegisteredRatio
	}

	return mau, dau, cumulative
}

// MetricMAUOrDAU alias for TargetMetric
type MetricMAUOrDAU = TargetMetric

func clampGrowth(val, ceiling float64) float64 {
	if math.IsNaN(val) || math.IsInf(val, 0) || val < 0 {
		return 0
	}
	if val > ceiling {
		return ceiling
	}
	return val
}
