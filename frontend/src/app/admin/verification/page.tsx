'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { BadgeCheck, XCircle, Clock, Eye, MessageSquare, RefreshCw } from 'lucide-react';

const STATUS_TABS = ['pending', 'approved', 'rejected'];

export default function VerificationPage() {
  const { user } = useAuth();
  const [docs, setDocs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [status, setStatus] = useState('pending');
  const [selected, setSelected] = useState<any | null>(null);
  const [notes, setNotes] = useState('');
  const [reason, setReason] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => { fetchDocs(); }, [status]); // eslint-disable-line react-hooks/exhaustive-deps

  const fetchDocs = async () => {
    if (!user?.token) return;
    setLoading(true);
    try {
      const data = await api.adminGetVerificationDocuments(status, user.token);
      setDocs(Array.isArray(data) ? data : []);
    } catch { setDocs([]); } finally { setLoading(false); }
  };

  // Backend expects: { approved: bool, reason: string, notes: string }
  const handleReview = async (approved: boolean) => {
    if (!selected || !user?.token) return;
    if (!approved && !reason.trim()) { setError('Rejection reason is required.'); return; }
    setSubmitting(true);
    setError('');
    try {
      await api.adminReviewVerificationDocument(selected.id, {
        approved,
        reason: reason.trim(),
        notes: notes.trim(),
      }, user.token);
      setSelected(null);
      setNotes('');
      setReason('');
      fetchDocs();
    } catch (err: any) {
      setError(err?.message ?? 'Action failed. Please try again.');
    } finally { setSubmitting(false); }
  };

  const statusBadge = (s: string) => ({
    pending:  'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400',
    approved: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-400',
    rejected: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
  }[s] ?? 'bg-zinc-700 text-zinc-300');

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-emerald-600/10 rounded-xl border border-emerald-500/20">
            <BadgeCheck className="w-6 h-6 text-emerald-400" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Seller Verification</h1>
            <p className="text-xs text-zinc-400">Review submitted identity and business documents</p>
          </div>
        </div>
        <button onClick={fetchDocs} className="p-2 bg-zinc-800 rounded-lg hover:bg-zinc-700 transition-colors">
          <RefreshCw className="w-4 h-4 text-zinc-400" />
        </button>
      </div>

      <div className="flex gap-2 flex-wrap">
        {STATUS_TABS.map(tab => (
          <button key={tab} onClick={() => setStatus(tab)}
            className={`px-4 py-2 rounded-xl text-xs font-bold capitalize transition-all ${
              status === tab ? 'bg-orange-600 text-white' : 'bg-zinc-800 text-zinc-400 hover:bg-zinc-700 hover:text-white'
            }`}>
            {tab}
          </button>
        ))}
      </div>

      <div className="bg-zinc-900 rounded-2xl border border-zinc-800 overflow-hidden">
        {loading ? (
          <div className="p-8 text-center text-zinc-500">
            <div className="animate-spin w-8 h-8 border-2 border-orange-500 border-t-transparent rounded-full mx-auto mb-3" />
            Loading documents...
          </div>
        ) : docs.length === 0 ? (
          <div className="p-8 text-center text-zinc-500">
            <BadgeCheck className="w-12 h-12 mx-auto mb-3 opacity-20" />
            <p>No {status} verification documents</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-zinc-800 text-zinc-500 text-xs uppercase font-bold">
                <th className="px-6 py-4 text-left">Seller</th>
                <th className="px-6 py-4 text-left">Document Type</th>
                <th className="px-6 py-4 text-left">Submitted</th>
                <th className="px-6 py-4 text-left">Status</th>
                <th className="px-6 py-4 text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/50">
              {docs.map(doc => (
                <tr key={doc.id} className="hover:bg-zinc-800/40 transition-colors">
                  <td className="px-6 py-4">
                    <p className="font-semibold text-white">{doc.seller_name || `Seller #${doc.seller_id}`}</p>
                  </td>
                  <td className="px-6 py-4 text-zinc-300 capitalize">
                    {doc.document_type?.replace(/_/g, ' ')}
                  </td>
                  <td className="px-6 py-4 text-zinc-400 text-xs">
                    {doc.created_at ? new Date(doc.created_at).toLocaleDateString() : '—'}
                  </td>
                  <td className="px-6 py-4">
                    <span className={`px-2 py-1 rounded-full text-xs font-bold capitalize ${statusBadge(doc.status)}`}>
                      {doc.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button
                      onClick={() => { setSelected(doc); setNotes(''); setReason(''); setError(''); }}
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
              <h2 className="text-lg font-black text-white">Review Verification Document</h2>
              <button onClick={() => setSelected(null)} className="text-zinc-500 hover:text-white text-xl leading-none">✕</button>
            </div>
            <div className="p-6 space-y-4">
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div className="bg-zinc-800 rounded-xl p-4">
                  <p className="text-xs text-zinc-500 mb-1">Seller</p>
                  <p className="font-bold text-white">{selected.seller_name || `Seller #${selected.seller_id}`}</p>
                </div>
                <div className="bg-zinc-800 rounded-xl p-4">
                  <p className="text-xs text-zinc-500 mb-1">Document Type</p>
                  <p className="font-bold text-white capitalize">{selected.document_type?.replace(/_/g, ' ')}</p>
                  <p className="text-zinc-400 text-xs mt-0.5">
                    Submitted {selected.created_at ? new Date(selected.created_at).toLocaleDateString() : '—'}
                  </p>
                </div>
              </div>

              {selected.file_url && (
                <a href={selected.file_url} target="_blank" rel="noopener noreferrer"
                  className="flex items-center gap-2 px-4 py-3 bg-blue-600/20 border border-blue-500/30 rounded-xl text-blue-400 text-sm hover:bg-blue-600/30 transition-colors">
                  <Eye className="w-4 h-4" /> View Document
                </a>
              )}

              <div>
                <label className="block text-xs font-bold text-zinc-400 uppercase mb-2">
                  <MessageSquare className="w-3.5 h-3.5 inline mr-1" /> Notes (internal)
                </label>
                <textarea value={notes} onChange={e => setNotes(e.target.value)} rows={2}
                  placeholder="Internal notes..."
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-3 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50" />
              </div>

              <div>
                <label className="block text-xs font-bold text-zinc-400 uppercase mb-2">
                  Rejection Reason <span className="normal-case font-normal text-zinc-600">(required if rejecting)</span>
                </label>
                <input value={reason} onChange={e => setReason(e.target.value)}
                  placeholder="e.g. Document unreadable, expired ID..."
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-3 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50" />
              </div>

              {error && <p className="text-red-400 text-sm bg-red-900/20 border border-red-500/30 rounded-xl px-3 py-2">{error}</p>}

              <div className="flex gap-3 pt-2">
                <button onClick={() => handleReview(true)} disabled={submitting}
                  className="flex-1 py-3 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-sm font-bold disabled:opacity-50 transition-colors flex items-center justify-center gap-2">
                  <BadgeCheck className="w-4 h-4" /> {submitting ? 'Saving…' : 'Approve'}
                </button>
                <button onClick={() => handleReview(false)} disabled={submitting}
                  className="flex-1 py-3 bg-red-600/20 border border-red-500/30 hover:bg-red-600/30 text-red-400 rounded-xl text-sm font-bold disabled:opacity-50 transition-colors flex items-center justify-center gap-2">
                  <XCircle className="w-4 h-4" /> Reject
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
