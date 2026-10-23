'use client';

import { useState, useRef, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { CapabilityActivationModal } from '@/components/auth/CapabilityActivationModal';
import {
  User,
  ShoppingBag,
  Home,
  ChevronDown,
  LogOut,
  Sparkles,
  CheckCircle2,
  PlusCircle,
  Shield,
  Layers,
} from 'lucide-react';

export function UserAccountMenu() {
  const {
    user,
    loading,
    isSeller,
    isLandlord,
    isAdmin,
    capabilities,
    logout,
  } = useAuth();
  const router = useRouter();

  const [isOpen, setIsOpen] = useState(false);
  const [activationCap, setActivationCap] = useState<'seller' | 'landlord' | null>(null);
  const menuRef = useRef<HTMLDivElement>(null);

  // Close on outside click
  useEffect(() => {
    if (!isOpen) return;
    const handleOutsideClick = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleOutsideClick);
    return () => document.removeEventListener('mousedown', handleOutsideClick);
  }, [isOpen]);

  if (loading) {
    return <div className="w-8 h-8 rounded-full bg-zinc-200 dark:bg-zinc-800 animate-pulse" />;
  }

  if (!user) {
    return (
      <Link
        href="/auth"
        className="flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-black uppercase tracking-wider text-orange-600 dark:text-orange-400 hover:bg-orange-50 dark:hover:bg-orange-950/30 transition-all border border-orange-600/30"
      >
        <User size={14} />
        <span>Sign In</span>
      </Link>
    );
  }

  const handleSignOut = async () => {
    setIsOpen(false);
    await logout();
    router.push('/');
  };

  return (
    <>
      <div className="relative" ref={menuRef}>
        {/* User Button */}
        <button
          onClick={() => setIsOpen((prev) => !prev)}
          className="flex items-center gap-2 p-1 pl-2 pr-2.5 rounded-full bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 transition-colors border border-zinc-200 dark:border-zinc-700 text-xs font-bold"
        >
          <div className="w-6 h-6 rounded-full bg-orange-600 text-white flex items-center justify-center font-black text-xs">
            {user.role?.[0]?.toUpperCase() || 'U'}
          </div>
          <span className="hidden sm:inline text-zinc-800 dark:text-zinc-200">
            Account
          </span>
          <ChevronDown
            size={12}
            className={`text-zinc-400 transition-transform ${isOpen ? 'rotate-180' : ''}`}
          />
        </button>

        {/* Dropdown Menu */}
        {isOpen && (
          <div className="absolute right-0 mt-2 w-72 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl shadow-2xl z-50 p-3 overflow-hidden animate-in fade-in slide-in-from-top-2 duration-150">
            {/* User Header */}
            <div className="p-3 bg-zinc-50 dark:bg-white/[0.03] rounded-2xl mb-3 border border-zinc-200/60 dark:border-white/5">
              <div className="flex items-center justify-between mb-1">
                <span className="text-[10px] font-black uppercase tracking-wider text-zinc-400">
                  Unified Identity
                </span>
                <span className="text-[10px] bg-green-500/10 text-green-500 font-bold px-2 py-0.5 rounded-full flex items-center gap-1">
                  <CheckCircle2 size={10} /> Active
                </span>
              </div>
              <p className="text-xs font-bold text-zinc-900 dark:text-white truncate">
                User #{user.id}
              </p>
              <div className="flex flex-wrap gap-1 mt-2">
                {capabilities.map((cap) => (
                  <span
                    key={cap}
                    className="text-[9px] font-black uppercase tracking-wider px-2 py-0.5 rounded-lg bg-orange-500/10 text-orange-600 dark:text-orange-400 border border-orange-500/20"
                  >
                    {cap}
                  </span>
                ))}
              </div>
            </div>

            {/* Marketplace Capabilities Switcher */}
            <div className="space-y-1 mb-3">
              <p className="text-[9px] font-black uppercase tracking-widest text-zinc-400 px-3 py-1">
                Marketplace Spaces
              </p>

              {/* Buyer Portal */}
              <Link
                href="/buyer/dashboard"
                onClick={() => setIsOpen(false)}
                className="flex items-center justify-between px-3 py-2 rounded-xl text-xs font-semibold text-zinc-700 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
              >
                <div className="flex items-center gap-2.5">
                  <User size={15} className="text-zinc-400" />
                  <span>Buyer Portal</span>
                </div>
                <span className="text-[10px] text-zinc-400 font-normal">Active</span>
              </Link>

              {/* Seller Studio */}
              {isSeller ? (
                <Link
                  href="/seller/dashboard"
                  onClick={() => setIsOpen(false)}
                  className="flex items-center justify-between px-3 py-2 rounded-xl text-xs font-semibold text-zinc-700 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
                >
                  <div className="flex items-center gap-2.5">
                    <ShoppingBag size={15} className="text-orange-500" />
                    <span>Seller Studio</span>
                  </div>
                  <span className="text-[10px] text-green-500 font-bold">Enabled</span>
                </Link>
              ) : (
                <button
                  type="button"
                  onClick={() => {
                    setIsOpen(false);
                    setActivationCap('seller');
                  }}
                  className="w-full flex items-center justify-between px-3 py-2 rounded-xl text-xs font-semibold text-orange-600 dark:text-orange-400 hover:bg-orange-50 dark:hover:bg-orange-950/20 transition-colors text-left"
                >
                  <div className="flex items-center gap-2.5">
                    <Sparkles size={15} />
                    <span>Activate Seller Studio</span>
                  </div>
                  <span className="text-[9px] uppercase tracking-wider font-bold bg-orange-600 text-white px-1.5 py-0.5 rounded">
                    New
                  </span>
                </button>
              )}

              {/* Landlord Console */}
              {isLandlord ? (
                <Link
                  href="/seller/products?preference=rent"
                  onClick={() => setIsOpen(false)}
                  className="flex items-center justify-between px-3 py-2 rounded-xl text-xs font-semibold text-zinc-700 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
                >
                  <div className="flex items-center gap-2.5">
                    <Home size={15} className="text-teal-500" />
                    <span>Landlord Console</span>
                  </div>
                  <span className="text-[10px] text-green-500 font-bold">Enabled</span>
                </Link>
              ) : (
                <button
                  type="button"
                  onClick={() => {
                    setIsOpen(false);
                    setActivationCap('landlord');
                  }}
                  className="w-full flex items-center justify-between px-3 py-2 rounded-xl text-xs font-semibold text-teal-600 dark:text-teal-400 hover:bg-teal-50 dark:hover:bg-teal-950/20 transition-colors text-left"
                >
                  <div className="flex items-center gap-2.5">
                    <Home size={15} />
                    <span>Activate Landlord Console</span>
                  </div>
                  <span className="text-[9px] uppercase tracking-wider font-bold bg-teal-600 text-white px-1.5 py-0.5 rounded">
                    Rent
                  </span>
                </button>
              )}

              {/* Admin Panel (Elevated permissions only) */}
              {isAdmin && (
                <Link
                  href="/admin/dashboard"
                  onClick={() => setIsOpen(false)}
                  className="flex items-center justify-between px-3 py-2 rounded-xl text-xs font-semibold text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-950/20 transition-colors"
                >
                  <div className="flex items-center gap-2.5">
                    <Shield size={15} />
                    <span>Admin Control Center</span>
                  </div>
                  <span className="text-[9px] uppercase font-bold bg-red-600 text-white px-1.5 py-0.5 rounded">
                    Admin
                  </span>
                </Link>
              )}
            </div>

            {/* Quick Actions */}
            <div className="pt-2 border-t border-zinc-100 dark:border-zinc-800 space-y-1">
              <Link
                href="/seller/product/new"
                onClick={() => setIsOpen(false)}
                className="flex items-center gap-2.5 px-3 py-2 rounded-xl text-xs font-bold text-orange-600 dark:text-orange-400 hover:bg-orange-50 dark:hover:bg-orange-950/20 transition-colors"
              >
                <PlusCircle size={15} />
                <span>Create New Listing</span>
              </Link>
              <button
                type="button"
                onClick={handleSignOut}
                className="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-xs font-semibold text-zinc-500 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-950/20 transition-colors text-left"
              >
                <LogOut size={15} />
                <span>Sign Out</span>
              </button>
            </div>
          </div>
        )}
      </div>

      {/* 1-Click Activation Modal */}
      {activationCap && (
        <CapabilityActivationModal
          isOpen={!!activationCap}
          capability={activationCap}
          onClose={() => setActivationCap(null)}
          onSuccess={() => {
            if (activationCap === 'seller') {
              router.push('/seller/dashboard');
            } else if (activationCap === 'landlord') {
              router.push('/seller/products?preference=rent');
            }
          }}
        />
      )}
    </>
  );
}
