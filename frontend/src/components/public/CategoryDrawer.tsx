'use client';

import { useState, useEffect, useRef, useMemo } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useLocale } from '@/hooks/useLocale';
import { CATEGORY_TREE } from '@/lib/constants/categoryTree';
import {
  X,
  Search,
  ChevronRight,
  ChevronDown,
  Layers,
  Sparkles,
  ArrowRight,
  Store,
  ShoppingBag,
  Home,
  Tag,
} from 'lucide-react';

const CATEGORY_META: Record<string, { emoji: string; color: string }> = {
  electronics: { emoji: '📱', color: 'from-blue-500 to-indigo-600' },
  vehicles: { emoji: '🚗', color: 'from-amber-500 to-orange-600' },
  property: { emoji: '🏢', color: 'from-emerald-500 to-teal-600' },
  home_appliances: { emoji: '🍳', color: 'from-rose-500 to-red-600' },
  furniture: { emoji: '🛋️', color: 'from-yellow-600 to-amber-700' },
  clothing: { emoji: '👗', color: 'from-purple-500 to-pink-600' },
  agriculture: { emoji: '🌾', color: 'from-lime-600 to-green-700' },
  food: { emoji: '🍎', color: 'from-red-500 to-rose-600' },
  beauty_health: { emoji: '💄', color: 'from-pink-500 to-rose-500' },
  sports: { emoji: '⚽', color: 'from-cyan-500 to-blue-600' },
  tools: { emoji: '🔧', color: 'from-zinc-600 to-zinc-800' },
  books: { emoji: '📚', color: 'from-indigo-600 to-violet-700' },
  services: { emoji: '💼', color: 'from-teal-600 to-emerald-700' },
  jobs: { emoji: '🤝', color: 'from-blue-600 to-cyan-700' },
  industrial: { emoji: '🏗️', color: 'from-orange-600 to-amber-700' },
};

interface CategoryDrawerProps {
  isOpen: boolean;
  onClose: () => void;
}

export function CategoryDrawer({ isOpen, onClose }: CategoryDrawerProps) {
  const { t } = useLocale();
  const router = useRouter();
  const [activeCategory, setActiveCategory] = useState<string>('electronics');
  const [searchQuery, setSearchQuery] = useState('');
  const [expandedMobileCat, setExpandedMobileCat] = useState<string | null>(null);
  const drawerRef = useRef<HTMLDivElement>(null);

  // Category labels helper
  const catLabel = (slug: string) =>
    (t.categories as any)?.[slug] ?? slug.replace(/_/g, ' ');

  const subLabel = (catSlug: string, subKey: string) => {
    const fromSub = (t.sub_categories as any)?.[catSlug]?.[subKey];
    if (fromSub) return fromSub;
    const fromItems = (t.items as any)?.[subKey];
    if (fromItems) return fromItems;
    return subKey.replace(/_/g, ' ');
  };

  const itemLabel = (slug: string) =>
    (t.items as any)?.[slug] ?? slug.replace(/_/g, ' ');

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
  const categoriesList = useMemo(() => {
    const keys = Object.keys(CATEGORY_TREE);
    if (!searchQuery.trim()) return keys;
    const q = searchQuery.toLowerCase();
    return keys.filter(cat => {
      const label = catLabel(cat).toLowerCase();
      if (label.includes(q) || cat.includes(q)) return true;
      const subs = CATEGORY_TREE[cat] || {};
      return Object.keys(subs).some(sub => {
        const sLabel = subLabel(cat, sub).toLowerCase();
        return sLabel.includes(q) || sub.includes(q);
      });
    });
  }, [searchQuery, t]);

  if (!isOpen) return null;

  const currentSubs = CATEGORY_TREE[activeCategory] || {};

  const handleNavigate = (url: string) => {
    onClose();
    router.push(url);
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
                Alibaba-style universal catalog & verified Ethiopian listings
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
              onChange={e => setSearchQuery(e.target.value)}
              placeholder="Filter categories or subcategories..."
              className="w-full pl-10 pr-4 py-2 text-xs rounded-xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900 focus:outline-none focus:ring-2 focus:ring-orange-600 focus:bg-white dark:focus:bg-zinc-950 transition-all"
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

        {/* Desktop 2-column view (hidden on mobile, flex on md+) */}
        <div className="hidden md:flex flex-1 overflow-hidden">
          {/* Left: Main Categories List */}
          <div className="w-72 border-r border-zinc-200 dark:border-zinc-800 overflow-y-auto p-2 space-y-1 bg-zinc-50/50 dark:bg-zinc-900/30">
            {categoriesList.map(cat => {
              const meta = CATEGORY_META[cat] || { emoji: '📦', color: 'from-zinc-500 to-zinc-700' };
              const isActive = cat === activeCategory;
              const subCount = Object.keys(CATEGORY_TREE[cat] || {}).length;

              return (
                <button
                  key={cat}
                  onClick={() => setActiveCategory(cat)}
                  className={`w-full text-left px-3.5 py-2.5 rounded-xl text-xs font-bold transition-all flex items-center justify-between group ${
                    isActive
                      ? 'bg-orange-600 text-white shadow-md shadow-orange-600/20'
                      : 'text-zinc-700 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800'
                  }`}
                >
                  <div className="flex items-center gap-2.5 truncate">
                    <span className="text-base shrink-0">{meta.emoji}</span>
                    <span className="truncate">{catLabel(cat)}</span>
                  </div>
                  <div className="flex items-center gap-1 shrink-0">
                    <span
                      className={`text-[10px] px-1.5 py-0.5 rounded-md font-semibold ${
                        isActive
                          ? 'bg-white/20 text-white'
                          : 'bg-zinc-200 dark:bg-zinc-800 text-zinc-500'
                      }`}
                    >
                      {subCount}
                    </span>
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
          <div className="flex-1 overflow-y-auto p-6 space-y-6 bg-white dark:bg-zinc-950">
            {/* Header for Active Category */}
            <div className="flex items-center justify-between pb-3 border-b border-zinc-100 dark:border-zinc-800">
              <div className="flex items-center gap-2">
                <span className="text-2xl">{CATEGORY_META[activeCategory]?.emoji ?? '📦'}</span>
                <div>
                  <h3 className="text-base font-black text-zinc-900 dark:text-white">
                    {catLabel(activeCategory)}
                  </h3>
                  <p className="text-[11px] text-zinc-400">
                    Browse all verified {catLabel(activeCategory)} items
                  </p>
                </div>
              </div>
              <button
                onClick={() => handleNavigate(`/products?cat=${activeCategory}`)}
                className="flex items-center gap-1 px-3 py-1.5 rounded-lg text-xs font-black text-orange-600 hover:bg-orange-50 dark:hover:bg-orange-950/20 transition-colors"
              >
                <span>View All</span>
                <ArrowRight size={12} />
              </button>
            </div>

            {/* Subcategories Grid */}
            <div className="grid grid-cols-2 gap-6">
              {Object.entries(currentSubs).map(([sub, items]) => (
                <div key={sub} className="space-y-2">
                  <button
                    onClick={() => handleNavigate(`/products?cat=${activeCategory}&subcategory=${sub}`)}
                    className="text-xs font-black uppercase tracking-wider text-zinc-800 dark:text-zinc-200 hover:text-orange-600 dark:hover:text-orange-400 flex items-center gap-1.5 group text-left"
                  >
                    <span className="w-1.5 h-1.5 rounded-full bg-orange-600" />
                    <span>{subLabel(activeCategory, sub)}</span>
                    <ChevronRight size={11} className="opacity-0 group-hover:opacity-100 transition-opacity text-orange-600" />
                  </button>
                  <div className="flex flex-wrap gap-1.5 pl-3">
                    {items.map(item => (
                      <button
                        key={item}
                        onClick={() => handleNavigate(`/products?cat=${activeCategory}&subcategory=${sub}&item=${item}`)}
                        className="text-[11px] text-zinc-500 dark:text-zinc-400 hover:text-orange-600 dark:hover:text-white px-2 py-1 rounded-md hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors capitalize"
                      >
                        {itemLabel(item)}
                      </button>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Mobile Accordion View (visible on md- screens) */}
        <div className="md:hidden flex-1 overflow-y-auto p-3 space-y-2 bg-zinc-50/50 dark:bg-zinc-900/30">
          {categoriesList.map(cat => {
            const meta = CATEGORY_META[cat] || { emoji: '📦', color: 'from-zinc-500 to-zinc-700' };
            const isExpanded = expandedMobileCat === cat;
            const subs = CATEGORY_TREE[cat] || {};

            return (
              <div
                key={cat}
                className="bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200/70 dark:border-zinc-800 overflow-hidden"
              >
                <div className="flex items-center justify-between p-3.5">
                  <button
                    onClick={() => handleNavigate(`/products?cat=${cat}`)}
                    className="flex items-center gap-2.5 text-xs font-black text-zinc-900 dark:text-white hover:text-orange-600 text-left flex-1"
                  >
                    <span className="text-xl">{meta.emoji}</span>
                    <span>{catLabel(cat)}</span>
                  </button>
                  <button
                    onClick={() => setExpandedMobileCat(isExpanded ? null : cat)}
                    className="p-1.5 rounded-lg text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800"
                    aria-label={`Expand ${catLabel(cat)}`}
                  >
                    <ChevronDown
                      size={16}
                      className={`transition-transform duration-200 ${isExpanded ? 'rotate-180 text-orange-600' : ''}`}
                    />
                  </button>
                </div>

                {isExpanded && (
                  <div className="px-4 pb-4 pt-1 border-t border-zinc-100 dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-950/40 space-y-3">
                    {Object.entries(subs).map(([sub, items]) => (
                      <div key={sub} className="space-y-1.5">
                        <button
                          onClick={() => handleNavigate(`/products?cat=${cat}&subcategory=${sub}`)}
                          className="text-xs font-bold text-orange-600 flex items-center gap-1 text-left"
                        >
                          <span>{subLabel(cat, sub)}</span>
                          <ArrowRight size={10} />
                        </button>
                        <div className="flex flex-wrap gap-1 pl-2">
                          {items.slice(0, 6).map(item => (
                            <button
                              key={item}
                              onClick={() => handleNavigate(`/products?cat=${cat}&subcategory=${sub}&item=${item}`)}
                              className="text-[10px] text-zinc-500 dark:text-zinc-400 bg-white dark:bg-zinc-900 px-2 py-0.5 rounded border border-zinc-200/60 dark:border-zinc-800"
                            >
                              {itemLabel(item)}
                            </button>
                          ))}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            );
          })}
        </div>

        {/* Drawer Footer with Quick Mode Actions */}
        <div className="p-4 border-t border-zinc-200 dark:border-zinc-800 bg-zinc-50/90 dark:bg-zinc-900/90 flex flex-wrap items-center justify-between gap-3 shrink-0">
          <div className="flex items-center gap-2">
            <Link
              href="/products?type=sale"
              onClick={onClose}
              className="px-3 py-1.5 rounded-xl bg-orange-600/10 text-orange-600 text-xs font-black uppercase tracking-wider hover:bg-orange-600 hover:text-white transition-colors"
            >
              🛒 Buy
            </Link>
            <Link
              href="/products?type=rent"
              onClick={onClose}
              className="px-3 py-1.5 rounded-xl bg-teal-600/10 text-teal-600 text-xs font-black uppercase tracking-wider hover:bg-teal-600 hover:text-white transition-colors"
            >
              🏠 Rental
            </Link>
            <Link
              href="/products?cat=property"
              onClick={onClose}
              className="px-3 py-1.5 rounded-xl bg-emerald-600/10 text-emerald-600 text-xs font-black uppercase tracking-wider hover:bg-emerald-600 hover:text-white transition-colors"
            >
              🏢 Land & Property
            </Link>
          </div>
          <Link
            href="/sell"
            onClick={onClose}
            className="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-orange-600 text-white text-xs font-black uppercase tracking-wider hover:bg-orange-500 shadow-md shadow-orange-600/20 transition-all active:scale-95"
          >
            <Store size={14} />
            <span>List An Item / Property</span>
          </Link>
        </div>
      </div>
    </div>
  );
}
