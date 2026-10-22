'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Activity, AlertCircle, CheckCircle, Clock, XCircle, RefreshCw, Database, Package, Users } from 'lucide-react';

// Matches backend models.OperationsHealth
interface OperationsHealth {
  total_jobs_24h: number;
  failed_jobs_24h: number;
  stuck_jobs: number;
  inactive_users: number;
  stale_pending_products: number;
  expired_otps: number;
  stale_upgrade_requests: number;
  database_size: string;
  products_table_size: string;
  users_table_size: string;
}

// Matches backend models.JobStatusSummary
interface JobStatusSummary {
  job_name: string;
  total_runs: number;
  successful_runs: number;
  failed_runs: number;
  currently_running: number;
  last_run_at: string | null;
  last_success_at: string | null;
  last_failure_at: string | null;
  avg_duration_ms: number | null;
  total_records_processed: number;
}

// Matches backend models.JobExecutionLog
interface JobExecutionLog {
  id: number;
  job_name: string;
  status: string;
  started_at: string;
  completed_at: string | null;
  duration_ms: number | null;
  records_processed: number;
  error_message: string;
}

function StatusIcon({ status }: { status: string }) {
  if (status === 'success') return <CheckCircle className="w-4 h-4 text-emerald-400" />;
  if (status === 'failed') return <XCircle className="w-4 h-4 text-red-400" />;
  if (status === 'running') return <Clock className="w-4 h-4 text-blue-400 animate-spin" />;
  return <AlertCircle className="w-4 h-4 text-zinc-500" />;
}

function StatusBadge({ status }: { status: string }) {
  const map: Record<string, string> = {
    success: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
    failed:  'bg-red-500/10 text-red-400 border-red-500/20',
    running: 'bg-blue-500/10 text-blue-400 border-blue-500/20',
  };
  return (
    <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase border ${map[status] ?? 'bg-zinc-800 text-zinc-400 border-zinc-700'}`}>
      {status}
    </span>
  );
}

export default function OperationsPage() {
  const { user } = useAuth();
  const [health, setHealth] = useState<OperationsHealth | null>(null);
  const [jobs, setJobs] = useState<JobStatusSummary[]>([]);
  const [failedJobs, setFailedJobs] = useState<JobExecutionLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const fetchData = async (silent = false) => {
    if (!user?.token) return;
    if (!silent) setLoading(true);
    else setRefreshing(true);
    try {
      const [healthData, jobsData, failedData] = await Promise.allSettled([
        api.adminGetOperationsHealth(user.token),
        api.adminGetJobStatus(user.token),
        api.adminGetFailedJobs(user.token),
      ]);
      if (healthData.status === 'fulfilled' && healthData.value) setHealth(healthData.value as OperationsHealth);
      if (jobsData.status === 'fulfilled' && Array.isArray(jobsData.value)) setJobs(jobsData.value);
      if (failedData.status === 'fulfilled' && Array.isArray(failedData.value)) setFailedJobs(failedData.value);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(() => fetchData(true), 30_000);
    return () => clearInterval(interval);
  }, [user]); // eslint-disable-line react-hooks/exhaustive-deps

  const fmt = (v: number | undefined | null) => v?.toLocaleString() ?? '—';
  const fmtTime = (s: string | null) => s ? new Date(s).toLocaleString() : '—';
  const fmtMs = (ms: number | null) => ms != null ? `${ms.toFixed(0)}ms` : '—';

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-teal-600 rounded-2xl flex items-center justify-center animate-pulse">
            <Activity size={24} className="text-white" />
          </div>
          <div>
            <div className="h-5 w-48 bg-zinc-800 rounded animate-pulse" />
            <div className="h-3 w-32 bg-zinc-800 rounded animate-pulse mt-2" />
          </div>
        </div>
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="h-28 bg-zinc-900 border border-zinc-800 rounded-2xl animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-teal-600 rounded-2xl flex items-center justify-center shadow-lg shadow-teal-600/20">
            <Activity size={24} className="text-white" />
          </div>
          <div>
            <h1 className="text-xl font-black text-white">Operations Monitor</h1>
            <p className="text-xs text-zinc-400">Background jobs, data health, and database metrics</p>
          </div>
        </div>
        <button
          onClick={() => fetchData(true)}
          disabled={refreshing}
          className="flex items-center gap-2 px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl text-sm font-bold transition-colors disabled:opacity-50"
        >
          <RefreshCw size={14} className={refreshing ? 'animate-spin' : ''} />
          Refresh
        </button>
      </div>

      {/* Health KPIs */}
      {health && (
        <>
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
            <KPIBox
              label="Jobs (24h)"
              value={fmt(health.total_jobs_24h)}
              sub={`${health.failed_jobs_24h} failed`}
              icon={<Activity size={18} />}
              alert={health.failed_jobs_24h > 0}
            />
            <KPIBox
              label="Stuck Jobs"
              value={fmt(health.stuck_jobs)}
              sub="running > 1 hour"
              icon={<Clock size={18} />}
              alert={health.stuck_jobs > 0}
            />
            <KPIBox
              label="Stale Listings"
              value={fmt(health.stale_pending_products)}
              sub="pending > 7 days"
              icon={<Package size={18} />}
              alert={health.stale_pending_products > 0}
            />
            <KPIBox
              label="Stale Upgrades"
              value={fmt(health.stale_upgrade_requests)}
              sub="pending > 14 days"
              icon={<Users size={18} />}
              alert={health.stale_upgrade_requests > 0}
            />
          </div>

          {/* Database sizes */}
          <div className="grid grid-cols-3 gap-4">
            {[
              { label: 'Database Total', value: health.database_size },
              { label: 'Products Table', value: health.products_table_size },
              { label: 'Users Table',    value: health.users_table_size },
            ].map(({ label, value }) => (
              <div key={label} className="bg-zinc-900 border border-zinc-800 rounded-2xl p-4 flex items-center gap-3">
                <div className="p-2 bg-zinc-800 rounded-xl">
                  <Database size={16} className="text-zinc-400" />
                </div>
                <div>
                  <p className="text-xs text-zinc-500 font-bold uppercase tracking-wider">{label}</p>
                  <p className="text-lg font-black text-white font-mono">{value ?? '—'}</p>
                </div>
              </div>
            ))}
          </div>
        </>
      )}

      {/* Job Status Summary */}
      <div className="bg-zinc-900 border border-zinc-800 rounded-2xl overflow-hidden">
        <div className="px-5 py-4 border-b border-zinc-800 flex items-center gap-2">
          <Activity size={15} className="text-orange-500" />
          <h2 className="font-black text-sm text-white uppercase tracking-wider">Background Jobs (Last 7 Days)</h2>
        </div>
        {jobs.length === 0 ? (
          <div className="p-8 text-center text-zinc-500">
            <Activity size={32} className="mx-auto mb-2 opacity-20" />
            <p className="text-sm">No job data — jobs run automatically on a schedule</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="text-zinc-500 text-xs uppercase font-bold border-b border-zinc-800">
                <th className="px-5 py-3 text-left">Job Name</th>
                <th className="px-4 py-3 text-right">Runs</th>
                <th className="px-4 py-3 text-right">Success</th>
                <th className="px-4 py-3 text-right">Failed</th>
                <th className="px-4 py-3 text-right">Avg Duration</th>
                <th className="px-5 py-3 text-right">Last Run</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/50">
              {jobs.map((job, i) => (
                <tr key={i} className="hover:bg-zinc-800/30 transition-colors">
                  <td className="px-5 py-3">
                    <div className="flex items-center gap-2">
                      <StatusIcon status={job.failed_runs > 0 ? 'failed' : 'success'} />
                      <span className="font-mono text-xs text-white">{job.job_name}</span>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-right font-mono text-zinc-300">{fmt(job.total_runs)}</td>
                  <td className="px-4 py-3 text-right font-mono text-emerald-400">{fmt(job.successful_runs)}</td>
                  <td className={`px-4 py-3 text-right font-mono font-bold ${job.failed_runs > 0 ? 'text-red-400' : 'text-zinc-500'}`}>
                    {fmt(job.failed_runs)}
                  </td>
                  <td className="px-4 py-3 text-right font-mono text-zinc-400 text-xs">{fmtMs(job.avg_duration_ms)}</td>
                  <td className="px-5 py-3 text-right text-zinc-400 text-xs">{fmtTime(job.last_run_at)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Failed Jobs */}
      {failedJobs.length > 0 && (
        <div className="bg-zinc-900 border border-red-500/20 rounded-2xl overflow-hidden">
          <div className="px-5 py-4 border-b border-red-500/20 bg-red-900/10 flex items-center gap-2">
            <XCircle size={15} className="text-red-400" />
            <h2 className="font-black text-sm text-red-300 uppercase tracking-wider">
              Recent Failures ({failedJobs.length})
            </h2>
          </div>
          <div className="divide-y divide-zinc-800/50">
            {failedJobs.map((job) => (
              <div key={job.id} className="px-5 py-4">
                <div className="flex items-start justify-between gap-4">
                  <div className="flex items-start gap-3 min-w-0">
                    <XCircle className="w-4 h-4 text-red-400 mt-0.5 shrink-0" />
                    <div className="min-w-0">
                      <p className="font-mono font-bold text-sm text-white">{job.job_name}</p>
                      {job.error_message && (
                        <p className="text-xs text-red-300 mt-1 break-words">{job.error_message}</p>
                      )}
                    </div>
                  </div>
                  <div className="text-right shrink-0">
                    <StatusBadge status={job.status} />
                    <p className="text-[10px] text-zinc-500 mt-1">{fmtTime(job.started_at)}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

function KPIBox({ label, value, sub, icon, alert }: {
  label: string; value: string; sub?: string; icon: React.ReactNode; alert?: boolean;
}) {
  return (
    <div className={`p-4 bg-zinc-900 border rounded-2xl space-y-2 ${alert ? 'border-red-500/30' : 'border-zinc-800'}`}>
      <div className="flex items-center justify-between">
        <span className="text-xs font-bold uppercase tracking-wider text-zinc-400">{label}</span>
        <div className={`p-1.5 rounded-lg ${alert ? 'bg-red-600/10 text-red-400' : 'bg-zinc-800 text-zinc-400'}`}>
          {icon}
        </div>
      </div>
      <p className={`text-2xl font-black font-mono ${alert ? 'text-red-300' : 'text-white'}`}>{value}</p>
      {sub && <p className="text-[10px] text-zinc-500 font-medium">{sub}</p>}
    </div>
  );
}
