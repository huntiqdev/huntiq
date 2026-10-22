'use client';

import React, { useState, useEffect, useRef, useCallback } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Search, X, Package, Home, FileText, User, ArrowRight, CornerDownLeft, Loader2 } from 'lucide-react';

export interface SearchResult {
  id: string;
  title: string;
  subtitle: string;
  category: 'product' | 'rental' | 'lease' | 'user';
  badge: string;
  badgeColor: string;
  path: string;
}

interface GlobalSearchModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSelectResult?: (result: SearchResult) => void;
}

export function GlobalSearchModal({ isOpen, onClose, onSelectResult }: GlobalSearchModalProps) {
  const { user } = useAuth();
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState<'all' | 'product' | 'rental' | 'lease' | 'user'>('all');
  const [selectedIndex, setSelectedIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);

  const fetchLiveSearchResults = useCallback(async (q: string) => {
    if (!user) return;
    setLoading(true);
    try {
      const searchItems: SearchResult[] = [];

      // 1. Products & Inventory
      const prods = await api.adminGetAllProducts(user.token);
      if (Array.isArray(prods)) {
        prods.forEach((p: any) => {
          const titleMatches = (p.title || '').toLowerCase().includes(q.toLowerCase());
          const catMatches = (p.category_id || '').toString().includes(q);
          const typeMatches = (p.listing_type || 'sell').toLowerCase().includes(q.toLowerCase());
          
          if (!q || titleMatches || catMatches || typeMatches) {
            const isRental = (p.listing_type || '').toLowerCase() === 'rent';
            const isLease = (p.listing_type || '').toLowerCase() === 'lease';
            const cat = isRental ? 'rental' : isLease ? 'lease' : 'product';

            searchItems.push({
              id: `PRD-${p.id}`,
              title: p.title || 'Untitled Product',
              subtitle: `ETB ${Number(p.price || 0).toLocaleString()} • Seller #${p.seller_id} • Status: ${p.status}`,
              category: cat,
              badge: p.status === 'approved' ? 'Active Listing' : p.status === 'pending' ? 'Pending Moderation' : 'Rejected',
              badgeColor: p.status === 'approved' 
                ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' 
                : 'bg-amber-500/10 text-amber-400 border-amber-500/20',
              path: '/admin/products'
            });
          }
        });
      }

      // 2. Sellers & Users
      const sellers = await api.adminGetSellers(user.token);
      if (Array.isArray(sellers)) {
        sellers.forEach((u: any) => {
          const nameMatches = (u.name || '').toLowerCase().includes(q.toLowerCase());
          const phoneMatches = (u.phone || '').includes(q);
          const emailMatches = (u.email || '').toLowerCase().includes(q.toLowerCase());

          if (!q || nameMatches || phoneMatches || emailMatches) {
            searchItems.push({
              id: `USR-${u.id}`,
              title: u.name || 'Platform User',
              subtitle: `${u.role.toUpperCase()} • Tier: ${(u.tier || 'normal').toUpperCase()} • ${u.phone}`,
              category: 'user',
              badge: u.is_active ? 'Active User' : 'Deactivated',
              badgeColor: u.is_active 
                ? 'bg-purple-500/10 text-purple-400 border-purple-500/20' 
                : 'bg-red-500/10 text-red-400 border-red-500/20',
              path: '/admin/users'
            });
          }
        });
      }

      setResults(searchItems.slice(0, 30));
    } catch (err) {
      console.error('Global search error:', err);
    } finally {
      setLoading(false);
    }
  }, [user]);

  useEffect(() => {
    if (isOpen) {
      setTimeout(() => inputRef.current?.focus(), 50);
      fetchLiveSearchResults(query);
    } else {
      setQuery('');
      setSelectedIndex(0);
    }
  }, [isOpen, fetchLiveSearchResults]);

  useEffect(() => {
    const timer = setTimeout(() => {
      if (isOpen) fetchLiveSearchResults(query);
    }, 200);
    return () => clearTimeout(timer);
  }, [query, isOpen, fetchLiveSearchResults]);

  const filteredResults = results.filter((item) => {
    return activeTab === 'all' || item.category === activeTab;
  });

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!isOpen) return;
      if (e.key === 'ArrowDown') {
        e.preventDefault();
        setSelectedIndex((prev) => (prev + 1) % (filteredResults.length || 1));
      } else if (e.key === 'ArrowUp') {
        e.preventDefault();
        setSelectedIndex((prev) => (prev - 1 + filteredResults.length) % (filteredResults.length || 1));
      } else if (e.key === 'Enter' && filteredResults[selectedIndex]) {
        e.preventDefault();
        handleSelect(filteredResults[selectedIndex]);
      } else if (e.key === 'Escape') {
        onClose();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, filteredResults, selectedIndex]);

  const handleSelect = (item: SearchResult) => {
    if (onSelectResult) onSelectResult(item);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center pt-16 sm:pt-24 bg-black/75 backdrop-blur-md transition-opacity">
      <div className="w-full max-w-2xl bg-zinc-900 border border-zinc-800 rounded-2xl shadow-2xl overflow-hidden flex flex-col animate-in fade-in zoom-in duration-150">
        
        {/* Input Header */}
        <div className="flex items-center px-4 py-3.5 border-b border-zinc-800 bg-zinc-950/60">
          <Search className="text-orange-500 shrink-0 mr-3" size={20} />
          <input
            ref={inputRef}
            type="text"
            value={query}
            onChange={(e) => {
              setQuery(e.target.value);
              setSelectedIndex(0);
            }}
            placeholder="Live Search Products, Rentals, Landlords, Users... (Cmd+K)"
            className="w-full bg-transparent text-white placeholder-zinc-500 text-sm focus:outline-none"
          />
          {query && (
            <button
              onClick={() => setQuery('')}
              className="text-zinc-500 hover:text-zinc-300 mr-2 p-1 rounded-md hover:bg-zinc-800"
            >
              <X size={16} />
            </button>
          )}
          <kbd className="hidden sm:inline-flex items-center gap-1 px-2 py-0.5 text-[10px] font-mono font-medium text-zinc-400 bg-zinc-800 border border-zinc-700 rounded-md">
            ESC
          </kbd>
        </div>

        {/* Filter Tabs */}
        <div className="flex items-center gap-2 px-4 py-2 bg-zinc-900 border-b border-zinc-800/80 text-xs font-semibold text-zinc-400 overflow-x-auto">
          {[
            { id: 'all', label: 'All Results' },
            { id: 'product', label: 'Products', icon: Package },
            { id: 'rental', label: 'Property Rentals', icon: Home },
            { id: 'lease', label: 'Leases', icon: FileText },
            { id: 'user', label: 'Users & Landlords', icon: User },
          ].map((tab) => {
            const Icon = tab.icon;
            const active = activeTab === tab.id;
            return (
              <button
                key={tab.id}
                onClick={() => {
                  setActiveTab(tab.id as any);
                  setSelectedIndex(0);
                }}
                className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg transition-all ${
                  active
                    ? 'bg-orange-600/15 text-orange-400 border border-orange-500/30'
                    : 'hover:bg-zinc-800/60 hover:text-zinc-200'
                }`}
              >
                {Icon && <Icon size={14} />}
                {tab.label}
              </button>
            );
          })}
        </div>

        {/* Results List */}
        <div className="max-h-96 overflow-y-auto p-2 divide-y divide-zinc-800/40">
          {loading ? (
            <div className="py-12 text-center text-zinc-500 text-sm flex items-center justify-center gap-2">
              <Loader2 className="animate-spin text-orange-500" size={18} />
              Searching database...
            </div>
          ) : filteredResults.length === 0 ? (
            <div className="py-12 text-center text-zinc-500 text-sm">
              No matching live database entries found for "{query}".
            </div>
          ) : (
            filteredResults.map((item, index) => {
              const isSelected = index === selectedIndex;
              let CategoryIcon = Package;
              if (item.category === 'rental') CategoryIcon = Home;
              if (item.category === 'lease') CategoryIcon = FileText;
              if (item.category === 'user') CategoryIcon = User;

              return (
                <div
                  key={item.id}
                  onClick={() => handleSelect(item)}
                  onMouseEnter={() => setSelectedIndex(index)}
                  className={`flex items-center justify-between p-3 rounded-xl cursor-pointer transition-all ${
                    isSelected ? 'bg-zinc-800/90 text-white' : 'text-zinc-300 hover:bg-zinc-800/50'
                  }`}
                >
                  <div className="flex items-center gap-3 min-w-0">
                    <div className={`p-2.5 rounded-xl shrink-0 ${isSelected ? 'bg-orange-500 text-white' : 'bg-zinc-800 text-zinc-400'}`}>
                      <CategoryIcon size={18} />
                    </div>
                    <div className="min-w-0">
                      <div className="flex items-center gap-2">
                        <span className="font-bold text-sm text-zinc-100 truncate">{item.title}</span>
                        <span className="text-[10px] font-mono text-zinc-500 uppercase">{item.id}</span>
                      </div>
                      <p className="text-xs text-zinc-400 truncate mt-0.5">{item.subtitle}</p>
                    </div>
                  </div>

                  <div className="flex items-center gap-3 shrink-0 ml-3">
                    <span className={`text-[11px] font-semibold px-2.5 py-1 rounded-md border ${item.badgeColor}`}>
                      {item.badge}
                    </span>
                    {isSelected && (
                      <div className="flex items-center text-xs font-semibold text-orange-400 gap-1">
                        <span>Select</span>
                        <CornerDownLeft size={14} />
                      </div>
                    )}
                  </div>
                </div>
              );
            })
          )}
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between px-4 py-2.5 bg-zinc-950 border-t border-zinc-800 text-[11px] text-zinc-500 font-mono">
          <div className="flex items-center gap-3">
            <span><kbd className="px-1.5 py-0.5 bg-zinc-800 border border-zinc-700 rounded text-zinc-300">↑↓</kbd> Navigate</span>
            <span><kbd className="px-1.5 py-0.5 bg-zinc-800 border border-zinc-700 rounded text-zinc-300">↵</kbd> Open</span>
          </div>
          <div>SellBuy Live Search Engine</div>
        </div>
      </div>
    </div>
  );
}
