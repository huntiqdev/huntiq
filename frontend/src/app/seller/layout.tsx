'use client';
import { useEffect, useState, useRef } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { SellerHelp } from '@/components/seller/SellerHelp';
import { LanguageSelector } from '@/components/layout/LanguageSelector';
import { useLocale } from '@/hooks/useLocale';
import {
  LayoutDashboard, Package, DollarSign, User,
  LogOut, Menu, ChevronRight, ShoppingBag, HelpCircle,
  ShoppingCart, BarChart2, Bell, Home, Layers, PhoneCall, Plus,
} from 'lucide-react';
import { TIER_RANK, TIER_ICON, TIER_COLOR } from '@/lib/constants/tiers';

// Each nav item has a minTier — only shown when seller's tier >= minTier
const NAV_ITEMS = [
  { href: '/seller/dashboard',    labelKey: 'nav_dashboard',    icon: LayoutDashboard, minTier: 0 },
  { href: '/seller/products',     labelKey: 'nav_my_products',  icon: Package,         minTier: 0 },
  { href: '/seller/orders',       labelKey: 'nav_orders',       icon: PhoneCall,       minTier: 1 }, // gold+
  { href: '/seller/transactions', labelKey: 'nav_sales',        icon: DollarSign,      minTier: 1 }, // gold+
  { href: '/seller/analytics',    labelKey: 'nav_analytics',    icon: BarChart2,       minTier: 2 }, // premium only
  { href: '/seller/profile',      labelKey: 'nav_profile',      icon: User,            minTier: 0 },
] as const;

const NAV_LABELS: Record<string, string> = {
  nav_dashboard:   'Dashboard',
  nav_my_products: 'My Products',
  nav_orders:      'Inquiries & Leads',
  nav_sales:       'Sales',
  nav_analytics:   'Analytics',
  nav_profile:     'Profile',
};

// ── Sidebar content — defined outside render to prevent remount ───────────────
interface SidebarProps {
  sellerName: string | null;
  sellerTier: string;
  userId?: number;
  avatarURL?: string;
  pathname: string;
  visibleNav: typeof NAV_ITEMS[number][];
  lockedNav: typeof NAV_ITEMS[number][];
  sellerPreference?: string;
  isLandlord?: boolean;
  onClose: () => void;
  onHelp: () => void;
  onLogout: () => void;
}

function SidebarContent({
  sellerName, sellerTier, userId, avatarURL, pathname,
  visibleNav, lockedNav, sellerPreference, isLandlord, onClose, onHelp, onLogout,
}: SidebarProps) {
  const { t } = useLocale();
  const sn = t.seller_nav;
  const TierIcon = TIER_ICON[sellerTier];

  const isRentMode = sellerPreference === 'rent' || (!sellerPreference && isLandlord);
  const isDualMode = sellerPreference === 'both';

  const panelTitle = 'HUNTIQ';

  const panelSubtitle = isRentMode
    ? (sn?.landlord_console ?? 'Landlord Console')
    : isDualMode
    ? (sn?.provider_studio ?? 'Provider Studio')
    : (sn?.seller_studio ?? 'Seller Studio');

  const activeColor = isRentMode
    ? 'bg-teal-600 text-white shadow-md shadow-teal-600/20'
    : isDualMode
    ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/20'
    : 'bg-[#1769E0] text-white shadow-md shadow-[#1769E0]/25';

  const logoColor = isRentMode
    ? 'bg-teal-600 shadow-teal-600/20'
    : isDualMode
    ? 'bg-indigo-600 shadow-indigo-600/20'
    : 'bg-[#1769E0] shadow-[#1769E0]/25';

  const getNavLabel = (key: string) => {
    if (key === 'nav_my_products') {
      return isRentMode ? (sn?.properties_rentals ?? 'Properties & Rentals') : isDualMode ? (sn?.my_products ?? 'Listings & Products') : (sn?.my_products ?? 'My Products');
    }
    if (key === 'nav_orders') {
      return isRentMode ? (sn?.lease_requests ?? 'Lease Requests') : (sn?.orders ?? 'Order Requests');
    }
    if (key === 'nav_sales') {
      return isRentMode ? (sn?.rental_income ?? 'Rental Income') : (sn?.sales ?? 'Sales');
    }
    if (key === 'nav_dashboard') return sn?.dashboard ?? 'Dashboard';
    if (key === 'nav_analytics') return sn?.analytics ?? 'Analytics';
    if (key === 'nav_profile') return sn?.profile ?? 'Profile';
    return NAV_LABELS[key] || key;
  };

  return (
    <div className="flex flex-col h-full">
      {/* Logo */}
      <div className="px-4 py-5 border-b border-zinc-100 dark:border-zinc-800">
        <div className="flex items-center gap-3">
          <div className={`w-8 h-8 ${logoColor} rounded-lg flex items-center justify-center shadow-md text-white`}>
            {isRentMode ? <Home size={16} /> : isDualMode ? <Layers size={16} /> : <ShoppingBag size={16} />}
          </div>
          <div className="flex-1 min-w-0">
            <p className="font-black text-sm text-zinc-900 dark:text-white leading-none">{panelTitle}</p>
            <p className="text-[10px] text-zinc-400 mt-0.5 uppercase tracking-widest">{panelSubtitle}</p>
          </div>
          <button onClick={onHelp} title="Help"
            className="w-7 h-7 rounded-lg flex items-center justify-center text-zinc-400 hover:text-[#1769E0] hover:bg-[#EAF3FF] dark:hover:bg-[#1769E0]/20 transition-colors shrink-0">
            <HelpCircle size={16} />
          </button>
        </div>
      </div>

      {/* Seller info + tier badge */}
      <div className="px-4 py-3 border-b border-zinc-100 dark:border-zinc-800">
        <div className="flex items-center gap-2.5">
          {avatarURL ? (
            <img src={avatarURL} alt={sellerName || 'Avatar'} className="w-8 h-8 rounded-lg object-cover ring-1 ring-zinc-300 dark:ring-zinc-700 shrink-0" />
          ) : (
            <div className={`w-8 h-8 rounded-lg ${isRentMode ? 'bg-teal-100 dark:bg-teal-900/30 text-teal-600' : isDualMode ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600' : 'bg-[#EAF3FF] dark:bg-[#1769E0]/20 text-[#1769E0]'} flex items-center justify-center font-black text-sm shrink-0`}>
              {sellerName?.[0]?.toUpperCase() ?? (isRentMode ? 'L' : 'S')}
            </div>
          )}
          <div className="min-w-0 flex-1">
            <p className="text-xs font-bold text-zinc-700 dark:text-zinc-300 truncate">
              {sellerName ?? (isRentMode ? `Landlord #${userId}` : isDualMode ? `Provider #${userId}` : `Seller #${userId}`)}
            </p>
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
                  ? activeColor
                  : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white hover:bg-zinc-100 dark:hover:bg-zinc-800'
              }`}>
              <Icon size={16} className={active ? 'text-white' : 'text-zinc-400 group-hover:text-zinc-600 dark:group-hover:text-zinc-300'} />
              <span className="flex-1">{getNavLabel(labelKey)}</span>
              {active && <ChevronRight size={12} className="text-white/60" />}
            </Link>
          );
        })}

        {/* Locked items */}
        {lockedNav.length > 0 && (
          <div className="pt-2 mt-1 border-t border-zinc-100 dark:border-zinc-800">
            <p className="text-[9px] font-black uppercase tracking-widest text-zinc-400 px-3 mb-1">
              {sn?.upgrade_hint ?? 'Upgrade to unlock'}
            </p>
            {lockedNav.map(({ labelKey, icon: Icon, minTier }) => (
              <Link key={labelKey} href="/seller/dashboard" onClick={onClose}
                className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-semibold text-zinc-300 dark:text-zinc-600 hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-all">
                <Icon size={16} className="text-zinc-300 dark:text-zinc-600" />
                <span className="flex-1 line-through opacity-50">{getNavLabel(labelKey)}</span>
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
      <div className="px-3 py-3 border-t border-zinc-100 dark:border-zinc-800 space-y-2">
        <div className="px-1">
          <LanguageSelector variant="full" dropDirection="up" />
        </div>
        <Link href="/" className="flex items-center gap-3 px-3 py-2 rounded-xl text-xs font-semibold text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-all">
          ← {sn?.back_to_marketplace ?? 'Back to Marketplace'}
        </Link>
        <button onClick={onLogout}
          className="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-semibold text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-all">
          <LogOut size={15} /> {sn?.sign_out ?? 'Sign Out'}
        </button>
      </div>
    </div>
  );
}

export default function SellerLayout({ children }: { children: React.ReactNode }) {
  const { t } = useLocale();
  const sn = t.seller_nav;
  const { user, isSeller, isLandlord, sellerPreference, loading, logout, activateCapability } = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [helpOpen, setHelpOpen] = useState(false);
  const [sellerName, setSellerName] = useState<string | null>(null);
  const [sellerTier, setSellerTier] = useState<string>('normal');
  const [unreadCount, setUnreadCount] = useState(0);

  const isProvider = isSeller || isLandlord;
  const sellerFetched = useRef(false);

  const isAuthPage = pathname === '/seller/login' || pathname === '/seller/register';

  useEffect(() => {
    if (!loading && !isAuthPage && !user) {
      const isRent = pathname.includes('rent') || (typeof window !== 'undefined' && window.location.search.includes('type=rent'));
      const intent = isRent ? 'rent' : 'sell';
      const redirect = typeof window !== 'undefined' ? `${pathname}${window.location.search}` : pathname;
      router.replace(`/auth?intent=${intent}&redirect=${encodeURIComponent(redirect)}`);
    }
  }, [loading, user, isAuthPage, pathname, router]);

  useEffect(() => {
    if (!user || !isProvider) return;
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
  }, [user?.token, isProvider, pathname]);

  if (isAuthPage) return <>{children}</>;
  if (loading || !user) {
    return (
      <div className="flex h-screen bg-inherit dark:bg-zinc-950 overflow-hidden items-center justify-center">
        <div className="w-8 h-8 border-2 border-[#1769E0] border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  // If user is logged in but hasn't activated seller or landlord capability yet
  if (!isProvider) {
    if (pathname.startsWith('/seller/product/new')) {
      return <>{children}</>;
    }
    return (
      <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 flex items-center justify-center p-4">
        <div className="max-w-lg w-full bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 shadow-2xl text-center space-y-6">
          <div className="flex items-center justify-between">
            <div className="w-12 h-12 bg-gradient-to-tr from-[#1769E0] to-teal-500 rounded-2xl flex items-center justify-center text-white shadow-lg shadow-[#1769E0]/20">
              <Layers size={24} />
            </div>
            <LanguageSelector variant="pill" />
          </div>
          <div>
            <h2 className="text-xl font-black text-zinc-900 dark:text-white">
              {sn?.choose_console ?? 'Choose Your Provider Console'}
            </h2>
            <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1 max-w-sm mx-auto">
              HUNTIQ separates Goods Merchants from Rental Hosts. Activate the capability matching your operational model:
            </p>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-left">
            {/* Goods Seller */}
            <div className="border border-[#1769E0]/30 bg-[#1769E0]/5 dark:bg-[#1769E0]/20 rounded-2xl p-4 flex flex-col justify-between">
              <div>
                <div className="w-8 h-8 rounded-lg bg-[#1769E0] text-white flex items-center justify-center mb-2">
                  <ShoppingBag size={16} />
                </div>
                <h3 className="font-bold text-sm text-zinc-900 dark:text-white">
                  {sn?.seller_studio ?? 'Seller Studio'}
                </h3>
                <p className="text-[11px] text-zinc-500 dark:text-zinc-400 mt-1">
                  For merchants selling physical goods, electronics, vehicles, and catalog inventory.
                </p>
              </div>
              <button
                onClick={async () => {
                  await activateCapability('seller', undefined, 'sale');
                }}
                className="mt-4 w-full py-2.5 bg-[#1769E0] hover:bg-[#1055BD] text-white rounded-xl font-bold text-xs shadow-md shadow-[#1769E0]/20 transition-all cursor-pointer"
              >
                {sn?.activate ?? 'Activate'} {sn?.seller_studio ?? 'Seller Studio'}
              </button>
            </div>

            {/* Rental Landlord */}
            <div className="border border-teal-200 dark:border-teal-900/40 bg-teal-50/50 dark:bg-teal-950/20 rounded-2xl p-4 flex flex-col justify-between">
              <div>
                <div className="w-8 h-8 rounded-lg bg-teal-600 text-white flex items-center justify-center mb-2">
                  <Home size={16} />
                </div>
                <h3 className="font-bold text-sm text-zinc-900 dark:text-white">
                  {sn?.landlord_console ?? 'Landlord Console'}
                </h3>
                <p className="text-[11px] text-zinc-500 dark:text-zinc-400 mt-1">
                  For property landlords, vehicle rental agencies, equipment leasing, and booking hosts.
                </p>
              </div>
              <button
                onClick={async () => {
                  await activateCapability('landlord', undefined, 'rent');
                }}
                className="mt-4 w-full py-2.5 bg-teal-600 hover:bg-teal-500 text-white rounded-xl font-bold text-xs shadow-md shadow-teal-600/20 transition-all cursor-pointer"
              >
                {sn?.activate ?? 'Activate'} {sn?.landlord_console ?? 'Landlord Console'}
              </button>
            </div>
          </div>

          <div className="pt-2 flex flex-col items-center gap-2">
            <button
              onClick={async () => {
                await activateCapability('seller', undefined, 'both');
              }}
              className="text-xs font-semibold text-indigo-600 dark:text-indigo-400 hover:underline cursor-pointer"
            >
              Or activate Dual Provider (Sell Goods & Rent Properties) →
            </button>
            <Link
              href="/"
              className="text-xs text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 font-semibold pt-1"
            >
              Back to Marketplace
            </Link>
          </div>
        </div>
      </div>
    );
  }

  const isRentMode = sellerPreference === 'rent' || (!sellerPreference && isLandlord);
  const isDualMode = sellerPreference === 'both';

  const tierRank = TIER_RANK[sellerTier] ?? 0;

  const visibleNav = NAV_ITEMS.filter(item => tierRank >= item.minTier);
  const lockedNav  = NAV_ITEMS.filter(item => tierRank < item.minTier);

  const sidebarProps: SidebarProps = {
    sellerName, sellerTier, userId: user?.id, avatarURL: user?.avatarURL, pathname,
    visibleNav, lockedNav, sellerPreference, isLandlord,
    onClose:  () => setSidebarOpen(false),
    onHelp:   () => setHelpOpen(true),
    onLogout: () => { logout(); router.push('/auth'); },
  };

  const headerButtonText = isRentMode
    ? '+ Add Property'
    : isDualMode
    ? '+ Add Listing'
    : '+ Add Product';

  const headerAddHref = isRentMode
    ? '/seller/products?preference=rent&new=1'
    : isDualMode
    ? '/seller/products?new=1'
    : '/seller/products?preference=sale&new=1';

  const headerColorClass = isRentMode
    ? 'bg-teal-600 hover:bg-teal-500'
    : isDualMode
    ? 'bg-indigo-600 hover:bg-indigo-500'
    : 'bg-[#1769E0] hover:bg-[#1055BD]';

  const headerTitle = isRentMode
    ? 'Landlord Console'
    : isDualMode
    ? 'Provider Studio'
    : 'Seller Panel';

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
        <header className="md:hidden h-14 bg-white dark:bg-zinc-900 border-b border-zinc-200 dark:border-zinc-800 flex items-center px-3 gap-2 shrink-0">
          <button onClick={() => setSidebarOpen(true)}
            className="p-2 rounded-xl text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors tap-bounce active:scale-90 shrink-0"
            aria-label="Open menu">
            <Menu size={18} />
          </button>
          <span className={`font-black text-xs sm:text-sm truncate min-w-0 ${isRentMode ? 'text-teal-600' : isDualMode ? 'text-indigo-600' : 'text-[#1769E0]'}`}>
            {headerTitle}
          </span>
          <div className="ml-auto flex items-center gap-1.5 shrink-0">
            <LanguageSelector variant="compact" className="shrink-0" />
            {unreadCount > 0 && (
              <Link href="/seller/dashboard" className="relative p-1.5 rounded-xl text-[#1769E0] hover:bg-[#1769E0]/10 dark:hover:bg-[#1769E0]/20 transition-colors tap-bounce active:scale-90 shrink-0" aria-label="Notifications">
                <Bell size={16} />
                <span className="absolute -top-0.5 -right-0.5 w-3.5 h-3.5 bg-[#1769E0] text-white text-[8px] font-black rounded-full flex items-center justify-center">
                  {unreadCount > 9 ? '9+' : unreadCount}
                </span>
              </Link>
            )}
            <button onClick={() => setHelpOpen(true)}
              className="p-1.5 rounded-xl text-zinc-400 hover:text-[#1769E0] hover:bg-[#1769E0]/10 dark:hover:bg-[#1769E0]/20 transition-colors tap-bounce active:scale-90 shrink-0"
              aria-label="Help">
              <HelpCircle size={16} />
            </button>
            <Link href={headerAddHref}
              className={`flex items-center gap-1 text-white px-2.5 py-1.5 min-h-[34px] rounded-xl text-xs font-bold shadow-sm tap-bounce active:scale-95 shrink-0 ${headerColorClass}`}>
              <Plus size={14} className="stroke-[2.5px]" />
              <span className="hidden xs:inline">{headerButtonText.replace('+ ', '')}</span>
            </Link>
          </div>
        </header>

        <header className="hidden md:flex h-14 bg-white dark:bg-zinc-900 border-b border-zinc-200 dark:border-zinc-800 items-center px-6 gap-3 shrink-0">
          <div className="flex-1">
            <p className="text-sm font-bold text-zinc-500 dark:text-zinc-400 capitalize">
              {pathname.split('/').pop()?.replace(/-/g, ' ') || 'Dashboard'}
            </p>
          </div>
          <LanguageSelector variant="pill" />
          {unreadCount > 0 && (
            <Link href="/seller/dashboard" className="relative p-1.5 rounded-lg text-[#1769E0] hover:bg-[#1769E0]/10 dark:hover:bg-[#1769E0]/20 transition-colors" title={`${unreadCount} unread notification${unreadCount > 1 ? 's' : ''}`}>
              <Bell size={16} />
              <span className="absolute -top-1 -right-1 w-4 h-4 bg-[#1769E0] text-white text-[9px] font-black rounded-full flex items-center justify-center">
                {unreadCount > 9 ? '9+' : unreadCount}
              </span>
            </Link>
          )}
          <Link href={headerAddHref}
            className={`flex items-center gap-1.5 text-white px-3 py-1.5 rounded-lg text-xs font-bold transition-all shadow-sm ${headerColorClass}`}>
            {headerButtonText}
          </Link>
          <Link href="/" className="text-xs text-zinc-400 hover:text-zinc-600 transition-colors">
            View Marketplace ↗
          </Link>
        </header>

        <main className="flex-1 overflow-y-auto p-4 md:p-6 pb-[max(24px,env(safe-area-inset-bottom,24px))]">{children}</main>
      </div>
    </div>
  );
}
