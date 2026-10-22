import { fetchAPI, authHeaders } from './core';
import type { ReferralStats, ReferralBalance, CommissionEntry, WithdrawalRequest } from '@/types';

type ReferralRole = 'seller' | 'buyer';

function referralFor(role: ReferralRole) {
  const base = `/${role}/referral`;
  return {
    getStats:          (token: string)                    => fetchAPI<ReferralStats>(base, { headers: authHeaders(token) }),
    getCommissions:    (token: string, page = 1)          => fetchAPI<{ commissions: CommissionEntry[]; total: number; page: number }>(`${base}/commissions?page=${page}`, { headers: authHeaders(token) }),
    acceptInvitation:  (referralCode: string, token: string) =>
      fetchAPI<any>(`${base}/accept`, { method: 'POST', body: JSON.stringify({ referral_code: referralCode }), headers: authHeaders(token) }),
    getBalance:        (token: string)                    => fetchAPI<ReferralBalance>(`${base}/balance`, { headers: authHeaders(token) }),
    requestWithdrawal: (data: { amount: number; payment_method: string; payer_phone: string }, token: string) =>
      fetchAPI<{ id: number; message: string }>(`${base}/withdraw`, { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
    getWithdrawals:    (token: string)                    => fetchAPI<WithdrawalRequest[]>(`${base}/withdrawals`, { headers: authHeaders(token) }).then(r => r ?? []),
  };
}

const seller = referralFor('seller');
const buyer  = referralFor('buyer');

export const referralApi = {
  getReferralStats:              seller.getStats,
  getReferralCommissions:        seller.getCommissions,
  acceptReferralInvitation:      seller.acceptInvitation,
  getReferralBalance:            seller.getBalance,
  requestReferralWithdrawal:     seller.requestWithdrawal,
  getMyReferralWithdrawals:      seller.getWithdrawals,

  getBuyerReferralStats:             buyer.getStats,
  getBuyerReferralCommissions:       buyer.getCommissions,
  acceptBuyerReferralInvitation:     buyer.acceptInvitation,
  getBuyerReferralBalance:           buyer.getBalance,
  requestBuyerReferralWithdrawal:    buyer.requestWithdrawal,
  getMyBuyerReferralWithdrawals:     buyer.getWithdrawals,

  adminGetReferralOverview: (token: string) =>
    fetchAPI<any>('/admin/referral/overview', { headers: authHeaders(token) }),
  adminGetReferralList: (params: { page?: number; limit?: number; search?: string }, token: string) => {
    const q = '?' + new URLSearchParams(Object.fromEntries(Object.entries(params).map(([k, v]) => [k, String(v)]))).toString();
    return fetchAPI<any>(`/admin/referral/list${q}`, { headers: authHeaders(token) });
  },
  adminSetReferralSuspended: (userID: number, suspended: boolean, token: string) =>
    fetchAPI<any>(`/admin/referral/users/${userID}/suspend`, { method: 'PUT', body: JSON.stringify({ suspended }), headers: authHeaders(token) }),
  adminSetReferralBlocked: (userID: number, blocked: boolean, reason: string, token: string) =>
    fetchAPI<any>(`/admin/referral/users/${userID}/block`, { method: 'PUT', body: JSON.stringify({ blocked, reason }), headers: authHeaders(token) }),
  adminListReferralWithdrawals: (status: string, token: string, page = 1) =>
    fetchAPI<any>(`/admin/referral/withdrawals?status=${encodeURIComponent(status)}&page=${page}`, { headers: authHeaders(token) }),
  adminApproveReferralWithdrawal: (id: number, note: string, token: string) =>
    fetchAPI<any>(`/admin/referral/withdrawals/${id}/approve`, { method: 'PUT', body: JSON.stringify({ note }), headers: authHeaders(token) }),
  adminRejectReferralWithdrawal: (id: number, reason: string, token: string) =>
    fetchAPI<any>(`/admin/referral/withdrawals/${id}/reject`, { method: 'PUT', body: JSON.stringify({ reason }), headers: authHeaders(token) }),
  adminMarkReferralWithdrawalPaid: (id: number, reference: string, token: string) =>
    fetchAPI<any>(`/admin/referral/withdrawals/${id}/paid`, { method: 'PUT', body: JSON.stringify({ reference }), headers: authHeaders(token) }),
};
