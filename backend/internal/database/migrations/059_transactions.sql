-- ============================================================================
-- 059_transactions.sql
-- Table: transactions (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for transactions
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS user_id integer DEFAULT 0 NOT NULL;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS product_id integer;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS amount numeric(12,2) DEFAULT 0 NOT NULL;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS fee_amount numeric(12,2) DEFAULT 0;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS transaction_type character varying(50) DEFAULT 'listing_fee'::character varying NOT NULL;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'pending'::character varying NOT NULL;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS payment_method character varying(50);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS payment_reference character varying(200);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS payer_phone character varying(20);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS notes text;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS seller_id integer;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS buyer_name character varying(100);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS buyer_phone character varying(20);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS broker_fee numeric(12,2) DEFAULT 0;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS seller_gets numeric(12,2) DEFAULT 0;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

-- Table Constraints for transactions
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'transactions' AND constraint_name = 'transactions_amount_check'
    ) THEN
        ALTER TABLE transactions ADD CONSTRAINT transactions_amount_check CHECK ((amount > (0)::numeric));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'transactions' AND constraint_name = 'transactions_pkey'
    ) THEN
        ALTER TABLE ONLY transactions ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_transactions_product ON transactions USING btree (product_id);

CREATE INDEX IF NOT EXISTS idx_transactions_seller ON transactions USING btree (seller_id) WHERE (seller_id IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions USING btree (status);

CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions USING btree (transaction_type);

CREATE INDEX IF NOT EXISTS idx_transactions_user ON transactions USING btree (user_id, created_at DESC);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'transactions' AND constraint_name = 'transactions_product_id_fkey'
    ) THEN
        ALTER TABLE ONLY transactions ADD CONSTRAINT transactions_product_id_fkey FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'transactions' AND constraint_name = 'transactions_seller_id_fkey'
    ) THEN
        ALTER TABLE ONLY transactions ADD CONSTRAINT transactions_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'transactions' AND constraint_name = 'transactions_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY transactions ADD CONSTRAINT transactions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
