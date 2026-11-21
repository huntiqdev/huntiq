import { Metadata } from 'next';
import { AdminPortalLogin } from '@/components/admin/enterprise/AdminPortalLogin';

export const metadata: Metadata = {
  title: 'Moderation Console Login | HUNTIQ Enterprise',
  description: 'Authorized administrative portal for HUNTIQ Content Moderators and Trust & Safety Officers.',
};

export default function ModerationConsoleLoginPage() {
  return (
    <AdminPortalLogin
      portal="moderation"
      portalTitle="HUNTIQ Moderation Console"
      portalBadge="Trust & Moderation"
      portalTheme="purple"
      portalSubtitle="Authorized content governance, listing moderation, and marketplace trust enforcement"
      portalDescription="Cryptographic hardware passkey sign-in with verified Staff MFA fallback."
      allowPasswordMfaFallback={true}
    />
  );
}
