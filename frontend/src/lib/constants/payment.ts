/** Shared payment method definitions used across seller pay, tier upgrade, and buyer upgrade. */
export interface PaymentMethod {
  value: 'telebirr' | 'safaricom' | 'cbe';
  label: string;
  emoji: string;
  /** USSD dial code — null for CBE */
  ussd: string | null;
  appLink: string;
  appName: string;
  refPlaceholder: string;
  /** Key in platform_settings to look up the merchant phone/account */
  settingsKey: 'telebirr_merchant_phone' | 'safaricom_merchant_phone' | 'cbe_account_number';
  borderColor: string;
  activeBg: string;
  accentColor: string;
}

export const PAYMENT_METHODS: PaymentMethod[] = [
  {
    value: 'telebirr',
    label: 'Telebirr',
    emoji: '📱',
    ussd: '*127*6*1#',
    appLink: 'https://play.google.com/store/apps/details?id=com.ethiotelecom.telebirr',
    appName: 'Telebirr App',
    refPlaceholder: 'e.g. TXN123456789',
    settingsKey: 'telebirr_merchant_phone',
    borderColor: 'border-orange-400',
    activeBg: 'bg-orange-50 dark:bg-orange-900/20',
    accentColor: 'text-orange-600',
  },
  {
    value: 'safaricom',
    label: 'M-PESA (Safaricom)',
    emoji: '🟢',
    ussd: '*880#',
    appLink: 'https://play.google.com/store/apps/details?id=com.safaricom.mpesa.lifestyle',
    appName: 'M-PESA App',
    refPlaceholder: 'e.g. MPE123456789',
    settingsKey: 'safaricom_merchant_phone',
    borderColor: 'border-green-400',
    activeBg: 'bg-green-50 dark:bg-green-900/20',
    accentColor: 'text-green-600',
  },
  {
    value: 'cbe',
    label: 'CBE Birr',
    emoji: '🏦',
    ussd: null,
    appLink: 'https://play.google.com/store/apps/details?id=com.cbe.cbebirr',
    appName: 'CBE Birr App',
    refPlaceholder: 'e.g. FT12345678',
    settingsKey: 'cbe_account_number',
    borderColor: 'border-blue-400',
    activeBg: 'bg-blue-50 dark:bg-blue-900/20',
    accentColor: 'text-blue-600',
  },
];

export type PaymentMethodValue = PaymentMethod['value'];

/** Default merchant targets — empty strings; must be configured in platform_settings via admin panel */
export const DEFAULT_MERCHANT: Record<PaymentMethod['settingsKey'], string> = {
  telebirr_merchant_phone:  '',
  safaricom_merchant_phone: '',
  cbe_account_number:       '',
};

/** Extract merchant target string from live settings for a given method */
export function getMerchantTarget(
  method: PaymentMethodValue,
  settings: Partial<Record<PaymentMethod['settingsKey'], string>>
): string {
  const m = PAYMENT_METHODS.find(p => p.value === method)!;
  return settings[m.settingsKey] ?? DEFAULT_MERCHANT[m.settingsKey];
}
