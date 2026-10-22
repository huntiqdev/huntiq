'use client';
import { useEffect, useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

function RentContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    // Redirect to main marketplace with rent filter
    const params = new URLSearchParams(searchParams?.toString() || '');
    params.set('type', 'rent');
    router.replace(`/products?${params.toString()}`);
  }, [router, searchParams]);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mx-auto"></div>
        <p className="mt-4 text-slate-600 dark:text-slate-300">Loading rentals...</p>
      </div>
    </div>
  );
}

// This routes to the marketplace with RENT filter
export default function RentPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mx-auto"></div>
      </div>
    }>
      <RentContent />
    </Suspense>
  );
}
