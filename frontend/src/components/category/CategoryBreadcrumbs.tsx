'use client';

import { useEffect, useState } from 'react';
import { ChevronRight, Home } from 'lucide-react';
import Link from 'next/link';

interface Breadcrumb {
  id: number;
  name: string;
  slug: string;
  level: number;
}

interface CategoryBreadcrumbsProps {
  categoryId: number;
  language?: string;
  className?: string;
}

export default function CategoryBreadcrumbs({
  categoryId,
  language = 'en',
  className = ''
}: CategoryBreadcrumbsProps) {
  const [breadcrumbs, setBreadcrumbs] = useState<Breadcrumb[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchBreadcrumbs();
  }, [categoryId, language]);

  const fetchBreadcrumbs = async () => {
    try {
      const res = await fetch(`/api/categories/${categoryId}/ancestors?lang=${language}`);
      const data = await res.json();
      if (data.success) {
        setBreadcrumbs(data.data || []);
      }
    } catch (error) {
      console.error('Failed to fetch breadcrumbs:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className={`flex items-center gap-2 ${className}`}>
        <div className="h-4 w-32 bg-gray-200 rounded animate-pulse" />
      </div>
    );
  }

  if (breadcrumbs.length === 0) {
    return null;
  }

  return (
    <nav className={`flex items-center gap-2 text-sm ${className}`} aria-label="Breadcrumb">
      <Link
        href="/find"
        className="text-gray-500 hover:text-gray-700 flex items-center gap-1"
      >
        <Home className="w-4 h-4" />
        <span>Home</span>
      </Link>

      {breadcrumbs.map((crumb, index) => (
        <div key={crumb.id} className="flex items-center gap-2">
          <ChevronRight className="w-4 h-4 text-gray-400" />
          {index === breadcrumbs.length - 1 ? (
            <span className="text-gray-900 font-medium">{crumb.name}</span>
          ) : (
            <Link
              href={`/find?category=${crumb.id}`}
              className="text-gray-500 hover:text-gray-700"
            >
              {crumb.name}
            </Link>
          )}
        </div>
      ))}
    </nav>
  );
}
