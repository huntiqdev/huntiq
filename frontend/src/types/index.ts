export interface Category {
  id: number;
  name: string;
  slug: string;
  name_en?: string;
  name_am?: string;
  name_or?: string;
  name_ti?: string;
  sort_order?: number;
  is_active?: boolean;
  parent_id?: number;
  children?: Category[];
}

export interface ProductImage {
  id: number;
  product_id: number;
  image_url: string;
  order_number: number;
}

export interface Location {
  id: number;
  region: string;
  zone?: string;
  subcity?: string;
  woreda?: string;
}

export interface Product {
  id: number;
  seller_id: number;
  title: string;
  price: number;
  description: string;
  condition: string;
  status: string;
  admin_notes?: string;
  subcategory?: string;
  item?: string;
  kebele?: string;
  specific_area?: string;
  latitude?: number;
  longitude?: number;
  location_hierarchy_id?: number;
  expires_at?: string;
  listing_duration_days?: number;
  total_fee_percent?: number;
  total_fee_amount?: number;
  fee_paid?: boolean;
  seller_tier: string; // normal | gold | premium — used for image blur logic
  required_buyer_tier: string; // normal | gold | premium — minimum buyer tier to order
  specifications_locked?: boolean;
  
  // Rental fields (only populated when listing_type='rent')
  listing_type: 'sale' | 'rent';
  rental_price?: number;
  rental_duration_value?: number;
  rental_duration_unit?: 'hour' | 'day' | 'week' | 'month' | 'year';
  deposit_amount?: number;
  available_from?: string; // DATE
  available_until?: string; // DATE
  is_currently_rented?: boolean;
  current_tenant_id?: number;
  rented_until?: string; // DATE
  
  category?: Category;
  location?: Location;
  images?: ProductImage[];
  created_at: string;
  updated_at?: string;
}

export interface ProductFilters {
  category_id?: number;
  category_hierarchy_id?: number;
  include_descendants?: boolean;
  subcategory?: string;
  item?: string;
  condition?: string;
  min_price?: number;
  max_price?: number;
  location_id?: number;
  location_hierarchy_id?: number;
  region?: string;
  zone?: string;
  subcity?: string;
  woreda?: string;
  search?: string;
  listing_type?: 'sale' | 'rent'; // Filter by listing type
  page?: number;
  limit?: number;
}
// LocationFilters is merged into ProductFilters above — kept as alias for compatibility
export type LocationFilters = Pick<ProductFilters, 'region' | 'zone' | 'subcity' | 'woreda'>;
export interface Transaction {
  id: number;
  product_id: number;
  buyer_name: string;
  buyer_phone: string;
  amount: number;
  broker_fee: number;
  seller_gets: number;
  status: string;
  notes?: string;
  created_at: string;
  updated_at?: string;
}

export interface OrderRequest {
  id: number;
  product_id: number;
  seller_id: number;
  buyer_name: string;
  buyer_phone: string;
  buyer_email?: string;
  message?: string;
  offer_price?: number;
  status: 'pending' | 'accepted' | 'declined' | 'completed' | 'cancelled';
  buyer_tier_at_order?: string;
  created_at: string;
  updated_at?: string;
}

export type { I18nLocation } from './location';

// ── Referral & Commission Types ──────────────────────────────────────────────

export interface ReferralBalance {
  total_earned: number;
  total_withdrawn: number;
  available_balance: number;
  pending_amount: number; // Commission not yet withdrawable (30-day pending period)
}

export interface CommissionEntry {
  id: number;
  invited_user_id: number;
  invited_user_name: string;
  transaction_amount: number;
  commission_amount: number;
  transaction_type: string;
  status: 'pending' | 'confirmed' | 'reversed'; // Chargeback protection status
  confirmed_at?: string;
  pending_until?: string;
  created_at: string;
}

export interface ReferralStats {
  referral_code: string;
  referral_link: string; // Frontend-generated full URL
  referral_suspended: boolean;
  invitation_accepted_at?: string;
  inviter_code?: string;
  inviter_name?: string;
  invited_count: number;
  total_commission: number;
}

export interface WithdrawalRequest {
  id: number;
  user_id: number;
  user_name?: string; // Only in admin view
  user_phone?: string; // Only in admin view
  amount: number;
  status: 'pending' | 'approved' | 'rejected' | 'paid';
  payment_method: string;
  payer_phone: string;
  reference?: string;
  admin_note?: string;
  requested_at: string;
  reviewed_at?: string;
  reviewed_by?: number;
  paid_at?: string;
}

// ── Validation error shapes ───────────────────────────────────────────────────
// All form validation interfaces live here so every part of the system
// uses the same shape. Values are i18n keys (from t.errors / t.ui),
// never hardcoded English strings.

export interface LoginErrors {
  phone?: string;
  password?: string;
}

export interface RegisterErrors {
  name?: string;
  phone?: string;
  email?: string;
  password?: string;
  seller_preference?: string;
}

export interface ProductErrors {
  title?: string;
  price?: string;
  images?: string;
  condition?: string;
}

export interface OrderErrors {
  name?: string;
  phone?: string;
}
