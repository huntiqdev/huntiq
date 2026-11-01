'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { useLocaleContext } from '@/context/LocaleContext';
import { api } from '@/lib/api';
import { ShoppingBag, Search, ArrowRight, Sparkles } from 'lucide-react';

export default function DemandPage() {
  const router = useRouter();
  const { user } = useAuth();
  const { t } = useLocaleContext();
  const demand = (t as any).demand as Record<string, string>;
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleDemand = async (demandType: 'FIND' | 'LIST') => {
    if (!user) {
      router.push('/auth');
      return;
    }

    setLoading(true);
    setError('');

    try {
      await api.post('/auth/demand', { demand: demandType }, {
        headers: { Authorization: `Bearer ${user.token}` }
      });

      if (demandType === 'FIND') {
        router.push('/find');
      } else {
        router.push('/list');
      }
    } catch (err: any) {
      setError(err.message || 'Something went wrong');
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-blue-50 dark:from-slate-900 dark:via-slate-800 dark:to-slate-900">
      <div className="container mx-auto px-4 py-12 max-w-5xl">
        {/* Header */}
        <div className="text-center mb-12 animate-fade-in">
          <div className="inline-flex items-center gap-2 px-4 py-2 bg-blue-100 dark:bg-blue-900/30 rounded-full mb-4">
            <Sparkles className="w-4 h-4 text-blue-600 dark:text-blue-400" />
            <span className="text-sm font-medium text-blue-700 dark:text-blue-300">
              {demand.welcome}
            </span>
          </div>
          <h1 className="text-4xl md:text-5xl font-bold text-slate-900 dark:text-white mb-3">
            {demand.title}
          </h1>
          <p className="text-lg text-slate-600 dark:text-slate-300">
            {demand.subtitle}
          </p>
        </div>

        {error && (
          <div className="mb-8 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-700 dark:text-red-300 text-center">
            {error}
          </div>
        )}

        {/* Two Pillars */}
        <div className="grid md:grid-cols-2 gap-4 sm:gap-6 mb-8">
          {/* FIND Pillar */}
          <button
            onClick={() => handleDemand('FIND')}
            disabled={loading}
            className="group relative overflow-hidden bg-white dark:bg-zinc-900 rounded-2xl p-5 sm:p-8 shadow-lg hover:shadow-2xl transition-all duration-300 border-2 border-zinc-200 dark:border-zinc-800 hover:border-orange-500 dark:hover:border-orange-500 disabled:opacity-50 disabled:cursor-not-allowed tap-bounce active:scale-[0.98] text-left"
          >
            <div className="absolute top-0 right-0 w-32 h-32 bg-orange-500/10 rounded-full -mr-16 -mt-16 group-hover:scale-150 transition-transform duration-500"></div>
            
            <div className="relative z-10">
              <div className="w-14 h-14 sm:w-16 sm:h-16 bg-orange-100 dark:bg-orange-950/40 rounded-2xl flex items-center justify-center mb-4 sm:mb-6 group-hover:scale-110 transition-transform">
                <Search className="w-7 h-7 sm:w-8 sm:h-8 text-orange-600 dark:text-orange-400" />
              </div>
              
              <h2 className="text-2xl sm:text-3xl font-bold text-zinc-900 dark:text-white mb-1.5 sm:mb-2">
                {demand.find}
              </h2>
              <p className="text-base sm:text-lg font-medium text-orange-600 dark:text-orange-400 mb-3 sm:mb-4">
                {demand.find_title}
              </p>
              <p className="text-sm text-zinc-600 dark:text-zinc-300 mb-5 leading-relaxed">
                {demand.find_description}
              </p>
              
              <div className="flex flex-col gap-2 text-xs sm:text-sm text-zinc-700 dark:text-zinc-300 mb-6">
                <div className="flex items-center gap-2">
                  <div className="w-1.5 h-1.5 bg-orange-500 rounded-full"></div>
                  <span>{demand.find_buy}</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-1.5 h-1.5 bg-teal-500 rounded-full"></div>
                  <span>{demand.find_rent}</span>
                </div>
              </div>

              <div className="flex items-center gap-2 text-orange-600 dark:text-orange-400 font-semibold group-hover:gap-3 transition-all text-sm">
                <span>{demand.continue}</span>
                <ArrowRight className="w-4 h-4 sm:w-5 sm:h-5" />
              </div>
            </div>
          </button>

          {/* LIST Pillar */}
          <button
            onClick={() => handleDemand('LIST')}
            disabled={loading}
            className="group relative overflow-hidden bg-white dark:bg-zinc-900 rounded-2xl p-5 sm:p-8 shadow-lg hover:shadow-2xl transition-all duration-300 border-2 border-zinc-200 dark:border-zinc-800 hover:border-orange-500 dark:hover:border-orange-500 disabled:opacity-50 disabled:cursor-not-allowed tap-bounce active:scale-[0.98] text-left"
          >
            <div className="absolute top-0 right-0 w-32 h-32 bg-orange-500/10 rounded-full -mr-16 -mt-16 group-hover:scale-150 transition-transform duration-500"></div>
            
            <div className="relative z-10">
              <div className="w-14 h-14 sm:w-16 sm:h-16 bg-orange-100 dark:bg-orange-950/40 rounded-2xl flex items-center justify-center mb-4 sm:mb-6 group-hover:scale-110 transition-transform">
                <ShoppingBag className="w-7 h-7 sm:w-8 sm:h-8 text-orange-600 dark:text-orange-400" />
              </div>
              
              <h2 className="text-2xl sm:text-3xl font-bold text-zinc-900 dark:text-white mb-1.5 sm:mb-2">
                {demand.list}
              </h2>
              <p className="text-base sm:text-lg font-medium text-orange-600 dark:text-orange-400 mb-3 sm:mb-4">
                {demand.list_title}
              </p>
              <p className="text-sm text-zinc-600 dark:text-zinc-300 mb-5 leading-relaxed">
                {demand.list_description}
              </p>
              
              <div className="flex flex-col gap-2 text-xs sm:text-sm text-zinc-700 dark:text-zinc-300 mb-6">
                <div className="flex items-center gap-2">
                  <div className="w-1.5 h-1.5 bg-orange-500 rounded-full"></div>
                  <span>{demand.list_sell}</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-1.5 h-1.5 bg-teal-500 rounded-full"></div>
                  <span>{demand.list_landlord}</span>
                </div>
              </div>

              <div className="flex items-center gap-2 text-orange-600 dark:text-orange-400 font-semibold group-hover:gap-3 transition-all text-sm">
                <span>{demand.continue}</span>
                <ArrowRight className="w-4 h-4 sm:w-5 sm:h-5" />
              </div>
            </div>
          </button>
        </div>

        {/* Bottom hint */}
        <p className="text-center text-sm text-slate-500 dark:text-slate-400">
          You can always switch between buying, renting, selling, and listing properties
        </p>
      </div>

      <style jsx>{`
        @keyframes fade-in {
          from { opacity: 0; transform: translateY(-10px); }
          to { opacity: 1; transform: translateY(0); }
        }
        .animate-fade-in {
          animation: fade-in 0.5s ease-out;
        }
      `}</style>
    </div>
  );
}
