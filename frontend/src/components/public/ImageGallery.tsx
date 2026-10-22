'use client';
import { useState } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { ProductImage } from '@/types';
import { useLocale } from '@/hooks/useLocale';
import { useAuth } from '@/context/AuthContext';
import { Star, Crown, Lock } from 'lucide-react';
import { resolveImageUrl } from '@/lib/utils/imageUrl';
import { TIER_RANK } from '@/lib/constants/tiers';

interface Props {
  images: ProductImage[];
  /** The seller's tier — determines which images are gated */
  sellerTier?: string;
}

export function ImageGallery({ images, sellerTier = 'normal' }: Props) {
  const { t } = useLocale();
  const { buyerTier, isBuyer } = useAuth();
  const [active, setActive] = useState(0);

  if (!images.length) {
    return (
      <div className="aspect-[4/3] bg-zinc-100 dark:bg-zinc-800 rounded-2xl flex items-center justify-center opacity-40">
        <p className="text-sm">{t.ui.no_images}</p>
      </div>
    );
  }

  // Determine if the current buyer can view images from this seller tier
  const sellerRank = TIER_RANK[sellerTier] ?? 0;
  const buyerRank = TIER_RANK[buyerTier] ?? 0;
  // Normal seller images are always visible; gold/premium seller images require matching buyer tier
  const canView = sellerRank === 0 || buyerRank >= sellerRank;

  const TierIcon = sellerTier === 'premium' ? Crown : Star;
  const tierLabel = sellerTier === 'premium' ? 'Premium' : 'Gold';
  const tierGradient = sellerTier === 'premium'
    ? 'from-purple-500 via-violet-500 to-indigo-500'
    : 'from-yellow-400 via-amber-400 to-orange-400';
  const tierShadow = sellerTier === 'premium' ? 'shadow-purple-500/30' : 'shadow-yellow-500/30';

  return (
    <div className="flex flex-col gap-3">
      {/* Main image */}
      <div className="relative w-full aspect-[4/3] rounded-2xl overflow-hidden bg-zinc-100 dark:bg-zinc-800">
        {canView ? (
          <Image
            src={resolveImageUrl(images[active].image_url)}
            alt={`${active + 1}`}
            fill
            unoptimized={false}
            className="object-cover"
            sizes="(max-width: 768px) 100vw, 50vw"
            priority
          />
        ) : (
          <>
            {/* Blurred image underneath */}
            <Image
              src={resolveImageUrl(images[active].image_url)}
              alt="locked"
              fill
              unoptimized={false}
              className="object-cover blur-xl scale-110 opacity-60"
              sizes="(max-width: 768px) 100vw, 50vw"
            />
            {/* Lock overlay */}
            <div className="absolute inset-0 flex flex-col items-center justify-center gap-3 bg-black/40 backdrop-blur-sm">
              <div className={`w-14 h-14 rounded-2xl bg-gradient-to-br ${tierGradient} flex items-center justify-center shadow-xl ${tierShadow}`}>
                <TierIcon size={24} className="text-white" />
              </div>
              <div className="text-center px-4">
                <p className="font-black text-white text-sm">{tierLabel} Images Locked</p>
                <p className="text-white/70 text-xs mt-1">
                  Upgrade to {tierLabel} buyer to view
                </p>
              </div>
              <Link
                href={isBuyer ? '/buyer/upgrade' : '/buyer/register'}
                className={`flex items-center gap-1.5 text-xs font-black px-4 py-2 rounded-xl bg-gradient-to-r ${tierGradient} text-white shadow-lg ${tierShadow} hover:opacity-90 transition-all`}>
                <Lock size={11} />
                {isBuyer ? `Upgrade to ${tierLabel}` : 'Register as Buyer'}
              </Link>
            </div>
          </>
        )}

        {/* Tier badge on image (top-right) — only for gold/premium sellers */}
        {sellerRank > 0 && (
          <div className={`absolute top-2 right-2 flex items-center gap-1 text-[10px] font-black px-2 py-1 rounded-lg bg-gradient-to-r ${tierGradient} text-white shadow-md`}>
            <TierIcon size={9} />
            {tierLabel}
          </div>
        )}
      </div>

      {/* Thumbnails */}
      {images.length > 1 && (
        <div className="flex gap-2 overflow-x-auto pb-1 no-scrollbar">
          {images.map((img, i) => (
            <button key={img.id} onClick={() => canView && setActive(i)}
              className={`relative shrink-0 w-14 h-14 md:w-16 md:h-16 rounded-xl overflow-hidden border-2 transition-all ${
                i === active ? 'border-orange-600 shadow-lg shadow-orange-600/20' : 'border-transparent opacity-60 hover:opacity-100'
              } ${!canView ? 'cursor-not-allowed' : ''}`}>
              <Image src={resolveImageUrl(img.image_url)} alt={`${i + 1}`} fill unoptimized={false} className={`object-cover ${!canView ? 'blur-sm' : ''}`} sizes="64px" />
              {!canView && (
                <div className="absolute inset-0 flex items-center justify-center bg-black/30">
                  <Lock size={10} className="text-white" />
                </div>
              )}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
