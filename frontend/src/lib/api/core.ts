export class ApiError extends Error {
  code?: string;
  status?: number;

  constructor(message: string, code?: string, status?: number) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.status = status;
  }
}

export const API_URL = (() => {
  const url = process.env.NEXT_PUBLIC_API_URL || '/api';
  return url.replace(/\/$/, ''); // strip trailing slash
})();

// Guard against multiple concurrent 401 redirects causing a redirect loop.
let lastRedirect = 0;

// Parse a JWT payload without verifying the signature (client-side only).
function parseJWTPayload(token: string): { exp?: number } | null {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) return null;
    return JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')));
  } catch { return null; }
}

function isTokenExpired(token: string): boolean {
  const payload = parseJWTPayload(token);
  if (!payload?.exp) return false;
  // 10s grace period after expiry to account for clock skew.
  // Backend isActiveCache TTL is also set to 10s — keep these in sync.
  return payload.exp * 1000 < Date.now() - 10_000;
}

// Check if online
function isOnline(): boolean {
  return typeof navigator !== 'undefined' ? navigator.onLine : true;
}

// Sleep utility for retry delays
const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

// Retry configuration
const RETRY_CONFIG = {
  maxRetries: 2,
  retryDelay: 1000,
  retryableStatuses: [408, 429, 500, 502, 503, 504],
};

export async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  // Check network before starting
  if (!isOnline()) {
    throw new ApiError('No internet connection', 'NETWORK_OFFLINE', 0);
  }

  const authHeader = (options?.headers as Record<string, string>)?.['Authorization'];
  if (authHeader?.startsWith('Bearer ')) {
    const token = authHeader.slice(7);
    if (isTokenExpired(token) && typeof window !== 'undefined' && Date.now() - lastRedirect > 5000) {
      lastRedirect = Date.now();
      await fetch('/api/auth/session', { method: 'DELETE' }).catch(() => {});
      const path = window.location.pathname;
      if (path.startsWith('/admin')) window.location.href = '/admin/login';
      else if (path.startsWith('/seller')) window.location.href = '/auth?as=seller';
      else if (path.startsWith('/buyer')) window.location.href = '/auth?as=buyer';
      else window.location.href = '/auth';
      throw new ApiError('Session expired — please log in again', 'AUTH_SESSION_EXPIRED', 401);
    }
  }

  let lastError: Error | null = null;
  
  // Retry loop
  for (let attempt = 0; attempt <= RETRY_CONFIG.maxRetries; attempt++) {
    try {
      const res = await fetch(`${API_URL}${endpoint}`, {
        ...options,
        headers: { 'Content-Type': 'application/json', ...options?.headers },
      });
      
      if (!res.ok) {
        // Handle 401 - unauthorized
        if (res.status === 401 && typeof window !== 'undefined' && Date.now() - lastRedirect > 5000) {
          lastRedirect = Date.now();
          await fetch('/api/auth/session', { method: 'DELETE' }).catch(() => {});
          const path = window.location.pathname;
          if (path.startsWith('/admin')) window.location.href = '/admin/login';
          else if (path.startsWith('/seller')) window.location.href = '/auth?as=seller';
          else if (path.startsWith('/buyer')) window.location.href = '/auth?as=buyer';
          else window.location.href = '/auth';
        }
        
        // Handle 503 - maintenance mode
        if (res.status === 503 && typeof window !== 'undefined' && window.location.pathname !== '/maintenance') {
          window.location.href = '/maintenance';
          throw new ApiError('Maintenance mode', 'MAINTENANCE', 503);
        }
        
        // Parse error response
        let body: any = {};
        const text = await res.text().catch(() => '');
        try {
          body = JSON.parse(text);
        } catch {
          body = { error: text ? `API error: ${res.status} — ${text.slice(0, 100).replace(/\s+/g, ' ')}` : `API error: ${res.status}` };
        }
        
        const error = new ApiError(
          body.error || `API error: ${res.status}`,
          body.code,
          res.status
        );
        
        // Retry on retryable status codes
        if (RETRY_CONFIG.retryableStatuses.includes(res.status) && attempt < RETRY_CONFIG.maxRetries) {
          await sleep(RETRY_CONFIG.retryDelay * (attempt + 1));
          continue;
        }
        
        throw error;
      }
      
      // Success - parse and return
      const json = await res.json();
      const result = ('data' in json) ? json.data : json;
      return (result ?? null) as T;
      
    } catch (error: any) {
      lastError = error;
      
      // Network errors - retry if online
      if (error instanceof TypeError && error.message.includes('fetch')) {
        if (!isOnline()) {
          throw new ApiError('Lost internet connection', 'NETWORK_OFFLINE', 0);
        }
        if (attempt < RETRY_CONFIG.maxRetries) {
          await sleep(RETRY_CONFIG.retryDelay * (attempt + 1));
          continue;
        }
        throw new ApiError('Network request failed', 'NETWORK_ERROR', 0);
      }
      
      // ApiError - don't retry unless it's a retryable status
      if (error instanceof ApiError) {
        if (error.status && RETRY_CONFIG.retryableStatuses.includes(error.status) && attempt < RETRY_CONFIG.maxRetries) {
          await sleep(RETRY_CONFIG.retryDelay * (attempt + 1));
          continue;
        }
        throw error;
      }
      
      // Other errors - retry
      if (attempt < RETRY_CONFIG.maxRetries) {
        await sleep(RETRY_CONFIG.retryDelay * (attempt + 1));
        continue;
      }
      
      throw error;
    }
  }
  
  // If we exhausted all retries
  throw lastError || new ApiError('Request failed after retries', 'REQUEST_FAILED', 0);
}

export function authHeaders(token: string) {
  return { Authorization: `Bearer ${token}` };
}
