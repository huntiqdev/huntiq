'use client';

import Link from 'next/link';
import { useState, useEffect, Suspense } from 'react';
import { usePathname, useSearchParams } from 'next/navigation';
import { useLocale } from '@/hooks/useLocale';
import { useAuth } from '@/context/AuthContext';
import { useTheme } from '@/context/ThemeContext';
import { SearchBar } from '@/components/public/SearchBar';
import { UserAccountMenu } from '@/components/layout/UserAccountMenu';
import {
  Globe,
  Moon,
  Sun,
  ShoppingBag,
  ChevronDown,
  HelpCircle,
  X,
  ChevronRight,
  Shield,
  Zap,
  ImageIcon,
  MapPin,
  Layers,
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';

// ── Help panel modal ─────────────────────────────────────────────────────────
function HelpPanel({ onClose }: { onClose: () => void }) {
  const { t } = useLocale();
  const h = t.sellerHelp ?? {};
  const [openFaq, setOpenFaq] = useState<number | null>(null);

  const steps = [
    { icon: ShoppingBag, title: h.step1Title ?? 'Register & Login',  desc: h.step1Desc  ?? 'Create your account and log in.' },
    { icon: ImageIcon,   title: h.step2Title ?? 'Add Images',        desc: h.step2Desc  ?? 'Upload 2–4 clear photos.' },
    { icon: MapPin,      title: h.step3Title ?? 'Set Location',      desc: h.step3Desc  ?? 'Select your region and area.' },
    { icon: Zap,         title: h.step4Title ?? 'Wait for Approval', desc: h.step4Desc  ?? 'Admin reviews and publishes.' },
  ];

  const faqs = [
    { q: h.faq1q, a: h.faq1a },
    { q: h.faq2q, a: h.faq2a },
    { q: h.faq3q, a: h.faq3a },
    { q: h.faq4q, a: h.faq4a },
    { q: h.faq5q, a: h.faq5a },
  ].filter(f => f.q);

  return (
    <div className="fixed inset-0 z-[60] flex items-end sm:items-center justify-center p-0 sm:p-4 animate-fade-in">
      <div className="absolute inset-0 bg-slate-950/60 backdrop-blur-sm" onClick={onClose} />
      <div className="relative z-10 w-full sm:max-w-md bg-white dark:bg-[#0E1B33] rounded-t-3xl sm:rounded-3xl shadow-2xl border border-slate-200 dark:border-slate-800 max-h-[88vh] flex flex-col">

        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-slate-100 dark:border-slate-800 shrink-0">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-xl bg-[#1769E0]/10 flex items-center justify-center">
              <HelpCircle size={16} className="text-[#1769E0]" />
            </div>
            <div>
              <p className="font-black text-sm text-slate-900 dark:text-white">{h.title ?? 'Help Center'}</p>
              <p className="text-[10px] text-slate-500 dark:text-slate-400">{h.subtitle ?? 'How to use HUNTIQ'}</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
          >
            <X size={16} />
          </button>
        </div>

        <div className="overflow-y-auto flex-1 px-5 py-4 space-y-5 admin-scrollbar">

          {/* How it works */}
          <div>
            <p className="text-[10px] font-black uppercase tracking-widest text-slate-400 dark:text-slate-500 mb-3">
              {h.howItWorks ?? 'How It Works'}
            </p>
            <div className="grid grid-cols-2 gap-2">
              {steps.map(({ icon: Icon, title, desc }, i) => (
                <div key={i} className="bg-slate-50 dark:bg-slate-800/60 rounded-2xl p-3 space-y-1.5 border border-slate-100 dark:border-slate-800">
                  <div className="flex items-center gap-2">
                    <div className="w-6 h-6 rounded-lg bg-[#EAF3FF] dark:bg-[#1769E0]/20 flex items-center justify-center">
                      <Icon size={12} className="text-[#1769E0]" />
                    </div>
                    <span className="text-[10px] font-black text-[#1769E0]">0{i + 1}</span>
                  </div>
                  <p className="text-xs font-bold text-slate-800 dark:text-slate-200 leading-tight">{title}</p>
                  <p className="text-[10px] text-slate-500 dark:text-slate-400 leading-relaxed">{desc}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Image rules */}
          <div className="bg-[#EAF3FF] dark:bg-[#1769E0]/15 border border-[#1769E0]/30 rounded-2xl p-3.5 space-y-3">
            <div>
              <div className="flex items-center gap-2 mb-2">
                <ImageIcon size={13} className="text-[#1769E0]" />
                <p className="text-[10px] font-black uppercase tracking-widest text-[#1055BD] dark:text-[#38A5FF]">
                  {h.imageRules ?? 'Image Rules'}
                </p>
              </div>
              <ul className="space-y-1">
                {[
                  h.imageRule1 ?? 'Min 2 images required, max 4.',
                  h.imageRule2 ?? 'Use clear, well-lit photos.',
                  h.imageRule3 ?? 'Products need 2+ images to be approved.',
                ].map((r, i) => (
                  <li key={i} className="flex items-start gap-2 text-xs text-slate-800 dark:text-slate-200">
                    <Shield size={11} className="text-[#1769E0] mt-0.5 shrink-0" />
                    {r}
                  </li>
                ))}
              </ul>
            </div>

            {/* Pro Tips */}
            {h.tip1 && (
              <div className="pt-2 border-t border-[#1769E0]/20">
                <p className="text-[10px] font-black uppercase tracking-widest text-[#1055BD] dark:text-[#38A5FF] mb-1.5">
                  {h.tips ?? 'Pro Tips'}
                </p>
                <ul className="space-y-1">
                  {[h.tip1, h.tip2, h.tip3, h.tip4].filter(Boolean).map((tip, i) => (
                    <li key={i} className="flex items-start gap-2 text-xs text-slate-700 dark:text-slate-300">
                      <span className="text-[#1769E0] text-[10px] font-bold mt-0.5">•</span>
                      {tip}
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>

          {/* FAQ */}
          {faqs.length > 0 && (
            <div>
              <p className="text-[10px] font-black uppercase tracking-widest text-slate-400 dark:text-slate-500 mb-3">
                {h.faq ?? 'FAQ'}
              </p>
              <div className="space-y-1.5">
                {faqs.map((faq, i) => (
                  <div key={i} className="border border-slate-200 dark:border-slate-800 rounded-xl overflow-hidden">
                    <button
                      onClick={() => setOpenFaq(openFaq === i ? null : i)}
                      className="w-full flex items-center justify-between px-3.5 py-2.5 text-left hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors"
                    >
                      <span className="text-xs font-semibold text-slate-800 dark:text-slate-200 pr-2">{faq.q}</span>
                      <ChevronRight size={13} className={`text-slate-400 shrink-0 transition-transform ${openFaq === i ? 'rotate-90' : ''}`} />
                    </button>
                    {openFaq === i && (
                      <div className="px-3.5 pb-3 pt-2 text-xs text-slate-600 dark:text-slate-400 border-t border-slate-100 dark:border-slate-800">
                        {faq.a}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

// ── Header ───────────────────────────────────────────────────────────────────
export function Header() {
  const { t, locale, setLocale, locales, localeLabels } = useLocale();
  const { resolvedTheme, toggleTheme } = useTheme();
  const [langOpen, setLangOpen] = useState(false);
  const [helpOpen, setHelpOpen] = useState(false);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  // Close lang dropdown on outside click
  useEffect(() => {
    if (!langOpen) return;
    const id = setTimeout(() => {
      const close = () => setLangOpen(false);
      document.addEventListener('click', close, { once: true });
    }, 0);
    return () => clearTimeout(id);
  }, [langOpen]);

  const isDark = mounted && resolvedTheme === 'dark';

  return (
    <>
      {helpOpen && <HelpPanel onClose={() => setHelpOpen(false)} />}

      <header className="sticky top-0 z-50 w-full border-b backdrop-blur-xl bg-white/95 dark:bg-[#091426]/95 border-slate-200 dark:border-slate-800 transition-colors">
        
        {/* ── Top Utility Row ──────────────────────────────────────────────── */}
        <div className="hidden sm:block border-b border-slate-100 dark:border-slate-800/80 bg-slate-50/80 dark:bg-[#070F1E]/80 text-[11px] text-slate-500 dark:text-slate-400">
          <div className="max-w-[1520px] mx-auto px-3 sm:px-6 lg:px-8 h-7 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <span className="flex items-center gap-1.5 font-semibold text-emerald-600 dark:text-emerald-400">
                <Shield size={12} className="text-[#16A36A]" />
                <span>{(t.nav as any)?.verified_ethiopian_listings ?? 'Verified Marketplace • 100% Ethiopian Listings'}</span>
              </span>
            </div>

            <div className="flex items-center gap-3">
              {/* Language Selector */}
              <div className="relative">
                <button
                  onClick={() => setLangOpen(v => !v)}
                  className="flex items-center gap-1 hover:text-[#1769E0] transition-colors font-bold uppercase tracking-wider text-[10px]"
                >
                  <Globe size={12} />
                  <span>{localeLabels[locale] ?? locale.toUpperCase()}</span>
                  <ChevronDown size={10} className={`transition-transform ${langOpen ? 'rotate-180' : ''}`} />
                </button>
                {langOpen && (
                  <div className="absolute right-0 mt-1.5 w-40 rounded-2xl border shadow-2xl z-50 overflow-hidden bg-white dark:bg-[#0E1B33] border-slate-200 dark:border-slate-800">
                    {locales.map((l: string) => (
                      <button
                        key={l}
                        onClick={() => { setLocale(l); setLangOpen(false); }}
                        className={`w-full text-left px-3.5 py-2 text-xs font-bold hover:bg-[#1769E0] hover:text-white transition-colors flex items-center justify-between ${
                          locale === l ? 'text-[#1769E0] dark:text-[#38A5FF]' : 'text-slate-700 dark:text-slate-300'
                        }`}
                      >
                        {localeLabels[l] ?? l.toUpperCase()}
                        {locale === l && <span className="w-1.5 h-1.5 rounded-full bg-[#1769E0] inline-block" />}
                      </button>
                    ))}
                  </div>
                )}
              </div>

              {/* Dark / Light Toggle */}
              <button
                onClick={toggleTheme}
                className="p-1 rounded-md hover:text-[#1769E0] transition-colors"
                aria-label="Toggle dark mode"
              >
                {mounted ? (isDark ? <Sun size={13} /> : <Moon size={13} />) : <span className="w-3.5 h-3.5 block" />}
              </button>
            </div>
          </div>
        </div>

        {/* ── Main Header Row ──────────────────────────────────────────────── */}
        <div className="max-w-[1520px] mx-auto px-3 sm:px-6 lg:px-8 h-15 py-2 flex items-center gap-2 sm:gap-4">

          {/* Logo */}
          <HuntiqLogo
            size="sm"
            showText
            subtitle="Marketplace"
            href="/"
            priority
          />

          {/* Global Search Bar (Balanced Width on Large Screens) */}
          <div className="flex-1 min-w-0 max-w-lg lg:max-w-xl mx-2 sm:mx-6">
            <SearchBar />
          </div>

          {/* Right Header Controls */}
          <div className="flex items-center gap-1.5 sm:gap-2 ml-auto shrink-0">

            {/* Mobile Language Selector (sm:hidden) */}
            <div className="relative sm:hidden">
              <button
                type="button"
                onClick={() => setLangOpen(v => !v)}
                className="flex items-center gap-1 p-2 rounded-xl border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-black text-[11px] uppercase tracking-wider bg-slate-50 dark:bg-slate-900 tap-bounce active:scale-95"
                aria-label="Select language"
              >
                <Globe size={13} className="text-[#1769E0]" />
                <span>{locale === 'am' ? 'አማ' : locale === 'ti' ? 'ትግ' : locale.toUpperCase()}</span>
              </button>
              {langOpen && (
                <div className="absolute right-0 mt-1.5 w-40 rounded-2xl border shadow-2xl z-50 overflow-hidden bg-white dark:bg-[#0E1B33] border-slate-200 dark:border-slate-800">
                  {locales.map((l: string) => (
                    <button
                      key={l}
                      onClick={() => { setLocale(l); setLangOpen(false); }}
                      className={`w-full text-left px-3.5 py-2 text-xs font-bold hover:bg-[#1769E0] hover:text-white transition-colors flex items-center justify-between tap-bounce ${
                        locale === l ? 'text-[#1769E0]' : ''
                      }`}
                    >
                      {localeLabels[l] ?? l.toUpperCase()}
                      {locale === l && <span className="w-1.5 h-1.5 rounded-full bg-[#1769E0] inline-block" />}
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* Mobile Dark / Light Toggle (sm:hidden) */}
            <button
              onClick={toggleTheme}
              className="sm:hidden p-2 rounded-xl border border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400 hover:text-[#1769E0] bg-slate-50 dark:bg-slate-900 transition-colors tap-bounce active:scale-95"
              aria-label="Toggle dark mode"
            >
              {mounted ? (isDark ? <Sun size={13} /> : <Moon size={13} />) : <span className="w-3.5 h-3.5 block" />}
            </button>

            {/* User Account / Role Menu */}
            <UserAccountMenu />
          </div>
        </div>
      </header>
    </>
  );
}
