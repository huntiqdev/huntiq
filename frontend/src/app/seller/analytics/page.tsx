'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import Link from 'next/link';
import {
  BarChart2, TrendingUp, Package, DollarSign, ShoppingCart,
  Crown, Lock, Zap, ArrowUpRight, Clock,
} from 'lucide-react';
import { TierUpgradeModal } from '@/components/seller/TierUpgradeModal';

function StatCard({ icon: Icon, label, value, sub, color }: {
  icon: any; label: string; value: string | number; sub?: string; color: string;
}) {
  return (
    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4">
      <div className={`w-9 h-9 rounded-xl flex items-center justify-center mb-3 ${color}`}>
        <Icon size={16} />
      </div>
      <p className="text-2xl font-black text-zinc-900 dark:text-white">{value}</p>
      <p className="text-xs text-zinc-500 mt-0.5">{label}</p>
      {sub && <p className="text-[10px] text-zinc-400 mt-1">{sub}</p>}
    </div>
  );
}

export default function SellerAnalyticsPage() {
  const { user } = useAuth();
  const [stats, setStats]       = useState<any>(null);
  const [products, setProducts] = useState<any[]>([]);
  const [loading, setLoading]   = useState(true);
  const [showUpgrade, setShowUpgrade] = useState(false);

  useEffect(() => {
    if (!user) return;
    Promise.all([
      api.getSellerStats(user.token),
      api.getSellerProductsDetailed(user.token),
    ]).then(([s, p]) => {
      setStats(s ?? {});
      setProducts(p ?? []);
    }).finally(() => setLoading(false));
  }, [user]);

  if (loading) return <div className="flex items-center justify-center h-64"><Spinner /></div>;

  const tier = stats?.tier ?? 'normal';

  // Non-premium sellers see a locked screen
  if (tier !== 'premium') {
    return (
      <div className="max-w-lg mx-auto py-12 px-4 text-center space-y-6">
        <div className="w-20 h-20 rounded-3xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center mx-auto">
          <Lock size={32} className="text-purple-400" />
        </div>
        <div>
          <h1 className="text-2xl font-black text-zinc-900 dark:text-white">Analytics</h1>
          <p className="text-sm text-zinc-500 mt-2">
            Advanced analytics is a <span className="font-bold text-purple-600 dark:text-purple-400">Premium</span> feature.
            Upgrade to unlock detailed insights about your listings, views, and revenue.
          </p>
        </div>

        {/* Preview of what's inside */}
        <div className="grid grid-cols-2 gap-3 text-left">
          {[
            { icon: TrendingUp,   label: 'Revenue trends',       desc: 'Weekly & monthly charts' },
            { icon: Package,      label: 'Listing performance',  desc: 'Views, clicks, conversions' },
            { icon: ShoppingCart, label: 'Order analytics',      desc: 'Acceptance rate & timing' },
            { icon: BarChart2,    label: 'Category insights',    desc: 'Best performing categories' },
          ].map(({ icon: Icon, label, desc }) => (
            <div key={label} className="bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl p-3 opacity-60">
              <Icon size={14} className="text-purple-400 mb-1.5" />
              <p className="text-xs font-bold text-zinc-700 dark:text-zinc-300">{label}</p>
              <p className="text-[10px] text-zinc-400 mt-0.5">{desc}</p>
            </div>
          ))}
        </div>

        <button onClick={() => setShowUpgrade(true)}
          className="w-full flex items-center justify-center gap-2 py-3.5 rounded-2xl bg-gradient-to-r from-purple-500 via-violet-500 to-indigo-500 text-white font-black text-sm shadow-lg shadow-purple-500/20 hover:opacity-90 transition-all">
          <Crown size={16} /> Upgrade to Premium — ETB 2,000
        </button>

        <Link href="/seller/dashboard" className="block text-xs text-zinc-400 hover:text-zinc-600 transition-colors">
          ← Back to Dashboard
        </Link>

        {showUpgrade && (
          <TierUpgradeModal
            currentTier={tier}
            onClose={() => setShowUpgrade(false)}
            onSuccess={() => setShowUpgrade(false)}
          />
        )}
      </div>
    );
  }

  // Premium analytics dashboard
  const approved  = products.filter(p => p.status === 'approved').length;
  const expired   = products.filter(p => p.status === 'expired').length;
  const totalFees = products.reduce((s, p) => s + (Number(p.total_fee_amount) || 0), 0);
  const paidFees  = products.filter(p => p.fee_paid).reduce((s, p) => s + (Number(p.total_fee_amount) || 0), 0);

  // Category breakdown
  const catMap: Record<string, number> = {};
  products.forEach(p => {
    const cat = p.category?.name ?? `Category #${p.category_id}`;
    catMap[cat] = (catMap[cat] ?? 0) + 1;
  });
  const topCats = Object.entries(catMap).sort((a, b) => b[1] - a[1]).slice(0, 5);

  return (
    <div className="space-y-6 max-w-4xl">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-black text-zinc-900 dark:text-white flex items-center gap-2">
            <BarChart2 size={20} className="text-purple-500" /> Analytics
          </h1>
          <p className="text-xs text-zinc-500 mt-0.5">Premium seller insights</p>
        </div>
        <span className="flex items-center gap-1.5 text-[10px] font-black px-3 py-1.5 rounded-full bg-purple-500/10 text-purple-500 border border-purple-500/20">
          <Crown size={10} /> Premium
        </span>
      </div>

      {/* KPI grid */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
        <StatCard icon={Package}      label="Total Listings"  value={stats?.total_products ?? 0}
          color="bg-blue-100 dark:bg-blue-900/30 text-blue-600" />
        <StatCard icon={TrendingUp}   label="Total Revenue"   value={`ETB ${(stats?.total_revenue ?? 0).toLocaleString()}`}
          color="bg-green-100 dark:bg-green-900/30 text-green-600" />
        <StatCard icon={DollarSign}   label="Fees Paid"       value={`ETB ${paidFees.toLocaleString()}`}
          sub={`of ETB ${totalFees.toLocaleString()} total`}
          color="bg-orange-100 dark:bg-orange-900/30 text-orange-600" />
        <StatCard icon={Clock}        label="Expired Listings" value={expired}
          sub={`${approved} currently active`}
          color="bg-zinc-100 dark:bg-zinc-800 text-zinc-600" />
      </div>

      {/* Status breakdown */}
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
        <h2 className="font-bold text-zinc-900 dark:text-white mb-4 flex items-center gap-2">
          <Package size={15} className="text-orange-600" /> Listing Status Breakdown
        </h2>
        <div className="space-y-2">
          {(['approved', 'pending', 'rejected', 'expired', 'sold'] as const).map(status => {
            const count = products.filter(p => p.status === status).length;
            const pct   = products.length > 0 ? Math.round((count / products.length) * 100) : 0;
            const colors: Record<string, string> = {
              approved: 'bg-green-500', pending: 'bg-yellow-500',
              rejected: 'bg-red-500',  expired: 'bg-zinc-400', sold: 'bg-emerald-500',
            };
            return (
              <div key={status} className="flex items-center gap-3">
                <span className="text-xs text-zinc-500 capitalize w-16 shrink-0">{status}</span>
                <div className="flex-1 h-2 bg-zinc-100 dark:bg-zinc-800 rounded-full overflow-hidden">
                  <div className={`h-full rounded-full ${colors[status]} transition-all`} style={{ width: `${pct}%` }} />
                </div>
                <span className="text-xs font-bold text-zinc-700 dark:text-zinc-300 w-8 text-right">{count}</span>
              </div>
            );
          })}
        </div>
      </div>

      {/* Top categories */}
      {topCats.length > 0 && (
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
          <h2 className="font-bold text-zinc-900 dark:text-white mb-4 flex items-center gap-2">
            <BarChart2 size={15} className="text-orange-600" /> Top Categories
          </h2>
          <div className="space-y-2">
            {topCats.map(([cat, count]) => {
              const pct = products.length > 0 ? Math.round((count / products.length) * 100) : 0;
              return (
                <div key={cat} className="flex items-center gap-3">
                  <span className="text-xs text-zinc-500 truncate flex-1 min-w-0">{cat}</span>
                  <div className="w-24 h-2 bg-zinc-100 dark:bg-zinc-800 rounded-full overflow-hidden shrink-0">
                    <div className="h-full rounded-full bg-orange-500 transition-all" style={{ width: `${pct}%` }} />
                  </div>
                  <span className="text-xs font-bold text-zinc-700 dark:text-zinc-300 w-6 text-right shrink-0">{count}</span>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Recent products table */}
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl overflow-hidden">
        <div className="flex items-center justify-between px-5 py-4 border-b border-zinc-100 dark:border-zinc-800">
          <h2 className="font-bold text-zinc-900 dark:text-white">Recent Listings</h2>
          <Link href="/seller/products" className="text-xs text-orange-600 font-bold hover:underline flex items-center gap-1">
            View all <ArrowUpRight size={11} />
          </Link>
        </div>
        <div className="divide-y divide-zinc-100 dark:divide-zinc-800">
          {products.slice(0, 8).map(p => (
            <div key={p.id} className="flex items-center gap-3 px-5 py-3">
              <div className="flex-1 min-w-0">
                <p className="text-sm font-semibold text-zinc-800 dark:text-zinc-200 truncate">{p.title}</p>
                <p className="text-xs text-zinc-400">ETB {Number(p.price).toLocaleString()}</p>
              </div>
              <span className={`text-[10px] font-black px-2 py-0.5 rounded-full uppercase shrink-0 ${
                p.status === 'approved' ? 'bg-green-100 dark:bg-green-900/30 text-green-600' :
                p.status === 'pending'  ? 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-600' :
                p.status === 'sold'     ? 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600' :
                'bg-zinc-100 dark:bg-zinc-800 text-zinc-500'
              }`}>{p.status}</span>
            </div>
          ))}
          {products.length === 0 && (
            <p className="px-5 py-8 text-center text-sm text-zinc-400">No listings yet</p>
          )}
        </div>
      </div>
    </div>
  );
}
