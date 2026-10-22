'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Tag, Plus, X, Save, RefreshCw, Hash, Type, ToggleLeft, List } from 'lucide-react';

const DATA_TYPES = [
  { value: 'text', label: 'Text', icon: Type },
  { value: 'number', label: 'Number', icon: Hash },
  { value: 'boolean', label: 'Yes/No', icon: ToggleLeft },
  { value: 'select', label: 'Single Select', icon: List },
  { value: 'multiselect', label: 'Multi Select', icon: List },
  { value: 'range', label: 'Range', icon: Hash },
  { value: 'date', label: 'Date', icon: Type },
];

const EMPTY_FORM = {
  name_en: '', name_am: '', name_om: '', slug: '',
  description: '', data_type: 'text', unit: '',
  is_required: false, is_searchable: true, is_filterable: true,
  display_order: 0, options: '' // options as comma-separated string for select types
};

export default function AttributesAdminPage() {
  const { user, isSuperAdmin } = useAuth();
  const [attributes, setAttributes] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [form, setForm] = useState({ ...EMPTY_FORM });
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => { fetchAttributes(); }, []);

  const fetchAttributes = async () => {
    setLoading(true);
    try {
      const res = await fetch('/api/attributes');
      const data = await res.json();
      if (data.success) setAttributes(data.data || []);
    } catch { setAttributes([]); } finally { setLoading(false); }
  };

  const autoSlug = (name: string) =>
    name.toLowerCase().replace(/[^a-z0-9\s-]/g, '').replace(/\s+/g, '-').replace(/-+/g, '-').trim();

  const handleSubmit = async () => {
    if (!user?.token || !form.name_en || !form.slug) { setError('Name and slug are required.'); return; }
    setSubmitting(true);
    setError('');
    try {
      // Build validation rules from options
      const validation_rules: Record<string, any> = {};
      if ((form.data_type === 'select' || form.data_type === 'multiselect') && form.options) {
        validation_rules.options = form.options.split(',').map(o => o.trim()).filter(Boolean);
      }
      await api.adminCreateAttribute({
        name_en: form.name_en, name_am: form.name_am, name_om: form.name_om,
        slug: form.slug, description: form.description, data_type: form.data_type,
        unit: form.unit, is_required: form.is_required, is_searchable: form.is_searchable,
        is_filterable: form.is_filterable, display_order: form.display_order,
        ...(Object.keys(validation_rules).length > 0 ? { validation_rules } : {})
      }, user.token);
      setShowModal(false);
      setForm({ ...EMPTY_FORM });
      fetchAttributes();
    } catch (e: any) {
      setError(e?.message || 'Create failed.');
    } finally { setSubmitting(false); }
  };

  const dataTypeBadge = (dt: string) => {
    const colors: Record<string, string> = {
      text: 'bg-blue-900/30 text-blue-400', number: 'bg-purple-900/30 text-purple-400',
      boolean: 'bg-emerald-900/30 text-emerald-400', select: 'bg-amber-900/30 text-amber-400',
      multiselect: 'bg-orange-900/30 text-orange-400', range: 'bg-cyan-900/30 text-cyan-400',
      date: 'bg-zinc-700 text-zinc-300',
    };
    return colors[dt] ?? 'bg-zinc-700 text-zinc-300';
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-purple-600/10 rounded-xl border border-purple-500/20">
            <Tag className="w-6 h-6 text-purple-400" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Attribute Definitions</h1>
            <p className="text-xs text-zinc-400">Manage dynamic product attributes (Brand, Model, Storage, etc.)</p>
          </div>
        </div>
        <div className="flex gap-2">
          <button onClick={fetchAttributes} className="p-2 bg-zinc-800 rounded-lg hover:bg-zinc-700 transition-colors">
            <RefreshCw className="w-4 h-4 text-zinc-400" />
          </button>
          {isSuperAdmin && (
            <button onClick={() => { setForm({ ...EMPTY_FORM }); setError(''); setShowModal(true); }}
              className="flex items-center gap-2 px-4 py-2 bg-orange-600 hover:bg-orange-700 text-white rounded-xl text-sm font-bold transition-colors">
              <Plus className="w-4 h-4" /> New Attribute
            </button>
          )}
        </div>
      </div>

      <div className="bg-zinc-900 rounded-2xl border border-zinc-800 overflow-hidden">
        {loading ? (
          <div className="p-8 text-center text-zinc-500">
            <div className="animate-spin w-8 h-8 border-2 border-orange-500 border-t-transparent rounded-full mx-auto mb-3" />
            Loading attributes...
          </div>
        ) : attributes.length === 0 ? (
          <div className="p-8 text-center text-zinc-500">
            <Tag className="w-12 h-12 mx-auto mb-3 opacity-20" />
            <p>No attributes defined yet.</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-zinc-800 text-zinc-500 text-xs uppercase font-bold">
                <th className="px-6 py-4 text-left">Name</th>
                <th className="px-6 py-4 text-left">Slug</th>
                <th className="px-6 py-4 text-left">Data Type</th>
                <th className="px-6 py-4 text-left">Unit</th>
                <th className="px-6 py-4 text-left">Flags</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/50">
              {attributes.map(attr => (
                <tr key={attr.id} className="hover:bg-zinc-800/40 transition-colors">
                  <td className="px-6 py-4">
                    <p className="font-semibold text-white">{attr.name_en}</p>
                    {attr.name_am && <p className="text-xs text-zinc-500">{attr.name_am}</p>}
                  </td>
                  <td className="px-6 py-4 font-mono text-zinc-400 text-xs">{attr.slug}</td>
                  <td className="px-6 py-4">
                    <span className={`px-2.5 py-1 rounded-full text-xs font-bold ${dataTypeBadge(attr.data_type)}`}>
                      {attr.data_type}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-zinc-400">{attr.unit || '—'}</td>
                  <td className="px-6 py-4">
                    <div className="flex gap-1.5 flex-wrap">
                      {attr.is_required && <span className="text-[10px] px-1.5 py-0.5 bg-red-900/30 text-red-400 rounded font-bold">REQ</span>}
                      {attr.is_filterable && <span className="text-[10px] px-1.5 py-0.5 bg-blue-900/30 text-blue-400 rounded font-bold">FILTER</span>}
                      {attr.is_searchable && <span className="text-[10px] px-1.5 py-0.5 bg-emerald-900/30 text-emerald-400 rounded font-bold">SEARCH</span>}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Create Modal */}
      {showModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4">
          <div className="bg-zinc-900 border border-zinc-700 rounded-2xl w-full max-w-lg shadow-2xl overflow-y-auto max-h-[90vh]">
            <div className="px-6 py-5 border-b border-zinc-800 flex items-center justify-between sticky top-0 bg-zinc-900 z-10">
              <h2 className="text-lg font-black text-white">New Attribute</h2>
              <button onClick={() => setShowModal(false)}><X className="w-5 h-5 text-zinc-500 hover:text-white" /></button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Name (English) *</label>
                <input value={form.name_en} onChange={e => setForm(f => ({ ...f, name_en: e.target.value, slug: autoSlug(e.target.value) }))}
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-2.5 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                  placeholder="Brand" />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Amharic</label>
                  <input value={form.name_am} onChange={e => setForm(f => ({ ...f, name_am: e.target.value }))}
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-2.5 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                    placeholder="ብራንድ" />
                </div>
                <div>
                  <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Oromo</label>
                  <input value={form.name_om} onChange={e => setForm(f => ({ ...f, name_om: e.target.value }))}
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-2.5 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                    placeholder="Maqaa Daldalaa" />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Slug *</label>
                  <input value={form.slug} onChange={e => setForm(f => ({ ...f, slug: autoSlug(e.target.value) }))}
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-2.5 text-sm text-white font-mono placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                    placeholder="brand" />
                </div>
                <div>
                  <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Unit</label>
                  <input value={form.unit} onChange={e => setForm(f => ({ ...f, unit: e.target.value }))}
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-2.5 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                    placeholder="GB, km, HP..." />
                </div>
              </div>
              <div>
                <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Data Type *</label>
                <div className="grid grid-cols-4 gap-2">
                  {DATA_TYPES.map(dt => (
                    <button key={dt.value} type="button"
                      onClick={() => setForm(f => ({ ...f, data_type: dt.value }))}
                      className={`py-2 px-2 rounded-xl text-xs font-bold border transition-all ${
                        form.data_type === dt.value
                          ? 'bg-orange-600/20 border-orange-500/40 text-orange-400'
                          : 'bg-zinc-800 border-zinc-700 text-zinc-400 hover:border-zinc-500'
                      }`}>
                      {dt.label}
                    </button>
                  ))}
                </div>
              </div>

              {(form.data_type === 'select' || form.data_type === 'multiselect') && (
                <div>
                  <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Options (comma-separated)</label>
                  <input value={form.options} onChange={e => setForm(f => ({ ...f, options: e.target.value }))}
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-2.5 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                    placeholder="Samsung, Apple, Huawei, Nokia..." />
                </div>
              )}

              <div className="grid grid-cols-3 gap-3">
                {[
                  { key: 'is_required', label: 'Required' },
                  { key: 'is_filterable', label: 'Filterable' },
                  { key: 'is_searchable', label: 'Searchable' },
                ].map(({ key, label }) => (
                  <label key={key} className="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" checked={(form as any)[key]}
                      onChange={e => setForm(f => ({ ...f, [key]: e.target.checked }))}
                      className="rounded border-zinc-600 bg-zinc-800 text-orange-600 focus:ring-orange-500" />
                    <span className="text-sm text-zinc-300">{label}</span>
                  </label>
                ))}
              </div>

              {error && <p className="text-red-400 text-sm">{error}</p>}

              <div className="flex gap-3 pt-2">
                <button onClick={() => setShowModal(false)} className="flex-1 py-3 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-sm font-bold">Cancel</button>
                <button onClick={handleSubmit} disabled={submitting}
                  className="flex-1 py-3 bg-orange-600 hover:bg-orange-700 text-white rounded-xl text-sm font-bold disabled:opacity-50 flex items-center justify-center gap-2">
                  <Save className="w-4 h-4" /> {submitting ? 'Saving...' : 'Create'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
