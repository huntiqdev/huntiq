-- ============================================================================
-- 043_upgrade_requests.sql
-- Table: upgrade_requests (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for upgrade_requests
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS user_id integer DEFAULT 0 NOT NULL;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS role character varying(20) DEFAULT 'seller'::character varying NOT NULL;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS from_tier character varying(20);
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS requested_tier character varying(20) DEFAULT 'gold'::character varying NOT NULL;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS amount numeric(12,2) DEFAULT 0 NOT NULL;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS duration_days integer;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS payment_method character varying(50) DEFAULT ''::character varying NOT NULL;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS reference character varying(200) DEFAULT ''::character varying NOT NULL;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS payer_phone character varying(20) DEFAULT ''::character varying NOT NULL;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'pending'::character varying NOT NULL;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS reviewed_by integer;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS reviewed_at timestamp with time zone;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS reject_reason text;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS notes text;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE upgrade_requests ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

-- Table Constraints for upgrade_requests
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'upgrade_requests' AND constraint_name = 'upgrade_requests_amount_check'
    ) THEN
        ALTER TABLE upgrade_requests ADD CONSTRAINT upgrade_requests_amount_check CHECK ((amount > (0)::numeric));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'upgrade_requests' AND constraint_name = 'upgrade_requests_pkey'
    ) THEN
        ALTER TABLE ONLY upgrade_requests ADD CONSTRAINT upgrade_requests_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_upgrade_requests_pending ON upgrade_requests USING btree (status) WHERE ((status)::text = 'pending'::text);

CREATE INDEX IF NOT EXISTS idx_upgrade_requests_role ON upgrade_requests USING btree (role);

CREATE INDEX IF NOT EXISTS idx_upgrade_requests_status ON upgrade_requests USING btree (status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_upgrade_requests_user ON upgrade_requests USING btree (user_id, status, created_at DESC);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'upgrade_requests' AND constraint_name = 'upgrade_requests_reviewed_by_fkey'
    ) THEN
        ALTER TABLE ONLY upgrade_requests ADD CONSTRAINT upgrade_requests_reviewed_by_fkey FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'upgrade_requests' AND constraint_name = 'upgrade_requests_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY upgrade_requests ADD CONSTRAINT upgrade_requests_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
