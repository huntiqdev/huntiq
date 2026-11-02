'use client';

import { useState, useEffect, useRef, useMemo, useCallback } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import { useLocale } from '@/hooks/useLocale';
import { getCategoryLabel } from '@/lib/taxonomy';
import { api } from '@/lib/api';
import {
  X,
  Search,
  ChevronRight,
  ChevronDown,
  Layers,
  ArrowRight,
  Store,
  Tag,
  AlertCircle,
  RefreshCw,
} from 'lucide-react';

interface CategoryNode {
  id: number;
  name: string;
  name_en?: string;
  name_am?: string;
  name_om?: string;
  name_ti?: string;
  slug: string;
  level: number;
  icon?: string;
  image_url?: string;
  is_leaf: boolean;
  children?: CategoryNode[];
}

interface CategoryDrawerProps {
  isOpen: boolean;
  onClose: () => void;
}

export function CategoryDrawer({ isOpen, onClose }: CategoryDrawerProps) {
  const { t, locale } = useLocale();
  const router = useRouter();
  const [tree, setTree] = useState<CategoryNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeCategoryId, setActiveCategoryId] = useState<number | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [expandedMobileCat, setExpandedMobileCat] = useState<number | null>(null);
  const drawerRef = useRef<HTMLDivElement>(null);

  // Load canonical category tree from the Go backend API
  const loadCategories = useCallback(() => {
    setLoading(true);
    setError(null);
    api.getCategoryTree(locale)
      .then((data) => {
        if (Array.isArray(data) && data.length > 0) {
          setTree(data);
          setActiveCategoryId(data[0].id);
        } else {
          setTree([]);
        }
      })
      .catch(() => {
        setError('Failed to load categories from marketplace catalog');
        setTree([]);
      })
      .finally(() => {
        setLoading(false);
      });
  }, [locale]);

  useEffect(() => {
    if (!isOpen || tree.length > 0) return;
    let active = true;
    api.getCategoryTree(locale)
      .then((data) => {
        if (!active) return;
        if (Array.isArray(data) && data.length > 0) {
          setTree(data);
          setActiveCategoryId(data[0].id);
        } else {
          setTree([]);
        }
      })
      .catch(() => {
        if (!active) return;
        setError('Failed to load categories from marketplace catalog');
        setTree([]);
      })
      .finally(() => {
        if (active) {
          setLoading(false);
        }
      });
    return () => {
      active = false;
    };
  }, [isOpen, tree.length, locale]);

  // Prevent background scrolling when drawer is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  // Handle escape key
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);

  // Filtered categories based on search
  const filteredTree = useMemo(() => {
    if (!searchQuery.trim()) return tree;
    const q = searchQuery.toLowerCase();
    return tree.filter((cat) => {
      const matchCat = (cat.name || cat.slug).toLowerCase().includes(q);
      if (matchCat) return true;
      const matchChildren = cat.children?.some((sub) => {
        const matchSub = (sub.name || sub.slug).toLowerCase().includes(q);
        if (matchSub) return true;
        return sub.children?.some((item) => (item.name || item.slug).toLowerCase().includes(q));
      });
      return matchChildren;
    });
  }, [tree, searchQuery]);

  if (!isOpen) return null;

  const activeCategory = tree.find((c) => c.id === activeCategoryId) || tree[0];

  const handleNavigate = (url: string) => {
    onClose();
    router.push(url);
  };

  const renderIcon = (cat: CategoryNode) => {
    if (cat.icon) {
      return <span className="text-lg">{cat.icon}</span>;
    }
    if (cat.image_url) {
      return (
        <Image
          src={cat.image_url}
          alt={cat.name || 'Category'}
          width={20}
          height={20}
          className="w-5 h-5 object-contain"
          unoptimized
        />
      );
    }
    return <Tag size={16} className="text-zinc-400" />;
  };

  const getCatLabel = (node: CategoryNode) => {
    if (node.id) {
      const translated = getCategoryLabel(node.id, locale);
      if (translated) return translated;
    }
    if (node.slug && (t.categories as any)?.[node.slug]) {
      return (t.categories as any)[node.slug];
    }
    if (locale === 'am' && node.name_am) return node.name_am;
    if ((locale === 'or' || locale === 'om') && node.name_om) return node.name_om;
    if (locale === 'ti' && node.name_ti) return node.name_ti;
    if (locale === 'en' && node.name_en) return node.name_en;
    return node.name || node.name_en || node.slug?.replace(/_/g, ' ') || '';
  };

  const getSubLabel = (sub: CategoryNode, parentSlug?: string) => {
    if (sub.id) {
      const translated = getCategoryLabel(sub.id, locale);
      if (translated) return translated;
    }
    if (locale === 'am' && sub.name_am) return sub.name_am;
    if ((locale === 'or' || locale === 'om') && sub.name_om) return sub.name_om;
    if (locale === 'ti' && sub.name_ti) return sub.name_ti;
    if (locale === 'en' && sub.name_en) return sub.name_en;
    return sub.name || sub.name_en || sub.slug?.replace(/_/g, ' ') || '';
  };

  return (
    <div className="fixed inset-0 z-[70] flex items-stretch">
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/60 backdrop-blur-sm transition-opacity duration-300"
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Drawer Container */}
      <div
        ref={drawerRef}
        role="dialog"
        aria-modal="true"
        aria-label="Categories Navigation"
        className="relative z-10 w-full max-w-4xl bg-white dark:bg-zinc-950 shadow-2xl flex flex-col h-full border-r border-zinc-200 dark:border-zinc-800 animate-in slide-in-from-left duration-300"
      >
        {/* Drawer Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-zinc-200 dark:border-zinc-800 bg-zinc-50/70 dark:bg-zinc-900/50 shrink-0">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-xl bg-orange-600 flex items-center justify-center text-white shadow-md shadow-orange-600/20">
              <Layers size={18} />
            </div>
            <div>
              <h2 className="text-base font-black tracking-tight text-zinc-900 dark:text-white">
                {(t.nav as any).categories ?? 'All Categories'}
              </h2>
              <p className="text-[11px] text-zinc-500 dark:text-zinc-400">
                Authoritative database hierarchy & verified Ethiopian listings
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 rounded-xl text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-zinc-200/60 dark:hover:bg-zinc-800 transition-colors"
            aria-label="Close categories drawer"
          >
            <X size={20} />
          </button>
        </div>

        {/* Drawer Search */}
        <div className="p-4 border-b border-zinc-100 dark:border-zinc-800/80 bg-white dark:bg-zinc-950 shrink-0">
          <div className="relative">
            <Search size={16} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search marketplace categories..."
              className="w-full pl-10 pr-4 py-2.5 text-base sm:text-xs min-h-[44px] rounded-xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900 focus:outline-none focus:ring-2 focus:ring-orange-600 focus:bg-white dark:focus:bg-zinc-950 transition-all"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery('')}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-400 hover:text-zinc-600"
              >
                <X size={14} />
              </button>
            )}
          </div>
        </div>

        {/* Loading and Error States */}
        {loading ? (
          <div className="flex-1 flex items-center justify-center p-8">
            <div className="flex flex-col items-center gap-3">
              <RefreshCw size={24} className="text-orange-600 animate-spin" />
              <p className="text-xs text-zinc-400 font-semibold">Loading taxonomy catalog...</p>
            </div>
          </div>
        ) : error ? (
          <div className="flex-1 flex items-center justify-center p-8">
            <div className="max-w-sm text-center space-y-3">
              <AlertCircle size={28} className="text-red-500 mx-auto" />
              <p className="text-sm font-bold text-zinc-800 dark:text-zinc-200">{error}</p>
              <button
                onClick={loadCategories}
                className="px-4 py-2 bg-orange-600 hover:bg-orange-500 text-white rounded-xl text-xs font-bold shadow-sm transition-colors"
              >
                Retry Loading
              </button>
            </div>
          </div>
        ) : filteredTree.length === 0 ? (
          <div className="flex-1 flex items-center justify-center p-8">
            <p className="text-xs text-zinc-400">No categories found matching &quot;{searchQuery}&quot;</p>
          </div>
        ) : (
          <>
            {/* Desktop 2-column view (hidden on mobile, flex on md+) */}
            <div className="hidden md:flex flex-1 overflow-hidden">
              {/* Left: Main Categories List */}
              <div className="w-72 border-r border-zinc-200 dark:border-zinc-800 overflow-y-auto p-2 space-y-1 bg-zinc-50/50 dark:bg-zinc-900/30">
                {filteredTree.map((cat) => {
                  const isActive = activeCategory && cat.id === activeCategory.id;
                  const childCount = cat.children?.length ?? 0;

                  return (
                    <button
                      key={cat.id}
                      onClick={() => setActiveCategoryId(cat.id)}
                      className={`w-full text-left px-3.5 py-2.5 rounded-xl text-xs font-bold transition-all flex items-center justify-between group ${
                        isActive
                          ? 'bg-orange-600 text-white shadow-md shadow-orange-600/20'
                          : 'text-zinc-700 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800'
                      }`}
                    >
                      <div className="flex items-center gap-2.5 truncate">
                        <span className="shrink-0">{renderIcon(cat)}</span>
                        <span className="truncate">{getCatLabel(cat)}</span>
                      </div>
                      <div className="flex items-center gap-1 shrink-0">
                        {childCount > 0 && (
                          <span
                            className={`text-[10px] px-1.5 py-0.5 rounded-md font-semibold ${
                              isActive
                                ? 'bg-white/20 text-white'
                                : 'bg-zinc-200 dark:bg-zinc-800 text-zinc-500'
                            }`}
                          >
                            {childCount}
                          </span>
                        )}
                        <ChevronRight
                          size={14}
                          className={isActive ? 'text-white' : 'text-zinc-400 group-hover:translate-x-0.5 transition-transform'}
                        />
                      </div>
                    </button>
                  );
                })}
              </div>

              {/* Right: Subcategories and Specific Items */}
              {activeCategory && (
                <div className="flex-1 overflow-y-auto p-6 space-y-6 bg-white dark:bg-zinc-950">
                  {/* Header for Active Category */}
                  <div className="flex items-center justify-between pb-3 border-b border-zinc-100 dark:border-zinc-800">
                    <div className="flex items-center gap-2.5">
                      <span className="text-2xl">{renderIcon(activeCategory)}</span>
                      <div>
                        <h3 className="text-base font-black text-zinc-900 dark:text-white">
                          {getCatLabel(activeCategory)}
                        </h3>
                        <p className="text-[11px] text-zinc-400">
                          {getCatLabel(activeCategory)}
                        </p>
                      </div>
                    </div>
                    <button
                      onClick={() =>
                        handleNavigate(
                          `/products?category_hierarchy_id=${activeCategory.id}&cat=${activeCategory.slug}`
                        )
                      }
                      className="flex items-center gap-1 px-3 py-1.5 rounded-lg text-xs font-black text-orange-600 hover:bg-orange-50 dark:hover:bg-orange-950/20 transition-colors"
                    >
                      <span>{locale === 'am' ? 'ሁሉንም እይ' : locale === 'or' ? 'Hunda Ilaali' : locale === 'ti' ? 'ንኹሉ ርኣይ' : 'View All'}</span>
                      <ArrowRight size={12} />
                    </button>
                  </div>

                  {/* Subcategories Grid */}
                  {activeCategory.children && activeCategory.children.length > 0 ? (
                    <div className="grid grid-cols-2 gap-6">
                      {activeCategory.children.map((sub) => (
                        <div key={sub.id} className="space-y-2">
                          <button
                            onClick={() =>
                              handleNavigate(
                                `/products?category_hierarchy_id=${sub.id}&cat=${activeCategory.slug}&subcategory=${sub.slug}`
                              )
                            }
                            className="text-xs font-black uppercase tracking-wider text-zinc-800 dark:text-zinc-200 hover:text-orange-600 dark:hover:text-orange-400 flex items-center gap-1.5 group text-left"
                          >
                            <span className="w-1.5 h-1.5 rounded-full bg-orange-600" />
                            <span>{getSubLabel(sub, activeCategory.slug)}</span>
                            <ChevronRight
                              size={11}
                              className="opacity-0 group-hover:opacity-100 transition-opacity text-orange-600"
                            />
                          </button>
                          {sub.children && sub.children.length > 0 && (
                            <div className="flex flex-wrap gap-1.5 pl-3">
                              {sub.children.map((item) => (
                                <button
                                  key={item.id}
                                  onClick={() =>
                                    handleNavigate(
                                      `/products?category_hierarchy_id=${item.id}&cat=${activeCategory.slug}&subcategory=${sub.slug}&item=${item.slug}`
                                    )
                                  }
                                  className="text-[11px] text-zinc-500 dark:text-zinc-400 hover:text-orange-600 dark:hover:text-white px-2 py-1 rounded-md hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors capitalize"
                                >
                                  {item.name || item.name_en || item.slug}
                                </button>
                              ))}
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="text-center py-12 text-zinc-400 text-xs">
                      No further subcategories. Click &quot;View All&quot; to see all listings.
                    </div>
                  )}
                </div>
              )}
            </div>

            {/* Mobile Accordion View (visible on md- screens) */}
            <div className="md:hidden flex-1 overflow-y-auto p-3 space-y-2 bg-zinc-50/50 dark:bg-zinc-900/30">
              {filteredTree.map((cat) => {
                const isExpanded = expandedMobileCat === cat.id;

                return (
                  <div
                    key={cat.id}
                    className="bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200/70 dark:border-zinc-800 overflow-hidden"
                  >
                    <div className="flex items-center justify-between p-3.5">
                      <button
                        onClick={() =>
                          handleNavigate(
                            `/products?category_hierarchy_id=${cat.id}&cat=${cat.slug}`
                          )
                        }
                        className="flex items-center gap-2.5 text-xs font-black text-zinc-900 dark:text-white hover:text-orange-600 text-left flex-1"
                      >
                        <span className="shrink-0">{renderIcon(cat)}</span>
                        <span>{getCatLabel(cat)}</span>
                      </button>
                      {cat.children && cat.children.length > 0 && (
                        <button
                          onClick={() => setExpandedMobileCat(isExpanded ? null : cat.id)}
                          className="p-1.5 rounded-lg text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800"
                          aria-label={`Expand ${getCatLabel(cat)}`}
                        >
                          <ChevronDown
                            size={16}
                            className={`transition-transform duration-200 ${
                              isExpanded ? 'rotate-180 text-orange-600' : ''
                            }`}
                          />
                        </button>
                      )}
                    </div>

                    {isExpanded && cat.children && (
                      <div className="px-4 pb-4 pt-1 border-t border-zinc-100 dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-950/40 space-y-3">
                        {cat.children.map((sub) => (
                          <div key={sub.id} className="space-y-1.5">
                            <button
                              onClick={() =>
                                handleNavigate(
                                  `/products?category_hierarchy_id=${sub.id}&cat=${cat.slug}&subcategory=${sub.slug}`
                                )
                              }
                              className="text-xs font-bold text-orange-600 flex items-center gap-1 text-left"
                            >
                              <span>{getSubLabel(sub, cat.slug)}</span>
                              <ArrowRight size={10} />
                            </button>
                            {sub.children && sub.children.length > 0 && (
                              <div className="flex flex-wrap gap-1 pl-2">
                                {sub.children.slice(0, 8).map((item) => (
                                  <button
                                    key={item.id}
                                    onClick={() =>
                                      handleNavigate(
                                        `/products?category_hierarchy_id=${item.id}&cat=${cat.slug}&subcategory=${sub.slug}&item=${item.slug}`
                                      )
                                    }
                                    className="text-[10px] text-zinc-500 dark:text-zinc-400 bg-white dark:bg-zinc-900 px-2 py-0.5 rounded border border-zinc-200/60 dark:border-zinc-800"
                                  >
                                    {item.name || item.name_en || item.slug}
                                  </button>
                                ))}
                              </div>
                            )}
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          </>
        )}

        {/* Drawer Footer with Quick Mode Actions */}
        <div className="p-4 border-t border-zinc-200 dark:border-zinc-800 bg-zinc-50/90 dark:bg-zinc-900/90 pb-[max(16px,env(safe-area-inset-bottom,16px))] flex flex-wrap items-center justify-between gap-3 shrink-0">
          <div className="flex items-center gap-2">
            <Link
              href="/products?type=sale"
              onClick={onClose}
              className="px-3.5 py-2 rounded-xl bg-orange-600/10 text-orange-600 text-xs font-black uppercase tracking-wider hover:bg-orange-600 hover:text-white transition-colors tap-bounce active:scale-95"
            >
              🛒 {t.nav?.buy ?? 'Buy'}
            </Link>
            <Link
              href="/products?type=rent"
              onClick={onClose}
              className="px-3.5 py-2 rounded-xl bg-teal-600/10 text-teal-600 text-xs font-black uppercase tracking-wider hover:bg-teal-600 hover:text-white transition-colors tap-bounce active:scale-95"
            >
              🏠 {t.nav?.rental ?? 'Rental'}
            </Link>
          </div>
          <Link
            href="/sell"
            onClick={onClose}
            className="flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-orange-600 text-white text-xs font-black uppercase tracking-wider hover:bg-orange-500 shadow-md shadow-orange-600/20 transition-all tap-bounce active:scale-95 min-h-[40px]"
          >
            <Store size={14} />
            <span>{t.nav?.sell?.split('/')[0]?.trim() ?? 'Post A Listing'}</span>
          </Link>
        </div>
      </div>
    </div>
  );
}
