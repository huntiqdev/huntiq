'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Gift, RefreshCw, CheckCircle, XCircle, DollarSign } from 'lucide-react';
import { Spinner } from '@/components/ui/Spinner';

const WITHDRAWAL_TABS = ['pending', 'approved', 'rejected', 'paid'];

export default function AdminReferralsPage() {
  const { user, isSuperAdmin } = useAuth();
  const [overview, setOverview] = useState<any>(null);
  const [withdrawals, setWithdrawals] = useState<any[]>([]);
  const [wStatus, setWStatus] = useState('pending');
  const [loading, setLoading] = useState(true);

  const load = async () => {
    if (!user?.token || !isSuperAdmin) return;
    setLoading(true);
    try {
      const [ov, wd] = await Promise.allSettled([
        api.adminGetReferralOverview(user.token),
        api.adminListReferralWithdrawals(wStatus, user.token),
      ]);
      if (ov.status === 'fulfilled') setOverview(ov.value);
      if (wd.status === 'fulfilled') setWithdrawals(wd.value ?? []);
    } finally { setLoading(false); }
  };

  useEffect(() => { load(); }, [user, wStatus]); // eslint-disable-line

  const approve = async (id: number) => {
    if (!user?.token) return;
    await api.adminApproveReferralWithdrawal(id, '', user.token);
    load();
  };

  const reject = async (id: number) => {
    if (!user?.token) return;
    const reason = prompt('Rejection reason:');
    if (reason === null) return;
    await api.adminRejectReferralWithdrawal(id, reason, user.token);
    load();
  };

  const markPaid = async (id: number) => {
    if (!user?.token) return;
    const ref = prompt('Payment reference / transaction ID:') || 'MANUAL-PAID';
    await api.adminMarkReferralWithdrawalPaid(id, ref, user.token);
    load();
  };

  if (!isSuperAdmin) return (
    <div className="flex flex-col items-center justify-center h-64 gap-3 text-zinc-500">
      <Gift size={32} className="opacity-20" />
      <p>Super admin access required</p>
    </div>
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-green-600/10 border border-green-500/20 rounded-2xl flex items-center justify-center">
            <Gift size={22} className="text-green-400" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Referral Program</h1>
            <p className="text-xs text-zinc-400">Commissions, withdrawals, and referral management</p>
          </div>
        </div>
        <button onClick={load} className="flex items-center gap-2 px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-sm font-bold transition-colors">
          <RefreshCw size={14} /> Refresh
        </button>
      </div>

      {/* Overview KPIs */}
      {overview && (
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
          {[
            { label: 'Total Referrers', value: overview.total_referrers?.toLocaleString() ?? '—', color: 'text-blue-400' },
            { label: 'Total Commissions', value: `ETB ${Number(overview.total_commissions ?? 0).toLocaleString()}`, color: 'text-green-400' },
            { label: 'Pending Payout', value: `ETB ${Number(overview.pending_payout ?? 0).toLocaleString()}`, color: 'text-amber-400' },
            { label: 'Total Withdrawn', value: `ETB ${Number(overview.total_withdrawn ?? 0).toLocaleString()}`, color: 'text-emerald-400' },
          ].map(({ label, value, color }) => (
            <div key={label} className="bg-zinc-900 border border-zinc-800 rounded-2xl p-4">
              <DollarSign size={16} className={`${color} mb-2`} />
              <p className={`text-xl font-black font-mono ${color}`}>{value}</p>
              <p className="text-xs text-zinc-500 mt-1">{label}</p>
            </div>
          ))}
        </div>
      )}

      {/* Withdrawal requests */}
      <div>
        <div className="flex items-center justify-between mb-3">
          <h2 className="font-black text-sm text-white uppercase tracking-wider">Withdrawal Requests</h2>
          <div className="flex gap-1">
            {WITHDRAWAL_TABS.map(s => (
              <button key={s} onClick={() => setWStatus(s)}
                className={`px-3 py-1 rounded-lg text-xs font-bold capitalize transition-all ${wStatus === s ? 'bg-orange-600 text-white' : 'bg-zinc-800 text-zinc-400 hover:bg-zinc-700'}`}>
                {s}
              </button>
            ))}
          </div>
        </div>

        <div className="bg-zinc-900 border border-zinc-800 rounded-2xl overflow-hidden">
          {loading ? (
            <div className="p-8 flex justify-center"><Spinner /></div>
          ) : withdrawals.length === 0 ? (
            <div className="p-8 text-center text-zinc-500">
              <Gift size={28} className="mx-auto mb-2 opacity-20" />
              <p>No {wStatus} withdrawal requests</p>
            </div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-800 text-zinc-500 text-xs uppercase font-bold">
                  <th className="px-5 py-3 text-left">User</th>
                  <th className="px-5 py-3 text-left">Amount</th>
                  <th className="px-5 py-3 text-left">Method</th>
                  <th className="px-5 py-3 text-left">Phone</th>
                  <th className="px-5 py-3 text-left">Status</th>
                  <th className="px-5 py-3 text-left">Date</th>
                  <th className="px-5 py-3 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-800/50">
                {withdrawals.map((w: any, i: number) => (
                  <tr key={i} className="hover:bg-zinc-800/30">
                    <td className="px-5 py-3">
                      <p className="font-bold text-white">{w.user_name || `User #${w.user_id}`}</p>
                      {w.user_phone && <p className="text-xs text-zinc-500 font-mono">{w.user_phone}</p>}
                    </td>
                    <td className="px-5 py-3 font-mono font-bold text-emerald-400">ETB {Number(w.amount).toLocaleString()}</td>
                    <td className="px-5 py-3 text-zinc-400 capitalize text-xs">{w.payment_method}</td>
                    <td className="px-5 py-3 font-mono text-xs text-zinc-400">{w.payer_phone}</td>
                    <td className="px-5 py-3">
                      <span className={`px-2 py-0.5 rounded-full text-xs font-bold capitalize ${
                        w.status === 'paid'     ? 'bg-emerald-900/40 text-emerald-400' :
                        w.status === 'approved' ? 'bg-blue-900/40 text-blue-400' :
                        w.status === 'rejected' ? 'bg-red-900/40 text-red-400' :
                        'bg-amber-900/40 text-amber-400'
                      }`}>{w.status}</span>
                    </td>
                    <td className="px-5 py-3 text-xs text-zinc-500">
                      {w.requested_at ? new Date(w.requested_at).toLocaleDateString() : '—'}
                    </td>
                    <td className="px-5 py-3 text-right">
                      <div className="flex items-center justify-end gap-1">
                        {w.status === 'pending' && <>
                          <button onClick={() => approve(w.id)} className="px-2 py-1 bg-emerald-600/20 text-emerald-400 border border-emerald-500/20 rounded-lg text-xs font-bold hover:bg-emerald-600/30 flex items-center gap-1">
                            <CheckCircle size={10} /> Approve
                          </button>
                          <button onClick={() => reject(w.id)} className="px-2 py-1 bg-red-600/20 text-red-400 border border-red-500/20 rounded-lg text-xs font-bold hover:bg-red-600/30 flex items-center gap-1">
                            <XCircle size={10} /> Reject
                          </button>
                        </>}
                        {w.status === 'approved' && (
                          <button onClick={() => markPaid(w.id)} className="px-2 py-1 bg-blue-600/20 text-blue-400 border border-blue-500/20 rounded-lg text-xs font-bold hover:bg-blue-600/30 flex items-center gap-1">
                            <DollarSign size={10} /> Mark Paid
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}
