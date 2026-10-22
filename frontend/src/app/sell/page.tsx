'use client';
import Link from 'next/link';
import { ShoppingBag, Home, Repeat, ArrowRight } from 'lucide-react';
import { useLocale } from '@/hooks/useLocale';

export default function SellPage() {
  const { t } = useLocale();
  const u = t.ui;

  const HOW_IT_WORKS = [
    { n: '01', text: u.sell_step_1 ?? 'Choose what you want to list below' },
    { n: '02', text: u.sell_step_2 ?? 'Sign in or create a free account' },
    { n: '03', text: u.sell_step_3 ?? 'Add your listing details & photos' },
    { n: '04', text: u.sell_step_4 ?? 'Wait for admin approval — then go live' },
  ];

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4 py-10">
      <div className="w-full max-w-2xl">

        {/* Header */}
        <div className="text-center mb-10">
          <h1 className="text-3xl font-black mb-2">{u.what_to_list ?? 'What do you want to list?'}</h1>
          <p className="text-zinc-500 text-base">{u.choose_listing_option ?? 'Choose one option below to get started'}</p>
        </div>

        {/* Three main options */}
        <div className="grid md:grid-cols-3 gap-4 mb-8">
          {/* SELL PRODUCTS */}
          <Link
            href="/auth?as=seller&preference=sale"
            className="group relative overflow-hidden rounded-2xl border-2 border-zinc-200 dark:border-zinc-700 hover:border-orange-600 p-6 flex flex-col items-center gap-4 transition-all hover:shadow-xl active:scale-95 bg-white dark:bg-zinc-900"
          >
            <div className="text-5xl">🛒</div>
            <div className="text-center">
              <p className="font-black text-lg mb-1">{u.sell_products ?? 'Sell Products'}</p>
              <p className="text-xs text-zinc-500 leading-relaxed">{u.sell_products_desc ?? 'List items for one-time purchase'}</p>
            </div>
            <div className="mt-auto pt-2 flex items-center gap-1 text-orange-600 text-xs font-bold group-hover:gap-2 transition-all">
              {u.register ?? 'Register'} <ArrowRight size={12} />
            </div>
          </Link>

          {/* RENT PROPERTY (LANDLORD) */}
          <Link
            href="/auth?as=seller&preference=rent"
            className="group relative overflow-hidden rounded-2xl border-2 border-zinc-200 dark:border-zinc-700 hover:border-blue-600 p-6 flex flex-col items-center gap-4 transition-all hover:shadow-xl active:scale-95 bg-white dark:bg-zinc-900"
          >
            <div className="text-5xl">🏠</div>
            <div className="text-center">
              <p className="font-black text-lg mb-1">{u.rent_property ?? 'Rent Property'}</p>
              <p className="text-xs text-zinc-500 leading-relaxed">{u.rent_property_desc ?? 'List properties for rental'}</p>
            </div>
            <div className="mt-auto pt-2 flex items-center gap-1 text-blue-600 text-xs font-bold group-hover:gap-2 transition-all">
              {u.register ?? 'Register'} <ArrowRight size={12} />
            </div>
          </Link>

          {/* BOTH */}
          <Link
            href="/auth?as=seller&preference=both"
            className="group relative overflow-hidden rounded-2xl border-2 border-zinc-200 dark:border-zinc-700 hover:border-purple-600 p-6 flex flex-col items-center gap-4 transition-all hover:shadow-xl active:scale-95 bg-white dark:bg-zinc-900"
          >
            <div className="text-5xl">🔄</div>
            <div className="text-center">
              <p className="font-black text-lg mb-1">{u.both_option ?? 'Both'}</p>
              <p className="text-xs text-zinc-500 leading-relaxed">{u.both_option_desc ?? 'Sell products AND rent properties'}</p>
            </div>
            <div className="mt-auto pt-2 flex items-center gap-1 text-purple-600 text-xs font-bold group-hover:gap-2 transition-all">
              {u.register ?? 'Register'} <ArrowRight size={12} />
            </div>
          </Link>
        </div>

        {/* How it works */}
        <div className="bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-6">
          <p className="text-xs font-black uppercase tracking-widest text-zinc-400 mb-4">
            {u.how_it_works ?? 'How it works'}
          </p>
          <div className="grid md:grid-cols-4 gap-4">
            {HOW_IT_WORKS.map(s => (
              <div key={s.n} className="flex flex-col items-center text-center gap-2">
                <span className="w-10 h-10 rounded-xl bg-orange-600/10 text-orange-600 text-sm font-black flex items-center justify-center">{s.n}</span>
                <p className="text-xs text-zinc-600 dark:text-zinc-400 leading-relaxed">{s.text}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Already have account */}
        <p className="text-center text-sm text-zinc-400 mt-6">
          {u.already_have_account ?? 'Already have an account?'}{' '}
          <Link href="/auth?as=seller" className="text-orange-600 font-bold hover:underline">
            {u.sign_in ?? 'Sign in'} <ArrowRight size={12} className="inline" />
          </Link>
        </p>
      </div>
    </div>
  );
}
