-- ============================================================
-- 002_users.sql
-- PURPOSE : Single table for ALL user roles: seller, buyer, admin.
--           Role column discriminates. Tier/access columns are
--           NULL when not applicable to the role.
-- CHILDREN: seller_phones — merged at bottom of this file.
-- NOTE    : All columns already exist in 001_schema.sql for fresh
--           installs. This file uses IF NOT EXISTS / ON CONFLICT
--           so it is safe to re-run on existing databases.
-- ============================================================

ALTER TABLE users
    -- ── Identity ─────────────────────────────────────────────
    ADD COLUMN IF NOT EXISTS name                         VARCHAR(100)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS phone                        VARCHAR(20)   NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS business_name                VARCHAR(150),
    ADD COLUMN IF NOT EXISTS email                        VARCHAR(150),
    ADD COLUMN IF NOT EXISTS password_hash                TEXT          NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS rating                       NUMERIC(3,2)  NOT NULL DEFAULT 0.00,
    ADD COLUMN IF NOT EXISTS role                         VARCHAR(20)   NOT NULL DEFAULT 'seller',
    ADD COLUMN IF NOT EXISTS is_active                    BOOLEAN       NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS deactivation_reason          TEXT,

    -- ── Seller fields ─────────────────────────────────────────
    ADD COLUMN IF NOT EXISTS tier                         seller_tier   NOT NULL DEFAULT 'normal',
    ADD COLUMN IF NOT EXISTS total_value_sold             NUMERIC(15,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS additional_phones            TEXT[]        DEFAULT NULL,

    -- ── Buyer fields ──────────────────────────────────────────
    ADD COLUMN IF NOT EXISTS buyer_tier                   buyer_tier    NOT NULL DEFAULT 'normal',
    ADD COLUMN IF NOT EXISTS buyer_tier_paid_at           TIMESTAMPTZ   DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS buyer_tier_expires_at        TIMESTAMPTZ   DEFAULT NULL,

    -- ── Admin fields ──────────────────────────────────────────
    ADD COLUMN IF NOT EXISTS admin_role                   admin_role    DEFAULT NULL,

    -- ── Auth / Security ───────────────────────────────────────
    ADD COLUMN IF NOT EXISTS otp_verified                 BOOLEAN       NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_banned                    BOOLEAN       NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS last_login_at                TIMESTAMPTZ   DEFAULT NULL,

    -- ── Referral ──────────────────────────────────────────────
    ADD COLUMN IF NOT EXISTS invited_by                   INT           REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS invitation_accepted_at       TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS referral_suspended           BOOLEAN       NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS referral_permanently_blocked BOOLEAN       NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS referral_block_reason        TEXT,

    -- ── Timestamps ────────────────────────────────────────────
    ADD COLUMN IF NOT EXISTS created_at                   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at                   TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- Unique constraint: same phone can register as both seller and buyer
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name='users' AND constraint_name='users_phone_role_key'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT users_phone_role_key UNIQUE (phone, role);
    END IF;
END $$;

-- referral_code column — needs special handling (UNIQUE + NOT NULL, no default)
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='users' AND column_name='referral_code'
    ) THEN
        ALTER TABLE users ADD COLUMN referral_code VARCHAR(11) UNIQUE NOT NULL DEFAULT 'PLACEHOLDER';
        ALTER TABLE users ALTER COLUMN referral_code DROP DEFAULT;
    END IF;
END $$;

-- ── Indexes (all idempotent) ──────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_users_email            ON users (email)                           WHERE email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_phone            ON users (phone);
CREATE INDEX IF NOT EXISTS idx_users_phone_role       ON users (phone, role);
CREATE INDEX IF NOT EXISTS idx_users_role             ON users (role);
CREATE INDEX IF NOT EXISTS idx_users_tier             ON users (tier);
CREATE INDEX IF NOT EXISTS idx_users_total_value_sold ON users (total_value_sold);
CREATE INDEX IF NOT EXISTS idx_users_last_login       ON users (last_login_at DESC)              WHERE last_login_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_invited_by       ON users (invited_by)                      WHERE invited_by IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_ref_suspended    ON users (referral_suspended)              WHERE referral_suspended = TRUE;
CREATE INDEX IF NOT EXISTS idx_users_ref_blocked      ON users (referral_permanently_blocked)    WHERE referral_permanently_blocked = TRUE;
