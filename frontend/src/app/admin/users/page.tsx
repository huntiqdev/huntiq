'use client';

import { useEffect, useState, useCallback } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { Users, ShieldAlert, Check, X, Shield, RefreshCw, UserPlus } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { AdminRoleManagerModal } from '@/components/admin/enterprise/AdminRoleManagerModal';

export default function AdminUsersPage() {
  const { user, isSuperAdmin, loading: authLoading } = useAuth();
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
    if (!authLoading && !isSuperAdmin) {
      router.push('/admin/dashboard');
      return;
    }
    if (isSuperAdmin) {
      loadAdmins();
    }
  }, [isSuperAdmin, authLoading, router, loadAdmins]);

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
                      <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-wider ${
                        a.admin_role === 'super'
                          ? 'bg-red-500/15 text-red-400 border border-red-500/30'
                          : 'bg-blue-500/15 text-blue-400 border border-blue-500/30'
                      }`}>
                        <Shield size={12} />
                        {a.admin_role === 'super' ? 'Super Admin' : 'Sub-Admin Manager'}
                      </span>
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
