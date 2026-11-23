'use client';
import { useRef, useState, useCallback } from 'react';
import Link from 'next/link';
import { X, Upload, AlertCircle, Loader2, ImageOff, Globe, Sparkles, Lock } from 'lucide-react';
import { MIN_IMAGES } from '@/lib/constants/fees';
import { useLocale } from '@/hooks/useLocale';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { ImageCropperModal } from '@/components/ui/ImageCropperModal';

const MAX_FILE_SIZE = 5 * 1024 * 1024; // Strictly 5 MB limit
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
  const [activeTab, setActiveTab] = useState<'upload' | 'url'>('upload');
  const [urlInput, setUrlInput] = useState('');
  const [urlImporting, setUrlImporting] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState('');
  const [brokenIndexes, setBrokenIndexes] = useState<Set<number>>(new Set());
  const [pendingCropFile, setPendingCropFile] = useState<File | null>(null);

  const isPrivileged =
    user?.role === 'admin' ||
    user?.tier === 'gold' ||
    user?.tier === 'premium' ||
    user?.buyerTier === 'gold' ||
    user?.buyerTier === 'premium';

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

  const uploadSingleFile = async (f: File) => {
    if (!user?.token) return;
    setUploading(true);
    setError('');
    try {
      const url = await api.uploadImage(f, user.token);
      if (!url || typeof url !== 'string' || url.trim() === '') {
        throw new Error('Upload failed');
      }
      onChange([...urls, url]);
    } catch (e: any) {
      setError(e.message ?? 'Upload failed');
    } finally {
      setUploading(false);
    }
  };

  const handleFiles = async (files: FileList | null) => {
    if (!files || files.length === 0) return;
    if (!user?.token) { setError('Not authenticated'); return; }

    const remaining = maxImages - urls.length;
    if (remaining <= 0) { setError(`Maximum ${maxImages} images reached`); return; }
    const toUpload = Array.from(files).slice(0, remaining);

    for (const f of toUpload) {
      if (f.size > MAX_FILE_SIZE) {
        const sz = (f.size / (1024 * 1024)).toFixed(1);
        setPendingCropFile(f);
        setError(`"${f.name}" is ${sz} MB (max 5 MB). Opening crop tool to resize, or take a phone screenshot.`);
        return;
      }
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

  const handleUrlImport = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user?.token) {
      setError('Not authenticated');
      return;
    }
    const cleanUrl = urlInput.trim();
    if (!cleanUrl) return;

    if (!isPrivileged) {
      setError('Importing from URL is exclusively for Gold and Premium members');
      return;
    }

    if (urls.length >= maxImages) {
      setError(`Maximum ${maxImages} images reached`);
      return;
    }

    setUrlImporting(true);
    setError('');
    try {
      const uploadedUrl = await api.importImageFromURL(cleanUrl, user.token);
      if (!uploadedUrl) throw new Error('Failed to import image');
      onChange([...urls, uploadedUrl]);
      setUrlInput('');
    } catch (err: any) {
      setError(err.message ?? 'Failed to import image from URL');
    } finally {
      setUrlImporting(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    handleFiles(e.dataTransfer.files);
  };

  const canAdd = urls.length < maxImages && !uploading && !urlImporting;

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <p className="text-xs font-black uppercase tracking-wider text-zinc-500 dark:text-white/40">
          {t.ui.upload_image ?? 'Images'}{' '}
          <span className="text-[#1769E0] normal-case font-normal">
            ({urls.length}/{maxImages} · {MIN_IMAGES} min)
          </span>
        </p>

        {/* Tab Selector */}
        <div className="inline-flex p-0.5 rounded-lg bg-zinc-100 dark:bg-zinc-800 text-[11px] font-medium">
          <button
            type="button"
            onClick={() => setActiveTab('upload')}
            className={`px-2.5 py-1 rounded-md transition-all ${
              activeTab === 'upload'
                ? 'bg-white dark:bg-zinc-700 text-zinc-900 dark:text-white shadow-xs font-bold'
                : 'text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'
            }`}
          >
            File Upload
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('url')}
            className={`px-2.5 py-1 rounded-md inline-flex items-center gap-1 transition-all ${
              activeTab === 'url'
                ? 'bg-white dark:bg-zinc-700 text-zinc-900 dark:text-white shadow-xs font-bold'
                : 'text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'
            }`}
          >
            <Globe size={11} />
            URL Import
            {!isPrivileged && <Lock size={10} className="text-amber-500" />}
          </button>
        </div>
      </div>

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

      {canAdd && activeTab === 'upload' && (
        <div onDrop={handleDrop} onDragOver={e => e.preventDefault()}
          onClick={() => inputRef.current?.click()}
          className="flex flex-col items-center justify-center gap-2 border-2 border-dashed border-zinc-300 dark:border-zinc-700 rounded-xl p-6 cursor-pointer hover:border-[#1769E0]/60 hover:bg-[#1769E0]/5 transition-all">
          {uploading ? (
            <>
              <Loader2 size={20} className="text-[#1769E0] animate-spin" />
              <p className="text-xs text-zinc-500 dark:text-white/40">Uploading…</p>
            </>
          ) : (
            <>
              <Upload size={20} className="text-zinc-400 dark:text-white/30" />
              <p className="text-xs font-bold text-zinc-600 dark:text-white/60">Click or drag images here</p>
              <p className="text-[10px] text-zinc-400 dark:text-white/25">JPEG · PNG · WebP · Max 5 MB each</p>
            </>
          )}
        </div>
      )}

      {canAdd && activeTab === 'url' && (
        <div className="border border-zinc-200 dark:border-zinc-800 rounded-xl p-4 bg-zinc-50/60 dark:bg-zinc-900/60 space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1.5 text-xs font-bold text-zinc-800 dark:text-zinc-200">
              <Globe size={14} className="text-[#1769E0]" />
              Import Image from External Web URL
            </div>
            <span className="inline-flex items-center gap-1 text-[10px] font-semibold text-amber-600 bg-amber-500/10 px-2 py-0.5 rounded-full border border-amber-500/20">
              <Sparkles size={10} /> Gold & Premium Feature
            </span>
          </div>

          {!isPrivileged ? (
            <div className="p-3 bg-amber-50 dark:bg-amber-950/30 border border-amber-200 dark:border-amber-900/40 rounded-lg text-xs text-amber-800 dark:text-amber-300 space-y-1">
              <p className="font-semibold flex items-center gap-1">
                <Lock size={12} /> Member Tier Required
              </p>
              <p className="text-[11px] leading-relaxed">
                Direct URL image importing with automated SSRF protection is reserved for Gold and Premium members.
              </p>
              <div className="pt-1">
                <Link
                  href="/seller/pay"
                  className="inline-block text-[11px] font-bold text-[#1769E0] hover:underline"
                >
                  Upgrade to Gold/Premium →
                </Link>
              </div>
            </div>
          ) : (
            <form onSubmit={handleUrlImport} className="flex gap-2">
              <input
                type="url"
                value={urlInput}
                onChange={e => setUrlInput(e.target.value)}
                placeholder="https://example.com/item-photo.jpg"
                disabled={urlImporting}
                className="flex-1 text-xs px-3 py-2 rounded-lg border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-[#1769E0]/40"
              />
              <button
                type="submit"
                disabled={urlImporting || !urlInput.trim()}
                className="bg-[#1769E0] hover:bg-[#1055BD] text-white text-xs font-bold px-3 py-2 rounded-lg disabled:opacity-50 inline-flex items-center gap-1 min-h-[36px]"
              >
                {urlImporting ? <Loader2 size={13} className="animate-spin" /> : 'Fetch & Save'}
              </button>
            </form>
          )}
        </div>
      )}

      {/* Advisory tip */}
      <p className="text-[10px] text-zinc-400 dark:text-zinc-500 leading-relaxed">
        💡 <strong>Pro Tip:</strong> Mobile screenshots maintain crisp visual detail while keeping file sizes small (typically under 1 MB). Photos larger than 5 MB will open the crop & resize tool automatically.
      </p>

      <input ref={inputRef} type="file" accept={ACCEPTED_ATTR} multiple className="hidden"
        onChange={e => handleFiles(e.target.files)} />

      {/* Crop Modal when a file needs cropping or is oversized */}
      {pendingCropFile && (
        <ImageCropperModal
          file={pendingCropFile}
          cropShape="rect"
          aspectRatio={4 / 3}
          onCropComplete={(croppedFile) => {
            setPendingCropFile(null);
            uploadSingleFile(croppedFile);
          }}
          onClose={() => setPendingCropFile(null)}
        />
      )}

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
