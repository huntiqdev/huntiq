-- ============================================================
-- 010_tier_and_upgrades.sql
-- CONSOLIDATED: upgrade_requests + user_tier_history
-- PURPOSE: Tier management and upgrade tracking
-- NOTE: Uses ALTER TABLE to extend shell tables from 001_schema.sql.
--        Column names ALIGNED with Go admin_repo.go queries:
--        - 'role' not 'user_role'
--        - 'requested_tier' not 'to_tier'
--        - 'reference' not 'payment_reference'
--        - 'reject_reason' not 'rejection_reason'
--        - 'tier_type' not 'user_role' (in user_tier_history)
-- ============================================================

-- ── Upgrade Requests ────────────────────────────────────────────────────────
ALTER TABLE upgrade_requests
    ADD COLUMN IF NOT EXISTS user_id          INT           NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS role             VARCHAR(20)   NOT NULL DEFAULT 'seller',
    ADD COLUMN IF NOT EXISTS from_tier        VARCHAR(20),
    ADD COLUMN IF NOT EXISTS requested_tier   VARCHAR(20)   NOT NULL DEFAULT 'gold',
    ADD COLUMN IF NOT EXISTS amount           NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (amount > 0),
    ADD COLUMN IF NOT EXISTS duration_days    INT,
    ADD COLUMN IF NOT EXISTS payment_method   VARCHAR(50)   NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS reference        VARCHAR(200)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS payer_phone      VARCHAR(20)   NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS status           VARCHAR(20)   NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS reviewed_by      INT           REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS reviewed_at      TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS reject_reason    TEXT,
    ADD COLUMN IF NOT EXISTS notes            TEXT,
    ADD COLUMN IF NOT EXISTS created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── User Tier History (audit trail) ─────────────────────────────────────────
-- Go code uses 'tier_type' column name (not 'user_role')
ALTER TABLE user_tier_history
    ADD COLUMN IF NOT EXISTS user_id            INT           NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS tier_type          VARCHAR(20)   NOT NULL DEFAULT 'seller',
    ADD COLUMN IF NOT EXISTS old_tier           VARCHAR(20)   NOT NULL DEFAULT 'normal',
    ADD COLUMN IF NOT EXISTS new_tier           VARCHAR(20)   NOT NULL DEFAULT 'normal',
    ADD COLUMN IF NOT EXISTS reason             VARCHAR(100)  NOT NULL DEFAULT 'system',
    ADD COLUMN IF NOT EXISTS changed_by         INT           REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS upgrade_request_id INT           REFERENCES upgrade_requests(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS duration_days      INT,
    ADD COLUMN IF NOT EXISTS expires_at         TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS created_at         TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Indexes ─────────────────────────────────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_upgrade_requests_user    ON upgrade_requests(user_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_upgrade_requests_status  ON upgrade_requests(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_upgrade_requests_pending ON upgrade_requests(status) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_upgrade_requests_role    ON upgrade_requests(role);
CREATE INDEX IF NOT EXISTS idx_tier_history_user        ON user_tier_history(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tier_history_type        ON user_tier_history(tier_type, created_at DESC);

-- ── Triggers ────────────────────────────────────────────────────────────────
DROP TRIGGER IF EXISTS update_upgrade_requests_updated_at ON upgrade_requests;
CREATE TRIGGER update_upgrade_requests_updated_at BEFORE UPDATE ON upgrade_requests
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to record tier changes automatically
CREATE OR REPLACE FUNCTION record_tier_change()
RETURNS TRIGGER AS $$
BEGIN
    IF (OLD.tier IS DISTINCT FROM NEW.tier) OR (OLD.buyer_tier IS DISTINCT FROM NEW.buyer_tier) THEN
        INSERT INTO user_tier_history (user_id, tier_type, old_tier, new_tier, reason)
        VALUES (
            NEW.id,
            CASE WHEN OLD.tier IS DISTINCT FROM NEW.tier THEN 'seller' ELSE 'buyer' END,
            COALESCE(CASE WHEN OLD.tier IS DISTINCT FROM NEW.tier THEN OLD.tier::text ELSE OLD.buyer_tier::text END, 'normal'),
            COALESCE(CASE WHEN OLD.tier IS DISTINCT FROM NEW.tier THEN NEW.tier::text ELSE NEW.buyer_tier::text END, 'normal'),
            'system'
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS track_tier_changes ON users;
CREATE TRIGGER track_tier_changes AFTER UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION record_tier_change();
