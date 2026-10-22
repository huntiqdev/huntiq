// I18nLocation: en is required, any other locale key is optional.
// This is open — adding a new language (e.g. zh, ar) just works automatically.
export interface I18nLocation {
  en: string;
  [locale: string]: string | undefined;
}

/**
 * Returns the best display name for a location item given the current locale.
 * Falls back to English if the locale translation is missing.
 * Extensible: works for any locale key (am, or, ti, zh, ar, fr, etc.)
 */
export function locName(item: I18nLocation, locale: string): string {
  if (locale && locale !== 'en' && item[locale]) return item[locale] as string;
  return item.en;
}

/**
 * Build a human-readable location path string from a search result row.
 * Picks the right language for each level.
 */
export function locPath(row: Record<string, any>, locale: string): string {
  const pick = (en: string | null, suffix: string): string => {
    if (!en) return '';
    const locKey = `${suffix}_${locale}`;
    return (locale !== 'en' && row[locKey]) ? row[locKey] : en;
  };
  return [
    pick(row.region, 'region'),
    pick(row.zone, 'zone'),
    pick(row.subcity, 'subcity'),
    pick(row.woreda, 'woreda'),
  ].filter(Boolean).join(' › ');
}
