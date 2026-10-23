import { fetchAPI, API_URL, ApiError } from './core';

type UserRole = 'seller' | 'buyer';

async function loginAs(role: UserRole, data: any): Promise<{ token: string }> {
  const res = await fetchAPI<any>(`/${role}/login`, { method: 'POST', body: JSON.stringify(data) });
  if ((res as any)?.otp_pending) {
    throw new ApiError('otp_pending: phone not verified — please complete OTP verification first', 'OTP_PENDING', 200);
  }
  return res as { token: string };
}

async function checkExists(role: UserRole, phone: string): Promise<boolean> {
  const res = await fetch(`${API_URL}/${role}/check`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ phone }),
  }).catch(() => null);
  if (!res) return role === 'seller'; // seller: fail-safe true; buyer: false
  return res.status !== 404;
}

export const authApi = {
  // ── Unified Authentication & Capabilities ────────────────────────────────
  unifiedRegister: (data: {
    name: string;
    phone: string;
    password: string;
    email?: string;
    business_name?: string;
    role?: string;
    seller_preference?: string;
    referral_code?: string;
  }) => fetchAPI<any>('/auth/register', { method: 'POST', body: JSON.stringify(data) }),

  unifiedLogin: async (data: { phone?: string; email?: string; password: string }) => {
    const res = await fetchAPI<any>('/auth/login', { method: 'POST', body: JSON.stringify(data) });
    if ((res as any)?.otp_pending) {
      throw new ApiError('otp_pending: phone not verified — please complete OTP verification first', 'OTP_PENDING', 200);
    }
    return res as { token: string; user: any; capabilities: string[] };
  },

  getMe: () => fetchAPI<any>('/auth/me'),

  getCapabilities: () => fetchAPI<{ capabilities: string[] }>('/auth/capabilities'),

  activateCapability: (capability: string, businessName?: string, preference?: string, metadata?: Record<string, any>) =>
    fetchAPI<any>('/auth/capabilities/activate', {
      method: 'POST',
      body: JSON.stringify({
        capability,
        business_name: businessName,
        preference,
        metadata,
      }),
    }),

  forgotPassword: (phone: string) =>
    fetchAPI<{ message: string }>('/auth/forgot-password', { method: 'POST', body: JSON.stringify({ phone }) }),
  resetPassword: (phone: string, code: string, new_password: string) =>
    fetchAPI<{ message: string }>('/auth/reset-password', { method: 'POST', body: JSON.stringify({ phone, code, new_password }) }),

  // Legacy endpoints maintained for backward compatibility
  sellerRegister: (data: any) =>
    fetchAPI<any>('/seller/register', { method: 'POST', body: JSON.stringify(data) }),
  buyerRegister: (data: any) =>
    fetchAPI<any>('/buyer/register', { method: 'POST', body: JSON.stringify(data) }),

  sellerLogin: (data: any) => loginAs('seller', data),
  buyerLogin:  (data: any) => loginAs('buyer',  data),

  checkSellerExists: (phone: string) => checkExists('seller', phone),
  checkBuyerExists:  (phone: string) => checkExists('buyer',  phone),

  sendOTP: (phone: string, purpose: string, userType: UserRole = 'seller') =>
    fetchAPI<any>('/auth/otp/send', { method: 'POST', body: JSON.stringify({ phone, purpose, user_type: userType }) }),
  verifyOTP: (phone: string, code: string, purpose: string) =>
    fetchAPI<any>('/auth/otp/verify', { method: 'POST', body: JSON.stringify({ phone, code, purpose }) }),

  sellerForgotPassword: (phone: string) =>
    fetchAPI<{ message: string }>('/seller/forgot-password', { method: 'POST', body: JSON.stringify({ phone }) }),
  sellerResetPassword: (phone: string, code: string, new_password: string) =>
    fetchAPI<{ message: string }>('/seller/reset-password', { method: 'POST', body: JSON.stringify({ phone, code, new_password }) }),
  buyerForgotPassword: (phone: string) =>
    fetchAPI<{ message: string }>('/buyer/forgot-password', { method: 'POST', body: JSON.stringify({ phone }) }),
  buyerResetPassword: (phone: string, code: string, new_password: string) =>
    fetchAPI<{ message: string }>('/buyer/reset-password', { method: 'POST', body: JSON.stringify({ phone, code, new_password }) }),

  adminLogin: (data: { email: string; password: string }) =>
    fetchAPI<{ email?: string; token?: string; is_super: boolean; message: string }>('/admin/login', { method: 'POST', body: JSON.stringify(data) }),
  adminLoginVerifyOTP: (data: { email: string; code: string }) =>
    fetchAPI<{ token: string }>('/admin/login/verify', { method: 'POST', body: JSON.stringify(data) }),
  adminForgotPassword: (email: string) =>
    fetchAPI<{ message: string }>('/admin/forgot-password', { method: 'POST', body: JSON.stringify({ email }) }),
  adminResetPassword: (email: string, code: string, new_password: string) =>
    fetchAPI<{ message: string }>('/admin/reset-password', { method: 'POST', body: JSON.stringify({ email, code, new_password }) }),

  // ── Enterprise Passkey / WebAuthn & Session Management ─────────────────────
  adminPasskeyLoginBegin: (data: { email: string; portal: string }) =>
    fetchAPI<{ challenge_id: string; options: any }>('/admin/auth/passkey/login/begin', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  adminPasskeyLoginVerify: (data: { challenge_id: string; response: any }) =>
    fetchAPI<{ token: string; session_token: string; role: string; redirect_url: string; user: any }>('/admin/auth/passkey/login/verify', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  adminPasskeyRegisterBegin: (token: string, friendlyName: string) =>
    fetchAPI<{ challenge_id: string; options: any }>('/admin/auth/passkey/register/begin', {
      method: 'POST',
      body: JSON.stringify({ friendly_name: friendlyName }),
      headers: { Authorization: `Bearer ${token}` },
    }),
  adminPasskeyRegisterVerify: (token: string, challengeId: string, friendlyName: string, response: any) =>
    fetchAPI<{ success: boolean; credential_id: number; message: string }>('/admin/auth/passkey/register/verify', {
      method: 'POST',
      body: JSON.stringify({ challenge_id: challengeId, friendly_name: friendlyName, response }),
      headers: { Authorization: `Bearer ${token}` },
    }),
  adminPasskeyList: (token: string) =>
    fetchAPI<{ authenticators: any[] }>('/admin/auth/passkey/list', {
      headers: { Authorization: `Bearer ${token}` },
    }),
  adminPasskeyDelete: (token: string, id: number) =>
    fetchAPI<{ success: boolean; message: string }>(`/admin/auth/passkey/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    }),
  adminReauth: (token: string, password: string, purpose: string) =>
    fetchAPI<{ success: boolean; reauth_token: string; expires_in_seconds: number }>('/admin/auth/reauth', {
      method: 'POST',
      body: JSON.stringify({ password, purpose }),
      headers: { Authorization: `Bearer ${token}` },
    }),
  adminSessionRevoke: (token: string, sessionId: string) =>
    fetchAPI<{ success: boolean; message: string }>('/admin/auth/session/revoke', {
      method: 'POST',
      body: JSON.stringify({ session_id: sessionId }),
      headers: { Authorization: `Bearer ${token}` },
    }),
};

