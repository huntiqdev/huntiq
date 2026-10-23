import { Metadata } from 'next';
import { AdminPortalLogin } from '@/components/admin/enterprise/AdminPortalLogin';

export const metadata: Metadata = {
  title: 'Platform Owner Console | HUNTIQ Enterprise',
  description: 'Tier-0 authorized governance login portal for HUNTIQ Platform Owners.',
};

export default function ExecutiveConsoleLoginPage() {
  return (
    <AdminPortalLogin
      portal="executive"
      portalTitle="HUNTIQ Platform Owner"
      portalBadge="Tier-0 Governance"
      portalTheme="crimson"
      portalSubtitle="Phishing-resistant WebAuthn hardware authentication for Platform Owners"
      portalDescription="Mandatory user verification and hardware-backed biometric security."
      allowPasswordMfaFallback={false}
    />
  );
}
