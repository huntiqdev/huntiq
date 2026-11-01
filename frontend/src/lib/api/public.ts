import { fetchAPI, fetchAPIWithMeta, authHeaders, API_URL, ApiError } from './core';

export interface RealTimePriceStats {
  min_price: number;
  avg_price: number;
  max_price: number;
  currency: string;
}

export interface RealTimeTopPick {
  id: number;
  title: string;
  price: number;
  seller_tier: string;
  badge: string;
  reason: string;
}

export interface RealTimeProductAnalysis {
  total_analyzed: number;
  total_matching?: number;
  premium_count: number;
  gold_count: number;
  normal_count: number;
  price_stats: RealTimePriceStats;
  top_premium_pick?: RealTimeTopPick;
  top_gold_pick?: RealTimeTopPick;
  market_verdict: string;
  deal_ratings?: Array<{
    product_id: number;
    deal_score: number;
    badge: string;
    value_verdict: string;
  }>;
  execution_time_ms?: number;
}

export interface ReasoningStep {
  step: number;
  title: string;
  detail: string;
}

export interface MarketplaceReasoningChain {
  intent_breakdown: string;
  market_feasibility: string;
  trade_off_analysis: string;
  inspection_protocol: string;
  steps: ReasoningStep[];
}

export interface SearchAIMeta {
  ai_mode: boolean;
  normalized_query?: string;
  cleaned_search_term?: string;
  sort_by?: string;
  primary_intent?: string;
  localized_answer?: string;
  detected_locale?: string;
  deep_parameters?: Record<string, any>;
  pii_detected?: boolean;
  is_safe?: boolean;
  safety_flags?: string[];
  real_time_analysis?: RealTimeProductAnalysis;
  market_verdict?: string;
  reasoning_chain?: MarketplaceReasoningChain;
  total_matching?: number;
  returned_count?: number;
  page?: number;
  limit?: number;
}

export interface PublicPaymentSettings {
  telebirr_merchant_phone: string;
  safaricom_merchant_phone: string;
  cbe_account_number: string;
  gold_upgrade_price: string;
  premium_upgrade_price: string;
  buyer_gold_upgrade_price: string;
  buyer_premium_upgrade_price: string;
  min_referral_withdrawal: string;
  max_referral_withdrawal: string;
  referral_commission_rate: string;
  broker_fee_rate: string;
  normal_startup_fee: string;
  gold_startup_fee: string;
  premium_startup_fee: string;
  normal_weekly_increment: string;
  gold_weekly_increment: string;
  premium_weekly_increment: string;
  company_phone: string;
  company_email: string;
  company_location: string;
  gold_threshold: string;
  premium_threshold: string;
  otp_gateway: string;
}

export interface AutoCategorizeResult {
  detected_title: string;
  primary_intent: 'sale' | 'rent';
  category_hierarchy_id: number;
  category_slug: string;
  root_category_name: string;
  subcategory_name: string;
  leaf_destination_name: string;
  full_branch_path: string;
  branch_path_localized: string;
  legacy_category_id: number;
  condition: 'new' | 'used';
  suggested_price_etb: number;
  price_range_min: number;
  price_range_max: number;
  suggested_description: string;
  detected_locale: string;
  brand?: string;
  model?: string;
  confidence: number;
  human_error_corrected: boolean;
  original_term?: string;
  corrected_term?: string;
}

export const publicApi = {
  getProductsWithMeta: (params?: Record<string, any>) => {
    // Clean undefined values and map filter keys to backend format
    const cleanParams: Record<string, string> = {};
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined && v !== null && v !== '') {
          cleanParams[k] = String(v);
        }
      });
    }
    const query = Object.keys(cleanParams).length > 0 ? '?' + new URLSearchParams(cleanParams).toString() : '';
    return fetchAPIWithMeta<any[], SearchAIMeta>(`/products${query}`).then(r => ({
      products: r.data ?? [],
      meta: r.meta,
    }));
  },
  getProducts: (params?: Record<string, any>) => {
    return publicApi.getProductsWithMeta(params).then(r => r.products);
  },
  getProduct: (id: number) => fetchAPI<any>(`/products/${id}`),
  getProductAttributes: (id: number) => fetchAPI<any[]>(`/products/${id}/attributes`).then(r => r ?? []),
  getRootCategories: (lang = 'en') => fetchAPI<any[]>(`/categories/roots?lang=${encodeURIComponent(lang)}`).then(r => r ?? []),
  getCategoryTree: (lang = 'en') => fetchAPI<any[]>(`/categories/tree?lang=${encodeURIComponent(lang)}`).then(r => r ?? []),
  getCategoryChildren: (id: number, lang = 'en') => fetchAPI<any[]>(`/categories/${id}/children?lang=${encodeURIComponent(lang)}`).then(r => r ?? []),
  getCategoryBySlug: (slug: string, lang = 'en') => fetchAPI<any>(`/categories/slug/${encodeURIComponent(slug)}?lang=${encodeURIComponent(lang)}`),
  getCategoryById: (id: number, lang = 'en') => fetchAPI<any>(`/categories/${id}?lang=${encodeURIComponent(lang)}`),
  checkCategorySlug: (slug: string, excludeId?: number) =>
    fetchAPI<{ available: boolean; slug: string }>(`/categories/check-slug?slug=${encodeURIComponent(slug)}${excludeId ? `&exclude_id=${excludeId}` : ''}`),
  getCategoryBreadcrumbs: (id: number, lang = 'en') => fetchAPI<any[]>(`/categories/${id}/ancestors?lang=${encodeURIComponent(lang)}`).then(r => r ?? []),
  getCategoryAttributes: (id: number) => fetchAPI<any[]>(`/categories/${id}/attributes`).then(r => r ?? []),
  getCategories: () => fetchAPI<any[]>('/categories').then(r => r ?? []),
  getPopularCategories: (days = 7, limit = 8) =>
    fetchAPI<any[]>(`/categories/popular?days=${days}&limit=${limit}`).then(r => r ?? []),
  getRegions: () => fetchAPI<Array<{en:string;am?:string;or?:string;ti?:string}>>('/locations/regions').then(r => r ?? []),
  getZones: (region: string) => fetchAPI<Array<{en:string;am?:string;or?:string;ti?:string}>>(`/locations/zones?region=${encodeURIComponent(region)}`).then(r => r ?? []),
  getSubcities: (region: string, zone?: string) => {
    const q = zone ? `region=${encodeURIComponent(region)}&zone=${encodeURIComponent(zone)}` : `region=${encodeURIComponent(region)}`;
    return fetchAPI<Array<{en:string;am?:string;or?:string;ti?:string}>>(`/locations/subcities?${q}`).then(r => r ?? []);
  },
  getWoredas: (region: string, subcity: string, zone?: string) => {
    const parts: string[] = [`region=${encodeURIComponent(region)}`];
    if (zone) parts.push(`zone=${encodeURIComponent(zone)}`);
    if (subcity) parts.push(`subcity=${encodeURIComponent(subcity)}`);
    return fetchAPI<Array<{en:string;am?:string;or?:string;ti?:string}>>(`/locations/woredas?${parts.join('&')}`).then(r => r ?? []);
  },
  searchLocations: (q: string) => fetchAPI<any[]>(`/locations/search?q=${encodeURIComponent(q)}`).then(r => r ?? []),
  createOrderRequest: (data: any, token?: string) =>
    fetchAPI<any>('/orders', { method: 'POST', body: JSON.stringify(data), headers: token ? authHeaders(token) : undefined }),
  getPublicPaymentSettings: () =>
    fetchAPI<PublicPaymentSettings>('/public/payment-settings'),
  /** 
   * Fetch ALL public platform settings (fees, limits, thresholds, etc.)
   * Used by SettingsProvider for frontend configuration
   */
  getPublicSettings: () =>
    fetchAPI<Record<string, string>>('/public/settings'),
  submitContactMessage: (data: { name: string; email: string; message: string }) =>
    fetchAPI<{ success: boolean; message: string }>('/contact', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  getBanners: (type: 'buyer' | 'seller') =>
    fetchAPI<any[]>(`/banners?type=${type}`).then(r => r ?? []),
  reverseGeocode: (latitude: number, longitude: number, language = 'en') =>
    fetchAPI<any>('/locations/reverse-geocode', {
      method: 'POST',
      body: JSON.stringify({ latitude, longitude, language }),
    }),
  getLocationTree: (lang = 'en') =>
    fetchAPI<any[]>(`/locations/tree?lang=${encodeURIComponent(lang)}`).then(r => r ?? []),
  getLocationChildren: (parentId?: number, lang = 'en') => {
    const q = parentId !== undefined ? `?parent_id=${parentId}&lang=${encodeURIComponent(lang)}` : `?lang=${encodeURIComponent(lang)}`;
    return fetchAPI<any[]>(`/locations/children${q}`).then(r => r ?? []);
  },
  searchLocationHierarchy: (query: string, lang = 'en') =>
    fetchAPI<any[]>(`/locations/search-hierarchy?q=${encodeURIComponent(query)}&lang=${encodeURIComponent(lang)}`).then(r => r ?? []),
  autoCategorize: (data: { title?: string; description?: string; locale?: string }) =>
    fetchAPI<AutoCategorizeResult>('/products/auto-categorize', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
};

