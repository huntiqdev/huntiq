'use client';

import { useEffect, useState, useCallback } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Product } from '@/types';
import { Spinner } from '@/components/ui/Spinner';
import { ProductReview } from '@/components/admin/ProductReview';
import { 
  Package, ExternalLink, RefreshCw, Filter, CheckCircle2, XCircle, 
  DollarSign, ShieldCheck, Tag, LayoutGrid, List, AlertTriangle, Layers
} from 'lucide-react';
import Link from 'next/link';

type TabType = 'pending' | 'all' | 'fee_claims';
type DomainFilter = 'all' | 'sell' | 'rent' | 'lease';

export default function AdminProductsPage() {
  const { user } = useAuth();
  const [activeTab, setActiveTab] = useState<TabType>('pending');
  const [domainFilter, setDomainFilter] = useState<DomainFilter>('all');
  const [products, setProducts] = useState<Product[]>([]);
  const [feeClaims, setFeeClaims] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [viewMode, setViewMode] = useState<'grid' | 'table'>('grid');
  const [selectedIds, setSelectedIds] = useState<number[]>([]);
  const [processingBulk, setProcessingBulk] = useState(false);

  const loadData = useCallback(async () => {
    if (!user) return;
    setLoading(true);
    try {
      if (activeTab === 'pending') {
        const res = await api.adminGetPending(user.token);
        setProducts(res || []);
      } else if (activeTab === 'all') {
        const res = await api.adminGetAllProducts(user.token);
        setProducts(res || []);
      } else if (activeTab === 'fee_claims') {
        const res = await api.adminGetPendingFeePayments(user.token);
        setFeeClaims(res || []);
      }
    } catch (err) {
      console.error('Failed to load admin product data', err);
    } finally {
      setLoading(false);
    }
  }, [user, activeTab]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleReview = async (id: number, status: 'approved' | 'rejected', notes: string) => {
    if (!user) return;
    try {
      await api.adminReviewProduct(id, { status, notes }, user.token);
      setProducts((prev) => prev.filter((p) => p.id !== id));
      setSelectedIds((prev) => prev.filter((i) => i !== id));
    } catch (err) {
      alert('Failed to review product');
    }
  };

  const handleVerifyFee = async (id: number) => {
    if (!user) return;
    try {
      await api.adminVerifyFeePayment(id, user.token);
      setFeeClaims((prev) => prev.filter((c) => c.product_id !== id));
    } catch (err) {
      alert('Failed to verify fee payment');
    }
  };

  const handleRejectFee = async (id: number) => {
    if (!user) return;
    const reason = prompt('Reason for rejecting fee payment claim:');
    if (reason === null) return;
    try {
      await api.adminRejectFeePayment(id, reason, user.token);
      setFeeClaims((prev) => prev.filter((c) => c.product_id !== id));
    } catch (err) {
      alert('Failed to reject fee payment');
    }
  };

  const handleBulkAction = async (status: 'approved' | 'rejected') => {
    if (!user || selectedIds.length === 0) return;
    const notes = status === 'rejected' ? (prompt('Rejection note for bulk action:') || 'Bulk rejected') : 'Bulk approved';
    setProcessingBulk(true);
    try {
      for (const id of selectedIds) {
        await api.adminReviewProduct(id, { status, notes }, user.token);
      }
      setProducts((prev) => prev.filter((p) => !selectedIds.includes(p.id)));
      setSelectedIds([]);
    } catch (err) {
      alert('Some items failed during bulk action');
    } finally {
      setProcessingBulk(false);
    }
  };

  const filteredProducts = products.filter((p) => {
    if (domainFilter === 'all') return true;
    return (p.listing_type || 'sell').toLowerCase() === domainFilter;
  });

  const toggleSelect = (id: number) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const toggleSelectAll = () => {
    if (selectedIds.length === filteredProducts.length) {
      setSelectedIds([]);
    } else {
      setSelectedIds(filteredProducts.map((p) => p.id));
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-blue-600 rounded-2xl flex items-center justify-center text-white shadow-lg shadow-blue-600/20">
            <Package size={24} />
          </div>
          <div>
            <h1 className="text-2xl font-black text-white leading-none">Catalog & Listing Moderation</h1>
            <p className="text-xs font-bold text-zinc-400 uppercase tracking-widest mt-1">
              Multi-Vertical Marketplace & Rental Inventory
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-1 flex items-center gap-1">
            <button
              onClick={() => setViewMode('grid')}
              className={`p-2 rounded-lg transition-colors ${viewMode === 'grid' ? 'bg-zinc-800 text-white' : 'text-zinc-400 hover:text-white'}`}
              title="Grid View"
            >
              <LayoutGrid size={18} />
            </button>
            <button
              onClick={() => setViewMode('table')}
              className={`p-2 rounded-lg transition-colors ${viewMode === 'table' ? 'bg-zinc-800 text-white' : 'text-zinc-400 hover:text-white'}`}
              title="Table View"
            >
              <List size={18} />
            </button>
          </div>

          <button
            onClick={loadData}
            className="p-2.5 rounded-xl bg-zinc-900 border border-zinc-800 text-zinc-300 hover:bg-zinc-800 transition-colors"
          >
            <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
          </button>
        </div>
      </div>

      {/* Primary Navigation Tabs */}
      <div className="flex flex-wrap items-center justify-between gap-4 border-b border-zinc-800 pb-3">
        <div className="flex items-center gap-2">
          <button
            onClick={() => setActiveTab('pending')}
            className={`px-4 py-2.5 rounded-xl font-extrabold text-xs uppercase tracking-wider flex items-center gap-2 transition-all ${
              activeTab === 'pending'
                ? 'bg-orange-600 text-white shadow-lg shadow-orange-600/20'
                : 'bg-zinc-900 text-zinc-400 hover:bg-zinc-800 hover:text-white'
            }`}
          >
            <AlertTriangle size={15} />
            Pending Review
          </button>

          <button
            onClick={() => setActiveTab('fee_claims')}
            className={`px-4 py-2.5 rounded-xl font-extrabold text-xs uppercase tracking-wider flex items-center gap-2 transition-all ${
              activeTab === 'fee_claims'
                ? 'bg-amber-600 text-white shadow-lg shadow-amber-600/20'
                : 'bg-zinc-900 text-zinc-400 hover:bg-zinc-800 hover:text-white'
            }`}
          >
            <DollarSign size={15} />
            Fee Payment Claims
          </button>

          <button
            onClick={() => setActiveTab('all')}
            className={`px-4 py-2.5 rounded-xl font-extrabold text-xs uppercase tracking-wider flex items-center gap-2 transition-all ${
              activeTab === 'all'
                ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/20'
                : 'bg-zinc-900 text-zinc-400 hover:bg-zinc-800 hover:text-white'
            }`}
          >
            <Layers size={15} />
            All Inventory
          </button>
        </div>

        {/* Domain Vertical Filter */}
        {activeTab !== 'fee_claims' && (
          <div className="flex items-center gap-1.5 bg-zinc-900 border border-zinc-800 rounded-xl p-1">
            <span className="text-[10px] font-black uppercase text-zinc-500 px-2 flex items-center gap-1">
              <Filter size={12} /> Domain:
            </span>
            {(['all', 'sell', 'rent', 'lease'] as DomainFilter[]).map((df) => (
              <button
                key={df}
                onClick={() => setDomainFilter(df)}
                className={`px-3 py-1.5 rounded-lg text-xs font-black uppercase tracking-wider transition-colors ${
                  domainFilter === df
                    ? 'bg-zinc-800 text-orange-400 border border-zinc-700'
                    : 'text-zinc-400 hover:text-white'
                }`}
              >
                {df}
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Floating Bulk Action Bar */}
      {selectedIds.length > 0 && activeTab === 'pending' && (
        <div className="bg-orange-950/80 border border-orange-500/30 rounded-2xl p-4 flex items-center justify-between shadow-2xl backdrop-blur-md">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-full bg-orange-600 text-white font-black flex items-center justify-center text-xs">
              {selectedIds.length}
            </div>
            <span className="text-xs font-black uppercase text-orange-200 tracking-wider">
              {selectedIds.length} listings selected for bulk review
            </span>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={() => handleBulkAction('approved')}
              disabled={processingBulk}
              className="bg-emerald-600 hover:bg-emerald-500 text-white font-black text-xs uppercase tracking-wider px-4 py-2 rounded-xl flex items-center gap-1.5 shadow-lg shadow-emerald-600/20"
            >
              <CheckCircle2 size={15} />
              Bulk Approve
            </button>
            <button
              onClick={() => handleBulkAction('rejected')}
              disabled={processingBulk}
              className="bg-red-600 hover:bg-red-500 text-white font-black text-xs uppercase tracking-wider px-4 py-2 rounded-xl flex items-center gap-1.5 shadow-lg shadow-red-600/20"
            >
              <XCircle size={15} />
              Bulk Reject
            </button>
          </div>
        </div>
      )}

      {/* Content Rendering */}
      {loading ? (
        <div className="flex items-center justify-center h-64">
          <Spinner />
        </div>
      ) : activeTab === 'fee_claims' ? (
        /* Fee Claims Table */
        feeClaims.length === 0 ? (
          <div className="bg-zinc-900 border border-zinc-800 rounded-3xl p-12 flex flex-col items-center justify-center text-zinc-500">
            <ShieldCheck size={48} className="opacity-20 mb-4 text-emerald-400" />
            <p className="font-bold text-zinc-400">No pending fee payment verification claims</p>
          </div>
        ) : (
          <div className="bg-zinc-900 border border-zinc-800 rounded-3xl overflow-hidden shadow-xl">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-zinc-950/80 text-zinc-400 text-[10px] uppercase font-black tracking-widest border-b border-zinc-800">
                  <th className="py-4 px-6">Product</th>
                  <th className="py-4 px-4">Price</th>
                  <th className="py-4 px-4">Required Fee</th>
                  <th className="py-4 px-4">Seller Claim Details</th>
                  <th className="py-4 px-6 text-right">Verification Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-800 text-xs">
                {feeClaims.map((claim) => (
                  <tr key={claim.product_id} className="hover:bg-zinc-800/40 transition-colors">
                    <td className="py-4 px-6">
                      <p className="font-black text-white line-clamp-1">{claim.title}</p>
                      <p className="text-[10px] font-mono text-zinc-500">ID: #{claim.product_id}</p>
                    </td>
                    <td className="py-4 px-4 font-mono font-bold text-zinc-300">
                      ETB {Number(claim.price).toLocaleString()}
                    </td>
                    <td className="py-4 px-4 font-mono font-black text-emerald-400">
                      ETB {Number(claim.fee_amount || 0).toLocaleString()}
                    </td>
                    <td className="py-4 px-4 max-w-xs">
                      <p className="text-zinc-300 line-clamp-2">{claim.payment_claim || 'Seller marked payment as completed'}</p>
                      <p className="text-[10px] text-zinc-500 mt-1">Seller: {claim.seller_name} ({claim.seller_phone})</p>
                    </td>
                    <td className="py-4 px-6 text-right">
                      <div className="flex items-center justify-end gap-2">
                        <button
                          onClick={() => handleVerifyFee(claim.product_id)}
                          className="bg-emerald-600 hover:bg-emerald-500 text-white font-black text-xs px-3 py-1.5 rounded-lg transition-colors flex items-center gap-1 shadow-md shadow-emerald-600/20"
                        >
                          <CheckCircle2 size={14} /> Verify & Approve Fee
                        </button>
                        <button
                          onClick={() => handleRejectFee(claim.product_id)}
                          className="bg-zinc-800 hover:bg-red-600 text-zinc-400 hover:text-white font-bold text-xs px-3 py-1.5 rounded-lg transition-colors"
                        >
                          Reject Claim
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )
      ) : filteredProducts.length === 0 ? (
        <div className="bg-zinc-900 border border-zinc-800 rounded-3xl p-12 flex flex-col items-center justify-center text-zinc-500">
          <Package size={48} className="opacity-20 mb-4" />
          <p className="font-bold text-zinc-400">No products matching current tab and domain filter</p>
        </div>
      ) : viewMode === 'grid' ? (
        /* Grid View */
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {filteredProducts.map((p) => {
            const isSelected = selectedIds.includes(p.id);
            return (
              <div
                key={p.id}
                className={`bg-zinc-900 border rounded-3xl overflow-hidden flex flex-col transition-all ${
                  isSelected ? 'border-orange-500 ring-2 ring-orange-500/20' : 'border-zinc-800'
                }`}
              >
                <div className="aspect-square bg-zinc-950 relative overflow-hidden">
                  {p.images && p.images.length > 0 ? (
                    <img src={p.images[0].image_url} alt={p.title} className="w-full h-full object-cover" />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-zinc-700">
                      <Package size={32} />
                    </div>
                  )}

                  {activeTab === 'pending' && (
                    <button
                      onClick={() => toggleSelect(p.id)}
                      className={`absolute top-3 left-3 w-6 h-6 rounded-md flex items-center justify-center border transition-all ${
                        isSelected
                          ? 'bg-orange-600 border-orange-500 text-white'
                          : 'bg-black/60 border-zinc-700 text-transparent hover:border-orange-400'
                      }`}
                    >
                      ✓
                    </button>
                  )}

                  <span className="absolute bottom-3 left-3 bg-zinc-950/80 backdrop-blur px-2.5 py-1 rounded-lg text-[10px] font-black uppercase text-orange-400 border border-zinc-800">
                    {p.listing_type || 'SELL'}
                  </span>

                  <Link
                    href={`/products/${p.id}`}
                    target="_blank"
                    className="absolute top-3 right-3 w-8 h-8 bg-black/60 backdrop-blur rounded-full flex items-center justify-center text-white hover:bg-black/80 transition-colors"
                  >
                    <ExternalLink size={14} />
                  </Link>
                </div>

                <div className="p-5 flex-1 flex flex-col gap-3">
                  <div>
                    <h3 className="font-black text-white text-base line-clamp-1">{p.title}</h3>
                    <p className="text-sm font-mono font-bold text-orange-500 mt-0.5">
                      ETB {Number(p.price).toLocaleString()}
                    </p>
                  </div>

                  <div className="text-xs text-zinc-400 space-y-1 bg-zinc-950/50 p-3 rounded-2xl border border-zinc-800/80">
                    <p className="flex justify-between">
                      <span className="text-zinc-500">Seller ID:</span>
                      <span className="font-mono text-zinc-300">#{p.seller_id}</span>
                    </p>
                    <p className="flex justify-between">
                      <span className="text-zinc-500">Condition:</span>
                      <span className="font-bold text-zinc-300 uppercase">{p.condition || 'USED'}</span>
                    </p>
                    <p className="flex justify-between">
                      <span className="text-zinc-500">Status:</span>
                      <span className={`font-bold uppercase ${p.status === 'approved' ? 'text-emerald-400' : p.status === 'pending' ? 'text-amber-400' : 'text-red-400'}`}>
                        {p.status}
                      </span>
                    </p>
                  </div>

                  {activeTab === 'pending' && (
                    <div className="mt-auto pt-3 border-t border-zinc-800">
                      <ProductReview productId={p.id} onReview={handleReview} />
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      ) : (
        /* Table View */
        <div className="bg-zinc-900 border border-zinc-800 rounded-3xl overflow-hidden shadow-xl">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-zinc-950/80 text-zinc-400 text-[10px] uppercase font-black tracking-widest border-b border-zinc-800">
                {activeTab === 'pending' && (
                  <th className="py-4 px-4 w-10">
                    <input
                      type="checkbox"
                      checked={selectedIds.length === filteredProducts.length && filteredProducts.length > 0}
                      onChange={toggleSelectAll}
                      className="rounded border-zinc-700 bg-zinc-900"
                    />
                  </th>
                )}
                <th className="py-4 px-6">Product & Domain</th>
                <th className="py-4 px-4">Price</th>
                <th className="py-4 px-4">Seller ID</th>
                <th className="py-4 px-4">Status</th>
                {activeTab === 'pending' && <th className="py-4 px-6 text-right">Moderation</th>}
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800 text-xs">
              {filteredProducts.map((p) => (
                <tr key={p.id} className="hover:bg-zinc-800/40 transition-colors">
                  {activeTab === 'pending' && (
                    <td className="py-4 px-4">
                      <input
                        type="checkbox"
                        checked={selectedIds.includes(p.id)}
                        onChange={() => toggleSelect(p.id)}
                        className="rounded border-zinc-700 bg-zinc-900"
                      />
                    </td>
                  )}
                  <td className="py-4 px-6">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-xl bg-zinc-950 overflow-hidden flex-shrink-0 border border-zinc-800">
                        {p.images && p.images.length > 0 ? (
                          <img src={p.images[0].image_url} alt={p.title} className="w-full h-full object-cover" />
                        ) : (
                          <div className="w-full h-full flex items-center justify-center text-zinc-600">
                            <Package size={16} />
                          </div>
                        )}
                      </div>
                      <div>
                        <p className="font-black text-white line-clamp-1">{p.title}</p>
                        <span className="text-[9px] font-black uppercase text-orange-400 tracking-wider">
                          {p.listing_type || 'SELL'}
                        </span>
                      </div>
                    </div>
                  </td>
                  <td className="py-4 px-4 font-mono font-bold text-zinc-300">
                    ETB {Number(p.price).toLocaleString()}
                  </td>
                  <td className="py-4 px-4 font-mono text-zinc-400">#{p.seller_id}</td>
                  <td className="py-4 px-4">
                    <span className={`font-bold uppercase text-[10px] px-2.5 py-1 rounded-lg border ${
                      p.status === 'approved' 
                        ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' 
                        : p.status === 'pending' 
                        ? 'bg-amber-500/10 text-amber-400 border-amber-500/20' 
                        : 'bg-red-500/10 text-red-400 border-red-500/20'
                    }`}>
                      {p.status}
                    </span>
                  </td>
                  {activeTab === 'pending' && (
                    <td className="py-4 px-6 text-right">
                      <ProductReview productId={p.id} onReview={handleReview} />
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
