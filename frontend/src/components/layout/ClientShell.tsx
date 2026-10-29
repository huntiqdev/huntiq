'use client';
import { usePathname } from 'next/navigation';
import { Header } from '@/components/layout/Header';
import { Footer } from '@/components/layout/Footer';
import Link from 'next/link';
import { Home, Search, PlusCircle, User, Sparkles } from 'lucide-react';

export function ClientShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const isAdmin  = pathname.startsWith('/admin');
  const isSeller = pathname.startsWith('/seller');
  const isBuyer  = pathname.startsWith('/buyer');
  const isAI     = pathname.startsWith('/ai');

  // Admin, seller, buyer, and AI workspaces have their own dedicated full-screen layouts
  if (isAdmin || isSeller || isBuyer || isAI) {
    return <>{children}</>;
  }

  const isProductDetail = pathname.startsWith('/products/') && pathname !== '/products';

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
        className="fixed bottom-20 md:bottom-8 right-4 md:right-8 z-40 w-12 h-12 md:w-14 md:h-14 rounded-full bg-gradient-to-tr from-orange-600 via-amber-600 to-orange-500 text-white shadow-xl shadow-orange-600/35 hover:scale-110 active:scale-95 transition-all flex items-center justify-center border-2 border-white/20 group"
      >
        <Sparkles size={22} className="animate-pulse" />
        <span className="sr-only">HUNTIQ AI</span>
      </Link>

      {/* Mobile Bottom Navigation — hidden on product detail pages to prioritize the sticky Buy/Contact CTA */}
      {!isProductDetail && (
        <nav className="md:hidden fixed bottom-0 left-0 right-0 bg-white/95 dark:bg-zinc-950/95 backdrop-blur-lg border-t border-zinc-200 dark:border-zinc-800 z-50 px-6 py-3 flex items-center justify-between pb-[env(safe-area-inset-bottom,12px)]">
          <Link href="/" className={`flex flex-col items-center gap-1 ${pathname === '/' ? 'text-orange-600' : 'text-zinc-500 dark:text-zinc-400'}`}>
            <Home size={20} className={pathname === '/' ? 'fill-current' : ''} />
            <span className="text-[10px] font-bold">Home</span>
          </Link>
          <Link href="/products" className={`flex flex-col items-center gap-1 ${pathname === '/products' || (pathname.startsWith('/products') && !isProductDetail) ? 'text-orange-600' : 'text-zinc-500 dark:text-zinc-400'}`}>
            <Search size={20} className={pathname === '/products' ? 'stroke-[3px]' : ''} />
            <span className="text-[10px] font-bold">Find</span>
          </Link>
          <Link href="/sell" className={`flex flex-col items-center gap-1 ${pathname === '/sell' ? 'text-orange-600' : 'text-zinc-500 dark:text-zinc-400'}`}>
            <div className="w-10 h-10 -mt-5 bg-orange-600 rounded-full flex items-center justify-center text-white shadow-lg shadow-orange-600/30 border-4 border-white dark:border-zinc-950">
              <PlusCircle size={20} className="fill-current" />
            </div>
            <span className="text-[10px] font-bold">Sell</span>
          </Link>
          <Link href="/seller/login" className={`flex flex-col items-center gap-1 ${pathname.startsWith('/seller') ? 'text-orange-600' : 'text-zinc-500 dark:text-zinc-400'}`}>
            <User size={20} className={pathname.startsWith('/seller') ? 'fill-current' : ''} />
            <span className="text-[10px] font-bold">Profile</span>
          </Link>
        </nav>
      )}
    </>
  );
}
