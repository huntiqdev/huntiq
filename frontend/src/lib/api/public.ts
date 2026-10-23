import { fetchAPI, authHeaders, API_URL, ApiError } from './core';

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

export const publicApi = {
  getProducts: (params?: Record<string, any>) => {
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
    return fetchAPI<any[]>(`/products${query}`).then(r => r ?? []);
  },
  getProduct: (id: number) => fetchAPI<any>(`/products/${id}`),
  getProductAttributes: (id: number) => fetchAPI<any[]>(`/products/${id}/attributes`).then(r => r ?? []),
  getCategoryTree: (lang = 'en') => fetchAPI<any[]>(`/categories/tree?lang=${lang}`).then(r => r ?? []),
  getCategoryBreadcrumbs: (id: number, lang = 'en') => fetchAPI<any[]>(`/categories/${id}/ancestors?lang=${lang}`).then(r => r ?? []),
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
};

