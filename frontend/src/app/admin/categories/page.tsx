'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import {
  FolderTree, Plus, Pencil, Trash2, ChevronRight, ChevronDown,
  RefreshCw, Folder, FolderOpen, X, Save, CheckCircle2, AlertCircle, Eye, EyeOff
} from 'lucide-react';

interface CategoryNode {
  id: number;
  name_en: string;
  name_am?: string;
  name_om?: string;
  name_ti?: string;
  slug: string;
  level: number;
  icon?: string;
  image_url?: string;
  is_leaf: boolean;
  is_active: boolean;
  display_order: number;
  parent_id?: number;
  children?: CategoryNode[];
}

const EMPTY_FORM = {
  name_en: '',
  name_am: '',
  name_om: '',
  name_ti: '',
  slug: '',
  icon: '',
  image_url: '',
  description: '',
  is_leaf: false,
  is_active: true,
  display_order: 0,
  parent_id: null as number | null,
};

export default function CategoriesAdminPage() {
  const { user, isSuperAdmin, isPlatformOwner, isPlatformAdmin, hasAdminPermission } = useAuth();
  const canManage = isSuperAdmin || isPlatformOwner || isPlatformAdmin || hasAdminPermission('categories');
  const [tree, setTree] = useState<CategoryNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState<Set<number>>(new Set());
  const [showModal, setShowModal] = useState(false);
  const [editItem, setEditItem] = useState<CategoryNode | null>(null);
  const [form, setForm] = useState({ ...EMPTY_FORM });
  const [slugChecking, setSlugChecking] = useState(false);
  const [slugAvailable, setSlugAvailable] = useState<boolean | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [deleteConfirm, setDeleteConfirm] = useState<number | null>(null);

  useEffect(() => { fetchTree(); }, []);

  const fetchTree = async () => {
    setLoading(true);
    try {
      const data = await api.getCategoryTree('en');
      setTree(Array.isArray(data) ? data : []);
    } catch {
      setTree([]);
    } finally {
      setLoading(false);
    }
  };

  const openCreate = (parentId?: number) => {
    setEditItem(null);
    setForm({ ...EMPTY_FORM, parent_id: parentId ?? null });
    setSlugAvailable(null);
    setError('');
    setShowModal(true);
  };

  const openEdit = (cat: CategoryNode) => {
    setEditItem(cat);
    setForm({
      name_en: cat.name_en,
      name_am: cat.name_am || '',
      name_om: cat.name_om || '',
      name_ti: cat.name_ti || '',
      slug: cat.slug,
      icon: cat.icon || '',
      image_url: cat.image_url || '',
      description: '',
      is_leaf: cat.is_leaf,
      is_active: cat.is_active,
      display_order: cat.display_order,
      parent_id: cat.parent_id ?? null,
    });
    setSlugAvailable(true);
    setError('');
    setShowModal(true);
  };

  const handleSlugBlur = async () => {
    if (!form.slug.trim()) return;
    setSlugChecking(true);
    try {
      const res = await api.checkCategorySlug(form.slug.trim(), editItem?.id);
      setSlugAvailable(res.available);
    } catch {
      setSlugAvailable(null);
    } finally {
      setSlugChecking(false);
    }
  };

  const handleSubmit = async () => {
    if (!user?.token || !form.name_en || !form.slug) {
      setError('Name (English) and slug are required.');
      return;
    }
    if (slugAvailable === false) {
      setError('Slug is already taken by another category. Please choose a unique slug.');
      return;
    }
    setSubmitting(true);
    setError('');
    try {
      const payload = {
        name_en: form.name_en,
        name_am: form.name_am || undefined,
        name_om: form.name_om || undefined,
        name_ti: form.name_ti || undefined,
        slug: form.slug,
        icon: form.icon || undefined,
        image_url: form.image_url || undefined,
        is_leaf: form.is_leaf,
        is_active: form.is_active,
        display_order: form.display_order,
        parent_id: form.parent_id || undefined,
      };
      if (editItem) {
        await api.adminUpdateCategoryHierarchy(editItem.id, payload, user.token);
      } else {
        await api.adminCreateCategoryHierarchy(payload, user.token);
      }
      setShowModal(false);
      fetchTree();
    } catch (e: any) {
      setError(e?.message || 'Save failed.');
    } finally {
      setSubmitting(false);
    }
  };

  const handleToggleStatus = async (node: CategoryNode) => {
    if (!user?.token) return;
    try {
      await api.adminSetCategoryStatus(node.id, !node.is_active, user.token);
      fetchTree();
    } catch (e: any) {
      alert(e?.message || 'Failed to toggle category status.');
    }
  };

  const handleDelete = async (id: number) => {
    if (!user?.token) return;
    try {
      await api.adminDeleteCategoryHierarchy(id, user.token);
      setDeleteConfirm(null);
      fetchTree();
    } catch (e: any) {
      alert(e?.message || 'Delete rejected — category may have children or referenced products. Use deactivation instead.');
    }
  };

  const toggle = (id: number) => {
    const s = new Set(expanded);
    s.has(id) ? s.delete(id) : s.add(id);
    setExpanded(s);
  };

  const autoSlug = (name: string) =>
    name.toLowerCase().replace(/[^a-z0-9\s-]/g, '').replace(/\s+/g, '-').replace(/-+/g, '-').trim();

  const renderNode = (node: CategoryNode, depth = 0): React.ReactNode => {
    const open = expanded.has(node.id);
    const hasChildren = node.children && node.children.length > 0;
    return (
      <div key={node.id}>
        <div
          className={`flex items-center gap-2 px-4 py-3 hover:bg-zinc-800/60 rounded-xl transition-colors group`}
          style={{ paddingLeft: `${16 + depth * 24}px` }}
        >
          <button onClick={() => hasChildren && toggle(node.id)} className="w-5 flex items-center justify-center">
            {hasChildren ? (open ? <ChevronDown className="w-4 h-4 text-zinc-400" /> : <ChevronRight className="w-4 h-4 text-zinc-400" />) : null}
          </button>
          {open ? <FolderOpen className="w-4 h-4 text-amber-400 shrink-0" /> : <Folder className="w-4 h-4 text-zinc-500 shrink-0" />}

          {node.icon && <span className="text-base">{node.icon}</span>}
          {node.image_url && <img src={node.image_url} alt={node.name_en} className="w-4 h-4 object-contain" />}

          <div className="flex-1 min-w-0">
            <span className={`font-semibold text-sm ${!node.is_active ? 'text-zinc-500 line-through' : 'text-white'}`}>
              {node.name_en}
            </span>
            <span className="ml-2 text-xs text-zinc-500 font-mono">/{node.slug}</span>
            {node.is_leaf && <span className="ml-2 text-[10px] bg-emerald-900/30 text-emerald-400 border border-emerald-500/20 rounded px-1.5 py-0.5 font-bold">LEAF</span>}
            {!node.is_active && <span className="ml-2 text-[10px] bg-red-900/30 text-red-400 border border-red-500/20 rounded px-1.5 py-0.5 font-bold">INACTIVE</span>}
          </div>

          <div className="flex items-center gap-1.5 opacity-0 group-hover:opacity-100 transition-opacity">
            {canManage && (
              <>
                <button
                  onClick={() => handleToggleStatus(node)}
                  className={`p-1.5 rounded-lg transition-colors ${node.is_active ? 'bg-zinc-700 hover:bg-zinc-600 text-zinc-300' : 'bg-red-900/40 text-red-400 hover:bg-red-800/50'}`}
                  title={node.is_active ? 'Deactivate category' : 'Activate category'}
                >
                  {node.is_active ? <Eye className="w-3.5 h-3.5" /> : <EyeOff className="w-3.5 h-3.5" />}
                </button>
                <button
                  onClick={() => openCreate(node.id)}
                  className="p-1.5 bg-zinc-700 hover:bg-zinc-600 rounded-lg text-zinc-400 hover:text-white transition-colors"
                  title="Add child category"
                >
                  <Plus className="w-3.5 h-3.5" />
                </button>
                <button
                  onClick={() => openEdit(node)}
                  className="p-1.5 bg-zinc-700 hover:bg-blue-600/30 rounded-lg text-zinc-400 hover:text-blue-400 transition-colors"
                  title="Edit category"
                >
                  <Pencil className="w-3.5 h-3.5" />
                </button>
                <button
                  onClick={() => setDeleteConfirm(node.id)}
                  className="p-1.5 bg-zinc-700 hover:bg-red-600/30 rounded-lg text-zinc-400 hover:text-red-400 transition-colors"
                  title="Delete category (safe check)"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                </button>
              </>
            )}
          </div>
        </div>
        {open && hasChildren && node.children!.map(child => renderNode(child, depth + 1))}
      </div>
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-amber-600/10 rounded-xl border border-amber-500/20">
            <FolderTree className="w-6 h-6 text-amber-400" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Database Category Hierarchy</h1>
            <p className="text-xs text-zinc-400">Canonical PostgreSQL taxonomy tree with multi-language & media support</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <button onClick={fetchTree} className="p-2 bg-zinc-800 rounded-lg hover:bg-zinc-700 transition-colors">
            <RefreshCw className="w-4 h-4 text-zinc-400" />
          </button>
          {canManage && (
            <button
              onClick={() => openCreate()}
              className="flex items-center gap-2 px-4 py-2 bg-orange-600 hover:bg-orange-700 text-white rounded-xl text-sm font-bold transition-colors"
            >
              <Plus className="w-4 h-4" /> Add Root Category
            </button>
          )}
        </div>
      </div>

      <div className="bg-zinc-900 rounded-2xl border border-zinc-800 p-4">
        {loading ? (
          <div className="p-8 text-center text-zinc-500">
            <div className="animate-spin w-8 h-8 border-2 border-orange-500 border-t-transparent rounded-full mx-auto mb-3" />
            Loading categories from canonical database...
          </div>
        ) : tree.length === 0 ? (
          <div className="p-8 text-center text-zinc-500">
            <FolderTree className="w-12 h-12 mx-auto mb-3 opacity-20" />
            <p>No categories found in taxonomy.</p>
            {canManage && (
              <button onClick={() => openCreate()} className="mt-4 px-4 py-2 bg-orange-600 text-white rounded-lg text-sm font-bold">
                Create First Category
              </button>
            )}
          </div>
        ) : (
          <div className="space-y-1">{tree.map(node => renderNode(node))}</div>
        )}
      </div>

      {/* Create/Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4">
          <div className="bg-zinc-900 border border-zinc-700 rounded-2xl w-full max-w-lg shadow-2xl max-h-[90vh] overflow-y-auto">
            <div className="px-6 py-5 border-b border-zinc-800 flex items-center justify-between sticky top-0 bg-zinc-900 z-10">
              <h2 className="text-lg font-black text-white">{editItem ? 'Edit Category' : 'New Category'}</h2>
              <button onClick={() => setShowModal(false)}><X className="w-5 h-5 text-zinc-500 hover:text-white" /></button>
            </div>
            <div className="p-6 space-y-4">
              <div className="grid grid-cols-1 gap-4">
                <div>
                  <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Name (English) *</label>
                  <input
                    value={form.name_en}
                    onChange={e => {
                      const val = e.target.value;
                      setForm(f => ({ ...f, name_en: val, slug: editItem ? f.slug : autoSlug(val) }));
                    }}
                    onBlur={handleSlugBlur}
                    className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-2.5 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                    placeholder="Electronics"
                  />
                </div>
                <div className="grid grid-cols-3 gap-3">
                  <div>
                    <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Amharic (አማርኛ)</label>
                    <input value={form.name_am} onChange={e => setForm(f => ({ ...f, name_am: e.target.value }))}
                      className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-3 py-2 text-xs text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                      placeholder="ኤሌክትሮኒክስ" />
                  </div>
                  <div>
                    <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Oromo (Oromoo)</label>
                    <input value={form.name_om} onChange={e => setForm(f => ({ ...f, name_om: e.target.value }))}
                      className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-3 py-2 text-xs text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                      placeholder="Elektirooniksii" />
                  </div>
                  <div>
                    <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Tigrinya (ትግርኛ)</label>
                    <input value={form.name_ti} onChange={e => setForm(f => ({ ...f, name_ti: e.target.value }))}
                      className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-3 py-2 text-xs text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                      placeholder="ኤለክትሮኒክስ" />
                  </div>
                </div>

                <div>
                  <div className="flex items-center justify-between mb-1.5">
                    <label className="block text-xs font-bold text-zinc-400 uppercase">Slug (Canonical URL Identifier) *</label>
                    {slugChecking && <span className="text-[10px] text-zinc-400">Checking availability...</span>}
                    {!slugChecking && slugAvailable === true && <span className="text-[10px] text-green-400 flex items-center gap-1"><CheckCircle2 size={12} /> Available</span>}
                    {!slugChecking && slugAvailable === false && <span className="text-[10px] text-red-400 flex items-center gap-1"><AlertCircle size={12} /> Already in use</span>}
                  </div>
                  <input
                    value={form.slug}
                    onChange={e => {
                      setForm(f => ({ ...f, slug: autoSlug(e.target.value) }));
                      setSlugAvailable(null);
                    }}
                    onBlur={handleSlugBlur}
                    className={`w-full bg-zinc-800 border rounded-xl px-4 py-2.5 text-sm text-white font-mono placeholder-zinc-500 focus:outline-none ${slugAvailable === false ? 'border-red-500' : 'border-zinc-700 focus:border-orange-500/50'}`}
                    placeholder="electronics"
                  />
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Icon (Emoji / Glyph)</label>
                    <input value={form.icon} onChange={e => setForm(f => ({ ...f, icon: e.target.value }))}
                      className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-3 py-2 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                      placeholder="📱 or 🚗" />
                  </div>
                  <div>
                    <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Image URL</label>
                    <input value={form.image_url} onChange={e => setForm(f => ({ ...f, image_url: e.target.value }))}
                      className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-3 py-2 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-orange-500/50"
                      placeholder="https://... or /uploads/..." />
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-xs font-bold text-zinc-400 uppercase mb-1.5">Display Order</label>
                    <input type="number" value={form.display_order} onChange={e => setForm(f => ({ ...f, display_order: parseInt(e.target.value) || 0 }))}
                      className="w-full bg-zinc-800 border border-zinc-700 rounded-xl px-4 py-2.5 text-sm text-white focus:outline-none focus:border-orange-500/50" />
                  </div>
                  <div className="flex flex-col gap-3 pt-5">
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input type="checkbox" checked={form.is_leaf} onChange={e => setForm(f => ({ ...f, is_leaf: e.target.checked }))}
                        className="rounded border-zinc-600 bg-zinc-800 text-orange-600 focus:ring-orange-500" />
                      <span className="text-xs text-zinc-300">Is Leaf (products attach here)</span>
                    </label>
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input type="checkbox" checked={form.is_active} onChange={e => setForm(f => ({ ...f, is_active: e.target.checked }))}
                        className="rounded border-zinc-600 bg-zinc-800 text-orange-600 focus:ring-orange-500" />
                      <span className="text-xs text-zinc-300">Active Status</span>
                    </label>
                  </div>
                </div>
              </div>

              {error && <p className="text-red-400 text-sm">{error}</p>}

              <div className="flex gap-3 pt-2">
                <button onClick={() => setShowModal(false)} className="flex-1 py-3 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-sm font-bold transition-colors">
                  Cancel
                </button>
                <button onClick={handleSubmit} disabled={submitting}
                  className="flex-1 py-3 bg-orange-600 hover:bg-orange-700 text-white rounded-xl text-sm font-bold disabled:opacity-50 transition-colors flex items-center justify-center gap-2">
                  <Save className="w-4 h-4" /> {submitting ? 'Saving...' : editItem ? 'Update Category' : 'Create Category'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Delete confirmation */}
      {deleteConfirm !== null && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4">
          <div className="bg-zinc-900 border border-zinc-700 rounded-2xl p-6 max-w-sm w-full shadow-2xl">
            <h3 className="text-lg font-black text-white mb-2">Delete Category?</h3>
            <p className="text-zinc-400 text-sm mb-6">Permanent deletion is blocked if listings or child categories reference it. Consider deactivating instead.</p>
            <div className="flex gap-3">
              <button onClick={() => setDeleteConfirm(null)} className="flex-1 py-2.5 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-sm font-bold">Cancel</button>
              <button onClick={() => handleDelete(deleteConfirm)} className="flex-1 py-2.5 bg-red-600 hover:bg-red-700 text-white rounded-xl text-sm font-bold">Delete</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
