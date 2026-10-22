import en from './language/en.json';
import am from './language/am.json';
import or_ from './language/or.json';
import ti from './language/ti.json';
// ── To add a new language ────────────────────────────────────────────────────
// 1. Create frontend/src/lib/language/zh.json  (copy en.json, translate values)
// 2. Add:  import zh from './language/zh.json';
// 3. Add:  zh: zh as DeepPartial<Translation>,  to the LANGUAGE_REGISTRY below
// 4. Add the locale key to SUPPORTED_LOCALES
// That's it — the entire system (UI, location search, categories) picks it up.
// ─────────────────────────────────────────────────────────────────────────────

export type Translation = typeof en;

// DeepPartial: non-English translations may be incomplete
type DeepPartial<T> = { [K in keyof T]?: T[K] extends object ? DeepPartial<T[K]> : T[K] };

// Registry: add new languages here
const LANGUAGE_REGISTRY: Record<string, DeepPartial<Translation>> = {
  am: am as DeepPartial<Translation>,
  or: or_ as DeepPartial<Translation>,
  ti: ti as DeepPartial<Translation>,
};

// Supported locale keys — drives language switcher UI
export const SUPPORTED_LOCALES = ['en', 'am', 'or', 'ti'] as const;
export type Locale = typeof SUPPORTED_LOCALES[number] | string; // string = future locales

// Human-readable locale labels for the language switcher
export const LOCALE_LABELS: Record<string, string> = {
  en: 'English',
  am: 'አማርኛ',
  or: 'Afaan Oromoo',
  ti: 'ትግርኛ',
  // Add new entries here when adding a language, e.g.:
  // zh: '中文',
  // ar: 'العربية',
};

// Deep merge: fills missing keys from base (English)
function deepMerge<T extends object>(base: T, override: DeepPartial<T>): T {
  const result = { ...base };
  for (const [key, ov] of Object.entries(override)) {
    const k = key as keyof T;
    const bv = (base as any)[k];
    if (ov !== undefined && ov !== null) {
      if (typeof bv === 'object' && !Array.isArray(bv) && typeof ov === 'object') {
        (result as any)[k] = deepMerge(bv, ov);
      } else {
        (result as any)[k] = ov;
      }
    }
  }
  return result;
}

// Build merged translations — English is the authoritative base
export const translations: Record<string, Translation> = {
  en,
  ...Object.fromEntries(
    Object.entries(LANGUAGE_REGISTRY).map(([locale, partial]) => [
      locale,
      deepMerge(en, partial),
    ])
  ),
};

// Safe dot-path getter with English fallback
export function tKey(locale: string, path: string): string {
  const parts = path.split('.');
  const walk = (obj: any) => parts.reduce((acc, k) => acc?.[k], obj);
  const val = walk(translations[locale] ?? en);
  if (typeof val === 'string') return val;
  const enVal = walk(en);
  return typeof enVal === 'string' ? enVal : path;
}
