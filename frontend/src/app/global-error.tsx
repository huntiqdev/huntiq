'use client';

import { useEffect } from 'react';
import { AlertCircle, RotateCcw } from 'lucide-react';

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error('Critical global error captured:', error);
  }, [error]);

  return (
    <html lang="en">
      <body className="bg-white dark:bg-zinc-950 text-zinc-900 dark:text-zinc-100 min-h-screen flex items-center justify-center p-4">
        <div className="max-w-md w-full text-center space-y-6">
          <div className="w-20 h-20 mx-auto rounded-3xl bg-red-100 dark:bg-red-950 text-red-600 dark:text-red-400 flex items-center justify-center shadow-lg">
            <AlertCircle size={40} />
          </div>
          <div className="space-y-2">
            <h1 className="text-2xl font-black">Critical Error</h1>
            <p className="text-sm text-zinc-500 dark:text-zinc-400">
              An unexpected system error occurred. Please refresh or try again.
            </p>
          </div>
          <button
            onClick={() => reset()}
            className="inline-flex items-center gap-2 px-5 py-2.5 min-h-[44px] rounded-xl bg-[#1769E0] hover:bg-[#1055BD] text-white font-bold text-sm shadow-md shadow-[#1769E0]/20 transition-all tap-bounce active:scale-95"
          >
            <RotateCcw size={16} />
            <span>Reload Application</span>
          </button>
        </div>
      </body>
    </html>
  );
}

