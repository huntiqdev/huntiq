import { fetchAPI, authHeaders } from './core';

export const adminApi = {
  adminGetPending: (token: string) =>
    fetchAPI<any[]>('/admin/products/pending', { headers: authHeaders(token) }).then(r => r ?? []),
  adminGetAllProducts: (token: string) =>
    fetchAPI<any>('/admin/products', { headers: authHeaders(token) }).then(r => Array.isArray(r) ? r : (r?.data ?? [])),
  adminGetProductsByStatus: (status: string, token: string) =>
    fetchAPI<any>(`/admin/products?status=${encodeURIComponent(status)}`, { headers: authHeaders(token) }).then(r => Array.isArray(r) ? r : (r?.data ?? [])),
  adminReviewProduct: (id: number, data: { status: string; notes: string }, token: string) =>
    fetchAPI<any>(`/admin/products/${id}/review`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminDeleteProduct: (id: number, token: string) =>
    fetchAPI<any>(`/admin/products/${id}`, { method: 'DELETE', headers: authHeaders(token) }),
  adminGetSellers: (token: string) =>
    fetchAPI<any>('/admin/sellers', { headers: authHeaders(token) }).then(r => Array.isArray(r) ? r : (r?.data ?? [])),
  adminSetSellerActive: (id: number, active: boolean, reason: string, token: string) =>
    fetchAPI<any>(`/admin/sellers/${id}/active`, { method: 'PUT', body: JSON.stringify({ active, reason }), headers: authHeaders(token) }),
  adminGetTransactions: (token: string) =>
    fetchAPI<any>('/admin/transactions', { headers: authHeaders(token) }).then(r => Array.isArray(r) ? r : (r?.data ?? [])),
  adminUpdateTransactionStatus: (id: number, status: string, token: string) =>
    fetchAPI<any>(`/admin/transactions/${id}/status`, { method: 'PUT', body: JSON.stringify({ status }), headers: authHeaders(token) }),
  adminGetLocations: (token: string) =>
    fetchAPI<any[]>('/admin/locations', { headers: authHeaders(token) }).then(r => r ?? []),
  adminAddLocation: (data: any, token: string) =>
    fetchAPI<any>('/admin/locations', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminUpdateLocation: (id: number, data: any, token: string) =>
    fetchAPI<any>(`/admin/locations/${id}`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminDeleteLocation: (id: number, token: string) =>
    fetchAPI<any>(`/admin/locations/${id}`, { method: 'DELETE', headers: authHeaders(token) }),
  adminAnalyticsOverview: (token: string) =>
    fetchAPI<any>('/admin/analytics/overview', { headers: authHeaders(token) }),
  adminAnalyticsHeatmap: (token: string) =>
    fetchAPI<any[]>('/admin/analytics/heatmap', { headers: authHeaders(token) }),
  adminAnalyticsMonthly: (months: number, token: string) =>
    fetchAPI<any[]>(`/admin/analytics/monthly?months=${months}`, { headers: authHeaders(token) }),
  adminAnalyticsGantt: (limit: number, token: string) =>
    fetchAPI<any[]>(`/admin/analytics/gantt?limit=${limit}`, { headers: authHeaders(token) }),
  adminAnalyticsFeeRevenue: (token: string) =>
    fetchAPI<any[]>('/admin/analytics/fee-revenue', { headers: authHeaders(token) }),
  adminGetSettings: (token: string) =>
    fetchAPI<Record<string,string>>('/admin/settings', { headers: authHeaders(token) }),
  adminUpdateSettings: (data: Record<string,string>, token: string) =>
    fetchAPI<any>('/admin/settings', { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminUpdateSettingsWithReason: (settings: Record<string,string>, reason: string, token: string) =>
    fetchAPI<any>('/admin/settings/with-reason', { method: 'PUT', body: JSON.stringify({ settings, reason }), headers: authHeaders(token) }),
  adminGetSettingsChangeLog: (key: string, token: string) =>
    fetchAPI<any[]>(`/admin/settings/log${key ? `?key=${encodeURIComponent(key)}` : ''}`, { headers: authHeaders(token) }),
  adminGetNotifications: (token: string) =>
    fetchAPI<any[]>('/admin/notifications', { headers: authHeaders(token) }),
  adminMarkNotificationRead: (id: number, token: string) =>
    fetchAPI<any>(`/admin/notifications/${id}/read`, { method: 'PUT', headers: authHeaders(token) }),
  adminCreateNotification: (data: any, token: string) =>
    fetchAPI<any>('/admin/notifications', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminProfitSummary: (token: string) =>
    fetchAPI<any>('/admin/profit/summary', { headers: authHeaders(token) }),
  adminProfitByTier: (token: string) =>
    fetchAPI<any[]>('/admin/profit/by-tier', { headers: authHeaders(token) }),
  adminProfitMonthly: (token: string) =>
    fetchAPI<any[]>('/admin/profit/monthly', { headers: authHeaders(token) }),
  adminProfitWeekly: (token: string) =>
    fetchAPI<any[]>('/admin/profit/weekly', { headers: authHeaders(token) }),
  adminProfitAnnual: (token: string) =>
    fetchAPI<any[]>('/admin/profit/annual', { headers: authHeaders(token) }),
  adminGetCategories: (token: string) =>
    fetchAPI<any[]>('/admin/categories', { headers: authHeaders(token) }),
  adminCreateCategory: (data: any, token: string) =>
    fetchAPI<any>('/admin/categories', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminUpdateCategory: (id: number, data: any, token: string) =>
    fetchAPI<any>(`/admin/categories/${id}`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminDeleteCategory: (id: number, token: string) =>
    fetchAPI<any>(`/admin/categories/${id}`, { method: 'DELETE', headers: authHeaders(token) }),
  adminGetBanners: (token: string) =>
    fetchAPI<any[]>('/admin/banners', { headers: authHeaders(token) }).then(r => r ?? []),
  adminCreateBanner: (data: any, token: string) =>
    fetchAPI<any>('/admin/banners', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminUpdateBanner: (id: number, data: any, token: string) =>
    fetchAPI<any>(`/admin/banners/${id}`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminDeleteBanner: (id: number, token: string) =>
    fetchAPI<any>(`/admin/banners/${id}`, { method: 'DELETE', headers: authHeaders(token) }),
  adminListAdmins: (token: string) =>
    fetchAPI<any[]>('/admin/admins', { headers: authHeaders(token) }).then(r => r ?? []),
  adminCreateAdmin: (data: { name: string; phone: string; email: string; password: string; admin_role: string }, token: string) =>
    fetchAPI<any>('/admin/admins', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminDeleteAdmin: (id: number, token: string) =>
    fetchAPI<any>(`/admin/admins/${id}`, { method: 'DELETE', headers: authHeaders(token) }),
  adminGetAdminPermissions: (id: number, token: string) =>
    fetchAPI<string[]>(`/admin/admins/${id}/permissions`, { headers: authHeaders(token) }),
  adminSetAdminPermissions: (id: number, permissions: string[], token: string) =>
    fetchAPI<any>(`/admin/admins/${id}/permissions`, { method: 'PUT', body: JSON.stringify({ permissions }), headers: authHeaders(token) }),
  adminGetMyPermissions: (token: string) =>
    fetchAPI<{ permissions: string[]; is_super: boolean }>('/admin/me/permissions', { headers: authHeaders(token) }),
  adminGetMe: (token: string) =>
    fetchAPI<{ id: number; name: string; phone: string; email: string; admin_role: string }>('/admin/me', { headers: authHeaders(token) }),
  adminChangePassword: (data: { current_password: string; new_password: string }, token: string) =>
    fetchAPI<any>('/admin/password', { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminGetLogs: (params: Record<string, any>, token: string) => {
    const q = '?' + new URLSearchParams(params).toString();
    return fetchAPI<{ logs: any[]; total: number }>(`/admin/logs${q}`, { headers: authHeaders(token) });
  },
  adminDeleteLogs: (before: string, token: string) => {
    const q = before ? `?before=${encodeURIComponent(before)}` : '';
    return fetchAPI<{ deleted: number }>(`/admin/logs${q}`, { method: 'DELETE', headers: authHeaders(token) });
  },
  adminGetOTPLogs: (params: Record<string, any>, token: string) => {
    const q = '?' + new URLSearchParams(params).toString();
    return fetchAPI<{ logs: any[]; total: number }>(`/admin/otp-logs${q}`, { headers: authHeaders(token) });
  },
  adminDeleteOTPLogs: (status: 'done' | 'all', token: string) =>
    fetchAPI<{ deleted: number }>(`/admin/otp-logs?status=${status}`, { method: 'DELETE', headers: authHeaders(token) }),
  adminGetActivityLogs: (params: Record<string, any>, token: string) => {
    const q = '?' + new URLSearchParams(params).toString();
    return fetchAPI<{ logs: any[]; total: number }>(`/admin/activity-logs${q}`, { headers: authHeaders(token) });
  },
  adminSetSellerTier: (userId: number, tier: string, reason: string, token: string) =>
    fetchAPI<any>(`/admin/users/${userId}/seller-tier`, { method: 'PUT', body: JSON.stringify({ tier, reason }), headers: authHeaders(token) }),
  adminSetBuyerTier: (userId: number, tier: string, reason: string, token: string) =>
    fetchAPI<any>(`/admin/users/${userId}/buyer-tier`, { method: 'PUT', body: JSON.stringify({ tier, reason }), headers: authHeaders(token) }),
  adminGetTierHistory: (userId: number, token: string) =>
    fetchAPI<any[]>(`/admin/users/${userId}/tier-history`, { headers: authHeaders(token) }).then(r => r ?? []),
  adminGetPremiumUsers: (params: Record<string, string>, token: string) => {
    const q = '?' + new URLSearchParams(params).toString();
    return fetchAPI<any[]>(`/admin/premium-users${q}`, { headers: authHeaders(token) }).then(r => r ?? []);
  },
  adminListUpgradeRequests: (status: string, token: string) =>
    fetchAPI<any[]>(`/admin/upgrade-requests?status=${encodeURIComponent(status)}`, { headers: authHeaders(token) }).then(r => r ?? []),
  adminApproveUpgrade: (id: number, token: string) =>
    fetchAPI<any>(`/admin/upgrade-requests/${id}/approve`, { method: 'PUT', headers: authHeaders(token) }),
  adminRejectUpgrade: (id: number, reason: string, token: string) =>
    fetchAPI<any>(`/admin/upgrade-requests/${id}/reject`, { method: 'PUT', body: JSON.stringify({ reason }), headers: authHeaders(token) }),
  adminListBuyerUpgradeRequests: (status: string, token: string) =>
    fetchAPI<{ requests: any[]; total: number }>(`/admin/buyer-upgrade-requests?status=${encodeURIComponent(status)}`, { headers: authHeaders(token) }),
  adminApproveBuyerUpgrade: (id: number, token: string) =>
    fetchAPI<any>(`/admin/buyer-upgrade-requests/${id}/approve`, { method: 'PUT', headers: authHeaders(token) }),
  adminRejectBuyerUpgrade: (id: number, reason: string, token: string) =>
    fetchAPI<any>(`/admin/buyer-upgrade-requests/${id}/reject`, { method: 'PUT', body: JSON.stringify({ reason }), headers: authHeaders(token) }),
  adminCategoryPopularity: (days: number, token: string) =>
    fetchAPI<any[]>(`/admin/categories/popularity?days=${days}&limit=50`, { headers: authHeaders(token) }).then(r => r ?? []),
  adminGetSystemHealth: (token: string) =>
    fetchAPI<any>('/admin/system/health', { headers: authHeaders(token) }),
  adminGetPendingFeePayments: (token: string) =>
    fetchAPI<any[]>('/admin/fee-payments/pending', { headers: authHeaders(token) }).then(r => r ?? []),
  adminVerifyFeePayment: (id: number, token: string) =>
    fetchAPI<any>(`/admin/fee-payments/${id}/verify`, { method: 'PUT', headers: authHeaders(token) }),
  adminRejectFeePayment: (id: number, reason: string, token: string) =>
    fetchAPI<any>(`/admin/fee-payments/${id}/reject`, { method: 'PUT', body: JSON.stringify({ reason }), headers: authHeaders(token) }),
  adminIssueWarning: (id: number, reason: string, token: string) =>
    fetchAPI<any>(`/admin/users/${id}/warn`, { method: 'POST', body: JSON.stringify({ reason }), headers: authHeaders(token) }),
  adminBanUser: (id: number, reason: string, token: string) =>
    fetchAPI<any>(`/admin/users/${id}/ban`, { method: 'PUT', body: JSON.stringify({ reason }), headers: authHeaders(token) }),
  
  // Referral Admin Endpoints
  adminGetReferralOverview: (token: string) =>
    fetchAPI<any>('/admin/referral/overview', { headers: authHeaders(token) }),
  adminGetReferralList: (token: string) =>
    fetchAPI<any[]>('/admin/referral/list', { headers: authHeaders(token) }).then(r => r ?? []),
  adminSuspendReferralUser: (id: number, suspended: boolean, token: string) =>
    fetchAPI<any>(`/admin/referral/users/${id}/suspend`, { method: 'PUT', body: JSON.stringify({ suspended }), headers: authHeaders(token) }),
  adminBlockReferralUser: (id: number, blocked: boolean, token: string) =>
    fetchAPI<any>(`/admin/referral/users/${id}/block`, { method: 'PUT', body: JSON.stringify({ blocked }), headers: authHeaders(token) }),
  adminListReferralWithdrawals: (status: string, token: string) =>
    fetchAPI<any[]>(`/admin/referral/withdrawals?status=${encodeURIComponent(status)}`, { headers: authHeaders(token) }).then(r => r ?? []),
  adminApproveReferralWithdrawal: (id: number, token: string) =>
    fetchAPI<any>(`/admin/referral/withdrawals/${id}/approve`, { method: 'PUT', headers: authHeaders(token) }),
  adminRejectReferralWithdrawal: (id: number, reason: string, token: string) =>
    fetchAPI<any>(`/admin/referral/withdrawals/${id}/reject`, { method: 'PUT', body: JSON.stringify({ reason }), headers: authHeaders(token) }),
  adminMarkReferralWithdrawalPaid: (id: number, token: string) =>
    fetchAPI<any>(`/admin/referral/withdrawals/${id}/paid`, { method: 'PUT', headers: authHeaders(token) }),

  // ── Operations Monitoring ────────────────────────────────────────────────
  adminGetOperationsHealth: (token: string) =>
    fetchAPI<any>('/admin/operations/health', { headers: authHeaders(token) }),
  adminGetJobStatus: (token: string) =>
    fetchAPI<any[]>('/admin/operations/jobs', { headers: authHeaders(token) }).then(r => r ?? []),
  adminGetFailedJobs: (token: string) =>
    fetchAPI<any[]>('/admin/operations/failed-jobs', { headers: authHeaders(token) }).then(r => r ?? []),

  // ── Seller Verification ──────────────────────────────────────────────────
  adminGetVerificationDocuments: (status: string, token: string) =>
    fetchAPI<any>(`/admin/verification/documents?status=${encodeURIComponent(status)}`, { headers: authHeaders(token) })
      .then(r => Array.isArray(r) ? r : (r?.data ?? [])),
  // Backend expects: { approved: bool, reason: string, notes: string }
  adminReviewVerificationDocument: (id: number, data: { approved: boolean; reason: string; notes: string }, token: string) =>
    fetchAPI<any>(`/admin/verification/documents/${id}/review`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),

  // ── Trust & Safety ───────────────────────────────────────────────────────
  adminGetAbuseReports: (status: string, token: string) =>
    fetchAPI<any[]>(`/admin/trust/reports?status=${encodeURIComponent(status)}`, { headers: authHeaders(token) }).then(r => r ?? []),
  adminReviewAbuseReport: (id: number, data: { status: string; action: string; notes: string }, token: string) =>
    fetchAPI<any>(`/admin/trust/reports/${id}/review`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminUpdateAbuseReportPriority: (id: number, priority: string, token: string) =>
    fetchAPI<any>(`/admin/trust/reports/${id}/priority`, { method: 'PUT', body: JSON.stringify({ priority }), headers: authHeaders(token) }),

  // ── Refund Management ───────────────────────────────────────────────────
  adminGetRefundRequests: (status: string, token: string) =>
    fetchAPI<any>(`/admin/refunds?status=${encodeURIComponent(status)}`, { headers: authHeaders(token) })
      .then(r => Array.isArray(r) ? r : (r?.data ?? [])),
  // Backend expects: { approved: bool, reason: string }
  adminReviewRefund: (id: number, data: { approved: boolean; reason: string }, token: string) =>
    fetchAPI<any>(`/admin/refunds/${id}/review`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  // Backend expects: { reference: string } — reference is required
  adminMarkRefundProcessed: (id: number, reference: string, token: string) =>
    fetchAPI<any>(`/admin/refunds/${id}/processed`, { method: 'PUT', body: JSON.stringify({ reference }), headers: authHeaders(token) }),

  // ── Account Appeals ──────────────────────────────────────────────────────
  adminGetAccountAppeals: (status: string, token: string) =>
    fetchAPI<any>(`/admin/appeals?status=${encodeURIComponent(status)}`, { headers: authHeaders(token) })
      .then(r => Array.isArray(r) ? r : (r?.data ?? [])),
  // Backend expects: { approved: bool, decision: string, action_taken?: string }
  adminReviewAppeal: (id: number, data: { approved: boolean; decision: string; action_taken?: string }, token: string) =>
    fetchAPI<any>(`/admin/appeals/${id}/review`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminUpdateAppealStatus: (id: number, data: { status: string }, token: string) =>
    fetchAPI<any>(`/admin/appeals/${id}/status`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),

  // ── Search Analytics ────────────────────────────────────────────────────────
  adminGetPopularSearches: (limit: number, token: string) =>
    fetchAPI<any[]>(`/admin/search/popular?limit=${limit}`, { headers: authHeaders(token) }).then(r => r ?? []),
  adminGetZeroResultQueries: (limit: number, token: string) =>
    fetchAPI<any[]>(`/admin/search/zero-results?limit=${limit}`, { headers: authHeaders(token) }).then(r => r ?? []),

  // ── Category Hierarchy ───────────────────────────────────────────────────
  adminGetCategoryHierarchy: (params: Record<string, any>, token: string) => {
    const q = '?' + new URLSearchParams(params).toString();
    return fetchAPI<any[]>(`/admin/category-hierarchy${q}`, { headers: authHeaders(token) }).then(r => r ?? []);
  },
  adminCreateCategoryHierarchy: (data: any, token: string) =>
    fetchAPI<any>('/admin/category-hierarchy', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminUpdateCategoryHierarchy: (id: number, data: any, token: string) =>
    fetchAPI<any>(`/admin/category-hierarchy/${id}`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminDeleteCategoryHierarchy: (id: number, token: string) =>
    fetchAPI<any>(`/admin/category-hierarchy/${id}`, { method: 'DELETE', headers: authHeaders(token) }),
  adminGetCategoryHierarchyStats: (id: number, token: string) =>
    fetchAPI<any>(`/admin/category-hierarchy/${id}/stats`, { headers: authHeaders(token) }),

  // ── Attribute Definitions ────────────────────────────────────────────────
  adminGetAttributes: (token: string) =>
    fetchAPI<any[]>('/admin/attributes', { headers: authHeaders(token) }).then(r => r ?? []),
  adminCreateAttribute: (data: any, token: string) =>
    fetchAPI<any>('/admin/attributes', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminAssignAttributeToCategory: (categoryId: number, data: any, token: string) =>
    fetchAPI<any>(`/admin/categories/${categoryId}/attributes`, { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminRemoveAttributeFromCategory: (categoryId: number, attributeId: number, token: string) =>
    fetchAPI<any>(`/admin/categories/${categoryId}/attributes/${attributeId}`, { method: 'DELETE', headers: authHeaders(token) }),

  // ── Location Hierarchy Management ──────────────────────────────────────────
  adminGetLocationHierarchy: (params?: { parent_id?: number; level?: number; search?: string }, token?: string) => {
    const cleanParams: Record<string, string> = {};
    if (params?.parent_id !== undefined) cleanParams.parent_id = String(params.parent_id);
    if (params?.level !== undefined) cleanParams.level = String(params.level);
    if (params?.search) cleanParams.search = params.search;
    const q = Object.keys(cleanParams).length > 0 ? '?' + new URLSearchParams(cleanParams).toString() : '';
    return fetchAPI<any>(`/admin/location-hierarchy${q}`, { headers: token ? authHeaders(token) : undefined })
      .then(r => Array.isArray(r) ? r : (r?.data ?? []));
  },
  adminGetLocationHierarchyItem: (id: number, token?: string) =>
    fetchAPI<any>(`/admin/location-hierarchy/${id}`, { headers: token ? authHeaders(token) : undefined })
      .then(r => r?.data ?? r),
  adminCreateLocationHierarchy: (data: any, token: string) =>
    fetchAPI<any>('/admin/location-hierarchy', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminUpdateLocationHierarchy: (id: number, data: any, token: string) =>
    fetchAPI<any>(`/admin/location-hierarchy/${id}`, { method: 'PUT', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminDeleteLocationHierarchy: (id: number, token: string) =>
    fetchAPI<any>(`/admin/location-hierarchy/${id}`, { method: 'DELETE', headers: authHeaders(token) }),
  adminImportLocationHierarchy: (data: { items: any[] }, token: string) =>
    fetchAPI<any>('/admin/location-hierarchy/import', { method: 'POST', body: JSON.stringify(data), headers: authHeaders(token) }),
  adminTestGeocode: (latitude: number, longitude: number, language = 'en', token?: string) =>
    fetchAPI<any>('/admin/location-hierarchy/test-geocode', {
      method: 'POST',
      body: JSON.stringify({ latitude, longitude, language }),
      headers: token ? authHeaders(token) : undefined,
    }).then(r => r?.data ?? r),
};

