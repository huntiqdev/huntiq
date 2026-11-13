-- ============================================================================
-- 051_referral_withdrawals.sql
-- Table: referral_withdrawals (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for referral_withdrawals
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS user_id integer DEFAULT 0 NOT NULL;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS amount numeric(12,2) DEFAULT 0 NOT NULL;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS withdrawal_method character varying(50) DEFAULT ''::character varying NOT NULL;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS payment_method character varying(50);
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS account_info character varying(200) DEFAULT ''::character varying NOT NULL;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS payer_phone character varying(20);
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS payment_phone character varying(20);
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS reference character varying(200);
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'pending'::character varying NOT NULL;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS approved_by integer;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS reviewed_by integer;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS approved_at timestamp with time zone;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS reviewed_at timestamp with time zone;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS paid_at timestamp with time zone;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS rejection_reason text;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS rejection_note text;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS admin_note text;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS notes text;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE referral_withdrawals ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

-- Table Constraints for referral_withdrawals
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'referral_withdrawals' AND constraint_name = 'referral_withdrawals_amount_check'
    ) THEN
        ALTER TABLE referral_withdrawals ADD CONSTRAINT referral_withdrawals_amount_check CHECK ((amount > (0)::numeric));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'referral_withdrawals' AND constraint_name = 'referral_withdrawals_pkey'
    ) THEN
        ALTER TABLE ONLY referral_withdrawals ADD CONSTRAINT referral_withdrawals_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_referral_withdrawals_status ON referral_withdrawals USING btree (status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_referral_withdrawals_user ON referral_withdrawals USING btree (user_id, status, created_at DESC);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'referral_withdrawals' AND constraint_name = 'referral_withdrawals_approved_by_fkey'
    ) THEN
        ALTER TABLE ONLY referral_withdrawals ADD CONSTRAINT referral_withdrawals_approved_by_fkey FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'referral_withdrawals' AND constraint_name = 'referral_withdrawals_reviewed_by_fkey'
    ) THEN
        ALTER TABLE ONLY referral_withdrawals ADD CONSTRAINT referral_withdrawals_reviewed_by_fkey FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'referral_withdrawals' AND constraint_name = 'referral_withdrawals_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY referral_withdrawals ADD CONSTRAINT referral_withdrawals_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
