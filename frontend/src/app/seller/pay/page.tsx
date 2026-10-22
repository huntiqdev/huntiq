'use client';
import { useEffect, useState, Suspense, FormEvent } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { useLocale } from '@/hooks/useLocale';
import { usePaymentSettings } from '@/hooks/usePaymentSettings';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { validateEthioPhone } from '@/lib/validations/common';
import { PAYMENT_METHODS, getMerchantTarget, type PaymentMethodValue } from '@/lib/constants/payment';
import { INPUT_CLS } from '@/lib/constants/styles';
import { CheckCircle, AlertCircle, Phone, ExternalLink, Copy } from 'lucide-react';

type Step = 'choose' | 'confirm' | 'done';


// ── Pay page content ─────────────────────────────────────────────────────────
function PayContent() {
  const { user } = useAuth();
  const { t } = useLocale();
  const { settings } = usePaymentSettings();
  const router = useRouter();
  const params = useSearchParams();

  const productId    = parseInt(params.get('product_id') ?? '0', 10) || 0;
  const feeAmount    = parseFloat(params.get('fee') ?? '0') || 0;
  const productTitle = decodeURIComponent(params.get('title') ?? '');

  const [methodId, setMethodId] = useState<PaymentMethodValue>('telebirr');
  const [step, setStep]         = useState<Step>('choose');
  const [reference, setRef]     = useState('');
  const [payerPhone, setPhone]  = useState('');
  const [phoneError, setPhoneError] = useState('');
  const [refError, setRefError] = useState('');
  const [submitting, setSub]    = useState(false);
  const [error, setError]       = useState('');
  const [copied, setCopied]     = useState<string | null>(null);

  useEffect(() => {
    if (!productId || feeAmount <= 0) router.replace('/seller/products');
  }, [productId, feeAmount, router]);

  const method = PAYMENT_METHODS.find(m => m.value === methodId)!;
  const merchantTarget = getMerchantTarget(methodId, settings);

  const copy = (text: string, key: string) => {
    navigator.clipboard.writeText(text).then(() => {
      setCopied(key);
      setTimeout(() => setCopied(null), 2000);
    });
  };

  const phoneResult = payerPhone ? validateEthioPhone(payerPhone) : null;

  const handleConfirm = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    let ok = true;

    if (!reference.trim()) {
      setRefError(t.ui.txn_ref_required);
      ok = false;
    } else setRefError('');

    const pr = validateEthioPhone(payerPhone);
    if (!pr.valid) {
      setPhoneError(pr.error ?? t.ui.validation_phone_invalid);
      ok = false;
    } else setPhoneError('');

    if (!ok || !user) return;

    setSub(true); setError('');
    try {
      await api.sellerPayListingFee({
        product_id: productId,
        amount: feeAmount,
        payment_method: methodId,
        payer_phone: pr.normalized,
        reference: reference.trim(),
        notes: `Listing fee via ${method.label} — ref: ${reference.trim()}`,
      }, user.token);
      setStep('done');
    } catch (err: any) {
      setError(err.message ?? t.ui.purchase_error);
    } finally {
      setSub(false);
    }
  };

  const cls = INPUT_CLS;

  // ── Done screen ────────────────────────────────────────────────────────────
  if (step === 'done') {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-5 text-center px-4">
        <div className="w-20 h-20 bg-green-500 rounded-full flex items-center justify-center shadow-xl shadow-green-500/30">
          <CheckCircle size={40} className="text-white" strokeWidth={2.5} />
        </div>
        <div>
          <h1 className="text-2xl font-black text-zinc-900 dark:text-white">{t.ui.payment_submitted_title}</h1>
          <p className="text-zinc-500 mt-2 text-sm max-w-xs mx-auto">{t.ui.payment_submitted_desc}</p>
          <p className="text-xs text-zinc-400 mt-1">{t.ui.payment_submitted_note}</p>
        </div>
        <button onClick={() => router.push('/seller/products')}
          className="bg-orange-600 hover:bg-orange-700 text-white px-8 py-3 rounded-2xl font-bold transition-colors">
          {t.ui.view_my_products}
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto space-y-5 py-6 px-4">
      <div>
        <h1 className="text-xl font-black text-zinc-900 dark:text-white">{t.ui.pay_listing_fee_title}</h1>
        {productTitle && <p className="text-xs text-zinc-500 mt-0.5 truncate">{productTitle}</p>}
      </div>

      {/* Fee summary */}
      <div className="bg-orange-50 dark:bg-orange-900/20 border border-orange-200 dark:border-orange-800 rounded-2xl p-4">
        <p className="text-xs text-orange-600 font-black uppercase tracking-widest mb-1">{t.ui.amount_due}</p>
        <p className="text-3xl font-black text-orange-600">ETB {feeAmount.toLocaleString()}</p>
        <p className="text-xs text-zinc-500 mt-1">{t.ui.amount_due_desc}</p>
      </div>

      {/* Step: choose method */}
      {step === 'choose' && (
        <div className="space-y-3">
          <p className="text-sm font-bold text-zinc-700 dark:text-zinc-300">{t.ui.choose_payment_method}</p>
          {PAYMENT_METHODS.map(m => (
            <button key={m.value}
              onClick={() => { setMethodId(m.value); setStep('confirm'); }}
              className="w-full flex items-center gap-4 bg-white dark:bg-zinc-900 border-2 border-zinc-200 dark:border-zinc-700 hover:border-orange-400 rounded-2xl p-4 transition-all text-left">
              <span className="text-3xl">{m.emoji}</span>
              <div className="flex-1 min-w-0">
                <p className="font-black text-zinc-900 dark:text-white">{m.label}</p>
              </div>
              <a href={m.appLink} target="_blank" rel="noopener noreferrer"
                onClick={e => e.stopPropagation()}
                className="text-xs text-zinc-400 hover:text-orange-600 flex items-center gap-1 shrink-0">
                <ExternalLink size={12} /> {t.ui.download_app}
              </a>
            </button>
          ))}
          <button onClick={() => router.push('/seller/products')}
            className="w-full text-xs text-zinc-400 hover:text-zinc-600 text-center py-2">
            {t.ui.skip_pay_now}
          </button>
        </div>
      )}

      {/* Step: confirm payment */}
      {step === 'confirm' && (
        <form onSubmit={handleConfirm} className="space-y-4">

          {/* Instructions card */}
          <div className={`border rounded-2xl p-4 space-y-3 ${method.activeBg} border-zinc-200 dark:border-zinc-700`}>
            <div className="flex items-center justify-between">
              <p className="text-sm font-black text-zinc-800 dark:text-zinc-200">
                {method.emoji} {method.label} {t.ui.instructions_label}
              </p>
              <a href={method.appLink} target="_blank" rel="noopener noreferrer"
                className="flex items-center gap-1 text-xs font-bold text-orange-600 hover:text-orange-700">
                <ExternalLink size={11} /> {method.appName}
              </a>
            </div>

            {/* USSD option */}
            {method.ussd && (
              <div className="bg-white dark:bg-zinc-900 rounded-xl p-3 space-y-1.5">
                <p className="text-xs font-bold text-zinc-600 dark:text-zinc-400">{t.ui.ussd_dial_label}</p>
                <div className="flex items-center gap-2">
                  <code className={`text-sm font-black ${method.accentColor} bg-orange-50 dark:bg-orange-900/30 px-2 py-1 rounded-lg`}>
                    {method.ussd}
                  </code>
                  <button type="button" onClick={() => copy(method.ussd!, 'ussd')}
                    className="text-zinc-400 hover:text-orange-600 transition-colors">
                    <Copy size={13} />
                  </button>
                  {copied === 'ussd' && <span className="text-xs text-green-600">{t.ui.copied_label}</span>}
                </div>
              </div>
            )}

            {/* Step-by-step instructions */}
            <ol className="text-xs text-zinc-600 dark:text-zinc-400 space-y-1.5 list-none">
              {methodId === 'telebirr' && <>
                <li className="flex gap-2">
                  <span className={`font-black ${method.accentColor} shrink-0`}>1.</span>
                  {t.ui.step_dial_or_app.replace('{ussd}', method.ussd!).replace('{app}', 'Telebirr')}
                </li>
                <li className="flex gap-2">
                  <span className={`font-black ${method.accentColor} shrink-0`}>2.</span>
                  {t.ui.step_send_money.replace('{amount}', feeAmount.toLocaleString())}
                </li>
                <li className="flex gap-2">
                  <span className={`font-black ${method.accentColor} shrink-0`}>3.</span>
                  <span>{t.ui.step_enter_merchant}
                    <span className="inline-flex items-center gap-1 ml-1">
                      <strong>{merchantTarget}</strong>
                      <button type="button" onClick={() => copy(merchantTarget, 'phone')}
                        className="text-zinc-400 hover:text-orange-600"><Copy size={11} /></button>
                      {copied === 'phone' && <span className="text-green-600">{t.ui.copied_label}</span>}
                    </span>
                  </span>
                </li>
                <li className="flex gap-2">
                  <span className={`font-black ${method.accentColor} shrink-0`}>4.</span>
                  {t.ui.step_copy_txn_id}
                </li>
              </>}

              {methodId === 'safaricom' && <>
                <li className="flex gap-2">
                  <span className={`font-black ${method.accentColor} shrink-0`}>1.</span>
                  {t.ui.step_dial_or_app.replace('{ussd}', method.ussd!).replace('{app}', 'M-PESA')}
                </li>
                <li className="flex gap-2">
                  <span className={`font-black ${method.accentColor} shrink-0`}>2.</span>
                  {t.ui.step_send_money.replace('{amount}', feeAmount.toLocaleString())}
                </li>
                <li className="flex gap-2">
                  <span className={`font-black ${method.accentColor} shrink-0`}>3.</span>
                  <span>{t.ui.step_enter_merchant}
                    <span className="inline-flex items-center gap-1 ml-1">
                      <strong>{merchantTarget}</strong>
                      <button type="button" onClick={() => copy(merchantTarget, 'phone')}
                        className="text-zinc-400 hover:text-green-600"><Copy size={11} /></button>
                      {copied === 'phone' && <span className="text-green-600">{t.ui.copied_label}</span>}
                    </span>
                  </span>
                </li>
                <li className="flex gap-2">
                  <span className={`font-black ${method.accentColor} shrink-0`}>4.</span>
                  {t.ui.step_copy_txn_id}
                </li>
              </>}

              {methodId === 'cbe' && <>
                <li className="flex gap-2">
                  <span className={`font-black ${method.accentColor} shrink-0`}>1.</span>
                  {t.ui.step_open_cbe}
                </li>
                <li className="flex gap-2">
                  <span className={`font-black ${method.accentColor} shrink-0`}>2.</span>
                  {t.ui.step_transfer_to.replace('{amount}', feeAmount.toLocaleString())}
                </li>
                <li className="flex gap-2 pl-4">
                  <span className="inline-flex items-center gap-1">
                    <code className={`font-black ${method.accentColor}`}>{merchantTarget}</code>
                    <button type="button" onClick={() => copy(merchantTarget, 'account')}
                      className="text-zinc-400 hover:text-blue-600"><Copy size={11} /></button>
                    {copied === 'account' && <span className="text-green-600">{t.ui.copied_label}</span>}
                  </span>
                </li>
                <li className="flex gap-2">
                  <span className={`font-black ${method.accentColor} shrink-0`}>3.</span>
                  {t.ui.step_note_ref}
                </li>
              </>}
            </ol>
          </div>

          {/* Payer phone */}
          <div>
            <label className="text-xs font-bold text-zinc-500 block mb-1">
              <Phone size={11} className="inline mr-1" />
              {t.ui.payer_phone_label}
            </label>
            <input type="tel" value={payerPhone}
              onChange={e => { setPhone(e.target.value); setPhoneError(''); }}
              placeholder={t.ui.payer_phone_placeholder}
              className={cls} />
            {phoneResult?.valid && (
              <p className="text-xs text-green-600 mt-1">
                {phoneResult.network === 'safaricom' ? t.ui.network_safaricom : t.ui.network_telebirr}
                {' — '}{phoneResult.normalized}
              </p>
            )}
            {phoneError && (
              <p className="text-xs text-red-500 mt-1 flex items-center gap-1">
                <AlertCircle size={11} /> {phoneError}
              </p>
            )}
            <p className="text-[10px] text-zinc-400 mt-1">{t.ui.payer_phone_hint}</p>
          </div>

          {/* Transaction reference */}
          <div>
            <label className="text-xs font-bold text-zinc-500 block mb-1">
              {t.ui.txn_ref_label}
            </label>
            <input value={reference}
              onChange={e => { setRef(e.target.value); setRefError(''); }}
              placeholder={method.refPlaceholder}
              className={cls} />
            {refError && (
              <p className="text-xs text-red-500 mt-1 flex items-center gap-1">
                <AlertCircle size={11} /> {refError}
              </p>
            )}
            <p className="text-[10px] text-zinc-400 mt-1">{t.ui.txn_ref_hint}</p>
          </div>

          {error && (
            <div className="flex items-start gap-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl px-3 py-2">
              <AlertCircle size={13} className="text-red-500 shrink-0 mt-0.5" />
              <p className="text-xs text-red-600 dark:text-red-400">{error}</p>
            </div>
          )}

          <div className="flex gap-2">
            <button type="button"
              onClick={() => { setStep('choose'); setError(''); setRef(''); setPhone(''); setPhoneError(''); setRefError(''); }}
              className="flex-1 py-3 rounded-xl border border-zinc-200 dark:border-zinc-700 text-sm font-bold text-zinc-600 dark:text-zinc-400 hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors">
              ← {t.ui.back}
            </button>
            <button type="submit" disabled={submitting}
              className="flex-1 py-3 rounded-xl bg-orange-600 hover:bg-orange-700 text-white font-bold text-sm transition-colors disabled:opacity-50 flex items-center justify-center gap-2">
              {submitting ? <><Spinner /> {t.ui.confirming}</> : t.ui.confirm_payment}
            </button>
          </div>

          <button type="button" onClick={() => router.push('/seller/products')}
            className="w-full text-xs text-zinc-400 hover:text-zinc-600 text-center py-1">
            {t.ui.skip_pay_later}
          </button>
        </form>
      )}
    </div>
  );
}

export default function SellerPayPage() {
  return (
    <Suspense fallback={<div className="flex items-center justify-center h-64"><Spinner /></div>}>
      <PayContent />
    </Suspense>
  );
}
