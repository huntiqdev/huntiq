import { fetchAPI, authHeaders } from './core';

export const buyerApi = {
  buyerUpgradeTier: (data: { tier: string; reference: string; payer_phone: string; payment_method: string; amount: number; duration_days: number }, token: string) =>
    fetchAPI<{ message: string }>('/buyer/upgrade-tier', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  getBuyerProfile: (token: string) =>
    fetchAPI<any>('/buyer/profile', { headers: authHeaders(token) }),
  updateBuyerProfile: (data: any, token: string) =>
    fetchAPI<any>('/buyer/profile', { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  buyerChangePassword: (data: { current_password: string; new_password: string }, token: string) =>
    fetchAPI<any>('/buyer/password', { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  getBuyerOrders: (token: string) =>
    fetchAPI<any[]>('/buyer/orders', { headers: authHeaders(token) }).then(r => r ?? []),
  buyerGetNotifications: (token: string) =>
    fetchAPI<any[]>('/buyer/notifications', { headers: authHeaders(token) }).then(r => r ?? []),
  buyerMarkNotificationRead: (id: number, token: string) =>
    fetchAPI<any>(`/buyer/notifications/${id}/read`, { method: 'PUT', headers: authHeaders(token) }),
  buyerGetUpgradeRequests: (token: string) =>
    fetchAPI<any[]>('/buyer/upgrade-tier', { headers: authHeaders(token) }).then(r => r ?? []),
  buyerGetReferralStats: (token: string) =>
    fetchAPI<any>('/buyer/referral', { headers: authHeaders(token) }),
  buyerGetReferralCommissions: (token: string) =>
    fetchAPI<any[]>('/buyer/referral/commissions', { headers: authHeaders(token) }).then(r => r ?? []),
  buyerGetReferralBalance: (token: string) =>
    fetchAPI<any>('/buyer/referral/balance', { headers: authHeaders(token) }),
  buyerGetReferralWithdrawals: (token: string) =>
    fetchAPI<any[]>('/buyer/referral/withdrawals', { headers: authHeaders(token) }).then(r => r ?? []),
  buyerRequestWithdrawal: (data: { amount: number; payment_method: string; payer_phone: string }, token: string) =>
    fetchAPI<any>('/buyer/referral/withdraw', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  buyerAcceptReferralInvitation: (data: { referral_code: string }, token: string) =>
    fetchAPI<any>('/buyer/referral/accept', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
};
