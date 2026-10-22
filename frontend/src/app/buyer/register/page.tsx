'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function BuyerRegisterRedirect() {
  const router = useRouter();
  useEffect(() => { router.replace('/auth?as=buyer'); }, [router]);
  return null;
}
