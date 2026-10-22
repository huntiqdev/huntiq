import { ProductImage } from '@/types';
import { ShieldCheck, AlertCircle } from 'lucide-react';
import { resolveImageUrl } from '@/lib/utils/imageUrl';

interface Props { images: ProductImage[] }

export function ImageVerification({ images }: Props) {
  const hasMin = images.length >= 2;

  return (
    <div className="flex flex-col gap-2">
      <div className={`flex items-center gap-2 text-xs font-bold ${hasMin ? 'text-green-600' : 'text-red-500'}`}>
        {hasMin ? <ShieldCheck size={14} /> : <AlertCircle size={14} />}
        {images.length} image{images.length !== 1 ? 's' : ''} {hasMin ? '— meets minimum' : '— needs 2+'}
      </div>
      <div className="flex gap-2 flex-wrap">
        {images.map(img => (
          <img key={img.id} src={resolveImageUrl(img.image_url)} alt={`Image ${img.order_number}`}
            className="w-16 h-16 rounded-xl object-cover border border-zinc-200 dark:border-zinc-700" />
        ))}
      </div>
    </div>
  );
}
