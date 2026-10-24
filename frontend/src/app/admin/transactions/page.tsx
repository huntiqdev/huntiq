'use client';

import { useEffect, useState, useCallback, useMemo } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import Link from 'next/link';
import {
  Receipt,
  Search,
  Filter,
  RefreshCw,
  CheckCircle2,
  Clock,
  XCircle,
  AlertTriangle,
  Eye,
  ShieldCheck,
  ShieldAlert,
  Phone,
  Calendar,
  DollarSign,
  Copy,
  Check,
  ChevronLeft,
  ChevronRight,
  ExternalLink,
  X,
  CreditCard,
  Building2,
  Lock,
  ArrowUpDown,
  Coins,
  FileSpreadsheet,
} from 'lucide-react';

interface Transaction {
  id: number;
  user_id: number;
  product_id: number;
  seller_id: number;
  buyer_name: string;
  buyer_phone: string;
  amount: number;
  fee_amount: number;
  broker_fee: number;
  seller_gets: number;
  transaction_type: string;
  payment_method?: string;
  payment_reference?: string;
  payer_phone?: string;
  status: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

function fmtETB(n: number | undefined | null): string {
  if (n == null || isNaN(n)) return '0.00 ETB';
  return `${Number(n).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })} ETB`;
}

function formatDate(iso: string | undefined | null): string {
  if (!iso) return '—';
  try {
    const d = new Date(iso);
    return d.toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: true,
    });
  } catch {
    return iso;
  }
}

export default function AdminTransactionsPage() {
  const {
    user,
    isPlatformOwner,
    isPlatformAdmin,
    isOperationsAdmin,
    isSuperAdmin,
    hasAdminPermission,
  } = useAuth();

  const canViewTransactions =
    isPlatformOwner ||
    isPlatformAdmin ||
    isOperationsAdmin ||
    isSuperAdmin ||
    hasAdminPermission('transactions');

  // Search & Filter state
  const [searchInput, setSearchInput] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<'all' | 'completed' | 'pending' | 'cancelled'>('all');
  const [methodFilter, setMethodFilter] = useState<'all' | 'telebirr' | 'safaricom' | 'cbe' | 'bank'>('all');
  const [typeFilter, setTypeFilter] = useState<'all' | 'listing_fee' | 'product_sale'>('all');
  const [fromDate, setFromDate] = useState('');
  const [toDate, setToDate] = useState('');
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(25);

  // Data state
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Inspection Drawer
  const [inspectTx, setInspectTx] = useState<Transaction | null>(null);
  const [copiedKey, setCopiedKey] = useState<string | null>(null);

  // Toast feedback
  const [toast, setToast] = useState<{ type: 'success' | 'error'; message: string } | null>(null);

  const showToast = (type: 'success' | 'error', message: string) => {
    setToast({ type, message });
    setTimeout(() => setToast(null), 3000);
  };

  // Debounce search input
  useEffect(() => {
    const timer = setTimeout(() => {
      setSearchQuery(searchInput.trim());
      setPage(1);
    }, 350);
    return () => clearTimeout(timer);
  }, [searchInput]);

  const loadTransactions = useCallback(async () => {
    if (!user?.token || !canViewTransactions) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.adminListTransactions(
        {
          page,
          limit,
          search: searchQuery || undefined,
          status: statusFilter === 'all' ? undefined : statusFilter,
          method: methodFilter === 'all' ? undefined : methodFilter,
          type: typeFilter === 'all' ? undefined : typeFilter,
          from_date: fromDate || undefined,
          to_date: toDate || undefined,
        },
        user.token
      );
      setTransactions(Array.isArray(res?.data) ? res.data : []);
      setTotalCount(res?.total ?? (Array.isArray(res?.data) ? res.data.length : 0));
    } catch (err: any) {
      console.error('Failed to load transactions ledger', err);
      setError(err?.message || 'Unable to connect to financial ledger service. Please retry.');
      setTransactions([]);
    } finally {
      setLoading(false);
    }
  }, [user?.token, canViewTransactions, page, limit, searchQuery, statusFilter, methodFilter, typeFilter, fromDate, toDate]);

  useEffect(() => {
    loadTransactions();
  }, [loadTransactions]);

  // Accurate page calculations (derived purely from current loaded dataset)
  const pageMetrics = useMemo(() => {
    let grossVolume = 0;
    let brokerFees = 0;
    let completed = 0;
    let pending = 0;
    let cancelled = 0;

    for (const t of transactions) {
      grossVolume += Number(t.amount || 0);
      brokerFees += Number(t.broker_fee || 0);
      if (t.status === 'completed') completed++;
      else if (t.status === 'pending') pending++;
      else if (t.status === 'cancelled') cancelled++;
    }

    return {
      grossVolume,
      brokerFees,
      completed,
      pending,
      cancelled,
    };
  }, [transactions]);

  const handleCopy = (text: string, key: string) => {
    if (!text) return;
    navigator.clipboard.writeText(text);
    setCopiedKey(key);
    setTimeout(() => setCopiedKey(null), 2000);
  };

  const handleResetFilters = () => {
    setSearchInput('');
    setSearchQuery('');
    setStatusFilter('all');
    setMethodFilter('all');
    setTypeFilter('all');
    setFromDate('');
    setToDate('');
    setPage(1);
  };

  const isFiltered = Boolean(
    searchQuery ||
    statusFilter !== 'all' ||
    methodFilter !== 'all' ||
    typeFilter !== 'all' ||
    fromDate ||
    toDate
  );

  const totalPages = Math.max(1, Math.ceil(totalCount / limit));

  // Render provider pill
  const renderProviderBadge = (provider?: string) => {
    const p = (provider || '').toLowerCase();
    switch (p) {
      case 'telebirr':
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-[11px] font-bold bg-cyan-950/60 text-cyan-300 border border-cyan-800/40">
            <span className="w-1.5 h-1.5 rounded-full bg-cyan-400"></span>
            Telebirr
          </span>
        );
      case 'cbe':
      case 'cbe_birr':
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-[11px] font-bold bg-purple-950/60 text-purple-300 border border-purple-800/40">
            <span className="w-1.5 h-1.5 rounded-full bg-purple-400"></span>
            CBE Birr
          </span>
        );
      case 'safaricom':
      case 'mpesa':
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-[11px] font-bold bg-emerald-950/60 text-emerald-300 border border-emerald-800/40">
            <span className="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
            Safaricom M-Pesa
          </span>
        );
      case 'bank':
      case 'bank_transfer':
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-[11px] font-bold bg-amber-950/60 text-amber-300 border border-amber-800/40">
            <Building2 size={11} className="text-amber-400" />
            Bank Transfer
          </span>
        );
      default:
        return (
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-medium bg-zinc-800 text-zinc-400 border border-zinc-700/50">
            {provider || 'Unspecified'}
          </span>
        );
    }
  };

  // Render status badge
  const renderStatusBadge = (status: string) => {
    switch (status) {
      case 'completed':
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-[11px] font-extrabold uppercase tracking-wide bg-emerald-950/60 text-emerald-300 border border-emerald-500/30">
            <CheckCircle2 size={12} className="text-emerald-400" />
            Completed
          </span>
        );
      case 'pending':
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-[11px] font-extrabold uppercase tracking-wide bg-amber-950/60 text-amber-300 border border-amber-500/30">
            <Clock size={12} className="text-amber-400" />
            Pending
          </span>
        );
      case 'cancelled':
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-[11px] font-extrabold uppercase tracking-wide bg-zinc-800 text-zinc-400 border border-zinc-700">
            <XCircle size={12} className="text-zinc-500" />
            Cancelled
          </span>
        );
      default:
        return (
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-bold uppercase bg-zinc-800 text-zinc-300">
            {status}
          </span>
        );
    }
  };

  // Render transaction type badge
  const renderTypeBadge = (tType: string) => {
    if (tType === 'listing_fee') {
      return (
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-extrabold uppercase tracking-wider bg-blue-950/50 text-blue-400 border border-blue-800/40">
          Listing Fee
        </span>
      );
    }
    if (tType === 'product_sale') {
      return (
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-extrabold uppercase tracking-wider bg-indigo-950/50 text-indigo-400 border border-indigo-800/40">
          Product Sale
        </span>
      );
    }
    return (
      <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-medium uppercase tracking-wider bg-zinc-800 text-zinc-400">
        {tType.replace('_', ' ')}
      </span>
    );
  };

  // Authorization check guard
  if (!canViewTransactions) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[460px] p-8 text-center bg-zinc-900/40 border border-zinc-800/80 rounded-3xl">
        <div className="p-4 bg-red-500/10 text-red-400 rounded-2xl border border-red-500/20 mb-4 shadow-xl">
          <ShieldAlert size={36} />
        </div>
        <h2 className="text-lg font-black text-white">Access Denied — Transactions Ledger Restricted</h2>
        <p className="text-xs text-zinc-400 max-w-md mt-1.5 leading-relaxed">
          Access to the platform-wide financial transactions ledger requires Platform Owner, Platform Administrator, Operations Administrator, or an explicit <code className="text-orange-400 font-mono">transactions</code> permission.
        </p>
        <div className="mt-5">
          <Link
            href="/admin/dashboard"
            className="px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-xs font-bold rounded-xl transition-colors"
          >
            Return to Dashboard
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Toast Feedback */}
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
            <AlertTriangle size={18} className="text-red-400 shrink-0" />
          )}
          <span>{toast.message}</span>
        </div>
      )}

      {/* Header Banner */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div className="flex items-center gap-3.5">
          <div className="p-3 bg-gradient-to-br from-emerald-600 to-teal-700 text-white rounded-2xl shadow-lg shadow-emerald-600/20">
            <Receipt size={26} />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-2xl font-black text-white tracking-tight">Financial Transactions Ledger</h1>
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider bg-emerald-600/10 text-emerald-400 border border-emerald-500/20">
                Immutable Audit Trail
              </span>
            </div>
            <p className="text-xs text-zinc-400 mt-0.5">
              Authoritative ledger of marketplace transactions, listing fee settlements, platform commissions, and mobile payment gateway reconciliation across Ethiopia.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={() => loadTransactions()}
            disabled={loading}
            className="flex items-center gap-2 px-3.5 py-2 bg-zinc-900 hover:bg-zinc-800 text-zinc-300 border border-zinc-800 rounded-xl text-xs font-bold transition-all disabled:opacity-50"
            title="Refresh transactions ledger"
          >
            <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
            <span>Refresh</span>
          </button>
        </div>
      </div>

      {/* Accurate Summary Metric Tiles (based strictly on verified API data) */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="p-4 bg-zinc-900/60 border border-zinc-800/80 rounded-2xl">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-bold text-zinc-500 uppercase tracking-wider">
              Total Records
            </span>
            <Receipt size={16} className="text-zinc-400" />
          </div>
          <div className="text-2xl font-black text-white mt-2">
            {loading ? '—' : totalCount.toLocaleString()}
          </div>
          <div className="text-[10px] text-zinc-500 mt-1 font-medium">Platform-wide audit entries</div>
        </div>

        <div className="p-4 bg-zinc-900/60 border border-zinc-800/80 rounded-2xl">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-bold text-emerald-400 uppercase tracking-wider">
              Page Gross Volume
            </span>
            <DollarSign size={16} className="text-emerald-400" />
          </div>
          <div className="text-2xl font-black text-emerald-400 mt-2 truncate">
            {loading ? '—' : fmtETB(pageMetrics.grossVolume)}
          </div>
          <div className="text-[10px] text-zinc-500 mt-1 font-medium">Current page gross total</div>
        </div>

        <div className="p-4 bg-zinc-900/60 border border-zinc-800/80 rounded-2xl">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-bold text-indigo-400 uppercase tracking-wider">
              Page Platform Commission
            </span>
            <Coins size={16} className="text-indigo-400" />
          </div>
          <div className="text-2xl font-black text-indigo-400 mt-2 truncate">
            {loading ? '—' : fmtETB(pageMetrics.brokerFees)}
          </div>
          <div className="text-[10px] text-zinc-500 mt-1 font-medium">HUNTIQ retained broker revenue</div>
        </div>

        <div className="p-4 bg-zinc-900/60 border border-zinc-800/80 rounded-2xl">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-bold text-amber-400 uppercase tracking-wider">
              Page Status Split
            </span>
            <CreditCard size={16} className="text-amber-400" />
          </div>
          <div className="flex items-baseline gap-2 mt-2">
            <span className="text-sm font-black text-emerald-400">{pageMetrics.completed} <span className="text-[10px] font-normal text-zinc-500">done</span></span>
            <span className="text-zinc-600">•</span>
            <span className="text-sm font-black text-amber-400">{pageMetrics.pending} <span className="text-[10px] font-normal text-zinc-500">pend</span></span>
            <span className="text-zinc-600">•</span>
            <span className="text-sm font-black text-zinc-400">{pageMetrics.cancelled} <span className="text-[10px] font-normal text-zinc-500">canc</span></span>
          </div>
          <div className="text-[10px] text-zinc-500 mt-1 font-medium">Across {transactions.length} rows loaded</div>
        </div>
      </div>

      {/* Search & Filter Controls */}
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
              placeholder="Search by TXN ID, payment reference, buyer name, phone..."
              className="w-full bg-zinc-950 border border-zinc-800 focus:border-emerald-500 text-zinc-200 placeholder-zinc-500 text-xs rounded-xl pl-10 pr-9 py-2.5 outline-none transition-colors"
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

          {/* Date Range Inputs */}
          <div className="flex items-center gap-2">
            <div className="flex items-center bg-zinc-950 border border-zinc-800 rounded-xl px-2.5 py-1.5 text-xs text-zinc-300">
              <span className="text-[10px] font-bold text-zinc-500 uppercase mr-1.5">From:</span>
              <input
                type="date"
                value={fromDate}
                onChange={(e) => {
                  setFromDate(e.target.value);
                  setPage(1);
                }}
                className="bg-transparent text-xs text-zinc-200 outline-none cursor-pointer"
              />
            </div>
            <div className="flex items-center bg-zinc-950 border border-zinc-800 rounded-xl px-2.5 py-1.5 text-xs text-zinc-300">
              <span className="text-[10px] font-bold text-zinc-500 uppercase mr-1.5">To:</span>
              <input
                type="date"
                value={toDate}
                onChange={(e) => {
                  setToDate(e.target.value);
                  setPage(1);
                }}
                className="bg-transparent text-xs text-zinc-200 outline-none cursor-pointer"
              />
            </div>
          </div>
        </div>

        {/* Filter Pills */}
        <div className="flex flex-wrap items-center justify-between gap-3 pt-1 border-t border-zinc-800/60">
          <div className="flex flex-wrap items-center gap-2">
            {/* Status Filters */}
            <div className="flex items-center bg-zinc-950 border border-zinc-800 rounded-xl p-1 gap-1">
              {(['all', 'completed', 'pending', 'cancelled'] as const).map((s) => (
                <button
                  key={s}
                  onClick={() => {
                    setStatusFilter(s);
                    setPage(1);
                  }}
                  className={`px-3 py-1 rounded-lg text-xs font-bold capitalize transition-all ${
                    statusFilter === s
                      ? 'bg-emerald-600 text-white shadow-sm'
                      : 'text-zinc-400 hover:text-zinc-200 hover:bg-zinc-900'
                  }`}
                >
                  {s === 'all' ? 'All Status' : s}
                </button>
              ))}
            </div>

            {/* Payment Method Filters */}
            <div className="flex items-center bg-zinc-950 border border-zinc-800 rounded-xl p-1 gap-1">
              {(['all', 'telebirr', 'cbe', 'safaricom', 'bank'] as const).map((m) => (
                <button
                  key={m}
                  onClick={() => {
                    setMethodFilter(m);
                    setPage(1);
                  }}
                  className={`px-3 py-1 rounded-lg text-xs font-bold capitalize transition-all ${
                    methodFilter === m
                      ? 'bg-indigo-600 text-white shadow-sm'
                      : 'text-zinc-400 hover:text-zinc-200 hover:bg-zinc-900'
                  }`}
                >
                  {m === 'all'
                    ? 'All Providers'
                    : m === 'cbe'
                    ? 'CBE'
                    : m === 'safaricom'
                    ? 'M-Pesa'
                    : m}
                </button>
              ))}
            </div>

            {/* Transaction Type Filters */}
            <div className="flex items-center bg-zinc-950 border border-zinc-800 rounded-xl p-1 gap-1">
              {(['all', 'listing_fee', 'product_sale'] as const).map((t) => (
                <button
                  key={t}
                  onClick={() => {
                    setTypeFilter(t);
                    setPage(1);
                  }}
                  className={`px-3 py-1 rounded-lg text-xs font-bold transition-all ${
                    typeFilter === t
                      ? 'bg-orange-600 text-white shadow-sm'
                      : 'text-zinc-400 hover:text-zinc-200 hover:bg-zinc-900'
                  }`}
                >
                  {t === 'all' ? 'All Types' : t === 'listing_fee' ? 'Listing Fees' : 'Product Sales'}
                </button>
              ))}
            </div>
          </div>

          {/* Reset Filters */}
          {isFiltered && (
            <button
              onClick={handleResetFilters}
              className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold text-zinc-400 hover:text-white bg-zinc-800/80 hover:bg-zinc-700/80 rounded-xl transition-colors"
            >
              <X size={13} />
              <span>Reset Filters</span>
            </button>
          )}
        </div>
      </div>

      {/* Error State */}
      {error && (
        <div className="p-4 bg-red-950/40 border border-red-500/30 rounded-2xl flex items-center justify-between gap-3 text-red-300">
          <div className="flex items-center gap-3">
            <AlertTriangle size={20} className="text-red-400 shrink-0" />
            <span className="text-xs font-medium">{error}</span>
          </div>
          <button
            onClick={() => loadTransactions()}
            className="px-3 py-1 bg-red-900/60 hover:bg-red-800 text-white text-xs font-bold rounded-lg transition-colors shrink-0"
          >
            Retry
          </button>
        </div>
      )}

      {/* Main Ledger Table */}
      <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl overflow-hidden shadow-xl">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-zinc-950/80 text-zinc-400 uppercase text-[10px] font-black tracking-wider border-b border-zinc-800/80">
              <tr>
                <th className="py-3 px-4">TXN ID & Date</th>
                <th className="py-3 px-4">Type</th>
                <th className="py-3 px-4">Parties & Contacts</th>
                <th className="py-3 px-4">Gross Amount</th>
                <th className="py-3 px-4">Platform Fee / Net</th>
                <th className="py-3 px-4">Payment Method & Ref</th>
                <th className="py-3 px-4">Status</th>
                <th className="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/40">
              {loading ? (
                // Loading Skeleton Rows
                Array.from({ length: 6 }).map((_, idx) => (
                  <tr key={idx} className="animate-pulse">
                    <td className="py-4 px-4"><div className="h-4 bg-zinc-800 rounded w-24 mb-1.5" /><div className="h-3 bg-zinc-800/60 rounded w-16" /></td>
                    <td className="py-4 px-4"><div className="h-5 bg-zinc-800 rounded-md w-20" /></td>
                    <td className="py-4 px-4"><div className="h-4 bg-zinc-800 rounded w-28 mb-1.5" /><div className="h-3 bg-zinc-800/60 rounded w-20" /></td>
                    <td className="py-4 px-4"><div className="h-4 bg-zinc-800 rounded w-20" /></td>
                    <td className="py-4 px-4"><div className="h-4 bg-zinc-800 rounded w-24 mb-1" /><div className="h-3 bg-zinc-800/60 rounded w-16" /></td>
                    <td className="py-4 px-4"><div className="h-5 bg-zinc-800 rounded-md w-24 mb-1" /><div className="h-3 bg-zinc-800/60 rounded w-28" /></td>
                    <td className="py-4 px-4"><div className="h-5 bg-zinc-800 rounded-lg w-20" /></td>
                    <td className="py-4 px-4 text-right"><div className="h-7 bg-zinc-800 rounded-lg w-16 ml-auto" /></td>
                  </tr>
                ))
              ) : transactions.length === 0 ? (
                // Empty State
                <tr>
                  <td colSpan={8} className="py-16 text-center text-zinc-500">
                    <div className="flex flex-col items-center justify-center gap-3">
                      <div className="p-4 bg-zinc-800/50 rounded-2xl text-zinc-400">
                        <Receipt size={36} />
                      </div>
                      <div className="text-sm font-bold text-zinc-300">No Transactions Found</div>
                      <p className="text-xs text-zinc-500 max-w-sm">
                        {isFiltered
                          ? 'No financial records match the specified search query or filters. Try adjusting your parameters.'
                          : 'No financial transaction records have been recorded on HUNTIQ yet.'}
                      </p>
                      {isFiltered && (
                        <button
                          onClick={handleResetFilters}
                          className="mt-2 px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-xs font-bold rounded-xl transition-colors"
                        >
                          Clear All Filters
                        </button>
                      )}
                    </div>
                  </td>
                </tr>
              ) : (
                // Data Rows
                transactions.map((tx) => (
                  <tr
                    key={tx.id}
                    onClick={() => setInspectTx(tx)}
                    className="hover:bg-zinc-800/30 transition-colors cursor-pointer group"
                  >
                    {/* ID & Date */}
                    <td className="py-3.5 px-4">
                      <div className="font-mono font-black text-white text-xs">#{tx.id}</div>
                      <div className="text-[10px] text-zinc-500 mt-0.5">{formatDate(tx.created_at)}</div>
                    </td>

                    {/* Type */}
                    <td className="py-3.5 px-4">
                      {renderTypeBadge(tx.transaction_type)}
                    </td>

                    {/* Parties */}
                    <td className="py-3.5 px-4">
                      <div className="font-bold text-zinc-200 truncate max-w-[160px]">
                        {tx.buyer_name || `Buyer (User #${tx.user_id})`}
                      </div>
                      <div className="text-[10px] text-zinc-500 font-mono flex items-center gap-1 mt-0.5">
                        {tx.buyer_phone || 'No phone recorded'}
                      </div>
                      {tx.seller_id > 0 && (
                        <div className="text-[9px] text-zinc-500 mt-0.5">
                          Merchant ID: <span className="text-zinc-400 font-mono">#{tx.seller_id}</span>
                        </div>
                      )}
                    </td>

                    {/* Gross Amount */}
                    <td className="py-3.5 px-4 font-mono font-black text-white text-xs">
                      {fmtETB(tx.amount)}
                    </td>

                    {/* Fees & Net */}
                    <td className="py-3.5 px-4 text-xs font-mono">
                      <div className="text-indigo-400 font-bold">
                        Fee: {fmtETB(tx.broker_fee)}
                      </div>
                      <div className="text-[10px] text-zinc-500 mt-0.5">
                        Net: {fmtETB(tx.seller_gets)}
                      </div>
                    </td>

                    {/* Payment Provider & Reference */}
                    <td className="py-3.5 px-4">
                      <div>{renderProviderBadge(tx.payment_method)}</div>
                      {tx.payment_reference ? (
                        <div className="mt-1 flex items-center gap-1.5 text-[10px] font-mono text-zinc-400 group/ref">
                          <span className="truncate max-w-[130px]">{tx.payment_reference}</span>
                          <button
                            type="button"
                            onClick={(e) => {
                              e.stopPropagation();
                              handleCopy(tx.payment_reference!, `row-${tx.id}`);
                            }}
                            className="text-zinc-500 hover:text-white transition-colors"
                            title="Copy reference code"
                          >
                            {copiedKey === `row-${tx.id}` ? (
                              <Check size={11} className="text-emerald-400" />
                            ) : (
                              <Copy size={11} />
                            )}
                          </button>
                        </div>
                      ) : (
                        <div className="text-[10px] text-zinc-600 mt-0.5 italic">No reference</div>
                      )}
                    </td>

                    {/* Status */}
                    <td className="py-3.5 px-4">
                      {renderStatusBadge(tx.status)}
                    </td>

                    {/* Actions */}
                    <td className="py-3.5 px-4 text-right">
                      <button
                        type="button"
                        onClick={(e) => {
                          e.stopPropagation();
                          setInspectTx(tx);
                        }}
                        className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 hover:text-white rounded-xl text-xs font-bold transition-all shadow-sm"
                      >
                        <Eye size={13} />
                        <span>Inspect</span>
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Server-Side Pagination Bar */}
        <div className="p-4 bg-zinc-950/80 border-t border-zinc-800/80 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs text-zinc-400">
          <div className="flex items-center gap-3">
            <span>
              Showing <span className="font-bold text-white">{transactions.length > 0 ? (page - 1) * limit + 1 : 0}</span> to{' '}
              <span className="font-bold text-white">{Math.min(page * limit, totalCount)}</span> of{' '}
              <span className="font-bold text-white">{totalCount}</span> transactions
            </span>

            {/* Page Size Selector */}
            <div className="flex items-center gap-1.5">
              <span className="text-[10px] text-zinc-500 uppercase font-bold">Per page:</span>
              <select
                value={limit}
                onChange={(e) => {
                  setLimit(Number(e.target.value));
                  setPage(1);
                }}
                className="bg-zinc-900 border border-zinc-800 rounded-lg px-2 py-1 text-xs text-zinc-200 outline-none cursor-pointer"
              >
                <option value={15}>15</option>
                <option value={25}>25</option>
                <option value={50}>50</option>
                <option value={100}>100</option>
              </select>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page <= 1 || loading}
              className="flex items-center gap-1 px-3 py-1.5 bg-zinc-900 hover:bg-zinc-800 disabled:opacity-40 disabled:hover:bg-zinc-900 border border-zinc-800 rounded-xl text-xs font-bold transition-all text-zinc-300"
            >
              <ChevronLeft size={14} />
              <span>Previous</span>
            </button>

            <span className="px-3 py-1 font-mono text-zinc-400">
              Page <span className="font-bold text-white">{page}</span> of{' '}
              <span className="font-bold text-white">{totalPages}</span>
            </span>

            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page >= totalPages || loading}
              className="flex items-center gap-1 px-3 py-1.5 bg-zinc-900 hover:bg-zinc-800 disabled:opacity-40 disabled:hover:bg-zinc-900 border border-zinc-800 rounded-xl text-xs font-bold transition-all text-zinc-300"
            >
              <span>Next</span>
              <ChevronRight size={14} />
            </button>
          </div>
        </div>
      </div>

      {/* Transaction Detail Slide-Over Drawer */}
      {inspectTx && (
        <div className="fixed inset-0 z-50 flex justify-end bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
          <div
            className="w-full max-w-xl bg-zinc-950 border-l border-zinc-800 h-full overflow-y-auto flex flex-col p-6 shadow-2xl animate-in slide-in-from-right duration-200"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Drawer Header */}
            <div className="flex items-center justify-between pb-4 border-b border-zinc-800">
              <div className="flex items-center gap-3">
                <div className="p-2.5 bg-emerald-600/10 text-emerald-400 rounded-xl border border-emerald-500/20">
                  <Receipt size={22} />
                </div>
                <div>
                  <div className="flex items-center gap-2">
                    <h2 className="text-base font-black text-white">Transaction Audit Record</h2>
                    <span className="font-mono text-xs text-zinc-500">#{inspectTx.id}</span>
                  </div>
                  <div className="text-[11px] text-zinc-400 mt-0.5">
                    Settled on {formatDate(inspectTx.created_at)}
                  </div>
                </div>
              </div>

              <button
                onClick={() => setInspectTx(null)}
                className="p-2 text-zinc-400 hover:text-white rounded-xl hover:bg-zinc-900 transition-colors"
                title="Close Drawer"
              >
                <X size={18} />
              </button>
            </div>

            {/* Drawer Content */}
            <div className="flex-1 py-6 space-y-6">
              {/* Status & Gateway Summary Card */}
              <div className="p-4 bg-zinc-900/60 border border-zinc-800 rounded-2xl space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-[11px] font-bold text-zinc-500 uppercase tracking-wider">
                    Settlement State
                  </span>
                  {renderStatusBadge(inspectTx.status)}
                </div>

                <div className="flex items-center justify-between pt-2 border-t border-zinc-800/60">
                  <span className="text-xs text-zinc-400">Payment Gateway:</span>
                  <div>{renderProviderBadge(inspectTx.payment_method)}</div>
                </div>

                <div className="flex items-center justify-between pt-2 border-t border-zinc-800/60">
                  <span className="text-xs text-zinc-400">Transaction Type:</span>
                  <div>{renderTypeBadge(inspectTx.transaction_type)}</div>
                </div>
              </div>

              {/* Financial Settlement Breakdown */}
              <div className="p-4 bg-zinc-900/60 border border-zinc-800 rounded-2xl space-y-3">
                <div className="text-[11px] font-bold text-zinc-500 uppercase tracking-wider">
                  Financial Settlement Breakdown
                </div>

                <div className="flex items-center justify-between text-xs">
                  <span className="text-zinc-400">Gross Settlement Amount:</span>
                  <span className="font-mono font-black text-white text-sm">
                    {fmtETB(inspectTx.amount)}
                  </span>
                </div>

                <div className="flex items-center justify-between text-xs pt-2 border-t border-zinc-800/60">
                  <span className="text-zinc-400">Platform Broker Fee:</span>
                  <span className="font-mono font-bold text-indigo-400">
                    {fmtETB(inspectTx.broker_fee)}
                  </span>
                </div>

                {inspectTx.fee_amount > 0 && (
                  <div className="flex items-center justify-between text-xs pt-2 border-t border-zinc-800/60">
                    <span className="text-zinc-400">Additional Fee Surcharge:</span>
                    <span className="font-mono text-zinc-300">
                      {fmtETB(inspectTx.fee_amount)}
                    </span>
                  </div>
                )}

                <div className="flex items-center justify-between text-xs pt-2 border-t border-zinc-800/60">
                  <span className="text-emerald-400 font-bold">Net Merchant Proceeds:</span>
                  <span className="font-mono font-black text-emerald-400 text-sm">
                    {fmtETB(inspectTx.seller_gets)}
                  </span>
                </div>
              </div>

              {/* Gateway References & Reconciliation */}
              <div className="p-4 bg-zinc-900/60 border border-zinc-800 rounded-2xl space-y-3">
                <div className="text-[11px] font-bold text-zinc-500 uppercase tracking-wider">
                  Payment Reference & Audit Codes
                </div>

                <div className="space-y-1">
                  <div className="text-[10px] text-zinc-500 uppercase font-bold">Payment Gateway Reference:</div>
                  <div className="flex items-center justify-between p-2.5 bg-zinc-950 border border-zinc-800 rounded-xl">
                    <span className="font-mono text-xs text-white select-all">
                      {inspectTx.payment_reference || 'N/A — No Reference Logged'}
                    </span>
                    {inspectTx.payment_reference && (
                      <button
                        onClick={() => handleCopy(inspectTx.payment_reference!, 'drawer-ref')}
                        className="flex items-center gap-1 text-[11px] text-zinc-400 hover:text-white px-2 py-1 rounded bg-zinc-900 transition-colors"
                      >
                        {copiedKey === 'drawer-ref' ? (
                          <>
                            <Check size={12} className="text-emerald-400" />
                            <span className="text-emerald-400">Copied</span>
                          </>
                        ) : (
                          <>
                            <Copy size={12} />
                            <span>Copy</span>
                          </>
                        )}
                      </button>
                    )}
                  </div>
                </div>

                {inspectTx.payer_phone && (
                  <div className="space-y-1 pt-2">
                    <div className="text-[10px] text-zinc-500 uppercase font-bold">Payer Mobile Contact:</div>
                    <div className="flex items-center justify-between p-2.5 bg-zinc-950 border border-zinc-800 rounded-xl">
                      <span className="font-mono text-xs text-zinc-200">
                        {inspectTx.payer_phone}
                      </span>
                      <button
                        onClick={() => handleCopy(inspectTx.payer_phone!, 'drawer-payer-phone')}
                        className="flex items-center gap-1 text-[11px] text-zinc-400 hover:text-white px-2 py-1 rounded bg-zinc-900 transition-colors"
                      >
                        {copiedKey === 'drawer-payer-phone' ? (
                          <>
                            <Check size={12} className="text-emerald-400" />
                            <span className="text-emerald-400">Copied</span>
                          </>
                        ) : (
                          <>
                            <Copy size={12} />
                            <span>Copy</span>
                          </>
                        )}
                      </button>
                    </div>
                  </div>
                )}
              </div>

              {/* Transaction Parties */}
              <div className="p-4 bg-zinc-900/60 border border-zinc-800 rounded-2xl space-y-3">
                <div className="text-[11px] font-bold text-zinc-500 uppercase tracking-wider">
                  Associated Parties & Entities
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div className="p-3 bg-zinc-950 border border-zinc-800/80 rounded-xl space-y-1">
                    <div className="text-[10px] text-zinc-500 uppercase font-bold">Buyer / Payer</div>
                    <div className="text-xs font-bold text-white truncate">
                      {inspectTx.buyer_name || `User #${inspectTx.user_id}`}
                    </div>
                    <div className="text-[10px] font-mono text-zinc-400">
                      {txPhone(inspectTx.buyer_phone)}
                    </div>
                    <div className="text-[9px] text-zinc-600 font-mono">User ID: #{inspectTx.user_id}</div>
                  </div>

                  <div className="p-3 bg-zinc-950 border border-zinc-800/80 rounded-xl space-y-1">
                    <div className="text-[10px] text-zinc-500 uppercase font-bold">Merchant / Seller</div>
                    <div className="text-xs font-bold text-white truncate">
                      {inspectTx.seller_id > 0 ? `Seller #${inspectTx.seller_id}` : 'Platform Direct'}
                    </div>
                    {inspectTx.seller_id > 0 && (
                      <Link
                        href={`/admin/sellers`}
                        className="inline-flex items-center gap-1 text-[10px] text-orange-400 hover:underline font-bold mt-1"
                      >
                        <span>Inspect in Directory</span>
                        <ExternalLink size={10} />
                      </Link>
                    )}
                  </div>
                </div>

                {inspectTx.product_id > 0 && (
                  <div className="flex items-center justify-between p-3 bg-zinc-950 border border-zinc-800/80 rounded-xl text-xs">
                    <div>
                      <div className="text-[10px] text-zinc-500 uppercase font-bold">Product / Listing</div>
                      <div className="text-zinc-200 font-mono font-bold mt-0.5">Listing Item #{inspectTx.product_id}</div>
                    </div>
                    <Link
                      href={`/admin/products`}
                      className="inline-flex items-center gap-1 px-2.5 py-1 bg-zinc-900 hover:bg-zinc-800 text-zinc-300 text-[11px] font-bold rounded-lg border border-zinc-800 transition-colors"
                    >
                      <span>Find Product</span>
                      <ExternalLink size={11} />
                    </Link>
                  </div>
                )}
              </div>

              {/* Operational Audit Notes */}
              {inspectTx.notes && (
                <div className="p-4 bg-zinc-900/60 border border-zinc-800 rounded-2xl space-y-1.5">
                  <div className="text-[11px] font-bold text-zinc-500 uppercase tracking-wider">
                    Ledger Audit Notes
                  </div>
                  <p className="text-xs text-zinc-300 bg-zinc-950 p-3 rounded-xl border border-zinc-800/80 leading-relaxed font-mono">
                    {inspectTx.notes}
                  </p>
                </div>
              )}

              {/* Immutability & Financial Compliance Notice */}
              <div className="p-3.5 bg-emerald-950/20 border border-emerald-500/20 rounded-xl flex items-start gap-2.5 text-zinc-400 text-xs">
                <Lock size={16} className="text-emerald-400 shrink-0 mt-0.5" />
                <div className="leading-relaxed text-[11px]">
                  <strong className="text-emerald-300">Immutable Financial Audit Standard:</strong> This transaction record represents an authoritative journal entry. Manual modifications and row deletions are cryptographically prohibited under Ethiopian commercial accounting controls.
                </div>
              </div>
            </div>

            {/* Drawer Footer */}
            <div className="pt-4 border-t border-zinc-800 flex items-center justify-between">
              <span className="text-[10px] font-mono text-zinc-500">
                Last modified: {formatDate(inspectTx.updated_at)}
              </span>
              <button
                onClick={() => setInspectTx(null)}
                className="px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-xs font-bold rounded-xl transition-colors"
              >
                Close Ledger Details
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

function txPhone(phone?: string) {
  if (!phone) return '—';
  return phone;
}
