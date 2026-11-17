'use client';

import { useEffect } from 'react';
import Link from 'next/link';
import { AlertTriangle, RotateCcw, Home } from 'lucide-react';

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function ErrorPage({ error, reset }: ErrorProps) {
  useEffect(() => {
    // Log error to client-side diagnostics
    console.error('App-level render error captured:', error);
  }, [error]);

  return (
    <div className="min-h-[70vh] flex items-center justify-center px-4 py-16">
      <div className="max-w-md w-full text-center space-y-6">
        <div className="w-20 h-20 mx-auto rounded-3xl bg-red-100 dark:bg-red-950/40 text-red-600 dark:text-red-400 flex items-center justify-center shadow-lg shadow-red-600/10">
          <AlertTriangle size={40} />
        </div>

        <div className="space-y-2">
          <span className="text-xs font-black tracking-widest uppercase text-red-600 dark:text-red-400">
            System Notice
          </span>
          <h1 className="text-2xl font-black text-zinc-900 dark:text-zinc-100">
            Unable to Load Content
          </h1>
          <p className="text-sm text-zinc-500 dark:text-zinc-400 leading-relaxed">
            {error?.message && !error.message.includes('Minified React')
              ? error.message
              : 'A temporary rendering or connection issue occurred while displaying this page.'}
          </p>
        </div>

        <div className="pt-2 flex flex-col sm:flex-row gap-2.5 justify-center">
          <button
            onClick={() => reset()}
            className="inline-flex items-center justify-center gap-2 px-5 py-2.5 min-h-[44px] rounded-xl bg-[#1769E0] hover:bg-[#1055BD] text-white font-bold text-sm shadow-md shadow-[#1769E0]/20 transition-all tap-bounce active:scale-95"
          >
            <RotateCcw size={16} />
            <span>Try Again</span>
          </button>
          <Link
            href="/"
            className="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-800 dark:text-zinc-200 font-bold text-sm transition-all"
          >
            <Home size={16} />
            <span>Go to Home</span>
          </Link>
        </div>
      </div>
    </div>
  );
}

