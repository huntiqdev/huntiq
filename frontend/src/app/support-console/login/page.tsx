import { Metadata } from 'next';
import { AdminPortalLogin } from '@/components/admin/enterprise/AdminPortalLogin';

export const metadata: Metadata = {
  title: 'Support Console Login | HUNTIQ Enterprise',
  description: 'Authorized administrative portal for HUNTIQ Customer Support Agents and Concierge Leads.',
};

export default function SupportConsoleLoginPage() {
  return (
    <AdminPortalLogin
      portal="support"
      portalTitle="HUNTIQ Support Console"
      portalBadge="Support & Concierge"
      portalTheme="blue"
      portalSubtitle="Authorized buyer inquiries, seller support, dispute triage, and VIP concierge"
      portalDescription="Hardware passkey sign-in with verified Staff Password & MFA fallback."
      allowPasswordMfaFallback={true}
    />
  );
}
