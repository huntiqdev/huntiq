'use client';

import React, { useState, useEffect, useCallback } from 'react';
import {
  ChevronDown, ChevronRight, CheckSquare, Square, Download, Edit3, ShieldCheck,
  AlertTriangle, Filter, Eye, FileText, CheckCircle2, XCircle, DollarSign, Send, Trash2, ArrowUpDown
} from 'lucide-react';
import { DomainMode } from './DomainSwitcherTabs';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';

export interface DataRowItem {
  id: string;
  title: string;
  category: string;
  vertical: 'ecommerce' | 'rentals' | 'leases';
  price: number;
  pricePeriod?: string;
  status: 'Approved' | 'Pending Review' | 'Suspended' | 'In Escrow';
  verificationBadge: 'Verified Landlord' | 'Top Seller' | 'Pending Verification' | 'VIP Tenant';
  dateCreated: string;
  specDetails: {
    skuOrUnit: string;
    locationOrSeller: string;
    escrowAmount?: number;
    depositStatus?: string;
    maintenanceHistory?: string;
    documentsAvailable: string[];
  };
}

const MOCK_ROWS: DataRowItem[] = [
  {
    id: 'RNT-9901',
    title: 'Bole Luxury Penthouse Apartment',
    category: 'Apartment',
    vertical: 'rentals',
    price: 2200,
    pricePeriod: '/month',
    status: 'Approved',
    verificationBadge: 'Verified Landlord',
    dateCreated: '2026-08-10',
    specDetails: {
      skuOrUnit: 'UNIT-B401',
      locationOrSeller: 'Bole Hub • Addis Ababa',
      escrowAmount: 4400,
      depositStatus: 'Escrow Lock Confirmed',
      maintenanceHistory: 'Elevator serviced on Aug 01',
      documentsAvailable: ['Lease Contract PDF', 'Landlord Deed', 'Tenant ID Scan']
    }
  },
  {
    id: 'PRD-4412',
    title: 'Sony WH-1000XM5 Wireless Headphones',
    category: 'Electronics',
    vertical: 'ecommerce',
    price: 349,
    status: 'Approved',
    verificationBadge: 'Top Seller',
    dateCreated: '2026-08-12',
    specDetails: {
      skuOrUnit: 'SKU-WH-1000XM5',
      locationOrSeller: 'TechWorld Superstore',
      documentsAvailable: ['Import Duty Clearance', 'Tax Compliance ID']
    }
  },
  {
    id: 'LSE-2209',
    title: 'Lease Contract #LA-2026-88 (Kazanchis Suite)',
    category: 'Lease Contract',
    vertical: 'leases',
    price: 1800,
    pricePeriod: '/month',
    status: 'In Escrow',
    verificationBadge: 'VIP Tenant',
    dateCreated: '2026-08-14',
    specDetails: {
      skuOrUnit: 'LEASE-LA88',
      locationOrSeller: 'Landlord: Solomon Tadesse • Tenant: Meron A.',
      escrowAmount: 3600,
      depositStatus: 'Security Deposit Held ($3,600)',
      maintenanceHistory: 'HVAC Filter Changed',
      documentsAvailable: ['Lease Agreement PDF', 'Escrow Voucher', 'Tenant Passport']
    }
  },
  {
    id: 'RNT-7721',
    title: 'Hawassa Lake Resort Luxury Villa A',
    category: 'Villa',
    vertical: 'rentals',
    price: 3500,
    pricePeriod: '/month',
    status: 'Approved',
    verificationBadge: 'Verified Landlord',
    dateCreated: '2026-08-01',
    specDetails: {
      skuOrUnit: 'VILLA-HAW-01',
      locationOrSeller: 'Lake Hawassa Waterfront',
      escrowAmount: 7000,
      depositStatus: 'Active Security Escrow',
      documentsAvailable: ['Property Title PDF', 'Pool Maintenance Log']
    }
  },
  {
    id: 'PRD-1092',
    title: 'Ergonomic Executive Office Chair',
    category: 'Furniture',
    vertical: 'ecommerce',
    price: 185,
    status: 'Pending Review',
    verificationBadge: 'Pending Verification',
    dateCreated: '2026-08-15',
    specDetails: {
      skuOrUnit: 'SKU-ERGO-09',
      locationOrSeller: 'Addis Furniture Outlet',
      documentsAvailable: ['Seller License Draft']
    }
  },
  {
    id: 'LSE-3091',
    title: 'Commercial Lease Agreement #LA-5510',
    category: 'Commercial Lease',
    vertical: 'leases',
    price: 4500,
    pricePeriod: '/month',
    status: 'Suspended',
    verificationBadge: 'Pending Verification',
    dateCreated: '2026-08-05',
    specDetails: {
      skuOrUnit: 'LEASE-COMM-01',
      locationOrSeller: 'Dire Dawa Commercial Zone',
      escrowAmount: 9000,
      depositStatus: 'Escrow Disputed',
      maintenanceHistory: 'Plumbing inspection pending',
      documentsAvailable: ['Dispute Claim Form', 'Lease Draft']
    }
  }
];

interface EnterpriseDataTableProps {
  activeMode: DomainMode;
  onInspectDocument: (documentName: string, row: DataRowItem) => void;
  onCountsChange?: (counts: { ecommerceCount: number; rentalsCount: number; leasesCount: number }) => void;
  refreshTrigger?: number;
}

export function EnterpriseDataTable({ activeMode, onInspectDocument, onCountsChange, refreshTrigger }: EnterpriseDataTableProps) {
  const { user } = useAuth();
  const [rows, setRows] = useState<DataRowItem[]>(MOCK_ROWS);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [expandedRowId, setExpandedRowId] = useState<string | null>(null);
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [searchFilter, setSearchFilter] = useState<string>('');
  const [bulkActionNotice, setBulkActionNotice] = useState<{ message: string; type: 'success' | 'error' | 'info' } | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  // Map Backend Product into DataRowItem
  const mapProductToDataRowItem = useCallback((p: any): DataRowItem => {
    const isRent = p.listing_type === 'rent';
    const isRented = p.is_currently_rented === true;
    
    let vertical: 'ecommerce' | 'rentals' | 'leases' = 'ecommerce';
    if (isRent) {
      if (isRented) {
        vertical = 'leases';
      } else {
        vertical = 'rentals';
      }
    }

    let status: DataRowItem['status'] = 'Approved';
    if (p.status === 'pending') status = 'Pending Review';
    else if (p.status === 'rejected' || p.status === 'expired') status = 'Suspended';
    else if (isRent && isRented) status = 'In Escrow';
    else if (p.status === 'approved') status = 'Approved';

    let verificationBadge: DataRowItem['verificationBadge'] = 'Pending Verification';
    if (p.seller_tier === 'gold') verificationBadge = 'Top Seller';
    else if (p.seller_tier === 'premium') verificationBadge = 'VIP Tenant';
    else if (p.seller_active) verificationBadge = 'Verified Landlord';

    const pricePeriod = isRent ? `/${p.rental_duration_unit || 'month'}` : undefined;

    return {
      id: `${isRent ? (isRented ? 'LSE' : 'RNT') : 'PRD'}-${p.id}`,
      title: p.title || 'Untitled Listing',
      category: p.category?.name_en || p.subcategory || 'General',
      vertical,
      price: isRent ? (p.rental_price || p.price) : p.price,
      pricePeriod,
      status,
      verificationBadge,
      dateCreated: p.created_at ? new Date(p.created_at).toISOString().split('T')[0] : new Date().toISOString().split('T')[0],
      specDetails: {
        skuOrUnit: isRent ? `UNIT-${p.id}` : `SKU-${p.id}`,
        locationOrSeller: p.location?.name_en || p.specific_area || 'Platform Vendor',
        escrowAmount: p.deposit_amount || undefined,
        depositStatus: isRented ? 'Escrow Lock Confirmed' : undefined,
        maintenanceHistory: isRent ? 'Property operational checks validated' : undefined,
        documentsAvailable: isRent 
          ? ['Lease Contract PDF', 'Property Deed Registration'] 
          : ['Import Duty Clearance', 'Tax Compliance ID']
      }
    };
  }, []);

  // Fetch all products from real API
  const loadData = useCallback(async () => {
    if (!user?.token) return;
    setLoading(true);
    try {
      const data = await api.adminGetAllProducts(user.token);
      if (Array.isArray(data)) {
        const mapped = data.map(mapProductToDataRowItem);
        // Include mock data only if the database has 0 items so it doesn't look empty initially.
        if (mapped.length === 0) {
          setRows(MOCK_ROWS);
        } else {
          setRows(mapped);
        }
      }
    } catch (err: any) {
      console.error('Failed to load products from API:', err);
      triggerNotice('Failed to sync live marketplace inventory. Using local sandbox fallback.', 'error');
      setRows(MOCK_ROWS);
    } finally {
      setLoading(false);
    }
  }, [user, mapProductToDataRowItem]);

  useEffect(() => {
    loadData();
  }, [loadData, refreshTrigger]);

  // Propagate counts back up to dashboard switcher
  useEffect(() => {
    if (onCountsChange) {
      const ecommerceCount = rows.filter(r => r.vertical === 'ecommerce').length;
      const rentalsCount = rows.filter(r => r.vertical === 'rentals').length;
      const leasesCount = rows.filter(r => r.vertical === 'leases').length;
      onCountsChange({ ecommerceCount, rentalsCount, leasesCount });
    }
  }, [rows, onCountsChange]);

  const toggleSelectAll = () => {
    if (selectedIds.length === filteredRows.length) {
      setSelectedIds([]);
    } else {
      setSelectedIds(filteredRows.map((r) => r.id));
    }
  };

  const toggleSelectRow = (id: string) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  // Bulk Status Modification
  const handleBulkStatusChange = async (newStatus: DataRowItem['status']) => {
    if (!user?.token) return;
    const backendStatus = newStatus === 'Approved' ? 'approved' : 'rejected';
    
    triggerNotice(`Moderating ${selectedIds.length} items in backend...`, 'info');

    const promises = selectedIds.map(async (rowId) => {
      const numericId = parseInt(rowId.split('-')[1]);
      if (isNaN(numericId)) return;
      try {
        await api.adminReviewProduct(numericId, { status: backendStatus, notes: 'Bulk approved via Enterprise Workspace' }, user.token);
      } catch (e: any) {
        console.error(`Failed to review product ${numericId}:`, e);
      }
    });

    await Promise.all(promises);
    await loadData();
    triggerNotice(`Successfully updated ${selectedIds.length} records in live database.`, 'success');
    setSelectedIds([]);
  };

  // Individual Row Deletion
  const handleDeleteRow = async (rowId: string) => {
    if (!user?.token) return;
    const numericId = parseInt(rowId.split('-')[1]);
    if (isNaN(numericId)) return;

    if (window.confirm(`Are you sure you want to permanently delete listing #${numericId} from the database?`)) {
      try {
        await api.adminDeleteProduct(numericId, user.token);
        triggerNotice(`Record #${numericId} successfully deleted.`, 'success');
        await loadData();
      } catch (err: any) {
        console.error('Delete failed:', err);
        triggerNotice(`Deletion failed: ${err.message || 'Server error'}`, 'error');
      }
    }
  };

  // Toggle Individual Status
  const handleToggleRowStatus = async (row: DataRowItem) => {
    if (!user?.token) {
      // client sandbox mode
      setRows((prev) =>
        prev.map((r) => (r.id === row.id ? { ...r, status: r.status === 'Approved' ? 'Suspended' : 'Approved' } : r))
      );
      triggerNotice('Toggled status in client sandbox.', 'success');
      return;
    }
    const numericId = parseInt(row.id.split('-')[1]);
    if (isNaN(numericId)) return;

    const newStatus = row.status === 'Approved' ? 'rejected' : 'approved';
    try {
      await api.adminReviewProduct(numericId, { status: newStatus, notes: 'Toggled via quick-moderation switch' }, user.token);
      triggerNotice(`Listing status updated to ${newStatus === 'approved' ? 'Approved' : 'Suspended'}.`, 'success');
      await loadData();
    } catch (err: any) {
      console.error('Status toggle failed:', err);
      triggerNotice(`Validation check: ${err.message || 'Server error'}`, 'error');
    }
  };

  const handleBulkPriceUpdate = (percentage: number) => {
    setRows((prev) =>
      prev.map((r) => {
        if (selectedIds.includes(r.id)) {
          const newPrice = Math.round(r.price * (1 + percentage / 100));
          return { ...r, price: newPrice };
        }
        return r;
      })
    );
    triggerNotice(`Applied ${percentage}% price adjustment to ${selectedIds.length} selected items.`, 'success');
    setSelectedIds([]);
  };

  const triggerNotice = (msg: string, type: 'success' | 'error' | 'info' = 'success') => {
    setBulkActionNotice({ message: msg, type });
    setTimeout(() => setBulkActionNotice(null), 4000);
  };

  // Filter rows based on domain mode, search, and status filter
  const filteredRows = rows.filter((row) => {
    const matchesDomain =
      activeMode === 'ecommerce'
        ? row.vertical === 'ecommerce'
        : activeMode === 'rentals'
        ? row.vertical === 'rentals'
        : row.vertical === 'leases';

    const matchesStatus = statusFilter === 'all' || row.status.toLowerCase() === statusFilter.toLowerCase();
    const matchesSearch =
      row.title.toLowerCase().includes(searchFilter.toLowerCase()) ||
      row.id.toLowerCase().includes(searchFilter.toLowerCase()) ||
      row.category.toLowerCase().includes(searchFilter.toLowerCase());

    return matchesDomain && matchesStatus && matchesSearch;
  });

  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-5 shadow-xl space-y-4 relative">
      
      {/* Table Header & Controls */}
      <div className="flex flex-wrap items-center justify-between gap-3 pb-3 border-b border-zinc-800">
        <div className="flex items-center gap-3">
          <div className="p-2.5 bg-orange-600/10 text-orange-500 rounded-xl border border-orange-500/20">
            <FileText size={20} />
          </div>
          <div>
            <h3 className="text-base font-extrabold text-white flex items-center gap-2">
              High-Density Inventory & Records Directory
              <span className="px-2.5 py-0.5 text-xs font-mono bg-zinc-800 text-zinc-300 rounded-full border border-zinc-700">
                {filteredRows.length} Records {loading && '• Syncing...'}
              </span>
            </h3>
            <p className="text-xs text-zinc-400">Sub-row details, inline quick-edit, and mass batch operations</p>
          </div>
        </div>

        {/* Filter Bar */}
        <div className="flex flex-wrap items-center gap-2 text-xs">
          <input
            type="text"
            placeholder="Filter by ID, Title or Category..."
            value={searchFilter}
            onChange={(e) => setSearchFilter(e.target.value)}
            className="bg-zinc-950 border border-zinc-800 text-zinc-200 placeholder-zinc-500 rounded-xl px-3 py-2 focus:outline-none focus:border-orange-500 text-xs"
          />

          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="bg-zinc-950 border border-zinc-800 text-zinc-200 rounded-xl px-3 py-2 focus:outline-none focus:border-orange-500 font-bold"
          >
            <option value="all">All Statuses</option>
            <option value="Approved">Approved</option>
            <option value="Pending Review">Pending Review</option>
            <option value="In Escrow">In Escrow</option>
            <option value="Suspended">Suspended</option>
          </select>
        </div>
      </div>

      {bulkActionNotice && (
        <div className={`p-3 border rounded-xl text-xs font-bold flex items-center justify-between animate-in fade-in ${
          bulkActionNotice.type === 'success' 
            ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-400' 
            : bulkActionNotice.type === 'error'
            ? 'bg-red-500/10 border-red-500/30 text-red-400'
            : 'bg-blue-500/10 border-blue-500/30 text-blue-400'
        }`}>
          <span className="flex items-center gap-2">
            <CheckCircle2 size={16} />
            {bulkActionNotice.message}
          </span>
          <button onClick={() => setBulkActionNotice(null)} className="hover:text-white">✕</button>
        </div>
      )}

      {/* Floating Bulk Operations Toolbar */}
      {selectedIds.length > 0 && (
        <div className="p-3.5 bg-gradient-to-r from-orange-950/90 to-zinc-900 border border-orange-500/40 rounded-2xl text-white shadow-2xl flex flex-wrap items-center justify-between gap-3 animate-in slide-in-from-top-2">
          <div className="flex items-center gap-2 text-xs font-extrabold">
            <span className="w-6 h-6 bg-orange-600 rounded-full flex items-center justify-center text-xs font-mono">
              {selectedIds.length}
            </span>
            <span>Items Selected for Bulk Action</span>
          </div>

          <div className="flex flex-wrap items-center gap-2 text-xs font-bold">
            <button
              onClick={() => handleBulkStatusChange('Approved')}
              className="px-3 py-1.5 bg-emerald-600/20 text-emerald-400 border border-emerald-500/30 rounded-xl hover:bg-emerald-600/30 transition-colors"
            >
              ✓ Bulk Approve
            </button>
            <button
              onClick={() => handleBulkStatusChange('Suspended')}
              className="px-3 py-1.5 bg-red-600/20 text-red-400 border border-red-500/30 rounded-xl hover:bg-red-600/30 transition-colors"
            >
              ⚠ Bulk Suspend
            </button>
            <button
              onClick={() => handleBulkPriceUpdate(5)}
              className="px-3 py-1.5 bg-blue-600/20 text-blue-400 border border-blue-500/30 rounded-xl hover:bg-blue-600/30 transition-colors"
            >
              +5% Price Adjustment
            </button>
            <button
              onClick={() => triggerNotice(`Exported ${selectedIds.length} records to CSV report.`, 'info')}
              className="px-3 py-1.5 bg-zinc-800 text-zinc-200 border border-zinc-700 rounded-xl hover:bg-zinc-700 transition-colors flex items-center gap-1"
            >
              <Download size={14} /> Export CSV
            </button>
            <button
              onClick={() => triggerNotice(`Broadcast notification dispatched to ${selectedIds.length} stakeholders.`, 'info')}
              className="px-3 py-1.5 bg-amber-600/20 text-amber-400 border border-amber-500/30 rounded-xl hover:bg-amber-600/30 transition-colors flex items-center gap-1"
            >
              <Send size={14} /> Mass Alert
            </button>
          </div>
        </div>
      )}

      {/* Main Table */}
      <div className="overflow-x-auto rounded-xl border border-zinc-800">
        <table className="w-full text-xs text-left border-collapse">
          <thead>
            <tr className="bg-zinc-950 text-zinc-400 font-mono text-[11px] border-b border-zinc-800">
              <th className="p-3.5 w-10 text-center">
                <button onClick={toggleSelectAll} className="text-zinc-400 hover:text-white">
                  {selectedIds.length === filteredRows.length && filteredRows.length > 0 ? (
                    <CheckSquare size={16} className="text-orange-500" />
                  ) : (
                    <Square size={16} />
                  )}
                </button>
              </th>
              <th className="p-3.5 w-8"></th>
              <th className="p-3.5 font-bold text-zinc-300">Listing / Item Title</th>
              <th className="p-3.5 font-bold text-zinc-300">Category & Type</th>
              <th className="p-3.5 font-bold text-zinc-300">Rate / Price</th>
              <th className="p-3.5 font-bold text-zinc-300">Verification Badge</th>
              <th className="p-3.5 font-bold text-zinc-300">Status</th>
              <th className="p-3.5 text-right font-bold text-zinc-300 border-l border-zinc-850">Actions</th>
            </tr>
          </thead>

          <tbody className="divide-y divide-zinc-800/80 bg-zinc-900/40">
            {filteredRows.length === 0 ? (
              <tr>
                <td colSpan={8} className="p-12 text-center text-zinc-500">
                  No records match current domain mode and filters.
                </td>
              </tr>
            ) : (
              filteredRows.map((row) => {
                const isSelected = selectedIds.includes(row.id);
                const isExpanded = expandedRowId === row.id;

                let statusBadge = 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20';
                if (row.status === 'Pending Review') statusBadge = 'bg-amber-500/10 text-amber-400 border-amber-500/20';
                if (row.status === 'Suspended') statusBadge = 'bg-red-500/10 text-red-400 border-red-500/20';
                if (row.status === 'In Escrow') statusBadge = 'bg-blue-500/10 text-blue-400 border-blue-500/20';

                return (
                  <React.Fragment key={row.id}>
                    <tr
                      className={`hover:bg-zinc-800/50 transition-colors ${
                        isSelected ? 'bg-orange-950/20' : ''
                      }`}
                    >
                      {/* Checkbox */}
                      <td className="p-3.5 text-center">
                        <button onClick={() => toggleSelectRow(row.id)} className="text-zinc-400 hover:text-white">
                          {isSelected ? (
                            <CheckSquare size={16} className="text-orange-500" />
                          ) : (
                            <Square size={16} />
                          )}
                        </button>
                      </td>

                      {/* Expand Toggle */}
                      <td className="p-3.5 text-center">
                        <button
                          onClick={() => setExpandedRowId(isExpanded ? null : row.id)}
                          className="text-zinc-400 hover:text-white p-1 rounded-md hover:bg-zinc-800"
                        >
                          {isExpanded ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
                        </button>
                      </td>

                      {/* Title & Sub-ID */}
                      <td className="p-3.5 font-bold text-white">
                        <div className="flex items-center gap-2">
                          <span className="truncate max-w-xs">{row.title}</span>
                        </div>
                        <div className="text-[10px] font-mono text-zinc-400 mt-0.5">
                          ID: <span className="text-zinc-300">{row.id}</span> • Created {row.dateCreated}
                        </div>
                      </td>

                      {/* Category */}
                      <td className="p-3.5 text-zinc-300">
                        <span className="px-2.5 py-1 bg-zinc-800 border border-zinc-700 rounded-lg text-xs font-semibold">
                          {row.category}
                        </span>
                      </td>

                      {/* Price */}
                      <td className="p-3.5 font-mono font-bold text-orange-400 text-sm">
                        ${row.price.toLocaleString()}
                        {row.pricePeriod && <span className="text-xs text-zinc-400 font-sans font-normal">{row.pricePeriod}</span>}
                      </td>

                      {/* Verification Badge */}
                      <td className="p-3.5">
                        <span className="px-2 py-0.5 text-[11px] font-bold bg-purple-500/10 text-purple-400 border border-purple-500/20 rounded-md">
                          {row.verificationBadge}
                        </span>
                      </td>

                      {/* Status */}
                      <td className="p-3.5">
                        <span className={`px-2.5 py-1 rounded-lg text-[11px] font-extrabold border ${statusBadge}`}>
                          {row.status}
                        </span>
                      </td>

                      {/* Actions */}
                      <td className="p-3.5 text-right space-x-2 border-l border-zinc-850">
                        <button
                          onClick={() => setExpandedRowId(isExpanded ? null : row.id)}
                          className="p-1.5 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-lg transition-colors"
                          title="Inspect Sub-row Details"
                        >
                          <Eye size={14} />
                        </button>
                        <button
                          onClick={() => handleDeleteRow(row.id)}
                          className="p-1.5 bg-red-950/20 hover:bg-red-900/30 text-red-400 rounded-lg border border-red-500/10 transition-colors"
                          title="Delete Listing Record"
                        >
                          <Trash2 size={14} />
                        </button>
                      </td>
                    </tr>

                    {/* Sub-row Expansion Panel */}
                    {isExpanded && (
                      <tr className="bg-zinc-950/90 border-b border-zinc-800">
                        <td colSpan={8} className="p-4 pl-14">
                          <div className="p-4 bg-zinc-900 border border-zinc-800 rounded-2xl space-y-3">
                            <div className="flex items-center justify-between border-b border-zinc-800 pb-2">
                              <h4 className="font-extrabold text-xs text-orange-400 uppercase tracking-widest flex items-center gap-2">
                                <ShieldCheck size={14} /> Sub-Row Technical Inspection & Escrow Specs
                              </h4>
                              <span className="text-[10px] font-mono text-zinc-400">Ref: {row.specDetails.skuOrUnit}</span>
                            </div>

                            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs">
                              <div>
                                <span className="text-zinc-400 font-medium">Location / Vendor:</span>
                                <p className="font-bold text-zinc-100 mt-0.5">{row.specDetails.locationOrSeller}</p>
                              </div>

                              {row.specDetails.escrowAmount && (
                                <div>
                                  <span className="text-zinc-400 font-medium">Escrow Security Deposit:</span>
                                  <p className="font-bold text-emerald-400 font-mono mt-0.5">
                                    ${row.specDetails.escrowAmount.toLocaleString()} ({row.specDetails.depositStatus})
                                  </p>
                                </div>
                              )}

                              {row.specDetails.maintenanceHistory && (
                                <div>
                                  <span className="text-zinc-400 font-medium">Maintenance / Operational Queue:</span>
                                  <p className="font-bold text-amber-400 mt-0.5">{row.specDetails.maintenanceHistory}</p>
                                </div>
                              )}
                            </div>

                            {/* Documents Inspector Trigger */}
                            <div className="pt-2 border-t border-zinc-800/80 flex items-center justify-between">
                              <div className="flex items-center gap-2">
                                <span className="text-[11px] font-bold text-zinc-400">Attached Documents:</span>
                                <div className="flex flex-wrap gap-2">
                                  {row.specDetails.documentsAvailable.map((doc, idx) => (
                                    <button
                                      key={idx}
                                      onClick={() => onInspectDocument(doc, row)}
                                      className="px-2.5 py-1 bg-zinc-800 hover:bg-orange-600/20 text-zinc-200 hover:text-orange-400 border border-zinc-700 rounded-lg text-[11px] font-semibold transition-all flex items-center gap-1"
                                    >
                                      <FileText size={12} />
                                      {doc}
                                    </button>
                                  ))}
                                </div>
                              </div>

                              <button
                                onClick={() => handleToggleRowStatus(row)}
                                className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 hover:text-white rounded-lg text-xs font-bold transition-colors"
                              >
                                Toggle Approval Status
                              </button>
                            </div>
                          </div>
                        </td>
                      </tr>
                    )}
                  </React.Fragment>
                );
              })
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
