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
  avatarURL?: string;
  token: string;
}

interface AuthContextValue {
  user: AuthUser | null;
  loading: boolean;
  login: (token: string, sessionToken?: string, profile?: { name?: string; phone?: string; [key: string]: any }) => Promise<void>;
  logout: () => Promise<void>;
  updateUserProfile: (data: Partial<AuthUser>) => void;
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
  user: null, loading: true, login: async () => {}, logout: async () => {}, updateUserProfile: () => {},
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

function userFromToken(token: string, fallbackProfile?: Partial<AuthUser>): AuthUser | null {
  const p = parseJWT(token);
  if (!p) return null;
  if (p.exp && p.exp * 1000 < Date.now() - 10_000) return null;

  // Derive additive capabilities: Seller and Landlord are unified on the same merchant account
  let caps: string[] = Array.isArray(p.capabilities) && p.capabilities.length > 0 ? [...p.capabilities] : [];
  if (caps.length === 0) {
    caps = ['buyer', 'renter'];
    if (p.role === 'seller') {
      caps.push('seller', 'landlord');
    }
  } else {
    if (!caps.includes('renter')) caps.push('renter');
    if (caps.includes('seller') && !caps.includes('landlord')) caps.push('landlord');
    if (caps.includes('landlord') && !caps.includes('seller')) caps.push('seller');
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
    avatarURL: (p as any).avatar_url || (fallbackProfile as any)?.avatar_url || fallbackProfile?.avatarURL,
    token,
  };
}

const TOKEN_STORAGE_KEY = 'huntiq_token';
const USER_STORAGE_KEY = 'huntiq_user';
const LOGGED_OUT_KEY = 'huntiq_logged_out';

function getSavedToken(): string | null {
  if (typeof window === 'undefined') return null;
  try {
    if (localStorage.getItem(LOGGED_OUT_KEY) === 'true') return null;
    return localStorage.getItem(TOKEN_STORAGE_KEY);
  } catch {
    return null;
  }
}

function getSavedUser(): Partial<AuthUser> | null {
  if (typeof window === 'undefined') return null;
  try {
    if (localStorage.getItem(LOGGED_OUT_KEY) === 'true') return null;
    const str = localStorage.getItem(USER_STORAGE_KEY);
    return str ? JSON.parse(str) : null;
  } catch {
    return null;
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  // Synchronous initial state from safe browser storage to avoid flicker & unwanted redirects
  const [user, setUser] = useState<AuthUser | null>(() => {
    const token = getSavedToken();
    if (!token) return null;
    const parsed = parseJWT(token);
    if (!parsed?.exp || parsed.exp * 1000 <= Date.now()) return null;
    const cachedProfile = getSavedUser();
    return userFromToken(token, cachedProfile || undefined);
  });

  const [loading, setLoading] = useState<boolean>(() => {
    const token = getSavedToken();
    if (!token) return true;
    const parsed = parseJWT(token);
    return !(parsed?.exp && parsed.exp * 1000 > Date.now());
  });

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
      if (typeof window !== 'undefined') {
        try {
          localStorage.removeItem(TOKEN_STORAGE_KEY);
          localStorage.removeItem(USER_STORAGE_KEY);
        } catch {}
      }
      await fetch('/api/auth/session', { method: 'DELETE' });
      setUser(null);
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
        if (typeof window !== 'undefined') {
          try {
            localStorage.removeItem(TOKEN_STORAGE_KEY);
            localStorage.removeItem(USER_STORAGE_KEY);
          } catch {}
        }
        await fetch('/api/auth/session', { method: 'DELETE' });
        setUser(null);
        setAdminPermissions([]);
        setIsSuperAdmin(false);
      }
      // On network errors, keep cached permissions
    }
  }, []);

  useEffect(() => {
    let isCancelled = false;

    const initAuth = async () => {
      // 1. If explicitly logged out, keep logged out
      const explicitlyLoggedOut = typeof window !== 'undefined' && localStorage.getItem(LOGGED_OUT_KEY) === 'true';
      if (explicitlyLoggedOut) {
        setUser(null);
        setLoading(false);
        return;
      }

      // 2. Read cached browser token
      const cachedToken = getSavedToken();
      const cachedProfile = getSavedUser();
      let activeUser: AuthUser | null = null;
      if (cachedToken) {
        activeUser = userFromToken(cachedToken, cachedProfile || undefined);
        if (activeUser && !isCancelled) {
          setUser(activeUser);
          setLoading(false); // Unblock immediately for smooth navigation
        }
      }

      // 3. Hydrate & verify with server session endpoint (/api/auth/session)
      try {
        const sessionRes = await fetch('/api/auth/session').catch(() => null);
        const sessionData = sessionRes && sessionRes.ok ? await sessionRes.json().catch(() => null) : null;
        if (isCancelled) return;

        const serverToken = sessionData?.token;
        const effectiveToken = serverToken || cachedToken;

        if (effectiveToken) {
          const u = userFromToken(effectiveToken, cachedProfile || undefined);
          if (u) {
            // Enrich with /auth/me to verify token validity and get freshest profile
            try {
              const res = await fetch(`${API_URL}/auth/me`, {
                credentials: 'include',
                headers: { Authorization: `Bearer ${effectiveToken}` },
              });
              if (res.ok) {
                const me = await res.json();
                if (me?.name) u.name = me.name;
                if (me?.phone) u.phone = me.phone;
                if (typeof window !== 'undefined') {
                  try {
                    localStorage.removeItem(LOGGED_OUT_KEY);
                    localStorage.setItem(TOKEN_STORAGE_KEY, effectiveToken);
                    localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(u));
                  } catch {}
                }
                if (!isCancelled) {
                  setUser(u);
                  await fetchPermissions(u);
                }
                return;
              } else if (res.status === 401 || res.status === 403) {
                // Token revoked on server
                if (typeof window !== 'undefined') {
                  try {
                    localStorage.removeItem(TOKEN_STORAGE_KEY);
                    localStorage.removeItem(USER_STORAGE_KEY);
                  } catch {}
                }
                await fetch('/api/auth/session', { method: 'DELETE' }).catch(() => {});
                if (!isCancelled) setUser(null);
                return;
              }
            } catch {
              // Network offline: retain activeUser from browser cache
              if (activeUser && !isCancelled) {
                setUser(activeUser);
                return;
              }
            }
          }
        }

        // If no server token and no cached token
        if (!activeUser && !isCancelled) {
          setUser(null);
        }
      } catch {
        if (!activeUser && !isCancelled) {
          setUser(null);
        }
      } finally {
        if (!isCancelled) setLoading(false);
      }
    };

    initAuth();

    return () => {
      isCancelled = true;
    };
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
    // 1. Persist to browser storage (cache setup)
    if (typeof window !== 'undefined') {
      try {
        localStorage.removeItem(LOGGED_OUT_KEY);
        localStorage.setItem(TOKEN_STORAGE_KEY, token);
      } catch {}
    }

    // 2. Persist to HttpOnly cookie session endpoint
    await fetch('/api/auth/session', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ token, session_token: sessionToken }),
    }).catch(() => {});

    // 3. Construct user model
    const u = userFromToken(token, profile);
    if (u) {
      if (profile?.name) u.name = profile.name;
      if (profile?.phone) u.phone = profile.phone;
      if (profile?.avatar_url || profile?.avatarURL) u.avatarURL = profile.avatar_url || profile.avatarURL;
      if (typeof window !== 'undefined') {
        try {
          localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(u));
        } catch {}
      }
      setUser(u);
      setLoading(false);
      await fetchPermissions(u);
    }
  };

  const updateUserProfile = useCallback((data: Partial<AuthUser>) => {
    setUser(prev => {
      if (!prev) return null;
      const updated = { ...prev, ...data };
      if (typeof window !== 'undefined') {
        try {
          localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(updated));
        } catch {}
      }
      return updated;
    });
  }, []);

  const logout = async () => {
    // 1. Mark explicit logout & purge browser cache
    if (typeof window !== 'undefined') {
      try {
        localStorage.setItem(LOGGED_OUT_KEY, 'true');
        localStorage.removeItem(TOKEN_STORAGE_KEY);
        localStorage.removeItem(USER_STORAGE_KEY);
        sessionStorage.clear();
      } catch {}
    }

    // 2. Invalidate backend session
    try {
      const csrf = typeof document !== 'undefined' ? (document.cookie.match(/(?:^|;\s*)huntiq_csrf=([^;]+)/)?.[1] || '') : '';
      await fetch(`${API_URL}/auth/logout`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          ...(csrf ? { 'X-CSRF-Token': decodeURIComponent(csrf) } : {}),
        },
      }).catch(() => {});
    } catch {}

    // 3. Clear cookie session
    await fetch('/api/auth/session', { method: 'DELETE' }).catch(() => {});

    // 4. Reset React state
    setUser(null);
    setAdminPermissions([]);
    setIsSuperAdmin(false);
  };

  const refreshAdminPermissions = async () => {
    if (user) await fetchPermissions(user);
  };

  const buyerTier = user?.buyerTier ?? 'normal';
  const sellerPreference = user?.sellerPreference ?? 'sale';
  const capabilities = useMemo(() => user?.capabilities ?? [], [user?.capabilities]);
  const adminRole = user?.adminRole ?? '';

  const isPlatformOwner = adminRole === 'platform_owner' || adminRole === 'super';
  const isPlatformAdmin = adminRole === 'platform_admin' || isPlatformOwner;
  const isSecurityAdmin = adminRole === 'security_admin' || adminRole === 'ciso' || adminRole === 'soc_analyst' || adminRole === 'fraud_investigator' || adminRole === 'compliance_auditor' || isPlatformOwner;
  const isOperationsAdmin = adminRole === 'operations_admin' || adminRole === 'kyc_verification_officer' || adminRole === 'product_inspection_officer' || adminRole === 'dispute_refund_officer' || adminRole === 'appeals_officer' || isPlatformOwner;
  const isModerator = adminRole === 'moderator' || adminRole === 'content_moderator' || isPlatformOwner;
  const isSupportAgent = adminRole === 'support_agent' || adminRole === 'customer_support_lead' || isPlatformOwner;

  const hasCapability = useCallback((cap: string): boolean => {
    if (!user) return false;
    if (capabilities.includes(cap)) return true;
    if (cap === 'buyer' || cap === 'renter') return true; // all authenticated marketplace users can buy and rent
    if ((cap === 'seller' || cap === 'landlord') && (user.role === 'seller' || user.role === 'admin')) return true;
    return false;
  }, [user, capabilities]);

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
      user, loading, login, logout, updateUserProfile,
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
