'use client';
import { useRef, useState, useCallback } from 'react';
import { X, Upload, AlertCircle, Loader2, ImageOff } from 'lucide-react';
import { MIN_IMAGES } from '@/lib/constants/fees';
import { useLocale } from '@/hooks/useLocale';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';

const MAX_FILE_SIZE = 20 * 1024 * 1024; // 20 MB
const ACCEPTED_MIME = ['image/jpeg', 'image/png', 'image/webp'];
const ACCEPTED_ATTR = 'image/jpeg,image/png,image/webp';

function BrokenImage({ index }: { index: number }) {
  return (
    <div className="w-full h-full flex flex-col items-center justify-center gap-1 bg-zinc-100 dark:bg-zinc-800 text-zinc-400 dark:text-zinc-600">
      <ImageOff size={18} />
      <span className="text-[9px]">#{index + 1}</span>
    </div>
  );
}

interface Props {
  urls: string[];
  onChange: (urls: string[]) => void;
  maxImages?: number;
}

export function ImageUploader({ urls, onChange, maxImages = 20 }: Props) {
  const { t } = useLocale();
  const { user } = useAuth();
  const inputRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState('');
  const [brokenIndexes, setBrokenIndexes] = useState<Set<number>>(new Set());

  const markBroken = useCallback((i: number) => {
    setBrokenIndexes(prev => new Set(prev).add(i));
  }, []);

  const remove = (i: number) => {
    if (urls.length <= MIN_IMAGES) return;
    onChange(urls.filter((_, idx) => idx !== i));
    setBrokenIndexes(prev => {
      const next = new Set<number>();
      prev.forEach(idx => { if (idx !== i) next.add(idx > i ? idx - 1 : idx); });
      return next;
    });
  };

  const handleFiles = async (files: FileList | null) => {
    if (!files || files.length === 0) return;
    if (!user?.token) { setError('Not authenticated'); return; }

    const remaining = maxImages - urls.length;
    if (remaining <= 0) { setError(`Maximum ${maxImages} images reached`); return; }
    const toUpload = Array.from(files).slice(0, remaining);

    for (const f of toUpload) {
      if (f.size > MAX_FILE_SIZE) { setError(`"${f.name}" exceeds 20 MB limit`); return; }
      if (!ACCEPTED_MIME.includes(f.type)) { setError(`"${f.name}" is not a JPEG, PNG or WebP image`); return; }
      if (f.size === 0) { setError(`"${f.name}" is empty`); return; }
    }

    setError('');
    setUploading(true);
    try {
      const uploaded: string[] = [];
      for (const f of toUpload) {
        const url = await api.uploadImage(f, user.token);
        if (!url || typeof url !== 'string' || url.trim() === '') {
          throw new Error('Server returned an empty URL — upload may have failed');
        }
        uploaded.push(url);
      }
      onChange([...urls, ...uploaded]);
    } catch (e: any) {
      setError(e.message ?? 'Upload failed');
    } finally {
      setUploading(false);
      if (inputRef.current) inputRef.current.value = '';
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    handleFiles(e.dataTransfer.files);
  };

  const canAdd = urls.length < maxImages && !uploading;

  return (
    <div className="space-y-3">
      <p className="text-xs font-black uppercase tracking-wider text-zinc-500 dark:text-white/40">
        {t.ui.upload_image ?? 'Images'}{' '}
        <span className="text-orange-600 normal-case font-normal">
          ({urls.length}/{maxImages} · {MIN_IMAGES} min)
        </span>
      </p>

      {urls.length > 0 && (
        <div className="grid grid-cols-3 sm:grid-cols-4 gap-2">
          {urls.map((url, i) => (
            <div key={`${url}-${i}`}
              className="relative group aspect-square rounded-xl overflow-hidden border border-zinc-200 dark:border-zinc-700 bg-zinc-100 dark:bg-zinc-800">
              {brokenIndexes.has(i) ? (
                <BrokenImage index={i} />
              ) : (
                // eslint-disable-next-line @next/next/no-img-element
                <img src={url} alt={`Product image ${i + 1}`}
                  className="w-full h-full object-cover"
                  onError={() => markBroken(i)}
                  loading="lazy" crossOrigin="anonymous" />
              )}
              {urls.length > MIN_IMAGES && (
                <button type="button" onClick={() => remove(i)}
                  className="absolute top-1 right-1 w-5 h-5 bg-red-500 text-white rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity shadow-md"
                  aria-label={`Remove image ${i + 1}`}>
                  <X size={10} />
                </button>
              )}
              {brokenIndexes.has(i) && urls.length > MIN_IMAGES && (
                <button type="button" onClick={() => remove(i)}
                  className="absolute top-1 right-1 w-5 h-5 bg-red-500 text-white rounded-full flex items-center justify-center shadow-md"
                  aria-label={`Remove broken image ${i + 1}`}>
                  <X size={10} />
                </button>
              )}
              <div className="absolute bottom-0 left-0 right-0 bg-black/40 text-white text-[9px] text-center py-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
                {i + 1}
              </div>
            </div>
          ))}
        </div>
      )}

      {canAdd && (
        <div onDrop={handleDrop} onDragOver={e => e.preventDefault()}
          onClick={() => inputRef.current?.click()}
          className="flex flex-col items-center justify-center gap-2 border-2 border-dashed border-zinc-300 dark:border-zinc-700 rounded-xl p-6 cursor-pointer hover:border-orange-500/60 hover:bg-orange-500/5 transition-all">
          {uploading ? (
            <>
              <Loader2 size={20} className="text-orange-500 animate-spin" />
              <p className="text-xs text-zinc-500 dark:text-white/40">Uploading…</p>
            </>
          ) : (
            <>
              <Upload size={20} className="text-zinc-400 dark:text-white/30" />
              <p className="text-xs font-bold text-zinc-600 dark:text-white/60">Click or drag images here</p>
              <p className="text-[10px] text-zinc-400 dark:text-white/25">JPEG · PNG · WebP · up to 20 MB each</p>
            </>
          )}
        </div>
      )}

      <input ref={inputRef} type="file" accept={ACCEPTED_ATTR} multiple className="hidden"
        onChange={e => handleFiles(e.target.files)} />

      {error && (
        <div className="flex items-center gap-2 text-xs text-red-500 bg-red-500/10 border border-red-500/20 rounded-xl px-3 py-2">
          <AlertCircle size={12} className="shrink-0" /> {error}
        </div>
      )}

      {brokenIndexes.size > 0 && (
        <div className="flex items-center gap-2 text-xs text-amber-600 bg-amber-500/10 border border-amber-500/20 rounded-xl px-3 py-2">
          <AlertCircle size={12} className="shrink-0" />
          {brokenIndexes.size} image{brokenIndexes.size > 1 ? 's' : ''} could not be loaded — remove and re-upload
        </div>
      )}
    </div>
  );
}
