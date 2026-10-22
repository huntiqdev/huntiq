import { isValidPrice } from '@/lib/validations/common';
import { MIN_IMAGES } from '@/lib/constants/fees';
import type { ProductErrors } from '@/types';

export type { ProductErrors };

export function validateProduct(data: {
  title: string;
  price: number | string;
  imageURLs: string[];
  condition: string;
}): ProductErrors {
  const errors: ProductErrors = {};
  if (!data.title.trim()) errors.title = 'validation_title_required';
  const price = Number(data.price);
  if (!isValidPrice(price)) errors.price = 'validation_price_required';
  const validImages = data.imageURLs.filter(u => u.trim());
  if (validImages.length < MIN_IMAGES) errors.images = 'product_min_images';
  if (!data.condition) errors.condition = 'validation_required';
  return errors;
}
