import en from './translations/en.json';
import am from './translations/am.json';
import or_ from './translations/or.json';
import ti from './translations/ti.json';

export const SUPPORTED_TAXONOMY_LOCALES = ['en', 'am', 'or', 'ti'] as const;
export type TaxonomyLocale = typeof SUPPORTED_TAXONOMY_LOCALES[number];

/**
 * Authoritative category translation dictionaries indexed strictly by canonical PostgreSQL category ID.
 * Exactly 2,575 categories covered per language with zero hierarchy duplication.
 */
export const categoryTranslations: Record<string, Record<string, string>> = {
  en: en as Record<string, string>,
  am: am as Record<string, string>,
  or: or_ as Record<string, string>,
  om: or_ as Record<string, string>, // Alias for Afaan Oromoo locale code
  ti: ti as Record<string, string>,
};

/**
 * Resolves a category's localized display title by its canonical integer ID.
 * Fallback precedence: Requested locale -> English base -> Empty string.
 */
export function getCategoryLabel(id: number | string | undefined | null, locale: string = 'en'): string {
  if (id === undefined || id === null) return '';
  const key = String(id);
  const langDict = categoryTranslations[locale] || categoryTranslations['en'];
  if (langDict && langDict[key]) {
    return langDict[key];
  }
  // Fallback to English
  const enDict = categoryTranslations['en'];
  if (enDict && enDict[key]) {
    return enDict[key];
  }
  return '';
}

/**
 * Returns true if a translation exists for the specified category ID in the given language.
 */
export function hasCategoryTranslation(id: number | string, locale: string = 'en'): boolean {
  const key = String(id);
  const dict = categoryTranslations[locale];
  return Boolean(dict && dict[key]);
}
