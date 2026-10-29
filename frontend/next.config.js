/** @type {import('next').NextConfig} */
// BACKEND_URL is only used for local dev proxy (next.config.js rewrites).
// In production on Vercel, rewrites are handled by vercel.json.
const backendURL = process.env.BACKEND_URL || 'http://localhost:8080';

// NEXT_PUBLIC_BACKEND_HOST is the hostname for remote image patterns (e.g. your-api.onrender.com)
const backendHost = process.env.NEXT_PUBLIC_BACKEND_HOST || '';

const remotePatterns = [
  { protocol: 'http', hostname: 'localhost', port: '8080', pathname: '/uploads/**' },
  { protocol: 'https', hostname: 'images.unsplash.com' },
  { protocol: 'https', hostname: 'plus.unsplash.com' },
  { protocol: 'https', hostname: '**.unsplash.com' },
];
if (backendHost) {
  remotePatterns.push({ protocol: 'https', hostname: backendHost, pathname: '/uploads/**' });
}

const nextConfig = {
  async rewrites() {
    if (process.env.NODE_ENV === 'production') return [];
    return [
      { source: '/api/:path*',     destination: `${backendURL}/api/:path*` },
      { source: '/uploads/:path*', destination: `${backendURL}/uploads/:path*` },
    ];
  },
  images: {
    remotePatterns,
  },
};

module.exports = nextConfig;
