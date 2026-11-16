'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import {
  Server, Cpu, HardDrive, Activity, RefreshCw, Filter,
  CheckCircle2, AlertTriangle, Layers, Zap, Database,
  Globe2, ShieldCheck, ShieldAlert, FolderTree, ArrowRight,
  Sparkles, Radio, Lock
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';

export default function PlatformConsolePage() {
  const [clusterFilter, setClusterFilter] = useState<'all' | 'addis' | 'bahrain' | 'edge'>('all');
  const [serviceFilter, setServiceFilter] = useState<'all' | 'api' | 'db' | 'redis' | 'search'>('all');
  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);

  const handleRefresh = () => {
    setIsRefreshing(true);
    setTimeout(() => setIsRefreshing(false), 600);
  };

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100 flex flex-col font-sans selection:bg-indigo-500/30 selection:text-indigo-200">
      
      {/* ── TOP PLATFORM OPS HEADER ── */}
      <header className="sticky top-0 z-40 bg-zinc-950/90 backdrop-blur-xl border-b border-zinc-800/80 px-4 lg:px-8 py-3.5 transition-all">
        <div className="max-w-[1720px] mx-auto flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
          
          <div className="flex items-center gap-3">
            <HuntiqLogo size="sm" href="/platform-console" />
            <div className="h-6 w-px bg-zinc-800 hidden sm:block" />
            <div>
              <div className="flex items-center gap-2">
                <span className="text-[10px] font-black uppercase tracking-widest px-2.5 py-0.5 rounded-full bg-indigo-950/60 border border-indigo-500/40 text-indigo-400">
                  Tier-1 Platform Operations
                </span>
                <span className="text-zinc-500 text-xs hidden sm:inline">•</span>
                <span className="text-xs text-zinc-400 font-medium hidden sm:inline">Platform Administrator & Infrastructure Telemetry</span>
              </div>
              <h1 className="text-lg lg:text-xl font-black text-white tracking-tight leading-none mt-1">
                Platform Infrastructure & Catalog Governance
              </h1>
            </div>
          </div>

          {/* Console Cross-Navigator Pills */}
          <div className="flex items-center flex-wrap gap-2 text-xs">
            <Link
              href="/platform-console"
              className="px-3 py-1.5 rounded-xl font-bold bg-indigo-600/20 text-indigo-300 border border-indigo-500/40 shadow-sm flex items-center gap-1.5"
            >
              <Activity size={14} className="text-indigo-400" />
              <span>Platform Ops (Active)</span>
            </Link>
            <Link
              href="/executive-console"
              className="px-3 py-1.5 rounded-xl font-bold bg-zinc-900 hover:bg-zinc-800 text-zinc-400 hover:text-white border border-zinc-800 transition-colors flex items-center gap-1.5"
            >
              <ShieldAlert size={14} className="text-red-400" />
              <span>Executive / CEO</span>
            </Link>
            <Link
              href="/security-console"
              className="px-3 py-1.5 rounded-xl font-bold bg-zinc-900 hover:bg-zinc-800 text-zinc-400 hover:text-white border border-zinc-800 transition-colors flex items-center gap-1.5"
            >
              <ShieldCheck size={14} className="text-emerald-400" />
              <span>Security SOC</span>
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
              <RefreshCw size={14} className={isRefreshing ? 'animate-spin text-indigo-400' : ''} />
            </button>
          </div>
        </div>
      </header>

      {/* ── MASTER PLATFORM FILTER TOOLBAR ── */}
      <section className="bg-zinc-900/60 border-b border-zinc-800/80 px-4 lg:px-8 py-3.5 backdrop-blur-md">
        <div className="max-w-[1720px] mx-auto flex flex-wrap items-center justify-between gap-3 text-xs">
          
          <div className="flex items-center gap-2 text-zinc-400 font-bold uppercase tracking-wider text-[11px]">
            <Filter size={14} className="text-indigo-400" />
            <span>Infrastructure Scope:</span>
          </div>

          <div className="flex flex-wrap items-center gap-2.5">
            <select
              value={clusterFilter}
              onChange={(e) => setClusterFilter(e.target.value as any)}
              className="bg-zinc-950 border border-zinc-800 text-zinc-200 rounded-xl px-3 py-1.5 text-xs font-bold focus:outline-none focus:ring-1 focus:ring-indigo-500"
            >
              <option value="all">🌐 All Hybrid Node Clusters</option>
              <option value="addis">🇪🇹 Local Addis Tier-3 Facility</option>
              <option value="bahrain">☁️ AWS Middle East (Bahrain)</option>
              <option value="edge">⚡ Cloudflare Global Edge</option>
            </select>

            <select
              value={serviceFilter}
              onChange={(e) => setServiceFilter(e.target.value as any)}
              className="bg-zinc-950 border border-zinc-800 text-zinc-200 rounded-xl px-3 py-1.5 text-xs font-bold focus:outline-none focus:ring-1 focus:ring-indigo-500"
            >
              <option value="all">🧩 All Microservice Clusters</option>
              <option value="api">🚀 Go Fiber Core Gateway</option>
              <option value="db">🐘 PostgreSQL Database Cluster</option>
              <option value="redis">⚡ Redis Shared Cache</option>
              <option value="search">🔍 Search Engine Cluster</option>
            </select>

            <div className="px-3 py-1.5 rounded-xl bg-zinc-950/80 border border-zinc-800 text-zinc-400 text-xs font-mono">
              Postgres Active Pool: <span className="text-emerald-400 font-bold">32 / 100 Conn</span>
            </div>

            <div className="px-3 py-1.5 rounded-xl bg-zinc-950/80 border border-zinc-800 text-zinc-400 text-xs font-mono">
              Redis Hit Ratio: <span className="text-indigo-400 font-bold">98.4%</span>
            </div>
          </div>

          <div className="text-[11px] font-mono text-zinc-400">
            System Uptime: <span className="text-emerald-400 font-bold">99.99%</span>
          </div>
        </div>
      </section>

      {/* ── MAIN PLATFORM VIEWPORT ── */}
      <main className="flex-1 max-w-[1720px] w-full mx-auto p-4 lg:p-8 space-y-8">
        
        {/* 1. TOPOLOGY & HEALTH GAUGES */}
        <section className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
          
          <div className="bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-indigo-950/20 border border-zinc-800/90 rounded-3xl p-6 shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Core Go Gateway API</span>
              <div className="p-2.5 bg-indigo-500/10 text-indigo-400 rounded-2xl border border-indigo-500/20">
                <Cpu size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              24 <span className="text-xs font-normal text-zinc-400">ms P99</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-emerald-400 font-bold">
              <CheckCircle2 size={14} />
              <span>4 Pods Active • Zero Errors</span>
            </div>
          </div>

          <div className="bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-emerald-950/20 border border-zinc-800/90 rounded-3xl p-6 shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Database Connection Pool</span>
              <div className="p-2.5 bg-emerald-500/10 text-emerald-400 rounded-2xl border border-emerald-500/20">
                <Database size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              32 <span className="text-xs font-normal text-zinc-400">/ 100 Conn</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-emerald-400 font-bold">
              <CheckCircle2 size={14} />
              <span>0 Queued Waiters • 58 Tables</span>
            </div>
          </div>

          <div className="bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-amber-950/20 border border-zinc-800/90 rounded-3xl p-6 shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Category Hierarchy Integrity</span>
              <div className="p-2.5 bg-amber-500/10 text-amber-400 rounded-2xl border border-amber-500/20">
                <FolderTree size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              4 <span className="text-xs font-normal text-zinc-400">Levels Clamped</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-amber-400 font-bold">
              <CheckCircle2 size={14} />
              <span>1,280 Leaves • 0 Orphaned</span>
            </div>
          </div>

          <div className="bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-sky-950/20 border border-zinc-800/90 rounded-3xl p-6 shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Global CDN & Edge Delivery</span>
              <div className="p-2.5 bg-sky-500/10 text-sky-400 rounded-2xl border border-sky-500/20">
                <Globe2 size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              99.8% <span className="text-xs font-normal text-zinc-400">Edge Hit</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-sky-400 font-bold">
              <CheckCircle2 size={14} />
              <span>Addis PoP Active</span>
            </div>
          </div>

        </section>

        {/* 2. CLUSTER TOPOLOGY PIPELINE */}
        <section className="bg-zinc-900/80 border border-zinc-800/80 rounded-3xl p-6 lg:p-8 shadow-2xl">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-lg font-black text-white tracking-tight">Active Microservices Cluster Map</h2>
              <p className="text-xs text-zinc-400 mt-1">High-availability node distribution across localized and cloud infrastructure.</p>
            </div>
            <span className="px-3 py-1 rounded-xl bg-emerald-500/10 text-emerald-400 border border-emerald-500/30 text-xs font-bold">
              ● All Clusters Healthy
            </span>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="bg-zinc-950/80 border border-zinc-800 rounded-2xl p-5">
              <div className="text-xs font-bold text-zinc-400 uppercase tracking-wider mb-2">Tier-3 Addis Ababa PoP</div>
              <div className="text-lg font-black text-white mb-2">Core Escrow & Local Media Vault</div>
              <p className="text-xs text-zinc-500">Low-latency direct interconnect with Ethio Telecom and Commercial Bank of Ethiopia.</p>
            </div>
            <div className="bg-zinc-950/80 border border-zinc-800 rounded-2xl p-5">
              <div className="text-xs font-bold text-zinc-400 uppercase tracking-wider mb-2">AWS Middle East (Bahrain)</div>
              <div className="text-lg font-black text-white mb-2">Primary App Cluster & Database Standby</div>
              <p className="text-xs text-zinc-500">Auto-scaling Go Fiber workers, automated transaction rollback journal, backup replicas.</p>
            </div>
            <div className="bg-zinc-950/80 border border-zinc-800 rounded-2xl p-5">
              <div className="text-xs font-bold text-zinc-400 uppercase tracking-wider mb-2">Cloudflare Global Anycast</div>
              <div className="text-lg font-black text-white mb-2">WAF, DDoS Shield & Asset CDN</div>
              <p className="text-xs text-zinc-500">Zero-Trust access policies, WebAuthn edge challenges, static image cache.</p>
            </div>
          </div>
        </section>

      </main>

      {/* ── FOOTER ── */}
      <footer className="border-t border-zinc-900/80 py-6 px-4 lg:px-8 mt-12 bg-zinc-950">
        <div className="max-w-[1720px] mx-auto flex flex-col md:flex-row items-center justify-between gap-4 text-xs text-zinc-500">
          <div className="flex items-center gap-3">
            <Lock size={14} className="text-indigo-400" />
            <span>HUNTIQ Platform Governance • Distributed Systems Architecture</span>
          </div>
          <div className="flex items-center gap-4">
            <Link href="/platform-console/login" className="hover:text-white transition-colors">
              Platform Authentication
            </Link>
            <Link href="/admin/settings" className="hover:text-white transition-colors">
              System Settings
            </Link>
          </div>
        </div>
      </footer>

    </div>
  );
}
