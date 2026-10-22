'use client';
import { useState, FormEvent } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { validateEthioPhone } from '@/lib/validations/common';
import { usePaymentSettings } from '@/hooks/usePaymentSettings';
import { PAYMENT_METHODS, getMerchantTarget } from '@/lib/constants/payment';
import { Crown, Star, Check, AlertCircle, Clock } from 'lucide-react';
import { INPUT_CLS_LIGHT as inputCls } from '@/lib/constants/styles';

const DURATIONS = [
  { label: '1 Week',   days: 7 },
  { label: '1 Month',  days: 30 },
  { label: '6 Months', days: 180 },
  { label: '1 Year',   days: 365 },
];

export default function BuyerUpgradePage() {
  const { user, buyerTier } = useAuth();
  const { settings } = usePaymentSettings();

  const availableTiers = (['gold', 'premium'] as const).filter(t =>
    buyerTier === 'gold' ? t === 'premium' : true
  );

  const [tier, setTier]         = useState<'gold' | 'premium'>(buyerTier === 'gold' ? 'premium' : 'gold');
  const [duration, setDuration] = useState(30);
  const [method, setMethod]     = useState<'telebirr' | 'safaricom' | 'cbe'>('telebirr');
  const [payerPhone, setPayerPhone] = useState('');
  const [phoneError, setPhoneError] = useState('');
  const [reference, setReference]   = useState('');
  const [loading, setLoading]   = useState(false);
  const [done, setDone]         = useState(false);
  const [error, setError]       = useState('');

  const goldPrice    = settings.buyer_gold_upgrade_price;
  const premiumPrice = settings.buyer_premium_upgrade_price;
  const basePrice    = tier === 'gold' ? goldPrice : premiumPrice;
  const price        = Math.round(basePrice * (duration / 30.0));
  const merchantNumber = getMerchantTarget(method, settings);

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!user) return;

    const pr = validateEthioPhone(payerPhone);
    if (!pr.valid) { setPhoneError(pr.error ?? 'Enter a valid Ethiopian phone number'); return; }
    setPhoneError('');
    if (!reference.trim()) { setError('Payment reference is required'); return; }

    setLoading(true); setError('');
    try {
      await api.buyerUpgradeTier({
        tier,
        reference: reference.trim(),
        payer_phone: pr.normalized,
        payment_method: method,
        amount: price,
        duration_days: duration,
      }, user.token);
      setDone(true);
    } catch (err: any) {
      setError(err.message ?? 'Something went wrong. Please try again.');
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

  if (done) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-4 max-w-sm mx-auto text-center">
        <div className="w-14 h-14 rounded-2xl bg-green-100 dark:bg-green-900/30 flex items-center justify-center">
          <Clock size={24} className="text-green-600" />
        </div>
        <div>
          <p className="font-black text-zinc-900 dark:text-white text-lg">Request Submitted</p>
          <p className="text-sm text-zinc-500 mt-1">
            Your upgrade request is pending admin review. Your tier will be updated once payment is verified.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-5 max-w-md">
      <div>
        <h1 className="text-xl font-black text-zinc-900 dark:text-white">Upgrade Tier</h1>
        <p className="text-xs text-zinc-500 mt-0.5">Current tier: <span className="font-bold capitalize">{buyerTier}</span></p>
      </div>

      {/* Tier selector */}
      <div className={`grid gap-3 ${availableTiers.length === 1 ? 'grid-cols-1' : 'grid-cols-2'}`}>
        {availableTiers.map(t => (
          <button key={t} onClick={() => setTier(t)}
            className={`flex flex-col items-center gap-2 py-4 px-3 rounded-2xl border-2 transition-all ${
              tier === t
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
        <p className="text-[10px] font-black uppercase tracking-widest text-zinc-400 mb-2">Duration</p>
        <div className="grid grid-cols-4 gap-2">
          {DURATIONS.map(d => (
            <button key={d.days} type="button" onClick={() => setDuration(d.days)}
              className={`py-2 px-1 rounded-xl text-[10px] font-bold border transition-all ${
                duration === d.days
                  ? 'border-orange-500 bg-orange-500/10 text-orange-600'
                  : 'border-zinc-200 dark:border-zinc-700 text-zinc-500 hover:border-orange-400'
              }`}>{d.label}</button>
          ))}
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">

        {/* Payment method */}
        <div>
          <p className="text-[10px] font-black uppercase tracking-widest text-zinc-400 mb-2">Payment Method</p>
          <div className="grid grid-cols-3 gap-2">
            {PAYMENT_METHODS.map(m => (
              <button key={m.value} type="button" onClick={() => setMethod(m.value)}
                className={`py-2 rounded-xl text-xs font-bold border transition-all ${
                  method === m.value
                    ? 'border-orange-500 bg-orange-500/10 text-orange-600'
                    : 'border-zinc-200 dark:border-zinc-700 text-zinc-500'
                }`}>{m.label}</button>
            ))}
          </div>
        </div>

        {/* Merchant number */}
        <div className="bg-zinc-50 dark:bg-zinc-800 rounded-xl px-4 py-3">
          <p className="text-[10px] text-zinc-400 uppercase tracking-widest">Send ETB {price.toLocaleString()} to</p>
          <p className="text-lg font-black text-zinc-900 dark:text-white mt-0.5">{merchantNumber}</p>
          <p className="text-[10px] text-zinc-400 mt-0.5">via {PAYMENT_METHODS.find(m => m.value === method)?.label}</p>
        </div>

        {/* Payer phone */}
        <div>
          <label className="text-[10px] font-black uppercase tracking-widest text-zinc-400 block mb-1.5">Your Payment Phone *</label>
          <input
            type="tel"
            value={payerPhone}
            onChange={e => { setPayerPhone(e.target.value); setPhoneError(''); }}
            placeholder="09XXXXXXXX or 07XXXXXXXX"
            className={inputCls}
          />
          {phoneError && (
            <p className="flex items-center gap-1 text-[11px] text-red-500 mt-1">
              <AlertCircle size={11} /> {phoneError}
            </p>
          )}
          <p className="text-[10px] text-zinc-400 mt-1">Ethio Telecom: 09XXXXXXXX · Safaricom: 07XXXXXXXX</p>
        </div>

        {/* Reference */}
        <div>
          <label className="text-[10px] font-black uppercase tracking-widest text-zinc-400 block mb-1.5">Transaction Reference *</label>
          <input
            value={reference}
            onChange={e => setReference(e.target.value)}
            placeholder="e.g. TXN123456789"
            className={inputCls}
          />
          <p className="text-[10px] text-zinc-400 mt-1">The reference number from your payment receipt</p>
        </div>

        {error && (
          <div className="flex items-center gap-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 text-xs px-3 py-2.5 rounded-xl">
            <AlertCircle size={13} /> {error}
          </div>
        )}

        <button type="submit" disabled={loading}
          className="w-full flex items-center justify-center gap-2 bg-orange-600 hover:bg-orange-700 text-white py-3 rounded-xl font-bold text-sm transition-colors disabled:opacity-50">
          {loading ? <Spinner /> : <><Check size={15} /> Submit Upgrade Request</>}
        </button>

        <p className="text-[10px] text-zinc-400 text-center">
          Your request will be reviewed by an admin. Tier is applied after payment verification.
        </p>
      </form>
    </div>
  );
}

