import { Metadata } from 'next';
import { AdminPortalLogin } from '@/components/admin/enterprise/AdminPortalLogin';

export const metadata: Metadata = {
  title: 'Platform Administrator Console | HUNTIQ Enterprise',
  description: 'Authorized administrative login portal for HUNTIQ Platform Administrators.',
};

export default function PlatformConsoleLoginPage() {
  return (
    <AdminPortalLogin
      portal="platform"
      portalTitle="HUNTIQ Platform Administrator"
      portalBadge="Platform Operations"
      portalTheme="indigo"
      portalSubtitle="Authorized platform configuration, system management, and catalog governance"
      portalDescription="Mandatory WebAuthn passkey authentication with user verification."
      allowPasswordMfaFallback={false}
    />
  );
}
