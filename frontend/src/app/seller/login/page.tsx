'use client';
import { useEffect, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

function Redirect() {
  const router = useRouter();
  const params = useSearchParams();
  useEffect(() => {
    const qs = new URLSearchParams({ as: 'seller' });
    const intent = params.get('intent');
    const condition = params.get('condition');
    if (intent) qs.set('intent', intent);
    if (condition) qs.set('condition', condition);
    router.replace(`/auth?${qs.toString()}`);
  }, [router, params]);
  return null;
}

export default function SellerLoginRedirect() {
  return <Suspense fallback={null}><Redirect /></Suspense>;
}
