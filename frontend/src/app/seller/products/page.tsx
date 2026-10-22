'use client';
import { useEffect, useState, useCallback, Suspense } from 'react';
import { useAuth } from '@/context/AuthContext';
import { useSearchParams, useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { Modal } from '@/components/ui/Modal';
import { ProductForm } from '@/components/seller/ProductForm';
import { Plus, Pencil, Trash2, Search, Image as ImageIcon, AlertCircle, Clock, DollarSign } from 'lucide-react';
import { resolveImageUrl } from '@/lib/utils/imageUrl';

const STATUS_STYLE: Record<string, string> = {
  pending:  'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-400',
  approved: 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400',
  rejected: 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400',
  expired:  'bg-zinc-100 dark:bg-zinc-800 text-zinc-500',
  sold:     'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400',
};

function daysRemaining(expiresAt: string | null | undefined): number | null {
  if (!expiresAt) return null;
  const diff = new Date(expiresAt).getTime() - Date.now();
  return Math.max(0, Math.ceil(diff / 86_400_000));
}

function SellerProductsContent() {
  const { user } = useAuth();
  const searchParams = useSearchParams();
  const router = useRouter();
  const [products, setProducts] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [intentCondition, setIntentCondition] = useState<string | undefined>();
  const [editing, setEditing] = useState<any | null>(null);
  const [filter, setFilter] = useState('all');
  const [search, setSearch] = useState('');
  const [error, setError] = useState('');

  const [confirmDelete, setConfirmDelete] = useState<number | null>(null);

  const load = useCallback(() => {
    if (!user) return;
    api.getSellerProductsDetailed(user.token)
      .then(p => setProducts(p ?? []))
      .finally(() => setLoading(false));
  }, [user]);

  useEffect(() => { load(); }, [load]);

  // Auto-open add product modal when coming from /sell with intent
  useEffect(() => {
    if (searchParams.get('new') === '1') {
      const cond = searchParams.get('condition') ?? undefined;
      setIntentCondition(cond);
      setShowCreate(true);
    }
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const handleDelete = async (id: number) => {
    if (!user) return;
    if (confirmDelete !== id) { setConfirmDelete(id); return; }
    setConfirmDelete(null);
    try {
      await api.deleteSellerProduct(id, user.token);
      setProducts(prev => prev.filter(p => p.id !== id));
    } catch (err: any) { setError(err.message); }
  };

  const FILTERS = ['all', 'pending', 'approved', 'rejected', 'expired'];
  const filtered = products
    .filter(p => filter === 'all' || p.status === filter)
    .filter(p => !search || p.title.toLowerCase().includes(search.toLowerCase()));

  const counts = FILTERS.reduce((acc, s) => {
    acc[s] = s === 'all' ? products.length : products.filter(p => p.status === s).length;
    return acc;
  }, {} as Record<string, number>);

  if (loading) return <div className="flex items-center justify-center h-64"><Spinner /></div>;

  return (
    <div className="space-y-4 max-w-4xl">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-black text-zinc-900 dark:text-white">My Products</h1>
          <p className="text-xs text-zinc-500 mt-0.5">{products.length} total listings</p>
        </div>
        <button onClick={() => setShowCreate(true)}
          className="flex items-center gap-2 bg-orange-600 hover:bg-orange-700 text-white px-4 py-2 rounded-xl text-sm font-bold transition-colors shadow-md shadow-orange-600/20">
          <Plus size={15} /> Add Product
        </button>
      </div>

      {error && (
        <div className="flex items-center gap-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl px-3 py-2">
          <AlertCircle size={14} className="text-red-500 shrink-0" />
          <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
        </div>
      )}

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-3">
        <div className="relative flex-1 max-w-xs">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none" />
          <input value={search} onChange={e => setSearch(e.target.value)}
            placeholder="Search products..."
            className="w-full pl-9 pr-3 py-2 border border-zinc-200 dark:border-zinc-700 rounded-xl text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 transition-all" />
        </div>
        <div className="flex gap-1.5 flex-wrap">
          {FILTERS.map(s => (
            <button key={s} onClick={() => setFilter(s)}
              className={`px-3 py-1.5 rounded-xl text-xs font-bold capitalize transition-all ${
                filter === s ? 'bg-orange-600 text-white' : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 hover:bg-zinc-200 dark:hover:bg-zinc-700'
              }`}>
              {s} ({counts[s]})
            </button>
          ))}
        </div>
      </div>

      {/* Product list */}
      {filtered.length === 0 ? (
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-12 text-center">
          <p className="text-zinc-400">No products found</p>
        </div>
      ) : (
        <div className="space-y-3">
          {filtered.map(p => {
            const days = daysRemaining(p.expires_at);
            const urgent = days !== null && days <= 3 && p.status === 'approved';
            const canDelete = p.status === 'pending' || p.status === 'rejected' || p.status === 'expired';
            const canEdit = p.status === 'pending' || p.status === 'rejected';

            return (
              <div key={p.id} className={`bg-white dark:bg-zinc-900 border rounded-2xl p-4 hover:shadow-md transition-all ${
                urgent ? 'border-red-200 dark:border-red-800' : 'border-zinc-200 dark:border-zinc-800'
              }`}>
                <div className="flex gap-4">
                  {p.images?.[0] ? (
                    <img src={resolveImageUrl(p.images[0].image_url)} alt={p.title}
                      className="w-16 h-16 rounded-xl object-cover bg-zinc-100 shrink-0" />
                  ) : (
                    <div className="w-16 h-16 rounded-xl bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center shrink-0">
                      <ImageIcon size={20} className="text-zinc-300" />
                    </div>
                  )}

                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2">
                      <div className="min-w-0">
                        <p className="font-bold text-zinc-900 dark:text-white truncate">{p.title}</p>
                        <p className="text-sm text-orange-600 font-bold mt-0.5">ETB {Number(p.price).toLocaleString()}</p>
                      </div>
                      <span className={`text-[10px] font-black px-2 py-0.5 rounded-full uppercase shrink-0 ${STATUS_STYLE[p.status] ?? 'bg-zinc-100 text-zinc-500'}`}>
                        {p.status}
                      </span>
                    </div>

                    {/* Time remaining + fee row */}
                    <div className="flex flex-wrap items-center gap-3 mt-1.5">
                      {days !== null && p.status === 'approved' && (
                        <span className={`flex items-center gap-1 text-xs font-bold ${urgent ? 'text-red-500' : 'text-green-600'}`}>
                          <Clock size={11} />
                          {days === 0 ? 'Expires today!' : `${days} days left`}
                        </span>
                      )}
                      {p.status === 'expired' && (
                        <span className="flex items-center gap-1 text-xs text-zinc-400">
                          <Clock size={11} /> Expired {p.expires_at ? new Date(p.expires_at).toLocaleDateString() : ''}
                        </span>
                      )}
                      {p.total_fee_amount != null && (
                        <span className={`flex items-center gap-1 text-xs font-bold ${
                          p.fee_paid ? 'text-green-500' : p.fee_payment_pending ? 'text-yellow-500' : 'text-red-500'
                        }`}>
                          <DollarSign size={11} />
                          Fee: ETB {Number(p.total_fee_amount).toLocaleString()}
                          {p.total_fee_percent != null && ` (${p.total_fee_percent}%)`}
                          {p.fee_paid ? ' ✓ verified' : p.fee_payment_pending ? ' ⏳ pending verification' : ' — unpaid'}
                        </span>
                      )}
                      <span className="text-xs text-zinc-400">{new Date(p.created_at).toLocaleDateString()}</span>
                    </div>

                    {p.admin_notes && (
                      <div className="mt-2 flex items-start gap-1.5 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-2.5 py-1.5">
                        <AlertCircle size={12} className="text-amber-600 shrink-0 mt-0.5" />
                        <p className="text-xs text-amber-700 dark:text-amber-400">{p.admin_notes}</p>
                      </div>
                    )}
                  </div>
                </div>

                {/* Actions */}
                <div className="flex gap-2 mt-3 pt-3 border-t border-zinc-100 dark:border-zinc-800">
                  {canEdit && (
                    <button onClick={() => setEditing(p)}
                      className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-zinc-100 dark:bg-zinc-800 hover:bg-orange-100 dark:hover:bg-orange-900/30 text-zinc-600 dark:text-zinc-400 hover:text-orange-600 text-xs font-bold transition-all">
                      <Pencil size={12} /> Edit
                    </button>
                  )}
                  {canDelete && (
                    <button onClick={() => handleDelete(p.id)}
                      className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-bold transition-all ${
                        confirmDelete === p.id
                          ? 'bg-red-500 text-white hover:bg-red-600'
                          : 'bg-zinc-100 dark:bg-zinc-800 hover:bg-red-100 dark:hover:bg-red-900/30 text-zinc-600 dark:text-zinc-400 hover:text-red-600'
                      }`}>
                      <Trash2 size={12} />
                      {confirmDelete === p.id ? 'Confirm delete?' : 'Delete'}
                    </button>
                  )}
                  {confirmDelete === p.id && (
                    <button onClick={() => setConfirmDelete(null)}
                      className="px-3 py-1.5 rounded-lg text-xs font-bold text-zinc-400 hover:text-zinc-600 transition-colors">
                      Cancel
                    </button>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}

      {showCreate && (
        <Modal title="Add New Product" onClose={() => { setShowCreate(false); setIntentCondition(undefined); }}>
          <ProductForm
            initialCondition={intentCondition}
            onSuccess={() => {
              setShowCreate(false);
              setIntentCondition(undefined);
              // If opened via ?new=1 intent (from dashboard/sell page), go to dashboard after success
              if (searchParams.get('new') === '1') {
                router.push('/seller/dashboard');
              } else {
                load();
              }
            }}
          />
        </Modal>
      )}
      {editing && (
        <Modal title="Edit Product" onClose={() => setEditing(null)}>
          <ProductForm product={editing} onSuccess={() => { setEditing(null); load(); }} />
        </Modal>
      )}
    </div>
  );
}

export default function SellerProductsPage() {
  return (
    <Suspense fallback={<div className="flex items-center justify-center h-64"><Spinner /></div>}>
      <SellerProductsContent />
    </Suspense>
  );
}
