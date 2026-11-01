/**
 * Shared Tailwind class strings for consistent form UI across the app.
 * Import the constant you need — never copy-paste these strings.
 */

/** Standard light/dark input field */
export const INPUT_CLS =
  'w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2.5 text-base sm:text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 transition-all';

/** Input with orange focus ring — used on light backgrounds */
export const INPUT_CLS_LIGHT =
  'w-full bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2.5 text-base sm:text-sm text-zinc-900 dark:text-white placeholder-zinc-400 outline-none focus:border-orange-500 transition-all';

/** Input for dark modal backgrounds (TierUpgradeModal) */
export const INPUT_CLS_DARK =
  'w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-base sm:text-sm text-white placeholder-white/30 outline-none focus:border-orange-500/60 focus:ring-1 focus:ring-orange-500/30 transition-all';

/** Shared label style for form sections */
export const LABEL_CLS = 'text-xs font-black uppercase tracking-wider opacity-60';

/** Select / dropdown — same border as INPUT_CLS but py-2 */
export const SELECT_CLS =
  'w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-base sm:text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 transition-all';
