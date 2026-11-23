'use client';
import { useEffect, useState, Suspense, FormEvent } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { useLocale } from '@/hooks/useLocale';
import { api, PaymentOrder } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { validateEthioPhone } from '@/lib/validations/common';
import { INPUT_CLS } from '@/lib/constants/styles';
import {
  CheckCircle,
  AlertCircle,
  Phone,
  ExternalLink,
  ShieldCheck,
  Zap,
  Smartphone,
  Radio,
  Cpu,
  RefreshCw,
  XCircle,
  Clock,
} from 'lucide-react';

type Step = 'select_provider' | 'otp_verify' | 'processing' | 'done' | 'failed';

function PayContent() {
  const { user } = useAuth();
  const { t } = useLocale();
  const router = useRouter();
  const params = useSearchParams();

  const productId = parseInt(params.get('product_id') ?? '0', 10) || 0;
  const initialFee = parseFloat(params.get('fee') ?? '0') || 0;
  const productTitle = decodeURIComponent(params.get('title') ?? '');

  const [step, setStep] = useState<Step>('select_provider');
  const [provider, setProvider] = useState<string>('demo_mock');
  const [operatingMode, setOperatingMode] = useState<'demo' | 'production'>('demo');
  const [authoritativeFee, setAuthoritativeFee] = useState<number>(initialFee || 50);
  const [payerPhone, setPhone] = useState(user?.phone ?? '');
  const [phoneError, setPhoneError] = useState('');
  const [otpCode, setOtpCode] = useState('');
  const [otpSending, setOtpSending] = useState(false);
  const [otpSent, setOtpSent] = useState(false);
  const [otpVerified, setOtpVerified] = useState(false);
  const [submitting, setSub] = useState(false);
  const [activeOrder, setActiveOrder] = useState<PaymentOrder | null>(null);
  const [error, setError] = useState('');

  // Fetch authoritative pricing and mode
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
        if (res.pricing?.listing_fee_etb) {
          setAuthoritativeFee(res.pricing.listing_fee_etb);
        }
      })
      .catch(() => {
        // Fallback default
      });
    return () => { isMounted = false; };
  }, []);

  useEffect(() => {
    if (!productId) router.replace('/seller/products');
  }, [productId, router]);

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
      await api.sendMembershipOTP({ phone: pr.normalized, target_tier: 'gold' }, user.token);
      setOtpSent(true);
    } catch (err: any) {
      setError(err.message ?? 'Failed to send Fida OTP');
    } finally {
      setOtpSending(false);
    }
  };

  const handleVerifyOTP = async () => {
    if (!user?.token || !otpCode.trim()) return;
    setSub(true);
    setError('');
    try {
      const pr = validateEthioPhone(payerPhone);
      await api.verifyMembershipOTP({ phone: pr.normalized, code: otpCode.trim() }, user.token);
      setOtpVerified(true);
      setError('');
    } catch (err: any) {
      setError(err.message ?? 'Invalid or expired OTP code');
    } finally {
      setSub(false);
    }
  };

  const handleInitiatePayment = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!user?.token) return;

    const pr = validateEthioPhone(payerPhone);
    if (!pr.valid) {
      setPhoneError(pr.error ?? 'Please enter a valid Ethiopian phone number');
      return;
    }
    setPhoneError('');
    setSub(true);
    setError('');

    try {
      const order = await api.checkoutPayment({
        order_type: 'listing_fee',
        product_id: productId,
        provider,
        payer_phone: pr.normalized,
        otp_verified: otpVerified,
      }, user.token);

      setActiveOrder(order);
      setStep('processing');

      // If simulated or demo mode, user can simulate outcome directly
    } catch (err: any) {
      setError(err.message ?? 'Failed to initiate payment');
    } finally {
      setSub(false);
    }
  };

  const handleSimulate = async (action: 'success' | 'failed' | 'cancelled' | 'expired') => {
    if (!user?.token || !activeOrder) return;
    setSub(true);
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
      setSub(false);
    }
  };

  const handleVerifyNow = async () => {
    if (!user?.token || !activeOrder) return;
    setSub(true);
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
      setSub(false);
    }
  };

  // ── Done Screen ────────────────────────────────────────────────────────────
  if (step === 'done') {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-5 text-center px-4">
        <div className="w-20 h-20 bg-emerald-500 rounded-full flex items-center justify-center shadow-xl shadow-emerald-500/30">
          <CheckCircle size={40} className="text-white" strokeWidth={2.5} />
        </div>
        <div>
          <h1 className="text-2xl font-black text-zinc-900 dark:text-white">Payment Verified & Activated!</h1>
          <p className="text-zinc-500 mt-2 text-sm max-w-sm mx-auto">
            Your listing fee of <strong>ETB {authoritativeFee.toLocaleString()}</strong> has been verified via {activeOrder?.provider?.toUpperCase() || 'HUNTIQ Pay'}.
          </p>
          <div className="inline-flex items-center gap-1.5 mt-3 px-3 py-1 bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-800 rounded-full text-xs font-bold text-emerald-600 dark:text-emerald-400">
            <ShieldCheck size={14} /> Instant System Activation · Listing is Live
          </div>
          {activeOrder?.provider_tx_id && (
            <p className="text-[11px] text-zinc-400 font-mono mt-2">Ref: {activeOrder.provider_tx_id}</p>
          )}
        </div>
        <button
          onClick={() => router.push('/seller/products')}
          className="bg-[#1769E0] hover:bg-[#1055BD] text-white px-8 py-3 rounded-2xl font-bold min-h-[44px] transition-all shadow-md shadow-[#1769E0]/20"
        >
          View My Products
        </button>
      </div>
    );
  }

  // ── Failed Screen ──────────────────────────────────────────────────────────
  if (step === 'failed') {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-5 text-center px-4">
        <div className="w-20 h-20 bg-red-500 rounded-full flex items-center justify-center shadow-xl shadow-red-500/30">
          <XCircle size={40} className="text-white" strokeWidth={2.5} />
        </div>
        <div>
          <h1 className="text-2xl font-black text-zinc-900 dark:text-white">Payment Unsuccessful</h1>
          <p className="text-zinc-500 mt-2 text-sm max-w-sm mx-auto">
            {activeOrder?.failure_reason || 'The transaction was declined, cancelled, or expired.'}
          </p>
        </div>
        <button
          onClick={() => setStep('select_provider')}
          className="bg-zinc-800 hover:bg-zinc-700 text-white px-6 py-2.5 rounded-xl font-bold text-xs"
        >
          Try Again
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto space-y-5 py-6 px-4">
      <div>
        <h1 className="text-xl font-black text-zinc-900 dark:text-white">Pay Listing Fee</h1>
        {productTitle && <p className="text-xs text-zinc-500 mt-0.5 truncate">{productTitle}</p>}
      </div>

      {/* Fee summary */}
      <div className="bg-[#1769E0]/5 dark:bg-[#1769E0]/20 border border-[#1769E0]/20 rounded-2xl p-4">
        <div className="flex items-center justify-between mb-1">
          <p className="text-xs text-[#1769E0] dark:text-[#38A5FF] font-black uppercase tracking-widest">Amount Due</p>
          <span className="flex items-center gap-1 text-[11px] font-bold text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/40 px-2 py-0.5 rounded-full">
            <Zap size={11} /> Automatic Verification
          </span>
        </div>
        <p className="text-3xl font-black text-[#1769E0] dark:text-[#38A5FF]">ETB {authoritativeFee.toLocaleString()}</p>
        <p className="text-xs text-zinc-500 mt-1">Authoritative fee for publishing Product #{productId}.</p>
      </div>

      {/* Step: Processing Order / Verification Sandbox */}
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
                The application is operating in controlled DEMO mode. Simulate provider response:
              </p>
              <div className="grid grid-cols-2 gap-2">
                <button
                  type="button"
                  onClick={() => handleSimulate('success')}
                  disabled={submitting}
                  className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold py-2 rounded-lg flex items-center justify-center gap-1"
                >
                  <CheckCircle size={12} /> Simulate Success
                </button>
                <button
                  type="button"
                  onClick={() => handleSimulate('failed')}
                  disabled={submitting}
                  className="bg-red-600 hover:bg-red-700 text-white text-xs font-bold py-2 rounded-lg flex items-center justify-center gap-1"
                >
                  <XCircle size={12} /> Simulate Decline
                </button>
                <button
                  type="button"
                  onClick={() => handleSimulate('cancelled')}
                  disabled={submitting}
                  className="bg-zinc-600 hover:bg-zinc-700 text-white text-xs font-bold py-2 rounded-lg flex items-center justify-center gap-1"
                >
                  Simulate Cancel
                </button>
                <button
                  type="button"
                  onClick={() => handleSimulate('expired')}
                  disabled={submitting}
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
              disabled={submitting}
              className="flex-1 bg-[#1769E0] hover:bg-[#1055BD] text-white py-2.5 rounded-xl text-xs font-bold flex items-center justify-center gap-1.5"
            >
              <RefreshCw size={12} className={submitting ? 'animate-spin' : ''} />
              Verify Transaction Status
            </button>
            <button
              type="button"
              onClick={() => setStep('select_provider')}
              className="px-4 py-2.5 rounded-xl border border-zinc-200 dark:border-zinc-800 text-xs font-bold text-zinc-600 dark:text-zinc-400"
            >
              Back
            </button>
          </div>
        </div>
      ) : (
        /* Step: Select Provider and Enter Phone */
        <form onSubmit={handleInitiatePayment} className="space-y-4">
          <div>
            <label className="text-xs font-black uppercase tracking-wider text-zinc-500 mb-2 block">
              Select Payment Gateway
            </label>
            <div className="space-y-2">
              {operatingMode === 'demo' && (
                <label className={`flex items-center justify-between p-3.5 rounded-xl border cursor-pointer transition-all ${
                  provider === 'demo_mock'
                    ? 'border-[#1769E0] bg-[#1769E0]/5 dark:bg-[#1769E0]/10'
                    : 'border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900'
                }`}>
                  <div className="flex items-center gap-3">
                    <Radio size={16} className={provider === 'demo_mock' ? 'text-[#1769E0]' : 'text-zinc-400'} />
                    <div>
                      <p className="text-xs font-black text-zinc-900 dark:text-white flex items-center gap-1.5">
                        HUNTIQ FinTech Demo Sandbox
                        <span className="px-1.5 py-0.2 bg-amber-500/10 text-amber-600 text-[10px] rounded-md font-bold">
                          DEMO ACTIVE
                        </span>
                      </p>
                      <p className="text-[11px] text-zinc-500">Simulate success, failure, cancel, and timeout</p>
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

              <label className={`flex items-center justify-between p-3.5 rounded-xl border cursor-pointer transition-all ${
                provider === 'telebirr'
                  ? 'border-[#1769E0] bg-[#1769E0]/5 dark:bg-[#1769E0]/10'
                  : 'border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900'
              }`}>
                <div className="flex items-center gap-3">
                  <Smartphone size={16} className={provider === 'telebirr' ? 'text-[#1769E0]' : 'text-zinc-400'} />
                  <div>
                    <p className="text-xs font-black text-zinc-900 dark:text-white">Telebirr SuperApp Direct</p>
                    <p className="text-[11px] text-zinc-500">Official Mobile Wallet H5 Checkout</p>
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

              <label className={`flex items-center justify-between p-3.5 rounded-xl border cursor-pointer transition-all ${
                provider === 'ethio_telecom'
                  ? 'border-[#1769E0] bg-[#1769E0]/5 dark:bg-[#1769E0]/10'
                  : 'border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900'
              }`}>
                <div className="flex items-center gap-3">
                  <Phone size={16} className={provider === 'ethio_telecom' ? 'text-[#1769E0]' : 'text-zinc-400'} />
                  <div>
                    <p className="text-xs font-black text-zinc-900 dark:text-white">Ethio Telecom TelePay</p>
                    <p className="text-[11px] text-zinc-500">Carrier Billing & SIM prompt validation</p>
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

          <div>
            <label className="text-xs font-bold text-zinc-700 dark:text-zinc-300 mb-1.5 block">
              Payer Phone Number
            </label>
            <input
              type="tel"
              value={payerPhone}
              onChange={e => { setPhone(e.target.value); setPhoneError(''); }}
              placeholder="0911000000 or +251911000000"
              className={INPUT_CLS}
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
            disabled={submitting}
            className="w-full bg-[#1769E0] hover:bg-[#1055BD] text-white py-3 rounded-2xl font-bold text-sm min-h-[44px] flex items-center justify-center gap-2 shadow-md shadow-[#1769E0]/20"
          >
            {submitting ? <Spinner size={14} /> : `Proceed to Pay ETB ${authoritativeFee.toLocaleString()}`}
          </button>
        </form>
      )}
    </div>
  );
}

export default function PayListingFeePage() {
  return (
    <Suspense fallback={<div className="flex justify-center py-20"><Spinner size={32} /></div>}>
      <PayContent />
    </Suspense>
  );
}
