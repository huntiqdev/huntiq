import { fetchAPI } from './core';
import type { Organization, OrganizationMember, AccountEntityType } from '@/types';

export const organizationApi = {
  registerOrganization: (data: {
    name: string;
    entity_type: AccountEntityType;
    registration_number?: string;
    tin_number?: string;
    official_address?: string;
    official_email?: string;
    official_phone?: string;
    publishing_authority_evidence?: string;
  }) =>
    fetchAPI<{ organization: Organization; message: string }>('/organizations', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  getMyOrganization: () =>
    fetchAPI<{
      has_organization: boolean;
      organization?: Organization;
      member?: OrganizationMember;
    }>('/organizations/me'),

  adminListOrganizations: (params?: {
    status?: string;
    entity_type?: string;
    page?: number;
    limit?: number;
  }) => {
    const query = new URLSearchParams();
    if (params?.status) query.set('status', params.status);
    if (params?.entity_type) query.set('entity_type', params.entity_type);
    if (params?.page) query.set('page', String(params.page));
    if (params?.limit) query.set('limit', String(params.limit));
    return fetchAPI<{ organizations: Organization[]; total: number }>(`/admin/organizations?${query.toString()}`);
  },

  adminReviewOrganization: (id: number, approved: boolean, notes: string) =>
    fetchAPI<{ message: string; organization: Organization }>(`/admin/organizations/${id}/review`, {
      method: 'PUT',
      body: JSON.stringify({ approved, notes }),
    }),
};

