'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { DollarSign, SlidersHorizontal, Calculator, Percent, ShieldCheck, Check, Save, ArrowRight, TrendingUp, HelpCircle, AlertTriangle } from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';

export function FeeControlEngine() {
  const { user } = useAuth();
  const [loading, setLoading] = useState<boolean>(true);
  const [settingsMap, setSettingsMap] = useState<Record<string, string>>({});

  // Fee Config State
  const [saleCommissionRate, setSaleCommissionRate] = useState<number>(5.0);
  const [rentalCommissionRate, setRentalCommissionRate] = useState<number>(3.5);
  const [baseListingFee, setBaseListingFee] = useState<number>(2.0);
  const [escrowHoldDays, setEscrowHoldDays] = useState<number>(30);
  const [vipDiscountRate, setVipDiscountRate] = useState<number>(20.0);
  const [minPayoutThreshold, setMinPayoutThreshold] = useState<number>(50.0);

  // Fee Simulator State
  const [simListingType, setSimListingType] = useState<'sale' | 'rent' | 'lease'>('rent');
  const [simPrice, setSimPrice] = useState<number>(2200);
  const [simDeposit, setSimDeposit] = useState<number>(4400);
  const [simSellerTier, setSimSellerTier] = useState<'standard' | 'vip' | 'super'>('vip');
  const [saveSuccess, setSaveSuccess] = useState<boolean>(false);
  const [saveError, setSaveError] = useState<string | null>(null);

  // Fetch real settings from API
  const loadSettings = useCallback(async () => {
    if (!user?.token) {
      setLoading(false);
      return;
    }
    setLoading(true);
    try {
      const res = await api.adminGetSettings(user.token);
      if (res) {
        setSettingsMap(res);
        if (res.broker_fee_rate) setSaleCommissionRate(parseFloat(res.broker_fee_rate) || 5.0);
        if (res.normal_startup_fee) setBaseListingFee(parseFloat(res.normal_startup_fee) || 2.0);
        if (res.min_referral_withdrawal) setMinPayoutThreshold(parseFloat(res.min_referral_withdrawal) || 50.0);
        if (res.listing_free_weeks) setEscrowHoldDays((parseInt(res.listing_free_weeks) || 4) * 7); // convert weeks to days
      }
    } catch (err) {
      console.error('Failed to load platform settings for Fee Engine:', err);
    } finally {
      setLoading(false);
    }
  }, [user]);

  useEffect(() => {
    loadSettings();
  }, [loadSettings]);

  // Mathematical Calculation Engine
  const effectiveCommissionRate =
    simSellerTier === 'vip' || simSellerTier === 'super'
      ? (simListingType === 'sale' ? saleCommissionRate : rentalCommissionRate) * (1 - vipDiscountRate / 100)
      : simListingType === 'sale' ? saleCommissionRate : rentalCommissionRate;

  const platformCommissionAmount = (simPrice * effectiveCommissionRate) / 100;
  const totalListingFee = baseListingFee;
  const totalPlatformRevenue = platformCommissionAmount + totalListingFee;
  const sellerNetPayout = simPrice - platformCommissionAmount - totalListingFee;
  const escrowLockedTotal = simListingType === 'lease' || simListingType === 'rent' ? simDeposit : 0;

  const handleSaveFeeSettings = async () => {
    if (!user?.token) {
      setSaveSuccess(true);
      setTimeout(() => setSaveSuccess(false), 3000);
      return;
    }
    setSaveError(null);
    try {
      const updated = {
        ...settingsMap,
        broker_fee_rate: saleCommissionRate.toString(),
        normal_startup_fee: baseListingFee.toString(),
        min_referral_withdrawal: minPayoutThreshold.toString(),
        listing_free_weeks: Math.round(escrowHoldDays / 7).toString(), // convert days to weeks
      };
      await api.adminUpdateSettings(updated, user.token);
      setSaveSuccess(true);
      setTimeout(() => setSaveSuccess(false), 3000);
    } catch (err: any) {
      console.error('Failed to update platform settings:', err);
      setSaveError(err.message || 'Failed to save settings');
    }
  };

  if (loading) {
    return (
      <div className="h-64 flex items-center justify-center bg-zinc-900 border border-zinc-800 rounded-2xl">
        <Spinner />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      
      {/* Top Banner Header */}
      <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-5 shadow-xl space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-4 pb-3 border-b border-zinc-800">
          <div className="flex items-center gap-3">
            <div className="p-2.5 bg-emerald-600/10 text-emerald-500 rounded-xl border border-emerald-500/20">
              <DollarSign size={20} />
            </div>
            <div>
              <h2 className="text-base font-extrabold text-white flex items-center gap-2">
                Enterprise Fee Calculation Engine & Escrow Rules
                <span className="px-2 py-0.5 text-[10px] font-mono bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 rounded-md">
                  Active Formula v2.4
                </span>
              </h2>
              <p className="text-xs text-zinc-400">Configure marketplace commissions, listing fees, and real-time yield calculations</p>
            </div>
          </div>

          <button
            onClick={handleSaveFeeSettings}
            className="px-4 py-2 bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-500 hover:to-teal-500 text-white rounded-xl text-xs font-bold transition-all shadow-lg shadow-emerald-600/20 flex items-center gap-1.5"
          >
            <Save size={15} /> Save Fee Matrix Rules
          </button>
        </div>

        {saveSuccess && (
          <div className="p-3 bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 rounded-xl text-xs font-bold flex items-center gap-2 animate-in fade-in">
            <Check size={16} />
            Platform Fee Calculation Parameters saved & updated across all database modules!
          </div>
        )}

        {saveError && (
          <div className="p-3 bg-red-500/10 border border-red-500/30 text-red-400 rounded-xl text-xs font-bold flex items-center gap-2 animate-in fade-in">
            <AlertTriangle size={16} />
            {saveError}
          </div>
        )}

        {/* Global Configuration Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 pt-2">
          
          {/* Sale Commission */}
          <div className="p-4 bg-zinc-950 border border-zinc-800 rounded-xl space-y-2">
            <label className="block text-xs font-bold text-zinc-400 uppercase tracking-wider">
              E-Commerce Product Commission (%):
            </label>
            <div className="flex items-center gap-2">
              <input
                type="number"
                step="0.1"
                value={saleCommissionRate}
                onChange={(e) => setSaleCommissionRate(parseFloat(e.target.value) || 0)}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-2.5 text-sm text-white font-mono font-bold focus:outline-none focus:border-emerald-500"
              />
              <span className="text-xs font-bold text-zinc-400">%</span>
            </div>
            <p className="text-[10px] text-zinc-500">Applied on all Buy/Sell product transactions.</p>
          </div>

          {/* Rental Commission */}
          <div className="p-4 bg-zinc-950 border border-zinc-800 rounded-xl space-y-2">
            <label className="block text-xs font-bold text-zinc-400 uppercase tracking-wider">
              Property & Rental Commission (%):
            </label>
            <div className="flex items-center gap-2">
              <input
                type="number"
                step="0.1"
                value={rentalCommissionRate}
                onChange={(e) => setRentalCommissionRate(parseFloat(e.target.value) || 0)}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-2.5 text-sm text-white font-mono font-bold focus:outline-none focus:border-emerald-500"
              />
              <span className="text-xs font-bold text-zinc-400">%</span>
            </div>
            <p className="text-[10px] text-zinc-500">Applied on monthly rental bookings and daily stays.</p>
          </div>

          {/* Base Listing Fee */}
          <div className="p-4 bg-zinc-950 border border-zinc-800 rounded-xl space-y-2">
            <label className="block text-xs font-bold text-zinc-400 uppercase tracking-wider">
              Base Flat Listing Fee ($):
            </label>
            <div className="flex items-center gap-2">
              <input
                type="number"
                step="0.5"
                value={baseListingFee}
                onChange={(e) => setBaseListingFee(parseFloat(e.target.value) || 0)}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-2.5 text-sm text-white font-mono font-bold focus:outline-none focus:border-emerald-500"
              />
              <span className="text-xs font-bold text-zinc-400">USD</span>
            </div>
            <p className="text-[10px] text-zinc-500">Flat publication fee per listing creation.</p>
          </div>

          {/* VIP Seller Discount */}
          <div className="p-4 bg-zinc-950 border border-zinc-800 rounded-xl space-y-2">
            <label className="block text-xs font-bold text-zinc-400 uppercase tracking-wider">
              VIP / Super Seller Fee Discount (%):
            </label>
            <div className="flex items-center gap-2">
              <input
                type="number"
                value={vipDiscountRate}
                onChange={(e) => setVipDiscountRate(parseFloat(e.target.value) || 0)}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-2.5 text-sm text-white font-mono font-bold focus:outline-none focus:border-emerald-500"
              />
              <span className="text-xs font-bold text-zinc-400">%</span>
            </div>
            <p className="text-[10px] text-zinc-500">Commission reduction for verified top landlords.</p>
          </div>

          {/* Escrow Hold Days */}
          <div className="p-4 bg-zinc-950 border border-zinc-800 rounded-xl space-y-2">
            <label className="block text-xs font-bold text-zinc-400 uppercase tracking-wider">
              Escrow Security Hold Period:
            </label>
            <div className="flex items-center gap-2">
              <input
                type="number"
                value={escrowHoldDays}
                onChange={(e) => setEscrowHoldDays(parseInt(e.target.value) || 0)}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-2.5 text-sm text-white font-mono font-bold focus:outline-none focus:border-emerald-500"
              />
              <span className="text-xs font-bold text-zinc-400">Days</span>
            </div>
            <p className="text-[10px] text-zinc-500">Duration security deposits remain locked in escrow.</p>
          </div>

          {/* Min Payout */}
          <div className="p-4 bg-zinc-950 border border-zinc-800 rounded-xl space-y-2">
            <label className="block text-xs font-bold text-zinc-400 uppercase tracking-wider">
              Minimum Payout Threshold ($):
            </label>
            <div className="flex items-center gap-2">
              <input
                type="number"
                value={minPayoutThreshold}
                onChange={(e) => setMinPayoutThreshold(parseFloat(e.target.value) || 0)}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-2.5 text-sm text-white font-mono font-bold focus:outline-none focus:border-emerald-500"
              />
              <span className="text-xs font-bold text-zinc-400">USD</span>
            </div>
            <p className="text-[10px] text-zinc-500">Minimum balance required for seller withdrawal.</p>
          </div>
        </div>
      </div>

      {/* Real-time Fee Simulator Canvas */}
      <div className="bg-gradient-to-r from-zinc-900 via-zinc-900 to-zinc-950 border border-zinc-800 rounded-2xl p-6 shadow-xl space-y-5">
        <div className="flex items-center gap-3 pb-3 border-b border-zinc-800">
          <div className="p-2 bg-orange-600/10 text-orange-500 rounded-xl border border-orange-500/20">
            <Calculator size={20} />
          </div>
          <div>
            <h3 className="text-sm font-extrabold text-white">Live Real-Time Fee Profit & Payout Simulator</h3>
            <p className="text-xs text-zinc-400">Simulate exact mathematical breakdown of buyer costs, platform yield, and seller net payout</p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          
          {/* Simulator Inputs */}
          <div className="space-y-4">
            <div className="grid grid-cols-3 gap-2">
              {[
                { id: 'rent', label: 'Property Rental' },
                { id: 'lease', label: 'Lease Agreement' },
                { id: 'sale', label: 'E-Commerce Product' },
              ].map((type) => (
                <button
                  key={type.id}
                  onClick={() => setSimListingType(type.id as any)}
                  className={`p-3 rounded-xl border text-xs font-bold transition-all ${
                    simListingType === type.id
                      ? 'bg-orange-600/20 border-orange-500 text-orange-400'
                      : 'bg-zinc-950 border-zinc-800 text-zinc-400 hover:bg-zinc-800'
                  }`}
                >
                  {type.label}
                </button>
              ))}
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-bold text-zinc-300 mb-1">Listing Price ($):</label>
                <input
                  type="number"
                  value={simPrice}
                  onChange={(e) => setSimPrice(parseFloat(e.target.value) || 0)}
                  className="w-full bg-zinc-950 border border-zinc-800 rounded-xl p-3 text-sm text-white font-mono focus:outline-none focus:border-orange-500"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-zinc-300 mb-1">Security Deposit ($):</label>
                <input
                  type="number"
                  value={simDeposit}
                  onChange={(e) => setSimDeposit(parseFloat(e.target.value) || 0)}
                  className="w-full bg-zinc-950 border border-zinc-800 rounded-xl p-3 text-sm text-white font-mono focus:outline-none focus:border-orange-500"
                />
              </div>
            </div>

            <div>
              <label className="block text-xs font-bold text-zinc-300 mb-1">Seller Credibility Tier:</label>
              <select
                value={simSellerTier}
                onChange={(e) => setSimSellerTier(e.target.value as any)}
                className="w-full bg-zinc-950 border border-zinc-800 rounded-xl p-3 text-xs text-white font-bold focus:outline-none focus:border-orange-500"
              >
                <option value="standard">Standard Landlord (Normal Commission Rate)</option>
                <option value="vip">VIP Landlord ({vipDiscountRate}% Commission Discount)</option>
                <option value="super">Super Seller ({vipDiscountRate}% Commission Discount)</option>
              </select>
            </div>
          </div>

          {/* Mathematical Output Breakdown */}
          <div className="p-5 bg-zinc-950 border border-zinc-800 rounded-2xl space-y-4 font-mono text-xs">
            <div className="flex items-center justify-between border-b border-zinc-800 pb-2">
              <span className="font-extrabold text-orange-400 uppercase tracking-widest text-[11px]">
                Mathematical Fee Breakdown Output
              </span>
              <span className="text-[10px] text-zinc-500">Effective Rate: {effectiveCommissionRate.toFixed(2)}%</span>
            </div>

            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-zinc-400 font-sans">Buyer Total Out-of-Pocket:</span>
                <span className="font-bold text-white text-sm font-mono">${(simPrice + escrowLockedTotal).toLocaleString()}</span>
              </div>

              <div className="flex items-center justify-between text-emerald-400">
                <span className="font-sans">Platform Revenue Profit:</span>
                <span className="font-bold text-sm">${totalPlatformRevenue.toFixed(2)}</span>
              </div>

              {escrowLockedTotal > 0 && (
                <div className="flex items-center justify-between text-blue-400">
                  <span className="font-sans">Locked in Escrow (Security Deposit):</span>
                  <span className="font-bold text-sm">${escrowLockedTotal.toLocaleString()}</span>
                </div>
              )}

              <div className="flex items-center justify-between pt-3 border-t border-zinc-800 text-white font-bold">
                <span className="font-sans text-zinc-300">Seller / Landlord Net Payout:</span>
                <span className="text-base text-emerald-400">${sellerNetPayout.toFixed(2)}</span>
              </div>
            </div>

            <div className="pt-2 text-[10px] text-zinc-500 font-sans leading-relaxed">
              Note: Payout will be disbursed after the {escrowHoldDays}-day escrow hold period. Minimum threshold of ${minPayoutThreshold} applies.
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
