'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Search, TrendingUp, AlertCircle, RefreshCw, BarChart2 } from 'lucide-react';

export default function AdminSearchPage() {
  const { user } = useAuth();
  const [popular, setPopular] = useState<any[]>([]);
  const [zeroResults, setZeroResults] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const load = async (silent = false) => {
    if (!user?.token) return;
    if (!silent) setLoading(true);
    else setRefreshing(true);
    try {
      const [p, z] = await Promise.allSettled([
        api.adminGetPopularSearches(30, user.token),
        api.adminGetZeroResultQueries(50, user.token),
      ]);
      if (p.status === 'fulfilled') setPopular(p.value ?? []);
      if (z.status === 'fulfilled') setZeroResults(z.value ?? []);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => { load(); }, [user]); // eslint-disable-line

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-zinc-800 rounded-2xl animate-pulse" />
          <div className="space-y-2">
            <div className="h-5 w-40 bg-zinc-800 rounded animate-pulse" />
            <div className="h-3 w-56 bg-zinc-800 rounded animate-pulse" />
          </div>
        </div>
        <div className="grid grid-cols-2 gap-4">
          {[...Array(2)].map((_, i) => (
            <div key={i} className="h-64 bg-zinc-900 border border-zinc-800 rounded-2xl animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-blue-600/10 border border-blue-500/20 rounded-2xl flex items-center justify-center">
            <Search size={22} className="text-blue-400" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Search Analytics</h1>
            <p className="text-xs text-zinc-400">Popular queries and zero-result searches</p>
          </div>
        </div>
        <button
          onClick={() => load(true)}
          disabled={refreshing}
          className="flex items-center gap-2 px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-sm font-bold transition-colors disabled:opacity-50"
        >
          <RefreshCw size={14} className={refreshing ? 'animate-spin' : ''} />
          Refresh
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Popular searches */}
        <div className="bg-zinc-900 border border-zinc-800 rounded-2xl overflow-hidden">
          <div className="px-5 py-4 border-b border-zinc-800 flex items-center gap-2">
            <TrendingUp size={15} className="text-orange-500" />
            <h2 className="font-black text-sm text-white uppercase tracking-wider">
              Popular Searches ({popular.length})
            </h2>
          </div>
          {popular.length === 0 ? (
            <div className="p-8 text-center text-zinc-500">
              <Search size={28} className="mx-auto mb-2 opacity-20" />
              <p className="text-sm">No search data yet — queries are recorded when users search</p>
            </div>
          ) : (
            <div className="divide-y divide-zinc-800/50 max-h-96 overflow-y-auto">
              {popular.map((item, i) => (
                <div key={i} className="flex items-center gap-4 px-5 py-3 hover:bg-zinc-800/30">
                  <span className="text-xs font-black text-zinc-600 w-6 shrink-0">{i + 1}</span>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-bold text-white truncate">{item.query}</p>
                    <p className="text-xs text-zinc-500 mt-0.5">
                      avg {Number(item.avg_results ?? 0).toFixed(0)} results
                    </p>
                  </div>
                  <div className="text-right shrink-0">
                    <p className="text-sm font-black text-orange-400 font-mono">{item.search_count?.toLocaleString()}</p>
                    <p className="text-[10px] text-zinc-500">searches</p>
                  </div>
                  {item.ctr_percent != null && (
                    <div className="text-right shrink-0 w-12">
                      <p className="text-xs font-bold text-emerald-400">{Number(item.ctr_percent).toFixed(1)}%</p>
                      <p className="text-[10px] text-zinc-500">CTR</p>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Zero-result queries */}
        <div className="bg-zinc-900 border border-red-500/20 rounded-2xl overflow-hidden">
          <div className="px-5 py-4 border-b border-red-500/20 bg-red-900/10 flex items-center gap-2">
            <AlertCircle size={15} className="text-red-400" />
            <h2 className="font-black text-sm text-red-300 uppercase tracking-wider">
              Zero-Result Queries ({zeroResults.length})
            </h2>
          </div>
          {zeroResults.length === 0 ? (
            <div className="p-8 text-center text-zinc-500">
              <AlertCircle size={28} className="mx-auto mb-2 opacity-20" />
              <p className="text-sm">No zero-result queries — great search coverage!</p>
            </div>
          ) : (
            <>
              <div className="px-5 py-3 bg-red-900/5 border-b border-zinc-800">
                <p className="text-xs text-zinc-400">
                  These queries returned no results. Consider adding synonyms, expanding category coverage, or seeding matching products.
                </p>
              </div>
              <div className="divide-y divide-zinc-800/50 max-h-96 overflow-y-auto">
                {zeroResults.map((item, i) => (
                  <div key={i} className="flex items-center gap-4 px-5 py-3 hover:bg-zinc-800/30">
                    <span className="text-xs font-black text-zinc-600 w-6 shrink-0">{i + 1}</span>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-bold text-white truncate">{item.query}</p>
                      {item.normalized_query && item.normalized_query !== item.query && (
                        <p className="text-xs text-zinc-500 mt-0.5 font-mono">→ {item.normalized_query}</p>
                      )}
                    </div>
                    <div className="text-right shrink-0">
                      <p className="text-sm font-black text-red-400 font-mono">{item.occurrence_count?.toLocaleString()}</p>
                      <p className="text-[10px] text-zinc-500">times</p>
                    </div>
                  </div>
                ))}
              </div>
            </>
          )}
        </div>
      </div>

      {/* Search quality summary */}
      <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-5">
        <div className="flex items-center gap-2 mb-4">
          <BarChart2 size={15} className="text-orange-500" />
          <h2 className="font-black text-sm text-white uppercase tracking-wider">Search Quality Summary</h2>
        </div>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
          {[
            { label: 'Popular Terms', value: popular.length, color: 'text-orange-400' },
            { label: 'Zero-Result Terms', value: zeroResults.length, color: 'text-red-400' },
            {
              label: 'Total Searches',
              value: popular.reduce((s: number, p: any) => s + (p.search_count || 0), 0).toLocaleString(),
              color: 'text-blue-400',
            },
            {
              label: 'Avg CTR',
              value: popular.length > 0
                ? `${(popular.reduce((s: number, p: any) => s + (Number(p.ctr_percent) || 0), 0) / popular.length).toFixed(1)}%`
                : '—',
              color: 'text-emerald-400',
            },
          ].map(({ label, value, color }) => (
            <div key={label} className="bg-zinc-800 rounded-xl p-3">
              <p className="text-xs text-zinc-500 font-bold uppercase tracking-wider">{label}</p>
              <p className={`text-2xl font-black font-mono mt-1 ${color}`}>{value}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
