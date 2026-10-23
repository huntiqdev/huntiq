import { Metadata } from 'next';
import { AdminPortalLogin } from '@/components/admin/enterprise/AdminPortalLogin';

export const metadata: Metadata = {
  title: 'Security Administrator Console | HUNTIQ Enterprise',
  description: 'Privileged security monitoring and audit investigation portal for HUNTIQ Security Administrators.',
};

export default function SecurityConsoleLoginPage() {
  return (
    <AdminPortalLogin
      portal="security"
      portalTitle="HUNTIQ Security Administrator"
      portalBadge="Security & Audit"
      portalTheme="emerald"
      portalSubtitle="Privileged access for security investigation, audit logs, and incident response"
      portalDescription="Mandatory cryptographic user verification via hardware authenticators."
      allowPasswordMfaFallback={false}
    />
  );
}
