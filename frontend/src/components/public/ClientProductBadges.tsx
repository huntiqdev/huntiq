'use client';
import { useLocale } from '@/hooks/useLocale';

interface Props {
  categorySlug: string;
  subcategory?: string;
  item?: string;
  condition?: string;
}

const CONDITION_STYLE: Record<string, string> = {
  new:  'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400',
  used: 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400',
};

const CONDITION_UI_KEY: Record<string, string> = {
  new:  'condition_new',
  used: 'condition_used',
};

export function ClientProductBadges({ categorySlug, subcategory, item, condition }: Props) {
  const { t } = useLocale();

  const subLabel = subcategory
    ? ((t.sub_categories as any)?.[categorySlug]?.[subcategory]
        ?? (t.items as any)?.[subcategory]
        ?? subcategory.replace(/_/g, ' '))
    : null;

  const itemLabel = item
    ? ((t.items as any)?.[item] ?? item.replace(/_/g, ' '))
    : null;

  const condLabel = condition
    ? ((t.ui as any)?.[CONDITION_UI_KEY[condition]] ?? condition.replace(/_/g, ' '))
    : null;

  return (
    <>
      {subLabel && (
        <span className="text-[10px] font-black uppercase tracking-widest text-purple-600 bg-purple-600/10 px-2 py-1 rounded-lg">
          {subLabel}
        </span>
      )}
      {itemLabel && (
        <span className="text-[10px] font-black uppercase tracking-widest text-green-600 bg-green-600/10 px-2 py-1 rounded-lg">
          {itemLabel}
        </span>
      )}
      {condLabel && (
        <span className={`text-[10px] font-black uppercase tracking-widest px-2 py-1 rounded-lg ${CONDITION_STYLE[condition!] ?? 'bg-zinc-100 dark:bg-zinc-800'}`}>
          {condLabel}
        </span>
      )}
    </>
  );
}
