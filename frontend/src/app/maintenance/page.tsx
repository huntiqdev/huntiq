'use client';
import { Hammer, Clock, Phone } from 'lucide-react';
import Link from 'next/link';

export default function MaintenancePage() {
  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 flex items-center justify-center p-6">
      <div className="max-w-md w-full text-center space-y-8">
        {/* Animated Icon */}
        <div className="relative inline-block">
          <div className="absolute inset-0 bg-orange-500/20 blur-3xl rounded-full" />
          <div className="relative bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 p-8 rounded-3xl shadow-2xl">
            <Hammer size={48} className="text-orange-600 mx-auto animate-bounce" />
          </div>
        </div>

        <div className="space-y-4">
          <h1 className="text-4xl font-black text-zinc-900 dark:text-white tracking-tight">
            Under Maintenance
          </h1>
          <p className="text-zinc-500 dark:text-zinc-400 text-lg leading-relaxed">
            We're currently upgrading the platform to serve you better. We'll be back online shortly!
          </p>
        </div>

        {/* Status Card */}
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 shadow-sm flex flex-col gap-4">
          <div className="flex items-center gap-3 text-sm font-bold text-zinc-700 dark:text-zinc-300">
            <Clock size={18} className="text-orange-600" />
            Estimated Downtime: ~15 minutes
          </div>
          <div className="flex items-center gap-3 text-sm font-bold text-zinc-700 dark:text-zinc-300">
            <Phone size={18} className="text-orange-600" />
            Urgent? Call us: +251 9XX XXX XXX
          </div>
        </div>

        <div className="pt-4">
          <Link 
            href="/"
            className="inline-flex items-center justify-center px-8 py-3 bg-orange-600 hover:bg-orange-700 text-white font-black rounded-2xl transition-all shadow-lg shadow-orange-600/20 active:scale-95"
          >
            Refresh Page
          </Link>
        </div>

        <p className="text-[10px] text-zinc-400 font-bold uppercase tracking-widest">
          Hulubedeje Marketplace
        </p>
      </div>
    </div>
  );
}
