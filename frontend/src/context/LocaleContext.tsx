'use client';
import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { Locale, Translation, translations, SUPPORTED_LOCALES, LOCALE_LABELS } from '@/lib/i18n';

function resolveDot(obj: Record<string, any>, key: string): string {
  const val = key.split('.').reduce((acc: any, k: string) => acc?.[k], obj);
  return typeof val === 'string' ? val : key;
}

function toSnake(s: string): string {
  return s.replace(/([A-Z])/g, '_$1').toLowerCase();
}

function buildFallback(
  enObj: Record<string, any>,
  localeObj: Record<string, any>
): Record<string, any> {
  const out: Record<string, any> = {};
  for (const key of Object.keys(enObj)) {
    const enVal = enObj[key];
    if (typeof enVal === 'object' && enVal !== null) {
      const nested = (localeObj?.[key] ?? localeObj?.[toSnake(key)] ?? {}) as Record<string, any>;
      out[key] = buildFallback(enVal, nested);
    } else {
      out[key] = localeObj?.[key] ?? localeObj?.[toSnake(key)] ?? enVal;
    }
  }
  return out;
}

interface LocaleContextValue {
  locale: Locale;
  setLocale: (l: Locale) => void;
  t: Translation;
  t_: (key: string) => string;
  locales: readonly Locale[];
  localeLabels: Record<string, string>;
}

const LocaleContext = createContext<LocaleContextValue>({
  locale: 'en',
  setLocale: () => {},
  t: translations.en,
  t_: (k) => resolveDot(translations.en as Record<string, any>, k),
  locales: ['en', 'am', 'or', 'ti'],
  localeLabels: { en: 'English', am: 'አማርኛ', or: 'Afaan Oromoo', ti: 'ትግርኛ' },
});

export function LocaleProvider({ children }: { children: ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>('en');

  // Restore persisted locale on mount
  useEffect(() => {
    try {
      const saved = localStorage.getItem('huntiq_locale') || localStorage.getItem('habesha_locale');
      if (saved && ['en', 'am', 'or', 'ti'].includes(saved)) {
        setLocaleState(saved as Locale);
      }
    } catch {}
  }, []);

  const setLocale = (l: Locale) => {
    setLocaleState(l);
    try {
      localStorage.setItem('huntiq_locale', l);
      localStorage.setItem('habesha_locale', l);
    } catch {}
  };
  const enDict = translations.en as Record<string, any>;
  const localeDict = translations[locale] as Record<string, any>;

  const t = (locale === 'en'
    ? translations.en
    : buildFallback(enDict, localeDict)) as Translation;

  const t_ = (key: string): string => {
    const v = resolveDot(localeDict, key);
    if (v !== key) return v;
    const snake = key.split('.').map(toSnake).join('.');
    const v2 = resolveDot(localeDict, snake);
    if (v2 !== snake) return v2;
    return resolveDot(enDict, key);
  };

  return (
    <LocaleContext.Provider value={{ locale, setLocale, t, t_, locales: SUPPORTED_LOCALES, localeLabels: LOCALE_LABELS }}>
      {children}
    </LocaleContext.Provider>
  );
}

export const useLocaleContext = () => useContext(LocaleContext);
export const useLocale = useLocaleContext;
