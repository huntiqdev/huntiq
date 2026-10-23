'use client';
import { useEffect, useState, FormEvent } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import {
  Star, Crown, Phone, User, Mail, ArrowUpRight,
  Lock, Eye, EyeOff, Check, AlertCircle, Calendar,
  ShoppingBag, Home, Sparkles,
} from 'lucide-react';
import Link from 'next/link';
import { ReferralCard } from '@/components/referral/ReferralCard';
import { CapabilityActivationModal } from '@/components/auth/CapabilityActivationModal';

const TIER_CONFIG = {
  normal:  { label: 'Normal',  icon: null,  color: 'text-zinc-500',   badge: 'bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400' },
  gold:    { label: 'Gold',    icon: Star,  color: 'text-yellow-500', badge: 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-400' },
  premium: { label: 'Premium', icon: Crown, color: 'text-purple-500', badge: 'bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-400' },
} as const;

export default function BuyerProfilePage() {
  const { user, buyerTier, isSeller, isLandlord } = useAuth();
  const [profile, setProfile] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [activationCap, setActivationCap] = useState<'seller' | 'landlord' | null>(null);

  // Password change state
  const [pwForm, setPwForm] = useState({ current: '', next: '', confirm: '' });
  const [showPw, setShowPw] = useState(false);
  const [pwLoading, setPwLoading] = useState(false);
  const [pwError, setPwError] = useState('');
  const [pwSuccess, setPwSuccess] = useState('');

  useEffect(() => {
    if (!user) return;
    api.getBuyerProfile(user.token)
      .then(setProfile)
      .finally(() => setLoading(false));
  }, [user]);

  const handlePasswordChange = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setPwError(''); setPwSuccess('');
    if (!user) return;
    if (pwForm.next.length < 8) { setPwError('New password must be at least 8 characters'); return; }
    if (pwForm.next !== pwForm.confirm) { setPwError('Passwords do not match'); return; }
    setPwLoading(true);
    try {
      await api.buyerChangePassword(
        { current_password: pwForm.current, new_password: pwForm.next },
        user.token,
      );
      setPwSuccess('Password updated successfully');
      setPwForm({ current: '', next: '', confirm: '' });
    } catch (err: any) {
      setPwError(err.message);
    } finally { setPwLoading(false); }
  };

  if (loading) return <div className="flex items-center justify-center h-64"><Spinner /></div>;

  const tier = (buyerTier || 'normal') as keyof typeof TIER_CONFIG;
  const cfg = TIER_CONFIG[tier];
  const TierIcon = cfg.icon;

  const cls = 'w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 transition-all';

  return (
    <div className="space-y-5 max-w-md">
      <h1 className="text-xl font-black">My Profile</h1>

      {/* Profile card */}
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5 space-y-4">
        {/* Avatar + tier */}
        <div className="flex items-center gap-4">
          <div className="w-14 h-14 rounded-2xl bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center text-orange-600 font-black text-xl">
            {profile?.name?.[0]?.toUpperCase() ?? 'B'}
          </div>
          <div>
            <p className="font-black text-lg text-zinc-900 dark:text-white">{profile?.name ?? '—'}</p>
            <div className={`inline-flex items-center gap-1 text-xs font-black px-2 py-0.5 rounded-full mt-1 ${cfg.badge}`}>
              {TierIcon && <TierIcon size={10} />}
              {cfg.label} Buyer
            </div>
          </div>
        </div>

        {/* Details */}
        <div className="space-y-3 pt-2 border-t border-zinc-100 dark:border-zinc-800">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-lg bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center shrink-0">
              <Phone size={13} className="text-zinc-500" />
            </div>
            <div>
              <p className="text-[10px] font-black uppercase tracking-widest text-zinc-400">Phone (auto-call)</p>
              <p className="text-sm font-semibold text-zinc-700 dark:text-zinc-300">{profile?.phone ?? '—'}</p>
            </div>
          </div>

          {profile?.email && !profile.email.includes('@noemail.local') && (
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-lg bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center shrink-0">
                <Mail size={13} className="text-zinc-500" />
              </div>
              <div>
                <p className="text-[10px] font-black uppercase tracking-widest text-zinc-400">Email</p>
                <p className="text-sm font-semibold text-zinc-700 dark:text-zinc-300">{profile.email}</p>
              </div>
            </div>
          )}

          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-lg bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center shrink-0">
              <Calendar size={13} className="text-zinc-500" />
            </div>
            <div>
              <p className="text-[10px] font-black uppercase tracking-widest text-zinc-400">Member since</p>
              <p className="text-sm font-semibold text-zinc-700 dark:text-zinc-300">
                {profile?.created_at
                  ? new Date(profile.created_at).toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })
                  : '—'}
              </p>
            </div>
          </div>

          {/* Tier paid date */}
          {profile?.buyer_tier_paid_at && tier !== 'normal' && (
            <div className="flex items-center gap-3">
              <div className={`w-8 h-8 rounded-lg flex items-center justify-center shrink-0 ${
                tier === 'gold' ? 'bg-yellow-100 dark:bg-yellow-900/30' : 'bg-purple-100 dark:bg-purple-900/30'
              }`}>
                {TierIcon && <TierIcon size={13} className={cfg.color} />}
              </div>
              <div>
                <p className="text-[10px] font-black uppercase tracking-widest text-zinc-400">
                  {cfg.label} since
                </p>
                <p className="text-sm font-semibold text-zinc-700 dark:text-zinc-300">
                  {new Date(profile.buyer_tier_paid_at).toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })}
                </p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Marketplace Capabilities */}
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5 space-y-3">
        <div className="flex items-center justify-between">
          <p className="text-xs font-black uppercase tracking-wider text-zinc-400">
            Marketplace Capabilities
          </p>
          <span className="text-[10px] text-zinc-500 font-medium">Unified Account</span>
        </div>

        <div className="space-y-2">
          {/* Buyer */}
          <div className="flex items-center justify-between p-3 rounded-xl bg-zinc-50 dark:bg-zinc-800/50 border border-zinc-100 dark:border-zinc-800">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-lg bg-orange-100 dark:bg-orange-950/40 text-orange-600 flex items-center justify-center">
                <User size={15} />
              </div>
              <div>
                <p className="text-xs font-bold text-zinc-800 dark:text-zinc-200">Buyer Activity</p>
                <p className="text-[10px] text-zinc-400">Discover and purchase products across Ethiopia</p>
              </div>
            </div>
            <span className="text-[10px] font-bold text-green-500 bg-green-500/10 px-2 py-0.5 rounded-full">
              Active
            </span>
          </div>

          {/* Seller */}
          <div className="flex items-center justify-between p-3 rounded-xl bg-zinc-50 dark:bg-zinc-800/50 border border-zinc-100 dark:border-zinc-800">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-lg bg-orange-100 dark:bg-orange-950/40 text-orange-600 flex items-center justify-center">
                <ShoppingBag size={15} />
              </div>
              <div>
                <p className="text-xs font-bold text-zinc-800 dark:text-zinc-200">Seller Studio</p>
                <p className="text-[10px] text-zinc-400">List items, manage stock, and receive orders</p>
              </div>
            </div>
            {isSeller ? (
              <Link
                href="/seller/dashboard"
                className="text-xs font-bold text-orange-600 dark:text-orange-400 hover:underline flex items-center gap-1"
              >
                <span>Dashboard</span>
                <ArrowUpRight size={12} />
              </Link>
            ) : (
              <button
                type="button"
                onClick={() => setActivationCap('seller')}
                className="text-[11px] font-bold bg-orange-600 hover:bg-orange-500 text-white px-2.5 py-1 rounded-lg transition-all"
              >
                Activate
              </button>
            )}
          </div>

          {/* Landlord */}
          <div className="flex items-center justify-between p-3 rounded-xl bg-zinc-50 dark:bg-zinc-800/50 border border-zinc-100 dark:border-zinc-800">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-lg bg-teal-100 dark:bg-teal-950/40 text-teal-600 flex items-center justify-center">
                <Home size={15} />
              </div>
              <div>
                <p className="text-xs font-bold text-zinc-800 dark:text-zinc-200">Landlord Console</p>
                <p className="text-[10px] text-zinc-400">List properties and rentable assets</p>
              </div>
            </div>
            {isLandlord ? (
              <Link
                href="/seller/products?preference=rent"
                className="text-xs font-bold text-teal-600 dark:text-teal-400 hover:underline flex items-center gap-1"
              >
                <span>Console</span>
                <ArrowUpRight size={12} />
              </Link>
            ) : (
              <button
                type="button"
                onClick={() => setActivationCap('landlord')}
                className="text-[11px] font-bold bg-teal-600 hover:bg-teal-500 text-white px-2.5 py-1 rounded-lg transition-all"
              >
                Activate
              </button>
            )}
          </div>
        </div>
      </div>

      {/* Referral */}
      <ReferralCard role="buyer" />

      {/* 1-Click Activation Modal */}
      {activationCap && (
        <CapabilityActivationModal
          isOpen={!!activationCap}
          capability={activationCap}
          onClose={() => setActivationCap(null)}
          onSuccess={() => {
            // Updated session automatically populates via AuthContext
          }}
        />
      )}

      {/* Upgrade CTA */}
      {tier !== 'premium' && (        <Link href="/buyer/upgrade"
          className="flex items-center justify-between bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4 hover:shadow-md transition-all group">
          <div className="flex items-center gap-3">
            <div className={`w-9 h-9 rounded-xl flex items-center justify-center ${
              tier === 'gold' ? 'bg-purple-100 dark:bg-purple-900/30' : 'bg-yellow-100 dark:bg-yellow-900/30'
            }`}>
              {tier === 'gold'
                ? <Crown size={16} className="text-purple-500" />
                : <Star size={16} className="text-yellow-500" />}
            </div>
            <div>
              <p className="font-bold text-sm">{tier === 'gold' ? 'Upgrade to Premium' : 'Upgrade to Gold or Premium'}</p>
              <p className="text-xs text-zinc-400">
                {tier === 'gold' ? 'Unlock all seller images' : 'Unlock gold & premium seller images'}
              </p>
            </div>
          </div>
          <ArrowUpRight size={15} className="text-zinc-400 group-hover:text-orange-600 transition-colors" />
        </Link>
      )}

      {/* Change password */}
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
        <div className="flex items-center gap-2 mb-4">
          <Lock size={15} className="text-zinc-500" />
          <p className="font-bold text-sm">Change Password</p>
        </div>
        <form onSubmit={handlePasswordChange} className="space-y-3">
          <div>
            <label className="text-xs font-bold text-zinc-500 block mb-1">Current Password</label>
            <div className="relative">
              <input
                required
                type={showPw ? 'text' : 'password'}
                value={pwForm.current}
                onChange={e => setPwForm(p => ({ ...p, current: e.target.value }))}
                placeholder="••••••••"
                className={cls + ' pr-10'}
              />
              <button type="button" onClick={() => setShowPw(v => !v)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-400 hover:text-zinc-600">
                {showPw ? <EyeOff size={14} /> : <Eye size={14} />}
              </button>
            </div>
          </div>
          <div>
            <label className="text-xs font-bold text-zinc-500 block mb-1">New Password</label>
            <input
              required
              type={showPw ? 'text' : 'password'}
              value={pwForm.next}
              onChange={e => setPwForm(p => ({ ...p, next: e.target.value }))}
              placeholder="Min 8 characters"
              className={cls}
            />
          </div>
          <div>
            <label className="text-xs font-bold text-zinc-500 block mb-1">Confirm New Password</label>
            <input
              required
              type={showPw ? 'text' : 'password'}
              value={pwForm.confirm}
              onChange={e => setPwForm(p => ({ ...p, confirm: e.target.value }))}
              placeholder="Repeat new password"
              className={cls}
            />
          </div>

          {pwError && (
            <div className="flex items-center gap-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl px-3 py-2">
              <AlertCircle size={13} className="text-red-500 shrink-0" />
              <p className="text-xs text-red-600 dark:text-red-400">{pwError}</p>
            </div>
          )}
          {pwSuccess && (
            <div className="flex items-center gap-2 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-xl px-3 py-2">
              <Check size={13} className="text-green-500 shrink-0" />
              <p className="text-xs text-green-600 dark:text-green-400">{pwSuccess}</p>
            </div>
          )}

          <button
            type="submit"
            disabled={pwLoading}
            className="w-full py-2.5 rounded-xl bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 font-bold text-sm hover:opacity-90 transition-all disabled:opacity-50"
          >
            {pwLoading ? 'Updating…' : 'Update Password'}
          </button>
        </form>
      </div>
    </div>
  );
}
