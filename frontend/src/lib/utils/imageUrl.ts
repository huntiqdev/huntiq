/**
 * Resolves a product image URL to an absolute URL.
 * - Absolute URLs (https://... or http://...) returned as-is.
 * - Relative /uploads/ paths prefixed with NEXT_PUBLIC_BACKEND_ORIGIN.
 * - Empty/null returns placeholder SVG.
 *
 * Set NEXT_PUBLIC_BACKEND_ORIGIN to the backend base URL so /uploads/ paths resolve correctly.
 * In production set NEXT_PUBLIC_API_URL or NEXT_PUBLIC_BACKEND_ORIGIN via environment variables.
 */
const BACKEND_ORIGIN = (() => {
  const api = process.env.NEXT_PUBLIC_API_URL ?? '';
  if (api.startsWith('http://') || api.startsWith('https://')) {
    return api.replace(/\/api$/, '');
  }
  return process.env.NEXT_PUBLIC_BACKEND_ORIGIN ?? '';
})();

export function resolveImageUrl(url: string | undefined | null): string {
  if (!url || url.trim() === '') {
    return 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="400" height="300" viewBox="0 0 400 300"%3E%3Crect width="400" height="300" fill="%23f4f4f5"/%3E%3Ctext x="50%25" y="50%25" dominant-baseline="middle" text-anchor="middle" font-family="sans-serif" font-size="14" fill="%23a1a1aa"%3ENo Image%3C/text%3E%3C/svg%3E';
  }
  if (url.startsWith('http://') || url.startsWith('https://')) return url;
  if (url.startsWith('/uploads/')) {
    if (!BACKEND_ORIGIN) return url;
    return `${BACKEND_ORIGIN}${url}`;
  }
  return url;
}

export function isUnoptimized(_url: string): boolean {
  // Always false — Next.js handles all image optimization
  return false;
}
