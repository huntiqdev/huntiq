'use client';
import { useEffect, useState, FormEvent } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import {
  User, Phone, Mail, Building2, Star, CheckCircle, AlertCircle,
  Lock, Eye, EyeOff, TrendingUp, Plus, X, Info,
  ShoppingBag, Home, ArrowUpRight,
} from 'lucide-react';
import Link from 'next/link';
import { SellerTier, getTierThresholds, getTierFeeParams } from '@/lib/constants/fees';
import { ReferralCard } from '@/components/referral/ReferralCard';
import { CapabilityActivationModal } from '@/components/auth/CapabilityActivationModal';

export default function SellerProfilePage() {
  const { user, isSeller, isLandlord } = useAuth();
  const [profile, setProfile] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');
  const [additionalPhones, setAdditionalPhones] = useState<string[]>([]);
  const [activationCap, setActivationCap] = useState<'seller' | 'landlord' | null>(null);

  const [pwForm, setPwForm] = useState({ current: '', next: '', confirm: '' });
  const [showPw, setShowPw] = useState({ current: false, next: false, confirm: false });
  const [pwSaving, setPwSaving] = useState(false);
  const [pwError, setPwError] = useState('');
  const [pwSuccess, setPwSuccess] = useState(false);
  const [settings, setSettings] = useState<any>(null);

  useEffect(() => {
    if (!user) return;
    api.getSellerProfile(user.token).then(p => {
      setProfile(p);
      setAdditionalPhones(p?.additional_phones ?? []);
    }).finally(() => setLoading(false));

    api.getPublicPaymentSettings().then(setSettings).catch(() => {});
  }, [user]);

  const tier = (profile?.tier ?? 'normal') as SellerTier;
  const thresholds = getTierThresholds(settings);
  const normalParams = getTierFeeParams('normal', settings);
  const goldParams = getTierFeeParams('gold', settings);
  const premiumParams = getTierFeeParams('premium', settings);
  
  const TIER_INFO = {
    normal:  { label: 'Normal',  color: 'text-zinc-600',   bg: 'bg-zinc-100 dark:bg-zinc-800',        bar: 'bg-zinc-400',   desc: `${normalParams.startupPct}% startup · +${normalParams.weeklyIncrementPct}%/week`, maxPhones: 0 },
    gold:    { label: 'Gold',    color: 'text-yellow-600', bg: 'bg-yellow-100 dark:bg-yellow-900/30', bar: 'bg-yellow-500', desc: `${goldParams.startupPct}% startup · +${goldParams.weeklyIncrementPct}%/week`, maxPhones: 3 },
    premium: { label: 'Premium', color: 'text-purple-600', bg: 'bg-purple-100 dark:bg-purple-900/30', bar: 'bg-purple-500', desc: `${premiumParams.startupPct}% startup · +${premiumParams.weeklyIncrementPct}%/week`,  maxPhones: -1 },
  };

  const tierCfg = TIER_INFO[tier] ?? TIER_INFO.normal;
  const maxPhones = tierCfg.maxPhones; // 0 = none, -1 = unlimited
  const canAddPhone = maxPhones === -1 || additionalPhones.length < maxPhones;

  const handleProfileSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!user) return;
    setSaving(true); setError(''); setSuccess('');
    try {
      await api.updateSellerProfile({
        name: profile.name,
        business_name: profile.business_name,
        email: profile.email,
        additional_phones: additionalPhones.filter(p => p.trim()),
      }, user.token);
      setSuccess('Profile updated successfully!');
    } catch (err: any) { setError(err.message); }
    finally { setSaving(false); }
  };

  const handlePasswordSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setPwError('');
    if (pwForm.next.length < 8) { setPwError('New password must be at least 8 characters'); return; }
    if (pwForm.next !== pwForm.confirm) { setPwError('Passwords do not match'); return; }
    if (!user) return;
    setPwSaving(true);
    try {
      await api.sellerChangePassword({ current_password: pwForm.current, new_password: pwForm.next }, user.token);
      setPwSuccess(true);
      setPwForm({ current: '', next: '', confirm: '' });
    } catch (err: any) { setPwError(String(err?.message ?? '').replace(/<[^>]*>/g, '')); }
    finally { setPwSaving(false); }
  };

  const cls = 'w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 transition-all';
  const lbl = 'text-xs font-bold text-zinc-500 dark:text-zinc-400 block mb-1';

  const PW_LABELS: Record<'current' | 'next' | 'confirm', string> = {
    current: 'Current Password',
    next:    'New Password',
    confirm: 'Confirm New Password',
  };

  if (loading) return <div className="flex items-center justify-center h-64"><Spinner /></div>;
  if (!profile) return <p className="text-zinc-400">Profile not found.</p>;

  return (
    <div className="space-y-5 max-w-lg">
      <div>
        <h1 className="text-xl font-black text-zinc-900 dark:text-white">Profile</h1>
        <p className="text-xs text-zinc-500 mt-0.5">Manage your seller account</p>
      </div>

      {/* Profile card */}
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl overflow-hidden">
        <div className="bg-gradient-to-r from-orange-600 to-orange-500 px-5 py-6">
          <div className="flex items-center gap-4">
            <div className="w-14 h-14 rounded-2xl bg-white/20 flex items-center justify-center text-white font-black text-xl">
              {profile.name?.[0]?.toUpperCase() ?? 'S'}
            </div>
            <div>
              <p className="font-black text-white text-lg">{profile.name}</p>
              {profile.business_name && <p className="text-white/70 text-sm">{profile.business_name}</p>}
              <div className="flex items-center gap-1 mt-1">
                <Star size={12} className="text-yellow-300 fill-yellow-300" />
                <span className="text-white/80 text-xs font-bold">{Number(profile.rating).toFixed(1)} rating</span>
              </div>
            </div>
          </div>
        </div>

        <form onSubmit={handleProfileSubmit} className="p-5 space-y-4">
          <div>
            <label className={lbl}><User size={11} className="inline mr-1" />Full Name</label>
            <input value={profile.name} onChange={e => setProfile((p: any) => ({ ...p, name: e.target.value }))}
              className={cls} placeholder="Your name" />
          </div>
          <div>
            <label className={lbl}><Building2 size={11} className="inline mr-1" />Business Name</label>
            <input value={profile.business_name ?? ''} onChange={e => setProfile((p: any) => ({ ...p, business_name: e.target.value }))}
              className={cls} placeholder="Optional" />
          </div>
          <div>
            <label className={lbl}><Mail size={11} className="inline mr-1" />Email</label>
            <input type="email" value={profile.email} onChange={e => setProfile((p: any) => ({ ...p, email: e.target.value }))}
              className={cls} />
          </div>

          {/* Primary phone — always read-only */}
          <div>
            <label className={lbl}><Phone size={11} className="inline mr-1" />Primary Phone (login, cannot change)</label>
            <input value={profile.phone} disabled className={cls + ' opacity-50 cursor-not-allowed'} />
          </div>

          {/* Additional phones — tier-gated */}
          <div>
            <div className="flex items-center justify-between mb-1">
              <label className={lbl + ' mb-0'}>
                <Phone size={11} className="inline mr-1" />
                Additional Contact Phones
              </label>
              <span className={`text-[10px] font-black px-2 py-0.5 rounded-full ${tierCfg.bg} ${tierCfg.color}`}>
                {tier === 'normal' ? 'Normal — upgrade to add' : tier === 'gold' ? 'Gold — up to 3' : 'Premium — unlimited'}
              </span>
            </div>

            {maxPhones === 0 ? (
              <div className="flex items-start gap-2 bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2.5">
                <Info size={13} className="text-zinc-400 shrink-0 mt-0.5" />
                <p className="text-xs text-zinc-500">
                  Upgrade to <span className="font-bold text-yellow-600">Gold</span> (ETB {(thresholds.gold/1000).toFixed(0)}k sold) to add up to 3 contact phones,
                  or <span className="font-bold text-purple-600">Premium</span> (ETB {(thresholds.premium/1000000).toFixed(1)}M sold) for unlimited.
                </p>
              </div>
            ) : (
              <div className="space-y-2">
                {additionalPhones.map((ph, i) => (
                  <div key={i} className="flex gap-2">
                    <input
                      type="tel"
                      value={ph}
                      onChange={e => {
                        const next = [...additionalPhones];
                        next[i] = e.target.value;
                        setAdditionalPhones(next);
                      }}
                      placeholder="+251 9X XXX XXXX"
                      className="flex-1 border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 transition-all"
                    />
                    <button type="button" onClick={() => setAdditionalPhones(additionalPhones.filter((_, j) => j !== i))}
                      className="p-2.5 rounded-xl text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 hover:text-red-600 transition-colors">
                      <X size={14} />
                    </button>
                  </div>
                ))}
                {canAddPhone && (
                  <button type="button" onClick={() => setAdditionalPhones([...additionalPhones, ''])}
                    className="flex items-center gap-1.5 text-xs font-bold text-orange-600 hover:text-orange-700 transition-colors">
                    <Plus size={13} /> Add phone number
                    {maxPhones > 0 && (
                      <span className="text-zinc-400 font-normal">({additionalPhones.length}/{maxPhones})</span>
                    )}
                  </button>
                )}
              </div>
            )}
          </div>

          {error && (
            <div className="flex items-center gap-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl px-3 py-2">
              <AlertCircle size={13} className="text-red-500 shrink-0" />
              <p className="text-xs text-red-600 dark:text-red-400">{error}</p>
            </div>
          )}
          {success && (
            <div className="flex items-center gap-2 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-xl px-3 py-2">
              <CheckCircle size={13} className="text-green-600 shrink-0" />
              <p className="text-xs text-green-700 dark:text-green-400">{success}</p>
            </div>
          )}

          <button type="submit" disabled={saving}
            className="w-full bg-orange-600 hover:bg-orange-700 text-white py-3 rounded-xl font-bold text-sm transition-colors disabled:opacity-50">
            {saving ? 'Saving…' : 'Save Changes'}
          </button>
        </form>
      </div>

      {/* Tier & upgrade */}
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
        <div className="flex items-center gap-2 mb-4 pb-3 border-b border-zinc-100 dark:border-zinc-800">
          <TrendingUp size={15} className="text-orange-600" />
          <h2 className="font-bold text-zinc-900 dark:text-white">Seller Tier</h2>
        </div>
        {(() => {
          const sold = profile?.total_value_sold ?? 0;
          const nextTier = tier === 'normal' ? 'gold' : tier === 'gold' ? 'premium' : null;
          const nextThreshold = nextTier ? thresholds[nextTier as keyof typeof thresholds] : 0;
          const progress = nextTier ? Math.min(100, (sold / nextThreshold) * 100) : 100;
          return (
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className={`flex items-center gap-1.5 text-sm font-black px-3 py-1 rounded-full ${tierCfg.bg} ${tierCfg.color}`}>
                  <Star size={12} /> {tierCfg.label} Tier
                </span>
                <span className="text-xs text-zinc-500">{tierCfg.desc}</span>
              </div>
              <div className="text-sm space-y-1">
                <div className="flex justify-between">
                  <span className="opacity-50">Total Value Sold</span>
                  <span className="font-black">ETB {Number(sold).toLocaleString()}</span>
                </div>
                {nextTier && (
                  <div className="flex justify-between">
                    <span className="opacity-50">Next tier ({nextTier})</span>
                    <span className="font-bold">ETB {nextThreshold.toLocaleString()}</span>
                  </div>
                )}
              </div>
              {nextTier && (
                <div>
                  <div className="h-2 bg-zinc-100 dark:bg-zinc-800 rounded-full overflow-hidden">
                    <div className={`h-full rounded-full transition-all ${tierCfg.bar}`} style={{ width: `${progress}%` }} />
                  </div>
                  <p className="text-[10px] text-zinc-400 mt-1">
                    ETB {Math.max(0, nextThreshold - sold).toLocaleString()} more to reach {nextTier.charAt(0).toUpperCase() + nextTier.slice(1)}
                  </p>
                </div>
              )}
              {tier === 'premium' && (
                <p className="text-xs text-purple-600 dark:text-purple-400 font-bold">
                  ✦ You are at the highest tier — enjoy the lowest fees and unlimited contact phones!
                </p>
              )}
              {/* Phone perks summary */}
              <div className="grid grid-cols-3 gap-2 pt-2 border-t border-zinc-100 dark:border-zinc-800 text-center text-[10px]">
                {(['normal', 'gold', 'premium'] as const).map(t => (
                  <div key={t} className={`rounded-xl p-2 ${t === tier ? TIER_INFO[t].bg : 'bg-zinc-50 dark:bg-zinc-800/50 opacity-50'}`}>
                    <p className={`font-black ${t === tier ? TIER_INFO[t].color : 'text-zinc-500'}`}>{TIER_INFO[t].label}</p>
                    <p className="text-zinc-500 mt-0.5">
                      {t === 'normal' ? '1 phone' : t === 'gold' ? `1+${TIER_INFO[t].maxPhones} phones` : 'unlimited phones'}
                    </p>
                  </div>
                ))}
              </div>
            </div>
          );
        })()}
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
            <Link
              href="/buyer/dashboard"
              className="text-xs font-bold text-orange-600 dark:text-orange-400 hover:underline flex items-center gap-1"
            >
              <span>Portal</span>
              <ArrowUpRight size={12} />
            </Link>
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
            <span className="text-[10px] font-bold text-green-500 bg-green-500/10 px-2 py-0.5 rounded-full">
              Active
            </span>
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

      {/* 1-Click Activation Modal */}
      {activationCap && (
        <CapabilityActivationModal
          isOpen={!!activationCap}
          capability={activationCap}
          onClose={() => setActivationCap(null)}
          onSuccess={() => {
            setActivationCap(null);
            setSuccess('Capability activated successfully!');
          }}
        />
      )}

      {/* Referral */}
      <ReferralCard role="seller" />

      {/* Change password */}
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">        <div className="flex items-center gap-2 mb-4 pb-3 border-b border-zinc-100 dark:border-zinc-800">
          <Lock size={15} className="text-orange-600" />
          <h2 className="font-bold text-zinc-900 dark:text-white">Change Password</h2>
        </div>
        {pwSuccess && (
          <div className="flex items-center gap-2 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-xl px-3 py-2 mb-4">
            <CheckCircle size={13} className="text-green-600 shrink-0" />
            <p className="text-xs text-green-700 dark:text-green-400">Password updated!</p>
          </div>
        )}
        <form onSubmit={handlePasswordSubmit} className="space-y-3">
          {(['current', 'next', 'confirm'] as const).map(field => (
            <div key={field}>
              <label className={lbl}>{PW_LABELS[field]}</label>
              <div className="relative">
                <input required type={showPw[field] ? 'text' : 'password'}
                  value={pwForm[field]} onChange={e => setPwForm(f => ({ ...f, [field]: e.target.value }))}
                  className={cls + ' pr-10'} placeholder="••••••••" />
                <button type="button" onClick={() => setShowPw(s => ({ ...s, [field]: !s[field] }))}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-400 hover:text-zinc-600 transition-colors">
                  {showPw[field] ? <EyeOff size={14} /> : <Eye size={14} />}
                </button>
              </div>
            </div>
          ))}
          {pwError && <p className="text-xs text-red-500 bg-red-50 dark:bg-red-900/20 px-3 py-2 rounded-xl">{pwError}</p>}
          <button type="submit" disabled={pwSaving}
            className="w-full bg-zinc-800 dark:bg-zinc-700 hover:bg-zinc-700 dark:hover:bg-zinc-600 text-white py-2.5 rounded-xl font-bold text-sm transition-colors disabled:opacity-50">
            {pwSaving ? 'Updating…' : 'Update Password'}
          </button>
        </form>
      </div>
    </div>
  );
}
