'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function SellerRegisterRedirect() {
  const router = useRouter();
  useEffect(() => { router.replace('/auth?as=seller'); }, [router]);
  return null;
}
