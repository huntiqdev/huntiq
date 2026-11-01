import { notFound } from 'next/navigation';
import { api } from '@/lib/api';
import { Product } from '@/types';
import { ImageGallery } from '@/components/public/ImageGallery';
import { BuyButton } from '@/components/public/BuyButton';
import Link from 'next/link';
import { MapPin, ShieldCheck, Navigation, ShieldAlert, Sparkles } from 'lucide-react';
import { ClientLabel } from '@/components/public/ClientLabel';
import { ClientProductBadges } from '@/components/public/ClientProductBadges';
import CategoryBreadcrumbs from '@/components/category/CategoryBreadcrumbs';
import { ProductAttributes } from '@/components/public/ProductAttributes';

export const dynamic = 'force-dynamic';

interface Props { params: Promise<{ id: string }> }

export default async function ProductDetailPage({ params }: Props) {
  const { id } = await params;
  let product: Product;
  try { product = await api.getProduct(Number(id)); }
  catch (err: any) {
    // Only treat a real 404 as not-found; let other errors surface
    if (err?.status === 404) notFound();
    throw err;
  }

  const images = product.images ?? [];
  const lat = product.latitude;
  const lng = product.longitude;
  // Parse as floats and validate range to prevent URL injection from seller-supplied coords
  const safeLat = lat != null ? parseFloat(String(lat)) : null;
  const safeLng = lng != null ? parseFloat(String(lng)) : null;
  const hasCoords =
    safeLat != null && safeLng != null &&
    isFinite(safeLat) && isFinite(safeLng) &&
    safeLat >= -90 && safeLat <= 90 &&
    safeLng >= -180 && safeLng <= 180;
  const mapsUrl = hasCoords
    ? `https://www.google.com/maps?q=${safeLat},${safeLng}`
    : null;
  // Embed URL for the iframe (no API key needed for basic embed)
  const embedUrl = hasCoords
    ? `https://maps.google.com/maps?q=${safeLat},${safeLng}&z=16&output=embed`
    : null;

  const condition = product.condition ?? '';

  return (
    <div className="max-w-5xl mx-auto px-4 pt-6 pb-24 md:py-8">
      {/* Breadcrumbs */}
      {(product as any).category_hierarchy_id && (
        <div className="mb-4">
          <CategoryBreadcrumbs categoryId={(product as any).category_hierarchy_id} />
        </div>
      )}

      <div className="grid md:grid-cols-2 gap-6 md:gap-8">
        {/* Gallery */}
        <div className="flex flex-col gap-4">
          <ImageGallery images={images} sellerTier={product.seller_tier ?? 'normal'} />

          {/* Google Maps embed — only for Gold/Premium sellers with coords */}
          {embedUrl && (
            <div className="rounded-2xl overflow-hidden border border-zinc-200 dark:border-zinc-700 shadow-sm">
              <div className="flex items-center justify-between px-3 py-2 bg-zinc-50 dark:bg-zinc-800 border-b border-zinc-200 dark:border-zinc-700">
                <div className="flex items-center gap-1.5 text-xs font-bold text-zinc-600 dark:text-zinc-300">
                  <MapPin size={12} className="text-orange-500" />
                  <ClientLabel k="ui.exact_location" />
                </div>
                <a
                  href={mapsUrl!}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-1 text-[10px] font-black text-orange-600 hover:text-orange-500 transition-colors"
                >
                  <Navigation size={10} /> <ClientLabel k="ui.open_in_maps" />
                </a>
              </div>
              <iframe
                src={embedUrl}
                width="100%"
                height="220"
                style={{ border: 0 }}
                allowFullScreen
                loading="lazy"
                referrerPolicy="no-referrer-when-downgrade"
                title="Product location map"
              />
            </div>
          )}
        </div>

        {/* Info */}
        <div className="flex flex-col gap-4">
          {/* Badges */}
          <div className="flex flex-wrap items-center gap-2">
            {product.category && (
              <span className="text-[10px] font-black uppercase tracking-widest text-orange-600 bg-orange-600/10 px-2 py-1 rounded-lg">
                {product.category.name}
              </span>
            )}
            <ClientProductBadges
              categorySlug={product.category?.slug ?? ''}
              subcategory={product.subcategory}
              item={product.item}
              condition={condition}
            />
            <span className="text-[10px] font-black uppercase tracking-widest text-blue-600 bg-blue-600/10 px-2 py-1 rounded-lg flex items-center gap-1">
              <ShieldCheck size={11} /><ClientLabel k="ui.verified_listing" />
            </span>
          </div>

          <h1 className="text-xl md:text-2xl font-black leading-tight">{product.title}</h1>

          <div className={`text-2xl md:text-3xl font-black ${product.listing_type === 'rent' ? 'text-teal-600' : 'text-orange-600'}`}>
            ETB {product.listing_type === 'rent' && product.rental_price != null
              ? `${product.rental_price.toLocaleString()} / ${product.rental_duration_value || 1} ${product.rental_duration_unit || 'month'}` 
              : (product.price != null ? product.price.toLocaleString() : '—')}
          </div>

          {product.location && (
            <div className="flex flex-col gap-1">
              <div className="flex items-center gap-1.5 text-sm opacity-60">
                <MapPin size={14} />
                {[product.location.woreda, product.location.subcity, product.location.zone, product.location.region]
                  .filter(Boolean).join(', ')}
              </div>
              {product.kebele && (
                <p className="text-xs opacity-50 pl-5">Kebele: {product.kebele}</p>
              )}
              {product.specific_area && (
                <p className="text-xs opacity-50 pl-5">{product.specific_area}</p>
              )}
            </div>
          )}

          {product.description && (
            <div>
              <p className="text-xs font-black uppercase tracking-wider opacity-50 mb-1.5">
                <ClientLabel k="ui.description" />
              </p>
              <p className="text-sm leading-relaxed opacity-80">{product.description}</p>
            </div>
          )}

          {/* Dynamic product attributes (specs) with server-enforced tier gating */}
          <ProductAttributes
            productId={product.id}
            requiredBuyerTier={product.required_buyer_tier}
            specificationsLocked={product.specifications_locked}
          />

          <BuyButton product={product} />

          {/* Safe Exchange Shield */}
          <div className="mt-2 p-3.5 rounded-xl border border-amber-500/20 bg-amber-500/5 text-amber-900 dark:text-amber-200">
            <div className="flex items-center gap-2 mb-1.5 font-bold text-xs">
              <ShieldAlert size={14} className="text-amber-600 dark:text-amber-400 shrink-0" />
              <span><ClientLabel k="ui.safe_exchange_title" /></span>
            </div>
            <p className="text-[11px] leading-relaxed text-zinc-600 dark:text-zinc-400">
              <ClientLabel k="ui.safe_exchange_desc" />
            </p>
          </div>

          {/* AI Related Search Discovery */}
          {product.category && (
            <div className="mt-1 flex flex-wrap items-center gap-2">
              <span className="text-[11px] font-semibold text-zinc-500 flex items-center gap-1">
                <Sparkles size={11} className="text-orange-500" />
                <ClientLabel k="ui.related_ai_search" />:
              </span>
              <Link
                href={`/products?q=${encodeURIComponent(product.category.name)}&mode=ai`}
                className="text-[11px] font-medium bg-zinc-100 hover:bg-orange-50 dark:bg-zinc-800 dark:hover:bg-orange-950/40 text-zinc-700 dark:text-zinc-300 hover:text-orange-600 px-2.5 py-1 rounded-full border border-zinc-200 dark:border-zinc-700 transition-colors"
              >
                <ClientLabel k="ui.more_listings_in" /> {product.category.name} →
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
