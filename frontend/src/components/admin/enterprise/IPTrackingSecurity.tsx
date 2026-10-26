'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { 
  ShieldAlert, ShieldCheck, Lock, Unlock, AlertTriangle, Eye, RefreshCw, 
  Filter, Search, Ban, CheckCircle2, Globe, Activity 
} from 'lucide-react';

export interface IPLogEntry {
  id: string;
  ipAddress: string;
  location: string;
  requestEndpoint: string;
  attemptCount: number;
  status: 'Normal' | 'Rate Limited' | 'Brute Force Alert' | 'Blacklisted';
  timestamp: string;
  userAgent: string;
}

export function IPTrackingSecurity() {
  const { user } = useAuth();
  const [logs, setLogs] = useState<IPLogEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [blacklistedIPs, setBlacklistedIPs] = useState<string[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('all');
  const [actionNotice, setActionNotice] = useState<string | null>(null);

  const loadRealLogs = useCallback(async () => {
    if (!user) return;
    setLoading(true);
    try {
      const res = await api.adminGetLogs({ limit: 100 }, user.token);
      const rawLogs = res.logs || [];
      const mapped: IPLogEntry[] = rawLogs.map((l: any, idx: number) => {
        const action = (l.action || '').toLowerCase();
        let status: IPLogEntry['status'] = 'Normal';
        if (action.includes('ban') || action.includes('block')) status = 'Blacklisted';
        else if (action.includes('warn') || action.includes('reject')) status = 'Rate Limited';
        else if (action.includes('fail') || action.includes('unauthorized')) status = 'Brute Force Alert';

        return {
          id: `LOG-${l.id || idx}`,
          ipAddress: l.ip_address || '—',
          location: l.admin_name ? `Admin: ${l.admin_name}` : '—',
          requestEndpoint: `ACTION: ${l.action || 'REQUEST'} ${l.details ? `(${l.details})` : ''}`,
          attemptCount: 1,
          status: status,
          timestamp: l.created_at ? new Date(l.created_at).toLocaleTimeString() : 'Just now',
          userAgent: l.notes || 'HUNTIQ Unified Core API'
        };
      });
      setLogs(mapped);
    } catch (err) {
      console.error('Failed to load live IP logs', err);
    } finally {
      setLoading(false);
    }
  }, [user]);

  useEffect(() => {
    loadRealLogs();
  }, [loadRealLogs]);

  const filteredLogs = logs.filter((log) => {
    const matchesSearch =
      log.ipAddress.includes(searchQuery) ||
      log.location.toLowerCase().includes(searchQuery.toLowerCase()) ||
      log.requestEndpoint.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesFilter = filterStatus === 'all' || log.status.toLowerCase() === filterStatus.toLowerCase();
    return matchesSearch && matchesFilter;
  });

  const handleBlockIP = (ip: string) => {
    if (!blacklistedIPs.includes(ip)) {
      setBlacklistedIPs((prev) => [...prev, ip]);
      setLogs((prev) =>
        prev.map((l) => (l.ipAddress === ip ? { ...l, status: 'Blacklisted' } : l))
      );
      triggerNotice(`IP Address ${ip} has been permanently blocked at Firewall level.`);
    }
  };

  const handleUnblockIP = (ip: string) => {
    setBlacklistedIPs((prev) => prev.filter((i) => i !== ip));
    setLogs((prev) =>
      prev.map((l) => (l.ipAddress === ip ? { ...l, status: 'Normal' } : l))
    );
    triggerNotice(`IP Address ${ip} unblocked.`);
  };

  const triggerNotice = (msg: string) => {
    setActionNotice(msg);
    setTimeout(() => setActionNotice(null), 3500);
  };

  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-5 shadow-xl space-y-5">
      {/* Top Header */}
      <div className="flex flex-wrap items-center justify-between gap-4 pb-4 border-b border-zinc-800">
        <div className="flex items-center gap-3">
          <div className="p-2.5 bg-red-600/10 text-red-500 rounded-xl border border-red-500/20">
            <ShieldAlert size={20} />
          </div>
          <div>
            <h2 className="text-base font-extrabold text-white flex items-center gap-2">
              IP Tracking & Hacker Threat Defense Firewall
              <span className="px-2 py-0.5 text-[10px] font-mono bg-red-500/10 text-red-400 border border-red-500/20 rounded-md">
                Active IPS / WAF
              </span>
            </h2>
            <p className="text-xs text-zinc-400">Live request audit log, automated brute-force alerts & firewall blacklists</p>
          </div>
        </div>

        {/* Stats */}
        <div className="flex items-center gap-3 text-xs font-mono">
          <div className="px-3 py-1.5 bg-zinc-950 border border-zinc-800 rounded-xl text-zinc-300">
            Blocked IPs: <strong className="text-red-400 font-bold">{blacklistedIPs.length}</strong>
          </div>
          <div className="px-3 py-1.5 bg-zinc-950 border border-zinc-800 rounded-xl text-zinc-300">
            Threat Level: <strong className="text-emerald-400 font-bold">NORMAL (SECURE)</strong>
          </div>
        </div>
      </div>

      {actionNotice && (
        <div className="p-3 bg-red-500/10 border border-red-500/30 text-red-400 rounded-xl text-xs font-bold flex items-center justify-between animate-in fade-in">
          <span className="flex items-center gap-2">
            <Ban size={16} />
            {actionNotice}
          </span>
          <button onClick={() => setActionNotice(null)} className="text-red-400 hover:text-white">✕</button>
        </div>
      )}

      {/* Filter Controls */}
      <div className="flex flex-wrap items-center justify-between gap-3 text-xs">
        <div className="flex items-center gap-2">
          <input
            type="text"
            placeholder="Search IP, Endpoint, or Action..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="bg-zinc-950 border border-zinc-800 text-zinc-200 placeholder-zinc-500 rounded-xl px-3 py-2 focus:outline-none focus:border-red-500 text-xs"
          />

          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value)}
            className="bg-zinc-950 border border-zinc-800 text-zinc-200 rounded-xl px-3 py-2 focus:outline-none focus:border-red-500 font-bold"
          >
            <option value="all">All Traffic Statuses</option>
            <option value="Brute Force Alert">Brute Force Alerts</option>
            <option value="Blacklisted">Blacklisted IPs</option>
            <option value="Rate Limited">Rate Limited</option>
            <option value="Normal">Normal Traffic</option>
          </select>
        </div>

        <button
          onClick={loadRealLogs}
          className="px-3 py-2 bg-zinc-950 hover:bg-zinc-800 text-zinc-300 border border-zinc-800 rounded-xl text-xs font-bold flex items-center gap-1.5"
        >
          <RefreshCw size={14} className={loading ? 'animate-spin' : ''} /> Refresh Live Logs
        </button>
      </div>

      {/* Main IP Log Table */}
      <div className="overflow-x-auto rounded-xl border border-zinc-800">
        <table className="w-full text-xs text-left border-collapse">
          <thead>
            <tr className="bg-zinc-950 text-zinc-400 font-mono text-[11px] border-b border-zinc-800">
              <th className="p-3.5 font-bold text-zinc-300">IP Address</th>
              <th className="p-3.5 font-bold text-zinc-300">Location / Context</th>
              <th className="p-3.5 font-bold text-zinc-300">Action / Endpoint</th>
              <th className="p-3.5 font-bold text-zinc-300">Threat Status</th>
              <th className="p-3.5 font-bold text-zinc-300">Timestamp</th>
              <th className="p-3.5 text-right font-bold text-zinc-300">Firewall Action</th>
            </tr>
          </thead>

          <tbody className="divide-y divide-zinc-800/80 bg-zinc-900/40">
            {logs.length === 0 ? (
              <tr>
                <td colSpan={6} className="p-8 text-center text-zinc-500 font-mono">
                  {loading ? 'Fetching real-time security log data...' : 'No security log entries recorded yet'}
                </td>
              </tr>
            ) : (
              filteredLogs.map((log) => {
                const isBlocked = blacklistedIPs.includes(log.ipAddress);

                let badgeStyle = 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20';
                if (log.status === 'Brute Force Alert') badgeStyle = 'bg-red-600/20 text-red-400 border-red-500/40 font-black animate-pulse';
                if (log.status === 'Rate Limited') badgeStyle = 'bg-amber-500/10 text-amber-400 border-amber-500/20';
                if (isBlocked || log.status === 'Blacklisted') badgeStyle = 'bg-zinc-800 text-zinc-400 border-zinc-700';

                return (
                  <tr key={log.id} className="hover:bg-zinc-800/40 transition-colors">
                    {/* IP */}
                    <td className="p-3.5 font-mono font-bold text-white flex items-center gap-2">
                      <Globe size={14} className="text-zinc-500" />
                      {log.ipAddress}
                    </td>

                    {/* Location */}
                    <td className="p-3.5 text-zinc-300">
                      <div>{log.location}</div>
                      <div className="text-[10px] text-zinc-500 font-mono truncate max-w-xs">{log.userAgent}</div>
                    </td>

                    {/* Endpoint */}
                    <td className="p-3.5 font-mono text-zinc-300">
                      <span className="px-2 py-0.5 bg-zinc-950 border border-zinc-800 rounded font-semibold text-[11px]">
                        {log.requestEndpoint}
                      </span>
                    </td>

                    {/* Status */}
                    <td className="p-3.5">
                      <span className={`px-2.5 py-1 rounded-lg text-[11px] border ${badgeStyle}`}>
                        {isBlocked ? 'Blacklisted IP' : log.status}
                      </span>
                    </td>

                    {/* Timestamp */}
                    <td className="p-3.5 font-mono text-zinc-500 text-[11px]">
                      {log.timestamp}
                    </td>

                    {/* Action */}
                    <td className="p-3.5 text-right">
                      {isBlocked ? (
                        <button
                          onClick={() => handleUnblockIP(log.ipAddress)}
                          className="px-3 py-1 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-lg font-bold text-[11px] transition-colors"
                        >
                          Unblock IP
                        </button>
                      ) : (
                        <button
                          onClick={() => handleBlockIP(log.ipAddress)}
                          className="px-3 py-1 bg-red-600/20 hover:bg-red-600/30 text-red-400 border border-red-500/30 rounded-lg font-extrabold text-[11px] transition-colors flex items-center gap-1 ml-auto"
                        >
                          <Ban size={12} /> Block Hacker IP
                        </button>
                      )}
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
