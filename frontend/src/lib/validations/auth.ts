import { isValidPhone, isValidEmail } from '@/lib/validations/common';
import { MIN_PASSWORD_LENGTH } from '@/lib/constants/validation';
import type { LoginErrors, RegisterErrors } from '@/types';

export type { LoginErrors, RegisterErrors };

export function validateLogin(phone: string, password: string): LoginErrors {
  const errors: LoginErrors = {};
  if (!phone.trim()) errors.phone = 'validation_phone_required';
  else if (!isValidPhone(phone)) errors.phone = 'validation_phone_invalid';
  if (!password || password.length < MIN_PASSWORD_LENGTH) errors.password = 'validation_password_min';
  return errors;
}

export function validateRegister(data: {
  name: string;
  phone: string;
  email: string;
  password: string;
}): RegisterErrors {
  const errors: RegisterErrors = {};
  if (!data.name.trim()) errors.name = 'validation_name_required';
  if (!isValidPhone(data.phone)) errors.phone = 'validation_phone_invalid';
  if (data.email && !isValidEmail(data.email)) errors.email = 'validation_email_invalid';
  if (data.password.length < MIN_PASSWORD_LENGTH) errors.password = 'validation_password_min';
  return errors;
}
