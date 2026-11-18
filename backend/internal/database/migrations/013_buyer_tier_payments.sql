-- ============================================================================
-- 013_buyer_tier_payments.sql
-- Table: buyer_tier_payments (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for buyer_tier_payments
ALTER TABLE buyer_tier_payments ADD COLUMN IF NOT EXISTS buyer_id integer DEFAULT 0 NOT NULL;
ALTER TABLE buyer_tier_payments ADD COLUMN IF NOT EXISTS tier character varying(20) DEFAULT 'gold'::character varying NOT NULL;
ALTER TABLE buyer_tier_payments ADD COLUMN IF NOT EXISTS amount numeric(12,2) DEFAULT 0 NOT NULL;
ALTER TABLE buyer_tier_payments ADD COLUMN IF NOT EXISTS reference character varying(200);
ALTER TABLE buyer_tier_payments ADD COLUMN IF NOT EXISTS payer_phone character varying(20);
ALTER TABLE buyer_tier_payments ADD COLUMN IF NOT EXISTS payment_method character varying(50);
ALTER TABLE buyer_tier_payments ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'buyer_tier_payments' AND constraint_name = 'buyer_tier_payments_pkey'
    ) THEN
        ALTER TABLE ONLY buyer_tier_payments ADD CONSTRAINT buyer_tier_payments_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_buyer_tier_payments_buyer ON buyer_tier_payments USING btree (buyer_id, created_at DESC);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'buyer_tier_payments' AND constraint_name = 'buyer_tier_payments_buyer_id_fkey'
    ) THEN
        ALTER TABLE ONLY buyer_tier_payments ADD CONSTRAINT buyer_tier_payments_buyer_id_fkey FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
