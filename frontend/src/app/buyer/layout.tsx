'use client';
import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import {
  LayoutDashboard, User, LogOut, Menu, ChevronRight,
  ShoppingBag, Lock, ShoppingCart, Gift,
} from 'lucide-react';
import { TIER_RANK, TIER_ICON, TIER_COLOR } from '@/lib/constants/tiers';

const NAV_ITEMS = [
  { href: '/buyer/dashboard', label: 'Dashboard',      icon: LayoutDashboard, minTier: 0 },
  { href: '/buyer/orders',    label: 'My Requests',    icon: ShoppingCart,    minTier: 0 },
  { href: '/buyer/referral',  label: 'Referral',       icon: Gift,            minTier: 0 },
  { href: '/buyer/profile',   label: 'Profile',        icon: User,            minTier: 0 },
] as const;

interface BuyerSidebarProps {
  buyerName: string | null;
  buyerTier: string;
  userId?: number;
  pathname: string;
  tierRank: number;
  onClose: () => void;
  onLogout: () => void;
}

// Defined outside BuyerLayout to prevent remount on every render (avoids hydration issues)
function BuyerSidebarContent({ buyerName, buyerTier, userId, pathname, tierRank, onClose, onLogout }: BuyerSidebarProps) {
  const TierIcon = TIER_ICON[buyerTier];
  return (
    <div className="flex flex-col h-full">
      {/* Logo */}
      <div className="px-4 py-5 border-b border-zinc-100 dark:border-zinc-800">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-orange-600 rounded-lg flex items-center justify-center shadow-md">
            <ShoppingBag size={16} className="text-white" />
          </div>
          <div className="flex-1 min-w-0">
            <p className="font-black text-sm text-zinc-900 dark:text-white leading-none">HUNTIQ</p>
            <p className="text-[10px] text-zinc-400 mt-0.5 uppercase tracking-widest">Buyer Panel</p>
          </div>
        </div>
      </div>

      {/* Buyer info + tier badge */}
      <div className="px-4 py-3 border-b border-zinc-100 dark:border-zinc-800">
        <div className="flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-lg bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center text-orange-600 font-black text-sm">
            {buyerName?.[0]?.toUpperCase() ?? 'B'}
          </div>
          <div className="min-w-0 flex-1">
            <p className="text-xs font-bold text-zinc-700 dark:text-zinc-300 truncate">{buyerName ?? `Buyer #${userId}`}</p>
            <div className="flex items-center gap-1 mt-0.5">
              {TierIcon && <TierIcon size={9} className={TIER_COLOR[buyerTier]} />}
              <p className={`text-[10px] font-bold capitalize ${TIER_COLOR[buyerTier]}`}>{buyerTier} buyer</p>
            </div>
          </div>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-3 space-y-0.5 overflow-y-auto">
        {NAV_ITEMS.filter(item => tierRank >= item.minTier).map(({ href, label, icon: Icon }) => {
          const active = pathname === href;
          return (
            <Link key={href} href={href} onClick={onClose}
              className={`flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-semibold transition-all group ${
                active
                  ? 'bg-orange-600 text-white shadow-md shadow-orange-600/20'
                  : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white hover:bg-zinc-100 dark:hover:bg-zinc-800'
              }`}>
              <Icon size={16} className={active ? 'text-white' : 'text-zinc-400 group-hover:text-zinc-600 dark:group-hover:text-zinc-300'} />
              <span className="flex-1">{label}</span>
              {active && <ChevronRight size={12} className="text-white/60" />}
            </Link>
          );
        })}

        {/* Upgrade CTA */}
        {buyerTier === 'normal' && (
          <div className="pt-2 mt-1 border-t border-zinc-100 dark:border-zinc-800">
            <p className="text-[9px] font-black uppercase tracking-widest text-zinc-400 px-3 mb-1">Upgrade to unlock</p>
            <Link href="/buyer/upgrade" onClick={onClose}
              className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-semibold text-yellow-600 dark:text-yellow-500 hover:bg-yellow-50 dark:hover:bg-yellow-900/20 transition-all">
              {(() => { const I = TIER_ICON['gold']; return <I size={16} className="text-yellow-500" />; })()}
              <span className="flex-1">Gold / Premium</span>
              <Lock size={10} className="text-yellow-500 shrink-0" />
            </Link>
          </div>
        )}
        {buyerTier === 'gold' && (
          <div className="pt-2 mt-1 border-t border-zinc-100 dark:border-zinc-800">
            <Link href="/buyer/upgrade" onClick={onClose}
              className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-semibold text-purple-600 dark:text-purple-400 hover:bg-purple-50 dark:hover:bg-purple-900/20 transition-all">
              {(() => { const I = TIER_ICON['premium']; return <I size={16} className="text-purple-500" />; })()}
              <span className="flex-1">Upgrade to Premium</span>
              <Lock size={10} className="text-purple-500 shrink-0" />
            </Link>
          </div>
        )}
      </nav>

      {/* Footer */}
      <div className="px-3 py-3 border-t border-zinc-100 dark:border-zinc-800 space-y-1">
        <Link href="/" onClick={onClose}
          className="flex items-center gap-3 px-3 py-2 rounded-xl text-xs font-semibold text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-all">
          ← Back to Marketplace
        </Link>
        <button onClick={onLogout}
          className="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-semibold text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-all">
          <LogOut size={15} /> Sign Out
        </button>
      </div>
    </div>
  );
}

const NOOP = () => {};

export default function BuyerLayout({ children }: { children: React.ReactNode }) {
  const { user, isBuyer, loading, logout, buyerTier: jwtTier } = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [buyerName, setBuyerName] = useState<string | null>(null);
  // Live tier from profile — may differ from JWT if admin changed it since last login
  const [liveTier, setLiveTier] = useState<string | null>(null);
  const buyerTier = liveTier ?? jwtTier;

  const isAuthPage = pathname === '/buyer/login' || pathname === '/buyer/register';

  useEffect(() => {
    if (!loading && !isAuthPage && !isBuyer) router.replace('/auth?as=buyer');
  }, [loading, isBuyer, isAuthPage, router]);

  useEffect(() => {
    if (!user || !isBuyer) return;
    api.getBuyerProfile(user.token)
      .then(p => {
        if (p?.name) setBuyerName(p.name);
        if (p?.buyer_tier) setLiveTier(p.buyer_tier);
      })
      .catch(() => {});
  }, [user?.token, isBuyer]);

  if (isAuthPage) return <>{children}</>;
  if (loading) return null;
  if (!isBuyer) return null;

  const tierRank = TIER_RANK[buyerTier] ?? 0;
  const TierIcon = TIER_ICON[buyerTier] as React.ElementType | undefined;

  const handleLogout = () => { logout(); router.push('/auth?as=buyer'); };

  const sidebarProps: BuyerSidebarProps = {
    buyerName, buyerTier, userId: user?.id, pathname, tierRank,
    onClose: () => setSidebarOpen(false),
    onLogout: handleLogout,
  };

  return (
    <div className="flex h-screen bg-inherit dark:bg-zinc-950 overflow-hidden">
      {/* Desktop sidebar */}
      <aside className="hidden md:flex w-56 bg-white dark:bg-zinc-900 border-r border-zinc-200 dark:border-zinc-800 flex-col shrink-0">
        <BuyerSidebarContent {...sidebarProps} onClose={NOOP} />
      </aside>

      {/* Mobile overlay */}
      {sidebarOpen && (
        <div className="fixed inset-0 z-50 md:hidden">
          <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={() => setSidebarOpen(false)} />
          <aside className="absolute left-0 top-0 bottom-0 w-60 bg-white dark:bg-zinc-900 border-r border-zinc-200 dark:border-zinc-800 flex flex-col">
            <BuyerSidebarContent {...sidebarProps} />
          </aside>
        </div>
      )}

      {/* Main */}
      <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
        <header className="md:hidden h-14 bg-white dark:bg-zinc-900 border-b border-zinc-200 dark:border-zinc-800 flex items-center px-4 gap-3 shrink-0">
          <button onClick={() => setSidebarOpen(true)}
            className="p-2 rounded-lg text-zinc-500 hover:text-zinc-700 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors">
            <Menu size={18} />
          </button>
          <span className="font-black text-orange-600 text-sm">Buyer Panel</span>
          <div className="ml-auto flex items-center gap-2">
            {TierIcon && <TierIcon size={14} className={TIER_COLOR[buyerTier]} />}
            <span className={`text-xs font-black capitalize ${TIER_COLOR[buyerTier]}`}>{buyerTier}</span>
          </div>
        </header>

        <header className="hidden md:flex h-14 bg-white dark:bg-zinc-900 border-b border-zinc-200 dark:border-zinc-800 items-center px-6 gap-3 shrink-0">
          <div className="flex-1">
            <p className="text-sm font-bold text-zinc-500 dark:text-zinc-400 capitalize">
              {pathname.split('/').pop()?.replace(/-/g, ' ') || 'Dashboard'}
            </p>
          </div>
          <Link href="/" className="text-xs text-zinc-400 hover:text-zinc-600 transition-colors">
            View Marketplace ↗
          </Link>
        </header>

        <main className="flex-1 overflow-y-auto p-4 md:p-6">{children}</main>
      </div>
    </div>
  );
}
