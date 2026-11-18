import Link from 'next/link';
import { FileQuestion, Home, ArrowLeft, Search } from 'lucide-react';

export default function NotFound() {
  return (
    <div className="min-h-[75vh] flex items-center justify-center px-4 py-16">
      <div className="max-w-md w-full text-center space-y-6">
        <div className="w-20 h-20 mx-auto rounded-3xl bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 flex items-center justify-center shadow-lg">
          <FileQuestion size={40} />
        </div>

        <div className="space-y-2">
          <span className="text-xs font-black tracking-widest uppercase text-[#1769E0] dark:text-[#38A5FF]">
            404 Error
          </span>
          <h1 className="text-3xl font-black text-zinc-900 dark:text-zinc-100">
            Page Not Found
          </h1>
          <p className="text-sm text-zinc-500 dark:text-zinc-400 leading-relaxed">
            The page you are looking for does not exist, has been moved, or is temporarily unavailable.
          </p>
        </div>

        <div className="pt-2 flex flex-col sm:flex-row gap-2.5 justify-center">
          <Link
            href="/"
            className="inline-flex items-center justify-center gap-2 px-5 py-2.5 min-h-[44px] rounded-xl bg-[#1769E0] hover:bg-[#1055BD] text-white font-bold text-sm shadow-md shadow-[#1769E0]/20 transition-all tap-bounce active:scale-95"
          >
            <Home size={16} />
            <span>Return to Home</span>
          </Link>
          <Link
            href="/products"
            className="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-800 dark:text-zinc-200 font-bold text-sm transition-all"
          >
            <Search size={16} />
            <span>Browse Marketplace</span>
          </Link>
        </div>
      </div>
    </div>
  );
}

