'use client';
import { useState } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { Product } from '@/types';
import {
  MapPin, ShieldCheck, ShoppingBag, X,
  ChevronLeft, ChevronRight, Lock,
} from 'lucide-react';
import { BuyModal } from '@/components/public/BuyModal';
import { useLocale } from '@/hooks/useLocale';
import { useAuth } from '@/context/AuthContext';
import { resolveImageUrl, isUnoptimized } from '@/lib/utils/imageUrl';

interface Props { product: Product }

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
        isGold
          ? 'bg-amber-500 text-white'
          : 'bg-purple-600 text-white'
      }`}
    >
      <Lock size={7} />
      <span>{isGold ? 'Gold' : 'Premium'}</span>
    </div>
  );
}

// ── Detail Drawer ─────────────────────────────────────────────────────────────
function DetailDrawer({ product, onClose, onBuy }: {
  product: Product; onClose: () => void; onBuy: () => void;
}) {
  const { t } = useLocale();
  const images = product.images ?? [];
  const [imgIdx, setImgIdx] = useState(0);
  const total = images.length;
  const sellerTier = product.seller_tier;

  return (
    <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4">
      <div className="absolute inset-0 bg-black/70 backdrop-blur-sm" onClick={onClose} />
      <div className="relative z-10 w-full sm:max-w-lg bg-white dark:bg-zinc-900 rounded-t-3xl sm:rounded-3xl overflow-hidden shadow-2xl max-h-[92vh] flex flex-col">

        {/* Mobile Grab Handle Bar */}
        <div className="w-12 h-1.5 rounded-full bg-zinc-300 dark:bg-zinc-700 mx-auto my-2.5 sm:hidden shrink-0" />

        {/* Image carousel */}
        <div className="relative aspect-[4/3] bg-zinc-100 dark:bg-zinc-800 shrink-0">
          {images.length > 0 ? (
          <Image src={resolveImageUrl(images[imgIdx]?.image_url)}
              alt={product.title} fill
              unoptimized={isUnoptimized(resolveImageUrl(images[imgIdx]?.image_url))}
              className="object-cover"
              sizes="(max-width: 640px) 100vw, 512px" />
          ) : (
            <div className="absolute inset-0 flex items-center justify-center opacity-20">
              <ShoppingBag size={48} />
            </div>
          )}
          <TierSpecBadge tier={product.required_buyer_tier} />

          {total > 1 && (
            <>
              <button onClick={e => { e.stopPropagation(); setImgIdx(i => (i - 1 + total) % total); }}
                className="absolute left-2 top-1/2 -translate-y-1/2 w-8 h-8 rounded-full bg-black/40 text-white flex items-center justify-center hover:bg-black/60 transition-colors tap-bounce active:scale-90">
                <ChevronLeft size={18} />
              </button>
              <button onClick={e => { e.stopPropagation(); setImgIdx(i => (i + 1) % total); }}
                className="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 rounded-full bg-black/40 text-white flex items-center justify-center hover:bg-black/60 transition-colors tap-bounce active:scale-90">
                <ChevronRight size={18} />
              </button>
              <div className="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1.5">
                {images.map((_, i) => (
                  <button key={i} onClick={e => { e.stopPropagation(); setImgIdx(i); }}
                    className={`h-1.5 rounded-full transition-all ${i === imgIdx ? 'bg-white w-5' : 'bg-white/50 w-1.5'}`} />
                ))}
              </div>
            </>
          )}

          <button onClick={onClose}
            className="absolute top-3 right-3 w-8 h-8 rounded-full bg-black/40 text-white flex items-center justify-center hover:bg-black/60 transition-colors tap-bounce active:scale-90">
            <X size={16} />
          </button>
          <div className="absolute top-3 left-3 bg-orange-600 text-white px-3 py-1 rounded-full text-sm font-black shadow-lg">
            ETB {product.price.toLocaleString()}
          </div>
        </div>

        {/* Info */}
        <div className="overflow-y-auto p-5 flex flex-col gap-3 flex-1">
          <div className="flex flex-wrap gap-1.5">
            {product.category && (
              <span className="text-[10px] font-black uppercase tracking-widest text-orange-600 bg-orange-600/10 px-2.5 py-1 rounded-lg">
                {product.category.name}
              </span>
            )}
            {product.subcategory && (
              <span className="text-[10px] font-black uppercase tracking-widest text-purple-600 bg-purple-600/10 px-2.5 py-1 rounded-lg">
                {(t.sub_categories as any)?.[(product.category as any)?.slug ?? '']?.[product.subcategory]
                  ?? (t.items as any)?.[product.subcategory]
                  ?? product.subcategory.replace(/_/g, ' ')}
              </span>
            )}
            <span className="text-[10px] font-black uppercase tracking-widest bg-zinc-100 dark:bg-zinc-800 px-2.5 py-1 rounded-lg capitalize">
              {(t.ui as any)?.[`condition_${product.condition}`] ?? product.condition.replace(/_/g, ' ')}
            </span>
            <span className="text-[10px] font-black uppercase tracking-widest text-blue-600 bg-blue-600/10 px-2.5 py-1 rounded-lg flex items-center gap-1">
              <ShieldCheck size={10} /> {t.ui.verified}
            </span>
          </div>

          <h2 className="text-lg font-black leading-tight text-zinc-900 dark:text-zinc-100">{product.title}</h2>

          {product.location && (
            <div className="flex items-center gap-1.5 text-xs text-zinc-500 dark:text-zinc-400">
              <MapPin size={13} className="text-orange-500 shrink-0" />
              <span>{[product.location.woreda, product.location.subcity, product.location.region].filter(Boolean).join(', ')}</span>
            </div>
          )}

          {product.description && (
            <p className="text-xs sm:text-sm leading-relaxed text-zinc-600 dark:text-zinc-300 opacity-90">{product.description}</p>
          )}
        </div>

        {/* Pinned Priority Action Bar */}
        <div className="shrink-0 p-4 border-t border-zinc-100 dark:border-zinc-800 bg-white/95 dark:bg-zinc-900/95 backdrop-blur flex gap-3 pb-[max(16px,env(safe-area-inset-bottom,16px))]">
          <button type="button" onClick={onBuy}
            className={`flex-1 text-white rounded-2xl py-3 font-bold text-sm flex items-center justify-center gap-2 transition-all tap-bounce active:scale-95 shadow-md ${
              product.listing_type === 'rent'
                ? 'bg-teal-600 hover:bg-teal-700 shadow-teal-600/25'
                : 'bg-orange-600 hover:bg-orange-700 shadow-orange-600/25'
            }`}>
            <ShoppingBag size={16} />
            {product.listing_type === 'rent' ? (t.ui.rent_now ?? 'Rent / Contact Landlord') : t.ui.buy_now}
          </button>
          <Link href={`/products/${product.id}`}
            className="px-4 py-3 rounded-2xl border border-zinc-200 dark:border-zinc-700 font-bold text-sm hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors tap-bounce active:scale-95 flex items-center justify-center">
            {t.ui.details}
          </Link>
        </div>
      </div>
    </div>
  );
}

// ── Card ──────────────────────────────────────────────────────────────────────
function expiryDaysLeft(expiresAt?: string): number | null {
  if (!expiresAt) return null;
  const diff = new Date(expiresAt).getTime() - Date.now();
  if (diff <= 0) return 0;
  return Math.ceil(diff / 86_400_000);
}

export function ProductCard({ product }: Props) {
  const [imgIdx, setImgIdx] = useState(0);
  const [detailOpen, setDetailOpen] = useState(false);
  const [buyOpen, setBuyOpen] = useState(false);
  const sellerTier = product.seller_tier;
  const tierBadge = sellerTier ? TIER_BADGE[sellerTier] : null;
  const daysLeft = expiryDaysLeft(product.expires_at);

  const images = product.images ?? [];
  const total = Math.min(images.length, 3);
  const thumb = resolveImageUrl(images[imgIdx]?.image_url);

  return (
    <>
      {/* Card — smartphone thumb-friendly, tactile bounce */}
      <div
        className="group relative bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-100 dark:border-zinc-800/80 overflow-hidden flex flex-col cursor-pointer transition-all duration-200 hover:shadow-lg hover:shadow-black/10 hover:-translate-y-0.5 hover:border-orange-300 dark:hover:border-orange-800/60 tap-bounce active:scale-[0.97]"
        onClick={() => setDetailOpen(true)}
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

          {/* Price — bottom left with high contrast shadow */}
          <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent px-2.5 pb-2 pt-6">
            <p className="text-white text-xs sm:text-sm font-black leading-tight drop-shadow-md">
              ETB {product.price >= 1000
                ? `${(product.price / 1000).toFixed(product.price % 1000 === 0 ? 0 : 1)}k`
                : product.price.toLocaleString()}
            </p>
          </div>

          {/* Sale Badge — top left */}
          <div className="absolute top-2 left-2 bg-orange-600 text-white text-[9px] font-black px-2 py-0.5 rounded-full flex items-center gap-1 leading-none shadow-sm tracking-wide">
            <ShoppingBag size={8} />
            SALE
          </div>

          {/* Tier badge — just below sale badge */}
          {tierBadge && (
            <div className={`absolute top-7 left-2 w-4 h-4 rounded-full flex items-center justify-center text-[8px] font-black shadow-sm ${tierBadge.cls}`}>
              {tierBadge.label}
            </div>
          )}

          {/* Condition — top right (offset if tier badge is present) */}
          <div className={`absolute ${product.required_buyer_tier && product.required_buyer_tier !== 'normal' ? 'top-7 right-2' : 'top-2 right-2'} bg-black/60 backdrop-blur-sm text-white text-[9px] font-bold px-1.5 py-0.5 rounded-md capitalize leading-none shadow-sm`}>
            {product.condition === 'new' ? 'New' : 'Used'}
          </div>

          {/* Expiry warning — bottom right, only when ≤3 days left */}
          {daysLeft !== null && daysLeft <= 3 && (
            <div className="absolute bottom-7 right-2 bg-red-600/90 text-white text-[8px] font-black px-1.5 py-0.5 rounded leading-none shadow-sm">
              {daysLeft === 0 ? 'Expires today' : `${daysLeft}d left`}
            </div>
          )}

          {/* Image dots — only if multiple */}
          {total > 1 && (
            <div className="absolute bottom-7 left-1/2 -translate-x-1/2 flex gap-1">
              {Array.from({ length: total }).map((_, i) => (
                <button key={i}
                  onClick={e => { e.stopPropagation(); setImgIdx(i); }}
                  className={`h-1 rounded-full transition-all ${i === imgIdx ? 'bg-white w-3.5' : 'bg-white/50 w-1'}`} />
              ))}
            </div>
          )}
        </div>

        {/* Body — smartphone readable typography */}
        <div className="p-2.5 flex flex-col gap-1">
          <p className="text-xs sm:text-sm font-semibold leading-snug line-clamp-2 text-zinc-900 dark:text-zinc-100 group-hover:text-orange-600 transition-colors">
            {product.title}
          </p>
          {product.location && (
            <p className="text-[11px] sm:text-xs text-zinc-500 dark:text-zinc-400 flex items-center gap-1 leading-none truncate font-medium">
              <MapPin size={10} className="shrink-0 text-orange-500" />
              <span>{product.location.woreda ?? product.location.subcity ?? product.location.region}</span>
            </p>
          )}
        </div>
      </div>

      {detailOpen && (
        <DetailDrawer product={product} onClose={() => setDetailOpen(false)}
          onBuy={() => { setDetailOpen(false); setBuyOpen(true); }} />
      )}
      {buyOpen && <BuyModal product={product} onClose={() => setBuyOpen(false)} />}
    </>
  );
}
