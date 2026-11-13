-- ============================================================================
-- 052_refund_requests.sql
-- Table: refund_requests (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for refund_requests
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS transaction_id integer;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS product_id integer;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS requester_id integer NOT NULL;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS requester_role character varying(20) NOT NULL;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS amount numeric(12,2) NOT NULL;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS reason character varying(100) NOT NULL;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS description text NOT NULL;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS evidence_urls text[];
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'pending'::character varying NOT NULL;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS reviewed_by integer;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS reviewed_at timestamp with time zone;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS reject_reason text;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS refund_reference character varying(100);
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS processed_at timestamp with time zone;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE refund_requests ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

-- Table Constraints for refund_requests
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'refund_requests' AND constraint_name = 'refund_requests_amount_check'
    ) THEN
        ALTER TABLE refund_requests ADD CONSTRAINT refund_requests_amount_check CHECK ((amount > (0)::numeric));
    END IF;
END $$;

COMMENT ON TABLE refund_requests IS 'Refund and dispute requests for transactions';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'refund_requests' AND constraint_name = 'refund_requests_pkey'
    ) THEN
        ALTER TABLE ONLY refund_requests ADD CONSTRAINT refund_requests_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_refund_requests_requester ON refund_requests USING btree (requester_id, status);

CREATE INDEX IF NOT EXISTS idx_refund_requests_reviewer ON refund_requests USING btree (reviewed_by) WHERE (reviewed_by IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_refund_requests_status ON refund_requests USING btree (status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_refund_requests_transaction ON refund_requests USING btree (transaction_id) WHERE (transaction_id IS NOT NULL);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'refund_requests' AND constraint_name = 'refund_requests_product_id_fkey'
    ) THEN
        ALTER TABLE ONLY refund_requests ADD CONSTRAINT refund_requests_product_id_fkey FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'refund_requests' AND constraint_name = 'refund_requests_requester_id_fkey'
    ) THEN
        ALTER TABLE ONLY refund_requests ADD CONSTRAINT refund_requests_requester_id_fkey FOREIGN KEY (requester_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'refund_requests' AND constraint_name = 'refund_requests_reviewed_by_fkey'
    ) THEN
        ALTER TABLE ONLY refund_requests ADD CONSTRAINT refund_requests_reviewed_by_fkey FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'refund_requests' AND constraint_name = 'refund_requests_transaction_id_fkey'
    ) THEN
        ALTER TABLE ONLY refund_requests ADD CONSTRAINT refund_requests_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE;
    END IF;
END $$;
