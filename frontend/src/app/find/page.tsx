'use client';
import { useRouter } from 'next/navigation';
import { useLocaleContext } from '@/context/LocaleContext';
import { ShoppingCart, Home, ArrowRight } from 'lucide-react';

export default function FindPage() {
  const router = useRouter();
  const { t } = useLocaleContext();
  const demand = (t as any).demand as Record<string, string>;

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-blue-50 dark:from-slate-900 dark:via-slate-800 dark:to-slate-900">
      <div className="container mx-auto px-4 py-12 max-w-4xl">
        <div className="text-center mb-12">
          <h1 className="text-4xl md:text-5xl font-bold text-slate-900 dark:text-white mb-3">
            {demand.find_title}
          </h1>
          <p className="text-lg text-slate-600 dark:text-slate-300">
            {demand.find_subtitle}
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-4 sm:gap-6">
          {/* BUY */}
          <button
            onClick={() => router.push('/buy')}
            className="group bg-white dark:bg-zinc-900 rounded-2xl p-5 sm:p-8 shadow-lg hover:shadow-2xl transition-all duration-300 border-2 border-zinc-200 dark:border-zinc-800 hover:border-orange-500 dark:hover:border-orange-500 tap-bounce active:scale-[0.98] text-left"
          >
            <div className="w-14 h-14 sm:w-16 sm:h-16 bg-orange-100 dark:bg-orange-950/40 rounded-2xl flex items-center justify-center mb-4 sm:mb-6 group-hover:scale-110 transition-transform">
              <ShoppingCart className="w-7 h-7 sm:w-8 sm:h-8 text-orange-600 dark:text-orange-400" />
            </div>
            
            <h2 className="text-xl sm:text-2xl font-bold text-zinc-900 dark:text-white mb-2">
              {demand.buy_title}
            </h2>
            <p className="text-sm text-zinc-600 dark:text-zinc-300 mb-5 leading-relaxed">
              {demand.buy_description}
            </p>
            
            <div className="flex items-center gap-2 text-orange-600 dark:text-orange-400 font-semibold group-hover:gap-3 transition-all text-sm">
              <span>{demand.buy_cta}</span>
              <ArrowRight className="w-4 h-4 sm:w-5 sm:h-5" />
            </div>
          </button>

          {/* RENT */}
          <button
            onClick={() => router.push('/rent')}
            className="group bg-white dark:bg-zinc-900 rounded-2xl p-5 sm:p-8 shadow-lg hover:shadow-2xl transition-all duration-300 border-2 border-zinc-200 dark:border-zinc-800 hover:border-teal-500 dark:hover:border-teal-500 tap-bounce active:scale-[0.98] text-left"
          >
            <div className="w-14 h-14 sm:w-16 sm:h-16 bg-teal-100 dark:bg-teal-950/40 rounded-2xl flex items-center justify-center mb-4 sm:mb-6 group-hover:scale-110 transition-transform">
              <Home className="w-7 h-7 sm:w-8 sm:h-8 text-teal-600 dark:text-teal-400" />
            </div>
            
            <h2 className="text-xl sm:text-2xl font-bold text-zinc-900 dark:text-white mb-2">
              {demand.rent_title}
            </h2>
            <p className="text-sm text-zinc-600 dark:text-zinc-300 mb-5 leading-relaxed">
              {demand.rent_description}
            </p>
            
            <div className="flex items-center gap-2 text-teal-600 dark:text-teal-400 font-semibold group-hover:gap-3 transition-all text-sm">
              <span>{demand.rent_cta}</span>
              <ArrowRight className="w-4 h-4 sm:w-5 sm:h-5" />
            </div>
          </button>
        </div>
      </div>
    </div>
  );
}
