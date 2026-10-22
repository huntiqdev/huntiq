import { SelectHTMLAttributes, forwardRef } from 'react';
import { cn } from '@/lib/utils/cn';
import { FieldWrapper } from './FieldWrapper';

interface Props extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  error?: string;
}

export const Select = forwardRef<HTMLSelectElement, Props>(({ label, error, className, children, ...props }, ref) => (
  <FieldWrapper label={label} error={error}>
    <select
      ref={ref}
      {...props}
      className={cn(
        'w-full border rounded-xl px-3 py-2.5 text-sm outline-none transition-all',
        'bg-zinc-50 dark:bg-zinc-800 border-zinc-200 dark:border-zinc-700',
        'focus:ring-2 focus:ring-orange-600',
        error && 'border-red-500',
        className,
      )}
    >
      {children}
    </select>
  </FieldWrapper>
));
Select.displayName = 'Select';
