/**
 * Ethiopian phone number validation.
 *
 * Supported formats (all normalize to +251XXXXXXXXX):
 *   Ethio Telecom (Telebirr):  09XXXXXXXX  or  +2519XXXXXXXX
 *   Safaricom Ethiopia:        07XXXXXXXX  or  +2517XXXXXXXX
 *
 * Returns the detected network or null if invalid.
 */
import { ETHIO_PHONE_PATTERN, EMAIL_PATTERN } from '@/lib/constants/validation';

export type EthioNetwork = 'telebirr' | 'safaricom' | 'unknown';

export interface PhoneValidation {
  valid: boolean;
  normalized: string;   // always +251XXXXXXXXX
  network: EthioNetwork;
  error?: string;
}

export function validateEthioPhone(raw: string): PhoneValidation {
  const cleaned = raw.replace(/[\s\-().]/g, '');

  // Normalize to local 10-digit form first
  let local = cleaned;
  if (local.startsWith('+251')) {
    local = '0' + local.slice(4);
  } else if (local.startsWith('251')) {
    local = '0' + local.slice(3);
  }

  // Must be exactly 10 digits starting with 09 or 07
  if (!ETHIO_PHONE_PATTERN.test(local)) {
    return {
      valid: false,
      normalized: '',
      network: 'unknown',
      error: 'Enter a valid Ethiopian phone number (09XXXXXXXX or 07XXXXXXXX)',
    };
  }

  const normalized = '+251' + local.slice(1);
  const network: EthioNetwork = local.startsWith('07') ? 'safaricom' : 'telebirr';

  return { valid: true, normalized, network };
}

/** Quick boolean check — use validateEthioPhone for details */
export function isValidEthioPhone(phone: string): boolean {
  return validateEthioPhone(phone).valid;
}

/** Legacy alias — kept for backward compat */
export function isValidPhone(phone: string): boolean {
  return isValidEthioPhone(phone);
}

/** Email validation (pattern matches backend pkg/validations/email.go) */
export function isValidEmail(email: string): boolean {
  return EMAIL_PATTERN.test(email.trim());
}

export function isValidPrice(price: number): boolean {
  return price > 0 && isFinite(price);
}

/** Validates that a URL is a safe image path accepted by the backend.
 *  Accepts: relative /uploads/ paths, Cloudinary URLs, and other trusted CDN prefixes.
 *  Rejects: javascript:, data: URIs, and arbitrary external URLs.
 */
export function isValidImageURL(url: string): boolean {
  if (typeof url !== 'string' || url.length < 5) return false;
  if (url.startsWith('/uploads/')) return true;
  const trustedPrefixes = [
    'https://picsum.photos/',
    'https://images.unsplash.com/',
    'https://lh3.googleusercontent.com/',
    'https://encrypted-tbn0.gstatic.com/',
    'https://res.cloudinary.com/',
  ];
  return trustedPrefixes.some(p => url.startsWith(p));
}
