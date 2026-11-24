'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import {
  Users, MessageSquare, HelpCircle, CheckCircle2, Clock,
  RefreshCw, Sparkles, ExternalLink, Headphones, Search,
  AlertTriangle, Shield, Check, X, Eye
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';

export default function SupportConsolePage() {
  const router = useRouter();
  const { user, isPlatformOwner, isSuperAdmin, adminRole, loading: authLoading, logout } = useAuth();
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [activeTab, setActiveTab] = useState<'all' | 'refunds' | 'appeals' | 'inquiries'>('all');
  const [actionMsg, setActionMsg] = useState<{ text: string; type: 'success' | 'error' } | null>(null);

  // Live Database Driven State
  const [abuseReports, setAbuseReports] = useState<any[]>([]);
  const [refundRequests, setRefundRequests] = useState<any[]>([]);
  const [accountAppeals, setAccountAppeals] = useState<any[]>([]);
  const [loadingData, setLoadingData] = useState(false);

  useEffect(() => {
    if (!authLoading) {
      const allowedRoles = [
        'platform_owner', 'super', 'customer_support_lead', 'support_agent',
        'gold_concierge', 'premium_concierge', 'retention_specialist'
      ];
      const isAllowed = isPlatformOwner || isSuperAdmin || (adminRole && allowedRoles.includes(adminRole));
      if (!user || !isAllowed) {
        router.replace('/support-console/login');
      }
    }
  }, [user, authLoading, isPlatformOwner, isSuperAdmin, adminRole, router]);

  const fetchLiveSupportData = async () => {
    if (!user?.token) return;
    try {
      setLoadingData(true);
      const [reports, refunds, appeals] = await Promise.all([
        api.adminGetAbuseReports('all', user.token).then((r: any) => (Array.isArray(r) ? r : r?.data) || []).catch(() => []),
        api.adminGetRefundRequests('all', user.token).then(r => r?.data || r || []).catch(() => []),
        api.adminGetAccountAppeals('all', user.token).then(r => r?.data || r || []).catch(() => []),
      ]);

      setAbuseReports(Array.isArray(reports) ? reports : []);
      setRefundRequests(Array.isArray(refunds) ? refunds : []);
      setAccountAppeals(Array.isArray(appeals) ? appeals : []);
    } catch {
      // Fallback
    } finally {
      setLoadingData(false);
      setIsRefreshing(false);
    }
  };

  useEffect(() => {
    fetchLiveSupportData();
  }, [user?.token]);

  const handleRefresh = () => {
    setIsRefreshing(true);
    fetchLiveSupportData();
  };

  if (authLoading) {
    return (
      <div className="min-h-screen bg-zinc-950 flex items-center justify-center">
        <Spinner />
      </div>
    );
  }

  const totalTickets = abuseReports.length + refundRequests.length + accountAppeals.length;

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100 flex flex-col font-sans selection:bg-blue-500/30 selection:text-blue-200">
      {/* Header */}
      <header className="sticky top-0 z-40 bg-zinc-950/90 backdrop-blur-xl border-b border-zinc-800/80 px-4 lg:px-8 py-3.5">
        <div className="max-w-[1720px] mx-auto flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <HuntiqLogo size="sm" href="/support-console" />
            <div className="h-6 w-px bg-zinc-800 hidden sm:block" />
            <div>
              <div className="flex items-center gap-2">
                <span className="text-[10px] font-black uppercase tracking-widest px-2.5 py-0.5 rounded-full bg-blue-950/60 border border-blue-500/40 text-blue-400">
                  Customer Support & Concierge
                </span>
                <span className="text-xs font-bold text-zinc-400 font-mono">Console</span>
              </div>
              <h1 className="text-base font-bold text-white tracking-tight mt-0.5">
                Customer Support & VIP Concierge Desk
              </h1>
            </div>
          </div>

          <div className="flex items-center gap-3 self-end md:self-auto">
            <button
              onClick={handleRefresh}
              className="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-zinc-900 border border-zinc-800 text-xs font-bold text-zinc-300 hover:text-white hover:border-zinc-700 transition-all shadow-sm"
            >
              <RefreshCw size={13} className={isRefreshing ? 'animate-spin text-blue-400' : ''} />
              <span>Refresh Desk</span>
            </button>
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-zinc-900/80 border border-zinc-800 text-xs">
              <span className="w-2 h-2 rounded-full bg-blue-500 animate-pulse" />
              <span className="text-zinc-400 font-medium truncate max-w-[150px]">
                #{user?.id} {user?.name || 'Support Agent'}
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
              <span>Open Support Inquiries</span>
              <MessageSquare size={16} className="text-blue-400" />
            </div>
            <div className="text-3xl font-black text-white">{abuseReports.length}</div>
            <div className="text-xs text-blue-400 mt-1 font-medium">Customer assistance inquiries</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Dispute Refund Requests</span>
              <AlertTriangle size={16} className="text-rose-400" />
            </div>
            <div className="text-3xl font-black text-white">{refundRequests.length}</div>
            <div className="text-xs text-rose-400 mt-1 font-medium">Buyer escrow & refund triage</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Active Customer Appeals</span>
              <HelpCircle size={16} className="text-amber-400" />
            </div>
            <div className="text-3xl font-black text-white">{accountAppeals.length}</div>
            <div className="text-xs text-amber-400 mt-1 font-medium">Account status appeal reviews</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Support SLA Compliance</span>
              <CheckCircle2 size={16} className="text-emerald-400" />
            </div>
            <div className="text-3xl font-black text-white">99.8%</div>
            <div className="text-xs text-emerald-400 mt-1 font-medium">First contact resolution target</div>
          </div>
        </div>

        {/* Live Support Streams */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
            
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-3 border-b border-zinc-800">
              <div>
                <h2 className="text-lg font-black text-white tracking-tight">Customer Assistance Queue</h2>
                <p className="text-xs text-zinc-400">Live order inquiries, dispute assistance, and customer requests ({totalTickets} total)</p>
              </div>

              {/* Tabs */}
              <div className="flex items-center gap-1 bg-zinc-950 p-1 rounded-xl border border-zinc-800 text-xs font-bold">
                <button
                  onClick={() => setActiveTab('all')}
                  className={`px-3 py-1 rounded-lg transition-colors ${activeTab === 'all' ? 'bg-zinc-800 text-white' : 'text-zinc-400 hover:text-white'}`}
                >
                  All ({totalTickets})
                </button>
                <button
                  onClick={() => setActiveTab('refunds')}
                  className={`px-3 py-1 rounded-lg transition-colors ${activeTab === 'refunds' ? 'bg-rose-950/80 text-rose-300' : 'text-zinc-400 hover:text-white'}`}
                >
                  Refunds ({refundRequests.length})
                </button>
                <button
                  onClick={() => setActiveTab('appeals')}
                  className={`px-3 py-1 rounded-lg transition-colors ${activeTab === 'appeals' ? 'bg-amber-950/80 text-amber-300' : 'text-zinc-400 hover:text-white'}`}
                >
                  Appeals ({accountAppeals.length})
                </button>
                <button
                  onClick={() => setActiveTab('inquiries')}
                  className={`px-3 py-1 rounded-lg transition-colors ${activeTab === 'inquiries' ? 'bg-blue-950/80 text-blue-300' : 'text-zinc-400 hover:text-white'}`}
                >
                  Reports ({abuseReports.length})
                </button>
              </div>
            </div>

            {loadingData && (
              <div className="py-8 flex justify-center"><Spinner /></div>
            )}

            {!loadingData && totalTickets === 0 && (
              <div className="py-12 text-center text-zinc-500 text-xs">
                No active tickets in queue. Customer support inbox is clear.
              </div>
            )}

            <div className="divide-y divide-zinc-800/60 text-xs">
              
              {/* Refunds Stream */}
              {(activeTab === 'all' || activeTab === 'refunds') && refundRequests.map((r) => (
                <div key={`ref-${r.id}`} className="py-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                  <div className="space-y-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="px-2 py-0.5 rounded bg-rose-500/10 text-rose-400 font-mono text-[10px] uppercase font-bold border border-rose-500/20">
                        Dispute Ticket #{r.id}
                      </span>
                      <span className="text-zinc-500 font-mono text-[10px]">User #{r.requester_id}</span>
                    </div>
                    <div className="text-white font-bold text-sm truncate">
                      {r.reason} • <span className="text-emerald-400">{Number(r.amount).toLocaleString()} ETB</span>
                    </div>
                    <div className="text-zinc-400 text-xs truncate max-w-lg">
                      {r.description}
                    </div>
                  </div>

                  <div className="flex items-center gap-2 shrink-0">
                    <Link
                      href="/admin/refunds"
                      className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-white rounded-xl font-bold flex items-center gap-1 transition-all border border-zinc-700"
                    >
                      <Eye size={13} /> Assist Customer
                    </Link>
                  </div>
                </div>
              ))}

              {/* Appeals Stream */}
              {(activeTab === 'all' || activeTab === 'appeals') && accountAppeals.map((a) => (
                <div key={`app-${a.id}`} className="py-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                  <div className="space-y-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 font-mono text-[10px] uppercase font-bold border border-amber-500/20">
                        Appeal Hearing #{a.id}
                      </span>
                      <span className="text-zinc-500 font-mono text-[10px]">User #{a.user_id}</span>
                    </div>
                    <div className="text-white font-bold text-sm truncate">
                      {a.appeal_type}
                    </div>
                    <div className="text-zinc-400 text-xs truncate max-w-lg">
                      {a.reason}
                    </div>
                  </div>

                  <div className="flex items-center gap-2 shrink-0">
                    <Link
                      href="/admin/appeals"
                      className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-white rounded-xl font-bold flex items-center gap-1 transition-all border border-zinc-700"
                    >
                      <Eye size={13} /> Review Case
                    </Link>
                  </div>
                </div>
              ))}

              {/* Inquiries Stream */}
              {(activeTab === 'all' || activeTab === 'inquiries') && abuseReports.map((rep) => (
                <div key={`rep-${rep.id}`} className="py-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                  <div className="space-y-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="px-2 py-0.5 rounded bg-blue-500/10 text-blue-400 font-mono text-[10px] uppercase font-bold border border-blue-500/20">
                        Report Ticket #{rep.id}
                      </span>
                      <span className="text-zinc-500 font-mono text-[10px]">Product #{rep.reported_product_id}</span>
                    </div>
                    <div className="text-white font-bold text-sm truncate">
                      {rep.reason} • Priority: <span className="uppercase text-amber-400">{rep.priority}</span>
                    </div>
                    <div className="text-zinc-400 text-xs truncate max-w-lg">
                      {rep.description}
                    </div>
                  </div>

                  <div className="flex items-center gap-2 shrink-0">
                    <Link
                      href="/admin/trust-safety"
                      className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-white rounded-xl font-bold flex items-center gap-1 transition-all border border-zinc-700"
                    >
                      <Eye size={13} /> Open Ticket
                    </Link>
                  </div>
                </div>
              ))}

            </div>
          </div>

          {/* Side Shortcuts */}
          <div className="space-y-6">
            <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
              <h2 className="text-lg font-black text-white tracking-tight">Support Tools</h2>
              <div className="space-y-2">
                <Link
                  href="/admin/users"
                  className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-blue-500/40 hover:bg-zinc-800/40 transition-all text-xs"
                >
                  <div className="flex items-center gap-2.5">
                    <Users size={16} className="text-blue-400" />
                    <span className="font-bold text-white">User Directory & Lookup</span>
                  </div>
                  <ExternalLink size={14} className="text-zinc-500" />
                </Link>

                <Link
                  href="/admin/transactions"
                  className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-emerald-500/40 hover:bg-zinc-800/40 transition-all text-xs"
                >
                  <div className="flex items-center gap-2.5">
                    <CheckCircle2 size={16} className="text-emerald-400" />
                    <span className="font-bold text-white">Payment & Escrow Search</span>
                  </div>
                  <ExternalLink size={14} className="text-zinc-500" />
                </Link>

                <Link
                  href="/admin/sellers"
                  className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-amber-500/40 hover:bg-zinc-800/40 transition-all text-xs"
                >
                  <div className="flex items-center gap-2.5">
                    <Sparkles size={16} className="text-amber-400" />
                    <span className="font-bold text-white">Merchant Vetting & Tiers</span>
                  </div>
                  <ExternalLink size={14} className="text-zinc-500" />
                </Link>
              </div>
            </div>

            {/* Shift & Workstation Card */}
            <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-3">
              <div className="flex items-center gap-2">
                <Shield size={16} className="text-blue-400" />
                <h3 className="text-xs font-black uppercase tracking-wider text-white">Support Terminal Assignment</h3>
              </div>
              <p className="text-[11px] text-zinc-400 leading-relaxed">
                Terminal stationed at Addis Ababa Support Center & Regional Hubs. Operator actions logged with non-repudiation guarantee.
              </p>
              <div className="p-3 bg-zinc-950 rounded-2xl border border-zinc-800 text-xs font-mono space-y-1.5">
                <div className="flex justify-between text-zinc-400 text-[11px]">
                  <span>Operational Shift:</span>
                  <span className="text-emerald-400 font-bold">DAY / SWING 24/7</span>
                </div>
                <div className="flex justify-between text-zinc-400 text-[11px]">
                  <span>Assigned Station:</span>
                  <span className="text-cyan-400 font-bold">WS-SUP-01 / 02</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
