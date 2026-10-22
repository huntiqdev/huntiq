/**
 * errors.ts — Central error translation and classification system.
 *
 * Supports two resolution strategies (in priority order):
 *   1. Structured error codes from backend (e.g. "AUTH_INVALID_CREDENTIALS")
 *   2. Substring matching on the raw message string (legacy fallback)
 *
 * Usage:
 *   import { translateError } from '@/lib/errors';
 *   const { message, type } = translateError(err.message, t.errors, err.code);
 *
 * The `t.errors` object comes from useLocaleContext() → t.errors.
 */

export type ErrorType =
  | 'auth'
  | 'auth_account_inactive'
  | 'auth_account_not_found'
  | 'auth_buyer_registration_closed'
  | 'auth_email_already_in_use'
  | 'auth_invalid_credentials'
  | 'auth_phone_already_registered'
  | 'auth_phone_not_verified'
  | 'auth_registration_closed'
  | 'auth_session_expired'
  | 'auth_token_invalid'
  | 'auth_wrong_password'
  | 'network'
  | 'network_error'
  | 'network_no_connection'
  | 'network_offline'
  | 'network_packet_loss'
  | 'network_server_down'
  | 'network_slow'
  | 'network_timeout'
  | 'network_unreachable'
  | 'not_found'
  | 'order'
  | 'order_not_active'
  | 'order_not_available'
  | 'order_not_found'
  | 'order_not_owned'
  | 'otp'
  | 'otp_expired'
  | 'otp_invalid'
  | 'otp_locked'
  | 'otp_no_active'
  | 'otp_rate_limit'
  | 'otp_required'
  | 'otp_send_failed'
  | 'payment'
  | 'payment_already_paid'
  | 'payment_insufficient_balance'
  | 'payment_invalid_method'
  | 'payment_max_withdrawal'
  | 'payment_min_withdrawal'
  | 'payment_pending_exists'
  | 'payment_pending_withdrawal'
  | 'payment_reference_required'
  | 'payment_rejection_reason_required'
  | 'payment_sold_product'
  | 'payment_upgrade_closed'
  | 'permission'
  | 'permission_admin_required'
  | 'permission_cannot_change_own_role'
  | 'permission_cannot_delete_self'
  | 'permission_denied'
  | 'permission_invitation_already_accepted'
  | 'permission_last_super_admin'
  | 'permission_own_referral_code'
  | 'permission_referral_blocked'
  | 'permission_referral_suspended'
  | 'permission_seller_tier_required'
  | 'permission_super_admin_delete'
  | 'permission_super_admin_required'
  | 'product'
  | 'product_not_active'
  | 'product_not_available'
  | 'product_not_found'
  | 'product_not_owned'
  | 'rate_limit'
  | 'rate_limit_deposit'
  | 'rate_limit_general'
  | 'rate_limit_message'
  | 'rate_limit_notification'
  | 'rate_limit_order'
  | 'rate_limit_otp'
  | 'rate_limit_payment'
  | 'rate_limit_product'
  | 'rate_limit_referral'
  | 'rate_limit_transaction'
  | 'rate_limit_upload'
  | 'rate_limit_withdrawal'
  | 'server'
  | 'server_down'
  | 'server_error'
  | 'server_internal_error'
  | 'server_maintenance'
  | 'server_overloaded'
  | 'server_service_unavailable'
  | 'server_timeout'
  | 'server_unavailable'
  | 'unknown'
  | 'upload'
  | 'upload_cloud_failed'
  | 'upload_failed'
  | 'upload_file_too_large'
  | 'upload_limit_reached'
  | 'upload_url_not_loadable'
  | 'upload_wrong_type'
  | 'validation'
  | 'validation_amount_positive'
  | 'validation_description_too_long'
  | 'validation_location_invalid'
  | 'validation_location_required'
  | 'validation_name_required'
  | 'validation_password_min'
  | 'validation_phone_invalid'
  | 'validation_phone_required'
  | 'validation_price_required'
  | 'validation_required'
  | 'validation_setting_key_unknown'
  | 'validation_setting_numeric'
  | 'validation_tier_invalid'
  | 'validation_title_required'
  | 'validation_title_too_long';


export interface TranslatedError {
  message: string;
  type: ErrorType;
  /** The raw backend message, preserved for logging */
  raw: string;
  /** The structured error code from backend, if present */
  code?: string;
}

// ── Structured error code → (translation_key, error_type) ────────────────────
// These map the machine-readable codes sent by the backend to UI keys.
// Priority: code-based lookup runs BEFORE substring matching.

type CodeEntry = [string, ErrorType];

const CODE_MAP: Record<string, CodeEntry> = {
  // Auth
  AUTH_INVALID_CREDENTIALS: ['auth_invalid_credentials', 'auth'],
  AUTH_WRONG_PASSWORD: ['auth_wrong_password', 'auth'],
  AUTH_ACCOUNT_INACTIVE: ['auth_account_inactive', 'auth'],
  AUTH_PHONE_NOT_VERIFIED: ['auth_phone_not_verified', 'auth'],
  AUTH_SESSION_EXPIRED: ['auth_session_expired', 'auth'],
  AUTH_TOKEN_INVALID: ['auth_token_invalid', 'auth'],
  AUTH_REGISTRATION_CLOSED: ['auth_registration_closed', 'auth'],
  AUTH_BUYER_REGISTRATION_CLOSED: ['auth_buyer_registration_closed', 'auth'],
  AUTH_PHONE_ALREADY_REGISTERED: ['auth_phone_already_registered', 'auth'],
  AUTH_EMAIL_ALREADY_IN_USE: ['auth_email_already_in_use', 'auth'],
  AUTH_ACCOUNT_NOT_FOUND: ['auth_account_not_found', 'auth'],

  // OTP
  OTP_LOCKED: ['otp_locked', 'otp'],
  OTP_EXPIRED: ['otp_expired', 'otp'],
  OTP_REQUIRED: ['otp_required', 'otp'],
  OTP_SEND_FAILED: ['otp_send_failed', 'otp'],
  OTP_INVALID: ['otp_invalid', 'otp'],
  OTP_NO_ACTIVE: ['no_active_otp', 'otp'],
  OTP_RATE_LIMIT: ['rate_limit_otp', 'rate_limit'],

  // Rate limiting
  RATE_LIMIT_GENERAL: ['rate_limit_general', 'rate_limit'],

  // Permissions
  PERM_SUPER_ADMIN_REQUIRED: ['permission_super_admin_required', 'permission'],
  PERM_ADMIN_REQUIRED: ['permission_admin_required', 'permission'],
  PERM_DENIED: ['permission_denied', 'permission'],
  PERM_SELLER_TIER_REQUIRED: ['permission_seller_tier_required', 'permission'],
  PERM_CANNOT_DELETE_SELF: ['permission_cannot_delete_self', 'permission'],
  PERM_CANNOT_CHANGE_OWN_ROLE: ['permission_cannot_change_own_role', 'permission'],
  PERM_LAST_SUPER_ADMIN: ['permission_last_super_admin', 'permission'],
  PERM_SUPER_ADMIN_DELETE: ['permission_super_admin_delete', 'permission'],
  PERM_REFERRAL_BLOCKED: ['permission_referral_blocked', 'permission'],
  PERM_REFERRAL_SUSPENDED: ['permission_referral_suspended', 'permission'],
  PERM_OWN_REFERRAL_CODE: ['permission_own_referral_code', 'permission'],
  PERM_INVITATION_ALREADY_ACCEPTED: ['permission_invitation_already_accepted', 'permission'],

  // Upload
  UPLOAD_CLOUD_FAILED: ['upload_cloud_failed', 'upload'],
  UPLOAD_FAILED: ['upload_failed', 'upload'],
  UPLOAD_FILE_TOO_LARGE: ['upload_file_too_large', 'upload'],
  UPLOAD_WRONG_TYPE: ['upload_wrong_type', 'upload'],
  UPLOAD_URL_NOT_LOADABLE: ['upload_url_not_loadable', 'upload'],
  UPLOAD_LIMIT_REACHED: ['upload_limit_reached', 'rate_limit'],

  // Payment
  PAY_PENDING_EXISTS: ['payment_pending_exists', 'payment'],
  PAY_ALREADY_PAID: ['payment_already_paid', 'payment'],
  PAY_SOLD_PRODUCT: ['payment_sold_product', 'payment'],
  PAY_UPGRADE_CLOSED: ['payment_upgrade_closed', 'payment'],
  PAY_INSUFFICIENT_BALANCE: ['payment_insufficient_balance', 'payment'],
  PAY_PENDING_WITHDRAWAL: ['payment_pending_withdrawal', 'payment'],
  PAY_MIN_WITHDRAWAL: ['payment_min_withdrawal', 'payment'],
  PAY_MAX_WITHDRAWAL: ['payment_max_withdrawal', 'payment'],
  PAY_INVALID_METHOD: ['payment_invalid_method', 'payment'],
  PAY_REJECTION_REASON_REQUIRED: ['payment_rejection_reason_required', 'payment'],
  PAY_REFERENCE_REQUIRED: ['payment_reference_required', 'payment'],

  // Product
  PROD_NOT_AVAILABLE: ['product_not_available', 'not_found'],
  PROD_NOT_FOUND: ['product_not_found', 'not_found'],
  PROD_LISTING_EXPIRED: ['product_listing_expired', 'not_found'],
  PROD_PRICE_MISMATCH: ['product_price_mismatch', 'validation'],
  PROD_CANNOT_EDIT: ['product_cannot_edit', 'validation'],
  PROD_CANNOT_DELETE: ['product_cannot_delete', 'validation'],
  PROD_MIN_IMAGES: ['product_min_images', 'validation'],

  // Validation
  VAL_PHONE_REQUIRED: ['validation_phone_required', 'validation'],
  VAL_PHONE_INVALID: ['validation_phone_invalid', 'validation'],
  VAL_PASSWORD_MIN: ['validation_password_min', 'validation'],
  VAL_NAME_REQUIRED: ['validation_name_required', 'validation'],
  VAL_TITLE_REQUIRED: ['validation_title_required', 'validation'],
  VAL_TITLE_TOO_LONG: ['validation_title_too_long', 'validation'],
  VAL_DESCRIPTION_TOO_LONG: ['validation_description_too_long', 'validation'],
  VAL_LOCATION_REQUIRED: ['validation_location_required', 'validation'],
  VAL_LOCATION_INVALID: ['location_invalid', 'validation'],
  VAL_TIER_INVALID: ['validation_tier_invalid', 'validation'],
  VAL_AMOUNT_POSITIVE: ['validation_amount_positive', 'validation'],
  VAL_PRICE_REQUIRED: ['validation_price_required', 'validation'],
  VAL_REQUIRED: ['validation_required', 'validation'],
  VAL_SETTING_KEY_UNKNOWN: ['validation_setting_key_unknown', 'validation'],
  VAL_SETTING_NUMERIC: ['validation_setting_numeric', 'validation'],
  VAL_SETTING_BOOLEAN: ['validation_setting_boolean', 'validation'],
  VAL_DATE_FORMAT: ['validation_date_format', 'validation'],
  VAL_PAYER_PHONE_REQUIRED: ['validation_payer_phone_required', 'validation'],
  VAL_MAX_IMAGES: ['validation_max_images', 'validation'],
  VAL_INVALID_IMAGE_URL: ['validation_invalid_image_url', 'validation'],
  VAL_PASSWORDS_MISMATCH: ['validation_passwords_mismatch', 'validation'],

  // Not found
  NF_USER: ['user_not_found', 'not_found'],
  NF_REFERRAL_CODE: ['referral_code_not_found', 'not_found'],
  NF_UPGRADE_REQUEST: ['upgrade_request_not_found', 'not_found'],
  NF_WITHDRAWAL: ['withdrawal_not_found', 'not_found'],
  NF_TRANSACTION: ['transaction_not_found', 'not_found'],
  NF_NOTIFICATION: ['notification_not_found', 'not_found'],
  NF_CATEGORY: ['category_not_found', 'not_found'],
  NF_ADMIN: ['admin_not_found', 'not_found'],

  // Server
  SRV_INTERNAL: ['server_error', 'server'],
  SRV_DATABASE: ['server_database_error', 'server'],

  // Network
  NET_OFFLINE: ['network_offline', 'network'],
  NET_TIMEOUT: ['network_timeout', 'network'],
};

// ── Substring-based fallback rules ───────────────────────────────────────────
// Each entry: [substring_to_match, error_key_in_errors_section, error_type]
// Rules are checked in order — first match wins.
// Substrings are lowercased before matching.

type Rule = [string, string, ErrorType];

const RULES: Rule[] = [
  // ── Rate limiting ──────────────────────────────────────────────────────────
  ['too many otp requests', 'rate_limit_otp', 'rate_limit'],
  ['too many requests', 'rate_limit_general', 'rate_limit'],
  ['upload limit reached', 'upload_limit_reached', 'rate_limit'],

  // ── OTP ───────────────────────────────────────────────────────────────────
  ['locked for 30 minutes', 'otp_locked', 'otp'],
  ['too many failed attempts', 'otp_locked', 'otp'],
  ['invalid or expired otp', 'otp_expired', 'otp'],
  ['otp code is required', 'otp_required', 'otp'],
  ['verification code', 'otp_required', 'otp'],
  ['failed to send otp', 'otp_send_failed', 'otp'],
  ['attempt(s) remaining', 'otp_invalid', 'otp'],
  ['invalid otp', 'otp_invalid', 'otp'],

  // ── Auth ──────────────────────────────────────────────────────────────────
  ['invalid credentials', 'auth_invalid_credentials', 'auth'],
  ['wrong password', 'auth_wrong_password', 'auth'],
  ['current password is incorrect', 'auth_wrong_password', 'auth'],
  ['account is inactive', 'auth_account_inactive', 'auth'],
  ['phone not verified', 'auth_phone_not_verified', 'auth'],
  ['complete otp verification', 'auth_phone_not_verified', 'auth'],
  ['session expired', 'auth_session_expired', 'auth'],
  ['invalid token', 'auth_token_invalid', 'auth'],
  ['missing or invalid authorization', 'auth_token_invalid', 'auth'],
  ['seller registration is currently closed', 'auth_registration_closed', 'auth'],
  ['buyer registration is currently closed', 'auth_buyer_registration_closed', 'auth'],
  ['phone number already registered', 'auth_phone_already_registered', 'auth'],
  ['already registered', 'auth_phone_already_registered', 'auth'],
  ['email address already in use', 'auth_email_already_in_use', 'auth'],
  ['account not found', 'auth_account_not_found', 'auth'],
  ['user not found', 'user_not_found', 'not_found'],

  // ── Permission ────────────────────────────────────────────────────────────
  ['insufficient admin privileges', 'permission_super_admin_required', 'permission'],
  ['super admin access required', 'permission_super_admin_required', 'permission'],
  ['admin access required', 'permission_admin_required', 'permission'],
  ['do not have permission', 'permission_denied', 'permission'],
  ['insufficient permissions', 'permission_denied', 'permission'],
  ['requires gold tier', 'permission_seller_tier_required', 'permission'],
  ['requires premium tier', 'permission_seller_tier_required', 'permission'],
  ['requires gold tier or higher', 'permission_seller_tier_required', 'permission'],
  ['requires premium tier or higher', 'permission_seller_tier_required', 'permission'],
  ['this feature requires', 'permission_seller_tier_required', 'permission'],
  ['cannot delete yourself', 'permission_cannot_delete_self', 'permission'],
  ['cannot change your own role', 'permission_cannot_change_own_role', 'permission'],
  ['cannot demote the last super admin', 'permission_last_super_admin', 'permission'],
  ['super admins cannot be deleted', 'permission_super_admin_delete', 'permission'],
  ['permanently blocked', 'permission_referral_blocked', 'permission'],
  ['referral code is not active', 'permission_referral_suspended', 'permission'],
  ['cannot use your own referral', 'permission_own_referral_code', 'permission'],
  ['already accepted an invitation', 'permission_invitation_already_accepted', 'permission'],
  ['invitation already accepted', 'permission_invitation_already_accepted', 'permission'],

  // ── Upload ────────────────────────────────────────────────────────────────
  ['cloud upload failed', 'upload_cloud_failed', 'upload'],
  ['upload failed', 'upload_failed', 'upload'],
  ['image must be smaller', 'upload_file_too_large', 'upload'],
  ['exceeds 10 mb', 'upload_file_too_large', 'upload'],
  ['only jpeg, png and webp', 'upload_wrong_type', 'upload'],
  ['only jpeg, png or webp', 'upload_wrong_type', 'upload'],
  ['not a jpeg, png or webp', 'upload_wrong_type', 'upload'],
  ['image file is required', 'upload_failed', 'upload'],
  ['could not load image', 'upload_url_not_loadable', 'upload'],
  ['image url could not be loaded', 'upload_url_not_loadable', 'upload'],

  // ── Payment / Upgrade ─────────────────────────────────────────────────────
  ['already have a pending upgrade request', 'payment_pending_exists', 'payment'],
  ['listing fee has already been recorded', 'payment_already_paid', 'payment'],
  ['cannot pay listing fee for a', 'payment_sold_product', 'payment'],
  ['buyer tier upgrades are currently closed', 'payment_upgrade_closed', 'payment'],
  ['insufficient balance', 'payment_insufficient_balance', 'payment'],
  ['already have a pending withdrawal', 'payment_pending_withdrawal', 'payment'],
  ['minimum withdrawal amount', 'payment_min_withdrawal', 'payment'],
  ['maximum withdrawal amount', 'payment_max_withdrawal', 'payment'],
  ['payment_method must be', 'payment_invalid_method', 'payment'],
  ['rejection reason is required', 'payment_rejection_reason_required', 'payment'],
  ['payment reference is required', 'payment_reference_required', 'payment'],
  ['payer_phone is required', 'validation_payer_phone_required', 'validation'],

  // ── Product / Listing ─────────────────────────────────────────────────────
  ['product is no longer available', 'product_not_available', 'not_found'],
  ['product not found', 'product_not_found', 'not_found'],
  ['listing has expired', 'product_listing_expired', 'not_found'],
  ['transaction amount', 'product_price_mismatch', 'validation'],
  ['does not match the listed price', 'product_price_mismatch', 'validation'],
  ['only pending or rejected', 'product_cannot_edit', 'validation'],
  ['cannot delete: product', 'product_cannot_delete', 'validation'],
  ['minimum 2 images required', 'product_min_images', 'validation'],
  ['too many images for your tier', 'validation_max_images', 'validation'],
  ['invalid image url', 'validation_invalid_image_url', 'validation'],
  ['only uploaded or cdn images', 'validation_invalid_image_url', 'validation'],

  // ── Validation ────────────────────────────────────────────────────────────
  ['phone is required', 'validation_phone_required', 'validation'],
  ['invalid phone', 'validation_phone_invalid', 'validation'],
  ['phone number', 'validation_phone_invalid', 'validation'],
  ['password must be at least 8', 'validation_password_min', 'validation'],
  ['name, phone and password', 'validation_name_required', 'validation'],
  ['title and a positive price', 'validation_title_required', 'validation'],
  ['title must be 200', 'validation_title_too_long', 'validation'],
  ['description must be 5000', 'validation_description_too_long', 'validation'],
  ['location (region) is required', 'validation_location_required', 'validation'],
  ['invalid location', 'location_invalid', 'validation'],
  ['invalid tier', 'validation_tier_invalid', 'validation'],
  ['must be gold or premium', 'validation_tier_invalid', 'validation'],
  ['amount must be greater than zero', 'validation_amount_positive', 'validation'],
  ['price must be greater than zero', 'validation_price_required', 'validation'],
  ['buyer_name and buyer_phone', 'validation_phone_required', 'validation'],
  ['product_id is required', 'validation_required', 'validation'],
  ['email and code are required', 'validation_required', 'validation'],
  ['email is required', 'validation_required', 'validation'],
  ['code is required', 'otp_required', 'otp'],
  ['invalid purpose', 'validation_required', 'validation'],
  ['invalid user_type', 'validation_required', 'validation'],
  ['invalid request body', 'validation_required', 'validation'],
  ['unknown setting key', 'validation_setting_key_unknown', 'validation'],
  ['must be a non-negative number', 'validation_setting_numeric', 'validation'],
  ['must be "true" or "false"', 'validation_setting_boolean', 'validation'],
  ['invalid date', 'validation_date_format', 'validation'],
  ['passwords do not match', 'validation_passwords_mismatch', 'validation'],

  // ── Not found ─────────────────────────────────────────────────────────────
  ['referral code not found', 'referral_code_not_found', 'not_found'],
  ['upgrade request not found', 'upgrade_request_not_found', 'not_found'],
  ['withdrawal not found', 'withdrawal_not_found', 'not_found'],
  ['transaction not found', 'transaction_not_found', 'not_found'],
  ['notification not found', 'notification_not_found', 'not_found'],
  ['category not found', 'category_not_found', 'not_found'],
  ['admin with id', 'admin_not_found', 'not_found'],
  ['seller with id', 'user_not_found', 'not_found'],
  ['no active otp', 'no_active_otp', 'otp'],

  // ── Network ───────────────────────────────────────────────────────────────
  ['network', 'network_offline', 'network'],
  ['failed to fetch', 'network_offline', 'network'],
  ['load failed', 'network_offline', 'network'],
  ['timed out', 'network_timeout', 'network'],
  ['timeout', 'network_timeout', 'network'],
  ['session expired — please log in', 'auth_session_expired', 'auth'],

  // ── Server ────────────────────────────────────────────────────────────────
  ['api error: 5', 'server_error', 'server'],
  ['internal error', 'server_error', 'server'],
  ['an internal error', 'server_error', 'server'],
];

// ── Public API ────────────────────────────────────────────────────────────────

/**
 * Translate a raw error message into a user-friendly localized string.
 *
 * Resolution order:
 *   1. Structured error code (from backend `code` field) — most precise
 *   2. Substring matching on the raw message — legacy fallback
 *   3. Raw message as-is — last resort
 *
 * @param raw     The raw error string from the API or caught exception
 * @param errors  The `t.errors` object from the current locale translation
 * @param code    Optional structured error code from backend response
 * @returns       { message, type, raw, code }
 */
export function translateError(
  raw: string | undefined | null,
  errors: Record<string, string>,
  code?: string | null
): TranslatedError {
  const input = (raw ?? '').trim();

  if (!input && !code) {
    return {
      message: errors['server_error'] ?? 'Something went wrong.',
      type: 'unknown',
      raw: input,
    };
  }

  // ── Strategy 1: Structured error code lookup ──────────────────────────────
  if (code) {
    const entry = CODE_MAP[code];
    if (entry) {
      const [key, type] = entry;
      let msg = errors[key] ?? errors['server_error'] ?? input;

      // Interpolate context values from the raw message
      msg = interpolate(msg, key, input);

      return { message: msg, type, raw: input, code };
    }
  }

  // ── Strategy 2: Substring matching ───────────────────────────────────────
  const lower = input.toLowerCase();

  for (const [substr, key, type] of RULES) {
    if (lower.includes(substr.toLowerCase())) {
      let msg = errors[key] ?? errors['server_error'] ?? input;
      msg = interpolate(msg, key, input);
      return { message: msg, type, raw: input, code: code ?? undefined };
    }
  }

  // ── Strategy 3: Return raw message ───────────────────────────────────────
  // Better than a generic "something went wrong" for unmapped messages
  return { message: input || (errors['server_error'] ?? 'Something went wrong.'), type: 'unknown', raw: input };
}

/**
 * Interpolate dynamic values into error message templates.
 * Handles: {n} for attempt counts, {tier} for tier names, {amount} for ETB amounts,
 *          {min}/{max} for image counts, {minutes} for lockout duration.
 */
function interpolate(msg: string, key: string, raw: string): string {
  // {n} — attempt count
  if (key === 'otp_invalid' || key === 'otp_attempts_remaining') {
    const match = raw.match(/(\d+)\s+attempt/i);
    if (match) msg = msg.replace('{n}', match[1]);
  }

  // {minutes} — lockout duration
  if (key === 'otp_locked') {
    const match = raw.match(/(\d+)\s+minute/i);
    if (match) msg = msg.replace('{minutes}', match[1]);
  }

  // {tier} — seller tier requirement
  if (key === 'permission_seller_tier_required') {
    const tierMatch = raw.match(/\b(gold|premium)\b/i);
    if (tierMatch) msg = msg.replace('{tier}', tierMatch[1]);
  }

  // {amount} — ETB balance
  if (key === 'payment_insufficient_balance') {
    const amtMatch = raw.match(/ETB\s*([\d,]+)/i);
    if (amtMatch) msg = msg.replace('{amount}', amtMatch[1]);
  }

  // {min} — minimum image count
  if (key === 'validation_min_images' || key === 'product_min_images') {
    const numMatch = raw.match(/(\d+)/);
    if (numMatch) msg = msg.replace('{min}', numMatch[1]);
  }

  // {max} — maximum image count
  if (key === 'validation_max_images') {
    const numMatch = raw.match(/(\d+)/);
    if (numMatch) msg = msg.replace('{max}', numMatch[1]);
  }

  return msg;
}

/**
 * Get the display label for an error type in the current locale.
 * Used for error banners that show a type badge (e.g. "Authentication Error").
 */
export function errorTypeLabel(type: ErrorType, errors: Record<string, string>): string {
  const key = `type_${type}`;
  return errors[key] ?? errors['type_unknown'] ?? 'Error';
}

/**
 * Convenience: returns just the translated message string.
 * Use when you don't need the type or raw value.
 */
export function te(raw: string | undefined | null, errors: Record<string, string>, code?: string | null): string {
  return translateError(raw, errors, code).message;
}

/**
 * Parse a backend API error response and extract both message and code.
 * Handles the standard { error: string, code?: string } response shape.
 */
export function parseApiError(err: unknown): { message: string; code?: string } {
  if (!err) return { message: 'Something went wrong.' };

  // Standard Error object
  if (err instanceof Error) {
    return { message: err.message };
  }

  // Object with message and optional code
  if (typeof err === 'object' && err !== null) {
    const e = err as Record<string, unknown>;
    return {
      message: typeof e.message === 'string' ? e.message : String(e.message ?? ''),
      code: typeof e.code === 'string' ? e.code : undefined,
    };
  }

  return { message: String(err) };
}
