-- ============================================================================
-- 049_users.sql
-- Table: users (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for users
ALTER TABLE users ADD COLUMN IF NOT EXISTS name character varying(100) DEFAULT ''::character varying NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone character varying(20) DEFAULT ''::character varying NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS business_name character varying(150);
ALTER TABLE users ADD COLUMN IF NOT EXISTS email character varying(150);
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash text DEFAULT ''::text NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS rating numeric(3,2) DEFAULT 0.00 NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS role character varying(20) DEFAULT 'seller'::character varying NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS deactivation_reason text;
ALTER TABLE users ADD COLUMN IF NOT EXISTS tier seller_tier DEFAULT 'normal'::seller_tier NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS total_value_sold numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS additional_phones text[];
ALTER TABLE users ADD COLUMN IF NOT EXISTS buyer_tier buyer_tier DEFAULT 'normal'::buyer_tier NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS buyer_tier_paid_at timestamp with time zone;
ALTER TABLE users ADD COLUMN IF NOT EXISTS buyer_tier_expires_at timestamp with time zone;
ALTER TABLE users ADD COLUMN IF NOT EXISTS admin_role admin_role;
ALTER TABLE users ADD COLUMN IF NOT EXISTS otp_verified boolean DEFAULT false NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_banned boolean DEFAULT false NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_at timestamp with time zone;
ALTER TABLE users ADD COLUMN IF NOT EXISTS invited_by integer;
ALTER TABLE users ADD COLUMN IF NOT EXISTS invitation_accepted_at timestamp with time zone;
ALTER TABLE users ADD COLUMN IF NOT EXISTS referral_suspended boolean DEFAULT false NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS referral_permanently_blocked boolean DEFAULT false NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS referral_block_reason text;
ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS referral_code character varying(11) NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS seller_preference character varying(10) DEFAULT 'sale'::character varying;
ALTER TABLE users ADD COLUMN IF NOT EXISTS current_demand character varying(10);

-- Table Constraints for users
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'users' AND constraint_name = 'users_current_demand_check'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT users_current_demand_check CHECK ((((current_demand)::text = ANY ((ARRAY['FIND'::character varying, 'LIST'::character varying])::text[])) OR (current_demand IS NULL)));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'users' AND constraint_name = 'users_seller_preference_check'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT users_seller_preference_check CHECK (((seller_preference)::text = ANY ((ARRAY['sale'::character varying, 'rent'::character varying, 'both'::character varying])::text[])));
    END IF;
END $$;

COMMENT ON COLUMN users.seller_preference IS 'Seller listing preference: sale (sell products), rent (rent properties), both (can do either)';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'users' AND constraint_name = 'users_phone_role_key'
    ) THEN
        ALTER TABLE ONLY users ADD CONSTRAINT users_phone_role_key UNIQUE (phone, role);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'users' AND constraint_name = 'users_pkey'
    ) THEN
        ALTER TABLE ONLY users ADD CONSTRAINT users_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'users' AND constraint_name = 'users_referral_code_key'
    ) THEN
        ALTER TABLE ONLY users ADD CONSTRAINT users_referral_code_key UNIQUE (referral_code);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_users_current_demand ON users USING btree (current_demand) WHERE (current_demand IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_users_email ON users USING btree (email) WHERE (email IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_users_invited_by ON users USING btree (invited_by) WHERE (invited_by IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_users_last_login ON users USING btree (last_login_at DESC) WHERE (last_login_at IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_users_phone ON users USING btree (phone);

CREATE INDEX IF NOT EXISTS idx_users_ref_blocked ON users USING btree (referral_permanently_blocked) WHERE (referral_permanently_blocked = true);

CREATE INDEX IF NOT EXISTS idx_users_ref_suspended ON users USING btree (referral_suspended) WHERE (referral_suspended = true);

CREATE INDEX IF NOT EXISTS idx_users_role ON users USING btree (role);

CREATE INDEX IF NOT EXISTS idx_users_seller_preference ON users USING btree (seller_preference) WHERE ((role)::text = 'seller'::text);

CREATE INDEX IF NOT EXISTS idx_users_tier ON users USING btree (tier);

CREATE INDEX IF NOT EXISTS idx_users_total_value_sold ON users USING btree (total_value_sold);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'users' AND constraint_name = 'users_invited_by_fkey'
    ) THEN
        ALTER TABLE ONLY users ADD CONSTRAINT users_invited_by_fkey FOREIGN KEY (invited_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

ALTER TABLE users ADD COLUMN IF NOT EXISTS custom_url character varying(500);
ALTER TABLE users ADD COLUMN IF NOT EXISTS social_links jsonb NOT NULL DEFAULT '{}'::jsonb;
