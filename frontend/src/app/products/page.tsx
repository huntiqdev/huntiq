'use client';
import { useState, useEffect, Suspense } from 'react';

export const dynamic = 'force-dynamic';
import { useSearchParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';
import { Category, Product, ProductFilters } from '@/types';
import { ProductCard } from '@/components/public/ProductCard';
import { RentalProductCard } from '@/components/public/RentalProductCard';
import { FilterSidebar } from '@/components/public/FilterSidebar';
import { Spinner } from '@/components/ui/Spinner';
import { useLocale } from '@/hooks/useLocale';
import { resolveCategoryTitle } from '@/lib/taxonomy';
import { useAuth } from '@/context/AuthContext';
import {
  ShoppingBag, ChevronDown, SlidersHorizontal, X,
  LayoutGrid, ChevronRight, TrendingUp, MapPin, Star,
  ShieldCheck, Sliders, ShieldAlert, Navigation, AlertTriangle, RotateCcw,
} from 'lucide-react';
import { formatLocationParts } from '@/types/location';

// ── Tier banner shown to buyers ───────────────────────────────────────────────
function TierBanner() {
  const { isBuyer, buyerTier } = useAuth();
  const { t } = useLocale();
  if (!isBuyer || buyerTier === 'premium') return null;
  const u = t.ui;
  return (
    <div className={`mx-[2%] sm:mx-[4%] lg:mx-[6%] mt-3 rounded-2xl px-4 py-2.5 flex items-center gap-3 text-xs font-semibold border ${buyerTier === 'gold'
      ? 'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800 text-yellow-800 dark:text-yellow-300'
      : 'bg-zinc-50 dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 text-zinc-600 dark:text-zinc-400'
      }`}>
      <Star size={13} className={buyerTier === 'gold' ? 'text-yellow-500' : 'text-zinc-400'} />
      <span>
        {buyerTier === 'gold'
          ? (u.gold_buyer_banner ?? 'Gold buyer — you can see Gold seller images. Upgrade to Premium for all images.')
          : (u.normal_buyer_banner ?? 'Normal buyer — some product images are blurred. Upgrade to Gold or Premium to unlock.')}
      </span>
      <Link href="/buyer/dashboard" className="ml-auto shrink-0 font-black text-orange-600 hover:underline">
        {u.upgrade_arrow ?? 'Upgrade →'}
      </Link>
    </div>
  );
}

// ── Popular Categories strip ──────────────────────────────────────────────────
function PopularCategoriesStrip({ onSelect }: {
  categories?: { id: number; name: string; slug: string }[];
  onSelect: (catId: number) => void;
}) {
  const [popular, setPopular] = useState<any[]>([]);
  const { t } = useLocale();

  useEffect(() => {
    api.getPopularCategories(7, 8).then(setPopular).catch(() => { });
  }, []);

  if (popular.length === 0) return null;

  const catLabel = (slug: string) => slug.replace(/-/g, ' ');

  const EMOJI: Record<string, string> = {
    agriculture: '🌾', 'food-beverage': '🍎', electronics: '📱', 'consumer-electronics': '📱',
    'computing-it': '💻', 'transport-mobility': '🚗', 'construction-infrastructure': '🔨',
    'property-real-estate': '🏠', 'healthcare-medicine': '🏥', 'energy-utilities': '⚡',
    'industrial-manufacturing': '🏭', 'clothing-apparel': '👗', 'home-appliances-kitchenware': '🍳',
    'furniture-furnishings': '🛋️', 'beauty-personal-care': '💄', 'sports-fitness': '⚽',
    'education-training': '📚', 'finance-insurance': '💳', 'commerce-trade': '🛒',
    'workforce-employment': '💼', 'logistics-supply-chain': '🚚', 'other': '📦',
  };

  return (
    <div className="px-4 py-3 border-b border-zinc-100 dark:border-zinc-800 bg-gradient-to-r from-orange-50/60 to-transparent dark:from-orange-950/20">
      <div className="flex items-center gap-2 mb-2">
        <TrendingUp size={12} className="text-orange-500" />
        <p className="text-[10px] font-black uppercase tracking-widest text-orange-600">{t.ui.trending_this_week}</p>
      </div>
      <div className="flex gap-2 overflow-x-auto no-scrollbar pb-1">
        {popular.map(p => (
          <button
            key={p.id}
            type="button"
            onClick={() => onSelect(p.id)}
            className="shrink-0 flex flex-col items-center gap-1 bg-white dark:bg-zinc-900 rounded-xl border border-zinc-100 dark:border-zinc-800 hover:border-orange-400 hover:shadow-md transition-all px-3 py-2 min-w-[64px]"
          >
            <span className="text-lg">{EMOJI[p.slug] ?? '📦'}</span>
            <span className="text-[9px] font-bold text-zinc-600 dark:text-zinc-400 text-center leading-tight capitalize whitespace-nowrap">
              {catLabel(p.slug)}
            </span>
            {p.activity_score > 0 && (
              <span className="text-[8px] text-orange-500 font-black">
                {p.new_listings > 0 ? `+${p.new_listings}` : ''}{p.order_count > 0 ? ` ${p.order_count}🛒` : ''}
              </span>
            )}
          </button>
        ))}
      </div>
    </div>
  );
}

// ── Main page ─────────────────────────────────────────────────────────────────
function ProductsContent() {
  const { t, locale } = useLocale();
  const router = useRouter();
  const searchParams = useSearchParams();

  const [categories, setCategories] = useState<Category[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState<string | null>(null);
  const [retryKey, setRetryKey] = useState(0);
  const [filters, setFilters] = useState<ProductFilters>({});
  const [visibleCount, setVisibleCount] = useState(24);
  const [filterOpen, setFilterOpen] = useState(false);
  const [activeCatId, setActiveCatId] = useState<number | undefined>();
  const [activeSub, setActiveSub] = useState<string | undefined>();
  const [activeItem, setActiveItem] = useState<string | undefined>();
  const [activeRegion, setActiveRegion] = useState<string>('');
  const [activeZone, setActiveZone] = useState<string>('');
  const [activeSubcity, setActiveSubcity] = useState<string>('');
  const [activeWoreda, setActiveWoreda] = useState<string>('');
  const [activeCondition, setActiveCondition] = useState<'all' | 'new' | 'used'>('all');
  const [listingType, setListingType] = useState<'all' | 'sale' | 'rent'>('all');

  // Master Filter State: Price & Proximity
  const [minPriceInput, setMinPriceInput] = useState<string>('');
  const [maxPriceInput, setMaxPriceInput] = useState<string>('');
  const [activeMinPrice, setActiveMinPrice] = useState<number | undefined>(undefined);
  const [activeMaxPrice, setActiveMaxPrice] = useState<number | undefined>(undefined);
  const [pricePopoverOpen, setPricePopoverOpen] = useState(false);
  const [userCoords, setUserCoords] = useState<{ lat: number; lng: number } | null>(null);
  const [isLocating, setIsLocating] = useState(false);

  useEffect(() => {
    api.getCategoryTree(locale)
      .then(tree => {
        if (Array.isArray(tree) && tree.length > 0) {
          setCategories(tree);
        } else {
          api.getCategories().then(setCategories).catch(() => {});
        }
      })
      .catch(() => {
        api.getCategories().then(setCategories).catch(() => {});
      });
  }, [locale]);

  useEffect(() => {
    if (categories.length === 0) return;
    const f: ProductFilters = {};
    const search = searchParams.get('search');
    const catId = searchParams.get('category_id');
    const catHierarchyId = searchParams.get('category_hierarchy_id');
    const catSlugParam = searchParams.get('cat');
    const sub = searchParams.get('subcategory');
    const item = searchParams.get('item');

    const typeParam = searchParams.get('type');
    if (typeParam === 'sale' || typeParam === 'rent' || typeParam === 'all') {
      setListingType(typeParam);
    }

    if (search) f.search = search;

    const findCategoryInTree = (nodes: Category[], slugOrId: string | number): Category | undefined => {
      for (const n of nodes) {
        if (typeof slugOrId === 'number' && n.id === slugOrId) return n;
        if (typeof slugOrId === 'string' && (n.slug === slugOrId || String(n.id) === slugOrId)) return n;
        if (n.children && n.children.length > 0) {
          const match = findCategoryInTree(n.children, slugOrId);
          if (match) return match;
        }
      }
      return undefined;
    };

    let resolvedCatId: number | undefined;
    if (catHierarchyId) {
      resolvedCatId = Number(catHierarchyId);
    } else if (catId) {
      resolvedCatId = Number(catId);
    } else if (catSlugParam) {
      const found = findCategoryInTree(categories, catSlugParam);
      resolvedCatId = found?.id;
    }

    if (resolvedCatId) {
      f.category_id = resolvedCatId;
      f.category_hierarchy_id = resolvedCatId;
      setActiveCatId(resolvedCatId);
    } else {
      setActiveCatId(undefined);
    }

    if (sub) { f.subcategory = sub; setActiveSub(sub); } else setActiveSub(undefined);
    if (item) { f.item = item; setActiveItem(item); } else setActiveItem(undefined);

    const woreda = searchParams.get('woreda');
    const subcity = searchParams.get('subcity');
    const zone = searchParams.get('zone');
    const region = searchParams.get('region');

    if (woreda) { f.woreda = woreda; setActiveWoreda(woreda); } else { setActiveWoreda(''); }
    if (subcity) { f.subcity = subcity; setActiveSubcity(subcity); } else { setActiveSubcity(''); }
    if (zone) { f.zone = zone; setActiveZone(zone); } else { setActiveZone(''); }
    if (region) { f.region = region; setActiveRegion(region); } else { setActiveRegion(''); }

    // Master Filters: Price min/max
    const minP = searchParams.get('min_price');
    const maxP = searchParams.get('max_price');
    if (minP) {
      const p = parseFloat(minP);
      if (!isNaN(p)) {
        f.min_price = p;
        setActiveMinPrice(p);
        setMinPriceInput(minP);
      }
    } else {
      setActiveMinPrice(undefined);
      setMinPriceInput('');
    }
    if (maxP) {
      const p = parseFloat(maxP);
      if (!isNaN(p)) {
        f.max_price = p;
        setActiveMaxPrice(p);
        setMaxPriceInput(maxP);
      }
    } else {
      setActiveMaxPrice(undefined);
      setMaxPriceInput('');
    }

    // Master Filters: Condition
    const condParam = searchParams.get('condition');
    if (condParam === 'new' || condParam === 'used') {
      setActiveCondition(condParam);
    } else {
      setActiveCondition('all');
    }

    // Master Filters: GPS proximity coordinates
    const latParam = searchParams.get('lat');
    const lngParam = searchParams.get('lng');
    if (latParam && lngParam) {
      const lat = parseFloat(latParam);
      const lng = parseFloat(lngParam);
      if (!isNaN(lat) && !isNaN(lng)) {
        f.lat = lat;
        f.lng = lng;
        setUserCoords({ lat, lng });
      }
    } else {
      setUserCoords(null);
    }

    setFilters(f);
  }, [searchParams, categories]);

  useEffect(() => {
    setLoading(true);
    setFetchError(null);
    setVisibleCount(24);
    const apiFilters: ProductFilters = { ...filters };
    if (activeCondition !== 'all') apiFilters.condition = activeCondition;
    if (listingType !== 'all') apiFilters.listing_type = listingType;
    if (activeMinPrice !== undefined) apiFilters.min_price = activeMinPrice;
    if (activeMaxPrice !== undefined) apiFilters.max_price = activeMaxPrice;
    if (userCoords) {
      apiFilters.lat = userCoords.lat;
      apiFilters.lng = userCoords.lng;
      apiFilters.sort_by = 'nearest';
    }
    apiFilters.locale = locale;
    apiFilters.preferred_language = locale;

    api.getProducts(apiFilters)
      .then(data => {
        setProducts(data);
        setFetchError(null);
      })
      .catch((err: any) => {
        setProducts([]);
        setFetchError(err?.message || 'Failed to load products. Please check your network connection.');
      })
      .finally(() => setLoading(false));
  }, [filters, activeCondition, listingType, activeMinPrice, activeMaxPrice, userCoords, locale, retryKey]);

  const findNodeById = (nodes: Category[], id?: number): Category | undefined => {
    if (!id) return undefined;
    for (const n of nodes) {
      if (n.id === id) return n;
      if (n.children && n.children.length > 0) {
        const found = findNodeById(n.children, id);
        if (found) return found;
      }
    }
    return undefined;
  };

  const activeCatNode = findNodeById(categories, activeCatId);
  const activeCatSlug = activeCatNode?.slug ?? '';
  const subcategoryNodes = activeCatNode?.children ?? [];
  const activeSubNode = subcategoryNodes.find(s => s.slug === activeSub || String(s.id) === activeSub);
  const itemNodes = activeSubNode?.children ?? [];

  const activeLocParts = formatLocationParts({ region: activeRegion, zone: activeZone, subcity: activeSubcity, woreda: activeWoreda }, locale);
  const activeLocLabel = activeLocParts.join(' › ');

  const buildUrl = (
    optsOrCatId?: number | {
      catId?: number;
      sub?: string;
      item?: string;
      search?: string;
      loc?: { region?: string; zone?: string; subcity?: string; woreda?: string };
      overrideType?: 'all' | 'sale' | 'rent';
      overrideCond?: 'all' | 'new' | 'used';
      minP?: string;
      maxP?: string;
      coords?: { lat: number; lng: number } | null;
    },
    subArg?: string,
    itemArg?: string,
    searchArg?: string,
    locArg?: { region?: string; zone?: string; subcity?: string; woreda?: string },
    overrideTypeArg?: 'all' | 'sale' | 'rent'
  ) => {
    let opts: {
      catId?: number;
      sub?: string;
      item?: string;
      search?: string;
      loc?: { region?: string; zone?: string; subcity?: string; woreda?: string };
      overrideType?: 'all' | 'sale' | 'rent';
      overrideCond?: 'all' | 'new' | 'used';
      minP?: string;
      maxP?: string;
      coords?: { lat: number; lng: number } | null;
    } = {};

    if (typeof optsOrCatId === 'number') {
      opts = {
        catId: optsOrCatId,
        sub: subArg,
        item: itemArg,
        search: searchArg,
        loc: locArg,
        overrideType: overrideTypeArg,
      };
    } else if (optsOrCatId && typeof optsOrCatId === 'object') {
      opts = optsOrCatId;
    }
    const p = new URLSearchParams();
    const typeToUse = opts?.overrideType !== undefined ? opts.overrideType : listingType;
    if (typeToUse && typeToUse !== 'all') p.set('type', typeToUse);

    const condToUse = opts?.overrideCond !== undefined ? opts.overrideCond : activeCondition;
    if (condToUse && condToUse !== 'all') p.set('condition', condToUse);

    const catToUse = opts?.catId !== undefined ? opts.catId : activeCatId;
    if (catToUse) {
      p.set('category_id', String(catToUse));
      p.set('category_hierarchy_id', String(catToUse));
    }

    const subToUse = opts?.sub !== undefined ? opts.sub : activeSub;
    if (subToUse) p.set('subcategory', subToUse);

    const itemToUse = opts?.item !== undefined ? opts.item : activeItem;
    if (itemToUse) p.set('item', itemToUse);

    const minPToUse = opts?.minP !== undefined ? opts.minP : (activeMinPrice !== undefined ? String(activeMinPrice) : '');
    if (minPToUse) p.set('min_price', minPToUse);

    const maxPToUse = opts?.maxP !== undefined ? opts.maxP : (activeMaxPrice !== undefined ? String(activeMaxPrice) : '');
    if (maxPToUse) p.set('max_price', maxPToUse);

    const coordsToUse = opts?.coords !== undefined ? opts.coords : userCoords;
    if (coordsToUse) {
      p.set('lat', coordsToUse.lat.toFixed(6));
      p.set('lng', coordsToUse.lng.toFixed(6));
      p.set('sort_by', 'nearest');
    }

    const l = opts?.loc ?? { region: activeRegion || undefined, zone: activeZone || undefined, subcity: activeSubcity || undefined, woreda: activeWoreda || undefined };
    if (l.woreda) p.set('woreda', l.woreda);
    else if (l.subcity) p.set('subcity', l.subcity);
    else if (l.zone) p.set('zone', l.zone);
    else if (l.region) p.set('region', l.region);

    return `/products${p.toString() ? '?' + p.toString() : ''}`;
  };

  const detectNearMe = () => {
    if (typeof window === 'undefined' || !navigator?.geolocation) {
      alert('Geolocation is not available on this device');
      return;
    }
    setIsLocating(true);
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        setIsLocating(false);
        const lat = pos.coords.latitude;
        const lng = pos.coords.longitude;
        setUserCoords({ lat, lng });
        router.push(buildUrl({ coords: { lat, lng } }));
      },
      (err) => {
        setIsLocating(false);
        console.warn('Geolocation error:', err);
      },
      { timeout: 8000, enableHighAccuracy: true }
    );
  };

  const getCatDisplayName = (cat?: Category) => {
    if (!cat) return '';
    return resolveCategoryTitle(cat as any, locale) || cat.name || (cat as any).name_en || cat.slug?.replace(/_/g, ' ') || '';
  };

  const getSubDisplayName = (sub?: Category) => {
    if (!sub) return '';
    return resolveCategoryTitle(sub as any, locale) || sub.name || (sub as any).name_en || sub.slug?.replace(/_/g, ' ') || '';
  };

  const getItemDisplayName = (item?: Category | string) => {
    if (!item) return '';
    if (typeof item === 'object') {
      return resolveCategoryTitle(item as any, locale) || item.name || (item as any).name_en || item.slug?.replace(/_/g, ' ') || '';
    }
    return item.replace(/_/g, ' ');
  };

  return (
    <>
      {/* ── 3rd Row: Categories list (Compact sticky bar, single row) ── */}
      <div className="sticky top-14 z-30 bg-white/95 dark:bg-zinc-950/95 backdrop-blur border-b border-zinc-100 dark:border-zinc-800">
        <div className="flex items-center gap-1.5 overflow-x-auto no-scrollbar py-2 px-[2%] sm:px-[4%] lg:px-[6%]">
          <button onClick={() => { setActiveCatId(undefined); setActiveSub(undefined); setActiveItem(undefined); router.push(buildUrl()); }}
            className={`flex items-center gap-1 px-3 py-1.5 rounded-full whitespace-nowrap font-bold text-xs transition-all shrink-0 ${!activeCatId ? 'bg-orange-600 text-white shadow-md shadow-orange-600/20' : 'bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700'
              }`}>
            <LayoutGrid size={12} /> {t.categories.all}
          </button>
          {categories.map(cat => (
            <button key={cat.id}
              onClick={() => { setActiveCatId(cat.id); setActiveSub(undefined); setActiveItem(undefined); router.push(buildUrl(cat.id)); }}
              className={`px-3 py-1.5 rounded-full whitespace-nowrap font-bold text-xs transition-all shrink-0 ${activeCatId === cat.id ? 'bg-orange-600 text-white shadow-md shadow-orange-600/20' : 'bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700'
                }`}>
              {getCatDisplayName(cat)}
            </button>
          ))}
        </div>
      </div>

      {/* Systematic Subcategory & Item Navigation Banner (Non-sticky, compact — never covers the mobile screen on scroll) */}
      {activeCatNode && (
        <div className="mx-[2%] sm:mx-[4%] lg:mx-[6%] mt-2.5 p-3 rounded-2xl bg-zinc-50 dark:bg-zinc-900 border border-zinc-200/80 dark:border-zinc-800/80 flex flex-col gap-2 shadow-sm">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1.5 text-xs font-black">
              <span className="text-orange-600 font-extrabold">{getCatDisplayName(activeCatNode)}</span>
              {activeSubNode && (
                <>
                  <ChevronRight size={11} className="text-zinc-400" />
                  <span className="text-zinc-800 dark:text-zinc-200">{getSubDisplayName(activeSubNode)}</span>
                </>
              )}
              {activeItem && (
                <>
                  <ChevronRight size={11} className="text-zinc-400" />
                  <span className="text-zinc-600 dark:text-zinc-400">{getItemDisplayName(activeItem)}</span>
                </>
              )}
            </div>
            <button
              onClick={() => { setActiveCatId(undefined); setActiveSub(undefined); setActiveItem(undefined); router.push(buildUrl()); }}
              className="text-[11px] font-bold text-zinc-400 hover:text-orange-600 flex items-center gap-1 transition-colors"
            >
              <X size={12} />
              <span>{t.ui.clear_filter ?? 'Clear'}</span>
            </button>
          </div>

          {subcategoryNodes.length > 0 && (
            <div className="flex items-center gap-1.5 overflow-x-auto no-scrollbar pt-1">
              {subcategoryNodes.map(s => (
                <button
                  key={s.id || s.slug}
                  onClick={() => { setActiveSub(activeSub === s.slug ? undefined : s.slug); setActiveItem(undefined); router.push(buildUrl(activeCatId, activeSub === s.slug ? undefined : s.slug)); }}
                  className={`px-3 py-1 rounded-xl text-xs font-bold whitespace-nowrap transition-all ${
                    activeSub === s.slug
                      ? 'bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 shadow-sm'
                      : 'bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 text-zinc-700 dark:text-zinc-300 hover:border-orange-500'
                  }`}
                >
                  {getSubDisplayName(s)}
                </button>
              ))}
            </div>
          )}

          {itemNodes.length > 0 && (
            <div className="flex items-center gap-1.5 overflow-x-auto no-scrollbar pt-1 border-t border-zinc-200/60 dark:border-zinc-800/60">
              {itemNodes.map(i => (
                <button
                  key={i.id || i.slug}
                  onClick={() => { setActiveItem(activeItem === i.slug ? undefined : i.slug); router.push(buildUrl(activeCatId, activeSub, activeItem === i.slug ? undefined : i.slug)); }}
                  className={`px-2.5 py-0.5 rounded-lg text-[10px] font-bold whitespace-nowrap transition-all ${
                    activeItem === i.slug
                      ? 'bg-orange-600 text-white shadow-sm'
                      : 'bg-white dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 hover:text-orange-600'
                  }`}
                >
                  {getItemDisplayName(i)}
                </button>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Tier banner */}
      <TierBanner />

      {/* Location filter indicator */}
      {activeLocLabel && (
        <div className="mx-[2%] sm:mx-[4%] lg:mx-[6%] mt-2 flex items-center gap-2 text-xs text-zinc-500 dark:text-zinc-400">
          <MapPin size={11} className="text-orange-500" />
          <span>{t.ui.showing_results_in} <strong className="text-zinc-700 dark:text-zinc-300">{activeLocLabel}</strong></span>
          <button onClick={() => { setActiveRegion(''); setActiveZone(''); setActiveSubcity(''); setActiveWoreda(''); router.push(buildUrl(activeCatId, activeSub, activeItem, filters.search, {})); }}
            className="ml-1 text-orange-600 hover:underline font-bold">{t.ui.clear_filter}</button>
        </div>
      )}

      {/* Popular categories strip (clean full-width, single floating AI button used instead) */}
      {!loading && (
        <div className="px-[2%] sm:px-[4%] lg:px-[6%] pt-1">
          <PopularCategoriesStrip
            categories={categories}
            onSelect={catId => { setActiveCatId(catId); setActiveSub(undefined); setActiveItem(undefined); router.push(buildUrl(catId)); }}
          />
        </div>
      )}

      {/* ── Page body ── */}
      <div className="max-w-screen-2xl mx-auto px-[2%] sm:px-[4%] lg:px-[6%] py-4">
        <div className="flex gap-4">

          {/* Sidebar — desktop only, narrow */}
          <aside className="hidden lg:block w-52 shrink-0">
            <div className="sticky top-36">
              <FilterSidebar categories={categories} onFilter={f => { setFilters(f); setFilterOpen(false); }} />
            </div>
          </aside>

          {/* Products area */}
          <div className="flex-1 min-w-0">

            {/* ── MASTER FILTER BAR ── */}
            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-2.5 mb-3 shadow-xs flex flex-wrap items-center justify-between gap-2.5">
              <div className="flex flex-wrap items-center gap-2">
                {/* 1. Transaction Master Filter: Buy / Rent / All */}
                <div className="flex items-center bg-zinc-100 dark:bg-zinc-800 p-0.5 rounded-xl border border-zinc-200/60 dark:border-zinc-700/60">
                  <button
                    type="button"
                    onClick={() => { setListingType('all'); router.push(buildUrl({ overrideType: 'all' })); }}
                    className={`px-2.5 py-1 text-[11px] font-black uppercase rounded-lg transition-all ${
                      listingType === 'all'
                        ? 'bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 shadow-xs'
                        : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'
                    }`}
                  >
                    {(t.ui as any)?.all_modes ?? 'All'}
                  </button>
                  <button
                    type="button"
                    onClick={() => { setListingType('sale'); router.push(buildUrl({ overrideType: 'sale' })); }}
                    className={`px-2.5 py-1 text-[11px] font-black uppercase rounded-lg transition-all flex items-center gap-1 ${
                      listingType === 'sale'
                        ? 'bg-orange-600 text-white shadow-xs'
                        : 'text-zinc-500 hover:text-orange-600'
                    }`}
                  >
                    🏷️ {t.nav?.buy ?? 'Buy'}
                  </button>
                  <button
                    type="button"
                    onClick={() => { setListingType('rent'); router.push(buildUrl({ overrideType: 'rent' })); }}
                    className={`px-2.5 py-1 text-[11px] font-black uppercase rounded-lg transition-all flex items-center gap-1 ${
                      listingType === 'rent'
                        ? 'bg-teal-600 text-white shadow-xs'
                        : 'text-zinc-500 hover:text-teal-600'
                    }`}
                  >
                    🏠 {t.nav?.rental ?? 'Rent'}
                  </button>
                </div>

                {/* 2. Condition Master Filter: New / Used / All */}
                <div className="flex items-center bg-zinc-100 dark:bg-zinc-800 p-0.5 rounded-xl border border-zinc-200/60 dark:border-zinc-700/60">
                  {(['all', 'new', 'used'] as const).map(c => (
                    <button
                      key={c}
                      type="button"
                      onClick={() => { setActiveCondition(c); router.push(buildUrl({ overrideCond: c })); }}
                      className={`px-2 py-1 text-[10px] font-black uppercase rounded-lg transition-all ${
                        activeCondition === c
                          ? c === 'new'
                            ? 'bg-emerald-600 text-white shadow-xs'
                            : c === 'used'
                            ? 'bg-blue-600 text-white shadow-xs'
                            : 'bg-zinc-700 dark:bg-zinc-200 text-white dark:text-zinc-900 shadow-xs'
                          : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'
                      }`}
                    >
                      {c === 'all' ? (t.ui?.condition_all ?? 'All') : c === 'new' ? (t.ui?.condition_new_badge ?? '🆕 New') : (t.ui?.condition_used_badge ?? '♻️ Used')}
                    </button>
                  ))}
                </div>

                {/* 3. Price Filter Pill (Category-scoped or Root-wide) */}
                <div className="relative">
                  <button
                    type="button"
                    onClick={() => setPricePopoverOpen(!pricePopoverOpen)}
                    className={`px-2.5 py-1 rounded-xl border text-[11px] font-bold flex items-center gap-1.5 transition-all ${
                      activeMinPrice !== undefined || activeMaxPrice !== undefined
                        ? 'bg-orange-50 dark:bg-orange-950/40 border-orange-500 text-orange-600 dark:text-orange-400 font-extrabold'
                        : 'bg-zinc-100 dark:bg-zinc-800 border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-300 hover:border-zinc-300'
                    }`}
                  >
                    <span>💰 {activeMinPrice !== undefined || activeMaxPrice !== undefined ? `${activeMinPrice ?? 0} - ${activeMaxPrice ?? '∞'} ETB` : (t.ui?.price ?? 'Price (ETB)')}</span>
                    <ChevronDown size={11} className={pricePopoverOpen ? 'rotate-180 transition-transform' : 'transition-transform'} />
                  </button>

                  {pricePopoverOpen && (
                    <div className="absolute left-0 top-full mt-2 z-50 w-64 p-3 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-700 rounded-2xl shadow-xl space-y-2">
                      <p className="text-[10px] font-black uppercase text-zinc-500">
                        {activeCatNode ? `${getCatDisplayName(activeCatNode)} Price (ETB)` : 'Marketplace Price Range (ETB)'}
                      </p>
                      <div className="flex items-center gap-2">
                        <input
                          type="number"
                          placeholder="Min"
                          value={minPriceInput}
                          onChange={e => setMinPriceInput(e.target.value)}
                          className="w-1/2 px-2.5 py-1.5 text-xs bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl outline-none focus:ring-1 focus:ring-orange-600"
                        />
                        <span className="text-zinc-400 text-xs">-</span>
                        <input
                          type="number"
                          placeholder="Max"
                          value={maxPriceInput}
                          onChange={e => setMaxPriceInput(e.target.value)}
                          className="w-1/2 px-2.5 py-1.5 text-xs bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl outline-none focus:ring-1 focus:ring-orange-600"
                        />
                      </div>
                      <div className="flex items-center justify-between pt-1">
                        <button
                          type="button"
                          onClick={() => {
                            setMinPriceInput('');
                            setMaxPriceInput('');
                            setActiveMinPrice(undefined);
                            setActiveMaxPrice(undefined);
                            setPricePopoverOpen(false);
                            router.push(buildUrl({ minP: '', maxP: '' }));
                          }}
                          className="text-[10px] font-bold text-zinc-400 hover:text-red-500 cursor-pointer"
                        >
                          {(t.ui as any).reset ?? 'Reset'}
                        </button>
                        <button
                          type="button"
                          onClick={() => {
                            setActiveMinPrice(minPriceInput ? parseFloat(minPriceInput) : undefined);
                            setActiveMaxPrice(maxPriceInput ? parseFloat(maxPriceInput) : undefined);
                            setPricePopoverOpen(false);
                            router.push(buildUrl({ minP: minPriceInput, maxP: maxPriceInput }));
                          }}
                          className="px-3 py-1 bg-orange-600 text-white rounded-lg text-xs font-bold hover:bg-orange-700 transition-colors cursor-pointer"
                        >
                          {(t.ui as any).apply ?? t.ui?.submit ?? 'Apply'}
                        </button>
                      </div>
                    </div>
                  )}
                </div>

                {/* 4. Proximity / Near Me Master Filter */}
                <button
                  type="button"
                  onClick={detectNearMe}
                  disabled={isLocating}
                  className={`px-2.5 py-1 rounded-xl border text-[11px] font-bold flex items-center gap-1.5 transition-all ${
                    userCoords
                      ? 'bg-teal-50 dark:bg-teal-950/40 border-teal-500 text-teal-600 dark:text-teal-400 ring-1 ring-teal-500/30'
                      : 'bg-zinc-100 dark:bg-zinc-800 border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-300 hover:border-teal-500'
                  }`}
                >
                  <Navigation size={11} className={isLocating ? 'animate-spin text-teal-500' : userCoords ? 'text-teal-600 fill-teal-600' : 'text-zinc-400'} />
                  <span>
                    {isLocating
                      ? ((t.ui as any).locating ?? 'Locating…')
                      : userCoords
                      ? ((t.ui as any).nearest_to_me ?? '📍 Nearest to Me')
                      : ((t.ui as any).near_me ?? 'Near Me')}
                  </span>
                  {userCoords && (
                    <span
                      onClick={(e) => {
                        e.stopPropagation();
                        setUserCoords(null);
                        router.push(buildUrl({ coords: null }));
                      }}
                      className="ml-1 hover:text-red-500 font-black cursor-pointer"
                    >
                      ×
                    </span>
                  )}
                </button>
              </div>

              {/* Mobile Filter Toggle & Count */}
              <div className="flex items-center gap-2">
                <p className="text-xs font-bold text-zinc-400">
                  {loading ? '…' : `${products.length} ${t.ui?.results_count ?? 'results'}`}
                </p>
                <button
                  onClick={() => setFilterOpen(true)}
                  className="lg:hidden flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-zinc-100 dark:bg-zinc-800 font-bold text-xs hover:bg-orange-600 hover:text-white transition-all cursor-pointer"
                >
                  <SlidersHorizontal size={12} /> {t.ui?.filters ?? 'Filters'}
                </button>
              </div>
            </div>

            {loading ? (
              <div className="flex items-center justify-center py-20"><Spinner /></div>
            ) : fetchError ? (
              <div className="flex flex-col items-center justify-center py-16 px-4 text-center max-w-md mx-auto space-y-4">
                <div className="w-14 h-14 rounded-2xl bg-red-100 dark:bg-red-950/40 text-red-600 dark:text-red-400 flex items-center justify-center shadow-md">
                  <AlertTriangle size={28} />
                </div>
                <div>
                  <h3 className="text-base font-bold text-zinc-900 dark:text-zinc-100">
                    {(t.ui as any).unable_to_load_listings ?? 'Unable to Load Listings'}
                  </h3>
                  <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">{fetchError}</p>
                </div>
                <button
                  onClick={() => setRetryKey(k => k + 1)}
                  className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-orange-600 hover:bg-orange-700 text-white text-xs font-bold shadow-md transition-all active:scale-95 cursor-pointer"
                >
                  <RotateCcw size={14} />
                  <span>{(t.ui as any).retry ?? 'Retry'}</span>
                </button>
              </div>
            ) : products.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-16 px-4 text-center max-w-md mx-auto space-y-4">
                <div className="w-16 h-16 rounded-3xl bg-zinc-100 dark:bg-zinc-800 text-zinc-400 flex items-center justify-center">
                  <ShoppingBag size={32} />
                </div>
                <div>
                  <h3 className="text-base font-black text-zinc-900 dark:text-zinc-100">{t.ui.no_products ?? 'No listings found'}</h3>
                  <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">
                    {(t.ui as any).no_active_listings_match ?? 'No active listings match your current search or filter criteria.'}
                  </p>
                </div>
                <div className="flex flex-wrap gap-2 justify-center pt-1">
                  <button
                    onClick={() => {
                      setActiveCatId(undefined);
                      setActiveSub(undefined);
                      setActiveItem(undefined);
                      setActiveMinPrice(undefined);
                      setActiveMaxPrice(undefined);
                      setActiveCondition('all');
                      setListingType('all');
                      router.push('/products');
                    }}
                    className="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-orange-600 hover:bg-orange-700 text-white text-xs font-bold transition-all shadow-md active:scale-95 cursor-pointer"
                  >
                    <span>{(t as any).search?.clear_filters ?? (t.ui as any).clear_filter ?? 'Clear All Filters'}</span>
                  </button>
                </div>
                <div className="pt-4 border-t border-zinc-200 dark:border-zinc-800 w-full">
                  <p className="text-[11px] font-bold text-zinc-400 uppercase tracking-wider mb-2.5">
                    {(t.ui as any).explore_categories ?? 'Explore Categories'}
                  </p>
                  <div className="flex flex-wrap gap-1.5 justify-center">
                    {categories.slice(0, 8).map(c => (
                      <Link
                        key={c.id}
                        href={`/products?category_hierarchy_id=${c.id}&cat=${c.slug}`}
                        className="px-2.5 py-1 rounded-lg bg-zinc-100 dark:bg-zinc-800 text-zinc-700 dark:text-zinc-300 hover:text-orange-600 text-xs font-medium transition-colors"
                      >
                        {c.icon ? `${c.icon} ` : '🏷️ '}
                        {resolveCategoryTitle(c, locale)}
                      </Link>
                    ))}
                  </div>
                </div>
              </div>
            ) : (
              <>
                {/*
                  Smartphone-optimized Grid:
                  - 2 columns on mobile (≥320px) for generous touch targets & readable photos
                  - 3 columns on sm (≥640px)
                  - 4 columns on md (≥768px)
                  - 5 columns on lg (≥1024px)
                  - 6 columns on xl (≥1280px)
                */}
                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-3 sm:gap-3.5">
                  {products.slice(0, visibleCount).map(p => (
                    p.listing_type === 'rent'
                      ? <RentalProductCard key={p.id} product={p} />
                      : <ProductCard key={p.id} product={p} />
                  ))}
                </div>

                {products.length > visibleCount && (
                  <div className="flex justify-center mt-8">
                    <button onClick={() => setVisibleCount(v => v + 24)}
                      className="flex flex-col items-center gap-2 group tap-bounce">
                      <span className="text-[10px] font-black uppercase tracking-[0.2em] opacity-30 group-hover:opacity-80 transition-opacity">
                        {t.ui.get_more}
                      </span>
                      <div className="w-10 h-10 rounded-2xl bg-orange-600/10 text-orange-600 flex items-center justify-center group-hover:bg-orange-600 group-hover:text-white transition-all shadow-md">
                        <ChevronDown size={20} />
                      </div>
                    </button>
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </div>

      {/* Mobile filter drawer */}
      {filterOpen && (
        <div className="fixed inset-0 z-50 lg:hidden">
          <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={() => setFilterOpen(false)} />
          <div className="absolute bottom-0 left-0 right-0 max-h-[88vh] overflow-y-auto rounded-t-3xl bg-zinc-50 dark:bg-zinc-950 p-4 pb-[max(24px,env(safe-area-inset-bottom,24px))] shadow-2xl animate-in slide-in-from-bottom duration-200">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-2 font-black text-sm">
                <SlidersHorizontal size={15} className="text-orange-600" />
                {t.ui.filters}
              </div>
              <button onClick={() => setFilterOpen(false)}
                className="p-2 rounded-full hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors">
                <X size={16} />
              </button>
            </div>
            <FilterSidebar categories={categories} onFilter={f => { setFilters(f); setFilterOpen(false); }} />
          </div>
        </div>
      )}
    </>
  );
}

export default function ProductsPage() {
  return (
    <Suspense fallback={<div className="flex items-center justify-center h-64"><Spinner /></div>}>
      <ProductsContent />
    </Suspense>
  );
}
