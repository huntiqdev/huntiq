export const SUPPORTED_TAXONOMY_LOCALES = ['en', 'am', 'or', 'ti'] as const;
export type TaxonomyLocale = typeof SUPPORTED_TAXONOMY_LOCALES[number];

export interface LocalizedCategoryEntity {
  id?: number;
  name?: string;
  name_en?: string;
  name_am?: string;
  name_om?: string;
  name_ti?: string;
  slug?: string;
  i18n?: Record<string, string>;
}

/**
 * Resolves a category's localized display title dynamically from its database record.
 * Authoritative source: PostgreSQL category_hierarchy via backend API.
 * Precedence: Dynamic i18n[locale] -> requested locale field -> name -> name_en -> formatted slug.
 */
export function resolveCategoryTitle(cat?: LocalizedCategoryEntity | null, locale: string = 'en'): string {
  if (!cat) return '';
  if (cat.i18n && cat.i18n[locale]) {
    return cat.i18n[locale];
  }
  switch (locale) {
    case 'am':
      if (cat.name_am) return cat.name_am;
      break;
    case 'or':
    case 'om':
      if (cat.name_om) return cat.name_om;
      break;
    case 'ti':
      if (cat.name_ti) return cat.name_ti;
      break;
    case 'en':
      if (cat.name_en) return cat.name_en;
      break;
  }
  return cat.i18n?.['en'] || cat.name || cat.name_en || cat.slug?.replace(/_/g, ' ') || '';
}

/**
 * Compatibility helper for components resolving category labels.
 * When passed a Category entity, extracts the dynamic localized title.
 */
export function getCategoryLabel(entityOrId: LocalizedCategoryEntity | number | string | undefined | null, locale: string = 'en'): string {
  if (!entityOrId) return '';
  if (typeof entityOrId === 'object') {
    return resolveCategoryTitle(entityOrId, locale);
  }
  return '';
}

/**
 * Returns true if a valid localized name exists on the category record.
 */
export function hasCategoryTranslation(cat?: LocalizedCategoryEntity | null, locale: string = 'en'): boolean {
  if (!cat) return false;
  switch (locale) {
    case 'am':
      return Boolean(cat.name_am);
    case 'or':
    case 'om':
      return Boolean(cat.name_om);
    case 'ti':
      return Boolean(cat.name_ti);
    default:
      return Boolean(cat.name_en || cat.name);
  }
}
