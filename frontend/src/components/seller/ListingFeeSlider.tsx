'use client';
import { useEffect, useMemo, useState } from 'react';
import { api } from '@/lib/api';
import {
  LISTING_DURATION_OPTIONS,
  DEFAULT_TIER_FEE_PARAMS,
  calcListingFeeHourly,
  type SellerTier,
  type ListingDurationOption,
} from '@/lib/constants/fees';
import { TrendingUp, Clock, Star, Pencil, Zap } from 'lucide-react';
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

function resolveTier(tier: SellerTier): SellerTier { return tier; }

const MIN_HOURS = 1;
const MAX_HOURS = 8760; // 1 year in hours

export function ListingFeeSlider({ price, tier = 'normal', selectedDays, onChange }: Props) {
  const { t } = useLocale();
  const [selectedHours, setSelectedHours] = useState(() => Math.max(4, selectedDays * 24));
  const [manualMode, setManualMode] = useState(false);
  const [manualInput, setManualInput] = useState(String(selectedHours));
  const [settings, setSettings] = useState<any>(null);

  useEffect(() => {
    api.getPublicPaymentSettings().then(setSettings).catch(() => {});
  }, []);

  // Synchronize when parent passes updated selectedDays
  useEffect(() => {
    const h = Math.max(1, selectedDays * 24);
    setSelectedHours(h);
    setManualInput(String(h));
  }, [selectedDays]);

  const effectiveTier = resolveTier(tier);
  const breakdown = useMemo(
    () => calcListingFeeHourly(effectiveTier, price, selectedHours, settings),
    [effectiveTier, price, selectedHours, settings],
  );

  const badge = TIER_BADGE[effectiveTier];
  const params = DEFAULT_TIER_FEE_PARAMS[effectiveTier];

  const applyManualHours = (raw: string) => {
    const n = parseInt(raw, 10);
    if (!isNaN(n)) {
      const clamped = Math.max(MIN_HOURS, Math.min(MAX_HOURS, n));
      setSelectedHours(clamped);
      onChange(Math.max(1, Math.ceil(clamped / 24)));
      setManualInput(String(clamped));
    }
  };

  const selectPreset = (opt: ListingDurationOption) => {
    setSelectedHours(opt.hours);
    onChange(Math.max(1, Math.ceil(opt.hours / 24)));
    setManualInput(String(opt.hours));
    setManualMode(false);
  };

  const durationLabel = (hours: number) => {
    if (hours < 24) return `${hours} Hours (Demand Flash)`;
    const days = Math.round(hours / 24);
    if (days === 1) return '24 Hours (1 Day)';
    if (days < 7)  return `${hours} Hours (${days} Days)`;
    if (days === 7) return '1 Week (168 Hours)';
    if (days < 30) return `${days} Days (${Math.round(days / 7)} Weeks)`;
    if (days === 30) return '1 Month (720 Hours)';
    return `${days} Days`;
  };

  const [baseTime] = useState(() => Date.now());
  const expiryDateStr = useMemo(() => {
    return new Date(baseTime + selectedHours * 3_600_000).toLocaleString();
  }, [baseTime, selectedHours]);

  return (
    <div className="flex flex-col gap-3 bg-zinc-50 dark:bg-zinc-800/60 rounded-2xl p-4 border border-zinc-200 dark:border-zinc-700">

      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Clock size={14} className="text-orange-600" />
          <p className="text-xs font-black uppercase tracking-widest opacity-60">Demand-Based Duration (Hourly)</p>
        </div>
        <span className={`flex items-center gap-1 text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-full ${badge.color}`}>
          <Star size={9} /> {badge.label}
        </span>
      </div>

      {/* Preset pills */}
      <div className="grid grid-cols-3 sm:grid-cols-6 gap-1.5">
        {LISTING_DURATION_OPTIONS.map((opt) => (
          <button
            key={opt.hours}
            type="button"
            onClick={() => selectPreset(opt)}
            className={`py-2 px-1 rounded-xl text-[10px] font-black text-center transition-all ${
              !manualMode && selectedHours === opt.hours
                ? 'bg-orange-600 text-white shadow-lg shadow-orange-600/20'
                : 'bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 hover:border-orange-400'
            }`}
          >
            <div>{opt.label}</div>
            {opt.badge && (
              <div className="text-[8px] opacity-80 font-normal truncate mt-0.5">{opt.badge}</div>
            )}
          </button>
        ))}
      </div>

      {/* Manual hours toggle */}
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
          <Pencil size={11} /> Custom Hours
        </button>
        {manualMode && (
          <div className="flex items-center gap-1.5 flex-1">
            <input
              type="number"
              min={MIN_HOURS}
              max={MAX_HOURS}
              value={manualInput}
              onChange={e => setManualInput(e.target.value)}
              onBlur={e => applyManualHours(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && applyManualHours(manualInput)}
              className="w-24 border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-1.5 text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-white dark:bg-zinc-800 transition-all"
              placeholder="e.g. 4, 24, 43"
            />
            <span className="text-xs text-zinc-400">hours (e.g. 4h, 24h, 43h)</span>
          </div>
        )}
      </div>

      {/* Selected duration summary */}
      <div className="flex items-center gap-2 text-xs text-zinc-500">
        <Clock size={11} className="text-orange-600 shrink-0" />
        <span>
          Selected: <strong className="text-zinc-800 dark:text-zinc-200">{durationLabel(selectedHours)}</strong>
          {' '}— active until {expiryDateStr}
        </span>
      </div>

      {/* Fee breakdown */}
      {price > 0 && (
        <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-700 p-3 space-y-2">
          <div className="flex items-center gap-1.5 mb-1">
            <TrendingUp size={12} className="text-orange-600" />
            <p className="text-[10px] font-black uppercase tracking-widest text-orange-600">Demand-Based Dynamic Fee</p>
          </div>

          <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs">
            <span className="opacity-50">Startup Rate</span>
            <span className="font-bold text-right">{breakdown.startupPct}%</span>

            <span className="opacity-50">Duration Hours</span>
            <span className="font-bold text-right">{selectedHours} Hours</span>

            <span className="opacity-50">Rate Calculation</span>
            <span className="font-bold text-right text-sky-600">
              {selectedHours <= 4 ? '⚡ Flash Discounted' : selectedHours <= 24 ? '🔥 Daily Rate' : 'Hourly Percentile'}
            </span>
          </div>

          <div className="border-t border-zinc-100 dark:border-zinc-800 pt-2 mt-1 space-y-1">
            <div className="flex justify-between text-xs">
              <span className="opacity-50">Calculated Fee Rate</span>
              <span className="font-black text-orange-600">{breakdown.totalFeePercent}%</span>
            </div>
            <div className="flex justify-between text-xs">
              <span className="opacity-50">Total Fee Amount</span>
              <span className="font-bold text-red-500">− ETB {breakdown.totalFeeAmount.toLocaleString()}</span>
            </div>
            <div className="flex justify-between text-sm font-black">
              <span>You Receive</span>
              <span className="text-green-600">ETB {breakdown.sellerGets.toLocaleString()}</span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
