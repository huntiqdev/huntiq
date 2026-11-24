'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import {
  Server, Cpu, HardDrive, Activity, RefreshCw, Filter,
  CheckCircle2, AlertTriangle, Layers, Zap, Database,
  Globe2, ShieldCheck, ShieldAlert, FolderTree, ArrowRight,
  Sparkles, Radio, Lock, MapPin, Terminal, Monitor, Laptop
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';

export default function PlatformConsolePage() {
  const router = useRouter();
  const { user, isPlatformOwner, isSuperAdmin, isPlatformAdmin, adminRole, loading: authLoading } = useAuth();
  const [clusterFilter, setClusterFilter] = useState<'all' | 'addis' | 'bahrain' | 'edge'>('all');
  const [serviceFilter, setServiceFilter] = useState<'all' | 'api' | 'db' | 'redis' | 'search'>('all');
  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);
  const [masterDataStatus, setMasterDataStatus] = useState<any>(null);

  useEffect(() => {
    if (!authLoading) {
      const isAllowed = isPlatformOwner || isSuperAdmin || isPlatformAdmin ||
        adminRole === 'taxonomy_architect' || adminRole === 'gis_location_specialist';
      if (!user || !isAllowed) {
        router.replace('/platform-console/login');
      }
    }
  }, [user, authLoading, isPlatformOwner, isSuperAdmin, isPlatformAdmin, adminRole, router]);

  const fetchMasterDataStatus = () => {
    if (!user?.token) return;
    api.adminMasterDataStatus(user.token)
      .then((res: any) => setMasterDataStatus(res?.data || res))
      .catch(() => {});
  };

  useEffect(() => {
    fetchMasterDataStatus();
  }, [user?.token]);

  const handleRefresh = () => {
    setIsRefreshing(true);
    fetchMasterDataStatus();
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
              <option value="api">🚀 Go Fiber Core Gateway (Port 8080)</option>
              <option value="db">🐘 PostgreSQL Database Cluster (Port 5432)</option>
              <option value="redis">⚡ Redis Shared Cache</option>
              <option value="search">🔍 Search Engine Cluster</option>
            </select>

            <div className="px-3 py-1.5 rounded-xl bg-zinc-950/80 border border-zinc-800 text-zinc-400 text-xs font-mono">
              Postgres Active Pool: <span className="text-emerald-400 font-bold">25 / 100 Conn</span>
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
              852 <span className="text-xs font-normal text-zinc-400">Active Handlers</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-emerald-400 font-bold">
              <CheckCircle2 size={14} />
              <span>Go Fiber v2.52 • Zero Errors</span>
            </div>
          </div>

          <div className="bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-emerald-950/20 border border-zinc-800/90 rounded-3xl p-6 shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">PostgreSQL Schema Integrity</span>
              <div className="p-2.5 bg-emerald-500/10 text-emerald-400 rounded-2xl border border-emerald-500/20">
                <Database size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              80 <span className="text-xs font-normal text-zinc-400">Base Tables</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-emerald-400 font-bold">
              <CheckCircle2 size={14} />
              <span>Migration 067 Applied • Healthy</span>
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
              5+ <span className="text-xs font-normal text-zinc-400">Levels Deep</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-amber-400 font-bold">
              <CheckCircle2 size={14} />
              <span>Multi-Tier Clamped • 0 Orphaned</span>
            </div>
          </div>

          <div className="bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-sky-950/20 border border-zinc-800/90 rounded-3xl p-6 shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Operational Shifts & Workstations</span>
              <div className="p-2.5 bg-sky-500/10 text-sky-400 rounded-2xl border border-sky-500/20">
                <Monitor size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              13 <span className="text-xs font-normal text-zinc-400">Dedicated Terminals</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-sky-400 font-bold">
              <CheckCircle2 size={14} />
              <span>3 Shifts Active (24/7 Rotations)</span>
            </div>
          </div>

        </section>

        {/* 2. MASTER DATA & WORKSTATION CLUSTERS */}
        <section className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-zinc-900/80 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <FolderTree size={18} className="text-indigo-400" />
                <h3 className="text-base font-bold text-white">Master Data Synchronization</h3>
              </div>
              <span className="text-[10px] font-mono uppercase px-2 py-0.5 rounded bg-emerald-950/60 text-emerald-400 border border-emerald-500/30">
                DB-Driven Synced
              </span>
            </div>
            <p className="text-xs text-zinc-400 leading-relaxed">
              Automatic validation and live database synchronization for four-language taxonomy manifests, dynamic attribute schemas, and regional boundary hierarchies.
            </p>
            <div className="space-y-2 text-xs font-mono">
              <div className="p-3 bg-zinc-950 rounded-xl border border-zinc-800 flex justify-between">
                <span className="text-zinc-400">Total PostgreSQL Base Tables:</span>
                <span className="text-emerald-400 font-bold">80 Tables</span>
              </div>
              <div className="p-3 bg-zinc-950 rounded-xl border border-zinc-800 flex justify-between">
                <span className="text-zinc-400">Category Depth Clamping:</span>
                <span className="text-amber-400 font-bold">Levels 1 to 5+</span>
              </div>
              <div className="p-3 bg-zinc-950 rounded-xl border border-zinc-800 flex justify-between">
                <span className="text-zinc-400">Supported National Languages:</span>
                <span className="text-cyan-400 font-bold">EN • AM • OR • TI (4 Languages)</span>
              </div>
            </div>
          </div>

          <div className="bg-zinc-900/80 border border-zinc-800/80 rounded-3xl p-6 shadow-xl space-y-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Monitor size={18} className="text-emerald-400" />
                <h3 className="text-base font-bold text-white">Dedicated Workstation Fleet</h3>
              </div>
              <span className="text-[10px] font-mono uppercase px-2 py-0.5 rounded bg-blue-950/60 text-blue-400 border border-blue-500/30">
                13 Terminals
              </span>
            </div>
            <p className="text-xs text-zinc-400 leading-relaxed">
              Physical terminals in Addis Ababa Headquarters, Bole Regional Data Hub, Hawassa, and Dire Dawa. Hardware signatures bound to individual admin sessions.
            </p>
            <div className="space-y-2 text-xs font-mono">
              <div className="p-3 bg-zinc-950 rounded-xl border border-zinc-800 flex justify-between">
                <span className="text-zinc-400">Cybersecurity SOC Consoles:</span>
                <span className="text-cyan-400 font-bold">WS-SOC-01, WS-SOC-02</span>
              </div>
              <div className="p-3 bg-zinc-950 rounded-xl border border-zinc-800 flex justify-between">
                <span className="text-zinc-400">KYC & Verification Terminals:</span>
                <span className="text-blue-400 font-bold">WS-OPS-01, WS-OPS-02, WS-OPS-03</span>
              </div>
              <div className="p-3 bg-zinc-950 rounded-xl border border-zinc-800 flex justify-between">
                <span className="text-zinc-400">Treasury & Escrow Desks:</span>
                <span className="text-amber-400 font-bold">WS-FIN-01, WS-FIN-02</span>
              </div>
            </div>
          </div>
        </section>

      </main>
    </div>
  );
}
