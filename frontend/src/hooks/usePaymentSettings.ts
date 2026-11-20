'use client';
import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { DEFAULT_MERCHANT } from '@/lib/constants/payment';

export interface PaymentSettings {
  telebirr_merchant_phone: string;
  gold_upgrade_price: number;
  premium_upgrade_price: number;
  buyer_gold_upgrade_price: number;
  buyer_premium_upgrade_price: number;
}

const DEFAULTS: PaymentSettings = {
  ...DEFAULT_MERCHANT,
  gold_upgrade_price: 500,
  premium_upgrade_price: 2000,
  buyer_gold_upgrade_price: 200,
  buyer_premium_upgrade_price: 800,
};

export function usePaymentSettings(): { settings: PaymentSettings; loading: boolean } {
  const [settings, setSettings] = useState<PaymentSettings>(DEFAULTS);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.getPublicPaymentSettings()
      .then(s => {
        if (!s) return;
        setSettings({
          telebirr_merchant_phone:  s.telebirr_merchant_phone  || DEFAULTS.telebirr_merchant_phone,
          gold_upgrade_price:       Number(s.gold_upgrade_price)             || DEFAULTS.gold_upgrade_price,
          premium_upgrade_price:    Number(s.premium_upgrade_price)          || DEFAULTS.premium_upgrade_price,
          buyer_gold_upgrade_price:    Number(s.buyer_gold_upgrade_price)    || DEFAULTS.buyer_gold_upgrade_price,
          buyer_premium_upgrade_price: Number(s.buyer_premium_upgrade_price) || DEFAULTS.buyer_premium_upgrade_price,
        });
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  return { settings, loading };
}
