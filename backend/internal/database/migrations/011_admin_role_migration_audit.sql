-- ============================================================================
-- 011_admin_role_migration_audit.sql
-- Table: admin_role_migration_audit (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for admin_role_migration_audit
ALTER TABLE admin_role_migration_audit ADD COLUMN IF NOT EXISTS user_id integer NOT NULL;
ALTER TABLE admin_role_migration_audit ADD COLUMN IF NOT EXISTS old_role character varying(50) NOT NULL;
ALTER TABLE admin_role_migration_audit ADD COLUMN IF NOT EXISTS new_role character varying(50) NOT NULL;
ALTER TABLE admin_role_migration_audit ADD COLUMN IF NOT EXISTS migrated_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE admin_role_migration_audit ADD COLUMN IF NOT EXISTS notes text;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_role_migration_audit' AND constraint_name = 'admin_role_migration_audit_pkey'
    ) THEN
        ALTER TABLE ONLY admin_role_migration_audit ADD CONSTRAINT admin_role_migration_audit_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_role_migration_audit' AND constraint_name = 'admin_role_migration_audit_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY admin_role_migration_audit ADD CONSTRAINT admin_role_migration_audit_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
