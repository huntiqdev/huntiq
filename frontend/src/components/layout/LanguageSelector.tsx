'use client';
import { useState, useRef, useEffect } from 'react';
import { useLocaleContext } from '@/context/LocaleContext';
import { Globe, ChevronDown, Check } from 'lucide-react';

interface LanguageSelectorProps {
  variant?: 'compact' | 'full' | 'pill';
  className?: string;
  dropDirection?: 'down' | 'up';
}

export function LanguageSelector({
  variant = 'compact',
  className = '',
  dropDirection = 'down',
}: LanguageSelectorProps) {
  const { locale, setLocale, locales, localeLabels } = useLocaleContext();
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }
    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen]);

  const shortLabels: Record<string, string> = {
    en: 'EN',
    am: 'አማ',
    or: 'ORO',
    ti: 'ትግ',
  };

  const fullLabels: Record<string, string> = {
    en: 'English',
    am: 'አማርኛ (Amharic)',
    or: 'Afaan Oromoo',
    ti: 'ትግርኛ (Tigrinya)',
  };

  return (
    <div ref={containerRef} className={`relative inline-block text-left ${className}`}>
      {variant === 'pill' ? (
        <button
          type="button"
          onClick={() => setIsOpen(v => !v)}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-full border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-800/80 hover:border-[#1769E0] dark:hover:border-[#1769E0] text-zinc-700 dark:text-zinc-200 text-xs font-bold transition-all shadow-sm cursor-pointer"
          aria-label="Change Language"
        >
          <Globe size={13} className="text-[#1769E0]" />
          <span>{shortLabels[locale] ?? locale.toUpperCase()}</span>
          <ChevronDown size={11} className={`transition-transform duration-200 ${isOpen ? 'rotate-180' : ''}`} />
        </button>
      ) : variant === 'full' ? (
        <button
          type="button"
          onClick={() => setIsOpen(v => !v)}
          className="w-full flex items-center justify-between px-3 py-2 rounded-xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-800/50 hover:bg-zinc-50 dark:hover:bg-zinc-800 text-zinc-700 dark:text-zinc-200 text-xs font-bold transition-all cursor-pointer"
          aria-label="Change Language"
        >
          <div className="flex items-center gap-2">
            <Globe size={14} className="text-[#1769E0]" />
            <span>{fullLabels[locale] ?? localeLabels[locale] ?? locale.toUpperCase()}</span>
          </div>
          <ChevronDown size={12} className={`transition-transform duration-200 ${isOpen ? 'rotate-180' : ''}`} />
        </button>
      ) : (
        <button
          type="button"
          onClick={() => setIsOpen(v => !v)}
          className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl text-zinc-600 dark:text-zinc-300 hover:text-[#1769E0] dark:hover:text-[#38A5FF] hover:bg-zinc-100 dark:hover:bg-zinc-800 text-xs font-bold transition-all cursor-pointer border border-transparent hover:border-zinc-200 dark:hover:border-zinc-700"
          aria-label="Change Language"
        >
          <Globe size={13} className="text-[#1769E0]" />
          <span className="font-extrabold uppercase tracking-wider text-[11px]">{shortLabels[locale] ?? locale.toUpperCase()}</span>
          <ChevronDown size={10} className={`transition-transform duration-200 text-zinc-400 ${isOpen ? 'rotate-180' : ''}`} />
        </button>
      )}

      {isOpen && (
        <div
          className={`absolute right-0 ${
            dropDirection === 'up' ? 'bottom-full mb-1.5' : 'top-full mt-1.5'
          } w-44 rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-2xl z-50 overflow-hidden py-1.5 animate-in fade-in-50 zoom-in-95 duration-150`}
        >
          <div className="px-3 py-1.5 border-b border-zinc-100 dark:border-zinc-800/80">
            <p className="text-[10px] font-black uppercase tracking-widest text-zinc-400">
              Select Language / ቋንቋ
            </p>
          </div>
          <div className="py-1">
            {locales.map((l: string) => {
              const active = locale === l;
              return (
                <button
                  key={l}
                  type="button"
                  onClick={() => {
                    setLocale(l as any);
                    setIsOpen(false);
                  }}
                  className={`w-full text-left px-3.5 py-2 text-xs font-bold transition-colors flex items-center justify-between cursor-pointer ${
                    active
                      ? 'bg-[#1769E0]/10 text-[#1769E0] dark:text-[#38A5FF]'
                      : 'text-zinc-700 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800'
                  }`}
                >
                  <div className="flex items-center gap-2">
                    <span className="font-mono text-[10px] uppercase font-black text-zinc-400">
                      {shortLabels[l] ?? l.toUpperCase()}
                    </span>
                    <span>{fullLabels[l] ?? localeLabels[l] ?? l.toUpperCase()}</span>
                  </div>
                  {active && <Check size={13} className="text-[#1769E0] dark:text-[#38A5FF]" />}
                </button>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}
