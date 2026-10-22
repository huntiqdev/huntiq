'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import Link from 'next/link';
import {
  ShoppingCart, Crown, Star, ArrowUpRight, Package,
  CheckCircle, Clock, Sparkles, Lock, Images, PhoneCall,
  ShoppingBag, Zap,
} from 'lucide-react';

const TIER_COLOR: Record<string, string> = {
  normal:  'text-zinc-400',
  gold:    'text-yellow-500',
  premium: 'text-purple-500',
};

const TIER_GRADIENT: Record<string, string> = {
  normal:  'from-zinc-500 to-zinc-600',
  gold:    'from-yellow-400 via-amber-400 to-orange-400',
  premium: 'from-purple-500 via-violet-500 to-indigo-500',
};

// What each tier unlocks for buyers
const TIER_PERKS: Record<string, { icon: any; text: string; color: string }[]> = {
  normal: [
    { icon: Star,      text: 'Gold: access exclusive gold-tier products',    color: 'text-yellow-400' },
    { icon: Crown,     text: 'Premium: access all premium-tier products',    color: 'text-purple-400' },
    { icon: Images,    text: 'Premium: priority support & early access',     color: 'text-purple-400' },
  ],
  gold: [
    { icon: CheckCircle, text: 'Gold products unlocked',                     color: 'text-yellow-400' },
    { icon: Crown,       text: 'Premium: access all premium-tier products',  color: 'text-purple-400' },
    { icon: Sparkles,    text: 'Premium: priority support & early access',   color: 'text-purple-400' },
  ],
  premium: [
    { icon: CheckCircle, text: 'All products unlocked',                      color: 'text-purple-400' },
    { icon: Sparkles,    text: 'Priority support & early access',            color: 'text-purple-400' },
    { icon: Zap,         text: 'Maximum tier — all features active',         color: 'text-purple-400' },
  ],
};

export default function BuyerDashboardPage() {
  const { user, buyerTier } = useAuth();
  const [profile, setProfile] = useState<any>(null);
  const [orders, setOrders]   = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!user) return;
    Promise.all([
      api.getBuyerProfile(user.token),
      api.getBuyerOrders(user.token),
    ]).then(([p, o]) => {
      setProfile(p);
      setOrders(o ?? []);
    }).catch(() => {}).finally(() => setLoading(false));
  }, [user]);

  if (loading) return <div className="flex items-center justify-center h-64"><Spinner /></div>;

  const pending   = orders.filter(o => o.status === 'pending').length;
  const accepted  = orders.filter(o => o.status === 'accepted').length;
  const completed = orders.filter(o => o.status === 'completed').length;
  const perks     = TIER_PERKS[buyerTier] ?? TIER_PERKS.normal;

  return (
    <div className="space-y-5 max-w-2xl">
      {/* Header */}
      <div>
        <h1 className="text-xl font-black text-zinc-900 dark:text-white">
          Welcome{profile?.name ? `, ${profile.name}` : ''}
        </h1>
        <div className="flex items-center gap-1.5 mt-0.5">
          {buyerTier === 'gold'    && <Star  size={11} className="text-yellow-500" />}
          {buyerTier === 'premium' && <Crown size={11} className="text-purple-500" />}
          <p className={`text-xs font-bold capitalize ${TIER_COLOR[buyerTier]}`}>{buyerTier} buyer</p>
        </div>
      </div>

      {/* KPIs */}
      <div className="grid grid-cols-3 gap-3">
        {[
          { label: 'Pending',   value: pending,   color: 'text-blue-500',    bg: 'bg-blue-100 dark:bg-blue-900/30' },
          { label: 'Accepted',  value: accepted,  color: 'text-green-500',   bg: 'bg-green-100 dark:bg-green-900/30' },
          { label: 'Completed', value: completed, color: 'text-emerald-500', bg: 'bg-emerald-100 dark:bg-emerald-900/30' },
        ].map(({ label, value, color, bg }) => (
          <div key={label} className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4">
            <div className={`w-8 h-8 rounded-xl ${bg} flex items-center justify-center mb-2`}>
              {label === 'Pending'   && <Clock       size={14} className={color} />}
              {label === 'Accepted'  && <CheckCircle size={14} className={color} />}
              {label === 'Completed' && <Sparkles    size={14} className={color} />}
            </div>
            <p className={`text-2xl font-black ${color}`}>{value}</p>
            <p className="text-xs text-zinc-500 mt-0.5">{label}</p>
          </div>
        ))}
      </div>

      {/* Tier card */}
      <div className="bg-[#0d1117] border border-white/8 rounded-2xl overflow-hidden">
        <div className={`h-1 bg-gradient-to-r ${TIER_GRADIENT[buyerTier]}`} />
        <div className="p-4 space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2.5">
              <div className={`w-9 h-9 rounded-xl bg-gradient-to-br ${TIER_GRADIENT[buyerTier]} flex items-center justify-center shadow-lg`}>
                {buyerTier === 'premium' ? <Crown size={16} className="text-white" /> : <Star size={16} className="text-white" />}
              </div>
              <div>
                <p className="text-[10px] font-black uppercase tracking-widest text-white/30">Buyer Tier</p>
                <p className="text-sm font-black text-white capitalize">{buyerTier}</p>
              </div>
            </div>
            {buyerTier !== 'premium' && (
              <Link href="/buyer/upgrade"
                className="flex items-center gap-1.5 bg-orange-600 hover:bg-orange-700 text-white px-3 py-1.5 rounded-xl text-xs font-bold transition-colors">
                Upgrade <ArrowUpRight size={11} />
              </Link>
            )}
          </div>

          <div className="grid grid-cols-1 gap-1.5">
            {perks.map((p, i) => (
              <div key={i} className="flex items-center gap-2 bg-white/[0.03] border border-white/5 rounded-xl px-3 py-2">
                {buyerTier === 'premium' || (buyerTier === 'gold' && p.color === 'text-yellow-400')
                  ? <CheckCircle size={11} className={p.color} />
                  : <Lock size={9} className="text-white/20 shrink-0" />
                }
                <p.icon size={11} className={`${p.color} shrink-0`} />
                <span className="text-[10px] text-white/50">{p.text}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Quick links */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <Link href="/buyer/orders"
          className="flex items-center gap-3 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4 hover:shadow-md transition-all">
          <div className="w-9 h-9 rounded-xl bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center text-blue-600">
            <ShoppingCart size={16} />
          </div>
          <div className="flex-1">
            <p className="font-bold text-zinc-900 dark:text-white text-sm">My Requests</p>
            <p className="text-xs text-zinc-400">{orders.length} total order{orders.length !== 1 ? 's' : ''}</p>
          </div>
          <ArrowUpRight size={14} className="text-zinc-300" />
        </Link>
        <Link href="/products"
          className="flex items-center gap-3 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4 hover:shadow-md transition-all">
          <div className="w-9 h-9 rounded-xl bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center text-orange-600">
            <Package size={16} />
          </div>
          <div className="flex-1">
            <p className="font-bold text-zinc-900 dark:text-white text-sm">Browse Products</p>
            <p className="text-xs text-zinc-400">Find tools &amp; equipment</p>
          </div>
          <ArrowUpRight size={14} className="text-zinc-300" />
        </Link>
      </div>

      {/* Recent orders — with product title and empty state */}
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl overflow-hidden">
        <div className="flex items-center justify-between px-5 py-4 border-b border-zinc-100 dark:border-zinc-800">
          <h2 className="font-bold text-zinc-900 dark:text-white text-sm">Recent Requests</h2>
          {orders.length > 0 && (
            <Link href="/buyer/orders" className="text-xs text-orange-600 font-bold hover:underline">View all →</Link>
          )}
        </div>

        {orders.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 gap-3 text-zinc-400">
            <ShoppingBag size={28} className="opacity-30" />
            <p className="text-sm font-semibold">No order requests yet</p>
            <Link href="/products"
              className="text-xs text-orange-600 font-bold hover:underline flex items-center gap-1">
              Browse products <ArrowUpRight size={11} />
            </Link>
          </div>
        ) : (
          <div className="divide-y divide-zinc-100 dark:divide-zinc-800">
            {orders.slice(0, 5).map((o: any) => (
              <div key={o.id} className="flex items-center gap-3 px-5 py-3">
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-semibold text-zinc-800 dark:text-zinc-200 truncate">
                    {o.product_title || `Product #${o.product_id}`}
                  </p>
                  <p className="text-xs text-zinc-400">
                    {o.product_price ? `ETB ${Number(o.product_price).toLocaleString()} · ` : ''}
                    {new Date(o.created_at).toLocaleDateString()}
                  </p>
                </div>
                <span className={`text-[10px] font-black px-2 py-0.5 rounded-full uppercase shrink-0 ${
                  o.status === 'pending'   ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-600' :
                  o.status === 'accepted'  ? 'bg-green-100 dark:bg-green-900/30 text-green-600' :
                  o.status === 'completed' ? 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600' :
                  o.status === 'rejected'  ? 'bg-red-100 dark:bg-red-900/30 text-red-500' :
                  'bg-zinc-100 dark:bg-zinc-800 text-zinc-500'
                }`}>{o.status}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
