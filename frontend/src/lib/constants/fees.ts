/** Platform broker fee: 2% charged on each transaction - loaded from backend at runtime */
export const BROKER_FEE_RATE = 0.02;

/** Minimum number of images required per product listing - loaded from backend */
export const MIN_IMAGES = 2;

/** Maximum image upload size in MB - loaded from backend at runtime */
export const MAX_IMAGE_SIZE_MB = 20;

/** Company contact — loaded from platform_settings at runtime, not hardcoded */
export const COMPANY_PHONE = '';
export const COMPANY_LOCATION = '';
export const COMPANY_EMAIL = '';

export interface ListingDurationOption {
  hours: number;
  label: string;
  demandTier?: 'flash' | 'urgent' | 'standard' | 'extended';
  badge?: string;
}

/** Dynamic listing durations measured in hours with demand-based tiers */
export const LISTING_DURATION_OPTIONS: ListingDurationOption[] = [
  { hours: 4,   label: '4h Flash',     demandTier: 'flash',    badge: '⚡ High Demand' },
  { hours: 24,  label: '24h Daily',    demandTier: 'urgent',   badge: '🔥 Urgent' },
  { hours: 48,  label: '48h (2 Days)', demandTier: 'standard' },
  { hours: 72,  label: '72h (3 Days)', demandTier: 'standard' },
  { hours: 168, label: '1 Week',       demandTier: 'standard' },
  { hours: 720, label: '1 Month',      demandTier: 'extended' },
];

/** Listing duration options in days (backward-compatible) */
export const LISTING_DURATIONS: { days: number; label: string }[] = [
  { days: 1,   label: '24 Hours' },
  { days: 2,   label: '2 Days' },
  { days: 7,   label: '1 Week' },
  { days: 30,  label: '1 Month' },
  { days: 180, label: '6 Months' },
  { days: 365, label: '1 Year' },
];

/** Default Seller tier thresholds (total value sold in ETB) - MUST MATCH BACKEND */
export const DEFAULT_TIER_THRESHOLDS = {
  gold:    5900,      // ✅ Matches backend DefaultGoldThreshold
  premium: 15900,     // ✅ Matches backend DefaultPremiumThreshold
} as const;

export function getTierThresholds(settings?: any) {
  if (!settings) return DEFAULT_TIER_THRESHOLDS;
  return {
    gold:    parseFloat(settings.gold_threshold ?? DEFAULT_TIER_THRESHOLDS.gold),
    premium: parseFloat(settings.premium_threshold ?? DEFAULT_TIER_THRESHOLDS.premium),
  };
}

export type SellerTier = 'normal' | 'gold' | 'premium';

/** Default Fee params per tier — MUST MATCH BACKEND platform_settings defaults */
export const DEFAULT_TIER_FEE_PARAMS: Record<SellerTier, { startupPct: number; weeklyIncrementPct: number }> = {
  normal:  { startupPct: 1.0,  weeklyIncrementPct: 0.5  },  // ✅ Matches backend defaults
  gold:    { startupPct: 0.5,  weeklyIncrementPct: 0.25 },  // ✅ Matches backend defaults
  premium: { startupPct: 0.25, weeklyIncrementPct: 0.1  },  // ✅ Matches backend defaults
};

export function getTierFeeParams(tier: SellerTier, settings?: any) {
  const d = DEFAULT_TIER_FEE_PARAMS[tier] ?? DEFAULT_TIER_FEE_PARAMS.normal;
  if (!settings) return d;
  return {
    startupPct:         parseFloat(settings[`${tier}_startup_fee`]      ?? d.startupPct),
    weeklyIncrementPct: parseFloat(settings[`${tier}_weekly_increment`] ?? d.weeklyIncrementPct),
  };
}

export function getMaxImages(tier: SellerTier, settings?: any) {
  const defaults: Record<SellerTier, number> = { normal: 2, gold: 7, premium: 20 };
  if (!settings) return defaults[tier];
  return parseInt(settings[`${tier}_max_images`] ?? defaults[tier]);
}

/** Maximum images per tier (for UI reference) */
export const MAX_IMAGES_PER_TIER: Record<SellerTier, number> = {
  normal: 2,
  gold: 7,
  premium: 20,
};

export interface FeeBreakdown {
  tier: SellerTier;
  productValue: number;
  durationDays: number;
  durationHours?: number;
  startupPct: number;
  weeklyIncrementPct: number;
  freeWeeks: number;
  billableWeeks: number;
  totalWeeks: number;
  totalFeePercent: number;
  totalFeeAmount: number;
  sellerGets: number;
}

/** Number of free weeks before weekly increment starts (matches platform_settings.listing_free_weeks) */
export const FREE_WEEKS = 1;

/**
 * Calculate demand-based listing fee measured by hours.
 * Supports custom hourly increments (e.g. 4h, 24h, 43h) with admin-controlled rate percentile.
 */
export function calcListingFeeHourly(
  tier: SellerTier,
  productValue: number,
  durationHours: number,
  settings?: any
): FeeBreakdown {
  const { startupPct, weeklyIncrementPct } = getTierFeeParams(tier, settings);

  // Hourly rate percentile is configured by admin or derived dynamically (per hour = weeklyIncrementPct / 168h)
  const hourlyRatePercentile = settings?.hourly_rate_percentile
    ? parseFloat(settings.hourly_rate_percentile)
    : (weeklyIncrementPct / 168);

  // Demand multiplier: flash/high demand (under 24h) gets optimized demand percentile
  let demandFactor = 1.0;
  if (durationHours <= 4) demandFactor = 0.6; // Discounted flash micro-duration
  else if (durationHours <= 24) demandFactor = 0.85; // Urgent daily duration

  // Startup fee covers the first 24 hours; additional hours billed by rate percentile
  const billableHours = Math.max(0, durationHours - 24);
  const incrementalFee = billableHours * hourlyRatePercentile * demandFactor;
  const totalFeePercent = Math.max(0.1, Math.round((startupPct + incrementalFee) * 100) / 100);
  const totalFeeAmount  = +(productValue * totalFeePercent / 100).toFixed(2);
  const sellerGets      = +(productValue - totalFeeAmount).toFixed(2);

  const durationDays = Math.max(1, Math.ceil(durationHours / 24));
  const totalWeeks = Math.ceil(durationDays / 7);
  const billableWeeks = Math.max(0, totalWeeks - FREE_WEEKS);

  return {
    tier,
    productValue,
    durationDays,
    durationHours,
    startupPct,
    weeklyIncrementPct,
    freeWeeks: FREE_WEEKS,
    billableWeeks,
    totalWeeks,
    totalFeePercent,
    totalFeeAmount,
    sellerGets,
  };
}

/**
 * Calculate listing fee client-side (mirrors backend CalculateListingFeeWithSettings).
 */
export function calcListingFee(
  tier: SellerTier,
  productValue: number,
  durationDays: number,
  settings?: any
): FeeBreakdown {
  return calcListingFeeHourly(tier, productValue, durationDays * 24, settings);
}
