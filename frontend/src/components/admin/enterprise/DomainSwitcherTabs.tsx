'use client';

import React from 'react';
import { ShoppingBag, Home, ShieldCheck, FileText, ArrowUpRight, TrendingUp, Layers, CheckCircle2 } from 'lucide-react';

export type DomainMode = 'ecommerce' | 'rentals' | 'leases';

interface DomainSwitcherTabsProps {
  activeMode: DomainMode;
  onModeChange: (mode: DomainMode) => void;
  counts: {
    ecommerceCount: number;
    rentalsCount: number;
    leasesCount: number;
  };
}

export function DomainSwitcherTabs({ activeMode, onModeChange, counts }: DomainSwitcherTabsProps) {
  const tabs = [
    {
      id: 'ecommerce' as DomainMode,
      title: 'E-Commerce Marketplace',
      subtitle: 'Products, SKUs & Moderation',
      icon: ShoppingBag,
      count: counts.ecommerceCount,
      activeColor: 'from-orange-600 to-amber-600 text-white shadow-orange-600/20',
      badge: 'Escrow Active',
    },
    {
      id: 'rentals' as DomainMode,
      title: 'Property & Unit Rentals',
      subtitle: 'Daily/Monthly Rent & Grids',
      icon: Home,
      count: counts.rentalsCount,
      activeColor: 'from-blue-600 to-indigo-600 text-white shadow-blue-600/20',
      badge: 'Real-time Matrix',
    },
    {
      id: 'leases' as DomainMode,
      title: 'Landlord & Lease Contracts',
      subtitle: 'Rent Collection & Deposits',
      icon: FileText,
      count: counts.leasesCount,
      activeColor: 'from-emerald-600 to-teal-600 text-white shadow-emerald-600/20',
      badge: 'Escrow Secured',
    },
  ];

  return (
    <div className="space-y-4">
      {/* Top Banner Mode Selector */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-3 p-1.5 bg-zinc-950/80 border border-zinc-800/80 rounded-2xl">
        {tabs.map((tab) => {
          const Icon = tab.icon;
          const isActive = activeMode === tab.id;

          return (
            <button
              key={tab.id}
              onClick={() => onModeChange(tab.id)}
              className={`flex items-start justify-between p-4 rounded-xl text-left transition-all duration-200 ${
                isActive
                  ? `bg-gradient-to-r ${tab.activeColor} shadow-xl scale-[1.01]`
                  : 'bg-zinc-900/50 hover:bg-zinc-900 border border-zinc-800/60 text-zinc-400 hover:text-zinc-200'
              }`}
            >
              <div className="space-y-2 min-w-0">
                <div className="flex items-center gap-2">
                  <div className={`p-2 rounded-lg ${isActive ? 'bg-white/10 text-white' : 'bg-zinc-800 text-orange-400'}`}>
                    <Icon size={18} />
                  </div>
                  <span className={`font-black text-sm tracking-tight ${isActive ? 'text-white' : 'text-zinc-200'}`}>
                    {tab.title}
                  </span>
                </div>
                <p className={`text-xs ${isActive ? 'text-white/80' : 'text-zinc-400'}`}>
                  {tab.subtitle}
                </p>
              </div>

              <div className="flex flex-col items-end gap-1.5 shrink-0 ml-2">
                <span
                  className={`text-[10px] font-extrabold uppercase px-2 py-0.5 rounded-full border ${
                    isActive
                      ? 'bg-white/20 text-white border-white/30'
                      : 'bg-zinc-800/80 text-zinc-400 border-zinc-700/50'
                  }`}
                >
                  {tab.badge}
                </span>
                <span className={`text-lg font-black font-mono ${isActive ? 'text-white' : 'text-zinc-300'}`}>
                  {tab.count.toLocaleString()}
                </span>
              </div>
            </button>
          );
        })}
      </div>

      {/* Mode Sub-Header & Live Domain Metrics */}
      <div className="flex flex-wrap items-center justify-between gap-3 px-4 py-3 bg-zinc-900/60 border border-zinc-800/80 rounded-xl text-xs text-zinc-400">
        <div className="flex items-center gap-2">
          <Layers size={15} className="text-orange-500" />
          <span className="font-bold text-zinc-200 uppercase tracking-wider text-[11px]">
            {activeMode === 'ecommerce' && 'Active View: Buy/Sell Product Inventory & Seller Verification Desk'}
            {activeMode === 'rentals' && 'Active View: Property Listings, Furnished Status & Availability Matrix'}
            {activeMode === 'leases' && 'Active View: Landlord Agreements, Security Deposits & Maintenance Queue'}
          </span>
        </div>

        <div className="flex items-center gap-4 text-zinc-400 text-[11px] font-mono">
          <span className="flex items-center gap-1">
            <CheckCircle2 size={13} className="text-emerald-400" /> Auto-Synced API
          </span>
          <span className="flex items-center gap-1 text-orange-400 font-bold">
            <TrendingUp size={13} /> 99.8% Platform Uptime
          </span>
        </div>
      </div>
    </div>
  );
}
