'use client';
import { useEffect, useState, useRef } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { SellerHelp } from '@/components/seller/SellerHelp';
import {
  LayoutDashboard, Package, DollarSign, User,
  LogOut, Menu, ChevronRight, ShoppingBag, HelpCircle,
  ShoppingCart, BarChart2, Bell,
} from 'lucide-react';
import { TIER_RANK, TIER_ICON, TIER_COLOR } from '@/lib/constants/tiers';

// Each nav item has a minTier — only shown when seller's tier >= minTier
const NAV_ITEMS = [
  { href: '/seller/dashboard',    labelKey: 'nav_dashboard',    icon: LayoutDashboard, minTier: 0 },
  { href: '/seller/products',     labelKey: 'nav_my_products',  icon: Package,         minTier: 0 },
  { href: '/seller/orders',       labelKey: 'nav_orders',       icon: ShoppingCart,    minTier: 1 }, // gold+
  { href: '/seller/transactions', labelKey: 'nav_sales',        icon: DollarSign,      minTier: 1 }, // gold+
  { href: '/seller/analytics',    labelKey: 'nav_analytics',    icon: BarChart2,       minTier: 2 }, // premium only
  { href: '/seller/profile',      labelKey: 'nav_profile',      icon: User,            minTier: 0 },
] as const;

const NAV_LABELS: Record<string, string> = {
  nav_dashboard:   'Dashboard',
  nav_my_products: 'My Products',
  nav_orders:      'Order Requests',
  nav_sales:       'Sales',
  nav_analytics:   'Analytics',
  nav_profile:     'Profile',
};

// ── Sidebar content — defined outside render to prevent remount ───────────────
interface SidebarProps {
  sellerName: string | null;
  sellerTier: string;
  userId?: number;
  pathname: string;
  visibleNav: typeof NAV_ITEMS[number][];
  lockedNav: typeof NAV_ITEMS[number][];
  onClose: () => void;
  onHelp: () => void;
  onLogout: () => void;
}

function SidebarContent({
  sellerName, sellerTier, userId, pathname,
  visibleNav, lockedNav, onClose, onHelp, onLogout,
}: SidebarProps) {
  const TierIcon = TIER_ICON[sellerTier];
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
            <p className="text-[10px] text-zinc-400 mt-0.5 uppercase tracking-widest">Seller Panel</p>
          </div>
          <button onClick={onHelp} title="Help"
            className="w-7 h-7 rounded-lg flex items-center justify-center text-zinc-400 hover:text-orange-600 hover:bg-orange-50 dark:hover:bg-orange-900/20 transition-colors shrink-0">
            <HelpCircle size={16} />
          </button>
        </div>
      </div>

      {/* Seller info + tier badge */}
      <div className="px-4 py-3 border-b border-zinc-100 dark:border-zinc-800">
        <div className="flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-lg bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center text-orange-600 font-black text-sm">
            {sellerName?.[0]?.toUpperCase() ?? 'S'}
          </div>
          <div className="min-w-0 flex-1">
            <p className="text-xs font-bold text-zinc-700 dark:text-zinc-300 truncate">{sellerName ?? `Seller #${userId}`}</p>
            <div className="flex items-center gap-1 mt-0.5">
              {TierIcon && <TierIcon size={9} className={TIER_COLOR[sellerTier]} />}
              <p className={`text-[10px] font-bold capitalize ${TIER_COLOR[sellerTier]}`}>{sellerTier} tier</p>
            </div>
          </div>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-3 space-y-0.5 overflow-y-auto">
        {visibleNav.map(({ href, labelKey, icon: Icon }) => {
          const active = pathname === href;
          return (
            <Link key={href} href={href} onClick={onClose}
              className={`flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-semibold transition-all group ${
                active
                  ? 'bg-orange-600 text-white shadow-md shadow-orange-600/20'
                  : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white hover:bg-zinc-100 dark:hover:bg-zinc-800'
              }`}>
              <Icon size={16} className={active ? 'text-white' : 'text-zinc-400 group-hover:text-zinc-600 dark:group-hover:text-zinc-300'} />
              <span className="flex-1">{NAV_LABELS[labelKey]}</span>
              {active && <ChevronRight size={12} className="text-white/60" />}
            </Link>
          );
        })}

        {/* Locked items */}
        {lockedNav.length > 0 && (
          <div className="pt-2 mt-1 border-t border-zinc-100 dark:border-zinc-800">
            <p className="text-[9px] font-black uppercase tracking-widest text-zinc-400 px-3 mb-1">Upgrade to unlock</p>
            {lockedNav.map(({ labelKey, icon: Icon, minTier }) => (
              <Link key={labelKey} href="/seller/dashboard" onClick={onClose}
                className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-semibold text-zinc-300 dark:text-zinc-600 hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-all">
                <Icon size={16} className="text-zinc-300 dark:text-zinc-600" />
                <span className="flex-1 line-through opacity-50">{NAV_LABELS[labelKey]}</span>
                <span className={`text-[8px] font-black px-1.5 py-0.5 rounded-full shrink-0 ${
                  minTier === 2
                    ? 'bg-purple-500/20 text-purple-400'
                    : 'bg-yellow-500/20 text-yellow-500'
                }`}>
                  {minTier === 2 ? 'Premium' : 'Gold'}
                </span>
              </Link>
            ))}
          </div>
        )}
      </nav>

      {/* Footer */}
      <div className="px-3 py-3 border-t border-zinc-100 dark:border-zinc-800 space-y-1">
        <Link href="/" className="flex items-center gap-3 px-3 py-2 rounded-xl text-xs font-semibold text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-all">
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

export default function SellerLayout({ children }: { children: React.ReactNode }) {
  const { user, isSeller, isLandlord, loading, logout, activateCapability } = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [helpOpen, setHelpOpen] = useState(false);
  const [sellerName, setSellerName] = useState<string | null>(null);
  const [sellerTier, setSellerTier] = useState<string>('normal');
  const [unreadCount, setUnreadCount] = useState(0);

  const sellerFetched = useRef(false);

  const isAuthPage = pathname === '/seller/login' || pathname === '/seller/register';

  useEffect(() => {
    if (!loading && !isAuthPage && !user) router.replace('/auth?intent=sell');
  }, [loading, user, isAuthPage, router]);

  useEffect(() => {
    if (!user || !isSeller) return;
    const isDashboard = pathname === '/seller/dashboard';
    if (!sellerFetched.current || isDashboard) {
      sellerFetched.current = true;
      Promise.all([
        api.getSellerStats(user.token).catch(() => null),
        api.getSellerProfile(user.token).catch(() => null),
        api.sellerGetNotifications(user.token).catch(() => []),
      ]).then(([s, p, notifs]) => {
        if (p?.name) setSellerName(p.name);
        const tier = s?.tier || p?.tier || 'normal';
        setSellerTier(tier);
        const unread = (notifs as any[]).filter((n: any) => !n.is_read).length;
        setUnreadCount(unread);
      });
    }
  }, [user?.token, isSeller, pathname]);

  if (isAuthPage) return <>{children}</>;
  if (loading) return null;
  if (!user) return null;

  // If user is logged in but hasn't activated seller capability yet
  if (!isSeller) {
    return (
      <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 flex items-center justify-center p-4">
        <div className="max-w-md w-full bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 shadow-2xl text-center space-y-4">
          <div className="w-14 h-14 bg-orange-600 rounded-2xl flex items-center justify-center mx-auto text-white shadow-lg shadow-orange-600/30">
            <ShoppingBag size={28} />
          </div>
          <div>
            <h2 className="text-xl font-black text-zinc-900 dark:text-white">Activate Seller Studio</h2>
            <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">
              You are signed in to your SellBuy account. Activate the seller capability to start listing products, managing inventory, and receiving orders.
            </p>
          </div>
          <div className="pt-2 flex flex-col gap-2">
            <button
              onClick={async () => {
                await activateCapability('seller');
              }}
              className="w-full py-3.5 bg-orange-600 hover:bg-orange-500 text-white rounded-xl font-black text-sm shadow-lg shadow-orange-600/25 transition-all active:scale-95"
            >
              Activate & Enter Seller Studio
            </button>
            <Link
              href="/"
              className="text-xs text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 font-semibold py-2"
            >
              Back to Marketplace
            </Link>
          </div>
        </div>
      </div>
    );
  }

  const tierRank = TIER_RANK[sellerTier] ?? 0;

  const visibleNav = NAV_ITEMS.filter(item => tierRank >= item.minTier);
  const lockedNav  = NAV_ITEMS.filter(item => tierRank < item.minTier);

  const sidebarProps: SidebarProps = {
    sellerName, sellerTier, userId: user?.id, pathname,
    visibleNav, lockedNav,
    onClose:  () => setSidebarOpen(false),
    onHelp:   () => setHelpOpen(true),
    onLogout: () => { logout(); router.push('/auth'); },
  };

  return (
    <div className="flex h-screen bg-inherit dark:bg-zinc-950 overflow-hidden">
      {helpOpen && <SellerHelp onClose={() => setHelpOpen(false)} />}

      {/* Desktop sidebar */}
      <aside className="hidden md:flex w-56 bg-white dark:bg-zinc-900 border-r border-zinc-200 dark:border-zinc-800 flex-col shrink-0">
        <SidebarContent {...sidebarProps} />
      </aside>

      {/* Mobile overlay */}
      {sidebarOpen && (
        <div className="fixed inset-0 z-50 md:hidden">
          <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={() => setSidebarOpen(false)} />
          <aside className="absolute left-0 top-0 bottom-0 w-60 bg-white dark:bg-zinc-900 border-r border-zinc-200 dark:border-zinc-800 flex flex-col">
            <SidebarContent {...sidebarProps} onClose={() => setSidebarOpen(false)} />
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
          <span className="font-black text-orange-600 text-sm">Seller Panel</span>
          <div className="ml-auto flex items-center gap-2">
            {unreadCount > 0 && (
              <Link href="/seller/dashboard" className="relative p-1.5 rounded-lg text-orange-500 hover:bg-orange-50 dark:hover:bg-orange-900/20 transition-colors">
                <Bell size={16} />
                <span className="absolute -top-0.5 -right-0.5 w-4 h-4 bg-orange-500 text-white text-[9px] font-black rounded-full flex items-center justify-center">
                  {unreadCount > 9 ? '9+' : unreadCount}
                </span>
              </Link>
            )}
            <button onClick={() => setHelpOpen(true)}
              className="p-1.5 rounded-lg text-zinc-400 hover:text-orange-600 hover:bg-orange-50 dark:hover:bg-orange-900/20 transition-colors">
              <HelpCircle size={16} />
            </button>
            <Link href="/seller/dashboard"
              className="flex items-center gap-1.5 bg-orange-600 text-white px-3 py-1.5 rounded-lg text-xs font-bold">
              + Add
            </Link>
          </div>
        </header>

        <header className="hidden md:flex h-14 bg-white dark:bg-zinc-900 border-b border-zinc-200 dark:border-zinc-800 items-center px-6 gap-3 shrink-0">
          <div className="flex-1">
            <p className="text-sm font-bold text-zinc-500 dark:text-zinc-400 capitalize">
              {pathname.split('/').pop()?.replace(/-/g, ' ') || 'Dashboard'}
            </p>
          </div>
          {unreadCount > 0 && (
            <Link href="/seller/dashboard" className="relative p-1.5 rounded-lg text-orange-500 hover:bg-orange-50 dark:hover:bg-orange-900/20 transition-colors" title={`${unreadCount} unread notification${unreadCount > 1 ? 's' : ''}`}>
              <Bell size={16} />
              <span className="absolute -top-1 -right-1 w-4 h-4 bg-orange-500 text-white text-[9px] font-black rounded-full flex items-center justify-center">
                {unreadCount > 9 ? '9+' : unreadCount}
              </span>
            </Link>
          )}
          <Link href="/" className="text-xs text-zinc-400 hover:text-zinc-600 transition-colors">
            View Marketplace ↗
          </Link>
        </header>

        <main className="flex-1 overflow-y-auto p-4 md:p-6">{children}</main>
      </div>
    </div>
  );
}
