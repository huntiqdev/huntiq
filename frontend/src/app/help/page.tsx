'use client';

import { useState, useMemo } from 'react';
import Link from 'next/link';
import { useLocale } from '@/hooks/useLocale';
import {
  HelpCircle,
  Search,
  X,
  ChevronDown,
  ChevronUp,
  Shield,
  ShieldCheck,
  ShoppingBag,
  Tag,
  Phone,
  Send,
  MessageCircle,
  CreditCard,
  Home,
  CheckCircle2,
  AlertTriangle,
  Sparkles,
} from 'lucide-react';

interface FaqItem {
  id: string;
  category: 'buying' | 'selling' | 'renting' | 'payments' | 'safety';
  q: string;
  a: string;
}

export default function HelpPage() {
  const { t, locale } = useLocale();
  const h = (t as any).help_page ?? {};

  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCat, setSelectedCat] = useState<string>('all');
  const [openFaq, setOpenFaq] = useState<string | null>('faq1');

  // Categories definitions with localized labels
  const categories = useMemo(() => [
    { id: 'all',      label: h.all_categories || 'All Questions',           icon: Sparkles },
    { id: 'buying',   label: h.cat_buying || 'Buying & Orders',             icon: ShoppingBag },
    { id: 'selling',  label: h.cat_selling || 'Selling & Shops',            icon: Tag },
    { id: 'renting',  label: h.cat_renting || 'Renting & Real Estate',      icon: Home },
    { id: 'payments', label: h.cat_payments || 'Payments (Telebirr & CBE)', icon: CreditCard },
    { id: 'safety',   label: h.cat_safety || 'Safety & Verification',       icon: ShieldCheck },
  ], [h]);

  // Master FAQ catalogue covering all Ethiopian marketplace scenarios
  const allFaqs: FaqItem[] = useMemo(() => [
    {
      id: 'faq1',
      category: 'buying',
      q: h.faq1q || 'How do I buy a product?',
      a: h.faq1a || 'Browse products on the Buy page, click on any item to view details, then click the Buy button to confirm your order. The seller will contact you to arrange delivery.',
    },
    {
      id: 'faq2',
      category: 'selling',
      q: h.faq2q || 'How do I list a product for sale?',
      a: h.faq2a || 'Click the Sell button in the header or on the home page. Fill in the product title, description, price, category, condition, and upload photos. Submit to publish your listing instantly.',
    },
    {
      id: 'faq3',
      category: 'payments',
      q: h.faq3q || 'Is payment done online or in person?',
      a: h.faq3a || 'Payment confirmation is done through the platform, but final delivery and cash exchange is arranged directly between buyer and seller. Always meet in a safe public place.',
    },
    {
      id: 'faq_pay',
      category: 'payments',
      q: h.faq_pay_q || 'What payment methods are supported for orders?',
      a: h.faq_pay_a || 'We support Ethiopian mobile payments including Telebirr, CBE Birr, Bank of Abyssinia, Awash Bank, and Cash on Delivery upon physical inspection.',
    },
    {
      id: 'faq_rent',
      category: 'renting',
      q: h.faq_rent_q || 'How does property and vehicle renting work on HUNTIQ?',
      a: h.faq_rent_a || 'Landlords and verified rental agencies list apartments, villas, and vehicles with transparent daily or monthly rates. Renters can call the host directly, arrange an in-person viewing, and sign agreements without middleman markups.',
    },
    {
      id: 'faq4',
      category: 'buying',
      q: h.faq4q || 'How do I contact the seller?',
      a: h.faq4a || 'Each product listing shows the seller\'s phone number. You can call or message them directly. Our customer support team is also available on the Contact page.',
    },
    {
      id: 'faq5',
      category: 'selling',
      q: h.faq5q || 'Can I sell used items?',
      a: h.faq5a || 'Yes! HUNTIQ supports both new and used items. When listing, select Used under the Condition field so buyers know the item\'s exact state.',
    },
    {
      id: 'faq_verify',
      category: 'safety',
      q: h.faq_verify_q || 'What does the Verified Seller badge mean?',
      a: h.faq_verify_a || 'Verified sellers have provided authenticated Ethiopian National ID (Fayda / Kebele), commercial TIN registration, or physical store location validation by our verification team.',
    },
    {
      id: 'faq6',
      category: 'safety',
      q: h.faq6q || 'How do I report a suspicious listing?',
      a: h.faq6a || 'Click the flag icon on any product card or inside the product detail view. Fill in the reason and submit. Our moderation team reviews all reports within 24 hours.',
    },
    {
      id: 'faq_dispute',
      category: 'safety',
      q: h.faq_dispute_q || 'What should I do if an item does not match its description?',
      a: h.faq_dispute_a || 'Inspect all items before handing over money or releasing transfer. If an issue occurs, report the listing immediately via the Flag icon or contact customer support with your order reference.',
    },
    {
      id: 'faq7',
      category: 'buying',
      q: h.faq7q || 'What categories are supported?',
      a: h.faq7a || 'We support Electronics, Vehicles, Property, Home Appliances, Furniture, Clothing, Agriculture, Industry, Construction, Services, Jobs, Animals, Food, Health & Beauty, Sports, Education, Entertainment, Business, Luxury, and Miscellaneous.',
    },
    {
      id: 'faq8',
      category: 'buying',
      q: h.faq8q || 'How do I change the language?',
      a: h.faq8a || 'Click the globe icon in the top navigation bar. You can switch between Amharic (አማርኛ), English, Afaan Oromoo, and Tigrinya (ትግርኛ) at any time.',
    },
    {
      id: 'faq9',
      category: 'buying',
      q: h.faq9q || 'How do I save a search?',
      a: h.faq9a || 'After applying filters on the Browse page, click Save Search. Give it a name and it will appear in your Profile under Saved Searches for quick access later.',
    },
    {
      id: 'faq10',
      category: 'safety',
      q: h.faq10q || 'Is my personal data safe?',
      a: h.faq10a || 'We only store the minimum data needed to run the platform. Your phone number is only visible on listings you create. Read our Privacy Policy for full details.',
    },
  ], [h]);

  // Filter FAQs based on active category and search term
  const filteredFaqs = useMemo(() => {
    const q = searchQuery.trim().toLowerCase();
    return allFaqs.filter(faq => {
      const matchCat = selectedCat === 'all' || faq.category === selectedCat;
      if (!matchCat) return false;
      if (!q) return true;
      return faq.q.toLowerCase().includes(q) || faq.a.toLowerCase().includes(q);
    });
  }, [allFaqs, selectedCat, searchQuery]);

  const steps = [
    { icon: ShoppingBag, title: h.step1_title, desc: h.step1_desc },
    { icon: Phone,       title: h.step2_title, desc: h.step2_desc },
    { icon: CheckCircle2, title: h.step3_title, desc: h.step3_desc },
  ].filter(s => s.title);

  const tips = [h.tip1, h.tip2, h.tip3, h.tip4, h.tip5].filter(Boolean);

  return (
    <div className="min-h-screen bg-zinc-50/50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-100 transition-colors">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 space-y-10 sm:space-y-12">

        {/* ── 1. Hero & Live Search ── */}
        <section className="text-center space-y-4">
          <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-orange-500/10 border border-orange-500/20 text-orange-600 dark:text-orange-400 text-xs font-bold tracking-wide">
            <HelpCircle size={14} />
            <span>HUNTIQ Ethiopia Support Hub</span>
          </div>

          <h1 className="text-3xl sm:text-4xl lg:text-5xl font-black tracking-tight bg-gradient-to-r from-zinc-900 via-orange-600 to-amber-600 dark:from-white dark:via-orange-400 dark:to-amber-400 bg-clip-text text-transparent">
            {h.title || 'Help Center'}
          </h1>
          <p className="text-sm sm:text-base text-zinc-600 dark:text-zinc-400 max-w-xl mx-auto leading-relaxed">
            {h.subtitle || 'Find quick answers to common questions or reach out to our dedicated Ethiopian support team.'}
          </p>

          {/* Interactive Search Bar */}
          <div className="max-w-xl mx-auto relative pt-2">
            <div className="relative flex items-center">
              <Search size={18} className="absolute left-4 text-zinc-400 pointer-events-none" />
              <input
                type="text"
                value={searchQuery}
                onChange={e => setSearchQuery(e.target.value)}
                placeholder={h.search_placeholder || 'Search questions, topics, or keywords...'}
                className="w-full pl-11 pr-10 py-3.5 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-sm hover:border-orange-500/40 focus:border-orange-500 focus:ring-4 focus:ring-orange-500/10 text-sm outline-none transition-all"
              />
              {searchQuery && (
                <button
                  type="button"
                  onClick={() => setSearchQuery('')}
                  className="absolute right-3.5 p-1 rounded-full text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
                >
                  <X size={15} />
                </button>
              )}
            </div>
            {searchQuery && (
              <p className="text-xs text-left text-zinc-500 mt-2 px-1">
                Showing {filteredFaqs.length} results for &quot;<span className="text-orange-600 dark:text-orange-400 font-bold">{searchQuery}</span>&quot;
              </p>
            )}
          </div>
        </section>

        {/* ── 2. Ethiopia Safe Trade Protocol (Golden Rules) ── */}
        <section className="bg-gradient-to-br from-amber-500/10 via-orange-500/5 to-transparent border border-amber-500/30 rounded-3xl p-5 sm:p-7 backdrop-blur-sm relative overflow-hidden">
          <div className="flex flex-col sm:flex-row items-start gap-4">
            <div className="w-12 h-12 rounded-2xl bg-gradient-to-tr from-amber-500 to-orange-500 text-white flex items-center justify-center shrink-0 shadow-lg shadow-orange-500/20">
              <ShieldCheck size={26} />
            </div>
            <div className="space-y-2 flex-1">
              <h2 className="text-lg sm:text-xl font-black text-zinc-900 dark:text-white flex items-center gap-2">
                <span>{h.safe_exchange_title || 'Ethiopia Safe Trade Protocol'}</span>
                <span className="text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-full bg-amber-500/20 text-amber-700 dark:text-amber-300">
                  Essential
                </span>
              </h2>
              <p className="text-xs sm:text-sm text-zinc-600 dark:text-zinc-300 leading-relaxed">
                {h.safe_exchange_p1 || 'Always meet in busy, well-lit public commercial hubs (Bole, Piassa, Mexico, Megenagna, Hawassa, Adama, Mekelle, Bahir Dar) during daylight hours.'}
              </p>
              <p className="text-xs sm:text-sm text-zinc-600 dark:text-zinc-300 leading-relaxed font-semibold text-amber-900 dark:text-amber-200">
                {h.safe_exchange_p2 || 'Inspect devices, vehicles, or items thoroughly before completing Telebirr or CBE Birr transfers. Never share SMS OTP codes or mobile banking PINs.'}
              </p>
            </div>
          </div>
        </section>

        {/* ── 3. How It Works Flow ── */}
        {steps.length > 0 && (
          <section className="space-y-4">
            <div className="flex items-center justify-between">
              <h2 className="text-xs font-black uppercase tracking-widest text-zinc-400">
                {h.how_it_works || 'How It Works'}
              </h2>
              <span className="text-xs font-bold text-orange-600">3 Simple Steps</span>
            </div>
            <div className="grid sm:grid-cols-3 gap-4">
              {steps.map(({ icon: Icon, title, desc }, i) => (
                <div
                  key={i}
                  className="bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 p-5 shadow-sm hover:shadow-md transition-shadow relative overflow-hidden group"
                >
                  <div className="flex items-center justify-between mb-3">
                    <div className="w-10 h-10 rounded-xl bg-orange-500/10 text-orange-600 flex items-center justify-center group-hover:scale-110 transition-transform">
                      <Icon size={18} />
                    </div>
                    <span className="text-sm font-black text-orange-600/60 dark:text-orange-400/60">
                      0{i + 1}
                    </span>
                  </div>
                  <h3 className="font-bold text-base mb-1.5 text-zinc-900 dark:text-white">{title}</h3>
                  <p className="text-xs text-zinc-500 dark:text-zinc-400 leading-relaxed">{desc}</p>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* ── 4. Categorized FAQ Section ── */}
        <section className="space-y-6">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
            <div>
              <h2 className="text-2xl font-black text-zinc-900 dark:text-white">
                {h.faq || 'Frequently Asked Questions'}
              </h2>
              <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-0.5">
                Categorized answers for Ethiopian buyers, sellers, landlords, and shoppers
              </p>
            </div>
          </div>

          {/* Category Filter Tabs */}
          <div className="flex items-center gap-2 overflow-x-auto pb-2 scrollbar-none">
            {categories.map(cat => {
              const Icon = cat.icon;
              const isActive = selectedCat === cat.id;
              const count = cat.id === 'all'
                ? allFaqs.length
                : allFaqs.filter(f => f.category === cat.id).length;

              return (
                <button
                  key={cat.id}
                  type="button"
                  onClick={() => setSelectedCat(cat.id)}
                  className={`flex items-center gap-2 px-3.5 py-2 rounded-xl text-xs font-bold whitespace-nowrap transition-all shrink-0 border ${
                    isActive
                      ? 'bg-orange-600 text-white border-orange-600 shadow-sm shadow-orange-600/30'
                      : 'bg-white dark:bg-zinc-900 text-zinc-600 dark:text-zinc-400 border-zinc-200 dark:border-zinc-800 hover:border-orange-500/40 hover:text-orange-600'
                  }`}
                >
                  <Icon size={13} className={isActive ? 'text-white' : 'text-orange-600 dark:text-orange-400'} />
                  <span>{cat.label}</span>
                  <span className={`text-[10px] px-1.5 py-0.2 rounded-full font-black ${
                    isActive ? 'bg-white/20 text-white' : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-400'
                  }`}>
                    {count}
                  </span>
                </button>
              );
            })}
          </div>

          {/* FAQ Accordion List */}
          <div className="space-y-3">
            {filteredFaqs.length === 0 ? (
              <div className="bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 p-8 text-center space-y-3">
                <AlertTriangle size={24} className="mx-auto text-amber-500 opacity-60" />
                <p className="text-sm font-bold text-zinc-700 dark:text-zinc-300">
                  {h.no_faq_found || 'No questions match your search. Try different terms or contact our support.'}
                </p>
                <button
                  type="button"
                  onClick={() => { setSearchQuery(''); setSelectedCat('all'); }}
                  className="text-xs font-bold text-orange-600 hover:underline"
                >
                  View all questions →
                </button>
              </div>
            ) : (
              filteredFaqs.map(faq => {
                const isOpen = openFaq === faq.id;
                return (
                  <div
                    key={faq.id}
                    className="bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 overflow-hidden shadow-xs hover:border-orange-500/40 transition-colors"
                  >
                    <button
                      type="button"
                      onClick={() => setOpenFaq(isOpen ? null : faq.id)}
                      className="w-full flex items-center justify-between px-5 py-4 text-left transition-colors hover:bg-zinc-50 dark:hover:bg-zinc-800/40"
                    >
                      <div className="flex items-center gap-3 pr-4">
                        <span className="w-2 h-2 rounded-full bg-orange-600 shrink-0" />
                        <span className="font-bold text-sm text-zinc-900 dark:text-white leading-snug">
                          {faq.q}
                        </span>
                      </div>
                      <div className="shrink-0 text-zinc-400">
                        {isOpen ? <ChevronUp size={18} /> : <ChevronDown size={18} />}
                      </div>
                    </button>
                    {isOpen && (
                      <div className="px-5 pb-5 text-sm text-zinc-600 dark:text-zinc-300 leading-relaxed border-t border-zinc-100 dark:border-zinc-800 pt-3.5 bg-zinc-50/50 dark:bg-zinc-900/50">
                        {faq.a}
                      </div>
                    )}
                  </div>
                );
              })
            )}
          </div>
        </section>

        {/* ── 5. Safety Tips Quick Checklist ── */}
        {tips.length > 0 && (
          <section className="space-y-4">
            <h2 className="text-xs font-black uppercase tracking-widest text-zinc-400 flex items-center gap-1.5">
              <Shield size={13} className="text-orange-600" />
              <span>{h.safety_tips ?? h.safety}</span>
            </h2>
            <div className="grid sm:grid-cols-2 gap-3">
              {tips.map((tip, i) => (
                <div
                  key={i}
                  className="flex items-start gap-3 bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-3.5 shadow-xs"
                >
                  <span className="w-5 h-5 rounded-full bg-orange-500/10 text-orange-600 text-xs font-black flex items-center justify-center shrink-0 mt-0.5">
                    {i + 1}
                  </span>
                  <p className="text-xs sm:text-sm text-zinc-700 dark:text-zinc-300 leading-relaxed font-medium">
                    {tip}
                  </p>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* ── 6. Still Have Questions? Direct Contact ── */}
        <section className="bg-gradient-to-r from-orange-600 to-amber-600 rounded-3xl p-6 sm:p-8 text-white shadow-xl shadow-orange-600/20">
          <div className="flex flex-col md:flex-row items-center justify-between gap-6">
            <div className="space-y-2 text-center md:text-left">
              <h2 className="text-2xl sm:text-3xl font-black">
                {h.still_have_questions || 'Still have questions?'}
              </h2>
              <p className="text-sm opacity-90 max-w-lg leading-relaxed">
                {h.contact_support_desc || 'Our Ethiopian customer care team is available 7 days a week to assist you.'}
              </p>
            </div>

            <div className="flex flex-wrap items-center justify-center gap-3 shrink-0">
              <a
                href="tel:+251911000000"
                className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-white text-orange-600 font-bold text-xs sm:text-sm shadow-md hover:bg-zinc-100 transition-colors"
              >
                <Phone size={15} />
                <span>{h.call_support || 'Call Helpline'}</span>
              </a>

              <a
                href="https://t.me/huntiq_support"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-white/20 hover:bg-white/30 text-white font-bold text-xs sm:text-sm backdrop-blur-sm transition-colors border border-white/20"
              >
                <Send size={15} />
                <span>{h.telegram_support || 'Telegram'}</span>
              </a>

              <Link
                href="/contact"
                className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-zinc-900 text-white hover:bg-black font-bold text-xs sm:text-sm transition-colors"
              >
                <MessageCircle size={15} />
                <span>{h.contact_form || 'Contact Form'}</span>
              </Link>
            </div>
          </div>
        </section>

      </div>
    </div>
  );
}
