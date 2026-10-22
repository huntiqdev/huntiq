'use client';
import {
  AlertCircle, WifiOff, ShieldOff, Lock, Upload,
  CreditCard, KeyRound, Clock, ServerCrash, MapPin,
} from 'lucide-react';
import type { ErrorType } from '@/lib/errors';
import { errorTypeLabel } from '@/lib/errors';
import { useLocaleContext } from '@/context/LocaleContext';

interface Props {
  message: string;
  type?: ErrorType;
  /** Extra Tailwind classes for the outer wrapper */
  className?: string;
  /** Show the error type badge (default true) */
  showType?: boolean;
  /** Optional action button */
  action?: { label: string; onClick: () => void };
}

const TYPE_STYLES: Record<string, {
  bg: string;
  border: string;
  text: string;
  icon: React.ReactNode;
}> = {
  auth:       { bg: 'bg-red-500/10',    border: 'border-red-500/20',    text: 'text-red-400',    icon: <Lock size={13} /> },
  validation: { bg: 'bg-orange-500/10', border: 'border-orange-500/20', text: 'text-orange-400', icon: <AlertCircle size={13} /> },
  network:    { bg: 'bg-blue-500/10',   border: 'border-blue-500/20',   text: 'text-blue-400',   icon: <WifiOff size={13} /> },
  server:     { bg: 'bg-red-500/10',    border: 'border-red-500/20',    text: 'text-red-400',    icon: <ServerCrash size={13} /> },
  permission: { bg: 'bg-yellow-500/10', border: 'border-yellow-500/20', text: 'text-yellow-400', icon: <ShieldOff size={13} /> },
  upload:     { bg: 'bg-orange-500/10', border: 'border-orange-500/20', text: 'text-orange-400', icon: <Upload size={13} /> },
  payment:    { bg: 'bg-purple-500/10', border: 'border-purple-500/20', text: 'text-purple-400', icon: <CreditCard size={13} /> },
  otp:        { bg: 'bg-yellow-500/10', border: 'border-yellow-500/20', text: 'text-yellow-400', icon: <KeyRound size={13} /> },
  rate_limit: { bg: 'bg-zinc-500/10',   border: 'border-zinc-500/20',   text: 'text-zinc-400',   icon: <Clock size={13} /> },
  not_found:  { bg: 'bg-zinc-500/10',   border: 'border-zinc-500/20',   text: 'text-zinc-400',   icon: <MapPin size={13} /> },
  unknown:    { bg: 'bg-red-500/10',    border: 'border-red-500/20',    text: 'text-red-400',    icon: <AlertCircle size={13} /> },
};

export function ErrorBanner({ message, type = 'unknown', className = '', showType = true, action }: Props) {
  const { t } = useLocaleContext();
  const errors = (t as any).errors as Record<string, string>;
  const style = TYPE_STYLES[type] ?? TYPE_STYLES.unknown;
  const label = showType ? errorTypeLabel(type, errors) : null;

  return (
    <div
      role="alert"
      aria-live="polite"
      className={`flex items-start gap-2.5 ${style.bg} border ${style.border} rounded-xl px-3 py-2.5 ${className}`}
    >
      <span className={`${style.text} shrink-0 mt-0.5`} aria-hidden="true">
        {style.icon}
      </span>
      <div className="min-w-0 flex-1">
        {label && (
          <p className={`text-[10px] font-black uppercase tracking-widest ${style.text} opacity-70 mb-0.5`}>
            {label}
          </p>
        )}
        <p className={`text-xs ${style.text} leading-snug`}>{message}</p>
        {action && (
          <button
            type="button"
            onClick={action.onClick}
            className={`mt-1.5 text-[10px] font-bold underline underline-offset-2 ${style.text} hover:opacity-80 transition-opacity`}
          >
            {action.label}
          </button>
        )}
      </div>
    </div>
  );
}

/**
 * Inline field-level error — smaller, no border box, just text + icon.
 * Use under form inputs for field-specific validation errors.
 */
export function FieldError({ message, className = '' }: { message: string; className?: string }) {
  if (!message) return null;
  return (
    <p role="alert" className={`flex items-center gap-1 text-[10px] text-red-400 mt-1 ${className}`}>
      <AlertCircle size={9} aria-hidden="true" />
      {message}
    </p>
  );
}
