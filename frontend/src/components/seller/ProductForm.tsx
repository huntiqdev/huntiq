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
import CategoryTreeSelector, { type CategoryNode } from '@/components/category/CategoryTreeSelector';
import { Check, X, RefreshCw, Layers, Globe, Sparkles, ExternalLink, CheckCircle2, Tag, ChevronDown, ChevronUp, Search, Link as LinkIcon } from 'lucide-react';

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

interface CanonicalCategoryPreset {
  id: number;
  slug: string;
  nameEn: string;
  nameAm: string;
  icon: string;
  activeBorder: string;
  keywords: string[];
  registeredItems: {
    title: string;
    subcategory: string;
    item: string;
  }[];
}

const CANONICAL_PRESETS: CanonicalCategoryPreset[] = [
  {
    id: 21,
    slug: 'phones',
    nameEn: 'Phones & Tablets',
    nameAm: 'ስልኮች እና ታብሌቶች',
    icon: '📱',
    activeBorder: 'border-blue-500 bg-blue-50/70 dark:bg-blue-950/40 ring-2 ring-blue-500/20 text-blue-700 dark:text-blue-300',
    keywords: ['phone', 'phones', 'iphone', 'samsung', 'mobile', 'cellphone', 'tablet', 'ipad', 'charger', 'power bank', 'ስልክ', 'ሞባይል', 'ታብሌት', 'mobayila', 'bilbila', 'accessories'],
    registeredItems: [
      { title: 'Apple iPhone 15 Pro Max', subcategory: 'Smartphones', item: 'Apple iPhone 15 Pro Max' },
      { title: 'Apple iPhone 14 Pro', subcategory: 'Smartphones', item: 'Apple iPhone 14 Pro' },
      { title: 'Apple iPhone 13', subcategory: 'Smartphones', item: 'Apple iPhone 13' },
      { title: 'Samsung Galaxy S24 Ultra', subcategory: 'Smartphones', item: 'Samsung Galaxy S24 Ultra' },
      { title: 'Samsung Galaxy A54 5G', subcategory: 'Smartphones', item: 'Samsung Galaxy A54' },
      { title: 'Xiaomi Redmi Note 13 Pro', subcategory: 'Smartphones', item: 'Xiaomi Redmi Note 13' },
      { title: 'Tecno Spark 20 Pro', subcategory: 'Smartphones', item: 'Tecno Spark 20' },
      { title: 'Tecno Camon 30', subcategory: 'Smartphones', item: 'Tecno Camon 30' },
      { title: 'Infinix Hot 40 Pro', subcategory: 'Smartphones', item: 'Infinix Hot 40' },
      { title: 'Apple iPad Air / Pro', subcategory: 'Tablets', item: 'Apple iPad' },
      { title: 'Original Fast Charger & Type-C Cable', subcategory: 'Accessories', item: 'Phone Charger' },
      { title: 'Power Bank 20,000mAh Fast Charge', subcategory: 'Accessories', item: 'Power Bank' },
    ],
  },
  {
    id: 18,
    slug: 'computers',
    nameEn: 'Laptops & Computers',
    nameAm: 'ላፕቶፖች እና ኮምፒውተሮች',
    icon: '💻',
    activeBorder: 'border-cyan-500 bg-cyan-50/70 dark:bg-cyan-950/40 ring-2 ring-cyan-500/20 text-cyan-700 dark:text-cyan-300',
    keywords: ['computer', 'laptop', 'macbook', 'dell', 'hp', 'lenovo', 'desktop', 'monitor', 'pc', 'ssd', 'ram', 'ላፕቶፕ', 'ኮምፒውተር', 'kompiitara'],
    registeredItems: [
      { title: 'Apple MacBook Pro M3', subcategory: 'Laptops', item: 'Apple MacBook Pro' },
      { title: 'Apple MacBook Air M2', subcategory: 'Laptops', item: 'Apple MacBook Air' },
      { title: 'Dell XPS 15 / Inspiron', subcategory: 'Laptops', item: 'Dell Laptop' },
      { title: 'HP Pavilion 15 / EliteBook', subcategory: 'Laptops', item: 'HP Laptop' },
      { title: 'Lenovo ThinkPad X1 Carbon', subcategory: 'Laptops', item: 'Lenovo ThinkPad' },
      { title: 'Gaming PC Desktop Full Set', subcategory: 'Desktops', item: 'Gaming PC' },
      { title: 'Computer Monitor 27-inch FHD', subcategory: 'Displays', item: 'PC Monitor' },
      { title: 'External 1TB Portable SSD/HDD', subcategory: 'Storage', item: 'External Storage' },
    ],
  },
  {
    id: 76,
    slug: 'vehicles',
    nameEn: 'Vehicles & Automotive',
    nameAm: 'ተሽከርካሪዎች እና መኪናዎች',
    icon: '🚗',
    activeBorder: 'border-amber-500 bg-amber-50/70 dark:bg-amber-950/40 ring-2 ring-amber-500/20 text-amber-700 dark:text-amber-300',
    keywords: ['car', 'cars', 'vehicle', 'vehicles', 'toyota', 'vitz', 'corolla', 'rav4', 'automobile', 'motorcycle', 'bajaj', 'truck', 'መኪና', 'ተሽከርካሪ', 'ባጃጅ', 'ሞተር', 'konkolaataa'],
    registeredItems: [
      { title: 'Toyota Vitz Compact Car', subcategory: 'Automobiles', item: 'Toyota Vitz' },
      { title: 'Toyota Corolla Executive', subcategory: 'Automobiles', item: 'Toyota Corolla' },
      { title: 'Toyota RAV4 / Land Cruiser', subcategory: 'Automobiles', item: 'Toyota SUV' },
      { title: 'Hyundai Tucson / Creta', subcategory: 'Automobiles', item: 'Hyundai Tucson' },
      { title: 'Suzuki Swift / Dzire', subcategory: 'Automobiles', item: 'Suzuki Swift' },
      { title: 'Bajaj / TVS King Three-Wheeler', subcategory: 'Motorcycles', item: 'Bajaj Three-Wheeler' },
      { title: 'Motorcycle Boxer / TVS 125', subcategory: 'Motorcycles', item: 'Motorcycle' },
      { title: 'Isuzu NPR Commercial Truck', subcategory: 'Commercial', item: 'Isuzu Truck' },
      { title: 'Auto Spare Parts & Brand New Tyres', subcategory: 'Parts', item: 'Auto Spare Parts' },
    ],
  },
  {
    id: 15,
    slug: 'clothing',
    nameEn: 'Clothing & Fashion',
    nameAm: 'አልባሳት እና ፋሽን',
    icon: '👗',
    activeBorder: 'border-rose-500 bg-rose-50/70 dark:bg-rose-950/40 ring-2 ring-rose-500/20 text-rose-700 dark:text-rose-300',
    keywords: ['clothing', 'cloth', 'clothes', 'fashion', 'habesha', 'kemis', 'dress', 'suit', 'jeans', 'shirt', 'shoes', 'sneakers', 'jacket', 'handbag', 'አልባሳት', 'ቀሚስ', 'ሱፍ', 'ጫማ', 'uffata'],
    registeredItems: [
      { title: 'Habesha Traditional Dress (የሀበሻ ቀሚስ)', subcategory: 'Traditional', item: 'Habesha Kemis' },
      { title: "Men's Italian Formal Suit & Blazer", subcategory: 'Men Fashion', item: 'Formal Suit' },
      { title: "Men's Casual Shirt & Denim Jeans", subcategory: 'Men Fashion', item: 'Casual Wear' },
      { title: "Women's Modern Dress & Skirt", subcategory: 'Women Fashion', item: 'Women Dress' },
      { title: 'Nike / Adidas Sports Sneakers', subcategory: 'Footwear', item: 'Sports Sneakers' },
      { title: 'Genuine Ethiopian Leather Shoes', subcategory: 'Footwear', item: 'Leather Shoes' },
      { title: 'Warm Leather Jacket & Winter Coat', subcategory: 'Outerwear', item: 'Leather Jacket' },
      { title: 'Designer Handbag & Travel Backpack', subcategory: 'Bags', item: 'Handbag' },
    ],
  },
  {
    id: 9,
    slug: 'books',
    nameEn: 'Books & Education',
    nameAm: 'መጽሐፍት እና ትምህርት',
    icon: '📚',
    activeBorder: 'border-emerald-500 bg-emerald-50/70 dark:bg-emerald-950/40 ring-2 ring-emerald-500/20 text-emerald-700 dark:text-emerald-300',
    keywords: ['book', 'books', 'education', 'textbook', 'history', 'novel', 'stationery', 'መጽሐፍ', 'ደብተር', 'ትምህርት', 'kitaaba', 'kitaabota'],
    registeredItems: [
      { title: 'Ethiopian History & Heritage Books', subcategory: 'History', item: 'History Book' },
      { title: 'High School & University Textbooks', subcategory: 'Academic', item: 'Textbook' },
      { title: 'English / Amharic Learning & Dictionary', subcategory: 'Languages', item: 'Language Book' },
      { title: 'Business, Economics & Investing Guide', subcategory: 'Business', item: 'Finance Book' },
      { title: 'Self-Help, Leadership & Psychology Books', subcategory: 'Self-Help', item: 'Self-Help Book' },
      { title: 'Spiritual, Religious & Wisdom Books', subcategory: 'Religious', item: 'Spiritual Book' },
      { title: "Children's Illustrated Story & Reading Books", subcategory: 'Kids', item: 'Children Book' },
      { title: 'Notebooks, Pens & Office Stationery', subcategory: 'Stationery', item: 'Stationery' },
    ],
  },
  {
    id: 31,
    slug: 'furniture',
    nameEn: 'Furniture & Home',
    nameAm: 'የቤት እና የቢሮ እቃዎች',
    icon: '🛋️',
    activeBorder: 'border-purple-500 bg-purple-50/70 dark:bg-purple-950/40 ring-2 ring-purple-500/20 text-purple-700 dark:text-purple-300',
    keywords: ['furniture', 'sofa', 'bed', 'mattress', 'table', 'chair', 'desk', 'wardrobe', 'closet', 'carpet', 'decor', 'የቤት እቃዎች', 'ሶፋ', 'አልጋ', 'ጠረጴዛ', 'ወንበር', 'mi\'a'],
    registeredItems: [
      { title: 'Living Room L-Shape Fabric Sofa Set', subcategory: 'Living Room', item: 'L-Shape Sofa' },
      { title: 'Queen Size Wooden Bed with Spring Mattress', subcategory: 'Bedroom', item: 'Queen Bed' },
      { title: 'Modern Dining Table with 6 Chairs', subcategory: 'Dining Room', item: 'Dining Table' },
      { title: 'Ergonomic Office Mesh Chair & Study Desk', subcategory: 'Office', item: 'Office Chair' },
      { title: '3-Door Wooden Wardrobe / Clothes Closet', subcategory: 'Bedroom', item: 'Wardrobe' },
      { title: 'Modern TV Stand & Living Room Wall Unit', subcategory: 'Living Room', item: 'TV Stand' },
      { title: 'Kitchen Storage Cabinet & Pantry Shelves', subcategory: 'Kitchen', item: 'Kitchen Cabinet' },
      { title: 'Living Room Turkish Carpet & Floor Rug', subcategory: 'Decor', item: 'Carpet' },
    ],
  },
  {
    id: 60,
    slug: 'property',
    nameEn: 'Property & Real Estate',
    nameAm: 'ቤት እና ሪል እስቴት',
    icon: '🏠',
    activeBorder: 'border-teal-500 bg-teal-50/70 dark:bg-teal-950/40 ring-2 ring-teal-500/20 text-teal-700 dark:text-teal-300',
    keywords: ['property', 'real estate', 'apartment', 'condo', 'condominium', 'villa', 'house', 'rent', 'rental', 'office', 'shop', 'land', 'plot', 'ቤት', 'ኪራይ', 'ኮንዶሚኒየም', 'ቪላ', 'መሬት', 'mana'],
    registeredItems: [
      { title: '2-Bedroom Modern Condominium Apartment', subcategory: 'Apartments', item: 'Condo 2-Bedroom' },
      { title: '3-Bedroom Luxury Apartment', subcategory: 'Apartments', item: 'Luxury Apartment' },
      { title: 'Furnished Studio Guest Apartment', subcategory: 'Apartments', item: 'Furnished Studio' },
      { title: 'G+1 / G+2 Modern Residential Villa', subcategory: 'Houses', item: 'Residential Villa' },
      { title: 'Commercial Retail Shop / Ground Floor Store', subcategory: 'Commercial', item: 'Commercial Shop' },
      { title: 'Prime Office Space / Business Hall', subcategory: 'Commercial', item: 'Office Space' },
      { title: 'Residential Land Plot with Title Deed', subcategory: 'Land', item: 'Land Plot' },
    ],
  },
  {
    id: 71,
    slug: 'tools',
    nameEn: 'Tools & Machinery',
    nameAm: 'መሳሪያዎች እና ማሽነሪዎች',
    icon: '🔧',
    activeBorder: 'border-slate-500 bg-slate-50/70 dark:bg-slate-900/60 ring-2 ring-slate-500/20 text-slate-700 dark:text-slate-300',
    keywords: ['tool', 'tools', 'machinery', 'drill', 'grinder', 'welding', 'generator', 'pump', 'equipment', 'safety', 'መሳሪያ', 'ማሽን', 'meeshaalee'],
    registeredItems: [
      { title: 'Electric Power Drill & Impact Driver Set', subcategory: 'Power Tools', item: 'Electric Drill' },
      { title: 'Angle Grinder with Cutting & Grinding Discs', subcategory: 'Power Tools', item: 'Angle Grinder' },
      { title: 'Professional Mechanics Hand Tool Kit (150 Pcs)', subcategory: 'Hand Tools', item: 'Tool Kit' },
      { title: 'Inverter ARC Welding Machine 250A', subcategory: 'Machinery', item: 'Welding Machine' },
      { title: 'Heavy Duty Power Generator (5kVA - 10kVA)', subcategory: 'Power', item: 'Power Generator' },
      { title: 'Submersible Water Pump & High Pressure Motor', subcategory: 'Pumps', item: 'Water Pump' },
      { title: 'Industrial Safety Helmet, Boots & Harness', subcategory: 'Safety', item: 'Safety Gear' },
    ],
  },
];

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
  const [selectedPresetSlug, setSelectedPresetSlug] = useState<string | null>(null);
  const [showAdvancedTree, setShowAdvancedTree] = useState(false);
  const [categorySearch, setCategorySearch] = useState('');

  useEffect(() => {
    if (form.category_hierarchy_id) {
      const match = CANONICAL_PRESETS.find(p => p.id === form.category_hierarchy_id);
      if (match) setSelectedPresetSlug(match.slug);
    }
  }, [form.category_hierarchy_id]);

  const activePreset = CANONICAL_PRESETS.find(p => p.slug === selectedPresetSlug || p.id === form.category_hierarchy_id);

  const handleSelectPreset = (preset: CanonicalCategoryPreset) => {
    setSelectedPresetSlug(preset.slug);
    setForm(f => ({
      ...f,
      category_hierarchy_id: preset.id,
      category_id: String(preset.id),
      subcategory: preset.registeredItems[0]?.subcategory ?? f.subcategory,
    }));
    setAttributeValues({});
  };

  const handleSelectItem = (item: { title: string; subcategory: string; item: string }) => {
    setForm(f => ({
      ...f,
      title: item.title,
      subcategory: item.subcategory,
      item: item.item,
    }));
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

      {/* ── STEP 1: CHOOSE PRODUCT CATEGORY & REGISTERED NAME ── */}
      <div className="space-y-3.5 p-4 bg-gradient-to-br from-blue-50/60 via-white to-slate-50 dark:from-zinc-900 dark:via-zinc-900 dark:to-zinc-950 border border-blue-200/80 dark:border-zinc-800 rounded-3xl shadow-sm">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2.5">
            <span className="w-6 h-6 rounded-full bg-[#1769E0] text-white text-xs font-black flex items-center justify-center shrink-0">1</span>
            <div>
              <h3 className="text-sm font-black text-zinc-900 dark:text-white flex items-center gap-2">
                <span>{locale === 'am' ? 'ደረጃ 1፦ የምርት አይነት እና ምድብ ይምረጡ' : 'Step 1: Choose Product Category & Type'}</span>
                <span className="text-[10px] font-bold text-red-500 uppercase tracking-wider">{locale === 'am' ? '(ግዴታ)' : '(Required)'}</span>
              </h3>
              <p className="text-[11px] text-zinc-500 dark:text-zinc-400">
                {locale === 'am'
                  ? 'በፍለጋ እና በማጣሪያ በቀላሉ እንዲገኝ ይፈልጉ ወይም ከተመዘገቡት ፈጣን ምድቦች አንዱን ይምረጡ'
                  : 'Search what you are listing or tap a popular category for instant indexing & filters'}
              </p>
            </div>
          </div>
          {form.category_hierarchy_id && (
            <span className="inline-flex items-center gap-1 text-[11px] font-bold text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/40 px-2.5 py-1 rounded-full border border-emerald-200 dark:border-emerald-800 shrink-0">
              <CheckCircle2 size={13} /> {locale === 'am' ? 'ምድብ ተመርጧል' : 'Category Selected'}
            </span>
          )}
        </div>

        {/* 🔍 Unified Search What You Are Listing Input */}
        <div className="relative">
          <div className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-[#1769E0]">
            <Search size={16} />
          </div>
          <input
            type="text"
            value={categorySearch}
            onChange={e => {
              const val = e.target.value;
              setCategorySearch(val);
              if (val.trim()) {
                const lower = val.toLowerCase().trim();
                const match = CANONICAL_PRESETS.find(p =>
                  p.slug.includes(lower) ||
                  p.nameEn.toLowerCase().includes(lower) ||
                  p.nameAm.includes(lower) ||
                  p.keywords.some(k => k.toLowerCase().includes(lower) || lower.includes(k.toLowerCase()))
                );
                if (match && match.id !== form.category_hierarchy_id) {
                  handleSelectPreset(match);
                }
              }
            }}
            placeholder={
              locale === 'am'
                ? 'ምን ሊሸጡ ወይም ሊያከራዩ ነው? (ለምሳሌ፦ ስልክ፣ ላፕቶፕ፣ መኪና፣ ቀሚስ፣ መጽሐፍ፣ ሶፋ፣ ቤት...)'
                : 'Search what you are listing (e.g. Phone, Laptop, Car, Dress, Book, Sofa, Apartment, Shoes...)'
            }
            className="w-full pl-10 pr-9 py-2.5 bg-white dark:bg-zinc-800 border-2 border-blue-300 dark:border-blue-900/60 focus:border-[#1769E0] focus:ring-2 focus:ring-[#1769E0]/20 rounded-2xl text-xs sm:text-sm font-medium text-zinc-900 dark:text-zinc-100 outline-none transition-all placeholder:text-zinc-400 shadow-xs"
          />
          {categorySearch && (
            <button
              type="button"
              onClick={() => setCategorySearch('')}
              className="absolute inset-y-0 right-0 pr-3 flex items-center text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200"
            >
              <X size={14} />
            </button>
          )}
        </div>

        {/* Category Preset 1-Click Buttons */}
        <div className="space-y-1.5">
          <p className="text-[11px] font-bold text-zinc-600 dark:text-zinc-400">
            {locale === 'am' ? 'ወይም ፈጣን ዋና ምድብ ይምረጡ፦' : 'Or choose popular marketplace category:'}
          </p>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-2">
            {CANONICAL_PRESETS.map((preset) => {
              const isSelected = form.category_hierarchy_id === preset.id || selectedPresetSlug === preset.slug;
              return (
                <button
                  key={preset.id}
                  type="button"
                  onClick={() => {
                    handleSelectPreset(preset);
                    if (categorySearch) setCategorySearch('');
                  }}
                  className={`p-2.5 rounded-2xl border text-left transition-all tap-bounce cursor-pointer flex items-center gap-2.5 min-h-[58px] ${
                    isSelected
                      ? preset.activeBorder
                      : 'border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-800/60 hover:border-zinc-300 dark:hover:border-zinc-700'
                  }`}
                >
                  <span className="text-xl sm:text-2xl shrink-0">{preset.icon}</span>
                  <div className="min-w-0 flex-1">
                    <p className="text-xs font-black text-zinc-900 dark:text-white leading-tight truncate">
                      {locale === 'am' ? preset.nameAm : preset.nameEn}
                    </p>
                    <p className="text-[10px] text-zinc-400 dark:text-zinc-500 font-medium truncate">
                      {locale === 'am' ? preset.nameEn : preset.nameAm}
                    </p>
                  </div>
                  {isSelected && <span className="w-2 h-2 rounded-full bg-[#1769E0] shrink-0" />}
                </button>
              );
            })}
          </div>
        </div>

        {/* Registered Catalog Items Chips (1-Click to Insert Name) */}
        {activePreset && (
          <div className="mt-2 pt-3 border-t border-zinc-200/70 dark:border-zinc-800 space-y-2">
            <div className="flex items-center justify-between">
              <p className="text-xs font-bold text-zinc-700 dark:text-zinc-300 flex items-center gap-1.5">
                <Tag size={13} className="text-[#1769E0]" />
                <span>
                  {locale === 'am'
                    ? `የተመዘገቡ የ${activePreset.nameAm} ስሞች (ጠቅ በማድረግ ርዕስ ይሙሉት):`
                    : `Registered ${activePreset.nameEn} Catalog Names (Click to insert):`}
                </span>
              </p>
              <span className="text-[10px] font-semibold text-zinc-400">
                {activePreset.registeredItems.length} {locale === 'am' ? 'እቃዎች' : 'items'}
              </span>
            </div>
            <div className="flex flex-wrap gap-1.5 max-h-36 overflow-y-auto p-1 bg-zinc-50/50 dark:bg-zinc-900/50 rounded-2xl border border-zinc-200/60 dark:border-zinc-800/60">
              {activePreset.registeredItems
                .filter(item => !categorySearch || item.title.toLowerCase().includes(categorySearch.toLowerCase()) || item.subcategory.toLowerCase().includes(categorySearch.toLowerCase()))
                .map((item, idx) => {
                  const isItemActive = form.title.startsWith(item.title) || form.item === item.item;
                  return (
                    <button
                      key={idx}
                      type="button"
                      onClick={() => handleSelectItem(item)}
                      className={`px-2.5 py-1.5 rounded-xl text-xs font-semibold transition-all tap-bounce cursor-pointer flex items-center gap-1.5 border ${
                        isItemActive
                          ? 'bg-[#1769E0] text-white border-[#1769E0] shadow-sm shadow-[#1769E0]/20'
                          : 'bg-white dark:bg-zinc-800 text-zinc-700 dark:text-zinc-300 border-zinc-200 dark:border-zinc-700 hover:bg-blue-50 dark:hover:bg-blue-950/40 hover:border-blue-300'
                      }`}
                    >
                      {isItemActive && <Check size={11} />}
                      <span>{item.title}</span>
                    </button>
                  );
                })}
            </div>
          </div>
        )}

        {/* Toggle for Advanced / All 80+ Categories */}
        <div className="pt-1">
          <button
            type="button"
            onClick={() => setShowAdvancedTree(prev => !prev)}
            className="text-xs font-semibold text-[#1769E0] dark:text-[#38A5FF] hover:underline flex items-center gap-1 cursor-pointer"
          >
            <span>
              {showAdvancedTree
                ? (locale === 'am' ? '▲ ሙሉውን የምድብ ዛፍ ደብቅ' : '▲ Hide complete category tree')
                : (locale === 'am' ? '▼ ሌላ እቃ ነው? ሙሉውን የምድብ ዝርዝር (80+ ምድቦች) ይፈልጉ' : '▼ Not in quick list? Search all 80+ category branches')}
            </span>
          </button>
          {showAdvancedTree && (
            <div className="mt-2.5 p-3 bg-zinc-50 dark:bg-zinc-800/50 rounded-2xl border border-zinc-200 dark:border-zinc-700">
              <label className={lbl}>{t.ui.selectCategory}</label>
              <CategoryTreeSelector
                value={form.category_hierarchy_id ?? (form.category_id ? Number(form.category_id) : undefined)}
                language={locale}
                onChange={(id, node: CategoryNode) => {
                  setForm(f => ({
                    ...f,
                    category_hierarchy_id: id,
                    category_id: id ? String(id) : '',
                    subcategory: node?.name_en || f.subcategory,
                  }));
                  setAttributeValues({});
                }}
                placeholder={t.ui.selectCategoryOption || 'Select a category branch...'}
              />
            </div>
          )}
        </div>
      </div>

      {/* ── STEP 2: PRODUCT TITLE & PRICING ── */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <label className={lbl}>
            <span className="w-4 h-4 rounded-full bg-zinc-700 text-white text-[10px] font-black inline-flex items-center justify-center mr-1.5">2</span>
            {t.ui.title} <span className="text-red-500">*</span>
          </label>
          {form.category_hierarchy_id ? (
            <span className="text-[11px] font-bold text-[#1769E0] dark:text-[#38A5FF] flex items-center gap-1">
              <CheckCircle2 size={12} /> Linked to {activePreset ? (locale === 'am' ? activePreset.nameAm : activePreset.nameEn) : 'Category'}
            </span>
          ) : (
            <span className="text-[11px] font-medium text-amber-600 dark:text-amber-400">
              ⚠️ Select category above first
            </span>
          )}
        </div>
        <Input
          required
          placeholder={isAnnouncement ? 'Announcement / Vacancy / Service Title' : isRental ? 'Property or Item title (e.g. 2 Bedroom Condominium)' : (t.ui.product_title || 'Product title (e.g. Apple iPhone 15 Pro Max 256GB, Toyota Vitz)')}
          value={form.title}
          onChange={e => setForm(f => ({ ...f, title: e.target.value }))}
        />
        
        {/* Real-time Auto-Slug & Discovery Indicator */}
        <div className="flex items-center justify-between text-[11px] text-zinc-500 dark:text-zinc-400 px-1">
          <span className="flex items-center gap-1 truncate font-mono text-[10px] text-zinc-600 dark:text-zinc-400">
            <LinkIcon size={11} className="text-[#1769E0] shrink-0" />
            <span className="text-zinc-400">huntiq.et/products/</span>
            <span className="font-bold text-zinc-800 dark:text-zinc-200 truncate">
              {form.title
                ? form.title.toLowerCase().trim().replace(/[^\w\s-]/g, '').replace(/[\s_-]+/g, '-').replace(/^-+|-+$/g, '') || 'listing'
                : 'your-product-slug'}
            </span>
          </span>
          <span className="text-[10px] text-emerald-600 dark:text-emerald-400 font-semibold shrink-0 ml-2">
            ✓ Auto-Indexed
          </span>
        </div>
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
