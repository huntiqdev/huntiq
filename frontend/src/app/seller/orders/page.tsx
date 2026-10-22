'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { Phone, Package, MessageSquare, Check, X, AlertCircle, RefreshCw, CheckCircle } from 'lucide-react';
import { ORDER_STATUS } from '@/lib/constants/statuses';

export default function SellerOrdersPage() {
  const { user } = useAuth();
  const [orders, setOrders] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('all');
  const [acting, setActing] = useState<number | null>(null);
  const [error, setError] = useState('');
  const [note, setNote] = useState<Record<number, string>>({});

  const load = () => {
    if (!user) return;
    setLoading(true);
    api.getSellerOrders(user.token)
      .then(r => setOrders(r ?? []))
      .catch(() => {})
      .finally(() => setLoading(false));
  };

  useEffect(() => { load(); }, [user?.token]); // eslint-disable-line

  const act = async (id: number, status: 'accepted' | 'declined' | 'completed') => {
    if (!user) return;
    setActing(id); setError('');
    try {
      await api.updateSellerOrder(id, { status, seller_note: note[id] ?? '' }, user.token);
      setOrders(prev => prev.map(o => o.id === id ? { ...o, status, seller_note: note[id] ?? o.seller_note } : o));
    } catch (err: any) {
      setError(err.message || 'Failed to update order');
      setTimeout(() => setError(''), 4000);
    } finally { setActing(null); }
  };

  const filtered = filter === 'all' ? orders : orders.filter(o => o.status === filter);
  const counts = orders.reduce((acc: Record<string, number>, o) => {
    acc[o.status] = (acc[o.status] ?? 0) + 1;
    return acc;
  }, {});

  if (loading) return <div className="flex items-center justify-center h-64"><Spinner /></div>;

  return (
    <div className="space-y-5 max-w-3xl">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-black text-zinc-900 dark:text-white">Order Requests</h1>
          <p className="text-xs text-zinc-500 mt-0.5">
            {orders.length} total · {counts['pending'] ?? 0} pending response
          </p>
        </div>
        <button onClick={load} className="p-2 rounded-lg text-zinc-400 hover:text-zinc-600 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors" title="Refresh">
          <RefreshCw size={15} />
        </button>
      </div>

      {error && (
        <div className="flex items-center gap-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl px-3 py-2">
          <AlertCircle size={13} className="text-red-500 shrink-0" />
          <p className="text-xs text-red-600 dark:text-red-400">{error}</p>
        </div>
      )}

      {/* Filter tabs */}
      {orders.length > 0 && (
        <div className="flex gap-2 overflow-x-auto pb-1 no-scrollbar">
          {['all', 'pending', 'accepted', 'completed', 'declined'].map(s => (
            <button key={s} onClick={() => setFilter(s)}
              className={`shrink-0 px-3 py-1.5 rounded-xl text-xs font-bold transition-all ${
                filter === s
                  ? 'bg-orange-600 text-white'
                  : 'bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 text-zinc-600 dark:text-zinc-400 hover:border-orange-400'
              }`}>
              {s === 'all' ? `All (${orders.length})` : `${s.charAt(0).toUpperCase() + s.slice(1)}${counts[s] ? ` (${counts[s]})` : ''}`}
            </button>
          ))}
        </div>
      )}

      {filtered.length === 0 ? (
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-12 text-center">
          <Package size={32} className="mx-auto mb-3 text-zinc-300 dark:text-zinc-600" />
          <p className="font-bold text-zinc-500">{filter === 'all' ? 'No order requests yet' : `No ${filter} orders`}</p>
          <p className="text-sm text-zinc-400 mt-1">Buyers will contact you when they want to purchase your products</p>
        </div>
      ) : (
        <div className="space-y-3">
          {filtered.map(o => {
            const sc = ORDER_STATUS[o.status] ?? ORDER_STATUS.pending;
            const StatusIcon = sc.icon;
            const isPending = o.status === 'pending';
            const isAccepted = o.status === 'accepted';

            return (
              <div key={o.id} className={`bg-white dark:bg-zinc-900 border rounded-2xl p-4 space-y-3 ${
                isPending ? 'border-blue-200 dark:border-blue-800' : 'border-zinc-200 dark:border-zinc-800'
              }`}>
                {/* Header */}
                <div className="flex items-start justify-between gap-3">
                  <div className="flex items-start gap-3 flex-1 min-w-0">
                    <div className="w-9 h-9 rounded-xl bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center shrink-0">
                      <Package size={15} className="text-orange-600" />
                    </div>
                    <div className="min-w-0">
                      <p className="font-bold text-sm text-zinc-800 dark:text-zinc-200 truncate">{o.product_title ?? `Product #${o.product_id}`}</p>
                      {o.product_price > 0 && (
                        <p className="text-xs font-black text-orange-600">ETB {Number(o.product_price).toLocaleString()}</p>
                      )}
                    </div>
                  </div>
                  <span className={`flex items-center gap-1 text-[10px] font-black px-2 py-0.5 rounded-full uppercase shrink-0 ${sc.color}`}>
                    <StatusIcon size={9} /> {sc.label}
                  </span>
                </div>

                {/* Buyer info */}
                <div className="flex items-center gap-3 bg-zinc-50 dark:bg-zinc-800 rounded-xl px-3 py-2.5">
                  <Phone size={13} className="text-zinc-400 shrink-0" />
                  <div>
                    <p className="text-sm font-bold text-zinc-800 dark:text-zinc-200">{o.buyer_name}</p>
                    <p className="text-xs text-zinc-500">{o.buyer_phone}</p>
                  </div>
                  {isPending && (
                    <a href={`tel:${o.buyer_phone}`}
                      className="ml-auto flex items-center gap-1 text-xs font-bold text-orange-600 hover:underline shrink-0">
                      <Phone size={11} /> Call
                    </a>
                  )}
                </div>

                {/* Buyer message */}
                {o.message && (
                  <div className="flex items-start gap-2 bg-blue-50 dark:bg-blue-900/20 border border-blue-100 dark:border-blue-800 rounded-xl px-3 py-2">
                    <MessageSquare size={11} className="text-blue-400 shrink-0 mt-0.5" />
                    <p className="text-xs text-blue-700 dark:text-blue-300 italic">"{o.message}"</p>
                  </div>
                )}

                {/* Seller note (existing) */}
                {o.seller_note && !isPending && (
                  <p className="text-xs text-zinc-400 italic pl-1">Your note: "{o.seller_note}"</p>
                )}

                {/* Actions for pending orders */}
                {isPending && (
                  <div className="space-y-2 pt-1">
                    <input
                      type="text"
                      placeholder="Optional note to buyer (e.g. 'Call me after 5pm')"
                      value={note[o.id] ?? ''}
                      onChange={e => setNote(prev => ({ ...prev, [o.id]: e.target.value }))}
                      className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-xs outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 transition-all"
                    />
                    <div className="flex gap-2">
                      <button onClick={() => act(o.id, 'accepted')} disabled={acting === o.id}
                        className="flex-1 flex items-center justify-center gap-1.5 py-2 rounded-xl bg-green-500 hover:bg-green-400 text-white text-xs font-bold transition-colors disabled:opacity-50">
                        {acting === o.id ? <Spinner /> : <Check size={13} />} Accept
                      </button>
                      <button onClick={() => act(o.id, 'declined')} disabled={acting === o.id}
                        className="flex-1 flex items-center justify-center gap-1.5 py-2 rounded-xl bg-red-100 dark:bg-red-900/30 hover:bg-red-200 text-red-600 text-xs font-bold transition-colors disabled:opacity-50">
                        {acting === o.id ? <Spinner /> : <X size={13} />} Decline
                      </button>
                    </div>
                  </div>
                )}

                {/* Mark completed for accepted orders */}
                {isAccepted && (
                  <button onClick={() => act(o.id, 'completed')} disabled={acting === o.id}
                    className="w-full flex items-center justify-center gap-1.5 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-white text-xs font-bold transition-colors disabled:opacity-50">
                    {acting === o.id ? <Spinner /> : <CheckCircle size={13} />} Mark as Completed
                  </button>
                )}

                <p className="text-[10px] text-zinc-300 dark:text-zinc-600">
                  {new Date(o.created_at).toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' })}
                </p>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
