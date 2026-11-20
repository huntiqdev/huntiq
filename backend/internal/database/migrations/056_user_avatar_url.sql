-- ============================================================================
-- 056_user_avatar_url.sql
-- Table: users (Add avatar_url column for circular profile photo support)
-- ============================================================================

ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url text;
