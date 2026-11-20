'use client';
import { useState } from 'react';
import { Product } from '@/types';
import { BuyModal } from './BuyModal';
import { useLocale } from '@/hooks/useLocale';
import { PhoneCall, MessageSquare, Key } from 'lucide-react';

export function BuyButton({ product }: { product: Product }) {
  const { t } = useLocale();
  const [open, setOpen] = useState(false);
  const isRental = product.listing_type === 'rent';

  const actionLabel = isRental 
    ? (t.ui.rent_now ?? 'Contact Landlord / Inquire') 
    : (t.ui.contact_seller ?? 'Contact Seller / Direct Inquiry');

  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        className={`mt-2 w-full text-white rounded-2xl py-4 font-bold text-lg transition-all shadow-lg tap-bounce active:scale-[0.98] flex items-center justify-center gap-2 min-h-[48px] ${
          isRental 
            ? 'bg-teal-600 hover:bg-teal-700 shadow-teal-600/20' 
            : 'bg-[#1769E0] hover:bg-[#1055BD] shadow-[#1769E0]/25'
        }`}
      >
        {isRental ? <Key size={20} /> : <PhoneCall size={20} />}
        {actionLabel}
      </button>

      {/* Mobile Sticky Priority Bottom Bar */}
      <div className="md:hidden fixed bottom-0 left-0 right-0 z-40 bg-white/95 dark:bg-[#091426]/95 backdrop-blur-xl border-t border-slate-200 dark:border-slate-800 px-4 py-2.5 pb-[calc(env(safe-area-inset-bottom,12px)+6px)] shadow-2xl flex items-center justify-between gap-3">
        <div className="flex flex-col min-w-0 pr-1">
          <span className="text-[10px] font-bold text-slate-400 uppercase tracking-wider leading-none">
            {isRental ? (t.ui.rental_rate ?? 'Rental Rate') : (t.ui.total_price ?? 'Listing Price')}
          </span>
          <span className={`text-base font-black truncate leading-tight mt-0.5 ${isRental ? 'text-teal-600' : 'text-[#1769E0] dark:text-[#38A5FF]'}`}>
            ETB {isRental && product.rental_price != null ? product.rental_price.toLocaleString() : (product.price != null ? product.price.toLocaleString() : '—')}
          </span>
        </div>
        <button
          type="button"
          onClick={() => setOpen(true)}
          className={`flex-1 text-white rounded-xl py-3 px-4 font-black text-sm shadow-md transition-all tap-bounce active:scale-95 flex items-center justify-center gap-2 min-h-[44px] ${
            isRental 
              ? 'bg-teal-600 hover:bg-teal-700 shadow-teal-600/25' 
              : 'bg-[#1769E0] hover:bg-[#1055BD] shadow-[#1769E0]/25'
          }`}
        >
          {isRental ? <Key size={16} /> : <PhoneCall size={16} />}
          <span>{actionLabel}</span>
        </button>
      </div>

      {open && <BuyModal product={product} onClose={() => setOpen(false)} />}
    </>
  );
}
