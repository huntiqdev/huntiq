'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Scale, Eye, RefreshCw, CheckCircle, XCircle, UserCheck } from 'lucide-react';

const STATUS_TABS = ['pending', 'under_review', 'approved', 'rejected'];

export default function AppealsPage() {
  const { user } = useAuth();
  const [appeals, setAppeals] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [status, setStatus] = useState('pending');
  const [selected, setSelected] = useState<any | null>(null);
  const [decision, setDecision] = useState('');
  const [actionTaken, setActionTaken] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => { fetchAppeals(); }, [status]); // eslint-disable-line react-hooks/exhaustive-deps

  const fetchAppeals = async () => {
    if (!user?.token) return;
    setLoading(true);
    try {
      const data = await api.adminGetAccountAppeals(status, user.token);
      setAppeals(Array.isArray(data) ? data : []);
    } catch { setAppeals([]); } finally { setLoading(false); }
  };

  // Backend expects: { approved: bool, decision: string, action_taken?: string }
  const handleDecision = async (approved: boolean) => {
    if (!selected || !user?.token) return;
    if (!decision.trim()) { setError('Decision notes are required.'); return; }
    setSubmitting(true);
    setError('');
    try {
      await api.adminReviewAppeal(selected.id, {
        approved,
        decision: decision.trim(),
        action_taken: actionTaken.trim() || undefined,
      }, user.token);
      setSelected(null);
      setDecision('');
      setActionTaken('');
      fetchAppeals();
    } catch (err: any) {
      setError(err?.message ?? 'Action failed.');
    } finally { setSubmitting(false); }
  };

  const appealTypeBadge = (type: string) => {
    const colors: Record<string, string> = {
      ban: 'bg-red-900/40 text-red-400',
      deactivation: 'bg-orange-900/40 text-orange-400',
      listing_rejection: 'bg-amber-900/40 text-amber-400',
      verification_rejection: 'bg-blue-900/40 text-blue-400',
      warning: 'bg-zinc-700 text-zinc-300',
    };
    return colors[type] ?? 'bg-zinc-700 text-zinc-300';
  };

  const statusBadge = (s: string) => ({
    pending:      'bg-amber-900/40 text-amber-400',
    under_review: 'bg-blue-900/40 text-blue-400',
    approved:     'bg-emerald-900/40 text-emerald-400',
    rejected:     'bg-red-900/40 text-red-400',
  }[s] ?? 'bg-zinc-700 text-zinc-300');

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-purple-600/10 rounded-xl border border-purple-500/20">
            <Scale className="w-6 h-6 text-purple-400" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Account Appeals</h1>
            <p className="text-xs text-zinc-400">Review ban/suspension appeals and account recovery requests</p>
          </div>
        </div>
        <button onClick={fetchAppeals} className="p-2 bg-zinc-800 rounded-lg hover:bg-zinc-700 transition-colors">
          <RefreshCw className="w-4 h-4 text-zinc-400" />
        </button>
      </div>

      <div className="flex gap-2 flex-wrap">
        {STATUS_TABS.map(tab => (
          <button key={tab} onClick={() => setStatus(tab)}
            className={`px-4 py-2 rounded-xl text-xs font-bold capitalize transition-all ${
              status === tab ? 'bg-orange-600 text-white' : 'bg-zinc-800 text-zinc-400 hover:bg-zinc-700 hover:text-white'
            }`}>
            {tab.replace('_', ' ')}
          </button>
        ))}
      </div>

      <div className="bg-zinc-900 rounded-2xl border border-zinc-800 overflow-hidden">
        {loading ? (
          <div className="p-8 text-center text-zinc-500">
            <div className="animate-spin w-8 h-8 border-2 border-orange-500 border-t-transparent rounded-full mx-auto mb-3" />
            Loading appeals...
          </div>
        ) : appeals.length === 0 ? (
          <div className="p-8 text-center text-zinc-500">
            <Scale className="w-12 h-12 mx-auto mb-3 opacity-20" />
            <p>No {status.replace('_', ' ')} appeals</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-zinc-800 text-zinc-500 text-xs uppercase font-bold">
                <th className="px-6 py-4 text-left">Appellant</th>
                <th className="px-6 py-4 text-left">Appeal Type</th>
                <th className="px-6 py-4 text-left">Submitted</th>
                <th className="px-6 py-4 text-left">Status</th>
                <th className="px-6 py-4 text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/50">
              {appeals.map(a => (
                <tr key={a.id} className="hover:bg-zinc-800/40 transition-colors">
                  <td className="px-6 py-4">
                    <p className="font-semibold text-white">{a.user_name || `User #${a.user_id}`}</p>
                  </td>
                  <td className="px-6 py-4">
                    <span className={`px-2 py-1 rounded-full text-xs font-bold capitalize ${appealTypeBadge(a.appeal_type)}`}>
                      {a.appeal_type?.replace(/_/g, ' ')}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-zinc-400 text-xs">
                    {a.created_at ? new Date(a.created_at).toLocaleDateString() : '—'}
                  </td>
                  <td className="px-6 py-4">
                    <span className={`px-2 py-1 rounded-full text-xs font-bold capitalize ${statusBadge(a.status)}`}>
                      {a.status?.replace('_', ' ')}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button
                      onClick={() => { setSelected(a); setDecision(''); setActionTaken(''); setError(''); }}
                      className="px-3 py-1.5 bg-orange-600/20 text-orange-400 border border-orange-500/20 rounded-lg text-xs font-bold hover:bg-orange-600/30 transition-colors flex items-center gap-1.5 ml-auto"
                    >
                      <Eye className="w-3.5 h-3.5" /> Review
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Review Modal */}
      {selected && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4">
          <div className="bg-zinc-900 border border-zinc-700 rounded-2xl w-full max-w-lg shadow-2xl">
            <div className="px-6 py-5 border-b border-zinc-800 flex items-center justify-between">
              <h2 className="text-lg font-black text-white flex items-center gap-2">
                <Scale className="w-5 h-5 text-purple-400" /> Review Appeal
              </h2>
              <button onClick={() => setSelected(null)} className="text-zinc-500 hover:text-white text-xl leading-none">✕</button>
            </div>
            <div className="p-6 space-y-4">
              <div className="bg-zinc-800 rounded-xl p-4 space-y-3 text-sm">
                <div className="flex justify-between">
                  <span className="text-zinc-500">Appellant</span>
                  <span className="text-white font-medium">{selected.user_name || `User #${selected.user_id}`}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-500">Appeal Type</span>
                  <span className="text-white capitalize">{selected.appeal_type?.replace(/_/g, ' ')}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-500">Current Status</span>
                  <span className={`px-2 py-0.5 rounded-full text-xs font-bold capitalize ${statusBadge(selected.status)}`}>
                    {selected.status?.replace('_', ' ')}
                  </span>
                </div>
              </div>

              {selected.reason && (
                <div className="bg-zinc-800 rounded-xl p-4">
                  <p className="text-xs text-zinc-500 mb-2">Appellant's Reason</p>
                  <p className="text-zinc-300 text-sm">{selected.reason}</p>
                </div>
              )}

              <div>
                <label className="block text-xs font-bold text-zinc-400 uppercase mb-2">Decision Notes *</label>
                <textarea
                  value={decision}
                  onChange={e => setDecision(e.target.value)}
                  rows={3}
                  placeholder="Explain your decision..."
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-3 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-zinc-400 uppercase mb-2">Action Taken (optional)</label>
                <input
                  value={actionTaken}
                  onChange={e => setActionTaken(e.target.value)}
                  placeholder="e.g. Account reactivated, listing approved..."
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-3 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                />
              </div>

              {error && <p className="text-red-400 text-sm bg-red-900/20 border border-red-500/30 rounded-xl px-3 py-2">{error}</p>}

              <div className="flex gap-3 pt-2">
                <button
                  onClick={() => handleDecision(true)}
                  disabled={submitting}
                  className="flex-1 py-3 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-sm font-bold disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
                >
                  <UserCheck className="w-4 h-4" /> {submitting ? 'Saving…' : 'Grant Appeal'}
                </button>
                <button
                  onClick={() => handleDecision(false)}
                  disabled={submitting}
                  className="flex-1 py-3 bg-red-600/20 border border-red-500/30 hover:bg-red-600/30 text-red-400 rounded-xl text-sm font-bold disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
                >
                  <XCircle className="w-4 h-4" /> Deny Appeal
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
