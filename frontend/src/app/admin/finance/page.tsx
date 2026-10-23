'use client';

import { useEffect, useState, useMemo } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import {
  TrendingUp, DollarSign, Activity, ShieldAlert, FileText, Download,
  RefreshCw, CheckCircle2, AlertTriangle, Scale, Server, Cpu, Database,
  Layers, Users, HelpCircle, ArrowUpRight, ArrowDownRight, RotateCcw,
  Clock, Shield, BarChart3, Calculator
} from 'lucide-react';

function fmtETB(n: any) {
  if (n === null || n === undefined || isNaN(Number(n))) return 'ETB 0.00';
  const val = Number(n);
  return 'ETB ' + val.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function fmtNum(n: any) {
  if (n === null || n === undefined || isNaN(Number(n))) return '0';
  return Number(n).toLocaleString('en-US');
}

export default function AdminFinancePage() {
  const { user, isPlatformOwner, isPlatformAdmin, hasAdminPermission } = useAuth();

  const isAuthorized = isPlatformOwner || isPlatformAdmin || hasAdminPermission('profit') || hasAdminPermission('finance');
  const canManageFinance = isPlatformOwner || hasAdminPermission('finance:manage') || hasAdminPermission('finance_manage');

  const [scenario, setScenario] = useState<'base' | 'conservative' | 'high_growth'>('base');
  const [horizon, setHorizon] = useState<number>(60);
  const [growthMode, setGrowthMode] = useState<'smooth_quarterly' | 'monthly'>('smooth_quarterly');
  const [activeTab, setActiveTab] = useState<'overview' | 'growth' | 'revenue' | 'opex' | 'tax' | 'cashflow' | 'sensitivity' | 'assumptions' | 'actuals'>('overview');

  const [report, setReport] = useState<any>(null);
  const [scenariosData, setScenariosData] = useState<any>(null);
  const [capacity20M, setCapacity20M] = useState<any>(null);
  const [taxRules, setTaxRules] = useState<any[]>([]);
  const [actuals, setActuals] = useState<any[]>([]);
  const [auditLogs, setAuditLogs] = useState<any[]>([]);

  const [loading, setLoading] = useState<boolean>(true);
  const [refreshing, setRefreshing] = useState<boolean>(false);
  const [feedbackMsg, setFeedbackMsg] = useState<string | null>(null);

  const loadData = async (silent = false) => {
    if (!user?.token || !isAuthorized) return;
    if (!silent) setLoading(true);
    else setRefreshing(true);

    try {
      const [rep, scs, cap, rules, acts, logs] = await Promise.allSettled([
        api.adminGetFinanceForecast({ scenario, horizon }, user.token),
        api.adminGetFinanceScenarios({ horizon }, user.token),
        api.adminGetFinanceCapacity20M(user.token),
        api.adminGetFinanceTaxRules(user.token),
        api.adminGetFinanceActuals(user.token),
        api.adminGetFinanceAuditLogs(user.token),
      ]);

      if (rep.status === 'fulfilled') setReport(rep.value);
      if (scs.status === 'fulfilled') setScenariosData(scs.value);
      if (cap.status === 'fulfilled') setCapacity20M(cap.value);
      if (rules.status === 'fulfilled') setTaxRules(rules.value ?? []);
      if (acts.status === 'fulfilled') setActuals(acts.value ?? []);
      if (logs.status === 'fulfilled') setAuditLogs(logs.value ?? []);
    } catch (err) {
      console.error('Error loading financial forecasting data:', err);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => {
    loadData();
  }, [user, scenario, horizon]); // eslint-disable-line

  const handleResetAssumptions = async () => {
    if (!user?.token) return;
    if (!canManageFinance) {
      alert('Unauthorized: Chief Financial Officer or Platform Owner authorization is required to modify or reset financial assumptions.');
      return;
    }
    if (!confirm('Revert financial assumptions for ' + scenario + ' scenario to system defaults?')) return;
    try {
      await api.adminResetFinanceAssumptions({ scenario_id: scenario }, user.token);
      setFeedbackMsg('Assumptions reverted to verified system defaults.');
      setTimeout(() => setFeedbackMsg(null), 4000);
      loadData(true);
    } catch (err: any) {
      alert('Failed to reset: ' + (err.message || 'Error'));
    }
  };

  if (!isAuthorized) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[500px] gap-4 text-zinc-500">
        <div className="p-4 bg-red-950/20 text-red-400 rounded-2xl border border-red-500/20">
          <ShieldAlert size={36} />
        </div>
        <h2 className="text-lg font-black text-white">Financial Authorization Required</h2>
        <p className="text-xs text-zinc-400 text-center max-w-md">
          Access to HUNTIQ Financial Planning & Analysis (FP&A) is strictly restricted to Platform Owners and authorized Finance personnel.
        </p>
      </div>
    );
  }

  if (loading && !report) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[500px] gap-3">
        <Spinner />
        <p className="text-xs font-mono text-zinc-500">Calculating financial models and capacity forecasts...</p>
      </div>
    );
  }

  const latestMonth = report?.MonthlySeries?.[report?.MonthlySeries?.length - 1] || {};
  const month12 = report?.MonthlySeries?.[11] || {};
  const milestones = report?.Milestones || {};

  return (
    <div className="space-y-6 pb-16">
      {/* Header Bar */}
      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4 bg-zinc-950 border border-zinc-900 rounded-3xl p-6">
        <div>
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-emerald-600/10 border border-emerald-500/20 rounded-2xl flex items-center justify-center">
              <TrendingUp size={24} className="text-emerald-400" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h1 className="text-2xl font-black text-white tracking-tight">Financial Planning & Analysis (FP&A)</h1>
                <span className="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                  {report?.Currency || 'ETB'} Default
                </span>
                <span className="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider bg-indigo-500/10 text-indigo-400 border border-indigo-500/20">
                  {horizon} Months Model
                </span>
              </div>
              <p className="text-xs text-zinc-400 mt-0.5">
                Multi-channel marketplace forecasting, infrastructure capacity planning, and Ethiopian tax estimation.
              </p>
            </div>
          </div>
        </div>

        {/* Global Controls */}
        <div className="flex flex-wrap items-center gap-3">
          {/* Scenario Selector */}
          <div className="flex bg-zinc-900 border border-zinc-800 rounded-2xl p-1 text-xs font-bold">
            {(['conservative', 'base', 'high_growth'] as const).map((sc) => (
              <button
                key={sc}
                onClick={() => setScenario(sc)}
                className={`px-3 py-1.5 rounded-xl transition-all capitalize ${
                  scenario === sc
                    ? 'bg-orange-600 text-white shadow-lg shadow-orange-600/30'
                    : 'text-zinc-400 hover:text-white'
                }`}
              >
                {sc.replace('_', ' ')}
              </button>
            ))}
          </div>

          {/* Horizon Selector */}
          <select
            value={horizon}
            onChange={(e) => setHorizon(Number(e.target.value))}
            className="bg-zinc-900 border border-zinc-800 text-zinc-300 text-xs font-bold rounded-2xl px-3 py-2 outline-none focus:border-orange-500"
          >
            <option value={12}>12 Months (1 Year)</option>
            <option value={24}>24 Months (2 Years)</option>
            <option value={36}>36 Months (3 Years)</option>
            <option value={60}>60 Months (5 Years)</option>
          </select>

          {/* CSV Export Button */}
          <a
            href={api.adminExportFinanceURL(scenario, horizon, 'csv')}
            download
            className="flex items-center gap-2 px-3.5 py-2 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 text-zinc-300 rounded-2xl text-xs font-bold transition-all"
          >
            <Download size={14} className="text-emerald-400" />
            Export CSV
          </a>

          {/* Refresh Button */}
          <button
            onClick={() => loadData(true)}
            disabled={refreshing}
            className="p-2.5 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 text-zinc-300 rounded-2xl transition-all"
            title="Refresh Financial Forecast"
          >
            <RefreshCw size={14} className={refreshing ? 'animate-spin text-orange-400' : ''} />
          </button>
        </div>
      </div>

      {feedbackMsg && (
        <div className="flex items-center gap-3 p-4 bg-emerald-950/30 border border-emerald-500/30 text-emerald-300 rounded-2xl text-xs font-bold">
          <CheckCircle2 size={16} />
          {feedbackMsg}
        </div>
      )}

      {/* Statutory Disclaimers Notice */}
      <div className="flex items-start gap-3 p-4 bg-amber-950/20 border border-amber-500/30 rounded-2xl text-amber-300 text-xs">
        <AlertTriangle size={18} className="shrink-0 text-amber-400 mt-0.5" />
        <div className="space-y-1">
          <p className="font-bold">
            MANAGEMENT PLANNING AND CAPACITY FORECAST ONLY — NOT CERTIFIED ACCOUNTING RECORDS
          </p>
          <p className="text-amber-200/80 leading-relaxed text-[11px]">
            {report?.DisclaimerNotice || 'Projections represent planning models under simulated adoption assumptions. Market funds held for sellers (GMV) are segregated from HUNTIQ platform revenue. Ethiopian Corporate Income Tax (30%) and VAT (15%) estimates require validation by an authorized Ethiopian tax consultant under Proclamations 979/2016 and 285/2002 & 1341/2024.'}
          </p>
        </div>
      </div>

      {/* KPI Overview Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-zinc-950 border border-zinc-900 rounded-2xl p-5">
          <div className="flex items-center justify-between text-zinc-400 mb-2">
            <span className="text-xs font-extrabold uppercase tracking-wider">Gross Merchandise Volume (GMV)</span>
            <Scale size={16} className="text-indigo-400" />
          </div>
          <div className="text-xl font-black font-mono text-white">{fmtETB(latestMonth.gross_merchandise_volume)}</div>
          <p className="text-[10px] text-zinc-500 mt-1">Escrow funds belonging to sellers (Month {latestMonth.month_index})</p>
        </div>

        <div className="bg-zinc-950 border border-zinc-900 rounded-2xl p-5">
          <div className="flex items-center justify-between text-zinc-400 mb-2">
            <span className="text-xs font-extrabold uppercase tracking-wider">Net Recognized Revenue</span>
            <DollarSign size={16} className="text-emerald-400" />
          </div>
          <div className="text-xl font-black font-mono text-emerald-400">{fmtETB(latestMonth.net_recognized_revenue)}</div>
          <p className="text-[10px] text-zinc-500 mt-1">Commission, listings & tiers (Month {latestMonth.month_index})</p>
        </div>

        <div className="bg-zinc-950 border border-zinc-900 rounded-2xl p-5">
          <div className="flex items-center justify-between text-zinc-400 mb-2">
            <span className="text-xs font-extrabold uppercase tracking-wider">Total Operating Expenses</span>
            <Layers size={16} className="text-red-400" />
          </div>
          <div className="text-xl font-black font-mono text-white">{fmtETB(latestMonth.total_operating_expense)}</div>
          <p className="text-[10px] text-zinc-500 mt-1">Cloud infra + staffing + OpEx</p>
        </div>

        <div className="bg-zinc-950 border border-zinc-900 rounded-2xl p-5">
          <div className="flex items-center justify-between text-zinc-400 mb-2">
            <span className="text-xs font-extrabold uppercase tracking-wider">Net Profit After Tax</span>
            <Activity size={16} className={latestMonth.net_profit_after_tax >= 0 ? 'text-emerald-400' : 'text-red-400'} />
          </div>
          <div className={`text-xl font-black font-mono ${latestMonth.net_profit_after_tax >= 0 ? 'text-emerald-400' : 'text-red-400'}`}>
            {fmtETB(latestMonth.net_profit_after_tax)}
          </div>
          <p className="text-[10px] text-zinc-500 mt-1">After 30% CIT & VAT credits</p>
        </div>

        <div className="bg-zinc-950 border border-zinc-900 rounded-2xl p-5">
          <div className="flex items-center justify-between text-zinc-400 mb-2">
            <span className="text-xs font-extrabold uppercase tracking-wider">Closing Cash Balance</span>
            <DollarSign size={16} className="text-blue-400" />
          </div>
          <div className="text-xl font-black font-mono text-blue-400">{fmtETB(latestMonth.closing_cash_balance)}</div>
          <p className="text-[10px] text-zinc-500 mt-1">Projected liquidity balance</p>
        </div>

        <div className="bg-zinc-950 border border-zinc-900 rounded-2xl p-5">
          <div className="flex items-center justify-between text-zinc-400 mb-2">
            <span className="text-xs font-extrabold uppercase tracking-wider">Monthly Cash Burn</span>
            <TrendingUp size={16} className="text-orange-400" />
          </div>
          <div className="text-xl font-black font-mono text-orange-400">{fmtETB(latestMonth.monthly_burn_rate)}</div>
          <p className="text-[10px] text-zinc-500 mt-1">Net monthly cash deficit</p>
        </div>

        <div className="bg-zinc-950 border border-zinc-900 rounded-2xl p-5">
          <div className="flex items-center justify-between text-zinc-400 mb-2">
            <span className="text-xs font-extrabold uppercase tracking-wider">Break-Even Milestone</span>
            <CheckCircle2 size={16} className="text-emerald-400" />
          </div>
          <div className="text-xl font-black font-mono text-white">
            {milestones.break_even_month_index > 0 ? `Month ${milestones.break_even_month_index}` : 'Deferred'}
          </div>
          <p className="text-[10px] text-zinc-500 mt-1">
            {milestones.break_even_period || 'Not reached within horizon'}
          </p>
        </div>

        <div className="bg-zinc-950 border border-zinc-900 rounded-2xl p-5">
          <div className="flex items-center justify-between text-zinc-400 mb-2">
            <span className="text-xs font-extrabold uppercase tracking-wider">Active Users (MAU)</span>
            <Users size={16} className="text-purple-400" />
          </div>
          <div className="text-xl font-black font-mono text-purple-400">{fmtNum(latestMonth.mau)}</div>
          <p className="text-[10px] text-zinc-500 mt-1">DAU: {fmtNum(latestMonth.dau)} ({((latestMonth.dau / (latestMonth.mau || 1)) * 100).toFixed(0)}%)</p>
        </div>
      </div>

      {/* Navigation Tabs */}
      <div className="flex items-center gap-2 border-b border-zinc-900 overflow-x-auto pb-2 text-xs font-bold">
        {[
          { id: 'overview', label: 'Overview & P&L', icon: BarChart3 },
          { id: 'growth', label: 'Growth & Workload', icon: Activity },
          { id: 'revenue', label: 'Revenue Streams', icon: DollarSign },
          { id: 'opex', label: 'OpEx & Cloud Infra', icon: Server },
          { id: 'tax', label: 'Ethiopia Tax Estimator', icon: Calculator },
          { id: 'cashflow', label: 'Cash Flow & Runway', icon: TrendingUp },
          { id: 'sensitivity', label: 'Sensitivity & Stress Tests', icon: ShieldAlert },
          { id: 'assumptions', label: 'Assumptions & Audit', icon: Layers },
          { id: 'actuals', label: 'Historical Actuals', icon: FileText },
        ].map((t) => (
          <button
            key={t.id}
            onClick={() => setActiveTab(t.id as any)}
            className={`flex items-center gap-2 px-4 py-2.5 rounded-xl transition-all shrink-0 ${
              activeTab === t.id
                ? 'bg-orange-600/15 text-orange-400 border border-orange-500/20'
                : 'text-zinc-400 hover:text-white hover:bg-zinc-900'
            }`}
          >
            <t.icon size={15} />
            {t.label}
          </button>
        ))}
      </div>

      {/* TAB CONTENT 1: Overview & P&L */}
      {activeTab === 'overview' && (
        <div className="space-y-6">
          {/* Revenue vs OpEx SVG Chart */}
          <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6">
            <div className="flex items-center justify-between mb-6">
              <div>
                <h3 className="text-base font-black text-white">Monthly Revenue vs Operating Expenses Trajectory</h3>
                <p className="text-xs text-zinc-400">Net recognized platform revenue compared to total hosting, staffing, and operating expenses.</p>
              </div>
              <div className="flex items-center gap-4 text-xs font-bold">
                <span className="flex items-center gap-1.5 text-emerald-400">
                  <span className="w-3 h-3 bg-emerald-500 rounded-full inline-block"></span> Net Revenue
                </span>
                <span className="flex items-center gap-1.5 text-red-400">
                  <span className="w-3 h-3 bg-red-500 rounded-full inline-block"></span> Total OpEx
                </span>
                <span className="flex items-center gap-1.5 text-blue-400">
                  <span className="w-3 h-3 bg-blue-500 rounded-full inline-block"></span> Net Profit
                </span>
              </div>
            </div>

            {/* Custom SVG Trend Visualization */}
            <div className="h-64 w-full">
              <svg viewBox="0 0 1000 240" className="w-full h-full overflow-visible">
                <defs>
                  <linearGradient id="revGrad" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="#10b981" stopOpacity="0.3" />
                    <stop offset="100%" stopColor="#10b981" stopOpacity="0.0" />
                  </linearGradient>
                </defs>
                {/* Horizontal Grid lines */}
                {[0.25, 0.5, 0.75, 1.0].map((pct, i) => (
                  <line
                    key={i}
                    x1="0"
                    y1={240 - pct * 220}
                    x2="1000"
                    y2={240 - pct * 220}
                    stroke="#27272a"
                    strokeDasharray="4 4"
                    strokeWidth="1"
                  />
                ))}

                {/* Plot Series Path */}
                {(() => {
                  const series = report?.MonthlySeries || [];
                  if (series.length < 2) return null;
                  const maxVal = Math.max(...series.map((s: any) => Math.max(s.net_recognized_revenue, s.total_operating_expense))) * 1.15 || 100000;

                  const revPoints = series.map((s: any, i: number) => {
                    const x = (i / (series.length - 1)) * 1000;
                    const y = 240 - (s.net_recognized_revenue / maxVal) * 220;
                    return `${x},${y}`;
                  }).join(' ');

                  const opexPoints = series.map((s: any, i: number) => {
                    const x = (i / (series.length - 1)) * 1000;
                    const y = 240 - (s.total_operating_expense / maxVal) * 220;
                    return `${x},${y}`;
                  }).join(' ');

                  const profitPoints = series.map((s: any, i: number) => {
                    const x = (i / (series.length - 1)) * 1000;
                    const clampedProfit = Math.max(0, s.net_profit_after_tax);
                    const y = 240 - (clampedProfit / maxVal) * 220;
                    return `${x},${y}`;
                  }).join(' ');

                  return (
                    <>
                      <polygon points={`0,240 ${revPoints} 1000,240`} fill="url(#revGrad)" />
                      <polyline points={revPoints} fill="none" stroke="#10b981" strokeWidth="2.5" />
                      <polyline points={opexPoints} fill="none" stroke="#ef4444" strokeWidth="2" strokeDasharray="3 3" />
                      <polyline points={profitPoints} fill="none" stroke="#3b82f6" strokeWidth="2" />
                    </>
                  );
                })()}
              </svg>
            </div>
            <div className="flex justify-between text-[11px] font-mono text-zinc-500 mt-2">
              <span>Month 1 (Nov 2026)</span>
              <span>Month {Math.round(horizon / 2)}</span>
              <span>Month {horizon} End</span>
            </div>
          </div>

          {/* 5-Year Annual Summary Table */}
          <div className="bg-zinc-950 border border-zinc-900 rounded-3xl overflow-hidden">
            <div className="p-5 border-b border-zinc-900 flex items-center justify-between">
              <div>
                <h3 className="text-base font-black text-white">Annual Pro-Forma Income Statement Summary</h3>
                <p className="text-xs text-zinc-400">Multi-year overview of gross merchandise volume, revenue, EBITDA, tax accruals, and cash position.</p>
              </div>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead>
                  <tr className="border-b border-zinc-900 text-zinc-400 uppercase font-extrabold bg-zinc-900/40">
                    <th className="p-4">Financial Metric</th>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <th key={ann.year_number} className="p-4 text-right">
                        Year {ann.year_number} ({ann.calendar_year})
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody className="divide-y divide-zinc-900/60 font-mono text-zinc-300">
                  <tr className="hover:bg-zinc-900/30">
                    <td className="p-4 font-sans font-bold text-white">Active Users (Ending MAU)</td>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <td key={ann.year_number} className="p-4 text-right">{fmtNum(ann.ending_mau)}</td>
                    ))}
                  </tr>
                  <tr className="hover:bg-zinc-900/30">
                    <td className="p-4 font-sans font-bold text-indigo-400">Gross Merchandise Volume (GMV)</td>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <td key={ann.year_number} className="p-4 text-right font-bold text-indigo-400">{fmtETB(ann.total_gmv)}</td>
                    ))}
                  </tr>
                  <tr className="hover:bg-zinc-900/30 bg-emerald-950/10">
                    <td className="p-4 font-sans font-bold text-emerald-400">Net Recognized Platform Revenue</td>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <td key={ann.year_number} className="p-4 text-right font-bold text-emerald-400">{fmtETB(ann.net_revenue)}</td>
                    ))}
                  </tr>
                  <tr className="hover:bg-zinc-900/30">
                    <td className="p-4 font-sans text-zinc-400">Cloud & Infrastructure Hosting</td>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <td key={ann.year_number} className="p-4 text-right text-red-300">{fmtETB(ann.total_infra_cost)}</td>
                    ))}
                  </tr>
                  <tr className="hover:bg-zinc-900/30">
                    <td className="p-4 font-sans text-zinc-400">Staffing & Employer Pension (11%)</td>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <td key={ann.year_number} className="p-4 text-right text-red-300">{fmtETB(ann.total_staffing_cost)}</td>
                    ))}
                  </tr>
                  <tr className="hover:bg-zinc-900/30">
                    <td className="p-4 font-sans font-bold text-zinc-300">Total Operating Expenses (OpEx)</td>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <td key={ann.year_number} className="p-4 text-right font-bold text-red-400">{fmtETB(ann.total_opex)}</td>
                    ))}
                  </tr>
                  <tr className="hover:bg-zinc-900/30">
                    <td className="p-4 font-sans font-bold text-purple-400">EBITDA</td>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <td key={ann.year_number} className="p-4 text-right font-bold text-purple-400">{fmtETB(ann.ebitda)}</td>
                    ))}
                  </tr>
                  <tr className="hover:bg-zinc-900/30">
                    <td className="p-4 font-sans font-bold text-blue-400">Operating Profit (EBIT)</td>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <td key={ann.year_number} className="p-4 text-right font-bold text-blue-400">{fmtETB(ann.operating_profit)}</td>
                    ))}
                  </tr>
                  <tr className="hover:bg-zinc-900/30">
                    <td className="p-4 font-sans text-amber-300">Estimated Corporate Income Tax (30%)</td>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <td key={ann.year_number} className="p-4 text-right text-amber-300">{fmtETB(ann.estimated_cit)}</td>
                    ))}
                  </tr>
                  <tr className="hover:bg-zinc-900/30">
                    <td className="p-4 font-sans text-amber-300">Estimated Net VAT Payable (15%)</td>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <td key={ann.year_number} className="p-4 text-right text-amber-300">{fmtETB(ann.net_vat_payable)}</td>
                    ))}
                  </tr>
                  <tr className="hover:bg-zinc-900/30 bg-emerald-950/20 font-bold">
                    <td className="p-4 font-sans text-emerald-400">Net Profit After Estimated Tax</td>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <td key={ann.year_number} className={`p-4 text-right font-bold ${ann.net_profit_after_tax >= 0 ? 'text-emerald-400' : 'text-red-400'}`}>
                        {fmtETB(ann.net_profit_after_tax)}
                      </td>
                    ))}
                  </tr>
                  <tr className="hover:bg-zinc-900/30">
                    <td className="p-4 font-sans font-bold text-blue-400">Ending Cash Balance</td>
                    {report?.AnnualSummaries?.map((ann: any) => (
                      <td key={ann.year_number} className="p-4 text-right font-bold text-blue-400">{fmtETB(ann.ending_cash_balance)}</td>
                    ))}
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {/* TAB CONTENT 2: User Growth, Workload & 20M DAU Stress Case */}
      {activeTab === 'growth' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Growth Parameters */}
            <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 space-y-4">
              <h3 className="text-base font-black text-white flex items-center gap-2">
                <Users size={18} className="text-purple-400" />
                Growth Modeling Levers & Quarterly Reconciliation
              </h3>
              <p className="text-xs text-zinc-400">Quarterly growth targets reconciled into smooth monthly compounding to prevent artificial operational spikes.</p>
              <div className="space-y-3 text-xs">
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="text-zinc-400">Model Algorithm:</span>
                  <span className="font-bold text-white uppercase">Quarterly Smooth Transition</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="text-zinc-400">Target Quarterly Rate:</span>
                  <span className="font-mono text-emerald-400 font-bold">{scenario === 'high_growth' ? '28%' : scenario === 'conservative' ? '10%' : '20%'} / quarter</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="text-zinc-400">Smooth Monthly Equivalent:</span>
                  <span className="font-mono text-emerald-400 font-bold">{scenario === 'high_growth' ? '8.57%' : scenario === 'conservative' ? '3.23%' : '6.27%'} / month</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="text-zinc-400">Early Stage Duration:</span>
                  <span className="font-mono text-white">4 Quarters (12 Months)</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="text-zinc-400">Mature Growth Rate:</span>
                  <span className="font-mono text-white font-bold">{scenario === 'high_growth' ? '3.5%' : scenario === 'conservative' ? '1.5%' : '2.0%'} / month</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="text-zinc-400">Market Saturation Ceiling:</span>
                  <span className="font-mono text-white font-bold">20,000,000 Active Users</span>
                </div>
                <div className="flex justify-between py-2">
                  <span className="text-zinc-400">DAU / MAU Ratio:</span>
                  <span className="font-mono text-white font-bold">{scenario === 'high_growth' ? '30%' : scenario === 'conservative' ? '20%' : '25%'}</span>
                </div>
              </div>
              <div className="p-3 bg-indigo-950/20 border border-indigo-500/20 rounded-2xl text-[11px] text-indigo-300">
                💡 <strong>Compounding Calibration:</strong> Applying 20% quarterly yields 2.07× annual expansion ($1.20^4$), preserving sustainable staffing runway versus uncalibrated monthly compounding ($1.20^{12} = 8.92\times$).
              </div>
            </div>

            {/* Workload Technical Demand Output */}
            <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 lg:col-span-2 space-y-4">
              <h3 className="text-base font-black text-white flex items-center gap-2">
                <Server size={18} className="text-indigo-400" />
                Technical Workload Demand (Month {horizon})
              </h3>
              <p className="text-xs text-zinc-400">Derived from user interaction density: search, listings, image uploads, and background jobs.</p>

              <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
                <div className="p-4 bg-zinc-900/50 rounded-2xl border border-zinc-800">
                  <span className="text-[10px] text-zinc-500 uppercase font-extrabold">Peak RPS Target</span>
                  <div className="text-lg font-black font-mono text-orange-400 mt-1">{fmtNum(latestMonth.peak_rps)}</div>
                  <span className="text-[10px] text-zinc-400">3x Peak/Avg ratio</span>
                </div>
                <div className="p-4 bg-zinc-900/50 rounded-2xl border border-zinc-800">
                  <span className="text-[10px] text-zinc-500 uppercase font-extrabold">Daily API Requests</span>
                  <div className="text-lg font-black font-mono text-white mt-1">{fmtNum(latestMonth.daily_requests)}</div>
                  <span className="text-[10px] text-zinc-400">Across sessions</span>
                </div>
                <div className="p-4 bg-zinc-900/50 rounded-2xl border border-zinc-800">
                  <span className="text-[10px] text-zinc-500 uppercase font-extrabold">Object Storage</span>
                  <div className="text-lg font-black font-mono text-blue-400 mt-1">{latestMonth.accumulated_storage_tb?.toFixed(1) || '0'} TB</div>
                  <span className="text-[10px] text-zinc-400">Images + 3x backup</span>
                </div>
                <div className="p-4 bg-zinc-900/50 rounded-2xl border border-zinc-800">
                  <span className="text-[10px] text-zinc-500 uppercase font-extrabold">CDN Egress Bandwidth</span>
                  <div className="text-lg font-black font-mono text-emerald-400 mt-1">{latestMonth.monthly_bandwidth_tb?.toFixed(1) || '0'} TB</div>
                  <span className="text-[10px] text-zinc-400">Monthly egress</span>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3 text-xs font-mono pt-2">
                <div className="flex justify-between p-3 bg-zinc-900/30 rounded-xl border border-zinc-800">
                  <span className="text-zinc-400">Database Read IOPS:</span>
                  <span className="text-white font-bold">{fmtNum(latestMonth.db_read_iops)} IOPS</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/30 rounded-xl border border-zinc-800">
                  <span className="text-zinc-400">Database Write IOPS:</span>
                  <span className="text-white font-bold">{fmtNum(latestMonth.db_write_iops)} IOPS</span>
                </div>
              </div>
            </div>
          </div>

          {/* Dedicated 20 Million DAU Capacity Stress Planner */}
          <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <div className="flex items-center gap-2">
                  <h3 className="text-base font-black text-white">20,000,000 DAU Capacity Stress Assessment</h3>
                  <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-purple-950 text-purple-300 border border-purple-500/30">
                    High-Watermark Stress Test
                  </span>
                </div>
                <p className="text-xs text-zinc-400">
                  Capacity planning analysis for 20M daily active users (~80M MAU). Note: Single-server or simple cloud setups cannot support this scale.
                </p>
              </div>
              <div className="text-right">
                <span className="text-xs text-zinc-400">Est. Monthly Cloud Budget</span>
                <div className="text-lg font-black font-mono text-orange-400">
                  {fmtETB(capacity20M?.estimated_monthly_cloud_etb)}
                </div>
              </div>
            </div>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 pt-2">
              <div className="p-4 bg-zinc-900/40 rounded-2xl border border-zinc-800">
                <span className="text-[10px] text-zinc-500 uppercase font-extrabold">App Server Cores</span>
                <div className="text-xl font-black font-mono text-white mt-1">{capacity20M?.required_app_cores || '240+'} Cores</div>
                <span className="text-[10px] text-zinc-400">vCPU with 50% buffer</span>
              </div>
              <div className="p-4 bg-zinc-900/40 rounded-2xl border border-zinc-800">
                <span className="text-[10px] text-zinc-500 uppercase font-extrabold">Memory Footprint</span>
                <div className="text-xl font-black font-mono text-white mt-1">{capacity20M?.required_ram_gb || '960+'} GB RAM</div>
                <span className="text-[10px] text-zinc-400">App + Redis caching</span>
              </div>
              <div className="p-4 bg-zinc-900/40 rounded-2xl border border-zinc-800">
                <span className="text-[10px] text-zinc-500 uppercase font-extrabold">DB Read IOPS Needed</span>
                <div className="text-xl font-black font-mono text-indigo-400 mt-1">{fmtNum(capacity20M?.database_read_iops)}</div>
                <span className="text-[10px] text-zinc-400">Requires read replicas</span>
              </div>
              <div className="p-4 bg-zinc-900/40 rounded-2xl border border-zinc-800">
                <span className="text-[10px] text-zinc-500 uppercase font-extrabold">CDN Egress Demand</span>
                <div className="text-xl font-black font-mono text-emerald-400 mt-1">{capacity20M?.monthly_egress_bandwidth_tb?.toFixed(0) || '1,800'} TB</div>
                <span className="text-[10px] text-zinc-400">Per month bandwidth</span>
              </div>
            </div>

            {/* Architecture Bottlenecks and Mitigations */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
              <div className="p-5 bg-red-950/10 border border-red-500/20 rounded-2xl space-y-2">
                <h4 className="text-xs font-black text-red-400 uppercase tracking-wider flex items-center gap-2">
                  <AlertTriangle size={14} /> Architectural Bottlenecks at 20M DAU
                </h4>
                <ul className="text-xs text-zinc-300 space-y-1.5 list-disc pl-4">
                  {capacity20M?.technical_bottlenecks?.map((b: string, i: number) => (
                    <li key={i}>{b}</li>
                  )) || (
                    <>
                      <li>Single-instance PostgreSQL cannot sustain ~250k+ write IOPS and ~2M+ read IOPS.</li>
                      <li>Peak RPS requires Cloudflare CDN edge caching with &gt;85% hit ratio to protect origins.</li>
                      <li>Centralized Redis must be migrated to a partitioned Redis Cluster with replica read offload.</li>
                    </>
                  )}
                </ul>
              </div>

              <div className="p-5 bg-emerald-950/10 border border-emerald-500/20 rounded-2xl space-y-2">
                <h4 className="text-xs font-black text-emerald-400 uppercase tracking-wider flex items-center gap-2">
                  <CheckCircle2 size={14} /> Mandatory Scale Engineering Measures
                </h4>
                <ul className="text-xs text-zinc-300 space-y-1.5 list-disc pl-4">
                  {capacity20M?.required_scale_measures?.map((m: string, i: number) => (
                    <li key={i}>{m}</li>
                  )) || (
                    <>
                      <li>Deploy PostgreSQL read replicas for read-heavy search and catalog listing queries.</li>
                      <li>Implement horizontal connection pooling (PgBouncer) with transaction-level pooling.</li>
                      <li>Conduct rigorous end-to-end distributed load tests before claiming production readiness.</li>
                    </>
                  )}
                </ul>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* TAB CONTENT 3: Revenue Streams & Escrow Separation */}
      {activeTab === 'revenue' && (
        <div className="space-y-6">
          {/* Escrow Segregation Notice */}
          <div className="p-5 bg-blue-950/20 border border-blue-500/30 rounded-3xl flex items-start gap-3">
            <Scale size={20} className="text-blue-400 shrink-0 mt-0.5" />
            <div className="space-y-1 text-xs">
              <h4 className="font-bold text-white">Statutory Separation of Marketplace Escrow from Company Revenue</h4>
              <p className="text-blue-200/80 leading-relaxed">
                Gross Merchandise Volume (GMV) represents customer payments held in trust for third-party sellers. Under Ethiopian commercial regulations, seller funds are escrow liabilities and must never be recorded as HUNTIQ platform revenue. Only explicit broker commissions, listing fees, and membership subscriptions constitute recognizable platform revenue.
              </p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 space-y-4">
              <h3 className="text-base font-black text-white">Revenue Channels Breakdown (Month {latestMonth.month_index})</h3>
              <div className="space-y-3 font-mono text-xs">
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">1. Broker Commission (2.0%):</span>
                  <span className="font-bold text-emerald-400">{fmtETB(latestMonth.commission_revenue)}</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">2. Paid Listing Fees (ETB 50/item):</span>
                  <span className="font-bold text-emerald-400">{fmtETB(latestMonth.listing_fee_revenue)}</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">3. Membership Tiers (Gold/Premium):</span>
                  <span className="font-bold text-emerald-400">{fmtETB(latestMonth.subscription_revenue)}</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">4. Featured & Sponsored Listings:</span>
                  <span className="font-bold text-emerald-400">{fmtETB(latestMonth.promoted_listing_revenue)}</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">5. Verified Business Subscriptions:</span>
                  <span className="font-bold text-emerald-400">{fmtETB(latestMonth.business_tools_revenue)}</span>
                </div>
                <div className="flex justify-between p-3 bg-red-950/20 rounded-xl border border-red-500/20 text-red-300">
                  <span className="font-sans">Less: Discounts & Refund Allowance:</span>
                  <span>- {fmtETB(latestMonth.discounts_and_refunds)}</span>
                </div>
                <div className="flex justify-between p-3 bg-emerald-950/30 rounded-xl border border-emerald-500/30 text-emerald-300 font-bold text-sm">
                  <span className="font-sans">Net Recognized Revenue:</span>
                  <span>{fmtETB(latestMonth.net_recognized_revenue)}</span>
                </div>
              </div>
            </div>

            <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 space-y-4">
              <h3 className="text-base font-black text-white">Marketplace Activity & Density Metrics</h3>
              <div className="space-y-3 font-mono text-xs">
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">Active Sellers:</span>
                  <span className="font-bold text-white">{fmtNum(latestMonth.mau * 0.15)} sellers (15% of MAU)</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">Completed Transactions:</span>
                  <span className="font-bold text-white">{fmtNum(latestMonth.completed_transactions)} orders</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">Average Transaction Value (ATV):</span>
                  <span className="font-bold text-white">ETB 4,500.00</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">ARPU (Revenue per Active User):</span>
                  <span className="font-bold text-emerald-400">{fmtETB(latestMonth.revenue_per_active_user)}</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">Infra Cost per Completed Order:</span>
                  <span className="font-bold text-zinc-300">{fmtETB(latestMonth.infra_cost_per_transaction)}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* TAB CONTENT 4: OpEx & Cloud Infrastructure */}
      {activeTab === 'opex' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Cloud Infrastructure Categories */}
            <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 space-y-4">
              <h3 className="text-base font-black text-white flex items-center gap-2">
                <Server size={18} className="text-indigo-400" />
                Cloud Infrastructure Cost Breakdown (Month {latestMonth.month_index})
              </h3>
              <p className="text-xs text-zinc-400">Includes autoscaling compute, database clusters, S3 object storage, and CDN traffic.</p>

              <div className="space-y-2.5 font-mono text-xs">
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">Compute & Load Balancers:</span>
                  <span className="text-white font-bold">{fmtETB(latestMonth.compute_and_network_cost)}</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">PostgreSQL DB Clusters & Storage:</span>
                  <span className="text-white font-bold">{fmtETB(latestMonth.database_cost)}</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">Object Storage (Images) & CDN Egress:</span>
                  <span className="text-white font-bold">{fmtETB(latestMonth.storage_and_cdn_cost)}</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">Security, WAF, SMS OTP & Auth:</span>
                  <span className="text-white font-bold">{fmtETB(latestMonth.security_and_auth_cost)}</span>
                </div>
                <div className="flex justify-between p-3 bg-indigo-950/20 rounded-xl border border-indigo-500/20 text-indigo-300 font-bold">
                  <span className="font-sans">Total Monthly Cloud Hosting (+20% reserve):</span>
                  <span>{fmtETB(latestMonth.total_infrastructure_cost)}</span>
                </div>
              </div>
            </div>

            {/* Staffing and Business OpEx */}
            <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 space-y-4">
              <h3 className="text-base font-black text-white flex items-center gap-2">
                <Users size={18} className="text-emerald-400" />
                Staffing & Organizational OpEx (Month {latestMonth.month_index})
              </h3>
              <p className="text-xs text-zinc-400">Engineering compensation, statutory employer pension, office, SaaS, and marketing.</p>

              <div className="space-y-2.5 font-mono text-xs">
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">Gross Salaries (Engineering, Ops, Support):</span>
                  <span className="text-white font-bold">{fmtETB(latestMonth.staff_salaries_cost)}</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">Employer Pension Contribution (11% Proc. 715/2011):</span>
                  <span className="text-amber-300 font-bold">{fmtETB(latestMonth.employer_pension_cost)}</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">Marketing & User Acquisition (CAC):</span>
                  <span className="text-white font-bold">{fmtETB(latestMonth.marketing_expense)}</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">Office Rent, Legal, Compliance & Insurance:</span>
                  <span className="text-white font-bold">{fmtETB(latestMonth.office_admin_legal_cost)}</span>
                </div>
                <div className="flex justify-between p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                  <span className="font-sans text-zinc-400">Software Tools, SaaS & Equipment Depreciation:</span>
                  <span className="text-white font-bold">{fmtETB(latestMonth.software_saas_cost + latestMonth.accounting_depreciation)}</span>
                </div>
                <div className="flex justify-between p-3 bg-red-950/20 rounded-xl border border-red-500/20 text-red-300 font-bold">
                  <span className="font-sans">Total Monthly Organizational OpEx:</span>
                  <span>{fmtETB(latestMonth.total_operating_expense - latestMonth.total_infrastructure_cost)}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* TAB CONTENT 5: Ethiopian Statutory Tax Estimator */}
      {activeTab === 'tax' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 space-y-4">
              <h3 className="text-base font-black text-white flex items-center gap-2">
                <Calculator size={18} className="text-amber-400" />
                Tax Calculation Schedule
              </h3>
              <p className="text-xs text-zinc-400">Statutory tax liabilities calculated under Ethiopian proclamations.</p>

              <div className="space-y-3 font-mono text-xs">
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="font-sans text-zinc-400">Accounting PBT:</span>
                  <span className="text-white font-bold">{fmtETB(latestMonth.accounting_profit_before_tax)}</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="font-sans text-zinc-400">Non-Deductible Addback (5%):</span>
                  <span className="text-amber-300 font-bold">+ {fmtETB(latestMonth.taxable_income_adjustments)}</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="font-sans text-zinc-400">Loss Relief Applied (Art. 26):</span>
                  <span className="text-emerald-400 font-bold">- {fmtETB(latestMonth.loss_relief_applied)}</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="font-sans text-zinc-400">Adjusted Taxable Income:</span>
                  <span className="text-white font-bold">{fmtETB(latestMonth.taxable_income)}</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900 bg-amber-950/20 px-2 rounded-lg">
                  <span className="font-sans text-amber-300 font-bold">Estimated Corporate Tax (30%):</span>
                  <span className="text-amber-300 font-bold">{fmtETB(latestMonth.estimated_income_tax)}</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="font-sans text-zinc-400">Output VAT on Platform Fees (15%):</span>
                  <span className="text-white font-bold">{fmtETB(latestMonth.output_vat)}</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="font-sans text-zinc-400">Eligible Input VAT Credits (85%):</span>
                  <span className="text-emerald-400 font-bold">- {fmtETB(latestMonth.eligible_input_vat)}</span>
                </div>
                <div className="flex justify-between py-2 bg-amber-950/30 px-2 rounded-lg text-amber-300 font-bold">
                  <span className="font-sans">Net VAT Payable (Monthly):</span>
                  <span>{fmtETB(latestMonth.net_vat_payable)}</span>
                </div>
              </div>
            </div>

            {/* Official Tax Rules Catalog */}
            <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 lg:col-span-2 space-y-4">
              <h3 className="text-base font-black text-white flex items-center gap-2">
                <FileText size={18} className="text-indigo-400" />
                Ethiopian Statutory Tax Rules Catalog
              </h3>
              <p className="text-xs text-zinc-400">
                Official statutory provisions with enacted legal references. All rules are marked requiring independent confirmation by an Ethiopian certified tax accountant.
              </p>

              <div className="space-y-3 max-h-[380px] overflow-y-auto pr-1">
                {taxRules.map((rule) => (
                  <div key={rule.rule_code} className="p-4 bg-zinc-900/40 rounded-2xl border border-zinc-800 space-y-2">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <span className="font-mono font-bold text-xs text-white">{rule.rule_name}</span>
                        <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-zinc-800 text-zinc-400">
                          {rule.rule_code}
                        </span>
                      </div>
                      <span className="text-xs font-mono font-bold text-amber-400">
                        {rule.standard_rate ? `${(rule.standard_rate * 100).toFixed(0)}%` : `ETB ${Number(rule.threshold_amount).toLocaleString()}`}
                      </span>
                    </div>
                    <div className="text-[11px] text-zinc-400">
                      <span className="font-semibold text-zinc-300">Enacted Source:</span> {rule.source_citation}
                    </div>
                    <div className="flex items-center justify-between text-[10px] pt-1 border-t border-zinc-800/60">
                      <span className="text-amber-400/90 font-medium">⚠️ {rule.reviewer_notes}</span>
                      <span className="text-zinc-500 font-mono">Effective: {rule.effective_date}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* TAB CONTENT 6: Cash Flow, Runway & Break-Even */}
      {activeTab === 'cashflow' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 space-y-4">
              <h3 className="text-base font-black text-white">Monthly Cash Flow Movement</h3>
              <div className="space-y-3 font-mono text-xs">
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="font-sans text-zinc-400">Cash Inflows (Revenue):</span>
                  <span className="text-emerald-400 font-bold">+ {fmtETB(latestMonth.operating_cash_inflow)}</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="font-sans text-zinc-400">Cash Outflows (OpEx + Cloud):</span>
                  <span className="text-red-400 font-bold">- {fmtETB(latestMonth.operating_cash_outflow)}</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="font-sans text-zinc-400">Cash Tax Outflow (VAT/CIT):</span>
                  <span className="text-amber-300 font-bold">- {fmtETB(latestMonth.tax_cash_outflow)}</span>
                </div>
                <div className="flex justify-between py-2 border-b border-zinc-900">
                  <span className="font-sans text-zinc-400">CapEx Investments:</span>
                  <span className="text-white font-bold">{fmtETB(latestMonth.capex_cash_outflow)}</span>
                </div>
                <div className="flex justify-between py-2 bg-zinc-900/50 px-2 rounded-lg text-white font-bold">
                  <span className="font-sans">Net Cash Change:</span>
                  <span className={latestMonth.net_monthly_cash_change >= 0 ? 'text-emerald-400' : 'text-red-400'}>
                    {fmtETB(latestMonth.net_monthly_cash_change)}
                  </span>
                </div>
                <div className="flex justify-between py-2 bg-blue-950/20 px-2 rounded-lg text-blue-300 font-bold text-sm">
                  <span className="font-sans">Closing Cash Balance:</span>
                  <span>{fmtETB(latestMonth.closing_cash_balance)}</span>
                </div>
              </div>
            </div>

            <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 lg:col-span-2 space-y-4">
              <h3 className="text-base font-black text-white">Break-Even & Runway Analysis</h3>
              <p className="text-xs text-zinc-400">Contribution margin analysis decomposes costs into fixed vs variable elements to determine revenue targets.</p>

              <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
                <div className="p-4 bg-zinc-900/50 rounded-2xl border border-zinc-800">
                  <span className="text-[10px] text-zinc-500 uppercase font-extrabold">Contribution Margin Ratio</span>
                  <div className="text-lg font-black font-mono text-emerald-400 mt-1">
                    {((latestMonth.contribution_margin_ratio || 0.4) * 100).toFixed(1)}%
                  </div>
                  <span className="text-[10px] text-zinc-400">(Revenue - Variable Costs) / Rev</span>
                </div>
                <div className="p-4 bg-zinc-900/50 rounded-2xl border border-zinc-800">
                  <span className="text-[10px] text-zinc-500 uppercase font-extrabold">Monthly Break-Even Target</span>
                  <div className="text-lg font-black font-mono text-white mt-1">
                    {fmtETB(latestMonth.monthly_break_even_revenue)}
                  </div>
                  <span className="text-[10px] text-zinc-400">Fixed Costs / Contribution Ratio</span>
                </div>
                <div className="p-4 bg-zinc-900/50 rounded-2xl border border-zinc-800">
                  <span className="text-[10px] text-zinc-500 uppercase font-extrabold">Projected Cash Runway</span>
                  <div className="text-lg font-black font-mono text-blue-400 mt-1">
                    {latestMonth.cash_runway_months === -1 ? 'Self-Sustaining' : `${latestMonth.cash_runway_months?.toFixed(1)} Months`}
                  </div>
                  <span className="text-[10px] text-zinc-400">Available cash / burn rate</span>
                </div>
              </div>

              <div className="p-4 bg-zinc-900/30 rounded-2xl border border-zinc-800 space-y-2 text-xs">
                <div className="flex justify-between">
                  <span className="text-zinc-400">First Break-Even Month Milestone:</span>
                  <span className="font-bold text-emerald-400">{milestones.break_even_period} (Month {milestones.break_even_month_index})</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-400">Peak Funding Need (Max Drawdown):</span>
                  <span className="font-bold text-orange-400">{fmtETB(milestones.peak_funding_need_etb)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-400">Lowest Liquidity Point:</span>
                  <span className="font-bold text-blue-400">{fmtETB(milestones.lowest_cash_balance)} ({milestones.lowest_cash_month})</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* TAB CONTENT 7: Scenarios & Sensitivity Shocks */}
      {activeTab === 'sensitivity' && (
        <div className="space-y-6">
          <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 space-y-4">
            <h3 className="text-base font-black text-white">Adversarial Sensitivity Shocks & Stress Tests</h3>
            <p className="text-xs text-zinc-400">
              Evaluates discrete macro shocks against the planning baseline to measure solvency risk, runway erosion, and capital requirements.
            </p>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
              {report?.SensitivityShocks?.map((shock: any, idx: number) => (
                <div key={idx} className="p-5 bg-zinc-900/40 rounded-2xl border border-zinc-800 space-y-3">
                  <div className="flex items-center justify-between">
                    <h4 className="font-bold text-sm text-white flex items-center gap-2">
                      <ShieldAlert size={16} className="text-orange-400" />
                      {shock.shock_name}
                    </h4>
                    <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-zinc-800 text-zinc-300">
                      Break-Even: M{shock.shock_break_even_mo > 0 ? shock.shock_break_even_mo : 'Delayed'}
                    </span>
                  </div>
                  <p className="text-xs text-zinc-400 leading-relaxed">{shock.description}</p>
                  
                  <div className="grid grid-cols-2 gap-2 text-xs font-mono pt-1">
                    <div className="p-2.5 bg-zinc-900/60 rounded-xl">
                      <span className="text-[10px] text-zinc-500 uppercase font-sans">Year 1 Net Profit</span>
                      <div className={`font-bold mt-0.5 ${shock.shock_year_1_profit >= 0 ? 'text-emerald-400' : 'text-red-400'}`}>
                        {fmtETB(shock.shock_year_1_profit)}
                      </div>
                    </div>
                    <div className="p-2.5 bg-zinc-900/60 rounded-xl">
                      <span className="text-[10px] text-zinc-500 uppercase font-sans">Survival Capital Needed</span>
                      <div className="font-bold text-orange-400 mt-0.5">
                        {fmtETB(shock.survival_cash_needed_etb)}
                      </div>
                    </div>
                  </div>

                  <p className="text-[11px] text-amber-300/90 font-medium bg-amber-950/20 p-2 rounded-lg border border-amber-500/20">
                    💡 {shock.key_takeaway}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* TAB CONTENT 8: Assumptions & Audit Trail */}
      {activeTab === 'assumptions' && (
        <div className="space-y-6">
          <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-base font-black text-white">Financial Assumptions & Configuration Parameters</h3>
                <p className="text-xs text-zinc-400">Baseline values driving growth, workload demand, monetization, tax rates, and payroll costs.</p>
              </div>
              {canManageFinance ? (
                <button
                  onClick={handleResetAssumptions}
                  className="flex items-center gap-2 px-3.5 py-2 bg-red-950/30 hover:bg-red-900/40 text-red-300 border border-red-500/30 rounded-2xl text-xs font-bold transition-all"
                >
                  <RotateCcw size={14} />
                  Reset {scenario} to Defaults
                </button>
              ) : (
                <div className="flex items-center gap-2 px-3 py-1.5 bg-zinc-900 border border-zinc-800 rounded-xl text-zinc-400 text-xs">
                  <Shield size={14} className="text-amber-400" />
                  <span>Read-Only Mode (Owner / CFO only)</span>
                </div>
              )}
            </div>

            {/* Audit Trail Table */}
            <div className="pt-4 border-t border-zinc-900">
              <h4 className="text-xs font-extrabold text-zinc-400 uppercase tracking-wider mb-3">Audit Trail of Assumptions Changes</h4>
              {auditLogs.length === 0 ? (
                <p className="text-xs text-zinc-500 py-4 text-center">No manual assumption changes recorded yet. System using verified baseline defaults.</p>
              ) : (
                <div className="divide-y divide-zinc-900 text-xs">
                  {auditLogs.map((log) => (
                    <div key={log.id} className="py-3 flex items-center justify-between">
                      <div>
                        <span className="font-mono font-bold text-white">{log.action}</span>
                        <span className="text-zinc-500 ml-2">Scenario: {log.scenario_id || 'base'}</span>
                      </div>
                      <div className="text-zinc-500 font-mono text-[11px]">{new Date(log.created_at).toLocaleString()}</div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* TAB CONTENT 9: Historical Actuals */}
      {activeTab === 'actuals' && (
        <div className="space-y-6">
          <div className="bg-zinc-950 border border-zinc-900 rounded-3xl p-6 space-y-4">
            <h3 className="text-base font-black text-white">Historical Accounting Records (Actuals vs Projections)</h3>
            <p className="text-xs text-zinc-400">
              Actual financial results recorded in the marketplace database. When historical data is recorded, it serves as the grounding benchmark against forecast estimates.
            </p>

            {actuals.length === 0 ? (
              <div className="p-8 border border-zinc-900 rounded-2xl bg-zinc-900/20 text-center space-y-2">
                <FileText size={28} className="mx-auto text-zinc-600" />
                <h4 className="text-xs font-bold text-zinc-400">Zero Historical Accounting Actuals In Database</h4>
                <p className="text-[11px] text-zinc-500 max-w-md mx-auto">
                  The local development database currently has no closed accounting months. The FP&A engine faithfully reports zero actuals rather than fabricating artificial accounting numbers.
                </p>
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-xs text-left">
                  <thead>
                    <tr className="border-b border-zinc-900 text-zinc-500 uppercase font-extrabold">
                      <th className="p-3">Period</th>
                      <th className="p-3 text-right">Transactions</th>
                      <th className="p-3 text-right">GMV</th>
                      <th className="p-3 text-right">Commission</th>
                      <th className="p-3 text-right">Infra Actual</th>
                      <th className="p-3 text-right">Staff Actual</th>
                      <th className="p-3 text-right">Closing Cash</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-zinc-900 font-mono text-zinc-300">
                    {actuals.map((act) => (
                      <tr key={act.id} className="hover:bg-zinc-900/30">
                        <td className="p-3 font-sans font-bold text-white">{act.period_month}</td>
                        <td className="p-3 text-right">{fmtNum(act.completed_transactions)}</td>
                        <td className="p-3 text-right">{fmtETB(act.gross_merchandise_volume)}</td>
                        <td className="p-3 text-right text-emerald-400">{fmtETB(act.commission_revenue)}</td>
                        <td className="p-3 text-right text-red-300">{fmtETB(act.infra_cost_actual)}</td>
                        <td className="p-3 text-right text-red-300">{fmtETB(act.staffing_cost_actual)}</td>
                        <td className="p-3 text-right text-blue-400">{fmtETB(act.closing_cash_balance)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
