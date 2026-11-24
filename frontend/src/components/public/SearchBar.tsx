'use client';
import { useState, useRef, useEffect, useCallback, FormEvent } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useLocale } from '@/hooks/useLocale';
import { api, CategoryIntentData, CategoryChildChip } from '@/lib/api';
import {
  Search,
  X,
  ChevronDown,
  ChevronRight,
  ChevronLeft,
  Tag,
  Layers,
  Hash,
  MapPin,
  Sliders,
  FolderOpen,
  Folder,
  ArrowRight,
  ExternalLink,
  CornerDownRight,
  Sparkles,
} from 'lucide-react';
import { locPath } from '@/lib/location';

interface SearchEntry {
  id: number;
  type: 'category' | 'subcategory' | 'item';
  label: string;
  labelEn: string;
  catId: number;
  catSlug: string;
  subSlug?: string;
  itemSlug?: string;
  breadcrumb?: string;
  hierarchy?: string;
  level?: number;
  isLeaf?: boolean;
}

interface DrilldownCrumb {
  id: number;
  name: string;
  slug: string;
  level: number;
}

interface DrilldownChild {
  id: number;
  name: string;
  slug: string;
  level: number;
  is_leaf: boolean;
  icon?: string;
  category_system_code?: string;
}

interface RelatedCategory {
  id: number;
  name: string;
  slug: string;
  level: number;
  breadcrumbs?: any[];
}

const DRILLDOWN_I18N = {
  en: {
    back: 'Back',
    categories: 'Categories',
    searchAll: 'Search all listings in',
    subcategories: 'Subcategories',
    relatedCategories: 'Cross-Category Discoverability',
    loading: 'Loading categories...',
    noChildren: 'No subcategories found. Search listings in this category above.',
    leafCategory: 'Selectable Leaf Category',
    explore: 'Explore subcategories',
    searchThis: 'Search listings',
    directSearch: 'Search products for',
    rootHeading: 'Browse Categories',
    resultsHeading: 'Results for',
  },
  am: {
    back: 'ተመለስ',
    categories: 'ምድቦች',
    searchAll: 'በዚህ ውስጥ ያሉ ሁሉንም ዝርዝሮች ፈልግ፦',
    subcategories: 'ንዑስ ምድቦች',
    relatedCategories: 'ተዛማጅ ምድቦች (ተሻጋሪ ግኝት)',
    loading: 'ምድቦችን በመጫን ላይ...',
    noChildren: 'ምንም ንዑስ ምድቦች አልተገኙም። ከላይ በዚህ ምድብ ውስጥ ይፈልጉ።',
    leafCategory: 'ቀጥታ የዕቃ ምድብ',
    explore: 'ንዑስ ምድቦችን ያስሱ',
    searchThis: 'ዝርዝሮችን ፈልግ',
    directSearch: 'ዕቃዎችን ፈልግ፦',
    rootHeading: 'ምድቦችን ያስሱ',
    resultsHeading: 'ውጤቶች ለ',
  },
  om: {
    back: 'Deebi’i',
    categories: 'Gareewwan',
    searchAll: 'Tarree hunda keessatti barbaadi:',
    subcategories: 'Gareewwan xixiqqoo',
    relatedCategories: 'Gareewwan walqabatan',
    loading: 'Gareewwan fe’aa jira...',
    noChildren: 'Gareen xiqqaan hin argamne. Olitti barbaadi.',
    leafCategory: 'Garee Kallattii',
    explore: 'Gareewwan xixiqqoo qoradhu',
    searchThis: 'Tarree barbaadi',
    directSearch: 'Meeshaalee barbaadi:',
    rootHeading: 'Gareewwan qoradhu',
    resultsHeading: 'Bu’aa barbaachaa:',
  },
  ti: {
    back: 'ተመለስ',
    categories: 'ምድባት',
    searchAll: 'ኣብዚ ዘለዉ ኩሎም ዝርዝራት ድለይ፦',
    subcategories: 'ንኡሳን ምድባት',
    relatedCategories: 'ተዛመድቲ ምድባት',
    loading: 'ምድባት ይጽዕን ኣሎ...',
    noChildren: 'ዝኾነ ንኡስ ምድብ ኣይተረኽበን። ኣብ ላዕሊ ድለዩ።',
    leafCategory: 'ቀጥታ ምድብ',
    explore: 'ንኡሳን ምድባት ዳህስስ',
    searchThis: 'ዝርዝራት ድለይ',
    directSearch: 'ንብረታት ድለይ፦',
    rootHeading: 'ምድባት ዞር',
    resultsHeading: 'ውጽኢታት ን',
  },
};

const TYPE_ICON = {
  category: Tag,
  subcategory: Layers,
  item: Hash,
} as const;

const TYPE_COLOR = {
  category: 'text-[#1769E0]',
  subcategory: 'text-zinc-500',
  item: 'text-zinc-400',
} as const;

export function SearchBar() {
  const { t, locale } = useLocale();
  const u = (t as any)?.ui ?? {};
  const router = useRouter();
  const pathname = usePathname();

  const langKey = (locale === 'am' || locale === 'om' || locale === 'or' || locale === 'ti')
    ? (locale === 'or' ? 'om' : locale)
    : 'en';
  const i18n = DRILLDOWN_I18N[langKey as keyof typeof DRILLDOWN_I18N] || DRILLDOWN_I18N.en;

  const [q, setQ] = useState('');
  const [open, setOpen] = useState(false);
  const [highlighted, setHighlighted] = useState(-1);
  const [rootCategories, setRootCategories] = useState<SearchEntry[]>([]);
  const [liveResults, setLiveResults] = useState<SearchEntry[]>([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [locationSuggestions, setLocationSuggestions] = useState<any[]>([]);
  const [categoryIntent, setCategoryIntent] = useState<CategoryIntentData | null>(null);

  // ── Recursive In-Popover Drilldown State ─────────────────────────────────────
  const [drilldownStack, setDrilldownStack] = useState<DrilldownCrumb[]>([]);
  const [drilldownChildren, setDrilldownChildren] = useState<DrilldownChild[]>([]);
  const [drilldownRelated, setDrilldownRelated] = useState<RelatedCategory[]>([]);
  const [drilldownLoading, setDrilldownLoading] = useState(false);

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const searchDebounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const intentDebounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const locationDebounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  // Helper to extract localized category name
  const extractLocalizedName = useCallback(
    (item: any): string => {
      if (!item) return '';
      if (item.i18n && item.i18n[langKey]) return item.i18n[langKey];
      if (langKey === 'am' && item.name_am) return item.name_am;
      if (langKey === 'om' && item.name_om) return item.name_om;
      if (langKey === 'ti' && item.name_ti) return item.name_ti;
      return item.name || item.name_en || item.slug || '';
    },
    [langKey]
  );

  // Load root categories dynamically
  useEffect(() => {
    api.getRootCategories(locale)
      .then(roots => {
        if (!Array.isArray(roots)) return;
        const idx: SearchEntry[] = roots.map(root => {
          const rootLabel = extractLocalizedName(root);
          const rootEn = root.name_en || root.name || root.slug;
          return {
            id: root.id,
            type: 'category',
            label: rootLabel,
            labelEn: rootEn,
            catId: root.id,
            catSlug: root.slug,
            level: 0,
            isLeaf: root.is_leaf ?? false,
          };
        });
        setRootCategories(idx);
      })
      .catch(() => {});
  }, [locale, extractLocalizedName]);

  // Close on outside click
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, []);

  // Fetch children and related categories whenever drilldownStack changes
  useEffect(() => {
    if (drilldownStack.length === 0) {
      setDrilldownChildren([]);
      setDrilldownRelated([]);
      setDrilldownLoading(false);
      return;
    }

    const current = drilldownStack[drilldownStack.length - 1];
    setDrilldownLoading(true);

    Promise.all([
      api.getCategoryChildren(current.id, locale).catch(() => []),
      api.getCategoryRelated(current.id, locale).catch(() => []),
    ]).then(([children, related]) => {
      if (Array.isArray(children)) {
        setDrilldownChildren(
          children.map(c => ({
            id: c.id,
            name: extractLocalizedName(c),
            slug: c.slug,
            level: c.level,
            is_leaf: c.is_leaf,
            icon: c.icon,
          }))
        );
      } else {
        setDrilldownChildren([]);
      }

      if (Array.isArray(related)) {
        setDrilldownRelated(
          related.map(r => ({
            id: r.id,
            name: extractLocalizedName(r),
            slug: r.slug,
            level: r.level,
            breadcrumbs: r.breadcrumbs,
          }))
        );
      } else {
        setDrilldownRelated([]);
      }
      setDrilldownLoading(false);
      setHighlighted(-1);
    });
  }, [drilldownStack, locale, extractLocalizedName]);

  const buildLocationUrl = (r: any) => {
    const p = new URLSearchParams();
    if (r.region) p.set('region', r.region);
    if (r.zone) p.set('zone', r.zone);
    if (r.subcity) p.set('subcity', r.subcity);
    if (r.woreda) p.set('woreda', r.woreda);
    if (r.latitude && r.longitude) {
      p.set('lat', String(r.latitude));
      p.set('lng', String(r.longitude));
    }
    return `/products?${p.toString()}`;
  };

  const buildUrl = (e: SearchEntry) => {
    const p = new URLSearchParams();
    p.set('category_hierarchy_id', String(e.id));
    p.set('cat', e.catSlug);
    if (e.subSlug) p.set('subcategory', e.subSlug);
    if (e.itemSlug) p.set('item', e.itemSlug);
    return `/products?${p.toString()}`;
  };

  const selectEntry = (e: SearchEntry) => {
    setQ(e.label);
    setOpen(false);
    setDrilldownStack([]);
    router.push(buildUrl(e));
  };

  const startDrilldown = (cat: { id: number; name: string; slug: string; level: number }) => {
    setDrilldownStack(prev => [...prev, cat]);
    setHighlighted(-1);
  };

  const popDrilldown = () => {
    setDrilldownStack(prev => prev.slice(0, -1));
    setHighlighted(-1);
  };

  const jumpToDrilldownIndex = (idx: number) => {
    setDrilldownStack(prev => prev.slice(0, idx + 1));
    setHighlighted(-1);
  };

  const searchAllInCurrentCategory = (current: DrilldownCrumb) => {
    setOpen(false);
    setDrilldownStack([]);
    router.push(`/products?category_hierarchy_id=${current.id}&cat=${current.slug}`);
  };

  const selectChildCategory = (child: DrilldownChild) => {
    if (child.is_leaf) {
      setOpen(false);
      setDrilldownStack([]);
      router.push(`/products?category_hierarchy_id=${child.id}&cat=${child.slug}`);
    } else {
      startDrilldown({
        id: child.id,
        name: child.name,
        slug: child.slug,
        level: child.level,
      });
    }
  };

  const handleChange = (ev: React.ChangeEvent<HTMLInputElement>) => {
    const val = ev.target.value;
    setQ(val);
    setHighlighted(-1);
    setOpen(true);
    // Reset drilldown stack if typing free query
    if (drilldownStack.length > 0) {
      setDrilldownStack([]);
    }

    if (debounceRef.current) clearTimeout(debounceRef.current);
    if (searchDebounceRef.current) clearTimeout(searchDebounceRef.current);
    if (intentDebounceRef.current) clearTimeout(intentDebounceRef.current);
    if (locationDebounceRef.current) clearTimeout(locationDebounceRef.current);

    const clean = val.trim();
    if (clean) {
      setSearchLoading(true);

      // Category Intent Auto-Discovery
      if (clean.length >= 2) {
        intentDebounceRef.current = setTimeout(async () => {
          try {
            const intent = await api.resolveCategoryIntent(clean, locale);
            if (intent && intent.has_intent) {
              setCategoryIntent(intent);
            } else {
              setCategoryIntent(null);
            }
          } catch {
            setCategoryIntent(null);
          }
        }, 120);
      } else {
        setCategoryIntent(null);
      }

      searchDebounceRef.current = setTimeout(async () => {
        try {
          const results = await api.searchCategories(clean, locale, 20);
          if (Array.isArray(results)) {
            const mapped: SearchEntry[] = results.map(r => {
              const bcs = Array.isArray(r.breadcrumbs) ? r.breadcrumbs : [];
              const rootSlug = bcs[0]?.slug || r.slug;
              const subSlug = bcs[1]?.slug;
              const itemSlug = bcs[2]?.slug;
              const parentBcs = bcs.length > 1 ? bcs.slice(0, -1) : [];
              const parentPath = parentBcs.map((b: any) => b.name).join(' › ');
              return {
                id: r.id,
                type: r.level === 0 ? 'category' : r.level === 1 ? 'subcategory' : 'item',
                label: extractLocalizedName(r),
                labelEn: r.name_en || r.name || r.slug,
                catId: bcs[0]?.id || r.id,
                catSlug: rootSlug,
                subSlug: subSlug,
                itemSlug: itemSlug,
                breadcrumb: parentPath ? `in ${parentPath}` : '',
                hierarchy: parentPath,
                level: r.level,
                isLeaf: r.is_leaf,
              };
            });
            setLiveResults(mapped);
          }
        } catch {
          setLiveResults([]);
        } finally {
          setSearchLoading(false);
        }
      }, 150);

      locationDebounceRef.current = setTimeout(async () => {
        try {
          const results = await api.searchLocations(clean);
          setLocationSuggestions(results.slice(0, 3));
        } catch {
          setLocationSuggestions([]);
        }
      }, 300);

      if (pathname === '/products') {
        debounceRef.current = setTimeout(() => {
          router.replace(`/products?search=${encodeURIComponent(clean)}`);
        }, 500);
      }
    } else {
      setLiveResults([]);
      setSearchLoading(false);
      setLocationSuggestions([]);
      setCategoryIntent(null);
      if (pathname === '/products') {
        debounceRef.current = setTimeout(() => {
          router.replace('/products');
        }, 300);
      }
    }
  };

  const displayList = q.trim().length >= 1 ? liveResults : rootCategories;
  const isDrilldownActive = drilldownStack.length > 0;
  const currentCrumb = isDrilldownActive ? drilldownStack[drilldownStack.length - 1] : null;

  const handleSubmit = (ev: FormEvent<HTMLFormElement>) => {
    ev.preventDefault();
    if (debounceRef.current) clearTimeout(debounceRef.current);
    setOpen(false);

    if (isDrilldownActive && currentCrumb) {
      if (highlighted === 0) {
        searchAllInCurrentCategory(currentCrumb);
        return;
      } else if (highlighted > 0 && drilldownChildren[highlighted - 1]) {
        selectChildCategory(drilldownChildren[highlighted - 1]);
        return;
      }
    }

    if (highlighted >= 0 && displayList[highlighted]) {
      selectEntry(displayList[highlighted]);
    } else if (q.trim()) {
      router.push(`/products?search=${encodeURIComponent(q.trim())}`);
    } else {
      router.push('/products');
    }
  };

  const handleKeyDown = (ev: React.KeyboardEvent) => {
    if (!open) return;

    if (ev.key === 'ArrowDown') {
      ev.preventDefault();
      const maxLen = isDrilldownActive ? drilldownChildren.length + 1 : displayList.length;
      setHighlighted(h => Math.min(h + 1, maxLen - 1));
    } else if (ev.key === 'ArrowUp') {
      ev.preventDefault();
      setHighlighted(h => Math.max(h - 1, -1));
    } else if (ev.key === 'Escape') {
      setOpen(false);
    } else if (ev.key === 'Backspace' && q === '' && isDrilldownActive) {
      ev.preventDefault();
      popDrilldown();
    }
  };

  const clear = () => {
    setQ('');
    setOpen(false);
    setDrilldownStack([]);
    setLocationSuggestions([]);
    setCategoryIntent(null);
    if (debounceRef.current) clearTimeout(debounceRef.current);
    if (intentDebounceRef.current) clearTimeout(intentDebounceRef.current);
    router.push('/products');
    inputRef.current?.focus();
  };

  const highlight = (label: string, query: string) => {
    if (!query.trim()) return <span>{label}</span>;
    const idx = label.toLowerCase().indexOf(query.toLowerCase());
    if (idx === -1) return <span>{label}</span>;
    return (
      <>
        {label.slice(0, idx)}
        <mark className="bg-orange-200 dark:bg-orange-800 text-inherit rounded px-0.5">
          {label.slice(idx, idx + query.length)}
        </mark>
        {label.slice(idx + query.length)}
      </>
    );
  };

  return (
    <div ref={containerRef} className="relative w-full">
      <form onSubmit={handleSubmit} className="relative flex w-full">
        <Search
          className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none"
          size={15}
        />
        <input
          ref={inputRef}
          type="text"
          value={q}
          onChange={handleChange}
          onFocus={() => setOpen(true)}
          onKeyDown={handleKeyDown}
          placeholder={t.nav?.search || 'Search products, vehicles, electronics...'}
          autoComplete="off"
          className="w-full pl-9 pr-9 py-2 rounded-full border border-zinc-200 dark:border-zinc-800 bg-zinc-100 dark:bg-zinc-900 text-base sm:text-sm outline-none focus:ring-2 focus:ring-[#1769E0] transition-all"
        />
        <div className="absolute right-2.5 top-1/2 -translate-y-1/2 flex items-center gap-1.5">
          {q && (
            <button
              type="button"
              onClick={clear}
              className="text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200 transition-colors p-1.5 tap-bounce active:scale-90"
              aria-label="Clear search"
            >
              <X size={14} />
            </button>
          )}
        </div>
      </form>

      {/* Popover Dropdown */}
      {open && (
        <div className="absolute top-full left-0 right-0 mt-2 rounded-2xl border shadow-2xl z-[200] bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 max-h-96 overflow-y-auto">
          {/* VIEW A: RECURSIVE DRILLDOWN VIEW */}
          {isDrilldownActive && currentCrumb ? (
            <div className="flex flex-col">
              {/* Drilldown Navigation Breadcrumbs */}
              <div className="flex items-center gap-1.5 px-3 py-2 bg-zinc-50 dark:bg-zinc-800/80 border-b border-zinc-200 dark:border-zinc-800 text-xs">
                <button
                  type="button"
                  onMouseDown={popDrilldown}
                  className="flex items-center gap-1 text-[#1769E0] font-bold hover:underline px-1.5 py-0.5 rounded hover:bg-[#EAF3FF] dark:hover:bg-[#1769E0]/20 transition-colors"
                >
                  <ChevronLeft size={14} />
                  <span>{i18n.back}</span>
                </button>
                <span className="text-zinc-300 dark:text-zinc-700">|</span>
                <div className="flex items-center gap-1 overflow-x-auto text-[11px] font-semibold text-zinc-600 dark:text-zinc-300">
                  <button
                    type="button"
                    onMouseDown={() => setDrilldownStack([])}
                    className="hover:text-[#1769E0] transition-colors whitespace-nowrap"
                  >
                    {i18n.categories}
                  </button>
                  {drilldownStack.map((crumb, idx) => (
                    <div key={crumb.id} className="flex items-center gap-1 whitespace-nowrap">
                      <ChevronRight size={10} className="text-zinc-400" />
                      {idx === drilldownStack.length - 1 ? (
                        <span className="font-bold text-[#1769E0]">{crumb.name}</span>
                      ) : (
                        <button
                          type="button"
                          onMouseDown={() => jumpToDrilldownIndex(idx)}
                          className="hover:text-[#1769E0] transition-colors"
                        >
                          {crumb.name}
                        </button>
                      )}
                    </div>
                  ))}
                </div>
              </div>

              {/* Action Banner: Search All Listings in Current Category */}
              <button
                type="button"
                onMouseDown={() => searchAllInCurrentCategory(currentCrumb)}
                className={`w-full flex items-center justify-between px-3 py-2.5 text-left transition-colors border-b border-zinc-100 dark:border-zinc-800 ${
                  highlighted === 0
                    ? 'bg-[#EAF3FF] dark:bg-[#1769E0]/20'
                    : 'bg-blue-50/50 dark:bg-blue-950/20 hover:bg-[#EAF3FF] dark:hover:bg-[#1769E0]/15'
                }`}
              >
                <div className="flex items-center gap-2.5 min-w-0">
                  <div className="w-7 h-7 rounded-full bg-[#1769E0] text-white flex items-center justify-center shrink-0 shadow-sm">
                    <Search size={13} />
                  </div>
                  <div className="min-w-0">
                    <p className="text-xs text-zinc-500 dark:text-zinc-400 font-medium">
                      {i18n.searchAll}
                    </p>
                    <p className="text-sm font-extrabold text-[#1769E0] truncate">
                      &quot;{currentCrumb.name}&quot;
                    </p>
                  </div>
                </div>
                <ArrowRight size={14} className="text-[#1769E0] shrink-0 ml-2" />
              </button>

              {/* Child Subcategories Section */}
              <div className="px-3 pt-2.5 pb-1 flex items-center justify-between">
                <p className="text-[10px] font-black uppercase tracking-widest text-zinc-400">
                  {i18n.subcategories} ({drilldownChildren.length})
                </p>
                <span className="text-[10px] text-zinc-400 font-semibold">
                  Level {currentCrumb.level + 1}
                </span>
              </div>

              {drilldownLoading ? (
                <div className="py-8 text-center">
                  <div className="inline-block animate-spin w-5 h-5 border-2 border-[#1769E0] border-t-transparent rounded-full mb-2" />
                  <p className="text-xs text-zinc-400 font-medium">{i18n.loading}</p>
                </div>
              ) : drilldownChildren.length === 0 ? (
                <div className="px-4 py-6 text-center">
                  <p className="text-xs text-zinc-400">{i18n.noChildren}</p>
                  <button
                    type="button"
                    onMouseDown={() => searchAllInCurrentCategory(currentCrumb)}
                    className="mt-3 text-xs font-bold text-[#1769E0] hover:underline"
                  >
                    {i18n.searchThis} →
                  </button>
                </div>
              ) : (
                <div className="divide-y divide-zinc-50 dark:divide-zinc-800/50">
                  {drilldownChildren.map((child, i) => {
                    const isActive = highlighted === i + 1;
                    return (
                      <div
                        key={child.id}
                        className={`flex items-center justify-between px-3 py-2 text-left transition-colors ${
                          isActive
                            ? 'bg-[#EAF3FF] dark:bg-[#1769E0]/15'
                            : 'hover:bg-zinc-50 dark:hover:bg-zinc-800/60'
                        }`}
                      >
                        <button
                          type="button"
                          onMouseDown={() => selectChildCategory(child)}
                          className="flex-1 flex items-center gap-2.5 min-w-0 pr-2 text-left"
                        >
                          <div
                            className={`shrink-0 ${
                              child.is_leaf ? 'text-emerald-500' : 'text-[#1769E0]'
                            }`}
                          >
                            {child.is_leaf ? <Tag size={13} /> : <FolderOpen size={14} />}
                          </div>
                          <div className="min-w-0 flex-1">
                            <p className="text-sm font-semibold truncate text-zinc-800 dark:text-zinc-100">
                              {child.name}
                            </p>
                            <span className="text-[9px] font-bold text-zinc-400">
                              {child.is_leaf
                                ? i18n.leafCategory
                                : `Level ${child.level + 1} • Click to explore deeper`}
                            </span>
                          </div>
                        </button>

                        <div className="flex items-center gap-1 shrink-0">
                          {/* Direct search shortcut button */}
                          <button
                            type="button"
                            title={i18n.searchThis}
                            onMouseDown={e => {
                              e.stopPropagation();
                              setOpen(false);
                              setDrilldownStack([]);
                              router.push(
                                `/products?category_hierarchy_id=${child.id}&cat=${child.slug}`
                              );
                            }}
                            className="p-1 text-zinc-400 hover:text-[#1769E0] hover:bg-zinc-200/50 dark:hover:bg-zinc-700/50 rounded transition-colors"
                          >
                            <Search size={13} />
                          </button>

                          {/* Drill-down chevron for non-leaf categories */}
                          {!child.is_leaf && (
                            <button
                              type="button"
                              title={i18n.explore}
                              onMouseDown={() => selectChildCategory(child)}
                              className="p-1 text-zinc-400 hover:text-[#1769E0] rounded transition-colors"
                            >
                              <ChevronRight size={15} />
                            </button>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}

              {/* Cross-Category Discoverability (Related Categories) */}
              {drilldownRelated.length > 0 && (
                <div className="mt-2 border-t border-zinc-100 dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-950/30 p-3">
                  <div className="flex items-center gap-1.5 mb-2">
                    <Sparkles size={12} className="text-amber-500" />
                    <p className="text-[10px] font-black uppercase tracking-wider text-zinc-500 dark:text-zinc-400">
                      {i18n.relatedCategories}
                    </p>
                  </div>
                  <div className="flex flex-wrap gap-1.5">
                    {drilldownRelated.map(rel => (
                      <button
                        key={rel.id}
                        type="button"
                        onMouseDown={() => {
                          startDrilldown({
                            id: rel.id,
                            name: rel.name,
                            slug: rel.slug,
                            level: rel.level,
                          });
                        }}
                        className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 text-zinc-700 dark:text-zinc-200 hover:border-[#1769E0] hover:text-[#1769E0] transition-colors"
                      >
                        <Tag size={10} className="text-[#1769E0]" />
                        <span>{rel.name}</span>
                        <CornerDownRight size={10} className="opacity-40" />
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>
          ) : (
            /* VIEW B: SEARCH RESULTS & INITIAL BROWSING VIEW */
            <div>
              {/* Category Intent Auto-Discovery & 1-Tap Subcategories */}
              {categoryIntent && categoryIntent.has_intent && categoryIntent.category && (
                <div className="p-3 bg-gradient-to-r from-blue-50 to-indigo-50/80 dark:from-blue-950/40 dark:to-indigo-950/40 border-b border-blue-200/80 dark:border-blue-900/60">
                  <div className="flex items-center justify-between gap-2">
                    <div className="flex items-center gap-2 min-w-0">
                      <div className="w-7 h-7 rounded-lg bg-blue-600 text-white flex items-center justify-center shrink-0 shadow-sm">
                        <FolderOpen size={14} />
                      </div>
                      <div className="min-w-0">
                        <p className="text-[10px] font-semibold uppercase tracking-wider text-blue-700 dark:text-blue-300">
                          {i18n.categories}
                        </p>
                        <p className="text-xs font-bold text-zinc-900 dark:text-zinc-100 truncate">
                          {categoryIntent.breadcrumbs && categoryIntent.breadcrumbs.length > 0
                            ? categoryIntent.breadcrumbs.map(b => b.name).join(' › ')
                            : (categoryIntent.category.name_en || categoryIntent.category.slug)}
                        </p>
                      </div>
                    </div>
                    <button
                      type="button"
                      onMouseDown={() => {
                        setOpen(false);
                        router.push(`/products?category_hierarchy_id=${categoryIntent.category!.id}&cat=${categoryIntent.category!.slug}`);
                      }}
                      className="px-2.5 py-1 text-xs font-bold text-white bg-blue-600 hover:bg-blue-700 rounded-lg shrink-0 flex items-center gap-1 shadow-sm transition-colors"
                    >
                      <span>Explore</span>
                      <ArrowRight size={12} />
                    </button>
                  </div>

                  {/* 1-Tap Subcategories inside this category */}
                  {categoryIntent.children && categoryIntent.children.length > 0 && (
                    <div className="mt-2 pt-2 border-t border-blue-200/60 dark:border-blue-900/40">
                      <p className="text-[10px] font-medium text-zinc-500 dark:text-zinc-400 mb-1.5 flex items-center gap-1">
                        <Layers size={11} className="text-blue-500" />
                        <span>Inside {categoryIntent.category.name_en}:</span>
                      </p>
                      <div className="flex flex-wrap gap-1.5">
                        {categoryIntent.children.slice(0, 8).map(ch => (
                          <button
                            key={ch.id}
                            type="button"
                            onMouseDown={() => {
                              setOpen(false);
                              router.push(`/products?category_hierarchy_id=${ch.id}&cat=${ch.slug}`);
                            }}
                            className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-white dark:bg-zinc-800 text-zinc-800 dark:text-zinc-200 hover:text-blue-600 dark:hover:text-blue-400 hover:border-blue-400 border border-zinc-200 dark:border-zinc-700 shadow-2xs transition-all"
                          >
                            <span>{ch.name}</span>
                            {ch.product_count !== undefined && ch.product_count > 0 && (
                              <span className="text-[10px] text-zinc-400">({ch.product_count})</span>
                            )}
                          </button>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Multi-Branch Ambiguity Notice (e.g., 'drill' in tools vs mining) */}
                  {categoryIntent.is_ambiguous && categoryIntent.alternatives && categoryIntent.alternatives.length > 0 && (
                    <div className="mt-2 pt-2 border-t border-blue-200/60 dark:border-blue-900/40">
                      <p className="text-[10px] font-semibold text-amber-700 dark:text-amber-400 mb-1 flex items-center gap-1">
                        <Sparkles size={11} className="text-amber-500" />
                        <span>Also found in other branches:</span>
                      </p>
                      <div className="flex flex-wrap gap-1.5">
                        {categoryIntent.alternatives.slice(0, 3).map(alt => (
                          <button
                            key={alt.id}
                            type="button"
                            onMouseDown={() => {
                              setOpen(false);
                              router.push(`/products?category_hierarchy_id=${alt.id}&cat=${alt.slug}`);
                            }}
                            className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-medium bg-amber-50 dark:bg-amber-950/40 text-amber-800 dark:text-amber-300 border border-amber-200/80 dark:border-amber-800/60 hover:bg-amber-100 transition-colors"
                          >
                            <span>{alt.name}</span>
                            {alt.hierarchy && (
                              <span className="text-[10px] text-amber-600/80 dark:text-amber-400/80">({alt.hierarchy})</span>
                            )}
                          </button>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Quick Filters (Buy / Rent) */}
                  <div className="mt-2 pt-1.5 border-t border-blue-200/40 dark:border-blue-900/30 flex items-center gap-2">
                    <span className="text-[10px] text-zinc-400 font-medium">Quick Filter:</span>
                    <button
                      type="button"
                      onMouseDown={() => {
                        setOpen(false);
                        const cat = categoryIntent.category!;
                        router.push(`/products?category_hierarchy_id=${cat.id}&cat=${cat.slug}&type=sale`);
                      }}
                      className="px-2 py-0.5 text-[11px] font-semibold rounded bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 hover:text-blue-600 text-zinc-700 dark:text-zinc-300 transition-colors"
                    >
                      🛍️ Buy
                    </button>
                    <button
                      type="button"
                      onMouseDown={() => {
                        setOpen(false);
                        const cat = categoryIntent.category!;
                        router.push(`/products?category_hierarchy_id=${cat.id}&cat=${cat.slug}&type=rent`);
                      }}
                      className="px-2 py-0.5 text-[11px] font-semibold rounded bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 hover:text-blue-600 text-zinc-700 dark:text-zinc-300 transition-colors"
                    >
                      🔑 Rent
                    </button>
                  </div>
                </div>
              )}

              {/* Location Suggestions (When query matches a location) */}
              {(() => {
                const seen = new Set<string>();
                const uniqueLocs = locationSuggestions.filter(r => {
                  const p = locPath(r, locale);
                  if (!p || seen.has(p)) return false;
                  seen.add(p);
                  return true;
                });
                if (uniqueLocs.length === 0) return null;

                const locBadge =
                  locale === 'am'
                    ? 'አካባቢ'
                    : locale === 'ti'
                    ? 'ቦታ'
                    : locale === 'or'
                    ? 'bakka'
                    : 'location';

                return (
                  <div className="border-b border-zinc-100 dark:border-zinc-800">
                    <div className="px-3 pt-2 pb-1">
                      <p className="text-[10px] font-black uppercase tracking-widest opacity-30">
                        {u.search_locations_heading || 'Locations'}
                      </p>
                    </div>
                    {uniqueLocs.map((r, i) => {
                      const displayPath = locPath(r, locale);
                      return (
                        <button
                          key={i}
                          onMouseDown={() => {
                            setOpen(false);
                            router.push(buildLocationUrl(r));
                          }}
                          className="w-full flex items-center gap-3 px-3 py-2.5 text-left hover:bg-[#EAF3FF] dark:hover:bg-[#1769E0]/15 transition-colors"
                        >
                          <MapPin size={14} className="text-[#1769E0] shrink-0" />
                          <p className="text-sm font-bold truncate">{displayPath}</p>
                          <span className="text-[9px] font-black uppercase tracking-wider px-1.5 py-0.5 rounded-full shrink-0 bg-[#EAF3FF] dark:bg-[#1769E0]/20 text-[#1769E0]">
                            {locBadge}
                          </span>
                        </button>
                      );
                    })}
                  </div>
                );
              })()}

              {/* Header Label */}
              <div className="px-3 pt-2.5 pb-1 flex items-center justify-between">
                <p className="text-[10px] font-black uppercase tracking-widest opacity-40">
                  {q.trim()
                    ? `${i18n.resultsHeading} "${q}"`
                    : i18n.rootHeading}
                </p>
                <span className="text-[10px] font-semibold text-zinc-400">
                  {displayList.length} items
                </span>
              </div>

              {/* No match view */}
              {displayList.length === 0 && q.trim().length >= 1 && (
                <div className="px-4 py-6 text-center">
                  <p className="text-sm font-bold opacity-40">
                    {u.no_match_for || 'No match for'} &quot;{q}&quot;
                  </p>
                  <p className="text-xs opacity-30 mt-1">
                    {u.try_different_word || 'Try a different word or browse categories below'}
                  </p>
                  <button
                    onClick={() => {
                      setQ('');
                      setOpen(true);
                    }}
                    className="mt-3 text-xs font-bold text-[#1769E0] hover:underline"
                  >
                    {u.browse_all_categories_arrow || 'Browse all categories →'}
                  </button>
                </div>
              )}

              {/* List of Matched or Root Categories */}
              {displayList.map((entry, i) => {
                const Icon = TYPE_ICON[entry.type];
                const isActive = i === highlighted;

                return (
                  <div
                    key={
                      entry.id
                        ? `search-cat-${entry.id}-${i}`
                        : `search-entry-${entry.catSlug}-${i}`
                    }
                    className={`flex items-center justify-between px-3 py-2 text-left transition-colors ${
                      isActive
                        ? 'bg-[#EAF3FF] dark:bg-[#1769E0]/15'
                        : 'hover:bg-zinc-50 dark:hover:bg-zinc-800/60'
                    }`}
                  >
                    {/* Main Category Click: selects or opens */}
                    <button
                      type="button"
                      onMouseDown={() => {
                        if (entry.isLeaf) {
                          selectEntry(entry);
                        } else {
                          startDrilldown({
                            id: entry.id,
                            name: entry.label,
                            slug: entry.catSlug,
                            level: entry.level ?? 0,
                          });
                        }
                      }}
                      onMouseEnter={() => setHighlighted(i)}
                      className="flex-1 flex items-center gap-3 min-w-0 pr-2 text-left"
                    >
                      <div className={`shrink-0 ${TYPE_COLOR[entry.type]}`}>
                        <Icon size={14} />
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-bold truncate">
                          {highlight(entry.label, q)}
                        </p>
                        {entry.breadcrumb && (
                          <p className="text-[10px] opacity-40 truncate">{entry.breadcrumb}</p>
                        )}
                      </div>
                    </button>

                    {/* Drill-down / Action icons */}
                    <div className="flex items-center gap-1.5 shrink-0">
                      <span
                        className={`text-[9px] font-black uppercase tracking-wider px-2 py-0.5 rounded-full ${
                          entry.level === 0 || entry.type === 'category'
                            ? 'bg-[#EAF3FF] dark:bg-[#1769E0]/20 text-[#1769E0]'
                            : entry.level === 1
                            ? 'bg-blue-100 dark:bg-blue-900/40 text-blue-600'
                            : entry.level === 2
                            ? 'bg-emerald-100 dark:bg-emerald-900/40 text-emerald-600'
                            : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-500'
                        }`}
                      >
                        {entry.level === 0
                          ? 'Sector'
                          : entry.level === 1
                          ? 'Superfamily'
                          : entry.level === 2
                          ? 'Category'
                          : entry.level === 3
                          ? 'Subcategory'
                          : 'Specific'}
                      </span>

                      {/* Interactive Drilldown Icon */}
                      {!entry.isLeaf && (
                        <button
                          type="button"
                          title={i18n.explore}
                          onMouseDown={e => {
                            e.stopPropagation();
                            startDrilldown({
                              id: entry.id,
                              name: entry.label,
                              slug: entry.catSlug,
                              level: entry.level ?? 0,
                            });
                          }}
                          className="p-1 text-zinc-400 hover:text-[#1769E0] rounded transition-colors"
                        >
                          <ChevronRight size={15} />
                        </button>
                      )}
                    </div>
                  </div>
                );
              })}

              {/* Free-text search option when typing */}
              {q.trim().length >= 1 && (
                <div className="border-t border-zinc-100 dark:border-zinc-800">
                  <button
                    type="button"
                    onMouseDown={() => {
                      setOpen(false);
                      router.push(`/products?search=${encodeURIComponent(q.trim())}`);
                    }}
                    className="w-full flex items-center gap-3 px-3 py-2.5 hover:bg-zinc-50 dark:hover:bg-zinc-800/60 transition-colors"
                  >
                    <Search size={14} className="text-[#1769E0] shrink-0" />
                    <p className="text-sm font-bold">
                      {i18n.directSearch} &quot;
                      <span className="text-[#1769E0]">{q}</span>&quot;
                    </p>
                  </button>
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
