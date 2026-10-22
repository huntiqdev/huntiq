'use client';
import { useState } from 'react';
import { useLocale } from '@/hooks/useLocale';
import { X, ChevronDown, ChevronUp, CheckCircle, AlertCircle, ImageIcon, MapPin, ShoppingBag, Clock } from 'lucide-react';

interface Props { onClose: () => void }

export function SellerHelp({ onClose }: Props) {
  const { t } = useLocale();
  const h = (t as any).sellerHelp ?? {};
  const [openFaq, setOpenFaq] = useState<number | null>(null);

  const steps = [
    { icon: ShoppingBag, key: 'step1' },
    { icon: ImageIcon,   key: 'step2' },
    { icon: MapPin,      key: 'step3' },
    { icon: Clock,       key: 'step4' },
  ];

  const faqs = [
    { q: h.faq1q, a: h.faq1a },
    { q: h.faq2q, a: h.faq2a },
    { q: h.faq3q, a: h.faq3a },
    { q: h.faq4q, a: h.faq4a },
    { q: h.faq5q, a: h.faq5a },
  ].filter(f => f.q);

  const tips = [h.tip1, h.tip2, h.tip3, h.tip4].filter(Boolean);

  return (
    <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4">
      <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose} />
      <div className="relative z-10 w-full sm:max-w-lg bg-white dark:bg-zinc-900 rounded-t-3xl sm:rounded-3xl shadow-2xl border border-zinc-200 dark:border-zinc-800 max-h-[90vh] flex flex-col">

        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-zinc-100 dark:border-zinc-800 shrink-0">
          <div>
            <h2 className="font-black text-zinc-900 dark:text-white text-base">{h.title ?? 'Seller Help'}</h2>
            <p className="text-xs text-zinc-400 mt-0.5">{h.subtitle ?? 'How to use the seller panel'}</p>
          </div>
          <button onClick={onClose} className="p-2 rounded-xl text-zinc-400 hover:text-zinc-600 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors">
            <X size={18} />
          </button>
        </div>

        <div className="overflow-y-auto flex-1 px-5 py-4 space-y-6">

          {/* How it works */}
          <section>
            <h3 className="text-xs font-black uppercase tracking-widest text-zinc-400 mb-3">{h.howItWorks ?? 'How It Works'}</h3>
            <div className="grid grid-cols-2 gap-2">
              {steps.map(({ icon: Icon, key }, i) => (
                <div key={key} className="bg-zinc-50 dark:bg-zinc-800 rounded-2xl p-3 flex flex-col gap-2">
                  <div className="flex items-center gap-2">
                    <div className="w-6 h-6 rounded-lg bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center">
                      <Icon size={13} className="text-orange-600" />
                    </div>
                    <span className="text-[10px] font-black text-orange-600">0{i + 1}</span>
                  </div>
                  <p className="text-xs font-bold text-zinc-800 dark:text-zinc-200 leading-tight">{h[`${key}Title`] ?? key}</p>
                  <p className="text-[11px] text-zinc-500 leading-relaxed">{h[`${key}Desc`]}</p>
                </div>
              ))}
            </div>
          </section>

          {/* Image rules */}
          <section className="bg-orange-50 dark:bg-orange-900/20 border border-orange-200 dark:border-orange-800 rounded-2xl p-4">
            <div className="flex items-center gap-2 mb-2">
              <ImageIcon size={14} className="text-orange-600" />
              <h3 className="text-xs font-black text-orange-700 dark:text-orange-400 uppercase tracking-widest">{h.imageRules ?? 'Image Rules'}</h3>
            </div>
            <ul className="space-y-1">
              {[h.imageRule1, h.imageRule2, h.imageRule3].filter(Boolean).map((rule, i) => (
                <li key={i} className="flex items-start gap-2 text-xs text-orange-800 dark:text-orange-300">
                  <CheckCircle size={12} className="text-orange-500 mt-0.5 shrink-0" />
                  {rule}
                </li>
              ))}
            </ul>
          </section>

          {/* Tips */}
          {tips.length > 0 && (
            <section>
              <h3 className="text-xs font-black uppercase tracking-widest text-zinc-400 mb-3">{h.tips ?? 'Tips'}</h3>
              <div className="space-y-2">
                {tips.map((tip, i) => (
                  <div key={i} className="flex items-start gap-2.5 bg-zinc-50 dark:bg-zinc-800 rounded-xl px-3 py-2.5">
                    <AlertCircle size={13} className="text-amber-500 mt-0.5 shrink-0" />
                    <p className="text-xs text-zinc-600 dark:text-zinc-400">{tip}</p>
                  </div>
                ))}
              </div>
            </section>
          )}

          {/* FAQ */}
          {faqs.length > 0 && (
            <section>
              <h3 className="text-xs font-black uppercase tracking-widest text-zinc-400 mb-3">{h.faq ?? 'FAQ'}</h3>
              <div className="space-y-2">
                {faqs.map((faq, i) => (
                  <div key={i} className="border border-zinc-200 dark:border-zinc-700 rounded-xl overflow-hidden">
                    <button onClick={() => setOpenFaq(openFaq === i ? null : i)}
                      className="w-full flex items-center justify-between px-4 py-3 text-left hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors">
                      <span className="text-sm font-semibold text-zinc-800 dark:text-zinc-200 pr-2">{faq.q}</span>
                      {openFaq === i ? <ChevronUp size={14} className="text-zinc-400 shrink-0" /> : <ChevronDown size={14} className="text-zinc-400 shrink-0" />}
                    </button>
                    {openFaq === i && (
                      <div className="px-4 pb-3 text-sm text-zinc-500 dark:text-zinc-400 border-t border-zinc-100 dark:border-zinc-800 pt-2">
                        {faq.a}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </section>
          )}
        </div>
      </div>
    </div>
  );
}
