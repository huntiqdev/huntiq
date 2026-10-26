-- ============================================================================
-- Migration 033: Admin Notifications Type and Read Tracking
-- ============================================================================
ALTER TABLE admin_notifications
    ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'info',
    ADD COLUMN IF NOT EXISTS read_by INT[] NOT NULL DEFAULT '{}';

CREATE INDEX IF NOT EXISTS idx_admin_notifications_type 
    ON admin_notifications(type);

CREATE INDEX IF NOT EXISTS idx_admin_notifications_read_by 
    ON admin_notifications USING GIN (read_by);
