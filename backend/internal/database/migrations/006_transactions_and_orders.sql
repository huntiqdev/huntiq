-- ============================================================
-- 006_transactions_and_orders.sql
-- CONSOLIDATED: transactions + order_requests
-- PURPOSE: All transaction and order tracking in one migration
-- NOTE: Uses ALTER TABLE pattern to add columns to shell tables
--        created in 001_schema.sql (same pattern as 002-005).
-- ============================================================

-- ── Transactions ────────────────────────────────────────────────────────────
ALTER TABLE transactions
    ADD COLUMN IF NOT EXISTS user_id              INT           NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS product_id           INT           REFERENCES products(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS amount               NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (amount > 0),
    ADD COLUMN IF NOT EXISTS fee_amount           NUMERIC(12,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS transaction_type     VARCHAR(50)   NOT NULL DEFAULT 'listing_fee',
    ADD COLUMN IF NOT EXISTS status               VARCHAR(20)   NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS payment_method       VARCHAR(50),
    ADD COLUMN IF NOT EXISTS payment_reference    VARCHAR(200),
    ADD COLUMN IF NOT EXISTS payer_phone          VARCHAR(20),
    ADD COLUMN IF NOT EXISTS notes                TEXT,
    -- Legacy compat columns used by Go Transaction model
    ADD COLUMN IF NOT EXISTS seller_id            INT           REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS buyer_name           VARCHAR(100),
    ADD COLUMN IF NOT EXISTS buyer_phone          VARCHAR(20),
    ADD COLUMN IF NOT EXISTS broker_fee           NUMERIC(12,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS seller_gets          NUMERIC(12,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Order Requests ──────────────────────────────────────────────────────────
ALTER TABLE order_requests
    ADD COLUMN IF NOT EXISTS product_id           INT           NOT NULL DEFAULT 0 REFERENCES products(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS seller_id            INT           NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS buyer_name           VARCHAR(100)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS buyer_phone          VARCHAR(20)   NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS buyer_email          VARCHAR(150),
    ADD COLUMN IF NOT EXISTS buyer_message        TEXT,
    ADD COLUMN IF NOT EXISTS message              TEXT,
    ADD COLUMN IF NOT EXISTS offer_price          NUMERIC(12,2),
    ADD COLUMN IF NOT EXISTS status               VARCHAR(20)   NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS seller_note          TEXT,
    ADD COLUMN IF NOT EXISTS buyer_tier_at_order  VARCHAR(20)   DEFAULT 'normal',
    ADD COLUMN IF NOT EXISTS created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Indexes ─────────────────────────────────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_transactions_user          ON transactions(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_product       ON transactions(product_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status        ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_type          ON transactions(transaction_type);
CREATE INDEX IF NOT EXISTS idx_transactions_seller        ON transactions(seller_id)          WHERE seller_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_order_requests_product     ON order_requests(product_id);
CREATE INDEX IF NOT EXISTS idx_order_requests_seller      ON order_requests(seller_id, status);
CREATE INDEX IF NOT EXISTS idx_order_requests_buyer_phone ON order_requests(buyer_phone);

-- ── Auto-update timestamps ──────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_order_requests_updated_at ON order_requests;
CREATE TRIGGER update_order_requests_updated_at BEFORE UPDATE ON order_requests
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
