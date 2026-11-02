'use client';
import { useEffect, useState, useRef, useCallback } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import {
  ShoppingBag, Store, Tag, Shield, Zap, Globe,
  HelpCircle, X, ChevronRight,
  Trophy, TrendingUp, CheckCircle2,
  Package, ArrowRight,
  Gift, Users, Wallet, ChevronLeft,
} from 'lucide-react';
import { useLocale } from '@/hooks/useLocale';
import { api } from '@/lib/api';
import { ProductCard } from '@/components/public/ProductCard';
import { RentalProductCard } from '@/components/public/RentalProductCard';



interface Banner {
  id: number;
  media_type: 'image' | 'video';
  media_url: string;
  title: string;
  subtitle: string;
  link_url: string;
}

// ── Mode slide data ──────────────────────────────────────────────────────────
// Each mode has a background (image URL, video URL, or emoji fallback),
// accent colours, and CTA. The background transitions with a CSS crossfade.
interface SlideMode {
  key: 'sell' | 'buy';
  emoji: string;
  // image or video URL — leave empty to show emoji-only background
  mediaSrc: string;
  mediaType: 'image' | 'video' | 'emoji';
  gradientFrom: string; // tailwind-compatible inline style colour
  gradientTo: string;
  accentBg: string;     // button active bg
  accentText: string;   // button active text
  href: string;
  external?: boolean;
}

const SLIDE_MODES: SlideMode[] = [
  {
    key: 'sell',
    emoji: '🏪',
    mediaSrc: 'https://images.unsplash.com/photo-1556742031-c6961e8560b0?w=1600&q=85',
    mediaType: 'image',
    gradientFrom: '#1a0a00',
    gradientTo: '#7c2d12',
    accentBg: '#ea580c',
    accentText: '#ffffff',
    href: '/sell',
  },
  {
    key: 'buy',
    emoji: '🛍️',
    mediaSrc: 'https://images.unsplash.com/photo-1542838132-92c53300491e?w=1600&q=85',
    mediaType: 'image',
    gradientFrom: '#00091a',
    gradientTo: '#1e3a5f',
    accentBg: '#2563eb',
    accentText: '#ffffff',
    href: '/products',
  },
];

// ── Hero Slider ───────────────────────────────────────────────────────────────
function HeroSlider() {
  const { t } = useLocale();
  const [active, setActive] = useState(0);
  const [prev, setPrev] = useState<number | null>(null);
  const [transitioning, setTransitioning] = useState(false);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const nav = t.nav as any;
  const ui  = t.ui as any;

  const LABELS = [
    { title: nav.venn_sell_title ?? 'SELL / LANDLORD', sub: ui.hero_sell_sub ?? 'Sell products or rent properties as a landlord — reach thousands across Ethiopia.', btn: nav.venn_sell_btn ?? 'Start Listing →' },
    { title: nav.venn_buy_title  ?? 'BUY / RENT',    sub: ui.hero_buy_sub    ?? 'Browse thousands of verified listings — buy products or rent properties.', btn: nav.venn_buy_btn   ?? 'Browse Now →' },
  ];

  const goTo = useCallback((idx: number) => {
    if (idx === active || transitioning) return;
    setPrev(active);
    setActive(idx);
    setTransitioning(true);
    setTimeout(() => { setPrev(null); setTransitioning(false); }, 600);
  }, [active, transitioning]);

  // Auto-advance every 5 s
  useEffect(() => {
    timerRef.current = setTimeout(() => goTo((active + 1) % SLIDE_MODES.length), 5000);
    return () => { if (timerRef.current) clearTimeout(timerRef.current); };
  }, [active, goTo]);

  const mode  = SLIDE_MODES[active];
  const pMode = prev !== null ? SLIDE_MODES[prev] : null;

  return (
    <div className="relative w-full overflow-hidden"
      style={{ aspectRatio: '16/7', minHeight: '280px', maxHeight: '600px' }}>

      {/* ── Background layers (crossfade) ── */}
      {/* Previous slide fades out */}
      {pMode && (
        <div key={`prev-${prev}`} className="absolute inset-0 z-0"
          style={{ animation: 'heroFadeOut 0.6s ease forwards' }}>
          <SlideBackground mode={pMode} />
        </div>
      )}
      {/* Active slide fades in */}
      <div key={`active-${active}`} className="absolute inset-0 z-0"
        style={{ animation: 'heroFadeIn 0.6s ease forwards' }}>
        <SlideBackground mode={mode} />
      </div>

      {/* ── Gradient overlay ── */}
      <div className="absolute inset-0 z-10 pointer-events-none"
        style={{ background: `linear-gradient(to right, ${mode.gradientFrom}ee 0%, ${mode.gradientFrom}99 40%, transparent 100%)` }} />
      <div className="absolute inset-0 z-10 pointer-events-none"
        style={{ background: 'linear-gradient(to top, rgba(0,0,0,0.55) 0%, transparent 50%)' }} />

      {/* ── Content ── */}
      <div className="absolute inset-0 z-20 flex flex-col justify-between p-5 sm:p-8 lg:p-12">
        {/* Top: mode tabs */}
        <div className="flex items-center gap-2">
          {SLIDE_MODES.map((m, i) => {
            const label = i === 0 ? (nav.venn_sell_title ?? 'SELL / LIST') : (nav.venn_buy_title ?? 'BUY / RENT');
            const isActive = i === active;
            const Tag = m.external ? 'a' : Link;
            const linkProps = m.external
              ? { href: m.href, target: '_blank', rel: 'noopener noreferrer' }
              : { href: m.href };
            return (
              <Tag
                key={m.key}
                {...(linkProps as any)}
                onClick={(e: React.MouseEvent) => {
                  // If already on this slide, let the link navigate.
                  // If switching slides, just switch — don't navigate yet.
                  if (!isActive) { e.preventDefault(); goTo(i); }
                }}
                className="relative flex items-center gap-1.5 px-4 py-2 rounded-full text-xs sm:text-sm font-black uppercase tracking-wider transition-all duration-300 select-none cursor-pointer"
                style={{
                  background: isActive ? mode.accentBg : 'rgba(255,255,255,0.12)',
                  color: isActive ? mode.accentText : 'rgba(255,255,255,0.65)',
                  backdropFilter: 'blur(8px)',
                  border: isActive ? `2px solid ${mode.accentBg}` : '2px solid rgba(255,255,255,0.15)',
                  boxShadow: isActive ? `0 4px 20px ${mode.accentBg}55` : 'none',
                  transform: isActive ? 'scale(1.05)' : 'scale(1)',
                }}
              >
                <span>{SLIDE_MODES[i].emoji}</span>
                <span>{label}</span>
                {isActive && (
                  <span className="absolute -bottom-1 left-1/2 -translate-x-1/2 w-1 h-1 rounded-full bg-white" />
                )}
              </Tag>
            );
          })}
        </div>

        {/* Bottom: headline + CTA */}
        <div className="space-y-3 sm:space-y-4 max-w-lg">
          <h2 className="text-3xl sm:text-4xl lg:text-5xl xl:text-6xl font-black text-white leading-none tracking-tight drop-shadow-lg"
            style={{ textShadow: '0 2px 20px rgba(0,0,0,0.6)' }}>
            {LABELS[active].title}
          </h2>
          <p className="text-sm sm:text-base text-white/80 leading-relaxed max-w-sm drop-shadow"
            style={{ textShadow: '0 1px 8px rgba(0,0,0,0.8)' }}>
            {LABELS[active].sub}
          </p>
          {mode.external ? (
            <a href={mode.href} target="_blank" rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-5 py-2.5 sm:px-6 sm:py-3 rounded-xl font-black text-sm sm:text-base transition-all hover:-translate-y-0.5 shadow-lg"
              style={{ background: mode.accentBg, color: mode.accentText, boxShadow: `0 8px 24px ${mode.accentBg}55` }}>
              {LABELS[active].btn} <ArrowRight size={16} />
            </a>
          ) : (
            <Link href={mode.href}
              className="inline-flex items-center gap-2 px-5 py-2.5 sm:px-6 sm:py-3 rounded-xl font-black text-sm sm:text-base transition-all hover:-translate-y-0.5 shadow-lg"
              style={{ background: mode.accentBg, color: mode.accentText, boxShadow: `0 8px 24px ${mode.accentBg}55` }}>
              {LABELS[active].btn} <ArrowRight size={16} />
            </Link>
          )}
        </div>
      </div>

      {/* ── Prev / Next arrows ── */}
      <button onClick={() => goTo((active - 1 + SLIDE_MODES.length) % SLIDE_MODES.length)}
        className="absolute left-3 top-1/2 -translate-y-1/2 z-30 w-9 h-9 rounded-full flex items-center justify-center transition-all hover:scale-110"
        style={{ background: 'rgba(255,255,255,0.15)', backdropFilter: 'blur(8px)', border: '1px solid rgba(255,255,255,0.2)' }}
        aria-label="Previous">
        <ChevronLeft size={18} className="text-white" />
      </button>
      <button onClick={() => goTo((active + 1) % SLIDE_MODES.length)}
        className="absolute right-3 top-1/2 -translate-y-1/2 z-30 w-9 h-9 rounded-full flex items-center justify-center transition-all hover:scale-110"
        style={{ background: 'rgba(255,255,255,0.15)', backdropFilter: 'blur(8px)', border: '1px solid rgba(255,255,255,0.2)' }}
        aria-label="Next">
        <ChevronRight size={18} className="text-white" />
      </button>

      {/* ── Progress dots ── */}
      <div className="absolute bottom-4 right-5 z-30 flex items-center gap-1.5">
        {SLIDE_MODES.map((_, i) => (
          <button key={i} onClick={() => goTo(i)}
            className="rounded-full transition-all duration-300"
            style={{
              width: i === active ? '20px' : '6px',
              height: '6px',
              background: i === active ? mode.accentBg : 'rgba(255,255,255,0.35)',
            }}
            aria-label={`Go to slide ${i + 1}`}
          />
        ))}
      </div>

      {/* ── Auto-progress bar ── */}
      <div className="absolute bottom-0 left-0 right-0 z-30 h-0.5 bg-white/10">
        <div key={active} className="h-full rounded-full"
          style={{ background: mode.accentBg, animation: 'heroProgress 5s linear forwards' }} />
      </div>

      {/* Keyframes injected once */}
      <style>{`
        @keyframes heroFadeIn  { from { opacity: 0; } to { opacity: 1; } }
        @keyframes heroFadeOut { from { opacity: 1; } to { opacity: 0; } }
        @keyframes heroProgress { from { width: 0%; } to { width: 100%; } }
      `}</style>
    </div>
  );
}

// Renders the background layer for a single slide mode
function SlideBackground({ mode }: { mode: SlideMode }) {
  if (mode.mediaType === 'video') {
    return (
      <video
        src={mode.mediaSrc}
        autoPlay muted loop playsInline
        className="absolute inset-0 w-full h-full object-cover"
      />
    );
  }
  if (mode.mediaType === 'image' && mode.mediaSrc) {
    return (
      // eslint-disable-next-line @next/next/no-img-element
      <img
        src={mode.mediaSrc}
        alt=""
        className="absolute inset-0 w-full h-full object-cover"
        loading="eager"
      />
    );
  }
  // Emoji fallback
  return (
    <div className="absolute inset-0 flex items-center justify-center text-[20vw] select-none"
      style={{ background: `linear-gradient(135deg, ${mode.gradientFrom}, ${mode.gradientTo})` }}>
      {mode.emoji}
    </div>
  );
}

function HelpPanel({ onClose }: { onClose: () => void }) {
  const { t } = useLocale();
  const h = t.sellerHelp ?? {};
  const [openFaq, setOpenFaq] = useState<number | null>(null);

  const STEP_COLORS = [
    { bg: 'from-orange-500 to-amber-500',   ring: 'ring-orange-400/40', light: 'bg-orange-50 dark:bg-orange-900/20',  text: 'text-orange-600 dark:text-orange-400' },
    { bg: 'from-blue-500 to-cyan-500',      ring: 'ring-blue-400/40',   light: 'bg-blue-50 dark:bg-blue-900/20',      text: 'text-blue-600 dark:text-blue-400' },
    { bg: 'from-violet-500 to-purple-600',  ring: 'ring-violet-400/40', light: 'bg-violet-50 dark:bg-violet-900/20',  text: 'text-violet-600 dark:text-violet-400' },
    { bg: 'from-green-500 to-emerald-500',  ring: 'ring-green-400/40',  light: 'bg-green-50 dark:bg-green-900/20',    text: 'text-green-600 dark:text-green-400' },
  ];
  const steps = [
    { emoji: '👤', title: h.step1Title ?? 'Register & Login',  desc: h.step1Desc  ?? 'Create your free account in seconds.' },
    { emoji: '📸', title: h.step2Title ?? 'Upload Photos',     desc: h.step2Desc  ?? 'Add 2–4 clear, well-lit product photos.' },
    { emoji: '📍', title: h.step3Title ?? 'Set Location',      desc: h.step3Desc  ?? 'Pick your region so buyers can find you.' },
    { emoji: '✅', title: h.step4Title ?? 'Go Live',           desc: h.step4Desc  ?? 'Admin reviews and publishes within 24h.' },
  ];
  const faqs = [
    { q: h.faq1q, a: h.faq1a }, { q: h.faq2q, a: h.faq2a },
    { q: h.faq3q, a: h.faq3a }, { q: h.faq4q, a: h.faq4a },
    { q: h.faq5q, a: h.faq5a },
  ].filter(f => f.q);

  return (
    <div className="fixed inset-0 z-[60] flex items-end sm:items-center justify-center p-0 sm:p-4 lg:p-8">
      <div className="absolute inset-0 bg-black/70 backdrop-blur-md" onClick={onClose} />
      <div className="relative z-10 w-full sm:max-w-lg lg:max-w-2xl bg-white dark:bg-zinc-950 rounded-t-[2rem] sm:rounded-[2rem] shadow-[0_32px_80px_-12px_rgba(0,0,0,0.5)] border border-zinc-200/60 dark:border-zinc-800 max-h-[92vh] sm:max-h-[85vh] flex flex-col overflow-hidden">
        <div className="relative shrink-0 bg-gradient-to-br from-orange-500 via-orange-600 to-red-600 px-6 py-6 sm:py-8 overflow-hidden">
          <div className="absolute -top-8 -right-8 w-32 h-32 rounded-full bg-white/10 blur-2xl" />
          <div className="absolute -bottom-6 -left-6 w-24 h-24 rounded-full bg-yellow-300/20 blur-xl" />
          <div className="relative flex items-start justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 sm:w-14 sm:h-14 rounded-2xl bg-white/20 backdrop-blur-sm ring-2 ring-white/30 flex items-center justify-center shadow-xl">
                <HelpCircle size={22} className="text-white" />
              </div>
              <div>
                <p className="font-black text-white text-lg sm:text-xl leading-tight">{h.title ?? 'Help Center'}</p>
                <p className="text-white/70 text-xs sm:text-sm mt-0.5">{h.subtitle ?? 'How to use HUNTIQ'}</p>
              </div>
            </div>
            <button onClick={onClose} className="shrink-0 w-9 h-9 rounded-xl bg-white/20 hover:bg-white/30 text-white flex items-center justify-center transition-colors backdrop-blur-sm">
              <X size={16} />
            </button>
          </div>
          <div className="relative flex flex-wrap gap-2 mt-5">
            {[
              { icon: '🚀', label: t.ui.feat_free_to_list ?? 'Free to list' },
              { icon: '⚡', label: t.ui.feat_fast_approval ?? 'Fast approval' },
              { icon: '🔒', label: t.ui.feat_secure ?? 'Secure' },
              { icon: '🌍', label: t.ui.feat_4_languages ?? '4 languages' },
            ].map(({ icon, label }) => (
              <div key={label} className="flex items-center gap-1.5 bg-white/15 backdrop-blur-sm text-white text-[10px] sm:text-xs font-bold px-2.5 py-1.5 rounded-full border border-white/20">
                <span>{icon}</span>{label}
              </div>
            ))}
          </div>
        </div>
        <div className="overflow-y-auto flex-1 px-5 sm:px-6 py-5 sm:py-6 space-y-6">
          <div>
            <div className="flex items-center gap-2 mb-4">
              <span className="text-base">🗺️</span>
              <p className="text-xs font-black uppercase tracking-widest text-zinc-400">{h.howItWorks ?? 'How It Works'}</p>
            </div>
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
              {steps.map(({ emoji, title, desc }, i) => {
                const c = STEP_COLORS[i];
                return (
                  <div key={i} className={`relative rounded-2xl p-4 ${c.light} border border-white dark:border-zinc-800 overflow-hidden`}>
                    <span className="absolute -top-2 -right-1 text-6xl font-black opacity-[0.06] select-none leading-none">{i + 1}</span>
                    <div className={`w-10 h-10 rounded-xl bg-gradient-to-br ${c.bg} ring-2 ${c.ring} flex items-center justify-center text-lg shadow-lg mb-3`}>{emoji}</div>
                    <p className={`text-xs font-black leading-tight mb-1 ${c.text}`}>{String(i + 1).padStart(2, '0')} — {title}</p>
                    <p className="text-[10px] sm:text-xs text-zinc-500 dark:text-zinc-400 leading-relaxed">{desc}</p>
                  </div>
                );
              })}
            </div>
          </div>
          <div className="rounded-2xl overflow-hidden border border-orange-200 dark:border-orange-800/50">
            <div className="bg-gradient-to-r from-orange-500/10 to-amber-500/5 dark:from-orange-900/30 dark:to-amber-900/10 px-4 py-3 flex items-center gap-2 border-b border-orange-200 dark:border-orange-800/50">
              <span className="text-base">📷</span>
              <p className="text-xs font-black uppercase tracking-widest text-orange-700 dark:text-orange-400">{h.imageRules ?? 'Image Rules'}</p>
            </div>
            <div className="px-4 py-3 space-y-2 bg-white dark:bg-zinc-900/50">
              {[h.imageRule1 ?? '📌 Minimum 2 images, maximum 4 per listing.', h.imageRule2 ?? '💡 Use clear, well-lit photos — no blurry shots.', h.imageRule3 ?? '🚫 No watermarks, text overlays, or collages.'].map((r, i) => (
                <div key={i} className="flex items-start gap-2.5 text-xs text-zinc-600 dark:text-zinc-300">
                  <CheckCircle2 size={13} className="text-orange-500 mt-0.5 shrink-0" /><span>{r}</span>
                </div>
              ))}
            </div>
          </div>
          {faqs.length > 0 && (
            <div>
              <div className="flex items-center gap-2 mb-3"><span className="text-base">💬</span><p className="text-xs font-black uppercase tracking-widest text-zinc-400">{h.faq ?? 'FAQ'}</p></div>
              <div className="space-y-2">
                {faqs.map((faq, i) => (
                  <div key={i} className="rounded-2xl border border-zinc-200 dark:border-zinc-800 overflow-hidden bg-white dark:bg-zinc-900/60">
                    <button onClick={() => setOpenFaq(openFaq === i ? null : i)} className="w-full flex items-center justify-between px-4 py-3 text-left hover:bg-zinc-50 dark:hover:bg-zinc-800/60 transition-colors gap-3">
                      <span className="text-xs sm:text-sm font-semibold text-zinc-800 dark:text-zinc-200 leading-snug">{faq.q}</span>
                      <div className={`shrink-0 w-6 h-6 rounded-full flex items-center justify-center transition-all ${openFaq === i ? 'bg-orange-600 text-white rotate-90' : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-400'}`}><ChevronRight size={12} /></div>
                    </button>
                    {openFaq === i && <div className="px-4 pb-4 pt-1 text-xs sm:text-sm text-zinc-500 dark:text-zinc-400 leading-relaxed border-t border-zinc-100 dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-800/30">{faq.a}</div>}
                  </div>
                ))}
              </div>
            </div>
          )}
          <div className="rounded-2xl bg-gradient-to-br from-zinc-900 to-zinc-800 dark:from-zinc-800 dark:to-zinc-900 p-5 flex flex-col sm:flex-row items-center gap-4 border border-zinc-700">
            <div className="flex-1 text-center sm:text-left">
              <p className="text-white font-black text-sm sm:text-base">{t.ui.ready_to_sell ?? 'Ready to start selling?'}</p>
              <p className="text-zinc-400 text-xs mt-0.5">{t.ui.join_sellers ?? 'Join thousands of sellers on HUNTIQ'}</p>
            </div>
            <Link href="/sell" onClick={onClose} className="shrink-0 flex items-center gap-2 bg-orange-600 hover:bg-orange-500 text-white font-black text-sm px-5 py-2.5 rounded-xl transition-colors shadow-lg shadow-orange-600/30">
              <Store size={15} /> {t.ui.start_selling ?? 'Start Selling'}
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

// ── Main Hero ────────────────────────────────────────────────────────────────
export function HomeHero() {
  const { t, locale } = useLocale();
  const [rootCategories, setRootCategories] = useState<any[]>([]);
  const [loadingCategories, setLoadingCategories] = useState(true);
  const [categoriesError, setCategoriesError] = useState<string | null>(null);
  const [settings, setSettings] = useState<any>(null);
  const [recentProducts, setRecentProducts] = useState<any[]>([]);
  const [rentalProducts, setRentalProducts] = useState<any[]>([]);
  const [conditionFilter, setConditionFilter] = useState<'all' | 'new' | 'used'>('all');

  useEffect(() => {
    let active = true;
    const fetchCats = async () => {
      try {
        const cats = await api.getRootCategories(locale);
        if (!active) return;
        if (Array.isArray(cats) && cats.length > 0) {
          setRootCategories(cats);
        } else {
          const fallback = await api.getCategories();
          if (active) setRootCategories(Array.isArray(fallback) ? fallback : []);
        }
      } catch {
        if (active) setCategoriesError('Unable to load categories');
      } finally {
        if (active) setLoadingCategories(false);
      }
    };
    fetchCats();
    return () => { active = false; };
  }, [locale]);

  useEffect(() => {
    api.getPublicPaymentSettings().then(setSettings).catch(() => {});
    // Load recent approved listings for homepage preview
    api.getProducts({ limit: 8, page: 1 })
      .then(data => setRecentProducts(Array.isArray(data) ? data : []))
      .catch(() => {});

    // Load active rental and leasing opportunities
    api.getProducts({ listing_type: 'rent', limit: 4, page: 1 })
      .then(data => setRentalProducts(Array.isArray(data) ? data : []))
      .catch(() => {});
  }, []);

  const filteredRecent = recentProducts.filter(p => {
    if (conditionFilter === 'all') return true;
    if (conditionFilter === 'new') return p.condition === 'new';
    return p.condition !== 'new';
  });

  return (
    <>
      {/* ── HERO: full-width slider, no circles ── */}
      <section className="w-full">
        <HeroSlider />
      </section>

      {/* ── Everything below the fold ── */}
      <section className="relative w-full overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-br from-violet-50/80 via-indigo-50/60 to-blue-50/40 dark:from-zinc-950 dark:via-zinc-900 dark:to-zinc-950 pointer-events-none" />

        <div className="relative z-10 w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 flex flex-col items-center gap-12 lg:gap-16">

          {/* Stats bar — values from real backend */}
          <div className="w-full max-w-3xl">
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 sm:gap-4">
              {[
                { v: recentProducts.length > 0 ? `${recentProducts.length}+` : '100+', l: t.ui.stat_products ?? 'Products', icon: <Package size={18} />, color: 'text-orange-600', bg: 'bg-orange-50 dark:bg-orange-900/20' },
                { v: '23+', l: t.ui.stat_categories ?? 'Categories', icon: <Tag size={18} />, color: 'text-blue-600', bg: 'bg-blue-50 dark:bg-blue-900/20' },
                { v: '4', l: t.ui.stat_languages ?? 'Languages', icon: <Globe size={18} />, color: 'text-purple-600', bg: 'bg-purple-50 dark:bg-purple-900/20' },
                { v: '100%', l: t.ui.stat_verified ?? 'Verified', icon: <CheckCircle2 size={18} />, color: 'text-green-600', bg: 'bg-green-50 dark:bg-green-900/20' },
              ].map(({ v, l, icon, color, bg }) => (
                <div key={l} className={`flex flex-col items-center gap-1.5 sm:gap-2 py-3.5 sm:py-4 px-2 rounded-2xl border border-zinc-200/60 dark:border-zinc-800 shadow-sm ${bg} transition-all hover:-translate-y-0.5 hover:shadow-md`}>
                  <span className={color}>{icon}</span>
                  <span className={`text-xl sm:text-2xl lg:text-3xl font-black leading-none ${color}`}>{v}</span>
                  <span className="text-[10px] sm:text-[10px] font-bold uppercase tracking-widest text-zinc-500 dark:text-zinc-400 text-center">{l}</span>
                </div>
              ))}
            </div>
          </div>

          {/* Trust strip */}
          <div className="flex flex-wrap items-center justify-center gap-3 sm:gap-5 lg:gap-8 text-xs sm:text-sm text-zinc-500 dark:text-zinc-400">
            <span className="flex items-center gap-1.5 font-medium"><Shield size={15} className="text-green-500" /> {t.ui.secure_platform_short}</span>
            <span className="w-1 h-1 rounded-full bg-zinc-300 dark:bg-zinc-700 hidden sm:block" />
            <span className="flex items-center gap-1.5 font-medium"><Globe size={15} className="text-orange-500" /> {t.ui.four_languages}</span>
            <span className="w-1 h-1 rounded-full bg-zinc-300 dark:bg-zinc-700 hidden sm:block" />
            <span className="flex items-center gap-1.5 font-medium"><Zap size={15} className="text-blue-500" /> {t.ui.fast_and_free}</span>
            <span className="w-1 h-1 rounded-full bg-zinc-300 dark:bg-zinc-700 hidden sm:block" />
            <span className="flex items-center gap-1.5 font-medium"><Trophy size={15} className="text-yellow-500" /> {t.ui.tier_rewards ?? 'Tier Rewards'}</span>
          </div>

          {/* Category quick links */}
          <div className="w-full">
            <div className="flex items-center justify-between mb-5">
              <div className="flex items-center gap-2">
                <div className="w-6 h-6 rounded-lg bg-orange-600/10 flex items-center justify-center"><Tag size={13} className="text-orange-600" /></div>
                <span className="text-xs font-black uppercase tracking-widest text-zinc-400">{t.ui.browse_by_category}</span>
              </div>
              <Link
                href="/products"
                className="text-xs font-bold text-orange-600 hover:underline flex items-center gap-1 tap-bounce"
              >
                <span>Browse All</span> <ArrowRight size={11} />
              </Link>
            </div>
            {loadingCategories ? (
              <div className="grid grid-cols-4 sm:grid-cols-8 gap-3">
                {Array.from({ length: 8 }).map((_, i) => (
                  <div key={i} className="animate-pulse flex flex-col items-center gap-2 p-3 sm:p-4 rounded-2xl bg-zinc-100 dark:bg-zinc-800/60 border border-zinc-200/50 dark:border-zinc-800">
                    <div className="w-8 h-8 rounded-full bg-zinc-200 dark:bg-zinc-700" />
                    <div className="w-12 h-2.5 rounded bg-zinc-200 dark:bg-zinc-700" />
                  </div>
                ))}
              </div>
            ) : categoriesError && rootCategories.length === 0 ? (
              <div className="text-center py-6 border border-dashed border-zinc-200 dark:border-zinc-800 rounded-2xl">
                <p className="text-xs text-zinc-500">{categoriesError}</p>
                <button
                  onClick={() => {
                    setLoadingCategories(true);
                    api.getRootCategories(locale)
                      .then(setRootCategories)
                      .catch(() => setCategoriesError('Unable to load categories'))
                      .finally(() => setLoadingCategories(false));
                  }}
                  className="mt-2 text-xs font-bold text-orange-600 hover:underline tap-bounce"
                >
                  Retry
                </button>
              </div>
            ) : (
              <div className="grid grid-cols-4 sm:grid-cols-8 gap-2.5 sm:gap-3">
                {rootCategories.slice(0, 16).map((cat) => {
                  const label =
                    (cat.slug && (t.categories as any)?.[cat.slug]) ||
                    (locale === 'am' && (cat.name_am || (cat as any).NameAm)) ||
                    ((locale === 'or' || locale === 'om') && (cat.name_om || cat.name_or || (cat as any).NameOm)) ||
                    (locale === 'ti' && (cat.name_ti || (cat as any).NameTi)) ||
                    (locale === 'en' && (cat.name_en || (cat as any).NameEn)) ||
                    cat.name ||
                    cat.name_en ||
                    cat.slug?.replace(/_/g, ' ');
                  const icon = cat.icon || (cat.image_url ? (
                    <Image
                      src={cat.image_url}
                      alt={label || 'Category'}
                      width={32}
                      height={32}
                      className="w-7 h-7 sm:w-8 sm:h-8 object-contain"
                      unoptimized
                    />
                  ) : '📦');
                  return (
                    <Link
                      key={cat.id || cat.slug}
                      href={`/products?category_hierarchy_id=${cat.id}&cat=${cat.slug}`}
                      className="group flex flex-col items-center gap-2 p-2.5 sm:p-4 rounded-2xl bg-white dark:bg-zinc-900 border border-zinc-100 dark:border-zinc-800 hover:border-orange-400 hover:shadow-xl hover:shadow-orange-600/10 hover:-translate-y-1 transition-all duration-200 tap-bounce active:scale-95"
                    >
                      <span className="text-2xl sm:text-3xl group-hover:scale-125 transition-transform duration-200 flex items-center justify-center min-h-[36px]">
                        {typeof icon === 'string' ? icon : icon}
                      </span>
                      <span className="text-[10px] sm:text-xs font-bold text-zinc-700 dark:text-zinc-300 text-center leading-tight line-clamp-1 capitalize">
                        {label}
                      </span>
                    </Link>
                  );
                })}
              </div>
            )}
          </div>

          {/* Featured Marketplace Listings (with condition tabs) */}
          <div className="w-full">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-6">
              <div className="flex items-center gap-2">
                <div className="w-7 h-7 rounded-lg bg-green-600/10 flex items-center justify-center">
                  <TrendingUp size={15} className="text-green-600" />
                </div>
                <div>
                  <h3 className="text-base sm:text-lg font-black tracking-tight text-zinc-900 dark:text-white leading-none">
                    {(t.ui as any).latest_listings ?? t.ui.recent_products ?? 'Featured Listings'}
                  </h3>
                  <p className="text-[11px] text-zinc-400 mt-0.5">
                    Verified items for sale across Ethiopia
                  </p>
                </div>
              </div>

              {/* Condition Filter Pills */}
              <div className="flex items-center gap-1 bg-zinc-100 dark:bg-zinc-800 p-1 rounded-xl text-xs font-bold self-start sm:self-auto">
                <button
                  onClick={() => setConditionFilter('all')}
                  className={`px-3 py-1 rounded-lg transition-all tap-bounce active:scale-95 ${
                    conditionFilter === 'all'
                      ? 'bg-white dark:bg-zinc-900 text-orange-600 shadow-sm'
                      : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'
                  }`}
                >
                  All Items
                </button>
                <button
                  onClick={() => setConditionFilter('new')}
                  className={`px-3 py-1 rounded-lg transition-all tap-bounce active:scale-95 ${
                    conditionFilter === 'new'
                      ? 'bg-white dark:bg-zinc-900 text-orange-600 shadow-sm'
                      : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'
                  }`}
                >
                  Brand New
                </button>
                <button
                  onClick={() => setConditionFilter('used')}
                  className={`px-3 py-1 rounded-lg transition-all tap-bounce active:scale-95 ${
                    conditionFilter === 'used'
                      ? 'bg-white dark:bg-zinc-900 text-orange-600 shadow-sm'
                      : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-white'
                  }`}
                >
                  Pre-Owned
                </button>
              </div>
            </div>

            {filteredRecent.length > 0 ? (
              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-4 gap-3 sm:gap-4">
                {filteredRecent.map(product => (
                  <ProductCard key={product.id} product={product} />
                ))}
              </div>
            ) : (
              <div className="text-center py-12 px-4 rounded-3xl border border-dashed border-zinc-200 dark:border-zinc-800 bg-white/50 dark:bg-zinc-900/30">
                <Package size={36} className="mx-auto text-zinc-300 dark:text-zinc-600 mb-2" />
                <p className="text-sm font-bold text-zinc-600 dark:text-zinc-400">
                  No listings found for this condition
                </p>
                <Link
                  href="/products"
                  className="mt-3 inline-flex items-center gap-1.5 px-4 py-2 rounded-xl text-xs font-bold bg-orange-600 text-white hover:bg-orange-700 transition-colors"
                >
                  <span>Explore All Products</span>
                  <ArrowRight size={12} />
                </Link>
              </div>
            )}
          </div>

          {/* Rental & Leasing Section (Real Data from API) */}
          {rentalProducts.length > 0 && (
            <div className="w-full pt-6 border-t border-zinc-200/60 dark:border-zinc-800">
              <div className="flex items-center justify-between mb-6">
                <div className="flex items-center gap-2.5">
                  <div className="w-8 h-8 rounded-xl bg-teal-600/10 flex items-center justify-center">
                    <span className="text-base">🏠</span>
                  </div>
                  <div>
                    <h3 className="text-base sm:text-lg font-black tracking-tight text-zinc-900 dark:text-white leading-none">
                      Rental & Leasing Opportunities
                    </h3>
                    <p className="text-[11px] text-zinc-400 mt-0.5">
                      Homes, apartments, commercial spaces, and vehicles for rent
                    </p>
                  </div>
                </div>
                <Link
                  href="/products?type=rent"
                  className="text-xs font-bold text-teal-600 hover:underline flex items-center gap-1"
                >
                  <span>Browse All Rentals</span> <ArrowRight size={11} />
                </Link>
              </div>
              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 sm:gap-4">
                {rentalProducts.map(product => (
                  <RentalProductCard key={product.id} product={product} />
                ))}
              </div>
            </div>
          )}

          {/* Seller & Landlord Onboarding Section */}
          <div className="w-full rounded-3xl bg-gradient-to-br from-orange-600 via-orange-500 to-amber-600 text-white p-6 sm:p-10 shadow-xl overflow-hidden relative">
            <div className="absolute top-0 right-0 w-80 h-80 bg-white/10 rounded-full blur-3xl pointer-events-none" />
            <div className="relative z-10 max-w-2xl space-y-4">
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/20 backdrop-blur-sm text-xs font-black uppercase tracking-wider text-white">
                <Store size={13} />
                <span>Join HUNTIQ Marketplace</span>
              </div>
              <h3 className="text-2xl sm:text-3xl lg:text-4xl font-black leading-tight">
                Start Selling or Renting Properties Across Ethiopia
              </h3>
              <p className="text-sm sm:text-base text-white/90 leading-relaxed">
                Connect with thousands of buyers and tenants in Addis Ababa and all regional cities. Fast listing verification and direct buyer calls.
              </p>

              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 pt-2 pb-2">
                {[
                  { step: '01', title: 'Free Account', desc: 'Register in seconds with your Ethiopian phone number.' },
                  { step: '02', title: 'Post Listings', desc: 'Add 2–4 clear photos, set price, and pick your location.' },
                  { step: '03', title: 'Connect & Earn', desc: 'Fast review within 24h — receive orders directly.' },
                ].map(({ step, title, desc }) => (
                  <div key={step} className="bg-white/15 backdrop-blur-sm rounded-2xl p-3 border border-white/20">
                    <span className="text-[10px] font-black text-amber-200">{step}</span>
                    <p className="text-xs font-black text-white mt-0.5">{title}</p>
                    <p className="text-[11px] text-white/80 leading-relaxed mt-0.5">{desc}</p>
                  </div>
                ))}
              </div>

              <div className="pt-2 flex flex-wrap items-center gap-3">
                <Link
                  href="/sell"
                  className="px-6 py-3 rounded-xl bg-white text-orange-600 font-black text-xs sm:text-sm shadow-lg hover:bg-orange-50 transition-all hover:scale-[1.02] active:scale-95 flex items-center gap-2"
                >
                  <Store size={16} />
                  <span>List An Item Or Property</span>
                </Link>
                <Link
                  href="/help"
                  className="px-5 py-3 rounded-xl bg-white/20 hover:bg-white/30 text-white font-bold text-xs sm:text-sm backdrop-blur-sm transition-all flex items-center gap-1.5"
                >
                  <HelpCircle size={15} />
                  <span>How It Works</span>
                </Link>
              </div>
            </div>
          </div>

          {/* Tier System */}
          <div className="w-full pt-12 border-t border-yellow-200/60 dark:border-zinc-800">
            <div className="text-center mb-10">
              <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full text-xs font-black uppercase tracking-widest border mb-4"
                style={{ background: 'linear-gradient(90deg,#fef08a,#e9d5ff)', borderColor: '#d8b4fe', color: '#7c3aed' }}>
                <Trophy size={13} />{t.ui.dual_tier_system ?? 'Tier System'}
              </div>
              <h2 className="text-3xl lg:text-4xl xl:text-5xl font-black mb-3">
                <span className="tier-gradient-text">{t.ui.dual_tier_system ?? 'Dual Tier System'}</span>
              </h2>
              <p className="text-base lg:text-lg text-zinc-500 dark:text-zinc-400 max-w-xl mx-auto">
                {t.ui.dual_tier_sub ?? 'Separate tiers for sellers and buyers — upgrade independently'}
              </p>
            </div>
            <div className="grid md:grid-cols-2 gap-6 lg:gap-8">
              {/* Seller tiers */}
              <div className="relative bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 lg:p-8 shadow-xl overflow-hidden">
                <div className="absolute top-0 left-0 right-0 h-1.5 rounded-t-3xl" style={{ background: 'linear-gradient(90deg,#eab308,#f97316,#ea580c)' }} />
                <div className="absolute top-0 right-0 w-40 h-40 bg-yellow-400/8 rounded-full blur-2xl" />
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-12 h-12 rounded-2xl flex items-center justify-center shadow-lg" style={{ background: 'linear-gradient(135deg,#eab308,#ea580c)' }}><Store size={20} className="text-white" /></div>
                  <div><h3 className="font-black text-lg">{t.ui.seller_tiers ?? 'Seller Tiers'}</h3><p className="text-xs text-zinc-400">{t.ui.list_and_earn ?? 'List & earn more'}</p></div>
                </div>
                <div className="space-y-3">
                  {[
                    { tier: t.ui.tier_normal ?? 'Normal',  icon: '🥉', req: t.ui.tier_normal_req ?? 'Any seller',            fee: settings ? `${settings.normal_startup_fee}% startup + ${settings.normal_weekly_increment}%/week` : (t.ui.tier_normal_fee ?? '2% startup + 0.25%/week'),  color: 'from-zinc-100 to-zinc-50 dark:from-zinc-800 dark:to-zinc-800/50 border-zinc-200 dark:border-zinc-700',  badge: 'bg-zinc-200 dark:bg-zinc-700 text-zinc-600 dark:text-zinc-300',  href: '/seller/login', btnLabel: t.ui.register_as_seller ?? 'Register as Seller', btnStyle: 'bg-orange-600 hover:bg-orange-700 text-white' },
                    { tier: t.ui.tier_gold ?? 'Gold',      icon: '🥇', req: t.ui.tier_gold_req ?? '≥ ETB 100,000 sold',      fee: settings ? `${settings.gold_startup_fee}% startup + ${settings.gold_weekly_increment}%/week` : (t.ui.tier_gold_fee ?? '1.5% startup + 0.15%/week'),   color: 'from-yellow-50 to-amber-50 dark:from-yellow-900/20 dark:to-amber-900/10 border-yellow-200 dark:border-yellow-800/50', badge: 'bg-yellow-100 dark:bg-yellow-900/40 text-yellow-700 dark:text-yellow-400', href: '/seller/login', btnLabel: t.ui.start_selling ?? 'Start Selling', btnStyle: 'bg-yellow-500 hover:bg-yellow-600 text-white' },
                    { tier: t.ui.tier_premium ?? 'Premium',icon: '👑', req: t.ui.tier_premium_req ?? '≥ ETB 5,000,000 sold', fee: settings ? `${settings.premium_startup_fee}% startup + ${settings.premium_weekly_increment}%/week` : (t.ui.tier_premium_fee ?? '1.2% startup + 0.10%/week'), color: 'from-purple-50 to-violet-50 dark:from-purple-900/20 dark:to-violet-900/10 border-purple-200 dark:border-purple-800/50', badge: 'bg-purple-100 dark:bg-purple-900/40 text-purple-700 dark:text-purple-400', href: '/seller/login', btnLabel: t.ui.start_selling ?? 'Start Selling', btnStyle: 'bg-purple-600 hover:bg-purple-700 text-white' },
                  ].map(({ tier, icon, req, fee, color, badge, href, btnLabel, btnStyle }) => (
                    <div key={tier} className={`flex items-start gap-3 p-4 bg-gradient-to-r ${color} rounded-2xl border`}>
                      <span className="text-xl shrink-0 mt-0.5">{icon}</span>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-1"><span className={`text-[10px] font-black px-2 py-0.5 rounded-full uppercase ${badge}`}>{tier}</span></div>
                        <p className="text-sm font-bold text-zinc-700 dark:text-zinc-300">{req}</p>
                        <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-0.5">{fee}</p>
                        <Link href={href} className={`inline-flex items-center gap-1 mt-2.5 px-3 py-1.5 rounded-xl text-xs font-black transition-all ${btnStyle}`}>{btnLabel} <ArrowRight size={11} /></Link>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Buyer tiers */}
              <div className="relative bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 lg:p-8 shadow-xl overflow-hidden">
                <div className="absolute top-0 left-0 right-0 h-1.5 rounded-t-3xl" style={{ background: 'linear-gradient(90deg,#a855f7,#6366f1,#3b82f6)' }} />
                <div className="absolute top-0 right-0 w-40 h-40 bg-purple-400/8 rounded-full blur-2xl" />
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-12 h-12 rounded-2xl flex items-center justify-center shadow-lg" style={{ background: 'linear-gradient(135deg,#a855f7,#6366f1)' }}><ShoppingBag size={20} className="text-white" /></div>
                  <div><h3 className="font-black text-lg">{t.ui.buyer_tiers ?? 'Buyer Tiers'}</h3><p className="text-xs text-zinc-400">{t.ui.unlock_access ?? 'Unlock more access'}</p></div>
                </div>
                <div className="space-y-3">
                  {[
                    { tier: t.ui.tier_normal ?? 'Normal',  icon: '🥉', cost: t.ui.tier_buyer_normal_cost ?? 'Free',        access: t.ui.tier_buyer_normal_access ?? 'View Normal seller images',            color: 'from-zinc-100 to-zinc-50 dark:from-zinc-800 dark:to-zinc-800/50 border-zinc-200 dark:border-zinc-700',  badge: 'bg-zinc-200 dark:bg-zinc-700 text-zinc-600 dark:text-zinc-300',  href: '/products',        btnLabel: t.ui.browse_products ?? 'Browse Products',      btnStyle: 'bg-orange-600 hover:bg-orange-700 text-white' },
                    { tier: t.ui.tier_gold ?? 'Gold',      icon: '🥇', cost: settings ? `ETB ${settings.buyer_gold_upgrade_price}` : (t.ui.tier_buyer_gold_cost ?? 'ETB 10,000'),    access: t.ui.tier_buyer_gold_access ?? 'Unlock Gold seller images + dashboard',  color: 'from-yellow-50 to-amber-50 dark:from-yellow-900/20 dark:to-amber-900/10 border-yellow-200 dark:border-yellow-800/50', badge: 'bg-yellow-100 dark:bg-yellow-900/40 text-yellow-700 dark:text-yellow-400', href: '/buyer/dashboard', btnLabel: t.ui.upgrade_to_gold ?? 'Upgrade to Gold',      btnStyle: 'bg-yellow-500 hover:bg-yellow-600 text-white' },
                    { tier: t.ui.tier_premium ?? 'Premium',icon: '👑', cost: settings ? `ETB ${settings.buyer_premium_upgrade_price}` : (t.ui.tier_buyer_premium_cost ?? 'ETB 50,000'), access: t.ui.tier_buyer_premium_access ?? 'Unlock ALL images + advanced filters', color: 'from-purple-50 to-violet-50 dark:from-purple-900/20 dark:to-violet-900/10 border-purple-200 dark:border-purple-800/50', badge: 'bg-purple-100 dark:bg-purple-900/40 text-purple-700 dark:text-purple-400', href: '/buyer/dashboard', btnLabel: t.ui.upgrade_to_premium ?? 'Upgrade to Premium', btnStyle: 'bg-purple-600 hover:bg-purple-700 text-white' },
                  ].map(({ tier, icon, cost, access, color, badge, href, btnLabel, btnStyle }) => (
                    <div key={tier} className={`flex items-start gap-3 p-4 bg-gradient-to-r ${color} rounded-2xl border`}>
                      <span className="text-xl shrink-0 mt-0.5">{icon}</span>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-1">
                          <span className={`text-[10px] font-black px-2 py-0.5 rounded-full uppercase ${badge}`}>{tier}</span>
                          <span className="text-xs font-black text-zinc-700 dark:text-zinc-300">{cost}</span>
                        </div>
                        <p className="text-xs text-zinc-500 dark:text-zinc-400">{access}</p>
                        <Link href={href} className={`inline-flex items-center gap-1 mt-2.5 px-3 py-1.5 rounded-xl text-xs font-black transition-all ${btnStyle}`}>{btnLabel} <ArrowRight size={11} /></Link>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
            <div className="mt-6 flex items-center justify-center gap-2 text-sm text-zinc-400">
              <TrendingUp size={14} className="text-orange-500" />
              <p>{t.ui.tier_independence_note ?? '💡 Tiers are independent — be a Normal seller + Premium buyer at the same time'}</p>
            </div>
          </div>

          {/* ── Referral / Invite Promo ── */}
          <div className="w-full pt-12 border-t border-zinc-200/60 dark:border-zinc-800">
            <div className="relative rounded-3xl overflow-hidden border border-orange-200/60 dark:border-orange-800/40 shadow-2xl shadow-orange-600/10">
              {/* Background gradient */}
              <div className="absolute inset-0 bg-gradient-to-br from-orange-50 via-amber-50 to-yellow-50 dark:from-zinc-900 dark:via-zinc-900 dark:to-zinc-900" />
              <div className="absolute top-0 right-0 w-80 h-80 rounded-full bg-orange-400/10 blur-3xl pointer-events-none" />
              <div className="absolute bottom-0 left-0 w-64 h-64 rounded-full bg-yellow-400/10 blur-3xl pointer-events-none" />

              <div className="relative z-10 px-6 py-10 lg:px-12 lg:py-14 flex flex-col lg:flex-row items-center gap-10 lg:gap-16">

                {/* Left: copy */}
                <div className="flex-1 text-center lg:text-left space-y-4">
                  <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-black uppercase tracking-widest bg-orange-100 dark:bg-orange-900/40 text-orange-700 dark:text-orange-400 border border-orange-200 dark:border-orange-800/50">
                    <Gift size={12} /> {t.ui.referral_promo_badge ?? 'Invite & Earn'}
                  </div>
                  <h2 className="text-3xl lg:text-4xl font-black text-zinc-900 dark:text-white leading-tight">
                    {t.ui.referral_promo_title ?? 'Invite friends.'}<br />
                    <span className="text-orange-600">{settings ? `Earn ${settings.referral_commission_rate}% forever.` : (t.ui.referral_promo_title2 ?? 'Earn 2% forever.')}</span>
                  </h2>
                  <p className="text-sm lg:text-base text-zinc-500 dark:text-zinc-400 max-w-md mx-auto lg:mx-0 leading-relaxed">
                    {settings ? `Share your unique invite link. Every time someone you invited makes a paid action — product sale, tier upgrade — you earn ${settings.referral_commission_rate}% commission automatically.` : (t.ui.referral_promo_desc ?? 'Share your unique invite link. Every time someone you invited makes a paid action — product sale, tier upgrade — you earn 2% commission automatically.')}
                  </p>

                  {/* How it works steps */}
                  <div className="flex flex-col sm:flex-row gap-3 pt-2">
                    {[
                      { icon: <Users size={14} />, label: t.ui.referral_step1 ?? 'Register as a seller', color: 'text-orange-600 bg-orange-100 dark:bg-orange-900/30' },
                      { icon: <Gift size={14} />,  label: t.ui.referral_step2 ?? 'Get your invite link', color: 'text-amber-600 bg-amber-100 dark:bg-amber-900/30' },
                      { icon: <Wallet size={14} />,label: settings ? `Earn ${settings.referral_commission_rate}% per action` : (t.ui.referral_step3 ?? 'Earn 2% per action'),  color: 'text-green-600 bg-green-100 dark:bg-green-900/30' },
                    ].map(({ icon, label, color }, i) => (
                      <div key={i} className="flex items-center gap-2.5 flex-1">
                        <div className={`w-8 h-8 rounded-xl flex items-center justify-center shrink-0 ${color}`}>{icon}</div>
                        <span className="text-xs font-bold text-zinc-600 dark:text-zinc-300 leading-tight">{label}</span>
                        {i < 2 && <ArrowRight size={12} className="text-zinc-300 dark:text-zinc-600 shrink-0 hidden sm:block" />}
                      </div>
                    ))}
                  </div>

                  <div className="flex flex-col sm:flex-row items-center lg:items-start gap-3 pt-2">
                    <Link
                      href="/auth?as=seller"
                      className="inline-flex items-center gap-2 bg-orange-600 hover:bg-orange-500 text-white font-black text-sm px-6 py-3 rounded-xl transition-all shadow-lg shadow-orange-600/30 hover:-translate-y-0.5"
                    >
                      <Store size={15} /> {t.ui.referral_cta_register ?? 'Register as Seller'}
                    </Link>
                    <Link
                      href="/seller/dashboard"
                      className="inline-flex items-center gap-2 bg-white dark:bg-zinc-800 hover:bg-zinc-50 dark:hover:bg-zinc-700 text-zinc-700 dark:text-zinc-200 font-bold text-sm px-6 py-3 rounded-xl transition-all border border-zinc-200 dark:border-zinc-700 hover:-translate-y-0.5"
                    >
                      {t.ui.referral_cta_dashboard ?? 'View My Referrals'} <ArrowRight size={13} />
                    </Link>
                  </div>
                  <p className="text-[10px] text-zinc-400 dark:text-zinc-500 lg:text-left text-center">
                    Invite link opens seller registration — both sellers and buyers can register there.
                  </p>
                </div>

                {/* Right: stat cards */}
                <div className="shrink-0 grid grid-cols-2 gap-3 w-full max-w-xs">
                  {[
                    { icon: <Wallet size={20} />,  value: settings ? `${settings.referral_commission_rate}%` : '1.2%', label: t.ui.referral_stat1 ?? 'Commission per action', color: 'text-orange-600', bg: 'bg-orange-50 dark:bg-orange-900/20', border: 'border-orange-200 dark:border-orange-800/40' },
                    { icon: <Users size={20} />,   value: '∞',        label: t.ui.referral_stat2 ?? 'People you can invite',  color: 'text-blue-600',   bg: 'bg-blue-50 dark:bg-blue-900/20',   border: 'border-blue-200 dark:border-blue-800/40' },
                    { icon: <Gift size={20} />,    value: settings ? `ETB ${settings.min_referral_withdrawal}+` : 'ETB 500+', label: t.ui.referral_stat3 ?? 'Min. withdrawal',        color: 'text-green-600',  bg: 'bg-green-50 dark:bg-green-900/20', border: 'border-green-200 dark:border-green-800/40' },
                    { icon: <CheckCircle2 size={20} />, value: '1×', label: t.ui.referral_stat4 ?? 'Invite link per account', color: 'text-purple-600', bg: 'bg-purple-50 dark:bg-purple-900/20', border: 'border-purple-200 dark:border-purple-800/40' },
                  ].map(({ icon, value, label, color, bg, border }) => (
                    <div key={label} className={`flex flex-col items-center gap-2 p-4 rounded-2xl border ${bg} ${border} text-center`}>
                      <span className={color}>{icon}</span>
                      <span className={`text-2xl font-black leading-none ${color}`}>{value}</span>
                      <span className="text-[10px] font-bold text-zinc-400 leading-tight">{label}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Bottom note */}
              <div className="relative z-10 border-t border-orange-200/60 dark:border-orange-800/30 px-6 py-3 lg:px-12 bg-orange-50/50 dark:bg-orange-900/10">
                <p className="text-[11px] text-zinc-400 dark:text-zinc-500 text-center">
                  {t.ui.referral_promo_note ?? '⚡ Invite link is only usable at first registration — one inviter per account, forever.'}
                </p>
              </div>
            </div>
          </div>

        </div>
      </section>
    </>
  );
}
