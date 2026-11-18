/**
 * Dynamic category taxonomy utilities
 * All categories are fetched from database - NO hardcoded data
 */

interface CategoryLike {
  name?: string;
  name_en?: string;
  name_am?: string;
  name_om?: string;
  name_ti?: string;
  slug?: string;
  i18n?: Record<string, string>;
}

/**
 * Resolves the category title in the requested language with proper fallback
 * @param category - Category object from database
 * @param locale - Language code (en, am, om/or, ti)
 * @returns Localized category name or slug as fallback
 */
export function resolveCategoryTitle(
  category: CategoryLike | null | undefined,
  locale: string = 'en'
): string {
  if (!category) return '';

  // Check i18n field first if available
  if (category.i18n && category.i18n[locale]) {
    return category.i18n[locale];
  }

  // Fallback to specific language fields based on locale
  switch (locale) {
    case 'am':
      return category.name_am || category.name || category.name_en || category.slug || '';
    case 'om':
    case 'or': // Oromo can be 'om' or 'or'
      return category.name_om || category.name || category.name_en || category.slug || '';
    case 'ti':
      return category.name_ti || category.name || category.name_en || category.slug || '';
    case 'en':
    default:
      return category.name_en || category.name || category.slug || '';
  }
}

/**
 * Gets the category name with proper formatting
 * @param category - Category object from database
 * @param locale - Language code
 * @returns Formatted category name
 */
export function getCategoryName(
  category: CategoryLike | null | undefined,
  locale: string = 'en'
): string {
  const title = resolveCategoryTitle(category, locale);
  // Clean up slug formatting if that's what we got
  return title.replace(/_/g, ' ').replace(/-/g, ' ');
}
