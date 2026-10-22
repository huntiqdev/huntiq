import { StatusBadge } from '@/components/ui/StatusBadge';

interface Props {
  status: string;
  adminNotes?: string;
}

export function ApprovalStatus({ status, adminNotes }: Props) {
  return (
    <div className="flex flex-col gap-1">
      <StatusBadge status={status} />
      {adminNotes && status === 'rejected' && (
        <p className="text-xs text-red-500 mt-1 bg-red-50 dark:bg-red-900/20 px-2 py-1 rounded-lg">
          {adminNotes}
        </p>
      )}
    </div>
  );
}
