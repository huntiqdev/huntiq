'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { 
  Settings, Save, AlertCircle, Shield, DollarSign, 
  Layers, CreditCard, Megaphone, CheckCircle2, ShieldAlert, Calculator 
} from 'lucide-react';
import { useRouter } from 'next/navigation';
import { IPTrackingSecurity } from '@/components/admin/enterprise/IPTrackingSecurity';
import { FeeControlEngine } from '@/components/admin/enterprise/FeeControlEngine';

interface SettingMeta {
  label: string;
  category: 'general' | 'fees' | 'tiers' | 'merchants' | 'promotions';
  type: 'boolean' | 'number' | 'text';
  description?: string;
}

const SETTING_METADATA: Record<string, SettingMeta> = {
  // General & Security
  maintenance_mode: { label: 'Maintenance Mode', category: 'general', type: 'boolean', description: 'Temporarily restrict platform access for users' },
  otp_required: { label: 'Enforce OTP Verification', category: 'general', type: 'boolean', description: 'Require SMS/Email OTP code on registration & login' },
  otp_gateway: { label: 'OTP Gateway', category: 'general', type: 'text', description: 'Delivery method: "log" (dev mode), "sms", or "email"' },
  otp_expiry_minutes: { label: 'OTP Expiry (Minutes)', category: 'general', type: 'number', description: 'Time before an OTP code expires' },
  ai_recommendations_enabled: { label: 'AI Recommendations', category: 'general', type: 'boolean', description: 'Enable personalized product recommendations' },

  // Fees & Commissions
  broker_fee_rate: { label: 'Broker Fee Rate (%)', category: 'fees', type: 'number', description: 'Platform fee percent charged on completed sales' },
  referral_commission_rate: { label: 'Referral Commission (%)', category: 'fees', type: 'number', description: 'Commission paid to referrers' },
  min_referral_withdrawal: { label: 'Min Referral Withdrawal (ETB)', category: 'fees', type: 'number', description: 'Minimum balance required to request cashout' },
  max_referral_withdrawal: { label: 'Max Referral Withdrawal (ETB)', category: 'fees', type: 'number', description: 'Maximum amount allowed per withdrawal request' },

  // Seller Tiers & Listings
  listing_free_weeks: { label: 'Free Listing Grace Period (Weeks)', category: 'tiers', type: 'number', description: 'Initial duration before weekly fee increments begin' },
  normal_startup_fee: { label: 'Normal Tier Startup Fee (%)', category: 'tiers', type: 'number' },
  gold_startup_fee: { label: 'Gold Tier Startup Fee (%)', category: 'tiers', type: 'number' },
  premium_startup_fee: { label: 'Premium Tier Startup Fee (%)', category: 'tiers', type: 'number' },

  // Merchant Payments
  telebirr_merchant_phone: { label: 'Telebirr Merchant Phone', category: 'merchants', type: 'text', description: 'Official Telebirr phone number for listing payments' },
  safaricom_merchant_phone: { label: 'M-PESA Merchant Phone', category: 'merchants', type: 'text', description: 'Official Safaricom M-PESA phone number' },
  cbe_account_number: { label: 'CBE Bank Account Number', category: 'merchants', type: 'text', description: 'Commercial Bank of Ethiopia account number' },
};

const CATEGORIES = [
  { id: 'general', label: 'General & Security', icon: Shield },
  { id: 'ip_security', label: 'IP Defense Firewall', icon: ShieldAlert },
  { id: 'fee_engine', label: 'Deep Fee Engine', icon: Calculator },
  { id: 'fees', label: 'Fees & Commissions', icon: DollarSign },
  { id: 'tiers', label: 'Tiers & Limits', icon: Layers },
  { id: 'merchants', label: 'Merchant Accounts', icon: CreditCard },
  { id: 'promotions', label: 'Promotions Banner', icon: Megaphone },
] as const;

export default function AdminSettingsPage() {
  const { user, isSuperAdmin, hasAdminPermission, loading: authLoading } = useAuth();
  const router = useRouter();
  
  const [settingsMap, setSettingsMap] = useState<Record<string, string>>({});
  const [activeCategory, setActiveCategory] = useState<string>('general');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const canAccess = isSuperAdmin || hasAdminPermission('settings');

  useEffect(() => {
    if (!authLoading && !canAccess) {
      router.push('/admin/dashboard');
      return;
    }

    const load = async () => {
      if (!user) return;
      try {
        const res = await api.adminGetSettings(user.token);
        setSettingsMap(res || {});
      } catch (err: any) {
        setError('Failed to load platform settings');
      } finally {
        setLoading(false);
      }
    };
    if (canAccess) load();
  }, [user, canAccess, authLoading, router]);

  const handleChange = (key: string, value: string) => {
    setSettingsMap((prev) => ({ ...prev, [key]: value }));
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) return;
    setSaving(true);
    setError('');
    setSuccess('');
    try {
      await api.adminUpdateSettings(settingsMap, user.token);
      setSuccess('Settings updated successfully');
      setTimeout(() => setSuccess(''), 3000);
    } catch (err: any) {
      setError(err?.message || 'Failed to save settings');
    } finally {
      setSaving(false);
    }
  };

  if (authLoading || loading) {
    return <div className="h-64 flex items-center justify-center"><Spinner /></div>;
  }

  return (
    <div className="space-y-6 max-w-6xl">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-purple-600 rounded-2xl flex items-center justify-center text-white shadow-lg shadow-purple-600/20">
            <Settings size={24} />
          </div>
          <div>
            <h1 className="text-2xl font-black text-white leading-none">Enterprise Platform Settings & Security</h1>
            <p className="text-xs font-bold text-zinc-400 uppercase tracking-widest mt-1">
              {isSuperAdmin ? 'Super Admin Full Control' : 'Manager Scope'}
            </p>
          </div>
        </div>
      </div>

      {/* Category Navigation Tabs */}
      <div className="flex flex-wrap gap-2 border-b border-zinc-800 pb-3">
        {CATEGORIES.map((cat) => {
          const Icon = cat.icon;
          const isActive = activeCategory === cat.id;
          return (
            <button
              key={cat.id}
              type="button"
              onClick={() => setActiveCategory(cat.id)}
              className={`flex items-center gap-2 px-4 py-2.5 rounded-xl font-extrabold text-xs uppercase tracking-wider transition-all ${
                isActive 
                  ? 'bg-orange-600 text-white shadow-lg shadow-orange-600/20' 
                  : 'bg-zinc-900 text-zinc-400 hover:bg-zinc-800 hover:text-white'
              }`}
            >
              <Icon size={16} />
              {cat.label}
            </button>
          );
        })}
      </div>

      {/* Special Module Views */}
      {activeCategory === 'ip_security' && (
        <IPTrackingSecurity />
      )}

      {activeCategory === 'fee_engine' && (
        <FeeControlEngine />
      )}

      {/* Standard Settings View */}
      {activeCategory !== 'ip_security' && activeCategory !== 'fee_engine' && (
        <form onSubmit={handleSave} className="bg-zinc-900 border border-zinc-800 rounded-3xl overflow-hidden shadow-xl">
          <div className="p-6 space-y-6 divide-y divide-zinc-800">
            {Object.keys(settingsMap)
              .filter((key) => (SETTING_METADATA[key]?.category || 'general') === activeCategory)
              .map((key, idx) => {
                const meta = SETTING_METADATA[key] || {
                  label: key.replace(/_/g, ' ').toUpperCase(),
                  type: 'text',
                  category: 'general'
                };
                const rawVal = settingsMap[key] ?? '';
                const isBool = meta.type === 'boolean' || rawVal === 'true' || rawVal === 'false';

                return (
                  <div key={key} className={`${idx === 0 ? '' : 'pt-6'} flex flex-col md:flex-row md:items-center justify-between gap-4`}>
                    <div className="max-w-md">
                      <label className="block text-sm font-black text-zinc-100">
                        {meta.label}
                      </label>
                      {meta.description && (
                        <p className="text-xs text-zinc-400 font-medium mt-1">
                          {meta.description}
                        </p>
                      )}
                    </div>

                    <div className="w-full md:w-64">
                      {isBool ? (
                        <button
                          type="button"
                          onClick={() => handleChange(key, rawVal === 'true' ? 'false' : 'true')}
                          className={`w-full py-3 px-4 rounded-xl font-black text-xs uppercase tracking-wider flex items-center justify-between border transition-all ${
                            rawVal === 'true'
                              ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-400'
                              : 'bg-zinc-950 border-zinc-800 text-zinc-500'
                          }`}
                        >
                          <span>{rawVal === 'true' ? 'Enabled' : 'Disabled'}</span>
                        </button>
                      ) : (
                        <input
                          type={meta.type === 'number' ? 'number' : 'text'}
                          value={rawVal}
                          onChange={(e) => handleChange(key, e.target.value)}
                          className="w-full bg-zinc-950 border border-zinc-800 text-white rounded-xl px-4 py-3 font-mono text-sm focus:outline-none focus:border-orange-500"
                        />
                      )}
                    </div>
                  </div>
                );
              })}
          </div>

          <div className="bg-zinc-950 p-6 border-t border-zinc-800 flex items-center justify-between">
            <button
              type="submit"
              disabled={saving}
              className="flex items-center gap-2 bg-orange-600 hover:bg-orange-500 text-white font-black px-6 py-3 rounded-xl transition-all shadow-lg shadow-orange-600/20"
            >
              {saving ? <Spinner size={16} className="text-white" /> : <Save size={16} />}
              Save Configuration
            </button>
          </div>
        </form>
      )}
    </div>
  );
}
