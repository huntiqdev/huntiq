import { Metadata } from 'next';
import { AdminPortalLogin } from '@/components/admin/enterprise/AdminPortalLogin';

export const metadata: Metadata = {
  title: 'Staff Access Portal | HUNTIQ Enterprise',
  description: 'Operational and content moderation access portal for HUNTIQ administrative staff.',
};

export default function StaffLoginPage() {
  return (
    <AdminPortalLogin
      portal="staff"
      portalTitle="HUNTIQ Staff Access"
      portalBadge="Staff Operations"
      portalTheme="orange"
      portalSubtitle="Operational administration, listing moderation, and customer support portal"
      portalDescription="WebAuthn Passkey sign-in with verified Staff Password & MFA fallback."
      allowPasswordMfaFallback={true}
    />
  );
}
