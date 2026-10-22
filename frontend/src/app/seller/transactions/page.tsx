'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { TrendingUp, Search, DollarSign, Package, Lock, Star } from 'lucide-react';
import Link from 'next/link';
import { INPUT_CLS } from '@/lib/constants/styles';

export default function SellerTransactionsPage() {
  const { user } = useAuth();
  const [txns, setTxns] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [tierError, setTierError] = useState(false);

  useEffect(() => {
    if (!user) return;
    api.getSellerTransactions(user.token)
      .then(r => setTxns(r ?? []))
      .catch(err => {
        if (err?.message?.includes('gold tier') || err?.message?.includes('requires gold')) {
          setTierError(true);
        }
      })
      .finally(() => setLoading(false));
  }, [user]);

  if (loading) return <div className="flex items-center justify-center h-64"><Spinner /></div>;

  if (tierError) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[50vh] gap-4 text-center px-4">
        <div className="w-16 h-16 rounded-2xl bg-yellow-100 dark:bg-yellow-900/30 flex items-center justify-center">
          <Lock size={28} className="text-yellow-500" />
        </div>
        <div>
          <h2 className="text-lg font-black text-zinc-900 dark:text-white">Gold Tier Required</h2>
          <p className="text-sm text-zinc-500 mt-1 max-w-xs">
            Sales history is available for Gold and Premium sellers. Upgrade your account to unlock this feature.
          </p>
        </div>
        <Link href="/seller/pay?upgrade=1"
          className="flex items-center gap-2 bg-yellow-500 hover:bg-yellow-400 text-white px-6 py-3 rounded-xl font-bold text-sm transition-colors">
          <Star size={15} /> Upgrade to Gold
        </Link>
      </div>
    );
  }

  const filtered = txns.filter(tx =>
    !search || [tx.buyer_name, tx.buyer_phone].some((v: string) => v?.toLowerCase().includes(search.toLowerCase()))
  );

  const totalEarnings = txns.reduce((s, tx) => s + (Number(tx.seller_gets) || 0), 0);
  const totalVolume = txns.reduce((s, tx) => s + (Number(tx.amount) || 0), 0);
  const totalFees = txns.reduce((s, tx) => s + (Number(tx.broker_fee) || 0), 0);

  return (
    <div className="space-y-5 max-w-4xl">
      <div>
        <h1 className="text-xl font-black text-zinc-900 dark:text-white">Sales History</h1>
        <p className="text-xs text-zinc-500 mt-0.5">{txns.length} total transactions</p>
      </div>

      {/* Summary */}
      <div className="grid grid-cols-3 gap-3">
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4">
          <div className="flex items-center gap-2 mb-1">
            <TrendingUp size={14} className="text-green-600" />
            <p className="text-xs text-zinc-500">Your Earnings</p>
          </div>
          <p className="text-xl font-black text-green-600">ETB {totalEarnings.toLocaleString()}</p>
          <p className="text-[10px] text-zinc-400 mt-0.5">After 2% broker fee</p>
        </div>
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4">
          <div className="flex items-center gap-2 mb-1">
            <DollarSign size={14} className="text-blue-600" />
            <p className="text-xs text-zinc-500">Total Volume</p>
          </div>
          <p className="text-xl font-black text-zinc-900 dark:text-white">ETB {totalVolume.toLocaleString()}</p>
        </div>
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4">
          <div className="flex items-center gap-2 mb-1">
            <Package size={14} className="text-orange-600" />
            <p className="text-xs text-zinc-500">Broker Fees</p>
          </div>
          <p className="text-xl font-black text-orange-600">ETB {totalFees.toLocaleString()}</p>
        </div>
      </div>

      {/* Search */}
      <div className="relative max-w-xs">
        <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none" />
        <input value={search} onChange={e => setSearch(e.target.value)}
          placeholder="Search by buyer..."
          className={`pl-9 pr-3 py-2 ${INPUT_CLS}`} />
      </div>

      {/* Transactions */}
      {filtered.length === 0 ? (
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-12 text-center">
          <DollarSign size={32} className="mx-auto mb-3 text-zinc-300 dark:text-zinc-600" />
          <p className="text-zinc-400">No sales yet</p>
        </div>
      ) : (
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-zinc-50 dark:bg-zinc-800 border-b border-zinc-100 dark:border-zinc-700">
                <tr>
                  {['Product', 'Buyer', 'Phone', 'Amount', 'Fee', 'You Get', 'Status', 'Date'].map(h => (
                    <th key={h} className="px-4 py-3 text-left text-[10px] font-black uppercase tracking-widest text-zinc-400">{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-100 dark:divide-zinc-800">
                {filtered.map(tx => (
                  <tr key={tx.id} className="hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors">
                    <td className="px-4 py-3 text-zinc-400 text-xs">#{tx.product_id}</td>
                    <td className="px-4 py-3 font-semibold text-zinc-800 dark:text-zinc-200">{tx.buyer_name}</td>
                    <td className="px-4 py-3 text-zinc-500 text-xs">{tx.buyer_phone}</td>
                    <td className="px-4 py-3 font-bold text-zinc-700 dark:text-zinc-300">ETB {Number(tx.amount).toLocaleString()}</td>
                    <td className="px-4 py-3 text-red-500 text-xs">-ETB {Number(tx.broker_fee).toLocaleString()}</td>
                    <td className="px-4 py-3 font-black text-green-600">ETB {Number(tx.seller_gets).toLocaleString()}</td>
                    <td className="px-4 py-3">
                      <span className={`text-[10px] font-black px-2 py-0.5 rounded-full uppercase ${
                        tx.status === 'completed' ? 'bg-green-100 dark:bg-green-900/30 text-green-600' : 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-600'
                      }`}>{tx.status}</span>
                    </td>
                    <td className="px-4 py-3 text-zinc-400 text-xs">{new Date(tx.created_at).toLocaleDateString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
}
