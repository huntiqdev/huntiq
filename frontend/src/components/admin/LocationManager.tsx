'use client';
import { useState } from 'react';
import { Plus, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';

interface Location { id: number; region: string; zone?: string; subcity?: string; woreda?: string }

interface Props {
  locations: Location[];
  onAdd: (data: Omit<Location, 'id'>) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
}

export function LocationManager({ locations, onAdd, onDelete }: Props) {
  const [form, setForm] = useState({ region: '', zone: '', subcity: '', woreda: '' });
  const [loading, setLoading] = useState(false);

  const handleAdd = async () => {
    if (!form.region.trim()) return;
    setLoading(true);
    await onAdd({ region: form.region, zone: form.zone || undefined, subcity: form.subcity || undefined, woreda: form.woreda || undefined });
    setForm({ region: '', zone: '', subcity: '', woreda: '' });
    setLoading(false);
  };

  return (
    <div className="flex flex-col gap-4">
      <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
        {(['region', 'zone', 'subcity', 'woreda'] as const).map(field => (
          <Input key={field} placeholder={field.charAt(0).toUpperCase() + field.slice(1)}
            value={form[field]} onChange={e => setForm(p => ({ ...p, [field]: e.target.value }))} />
        ))}
      </div>
      <Button size="sm" disabled={loading || !form.region.trim()} onClick={handleAdd}>
        <Plus size={14} />Add Location
      </Button>
      <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="bg-zinc-50 dark:bg-zinc-800 text-xs uppercase text-zinc-500">
            <tr>{['ID', 'Region', 'Zone', 'Subcity', 'Woreda', ''].map(h => <th key={h} className="px-3 py-2 text-left">{h}</th>)}</tr>
          </thead>
          <tbody className="divide-y divide-zinc-100 dark:divide-zinc-800">
            {locations.map(l => (
              <tr key={l.id}>
                <td className="px-3 py-2 text-zinc-400">#{l.id}</td>
                <td className="px-3 py-2 font-medium">{l.region}</td>
                <td className="px-3 py-2 text-zinc-500">{l.zone ?? '—'}</td>
                <td className="px-3 py-2 text-zinc-500">{l.subcity ?? '—'}</td>
                <td className="px-3 py-2 text-zinc-500">{l.woreda ?? '—'}</td>
                <td className="px-3 py-2">
                  <button onClick={() => onDelete(l.id)} className="text-red-500 hover:text-red-700 transition-colors">
                    <Trash2 size={14} />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
