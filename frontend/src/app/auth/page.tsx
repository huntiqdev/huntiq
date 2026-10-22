'use client';
// auth page
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
import { usePaymentSettings } from '@/hooks/usePaymentSettings';
import type { ErrorType } from '@/lib/errors';
import {
  ShoppingBag, Phone, Lock, User, Building2, Mail,
  Eye, EyeOff, AlertCircle, Star, CheckCircle2, ArrowLeft, Zap,
} from 'lucide-react';
import PasswordStrength, { validatePasswordStrength } from '@/components/forms/PasswordStrength';

type Role = 'seller' | 'buyer';
type Tab  = 'login' | 'register';
type Step = 'form' | 'otp' | 'forgot' | 'forgot-otp' | 'forgot-reset';

const ROLE_ICONS: Record<Role, React.ReactNode> = {
  seller: <ShoppingBag size={14} />,
  buyer:  <Star size={14} />,
};
const ROLE_ACCENT: Record<Role, string> = {
  seller: 'bg-orange-600 shadow-orange-600/20',
  buyer:  'bg-yellow-500 shadow-yellow-500/20',
};
const ROLE_RING: Record<Role, string> = {
  seller: 'focus:ring-orange-500/30 focus:border-orange-500/60',
  buyer:  'focus:ring-yellow-500/30 focus:border-yellow-500/60',
};

function AuthContent() {
  const { login } = useAuth();
  const router = useRouter();
  const params = useSearchParams();
  const { t } = useLocaleContext();
  const u = (t as any).ui as Record<string, string>;
  const errors = (t as any).errors as Record<string, string>;

  const asParam   = params.get('as') as Role | null;
  const intent    = params.get('intent');
  const condition = params.get('condition');
  const refCode   = params.get('ref') ?? '';
  const preferenceParam = params.get('preference') as 'sale' | 'rent' | 'both' | null;

  const initialRole: Role = asParam === 'buyer' ? 'buyer' : 'seller';
  const [role, setRole] = useState<Role>(initialRole);
  const [tab,  setTab]  = useState<Tab>(refCode ? 'register' : 'login');
  const [step, setStep] = useState<Step>('form');
  const [error,     setError]     = useState('');
  const [errorType, setErrorType] = useState<ErrorType>('unknown');
  const [loading,   setLoading]   = useState(false);

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

  const autoDetectRef = useRef<string>('');

  const [loginPhone,    setLoginPhone]    = useState('');
  const [loginPassword, setLoginPassword] = useState('');
  const [showLoginPw,   setShowLoginPw]   = useState(false);
  const [loginPhoneErr, setLoginPhoneErr] = useState('');
  const [inactiveReason, setInactiveReason] = useState('');

  const [regPhone,    setRegPhone]    = useState('');
  const [regPhoneErr, setRegPhoneErr] = useState('');
  const [showRegPw,   setShowRegPw]   = useState(false);
  // regForm.password is also used after OTP verify to auto-login
  const [regForm, setRegForm] = useState({
    name: '', business_name: '', email: '', password: '', confirm: '', 
    seller_preference: preferenceParam || 'sale',
  });

  const [forgotPhone,     setForgotPhone]     = useState('');
  const [forgotPhoneErr,  setForgotPhoneErr]  = useState('');
  const [forgotCode,      setForgotCode]      = useState('');
  const [forgotNewPw,     setForgotNewPw]     = useState('');
  const [forgotConfirmPw, setForgotConfirmPw] = useState('');
  const [showForgotPw,    setShowForgotPw]    = useState(false);

  const { settings: paySettings } = usePaymentSettings();
  const tierPrices = {
    gold:    String(paySettings.gold_upgrade_price),
    premium: String(paySettings.premium_upgrade_price),
  };

  // Don't reset form state while user is on the OTP step — password would be wiped
  useEffect(() => {
    if (step === 'otp') return;
    clearErr(); setStep('form');
    setLoginPhone(''); setLoginPassword(''); setLoginPhoneErr(''); setRegPhoneErr('');
    setRegPhone('');
    setRegForm({ name: '', business_name: '', email: '', password: '', confirm: '', seller_preference: 'sale' });
    autoDetectRef.current = '';
    if (refCode) setTab('register');
  }, [role]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    clearErr(); setStep('form');
    autoDetectRef.current = '';
  }, [tab]);

  // Auto-detect role from phone number on login tab (debounced 400ms)
  useEffect(() => {
    if (tab !== 'login' || refCode) return;
    const pr = validateEthioPhone(loginPhone);
    if (!pr.valid) return;
    const normalized = pr.normalized;
    if (autoDetectRef.current === normalized) return;
    const timer = setTimeout(() => {
      autoDetectRef.current = normalized;
      let cancelled = false;
      Promise.all([
        api.checkSellerExists(normalized).catch(() => null),
        api.checkBuyerExists(normalized).catch(() => null),
      ]).then(([isSeller, isBuyer]) => {
        if (cancelled) return;
        if (isBuyer && !isSeller) setRole('buyer');
        else if (isSeller && !isBuyer) setRole('seller');
      });
      return () => { cancelled = true; };
    }, 400);
    return () => clearTimeout(timer);
  }, [loginPhone, tab, refCode]); // eslint-disable-line react-hooks/exhaustive-deps

  const inputCls = `w-full bg-zinc-100 dark:bg-white/5 border border-zinc-200 dark:border-white/10 rounded-xl px-4 py-3 text-sm text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-white/30 outline-none transition-all ${ROLE_RING[role]}`;
  const lbl = 'text-[10px] font-black uppercase tracking-widest text-zinc-500 dark:text-white/40 block mb-1.5';

  const afterAuth = (userRole: Role, isNewRegistration = false) => {
    // New registrations go to demand selection page
    if (isNewRegistration) {
      router.push('/demand');
      return;
    }
    
    // Existing login flow
    if (userRole === 'buyer') { router.push('/buyer/dashboard'); return; }
    if (intent === 'sell') {
      router.push(`/seller/products?new=1${condition ? `&condition=${condition}` : ''}`);
    } else {
      router.push('/seller/dashboard');
    }
  };

  // ── Forgot password ────────────────────────────────────────────────────────
  const handleForgotSend = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault(); clearErr();
    const pr = validateEthioPhone(forgotPhone);
    if (!pr.valid) { setForgotPhoneErr(pr.error ?? u.invalid_phone); return; }
    setForgotPhoneErr(''); setLoading(true);
    try {
      if (role === 'buyer') await api.buyerForgotPassword(pr.normalized);
      else await api.sellerForgotPassword(pr.normalized);
      setForgotPhone(pr.normalized);
      setStep('forgot-otp');
    } catch (err) { setErrFrom(err); }
    finally { setLoading(false); }
  };

  // After reset succeeds: go back to login form (do NOT auto-login — user should sign in fresh)
  const handleForgotReset = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault(); clearErr();
    if (!forgotCode.trim()) { setErr(errors['otp_required'] ?? u.otp_required); return; }
    if (forgotNewPw.length < 8) { setErr(errors['validation_password_min'] ?? u.password_min_8); return; }
    if (forgotNewPw !== forgotConfirmPw) { setErr(errors['validation_passwords_mismatch'] ?? u.passwords_do_not_match); return; }
    setLoading(true);
    try {
      if (role === 'buyer') await api.buyerResetPassword(forgotPhone, forgotCode.trim(), forgotNewPw);
      else await api.sellerResetPassword(forgotPhone, forgotCode.trim(), forgotNewPw);
      // Reset succeeded — pre-fill login phone and go to login tab
      setLoginPhone(forgotPhone);
      setForgotPhone(''); setForgotCode(''); setForgotNewPw(''); setForgotConfirmPw('');
      setStep('form');
      setTab('login');
      clearErr();
    } catch (err) { setErrFrom(err); }
    finally { setLoading(false); }
  };

  // ── Login ──────────────────────────────────────────────────────────────────
  const handleLogin = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault(); clearErr(); setInactiveReason('');
    const pr = validateEthioPhone(loginPhone);
    if (!pr.valid) { setLoginPhoneErr(pr.error ?? u.invalid_phone); return; }
    setLoginPhoneErr(''); setLoading(true);
    try {
      const res = role === 'buyer'
        ? await api.buyerLogin({ phone: pr.normalized, password: loginPassword })
        : await api.sellerLogin({ phone: pr.normalized, password: loginPassword });
      await login(res.token);
      afterAuth(role);
    } catch (err: any) {
      const msg: string = err?.message ?? '';
      if (msg.startsWith('account is inactive:')) {
        setInactiveReason(msg.replace('account is inactive:', '').trim());
        setError('');
      } else if (msg === 'account is inactive') {
        setInactiveReason('Your account has been deactivated. Please contact support.');
        setError('');
      } else if (msg.startsWith('otp_pending:')) {
        // User registered but never completed OTP — send them through registration OTP flow.
        // OTP purpose is always 'register' (the only OTP type for sellers/buyers).
        setRegPhone(pr.normalized);
        setRegForm(prev => ({ ...prev, password: loginPassword }));
        api.sendOTP(pr.normalized, 'register', role as 'seller' | 'buyer').catch(() => {});
        setStep('otp');
      } else {
        setErrFrom(err);
      }
    } finally { setLoading(false); }
  };

  // ── Register ───────────────────────────────────────────────────────────────
  const handleRegister = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault(); clearErr();
    const pr = validateEthioPhone(regPhone);
    if (!pr.valid) { setRegPhoneErr(pr.error ?? u.invalid_phone); return; }
    setRegPhoneErr('');
    if (!regForm.name.trim()) { setErr(errors['validation_name_required'] ?? (u.full_name_label + ' is required')); return; }
    if (regForm.password !== regForm.confirm) { setErr(errors['validation_passwords_mismatch'] ?? u.passwords_do_not_match); return; }
    const pwStrength = validatePasswordStrength(regForm.password);
    if (!pwStrength.valid) { setErr(pwStrength.message); return; }
    setLoading(true);
    if (regForm.password.length < 8) { setErr(errors['validation_password_min'] ?? u.password_min_8); return; }
    try {
      const regRes: any = role === 'buyer'
        ? await api.buyerRegister({ ...regForm, phone: pr.normalized, role: 'buyer', referral_code: refCode })
        : await api.sellerRegister({ ...regForm, phone: pr.normalized, role: 'seller', referral_code: refCode, seller_preference: regForm.seller_preference });
      setRegPhone(pr.normalized);

      // otp_required=false means OTP is disabled globally — log in immediately
      if (regRes?.otp_required === false) {
        const loginRes = role === 'buyer'
          ? await api.buyerLogin({ phone: pr.normalized, password: regForm.password })
          : await api.sellerLogin({ phone: pr.normalized, password: regForm.password });
        await login(loginRes.token);
        afterAuth(role);
        return;
      }

      // OTP is required (default) — backend already auto-sent OTP on register.
      // Always go to OTP step; never skip it when otp_required=true.
      setStep('otp');
    } catch (err: any) {
      const msg: string = err?.message ?? '';
      if (msg.startsWith('otp_pending:')) {
        // Phone already registered but unverified — resend OTP and go to OTP step
        api.sendOTP(pr.normalized, 'register', role as 'seller' | 'buyer').catch(() => {});
        setRegPhone(pr.normalized);
        setStep('otp');
      } else if (msg.toLowerCase().includes('already registered') || msg.toLowerCase().includes('phone number already')) {
        setLoginPhone(pr.normalized);
        setTab('login');
        setErrFrom(err);
      } else {
        setErrFrom(err);
      }
    } finally { setLoading(false); }
  };

  // ── OTP verify → auto-login ────────────────────────────────────────────────
  // OTP is only used at registration. After verify, immediately log in.
  const handleOTP = async (code: string) => {
    if (!regForm.password) {
      // Password was lost (e.g. page refresh) — redirect to login with phone pre-filled
      // so the user can sign in normally (their account already exists and is verified).
      setLoginPhone(regPhone);
      setTab('login');
      setStep('form');
      setErrFrom(new Error('Your session expired — your account is registered, please sign in with your password'));
      return;
    }
    setLoading(true); clearErr();
    try {
      // Step 1: verify OTP — marks phone as verified in DB
      await api.verifyOTP(regPhone, code, 'register');
      // Step 2: login — now that phone is verified, this will succeed
      const res = role === 'buyer'
        ? await api.buyerLogin({ phone: regPhone, password: regForm.password })
        : await api.sellerLogin({ phone: regPhone, password: regForm.password });
      await login(res.token);
      afterAuth(role, true); // true = new registration
    } catch (err: any) {
      // If login still returns otp_pending after verify, the verify failed silently
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
  const regPhoneResult   = regPhone   ? validateEthioPhone(regPhone)   : null;

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
        {/* OTP is only for registration — resend always uses purpose 'register' */}
        <OTPInput
          recipient={regPhone}
          onVerified={handleOTP}
          onResend={() => api.sendOTP(regPhone, 'register', role as 'seller' | 'buyer')}
          error={error}
          loading={loading}
        />
      </div>
    );
  }

  // ── Forgot password — enter phone ──────────────────────────────────────────
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
            inputCls={inputCls} lbl={lbl}
            phoneLabel={u.phone_number_required} phoneHint={u.phone_hint}
          />
          {error && <ErrorBanner message={error} type={errorType} />}
          <SubmitBtn loading={loading} label={u.send_reset_code} accent={ROLE_ACCENT[role]} waitLabel={u.please_wait} />
        </form>
      </div>
    );
  }

  // ── Forgot password — enter OTP + new password ─────────────────────────────
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
            <label className={lbl}>{u.reset_code_label}</label>
            <input required value={forgotCode} onChange={e => setForgotCode(e.target.value)}
              placeholder={u.reset_code_placeholder} className={inputCls} maxLength={6} />
          </div>
          <div>
            <label className={lbl}><Lock size={9} className="inline mr-1" />{u.new_password_label}</label>
            <div className="relative">
              <Lock size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none" />
              <input required type={showForgotPw ? 'text' : 'password'}
                value={forgotNewPw} onChange={e => setForgotNewPw(e.target.value)}
                placeholder="••••••••" className={inputCls + ' pl-10 pr-10'} />
              <button type="button" onClick={() => setShowForgotPw(v => !v)}
                className="absolute right-3.5 top-1/2 -translate-y-1/2 text-white/30 hover:text-white/60">
                {showForgotPw ? <EyeOff size={14} /> : <Eye size={14} />}
              </button>
            </div>
          </div>
          <div>
            <label className={lbl}><Lock size={9} className="inline mr-1" />{u.confirm_new_password_label}</label>
            <div className="relative">
              <Lock size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none" />
              <input required type="password" value={forgotConfirmPw}
                onChange={e => setForgotConfirmPw(e.target.value)}
                placeholder="••••••••" className={inputCls + ' pl-10'} />
            </div>
          </div>
          {error && <ErrorBanner message={error} type={errorType} />}
          <SubmitBtn loading={loading} label={u.reset_password_btn} accent={ROLE_ACCENT[role]} waitLabel={u.please_wait} />
        </form>
      </div>
    );
  }

  // ── Main form (login / register tabs) ─────────────────────────────────────
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

      {/* Role switcher */}
      <div className="flex bg-zinc-100 dark:bg-white/5 rounded-xl p-1 gap-1">
        {(['seller', 'buyer'] as Role[]).map(r => (
          <button key={r} type="button" onClick={() => setRole(r)}
            className={`flex-1 flex items-center justify-center gap-1.5 py-2 rounded-lg text-xs font-black transition-all ${
              role === r ? `${ROLE_ACCENT[r]} text-white shadow-md` : 'text-zinc-500 dark:text-white/40 hover:text-zinc-700 dark:hover:text-white/70'
            }`}>
            {ROLE_ICONS[r]}
            {r === 'seller' ? u.role_seller : u.role_buyer}
          </button>
        ))}
      </div>

      {/* Tab switcher */}
      <div className="flex bg-zinc-100 dark:bg-white/5 rounded-xl p-1">
        {(['login', 'register'] as Tab[]).map(tb => (
          <button key={tb} type="button" onClick={() => setTab(tb)}
            className={`flex-1 py-2 rounded-lg text-sm font-bold transition-all ${
              tab === tb ? `${ROLE_ACCENT[role]} text-white shadow-md` : 'text-zinc-500 dark:text-white/40 hover:text-zinc-700 dark:hover:text-white/70'
            }`}>
            {tb === 'login' ? u.tab_sign_in : u.tab_create_account}
          </button>
        ))}
      </div>

      {/* Login form */}
      {tab === 'login' && (
        <form onSubmit={handleLogin} className="space-y-4">
          {role === 'buyer' && (
            <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-xl p-3 mb-4">
              <p className="text-xs text-yellow-400 font-bold">Premium & Gold Login</p>
              <p className="text-[10px] text-yellow-400/80 mt-1">
                Normal buyers do not need an account. Login here to access your Gold or Premium features.
              </p>
            </div>
          )}
          <PhoneField value={loginPhone} onChange={e => { setLoginPhone(e.target.value); setLoginPhoneErr(''); }}
            error={loginPhoneErr} result={loginPhoneResult} inputCls={inputCls} lbl={lbl}
            phoneLabel={u.phone_number_required} phoneHint={u.phone_hint} />
          <div>
            <label className={lbl}><Lock size={9} className="inline mr-1" />{u.password_label}</label>
            <div className="relative">
              <Lock size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none" />
              <input required type={showLoginPw ? 'text' : 'password'}
                value={loginPassword} onChange={e => setLoginPassword(e.target.value)}
                placeholder="••••••••" className={inputCls + ' pl-10 pr-10'} />
              <button type="button" onClick={() => setShowLoginPw(v => !v)}
                className="absolute right-3.5 top-1/2 -translate-y-1/2 text-white/30 hover:text-white/60">
                {showLoginPw ? <EyeOff size={14} /> : <Eye size={14} />}
              </button>
            </div>
          </div>
          {inactiveReason && (
            <div className="flex items-start gap-3 bg-red-500/10 border border-red-500/30 rounded-xl px-4 py-3">
              <AlertCircle size={16} className="text-red-400 shrink-0 mt-0.5" />
              <div>
                <p className="text-sm font-black text-red-400">Account Deactivated</p>
                <p className="text-xs text-red-300/80 mt-0.5">{inactiveReason}</p>
                <p className="text-[10px] text-red-400/60 mt-1">Contact support if you believe this is a mistake.</p>
              </div>
            </div>
          )}
          {error && <ErrorBanner message={error} type={errorType} />}
          <SubmitBtn loading={loading} label={u.sign_in_arrow} accent={ROLE_ACCENT[role]} waitLabel={u.please_wait} />
          <div className="flex items-center justify-between text-xs text-zinc-400 dark:text-white/30">
            <p>
              {u.no_account_short}{' '}
              <button type="button" onClick={() => setTab('register')}
                className="text-orange-400 font-bold hover:text-orange-300">
                {refCode ? u.create_one_to_use_invite : u.create_one_free}
              </button>
            </p>
            <button type="button"
              onClick={() => { setStep('forgot'); setForgotPhone(loginPhone); clearErr(); }}
              className="text-zinc-400 dark:text-white/40 hover:text-zinc-700 dark:hover:text-white/70 transition-colors">
              {u.forgot_password}
            </button>
          </div>
        </form>
      )}

      {/* Register form */}
      {tab === 'register' && (
        <form onSubmit={handleRegister} className="space-y-3">
          {role === 'buyer' && (
            <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-xl p-3 mb-4">
              <p className="text-xs text-yellow-400 font-bold">Premium & Gold Upgrades</p>
              <p className="text-[10px] text-yellow-400/80 mt-1">
                Normal buyers do not need an account. Register here only to upgrade to Gold or Premium to unlock seller contacts and all product images.
              </p>
            </div>
          )}
          <PhoneField value={regPhone} onChange={e => { setRegPhone(e.target.value); setRegPhoneErr(''); }}
            error={regPhoneErr} result={regPhoneResult} inputCls={inputCls} lbl={lbl}
            phoneLabel={u.phone_number_required} phoneHint={u.phone_hint} />
          <div>
            <label className={lbl}><User size={9} className="inline mr-1" />{u.full_name_label}</label>
            <div className="relative">
              <User size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none" />
              <input required className={inputCls + ' pl-10'} placeholder={u.full_name_placeholder}
                value={regForm.name} onChange={setReg('name')} />
            </div>
          </div>
          {role === 'seller' && (
            <div>
              <label className={lbl}><Building2 size={9} className="inline mr-1" />{u.business_name_label} <span className="normal-case font-normal text-white/20">{u.optional_label}</span></label>
              <div className="relative">
                <Building2 size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none" />
                <input className={inputCls + ' pl-10'} placeholder={u.optional_label}
                  value={regForm.business_name} onChange={setReg('business_name')} />
              </div>
            </div>
          )}
          <div>
            <label className={lbl}><Mail size={9} className="inline mr-1" />{u.email_label} <span className="normal-case font-normal text-white/20">{u.optional_label}</span></label>
            <div className="relative">
              <Mail size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none" />
              <input type="email" className={inputCls + ' pl-10'} placeholder={u.email_placeholder}
                value={regForm.email} onChange={setReg('email')} />
            </div>
          </div>
          {role === 'seller' && (
            <div>
              <label className={lbl}><ShoppingBag size={9} className="inline mr-1" />{u.what_to_list}</label>
              <div className="space-y-2">
                {[
                  { value: 'sale', icon: '🛒', label: u.sell_products, desc: u.sell_products_desc },
                  { value: 'rent', icon: '🏠', label: u.rent_property, desc: u.rent_property_desc },
                  { value: 'both', icon: '🔄', label: u.both_option, desc: u.both_option_desc },
                ].map(pref => (
                  <label key={pref.value} className={`block cursor-pointer rounded-xl border-2 transition-all ${
                    regForm.seller_preference === pref.value
                      ? 'border-orange-500 bg-orange-500/10'
                      : 'border-zinc-200 dark:border-white/10 bg-zinc-50 dark:bg-white/5 hover:border-orange-300'
                  }`}>
                    <input type="radio" name="seller_preference" value={pref.value}
                      checked={regForm.seller_preference === pref.value}
                      onChange={() => setRegForm(p => ({ ...p, seller_preference: pref.value as 'sale' | 'rent' | 'both' }))}
                      className="sr-only" />
                    <div className="flex items-start gap-3 p-3">
                      <span className="text-2xl shrink-0">{pref.icon}</span>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-bold text-zinc-900 dark:text-white">{pref.label}</p>
                        <p className="text-xs text-zinc-500 dark:text-white/40 mt-0.5">{pref.desc}</p>
                      </div>
                      {regForm.seller_preference === pref.value && (
                        <CheckCircle2 size={16} className="text-orange-500 shrink-0 mt-0.5" />
                      )}
                    </div>
                  </label>
                ))}
              </div>
            </div>
          )}
          <div>
            <label className={lbl}><Lock size={9} className="inline mr-1" />{u.password_min_label}</label>
            <div className="relative">
              <Lock size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none" />
              <input required type={showRegPw ? 'text' : 'password'}
                className={inputCls + ' pl-10 pr-10'} placeholder="••••••••"
                value={regForm.password} onChange={setReg('password')} />
              <button type="button" onClick={() => setShowRegPw(v => !v)}
                className="absolute right-3.5 top-1/2 -translate-y-1/2 text-white/30 hover:text-white/60">
                {showRegPw ? <EyeOff size={14} /> : <Eye size={14} />}
              </button>
            </div>
            {regForm.password && (
              <PasswordStrength password={regForm.password} className="mt-2" />
            )}
          </div>
          <div>
            <label className={lbl}><Lock size={9} className="inline mr-1" />{u.confirm_password_label}</label>
            <div className="relative">
              <Lock size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none" />
              <input required type="password" className={inputCls + ' pl-10'} placeholder="••••••••"
                value={regForm.confirm} onChange={setReg('confirm')} />
            </div>
          </div>
          {refCode && tab === 'register' && (
            <div className="flex items-center gap-2 bg-orange-500/8 border border-orange-500/20 rounded-xl px-3 py-2">
              <CheckCircle2 size={12} className="text-orange-400 shrink-0" />
              <span className="text-[11px] text-orange-400">
                {u.invite_code_applied} <span className="font-mono font-bold">{refCode}</span>
              </span>
            </div>
          )}
          {error && <ErrorBanner message={error} type={errorType} />}
          <SubmitBtn loading={loading} label={u.create_account_arrow} accent={ROLE_ACCENT[role]} waitLabel={u.please_wait} />
          <p className="text-center text-xs text-zinc-400 dark:text-white/30">
            {u.already_have_account}{' '}
            <button type="button" onClick={() => setTab('login')}
              className="text-orange-400 font-bold hover:text-orange-300">{u.sign_in_link}</button>
          </p>
        </form>
      )}

      {role === 'seller' && (
        <div className="flex items-start gap-2 bg-orange-500/5 border border-orange-500/15 rounded-xl p-3 mt-2">
          <Zap size={12} className="text-orange-400 shrink-0 mt-0.5" />
          <p className="text-[10px] text-white/35 leading-relaxed">
            {u.seller_tier_info
              .replace('{gold}', Number(tierPrices.gold).toLocaleString())
              .replace('{premium}', Number(tierPrices.premium).toLocaleString())}
          </p>
        </div>
      )}
      {role === 'buyer' && (
        <div className="flex items-start gap-2 bg-yellow-500/5 border border-yellow-500/15 rounded-xl p-3 mt-2">
          <Star size={12} className="text-yellow-400 shrink-0 mt-0.5" />
          <p className="text-[10px] text-white/35 leading-relaxed">
            {u.buyer_tier_info
              .replace('{gold}', Number(tierPrices.gold).toLocaleString())
              .replace('{premium}', Number(tierPrices.premium).toLocaleString())}
          </p>
        </div>
      )}
    </div>
  );
}

function PhoneField({ value, onChange, error, result, inputCls, lbl, phoneLabel, phoneHint }: {
  value: string; onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  error: string; result: ReturnType<typeof validateEthioPhone> | null;
  inputCls: string; lbl: string; phoneLabel: string; phoneHint: string;
}) {
  return (
    <div>
      <label className={lbl}><Phone size={9} className="inline mr-1" />{phoneLabel}</label>
      <div className="relative">
        <Phone size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-white/30 pointer-events-none" />
        <input required type="tel" value={value} onChange={onChange}
          placeholder="09XXXXXXXX or 07XXXXXXXX"
          className={inputCls + ' pl-10'} autoFocus />
      </div>
      {result?.valid && (
        <p className="text-[10px] text-green-400 mt-1 flex items-center gap-1">
          <CheckCircle2 size={9} /> {result.normalized}
        </p>
      )}
      {error && (
        <p className="text-[10px] text-red-400 mt-1 flex items-center gap-1">
          <AlertCircle size={9} /> {error}
        </p>
      )}
      <p className="text-[10px] text-zinc-400 dark:text-white/20 mt-1">{phoneHint}</p>
    </div>
  );
}

function SubmitBtn({ loading, label, accent, waitLabel }: {
  loading: boolean; label: string; accent: string; waitLabel: string;
}) {
  return (
    <button type="submit" disabled={loading}
      className={`w-full ${accent} text-white py-3.5 rounded-xl font-black text-sm transition-all disabled:opacity-50 shadow-lg`}>
      {loading ? waitLabel : label}
    </button>
  );
}

export default function AuthPage() {
  const { t } = useLocaleContext();
  const u = (t as any).ui as Record<string, string>;
  return (
    <div className="min-h-screen bg-white dark:bg-zinc-950 flex items-center justify-center p-4">
      <div className="w-full max-w-sm">
        <div className="flex items-center gap-3 mb-8">
          <div className="w-10 h-10 bg-orange-600 rounded-xl flex items-center justify-center shadow-lg shadow-orange-600/30">
            <ShoppingBag size={20} className="text-white" />
          </div>
          <div>
            <p className="font-black text-zinc-900 dark:text-white text-base leading-none">{u.app_name}</p>
            <p className="text-[10px] text-zinc-400 dark:text-white/30 uppercase tracking-widest mt-0.5">{u.sign_in_or_create}</p>
          </div>
        </div>
        <Suspense fallback={null}>
          <AuthContent />
        </Suspense>
        <p className="text-center text-xs text-zinc-400 dark:text-white/20 mt-6">
          <Link href="/" className="hover:text-zinc-700 dark:hover:text-white/50 transition-colors">{u.back_to_marketplace}</Link>
        </p>
      </div>
    </div>
  );
}
