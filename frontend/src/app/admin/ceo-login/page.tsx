import { Metadata } from 'next';
import { AdminPortalLogin } from '@/components/admin/enterprise/AdminPortalLogin';

export const metadata: Metadata = {
  title: 'CEO & Executive Login | HUNTIQ Enterprise',
  description: 'Tier-0 authorized governance login portal for HUNTIQ Executive Leadership.',
};

export default function AdminCeoLoginPage() {
  return (
    <AdminPortalLogin
      portal="executive"
      portalTitle="HUNTIQ CEO Executive Portal"
      portalBadge="Tier-0 Governance"
      portalTheme="crimson"
      portalSubtitle="Phishing-resistant WebAuthn hardware authentication for Executive Leadership"
      portalDescription="Mandatory cryptographic user verification and hardware-backed biometric security."
      allowPasswordMfaFallback={true}
    />
  );
}
