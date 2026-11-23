import { fetchAPI, authHeaders } from './core';

export interface PaymentOrder {
  id: number;
  order_id: string;
  user_id: number;
  user_name?: string;
  user_phone?: string;
  provider: string;
  environment: 'demo' | 'production';
  order_type: 'tier_upgrade' | 'listing_fee' | 'membership';
  target_tier?: string;
  user_role: string;
  duration_days: number;
  product_id?: number;
  product_title?: string;
  amount: number;
  currency: string;
  title: string;
  status: 'CREATED' | 'PENDING' | 'PROCESSING' | 'SUCCESSFUL' | 'FAILED' | 'CANCELLED' | 'EXPIRED';
  payer_phone?: string;
  provider_tx_id?: string;
  checkout_url?: string;
  is_demo: boolean;
  otp_verified: boolean;
  failure_reason?: string;
  created_at: string;
  updated_at: string;
}

export interface PaymentPricing {
  pricing: {
    seller_gold_price_etb: number;
    seller_premium_price_etb: number;
    buyer_gold_price_etb: number;
    buyer_premium_price_etb: number;
    listing_fee_etb: number;
  };
  operating_mode: 'demo' | 'production';
}

export interface PaymentProviderConfig {
  id: number;
  provider_code: string;
  provider_name: string;
  environment_mode: string;
  is_enabled: boolean;
  api_endpoint: string;
  app_id: string;
  app_key_masked: string;
  short_code: string;
  notify_url: string;
  return_url: string;
}

export const paymentApi = {
  getPaymentPricing: () =>
    fetchAPI<PaymentPricing>('/payments/pricing'),

  checkoutPayment: (
    data: {
      order_type?: string;
      target_tier?: string;
      duration_days?: number;
      product_id?: number;
      provider?: string;
      payer_phone?: string;
      otp_verified?: boolean;
    },
    token: string
  ) =>
    fetchAPI<PaymentOrder>('/payments/checkout', {
      method: 'POST',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  getPaymentOrder: (orderId: string, token: string) =>
    fetchAPI<PaymentOrder>(`/payments/orders/${orderId}`, {
      headers: authHeaders(token),
    }),

  verifyPaymentOrder: (orderId: string, token: string) =>
    fetchAPI<PaymentOrder>(`/payments/verify/${orderId}`, {
      method: 'POST',
      headers: authHeaders(token),
    }),

  simulateDemoPayment: (data: { order_id: string; action: string }, token: string) =>
    fetchAPI<PaymentOrder>('/payments/demo/simulate', {
      method: 'POST',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  sendMembershipOTP: (data: { phone: string; target_tier: string }, token: string) =>
    fetchAPI<{ message: string }>('/membership/otp/send', {
      method: 'POST',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  verifyMembershipOTP: (data: { phone: string; code: string }, token: string) =>
    fetchAPI<{ verified: boolean; message: string }>('/membership/otp/verify', {
      method: 'POST',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  // Admin APIs
  adminListPaymentOrders: (params: Record<string, any>, token: string) => {
    const qs = new URLSearchParams();
    Object.entries(params).forEach(([k, v]) => {
      if (v !== undefined && v !== null && v !== '') qs.set(k, String(v));
    });
    return fetchAPI<{ orders: PaymentOrder[]; total: number; limit: number; offset: number }>(
      `/admin/payments/orders?${qs.toString()}`,
      { headers: authHeaders(token) }
    );
  },

  adminGetPaymentOrder: (id: string | number, token: string) =>
    fetchAPI<PaymentOrder>(`/admin/payments/orders/${id}`, {
      headers: authHeaders(token),
    }),

  adminGetPaymentProviders: (token: string) =>
    fetchAPI<{ providers: PaymentProviderConfig[]; operating_mode: 'demo' | 'production' }>(
      '/admin/payments/providers',
      { headers: authHeaders(token) }
    ),

  adminUpdatePaymentProvider: (code: string, data: any, token: string) =>
    fetchAPI<{ message: string }>(`/admin/payments/providers/${code}`, {
      method: 'PUT',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  adminSwitchPaymentMode: (mode: 'demo' | 'production', token: string) =>
    fetchAPI<{ message: string; operating_mode: 'demo' | 'production' }>('/admin/payments/switch-mode', {
      method: 'POST',
      body: JSON.stringify({ mode }),
      headers: authHeaders(token),
    }),
};
