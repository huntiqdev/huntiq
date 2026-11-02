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
import { Sparkles, Check, Wand2, Lightbulb, X, RefreshCw, Layers } from 'lucide-react';
import type { AutoCategorizeResult } from '@/lib/api/public';

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
  const { t, locale } = useLocale();
  const router = useRouter();
  const [selectedType, setSelectedType] = useState<'sale' | 'rent'>(product?.listing_type ?? listingType);
  const isRental = selectedType === 'rent';
  
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

  // AI Auto-Categorize & Reasoning States
  const [aiResult, setAiResult] = useState<AutoCategorizeResult | null>(null);
  const [aiLoading, setAiLoading] = useState(false);
  const [aiApplied, setAiApplied] = useState(false);
  const [aiDismissed, setAiDismissed] = useState(false);

  // Trigger AI Auto-Categorize with human error tolerance
  const triggerAutoCategorize = useCallback(async (textToAnalyze?: string) => {
    const query = (textToAnalyze ?? form.title).trim();
    if (!query || query.length < 3) return;

    setAiLoading(true);
    try {
      const res = await api.autoCategorize({
        title: query,
        description: form.description,
        locale,
      });
      if (res) {
        setAiResult(res);
        setAiDismissed(false);
      }
    } catch (err) {
      console.warn('Auto-categorize suggestion error:', err);
    } finally {
      setAiLoading(false);
    }
  }, [form.title, form.description, locale]);

  // Debounced auto-categorize as seller types title
  useEffect(() => {
    if (product || !form.title.trim() || form.title.trim().length < 3) return;
    const timer = setTimeout(() => {
      triggerAutoCategorize(form.title);
    }, 700);
    return () => clearTimeout(timer);
  }, [form.title, product, triggerAutoCategorize]);

  // Apply AI Suggestions (category, condition, suggested price, description, listing type)
  const applyAiSuggestions = (resToApply?: AutoCategorizeResult) => {
    const res = resToApply || aiResult;
    if (!res) return;

    setForm(f => ({
      ...f,
      category_hierarchy_id: res.category_hierarchy_id || f.category_hierarchy_id,
      category_id: res.legacy_category_id ? String(res.legacy_category_id) : f.category_id,
      condition: res.condition || f.condition,
      price: !f.price && res.suggested_price_etb ? String(res.suggested_price_etb) : f.price,
      description: !f.description && res.suggested_description ? res.suggested_description : f.description,
    }));

    if (res.primary_intent === 'rent' && selectedType !== 'rent') {
      setSelectedType('rent');
      if (!rentalForm.rental_price && res.suggested_price_etb) {
        setRentalForm(rf => ({ ...rf, rental_price: String(res.suggested_price_etb) }));
      }
    } else if (res.primary_intent === 'sale' && selectedType !== 'sale') {
      setSelectedType('sale');
    }

    setAiApplied(true);
  };

  useEffect(() => {
    if (!user) return;
    api.getSellerStats(user.token).then((s: any) => {
      if (s?.tier && ['normal', 'gold', 'premium'].includes(s.tier)) {
        setSellerTier(s.tier as SellerTier);
      }
      setTierLoaded(true);
    }).catch(() => { setTierLoaded(true); });
  }, [user]);

  // Fetch applicable attributes when category_hierarchy_id changes
  useEffect(() => {
    let active = true;
    if (!form.category_hierarchy_id) {
      setApplicableAttributes(prev => (prev.length === 0 ? prev : []));
      return;
    }
    api.getCategoryAttributes(form.category_hierarchy_id)
      .then(data => {
        if (active) setApplicableAttributes(Array.isArray(data) ? data : []);
      })
      .catch(() => {
        if (active) setApplicableAttributes([]);
      });
    return () => { active = false; };
  }, [form.category_hierarchy_id]);

  const maxImages = MAX_IMAGES_PER_TIER[sellerTier] ?? 2;

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
    if (!form.category_id && !form.category_hierarchy_id) {
      setError('Please select a valid marketplace category');
      return;
    }

    setLoading(true); setError('');
    try {
      const legacyIdNum = form.category_id ? Number(form.category_id) : undefined;
      const validLegacyId = legacyIdNum && legacyIdNum >= 1 && legacyIdNum <= 11 ? legacyIdNum : undefined;

      const basePayload = {
        title: form.title,
        description: form.description,
        category_id: validLegacyId,
        category_hierarchy_id: form.category_hierarchy_id ? Number(form.category_hierarchy_id) : (validLegacyId ?? undefined),
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
      const rentalPriceNum = parseFloat(String(rentalForm.rental_price)) || 0;
      const payload = isRental ? {
        ...basePayload,
        price: rentalPriceNum,
        rental_price: rentalPriceNum,
        rental_duration_value: rentalForm.min_duration ? parseInt(rentalForm.min_duration) : 1,
        rental_duration_unit: 'month',
        deposit_amount: rentalForm.deposit ? parseFloat(String(rentalForm.deposit)) : undefined,
        available_from: rentalForm.available_from || todayStr,
      } : { ...basePayload, price: parseFloat(String(form.price)) };

      let productId: number;
      if (product) {
        await api.updateSellerProduct(product.id, payload, user.token);
        productId = product.id;

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
          if (attrPayload.length > 0 && form.category_hierarchy_id) {
            try {
              await api.saveProductAttributes(
                productId,
                { category_id: form.category_hierarchy_id, attributes: attrPayload },
                user.token
              );
            } catch { /* non-critical, don't fail the whole submission */ }
          }
        }

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
          if (attrPayload.length > 0 && form.category_hierarchy_id) {
            try {
              await api.saveProductAttributes(
                productId,
                { category_id: form.category_hierarchy_id, attributes: attrPayload },
                user.token
              );
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

  const renderAttributeInput = (attr: AttributeDefinition) => {
    const val = attributeValues[attr.slug] ?? '';
    const baseClass = 'w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-base sm:text-sm min-h-[42px] outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-100';
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
      <div className="flex bg-zinc-100 dark:bg-zinc-800 p-1.5 rounded-2xl mb-1 border border-zinc-200 dark:border-zinc-700 gap-1.5">
        <button type="button" onClick={() => setSelectedType('sale')}
          className={`flex-1 py-2.5 px-4 min-h-[44px] rounded-xl text-xs font-black uppercase tracking-wider transition-all tap-bounce active:scale-95 flex items-center justify-center gap-2 ${selectedType === 'sale' ? 'bg-orange-600 text-white shadow-md shadow-orange-600/20' : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'}`}>
          🏷️ {(t.ui as any).mode_sale ?? 'For Sale'}
        </button>
        <button type="button" onClick={() => setSelectedType('rent')}
          className={`flex-1 py-2.5 px-4 min-h-[44px] rounded-xl text-xs font-black uppercase tracking-wider transition-all tap-bounce active:scale-95 flex items-center justify-center gap-2 ${selectedType === 'rent' ? 'bg-teal-600 text-white shadow-md shadow-teal-600/20' : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'}`}>
          🏠 {(t.ui as any).mode_rent ?? 'For Rent'}
        </button>
      </div>

      <div className="space-y-2">
        <div className="flex items-end gap-2">
          <div className="flex-1">
            <Input
              required
              label={t.ui.title}
              placeholder={isRental ? 'Property title (e.g. 2 Bedroom Condominium)' : (t.ui.product_title || 'Product title (e.g. iPhone 13, Toyota Vitz, Sofa)')}
              value={form.title}
              onChange={e => {
                setForm(f => ({ ...f, title: e.target.value }));
                setAiApplied(false);
              }}
            />
          </div>
          <button
            type="button"
            onClick={() => triggerAutoCategorize()}
            disabled={aiLoading || !form.title.trim()}
            className="mb-0.5 px-3 py-2.5 min-h-[44px] rounded-xl border border-orange-500/30 bg-orange-50 dark:bg-orange-950/30 hover:bg-orange-100 dark:hover:bg-orange-900/40 text-orange-600 dark:text-orange-400 text-xs font-bold transition-all tap-bounce active:scale-95 flex items-center gap-1.5 shrink-0 shadow-sm disabled:opacity-50"
            title="Auto-detect category, price, and condition with AI"
          >
            {aiLoading ? (
              <RefreshCw className="w-3.5 h-3.5 animate-spin" />
            ) : (
              <Sparkles className="w-3.5 h-3.5" />
            )}
            <span className="hidden sm:inline">AI Auto-Detect</span>
          </button>
        </div>

        {/* AI Auto-Categorize Suggestion Banner */}
        {aiResult && !aiDismissed && (
          <div className="p-3.5 rounded-2xl border border-orange-500/30 bg-gradient-to-r from-orange-50/70 via-amber-50/50 to-orange-50/40 dark:from-orange-950/30 dark:via-zinc-900 dark:to-orange-950/20 shadow-sm animate-in fade-in-50 duration-200 space-y-2.5">
            <div className="flex items-start justify-between gap-2">
              <div className="space-y-1">
                <div className="flex items-center gap-1.5 text-xs font-black uppercase tracking-wider text-orange-600 dark:text-orange-400">
                  <Sparkles className="w-3.5 h-3.5" />
                  <span>AI Auto-Detected Branch Destination</span>
                  {aiApplied && (
                    <span className="text-[10px] bg-emerald-100 text-emerald-800 dark:bg-emerald-900/40 dark:text-emerald-300 px-2 py-0.5 rounded-full font-bold ml-1">
                      ✓ Applied
                    </span>
                  )}
                </div>
                <div className="text-xs font-bold text-zinc-900 dark:text-zinc-100 flex items-center gap-1.5 flex-wrap">
                  <span className="text-zinc-600 dark:text-zinc-300 font-normal">Path:</span>
                  <span>{aiResult.branch_path_localized || aiResult.full_branch_path}</span>
                </div>
              </div>

              <button
                type="button"
                onClick={() => setAiDismissed(true)}
                className="p-1 text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200 rounded-lg"
                title="Dismiss suggestion"
              >
                <X className="w-3.5 h-3.5" />
              </button>
            </div>

            {/* Smart tags: Leaf destination, condition, intent, suggested price */}
            <div className="flex items-center gap-2 flex-wrap text-xs">
              <span className="px-2 py-0.5 rounded-lg bg-emerald-100 dark:bg-emerald-950/40 text-emerald-800 dark:text-emerald-300 font-bold border border-emerald-300 dark:border-emerald-800">
                🎯 {aiResult.leaf_destination_name}
              </span>

              <span className="px-2 py-0.5 rounded-lg bg-blue-100 dark:bg-blue-950/40 text-blue-800 dark:text-blue-300 font-semibold">
                🏷️ {aiResult.primary_intent === 'rent' ? 'For Rent' : 'For Sale'}
              </span>

              <span className="px-2 py-0.5 rounded-lg bg-purple-100 dark:bg-purple-950/40 text-purple-800 dark:text-purple-300 font-semibold">
                📦 {aiResult.condition === 'new' ? 'New / አዲስ' : 'Used / ያገለገለ'}
              </span>

              {aiResult.suggested_price_etb > 0 && (
                <span className="px-2 py-0.5 rounded-lg bg-amber-100 dark:bg-amber-950/40 text-amber-900 dark:text-amber-300 font-semibold font-mono">
                  💰 ~{aiResult.suggested_price_etb.toLocaleString()} ETB
                </span>
              )}
            </div>

            {/* Human error correction highlight */}
            {aiResult.human_error_corrected && aiResult.original_term && (
              <div className="flex items-center gap-1.5 text-[11px] text-amber-700 dark:text-amber-300 bg-amber-500/10 px-2.5 py-1 rounded-xl">
                <Lightbulb className="w-3 h-3 text-amber-500 shrink-0" />
                <span>
                  Understood colloquial typo: <strong className="line-through">{aiResult.original_term}</strong> → <strong className="text-orange-600 dark:text-orange-400">{aiResult.corrected_term}</strong>
                </span>
              </div>
            )}

            {/* Action buttons */}
            <div className="flex items-center gap-2 pt-0.5 flex-wrap">
              {!aiApplied ? (
                <button
                  type="button"
                  onClick={() => applyAiSuggestions()}
                  className="px-3 py-2 min-h-[38px] rounded-xl bg-orange-600 hover:bg-orange-700 text-white text-xs font-bold transition-all tap-bounce active:scale-95 shadow-sm flex items-center gap-1.5"
                >
                  <Check className="w-3.5 h-3.5" />
                  <span>Apply AI Suggestions</span>
                </button>
              ) : (
                <button
                  type="button"
                  onClick={() => applyAiSuggestions()}
                  className="px-3 py-2 min-h-[38px] rounded-xl bg-zinc-200 dark:bg-zinc-800 text-zinc-700 dark:text-zinc-300 text-xs font-bold transition-all tap-bounce active:scale-95 flex items-center gap-1.5"
                >
                  <RefreshCw className="w-3 h-3" />
                  <span>Re-apply Suggestions</span>
                </button>
              )}

              {aiResult.suggested_description && (
                <button
                  type="button"
                  onClick={() => {
                    setForm(f => ({ ...f, description: aiResult.suggested_description }));
                  }}
                  className="px-3 py-2 min-h-[38px] rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 hover:bg-zinc-50 dark:hover:bg-zinc-700 text-zinc-700 dark:text-zinc-300 text-xs font-medium transition-all tap-bounce active:scale-95 flex items-center gap-1.5"
                >
                  <Wand2 className="w-3.5 h-3.5 text-orange-500" />
                  <span>Draft Description</span>
                </button>
              )}
            </div>
          </div>
        )}
      </div>

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
        <div className="flex items-center justify-between">
          <label className={lbl}>{t.ui.description}</label>
          {aiResult?.suggested_description && (
            <button
              type="button"
              onClick={() => setForm(f => ({ ...f, description: aiResult.suggested_description }))}
              className="text-xs text-orange-600 dark:text-orange-400 hover:underline flex items-center gap-1"
            >
              <Wand2 className="w-3 h-3" />
              <span>AI Auto-Draft</span>
            </button>
          )}
        </div>
        <textarea
          placeholder={t.ui.description || 'Describe the item or property...'}
          value={form.description}
          onChange={e => setForm(f => ({ ...f, description: e.target.value }))}
          rows={3}
          className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3.5 py-2.5 text-base sm:text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 transition-all resize-none text-zinc-900 dark:text-zinc-100"
        />
      </div>

      {/* Unified Category Selector */}
      <div className="flex flex-col gap-1">
        <label className={lbl}>{t.ui.selectCategory}</label>
        <CategoryTreeSelector
          value={form.category_hierarchy_id ?? (form.category_id ? Number(form.category_id) : undefined)}
          language={locale}
          onChange={(id) => {
            setForm(f => ({
              ...f,
              category_hierarchy_id: id,
              category_id: id && id <= 11 ? String(id) : '',
            }));
            setAttributeValues({});
          }}
          placeholder={t.ui.selectCategoryOption || 'Select a category...'}
        />
        {form.category_hierarchy_id && (
          <p className="text-xs text-emerald-600 dark:text-emerald-400 mt-1">
            ✓ Category selected — dynamic attributes loaded below
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

      <Button type="submit" size="lg" disabled={loading} className={`w-full mt-2 tap-bounce min-h-[48px] ${isRental ? 'bg-teal-600 hover:bg-teal-700 shadow-teal-600/20' : 'bg-orange-600 hover:bg-orange-700 shadow-orange-600/20'}`}>
        {loading ? t.ui.saving : product ? t.ui.update_product : (isRental ? 'Post for Rent' : t.ui.createProduct)}
      </Button>
    </form>
  );
}
