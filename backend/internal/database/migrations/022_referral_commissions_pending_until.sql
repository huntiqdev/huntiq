-- Migration 022: Add pending_until column to referral_commissions
-- Required by commission_confirm_job.go and referral_repo.go RecordPendingCommission
-- Also add unique constraint on (referrer_id, transaction_id, transaction_type)
-- to support the ON CONFLICT DO NOTHING guard in RecordCommission/RecordPendingCommission.

-- Step 1: Add the missing column
ALTER TABLE referral_commissions
    ADD COLUMN IF NOT EXISTS pending_until TIMESTAMPTZ;

-- Step 2: Add the unique constraint referenced by uq_referral_commission_source
--         (used in RecordCommission & RecordPendingCommission ON CONFLICT clauses)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'uq_referral_commission_source'
    ) THEN
        ALTER TABLE referral_commissions
            ADD CONSTRAINT uq_referral_commission_source
            UNIQUE (referrer_id, transaction_id, transaction_type);
    END IF;
END;
$$;

-- Step 3: Add an index to accelerate the commission-confirm job query
--         (finds pending rows where pending_until has passed)
CREATE INDEX IF NOT EXISTS idx_referral_commissions_pending_until
    ON referral_commissions (status, pending_until)
    WHERE status = 'pending';
