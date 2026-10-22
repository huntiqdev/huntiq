'use client';
import { useLocale } from '@/hooks/useLocale';

interface Props { k: string; className?: string }

/** Renders a translation value by dot-notation key from the current locale */
export function ClientLabel({ k, className }: Props) {
  const { t_ } = useLocale();
  const text = t_(k);
  return className ? <span className={className}>{text}</span> : <>{text}</>;
}
