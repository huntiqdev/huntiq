'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import Link from 'next/link';
import {
  ShieldAlert,
  Server,
  ShieldCheck,
  Users,
  Fingerprint,
  ArrowRight,
  Lock,
  ChevronRight,
  Shield,
} from 'lucide-react';

import { HuntiqLogo } from '@/components/ui/HuntiqLogo';

interface PortalCard {
  title: string;
  badge: string;
  badgeClass: string;
  desc: string;
  href: string;
  dedicatedUrl: string;
  icon: any;
  iconColor: string;
  borderHover: string;
  roles: string;
  securityLevel: string;
}

const PORTALS: PortalCard[] = [
  {
    title: 'Platform Owner & CEO',
    badge: 'Tier-0 Governance',
    badgeClass: 'bg-red-500/10 text-red-400 border-red-500/30',
    desc: 'Ultimate authorized platform governance, root configuration, and administrative provisioning.',
    href: '/admin/owner-login',
    dedicatedUrl: '/admin/owner-login (or /executive-console/login)',
    icon: ShieldAlert,
    iconColor: 'text-red-400 bg-red-950/40 border-red-500/20',
    borderHover: 'hover:border-red-500/50 hover:shadow-red-500/10',
    roles: 'Platform Owner, CEO',
    securityLevel: 'Passkey Required (User Verification Enforced)',
  },
  {
    title: 'Platform Administrator',
    badge: 'Operations & Config',
    badgeClass: 'bg-indigo-500/10 text-indigo-400 border-indigo-500/30',
    desc: 'Platform configuration, service administration, system settings, and catalog structure management.',
    href: '/admin/platform-login',
    dedicatedUrl: '/admin/platform-login (or /platform-console/login)',
    icon: Server,
    iconColor: 'text-indigo-400 bg-indigo-950/40 border-indigo-500/20',
    borderHover: 'hover:border-indigo-500/50 hover:shadow-indigo-500/10',
    roles: 'Platform Administrator',
    securityLevel: 'Passkey Required (User Verification Enforced)',
  },
  {
    title: 'Security Administrator',
    badge: 'Security & Audit',
    badgeClass: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30',
    desc: 'Security incident monitoring, tamper-resistant audit logs, session revocation, and access reviews.',
    href: '/admin/security-login',
    dedicatedUrl: '/admin/security-login (or /security-console/login)',
    icon: ShieldCheck,
    iconColor: 'text-emerald-400 bg-emerald-950/40 border-emerald-500/20',
    borderHover: 'hover:border-emerald-500/50 hover:shadow-emerald-500/10',
    roles: 'Security Administrator',
    securityLevel: 'Passkey Required (User Verification Enforced)',
  },
  {
    title: 'Enterprise Staff Access',
    badge: 'Operations & Support',
    badgeClass: 'bg-blue-500/10 text-blue-400 border-blue-500/30',
    desc: 'Routine catalog management, listing moderation, seller verification, and customer inquiries.',
    href: '/admin/staff-login',
    dedicatedUrl: '/admin/staff-login (or /staff/login)',
    icon: Users,
    iconColor: 'text-blue-400 bg-blue-950/40 border-blue-500/20',
    borderHover: 'hover:border-blue-500/50 hover:shadow-blue-500/10',
    roles: 'Operations Admin, Content Moderator, Support Agent',
    securityLevel: 'Passkey Recommended / Mandatory Staff MFA',
  },
];

export default function AdminPortalDirectoryPage() {
  const router = useRouter();
  const { user } = useAuth();

  useEffect(() => {
    if (user?.role === 'admin') {
      router.replace('/admin/dashboard');
    }
  }, [user, router]);

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100 flex flex-col justify-between py-12 px-4 sm:px-6 lg:px-8 selection:bg-[#1769E0] selection:text-white">
      {/* Background ambient gradient glow */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 left-1/2 -translate-x-1/2 w-[900px] h-[500px] bg-[#1769E0]/5 blur-[160px] rounded-full" />
      </div>

      <div className="max-w-5xl mx-auto w-full relative z-10 flex-1 flex flex-col justify-center">
        {/* Header */}
        <div className="text-center space-y-3 mb-12 flex flex-col items-center">
          <HuntiqLogo size="lg" className="mb-2" />
          <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-zinc-900 border border-zinc-800 text-xs font-bold text-zinc-400">
            <Shield size={14} className="text-[#1769E0]" />
            <span>HUNTIQ Privileged Security Architecture</span>
          </div>
          <h1 className="text-3xl sm:text-4xl lg:text-5xl font-black tracking-tight text-white">
            Administrative Access Portals
          </h1>
          <p className="text-zinc-400 text-sm sm:text-base max-w-2xl mx-auto font-normal">
            Select your assigned administrative portal. Privileged roles require phishing-resistant
            hardware passkey verification with zero biometric exposure.
          </p>
        </div>

        {/* Portal Selection Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-5 mb-12">
          {PORTALS.map((portal) => {
            const Icon = portal.icon;
            return (
              <Link
                key={portal.title}
                href={portal.href}
                className={`group relative bg-zinc-900/80 border border-zinc-800/90 rounded-3xl p-6 sm:p-7 transition-all duration-200 hover:-translate-y-1 shadow-lg ${portal.borderHover}`}
              >
                <div className="flex items-start justify-between mb-4">
                  <div className={`p-3 rounded-2xl border ${portal.iconColor}`}>
                    <Icon size={24} />
                  </div>
                  <span className={`text-[10px] font-black uppercase tracking-wider px-2.5 py-1 rounded-full border ${portal.badgeClass}`}>
                    {portal.badge}
                  </span>
                </div>

                <h2 className="text-lg font-black text-white group-hover:text-[#38A5FF] transition-colors flex items-center gap-2">
                  <span>{portal.title}</span>
                  <ChevronRight size={16} className="opacity-0 group-hover:opacity-100 -translate-x-2 group-hover:translate-x-0 transition-all text-[#38A5FF]" />
                </h2>

                <p className="text-zinc-400 text-xs leading-relaxed mt-2 mb-4 font-normal">
                  {portal.desc}
                </p>

                <div className="pt-4 border-t border-zinc-800/60 space-y-1.5 text-[11px]">
                  <div className="flex items-center justify-between text-zinc-400">
                    <span className="text-zinc-500 font-semibold">Direct Login:</span>
                    <span className="font-mono text-zinc-300 text-[10px] bg-zinc-800/60 px-2 py-0.5 rounded-md">{portal.href}</span>
                  </div>
                  <div className="flex items-center justify-between text-zinc-400">
                    <span className="text-zinc-500 font-semibold">Authorized:</span>
                    <span className="font-bold text-zinc-300">{portal.roles}</span>
                  </div>
                  <div className="flex items-center justify-between text-zinc-400">
                    <span className="text-zinc-500 font-semibold">Assurance:</span>
                    <span className="font-medium text-zinc-400 flex items-center gap-1">
                      <Fingerprint size={12} className="text-zinc-400" />
                      {portal.securityLevel}
                    </span>
                  </div>
                </div>
              </Link>
            );
          })}
        </div>

        {/* Security Policy Footnote */}
        <div className="bg-zinc-900/40 border border-zinc-900 rounded-2xl p-4 text-center max-w-2xl mx-auto">
          <p className="text-[11px] text-zinc-500 leading-relaxed">
            HUNTIQ enforces strict separation of duties and zero-trust backend authorization.
            All administrative access ceremonies, session rotations, and audit events are cryptographically recorded.
          </p>
        </div>
      </div>

      {/* Footer */}
      <footer className="text-center text-xs text-zinc-600 mt-8 relative z-10">
        &copy; {new Date().getFullYear()} HUNTIQ Universal Marketplace Platform. All rights reserved.
      </footer>
    </div>
  );
}
