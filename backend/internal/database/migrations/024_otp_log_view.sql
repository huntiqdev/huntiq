-- ============================================================
-- 024_otp_log_view.sql
-- PURPOSE: Create otp_log_view for admin OTP audit logs
-- ============================================================

CREATE OR REPLACE VIEW otp_log_view AS
SELECT
    id,
    COALESCE(phone, '') AS phone,
    COALESCE(email, '') AS email,
    purpose,
    user_type,
    used,
    attempts,
    expires_at,
    created_at,
    status,
    COALESCE(code_encrypted, code, '') AS encrypted_code
FROM otp_codes;
