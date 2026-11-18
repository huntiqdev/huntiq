'use client';

import React, { useState } from 'react';
import { X, Shield, ShieldAlert, Check, UserPlus, Key, Mail, Phone, Lock, Sparkles, CheckSquare, Square, Building2 } from 'lucide-react';
import { api } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';

interface AdminRoleManagerModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

interface RoleOption {
  id: string;
  name: string;
  desc: string;
  dept: string;
  defaultPerms: string[];
}

const ENTERPRISE_ROLES: RoleOption[] = [
  // Dept 1: Executive & Governance
  { id: 'gov_compliance_officer', name: 'Gov Compliance Officer', desc: 'Regulatory filings, legal reporting & audit oversight', dept: 'Executive Governance', defaultPerms: ['dashboard', 'finance', 'profit', 'transactions', 'operations', 'logs', 'seller_verification', 'trust_safety', 'appeals', 'refunds'] },
  { id: 'legal_counsel', name: 'Legal Counsel', desc: 'Dispute arbitration, terms enforcement & compliance', dept: 'Executive Governance', defaultPerms: ['dashboard', 'trust_safety', 'appeals', 'refunds', 'seller_verification', 'logs', 'transactions'] },
  { id: 'hr_director', name: 'HR & Staff Director', desc: 'Staff credential governance & internal access oversight', dept: 'Executive Governance', defaultPerms: ['dashboard', 'logs', 'notifications', 'operations'] },

  // Dept 2: Cybersecurity & SOC Defense
  { id: 'ciso', name: 'Chief Info Security Officer', desc: 'SOC defense oversight, security architecture & policies', dept: 'Cybersecurity & SOC', defaultPerms: ['dashboard', 'logs', 'trust_safety', 'operations', 'notifications', 'settings', 'analytics'] },
  { id: 'security_admin', name: 'Security Administrator', desc: 'Threat mitigation, IP blackholing & active defense', dept: 'Cybersecurity & SOC', defaultPerms: ['dashboard', 'logs', 'trust_safety', 'operations', 'analytics'] },
  { id: 'soc_analyst', name: 'SOC Threat Analyst', desc: 'Real-time telemetry, brute-force & suspicious traffic', dept: 'Cybersecurity & SOC', defaultPerms: ['dashboard', 'logs', 'trust_safety', 'operations'] },
  { id: 'fraud_investigator', name: 'Fraud & AML Investigator', desc: 'Scam forensics, transaction tampering & sanctions', dept: 'Cybersecurity & SOC', defaultPerms: ['dashboard', 'trust_safety', 'transactions', 'logs', 'seller_verification', 'refunds', 'appeals'] },
  { id: 'compliance_auditor', name: 'Compliance Auditor', desc: 'Security, operational & regulatory compliance audits', dept: 'Cybersecurity & SOC', defaultPerms: ['dashboard', 'logs', 'transactions', 'profit', 'operations', 'finance'] },

  // Dept 3: Banking, Treasury & Financial Operations
  { id: 'finance_director', name: 'Finance & Treasury Director', desc: 'Treasury management, banking settlements & FP&A', dept: 'Finance & Treasury', defaultPerms: ['dashboard', 'finance', 'profit', 'transactions', 'referrals', 'analytics'] },
  { id: 'transaction_auditor', name: 'Transaction & Ledger Auditor', desc: 'Double-entry audit verification & reconciliation', dept: 'Finance & Treasury', defaultPerms: ['dashboard', 'transactions', 'profit', 'finance'] },
  { id: 'tax_controller', name: 'Tax & NBE Controller', desc: 'Withholding tax compliance & central bank reporting', dept: 'Finance & Treasury', defaultPerms: ['dashboard', 'finance', 'profit', 'transactions', 'analytics'] },
  { id: 'escrow_reconciler', name: 'Escrow Reconciler', desc: 'Order escrow release, dispute holds & bank settlement', dept: 'Finance & Treasury', defaultPerms: ['dashboard', 'transactions', 'refunds', 'profit'] },
  { id: 'finance_support_agent', name: 'Finance Support Agent', desc: 'Merchant payout inquiries, Telebirr/CBE receipts', dept: 'Finance & Treasury', defaultPerms: ['dashboard', 'transactions', 'refunds', 'seller_upgrades', 'buyer_upgrades'] },

  // Dept 4: VIP Concierge & Growth Services
  { id: 'vip_account_director', name: 'VIP Account Director', desc: 'High-volume merchants & enterprise landlord onboarding', dept: 'VIP & Growth', defaultPerms: ['dashboard', 'sellers', 'seller_upgrades', 'buyer_upgrades', 'premium_users', 'analytics', 'referrals'] },
  { id: 'gold_concierge', name: 'Gold Concierge', desc: 'Tier-1 seller support, priority listings & direct contact', dept: 'VIP & Growth', defaultPerms: ['dashboard', 'sellers', 'seller_upgrades', 'products', 'premium_users'] },
  { id: 'premium_concierge', name: 'Premium Concierge', desc: 'Elite seller account management & marketing promotions', dept: 'VIP & Growth', defaultPerms: ['dashboard', 'sellers', 'seller_upgrades', 'buyer_upgrades', 'products', 'banners', 'premium_users'] },
  { id: 'referral_director', name: 'Referral Director', desc: 'Affiliate commission tiers, viral loops & payout audits', dept: 'VIP & Growth', defaultPerms: ['dashboard', 'referrals', 'analytics'] },
  { id: 'retention_specialist', name: 'Retention Specialist', desc: 'Seller churn mitigation, reactivation & loyalty offers', dept: 'VIP & Growth', defaultPerms: ['dashboard', 'sellers', 'analytics', 'premium_users'] },

  // Dept 5: Marketplace Operations & KYC Verification
  { id: 'operations_admin', name: 'Operations Administrator', desc: 'Daily marketplace flow, seller approvals & support', dept: 'Marketplace Operations', defaultPerms: ['dashboard', 'operations', 'seller_verification', 'refunds', 'appeals', 'pending', 'products', 'sellers', 'seller_upgrades', 'buyer_upgrades', 'transactions'] },
  { id: 'kyc_verification_officer', name: 'KYC Verification Officer', desc: 'Fayda ID, trade license & TIN document verification', dept: 'Marketplace Operations', defaultPerms: ['dashboard', 'seller_verification', 'sellers', 'trust_safety'] },
  { id: 'product_inspection_officer', name: 'Product Inspection Officer', desc: 'Quality assurance, duplicate detection & item inspection', dept: 'Marketplace Operations', defaultPerms: ['dashboard', 'pending', 'products'] },
  { id: 'dispute_refund_officer', name: 'Dispute & Refund Officer', desc: 'Buyer protection claims, order mediation & refunds', dept: 'Marketplace Operations', defaultPerms: ['dashboard', 'refunds', 'transactions', 'appeals'] },
  { id: 'appeals_officer', name: 'Appeals Hearing Officer', desc: 'Sanctions review, suspension appeals & due process', dept: 'Marketplace Operations', defaultPerms: ['dashboard', 'appeals', 'trust_safety', 'sellers'] },

  // Dept 6: Catalog, Taxonomy, Logistics & Support
  { id: 'platform_admin', name: 'Platform Administrator', desc: 'Infrastructure telemetry, server health & system config', dept: 'Catalog & Infrastructure', defaultPerms: ['dashboard', 'settings', 'categories', 'analytics', 'locations', 'banners', 'notifications', 'profit', 'finance', 'sellers', 'transactions', 'operations'] },
  { id: 'taxonomy_architect', name: 'Taxonomy Architect', desc: 'Category hierarchy, attribute schemas & catalog tags', dept: 'Catalog & Taxonomy', defaultPerms: ['dashboard', 'categories', 'analytics'] },
  { id: 'gis_location_specialist', name: 'GIS Location Specialist', desc: 'Ethiopian administrative zones, woredas & GPS nodes', dept: 'Logistics & GIS', defaultPerms: ['dashboard', 'locations', 'analytics'] },
  { id: 'content_moderator', name: 'Content & Media Moderator', desc: 'Image screening, prohibited item takedowns & text hygiene', dept: 'Moderation & Safety', defaultPerms: ['dashboard', 'pending', 'products', 'trust_safety', 'banners'] },
  { id: 'moderator', name: 'Listing Moderator', desc: 'Queue moderation & policy enforcement', dept: 'Moderation & Safety', defaultPerms: ['dashboard', 'pending', 'products', 'trust_safety', 'banners'] },
  { id: 'customer_support_lead', name: 'Customer Support Lead', desc: 'Support ticket routing, agent escalation & resolution', dept: 'Support & Helpdesk', defaultPerms: ['dashboard', 'sellers', 'seller_upgrades', 'buyer_upgrades', 'refunds', 'appeals', 'products'] },
  { id: 'support_agent', name: 'Support Agent', desc: 'Tier-1 customer assistance, seller inquiry routing', dept: 'Support & Helpdesk', defaultPerms: ['dashboard', 'sellers', 'seller_upgrades', 'buyer_upgrades', 'products'] },
];

const AVAILABLE_PERMISSIONS = [
  { id: 'pending',             name: 'Listing Moderation',       desc: 'Review and approve/reject pending listings' },
  { id: 'products',            name: 'All Products',             desc: 'View and manage all marketplace listings' },
  { id: 'sellers',             name: 'Seller Management',        desc: 'View sellers, activate/deactivate accounts' },
  { id: 'analytics',           name: 'Analytics',                desc: 'View platform analytics and charts' },
  { id: 'settings',            name: 'System Settings',          desc: 'Modify platform settings and fee rates' },
  { id: 'transactions',        name: 'Transactions',             desc: 'View and update transaction status' },
  { id: 'notifications',       name: 'Notifications',            desc: 'View and manage admin notifications' },
  { id: 'categories',          name: 'Categories',               desc: 'Manage category tree and attributes' },
  { id: 'seller_verification', name: 'Seller Verification',      desc: 'Review identity and business documents' },
  { id: 'trust_safety',        name: 'Trust & Safety',           desc: 'Manage abuse reports and enforcement' },
  { id: 'refunds',             name: 'Refund Management',        desc: 'Review and process refund requests' },
  { id: 'appeals',             name: 'Account Appeals',          desc: 'Review ban/suspension appeals' },
  { id: 'operations',          name: 'Operations Monitor',       desc: 'View background jobs and system health' },
  { id: 'seller_upgrades',     name: 'Seller Upgrades',          desc: 'Approve/reject seller tier upgrade requests' },
  { id: 'buyer_upgrades',      name: 'Buyer Upgrades',           desc: 'Approve/reject buyer tier upgrade requests' },
  { id: 'locations',           name: 'Locations & GIS',          desc: 'Manage administrative locations and GPS boundaries' },
  { id: 'profit',              name: 'Profit Analytics',         desc: 'View financial planning and historical profit' },
  { id: 'finance',             name: 'Financial Ledger',         desc: 'Manage financial accounts and treasury' },
  { id: 'referrals',           name: 'Referral System',          desc: 'Manage referral commissions and payouts' },
  { id: 'logs',                name: 'Security Logs',            desc: 'Inspect security audit logs and events' },
  { id: 'premium_users',       name: 'Premium Subscribers',      desc: 'Manage VIP and premium memberships' },
  { id: 'banners',             name: 'Marketing Banners',        desc: 'Manage platform promotional banners' },
];

const DEPARTMENTS = [
  'All',
  'Executive Governance',
  'Cybersecurity & SOC',
  'Finance & Treasury',
  'VIP & Growth',
  'Marketplace Operations',
  'Catalog & Infrastructure',
  'Support & Helpdesk',
];

export function AdminRoleManagerModal({ isOpen, onClose, onSuccess }: AdminRoleManagerModalProps) {
  const { user } = useAuth();
  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [selectedDept, setSelectedDept] = useState<string>('All');
  const [adminRole, setAdminRole] = useState<string>('operations_admin');
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([
    'dashboard', 'operations', 'seller_verification', 'refunds', 'appeals', 'pending', 'products', 'sellers', 'seller_upgrades', 'buyer_upgrades', 'transactions'
  ]);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleRoleSelect = (role: RoleOption) => {
    setAdminRole(role.id);
    setSelectedPermissions(role.defaultPerms);
  };

  const togglePermission = (id: string) => {
    setSelectedPermissions((prev) =>
      prev.includes(id) ? prev.filter((p) => p !== id) : [...prev, id]
    );
  };

  const filteredRoles = selectedDept === 'All'
    ? ENTERPRISE_ROLES
    : ENTERPRISE_ROLES.filter(r => r.dept.toLowerCase().includes(selectedDept.toLowerCase()) || (selectedDept.includes('Catalog') && r.dept.includes('Catalog')));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user?.token) return;
    setError(null);
    setSubmitting(true);

    try {
      const created = await api.adminCreateAdmin(
        { name, phone, email, password, admin_role: adminRole },
        user.token
      );

      // Apply the selected RBAC permissions
      if (adminRole !== 'super' && selectedPermissions.length > 0 && created?.id) {
        try {
          await api.adminSetAdminPermissions(created.id, selectedPermissions, user.token);
        } catch {
          // Non-fatal — default permissions are already set by the backend
        }
      }

      onSuccess();
      onClose();
    } catch (err: any) {
      setError(err?.message || 'Failed to create administrative account.');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-md transition-opacity">
      <div className="w-full max-w-3xl bg-zinc-950 border border-zinc-800 rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[92vh] animate-in zoom-in-95">
        
        {/* Header */}
        <div className="p-5 bg-zinc-900 border-b border-zinc-800 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-2.5 bg-orange-600/10 text-orange-500 rounded-xl border border-orange-500/20">
              <UserPlus size={20} />
            </div>
            <div>
              <h2 className="text-base font-extrabold text-white">Provision Enterprise Administrative Account</h2>
              <p className="text-xs text-zinc-400">Deploy specialized personas across 6 enterprise governance departments</p>
            </div>
          </div>

          <button onClick={onClose} className="text-zinc-400 hover:text-white p-2 rounded-xl hover:bg-zinc-800">
            <X size={20} />
          </button>
        </div>

        {error && (
          <div className="m-5 p-3 bg-red-500/10 border border-red-500/30 text-red-400 rounded-xl text-xs font-bold flex items-center gap-2">
            <ShieldAlert size={16} />
            {error}
          </div>
        )}

        {/* Form Body */}
        <form onSubmit={handleSubmit} className="p-6 flex-1 overflow-y-auto space-y-5">
          
          {/* Department Filter Tabs */}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <label className="block text-xs font-bold text-zinc-300 uppercase tracking-wider">
                1. Select Governance Department:
              </label>
              <span className="text-[11px] text-zinc-500 font-mono">
                {filteredRoles.length} Roles in Scope
              </span>
            </div>

            <div className="flex flex-wrap gap-1.5 p-1 bg-zinc-900/80 border border-zinc-800/80 rounded-xl">
              {DEPARTMENTS.map((dept) => (
                <button
                  key={dept}
                  type="button"
                  onClick={() => setSelectedDept(dept)}
                  className={`px-3 py-1.5 rounded-lg text-xs font-bold transition-all ${
                    selectedDept === dept
                      ? 'bg-orange-600 text-white shadow-md'
                      : 'text-zinc-400 hover:text-zinc-200 hover:bg-zinc-800/60'
                  }`}
                >
                  {dept}
                </button>
              ))}
            </div>
          </div>

          {/* Role Grid */}
          <div className="space-y-2">
            <label className="block text-xs font-bold text-zinc-300 uppercase tracking-wider">
              2. Select Specialized Role Persona:
            </label>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-2 max-h-56 overflow-y-auto pr-1">
              {filteredRoles.map((role) => {
                const isSelected = adminRole === role.id;
                return (
                  <button
                    key={role.id}
                    type="button"
                    onClick={() => handleRoleSelect(role)}
                    className={`p-3 rounded-xl border text-left transition-all ${
                      isSelected
                        ? 'bg-orange-600/15 border-orange-500 text-white shadow-md'
                        : 'bg-zinc-900/70 border-zinc-800/80 text-zinc-400 hover:bg-zinc-850 hover:text-white'
                    }`}
                  >
                    <div className="flex items-center gap-1.5 font-bold text-xs text-white">
                      <Shield size={14} className={isSelected ? 'text-orange-400' : 'text-zinc-500'} />
                      <span className="truncate">{role.name}</span>
                    </div>
                    <div className="text-[10px] text-orange-400/80 font-medium mt-0.5">{role.dept}</div>
                    <p className="text-[10px] text-zinc-400 mt-1 line-clamp-2 leading-snug">{role.desc}</p>
                  </button>
                );
              })}
            </div>
            <p className="text-[10px] text-zinc-500 italic mt-1">
              * Note: Platform Owner & Super Admin are Tier-0 root governance and can only be provisioned via the secure offline bootstrap utility.
            </p>
          </div>

          {/* Account Details */}
          <div className="space-y-3 pt-1">
            <label className="block text-xs font-bold text-zinc-300 uppercase tracking-wider">
              3. Administrative Credentials:
            </label>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-semibold text-zinc-400 mb-1">Full Staff Name:</label>
                <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="e.g. Dawit Haile"
                  className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-3 text-sm text-white focus:outline-none focus:border-orange-500"
                  required
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-zinc-400 mb-1">Phone Number:</label>
                <input
                  type="text"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  placeholder="0911234567"
                  className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-3 text-sm text-white focus:outline-none focus:border-orange-500"
                  required
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-semibold text-zinc-400 mb-1">Staff Email Address:</label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="staff@huntiq.et"
                  className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-3 text-sm text-white focus:outline-none focus:border-orange-500"
                  required
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-zinc-400 mb-1">Access Password:</label>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••••••"
                  className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-3 text-sm text-white focus:outline-none focus:border-orange-500"
                  required
                />
              </div>
            </div>
          </div>

          {/* RBAC Permissions Grid */}
          <div className="space-y-3 pt-2">
            <div className="flex items-center justify-between">
              <label className="block text-xs font-bold text-zinc-300 uppercase tracking-wider">
                4. Granular RBAC Permissions ({selectedPermissions.length} Active):
              </label>
              <button
                type="button"
                onClick={() => {
                  const role = ENTERPRISE_ROLES.find(r => r.id === adminRole);
                  if (role) setSelectedPermissions(role.defaultPerms);
                }}
                className="text-[11px] text-orange-400 hover:text-orange-300 font-bold underline"
              >
                Reset to Role Defaults
              </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-2 max-h-48 overflow-y-auto pr-1">
              {AVAILABLE_PERMISSIONS.map((perm) => {
                const isSelected = selectedPermissions.includes(perm.id);
                return (
                  <div
                    key={perm.id}
                    onClick={() => togglePermission(perm.id)}
                    className={`p-2.5 rounded-xl border cursor-pointer transition-all flex items-start gap-2.5 ${
                      isSelected
                        ? 'bg-orange-950/20 border-orange-500/40 text-white'
                        : 'bg-zinc-900/60 border-zinc-800 text-zinc-400 hover:bg-zinc-800'
                    }`}
                  >
                    <div className="mt-0.5 text-orange-400 shrink-0">
                      {isSelected ? <CheckSquare size={15} /> : <Square size={15} />}
                    </div>
                    <div>
                      <div className="text-xs font-bold text-zinc-200">{perm.name}</div>
                      <div className="text-[10px] text-zinc-400 leading-tight mt-0.5">{perm.desc}</div>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          {/* Footer Submit */}
          <div className="pt-4 border-t border-zinc-800 flex items-center justify-end gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2.5 bg-zinc-900 hover:bg-zinc-800 text-zinc-300 rounded-xl text-xs font-bold"
            >
              Cancel
            </button>

            <button
              type="submit"
              disabled={submitting}
              className="px-6 py-2.5 bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-500 hover:to-amber-500 text-white rounded-xl text-xs font-black shadow-lg shadow-orange-600/20 transition-all flex items-center gap-1.5"
            >
              {submitting ? 'Creating Account...' : 'Provision Admin Account'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
