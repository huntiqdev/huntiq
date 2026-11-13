-- ============================================================================
-- 034_invitation_requests.sql
-- Table: invitation_requests (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for invitation_requests
ALTER TABLE invitation_requests ADD COLUMN IF NOT EXISTS email character varying(150) DEFAULT ''::character varying NOT NULL;
ALTER TABLE invitation_requests ADD COLUMN IF NOT EXISTS phone character varying(20) DEFAULT ''::character varying NOT NULL;
ALTER TABLE invitation_requests ADD COLUMN IF NOT EXISTS name character varying(100) DEFAULT ''::character varying NOT NULL;
ALTER TABLE invitation_requests ADD COLUMN IF NOT EXISTS message text;
ALTER TABLE invitation_requests ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'pending'::character varying NOT NULL;
ALTER TABLE invitation_requests ADD COLUMN IF NOT EXISTS invited_by integer;
ALTER TABLE invitation_requests ADD COLUMN IF NOT EXISTS approved_by integer;
ALTER TABLE invitation_requests ADD COLUMN IF NOT EXISTS approved_at timestamp with time zone;
ALTER TABLE invitation_requests ADD COLUMN IF NOT EXISTS invitation_code character varying(100);
ALTER TABLE invitation_requests ADD COLUMN IF NOT EXISTS used_at timestamp with time zone;
ALTER TABLE invitation_requests ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'invitation_requests' AND constraint_name = 'invitation_requests_invitation_code_key'
    ) THEN
        ALTER TABLE ONLY invitation_requests ADD CONSTRAINT invitation_requests_invitation_code_key UNIQUE (invitation_code);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'invitation_requests' AND constraint_name = 'invitation_requests_pkey'
    ) THEN
        ALTER TABLE ONLY invitation_requests ADD CONSTRAINT invitation_requests_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_invitation_code ON invitation_requests USING btree (invitation_code) WHERE (invitation_code IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_invitation_status ON invitation_requests USING btree (status, created_at DESC);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'invitation_requests' AND constraint_name = 'invitation_requests_approved_by_fkey'
    ) THEN
        ALTER TABLE ONLY invitation_requests ADD CONSTRAINT invitation_requests_approved_by_fkey FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'invitation_requests' AND constraint_name = 'invitation_requests_invited_by_fkey'
    ) THEN
        ALTER TABLE ONLY invitation_requests ADD CONSTRAINT invitation_requests_invited_by_fkey FOREIGN KEY (invited_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;
