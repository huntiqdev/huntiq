'use client';
import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api, PaymentOrder, PaymentProviderConfig } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import {
  CreditCard,
  ShieldAlert,
  ShieldCheck,
  CheckCircle,
  XCircle,
  AlertTriangle,
  RefreshCw,
  Search,
  Filter,
  ExternalLink,
  Settings,
  Cpu,
  Smartphone,
  Phone,
  Clock,
  Eye,
} from 'lucide-react';

export default function AdminPaymentsPage() {
  const { user, isPlatformOwner } = useAuth();
  const [activeTab, setActiveTab] = useState<'orders' | 'providers'>('orders');

  // Operating mode & providers state
  const [operatingMode, setOperatingMode] = useState<'demo' | 'production'>('demo');
  const [providers, setProviders] = useState<PaymentProviderConfig[]>([]);
  const [switchingMode, setSwitchingMode] = useState(false);
  const [modeError, setModeError] = useState('');
  const [modeSuccess, setModeSuccess] = useState('');

  // Orders state
  const [orders, setOrders] = useState<PaymentOrder[]>([]);
  const [totalOrders, setTotalOrders] = useState(0);
  const [loadingOrders, setLoadingOrders] = useState(false);
  const [statusFilter, setStatusFilter] = useState('');
  const [providerFilter, setProviderFilter] = useState('');
  const [envFilter, setEnvFilter] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [page, setPage] = useState(1);
  const limit = 15;

  // Selected order modal
  const [selectedOrder, setSelectedOrder] = useState<PaymentOrder | null>(null);

  // Edit provider modal
  const [editingProvider, setEditingProvider] = useState<PaymentProviderConfig | null>(null);
  const [editEndpoint, setEditEndpoint] = useState('');
  const [editAppID, setEditAppID] = useState('');
  const [editAppKey, setEditAppKey] = useState('');
  const [editShortCode, setEditShortCode] = useState('');
  const [editEnabled, setEditEnabled] = useState(true);
  const [savingProvider, setSavingProvider] = useState(false);

  // Load providers and operating mode
  const loadProviders = useCallback(async () => {
    if (!user?.token) return;
    try {
      const res = await api.adminGetPaymentProviders(user.token);
      if (res) {
        setProviders(res.providers || []);
        setOperatingMode(res.operating_mode || 'demo');
      }
    } catch {
      // ignore
    }
  }, [user?.token]);

  // Load orders
  const loadOrders = useCallback(async () => {
    if (!user?.token) return;
    setLoadingOrders(true);
    try {
      const res = await api.adminListPaymentOrders(
        {
          status: statusFilter,
          provider: providerFilter,
          environment: envFilter,
          search: searchQuery,
          limit,
          offset: (page - 1) * limit,
        },
        user.token
      );
      if (res) {
        setOrders(res.orders || []);
        setTotalOrders(res.total || 0);
      }
    } catch {
      // ignore
    } finally {
      setLoadingOrders(false);
    }
  }, [user?.token, statusFilter, providerFilter, envFilter, searchQuery, page]);

  useEffect(() => {
    loadProviders();
  }, [loadProviders]);

  useEffect(() => {
    loadOrders();
  }, [loadOrders]);

  const handleSwitchMode = async (targetMode: 'demo' | 'production') => {
    if (!user?.token || !isPlatformOwner) return;
    setSwitchingMode(true);
    setModeError('');
    setModeSuccess('');
    try {
      const res = await api.adminSwitchPaymentMode(targetMode, user.token);
      setOperatingMode(res.operating_mode);
      setModeSuccess(`Payment operating mode set to ${res.operating_mode.toUpperCase()}`);
      loadProviders();
      loadOrders();
    } catch (err: any) {
      setModeError(err.message ?? 'Mode switch failed');
    } finally {
      setSwitchingMode(false);
    }
  };

  const handleSaveProvider = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user?.token || !editingProvider) return;
    setSavingProvider(true);
    try {
      await api.adminUpdatePaymentProvider(
        editingProvider.provider_code,
        {
          environment_mode: editingProvider.environment_mode,
          is_enabled: editEnabled,
          api_endpoint: editEndpoint,
          app_id: editAppID,
          app_key: editAppKey,
          short_code: editShortCode,
          notify_url: editingProvider.notify_url,
          return_url: editingProvider.return_url,
        },
        user.token
      );
      setEditingProvider(null);
      loadProviders();
    } catch (err: any) {
      alert(err.message ?? 'Failed to update provider');
    } finally {
      setSavingProvider(false);
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'SUCCESSFUL':
        return 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20';
      case 'PENDING':
      case 'PROCESSING':
        return 'bg-amber-500/10 text-amber-400 border-amber-500/20';
      case 'FAILED':
      case 'CANCELLED':
      case 'EXPIRED':
        return 'bg-red-500/10 text-red-400 border-red-500/20';
      default:
        return 'bg-slate-500/10 text-slate-400 border-slate-500/20';
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <CreditCard className="text-blue-400" size={24} />
            <h1 className="text-xl font-black text-white">Payment Gateway Governance</h1>
          </div>
          <p className="text-xs text-slate-400 mt-1">
            Authoritative order verification, provider credentials, and fail-closed demo/production switching.
          </p>
        </div>

        {/* Operating Mode Status Pill */}
        <div className="flex items-center gap-3 bg-slate-900 border border-slate-800 p-2 rounded-2xl">
          <div className="flex items-center gap-2 px-3 py-1 rounded-xl bg-slate-950 border border-slate-800">
            <span className={`w-2 h-2 rounded-full ${operatingMode === 'production' ? 'bg-emerald-400 animate-pulse' : 'bg-amber-400'}`} />
            <span className="text-xs font-mono font-bold text-white uppercase tracking-wider">
              Mode: {operatingMode}
            </span>
          </div>

          {isPlatformOwner && (
            <button
              onClick={() => handleSwitchMode(operatingMode === 'demo' ? 'production' : 'demo')}
              disabled={switchingMode}
              className={`px-3 py-1 rounded-xl text-xs font-bold transition-all ${
                operatingMode === 'demo'
                  ? 'bg-blue-600 hover:bg-blue-500 text-white'
                  : 'bg-amber-600 hover:bg-amber-500 text-white'
              }`}
            >
              {switchingMode ? (
                <RefreshCw size={12} className="animate-spin inline mr-1" />
              ) : null}
              Switch to {operatingMode === 'demo' ? 'PRODUCTION' : 'DEMO'}
            </button>
          )}
        </div>
      </div>

      {modeError && (
        <div className="p-3 bg-red-950/40 border border-red-800 rounded-xl flex items-center gap-2 text-xs text-red-300">
          <AlertTriangle size={14} className="shrink-0 text-red-400" />
          <span>{modeError}</span>
        </div>
      )}

      {modeSuccess && (
        <div className="p-3 bg-emerald-950/40 border border-emerald-800 rounded-xl flex items-center gap-2 text-xs text-emerald-300">
          <CheckCircle size={14} className="shrink-0 text-emerald-400" />
          <span>{modeSuccess}</span>
        </div>
      )}

      {/* Tabs */}
      <div className="flex gap-2 border-b border-slate-800 pb-2">
        <button
          onClick={() => setActiveTab('orders')}
          className={`px-4 py-2 rounded-xl text-xs font-bold transition-all ${
            activeTab === 'orders'
              ? 'bg-blue-600 text-white shadow-xs'
              : 'text-slate-400 hover:text-white hover:bg-slate-900'
          }`}
        >
          Payment Orders Ledger ({totalOrders})
        </button>
        <button
          onClick={() => setActiveTab('providers')}
          className={`px-4 py-2 rounded-xl text-xs font-bold transition-all ${
            activeTab === 'providers'
              ? 'bg-blue-600 text-white shadow-xs'
              : 'text-slate-400 hover:text-white hover:bg-slate-900'
          }`}
        >
          Provider Configurations ({providers.length})
        </button>
      </div>

      {/* Tab: Orders */}
      {activeTab === 'orders' && (
        <div className="space-y-4">
          {/* Filters Bar */}
          <div className="grid grid-cols-1 sm:grid-cols-4 gap-3 bg-slate-900/60 p-3.5 rounded-2xl border border-slate-800">
            <div className="relative">
              <Search size={14} className="absolute left-3 top-3 text-slate-500" />
              <input
                type="text"
                value={searchQuery}
                onChange={e => { setSearchQuery(e.target.value); setPage(1); }}
                placeholder="Search Order ID, user, phone..."
                className="w-full pl-8 pr-3 py-1.5 text-xs rounded-xl bg-slate-950 border border-slate-800 text-white placeholder-slate-500 focus:outline-none focus:border-blue-500"
              />
            </div>

            <select
              value={statusFilter}
              onChange={e => { setStatusFilter(e.target.value); setPage(1); }}
              className="py-1.5 px-3 text-xs rounded-xl bg-slate-950 border border-slate-800 text-white focus:outline-none focus:border-blue-500"
            >
              <option value="">All Statuses</option>
              <option value="SUCCESSFUL">SUCCESSFUL</option>
              <option value="PENDING">PENDING</option>
              <option value="PROCESSING">PROCESSING</option>
              <option value="FAILED">FAILED</option>
              <option value="CANCELLED">CANCELLED</option>
              <option value="EXPIRED">EXPIRED</option>
            </select>

            <select
              value={providerFilter}
              onChange={e => { setProviderFilter(e.target.value); setPage(1); }}
              className="py-1.5 px-3 text-xs rounded-xl bg-slate-950 border border-slate-800 text-white focus:outline-none focus:border-blue-500"
            >
              <option value="">All Providers</option>
              <option value="demo_mock">Demo Mock Simulator</option>
              <option value="telebirr">Ethio Telecom Telebirr</option>
              <option value="ethio_telecom">Ethio Telecom TelePay</option>
            </select>

            <select
              value={envFilter}
              onChange={e => { setEnvFilter(e.target.value); setPage(1); }}
              className="py-1.5 px-3 text-xs rounded-xl bg-slate-950 border border-slate-800 text-white focus:outline-none focus:border-blue-500"
            >
              <option value="">All Environments</option>
              <option value="demo">Demo Sandbox</option>
              <option value="production">Live Production</option>
            </select>
          </div>

          {/* Orders Table */}
          <div className="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead className="bg-slate-950 text-slate-400 uppercase tracking-wider font-mono text-[10px] border-b border-slate-800">
                  <tr>
                    <th className="p-3">Order ID</th>
                    <th className="p-3">User</th>
                    <th className="p-3">Type</th>
                    <th className="p-3">Amount</th>
                    <th className="p-3">Provider</th>
                    <th className="p-3">Env</th>
                    <th className="p-3">Status</th>
                    <th className="p-3">Date</th>
                    <th className="p-3 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800">
                  {loadingOrders ? (
                    <tr>
                      <td colSpan={9} className="p-8 text-center text-slate-500">
                        <Spinner size={24} />
                      </td>
                    </tr>
                  ) : orders.length === 0 ? (
                    <tr>
                      <td colSpan={9} className="p-8 text-center text-slate-500">
                        No payment orders found matching criteria.
                      </td>
                    </tr>
                  ) : (
                    orders.map(o => (
                      <tr key={o.id} className="hover:bg-slate-800/40 transition-colors">
                        <td className="p-3 font-mono font-bold text-white whitespace-nowrap">
                          {o.order_id}
                        </td>
                        <td className="p-3 whitespace-nowrap">
                          <p className="font-semibold text-slate-200">{o.user_name || `User #${o.user_id}`}</p>
                          <p className="text-[10px] text-slate-400">{o.user_phone || o.payer_phone || '—'}</p>
                        </td>
                        <td className="p-3 whitespace-nowrap">
                          <span className="capitalize text-slate-300 font-medium">{o.order_type.replace('_', ' ')}</span>
                          {o.target_tier && (
                            <span className="ml-1 text-[10px] font-bold text-purple-400 uppercase">
                              ({o.target_tier})
                            </span>
                          )}
                        </td>
                        <td className="p-3 font-mono font-bold text-white whitespace-nowrap">
                          {o.amount.toLocaleString()} {o.currency}
                        </td>
                        <td className="p-3 font-mono text-[11px] text-slate-300 whitespace-nowrap">
                          {o.provider}
                        </td>
                        <td className="p-3 whitespace-nowrap">
                          <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase ${
                            o.environment === 'production'
                              ? 'bg-emerald-950 text-emerald-400 border border-emerald-800'
                              : 'bg-amber-950 text-amber-400 border border-amber-800'
                          }`}>
                            {o.environment}
                          </span>
                        </td>
                        <td className="p-3 whitespace-nowrap">
                          <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold border ${getStatusBadge(o.status)}`}>
                            {o.status}
                          </span>
                        </td>
                        <td className="p-3 text-slate-400 whitespace-nowrap text-[11px]">
                          {new Date(o.created_at).toLocaleDateString()} {new Date(o.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                        </td>
                        <td className="p-3 text-right whitespace-nowrap">
                          <button
                            onClick={() => setSelectedOrder(o)}
                            className="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
                            title="View Details"
                          >
                            <Eye size={14} />
                          </button>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {totalOrders > limit && (
              <div className="p-3 bg-slate-950 border-t border-slate-800 flex items-center justify-between text-xs text-slate-400">
                <span>Total {totalOrders} orders</span>
                <div className="flex gap-2">
                  <button
                    disabled={page <= 1}
                    onClick={() => setPage(p => Math.max(1, p - 1))}
                    className="px-3 py-1 rounded-lg border border-slate-800 disabled:opacity-50 hover:bg-slate-900"
                  >
                    Prev
                  </button>
                  <span className="px-3 py-1 font-mono">{page} / {Math.ceil(totalOrders / limit)}</span>
                  <button
                    disabled={page * limit >= totalOrders}
                    onClick={() => setPage(p => p + 1)}
                    className="px-3 py-1 rounded-lg border border-slate-800 disabled:opacity-50 hover:bg-slate-900"
                  >
                    Next
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Tab: Providers */}
      {activeTab === 'providers' && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {providers.map(p => (
            <div key={p.id} className="bg-slate-900 border border-slate-800 rounded-2xl p-5 space-y-4 flex flex-col justify-between">
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    {p.provider_code === 'telebirr' ? (
                      <Smartphone className="text-blue-400" size={20} />
                    ) : p.provider_code === 'ethio_telecom' ? (
                      <Phone className="text-emerald-400" size={20} />
                    ) : (
                      <Cpu className="text-amber-400" size={20} />
                    )}
                    <div>
                      <h3 className="text-sm font-bold text-white">{p.provider_name}</h3>
                      <p className="text-[10px] font-mono text-slate-500">{p.provider_code}</p>
                    </div>
                  </div>
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold ${
                    p.is_enabled ? 'bg-emerald-500/10 text-emerald-400' : 'bg-red-500/10 text-red-400'
                  }`}>
                    {p.is_enabled ? 'ENABLED' : 'DISABLED'}
                  </span>
                </div>

                <div className="space-y-1.5 text-xs text-slate-400">
                  <div className="flex justify-between border-b border-slate-800/60 pb-1">
                    <span className="text-slate-500">Mode</span>
                    <span className="font-mono text-white capitalize">{p.environment_mode}</span>
                  </div>
                  <div className="flex justify-between border-b border-slate-800/60 pb-1">
                    <span className="text-slate-500">App / Merchant ID</span>
                    <span className="font-mono text-slate-300">{p.app_id || '—'}</span>
                  </div>
                  <div className="flex justify-between border-b border-slate-800/60 pb-1">
                    <span className="text-slate-500">Short Code</span>
                    <span className="font-mono text-slate-300">{p.short_code || '—'}</span>
                  </div>
                  <div className="flex justify-between border-b border-slate-800/60 pb-1">
                    <span className="text-slate-500">App Key</span>
                    <span className="font-mono text-slate-300">{p.app_key_masked || '—'}</span>
                  </div>
                  <div className="pt-1">
                    <span className="text-[10px] text-slate-500 block mb-0.5">Endpoint</span>
                    <p className="text-[10px] font-mono text-slate-400 truncate">{p.api_endpoint || '—'}</p>
                  </div>
                </div>
              </div>

              {isPlatformOwner && (
                <button
                  onClick={() => {
                    setEditingProvider(p);
                    setEditEndpoint(p.api_endpoint || '');
                    setEditAppID(p.app_id || '');
                    setEditAppKey('');
                    setEditShortCode(p.short_code || '');
                    setEditEnabled(p.is_enabled);
                  }}
                  className="w-full bg-slate-800 hover:bg-slate-700 text-white text-xs font-bold py-2 rounded-xl flex items-center justify-center gap-1.5 transition-colors"
                >
                  <Settings size={12} /> Configure Credentials
                </button>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Detail Modal */}
      {selectedOrder && (
        <div className="fixed inset-0 z-50 bg-black/80 flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-lg w-full p-6 space-y-4">
            <div className="flex items-center justify-between">
              <h2 className="text-base font-bold text-white font-mono">{selectedOrder.order_id}</h2>
              <button
                onClick={() => setSelectedOrder(null)}
                className="text-slate-400 hover:text-white text-sm"
              >
                ✕
              </button>
            </div>

            <div className="space-y-2 text-xs">
              <div className="grid grid-cols-2 gap-2 p-3 bg-slate-950 rounded-xl">
                <div>
                  <span className="text-slate-500 block">User</span>
                  <span className="font-bold text-white">{selectedOrder.user_name || `User #${selectedOrder.user_id}`}</span>
                </div>
                <div>
                  <span className="text-slate-500 block">Amount</span>
                  <span className="font-bold text-white">{selectedOrder.amount} {selectedOrder.currency}</span>
                </div>
                <div>
                  <span className="text-slate-500 block">Provider Ref</span>
                  <span className="font-mono text-slate-300">{selectedOrder.provider_tx_id || '—'}</span>
                </div>
                <div>
                  <span className="text-slate-500 block">OTP Verified</span>
                  <span className="font-mono text-emerald-400">{selectedOrder.otp_verified ? 'YES' : 'NO'}</span>
                </div>
              </div>

              {selectedOrder.failure_reason && (
                <div className="p-3 bg-red-950/40 border border-red-800 rounded-xl text-red-300 text-xs">
                  <strong>Failure Reason:</strong> {selectedOrder.failure_reason}
                </div>
              )}
            </div>

            <button
              onClick={() => setSelectedOrder(null)}
              className="w-full bg-slate-800 hover:bg-slate-700 text-white font-bold py-2 rounded-xl text-xs"
            >
              Close
            </button>
          </div>
        </div>
      )}

      {/* Edit Provider Modal */}
      {editingProvider && (
        <div className="fixed inset-0 z-50 bg-black/80 flex items-center justify-center p-4">
          <form onSubmit={handleSaveProvider} className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 space-y-4">
            <div className="flex items-center justify-between">
              <h2 className="text-base font-bold text-white">Configure {editingProvider.provider_name}</h2>
              <button
                type="button"
                onClick={() => setEditingProvider(null)}
                className="text-slate-400 hover:text-white"
              >
                ✕
              </button>
            </div>

            <div className="space-y-3 text-xs">
              <div>
                <label className="text-slate-400 block mb-1">API Endpoint</label>
                <input
                  type="text"
                  value={editEndpoint}
                  onChange={e => setEditEndpoint(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-white focus:outline-none focus:border-blue-500"
                />
              </div>

              <div>
                <label className="text-slate-400 block mb-1">App / Merchant ID</label>
                <input
                  type="text"
                  value={editAppID}
                  onChange={e => setEditAppID(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-white focus:outline-none focus:border-blue-500"
                />
              </div>

              <div>
                <label className="text-slate-400 block mb-1">Short Code</label>
                <input
                  type="text"
                  value={editShortCode}
                  onChange={e => setEditShortCode(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-white focus:outline-none focus:border-blue-500"
                />
              </div>

              <div>
                <label className="text-slate-400 block mb-1">App Key / Secret (Leave blank to keep existing)</label>
                <input
                  type="password"
                  value={editAppKey}
                  onChange={e => setEditAppKey(e.target.value)}
                  placeholder="••••••••••••••••"
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-white focus:outline-none focus:border-blue-500"
                />
              </div>

              <div className="flex items-center gap-2 pt-1">
                <input
                  type="checkbox"
                  id="providerEnabled"
                  checked={editEnabled}
                  onChange={e => setEditEnabled(e.target.checked)}
                  className="rounded text-blue-600"
                />
                <label htmlFor="providerEnabled" className="text-slate-300">Enable this provider</label>
              </div>
            </div>

            <div className="flex gap-2 pt-2">
              <button
                type="submit"
                disabled={savingProvider}
                className="flex-1 bg-blue-600 hover:bg-blue-500 text-white font-bold py-2 rounded-xl text-xs flex items-center justify-center gap-1"
              >
                {savingProvider ? <Spinner size={14} /> : 'Save Configuration'}
              </button>
              <button
                type="button"
                onClick={() => setEditingProvider(null)}
                className="px-4 py-2 border border-slate-800 rounded-xl text-slate-400 text-xs font-bold"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}
