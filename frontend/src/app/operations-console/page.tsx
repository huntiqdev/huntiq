'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { Spinner } from '@/components/ui/Spinner';
import {
  Activity, CheckCircle2, Clock, AlertTriangle, ShieldCheck,
  FileCheck, Users, Search, RefreshCw, Filter, ArrowUpRight,
  ClipboardList, Package, ExternalLink, ShieldAlert
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';

export default function OperationsConsolePage() {
  const router = useRouter();
  const { user, isPlatformOwner, isSuperAdmin, adminRole, loading: authLoading, logout } = useAuth();
  const [filterStatus, setFilterStatus] = useState<'all' | 'pending' | 'in_review' | 'resolved'>('all');
  const [isRefreshing, setIsRefreshing] = useState(false);

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

  const handleRefresh = () => {
    setIsRefreshing(true);
    setTimeout(() => setIsRefreshing(false), 600);
  };

  if (authLoading) {
    return (
      <div className="min-h-screen bg-zinc-950 flex items-center justify-center">
        <Spinner />
      </div>
    );
  }

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
                  Tier-1 Operations & Inspection
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
              <span>Refresh</span>
            </button>
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-zinc-900/80 border border-zinc-800 text-xs">
              <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
              <span className="text-zinc-400 font-medium truncate max-w-[150px]">{user?.name || user?.phone || 'Authorized Staff'}</span>
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
        {/* KPI Summary Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Pending Inspections</span>
              <Package size={16} className="text-amber-400" />
            </div>
            <div className="text-3xl font-black text-white">14</div>
            <div className="text-xs text-amber-400 mt-1 font-medium">Physical & spec verification queue</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Seller KYC Reviews</span>
              <FileCheck size={16} className="text-blue-400" />
            </div>
            <div className="text-3xl font-black text-white">8</div>
            <div className="text-xs text-blue-400 mt-1 font-medium">TIN & Government ID verifications</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Active Disputes</span>
              <AlertTriangle size={16} className="text-rose-400" />
            </div>
            <div className="text-3xl font-black text-white">3</div>
            <div className="text-xs text-rose-400 mt-1 font-medium">Buyer escrow & refund claims</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Verification SLAs</span>
              <ShieldCheck size={16} className="text-emerald-400" />
            </div>
            <div className="text-3xl font-black text-white">99.2%</div>
            <div className="text-xs text-emerald-400 mt-1 font-medium">Under 4-hour review guarantee</div>
          </div>
        </div>

        {/* Action Queues & Workflow Section */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-lg font-black text-white tracking-tight">Active Operations Queue</h2>
                <p className="text-xs text-zinc-400">Live operational tasks requiring manual verification or decision</p>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-xs px-2.5 py-1 rounded-lg bg-zinc-800 text-zinc-300 font-bold border border-zinc-700">
                  Priority Sort: Urgency
                </span>
              </div>
            </div>

            <div className="divide-y divide-zinc-800/60 text-xs">
              <div className="py-3 flex items-center justify-between">
                <div className="space-y-1">
                  <span className="px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 font-mono text-[10px] uppercase font-bold border border-amber-500/20">
                    Product Inspection
                  </span>
                  <div className="text-white font-bold text-sm">Industrial Heavy Generator CAT-3406</div>
                  <div className="text-zinc-500">Seller: Addis Equipment PLC • Bole Subcity Warehouse</div>
                </div>
                <Link
                  href="/admin/products"
                  className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-white rounded-xl font-bold transition-all border border-zinc-700"
                >
                  Inspect Listing
                </Link>
              </div>

              <div className="py-3 flex items-center justify-between">
                <div className="space-y-1">
                  <span className="px-2 py-0.5 rounded bg-blue-500/10 text-blue-400 font-mono text-[10px] uppercase font-bold border border-blue-500/20">
                    KYC Verification
                  </span>
                  <div className="text-white font-bold text-sm">Mekelle Machinery Trading - License Verification</div>
                  <div className="text-zinc-500">Document Type: Commercial Registration Certificate</div>
                </div>
                <Link
                  href="/admin/verification"
                  className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-white rounded-xl font-bold transition-all border border-zinc-700"
                >
                  Verify ID
                </Link>
              </div>

              <div className="py-3 flex items-center justify-between">
                <div className="space-y-1">
                  <span className="px-2 py-0.5 rounded bg-rose-500/10 text-rose-400 font-mono text-[10px] uppercase font-bold border border-rose-500/20">
                    Escrow Refund Claim
                  </span>
                  <div className="text-white font-bold text-sm">Transaction #TX-88319 - Spec Mismatch Inquiry</div>
                  <div className="text-zinc-500">Buyer requested refund: 25,000 ETB • Status: Evidence Submitted</div>
                </div>
                <Link
                  href="/admin/refunds"
                  className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-white rounded-xl font-bold transition-all border border-zinc-700"
                >
                  Review Claim
                </Link>
              </div>
            </div>
          </div>

          {/* Quick Operations Nav */}
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
        </div>
      </main>
    </div>
  );
}
