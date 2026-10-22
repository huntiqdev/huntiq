'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Crown, Star, RefreshCw, ShieldAlert } from 'lucide-react';
import { Spinner } from '@/components/ui/Spinner';

export default function AdminPremiumUsersPage() {
  const { user, isSuperAdmin, hasAdminPermission } = useAuth();
  const [users, setUsers] = useState<any[]>([]);
  const [role, setRole] = useState('seller');
  const [tierType, setTierType] = useState('seller');
  const [tier, setTier] = useState('gold');
  const [loading, setLoading] = useState(true);

  const canAccess = isSuperAdmin || hasAdminPermission('premium_users');

  const load = async () => {
    if (!user?.token || !canAccess) return;
    setLoading(true);
    try {
      const data = await api.adminGetPremiumUsers({ role, tier_type: tierType, tier }, user.token);
      setUsers(data ?? []);
    } catch { setUsers([]); }
    finally { setLoading(false); }
  };

  useEffect(() => { load(); }, [user, role, tierType, tier]); // eslint-disable-line

  if (!canAccess) return (
    <div className="flex flex-col items-center justify-center h-64 gap-3 text-zinc-500">
      <ShieldAlert size={32} className="opacity-20" />
      <p>Insufficient permissions</p>
    </div>
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-yellow-600/10 border border-yellow-500/20 rounded-2xl flex items-center justify-center">
            <Crown size={22} className="text-yellow-400" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Premium Users</h1>
            <p className="text-xs text-zinc-400">Gold and premium tier accounts</p>
          </div>
        </div>
        <button onClick={load} className="flex items-center gap-2 px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-sm font-bold transition-colors">
          <RefreshCw size={14} /> Refresh
        </button>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap gap-3">
        <div className="flex bg-zinc-900 border border-zinc-800 rounded-xl p-1 gap-1">
          <span className="text-xs font-bold text-zinc-500 px-2 flex items-center">Role:</span>
          {['seller', 'buyer'].map(r => (
            <button key={r} onClick={() => setRole(r)}
              className={`px-3 py-1.5 rounded-lg text-xs font-bold capitalize transition-all ${role === r ? 'bg-orange-600 text-white' : 'text-zinc-400 hover:text-white'}`}>
              {r}
            </button>
          ))}
        </div>
        <div className="flex bg-zinc-900 border border-zinc-800 rounded-xl p-1 gap-1">
          <span className="text-xs font-bold text-zinc-500 px-2 flex items-center">Tier:</span>
          {['gold', 'premium'].map(t => (
            <button key={t} onClick={() => setTier(t)}
              className={`px-3 py-1.5 rounded-lg text-xs font-bold capitalize transition-all ${tier === t
                ? t === 'gold' ? 'bg-yellow-500 text-white' : 'bg-purple-600 text-white'
                : 'text-zinc-400 hover:text-white'}`}>
              {t === 'gold' ? <><Star size={10} className="inline mr-1" />Gold</> : <><Crown size={10} className="inline mr-1" />Premium</>}
            </button>
          ))}
        </div>
      </div>

      <div className="bg-zinc-900 border border-zinc-800 rounded-2xl overflow-hidden">
        {loading ? (
          <div className="p-8 flex justify-center"><Spinner /></div>
        ) : users.length === 0 ? (
          <div className="p-8 text-center text-zinc-500">
            <Crown size={28} className="mx-auto mb-2 opacity-20" />
            <p>No {tier} {role}s found</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-zinc-800 text-zinc-500 text-xs uppercase font-bold">
                <th className="px-5 py-3 text-left">User</th>
                <th className="px-5 py-3 text-left">Phone</th>
                <th className="px-5 py-3 text-left">Tier</th>
                <th className="px-5 py-3 text-right">Value Sold</th>
                <th className="px-5 py-3 text-left">Joined</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/50">
              {users.map((u: any, i: number) => (
                <tr key={i} className="hover:bg-zinc-800/30">
                  <td className="px-5 py-3">
                    <p className="font-bold text-white">{u.name || `User #${u.id}`}</p>
                    {u.email && !u.email.includes('@noemail.local') && (
                      <p className="text-xs text-zinc-500">{u.email}</p>
                    )}
                  </td>
                  <td className="px-5 py-3 font-mono text-xs text-zinc-400">{u.phone}</td>
                  <td className="px-5 py-3">
                    <span className={`px-2 py-0.5 rounded-full text-xs font-bold capitalize ${
                      (u.tier === 'premium' || u.buyer_tier === 'premium')
                        ? 'bg-purple-900/40 text-purple-400 border border-purple-500/20'
                        : 'bg-yellow-900/40 text-yellow-400 border border-yellow-500/20'
                    }`}>
                      {u.tier || u.buyer_tier}
                    </span>
                  </td>
                  <td className="px-5 py-3 text-right font-mono text-emerald-400">
                    {u.total_value_sold ? `ETB ${Number(u.total_value_sold).toLocaleString()}` : '—'}
                  </td>
                  <td className="px-5 py-3 text-xs text-zinc-500">
                    {u.created_at ? new Date(u.created_at).toLocaleDateString() : '—'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
