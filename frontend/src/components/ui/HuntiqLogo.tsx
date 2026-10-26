'use client';

import Link from 'next/link';
import Image from 'next/image';

export interface HuntiqLogoProps {
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  showText?: boolean;
  subtitle?: string;
  className?: string;
  href?: string;
  priority?: boolean;
}

const SIZE_MAP = {
  xs: { dim: 24, px: 'w-6 h-6', text: 'text-sm', sub: 'text-[8px]' },
  sm: { dim: 32, px: 'w-8 h-8', text: 'text-base', sub: 'text-[9px]' },
  md: { dim: 40, px: 'w-10 h-10', text: 'text-lg', sub: 'text-[10px]' },
  lg: { dim: 48, px: 'w-12 h-12', text: 'text-xl', sub: 'text-xs' },
  xl: { dim: 64, px: 'w-16 h-16', text: 'text-2xl', sub: 'text-sm' },
};

export function HuntiqLogo({
  size = 'md',
  showText = false,
  subtitle,
  className = '',
  href,
  priority = false,
}: HuntiqLogoProps) {
  const config = SIZE_MAP[size] || SIZE_MAP.md;

  const content = (
    <div className={`flex items-center gap-2.5 shrink-0 group ${className}`}>
      <div
        className={`relative ${config.px} rounded-xl overflow-hidden flex items-center justify-center bg-zinc-950 border border-zinc-800 shadow-sm group-hover:scale-105 transition-transform duration-200 shrink-0`}
      >
        <Image
          src="/huntiq.png"
          alt="HUNTIQ"
          width={config.dim}
          height={config.dim}
          priority={priority}
          className="w-full h-full object-contain p-0.5"
        />
      </div>

      {showText && (
        <div className="flex flex-col select-none">
          <span
            className={`${config.text} font-black tracking-tight leading-none text-zinc-900 dark:text-white group-hover:text-orange-600 transition-colors`}
          >
            HUNTIQ
          </span>
          {subtitle && (
            <span
              className={`${config.sub} font-bold text-orange-600 uppercase tracking-widest leading-tight mt-0.5`}
            >
              {subtitle}
            </span>
          )}
        </div>
      )}
    </div>
  );

  if (href) {
    return (
      <Link href={href} className="inline-flex shrink-0">
        {content}
      </Link>
    );
  }

  return content;
}
