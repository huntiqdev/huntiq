'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Scale, Eye, RefreshCw, CheckCircle, XCircle, UserCheck, ShieldAlert, Clock, CheckCircle2, TrendingUp, AlertTriangle } from 'lucide-react';

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

  const fetchAppeals = async () => {
    if (!user?.token) return;
    setLoading(true);
    try {
      const data = await api.adminGetAccountAppeals(status, user.token);
      setAppeals(Array.isArray(data) ? data : []);
    } catch { setAppeals([]); } finally { setLoading(false); }
  };

  useEffect(() => { fetchAppeals(); }, [status]); // eslint-disable-line react-hooks/exhaustive-deps

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
      ban: 'bg-rose-500/10 text-rose-400 border border-rose-500/20',
      deactivation: 'bg-orange-500/10 text-orange-400 border border-orange-500/20',
      listing_rejection: 'bg-amber-500/10 text-amber-400 border border-amber-500/20',
      verification_rejection: 'bg-blue-500/10 text-blue-400 border border-blue-500/20',
      warning: 'bg-zinc-800 text-zinc-300 border border-zinc-700',
    };
    return colors[type] ?? 'bg-zinc-800 text-zinc-300 border border-zinc-700';
  };

  const statusBadge = (s: string) => ({
    pending:      'bg-amber-500/10 text-amber-400 border border-amber-500/20',
    under_review: 'bg-blue-500/10 text-blue-400 border border-blue-500/20',
    approved:     'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20',
    rejected:     'bg-rose-500/10 text-rose-400 border border-rose-500/20',
  }[s] ?? 'bg-zinc-800 text-zinc-300 border border-zinc-700');

  return (
    <div className="space-y-6 max-w-[1600px] mx-auto pb-10">

      {/* ── BREADCRUMB & HEADER (Image 4 Style) ── */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-2 border-b border-zinc-800/80">
        <div>
          <div className="flex items-center gap-2 text-xs font-semibold text-zinc-400">
            <span>Pages</span>
            <span>/</span>
            <span>Trust & Compliance</span>
            <span>/</span>
            <span className="text-zinc-200">Account Appeals</span>
          </div>
          <h1 className="text-2xl font-black text-white tracking-tight mt-1 flex items-center gap-3">
            Account Appeals
            <span className="text-[10px] font-extrabold uppercase px-2.5 py-1 bg-purple-500/10 text-purple-400 border border-purple-500/20 rounded-full">
              Legal & Safety Desk
            </span>
          </h1>
        </div>

        <button
          onClick={fetchAppeals}
          className="flex items-center gap-2 px-3.5 py-2 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 text-zinc-300 rounded-xl text-xs font-bold transition-all w-fit"
        >
          <RefreshCw size={13} className={loading ? 'animate-spin text-orange-500' : 'text-zinc-400'} />
          <span>Refresh Queue</span>
        </button>
      </div>

      {/* ── 4-KPI RIBBON (Image 4 Style) ── */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3.5">
        <div className="bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-4 flex items-center gap-3.5 shadow-sm">
          <div className="w-10 h-10 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400 shrink-0">
            <Scale size={18} />
          </div>
          <div>
            <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-wider">Queue Total</p>
            <p className="text-lg font-black text-white font-mono mt-0.5">{appeals.length} Cases</p>
          </div>
        </div>

        <div className="bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-4 flex items-center gap-3.5 shadow-sm">
          <div className="w-10 h-10 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-400 shrink-0">
            <Clock size={18} />
          </div>
          <div>
            <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-wider">Active Status</p>
            <p className="text-lg font-black text-amber-400 font-mono mt-0.5 capitalize">{status.replace('_', ' ')}</p>
          </div>
        </div>

        <div className="bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-4 flex items-center gap-3.5 shadow-sm">
          <div className="w-10 h-10 rounded-xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400 shrink-0">
            <CheckCircle2 size={18} />
          </div>
          <div>
            <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-wider">Resolution SLA</p>
            <p className="text-lg font-black text-emerald-400 font-mono mt-0.5">99.4%</p>
          </div>
        </div>

        <div className="bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-4 flex items-center gap-3.5 shadow-sm">
          <div className="w-10 h-10 rounded-xl bg-blue-500/10 border border-blue-500/20 flex items-center justify-center text-blue-400 shrink-0">
            <TrendingUp size={18} />
          </div>
          <div>
            <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-wider">Avg Decision Time</p>
            <p className="text-lg font-black text-white font-mono mt-0.5">&lt; 3.5 hrs</p>
          </div>
        </div>
      </div>

      {/* ── FILTER TABS ── */}
      <div className="flex gap-2 flex-wrap bg-zinc-900/60 p-1.5 rounded-2xl border border-zinc-800/80 w-fit">
        {STATUS_TABS.map(tab => (
          <button
            key={tab}
            onClick={() => setStatus(tab)}
            className={`px-4 py-2 rounded-xl text-xs font-bold capitalize transition-all ${
              status === tab
                ? 'bg-zinc-800 text-white shadow-sm border border-zinc-700'
                : 'text-zinc-400 hover:text-zinc-200'
            }`}
          >
            {tab.replace('_', ' ')}
          </button>
        ))}
      </div>

      {/* ── APPEALS DATA TABLE ── */}
      <div className="bg-zinc-900/90 rounded-2xl border border-zinc-800/80 overflow-hidden shadow-sm">
        {loading ? (
          <div className="p-12 text-center text-zinc-400">
            <div className="animate-spin w-8 h-8 border-2 border-orange-500 border-t-transparent rounded-full mx-auto mb-3" />
            <p className="text-xs font-bold">Synchronizing appeals audit trail...</p>
          </div>
        ) : appeals.length === 0 ? (
          <div className="p-12 text-center text-zinc-400">
            <Scale className="w-12 h-12 mx-auto mb-3 text-zinc-600" />
            <p className="text-sm font-bold text-zinc-300">No {status.replace('_', ' ')} appeals</p>
            <p className="text-xs text-zinc-500 mt-1">All cases in this category are fully resolved.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead>
                <tr className="border-b border-zinc-800 text-zinc-400 text-[10px] uppercase font-bold tracking-wider bg-zinc-950/40">
                  <th className="px-6 py-3.5">Appellant Identifier</th>
                  <th className="px-6 py-3.5">Incident Type</th>
                  <th className="px-6 py-3.5">Timestamp</th>
                  <th className="px-6 py-3.5">Review Status</th>
                  <th className="px-6 py-3.5 text-right">Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-800/60">
                {appeals.map(a => (
                  <tr key={a.id} className="hover:bg-zinc-800/40 transition-colors">
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-full bg-zinc-800 border border-zinc-700 flex items-center justify-center font-bold text-zinc-300 text-xs">
                          {(a.user_name || 'U').charAt(0).toUpperCase()}
                        </div>
                        <div>
                          <p className="font-bold text-white">{a.user_name || `User #${a.user_id}`}</p>
                          <p className="text-[10px] text-zinc-500 font-mono">UID: {a.user_id}</p>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`px-2.5 py-1 rounded-lg text-[10px] font-bold capitalize ${appealTypeBadge(a.appeal_type)}`}>
                        {a.appeal_type?.replace(/_/g, ' ')}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-zinc-400 font-mono text-[11px]">
                      {a.created_at ? new Date(a.created_at).toLocaleString() : '—'}
                    </td>
                    <td className="px-6 py-4">
                      <span className={`px-2.5 py-1 rounded-lg text-[10px] font-bold capitalize ${statusBadge(a.status)}`}>
                        {a.status?.replace('_', ' ')}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <button
                        onClick={() => { setSelected(a); setDecision(''); setActionTaken(''); setError(''); }}
                        className="px-3 py-1.5 bg-orange-600/15 text-orange-400 border border-orange-500/20 rounded-xl text-xs font-bold hover:bg-orange-600/25 transition-all inline-flex items-center gap-1.5"
                      >
                        <Eye className="w-3.5 h-3.5" /> Review Case
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* ── REVIEW MODAL ── */}
      {selected && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
          <div className="bg-zinc-900 border border-zinc-700/80 rounded-2xl w-full max-w-lg shadow-2xl overflow-hidden">
            <div className="px-6 py-5 border-b border-zinc-800 flex items-center justify-between bg-zinc-950/50">
              <h2 className="text-base font-black text-white flex items-center gap-2">
                <Scale className="w-5 h-5 text-purple-400" /> Appeal Docket #{selected.id}
              </h2>
              <button onClick={() => setSelected(null)} className="text-zinc-500 hover:text-white text-lg font-bold">✕</button>
            </div>
            <div className="p-6 space-y-4">
              <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-4 space-y-2.5 text-xs">
                <div className="flex justify-between">
                  <span className="text-zinc-500">Appellant Account</span>
                  <span className="text-white font-bold">{selected.user_name || `User #${selected.user_id}`}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-500">Violation / Appeal Category</span>
                  <span className="text-white capitalize font-semibold">{selected.appeal_type?.replace(/_/g, ' ')}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-500">Current Queue State</span>
                  <span className={`px-2 py-0.5 rounded-md text-[10px] font-bold capitalize ${statusBadge(selected.status)}`}>
                    {selected.status?.replace('_', ' ')}
                  </span>
                </div>
              </div>

              {selected.reason && (
                <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-4">
                  <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-wider mb-1.5">Appellant Statement</p>
                  <p className="text-zinc-300 text-xs leading-relaxed">{selected.reason}</p>
                </div>
              )}

              <div>
                <label className="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-2">Decision Notes *</label>
                <textarea
                  value={decision}
                  onChange={e => setDecision(e.target.value)}
                  rows={3}
                  placeholder="Record formal moderation findings and justification..."
                  className="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-4 py-3 text-xs text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-2">Corrective Action Taken</label>
                <input
                  value={actionTaken}
                  onChange={e => setActionTaken(e.target.value)}
                  placeholder="e.g. Account reinstated, warning issued, 30-day probation..."
                  className="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-4 py-2.5 text-xs text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500"
                />
              </div>

              {error && <p className="text-rose-400 text-xs bg-rose-500/10 border border-rose-500/20 rounded-xl px-3 py-2">{error}</p>}

              <div className="flex gap-3 pt-2">
                <button
                  onClick={() => handleDecision(true)}
                  disabled={submitting}
                  className="flex-1 py-2.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl text-xs font-bold disabled:opacity-50 transition-all flex items-center justify-center gap-1.5 shadow-sm"
                >
                  <UserCheck className="w-4 h-4" /> {submitting ? 'Saving…' : 'Approve & Restore'}
                </button>
                <button
                  onClick={() => handleDecision(false)}
                  disabled={submitting}
                  className="flex-1 py-2.5 bg-rose-500/15 border border-rose-500/30 hover:bg-rose-500/25 text-rose-400 rounded-xl text-xs font-bold disabled:opacity-50 transition-all flex items-center justify-center gap-1.5"
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
