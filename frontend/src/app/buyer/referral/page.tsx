'use client';
import { ReferralCard } from '@/components/referral/ReferralCard';

export default function BuyerReferralPage() {
  return (
    <div className="max-w-2xl mx-auto space-y-4">
      <div>
        <h1 className="text-lg font-black text-zinc-900 dark:text-white">Referral Program</h1>
        <p className="text-sm text-zinc-500 dark:text-zinc-400 mt-0.5">
          Invite others and earn 2% commission on their paid actions.
        </p>
      </div>
      <ReferralCard role="buyer" />
    </div>
  );
}
