'use client';
import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '@/context/AuthContext';
import { useLocale } from '@/hooks/useLocale';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import { Button } from '@/components/ui/Button';
import { LocationSelector, type LocationValue } from '@/components/location/LocationSelector';
import { ListingFeeSlider } from '@/components/seller/ListingFeeSlider';
import { ImageUploader } from '@/components/seller/ImageUploader';
import { type SellerTier, MAX_IMAGES_PER_TIER } from '@/lib/constants/fees';
import { SELECT_CLS as cls, LABEL_CLS as lbl } from '@/lib/constants/styles';
import CategoryTreeSelector from '@/components/category/CategoryTreeSelector';

interface AttributeDefinition {
  attribute_id: number;
  name_en: string;
  name_am?: string;
  slug: string;
  data_type: string;
  validation_rules?: Record<string, any>;
  unit?: string;
  is_required: boolean;
  display_order: number;
}

interface Props { 
  product?: any; 
  onSuccess: () => void; 
  initialCondition?: string;
  listingType?: 'sale' | 'rent';
}

const CONDITIONS = ['new', 'used'] as const;

export function ProductForm({ product, onSuccess, initialCondition, listingType = 'sale' }: Props) {
  const { user } = useAuth();
  const { t } = useLocale();
  const router = useRouter();
  const [selectedType, setSelectedType] = useState<'sale' | 'rent'>(product?.listing_type ?? listingType);
  const isRental = selectedType === 'rent';
  
  const [categories, setCategories] = useState<any[]>([]);
  const [sellerTier, setSellerTier] = useState<SellerTier>('normal');
  const [location, setLocation] = useState<LocationValue>(product ? {
    region: product.location?.region,
    zone: product.location?.zone,
    subcity: product.location?.subcity,
    woreda: product.location?.woreda,
    kebele: product.kebele,
    specific_area: product.specific_area,
    latitude: product.latitude,
    longitude: product.longitude,
  } : {});
  const handleLocationChange = useCallback((v: LocationValue) => setLocation(v), []);
  const [listingDays, setListingDays] = useState(product?.listing_duration_days ?? 30);
  const [form, setForm] = useState({
    title: product?.title ?? '',
    price: product?.price ?? '',
    description: product?.description ?? '',
    category_id: product?.category_id ?? product?.category?.id ?? '',
    category_hierarchy_id: product?.category_hierarchy_id ?? null as number | null,
    subcategory: product?.subcategory ?? '',
    item: product?.item ?? '',
    condition: product?.condition ?? initialCondition ?? 'used',
  });

  // Dynamic attributes per category
  const [applicableAttributes, setApplicableAttributes] = useState<AttributeDefinition[]>([]);
  const [attributeValues, setAttributeValues] = useState<Record<string, any>>({});
  
  const [rentalForm, setRentalForm] = useState({
    rental_price: product?.rental_price ?? '',
    min_duration: product?.rental_duration_value ?? '',
    max_duration: product?.rental_max_duration_months ?? '',
    deposit: product?.deposit_amount ?? '',
    available_from: product?.available_from ? String(product.available_from).split('T')[0] : '',
    rental_terms: product?.rental_terms ?? '',
  });
  
  const [imageURLs, setImageURLs] = useState<string[]>(
    product?.images?.map((i: any) => i.image_url) ?? []
  );
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [tierLoaded, setTierLoaded] = useState(false);

  useEffect(() => { api.getCategories().then(setCategories).catch(() => {}); }, []);
  useEffect(() => {
    if (!user) return;
    api.getSellerStats(user.token).then((s: any) => {
      if (s?.tier && ['normal', 'gold', 'premium'].includes(s.tier)) {
        setSellerTier(s.tier as SellerTier);
      }
      setTierLoaded(true);
    }).catch(() => { setTierLoaded(true); });
  }, [user]);

  useEffect(() => {
    if (!tierLoaded) return;
    if (product) return;
    if (imageURLs.some(u => u.trim())) return;
    setImageURLs([]);
  }, [tierLoaded]); // eslint-disable-line react-hooks/exhaustive-deps

  // Fetch applicable attributes when category_hierarchy_id changes
  useEffect(() => {
    if (!form.category_hierarchy_id) {
      setApplicableAttributes([]);
      return;
    }
    api.getCategoryAttributes(form.category_hierarchy_id)
      .then(data => setApplicableAttributes(Array.isArray(data) ? data : []))
      .catch(() => setApplicableAttributes([]));
  }, [form.category_hierarchy_id]);

  const maxImages = MAX_IMAGES_PER_TIER[sellerTier] ?? 2;

  useEffect(() => {
    if (imageURLs.length > maxImages) setImageURLs(prev => prev.slice(0, maxImages));
  }, [maxImages, imageURLs.length]);

  const setAttrValue = (slug: string, value: any) => {
    setAttributeValues(prev => ({ ...prev, [slug]: value }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!user) return;

    if (isRental) {
      const rp = parseFloat(String(rentalForm.rental_price));
      if (!rentalForm.rental_price || isNaN(rp) || rp <= 0) { setError('A valid rental price is required'); return; }
    } else {
      const p = parseFloat(String(form.price));
      if (!form.price || isNaN(p) || p <= 0) { setError('A valid positive price is required'); return; }
    }
    if (imageURLs.length < 2) { setError('Minimum 2 images are required'); return; }
    if (imageURLs.length > maxImages) { setError(`Your tier allows max ${maxImages} images`); return; }

    setLoading(true); setError('');
    try {
      const basePayload = {
        title: form.title,
        description: form.description,
        category_id: form.category_id ? Number(form.category_id) : undefined,
        category_hierarchy_id: form.category_hierarchy_id ? Number(form.category_hierarchy_id) : (form.category_id ? Number(form.category_id) : undefined),
        subcategory: form.subcategory,
        item: form.item,
        condition: form.condition,
        image_urls: imageURLs,
        region: location.region,
        zone: location.zone,
        subcity: location.subcity,
        woreda: location.woreda,
        kebele: location.kebele,
        specific_area: location.specific_area,
        latitude: location.latitude,
        longitude: location.longitude,
        location_hierarchy_id: location.location_hierarchy_id,
        listing_duration_days: listingDays,
        listing_type: isRental ? 'rent' : 'sale',
      };

      const todayStr = new Date().toISOString().split('T')[0];
      const payload = isRental ? {
        ...basePayload, price: 0,
        rental_price: parseFloat(String(rentalForm.rental_price)),
        rental_duration_value: rentalForm.min_duration ? parseInt(rentalForm.min_duration) : 1,
        rental_duration_unit: 'month',
        deposit_amount: rentalForm.deposit ? parseFloat(String(rentalForm.deposit)) : undefined,
        available_from: rentalForm.available_from || todayStr,
      } : { ...basePayload, price: parseFloat(String(form.price)) };

      let productId: number;
      if (product) {
        await api.updateSellerProduct(product.id, payload, user.token);
        productId = product.id;
        onSuccess();
      } else {
        const created = await api.createSellerProduct(payload, user.token);
        productId = created.id;

        // Save dynamic attributes if any were filled
        if (form.category_hierarchy_id && Object.keys(attributeValues).length > 0) {
          const attrPayload = applicableAttributes
            .filter(attr => attributeValues[attr.slug] !== undefined && attributeValues[attr.slug] !== '')
            .map(attr => {
              const val = attributeValues[attr.slug];
              const base: any = { attribute_id: attr.attribute_id };
              if (attr.data_type === 'number' || attr.data_type === 'range') base.value_number = parseFloat(val);
              else if (attr.data_type === 'boolean') base.value_boolean = val === true || val === 'true';
              else base.value_text = String(val);
              return base;
            });
          if (attrPayload.length > 0) {
            try {
              await fetch(`/api/seller/products/${productId}/attributes`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${user.token}` },
                body: JSON.stringify({ category_id: form.category_hierarchy_id, attributes: attrPayload }),
              });
            } catch { /* non-critical, don't fail the whole submission */ }
          }
        }

        const fee = Number(created?.total_fee_amount) || 0;
        if (fee > 0) {
          router.push(`/seller/pay?product_id=${productId}&fee=${fee}&title=${encodeURIComponent(basePayload.title)}`);
        } else {
          onSuccess();
        }
      }
    } catch (err: any) { setError(err.message); }
    finally { setLoading(false); }
  };

  const condLabel = (c: string) => ({ new: t.ui.condition_new ?? 'New', used: t.ui.condition_used ?? 'Used' }[c] ?? c);
  const label = (ns: 'categories' | 'items', key: string) => (t[ns] as any)?.[key] ?? key.replace(/_/g, ' ');

  const renderAttributeInput = (attr: AttributeDefinition) => {
    const val = attributeValues[attr.slug] ?? '';
    const baseClass = 'w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800';
    const labelText = attr.name_en + (attr.is_required ? ' *' : '') + (attr.unit ? ` (${attr.unit})` : '');

    switch (attr.data_type) {
      case 'select': {
        const options = attr.validation_rules?.options ?? [];
        return (
          <div key={attr.attribute_id} className="flex flex-col gap-1">
            <label className={lbl}>{labelText}</label>
            <select className={baseClass} value={val} onChange={e => setAttrValue(attr.slug, e.target.value)}>
              <option value="">— Select —</option>
              {options.map((o: string) => <option key={o} value={o}>{o}</option>)}
            </select>
          </div>
        );
      }
      case 'boolean':
        return (
          <div key={attr.attribute_id} className="flex items-center gap-3 py-1">
            <label className="flex items-center gap-2 cursor-pointer">
              <input type="checkbox" checked={val === true || val === 'true'}
                onChange={e => setAttrValue(attr.slug, e.target.checked)}
                className="rounded border-zinc-300 text-orange-600 focus:ring-orange-500" />
              <span className="text-sm text-zinc-700 dark:text-zinc-300">{labelText}</span>
            </label>
          </div>
        );
      case 'number':
      case 'range': {
        const min = attr.validation_rules?.min;
        const max = attr.validation_rules?.max;
        return (
          <div key={attr.attribute_id} className="flex flex-col gap-1">
            <label className={lbl}>{labelText}</label>
            <input type="number" className={baseClass} value={val}
              onChange={e => setAttrValue(attr.slug, e.target.value)}
              min={min} max={max} placeholder={`${min !== undefined ? `Min: ${min}` : ''}${max !== undefined ? ` Max: ${max}` : ''}`} />
          </div>
        );
      }
      default:
        return (
          <div key={attr.attribute_id} className="flex flex-col gap-1">
            <label className={lbl}>{labelText}</label>
            <input type="text" className={baseClass} value={val}
              onChange={e => setAttrValue(attr.slug, e.target.value)}
              placeholder={attr.name_en} />
          </div>
        );
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-3">
      {/* Listing type toggle */}
      <div className="flex bg-zinc-100 dark:bg-zinc-800 p-1.5 rounded-2xl mb-1 border border-zinc-200 dark:border-zinc-700">
        <button type="button" onClick={() => setSelectedType('sale')}
          className={`flex-1 py-2.5 px-4 rounded-xl text-xs font-black uppercase tracking-wider transition-all flex items-center justify-center gap-2 ${selectedType === 'sale' ? 'bg-orange-600 text-white shadow-md shadow-orange-600/20' : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'}`}>
          🏷️ {(t.ui as any).mode_sale ?? 'For Sale'}
        </button>
        <button type="button" onClick={() => setSelectedType('rent')}
          className={`flex-1 py-2.5 px-4 rounded-xl text-xs font-black uppercase tracking-wider transition-all flex items-center justify-center gap-2 ${selectedType === 'rent' ? 'bg-teal-600 text-white shadow-md shadow-teal-600/20' : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'}`}>
          🏠 {(t.ui as any).mode_rent ?? 'For Rent'}
        </button>
      </div>

      <Input required label={t.ui.title} placeholder={isRental ? 'Property title' : t.ui.product_title}
        value={form.title} onChange={e => setForm(f => ({ ...f, title: e.target.value }))} />

      {isRental ? (
        <>
          <Input required type="number" label="Monthly Rent (ETB)" placeholder="Rental price per month"
            value={rentalForm.rental_price} onChange={e => setRentalForm(f => ({ ...f, rental_price: e.target.value }))} min={1} />
          <div className="grid grid-cols-2 gap-3">
            <Input type="number" label="Min Duration (months)" placeholder="e.g., 3"
              value={rentalForm.min_duration} onChange={e => setRentalForm(f => ({ ...f, min_duration: e.target.value }))} min={1} />
            <Input type="number" label="Max Duration (months)" placeholder="e.g., 12"
              value={rentalForm.max_duration} onChange={e => setRentalForm(f => ({ ...f, max_duration: e.target.value }))} min={1} />
          </div>
          <Input type="number" label="Security Deposit (ETB)" placeholder="Optional"
            value={rentalForm.deposit} onChange={e => setRentalForm(f => ({ ...f, deposit: e.target.value }))} min={0} />
          <Input type="date" label="Available From" value={rentalForm.available_from}
            onChange={e => setRentalForm(f => ({ ...f, available_from: e.target.value }))} />
        </>
      ) : (
        <Input required type="number" label={`${t.ui.price} (ETB)`} placeholder="Price in ETB"
          value={form.price} onChange={e => setForm(f => ({ ...f, price: e.target.value }))} min={1} />
      )}

      <div className="flex flex-col gap-1">
        <label className={lbl}>{t.ui.description}</label>
        <textarea placeholder={t.ui.description} value={form.description}
          onChange={e => setForm(f => ({ ...f, description: e.target.value }))}
          rows={3} className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 transition-all resize-none" />
      </div>

      {/* Category selector — legacy flat categories for backward compat */}
      <div className="flex flex-col gap-1">
        <label className={lbl}>{t.ui.selectCategory}</label>
        <select className={cls} value={form.category_id}
          onChange={e => setForm(f => ({ ...f, category_id: e.target.value, subcategory: '', item: '' }))}>
          <option value="">{t.ui.selectCategoryOption}</option>
          {categories.map(c => <option key={c.id} value={c.id}>{label('categories', c.slug)}</option>)}
        </select>
      </div>

      {/* New hierarchical category selector */}
      <div className="flex flex-col gap-1">
        <label className={lbl}>Specific Category (Hierarchical)</label>
        <CategoryTreeSelector
          value={form.category_hierarchy_id ?? undefined}
          onChange={(id) => {
            setForm(f => ({ ...f, category_hierarchy_id: id }));
            setAttributeValues({});
          }}
          placeholder="Select specific subcategory..."
        />
        {form.category_hierarchy_id && (
          <p className="text-xs text-emerald-600 dark:text-emerald-400 mt-1">
            ✓ Hierarchical category selected — dynamic attributes loaded below
          </p>
        )}
      </div>

      {/* Dynamic Attributes */}
      {applicableAttributes.length > 0 && (
        <div className="border border-zinc-200 dark:border-zinc-700 rounded-2xl p-4 space-y-3 bg-zinc-50/50 dark:bg-zinc-800/30">
          <p className="text-xs font-black uppercase tracking-wider text-zinc-500">
            Product Specifications ({applicableAttributes.length} attributes)
          </p>
          {applicableAttributes
            .sort((a, b) => a.display_order - b.display_order)
            .map(attr => renderAttributeInput(attr))}
        </div>
      )}

      <Select label={t.ui.condition} value={form.condition}
        onChange={e => setForm(f => ({ ...f, condition: e.target.value }))}>
        {CONDITIONS.map(c => <option key={c} value={c}>{condLabel(c)}</option>)}
      </Select>

      <div className="flex flex-col gap-1">
        <label className={lbl}>{t.ui.location_label}</label>
        <LocationSelector onChange={handleLocationChange} searchable showPrecise
          showCoords={sellerTier === 'gold' || sellerTier === 'premium'}
          defaultValue={product ? {
            region: product.location?.region, zone: product.location?.zone,
            subcity: product.location?.subcity, woreda: product.location?.woreda,
            kebele: product.kebele, specific_area: product.specific_area,
            latitude: product.latitude, longitude: product.longitude,
            location_hierarchy_id: product.location_hierarchy_id,
          } : {}} />
      </div>

      <ListingFeeSlider price={Number(isRental ? rentalForm.rental_price : form.price) || 0}
        tier={sellerTier} selectedDays={listingDays} onChange={setListingDays} />

      <ImageUploader urls={imageURLs} onChange={setImageURLs} maxImages={maxImages} />

      {error && <p className="text-red-500 text-xs bg-red-50 dark:bg-red-900/20 px-3 py-2 rounded-xl">{error}</p>}

      <Button type="submit" disabled={loading} className={`w-full mt-1 ${isRental ? 'bg-teal-600 hover:bg-teal-700' : ''}`}>
        {loading ? t.ui.saving : product ? t.ui.update_product : (isRental ? 'Post for Rent' : t.ui.createProduct)}
      </Button>
    </form>
  );
}
