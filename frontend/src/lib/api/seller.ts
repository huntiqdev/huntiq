import { fetchAPI, authHeaders, API_URL, ApiError } from './core';

export const sellerApi = {
  getSellerStats: (token: string) =>
    fetchAPI<any>('/seller/stats', { headers: authHeaders(token) }),
  getSellerProductsDetailed: (token: string) =>
    fetchAPI<any[]>('/seller/products/detailed', { headers: authHeaders(token) }).then(r => r ?? []),
  createSellerProduct: (data: any, token: string) =>
    fetchAPI<any>('/seller/products', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  updateSellerProduct: (id: number, data: any, token: string) =>
    fetchAPI<any>(`/seller/products/${id}`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  deleteSellerProduct: (id: number, token: string) =>
    fetchAPI<any>(`/seller/products/${id}`, { method: 'DELETE', headers: authHeaders(token) }),
  getSellerTransactions: (token: string) =>
    fetchAPI<any[]>('/seller/transactions', { headers: authHeaders(token) }).then(r => r ?? []),
  getSellerOrders: (token: string) =>
    fetchAPI<any[]>('/seller/orders', { headers: authHeaders(token) }).then(r => r ?? []),
  updateSellerOrder: (id: number, data: { status: string; seller_note?: string }, token: string) =>
    fetchAPI<any>(`/seller/orders/${id}`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  getSellerProfile: (token: string) =>
    fetchAPI<any>('/seller/profile', { headers: authHeaders(token) }),
  updateSellerProfile: (data: any, token: string) =>
    fetchAPI<any>('/seller/profile', { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  sellerChangePassword: (data: { current_password: string; new_password: string }, token: string) =>
    fetchAPI<any>('/seller/password', { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  sellerPayListingFee: (data: {
    product_id: number; amount: number; payment_method: string;
    payer_phone: string; reference: string; notes?: string;
  }, token: string) =>
    fetchAPI<any>('/seller/listing-fee', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  sellerRequestUpgrade: (data: {
    requested_tier: 'gold' | 'premium'; payment_method: string;
    payer_phone: string; reference: string; amount: number; notes?: string;
  }, token: string) =>
    fetchAPI<any>('/seller/upgrade-tier', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  sellerGetNotifications: (token: string) =>
    fetchAPI<any[]>('/seller/notifications', { headers: authHeaders(token) }).then(r => r ?? []),
  sellerMarkNotificationRead: (id: number, token: string) =>
    fetchAPI<any>(`/seller/notifications/${id}/read`, { method: 'PUT', headers: authHeaders(token) }),
  sellerGetUpgradeRequests: (token: string) =>
    fetchAPI<any[]>('/seller/upgrade-tier', { headers: authHeaders(token) }).then(r => r ?? []),
  sellerGetReferralStats: (token: string) =>
    fetchAPI<any>('/seller/referral', { headers: authHeaders(token) }),
  sellerGetReferralCommissions: (token: string) =>
    fetchAPI<any[]>('/seller/referral/commissions', { headers: authHeaders(token) }).then(r => r ?? []),
  sellerGetReferralBalance: (token: string) =>
    fetchAPI<any>('/seller/referral/balance', { headers: authHeaders(token) }),
  sellerGetReferralWithdrawals: (token: string) =>
    fetchAPI<any[]>('/seller/referral/withdrawals', { headers: authHeaders(token) }).then(r => r ?? []),
  sellerRequestWithdrawal: (data: { amount: number; payment_method: string; payer_phone: string }, token: string) =>
    fetchAPI<any>('/seller/referral/withdraw', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  sellerAcceptReferralInvitation: (data: { referral_code: string }, token: string) =>
    fetchAPI<any>('/seller/referral/accept', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  sellerUpdatePreference: (data: { seller_preference: string }, token: string) =>
    fetchAPI<any>('/seller/preference', { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  sellerGetFeePreview: (price: number, tier: string, days: number, token: string) =>
    fetchAPI<any>(`/seller/fee-preview?price=${price}&tier=${tier}&days=${days}`, { headers: authHeaders(token) }),
  uploadImage: async (file: File, token: string): Promise<string> => {
    const form = new FormData();
    form.append('image', file);
    const res = await fetch(`${API_URL}/upload/image`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: form,
    });
    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      throw new ApiError(
        body.error || `Upload failed: ${res.status}`,
        body.code || 'UPLOAD_FAILED',
        res.status
      );
    }
    const json = await res.json();
    const result = ('data' in json) ? json.data : json;
    return result.url as string;
  },
};
