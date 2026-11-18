-- ============================================================================
-- 003_account_appeals.sql
-- Table: account_appeals (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for account_appeals
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS user_id integer NOT NULL;
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS appeal_type character varying(50) NOT NULL;
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS related_product_id integer;
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS related_document_id integer;
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS reason text NOT NULL;
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS evidence_urls text[];
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'pending'::character varying NOT NULL;
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS reviewed_by integer;
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS reviewed_at timestamp with time zone;
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS decision text;
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS action_taken text;
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE account_appeals ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

COMMENT ON TABLE account_appeals IS 'User appeals for account actions and listing decisions';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'account_appeals' AND constraint_name = 'account_appeals_pkey'
    ) THEN
        ALTER TABLE ONLY account_appeals ADD CONSTRAINT account_appeals_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_appeals_reviewer ON account_appeals USING btree (reviewed_by) WHERE (reviewed_by IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_appeals_status ON account_appeals USING btree (status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_appeals_type ON account_appeals USING btree (appeal_type, status);

CREATE INDEX IF NOT EXISTS idx_appeals_user ON account_appeals USING btree (user_id, status);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'account_appeals' AND constraint_name = 'account_appeals_related_document_id_fkey'
    ) THEN
        ALTER TABLE ONLY account_appeals ADD CONSTRAINT account_appeals_related_document_id_fkey FOREIGN KEY (related_document_id) REFERENCES seller_verification_documents(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'account_appeals' AND constraint_name = 'account_appeals_related_product_id_fkey'
    ) THEN
        ALTER TABLE ONLY account_appeals ADD CONSTRAINT account_appeals_related_product_id_fkey FOREIGN KEY (related_product_id) REFERENCES products(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'account_appeals' AND constraint_name = 'account_appeals_reviewed_by_fkey'
    ) THEN
        ALTER TABLE ONLY account_appeals ADD CONSTRAINT account_appeals_reviewed_by_fkey FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'account_appeals' AND constraint_name = 'account_appeals_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY account_appeals ADD CONSTRAINT account_appeals_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
