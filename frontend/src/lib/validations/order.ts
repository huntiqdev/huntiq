import { validateEthioPhone } from '@/lib/validations/common';
import type { OrderErrors } from '@/types';

export type { OrderErrors };

export function validateOrderRequest(data: {
  name: string;
  phone: string;
  /** Pass true when the buyer is logged in and phone is pre-filled from profile */
  phonePreFilled?: boolean;
}): OrderErrors {
  const errors: OrderErrors = {};
  if (!data.phonePreFilled) {
    if (!data.name.trim()) errors.name = 'validation_name_required';
    const result = validateEthioPhone(data.phone);
    if (!result.valid) errors.phone = 'validation_phone_invalid';
  }
  return errors;
}

export function hasErrors(errors: Record<string, string | undefined>): boolean {
  return Object.values(errors).some(Boolean);
}
