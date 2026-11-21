'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { Spinner } from '@/components/ui/Spinner';
import {
  Users, MessageSquare, HelpCircle, CheckCircle2, Clock,
  RefreshCw, Sparkles, ExternalLink, Headphones, Search
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';

export default function SupportConsolePage() {
  const router = useRouter();
  const { user, isPlatformOwner, isSuperAdmin, adminRole, loading: authLoading, logout } = useAuth();
  const [isRefreshing, setIsRefreshing] = useState(false);

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
              <span>Refresh</span>
            </button>
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-zinc-900/80 border border-zinc-800 text-xs">
              <span className="w-2 h-2 rounded-full bg-blue-500 animate-pulse" />
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
              <span>Open Tickets</span>
              <MessageSquare size={16} className="text-blue-400" />
            </div>
            <div className="text-3xl font-black text-white">7</div>
            <div className="text-xs text-blue-400 mt-1 font-medium">Customer questions & order triage</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>VIP Concierge Requests</span>
              <Sparkles size={16} className="text-amber-400" />
            </div>
            <div className="text-3xl font-black text-white">4</div>
            <div className="text-xs text-amber-400 mt-1 font-medium">Gold & Premium dedicated assistance</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Avg Response Time</span>
              <Clock size={16} className="text-emerald-400" />
            </div>
            <div className="text-3xl font-black text-white">8m</div>
            <div className="text-xs text-emerald-400 mt-1 font-medium">Live support desk latency</div>
          </div>

          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-2xl p-5 shadow-lg">
            <div className="flex items-center justify-between text-zinc-400 text-xs font-bold uppercase tracking-wider mb-2">
              <span>Customer Satisfaction</span>
              <CheckCircle2 size={16} className="text-teal-400" />
            </div>
            <div className="text-3xl font-black text-white">98.5%</div>
            <div className="text-xs text-teal-400 mt-1 font-medium">Positive feedback rating</div>
          </div>
        </div>

        {/* Support Queues */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-lg font-black text-white tracking-tight">Active Inquiries & Assistance</h2>
                <p className="text-xs text-zinc-400">Incoming tickets from buyers, sellers, and landlords</p>
              </div>
            </div>

            <div className="divide-y divide-zinc-800/60 text-xs">
              <div className="py-3 flex items-center justify-between">
                <div className="space-y-1">
                  <span className="px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 font-mono text-[10px] uppercase font-bold border border-amber-500/20">
                    VIP Concierge
                  </span>
                  <div className="text-white font-bold text-sm">Property Inspection Scheduling • Bole Atlas Villa</div>
                  <div className="text-zinc-500">Premium User: Dawit A. (+251911445566) • Priority: High</div>
                </div>
                <Link
                  href="/admin/premium-users"
                  className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-white rounded-xl font-bold transition-all border border-zinc-700"
                >
                  Assist
                </Link>
              </div>

              <div className="py-3 flex items-center justify-between">
                <div className="space-y-1">
                  <span className="px-2 py-0.5 rounded bg-blue-500/10 text-blue-400 font-mono text-[10px] uppercase font-bold border border-blue-500/20">
                    Payment Inquiry
                  </span>
                  <div className="text-white font-bold text-sm">Telebirr SuperApp Settlement Query</div>
                  <div className="text-zinc-500">User: Sara K. • Order #ORD-1092 • Listing fee verification</div>
                </div>
                <Link
                  href="/admin/transactions"
                  className="px-3 py-1.5 bg-zinc-800 hover:bg-zinc-700 text-white rounded-xl font-bold transition-all border border-zinc-700"
                >
                  View Details
                </Link>
              </div>
            </div>
          </div>

          {/* Quick Support Links */}
          <div className="bg-zinc-900/70 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
            <h2 className="text-lg font-black text-white tracking-tight">Support Tools</h2>
            <div className="space-y-2">
              <Link
                href="/admin/premium-users"
                className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-amber-500/40 hover:bg-zinc-800/40 transition-all text-xs"
              >
                <div className="flex items-center gap-2.5">
                  <Sparkles size={16} className="text-amber-400" />
                  <span className="font-bold text-white">VIP & Premium Subscribers</span>
                </div>
                <ExternalLink size={14} className="text-zinc-500" />
              </Link>

              <Link
                href="/admin/sellers"
                className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-blue-500/40 hover:bg-zinc-800/40 transition-all text-xs"
              >
                <div className="flex items-center gap-2.5">
                  <Users size={16} className="text-blue-400" />
                  <span className="font-bold text-white">Seller Accounts Directory</span>
                </div>
                <ExternalLink size={14} className="text-zinc-500" />
              </Link>

              <Link
                href="/admin/transactions"
                className="flex items-center justify-between p-3 rounded-2xl bg-zinc-950/60 border border-zinc-800/80 hover:border-emerald-500/40 hover:bg-zinc-800/40 transition-all text-xs"
              >
                <div className="flex items-center gap-2.5">
                  <CheckCircle2 size={16} className="text-emerald-400" />
                  <span className="font-bold text-white">Marketplace Transactions</span>
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
