import Link from 'next/link';
import { PackageX, ArrowLeft, Search, Camera, Smartphone, Car, Home } from 'lucide-react';

export default function ProductNotFound() {
  return (
    <div className="min-h-[70vh] flex items-center justify-center px-4 py-12">
      <div className="max-w-md w-full text-center space-y-6">
        <div className="w-20 h-20 mx-auto rounded-3xl bg-orange-100 dark:bg-orange-950/40 text-orange-600 dark:text-orange-400 flex items-center justify-center shadow-lg shadow-orange-600/10">
          <PackageX size={40} />
        </div>

        <div className="space-y-2">
          <h1 className="text-2xl font-black text-zinc-900 dark:text-zinc-100">
            Product Not Found
          </h1>
          <p className="text-sm text-zinc-500 dark:text-zinc-400 leading-relaxed">
            This listing may have been sold, expired, or removed by the seller.
          </p>
        </div>

        <div className="pt-2 flex flex-col sm:flex-row gap-2.5 justify-center">
          <Link
            href="/products"
            className="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-orange-600 hover:bg-orange-700 text-white font-bold text-sm shadow-md shadow-orange-600/20 transition-all"
          >
            <ArrowLeft size={16} />
            <span>Browse Products</span>
          </Link>
          <Link
            href="/"
            className="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-800 dark:text-zinc-200 font-bold text-sm transition-all"
          >
            <Home size={16} />
            <span>Go to Home</span>
          </Link>
        </div>

        <div className="pt-6 border-t border-zinc-200 dark:border-zinc-800">
          <p className="text-xs font-black uppercase tracking-wider text-zinc-400 mb-3">
            Explore Popular Categories
          </p>
          <div className="flex flex-wrap gap-2 justify-center">
            <Link
              href="/products?category_hierarchy_id=1380"
              className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-zinc-100 dark:bg-zinc-800/80 hover:bg-orange-50 hover:text-orange-600 dark:hover:bg-orange-950/30 text-xs font-semibold text-zinc-700 dark:text-zinc-300 transition-colors"
            >
              <Camera size={13} className="text-orange-500" />
              <span>Camera Lenses</span>
            </Link>
            <Link
              href="/products?cat=electronics"
              className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-zinc-100 dark:bg-zinc-800/80 hover:bg-orange-50 hover:text-orange-600 dark:hover:bg-orange-950/30 text-xs font-semibold text-zinc-700 dark:text-zinc-300 transition-colors"
            >
              <Smartphone size={13} className="text-orange-500" />
              <span>Smartphones</span>
            </Link>
            <Link
              href="/products?cat=transport-mobility"
              className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-zinc-100 dark:bg-zinc-800/80 hover:bg-orange-50 hover:text-orange-600 dark:hover:bg-orange-950/30 text-xs font-semibold text-zinc-700 dark:text-zinc-300 transition-colors"
            >
              <Car size={13} className="text-orange-500" />
              <span>Vehicles</span>
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

