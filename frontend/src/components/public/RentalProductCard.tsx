'use client';
import { useState } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { Product } from '@/types';
import {
  MapPin, Clock, Calendar, Shield, ChevronLeft, ChevronRight, Lock,
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { resolveImageUrl, isUnoptimized } from '@/lib/utils/imageUrl';

interface Props { product: Product }

const DURATION_LABELS: Record<string, string> = {
  hour: 'hr',
  day: 'day',
  week: 'wk',
  month: 'mo',
  year: 'yr',
};

// ── Tier badge helpers ────────────────────────────────────────────────────────
const TIER_BADGE: Record<string, { label: string; cls: string }> = {
  gold:    { label: 'G', cls: 'bg-yellow-500 text-white' },
  premium: { label: 'P', cls: 'bg-purple-600 text-white' },
};

// ── Tier indicator helpers for restricted specifications ──────────────────────
function TierSpecBadge({ tier }: { tier?: string }) {
  if (!tier || tier === 'normal') return null;
  const isGold = tier === 'gold';
  return (
    <div
      className={`absolute top-1 right-1 z-10 px-1.5 py-0.5 rounded text-[8px] font-black uppercase tracking-wider flex items-center gap-0.5 shadow-sm ${
        isGold ? 'bg-amber-500 text-white' : 'bg-purple-600 text-white'
      }`}
    >
      <Lock size={7} />
      <span>{isGold ? 'Gold' : 'Premium'}</span>
    </div>
  );
}

export function RentalProductCard({ product }: Props) {
  const [imgIdx, setImgIdx] = useState(0);
  const sellerTier = product.seller_tier;
  const tierBadge = sellerTier ? TIER_BADGE[sellerTier] : null;

  const images = product.images ?? [];
  const total = Math.min(images.length, 3);
  const thumb = resolveImageUrl(images[imgIdx]?.image_url);

  const durationLabel = product.rental_duration_unit 
    ? DURATION_LABELS[product.rental_duration_unit] || product.rental_duration_unit
    : 'period';

  const isCurrentlyRented = product.is_currently_rented;
  const rentedUntil = product.rented_until ? new Date(product.rented_until) : null;
  const availableFrom = product.available_from ? new Date(product.available_from) : null;
  const now = new Date();

  // Availability badge
  let availabilityStatus: 'available' | 'rented' | 'soon' = 'available';
  if (isCurrentlyRented && rentedUntil) availabilityStatus = 'rented';
  else if (availableFrom && availableFrom > now) availabilityStatus = 'soon';

  return (
    <Link href={`/products/${product.id}`} className="block">
      <div
        className="group relative bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-100 dark:border-zinc-800/80 overflow-hidden flex flex-col cursor-pointer transition-all duration-200 hover:shadow-lg hover:shadow-black/10 hover:-translate-y-0.5 hover:border-teal-300 dark:hover:border-teal-800/60 tap-bounce active:scale-[0.97]"
      >
        {/* Image */}
        <div className="relative aspect-square bg-zinc-100 dark:bg-zinc-800 overflow-hidden">
          <Image
            src={thumb}
            alt={product.title}
            fill
            unoptimized={isUnoptimized(thumb)}
            className="object-cover transition-transform duration-500 group-hover:scale-105"
            sizes="(max-width: 480px) 50vw, (max-width: 768px) 33vw, (max-width: 1024px) 25vw, 16vw"
          />
          <TierSpecBadge tier={product.required_buyer_tier} />

          {/* Rental Price — bottom left */}
          <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent px-2.5 pb-2 pt-6">
            <p className="text-white text-xs sm:text-sm font-black leading-tight drop-shadow-md">
              {product.rental_price
                ? `${product.rental_price >= 1000
                    ? `${(product.rental_price / 1000).toFixed(product.rental_price % 1000 === 0 ? 0 : 1)}k`
                    : product.rental_price.toLocaleString()} ETB/${product.rental_duration_value || 1}${durationLabel}`
                : 'Contact'}
            </p>
          </div>

          {/* Rental Badge — top left */}
          <div className="absolute top-2 left-2 bg-teal-600 text-white text-[9px] font-black px-2 py-0.5 rounded-full flex items-center gap-1 leading-none shadow-sm tracking-wide">
            <Clock size={8} />
            RENT
          </div>

          {/* Tier badge — just below rental badge */}
          {tierBadge && (
            <div className={`absolute top-7 left-2 w-4 h-4 rounded-full flex items-center justify-center text-[8px] font-black shadow-sm ${tierBadge.cls}`}>
              {tierBadge.label}
            </div>
          )}

          {/* Availability Status — top right */}
          <div className={`absolute top-2 right-2 text-white text-[9px] font-bold px-1.5 py-0.5 rounded-md leading-none flex items-center gap-1 shadow-sm ${
            availabilityStatus === 'rented' ? 'bg-red-600' :
            availabilityStatus === 'soon' ? 'bg-amber-600' :
            'bg-emerald-600'
          }`}>
            {availabilityStatus === 'rented' && rentedUntil && (
              <>Until {rentedUntil.getMonth() + 1}/{rentedUntil.getDate()}</>
            )}
            {availabilityStatus === 'soon' && availableFrom && (
              <>From {availableFrom.getMonth() + 1}/{availableFrom.getDate()}</>
            )}
            {availabilityStatus === 'available' && (
              <>
                <div className="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
                Now
              </>
            )}
          </div>

          {/* Deposit info — bottom right */}
          {product.deposit_amount && product.deposit_amount > 0 && (
            <div className="absolute bottom-7 right-2 bg-emerald-600/90 text-white text-[8px] font-black px-1.5 py-0.5 rounded flex items-center gap-1 leading-none shadow-sm">
              <Shield size={7} />
              {product.deposit_amount >= 1000 
                ? `${(product.deposit_amount / 1000).toFixed(0)}k` 
                : product.deposit_amount} dep
            </div>
          )}

          {/* Image dots — only if multiple */}
          {total > 1 && (
            <div className="absolute bottom-7 left-1/2 -translate-x-1/2 flex gap-1">
              {Array.from({ length: total }).map((_, i) => (
                <button key={i}
                  onClick={e => { e.preventDefault(); e.stopPropagation(); setImgIdx(i); }}
                  className={`h-1 rounded-full transition-all ${i === imgIdx ? 'bg-white w-3.5' : 'bg-white/50 w-1'}`} />
              ))}
            </div>
          )}
        </div>

        {/* Body — smartphone readable typography */}
        <div className="p-2.5 flex flex-col gap-1">
          <p className="text-xs sm:text-sm font-semibold leading-snug line-clamp-2 text-zinc-900 dark:text-zinc-100 group-hover:text-teal-600 transition-colors">
            {product.title}
          </p>
          {product.location && (
            <p className="text-[11px] sm:text-xs text-zinc-500 dark:text-zinc-400 flex items-center gap-1 leading-none truncate font-medium">
              <MapPin size={10} className="shrink-0 text-teal-600" />
              <span>{product.location.woreda ?? product.location.subcity ?? product.location.region}</span>
            </p>
          )}
        </div>
      </div>
    </Link>
  );
}
