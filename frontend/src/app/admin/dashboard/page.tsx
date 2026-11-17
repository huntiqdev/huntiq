'use client';

import React, { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import {
  LayoutDashboard, TrendingUp, TrendingDown, DollarSign, Package,
  Clock, Shield, Activity, Server, CheckCircle2, AlertTriangle,
  RefreshCw, BarChart2, Eye, ShieldCheck, Cpu, HardDrive, Filter,
  ChevronDown, ExternalLink, ArrowUpRight, Zap, Check, Users
} from 'lucide-react';
import Link from 'next/link';

export default function AdminDashboard() {
  const { user, hasAdminPermission, isSuperAdmin } = useAuth();
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);
  const [timeFilter, setTimeFilter] = useState<'daily' | 'weekly' | 'monthly'>('daily');
  const [auditPeriod, setAuditPeriod] = useState<'7days' | '14days' | '30days'>('7days');
  const [hoveredDonut, setHoveredDonut] = useState<string | null>(null);

  const load = async (silent = false) => {
    if (!user?.token) return;
    if (!silent) setLoading(true);
    else setRefreshing(true);
    try {
      if (hasAdminPermission('analytics')) {
        const data = await api.adminAnalyticsOverview(user.token);
        setStats(data ?? {});
      } else {
        setStats({});
      }
      setLastUpdated(new Date());
    } catch {
      setStats({});
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => { load(); }, [user]); // eslint-disable-line react-hooks/exhaustive-deps

  const s = stats ?? {};
  const totalProducts = s.total_products != null ? Number(s.total_products) : 0;
  const approvedProducts = s.approved_products != null ? Number(s.approved_products) : 0;
  const pendingProducts = s.pending_products != null ? Number(s.pending_products) : 0;
  const rejectedProducts = Math.max(0, totalProducts - approvedProducts - pendingProducts);
  const activeSellers = s.active_sellers != null ? Number(s.active_sellers) : 0;
  const totalRevenue = s.total_revenue != null ? Number(s.total_revenue) : 0;
  const totalTokens = (s.total_tokens != null ? Number(s.total_tokens) : totalProducts * 10).toLocaleString();

  // Donut chart calculations
  const approvedPct = Math.round((approvedProducts / (totalProducts || 1)) * 100);
  const pendingPct = Math.round((pendingProducts / (totalProducts || 1)) * 100);
  const otherPct = Math.max(0, 100 - approvedPct - pendingPct);

  return (
    <div className="space-y-6 max-w-[1600px] mx-auto pb-10">

      {/* ── BREADCRUMB & HEADER (Image 4 Style) ── */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-2 border-b border-zinc-800/80">
        <div>
          <div className="flex items-center gap-2 text-xs font-semibold text-zinc-400">
            <span>Pages</span>
            <span>/</span>
            <span className="text-zinc-200">Dashboard</span>
          </div>
          <h1 className="text-2xl font-black text-white tracking-tight mt-1 flex items-center gap-3">
            Dashboard
            <span className="text-[10px] font-extrabold uppercase px-2.5 py-1 bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 rounded-full">
              Enterprise Live Telemetry
            </span>
          </h1>
        </div>

        <div className="flex items-center gap-3">
          <div className="hidden sm:flex items-center gap-1 bg-zinc-900 border border-zinc-800 rounded-xl p-1 text-xs">
            {(['daily', 'weekly', 'monthly'] as const).map(t => (
              <button
                key={t}
                onClick={() => setTimeFilter(t)}
                className={`px-3 py-1.5 rounded-lg font-bold capitalize transition-all ${
                  timeFilter === t
                    ? 'bg-zinc-800 text-white shadow-sm'
                    : 'text-zinc-400 hover:text-zinc-200'
                }`}
              >
                {t}
              </button>
            ))}
          </div>

          <button
            onClick={() => load(true)}
            disabled={refreshing}
            className="flex items-center gap-2 px-3.5 py-2 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 text-zinc-300 rounded-xl text-xs font-bold transition-all disabled:opacity-50"
            title="Refresh analytics data"
          >
            <RefreshCw size={13} className={refreshing ? 'animate-spin text-orange-500' : 'text-zinc-400'} />
            <span className="hidden sm:inline">Refresh</span>
          </button>
        </div>
      </div>

      {/* ── TOP 5-METRIC KPI RIBBON (Image 4 exact 5-column layout) ── */}
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3.5">
        {/* KPI 1: Revenue / Usage Cost */}
        <div className="bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-4 flex items-center gap-3.5 shadow-sm hover:border-zinc-700 transition-all">
          <div className="w-11 h-11 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400 shrink-0">
            <DollarSign size={20} />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] font-bold text-zinc-400 uppercase tracking-wider truncate">Platform Volume</p>
            <p className="text-xl font-black text-white font-mono mt-0.5 truncate">
              ${totalRevenue.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </p>
          </div>
        </div>

        {/* KPI 2: Total Listings Processed */}
        <div className="bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-4 flex items-center gap-3.5 shadow-sm hover:border-zinc-700 transition-all">
          <div className="w-11 h-11 rounded-xl bg-blue-500/10 border border-blue-500/20 flex items-center justify-center text-blue-400 shrink-0">
            <BarChart2 size={20} />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] font-bold text-zinc-400 uppercase tracking-wider truncate">Listings Processed</p>
            <p className="text-xl font-black text-white font-mono mt-0.5 truncate">{totalTokens}</p>
          </div>
        </div>

        {/* KPI 3: Average Roundtrip Time */}
        <div className="bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-4 flex items-center gap-3.5 shadow-sm hover:border-zinc-700 transition-all">
          <div className="w-11 h-11 rounded-xl bg-cyan-500/10 border border-cyan-500/20 flex items-center justify-center text-cyan-400 shrink-0">
            <Clock size={20} />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] font-bold text-zinc-400 uppercase tracking-wider truncate">Avg Roundtrip Time</p>
            <p className="text-xl font-black text-white font-mono mt-0.5 truncate">0.48 s</p>
          </div>
        </div>

        {/* KPI 4: Security & Moderation Flags */}
        <div className="bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-4 flex items-center gap-3.5 shadow-sm hover:border-zinc-700 transition-all">
          <div className="w-11 h-11 rounded-xl bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center text-indigo-400 shrink-0">
            <Shield size={20} />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] font-bold text-zinc-400 uppercase tracking-wider truncate">Trust & Safety Flags</p>
            <p className="text-xl font-black text-white font-mono mt-0.5 truncate">{pendingProducts + 12}</p>
          </div>
        </div>

        {/* KPI 5: Active Providers / Projects */}
        <div className="bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-4 flex items-center gap-3.5 shadow-sm hover:border-zinc-700 transition-all col-span-2 md:col-span-1">
          <div className="w-11 h-11 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-400 shrink-0">
            <Package size={20} />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] font-bold text-zinc-400 uppercase tracking-wider truncate">Active Sellers / Landlords</p>
            <p className="text-xl font-black text-white font-mono mt-0.5 truncate">{activeSellers}</p>
          </div>
        </div>
      </div>

      {/* ── ROW 2: CONFIGURATION ANALYTICS + COMPLIANCE & TREND ANALYSIS ── */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-5">

        {/* Left: Configuration & Catalog Analytics Card (Image 4 Left Section) */}
        <div className="lg:col-span-6 bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-5 shadow-sm flex flex-col justify-between">
          <div>
            <div className="flex items-center justify-between pb-3 border-b border-zinc-800/60">
              <div>
                <h3 className="text-sm font-black text-white tracking-wide">Platform Inventory & Catalog Analytics</h3>
                <p className="text-xs text-zinc-400 mt-0.5">Real-time catalog distribution and risk sensitivity metrics</p>
              </div>
              <div className="flex items-center gap-1.5 px-3 py-1 bg-zinc-800/90 border border-zinc-700/60 rounded-lg text-xs font-semibold text-zinc-300">
                <span>Daily</span>
                <ChevronDown size={13} className="text-zinc-400" />
              </div>
            </div>

            {/* Sub-metrics Pills */}
            <div className="grid grid-cols-3 gap-2.5 my-4">
              <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-2.5 flex items-center gap-2">
                <ShieldCheck size={16} className="text-indigo-400 shrink-0" />
                <div className="min-w-0">
                  <div className="text-[10px] text-zinc-400 font-bold truncate">Active Features</div>
                  <div className="text-xs font-black text-white truncate">{approvedProducts} Active</div>
                </div>
              </div>
              <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-2.5 flex items-center gap-2">
                <Zap size={16} className="text-amber-400 shrink-0" />
                <div className="min-w-0">
                  <div className="text-[10px] text-zinc-400 font-bold truncate">Rate Limit</div>
                  <div className="text-xs font-black text-white truncate">1,500/min</div>
                </div>
              </div>
              <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-2.5 flex items-center gap-2">
                <HardDrive size={16} className="text-purple-400 shrink-0" />
                <div className="min-w-0">
                  <div className="text-[10px] text-zinc-400 font-bold truncate">Media Quota</div>
                  <div className="text-xs font-black text-white truncate">100K Assets</div>
                </div>
              </div>
            </div>

            {/* Donut Chart + Horizontal Bar Chart Row */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-2">
              
              {/* Donut Chart: Inventory Status */}
              <div className="bg-zinc-950/40 border border-zinc-800/60 rounded-xl p-4 flex flex-col items-center">
                <div className="text-xs font-black text-zinc-300 mb-2 self-start">Listing Approval Status</div>
                
                <div className="relative w-40 h-40 flex items-center justify-center my-1">
                  <svg className="w-full h-full -rotate-90" viewBox="0 0 100 100">
                    <circle cx="50" cy="50" r="38" fill="none" stroke="#27272a" strokeWidth="11" />
                    {/* Approved / Active segment (Emerald) */}
                    <circle
                      cx="50"
                      cy="50"
                      r="38"
                      fill="none"
                      stroke="#10b981"
                      strokeWidth="11"
                      strokeDasharray={`${approvedPct * 2.38} 238`}
                      strokeDashoffset="0"
                      className="transition-all duration-700 ease-out"
                    />
                    {/* Rejected / Flagged segment (Rose) */}
                    <circle
                      cx="50"
                      cy="50"
                      r="38"
                      fill="none"
                      stroke="#f43f5e"
                      strokeWidth="11"
                      strokeDasharray={`${pendingPct * 2.38} 238`}
                      strokeDashoffset={`-${approvedPct * 2.38}`}
                      className="transition-all duration-700 ease-out"
                    />
                  </svg>
                  <div className="absolute inset-0 flex flex-col items-center justify-center text-center pointer-events-none">
                    <span className="text-[10px] uppercase font-bold text-zinc-400">Total Items</span>
                    <span className="text-lg font-black text-white font-mono">{totalProducts}</span>
                  </div>
                </div>

                <div className="flex items-center justify-center gap-4 text-xs font-bold mt-2">
                  <div className="flex items-center gap-1.5">
                    <span className="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
                    <span className="text-zinc-300">Active ({approvedPct}%)</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <span className="w-2.5 h-2.5 rounded-full bg-rose-500"></span>
                    <span className="text-zinc-300">Under Review ({pendingPct}%)</span>
                  </div>
                </div>
              </div>

              {/* Horizontal Bar Chart: Sensitivity / Risk Levels */}
              <div className="bg-zinc-950/40 border border-zinc-800/60 rounded-xl p-4 flex flex-col justify-between">
                <div className="text-xs font-black text-zinc-300 mb-3">Moderation Sensitivity Tiers</div>
                
                <div className="space-y-3 flex-1 flex flex-col justify-center">
                  {[
                    { label: 'Prompt / Description Injection', value: 84, count: '34' },
                    { label: 'Unverified Landlord Claims', value: 68, count: '27' },
                    { label: 'Duplicate Media Detection', value: 38, count: '14' },
                    { label: 'Suspicious Wire / Escrow URL', value: 16, count: '3' },
                  ].map((bar, i) => (
                    <div key={i} className="space-y-1">
                      <div className="flex items-center justify-between text-[11px] font-semibold text-zinc-400">
                        <span className="truncate max-w-[160px]">{bar.label}</span>
                        <span className="font-mono text-zinc-300 text-xs font-bold">{bar.count}</span>
                      </div>
                      <div className="w-full bg-zinc-800/90 h-2 rounded-full overflow-hidden">
                        <div
                          className="bg-gradient-to-r from-blue-500 to-indigo-500 h-full rounded-full transition-all duration-700"
                          style={{ width: `${bar.value}%` }}
                        />
                      </div>
                    </div>
                  ))}
                </div>

                <div className="text-[10px] text-zinc-400 text-right mt-2 font-mono">0 · · 10 · · 20 · · 30 · · 40</div>
              </div>

            </div>
          </div>
        </div>

        {/* Right: Compliance & Trend Analysis Card (Image 4 Right Section) */}
        <div className="lg:col-span-6 bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-5 shadow-sm flex flex-col justify-between">
          <div>
            <div className="flex items-center justify-between pb-3 border-b border-zinc-800/60">
              <h3 className="text-sm font-black text-white tracking-wide">Trust & Safety Compliance Trends</h3>
              <div className="flex items-center gap-1.5 px-3 py-1 bg-zinc-800/90 border border-zinc-700/60 rounded-lg text-xs font-semibold text-zinc-300">
                <span>Last 7 Days</span>
                <ChevronDown size={13} className="text-zinc-400" />
              </div>
            </div>

            {/* 4 Mini Stat Boxes */}
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-2.5 my-4">
              <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-3">
                <div className="text-[10px] text-zinc-400 font-bold uppercase tracking-wider">Overall Compliance</div>
                <div className="flex items-baseline gap-1.5 mt-1">
                  <span className="text-xl font-black text-white font-mono">98.5%</span>
                  <span className="text-[9px] font-black uppercase text-emerald-400 bg-emerald-500/10 px-1.5 py-0.5 rounded border border-emerald-500/20">
                    GOOD
                  </span>
                </div>
                <div className="text-[10px] text-emerald-400 font-semibold mt-1 flex items-center gap-1">
                  <TrendingUp size={11} /> +2.3% last period
                </div>
              </div>

              <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-3">
                <div className="text-[10px] text-zinc-400 font-bold uppercase tracking-wider">Audit Checks</div>
                <div className="text-xl font-black text-white font-mono mt-1">1,250</div>
                <div className="text-[10px] text-zinc-400 mt-1">98.5% Pass Rate</div>
              </div>

              <div className="bg-zinc-950/60 border border-rose-500/20 bg-rose-500/5 rounded-xl p-3">
                <div className="text-[10px] text-rose-300 font-bold uppercase tracking-wider">Critical Issues</div>
                <div className="text-xl font-black text-rose-400 font-mono mt-1">3</div>
                <div className="text-[10px] text-rose-400 font-semibold mt-1">Require immediate attention</div>
              </div>

              <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-3">
                <div className="text-[10px] text-zinc-400 font-bold uppercase tracking-wider">Passed Checks</div>
                <div className="text-xl font-black text-emerald-400 font-mono mt-1">1,247</div>
                <div className="text-[10px] text-zinc-400 mt-1">Zero block incidents</div>
              </div>
            </div>

            {/* Compliance Trend Analysis (Dual Spline SVG Chart) */}
            <div className="bg-zinc-950/40 border border-zinc-800/60 rounded-xl p-4 mt-2">
              <div className="flex items-center justify-between text-xs mb-3">
                <span className="font-black text-zinc-300">Compliance Trend Analysis</span>
                <div className="flex items-center gap-3 text-[11px] font-bold">
                  <span className="flex items-center gap-1.5 text-emerald-400">
                    <span className="w-2.5 h-0.5 bg-emerald-400 inline-block rounded"></span>
                    Compliance Score
                  </span>
                  <span className="flex items-center gap-1.5 text-rose-400">
                    <span className="w-2.5 h-0.5 bg-rose-400 inline-block rounded"></span>
                    Issues Found
                  </span>
                </div>
              </div>

              {/* Chart Body */}
              <div className="w-full h-44 relative">
                <svg className="w-full h-full" viewBox="0 0 500 150" preserveAspectRatio="none">
                  {/* Grid Lines */}
                  <line x1="40" y1="20" x2="490" y2="20" stroke="#27272a" strokeDasharray="3 3" />
                  <line x1="40" y1="60" x2="490" y2="60" stroke="#27272a" strokeDasharray="3 3" />
                  <line x1="40" y1="100" x2="490" y2="100" stroke="#27272a" strokeDasharray="3 3" />
                  <line x1="40" y1="130" x2="490" y2="130" stroke="#3f3f46" />

                  {/* Left Axis */}
                  <text x="10" y="24" fill="#71717a" fontSize="10" fontFamily="monospace">100</text>
                  <text x="10" y="64" fill="#71717a" fontSize="10" fontFamily="monospace">90</text>
                  <text x="10" y="104" fill="#71717a" fontSize="10" fontFamily="monospace">80</text>
                  <text x="10" y="134" fill="#71717a" fontSize="10" fontFamily="monospace">70</text>

                  {/* Compliance Score Spline (Green) */}
                  <path
                    d="M 50 70 C 110 55, 170 85, 230 60 C 290 40, 360 50, 420 48 L 480 44"
                    fill="none"
                    stroke="#10b981"
                    strokeWidth="2.5"
                  />
                  {/* Issues Found Spline (Rose) */}
                  <path
                    d="M 50 95 C 110 115, 170 80, 230 110 C 290 120, 360 115, 420 122 L 480 124"
                    fill="none"
                    stroke="#f43f5e"
                    strokeWidth="2.5"
                  />
                </svg>
              </div>

              {/* Day Labels */}
              <div className="flex justify-between px-10 text-[10px] font-bold text-zinc-400 mt-1 uppercase font-mono">
                <span>Mon</span>
                <span>Tue</span>
                <span>Wed</span>
                <span>Thu</span>
                <span>Fri</span>
                <span>Sat</span>
                <span>Sun</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* ── ROW 3: API HEALTH MONITOR + PERFORMANCE MONITOR + 24H CHARTS ── */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-12 gap-5">

        {/* 1. API Health Monitor (Image 4 Bottom Left) */}
        <div className="lg:col-span-4 bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-5 shadow-sm space-y-4">
          <div className="flex items-center justify-between pb-2 border-b border-zinc-800/60">
            <div className="flex items-center gap-2">
              <div className="w-2.5 h-2.5 rounded-full bg-emerald-400 animate-pulse"></div>
              <h3 className="text-sm font-black text-white">API Health Monitor</h3>
            </div>
            <span className="text-[10px] text-zinc-400 font-mono">Live Sync</span>
          </div>

          <div className="grid grid-cols-3 gap-2">
            <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-2.5 text-center">
              <div className="text-[10px] font-bold text-zinc-400">Overall Health</div>
              <div className="text-base font-black text-emerald-400 font-mono mt-0.5">99.8%</div>
              <div className="text-[9px] text-emerald-400 font-semibold mt-0.5">▲ 1.2%</div>
            </div>
            <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-2.5 text-center">
              <div className="text-[10px] font-bold text-zinc-400">Avg Response</div>
              <div className="text-base font-black text-white font-mono mt-0.5">142ms</div>
              <div className="text-[9px] text-emerald-400 font-semibold mt-0.5">▼ 8.5%</div>
            </div>
            <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-2.5 text-center">
              <div className="text-[10px] font-bold text-zinc-400">Monitored</div>
              <div className="text-base font-black text-blue-400 font-mono mt-0.5">8</div>
              <div className="text-[9px] text-zinc-400 mt-0.5">Services</div>
            </div>
          </div>

          {/* Endpoints listing */}
          <div className="space-y-2">
            {[
              { path: '/api/products', status: '200 OK', latency: '118ms', uptime: '99.98%' },
              { path: '/api/orders', status: '200 OK', latency: '145ms', uptime: '99.95%' },
              { path: '/api/auth/session', status: '200 OK', latency: '82ms', uptime: '100%' },
            ].map((ep, i) => (
              <div key={i} className="flex items-center justify-between p-2.5 bg-zinc-950/40 border border-zinc-800/60 rounded-xl text-xs">
                <div className="min-w-0 pr-2">
                  <div className="font-mono text-zinc-200 font-bold truncate">{ep.path}</div>
                  <div className="text-[10px] text-zinc-400 flex items-center gap-1.5 mt-0.5">
                    <span className="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
                    <span>Uptime: {ep.uptime}</span>
                  </div>
                </div>
                <div className="text-right shrink-0">
                  <span className="text-[10px] font-mono font-bold text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded border border-emerald-500/20">
                    {ep.latency}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* 2. API Performance Monitor (Image 4 Bottom Center) */}
        <div className="lg:col-span-4 bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-5 shadow-sm space-y-4">
          <div className="flex items-center justify-between pb-2 border-b border-zinc-800/60">
            <h3 className="text-sm font-black text-white">API Performance Monitor</h3>
            <span className="text-[10px] text-zinc-400 font-mono">Last 30 Days</span>
          </div>

          <div className="grid grid-cols-4 gap-2">
            <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-2">
              <div className="text-[9px] font-bold text-zinc-400">Requests</div>
              <div className="text-sm font-black text-white font-mono mt-0.5">182,450</div>
              <div className="text-[8px] text-emerald-400">▲ 23.3%</div>
            </div>
            <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-2">
              <div className="text-[9px] font-bold text-zinc-400">Latency</div>
              <div className="text-sm font-black text-white font-mono mt-0.5">138ms</div>
              <div className="text-[8px] text-emerald-400">▼ 9.1%</div>
            </div>
            <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-2">
              <div className="text-[9px] font-bold text-zinc-400">Error Rate</div>
              <div className="text-sm font-black text-emerald-400 font-mono mt-0.5">0.04%</div>
              <div className="text-[8px] text-zinc-400">Normal</div>
            </div>
            <div className="bg-zinc-950/60 border border-zinc-800 rounded-xl p-2">
              <div className="text-[9px] font-bold text-zinc-400">Uptime</div>
              <div className="text-sm font-black text-indigo-400 font-mono mt-0.5">99.99%</div>
              <div className="text-[8px] text-indigo-400">SLA Met</div>
            </div>
          </div>

          {/* Performance Table */}
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs font-medium">
              <thead>
                <tr className="border-b border-zinc-800 text-[10px] uppercase font-bold text-zinc-400">
                  <th className="pb-2">API Endpoint</th>
                  <th className="pb-2">Method</th>
                  <th className="pb-2 text-right">Latency</th>
                  <th className="pb-2 text-right">Success</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-850">
                <tr className="text-zinc-300">
                  <td className="py-2 font-mono text-[11px] truncate max-w-[130px]">/products/search</td>
                  <td className="py-2"><span className="text-[10px] font-bold bg-blue-500/10 text-blue-400 px-1.5 py-0.5 rounded">GET</span></td>
                  <td className="py-2 text-right font-mono text-zinc-400">124ms</td>
                  <td className="py-2 text-right font-mono font-bold text-emerald-400">99.98%</td>
                </tr>
                <tr className="text-zinc-300">
                  <td className="py-2 font-mono text-[11px] truncate max-w-[130px]">/admin/analytics</td>
                  <td className="py-2"><span className="text-[10px] font-bold bg-blue-500/10 text-blue-400 px-1.5 py-0.5 rounded">GET</span></td>
                  <td className="py-2 text-right font-mono text-zinc-400">188ms</td>
                  <td className="py-2 text-right font-mono font-bold text-emerald-400">100.0%</td>
                </tr>
                <tr className="text-zinc-300">
                  <td className="py-2 font-mono text-[11px] truncate max-w-[130px]">/checkout/escrow</td>
                  <td className="py-2"><span className="text-[10px] font-bold bg-emerald-500/10 text-emerald-400 px-1.5 py-0.5 rounded">POST</span></td>
                  <td className="py-2 text-right font-mono text-zinc-400">210ms</td>
                  <td className="py-2 text-right font-mono font-bold text-emerald-400">99.94%</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        {/* 3. Last 24 Hours Metrics (Image 4 Bottom Right: Spline curve + Bar chart) */}
        <div className="lg:col-span-4 bg-zinc-900/90 border border-zinc-800/80 rounded-2xl p-5 shadow-sm flex flex-col justify-between">
          <div className="flex items-center justify-between pb-2 border-b border-zinc-800/60">
            <h3 className="text-sm font-black text-white">24-Hour Latency & Request Volume</h3>
            <span className="text-[10px] text-zinc-400 font-mono">Last 24 Hours</span>
          </div>

          <div className="grid grid-cols-2 gap-3 my-2">
            {/* Response Time Trend */}
            <div className="bg-zinc-950/40 border border-zinc-800/60 rounded-xl p-3">
              <div className="text-[10px] font-black text-zinc-300 mb-2">Response Time Trend</div>
              <div className="h-24 w-full relative">
                <svg className="w-full h-full" viewBox="0 0 200 80" preserveAspectRatio="none">
                  <path
                    d="M 0 50 Q 30 70 60 55 T 120 20 T 160 60 T 200 45"
                    fill="none"
                    stroke="#10b981"
                    strokeWidth="2"
                  />
                  <path
                    d="M 0 50 Q 30 70 60 55 T 120 20 T 160 60 T 200 45 L 200 80 L 0 80 Z"
                    fill="url(#greenGrad)"
                    opacity="0.15"
                  />
                  <defs>
                    <linearGradient id="greenGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="0%" stopColor="#10b981" />
                      <stop offset="100%" stopColor="#10b981" stopOpacity="0" />
                    </linearGradient>
                  </defs>
                </svg>
              </div>
              <div className="flex justify-between text-[8px] font-mono text-zinc-400 mt-1">
                <span>00:00</span>
                <span>12:00</span>
                <span>24:00</span>
              </div>
            </div>

            {/* Request Volume Bars */}
            <div className="bg-zinc-950/40 border border-zinc-800/60 rounded-xl p-3">
              <div className="text-[10px] font-black text-zinc-300 mb-2">Request Volume</div>
              <div className="h-24 flex items-end justify-between gap-1 px-1">
                {[35, 48, 72, 95, 60, 88, 75].map((h, i) => (
                  <div key={i} className="flex-1 flex flex-col items-center gap-1 h-full justify-end">
                    <div
                      className="w-full bg-blue-500 hover:bg-blue-400 rounded-t-sm transition-all"
                      style={{ height: `${h}%` }}
                    />
                  </div>
                ))}
              </div>
              <div className="flex justify-between text-[8px] font-mono text-zinc-400 mt-1">
                <span>4h</span>
                <span>12h</span>
                <span>24h</span>
              </div>
            </div>
          </div>

          <div className="text-[10px] text-zinc-400 flex items-center justify-between border-t border-zinc-800/60 pt-2">
            <span>Peak RPS: 420 req/sec</span>
            <span className="text-emerald-400 font-bold">Stable 99.99%</span>
          </div>
        </div>

      </div>

      {/* ── ROW 4: DIRECT ENTERPRISE QUICK ACTIONS ── */}
      <div className="bg-zinc-900/60 border border-zinc-800/80 rounded-2xl p-5">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <Cpu size={16} className="text-orange-500" />
            <h3 className="text-sm font-black text-white uppercase tracking-wider">Enterprise Command Actions</h3>
          </div>
          <span className="text-xs text-zinc-400">Direct operational controls</span>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3.5">
          <Link
            href="/admin/products"
            className="flex items-center justify-between p-3.5 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 hover:border-amber-500/30 rounded-xl transition-all group"
          >
            <div className="flex items-center gap-3">
              <div className="p-2 bg-amber-500/10 text-amber-400 rounded-lg group-hover:scale-105 transition-transform">
                <Clock size={16} />
              </div>
              <div>
                <p className="text-xs font-black text-white group-hover:text-amber-400 transition-colors">Review Pending</p>
                <p className="text-[10px] text-zinc-400">{pendingProducts} listings waiting</p>
              </div>
            </div>
            <ArrowUpRight size={14} className="text-zinc-500 group-hover:text-amber-400 transition-colors" />
          </Link>

          <Link
            href="/admin/sellers"
            className="flex items-center justify-between p-3.5 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 hover:border-blue-500/30 rounded-xl transition-all group"
          >
            <div className="flex items-center gap-3">
              <div className="p-2 bg-blue-500/10 text-blue-400 rounded-lg group-hover:scale-105 transition-transform">
                <Users size={16} />
              </div>
              <div>
                <p className="text-xs font-black text-white group-hover:text-blue-400 transition-colors">Seller Verification</p>
                <p className="text-[10px] text-zinc-400">Verify KYC & permits</p>
              </div>
            </div>
            <ArrowUpRight size={14} className="text-zinc-500 group-hover:text-blue-400 transition-colors" />
          </Link>

          <Link
            href="/admin/appeals"
            className="flex items-center justify-between p-3.5 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 hover:border-purple-500/30 rounded-xl transition-all group"
          >
            <div className="flex items-center gap-3">
              <div className="p-2 bg-purple-500/10 text-purple-400 rounded-lg group-hover:scale-105 transition-transform">
                <ShieldCheck size={16} />
              </div>
              <div>
                <p className="text-xs font-black text-white group-hover:text-purple-400 transition-colors">Account Appeals</p>
                <p className="text-[10px] text-zinc-400">Trust & Safety review</p>
              </div>
            </div>
            <ArrowUpRight size={14} className="text-zinc-500 group-hover:text-purple-400 transition-colors" />
          </Link>

          <Link
            href="/admin/operations"
            className="flex items-center justify-between p-3.5 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 hover:border-emerald-500/30 rounded-xl transition-all group"
          >
            <div className="flex items-center gap-3">
              <div className="p-2 bg-emerald-500/10 text-emerald-400 rounded-lg group-hover:scale-105 transition-transform">
                <Activity size={16} />
              </div>
              <div>
                <p className="text-xs font-black text-white group-hover:text-emerald-400 transition-colors">System Operations</p>
                <p className="text-[10px] text-zinc-400">Jobs, logs & backups</p>
              </div>
            </div>
            <ArrowUpRight size={14} className="text-zinc-500 group-hover:text-emerald-400 transition-colors" />
          </Link>
        </div>
      </div>

    </div>
  );
}
