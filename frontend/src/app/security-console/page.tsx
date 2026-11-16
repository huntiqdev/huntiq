'use client';

import React, { useState, useMemo } from 'react';
import Link from 'next/link';
import {
  ShieldAlert, ShieldCheck, Activity, AlertTriangle, Lock,
  RefreshCw, Filter, Globe2, Eye, Ban, Radio, Flame,
  Terminal, Server, UserX, CheckCircle2, XCircle, ArrowRight,
  SlidersHorizontal, Download, Power, Zap, Cpu, BellRing
} from 'lucide-react';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';

interface AttackEvent {
  id: string;
  time: string;
  ip: string;
  location: string;
  vector: string;
  endpoint: string;
  threatScore: number;
  severity: 'critical' | 'high' | 'warning' | 'info';
  actionTaken: string;
  autoMitigated: boolean;
}

const INITIAL_ATTACKS: AttackEvent[] = [
  {
    id: 'ATK-9812',
    time: 'Just now',
    ip: '185.220.101.44',
    location: 'Tor Exit Node / Netherlands',
    vector: 'Distributed Credential Stuffing',
    endpoint: 'POST /api/admin/login',
    threatScore: 98,
    severity: 'critical',
    actionTaken: 'Auto-Quarantine 24h & Session Drop',
    autoMitigated: true,
  },
  {
    id: 'ATK-9811',
    time: '42s ago',
    ip: '45.154.255.89',
    location: 'Hosting AS202425 / Frankfurt',
    vector: 'SQL Injection Union Probe',
    endpoint: "GET /api/products?cat=' OR 1=1--",
    threatScore: 94,
    severity: 'critical',
    actionTaken: 'WAF Rule #819 Triggered (TCP RST)',
    autoMitigated: true,
  },
  {
    id: 'ATK-9810',
    time: '1m ago',
    ip: '103.149.136.18',
    location: 'Residential Proxy / Singapore',
    vector: 'SMS / OTP Gate Flood Exhaustion',
    endpoint: 'POST /api/auth/otp/send',
    threatScore: 86,
    severity: 'high',
    actionTaken: 'Token Bucket Rate Limit Throttled',
    autoMitigated: true,
  },
  {
    id: 'ATK-9809',
    time: '3m ago',
    ip: '194.26.29.112',
    location: 'Known Botnet / Russia',
    vector: 'Admin Privilege Escalation Probe',
    endpoint: 'PUT /api/admin/users/12/role',
    threatScore: 99,
    severity: 'critical',
    actionTaken: 'Honeytoken Trap Sprung & IP Blackholed',
    autoMitigated: true,
  },
  {
    id: 'ATK-9808',
    time: '5m ago',
    ip: '82.102.23.41',
    location: 'VPN Gateway / London',
    vector: 'Cross-Site Scripting (Stored XSS)',
    endpoint: 'POST /api/seller/product/new',
    threatScore: 78,
    severity: 'high',
    actionTaken: 'Payload Stripped & Content Quarantined',
    autoMitigated: true,
  },
  {
    id: 'ATK-9807',
    time: '9m ago',
    ip: '196.188.12.90',
    location: 'Cellular EthioTelecom / Addis Ababa',
    vector: 'Excessive Failed Passkey Biometrics',
    endpoint: 'POST /api/admin/passkey/login/verify',
    threatScore: 62,
    severity: 'warning',
    actionTaken: 'Step-Up Verification Required',
    autoMitigated: false,
  },
];

export default function SecurityConsolePage() {
  const [attacks, setAttacks] = useState<AttackEvent[]>(INITIAL_ATTACKS);
  const [severityFilter, setSeverityFilter] = useState<'all' | 'critical' | 'high' | 'warning'>('all');
  const [vectorFilter, setVectorFilter] = useState<'all' | 'bruteforce' | 'sqli' | 'ddos' | 'privilege'>('all');
  const [autoDefenseMode, setAutoDefenseMode] = useState<boolean>(true);
  const [defconLevel, setDefconLevel] = useState<number>(2);
  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);

  // Manual IP Quarantine inputs
  const [manualIp, setManualIp] = useState('');
  const [manualDuration, setManualDuration] = useState('24h');
  const [manualReason, setManualReason] = useState('Manual SOC Threat Quarantine');
  const [blockedIps, setBlockedIps] = useState<string[]>([
    '185.220.101.44', '45.154.255.89', '194.26.29.112', '195.123.210.15'
  ]);

  // Killswitch state
  const [killswitchActive, setKillswitchActive] = useState<boolean>(false);
  const [killswitchModalOpen, setKillswitchModalOpen] = useState<boolean>(false);

  const filteredAttacks = useMemo(() => {
    return attacks.filter((a) => {
      if (severityFilter !== 'all' && a.severity !== severityFilter) return false;
      if (vectorFilter === 'bruteforce' && !a.vector.toLowerCase().includes('credential') && !a.vector.toLowerCase().includes('passkey')) return false;
      if (vectorFilter === 'sqli' && !a.vector.toLowerCase().includes('sql') && !a.vector.toLowerCase().includes('xss')) return false;
      if (vectorFilter === 'ddos' && !a.vector.toLowerCase().includes('flood')) return false;
      if (vectorFilter === 'privilege' && !a.vector.toLowerCase().includes('privilege')) return false;
      return true;
    });
  }, [attacks, severityFilter, vectorFilter]);

  const handleManualBan = (e: React.FormEvent) => {
    e.preventDefault();
    if (!manualIp.trim()) return;
    if (!blockedIps.includes(manualIp.trim())) {
      setBlockedIps([manualIp.trim(), ...blockedIps]);
    }
    const newAtk: AttackEvent = {
      id: `MAN-${Date.now().toString().slice(-4)}`,
      time: 'Just now',
      ip: manualIp.trim(),
      location: 'Manual SOC Operator Quarantine',
      vector: manualReason,
      endpoint: 'GLOBAL GATEWAY BLACKHOLE',
      threatScore: 100,
      severity: 'critical',
      actionTaken: `Manually Blocked (${manualDuration})`,
      autoMitigated: false,
    };
    setAttacks([newAtk, ...attacks]);
    setManualIp('');
  };

  const handleUnblock = (ipToUnblock: string) => {
    setBlockedIps(blockedIps.filter((ip) => ip !== ipToUnblock));
  };

  const handleRefresh = () => {
    setIsRefreshing(true);
    setTimeout(() => setIsRefreshing(false), 600);
  };

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100 flex flex-col font-sans selection:bg-emerald-500/30 selection:text-emerald-200">
      
      {/* ── TOP SECURITY SOC HEADER BAR ── */}
      <header className="sticky top-0 z-40 bg-zinc-950/90 backdrop-blur-xl border-b border-zinc-800/80 px-4 lg:px-8 py-3.5 transition-all">
        <div className="max-w-[1720px] mx-auto flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
          
          <div className="flex items-center gap-3">
            <HuntiqLogo size="sm" href="/security-console" />
            <div className="h-6 w-px bg-zinc-800 hidden sm:block" />
            <div>
              <div className="flex items-center gap-2">
                <span className="text-[10px] font-black uppercase tracking-widest px-2.5 py-0.5 rounded-full bg-emerald-950/60 border border-emerald-500/40 text-emerald-400">
                  Tier-1 Cyber Defense SOC
                </span>
                <span className="text-zinc-500 text-xs hidden sm:inline">•</span>
                <span className="text-xs text-zinc-400 font-medium hidden sm:inline">Security Administrator & Incident Response</span>
              </div>
              <h1 className="text-lg lg:text-xl font-black text-white tracking-tight leading-none mt-1">
                Security Operations & Threat Defense Console
              </h1>
            </div>
          </div>

          {/* Console Cross-Navigator Pills */}
          <div className="flex items-center flex-wrap gap-2 text-xs">
            <Link
              href="/security-console"
              className="px-3 py-1.5 rounded-xl font-bold bg-emerald-600/20 text-emerald-300 border border-emerald-500/40 shadow-sm flex items-center gap-1.5"
            >
              <ShieldCheck size={14} className="text-emerald-400" />
              <span>Security SOC (Active)</span>
            </Link>
            <Link
              href="/executive-console"
              className="px-3 py-1.5 rounded-xl font-bold bg-zinc-900 hover:bg-zinc-800 text-zinc-400 hover:text-white border border-zinc-800 transition-colors flex items-center gap-1.5"
            >
              <ShieldAlert size={14} className="text-red-400" />
              <span>Executive / CEO</span>
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
              <Radio size={14} className="text-amber-400" />
              <span>Staff Dashboard</span>
            </Link>
            <button
              onClick={handleRefresh}
              disabled={isRefreshing}
              className="p-2 bg-zinc-900 hover:bg-zinc-800 text-zinc-300 border border-zinc-800 rounded-xl transition-all disabled:opacity-50"
              title="Refresh Telemetry"
            >
              <RefreshCw size={14} className={isRefreshing ? 'animate-spin text-emerald-400' : ''} />
            </button>
          </div>
        </div>
      </header>

      {/* ── MASTER SECURITY FILTER TOOLBAR ── */}
      <section className="bg-zinc-900/60 border-b border-zinc-800/80 px-4 lg:px-8 py-3.5 backdrop-blur-md">
        <div className="max-w-[1720px] mx-auto flex flex-wrap items-center justify-between gap-3 text-xs">
          
          <div className="flex items-center gap-2 text-zinc-400 font-bold uppercase tracking-wider text-[11px]">
            <Filter size={14} className="text-emerald-400" />
            <span>SOC Master Filters:</span>
          </div>

          <div className="flex flex-wrap items-center gap-2.5">
            {/* Threat Severity Filter */}
            <div className="flex items-center bg-zinc-950/80 border border-zinc-800 rounded-xl p-1">
              {(['all', 'critical', 'high', 'warning'] as const).map((s) => (
                <button
                  key={s}
                  onClick={() => setSeverityFilter(s)}
                  className={`px-3 py-1 rounded-lg font-bold uppercase text-[10px] tracking-wider transition-all ${
                    severityFilter === s
                      ? 'bg-emerald-600 text-white shadow-md'
                      : 'text-zinc-400 hover:text-zinc-200'
                  }`}
                >
                  {s}
                </button>
              ))}
            </div>

            {/* Attack Vector Filter */}
            <select
              value={vectorFilter}
              onChange={(e) => setVectorFilter(e.target.value as any)}
              className="bg-zinc-950 border border-zinc-800 text-zinc-200 rounded-xl px-3 py-1.5 text-xs font-bold focus:outline-none focus:ring-1 focus:ring-emerald-500"
            >
              <option value="all">🛡️ All Threat Vectors</option>
              <option value="bruteforce">🔑 Credential Stuffing & Passkeys</option>
              <option value="sqli">💉 SQL Injection & XSS</option>
              <option value="ddos">🌊 Rate-Limit Flood & DDoS</option>
              <option value="privilege">🚨 Privilege Escalation</option>
            </select>

            {/* Auto-Mitigation Engine Toggle */}
            <button
              onClick={() => setAutoDefenseMode(!autoDefenseMode)}
              className={`flex items-center gap-2 px-3 py-1.5 rounded-xl text-xs font-bold transition-all border ${
                autoDefenseMode
                  ? 'bg-emerald-950/60 border-emerald-500/50 text-emerald-300'
                  : 'bg-amber-950/60 border-amber-500/50 text-amber-300'
              }`}
            >
              <Zap size={14} />
              <span>Auto-Defense: {autoDefenseMode ? 'ACTIVE (Autonomous AI Drop)' : 'MANUAL APPROVAL ONLY'}</span>
            </button>

            {/* Emergency Killswitch Action Button */}
            <button
              onClick={() => setKillswitchModalOpen(true)}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-black transition-all border ${
                killswitchActive
                  ? 'bg-red-600 text-white border-red-500 animate-pulse'
                  : 'bg-red-950/40 hover:bg-red-900/60 text-red-400 border-red-500/40'
              }`}
            >
              <Power size={14} />
              <span>{killswitchActive ? 'KILLSWITCH ENGAGED' : 'SYSTEM KILLSWITCH'}</span>
            </button>
          </div>

          <div className="flex items-center gap-2 font-mono text-[11px] text-zinc-400">
            <span>Threat Radar: <span className="text-emerald-400 font-bold">Scanning Ingress</span></span>
          </div>
        </div>
      </section>

      {/* ── MAIN SOC VIEWPORT ── */}
      <main className="flex-1 max-w-[1720px] w-full mx-auto p-4 lg:p-8 space-y-8">
        
        {/* 1. DEFCON & CYBER THREAT RADAR CARDS */}
        <section className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
          
          {/* Card 1: DEFCON Status */}
          <div className="bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-emerald-950/20 border border-zinc-800/90 rounded-3xl p-6 relative overflow-hidden shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Current Threat Status</span>
              <div className="p-2.5 bg-emerald-500/10 text-emerald-400 rounded-2xl border border-emerald-500/20">
                <Radio size={20} className="animate-pulse" />
              </div>
            </div>
            <div className="text-3xl font-black text-emerald-400 tracking-tight mb-2 font-mono flex items-center gap-2">
              <span>DEFCON {defconLevel}</span>
              <span className="text-xs px-2 py-0.5 rounded-full bg-emerald-500/20 text-emerald-300 font-sans font-bold">
                Elevated Shield
              </span>
            </div>
            <div className="flex items-center gap-2 text-xs text-zinc-400 font-medium">
              <span>Autonomous IPS Active</span>
              <span className="text-zinc-600">•</span>
              <span className="text-emerald-400 font-bold">0 Active Breaches</span>
            </div>
          </div>

          {/* Card 2: Intercepted Ingress Probes */}
          <div className="bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-red-950/20 border border-zinc-800/90 rounded-3xl p-6 relative overflow-hidden shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Attacks Intercepted (24h)</span>
              <div className="p-2.5 bg-red-500/10 text-red-400 rounded-2xl border border-red-500/20">
                <Flame size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              14,892 <span className="text-xs font-normal text-zinc-400">Probes</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-red-400 font-bold">
              <span>99.98% Dropped at Edge</span>
              <span className="text-zinc-600">•</span>
              <span className="text-zinc-400 font-normal">WAF + AI Shield</span>
            </div>
          </div>

          {/* Card 3: Active Quarantined IPs */}
          <div className="bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-amber-950/20 border border-zinc-800/90 rounded-3xl p-6 relative overflow-hidden shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Active Blackhole IPs</span>
              <div className="p-2.5 bg-amber-500/10 text-amber-400 rounded-2xl border border-amber-500/20">
                <Ban size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              {blockedIps.length + 1424} <span className="text-xs font-normal text-zinc-400">IPs</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-amber-400 font-bold">
              <span>Automatic BGP Sync</span>
              <span className="text-zinc-600">•</span>
              <span className="text-zinc-400 font-normal">Tor & Botnet Feeds</span>
            </div>
          </div>

          {/* Card 4: Honeytoken Decoy Traps */}
          <div className="bg-gradient-to-br from-zinc-900/90 via-zinc-900/60 to-indigo-950/20 border border-zinc-800/90 rounded-3xl p-6 relative overflow-hidden shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <span className="text-xs font-black uppercase tracking-wider text-zinc-400">Honeytoken Traps</span>
              <div className="p-2.5 bg-indigo-500/10 text-indigo-400 rounded-2xl border border-indigo-500/20">
                <Cpu size={20} />
              </div>
            </div>
            <div className="text-3xl font-black text-white tracking-tight mb-2 font-mono">
              43 <span className="text-xs font-normal text-zinc-400">Armed Traps</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-indigo-400 font-bold">
              <span>Zero False-Positive</span>
              <span className="text-zinc-600">•</span>
              <span className="text-zinc-400 font-normal">Decoy Admin Endpoints</span>
            </div>
          </div>

        </section>

        {/* 2. REAL-TIME CYBER ATTACK AUTO-DETECTION FEED & DEFENSE PIPELINE */}
        <section className="grid grid-cols-1 xl:grid-cols-12 gap-8">
          
          {/* LEFT: LIVE ATTACK AUTO-DETECTION STREAM (8 Cols) */}
          <div className="xl:col-span-8 bg-zinc-900/80 border border-zinc-800/80 rounded-3xl p-6 lg:p-7 shadow-2xl relative overflow-hidden flex flex-col justify-between">
            <div>
              <div className="flex items-center justify-between gap-3 mb-2">
                <div className="flex items-center gap-2">
                  <div className="p-1.5 bg-red-600/10 text-red-400 rounded-lg border border-red-500/20">
                    <Terminal size={16} />
                  </div>
                  <h2 className="text-base font-black text-white tracking-tight">Live Cyber Attack Auto-Detection Feed</h2>
                </div>
                <div className="flex items-center gap-2">
                  <span className="inline-block w-2 h-2 rounded-full bg-emerald-400 animate-ping" />
                  <span className="text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-md bg-zinc-800 text-zinc-300">
                    Real-Time Stream
                  </span>
                </div>
              </div>
              <p className="text-xs text-zinc-400 leading-relaxed mb-6">
                Intercepted hacker probes, credential stuffing surges, SQL injections, and unauthorized privilege escalation events.
              </p>

              {/* STREAM TABLE */}
              <div className="overflow-x-auto">
                <table className="w-full text-left text-xs">
                  <thead>
                    <tr className="border-b border-zinc-800 text-[10px] font-black uppercase tracking-wider text-zinc-500">
                      <th className="pb-3">Timestamp</th>
                      <th className="pb-3">Attacker IP & Origin</th>
                      <th className="pb-3">Threat Vector</th>
                      <th className="pb-3">Target Endpoint</th>
                      <th className="pb-3 text-center">Score</th>
                      <th className="pb-3 text-right">Defense Mitigation</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-zinc-800/60 font-mono">
                    {filteredAttacks.map((atk) => (
                      <tr key={atk.id} className="hover:bg-zinc-800/30 transition-colors">
                        <td className="py-3 text-zinc-400 text-[11px] whitespace-nowrap">
                          {atk.time}
                        </td>
                        <td className="py-3">
                          <div className="font-bold text-white text-xs">{atk.ip}</div>
                          <div className="text-[10px] text-zinc-500 font-sans">{atk.location}</div>
                        </td>
                        <td className="py-3 font-sans">
                          <span className={`inline-block px-2 py-0.5 rounded-md text-[10px] font-bold ${
                            atk.severity === 'critical'
                              ? 'bg-red-950/60 text-red-400 border border-red-500/30'
                              : atk.severity === 'high'
                              ? 'bg-amber-950/60 text-amber-400 border border-amber-500/30'
                              : 'bg-zinc-800 text-zinc-300'
                          }`}>
                            {atk.vector}
                          </span>
                        </td>
                        <td className="py-3 text-zinc-300 text-[11px] max-w-[200px] truncate" title={atk.endpoint}>
                          {atk.endpoint}
                        </td>
                        <td className="py-3 text-center">
                          <span className={`font-black text-xs ${
                            atk.threatScore >= 90 ? 'text-red-400' : atk.threatScore >= 75 ? 'text-amber-400' : 'text-emerald-400'
                          }`}>
                            {atk.threatScore}
                          </span>
                        </td>
                        <td className="py-3 text-right whitespace-nowrap font-sans">
                          <span className="text-[11px] font-bold text-emerald-400 flex items-center justify-end gap-1">
                            <CheckCircle2 size={13} />
                            {atk.actionTaken}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>

            <div className="mt-6 pt-4 border-t border-zinc-800 flex items-center justify-between text-[11px] text-zinc-500">
              <span>Showing {filteredAttacks.length} active intrusion alerts</span>
              <span className="text-zinc-400 font-bold">Autonomous Incident Dispatch: Enabled</span>
            </div>
          </div>

          {/* RIGHT: DEFENSE FLOWCHART & PIPELINE (4 Cols) */}
          <div className="xl:col-span-4 bg-zinc-900/80 border border-zinc-800/80 rounded-3xl p-6 lg:p-7 shadow-2xl relative overflow-hidden flex flex-col justify-between">
            <div>
              <div className="flex items-center justify-between gap-3 mb-2">
                <div className="flex items-center gap-2">
                  <div className="p-1.5 bg-emerald-600/10 text-emerald-400 rounded-lg border border-emerald-500/20">
                    <Activity size={16} />
                  </div>
                  <h2 className="text-base font-black text-white tracking-tight">Cyber Defense Flowchart</h2>
                </div>
                <span className="text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-md bg-emerald-950/60 border border-emerald-500/30 text-emerald-400">
                  Active Filter
                </span>
              </div>
              <p className="text-xs text-zinc-400 leading-relaxed mb-6">
                Packet filtration flow through Cloudflare WAF, Ingress ML Scorer, Token-Bucket Shaper, and Blackhole routing.
              </p>

              {/* FLOWCHART STEPS */}
              <div className="space-y-3.5 my-2 text-xs">
                
                <div className="bg-zinc-950/80 border border-zinc-800 rounded-2xl p-3.5">
                  <div className="flex items-center justify-between mb-1">
                    <span className="font-black text-white">01. Public Edge Ingress</span>
                    <span className="text-sky-400 font-mono text-[10px]">Cloudflare TLS 1.3</span>
                  </div>
                  <p className="text-[11px] text-zinc-400">Geo-IP restriction, known bad ASN filtering, and TCP SYN flood mitigation.</p>
                </div>

                <div className="flex justify-center text-zinc-600">
                  <ArrowRight size={14} className="rotate-90" />
                </div>

                <div className="bg-zinc-950/80 border border-zinc-800 rounded-2xl p-3.5">
                  <div className="flex items-center justify-between mb-1">
                    <span className="font-black text-white">02. AI Threat & Heuristic Engine</span>
                    <span className="text-emerald-400 font-mono text-[10px]">Anomaly Score &lt; 20</span>
                  </div>
                  <p className="text-[11px] text-zinc-400">Deep packet inspection, SQLi/XSS AST tokenizer, credential spray detection.</p>
                </div>

                <div className="flex justify-center text-zinc-600">
                  <ArrowRight size={14} className="rotate-90" />
                </div>

                <div className="bg-zinc-950/80 border border-zinc-800 rounded-2xl p-3.5">
                  <div className="flex items-center justify-between mb-1">
                    <span className="font-black text-white">03. Token-Bucket Rate Limiter</span>
                    <span className="text-amber-400 font-mono text-[10px]">Granular IP/User</span>
                  </div>
                  <p className="text-[11px] text-zinc-400">Burst limit enforcement (5 req/sec OTP, 10 req/sec Auth, 50 req/sec Read).</p>
                </div>

                <div className="flex justify-center text-zinc-600">
                  <ArrowRight size={14} className="rotate-90" />
                </div>

                <div className="bg-zinc-950/80 border border-emerald-500/40 rounded-2xl p-3.5 bg-emerald-950/10">
                  <div className="flex items-center justify-between mb-1">
                    <span className="font-black text-emerald-300">04. Zero-Trust Core Cluster</span>
                    <span className="text-emerald-400 font-mono text-[10px]">Clean Forwarding</span>
                  </div>
                  <p className="text-[11px] text-zinc-400">Passed traffic routed to Go Fiber backend and PostgreSQL database.</p>
                </div>

              </div>
            </div>

            <div className="bg-zinc-950/90 border border-zinc-800 rounded-2xl p-3.5 mt-4">
              <span className="text-[10px] text-zinc-500 uppercase tracking-wider font-extrabold block">Cluster Health</span>
              <span className="text-xs text-white font-bold">100% Ingress Inspection Coverage</span>
            </div>
          </div>

        </section>

        {/* 3. FULL CONTROL SYSTEM: MANUAL IP QUARANTINE & ACTIVE BLOCKLIST */}
        <section className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          
          {/* MANUAL BAN CONTROLLER (5 Cols) */}
          <div className="lg:col-span-5 bg-zinc-900/80 border border-zinc-800/80 rounded-3xl p-6 lg:p-7 shadow-2xl flex flex-col justify-between">
            <div>
              <div className="flex items-center gap-2 mb-2">
                <div className="p-1.5 bg-red-600/10 text-red-400 rounded-lg border border-red-500/20">
                  <Ban size={16} />
                </div>
                <h2 className="text-base font-black text-white tracking-tight">Manual Threat Quarantine & Ban Tool</h2>
              </div>
              <p className="text-xs text-zinc-400 leading-relaxed mb-6">
                Instantly inject IP addresses or CIDR ranges into the edge BGP blackhole filter to mitigate ongoing targeted attacks.
              </p>

              <form onSubmit={handleManualBan} className="space-y-4">
                <div>
                  <label className="text-[11px] font-black text-zinc-400 uppercase tracking-wider mb-1 block">
                    Target IP or CIDR Subnet
                  </label>
                  <input
                    type="text"
                    required
                    placeholder="e.g. 185.220.101.44 or 45.154.255.0/24"
                    value={manualIp}
                    onChange={(e) => setManualIp(e.target.value)}
                    className="w-full bg-zinc-950 border border-zinc-800 text-white rounded-xl px-4 py-3 text-xs font-mono focus:outline-none focus:ring-1 focus:ring-red-500"
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-[11px] font-black text-zinc-400 uppercase tracking-wider mb-1 block">
                      Quarantine Duration
                    </label>
                    <select
                      value={manualDuration}
                      onChange={(e) => setManualDuration(e.target.value)}
                      className="w-full bg-zinc-950 border border-zinc-800 text-zinc-200 rounded-xl px-3 py-2.5 text-xs font-bold focus:outline-none focus:ring-1 focus:ring-red-500"
                    >
                      <option value="1h">1 Hour (Cool-off)</option>
                      <option value="24h">24 Hours (Standard)</option>
                      <option value="7d">7 Days (Severe)</option>
                      <option value="permanent">Permanent Blackhole</option>
                    </select>
                  </div>

                  <div>
                    <label className="text-[11px] font-black text-zinc-400 uppercase tracking-wider mb-1 block">
                      Incident Severity
                    </label>
                    <div className="text-xs font-bold text-red-400 bg-red-950/40 border border-red-500/30 rounded-xl px-3 py-2.5">
                      CRITICAL BAN
                    </div>
                  </div>
                </div>

                <div>
                  <label className="text-[11px] font-black text-zinc-400 uppercase tracking-wider mb-1 block">
                    Reason / Incident Log Entry
                  </label>
                  <input
                    type="text"
                    value={manualReason}
                    onChange={(e) => setManualReason(e.target.value)}
                    className="w-full bg-zinc-950 border border-zinc-800 text-white rounded-xl px-4 py-2.5 text-xs font-medium focus:outline-none focus:ring-1 focus:ring-red-500"
                  />
                </div>

                <button
                  type="submit"
                  className="w-full py-3.5 rounded-xl bg-red-600 hover:bg-red-500 text-white font-black text-xs uppercase tracking-wider transition-all shadow-lg shadow-red-600/20 flex items-center justify-center gap-2"
                >
                  <Ban size={15} />
                  <span>Execute Instant IP Quarantine</span>
                </button>
              </form>
            </div>

            <div className="mt-6 pt-4 border-t border-zinc-800 text-[11px] text-zinc-500 flex items-center gap-2">
              <ShieldCheck size={14} className="text-emerald-400 shrink-0" />
              <span>Broadcasts real-time drop rules to Cloudflare WAF and Linux iptables.</span>
            </div>
          </div>

          {/* ACTIVE QUARANTINED IPS TABLE (7 Cols) */}
          <div className="lg:col-span-7 bg-zinc-900/80 border border-zinc-800/80 rounded-3xl p-6 lg:p-7 shadow-2xl flex flex-col justify-between">
            <div>
              <div className="flex items-center justify-between gap-3 mb-2">
                <div className="flex items-center gap-2">
                  <div className="p-1.5 bg-amber-600/10 text-amber-400 rounded-lg border border-amber-500/20">
                    <UserX size={16} />
                  </div>
                  <h2 className="text-base font-black text-white tracking-tight">Active Quarantined Entities & IP Blacklist</h2>
                </div>
                <span className="text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-md bg-zinc-800 text-zinc-300">
                  {blockedIps.length} Active
                </span>
              </div>
              <p className="text-xs text-zinc-400 leading-relaxed mb-6">
                Currently isolated addresses blocked from making RPC requests, creating listings, or initiating logins.
              </p>

              <div className="space-y-3 font-mono">
                {blockedIps.map((ip) => (
                  <div
                    key={ip}
                    className="flex items-center justify-between p-3.5 bg-zinc-950/80 border border-zinc-800 rounded-2xl hover:border-zinc-700 transition-all"
                  >
                    <div className="flex items-center gap-3">
                      <div className="p-2 rounded-xl bg-red-500/10 text-red-400 border border-red-500/20">
                        <Ban size={15} />
                      </div>
                      <div>
                        <div className="text-xs font-bold text-white">{ip}</div>
                        <div className="text-[10px] text-zinc-500 font-sans">Enforced by SOC Policy • Drop Packet</div>
                      </div>
                    </div>

                    <button
                      onClick={() => handleUnblock(ip)}
                      className="px-3 py-1.5 rounded-xl bg-zinc-800 hover:bg-zinc-700 text-zinc-300 text-xs font-bold font-sans transition-all border border-zinc-700"
                    >
                      Unblock / Whitelist
                    </button>
                  </div>
                ))}
              </div>
            </div>

            <div className="mt-6 pt-4 border-t border-zinc-800 text-[11px] text-zinc-500 flex items-center justify-between">
              <span>Automatic blacklist sync interval: 60s</span>
              <span className="text-emerald-400 font-bold">Quarantine Matrix: Synced</span>
            </div>
          </div>

        </section>

      </main>

      {/* ── EMERGENCY KILLSWITCH MODAL ── */}
      {killswitchModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-md animate-in fade-in">
          <div className="w-full max-w-lg bg-zinc-900 border border-red-500/50 rounded-3xl p-6 lg:p-8 shadow-2xl relative text-left">
            <div className="w-14 h-14 rounded-2xl bg-red-500/20 text-red-400 border border-red-500/30 flex items-center justify-center mb-4">
              <Power size={28} />
            </div>
            
            <h3 className="text-xl font-black text-white tracking-tight mb-2">
              {killswitchActive ? 'Disengage Platform Killswitch?' : 'ENGAGE EMERGENCY PLATFORM KILLSWITCH?'}
            </h3>
            <p className="text-xs text-zinc-400 leading-relaxed mb-6">
              {killswitchActive
                ? 'Deactivating killswitch will restore public write APIs, resume escrow payouts, and reopen administrative session modifications.'
                : 'CRITICAL SECURITY ACTION: Engaging the killswitch freezes high-value escrow payouts, sets all listing modifications to read-only, locks administrative privilege changes, and isolates the platform from external zero-day mutations.'}
            </p>

            <div className="flex items-center gap-3">
              <button
                type="button"
                onClick={() => setKillswitchModalOpen(false)}
                className="flex-1 py-3 rounded-xl bg-zinc-800 hover:bg-zinc-700 text-white font-bold text-xs uppercase tracking-wider transition-all"
              >
                Cancel
              </button>
              <button
                type="button"
                onClick={() => {
                  setKillswitchActive(!killswitchActive);
                  setKillswitchModalOpen(false);
                }}
                className="flex-1 py-3 rounded-xl bg-red-600 hover:bg-red-500 text-white font-black text-xs uppercase tracking-wider transition-all shadow-lg shadow-red-600/30"
              >
                {killswitchActive ? 'Confirm Deactivation' : 'Confirm Lockdown'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* ── FOOTER ── */}
      <footer className="border-t border-zinc-900/80 py-6 px-4 lg:px-8 mt-12 bg-zinc-950">
        <div className="max-w-[1720px] mx-auto flex flex-col md:flex-row items-center justify-between gap-4 text-xs text-zinc-500">
          <div className="flex items-center gap-3">
            <Lock size={14} className="text-emerald-400" />
            <span>HUNTIQ Security Operations Center (SOC) • Zero-Trust Perimeter Defense</span>
          </div>
          <div className="flex items-center gap-4">
            <Link href="/security-console/login" className="hover:text-white transition-colors">
              SOC Authenticator
            </Link>
            <Link href="/admin/logs" className="hover:text-white transition-colors">
              Security Audit Trails
            </Link>
          </div>
        </div>
      </footer>

    </div>
  );
}

