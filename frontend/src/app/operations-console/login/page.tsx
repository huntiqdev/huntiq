import { Metadata } from 'next';
import { AdminPortalLogin } from '@/components/admin/enterprise/AdminPortalLogin';

export const metadata: Metadata = {
  title: 'Operations Console Login | HUNTIQ Enterprise',
  description: 'Authorized administrative portal for HUNTIQ Operations Administrators, Inspection, and Verification.',
};

export default function OperationsConsoleLoginPage() {
  return (
    <AdminPortalLogin
      portal="operations"
      portalTitle="HUNTIQ Operations Console"
      portalBadge="Operations & Inspection"
      portalTheme="orange"
      portalSubtitle="Authorized logistics, seller inspection, KYC verification, and transaction operations"
      portalDescription="Mandatory WebAuthn passkey authentication with verified Staff MFA fallback."
      allowPasswordMfaFallback={true}
    />
  );
}
