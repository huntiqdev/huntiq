'use client';
import { useState, useEffect, useCallback } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { ProductFilters } from '@/types';
import { LocationSelector, LocationValue } from '@/components/location/LocationSelector';
import { useLocale } from '@/hooks/useLocale';
import { SlidersHorizontal, RotateCcw } from 'lucide-react';

interface Props {
  onFilter: (f: ProductFilters) => void;
}

const CONDITIONS = [
  { value: 'new', labelKey: 'condition_new' },
  { value: 'used', labelKey: 'condition_used' },
] as const;

export function FilterSidebar({ onFilter }: Props) {
  const { t } = useLocale();
  const router = useRouter();
  const searchParams = useSearchParams();

  const [minPrice, setMinPrice] = useState('');
  const [maxPrice, setMaxPrice] = useState('');
  const [condition, setCondition] = useState('');
  const [listingType, setListingType] = useState<'all' | 'sale' | 'rent'>('all');
  const [locationQuery, setLocationQuery] = useState<LocationValue>({});

  // Sync with searchParams if URL changes externally
  useEffect(() => {
    const type = searchParams.get('type');
    if (type === 'sale' || type === 'rent' || type === 'all') {
      setListingType(type);
    }
    const cond = searchParams.get('condition');
    if (cond) setCondition(cond);
    const minP = searchParams.get('min_price');
    if (minP) setMinPrice(minP);
    const maxP = searchParams.get('max_price');
    if (maxP) setMaxPrice(maxP);
  }, [searchParams]);

  const handleLocationChange = useCallback((v: LocationValue) => {
    setLocationQuery(v);
  }, []);

  const apply = () => {
    const f: ProductFilters = {};
    if (listingType !== 'all') f.listing_type = listingType;
    if (minPrice) f.min_price = Number(minPrice);
    if (maxPrice) f.max_price = Number(maxPrice);
    if (condition) f.condition = condition;
    if (locationQuery.location_hierarchy_id) f.location_hierarchy_id = locationQuery.location_hierarchy_id;
    if (locationQuery.woreda) f.woreda = locationQuery.woreda;
    else if (locationQuery.subcity) f.subcity = locationQuery.subcity;
    else if (locationQuery.zone) f.zone = locationQuery.zone;
    else if (locationQuery.region) f.region = locationQuery.region;
    onFilter(f);

    const p = new URLSearchParams(searchParams.toString());
    if (listingType !== 'all') p.set('type', listingType);
    else p.delete('type');

    if (minPrice) p.set('min_price', minPrice);
    else p.delete('min_price');

    if (maxPrice) p.set('max_price', maxPrice);
    else p.delete('max_price');

    if (condition) p.set('condition', condition);
    else p.delete('condition');

    if (f.location_hierarchy_id) p.set('location_hierarchy_id', String(f.location_hierarchy_id));
    if (f.woreda) p.set('woreda', f.woreda);
    else if (f.subcity) p.set('subcity', f.subcity);
    else if (f.zone) p.set('zone', f.zone);
    else if (f.region) p.set('region', f.region);

    router.push(`/products${p.toString() ? '?' + p.toString() : ''}`);
  };

  const clear = () => {
    setListingType('all');
    setMinPrice('');
    setMaxPrice('');
    setCondition('');
    setLocationQuery({} as LocationValue);
    onFilter({});

    const p = new URLSearchParams(searchParams.toString());
    p.delete('type');
    p.delete('min_price');
    p.delete('max_price');
    p.delete('condition');
    p.delete('location_hierarchy_id');
    p.delete('region');
    p.delete('zone');
    p.delete('subcity');
    p.delete('woreda');
    router.push(`/products${p.toString() ? '?' + p.toString() : ''}`);
  };

  const cls = 'w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-[#1769E0] focus:border-[#1769E0] bg-zinc-50 dark:bg-zinc-800 transition-all';
  const lbl = 'text-xs font-black text-zinc-500 uppercase tracking-wider';

  return (
    <aside className="bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 p-4 flex flex-col gap-4 h-fit">
      <div className="flex items-center gap-2">
        <SlidersHorizontal size={16} className="text-[#1769E0]" />
        <h2 className="font-bold">{t.ui.filters}</h2>
      </div>

      {/* Dimension 1: Marketplace Pillar (Transaction Mode) */}
      <div className="space-y-1.5 bg-zinc-50 dark:bg-zinc-800/60 p-2.5 rounded-2xl border border-zinc-200 dark:border-zinc-700/60">
        <p className={lbl}>{(t.ui as any)?.marketplace_pillar ?? 'Marketplace'}</p>
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
            {t.ui?.condition_all ?? 'All'}
          </button>
          <button
            type="button"
            onClick={() => setListingType('sale')}
            className={`py-1.5 text-[11px] font-black uppercase tracking-wider rounded-lg transition-all ${
              listingType === 'sale'
                ? 'bg-[#1769E0] text-white shadow-sm'
                : 'text-zinc-500 hover:text-[#1769E0]'
            }`}
          >
            🛒 {t.nav?.buy ?? 'Buy'}
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
            🏠 {t.nav?.rental ?? 'Rent'}
          </button>
        </div>
      </div>

      {/* Dimension 2: Item Condition */}
      <div className="space-y-1.5">
        <p className={lbl}>{t.ui.condition}</p>
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
          className="bg-[#1769E0] text-white rounded-xl py-2.5 text-sm font-bold hover:bg-[#1055BD] transition-colors"
        >
          {t.ui.submit}
        </button>
        <button
          onClick={clear}
          className="flex items-center justify-center gap-1.5 border border-zinc-200 dark:border-zinc-700 rounded-xl py-2.5 text-sm font-bold text-zinc-600 dark:text-zinc-400 hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors"
        >
          <RotateCcw size={14} /> {t.ui.clear_filter ?? 'Clear'}
        </button>
      </div>
    </aside>
  );
}
