'use client';
import { useState, useEffect, useCallback, FormEvent } from 'react';
import { api } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import type { ReferralStats, ReferralBalance, CommissionEntry, WithdrawalRequest } from '@/types';
import {
  Link2, Copy, CheckCircle2, Users, TrendingUp, UserPlus,
  AlertCircle, ChevronDown, ChevronUp, Gift, Loader2,
  Wallet, Clock, DollarSign,
} from 'lucide-react';

interface Props { role?: 'seller' | 'buyer'; }

const TX_TYPE_LABELS: Record<string, string> = {
  product_purchase: 'Product Sale', gold_upgrade: 'Gold Upgrade',
  premium_upgrade: 'Premium Upgrade', buyer_gold_upgrade: 'Buyer Gold Upgrade',
  buyer_premium_upgrade: 'Buyer Premium Upgrade',
  seller_registration: 'Seller Registration', seller_renewal: 'Seller Renewal',
};
const STATUS_COLOR: Record<string, string> = {
  pending:  'text-yellow-600 bg-yellow-100 dark:bg-yellow-900/30',
  approved: 'text-blue-600 bg-blue-100 dark:bg-blue-900/30',
  rejected: 'text-red-600 bg-red-100 dark:bg-red-900/30',
  paid:     'text-green-600 bg-green-100 dark:bg-green-900/30',
};
const COMMISSION_STATUS_COLOR: Record<string, string> = {
  pending:   'text-yellow-600 bg-yellow-100 dark:bg-yellow-900/30',
  confirmed: 'text-green-600 bg-green-100 dark:bg-green-900/30',
  reversed:  'text-red-600 bg-red-100 dark:bg-red-900/30',
};
const PM_LABEL: Record<string, string> = {
  telebirr: 'Telebirr', safaricom: 'Safaricom', cbe: 'CBE', bank: 'Bank',
};

export function ReferralCard({ role = 'seller' }: Props) {
  const { user } = useAuth();
  const [stats, setStats] = useState<ReferralStats | null>(null);
  const [balance, setBalance] = useState<ReferralBalance | null>(null);
  const [commissions, setCommissions] = useState<CommissionEntry[]>([]);
  const [commissionsLoaded, setCommissionsLoaded] = useState(false);
  const [withdrawals, setWithdrawals] = useState<WithdrawalRequest[]>([]);
  const [withdrawalsLoaded, setWithdrawalsLoaded] = useState(false);
  const [loading, setLoading] = useState(true);
  const [copied, setCopied] = useState(false);
  const [showCommissions, setShowCommissions] = useState(false);
  const [showWithdrawals, setShowWithdrawals] = useState(false);
  const [showWithdrawForm, setShowWithdrawForm] = useState(false);

  // Accept invitation form
  const [showAcceptForm, setShowAcceptForm] = useState(false);
  const [inviteCode, setInviteCode] = useState('');
  const [acceptLoading, setAcceptLoading] = useState(false);
  const [acceptError, setAcceptError] = useState('');
  const [acceptSuccess, setAcceptSuccess] = useState('');

  // Withdrawal form
  const [wAmount, setWAmount] = useState('');
  const [wMethod, setWMethod] = useState('telebirr');
  const [wPhone, setWPhone] = useState('');
  const [wLoading, setWLoading] = useState(false);
  const [wError, setWError] = useState('');
  const [wSuccess, setWSuccess] = useState('');
  const [minWithdrawal, setMinWithdrawal] = useState(500);
  const [maxWithdrawal, setMaxWithdrawal] = useState(1000000);
  const [commissionRate, setCommissionRate] = useState(1.2);

  const fetchStats = useCallback(async () => {
    if (!user) return;
    try {
      const [statsData, balData, settings] = await Promise.all([
        role === 'buyer' ? api.getBuyerReferralStats(user.token) : api.getReferralStats(user.token),
        role === 'buyer' ? api.getBuyerReferralBalance(user.token) : api.getReferralBalance(user.token),
        api.getPublicPaymentSettings(),
      ]);
      setStats(statsData);
      setBalance(balData);
      if (settings?.min_referral_withdrawal) {
        setMinWithdrawal(Number(settings.min_referral_withdrawal) || 500);
      }
      if (settings?.max_referral_withdrawal) {
        setMaxWithdrawal(Number(settings.max_referral_withdrawal) || 1000000);
      }
      if (settings?.referral_commission_rate) {
        setCommissionRate(Number(settings.referral_commission_rate) || 1.2);
      }
    } catch { /* non-critical */ }
    finally { setLoading(false); }
  }, [user, role]);

  useEffect(() => { fetchStats(); }, [fetchStats]);

  const loadCommissions = async () => {
    if (!user || commissionsLoaded) return;
    try {
      const data = role === 'buyer'
        ? await api.getBuyerReferralCommissions(user.token)
        : await api.getReferralCommissions(user.token);
      setCommissions(data?.commissions ?? []);
      setCommissionsLoaded(true);
    } catch { /* non-critical */ }
  };

  const loadWithdrawals = async () => {
    if (!user || withdrawalsLoaded) return;
    try {
      const data = role === 'buyer'
        ? await api.getMyBuyerReferralWithdrawals(user.token)
        : await api.getMyReferralWithdrawals(user.token);
      setWithdrawals(Array.isArray(data) ? data : []);
      setWithdrawalsLoaded(true);
    } catch { /* non-critical */ }
  };

  const handleCopy = async () => {
    if (!stats?.referral_link) return;
    try { await navigator.clipboard.writeText(stats.referral_link); setCopied(true); setTimeout(() => setCopied(false), 2000); } catch { /* ignore */ }
  };

  const handleAccept = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!user || !inviteCode.trim()) return;
    setAcceptLoading(true); setAcceptError('');
    try {
      const fn = role === 'buyer' ? api.acceptBuyerReferralInvitation : api.acceptReferralInvitation;
      await fn(inviteCode.trim().toUpperCase(), user.token);
      setAcceptSuccess(`Invitation accepted! Your inviter will earn ${commissionRate}% when you make a paid action.`);
      setShowAcceptForm(false); setInviteCode('');
      await fetchStats();
    } catch (err: any) {
      setAcceptError(err?.message ?? 'Failed to accept invitation — check the code and try again');
    } finally { setAcceptLoading(false); }
  };

  const handleWithdraw = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!user) return;
    const amount = Math.floor(parseFloat(wAmount)); // integer ETB only
    if (isNaN(amount) || amount < minWithdrawal) { setWError(`Minimum withdrawal is ETB ${minWithdrawal}`); return; }
    if (amount > maxWithdrawal) { setWError(`Maximum withdrawal per request is ETB ${maxWithdrawal.toLocaleString()}`); return; }
    if (balance && amount > Math.floor(balance.available_balance)) {
      setWError(`Amount exceeds available balance (ETB ${Math.floor(balance.available_balance).toLocaleString()})`);
      return;
    }
    if (wPhone.trim().length < 7) { setWError('Account/phone number is required (min 7 characters)'); return; }
    setWLoading(true); setWError(''); setWSuccess('');
    try {
      const fn = role === 'buyer' ? api.requestBuyerReferralWithdrawal : api.requestReferralWithdrawal;
      await fn({ amount, payment_method: wMethod, payer_phone: wPhone.trim() }, user.token);
      setWSuccess('Withdrawal request submitted — admin will process within 24 hours.');
      setShowWithdrawForm(false); setWAmount(''); setWPhone('');
      setWithdrawalsLoaded(false); // force reload
      await fetchStats();
    } catch (err: any) {
      setWError(err?.message ?? 'Failed to submit withdrawal request');
    } finally { setWLoading(false); }
  };

  // ── Loading skeleton ──────────────────────────────────────────────────────
  if (loading) {
    return (
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5 animate-pulse">
        <div className="h-4 bg-zinc-200 dark:bg-zinc-700 rounded w-1/3 mb-3" />
        <div className="h-10 bg-zinc-100 dark:bg-zinc-800 rounded-xl mb-3" />
        <div className="h-8 bg-zinc-100 dark:bg-zinc-800 rounded-xl w-2/3" />
      </div>
    );
  }

  // ── Error / unavailable state ─────────────────────────────────────────────
  if (!stats) {
    return (
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
        <div className="flex items-center gap-2 text-zinc-400">
          <Gift size={15} />
          <p className="text-sm">Referral program unavailable — try refreshing.</p>
        </div>
      </div>
    );
  }

  const hasInviter = !!stats.invitation_accepted_at;
  const hasEarnings = stats.total_commission > 0 || stats.invited_count > 0;
  const isSuspended = stats.referral_suspended;

  return (
    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl overflow-hidden">
      {/* Header */}
      <div className="bg-gradient-to-r from-orange-600 to-orange-500 px-5 py-4 flex items-center gap-3">
        <div className="w-9 h-9 rounded-xl bg-white/20 flex items-center justify-center shrink-0">
          <Gift size={18} className="text-white" />
        </div>
        <div>
          <p className="font-black text-white text-sm">Referral Program</p>
          <p className="text-white/70 text-xs">Earn {commissionRate}% when your invitees make a paid action</p>
        </div>
      </div>

      <div className="p-5 space-y-4">

        {/* Stats + balance row */}
        <div className="grid grid-cols-2 gap-3">
          <div className="bg-zinc-50 dark:bg-zinc-800 rounded-xl p-3 text-center">
            <div className="flex items-center justify-center gap-1 mb-1">
              <Users size={12} className="text-orange-500" />
              <span className="text-[10px] font-bold text-zinc-500 uppercase tracking-wide">People Invited</span>
            </div>
            <p className="text-2xl font-black text-zinc-900 dark:text-white">{stats.invited_count}</p>
          </div>
          <div className="bg-zinc-50 dark:bg-zinc-800 rounded-xl p-3 text-center">
            <div className="flex items-center justify-center gap-1 mb-1">
              <TrendingUp size={12} className="text-green-500" />
              <span className="text-[10px] font-bold text-zinc-500 uppercase tracking-wide">Total Earned</span>
            </div>
            <p className="text-lg font-black text-zinc-900 dark:text-white">
              ETB {Number(stats.total_commission).toLocaleString()}
            </p>
          </div>
        </div>

        {/* Balance card — only shown when user has earnings */}
        {balance && (balance.total_earned > 0 || balance.available_balance > 0) && (
          <div className="bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/10 border border-green-200 dark:border-green-800/40 rounded-xl p-4 space-y-3">
            <p className="text-[10px] font-black uppercase tracking-widest text-green-700 dark:text-green-400 flex items-center gap-1.5">
              <Wallet size={10} /> Referral Balance
            </p>
            <div className="grid grid-cols-3 gap-2 text-center">
              <div>
                <p className="text-xs font-black text-zinc-900 dark:text-white">ETB {Number(balance.total_earned).toLocaleString()}</p>
                <p className="text-[10px] text-zinc-500">Total Earned</p>
              </div>
              <div>
                <p className="text-xs font-black text-zinc-900 dark:text-white">ETB {Number(balance.total_withdrawn).toLocaleString()}</p>
                <p className="text-[10px] text-zinc-500">Withdrawn</p>
              </div>
              <div>
                <p className={`text-xs font-black ${balance.available_balance >= minWithdrawal ? 'text-green-600' : 'text-zinc-500'}`}>
                  ETB {Number(balance.available_balance).toLocaleString()}
                </p>
                <p className="text-[10px] text-zinc-500">Available</p>
              </div>
            </div>
            {balance.pending_amount > 0 && (
              <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800/40 rounded-xl px-3 py-2 space-y-1">
                <div className="flex items-center gap-1.5 text-xs font-bold text-yellow-700 dark:text-yellow-400">
                  <Clock size={11} /> Pending Commission
                </div>
                <p className="text-xs text-yellow-700 dark:text-yellow-300">
                  ETB {Number(balance.pending_amount).toLocaleString()}
                </p>
                <p className="text-[10px] text-yellow-600/70 dark:text-yellow-400/70">
                  New commissions stay pending for 30 days before becoming withdrawable. This protects against payment disputes and chargebacks.
                </p>
              </div>
            )}
            {/* Withdraw button / form */}
            {!isSuspended && balance.available_balance >= minWithdrawal && balance.pending_amount === 0 && !showWithdrawForm && (
              <button onClick={() => { setShowWithdrawForm(true); setWError(''); setWSuccess(''); }}
                className="w-full flex items-center justify-center gap-1.5 py-2 rounded-xl bg-green-600 hover:bg-green-700 text-white text-xs font-bold transition-colors">
                <DollarSign size={12} /> Request Withdrawal
              </button>
            )}
            {isSuspended && balance.available_balance >= 500 && (
              <div className="flex items-center gap-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl px-3 py-2">
                <AlertCircle size={12} className="text-red-500 shrink-0" />
                <p className="text-[10px] text-red-600 dark:text-red-400">Withdrawals blocked — referral code is suspended</p>
              </div>
            )}
            {!isSuspended && balance.available_balance < minWithdrawal && balance.pending_amount === 0 && (
              <p className="text-[10px] text-zinc-400 text-center">
                Minimum withdrawal is ETB {minWithdrawal} — earn ETB {Math.max(0, Math.ceil(minWithdrawal - balance.available_balance)).toLocaleString()} more to withdraw
              </p>
            )}
            {wSuccess && (
              <div className="flex items-start gap-2 bg-green-100 dark:bg-green-900/30 rounded-xl px-3 py-2">
                <CheckCircle2 size={12} className="text-green-600 shrink-0 mt-0.5" />
                <p className="text-xs text-green-700 dark:text-green-400">{wSuccess}</p>
              </div>
            )}
            {!isSuspended && showWithdrawForm && (
              <form onSubmit={handleWithdraw} className="space-y-2 pt-1 border-t border-green-200 dark:border-green-800/40">
                <p className="text-xs font-bold text-zinc-700 dark:text-zinc-300">Withdrawal Request</p>
                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <label className="text-[10px] font-bold text-zinc-500 block mb-1">
                      Amount (ETB) * <span className="font-normal text-zinc-400">min {minWithdrawal}, max {balance ? Math.min(maxWithdrawal, Math.floor(balance.available_balance)).toLocaleString() : maxWithdrawal.toLocaleString()}</span>
                    </label>
                    <input type="number" min={minWithdrawal} max={balance ? Math.min(maxWithdrawal, Math.floor(balance.available_balance)) : maxWithdrawal} step="1"
                      value={wAmount} onChange={e => setWAmount(e.target.value)}
                      placeholder={minWithdrawal.toString()}
                      className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-sm bg-white dark:bg-zinc-800 outline-none focus:ring-2 focus:ring-green-500/30" />
                  </div>
                  <div>
                    <label className="text-[10px] font-bold text-zinc-500 block mb-1">Method *</label>
                    <select value={wMethod} onChange={e => setWMethod(e.target.value)}
                      className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-sm bg-white dark:bg-zinc-800 outline-none focus:ring-2 focus:ring-green-500/30">
                      <option value="telebirr">Telebirr</option>
                      <option value="safaricom">Safaricom</option>
                      <option value="cbe">CBE</option>
                      <option value="bank">Bank</option>
                    </select>
                  </div>
                </div>
                <div>
                  <label className="text-[10px] font-bold text-zinc-500 block mb-1">Account / Phone *</label>
                  <input value={wPhone} onChange={e => setWPhone(e.target.value)} placeholder="0912345678 or account number"
                    className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-sm bg-white dark:bg-zinc-800 outline-none focus:ring-2 focus:ring-green-500/30" />
                </div>
                {wError && (
                  <div className="flex items-center gap-1.5 text-xs text-red-500 bg-red-50 dark:bg-red-900/20 px-3 py-2 rounded-xl">
                    <AlertCircle size={11} className="shrink-0" /> {wError}
                  </div>
                )}
                <div className="flex gap-2">
                  <button type="submit" disabled={wLoading}
                    className="flex-1 py-2 rounded-xl bg-green-600 hover:bg-green-700 text-white text-xs font-bold transition-colors disabled:opacity-50 flex items-center justify-center gap-1.5">
                    {wLoading ? <Loader2 size={12} className="animate-spin" /> : <DollarSign size={12} />}
                    {wLoading ? 'Submitting…' : 'Submit Request'}
                  </button>
                  <button type="button" onClick={() => { setShowWithdrawForm(false); setWError(''); }}
                    className="px-3 py-2 text-zinc-400 hover:text-zinc-600 text-xs transition-colors">Cancel</button>
                </div>
              </form>
            )}
          </div>
        )}

        {/* Referral link */}
        <div>
          <p className="text-[10px] font-black uppercase tracking-widest text-zinc-500 dark:text-zinc-400 mb-1.5 flex items-center gap-1">
            <Link2 size={9} /> Your Referral Link
          </p>
          {isSuspended ? (
            <div className="flex items-start gap-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl px-3 py-2.5">
              <AlertCircle size={13} className="text-red-500 shrink-0 mt-0.5" />
              <div>
                <p className="text-xs font-bold text-red-700 dark:text-red-400">Referral code suspended</p>
                <p className="text-[10px] text-red-500/70 mt-0.5">
                  Your referral code has been suspended by an admin. New invitations cannot be accepted and no commissions will be earned. Contact support if you believe this is an error.
                </p>
              </div>
            </div>
          ) : (
            <>
              <div className="flex gap-2">
                <div className="flex-1 bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2.5 text-xs text-zinc-600 dark:text-zinc-300 truncate font-mono select-all">
                  {stats.referral_link}
                </div>
                <button
                  onClick={handleCopy}
                  className="flex items-center gap-1.5 px-3 py-2.5 bg-orange-600 hover:bg-orange-700 text-white rounded-xl text-xs font-bold transition-colors shrink-0"
                >
                  {copied ? <CheckCircle2 size={13} /> : <Copy size={13} />}
                  {copied ? 'Copied!' : 'Copy'}
                </button>
              </div>
              <p className="text-[10px] text-zinc-400 mt-1">
                Code: <span className="font-mono font-bold text-zinc-600 dark:text-zinc-300">{stats.referral_code}</span>
                <span className="ml-2 text-zinc-300 dark:text-zinc-600">· Share this link to invite people</span>
              </p>
            </>
          )}
        </div>

        {/* ── Inviter section ─────────────────────────────────────────────── */}
        {hasInviter ? (
          // Already has an inviter — show who invited them
          <div className="flex items-center gap-2 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-xl px-3 py-2.5">
            <CheckCircle2 size={13} className="text-green-600 shrink-0" />
            <div>
              <p className="text-xs text-green-700 dark:text-green-400">
                Invited by <span className="font-bold">{stats.inviter_name ?? stats.inviter_code}</span>
              </p>
              <p className="text-[10px] text-green-600/60 dark:text-green-500/50 mt-0.5">
                They earn {commissionRate}% when you make a paid action
              </p>
            </div>
          </div>
        ) : (
          // No inviter yet — show accept form or success message
          <div className="space-y-2">
            {acceptSuccess ? (
              <div className="flex items-start gap-2 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-xl px-3 py-2.5">
                <CheckCircle2 size={13} className="text-green-600 shrink-0 mt-0.5" />
                <p className="text-xs text-green-700 dark:text-green-400">{acceptSuccess}</p>
              </div>
            ) : !showAcceptForm ? (
              <button
                onClick={() => { setShowAcceptForm(true); setAcceptError(''); }}
                className="flex items-center gap-2 w-full border border-dashed border-zinc-300 dark:border-zinc-700 rounded-xl px-3 py-2.5 text-xs text-zinc-500 dark:text-zinc-400 hover:border-orange-400 hover:text-orange-500 transition-colors"
              >
                <UserPlus size={13} />
                Accept an Invitation — enter someone's referral code
              </button>
            ) : (
              <form onSubmit={handleAccept} className="space-y-2">
                <p className="text-xs font-bold text-zinc-700 dark:text-zinc-300">
                  Enter inviter's referral code
                </p>
                <div className="flex gap-2">
                  <input
                    autoFocus
                    value={inviteCode}
                    onChange={e => { setInviteCode(e.target.value.toUpperCase()); setAcceptError(''); }}
                    placeholder="REFxxxxxxxx"
                    maxLength={11}
                    className="flex-1 border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-sm font-mono outline-none focus:ring-2 focus:ring-orange-500/30 focus:border-orange-500/60 bg-zinc-50 dark:bg-zinc-800 transition-all"
                  />
                  <button
                    type="submit"
                    disabled={acceptLoading || inviteCode.length < 11}
                    className="flex items-center gap-1.5 px-3 py-2 bg-orange-600 hover:bg-orange-700 text-white rounded-xl text-xs font-bold transition-colors disabled:opacity-50"
                  >
                    {acceptLoading ? <Loader2 size={12} className="animate-spin" /> : 'Accept'}
                  </button>
                  <button
                    type="button"
                    onClick={() => { setShowAcceptForm(false); setInviteCode(''); setAcceptError(''); }}
                    className="px-3 py-2 text-zinc-400 hover:text-zinc-600 text-xs transition-colors"
                  >
                    Cancel
                  </button>
                </div>
                {acceptError && (
                  <div className="flex items-center gap-1.5 text-xs text-red-500 bg-red-50 dark:bg-red-900/20 px-3 py-2 rounded-xl">
                    <AlertCircle size={11} className="shrink-0" /> {acceptError}
                  </div>
                )}
                <p className="text-[10px] text-zinc-400">
                  This links you to your inviter permanently — cannot be changed later.
                </p>
              </form>
            )}
          </div>
        )}

        {/* ── Withdrawal history ───────────────────────────────────────────── */}
        {balance && balance.total_withdrawn > 0 && (
          <div className="border-t border-zinc-100 dark:border-zinc-800 pt-3">
            <button
              onClick={() => { const next = !showWithdrawals; setShowWithdrawals(next); if (next) loadWithdrawals(); }}
              className="flex items-center justify-between w-full text-xs font-bold text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">
              <span className="flex items-center gap-1.5"><Wallet size={12} /> Withdrawal History</span>
              {showWithdrawals ? <ChevronUp size={13} /> : <ChevronDown size={13} />}
            </button>
            {showWithdrawals && (
              <div className="mt-2 space-y-1.5 max-h-48 overflow-y-auto pr-1">
                {!withdrawalsLoaded ? (
                  <div className="flex items-center justify-center py-4"><Loader2 size={16} className="animate-spin text-zinc-400" /></div>
                ) : withdrawals.length === 0 ? (
                  <p className="text-xs text-zinc-400 text-center py-3">No withdrawals yet.</p>
                ) : (
                  withdrawals.map(w => (
                    <div key={w.id} className="flex items-center justify-between bg-zinc-50 dark:bg-zinc-800 rounded-xl px-3 py-2">
                      <div className="min-w-0 flex-1">
                        <div className="flex items-center gap-1.5">
                          <span className={`text-[9px] font-black px-1.5 py-0.5 rounded-full uppercase ${STATUS_COLOR[w.status] ?? ''}`}>{w.status}</span>
                          <span className="text-[10px] text-zinc-400">{PM_LABEL[w.payment_method] ?? w.payment_method} · {w.payer_phone}</span>
                        </div>
                        {w.admin_note && <p className="text-[10px] text-zinc-400 mt-0.5 truncate">{w.admin_note}</p>}
                        {w.reference && <p className="text-[10px] text-green-600 font-mono mt-0.5">Ref: {w.reference}</p>}
                      </div>
                      <div className="text-right shrink-0 ml-3">
                        <p className="text-xs font-black text-zinc-800 dark:text-zinc-200">ETB {Number(w.amount).toLocaleString()}</p>
                        <p className="text-[10px] text-zinc-400">{new Date(w.requested_at).toLocaleDateString()}</p>
                      </div>
                    </div>
                  ))
                )}
              </div>
            )}
          </div>
        )}

        {/* ── Commission history ───────────────────────────────────────────── */}
        {hasEarnings && (
          <div className="border-t border-zinc-100 dark:border-zinc-800 pt-3">
            <button
              onClick={() => {
                const next = !showCommissions;
                setShowCommissions(next);
                if (next) loadCommissions();
              }}
              className="flex items-center justify-between w-full text-xs font-bold text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors"
            >
              <span>Commission History</span>
              {showCommissions ? <ChevronUp size={13} /> : <ChevronDown size={13} />}
            </button>

            {showCommissions && (
              <div className="mt-2 space-y-1.5 max-h-52 overflow-y-auto pr-1">
                {!commissionsLoaded ? (
                  <div className="flex items-center justify-center py-4">
                    <Loader2 size={16} className="animate-spin text-zinc-400" />
                  </div>
                ) : commissions.length === 0 ? (
                  <div className="text-center py-4">
                    <p className="text-xs text-zinc-400">No commissions yet.</p>
                    <p className="text-[10px] text-zinc-300 dark:text-zinc-600 mt-1">
                      Commissions are recorded when your invitees make a paid action.
                    </p>
                  </div>
                ) : (
                  commissions.map(c => (
                    <div key={c.id} className="flex items-center justify-between bg-zinc-50 dark:bg-zinc-800 rounded-xl px-3 py-2">
                      <div className="min-w-0 flex-1">
                        <div className="flex items-center gap-1.5">
                          <p className="text-xs font-bold text-zinc-800 dark:text-zinc-200 truncate">
                            {c.invited_user_name || 'User'}
                          </p>
                          {c.status && c.status !== 'confirmed' && (
                            <span className={`text-[9px] font-black px-1.5 py-0.5 rounded-full uppercase ${COMMISSION_STATUS_COLOR[c.status] ?? ''}`}>
                              {c.status}
                            </span>
                          )}
                        </div>
                        <p className="text-[10px] text-zinc-400">
                          {TX_TYPE_LABELS[c.transaction_type] ?? c.transaction_type.replace(/_/g, ' ')}
                          {' · '}ETB {Number(c.transaction_amount).toLocaleString()}
                        </p>
                        {c.status === 'pending' && c.pending_until && (
                          <p className="text-[9px] text-yellow-600 dark:text-yellow-400 mt-0.5 flex items-center gap-1">
                            <Clock size={8} /> Withdrawable after {new Date(c.pending_until).toLocaleDateString()}
                          </p>
                        )}
                        {c.status === 'reversed' && (
                          <p className="text-[9px] text-red-600 dark:text-red-400 mt-0.5">
                            Commission reversed (payment dispute)
                          </p>
                        )}
                      </div>
                      <div className="text-right shrink-0 ml-3">
                        <p className={`text-xs font-black ${c.status === 'reversed' ? 'text-red-500 line-through' : 'text-green-600'}`}>
                          +ETB {Number(c.commission_amount).toLocaleString()}
                        </p>
                        <p className="text-[10px] text-zinc-400">
                          {new Date(c.created_at).toLocaleDateString()}
                        </p>
                      </div>
                    </div>
                  ))
                )}
              </div>
            )}
          </div>
        )}

        {/* How it works — shown when no invitees yet */}
        {stats.invited_count === 0 && stats.total_commission === 0 && (
          <div className="bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-3 py-3 space-y-1.5">
            <p className="text-[10px] font-black uppercase tracking-widest text-zinc-400">How it works</p>
            {[
              'Share your referral link with friends',
              'They register using your link',
              `You earn ${commissionRate}% every time they make a paid action`,
              'Paid actions: product sale, tier upgrade, seller registration',
            ].map((step, i) => (
              <div key={i} className="flex items-start gap-2">
                <span className="text-[10px] font-black text-orange-500 shrink-0 mt-0.5">{i + 1}.</span>
                <p className="text-[10px] text-zinc-500 dark:text-zinc-400">{step}</p>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
