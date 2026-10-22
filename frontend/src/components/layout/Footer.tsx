'use client';
import Link from 'next/link';
import { useLocale } from '@/hooks/useLocale';
import { Locale } from '@/lib/i18n';
import { ShoppingBag, Globe, Check, Shield, ChevronRight, Mail, Phone, ShieldCheck } from 'lucide-react';

const LOCALES: Locale[] = ['en', 'am', 'or', 'ti'];
const LOCALE_NAMES: Record<Locale, string> = { en: 'English', am: 'አማርኛ', or: 'Afaan Oromoo', ti: 'ትግርኛ' };

export function Footer() {
  const { t, locale, setLocale } = useLocale();

  return (
    <footer className="mt-16 border-t border-purple-200/40 dark:border-zinc-800 relative overflow-hidden">
      {/* Gradient background */}
      <div className="absolute inset-0 dark:hidden"
        style={{ background: 'linear-gradient(135deg, #ede9fe 0%, #ddd6fe 25%, #c7d2fe 60%, #bfdbfe 100%)' }} />
      <div className="absolute inset-0 hidden dark:block bg-zinc-950" />

      <div className="relative max-w-7xl mx-auto px-4 pt-12 pb-8">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mb-10">

          {/* Brand */}
          <div className="col-span-2 md:col-span-1 space-y-4">
            <div className="flex items-center gap-2.5">
              <div className="relative w-9 h-9 rounded-xl overflow-hidden flex items-center justify-center bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-md">
                <img src="/huntiq.png" alt="HUNTIQ" className="w-full h-full object-contain p-0.5" />
              </div>
              <span className="text-xl font-black text-zinc-900 dark:text-white">{t.ui.app_name}</span>
            </div>
            <p className="text-sm leading-relaxed text-zinc-600 dark:text-zinc-400">{t.footer.tagline}</p>
            <div className="flex gap-2">
              {[
                { href: 'https://www.facebook.com', label: 'f' },
                { href: 'https://www.twitter.com', label: '𝕏' },
                { href: 'https://www.instagram.com', label: '◎' },
                { href: 'https://www.youtube.com', label: '▶' },
              ].map(({ href, label }) => (
                <a key={href} href={href} target="_blank" rel="noopener noreferrer"
                  className="w-9 h-9 rounded-xl text-white flex items-center justify-center hover:opacity-90 transition-all font-bold text-sm shadow-md"
                  style={{ background: 'linear-gradient(135deg,#7c3aed,#3b82f6)' }}>
                  {label}
                </a>
              ))}
            </div>
          </div>

          {/* Quick Links — no admin, no developer */}
          <div className="space-y-3">
            <h4 className="text-xs font-black uppercase tracking-widest text-violet-700 dark:text-zinc-400">{t.footer.quick_links}</h4>
            <nav className="flex flex-col gap-2">
              {[
                { href: '/', label: t.footer.home },
                { href: '/products', label: t.nav.buy },
                { href: '/sell', label: t.nav.sell },
                { href: '/help', label: t.footer.help },
                { href: '/contact', label: t.footer.contact },
              ].map(({ href, label }) => (
                <Link key={href} href={href}
                  className="text-sm font-semibold text-zinc-700 dark:text-zinc-300 hover:text-violet-700 dark:hover:text-violet-400 transition-colors flex items-center gap-1.5 group">
                  <ChevronRight size={12} className="opacity-0 group-hover:opacity-100 transition-opacity text-violet-600" />
                  {label}
                </Link>
              ))}
            </nav>
          </div>

          {/* Legal only */}
          <div className="space-y-3">
            <h4 className="text-xs font-black uppercase tracking-widest text-violet-700 dark:text-zinc-400">{t.footer.legal}</h4>
            <nav className="flex flex-col gap-2">
              {[
                { href: '/terms', label: t.footer.terms },
                { href: '/privacy', label: t.footer.privacy },
              ].map(({ href, label }) => (
                <Link key={href} href={href}
                  className="text-sm font-semibold text-zinc-700 dark:text-zinc-300 hover:text-violet-700 dark:hover:text-violet-400 transition-colors flex items-center gap-1.5 group">
                  <ChevronRight size={12} className="opacity-0 group-hover:opacity-100 transition-opacity text-violet-600" />
                  {label}
                </Link>
              ))}
            </nav>
          </div>

          {/* Language */}
          <div className="space-y-3">
            <h4 className="text-xs font-black uppercase tracking-widest text-violet-700 dark:text-zinc-400">{t.ui.language}</h4>
            <div className="flex flex-col gap-1.5">
              {LOCALES.map(l => (
                <button key={l} onClick={() => setLocale(l)}
                  className={`flex items-center gap-2 px-3 py-2 rounded-xl border transition-all text-sm font-bold ${
                    locale === l
                      ? 'text-white border-transparent shadow-md'
                      : 'bg-white/60 dark:bg-zinc-900 border-violet-200 dark:border-zinc-700 text-zinc-700 dark:text-zinc-300 hover:border-violet-400 hover:text-violet-700'
                  }`}
                  style={locale === l ? { background: 'linear-gradient(135deg,#7c3aed,#3b82f6)', borderColor: 'transparent' } : {}}>
                  <Globe size={13} className={locale === l ? 'opacity-100' : 'opacity-50'} />
                  {LOCALE_NAMES[l]}
                  {locale === l && <Check size={12} className="ml-auto" />}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Bottom bar */}
        <div className="flex flex-col sm:flex-row items-center justify-between pt-6 border-t border-violet-200/60 dark:border-zinc-800 gap-3">
          <p className="text-xs font-semibold text-zinc-500 text-center sm:text-left">
            {t.footer.copyright}
          </p>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-1.5 text-xs font-bold text-zinc-500">
              <Shield size={12} className="text-violet-600" />{t.ui.secure_platform}
            </div>
            <div className="flex items-center gap-1.5 text-xs font-bold text-zinc-500">
              <Check size={12} className="text-blue-600" />{t.ui.verified_sellers}
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
