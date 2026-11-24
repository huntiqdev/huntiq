-- Migration 066: Enterprise Sub-Security Separation and Real-Time Admin Telemetry
-- National-scale financial and administrative governance for HUNTIQ

-- 1. Create Sub-Security Zones table
CREATE TABLE IF NOT EXISTS public.admin_security_zones (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '',
    required_role VARCHAR(50) NOT NULL DEFAULT 'security_admin',
    ip_whitelist TEXT[] DEFAULT '{}',
    require_trusted_device BOOLEAN NOT NULL DEFAULT true,
    require_passkey BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 2. Create Sub-Security Zone Assignments table
CREATE TABLE IF NOT EXISTS public.admin_security_zone_assignments (
    id SERIAL PRIMARY KEY,
    zone_id INT NOT NULL REFERENCES public.admin_security_zones(id) ON DELETE CASCADE,
    admin_id INT NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    granted_by INT REFERENCES public.users(id) ON DELETE SET NULL,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ,
    notes TEXT DEFAULT '',
    CONSTRAINT uq_admin_security_zone UNIQUE (zone_id, admin_id)
);

CREATE INDEX IF NOT EXISTS idx_sec_zone_assign_admin ON public.admin_security_zone_assignments(admin_id);
CREATE INDEX IF NOT EXISTS idx_sec_zone_assign_zone ON public.admin_security_zone_assignments(zone_id);

-- 3. Enhance admin_sessions and admin_logs with geographic and device fingerprint metadata if missing
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'admin_sessions' AND column_name = 'geo_location'
    ) THEN
        ALTER TABLE public.admin_sessions ADD COLUMN geo_location VARCHAR(100) DEFAULT 'Ethiopia / Addis Ababa Metro';
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'admin_sessions' AND column_name = 'device_fingerprint'
    ) THEN
        ALTER TABLE public.admin_sessions ADD COLUMN device_fingerprint VARCHAR(128);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'admin_logs' AND column_name = 'device_id'
    ) THEN
        ALTER TABLE public.admin_logs ADD COLUMN device_id VARCHAR(64);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'admin_logs' AND column_name = 'session_id'
    ) THEN
        ALTER TABLE public.admin_logs ADD COLUMN session_id INT;
    END IF;
END $$;

-- 4. Seed Canonical Sub-Security Zones
INSERT INTO public.admin_security_zones (code, name, description, required_role, require_trusted_device, require_passkey)
VALUES 
('SOC_DEFENSE', 'SOC Threat & Attack Defense Zone', 'Perimeter intrusion detection, DDoS mitigation, IP quarantine, and live threat telemetry.', 'ciso', true, true),
('FRAUD_AML', 'Fraud & Financial AML Investigation Zone', 'Suspicious transaction pattern analysis, bank ledger anomaly detection, and account freeze.', 'ciso', true, false),
('KYC_COMPLIANCE', 'Identity, KYC & FIDA Regulatory Zone', 'Verification of business licenses, taxpayer identification, and citizen digital ID records.', 'security_admin', true, false),
('GATEWAY_OPS', 'Payment Gateway & Telecom Operations Zone', 'Telebirr API gateway health, CBE Birr integration monitor, and SMS OTP dispatch infrastructure.', 'ciso', true, true),
('DISPUTE_ESCROW', 'Dispute Resolution & Escrow Audit Zone', 'Buyer-seller arbitration, escrow funds audit, and manual settlement reconciliation.', 'security_admin', false, false)
ON CONFLICT (code) DO NOTHING;

-- 5. Seed initial zone assignments for CISO (ID 206), Security Admin (ID 207), SOC Analyst (ID 208), Fraud Investigator (ID 209), Compliance Auditor (ID 210)
INSERT INTO public.admin_security_zone_assignments (zone_id, admin_id, granted_by, notes)
SELECT z.id, u.id, 202, 'Initial country-level security governance assignment'
FROM public.admin_security_zones z
CROSS JOIN public.users u
WHERE (z.code = 'SOC_DEFENSE' AND u.id IN (206, 207, 208))
   OR (z.code = 'FRAUD_AML' AND u.id IN (206, 209, 210))
   OR (z.code = 'KYC_COMPLIANCE' AND u.id IN (206, 207, 210, 222))
   OR (z.code = 'GATEWAY_OPS' AND u.id IN (206, 207))
   OR (z.code = 'DISPUTE_ESCROW' AND u.id IN (206, 207, 224, 225))
ON CONFLICT (zone_id, admin_id) DO NOTHING;
