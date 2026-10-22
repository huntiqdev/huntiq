'use client';

import { useState, useEffect } from 'react';
import { ChevronDown, ChevronRight, Folder, FolderOpen } from 'lucide-react';

interface CategoryNode {
  id: number;
  name: string;
  slug: string;
  level: number;
  is_leaf: boolean;
  count: number;
  children: CategoryNode[];
}

interface CategoryTreeSelectorProps {
  value?: number;
  onChange: (categoryId: number, category: CategoryNode) => void;
  language?: string;
  placeholder?: string;
}

export default function CategoryTreeSelector({
  value,
  onChange,
  language = 'en',
  placeholder = 'Select a category'
}: CategoryTreeSelectorProps) {
  const [categories, setCategories] = useState<CategoryNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState<Set<number>>(new Set());
  const [selected, setSelected] = useState<number | undefined>(value);
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    fetchCategories();
  }, [language]);

  useEffect(() => {
    setSelected(value);
  }, [value]);

  const fetchCategories = async () => {
    try {
      const res = await fetch(`/api/categories/tree?lang=${language}`);
      const data = await res.json();
      if (data.success) {
        setCategories(data.data || []);
      }
    } catch (error) {
      console.error('Failed to fetch categories:', error);
    } finally {
      setLoading(false);
    }
  };

  const toggleExpand = (id: number) => {
    const newExpanded = new Set(expanded);
    if (newExpanded.has(id)) {
      newExpanded.delete(id);
    } else {
      newExpanded.add(id);
    }
    setExpanded(newExpanded);
  };

  const handleSelect = (category: CategoryNode) => {
    if (category.is_leaf) {
      setSelected(category.id);
      onChange(category.id, category);
      setIsOpen(false);
    } else {
      toggleExpand(category.id);
    }
  };

  const getSelectedName = (): string => {
    if (!selected) return placeholder;
    
    const findName = (nodes: CategoryNode[]): string | null => {
      for (const node of nodes) {
        if (node.id === selected) return node.name;
        if (node.children) {
          const found = findName(node.children);
          if (found) return found;
        }
      }
      return null;
    };
    
    return findName(categories) || placeholder;
  };

  const renderNode = (node: CategoryNode, depth: number = 0) => {
    const isExpanded = expanded.has(node.id);
    const isSelected = selected === node.id;
    const hasChildren = node.children && node.children.length > 0;

    return (
      <div key={node.id} style={{ marginLeft: `${depth * 20}px` }}>
        <div
          className={`
            flex items-center gap-2 px-3 py-2 cursor-pointer rounded
            hover:bg-gray-100
            ${isSelected ? 'bg-blue-50 text-blue-700' : ''}
            ${!node.is_leaf ? 'font-medium' : ''}
          `}
          onClick={() => handleSelect(node)}
        >
          {hasChildren && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleExpand(node.id);
              }}
              className="p-0.5 hover:bg-gray-200 rounded"
            >
              {isExpanded ? (
                <ChevronDown className="w-4 h-4" />
              ) : (
                <ChevronRight className="w-4 h-4" />
              )}
            </button>
          )}
          
          {!hasChildren && <div className="w-5" />}
          
          {isExpanded ? (
            <FolderOpen className="w-4 h-4 text-yellow-600" />
          ) : (
            <Folder className="w-4 h-4 text-gray-400" />
          )}
          
          <span className="flex-1">{node.name}</span>
          
          {node.count > 0 && (
            <span className="text-xs text-gray-500 bg-gray-100 px-2 py-0.5 rounded">
              {node.count}
            </span>
          )}
          
          {!node.is_leaf && (
            <span className="text-xs text-gray-400">(folder)</span>
          )}
        </div>
        
        {isExpanded && hasChildren && (
          <div className="mt-1">
            {node.children.map(child => renderNode(child, depth + 1))}
          </div>
        )}
      </div>
    );
  };

  if (loading) {
    return (
      <div className="border border-gray-300 rounded-lg p-3">
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-1/3"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="relative">
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className="w-full border border-gray-300 rounded-lg px-4 py-2 text-left flex items-center justify-between hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <span className={selected ? 'text-gray-900' : 'text-gray-500'}>
          {getSelectedName()}
        </span>
        <ChevronDown
          className={`w-5 h-5 text-gray-400 transition-transform ${
            isOpen ? 'transform rotate-180' : ''
          }`}
        />
      </button>

      {isOpen && (
        <>
          <div
            className="fixed inset-0 z-10"
            onClick={() => setIsOpen(false)}
          />
          <div className="absolute z-20 w-full mt-2 bg-white border border-gray-300 rounded-lg shadow-lg max-h-96 overflow-auto">
            <div className="p-2">
              {categories.length === 0 ? (
                <div className="text-center text-gray-500 py-4">
                  No categories available
                </div>
              ) : (
                categories.map(node => renderNode(node))
              )}
            </div>
          </div>
        </>
      )}
    </div>
  );
}
