'use client';
import { usePathname } from 'next/navigation';
import { Header } from '@/components/layout/Header';
import { Footer } from '@/components/layout/Footer';
import Link from 'next/link';
import { Home, Search, PlusCircle, User, Sparkles } from 'lucide-react';
import { useLocale } from '@/hooks/useLocale';
import { useAuth } from '@/context/AuthContext';

export function ClientShell({ children }: { children: React.ReactNode }) {
  const { t, locale } = useLocale();
  const { user, isSeller, isBuyer } = useAuth();
  const pathname = usePathname();
  const isAdmin  = pathname.startsWith('/admin');
  const isSellerArea = pathname.startsWith('/seller');
  const isBuyerArea  = pathname.startsWith('/buyer');
  const isAI     = pathname.startsWith('/ai');

  // Admin, seller, buyer, and AI workspaces have their own dedicated full-screen layouts
  if (isAdmin || isSellerArea || isBuyerArea || isAI) {
    return <>{children}</>;
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
      
      {/* Floating AI Assistant Action Button — Dedicated single entrypoint at right bottom */}
      <Link
        href="/ai"
        aria-label="Open HUNTIQ AI Assistant"
        title="Open HUNTIQ AI Shopping Assistant"
        className={`fixed ${isProductDetail ? 'bottom-22 md:bottom-8' : 'bottom-20 md:bottom-8'} right-4 md:right-8 z-40 w-12 h-12 md:w-14 md:h-14 rounded-full bg-gradient-to-tr from-orange-600 via-amber-600 to-orange-500 text-white shadow-xl shadow-orange-600/35 hover:scale-110 active:scale-90 transition-all flex items-center justify-center border-2 border-white/20 group tap-bounce`}
      >
        <Sparkles size={22} className="animate-pulse" />
        <span className="sr-only">HUNTIQ AI</span>
      </Link>

      {/* Mobile Bottom Navigation — hidden on product detail pages to prioritize the sticky Buy/Contact CTA */}
      {!isProductDetail && (
        <nav className="md:hidden fixed bottom-0 left-0 right-0 bg-white/95 dark:bg-zinc-950/95 backdrop-blur-lg border-t border-zinc-200 dark:border-zinc-800 z-50 px-6 py-2.5 flex items-center justify-between pb-[max(12px,env(safe-area-inset-bottom,12px))] shadow-lg">
          <Link
            href="/"
            className={`flex flex-col items-center gap-1 tap-bounce active:scale-90 transition-all ${
              pathname === '/' ? 'text-orange-600 dark:text-orange-500 font-black' : 'text-zinc-500 dark:text-zinc-400'
            }`}
          >
            <Home size={20} className={pathname === '/' ? 'fill-current' : ''} />
            <span className="text-[10px] font-bold">
              {locale === 'am' ? 'መነሻ' : locale === 'or' ? 'Mana' : locale === 'ti' ? 'መበገሲ' : 'Home'}
            </span>
          </Link>

          <Link
            href="/products"
            className={`flex flex-col items-center gap-1 tap-bounce active:scale-90 transition-all ${
              pathname === '/products' || (pathname.startsWith('/products') && !isProductDetail)
                ? 'text-orange-600 dark:text-orange-500 font-black'
                : 'text-zinc-500 dark:text-zinc-400'
            }`}
          >
            <Search size={20} className={pathname === '/products' ? 'stroke-[2.5px]' : ''} />
            <span className="text-[10px] font-bold">{t.nav?.search ?? 'Find'}</span>
          </Link>

          <Link
            href="/sell"
            className={`flex flex-col items-center gap-1 tap-bounce active:scale-90 transition-all ${
              pathname === '/sell' ? 'text-orange-600 dark:text-orange-500 font-black' : 'text-zinc-500 dark:text-zinc-400'
            }`}
          >
            <div className="w-11 h-11 -mt-6 bg-gradient-to-tr from-orange-600 to-amber-500 rounded-full flex items-center justify-center text-white shadow-lg shadow-orange-600/35 border-4 border-white dark:border-zinc-950 transition-transform active:scale-95">
              <PlusCircle size={22} className="fill-current text-white" />
            </div>
            <span className="text-[10px] font-bold">{t.nav?.sell?.split('/')[0]?.trim() ?? 'Sell'}</span>
          </Link>

          <Link
            href={profileHref}
            className={`flex flex-col items-center gap-1 tap-bounce active:scale-90 transition-all ${
              isProfileActive ? 'text-orange-600 dark:text-orange-500 font-black' : 'text-zinc-500 dark:text-zinc-400'
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
