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
import { Check, X, RefreshCw, Layers, Globe, Sparkles, ExternalLink } from 'lucide-react';

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
  const [selectedType, setSelectedType] = useState<'sale' | 'rent' | 'announcement'>(
    product?.listing_type === 'rent' ? 'rent' : (product?.price === 1 ? 'announcement' : (listingType ?? 'sale'))
  );
  const isRental = selectedType === 'rent';
  const isAnnouncement = selectedType === 'announcement';
  const [isFreeOrInquiry, setIsFreeOrInquiry] = useState(product?.price === 1 && selectedType === 'announcement');
  
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

  const [publisherRole, setPublisherRole] = useState<'owner' | 'authorized_seller' | 'broker_agent' | 'representative' | 'government_officer'>(
    (product?.publisher_role as any) ?? 'owner'
  );
  const [brokerageDisclosureAgreed, setBrokerageDisclosureAgreed] = useState(
    product?.brokerage_disclosure_agreed ?? true
  );

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

  // Gold & Premium merchant external links
  const [customURL, setCustomURL] = useState(product?.custom_url ?? '');
  const [socialTelegram, setSocialTelegram] = useState(product?.social_links?.telegram ?? '');
  const [socialYoutube, setSocialYoutube] = useState(product?.social_links?.youtube ?? '');
  const [socialTiktok, setSocialTiktok] = useState(product?.social_links?.tiktok ?? '');
  
  const [imageURLs, setImageURLs] = useState<string[]>(
    product?.images?.map((i: any) => i.image_url) ?? []
  );
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [tierLoaded, setTierLoaded] = useState(false);

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
      const timer = setTimeout(() => {
        if (active) setApplicableAttributes(prev => (prev.length === 0 ? prev : []));
      }, 0);
      return () => { active = false; clearTimeout(timer); };
    }
    // Prioritize direct category hierarchy attributes, fallback to universal classification
    api.getCategoryAttributes(form.category_hierarchy_id)
      .then(catAttrs => {
        if (active && Array.isArray(catAttrs) && catAttrs.length > 0) {
          setApplicableAttributes(catAttrs);
          return null;
        }
        return api.getEffectiveAttributes(form.category_hierarchy_id!);
      })
      .then(classAttrs => {
        if (active && classAttrs && Array.isArray(classAttrs)) {
          setApplicableAttributes(classAttrs.map(d => ({
            attribute_id: d.attribute_id,
            name_en: d.name_en,
            slug: d.slug,
            data_type: d.data_type,
            validation_rules: d.validation_rules,
            unit: d.unit,
            is_required: d.is_required,
            is_filterable: d.is_filterable,
            display_order: d.display_order,
          })));
        }
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
    } else if (isAnnouncement) {
      if (!isFreeOrInquiry && form.price) {
        const p = parseFloat(String(form.price));
        if (isNaN(p) || p <= 0) { setError('Please enter a valid rate/salary or select Contact for Quote / Apply'); return; }
      }
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
    if (!location.region && !location.location_hierarchy_id) {
      setError(locale === 'am' ? 'እባክዎ አካባቢዎን (ክልል/ከተማ) ይምረጡ - አካባቢ የግዴታ ነው' : 'Location is mandatory. Please select your Region/City.');
      return;
    }

    setLoading(true); setError('');
    try {
      const legacyIdNum = form.category_id ? Number(form.category_id) : undefined;
      const validLegacyId = legacyIdNum && legacyIdNum > 0 ? legacyIdNum : undefined;

      const socialLinks: Record<string, string> = {};
      if (socialTelegram.trim()) socialLinks.telegram = socialTelegram.trim();
      if (socialYoutube.trim()) socialLinks.youtube = socialYoutube.trim();
      if (socialTiktok.trim()) socialLinks.tiktok = socialTiktok.trim();

      const isEligibleTier = sellerTier === 'gold' || sellerTier === 'premium';

      const basePayload = {
        title: form.title,
        description: form.description,
        category_id: validLegacyId,
        category_hierarchy_id: form.category_hierarchy_id ? Number(form.category_hierarchy_id) : (validLegacyId ?? undefined),
        placement_id: form.category_hierarchy_id ? Number(form.category_hierarchy_id) : undefined,
        subcategory: form.subcategory,
        item: form.item,
        condition: form.condition,
        custom_url: isEligibleTier && customURL.trim() ? customURL.trim() : undefined,
        social_links: isEligibleTier && Object.keys(socialLinks).length > 0 ? socialLinks : undefined,
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
        publisher_role: publisherRole,
        brokerage_disclosure_agreed: brokerageDisclosureAgreed,
      };

      const todayStr = new Date().toISOString().split('T')[0];
      const rentalPriceNum = parseFloat(String(rentalForm.rental_price)) || 0;
      const announcementPrice = isFreeOrInquiry || !form.price ? 1 : (parseFloat(String(form.price)) || 1);
      const payload = isRental ? {
        ...basePayload,
        price: rentalPriceNum,
        rental_price: rentalPriceNum,
        rental_duration_value: rentalForm.min_duration ? parseInt(rentalForm.min_duration) : 1,
        rental_duration_unit: 'month',
        deposit_amount: rentalForm.deposit ? parseFloat(String(rentalForm.deposit)) : undefined,
        available_from: rentalForm.available_from || todayStr,
      } : {
        ...basePayload,
        price: isAnnouncement ? announcementPrice : parseFloat(String(form.price)),
      };

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
    const baseClass = 'w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-base sm:text-sm min-h-[44px] outline-none focus:ring-2 focus:ring-[#1769E0] focus:border-[#1769E0] bg-zinc-50 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-100';
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
                className="rounded border-zinc-300 text-[#1769E0] focus:ring-[#1769E0]" />
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
      {/* Business Posting Mode Toggle */}
      <div className="space-y-1.5">
        <label className={lbl}>Business Posting Mode</label>
        <div className="grid grid-cols-3 bg-zinc-100 dark:bg-zinc-800 p-1.5 rounded-2xl mb-1 border border-zinc-200 dark:border-zinc-700 gap-1.5">
          <button type="button" onClick={() => setSelectedType('sale')}
            className={`py-2.5 px-3 min-h-[44px] rounded-xl text-xs font-black uppercase tracking-wider transition-all tap-bounce active:scale-95 flex items-center justify-center gap-1.5 ${selectedType === 'sale' ? 'bg-[#1769E0] text-white shadow-md shadow-[#1769E0]/20' : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'}`}>
            🛍️ {(t.ui as any).mode_sale ?? 'For Sale'}
          </button>
          <button type="button" onClick={() => setSelectedType('rent')}
            className={`py-2.5 px-3 min-h-[44px] rounded-xl text-xs font-black uppercase tracking-wider transition-all tap-bounce active:scale-95 flex items-center justify-center gap-1.5 ${selectedType === 'rent' ? 'bg-teal-600 text-white shadow-md shadow-teal-600/20' : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'}`}>
            🏠 {(t.ui as any).mode_rent ?? 'For Rent'}
          </button>
          <button type="button" onClick={() => setSelectedType('announcement')}
            className={`py-2.5 px-3 min-h-[44px] rounded-xl text-xs font-black uppercase tracking-wider transition-all tap-bounce active:scale-95 flex items-center justify-center gap-1.5 ${selectedType === 'announcement' ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/20' : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'}`}>
            📢 Announce / Info
          </button>
        </div>
      </div>

      <div className="space-y-2">
        <Input
          required
          label={t.ui.title}
          placeholder={isAnnouncement ? 'Announcement / Vacancy / Service Title (e.g. Senior Accountant Vacancy, Professional Plumbing Service)' : isRental ? 'Property or Item title (e.g. 2 Bedroom Condominium, Excavator Rental)' : (t.ui.product_title || 'Product title (e.g. iPhone 13, Toyota Vitz, Sofa)')}
          value={form.title}
          onChange={e => setForm(f => ({ ...f, title: e.target.value }))}
        />
      </div>

      {isAnnouncement ? (
        <div className="space-y-3 p-3.5 bg-indigo-50/70 dark:bg-indigo-950/30 border border-indigo-200/80 dark:border-indigo-800/60 rounded-2xl">
          <div className="flex items-center justify-between">
            <span className="text-xs font-bold text-indigo-900 dark:text-indigo-300">
              📢 Informational / Vacancy / Service Post
            </span>
            <label className="flex items-center gap-1.5 text-xs font-bold text-zinc-700 dark:text-zinc-200 cursor-pointer">
              <input
                type="checkbox"
                checked={isFreeOrInquiry}
                onChange={e => setIsFreeOrInquiry(e.target.checked)}
                className="w-4 h-4 rounded text-[#1769E0] focus:ring-[#1769E0]"
              />
              <span>Contact for Quote / Apply</span>
            </label>
          </div>
          {!isFreeOrInquiry ? (
            <Input
              type="number"
              label="Salary / Budget / Rate (ETB) (Optional)"
              placeholder="e.g. 20000 (or check 'Contact for Quote / Apply' above)"
              value={form.price}
              onChange={e => setForm(f => ({ ...f, price: e.target.value }))}
              min={1}
            />
          ) : (
            <p className="text-[11px] text-zinc-500 dark:text-zinc-400">
              Listing will be published as <strong>Contact / Free Quote / Apply</strong> without a mandatory purchase price.
            </p>
          )}
        </div>
      ) : isRental ? (
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
        <textarea
          placeholder={t.ui.description || 'Describe the item or property...'}
          value={form.description}
          onChange={e => setForm(f => ({ ...f, description: e.target.value }))}
          rows={3}
          className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3.5 py-2.5 text-base sm:text-sm outline-none focus:ring-2 focus:ring-[#1769E0] focus:border-[#1769E0] bg-zinc-50 dark:bg-zinc-800 transition-all resize-none text-zinc-900 dark:text-zinc-100"
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
              category_id: id ? String(id) : '',
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

      {!isAnnouncement && (
        <Select label={t.ui.condition} value={form.condition}
          onChange={e => setForm(f => ({ ...f, condition: e.target.value }))}>
          {CONDITIONS.map(c => <option key={c} value={c}>{condLabel(c)}</option>)}
        </Select>
      )}

      <div className="flex flex-col gap-1">
        <label className={lbl}>{t.ui.location_label}</label>
        <LocationSelector onChange={handleLocationChange} searchable showPrecise
          showCoords={true} enableGpsAutofill={true}
          defaultValue={product ? {
            region: product.location?.region, zone: product.location?.zone,
            subcity: product.location?.subcity, woreda: product.location?.woreda,
            kebele: product.kebele, specific_area: product.specific_area,
            latitude: product.latitude, longitude: product.longitude,
            location_hierarchy_id: product.location_hierarchy_id,
          } : {}} />
      </div>

      {/* Gold & Premium Merchant Features: Verified Portfolio URL & Social Channels */}
      {sellerTier === 'gold' || sellerTier === 'premium' ? (
        <div className="border border-amber-300 dark:border-amber-700/60 bg-gradient-to-br from-amber-50/60 to-yellow-50/30 dark:from-amber-950/20 dark:to-yellow-950/10 rounded-2xl p-4 space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Sparkles size={16} className="text-amber-600 dark:text-amber-400" />
              <p className="text-xs font-black uppercase tracking-wider text-amber-800 dark:text-amber-300">
                {sellerTier === 'premium' ? 'Premium Merchant Channels & Portfolio' : 'Gold Merchant Channels & Portfolio'}
              </p>
            </div>
            <span className="text-[10px] font-black uppercase px-2 py-0.5 rounded-full bg-amber-500 text-white shadow-sm">
              Unlocked
            </span>
          </div>
          <p className="text-xs text-zinc-600 dark:text-zinc-400">
            Link your official website, Telegram channel, YouTube showcase, or TikTok portfolio directly to your listing.
          </p>
          <div className="space-y-2.5">
            <Input
              type="url"
              label="Official Website or Portfolio URL"
              placeholder="https://yourwebsite.com or https://portfolio.et"
              value={customURL}
              onChange={e => setCustomURL(e.target.value)}
            />
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-2.5">
              <Input
                type="text"
                label="Telegram (@channel or link)"
                placeholder="@mychannel or t.me/..."
                value={socialTelegram}
                onChange={e => setSocialTelegram(e.target.value)}
              />
              <Input
                type="text"
                label="YouTube Channel"
                placeholder="@youtubechannel"
                value={socialYoutube}
                onChange={e => setSocialYoutube(e.target.value)}
              />
              <Input
                type="text"
                label="TikTok Handle"
                placeholder="@tiktokhandle"
                value={socialTiktok}
                onChange={e => setSocialTiktok(e.target.value)}
              />
            </div>
          </div>
        </div>
      ) : (
        <div className="border border-dashed border-zinc-300 dark:border-zinc-700 rounded-2xl p-4 bg-zinc-50/50 dark:bg-zinc-800/20 flex items-center justify-between gap-3">
          <div className="space-y-0.5">
            <div className="flex items-center gap-1.5 text-xs font-bold text-zinc-700 dark:text-zinc-300">
              <Sparkles size={14} className="text-amber-500" />
              <span>Gold & Premium Exclusive: Custom URL & Social Channels</span>
            </div>
            <p className="text-[11px] text-zinc-500 dark:text-zinc-400">
              Showcase external websites, verified TikTok/YouTube links, and direct Telegram channels on your listings.
            </p>
          </div>
          <button
            type="button"
            onClick={() => router.push('/pricing')}
            className="shrink-0 text-xs font-bold text-amber-700 dark:text-amber-300 bg-amber-100 hover:bg-amber-200 dark:bg-amber-900/40 dark:hover:bg-amber-900/60 px-3 py-1.5 rounded-xl transition-colors"
          >
            Upgrade →
          </button>
        </div>
      )}

      <ListingFeeSlider price={Number(isRental ? rentalForm.rental_price : form.price) || 0}
        tier={sellerTier} selectedDays={listingDays} onChange={setListingDays} />

      {/* Brokerage & Publisher Classification */}
      <div className="space-y-3 p-4 rounded-2xl bg-zinc-50 dark:bg-zinc-800/50 border border-zinc-200 dark:border-zinc-700/60">
        <label className="block text-xs font-black uppercase tracking-wider text-zinc-700 dark:text-zinc-300">
          Publisher Role & Intermediary Disclosure
        </label>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-2">
          {[
            { id: 'owner', label: 'Direct Owner' },
            { id: 'authorized_seller', label: 'Authorized Merchant' },
            { id: 'broker_agent', label: 'Broker / Agent' },
            { id: 'representative', label: 'Representative' },
          ].map(r => (
            <button
              key={r.id}
              type="button"
              onClick={() => {
                setPublisherRole(r.id as any);
                if (r.id === 'broker_agent' || r.id === 'representative') {
                  setBrokerageDisclosureAgreed(true);
                }
              }}
              className={`p-2.5 rounded-xl text-xs font-bold border text-center transition-all ${
                publisherRole === r.id
                  ? 'border-[#1769E0] bg-[#1769E0]/10 dark:bg-[#1769E0]/20 text-[#1769E0] dark:text-[#38A5FF] shadow-xs'
                  : 'border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-400 hover:border-zinc-300'
              }`}
            >
              {r.label}
            </button>
          ))}
        </div>
        {publisherRole === 'broker_agent' && (
          <div className="p-3 rounded-xl bg-blue-500/10 border border-blue-500/30 text-blue-900 dark:text-blue-200 text-xs">
            <span className="font-bold">Intermediary Transparency:</span> This listing will display a verified <span className="font-mono font-black">[BROKER]</span> badge and standard marketplace brokerage disclosure.
          </div>
        )}
      </div>

      <ImageUploader urls={imageURLs} onChange={setImageURLs} maxImages={maxImages} />

      {error && <p className="text-red-500 text-xs bg-red-50 dark:bg-red-900/20 px-3 py-2 rounded-xl">{error}</p>}

      <Button type="submit" size="lg" disabled={loading} className={`w-full mt-2 tap-bounce min-h-[48px] ${isRental ? 'bg-teal-600 hover:bg-teal-700 shadow-teal-600/20' : 'bg-[#1769E0] hover:bg-[#1055BD] shadow-[#1769E0]/20'}`}>
        {loading ? t.ui.saving : product ? t.ui.update_product : (isRental ? 'Post for Rent' : t.ui.createProduct)}
      </Button>
    </form>
  );
}
