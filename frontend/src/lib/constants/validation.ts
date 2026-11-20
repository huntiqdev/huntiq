/**
 * Validation rules and dynamic privilege engine.
 * Admin settings and user tiers (normal, gold, premium) determine limits at runtime.
 */

/** Minimum password length (matches backend pkg/validations/password.go) */
export const MIN_PASSWORD_LENGTH = 8;

/** Ethiopian phone number patterns */
export const ETHIO_PHONE_PATTERN = /^0[79]\d{8}$/;
export const ETHIO_PHONE_INTERNATIONAL_PATTERN = /^\+251[79]\d{8}$/;

/** Email validation pattern (matches backend pkg/validations/email.go) */
export const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;

/** Base fallback validation constants */
export const MIN_PRODUCT_TITLE_LENGTH = 3;
export const MAX_PRODUCT_TITLE_LENGTH = 200;
export const MIN_PRODUCT_DESCRIPTION_LENGTH = 10;
export const MAX_PRODUCT_DESCRIPTION_LENGTH = 5000;
export const MIN_PRODUCT_PRICE = 0.01;
export const MAX_PRODUCT_PRICE = 999_999_999.99;
export const MIN_IMAGES_PER_PRODUCT = 2;

export type UserTier = 'normal' | 'gold' | 'premium';

export interface TierPrivilegeRules {
  tier: UserTier;
  maxImages: number;
  minImages: number;
  maxLocations: number;
  maxTitleLength: number;
  maxDescriptionLength: number;
  maxImageSizeMB: number;
  hasGpsAccess: boolean;
  hasPriorityListing: boolean;
  hasBulkDiscounts: boolean;
}

/**
 * Resolves dynamic contextual validation rules based on user privilege tier
 * and platform_settings configured by admin in real-time.
 */
export function getTierValidationRules(
  tier: UserTier = 'normal',
  settings?: Record<string, any> | null
): TierPrivilegeRules {
  const defaultImages: Record<UserTier, number> = { normal: 2, gold: 7, premium: 20 };
  const defaultLocations: Record<UserTier, number> = { normal: 1, gold: 3, premium: 10 };

  const maxImages = settings?.[`${tier}_max_images`]
    ? parseInt(settings[`${tier}_max_images`], 10)
    : defaultImages[tier];

  const minImages = settings?.min_images_per_listing
    ? parseInt(settings.min_images_per_listing, 10)
    : 2;

  const maxLocations = settings?.[`${tier}_max_locations`]
    ? parseInt(settings[`${tier}_max_locations`], 10)
    : defaultLocations[tier];

  const maxTitleLength = settings?.max_title_length
    ? parseInt(settings.max_title_length, 10)
    : MAX_PRODUCT_TITLE_LENGTH;

  const maxDescriptionLength = settings?.max_description_length
    ? parseInt(settings.max_description_length, 10)
    : MAX_PRODUCT_DESCRIPTION_LENGTH;

  const maxImageSizeMB = settings?.max_image_size_mb
    ? parseInt(settings.max_image_size_mb, 10)
    : 5; // Strictly 5MB max

  return {
    tier,
    maxImages: Math.max(minImages, maxImages || 2),
    minImages,
    maxLocations: Math.max(1, maxLocations || 1),
    maxTitleLength,
    maxDescriptionLength,
    maxImageSizeMB,
    hasGpsAccess: true, // Free for all tiers
    hasPriorityListing: tier === 'premium',
    hasBulkDiscounts: tier === 'gold' || tier === 'premium',
  };
}
