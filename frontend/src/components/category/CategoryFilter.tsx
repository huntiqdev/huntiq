'use client';

import { useState, useEffect } from 'react';
import { ChevronDown, ChevronRight, Filter } from 'lucide-react';

interface CategoryNode {
  id: number;
  name: string;
  slug: string;
  level: number;
  is_leaf: boolean;
  count: number;
  children: CategoryNode[];
}

interface CategoryFilterProps {
  value?: number;
  includeDescendants?: boolean;
  onChange: (categoryId: number | undefined, includeDescendants: boolean) => void;
  language?: string;
}

export default function CategoryFilter({
  value,
  includeDescendants = true,
  onChange,
  language = 'en'
}: CategoryFilterProps) {
  const [categories, setCategories] = useState<CategoryNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState<Set<number>>(new Set());
  const [selected, setSelected] = useState<number | undefined>(value);
  const [includeChildren, setIncludeChildren] = useState(includeDescendants);

  useEffect(() => {
    fetchCategories();
  }, [language]);

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

  const handleCategoryClick = (categoryId: number) => {
    const newSelected = selected === categoryId ? undefined : categoryId;
    setSelected(newSelected);
    onChange(newSelected, includeChildren);
  };

  const handleIncludeChildrenToggle = () => {
    const newInclude = !includeChildren;
    setIncludeChildren(newInclude);
    if (selected) {
      onChange(selected, newInclude);
    }
  };

  const renderNode = (node: CategoryNode, depth: number = 0) => {
    const isExpanded = expanded.has(node.id);
    const isSelected = selected === node.id;
    const hasChildren = node.children && node.children.length > 0;

    return (
      <div key={node.id}>
        <div
          className={`
            flex items-center gap-2 px-3 py-2 rounded cursor-pointer
            hover:bg-gray-50
            ${isSelected ? 'bg-blue-50 border-l-4 border-blue-500' : ''}
          `}
          style={{ paddingLeft: `${depth * 16 + 12}px` }}
        >
          {hasChildren ? (
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
          ) : (
            <div className="w-5" />
          )}
          
          <div
            className="flex-1 flex items-center justify-between"
            onClick={() => handleCategoryClick(node.id)}
          >
            <span className={`text-sm ${isSelected ? 'font-medium text-blue-700' : 'text-gray-700'}`}>
              {node.name}
            </span>
            
            {node.count > 0 && (
              <span className="text-xs text-gray-500 bg-gray-100 px-2 py-0.5 rounded ml-2">
                {node.count}
              </span>
            )}
          </div>
        </div>
        
        {isExpanded && hasChildren && (
          <div>
            {node.children.map(child => renderNode(child, depth + 1))}
          </div>
        )}
      </div>
    );
  };

  if (loading) {
    return (
      <div className="bg-white rounded-lg border border-gray-200 p-4">
        <div className="flex items-center gap-2 mb-4">
          <Filter className="w-5 h-5 text-gray-400" />
          <h3 className="font-medium text-gray-900">Categories</h3>
        </div>
        <div className="space-y-2">
          {[1, 2, 3].map(i => (
            <div key={i} className="h-8 bg-gray-100 rounded animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-4">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <Filter className="w-5 h-5 text-gray-600" />
          <h3 className="font-medium text-gray-900">Categories</h3>
        </div>
        
        {selected && (
          <button
            onClick={() => {
              setSelected(undefined);
              onChange(undefined, includeChildren);
            }}
            className="text-sm text-blue-600 hover:text-blue-700"
          >
            Clear
          </button>
        )}
      </div>

      {selected && (
        <div className="mb-3 pb-3 border-b border-gray-200">
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={includeChildren}
              onChange={handleIncludeChildrenToggle}
              className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span className="text-sm text-gray-700">
              Include subcategories
            </span>
          </label>
        </div>
      )}

      <div className="space-y-1 max-h-96 overflow-y-auto">
        {categories.length === 0 ? (
          <div className="text-center text-gray-500 py-4 text-sm">
            No categories available
          </div>
        ) : (
          categories.map(node => renderNode(node))
        )}
      </div>
    </div>
  );
}
