'use client';
import { useState, useEffect, useRef, FormEvent } from 'react';
import { api } from '@/lib/api';
import { Product } from '@/types';
import { useLocale } from '@/hooks/useLocale';
import { useAuth } from '@/context/AuthContext';
import { validateEthioPhone, validateOrderRequest, resolveFieldError } from '@/lib/validations';
import { X, Check, Loader2, Phone, MessageCircle, AlertCircle, Star, Crown, Lock } from 'lucide-react';
import Link from 'next/link';
import { TIER_RANK } from '@/lib/constants/tiers';
import { INPUT_CLS as inputClass } from '@/lib/constants/styles';

interface Props {
  product: Product;
  onClose: () => void;
}

type Step = 'form' | 'success';

export function BuyModal({ product, onClose }: Props) {
  const { t } = useLocale();
  const u = t.ui;
  const { user, isBuyer, buyerTier } = useAuth();
  const isRental = product.listing_type === 'rent';

  // For logged-in buyers, pre-fill their registered phone (auto-call)
  const [step, setStep] = useState<Step>('form');
  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');
  const [phoneError, setPhoneError] = useState('');
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // Fetch buyer profile phone once — skip if already loaded
  const [buyerPhone, setBuyerPhone] = useState<string | null>(null);
  const [buyerName, setBuyerName] = useState<string | null>(null);
  const fetchedRef = useRef(false);

  useEffect(() => {
    if (!isBuyer || !user || fetchedRef.current) return;
    fetchedRef.current = true;
    api.getBuyerProfile(user.token).then(p => {
      if (p?.phone) setBuyerPhone(p.phone);
      if (p?.name) setBuyerName(p.name);
    }).catch(() => {});
  }, [isBuyer, user?.token]); // eslint-disable-line react-hooks/exhaustive-deps

  const handlePhoneChange = (v: string) => {
    setPhone(v);
    if (phoneError) setPhoneError('');
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');

    // Use buyer's registered phone if logged in, otherwise validate entered phone
    const effectivePhone = isBuyer && buyerPhone ? buyerPhone : phone;
    const effectiveName = isBuyer && buyerName ? buyerName : name.trim();

    const fieldErrors = validateOrderRequest({
      name,
      phone,
      phonePreFilled: !!(isBuyer && buyerPhone),
    });
    if (fieldErrors.phone) { setPhoneError(resolveFieldError(fieldErrors.phone, t) ?? ''); return; }
    if (fieldErrors.name) { setError(resolveFieldError(fieldErrors.name, t) ?? ''); return; }

    const phoneResult = validateEthioPhone(effectivePhone);
    const normalizedPhone = phoneResult.valid ? phoneResult.normalized : effectivePhone;

    setLoading(true);
    try {
      await api.createOrderRequest({
        product_id: product.id,
        buyer_name: effectiveName || name.trim() || (isRental ? 'Renter' : 'Buyer'),
        buyer_phone: normalizedPhone,
        message: message.trim() || (isRental ? 'Renter interested in property — please contact me.' : 'Buyer is interested — please contact me.'),
      }, user?.token);
      setStep('success');
    } catch (err: any) {
      setError(err.message ?? u.purchase_error);
    } finally {
      setLoading(false);
    }
  };

  // Tier-gated check — product's required_buyer_tier vs buyer's tier
  const requiredTier = product.required_buyer_tier ?? 'normal';
  const buyerRank = (isBuyer && buyerTier) ? (TIER_RANK[buyerTier] ?? 0) : 0;
  const requiredRank = TIER_RANK[requiredTier] ?? 0;
  const tierLocked = buyerRank < requiredRank;

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-3 sm:p-4">
      <div className="absolute inset-0" onClick={onClose} />
      <div className="relative z-10 bg-white dark:bg-zinc-900 rounded-3xl shadow-2xl w-full max-w-sm max-h-[92vh] overflow-y-auto p-5 sm:p-6">
        {/* Header */}
        <div className="flex items-center justify-between mb-5">
          <div className="flex items-center gap-2">
            <div className={`w-8 h-8 rounded-xl flex items-center justify-center ${isRental ? 'bg-teal-100 dark:bg-teal-900/30' : 'bg-orange-100 dark:bg-orange-900/30'}`}>
              <Phone size={16} className={isRental ? 'text-teal-600' : 'text-orange-600'} />
            </div>
            <h2 className="font-black text-base">
              {step === 'form' 
                ? (isRental ? (u.contact_landlord ?? 'Contact Landlord') : (u.contact_seller ?? 'Contact Seller')) 
                : (u.request_sent ?? 'Request Sent')}
            </h2>
          </div>
          <button type="button" onClick={onClose} className="p-2 rounded-full hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors">
            <X size={18} />
          </button>
        </div>

        {/* Tier-locked warning */}
        {tierLocked && (
          <div className="mb-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-2xl p-4 text-center space-y-2">
            <div className="flex items-center justify-center gap-2">
              <Lock size={16} className="text-yellow-600" />
              <p className="font-black text-sm text-yellow-700 dark:text-yellow-400">
                {requiredTier === 'gold' ? 'Gold' : 'Premium'} Buyer Required
              </p>
            </div>
            <p className="text-xs text-yellow-600 dark:text-yellow-500">
              This {isRental ? 'landlord' : 'seller'}'s listings require a{' '}
              {requiredTier === 'gold'
                ? <><Star size={10} className="inline" /> Gold</>
                : <><Crown size={10} className="inline" /> Premium</>}{' '}
              buyer tier to view full details.
            </p>
            <Link href="/buyer/login" onClick={onClose}
              className="inline-flex items-center gap-1.5 text-xs font-black px-3 py-1.5 rounded-xl bg-yellow-500 hover:bg-yellow-600 text-white transition-colors">
              {requiredTier === 'gold' ? <Star size={11} /> : <Crown size={11} />}
              Upgrade Now
            </Link>
          </div>
        )}

        {step === 'form' && !tierLocked ? (
          <>
            {/* Product strip */}
            <div className="bg-zinc-50 dark:bg-zinc-800 rounded-2xl px-4 py-3 mb-4 flex items-center gap-3">
              <div className="flex-1 min-w-0">
                <p className="font-bold text-sm truncate">{product.title}</p>
                <p className={`font-black text-sm ${isRental ? 'text-teal-600' : 'text-orange-600'}`}>
                  {isRental 
                    ? `ETB ${product.rental_price ? product.rental_price.toLocaleString() : product.price.toLocaleString()} / month`
                    : `ETB ${product.price.toLocaleString()}`}
                </p>
              </div>
              <div className="text-[10px] text-zinc-400 text-right leading-tight">
                <p className="font-black text-zinc-500">{u.buy_modal_agreement ?? 'Agreement'}</p>
                <p>{isRental ? (u.buy_modal_renter_landlord ?? 'Renter ↔ Landlord') : (u.buy_modal_buyer_seller ?? 'Buyer ↔ Seller')}</p>
              </div>
            </div>

            {/* Auto-call info for logged-in buyers */}
            {isBuyer && buyerPhone ? (
              <div className="flex items-start gap-2 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-xl px-3 py-2.5 mb-4">
                <Phone size={13} className="text-green-500 shrink-0 mt-0.5" />
                <p className="text-[11px] text-green-700 dark:text-green-300 leading-relaxed">
                  The {isRental ? 'landlord' : 'seller'} will call you on your registered number:{' '}
                  <span className="font-black">{buyerPhone}</span>
                </p>
              </div>
            ) : (
              <div className="flex items-start gap-2 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl px-3 py-2.5 mb-4">
                <MessageCircle size={13} className="text-blue-500 shrink-0 mt-0.5" />
                <p className="text-[11px] text-blue-700 dark:text-blue-300 leading-relaxed">
                  {isRental 
                    ? (u.contact_landlord_call ?? 'The landlord will call you directly on the number you provide to arrange the rental agreement.') 
                    : (u.buy_modal_seller_call ?? 'The seller will call you directly on the number you provide to arrange the purchase.')}
                </p>
              </div>
            )}

            <form onSubmit={handleSubmit} className="flex flex-col gap-3">
              {/* Full name — only if not logged in as buyer */}
              {(!isBuyer || !buyerName) && (
                <div>
                  <label className="text-[10px] font-black uppercase tracking-widest text-zinc-500 block mb-1">
                    {u.buy_modal_your_name ?? 'Your Full Name *'}
                  </label>
                  <input
                    required
                    type="text"
                    placeholder="e.g. Abebe Kebede"
                    value={name}
                    onChange={e => setName(e.target.value)}
                    className={inputClass}
                  />
                </div>
              )}

              {/* Phone — only if not logged in as buyer with phone */}
              {(!isBuyer || !buyerPhone) && (
                <div>
                  <label className="text-[10px] font-black uppercase tracking-widest text-zinc-500 block mb-1">
                    {u.buy_modal_your_phone ?? 'Your Phone Number *'}{' '}
                    <span className="normal-case font-normal text-zinc-400">
                      ({isRental ? 'landlord' : 'seller'} will call you)
                    </span>
                  </label>
                  <div className="relative">
                    <Phone size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none" />
                    <input
                      required
                      type="tel"
                      placeholder="09XXXXXXXX or 07XXXXXXXX"
                      value={phone}
                      onChange={e => handlePhoneChange(e.target.value)}
                      className={`${inputClass} pl-9 ${phoneError ? 'border-red-400 focus:ring-red-400' : ''}`}
                    />
                  </div>
                  {phoneError && (
                    <div className="flex items-center gap-1.5 mt-1.5">
                      <AlertCircle size={11} className="text-red-500 shrink-0" />
                      <p className="text-[11px] text-red-500">{phoneError}</p>
                    </div>
                  )}
                  <p className="text-[10px] text-zinc-400 mt-1">
                    {u.payer_phone_hint ?? 'Ethio Telecom: 09XXXXXXXX · Safaricom: 07XXXXXXXX · or +251 prefix'}
                  </p>
                </div>
              )}

              {/* Optional message */}
              <div>
                <label className="text-[10px] font-black uppercase tracking-widest text-zinc-500 block mb-1">
                  {u.optional_message ?? 'Message'}{' '}
                  <span className="normal-case font-normal text-zinc-400">({u.optional ?? 'optional'})</span>
                </label>
                <textarea
                  placeholder={isRental ? 'Any questions regarding deposit, duration, or property rules…' : (u.message_placeholder ?? 'Any questions or details for the seller…')}
                  value={message}
                  onChange={e => setMessage(e.target.value)}
                  rows={2}
                  className={inputClass + ' resize-none'}
                />
              </div>

              {error && (
                <div className="flex items-center gap-2 text-red-500 text-xs bg-red-50 dark:bg-red-900/20 px-3 py-2 rounded-xl">
                  <AlertCircle size={12} /> {error}
                </div>
              )}

              <button type="submit" disabled={loading}
                className={`w-full text-white rounded-xl py-3 text-sm font-bold shadow-lg transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2 ${
                  isRental 
                    ? 'bg-teal-600 hover:bg-teal-700 shadow-teal-600/20' 
                    : 'bg-orange-600 hover:bg-orange-700 shadow-orange-600/20'
                }`}>
                {loading ? <Loader2 size={16} className="animate-spin" /> : <Phone size={16} />}
                {loading ? (u.buy_modal_sending ?? 'Sending…') : (u.send_contact_request ?? 'Send Request')}
              </button>

              {/* Login prompt for non-buyers */}
              {!isBuyer && (
                <p className="text-center text-[11px] text-zinc-400">
                  Have a user account?{' '}
                  <Link href="/buyer/login" onClick={onClose} className="text-orange-600 font-bold hover:underline">
                    Sign in
                  </Link>{' '}
                  for auto-fill
                </p>
              )}
            </form>
          </>
        ) : step === 'success' ? (
          <div className="text-center py-4 space-y-4">
            <div className="w-16 h-16 bg-green-500 rounded-full flex items-center justify-center text-white mx-auto shadow-xl shadow-green-500/30">
              <Check size={32} strokeWidth={3} />
            </div>
            <div>
              <h2 className="text-lg font-black">{u.request_sent ?? 'Request Sent'}</h2>
              <p className="text-sm text-zinc-500 mt-1">
                {isRental ? 'Your rental request has been sent to the landlord.' : (u.request_sent_msg ?? 'Your request has been sent to the seller.')}
              </p>
              <p className="text-xs text-zinc-400 mt-2">
                {isRental ? 'The landlord will call you soon.' : (u.buy_modal_seller_will_call ?? 'The seller will call you soon.')}{' '}
                <span className="font-bold text-zinc-600 dark:text-zinc-300">
                  {isBuyer && buyerPhone ? buyerPhone : phone}
                </span>
              </p>
            </div>
            <button type="button" onClick={onClose}
              className="w-full bg-orange-600 hover:bg-orange-700 text-white rounded-xl py-3 font-bold text-sm transition-colors">
              {u.done ?? 'Done'}
            </button>
          </div>
        ) : null}
      </div>
    </div>
  );
}
