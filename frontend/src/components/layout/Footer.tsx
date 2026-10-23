'use client';
import Link from 'next/link';
import { useLocale } from '@/hooks/useLocale';
import { Locale } from '@/lib/i18n';
import {
  Globe,
  Check,
  Shield,
  ChevronRight,
  ShieldCheck,
  Store,
  ShoppingBag,
  HelpCircle,
  FileText,
  Lock,
} from 'lucide-react';

const LOCALES: Locale[] = ['en', 'am', 'or', 'ti'];
const LOCALE_NAMES: Record<Locale, string> = {
  en: 'English',
  am: 'አማርኛ',
  or: 'Afaan Oromoo',
  ti: 'ትግርኛ',
};

export function Footer() {
  const { t, locale, setLocale } = useLocale();

  return (
    <footer className="mt-16 border-t border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-950 transition-colors">
      <div className="max-w-7xl mx-auto px-4 pt-12 pb-8">
        {/* 5-Column Grid */}
        <div className="grid grid-cols-2 md:grid-cols-5 gap-8 mb-12">
          {/* Col 1: Brand & Identity */}
          <div className="col-span-2 md:col-span-1 space-y-4">
            <div className="flex items-center gap-2.5">
              <div className="relative w-9 h-9 rounded-xl overflow-hidden flex items-center justify-center bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-md">
                <img src="/huntiq.png" alt="SellBuy" className="w-full h-full object-contain p-0.5" />
              </div>
              <div className="flex flex-col">
                <span className="text-lg font-black text-zinc-900 dark:text-white leading-none">
                  {t.ui.app_name}
                </span>
                <span className="text-[10px] font-bold text-orange-600 uppercase tracking-wider">
                  Universal Marketplace
                </span>
              </div>
            </div>
            <p className="text-xs leading-relaxed text-zinc-500 dark:text-zinc-400">
              {t.footer.tagline ?? 'The premier universal marketplace for buying, selling, and renting across Ethiopia.'}
            </p>
            <div className="pt-2">
              <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-orange-600/10 text-orange-600 text-[11px] font-bold">
                <ShieldCheck size={13} />
                <span>Verified Platform</span>
              </div>
            </div>
          </div>

          {/* Col 2: Customer Support */}
          <div className="space-y-3">
            <h4 className="text-xs font-black uppercase tracking-widest text-zinc-900 dark:text-white flex items-center gap-1.5">
              <HelpCircle size={13} className="text-orange-600" />
              <span>Customer Support</span>
            </h4>
            <nav className="flex flex-col gap-2">
              {[
                { href: '/help', label: t.footer.help ?? 'Help Center' },
                { href: '/contact', label: t.footer.contact ?? 'Contact Us' },
                { href: '/help#faq', label: 'Frequently Asked Questions' },
                { href: '/contact', label: 'Support Ticket' },
              ].map(({ href, label }) => (
                <Link
                  key={label}
                  href={href}
                  className="text-xs font-semibold text-zinc-600 dark:text-zinc-400 hover:text-orange-600 dark:hover:text-orange-400 transition-colors flex items-center gap-1 group"
                >
                  <ChevronRight size={11} className="opacity-0 group-hover:opacity-100 transition-opacity text-orange-600" />
                  <span>{label}</span>
                </Link>
              ))}
            </nav>
          </div>

          {/* Col 3: Buying */}
          <div className="space-y-3">
            <h4 className="text-xs font-black uppercase tracking-widest text-zinc-900 dark:text-white flex items-center gap-1.5">
              <ShoppingBag size={13} className="text-orange-600" />
              <span>Buying</span>
            </h4>
            <nav className="flex flex-col gap-2">
              {[
                { href: '/products?type=sale', label: 'All Products for Sale' },
                { href: '/products?condition=new', label: 'Brand New Items' },
                { href: '/products?condition=used', label: 'Quality Pre-Owned' },
                { href: '/products?cat=property', label: 'Land & Property' },
                { href: '/help#buyer-protection', label: 'Buyer Protection' },
              ].map(({ href, label }) => (
                <Link
                  key={label}
                  href={href}
                  className="text-xs font-semibold text-zinc-600 dark:text-zinc-400 hover:text-orange-600 dark:hover:text-orange-400 transition-colors flex items-center gap-1 group"
                >
                  <ChevronRight size={11} className="opacity-0 group-hover:opacity-100 transition-opacity text-orange-600" />
                  <span>{label}</span>
                </Link>
              ))}
            </nav>
          </div>

          {/* Col 4: Selling & Renting */}
          <div className="space-y-3">
            <h4 className="text-xs font-black uppercase tracking-widest text-zinc-900 dark:text-white flex items-center gap-1.5">
              <Store size={13} className="text-orange-600" />
              <span>Selling & Renting</span>
            </h4>
            <nav className="flex flex-col gap-2">
              {[
                { href: '/sell', label: t.ui.start_selling ?? 'Start Selling' },
                { href: '/sell?type=rent', label: 'List Rental / Property' },
                { href: '/help#seller-guidelines', label: 'Seller Requirements' },
                { href: '/seller/login', label: 'Seller Dashboard' },
                { href: '/#tier-system', label: 'Gold & Premium Tiers' },
              ].map(({ href, label }) => (
                <Link
                  key={label}
                  href={href}
                  className="text-xs font-semibold text-zinc-600 dark:text-zinc-400 hover:text-orange-600 dark:hover:text-orange-400 transition-colors flex items-center gap-1 group"
                >
                  <ChevronRight size={11} className="opacity-0 group-hover:opacity-100 transition-opacity text-orange-600" />
                  <span>{label}</span>
                </Link>
              ))}
            </nav>
          </div>

          {/* Col 5: Policies, Safety & Languages */}
          <div className="space-y-3">
            <h4 className="text-xs font-black uppercase tracking-widest text-zinc-900 dark:text-white flex items-center gap-1.5">
              <FileText size={13} className="text-orange-600" />
              <span>Policies & Language</span>
            </h4>
            <nav className="flex flex-col gap-2 pb-2">
              {[
                { href: '/terms', label: t.footer.terms ?? 'Terms & Conditions' },
                { href: '/privacy', label: t.footer.privacy ?? 'Privacy Policy' },
                { href: '/help#safety', label: 'Safety Guidelines' },
              ].map(({ href, label }) => (
                <Link
                  key={label}
                  href={href}
                  className="text-xs font-semibold text-zinc-600 dark:text-zinc-400 hover:text-orange-600 dark:hover:text-orange-400 transition-colors flex items-center gap-1 group"
                >
                  <ChevronRight size={11} className="opacity-0 group-hover:opacity-100 transition-opacity text-orange-600" />
                  <span>{label}</span>
                </Link>
              ))}
            </nav>

            {/* Language Selection Buttons */}
            <div className="grid grid-cols-2 gap-1.5 pt-1">
              {LOCALES.map(l => (
                <button
                  key={l}
                  onClick={() => setLocale(l)}
                  className={`flex items-center justify-between px-2.5 py-1.5 rounded-lg border text-[11px] font-bold transition-all ${
                    locale === l
                      ? 'bg-orange-600 text-white border-orange-600 shadow-sm'
                      : 'bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 text-zinc-700 dark:text-zinc-300 hover:border-orange-400'
                  }`}
                >
                  <span>{LOCALE_NAMES[l]}</span>
                  {locale === l && <Check size={11} />}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Bottom Strip */}
        <div className="flex flex-col sm:flex-row items-center justify-between pt-6 border-t border-zinc-200 dark:border-zinc-800 gap-3">
          <p className="text-xs font-medium text-zinc-500 text-center sm:text-left">
            {t.footer.copyright ?? '© 2026 SellBuy. All rights reserved across Ethiopia.'}
          </p>
          <div className="flex items-center gap-4 text-xs font-bold text-zinc-500">
            <span className="flex items-center gap-1">
              <Shield size={12} className="text-green-600" />
              {t.ui.secure_platform_short ?? 'Protected Transactions'}
            </span>
            <span className="flex items-center gap-1">
              <Globe size={12} className="text-orange-600" />
              4 Ethiopian Languages
            </span>
          </div>
        </div>
      </div>
    </footer>
  );
}
