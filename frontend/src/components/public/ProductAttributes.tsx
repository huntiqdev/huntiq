'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';

interface AttributeValue {
  id: number;
  attribute_id: number;
  value_text?: string;
  value_number?: number;
  value_boolean?: boolean;
  value_date?: string;
  value_json?: Record<string, any>;
  attribute?: {
    name_en: string;
    name_am?: string;
    slug: string;
    data_type: string;
    unit?: string;
  };
}

interface ProductAttributesProps {
  productId: number;
  language?: string;
  className?: string;
}

export function ProductAttributes({ productId, language = 'en', className = '' }: ProductAttributesProps) {
  const [attributes, setAttributes] = useState<AttributeValue[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.getProductAttributes(productId)
      .then(data => setAttributes(Array.isArray(data) ? data : []))
      .catch(() => setAttributes([]))
      .finally(() => setLoading(false));
  }, [productId]);

  if (loading) {
    return (
      <div className={`space-y-2 ${className}`}>
        {[1, 2, 3].map(i => (
          <div key={i} className="h-6 bg-zinc-100 dark:bg-zinc-800 rounded animate-pulse" />
        ))}
      </div>
    );
  }

  if (attributes.length === 0) return null;

  const getDisplayValue = (attr: AttributeValue): string => {
    const def = attr.attribute;
    if (!def) return '—';
    
    switch (def.data_type) {
      case 'boolean':
        return attr.value_boolean ? 'Yes' : 'No';
      case 'number':
      case 'range':
        if (attr.value_number === undefined || attr.value_number === null) return '—';
        return `${attr.value_number.toLocaleString()}${def.unit ? ` ${def.unit}` : ''}`;
      case 'date':
        return attr.value_date ? new Date(attr.value_date).toLocaleDateString() : '—';
      case 'multiselect':
        if (attr.value_json?.values) {
          return Array.isArray(attr.value_json.values)
            ? attr.value_json.values.join(', ')
            : '—';
        }
        return '—';
      default:
        return attr.value_text || '—';
    }
  };

  const getAttrName = (attr: AttributeValue): string => {
    const def = attr.attribute;
    if (!def) return `Attribute ${attr.attribute_id}`;
    if (language === 'am' && def.name_am) return def.name_am;
    return def.name_en;
  };

  return (
    <div className={className}>
      <h3 className="text-sm font-bold text-zinc-900 dark:text-white mb-3">Specifications</h3>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
        {attributes.map(attr => (
          <div
            key={attr.id}
            className="flex items-center justify-between px-3 py-2 bg-zinc-50 dark:bg-zinc-800/50 rounded-lg"
          >
            <span className="text-xs text-zinc-500 font-medium">{getAttrName(attr)}</span>
            <span className="text-xs font-bold text-zinc-900 dark:text-white">
              {getDisplayValue(attr)}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
