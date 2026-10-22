'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { useLocaleContext } from '@/context/LocaleContext';
import { Tag, Building2, ArrowRight } from 'lucide-react';

export default function ListPage() {
  const router = useRouter();
  const { user, loading } = useAuth();
  const { t } = useLocaleContext();
  const demand = (t as any).demand as Record<string, string>;

  // Ensure user is authenticated for listing creation
  useEffect(() => {
    if (!loading && !user) {
      router.push('/auth');
    }
  }, [user, loading, router]);

  if (loading || !user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-orange-50 dark:from-slate-900 dark:via-slate-800 dark:to-slate-900">
      <div className="container mx-auto px-4 py-12 max-w-4xl">
        <div className="text-center mb-12">
          <h1 className="text-4xl md:text-5xl font-bold text-slate-900 dark:text-white mb-3">
            {demand.list_title}
          </h1>
          <p className="text-lg text-slate-600 dark:text-slate-300">
            {demand.list_subtitle}
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-6">
          {/* SELL */}
          <button
            onClick={() => router.push('/seller/product/new?type=sale')}
            className="group bg-white dark:bg-slate-800 rounded-2xl p-8 shadow-lg hover:shadow-2xl transition-all duration-300 border-2 border-slate-200 dark:border-slate-700 hover:border-orange-400 dark:hover:border-orange-500"
          >
            <div className="w-16 h-16 bg-orange-100 dark:bg-orange-900/30 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
              <Tag className="w-8 h-8 text-orange-600 dark:text-orange-400" />
            </div>
            
            <h2 className="text-2xl font-bold text-slate-900 dark:text-white mb-2">
              {demand.sell_title}
            </h2>
            <p className="text-slate-600 dark:text-slate-300 mb-6">
              {demand.sell_description}
            </p>
            
            <div className="flex items-center justify-center gap-2 text-orange-600 dark:text-orange-400 font-semibold group-hover:gap-4 transition-all">
              <span>{demand.sell_cta}</span>
              <ArrowRight className="w-5 h-5" />
            </div>
          </button>

          {/* RENT OUT / LANDLORD */}
          <button
            onClick={() => router.push('/seller/product/new?type=rent')}
            className="group bg-white dark:bg-slate-800 rounded-2xl p-8 shadow-lg hover:shadow-2xl transition-all duration-300 border-2 border-slate-200 dark:border-slate-700 hover:border-teal-400 dark:hover:border-teal-500"
          >
            <div className="w-16 h-16 bg-teal-100 dark:bg-teal-900/30 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
              <Building2 className="w-8 h-8 text-teal-600 dark:text-teal-400" />
            </div>
            
            <h2 className="text-2xl font-bold text-slate-900 dark:text-white mb-2">
              {demand.landlord_title}
            </h2>
            <p className="text-slate-600 dark:text-slate-300 mb-6">
              {demand.landlord_description}
            </p>
            
            <div className="flex items-center justify-center gap-2 text-teal-600 dark:text-teal-400 font-semibold group-hover:gap-4 transition-all">
              <span>{demand.landlord_cta}</span>
              <ArrowRight className="w-5 h-5" />
            </div>
          </button>
        </div>
      </div>
    </div>
  );
}
