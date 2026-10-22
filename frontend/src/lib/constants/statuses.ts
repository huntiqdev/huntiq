import { Clock, CheckCircle, XCircle, type LucideIcon } from 'lucide-react';

export interface StatusConfig {
  label: string;
  color: string;
  icon: LucideIcon;
}

/** Shared order/request status config used by buyer and seller order pages */
export const ORDER_STATUS: Record<string, StatusConfig> = {
  pending:   { label: 'Pending',   color: 'bg-blue-100 dark:bg-blue-900/30 text-blue-600',           icon: Clock },
  accepted:  { label: 'Accepted',  color: 'bg-green-100 dark:bg-green-900/30 text-green-600',         icon: CheckCircle },
  completed: { label: 'Completed', color: 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600',   icon: CheckCircle },
  declined:  { label: 'Declined',  color: 'bg-red-100 dark:bg-red-900/30 text-red-500',               icon: XCircle },
  cancelled: { label: 'Cancelled', color: 'bg-zinc-100 dark:bg-zinc-800 text-zinc-500',               icon: XCircle },
};
