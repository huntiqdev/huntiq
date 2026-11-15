'use client';
// Unified HUNTIQ Master Authentication Page
import { useState, useEffect, useRef, Suspense, FormEvent } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { api, ApiError } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import { useLocaleContext } from '@/context/LocaleContext';
import { OTPInput } from '@/components/forms/OTPInput';
import { validateEthioPhone } from '@/lib/validations/common';
import { translateError } from '@/lib/errors';
import { ErrorBanner } from '@/components/ui/ErrorBanner';
import type { ErrorType } from '@/lib/errors';
import {
  ShoppingBag, Phone, Lock, User, Building2, Mail,
  Eye, EyeOff, AlertCircle, CheckCircle2, ArrowLeft, Sparkles, Home,
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';
import PasswordStrength, { validatePasswordStrength } from '@/components/forms/PasswordStrength';

type Tab = 'login' | 'register';
type Step = 'form' | 'otp' | 'forgot' | 'forgot-otp' | 'forgot-reset';

const ACCENT_BTN = 'bg-orange-600 shadow-orange-600/25 hover:bg-orange-500';
const INPUT_CLS = 'w-full bg-zinc-100 dark:bg-white/5 border border-zinc-200 dark:border-white/10 rounded-xl px-4 py-3 text-base sm:text-sm text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-white/30 outline-none transition-all focus:ring-2 focus:ring-orange-500/30 focus:border-orange-500/60';
const LBL_CLS = 'text-[10px] font-black uppercase tracking-widest text-zinc-500 dark:text-white/40 block mb-1.5';

function AuthContent() {
  const { user, loading: authLoading, login, logout, activateCapability } = useAuth();
  const router = useRouter();
  const params = useSearchParams();
  const { t } = useLocaleContext();
  const u = (t as any).ui as Record<string, string>;
  const errors = (t as any).errors as Record<string, string>;

  type AccountType = 'buyer' | 'seller' | 'landlord';

  const asParam = params.get('as');
  const intent = params.get('intent');
  const condition = params.get('condition');
  const refCode = params.get('ref') ?? '';
  const rawRedirect = params.get('redirect') ?? '';
  const tabParam = params.get('tab');
  const typeParam = params.get('type');

  // Unified detection across all query styles (encoded or legacy unencoded)
  const isRental = asParam === 'landlord' || intent === 'rent' || typeParam === 'rent' || rawRedirect.includes('type=rent');
  const isSeller = asParam === 'seller' || intent === 'sell' || typeParam === 'sale' || rawRedirect.includes('type=sale');

  const initialAccountType: AccountType = isSeller
    ? 'seller'
    : isRental
    ? 'landlord'
    : 'buyer';

  // Defensive reconstruction of target redirect URL ensuring type=sale / type=rent is preserved
  let redirectUrl = rawRedirect;
  if (redirectUrl && !redirectUrl.includes('type=') && typeParam) {
    redirectUrl += (redirectUrl.includes('?') ? '&' : '?') + `type=${typeParam}`;
  }

  const [accountType, setAccountType] = useState<AccountType>(initialAccountType);
  const [sellerPref, setSellerPref] = useState<'sale' | 'rent' | 'both'>(
    isRental ? 'rent' : 'sale'
  );
  // Default to 'login' so returning users are never asked signup questions
  const [tab, setTab] = useState<Tab>(
    tabParam === 'register' || Boolean(refCode) ? 'register' : 'login'
  );
  const [step, setStep] = useState<Step>('form');
  const [error, setError] = useState('');
  const [errorType, setErrorType] = useState<ErrorType>('unknown');
  const [loading, setLoading] = useState(false);

  const setErr = (raw: string | null | undefined, code?: string | null) => {
    if (!raw && !code) { setError(''); return; }
    const { message, type } = translateError(raw, errors, code);
    setError(message);
    setErrorType(type);
  };
  const setErrFrom = (err: unknown) => {
    if (err instanceof ApiError) setErr(err.message, err.code);
    else if (err instanceof Error) setErr(err.message);
    else setErr(String(err));
  };
  const clearErr = () => { setError(''); setErrorType('unknown'); };

  const [loginPhone, setLoginPhone] = useState('');
  const [loginPassword, setLoginPassword] = useState('');
  const [showLoginPw, setShowLoginPw] = useState(false);
  const [loginPhoneErr, setLoginPhoneErr] = useState('');
  const [inactiveReason, setInactiveReason] = useState('');

  // Load saved phone from localStorage for returning registered users on this browser
  useEffect(() => {
    if (typeof window === 'undefined') return;
    try {
      const savedPhone = localStorage.getItem('huntiq_last_phone');
      if (savedPhone) {
        setLoginPhone(prev => prev || savedPhone);
        if (!tabParam && !refCode) {
          setTab('login');
        }
      }
    } catch {}
  }, [tabParam, refCode]);

  const [regPhone, setRegPhone] = useState('');
  const [regPhoneErr, setRegPhoneErr] = useState('');
  const [showRegPw, setShowRegPw] = useState(false);
  const [regForm, setRegForm] = useState({
    name: '',
    business_name: '',
    email: '',
    password: '',
    confirm: '',
  });

  const [forgotPhone, setForgotPhone] = useState('');
  const [forgotPhoneErr, setForgotPhoneErr] = useState('');
  const [forgotCode, setForgotCode] = useState('');
  const [forgotNewPw, setForgotNewPw] = useState('');
  const [forgotConfirmPw, setForgotConfirmPw] = useState('');
  const [showForgotPw, setShowForgotPw] = useState(false);

  useEffect(() => {
    clearErr();
    setStep('form');
  }, [tab]);

  // Contextual post-auth redirection
  const handlePostAuth = async (authUser: any) => {
    const userCaps = Array.isArray(authUser?.capabilities) ? authUser.capabilities : [];

    // If explicit redirect parameter provided, prioritize it
    if (redirectUrl) {
      if ((isSeller || redirectUrl.includes('type=sale')) && !userCaps.includes('seller')) {
        await activateCapability('seller').catch(() => {});
      } else if ((isRental || redirectUrl.includes('type=rent')) && !userCaps.includes('landlord')) {
        await activateCapability('landlord').catch(() => {});
      }
      router.replace(redirectUrl);
      return;
    }

    // If user arrived wanting to sell goods
    if (isSeller) {
      if (!userCaps.includes('seller')) {
        await activateCapability('seller').catch(() => {});
      }
      router.replace(`/seller/product/new?type=sale${condition ? `&condition=${condition}` : ''}`);
      return;
    }

    // If user arrived wanting to rent/lease property or equipment
    if (isRental) {
      if (!userCaps.includes('landlord')) {
        await activateCapability('landlord').catch(() => {});
      }
      router.replace('/seller/product/new?type=rent');
      return;
    }

    // Default route: if user has seller or landlord capability, go to seller dashboard
    if (userCaps.includes('seller') || userCaps.includes('landlord') || authUser?.role === 'seller') {
      router.replace('/seller/dashboard');
      return;
    }

    // Default regular buyer experience
    router.replace('/buyer/dashboard');
  };

  // If already authenticated on this browser, do NOT ask signup questions or render auth form — auto-route directly
  useEffect(() => {
    if (!authLoading && user) {
      handlePostAuth(user);
    }
  }, [user, authLoading]);

  // ── Forgot Password ────────────────────────────────────────────────────────
  const handleForgotSend = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    clearErr();
    const pr = validateEthioPhone(forgotPhone);
    if (!pr.valid) { setForgotPhoneErr(pr.error ?? u.invalid_phone); return; }
    setForgotPhoneErr('');
    setLoading(true);
    try {
      await api.forgotPassword(pr.normalized);
      setForgotPhone(pr.normalized);
      setStep('forgot-otp');
    } catch (err) {
      setErrFrom(err);
    } finally {
      setLoading(false);
    }
  };

  const handleForgotReset = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    clearErr();
    if (!forgotCode.trim()) { setErr(errors['otp_required'] ?? u.otp_required); return; }
    if (forgotNewPw.length < 8) { setErr(errors['validation_password_min'] ?? u.password_min_8); return; }
    if (forgotNewPw !== forgotConfirmPw) { setErr(errors['validation_passwords_mismatch'] ?? u.passwords_do_not_match); return; }
    setLoading(true);
    try {
      await api.resetPassword(forgotPhone, forgotCode.trim(), forgotNewPw);
      setLoginPhone(forgotPhone);
      setForgotPhone('');
      setForgotCode('');
      setForgotNewPw('');
      setForgotConfirmPw('');
      setStep('form');
      setTab('login');
      clearErr();
    } catch (err) {
      setErrFrom(err);
    } finally {
      setLoading(false);
    }
  };

  // ── Unified Login ──────────────────────────────────────────────────────────
  const handleLogin = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    clearErr();
    setInactiveReason('');
    const pr = validateEthioPhone(loginPhone);
    if (!pr.valid) { setLoginPhoneErr(pr.error ?? u.invalid_phone); return; }
    setLoginPhoneErr('');
    setLoading(true);
    try {
      const res = await api.unifiedLogin({ phone: pr.normalized, password: loginPassword });
      try {
        localStorage.setItem('huntiq_last_phone', pr.normalized);
      } catch {}
      await login(res.token, undefined, res.user);
      await handlePostAuth(res.user);
    } catch (err: any) {
      const msg: string = err?.message ?? '';
      if (msg.startsWith('account is inactive:')) {
        setInactiveReason(msg.replace('account is inactive:', '').trim());
        setError('');
      } else if (msg === 'account is inactive') {
        setInactiveReason('Your account has been deactivated. Please contact support.');
        setError('');
      } else if (msg.startsWith('otp_pending:')) {
        setRegPhone(pr.normalized);
      try {
        localStorage.setItem('huntiq_last_phone', pr.normalized);
      } catch {}
        setRegForm(prev => ({ ...prev, password: loginPassword }));
        api.sendOTP(pr.normalized, 'register', 'seller').catch(() => {});
        setStep('otp');
      } else {
        setErrFrom(err);
      }
    } finally {
      setLoading(false);
    }
  };

  // ── Unified Register ───────────────────────────────────────────────────────
  const handleRegister = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    clearErr();
    const pr = validateEthioPhone(regPhone);
    if (!pr.valid) { setRegPhoneErr(pr.error ?? u.invalid_phone); return; }
    setRegPhoneErr('');
    if (!regForm.name.trim()) { setErr(errors['validation_name_required'] ?? (u.full_name_label + ' is required')); return; }
    if (regForm.password !== regForm.confirm) { setErr(errors['validation_passwords_mismatch'] ?? u.passwords_do_not_match); return; }
    const pwStrength = validatePasswordStrength(regForm.password);
    if (!pwStrength.valid) { setErr(pwStrength.message); return; }
    if (regForm.password.length < 8) { setErr(errors['validation_password_min'] ?? u.password_min_8); return; }

    setLoading(true);
    try {
      const selectedRole = accountType === 'buyer' ? 'buyer' : 'seller';
      // Unified merchant identity: Single account serves both selling goods and renting properties
      const selectedSellerPref = accountType === 'buyer' 
        ? undefined 
        : (accountType === 'landlord' ? 'rent' : sellerPref);

      const regRes: any = await api.unifiedRegister({
        name: regForm.name.trim(),
        phone: pr.normalized,
        password: regForm.password,
        email: regForm.email.trim() || undefined,
        business_name: regForm.business_name.trim() || undefined,
        role: selectedRole,
        seller_preference: selectedSellerPref,
        referral_code: refCode || undefined,
      });

      setRegPhone(pr.normalized);
      try {
        localStorage.setItem('huntiq_last_phone', pr.normalized);
      } catch {}

      // If OTP disabled globally
      if (regRes?.otp_required === false) {
        if (regRes?.token && regRes?.user) {
          await login(regRes.token, undefined, regRes.user);
          await handlePostAuth(regRes.user);
          return;
        }
        const loginRes = await api.unifiedLogin({ phone: pr.normalized, password: regForm.password });
        await login(loginRes.token, undefined, loginRes.user);
        await handlePostAuth(loginRes.user);
        return;
      }

      // OTP required
      setStep('otp');
    } catch (err: any) {
      const msg: string = err?.message ?? '';
      if (msg.startsWith('otp_pending:')) {
        api.sendOTP(pr.normalized, 'register', 'seller').catch(() => {});
        setRegPhone(pr.normalized);
      try {
        localStorage.setItem('huntiq_last_phone', pr.normalized);
      } catch {}
        setStep('otp');
      } else if (msg.toLowerCase().includes('already registered') || msg.toLowerCase().includes('phone number already')) {
        setLoginPhone(pr.normalized);
        setTab('login');
        setErrFrom(err);
      } else {
        setErrFrom(err);
      }
    } finally {
      setLoading(false);
    }
  };

  // ── OTP Verify ─────────────────────────────────────────────────────────────
  const handleOTP = async (code: string) => {
    if (!regForm.password) {
      setLoginPhone(regPhone);
      setTab('login');
      setStep('form');
      setErrFrom(new Error('Your session expired. Please sign in with your password.'));
      return;
    }
    setLoading(true);
    clearErr();
    try {
      await api.verifyOTP(regPhone, code, 'register');
      const res = await api.unifiedLogin({ phone: regPhone, password: regForm.password });
      try {
        localStorage.setItem('huntiq_last_phone', regPhone);
      } catch {}
      await login(res.token, undefined, res.user);
      await handlePostAuth(res.user);
    } catch (err: any) {
      if (err?.message?.startsWith('otp_pending:')) {
        setErrFrom(new Error('OTP verification failed — please try again'));
      } else {
        setErrFrom(err);
      }
      setLoading(false);
    }
  };

  const setReg = (k: keyof typeof regForm) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setRegForm(p => ({ ...p, [k]: e.target.value }));

  const loginPhoneResult = loginPhone ? validateEthioPhone(loginPhone) : null;
  const regPhoneResult = regPhone ? validateEthioPhone(regPhone) : null;

  // ── OTP step UI ────────────────────────────────────────────────────────────
  if (step === 'otp') {
    return (
      <div className="space-y-4">
        <button onClick={() => { setStep('form'); clearErr(); }}
          className="flex items-center gap-1.5 text-xs text-zinc-400 dark:text-white/40 hover:text-zinc-700 dark:hover:text-white/70 transition-colors">
          <ArrowLeft size={12} /> {u.back}
        </button>
        <div>
          <h1 className="text-xl font-black text-zinc-900 dark:text-white">{u.verify_phone_title}</h1>
          <p className="text-sm text-zinc-500 dark:text-white/40 mt-1">{u.code_sent_to} {regPhone}</p>
        </div>
        <OTPInput
          recipient={regPhone}
          onVerified={handleOTP}
          onResend={() => api.sendOTP(regPhone, 'register', 'seller')}
          error={error}
          loading={loading}
        />
      </div>
    );
  }

  // ── Forgot password — phone ────────────────────────────────────────────────
  if (step === 'forgot') {
    return (
      <div className="space-y-4">
        <button onClick={() => { setStep('form'); clearErr(); setForgotPhone(''); setForgotPhoneErr(''); }}
          className="flex items-center gap-1.5 text-xs text-zinc-400 dark:text-white/40 hover:text-zinc-700 dark:hover:text-white/70 transition-colors">
          <ArrowLeft size={12} /> {u.back_to_sign_in}
        </button>
        <div>
          <h1 className="text-xl font-black text-zinc-900 dark:text-white">{u.reset_password_title}</h1>
          <p className="text-sm text-zinc-500 dark:text-white/40 mt-1">{u.reset_password_subtitle}</p>
        </div>
        <form onSubmit={handleForgotSend} className="space-y-4">
          <PhoneField
            value={forgotPhone}
            onChange={e => { setForgotPhone(e.target.value); setForgotPhoneErr(''); }}
            error={forgotPhoneErr}
            result={forgotPhone ? validateEthioPhone(forgotPhone) : null}
            phoneLabel={u.phone_number_required}
            phoneHint={u.phone_hint}
          />
          {error && <ErrorBanner message={error} type={errorType} />}
          <SubmitBtn loading={loading} label={u.send_reset_code} waitLabel={u.please_wait} />
        </form>
      </div>
    );
  }

  // ── Forgot password — reset ────────────────────────────────────────────────
  if (step === 'forgot-otp') {
    return (
      <div className="space-y-4">
        <button onClick={() => { setStep('forgot'); clearErr(); setForgotCode(''); }}
          className="flex items-center gap-1.5 text-xs text-zinc-400 dark:text-white/40 hover:text-zinc-700 dark:hover:text-white/70 transition-colors">
          <ArrowLeft size={12} /> {u.back}
        </button>
        <div>
          <h1 className="text-xl font-black text-zinc-900 dark:text-white">{u.enter_reset_code_title}</h1>
          <p className="text-sm text-zinc-500 dark:text-white/40 mt-1">{u.code_sent_to} {forgotPhone}</p>
        </div>
        <form onSubmit={handleForgotReset} className="space-y-4">
          <div>
            <label className={LBL_CLS}>{u.reset_code_label}</label>
            <input required value={forgotCode} onChange={e => setForgotCode(e.target.value)}
              placeholder={u.reset_code_placeholder} className={INPUT_CLS} maxLength={6} />
          </div>
          <div>
            <label className={LBL_CLS}><Lock size={9} className="inline mr-1" />{u.new_password_label}</label>
            <div className="relative">
              <Lock size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 pointer-events-none" />
              <input required type={showForgotPw ? 'text' : 'password'}
                value={forgotNewPw} onChange={e => setForgotNewPw(e.target.value)}
                placeholder="••••••••" className={INPUT_CLS + ' pl-10 pr-10'} />
              <button type="button" onClick={() => setShowForgotPw(v => !v)}
                className="absolute right-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 hover:text-zinc-700 dark:hover:text-white/60">
                {showForgotPw ? <EyeOff size={14} /> : <Eye size={14} />}
              </button>
            </div>
          </div>
          <div>
            <label className={LBL_CLS}><Lock size={9} className="inline mr-1" />{u.confirm_new_password_label}</label>
            <div className="relative">
              <Lock size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 pointer-events-none" />
              <input required type="password" value={forgotConfirmPw}
                onChange={e => setForgotConfirmPw(e.target.value)}
                placeholder="••••••••" className={INPUT_CLS + ' pl-10'} />
            </div>
          </div>
          {error && <ErrorBanner message={error} type={errorType} />}
          <SubmitBtn loading={loading} label={u.reset_password_btn} waitLabel={u.please_wait} />
        </form>
      </div>
    );
  }

  // If already authenticated on this browser or loading, show redirecting state — do NOT ask signup questions!
  if (authLoading || user) {
    return (
      <div className="space-y-6 text-center animate-fade-in py-12">
        <div className="w-12 h-12 border-4 border-orange-600 border-t-transparent rounded-full animate-spin mx-auto" />
        <div>
          <h2 className="text-base font-black text-zinc-900 dark:text-white">
            {user ? `${u.already_signed_in_title || 'Welcome back,'} ${user.name || user.phone}` : 'Loading...'}
          </h2>
          <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">
            {intent === 'sell'
              ? 'Opening listing creator for sale...'
              : intent === 'rent'
              ? 'Opening listing creator for rental...'
              : 'Opening your account...'}
          </p>
        </div>
      </div>
    );
  }

  // ── Main View (Unified Sign In / Create Account) ───────────────────────────
  return (
    <div className="space-y-5">
      {refCode && (
        <div className="flex items-start gap-2 bg-orange-500/10 border border-orange-500/20 rounded-xl px-3 py-2.5">
          <CheckCircle2 size={13} className="text-orange-400 shrink-0 mt-0.5" />
          <div>
            <p className="text-xs text-orange-400 font-bold">{u.invited_title}</p>
            <p className="text-[11px] text-orange-300/80 mt-0.5">
              {u.invite_code_label} <span className="font-mono font-bold">{refCode}</span> {u.invite_code_desc}
            </p>
          </div>
        </div>
      )}

      {/* Single Unified Tab Switcher */}
      <div className="flex bg-zinc-100 dark:bg-white/5 rounded-2xl p-1 border border-zinc-200 dark:border-white/10">
        {(['login', 'register'] as Tab[]).map(tb => (
          <button
            key={tb}
            type="button"
            onClick={() => setTab(tb)}
            className={`flex-1 py-2.5 rounded-xl text-sm font-black transition-all tap-bounce active:scale-95 ${
              tab === tb
                ? 'bg-orange-600 text-white shadow-md shadow-orange-600/20'
                : 'text-zinc-500 dark:text-white/40 hover:text-zinc-800 dark:hover:text-white/70'
            }`}
          >
            {tb === 'login' ? u.tab_sign_in : u.tab_create_account}
          </button>
        ))}
      </div>

      {inactiveReason && (
        <div className="bg-red-500/10 border border-red-500/30 rounded-2xl p-4 text-xs text-red-500 dark:text-red-400">
          <p className="font-bold flex items-center gap-1.5 mb-1">
            <AlertCircle size={14} /> Account Inactive
          </p>
          <p>{inactiveReason}</p>
        </div>
      )}

      {/* ── Unified Login Form ─────────────────────────────────────────────── */}
      {tab === 'login' && (
        <form onSubmit={handleLogin} className="space-y-4">
          <PhoneField
            value={loginPhone}
            onChange={e => { setLoginPhone(e.target.value); setLoginPhoneErr(''); }}
            error={loginPhoneErr}
            result={loginPhoneResult}
            phoneLabel={u.phone_number_required}
            phoneHint={u.phone_hint}
          />

          <div>
            <div className="flex items-center justify-between mb-1.5">
              <label className={LBL_CLS}><Lock size={9} className="inline mr-1" />{u.password_label}</label>
              <button
                type="button"
                onClick={() => { setStep('forgot'); setForgotPhone(loginPhone); clearErr(); }}
                className="text-[11px] text-orange-500 hover:text-orange-600 dark:hover:text-orange-400 transition-colors font-semibold"
              >
                {u.forgot_password}
              </button>
            </div>
            <div className="relative">
              <Lock size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 pointer-events-none" />
              <input
                required
                type={showLoginPw ? 'text' : 'password'}
                className={INPUT_CLS + ' pl-10 pr-10'}
                placeholder="••••••••"
                value={loginPassword}
                onChange={e => setLoginPassword(e.target.value)}
              />
              <button
                type="button"
                onClick={() => setShowLoginPw(v => !v)}
                className="absolute right-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 hover:text-zinc-700 dark:hover:text-white/60"
              >
                {showLoginPw ? <EyeOff size={14} /> : <Eye size={14} />}
              </button>
            </div>
          </div>

          {error && <ErrorBanner message={error} type={errorType} />}

          <SubmitBtn loading={loading} label={u.sign_in_arrow} waitLabel={u.please_wait} />

          <p className="text-center text-xs text-zinc-500 dark:text-white/40 pt-1">
            {u.dont_have_account}{' '}
            <button
              type="button"
              onClick={() => setTab('register')}
              className="text-orange-600 dark:text-orange-400 font-bold hover:underline"
            >
              {u.create_account_link}
            </button>
          </p>
        </form>
      )}

      {/* ── Unified Register Form ────────────────────────────────────────── */}
      {tab === 'register' && (
        <form onSubmit={handleRegister} className="space-y-3.5">
          {/* Account Type Selector (Disambiguating Buyer vs Seller vs Landlord) */}
          <div className="space-y-1.5">
            <div className="flex items-center justify-between">
              <label className={LBL_CLS}>{u.reg_select_type || 'Select Account Type'}</label>
              <span className="text-[10px] font-bold text-orange-600 dark:text-orange-400 uppercase tracking-wider">
                {u.reg_as_badge || 'Registering as'}: {accountType === 'seller' ? (u.reg_seller_title || 'Seller') : accountType === 'landlord' ? (u.reg_landlord_title || 'Landlord') : (u.reg_buyer_title || 'Buyer')}
              </span>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-2">
              {/* Buyer / Consumer Card */}
              <button
                type="button"
                onClick={() => setAccountType('buyer')}
                className={`p-2.5 rounded-2xl border text-left transition-all flex flex-col justify-between ${
                  accountType === 'buyer'
                    ? 'bg-orange-50/80 dark:bg-orange-950/40 border-orange-500 shadow-sm ring-1 ring-orange-500/30'
                    : 'bg-zinc-50 dark:bg-white/[0.02] border-zinc-200 dark:border-white/10 hover:border-zinc-300 dark:hover:border-white/20'
                }`}
              >
                <div className="flex items-center justify-between mb-1.5">
                  <div className={`w-6 h-6 rounded-lg flex items-center justify-center ${
                    accountType === 'buyer' ? 'bg-orange-600 text-white' : 'bg-zinc-200 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-300'
                  }`}>
                    <User size={13} />
                  </div>
                  {accountType === 'buyer' && <CheckCircle2 size={14} className="text-orange-600 dark:text-orange-400" />}
                </div>
                <div>
                  <p className="text-xs font-black text-zinc-900 dark:text-white leading-tight">
                    {u.reg_buyer_title || 'Buyer / Consumer'}
                  </p>
                  <p className="text-[10px] text-zinc-500 dark:text-zinc-400 mt-0.5 leading-snug">
                    {u.reg_buyer_desc || 'Shop, rent items, place orders, chat with sellers'}
                  </p>
                </div>
              </button>

              {/* Seller / Merchant Card */}
              <button
                type="button"
                onClick={() => setAccountType('seller')}
                className={`p-2.5 rounded-2xl border text-left transition-all flex flex-col justify-between ${
                  accountType === 'seller'
                    ? 'bg-orange-50/80 dark:bg-orange-950/40 border-orange-500 shadow-sm ring-1 ring-orange-500/30'
                    : 'bg-zinc-50 dark:bg-white/[0.02] border-zinc-200 dark:border-white/10 hover:border-zinc-300 dark:hover:border-white/20'
                }`}
              >
                <div className="flex items-center justify-between mb-1.5">
                  <div className={`w-6 h-6 rounded-lg flex items-center justify-center ${
                    accountType === 'seller' ? 'bg-orange-600 text-white' : 'bg-zinc-200 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-300'
                  }`}>
                    <ShoppingBag size={13} />
                  </div>
                  {accountType === 'seller' && <CheckCircle2 size={14} className="text-orange-600 dark:text-orange-400" />}
                </div>
                <div>
                  <p className="text-xs font-black text-zinc-900 dark:text-white leading-tight">
                    {u.reg_seller_title || 'Seller / Merchant'}
                  </p>
                  <p className="text-[10px] text-zinc-500 dark:text-zinc-400 mt-0.5 leading-snug">
                    {u.reg_seller_desc || 'Post items for sale and reach thousands across Ethiopia'}
                  </p>
                </div>
              </button>

              {/* Landlord / Rental Host Card */}
              <button
                type="button"
                onClick={() => setAccountType('landlord')}
                className={`p-2.5 rounded-2xl border text-left transition-all flex flex-col justify-between ${
                  accountType === 'landlord'
                    ? 'bg-teal-50/80 dark:bg-teal-950/40 border-teal-500 shadow-sm ring-1 ring-teal-500/30'
                    : 'bg-zinc-50 dark:bg-white/[0.02] border-zinc-200 dark:border-white/10 hover:border-zinc-300 dark:hover:border-white/20'
                }`}
              >
                <div className="flex items-center justify-between mb-1.5">
                  <div className={`w-6 h-6 rounded-lg flex items-center justify-center ${
                    accountType === 'landlord' ? 'bg-teal-600 text-white' : 'bg-zinc-200 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-300'
                  }`}>
                    <Home size={13} />
                  </div>
                  {accountType === 'landlord' && <CheckCircle2 size={14} className="text-teal-600 dark:text-teal-400" />}
                </div>
                <div>
                  <p className="text-xs font-black text-zinc-900 dark:text-white leading-tight">
                    {u.reg_landlord_title || 'Landlord / Host'}
                  </p>
                  <p className="text-[10px] text-zinc-500 dark:text-zinc-400 mt-0.5 leading-snug">
                    {u.reg_landlord_desc || 'Rent out tools, vehicles, equipment, or properties'}
                  </p>
                </div>
              </button>
            </div>

            {accountType !== 'buyer' && (
              <div className="p-3 bg-zinc-50 dark:bg-white/[0.03] border border-zinc-200 dark:border-white/10 rounded-2xl space-y-2 mt-2">
                <div className="flex items-center justify-between">
                  <label className={LBL_CLS}>Primary Posting Focus (Can list both anytime)</label>
                  <span className="text-[9px] font-bold text-orange-600 dark:text-orange-400 bg-orange-500/10 px-2 py-0.5 rounded-full">
                    Unified Account
                  </span>
                </div>
                <div className="grid grid-cols-3 gap-1.5">
                  <button
                    type="button"
                    onClick={() => { setSellerPref('sale'); setAccountType('seller'); }}
                    className={`py-2 px-2 rounded-xl text-center border text-xs font-bold transition-all ${
                      sellerPref === 'sale' && accountType === 'seller'
                        ? 'bg-orange-500 text-white border-orange-500 shadow-sm'
                        : 'bg-white dark:bg-zinc-800 border-zinc-200 dark:border-white/10 text-zinc-700 dark:text-zinc-300 hover:border-orange-300'
                    }`}
                  >
                    🛍️ Sell Goods
                  </button>
                  <button
                    type="button"
                    onClick={() => { setSellerPref('rent'); setAccountType('landlord'); }}
                    className={`py-2 px-2 rounded-xl text-center border text-xs font-bold transition-all ${
                      sellerPref === 'rent' || accountType === 'landlord'
                        ? 'bg-teal-600 text-white border-teal-600 shadow-sm'
                        : 'bg-white dark:bg-zinc-800 border-zinc-200 dark:border-white/10 text-zinc-700 dark:text-zinc-300 hover:border-teal-300'
                    }`}
                  >
                    🏠 Rent Properties
                  </button>
                  <button
                    type="button"
                    onClick={() => { setSellerPref('both'); setAccountType('seller'); }}
                    className={`py-2 px-2 rounded-xl text-center border text-xs font-bold transition-all ${
                      sellerPref === 'both' && accountType === 'seller'
                        ? 'bg-indigo-600 text-white border-indigo-600 shadow-sm'
                        : 'bg-white dark:bg-zinc-800 border-zinc-200 dark:border-white/10 text-zinc-700 dark:text-zinc-300 hover:border-indigo-300'
                    }`}
                  >
                    🔄 Both (Sell & Rent)
                  </button>
                </div>
              </div>
            )}
          </div>

          <PhoneField
            value={regPhone}
            onChange={e => { setRegPhone(e.target.value); setRegPhoneErr(''); }}
            error={regPhoneErr}
            result={regPhoneResult}
            phoneLabel={u.phone_number_required}
            phoneHint={u.phone_hint}
          />

          <div>
            <label className={LBL_CLS}><User size={9} className="inline mr-1" />{u.full_name_label}</label>
            <div className="relative">
              <User size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 pointer-events-none" />
              <input
                required
                className={INPUT_CLS + ' pl-10'}
                placeholder={u.full_name_placeholder}
                value={regForm.name}
                onChange={setReg('name')}
              />
            </div>
          </div>

          <div>
            <label className={LBL_CLS}>
              <Building2 size={9} className="inline mr-1" />{u.business_name_label}{' '}
              <span className="normal-case font-normal text-zinc-400 dark:text-white/30">
                ({accountType === 'buyer' ? u.optional_label : (accountType === 'seller' ? (u.shop_store_name || 'Shop / Store Name') : (u.agency_host_name || 'Agency / Host Name'))})
              </span>
            </label>
            <div className="relative">
              <Building2 size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 pointer-events-none" />
              <input
                className={INPUT_CLS + ' pl-10'}
                placeholder={accountType === 'seller' ? (u.shop_store_placeholder || 'Your Shop or Store Name (optional)') : accountType === 'landlord' ? (u.agency_host_placeholder || 'Your Property or Agency Name (optional)') : (u.business_name_placeholder || 'Business Name (optional)')}
                value={regForm.business_name}
                onChange={setReg('business_name')}
              />
            </div>
          </div>

          <div>
            <label className={LBL_CLS}>
              <Mail size={9} className="inline mr-1" />{u.email_label}{' '}
              <span className="normal-case font-normal text-zinc-400 dark:text-white/30">({u.optional_label})</span>
            </label>
            <div className="relative">
              <Mail size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 pointer-events-none" />
              <input
                type="email"
                className={INPUT_CLS + ' pl-10'}
                placeholder={u.email_placeholder}
                value={regForm.email}
                onChange={setReg('email')}
              />
            </div>
          </div>

          <div>
            <label className={LBL_CLS}><Lock size={9} className="inline mr-1" />{u.password_min_label}</label>
            <div className="relative">
              <Lock size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 pointer-events-none" />
              <input
                required
                type={showRegPw ? 'text' : 'password'}
                className={INPUT_CLS + ' pl-10 pr-10'}
                placeholder="••••••••"
                value={regForm.password}
                onChange={setReg('password')}
              />
              <button
                type="button"
                onClick={() => setShowRegPw(v => !v)}
                className="absolute right-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 hover:text-zinc-700 dark:hover:text-white/60"
              >
                {showRegPw ? <EyeOff size={14} /> : <Eye size={14} />}
              </button>
            </div>
            {regForm.password && (
              <PasswordStrength password={regForm.password} className="mt-2" />
            )}
          </div>

          <div>
            <label className={LBL_CLS}><Lock size={9} className="inline mr-1" />{u.confirm_password_label}</label>
            <div className="relative">
              <Lock size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 pointer-events-none" />
              <input
                required
                type="password"
                className={INPUT_CLS + ' pl-10'}
                placeholder="••••••••"
                value={regForm.confirm}
                onChange={setReg('confirm')}
              />
            </div>
          </div>

          {error && <ErrorBanner message={error} type={errorType} />}

          <SubmitBtn loading={loading} label={u.create_account_arrow} waitLabel={u.please_wait} />

          <p className="text-center text-xs text-zinc-500 dark:text-white/40 pt-1">
            {u.already_have_account}{' '}
            <button
              type="button"
              onClick={() => setTab('login')}
              className="text-orange-600 dark:text-orange-400 font-bold hover:underline"
            >
              {u.sign_in_link}
            </button>
          </p>
        </form>
      )}

      {/* Account flexibility footnote */}
      <div className="bg-zinc-100/70 dark:bg-white/[0.03] border border-zinc-200/70 dark:border-white/5 rounded-2xl p-3.5 mt-2 text-center">
        <p className="text-[11px] font-medium text-zinc-600 dark:text-white/50 leading-relaxed">
          {u.account_flexibility_note || '✨ One account for everything. Browse, buy, sell items, or rent property anytime without creating separate accounts.'}
        </p>
      </div>
    </div>
  );
}

function PhoneField({ value, onChange, error, result, phoneLabel, phoneHint }: {
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  error: string;
  result: ReturnType<typeof validateEthioPhone> | null;
  phoneLabel: string;
  phoneHint: string;
}) {
  return (
    <div>
      <label className={LBL_CLS}><Phone size={9} className="inline mr-1" />{phoneLabel}</label>
      <div className="relative">
        <Phone size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 pointer-events-none" />
        <input
          required
          type="tel"
          value={value}
          onChange={onChange}
          placeholder="09XXXXXXXX or 07XXXXXXXX"
          className={INPUT_CLS + ' pl-10'}
          autoFocus
        />
      </div>
      {result?.valid && (
        <p className="text-[10px] text-green-500 mt-1 flex items-center gap-1 font-semibold">
          <CheckCircle2 size={10} /> {result.normalized}
        </p>
      )}
      {error && (
        <p className="text-[10px] text-red-500 mt-1 flex items-center gap-1 font-semibold">
          <AlertCircle size={10} /> {error}
        </p>
      )}
      <p className="text-[10px] text-zinc-400 dark:text-white/30 mt-1">{phoneHint}</p>
    </div>
  );
}

function SubmitBtn({ loading, label, waitLabel }: {
  loading: boolean;
  label: string;
  waitLabel: string;
}) {
  return (
    <button
      type="submit"
      disabled={loading}
      className={`w-full ${ACCENT_BTN} text-white py-3.5 rounded-xl font-black text-sm transition-all disabled:opacity-50 shadow-lg tap-bounce active:scale-95`}
    >
      {loading ? waitLabel : label}
    </button>
  );
}

export default function AuthPage() {
  const { t } = useLocaleContext();
  const u = (t as any).ui as Record<string, string>;
  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 flex items-center justify-center p-4">
      <div className="w-full max-w-sm">
        <div className="mb-8">
          <HuntiqLogo
            size="lg"
            showText
            subtitle={u.sign_in_or_create ?? 'Sign in or create account'}
            href="/"
          />
        </div>
        <Suspense fallback={null}>
          <AuthContent />
        </Suspense>
        <p className="text-center text-xs text-zinc-400 dark:text-white/30 mt-6">
          <Link href="/" className="hover:text-zinc-700 dark:hover:text-white/70 transition-colors font-medium">
            {u.back_to_marketplace}
          </Link>
        </p>
      </div>
    </div>
  );
}
