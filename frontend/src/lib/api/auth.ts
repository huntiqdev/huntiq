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
};
