'use client';

import React, { useState, useRef, useEffect, useCallback } from 'react';
import { X, ZoomIn, ZoomOut, RotateCw, Check, AlertCircle, Scissors } from 'lucide-react';

interface Props {
  file: File;
  cropShape?: 'rect' | 'round'; // 'round' for profile photo, 'rect' for product images
  aspectRatio?: number; // e.g. 1 for square/round, 16/9, 4/3 for rect. Defaults: round=1, rect=4/3
  maxOutputMB?: number;
  onCropComplete: (croppedFile: File, previewUrl: string) => void;
  onClose: () => void;
}

export function ImageCropperModal({
  file,
  cropShape = 'rect',
  aspectRatio = cropShape === 'round' ? 1 : 4 / 3,
  maxOutputMB = 5,
  onCropComplete,
  onClose,
}: Props) {
  const [imageSrc, setImageSrc] = useState<string>('');
  const [scale, setScale] = useState<number>(1);
  const [rotation, setRotation] = useState<number>(0);
  const [position, setPosition] = useState<{ x: number; y: number }>({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState<{ x: number; y: number }>({ x: 0, y: 0 });
  const [processing, setProcessing] = useState(false);

  const containerRef = useRef<HTMLDivElement>(null);
  const imageRef = useRef<HTMLImageElement>(null);

  const fileSizeMB = (file.size / (1024 * 1024)).toFixed(1);
  const isLarge = file.size > maxOutputMB * 1024 * 1024;

  // Load image preview from file
  useEffect(() => {
    const objectUrl = URL.createObjectURL(file);
    setImageSrc(objectUrl);
    return () => URL.revokeObjectURL(objectUrl);
  }, [file]);

  // Handle Drag / Pan
  const handlePointerDown = (e: React.PointerEvent) => {
    setIsDragging(true);
    setDragStart({ x: e.clientX - position.x, y: e.clientY - position.y });
    (e.target as HTMLElement).setPointerCapture(e.pointerId);
  };

  const handlePointerMove = (e: React.PointerEvent) => {
    if (!isDragging) return;
    setPosition({
      x: e.clientX - dragStart.x,
      y: e.clientY - dragStart.y,
    });
  };

  const handlePointerUp = (e: React.PointerEvent) => {
    setIsDragging(false);
    try {
      (e.target as HTMLElement).releasePointerCapture(e.pointerId);
    } catch {}
  };

  // Perform HTML5 Canvas Crop
  const handleCrop = useCallback(async () => {
    if (!imageRef.current) return;
    setProcessing(true);

    try {
      const img = imageRef.current;
      const naturalWidth = img.naturalWidth;
      const naturalHeight = img.naturalHeight;

      // Desired output dimension
      const outputWidth = cropShape === 'round' ? 600 : 1200;
      const outputHeight = Math.round(outputWidth / aspectRatio);

      const canvas = document.createElement('canvas');
      canvas.width = outputWidth;
      canvas.height = outputHeight;
      const ctx = canvas.getContext('2d');

      if (!ctx) {
        throw new Error('Canvas context not available');
      }

      ctx.imageSmoothingEnabled = true;
      ctx.imageSmoothingQuality = 'high';

      // Transform context for zoom and rotation
      ctx.translate(outputWidth / 2, outputHeight / 2);
      ctx.rotate((rotation * Math.PI) / 180);
      ctx.scale(scale, scale);

      // Map drag translation to canvas space
      const renderRatio = naturalWidth / (img.width || 300);
      const canvasTranslateX = position.x * renderRatio;
      const canvasTranslateY = position.y * renderRatio;

      ctx.drawImage(
        img,
        -naturalWidth / 2 + canvasTranslateX,
        -naturalHeight / 2 + canvasTranslateY,
        naturalWidth,
        naturalHeight
      );

      // If round shape, we can keep the file rectangular (clean transparent or black background)
      // while displaying it circular via CSS, exactly as requested: "profile cercular cut and image is rectangular actual file"
      canvas.toBlob(
        (blob) => {
          if (!blob) {
            setProcessing(false);
            return;
          }

          const croppedFile = new File([blob], file.name.replace(/\.[^/.]+$/, '') + '_cropped.webp', {
            type: 'image/webp',
          });
          const previewUrl = URL.createObjectURL(blob);
          onCropComplete(croppedFile, previewUrl);
          setProcessing(false);
          onClose();
        },
        'image/webp',
        0.88
      );
    } catch (err) {
      setProcessing(false);
      console.error('Crop failed', err);
    }
  }, [aspectRatio, cropShape, file.name, onCropComplete, onClose, position.x, position.y, rotation, scale]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-4 bg-black/80 backdrop-blur-md animate-in fade-in duration-200">
      <div className="relative w-full max-w-lg bg-zinc-900 border border-zinc-800 rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[95vh]">
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-zinc-800">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-xl bg-orange-600/20 text-orange-500 flex items-center justify-center">
              <Scissors size={16} />
            </div>
            <div>
              <h3 className="text-sm font-black text-white">
                {cropShape === 'round' ? 'Circular Profile Photo Crop' : 'Rectangular Image Crop'}
              </h3>
              <p className="text-[10px] text-zinc-400">Position & crop to ensure crisp display</p>
            </div>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="w-8 h-8 rounded-full bg-zinc-800 hover:bg-zinc-700 text-zinc-400 hover:text-white flex items-center justify-center transition-colors"
          >
            <X size={16} />
          </button>
        </div>

        {/* Advisory banner for large files */}
        {isLarge && (
          <div className="bg-amber-950/40 border-b border-amber-800/60 px-4 py-2 flex items-start gap-2 text-xs text-amber-300">
            <AlertCircle size={14} className="shrink-0 mt-0.5 text-amber-400" />
            <p>
              Photo is <strong>{fileSizeMB} MB</strong>. Cropping here automatically produces an optimized, high-res file under 5 MB. Tip: Taking a phone screenshot is also a quick way to get images under 1 MB.
            </p>
          </div>
        )}

        {/* Viewport Canvas Frame */}
        <div className="relative flex-1 bg-black overflow-hidden flex items-center justify-center min-h-[300px] select-none touch-none">
          {/* Draggable Image */}
          {imageSrc && (
            <div
              onPointerDown={handlePointerDown}
              onPointerMove={handlePointerMove}
              onPointerUp={handlePointerUp}
              className="absolute inset-0 flex items-center justify-center cursor-move"
            >
              <img
                ref={imageRef}
                src={imageSrc}
                alt="To Crop"
                draggable={false}
                style={{
                  transform: `translate(${position.x}px, ${position.y}px) scale(${scale}) rotate(${rotation}deg)`,
                  maxWidth: '85%',
                  maxHeight: '85%',
                  objectFit: 'contain',
                  userSelect: 'none',
                }}
                className="transition-transform duration-75 pointer-events-none"
              />
            </div>
          )}

          {/* Mask Overlay */}
          <div className="absolute inset-0 pointer-events-none flex items-center justify-center">
            {cropShape === 'round' ? (
              <div
                style={{ width: '240px', height: '240px' }}
                className="rounded-full border-2 border-orange-500 shadow-[0_0_0_9999px_rgba(0,0,0,0.65)] ring-1 ring-white/30"
              />
            ) : (
              <div
                style={{ width: '280px', height: `${Math.round(280 / aspectRatio)}px` }}
                className="rounded-2xl border-2 border-orange-500 shadow-[0_0_0_9999px_rgba(0,0,0,0.65)] ring-1 ring-white/30"
              />
            )}
          </div>

          <div className="absolute bottom-3 left-3 bg-black/60 backdrop-blur-md px-2.5 py-1 rounded-lg text-[10px] text-zinc-300 pointer-events-none">
            Drag photo to reposition
          </div>
        </div>

        {/* Controls Toolbar */}
        <div className="px-5 py-3 bg-zinc-950 border-t border-zinc-800 space-y-3">
          <div className="flex items-center justify-between gap-4">
            <div className="flex items-center gap-2 flex-1">
              <ZoomOut size={14} className="text-zinc-500" />
              <input
                type="range"
                min="0.5"
                max="3"
                step="0.05"
                value={scale}
                onChange={(e) => setScale(parseFloat(e.target.value))}
                className="flex-1 accent-orange-500 h-1.5 bg-zinc-800 rounded-lg cursor-pointer"
              />
              <ZoomIn size={14} className="text-zinc-500" />
              <span className="text-[11px] font-mono text-zinc-400 w-10 text-right">
                {Math.round(scale * 100)}%
              </span>
            </div>

            <button
              type="button"
              onClick={() => setRotation((r) => (r + 90) % 360)}
              className="px-3 py-1.5 rounded-xl bg-zinc-800 hover:bg-zinc-700 text-zinc-300 text-xs font-bold flex items-center gap-1.5 transition-colors"
            >
              <RotateCw size={13} />
              <span>Rotate</span>
            </button>
          </div>
        </div>

        {/* Footer Actions */}
        <div className="px-5 py-4 border-t border-zinc-800 flex gap-3 bg-zinc-900">
          <button
            type="button"
            onClick={onClose}
            className="flex-1 py-3 rounded-xl bg-zinc-800 hover:bg-zinc-700 text-zinc-300 text-xs font-bold transition-colors"
          >
            Cancel
          </button>
          <button
            type="button"
            disabled={processing}
            onClick={handleCrop}
            className="flex-1 py-3 rounded-xl bg-orange-600 hover:bg-orange-500 disabled:opacity-50 text-white text-xs font-black shadow-lg shadow-orange-600/25 flex items-center justify-center gap-2 transition-all"
          >
            <Check size={14} />
            <span>{processing ? 'Applying Crop…' : 'Apply Crop'}</span>
          </button>
        </div>
      </div>
    </div>
  );
}
