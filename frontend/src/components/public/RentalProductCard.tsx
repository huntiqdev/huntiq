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

// ── Image blur logic ──────────────────────────────────────────────────────────
function useImageBlur(product: Product) {
  const { buyerTier } = useAuth();
  const sellerTier = product.seller_tier;
  if (!sellerTier || sellerTier === 'normal') return false;
  if (sellerTier === 'gold')    return buyerTier !== 'gold' && buyerTier !== 'premium';
  if (sellerTier === 'premium') return buyerTier !== 'premium';
  return false;
}

function BlurOverlay({ tier }: { tier: string }) {
  return (
    <div className="absolute inset-0 flex flex-col items-center justify-center bg-black/50 backdrop-blur-lg z-10 gap-1">
      <Lock size={14} className="text-white/80" />
      <p className="text-white text-[9px] font-black uppercase tracking-wider text-center px-1 leading-tight">
        {tier === 'premium' ? 'Premium' : 'Gold'}<br />buyer only
      </p>
      <Link href="/buyer/dashboard"
        className="text-[9px] font-black bg-orange-600 text-white px-2 py-0.5 rounded-full hover:bg-orange-500 transition-colors mt-0.5">
        Upgrade
      </Link>
    </div>
  );
}

export function RentalProductCard({ product }: Props) {
  const [imgIdx, setImgIdx] = useState(0);
  const blurred = useImageBlur(product);
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
    <Link href={`/products/${product.id}`}>
      <div
        className="group relative bg-white dark:bg-zinc-900 rounded-xl border border-zinc-100 dark:border-zinc-800 overflow-hidden flex flex-col cursor-pointer transition-all duration-200 hover:shadow-lg hover:shadow-black/10 hover:-translate-y-0.5 hover:border-blue-300 dark:hover:border-blue-700 active:scale-[0.98]"
      >
        {/* Image */}
        <div className="relative aspect-square bg-zinc-100 dark:bg-zinc-800 overflow-hidden">
          <Image
            src={thumb}
            alt={product.title}
            fill
            unoptimized={isUnoptimized(thumb)}
            className={`object-cover transition-transform duration-500 group-hover:scale-105 ${blurred ? 'blur-xl scale-110' : ''}`}
            sizes="(max-width: 480px) 33vw, (max-width: 768px) 25vw, (max-width: 1024px) 20vw, 16vw"
          />
          {blurred && <BlurOverlay tier={sellerTier ?? 'gold'} />}

          {/* Rental Price — bottom left */}
          <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 to-transparent px-1.5 pb-1 pt-4">
            <p className="text-white text-[10px] font-black leading-none">
              {product.rental_price
                ? `${product.rental_price >= 1000
                    ? `${(product.rental_price / 1000).toFixed(product.rental_price % 1000 === 0 ? 0 : 1)}k`
                    : product.rental_price.toLocaleString()} ETB/${product.rental_duration_value || 1}${durationLabel}`
                : 'Contact'}
            </p>
          </div>

          {/* Rental Badge — top left */}
          <div className="absolute top-1 left-1 bg-blue-600 text-white text-[8px] font-black px-1.5 py-0.5 rounded-full flex items-center gap-0.5 leading-none">
            <Clock size={7} />
            RENT
          </div>

          {/* Tier badge — just below rental badge */}
          {tierBadge && (
            <div className={`absolute top-6 left-1 w-4 h-4 rounded-full flex items-center justify-center text-[8px] font-black ${tierBadge.cls}`}>
              {tierBadge.label}
            </div>
          )}

          {/* Availability Status — top right */}
          <div className={`absolute top-1 right-1 text-white text-[8px] font-black px-1.5 py-0.5 rounded-full leading-none flex items-center gap-0.5 ${
            availabilityStatus === 'rented' ? 'bg-red-600' :
            availabilityStatus === 'soon' ? 'bg-yellow-600' :
            'bg-green-600'
          }`}>
            {availabilityStatus === 'rented' && rentedUntil && (
              <>Until {rentedUntil.getMonth() + 1}/{rentedUntil.getDate()}</>
            )}
            {availabilityStatus === 'soon' && availableFrom && (
              <>From {availableFrom.getMonth() + 1}/{availableFrom.getDate()}</>
            )}
            {availabilityStatus === 'available' && (
              <>
                <div className="w-1 h-1 bg-white rounded-full animate-pulse" />
                Now
              </>
            )}
          </div>

          {/* Deposit info — bottom right */}
          {product.deposit_amount && product.deposit_amount > 0 && (
            <div className="absolute bottom-5 right-1 bg-green-600/90 text-white text-[8px] font-black px-1 py-0.5 rounded flex items-center gap-0.5 leading-none">
              <Shield size={7} />
              {product.deposit_amount >= 1000 
                ? `${(product.deposit_amount / 1000).toFixed(0)}k` 
                : product.deposit_amount} dep
            </div>
          )}

          {/* Image dots — only if multiple */}
          {total > 1 && !blurred && (
            <div className="absolute bottom-5 left-1/2 -translate-x-1/2 flex gap-0.5">
              {Array.from({ length: total }).map((_, i) => (
                <button key={i}
                  onClick={e => { e.preventDefault(); e.stopPropagation(); setImgIdx(i); }}
                  className={`h-0.5 rounded-full transition-all ${i === imgIdx ? 'bg-white w-3' : 'bg-white/50 w-1'}`} />
              ))}
            </div>
          )}
        </div>

        {/* Body — ultra compact */}
        <div className="px-1.5 py-1.5 flex flex-col gap-0.5">
          <p className="text-[10px] font-bold leading-tight line-clamp-2 text-zinc-800 dark:text-zinc-200">
            {product.title}
          </p>
          {product.location && (
            <p className="text-[9px] text-zinc-400 flex items-center gap-0.5 leading-none truncate">
              <MapPin size={7} className="shrink-0" />
              {product.location.woreda ?? product.location.subcity ?? product.location.region}
            </p>
          )}
        </div>
      </div>
    </Link>
  );
}
