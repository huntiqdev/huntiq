-- ==============================================================================
-- HUNTIQ ENTERPRISE MIGRATION 067: SHIFTS, WORKSTATIONS & 24/7 OPERATIONAL TELEMETRY
-- Ensures strict individual operator attribution across time shifts & physical PCs.
-- ==============================================================================

-- 1. Canonical Operational Shifts
CREATE TABLE IF NOT EXISTS admin_shifts (
    id SERIAL PRIMARY KEY,
    code VARCHAR(30) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    start_hour INT NOT NULL,
    end_hour INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Physical Workstations / Terminal Inventory
CREATE TABLE IF NOT EXISTS admin_workstations (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    department VARCHAR(50) NOT NULL,
    facility_location VARCHAR(150) NOT NULL,
    allowed_ip_cidr VARCHAR(50) NOT NULL DEFAULT '127.0.0.1/32',
    hardware_uuid VARCHAR(128) NOT NULL,
    is_dedicated BOOLEAN NOT NULL DEFAULT true,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Enhance admin_sessions and users with shift and workstation attribution
ALTER TABLE admin_sessions
    ADD COLUMN IF NOT EXISTS shift_code VARCHAR(30),
    ADD COLUMN IF NOT EXISTS workstation_code VARCHAR(50);

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS default_shift VARCHAR(30) DEFAULT 'DAY_SHIFT',
    ADD COLUMN IF NOT EXISTS assigned_workstation VARCHAR(50);

-- Indexes for rapid operational lookup
CREATE INDEX IF NOT EXISTS idx_admin_sessions_shift ON admin_sessions(shift_code);
CREATE INDEX IF NOT EXISTS idx_admin_sessions_workstation ON admin_sessions(workstation_code);

-- 4. Seed Canonical 24/7 Operational Shifts
INSERT INTO admin_shifts (code, name, start_hour, end_hour, is_active)
VALUES
    ('DAY_SHIFT',   'Morning Day Shift (08:00 - 16:00)', 8, 16, true),
    ('SWING_SHIFT', 'Evening Swing Shift (16:00 - 00:00)', 16, 0, true),
    ('NIGHT_SHIFT', 'Graveyard Night Shift (00:00 - 08:00)', 0, 8, true)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name, start_hour = EXCLUDED.start_hour, end_hour = EXCLUDED.end_hour;

-- 5. Seed Enterprise Physical Workstations
INSERT INTO admin_workstations (code, department, facility_location, allowed_ip_cidr, hardware_uuid, is_dedicated)
VALUES
    ('WS-EXEC-01', 'Executive', 'Addis Ababa HQ - Boardroom Terminal', '127.0.0.1/32', 'HW-EXEC-AA-01-SECURE', true),
    ('WS-SOC-01',  'Cybersecurity & SOC', 'Addis Ababa HQ - SOC Command Room Console 1', '127.0.0.1/32', 'HW-SOC-AA-01-SECURE', true),
    ('WS-SOC-02',  'Cybersecurity & SOC', 'Addis Ababa HQ - SOC Command Room Console 2', '127.0.0.1/32', 'HW-SOC-AA-02-SECURE', true),
    ('WS-FRAUD-01', 'Cybersecurity & SOC', 'Addis Ababa HQ - AML Investigation Lab', '127.0.0.1/32', 'HW-FRAUD-AA-01-SECURE', true),
    ('WS-FIN-01',  'Finance & Treasury', 'Addis Ababa HQ - Treasury Terminal A', '127.0.0.1/32', 'HW-FIN-AA-01-SECURE', true),
    ('WS-FIN-02',  'Finance & Treasury', 'Addis Ababa HQ - Escrow Reconciliation Desk', '127.0.0.1/32', 'HW-FIN-AA-02-SECURE', true),
    ('WS-OPS-01',  'Operations & KYC', 'Addis Ababa HQ - KYC Identity Inspection Terminal 1', '127.0.0.1/32', 'HW-OPS-AA-01-SECURE', true),
    ('WS-OPS-02',  'Operations & KYC', 'Addis Ababa HQ - KYC Identity Inspection Terminal 2', '127.0.0.1/32', 'HW-OPS-AA-02-SECURE', true),
    ('WS-OPS-03',  'Operations & KYC', 'Bole Regional Hub - Dispute Resolution Center', '127.0.0.1/32', 'HW-OPS-BOLE-03-SECURE', true),
    ('WS-MOD-01',  'Catalog & Moderation', 'Addis Ababa HQ - Content Moderation Station 1', '127.0.0.1/32', 'HW-MOD-AA-01-SECURE', true),
    ('WS-MOD-02',  'Catalog & Moderation', 'Hawassa Regional Facility - Moderation Desk 2', '127.0.0.1/32', 'HW-MOD-HAW-02-SECURE', true),
    ('WS-SUP-01',  'Customer Support', 'Addis Ababa HQ - Support Desk Lead Station', '127.0.0.1/32', 'HW-SUP-AA-01-SECURE', true),
    ('WS-SUP-02',  'Customer Support', 'Dire Dawa Regional Support Center', '127.0.0.1/32', 'HW-SUP-DD-02-SECURE', true)
ON CONFLICT (code) DO UPDATE
SET facility_location = EXCLUDED.facility_location, department = EXCLUDED.department;
