'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { DollarSign, TrendingUp, BarChart2, RefreshCw } from 'lucide-react';
import { Spinner } from '@/components/ui/Spinner';

function fmt(n: any) { return n != null ? `ETB ${Number(n).toLocaleString()}` : '—'; }
function fmtNum(n: any) { return n != null ? Number(n).toLocaleString() : '—'; }

export default function AdminProfitPage() {
  const { user, isSuperAdmin } = useAuth();
  const [summary, setSummary] = useState<any>(null);
  const [byTier, setByTier] = useState<any[]>([]);
  const [monthly, setMonthly] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  const load = async () => {
    if (!user?.token || !isSuperAdmin) return;
    setLoading(true);
    try {
      const [s, t, m] = await Promise.allSettled([
        api.adminProfitSummary(user.token),
        api.adminProfitByTier(user.token),
        api.adminProfitMonthly(user.token),
      ]);
      if (s.status === 'fulfilled') setSummary(s.value);
      if (t.status === 'fulfilled') setByTier(t.value ?? []);
      if (m.status === 'fulfilled') setMonthly(m.value ?? []);
    } finally { setLoading(false); }
  };

  useEffect(() => { load(); }, [user]); // eslint-disable-line

  if (!isSuperAdmin) return (
    <div className="flex flex-col items-center justify-center h-64 gap-3 text-zinc-500">
      <DollarSign size={32} className="opacity-20" />
      <p>Super admin access required</p>
    </div>
  );

  if (loading) return <div className="flex items-center justify-center h-64"><Spinner /></div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-emerald-600/10 border border-emerald-500/20 rounded-2xl flex items-center justify-center">
            <DollarSign size={22} className="text-emerald-400" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Profit & Revenue</h1>
            <p className="text-xs text-zinc-400">Platform financial summary — super admin only</p>
          </div>
        </div>
        <button onClick={load} className="flex items-center gap-2 px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-sm font-bold transition-colors">
          <RefreshCw size={14} /> Refresh
        </button>
      </div>

      {/* Summary KPIs */}
      {summary && (
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
          {[
            { label: 'Listing Fee Revenue', value: fmt(summary.total_listing_fee_revenue ?? summary.listing_fee_revenue), color: 'text-orange-400', bg: 'bg-orange-600/10' },
            { label: 'Broker Fee Revenue', value: fmt(summary.total_broker_revenue ?? summary.broker_revenue), color: 'text-blue-400', bg: 'bg-blue-600/10' },
            { label: 'Net Platform Profit', value: fmt(summary.total_platform_profit ?? summary.net_profit), color: 'text-emerald-400', bg: 'bg-emerald-600/10' },
            { label: 'Total Transactions', value: fmtNum(summary.total_transactions ?? summary.transaction_count), color: 'text-purple-400', bg: 'bg-purple-600/10' },
          ].map(({ label, value, color, bg }) => (
            <div key={label} className="bg-zinc-900 border border-zinc-800 rounded-2xl p-5">
              <div className={`w-8 h-8 rounded-xl ${bg} flex items-center justify-center mb-3`}>
                <DollarSign size={16} className={color} />
              </div>
              <p className={`text-2xl font-black font-mono ${color}`}>{value}</p>
              <p className="text-xs text-zinc-500 mt-1">{label}</p>
            </div>
          ))}
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* By tier */}
        {byTier.length > 0 && (
          <div className="bg-zinc-900 border border-zinc-800 rounded-2xl overflow-hidden">
            <div className="px-5 py-4 border-b border-zinc-800 flex items-center gap-2">
              <BarChart2 size={15} className="text-orange-500" />
              <h2 className="font-black text-sm text-white uppercase tracking-wider">Revenue by Seller Tier</h2>
            </div>
            <table className="w-full text-sm">
              <thead><tr className="border-b border-zinc-800 text-zinc-500 text-xs uppercase font-bold">
                <th className="px-5 py-3 text-left">Tier</th>
                <th className="px-5 py-3 text-right">Sellers</th>
                <th className="px-5 py-3 text-right">Revenue</th>
              </tr></thead>
              <tbody className="divide-y divide-zinc-800/50">
                {byTier.map((row, i) => (
                  <tr key={i} className="hover:bg-zinc-800/30">
                    <td className="px-5 py-3 capitalize font-bold text-white">{row.tier || row.seller_tier || '—'}</td>
                    <td className="px-5 py-3 text-right font-mono text-zinc-300">{fmtNum(row.seller_count ?? row.count)}</td>
                    <td className="px-5 py-3 text-right font-mono font-bold text-emerald-400">{fmt(row.total_revenue ?? row.revenue)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* Monthly trend */}
        {monthly.length > 0 && (
          <div className="bg-zinc-900 border border-zinc-800 rounded-2xl overflow-hidden">
            <div className="px-5 py-4 border-b border-zinc-800 flex items-center gap-2">
              <TrendingUp size={15} className="text-orange-500" />
              <h2 className="font-black text-sm text-white uppercase tracking-wider">Monthly Trend</h2>
            </div>
            <table className="w-full text-sm">
              <thead><tr className="border-b border-zinc-800 text-zinc-500 text-xs uppercase font-bold">
                <th className="px-5 py-3 text-left">Month</th>
                <th className="px-5 py-3 text-right">Listings</th>
                <th className="px-5 py-3 text-right">Revenue</th>
              </tr></thead>
              <tbody className="divide-y divide-zinc-800/50">
                {monthly.map((row, i) => (
                  <tr key={i} className="hover:bg-zinc-800/30">
                    <td className="px-5 py-3 font-bold text-white">{row.month ?? row.period}</td>
                    <td className="px-5 py-3 text-right font-mono text-zinc-300">{fmtNum(row.new_products ?? row.listings)}</td>
                    <td className="px-5 py-3 text-right font-mono font-bold text-emerald-400">{fmt(row.revenue ?? row.total_revenue)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
