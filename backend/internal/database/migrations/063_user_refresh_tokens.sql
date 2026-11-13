-- ============================================================================
-- 063_user_refresh_tokens.sql
-- Table: user_refresh_tokens (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for user_refresh_tokens
ALTER TABLE user_refresh_tokens ADD COLUMN IF NOT EXISTS user_id integer NOT NULL;
ALTER TABLE user_refresh_tokens ADD COLUMN IF NOT EXISTS family_id uuid NOT NULL;
ALTER TABLE user_refresh_tokens ADD COLUMN IF NOT EXISTS token_hash character varying(64) NOT NULL;
ALTER TABLE user_refresh_tokens ADD COLUMN IF NOT EXISTS parent_token_hash character varying(64);
ALTER TABLE user_refresh_tokens ADD COLUMN IF NOT EXISTS is_revoked boolean DEFAULT false NOT NULL;
ALTER TABLE user_refresh_tokens ADD COLUMN IF NOT EXISTS revocation_reason character varying(100);
ALTER TABLE user_refresh_tokens ADD COLUMN IF NOT EXISTS revoked_at timestamp with time zone;
ALTER TABLE user_refresh_tokens ADD COLUMN IF NOT EXISTS expires_at timestamp with time zone NOT NULL;
ALTER TABLE user_refresh_tokens ADD COLUMN IF NOT EXISTS user_agent text;
ALTER TABLE user_refresh_tokens ADD COLUMN IF NOT EXISTS ip_address character varying(50);
ALTER TABLE user_refresh_tokens ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE user_refresh_tokens ADD COLUMN IF NOT EXISTS last_used_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_refresh_tokens' AND constraint_name = 'user_refresh_tokens_pkey'
    ) THEN
        ALTER TABLE ONLY user_refresh_tokens ADD CONSTRAINT user_refresh_tokens_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_refresh_tokens' AND constraint_name = 'user_refresh_tokens_token_hash_key'
    ) THEN
        ALTER TABLE ONLY user_refresh_tokens ADD CONSTRAINT user_refresh_tokens_token_hash_key UNIQUE (token_hash);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_user_refresh_tokens_expiry ON user_refresh_tokens USING btree (expires_at);

CREATE INDEX IF NOT EXISTS idx_user_refresh_tokens_family ON user_refresh_tokens USING btree (family_id);

CREATE INDEX IF NOT EXISTS idx_user_refresh_tokens_hash ON user_refresh_tokens USING btree (token_hash) WHERE (is_revoked IS FALSE);

CREATE INDEX IF NOT EXISTS idx_user_refresh_tokens_user ON user_refresh_tokens USING btree (user_id) WHERE (is_revoked IS FALSE);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_refresh_tokens' AND constraint_name = 'user_refresh_tokens_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY user_refresh_tokens ADD CONSTRAINT user_refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
