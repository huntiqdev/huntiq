'use client';

import { InputHTMLAttributes, forwardRef } from 'react';
import { cn } from '@/lib/utils/cn';
import { FieldWrapper } from './FieldWrapper';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  hint?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, hint, className, id, ...props }, ref) => {
    const inputId = id || (label ? `input-${label.toLowerCase().replace(/\s+/g, '-')}` : undefined);

    return (
      <FieldWrapper label={label} error={error}>
        <input
          ref={ref}
          id={inputId}
          {...props}
          className={cn(
            'w-full border rounded-xl px-3.5 py-2.5 text-base sm:text-sm min-h-[44px] outline-none transition-all',
            'bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-750 text-slate-900 dark:text-slate-100 placeholder:text-slate-400 dark:placeholder:text-slate-500',
            'focus:ring-2 focus:ring-[#1769E0] focus:border-transparent',
            error && 'border-[#DC4655] focus:ring-[#DC4655]',
            className,
          )}
        />
        {hint && !error && (
          <p className="text-[11px] text-slate-500 dark:text-slate-400 mt-1">{hint}</p>
        )}
      </FieldWrapper>
    );
  },
);

Input.displayName = 'Input';
