'use client';
import { useState, useEffect, FormEvent } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api, PaymentOrder } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { validateEthioPhone } from '@/lib/validations/common';
import {
  Crown,
  Star,
  CheckCircle,
  AlertCircle,
  Phone,
  ShieldCheck,
  Zap,
  Smartphone,
  Radio,
  Cpu,
  RefreshCw,
  XCircle,
  Clock,
  ExternalLink,
} from 'lucide-react';
import { INPUT_CLS_LIGHT as inputCls } from '@/lib/constants/styles';

const DURATIONS = [
  { label: '1 Week', days: 7 },
  { label: '1 Month', days: 30 },
  { label: '6 Months', days: 180 },
  { label: '1 Year', days: 365 },
];

type UpgradeStep = 'configure' | 'processing' | 'done' | 'failed';

export default function BuyerUpgradePage() {
  const { user, buyerTier } = useAuth();

  const availableTiers = (['gold', 'premium'] as const).filter(t =>
    buyerTier === 'gold' ? t === 'premium' : true
  );

  const [step, setStep] = useState<UpgradeStep>('configure');
  const [tier, setTier] = useState<'gold' | 'premium'>(buyerTier === 'gold' ? 'premium' : 'gold');
  const [duration, setDuration] = useState(30);
  const [provider, setProvider] = useState<string>('demo_mock');
  const [operatingMode, setOperatingMode] = useState<'demo' | 'production'>('demo');
  const [goldPrice, setGoldPrice] = useState(200);
  const [premiumPrice, setPremiumPrice] = useState(500);
  const [payerPhone, setPayerPhone] = useState(user?.phone ?? '');
  const [phoneError, setPhoneError] = useState('');
  const [otpCode, setOtpCode] = useState('');
  const [otpSent, setOtpSent] = useState(false);
  const [otpSending, setOtpSending] = useState(false);
  const [otpVerified, setOtpVerified] = useState(false);
  const [loading, setLoading] = useState(false);
  const [activeOrder, setActiveOrder] = useState<PaymentOrder | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    let isMounted = true;
    api.getPaymentPricing()
      .then(res => {
        if (!isMounted || !res) return;
        setOperatingMode(res.operating_mode);
        if (res.operating_mode === 'production') {
          setProvider('telebirr');
        } else {
          setProvider('demo_mock');
        }
        if (res.pricing?.buyer_gold_price_etb) setGoldPrice(res.pricing.buyer_gold_price_etb);
        if (res.pricing?.buyer_premium_price_etb) setPremiumPrice(res.pricing.buyer_premium_price_etb);
      })
      .catch(() => {});
    return () => { isMounted = false; };
  }, []);

  const basePrice = tier === 'gold' ? goldPrice : premiumPrice;
  const calculatedPrice = Math.round(basePrice * (duration / 30.0));

  const handleSendOTP = async () => {
    if (!user?.token) return;
    const pr = validateEthioPhone(payerPhone);
    if (!pr.valid) {
      setPhoneError(pr.error ?? 'Please enter a valid Ethiopian phone number');
      return;
    }
    setPhoneError('');
    setOtpSending(true);
    setError('');
    try {
      await api.sendMembershipOTP({ phone: pr.normalized, target_tier: tier }, user.token);
      setOtpSent(true);
    } catch (err: any) {
      setError(err.message ?? 'Failed to send verification code');
    } finally {
      setOtpSending(false);
    }
  };

  const handleVerifyOTP = async () => {
    if (!user?.token || !otpCode.trim()) return;
    setLoading(true);
    setError('');
    try {
      const pr = validateEthioPhone(payerPhone);
      await api.verifyMembershipOTP({ phone: pr.normalized, code: otpCode.trim() }, user.token);
      setOtpVerified(true);
    } catch (err: any) {
      setError(err.message ?? 'Invalid or expired OTP code');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!user?.token) return;

    const pr = validateEthioPhone(payerPhone);
    if (!pr.valid) {
      setPhoneError(pr.error ?? 'Enter a valid Ethiopian phone number');
      return;
    }
    setPhoneError('');

    setLoading(true);
    setError('');
    try {
      const order = await api.checkoutPayment({
        order_type: 'tier_upgrade',
        target_tier: tier,
        duration_days: duration,
        provider,
        payer_phone: pr.normalized,
        otp_verified: otpVerified,
      }, user.token);

      setActiveOrder(order);
      setStep('processing');
    } catch (err: any) {
      setError(err.message ?? 'Something went wrong initiating payment.');
    } finally {
      setLoading(false);
    }
  };

  const handleSimulate = async (action: 'success' | 'failed' | 'cancelled' | 'expired') => {
    if (!user?.token || !activeOrder) return;
    setLoading(true);
    setError('');
    try {
      const res = await api.simulateDemoPayment({
        order_id: activeOrder.order_id,
        action,
      }, user.token);
      setActiveOrder(res);
      if (res.status === 'SUCCESSFUL') {
        setStep('done');
      } else {
        setStep('failed');
      }
    } catch (err: any) {
      setError(err.message ?? 'Simulation error');
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyNow = async () => {
    if (!user?.token || !activeOrder) return;
    setLoading(true);
    setError('');
    try {
      const res = await api.verifyPaymentOrder(activeOrder.order_id, user.token);
      setActiveOrder(res);
      if (res.status === 'SUCCESSFUL') {
        setStep('done');
      } else if (res.status === 'FAILED' || res.status === 'CANCELLED' || res.status === 'EXPIRED') {
        setStep('failed');
      }
    } catch (err: any) {
      setError(err.message ?? 'Failed to verify transaction');
    } finally {
      setLoading(false);
    }
  };

  if (buyerTier === 'premium') {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-3 text-zinc-400">
        <Crown size={32} className="text-purple-500" />
        <p className="font-bold text-zinc-900 dark:text-white">You are already Premium</p>
        <p className="text-sm">Maximum tier — all features unlocked.</p>
      </div>
    );
  }

  if (step === 'done') {
    return (
      <div className="flex flex-col items-center justify-center min-h-[50vh] gap-4 max-w-sm mx-auto text-center px-4">
        <div className="w-16 h-16 rounded-full bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center text-emerald-600 shadow-xl shadow-emerald-500/20">
          <CheckCircle size={36} strokeWidth={2.5} />
        </div>
        <div>
          <p className="font-black text-zinc-900 dark:text-white text-xl">Tier Upgrade Activated!</p>
          <p className="text-sm text-zinc-500 mt-2">
            Your account is now upgraded to <strong className="capitalize text-zinc-800 dark:text-zinc-200">{tier}</strong> tier with instant automatic activation.
          </p>
          <div className="inline-flex items-center gap-1.5 mt-3 px-3 py-1 bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-800 rounded-full text-xs font-bold text-emerald-600 dark:text-emerald-400">
            <ShieldCheck size={14} /> System Auto-Approved · No Admin Wait
          </div>
        </div>
        <button
          onClick={() => { window.location.href = '/buyer/dashboard'; }}
          className="bg-[#1769E0] text-white font-bold px-6 py-2.5 rounded-xl hover:bg-[#1055BD] transition-all"
        >
          Go to Dashboard
        </button>
      </div>
    );
  }

  if (step === 'failed') {
    return (
      <div className="flex flex-col items-center justify-center min-h-[50vh] gap-4 max-w-sm mx-auto text-center px-4">
        <div className="w-16 h-16 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center text-red-600 shadow-xl shadow-red-500/20">
          <XCircle size={36} strokeWidth={2.5} />
        </div>
        <div>
          <p className="font-black text-zinc-900 dark:text-white text-xl">Upgrade Incomplete</p>
          <p className="text-sm text-zinc-500 mt-2">
            {activeOrder?.failure_reason || 'Payment could not be verified or was declined.'}
          </p>
        </div>
        <button
          onClick={() => setStep('configure')}
          className="bg-zinc-800 text-white font-bold px-6 py-2.5 rounded-xl text-xs hover:bg-zinc-700 transition-all"
        >
          Try Again
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-5 max-w-md mx-auto py-6 px-4">
      <div>
        <h1 className="text-xl font-black text-zinc-900 dark:text-white">Upgrade Buyer Tier</h1>
        <p className="text-xs text-zinc-500 mt-0.5">Current tier: <span className="font-bold capitalize">{buyerTier}</span></p>
      </div>

      {step === 'processing' && activeOrder ? (
        <div className="border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5 bg-white dark:bg-zinc-900 shadow-sm space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs text-zinc-400 uppercase font-mono">Order Reference</p>
              <p className="text-sm font-bold font-mono text-zinc-800 dark:text-zinc-200">{activeOrder.order_id}</p>
            </div>
            <span className="px-2.5 py-0.5 rounded-full text-xs font-bold uppercase bg-amber-500/10 text-amber-600 border border-amber-500/20">
              {activeOrder.status}
            </span>
          </div>

          {activeOrder.checkout_url && (
            <div className="p-3 bg-blue-50 dark:bg-blue-950/30 border border-blue-200 dark:border-blue-800/40 rounded-xl space-y-2">
              <p className="text-xs text-blue-900 dark:text-blue-200 font-semibold">Payment Link Ready</p>
              <a
                href={activeOrder.checkout_url}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center gap-1.5 text-xs font-bold text-[#1769E0] hover:underline"
              >
                Open Provider Checkout Page <ExternalLink size={12} />
              </a>
            </div>
          )}

          {/* DEMO Simulation Panel */}
          {activeOrder.is_demo && (
            <div className="border border-dashed border-amber-400 dark:border-amber-600 bg-amber-50/50 dark:bg-amber-950/20 rounded-xl p-3.5 space-y-3">
              <div className="flex items-center gap-1.5 text-xs font-black text-amber-700 dark:text-amber-300">
                <Cpu size={14} /> FinTech Sandbox Simulator
              </div>
              <p className="text-[11px] text-zinc-500 dark:text-zinc-400">
                Simulate payment gateway callback in DEMO mode:
              </p>
              <div className="grid grid-cols-2 gap-2">
                <button
                  type="button"
                  onClick={() => handleSimulate('success')}
                  disabled={loading}
                  className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold py-2 rounded-lg flex items-center justify-center gap-1"
                >
                  <CheckCircle size={12} /> Simulate Success
                </button>
                <button
                  type="button"
                  onClick={() => handleSimulate('failed')}
                  disabled={loading}
                  className="bg-red-600 hover:bg-red-700 text-white text-xs font-bold py-2 rounded-lg flex items-center justify-center gap-1"
                >
                  <XCircle size={12} /> Simulate Decline
                </button>
                <button
                  type="button"
                  onClick={() => handleSimulate('cancelled')}
                  disabled={loading}
                  className="bg-zinc-600 hover:bg-zinc-700 text-white text-xs font-bold py-2 rounded-lg flex items-center justify-center gap-1"
                >
                  Simulate Cancel
                </button>
                <button
                  type="button"
                  onClick={() => handleSimulate('expired')}
                  disabled={loading}
                  className="bg-amber-600 hover:bg-amber-700 text-white text-xs font-bold py-2 rounded-lg flex items-center justify-center gap-1"
                >
                  <Clock size={12} /> Simulate Timeout
                </button>
              </div>
            </div>
          )}

          <div className="pt-2 flex gap-2">
            <button
              type="button"
              onClick={handleVerifyNow}
              disabled={loading}
              className="flex-1 bg-[#1769E0] hover:bg-[#1055BD] text-white py-2.5 rounded-xl text-xs font-bold flex items-center justify-center gap-1.5"
            >
              <RefreshCw size={12} className={loading ? 'animate-spin' : ''} />
              Verify Transaction Status
            </button>
            <button
              type="button"
              onClick={() => setStep('configure')}
              className="px-4 py-2.5 rounded-xl border border-zinc-200 dark:border-zinc-800 text-xs font-bold text-zinc-600 dark:text-zinc-400"
            >
              Back
            </button>
          </div>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Tier selector */}
          <div className={`grid gap-3 ${availableTiers.length === 1 ? 'grid-cols-1' : 'grid-cols-2'}`}>
            {availableTiers.map(t => (
              <button key={t} type="button" onClick={() => setTier(t)}
                className={`flex flex-col items-center gap-2 py-4 px-3 rounded-2xl border-2 transition-all ${tier === t
                  ? t === 'gold' ? 'border-yellow-500 bg-yellow-500/10' : 'border-purple-500 bg-purple-500/10'
                  : 'border-zinc-200 dark:border-zinc-700 hover:border-zinc-300 dark:hover:border-zinc-600'
                  }`}>
                {t === 'gold'
                  ? <Star size={20} className="text-yellow-500" />
                  : <Crown size={20} className="text-purple-500" />}
                <span className={`font-black text-sm capitalize ${t === 'gold' ? 'text-yellow-600 dark:text-yellow-400' : 'text-purple-600 dark:text-purple-400'}`}>{t}</span>
                <span className="text-xs font-bold text-zinc-500">ETB {(t === 'gold' ? goldPrice : premiumPrice).toLocaleString()} / mo</span>
              </button>
            ))}
          </div>

          {/* Duration selector */}
          <div>
            <label className="text-xs font-black uppercase tracking-wider text-zinc-500 mb-2 block">Duration</label>
            <div className="grid grid-cols-4 gap-2">
              {DURATIONS.map(d => (
                <button key={d.days} type="button" onClick={() => setDuration(d.days)}
                  className={`py-2 text-xs font-bold rounded-xl border transition-all ${duration === d.days
                    ? 'border-[#1769E0] bg-[#1769E0] text-white shadow-xs'
                    : 'border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-400 hover:border-zinc-300'
                    }`}>
                  {d.label}
                </button>
              ))}
            </div>
          </div>

          {/* Price summary */}
          <div className="flex items-center justify-between p-3.5 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-200 dark:border-zinc-700">
            <div>
              <span className="text-xs font-bold text-zinc-500">Total Authoritative Price</span>
              <p className="text-xl font-black text-zinc-900 dark:text-white">ETB {calculatedPrice.toLocaleString()}</p>
            </div>
            <span className="text-xs font-bold text-emerald-600 bg-emerald-500/10 px-2.5 py-1 rounded-full border border-emerald-500/20 flex items-center gap-1">
              <Zap size={11} /> Auto-Activation
            </span>
          </div>

          {/* Gateway selector */}
          <div>
            <label className="text-xs font-black uppercase tracking-wider text-zinc-500 mb-2 block">Payment Rail</label>
            <div className="space-y-2">
              {operatingMode === 'demo' && (
                <label className={`flex items-center justify-between p-3 rounded-xl border cursor-pointer transition-all ${
                  provider === 'demo_mock'
                    ? 'border-[#1769E0] bg-[#1769E0]/5 dark:bg-[#1769E0]/10'
                    : 'border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900'
                }`}>
                  <div className="flex items-center gap-2.5">
                    <Radio size={16} className={provider === 'demo_mock' ? 'text-[#1769E0]' : 'text-zinc-400'} />
                    <div>
                      <p className="text-xs font-bold text-zinc-900 dark:text-white">FinTech Sandbox Simulator</p>
                      <p className="text-[10px] text-zinc-500">Test mode active</p>
                    </div>
                  </div>
                  <input
                    type="radio"
                    name="provider"
                    value="demo_mock"
                    checked={provider === 'demo_mock'}
                    onChange={() => setProvider('demo_mock')}
                    className="hidden"
                  />
                </label>
              )}

              <label className={`flex items-center justify-between p-3 rounded-xl border cursor-pointer transition-all ${
                provider === 'telebirr'
                  ? 'border-[#1769E0] bg-[#1769E0]/5 dark:bg-[#1769E0]/10'
                  : 'border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900'
              }`}>
                <div className="flex items-center gap-2.5">
                  <Smartphone size={16} className={provider === 'telebirr' ? 'text-[#1769E0]' : 'text-zinc-400'} />
                  <div>
                    <p className="text-xs font-bold text-zinc-900 dark:text-white">Telebirr SuperApp Direct</p>
                    <p className="text-[10px] text-zinc-500">Official Mobile Wallet</p>
                  </div>
                </div>
                <input
                  type="radio"
                  name="provider"
                  value="telebirr"
                  checked={provider === 'telebirr'}
                  onChange={() => setProvider('telebirr')}
                  className="hidden"
                />
              </label>

              <label className={`flex items-center justify-between p-3 rounded-xl border cursor-pointer transition-all ${
                provider === 'ethio_telecom'
                  ? 'border-[#1769E0] bg-[#1769E0]/5 dark:bg-[#1769E0]/10'
                  : 'border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900'
              }`}>
                <div className="flex items-center gap-2.5">
                  <Phone size={16} className={provider === 'ethio_telecom' ? 'text-[#1769E0]' : 'text-zinc-400'} />
                  <div>
                    <p className="text-xs font-bold text-zinc-900 dark:text-white">Ethio Telecom TelePay</p>
                    <p className="text-[10px] text-zinc-500">Carrier Billing</p>
                  </div>
                </div>
                <input
                  type="radio"
                  name="provider"
                  value="ethio_telecom"
                  checked={provider === 'ethio_telecom'}
                  onChange={() => setProvider('ethio_telecom')}
                  className="hidden"
                />
              </label>
            </div>
          </div>

          {/* Payer Phone */}
          <div>
            <label className="text-xs font-bold text-zinc-700 dark:text-zinc-300 mb-1.5 block">Payer Phone Number</label>
            <input
              type="tel"
              value={payerPhone}
              onChange={e => { setPayerPhone(e.target.value); setPhoneError(''); }}
              placeholder="0911000000 or +251911000000"
              className={inputCls}
            />
            {phoneError && <p className="text-xs text-red-500 mt-1">{phoneError}</p>}
          </div>

          {error && (
            <div className="flex items-center gap-2 text-xs text-red-500 bg-red-500/10 border border-red-500/20 rounded-xl px-3 py-2">
              <AlertCircle size={14} className="shrink-0" /> {error}
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-[#1769E0] hover:bg-[#1055BD] text-white font-bold py-3 rounded-2xl transition-all shadow-md shadow-[#1769E0]/20 flex items-center justify-center gap-2 text-sm"
          >
            {loading ? <Spinner size={14} /> : `Upgrade to ${tier.toUpperCase()} · ETB ${calculatedPrice.toLocaleString()}`}
          </button>
        </form>
      )}
    </div>
  );
}
