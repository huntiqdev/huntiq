'use client';

import { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api';
import {
  ClassificationScheme,
  PlacementTreeNode,
  ClassificationRelationship
} from '@/lib/api/classification';
import {
  FolderTree,
  Layers,
  RefreshCw,
  Tag,
  ChevronRight,
  ChevronDown,
  ArrowRight,
  Sliders,
  CheckCircle2,
  Info
} from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Spinner } from '@/components/ui/Spinner';

export default function AdminClassificationPage() {
  const [schemes, setSchemes] = useState<ClassificationScheme[]>([]);
  const [selectedScheme, setSelectedScheme] = useState<string>('SYS_TOOLS');
  const [tree, setTree] = useState<PlacementTreeNode[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [treeLoading, setTreeLoading] = useState<boolean>(false);
  const [selectedNode, setSelectedNode] = useState<PlacementTreeNode | null>(null);
  const [effectiveAttributes, setEffectiveAttributes] = useState<any[]>([]);
  const [attrLoading, setAttrLoading] = useState<boolean>(false);
  const [relationships, setRelationships] = useState<ClassificationRelationship[]>([]);
  const [expandedNodes, setExpandedNodes] = useState<Set<number>>(new Set());

  // Load active schemes
  useEffect(() => {
    async function loadSchemes() {
      try {
        const data = await api.getSchemes('ACTIVE');
        setSchemes(data || []);
        if (data && data.length > 0 && !selectedScheme) {
          setSelectedScheme(data[0].code);
        }
      } catch (err) {
        console.error('Failed to load schemes:', err);
      } finally {
        setLoading(false);
      }
    }
    loadSchemes();
  }, [selectedScheme]);

  // Load tree when selectedScheme changes
  const loadTree = useCallback(async (code: string) => {
    setTreeLoading(true);
    try {
      const data = await api.getTree(code, 'en', 6);
      setTree(data || []);
      // auto-expand root nodes
      if (data && data.length > 0) {
        setExpandedNodes(new Set(data.map((d: PlacementTreeNode) => d.id)));
      }
    } catch (err) {
      console.error('Failed to load tree:', err);
      setTree([]);
    } finally {
      setTreeLoading(false);
    }
  }, []);

  useEffect(() => {
    if (selectedScheme) {
      loadTree(selectedScheme);
      setSelectedNode(null);
      setEffectiveAttributes([]);
    }
  }, [selectedScheme, loadTree]);

  // Load attributes and relationships when selectedNode changes
  useEffect(() => {
    if (!selectedNode) {
      setEffectiveAttributes([]);
      setRelationships([]);
      return;
    }
    setAttrLoading(true);

    // Fetch effective inherited attributes
    api.getEffectiveAttributes(selectedNode.id)
      .then(attrs => setEffectiveAttributes(attrs || []))
      .catch(() => setEffectiveAttributes([]))
      .finally(() => setAttrLoading(false));

    // Fetch relationships if concept exists
    if (selectedNode.concept_id) {
      api.getRelationships(selectedNode.concept_id)
        .then((rels: ClassificationRelationship[]) => setRelationships(rels || []))
        .catch(() => setRelationships([]));
    }
  }, [selectedNode]);

  const toggleExpand = (id: number) => {
    setExpandedNodes(prev => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const renderTreeNode = (node: PlacementTreeNode) => {
    const isExpanded = expandedNodes.has(node.id);
    const isSelected = selectedNode?.id === node.id;
    const hasChildren = node.children && node.children.length > 0;

    return (
      <div key={node.id} className="select-none">
        <div
          onClick={() => setSelectedNode(node)}
          className={`flex items-center gap-2 py-1.5 px-2 rounded-lg cursor-pointer text-sm transition-all ${
            isSelected
              ? 'bg-orange-500/10 text-orange-600 dark:text-orange-400 font-semibold'
              : 'hover:bg-zinc-100 dark:hover:bg-zinc-800 text-zinc-700 dark:text-zinc-300'
          }`}
          style={{ paddingLeft: `${node.depth * 14 + 8}px` }}
        >
          {hasChildren ? (
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation();
                toggleExpand(node.id);
              }}
              className="p-0.5 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded text-zinc-400"
            >
              {isExpanded ? <ChevronDown className="w-3.5 h-3.5" /> : <ChevronRight className="w-3.5 h-3.5" />}
            </button>
          ) : (
            <span className="w-4" />
          )}

          <span className="text-xs px-1.5 py-0.5 rounded font-mono font-bold bg-zinc-200/60 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400">
            L{node.depth}
          </span>

          <span className="truncate">{node.name || node.name_en || node.slug}</span>

          {node.role && (
            <span className="ml-auto text-[10px] uppercase font-bold tracking-wider px-1.5 py-0.5 rounded bg-zinc-100 dark:bg-zinc-800 text-zinc-500">
              {node.role}
            </span>
          )}
        </div>

        {hasChildren && isExpanded && (
          <div>
            {node.children.map(renderTreeNode)}
          </div>
        )}
      </div>
    );
  };

  if (loading) {
    return (
      <div className="flex h-96 items-center justify-center">
        <Spinner size={36} />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header bar */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-zinc-200 dark:border-zinc-800 pb-5">
        <div>
          <h1 className="text-2xl font-black text-zinc-900 dark:text-zinc-100 flex items-center gap-2">
            <FolderTree className="w-6 h-6 text-orange-500" />
            Universal Classification Platform
          </h1>
          <p className="text-sm text-zinc-500 mt-1">
            Database-driven multi-scheme classification, 7-layer hierarchy, cross-sector relationships, and inherited attributes.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <Button
            size="sm"
            variant="secondary"
            onClick={() => loadTree(selectedScheme)}
            className="flex items-center gap-2"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </Button>
        </div>
      </div>

      {/* Scheme Selector Tabs */}
      <div className="flex flex-wrap gap-2">
        {schemes.map(s => (
          <button
            key={s.code}
            type="button"
            onClick={() => setSelectedScheme(s.code)}
            className={`px-3.5 py-2 rounded-xl text-xs font-bold transition-all flex items-center gap-2 ${
              selectedScheme === s.code
                ? 'bg-orange-500 text-white shadow-md shadow-orange-500/20'
                : 'bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-700 dark:text-zinc-300'
            }`}
          >
            <Layers className="w-3.5 h-3.5" />
            {s.name_en}
            <span className="text-[10px] opacity-75 font-mono">({s.code})</span>
          </button>
        ))}
      </div>

      {/* Main Grid: Left Tree, Right Detail Panel */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Left Column: Placement Tree */}
        <div className="lg:col-span-6 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4 shadow-sm flex flex-col h-[700px]">
          <div className="flex items-center justify-between pb-3 border-b border-zinc-100 dark:border-zinc-800 mb-3">
            <div className="flex items-center gap-2">
              <FolderTree className="w-4 h-4 text-orange-500" />
              <h2 className="text-sm font-bold text-zinc-900 dark:text-zinc-100">
                Hierarchy Tree (Depth 0 - 6)
              </h2>
            </div>
            <span className="text-xs text-zinc-400 font-mono">
              {tree.length} Root Placements
            </span>
          </div>

          <div className="flex-1 overflow-y-auto pr-1">
            {treeLoading ? (
              <div className="flex h-64 items-center justify-center">
                <Spinner size={24} />
              </div>
            ) : tree.length === 0 ? (
              <div className="text-center py-16 text-zinc-400 text-sm">
                No placements found for this scheme.
              </div>
            ) : (
              <div className="space-y-0.5">
                {tree.map(renderTreeNode)}
              </div>
            )}
          </div>
        </div>

        {/* Right Column: Node Details & Inherited Attributes & Cross-Sector Links */}
        <div className="lg:col-span-6 space-y-6">
          {selectedNode ? (
            <>
              {/* Concept / Placement Meta Box */}
              <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5 shadow-sm space-y-4">
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="text-xs font-mono font-bold px-2 py-0.5 rounded bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400">
                        Layer {selectedNode.depth} • {selectedNode.role}
                      </span>
                      {selectedNode.is_leaf && (
                        <span className="text-xs font-semibold px-2 py-0.5 rounded bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400 flex items-center gap-1">
                          <CheckCircle2 className="w-3 h-3" /> Leaf Node
                        </span>
                      )}
                    </div>
                    <h3 className="text-lg font-black text-zinc-900 dark:text-zinc-100 mt-1">
                      {selectedNode.name || selectedNode.name_en}
                    </h3>
                  </div>
                  <span className="text-xs font-mono text-zinc-400">
                    ID: #{selectedNode.id}
                  </span>
                </div>

                <div className="grid grid-cols-2 gap-3 text-xs border-t border-zinc-100 dark:border-zinc-800 pt-3">
                  <div>
                    <span className="text-zinc-400 block">Slug Path</span>
                    <span className="font-mono text-zinc-700 dark:text-zinc-300 truncate block">
                      {selectedNode.path || `/${selectedNode.slug}`}
                    </span>
                  </div>
                  <div>
                    <span className="text-zinc-400 block">Concept Type</span>
                    <span className="font-semibold text-zinc-700 dark:text-zinc-300">
                      {selectedNode.concept_type || 'PRODUCT'}
                    </span>
                  </div>
                </div>
              </div>

              {/* Effective Inherited Attributes */}
              <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5 shadow-sm space-y-3">
                <div className="flex items-center justify-between border-b border-zinc-100 dark:border-zinc-800 pb-3">
                  <h4 className="text-sm font-bold text-zinc-900 dark:text-zinc-100 flex items-center gap-2">
                    <Sliders className="w-4 h-4 text-orange-500" />
                    Effective Inherited Attributes
                  </h4>
                  <span className="text-xs text-zinc-400">
                    {effectiveAttributes.length} Rules Active
                  </span>
                </div>

                {attrLoading ? (
                  <div className="py-6 flex justify-center">
                    <Spinner size={20} />
                  </div>
                ) : effectiveAttributes.length === 0 ? (
                  <p className="text-xs text-zinc-400 py-3">
                    No attribute inheritance rules defined for this branch.
                  </p>
                ) : (
                  <div className="space-y-2">
                    {effectiveAttributes.map((attr, idx) => (
                      <div
                        key={idx}
                        className="flex items-center justify-between p-2.5 rounded-xl border border-zinc-100 dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-800/30 text-xs"
                      >
                        <div>
                          <span className="font-bold text-zinc-800 dark:text-zinc-200 block">
                            {attr.name_en || attr.slug}
                          </span>
                          <span className="text-[11px] text-zinc-400 font-mono">
                            Type: {attr.data_type} {attr.unit ? `(${attr.unit})` : ''}
                          </span>
                        </div>
                        <div className="flex items-center gap-1.5">
                          {attr.is_required && (
                            <span className="text-[10px] px-1.5 py-0.5 rounded bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 font-bold">
                              Required
                            </span>
                          )}
                          {attr.is_inherited && (
                            <span className="text-[10px] px-1.5 py-0.5 rounded bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 font-bold">
                              Inherited
                            </span>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              {/* Cross-Sector Relationships */}
              <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5 shadow-sm space-y-3">
                <div className="flex items-center justify-between border-b border-zinc-100 dark:border-zinc-800 pb-3">
                  <h4 className="text-sm font-bold text-zinc-900 dark:text-zinc-100 flex items-center gap-2">
                    <Tag className="w-4 h-4 text-orange-500" />
                    Cross-Sector Graph Relationships
                  </h4>
                  <span className="text-xs text-zinc-400">
                    {relationships.length} Links
                  </span>
                </div>

                {relationships.length === 0 ? (
                  <p className="text-xs text-zinc-400 py-3">
                    No explicit cross-domain links defined for this concept.
                  </p>
                ) : (
                  <div className="space-y-2">
                    {relationships.map((rel, idx) => (
                      <div
                        key={idx}
                        className="flex items-center justify-between p-2.5 rounded-xl border border-zinc-100 dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-800/30 text-xs"
                      >
                        <div className="flex items-center gap-2">
                          <span className="px-2 py-0.5 rounded font-mono font-bold text-[10px] bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400">
                            {rel.relationship_type_code}
                          </span>
                          <ArrowRight className="w-3.5 h-3.5 text-zinc-400" />
                          <span className="font-semibold text-zinc-800 dark:text-zinc-200">
                            {rel.target_concept?.canonical_name || `Concept #${rel.target_concept_id}`}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </>
          ) : (
            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-12 text-center shadow-sm">
              <Info className="w-8 h-8 text-zinc-400 mx-auto mb-3" />
              <h3 className="text-sm font-bold text-zinc-700 dark:text-zinc-300">
                No Node Selected
              </h3>
              <p className="text-xs text-zinc-400 mt-1 max-w-sm mx-auto">
                Click on any node in the left hierarchy tree to inspect its 7-layer depth, inherited attribute rules, and cross-sector graph relationships.
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
