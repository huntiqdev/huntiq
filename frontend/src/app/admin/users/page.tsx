'use client';

import { useEffect, useState, useCallback } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { Users, ShieldAlert, Check, X, Shield, RefreshCw, UserPlus } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { AdminRoleManagerModal } from '@/components/admin/enterprise/AdminRoleManagerModal';

interface RoleMeta {
  name: string;
  dept: string;
  badgeColor: string;
}

const ROLE_DICTIONARY: Record<string, RoleMeta> = {
  // Department 1: Executive Governance
  super: { name: 'Super Administrator', dept: 'Executive Governance', badgeColor: 'bg-red-500/15 text-red-400 border-red-500/30' },
  platform_owner: { name: 'Platform Owner (CEO)', dept: 'Executive Governance', badgeColor: 'bg-red-500/15 text-red-400 border-red-500/30' },
  gov_compliance_officer: { name: 'Gov Compliance Officer', dept: 'Executive Governance', badgeColor: 'bg-rose-500/15 text-rose-400 border-rose-500/30' },
  legal_counsel: { name: 'Legal Counsel', dept: 'Executive Governance', badgeColor: 'bg-rose-500/15 text-rose-400 border-rose-500/30' },
  hr_director: { name: 'HR & Staff Director', dept: 'Executive Governance', badgeColor: 'bg-rose-500/15 text-rose-400 border-rose-500/30' },

  // Department 2: Cybersecurity & SOC Defense
  ciso: { name: 'Chief Info Security Officer (CISO)', dept: 'Cybersecurity & SOC', badgeColor: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30' },
  security_admin: { name: 'Security Administrator', dept: 'Cybersecurity & SOC', badgeColor: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30' },
  soc_analyst: { name: 'SOC Threat Analyst', dept: 'Cybersecurity & SOC', badgeColor: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30' },
  fraud_investigator: { name: 'Fraud & AML Investigator', dept: 'Cybersecurity & SOC', badgeColor: 'bg-teal-500/15 text-teal-400 border-teal-500/30' },
  compliance_auditor: { name: 'Compliance Auditor', dept: 'Cybersecurity & SOC', badgeColor: 'bg-teal-500/15 text-teal-400 border-teal-500/30' },

  // Department 3: Banking, Treasury & Financial Operations
  finance_director: { name: 'Finance & Treasury Director', dept: 'Finance & Treasury', badgeColor: 'bg-cyan-500/15 text-cyan-400 border-cyan-500/30' },
  transaction_auditor: { name: 'Transaction & Ledger Auditor', dept: 'Finance & Treasury', badgeColor: 'bg-cyan-500/15 text-cyan-400 border-cyan-500/30' },
  tax_controller: { name: 'Tax & NBE Controller', dept: 'Finance & Treasury', badgeColor: 'bg-blue-500/15 text-blue-400 border-blue-500/30' },
  escrow_reconciler: { name: 'Escrow Reconciler', dept: 'Finance & Treasury', badgeColor: 'bg-blue-500/15 text-blue-400 border-blue-500/30' },
  finance_support_agent: { name: 'Finance Support Agent', dept: 'Finance & Treasury', badgeColor: 'bg-sky-500/15 text-sky-400 border-sky-500/30' },

  // Department 4: VIP Concierge & Growth Services
  vip_account_director: { name: 'VIP Account Director', dept: 'VIP & Growth', badgeColor: 'bg-amber-500/15 text-amber-400 border-amber-500/30' },
  gold_concierge: { name: 'Gold Concierge', dept: 'VIP & Growth', badgeColor: 'bg-amber-500/15 text-amber-400 border-amber-500/30' },
  premium_concierge: { name: 'Premium Concierge', dept: 'VIP & Growth', badgeColor: 'bg-amber-500/15 text-amber-400 border-amber-500/30' },
  referral_director: { name: 'Referral Director', dept: 'VIP & Growth', badgeColor: 'bg-yellow-500/15 text-yellow-400 border-yellow-500/30' },
  retention_specialist: { name: 'Retention Specialist', dept: 'VIP & Growth', badgeColor: 'bg-yellow-500/15 text-yellow-400 border-yellow-500/30' },

  // Department 5: Marketplace Operations & KYC Verification
  operations_admin: { name: 'Operations Administrator', dept: 'Marketplace Operations', badgeColor: 'bg-orange-500/15 text-orange-400 border-orange-500/30' },
  kyc_verification_officer: { name: 'KYC Verification Officer', dept: 'Marketplace Operations', badgeColor: 'bg-orange-500/15 text-orange-400 border-orange-500/30' },
  product_inspection_officer: { name: 'Product Inspection Officer', dept: 'Marketplace Operations', badgeColor: 'bg-orange-500/15 text-orange-400 border-orange-500/30' },
  dispute_refund_officer: { name: 'Dispute & Refund Officer', dept: 'Marketplace Operations', badgeColor: 'bg-orange-500/15 text-orange-400 border-orange-500/30' },
  appeals_officer: { name: 'Appeals Hearing Officer', dept: 'Marketplace Operations', badgeColor: 'bg-orange-500/15 text-orange-400 border-orange-500/30' },

  // Department 6: Catalog, Taxonomy, Logistics & Support
  platform_admin: { name: 'Platform Administrator', dept: 'Catalog & Infrastructure', badgeColor: 'bg-indigo-500/15 text-indigo-400 border-indigo-500/30' },
  taxonomy_architect: { name: 'Taxonomy Architect', dept: 'Catalog & Taxonomy', badgeColor: 'bg-indigo-500/15 text-indigo-400 border-indigo-500/30' },
  gis_location_specialist: { name: 'GIS Location Specialist', dept: 'Logistics & GIS', badgeColor: 'bg-indigo-500/15 text-indigo-400 border-indigo-500/30' },
  content_moderator: { name: 'Content & Media Moderator', dept: 'Moderation & Safety', badgeColor: 'bg-blue-500/15 text-blue-400 border-blue-500/30' },
  moderator: { name: 'Content Moderator', dept: 'Moderation & Safety', badgeColor: 'bg-blue-500/15 text-blue-400 border-blue-500/30' },
  customer_support_lead: { name: 'Customer Support Lead', dept: 'Support & Helpdesk', badgeColor: 'bg-purple-500/15 text-purple-400 border-purple-500/30' },
  support_agent: { name: 'Support Agent', dept: 'Support & Helpdesk', badgeColor: 'bg-purple-500/15 text-purple-400 border-purple-500/30' },
  normal: { name: 'Standard Verified Staff', dept: 'General Staff', badgeColor: 'bg-zinc-800 text-zinc-400 border-zinc-700' },
};

export default function AdminUsersPage() {
  const { user, isSuperAdmin, isPlatformOwner, loading: authLoading } = useAuth();
  const canManageAdmins = isPlatformOwner || isSuperAdmin;
  const router = useRouter();

  const [admins, setAdmins] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const loadAdmins = useCallback(async () => {
    if (!user) return;
    setLoading(true);
    try {
      const res = await api.adminListAdmins(user.token);
      setAdmins(res || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [user]);

  useEffect(() => {
    if (!authLoading && !canManageAdmins) {
      router.push('/admin/dashboard');
      return;
    }
    if (canManageAdmins) {
      loadAdmins();
    }
  }, [canManageAdmins, authLoading, router, loadAdmins]);

  const handleDelete = async (id: number) => {
    if (!user || !confirm('Are you sure you want to revoke access for this admin account?')) return;
    try {
      await api.adminDeleteAdmin(id, user.token);
      setAdmins((prev) => prev.filter((a) => a.id !== id));
    } catch (err) {
      alert('Failed to delete admin account');
    }
  };

  if (authLoading) return <div className="h-64 flex items-center justify-center"><Spinner /></div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-red-600 rounded-2xl flex items-center justify-center text-white shadow-lg shadow-red-600/20">
            <Users size={24} />
          </div>
          <div>
            <h1 className="text-2xl font-black text-white leading-none">Administrative Accounts & RBAC</h1>
            <p className="text-xs font-bold text-zinc-400 uppercase tracking-widest mt-1">Manage Sub-Admins, Super-Admins & Roles</p>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <button
            onClick={() => setIsModalOpen(true)}
            className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-500 hover:to-amber-500 text-white font-extrabold text-xs rounded-xl shadow-lg shadow-orange-600/20 transition-all"
          >
            <UserPlus size={16} />
            <span>Create Admin Account</span>
          </button>

          <button
            onClick={loadAdmins}
            className="p-2 rounded-xl bg-zinc-900 text-zinc-300 hover:bg-zinc-800 border border-zinc-800 transition-colors"
          >
            <RefreshCw size={20} className={loading ? 'animate-spin' : ''} />
          </button>
        </div>
      </div>

      <div className="bg-zinc-900 border border-zinc-800 rounded-3xl overflow-hidden shadow-xl">
        <div className="p-5 border-b border-zinc-800 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <ShieldAlert size={18} className="text-red-500" />
            <h2 className="font-extrabold text-sm text-white">Active System Administrator Directory</h2>
          </div>
          <span className="text-xs font-mono text-zinc-400">{admins.length} Total Accounts</span>
        </div>
        
        {loading ? (
          <div className="h-40 flex items-center justify-center"><Spinner /></div>
        ) : admins.length === 0 ? (
          <div className="p-8 text-center text-zinc-500 font-bold">No admin accounts registered yet.</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs whitespace-nowrap">
              <thead className="bg-zinc-950 text-zinc-400 font-mono text-[11px] uppercase tracking-widest border-b border-zinc-800">
                <tr>
                  <th className="px-6 py-4">ID</th>
                  <th className="px-6 py-4">Admin Name</th>
                  <th className="px-6 py-4">Email / Phone Contact</th>
                  <th className="px-6 py-4">Role Classification</th>
                  <th className="px-6 py-4 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-800/80">
                {admins.map((a) => (
                  <tr key={a.id} className="hover:bg-zinc-800/50 transition-colors">
                    <td className="px-6 py-4 font-mono text-zinc-400">#ADM-{a.id}</td>
                    <td className="px-6 py-4 font-bold text-white">{a.name}</td>
                    <td className="px-6 py-4 text-zinc-300">
                      <div>{a.email}</div>
                      <div className="text-[11px] font-mono text-zinc-500">{a.phone}</div>
                    </td>
                    <td className="px-6 py-4">
                      {(() => {
                        const roleKey = a.admin_role || 'operations_admin';
                        const meta = ROLE_DICTIONARY[roleKey] || {
                          name: roleKey.replace(/_/g, ' ').replace(/\b\w/g, (c: string) => c.toUpperCase()),
                          dept: 'General Operations',
                          badgeColor: 'bg-zinc-800 text-zinc-300 border-zinc-700',
                        };

                        return (
                          <div className="space-y-1">
                            <span className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider border ${meta.badgeColor}`}>
                              <Shield size={11} />
                              {meta.name}
                            </span>
                            <div className="text-[10px] font-medium text-zinc-500">
                              {meta.dept}
                            </div>
                          </div>
                        );
                      })()}
                    </td>
                    <td className="px-6 py-4 text-right">
                      {user?.id !== a.id ? (
                        <button
                          onClick={() => handleDelete(a.id)}
                          className="text-red-400 hover:text-red-300 font-bold text-xs p-2 hover:bg-red-500/10 rounded-lg transition-colors border border-red-500/20"
                        >
                          Revoke Access
                        </button>
                      ) : (
                        <span className="text-xs font-bold text-emerald-400 px-3 py-1 bg-emerald-500/10 rounded-lg border border-emerald-500/20">
                          Current Active Session
                        </span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Role Manager Modal */}
      <AdminRoleManagerModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSuccess={() => {
          loadAdmins();
        }}
      />
    </div>
  );
}
