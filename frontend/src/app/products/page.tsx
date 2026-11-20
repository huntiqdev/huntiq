'use client';
import { useState, useEffect, useMemo, Suspense } from 'react';

export const dynamic = 'force-dynamic';
import { useSearchParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';
import { Category, Product, ProductFilters } from '@/types';
import { ProductCard } from '@/components/public/ProductCard';
import { RentalProductCard } from '@/components/public/RentalProductCard';
import { LocationSelector } from '@/components/location/LocationSelector';
import { NearMeFindButton } from '@/components/location/NearMeFindButton';
import { Spinner } from '@/components/ui/Spinner';
import { useLocale } from '@/hooks/useLocale';
import { resolveCategoryTitle, findCategoryPath, getCategoryEmoji } from '@/lib/taxonomy';
import { useAuth } from '@/context/AuthContext';
import {
  ShoppingBag, ChevronDown, SlidersHorizontal, X,
  LayoutGrid, ChevronRight, TrendingUp, MapPin, Star,
  ShieldCheck, Sliders, ShieldAlert, Navigation, AlertTriangle, RotateCcw,
} from 'lucide-react';
import { formatLocationParts } from '@/lib/location';

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
      <Link href="/buyer/dashboard" className="ml-auto shrink-0 font-black text-[#1769E0] dark:text-[#38A5FF] hover:underline">
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
  const [activeCatId, setActiveCatId] = useState<number | undefined>();
  const [activeSub, setActiveSub] = useState<string | undefined>();
  const [activeItem, setActiveItem] = useState<string | undefined>();
  const [activeRegion, setActiveRegion] = useState<string>('');
  const [activeZone, setActiveZone] = useState<string>('');
  const [activeSubcity, setActiveSubcity] = useState<string>('');
  const [activeWoreda, setActiveWoreda] = useState<string>('');
  const [activeCondition, setActiveCondition] = useState<'all' | 'new' | 'used'>('all');
  const [listingType, setListingType] = useState<'all' | 'sale' | 'rent'>('all');
  const [sortBy, setSortBy] = useState<'trending' | 'newest' | 'nearest' | 'views' | 'price_asc' | 'price_desc'>('trending');
  const [maxDistance, setMaxDistance] = useState<number | undefined>(undefined);

  // Filter State: Price & Proximity
  const [minPriceInput, setMinPriceInput] = useState<string>('');
  const [maxPriceInput, setMaxPriceInput] = useState<string>('');
  const [activeMinPrice, setActiveMinPrice] = useState<number | undefined>(undefined);
  const [activeMaxPrice, setActiveMaxPrice] = useState<number | undefined>(undefined);
  const [userCoords, setUserCoords] = useState<{ lat: number; lng: number } | null>(null);
  const [isLocating, setIsLocating] = useState(false);

  // Slide-over Filter Drawer State & Drafts
  const [filterDrawerOpen, setFilterDrawerOpen] = useState(false);
  const [drawerType, setDrawerType] = useState<'all' | 'sale' | 'rent'>('all');
  const [drawerCondition, setDrawerCondition] = useState<'all' | 'new' | 'used'>('all');
  const [drawerMinPrice, setDrawerMinPrice] = useState<string>('');
  const [drawerMaxPrice, setDrawerMaxPrice] = useState<string>('');
  const [drawerRegion, setDrawerRegion] = useState<string>('');
  const [drawerZone, setDrawerZone] = useState<string>('');
  const [drawerSubcity, setDrawerSubcity] = useState<string>('');
  const [drawerWoreda, setDrawerWoreda] = useState<string>('');

  useEffect(() => {
    api.getCategoryTree(locale)
      .then(tree => {
        if (Array.isArray(tree) && tree.length > 0) {
          setCategories(tree);
        } else {
          api.getCategories().then(setCategories).catch(() => { });
        }
      })
      .catch(() => {
        api.getCategories().then(setCategories).catch(() => { });
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

    // Resolve category path from most specific identifier to least specific
    const targetCategoryKey = item || sub || catHierarchyId || catId || catSlugParam;
    const resolvedPath = targetCategoryKey ? findCategoryPath(categories, targetCategoryKey) : [];
    const activeTargetNode = resolvedPath.length > 0 ? resolvedPath[resolvedPath.length - 1] : undefined;

    if (activeTargetNode) {
      f.category_id = activeTargetNode.id;
      f.category_hierarchy_id = activeTargetNode.id;
      setActiveCatId(activeTargetNode.id);
    } else {
      setActiveCatId(undefined);
    }

    if (sub) { f.subcategory = sub; setActiveSub(sub); } else { setActiveSub(undefined); }
    if (item) { f.item = item; setActiveItem(item); } else { setActiveItem(undefined); }

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

    // Sort order: Trending (default), newest, nearest, views, price_asc, price_desc
    const sortParam = (searchParams.get('sort_by') || searchParams.get('sort')) as any;
    if (sortParam && ['trending', 'newest', 'nearest', 'views', 'price_asc', 'price_desc'].includes(sortParam)) {
      setSortBy(sortParam);
      f.sort_by = sortParam;
    } else {
      setSortBy('trending');
      f.sort_by = 'trending';
    }

    const maxDistParam = searchParams.get('max_distance');
    if (maxDistParam) {
      const d = parseFloat(maxDistParam);
      if (!isNaN(d) && d > 0) {
        setMaxDistance(d);
        f.max_distance = d;
      }
    } else {
      setMaxDistance(undefined);
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
    if (sortBy) apiFilters.sort_by = sortBy;
    if (maxDistance) apiFilters.max_distance = maxDistance;
    if (userCoords) {
      apiFilters.lat = userCoords.lat;
      apiFilters.lng = userCoords.lng;
      if (!sortBy || sortBy === 'trending') {
        apiFilters.sort_by = 'nearest';
      }
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
  }, [filters, activeCondition, listingType, activeMinPrice, activeMaxPrice, userCoords, sortBy, maxDistance, locale, retryKey]);

  // Active category lineage: full path from root down to currently active node
  const activeTargetKey = activeItem || activeSub || activeCatId || searchParams.get('item') || searchParams.get('subcategory') || searchParams.get('category_hierarchy_id') || searchParams.get('category_id') || searchParams.get('cat');
  const activePath = useMemo(() => {
    return findCategoryPath(categories, activeTargetKey);
  }, [categories, activeTargetKey]);

  const activeNode = activePath.length > 0 ? activePath[activePath.length - 1] : undefined;
  const parentNode = activePath.length > 1 ? activePath[activePath.length - 2] : undefined;
  const childOptions = activeNode ? (activeNode.children ?? []) : categories;

  const activeLocParts = formatLocationParts({ region: activeRegion, zone: activeZone, subcity: activeSubcity, woreda: activeWoreda }, locale);
  const activeLocLabel = activeLocParts.join(' › ');

  const buildUrl = (
    optsOrCatId?: number | {
      catId?: number;
      catSlug?: string;
      sub?: string;
      item?: string;
      search?: string;
      loc?: { region?: string; zone?: string; subcity?: string; woreda?: string };
      overrideType?: 'all' | 'sale' | 'rent';
      overrideCond?: 'all' | 'new' | 'used';
      minP?: string;
      maxP?: string;
      coords?: { lat: number; lng: number } | null;
      sort?: string;
      maxDist?: number;
    },
    subArg?: string,
    itemArg?: string,
    searchArg?: string,
    locArg?: { region?: string; zone?: string; subcity?: string; woreda?: string },
    overrideTypeArg?: 'all' | 'sale' | 'rent'
  ) => {
    let opts: {
      catId?: number;
      catSlug?: string;
      sub?: string;
      item?: string;
      search?: string;
      loc?: { region?: string; zone?: string; subcity?: string; woreda?: string };
      overrideType?: 'all' | 'sale' | 'rent';
      overrideCond?: 'all' | 'new' | 'used';
      minP?: string;
      maxP?: string;
      coords?: { lat: number; lng: number } | null;
      sort?: string;
      maxDist?: number;
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

    const catToUse = ('catId' in opts) ? opts.catId : activeCatId;
    if (catToUse) {
      p.set('category_id', String(catToUse));
      p.set('category_hierarchy_id', String(catToUse));
    }

    const catSlugToUse = ('catSlug' in opts) ? opts.catSlug : (activePath[0]?.slug);
    if (catSlugToUse && catToUse) p.set('cat', catSlugToUse);

    const subToUse = ('sub' in opts) ? opts.sub : activeSub;
    if (subToUse && catToUse) p.set('subcategory', subToUse);

    const itemToUse = ('item' in opts) ? opts.item : activeItem;
    if (itemToUse && catToUse) p.set('item', itemToUse);

    const searchToUse = ('search' in opts) ? opts.search : filters.search;
    if (searchToUse) p.set('search', searchToUse);

    const minPToUse = opts?.minP !== undefined ? opts.minP : (activeMinPrice !== undefined ? String(activeMinPrice) : '');
    if (minPToUse) p.set('min_price', minPToUse);

    const maxPToUse = opts?.maxP !== undefined ? opts.maxP : (activeMaxPrice !== undefined ? String(activeMaxPrice) : '');
    if (maxPToUse) p.set('max_price', maxPToUse);

    const coordsToUse = opts?.coords !== undefined ? opts.coords : userCoords;
    if (coordsToUse) {
      p.set('lat', coordsToUse.lat.toFixed(6));
      p.set('lng', coordsToUse.lng.toFixed(6));
    }

    const sortToUse = opts?.sort !== undefined ? opts.sort : sortBy;
    if (sortToUse && sortToUse !== 'trending') {
      p.set('sort_by', sortToUse);
    }

    const maxDistToUse = opts?.maxDist !== undefined ? opts.maxDist : maxDistance;
    if (maxDistToUse) {
      p.set('max_distance', String(maxDistToUse));
    }

    const l = opts?.loc ?? { region: activeRegion || undefined, zone: activeZone || undefined, subcity: activeSubcity || undefined, woreda: activeWoreda || undefined };
    if (l.woreda) p.set('woreda', l.woreda);
    else if (l.subcity) p.set('subcity', l.subcity);
    else if (l.zone) p.set('zone', l.zone);
    else if (l.region) p.set('region', l.region);

    return `/products${p.toString() ? '?' + p.toString() : ''}`;
  };

  const navigateToCategory = (target: Category) => {
    const targetPath = findCategoryPath(categories, target.id);
    const root = targetPath[0];
    const sub = targetPath.length > 1 ? targetPath[1] : undefined;
    const item = targetPath.length > 2 ? targetPath[2] : undefined;

    setActiveCatId(target.id);
    setActiveSub(sub?.slug);
    setActiveItem(item?.slug);

    router.push(buildUrl({
      catId: target.id,
      catSlug: root?.slug,
      sub: sub?.slug,
      item: item?.slug,
    }));
  };

  const clearCategorySelection = () => {
    setActiveCatId(undefined);
    setActiveSub(undefined);
    setActiveItem(undefined);
    router.push(buildUrl({
      catId: undefined,
      catSlug: undefined,
      sub: undefined,
      item: undefined,
    }));
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

  const openFilterDrawer = () => {
    setDrawerType(listingType);
    setDrawerCondition(activeCondition);
    setDrawerMinPrice(activeMinPrice !== undefined ? String(activeMinPrice) : (minPriceInput || ''));
    setDrawerMaxPrice(activeMaxPrice !== undefined ? String(activeMaxPrice) : (maxPriceInput || ''));
    setDrawerRegion(activeRegion);
    setDrawerZone(activeZone);
    setDrawerSubcity(activeSubcity);
    setDrawerWoreda(activeWoreda);
    setFilterDrawerOpen(true);
  };

  const applyDrawerFilters = () => {
    setListingType(drawerType);
    setActiveCondition(drawerCondition);
    const minPVal = drawerMinPrice.trim() ? parseFloat(drawerMinPrice) : undefined;
    const maxPVal = drawerMaxPrice.trim() ? parseFloat(drawerMaxPrice) : undefined;
    const resolvedMin = (minPVal !== undefined && !isNaN(minPVal)) ? minPVal : undefined;
    const resolvedMax = (maxPVal !== undefined && !isNaN(maxPVal)) ? maxPVal : undefined;
    setActiveMinPrice(resolvedMin);
    setActiveMaxPrice(resolvedMax);
    setMinPriceInput(drawerMinPrice);
    setMaxPriceInput(drawerMaxPrice);
    setActiveRegion(drawerRegion);
    setActiveZone(drawerZone);
    setActiveSubcity(drawerSubcity);
    setActiveWoreda(drawerWoreda);

    setFilterDrawerOpen(false);
    router.push(buildUrl({
      overrideType: drawerType,
      overrideCond: drawerCondition,
      minP: drawerMinPrice.trim(),
      maxP: drawerMaxPrice.trim(),
      loc: { region: drawerRegion, zone: drawerZone, subcity: drawerSubcity, woreda: drawerWoreda },
    }));
  };

  const resetDrawerFilters = () => {
    setDrawerType('all');
    setDrawerCondition('all');
    setDrawerMinPrice('');
    setDrawerMaxPrice('');
    setDrawerRegion('');
    setDrawerZone('');
    setDrawerSubcity('');
    setDrawerWoreda('');
    setUserCoords(null);

    setListingType('all');
    setActiveCondition('all');
    setActiveMinPrice(undefined);
    setActiveMaxPrice(undefined);
    setMinPriceInput('');
    setMaxPriceInput('');
    setActiveRegion('');
    setActiveZone('');
    setActiveSubcity('');
    setActiveWoreda('');

    setFilterDrawerOpen(false);
    router.push(buildUrl({
      overrideType: 'all',
      overrideCond: 'all',
      minP: '',
      maxP: '',
      coords: null,
      loc: {},
    }));
  };

  const activeFilterCount = useMemo(() => {
    let count = 0;
    if (listingType !== 'all') count++;
    if (activeCondition !== 'all') count++;
    if (activeMinPrice !== undefined || activeMaxPrice !== undefined) count++;
    if (activeRegion || activeZone || activeSubcity || activeWoreda) count++;
    if (userCoords) count++;
    return count;
  }, [listingType, activeCondition, activeMinPrice, activeMaxPrice, activeRegion, activeZone, activeSubcity, activeWoreda, userCoords]);

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
      {/* ── Progressive Drill-Down Category Bar (Sticky) ── */}
      <div className="sticky top-14 z-30 bg-white/95 dark:bg-[#091426]/95 backdrop-blur border-b border-slate-200/80 dark:border-slate-800 shadow-xs">
        <div className="px-[2%] sm:px-[4%] lg:px-[6%] py-2">
          {activeNode ? (
            <div className="flex flex-col gap-2">
              {/* Row 1: Lineage Breadcrumbs + Back Action + Clear */}
              <div className="flex items-center justify-between gap-2">
                <div className="flex items-center gap-1.5 overflow-x-auto no-scrollbar py-0.5 text-xs">
                  {/* Step back button */}
                  <button
                    onClick={() => {
                      if (parentNode) {
                        navigateToCategory(parentNode);
                      } else {
                        clearCategorySelection();
                      }
                    }}
                    className="flex items-center gap-1.5 px-2.5 py-1 rounded-xl bg-slate-100 dark:bg-slate-800 hover:bg-[#1769E0] hover:text-white dark:hover:bg-[#1769E0] text-slate-700 dark:text-slate-300 font-bold transition-all text-xs shrink-0 border border-slate-200/80 dark:border-slate-700/80"
                    title={parentNode ? `Back to ${getCatDisplayName(parentNode)}` : 'Back to All Categories'}
                  >
                    <ChevronRight size={13} className="rotate-180" />
                    <span>{parentNode ? `Back to ${getCatDisplayName(parentNode)}` : t.categories.all}</span>
                  </button>

                  <span className="text-slate-300 dark:text-slate-700 select-none">|</span>

                  {/* Root link */}
                  <button
                    onClick={clearCategorySelection}
                    className="font-semibold text-slate-500 hover:text-[#1769E0] dark:text-slate-400 dark:hover:text-[#38A5FF] transition-colors shrink-0 flex items-center gap-1"
                  >
                    <LayoutGrid size={12} />
                    <span>{t.categories.all}</span>
                  </button>

                  {/* Active Path Breadcrumb Trail */}
                  {activePath.map((stepNode: Category, idx: number) => {
                    const isLast = idx === activePath.length - 1;
                    return (
                      <div key={stepNode.id} className="flex items-center gap-1.5 shrink-0">
                        <ChevronRight size={11} className="text-slate-400 dark:text-slate-600" />
                        {isLast ? (
                          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-xl bg-[#1769E0] text-white font-black text-xs shadow-xs">
                            <span>{getCategoryEmoji(stepNode.slug, stepNode.icon)}</span>
                            <span>{getCatDisplayName(stepNode)}</span>
                          </span>
                        ) : (
                          <button
                            onClick={() => navigateToCategory(stepNode)}
                            className="font-bold text-slate-700 hover:text-[#1769E0] dark:text-slate-300 dark:hover:text-[#38A5FF] transition-colors flex items-center gap-1"
                          >
                            <span>{getCategoryEmoji(stepNode.slug, stepNode.icon)}</span>
                            <span>{getCatDisplayName(stepNode)}</span>
                          </button>
                        )}
                      </div>
                    );
                  })}
                </div>

                {/* Reset / Clear Button */}
                <button
                  onClick={clearCategorySelection}
                  className="text-xs font-bold text-slate-400 hover:text-red-500 flex items-center gap-1 transition-colors shrink-0 px-2 py-1 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
                  title="Clear category filter"
                >
                  <X size={13} />
                  <span className="hidden sm:inline">{t.ui.clear_filter ?? 'Clear'}</span>
                </button>
              </div>

              {/* Row 2: Choices BELOW the selected category */}
              {childOptions.length > 0 ? (
                <div className="flex items-center gap-1.5 overflow-x-auto no-scrollbar pt-1 border-t border-slate-150 dark:border-slate-800">
                  <span className="text-[10px] font-black uppercase tracking-widest text-slate-400 dark:text-slate-500 shrink-0 mr-1">
                    In {getCatDisplayName(activeNode)}:
                  </span>
                  {childOptions.map((child: Category) => (
                    <button
                      key={child.id}
                      onClick={() => navigateToCategory(child)}
                      className="flex items-center gap-1.5 px-3 py-1.5 rounded-full whitespace-nowrap font-bold text-xs bg-slate-100 dark:bg-slate-800 hover:bg-[#1769E0] hover:text-white dark:hover:bg-[#1769E0] dark:hover:text-white transition-all shrink-0 border border-slate-200/60 dark:border-slate-700/60 group"
                    >
                      <span className="text-sm">{getCategoryEmoji(child.slug, child.icon)}</span>
                      <span>{getCatDisplayName(child)}</span>
                      {child.children && child.children.length > 0 && (
                        <span className="text-[10px] text-slate-400 group-hover:text-white/80">›</span>
                      )}
                    </button>
                  ))}
                </div>
              ) : (
                <div className="flex items-center gap-2 py-0.5 text-xs text-slate-500 dark:text-slate-400">
                  <span className="w-2 h-2 rounded-full bg-emerald-500 inline-block animate-pulse" />
                  <span>Browsing verified listings under <strong>{getCatDisplayName(activeNode)}</strong></span>
                </div>
              )}
            </div>
          ) : (
            /* Root level: All root categories */
            <div className="flex items-center gap-1.5 overflow-x-auto no-scrollbar py-0.5">
              <button
                onClick={clearCategorySelection}
                className="flex items-center gap-1.5 px-3.5 py-1.5 rounded-full whitespace-nowrap font-bold text-xs bg-[#1769E0] text-white shadow-sm shadow-[#1769E0]/25 shrink-0"
              >
                <LayoutGrid size={13} />
                <span>{t.categories.all}</span>
              </button>

              {categories.map(cat => (
                <button
                  key={cat.id}
                  onClick={() => navigateToCategory(cat)}
                  className="flex items-center gap-1.5 px-3 py-1.5 rounded-full whitespace-nowrap font-bold text-xs bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 hover:text-[#1769E0] transition-all shrink-0 border border-slate-200/50 dark:border-slate-700/50"
                >
                  <span className="text-sm">{getCategoryEmoji(cat.slug, cat.icon)}</span>
                  <span>{getCatDisplayName(cat)}</span>
                </button>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Tier banner */}
      <TierBanner />

      {/* Popular categories strip (clean full-width, single floating AI button used instead) */}
      {!loading && (
        <div className="px-[2%] sm:px-[4%] lg:px-[6%] pt-1">
          <PopularCategoriesStrip
            categories={categories}
            onSelect={catId => navigateToCategory({ id: catId } as Category)}
          />
        </div>
      )}

      {/* ── Page body ── */}
      <div className="max-w-screen-2xl mx-auto px-[2%] sm:px-[4%] lg:px-[6%] py-4">
        {/* Products area — full-width, clean layout */}
        <div className="w-full">

          {/* ── HIGH-DEMAND ETHIOPIAN RECOMMENDATION CHIPS ── */}
          <div className="flex items-center gap-1.5 overflow-x-auto no-scrollbar pb-3 mb-1">
            <span className="text-[10px] font-black uppercase tracking-wider text-zinc-400 dark:text-zinc-500 shrink-0 mr-1 flex items-center gap-1">
              <TrendingUp size={12} className="text-[#1769E0]" />
              {locale === 'am' ? 'ተፈላጊዎች:' : 'Trending:'}
            </span>

            {/* Basic Tools & Machinery */}
            <button
              type="button"
              onClick={() => {
                const toolsCat = categories.find(c => c.slug.includes('tool') || c.slug.includes('construction') || c.slug.includes('machinery'));
                if (toolsCat) {
                  navigateToCategory(toolsCat);
                } else {
                  router.push(buildUrl({ search: 'tools' }));
                }
              }}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-amber-500/10 hover:bg-amber-500/20 text-amber-700 dark:text-amber-300 border border-amber-500/20 text-xs font-bold whitespace-nowrap shrink-0 transition-all cursor-pointer shadow-xs"
            >
              <span>🛠️</span>
              <span>{locale === 'am' ? 'መሣሪያዎችና ማሽኖች' : 'Tools & Equipment'}</span>
            </button>

            {/* Smart Phones & Electronics */}
            <button
              type="button"
              onClick={() => {
                const elecCat = categories.find(c => c.slug.includes('electronic') || c.slug.includes('phone') || c.slug.includes('computing'));
                if (elecCat) {
                  navigateToCategory(elecCat);
                } else {
                  router.push(buildUrl({ search: 'phone' }));
                }
              }}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-blue-500/10 hover:bg-blue-500/20 text-blue-700 dark:text-blue-300 border border-blue-500/20 text-xs font-bold whitespace-nowrap shrink-0 transition-all cursor-pointer shadow-xs"
            >
              <span>📱</span>
              <span>{locale === 'am' ? 'ስልኮችና ኤሌክትሮኒክስ' : 'Electronics & Phones'}</span>
            </button>

            {/* Properties & Rentals */}
            <button
              type="button"
              onClick={() => {
                setListingType('rent');
                router.push(buildUrl({ overrideType: 'rent' }));
              }}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold whitespace-nowrap shrink-0 transition-all cursor-pointer shadow-xs border ${
                listingType === 'rent'
                  ? 'bg-purple-600 text-white border-purple-600 shadow-purple-600/25'
                  : 'bg-purple-500/10 hover:bg-purple-500/20 text-purple-700 dark:text-purple-300 border-purple-500/20'
              }`}
            >
              <span>🏠</span>
              <span>{locale === 'am' ? 'ኪራይ ቤቶችና ቦታዎች' : 'Property & Rentals'}</span>
            </button>

            {/* Vehicles & Machinery */}
            <button
              type="button"
              onClick={() => {
                const vehCat = categories.find(c => c.slug.includes('vehicle') || c.slug.includes('transport') || c.slug.includes('car'));
                if (vehCat) {
                  navigateToCategory(vehCat);
                } else {
                  router.push(buildUrl({ search: 'car' }));
                }
              }}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-700 dark:text-emerald-300 border border-emerald-500/20 text-xs font-bold whitespace-nowrap shrink-0 transition-all cursor-pointer shadow-xs"
            >
              <span>🚗</span>
              <span>{locale === 'am' ? 'ተሽከርካሪዎችና ማሽኖች' : 'Vehicles & Transport'}</span>
            </button>

            {/* Agriculture & Farming */}
            <button
              type="button"
              onClick={() => {
                const agriCat = categories.find(c => c.slug.includes('agri') || c.slug.includes('farm') || c.slug.includes('food'));
                if (agriCat) {
                  navigateToCategory(agriCat);
                } else {
                  router.push(buildUrl({ search: 'agriculture' }));
                }
              }}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-lime-500/10 hover:bg-lime-500/20 text-lime-700 dark:text-lime-300 border border-lime-500/20 text-xs font-bold whitespace-nowrap shrink-0 transition-all cursor-pointer shadow-xs"
            >
              <span>🌾</span>
              <span>{locale === 'am' ? 'ግብርናና ምርቶች' : 'Agriculture & Produce'}</span>
            </button>
          </div>

          {/* ── CLEAN FILTER BAR: Results & Active Chips on Left, Sort & Filters on Right ── */}
          <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl px-4 py-3 mb-4 shadow-xs flex items-center justify-between gap-3">
            {/* Left: Results count & Active filter tags */}
            <div className="flex items-center flex-wrap gap-2 min-w-0">
              <p className="text-sm font-black text-zinc-800 dark:text-zinc-200 tracking-tight">
                {loading ? '…' : `${products.length} ${t.ui?.results_count ?? 'results'}`}
              </p>

              {/* Active filter badges with 1-click remove */}
              {listingType !== 'all' && (
                <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-orange-50 dark:bg-orange-950/40 text-orange-700 dark:text-orange-300 border border-orange-200 dark:border-orange-800 text-xs font-bold">
                  <span>{listingType === 'sale' ? `🛒 ${t.nav?.buy ?? 'Buy'}` : `🏠 ${t.nav?.rental ?? 'Rent'}`}</span>
                  <button
                    type="button"
                    onClick={() => { setListingType('all'); router.push(buildUrl({ overrideType: 'all' })); }}
                    className="hover:text-red-500 font-black ml-0.5 cursor-pointer"
                    title="Remove transaction filter"
                  >
                    ×
                  </button>
                </span>
              )}

              {activeCondition !== 'all' && (
                <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-blue-50 dark:bg-blue-950/40 text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800 text-xs font-bold">
                  <span>{activeCondition === 'new' ? (t.ui?.condition_new_badge ?? '🆕 New') : (t.ui?.condition_used_badge ?? '♻️ Used')}</span>
                  <button
                    type="button"
                    onClick={() => { setActiveCondition('all'); router.push(buildUrl({ overrideCond: 'all' })); }}
                    className="hover:text-red-500 font-black ml-0.5 cursor-pointer"
                    title="Remove condition filter"
                  >
                    ×
                  </button>
                </span>
              )}

              {(activeMinPrice !== undefined || activeMaxPrice !== undefined) && (
                <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800 text-xs font-bold">
                  <span>💰 {activeMinPrice ?? 0} - {activeMaxPrice ?? '∞'} ETB</span>
                  <button
                    type="button"
                    onClick={() => {
                      setMinPriceInput('');
                      setMaxPriceInput('');
                      setActiveMinPrice(undefined);
                      setActiveMaxPrice(undefined);
                      router.push(buildUrl({ minP: '', maxP: '' }));
                    }}
                    className="hover:text-red-500 font-black ml-0.5 cursor-pointer"
                    title="Remove price filter"
                  >
                    ×
                  </button>
                </span>
              )}

              {userCoords && (
                <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-teal-50 dark:bg-teal-950/40 text-teal-700 dark:text-teal-300 border border-teal-200 dark:border-teal-800 text-xs font-bold">
                  <span>📍 {(t.ui as any)?.near_me ?? 'Near Me'}</span>
                  <button
                    type="button"
                    onClick={() => { setUserCoords(null); router.push(buildUrl({ coords: null })); }}
                    className="hover:text-red-500 font-black ml-0.5 cursor-pointer"
                    title="Remove near me filter"
                  >
                    ×
                  </button>
                </span>
              )}

              {activeLocLabel && (
                <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-indigo-50 dark:bg-indigo-950/40 text-indigo-700 dark:text-indigo-300 border border-indigo-200 dark:border-indigo-800 text-xs font-bold">
                  <MapPin size={11} />
                  <span className="max-w-[150px] truncate">{activeLocLabel}</span>
                  <button
                    type="button"
                    onClick={() => {
                      setActiveRegion('');
                      setActiveZone('');
                      setActiveSubcity('');
                      setActiveWoreda('');
                      router.push(buildUrl({ loc: {} }));
                    }}
                    className="hover:text-red-500 font-black ml-0.5 cursor-pointer"
                    title="Remove location filter"
                  >
                    ×
                  </button>
                </span>
              )}

              {activeFilterCount > 1 && (
                <button
                  type="button"
                  onClick={resetDrawerFilters}
                  className="text-xs font-bold text-red-500 hover:underline cursor-pointer ml-1"
                >
                  {t.ui?.clear_filter ?? 'Clear all'}
                </button>
              )}
            </div>

            {/* Right Side: Near Me Find Button + Sleek Sort Selector + Filters Drawer Button */}
            <div className="flex items-center gap-2 shrink-0 flex-wrap sm:flex-nowrap">
              {/* Single "Near Me Find" Button: Automatic GPS + Manual Search */}
              <NearMeFindButton
                activeCoords={userCoords}
                activeLocationLabel={activeLocLabel || activeRegion || activeSubcity}
                activeMaxDistance={maxDistance}
                onSearch={(coords, label, dist) => {
                  setUserCoords(coords);
                  setMaxDistance(dist);
                  if (label && !coords) {
                    router.push(buildUrl({
                      search: label,
                      coords: null,
                      sort: sortBy,
                      maxDist: dist,
                    }));
                  } else {
                    router.push(buildUrl({
                      coords: coords || null,
                      sort: coords ? 'nearest' : sortBy,
                      maxDist: dist,
                    }));
                  }
                }}
              />

              {/* Sort Selector Dropdown */}
              <div className="relative">
                <select
                  value={sortBy}
                  onChange={(e) => {
                    const newSort = e.target.value as any;
                    setSortBy(newSort);
                    router.push(buildUrl({ sort: newSort }));
                  }}
                  className="appearance-none bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl pl-3 pr-7 py-2 text-xs font-extrabold text-zinc-800 dark:text-zinc-200 outline-none focus:ring-2 focus:ring-[#1769E0] cursor-pointer shadow-xs"
                >
                  <option value="trending">🔥 {locale === 'am' ? 'ተፈላጊ / ምርጥ' : 'Trending'}</option>
                  <option value="newest">🆕 {locale === 'am' ? 'አዳዲስ' : 'Newest'}</option>
                  <option value="nearest">📍 {locale === 'am' ? 'በአቅራቢያዎ' : 'Nearest GPS'}</option>
                  <option value="views">👁️ {locale === 'am' ? 'በጣም የታዩ' : 'Most Viewed'}</option>
                  <option value="price_asc">💰 {locale === 'am' ? 'ዋጋ፡ ዝቅተኛ' : 'Price: Low'}</option>
                  <option value="price_desc">💎 {locale === 'am' ? 'ዋጋ፡ ከፍተኛ' : 'Price: High'}</option>
                </select>
                <ChevronDown size={11} className="absolute right-2 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none" />
              </div>

              {/* Filters Button */}
              <button
                type="button"
                onClick={openFilterDrawer}
                className={`shrink-0 flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-extrabold transition-all cursor-pointer shadow-xs ${activeFilterCount > 0
                  ? 'bg-[#1769E0] text-white hover:bg-[#1055BD] shadow-md shadow-[#1769E0]/25'
                  : 'bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 text-zinc-800 dark:text-zinc-200 hover:border-[#1769E0] hover:text-[#1769E0]'
                  }`}
                aria-label="Open filters"
              >
                <SlidersHorizontal size={14} className={activeFilterCount > 0 ? 'text-white' : 'text-[#1769E0]'} />
                <span>{t.ui?.filters ?? 'Filters'}</span>
                {activeFilterCount > 0 && (
                  <span className="w-5 h-5 rounded-full bg-white text-[#1769E0] text-[10px] font-black flex items-center justify-center">
                    {activeFilterCount}
                  </span>
                )}
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
                <h3 className="text-base font-bold text-slate-900 dark:text-slate-100">
                  {(t.ui as any).unable_to_load_listings ?? 'Unable to Load Listings'}
                </h3>
                <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">{fetchError}</p>
              </div>
              <button
                onClick={() => setRetryKey(k => k + 1)}
                className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-[#1769E0] hover:bg-[#1055BD] text-white text-xs font-bold shadow-md transition-all active:scale-95 cursor-pointer"
              >
                <RotateCcw size={14} />
                <span>{(t.ui as any).retry ?? 'Retry'}</span>
              </button>
            </div>
          ) : products.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-16 px-4 text-center max-w-md mx-auto space-y-4">
              <div className="w-16 h-16 rounded-3xl bg-slate-100 dark:bg-slate-800 text-slate-400 flex items-center justify-center">
                <ShoppingBag size={32} />
              </div>
              <div>
                <h3 className="text-base font-black text-slate-900 dark:text-slate-100">{t.ui.no_products ?? 'No listings found'}</h3>
                <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">
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
                  className="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-[#1769E0] hover:bg-[#1055BD] text-white text-xs font-bold transition-all shadow-md active:scale-95 cursor-pointer"
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
                      className="px-2.5 py-1 rounded-lg bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:text-[#1769E0] text-xs font-medium transition-colors"
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
                    <div className="w-10 h-10 rounded-2xl bg-[#1769E0]/15 text-[#1769E0] dark:text-[#38A5FF] flex items-center justify-center group-hover:bg-[#1769E0] group-hover:text-white transition-all shadow-md">
                      <ChevronDown size={20} />
                    </div>
                  </button>
                </div>
              )}
            </>
          )}
        </div>
      </div>

      {/* ── Slide-Over Filter Drawer ── */}
      {filterDrawerOpen && (
        <div className="fixed inset-0 z-50 flex justify-end">
          {/* Backdrop */}
          <div
            className="fixed inset-0 bg-black/50 backdrop-blur-xs transition-opacity animate-fade-in"
            onClick={() => setFilterDrawerOpen(false)}
          />

          {/* Side Drawer Panel */}
          <div className="relative w-full max-w-md bg-white dark:bg-zinc-900 shadow-2xl h-full flex flex-col z-10 border-l border-zinc-200 dark:border-zinc-800 animate-slide-in-right overflow-hidden">
            {/* Drawer Header */}
            <div className="px-5 py-4 border-b border-zinc-200 dark:border-zinc-800 flex items-center justify-between">
              <div className="flex items-center gap-2">
                <SlidersHorizontal size={18} className="text-[#1769E0]" />
                <h2 className="text-base font-black text-zinc-900 dark:text-white">
                  {t.ui?.filters ?? 'Filters'}
                </h2>
                {activeFilterCount > 0 && (
                  <span className="px-2 py-0.5 rounded-full bg-[#1769E0]/10 text-[#1769E0] text-xs font-black">
                    {activeFilterCount} active
                  </span>
                )}
              </div>
              <button
                type="button"
                onClick={() => setFilterDrawerOpen(false)}
                className="p-2 rounded-xl text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors cursor-pointer"
                aria-label="Close filters"
              >
                <X size={18} />
              </button>
            </div>

            {/* Drawer Body (Scrollable) */}
            <div className="flex-1 overflow-y-auto p-5 space-y-5">
              {/* Filter 1: Marketplace Mode */}
              <div className="space-y-1.5">
                <label className="text-[11px] font-black uppercase tracking-wider text-zinc-500">
                  {(t.ui as any)?.marketplace_pillar ?? 'Marketplace Mode'}
                </label>
                <div className="grid grid-cols-3 gap-1.5 bg-zinc-100 dark:bg-zinc-800 p-1 rounded-xl">
                  <button
                    type="button"
                    onClick={() => setDrawerType('all')}
                    className={`py-2 text-xs font-black uppercase rounded-lg transition-all cursor-pointer ${drawerType === 'all'
                      ? 'bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white shadow-xs'
                      : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'
                      }`}
                  >
                    {(t.ui as any)?.all_modes ?? 'All'}
                  </button>
                  <button
                    type="button"
                    onClick={() => setDrawerType('sale')}
                    className={`py-2 text-xs font-black uppercase rounded-lg transition-all flex items-center justify-center gap-1 cursor-pointer ${drawerType === 'sale'
                      ? 'bg-orange-600 text-white shadow-xs'
                      : 'text-zinc-500 hover:text-orange-600'
                      }`}
                  >
                    🛒 {t.nav?.buy ?? 'Buy'}
                  </button>
                  <button
                    type="button"
                    onClick={() => setDrawerType('rent')}
                    className={`py-2 text-xs font-black uppercase rounded-lg transition-all flex items-center justify-center gap-1 cursor-pointer ${drawerType === 'rent'
                      ? 'bg-teal-600 text-white shadow-xs'
                      : 'text-zinc-500 hover:text-teal-600'
                      }`}
                  >
                    🏠 {t.nav?.rental ?? 'Rent'}
                  </button>
                </div>
              </div>

              {/* Filter 2: Condition */}
              <div className="space-y-1.5">
                <label className="text-[11px] font-black uppercase tracking-wider text-zinc-500">
                  {t.ui?.condition ?? 'Condition'}
                </label>
                <div className="grid grid-cols-3 gap-1.5 bg-zinc-100 dark:bg-zinc-800 p-1 rounded-xl">
                  {(['all', 'new', 'used'] as const).map(c => (
                    <button
                      key={c}
                      type="button"
                      onClick={() => setDrawerCondition(c)}
                      className={`py-2 text-xs font-black uppercase rounded-lg transition-all cursor-pointer ${drawerCondition === c
                        ? c === 'new'
                          ? 'bg-emerald-600 text-white shadow-xs'
                          : c === 'used'
                            ? 'bg-blue-600 text-white shadow-xs'
                            : 'bg-zinc-800 dark:bg-zinc-200 text-white dark:text-zinc-900 shadow-xs'
                        : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'
                        }`}
                    >
                      {c === 'all' ? (t.ui?.condition_all ?? 'All') : c === 'new' ? (t.ui?.condition_new_badge ?? '🆕 New') : (t.ui?.condition_used_badge ?? '♻️ Used')}
                    </button>
                  ))}
                </div>
              </div>

              {/* Filter 3: Price Range */}
              <div className="space-y-1.5">
                <label className="text-[11px] font-black uppercase tracking-wider text-zinc-500">
                  {t.ui?.price ?? 'Price'} (ETB)
                </label>
                <div className="flex items-center gap-2">
                  <input
                    type="number"
                    placeholder={t.ui?.min_price ?? 'Min ETB'}
                    value={drawerMinPrice}
                    onChange={e => setDrawerMinPrice(e.target.value)}
                    className="w-1/2 px-3 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl outline-none focus:ring-2 focus:ring-[#1769E0]"
                    min={0}
                  />
                  <span className="text-zinc-400 font-bold">-</span>
                  <input
                    type="number"
                    placeholder={t.ui?.max_price ?? 'Max ETB'}
                    value={drawerMaxPrice}
                    onChange={e => setDrawerMaxPrice(e.target.value)}
                    className="w-1/2 px-3 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl outline-none focus:ring-2 focus:ring-[#1769E0]"
                    min={0}
                  />
                </div>
              </div>

              {/* Filter 4: Proximity / Near Me */}
              <div className="space-y-1.5">
                <label className="text-[11px] font-black uppercase tracking-wider text-zinc-500">
                  {(t.ui as any)?.distance_label ?? 'Proximity'}
                </label>
                <button
                  type="button"
                  onClick={detectNearMe}
                  disabled={isLocating}
                  className={`w-full py-2.5 px-3 rounded-xl border text-xs font-bold flex items-center justify-center gap-2 transition-all cursor-pointer ${userCoords
                    ? 'bg-teal-50 dark:bg-teal-950/40 border-teal-500 text-teal-600 dark:text-teal-400 ring-1 ring-teal-500/30'
                    : 'bg-zinc-50 dark:bg-zinc-800 border-zinc-200 dark:border-zinc-700 text-zinc-700 dark:text-zinc-300 hover:border-teal-500'
                    }`}
                >
                  <Navigation size={14} className={isLocating ? 'animate-spin text-teal-500' : userCoords ? 'text-teal-600 fill-teal-600' : 'text-zinc-400'} />
                  <span>
                    {isLocating
                      ? ((t.ui as any).locating ?? 'Locating…')
                      : userCoords
                        ? ((t.ui as any).nearest_to_me ?? '📍 Nearest to Me (Active)')
                        : ((t.ui as any).near_me ?? 'Find Listings Near Me')}
                  </span>
                  {userCoords && (
                    <span
                      onClick={(e) => {
                        e.stopPropagation();
                        setUserCoords(null);
                        router.push(buildUrl({ coords: null }));
                      }}
                      className="ml-2 hover:text-red-500 font-black cursor-pointer text-sm"
                    >
                      ×
                    </span>
                  )}
                </button>
              </div>

              {/* Filter 5: Location Selector */}
              <div className="space-y-1.5">
                <label className="text-[11px] font-black uppercase tracking-wider text-zinc-500">
                  {t.ui?.location_label ?? 'Location (Ethiopia)'}
                </label>
                <div className="p-3 bg-zinc-50 dark:bg-zinc-800/60 rounded-xl border border-zinc-200 dark:border-zinc-700">
                  <LocationSelector
                    showPrecise={false}
                    defaultValue={{ region: drawerRegion, zone: drawerZone, subcity: drawerSubcity, woreda: drawerWoreda }}
                    onChange={(val) => {
                      setDrawerRegion(val.region || '');
                      setDrawerZone(val.zone || '');
                      setDrawerSubcity(val.subcity || '');
                      setDrawerWoreda(val.woreda || '');
                    }}
                  />
                </div>
              </div>
            </div>

            {/* Drawer Footer Actions */}
            <div className="p-4 border-t border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900/80 flex items-center gap-3">
              <button
                type="button"
                onClick={resetDrawerFilters}
                className="flex-1 py-2.5 rounded-xl border border-zinc-200 dark:border-zinc-700 text-xs font-bold text-zinc-600 dark:text-zinc-400 hover:bg-white dark:hover:bg-zinc-800 transition-colors flex items-center justify-center gap-1.5 cursor-pointer"
              >
                <RotateCcw size={13} />
                <span>{t.ui?.clear_filter ?? 'Reset All'}</span>
              </button>

              <button
                type="button"
                onClick={applyDrawerFilters}
                className="flex-1 py-2.5 rounded-xl bg-[#1769E0] text-white text-xs font-extrabold hover:bg-[#1055BD] shadow-md shadow-[#1769E0]/20 transition-all flex items-center justify-center gap-1.5 cursor-pointer"
              >
                <span>{t.ui?.submit ?? 'Apply Filters'}</span>
              </button>
            </div>
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
