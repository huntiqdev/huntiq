/**
 * image.ts — Image validation. Single source for all image rules.
 * Returns i18n keys from errors.* — resolved at render via resolveFieldError().
 */
import { isValidImageURL } from '@/lib/validations/common';

export interface ImageValidationResult {
  valid: boolean;
  /** i18n key */
  error?: string;
}

export interface ImageListResult {
  valid: boolean;
  /** i18n key */
  error?: string;
}

/** Validate a single image URL is an accepted format */
export function validateImageURL(url: string): ImageValidationResult {
  if (!url.trim()) return { valid: false, error: 'validation_required' };
  if (!isValidImageURL(url)) return { valid: false, error: 'validation_invalid_image_url' };
  return { valid: true };
}

/** Validate a list of image URLs against min/max constraints */
export function validateImageList(
  urls: string[],
  min: number,
  max: number
): ImageListResult {
  const filled = urls.filter(u => u.trim());
  if (filled.length < min) return { valid: false, error: 'validation_min_images' };
  if (filled.length > max) return { valid: false, error: 'validation_max_images' };
  const badUrl = filled.find(u => !isValidImageURL(u));
  if (badUrl) return { valid: false, error: 'validation_invalid_image_url' };
  return { valid: true };
}
