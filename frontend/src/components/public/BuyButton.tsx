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
      {open && <BuyModal product={product} onClose={() => setOpen(false)} />}
    </>
  );
}
