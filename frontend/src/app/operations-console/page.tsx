'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import {
  Activity, CheckCircle2, Clock, AlertTriangle, ShieldCheck,
  FileCheck, Users, Search, RefreshCw, Filter, ArrowUpRight,
  ClipboardList, Package, ExternalLink, ShieldAlert, XCircle,
  ThumbsUp, ThumbsDown, Check, X, Shield, Eye
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';

export default function OperationsConsolePage() {
  const router = useRouter();
  const { user, isPlatformOwner, isSuperAdmin, adminRole, loading: authLoading, logout } = useAuth();
  
  const [activeTab, setActiveTab] = useState<'all' | 'inspections' | 'kyc' | 'disputes' | 'appeals'>('all');
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [actionMsg, setActionMsg] = useState<{ text: string; type: 'success' | 'error' } | null>(null);

  // Live Database Driven State
  const [pendingProducts, setPendingProducts] = useState<any[]>([]);
  const [kycDocs, setKycDocs] = useState<any[]>([]);
  const [refunds, setRefunds] = useState<any[]>([]);
  const [appeals, setAppeals] = useState<any[]>([]);
  const [loadingData, setLoadingData] = useState(false);

  useEffect(() => {
    if (!authLoading) {
      const allowedRoles = [
        'platform_owner', 'super', 'operations_admin',
        'kyc_verification_officer', 'product_inspection_officer',
        'dispute_refund_officer', 'appeals_officer'
      ];
      const isAllowed = isPlatformOwner || isSuperAdmin || (adminRole && allowedRoles.includes(adminRole));
      if (!user || !isAllowed) {
        router.replace('/operations-console/login');
      }
    }
  }, [user, authLoading, isPlatformOwner, isSuperAdmin, adminRole, router]);

  const fetchLiveOperationsData = async () => {
    if (!user?.token) return;
    try {
      setLoadingData(true);
      const [prods, docs, refList, appList] = await Promise.all([
        api.adminGetPending(user.token).catch(() => []),
        api.adminGetVerificationDocuments('pending', user.token).then(r => r?.data || r || []).catch(() => []),
        api.adminGetRefundRequests('pending', user.token).then(r => r?.data || r || []).catch(() => []),
        api.adminGetAccountAppeals('pending', user.token).then(r => r?.data || r || []).catch(() => []),
      ]);

      setPendingProducts(Array.isArray(prods) ? prods : []);
      setKycDocs(Array.isArray(docs) ? docs : []);
      setRefunds(Array.isArray(refList) ? refList : []);
      setAppeals(Array.isArray(appList) ? appList : []);
    } catch {
      // Fallback
    } finally {
      setLoadingData(false);
      setIsRefreshing(false);
    }
  };

  useEffect(() => {
    fetchLiveOperationsData();
  }, [user?.token]);

  const handleRefresh = () => {
    setIsRefreshing(true);
    fetchLiveOperationsData();
  };

  // Action handlers
  const handleApproveProduct = async (id: number) => {
    if (!user?.token) return;
    try {
      await api.adminReviewProduct(id, { status: 'active', notes: 'Inspected and approved by Operations Officer' }, user.token);
      setActionMsg({ text: `Listing #${id} inspected and published successfully.`, type: 'success' });
      fetchLiveOperationsData();
    } catch (err: any) {
      setActionMsg({ text: err?.message || 'Failed to approve product', type: 'error' });
    }
  };

  const handleRejectProduct = async (id: number) => {
    if (!user?.token) return;
    const reason = prompt('Enter rejection reason for listing:', 'Does not meet specification standards');
    if (!reason) return;
    try {
      await api.adminReviewProduct(id, { status: 'rejected', notes: reason }, user.token);
      setActionMsg({ text: `Listing #${id} rejected with reason: ${reason}`, type: 'success' });
      fetchLiveOperationsData();
    } catch (err: any) {
      setActionMsg({ text: err?.message || 'Failed to reject product', type: 'error' });
    }
  };

  const handleReviewKYC = async (id: number, approve: boolean) => {
    if (!user?.token) return;
    const status = approve ? 'approved' : 'rejected';
    const reason = approve ? '' : (prompt('Enter reason for rejection:') || 'Invalid document');
    try {
      await api.adminReviewVerificationDocument(id, {
        approved: approve,
        reason: reason,
        notes: approve ? 'Verified by KYC Officer' : reason
      }, user.token);
      setActionMsg({ text: `KYC Document #${id} marked as ${status}.`, type: 'success' });
      fetchLiveOperationsData();
    } catch (err: any) {
      setActionMsg({ text: err?.message || 'Failed to review document', type: 'error' });
    }
  };

  const handleReviewRefund = async (id: number, approve: boolean) => {
    if (!user?.token) return;
    const decision = approve ? 'approved' : 'rejected';
    try {
      await api.adminReviewRefund(id, {
        approved: approve,
        reason: approve ? 'Approved by Dispute Officer' : 'Rejected after dispute review'
      }, user.token);
      setActionMsg({ text: `Dispute Refund #${id} set to ${decision}.`, type: 'success' });
      fetchLiveOperationsData();
    } catch (err: any) {
      setActionMsg({ text: err?.message || 'Failed to review refund', type: 'error' });
    }
  };

  if (authLoading) {
    return (
      <div className="min-h-screen bg-zinc-950 flex items-center justify-center">
        <Spinner />
      </div>
    );
  }

  const totalActionItems = pendingProducts.length + kycDocs.length + refunds.length + appeals.length;

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100 flex flex-col font-sans selection:bg-amber-500/30 selection:text-amber-200">
      {/* Header */}
      <header className="sticky top-0 z-40 bg-zinc-950/90 backdrop-blur-xl border-b border-zinc-800/80 px-4 lg:px-8 py-3.5">
        <div className="max-w-[1720px] mx-auto flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <HuntiqLogo size="sm" href="/operations-console" />
            <div className="h-6 w-px bg-zinc-800 hidden sm:block" />
            <div>
              <div className="flex items-center gap-2">
                <span className="text-[10px] font-black uppercase tracking-widest px-2.5 py-0.5 rounded-full bg-amber-950/60 border border-amber-500/40 text-amber-400">
                  Tier-1 Operations & Verification
                </span>
                <span className="text-xs font-bold text-zinc-400 font-mono">Console</span>
              </div>
              <h1 className="text-base font-bold text-white tracking-tight mt-0.5">
                Marketplace Operations & Verification Center
              </h1>
            </div>
          </div>

          <div className="flex items-center gap-3 self-end md:self-auto">
            <button
              onClick={handleRefresh}
              className="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-zinc-900 border border-zinc-800 text-xs font-bold text-zinc-300 hover:text-white hover:border-zinc-700 transition-all shadow-sm"
            >
              <RefreshCw size={13} className={isRefreshing ? 'animate-spin text-amber-400' : ''} />
              <span>Refresh Telemetry</span>
            </button>
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-zinc-900/80 border border-zinc-800 text-xs">
              <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
              <span className="text-zinc-400 font-medium truncate max-w-[150px]">
                #{user?.id} {user?.name || 'Operator'}
              </span>
            </div>
            <button
              onClick={() => logout()}
              className="px-3 py-1.5 rounded-xl bg-red-950/30 border border-red-800/40 text-red-400 text-xs font-bold hover:bg-red-900/40 transition-colors"
            >
              Sign Out
            </button>
          </div>
        </div>
      </header>

      {/* Main Container */}
      <main className="flex-1 max-w-[1720px] w-full mx-auto px-4 lg:px-8 py-8 space-y-8">
        
        {actionMsg && (
          <div className={`p-3.5 rounded-2xl text-xs font-bold flex items-center justify-between animate-in fade-in ${
            actionMsg.type === 'success' ? 'bg-emerald-950/70 border border-emerald-500/40 text-emerald-300' : 'bg-red-950/70 border border-red-500/40 text-red-300'
          }`}>
            <span>{actionMsg.text}</span>
            <button onClick={() => setActionMsg(null)} className="text-zinc-400 hover:text-white"><X size={14} /></button>
          </div>
        )}

        {/* Live Database-Driven KPI Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Pending Inspections</span>
              <Package size={16} className="text-amber-400" />
            </div>
            <div className="text-3xl font-black text-white">{pendingProducts.length}</div>
            <div className="text-xs text-amber-400 mt-1 font-medium">Physical & spec verification queue</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Seller KYC Reviews</span>
              <FileCheck size={16} className="text-blue-400" />
            </div>
            <div className="text-3xl font-black text-white">{kycDocs.length}</div>
            <div className="text-xs text-blue-400 mt-1 font-medium">TIN, License & Fayda ID documents</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Active Disputes & Refunds</span>
              <AlertTriangle size={16} className="text-rose-400" />
            </div>
            <div className="text-3xl font-black text-white">{refunds.length}</div>
            <div className="text-xs text-rose-400 mt-1 font-medium">Buyer escrow & refund claims</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Account Appeals</span>
              <ShieldCheck size={16} className="text-emerald-400" />
            </div>
            <div className="text-3xl font-black text-white">{appeals.length}</div>
            <div className="text-xs text-emerald-400 mt-1 font-medium">Contested merchant hearing queue</div>
          </div>
        </div>

        {/* Action Queues & Workflow Section */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
            
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-3 border-b border-zinc-800">
              <div>
                <h2 className="text-lg font-black text-white tracking-tight">Active Operations Queue</h2>
                <p className="text-xs text-zinc-400">Live operational tasks requiring manual verification or decision ({totalActionItems} total)</p>
              </div>

              {/* Tabs */}
              <div className="flex items-center flex-wrap gap-1 bg-zinc-950 p-1 rounded-xl border border-zinc-800 text-xs font-bold">
                <button
                  onClick={() => setActiveTab('all')}
                  className={`px-2.5 py-1 rounded-lg transition-colors ${activeTab === 'all' ? 'bg-zinc-800 text-white' : 'text-zinc-400 hover:text-white'}`}
                >
                  All ({totalActionItems})
                </button>
                <button
                  onClick={() => setActiveTab('inspections')}
                  className={`px-2.5 py-1 rounded-lg transition-colors ${activeTab === 'inspections' ? 'bg-amber-950/80 text-amber-300' : 'text-zinc-400 hover:text-white'}`}
                >
                  Products ({pendingProducts.length})
                </button>
                <button
                  onClick={() => setActiveTab('kyc')}
                  className={`px-2.5 py-1 rounded-lg transition-colors ${activeTab === 'kyc' ? 'bg-blue-950/80 text-blue-300' : 'text-zinc-400 hover:text-white'}`}
                >
                  KYC ({kycDocs.length})
                </button>
                <button
                  onClick={() => setActiveTab('disputes')}
                  className={`px-2.5 py-1 rounded-lg transition-colors ${activeTab === 'disputes' ? 'bg-rose-950/80 text-rose-300' : 'text-zinc-400 hover:text-white'}`}
                >
                  Disputes ({refunds.length})
                </button>
                <button
                  onClick={() => setActiveTab('appeals')}
                  className={`px-2.5 py-1 rounded-lg transition-colors ${activeTab === 'appeals' ? 'bg-emerald-950/80 text-emerald-300' : 'text-zinc-400 hover:text-white'}`}
                >
                  Appeals ({appeals.length})
                </button>
              </div>
            </div>

            {loadingData && (
              <div className="py-8 flex justify-center"><Spinner /></div>
            )}

            {!loadingData && totalActionItems === 0 && (
              <div className="py-12 text-center text-zinc-500 text-xs">
                No pending items in queue. All operations current.
              </div>
            )}

            <div className="divide-y divide-zinc-800/60 text-xs">
              
              {/* Product Inspections */}
              {(activeTab === 'all' || activeTab === 'inspections') && pendingProducts.map((p) => (
                <div key={`prod-${p.id}`} className="py-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                  <div className="space-y-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 font-mono text-[10px] uppercase font-bold border border-amber-500/20">
                        Listing Inspection
                      </span>
                      <span className="text-zinc-500 font-mono text-[10px]">#{p.id}</span>
                    </div>
                    <div className="text-white font-bold text-sm truncate">{p.title}</div>
                    <div className="text-zinc-400 text-xs">
                      Price: <span className="text-emerald-400 font-bold">{Number(p.price).toLocaleString()} ETB</span> • Seller ID: #{p.seller_id}
                    </div>
                  </div>

                  <div className="flex items-center gap-2 shrink-0">
                    <button
                      onClick={() => handleApproveProduct(p.id)}
                      className="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl font-bold flex items-center gap-1 transition-all"
                      title="Approve and Publish"
                    >
                      <Check size={14} /> Approve
                    </button>
                    <button
                      onClick={() => handleRejectProduct(p.id)}
                      className="px-3 py-1.5 bg-zinc-800 hover:bg-red-950 hover:text-red-400 text-zinc-300 rounded-xl font-bold flex items-center gap-1 transition-all border border-zinc-700"
                      title="Reject with Reason"
                    >
                      <X size={14} /> Reject
                    </button>
                  </div>
                </div>
              ))}

              {/* KYC Document Reviews */}
              {(activeTab === 'all' || activeTab === 'kyc') && kycDocs.map((d) => (
                <div key={`kyc-${d.id}`} className="py-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                  <div className="space-y-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="px-2 py-0.5 rounded bg-blue-500/10 text-blue-400 font-mono text-[10px] uppercase font-bold border border-blue-500/20">
                        KYC Verification
                      </span>
                      <span className="text-zinc-500 font-mono text-[10px]">Doc #{d.id}</span>
                    </div>
                    <div className="text-white font-bold text-sm truncate">
                      {d.document_type?.toUpperCase()} — Seller #{d.seller_id}
                    </div>
                    <div className="text-zinc-400 text-xs">
                      File: {d.file_name || 'Verification Document'} • Status: <span className="uppercase text-amber-400">{d.status}</span>
                    </div>
                  </div>

                  <div className="flex items-center gap-2 shrink-0">
                    <button
                      onClick={() => handleReviewKYC(d.id, true)}
                      className="px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-bold flex items-center gap-1 transition-all"
                    >
                      <Check size={14} /> Verify ID
                    </button>
                    <button
                      onClick={() => handleReviewKYC(d.id, false)}
                      className="px-3 py-1.5 bg-zinc-800 hover:bg-red-950 hover:text-red-400 text-zinc-300 rounded-xl font-bold flex items-center gap-1 transition-all border border-zinc-700"
                    >
                      <X size={14} /> Reject
                    </button>
                  </div>
                </div>
              ))}

              {/* Escrow Disputes & Refunds */}
              {(activeTab === 'all' || activeTab === 'disputes') && refunds.map((r) => (
                <div key={`ref-${r.id}`} className="py-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                  <div className="space-y-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="px-2 py-0.5 rounded bg-rose-500/10 text-rose-400 font-mono text-[10px] uppercase font-bold border border-rose-500/20">
                        Dispute & Refund Claim
                      </span>
                      <span className="text-zinc-500 font-mono text-[10px]">Claim #{r.id}</span>
                    </div>
                    <div className="text-white font-bold text-sm truncate">
                      Reason: {r.reason} • Amount: <span className="text-rose-400 font-bold">{Number(r.amount).toLocaleString()} ETB</span>
                    </div>
                    <div className="text-zinc-400 text-xs">
                      {r.description}
                    </div>
                  </div>

                  <div className="flex items-center gap-2 shrink-0">
                    <button
                      onClick={() => handleReviewRefund(r.id, true)}
                      className="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl font-bold flex items-center gap-1 transition-all"
                    >
                      <Check size={14} /> Authorize Refund
                    </button>
                    <button
                      onClick={() => handleReviewRefund(r.id, false)}
                      className="px-3 py-1.5 bg-zinc-800 hover:bg-rose-950 hover:text-rose-400 text-zinc-300 rounded-xl font-bold flex items-center gap-1 transition-all border border-zinc-700"
                    >
                      <X size={14} /> Dismiss
                    </button>
                  </div>
                </div>
              ))}

              {/* Account Appeals */}
              {(activeTab === 'all' || activeTab === 'appeals') && appeals.map((a) => (
                <div key={`app-${a.id}`} className="py-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                  <div className="space-y-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-400 font-mono text-[10px] uppercase font-bold border border-emerald-500/20">
                        Account Appeal Hearing
                      </span>
                      <span className="text-zinc-500 font-mono text-[10px]">Appeal #{a.id}</span>
                    </div>
                    <div className="text-white font-bold text-sm truncate">
                      User #{a.user_id} — {a.appeal_type}
                    </div>
                    <div className="text-zinc-400 text-xs">
                      {a.reason}
                    </div>
                  </div>

                  <div className="flex items-center gap-2 shrink-0">
                    <Link
                      href="/admin/appeals"
                      className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-white rounded-xl font-bold flex items-center gap-1 transition-all border border-zinc-700"
                    >
                      <Eye size={14} /> Review Appeal
                    </Link>
                  </div>
                </div>
              ))}

            </div>
          </div>

          {/* Quick Operations Nav & Sub-Security Links */}
          <div className="space-y-6">
            <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
              <h2 className="text-lg font-black text-white tracking-tight">Operations Modules</h2>
              <div className="space-y-2">
                <Link
                  href="/admin/products"
                  className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-amber-500/40 hover:bg-zinc-800/40 transition-all text-xs"
                >
                  <div className="flex items-center gap-2.5">
                    <Package size={16} className="text-amber-400" />
                    <span className="font-bold text-white">Listing & Inspection Queue</span>
                  </div>
                  <ExternalLink size={14} className="text-zinc-500" />
                </Link>

                <Link
                  href="/admin/verification"
                  className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-blue-500/40 hover:bg-zinc-800/40 transition-all text-xs"
                >
                  <div className="flex items-center gap-2.5">
                    <FileCheck size={16} className="text-blue-400" />
                    <span className="font-bold text-white">Seller KYC Document Desk</span>
                  </div>
                  <ExternalLink size={14} className="text-zinc-500" />
                </Link>

                <Link
                  href="/admin/refunds"
                  className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-rose-500/40 hover:bg-zinc-800/40 transition-all text-xs"
                >
                  <div className="flex items-center gap-2.5">
                    <AlertTriangle size={16} className="text-rose-400" />
                    <span className="font-bold text-white">Dispute & Refund Resolution</span>
                  </div>
                  <ExternalLink size={14} className="text-zinc-500" />
                </Link>

                <Link
                  href="/admin/appeals"
                  className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-indigo-500/40 hover:bg-zinc-800/40 transition-all text-xs"
                >
                  <div className="flex items-center gap-2.5">
                    <ClipboardList size={16} className="text-indigo-400" />
                    <span className="font-bold text-white">Listing Appeals Hearing</span>
                  </div>
                  <ExternalLink size={14} className="text-zinc-500" />
                </Link>
              </div>
            </div>

            {/* Department Operational Shift Telemetry */}
            <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-3">
              <div className="flex items-center gap-2">
                <Shield size={16} className="text-amber-400" />
                <h3 className="text-xs font-black uppercase tracking-wider text-white">Operational Shift & Workstation</h3>
              </div>
              <p className="text-[11px] text-zinc-400 leading-relaxed">
                Workstation terminal bound to Addis Ababa Command Center. Every decision logged to immutable PostgreSQL audit trail.
              </p>
              <div className="p-3 bg-zinc-950 rounded-2xl border border-zinc-800 text-xs font-mono space-y-1.5">
                <div className="flex justify-between text-zinc-400 text-[11px]">
                  <span>Active Shift:</span>
                  <span className="text-emerald-400 font-bold">SWING / DAY 24/7</span>
                </div>
                <div className="flex justify-between text-zinc-400 text-[11px]">
                  <span>Assigned Terminal:</span>
                  <span className="text-cyan-400 font-bold">WS-OPS-01 / 02</span>
                </div>
                <div className="flex justify-between text-zinc-400 text-[11px]">
                  <span>Perimeter Zone:</span>
                  <span className="text-amber-400 font-bold">KYC_COMPLIANCE</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
