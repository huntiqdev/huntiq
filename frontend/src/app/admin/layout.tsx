'use client';
import { ReactNode, useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { useRouter, usePathname } from 'next/navigation';
import { Spinner } from '@/components/ui/Spinner';
import Link from 'next/link';
import {
  LayoutDashboard, Settings, Users, Package, LogOut, FileText, Activity,
  BadgeCheck, AlertTriangle, Scale, FolderTree, Tag, Search, MapPin,
  DollarSign, TrendingUp, Crown, ArrowUpCircle, Share2, Store, Receipt, X,
  ShieldAlert, ShieldCheck, Server
} from 'lucide-react';
import { AdminHeaderTopBar } from '@/components/admin/enterprise/AdminHeaderTopBar';
import { MultiStepListingWizard } from '@/components/admin/enterprise/MultiStepListingWizard';
import { HuntiqLogo } from '@/components/ui/HuntiqLogo';

interface AdminSidebarProps {
  onClose?: () => void;
  className?: string;
}

function AdminSidebar({ onClose, className = '' }: AdminSidebarProps) {
  const {
    isSuperAdmin, isPlatformOwner, isPlatformAdmin, isSecurityAdmin,
    isOperationsAdmin, isModerator, isSupportAgent,
    hasAdminPermission, logout
  } = useAuth();
  const pathname = usePathname();

  const roleLabel = isPlatformOwner ? 'Platform Owner'
    : isPlatformAdmin ? 'Platform Administrator'
    : isSecurityAdmin ? 'Security Administrator'
    : isOperationsAdmin ? 'Operations Administrator'
    : isModerator ? 'Content Moderator'
    : isSupportAgent ? 'Customer Support Agent'
    : isSuperAdmin ? 'Super Admin'
    : 'Enterprise Staff';

  const roleBadgeStyle = isPlatformOwner
    ? 'text-rose-400 border-rose-500/30 bg-rose-950/30'
    : isPlatformAdmin
    ? 'text-blue-400 border-blue-500/30 bg-blue-950/30'
    : isSecurityAdmin
    ? 'text-emerald-400 border-emerald-500/30 bg-emerald-950/30'
    : 'text-[#38A5FF] border-[#1769E0]/30 bg-[#1769E0]/15';

  const navigationGroups = [
    {
      title: 'ENTERPRISE CONSOLES',
      links: [
        { href: '/executive-console', label: 'Executive Cockpit', icon: ShieldAlert, visible: isPlatformOwner || isSuperAdmin },
        { href: '/security-console', label: 'Security SOC Defense', icon: ShieldCheck, visible: isSecurityAdmin || isPlatformOwner || isSuperAdmin },
        { href: '/platform-console', label: 'Platform Infrastructure', icon: Server, visible: isPlatformAdmin || isPlatformOwner || isSuperAdmin },
      ],
    },
    {
      title: 'OPERATIONS & CONTROL',
      links: [
        { href: '/admin/dashboard', label: 'Overview', icon: LayoutDashboard, visible: true },
        { href: '/admin/sellers', label: 'Sellers & Landlords', icon: Store, visible: isPlatformOwner || isPlatformAdmin || isOperationsAdmin || hasAdminPermission('sellers') },
        { href: '/admin/operations', label: 'Operations Monitor', icon: Activity, visible: isPlatformOwner || isPlatformAdmin || isSecurityAdmin || isOperationsAdmin || hasAdminPermission('operations') },
      ],
    },
    {
      title: 'LISTINGS & CATALOG',
      links: [
        { href: '/admin/products',  label: 'Products & Rentals', icon: Package,   visible: isPlatformOwner || isPlatformAdmin || isOperationsAdmin || isModerator || hasAdminPermission('pending') || hasAdminPermission('products') },
        { href: '/admin/classification', label: 'Universal Classification', icon: FolderTree, visible: isPlatformOwner || isPlatformAdmin || isOperationsAdmin || hasAdminPermission('categories') },
        { href: '/admin/categories',label: 'Legacy Categories', icon: FolderTree, visible: isPlatformOwner || isPlatformAdmin || isOperationsAdmin || hasAdminPermission('categories') },
        { href: '/admin/attributes', label: 'Attributes',         icon: Tag,        visible: isPlatformOwner || isPlatformAdmin || isOperationsAdmin || hasAdminPermission('categories') },
        { href: '/admin/locations', label: 'Locations & GPS',     icon: MapPin,     visible: isPlatformOwner || isPlatformAdmin || isOperationsAdmin || hasAdminPermission('locations') || hasAdminPermission('categories') },
        { href: '/admin/search',    label: 'Search Analytics',    icon: Search,     visible: isPlatformOwner || isPlatformAdmin || isOperationsAdmin || hasAdminPermission('analytics') },
      ],
    },
    {
      title: 'MEMBERSHIPS & MONETIZATION',
      links: [
        { href: '/admin/upgrades', label: 'Tier Upgrades', icon: ArrowUpCircle, visible: isPlatformOwner || isPlatformAdmin || isOperationsAdmin || hasAdminPermission('seller_upgrades') || hasAdminPermission('buyer_upgrades') },
        { href: '/admin/premium-users', label: 'Premium Subscribers', icon: Crown, visible: isPlatformOwner || isPlatformAdmin || hasAdminPermission('premium_users') },
        { href: '/admin/referrals', label: 'Referral Payouts', icon: Share2, visible: isPlatformOwner || isPlatformAdmin || hasAdminPermission('referrals') },
      ],
    },
    {
      title: 'FINANCE & PLANNING',
      links: [
        { href: '/admin/transactions', label: 'Transactions Ledger', icon: Receipt, visible: isPlatformOwner || isPlatformAdmin || isOperationsAdmin || hasAdminPermission('transactions') },
        { href: '/admin/finance', label: 'Financial Planning & FP&A', icon: TrendingUp, visible: isPlatformOwner || isPlatformAdmin || hasAdminPermission('profit') || hasAdminPermission('finance') },
        { href: '/admin/profit', label: 'Historical Profit Analytics', icon: DollarSign, visible: isPlatformOwner || isPlatformAdmin || hasAdminPermission('profit') },
      ],
    },
    {
      title: 'GOVERNANCE & SYSTEM',
      links: [
        { href: '/admin/users', label: 'Admin Accounts', icon: Users, visible: isPlatformOwner },
        { href: '/admin/settings', label: 'Platform Settings', icon: Settings, visible: isPlatformOwner || isPlatformAdmin || hasAdminPermission('settings') },
        { href: '/admin/logs', label: 'Security & Audit Logs', icon: FileText, visible: isPlatformOwner || isSecurityAdmin || hasAdminPermission('logs') },
      ],
    },
    {
      title: 'TRUST & COMPLIANCE',
      links: [
        { href: '/admin/verification', label: 'Provider Verification', icon: BadgeCheck, visible: isPlatformOwner || isPlatformAdmin || isOperationsAdmin || hasAdminPermission('seller_verification') },
        { href: '/admin/trust-safety', label: 'Trust & Safety', icon: AlertTriangle, visible: isPlatformOwner || isPlatformAdmin || isSecurityAdmin || isOperationsAdmin || isModerator || hasAdminPermission('trust_safety') },
        { href: '/admin/refunds', label: 'Refund Requests', icon: FileText, visible: isPlatformOwner || isPlatformAdmin || isOperationsAdmin || hasAdminPermission('refunds') },
        { href: '/admin/appeals', label: 'Account Appeals', icon: Scale, visible: isPlatformOwner || isPlatformAdmin || isOperationsAdmin || hasAdminPermission('appeals') },
      ],
    },
  ];

  return (
    <aside className={`w-72 bg-[#060D1A] text-slate-300 min-h-screen flex flex-col shrink-0 border-r border-slate-800/80 ${className}`}>
      <div className="p-4 sm:p-5 flex items-center justify-between border-b border-slate-800/80">
        <div className="flex items-center gap-3">
          <HuntiqLogo size="sm" href="/admin/dashboard" />
          <div>
            <h2 className="text-sm font-black text-white leading-none">HUNTIQ Admin</h2>
            <span className={`inline-block text-[9px] font-black uppercase tracking-wider mt-1 px-2 py-0.5 rounded-md border ${roleBadgeStyle}`}>
              {roleLabel}
            </span>
          </div>
        </div>
        {onClose && (
          <button
            onClick={onClose}
            className="lg:hidden p-1.5 rounded-xl text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
            aria-label="Close menu"
          >
            <X size={18} />
          </button>
        )}
      </div>
      
      <nav className="flex-1 p-3.5 space-y-5 overflow-y-auto admin-scrollbar">
        {navigationGroups.map((group, idx) => {
          const visibleLinks = group.links.filter((l) => l.visible);
          if (visibleLinks.length === 0) return null;

          return (
            <div key={idx} className="space-y-1.5">
              <div className="px-3 text-[9px] font-extrabold text-slate-500 uppercase tracking-widest">
                {group.title}
              </div>
              <div className="space-y-0.5">
                {visibleLinks.map((link) => {
                  const active = pathname === link.href || (link.href !== '/admin/dashboard' && pathname.startsWith(link.href));
                  return (
                    <Link
                      key={link.href}
                      href={link.href}
                      onClick={onClose}
                      className={`flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-xs font-bold transition-all ${
                        active
                          ? 'bg-[#1769E0]/20 text-[#38A5FF] border border-[#1769E0]/30 shadow-sm'
                          : 'text-slate-400 hover:bg-slate-900/90 hover:text-white'
                      }`}
                    >
                      <link.icon size={15} />
                      <span className="truncate">{link.label}</span>
                    </Link>
                  );
                })}
              </div>
            </div>
          );
        })}
      </nav>

      <div className="p-3.5 border-t border-slate-800/80">
        <button
          onClick={logout}
          className="flex items-center gap-3 px-3.5 py-2.5 w-full rounded-xl text-xs font-bold text-rose-400 hover:bg-rose-500/10 transition-all border border-rose-500/20"
        >
          <LogOut size={15} /> Sign Out Enterprise
        </button>
      </div>
    </aside>
  );
}

export default function AdminLayout({ children }: { children: ReactNode }) {
  const { user, loading, isAdmin } = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const [selectedBranch, setSelectedBranch] = useState('all');
  const [isWizardOpen, setIsWizardOpen] = useState(false);
  const [isMobileSidebarOpen, setIsMobileSidebarOpen] = useState(false);

  useEffect(() => {
    if (!loading && (!user || !isAdmin) && pathname !== '/admin/login') {
      router.push('/admin/login');
    }
  }, [user, loading, isAdmin, pathname, router]);

  if (loading) return <div className="h-screen flex items-center justify-center bg-slate-900"><Spinner /></div>;

  // Login page doesn't have the sidebar
  if (pathname === '/admin/login') {
    return <div className="min-h-screen bg-[#091426] text-white">{children}</div>;
  }

  // Must be logged in as admin for all other routes
  if (!user || !isAdmin) return null;

  return (
    <div className="flex min-h-screen bg-[#091426] font-sans text-slate-100 selection:bg-[#1769E0] selection:text-white">
      {/* Desktop Persistent Sidebar */}
      <AdminSidebar className="hidden lg:flex" />

      {/* Mobile Drawer Sidebar */}
      {isMobileSidebarOpen && (
        <div className="fixed inset-0 z-50 lg:hidden flex animate-fade-in">
          <div
            className="fixed inset-0 bg-slate-950/70 backdrop-blur-sm"
            onClick={() => setIsMobileSidebarOpen(false)}
            aria-hidden="true"
          />
          <div className="relative z-10 w-72 h-full shadow-2xl animate-slide-in-left">
            <AdminSidebar onClose={() => setIsMobileSidebarOpen(false)} />
          </div>
        </div>
      )}

      <main className="flex-1 min-w-0 flex flex-col h-screen overflow-y-auto admin-scrollbar">
        <AdminHeaderTopBar
          onOpenWizard={() => setIsWizardOpen(true)}
          selectedBranch={selectedBranch}
          onBranchChange={setSelectedBranch}
          onToggleMobileSidebar={() => setIsMobileSidebarOpen((prev) => !prev)}
        />
        <div className="flex-1 p-3.5 sm:p-6 lg:p-8 space-y-6">
          {children}
        </div>
      </main>

      {/* Global Form Wizard */}
      <MultiStepListingWizard
        isOpen={isWizardOpen}
        onClose={() => setIsWizardOpen(false)}
        onSuccess={() => setIsWizardOpen(false)}
      />
    </div>
  );
}
