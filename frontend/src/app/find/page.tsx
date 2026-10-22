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

        <div className="grid md:grid-cols-2 gap-6">
          {/* BUY */}
          <button
            onClick={() => router.push('/buy')}
            className="group bg-white dark:bg-slate-800 rounded-2xl p-8 shadow-lg hover:shadow-2xl transition-all duration-300 border-2 border-slate-200 dark:border-slate-700 hover:border-blue-400 dark:hover:border-blue-500"
          >
            <div className="w-16 h-16 bg-blue-100 dark:bg-blue-900/30 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
              <ShoppingCart className="w-8 h-8 text-blue-600 dark:text-blue-400" />
            </div>
            
            <h2 className="text-2xl font-bold text-slate-900 dark:text-white mb-2">
              {demand.buy_title}
            </h2>
            <p className="text-slate-600 dark:text-slate-300 mb-6">
              {demand.buy_description}
            </p>
            
            <div className="flex items-center justify-center gap-2 text-blue-600 dark:text-blue-400 font-semibold group-hover:gap-4 transition-all">
              <span>{demand.buy_cta}</span>
              <ArrowRight className="w-5 h-5" />
            </div>
          </button>

          {/* RENT */}
          <button
            onClick={() => router.push('/rent')}
            className="group bg-white dark:bg-slate-800 rounded-2xl p-8 shadow-lg hover:shadow-2xl transition-all duration-300 border-2 border-slate-200 dark:border-slate-700 hover:border-purple-400 dark:hover:border-purple-500"
          >
            <div className="w-16 h-16 bg-purple-100 dark:bg-purple-900/30 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
              <Home className="w-8 h-8 text-purple-600 dark:text-purple-400" />
            </div>
            
            <h2 className="text-2xl font-bold text-slate-900 dark:text-white mb-2">
              {demand.rent_title}
            </h2>
            <p className="text-slate-600 dark:text-slate-300 mb-6">
              {demand.rent_description}
            </p>
            
            <div className="flex items-center justify-center gap-2 text-purple-600 dark:text-purple-400 font-semibold group-hover:gap-4 transition-all">
              <span>{demand.rent_cta}</span>
              <ArrowRight className="w-5 h-5" />
            </div>
          </button>
        </div>
      </div>
    </div>
  );
}
