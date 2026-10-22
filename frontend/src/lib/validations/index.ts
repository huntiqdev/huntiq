/**
 * validations/index.ts — Single entry point for all validation.
 *
 * RULE: validators return i18n keys (strings like 'validation_phone_invalid').
 * Components NEVER look up keys themselves. They call resolveFieldError() here.
 *
 * Resolution order for resolveFieldError:
 *   1. t.errors[key]   — structured error messages (primary)
 *   2. t.ui[key]       — ui fallback (some keys live there for legacy reasons)
 *   3. key itself      — last resort so nothing is ever blank
 */
import type { Translation } from '@/lib/i18n';

export { validateLogin, validateRegister } from './auth';
export { validateProduct } from './product';
export { validateOrderRequest, hasErrors } from './order';
export { validateImageURL, validateImageList } from './image';
export { validateLocation, isValidLocation } from './location';
export { validateEthioPhone, isValidPhone, isValidEmail, isValidPrice, isValidImageURL } from './common';

/**
 * Resolve a validation i18n key to a display string.
 * Pass the full `t` object from useLocale().
 *
 * @example
 *   const msg = resolveFieldError('validation_phone_invalid', t);
 *   // → "Invalid phone number. Use Ethiopian format: 09XXXXXXXX or 07XXXXXXXX."
 */
export function resolveFieldError(
  key: string | undefined,
  t: Translation
): string | undefined {
  if (!key) return undefined;
  const fromErrors = (t.errors as Record<string, string>)[key];
  if (fromErrors) return fromErrors;
  const fromUi = (t.ui as unknown as Record<string, string>)[key];
  if (fromUi && typeof fromUi === 'string') return fromUi;
  return key; // last resort — never blank
}

/**
 * Resolve a whole errors object to display strings in one call.
 * Useful for form error maps.
 *
 * @example
 *   const display = resolveErrors(fieldErrors, t);
 *   // → { phone: "Invalid phone number...", name: "Name is required." }
 */
export function resolveErrors<K extends string>(
  errors: Partial<Record<K, string>>,
  t: Translation
): Partial<Record<K, string>> {
  const out: Partial<Record<K, string>> = {};
  for (const key in errors) {
    const raw = errors[key];
    out[key] = resolveFieldError(raw, t);
  }
  return out;
}
