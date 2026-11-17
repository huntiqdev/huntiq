'use client';

import { useState, useEffect, useMemo, useRef, useCallback } from 'react';
import { api } from '@/lib/api';
import {
  ChevronDown,
  ChevronRight,
  Folder,
  FolderOpen,
  Search,
  X,
  Tag,
  Layers,
  Check,
  Compass
} from 'lucide-react';

export interface CategoryNode {
  id: number;
  name: string;
  name_en?: string;
  name_am?: string;
  name_om?: string;
  name_ti?: string;
  slug: string;
  level: number;
  is_leaf: boolean;
  count: number;
  i18n?: Record<string, string>;
  children?: CategoryNode[];
}

export interface FlattenedLeafBranch {
  id: number;
  leafName: string;
  rootName: string;
  subcatName?: string;
  fullPath: string;
  slug: string;
  count: number;
  node: CategoryNode;
}

interface CategoryTreeSelectorProps {
  value?: number;
  onChange: (categoryId: number, category: CategoryNode) => void;
  language?: string;
  schemeCode?: string;
  placeholder?: string;
  className?: string;
}

export default function CategoryTreeSelector({
  value,
  onChange,
  language = 'en',
  schemeCode = 'SYS_PROD',
  placeholder = 'Select a category branch...',
  className = '',
}: CategoryTreeSelectorProps) {
  const [categories, setCategories] = useState<CategoryNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState<Set<number>>(new Set());
  const [selected, setSelected] = useState<number | undefined>(value);
  const [prevValue, setPrevValue] = useState<number | undefined>(value);
  if (prevValue !== value) {
    setPrevValue(value);
    setSelected(value);
  }
  const [isOpen, setIsOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [viewMode, setViewMode] = useState<'leaves' | 'tree'>('leaves');
  const [serverResults, setServerResults] = useState<FlattenedLeafBranch[]>([]);
  const [searching, setSearching] = useState(false);
  const [selectedNodeCache, setSelectedNodeCache] = useState<CategoryNode | null>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const fetchCategories = useCallback(async () => {
    setLoading(true);
    try {
      let data: any[] = [];
      try {
        // Fast initial tree fetch from canonical category hierarchy
        data = await api.getCategoryTree(language, 3);
      } catch {
        data = await api.getTree(schemeCode, language, 3);
      }

      const normalize = (n: any): CategoryNode => ({
        id: n.id,
        name: n.name || n.name_en || n.canonical_name || '',
        name_en: n.name_en || n.canonical_name || '',
        name_am: n.name_am,
        name_om: n.name_om,
        name_ti: n.name_ti,
        slug: n.slug,
        level: n.depth ?? n.level ?? 0,
        is_leaf: n.is_leaf ?? (!n.children || n.children.length === 0),
        count: n.product_count ?? n.count ?? 0,
        children: Array.isArray(n.children) ? n.children.map(normalize) : undefined,
      });

      setCategories((Array.isArray(data) ? data : []).map(normalize));
    } catch (error) {
      console.error('Failed to fetch categories:', error);
      setCategories([]);
    } finally {
      setLoading(false);
    }
  }, [schemeCode, language]);

  useEffect(() => {
    fetchCategories();
  }, [fetchCategories]);

  // Close dropdown on outside click
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }
    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen]);

  // Helper to extract localized name with dynamic i18n support
  const getNodeName = useCallback((node: CategoryNode | { name?: string; name_en?: string; name_am?: string; name_om?: string; name_ti?: string; slug?: string; i18n?: Record<string, string> }): string => {
    if (node.i18n && node.i18n[language]) {
      return node.i18n[language];
    }
    switch (language) {
      case 'am':
        return node.name_am || node.name || node.name_en || node.slug || '';
      case 'om':
      case 'or':
        return node.name_om || node.name || node.name_en || node.slug || '';
      case 'ti':
        return node.name_ti || node.name || node.name_en || node.slug || '';
      default:
        return node.i18n?.['en'] || node.name_en || node.name || node.slug || '';
    }
  }, [language]);

  // Debounced server-side universal classification search across all 100k+ leaves & 4 languages
  useEffect(() => {
    const q = searchQuery.trim();
    if (q.length < 2) {
      setServerResults([]);
      setSearching(false);
      return;
    }

    setSearching(true);
    const timer = setTimeout(async () => {
      try {
        const results = await api.searchCategories(q, language, 40);
        if (Array.isArray(results)) {
          const mapped: FlattenedLeafBranch[] = results.map(r => {
            const leafName = r.name || r.name_en;
            const bcs = Array.isArray(r.breadcrumbs) ? r.breadcrumbs : [];
            const rootName = bcs[0]?.name || 'General';
            const subcatName = bcs[1]?.name || '';
            const fullPath = r.hierarchy || bcs.map((b: any) => b.name).join(' › ') || leafName;
            return {
              id: r.id,
              leafName,
              rootName,
              subcatName,
              fullPath,
              slug: r.slug,
              count: 0,
              node: {
                id: r.id,
                name: leafName,
                name_en: r.name_en,
                name_am: r.name_am,
                name_om: r.name_om,
                name_ti: r.name_ti,
                slug: r.slug,
                level: r.level,
                is_leaf: r.is_leaf,
                count: 0,
              },
            };
          });
          setServerResults(mapped);
        }
      } catch (err) {
        console.error('Classification search error:', err);
      } finally {
        setSearching(false);
      }
    }, 200);

    return () => clearTimeout(timer);
  }, [searchQuery, schemeCode, language, getNodeName]);

  // Extract all flattened leaf branches with their complete hierarchy path from initial loaded tree
  const allLeafBranches = useMemo<FlattenedLeafBranch[]>(() => {
    const leaves: FlattenedLeafBranch[] = [];

    const traverse = (node: CategoryNode, breadcrumbs: string[], rootName: string, subcatName?: string) => {
      const currentName = getNodeName(node);
      const currentPath = [...breadcrumbs, currentName];

      const hasChildren = Array.isArray(node.children) && node.children.length > 0;
      if (node.is_leaf) {
        leaves.push({
          id: node.id,
          leafName: currentName,
          rootName,
          subcatName,
          fullPath: currentPath.join(' › '),
          slug: node.slug,
          count: node.count || 0,
          node,
        });
      } else if (hasChildren) {
        const nextSubcat = breadcrumbs.length === 1 ? currentName : subcatName;
        node.children!.forEach(child => traverse(child, currentPath, rootName, nextSubcat));
      }
    };

    categories.forEach(root => {
      const rootTitle = getNodeName(root);
      const hasChildren = Array.isArray(root.children) && root.children.length > 0;
      if (root.is_leaf) {
        leaves.push({
          id: root.id,
          leafName: rootTitle,
          rootName: rootTitle,
          fullPath: rootTitle,
          slug: root.slug,
          count: root.count || 0,
          node: root,
        });
      } else if (hasChildren) {
        root.children!.forEach(child => traverse(child, [rootTitle], rootTitle));
      }
    });

    return leaves;
  }, [categories, getNodeName]);

  // Filter leaf branches by live search query across all 4 languages & slugs
  const filteredLeafBranches = useMemo(() => {
    if (!searchQuery.trim()) return allLeafBranches;
    const q = searchQuery.toLowerCase().trim();

    return allLeafBranches.filter(b => {
      const node = b.node;
      const nEn = (node.name_en || '').toLowerCase();
      const nAm = (node.name_am || '').toLowerCase();
      const nOm = (node.name_om || '').toLowerCase();
      const nTi = (node.name_ti || '').toLowerCase();
      const slug = (node.slug || '').toLowerCase();
      const fullPath = b.fullPath.toLowerCase();

      return (
        fullPath.includes(q) ||
        nEn.includes(q) ||
        nAm.includes(q) ||
        nOm.includes(q) ||
        nTi.includes(q) ||
        slug.includes(q)
      );
    });
  }, [allLeafBranches, searchQuery]);

  // Group leaf branches by Root category (prefers server search results when searching)
  const displayedBranches = useMemo(() => {
    if (searchQuery.trim().length >= 2 && serverResults.length > 0) {
      return serverResults;
    }
    return filteredLeafBranches.slice(0, 100);
  }, [searchQuery, serverResults, filteredLeafBranches]);

  const groupedBranches = useMemo(() => {
    const map = new Map<string, FlattenedLeafBranch[]>();
    displayedBranches.forEach(branch => {
      const list = map.get(branch.rootName) || [];
      list.push(branch);
      map.set(branch.rootName, list);
    });
    return Array.from(map.entries());
  }, [displayedBranches]);

  const toggleExpand = async (id: number, node?: CategoryNode) => {
    const newExpanded = new Set(expanded);
    if (newExpanded.has(id)) {
      newExpanded.delete(id);
      setExpanded(newExpanded);
      return;
    }

    newExpanded.add(id);
    setExpanded(newExpanded);

    // On-demand lazy load children for deeper levels (layers 4, 5, 6)
    if (node && !node.is_leaf && (!node.children || node.children.length === 0)) {
      try {
        let rawChildren: any[] = [];
        try {
          rawChildren = await api.getCategoryChildren(id, language);
        } catch {
          rawChildren = await api.getChildren(id, language);
        }
        const children = rawChildren;
        if (Array.isArray(children) && children.length > 0) {
          const normalize = (n: any): CategoryNode => ({
            id: n.id,
            name: n.name || n.name_en || n.canonical_name || '',
            name_en: n.name_en || n.canonical_name || '',
            name_am: n.name_am,
            name_om: n.name_om,
            name_ti: n.name_ti,
            slug: n.slug,
            level: n.depth ?? n.level ?? 0,
            is_leaf: n.is_leaf ?? (!n.children || n.children.length === 0),
            count: n.product_count ?? n.count ?? 0,
            children: Array.isArray(n.children) ? n.children.map(normalize) : undefined,
          });
          const mappedChildren = children.map(normalize);

          const attachChildren = (nodes: CategoryNode[]): CategoryNode[] => {
            return nodes.map(item => {
              if (item.id === id) {
                return { ...item, children: mappedChildren };
              }
              if (item.children && item.children.length > 0) {
                return { ...item, children: attachChildren(item.children) };
              }
              return item;
            });
          };
          setCategories(prev => attachChildren(prev));
        }
      } catch (err) {
        console.error('Failed to load children for node', id, err);
      }
    }
  };

  const handleSelect = (category: CategoryNode) => {
    setSelected(category.id);
    setSelectedNodeCache(category);
    onChange(category.id, category);
    setIsOpen(false);
  };

  const handleSelectLeaf = (branch: FlattenedLeafBranch) => {
    setSelected(branch.id);
    setSelectedNodeCache(branch.node);
    onChange(branch.id, branch.node);
    setIsOpen(false);
  };

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation();
    setSelected(undefined);
    setSelectedNodeCache(null);
  };

  // Find full path of currently selected category
  const getSelectedDisplay = (): string => {
    if (!selected) return placeholder;
    if (selectedNodeCache && selectedNodeCache.id === selected) {
      return getNodeName(selectedNodeCache);
    }
    const foundBranch = allLeafBranches.find(b => b.id === selected) || serverResults.find(b => b.id === selected);
    if (foundBranch) {
      return foundBranch.fullPath;
    }

    const findName = (nodes: CategoryNode[]): string | null => {
      for (const node of nodes) {
        if (node.id === selected) return getNodeName(node);
        if (node.children) {
          const found = findName(node.children);
          if (found) return found;
        }
      }
      return null;
    };

    return findName(categories) || placeholder;
  };

  // Localized texts
  const tLabels = {
    allBranches: language === 'am' ? 'ሁሉም መዳረሻ ቅርንጫፎች' : language === 'or' ? 'Dameewwan Hundaa' : language === 'ti' ? 'ኩሎም ቅርንጫፋት' : 'All Leaf Destinations',
    hierarchyTree: language === 'am' ? 'የደረጃ ዛፍ' : language === 'or' ? 'Muka Qindoominaa' : language === 'ti' ? 'ናይ ደረጃ ዛፍ' : 'Hierarchy Tree',
    searchPlaceholder: language === 'am' ? 'ቅርንጫፍ ወይም መዳረሻ ይፈልጉ...' : language === 'or' ? 'Damee ykn galma barbaadi...' : language === 'ti' ? 'ቅርንጫፍ ወይ መዕረፊ ድለዩ...' : 'Search leaf destinations, branches, or keywords...',
    noResults: language === 'am' ? 'ምንም አይነት ተዛማጅ ቅርንጫፍ አልተገኘም' : language === 'or' ? 'Dameen walsimu hin argamne' : language === 'ti' ? 'ዝተዛመደ ቅርንጫፍ ኣይተረኸበን' : 'No matching branch destinations found',
    leafBadge: language === 'am' ? 'መዳረሻ' : language === 'or' ? 'Galma' : language === 'ti' ? 'መዕረፊ' : 'Leaf Destination',
    selectPrompt: language === 'am' ? 'እባክዎ የተፈለገውን መዳረሻ ይምረጡ' : language === 'or' ? 'Mee galma barbaaddan filadhaa' : language === 'ti' ? 'በጃኹም ዝደለይዎ መዕረፊ ምረጹ' : 'Click any leaf destination below to select',
  };

  const renderTreeNode = (node: CategoryNode, depth: number = 0) => {
    const isExpanded = expanded.has(node.id);
    const isSelected = selected === node.id;
    const hasChildren = Array.isArray(node.children) && node.children.length > 0;
    const displayName = getNodeName(node);

    return (
      <div key={node.id} style={{ marginLeft: `${depth * 16}px` }}>
        <div
          className={`
            flex items-center gap-2 px-3 py-2 cursor-pointer rounded-xl transition-all my-0.5
            ${isSelected
              ? 'bg-[#1769E0]/10 text-[#1769E0] dark:text-[#38A5FF] font-bold border border-[#1769E0]/20'
              : 'hover:bg-zinc-100 dark:hover:bg-zinc-800 text-zinc-800 dark:text-zinc-200'}
            ${!node.is_leaf && hasChildren ? 'font-semibold' : ''}
          `}
          onClick={() => {
            if (node.is_leaf) {
              handleSelect(node);
            } else {
              toggleExpand(node.id, node);
            }
          }}
        >
          {hasChildren || !node.is_leaf ? (
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation();
                toggleExpand(node.id, node);
              }}
              className="p-1 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded-lg text-zinc-500"
            >
              {isExpanded ? (
                <ChevronDown className="w-4 h-4" />
              ) : (
                <ChevronRight className="w-4 h-4" />
              )}
            </button>
          ) : (
            <div className="w-6 flex items-center justify-center">
              <Tag className="w-3.5 h-3.5 text-[#1769E0]" />
            </div>
          )}

          {hasChildren || !node.is_leaf ? (
            isExpanded ? (
              <FolderOpen className="w-4 h-4 text-amber-500 shrink-0" />
            ) : (
              <Folder className="w-4 h-4 text-zinc-400 dark:text-zinc-500 shrink-0" />
            )
          ) : null}

          <span className="flex-1 text-sm truncate">{displayName}</span>

          {node.is_leaf && (
            <span className="text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-full bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300">
              {tLabels.leafBadge}
            </span>
          )}

          {node.count > 0 && (
            <span className="text-xs text-zinc-500 bg-zinc-200/60 dark:bg-zinc-700/60 px-2 py-0.5 rounded-full font-mono">
              {node.count}
            </span>
          )}

          {isSelected && (
            <Check className="w-4 h-4 text-[#1769E0] dark:text-[#38A5FF] ml-1 shrink-0" />
          )}
        </div>

        {isExpanded && hasChildren && (
          <div className="mt-0.5 border-l-2 border-zinc-200 dark:border-zinc-800 ml-3">
            {node.children!.map(child => renderTreeNode(child, depth + 1))}
          </div>
        )}
      </div>
    );
  };

  if (loading) {
    return (
      <div className={`border border-zinc-200 dark:border-zinc-700 rounded-2xl p-3 bg-zinc-50 dark:bg-zinc-800/50 ${className}`}>
        <div className="animate-pulse flex items-center gap-2">
          <div className="w-4 h-4 bg-zinc-300 dark:bg-zinc-700 rounded-full" />
          <div className="h-4 bg-zinc-200 dark:bg-zinc-700 rounded w-1/2" />
        </div>
      </div>
    );
  }

  const isItemSelected = selected !== undefined && selected !== null;

  return (
    <div className={`relative ${className}`} ref={dropdownRef}>
      {/* Trigger Button */}
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className={`w-full border rounded-2xl px-3.5 py-2.5 min-h-[44px] text-left flex items-center justify-between transition-all outline-none focus:ring-2 focus:ring-[#1769E0] tap-bounce active:scale-[0.99] ${
          isItemSelected
            ? 'border-[#1769E0]/40 bg-[#1769E0]/5 dark:bg-[#1769E0]/20 text-zinc-900 dark:text-zinc-100'
            : 'border-zinc-200 dark:border-zinc-700 bg-zinc-50 dark:bg-zinc-800 text-zinc-500 dark:text-zinc-400'
        }`}
      >
        <div className="flex items-center gap-2.5 overflow-hidden flex-1 mr-2">
          <Compass className={`w-4 h-4 shrink-0 ${isItemSelected ? 'text-[#1769E0] dark:text-[#38A5FF]' : 'text-zinc-400'}`} />
          <span className={`text-sm truncate ${isItemSelected ? 'font-semibold text-zinc-900 dark:text-zinc-100' : ''}`}>
            {getSelectedDisplay()}
          </span>
        </div>

        <div className="flex items-center gap-1.5 shrink-0">
          {isItemSelected && (
            <span
              onClick={handleClear}
              className="p-1 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded-full text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200 transition-colors"
              title="Clear selection"
            >
              <X className="w-3.5 h-3.5" />
            </span>
          )}
          <ChevronDown
            className={`w-4 h-4 text-zinc-400 transition-transform duration-200 ${
              isOpen ? 'transform rotate-180 text-[#1769E0]' : ''
            }`}
          />
        </div>
      </button>

      {/* Dropdown Menu */}
      {isOpen && (
        <div className="absolute z-50 w-full mt-2 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-700 rounded-2xl shadow-2xl overflow-hidden animate-in fade-in-50 zoom-in-95 duration-150">
          {/* Header with Search and Mode Toggle */}
          <div className="p-3 border-b border-zinc-100 dark:border-zinc-800 space-y-2.5 bg-zinc-50/70 dark:bg-zinc-800/40">
            {/* Live Search Input */}
            <div className="relative">
              <Search className="w-4 h-4 text-zinc-400 absolute left-3 top-1/2 -translate-y-1/2" />
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder={tLabels.searchPlaceholder}
                className="w-full pl-9 pr-8 py-2 text-base sm:text-xs min-h-[44px] border border-zinc-200 dark:border-zinc-700 rounded-xl bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-100 outline-none focus:ring-2 focus:ring-[#1769E0] focus:border-[#1769E0] transition-all"
                autoFocus
              />
              {searchQuery && (
                <div className="absolute right-2.5 top-1/2 -translate-y-1/2 flex items-center gap-1.5">
                  {searching && (
                    <div className="w-3.5 h-3.5 border-2 border-[#1769E0] border-t-transparent rounded-full animate-spin" />
                  )}
                  <button
                    type="button"
                    onClick={() => setSearchQuery('')}
                    className="text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200"
                  >
                    <X className="w-3.5 h-3.5" />
                  </button>
                </div>
              )}
            </div>

            {/* View Mode Toggle: All Branches vs Hierarchy Tree */}
            <div className="flex items-center justify-between text-xs">
              <div className="flex bg-zinc-200/70 dark:bg-zinc-800 p-0.5 rounded-xl border border-zinc-200/60 dark:border-zinc-700">
                <button
                  type="button"
                  onClick={() => setViewMode('leaves')}
                  className={`px-3 py-1.5 rounded-lg font-bold flex items-center gap-1.5 transition-all ${
                    viewMode === 'leaves'
                      ? 'bg-white dark:bg-zinc-700 text-[#1769E0] dark:text-[#38A5FF] shadow-sm'
                      : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-200'
                  }`}
                >
                  <Layers className="w-3.5 h-3.5" />
                  <span>{tLabels.allBranches}</span>
                  <span className="text-[10px] bg-zinc-100 dark:bg-zinc-600 px-1.5 rounded-full font-mono">
                    {displayedBranches.length}
                  </span>
                </button>

                <button
                  type="button"
                  onClick={() => setViewMode('tree')}
                  className={`px-3 py-1.5 rounded-lg font-bold flex items-center gap-1.5 transition-all ${
                    viewMode === 'tree'
                      ? 'bg-white dark:bg-zinc-700 text-[#1769E0] dark:text-[#38A5FF] shadow-sm'
                      : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-200'
                  }`}
                >
                  <Folder className="w-3.5 h-3.5" />
                  <span>{tLabels.hierarchyTree}</span>
                </button>
              </div>

              <span className="text-[11px] text-zinc-400 hidden sm:inline">
                {tLabels.selectPrompt}
              </span>
            </div>
          </div>

          {/* List Content */}
          <div className="max-h-[380px] overflow-y-auto p-2 scrollbar-thin scrollbar-thumb-zinc-300 dark:scrollbar-thumb-zinc-700">
            {viewMode === 'leaves' ? (
              // 🌿 Leaf Branches Listed Out One by One with Destination Breadcrumbs
              groupedBranches.length === 0 ? (
                <div className="text-center text-zinc-400 py-8 text-xs">
                  {tLabels.noResults}
                </div>
              ) : (
                <div className="space-y-4 p-1">
                  {groupedBranches.map(([rootTitle, branches]) => (
                    <div key={rootTitle} className="space-y-1">
                      {/* Root Category Domain Header */}
                      <div className="sticky top-0 bg-white/95 dark:bg-zinc-900/95 backdrop-blur-sm z-10 py-1 px-2 border-b border-zinc-100 dark:border-zinc-800 flex items-center justify-between">
                        <span className="text-[11px] font-black uppercase tracking-wider text-[#1769E0] dark:text-[#38A5FF] flex items-center gap-1.5">
                          <Folder className="w-3.5 h-3.5" />
                          {rootTitle}
                        </span>
                        <span className="text-[10px] text-zinc-400 font-mono">
                          {branches.length} {branches.length === 1 ? 'branch' : 'branches'}
                        </span>
                      </div>

                      {/* Individual Leaf Destinations */}
                      <div className="grid grid-cols-1 gap-1 pt-1">
                        {branches.map(branch => {
                          const isSelected = selected === branch.id;
                          return (
                            <div
                              key={branch.id}
                              onClick={() => handleSelectLeaf(branch)}
                              className={`
                                flex items-center justify-between p-2.5 rounded-xl cursor-pointer transition-all border
                                ${isSelected
                                  ? 'bg-[#1769E0]/10 border-[#1769E0]/30 text-[#1769E0] dark:text-[#38A5FF]'
                                  : 'border-transparent hover:border-zinc-200 dark:hover:border-zinc-700 hover:bg-zinc-50 dark:hover:bg-zinc-800/60 text-zinc-800 dark:text-zinc-200'}
                              `}
                            >
                              <div className="flex-1 min-w-0 pr-2">
                                {/* Subcategory Breadcrumb if present */}
                                {branch.subcatName && (
                                  <div className="text-[10px] text-zinc-400 truncate mb-0.5">
                                    {branch.subcatName}
                                  </div>
                                )}
                                {/* Leaf Destination Name */}
                                <div className="text-xs font-semibold truncate flex items-center gap-1.5">
                                  <Tag className="w-3 h-3 text-[#1769E0] shrink-0" />
                                  <span>{branch.leafName}</span>
                                </div>
                              </div>

                              <div className="flex items-center gap-2 shrink-0">
                                <span className="text-[10px] font-bold px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800/40">
                                  {tLabels.leafBadge}
                                </span>

                                {branch.count > 0 && (
                                  <span className="text-[11px] text-zinc-400 bg-zinc-100 dark:bg-zinc-800 px-2 py-0.5 rounded-full font-mono">
                                    {branch.count}
                                  </span>
                                )}

                                {isSelected && (
                                  <Check className="w-4 h-4 text-[#1769E0] dark:text-[#38A5FF]" />
                                )}
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    </div>
                  ))}
                </div>
              )
            ) : (
              // 📁 Hierarchy Tree View
              categories.length === 0 ? (
                <div className="text-center text-zinc-400 py-8 text-xs">
                  {tLabels.noResults}
                </div>
              ) : (
                <div className="space-y-0.5 p-1">
                  {categories.map(node => renderTreeNode(node))}
                </div>
              )
            )}
          </div>
          {/* Quick "Other / Can't Find" Bar */}
          <div className="p-2.5 border-t border-zinc-200 dark:border-zinc-800 bg-zinc-50/90 dark:bg-zinc-850/80 flex items-center justify-between text-xs">
            <span className="text-[11px] text-zinc-500 font-medium">
              {language === 'am' ? 'የሚፈልጉትን ምድብ አላገኙም?' : language === 'or' ? 'Ramaddii hin argannee?' : language === 'ti' ? 'ዝደለይዎ ምድብ ኣይረኸብኩምን?' : "Can't find your category?"}
            </span>
            <button
              type="button"
              onClick={() => {
                const otherBranch = allLeafBranches.find(b => b.slug.startsWith('other-')) || allLeafBranches[allLeafBranches.length - 1];
                if (otherBranch) {
                  handleSelectLeaf(otherBranch);
                }
              }}
              className="px-3 py-1.5 rounded-lg bg-[#1769E0] text-white font-bold hover:bg-[#1055BD] transition-all text-xs shadow-sm shadow-[#1769E0]/20 active:scale-95"
            >
              {language === 'am' ? 'ሌላ ምድብ ምረጥ (Other)' : language === 'or' ? 'Biroo Filadhu' : language === 'ti' ? 'ካልእ ምረጽ' : 'Select "Other"'}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
