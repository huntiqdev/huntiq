-- ============================================================================
-- 009_admin_permissions.sql
-- Table: admin_permissions (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for admin_permissions
ALTER TABLE admin_permissions ADD COLUMN IF NOT EXISTS enabled boolean DEFAULT true NOT NULL;
ALTER TABLE admin_permissions ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE admin_permissions ADD COLUMN IF NOT EXISTS granted_by integer;
ALTER TABLE admin_permissions ADD COLUMN IF NOT EXISTS granted_at timestamp with time zone DEFAULT now() NOT NULL;

COMMENT ON TABLE admin_permissions IS 'Granular admin permissions (now includes: operations, seller_verification, trust_safety, refunds, appeals)';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_permissions' AND constraint_name = 'admin_permissions_pkey'
    ) THEN
        ALTER TABLE ONLY admin_permissions ADD CONSTRAINT admin_permissions_pkey PRIMARY KEY (admin_id, permission);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_admin_permissions_admin ON admin_permissions USING btree (admin_id);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_permissions' AND constraint_name = 'admin_permissions_granted_by_fkey'
    ) THEN
        ALTER TABLE ONLY admin_permissions ADD CONSTRAINT admin_permissions_granted_by_fkey FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;
