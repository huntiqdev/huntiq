'use client';

import React, { createContext, useContext, useEffect, useState, useTransition } from 'react';

export type Theme = 'light' | 'dark' | 'system';
export type ResolvedTheme = 'light' | 'dark';

interface ThemeContextType {
  theme: Theme;
  resolvedTheme: ResolvedTheme;
  setTheme: (theme: Theme) => void;
  toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType>({
  theme: 'system',
  resolvedTheme: 'light',
  setTheme: () => {},
  toggleTheme: () => {},
});

const STORAGE_KEY = 'huntiq_theme';
const LEGACY_STORAGE_KEY = 'habesha_theme';

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setThemeState] = useState<Theme>('system');
  const [resolvedTheme, setResolvedTheme] = useState<ResolvedTheme>('light');
  const [, startTransition] = useTransition();

  // Initialize theme from storage or system preference
  useEffect(() => {
    let saved: string | null = null;
    try {
      saved = localStorage.getItem(STORAGE_KEY) || localStorage.getItem(LEGACY_STORAGE_KEY);
    } catch {
      // ignore
    }

    const initialTheme: Theme = saved === 'light' || saved === 'dark' ? saved : 'system';
    setThemeState(initialTheme);

    const isSystemDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const computedResolved: ResolvedTheme = initialTheme === 'system' ? (isSystemDark ? 'dark' : 'light') : initialTheme;
    setResolvedTheme(computedResolved);

    applyThemeClass(computedResolved);
  }, []);

  // Listen to OS system preference changes
  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleChange = (e: MediaQueryListEvent) => {
      if (theme === 'system') {
        const next: ResolvedTheme = e.matches ? 'dark' : 'light';
        setResolvedTheme(next);
        applyThemeClass(next);
      }
    };

    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, [theme]);

  const applyThemeClass = (target: ResolvedTheme) => {
    const root = document.documentElement;
    if (target === 'dark') {
      root.classList.add('dark');
    } else {
      root.classList.remove('dark');
    }
  };

  const setTheme = (next: Theme) => {
    startTransition(() => {
      setThemeState(next);
      try {
        localStorage.setItem(STORAGE_KEY, next);
        localStorage.setItem(LEGACY_STORAGE_KEY, next);
      } catch {
        // ignore
      }

      const isSystemDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      const computedResolved: ResolvedTheme = next === 'system' ? (isSystemDark ? 'dark' : 'light') : next;
      setResolvedTheme(computedResolved);
      applyThemeClass(computedResolved);
    });
  };

  const toggleTheme = () => {
    const next: Theme = resolvedTheme === 'dark' ? 'light' : 'dark';
    setTheme(next);
  };

  return (
    <ThemeContext.Provider value={{ theme, resolvedTheme, setTheme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  return useContext(ThemeContext);
}

