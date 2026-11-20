'use client';

import React, { useEffect, useState } from 'react';
import {
  Sparkles,
  Shield,
  AlertTriangle,
  Activity,
  DollarSign,
  Cpu,
  Clock,
  Layers,
  CheckCircle2,
  XCircle,
  RefreshCw,
} from 'lucide-react';
import { api, AIHealthStatus, AIMetricsReport } from '@/lib/api';

export function AIControlCenter() {
  const [health, setHealth] = useState<AIHealthStatus | null>(null);
  const [metrics, setMetrics] = useState<AIMetricsReport | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const fetchStats = async () => {
    try {
      setRefreshing(true);
      const [h, m] = await Promise.all([
        api.getAIHealth().catch(() => null),
        api.getAIAdminMetrics().catch(() => null),
      ]);
      if (h) setHealth(h);
      if (m) setMetrics(m);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, []);

  return (
    <div className="space-y-6">
      {/* ── Status Header Card ── */}
      <div className="p-6 rounded-3xl bg-gradient-to-r from-blue-900 via-indigo-900 to-zinc-900 text-white shadow-xl border border-blue-500/30 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-2xl bg-blue-500/20 backdrop-blur-md border border-blue-400/30 flex items-center justify-center">
            <Sparkles size={24} className="text-yellow-300 animate-pulse" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h3 className="text-xl font-black">HUNTIQ AI Control Center</h3>
              <span className="text-[10px] font-bold px-2 py-0.5 rounded-full bg-blue-500/30 text-blue-200 border border-blue-400/40 uppercase tracking-widest">
                Official Assistant
              </span>
            </div>
            <p className="text-xs text-blue-200 mt-0.5">
              Provider: {health?.provider || 'Google Gemini'} · Model: {health?.model || 'gemini-1.5-flash'}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2 px-3 py-1.5 rounded-2xl bg-white/10 border border-white/20 text-xs font-bold">
            {health?.status === 'operational' ? (
              <>
                <CheckCircle2 size={14} className="text-emerald-400" />
                <span className="text-emerald-300">Operational</span>
              </>
            ) : health?.status === 'disabled' ? (
              <>
                <XCircle size={14} className="text-zinc-400" />
                <span className="text-zinc-300">Disabled</span>
              </>
            ) : (
              <>
                <AlertTriangle size={14} className="text-amber-400" />
                <span className="text-amber-300">Degraded / Key Missing</span>
              </>
            )}
          </div>

          <button
            onClick={fetchStats}
            disabled={refreshing}
            className="p-2 rounded-xl bg-white/10 hover:bg-white/20 text-white transition-all disabled:opacity-50"
            title="Refresh metrics"
          >
            <RefreshCw size={15} className={refreshing ? 'animate-spin' : ''} />
          </button>
        </div>
      </div>

      {/* ── Metric Cards Grid ── */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Total Queries */}
        <div className="p-5 rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-sm">
          <div className="flex items-center justify-between text-zinc-400 mb-2">
            <span className="text-xs font-bold uppercase tracking-wider">Total Queries</span>
            <Activity size={18} className="text-blue-500" />
          </div>
          <p className="text-2xl font-black text-zinc-900 dark:text-white">
            {metrics?.total_queries?.toLocaleString() || 0}
          </p>
          <span className="text-[10px] text-zinc-400 mt-1 block">Past 30 days</span>
        </div>

        {/* Security / Adversarial Blocked */}
        <div className="p-5 rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-sm">
          <div className="flex items-center justify-between text-zinc-400 mb-2">
            <span className="text-xs font-bold uppercase tracking-wider">Attacks Blocked</span>
            <Shield size={18} className="text-emerald-500" />
          </div>
          <p className="text-2xl font-black text-emerald-600 dark:text-emerald-400">
            {metrics?.adversarial_blocks?.toLocaleString() || 0}
          </p>
          <span className="text-[10px] text-zinc-400 mt-1 block">Prompt injection & jailbreaks</span>
        </div>

        {/* Out of Domain Intercepted */}
        <div className="p-5 rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-sm">
          <div className="flex items-center justify-between text-zinc-400 mb-2">
            <span className="text-xs font-bold uppercase tracking-wider">Out of Domain</span>
            <AlertTriangle size={18} className="text-amber-500" />
          </div>
          <p className="text-2xl font-black text-amber-600 dark:text-amber-400">
            {metrics?.out_of_domain_queries?.toLocaleString() || 0}
          </p>
          <span className="text-[10px] text-zinc-400 mt-1 block">Politely redirected to HUNTIQ</span>
        </div>

        {/* Estimated API Cost */}
        <div className="p-5 rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-sm">
          <div className="flex items-center justify-between text-zinc-400 mb-2">
            <span className="text-xs font-bold uppercase tracking-wider">Estimated Cost</span>
            <DollarSign size={18} className="text-purple-500" />
          </div>
          <p className="text-2xl font-black text-purple-600 dark:text-purple-400">
            ${(metrics?.estimated_cost_usd || 0).toFixed(4)}
          </p>
          <span className="text-[10px] text-zinc-400 mt-1 block">
            {(metrics?.total_tokens || 0).toLocaleString()} tokens consumed
          </span>
        </div>
      </div>

      {/* ── Operational Architecture Guidelines ── */}
      <div className="p-6 rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-900/50 space-y-3">
        <h4 className="text-sm font-black text-zinc-900 dark:text-white flex items-center gap-2">
          <Shield size={16} className="text-blue-600" />
          <span>HUNTIQ AI Domain Policy & Non-Negotiable Safety Invariants</span>
        </h4>
        <ul className="text-xs text-zinc-600 dark:text-zinc-300 space-y-2 list-disc pl-5 leading-relaxed">
          <li>
            <strong>Strict Domain Restriction:</strong> Assistant answers solely for physical goods, property rentals, transport, services, gigs, tenders, events, and platform guidance. Academic math, quantum physics, coding tutorials, and general trivia are deterministically blocked.
          </li>
          <li>
            <strong>Server-Side Secret Isolation:</strong> <code>GEMINI_API_KEY</code> is loaded strictly in the Go backend. Keys are never bundled into client Next.js files or returned in API responses.
          </li>
          <li>
            <strong>4-Language Enforcement:</strong> English (en), Amharic (am), Tigrinya (ti), and Afaan Oromo (or). Unsupported languages receive an automatic polite invitation in the 4 supported languages.
          </li>
          <li>
            <strong>Database Authenticity:</strong> All listing recommendations and photographs are verified against canonical PostgreSQL records before display.
          </li>
        </ul>
      </div>
    </div>
  );
}
