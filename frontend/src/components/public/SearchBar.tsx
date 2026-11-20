'use client';
import { useState, useRef, useEffect, useCallback, FormEvent } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useLocale } from '@/hooks/useLocale';
import { api } from '@/lib/api';
import { Search, X, ChevronDown, Tag, Layers, Hash, MapPin, Sliders } from 'lucide-react';
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
}

// Simple fuzzy: every char of needle must appear in order in haystack
function fuzzyMatch(haystack: string, needle: string): boolean {
  const h = haystack.toLowerCase();
  const n = needle.toLowerCase();
  let hi = 0;
  for (let ni = 0; ni < n.length; ni++) {
    const found = h.indexOf(n[ni], hi);
    if (found === -1) return false;
    hi = found + 1;
  }
  return true;
}

// Score: exact prefix > contains > fuzzy
function score(label: string, q: string): number {
  const l = label.toLowerCase();
  const n = q.toLowerCase();
  if (l === n) return 100;
  if (l.startsWith(n)) return 80;
  if (l.includes(n)) return 60;
  if (fuzzyMatch(l, n)) return 30;
  return 0;
}

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

// ── Component ────────────────────────────────────────────────────────────────
export function SearchBar() {
  const { t, locale } = useLocale();
  const u = (t as any)?.ui ?? {};
  const router = useRouter();
  const pathname = usePathname();
  const [q, setQ] = useState('');
  const [open, setOpen] = useState(false);
  const [highlighted, setHighlighted] = useState(-1);
  const [rootCategories, setRootCategories] = useState<SearchEntry[]>([]);
  const [liveResults, setLiveResults] = useState<SearchEntry[]>([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [locationSuggestions, setLocationSuggestions] = useState<any[]>([]);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const searchDebounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  // Load popular root categories dynamically from authoritative database taxonomy
  useEffect(() => {
    api.getRootCategories(locale).then(roots => {
      if (!Array.isArray(roots)) return;
      const idx: SearchEntry[] = roots.map(root => {
        const rootLabel =
          (locale === 'am' && root.name_am) ||
          ((locale === 'or' || locale === 'om') && root.name_om) ||
          (locale === 'ti' && root.name_ti) ||
          (locale === 'en' && root.name_en) ||
          root.name ||
          root.name_en ||
          root.slug;
        const rootEn = root.name_en || root.name || root.slug;
        return {
          id: root.id,
          type: 'category',
          label: rootLabel,
          labelEn: rootEn,
          catId: root.id,
          catSlug: root.slug,
          level: 0,
        };
      });
      setRootCategories(idx);
    }).catch(() => { });
  }, [locale]);

  const getLabel = useCallback((e: SearchEntry): string => e.label, []);
  const getBreadcrumb = useCallback((e: SearchEntry): string => e.breadcrumb || e.hierarchy || '', []);

  const displayList = q.trim().length >= 1 ? liveResults : rootCategories;

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

  const buildLocationUrl = (r: any) => {
    const p = new URLSearchParams();
    if (r.region) p.set('region', r.region);
    if (r.zone) p.set('zone', r.zone);
    if (r.subcity) p.set('subcity', r.subcity);
    if (r.woreda) p.set('woreda', r.woreda);
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
    setQ(getLabel(e));
    setOpen(false);
    router.push(buildUrl(e));
  };

  const locationDebounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const handleChange = (ev: React.ChangeEvent<HTMLInputElement>) => {
    const val = ev.target.value;
    setQ(val);
    setHighlighted(-1);
    setOpen(true);
    if (debounceRef.current) clearTimeout(debounceRef.current);
    if (searchDebounceRef.current) clearTimeout(searchDebounceRef.current);
    if (locationDebounceRef.current) clearTimeout(locationDebounceRef.current);

    const clean = val.trim();
    if (clean) {
      setSearchLoading(true);
      // Fast, database-driven hierarchical search across all branches with full breadcrumbs
      searchDebounceRef.current = setTimeout(async () => {
        try {
          const results = await api.searchCategories(clean, locale, 20);
          if (Array.isArray(results)) {
            const mapped: SearchEntry[] = results.map(r => {
              const bcs = Array.isArray(r.breadcrumbs) ? r.breadcrumbs : [];
              const rootSlug = bcs[0]?.slug || r.slug;
              const subSlug = bcs[1]?.slug;
              const itemSlug = bcs[2]?.slug;
              // Clean parent path without repeating current item's own name
              const parentBcs = bcs.length > 1 ? bcs.slice(0, -1) : [];
              const parentPath = parentBcs.map((b: any) => b.name).join(' › ');
              return {
                id: r.id,
                type: r.level === 0 ? 'category' : r.level === 1 ? 'subcategory' : 'item',
                label: r.name || r.name_en,
                labelEn: r.name_en,
                catId: bcs[0]?.id || r.id,
                catSlug: rootSlug,
                subSlug: subSlug,
                itemSlug: itemSlug,
                breadcrumb: parentPath ? `in ${parentPath}` : '',
                hierarchy: parentPath,
                level: r.level,
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

      // Fetch location suggestions in parallel
      locationDebounceRef.current = setTimeout(async () => {
        try {
          const { api } = await import('@/lib/api');
          const results = await api.searchLocations(clean);
          setLocationSuggestions(results.slice(0, 3));
        } catch { setLocationSuggestions([]); }
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
      if (pathname === '/products') {
        debounceRef.current = setTimeout(() => {
          router.replace('/products');
        }, 300);
      }
    }
  };

  const handleSubmit = (ev: FormEvent<HTMLFormElement>) => {
    ev.preventDefault();
    if (debounceRef.current) clearTimeout(debounceRef.current);
    setOpen(false);
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
      setHighlighted(h => Math.min(h + 1, displayList.length - 1));
    } else if (ev.key === 'ArrowUp') {
      ev.preventDefault();
      setHighlighted(h => Math.max(h - 1, -1));
    } else if (ev.key === 'Escape') {
      setOpen(false);
    }
  };

  const clear = () => {
    setQ('');
    setOpen(false);
    setLocationSuggestions([]);
    if (debounceRef.current) clearTimeout(debounceRef.current);
    router.push('/products');
    inputRef.current?.focus();
  };

  // Highlight matching chars in label
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

      {/* Dropdown */}
      {open && (
        <div className="absolute top-full left-0 right-0 mt-2 rounded-2xl border shadow-2xl z-[200] bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 max-h-80 overflow-y-auto">

          {/* Location suggestions — shown when query matches a location */}
          {(() => {
            const seen = new Set<string>();
            const uniqueLocs = locationSuggestions.filter(r => {
              const p = locPath(r, locale);
              if (!p || seen.has(p)) return false;
              seen.add(p);
              return true;
            });
            if (uniqueLocs.length === 0) return null;

            const locBadge = locale === 'am' ? 'አካባቢ' : locale === 'ti' ? 'ቦታ' : locale === 'or' ? 'bakka' : 'location';

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
                      onMouseDown={() => { setOpen(false); router.push(buildLocationUrl(r)); }}
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

          {/* Header label */}
          <div className="px-3 pt-2.5 pb-1">
            <p className="text-[10px] font-black uppercase tracking-widest opacity-30">
              {q.trim() ? `${u.results_for || 'Results for'} "${q}"` : (u.browse_categories || 'Browse Categories')}
            </p>
          </div>

          {displayList.length === 0 && q.trim().length >= 1 && (
            <div className="px-4 py-6 text-center">
              <p className="text-sm font-bold opacity-40">{u.no_match_for || 'No match for'} &quot;{q}&quot;</p>
              <p className="text-xs opacity-30 mt-1">{u.try_different_word || 'Try a different word or browse categories below'}</p>
              <button
                onClick={() => { setQ(''); setOpen(true); }}
                className="mt-3 text-xs font-bold text-[#1769E0] hover:underline"
              >
                {u.browse_all_categories_arrow || 'Browse all categories →'}
              </button>
            </div>
          )}

          {displayList.map((entry, i) => {
            const Icon = TYPE_ICON[entry.type];
            const label = getLabel(entry);
            const breadcrumb = getBreadcrumb(entry);
            const isActive = i === highlighted;

            return (
              <button
                key={entry.id ? `search-cat-${entry.id}-${i}` : `search-entry-${entry.catSlug}-${entry.subSlug ?? ''}-${entry.itemSlug ?? ''}-${i}`}
                onMouseDown={() => selectEntry(entry)}
                onMouseEnter={() => setHighlighted(i)}
                className={`w-full flex items-center gap-3 px-3 py-2.5 text-left transition-colors ${isActive
                  ? 'bg-[#EAF3FF] dark:bg-[#1769E0]/15'
                  : 'hover:bg-zinc-50 dark:hover:bg-zinc-800/60'
                  }`}
              >
                <div className={`shrink-0 ${TYPE_COLOR[entry.type]}`}>
                  <Icon size={14} />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-bold truncate">
                    {highlight(label, q)}
                  </p>
                  {breadcrumb && (
                    <p className="text-[10px] opacity-40 truncate">{breadcrumb}</p>
                  )}
                </div>
                <span className={`text-[9px] font-black uppercase tracking-wider px-2 py-0.5 rounded-full shrink-0 ${entry.level === 0 || entry.type === 'category'
                  ? 'bg-[#EAF3FF] dark:bg-[#1769E0]/20 text-[#1769E0]'
                  : entry.level === 1
                    ? 'bg-blue-100 dark:bg-blue-900/40 text-blue-600'
                    : entry.level === 2
                      ? 'bg-emerald-100 dark:bg-emerald-900/40 text-emerald-600'
                      : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-500'
                  }`}>
                  {entry.level === 0
                    ? 'Sector'
                    : entry.level === 1
                      ? 'Superfamily'
                      : entry.level === 2
                        ? 'Category'
                        : entry.level === 3
                          ? 'Subcategory'
                          : 'Item'}
                </span>
              </button>
            );
          })}

          {/* Free-text search option when typing */}
          {q.trim().length >= 1 && (
            <div className="border-t border-zinc-100 dark:border-zinc-800">
              <button
                onMouseDown={() => {
                  setOpen(false);
                  router.push(`/products?search=${encodeURIComponent(q.trim())}`);
                }}
                className="w-full flex items-center gap-3 px-3 py-2.5 hover:bg-zinc-50 dark:hover:bg-zinc-800/60 transition-colors"
              >
                <Search size={14} className="text-[#1769E0] shrink-0" />
                <p className="text-sm font-bold">
                  {u.search_products_for_query || 'Search products for'} &quot;<span className="text-[#1769E0]">{q}</span>&quot;
                </p>
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
