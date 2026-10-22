'use client';
import { useState } from 'react';
import { useLocale } from '@/hooks/useLocale';
import { HelpCircle, ChevronDown, ChevronUp, Shield, ShoppingBag, Tag } from 'lucide-react';

export default function HelpPage() {
  const { t } = useLocale();
  const h = (t as any).help_page ?? {};
  const [openFaq, setOpenFaq] = useState<number | null>(null);

  const faqs = [
    { q: h.faq1q, a: h.faq1a },
    { q: h.faq2q, a: h.faq2a },
    { q: h.faq3q, a: h.faq3a },
    { q: h.faq4q, a: h.faq4a },
    { q: h.faq5q, a: h.faq5a },
    { q: h.faq6q, a: h.faq6a },
    { q: h.faq7q, a: h.faq7a },
    { q: h.faq8q, a: h.faq8a },
    { q: h.faq9q, a: h.faq9a },
    { q: h.faq10q, a: h.faq10a },
  ].filter(f => f.q && f.q !== 'Faq1q' && !f.q.startsWith('Faq'));

  const steps = [
    { icon: ShoppingBag, title: h.step1_title, desc: h.step1_desc },
    { icon: Tag,         title: h.step2_title, desc: h.step2_desc },
    { icon: ShoppingBag, title: h.step3_title, desc: h.step3_desc },
  ].filter(s => s.title);

  const tips = [h.tip1, h.tip2, h.tip3, h.tip4, h.tip5].filter(Boolean);

  return (
    <div className="max-w-3xl mx-auto px-4 py-12 space-y-10">
      {/* Header */}
      <div className="text-center">
        <div className="inline-flex items-center justify-center w-12 h-12 rounded-2xl bg-orange-600/10 text-orange-600 mb-4">
          <HelpCircle size={24} />
        </div>
        <h1 className="text-3xl font-black mb-2">{h.title}</h1>
        <p className="text-sm opacity-60 max-w-md mx-auto">{h.subtitle}</p>
      </div>

      {/* How it works */}
      {steps.length > 0 && (
        <section>
          <h2 className="text-xs font-black uppercase tracking-widest opacity-40 mb-4">{h.how_it_works}</h2>
          <div className="grid sm:grid-cols-3 gap-3">
            {steps.map(({ icon: Icon, title, desc }, i) => (
              <div key={i} className="bg-white/70 dark:bg-zinc-900 rounded-2xl border border-violet-100 dark:border-zinc-800 p-4 backdrop-blur-sm">
                <div className="flex items-center gap-2 mb-2">
                  <div className="w-8 h-8 rounded-xl bg-orange-600/10 text-orange-600 flex items-center justify-center">
                    <Icon size={16} />
                  </div>
                  <span className="text-xs font-black text-orange-600">0{i + 1}</span>
                </div>
                <p className="font-bold text-sm mb-1">{title}</p>
                <p className="text-xs opacity-60 leading-relaxed">{desc}</p>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* Safety tips */}
      {tips.length > 0 && (
        <section>
          <h2 className="text-xs font-black uppercase tracking-widest opacity-40 mb-4">
            <Shield size={12} className="inline mr-1" />{h.safety_tips ?? h.safety}
          </h2>
          <div className="space-y-2">
            {tips.map((tip, i) => (
              <div key={i} className="flex items-start gap-3 bg-white/70 dark:bg-zinc-900 rounded-xl border border-violet-100 dark:border-zinc-800 px-4 py-3 backdrop-blur-sm">
                <span className="text-orange-600 font-black text-xs mt-0.5 shrink-0">{i + 1}.</span>
                <p className="text-sm opacity-70">{tip}</p>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* FAQ */}
      {faqs.length > 0 && (
        <section>
          <h2 className="text-xs font-black uppercase tracking-widest opacity-40 mb-4">{h.faq}</h2>
          <div className="space-y-2">
            {faqs.map((faq, i) => (
              <div key={i} className="bg-white/70 dark:bg-zinc-900 rounded-2xl border border-violet-100 dark:border-zinc-800 overflow-hidden backdrop-blur-sm">
                <button
                  onClick={() => setOpenFaq(openFaq === i ? null : i)}
                  className="w-full flex items-center justify-between px-5 py-4 text-left hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors"
                >
                  <span className="font-semibold text-sm pr-4">{faq.q}</span>
                  {openFaq === i
                    ? <ChevronUp size={16} className="opacity-40 shrink-0" />
                    : <ChevronDown size={16} className="opacity-40 shrink-0" />}
                </button>
                {openFaq === i && (
                  <div className="px-5 pb-4 text-sm opacity-60 leading-relaxed border-t border-zinc-100 dark:border-zinc-800 pt-3">
                    {faq.a}
                  </div>
                )}
              </div>
            ))}
          </div>
        </section>
      )}
    </div>
  );
}
