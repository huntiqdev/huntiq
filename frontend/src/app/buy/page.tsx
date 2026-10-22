'use client';
import { useEffect, useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

function BuyContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    // Redirect to main marketplace with sale filter
    const params = new URLSearchParams(searchParams?.toString() || '');
    params.set('type', 'sale');
    router.replace(`/products?${params.toString()}`);
  }, [router, searchParams]);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        <p className="mt-4 text-slate-600 dark:text-slate-300">Loading marketplace...</p>
      </div>
    </div>
  );
}

// This is a placeholder - the actual buy page already exists
// This routes to the existing marketplace with SALE filter
export default function BuyPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
      </div>
    }>
      <BuyContent />
    </Suspense>
  );
}
