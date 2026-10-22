'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { AlertTriangle, Eye, RefreshCw, Flag, Shield } from 'lucide-react';

const STATUS_TABS = ['pending', 'investigating', 'resolved', 'dismissed'];
const PRIORITY_LEVELS = ['low', 'medium', 'high', 'critical'] as const;
type Priority = typeof PRIORITY_LEVELS[number];

const PRIORITY_STYLE: Record<Priority, string> = {
  low:      'bg-zinc-700 text-zinc-300',
  medium:   'bg-amber-900/40 text-amber-400 border border-amber-500/20',
  high:     'bg-orange-900/40 text-orange-400 border border-orange-500/20',
  critical: 'bg-red-900/40 text-red-400 border border-red-500/20',
};

// Possible actions for an abuse report
const ACTIONS = [
  { value: 'warn',            label: 'Issue Warning' },
  { value: 'content_removed', label: 'Remove Content' },
  { value: 'suspended',       label: 'Suspend Account' },
  { value: 'banned',          label: 'Ban Account' },
  { value: 'no_action',       label: 'No Action Required' },
];

export default function TrustSafetyPage() {
  const { user } = useAuth();
  const [reports, setReports] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [status, setStatus] = useState('pending');
  const [selected, setSelected] = useState<any | null>(null);
  // Backend ReviewAbuseReport expects: { status, action_taken, notes }
  const [actionTaken, setActionTaken] = useState('');
  const [notes, setNotes] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => { fetchReports(); }, [status]); // eslint-disable-line react-hooks/exhaustive-deps

  const fetchReports = async () => {
    if (!user?.token) return;
    setLoading(true);
    try {
      const data = await api.adminGetAbuseReports(status, user.token);
      setReports(Array.isArray(data) ? data : []);
    } catch { setReports([]); } finally { setLoading(false); }
  };

  // Backend PUT /admin/trust/reports/:id/review expects: { status, action_taken, notes }
  const handleReview = async () => {
    if (!selected || !user?.token || !actionTaken) return;
    setSubmitting(true);
    setError('');
    try {
      await api.adminReviewAbuseReport(selected.id, {
        status: 'resolved',
        action: actionTaken,   // field name the API client sends as action_taken
        notes,
      }, user.token);
      setSelected(null);
      setNotes('');
      setActionTaken('');
      fetchReports();
    } catch (err: any) {
      setError(err?.message ?? 'Action failed.');
    } finally { setSubmitting(false); }
  };

  const setPriority = async (id: number, priority: Priority) => {
    if (!user?.token) return;
    try {
      await api.adminUpdateAbuseReportPriority(id, priority, user.token);
      setReports(prev => prev.map(r => r.id === id ? { ...r, priority } : r));
    } catch { /* silent */ }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-orange-600/10 rounded-xl border border-orange-500/20">
            <AlertTriangle className="w-6 h-6 text-orange-400" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Trust & Safety</h1>
            <p className="text-xs text-zinc-400">Review abuse reports and take enforcement actions</p>
          </div>
        </div>
        <button onClick={fetchReports} className="p-2 bg-zinc-800 rounded-lg hover:bg-zinc-700 transition-colors">
          <RefreshCw className="w-4 h-4 text-zinc-400" />
        </button>
      </div>

      {/* Status tabs */}
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

      {/* Reports table */}
      <div className="bg-zinc-900 rounded-2xl border border-zinc-800 overflow-hidden">
        {loading ? (
          <div className="p-8 text-center text-zinc-500">
            <div className="animate-spin w-8 h-8 border-2 border-orange-500 border-t-transparent rounded-full mx-auto mb-3" />
            Loading reports...
          </div>
        ) : reports.length === 0 ? (
          <div className="p-8 text-center text-zinc-500">
            <AlertTriangle className="w-12 h-12 mx-auto mb-3 opacity-20" />
            <p>No {status} reports</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-zinc-800 text-zinc-500 text-xs uppercase font-bold">
                <th className="px-5 py-4 text-left">Reporter</th>
                <th className="px-5 py-4 text-left">Reported</th>
                <th className="px-5 py-4 text-left">Reason</th>
                <th className="px-5 py-4 text-left">Priority</th>
                <th className="px-5 py-4 text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/50">
              {reports.map(r => (
                <tr key={r.id} className="hover:bg-zinc-800/40 transition-colors">
                  <td className="px-5 py-4">
                    <p className="font-semibold text-white">{r.reporter_name || 'Anonymous'}</p>
                    <p className="text-xs text-zinc-500">{new Date(r.created_at).toLocaleDateString()}</p>
                  </td>
                  <td className="px-5 py-4 text-zinc-300">
                    {r.reported_user_name || (r.reported_user_id ? `User #${r.reported_user_id}` : '—')}
                    {r.reported_product_id && (
                      <p className="text-xs text-zinc-500">Listing #{r.reported_product_id}</p>
                    )}
                  </td>
                  <td className="px-5 py-4 text-zinc-400 capitalize text-xs">{r.reason?.replace(/_/g, ' ')}</td>
                  <td className="px-5 py-4">
                    {/* Priority toggle buttons instead of select */}
                    <div className="flex gap-1">
                      {PRIORITY_LEVELS.map(p => (
                        <button key={p}
                          onClick={() => setPriority(r.id, p)}
                          className={`px-2 py-0.5 rounded text-[9px] font-bold uppercase transition-all ${
                            (r.priority || 'low') === p
                              ? PRIORITY_STYLE[p]
                              : 'bg-zinc-800 text-zinc-600 hover:text-zinc-400'
                          }`}>
                          {p[0].toUpperCase()}
                        </button>
                      ))}
                    </div>
                    <p className="text-[9px] text-zinc-500 mt-1 capitalize">{r.priority || 'low'}</p>
                  </td>
                  <td className="px-5 py-4 text-right">
                    <button
                      onClick={() => { setSelected(r); setNotes(''); setActionTaken(''); setError(''); }}
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
          <div className="bg-zinc-900 border border-zinc-700 rounded-2xl w-full max-w-2xl shadow-2xl">
            <div className="px-6 py-5 border-b border-zinc-800 flex items-center justify-between">
              <h2 className="text-lg font-black text-white flex items-center gap-2">
                <Flag className="w-5 h-5 text-orange-400" /> Review Abuse Report
              </h2>
              <button onClick={() => setSelected(null)} className="text-zinc-500 hover:text-white text-xl leading-none">✕</button>
            </div>
            <div className="p-6 space-y-4">
              {/* Report details */}
              <div className="bg-zinc-800 rounded-xl p-4 space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-zinc-500">Reason</span>
                  <span className="text-white capitalize font-medium">{selected.reason?.replace(/_/g, ' ')}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-500">Reported</span>
                  <span className="text-white">{selected.reported_user_name || `User #${selected.reported_user_id}`}</span>
                </div>
                {selected.description && (
                  <div>
                    <p className="text-zinc-500 mb-1">Description</p>
                    <p className="text-zinc-300">{selected.description}</p>
                  </div>
                )}
              </div>

              {/* Action */}
              <div>
                <label className="block text-xs font-bold text-zinc-400 uppercase mb-2">
                  <Shield className="w-3 h-3 inline mr-1" /> Enforcement Action *
                </label>
                <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                  {ACTIONS.map(a => (
                    <button key={a.value} type="button"
                      onClick={() => setActionTaken(a.value)}
                      className={`px-3 py-2.5 rounded-xl text-xs font-bold text-left transition-all border ${
                        actionTaken === a.value
                          ? 'bg-orange-600/20 text-orange-400 border-orange-500/30'
                          : 'bg-zinc-800 text-zinc-400 border-zinc-700 hover:text-white hover:border-zinc-600'
                      }`}>
                      {a.label}
                    </button>
                  ))}
                </div>
              </div>

              {/* Notes */}
              <div>
                <label className="block text-xs font-bold text-zinc-400 uppercase mb-2">Internal Notes</label>
                <textarea
                  value={notes}
                  onChange={e => setNotes(e.target.value)}
                  rows={3}
                  placeholder="Document your decision..."
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-3 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                />
              </div>

              {error && <p className="text-red-400 text-sm bg-red-900/20 border border-red-500/30 rounded-xl px-3 py-2">{error}</p>}

              <div className="flex gap-3 pt-2">
                <button onClick={() => setSelected(null)}
                  className="flex-1 py-3 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-sm font-bold transition-colors">
                  Cancel
                </button>
                <button
                  onClick={handleReview}
                  disabled={submitting || !actionTaken}
                  className="flex-1 py-3 bg-orange-600 hover:bg-orange-700 text-white rounded-xl text-sm font-bold disabled:opacity-50 transition-colors">
                  {submitting ? 'Submitting…' : 'Submit Decision'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
