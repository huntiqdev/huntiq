-- ============================================================
-- 007_notifications_system.sql
-- CONSOLIDATED: notifications + admin_notifications
-- PURPOSE: Unified notification system for all users
-- NOTE: Uses ALTER TABLE to extend shell tables from 001_schema.sql
-- ============================================================

-- ── Notifications ───────────────────────────────────────────────────────────
ALTER TABLE notifications
    ADD COLUMN IF NOT EXISTS user_id    INT           NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS type       VARCHAR(50)   NOT NULL DEFAULT 'system_message',
    ADD COLUMN IF NOT EXISTS title      VARCHAR(200)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS message    TEXT          NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS link       VARCHAR(500),
    ADD COLUMN IF NOT EXISTS is_read    BOOLEAN       NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_admin   BOOLEAN       NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS priority   VARCHAR(20)   DEFAULT 'normal',
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Admin Notifications (extends shell from 001) ────────────────────────────
ALTER TABLE admin_notifications
    ADD COLUMN IF NOT EXISTS title      VARCHAR(200)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS message    TEXT          NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS priority   VARCHAR(20)   DEFAULT 'normal',
    ADD COLUMN IF NOT EXISTS target_role VARCHAR(50),
    ADD COLUMN IF NOT EXISTS is_read    BOOLEAN       NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Indexes ─────────────────────────────────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_notifications_user   ON notifications(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;
CREATE INDEX IF NOT EXISTS idx_notifications_admin  ON notifications(is_admin, created_at DESC) WHERE is_admin = TRUE;
CREATE INDEX IF NOT EXISTS idx_notifications_type   ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_admin_notifications_priority ON admin_notifications(priority, created_at DESC);

-- Auto-cleanup old read notifications (keep last 90 days)
CREATE OR REPLACE FUNCTION cleanup_old_notifications()
RETURNS void AS $$
BEGIN
    DELETE FROM notifications
    WHERE is_read = TRUE
      AND created_at < NOW() - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;
