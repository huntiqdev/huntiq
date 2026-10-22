import { Product } from '@/types';
import { ProductCard } from '@/components/public/ProductCard';
import { ShoppingBag } from 'lucide-react';

interface Props {
  products: Product[];
  emptyMessage?: string;
}

export function ProductGrid({ products, emptyMessage = 'No products found.' }: Props) {
  if (products.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-20 opacity-50">
        <ShoppingBag size={64} className="mb-4" />
        <p className="text-xl font-medium">{emptyMessage}</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-2 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {products.map(p => <ProductCard key={p.id} product={p} />)}
    </div>
  );
}
