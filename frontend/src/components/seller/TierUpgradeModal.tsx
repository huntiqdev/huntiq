'use client';
import { useState, FormEvent } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import {
  X, Crown, Star, Phone, CreditCard, CheckCircle, AlertCircle,
  Loader2, Images, PhoneCall, TrendingDown, ShoppingCart,
  DollarSign, BarChart2, Zap, Shield, Sparkles, Copy,
} from 'lucide-react';
import { validateEthioPhone } from '@/lib/validations/common';
import { usePaymentSettings } from '@/hooks/usePaymentSettings';
import { PAYMENT_METHODS, getMerchantTarget } from '@/lib/constants/payment';
import { INPUT_CLS_DARK as inputCls } from '@/lib/constants/styles';

interface Props {
  currentTier: string;
  initialTier?: 'gold' | 'premium';
  onClose: () => void;
  onSuccess: () => void;
}

interface Perk { icon: any; text: string; highlight?: boolean }

const TIER_DATA: Record<string, {
  gradient: string; border: string; glow: string;
  iconColor: string; textColor: string; badgeBg: string;
  perks: Perk[];
}> = {
  gold: {
    gradient: 'from-yellow-400 via-amber-400 to-orange-400',
    border: 'border-yellow-400/40',
    glow: 'shadow-yellow-500/20',
    iconColor: 'text-yellow-400',
    textColor: 'text-yellow-300',
    badgeBg: 'bg-yellow-500/20',
    perks: [
      { icon: Images,       text: 'Up to 7 product images',          highlight: true },
      { icon: PhoneCall,    text: '3 additional contact phones',      highlight: true },
      { icon: ShoppingCart, text: 'Order Requests page unlocked',     highlight: true },
      { icon: DollarSign,   text: 'Sales & Transactions page unlocked', highlight: true },
      { icon: TrendingDown, text: 'Lower listing fee — 1.5% startup' },
      { icon: Shield,       text: 'Priority listing visibility' },
    ],
  },
  premium: {
    gradient: 'from-purple-500 via-violet-500 to-indigo-500',
    border: 'border-purple-400/40',
    glow: 'shadow-purple-500/20',
    iconColor: 'text-purple-400',
    textColor: 'text-purple-300',
    badgeBg: 'bg-purple-500/20',
    perks: [
      { icon: Images,       text: 'Unlimited product images',         highlight: true },
      { icon: PhoneCall,    text: 'Unlimited contact phones',         highlight: true },
      { icon: ShoppingCart, text: 'Order Requests page unlocked',     highlight: true },
      { icon: DollarSign,   text: 'Sales & Transactions page unlocked', highlight: true },
      { icon: BarChart2,    text: 'Analytics dashboard unlocked',     highlight: true },
      { icon: TrendingDown, text: 'Lowest listing fee — 1.2% startup' },
      { icon: Sparkles,     text: 'Top listing priority + dedicated support' },
    ],
  },
};

export function TierUpgradeModal({ currentTier, initialTier, onClose, onSuccess }: Props) {
  const { user } = useAuth();
  const { settings } = usePaymentSettings();

  const availableTiers = currentTier === 'normal'
    ? (['gold', 'premium'] as const)
    : (['premium'] as const);

  const [selectedTier, setSelectedTier] = useState<'gold' | 'premium'>(
    initialTier ?? availableTiers[0]
  );
  const [step, setStep]       = useState<'pick' | 'pay' | 'done'>('pick');
  const [method, setMethod]   = useState('telebirr');
  const [phone, setPhone]     = useState('');
  const [reference, setRef]   = useState('');
  const [notes, setNotes]     = useState('');
  const [saving, setSaving]   = useState(false);
  const [error, setError]     = useState('');
  const [phoneErr, setPhoneErr] = useState('');
  const [copied, setCopied]   = useState(false);

  const tierPrices = { gold: settings.gold_upgrade_price, premium: settings.premium_upgrade_price };
  const merchantTarget = getMerchantTarget(method as any, settings);

  const copyTarget = () => {
    navigator.clipboard.writeText(merchantTarget).then(() => {
      setCopied(true); setTimeout(() => setCopied(false), 2000);
    });
  };

  const price = tierPrices[selectedTier] ?? 500;
  const td    = TIER_DATA[selectedTier];

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!user) return;
    const pr = validateEthioPhone(phone);
    if (!pr.valid) { setPhoneErr(pr.error ?? 'Invalid phone number'); return; }
    setPhoneErr('');
    if (!reference.trim()) { setError('Payment reference is required'); return; }
    setSaving(true); setError('');
    try {
      await api.sellerRequestUpgrade({
        requested_tier: selectedTier,
        payment_method: method,
        payer_phone: pr.normalized,
        reference: reference.trim(),
        amount: price,
        notes: notes.trim() || undefined,
      }, user.token);
      setStep('done');
    } catch (err: any) {
      setError(err.message ?? 'Failed to submit request');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4">
      <div className="absolute inset-0 bg-black/75 backdrop-blur-md" onClick={onClose} />

      <div className="relative z-10 w-full sm:max-w-lg bg-[#0d1117] border border-white/10 rounded-t-3xl sm:rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[95vh] sm:max-h-[90vh]">

        {/* Gradient header bar */}
        <div className={`h-1 w-full bg-gradient-to-r ${td.gradient}`} />

        {/* Header */}
        <div className="flex items-center justify-between px-5 pt-5 pb-3 shrink-0">
          <div className="flex items-center gap-2.5">
            <div className={`w-8 h-8 rounded-xl bg-gradient-to-br ${td.gradient} flex items-center justify-center shadow-lg ${td.glow}`}>
              {selectedTier === 'gold' ? <Star size={16} className="text-white" /> : <Crown size={16} className="text-white" />}
            </div>
            <div>
              <h2 className="text-base font-black text-white leading-none">Upgrade to {selectedTier === 'gold' ? 'Gold' : 'Premium'}</h2>
              <p className="text-[10px] text-white/40 mt-0.5">Admin activates within 24 hours</p>
            </div>
          </div>
          <button type="button" onClick={onClose}
            className="w-8 h-8 rounded-xl bg-white/5 hover:bg-white/10 flex items-center justify-center text-white/40 hover:text-white transition-all">
            <X size={15} />
          </button>
        </div>

        {/* Step indicator */}
        {step !== 'done' && (
          <div className="flex items-center gap-2 px-5 pb-3 shrink-0">
            {(['pick', 'pay'] as const).map((s, i) => (
              <div key={s} className="flex items-center gap-2">
                <div className={`w-6 h-6 rounded-full flex items-center justify-center text-[10px] font-black transition-all ${
                  step === s ? `bg-gradient-to-br ${td.gradient} text-white shadow-md` :
                  (step === 'pay' && s === 'pick') ? 'bg-green-500 text-white' :
                  'bg-white/10 text-white/30'
                }`}>
                  {step === 'pay' && s === 'pick' ? <CheckCircle size={12} /> : i + 1}
                </div>
                <span className={`text-[10px] font-bold ${step === s ? 'text-white' : 'text-white/30'}`}>
                  {s === 'pick' ? 'Choose Tier' : 'Payment'}
                </span>
                {i === 0 && <div className="w-6 h-px bg-white/10" />}
              </div>
            ))}
          </div>
        )}

        {/* ── STEP: DONE ── */}
        {step === 'done' && (
          <div className="px-5 pb-8 pt-4 text-center space-y-5 flex-1 flex flex-col items-center justify-center">
            <div className={`w-20 h-20 rounded-full bg-gradient-to-br ${td.gradient} flex items-center justify-center shadow-2xl ${td.glow}`}>
              <CheckCircle size={40} className="text-white" strokeWidth={2} />
            </div>
            <div>
              <h3 className="text-xl font-black text-white">Request Submitted!</h3>
              <p className="text-sm text-white/50 mt-2 leading-relaxed">
                Your upgrade to <span className={`font-black capitalize ${td.textColor}`}>{selectedTier}</span> is under review.
                Admin will verify your payment and activate within 24 hours.
              </p>
            </div>
            <div className={`w-full rounded-2xl border ${td.border} ${td.badgeBg} p-4 text-left space-y-2`}>
              {td.perks.filter(p => p.highlight).map((p, i) => (
                <div key={i} className="flex items-center gap-2.5">
                  <p.icon size={13} className={td.iconColor} />
                  <span className="text-xs text-white/70">{p.text}</span>
                </div>
              ))}
            </div>
            <button onClick={() => { onSuccess(); onClose(); }}
              className={`w-full py-3.5 rounded-2xl bg-gradient-to-r ${td.gradient} text-white font-black text-sm shadow-lg ${td.glow} hover:opacity-90 transition-all`}>
              Done
            </button>
          </div>
        )}

        {/* ── STEP: PICK TIER ── */}
        {step === 'pick' && (
          <div className="flex flex-col overflow-hidden flex-1">
            <div className="px-5 pb-4 space-y-4 overflow-y-auto flex-1">

              {/* Tier cards — only show selector if normal (both options) */}
              {availableTiers.length > 1 ? (
                <div className="grid grid-cols-2 gap-3">
                  {availableTiers.map(t => {
                    const d = TIER_DATA[t];
                    const active = selectedTier === t;
                    return (
                      <button key={t} type="button" onClick={() => setSelectedTier(t)}
                        className={`relative flex flex-col items-center gap-2 p-4 rounded-2xl border-2 transition-all text-left ${
                          active ? `${d.border} ${d.badgeBg} shadow-lg ${d.glow}` : 'border-white/5 bg-white/[0.02] hover:border-white/10'
                        }`}>
                        {active && (
                          <div className={`absolute top-2 right-2 w-4 h-4 rounded-full bg-gradient-to-br ${d.gradient} flex items-center justify-center`}>
                            <CheckCircle size={10} className="text-white" />
                          </div>
                        )}
                        <div className={`w-10 h-10 rounded-xl bg-gradient-to-br ${d.gradient} flex items-center justify-center shadow-md`}>
                          {t === 'gold' ? <Star size={18} className="text-white" /> : <Crown size={18} className="text-white" />}
                        </div>
                        <p className={`text-sm font-black capitalize ${active ? 'text-white' : 'text-white/50'}`}>{t}</p>
                        <p className={`text-xs font-black ${active ? d.textColor : 'text-white/30'}`}>ETB {tierPrices[t].toLocaleString()}</p>
                      </button>
                    );
                  })}
                </div>
              ) : (
                /* Single tier — premium only */
                <div className={`flex items-center gap-3 p-4 rounded-2xl border ${td.border} ${td.badgeBg}`}>
                  <div className={`w-10 h-10 rounded-xl bg-gradient-to-br ${td.gradient} flex items-center justify-center shadow-md shrink-0`}>
                    <Crown size={18} className="text-white" />
                  </div>
                  <div>
                    <p className="text-sm font-black text-white">Premium Tier</p>
                    <p className={`text-xs font-black ${td.textColor}`}>ETB {price.toLocaleString()} one-time</p>
                  </div>
                </div>
              )}

              {/* Perks list */}
              <div className="space-y-1.5">
                <p className="text-[10px] font-black uppercase tracking-widest text-white/30 px-0.5">What you unlock</p>
                {td.perks.map((p, i) => (
                  <div key={i} className={`flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all ${
                    p.highlight
                      ? `${td.badgeBg} border ${td.border}`
                      : 'bg-white/[0.02]'
                  }`}>
                    <div className={`w-7 h-7 rounded-lg flex items-center justify-center shrink-0 ${
                      p.highlight ? `bg-gradient-to-br ${td.gradient}` : 'bg-white/5'
                    }`}>
                      <p.icon size={13} className={p.highlight ? 'text-white' : 'text-white/30'} />
                    </div>
                    <span className={`text-xs font-semibold ${p.highlight ? 'text-white' : 'text-white/40'}`}>{p.text}</span>
                    {p.highlight && <Zap size={10} className={`ml-auto shrink-0 ${td.iconColor}`} />}
                  </div>
                ))}
              </div>
            </div>

            <div className="px-5 py-4 border-t border-white/5 flex gap-2.5 shrink-0">
              <button type="button" onClick={onClose}
                className="flex-1 py-3 rounded-2xl bg-white/5 text-white/50 text-sm font-bold hover:bg-white/10 transition-all">
                Cancel
              </button>
              <button type="button" onClick={() => setStep('pay')}
                className={`flex-2 flex-1 py-3 rounded-2xl bg-gradient-to-r ${td.gradient} text-white text-sm font-black shadow-lg ${td.glow} hover:opacity-90 transition-all flex items-center justify-center gap-2`}>
                Continue <Zap size={13} />
              </button>
            </div>
          </div>
        )}

        {/* ── STEP: PAY ── */}
        {step === 'pay' && (
          <form onSubmit={handleSubmit} className="flex flex-col overflow-hidden flex-1">
            <div className="px-5 pb-4 space-y-4 overflow-y-auto flex-1">

              {/* Price summary */}
              <div className={`flex items-center justify-between p-4 rounded-2xl border ${td.border} ${td.badgeBg}`}>
                <div>
                  <p className="text-[10px] text-white/40 uppercase tracking-widest">Upgrade fee</p>
                  <p className="text-2xl font-black text-white mt-0.5">ETB {price.toLocaleString()}</p>
                </div>
                <div className={`w-12 h-12 rounded-2xl bg-gradient-to-br ${td.gradient} flex items-center justify-center shadow-lg`}>
                  {selectedTier === 'gold' ? <Star size={22} className="text-white" /> : <Crown size={22} className="text-white" />}
                </div>
              </div>

              {/* Payment method */}
              <div>
                <p className="text-[10px] font-black uppercase tracking-widest text-white/30 mb-2">Payment Method</p>
                <div className="grid grid-cols-3 gap-2">
                  {PAYMENT_METHODS.map(m => (
                    <button key={m.value} type="button" onClick={() => setMethod(m.value)}
                      className={`flex flex-col items-center gap-1.5 py-3 px-2 rounded-xl border transition-all ${
                        method === m.value
                          ? `border-orange-500/50 bg-orange-500/10 text-orange-400`
                          : 'border-white/5 bg-white/[0.02] text-white/40 hover:border-white/10 hover:text-white/60'
                      }`}>
                      <span className="text-lg leading-none">{m.emoji}</span>
                      <span className="text-[10px] font-black">{m.label}</span>
                    </button>
                  ))}
                </div>
              </div>

              {/* Send-to instructions */}
              <div className="bg-white/[0.04] border border-white/8 rounded-xl p-3 space-y-1.5">
                <p className="text-[10px] font-black uppercase tracking-widest text-white/30">
                  {method === 'cbe' ? 'Send to account' : 'Send to phone'}
                </p>
                <div className="flex items-center gap-2">
                  <code className={`text-sm font-black ${td.textColor}`}>{merchantTarget}</code>
                  <button type="button" onClick={copyTarget}
                    className="text-white/30 hover:text-white/70 transition-colors">
                    <Copy size={12} />
                  </button>
                  {copied && <span className="text-[10px] text-green-400">Copied!</span>}
                </div>
                <p className="text-[10px] text-white/30">
                  Amount: <strong className="text-white">ETB {price.toLocaleString()}</strong> via {PAYMENT_METHODS.find(m2 => m2.value === method)?.label}
                </p>
              </div>

              {/* Phone */}
              <div>
                <label className="text-[10px] font-black uppercase tracking-widest text-white/30 block mb-1.5">
                  Phone Used for Payment
                </label>
                <div className="relative">
                  <Phone size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none" />
                  <input type="tel" value={phone} onChange={e => { setPhone(e.target.value); setPhoneErr(''); }}
                    placeholder="09XXXXXXXX or 07XXXXXXXX" className={inputCls + ' pl-10'} />
                </div>
                {phoneErr && (
                  <p className="text-[10px] text-red-400 mt-1 flex items-center gap-1">
                    <AlertCircle size={9} /> {phoneErr}
                  </p>
                )}
              </div>

              {/* Reference */}
              <div>
                <label className="text-[10px] font-black uppercase tracking-widest text-white/30 block mb-1.5">
                  Transaction Reference / Receipt No.
                </label>
                <div className="relative">
                  <CreditCard size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none" />
                  <input type="text" value={reference} onChange={e => setRef(e.target.value)}
                    placeholder="e.g. TXN123456789" className={inputCls + ' pl-10'} />
                </div>
                <p className="text-[10px] text-white/25 mt-1.5 leading-relaxed">
                  Send ETB {price.toLocaleString()} to our account first, then paste the reference number here.
                </p>
              </div>

              {/* Notes */}
              <div>
                <label className="text-[10px] font-black uppercase tracking-widest text-white/30 block mb-1.5">
                  Notes <span className="normal-case font-normal text-white/20">(optional)</span>
                </label>
                <textarea value={notes} onChange={e => setNotes(e.target.value)}
                  placeholder="Any additional info for the admin…" rows={2}
                  className={inputCls + ' resize-none'} />
              </div>

              {error && (
                <div className="flex items-start gap-2.5 bg-red-500/10 border border-red-500/20 rounded-xl px-3.5 py-3">
                  <AlertCircle size={14} className="text-red-400 shrink-0 mt-0.5" />
                  <p className="text-xs text-red-400">{error}</p>
                </div>
              )}
            </div>

            <div className="px-5 py-4 border-t border-white/5 flex gap-2.5 shrink-0">
              <button type="button" onClick={() => { setStep('pick'); setError(''); }}
                className="flex-1 py-3 rounded-2xl bg-white/5 text-white/50 text-sm font-bold hover:bg-white/10 transition-all">
                ← Back
              </button>
              <button type="submit" disabled={saving}
                className={`flex-1 py-3 rounded-2xl bg-gradient-to-r ${td.gradient} text-white text-sm font-black shadow-lg ${td.glow} hover:opacity-90 transition-all disabled:opacity-40 flex items-center justify-center gap-2`}>
                {saving ? <><Loader2 size={14} className="animate-spin" /> Submitting…</> : 'Submit Request'}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
