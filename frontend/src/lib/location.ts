import type { I18nLocation } from '@/types';

/**
 * Returns the best display name for a location item given the current locale.
 * Uses the localized string from backend if available, falling back to English.
 */
export function locName(item: I18nLocation | null | undefined, locale: string): string {
  if (!item) return '';
  const norm = (locale === 'om' ? 'or' : locale) || 'en';
  if (norm !== 'en' && item[norm]) {
    return item[norm] as string;
  }
  return item.en || '';
}

/**
 * Extract deduplicated and localized location parts from a backend database row.
 * Eliminates duplicate hierarchy labels (e.g., Addis Ababa › Addis Ababa › Bole)
 * and picks the active locale columns (e.g. region_am, zone_or, subcity_ti).
 */
export function formatLocationParts(row: Record<string, any> | null | undefined, locale: string): string[] {
  if (!row) return [];
  const normLocale = (locale === 'om' ? 'or' : locale) || 'en';

  const pick = (en: string | null | undefined, suffix: string): string => {
    if (!en) return '';
    const trimmed = typeof en === 'string' ? en.trim() : '';
    if (!trimmed) return '';

    const locKey = `${suffix}_${normLocale}`;
    if (normLocale !== 'en' && row[locKey]) {
      return typeof row[locKey] === 'string' ? row[locKey].trim() : '';
    }

    if (normLocale !== 'en' && (suffix === 'woreda' || /^woreda\s+/i.test(trimmed))) {
      if (normLocale === 'am' || normLocale === 'ti') {
        return trimmed.replace(/^woreda\s+/i, 'ወረዳ ');
      }
    }

    return trimmed;
  };

  const r = pick(row.region, 'region');
  const z = pick(row.zone, 'zone');
  const s = pick(row.subcity, 'subcity');
  const w = pick(row.woreda, 'woreda');
  const k = pick(row.kebele, 'kebele');

  const parts: string[] = [];
  const seen = new Set<string>();

  const add = (text: string) => {
    if (!text) return;
    const key = text.trim().toLowerCase();
    if (!seen.has(key)) {
      seen.add(key);
      parts.push(text.trim());
    }
  };

  add(r);
  add(z);
  add(s);
  add(w);
  add(k);

  return parts;
}

/**
 * Build a human-readable location path string from a backend search result row.
 */
export function locPath(row: Record<string, any> | null | undefined, locale: string): string {
  const parts = formatLocationParts(row, locale);
  return parts.join(' › ');
}
