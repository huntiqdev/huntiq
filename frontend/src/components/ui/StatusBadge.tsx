const colors: Record<string, string> = {
  pending:   'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-400',
  approved:  'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-400',
  rejected:  'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-400',
  completed: 'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-400',
  active:    'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-400',
  inactive:  'bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400',
  cancelled: 'bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400',
};

export function StatusBadge({ status }: { status: string }) {
  return (
    <span className={`inline-block text-xs font-bold px-2.5 py-0.5 rounded-full capitalize ${colors[status] ?? 'bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400'}`}>
      {status}
    </span>
  );
}
