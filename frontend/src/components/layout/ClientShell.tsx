'use client';
import { usePathname } from 'next/navigation';
import { Header } from '@/components/layout/Header';
import { Footer } from '@/components/layout/Footer';
import Link from 'next/link';
import { Home, Search, PlusCircle, User, Layers } from 'lucide-react';
import { useLocale } from '@/hooks/useLocale';
import { useAuth } from '@/context/AuthContext';
import { AIAssistant } from '@/components/ai/AIAssistant';

export function ClientShell({ children }: { children: React.ReactNode }) {
  const { t, locale } = useLocale();
  const { user, isSeller, isBuyer } = useAuth();
  const pathname = usePathname();
  const isAdmin = pathname.startsWith('/admin') || pathname.startsWith('/staff') || pathname.startsWith('/executive-console') || pathname.startsWith('/platform-console') || pathname.startsWith('/security-console');
  const isSellerArea = pathname.startsWith('/seller');
  const isBuyerArea = pathname.startsWith('/buyer');

  // Admin console has its own full-screen layout
  if (isAdmin) {
    return <>{children}</>;
  }

  // Seller and Buyer workspaces
  if (isSellerArea || isBuyerArea) {
    return (
      <>
        {children}
        <AIAssistant />
      </>
    );
  }

  const isProductDetail = pathname.startsWith('/products/') && pathname !== '/products';

  // Dynamic Profile destination for mobile navigation based on authenticated role
  const profileHref = user
    ? isSeller
      ? '/seller/dashboard'
      : isBuyer
      ? '/buyer/dashboard'
      : '/buyer/profile'
    : '/auth';
  const isProfileActive = pathname.startsWith('/seller') || pathname.startsWith('/buyer') || pathname === '/auth';

  return (
    <>
      <Header />
      <main className="min-h-screen pb-20 md:pb-0">{children}</main>
      <Footer />
      <AIAssistant />

      {/* Mobile Bottom Navigation — hidden on product detail pages to prioritize the sticky Buy/Contact CTA */}
      {!isProductDetail && (
        <nav
          aria-label="Mobile Bottom Navigation"
          className="md:hidden fixed bottom-0 left-0 right-0 bg-white/95 dark:bg-[#091426]/95 backdrop-blur-xl border-t border-slate-200 dark:border-slate-800 z-50 px-2 py-2 flex items-center justify-around pb-[max(12px,env(safe-area-inset-bottom,12px))] shadow-2xl"
        >
          {/* 1. Home */}
          <Link
            href="/"
            className={`flex flex-col items-center gap-1 tap-bounce active:scale-95 transition-all min-w-[56px] py-1 ${
              pathname === '/' ? 'text-[#1769E0] dark:text-[#38A5FF] font-black' : 'text-slate-500 dark:text-slate-400'
            }`}
          >
            <Home size={20} className={pathname === '/' ? 'fill-current' : ''} />
            <span className="text-[10px] font-bold">
              {locale === 'am' ? 'መነሻ' : locale === 'or' ? 'Mana' : locale === 'ti' ? 'መበገሲ' : 'Home'}
            </span>
          </Link>

          {/* 2. Categories Drawer */}
          <button
            type="button"
            onClick={() => {
              if (typeof window !== 'undefined') {
                window.dispatchEvent(new CustomEvent('open-category-drawer'));
              }
            }}
            className="flex flex-col items-center gap-1 tap-bounce active:scale-95 transition-all text-slate-500 dark:text-slate-400 hover:text-[#1769E0] min-w-[56px] py-1"
            aria-label="Open categories drawer"
          >
            <Layers size={20} />
            <span className="text-[10px] font-bold">{t.nav?.categories ?? 'Categories'}</span>
          </button>

          {/* 3. Sell / Rent Out (Elevated Center CTA) */}
          <Link
            href={user ? '/seller/product/new' : `/auth?as=seller&intent=sell&redirect=${encodeURIComponent('/seller/product/new')}`}
            className="flex flex-col items-center tap-bounce active:scale-95 transition-all group"
            aria-label="Post listing"
          >
            <div className="w-12 h-12 -mt-7 bg-gradient-to-tr from-[#1769E0] to-[#38A5FF] rounded-full flex items-center justify-center text-white shadow-lg shadow-[#1769E0]/35 border-4 border-white dark:border-[#091426] transition-transform group-hover:scale-105">
              <PlusCircle size={22} className="text-white stroke-[2.5px]" />
            </div>
            <span className="text-[10px] font-extrabold text-[#1769E0] dark:text-[#38A5FF] mt-0.5">
              {t.nav?.sell?.split('/')[0]?.trim() ?? 'Post'}
            </span>
          </Link>

          {/* 4. Buy / Rent */}
          <Link
            href="/products"
            className={`flex flex-col items-center gap-1 tap-bounce active:scale-95 transition-all min-w-[56px] py-1 ${
              pathname === '/products' || pathname === '/buy' || pathname === '/rent'
                ? 'text-[#1769E0] dark:text-[#38A5FF] font-black'
                : 'text-slate-500 dark:text-slate-400'
            }`}
          >
            <Search size={20} className={pathname === '/products' || pathname === '/buy' || pathname === '/rent' ? 'stroke-[2.5px]' : ''} />
            <span className="text-[10px] font-bold">{(t.nav as any)?.buy_rent ?? 'Buy / Rent'}</span>
          </Link>

          {/* 5. Account / Profile */}
          <Link
            href={profileHref}
            className={`flex flex-col items-center gap-1 tap-bounce active:scale-95 transition-all min-w-[56px] py-1 ${
              isProfileActive ? 'text-[#1769E0] dark:text-[#38A5FF] font-black' : 'text-slate-500 dark:text-slate-400'
            }`}
          >
            <User size={20} className={isProfileActive ? 'fill-current' : ''} />
            <span className="text-[10px] font-bold">{t.ui?.profile ?? 'Profile'}</span>
          </Link>
        </nav>
      )}
    </>
  );
}
