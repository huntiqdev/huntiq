export interface Category {
  id: number;
  name: string;
  slug: string;
  name_en?: string;
  name_am?: string;
  name_or?: string;
  name_ti?: string;
  sort_order?: number;
  display_order?: number;
  is_active?: boolean;
  parent_id?: number;
  icon?: string;
  image_url?: string;
  children?: Category[];
}

export interface CategoryHierarchyItem {
  id: number;
  name: string;
  name_en?: string;
  name_am?: string;
  name_om?: string;
  name_or?: string;
  name_ti?: string;
  slug: string;
  parent_id?: number | null;
  level: number;
  path?: string;
  display_order: number;
  icon?: string;
  image_url?: string;
  is_active: boolean;
  is_leaf?: boolean;
  product_count?: number;
  breadcrumbs?: Array<{ id: number; name: string; slug: string; level: number }>;
  children?: CategoryHierarchyItem[];
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

  // Gold & Premium Merchant Features
  custom_url?: string;
  social_links?: Record<string, string>;

  // Brokerage & Publisher Transparency
  publisher_role?: 'owner' | 'authorized_seller' | 'broker_agent' | 'representative' | 'government_officer';
  organization_id?: number;
  brokerage_disclosure_agreed?: boolean;

  // Engagement & Trending Metrics
  view_count?: number;
  contact_count?: number;
  
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
  sort_by?: 'trending' | 'popular' | 'views' | 'most_viewed' | 'price_asc' | 'price_desc' | 'newest' | 'nearest';
  lat?: number;
  lng?: number;
  max_distance?: number;
  mode?: 'ai' | 'manual';
  locale?: string;
  preferred_language?: string;
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

export interface I18nLocation {
  en: string;
  [locale: string]: string | undefined;
}

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

export interface WalletAccount {
  id: number;
  user_id: number;
  currency: string;
  available_balance: number;
  pending_balance: number;
  reserved_balance: number;
  total_withdrawn: number;
  lifetime_earned: number;
  is_locked: boolean;
  created_at: string;
  updated_at: string;
}

export interface LedgerTransaction {
  id: number;
  ledger_tx_id: string;
  user_id: number;
  debit_account: string;
  credit_account: string;
  amount: number;
  currency: string;
  entry_type: string;
  reference_id?: string;
  source_event: string;
  status: string;
  idempotency_key: string;
  notes?: string;
  created_at: string;
}

export interface TelebirrOrder {
  id: number;
  out_trade_no: string;
  user_id: number;
  amount: number;
  currency: string;
  title: string;
  order_type: string;
  status: string;
  prepay_id?: string;
  transaction_no?: string;
  payer_phone?: string;
  created_at: string;
  updated_at: string;
}

export interface ReferralRewardPolicy {
  id: number;
  policy_name: string;
  registration_reward_amount: number;
  transaction_reward_percent: number;
  min_qualifying_tx_amount: number;
  max_reward_per_tx: number;
  reward_holding_days: number;
  version: number;
  is_active: boolean;
  effective_from: string;
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

// ── National Identity Verification (Fayda eKYC) ──────────────────────────────

export type IdentityVerificationStatus = 'unverified' | 'pending' | 'verified' | 'rejected' | 'escalated';
export type PhoneMatchStatus = 'verified' | 'mismatch' | 'unsupported_by_provider' | 'exempted';

export interface IdentityVerification {
  id: number;
  user_id: number;
  provider_code: string;
  credential_type: string;
  id_masked: string;
  full_legal_name: string;
  phone_matched: boolean;
  phone_match_status: PhoneMatchStatus;
  auth_reference: string;
  status: IdentityVerificationStatus;
  verified_at?: string;
  created_at: string;
}

// ── Account Entity Classification & Organizations ─────────────────────────────

export type AccountEntityType = 'INDIVIDUAL' | 'COMPANY' | 'GOVERNMENT' | 'OTHER_LEGALLY_APPROVED_ENTITY';
export type OrgVerificationStatus = 'pending' | 'verified' | 'rejected' | 'suspended';
export type OrgMemberRole = 'owner' | 'admin' | 'publisher' | 'finance' | 'viewer';

export interface Organization {
  id: number;
  name: string;
  entity_type: AccountEntityType;
  registration_number?: string;
  tin_number?: string;
  official_address?: string;
  official_email?: string;
  official_phone?: string;
  verification_status: OrgVerificationStatus;
  publishing_mandate_verified: boolean;
  created_at: string;
  updated_at: string;
}

export interface OrganizationMember {
  id: number;
  organization_id: number;
  user_id: number;
  role: OrgMemberRole;
  joined_at: string;
}
