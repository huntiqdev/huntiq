'use client';
import { useState, FormEvent } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { validateEthioPhone } from '@/lib/validations/common';
import { usePaymentSettings } from '@/hooks/usePaymentSettings';
import { getMerchantTarget } from '@/lib/constants/payment';
import { Crown, Star, CheckCircle, AlertCircle, Phone, ShieldCheck, Zap } from 'lucide-react';
import { INPUT_CLS_LIGHT as inputCls } from '@/lib/constants/styles';

const DURATIONS = [
  { label: '1 Week', days: 7 },
  { label: '1 Month', days: 30 },
  { label: '6 Months', days: 180 },
  { label: '1 Year', days: 365 },
];

export default function BuyerUpgradePage() {
  const { user, buyerTier } = useAuth();
  const { settings } = usePaymentSettings();

  const availableTiers = (['gold', 'premium'] as const).filter(t =>
    buyerTier === 'gold' ? t === 'premium' : true
  );

  const [tier, setTier] = useState<'gold' | 'premium'>(buyerTier === 'gold' ? 'premium' : 'gold');
  const [duration, setDuration] = useState(30);
  const [payerPhone, setPayerPhone] = useState(user?.phone ?? '');
  const [phoneError, setPhoneError] = useState('');
  const [loading, setLoading] = useState(false);
  const [done, setDone] = useState(false);
  const [error, setError] = useState('');

  const goldPrice = settings.buyer_gold_upgrade_price;
  const premiumPrice = settings.buyer_premium_upgrade_price;
  const basePrice = tier === 'gold' ? goldPrice : premiumPrice;
  const price = Math.round(basePrice * (duration / 30.0));
  const merchantNumber = getMerchantTarget('telebirr', settings);

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!user) return;

    const pr = validateEthioPhone(payerPhone);
    if (!pr.valid) {
      setPhoneError(pr.error ?? 'Enter a valid Ethiopian Telebirr phone number');
      return;
    }
    setPhoneError('');

    setLoading(true);
    setError('');
    try {
      // 1. Create Telebirr payment order
      const res = await api.createTelebirrOrder({
        amount: price,
        order_type: 'buyer_tier_upgrade',
        title: `Buyer Tier Upgrade to ${tier.toUpperCase()}`,
      }, user.token);

      // 2. Register upgrade tier request with Telebirr out_trade_no
      await api.buyerUpgradeTier({
        tier,
        reference: res?.order?.out_trade_no || `TB-UPG-${Date.now()}`,
        payer_phone: pr.normalized,
        payment_method: 'telebirr',
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
      <div className="flex flex-col items-center justify-center min-h-[50vh] gap-4 max-w-sm mx-auto text-center px-4">
        <div className="w-16 h-16 rounded-full bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center text-emerald-600 shadow-xl shadow-emerald-500/20">
          <CheckCircle size={36} strokeWidth={2.5} />
        </div>
        <div>
          <p className="font-black text-zinc-900 dark:text-white text-xl">Tier Upgrade Activated!</p>
          <p className="text-sm text-zinc-500 mt-2">
            Your account is now upgraded to <strong className="capitalize text-zinc-800 dark:text-zinc-200">{tier}</strong> tier via Telebirr automated verification.
          </p>
          <div className="inline-flex items-center gap-1.5 mt-3 px-3 py-1 bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-800 rounded-full text-xs font-bold text-emerald-600 dark:text-emerald-400">
            <ShieldCheck size={14} /> System Auto-Approved
          </div>
        </div>
        <button
          onClick={() => window.location.href = '/buyer/dashboard'}
          className="bg-[#1769E0] text-white font-bold px-6 py-2.5 rounded-xl hover:bg-[#1055BD] transition-all"
        >
          Go to Dashboard
        </button>
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
        <p className="text-[10px] font-black uppercase tracking-widest text-zinc-400 mb-2">Duration</p>
        <div className="grid grid-cols-4 gap-2">
          {DURATIONS.map(d => (
            <button key={d.days} type="button" onClick={() => setDuration(d.days)}
              className={`py-2 px-1 rounded-xl text-[10px] font-bold border transition-all ${duration === d.days
                ? 'border-[#1769E0] bg-[#1769E0]/10 text-[#1769E0]'
                : 'border-zinc-200 dark:border-zinc-700 text-zinc-500 hover:border-[#1769E0]'
                }`}>{d.label}</button>
          ))}
        </div>
      </div>

      {/* Telebirr Integration Box */}
      <div className="rounded-2xl border-2 border-[#1769E0] bg-blue-50/50 dark:bg-blue-900/10 p-4 space-y-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-2xl">📱</span>
            <div>
              <p className="font-black text-sm text-zinc-900 dark:text-white">Telebirr SuperApp Direct</p>
              <p className="text-xs text-zinc-500">Official Mobile Payment Gateway</p>
            </div>
          </div>
          <span className="flex items-center gap-1 text-[11px] font-bold text-emerald-600 dark:text-emerald-400 bg-white dark:bg-zinc-800 px-2 py-0.5 rounded-full border border-blue-200 dark:border-blue-800">
            <Zap size={10} /> Auto-Approved
          </span>
        </div>

        {merchantNumber && (
          <div className="bg-white dark:bg-zinc-900 rounded-xl p-2.5 flex justify-between items-center text-xs">
            <span className="text-zinc-500">Merchant:</span>
            <code className="font-black text-[#1769E0]">{merchantNumber}</code>
          </div>
        )}
      </div>

      <form onSubmit={handleSubmit} className="space-y-4 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
        <div>
          <label className="text-xs font-bold text-zinc-500 block mb-1">
            <Phone size={11} className="inline mr-1" />
            Your Telebirr Phone Number
          </label>
          <input
            type="tel"
            value={payerPhone}
            onChange={e => { setPayerPhone(e.target.value); setPhoneError(''); }}
            placeholder="09XXXXXXXX"
            className={inputCls}
            required
          />
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
          disabled={loading}
          className="w-full bg-[#1769E0] hover:bg-[#1055BD] disabled:opacity-50 text-white font-bold py-3.5 px-4 rounded-xl min-h-[48px] flex items-center justify-center gap-2 transition-all shadow-md shadow-[#1769E0]/20"
        >
          {loading ? (
            <>
              <Spinner size={16} />
              <span>Verifying Telebirr Payment…</span>
            </>
          ) : (
            <span>Pay ETB {price.toLocaleString()} with Telebirr</span>
          )}
        </button>
      </form>
    </div>
  );
}
