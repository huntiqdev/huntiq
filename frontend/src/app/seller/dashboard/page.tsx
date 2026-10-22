'use client';
import { useEffect, useState, useCallback } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import Link from 'next/link';
import {
  Package, Clock, CheckCircle, TrendingUp, Plus, ArrowUpRight,
  AlertCircle, ShoppingCart, Star, AlertTriangle, DollarSign,
  Crown, Lock, Images, PhoneCall, BarChart2, Sparkles,
  Gift, Users, Copy, CheckCircle2, Wallet, ChevronRight,
} from 'lucide-react';
import { DEFAULT_TIER_THRESHOLDS } from '@/lib/constants/fees';
import { TierUpgradeModal } from '@/components/seller/TierUpgradeModal';

const TIER_CONFIG = {
  normal:  { label: 'Normal',  gradient: 'from-zinc-500 to-zinc-600',           bar: 'bg-zinc-400',   glow: '' },
  gold:    { label: 'Gold',    gradient: 'from-yellow-400 via-amber-400 to-orange-400', bar: 'bg-yellow-500', glow: 'shadow-yellow-500/20' },
  premium: { label: 'Premium', gradient: 'from-purple-500 via-violet-500 to-indigo-500', bar: 'bg-purple-500', glow: 'shadow-purple-500/20' },
} as const;

// Features locked per tier — shown as upgrade hints
const LOCKED_FEATURES = {
  normal: [
    { icon: ShoppingCart, text: 'Order Requests',    unlockAt: 'Gold',    color: 'text-yellow-400' },
    { icon: DollarSign,   text: 'Sales & Transactions', unlockAt: 'Gold', color: 'text-yellow-400' },
    { icon: Images,       text: 'Up to 7 images',    unlockAt: 'Gold',    color: 'text-yellow-400' },
    { icon: PhoneCall,    text: '3 extra phones',     unlockAt: 'Gold',    color: 'text-yellow-400' },
    { icon: BarChart2,    text: 'Analytics',          unlockAt: 'Premium', color: 'text-purple-400' },
    { icon: Sparkles,     text: 'Unlimited images',   unlockAt: 'Premium', color: 'text-purple-400' },
  ],
  gold: [
    { icon: BarChart2,  text: 'Analytics dashboard', unlockAt: 'Premium', color: 'text-purple-400' },
    { icon: Images,     text: 'Unlimited images',    unlockAt: 'Premium', color: 'text-purple-400' },
    { icon: PhoneCall,  text: 'Unlimited phones',    unlockAt: 'Premium', color: 'text-purple-400' },
    { icon: Sparkles,   text: 'Top listing priority + dedicated support', unlockAt: 'Premium', color: 'text-purple-400' },
  ],
  premium: [],
} as const;

function KPI({ icon: Icon, label, value, color, href, alert }: {
  icon: any; label: string; value: string | number; color: string; href?: string; alert?: boolean;
}) {
  const inner = (
    <div className={`bg-white dark:bg-zinc-900 border rounded-2xl p-4 flex items-start gap-3 hover:shadow-md transition-all ${
      alert ? 'border-red-200 dark:border-red-800' : 'border-zinc-200 dark:border-zinc-800'
    } ${href ? 'cursor-pointer' : ''}`}>
      <div className={`w-9 h-9 rounded-xl flex items-center justify-center shrink-0 ${color}`}>
        <Icon size={16} />
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-xl font-black text-zinc-900 dark:text-white">{value}</p>
        <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-0.5">{label}</p>
      </div>
      {href && <ArrowUpRight size={13} className="ml-auto text-zinc-300 dark:text-zinc-600 shrink-0" />}
    </div>
  );
  return href ? <Link href={href}>{inner}</Link> : inner;
}

function TierCard({ stats, upgradePrices, tierThresholds, onUpgrade }: {
  stats: any;
  upgradePrices: { gold: number; premium: number };
  tierThresholds: { gold: number; premium: number };
  onUpgrade: (target: 'gold' | 'premium') => void;
}) {
  const tier = (stats?.tier ?? 'normal') as keyof typeof TIER_CONFIG;
  const cfg  = TIER_CONFIG[tier] ?? TIER_CONFIG.normal;
  const sold = stats?.total_value_sold ?? 0;
  const locked = LOCKED_FEATURES[tier] ?? [];

  let nextTier = '', nextThreshold = 0, progress = 100;
  if (tier === 'normal') {
    nextTier = 'Gold'; nextThreshold = tierThresholds.gold;
    progress = Math.min(100, (sold / nextThreshold) * 100);
  } else if (tier === 'gold') {
    nextTier = 'Premium'; nextThreshold = tierThresholds.premium;
    progress = Math.min(100, (sold / nextThreshold) * 100);
  }

  return (
    <div className="bg-[#0d1117] border border-white/8 rounded-2xl overflow-hidden">
      {/* Gradient top bar */}
      <div className={`h-1 bg-gradient-to-r ${cfg.gradient}`} />

      <div className="p-4 space-y-4">
        {/* Tier badge + label */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2.5">
            <div className={`w-9 h-9 rounded-xl bg-gradient-to-br ${cfg.gradient} flex items-center justify-center shadow-lg ${cfg.glow}`}>
              {tier === 'premium' ? <Crown size={16} className="text-white" /> : <Star size={16} className="text-white" />}
            </div>
            <div>
              <p className="text-[10px] font-black uppercase tracking-widest text-white/30">Seller Tier</p>
              <p className="text-sm font-black text-white">{cfg.label}</p>
            </div>
          </div>
          <div className="text-right">
            <p className="text-[10px] text-white/30">Value Sold</p>
            <p className="text-sm font-black text-white">ETB {sold.toLocaleString()}</p>
          </div>
        </div>

        {/* Progress bar */}
        {nextTier && (
          <div>
            <div className="flex justify-between text-[10px] text-white/30 mb-1.5">
              <span>{cfg.label}</span>
              <span>{nextTier} at ETB {nextThreshold.toLocaleString()}</span>
            </div>
            <div className="h-2 bg-white/5 rounded-full overflow-hidden">
              <div className={`h-full rounded-full bg-gradient-to-r ${cfg.gradient} transition-all duration-700`}
                style={{ width: `${Math.max(2, progress)}%` }} />
            </div>
            <p className="text-[10px] text-white/20 mt-1">
              ETB {Math.max(0, nextThreshold - sold).toLocaleString()} more to reach {nextTier} automatically
            </p>
          </div>
        )}

        {/* Locked features */}
        {locked.length > 0 && (
          <div className="space-y-1.5">
            <p className="text-[10px] font-black uppercase tracking-widest text-white/25">Locked features</p>
            <div className="grid grid-cols-2 gap-1.5">
              {locked.map((f, i) => (
                <div key={i} className="flex items-center gap-2 bg-white/[0.03] border border-white/5 rounded-xl px-2.5 py-2">
                  <Lock size={9} className="text-white/20 shrink-0" />
                  <f.icon size={11} className={`${f.color} shrink-0`} />
                  <span className="text-[10px] text-white/40 truncate">{f.text}</span>
                  <span className={`text-[8px] font-black ml-auto shrink-0 ${f.color} opacity-70`}>{f.unlockAt}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Premium max tier message */}
        {tier === 'premium' && (
          <div className="flex items-center gap-2 bg-purple-500/10 border border-purple-500/20 rounded-xl px-3 py-2.5">
            <Sparkles size={13} className="text-purple-400 shrink-0" />
            <p className="text-xs text-purple-300 font-semibold">Maximum tier — all features unlocked</p>
          </div>
        )}

        {/* Upgrade buttons — normal sees both Gold & Premium, gold sees only Premium */}
        {tier !== 'premium' && (
          <div className={`grid gap-2 ${tier === 'normal' ? 'grid-cols-2' : 'grid-cols-1'}`}>
            {tier === 'normal' && (
              <button onClick={() => onUpgrade('gold')}
                className="flex flex-col items-center gap-1 py-3 px-2 rounded-xl bg-gradient-to-br from-yellow-400 via-amber-400 to-orange-400 text-white font-black shadow-lg shadow-yellow-500/20 hover:opacity-90 active:scale-[0.98] transition-all">
                <div className="flex items-center gap-1.5">
                  <Star size={13} />
                  <span className="text-xs">Gold</span>
                </div>
                <span className="text-[10px] font-bold opacity-80">ETB {upgradePrices.gold.toLocaleString()}</span>
              </button>
            )}
            <button onClick={() => onUpgrade('premium')}
              className={`flex flex-col items-center gap-1 py-3 px-2 rounded-xl bg-gradient-to-br from-purple-500 via-violet-500 to-indigo-500 text-white font-black shadow-lg shadow-purple-500/20 hover:opacity-90 active:scale-[0.98] transition-all ${tier === 'gold' ? 'col-span-1' : ''}`}>
              <div className="flex items-center gap-1.5">
                <Crown size={13} />
                <span className="text-xs">Premium</span>
              </div>
              <span className="text-[10px] font-bold opacity-80">ETB {upgradePrices.premium.toLocaleString()}</span>
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

// ── Referral Summary Card ─────────────────────────────────────────────────────
function ReferralSummary() {
  const { user } = useAuth();
  const [stats, setStats]     = useState<any>(null);
  const [balance, setBalance] = useState<any>(null);
  const [copied, setCopied]   = useState(false);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    if (!user) return;
    try {
      const [s, b] = await Promise.all([
        api.getReferralStats(user.token),
        api.getReferralBalance(user.token),
      ]);
      setStats(s); setBalance(b);
    } catch { /* non-critical */ }
    finally { setLoading(false); }
  }, [user]);

  useEffect(() => { load(); }, [load]);

  const handleCopy = async () => {
    if (!stats?.referral_link) return;
    try { await navigator.clipboard.writeText(stats.referral_link); setCopied(true); setTimeout(() => setCopied(false), 2000); }
    catch { /* ignore */ }
  };

  if (loading) return (
    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4 animate-pulse h-32" />
  );
  if (!stats) return null;

  const available = Math.floor(balance?.available_balance ?? 0);
  const pending   = Math.floor(balance?.pending_amount ?? 0);
  const earned    = Math.floor(balance?.total_earned ?? 0);

  return (
    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl overflow-hidden">
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-zinc-100 dark:border-zinc-800">
        <div className="flex items-center gap-2">
          <Gift size={14} className="text-orange-500" />
          <p className="text-xs font-black uppercase tracking-widest text-zinc-500 dark:text-zinc-400">Referral Program</p>
        </div>
        <Link href="/seller/profile" className="text-[10px] font-bold text-orange-600 hover:underline flex items-center gap-0.5">
          Full details <ChevronRight size={10} />
        </Link>
      </div>

      <div className="p-4 space-y-3">
        {/* Stats row: invited / earned / available */}
        <div className="grid grid-cols-3 gap-2">
          <div className="bg-zinc-50 dark:bg-zinc-800 rounded-xl p-2.5 text-center">
            <Users size={12} className="text-orange-500 mx-auto mb-1" />
            <p className="text-lg font-black text-zinc-900 dark:text-white">{stats.invited_count ?? 0}</p>
            <p className="text-[9px] font-bold text-zinc-400 uppercase tracking-wide">Invited</p>
          </div>
          <div className="bg-zinc-50 dark:bg-zinc-800 rounded-xl p-2.5 text-center">
            <TrendingUp size={12} className="text-green-500 mx-auto mb-1" />
            <p className="text-sm font-black text-zinc-900 dark:text-white">ETB {earned.toLocaleString()}</p>
            <p className="text-[9px] font-bold text-zinc-400 uppercase tracking-wide">Earned</p>
          </div>
          <div className="bg-zinc-50 dark:bg-zinc-800 rounded-xl p-2.5 text-center">
            <Wallet size={12} className={`mx-auto mb-1 ${available >= 500 ? 'text-green-500' : 'text-zinc-400'}`} />
            <p className={`text-sm font-black ${available >= 500 ? 'text-green-600' : 'text-zinc-500 dark:text-zinc-400'}`}>
              ETB {available.toLocaleString()}
            </p>
            <p className="text-[9px] font-bold text-zinc-400 uppercase tracking-wide">Available</p>
          </div>
        </div>

        {/* Pending withdrawal */}
        {pending > 0 && (
          <div className="flex items-center gap-1.5 text-[10px] text-yellow-600 dark:text-yellow-400 bg-yellow-50 dark:bg-yellow-900/20 rounded-xl px-3 py-1.5">
            <Clock size={10} /> ETB {pending.toLocaleString()} pending withdrawal
          </div>
        )}

        {/* Suspended warning */}
        {stats.referral_suspended && (
          <div className="flex items-center gap-1.5 text-[10px] text-red-500 bg-red-50 dark:bg-red-900/20 rounded-xl px-3 py-1.5">
            <AlertCircle size={10} /> Your referral code is suspended — contact support
          </div>
        )}

        {/* Invite link copy */}
        {!stats.referral_suspended && (
          <div className="flex gap-2">
            <div className="flex-1 bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-[10px] text-zinc-500 dark:text-zinc-400 truncate font-mono">
              {stats.referral_link}
            </div>
            <button onClick={handleCopy}
              className="flex items-center gap-1.5 px-3 py-2 bg-orange-600 hover:bg-orange-700 text-white rounded-xl text-[10px] font-bold transition-colors shrink-0">
              {copied ? <CheckCircle2 size={11} /> : <Copy size={11} />}
              {copied ? 'Copied!' : 'Copy'}
            </button>
          </div>
        )}

        {/* Who invited you */}
        {stats.invitation_accepted_at && (
          <p className="text-[10px] text-zinc-400">
            Invited by <span className="font-bold text-zinc-600 dark:text-zinc-300">{stats.inviter_name ?? stats.inviter_code}</span>
            {' '}· they earn 2% on your paid actions
          </p>
        )}
      </div>
    </div>
  );
}

function daysRemaining(expiresAt: string | null): number | null {
  if (!expiresAt) return null;
  const diff = new Date(expiresAt).getTime() - Date.now();
  return Math.max(0, Math.ceil(diff / 86_400_000));
}

export default function SellerDashboardPage() {
  const { user, sellerPreference } = useAuth();
  const [stats, setStats] = useState<any>(null);
  const [products, setProducts] = useState<any[]>([]);
  const [orders, setOrders] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [showUpgrade, setShowUpgrade] = useState(false);
  const [upgradeTo, setUpgradeTo]     = useState<'gold' | 'premium'>('gold');
  const [upgradePrices, setUpgradePrices] = useState({ gold: 500, premium: 2000 });
  const [tierThresholds, setTierThresholds] = useState<{ gold: number; premium: number }>({ gold: DEFAULT_TIER_THRESHOLDS.gold, premium: DEFAULT_TIER_THRESHOLDS.premium });

  // Context-aware labels based on seller_preference
  const labelMap = {
    sale: {
      dashboard: 'Seller Dashboard',
      products: 'My Products',
      addButton: '+ Add Product',
      revenue: 'Sales Revenue',
    },
    rent: {
      dashboard: 'Landlord Dashboard',
      products: 'My Properties',
      addButton: '+ Add Property',
      revenue: 'Rental Income',
    },
    both: {
      dashboard: 'Seller/Landlord Dashboard',
      products: 'My Listings',
      addButton: '+ Add Listing',
      revenue: 'Total Revenue',
    },
  };
  const labels = labelMap[sellerPreference as keyof typeof labelMap] ?? labelMap.sale;

  useEffect(() => {
    if (!user) return;
    Promise.allSettled([
      api.getSellerStats(user.token),
      api.getSellerProductsDetailed(user.token),
      api.getPublicPaymentSettings(),
    ]).then(([statsRes, productsRes, settingsRes]) => {
      const statsData = statsRes.status === 'fulfilled' ? (statsRes.value ?? {}) : {};
      setStats(statsData);
      setProducts(productsRes.status === 'fulfilled' ? (productsRes.value ?? []) : []);
      if (settingsRes.status === 'fulfilled' && settingsRes.value) {
        const s = settingsRes.value;
        setUpgradePrices({
          gold:    Number(s.gold_upgrade_price)    || 500,
          premium: Number(s.premium_upgrade_price) || 2000,
        });
        setTierThresholds({
          gold:    Number(s.gold_threshold)    || DEFAULT_TIER_THRESHOLDS.gold,
          premium: Number(s.premium_threshold) || DEFAULT_TIER_THRESHOLDS.premium,
        });
      }
      // Only fetch orders if seller is gold or premium (endpoint requires gold tier)
      const tier = statsData?.tier ?? 'normal';
      if (tier === 'gold' || tier === 'premium') {
        api.getSellerOrders(user.token)
          .then(o => setOrders(o ?? []))
          .catch(() => setOrders([]));
      }
    }).finally(() => setLoading(false));
  }, [user]);

  if (loading) return <div className="flex items-center justify-center h-64"><Spinner /></div>;

  const pendingOrders = orders.filter(o => o.status === 'pending');
  const unpaidFees = stats?.unpaid_fees ?? 0;

  return (
    <div className="space-y-5 max-w-4xl">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-black text-zinc-900 dark:text-white">{labels.dashboard}</h1>
          <p className="text-xs text-zinc-500 mt-0.5">
            {new Date().toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric' })}
          </p>
        </div>
        <Link href="/seller/products"
          className="flex items-center gap-2 bg-orange-600 hover:bg-orange-700 text-white px-4 py-2 rounded-xl text-sm font-bold transition-colors shadow-md shadow-orange-600/20">
          <Plus size={15} /> {labels.addButton}
        </Link>
      </div>

      {/* Alerts */}
      {pendingOrders.length > 0 && (
        <div className="flex items-center gap-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-2xl px-4 py-3">
          <ShoppingCart size={16} className="text-blue-500 shrink-0" />
          <p className="text-sm font-semibold text-blue-700 dark:text-blue-300">
            {pendingOrders.length} new order request{pendingOrders.length > 1 ? 's' : ''} waiting for your response.
          </p>
          <Link href="/seller/orders" className="ml-auto text-xs font-bold text-blue-600 hover:underline shrink-0">Review →</Link>
        </div>
      )}
      {unpaidFees > 0 && (
        <div className="flex items-center gap-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-2xl px-4 py-3">
          <AlertTriangle size={16} className="text-red-500 shrink-0" />
          <div className="flex-1">
            <p className="text-sm font-semibold text-red-700 dark:text-red-300">
              Outstanding listing fees: <strong>ETB {unpaidFees.toLocaleString()}</strong>
            </p>
            <p className="text-xs text-red-500/70 mt-0.5">Unpaid listings may be delisted. Pay to keep them active.</p>
          </div>
          <Link href="/seller/transactions" className="ml-auto text-xs font-bold text-red-600 hover:underline shrink-0">Pay now →</Link>
        </div>
      )}
      {(stats?.rejected ?? 0) > 0 && (
        <div className="flex items-center gap-3 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-2xl px-4 py-3">
          <AlertCircle size={16} className="text-amber-500 shrink-0" />
          <p className="text-sm font-semibold text-amber-700 dark:text-amber-300">
            {stats.rejected} product{stats.rejected > 1 ? 's were' : ' was'} rejected. Check admin notes and resubmit.
          </p>
          <Link href="/seller/products" className="ml-auto text-xs font-bold text-amber-600 hover:underline shrink-0">View →</Link>
        </div>
      )}

      {/* KPIs */}
      <div className="grid grid-cols-2 lg:grid-cols-5 gap-3">
        <KPI icon={TrendingUp}   label="Net Profit"       value={`ETB ${(stats?.total_profit ?? 0).toLocaleString()}`}
          color="bg-orange-100 dark:bg-orange-900/30 text-orange-600" />
        <KPI icon={Package}      label="Total Listings"   value={stats?.total_products ?? 0}
          color="bg-blue-100 dark:bg-blue-900/30 text-blue-600"    href="/seller/products" />
        <KPI icon={Clock}        label="Pending Review"   value={stats?.pending ?? 0}
          color="bg-yellow-100 dark:bg-yellow-900/30 text-yellow-600" href="/seller/products" />
        <KPI icon={CheckCircle}  label="Active Listings"  value={stats?.approved ?? 0}
          color="bg-green-100 dark:bg-green-900/30 text-green-600" />
        <KPI icon={ShoppingCart} label="Order Requests"   value={pendingOrders.length}
          color="bg-blue-100 dark:bg-blue-900/30 text-blue-500"    href="/seller/orders"
          alert={pendingOrders.length > 0} />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Tier card */}
        <TierCard stats={stats} upgradePrices={upgradePrices} tierThresholds={tierThresholds} onUpgrade={(target) => { setUpgradeTo(target); setShowUpgrade(true); }} />

        {/* Fee summary */}
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4">
          <div className="flex items-center gap-2 mb-3">
            <DollarSign size={14} className="text-orange-600" />
            <p className="text-xs font-black uppercase tracking-widest opacity-60">Fee Summary</p>
          </div>
          <div className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="opacity-50">Total Revenue (Gross)</span>
              <span className="font-black text-green-600">ETB {(stats?.total_revenue ?? 0).toLocaleString()}</span>
            </div>
            <div className="flex justify-between">
              <span className="opacity-50">Total Fees Paid</span>
              <span className="font-black text-red-400">ETB {(stats?.paid_fees ?? 0).toLocaleString()}</span>
            </div>
            <div className="border-t border-white/5 pt-2 flex justify-between">
              <span className="font-bold">Net Profit</span>
              <span className="font-black text-orange-500">ETB {(stats?.total_profit ?? 0).toLocaleString()}</span>
            </div>
            {(stats?.pending_revenue ?? 0) > 0 && (
              <div className="flex justify-between text-xs pt-1">
                <span className="opacity-50">Pending (in transit)</span>
                <span className="font-black text-yellow-500">ETB {(stats.pending_revenue).toLocaleString()}</span>
              </div>
            )}
            <div className="flex justify-between text-xs">
              <span className="opacity-50">Unpaid Listing Fees</span>
              <span className={`font-black ${unpaidFees > 0 ? 'text-red-500' : 'text-zinc-400'}`}>
                ETB {unpaidFees.toLocaleString()}
              </span>
            </div>
          </div>
          <Link href="/seller/transactions"
            className="mt-3 flex items-center gap-1 text-xs font-bold text-orange-600 hover:underline">
            <TrendingUp size={11} /> View full payment history
          </Link>
        </div>
      </div>

      {/* Referral summary — full width below tier/fee row */}
      <ReferralSummary />

      {/* Active listings with time remaining — or empty state if no products */}
      {products.length === 0 ? (
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-8 flex flex-col items-center gap-3 text-zinc-400">
          <Package size={28} className="opacity-30" />
          <p className="text-sm font-semibold">No listings yet</p>
          <Link href="/seller/products"
            className="flex items-center gap-1.5 bg-orange-600 hover:bg-orange-700 text-white px-4 py-2 rounded-xl text-xs font-bold transition-colors">
            <Plus size={13} /> Add your first product
          </Link>
        </div>
      ) : products.filter(p => p.status === 'approved' && p.expires_at).length > 0 ? (
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl overflow-hidden">
          <div className="flex items-center justify-between px-5 py-4 border-b border-zinc-100 dark:border-zinc-800">
            <h2 className="font-bold text-zinc-900 dark:text-white">Active Listings — Time Remaining</h2>
            <Link href="/seller/products" className="text-xs text-orange-600 font-bold hover:underline">View all →</Link>
          </div>
          <div className="divide-y divide-zinc-100 dark:divide-zinc-800">
            {products.filter(p => p.status === 'approved').slice(0, 5).map(p => {
              const days = daysRemaining(p.expires_at);
              const urgent = days !== null && days <= 3;
              return (
                <div key={p.id} className="flex items-center gap-3 px-5 py-3">
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-semibold text-zinc-800 dark:text-zinc-200 truncate">{p.title}</p>
                    <p className="text-xs text-zinc-400">ETB {Number(p.price).toLocaleString()}</p>
                  </div>
                  <div className="text-right shrink-0">
                    {days !== null ? (
                      <p className={`text-xs font-black ${urgent ? 'text-red-500' : 'text-green-600'}`}>
                        {days === 0 ? 'Expires today' : `${days}d left`}
                      </p>
                    ) : (
                      <p className="text-xs text-zinc-400">No expiry</p>
                    )}
                    {p.total_fee_amount && (
                      <p className={`text-[10px] ${p.fee_paid ? 'text-green-500' : 'text-red-400'}`}>
                        Fee: ETB {Number(p.total_fee_amount).toLocaleString()} {p.fee_paid ? '✓' : '(unpaid)'}
                      </p>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      ) : null}

      {/* Recent order requests */}
      {orders.length > 0 && (
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl overflow-hidden">
          <div className="flex items-center justify-between px-5 py-4 border-b border-zinc-100 dark:border-zinc-800">
            <h2 className="font-bold text-zinc-900 dark:text-white">Recent Order Requests</h2>
            <Link href="/seller/orders" className="text-xs text-orange-600 font-bold hover:underline">View all →</Link>
          </div>
          <div className="divide-y divide-zinc-100 dark:divide-zinc-800">
            {orders.slice(0, 4).map(o => (
              <div key={o.id} className="flex items-center gap-3 px-5 py-3">
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-semibold text-zinc-800 dark:text-zinc-200 truncate">{o.buyer_name}</p>
                  <p className="text-xs text-zinc-400 truncate">{o.product_title}</p>
                </div>
                <span className={`text-[10px] font-black px-2 py-0.5 rounded-full uppercase shrink-0 ${
                  o.status === 'pending'   ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-600' :
                  o.status === 'accepted'  ? 'bg-green-100 dark:bg-green-900/30 text-green-600' :
                  o.status === 'completed' ? 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600' :
                  'bg-zinc-100 dark:bg-zinc-800 text-zinc-500'
                }`}>{o.status}</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {showUpgrade && (
        <TierUpgradeModal
          currentTier={stats?.tier ?? 'normal'}
          initialTier={upgradeTo}
          onClose={() => setShowUpgrade(false)}
          onSuccess={() => {
            setShowUpgrade(false);
            if (user) api.getSellerStats(user.token).then(s => setStats(s ?? {}));
          }}
        />
      )}
    </div>
  );
}
