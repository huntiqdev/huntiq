'use client';

import React, { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import { Search, Bell, Plus, Building2, ChevronDown, Menu, ShieldCheck, ShieldAlert, Monitor, Terminal, Lock, ExternalLink, X, RefreshCw } from 'lucide-react';
import { GlobalSearchModal } from './GlobalSearchModal';
import { NotificationDrawer } from './NotificationDrawer';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';

interface AdminHeaderTopBarProps {
  onOpenWizard: () => void;
  selectedBranch: string;
  onBranchChange: (branch: string) => void;
  onToggleMobileSidebar?: () => void;
}

interface SecurityTelemetry {
  admin_id: number;
  admin_name?: string;
  name?: string;
  admin_role: string;
  client_ip?: string;
  ip_address?: string;
  client_device_id?: string;
  device_id?: string;
  device_name?: string;
  is_trusted_device?: boolean;
  is_trusted?: boolean;
  shift_code?: string;
  workstation_code?: string;
  geo_location?: string;
  session_id?: number;
  active_sessions_count?: number;
  security_zones: Array<{
    zone_code: string;
    zone_name: string;
    granted_at?: string;
  }>;
  evaluated_at?: string;
}

const BRANCHES = [
  { id: 'all',          name: 'All Regions',             badge: 'Global' },
  { id: 'addis-metro',  name: 'Addis Ababa',             badge: 'Metro' },
  { id: 'bole-hub',     name: 'Bole District',           badge: 'Hub' },
  { id: 'hawassa',      name: 'Hawassa',                 badge: 'South' },
  { id: 'dire-dawa',    name: 'Dire Dawa',               badge: 'East' },
];

export function AdminHeaderTopBar({
  onOpenWizard,
  selectedBranch,
  onBranchChange,
  onToggleMobileSidebar,
}: AdminHeaderTopBarProps) {
  const { user } = useAuth();
  const [isSearchOpen, setIsSearchOpen] = useState(false);
  const [isNotifOpen, setIsNotifOpen] = useState(false);
  const [isBranchDropdownOpen, setIsBranchDropdownOpen] = useState(false);
  const [isTelemetryOpen, setIsTelemetryOpen] = useState(false);
  const [telemetry, setTelemetry] = useState<SecurityTelemetry | null>(null);
  const [loadingTelemetry, setLoadingTelemetry] = useState(false);

  const telemetryDropdownRef = useRef<HTMLDivElement>(null);

  const currentBranchObj = BRANCHES.find((b) => b.id === selectedBranch) || BRANCHES[0];

  const fetchTelemetry = async () => {
    if (!user?.token) return;
    try {
      setLoadingTelemetry(true);
      const data = await api.adminGetSecurityTelemetryMe(user.token);
      if (data && typeof data === 'object') {
        setTelemetry(data);
      }
    } catch {
      // Fallback telemetry display
    } finally {
      setLoadingTelemetry(false);
    }
  };

  useEffect(() => {
    fetchTelemetry();
    const interval = setInterval(fetchTelemetry, 30000);
    return () => clearInterval(interval);
  }, [user?.token]);

  // Click outside to close telemetry dropdown
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (telemetryDropdownRef.current && !telemetryDropdownRef.current.contains(event.target as Node)) {
        setIsTelemetryOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <>
      <header className="sticky top-0 z-30 w-full bg-[#091426]/90 backdrop-blur-xl border-b border-slate-800/80 px-3 sm:px-6 lg:px-8 py-2.5 transition-all">
        <div className="flex items-center justify-between gap-3">
          
          {/* Left: Mobile Drawer Trigger + Sector / Branch Switcher */}
          <div className="flex items-center gap-2 sm:gap-3">
            {onToggleMobileSidebar && (
              <button
                type="button"
                onClick={onToggleMobileSidebar}
                className="lg:hidden p-2 rounded-xl border border-slate-800 bg-[#0E1B33] text-slate-300 hover:text-white hover:bg-slate-800 transition-colors tap-bounce active:scale-95"
                aria-label="Toggle navigation menu"
              >
                <Menu size={18} />
              </button>
            )}

            <div className="relative">
              <button
                onClick={() => setIsBranchDropdownOpen(!isBranchDropdownOpen)}
                className="flex items-center gap-2.5 px-3 py-1.5 sm:px-3.5 sm:py-2 bg-[#0E1B33] hover:bg-slate-800/90 border border-slate-800 rounded-xl text-slate-200 text-xs font-bold transition-all shadow-sm group"
              >
                <div className="p-1 bg-[#1769E0]/15 text-[#38A5FF] rounded-md group-hover:bg-[#1769E0] group-hover:text-white transition-colors">
                  <Building2 size={15} />
                </div>
                <div className="text-left hidden sm:block">
                  <div className="text-[10px] text-slate-400 uppercase tracking-widest leading-none font-bold">Scope</div>
                  <div className="text-xs text-white font-extrabold flex items-center gap-1.5 mt-0.5">
                    {currentBranchObj.name}
                    <ChevronDown size={13} className="text-slate-400" />
                  </div>
                </div>
              </button>

              {/* Branch Dropdown */}
              {isBranchDropdownOpen && (
                <div className="absolute top-full left-0 mt-2 w-72 bg-[#0E1B33] border border-slate-800 rounded-2xl shadow-2xl p-2 z-50 animate-in fade-in slide-in-from-top-2">
                  <div className="px-3 py-2 text-[10px] font-bold text-slate-400 uppercase tracking-wider border-b border-slate-800">
                    Switch Regional Scope
                  </div>
                  <div className="py-1 space-y-1">
                    {BRANCHES.map((b) => (
                      <button
                        key={b.id}
                        onClick={() => {
                          onBranchChange(b.id);
                          setIsBranchDropdownOpen(false);
                        }}
                        className={`w-full flex items-center justify-between p-2.5 rounded-xl text-xs font-semibold text-left transition-all ${
                          selectedBranch === b.id
                            ? 'bg-[#1769E0]/20 text-[#38A5FF] border border-[#1769E0]/30'
                            : 'text-slate-300 hover:bg-slate-800'
                        }`}
                      >
                        <span className="truncate">{b.name}</span>
                        <span className="text-[10px] text-slate-400 font-mono shrink-0 ml-2">{b.badge}</span>
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Center: Command Palette Global Search Trigger */}
          <div className="flex-1 max-w-md hidden md:block">
            <button
              onClick={() => setIsSearchOpen(true)}
              className="w-full flex items-center justify-between px-4 py-2 bg-[#0E1B33]/80 hover:bg-[#0E1B33] border border-slate-800/80 rounded-xl text-slate-400 text-xs font-medium transition-all group"
            >
              <div className="flex items-center gap-2.5">
                <Search size={15} className="text-slate-400 group-hover:text-[#38A5FF] transition-colors" />
                <span>Search Products, Properties, Users, Audit...</span>
              </div>
              <div className="flex items-center gap-1">
                <kbd className="px-2 py-0.5 text-[10px] font-mono text-slate-400 bg-slate-800 border border-slate-700 rounded-md">
                  ⌘K
                </kbd>
              </div>
            </button>
          </div>

          {/* Right: Security Telemetry Badge + Quick Action + Notifications */}
          <div className="flex items-center gap-2 sm:gap-3">
            
            {/* Real-time Sub-Security Audit Telemetry Badge */}
            <div className="relative" ref={telemetryDropdownRef}>
              <button
                type="button"
                onClick={() => setIsTelemetryOpen(!isTelemetryOpen)}
                className="flex items-center gap-2 px-2.5 py-1.5 sm:px-3 sm:py-2 bg-[#071324] hover:bg-[#0b1d38] border border-cyan-500/30 hover:border-cyan-400/50 rounded-xl text-xs transition-all shadow-sm group"
                title="Sovereign Audit Telemetry & Security Zones"
              >
                <div className="relative flex items-center justify-center">
                  <span className="w-2 h-2 rounded-full bg-emerald-400 animate-ping absolute opacity-75" />
                  <span className="w-2 h-2 rounded-full bg-emerald-500 relative" />
                </div>
                <div className="hidden sm:flex items-center gap-1.5 text-[11px] font-mono">
                  <span className="text-cyan-400 font-bold">#{telemetry?.admin_id ?? user?.id ?? 0}</span>
                  <span className="text-slate-600">•</span>
                  <span className="text-amber-400 font-semibold uppercase text-[10px] tracking-wide">
                    {telemetry?.admin_role ?? user?.adminRole ?? 'ADMIN'}
                  </span>
                  <span className="text-slate-600 hidden md:inline">•</span>
                  <span className="text-indigo-300 font-mono hidden md:inline text-[10px]">
                    {telemetry?.shift_code || 'DAY'}
                  </span>
                  <span className="text-slate-600 hidden xl:inline">•</span>
                  <span className="text-emerald-400 hidden xl:inline text-[10px]">
                    {telemetry?.workstation_code || 'WS-AA-HQ'}
                  </span>
                </div>
                <ShieldCheck size={14} className="text-cyan-400 shrink-0 ml-0.5 group-hover:scale-110 transition-transform" />
              </button>

              {/* Telemetry Details Popover */}
              {isTelemetryOpen && (
                <div className="absolute top-full right-0 mt-2 w-80 sm:w-96 bg-[#0a1628] border border-cyan-500/40 rounded-2xl shadow-2xl p-4 z-50 animate-in fade-in slide-in-from-top-2">
                  <div className="flex items-center justify-between pb-3 border-b border-slate-800">
                    <div className="flex items-center gap-2">
                      <ShieldCheck size={16} className="text-cyan-400" />
                      <span className="text-xs font-bold text-white uppercase tracking-wider">
                        Sovereign Audit Telemetry
                      </span>
                    </div>
                    <div className="flex items-center gap-1">
                      <button
                        onClick={fetchTelemetry}
                        className="p-1 text-slate-400 hover:text-white rounded hover:bg-slate-800"
                        title="Refresh Telemetry"
                      >
                        <RefreshCw size={13} className={loadingTelemetry ? 'animate-spin' : ''} />
                      </button>
                      <button
                        onClick={() => setIsTelemetryOpen(false)}
                        className="p-1 text-slate-400 hover:text-white rounded hover:bg-slate-800"
                      >
                        <X size={14} />
                      </button>
                    </div>
                  </div>

                  {/* Identity Tuple Matrix */}
                  <div className="mt-3 space-y-2 font-mono text-xs">
                    <div className="p-2.5 bg-[#050e1c] rounded-xl border border-slate-800/80 flex items-center justify-between">
                      <span className="text-slate-400 text-[11px]">Administrator</span>
                      <span className="text-cyan-300 font-bold">
                        #{telemetry?.admin_id ?? user?.id} ({telemetry?.name ?? telemetry?.admin_name ?? user?.name ?? 'Admin'})
                      </span>
                    </div>

                    <div className="p-2.5 bg-[#050e1c] rounded-xl border border-slate-800/80 flex items-center justify-between">
                      <span className="text-slate-400 text-[11px]">Governance Role</span>
                      <span className="text-amber-400 font-bold uppercase text-[11px]">
                        {telemetry?.admin_role ?? user?.adminRole}
                      </span>
                    </div>

                    <div className="p-2.5 bg-[#050e1c] rounded-xl border border-slate-800/80 flex items-center justify-between">
                      <span className="text-slate-400 text-[11px]">Operational Shift</span>
                      <span className="text-indigo-400 font-bold uppercase text-[11px] px-2 py-0.5 rounded bg-indigo-950/60 border border-indigo-500/30">
                        {telemetry?.shift_code ?? 'DAY_SHIFT'}
                      </span>
                    </div>

                    <div className="p-2.5 bg-[#050e1c] rounded-xl border border-slate-800/80 flex items-center justify-between">
                      <span className="text-slate-400 text-[11px]">Terminal Workstation</span>
                      <div className="flex items-center gap-1.5">
                        <Terminal size={12} className="text-emerald-400" />
                        <span className="text-emerald-400 font-bold text-[11px]">
                          {telemetry?.workstation_code ?? 'WS-CORE-01'}
                        </span>
                      </div>
                    </div>

                    <div className="p-2.5 bg-[#050e1c] rounded-xl border border-slate-800/80 flex items-center justify-between">
                      <span className="text-slate-400 text-[11px]">Client IP Address</span>
                      <span className="text-emerald-400 font-bold">
                        {telemetry?.ip_address ?? telemetry?.client_ip ?? '127.0.0.1'}
                      </span>
                    </div>

                    <div className="p-2.5 bg-[#050e1c] rounded-xl border border-slate-800/80 flex items-center justify-between">
                      <span className="text-slate-400 text-[11px]">Device Hardware</span>
                      <div className="flex items-center gap-1.5">
                        <Monitor size={12} className={telemetry?.is_trusted_device || telemetry?.is_trusted ? 'text-emerald-400' : 'text-amber-400'} />
                        <span className={`text-[10px] font-bold ${telemetry?.is_trusted_device || telemetry?.is_trusted ? 'text-emerald-300' : 'text-slate-300'}`}>
                          {telemetry?.device_name || (telemetry?.is_trusted_device || telemetry?.is_trusted ? 'Trusted Terminal Workstation' : 'Registered Device')}
                        </span>
                      </div>
                    </div>

                    {/* Active Sub-Security Zones */}
                    <div className="p-2.5 bg-[#050e1c] rounded-xl border border-slate-800/80">
                      <div className="flex items-center justify-between mb-1.5">
                        <span className="text-slate-400 text-[11px]">Delegated Sub-Security Zones</span>
                        <span className="text-[10px] text-cyan-400 font-bold">
                          {telemetry?.security_zones?.length ?? 0} Active
                        </span>
                      </div>
                      <div className="flex flex-wrap gap-1.5">
                        {telemetry?.security_zones && telemetry.security_zones.length > 0 ? (
                          telemetry.security_zones.map((z) => (
                            <span
                              key={z.zone_code}
                              className="px-2 py-0.5 bg-cyan-950/80 border border-cyan-500/40 rounded-lg text-[10px] text-cyan-300 font-bold"
                            >
                              {z.zone_code}
                            </span>
                          ))
                        ) : (
                          <span className="text-[10px] text-slate-500 italic">No specific sub-zones assigned (Standard Scope)</span>
                        )}
                      </div>
                    </div>
                  </div>

                  {/* Footer link to full security console */}
                  <div className="mt-3 pt-2.5 border-t border-slate-800/80 flex items-center justify-between">
                    <span className="text-[10px] text-slate-500 font-mono">
                      Telemetry Session: #{telemetry?.session_id ?? 'ACTIVE'}
                    </span>
                    <Link
                      href="/security-console"
                      onClick={() => setIsTelemetryOpen(false)}
                      className="inline-flex items-center gap-1 text-[11px] font-bold text-cyan-400 hover:text-cyan-300"
                    >
                      Security Console
                      <ExternalLink size={11} />
                    </Link>
                  </div>
                </div>
              )}
            </div>

            {/* Search Trigger for Mobile */}
            <button
              onClick={() => setIsSearchOpen(true)}
              className="md:hidden p-2 bg-[#0E1B33] hover:bg-slate-800 border border-slate-800 text-slate-300 rounded-xl"
              aria-label="Search"
            >
              <Search size={18} />
            </button>

            {/* Quick Action Button */}
            <button
              onClick={onOpenWizard}
              className="flex items-center gap-2 px-3 py-2 bg-[#1769E0] hover:bg-[#1055BD] text-white rounded-xl text-xs font-extrabold shadow-sm shadow-[#1769E0]/25 transition-all tap-bounce active:scale-95"
            >
              <Plus size={16} />
              <span className="hidden sm:inline">New Listing</span>
            </button>

            {/* Notification Drawer Bell */}
            <button
              onClick={() => setIsNotifOpen(true)}
              className="relative p-2 bg-[#0E1B33] hover:bg-slate-800 border border-slate-800 text-slate-300 hover:text-white rounded-xl transition-all"
              title="Notifications & Queue"
            >
              <Bell size={18} />
              <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-[#38A5FF] rounded-full ring-2 ring-[#091426] animate-pulse" />
            </button>
          </div>
        </div>
      </header>

      {/* Global Modals */}
      <GlobalSearchModal isOpen={isSearchOpen} onClose={() => setIsSearchOpen(false)} />
      <NotificationDrawer isOpen={isNotifOpen} onClose={() => setIsNotifOpen(false)} />
    </>
  );
}
