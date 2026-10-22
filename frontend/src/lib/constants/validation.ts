/**
 * Validation constants shared across the application.
 * These values MUST match backend validations for consistency.
 */

/** Minimum password length (matches backend pkg/validations/password.go) */
export const MIN_PASSWORD_LENGTH = 8;

/** Ethiopian phone number patterns */
export const ETHIO_PHONE_PATTERN = /^0[79]\d{8}$/;
export const ETHIO_PHONE_INTERNATIONAL_PATTERN = /^\+251[79]\d{8}$/;

/** Email validation pattern (matches backend pkg/validations/email.go) */
export const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;

/** Product validation constants - MUST MATCH BACKEND */
export const MIN_PRODUCT_TITLE_LENGTH = 3;
export const MAX_PRODUCT_TITLE_LENGTH = 200;  // Matches backend max_title_length setting
export const MIN_PRODUCT_DESCRIPTION_LENGTH = 10;
export const MAX_PRODUCT_DESCRIPTION_LENGTH = 5000;  // Matches backend max_description_length setting
export const MIN_PRODUCT_PRICE = 0.01;
export const MAX_PRODUCT_PRICE = 1_000_000;

/** 
 * Image validation - These are loaded from backend at runtime
 * Defaults here for initial render only
 * Actual limits vary by tier: Normal=2, Gold=7, Premium=20
 */
export const MIN_IMAGES_PER_PRODUCT = 2;  // Matches backend min_images_per_listing
