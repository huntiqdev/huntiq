'use client';
import { useEffect, useMemo, useState } from 'react';
import { api } from '@/lib/api';
import { LISTING_DURATIONS, DEFAULT_TIER_FEE_PARAMS, calcListingFee, type SellerTier } from '@/lib/constants/fees';
import { TrendingUp, Clock, Star, Pencil } from 'lucide-react';
import { useLocale } from '@/hooks/useLocale';

interface Props {
  price: number;
  /** Seller's actual tier from their account stats */
  tier?: SellerTier;
  selectedDays: number;
  onChange: (days: number) => void;
}

const TIER_BADGE: Record<SellerTier, { label: string; color: string }> = {
  normal:  { label: 'Normal',  color: 'bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-300' },
  gold:    { label: 'Gold',    color: 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-400' },
  premium: { label: 'Premium', color: 'bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-400' },
};

/** Resolve display tier from seller's total_value_sold (passed as tier prop from stats) */
function resolveTier(tier: SellerTier): SellerTier { return tier; }

const MIN_DAYS = 1;   // 24 hours
const MAX_DAYS = 365; // 1 year

export function ListingFeeSlider({ price, tier = 'normal', selectedDays, onChange }: Props) {
  const { t } = useLocale();
  const [manualMode, setManualMode] = useState(false);
  const [manualInput, setManualInput] = useState(String(selectedDays));
  const [settings, setSettings] = useState<any>(null);

  useEffect(() => {
    api.getPublicPaymentSettings().then(setSettings).catch(() => {});
  }, []);

  const effectiveTier = resolveTier(tier);
  const breakdown = useMemo(
    () => calcListingFee(effectiveTier, price, selectedDays, settings),
    [effectiveTier, price, selectedDays, settings],
  );

  const badge = TIER_BADGE[effectiveTier];
  const params = DEFAULT_TIER_FEE_PARAMS[effectiveTier];

  const applyManual = (raw: string) => {
    const n = parseInt(raw, 10);
    if (!isNaN(n)) {
      const clamped = Math.max(MIN_DAYS, Math.min(MAX_DAYS, n));
      onChange(clamped);
      setManualInput(String(clamped));
    }
  };

  const selectPreset = (days: number) => {
    onChange(days);
    setManualInput(String(days));
    setManualMode(false);
  };

  // Human-readable duration label
  const durationLabel = (days: number) => {
    if (days === 1) return '1 day (24h)';
    if (days < 7)  return `${days} days`;
    if (days === 7) return '1 week';
    if (days < 30) return `${days} days (${Math.round(days / 7)} wks)`;
    if (days === 30) return '1 month';
    if (days < 365) return `${days} days (~${Math.round(days / 30)} mo)`;
    return '1 year';
  };

  return (
    <div className="flex flex-col gap-3 bg-zinc-50 dark:bg-zinc-800/60 rounded-2xl p-4 border border-zinc-200 dark:border-zinc-700">

      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Clock size={14} className="text-orange-600" />
          <p className="text-xs font-black uppercase tracking-widest opacity-60">{t.ui.listing_duration}</p>
        </div>
        <span className={`flex items-center gap-1 text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-full ${badge.color}`}>
          <Star size={9} /> {badge.label}
        </span>
      </div>

      {/* Preset pills */}
      <div className="grid grid-cols-5 gap-1.5">
        {LISTING_DURATIONS.map(({ days, label }) => (
          <button
            key={days}
            type="button"
            onClick={() => selectPreset(days)}
            className={`py-2 px-1 rounded-xl text-[10px] font-black text-center transition-all ${
              !manualMode && selectedDays === days
                ? 'bg-orange-600 text-white shadow-lg shadow-orange-600/20'
                : 'bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 hover:border-orange-400'
            }`}
          >
            {label}
          </button>
        ))}
      </div>

      {/* Manual date toggle */}
      <div className="flex items-center gap-2">
        <button
          type="button"
          onClick={() => setManualMode(m => !m)}
          className={`flex items-center gap-1.5 text-xs font-bold px-3 py-1.5 rounded-xl border transition-all ${
            manualMode
              ? 'bg-orange-600 text-white border-orange-600'
              : 'border-zinc-200 dark:border-zinc-700 text-zinc-500 hover:border-orange-400 hover:text-orange-600'
          }`}
        >
          <Pencil size={11} /> {t.ui.custom_days}
        </button>
        {manualMode && (
          <div className="flex items-center gap-1.5 flex-1">
            <input
              type="number"
              min={MIN_DAYS}
              max={MAX_DAYS}
              value={manualInput}
              onChange={e => setManualInput(e.target.value)}
              onBlur={e => applyManual(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && applyManual(manualInput)}
              className="w-24 border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-1.5 text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-white dark:bg-zinc-800 transition-all"
              placeholder="days"
            />
            <span className="text-xs text-zinc-400">days (1–365)</span>
          </div>
        )}
      </div>

      {/* Selected duration summary */}
      <div className="flex items-center gap-2 text-xs text-zinc-500">
        <Clock size={11} className="text-orange-600 shrink-0" />
        <span>
          {t.ui.selected_label}: <strong className="text-zinc-800 dark:text-zinc-200">{durationLabel(selectedDays)}</strong>
          {' '}— expires {new Date(Date.now() + selectedDays * 86_400_000).toLocaleDateString()}
        </span>
      </div>

      {/* Fee breakdown */}
      {price > 0 && (
        <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-700 p-3 space-y-2">
          <div className="flex items-center gap-1.5 mb-1">
            <TrendingUp size={12} className="text-orange-600" />
            <p className="text-[10px] font-black uppercase tracking-widest text-orange-600">{t.ui.fee_breakdown}</p>
          </div>

          <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs">
            <span className="opacity-50">{t.ui.startup_fee}</span>
            <span className="font-bold text-right">{breakdown.startupPct}%</span>

            <span className="opacity-50">Free weeks</span>
            <span className="font-bold text-right text-green-600">1 week free 🎉</span>

            <span className="opacity-50">Total weeks</span>
            <span className="font-bold text-right">{breakdown.totalWeeks} wk{breakdown.totalWeeks !== 1 ? 's' : ''}</span>

            <span className="opacity-50">Billable weeks</span>
            <span className="font-bold text-right">
              {breakdown.billableWeeks === 0
                ? <span className="text-green-600">0 (within free week)</span>
                : `${breakdown.billableWeeks} × ${params.weeklyIncrementPct}%`}
            </span>
          </div>

          <div className="border-t border-zinc-100 dark:border-zinc-800 pt-2 mt-1 space-y-1">
            <div className="flex justify-between text-xs">
              <span className="opacity-50">{t.ui.total_fee_pct}</span>
              <span className="font-black text-orange-600">{breakdown.totalFeePercent}%</span>
            </div>
            <div className="flex justify-between text-xs">
              <span className="opacity-50">{t.ui.fee_amount}</span>
              <span className="font-bold text-red-500">− ETB {breakdown.totalFeeAmount.toLocaleString()}</span>
            </div>
            <div className="flex justify-between text-sm font-black">
              <span>{t.ui.you_receive}</span>
              <span className="text-green-600">ETB {breakdown.sellerGets.toLocaleString()}</span>
            </div>
          </div>
        </div>
      )}

      {/* Tier info */}
      <div className="text-[10px] text-zinc-400 leading-relaxed space-y-0.5">
        {effectiveTier === 'normal' && (
          <p>💡 Sell over <strong>ETB {settings?.gold_threshold ?? '100,000'}</strong> total to unlock <strong>Gold</strong> ({settings?.gold_startup_fee ?? '1.5'}% startup, +{settings?.gold_weekly_increment ?? '0.15'}%/wk).</p>
        )}
        {effectiveTier === 'gold' && (
          <p>⭐ Sell over <strong>ETB {settings?.premium_threshold ?? '5,000,000'}</strong> total to unlock <strong>Premium</strong> ({settings?.premium_startup_fee ?? '1.2'}% startup, +{settings?.premium_weekly_increment ?? '0.10'}%/wk).</p>
        )}
        {effectiveTier === 'premium' && (
          <p>✦ <strong>Premium</strong> — lowest fees, priority listing visibility.</p>
        )}
      </div>
    </div>
  );
}
