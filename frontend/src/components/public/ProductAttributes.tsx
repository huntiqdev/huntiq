'use client';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import { api } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import { useLocale } from '@/hooks/useLocale';
import { Lock, Shield, ArrowRight, Sparkles, User } from 'lucide-react';

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
  requiredBuyerTier?: string;
  specificationsLocked?: boolean;
}

export function ProductAttributes({
  productId,
  language = 'en',
  className = '',
  requiredBuyerTier = 'normal',
  specificationsLocked = false,
}: ProductAttributesProps) {
  const { user, buyerTier } = useAuth();
  const { t } = useLocale();
  const [attributes, setAttributes] = useState<AttributeValue[]>([]);
  const [loading, setLoading] = useState(true);
  const [isLocked, setIsLocked] = useState(specificationsLocked);

  const effectiveRequiredTier = requiredBuyerTier || 'normal';
  const isTierRestricted = effectiveRequiredTier === 'gold' || effectiveRequiredTier === 'premium';

  useEffect(() => {
    // If explicitly marked as locked from server, keep locked
    if (specificationsLocked) {
      setIsLocked(true);
      setLoading(false);
      return;
    }

    // Call API with token if authenticated
    api.getProductAttributes(productId)
      .then((data: any) => {
        if (data && typeof data === 'object' && data.locked) {
          setIsLocked(true);
          setAttributes([]);
        } else if (Array.isArray(data)) {
          setAttributes(data);
          setIsLocked(false);
        } else {
          setAttributes([]);
        }
      })
      .catch(() => setAttributes([]))
      .finally(() => setLoading(false));
  }, [productId, specificationsLocked, user]);

  if (loading) {
    return (
      <div className={`space-y-2 ${className}`}>
        {[1, 2, 3].map(i => (
          <div key={i} className="h-6 bg-zinc-100 dark:bg-zinc-800 rounded-xl animate-pulse" />
        ))}
      </div>
    );
  }

  // If specifications are locked by server-side tier policy
  if (isLocked || (isTierRestricted && specificationsLocked)) {
    const isGold = effectiveRequiredTier === 'gold';
    const tierName = isGold ? 'Gold' : 'Premium';
    const tierBadgeStyle = isGold
      ? 'bg-gradient-to-r from-amber-500 to-yellow-500 text-white shadow-amber-500/20'
      : 'bg-gradient-to-r from-purple-600 to-indigo-600 text-white shadow-purple-600/20';

    const cardBgStyle = isGold
      ? 'border-amber-200/80 dark:border-amber-900/50 bg-gradient-to-br from-amber-50/70 via-white to-yellow-50/40 dark:from-amber-950/20 dark:via-zinc-900 dark:to-yellow-950/10'
      : 'border-purple-200/80 dark:border-purple-900/50 bg-gradient-to-br from-purple-50/70 via-white to-indigo-50/40 dark:from-purple-950/20 dark:via-zinc-900 dark:to-indigo-950/10';

    return (
      <div className={`rounded-2xl border p-4 sm:p-5 shadow-sm overflow-hidden relative ${cardBgStyle} ${className}`}>
        <div className="flex items-start justify-between gap-3 mb-3">
          <div className="flex items-center gap-2">
            <div className={`w-8 h-8 rounded-xl flex items-center justify-center shadow-md ${tierBadgeStyle}`}>
              <Lock size={15} />
            </div>
            <div>
              <p className="text-xs font-black uppercase tracking-wider text-zinc-900 dark:text-white flex items-center gap-1.5">
                {isGold
                  ? (t.ui as any).gold_specs_only ?? 'Gold Members Only'
                  : (t.ui as any).premium_specs_only ?? 'Premium Members Only'}
              </p>
              <p className="text-[11px] text-zinc-500 dark:text-zinc-400">
                {(t.ui as any).specifications_locked ?? 'Restricted Product Specifications'}
              </p>
            </div>
          </div>
          <span className={`text-[10px] font-black uppercase tracking-widest px-2.5 py-0.5 rounded-full shadow-sm ${tierBadgeStyle}`}>
            {tierName}
          </span>
        </div>

        <p className="text-xs text-zinc-600 dark:text-zinc-300 leading-relaxed mb-4">
          {isGold
            ? (t.ui as any).specs_locked_desc_gold ?? 'Full verified technical specifications, dimensions, and model details for this listing are reserved for Gold members.'
            : (t.ui as any).specs_locked_desc_premium ?? 'Full verified technical specifications, commercial documentation, and model details for this listing are reserved for Premium members.'}
        </p>

        {/* CTA buttons */}
        {!user ? (
          <div className="flex flex-wrap items-center gap-2 pt-1">
            <Link
              href={`/auth?redirect=/products/${productId}`}
              className="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl text-xs font-black text-white bg-orange-600 hover:bg-orange-700 shadow-md shadow-orange-600/20 transition-all active:scale-95"
            >
              <User size={13} />
              <span>{(t.ui as any).sign_in_to_view_specs ?? 'Sign In to View'}</span>
            </Link>
            <Link
              href={`/auth?as=buyer&redirect=/products/${productId}`}
              className="inline-flex items-center gap-1 px-3 py-2 rounded-xl text-xs font-bold text-zinc-700 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
            >
              <span>{(t as any).auth?.register_now ?? 'Create Account'}</span>
              <ArrowRight size={12} />
            </Link>
          </div>
        ) : (
          <div className="flex items-center gap-3 pt-1">
            <Link
              href="/buyer/dashboard"
              className={`inline-flex items-center gap-1.5 px-4 py-2 rounded-xl text-xs font-black text-white shadow-md transition-all active:scale-95 ${tierBadgeStyle}`}
            >
              <Sparkles size={13} />
              <span>
                {isGold
                  ? (t.ui as any).upgrade_to_gold ?? 'Upgrade to Gold'
                  : (t.ui as any).upgrade_to_premium ?? 'Upgrade to Premium'}
              </span>
              <ArrowRight size={13} />
            </Link>
            <span className="text-[11px] text-zinc-500">
              Current: <strong className="capitalize">{buyerTier || 'Normal'}</strong>
            </span>
          </div>
        )}
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
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-black uppercase tracking-wider text-zinc-900 dark:text-white">
          {(t.ui as any).specifications ?? 'Specifications'}
        </h3>
        <span className="text-[10px] font-bold text-green-600 dark:text-green-400 bg-green-500/10 px-2 py-0.5 rounded-full flex items-center gap-1">
          <Shield size={10} /> Verified Specs
        </span>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
        {attributes.map(attr => (
          <div
            key={attr.id}
            className="flex items-center justify-between px-3.5 py-2.5 bg-zinc-50 dark:bg-zinc-850 rounded-xl border border-zinc-100 dark:border-zinc-800"
          >
            <span className="text-xs text-zinc-500 font-semibold">{getAttrName(attr)}</span>
            <span className="text-xs font-bold text-zinc-900 dark:text-white">
              {getDisplayValue(attr)}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
