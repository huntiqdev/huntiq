-- ============================================================================
-- 010_admin_reauth_tokens.sql
-- Table: admin_reauth_tokens (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for admin_reauth_tokens
ALTER TABLE admin_reauth_tokens ADD COLUMN IF NOT EXISTS user_id integer NOT NULL;
ALTER TABLE admin_reauth_tokens ADD COLUMN IF NOT EXISTS token_hash character varying(64) NOT NULL;
ALTER TABLE admin_reauth_tokens ADD COLUMN IF NOT EXISTS action_scope character varying(100) NOT NULL;
ALTER TABLE admin_reauth_tokens ADD COLUMN IF NOT EXISTS expires_at timestamp with time zone NOT NULL;
ALTER TABLE admin_reauth_tokens ADD COLUMN IF NOT EXISTS used_at timestamp with time zone;
ALTER TABLE admin_reauth_tokens ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_reauth_tokens' AND constraint_name = 'admin_reauth_tokens_pkey'
    ) THEN
        ALTER TABLE ONLY admin_reauth_tokens ADD CONSTRAINT admin_reauth_tokens_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_reauth_tokens' AND constraint_name = 'admin_reauth_tokens_token_hash_key'
    ) THEN
        ALTER TABLE ONLY admin_reauth_tokens ADD CONSTRAINT admin_reauth_tokens_token_hash_key UNIQUE (token_hash);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_admin_reauth_tokens_lookup ON admin_reauth_tokens USING btree (token_hash) WHERE (used_at IS NULL);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_reauth_tokens' AND constraint_name = 'admin_reauth_tokens_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY admin_reauth_tokens ADD CONSTRAINT admin_reauth_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
