/**
 * location.ts — Location validation. Single source for all location rules.
 * Returns i18n keys from errors.* — resolved at render via resolveFieldError().
 */

export interface LocationData {
  region?: string;
  zone?: string;
  subcity?: string;
  woreda?: string;
}

export interface LocationErrors {
  region?: string;
}

/** Region is the only required location field everywhere in the system */
export function validateLocation(data: LocationData): LocationErrors {
  const errors: LocationErrors = {};
  if (!data.region?.trim()) errors.region = 'validation_location_required';
  return errors;
}

/** Quick boolean check */
export function isValidLocation(data: LocationData): boolean {
  return !!data.region?.trim();
}
