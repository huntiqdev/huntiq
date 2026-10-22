'use client';
import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { api } from '@/lib/api';

/**
 * Platform settings loaded from backend /public/payment-settings endpoint.
 * These settings are admin-configurable and control business logic across the app.
 */
export interface PlatformSettings {
  // Fee rates
  broker_fee_rate?: string;
  referral_commission_rate?: string;
  normal_startup_fee?: string;
  gold_startup_fee?: string;
  premium_startup_fee?: string;
  normal_weekly_increment?: string;
  gold_weekly_increment?: string;
  premium_weekly_increment?: string;
  listing_free_weeks?: string;
  
  // Image limits
  normal_max_images?: string;
  gold_max_images?: string;
  premium_max_images?: string;
  min_images_per_listing?: string;
  max_image_size_mb?: string;
  
  // Content limits
  max_title_length?: string;
  max_description_length?: string;
  
  // Location limits
  normal_max_locations?: string;
  gold_max_locations?: string;
  premium_max_locations?: string;
  
  // Tier thresholds
  gold_threshold?: string;
  premium_threshold?: string;
  
  // Referral
  min_referral_withdrawal?: string;
  max_referral_withdrawal?: string;
  
  // Payment
  telebirr_merchant_phone?: string;
  safaricom_merchant_phone?: string;
  cbe_account_number?: string;
  
  // Upgrade pricing
  gold_upgrade_price?: string;
  premium_upgrade_price?: string;
  buyer_gold_upgrade_price?: string;
  buyer_premium_upgrade_price?: string;
  
  // Security
  otp_expiry_minutes?: string;
  currency_tolerance?: string;
  
  // Company info
  company_phone?: string;
  company_email?: string;
  company_location?: string;
}

interface SettingsContextValue {
  settings: PlatformSettings | null;
  loading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
}

const SettingsContext = createContext<SettingsContextValue>({
  settings: null,
  loading: true,
  error: null,
  refresh: async () => {},
});

/**
 * SettingsProvider - Loads platform settings from backend and makes them available app-wide.
 * 
 * Settings are fetched:
 * - On mount
 * - Every 5 minutes (auto-refresh for admin changes)
 * - On window focus (catches changes when user returns)
 * 
 * If fetch fails, uses null to allow components to fall back to default constants.
 */
export function SettingsProvider({ children }: { children: ReactNode }) {
  const [settings, setSettings] = useState<PlatformSettings | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchSettings = async () => {
    try {
      const data = await api.getPublicSettings();
      setSettings(data as PlatformSettings);
      setError(null);
    } catch (err: any) {
      console.error('[SettingsProvider] Failed to load settings:', err);
      setError(err?.message || 'Failed to load settings');
      // Don't clear existing settings on refresh errors
      if (!settings) {
        setSettings(null); // First load failed, set null to use defaults
      }
    } finally {
      setLoading(false);
    }
  };

  // Initial load
  useEffect(() => {
    fetchSettings();
  }, []);

  // Auto-refresh every 5 minutes
  useEffect(() => {
    const interval = setInterval(() => {
      fetchSettings();
    }, 5 * 60 * 1000); // 5 minutes

    return () => clearInterval(interval);
  }, [settings]);

  // Refresh on window focus (user returns to tab)
  useEffect(() => {
    const handleFocus = () => {
      // Only refresh if settings are already loaded (avoid double-load on mount)
      if (settings) {
        fetchSettings();
      }
    };

    window.addEventListener('focus', handleFocus);
    return () => window.removeEventListener('focus', handleFocus);
  }, [settings]);

  return (
    <SettingsContext.Provider value={{ settings, loading, error, refresh: fetchSettings }}>
      {children}
    </SettingsContext.Provider>
  );
}

/**
 * useSettings - Access platform settings in any component.
 * 
 * Returns null while loading or if fetch failed. Components should fall back to constants.
 * 
 * Example:
 * ```tsx
 * const settings = useSettings();
 * const maxImages = getMaxImages(userTier, settings);
 * ```
 */
export function useSettings() {
  const context = useContext(SettingsContext);
  if (!context) {
    throw new Error('useSettings must be used within SettingsProvider');
  }
  return context.settings;
}

/**
 * useSettingsWithLoading - Access settings with loading state.
 * 
 * Useful for components that need to show loading UI while settings load.
 */
export function useSettingsWithLoading() {
  const context = useContext(SettingsContext);
  if (!context) {
    throw new Error('useSettingsWithLoading must be used within SettingsProvider');
  }
  return context;
}
