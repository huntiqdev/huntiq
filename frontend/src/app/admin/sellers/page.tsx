'use client';

import { useEffect, useState, useCallback, useTransition } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import Link from 'next/link';
import {
  Store,
  Search,
  Filter,
  RefreshCw,
  CheckCircle2,
  XCircle,
  AlertTriangle,
  Eye,
  ShieldCheck,
  ShieldAlert,
  Phone,
  Mail,
  Calendar,
  DollarSign,
  Crown,
  Sparkles,
  ChevronLeft,
  ChevronRight,
  ExternalLink,
  X,
  AlertCircle,
  ArrowUpDown,
  Building2,
  Star,
  Home,
  Layers,
} from 'lucide-react';

interface Seller {
  id: number;
  name: string;
  phone: string;
  additional_phones?: string[];
  business_name?: string;
  email?: string;
  rating: number;
  role: string;
  is_active: boolean;
  deactivation_reason?: string;
  tier: string;
  total_value_sold: number;
  net_profit?: number;
  seller_preference?: string;
  capabilities?: string[];
  created_at: string;
  updated_at: string;
}

export default function AdminSellersPage() {
  const {
    user,
    isPlatformOwner,
    isPlatformAdmin,
    isOperationsAdmin,
    isSuperAdmin,
    hasAdminPermission,
  } = useAuth();

  const canManageSellers =
    isPlatformOwner ||
    isPlatformAdmin ||
    isOperationsAdmin ||
    isSuperAdmin ||
    hasAdminPermission('sellers');

  // Search & Filter state
  const [searchInput, setSearchInput] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [personaFilter, setPersonaFilter] = useState<'all' | 'sale' | 'rent' | 'both'>('all');
  const [statusFilter, setStatusFilter] = useState<'all' | 'active' | 'inactive'>('all');
  const [tierFilter, setTierFilter] = useState<'all' | 'normal' | 'gold' | 'premium'>('all');
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(25);

  // Data state
  const [sellers, setSellers] = useState<Seller[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Inspection Drawer & Status Modal state
  const [inspectSeller, setInspectSeller] = useState<Seller | null>(null);
  const [statusModal, setStatusModal] = useState<{ seller: Seller; targetActive: boolean } | null>(null);
  const [reasonInput, setReasonInput] = useState('');
  const [reasonError, setReasonError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  // Toast notifications
  const [toast, setToast] = useState<{ type: 'success' | 'error'; message: string } | null>(null);

  const showToast = (type: 'success' | 'error', message: string) => {
    setToast({ type, message });
    setTimeout(() => setToast(null), 4000);
  };

  // Debounce search input
  useEffect(() => {
    const timer = setTimeout(() => {
      setSearchQuery(searchInput.trim());
      setPage(1);
    }, 350);
    return () => clearTimeout(timer);
  }, [searchInput]);

  // Load Sellers
  const loadSellers = useCallback(async () => {
    if (!user?.token) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.adminListSellers(
        {
          page,
          limit,
          search: searchQuery || undefined,
          status: statusFilter === 'all' ? undefined : statusFilter,
          tier: tierFilter === 'all' ? undefined : tierFilter,
          strategy: personaFilter === 'all' ? undefined : personaFilter,
        },
        user.token
      );
      setSellers(Array.isArray(res?.data) ? res.data : []);
      setTotalCount(res?.total ?? (Array.isArray(res?.data) ? res.data.length : 0));
    } catch (err: any) {
      console.error('Failed to load sellers', err);
      setError(err?.message || 'Unable to connect to seller management service. Please retry.');
      setSellers([]);
    } finally {
      setLoading(false);
    }
  }, [user?.token, page, limit, searchQuery, statusFilter, tierFilter, personaFilter]);

  useEffect(() => {
    loadSellers();
  }, [loadSellers]);

  // Handle status toggle (Activation / Deactivation)
  const handleStatusChangeSubmit = async () => {
    if (!statusModal || !user?.token) return;
    const { seller, targetActive } = statusModal;

    if (!targetActive && !reasonInput.trim()) {
      setReasonError('A valid operational reason is required to deactivate a seller.');
      return;
    }

    setSubmitting(true);
    setReasonError('');
    try {
      await api.adminSetSellerActive(seller.id, targetActive, reasonInput.trim(), user.token);
      showToast(
        'success',
        `Seller #${seller.id} (${seller.name}) successfully ${
          targetActive ? 'reactivated' : 'deactivated'
        }.`
      );
      setStatusModal(null);
      setReasonInput('');
      // If the inspected seller was modified, update their local state too
      if (inspectSeller?.id === seller.id) {
        setInspectSeller({
          ...inspectSeller,
          is_active: targetActive,
          deactivation_reason: targetActive ? '' : reasonInput.trim(),
        });
      }
      loadSellers();
    } catch (err: any) {
      console.error('Status change error', err);
      setReasonError(err?.message || 'Failed to update seller status. Please try again.');
    } finally {
      setSubmitting(false);
    }
  };

  const totalPages = Math.max(1, Math.ceil(totalCount / limit));

  // Quick summary counts
  const activeCount = sellers.filter((s) => s.is_active).length;
  const deactivatedCount = sellers.filter((s) => !s.is_active).length;
  const goodsSellersCount = sellers.filter(
    (s) => (s.seller_preference || 'sale') === 'sale' || s.seller_preference === 'both'
  ).length;
  const landlordsCount = sellers.filter(
    (s) => s.seller_preference === 'rent' || s.seller_preference === 'both'
  ).length;
  const dualCount = sellers.filter((s) => s.seller_preference === 'both').length;

  // With server-side strategy filtering, sellers is already filtered by personaFilter
  const displayedSellers = sellers;

  return (
    <div className="space-y-6">
      {/* Toast Notification */}
      {toast && (
        <div
          className={`fixed bottom-6 right-6 z-50 flex items-center gap-3 px-4 py-3 rounded-2xl shadow-2xl border text-sm font-semibold transition-all animate-in fade-in slide-in-from-bottom-5 ${
            toast.type === 'success'
              ? 'bg-emerald-950/90 border-emerald-500/30 text-emerald-300'
              : 'bg-red-950/90 border-red-500/30 text-red-300'
          }`}
        >
          {toast.type === 'success' ? (
            <CheckCircle2 size={18} className="text-emerald-400 shrink-0" />
          ) : (
            <AlertCircle size={18} className="text-red-400 shrink-0" />
          )}
          <span>{toast.message}</span>
          <button
            onClick={() => setToast(null)}
            className="ml-2 text-zinc-400 hover:text-white transition-colors"
          >
            <X size={15} />
          </button>
        </div>
      )}

      {/* Header Banner */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-gradient-to-br from-orange-600 via-amber-600 to-teal-600 text-white rounded-2xl shadow-lg shadow-orange-600/20">
            <Store size={26} />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-2xl font-black text-white tracking-tight">Providers Directory: Sellers & Landlords</h1>
              <span className="px-2.5 py-0.5 rounded-full text-[11px] font-black uppercase tracking-wider bg-orange-600/10 text-orange-400 border border-orange-500/20">
                {totalCount} Total Providers
              </span>
            </div>
            <p className="text-xs text-zinc-400 mt-0.5">
              Discover, audit, inspect, and manage verified goods sellers, rental landlords, and hybrid providers across Ethiopia.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={() => loadSellers()}
            disabled={loading}
            className="flex items-center gap-2 px-3.5 py-2 bg-zinc-900 hover:bg-zinc-800 text-zinc-300 border border-zinc-800 rounded-xl text-xs font-bold transition-all disabled:opacity-50"
            title="Refresh directory data"
          >
            <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
            <span>Refresh</span>
          </button>
        </div>
      </div>

      {/* Metric Tiles */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="p-4 bg-zinc-900/60 border border-zinc-800/80 rounded-2xl">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-bold text-zinc-500 uppercase tracking-wider">
              Total Providers
            </span>
            <Store size={16} className="text-zinc-400" />
          </div>
          <div className="text-2xl font-black text-white mt-2">{totalCount}</div>
          <div className="text-[10px] text-zinc-500 mt-1 font-medium">All registered supply partners</div>
        </div>

        <div className="p-4 bg-zinc-900/60 border border-zinc-800/80 rounded-2xl">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-bold text-orange-400 uppercase tracking-wider">
              Goods Sellers (Sale)
            </span>
            <Store size={16} className="text-orange-400" />
          </div>
          <div className="text-2xl font-black text-orange-400 mt-2">
            {loading ? '-' : goodsSellersCount}
          </div>
          <div className="text-[10px] text-zinc-500 mt-1 font-medium">Retail & wholesale merchants</div>
        </div>

        <div className="p-4 bg-zinc-900/60 border border-zinc-800/80 rounded-2xl">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-bold text-teal-400 uppercase tracking-wider">
              Landlords & Hosts (Rent)
            </span>
            <Home size={16} className="text-teal-400" />
          </div>
          <div className="text-2xl font-black text-teal-400 mt-2">
            {loading ? '-' : landlordsCount}
          </div>
          <div className="text-[10px] text-zinc-500 mt-1 font-medium">Properties & rental assets ({dualCount} dual)</div>
        </div>

        <div className="p-4 bg-zinc-900/60 border border-zinc-800/80 rounded-2xl">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-bold text-emerald-400 uppercase tracking-wider">
              Active Authorized
            </span>
            <CheckCircle2 size={16} className="text-emerald-400" />
          </div>
          <div className="text-2xl font-black text-emerald-400 mt-2">
            {loading ? '-' : activeCount}
          </div>
          <div className="text-[10px] text-zinc-500 mt-1 font-medium">Authorized to list & trade</div>
        </div>
      </div>

      {/* Search & Filters Controls */}
      <div className="p-4 bg-zinc-900/80 border border-zinc-800 rounded-2xl space-y-3">
        <div className="flex flex-col md:flex-row gap-3">
          {/* Search Input */}
          <div className="relative flex-1">
            <Search
              size={16}
              className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-500"
            />
            <input
              type="text"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              placeholder="Search by merchant name, business name, phone, or email..."
              className="w-full bg-zinc-950 border border-zinc-800 focus:border-orange-500 text-zinc-200 placeholder-zinc-500 text-xs rounded-xl pl-10 pr-9 py-2.5 outline-none transition-colors"
            />
            {searchInput && (
              <button
                onClick={() => setSearchInput('')}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-white"
              >
                <X size={14} />
              </button>
            )}
          </div>

          {/* Provider Strategy Filter */}
          <div className="flex items-center gap-1.5 overflow-x-auto pb-1 md:pb-0">
            {(
              [
                { id: 'all', label: 'All Providers' },
                { id: 'sale', label: '🛒 Goods Sellers' },
                { id: 'rent', label: '🏠 Landlords' },
                { id: 'both', label: '🔄 Dual (Both)' },
              ] as const
            ).map((p) => (
              <button
                key={p.id}
                onClick={() => {
                  setPersonaFilter(p.id);
                  setPage(1);
                }}
                className={`px-3 py-2 rounded-xl text-xs font-bold whitespace-nowrap transition-all ${
                  personaFilter === p.id
                    ? p.id === 'rent'
                      ? 'bg-teal-600 text-white shadow-md shadow-teal-600/20'
                      : p.id === 'both'
                      ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/20'
                      : 'bg-orange-600 text-white shadow-md shadow-orange-600/20'
                    : 'bg-zinc-950 text-zinc-400 hover:text-white hover:bg-zinc-800 border border-zinc-800/80'
                }`}
              >
                {p.label}
              </button>
            ))}
          </div>

          {/* Status Filter */}
          <div className="flex items-center gap-1.5 overflow-x-auto pb-1 md:pb-0">
            {(
              [
                { id: 'all', label: 'All Status' },
                { id: 'active', label: 'Active' },
                { id: 'inactive', label: 'Deactivated' },
              ] as const
            ).map((tab) => (
              <button
                key={tab.id}
                onClick={() => {
                  setStatusFilter(tab.id);
                  setPage(1);
                }}
                className={`px-3 py-2 rounded-xl text-xs font-bold whitespace-nowrap transition-all ${
                  statusFilter === tab.id
                    ? 'bg-zinc-700 text-white shadow-md'
                    : 'bg-zinc-950 text-zinc-400 hover:text-white hover:bg-zinc-800 border border-zinc-800/80'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>

          {/* Tier Filter */}
          <div className="flex items-center gap-1.5 overflow-x-auto pb-1 md:pb-0">
            {(
              [
                { id: 'all', label: 'All Tiers' },
                { id: 'normal', label: 'Standard' },
                { id: 'gold', label: 'Gold' },
                { id: 'premium', label: 'Premium' },
              ] as const
            ).map((t) => (
              <button
                key={t.id}
                onClick={() => {
                  setTierFilter(t.id);
                  setPage(1);
                }}
                className={`px-3 py-2 rounded-xl text-xs font-bold whitespace-nowrap transition-all ${
                  tierFilter === t.id
                    ? 'bg-zinc-800 text-white border border-zinc-700'
                    : 'bg-zinc-950 text-zinc-400 hover:text-white hover:bg-zinc-800 border border-zinc-800/80'
                }`}
              >
                {t.label}
              </button>
            ))}
          </div>
        </div>

        {/* Active Filter Badges */}
        {(searchQuery || statusFilter !== 'all' || tierFilter !== 'all' || personaFilter !== 'all') && (
          <div className="flex items-center gap-2 pt-2 border-t border-zinc-800/60 text-xs">
            <span className="text-zinc-500 font-bold text-[11px] uppercase tracking-wider">
              Filtered By:
            </span>
            {searchQuery && (
              <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-zinc-800 text-zinc-300 text-xs font-medium">
                Query: &quot;{searchQuery}&quot;
                <button onClick={() => setSearchInput('')} className="hover:text-white">
                  <X size={12} />
                </button>
              </span>
            )}
            {personaFilter !== 'all' && (
              <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-zinc-800 text-zinc-300 text-xs font-medium capitalize">
                Strategy: {personaFilter === 'sale' ? 'Goods Sellers' : personaFilter === 'rent' ? 'Landlords' : 'Dual Providers'}
                <button onClick={() => setPersonaFilter('all')} className="hover:text-white">
                  <X size={12} />
                </button>
              </span>
            )}
            {statusFilter !== 'all' && (
              <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-zinc-800 text-zinc-300 text-xs font-medium capitalize">
                Status: {statusFilter}
                <button onClick={() => setStatusFilter('all')} className="hover:text-white">
                  <X size={12} />
                </button>
              </span>
            )}
            {tierFilter !== 'all' && (
              <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-zinc-800 text-zinc-300 text-xs font-medium capitalize">
                Tier: {tierFilter}
                <button onClick={() => setTierFilter('all')} className="hover:text-white">
                  <X size={12} />
                </button>
              </span>
            )}
            <button
              onClick={() => {
                setSearchInput('');
                setPersonaFilter('all');
                setStatusFilter('all');
                setTierFilter('all');
                setPage(1);
              }}
              className="text-xs text-orange-400 hover:underline font-bold ml-auto"
            >
              Reset Filters
            </button>
          </div>
        )}
      </div>

      {/* Recoverable Error Banner */}
      {error && (
        <div className="p-4 bg-red-950/40 border border-red-500/30 rounded-2xl flex items-center justify-between gap-3 text-red-300">
          <div className="flex items-center gap-3">
            <AlertCircle size={20} className="text-red-400 shrink-0" />
            <div>
              <div className="font-bold text-sm">Failed to retrieve merchant directory</div>
              <div className="text-xs text-red-300/80">{error}</div>
            </div>
          </div>
          <button
            onClick={() => loadSellers()}
            className="px-3 py-1.5 bg-red-900/60 hover:bg-red-800 border border-red-500/30 rounded-xl text-xs font-bold text-white transition-colors"
          >
            Retry
          </button>
        </div>
      )}

      {/* Main Sellers Table */}
      <div className="bg-zinc-900/90 border border-zinc-800 rounded-3xl overflow-hidden shadow-2xl">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs border-collapse">
            <thead>
              <tr className="border-b border-zinc-800 bg-zinc-950/50 text-zinc-400 font-bold uppercase tracking-wider text-[11px]">
                <th className="py-4 px-4">Provider</th>
                <th className="py-4 px-4">Contact Info</th>
                <th className="py-4 px-4">Strategy & Role</th>
                <th className="py-4 px-4">Business Profile</th>
                <th className="py-4 px-4">Tier & Rating</th>
                <th className="py-4 px-4">Account Status</th>
                <th className="py-4 px-4">Registered</th>
                <th className="py-4 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/60">
              {loading ? (
                // Skeleton loading rows
                Array.from({ length: 5 }).map((_, idx) => (
                  <tr key={idx} className="animate-pulse">
                    <td className="py-4 px-4">
                      <div className="flex items-center gap-3">
                        <div className="w-9 h-9 rounded-xl bg-zinc-800" />
                        <div className="space-y-1.5">
                          <div className="h-3 w-28 bg-zinc-800 rounded" />
                          <div className="h-2.5 w-16 bg-zinc-800/60 rounded" />
                        </div>
                      </div>
                    </td>
                    <td className="py-4 px-4">
                      <div className="space-y-1.5">
                        <div className="h-3 w-24 bg-zinc-800 rounded" />
                        <div className="h-2.5 w-32 bg-zinc-800/60 rounded" />
                      </div>
                    </td>
                    <td className="py-4 px-4">
                      <div className="h-4 w-24 bg-zinc-800 rounded" />
                    </td>
                    <td className="py-4 px-4">
                      <div className="h-3 w-20 bg-zinc-800 rounded" />
                    </td>
                    <td className="py-4 px-4">
                      <div className="h-3 w-16 bg-zinc-800 rounded" />
                    </td>
                    <td className="py-4 px-4">
                      <div className="h-5 w-20 bg-zinc-800 rounded-full" />
                    </td>
                    <td className="py-4 px-4">
                      <div className="h-3 w-16 bg-zinc-800 rounded" />
                    </td>
                    <td className="py-4 px-4 text-right">
                      <div className="h-7 w-20 bg-zinc-800 rounded ml-auto" />
                    </td>
                  </tr>
                ))
              ) : displayedSellers.length === 0 ? (
                <tr>
                  <td colSpan={8} className="py-14 text-center">
                    <div className="max-w-sm mx-auto space-y-3">
                      <div className="w-12 h-12 rounded-2xl bg-zinc-800 flex items-center justify-center text-zinc-500 mx-auto">
                        <Store size={24} />
                      </div>
                      <div className="text-sm font-bold text-white">No providers found</div>
                      <p className="text-xs text-zinc-400">
                        {searchQuery || statusFilter !== 'all' || tierFilter !== 'all' || personaFilter !== 'all'
                          ? 'No merchant or landlord records match your active search filters. Try resetting the filter criteria.'
                          : 'No registered sellers or landlords exist in the marketplace directory yet.'}
                      </p>
                      {(searchQuery || statusFilter !== 'all' || tierFilter !== 'all' || personaFilter !== 'all') && (
                        <button
                          onClick={() => {
                            setSearchInput('');
                            setPersonaFilter('all');
                            setStatusFilter('all');
                            setTierFilter('all');
                          }}
                          className="px-3.5 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-zinc-200 rounded-xl text-xs font-bold transition-colors"
                        >
                          Clear All Filters
                        </button>
                      )}
                    </div>
                  </td>
                </tr>
              ) : (
                displayedSellers.map((seller) => {
                  const isDeactivated = !seller.is_active;
                  const pref = seller.seller_preference || 'sale';

                  return (
                    <tr
                      key={seller.id}
                      className="hover:bg-zinc-800/30 transition-colors group"
                    >
                      {/* Merchant Name & ID */}
                      <td className="py-3.5 px-4">
                        <div className="flex items-center gap-3">
                          <div className={`w-9 h-9 rounded-xl border flex items-center justify-center font-black text-xs shrink-0 group-hover:scale-105 transition-transform ${
                            pref === 'rent'
                              ? 'bg-teal-600/10 border-teal-500/20 text-teal-400'
                              : pref === 'both'
                              ? 'bg-indigo-600/10 border-indigo-500/20 text-indigo-400'
                              : 'bg-orange-600/10 border-orange-500/20 text-orange-400'
                          }`}>
                            {seller.name.slice(0, 2).toUpperCase()}
                          </div>
                          <div>
                            <div className="font-extrabold text-white text-xs flex items-center gap-1.5">
                              {seller.name}
                              {seller.tier === 'premium' && (
                                <Crown size={12} className="text-amber-400" />
                              )}
                            </div>
                            <div className="text-[10px] text-zinc-500 font-mono">
                              #SLR-{seller.id}
                            </div>
                          </div>
                        </div>
                      </td>

                      {/* Contact Info */}
                      <td className="py-3.5 px-4">
                        <div className="space-y-0.5">
                          <div className="flex items-center gap-1 text-zinc-300 font-mono text-[11px]">
                            <Phone size={11} className="text-zinc-500" />
                            {seller.phone || 'No phone'}
                          </div>
                          {seller.email ? (
                            <div className="flex items-center gap-1 text-zinc-400 text-[11px] truncate max-w-[180px]">
                              <Mail size={11} className="text-zinc-500 shrink-0" />
                              {seller.email}
                            </div>
                          ) : (
                            <span className="text-[10px] text-zinc-600 italic">No email linked</span>
                          )}
                        </div>
                      </td>

                      {/* Strategy & Role Badge */}
                      <td className="py-3.5 px-4">
                        {pref === 'rent' ? (
                          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-teal-950/60 border border-teal-500/30 text-teal-300 text-[11px] font-bold">
                            <Home size={12} className="text-teal-400" /> Landlord (Rent)
                          </span>
                        ) : pref === 'both' ? (
                          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-indigo-950/60 border border-indigo-500/30 text-indigo-300 text-[11px] font-bold">
                            <Layers size={12} className="text-indigo-400" /> Dual (Sale & Rent)
                          </span>
                        ) : (
                          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-orange-950/60 border border-orange-500/30 text-orange-300 text-[11px] font-bold">
                            <Store size={12} className="text-orange-400" /> Seller (Sale)
                          </span>
                        )}
                      </td>

                      {/* Business Details */}
                      <td className="py-3.5 px-4">
                        <div className="space-y-0.5">
                          <div className="font-bold text-zinc-200 truncate max-w-[160px]">
                            {seller.business_name || (
                              <span className="text-zinc-500 italic">Individual Partner</span>
                            )}
                          </div>
                          <div className="text-[10px] text-zinc-500 font-mono">
                            {seller.capabilities?.join(', ') || 'merchant'}
                          </div>
                        </div>
                      </td>

                      {/* Tier & Rating */}
                      <td className="py-3.5 px-4">
                        <div className="space-y-1">
                          <span
                            className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-black uppercase tracking-wider border ${
                              seller.tier === 'premium'
                                ? 'bg-amber-950/40 border-amber-500/40 text-amber-300'
                                : seller.tier === 'gold'
                                ? 'bg-yellow-950/40 border-yellow-500/40 text-yellow-300'
                                : 'bg-zinc-800 border-zinc-700 text-zinc-300'
                            }`}
                          >
                            {seller.tier === 'premium' && <Sparkles size={10} />}
                            {seller.tier || 'standard'}
                          </span>
                          <div className="flex items-center gap-1 text-[11px] text-amber-400 font-bold">
                            <Star size={11} className="fill-amber-400 text-amber-400" />
                            <span>{(seller.rating || 0).toFixed(1)}</span>
                          </div>
                        </div>
                      </td>

                      {/* Account Status Badge */}
                      <td className="py-3.5 px-4">
                        {seller.is_active ? (
                          <div className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-950/40 border border-emerald-500/30 text-emerald-400 text-[11px] font-bold">
                            <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
                            Active
                          </div>
                        ) : (
                          <div className="space-y-1">
                            <div className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-red-950/40 border border-red-500/30 text-red-400 text-[11px] font-bold">
                              <span className="w-1.5 h-1.5 rounded-full bg-red-400" />
                              Deactivated
                            </div>
                            {seller.deactivation_reason && (
                              <div
                                className="text-[10px] text-red-300/80 truncate max-w-[150px]"
                                title={seller.deactivation_reason}
                              >
                                &quot;{seller.deactivation_reason}&quot;
                              </div>
                            )}
                          </div>
                        )}
                      </td>

                      {/* Registered Date */}
                      <td className="py-3.5 px-4 text-zinc-400 font-mono text-[11px]">
                        {seller.created_at
                          ? new Date(seller.created_at).toLocaleDateString(undefined, {
                              year: 'numeric',
                              month: 'short',
                              day: 'numeric',
                            })
                          : '—'}
                      </td>

                      {/* Actions */}
                      <td className="py-3.5 px-4 text-right">
                        <div className="flex items-center justify-end gap-1.5">
                          {/* Inspect Profile */}
                          <button
                            onClick={() => setInspectSeller(seller)}
                            className="p-1.5 rounded-lg bg-zinc-800 text-zinc-300 hover:text-white hover:bg-zinc-700 transition-colors"
                            title="Inspect Merchant Profile"
                          >
                            <Eye size={14} />
                          </button>

                          {/* Permission-guarded Status Toggle */}
                          {canManageSellers && (
                            <button
                              onClick={() => {
                                setReasonInput('');
                                setReasonError('');
                                setStatusModal({
                                  seller,
                                  targetActive: isDeactivated,
                                });
                              }}
                              className={`px-2.5 py-1 rounded-lg text-[11px] font-extrabold transition-all border ${
                                isDeactivated
                                  ? 'bg-emerald-950/60 hover:bg-emerald-900 border-emerald-500/30 text-emerald-300'
                                  : 'bg-red-950/60 hover:bg-red-900 border-red-500/30 text-red-300'
                              }`}
                            >
                              {isDeactivated ? 'Activate' : 'Deactivate'}
                            </button>
                          )}
                        </div>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination Controls */}
        <div className="p-4 border-t border-zinc-800 bg-zinc-950/50 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs text-zinc-400">
          <div className="flex items-center gap-2">
            <span>
              Showing{' '}
              <strong className="text-white font-bold">
                {totalCount === 0 ? 0 : (page - 1) * limit + 1}
              </strong>{' '}
              to{' '}
              <strong className="text-white font-bold">
                {Math.min(page * limit, totalCount)}
              </strong>{' '}
              of <strong className="text-white font-bold">{totalCount}</strong> merchants
            </span>

            <select
              value={limit}
              onChange={(e) => {
                setLimit(Number(e.target.value));
                setPage(1);
              }}
              className="bg-zinc-900 border border-zinc-800 rounded-lg px-2 py-1 text-zinc-300 text-xs outline-none ml-2"
            >
              <option value={10}>10 per page</option>
              <option value={25}>25 per page</option>
              <option value={50}>50 per page</option>
              <option value={100}>100 per page</option>
            </select>
          </div>

          <div className="flex items-center gap-1.5">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page <= 1 || loading}
              className="p-1.5 rounded-lg bg-zinc-900 border border-zinc-800 text-zinc-300 hover:text-white hover:bg-zinc-800 disabled:opacity-40 transition-colors"
            >
              <ChevronLeft size={16} />
            </button>
            <span className="px-3 py-1 bg-zinc-900 border border-zinc-800 rounded-lg text-white font-bold font-mono">
              {page} / {totalPages}
            </span>
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page >= totalPages || loading}
              className="p-1.5 rounded-lg bg-zinc-900 border border-zinc-800 text-zinc-300 hover:text-white hover:bg-zinc-800 disabled:opacity-40 transition-colors"
            >
              <ChevronRight size={16} />
            </button>
          </div>
        </div>
      </div>

      {/* Seller Detail Inspection Modal */}
      {inspectSeller && (
        <div className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-in fade-in">
          <div
            className="absolute inset-0"
            onClick={() => setInspectSeller(null)}
          />
          <div className="relative z-10 bg-zinc-900 border border-zinc-800 rounded-3xl shadow-2xl w-full max-w-xl max-h-[90vh] overflow-y-auto">
            {/* Header */}
            <div className="p-6 border-b border-zinc-800 flex items-center justify-between sticky top-0 bg-zinc-900/95 backdrop-blur-sm">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-2xl bg-orange-600/10 border border-orange-500/20 text-orange-400 flex items-center justify-center font-black text-lg">
                  {inspectSeller.name.slice(0, 2).toUpperCase()}
                </div>
                <div>
                  <h3 className="text-lg font-black text-white">{inspectSeller.name}</h3>
                  <div className="flex items-center gap-2 text-xs text-zinc-400">
                    <span className="font-mono text-zinc-500">#SLR-{inspectSeller.id}</span>
                    <span>•</span>
                    <span>{inspectSeller.business_name || 'Individual Merchant'}</span>
                  </div>
                </div>
              </div>
              <button
                onClick={() => setInspectSeller(null)}
                className="p-2 rounded-xl bg-zinc-800 text-zinc-400 hover:text-white hover:bg-zinc-700 transition-colors"
              >
                <X size={18} />
              </button>
            </div>

            {/* Content */}
            <div className="p-6 space-y-6 text-xs">
              {/* Account Status Card */}
              <div
                className={`p-4 rounded-2xl border ${
                  inspectSeller.is_active
                    ? 'bg-emerald-950/20 border-emerald-500/30'
                    : 'bg-red-950/20 border-red-500/30'
                }`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    {inspectSeller.is_active ? (
                      <CheckCircle2 size={18} className="text-emerald-400" />
                    ) : (
                      <XCircle size={18} className="text-red-400" />
                    )}
                    <span className="font-bold text-white text-sm">
                      Account Status: {inspectSeller.is_active ? 'Active' : 'Deactivated'}
                    </span>
                  </div>
                  {canManageSellers && (
                    <button
                      onClick={() => {
                        const target = !inspectSeller.is_active;
                        setInspectSeller(null);
                        setReasonInput('');
                        setReasonError('');
                        setStatusModal({
                          seller: inspectSeller,
                          targetActive: target,
                        });
                      }}
                      className="px-3 py-1 rounded-lg text-xs font-bold bg-zinc-800 hover:bg-zinc-700 text-white border border-zinc-700 transition-colors"
                    >
                      {inspectSeller.is_active ? 'Deactivate Account' : 'Reactivate Account'}
                    </button>
                  )}
                </div>

                {!inspectSeller.is_active && inspectSeller.deactivation_reason && (
                  <div className="mt-3 p-3 bg-red-950/40 rounded-xl border border-red-500/20 text-red-200">
                    <span className="font-bold block text-[11px] uppercase tracking-wider text-red-400">
                      Reason on Record:
                    </span>
                    <span className="italic mt-0.5 block">&quot;{inspectSeller.deactivation_reason}&quot;</span>
                  </div>
                )}
              </div>

              {/* Profile Details Grid */}
              <div className="grid grid-cols-2 gap-4">
                <div className="p-3.5 bg-zinc-950/60 rounded-xl border border-zinc-800/80">
                  <span className="text-[10px] text-zinc-500 uppercase tracking-wider font-bold block">
                    Telephone
                  </span>
                  <a
                    href={`tel:${inspectSeller.phone}`}
                    className="font-mono text-zinc-200 text-sm font-bold hover:text-orange-400 transition-colors mt-0.5 block"
                  >
                    {inspectSeller.phone || '—'}
                  </a>
                </div>

                <div className="p-3.5 bg-zinc-950/60 rounded-xl border border-zinc-800/80">
                  <span className="text-[10px] text-zinc-500 uppercase tracking-wider font-bold block">
                    Email Address
                  </span>
                  {inspectSeller.email ? (
                    <a
                      href={`mailto:${inspectSeller.email}`}
                      className="text-zinc-200 text-sm font-bold hover:text-orange-400 transition-colors mt-0.5 block truncate"
                    >
                      {inspectSeller.email}
                    </a>
                  ) : (
                    <span className="text-zinc-500 italic mt-0.5 block">Not provided</span>
                  )}
                </div>

                <div className="p-3.5 bg-zinc-950/60 rounded-xl border border-zinc-800/80">
                  <span className="text-[10px] text-zinc-500 uppercase tracking-wider font-bold block">
                    Marketplace Tier
                  </span>
                  <span className="text-sm font-black text-amber-400 uppercase tracking-wider mt-0.5 block">
                    {inspectSeller.tier || 'Standard'}
                  </span>
                </div>

                <div className="p-3.5 bg-zinc-950/60 rounded-xl border border-zinc-800/80">
                  <span className="text-[10px] text-zinc-500 uppercase tracking-wider font-bold block">
                    Customer Rating
                  </span>
                  <div className="flex items-center gap-1.5 text-sm font-black text-white mt-0.5">
                    <Star size={14} className="fill-amber-400 text-amber-400" />
                    <span>{(inspectSeller.rating || 0).toFixed(1)} / 5.0</span>
                  </div>
                </div>

                <div className="p-3.5 bg-zinc-950/60 rounded-xl border border-zinc-800/80 col-span-2">
                  <span className="text-[10px] text-zinc-500 uppercase tracking-wider font-bold block">
                    Marketplace Provider Strategy & Role
                  </span>
                  <div className="flex items-center gap-2 mt-1">
                    {inspectSeller.seller_preference === 'rent' ? (
                      <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-teal-950/60 border border-teal-500/30 text-teal-300 text-xs font-bold">
                        <Home size={13} className="text-teal-400" /> Rental Landlord / Host (Rentals)
                      </span>
                    ) : inspectSeller.seller_preference === 'both' ? (
                      <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-indigo-950/60 border border-indigo-500/30 text-indigo-300 text-xs font-bold">
                        <Layers size={13} className="text-indigo-400" /> Dual Provider (Goods Seller & Rental Landlord)
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-orange-950/60 border border-orange-500/30 text-orange-300 text-xs font-bold">
                        <Store size={13} className="text-orange-400" /> Goods Seller (Sales Only)
                      </span>
                    )}
                  </div>
                  <p className="text-[11px] text-zinc-400 mt-1.5">
                    {inspectSeller.seller_preference === 'rent'
                      ? 'Authorized to list residential/commercial properties, land, vehicles, or equipment for rent or lease.'
                      : inspectSeller.seller_preference === 'both'
                      ? 'Authorized for both direct item sales and long/short-term asset leasing across the platform.'
                      : 'Authorized to list physical retail and wholesale products for direct purchase by buyers.'}
                  </p>
                  <div className="flex items-center gap-1.5 mt-2 pt-2 border-t border-zinc-800/60">
                    <span className="text-[10px] text-zinc-500 uppercase font-mono">Capabilities:</span>
                    {(inspectSeller.capabilities && inspectSeller.capabilities.length > 0
                      ? inspectSeller.capabilities
                      : ['buyer', inspectSeller.seller_preference === 'rent' ? 'landlord' : 'seller']
                    ).map((cap) => (
                      <span
                        key={cap}
                        className="px-2 py-0.5 rounded text-[10px] font-mono font-bold uppercase bg-zinc-800 text-zinc-300 border border-zinc-700"
                      >
                        {cap}
                      </span>
                    ))}
                  </div>
                </div>

                <div className="p-3.5 bg-zinc-950/60 rounded-xl border border-zinc-800/80">
                  <span className="text-[10px] text-zinc-500 uppercase tracking-wider font-bold block">
                    Total Volume Sold
                  </span>
                  <span className="font-mono text-sm font-bold text-zinc-200 mt-0.5 block">
                    {(inspectSeller.total_value_sold || 0).toLocaleString()} ETB
                  </span>
                </div>
              </div>

              {/* Net Profit (Restricted to Platform Owner / Super) */}
              {(isPlatformOwner || isSuperAdmin) && typeof inspectSeller.net_profit === 'number' && (
                <div className="p-4 bg-zinc-950/80 rounded-2xl border border-zinc-800 flex items-center justify-between">
                  <div>
                    <span className="text-[10px] font-bold text-orange-400 uppercase tracking-wider block">
                      Platform Net Profit Generated
                    </span>
                    <span className="text-xs text-zinc-500">
                      Calculated merchant fee yield (Owner view only)
                    </span>
                  </div>
                  <div className="text-base font-black font-mono text-white">
                    {inspectSeller.net_profit.toLocaleString()} ETB
                  </div>
                </div>
              )}

              {/* Navigation Deep Link */}
              <div className="pt-2 flex items-center justify-between border-t border-zinc-800">
                <Link
                  href="/admin/products"
                  className="inline-flex items-center gap-2 text-xs font-bold text-orange-400 hover:text-orange-300 transition-colors"
                >
                  <ExternalLink size={14} />
                  <span>Inspect Listings on Products & Listings</span>
                </Link>

                <div className="text-[10px] text-zinc-500 font-mono">
                  Joined: {new Date(inspectSeller.created_at).toLocaleDateString()}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Confirmation & Status Change Modal */}
      {statusModal && (
        <div className="fixed inset-0 bg-black/75 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-in fade-in">
          <div
            className="absolute inset-0"
            onClick={() => !submitting && setStatusModal(null)}
          />
          <div className="relative z-10 bg-zinc-900 border border-zinc-800 rounded-3xl shadow-2xl w-full max-w-md overflow-hidden">
            {/* Modal Header */}
            <div className="p-6 border-b border-zinc-800 flex items-center gap-3">
              <div
                className={`p-2.5 rounded-2xl border ${
                  statusModal.targetActive
                    ? 'bg-emerald-600/10 border-emerald-500/20 text-emerald-400'
                    : 'bg-red-600/10 border-red-500/20 text-red-400'
                }`}
              >
                {statusModal.targetActive ? (
                  <ShieldCheck size={22} />
                ) : (
                  <AlertTriangle size={22} />
                )}
              </div>
              <div>
                <h3 className="text-base font-black text-white">
                  {statusModal.targetActive
                    ? 'Confirm Seller Reactivation'
                    : 'Confirm Account Deactivation'}
                </h3>
                <p className="text-xs text-zinc-400">
                  Target: {statusModal.seller.name} (#SLR-{statusModal.seller.id})
                </p>
              </div>
            </div>

            {/* Modal Body */}
            <div className="p-6 space-y-4 text-xs">
              <div className="text-zinc-300 leading-relaxed">
                {statusModal.targetActive ? (
                  <p>
                    Are you sure you want to reactivate{' '}
                    <strong className="text-white font-bold">{statusModal.seller.name}</strong>?
                    This will restore their ability to publish listings, respond to buyers, and
                    conduct marketplace operations.
                  </p>
                ) : (
                  <p>
                    Deactivating{' '}
                    <strong className="text-white font-bold">{statusModal.seller.name}</strong>{' '}
                    will terminate all active sessions, prohibit new listing creations, and block
                    marketplace merchant capabilities immediately.
                  </p>
                )}
              </div>

              {/* Mandatory Reason for Deactivation */}
              {!statusModal.targetActive && (
                <div className="space-y-1.5">
                  <label className="block text-[11px] font-bold text-zinc-400 uppercase tracking-wider">
                    Deactivation Reason <span className="text-red-400">*</span>
                  </label>
                  <textarea
                    rows={3}
                    value={reasonInput}
                    onChange={(e) => {
                      setReasonInput(e.target.value);
                      setReasonError('');
                    }}
                    placeholder="Enter explicit administrative justification (e.g., suspicious activity, unverified identity, merchant request)..."
                    className="w-full bg-zinc-950 border border-zinc-800 focus:border-red-500 text-zinc-200 placeholder-zinc-500 text-xs rounded-xl p-3 outline-none transition-colors"
                  />
                  {reasonError && (
                    <p className="text-[11px] text-red-400 font-bold">{reasonError}</p>
                  )}
                </div>
              )}

              {/* Optional Notes for Activation */}
              {statusModal.targetActive && (
                <div className="space-y-1.5">
                  <label className="block text-[11px] font-bold text-zinc-400 uppercase tracking-wider">
                    Administrative Notes (Optional)
                  </label>
                  <input
                    type="text"
                    value={reasonInput}
                    onChange={(e) => setReasonInput(e.target.value)}
                    placeholder="E.g., verified identity cleared, compliance resolved"
                    className="w-full bg-zinc-950 border border-zinc-800 focus:border-emerald-500 text-zinc-200 placeholder-zinc-500 text-xs rounded-xl px-3 py-2 outline-none transition-colors"
                  />
                </div>
              )}

              {/* Modal Actions */}
              <div className="flex items-center justify-end gap-3 pt-3 border-t border-zinc-800/80">
                <button
                  type="button"
                  onClick={() => setStatusModal(null)}
                  disabled={submitting}
                  className="px-4 py-2 rounded-xl text-xs font-bold text-zinc-400 hover:text-white hover:bg-zinc-800 transition-colors disabled:opacity-50"
                >
                  Cancel
                </button>
                <button
                  type="button"
                  onClick={handleStatusChangeSubmit}
                  disabled={submitting}
                  className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-black text-white shadow-lg transition-all disabled:opacity-50 ${
                    statusModal.targetActive
                      ? 'bg-emerald-600 hover:bg-emerald-500 shadow-emerald-600/20'
                      : 'bg-red-600 hover:bg-red-500 shadow-red-600/20'
                  }`}
                >
                  {submitting && <RefreshCw size={14} className="animate-spin" />}
                  <span>
                    {statusModal.targetActive ? 'Confirm Reactivation' : 'Confirm Deactivation'}
                  </span>
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
