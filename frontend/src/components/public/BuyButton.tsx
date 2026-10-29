'use client';
import { useState } from 'react';
import { Product } from '@/types';
import { BuyModal } from './BuyModal';
import { useLocale } from '@/hooks/useLocale';
import { ShoppingBag, Key } from 'lucide-react';

export function BuyButton({ product }: { product: Product }) {
  const { t } = useLocale();
  const [open, setOpen] = useState(false);
  const isRental = product.listing_type === 'rent';

  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        className={`mt-2 w-full text-white rounded-2xl py-4 font-bold text-lg transition-all shadow-lg active:scale-95 flex items-center justify-center gap-2 ${
          isRental 
            ? 'bg-teal-600 hover:bg-teal-700 shadow-teal-600/20' 
            : 'bg-orange-600 hover:bg-orange-700 shadow-orange-600/20'
        }`}
      >
        {isRental ? <Key size={20} /> : <ShoppingBag size={20} />}
        {isRental ? (t.ui.rent_now ?? 'Rent / Contact Landlord') : t.ui.buy_now}
      </button>

      {/* Mobile Sticky Priority Bottom Bar */}
      <div className="md:hidden fixed bottom-0 left-0 right-0 z-40 bg-white/95 dark:bg-zinc-950/95 backdrop-blur-lg border-t border-zinc-200 dark:border-zinc-800 px-4 py-2.5 pb-[calc(env(safe-area-inset-bottom,12px)+6px)] shadow-2xl flex items-center justify-between gap-3">
        <div className="flex flex-col min-w-0 pr-1">
          <span className="text-[10px] font-bold text-zinc-400 uppercase tracking-wider leading-none">
            {isRental ? 'Rental Rate' : 'Total Price'}
          </span>
          <span className={`text-base font-black truncate leading-tight mt-0.5 ${isRental ? 'text-teal-600' : 'text-orange-600'}`}>
            ETB {isRental && product.rental_price ? product.rental_price.toLocaleString() : product.price.toLocaleString()}
          </span>
        </div>
        <button
          type="button"
          onClick={() => setOpen(true)}
          className={`flex-1 text-white rounded-xl py-3 px-4 font-black text-sm shadow-md transition-all active:scale-95 flex items-center justify-center gap-2 ${
            isRental 
              ? 'bg-teal-600 hover:bg-teal-700 shadow-teal-600/25' 
              : 'bg-orange-600 hover:bg-orange-700 shadow-orange-600/25'
          }`}
        >
          {isRental ? <Key size={16} /> : <ShoppingBag size={16} />}
          <span>{isRental ? (t.ui.rent_now ?? 'Rent Now') : (t.ui.buy_now ?? 'Buy Now')}</span>
        </button>
      </div>

      {open && <BuyModal product={product} onClose={() => setOpen(false)} />}
    </>
  );
}
