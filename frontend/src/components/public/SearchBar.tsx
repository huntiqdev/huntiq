'use client';
import { useState, useRef, useEffect, useCallback, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { useLocale } from '@/hooks/useLocale';
import { api } from '@/lib/api';
import { Search, X, ChevronDown, Tag, Layers, Hash, MapPin, Sparkles, Sliders } from 'lucide-react';

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
  category: 'text-orange-600',
  subcategory: 'text-zinc-500',
  item: 'text-zinc-400',
} as const;

// ── Component ────────────────────────────────────────────────────────────────
export function SearchBar() {
  const { t, locale } = useLocale();
  const router = useRouter();
  const [q, setQ] = useState('');
  const [open, setOpen] = useState(false);
  const [highlighted, setHighlighted] = useState(-1);
  const [searchIndex, setSearchIndex] = useState<SearchEntry[]>([]);
  const [locationSuggestions, setLocationSuggestions] = useState<any[]>([]);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  // Build searchable index dynamically from authoritative database taxonomy
  useEffect(() => {
    api.getCategoryTree(locale).then(tree => {
      if (!Array.isArray(tree)) return;
      const idx: SearchEntry[] = [];
      for (const root of tree) {
        const rootLabel = root.name || root.name_en || root.slug;
        const rootEn = root.name_en || root.name || root.slug;
        idx.push({
          id: root.id,
          type: 'category',
          label: rootLabel,
          labelEn: rootEn,
          catId: root.id,
          catSlug: root.slug,
        });
        if (root.children) {
          for (const sub of root.children) {
            const subLabel = sub.name || sub.name_en || sub.slug;
            const subEn = sub.name_en || sub.name || sub.slug;
            idx.push({
              id: sub.id,
              type: 'subcategory',
              label: subLabel,
              labelEn: subEn,
              catId: root.id,
              catSlug: root.slug,
              subSlug: sub.slug,
              breadcrumb: rootLabel,
            });
            if (sub.children) {
              for (const item of sub.children) {
                const itemLabel = item.name || item.name_en || item.slug;
                const itemEn = item.name_en || item.name || item.slug;
                idx.push({
                  id: item.id,
                  type: 'item',
                  label: itemLabel,
                  labelEn: itemEn,
                  catId: root.id,
                  catSlug: root.slug,
                  subSlug: sub.slug,
                  itemSlug: item.slug,
                  breadcrumb: `${rootLabel} › ${subLabel}`,
                });
              }
            }
          }
        }
      }
      setSearchIndex(idx);
    }).catch(() => {});
  }, [locale]);

  const getLabel = useCallback((e: SearchEntry): string => e.label, []);
  const getLabelEn = useCallback((e: SearchEntry): string => e.labelEn, []);
  const getBreadcrumb = useCallback((e: SearchEntry): string => e.breadcrumb || '', []);

  // Score against both locale label AND English label — whichever is higher wins
  const scoreEntry = useCallback((e: SearchEntry, query: string): number => {
    return Math.max(score(e.label, query), score(e.labelEn, query));
  }, []);

  // Compute suggestions
  const suggestions = q.trim().length >= 1
    ? searchIndex
        .map(e => ({ e, s: scoreEntry(e, q.trim()) }))
        .filter(x => x.s > 0)
        .sort((a, b) => b.s - a.s || a.e.type.localeCompare(b.e.type))
        .slice(0, 12)
        .map(x => x.e)
    : [];

  // All root categories for the "browse" dropdown (shown when input is empty & focused)
  const allCategories = searchIndex.filter(e => e.type === 'category');

  const displayList = q.trim().length >= 1 ? suggestions : allCategories;

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
    if (r.woreda) p.set('woreda', r.woreda);
    else if (r.subcity) p.set('subcity', r.subcity);
    else if (r.zone) p.set('zone', r.zone);
    else if (r.region) p.set('region', r.region);
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
    setQ(ev.target.value);
    setHighlighted(-1);
    setOpen(true);
    if (debounceRef.current) clearTimeout(debounceRef.current);
    if (locationDebounceRef.current) clearTimeout(locationDebounceRef.current);
    const val = ev.target.value.trim();
    if (val) {
      // Fetch location suggestions in parallel
      locationDebounceRef.current = setTimeout(async () => {
        try {
          const { api } = await import('@/lib/api');
          const results = await api.searchLocations(val);
          setLocationSuggestions(results.slice(0, 3));
        } catch { setLocationSuggestions([]); }
      }, 300);
      debounceRef.current = setTimeout(() => {
        router.push(`/products?search=${encodeURIComponent(val)}`);
      }, 500);
    } else {
      setLocationSuggestions([]);
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
          className="w-full pl-9 pr-24 py-2 rounded-full border border-zinc-200 dark:border-zinc-800 bg-zinc-100 dark:bg-zinc-900 text-sm outline-none focus:ring-2 focus:ring-orange-600 transition-all"
        />
        <div className="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1.5">
          {q && (
            <button
              type="button"
              onClick={clear}
              className="text-zinc-400 hover:text-zinc-600 transition-colors p-1"
            >
              <X size={14} />
            </button>
          )}
          <button
            type="button"
            onClick={() => {
              if (q.trim()) {
                router.push(`/ai?q=${encodeURIComponent(q.trim())}`);
              } else {
                router.push('/ai');
              }
            }}
            title="Open HUNTIQ AI Conversational Workspace"
            className="flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-black uppercase tracking-wider transition-all bg-gradient-to-r from-orange-600 via-amber-600 to-orange-500 hover:from-orange-700 hover:to-amber-600 text-white shadow-sm shadow-orange-600/30 active:scale-95"
          >
            <Sparkles size={11} className="animate-pulse" />
            <span>AI Mode</span>
          </button>
        </div>
      </form>

      {/* Dropdown */}
      {open && (
        <div className="absolute top-full left-0 right-0 mt-2 rounded-2xl border shadow-2xl z-[200] bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 max-h-80 overflow-y-auto">

          {/* Location suggestions — shown when query matches a location */}
          {locationSuggestions.length > 0 && (
            <div className="border-b border-zinc-100 dark:border-zinc-800">
              <div className="px-3 pt-2 pb-1">
                <p className="text-[10px] font-black uppercase tracking-widest opacity-30">Locations</p>
              </div>
              {locationSuggestions.map((r, i) => {
                const parts = [r.region, r.zone, r.subcity, r.woreda].filter(Boolean);
                return (
                  <button
                    key={i}
                    onMouseDown={() => { setOpen(false); router.push(buildLocationUrl(r)); }}
                    className="w-full flex items-center gap-3 px-3 py-2.5 text-left hover:bg-orange-50 dark:hover:bg-orange-950/40 transition-colors"
                  >
                    <MapPin size={14} className="text-orange-600 shrink-0" />
                    <p className="text-sm font-bold truncate">{parts.join(' › ')}</p>
                    <span className="text-[9px] font-black uppercase tracking-wider px-1.5 py-0.5 rounded-full shrink-0 bg-orange-100 dark:bg-orange-900/40 text-orange-600">location</span>
                  </button>
                );
              })}
            </div>
          )}

          {/* Header label */}
          <div className="px-3 pt-2.5 pb-1">
            <p className="text-[10px] font-black uppercase tracking-widest opacity-30">
              {q.trim() ? `Results for "${q}"` : 'Browse Categories'}
            </p>
          </div>

          {displayList.length === 0 && q.trim().length >= 1 && (
            <div className="px-4 py-6 text-center">
              <p className="text-sm font-bold opacity-40">No match for &quot;{q}&quot;</p>
              <p className="text-xs opacity-30 mt-1">Try a different word or browse categories below</p>
              <button
                onClick={() => { setQ(''); setOpen(true); }}
                className="mt-3 text-xs font-bold text-orange-600 hover:underline"
              >
                Browse all categories →
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
                key={`${entry.catSlug}-${entry.subSlug ?? ''}-${entry.itemSlug ?? ''}`}
                onMouseDown={() => selectEntry(entry)}
                onMouseEnter={() => setHighlighted(i)}
                className={`w-full flex items-center gap-3 px-3 py-2.5 text-left transition-colors ${
                  isActive
                    ? 'bg-orange-50 dark:bg-orange-950/40'
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
                <span className={`text-[9px] font-black uppercase tracking-wider px-1.5 py-0.5 rounded-full shrink-0 ${
                  entry.type === 'category'
                    ? 'bg-orange-100 dark:bg-orange-900/40 text-orange-600'
                    : entry.type === 'subcategory'
                    ? 'bg-zinc-100 dark:bg-zinc-800 text-zinc-500'
                    : 'bg-zinc-50 dark:bg-zinc-800/50 text-zinc-400'
                }`}>
                  {entry.type}
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
                <Search size={14} className="text-orange-600 shrink-0" />
                <p className="text-sm font-bold">
                  Search products for &quot;<span className="text-orange-600">{q}</span>&quot;
                </p>
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
