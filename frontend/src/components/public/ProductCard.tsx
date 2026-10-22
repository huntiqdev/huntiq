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

// ── Detail Drawer ─────────────────────────────────────────────────────────────
function DetailDrawer({ product, onClose, onBuy }: {
  product: Product; onClose: () => void; onBuy: () => void;
}) {
  const { t } = useLocale();
  const images = product.images ?? [];
  const [imgIdx, setImgIdx] = useState(0);
  const total = images.length;
  const blurred = useImageBlur(product);
  const sellerTier = product.seller_tier;

  return (
    <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4">
      <div className="absolute inset-0 bg-black/70 backdrop-blur-sm" onClick={onClose} />
      <div className="relative z-10 w-full sm:max-w-lg bg-white dark:bg-zinc-900 rounded-t-3xl sm:rounded-3xl overflow-hidden shadow-2xl max-h-[92vh] flex flex-col">

        {/* Image carousel */}
        <div className="relative aspect-[4/3] bg-zinc-100 dark:bg-zinc-800 shrink-0">
          {images.length > 0 ? (
          <Image src={resolveImageUrl(images[imgIdx]?.image_url)}
              alt={product.title} fill
              unoptimized={isUnoptimized(resolveImageUrl(images[imgIdx]?.image_url))}
              className={`object-cover ${blurred ? 'blur-2xl scale-110' : ''}`}
              sizes="(max-width: 640px) 100vw, 512px" />
          ) : (
            <div className="absolute inset-0 flex items-center justify-center opacity-20">
              <ShoppingBag size={48} />
            </div>
          )}
          {blurred && <BlurOverlay tier={sellerTier ?? 'gold'} />}

          {total > 1 && !blurred && (
            <>
              <button onClick={e => { e.stopPropagation(); setImgIdx(i => (i - 1 + total) % total); }}
                className="absolute left-2 top-1/2 -translate-y-1/2 w-8 h-8 rounded-full bg-black/40 text-white flex items-center justify-center hover:bg-black/60 transition-colors">
                <ChevronLeft size={18} />
              </button>
              <button onClick={e => { e.stopPropagation(); setImgIdx(i => (i + 1) % total); }}
                className="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 rounded-full bg-black/40 text-white flex items-center justify-center hover:bg-black/60 transition-colors">
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
            className="absolute top-3 right-3 w-8 h-8 rounded-full bg-black/40 text-white flex items-center justify-center hover:bg-black/60 transition-colors">
            <X size={16} />
          </button>
          <div className="absolute top-3 left-3 bg-orange-600 text-white px-3 py-1 rounded-full text-sm font-black shadow-lg">
            ETB {product.price.toLocaleString()}
          </div>
        </div>

        {/* Info */}
        <div className="overflow-y-auto p-5 flex flex-col gap-3">
          <div className="flex flex-wrap gap-1.5">
            {product.category && (
              <span className="text-[10px] font-black uppercase tracking-widest text-orange-600 bg-orange-600/10 px-2 py-1 rounded-lg">
                {product.category.name}
              </span>
            )}
            {product.subcategory && (
              <span className="text-[10px] font-black uppercase tracking-widest text-purple-600 bg-purple-600/10 px-2 py-1 rounded-lg">
                {(t.sub_categories as any)?.[(product.category as any)?.slug ?? '']?.[product.subcategory]
                  ?? (t.items as any)?.[product.subcategory]
                  ?? product.subcategory.replace(/_/g, ' ')}
              </span>
            )}
            <span className="text-[10px] font-black uppercase tracking-widest bg-zinc-100 dark:bg-zinc-800 px-2 py-1 rounded-lg capitalize">
              {(t.ui as any)?.[`condition_${product.condition}`] ?? product.condition.replace(/_/g, ' ')}
            </span>
            <span className="text-[10px] font-black uppercase tracking-widest text-blue-600 bg-blue-600/10 px-2 py-1 rounded-lg flex items-center gap-1">
              <ShieldCheck size={10} /> {t.ui.verified}
            </span>
          </div>

          <h2 className="text-lg font-black leading-tight">{product.title}</h2>

          {product.location && (
            <div className="flex items-center gap-1.5 text-sm opacity-60">
              <MapPin size={13} />
              {[product.location.woreda, product.location.subcity, product.location.region].filter(Boolean).join(', ')}
            </div>
          )}

          {product.description && (
            <p className="text-sm leading-relaxed opacity-70">{product.description}</p>
          )}

          <div className="flex gap-3 pt-1">
            <button type="button" onClick={onBuy}
              className={`flex-1 text-white rounded-2xl py-3 font-bold text-sm flex items-center justify-center gap-2 transition-colors active:scale-95 ${
                product.listing_type === 'rent'
                  ? 'bg-teal-600 hover:bg-teal-700'
                  : 'bg-orange-600 hover:bg-orange-700'
              }`}>
              {product.listing_type === 'rent' ? <ShoppingBag size={16} /> : <ShoppingBag size={16} />}
              {product.listing_type === 'rent' ? (t.ui.rent_now ?? 'Rent / Contact Landlord') : t.ui.buy_now}
            </button>
            <Link href={`/products/${product.id}`}
              className="px-4 py-3 rounded-2xl border border-zinc-200 dark:border-zinc-700 font-bold text-sm hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors">
              {t.ui.details}
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

// ── Card ──────────────────────────────────────────────────────────────────────
// Compact design: fits 3 cols on mobile (≥360px), 4 on sm, 5 on md, 6 on lg+
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
  const blurred = useImageBlur(product);
  const sellerTier = product.seller_tier;
  const tierBadge = sellerTier ? TIER_BADGE[sellerTier] : null;
  const daysLeft = expiryDaysLeft(product.expires_at);

  const images = product.images ?? [];
  const total = Math.min(images.length, 3);
  const thumb = resolveImageUrl(images[imgIdx]?.image_url);

  return (
    <>
      {/* Card — square-ish, very compact */}
      <div
        className="group relative bg-white dark:bg-zinc-900 rounded-xl border border-zinc-100 dark:border-zinc-800 overflow-hidden flex flex-col cursor-pointer transition-all duration-200 hover:shadow-lg hover:shadow-black/10 hover:-translate-y-0.5 hover:border-zinc-200 dark:hover:border-zinc-700 active:scale-[0.98]"
        onClick={() => setDetailOpen(true)}
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

          {/* Price — bottom left */}
          <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 to-transparent px-1.5 pb-1 pt-4">
            <p className="text-white text-[10px] font-black leading-none">
              ETB {product.price >= 1000
                ? `${(product.price / 1000).toFixed(product.price % 1000 === 0 ? 0 : 1)}k`
                : product.price.toLocaleString()}
            </p>
          </div>

          {/* Sale Badge — top left */}
          <div className="absolute top-1 left-1 bg-orange-600 text-white text-[8px] font-black px-1.5 py-0.5 rounded-full flex items-center gap-0.5 leading-none shadow-sm">
            <ShoppingBag size={7} />
            SALE
          </div>

          {/* Tier badge — just below sale badge */}
          {tierBadge && (
            <div className={`absolute top-6 left-1 w-4 h-4 rounded-full flex items-center justify-center text-[8px] font-black ${tierBadge.cls}`}>
              {tierBadge.label}
            </div>
          )}

          {/* Condition — top right */}
          <div className="absolute top-1 right-1 bg-black/50 backdrop-blur-sm text-white text-[8px] font-black px-1 py-0.5 rounded capitalize leading-none">
            {product.condition === 'new' ? 'New' : 'Used'}
          </div>

          {/* Expiry warning — bottom right, only when ≤3 days left */}
          {daysLeft !== null && daysLeft <= 3 && (
            <div className="absolute bottom-5 right-1 bg-red-600/90 text-white text-[8px] font-black px-1 py-0.5 rounded leading-none">
              {daysLeft === 0 ? 'Expires today' : `${daysLeft}d left`}
            </div>
          )}

          {/* Image dots — only if multiple */}
          {total > 1 && !blurred && (
            <div className="absolute bottom-5 left-1/2 -translate-x-1/2 flex gap-0.5">
              {Array.from({ length: total }).map((_, i) => (
                <button key={i}
                  onClick={e => { e.stopPropagation(); setImgIdx(i); }}
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

      {detailOpen && (
        <DetailDrawer product={product} onClose={() => setDetailOpen(false)}
          onBuy={() => { setDetailOpen(false); setBuyOpen(true); }} />
      )}
      {buyOpen && <BuyModal product={product} onClose={() => setBuyOpen(false)} />}
    </>
  );
}
