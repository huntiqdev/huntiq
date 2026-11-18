import { fetchAPI } from './core';
import type { IdentityVerification } from '@/types';

export const identityApi = {
  verifyIdentity: (data: {
    fin: string;
    full_legal_name: string;
    phone: string;
    consent_given: boolean;
    auth_code?: string;
  }) =>
    fetchAPI<{ verification: IdentityVerification; message: string }>('/identity/verify', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  getMyIdentityStatus: () =>
    fetchAPI<{
      status: string;
      identity_verified: boolean;
      id_masked?: string;
      provider_code?: string;
      auth_reference?: string;
      phone_matched?: boolean;
      verified_at?: string;
    }>('/identity/status'),

  adminListPendingIdentity: (page: number = 1, limit: number = 50) =>
    fetchAPI<{ verifications: IdentityVerification[]; total: number }>(
      `/admin/identity/pending?page=${page}&limit=${limit}`
    ),

  adminResolveIdentityException: (id: number, status: 'verified' | 'rejected', reason: string) =>
    fetchAPI<{ message: string }>(`/admin/identity/${id}/resolve`, {
      method: 'PUT',
      body: JSON.stringify({ status, reason }),
    }),
};

