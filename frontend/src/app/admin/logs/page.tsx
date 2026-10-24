'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { FileText, RefreshCw, Trash2, AlertCircle } from 'lucide-react';
import { Spinner } from '@/components/ui/Spinner';

export default function AdminLogsPage() {
  const { user, isSuperAdmin, isSecurityAdmin, isPlatformOwner, hasAdminPermission } = useAuth();
  const [logs, setLogs] = useState<any[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [error, setError] = useState('');
  const limit = 50;

  const canViewLogs = isPlatformOwner || isSecurityAdmin || isSuperAdmin || hasAdminPermission('logs');

  const load = async (p = 1) => {
    if (!user?.token || !canViewLogs) return;
    setLoading(true);
    setError('');
    try {
      const res = await api.adminGetLogs({ page: p, limit }, user.token);
      setLogs(res?.logs ?? []);
      setTotal(res?.total ?? 0);
    } catch (err: any) { setError(err.message); }
    finally { setLoading(false); }
  };

  useEffect(() => { load(page); }, [user, page]); // eslint-disable-line

  if (!canViewLogs) return (
    <div className="flex flex-col items-center justify-center h-64 gap-3 text-zinc-500">
      <FileText size={32} className="opacity-20" />
      <p>Security Administrator or Platform Owner access required</p>
    </div>
  );

  const totalPages = Math.ceil(total / limit);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-zinc-800 border border-zinc-700 rounded-2xl flex items-center justify-center">
            <FileText size={22} className="text-zinc-400" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Security & Audit Logs</h1>
            <p className="text-xs text-zinc-400">{total.toLocaleString()} total records • Immutable Append-Only Audit Trail</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <button onClick={() => load(page)} className="flex items-center gap-2 px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-sm font-bold transition-colors">
            <RefreshCw size={14} /> Refresh
          </button>
        </div>
      </div>

      {error && (
        <div className="flex items-center gap-2 bg-red-900/20 border border-red-500/30 text-red-400 px-4 py-3 rounded-xl text-sm">
          <AlertCircle size={14} /> {error}
        </div>
      )}

      <div className="bg-zinc-900 border border-zinc-800 rounded-2xl overflow-hidden">
        {loading ? (
          <div className="p-8 flex justify-center"><Spinner /></div>
        ) : logs.length === 0 ? (
          <div className="p-8 text-center text-zinc-500">
            <FileText size={28} className="mx-auto mb-2 opacity-20" />
            <p>No activity logs found</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-zinc-800 text-zinc-500 text-xs uppercase font-bold">
                <th className="px-5 py-3 text-left">Admin</th>
                <th className="px-5 py-3 text-left">Action</th>
                <th className="px-5 py-3 text-left">Target</th>
                <th className="px-5 py-3 text-left">Notes</th>
                <th className="px-5 py-3 text-right">Time</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/50">
              {logs.map((log: any, i: number) => (
                <tr key={i} className="hover:bg-zinc-800/30">
                  <td className="px-5 py-3 font-mono text-xs text-zinc-400">#{log.admin_id}</td>
                  <td className="px-5 py-3">
                    <span className="font-bold text-white capitalize">{log.action?.replace(/_/g, ' ')}</span>
                  </td>
                  <td className="px-5 py-3 text-zinc-400 text-xs">
                    {log.target_user_id ? `User #${log.target_user_id}` : ''}
                    {log.product_id ? `Listing #${log.product_id}` : ''}
                  </td>
                  <td className="px-5 py-3 text-zinc-500 text-xs max-w-xs truncate">{log.notes ?? '—'}</td>
                  <td className="px-5 py-3 text-right text-zinc-500 text-xs whitespace-nowrap">
                    {log.created_at ? new Date(log.created_at).toLocaleString() : '—'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <p className="text-xs text-zinc-500">Page {page} of {totalPages}</p>
          <div className="flex gap-2">
            <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}
              className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-lg text-xs font-bold disabled:opacity-40">
              ← Prev
            </button>
            <button disabled={page >= totalPages} onClick={() => setPage(p => p + 1)}
              className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-lg text-xs font-bold disabled:opacity-40">
              Next →
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
