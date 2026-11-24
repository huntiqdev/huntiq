import { ApiError, fetchAPI, API_URL } from './core';
import { publicApi, PublicPaymentSettings } from './public';
import { authApi } from './auth';
import { sellerApi } from './seller';
import { buyerApi } from './buyer';
import { adminApi } from './admin';
import { referralApi } from './referral';
import { classificationApi } from './classification';
import { identityApi } from './identity';
import { organizationApi } from './organization';
import { paymentApi } from './payment';
export { ApiError, fetchAPI, API_URL };
export type { PublicPaymentSettings, CategoryIntentData, CategoryChildChip, CategoryIntentCandidate, IntentQuickFilters } from './public';
export * from './classification';
export * from './identity';
export * from './organization';
export * from './payment';

export const api = {
  ...publicApi,
  ...authApi,
  ...sellerApi,
  ...buyerApi,
  ...adminApi,
  ...referralApi,
  ...classificationApi,
  ...identityApi,
  ...organizationApi,
  ...paymentApi,
  // Generic methods for direct API calls
  get: <T>(endpoint: string, options?: RequestInit) => fetchAPI<T>(endpoint, { ...options, method: 'GET' }),
  post: <T>(endpoint: string, body?: any, options?: RequestInit) => fetchAPI<T>(endpoint, { ...options, method: 'POST', body: body ? JSON.stringify(body) : undefined }),
  put: <T>(endpoint: string, body?: any, options?: RequestInit) => fetchAPI<T>(endpoint, { ...options, method: 'PUT', body: body ? JSON.stringify(body) : undefined }),
  delete: <T>(endpoint: string, options?: RequestInit) => fetchAPI<T>(endpoint, { ...options, method: 'DELETE' }),
};
