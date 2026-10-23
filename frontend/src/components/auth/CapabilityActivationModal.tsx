'use client';

import { useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { ShoppingBag, Home, CheckCircle2, AlertCircle, X, Sparkles, ArrowRight } from 'lucide-react';

interface CapabilityActivationModalProps {
  isOpen: boolean;
  onClose: () => void;
  capability: 'seller' | 'landlord';
  onSuccess?: () => void;
}

export function CapabilityActivationModal({
  isOpen,
  onClose,
  capability,
  onSuccess,
}: CapabilityActivationModalProps) {
  const { user, activateCapability } = useAuth();
  const [businessName, setBusinessName] = useState('');
  const [preference, setPreference] = useState<'sale' | 'rent' | 'both'>(
    capability === 'landlord' ? 'rent' : 'sale'
  );
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  if (!isOpen) return null;

  const isLandlord = capability === 'landlord';
  const title = isLandlord ? 'Activate Landlord Capability' : 'Activate Seller Capability';
  const subtitle = isLandlord
    ? 'Start listing houses, apartments, shops, and land for rent on SellBuy.'
    : 'Start listing and selling products, tools, and equipment across Ethiopia.';

  const handleActivate = async () => {
    setLoading(true);
    setError('');
    try {
      const res = await activateCapability(
        capability,
        businessName.trim() || undefined,
        preference
      );
      if (!res.success) {
        setError(res.error || 'Failed to activate capability. Please try again.');
        setLoading(false);
        return;
      }
      onSuccess?.();
      onClose();
    } catch (err: any) {
      setError(err?.message || 'An unexpected error occurred.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
      <div className="relative w-full max-w-md bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl shadow-2xl overflow-hidden p-6 text-zinc-900 dark:text-white">
        {/* Close Button */}
        <button
          onClick={onClose}
          disabled={loading}
          className="absolute top-5 right-5 p-2 rounded-xl text-zinc-400 hover:text-zinc-700 dark:hover:text-white hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
        >
          <X size={18} />
        </button>

        {/* Header Icon */}
        <div className="flex items-center gap-3.5 mb-5">
          <div
            className={`w-12 h-12 rounded-2xl flex items-center justify-center shadow-lg ${
              isLandlord
                ? 'bg-teal-600 shadow-teal-600/30 text-white'
                : 'bg-orange-600 shadow-orange-600/30 text-white'
            }`}
          >
            {isLandlord ? <Home size={24} /> : <ShoppingBag size={24} />}
          </div>
          <div>
            <h2 className="text-lg font-black tracking-tight">{title}</h2>
            <p className="text-xs text-zinc-500 dark:text-zinc-400">{subtitle}</p>
          </div>
        </div>

        {/* Benefits List */}
        <div className="bg-zinc-50 dark:bg-white/[0.03] border border-zinc-200/80 dark:border-white/5 rounded-2xl p-4 mb-5 space-y-2.5">
          <p className="text-[10px] font-black uppercase tracking-wider text-zinc-400">
            Included with your unified identity:
          </p>
          <div className="flex items-center gap-2 text-xs text-zinc-700 dark:text-zinc-300">
            <CheckCircle2 size={14} className="text-green-500 shrink-0" />
            <span>No separate account or new login password needed</span>
          </div>
          <div className="flex items-center gap-2 text-xs text-zinc-700 dark:text-zinc-300">
            <CheckCircle2 size={14} className="text-green-500 shrink-0" />
            <span>Instant access to listing tools and seller studio</span>
          </div>
          <div className="flex items-center gap-2 text-xs text-zinc-700 dark:text-zinc-300">
            <CheckCircle2 size={14} className="text-green-500 shrink-0" />
            <span>Retain all your existing purchases and favorites</span>
          </div>
        </div>

        {/* Optional Business Name */}
        <div className="space-y-4 mb-6">
          <div>
            <label className="text-[11px] font-bold uppercase tracking-wider text-zinc-500 dark:text-zinc-400 block mb-1.5">
              {isLandlord ? 'Agency / Property Name' : 'Shop / Business Name'}{' '}
              <span className="normal-case font-normal text-zinc-400">(optional)</span>
            </label>
            <input
              type="text"
              value={businessName}
              onChange={(e) => setBusinessName(e.target.value)}
              placeholder={isLandlord ? 'e.g. Addis Real Estate' : 'e.g. Bikila Hardware'}
              className="w-full bg-zinc-100 dark:bg-white/5 border border-zinc-200 dark:border-white/10 rounded-xl px-4 py-3 text-sm outline-none focus:ring-2 focus:ring-orange-500/30 focus:border-orange-500/60 transition-all text-zinc-900 dark:text-white placeholder-zinc-400"
            />
          </div>
        </div>

        {error && (
          <div className="flex items-center gap-2 bg-red-500/10 border border-red-500/20 text-red-500 dark:text-red-400 text-xs rounded-xl p-3 mb-4">
            <AlertCircle size={14} className="shrink-0" />
            <span>{error}</span>
          </div>
        )}

        {/* Action Buttons */}
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={onClose}
            disabled={loading}
            className="flex-1 py-3 rounded-xl border border-zinc-200 dark:border-zinc-800 text-xs font-bold text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={handleActivate}
            disabled={loading}
            className={`flex-1 py-3 rounded-xl text-white text-xs font-black flex items-center justify-center gap-2 shadow-lg transition-all active:scale-95 disabled:opacity-50 ${
              isLandlord
                ? 'bg-teal-600 hover:bg-teal-500 shadow-teal-600/30'
                : 'bg-orange-600 hover:bg-orange-500 shadow-orange-600/30'
            }`}
          >
            {loading ? (
              'Activating...'
            ) : (
              <>
                <Sparkles size={14} />
                <span>Activate Capability</span>
                <ArrowRight size={13} />
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
