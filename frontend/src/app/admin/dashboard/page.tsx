'use client';

import React, { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import {
  LayoutDashboard, TrendingUp, ShoppingBag, FileText,
  DollarSign, Users, Clock, AlertTriangle, Package,
  CheckCircle2, RefreshCw, BarChart3, Activity,
} from 'lucide-react';
import Link from 'next/link';

interface KPI {
  label: string;
  value: string | number;
  sub?: string;
  icon: React.ReactNode;
  color: string;
  bg: string;
  href?: string;
}

export default function AdminDashboard() {
  const { user, hasAdminPermission, isSuperAdmin } = useAuth();
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);

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

  // Build KPIs from real backend data
  const kpis: KPI[] = [
    {
      label: 'Total Listings',
      value: s.total_products?.toLocaleString() ?? '—',
      sub: `${s.approved_products ?? 0} active`,
      icon: <Package size={20} />,
      color: 'text-orange-400',
      bg: 'bg-orange-600/10 border-orange-500/20',
      href: '/admin/products',
    },
    {
      label: 'Pending Review',
      value: s.pending_products?.toLocaleString() ?? '—',
      sub: 'awaiting moderation',
      icon: <Clock size={20} />,
      color: 'text-amber-400',
      bg: 'bg-amber-600/10 border-amber-500/20',
      href: '/admin/products',
    },
    {
      label: 'Total Sellers',
      value: s.total_sellers?.toLocaleString() ?? '—',
      sub: `${s.active_sellers ?? 0} active`,
      icon: <Users size={20} />,
      color: 'text-blue-400',
      bg: 'bg-blue-600/10 border-blue-500/20',
      href: isSuperAdmin ? '/admin/users' : undefined,
    },
    {
      label: 'Total Buyers',
      value: s.total_buyers?.toLocaleString() ?? '—',
      sub: `${s.premium_buyers ?? 0} gold/premium`,
      icon: <ShoppingBag size={20} />,
      color: 'text-purple-400',
      bg: 'bg-purple-600/10 border-purple-500/20',
    },
  ];

  // Super-admin-only financial KPIs
  const financialKpis: KPI[] = isSuperAdmin ? [
    {
      label: 'Total Revenue',
      value: s.total_revenue != null ? `ETB ${Number(s.total_revenue).toLocaleString()}` : '—',
      sub: 'all-time listing fees',
      icon: <DollarSign size={20} />,
      color: 'text-emerald-400',
      bg: 'bg-emerald-600/10 border-emerald-500/20',
      href: '/admin/settings',
    },
    {
      label: 'Fee Revenue',
      value: s.listing_fee_revenue != null ? `ETB ${Number(s.listing_fee_revenue).toLocaleString()}` : '—',
      sub: 'confirmed payments',
      icon: <BarChart3 size={20} />,
      color: 'text-teal-400',
      bg: 'bg-teal-600/10 border-teal-500/20',
    },
  ] : [];

  return (
    <div className="space-y-6">

      {/* Header */}
      <div className="flex flex-wrap items-center justify-between gap-4 bg-zinc-900 border border-zinc-800 p-5 rounded-2xl">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-gradient-to-tr from-orange-600 to-amber-600 rounded-2xl flex items-center justify-center shadow-lg shadow-orange-600/20">
            <LayoutDashboard size={24} className="text-white" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Platform Overview</h1>
            <p className="text-xs text-zinc-400 mt-0.5">
              {lastUpdated ? `Updated ${lastUpdated.toLocaleTimeString()}` : 'Loading…'}
            </p>
          </div>
        </div>
        <button
          onClick={() => load(true)}
          disabled={refreshing}
          className="flex items-center gap-2 px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-sm font-bold transition-colors disabled:opacity-50"
        >
          <RefreshCw size={14} className={refreshing ? 'animate-spin' : ''} />
          Refresh
        </button>
      </div>

      {/* KPI Grid */}
      {loading ? (
        <div className="flex items-center justify-center h-48">
          <Spinner />
        </div>
      ) : (
        <>
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
            {kpis.map(kpi => (
              <KPICard key={kpi.label} kpi={kpi} />
            ))}
          </div>

          {financialKpis.length > 0 && (
            <div className="grid grid-cols-2 gap-4">
              {financialKpis.map(kpi => (
                <KPICard key={kpi.label} kpi={kpi} />
              ))}
            </div>
          )}
        </>
      )}

      {/* Quick Actions */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {hasAdminPermission('pending') && (
          <QuickAction
            href="/admin/products"
            icon={<Clock size={18} className="text-amber-400" />}
            label="Review Pending"
            desc="Listings awaiting moderation"
            count={s.pending_products}
            accent="border-amber-500/20 hover:border-amber-500/40"
          />
        )}
        {hasAdminPermission('sellers') && (
          <QuickAction
            href="/admin/users"
            icon={<Users size={18} className="text-blue-400" />}
            label="Manage Users"
            desc="Sellers, buyers, and admins"
            accent="border-blue-500/20 hover:border-blue-500/40"
          />
        )}
        {hasAdminPermission('analytics') && (
          <QuickAction
            href="/admin/settings"
            icon={<BarChart3 size={18} className="text-purple-400" />}
            label="Platform Settings"
            desc="Fees, limits, and toggles"
            accent="border-purple-500/20 hover:border-purple-500/40"
          />
        )}
        {hasAdminPermission('seller_verification') && (
          <QuickAction
            href="/admin/verification"
            icon={<CheckCircle2 size={18} className="text-emerald-400" />}
            label="Verify Sellers"
            desc="Identity and business documents"
            accent="border-emerald-500/20 hover:border-emerald-500/40"
          />
        )}
        {hasAdminPermission('trust_safety') && (
          <QuickAction
            href="/admin/trust-safety"
            icon={<AlertTriangle size={18} className="text-orange-400" />}
            label="Trust & Safety"
            desc="Abuse reports and fraud"
            accent="border-orange-500/20 hover:border-orange-500/40"
          />
        )}
        {hasAdminPermission('operations') && (
          <QuickAction
            href="/admin/operations"
            icon={<Activity size={18} className="text-teal-400" />}
            label="Operations"
            desc="Background jobs and health"
            accent="border-teal-500/20 hover:border-teal-500/40"
          />
        )}
      </div>

      {/* Recent monthly trend — only for analytics permission */}
      {hasAdminPermission('analytics') && <MonthlyTrend token={user?.token ?? ''} />}
    </div>
  );
}

function KPICard({ kpi }: { kpi: KPI }) {
  const inner = (
    <div className={`p-5 bg-zinc-900 border rounded-2xl space-y-3 transition-all hover:border-zinc-600 ${kpi.href ? 'cursor-pointer' : ''}`}>
      <div className="flex items-center justify-between">
        <span className="text-xs font-bold uppercase tracking-wider text-zinc-400">{kpi.label}</span>
        <div className={`p-2 rounded-xl border ${kpi.bg}`}>
          <span className={kpi.color}>{kpi.icon}</span>
        </div>
      </div>
      <div className="text-2xl font-black text-white font-mono">{kpi.value}</div>
      {kpi.sub && (
        <div className="flex items-center gap-1.5 text-[11px] font-semibold text-zinc-500">
          <TrendingUp size={11} />
          {kpi.sub}
        </div>
      )}
    </div>
  );
  if (kpi.href) return <Link href={kpi.href}>{inner}</Link>;
  return inner;
}

function QuickAction({ href, icon, label, desc, count, accent }: {
  href: string; icon: React.ReactNode; label: string; desc: string; count?: number; accent: string;
}) {
  return (
    <Link href={href}
      className={`flex items-center gap-4 p-4 bg-zinc-900 border ${accent} rounded-2xl transition-all hover:bg-zinc-800`}>
      <div className="p-2.5 bg-zinc-800 rounded-xl shrink-0">{icon}</div>
      <div className="min-w-0">
        <p className="font-black text-sm text-white">{label}</p>
        <p className="text-xs text-zinc-500 truncate">{desc}</p>
      </div>
      {count != null && count > 0 && (
        <span className="ml-auto shrink-0 px-2 py-0.5 bg-orange-600/20 text-orange-400 border border-orange-500/20 rounded-full text-xs font-bold">
          {count}
        </span>
      )}
    </Link>
  );
}

function MonthlyTrend({ token }: { token: string }) {
  const [data, setData] = useState<any[]>([]);

  useEffect(() => {
    if (!token) return;
    api.adminAnalyticsMonthly(3, token).then(d => setData(d ?? [])).catch(() => {});
  }, [token]);

  if (data.length === 0) return null;

  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-5">
      <div className="flex items-center gap-2 mb-4">
        <BarChart3 size={16} className="text-orange-500" />
        <h2 className="font-black text-sm text-white uppercase tracking-wider">Monthly Activity</h2>
      </div>
      <div className="grid grid-cols-3 gap-3">
        {data.map((row: any, i: number) => (
          <div key={i} className="bg-zinc-800 rounded-xl p-3 space-y-1">
            <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-wider">{row.month}</p>
            <p className="text-lg font-black text-white font-mono">{row.new_products ?? 0}</p>
            <p className="text-[10px] text-zinc-500">new listings</p>
          </div>
        ))}
      </div>
    </div>
  );
}
