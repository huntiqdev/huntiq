'use client';
import { createContext, useContext, useState, useEffect, useRef, useCallback, useMemo, ReactNode } from 'react';
import { API_URL, api } from '@/lib/api';

export interface AuthUser {
  id: number;
  name?: string;
  phone?: string;
  role: string;
  adminRole: string; // normal | super (admins only)
  buyerTier: string;
  sellerPreference?: string; // sale | rent | both (sellers only)
  capabilities: string[];
  token: string;
}

interface AuthContextValue {
  user: AuthUser | null;
  loading: boolean;
  login: (token: string, sessionToken?: string, profile?: { name?: string; phone?: string; [key: string]: any }) => Promise<void>;
  logout: () => Promise<void>;
  isAdmin: boolean;
  isSuperAdmin: boolean;
  isPlatformOwner: boolean;
  isPlatformAdmin: boolean;
  isSecurityAdmin: boolean;
  isOperationsAdmin: boolean;
  isModerator: boolean;
  isSupportAgent: boolean;
  adminRole: string;
  isSeller: boolean;
  isBuyer: boolean;
  isLandlord: boolean;
  isRenter: boolean;
  capabilities: string[];
  hasCapability: (cap: string) => boolean;
  activateCapability: (capability: string, businessName?: string, preference?: string) => Promise<{ success: boolean; error?: string }>;
  buyerTier: string;
  sellerPreference: string;
  canViewGoldImages: boolean;
  canViewPremiumImages: boolean;
  adminPermissions: string[];
  hasAdminPermission: (perm: string) => boolean;
  refreshAdminPermissions: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue>({
  user: null, loading: true, login: async () => {}, logout: async () => {},
  isAdmin: false, isSuperAdmin: false,
  isPlatformOwner: false, isPlatformAdmin: false, isSecurityAdmin: false,
  isOperationsAdmin: false, isModerator: false, isSupportAgent: false,
  adminRole: '',
  isSeller: false, isBuyer: false, isLandlord: false, isRenter: false,
  capabilities: [], hasCapability: () => false, activateCapability: async () => ({ success: false }),
  buyerTier: 'normal', sellerPreference: 'sale', canViewGoldImages: false, canViewPremiumImages: false,
  adminPermissions: [], hasAdminPermission: () => false, refreshAdminPermissions: async () => {},
});

function parseJWT(token: string): {
  user_id: number;
  name?: string;
  phone?: string;
  role: string;
  admin_role?: string;
  buyer_tier?: string;
  seller_preference?: string;
  capabilities?: string[];
  exp?: number;
} | null {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) return null;
    return JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')));
  }
  catch { return null; }
}

const STORAGE_KEY = 'huntiq_auth_user';

function userFromToken(token: string, fallbackProfile?: Partial<AuthUser>): AuthUser | null {
  const p = parseJWT(token);
  if (!p) return null;
  if (p.exp && p.exp * 1000 < Date.now() - 10_000) return null;

  // Derive additive capabilities
  let caps: string[] = Array.isArray(p.capabilities) && p.capabilities.length > 0 ? p.capabilities : [];
  if (caps.length === 0) {
    caps = ['buyer'];
    if (p.role === 'seller') caps.push('seller');
    if (p.seller_preference === 'rent' || p.seller_preference === 'both') caps.push('landlord');
  }

  return {
    id: p.user_id,
    name: p.name || fallbackProfile?.name || '',
    phone: p.phone || fallbackProfile?.phone || '',
    role: p.role,
    adminRole: p.admin_role ?? '',
    buyerTier: p.buyer_tier ?? 'normal',
    sellerPreference: p.seller_preference,
    capabilities: caps,
    token,
  };
}

function getCachedAuthUser(): AuthUser | null {
  if (typeof window === 'undefined') return null;
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw);
    if (!parsed || !parsed.token) return null;
    const u = userFromToken(parsed.token, parsed);
    if (!u) {
      localStorage.removeItem(STORAGE_KEY);
      return null;
    }
    return u;
  } catch {
    return null;
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(() => getCachedAuthUser());
  const [loading, setLoading] = useState(() => (typeof window !== 'undefined' ? !getCachedAuthUser() : true));
  const [adminPermissions, setAdminPermissions] = useState<string[]>([]);
  const [isSuperAdmin, setIsSuperAdmin] = useState(false);

  const fetchPermissions = useCallback(async (u: AuthUser): Promise<void> => {
    if (u.role !== 'admin') {
      setAdminPermissions([]);
      setIsSuperAdmin(false);
      return;
    }
    const parsed = parseJWT(u.token);
    if (parsed?.exp && parsed.exp * 1000 < Date.now()) {
      await fetch('/api/auth/session', { method: 'DELETE' });
      setUser(null);
      try { localStorage.removeItem(STORAGE_KEY); } catch {}
      setAdminPermissions([]);
      setIsSuperAdmin(false);
      return;
    }
    try {
      const res = await api.adminGetMyPermissions(u.token);
      setAdminPermissions(res?.permissions ?? []);
      setIsSuperAdmin(!!res?.is_super);
    } catch (err: any) {
      if (err?.status === 401 || err?.status === 403) {
        await fetch('/api/auth/session', { method: 'DELETE' });
        setUser(null);
        try { localStorage.removeItem(STORAGE_KEY); } catch {}
        setAdminPermissions([]);
        setIsSuperAdmin(false);
      }
      // On network errors, keep cached permissions
    }
  }, []);

  useEffect(() => {
    fetch('/api/auth/session')
      .then(r => r.json())
      .then(async ({ token }) => {
        const cached = getCachedAuthUser();
        const activeToken = token || cached?.token;
        if (activeToken) {
          const u = userFromToken(activeToken, cached || undefined);
          if (u) {
            // Restore session cookie if cookie was missing but localStorage had valid credentials
            if (!token && cached?.token) {
              fetch('/api/auth/session', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ token: cached.token }),
              }).catch(() => {});
            }

            // Enrich name and phone if missing
            if (!u.name || !u.phone) {
              try {
                const res = await fetch(`${API_URL}/auth/me`, {
                  headers: { Authorization: `Bearer ${activeToken}` },
                });
                if (res.ok) {
                  const me = await res.json();
                  if (me?.name) u.name = me.name;
                  if (me?.phone) u.phone = me.phone;
                }
              } catch {}
            }
            setUser(u);
            try {
              localStorage.setItem(STORAGE_KEY, JSON.stringify(u));
            } catch {}
            await fetchPermissions(u);
            return;
          }
        }

        // Neither cookie nor cached storage had a valid, unexpired token
        setUser(null);
        try {
          localStorage.removeItem(STORAGE_KEY);
        } catch {}
      })
      .catch(() => {
        // In offline or intermittent network conditions, preserve valid cached user
        const cached = getCachedAuthUser();
        if (cached) {
          setUser(cached);
        }
      })
      .finally(() => setLoading(false));
  }, [fetchPermissions]);

  const userRef = useRef<AuthUser | null>(null);

  useEffect(() => {
    userRef.current = user;
  }, [user]);

  useEffect(() => {
    if (!user || user.role !== 'admin') return;
    let cancelled = false;
    const poll = () => {
      if (!cancelled && userRef.current) fetchPermissions(userRef.current);
    };
    const interval = setInterval(poll, 30_000);
    const onFocus = () => {
      if (!cancelled && userRef.current) fetchPermissions(userRef.current);
    };
    window.addEventListener('focus', onFocus);
    return () => {
      cancelled = true;
      clearInterval(interval);
      window.removeEventListener('focus', onFocus);
    };
  }, [user, fetchPermissions]);

  const login = async (token: string, sessionToken?: string, profile?: { name?: string; phone?: string; [key: string]: any }) => {
    await fetch('/api/auth/session', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ token, session_token: sessionToken }),
    }).catch(() => {});
    const u = userFromToken(token, profile);
    if (u) {
      if (profile?.name) u.name = profile.name;
      if (profile?.phone) u.phone = profile.phone;
      setUser(u);
      try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(u));
      } catch {}
      await fetchPermissions(u);
    }
  };

  const logout = async () => {
    await fetch('/api/auth/session', { method: 'DELETE' });
    setUser(null);
    setAdminPermissions([]);
    setIsSuperAdmin(false);
    try {
      localStorage.removeItem(STORAGE_KEY);
    } catch {}
  };

  const refreshAdminPermissions = async () => {
    if (user) await fetchPermissions(user);
  };

  const buyerTier = user?.buyerTier ?? 'normal';
  const sellerPreference = user?.sellerPreference ?? 'sale';
  const capabilities = useMemo(() => user?.capabilities ?? [], [user?.capabilities]);
  const adminRole = user?.adminRole ?? '';

  const isPlatformOwner = adminRole === 'platform_owner' || adminRole === 'super';
  const isPlatformAdmin = adminRole === 'platform_admin';
  const isSecurityAdmin = adminRole === 'security_admin';
  const isOperationsAdmin = adminRole === 'operations_admin';
  const isModerator = adminRole === 'moderator';
  const isSupportAgent = adminRole === 'support_agent';

  const hasCapability = useCallback((cap: string): boolean => {
    if (!user) return false;
    if (capabilities.includes(cap)) return true;
    if (cap === 'buyer') return true; // all authenticated marketplace users have buyer capabilities
    if (cap === 'seller' && user.role === 'seller') return true;
    if (cap === 'landlord' && (sellerPreference === 'rent' || sellerPreference === 'both')) return true;
    return false;
  }, [user, capabilities, sellerPreference]);

  const isSeller = hasCapability('seller');
  const isBuyer = hasCapability('buyer');
  const isLandlord = hasCapability('landlord');
  const isRenter = hasCapability('renter') || isBuyer;

  const activateCapability = async (capability: string, businessName?: string, preference?: string) => {
    if (!user) return { success: false, error: 'Not authenticated' };
    try {
      const data = await api.activateCapability(capability, businessName, preference);
      if (data?.token) {
        await login(data.token);
      }
      return { success: true };
    } catch (err: any) {
      return { success: false, error: err.message || 'Network error' };
    }
  };

  const hasAdminPermission = (perm: string) => {
    if (isPlatformOwner || isSuperAdmin) return true;
    return adminPermissions.includes(perm);
  };

  return (
    <AuthContext.Provider value={{
      user, loading, login, logout,
      isAdmin:  user?.role === 'admin',
      isSuperAdmin: isPlatformOwner || isSuperAdmin,
      isPlatformOwner,
      isPlatformAdmin,
      isSecurityAdmin,
      isOperationsAdmin,
      isModerator,
      isSupportAgent,
      adminRole,
      isSeller,
      isBuyer,
      isLandlord,
      isRenter,
      capabilities,
      hasCapability,
      activateCapability,
      buyerTier,
      sellerPreference,
      canViewGoldImages:    buyerTier === 'gold' || buyerTier === 'premium',
      canViewPremiumImages: buyerTier === 'premium',
      adminPermissions,
      hasAdminPermission,
      refreshAdminPermissions,
    }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
