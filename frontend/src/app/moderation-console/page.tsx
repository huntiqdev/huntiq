'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { Spinner } from '@/components/ui/Spinner';
import {
  ShieldAlert, ShieldCheck, Flag, AlertOctagon, CheckCircle2,
  XCircle, Eye, RefreshCw, Filter, ExternalLink, Search, Ban
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';

export default function ModerationConsolePage() {
  const router = useRouter();
  const { user, isPlatformOwner, isSuperAdmin, adminRole, loading: authLoading, logout } = useAuth();
  const [filterType, setFilterType] = useState<'all' | 'flagged' | 'pending' | 'prohibited'>('all');
  const [isRefreshing, setIsRefreshing] = useState(false);

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
                Content Governance & Community Safety
              </h1>
            </div>
          </div>

          <div className="flex items-center gap-3 self-end md:self-auto">
            <button
              onClick={handleRefresh}
              className="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-zinc-900 border border-zinc-800 text-xs font-bold text-zinc-300 hover:text-white hover:border-zinc-700 transition-all shadow-sm"
            >
              <RefreshCw size={13} className={isRefreshing ? 'animate-spin text-purple-400' : ''} />
              <span>Refresh</span>
            </button>
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-zinc-900/80 border border-zinc-800 text-xs">
              <span className="w-2 h-2 rounded-full bg-purple-500 animate-pulse" />
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
              <span>Flagged Items</span>
              <Flag size={16} className="text-purple-400" />
            </div>
            <div className="text-3xl font-black text-white">5</div>
            <div className="text-xs text-purple-400 mt-1 font-medium">Reported by community users</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>New Listings Queue</span>
              <Eye size={16} className="text-blue-400" />
            </div>
            <div className="text-3xl font-black text-white">22</div>
            <div className="text-xs text-blue-400 mt-1 font-medium">Awaiting publication approval</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Prohibited Item Blocks</span>
              <Ban size={16} className="text-rose-400" />
            </div>
            <div className="text-3xl font-black text-white">19</div>
            <div className="text-xs text-rose-400 mt-1 font-medium">Auto-intercepted this week</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Community Integrity</span>
              <ShieldCheck size={16} className="text-emerald-400" />
            </div>
            <div className="text-3xl font-black text-white">99.8%</div>
            <div className="text-xs text-emerald-400 mt-1 font-medium">Trust score across active listings</div>
          </div>
        </div>

        {/* Content Moderation Queues */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-lg font-black text-white tracking-tight">Active Moderation Queue</h2>
                <p className="text-xs text-zinc-400">Listings flagged for policy review or suspicious activity</p>
              </div>
            </div>

            <div className="divide-y divide-zinc-800/60 text-xs">
              <div className="py-3 flex items-center justify-between">
                <div className="space-y-1">
                  <span className="px-2 py-0.5 rounded bg-purple-500/10 text-purple-400 font-mono text-[10px] uppercase font-bold border border-purple-500/20">
                    Policy Check
                  </span>
                  <div className="text-white font-bold text-sm">Real Estate Villa in Bole Bulbula</div>
                  <div className="text-zinc-500">Reason: Price appears 90% below average district market value</div>
                </div>
                <Link
                  href="/admin/products"
                  className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-white rounded-xl font-bold transition-all border border-zinc-700"
                >
                  Review
                </Link>
              </div>

              <div className="py-3 flex items-center justify-between">
                <div className="space-y-1">
                  <span className="px-2 py-0.5 rounded bg-rose-500/10 text-rose-400 font-mono text-[10px] uppercase font-bold border border-rose-500/20">
                    Duplicate Images
                  </span>
                  <div className="text-white font-bold text-sm">Toyota Hilux Double Cab 2022</div>
                  <div className="text-zinc-500">Photos match 3 previous listings with different contact phones</div>
                </div>
                <Link
                  href="/admin/trust-safety"
                  className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-white rounded-xl font-bold transition-all border border-zinc-700"
                >
                  Investigate
                </Link>
              </div>
            </div>
          </div>

          {/* Quick Moderation Links */}
          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
            <h2 className="text-lg font-black text-white tracking-tight">Moderation Tools</h2>
            <div className="space-y-2">
              <Link
                href="/admin/trust-safety"
                className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-purple-500/40 hover:bg-zinc-800/40 transition-all text-xs"
              >
                <div className="flex items-center gap-2.5">
                  <ShieldAlert size={16} className="text-purple-400" />
                  <span className="font-bold text-white">Trust & Safety Console</span>
                </div>
                <ExternalLink size={14} className="text-zinc-500" />
              </Link>

              <Link
                href="/admin/products"
                className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-blue-500/40 hover:bg-zinc-800/40 transition-all text-xs"
              >
                <div className="flex items-center gap-2.5">
                  <Eye size={16} className="text-blue-400" />
                  <span className="font-bold text-white">All Marketplace Listings</span>
                </div>
                <ExternalLink size={14} className="text-zinc-500" />
              </Link>

              <Link
                href="/admin/appeals"
                className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-amber-500/40 hover:bg-zinc-800/40 transition-all text-xs"
              >
                <div className="flex items-center gap-2.5">
                  <CheckCircle2 size={16} className="text-amber-400" />
                  <span className="font-bold text-white">Seller Account Appeals</span>
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
