'use client';

import React, { useState } from 'react';
import { X, Shield, ShieldAlert, Check, UserPlus, Key, Mail, Phone, Lock, Sparkles, CheckSquare, Square } from 'lucide-react';
import { api } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';

interface AdminRoleManagerModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

const AVAILABLE_PERMISSIONS = [
  { id: 'pending',           name: 'Listing Moderation',       desc: 'Review and approve/reject pending listings' },
  { id: 'products',          name: 'All Products',             desc: 'View and manage all marketplace listings' },
  { id: 'sellers',           name: 'Seller Management',        desc: 'View sellers, activate/deactivate accounts' },
  { id: 'analytics',         name: 'Analytics',                desc: 'View platform analytics and charts' },
  { id: 'settings',          name: 'System Settings',          desc: 'Modify platform settings and fee rates' },
  { id: 'transactions',      name: 'Transactions',             desc: 'View and update transaction status' },
  { id: 'notifications',     name: 'Notifications',            desc: 'View and manage admin notifications' },
  { id: 'categories',        name: 'Categories',               desc: 'Manage category tree and attributes' },
  { id: 'seller_verification', name: 'Seller Verification',    desc: 'Review identity and business documents' },
  { id: 'trust_safety',      name: 'Trust & Safety',           desc: 'Manage abuse reports and enforcement' },
  { id: 'refunds',           name: 'Refund Management',        desc: 'Review and process refund requests' },
  { id: 'appeals',           name: 'Account Appeals',          desc: 'Review ban/suspension appeals' },
  { id: 'operations',        name: 'Operations Monitor',       desc: 'View background jobs and system health' },
  { id: 'seller_upgrades',   name: 'Seller Upgrades',          desc: 'Approve/reject seller tier upgrade requests' },
  { id: 'buyer_upgrades',    name: 'Buyer Upgrades',           desc: 'Approve/reject buyer tier upgrade requests' },
];

export function AdminRoleManagerModal({ isOpen, onClose, onSuccess }: AdminRoleManagerModalProps) {
  const { user } = useAuth();
  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [adminRole, setAdminRole] = useState<string>('operations_admin');
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>(['products', 'analytics']);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const togglePermission = (id: string) => {
    setSelectedPermissions((prev) =>
      prev.includes(id) ? prev.filter((p) => p !== id) : [...prev, id]
    );
  };

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

      // Apply the selected RBAC permissions (only for sub-admins — super admins have all)
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
      <div className="w-full max-w-2xl bg-zinc-950 border border-zinc-800 rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh] animate-in zoom-in-95">
        
        {/* Header */}
        <div className="p-5 bg-zinc-900 border-b border-zinc-800 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-2.5 bg-orange-600/10 text-orange-500 rounded-xl border border-orange-500/20">
              <UserPlus size={20} />
            </div>
            <div>
              <h2 className="text-base font-extrabold text-white">Create Administrative Account</h2>
              <p className="text-xs text-zinc-400">Configure Sub-Admin / Super-Admin roles & granular RBAC permissions</p>
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
          
          {/* Role Selection Tabs */}
          <div className="space-y-2">
            <label className="block text-xs font-bold text-zinc-300 uppercase tracking-wider">Enterprise Administrative Role:</label>
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-2.5">
              {[
                { id: 'platform_admin', name: 'Platform Admin', desc: 'System configuration & platform operations', color: 'indigo' },
                { id: 'security_admin', name: 'Security Admin', desc: 'Audit monitoring & incident investigation', color: 'emerald' },
                { id: 'operations_admin', name: 'Operations Admin', desc: 'Routine operational management', color: 'amber' },
                { id: 'moderator', name: 'Content Moderator', desc: 'Listing reviews & report handling', color: 'blue' },
                { id: 'support_agent', name: 'Support Agent', desc: 'Customer support & assistance', color: 'purple' },
              ].map((role) => (
                <button
                  key={role.id}
                  type="button"
                  onClick={() => setAdminRole(role.id as any)}
                  className={`p-3 rounded-xl border text-left transition-all ${
                    adminRole === role.id
                      ? 'bg-orange-600/15 border-orange-500 text-white shadow-md'
                      : 'bg-zinc-900 border-zinc-800 text-zinc-400 hover:bg-zinc-850 hover:text-white'
                  }`}
                >
                  <div className="flex items-center gap-1.5 font-bold text-xs text-white">
                    <Shield size={14} className={adminRole === role.id ? 'text-orange-400' : 'text-zinc-500'} />
                    <span>{role.name}</span>
                  </div>
                  <p className="text-[10px] text-zinc-400 mt-1 leading-snug">{role.desc}</p>
                </button>
              ))}
            </div>
            <p className="text-[10px] text-zinc-500 italic mt-1">
              * Note: Platform Owner is Tier-0 root governance and can only be provisioned via the secure offline bootstrap utility.
            </p>
          </div>

          {/* Account Details */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-bold text-zinc-300 mb-1">Full Name:</label>
              <div className="relative">
                <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="e.g. Solomon Tadesse"
                  className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-3 text-sm text-white focus:outline-none focus:border-orange-500"
                  required
                />
              </div>
            </div>

            <div>
              <label className="block text-xs font-bold text-zinc-300 mb-1">Phone Number:</label>
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

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-bold text-zinc-300 mb-1">Email Address:</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="admin@huntiq.et"
                className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-3 text-sm text-white focus:outline-none focus:border-orange-500"
                required
              />
            </div>

            <div>
              <label className="block text-xs font-bold text-zinc-300 mb-1">Access Password:</label>
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

          {/* RBAC Permissions Grid */}
          <div className="space-y-3 pt-2">
            <label className="block text-xs font-bold text-zinc-300 uppercase tracking-wider">
              Granular Role-Based Access Permissions (RBAC):
            </label>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-2.5">
              {AVAILABLE_PERMISSIONS.map((perm) => {
                const isSelected = selectedPermissions.includes(perm.id) || adminRole === 'super';
                return (
                  <div
                    key={perm.id}
                    onClick={() => adminRole !== 'super' && togglePermission(perm.id)}
                    className={`p-3 rounded-xl border cursor-pointer transition-all flex items-start gap-3 ${
                      isSelected
                        ? 'bg-orange-950/20 border-orange-500/40 text-white'
                        : 'bg-zinc-900/60 border-zinc-800 text-zinc-400 hover:bg-zinc-800'
                    }`}
                  >
                    <div className="mt-0.5 text-orange-400 shrink-0">
                      {isSelected ? <CheckSquare size={16} /> : <Square size={16} />}
                    </div>
                    <div>
                      <div className="text-xs font-bold text-zinc-200">{perm.name}</div>
                      <div className="text-[11px] text-zinc-400 leading-tight mt-0.5">{perm.desc}</div>
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
