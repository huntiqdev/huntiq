'use client';

import { useEffect, useState, useMemo } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { 
  MapPin, Plus, Pencil, Trash2, ChevronRight, ChevronDown, 
  RefreshCw, Globe, Layers, Navigation, Upload, CheckCircle2, 
  AlertCircle, Search, ExternalLink, X, Save, Compass
} from 'lucide-react';

interface LocationItem {
  id: number;
  parent_id?: number | null;
  location_type: string;
  name_en: string;
  name_am?: string;
  name_om?: string;
  name_ti?: string;
  slug: string;
  admin_code?: string;
  level: number;
  path: string;
  display_order: number;
  is_active: boolean;
  is_leaf: boolean;
  latitude?: number | null;
  longitude?: number | null;
  min_lat?: number | null;
  max_lat?: number | null;
  min_lng?: number | null;
  max_lng?: number | null;
  aliases?: string[];
  source?: string;
  verification_status?: string;
}

interface LocationTreeNode {
  id: number;
  name: string;
  slug: string;
  location_type: string;
  admin_code?: string;
  level: number;
  latitude?: number | null;
  longitude?: number | null;
  is_leaf: boolean;
  children?: LocationTreeNode[];
}

const LOCATION_TYPES = [
  { value: 'country', label: 'Country (Level 0)', level: 0 },
  { value: 'region', label: 'Region / Chartered City (Level 1)', level: 1 },
  { value: 'zone', label: 'Zone (Level 2)', level: 2 },
  { value: 'city', label: 'City / Municipality (Level 2/3)', level: 2 },
  { value: 'subcity', label: 'Subcity (Level 2/3)', level: 2 },
  { value: 'woreda', label: 'Woreda / District (Level 3/4)', level: 3 },
  { value: 'kebele', label: 'Kebele (Level 4/5)', level: 4 },
];

const EMPTY_FORM = {
  parent_id: null as number | null,
  location_type: 'region',
  name_en: '',
  name_am: '',
  name_om: '',
  name_ti: '',
  slug: '',
  admin_code: '',
  display_order: 0,
  is_active: true,
  is_leaf: false,
  latitude: '' as string | number,
  longitude: '' as string | number,
  min_lat: '' as string | number,
  max_lat: '' as string | number,
  min_lng: '' as string | number,
  max_lng: '' as string | number,
  aliases_str: '',
};

export default function AdminLocationsPage() {
  const { user, isSuperAdmin, hasAdminPermission } = useAuth();
  const token = user?.token || (typeof window !== 'undefined' ? localStorage.getItem('token') || '' : '');

  // Tabs: 'tree' | 'list' | 'geocoding' | 'import'
  const [activeTab, setActiveTab] = useState<'tree' | 'list' | 'geocoding' | 'import'>('tree');
  
  // Data states
  const [tree, setTree] = useState<LocationTreeNode[]>([]);
  const [flatList, setFlatList] = useState<LocationItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [levelFilter, setLevelFilter] = useState<string>('all');
  const [expanded, setExpanded] = useState<Set<number>>(new Set([1, 2])); // Auto expand Ethiopia & Addis Ababa
  
  // Modal states
  const [showModal, setShowModal] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [form, setForm] = useState({ ...EMPTY_FORM });
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState('');
  const [deleteConfirmId, setDeleteConfirmId] = useState<number | null>(null);

  // Geocoding sandbox states
  const [geoLat, setGeoLat] = useState('8.9950');
  const [geoLng, setGeoLng] = useState('38.7900');
  const [geoTesting, setGeoTesting] = useState(false);
  const [geoResult, setGeoResult] = useState<any>(null);
  const [geoError, setGeoError] = useState('');

  // Batch import states
  const [importJson, setImportJson] = useState('');
  const [importing, setImporting] = useState(false);
  const [importResult, setImportResult] = useState<any>(null);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      // 1. Fetch hierarchical tree
      const treeRes = await api.getLocationTree('en');
      setTree(treeRes || []);

      // 2. Fetch flat admin list
      const listRes = await api.adminGetLocationHierarchy({}, token);
      setFlatList(listRes || []);
    } catch (err: any) {
      console.error('Failed to load locations', err);
    } finally {
      setLoading(false);
    }
  };

  const toggleExpand = (id: number) => {
    setExpanded(prev => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const openCreateModal = (parentId?: number) => {
    setEditingId(null);
    let defaultType = 'region';
    if (parentId) {
      const parent = flatList.find(l => l.id === parentId);
      if (parent) {
        if (parent.location_type === 'country') defaultType = 'region';
        else if (parent.location_type === 'region') defaultType = 'subcity';
        else if (parent.location_type === 'subcity' || parent.location_type === 'zone') defaultType = 'woreda';
        else if (parent.location_type === 'woreda') defaultType = 'kebele';
      }
    }
    setForm({
      ...EMPTY_FORM,
      parent_id: parentId ?? 1,
      location_type: defaultType,
    });
    setFormError('');
    setShowModal(true);
  };

  const openEditModal = (item: LocationItem) => {
    setEditingId(item.id);
    setForm({
      parent_id: item.parent_id ?? null,
      location_type: item.location_type,
      name_en: item.name_en,
      name_am: item.name_am || '',
      name_om: item.name_om || '',
      name_ti: item.name_ti || '',
      slug: item.slug,
      admin_code: item.admin_code || '',
      display_order: item.display_order,
      is_active: item.is_active,
      is_leaf: item.is_leaf,
      latitude: item.latitude !== null && item.latitude !== undefined ? item.latitude : '',
      longitude: item.longitude !== null && item.longitude !== undefined ? item.longitude : '',
      min_lat: item.min_lat !== null && item.min_lat !== undefined ? item.min_lat : '',
      max_lat: item.max_lat !== null && item.max_lat !== undefined ? item.max_lat : '',
      min_lng: item.min_lng !== null && item.min_lng !== undefined ? item.min_lng : '',
      max_lng: item.max_lng !== null && item.max_lng !== undefined ? item.max_lng : '',
      aliases_str: (item.aliases || []).join(', '),
    });
    setFormError('');
    setShowModal(true);
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name_en.trim()) {
      setFormError('English name is required');
      return;
    }
    setSubmitting(true);
    setFormError('');

    try {
      const aliases = form.aliases_str
        ? form.aliases_str.split(',').map(s => s.trim()).filter(Boolean)
        : [];

      const payload: any = {
        parent_id: form.parent_id ? Number(form.parent_id) : null,
        location_type: form.location_type,
        name_en: form.name_en.trim(),
        name_am: form.name_am.trim(),
        name_om: form.name_om.trim(),
        name_ti: form.name_ti.trim(),
        slug: form.slug.trim(),
        admin_code: form.admin_code.trim(),
        display_order: Number(form.display_order) || 0,
        is_active: Boolean(form.is_active),
        is_leaf: Boolean(form.is_leaf),
        aliases,
      };

      if (form.latitude !== '') payload.latitude = parseFloat(String(form.latitude));
      if (form.longitude !== '') payload.longitude = parseFloat(String(form.longitude));
      if (form.min_lat !== '') payload.min_lat = parseFloat(String(form.min_lat));
      if (form.max_lat !== '') payload.max_lat = parseFloat(String(form.max_lat));
      if (form.min_lng !== '') payload.min_lng = parseFloat(String(form.min_lng));
      if (form.max_lng !== '') payload.max_lng = parseFloat(String(form.max_lng));

      if (editingId) {
        await api.adminUpdateLocationHierarchy(editingId, payload, token);
      } else {
        await api.adminCreateLocationHierarchy(payload, token);
      }

      setShowModal(false);
      fetchData();
    } catch (err: any) {
      setFormError(err?.message || 'Failed to save location');
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await api.adminDeleteLocationHierarchy(id, token);
      setDeleteConfirmId(null);
      fetchData();
    } catch (err: any) {
      alert(err?.message || 'Failed to delete location');
    }
  };

  const runTestGeocode = async (lat?: number, lng?: number) => {
    const testLat = lat !== undefined ? lat : parseFloat(geoLat);
    const testLng = lng !== undefined ? lng : parseFloat(geoLng);
    if (isNaN(testLat) || isNaN(testLng)) {
      setGeoError('Please enter valid numeric latitude and longitude');
      return;
    }
    setGeoTesting(true);
    setGeoError('');
    setGeoResult(null);

    try {
      const res = await api.adminTestGeocode(testLat, testLng, 'en', token);
      setGeoResult(res);
    } catch (err: any) {
      setGeoError(err?.message || 'Geocoding request failed');
    } finally {
      setGeoTesting(false);
    }
  };

  const handleBatchImport = async () => {
    if (!importJson.trim()) return;
    setImporting(true);
    setImportResult(null);

    try {
      const parsed = JSON.parse(importJson);
      const items = Array.isArray(parsed) ? parsed : (parsed.items || [parsed]);
      const res = await api.adminImportLocationHierarchy({ items }, token);
      setImportResult(res?.data || res);
      fetchData();
    } catch (err: any) {
      setImportResult({ error: err?.message || 'JSON Parse or Import error' });
    } finally {
      setImporting(false);
    }
  };

  const loadImportTemplate = () => {
    const template = [
      {
        parent_id: 2,
        location_type: 'woreda',
        name_en: 'Woreda 03',
        name_am: 'ወረዳ 03',
        name_om: 'Aanaa 03',
        slug: 'woreda-03',
        admin_code: 'ET-AA-BOL-W03',
        is_active: true,
        is_leaf: false,
        latitude: 8.9912,
        longitude: 38.7891,
        min_lat: 8.9800,
        max_lat: 9.0050,
        min_lng: 38.7750,
        max_lng: 38.8050,
        aliases: ['Bole Medhanealem Woreda', 'Woreda 3']
      }
    ];
    setImportJson(JSON.stringify(template, null, 2));
  };

  // Filtered flat list
  const filteredList = useMemo(() => {
    return flatList.filter(item => {
      const matchesSearch = !searchQuery || 
        item.name_en.toLowerCase().includes(searchQuery.toLowerCase()) ||
        (item.name_am && item.name_am.includes(searchQuery)) ||
        (item.admin_code && item.admin_code.toLowerCase().includes(searchQuery.toLowerCase())) ||
        (item.aliases && item.aliases.some(a => a.toLowerCase().includes(searchQuery.toLowerCase())));
      
      const matchesLevel = levelFilter === 'all' || String(item.level) === levelFilter;
      return matchesSearch && matchesLevel;
    });
  }, [flatList, searchQuery, levelFilter]);

  // Recursive Tree Node Renderer
  const renderTreeNode = (node: LocationTreeNode, depth = 0) => {
    const hasChildren = node.children && node.children.length > 0;
    const isExpanded = expanded.has(node.id);
    const itemDetail = flatList.find(f => f.id === node.id);

    return (
      <div key={node.id} className="select-none text-xs">
        <div 
          className={`flex items-center justify-between py-2 px-3 rounded-xl transition-colors hover:bg-zinc-900/60 ${
            depth === 0 ? 'bg-zinc-900/40 font-bold text-white mb-1 border border-zinc-800' : 'text-zinc-300'
          }`}
          style={{ paddingLeft: `${depth * 20 + 12}px` }}
        >
          <div className="flex items-center gap-2.5">
            {hasChildren ? (
              <button 
                onClick={() => toggleExpand(node.id)}
                className="p-1 hover:bg-zinc-800 rounded-md text-zinc-400 hover:text-white"
              >
                {isExpanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
              </button>
            ) : (
              <span className="w-5 text-center text-zinc-600">•</span>
            )}

            <span className={`inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider ${
              node.level === 0 ? 'bg-purple-900/40 text-purple-300 border border-purple-700/50' :
              node.level === 1 ? 'bg-blue-900/40 text-blue-300 border border-blue-700/50' :
              node.level === 2 ? 'bg-orange-900/40 text-orange-300 border border-orange-700/50' :
              node.level === 3 ? 'bg-emerald-900/40 text-emerald-300 border border-emerald-700/50' :
              'bg-zinc-800 text-zinc-400'
            }`}>
              {node.location_type}
            </span>

            <span className="font-semibold text-white">{node.name}</span>

            {node.admin_code && (
              <span className="text-[10px] font-mono text-zinc-500 bg-zinc-950 px-1.5 py-0.5 rounded border border-zinc-800">
                {node.admin_code}
              </span>
            )}

            {node.latitude && node.longitude && (
              <span className="text-[10px] text-zinc-500 flex items-center gap-1">
                <Navigation size={10} className="text-orange-400" />
                {node.latitude.toFixed(3)}, {node.longitude.toFixed(3)}
              </span>
            )}
          </div>

          <div className="flex items-center gap-1.5 opacity-80 hover:opacity-100">
            <button
              onClick={() => openCreateModal(node.id)}
              title="Add Child Location"
              className="p-1 hover:bg-zinc-800 text-zinc-400 hover:text-emerald-400 rounded-lg transition-colors"
            >
              <Plus size={14} />
            </button>
            {itemDetail && (
              <button
                onClick={() => openEditModal(itemDetail)}
                title="Edit Location"
                className="p-1 hover:bg-zinc-800 text-zinc-400 hover:text-orange-400 rounded-lg transition-colors"
              >
                <Pencil size={14} />
              </button>
            )}
            {node.level > 0 && (
              <button
                onClick={() => setDeleteConfirmId(node.id)}
                title="Delete Location"
                className="p-1 hover:bg-zinc-800 text-zinc-400 hover:text-red-400 rounded-lg transition-colors"
              >
                <Trash2 size={14} />
              </button>
            )}
          </div>
        </div>

        {hasChildren && isExpanded && (
          <div className="space-y-0.5 my-0.5">
            {node.children!.map(child => renderTreeNode(child, depth + 1))}
          </div>
        )}
      </div>
    );
  };

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100 p-6 space-y-6">
      {/* Top Banner & Stats */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-6 border-b border-zinc-900">
        <div>
          <div className="flex items-center gap-2.5">
            <div className="p-2.5 bg-orange-600/10 text-orange-500 rounded-2xl border border-orange-500/20">
              <MapPin size={26} />
            </div>
            <div>
              <h1 className="text-2xl font-black tracking-tight text-white">Location Hierarchy & GPS Automation</h1>
              <p className="text-xs text-zinc-400 mt-0.5">
                Data-driven Ethiopian administrative levels, reverse geocoding engine, spatial bounding boxes & listing autofill.
              </p>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-2.5">
          <button
            onClick={fetchData}
            disabled={loading}
            className="flex items-center gap-2 px-3 py-2 bg-zinc-900 hover:bg-zinc-800 text-zinc-300 text-xs font-bold rounded-xl border border-zinc-800 transition-all"
          >
            <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
            Sync
          </button>
          <button
            onClick={() => openCreateModal()}
            className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-500 hover:to-amber-500 text-white text-xs font-bold rounded-xl shadow-lg shadow-orange-600/20 transition-all"
          >
            <Plus size={15} />
            Add Location
          </button>
        </div>
      </div>

      {/* Metric Cards */}
      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
        <div className="bg-zinc-900/60 border border-zinc-800/80 rounded-2xl p-3.5">
          <p className="text-[10px] font-bold text-zinc-500 uppercase tracking-wider">Total Nodes</p>
          <p className="text-xl font-black text-white mt-1">{flatList.length}</p>
        </div>
        <div className="bg-zinc-900/60 border border-zinc-800/80 rounded-2xl p-3.5">
          <p className="text-[10px] font-bold text-blue-400 uppercase tracking-wider">Regions & Cities</p>
          <p className="text-xl font-black text-white mt-1">
            {flatList.filter(l => l.level === 1).length}
          </p>
        </div>
        <div className="bg-zinc-900/60 border border-zinc-800/80 rounded-2xl p-3.5">
          <p className="text-[10px] font-bold text-orange-400 uppercase tracking-wider">Zones & Subcities</p>
          <p className="text-xl font-black text-white mt-1">
            {flatList.filter(l => l.level === 2).length}
          </p>
        </div>
        <div className="bg-zinc-900/60 border border-zinc-800/80 rounded-2xl p-3.5">
          <p className="text-[10px] font-bold text-emerald-400 uppercase tracking-wider">Woredas</p>
          <p className="text-xl font-black text-white mt-1">
            {flatList.filter(l => l.level === 3).length}
          </p>
        </div>
        <div className="bg-zinc-900/60 border border-zinc-800/80 rounded-2xl p-3.5">
          <p className="text-[10px] font-bold text-purple-400 uppercase tracking-wider">GPS Coordinates</p>
          <p className="text-xl font-black text-white mt-1">
            {flatList.filter(l => l.latitude !== null && l.latitude !== undefined).length}
          </p>
        </div>
        <div className="bg-zinc-900/60 border border-zinc-800/80 rounded-2xl p-3.5">
          <p className="text-[10px] font-bold text-amber-400 uppercase tracking-wider">Bounding Boxes</p>
          <p className="text-xl font-black text-white mt-1">
            {flatList.filter(l => l.min_lat !== null && l.min_lat !== undefined).length}
          </p>
        </div>
      </div>

      {/* Navigation Tabs */}
      <div className="flex items-center gap-2 border-b border-zinc-900 pb-3">
        <button
          onClick={() => setActiveTab('tree')}
          className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all ${
            activeTab === 'tree'
              ? 'bg-orange-600/15 text-orange-400 border border-orange-500/20'
              : 'text-zinc-400 hover:bg-zinc-900 hover:text-white'
          }`}
        >
          <Layers size={15} />
          Hierarchy Tree
        </button>
        <button
          onClick={() => setActiveTab('list')}
          className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all ${
            activeTab === 'list'
              ? 'bg-orange-600/15 text-orange-400 border border-orange-500/20'
              : 'text-zinc-400 hover:bg-zinc-900 hover:text-white'
          }`}
        >
          <Globe size={15} />
          All Records ({flatList.length})
        </button>
        <button
          onClick={() => setActiveTab('geocoding')}
          className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all ${
            activeTab === 'geocoding'
              ? 'bg-orange-600/15 text-orange-400 border border-orange-500/20'
              : 'text-zinc-400 hover:bg-zinc-900 hover:text-white'
          }`}
        >
          <Compass size={15} />
          Geocoding Sandbox
        </button>
        <button
          onClick={() => setActiveTab('import')}
          className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all ${
            activeTab === 'import'
              ? 'bg-orange-600/15 text-orange-400 border border-orange-500/20'
              : 'text-zinc-400 hover:bg-zinc-900 hover:text-white'
          }`}
        >
          <Upload size={15} />
          Batch Import Tool
        </button>
      </div>

      {/* Tab 1: Hierarchy Tree Explorer */}
      {activeTab === 'tree' && (
        <div className="bg-zinc-900/40 border border-zinc-900 rounded-3xl p-5 space-y-4">
          <div className="flex items-center justify-between gap-4">
            <div>
              <h2 className="text-sm font-bold text-white">Interactive Administrative Tree</h2>
              <p className="text-[11px] text-zinc-400">
                Data-driven cascading nodes. Click arrow icons to expand child districts and kebeles.
              </p>
            </div>
            <div className="flex items-center gap-2">
              <button
                onClick={() => setExpanded(new Set(flatList.map(l => l.id)))}
                className="text-[11px] text-zinc-400 hover:text-white px-2.5 py-1 bg-zinc-900 rounded-lg border border-zinc-800"
              >
                Expand All
              </button>
              <button
                onClick={() => setExpanded(new Set([1]))}
                className="text-[11px] text-zinc-400 hover:text-white px-2.5 py-1 bg-zinc-900 rounded-lg border border-zinc-800"
              >
                Collapse All
              </button>
            </div>
          </div>

          <div className="bg-zinc-950/80 rounded-2xl p-4 border border-zinc-900/80 max-h-[600px] overflow-y-auto space-y-1">
            {tree.length === 0 ? (
              <div className="text-center py-12 text-zinc-500 text-xs">
                No location nodes found. Click "Add Location" to create Ethiopia root.
              </div>
            ) : (
              tree.map(node => renderTreeNode(node))
            )}
          </div>
        </div>
      )}

      {/* Tab 2: All Records Table & Filtering */}
      {activeTab === 'list' && (
        <div className="space-y-4">
          <div className="flex flex-col sm:flex-row items-center gap-3">
            <div className="relative flex-1 w-full">
              <Search size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-500" />
              <input
                type="text"
                value={searchQuery}
                onChange={e => setSearchQuery(e.target.value)}
                placeholder="Search by name in EN/AM/OM, slug, admin code, or alias..."
                className="w-full bg-zinc-900/70 border border-zinc-800 rounded-xl pl-9 pr-4 py-2 text-xs text-white placeholder-zinc-500 outline-none focus:ring-1 focus:ring-orange-500"
              />
            </div>
            <select
              value={levelFilter}
              onChange={e => setLevelFilter(e.target.value)}
              className="w-full sm:w-48 bg-zinc-900/70 border border-zinc-800 rounded-xl px-3 py-2 text-xs text-zinc-300 outline-none focus:ring-1 focus:ring-orange-500"
            >
              <option value="all">All Administrative Levels</option>
              <option value="0">Level 0 — Country</option>
              <option value="1">Level 1 — Region / City</option>
              <option value="2">Level 2 — Zone / Subcity</option>
              <option value="3">Level 3 — Woreda</option>
              <option value="4">Level 4 — Kebele</option>
            </select>
          </div>

          <div className="bg-zinc-900/40 border border-zinc-900 rounded-2xl overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead className="bg-zinc-950 text-zinc-400 font-extrabold uppercase tracking-wider text-[10px] border-b border-zinc-900">
                  <tr>
                    <th className="py-3 px-4">ID</th>
                    <th className="py-3 px-4">Name (English / Localized)</th>
                    <th className="py-3 px-4">Type & Level</th>
                    <th className="py-3 px-4">Admin Code</th>
                    <th className="py-3 px-4">Materialized Path</th>
                    <th className="py-3 px-4">GPS & Bounds</th>
                    <th className="py-3 px-4">Status</th>
                    <th className="py-3 px-4 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-zinc-900">
                  {filteredList.map(item => (
                    <tr key={item.id} className="hover:bg-zinc-900/50 transition-colors">
                      <td className="py-3 px-4 font-mono text-zinc-500">{item.id}</td>
                      <td className="py-3 px-4">
                        <div className="font-bold text-white">{item.name_en}</div>
                        {(item.name_am || item.name_om || item.name_ti) && (
                          <div className="text-[11px] text-zinc-400 mt-0.5">
                            {[item.name_am, item.name_om, item.name_ti].filter(Boolean).join(' • ')}
                          </div>
                        )}
                      </td>
                      <td className="py-3 px-4">
                        <span className="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider bg-zinc-800 text-zinc-300">
                          {item.location_type} (L{item.level})
                        </span>
                      </td>
                      <td className="py-3 px-4 font-mono text-zinc-400">
                        {item.admin_code || '—'}
                      </td>
                      <td className="py-3 px-4 font-mono text-[10px] text-zinc-500">
                        {item.path}
                      </td>
                      <td className="py-3 px-4 text-zinc-400">
                        {item.latitude && item.longitude ? (
                          <div className="flex items-center gap-1 text-[11px]">
                            <Navigation size={11} className="text-orange-400" />
                            {item.latitude.toFixed(2)}, {item.longitude.toFixed(2)}
                            {item.min_lat && <span className="text-[10px] text-amber-500">(BBox)</span>}
                          </div>
                        ) : (
                          <span className="text-zinc-600">—</span>
                        )}
                      </td>
                      <td className="py-3 px-4">
                        <span className={`inline-flex px-2 py-0.5 rounded text-[10px] font-bold ${
                          item.is_active ? 'bg-emerald-950 text-emerald-400 border border-emerald-800/50' : 'bg-red-950 text-red-400 border border-red-800/50'
                        }`}>
                          {item.is_active ? 'Active' : 'Inactive'}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-right">
                        <div className="flex items-center justify-end gap-1.5">
                          <button
                            onClick={() => openEditModal(item)}
                            className="p-1.5 hover:bg-zinc-800 text-zinc-400 hover:text-orange-400 rounded-lg transition-colors"
                          >
                            <Pencil size={13} />
                          </button>
                          {item.level > 0 && (
                            <button
                              onClick={() => setDeleteConfirmId(item.id)}
                              className="p-1.5 hover:bg-zinc-800 text-zinc-400 hover:text-red-400 rounded-lg transition-colors"
                            >
                              <Trash2 size={13} />
                            </button>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {/* Tab 3: Geocoding Sandbox */}
      {activeTab === 'geocoding' && (
        <div className="bg-zinc-900/40 border border-zinc-900 rounded-3xl p-6 space-y-6">
          <div className="max-w-xl">
            <h2 className="text-base font-bold text-white flex items-center gap-2">
              <Compass size={18} className="text-orange-500" />
              GPS Reverse Geocoding Sandbox
            </h2>
            <p className="text-xs text-zinc-400 mt-1">
              Test how the backend resolves GPS coordinates into verified administrative boundaries (Region › Subcity › Woreda) with privacy-preserving local spatial bounding boxes.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-zinc-950/80 border border-zinc-800/80 rounded-2xl p-5 space-y-4">
              <h3 className="text-xs font-bold text-zinc-300 uppercase tracking-wider">Test Coordinate Input</h3>
              
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="text-[11px] font-bold text-zinc-400">Latitude</label>
                  <input
                    type="text"
                    value={geoLat}
                    onChange={e => setGeoLat(e.target.value)}
                    placeholder="8.9950"
                    className="w-full mt-1 bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-xs text-white outline-none focus:ring-1 focus:ring-orange-500"
                  />
                </div>
                <div>
                  <label className="text-[11px] font-bold text-zinc-400">Longitude</label>
                  <input
                    type="text"
                    value={geoLng}
                    onChange={e => setGeoLng(e.target.value)}
                    placeholder="38.7900"
                    className="w-full mt-1 bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-xs text-white outline-none focus:ring-1 focus:ring-orange-500"
                  />
                </div>
              </div>

              <div>
                <p className="text-[11px] font-bold text-zinc-400 mb-2">Preset Ethiopian Locations:</p>
                <div className="flex flex-wrap gap-1.5">
                  {[
                    { name: 'Bole (Addis Ababa)', lat: 8.995, lng: 38.790 },
                    { name: 'Piassa / Arada', lat: 9.034, lng: 38.752 },
                    { name: 'Adama City', lat: 8.540, lng: 39.270 },
                    { name: 'Hawassa', lat: 7.050, lng: 38.470 },
                    { name: 'Bahir Dar', lat: 11.590, lng: 37.390 },
                    { name: 'Mekelle', lat: 13.490, lng: 39.470 },
                  ].map((p, idx) => (
                    <button
                      key={idx}
                      type="button"
                      onClick={() => {
                        setGeoLat(String(p.lat));
                        setGeoLng(String(p.lng));
                        runTestGeocode(p.lat, p.lng);
                      }}
                      className="px-2.5 py-1 bg-zinc-900 hover:bg-zinc-800 text-zinc-300 text-[10px] font-medium rounded-lg border border-zinc-800 transition-colors"
                    >
                      📍 {p.name}
                    </button>
                  ))}
                </div>
              </div>

              {geoError && (
                <div className="p-3 bg-red-950/40 border border-red-800/50 rounded-xl text-xs text-red-400 flex items-center gap-2">
                  <AlertCircle size={14} />
                  {geoError}
                </div>
              )}

              <button
                type="button"
                onClick={() => runTestGeocode()}
                disabled={geoTesting}
                className="w-full py-2.5 bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-500 hover:to-amber-500 text-white font-bold text-xs rounded-xl shadow-lg transition-all flex items-center justify-center gap-2"
              >
                {geoTesting ? <RefreshCw size={14} className="animate-spin" /> : <Navigation size={14} />}
                Resolve Coordinates
              </button>
            </div>

            <div className="bg-zinc-950/80 border border-zinc-800/80 rounded-2xl p-5 space-y-4">
              <h3 className="text-xs font-bold text-zinc-300 uppercase tracking-wider">Engine Resolution Output</h3>
              
              {!geoResult && !geoTesting && (
                <div className="text-center py-12 text-zinc-500 text-xs">
                  Click "Resolve Coordinates" or choose a preset location to run geocoding.
                </div>
              )}

              {geoTesting && (
                <div className="text-center py-12 text-orange-400 text-xs flex flex-col items-center gap-2">
                  <RefreshCw size={20} className="animate-spin" />
                  Resolving coordinates against spatial database...
                </div>
              )}

              {geoResult && (
                <div className="space-y-4 text-xs">
                  <div className="flex items-center justify-between p-3 bg-zinc-900/60 rounded-xl border border-zinc-800">
                    <div>
                      <p className="text-[10px] font-bold text-zinc-500 uppercase">Matched Location</p>
                      <p className="text-sm font-black text-white mt-0.5">{geoResult.name}</p>
                      <p className="text-[11px] text-orange-400 capitalize">{geoResult.location_type} (Level {geoResult.level})</p>
                    </div>
                    <div className="text-right">
                      <p className="text-[10px] font-bold text-zinc-500 uppercase">Confidence</p>
                      <div className="flex items-center gap-1.5 mt-0.5 justify-end">
                        <span className="font-mono font-bold text-emerald-400">{Math.round((geoResult.confidence || 0) * 100)}%</span>
                      </div>
                      <span className="text-[10px] text-zinc-500 font-mono">Source: {geoResult.source}</span>
                    </div>
                  </div>

                  {geoResult.breadcrumbs && geoResult.breadcrumbs.length > 0 && (
                    <div>
                      <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-wider mb-1.5">Administrative Breadcrumbs</p>
                      <div className="flex flex-wrap items-center gap-1.5 p-3 bg-zinc-900/40 rounded-xl border border-zinc-800">
                        {geoResult.breadcrumbs.map((b: any, i: number) => (
                          <div key={b.id} className="flex items-center gap-1.5">
                            {i > 0 && <span className="text-zinc-600">›</span>}
                            <span className="px-2 py-0.5 bg-zinc-800 text-zinc-200 rounded text-[11px] font-bold">
                              {b.name} <span className="text-[9px] text-zinc-400 capitalize">({b.location_type})</span>
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {geoResult.resolved_levels && (
                    <div>
                      <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-wider mb-1.5">Autofilled Listing Fields</p>
                      <pre className="p-3 bg-black/60 rounded-xl border border-zinc-900 text-[11px] font-mono text-zinc-300 overflow-x-auto">
                        {JSON.stringify(geoResult.resolved_levels, null, 2)}
                      </pre>
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Tab 4: Batch Import Tool */}
      {activeTab === 'import' && (
        <div className="bg-zinc-900/40 border border-zinc-900 rounded-3xl p-6 space-y-6">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <h2 className="text-base font-bold text-white flex items-center gap-2">
                <Upload size={18} className="text-orange-500" />
                Authoritative Dataset Batch Importer
              </h2>
              <p className="text-xs text-zinc-400 mt-1">
                Import official Ethiopian administrative zones, woredas, and kebeles with coordinates and bounding boxes in bulk.
              </p>
            </div>
            <button
              type="button"
              onClick={loadImportTemplate}
              className="px-3 py-1.5 bg-zinc-900 hover:bg-zinc-800 text-zinc-300 text-xs font-bold rounded-xl border border-zinc-800 transition-colors"
            >
              Load Sample Template
            </button>
          </div>

          <div className="space-y-4">
            <textarea
              rows={12}
              value={importJson}
              onChange={e => setImportJson(e.target.value)}
              placeholder="Paste JSON array of location objects here..."
              className="w-full bg-zinc-950 font-mono text-xs text-zinc-200 p-4 rounded-2xl border border-zinc-800 outline-none focus:ring-1 focus:ring-orange-500"
            />

            <div className="flex items-center justify-between">
              <span className="text-[11px] text-zinc-500">
                Supports automated deduplication, slug generation, path calculation, and spatial bounds.
              </span>
              <button
                type="button"
                onClick={handleBatchImport}
                disabled={importing || !importJson.trim()}
                className="px-6 py-2.5 bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-500 hover:to-amber-500 disabled:opacity-50 text-white font-bold text-xs rounded-xl shadow-lg transition-all flex items-center gap-2"
              >
                {importing ? <RefreshCw size={14} className="animate-spin" /> : <Upload size={14} />}
                Import & Sync
              </button>
            </div>

            {importResult && (
              <div className={`p-4 rounded-2xl border text-xs ${
                importResult.error ? 'bg-red-950/40 border-red-800/50 text-red-300' : 'bg-emerald-950/40 border-emerald-800/50 text-emerald-300'
              }`}>
                {importResult.error ? (
                  <div className="flex items-center gap-2">
                    <AlertCircle size={16} />
                    <span>{importResult.error}</span>
                  </div>
                ) : (
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 font-bold text-sm">
                      <CheckCircle2 size={16} className="text-emerald-400" />
                      Import Completed Successfully!
                    </div>
                    <p className="text-[11px]">
                      Created: <strong>{importResult.created || 0}</strong> • Updated: <strong>{importResult.updated || 0}</strong> • Skipped/Failed: <strong>{importResult.failed || 0}</strong>
                    </p>
                    {importResult.errors && importResult.errors.length > 0 && (
                      <div className="mt-2 text-[10px] text-red-400 font-mono space-y-1">
                        {importResult.errors.map((e: string, i: number) => (
                          <div key={i}>⚠️ {e}</div>
                        ))}
                      </div>
                    )}
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      )}

      {/* Create / Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 z-50 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-zinc-900 border border-zinc-800 rounded-3xl w-full max-w-2xl max-h-[90vh] overflow-y-auto p-6 space-y-6">
            <div className="flex items-center justify-between pb-4 border-b border-zinc-800">
              <h2 className="text-base font-bold text-white flex items-center gap-2">
                <MapPin size={18} className="text-orange-500" />
                {editingId ? 'Edit Administrative Location' : 'Create Administrative Location'}
              </h2>
              <button
                onClick={() => setShowModal(false)}
                className="p-1.5 hover:bg-zinc-800 text-zinc-400 hover:text-white rounded-xl"
              >
                <X size={16} />
              </button>
            </div>

            {formError && (
              <div className="p-3 bg-red-950/40 border border-red-800/50 rounded-xl text-xs text-red-400 flex items-center gap-2">
                <AlertCircle size={14} />
                {formError}
              </div>
            )}

            <form onSubmit={handleSave} className="space-y-4 text-xs">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="font-bold text-zinc-400">Parent Location</label>
                  <select
                    value={form.parent_id || ''}
                    onChange={e => setForm({ ...form, parent_id: e.target.value ? Number(e.target.value) : null })}
                    className="w-full mt-1 bg-zinc-950 border border-zinc-800 rounded-xl px-3 py-2 text-zinc-200 outline-none focus:ring-1 focus:ring-orange-500"
                  >
                    <option value="">No Parent (Root Country)</option>
                    {flatList.map(l => (
                      <option key={l.id} value={l.id}>
                        {l.name_en} ({l.location_type} • L{l.level})
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="font-bold text-zinc-400">Administrative Level Type</label>
                  <select
                    value={form.location_type}
                    onChange={e => setForm({ ...form, location_type: e.target.value })}
                    className="w-full mt-1 bg-zinc-950 border border-zinc-800 rounded-xl px-3 py-2 text-zinc-200 outline-none focus:ring-1 focus:ring-orange-500"
                  >
                    {LOCATION_TYPES.map(t => (
                      <option key={t.value} value={t.value}>{t.label}</option>
                    ))}
                  </select>
                </div>
              </div>

              {/* Localized Names */}
              <div className="p-4 bg-zinc-950 rounded-2xl border border-zinc-800/80 space-y-3">
                <p className="text-[11px] font-bold text-orange-400 uppercase tracking-wider">Localized Names (Official Ethiopian Languages)</p>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <div>
                    <label className="font-bold text-zinc-400">English Name *</label>
                    <input
                      type="text"
                      required
                      value={form.name_en}
                      onChange={e => setForm({ ...form, name_en: e.target.value })}
                      placeholder="e.g. Bole"
                      className="w-full mt-1 bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-white outline-none focus:ring-1 focus:ring-orange-500"
                    />
                  </div>
                  <div>
                    <label className="font-bold text-zinc-400">Amharic (አማርኛ)</label>
                    <input
                      type="text"
                      value={form.name_am}
                      onChange={e => setForm({ ...form, name_am: e.target.value })}
                      placeholder="e.g. ቦሌ"
                      className="w-full mt-1 bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-white outline-none focus:ring-1 focus:ring-orange-500"
                    />
                  </div>
                  <div>
                    <label className="font-bold text-zinc-400">Afaan Oromoo</label>
                    <input
                      type="text"
                      value={form.name_om}
                      onChange={e => setForm({ ...form, name_om: e.target.value })}
                      placeholder="e.g. Boolee"
                      className="w-full mt-1 bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-white outline-none focus:ring-1 focus:ring-orange-500"
                    />
                  </div>
                  <div>
                    <label className="font-bold text-zinc-400">Tigrinya (ትግርኛ)</label>
                    <input
                      type="text"
                      value={form.name_ti}
                      onChange={e => setForm({ ...form, name_ti: e.target.value })}
                      placeholder="e.g. ቦሌ"
                      className="w-full mt-1 bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-white outline-none focus:ring-1 focus:ring-orange-500"
                    />
                  </div>
                </div>
              </div>

              {/* Codes & Slug */}
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <div>
                  <label className="font-bold text-zinc-400">Slug (URL friendly)</label>
                  <input
                    type="text"
                    value={form.slug}
                    onChange={e => setForm({ ...form, slug: e.target.value })}
                    placeholder="e.g. bole"
                    className="w-full mt-1 bg-zinc-950 border border-zinc-800 rounded-xl px-3 py-2 text-white outline-none focus:ring-1 focus:ring-orange-500"
                  />
                </div>
                <div>
                  <label className="font-bold text-zinc-400">Admin Code</label>
                  <input
                    type="text"
                    value={form.admin_code}
                    onChange={e => setForm({ ...form, admin_code: e.target.value })}
                    placeholder="e.g. ET-AA-BOL"
                    className="w-full mt-1 bg-zinc-950 border border-zinc-800 rounded-xl px-3 py-2 text-white outline-none focus:ring-1 focus:ring-orange-500 font-mono"
                  />
                </div>
                <div>
                  <label className="font-bold text-zinc-400">Display Order</label>
                  <input
                    type="number"
                    value={form.display_order}
                    onChange={e => setForm({ ...form, display_order: Number(e.target.value) })}
                    className="w-full mt-1 bg-zinc-950 border border-zinc-800 rounded-xl px-3 py-2 text-white outline-none focus:ring-1 focus:ring-orange-500"
                  />
                </div>
              </div>

              {/* Spatial Coordinates & Bounding Box */}
              <div className="p-4 bg-zinc-950 rounded-2xl border border-zinc-800/80 space-y-3">
                <div className="flex items-center justify-between">
                  <p className="text-[11px] font-bold text-orange-400 uppercase tracking-wider">Spatial GPS Coordinates & Bounding Box</p>
                  <span className="text-[10px] text-zinc-500">Required for 1-click GPS Autofill</span>
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="font-bold text-zinc-400">Centroid Latitude</label>
                    <input
                      type="text"
                      value={form.latitude}
                      onChange={e => setForm({ ...form, latitude: e.target.value })}
                      placeholder="8.9950"
                      className="w-full mt-1 bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-white outline-none focus:ring-1 focus:ring-orange-500 font-mono"
                    />
                  </div>
                  <div>
                    <label className="font-bold text-zinc-400">Centroid Longitude</label>
                    <input
                      type="text"
                      value={form.longitude}
                      onChange={e => setForm({ ...form, longitude: e.target.value })}
                      placeholder="38.7900"
                      className="w-full mt-1 bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-white outline-none focus:ring-1 focus:ring-orange-500 font-mono"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-4 gap-2 pt-2 border-t border-zinc-900">
                  <div>
                    <label className="text-[10px] font-bold text-zinc-500">Min Lat</label>
                    <input
                      type="text"
                      value={form.min_lat}
                      onChange={e => setForm({ ...form, min_lat: e.target.value })}
                      placeholder="8.940"
                      className="w-full mt-0.5 bg-zinc-900 border border-zinc-800 rounded-lg px-2 py-1.5 text-white font-mono text-[11px]"
                    />
                  </div>
                  <div>
                    <label className="text-[10px] font-bold text-zinc-500">Max Lat</label>
                    <input
                      type="text"
                      value={form.max_lat}
                      onChange={e => setForm({ ...form, max_lat: e.target.value })}
                      placeholder="9.040"
                      className="w-full mt-0.5 bg-zinc-900 border border-zinc-800 rounded-lg px-2 py-1.5 text-white font-mono text-[11px]"
                    />
                  </div>
                  <div>
                    <label className="text-[10px] font-bold text-zinc-500">Min Lng</label>
                    <input
                      type="text"
                      value={form.min_lng}
                      onChange={e => setForm({ ...form, min_lng: e.target.value })}
                      placeholder="38.750"
                      className="w-full mt-0.5 bg-zinc-900 border border-zinc-800 rounded-lg px-2 py-1.5 text-white font-mono text-[11px]"
                    />
                  </div>
                  <div>
                    <label className="text-[10px] font-bold text-zinc-500">Max Lng</label>
                    <input
                      type="text"
                      value={form.max_lng}
                      onChange={e => setForm({ ...form, max_lng: e.target.value })}
                      placeholder="38.860"
                      className="w-full mt-0.5 bg-zinc-900 border border-zinc-800 rounded-lg px-2 py-1.5 text-white font-mono text-[11px]"
                    />
                  </div>
                </div>
              </div>

              {/* Aliases & Search keywords */}
              <div>
                <label className="font-bold text-zinc-400">Search Aliases & Common Spellings (comma-separated)</label>
                <input
                  type="text"
                  value={form.aliases_str}
                  onChange={e => setForm({ ...form, aliases_str: e.target.value })}
                  placeholder="e.g. Bole Medhanealem, Bole International Airport, Bole Atlas"
                  className="w-full mt-1 bg-zinc-950 border border-zinc-800 rounded-xl px-3 py-2 text-white outline-none focus:ring-1 focus:ring-orange-500"
                />
              </div>

              {/* Checkboxes */}
              <div className="flex items-center gap-6 pt-2">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={form.is_active}
                    onChange={e => setForm({ ...form, is_active: e.target.checked })}
                    className="rounded text-orange-600 focus:ring-orange-500 bg-zinc-950 border-zinc-800"
                  />
                  <span className="text-zinc-300 font-medium">Is Active</span>
                </label>

                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={form.is_leaf}
                    onChange={e => setForm({ ...form, is_leaf: e.target.checked })}
                    className="rounded text-orange-600 focus:ring-orange-500 bg-zinc-950 border-zinc-800"
                  />
                  <span className="text-zinc-300 font-medium">Is Leaf Node (Has no children)</span>
                </label>
              </div>

              <div className="flex items-center justify-end gap-3 pt-4 border-t border-zinc-800">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl font-bold transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={submitting}
                  className="px-6 py-2 bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-500 hover:to-amber-500 text-white font-bold rounded-xl shadow-lg transition-all flex items-center gap-2"
                >
                  {submitting ? <RefreshCw size={14} className="animate-spin" /> : <Save size={14} />}
                  {editingId ? 'Update Location' : 'Create Location'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Delete Confirmation Modal */}
      {deleteConfirmId && (
        <div className="fixed inset-0 z-50 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-zinc-900 border border-zinc-800 rounded-3xl w-full max-w-md p-6 space-y-4">
            <h3 className="text-base font-bold text-white flex items-center gap-2">
              <AlertCircle size={18} className="text-red-500" />
              Confirm Deletion
            </h3>
            <p className="text-xs text-zinc-400">
              Are you sure you want to delete this administrative location? 
              <br /><br />
              <strong className="text-amber-400">Safety Guard:</strong> Deletion will be rejected by database constraints if this location contains child administrative units or existing product listings.
            </p>
            <div className="flex items-center justify-end gap-2 pt-2">
              <button
                onClick={() => setDeleteConfirmId(null)}
                className="px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-xs font-bold transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={() => handleDelete(deleteConfirmId)}
                className="px-4 py-2 bg-red-600 hover:bg-red-500 text-white rounded-xl text-xs font-bold shadow-lg transition-all"
              >
                Delete Location
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
