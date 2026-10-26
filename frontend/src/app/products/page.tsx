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
import { useAuth } from '@/context/AuthContext';
import {
  ShoppingBag, ChevronDown, SlidersHorizontal, X,
  LayoutGrid, ChevronRight, TrendingUp, MapPin, Star,
} from 'lucide-react';

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

  const catLabel = (slug: string) => (t.categories as any)[slug] ?? slug.replace(/_/g, ' ');

  const EMOJI: Record<string, string> = {
    electronics: '📱', vehicles: '🚗', property: '🏠', clothing: '👗',
    agriculture: '🌾', furniture: '🛋️', food: '🍎', business: '💼',
    industry: '🏭', construction: '🔨', health_beauty: '💄', sports: '⚽',
    education: '📚', entertainment: '🎮', luxury: '💎', animals: '🐄',
    services: '🔧', jobs: '💼', other: '📦',
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
    if (woreda) { f.woreda = woreda; setActiveWoreda(woreda); setActiveSubcity(''); setActiveZone(''); setActiveRegion(''); }
    else if (subcity) { f.subcity = subcity; setActiveSubcity(subcity); setActiveWoreda(''); setActiveZone(''); setActiveRegion(''); }
    else if (zone) { f.zone = zone; setActiveZone(zone); setActiveSubcity(''); setActiveWoreda(''); setActiveRegion(''); }
    else if (region) { f.region = region; setActiveRegion(region); setActiveZone(''); setActiveSubcity(''); setActiveWoreda(''); }
    else { setActiveRegion(''); setActiveZone(''); setActiveSubcity(''); setActiveWoreda(''); }

    setFilters(f);
  }, [searchParams, categories]);

  useEffect(() => {
    setLoading(true);
    setVisibleCount(24);
    const apiFilters: ProductFilters = { ...filters };
    if (activeCondition !== 'all') apiFilters.condition = activeCondition;
    if (listingType !== 'all') apiFilters.listing_type = listingType;
    api.getProducts(apiFilters)
      .then(data => setProducts(data))
      .catch(() => setProducts([]))
      .finally(() => setLoading(false));
  }, [filters, activeCondition, listingType]);

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

  const activeLocLabel = activeWoreda || activeSubcity || activeZone || activeRegion;

  const buildUrl = (catId?: number, sub?: string, item?: string, search?: string, loc?: { region?: string; zone?: string; subcity?: string; woreda?: string }, overrideType?: 'all' | 'sale' | 'rent') => {
    const p = new URLSearchParams();
    const typeToUse = overrideType !== undefined ? overrideType : listingType;
    if (typeToUse && typeToUse !== 'all') p.set('type', typeToUse);
    if (catId) {
      p.set('category_id', String(catId));
      p.set('category_hierarchy_id', String(catId));
    }
    if (sub) p.set('subcategory', sub);
    if (item) p.set('item', item);
    if (search) p.set('search', search);
    const l = loc ?? { region: activeRegion || undefined, zone: activeZone || undefined, subcity: activeSubcity || undefined, woreda: activeWoreda || undefined };
    if (l.woreda) p.set('woreda', l.woreda);
    else if (l.subcity) p.set('subcity', l.subcity);
    else if (l.zone) p.set('zone', l.zone);
    else if (l.region) p.set('region', l.region);
    return `/products${p.toString() ? '?' + p.toString() : ''}`;
  };

  const catLabel = (slug: string) => (t.categories as any)[slug] ?? slug.replace(/_/g, ' ');
  const itemLabel = (key: string) => {
    // For subcategory keys: check sub_categories[activeCatSlug] first
    const fromSub = activeCatSlug ? (t.sub_categories as any)?.[activeCatSlug]?.[key] : undefined;
    if (fromSub) return fromSub;
    const fromItems = (t.items as any)?.[key];
    if (fromItems) return fromItems;
    return key.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase());
  };

  return (
    <>
      {/* ── Sticky nav ── */}
      <div className="sticky top-14 z-40 bg-white/95 dark:bg-zinc-950/95 backdrop-blur border-b border-zinc-100 dark:border-zinc-800">

        {/* Category pills */}
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
              {catLabel(cat.slug)}
            </button>
          ))}
        </div>

        {/* Dimension 1: Marketplace Pillar (Transaction Mode) */}
        <div className="flex items-center gap-1.5 px-[2%] sm:px-[4%] lg:px-[6%] py-1.5 border-t border-zinc-100 dark:border-zinc-800 bg-gradient-to-r from-orange-50/50 via-teal-50/30 to-transparent dark:from-zinc-900/60">
          <span className="text-[9px] font-black uppercase tracking-widest text-zinc-400 shrink-0 mr-1">{t.ui.listing_type_label ?? 'Marketplace Pillar:'}</span>
          {(['all', 'sale', 'rent'] as const).map(type => (
            <button key={type} onClick={() => { 
              setListingType(type); 
              setVisibleCount(24);
              router.push(buildUrl(activeCatId, activeSub, activeItem, filters.search, undefined, type));
            }}
              className={`px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-wide transition-all shrink-0 flex items-center gap-1 ${listingType === type
                ? type === 'sale' ? 'bg-orange-600 text-white shadow-md shadow-orange-600/20'
                  : type === 'rent' ? 'bg-teal-600 text-white shadow-md shadow-teal-600/20'
                    : 'bg-zinc-800 dark:bg-zinc-100 text-white dark:text-zinc-900 shadow-md'
                : 'bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'
                }`}>
              {type === 'all' ? (t.ui.listing_type_all ?? '🌐 All Listings') 
                : type === 'sale' ? (t.ui.listing_type_sale ?? '🛒 Buy (Seller Goods)') 
                : (t.ui.listing_type_rent ?? '🏠 Rent (Landlord Rentals)')}
            </button>
          ))}
        </div>

        {/* Dimension 2: Condition filter — Physical state of items (orthogonal to Pillar & Category) */}
        <div className="flex items-center gap-1.5 px-[2%] sm:px-[4%] lg:px-[6%] pb-2 pt-1 border-t border-zinc-100 dark:border-zinc-850">
          <span className="text-[9px] font-black uppercase tracking-widest text-zinc-400 shrink-0 mr-1">{t.ui.condition_label ?? 'Physical Condition:'}</span>
          {(['all', 'new', 'used'] as const).map(c => (
            <button key={c} onClick={() => { setActiveCondition(c); setVisibleCount(24); }}
              className={`px-3 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wide transition-all shrink-0 ${activeCondition === c
                ? c === 'new' ? 'bg-green-600 text-white'
                  : c === 'used' ? 'bg-blue-600 text-white'
                    : 'bg-zinc-700 dark:bg-zinc-200 text-white dark:text-zinc-900'
                : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'
                }`}>
              {c === 'all' ? (t.ui.condition_all ?? 'All') : c === 'new' ? (t.ui.condition_new_badge ?? '🆕 Brand New') : (t.ui.condition_used_badge ?? '♻️ Used / Pre-owned')}
            </button>
          ))}
        </div>

        {/* Subcategory pills */}
        {subcategoryNodes.length > 0 && (
          <div className="border-t border-zinc-100 dark:border-zinc-800/60 bg-zinc-50/80 dark:bg-zinc-900/60">
            <div className="flex items-center gap-1.5 overflow-x-auto no-scrollbar py-1.5 px-[2%] sm:px-[4%] lg:px-[6%]">
              <ChevronRight size={10} className="text-orange-500 shrink-0 opacity-60" />
              {subcategoryNodes.map(s => {
                const sSlug = s.slug;
                const sLabel = s.name || s.name_en || sSlug;
                return (
                  <button key={s.id || sSlug}
                    onClick={() => { setActiveSub(activeSub === sSlug ? undefined : sSlug); setActiveItem(undefined); router.push(buildUrl(activeCatId, activeSub === sSlug ? undefined : sSlug)); }}
                    className={`px-2.5 py-1 rounded-full whitespace-nowrap font-bold text-[10px] transition-all shrink-0 ${activeSub === sSlug ? 'bg-zinc-800 dark:bg-zinc-100 text-white dark:text-zinc-900' : 'bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 hover:border-zinc-400'
                      }`}>
                    {sLabel}
                  </button>
                );
              })}
            </div>
          </div>
        )}

        {/* Item pills */}
        {itemNodes.length > 0 && (
          <div className="border-t border-zinc-100 dark:border-zinc-800/60 bg-orange-50/50 dark:bg-orange-950/20">
            <div className="flex items-center gap-1.5 overflow-x-auto no-scrollbar py-1.5 px-[2%] sm:px-[4%] lg:px-[6%]">
              <ChevronRight size={10} className="text-orange-500 shrink-0 opacity-60" />
              <ChevronRight size={10} className="text-orange-500 shrink-0 opacity-40 -ml-2.5" />
              {itemNodes.map(i => {
                const iSlug = i.slug;
                const iLabel = i.name || i.name_en || iSlug;
                return (
                  <button key={i.id || iSlug}
                    onClick={() => { setActiveItem(activeItem === iSlug ? undefined : iSlug); router.push(buildUrl(activeCatId, activeSub, activeItem === iSlug ? undefined : iSlug)); }}
                    className={`px-2 py-0.5 rounded-full whitespace-nowrap font-bold text-[9px] transition-all shrink-0 ${activeItem === iSlug ? 'bg-orange-600 text-white shadow-sm' : 'bg-white dark:bg-zinc-800 border border-orange-200 dark:border-orange-900 text-zinc-600 dark:text-zinc-400 hover:border-orange-400'
                      }`}>
                    {iLabel}
                  </button>
                );
              })}
            </div>
          </div>
        )}
      </div>

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

      {/* Popular categories strip */}
      {!loading && (
        <PopularCategoriesStrip
          categories={categories}
          onSelect={catId => { setActiveCatId(catId); setActiveSub(undefined); setActiveItem(undefined); router.push(buildUrl(catId)); }}
        />
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

            {/* Top bar */}
            <div className="flex items-center justify-between mb-3">
              <p className="text-xs font-bold text-zinc-400">
                {loading ? '…' : `${products.length} ${t.ui.results_count ?? 'results'}`}
              </p>
              <button onClick={() => setFilterOpen(true)}
                className="lg:hidden flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-zinc-100 dark:bg-zinc-800 font-bold text-xs hover:bg-orange-600 hover:text-white transition-all">
                <SlidersHorizontal size={12} /> {t.ui.filters}
              </button>
            </div>

            {loading ? (
              <div className="flex items-center justify-center py-20"><Spinner /></div>
            ) : products.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-20 opacity-40">
                <ShoppingBag size={48} className="mb-3" />
                <p className="font-bold">{t.ui.no_products}</p>
              </div>
            ) : (
              <>
                {/*
                  Grid: 3 cols on mobile (≥360px), 4 on sm (≥480px),
                        5 on md (≥768px), 6 on lg+ (≥1024px)
                  With sidebar: effectively 5 on lg, 6 on xl
                */}
                <div className="grid grid-cols-3 min-[480px]:grid-cols-4 md:grid-cols-5 lg:grid-cols-5 xl:grid-cols-6 gap-2">
                  {products.slice(0, visibleCount).map(p => (
                    p.listing_type === 'rent'
                      ? <RentalProductCard key={p.id} product={p} />
                      : <ProductCard key={p.id} product={p} />
                  ))}
                </div>

                {products.length > visibleCount && (
                  <div className="flex justify-center mt-8">
                    <button onClick={() => setVisibleCount(v => v + 24)}
                      className="flex flex-col items-center gap-2 group">
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
          <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={() => setFilterOpen(false)} />
          <div className="absolute bottom-0 left-0 right-0 max-h-[88vh] overflow-y-auto rounded-t-3xl bg-zinc-50 dark:bg-zinc-950 p-4">
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
