import { InputHTMLAttributes, forwardRef } from 'react';
import { cn } from '@/lib/utils/cn';
import { FieldWrapper } from './FieldWrapper';

interface Props extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Input = forwardRef<HTMLInputElement, Props>(({ label, error, className, ...props }, ref) => (
  <FieldWrapper label={label} error={error}>
    <input
      ref={ref}
      {...props}
      className={cn(
        'w-full border rounded-xl px-3 py-2.5 text-sm outline-none transition-all',
        'bg-zinc-50 dark:bg-zinc-800 border-zinc-200 dark:border-zinc-700',
        'focus:ring-2 focus:ring-orange-600 focus:border-transparent',
        error && 'border-red-500 focus:ring-red-500',
        className,
      )}
    />
  </FieldWrapper>
));
Input.displayName = 'Input';
