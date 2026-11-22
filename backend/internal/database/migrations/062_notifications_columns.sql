-- ============================================================================
-- 062_notifications_columns.sql
-- Add missing columns to notifications table required by ListNotifications & MarkNotificationRead
-- ============================================================================

ALTER TABLE notifications
    ADD COLUMN IF NOT EXISTS user_id INT REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS title VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS message TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'system',
    ADD COLUMN IF NOT EXISTS is_read BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id, created_at DESC);
