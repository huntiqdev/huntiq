import { Metadata } from 'next';
import { AdminPortalLogin } from '@/components/admin/enterprise/AdminPortalLogin';

export const metadata: Metadata = {
  title: 'Platform Owner / CEO Login | HUNTIQ Enterprise',
  description: 'Tier-0 authorized governance login portal for HUNTIQ Platform Owners and Executive Leadership.',
};

export default function AdminOwnerLoginPage() {
  return (
    <AdminPortalLogin
      portal="executive"
      portalTitle="HUNTIQ Platform Owner & CEO"
      portalBadge="Tier-0 Governance"
      portalTheme="crimson"
      portalSubtitle="Phishing-resistant WebAuthn hardware authentication for Platform Owners & CEO"
      portalDescription="Mandatory cryptographic user verification and hardware-backed biometric security."
      allowPasswordMfaFallback={true}
    />
  );
}
