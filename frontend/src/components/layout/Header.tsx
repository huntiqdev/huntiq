'use client';

import Link from 'next/link';
import { useState, useEffect } from 'react';
import { usePathname } from 'next/navigation';
import { useLocale } from '@/hooks/useLocale';
import { SearchBar } from '@/components/public/SearchBar';
import { Globe, Moon, Sun, ShoppingBag, ChevronDown,
  HelpCircle, X, ChevronRight, Shield, Zap, ImageIcon, MapPin, User,
} from 'lucide-react';

// Language list is now driven by i18n registry — no hardcoding needed.
// To add a new language: add its JSON + register in i18n.ts. Header auto-updates.

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
    <div className="fixed inset-0 z-[60] flex items-end sm:items-center justify-center p-0 sm:p-4">
      <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />
      <div className="relative z-10 w-full sm:max-w-md bg-white dark:bg-zinc-900 rounded-t-3xl sm:rounded-3xl shadow-2xl border border-zinc-200 dark:border-zinc-800 max-h-[88vh] flex flex-col">

        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-zinc-100 dark:border-zinc-800 shrink-0">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-xl bg-orange-600/10 flex items-center justify-center">
              <HelpCircle size={16} className="text-orange-600" />
            </div>
            <div>
              <p className="font-black text-sm text-zinc-900 dark:text-white">{h.title ?? 'Help Center'}</p>
              <p className="text-[10px] text-zinc-400">{h.subtitle ?? 'How to use HUNTIQ'}</p>
            </div>
          </div>
          <button onClick={onClose}
            className="p-2 rounded-xl text-zinc-400 hover:text-zinc-600 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors">
            <X size={16} />
          </button>
        </div>

        <div className="overflow-y-auto flex-1 px-5 py-4 space-y-5">

          {/* How it works */}
          <div>
            <p className="text-[10px] font-black uppercase tracking-widest text-zinc-400 mb-3">
              {h.howItWorks ?? 'How It Works'}
            </p>
            <div className="grid grid-cols-2 gap-2">
              {steps.map(({ icon: Icon, title, desc }, i) => (
                <div key={i} className="bg-zinc-50 dark:bg-zinc-800 rounded-2xl p-3 space-y-1.5">
                  <div className="flex items-center gap-2">
                    <div className="w-6 h-6 rounded-lg bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center">
                      <Icon size={12} className="text-orange-600" />
                    </div>
                    <span className="text-[10px] font-black text-orange-600">0{i + 1}</span>
                  </div>
                  <p className="text-xs font-bold text-zinc-800 dark:text-zinc-200 leading-tight">{title}</p>
                  <p className="text-[10px] text-zinc-500 leading-relaxed">{desc}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Image rules */}
          <div className="bg-orange-50 dark:bg-orange-900/20 border border-orange-200 dark:border-orange-800 rounded-2xl p-3.5 space-y-3">
            <div>
              <div className="flex items-center gap-2 mb-2">
                <ImageIcon size={13} className="text-orange-600" />
                <p className="text-[10px] font-black uppercase tracking-widest text-orange-700 dark:text-orange-400">
                  {h.imageRules ?? 'Image Rules'}
                </p>
              </div>
              <ul className="space-y-1">
                {[
                  h.imageRule1 ?? 'Min 2 images required, max 4.',
                  h.imageRule2 ?? 'Use clear, well-lit photos.',
                  h.imageRule3 ?? 'Products need 2+ images to be approved.',
                ].map((r, i) => (
                  <li key={i} className="flex items-start gap-2 text-xs text-orange-800 dark:text-orange-300">
                    <Shield size={11} className="text-orange-500 mt-0.5 shrink-0" />
                    {r}
                  </li>
                ))}
              </ul>
            </div>

            {/* Pro Tips */}
            {h.tip1 && (
              <div className="pt-2 border-t border-orange-200/60 dark:border-orange-800/60">
                <p className="text-[10px] font-black uppercase tracking-widest text-orange-700 dark:text-orange-400 mb-1.5">
                  {h.tips ?? 'Pro Tips'}
                </p>
                <ul className="space-y-1">
                  {[h.tip1, h.tip2, h.tip3, h.tip4].filter(Boolean).map((tip, i) => (
                    <li key={i} className="flex items-start gap-2 text-xs text-orange-900 dark:text-orange-200">
                      <span className="text-orange-500 text-[10px] font-bold mt-0.5">•</span>
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
              <p className="text-[10px] font-black uppercase tracking-widest text-zinc-400 mb-3">
                {h.faq ?? 'FAQ'}
              </p>
              <div className="space-y-1.5">
                {faqs.map((faq, i) => (
                  <div key={i} className="border border-zinc-200 dark:border-zinc-700 rounded-xl overflow-hidden">
                    <button
                      onClick={() => setOpenFaq(openFaq === i ? null : i)}
                      className="w-full flex items-center justify-between px-3.5 py-2.5 text-left hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors">
                      <span className="text-xs font-semibold text-zinc-800 dark:text-zinc-200 pr-2">{faq.q}</span>
                      <ChevronRight size={13} className={`text-zinc-400 shrink-0 transition-transform ${openFaq === i ? 'rotate-90' : ''}`} />
                    </button>
                    {openFaq === i && (
                      <div className="px-3.5 pb-3 pt-2 text-xs text-zinc-500 dark:text-zinc-400 border-t border-zinc-100 dark:border-zinc-800">
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
  const pathname = usePathname();
  const [langOpen, setLangOpen] = useState(false);
  const [helpOpen, setHelpOpen] = useState(false);
  const [mounted, setMounted] = useState(false);
  const [dark, setDark] = useState(false);

  useEffect(() => {
    setDark(document.documentElement.classList.contains('dark'));
    setMounted(true);
  }, []);

  const toggleDark = () => {
    const next = !dark;
    setDark(next);
    document.documentElement.classList.toggle('dark', next);
    localStorage.setItem('habesha_theme', next ? 'dark' : 'light');
  };

  const isBuyPage = pathname === '/products' || pathname.startsWith('/products/');
  const isSellPage = pathname === '/sell';
  const isHomePage = pathname === '/';
  const ctaHref = isSellPage ? '/products' : '/sell';
  const ctaLabel = isSellPage ? t.nav.buy : t.nav.sell;

  // Close lang dropdown on outside click
  useEffect(() => {
    if (!langOpen) return;
    // Use setTimeout so this listener is added AFTER the current click event finishes
    const id = setTimeout(() => {
      const close = () => setLangOpen(false);
      document.addEventListener('click', close, { once: true });
    }, 0);
    return () => clearTimeout(id);
  }, [langOpen]);

  const [currentType, setCurrentType] = useState<string | null>(null);

  useEffect(() => {
    if (typeof window !== 'undefined') {
      const p = new URLSearchParams(window.location.search);
      setCurrentType(p.get('type'));
    }
  }, [pathname]);

  const isRentActive = pathname === '/rent' || (pathname === '/products' && currentType === 'rent');
  const isBuyActive = (pathname === '/products' && currentType !== 'rent') || pathname === '/buy';

  return (
    <>
      {helpOpen && <HelpPanel onClose={() => setHelpOpen(false)} />}

      <header className="sticky top-0 z-50 w-full border-b backdrop-blur-md bg-white/90 dark:bg-zinc-950/90 border-zinc-200 dark:border-zinc-800 overflow-visible">
        <div className="max-w-7xl mx-auto px-4 h-14 flex items-center gap-2">

          {/* Logo */}
          <Link href="/" className="flex items-center gap-2.5 shrink-0 group">
            <div className="relative w-8 h-8 rounded-xl overflow-hidden flex items-center justify-center bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-sm group-hover:scale-105 transition-transform">
              <img src="/huntiq.png" alt="HUNTIQ" className="w-full h-full object-contain p-0.5" />
            </div>
            <span className="text-base font-black tracking-tight text-zinc-900 dark:text-white hidden sm:block">
              {t.ui.app_name}
            </span>
          </Link>

          {/* Desktop nav */}
          <nav className="hidden lg:flex items-center gap-1 ml-2 border-r border-zinc-200 dark:border-zinc-800 pr-3 mr-1">
            <Link href="/"
              className={`px-3 py-1.5 rounded-lg text-xs font-black uppercase tracking-widest transition-all ${pathname === '/' ? 'bg-orange-600 text-white' : 'opacity-60 hover:opacity-100 hover:bg-zinc-100 dark:hover:bg-zinc-800'}`}>
              {t.footer.home}
            </Link>
            <Link href="/products?type=sale"
              className={`px-3 py-1.5 rounded-lg text-xs font-black uppercase tracking-widest transition-all ${isBuyActive ? 'bg-orange-600 text-white shadow-md shadow-orange-600/20' : 'opacity-60 hover:opacity-100 hover:bg-zinc-100 dark:hover:bg-zinc-800'}`}>
              🛒 {t.nav.buy ?? 'Buy'}
            </Link>
            <Link href="/products?type=rent"
              className={`px-3 py-1.5 rounded-lg text-xs font-black uppercase tracking-widest transition-all ${isRentActive ? 'bg-teal-600 text-white shadow-md shadow-teal-600/20' : 'opacity-60 hover:opacity-100 hover:bg-zinc-100 dark:hover:bg-zinc-800'}`}>
              🏠 {t.nav.rental ?? 'Rental'}
            </Link>
          </nav>

          {/* Search — hidden on home page */}
          <div className="flex-1 min-w-0 hidden md:block mx-2">
            {!isHomePage && <SearchBar />}
          </div>

          {/* Right actions */}
          <div className="flex items-center gap-1 ml-auto">

            {/* My Profile link for authenticated sellers */}
            {pathname.startsWith('/seller/') && !pathname.includes('/login') && !pathname.includes('/register') && (
              <Link href="/seller/profile"
                className="hidden sm:flex items-center gap-1.5 px-3 py-1.5 rounded-lg hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors text-xs font-bold text-zinc-600 dark:text-zinc-400 hover:text-orange-600 dark:hover:text-orange-500">
                <User size={14} />
                <span>My Profile</span>
              </Link>
            )}

            {/* Language */}
            <div className="relative">
              <button onClick={() => setLangOpen(v => !v)}
                className="flex items-center gap-1 px-2 py-1.5 rounded-full hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors text-xs font-black uppercase">
                <Globe size={14} />
                <span className="hidden sm:inline">{locale}</span>
                <ChevronDown size={10} className={`transition-transform ${langOpen ? 'rotate-180' : ''}`} />
              </button>
              {langOpen && (
                <div className="absolute right-0 mt-2 w-44 rounded-2xl border shadow-2xl z-50 overflow-hidden bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800">
                  {locales.map((l: string) => (
                    <button key={l} onClick={() => { setLocale(l); setLangOpen(false); }}
                      className={`w-full text-left px-4 py-2.5 text-sm font-bold hover:bg-orange-600 hover:text-white transition-colors flex items-center justify-between ${locale === l ? 'text-orange-600' : ''}`}>
                      {localeLabels[l] ?? l.toUpperCase()}
                      {locale === l && <span className="w-2 h-2 rounded-full bg-orange-600 inline-block" />}
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* Dark/Light */}
            <button onClick={toggleDark}
              className="p-2 rounded-full hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors">
              {mounted ? (dark ? <Sun size={16} /> : <Moon size={16} />) : <span className="w-4 h-4 block" />}
            </button>

            {/* Help */}
            <button onClick={() => setHelpOpen(true)}
              className="p-2 rounded-full hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors text-zinc-500 hover:text-orange-600 dark:text-zinc-400 dark:hover:text-orange-500">
              <HelpCircle size={16} />
            </button>

            {/* CTA — hidden on homepage (hero has the big cards) */}
            {!isHomePage && (
              <Link href={ctaHref}
                className="bg-orange-600 hover:bg-orange-700 text-white px-3 sm:px-4 py-1.5 sm:py-2 rounded-full font-bold text-xs sm:text-sm shadow-lg shadow-orange-600/20 transition-all active:scale-95 ml-1">
                {ctaLabel}
              </Link>
            )}
          </div>
        </div>

        {/* Mobile search — hidden on home page */}
        {!isHomePage && (
          <div className="md:hidden px-4 pb-3 bg-white/90 dark:bg-zinc-950/90 pb-[env(safe-area-inset-bottom,12px)]">
            <SearchBar />
          </div>
        )}
      </header>
    </>
  );
}
