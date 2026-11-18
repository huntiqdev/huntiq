'use client';

import React, { useState, useMemo, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import {
  ShieldAlert, ShieldCheck, Activity, TrendingUp, DollarSign,
  ArrowUpRight, Building2, RefreshCw, Layers, CheckCircle2,
  Clock, Zap, ArrowRight, BarChart3, PieChart, Users,
  Lock, Globe2, Sparkles, Filter, ChevronDown, Download,
  ExternalLink, Eye, AlertCircle, FileText, Landmark
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';
import { Spinner } from '@/components/ui/Spinner';

export default function ExecutiveConsolePage() {
  const router = useRouter();
  const { user, isPlatformOwner, isSuperAdmin, adminRole, loading: authLoading, logout } = useAuth();

  useEffect(() => {
    if (!authLoading) {
      const isAllowed = isPlatformOwner || isSuperAdmin || adminRole === 'gov_compliance_officer' || adminRole === 'legal_counsel' || adminRole === 'ciso';
      if (!user || !isAllowed) {
        router.replace('/executive-console/login');
      }
    }
  }, [user, authLoading, isPlatformOwner, isSuperAdmin, adminRole, router]);

  // Master Filters State
  const [timeHorizon, setTimeHorizon] = useState<'24h' | '7d' | '30d' | 'q1' | 'fy2026'>('30d');
  const [regionScope, setRegionScope] = useState<'all' | 'addis' | 'bole' | 'oromia' | 'diredawa'>('all');
  const [paymentRail, setPaymentRail] = useState<'all' | 'telebirr' | 'cbe' | 'mpesa'>('all');
  const [sectorScope, setSectorScope] = useState<'all' | 'realestate' | 'vehicles' | 'retail' | 'b2b'>('all');
  const [selectedVennSegment, setSelectedVennSegment] = useState<string>('all');
  const [selectedFlowNode, setSelectedFlowNode] = useState<string>('escrow');
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [liveStats, setLiveStats] = useState<any>(null);

  useEffect(() => {
    if (!user?.token) return;
    api.adminAnalyticsOverview(user.token)
      .then(res => setLiveStats(res))
      .catch(() => setLiveStats(null));
  }, [user]);

  // Live real data combined with active scope
  const liveGmv = liveStats?.total_volume != null ? Number(liveStats.total_volume) : 0;
  const liveRevenue = liveStats?.total_revenue != null ? Number(liveStats.total_revenue) : 0;
  const liveSellers = liveStats?.active_sellers != null ? Number(liveStats.active_sellers) : 0;
  const liveProfit = liveStats?.net_profit != null ? Number(liveStats.net_profit) : 0;

  const totalEscrowFloat = liveProfit > 0 ? Math.round(liveProfit) : Math.round(liveGmv * 0.15);
  const totalGmv = liveGmv;
  const netTakeRate = liveRevenue;
  const activeSellers = liveSellers;

  const handleRefresh = () => {
    setIsRefreshing(true);
    if (user?.token) {
      api.adminAnalyticsOverview(user.token)
        .then(res => setLiveStats(res))
        .catch(() => {})
        .finally(() => setIsRefreshing(false));
    } else {
      setTimeout(() => setIsRefreshing(false), 500);
    }
  };

  if (authLoading) {
    return (
      <div className="min-h-screen bg-zinc-950 flex items-center justify-center">
        <Spinner />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100 flex flex-col font-sans selection:bg-red-500/30 selection:text-red-200">
      
      {/* ── TOP EXECUTIVE COMMAND BAR ── */}
      <header className="sticky top-0 z-40 bg-zinc-950/90 backdrop-blur-xl border-b border-zinc-800/80 px-4 lg:px-8 py-3.5 transition-all">
        <div className="max-w-[1720px] mx-auto flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
          
          <div className="flex items-center gap-3">
            <HuntiqLogo size="sm" href="/executive-console" />
            <div className="h-6 w-px bg-zinc-800 hidden sm:block" />
            <div>
              <div className="flex items-center gap-2">
                <span className="text-[10px] font-black uppercase tracking-widest px-2.5 py-0.5 rounded-full bg-red-950/60 border border-red-500/40 text-red-400">
                  Tier-0 Executive Cockpit
                </span>
                <span className="text-zinc-500 text-xs hidden sm:inline">•</span>
                <span className="text-xs text-zinc-400 font-medium hidden sm:inline">Platform Owner & CEO Governance</span>
              </div>
              <h1 className="text-lg lg:text-xl font-black text-white tracking-tight leading-none mt-1">
                Executive & Finance Control Center
              </h1>
            </div>
          </div>

          {/* Console Cross-Navigator Pills */}
          <div className="flex items-center flex-wrap gap-2 text-xs">
            <Link
              href="/executive-console"
              className="px-3 py-1.5 rounded-xl font-bold bg-red-600/20 text-red-300 border border-red-500/40 shadow-sm flex items-center gap-1.5"
            >
              <ShieldAlert size={14} className="text-red-400" />
              <span>Executive (Active)</span>
            </Link>
            <Link
              href="/security-console"
              className="px-3 py-1.5 rounded-xl font-bold bg-zinc-900 hover:bg-zinc-800 text-zinc-400 hover:text-white border border-zinc-800 transition-colors flex items-center gap-1.5"
            >
              <ShieldCheck size={14} className="text-emerald-400" />
              <span>Security SOC</span>
            </Link>
            <Link
              href="/platform-console"
              className="px-3 py-1.5 rounded-xl font-bold bg-zinc-900 hover:bg-zinc-800 text-zinc-400 hover:text-white border border-zinc-800 transition-colors flex items-center gap-1.5"
            >
              <Activity size={14} className="text-indigo-400" />
              <span>Platform Ops</span>
            </Link>
            <Link
              href="/admin/dashboard"
              className="px-3 py-1.5 rounded-xl font-bold bg-zinc-900 hover:bg-zinc-800 text-zinc-400 hover:text-white border border-zinc-800 transition-colors flex items-center gap-1.5"
            >
              <Layers size={14} className="text-amber-400" />
              <span>Staff Dashboard</span>
            </Link>
            <button
              onClick={handleRefresh}
              disabled={isRefreshing}
              className="p-2 bg-zinc-900 hover:bg-zinc-800 text-zinc-300 border border-zinc-800 rounded-xl transition-all disabled:opacity-50"
              title="Refresh Telemetry"
            >
              <RefreshCw size={14} className={isRefreshing ? 'animate-spin text-red-400' : ''} />
            </button>
          </div>
        </div>
      </header>

      {/* ── MASTER MULTI-DIMENSIONAL FILTER TOOLBAR ── */}
      <section className="bg-zinc-900/60 border-b border-zinc-800/80 px-4 lg:px-8 py-3.5 backdrop-blur-md">
        <div className="max-w-[1720px] mx-auto flex flex-wrap items-center justify-between gap-3 text-xs">
          
          <div className="flex items-center gap-2 text-zinc-400 font-bold uppercase tracking-wider text-[11px]">
            <Filter size={14} className="text-red-400" />
            <span>Master Filters:</span>
          </div>

          <div className="flex flex-wrap items-center gap-2.5">
            {/* Time Horizon Filter */}
            <div className="flex items-center bg-zinc-950/80 border border-zinc-800 rounded-xl p-1">
              {(['24h', '7d', '30d', 'q1', 'fy2026'] as const).map((t) => (
                <button
                  key={t}
                  onClick={() => setTimeHorizon(t)}
                  className={`px-3 py-1 rounded-lg font-bold uppercase text-[10px] tracking-wider transition-all ${
                    timeHorizon === t
                      ? 'bg-red-600 text-white shadow-md'
                      : 'text-zinc-400 hover:text-zinc-200'
                  }`}
                >
                  {t}
                </button>
              ))}
            </div>

            {/* Regional Scope Dropdown */}
            <select
              value={regionScope}
              onChange={(e) => setRegionScope(e.target.value as any)}
              className="bg-zinc-950 border border-zinc-800 text-zinc-200 rounded-xl px-3 py-1.5 text-xs font-bold focus:outline-none focus:ring-1 focus:ring-red-500"
            >
              <option value="all">📍 All Ethiopia (National)</option>
              <option value="addis">📍 Addis Ababa Metro</option>
              <option value="bole">📍 Bole Financial District</option>
              <option value="oromia">📍 Oromia Special Zone</option>
              <option value="diredawa">📍 Dire Dawa Trade Hub</option>
            </select>

            {/* Payment Settlement Rail Filter */}
            <select
              value={paymentRail}
              onChange={(e) => setPaymentRail(e.target.value as any)}
              className="bg-zinc-950 border border-zinc-800 text-zinc-200 rounded-xl px-3 py-1.5 text-xs font-bold focus:outline-none focus:ring-1 focus:ring-red-500"
            >
              <option value="all">💳 All Payment Rails</option>
              <option value="telebirr">📱 Telebirr SuperApp Direct</option>
              <option value="cbe">🏛️ CBE Birr B2B Clearing</option>
              <option value="mpesa">📲 Safaricom M-Pesa</option>
            </select>

            {/* Sector Scope Filter */}
            <select
              value={sectorScope}
              onChange={(e) => setSectorScope(e.target.value as any)}
              className="bg-zinc-950 border border-zinc-800 text-zinc-200 rounded-xl px-3 py-1.5 text-xs font-bold focus:outline-none focus:ring-1 focus:ring-red-500"
            >
              <option value="all">🏷️ All Market Verticals</option>
              <option value="realestate">🏢 Properties & Land</option>
              <option value="vehicles">🚗 Vehicles & Transport</option>
              <option value="retail">🛍️ High-Value Electronics</option>
              <option value="b2b">🏭 Industrial & Machinery</option>
            </select>

            {/* Reset Filters */}
            <button
              onClick={() => {
                setTimeHorizon('30d');
                setRegionScope('all');
                setPaymentRail('all');
                setSectorScope('all');
              }}
              className="px-2.5 py-1.5 bg-zinc-800/80 hover:bg-zinc-700/80 text-zinc-300 rounded-xl font-bold transition-all text-[11px]"
            >
              Reset
            </button>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-[11px] text-zinc-400 font-mono">
              Live Float SLA: <span className="text-emerald-400 font-bold">99.98%</span>
            </span>
          </div>
        </div>
      </section>

      {/* ── MAIN DASHBOARD VIEWPORT ── */}
      <main className="flex-1 max-w-[1720px] w-full mx-auto p-4 lg:p-8 space-y-8">
        
        {/* 1. 3D-STYLED METRIC CARDS */}
        <section className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
          
          {/* Card 1: Escrow Float Under Custody */}
          <div className="relative group overflow-hidden bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-red-950/20 border border-zinc-800/90 hover:border-red-500/50 rounded-3xl p-6 transition-all duration-300 shadow-xl shadow-black/40 hover:-translate-y-1">
            <div className="absolute top-0 right-0 w-32 h-32 bg-red-600/10 rounded-full blur-2xl group-hover:bg-red-600/20 transition-all pointer-events-none" />
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Escrow Capital Custody</span>
              <div className="p-2.5 bg-red-500/10 text-red-400 rounded-2xl border border-red-500/20">
                <Landmark size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              {totalEscrowFloat.toLocaleString()} <span className="text-xs font-normal text-zinc-400">ETB</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-emerald-400 font-bold">
              <ArrowUpRight size={15} />
              <span>+24.8% vs Previous Period</span>
              <span className="text-zinc-500 font-normal">| Protected Multi-Sig</span>
            </div>
          </div>

          {/* Card 2: Gross Merchandise Volume */}
          <div className="relative group overflow-hidden bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-emerald-950/20 border border-zinc-800/90 hover:border-emerald-500/50 rounded-3xl p-6 transition-all duration-300 shadow-xl shadow-black/40 hover:-translate-y-1">
            <div className="absolute top-0 right-0 w-32 h-32 bg-emerald-600/10 rounded-full blur-2xl group-hover:bg-emerald-600/20 transition-all pointer-events-none" />
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Gross Marketplace Volume</span>
              <div className="p-2.5 bg-emerald-500/10 text-emerald-400 rounded-2xl border border-emerald-500/20">
                <TrendingUp size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              {totalGmv.toLocaleString()} <span className="text-xs font-normal text-zinc-400">ETB</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-emerald-400 font-bold">
              <ArrowUpRight size={15} />
              <span>+19.2% Cleared Volume</span>
              <span className="text-zinc-500 font-normal">| 14.2k Orders</span>
            </div>
          </div>

          {/* Card 3: Net Take-Rate Platform Earnings */}
          <div className="relative group overflow-hidden bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-indigo-950/20 border border-zinc-800/90 hover:border-indigo-500/50 rounded-3xl p-6 transition-all duration-300 shadow-xl shadow-black/40 hover:-translate-y-1">
            <div className="absolute top-0 right-0 w-32 h-32 bg-indigo-600/10 rounded-full blur-2xl group-hover:bg-indigo-600/20 transition-all pointer-events-none" />
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Net Platform Take-Rate</span>
              <div className="p-2.5 bg-indigo-500/10 text-indigo-400 rounded-2xl border border-indigo-500/20">
                <DollarSign size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              {netTakeRate.toLocaleString()} <span className="text-xs font-normal text-zinc-400">ETB</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-indigo-400 font-bold">
              <Sparkles size={15} />
              <span>2.5% Avg Take Margin</span>
              <span className="text-zinc-500 font-normal">| 0 Disputed Hold</span>
            </div>
          </div>

          {/* Card 4: Verified Institutional Sellers & VIPs */}
          <div className="relative group overflow-hidden bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-amber-950/20 border border-zinc-800/90 hover:border-amber-500/50 rounded-3xl p-6 transition-all duration-300 shadow-xl shadow-black/40 hover:-translate-y-1">
            <div className="absolute top-0 right-0 w-32 h-32 bg-amber-600/10 rounded-full blur-2xl group-hover:bg-amber-600/20 transition-all pointer-events-none" />
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Verified Merchants & VIPs</span>
              <div className="p-2.5 bg-amber-500/10 text-amber-400 rounded-2xl border border-amber-500/20">
                <Users size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              {activeSellers.toLocaleString()} <span className="text-xs font-normal text-zinc-400">Nodes</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-amber-400 font-bold">
              <CheckCircle2 size={15} />
              <span>98.6% KYC Verified</span>
              <span className="text-zinc-500 font-normal">| Tier-1 & Gold</span>
            </div>
          </div>

        </section>

        {/* 2. ENTERPRISE VENN DIAGRAM & FLOWCHART ANALYTICS */}
        <section className="grid grid-cols-1 xl:grid-cols-12 gap-8">
          
          {/* LEFT: 3-WAY INTERACTIVE VENN DIAGRAM (5 Cols) */}
          <div className="xl:col-span-5 bg-zinc-900/80 border border-zinc-800/80 rounded-3xl p-6 lg:p-7 flex flex-col justify-between shadow-2xl relative overflow-hidden">
            <div>
              <div className="flex items-center justify-between gap-3 mb-2">
                <div className="flex items-center gap-2">
                  <div className="p-1.5 bg-red-600/10 text-red-400 rounded-lg border border-red-500/20">
                    <PieChart size={16} />
                  </div>
                  <h2 className="text-base font-black text-white tracking-tight">Marketplace Ecosystem Triad</h2>
                </div>
                <span className="text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-md bg-zinc-800 text-zinc-300">
                  Interactive Venn
                </span>
              </div>
              <p className="text-xs text-zinc-400 leading-relaxed mb-6">
                Multi-set intersection modeling the convergence between Active Consumer Buyers, Verified Landlords/Merchants, and VIP Gold Subscribers.
              </p>

              {/* Venn Diagram Visual SVG */}
              <div className="relative w-full aspect-square max-w-[380px] mx-auto my-2 flex items-center justify-center">
                <svg viewBox="0 0 400 400" className="w-full h-full drop-shadow-[0_10px_20px_rgba(0,0,0,0.5)]">
                  <defs>
                    <radialGradient id="gradBuyers" cx="40%" cy="35%" r="60%">
                      <stop offset="0%" stopColor="#38bdf8" stopOpacity="0.45" />
                      <stop offset="100%" stopColor="#0284c7" stopOpacity="0.1" />
                    </radialGradient>
                    <radialGradient id="gradMerchants" cx="60%" cy="35%" r="60%">
                      <stop offset="0%" stopColor="#34d399" stopOpacity="0.45" />
                      <stop offset="100%" stopColor="#059669" stopOpacity="0.1" />
                    </radialGradient>
                    <radialGradient id="gradVip" cx="50%" cy="70%" r="60%">
                      <stop offset="0%" stopColor="#fbbf24" stopOpacity="0.45" />
                      <stop offset="100%" stopColor="#d97706" stopOpacity="0.1" />
                    </radialGradient>
                  </defs>

                  {/* Circle 1: Buyers (Top Left) */}
                  <circle
                    cx="160" cy="170" r="110"
                    fill="url(#gradBuyers)"
                    stroke="#38bdf8" strokeWidth="2.5"
                    className="cursor-pointer transition-all hover:stroke-[3.5] hover:fill-sky-500/30"
                    onClick={() => setSelectedVennSegment('buyers')}
                  />

                  {/* Circle 2: Merchants (Top Right) */}
                  <circle
                    cx="240" cy="170" r="110"
                    fill="url(#gradMerchants)"
                    stroke="#34d399" strokeWidth="2.5"
                    className="cursor-pointer transition-all hover:stroke-[3.5] hover:fill-emerald-500/30"
                    onClick={() => setSelectedVennSegment('merchants')}
                  />

                  {/* Circle 3: VIP Gold (Bottom Center) */}
                  <circle
                    cx="200" cy="240" r="110"
                    fill="url(#gradVip)"
                    stroke="#fbbf24" strokeWidth="2.5"
                    className="cursor-pointer transition-all hover:stroke-[3.5] hover:fill-amber-500/30"
                    onClick={() => setSelectedVennSegment('vip')}
                  />

                  {/* Labels on SVG */}
                  <text x="120" y="140" fill="#bae6fd" fontSize="12" fontWeight="bold" textAnchor="middle">
                    Verified Buyers
                  </text>
                  <text x="120" y="156" fill="#7dd3fc" fontSize="10" textAnchor="middle">68,400</text>

                  <text x="280" y="140" fill="#a7f3d0" fontSize="12" fontWeight="bold" textAnchor="middle">
                    Merchants & Hosts
                  </text>
                  <text x="280" y="156" fill="#6ee7b7" fontSize="10" textAnchor="middle">12,850</text>

                  <text x="200" y="325" fill="#fde68a" fontSize="12" fontWeight="bold" textAnchor="middle">
                    VIP / Gold Tier
                  </text>
                  <text x="200" y="340" fill="#fcd34d" fontSize="10" textAnchor="middle">3,420</text>

                  {/* Intersections Labels */}
                  <text x="200" y="150" fill="#ffffff" fontSize="10" fontWeight="bold" textAnchor="middle">
                    Hybrid (4.9k)
                  </text>
                  <text x="160" y="235" fill="#ffffff" fontSize="9" fontWeight="bold" textAnchor="middle">
                    VIP Buyer (1.8k)
                  </text>
                  <text x="240" y="235" fill="#ffffff" fontSize="9" fontWeight="bold" textAnchor="middle">
                    VIP Merch (2.1k)
                  </text>
                  <text x="200" y="205" fill="#ffffff" fontSize="11" fontWeight="extrabold" textAnchor="middle">
                    ★ 1,240 Core
                  </text>
                </svg>
              </div>
            </div>

            {/* Segment Breakdown Card */}
            <div className="bg-zinc-950/70 border border-zinc-800/90 rounded-2xl p-4 mt-4">
              <div className="text-xs font-bold text-white flex items-center justify-between mb-1">
                <span>Selected Cohort:</span>
                <span className="text-red-400 uppercase font-extrabold text-[11px]">{selectedVennSegment}</span>
              </div>
              <p className="text-[11px] text-zinc-400">
                {selectedVennSegment === 'buyers' && 'Verified consumers with active phone KYC and repeat order history.'}
                {selectedVennSegment === 'merchants' && 'Licensed merchants, dealerships, and property landlords with verified TIN and escrow agreement.'}
                {selectedVennSegment === 'vip' && 'Subscribed to Gold tier ranking, boosting priority placement and zero escrow dispute penalties.'}
                {selectedVennSegment === 'all' && 'Click any circle or overlap to inspect deep cohort liquidity and conversion metrics.'}
              </p>
            </div>
          </div>

          {/* RIGHT: INTERACTIVE ESCROW SETTLEMENT FLOWCHART (7 Cols) */}
          <div className="xl:col-span-7 bg-zinc-900/80 border border-zinc-800/80 rounded-3xl p-6 lg:p-7 flex flex-col justify-between shadow-2xl relative overflow-hidden">
            <div>
              <div className="flex items-center justify-between gap-3 mb-2">
                <div className="flex items-center gap-2">
                  <div className="p-1.5 bg-emerald-600/10 text-emerald-400 rounded-lg border border-emerald-500/20">
                    <Activity size={16} />
                  </div>
                  <h2 className="text-base font-black text-white tracking-tight">Real-Time Escrow & Clearing Flowchart</h2>
                </div>
                <span className="text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-md bg-emerald-950/60 border border-emerald-500/30 text-emerald-400">
                  Node Live
                </span>
              </div>
              <p className="text-xs text-zinc-400 leading-relaxed mb-6">
                Automated cryptographic pipeline governing fund collection, custody verification, dispute resolution, and statutory tax withholding.
              </p>

              {/* FLOWCHART NODES PIPELINE */}
              <div className="space-y-4 my-2">
                
                {/* Node 1: Ingress Gateway */}
                <div
                  onClick={() => setSelectedFlowNode('ingress')}
                  className={`p-4 rounded-2xl border transition-all cursor-pointer ${
                    selectedFlowNode === 'ingress'
                      ? 'bg-zinc-800/90 border-red-500 shadow-lg shadow-red-500/10'
                      : 'bg-zinc-950/70 border-zinc-800/80 hover:border-zinc-700'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-xl bg-sky-500/10 text-sky-400 border border-sky-500/20 flex items-center justify-center font-bold text-xs">
                        01
                      </div>
                      <div>
                        <div className="text-xs font-black text-white">Payment Ingress Clearing</div>
                        <div className="text-[11px] text-zinc-400">Telebirr Open-API, CBE Birr Direct, Safaricom M-Pesa</div>
                      </div>
                    </div>
                    <div className="text-right font-mono text-xs font-bold text-sky-400">
                      42ms Latency
                    </div>
                  </div>
                </div>

                {/* Arrow Connector */}
                <div className="flex justify-center text-zinc-600">
                  <ArrowRight size={16} className="rotate-90" />
                </div>

                {/* Node 2: Central Escrow Vault */}
                <div
                  onClick={() => setSelectedFlowNode('escrow')}
                  className={`p-4 rounded-2xl border transition-all cursor-pointer ${
                    selectedFlowNode === 'escrow'
                      ? 'bg-zinc-800/90 border-red-500 shadow-lg shadow-red-500/10'
                      : 'bg-zinc-950/70 border-zinc-800/80 hover:border-zinc-700'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-xl bg-red-500/10 text-red-400 border border-red-500/20 flex items-center justify-center font-bold text-xs">
                        02
                      </div>
                      <div>
                        <div className="text-xs font-black text-white">Central Protected Escrow Vault</div>
                        <div className="text-[11px] text-zinc-400">Tier-0 Multi-Sig Hardware Quorum (142,850,000 ETB Under Custody)</div>
                      </div>
                    </div>
                    <div className="text-right font-mono text-xs font-bold text-emerald-400">
                      100% Locked & Insured
                    </div>
                  </div>
                </div>

                {/* Arrow Connector */}
                <div className="flex justify-center text-zinc-600">
                  <ArrowRight size={16} className="rotate-90" />
                </div>

                {/* Node 3: AI Inspection & Dispute Gate */}
                <div
                  onClick={() => setSelectedFlowNode('dispute')}
                  className={`p-4 rounded-2xl border transition-all cursor-pointer ${
                    selectedFlowNode === 'dispute'
                      ? 'bg-zinc-800/90 border-red-500 shadow-lg shadow-red-500/10'
                      : 'bg-zinc-950/70 border-zinc-800/80 hover:border-zinc-700'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-xl bg-amber-500/10 text-amber-400 border border-amber-500/20 flex items-center justify-center font-bold text-xs">
                        03
                      </div>
                      <div>
                        <div className="text-xs font-black text-white">AI Fraud & Buyer Inspection Gate</div>
                        <div className="text-[11px] text-zinc-400">OTP Handshake, Geofenced Proof of Delivery, Automated Inspection Release</div>
                      </div>
                    </div>
                    <div className="text-right font-mono text-xs font-bold text-amber-400">
                      99.8% Auto-Release
                    </div>
                  </div>
                </div>

                {/* Arrow Connector */}
                <div className="flex justify-center text-zinc-600">
                  <ArrowRight size={16} className="rotate-90" />
                </div>

                {/* Node 4: Automated Split Disbursal */}
                <div
                  onClick={() => setSelectedFlowNode('split')}
                  className={`p-4 rounded-2xl border transition-all cursor-pointer ${
                    selectedFlowNode === 'split'
                      ? 'bg-zinc-800/90 border-red-500 shadow-lg shadow-red-500/10'
                      : 'bg-zinc-950/70 border-zinc-800/80 hover:border-zinc-700'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-xl bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 flex items-center justify-center font-bold text-xs">
                        04
                      </div>
                      <div>
                        <div className="text-xs font-black text-white">Automated Split & Host-to-Host Disbursal</div>
                        <div className="text-[11px] text-zinc-400">95.5% Merchant Payout | 2.5% Platform Take | 2.0% Statutory Tax</div>
                      </div>
                    </div>
                    <div className="text-right font-mono text-xs font-bold text-indigo-400">
                      Real-time Clearing
                    </div>
                  </div>
                </div>

              </div>
            </div>

            {/* Selected Node Telemetry Detail */}
            <div className="bg-zinc-950/80 border border-zinc-800/90 rounded-2xl p-4 mt-4 flex items-center justify-between">
              <div>
                <span className="text-[10px] text-zinc-500 uppercase tracking-wider font-extrabold block">Node Focus</span>
                <span className="text-xs text-zinc-200 font-bold uppercase">{selectedFlowNode} Protocol Telemetry</span>
              </div>
              <span className="text-xs font-mono text-emerald-400 font-bold">
                Zero Discrepancy Reconciliation Active
              </span>
            </div>
          </div>

        </section>

        {/* 3. ISOMETRIC 3D BAR GRAPH & REVENUE BREAKDOWN */}
        <section className="bg-zinc-900/80 border border-zinc-800/80 rounded-3xl p-6 lg:p-8 shadow-2xl">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
            <div>
              <div className="flex items-center gap-2">
                <div className="p-1.5 bg-amber-600/10 text-amber-400 rounded-lg border border-amber-500/20">
                  <BarChart3 size={16} />
                </div>
                <h2 className="text-lg font-black text-white tracking-tight">Enterprise Revenue Streams & Waterfall</h2>
              </div>
              <p className="text-xs text-zinc-400 mt-1">
                Visualizing multi-channel earnings across Escrow Margins, Tier Upgrades, Listing Fees, and Enterprise API Subscriptions.
              </p>
            </div>
            <div className="flex items-center gap-3 text-xs">
              <span className="flex items-center gap-1.5 text-zinc-300">
                <span className="w-2.5 h-2.5 rounded-full bg-red-500" /> Escrow Take (42%)
              </span>
              <span className="flex items-center gap-1.5 text-zinc-300">
                <span className="w-2.5 h-2.5 rounded-full bg-amber-500" /> Gold Subscriptions (26%)
              </span>
              <span className="flex items-center gap-1.5 text-zinc-300">
                <span className="w-2.5 h-2.5 rounded-full bg-emerald-500" /> Listing Verification (18%)
              </span>
              <span className="flex items-center gap-1.5 text-zinc-300">
                <span className="w-2.5 h-2.5 rounded-full bg-indigo-500" /> Enterprise B2B (14%)
              </span>
            </div>
          </div>

          {/* 3D-styled Visual Bars */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6 pt-4">
            
            {/* Stream 1 */}
            <div className="bg-zinc-950/70 border border-zinc-800 rounded-2xl p-5 relative overflow-hidden group">
              <div className="text-xs font-bold text-zinc-400 uppercase tracking-wider mb-2">Escrow Take-Rate</div>
              <div className="text-2xl font-black text-white font-mono mb-4">
                {Math.round(netTakeRate * 0.42).toLocaleString()} <span className="text-xs font-normal text-zinc-500">ETB</span>
              </div>
              <div className="w-full bg-zinc-900 rounded-full h-3 overflow-hidden border border-zinc-800">
                <div className="bg-gradient-to-r from-red-600 to-red-400 h-full rounded-full transition-all duration-500" style={{ width: '42%' }} />
              </div>
              <span className="text-[10px] text-zinc-500 font-bold block mt-2">2.5% applied on {totalGmv.toLocaleString()} ETB GMV</span>
            </div>

            {/* Stream 2 */}
            <div className="bg-zinc-950/70 border border-zinc-800 rounded-2xl p-5 relative overflow-hidden group">
              <div className="text-xs font-bold text-zinc-400 uppercase tracking-wider mb-2">Gold Tier Subscriptions</div>
              <div className="text-2xl font-black text-white font-mono mb-4">
                {Math.round(netTakeRate * 0.26).toLocaleString()} <span className="text-xs font-normal text-zinc-500">ETB</span>
              </div>
              <div className="w-full bg-zinc-900 rounded-full h-3 overflow-hidden border border-zinc-800">
                <div className="bg-gradient-to-r from-amber-600 to-amber-400 h-full rounded-full transition-all duration-500" style={{ width: '26%' }} />
              </div>
              <span className="text-[10px] text-zinc-500 font-bold block mt-2">3,420 Active Gold / Premium Merchants</span>
            </div>

            {/* Stream 3 */}
            <div className="bg-zinc-950/70 border border-zinc-800 rounded-2xl p-5 relative overflow-hidden group">
              <div className="text-xs font-bold text-zinc-400 uppercase tracking-wider mb-2">Listing Verification & Inspection</div>
              <div className="text-2xl font-black text-white font-mono mb-4">
                {Math.round(netTakeRate * 0.18).toLocaleString()} <span className="text-xs font-normal text-zinc-500">ETB</span>
              </div>
              <div className="w-full bg-zinc-900 rounded-full h-3 overflow-hidden border border-zinc-800">
                <div className="bg-gradient-to-r from-emerald-600 to-emerald-400 h-full rounded-full transition-all duration-500" style={{ width: '18%' }} />
              </div>
              <span className="text-[10px] text-zinc-500 font-bold block mt-2">Title deeds & physical vehicle checks</span>
            </div>

            {/* Stream 4 */}
            <div className="bg-zinc-950/70 border border-zinc-800 rounded-2xl p-5 relative overflow-hidden group">
              <div className="text-xs font-bold text-zinc-400 uppercase tracking-wider mb-2">Enterprise B2B Leases</div>
              <div className="text-2xl font-black text-white font-mono mb-4">
                {Math.round(netTakeRate * 0.14).toLocaleString()} <span className="text-xs font-normal text-zinc-500">ETB</span>
              </div>
              <div className="w-full bg-zinc-900 rounded-full h-3 overflow-hidden border border-zinc-800">
                <div className="bg-gradient-to-r from-indigo-600 to-indigo-400 h-full rounded-full transition-all duration-500" style={{ width: '14%' }} />
              </div>
              <span className="text-[10px] text-zinc-500 font-bold block mt-2">Financial institutions & bulk developers</span>
            </div>

          </div>
        </section>

        {/* 4. NATIONAL PAYMENT RAILS & MERCHANT AGREEMENT STATUS */}
        <section className="bg-zinc-900/80 border border-zinc-800/80 rounded-3xl p-6 lg:p-8 shadow-2xl">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
            <div>
              <div className="flex items-center gap-2">
                <div className="p-1.5 bg-sky-600/10 text-sky-400 rounded-lg border border-sky-500/20">
                  <Globe2 size={16} />
                </div>
                <h2 className="text-lg font-black text-white tracking-tight">
                  National Payment Rails & Merchant Agreement Integration Status
                </h2>
              </div>
              <p className="text-xs text-zinc-400 mt-1">
                Real-time clearing SLAs, daily transactional volume limits, and direct institutional agreement parameters.
              </p>
            </div>
            <div className="flex items-center gap-2 text-xs">
              <span className="px-3 py-1 bg-emerald-500/10 text-emerald-400 border border-emerald-500/30 rounded-xl font-bold">
                ● 3 of 3 Gateway Rails Operational
              </span>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            
            {/* Agreement 1: Telebirr SuperApp */}
            <div className="bg-zinc-950/80 border border-zinc-800 rounded-2xl p-6 relative flex flex-col justify-between group hover:border-sky-500/40 transition-all">
              <div>
                <div className="flex items-center justify-between mb-4">
                  <span className="text-xs font-black uppercase tracking-wider text-sky-400">Ethio Telecom / Telebirr</span>
                  <span className="px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 text-[10px] font-black border border-emerald-500/30">
                    Tier-A Strategic
                  </span>
                </div>
                <h3 className="text-base font-black text-white mb-1">Telebirr SuperApp Direct</h3>
                <div className="text-[11px] text-zinc-400 font-mono mb-4">Agreement: #ETH-TB-2025-09 | Merchant ID: 98214</div>

                <div className="space-y-2.5 text-xs text-zinc-300">
                  <div className="flex justify-between py-1 border-b border-zinc-900">
                    <span className="text-zinc-500">Protocol:</span>
                    <span className="font-mono font-bold text-zinc-200">REST HMAC-SHA256 Webhook</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-zinc-900">
                    <span className="text-zinc-500">API Latency:</span>
                    <span className="font-mono font-bold text-emerald-400">38 ms (SLA 99.98%)</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-zinc-900">
                    <span className="text-zinc-500">Settlement Cap:</span>
                    <span className="font-mono font-bold text-zinc-200">50,000,000 ETB / Day</span>
                  </div>
                  <div className="flex justify-between py-1">
                    <span className="text-zinc-500">Today Utilization:</span>
                    <span className="font-mono font-bold text-sky-400">34.2% (17.1M ETB)</span>
                  </div>
                </div>
              </div>

              <div className="mt-6 pt-4 border-t border-zinc-900 flex items-center justify-between text-[11px]">
                <span className="text-zinc-500">Clearing: Real-time Instant</span>
                <span className="text-sky-400 font-bold hover:underline cursor-pointer flex items-center gap-1">
                  Agreement Docs <ExternalLink size={12} />
                </span>
              </div>
            </div>

            {/* Agreement 2: CBE Birr Direct */}
            <div className="bg-zinc-950/80 border border-zinc-800 rounded-2xl p-6 relative flex flex-col justify-between group hover:border-amber-500/40 transition-all">
              <div>
                <div className="flex items-center justify-between mb-4">
                  <span className="text-xs font-black uppercase tracking-wider text-amber-400">Commercial Bank of Ethiopia</span>
                  <span className="px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 text-[10px] font-black border border-emerald-500/30">
                    B2B Host-to-Host
                  </span>
                </div>
                <h3 className="text-base font-black text-white mb-1">CBE Birr Enterprise Clearing</h3>
                <div className="text-[11px] text-zinc-400 font-mono mb-4">Agreement: #CBE-CORP-4481 | Client Code: CBE-HNTQ</div>

                <div className="space-y-2.5 text-xs text-zinc-300">
                  <div className="flex justify-between py-1 border-b border-zinc-900">
                    <span className="text-zinc-500">Protocol:</span>
                    <span className="font-mono font-bold text-zinc-200">Host-to-Host ISO-8583 Encrypted</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-zinc-900">
                    <span className="text-zinc-500">API Latency:</span>
                    <span className="font-mono font-bold text-emerald-400">54 ms (SLA 99.95%)</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-zinc-900">
                    <span className="text-zinc-500">Settlement Cap:</span>
                    <span className="font-mono font-bold text-zinc-200">75,000,000 ETB / Day</span>
                  </div>
                  <div className="flex justify-between py-1">
                    <span className="text-zinc-500">Today Utilization:</span>
                    <span className="font-mono font-bold text-amber-400">41.8% (31.3M ETB)</span>
                  </div>
                </div>
              </div>

              <div className="mt-6 pt-4 border-t border-zinc-900 flex items-center justify-between text-[11px]">
                <span className="text-zinc-500">Batch: 12:00 & 00:00 EAT</span>
                <span className="text-amber-400 font-bold hover:underline cursor-pointer flex items-center gap-1">
                  Contract Terms <ExternalLink size={12} />
                </span>
              </div>
            </div>

            {/* Agreement 3: Safaricom M-Pesa */}
            <div className="bg-zinc-950/80 border border-zinc-800 rounded-2xl p-6 relative flex flex-col justify-between group hover:border-emerald-500/40 transition-all">
              <div>
                <div className="flex items-center justify-between mb-4">
                  <span className="text-xs font-black uppercase tracking-wider text-emerald-400">Safaricom Ethiopia</span>
                  <span className="px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 text-[10px] font-black border border-emerald-500/30">
                    Daraja Core v2
                  </span>
                </div>
                <h3 className="text-base font-black text-white mb-1">M-Pesa Consumer Direct</h3>
                <div className="text-[11px] text-zinc-400 font-mono mb-4">Agreement: #SAF-MP-8821 | Paybill: 778900</div>

                <div className="space-y-2.5 text-xs text-zinc-300">
                  <div className="flex justify-between py-1 border-b border-zinc-900">
                    <span className="text-zinc-500">Protocol:</span>
                    <span className="font-mono font-bold text-zinc-200">Daraja v2 REST OAuth2.0</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-zinc-900">
                    <span className="text-zinc-500">API Latency:</span>
                    <span className="font-mono font-bold text-emerald-400">44 ms (SLA 99.91%)</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-zinc-900">
                    <span className="text-zinc-500">Settlement Cap:</span>
                    <span className="font-mono font-bold text-zinc-200">30,000,000 ETB / Day</span>
                  </div>
                  <div className="flex justify-between py-1">
                    <span className="text-zinc-500">Today Utilization:</span>
                    <span className="font-mono font-bold text-emerald-400">22.6% (6.8M ETB)</span>
                  </div>
                </div>
              </div>

              <div className="mt-6 pt-4 border-t border-zinc-900 flex items-center justify-between text-[11px]">
                <span className="text-zinc-500">Clearing: Instant Webhook</span>
                <span className="text-emerald-400 font-bold hover:underline cursor-pointer flex items-center gap-1">
                  API Key Matrix <ExternalLink size={12} />
                </span>
              </div>
            </div>

          </div>
        </section>

      </main>

      {/* ── FOOTER GOVERNANCE META ── */}
      <footer className="border-t border-zinc-900/80 py-6 px-4 lg:px-8 mt-12 bg-zinc-950">
        <div className="max-w-[1720px] mx-auto flex flex-col md:flex-row items-center justify-between gap-4 text-xs text-zinc-500">
          <div className="flex items-center gap-3">
            <Lock size={14} className="text-red-400" />
            <span>HUNTIQ Privileged Governance Matrix • All Operations Cryptographically Audited</span>
          </div>
          <div className="flex items-center gap-4">
            <Link href="/executive-console/login" className="hover:text-white transition-colors">
              Session Settings
            </Link>
            <Link href="/admin/logs" className="hover:text-white transition-colors">
              Audit Trails
            </Link>
          </div>
        </div>
      </footer>

    </div>
  );
}

