'use client';
import { useState, FormEvent, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { api, ApiError } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import { validateEthioPhone } from '@/lib/validations/common';
import PasswordStrength, { validatePasswordStrength } from '@/components/forms/PasswordStrength';
import { OTPInput } from '@/components/forms/OTPInput';
import { ErrorBanner } from '@/components/ui/ErrorBanner';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';
import { useLocale } from '@/hooks/useLocale';
import {
  Crown,
  Sparkles,
  ShieldCheck,
  CheckCircle2,
  Lock,
  Phone,
  User,
  Mail,
  ArrowRight,
  ArrowLeft,
  Building2,
  Zap,
  Star,
  Check,
  HelpCircle,
} from 'lucide-react';

const INPUT_CLS =
  'w-full bg-zinc-100 dark:bg-white/5 border border-zinc-200 dark:border-white/10 rounded-xl px-4 py-3 min-h-[44px] text-base sm:text-sm text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-white/30 outline-none transition-all focus:ring-2 focus:ring-[#1769E0]/30 focus:border-[#1769E0]';
const LBL_CLS = 'text-[10px] font-black uppercase tracking-widest text-zinc-500 dark:text-white/40 block mb-1.5';

export default function BuyerRegisterPage() {
  const router = useRouter();
  const { user, login } = useAuth();
  const { t } = useLocale();
  const btf = (t as any).buyer_tier_features ?? {};

  const [tier, setTier] = useState<'gold' | 'premium'>('gold');
  const [duration, setDuration] = useState<number>(30);
  const [paymentMethod, setPaymentMethod] = useState<'telebirr' | 'cbe'>('telebirr');

  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');
  const [phoneErr, setPhoneErr] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [payerPhone, setPayerPhone] = useState('');
  const [paymentRef, setPaymentRef] = useState(`TB-REG-${Date.now().toString().slice(-6)}`);

  const [step, setStep] = useState<'form' | 'otp'>('form');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // If already logged in, redirect to buyer dashboard
  useEffect(() => {
    if (user) {
      router.replace('/buyer/dashboard');
    }
  }, [user, router]);

  const phoneResult = phone ? validateEthioPhone(phone) : null;

  const tierPrices = {
    gold: { 30: 299, 90: 799, 365: 2899 },
    premium: { 30: 599, 90: 1599, 365: 5499 },
  };

  const currentPrice = tierPrices[tier][duration as 30 | 90 | 365] || (tier === 'gold' ? 299 : 599);

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');

    const pr = validateEthioPhone(phone);
    if (!pr.valid) {
      setPhoneErr(pr.error ?? 'Enter a valid Ethiopian phone number (e.g. 0911223344)');
      return;
    }
    setPhoneErr('');

    if (!name.trim()) {
      setError('Full Name is required');
      return;
    }

    if (password !== confirm) {
      setError('Passwords do not match');
      return;
    }

    const pwCheck = validatePasswordStrength(password);
    if (!pwCheck.valid) {
      setError(pwCheck.message);
      return;
    }

    setLoading(true);
    try {
      const regRes = await api.unifiedRegister({
        name: name.trim(),
        phone: pr.normalized,
        password,
        email: email.trim() || undefined,
        role: 'buyer',
        buyer_tier: tier,
        payment_method: paymentMethod,
        payment_ref: paymentRef.trim() || `TB-${Date.now()}`,
        payer_phone: payerPhone.trim() || pr.normalized,
        duration_days: duration,
      });

      try {
        localStorage.setItem('huntiq_last_phone', pr.normalized);
      } catch {}

      // If OTP disabled globally
      if (regRes?.otp_required === false) {
        if (regRes?.token && regRes?.user) {
          await login(regRes.token, undefined, regRes.user);
          router.replace('/buyer/dashboard');
          return;
        }
        const loginRes = await api.unifiedLogin({ phone: pr.normalized, password });
        await login(loginRes.token, undefined, loginRes.user);
        router.replace('/buyer/dashboard');
        return;
      }

      // OTP required
      setStep('otp');
    } catch (err: any) {
      console.error('Buyer register error:', err);
      if (err?.message?.startsWith('otp_pending:')) {
        api.sendOTP(pr.normalized, 'register', 'buyer').catch(() => {});
        setStep('otp');
      } else {
        setError(err?.message || 'Registration failed. Please verify your details.');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleOTP = async (code: string) => {
    setLoading(true);
    setError('');
    const pr = validateEthioPhone(phone);
    try {
      await api.verifyOTP(pr.normalized, code, 'register');
      const res = await api.unifiedLogin({ phone: pr.normalized, password });
      await login(res.token, undefined, res.user);
      router.replace('/buyer/dashboard');
    } catch (err: any) {
      setError(err?.message || 'Invalid or expired OTP code');
      setLoading(false);
    }
  };

  // ── OTP STEP ───────────────────────────────────────────────────────────────
  if (step === 'otp') {
    const pr = validateEthioPhone(phone);
    return (
      <div className="min-h-screen bg-zinc-950 text-zinc-100 flex items-center justify-center p-4">
        <div className="w-full max-w-md bg-zinc-900/90 border border-white/10 rounded-3xl p-6 sm:p-8 space-y-6">
          <button
            onClick={() => {
              setStep('form');
              setError('');
            }}
            className="flex items-center gap-1.5 text-xs text-zinc-400 hover:text-white transition-colors"
          >
            <ArrowLeft size={14} /> Back to details
          </button>
          <div>
            <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-amber-500/10 text-amber-400 border border-amber-500/20 mb-3">
              <Crown size={12} /> {tier === 'gold' ? 'Gold Buyer' : 'Premium VIP'} Membership
            </span>
            <h1 className="text-2xl font-black text-white">Verify Your Phone</h1>
            <p className="text-xs text-zinc-400 mt-1">
              Enter the 6-digit confirmation code sent to <strong className="text-white">{pr.normalized}</strong>.
            </p>
          </div>

          <OTPInput
            recipient={pr.normalized}
            onVerified={handleOTP}
            onResend={() => api.sendOTP(pr.normalized, 'register', 'buyer')}
            error={error}
            loading={loading}
          />

          <div className="bg-white/5 border border-white/10 rounded-2xl p-4 text-xs text-zinc-400 space-y-1">
            <p className="font-bold text-white flex items-center gap-1.5">
              <Zap size={14} className="text-amber-400" /> Dev Mode OTP Helper
            </p>
            <p>
              In local development, you can peek your plaintext OTP at:{' '}
              <a
                href={`http://localhost:8080/api/dev/otp?phone=${encodeURIComponent(pr.normalized)}`}
                target="_blank"
                rel="noreferrer"
                className="text-[#1769E0] dark:text-[#38A5FF] underline font-mono"
              >
                /api/dev/otp?phone={pr.normalized}
              </a>
            </p>
          </div>
        </div>
      </div>
    );
  }

  // ── REGISTRATION FORM ──────────────────────────────────────────────────────
  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100 flex flex-col justify-between py-10 px-4 sm:px-6 lg:px-8 selection:bg-amber-500 selection:text-black">
      <div className="w-full max-w-2xl mx-auto space-y-8">
        {/* Top Header */}
        <div className="text-center space-y-3">
          <Link href="/" className="inline-flex items-center gap-2 mb-2">
            <HuntiqLogo size="sm" />
            <span className="font-black text-2xl tracking-tight text-white">HUNTIQ</span>
          </Link>
          <div className="inline-flex items-center gap-2 px-3.5 py-1 rounded-full text-xs font-black tracking-wide uppercase bg-amber-500/10 text-amber-400 border border-amber-500/30">
            <Crown size={13} /> Exclusive Buyer Membership
          </div>
          <h1 className="text-3xl sm:text-4xl font-black text-white tracking-tight">
            Register as a Verified Buyer
          </h1>
          <p className="text-sm text-zinc-400 max-w-md mx-auto leading-relaxed">
            Free normal buyer registration is disabled. To guarantee trust, eliminate fraudulent spam, and unlock direct broker-verified wholesale pricing, buyers join HUNTIQ with an active privilege tier.
          </p>
        </div>

        {/* Tier Cards Selector */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {/* Gold Buyer */}
          <div
            onClick={() => setTier('gold')}
            className={`cursor-pointer rounded-3xl p-5 border transition-all relative ${
              tier === 'gold'
                ? 'bg-gradient-to-b from-amber-500/15 via-zinc-900 to-zinc-900 border-amber-500 shadow-xl shadow-amber-500/10 ring-2 ring-amber-500/30'
                : 'bg-zinc-900/50 border-white/10 hover:border-white/20 opacity-80 hover:opacity-100'
            }`}
          >
            {tier === 'gold' && (
              <span className="absolute top-4 right-4 w-6 h-6 rounded-full bg-amber-500 text-black flex items-center justify-center text-xs font-black shadow-md">
                <Check size={14} strokeWidth={3} />
              </span>
            )}
            <div className="flex items-center gap-2.5 mb-3 text-amber-400">
              <div className="w-9 h-9 rounded-xl bg-amber-500/20 flex items-center justify-center font-black">
                <Star size={18} />
              </div>
              <div>
                <h3 className="font-black text-base text-white">{btf.gold_buyer ?? 'Gold Buyer'}</h3>
                <p className="text-[10px] font-bold text-amber-400 uppercase tracking-widest">{btf.verified_trader ?? 'Verified Trader'}</p>
              </div>
            </div>
            <div className="mb-4">
              <span className="text-2xl font-black text-white">299 ETB</span>
              <span className="text-xs text-zinc-400"> / month</span>
            </div>
            <ul className="text-xs text-zinc-300 space-y-2">
              <li className="flex items-start gap-2">
                <CheckCircle2 size={13} className="text-amber-400 shrink-0 mt-0.5" />
                <span>{btf.gold_benefit_contact ?? 'Direct verified seller & landlord phone contacts'}</span>
              </li>
              <li className="flex items-start gap-2">
                <CheckCircle2 size={13} className="text-amber-400 shrink-0 mt-0.5" />
                <span>5% discount on all broker transactions</span>
              </li>
              <li className="flex items-start gap-2">
                <CheckCircle2 size={13} className="text-amber-400 shrink-0 mt-0.5" />
                <span>{btf.gold_benefit_sms ?? 'Priority SMS alerts for discounted listings'}</span>
              </li>
              <li className="flex items-start gap-2">
                <CheckCircle2 size={13} className="text-amber-400 shrink-0 mt-0.5" />
                <span>{btf.gold_benefit_badge ?? 'Official Gold Buyer Verified Badge'}</span>
              </li>
            </ul>
          </div>

          {/* Premium VIP Buyer */}
          <div
            onClick={() => setTier('premium')}
            className={`cursor-pointer rounded-3xl p-5 border transition-all relative ${
              tier === 'premium'
                ? 'bg-gradient-to-b from-purple-500/15 via-zinc-900 to-zinc-900 border-purple-500 shadow-xl shadow-purple-500/10 ring-2 ring-purple-500/30'
                : 'bg-zinc-900/50 border-white/10 hover:border-white/20 opacity-80 hover:opacity-100'
            }`}
          >
            {tier === 'premium' && (
              <span className="absolute top-4 right-4 w-6 h-6 rounded-full bg-purple-500 text-white flex items-center justify-center text-xs font-black shadow-md">
                <Check size={14} strokeWidth={3} />
              </span>
            )}
            <div className="flex items-center gap-2.5 mb-3 text-purple-400">
              <div className="w-9 h-9 rounded-xl bg-purple-500/20 flex items-center justify-center font-black">
                <Crown size={18} />
              </div>
              <div>
                <h3 className="font-black text-base text-white">{btf.premium_vip ?? 'Premium VIP'}</h3>
                <p className="text-[10px] font-bold text-purple-400 uppercase tracking-widest">{btf.maximum_privilege ?? 'Maximum Privilege'}</p>
              </div>
            </div>
            <div className="mb-4">
              <span className="text-2xl font-black text-white">599 ETB</span>
              <span className="text-xs text-zinc-400"> / month</span>
            </div>
            <ul className="text-xs text-zinc-300 space-y-2">
              <li className="flex items-start gap-2">
                <CheckCircle2 size={13} className="text-purple-400 shrink-0 mt-0.5" />
                <span>{btf.premium_benefit_all_gold ?? 'All Gold Buyer privileges included'}</span>
              </li>
              <li className="flex items-start gap-2">
                <CheckCircle2 size={13} className="text-purple-400 shrink-0 mt-0.5" />
                <span>Access to Private Off-Market Demand Boards</span>
              </li>
              <li className="flex items-start gap-2">
                <CheckCircle2 size={13} className="text-purple-400 shrink-0 mt-0.5" />
                <span>1-on-1 Dedicated Broker & AI Market Advisor</span>
              </li>
              <li className="flex items-start gap-2">
                <CheckCircle2 size={13} className="text-purple-400 shrink-0 mt-0.5" />
                <span>{btf.premium_benefit_legal_check ?? 'Free Title Deed & Vehicle Inspection Legal Checks'}</span>
              </li>
            </ul>
          </div>
        </div>

        {/* Registration & Payment Form */}
        <form onSubmit={handleSubmit} className="bg-zinc-900/80 border border-white/10 rounded-3xl p-6 sm:p-8 space-y-5">
          <div className="flex items-center justify-between pb-3 border-b border-white/10">
            <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Step 1: Your Account Details</span>
            <span className="text-xs font-bold text-amber-400">Tier: {tier === 'gold' ? 'Gold Buyer' : 'Premium VIP'}</span>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className={LBL_CLS}>Full Name *</label>
              <div className="relative">
                <User size={14} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none" />
                <input
                  required
                  type="text"
                  placeholder="Abebe Kebede"
                  value={name}
                  onChange={e => setName(e.target.value)}
                  className={INPUT_CLS + ' pl-10'}
                />
              </div>
            </div>

            <div>
              <label className={LBL_CLS}>Ethiopian Phone Number *</label>
              <div className="relative">
                <Phone size={14} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none" />
                <input
                  required
                  type="tel"
                  placeholder="0911223344"
                  value={phone}
                  onChange={e => {
                    setPhone(e.target.value);
                    setPhoneErr('');
                  }}
                  className={INPUT_CLS + ' pl-10'}
                />
              </div>
              {phoneErr && <p className="text-[11px] text-red-400 mt-1">{phoneErr}</p>}
            </div>
          </div>

          <div>
            <label className={LBL_CLS}>Email Address (Optional)</label>
            <div className="relative">
              <Mail size={14} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none" />
              <input
                type="email"
                placeholder="buyer@example.com"
                value={email}
                onChange={e => setEmail(e.target.value)}
                className={INPUT_CLS + ' pl-10'}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className={LBL_CLS}>Password *</label>
              <div className="relative">
                <Lock size={14} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none" />
                <input
                  required
                  type="password"
                  placeholder="••••••••"
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  className={INPUT_CLS + ' pl-10'}
                />
              </div>
              <PasswordStrength password={password} />
            </div>

            <div>
              <label className={LBL_CLS}>Confirm Password *</label>
              <div className="relative">
                <Lock size={14} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none" />
                <input
                  required
                  type="password"
                  placeholder="••••••••"
                  value={confirm}
                  onChange={e => setConfirm(e.target.value)}
                  className={INPUT_CLS + ' pl-10'}
                />
              </div>
            </div>
          </div>

          {/* Payment Section */}
          <div className="pt-4 border-t border-white/10 space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Step 2: Tier Payment & Activation</span>
              <span className="text-xs font-bold text-white bg-white/10 px-2.5 py-1 rounded-full">
                Total: {currentPrice} ETB / {duration} days
              </span>
            </div>

            {/* Payment Method Switcher */}
            <div className="grid grid-cols-2 gap-3">
              <button
                type="button"
                onClick={() => setPaymentMethod('telebirr')}
                className={`py-3 px-4 rounded-xl border text-xs font-black flex items-center justify-center gap-2 transition-all ${
                  paymentMethod === 'telebirr'
                    ? 'bg-amber-500/20 border-amber-500 text-amber-300 shadow-md shadow-amber-500/10'
                    : 'bg-white/5 border-white/10 text-zinc-400 hover:text-white'
                }`}
              >
                <Zap size={15} /> Telebirr (Instant Auto-Verify)
              </button>
              <button
                type="button"
                onClick={() => setPaymentMethod('cbe')}
                className={`py-3 px-4 rounded-xl border text-xs font-black flex items-center justify-center gap-2 transition-all ${
                  paymentMethod === 'cbe'
                    ? 'bg-purple-500/20 border-purple-500 text-purple-300 shadow-md shadow-purple-500/10'
                    : 'bg-white/5 border-white/10 text-zinc-400 hover:text-white'
                }`}
              >
                <Building2 size={15} /> CBE Bank Transfer
              </button>
            </div>

            {/* Payment Instructions */}
            <div className="bg-white/5 border border-white/10 rounded-2xl p-4 text-xs space-y-2 text-zinc-300">
              {paymentMethod === 'telebirr' ? (
                <>
                  <div className="flex items-center justify-between font-bold text-white">
                    <span>Telebirr Merchant Code:</span>
                    <span className="font-mono bg-amber-500/20 text-amber-300 px-2 py-0.5 rounded">123456</span>
                  </div>
                  <p className="text-[11px] text-zinc-400">
                    Pay <strong>{currentPrice} ETB</strong> via Telebirr app or *127#. In development mode, Telebirr tier payments are verified and activated instantly!
                  </p>
                </>
              ) : (
                <>
                  <div className="flex items-center justify-between font-bold text-white">
                    <span>CBE Account (HUNTIQ):</span>
                    <span className="font-mono bg-purple-500/20 text-purple-300 px-2 py-0.5 rounded">1000234567890</span>
                  </div>
                  <p className="text-[11px] text-zinc-400">
                    Transfer <strong>{currentPrice} ETB</strong> via CBE Birr or Commercial Bank of Ethiopia. Enter the bank reference code below.
                  </p>
                </>
              )}
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className={LBL_CLS}>Payer Phone Number</label>
                <input
                  type="text"
                  placeholder="Phone used to pay"
                  value={payerPhone}
                  onChange={e => setPayerPhone(e.target.value)}
                  className={INPUT_CLS}
                />
              </div>

              <div>
                <label className={LBL_CLS}>Payment Reference / Transaction ID</label>
                <input
                  required
                  type="text"
                  placeholder="e.g. TB-982341"
                  value={paymentRef}
                  onChange={e => setPaymentRef(e.target.value)}
                  className={INPUT_CLS}
                />
              </div>
            </div>
          </div>

          {error && <ErrorBanner message={error} />}

          <button
            disabled={loading}
            type="submit"
            className="w-full py-4 rounded-2xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-black font-black text-sm tracking-wide shadow-xl shadow-amber-500/20 transition-all flex items-center justify-center gap-2 active:scale-98 disabled:opacity-50"
          >
            {loading ? (
              'Processing Activation...'
            ) : (
              <>
                <Crown size={18} />
                Activate {tier === 'gold' ? 'Gold Buyer (299 ETB)' : 'Premium VIP (599 ETB)'}
                <ArrowRight size={16} />
              </>
            )}
          </button>

          <p className="text-center text-xs text-zinc-500 pt-2">
            Already have an active account?{' '}
            <Link href="/buyer/login" className="text-amber-400 font-bold hover:underline">
              Sign In Here
            </Link>
          </p>
        </form>
      </div>
    </div>
  );
}
