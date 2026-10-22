'use client';

import { useNetworkStatus } from '@/hooks/useNetworkStatus';

export function NetworkStatusBanner() {
  const { isOnline, wasOffline } = useNetworkStatus();

  if (isOnline && !wasOffline) return null;

  return (
    <div
      className={`fixed top-0 left-0 right-0 z-50 px-4 py-3 text-center font-medium transition-colors ${
        isOnline
          ? 'bg-green-500 text-white'
          : 'bg-red-500 text-white'
      }`}
    >
      {isOnline ? (
        <span>✓ Back online</span>
      ) : (
        <span>⚠️ No internet connection</span>
      )}
    </div>
  );
}
