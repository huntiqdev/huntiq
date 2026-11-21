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
  Eye, EyeOff, AlertCircle, CheckCircle2, ArrowLeft, Sparkles, Home, Crown,
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';
import PasswordStrength, { validatePasswordStrength } from '@/components/forms/PasswordStrength';

type Tab = 'login' | 'register';
type Step = 'form' | 'otp' | 'forgot' | 'forgot-otp' | 'forgot-reset';

const ACCENT_BTN = 'bg-[#1769E0] shadow-[#1769E0]/25 hover:bg-[#1055BD]';
const INPUT_CLS = 'w-full bg-zinc-100 dark:bg-white/5 border border-zinc-200 dark:border-white/10 rounded-xl px-4 py-3 min-h-[44px] text-base sm:text-sm text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-white/30 outline-none transition-all focus:ring-2 focus:ring-[#1769E0]/30 focus:border-[#1769E0]';
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

  // Defensive reconstruction of target redirect URL ensuring type=sale / type=rent is preserved
  let redirectUrl = rawRedirect;
  if (redirectUrl && !redirectUrl.includes('type=') && typeParam) {
    redirectUrl += (redirectUrl.includes('?') ? '&' : '?') + `type=${typeParam}`;
  }

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
  const [selectedBuyerTier, setSelectedBuyerTier] = useState<'gold' | 'premium'>('gold');
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
  const [devOtp, setDevOtp] = useState<string>('');

  const fetchDevOtp = async (phone: string) => {
    if (!phone) return;
    try {
      const res = await fetch(`http://localhost:8080/api/dev/otp?phone=${encodeURIComponent(phone)}`);
      const data = await res.json();
      if (data.otp) {
        setDevOtp(data.otp);
      }
    } catch {}
  };

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
      const selectedRole = isSeller || isRental ? 'seller' : 'buyer';
      const selectedSellerPref = isRental ? 'rent' : (isSeller ? 'sale' : 'both');

      const regRes: any = await api.unifiedRegister({
        name: regForm.name.trim(),
        phone: pr.normalized,
        password: regForm.password,
        email: regForm.email.trim() || undefined,
        business_name: regForm.business_name.trim() || undefined,
        role: selectedRole,
        seller_preference: selectedSellerPref,
        referral_code: refCode || undefined,
        buyer_tier: selectedRole === 'buyer' ? selectedBuyerTier : undefined,
        payment_method: selectedRole === 'buyer' ? 'telebirr' : undefined,
        payment_ref: selectedRole === 'buyer' ? `TB-REG-${Date.now().toString().slice(-6)}` : undefined,
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

        {/* Development Environment OTP Helper */}
        <div className="p-3 bg-amber-500/10 border border-amber-500/30 rounded-xl text-xs space-y-1">
          <div className="flex items-center justify-between">
            <span className="font-bold text-amber-500">🛠️ Dev Environment Helper:</span>
            <button
              type="button"
              onClick={() => fetchDevOtp(regPhone)}
              className="text-[10px] px-2 py-0.5 bg-amber-500/20 text-amber-400 rounded hover:bg-amber-500/30 font-medium transition-colors"
            >
              Fetch Dev OTP
            </button>
          </div>
          <p className="text-[11px] text-zinc-400">
            In development mode, OTP codes are logged to backend terminal or available via endpoint:
          </p>
          <p className="font-mono text-zinc-300">
            GET <a href={`http://localhost:8080/api/dev/otp?phone=${encodeURIComponent(regPhone)}`} target="_blank" rel="noreferrer" className="text-[#1769E0] underline">/api/dev/otp?phone={regPhone}</a>
          </p>
          {devOtp && (
            <div className="mt-1 p-1.5 bg-black/40 rounded border border-amber-500/40 text-center font-mono font-bold text-amber-300 text-sm">
              Active OTP: {devOtp}
            </div>
          )}
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

        {/* Development Environment OTP Helper */}
        <div className="p-3 bg-amber-500/10 border border-amber-500/30 rounded-xl text-xs space-y-1">
          <div className="flex items-center justify-between">
            <span className="font-bold text-amber-500">🛠️ Dev Environment Helper:</span>
            <button
              type="button"
              onClick={() => fetchDevOtp(forgotPhone)}
              className="text-[10px] px-2 py-0.5 bg-amber-500/20 text-amber-400 rounded hover:bg-amber-500/30 font-medium transition-colors"
            >
              Fetch Dev OTP
            </button>
          </div>
          <p className="text-[11px] text-zinc-400">
            In development mode, reset OTP codes are available via:
          </p>
          <p className="font-mono text-zinc-300">
            GET <a href={`http://localhost:8080/api/dev/otp?phone=${encodeURIComponent(forgotPhone)}`} target="_blank" rel="noreferrer" className="text-[#1769E0] underline">/api/dev/otp?phone={forgotPhone}</a>
          </p>
          {devOtp && (
            <div className="mt-1 p-1.5 bg-black/40 rounded border border-amber-500/40 text-center font-mono font-bold text-amber-300 text-sm">
              Active OTP: {devOtp}
            </div>
          )}
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
        <div className="w-12 h-12 border-4 border-[#1769E0] border-t-transparent rounded-full animate-spin mx-auto" />
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
        <div className="flex items-start gap-2 bg-[#1769E0]/10 border border-[#1769E0]/20 rounded-xl px-3 py-2.5">
          <CheckCircle2 size={13} className="text-[#1769E0] dark:text-[#38A5FF] shrink-0 mt-0.5" />
          <div>
            <p className="text-xs text-[#1769E0] dark:text-[#38A5FF] font-bold">{u.invited_title}</p>
            <p className="text-[11px] text-zinc-600 dark:text-zinc-300 mt-0.5">
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
                ? 'bg-[#1769E0] text-white shadow-md shadow-[#1769E0]/20'
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
                className="text-[11px] text-[#1769E0] hover:text-[#1055BD] dark:text-[#38A5FF] dark:hover:text-[#60B5FF] transition-colors font-semibold"
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
              className="text-[#1769E0] dark:text-[#38A5FF] font-bold hover:underline"
            >
              {u.create_account_link}
            </button>
          </p>
        </form>
      )}

      {/* ── Unified Register Form ────────────────────────────────────────── */}
      {tab === 'register' && (
        <form onSubmit={handleRegister} className="space-y-3.5">
          {/* Master Universal Account Banner */}
          <div className="flex items-center gap-3 p-3 bg-zinc-50 dark:bg-white/[0.03] border border-zinc-200/80 dark:border-white/10 rounded-2xl text-xs">
            <div className="w-9 h-9 rounded-xl bg-[#1769E0]/10 flex items-center justify-center shrink-0 text-[#1769E0]">
              <Sparkles size={18} />
            </div>
            <div>
              <p className="font-black text-zinc-900 dark:text-white leading-tight">
                Master Universal Account
              </p>
              <p className="text-[11px] text-zinc-500 dark:text-zinc-400 mt-0.5 leading-snug">
                One master account for everything: Buy, Rent, Sell, Host, or Announce across Ethiopia.
              </p>
            </div>
          </div>

          {/* Required Buyer Privilege Tier (Normal Buyer Registration Disabled) */}
          {!isSeller && !isRental && (
            <div className="space-y-2 p-3.5 bg-amber-500/10 border border-amber-500/30 rounded-2xl">
              <div className="flex items-center justify-between">
                <span className="text-[11px] font-black uppercase text-amber-500 flex items-center gap-1.5">
                  <Crown size={13} /> Required Buyer Privilege Tier
                </span>
                <span className="text-[10px] text-zinc-400">Normal registration disabled</span>
              </div>
              <div className="grid grid-cols-2 gap-2 text-xs">
                <button
                  type="button"
                  onClick={() => setSelectedBuyerTier('gold')}
                  className={`p-2.5 rounded-xl border text-left transition-all ${
                    selectedBuyerTier === 'gold'
                      ? 'bg-amber-500/20 border-amber-500 text-white font-bold ring-1 ring-amber-500/50'
                      : 'bg-white/5 border-white/10 text-zinc-400'
                  }`}
                >
                  <p className="font-black text-amber-400">Gold Buyer</p>
                  <p className="text-[10px] text-zinc-300">299 ETB / mo</p>
                  <p className="text-[9px] text-zinc-400 mt-0.5">Direct seller contact & broker discount</p>
                </button>
                <button
                  type="button"
                  onClick={() => setSelectedBuyerTier('premium')}
                  className={`p-2.5 rounded-xl border text-left transition-all ${
                    selectedBuyerTier === 'premium'
                      ? 'bg-purple-500/20 border-purple-500 text-white font-bold ring-1 ring-purple-500/50'
                      : 'bg-white/5 border-white/10 text-zinc-400'
                  }`}
                >
                  <p className="font-black text-purple-400">Premium VIP</p>
                  <p className="text-[10px] text-zinc-300">599 ETB / mo</p>
                  <p className="text-[9px] text-zinc-400 mt-0.5">Off-market demand & 1-on-1 advisor</p>
                </button>
              </div>
            </div>
          )}

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
              <Building2 size={9} className="inline mr-1" />{u.business_name_label || 'Business / Store / Brand Name'}{' '}
              <span className="normal-case font-normal text-zinc-400 dark:text-white/30">
                ({u.optional_label || 'optional'})
              </span>
            </label>
            <div className="relative">
              <Building2 size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 pointer-events-none" />
              <input
                className={INPUT_CLS + ' pl-10'}
                placeholder={u.business_name_placeholder || 'Business, shop, or agency name (optional)'}
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
              className="text-[#1769E0] dark:text-[#38A5FF] font-bold hover:underline"
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
