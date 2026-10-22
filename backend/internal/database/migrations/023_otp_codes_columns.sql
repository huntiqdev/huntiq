-- ============================================================
-- 023_otp_codes_columns.sql
-- PURPOSE: Add used, locked_until, and code columns to otp_codes table
--          to align database schema with UserRepo queries and maintain
--          bidirectional synchronization between code/code_encrypted and status/used.
-- ============================================================

ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS used BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS locked_until TIMESTAMPTZ;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS code TEXT NOT NULL DEFAULT '';

-- Populate any existing rows
UPDATE otp_codes SET code = code_encrypted WHERE code = '' AND code_encrypted != '';
UPDATE otp_codes SET code_encrypted = code WHERE code_encrypted = '' AND code != '';
UPDATE otp_codes SET used = TRUE WHERE status IN ('verified', 'used');

-- Trigger to automatically synchronize code <-> code_encrypted and status <-> used
CREATE OR REPLACE FUNCTION sync_otp_codes_columns()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.code != '' AND (NEW.code_encrypted IS NULL OR NEW.code_encrypted = '') THEN
        NEW.code_encrypted := NEW.code;
    ELSIF NEW.code_encrypted != '' AND (NEW.code IS NULL OR NEW.code = '') THEN
        NEW.code := NEW.code_encrypted;
    END IF;

    IF NEW.used = TRUE AND NEW.status = 'pending' THEN
        NEW.status := 'verified';
    ELSIF NEW.status = 'verified' AND NEW.used = FALSE THEN
        NEW.used := TRUE;
    END IF;

    IF NEW.locked_until IS NOT NULL AND NEW.locked_until > NOW() THEN
        NEW.status := 'locked';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_sync_otp_codes ON otp_codes;
CREATE TRIGGER trg_sync_otp_codes
BEFORE INSERT OR UPDATE ON otp_codes
FOR EACH ROW EXECUTE FUNCTION sync_otp_codes_columns();

CREATE INDEX IF NOT EXISTS idx_otp_used_expires ON otp_codes(phone, purpose, used, expires_at);
CREATE INDEX IF NOT EXISTS idx_otp_email_used_expires ON otp_codes(email, purpose, used, expires_at);
