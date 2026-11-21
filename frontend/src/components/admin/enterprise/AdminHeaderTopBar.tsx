'use client';

import React, { useState } from 'react';
import { Search, Bell, Plus, Building2, ChevronDown, Menu } from 'lucide-react';
import { GlobalSearchModal } from './GlobalSearchModal';
import { NotificationDrawer } from './NotificationDrawer';
import { LanguageSelector } from '@/components/layout/LanguageSelector';

interface AdminHeaderTopBarProps {
  onOpenWizard: () => void;
  selectedBranch: string;
  onBranchChange: (branch: string) => void;
  onToggleMobileSidebar?: () => void;
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
  const [isSearchOpen, setIsSearchOpen] = useState(false);
  const [isNotifOpen, setIsNotifOpen] = useState(false);
  const [isBranchDropdownOpen, setIsBranchDropdownOpen] = useState(false);

  const currentBranchObj = BRANCHES.find((b) => b.id === selectedBranch) || BRANCHES[0];

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

          {/* Right: Quick Action Wizard Trigger & Notification Drawer Bell */}
          <div className="flex items-center gap-2">
            {/* Search Trigger for Mobile */}
            <button
              onClick={() => setIsSearchOpen(true)}
              className="md:hidden p-2 bg-[#0E1B33] hover:bg-slate-800 border border-slate-800 text-slate-300 rounded-xl"
              aria-label="Search"
            >
              <Search size={18} />
            </button>

            {/* Language Selector */}
            <LanguageSelector variant="pill" />

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
