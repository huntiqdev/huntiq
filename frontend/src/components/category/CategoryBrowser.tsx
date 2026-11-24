'use client';

import React, { useState, useEffect, useCallback, useRef } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import {
  ChevronRight,
  ChevronLeft,
  Search,
  Grid,
  Layers,
  Sparkles,
  ArrowRight,
  FolderOpen,
  Folder,
  Tag,
  X,
  Compass
} from 'lucide-react';
import { useLocale } from '@/hooks/useLocale';
import { api } from '@/lib/api';

export interface CategoryItem {
  id: number;
  parent_id?: number | null;
  name: string;
  name_en?: string;
  name_am?: string;
  name_om?: string;
  name_ti?: string;
  slug: string;
  level: number;
  is_leaf: boolean;
  display_order?: number;
  icon?: string;
  image_url?: string;
  path?: string;
  classification_system_code?: string;
}

interface CategoryBrowserProps {
  className?: string;
  onSelectCategory?: (category: CategoryItem) => void;
  standalone?: boolean;
}

const CLASSIFICATION_SYSTEMS = [
  { code: 'ALL', labelEn: 'All Categories', labelAm: 'ሁሉም ምድቦች', icon: '🌐' },
  { code: 'SYS_PROD', labelEn: 'Products & Goods', labelAm: 'ሸቀጦች እና ምርቶች', icon: '🛍️' },
  { code: 'SYS_SRV', labelEn: 'Services', labelAm: 'አገልግሎቶች', icon: '🛠️' },
  { code: 'SYS_PRM', labelEn: 'Properties & Real Estate', labelAm: 'ንብረት እና ቤቶች', icon: '🏠' },
  { code: 'SYS_EMP', labelEn: 'Jobs & Careers', labelAm: 'ስራዎች', icon: '💼' },
  { code: 'SYS_GIG', labelEn: 'Freelance & Gigs', labelAm: 'ፍሪላንስ', icon: '⚡' },
  { code: 'SYS_MED', labelEn: 'Healthcare & Medical', labelAm: 'ጤና እና ህክምና', icon: '🏥' },
  { code: 'SYS_OPP', labelEn: 'Tenders & Business', labelAm: 'ጨረታዎች', icon: '📑' },
];

export function CategoryBrowser({
  className = '',
  onSelectCategory,
  standalone = true,
}: CategoryBrowserProps) {
  const { locale, t } = useLocale();
  const router = useRouter();

  // State
  const [selectedSystem, setSelectedSystem] = useState<string>('ALL');
  const [breadcrumbs, setBreadcrumbs] = useState<CategoryItem[]>([]);
  const [currentChildren, setCurrentChildren] = useState<CategoryItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [searching, setSearching] = useState<boolean>(false);

  // Helper to extract localized name
  const getCategoryName = useCallback((cat: CategoryItem | any): string => {
    if (!cat) return '';
    switch (locale) {
      case 'am':
        return cat.name_am || cat.name || cat.name_en || cat.slug;
      case 'om':
      case 'or':
        return cat.name_om || cat.name || cat.name_en || cat.slug;
      case 'ti':
        return cat.name_ti || cat.name || cat.name_en || cat.slug;
      default:
        return cat.name_en || cat.name || cat.slug;
    }
  }, [locale]);

  // Load root categories
  const loadRoots = useCallback(async () => {
    setLoading(true);
    try {
      const res = await fetch(`/api/categories/roots?lang=${encodeURIComponent(locale)}`);
      const data = await res.json();
      if (data.success && Array.isArray(data.data)) {
        setCurrentChildren(data.data);
      } else {
        // Fallback to general categories API
        const fallback = await api.getCategories();
        setCurrentChildren(fallback || []);
      }
    } catch (err) {
      console.error('Failed to load root categories:', err);
      setCurrentChildren([]);
    } finally {
      setLoading(false);
    }
  }, [locale]);

  // Load children for current drill-down target
  const loadChildren = useCallback(async (parentId: number) => {
    setLoading(true);
    try {
      const res = await fetch(`/api/categories/${parentId}/children?lang=${encodeURIComponent(locale)}`);
      const data = await res.json();
      if (data.success && Array.isArray(data.data)) {
        setCurrentChildren(data.data);
      } else {
        setCurrentChildren([]);
      }
    } catch (err) {
      console.error(`Failed to load children for category ${parentId}:`, err);
      setCurrentChildren([]);
    } finally {
      setLoading(false);
    }
  }, [locale]);

  // Initial load
  useEffect(() => {
    if (breadcrumbs.length === 0) {
      loadRoots();
    } else {
      const current = breadcrumbs[breadcrumbs.length - 1];
      loadChildren(current.id);
    }
  }, [locale, breadcrumbs, loadRoots, loadChildren]);

  // Debounced live hierarchical search
  useEffect(() => {
    const q = searchQuery.trim();
    if (q.length < 2) {
      setSearchResults([]);
      setSearching(false);
      return;
    }

    setSearching(true);
    const timer = setTimeout(async () => {
      try {
        const res = await fetch(`/api/categories/search?q=${encodeURIComponent(q)}&lang=${encodeURIComponent(locale)}&limit=30`);
        const data = await res.json();
        if (data.success && Array.isArray(data.data)) {
          setSearchResults(data.data);
        } else {
          setSearchResults([]);
        }
      } catch (err) {
        console.error('Category search failed:', err);
        setSearchResults([]);
      } finally {
        setSearching(false);
      }
    }, 200);

    return () => clearTimeout(timer);
  }, [searchQuery, locale]);

  // User clicked on a category node
  const handleCategoryClick = (cat: CategoryItem) => {
    if (cat.is_leaf) {
      if (onSelectCategory) {
        onSelectCategory(cat);
      } else {
        router.push(`/products?category_hierarchy_id=${cat.id}&cat=${cat.slug}`);
      }
      return;
    }

    // Drill down into children
    setBreadcrumbs((prev) => [...prev, cat]);
  };

  // Navigate back to a specific breadcrumb level
  const handleBreadcrumbClick = (index: number) => {
    if (index === -1) {
      // Return to roots
      setBreadcrumbs([]);
      return;
    }
    setBreadcrumbs((prev) => prev.slice(0, index + 1));
  };

  // Step up one level
  const handleStepUp = () => {
    if (breadcrumbs.length === 0) return;
    setBreadcrumbs((prev) => prev.slice(0, prev.length - 1));
  };

  // Filter children by classification system code if selected
  const filteredChildren = currentChildren.filter((cat) => {
    if (selectedSystem === 'ALL') return true;
    return !cat.classification_system_code || cat.classification_system_code === selectedSystem;
  });

  const activeParent = breadcrumbs.length > 0 ? breadcrumbs[breadcrumbs.length - 1] : null;

  return (
    <div className={`w-full rounded-3xl bg-white dark:bg-zinc-900 border border-zinc-200/80 dark:border-zinc-800 shadow-xl overflow-hidden ${className}`}>
      {/* Top Header & Search Bar */}
      <div className="p-5 sm:p-7 border-b border-zinc-150 dark:border-zinc-800 bg-gradient-to-r from-blue-50/50 via-indigo-50/30 to-purple-50/50 dark:from-zinc-900 dark:via-zinc-900 dark:to-zinc-900">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div className="flex items-center gap-2 mb-1.5">
              <span className="flex items-center justify-center w-7 h-7 rounded-xl bg-[#1769E0] text-white text-xs shadow-md shadow-[#1769E0]/20">
                <Compass size={16} />
              </span>
              <h3 className="text-lg sm:text-xl font-black text-zinc-900 dark:text-white tracking-tight">
                {locale === 'am' ? 'የምድብ አሳሽ (20,500+ ምድቦች)' : 'Marketplace Taxonomy Explorer (20,500+ Categories)'}
              </h3>
            </div>
            <p className="text-xs sm:text-sm text-zinc-500 dark:text-zinc-400">
              {locale === 'am'
                ? 'ምርቶችን፣ አገልግሎቶችን እና ንብረቶችን በደረጃ ይፈልጉ እና ያስሱ'
                : 'Drill down through deep canonical hierarchy across all 7 national classification systems'}
            </p>
          </div>

          {/* Search Input */}
          <div className="relative w-full md:w-80">
            <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder={locale === 'am' ? 'ምድብ ፈልግ (ምሳሌ፦ ሳምሰንግ)...' : 'Search any category (e.g. Samsung)...'}
              className="w-full pl-10 pr-9 py-2 text-xs sm:text-sm rounded-2xl bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 focus:outline-none focus:ring-2 focus:ring-[#1769E0] dark:text-white transition-all shadow-sm"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery('')}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200"
              >
                <X size={14} />
              </button>
            )}
          </div>
        </div>

        {/* Classification Systems Filter Pills (only at root level) */}
        {breadcrumbs.length === 0 && !searchQuery && (
          <div className="flex items-center gap-2 overflow-x-auto no-scrollbar pt-4 mt-2 border-t border-zinc-200/50 dark:border-zinc-800/50">
            {CLASSIFICATION_SYSTEMS.map((sys) => {
              const active = selectedSystem === sys.code;
              return (
                <button
                  key={sys.code}
                  onClick={() => setSelectedSystem(sys.code)}
                  className={`flex items-center gap-1.5 px-3 py-1.5 rounded-xl whitespace-nowrap text-xs font-bold transition-all shrink-0 ${
                    active
                      ? 'bg-[#1769E0] text-white shadow-md shadow-[#1769E0]/25'
                      : 'bg-white/80 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-700 border border-zinc-200/60 dark:border-zinc-700/60'
                  }`}
                >
                  <span>{sys.icon}</span>
                  <span>{locale === 'am' ? sys.labelAm : sys.labelEn}</span>
                </button>
              );
            })}
          </div>
        )}

        {/* Interactive Breadcrumb Navigation Trail */}
        {breadcrumbs.length > 0 && !searchQuery && (
          <div className="flex items-center justify-between gap-3 pt-3 mt-3 border-t border-zinc-200/60 dark:border-zinc-800">
            <div className="flex items-center gap-1.5 overflow-x-auto no-scrollbar py-1">
              <button
                onClick={() => handleBreadcrumbClick(-1)}
                className="text-xs font-bold text-zinc-500 hover:text-[#1769E0] dark:text-zinc-400 dark:hover:text-[#38A5FF] transition-colors flex items-center gap-1 shrink-0"
              >
                <Grid size={13} />
                <span>{locale === 'am' ? 'ዋና ምድቦች' : 'All Roots'}</span>
              </button>

              {breadcrumbs.map((crumb, idx) => {
                const isLast = idx === breadcrumbs.length - 1;
                return (
                  <div key={crumb.id} className="flex items-center gap-1.5 shrink-0">
                    <ChevronRight size={13} className="text-zinc-300 dark:text-zinc-600" />
                    {isLast ? (
                      <span className="text-xs font-black text-[#1769E0] dark:text-[#38A5FF] bg-blue-50 dark:bg-blue-900/30 px-2 py-0.5 rounded-lg border border-blue-100 dark:border-blue-800">
                        {getCategoryName(crumb)}
                      </span>
                    ) : (
                      <button
                        onClick={() => handleBreadcrumbClick(idx)}
                        className="text-xs font-semibold text-zinc-600 hover:text-[#1769E0] dark:text-zinc-300 dark:hover:text-[#38A5FF] transition-colors"
                      >
                        {getCategoryName(crumb)}
                      </button>
                    )}
                  </div>
                );
              })}
            </div>

            {/* Back Up Button */}
            <button
              onClick={handleStepUp}
              className="flex items-center gap-1 px-2.5 py-1 text-xs font-bold text-zinc-600 dark:text-zinc-300 bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded-lg transition-colors shrink-0"
              title="Step Up One Level"
            >
              <ChevronLeft size={14} />
              <span>{locale === 'am' ? 'ተመለስ' : 'Back'}</span>
            </button>
          </div>
        )}
      </div>

      {/* Main Content Area */}
      <div className="p-5 sm:p-7 min-h-[300px]">
        {/* Search Results Display */}
        {searchQuery ? (
          <div>
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-bold text-zinc-400 uppercase tracking-wider">
                {searching ? (locale === 'am' ? 'በመፈለግ ላይ...' : 'Searching canonical taxonomy...') : `${searchResults.length} matching categories found`}
              </span>
              <button
                onClick={() => setSearchQuery('')}
                className="text-xs font-bold text-[#1769E0] hover:underline"
              >
                {locale === 'am' ? 'ፍለጋ አጽዳ' : 'Clear Search'}
              </button>
            </div>

            {searchResults.length > 0 ? (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                {searchResults.map((result: any) => {
                  const bcs = Array.isArray(result.breadcrumbs) ? result.breadcrumbs : [];
                  const pathString = bcs.map((b: any) => b.name).join(' › ') || result.hierarchy;

                  return (
                    <Link
                      key={result.id}
                      href={`/products?category_hierarchy_id=${result.id}&cat=${result.slug}`}
                      className="group flex flex-col justify-between p-3.5 rounded-2xl border border-zinc-200/70 dark:border-zinc-800 bg-white dark:bg-zinc-800/40 hover:border-[#1769E0] dark:hover:border-[#1769E0] hover:shadow-md transition-all"
                    >
                      <div>
                        <div className="flex items-center justify-between gap-2 mb-1.5">
                          <span className="text-xs font-black text-zinc-900 dark:text-white group-hover:text-[#1769E0] transition-colors">
                            {result.name || result.name_en}
                          </span>
                          <span className="text-[10px] font-bold px-2 py-0.5 rounded-full bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 border border-blue-100 dark:border-blue-800">
                            Lvl {result.level}
                          </span>
                        </div>
                        {pathString && (
                          <p className="text-[11px] text-zinc-400 line-clamp-1 flex items-center gap-1">
                            <span>{pathString}</span>
                          </p>
                        )}
                      </div>
                      <div className="mt-3 flex items-center justify-between text-[11px] font-bold text-[#1769E0] group-hover:translate-x-1 transition-transform">
                        <span>{locale === 'am' ? 'ዝርዝሮችን ይመልከቱ' : 'View Listings'}</span>
                        <ArrowRight size={12} />
                      </div>
                    </Link>
                  );
                })}
              </div>
            ) : !searching ? (
              <div className="text-center py-12">
                <Tag size={32} className="mx-auto text-zinc-300 dark:text-zinc-600 mb-2" />
                <p className="text-sm font-bold text-zinc-500">
                  {locale === 'am' ? 'ምንም አይነት ምድብ አልተገኘም' : 'No matching categories found in canonical taxonomy'}
                </p>
              </div>
            ) : null}
          </div>
        ) : loading ? (
          /* Loading Skeleton */
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
            {Array.from({ length: 10 }).map((_, idx) => (
              <div key={idx} className="h-24 rounded-2xl bg-zinc-100 dark:bg-zinc-800 animate-pulse" />
            ))}
          </div>
        ) : filteredChildren.length > 0 ? (
          /* Category Cards Grid */
          <div>
            {/* Header when inside a branch */}
            {activeParent && (
              <div className="flex items-center justify-between mb-4 pb-2 border-b border-zinc-150 dark:border-zinc-800">
                <span className="text-xs font-bold text-zinc-400 uppercase tracking-wider">
                  {locale === 'am' ? 'የቅርንጫፍ ዝርዝሮች' : `Subcategories under ${getCategoryName(activeParent)}`} ({filteredChildren.length})
                </span>
                <Link
                  href={`/products?category_hierarchy_id=${activeParent.id}&cat=${activeParent.slug}`}
                  className="inline-flex items-center gap-1 text-xs font-bold text-[#1769E0] hover:underline"
                >
                  <span>{locale === 'am' ? 'በዚህ ምድብ ያሉትን ሁሉ ይመልከቱ' : `View All in ${getCategoryName(activeParent)}`}</span>
                  <ArrowRight size={12} />
                </Link>
              </div>
            )}

            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
              {filteredChildren.map((cat) => {
                const name = getCategoryName(cat);
                const isLeaf = cat.is_leaf;

                return (
                  <button
                    key={cat.id}
                    onClick={() => handleCategoryClick(cat)}
                    className="group flex flex-col justify-between p-3.5 sm:p-4 rounded-2xl text-left bg-zinc-50/80 dark:bg-zinc-800/40 hover:bg-white dark:hover:bg-zinc-800 border border-zinc-200/70 dark:border-zinc-700/60 hover:border-[#1769E0] dark:hover:border-[#1769E0] hover:shadow-lg transition-all"
                  >
                    <div>
                      <div className="flex items-center justify-between gap-1.5 mb-2.5">
                        <div className="w-8 h-8 rounded-xl bg-blue-100/70 dark:bg-blue-900/30 flex items-center justify-center text-sm shadow-xs group-hover:scale-110 transition-transform">
                          {isLeaf ? (
                            <Tag size={15} className="text-[#1769E0] dark:text-[#38A5FF]" />
                          ) : (
                            <Folder size={15} className="text-[#1769E0] dark:text-[#38A5FF]" />
                          )}
                        </div>
                        <span className="text-[10px] font-black uppercase tracking-wider px-1.5 py-0.5 rounded-md bg-zinc-200/60 dark:bg-zinc-700 text-zinc-500 dark:text-zinc-400">
                          Lvl {cat.level}
                        </span>
                      </div>

                      <h4 className="text-xs sm:text-sm font-bold text-zinc-800 dark:text-zinc-100 group-hover:text-[#1769E0] dark:group-hover:text-[#38A5FF] line-clamp-2 transition-colors">
                        {name}
                      </h4>
                    </div>

                    <div className="mt-3 pt-2 border-t border-zinc-200/40 dark:border-zinc-700/40 flex items-center justify-between text-[11px] font-bold text-zinc-400 group-hover:text-[#1769E0] transition-colors">
                      <span>{isLeaf ? (locale === 'am' ? 'ምርቶች' : 'View Items') : (locale === 'am' ? 'ያስሱ' : 'Explore')}</span>
                      <ChevronRight size={13} className="group-hover:translate-x-1 transition-transform" />
                    </div>
                  </button>
                );
              })}
            </div>
          </div>
        ) : (
          /* Empty Branch */
          <div className="text-center py-12">
            <FolderOpen size={36} className="mx-auto text-zinc-300 dark:text-zinc-600 mb-2" />
            <p className="text-sm font-bold text-zinc-600 dark:text-zinc-300">
              {locale === 'am' ? 'በዚህ ቅርንጫፍ ውስጥ ተጨማሪ ንዑስ ምድቦች የሉም' : 'No further subcategories in this branch'}
            </p>
            {activeParent && (
              <Link
                href={`/products?category_hierarchy_id=${activeParent.id}&cat=${activeParent.slug}`}
                className="mt-4 inline-flex items-center gap-2 px-5 py-2.5 rounded-xl text-xs font-bold bg-[#1769E0] text-white hover:bg-[#1055BD] transition-all shadow-md shadow-[#1769E0]/20"
              >
                <span>{locale === 'am' ? 'የተዘረዘሩ እቃዎችን ይመልከቱ' : `View Listings in ${getCategoryName(activeParent)}`}</span>
                <ArrowRight size={13} />
              </Link>
            )}
          </div>
        )}
      </div>

      {/* Footer CTA */}
      <div className="p-4 sm:p-5 bg-zinc-50 dark:bg-zinc-950 border-t border-zinc-200/80 dark:border-zinc-800 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs text-zinc-500 dark:text-zinc-400">
        <div className="flex items-center gap-2">
          <Sparkles size={14} className="text-amber-500" />
          <span>
            {locale === 'am'
              ? '20,500 ዝርዝር ምድቦች በ 4 ቋንቋዎች ተደራጅተዋል'
              : '20,500 canonical categories synchronized across English, Amharic, Oromo, and Tigrinya'}
          </span>
        </div>
        <Link
          href="/products"
          className="font-bold text-[#1769E0] hover:text-[#1055BD] dark:text-[#38A5FF] flex items-center gap-1 transition-colors"
        >
          <span>{locale === 'am' ? 'ሁሉንም ምርቶች ያስሱ' : 'Browse All Marketplace Listings'}</span>
          <ArrowRight size={12} />
        </Link>
      </div>
    </div>
  );
}
export default CategoryBrowser;
