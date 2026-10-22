'use client';

import React from 'react';
import { X, Package, Home, ArrowRight } from 'lucide-react';
import { useRouter } from 'next/navigation';

interface MultiStepListingWizardProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: (newListing: any) => void;
}

/**
 * MultiStepListingWizard — Admin quick-access to listing creation.
 * Redirects to the appropriate admin product page rather than generating fake data.
 */
export function MultiStepListingWizard({ isOpen, onClose, onSuccess }: MultiStepListingWizardProps) {
  const router = useRouter();

  if (!isOpen) return null;

  const options = [
    {
      label: 'Product / Sale Listing',
      description: 'Add a product for sale to the marketplace',
      icon: Package,
      href: '/admin/products',
      color: 'border-orange-500/30 hover:border-orange-500/60 text-orange-400',
      bg: 'bg-orange-600/10',
    },
    {
      label: 'Rental / Property Listing',
      description: 'Add a property or rental to the marketplace',
      icon: Home,
      href: '/admin/products',
      color: 'border-teal-500/30 hover:border-teal-500/60 text-teal-400',
      bg: 'bg-teal-600/10',
    },
  ];

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-md">
      <div className="w-full max-w-md bg-zinc-950 border border-zinc-800 rounded-2xl shadow-2xl overflow-hidden">
        <div className="flex items-center justify-between px-5 py-4 border-b border-zinc-800">
          <h2 className="font-black text-white text-sm">Add New Listing</h2>
          <button onClick={onClose} className="p-1.5 rounded-lg text-zinc-500 hover:text-white hover:bg-zinc-800 transition-colors">
            <X size={16} />
          </button>
        </div>

        <div className="p-5 space-y-3">
          <p className="text-xs text-zinc-400">Choose the type of listing to add:</p>
          {options.map(opt => {
            const Icon = opt.icon;
            return (
              <button key={opt.label} onClick={() => { router.push(opt.href); onClose(); }}
                className={`w-full flex items-center gap-4 p-4 rounded-xl border bg-zinc-900 transition-all ${opt.color}`}>
                <div className={`w-10 h-10 rounded-xl flex items-center justify-center shrink-0 ${opt.bg}`}>
                  <Icon size={18} className={opt.color.split(' ')[0].replace('border-', 'text-').replace('/30', '')} />
                </div>
                <div className="text-left flex-1 min-w-0">
                  <p className="font-bold text-white text-sm">{opt.label}</p>
                  <p className="text-xs text-zinc-400 mt-0.5">{opt.description}</p>
                </div>
                <ArrowRight size={14} className="text-zinc-600 shrink-0" />
              </button>
            );
          })}
        </div>

        <div className="px-5 pb-4">
          <p className="text-[10px] text-zinc-600 text-center">
            Listings added through admin panel go through normal moderation flow
          </p>
        </div>
      </div>
    </div>
  );
}

