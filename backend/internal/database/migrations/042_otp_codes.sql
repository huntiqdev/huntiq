-- ============================================================================
-- 042_otp_codes.sql
-- Table: otp_codes (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for otp_codes
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS phone character varying(20);
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS email character varying(150);
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS code_encrypted text DEFAULT ''::text NOT NULL;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS purpose character varying(50) DEFAULT 'register'::character varying NOT NULL;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS user_type character varying(20) DEFAULT 'seller'::character varying NOT NULL;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS attempts integer DEFAULT 0 NOT NULL;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS max_attempts integer DEFAULT 5 NOT NULL;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'pending'::character varying NOT NULL;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS expires_at timestamp with time zone DEFAULT (now() + '00:05:00'::interval) NOT NULL;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS verified_at timestamp with time zone;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS used boolean DEFAULT false NOT NULL;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS locked_until timestamp with time zone;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS code text DEFAULT ''::text NOT NULL;

-- Table Constraints for otp_codes
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'otp_codes' AND constraint_name = 'otp_contact_check'
    ) THEN
        ALTER TABLE otp_codes ADD CONSTRAINT otp_contact_check CHECK (((phone IS NOT NULL) OR (email IS NOT NULL)));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'otp_codes' AND constraint_name = 'otp_codes_pkey'
    ) THEN
        ALTER TABLE ONLY otp_codes ADD CONSTRAINT otp_codes_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_otp_email ON otp_codes USING btree (email, status, expires_at DESC);

CREATE INDEX IF NOT EXISTS idx_otp_email_used_expires ON otp_codes USING btree (email, purpose, used, expires_at);

CREATE INDEX IF NOT EXISTS idx_otp_phone ON otp_codes USING btree (phone, status, expires_at DESC);

CREATE INDEX IF NOT EXISTS idx_otp_status ON otp_codes USING btree (status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_otp_used_expires ON otp_codes USING btree (phone, purpose, used, expires_at);
