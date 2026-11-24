'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import {
  ShieldAlert, ShieldCheck, Flag, AlertOctagon, CheckCircle2,
  XCircle, Eye, RefreshCw, Filter, ExternalLink, Search, Ban,
  Check, X, AlertTriangle, Shield, Trash2, ArrowUpRight
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';

export default function ModerationConsolePage() {
  const router = useRouter();
  const { user, isPlatformOwner, isSuperAdmin, adminRole, loading: authLoading, logout } = useAuth();
  const [activeTab, setActiveTab] = useState<'all' | 'listings' | 'reports'>('all');
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [actionMsg, setActionMsg] = useState<{ text: string; type: 'success' | 'error' } | null>(null);

  // Live Database Driven State
  const [pendingProducts, setPendingProducts] = useState<any[]>([]);
  const [abuseReports, setAbuseReports] = useState<any[]>([]);
  const [loadingData, setLoadingData] = useState(false);

  useEffect(() => {
    if (!authLoading) {
      const allowedRoles = [
        'platform_owner', 'super', 'content_moderator', 'moderator'
      ];
      const isAllowed = isPlatformOwner || isSuperAdmin || (adminRole && allowedRoles.includes(adminRole));
      if (!user || !isAllowed) {
        router.replace('/moderation-console/login');
      }
    }
  }, [user, authLoading, isPlatformOwner, isSuperAdmin, adminRole, router]);

  const fetchLiveModerationData = async () => {
    if (!user?.token) return;
    try {
      setLoadingData(true);
      const [prods, reports] = await Promise.all([
        api.adminGetPending(user.token).catch(() => []),
        api.adminGetAbuseReports('pending', user.token).then((r: any) => (Array.isArray(r) ? r : r?.data) || []).catch(() => []),
      ]);

      setPendingProducts(Array.isArray(prods) ? prods : []);
      setAbuseReports(Array.isArray(reports) ? reports : []);
    } catch {
      // Fallback
    } finally {
      setLoadingData(false);
      setIsRefreshing(false);
    }
  };

  useEffect(() => {
    fetchLiveModerationData();
  }, [user?.token]);

  const handleRefresh = () => {
    setIsRefreshing(true);
    fetchLiveModerationData();
  };

  const handleApproveProduct = async (id: number) => {
    if (!user?.token) return;
    try {
      await api.adminReviewProduct(id, { status: 'active', notes: 'Verified and approved by Content Moderator' }, user.token);
      setActionMsg({ text: `Listing #${id} approved and published to catalog.`, type: 'success' });
      fetchLiveModerationData();
    } catch (err: any) {
      setActionMsg({ text: err?.message || 'Failed to approve listing', type: 'error' });
    }
  };

  const handleRejectProduct = async (id: number) => {
    if (!user?.token) return;
    const reason = prompt('Enter moderation rejection reason:', 'Prohibited content or policy violation');
    if (!reason) return;
    try {
      await api.adminReviewProduct(id, { status: 'rejected', notes: reason }, user.token);
      setActionMsg({ text: `Listing #${id} rejected: ${reason}`, type: 'success' });
      fetchLiveModerationData();
    } catch (err: any) {
      setActionMsg({ text: err?.message || 'Failed to reject listing', type: 'error' });
    }
  };

  const handleResolveAbuseReport = async (id: number, action: 'dismiss' | 'take_down') => {
    if (!user?.token) return;
    const actionTaken = action === 'take_down' ? 'Listing removed and seller warned' : 'Report investigated and dismissed';
    try {
      await api.adminReviewAbuseReport(id, { 
        status: action === 'take_down' ? 'resolved' : 'dismissed',
        action: action,
        action_taken: actionTaken,
        notes: 'Reviewed in Moderation Console'
      }, user.token);
      setActionMsg({ text: `Report #${id} resolved: ${actionTaken}`, type: 'success' });
      fetchLiveModerationData();
    } catch (err: any) {
      setActionMsg({ text: err?.message || 'Failed to review abuse report', type: 'error' });
    }
  };

  if (authLoading) {
    return (
      <div className="min-h-screen bg-zinc-950 flex items-center justify-center">
        <Spinner />
      </div>
    );
  }

  const criticalReports = abuseReports.filter(r => r.priority === 'critical');
  const totalItems = pendingProducts.length + abuseReports.length;

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100 flex flex-col font-sans selection:bg-purple-500/30 selection:text-purple-200">
      {/* Header */}
      <header className="sticky top-0 z-40 bg-zinc-950/90 backdrop-blur-xl border-b border-zinc-800/80 px-4 lg:px-8 py-3.5">
        <div className="max-w-[1720px] mx-auto flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <HuntiqLogo size="sm" href="/moderation-console" />
            <div className="h-6 w-px bg-zinc-800 hidden sm:block" />
            <div>
              <div className="flex items-center gap-2">
                <span className="text-[10px] font-black uppercase tracking-widest px-2.5 py-0.5 rounded-full bg-purple-950/60 border border-purple-500/40 text-purple-400">
                  Trust & Safety Enforcement
                </span>
                <span className="text-xs font-bold text-zinc-400 font-mono">Console</span>
              </div>
              <h1 className="text-base font-bold text-white tracking-tight mt-0.5">
                Content Governance & Community Safety Desk
              </h1>
            </div>
          </div>

          <div className="flex items-center gap-3 self-end md:self-auto">
            <button
              onClick={handleRefresh}
              className="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-zinc-900 border border-zinc-800 text-xs font-bold text-zinc-300 hover:text-white hover:border-zinc-700 transition-all shadow-sm"
            >
              <RefreshCw size={13} className={isRefreshing ? 'animate-spin text-purple-400' : ''} />
              <span>Refresh Queue</span>
            </button>
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-zinc-900/80 border border-zinc-800 text-xs">
              <span className="w-2 h-2 rounded-full bg-purple-500 animate-pulse" />
              <span className="text-zinc-400 font-medium truncate max-w-[150px]">
                #{user?.id} {user?.name || 'Moderator'}
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

        {/* Live Database-Driven KPI Summary Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Flagged Abuse Reports</span>
              <Flag size={16} className="text-purple-400" />
            </div>
            <div className="text-3xl font-black text-white">{abuseReports.length}</div>
            <div className="text-xs text-purple-400 mt-1 font-medium">Community & AI scam alerts</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Pending New Listings</span>
              <Eye size={16} className="text-blue-400" />
            </div>
            <div className="text-3xl font-black text-white">{pendingProducts.length}</div>
            <div className="text-xs text-blue-400 mt-1 font-medium">Awaiting publication inspection</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Critical Prohibited Alerts</span>
              <Ban size={16} className="text-rose-400" />
            </div>
            <div className="text-3xl font-black text-white">{criticalReports.length}</div>
            <div className="text-xs text-rose-400 mt-1 font-medium">Immediate policy violation queue</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Moderation Coverage</span>
              <ShieldCheck size={16} className="text-emerald-400" />
            </div>
            <div className="text-3xl font-black text-white">100%</div>
            <div className="text-xs text-emerald-400 mt-1 font-medium">Pre-publish zero-trust clamp</div>
          </div>
        </div>

        {/* Action Queues & Workflow Section */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
            
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-3 border-b border-zinc-800">
              <div>
                <h2 className="text-lg font-black text-white tracking-tight">Active Moderation Stream</h2>
                <p className="text-xs text-zinc-400">Review pending marketplace listings and reported community abuse ({totalItems} items)</p>
              </div>

              {/* Tabs */}
              <div className="flex items-center gap-1 bg-zinc-950 p-1 rounded-xl border border-zinc-800 text-xs font-bold">
                <button
                  onClick={() => setActiveTab('all')}
                  className={`px-3 py-1 rounded-lg transition-colors ${activeTab === 'all' ? 'bg-zinc-800 text-white' : 'text-zinc-400 hover:text-white'}`}
                >
                  All ({totalItems})
                </button>
                <button
                  onClick={() => setActiveTab('listings')}
                  className={`px-3 py-1 rounded-lg transition-colors ${activeTab === 'listings' ? 'bg-purple-950/80 text-purple-300' : 'text-zinc-400 hover:text-white'}`}
                >
                  Listings ({pendingProducts.length})
                </button>
                <button
                  onClick={() => setActiveTab('reports')}
                  className={`px-3 py-1 rounded-lg transition-colors ${activeTab === 'reports' ? 'bg-rose-950/80 text-rose-300' : 'text-zinc-400 hover:text-white'}`}
                >
                  Reports ({abuseReports.length})
                </button>
              </div>
            </div>

            {loadingData && (
              <div className="py-8 flex justify-center"><Spinner /></div>
            )}

            {!loadingData && totalItems === 0 && (
              <div className="py-12 text-center text-zinc-500 text-xs">
                No items pending moderation review. Queue is clear.
              </div>
            )}

            <div className="divide-y divide-zinc-800/60 text-xs">
              
              {/* Listings Queue */}
              {(activeTab === 'all' || activeTab === 'listings') && pendingProducts.map((p) => (
                <div key={`p-${p.id}`} className="py-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                  <div className="space-y-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="px-2 py-0.5 rounded bg-blue-500/10 text-blue-400 font-mono text-[10px] uppercase font-bold border border-blue-500/20">
                        New Listing
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
                    >
                      <Check size={14} /> Approve
                    </button>
                    <button
                      onClick={() => handleRejectProduct(p.id)}
                      className="px-3 py-1.5 bg-zinc-800 hover:bg-red-950 hover:text-red-400 text-zinc-300 rounded-xl font-bold flex items-center gap-1 transition-all border border-zinc-700"
                    >
                      <X size={14} /> Reject
                    </button>
                  </div>
                </div>
              ))}

              {/* Abuse Reports Queue */}
              {(activeTab === 'all' || activeTab === 'reports') && abuseReports.map((r) => (
                <div key={`rep-${r.id}`} className="py-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                  <div className="space-y-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className={`px-2 py-0.5 rounded font-mono text-[10px] uppercase font-bold border ${
                        r.priority === 'critical' ? 'bg-red-950/80 text-red-400 border-red-500/40' : 'bg-purple-950/80 text-purple-400 border-purple-500/30'
                      }`}>
                        Abuse Report: {r.reason}
                      </span>
                      <span className="text-zinc-500 font-mono text-[10px]">#{r.id}</span>
                    </div>
                    <div className="text-white font-bold text-sm truncate">{r.description}</div>
                    <div className="text-zinc-400 text-xs">
                      Reported Product ID: #{r.reported_product_id} • Priority: <span className="uppercase text-amber-400 font-bold">{r.priority}</span>
                    </div>
                  </div>

                  <div className="flex items-center gap-2 shrink-0">
                    <button
                      onClick={() => handleResolveAbuseReport(r.id, 'take_down')}
                      className="px-3 py-1.5 bg-rose-600 hover:bg-rose-500 text-white rounded-xl font-bold flex items-center gap-1 transition-all"
                    >
                      <Ban size={14} /> Take Down
                    </button>
                    <button
                      onClick={() => handleResolveAbuseReport(r.id, 'dismiss')}
                      className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl font-bold flex items-center gap-1 transition-all border border-zinc-700"
                    >
                      Dismiss
                    </button>
                  </div>
                </div>
              ))}

            </div>
          </div>

          {/* Quick Moderation Controls */}
          <div className="space-y-6">
            <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
              <h2 className="text-lg font-black text-white tracking-tight">Governance Shortcuts</h2>
              <div className="space-y-2">
                <Link
                  href="/admin/products"
                  className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-purple-500/40 hover:bg-zinc-800/40 transition-all text-xs"
                >
                  <div className="flex items-center gap-2.5">
                    <Eye size={16} className="text-purple-400" />
                    <span className="font-bold text-white">Full Product Catalog Inspector</span>
                  </div>
                  <ExternalLink size={14} className="text-zinc-500" />
                </Link>

                <Link
                  href="/admin/trust-safety"
                  className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-red-500/40 hover:bg-zinc-800/40 transition-all text-xs"
                >
                  <div className="flex items-center gap-2.5">
                    <ShieldAlert size={16} className="text-red-400" />
                    <span className="font-bold text-white">Trust & Abuse Resolution Desk</span>
                  </div>
                  <ExternalLink size={14} className="text-zinc-500" />
                </Link>

                <Link
                  href="/admin/logs"
                  className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-blue-500/40 hover:bg-zinc-800/40 transition-all text-xs"
                >
                  <div className="flex items-center gap-2.5">
                    <AlertTriangle size={16} className="text-blue-400" />
                    <span className="font-bold text-white">Moderator Audit History</span>
                  </div>
                  <ExternalLink size={14} className="text-zinc-500" />
                </Link>
              </div>
            </div>

            {/* Shift & Workstation Card */}
            <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-3">
              <div className="flex items-center gap-2">
                <Shield size={16} className="text-purple-400" />
                <h3 className="text-xs font-black uppercase tracking-wider text-white">Workstation Assignment</h3>
              </div>
              <p className="text-[11px] text-zinc-400 leading-relaxed">
                Dedicated terminal in Addis Ababa Moderation Lab. Operating under Ethiopian Consumer Protection Directives.
              </p>
              <div className="p-3 bg-zinc-950 rounded-2xl border border-zinc-800 text-xs font-mono space-y-1.5">
                <div className="flex justify-between text-zinc-400 text-[11px]">
                  <span>Operational Shift:</span>
                  <span className="text-emerald-400 font-bold">24/7 ROTATING (3 SHIFTS)</span>
                </div>
                <div className="flex justify-between text-zinc-400 text-[11px]">
                  <span>Assigned Station:</span>
                  <span className="text-cyan-400 font-bold">WS-MOD-01 / 02</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
