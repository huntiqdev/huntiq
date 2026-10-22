'use client';
import { ReactNode, useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { useRouter, usePathname } from 'next/navigation';
import { Spinner } from '@/components/ui/Spinner';
import Link from 'next/link';
import { LayoutDashboard, Settings, Users, Package, LogOut, ShieldCheck, FileText, Activity, BadgeCheck, AlertTriangle, Scale, FolderTree, Tag, Search, MapPin } from 'lucide-react';
import { AdminHeaderTopBar } from '@/components/admin/enterprise/AdminHeaderTopBar';
import { MultiStepListingWizard } from '@/components/admin/enterprise/MultiStepListingWizard';

function AdminSidebar() {
  const { isSuperAdmin, hasAdminPermission, logout } = useAuth();
  const pathname = usePathname();

  const navigationGroups = [
    {
      title: 'OPERATIONS & CONTROL',
      links: [
        { href: '/admin/dashboard', label: 'Overview', icon: LayoutDashboard, visible: true },
        { href: '/admin/operations', label: 'Operations Monitor', icon: Activity, visible: isSuperAdmin || hasAdminPermission('operations') },
      ],
    },
    {
      title: 'LISTINGS & CATALOG',
      links: [
        { href: '/admin/products',  label: 'Products & Listings', icon: Package,   visible: hasAdminPermission('pending') || hasAdminPermission('products') },
        { href: '/admin/categories',label: 'Categories',          icon: FolderTree, visible: isSuperAdmin || hasAdminPermission('categories') },
        { href: '/admin/attributes', label: 'Attributes',         icon: Tag,        visible: isSuperAdmin || hasAdminPermission('categories') },
        { href: '/admin/locations', label: 'Locations & GPS',     icon: MapPin,     visible: isSuperAdmin || hasAdminPermission('locations') || hasAdminPermission('categories') },
        { href: '/admin/search',    label: 'Search Analytics',    icon: Search,     visible: isSuperAdmin || hasAdminPermission('analytics') },
      ],
    },
    {
      title: 'GOVERNANCE & SYSTEM',
      links: [
        { href: '/admin/users', label: 'Admin Accounts', icon: Users, visible: isSuperAdmin },
        { href: '/admin/settings', label: 'Platform Settings', icon: Settings, visible: isSuperAdmin || hasAdminPermission('settings') },
      ],
    },
    {
      title: 'TRUST & COMPLIANCE',
      links: [
        { href: '/admin/verification', label: 'Seller Verification', icon: BadgeCheck, visible: isSuperAdmin || hasAdminPermission('seller_verification') },
        { href: '/admin/trust-safety', label: 'Trust & Safety', icon: AlertTriangle, visible: isSuperAdmin || hasAdminPermission('trust_safety') },
        { href: '/admin/refunds', label: 'Refund Requests', icon: FileText, visible: isSuperAdmin || hasAdminPermission('refunds') },
        { href: '/admin/appeals', label: 'Account Appeals', icon: Scale, visible: isSuperAdmin || hasAdminPermission('appeals') },
      ],
    },
  ];

  return (
    <aside className="w-64 bg-zinc-950 text-zinc-300 min-h-screen flex flex-col shrink-0 border-r border-zinc-900/80 hidden lg:flex">
      <div className="p-5 flex items-center gap-3 border-b border-zinc-900">
        <div className="p-2 bg-orange-600/10 text-orange-500 rounded-xl border border-orange-500/20">
          <ShieldCheck size={22} />
        </div>
        <div>
          <h2 className="text-base font-black text-white leading-none">Enterprise Portal</h2>
          <p className="text-[10px] text-orange-400 font-bold uppercase tracking-widest mt-1">
            {isSuperAdmin ? 'Super Admin Mode' : 'Manager Portal'}
          </p>
        </div>
      </div>
      
      <nav className="flex-1 p-4 space-y-6 overflow-y-auto">
        {navigationGroups.map((group, idx) => {
          const visibleLinks = group.links.filter((l) => l.visible);
          if (visibleLinks.length === 0) return null;

          return (
            <div key={idx} className="space-y-2">
              <div className="px-3 text-[10px] font-extrabold text-zinc-500 uppercase tracking-wider">
                {group.title}
              </div>
              <div className="space-y-1">
                {visibleLinks.map((link) => {
                  const active = pathname === link.href || (link.href !== '/admin/dashboard' && pathname.startsWith(link.href));
                  return (
                    <Link
                      key={link.href}
                      href={link.href}
                      className={`flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-xs font-bold transition-all ${
                        active
                          ? 'bg-orange-600/15 text-orange-400 border border-orange-500/20'
                          : 'text-zinc-400 hover:bg-zinc-900 hover:text-white'
                      }`}
                    >
                      <link.icon size={16} />
                      {link.label}
                    </Link>
                  );
                })}
              </div>
            </div>
          );
        })}
      </nav>

      <div className="p-4 border-t border-zinc-900">
        <button
          onClick={logout}
          className="flex items-center gap-3 px-4 py-2.5 w-full rounded-xl text-xs font-bold text-red-400 hover:bg-red-500/10 transition-all border border-red-500/10"
        >
          <LogOut size={16} /> Sign Out Enterprise
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

  useEffect(() => {
    if (!loading && (!user || !isAdmin) && pathname !== '/admin/login') {
      router.push('/admin/login');
    }
  }, [user, loading, isAdmin, pathname, router]);

  if (loading) return <div className="h-screen flex items-center justify-center bg-zinc-50 dark:bg-zinc-950"><Spinner /></div>;

  // Login page doesn't have the sidebar
  if (pathname === '/admin/login') {
    return <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-100">{children}</div>;
  }

  // Must be logged in as admin for all other routes
  if (!user || !isAdmin) return null;

  return (
    <div className="flex min-h-screen bg-zinc-950 font-sans text-zinc-900 dark:text-zinc-100 selection:bg-orange-600 selection:text-white">
      <AdminSidebar />
      <main className="flex-1 min-w-0 flex flex-col h-screen overflow-y-auto">
        <AdminHeaderTopBar
          onOpenWizard={() => setIsWizardOpen(true)}
          selectedBranch={selectedBranch}
          onBranchChange={setSelectedBranch}
        />
        <div className="flex-1 p-4 lg:p-8 space-y-6">
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

