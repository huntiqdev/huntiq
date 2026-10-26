'use client';
import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { Category, ProductFilters } from '@/types';
import { LocationSelector, LocationValue } from '@/components/location/LocationSelector';
import { useLocale } from '@/hooks/useLocale';
import { api } from '@/lib/api';
import { SlidersHorizontal, RotateCcw } from 'lucide-react';

interface Props {
  categories: Category[];
  onFilter: (f: ProductFilters) => void;
}

const CONDITIONS = [
  { value: 'new', labelKey: 'condition_new' },
  { value: 'used', labelKey: 'condition_used' },
] as const;

export function FilterSidebar({ categories, onFilter }: Props) {
  const { t, locale } = useLocale();
  const router = useRouter();
  const [categoryId, setCategoryId] = useState('');
  const [subcategory, setSubcategory] = useState('');
  const [item, setItem] = useState('');
  const [minPrice, setMinPrice] = useState('');
  const [maxPrice, setMaxPrice] = useState('');
  const [condition, setCondition] = useState('');
  const [listingType, setListingType] = useState<'all' | 'sale' | 'rent'>('all');
  const [locationQuery, setLocationQuery] = useState<LocationValue>({});

  // Dynamic child categories loaded from database hierarchy
  const [loadedChildren, setLoadedChildren] = useState<Category[]>([]);

  const handleLocationChange = useCallback((v: LocationValue) => {
    setLocationQuery(v);
  }, []);

  const selectedCat = categories.find(c => String(c.id) === categoryId);

  // If selectedCat already has children in the tree, use them. Otherwise fetch via API.
  useEffect(() => {
    if (!categoryId) {
      setLoadedChildren([]);
      setSubcategory('');
      setItem('');
      return;
    }

    if (selectedCat?.children && selectedCat.children.length > 0) {
      setLoadedChildren(selectedCat.children);
    } else {
      api.getCategoryChildren(Number(categoryId), locale)
        .then(children => setLoadedChildren(Array.isArray(children) ? children : []))
        .catch(() => setLoadedChildren([]));
    }
  }, [categoryId, selectedCat, locale]);

  const selectedSub = loadedChildren.find(sc => sc.slug === subcategory || String(sc.id) === subcategory);
  const itemOptions = selectedSub?.children ?? [];

  const apply = () => {
    const f: ProductFilters = {};
    if (listingType !== 'all') f.listing_type = listingType;
    if (categoryId) {
      f.category_id = Number(categoryId);
      f.category_hierarchy_id = selectedSub?.id ?? Number(categoryId);
    }
    if (subcategory) f.subcategory = subcategory;
    if (item) f.item = item;
    if (minPrice) f.min_price = Number(minPrice);
    if (maxPrice) f.max_price = Number(maxPrice);
    if (condition) f.condition = condition;
    if (locationQuery.location_hierarchy_id) f.location_hierarchy_id = locationQuery.location_hierarchy_id;
    if (locationQuery.woreda) f.woreda = locationQuery.woreda;
    else if (locationQuery.subcity) f.subcity = locationQuery.subcity;
    else if (locationQuery.zone) f.zone = locationQuery.zone;
    else if (locationQuery.region) f.region = locationQuery.region;
    onFilter(f);

    const p = new URLSearchParams();
    if (listingType !== 'all') p.set('type', listingType);
    if (f.category_id) p.set('category_id', String(f.category_id));
    if (f.category_hierarchy_id) p.set('category_hierarchy_id', String(f.category_hierarchy_id));
    if (f.subcategory) p.set('subcategory', f.subcategory);
    if (f.item) p.set('item', f.item);
    if (f.location_hierarchy_id) p.set('location_hierarchy_id', String(f.location_hierarchy_id));
    if (f.woreda) p.set('woreda', f.woreda);
    else if (f.subcity) p.set('subcity', f.subcity);
    else if (f.zone) p.set('zone', f.zone);
    else if (f.region) p.set('region', f.region);
    router.push(`/products${p.toString() ? '?' + p.toString() : ''}`);
  };

  const clear = () => {
    setListingType('all');
    setCategoryId('');
    setSubcategory('');
    setItem('');
    setMinPrice('');
    setMaxPrice('');
    setCondition('');
    setLocationQuery({} as LocationValue);
    onFilter({});
  };

  const cls = 'w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 transition-all';
  const lbl = 'text-xs font-black text-zinc-500 uppercase tracking-wider';

  return (
    <aside className="bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 p-4 flex flex-col gap-4 h-fit">
      <div className="flex items-center gap-2">
        <SlidersHorizontal size={16} className="text-orange-600" />
        <h2 className="font-bold">{t.ui.filters}</h2>
      </div>

      {/* Dimension 1: Marketplace Pillar (Transaction Mode) */}
      <div className="space-y-1.5 bg-zinc-50 dark:bg-zinc-800/60 p-2.5 rounded-2xl border border-zinc-200 dark:border-zinc-700/60">
        <p className={lbl}>Marketplace Pillar</p>
        <div className="grid grid-cols-3 gap-1 bg-zinc-200/70 dark:bg-zinc-900 p-1 rounded-xl">
          <button
            type="button"
            onClick={() => setListingType('all')}
            className={`py-1.5 text-[11px] font-black uppercase tracking-wider rounded-lg transition-all ${
              listingType === 'all'
                ? 'bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white shadow-sm'
                : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'
            }`}
          >
            All
          </button>
          <button
            type="button"
            onClick={() => setListingType('sale')}
            className={`py-1.5 text-[11px] font-black uppercase tracking-wider rounded-lg transition-all ${
              listingType === 'sale'
                ? 'bg-orange-600 text-white shadow-sm'
                : 'text-zinc-500 hover:text-orange-600'
            }`}
          >
            🛒 Buy
          </button>
          <button
            type="button"
            onClick={() => setListingType('rent')}
            className={`py-1.5 text-[11px] font-black uppercase tracking-wider rounded-lg transition-all ${
              listingType === 'rent'
                ? 'bg-teal-600 text-white shadow-sm'
                : 'text-zinc-500 hover:text-teal-600'
            }`}
          >
            🏠 Rent
          </button>
        </div>
      </div>

      {/* Dimension 2: Product & Asset Taxonomy (Categories) */}
      <div className="space-y-1.5">
        <p className={lbl}>{t.ui.selectCategory}</p>
        <select
          className={cls}
          value={categoryId}
          onChange={e => {
            setCategoryId(e.target.value);
            setSubcategory('');
            setItem('');
          }}
        >
          <option value="">{t.categories.all}</option>
          {categories.map(c => (
            <option key={c.id} value={c.id}>
              {c.name || c.name_en || c.slug}
            </option>
          ))}
        </select>
      </div>

      {/* Level 2 — Subcategory */}
      {loadedChildren.length > 0 && (
        <div className="space-y-1.5">
          <p className={lbl}>{t.ui.subCategoryLabel}</p>
          <select
            className={cls}
            value={subcategory}
            onChange={e => {
              setSubcategory(e.target.value);
              setItem('');
            }}
          >
            <option value="">— {t.ui.subCategoryLabel} —</option>
            {loadedChildren.map(sub => (
              <option key={sub.id} value={sub.slug}>
                {sub.name || sub.name_en || sub.slug}
              </option>
            ))}
          </select>
        </div>
      )}

      {/* Level 3 — Item */}
      {itemOptions.length > 0 && (
        <div className="space-y-1.5">
          <p className={lbl}>{(t.ui as any).itemLabel ?? 'Specific Item'}</p>
          <select className={cls} value={item} onChange={e => setItem(e.target.value)}>
            <option value="">— {(t.ui as any).itemLabel ?? 'All'} —</option>
            {itemOptions.map(it => (
              <option key={it.id} value={it.slug}>
                {it.name || it.name_en || it.slug}
              </option>
            ))}
          </select>
        </div>
      )}

      {/* Dimension 3: Item Condition (Physical State) */}
      <div className="space-y-1.5">
        <p className={lbl}>{t.ui.condition} (Physical State)</p>
        <select className={cls} value={condition} onChange={e => setCondition(e.target.value)}>
          <option value="">{t.ui.all_conditions}</option>
          {CONDITIONS.map(c => (
            <option key={c.value} value={c.value}>
              {(t.ui as any)[c.labelKey] ?? c.value.replace('_', ' ')}
            </option>
          ))}
        </select>
      </div>

      {/* Price range */}
      <div className="space-y-1.5">
        <p className={lbl}>{t.ui.price} (ETB)</p>
        <div className="flex gap-2">
          <input
            type="number"
            placeholder={t.ui.min_price}
            className={cls}
            value={minPrice}
            onChange={e => setMinPrice(e.target.value)}
            min={0}
          />
          <input
            type="number"
            placeholder={t.ui.max_price}
            className={cls}
            value={maxPrice}
            onChange={e => setMaxPrice(e.target.value)}
            min={0}
          />
        </div>
      </div>

      {/* Location */}
      <div className="space-y-1.5">
        <p className={lbl}>{t.ui.location_label}</p>
        <LocationSelector onChange={handleLocationChange} />
      </div>

      <div className="flex flex-col gap-2 pt-2">
        <button
          onClick={apply}
          className="bg-orange-600 text-white rounded-xl py-2.5 text-sm font-bold hover:bg-orange-700 transition-colors"
        >
          {t.ui.submit}
        </button>
        <button
          onClick={clear}
          className="flex items-center justify-center gap-1.5 border border-zinc-200 dark:border-zinc-700 rounded-xl py-2.5 text-sm font-bold text-zinc-600 dark:text-zinc-400 hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors"
        >
          <RotateCcw size={14} /> Clear
        </button>
      </div>
    </aside>
  );
}
