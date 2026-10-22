-- ============================================================
-- 008_admin_system.sql
-- CONSOLIDATED: admin_logs + admin_permissions + hero_banners
-- PURPOSE: Complete admin system in one migration
-- NOTE: Uses ALTER TABLE to extend shell tables from 001_schema.sql.
--        admin_permissions gets 'enabled' and 'updated_at' columns
--        that the Go code expects.
-- ============================================================

-- ── Admin Logs ──────────────────────────────────────────────────────────────
ALTER TABLE admin_logs
    ADD COLUMN IF NOT EXISTS admin_id       INT           NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS action         VARCHAR(100)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS target_type    VARCHAR(50),
    ADD COLUMN IF NOT EXISTS target_id      INT,
    ADD COLUMN IF NOT EXISTS product_id     INT,
    ADD COLUMN IF NOT EXISTS target_user_id INT,
    ADD COLUMN IF NOT EXISTS notes          TEXT,
    ADD COLUMN IF NOT EXISTS details        TEXT,
    ADD COLUMN IF NOT EXISTS ip_address     VARCHAR(50),
    ADD COLUMN IF NOT EXISTS user_agent     TEXT,
    ADD COLUMN IF NOT EXISTS created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Admin Permissions ───────────────────────────────────────────────────────
-- 001_schema.sql creates: (admin_id INT, permission VARCHAR(50), PRIMARY KEY(admin_id, permission))
-- Go code needs: enabled BOOLEAN, updated_at TIMESTAMPTZ, granted_by INT
ALTER TABLE admin_permissions
    ADD COLUMN IF NOT EXISTS enabled    BOOLEAN       NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS granted_by INT           REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS granted_at TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Hero Banners ────────────────────────────────────────────────────────────
-- Go HeroBanner model fields: BannerType, MediaType, MediaURL, Title, Subtitle,
-- LinkURL, PromotionText, TargetAudience, StartsAt, EndsAt, SortOrder, IsActive
ALTER TABLE hero_banners
    ADD COLUMN IF NOT EXISTS banner_type     VARCHAR(20)   NOT NULL DEFAULT 'seller',
    ADD COLUMN IF NOT EXISTS media_type      VARCHAR(20)   NOT NULL DEFAULT 'image',
    ADD COLUMN IF NOT EXISTS media_url       TEXT          NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS title           VARCHAR(200)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS subtitle        TEXT,
    ADD COLUMN IF NOT EXISTS link_url        TEXT,
    ADD COLUMN IF NOT EXISTS promotion_text  TEXT,
    ADD COLUMN IF NOT EXISTS target_audience VARCHAR(50)   NOT NULL DEFAULT 'all',
    ADD COLUMN IF NOT EXISTS sort_order      INT           NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS is_active       BOOLEAN       NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS starts_at       TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS ends_at         TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS created_by      INT           REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Indexes ─────────────────────────────────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_admin_logs_admin       ON admin_logs(admin_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_logs_action      ON admin_logs(action, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_logs_target      ON admin_logs(target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_admin_logs_target_user ON admin_logs(target_user_id)        WHERE target_user_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_admin_logs_product     ON admin_logs(product_id)             WHERE product_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_admin_permissions_admin ON admin_permissions(admin_id);
CREATE INDEX IF NOT EXISTS idx_hero_banners_active    ON hero_banners(is_active, sort_order) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_hero_banners_dates     ON hero_banners(starts_at, ends_at)    WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_hero_banners_type      ON hero_banners(banner_type)           WHERE is_active = TRUE;

-- ── Trigger for banner timestamps ───────────────────────────────────────────
DROP TRIGGER IF EXISTS update_hero_banners_updated_at ON hero_banners;
CREATE TRIGGER update_hero_banners_updated_at BEFORE UPDATE ON hero_banners
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
