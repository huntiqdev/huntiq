-- ============================================================================
-- 062_user_capabilities.sql
-- Table: user_capabilities (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for user_capabilities
ALTER TABLE user_capabilities ADD COLUMN IF NOT EXISTS user_id integer NOT NULL;
ALTER TABLE user_capabilities ADD COLUMN IF NOT EXISTS capability character varying(30) NOT NULL;
ALTER TABLE user_capabilities ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'active'::character varying NOT NULL;
ALTER TABLE user_capabilities ADD COLUMN IF NOT EXISTS verified_at timestamp with time zone;
ALTER TABLE user_capabilities ADD COLUMN IF NOT EXISTS activated_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE user_capabilities ADD COLUMN IF NOT EXISTS metadata jsonb DEFAULT '{}'::jsonb;

COMMENT ON TABLE user_capabilities IS 'Additive capabilities per user: buyer, seller, landlord, renter';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_capabilities' AND constraint_name = 'uq_user_capability'
    ) THEN
        ALTER TABLE ONLY user_capabilities ADD CONSTRAINT uq_user_capability UNIQUE (user_id, capability);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_capabilities' AND constraint_name = 'user_capabilities_pkey'
    ) THEN
        ALTER TABLE ONLY user_capabilities ADD CONSTRAINT user_capabilities_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_user_capabilities_active ON user_capabilities USING btree (user_id, capability) WHERE ((status)::text = 'active'::text);

CREATE INDEX IF NOT EXISTS idx_user_capabilities_user_id ON user_capabilities USING btree (user_id);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_capabilities' AND constraint_name = 'user_capabilities_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY user_capabilities ADD CONSTRAINT user_capabilities_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
