-- ============================================================================
-- 035_referral_commissions.sql
-- Table: referral_commissions (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for referral_commissions
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS referrer_id integer DEFAULT 0 NOT NULL;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS referred_id integer DEFAULT 0 NOT NULL;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS transaction_id integer;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS commission_amount numeric(12,2) DEFAULT 0 NOT NULL;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS commission_percent numeric(5,2) DEFAULT 0 NOT NULL;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS commission_rate numeric(5,2) DEFAULT 0;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS transaction_type character varying(50) DEFAULT ''::character varying NOT NULL;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS transaction_amount numeric(12,2) DEFAULT 0;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'pending'::character varying NOT NULL;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS confirmed_at timestamp with time zone;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS paid_at timestamp with time zone;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS notes text;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS pending_until timestamp with time zone;

-- Table Constraints for referral_commissions
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'referral_commissions' AND constraint_name = 'referral_commissions_commission_amount_check'
    ) THEN
        ALTER TABLE referral_commissions ADD CONSTRAINT referral_commissions_commission_amount_check CHECK ((commission_amount >= (0)::numeric));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'referral_commissions' AND constraint_name = 'referral_commissions_pkey'
    ) THEN
        ALTER TABLE ONLY referral_commissions ADD CONSTRAINT referral_commissions_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'referral_commissions' AND constraint_name = 'uq_referral_commission_source'
    ) THEN
        ALTER TABLE ONLY referral_commissions ADD CONSTRAINT uq_referral_commission_source UNIQUE (referrer_id, transaction_id, transaction_type);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_referral_commissions_pending ON referral_commissions USING btree (status, created_at) WHERE ((status)::text = 'pending'::text);

CREATE INDEX IF NOT EXISTS idx_referral_commissions_pending_until ON referral_commissions USING btree (status, pending_until) WHERE ((status)::text = 'pending'::text);

CREATE INDEX IF NOT EXISTS idx_referral_commissions_referred ON referral_commissions USING btree (referred_id);

CREATE INDEX IF NOT EXISTS idx_referral_commissions_referrer ON referral_commissions USING btree (referrer_id, status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_referral_commissions_status ON referral_commissions USING btree (status);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'referral_commissions' AND constraint_name = 'referral_commissions_referred_id_fkey'
    ) THEN
        ALTER TABLE ONLY referral_commissions ADD CONSTRAINT referral_commissions_referred_id_fkey FOREIGN KEY (referred_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'referral_commissions' AND constraint_name = 'referral_commissions_referrer_id_fkey'
    ) THEN
        ALTER TABLE ONLY referral_commissions ADD CONSTRAINT referral_commissions_referrer_id_fkey FOREIGN KEY (referrer_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'referral_commissions' AND constraint_name = 'referral_commissions_transaction_id_fkey'
    ) THEN
        ALTER TABLE ONLY referral_commissions ADD CONSTRAINT referral_commissions_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE SET NULL;
    END IF;
END $$;
