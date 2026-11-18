-- ============================================================================
-- 008_admin_passkeys.sql
-- Table: admin_passkeys (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for admin_passkeys
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS user_id integer NOT NULL;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS credential_id bytea NOT NULL;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS public_key bytea NOT NULL;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS attestation_type character varying(50) DEFAULT 'none'::character varying NOT NULL;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS aaguid bytea;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS sign_count bigint DEFAULT 0 NOT NULL;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS clone_warning boolean DEFAULT false NOT NULL;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS attachment character varying(20) DEFAULT 'platform'::character varying;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS name character varying(100) DEFAULT 'Passkey'::character varying NOT NULL;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS transports text[] DEFAULT '{}'::text[];
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS backup_eligible boolean DEFAULT false;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS backup_state boolean DEFAULT false;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS last_used_at timestamp with time zone;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS revoked_at timestamp with time zone;
ALTER TABLE admin_passkeys ADD COLUMN IF NOT EXISTS revocation_reason text;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_passkeys' AND constraint_name = 'admin_passkeys_credential_id_key'
    ) THEN
        ALTER TABLE ONLY admin_passkeys ADD CONSTRAINT admin_passkeys_credential_id_key UNIQUE (credential_id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_passkeys' AND constraint_name = 'admin_passkeys_pkey'
    ) THEN
        ALTER TABLE ONLY admin_passkeys ADD CONSTRAINT admin_passkeys_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_admin_passkeys_cred_id ON admin_passkeys USING btree (credential_id);

CREATE INDEX IF NOT EXISTS idx_admin_passkeys_user ON admin_passkeys USING btree (user_id) WHERE (revoked_at IS NULL);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_passkeys' AND constraint_name = 'admin_passkeys_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY admin_passkeys ADD CONSTRAINT admin_passkeys_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
