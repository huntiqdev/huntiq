import { ReactNode } from 'react';
import { cn } from '@/lib/utils/cn';

interface Props {
  children: ReactNode;
  className?: string;
  padding?: boolean;
}

export function Card({ children, className, padding = true }: Props) {
  return (
    <div className={cn(
      'bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl',
      padding && 'p-5',
      className,
    )}>
      {children}
    </div>
  );
}
