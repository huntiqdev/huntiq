import type { Locale } from '@/lib/i18n';

// A category as returned by the API (with 4-language names)
export interface CategoryI18n {
  id: number;
  name: string;
  slug: string;
  name_en?: string;
  name_am?: string;
  name_or?: string;
  name_ti?: string;
  sort_order?: number;
  is_active?: boolean;
}

// Pick the best display name for a category given the current locale
export function categoryName(cat: CategoryI18n, locale: Locale = 'en'): string {
  const map: Record<Locale, keyof CategoryI18n> = {
    en: 'name_en',
    am: 'name_am',
    or: 'name_or',
    ti: 'name_ti',
  };
  const key = map[locale];
  return (cat[key] as string | undefined) || cat.name_en || cat.name || cat.slug;
}

// Legacy: kept for backward compat with validators.ts — returns slug list from fetched categories
export function getCategorySlugs(categories: CategoryI18n[]): string[] {
  return categories.map(c => c.slug);
}
