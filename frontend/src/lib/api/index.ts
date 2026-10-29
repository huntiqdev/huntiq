import { ApiError, fetchAPI, API_URL } from './core';
import { publicApi, PublicPaymentSettings } from './public';
import { authApi } from './auth';
import { sellerApi } from './seller';
import { buyerApi } from './buyer';
import { adminApi } from './admin';
import { referralApi } from './referral';
import { aiApi } from './ai';

export { ApiError, fetchAPI, API_URL };
export type { PublicPaymentSettings };
export { aiApi };

export const api = {
  ...publicApi,
  ...authApi,
  ...sellerApi,
  ...buyerApi,
  ...adminApi,
  ...referralApi,
  ...aiApi,
  // Generic methods for direct API calls
  get: <T>(endpoint: string, options?: RequestInit) => fetchAPI<T>(endpoint, { ...options, method: 'GET' }),
  post: <T>(endpoint: string, body?: any, options?: RequestInit) => fetchAPI<T>(endpoint, { ...options, method: 'POST', body: body ? JSON.stringify(body) : undefined }),
  put: <T>(endpoint: string, body?: any, options?: RequestInit) => fetchAPI<T>(endpoint, { ...options, method: 'PUT', body: body ? JSON.stringify(body) : undefined }),
  delete: <T>(endpoint: string, options?: RequestInit) => fetchAPI<T>(endpoint, { ...options, method: 'DELETE' }),
};
