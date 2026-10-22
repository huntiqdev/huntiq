'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { RotateCcw, Eye, RefreshCw, CheckCircle, XCircle, DollarSign } from 'lucide-react';

const STATUS_TABS = ['pending', 'approved', 'rejected', 'processed'];

export default function RefundsPage() {
  const { user } = useAuth();
  const [refunds, setRefunds] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [status, setStatus] = useState('pending');
  const [selected, setSelected] = useState<any | null>(null);
  const [notes, setNotes] = useState('');
  const [reference, setReference] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => { fetchRefunds(); }, [status]); // eslint-disable-line react-hooks/exhaustive-deps

  const fetchRefunds = async () => {
    if (!user?.token) return;
    setLoading(true);
    try {
      const data = await api.adminGetRefundRequests(status, user.token);
      setRefunds(Array.isArray(data) ? data : []);
    } catch { setRefunds([]); } finally { setLoading(false); }
  };

  // Approve or reject — backend expects { approved: bool, reason: string }
  const handleReview = async (approved: boolean) => {
    if (!selected || !user?.token) return;
    setSubmitting(true);
    setError('');
    try {
      await api.adminReviewRefund(selected.id, { approved, reason: notes }, user.token);
      setSelected(null);
      setNotes('');
      fetchRefunds();
    } catch (err: any) {
      setError(err?.message ?? 'Action failed.');
    } finally { setSubmitting(false); }
  };

  // Mark as processed — backend requires { reference: string }
  const handleMarkProcessed = async () => {
    if (!selected || !user?.token) return;
    if (!reference.trim()) { setError('Payment reference is required to mark as processed.'); return; }
    setSubmitting(true);
    setError('');
    try {
      await api.adminMarkRefundProcessed(selected.id, reference.trim(), user.token);
      setSelected(null);
      setReference('');
      fetchRefunds();
    } catch (err: any) {
      setError(err?.message ?? 'Action failed.');
    } finally { setSubmitting(false); }
  };

  const statusColor = (s: string) => ({
    pending:  'bg-amber-900/40 text-amber-400',
    approved: 'bg-emerald-900/40 text-emerald-400',
    rejected: 'bg-red-900/40 text-red-400',
    processed:'bg-zinc-700 text-zinc-300',
  }[s] ?? 'bg-zinc-700 text-zinc-300');

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-blue-600/10 rounded-xl border border-blue-500/20">
            <RotateCcw className="w-6 h-6 text-blue-400" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Refund Management</h1>
            <p className="text-xs text-zinc-400">Review and process transaction refund requests</p>
          </div>
        </div>
        <button onClick={fetchRefunds} className="p-2 bg-zinc-800 rounded-lg hover:bg-zinc-700 transition-colors">
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
            Loading refund requests...
          </div>
        ) : refunds.length === 0 ? (
          <div className="p-8 text-center text-zinc-500">
            <RotateCcw className="w-12 h-12 mx-auto mb-3 opacity-20" />
            <p>No {status} refund requests</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-zinc-800 text-zinc-500 text-xs uppercase font-bold">
                <th className="px-6 py-4 text-left">Requester</th>
                <th className="px-6 py-4 text-left">Amount</th>
                <th className="px-6 py-4 text-left">Reason</th>
                <th className="px-6 py-4 text-left">Submitted</th>
                <th className="px-6 py-4 text-left">Status</th>
                <th className="px-6 py-4 text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/50">
              {refunds.map(r => (
                <tr key={r.id} className="hover:bg-zinc-800/40 transition-colors">
                  <td className="px-6 py-4">
                    <p className="font-semibold text-white">{r.requester_name || `User #${r.requester_id}`}</p>
                    <p className="text-xs text-zinc-500 capitalize">{r.requester_role}</p>
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-1">
                      <DollarSign className="w-3.5 h-3.5 text-emerald-400" />
                      <span className="font-mono font-bold text-white">ETB {r.amount?.toLocaleString()}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-zinc-400 text-xs max-w-xs truncate capitalize">
                    {r.reason?.replace(/_/g, ' ')}
                  </td>
                  <td className="px-6 py-4 text-zinc-400 text-xs">
                    {r.created_at ? new Date(r.created_at).toLocaleDateString() : '—'}
                  </td>
                  <td className="px-6 py-4">
                    <span className={`px-2 py-1 rounded-full text-xs font-bold capitalize ${statusColor(r.status)}`}>
                      {r.status?.replace('_', ' ')}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button
                      onClick={() => { setSelected(r); setNotes(''); setReference(''); setError(''); }}
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
              <h2 className="text-lg font-black text-white">Review Refund Request</h2>
              <button onClick={() => setSelected(null)} className="text-zinc-500 hover:text-white text-xl leading-none">✕</button>
            </div>
            <div className="p-6 space-y-4">
              <div className="bg-zinc-800 rounded-xl p-4 space-y-3 text-sm">
                <div className="flex justify-between">
                  <span className="text-zinc-500">Requester</span>
                  <span className="text-white font-medium">{selected.requester_name || `User #${selected.requester_id}`}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-500">Amount</span>
                  <span className="text-emerald-400 font-bold font-mono">ETB {selected.amount?.toLocaleString()}</span>
                </div>
                {selected.transaction_id && (
                  <div className="flex justify-between">
                    <span className="text-zinc-500">Transaction</span>
                    <span className="text-white">#{selected.transaction_id}</span>
                  </div>
                )}
                <div className="flex justify-between">
                  <span className="text-zinc-500">Reason</span>
                  <span className="text-white capitalize">{selected.reason?.replace(/_/g, ' ')}</span>
                </div>
              </div>

              {selected.description && (
                <div className="bg-zinc-800 rounded-xl p-4">
                  <p className="text-xs text-zinc-500 mb-2">Details</p>
                  <p className="text-zinc-300 text-sm">{selected.description}</p>
                </div>
              )}

              {/* Show reference input only when marking as processed */}
              {selected.status === 'approved' ? (
                <div>
                  <label className="block text-xs font-bold text-zinc-400 uppercase mb-2">Payment Reference *</label>
                  <input
                    value={reference}
                    onChange={e => setReference(e.target.value)}
                    placeholder="e.g. TXN-12345 or Telebirr reference"
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-3 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                  />
                  <p className="text-xs text-zinc-500 mt-1">Required to confirm the refund was paid</p>
                </div>
              ) : (
                <div>
                  <label className="block text-xs font-bold text-zinc-400 uppercase mb-2">Reason / Notes</label>
                  <textarea
                    value={notes}
                    onChange={e => setNotes(e.target.value)}
                    rows={3}
                    placeholder="Document your decision..."
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-3 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                  />
                </div>
              )}

              {error && <p className="text-red-400 text-sm bg-red-900/20 border border-red-500/30 rounded-xl px-3 py-2">{error}</p>}

              <div className="flex gap-3 pt-2">
                {selected.status === 'approved' ? (
                  <button
                    onClick={handleMarkProcessed}
                    disabled={submitting}
                    className="flex-1 py-3 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-sm font-bold disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
                  >
                    <CheckCircle className="w-4 h-4" /> {submitting ? 'Processing…' : 'Mark Processed'}
                  </button>
                ) : (
                  <>
                    <button
                      onClick={() => handleReview(true)}
                      disabled={submitting}
                      className="flex-1 py-3 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-sm font-bold disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
                    >
                      <CheckCircle className="w-4 h-4" /> {submitting ? 'Saving…' : 'Approve'}
                    </button>
                    <button
                      onClick={() => handleReview(false)}
                      disabled={submitting}
                      className="flex-1 py-3 bg-red-600/20 border border-red-500/30 hover:bg-red-600/30 text-red-400 rounded-xl text-sm font-bold disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
                    >
                      <XCircle className="w-4 h-4" /> Reject
                    </button>
                  </>
                )}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
