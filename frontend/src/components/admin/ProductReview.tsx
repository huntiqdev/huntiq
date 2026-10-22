'use client';
import { useState } from 'react';
import { Check, X } from 'lucide-react';
import { Button } from '@/components/ui/Button';

interface Props {
  productId: number;
  onReview: (id: number, status: 'approved' | 'rejected', notes: string) => Promise<void>;
}

export function ProductReview({ productId, onReview }: Props) {
  const [notes, setNotes] = useState('');
  const [loading, setLoading] = useState(false);

  const handle = async (status: 'approved' | 'rejected') => {
    setLoading(true);
    await onReview(productId, status, notes);
    setLoading(false);
    setNotes('');
  };

  return (
    <div className="flex flex-col gap-2">
      <textarea
        placeholder="Admin notes (optional for approval, recommended for rejection)"
        value={notes}
        onChange={e => setNotes(e.target.value)}
        rows={2}
        className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800"
      />
      <div className="flex gap-2">
        <Button variant="primary" size="sm" disabled={loading} onClick={() => handle('approved')}>
          <Check size={14} />Approve
        </Button>
        <Button variant="danger" size="sm" disabled={loading} onClick={() => handle('rejected')}>
          <X size={14} />Reject
        </Button>
      </div>
    </div>
  );
}
