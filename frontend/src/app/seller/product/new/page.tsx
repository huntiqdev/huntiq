'use client';
import { Suspense, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { useLocaleContext } from '@/context/LocaleContext';
import { ProductForm } from '@/components/seller/ProductForm';
import { ArrowLeft } from 'lucide-react';

function NewProductContent() {
  const router = useRouter();
  const { user, loading, isSeller, isLandlord, activateCapability } = useAuth();
  const { t } = useLocaleContext();
  const searchParams = useSearchParams();
  const listingType = searchParams?.get('type') || 'sale';
  const isRental = listingType === 'rent';

  useEffect(() => {
    if (loading) return;
    if (!user) {
      router.push(`/auth?intent=${isRental ? 'rent' : 'sell'}`);
      return;
    }
    // Contextual capability activation for authenticated users
    if (!isRental && !isSeller) {
      activateCapability('seller').catch(() => {});
    } else if (isRental && !isLandlord) {
      activateCapability('landlord').catch(() => {});
    }
  }, [user, loading, isRental, isSeller, isLandlord, activateCapability, router]);

  if (loading || !user) {
    return null;
  }

  const pageTitle = isRental 
    ? (t as any).demand?.landlord_title || 'List Property for Rent'
    : (t as any).demand?.sell_title || 'List Product for Sale';

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-slate-900">
      <div className="container mx-auto px-4 py-8 max-w-3xl">
        <button
          onClick={() => router.back()}
          className="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-300 hover:text-slate-900 dark:hover:text-white mb-6 transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          {(t as any).ui?.back || 'Back'}
        </button>

        <div className="bg-white dark:bg-slate-800 rounded-2xl shadow-lg p-6 md:p-8">
          <h1 className="text-2xl md:text-3xl font-bold text-slate-900 dark:text-white mb-6">
            {pageTitle}
          </h1>
          
          {isRental && (
            <div className="mb-6 p-4 bg-teal-50 dark:bg-teal-900/20 border border-teal-200 dark:border-teal-800 rounded-xl">
              <p className="text-sm text-teal-800 dark:text-teal-200">
                <strong>Rental Listing:</strong> You're creating a property listing for rent. 
                Fill in the rental price, duration, and availability details.
              </p>
            </div>
          )}

          <ProductForm 
            listingType={listingType as 'sale' | 'rent'}
            onSuccess={() => {
              router.push('/seller/products');
            }}
          />
        </div>
      </div>
    </div>
  );
}

export default function NewProductPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-600"></div>
      </div>
    }>
      <NewProductContent />
    </Suspense>
  );
}
