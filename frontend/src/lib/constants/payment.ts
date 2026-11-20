/** Authoritative payment method definition: Telebirr only (official Ethiopian mobile payment integration). */
export interface PaymentMethod {
  value: 'telebirr';
  label: string;
  emoji: string;
  /** USSD dial code */
  ussd: string;
  appLink: string;
  appName: string;
  refPlaceholder: string;
  settingsKey: 'telebirr_merchant_phone';
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
    appName: 'Telebirr SuperApp',
    refPlaceholder: 'Auto-verified by Telebirr system',
    settingsKey: 'telebirr_merchant_phone',
    borderColor: 'border-[#1769E0]',
    activeBg: 'bg-blue-50 dark:bg-blue-900/20',
    accentColor: 'text-[#1769E0]',
  },
];

export type PaymentMethodValue = 'telebirr';

/** Default merchant targets */
export const DEFAULT_MERCHANT: Record<PaymentMethod['settingsKey'], string> = {
  telebirr_merchant_phone: '',
};

/** Extract merchant target string from live settings for a given method */
export function getMerchantTarget(
  method: PaymentMethodValue = 'telebirr',
  settings?: Partial<Record<PaymentMethod['settingsKey'], string>>
): string {
  const m = PAYMENT_METHODS[0];
  if (!settings) return DEFAULT_MERCHANT[m.settingsKey];
  return settings[m.settingsKey] ?? DEFAULT_MERCHANT[m.settingsKey];
}
