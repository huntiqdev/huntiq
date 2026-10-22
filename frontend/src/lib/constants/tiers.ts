import { Star, Crown, type LucideIcon } from 'lucide-react';

export const TIER_RANK: Record<string, number> = {
  normal: 0,
  gold: 1,
  premium: 2,
};

/** Icon component for gold/premium tiers. Normal has no icon. */
export const TIER_ICON: Record<string, LucideIcon> = {
  gold: Star,
  premium: Crown,
};

/** Tailwind text color class per tier */
export const TIER_COLOR: Record<string, string> = {
  normal:  'text-zinc-400',
  gold:    'text-yellow-500',
  premium: 'text-purple-500',
};

/** Returns true if buyerTier meets or exceeds the required tier */
export function meetsRequiredTier(buyerTier: string, requiredTier: string): boolean {
  return (TIER_RANK[buyerTier] ?? 0) >= (TIER_RANK[requiredTier] ?? 0);
}
