'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { TrendingUp, CheckCircle, XCircle, RefreshCw, Eye } from 'lucide-react';
import { Spinner } from '@/components/ui/Spinner';

const STATUS_TABS = ['pending', 'approved', 'rejected'];

export default function AdminUpgradesPage() {
  const { user } = useAuth();
  const [sellerStatus, setSellerStatus] = useState('pending');
  const [buyerStatus, setBuyerStatus] = useState('pending');
  const [sellerReqs, setSellerReqs] = useState<any[]>([]);
  const [buyerReqs, setBuyerReqs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const load = async () => {
    if (!user?.token) return;
    setLoading(true);
    setError('');
    try {
      const [s, b] = await Promise.allSettled([
        api.adminListUpgradeRequests(sellerStatus, user.token),
        api.adminListBuyerUpgradeRequests(buyerStatus, user.token),
      ]);
      if (s.status === 'fulfilled') setSellerReqs(s.value ?? []);
      if (b.status === 'fulfilled') {
        const res = b.value as any;
        setBuyerReqs(Array.isArray(res) ? res : (res?.requests ?? []));
      }
    } catch (err: any) { setError(err.message); }
    finally { setLoading(false); }
  };

  useEffect(() => { load(); }, [user, sellerStatus, buyerStatus]); // eslint-disable-line

  const approveSeller = async (id: number) => {
    if (!user?.token) return;
    await api.adminApproveUpgrade(id, user.token);
    load();
  };

  const rejectSeller = async (id: number) => {
    if (!user?.token) return;
    const reason = prompt('Rejection reason:');
    if (reason === null) return;
    await api.adminRejectUpgrade(id, reason, user.token);
    load();
  };

  const approveBuyer = async (id: number) => {
    if (!user?.token) return;
    await api.adminApproveBuyerUpgrade(id, user.token);
    load();
  };

  const rejectBuyer = async (id: number) => {
    if (!user?.token) return;
    const reason = prompt('Rejection reason:');
    if (reason === null) return;
    await api.adminRejectBuyerUpgrade(id, reason, user.token);
    load();
  };

  const tierBadge = (tier: string) => ({
    gold:    'bg-yellow-900/40 text-yellow-400 border border-yellow-500/20',
    premium: 'bg-purple-900/40 text-purple-400 border border-purple-500/20',
  }[tier] ?? 'bg-zinc-700 text-zinc-300');

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-yellow-600/10 border border-yellow-500/20 rounded-2xl flex items-center justify-center">
            <TrendingUp size={22} className="text-yellow-400" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Upgrade Requests</h1>
            <p className="text-xs text-zinc-400">Seller and buyer tier upgrade approvals</p>
          </div>
        </div>
        <button onClick={load} className="flex items-center gap-2 px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-sm font-bold transition-colors">
          <RefreshCw size={14} /> Refresh
        </button>
      </div>

      {error && <p className="text-red-400 text-sm bg-red-900/20 border border-red-500/30 rounded-xl px-4 py-3">{error}</p>}

      {/* Seller upgrades */}
      <section>
        <div className="flex items-center justify-between mb-3">
          <h2 className="font-black text-sm text-white uppercase tracking-wider">Seller Upgrades</h2>
          <div className="flex gap-1">
            {STATUS_TABS.map(s => (
              <button key={s} onClick={() => setSellerStatus(s)}
                className={`px-3 py-1 rounded-lg text-xs font-bold capitalize transition-all ${sellerStatus === s ? 'bg-orange-600 text-white' : 'bg-zinc-800 text-zinc-400 hover:bg-zinc-700'}`}>
                {s}
              </button>
            ))}
          </div>
        </div>
        <RequestTable
          loading={loading} data={sellerReqs} mode="seller"
          onApprove={approveSeller} onReject={rejectSeller} tierBadge={tierBadge}
        />
      </section>

      {/* Buyer upgrades */}
      <section>
        <div className="flex items-center justify-between mb-3">
          <h2 className="font-black text-sm text-white uppercase tracking-wider">Buyer Upgrades</h2>
          <div className="flex gap-1">
            {STATUS_TABS.map(s => (
              <button key={s} onClick={() => setBuyerStatus(s)}
                className={`px-3 py-1 rounded-lg text-xs font-bold capitalize transition-all ${buyerStatus === s ? 'bg-orange-600 text-white' : 'bg-zinc-800 text-zinc-400 hover:bg-zinc-700'}`}>
                {s}
              </button>
            ))}
          </div>
        </div>
        <RequestTable
          loading={loading} data={buyerReqs} mode="buyer"
          onApprove={approveBuyer} onReject={rejectBuyer} tierBadge={tierBadge}
        />
      </section>
    </div>
  );
}

function RequestTable({ loading, data, mode, onApprove, onReject, tierBadge }: {
  loading: boolean; data: any[]; mode: 'seller' | 'buyer';
  onApprove: (id: number) => void; onReject: (id: number) => void;
  tierBadge: (tier: string) => string;
}) {
  if (loading) return <div className="p-6 flex justify-center"><Spinner /></div>;
  if (data.length === 0) return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-6 text-center text-zinc-500 text-sm">
      No {mode} upgrade requests
    </div>
  );

  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-2xl overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-zinc-800 text-zinc-500 text-xs uppercase font-bold">
            <th className="px-5 py-3 text-left">{mode === 'seller' ? 'Seller' : 'Buyer'}</th>
            <th className="px-5 py-3 text-left">Target Tier</th>
            <th className="px-5 py-3 text-left">Amount</th>
            <th className="px-5 py-3 text-left">Method</th>
            <th className="px-5 py-3 text-left">Reference</th>
            <th className="px-5 py-3 text-left">Date</th>
            <th className="px-5 py-3 text-left">Status</th>
            <th className="px-5 py-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-zinc-800/50">
          {data.map((r: any, i: number) => (
            <tr key={i} className="hover:bg-zinc-800/30">
              <td className="px-5 py-3 font-mono text-xs text-zinc-400">
                #{r.user_id ?? r.seller_id ?? r.buyer_id}
                {r.user_name && <span className="text-white ml-1">{r.user_name}</span>}
              </td>
              <td className="px-5 py-3">
                <span className={`px-2 py-0.5 rounded-full text-xs font-bold capitalize ${tierBadge(r.requested_tier)}`}>
                  {r.requested_tier}
                </span>
              </td>
              <td className="px-5 py-3 font-mono font-bold text-emerald-400">
                ETB {Number(r.amount).toLocaleString()}
              </td>
              <td className="px-5 py-3 text-zinc-400 capitalize text-xs">{r.payment_method}</td>
              <td className="px-5 py-3 font-mono text-xs text-zinc-400 max-w-[120px] truncate">{r.reference || '—'}</td>
              <td className="px-5 py-3 text-xs text-zinc-500">
                {r.created_at ? new Date(r.created_at).toLocaleDateString() : '—'}
              </td>
              <td className="px-5 py-3">
                <span className={`px-2 py-0.5 rounded-full text-xs font-bold capitalize ${
                  r.status === 'approved' ? 'bg-emerald-900/40 text-emerald-400' :
                  r.status === 'rejected' ? 'bg-red-900/40 text-red-400' :
                  'bg-amber-900/40 text-amber-400'
                }`}>{r.status}</span>
              </td>
              <td className="px-5 py-3 text-right">
                {r.status === 'pending' && (
                  <div className="flex items-center justify-end gap-2">
                    <button onClick={() => onApprove(r.id)}
                      className="px-3 py-1 bg-emerald-600/20 text-emerald-400 border border-emerald-500/20 rounded-lg text-xs font-bold hover:bg-emerald-600/30 flex items-center gap-1">
                      <CheckCircle size={11} /> Approve
                    </button>
                    <button onClick={() => onReject(r.id)}
                      className="px-3 py-1 bg-red-600/20 text-red-400 border border-red-500/20 rounded-lg text-xs font-bold hover:bg-red-600/30 flex items-center gap-1">
                      <XCircle size={11} /> Reject
                    </button>
                  </div>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
