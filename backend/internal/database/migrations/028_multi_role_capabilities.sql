-- ============================================================
-- 028_multi_role_capabilities.sql
-- PURPOSE: Unified User Identity with Additive Marketplace Capabilities
--          (Buyer, Seller, Landlord, Renter).
--          Enables one account per person with multiple simultaneous capabilities.
-- ============================================================

-- 1. Create user_capabilities table
CREATE TABLE IF NOT EXISTS user_capabilities (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    capability VARCHAR(30) NOT NULL, -- 'buyer', 'seller', 'landlord', 'renter'
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- 'active', 'pending', 'suspended', 'revoked'
    verified_at TIMESTAMPTZ,
    activated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'::jsonb,
    CONSTRAINT uq_user_capability UNIQUE (user_id, capability)
);

-- Fast lookup indexes
CREATE INDEX IF NOT EXISTS idx_user_capabilities_user_id ON user_capabilities (user_id);
CREATE INDEX IF NOT EXISTS idx_user_capabilities_active ON user_capabilities (user_id, capability) WHERE status = 'active';

-- 2. Helper functions to check capabilities
CREATE OR REPLACE FUNCTION user_has_capability(p_user_id INT, p_capability VARCHAR)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM user_capabilities
        WHERE user_id = p_user_id
          AND capability = p_capability
          AND status = 'active'
    );
END;
$$ LANGUAGE plpgsql STABLE;

CREATE OR REPLACE FUNCTION get_user_capabilities(p_user_id INT)
RETURNS TEXT[] AS $$
BEGIN
    RETURN ARRAY(
        SELECT capability FROM user_capabilities
        WHERE user_id = p_user_id AND status = 'active'
        ORDER BY capability ASC
    );
END;
$$ LANGUAGE plpgsql STABLE;

-- 3. Data Migration / Backfill for existing users
-- All registered users receive 'buyer' capability by default
INSERT INTO user_capabilities (user_id, capability, status, activated_at)
SELECT id, 'buyer', 'active', created_at
FROM users
ON CONFLICT (user_id, capability) DO NOTHING;

-- Users with role = 'seller' or seller_preference IN ('sale', 'both') receive 'seller' capability
INSERT INTO user_capabilities (user_id, capability, status, activated_at)
SELECT id, 'seller', 'active', created_at
FROM users
WHERE role = 'seller' OR seller_preference IN ('sale', 'both')
ON CONFLICT (user_id, capability) DO NOTHING;

-- Users with seller_preference IN ('rent', 'both') receive 'landlord' capability
INSERT INTO user_capabilities (user_id, capability, status, activated_at)
SELECT id, 'landlord', 'active', created_at
FROM users
WHERE seller_preference IN ('rent', 'both')
ON CONFLICT (user_id, capability) DO NOTHING;

-- View: user_capabilities_summary for easy joining
CREATE OR REPLACE VIEW user_capabilities_view AS
SELECT 
    u.id AS user_id,
    u.name,
    u.phone,
    u.email,
    u.role,
    u.tier AS seller_tier,
    u.buyer_tier,
    u.admin_role,
    u.is_active,
    COALESCE(
        ARRAY_AGG(uc.capability) FILTER (WHERE uc.status = 'active'),
        ARRAY['buyer']::VARCHAR[]
    ) AS active_capabilities
FROM users u
LEFT JOIN user_capabilities uc ON uc.user_id = u.id
GROUP BY u.id;

COMMENT ON TABLE user_capabilities IS 'Additive capabilities per user: buyer, seller, landlord, renter';
