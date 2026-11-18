'use client';

import { ButtonHTMLAttributes, ReactNode } from 'react';
import { cn } from '@/lib/utils/cn';
import { Loader2 } from 'lucide-react';

export type ButtonVariant =
  | 'primary'
  | 'secondary'
  | 'tertiary'
  | 'outline'
  | 'ghost'
  | 'destructive'
  | 'danger'
  | 'success'
  | 'premium'
  | 'icon-only';

export type ButtonSize = 'sm' | 'md' | 'lg';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  loading?: boolean;
  loadingText?: string;
  fullWidth?: boolean;
  children: ReactNode;
}

const variantStyles: Record<ButtonVariant, string> = {
  primary:
    'bg-[#1769E0] hover:bg-[#1055BD] active:bg-[#103B78] text-white shadow-sm shadow-[#1769E0]/25 focus-visible:ring-2 focus-visible:ring-[#1769E0] focus-visible:ring-offset-2 dark:focus-visible:ring-offset-slate-900',
  secondary:
    'bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-100 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700/80 shadow-sm focus-visible:ring-2 focus-visible:ring-slate-400',
  tertiary:
    'bg-[#EAF3FF] dark:bg-[#1769E0]/15 text-[#1769E0] dark:text-[#38A5FF] hover:bg-[#D5E7FF] dark:hover:bg-[#1769E0]/25 focus-visible:ring-2 focus-visible:ring-[#1769E0]',
  outline:
    'border-2 border-[#1769E0] dark:border-[#38A5FF] text-[#1769E0] dark:text-[#38A5FF] hover:bg-[#EAF3FF] dark:hover:bg-[#1769E0]/15 focus-visible:ring-2 focus-visible:ring-[#1769E0]',
  ghost:
    'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 focus-visible:ring-2 focus-visible:ring-slate-400',
  destructive:
    'bg-[#DC4655] hover:bg-[#B93845] active:bg-[#96232E] text-white shadow-sm shadow-[#DC4655]/25 focus-visible:ring-2 focus-visible:ring-[#DC4655]',
  danger:
    'bg-[#DC4655] hover:bg-[#B93845] active:bg-[#96232E] text-white shadow-sm shadow-[#DC4655]/25 focus-visible:ring-2 focus-visible:ring-[#DC4655]',
  success:
    'bg-[#16A36A] hover:bg-[#128354] active:bg-[#0E6B45] text-white shadow-sm shadow-[#16A36A]/25 focus-visible:ring-2 focus-visible:ring-[#16A36A]',
  premium:
    'bg-gradient-to-r from-[#8B5CF6] to-[#6366F1] hover:from-[#7C3AED] hover:to-[#4F46E5] text-white shadow-sm shadow-[#8B5CF6]/30 focus-visible:ring-2 focus-visible:ring-[#8B5CF6]',
  'icon-only':
    'p-2.5 rounded-xl border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 focus-visible:ring-2 focus-visible:ring-[#1769E0]',
};

const sizeStyles: Record<ButtonSize, string> = {
  sm: 'px-3 py-1.5 text-xs rounded-xl min-h-[38px]',
  md: 'px-4 py-2.5 text-sm rounded-xl min-h-[44px]',
  lg: 'px-6 py-3.5 text-base rounded-2xl min-h-[48px]',
};

export function Button({
  variant = 'primary',
  size = 'md',
  loading = false,
  loadingText,
  fullWidth = false,
  className,
  children,
  disabled,
  onClick,
  ...props
}: ButtonProps) {
  const isDisabled = disabled || loading;

  return (
    <button
      {...props}
      disabled={isDisabled}
      aria-busy={loading}
      onClick={(e) => {
        if (isDisabled) {
          e.preventDefault();
          return;
        }
        if (onClick) onClick(e);
      }}
      className={cn(
        'font-bold transition-all tap-bounce active:scale-97 select-none',
        'inline-flex items-center justify-center gap-2 outline-none',
        'disabled:opacity-50 disabled:cursor-not-allowed disabled:pointer-events-none',
        variantStyles[variant],
        variant !== 'icon-only' && sizeStyles[size],
        fullWidth && 'w-full',
        className,
      )}
    >
      {loading ? (
        <>
          <Loader2 className="w-4 h-4 animate-spin shrink-0" />
          <span>{loadingText || children}</span>
        </>
      ) : (
        children
      )}
    </button>
  );
}
