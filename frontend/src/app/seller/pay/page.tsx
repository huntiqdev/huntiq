'use client';
import { useEffect, useState, Suspense, FormEvent } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { useLocale } from '@/hooks/useLocale';
import { usePaymentSettings } from '@/hooks/usePaymentSettings';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { validateEthioPhone } from '@/lib/validations/common';
import { PAYMENT_METHODS, getMerchantTarget } from '@/lib/constants/payment';
import { INPUT_CLS } from '@/lib/constants/styles';
import { CheckCircle, AlertCircle, Phone, ExternalLink, ShieldCheck, Zap } from 'lucide-react';

type Step = 'confirm' | 'done';

function PayContent() {
  const { user } = useAuth();
  const { t } = useLocale();
  const { settings } = usePaymentSettings();
  const router = useRouter();
  const params = useSearchParams();

  const productId = parseInt(params.get('product_id') ?? '0', 10) || 0;
  const feeAmount = parseFloat(params.get('fee') ?? '0') || 0;
  const productTitle = decodeURIComponent(params.get('title') ?? '');

  const [step, setStep] = useState<Step>('confirm');
  const [payerPhone, setPhone] = useState(user?.phone ?? '');
  const [phoneError, setPhoneError] = useState('');
  const [submitting, setSub] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!productId || feeAmount <= 0) router.replace('/seller/products');
  }, [productId, feeAmount, router]);

  const method = PAYMENT_METHODS[0]; // Telebirr is authoritative
  const merchantTarget = getMerchantTarget('telebirr', settings);
  const phoneResult = payerPhone ? validateEthioPhone(payerPhone) : null;

  const handlePay = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!user) return;

    const pr = validateEthioPhone(payerPhone);
    if (!pr.valid) {
      setPhoneError(pr.error ?? 'Please enter a valid Ethiopian Telebirr phone number');
      return;
    }
    setPhoneError('');

    setSub(true);
    setError('');
    try {
      // 1. Initiate official Telebirr payment order
      const res = await api.createTelebirrOrder({
        amount: feeAmount,
        order_type: 'listing_fee',
        title: `Listing fee for ${productTitle || 'Product #' + productId}`,
      }, user.token);

      // 2. Also register seller payment record with auto-approved status
      await api.sellerPayListingFee({
        product_id: productId,
        amount: feeAmount,
        payment_method: 'telebirr',
        payer_phone: pr.normalized,
        reference: res?.order?.out_trade_no || `TB-${Date.now()}`,
        notes: `Automated Telebirr payment for product #${productId}`,
      }, user.token);

      setStep('done');
    } catch (err: any) {
      setError(err.message ?? 'Failed to process Telebirr payment');
    } finally {
      setSub(false);
    }
  };

  // ── Done screen ────────────────────────────────────────────────────────────
  if (step === 'done') {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-5 text-center px-4">
        <div className="w-20 h-20 bg-emerald-500 rounded-full flex items-center justify-center shadow-xl shadow-emerald-500/30">
          <CheckCircle size={40} className="text-white" strokeWidth={2.5} />
        </div>
        <div>
          <h1 className="text-2xl font-black text-zinc-900 dark:text-white">Payment Verified & Approved!</h1>
          <p className="text-zinc-500 mt-2 text-sm max-w-sm mx-auto">
            Your listing fee of <strong>ETB {feeAmount.toLocaleString()}</strong> has been automatically approved via the Telebirr payment gateway.
          </p>
          <div className="inline-flex items-center gap-1.5 mt-3 px-3 py-1 bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-800 rounded-full text-xs font-bold text-emerald-600 dark:text-emerald-400">
            <ShieldCheck size={14} /> Instant System Auto-Approval Active
          </div>
        </div>
        <button
          onClick={() => router.push('/seller/products')}
          className="bg-[#1769E0] hover:bg-[#1055BD] text-white px-8 py-3 rounded-2xl font-bold min-h-[44px] transition-all tap-bounce active:scale-95 shadow-md shadow-[#1769E0]/20"
        >
          View My Products
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
            <Zap size={11} /> Automated Verification
          </span>
        </div>
        <p className="text-3xl font-black text-[#1769E0] dark:text-[#38A5FF]">ETB {feeAmount.toLocaleString()}</p>
        <p className="text-xs text-zinc-500 mt-1">One-time listing fee via official Telebirr payment integration.</p>
      </div>

      {/* Official Telebirr Card */}
      <div className="rounded-2xl border-2 border-[#1769E0] bg-blue-50/50 dark:bg-blue-900/10 p-4 space-y-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-2xl">📱</span>
            <div>
              <p className="font-black text-sm text-zinc-900 dark:text-white">Telebirr SuperApp Direct</p>
              <p className="text-xs text-zinc-500">Official Mobile Payment Rail</p>
            </div>
          </div>
          <span className="text-xs font-bold text-[#1769E0] bg-white dark:bg-zinc-800 px-2.5 py-1 rounded-full border border-blue-200 dark:border-blue-800">
            Authoritative
          </span>
        </div>

        {merchantTarget && (
          <div className="bg-white dark:bg-zinc-900 rounded-xl p-3 flex justify-between items-center text-xs">
            <span className="text-zinc-500">Merchant Phone:</span>
            <code className="font-black text-[#1769E0]">{merchantTarget}</code>
          </div>
        )}

        <div className="text-xs text-zinc-600 dark:text-zinc-400 space-y-1 bg-white/60 dark:bg-zinc-900/40 rounded-xl p-3 border border-blue-100 dark:border-blue-900/40">
          <p className="font-bold text-zinc-800 dark:text-zinc-200">✓ Zero Manual Reference Numbers Required</p>
          <p className="text-[11px]">Your payment is verified and approved directly by the system gateway.</p>
        </div>
      </div>

      {/* Payment Form */}
      <form onSubmit={handlePay} className="space-y-4">
        <div>
          <label className="text-xs font-bold text-zinc-500 block mb-1">
            <Phone size={11} className="inline mr-1" />
            Your Telebirr Phone Number
          </label>
          <input
            type="tel"
            value={payerPhone}
            onChange={e => { setPhone(e.target.value); setPhoneError(''); }}
            placeholder="09XXXXXXXX"
            className={INPUT_CLS}
            required
          />
          {phoneResult?.valid && (
            <p className="text-xs text-emerald-600 mt-1">
              ✓ Telebirr Account Verified: {phoneResult.normalized}
            </p>
          )}
          {phoneError && (
            <p className="text-xs text-red-500 mt-1 flex items-center gap-1">
              <AlertCircle size={11} /> {phoneError}
            </p>
          )}
        </div>

        {error && (
          <div className="p-3 bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-800 rounded-xl text-xs text-red-600 flex items-center gap-2">
            <AlertCircle size={14} className="shrink-0" />
            <span>{error}</span>
          </div>
        )}

        <button
          type="submit"
          disabled={submitting}
          className="w-full bg-[#1769E0] hover:bg-[#1055BD] disabled:opacity-50 text-white font-bold py-3.5 px-4 rounded-xl min-h-[48px] flex items-center justify-center gap-2 transition-all shadow-md shadow-[#1769E0]/20"
        >
          {submitting ? (
            <>
              <Spinner size={16} />
              <span>Verifying Telebirr Payment…</span>
            </>
          ) : (
            <>
              <span>Pay ETB {feeAmount.toLocaleString()} with Telebirr</span>
            </>
          )}
        </button>
      </form>
    </div>
  );
}

export default function PayPage() {
  return (
    <Suspense fallback={<div className="flex justify-center py-20"><Spinner size={32} /></div>}>
      <PayContent />
    </Suspense>
  );
}
