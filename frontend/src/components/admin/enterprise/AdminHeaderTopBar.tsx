'use client';

import React, { useState } from 'react';
import { Search, Bell, Plus, Building2, ChevronDown, ShieldCheck, Sparkles, SlidersHorizontal, RefreshCw } from 'lucide-react';
import { GlobalSearchModal } from './GlobalSearchModal';
import { NotificationDrawer } from './NotificationDrawer';

interface AdminHeaderTopBarProps {
  onOpenWizard: () => void;
  selectedBranch: string;
  onBranchChange: (branch: string) => void;
}

const BRANCHES = [
  { id: 'all',          name: 'All Regions',             badge: 'Global' },
  { id: 'addis-metro',  name: 'Addis Ababa',             badge: 'Metro' },
  { id: 'bole-hub',     name: 'Bole District',           badge: 'Hub' },
  { id: 'hawassa',      name: 'Hawassa',                 badge: 'South' },
  { id: 'dire-dawa',    name: 'Dire Dawa',               badge: 'East' },
];

export function AdminHeaderTopBar({ onOpenWizard, selectedBranch, onBranchChange }: AdminHeaderTopBarProps) {
  const [isSearchOpen, setIsSearchOpen] = useState(false);
  const [isNotifOpen, setIsNotifOpen] = useState(false);
  const [isBranchDropdownOpen, setIsBranchDropdownOpen] = useState(false);

  const currentBranchObj = BRANCHES.find((b) => b.id === selectedBranch) || BRANCHES[0];

  return (
    <>
      <header className="sticky top-0 z-30 w-full bg-zinc-950/80 backdrop-blur-xl border-b border-zinc-800/80 px-4 lg:px-8 py-3 transition-all">
        <div className="flex items-center justify-between gap-4">
          
          {/* Left: Sector / Branch Switcher */}
          <div className="flex items-center gap-3">
            <div className="relative">
              <button
                onClick={() => setIsBranchDropdownOpen(!isBranchDropdownOpen)}
                className="flex items-center gap-2.5 px-3.5 py-2 bg-zinc-900/90 hover:bg-zinc-800/90 border border-zinc-800 rounded-xl text-zinc-200 text-xs font-bold transition-all shadow-sm group"
              >
                <div className="p-1 bg-orange-600/10 text-orange-500 rounded-md group-hover:bg-orange-500 group-hover:text-white transition-colors">
                  <Building2 size={15} />
                </div>
                <div className="text-left hidden sm:block">
                  <div className="text-[10px] text-zinc-500 uppercase tracking-widest leading-none font-bold">Active Branch Sector</div>
                  <div className="text-xs text-white font-extrabold flex items-center gap-1.5 mt-0.5">
                    {currentBranchObj.name}
                    <ChevronDown size={13} className="text-zinc-400" />
                  </div>
                </div>
              </button>

              {/* Branch Dropdown */}
              {isBranchDropdownOpen && (
                <div className="absolute top-full left-0 mt-2 w-72 bg-zinc-900 border border-zinc-800 rounded-2xl shadow-2xl p-2 z-50 animate-in fade-in slide-in-from-top-2">
                  <div className="px-3 py-2 text-[10px] font-bold text-zinc-500 uppercase tracking-wider border-b border-zinc-800">
                    Switch Property Sector Scope
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
                            ? 'bg-orange-600/15 text-orange-400 border border-orange-500/20'
                            : 'text-zinc-300 hover:bg-zinc-800'
                        }`}
                      >
                        <span className="truncate">{b.name}</span>
                        <span className="text-[10px] text-zinc-500 font-mono shrink-0 ml-2">{b.badge}</span>
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
              className="w-full flex items-center justify-between px-4 py-2 bg-zinc-900/60 hover:bg-zinc-900 border border-zinc-800/80 rounded-xl text-zinc-400 text-xs font-medium transition-all group"
            >
              <div className="flex items-center gap-2.5">
                <Search size={15} className="text-zinc-500 group-hover:text-orange-500 transition-colors" />
                <span>Search Products, Properties, Tenants, Leases...</span>
              </div>
              <div className="flex items-center gap-1">
                <kbd className="px-2 py-0.5 text-[10px] font-mono text-zinc-400 bg-zinc-800 border border-zinc-700 rounded-md">
                  ⌘K
                </kbd>
              </div>
            </button>
          </div>

          {/* Right: Quick Action Wizard Trigger & Notification Drawer Bell */}
          <div className="flex items-center gap-2.5">
            {/* Search Trigger for Mobile */}
            <button
              onClick={() => setIsSearchOpen(true)}
              className="md:hidden p-2.5 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 text-zinc-300 rounded-xl"
            >
              <Search size={18} />
            </button>

            {/* Quick Action Button */}
            <button
              onClick={onOpenWizard}
              className="flex items-center gap-2 px-3.5 py-2 bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-500 hover:to-amber-500 text-white rounded-xl text-xs font-extrabold shadow-lg shadow-orange-600/20 transition-all hover:scale-[1.02] active:scale-[0.98]"
            >
              <Plus size={16} />
              <span className="hidden sm:inline">New Listing / Lease</span>
            </button>

            {/* Notification Drawer Bell */}
            <button
              onClick={() => setIsNotifOpen(true)}
              className="relative p-2.5 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 text-zinc-300 hover:text-white rounded-xl transition-all"
              title="Notifications & Queue"
            >
              <Bell size={18} />
              <span className="absolute top-1.5 right-1.5 w-2.5 h-2.5 bg-orange-500 rounded-full ring-4 ring-zinc-950 animate-pulse" />
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
