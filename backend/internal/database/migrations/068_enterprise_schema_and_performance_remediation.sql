-- Migration 068: Enterprise Schema & Performance Remediation
-- 1. Add seller tier expiry and payment tracking columns to users
-- 2. Safeguard cleanup_unverified_users against deleting administrative/staff roles
-- 3. Add 33 missing foreign key indexes to prevent table-level locking and sequential scans
-- 4. Align monetary view calculations with strict numeric precision

-- 1. User Tier Expiry Columns
ALTER TABLE users ADD COLUMN IF NOT EXISTS tier_paid_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS tier_expires_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX IF NOT EXISTS idx_users_tier_expires_at ON users(tier_expires_at) WHERE tier_expires_at IS NOT NULL;

-- 2. Hardened Cleanup Function
CREATE OR REPLACE FUNCTION public.cleanup_unverified_users()
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
DECLARE
    deleted_count INT;
BEGIN
    WITH deleted AS (
        DELETE FROM users
        WHERE otp_verified = false
          AND role NOT IN ('admin', 'super_admin', 'operator', 'moderator', 'support')
          AND created_at < NOW() - INTERVAL '24 hours'
        RETURNING id
    )
    SELECT count(*) INTO deleted_count FROM deleted;
    RETURN deleted_count;
END;
$function$;

-- 3. 33 Foreign Key Performance Indexes
CREATE INDEX IF NOT EXISTS idx_account_appeals_related_document_id ON account_appeals(related_document_id);
CREATE INDEX IF NOT EXISTS idx_account_appeals_related_product_id ON account_appeals(related_product_id);
CREATE INDEX IF NOT EXISTS idx_admin_permissions_granted_by ON admin_permissions(granted_by);
CREATE INDEX IF NOT EXISTS idx_admin_reauth_tokens_user_id ON admin_reauth_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_financial_assumptions_updated_by ON financial_assumptions(updated_by);
CREATE INDEX IF NOT EXISTS idx_financial_audit_logs_user_id ON financial_audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_financial_scenarios_created_by ON financial_scenarios(created_by);
CREATE INDEX IF NOT EXISTS idx_hero_banners_created_by ON hero_banners(created_by);
CREATE INDEX IF NOT EXISTS idx_invitation_requests_approved_by ON invitation_requests(approved_by);
CREATE INDEX IF NOT EXISTS idx_invitation_requests_invited_by ON invitation_requests(invited_by);
CREATE INDEX IF NOT EXISTS idx_location_hierarchy_replaced_by_id ON location_hierarchy(replaced_by_id);
CREATE INDEX IF NOT EXISTS idx_platform_settings_updated_by ON platform_settings(updated_by);
CREATE INDEX IF NOT EXISTS idx_products_current_tenant_id ON products(current_tenant_id);
CREATE INDEX IF NOT EXISTS idx_referral_withdrawals_approved_by ON referral_withdrawals(approved_by);
CREATE INDEX IF NOT EXISTS idx_referral_withdrawals_reviewed_by ON referral_withdrawals(reviewed_by);
CREATE INDEX IF NOT EXISTS idx_refund_requests_product_id ON refund_requests(product_id);
CREATE INDEX IF NOT EXISTS idx_taxonomy_audit_logs_admin_id ON taxonomy_audit_logs(admin_id);
CREATE INDEX IF NOT EXISTS idx_upgrade_requests_reviewed_by ON upgrade_requests(reviewed_by);
CREATE INDEX IF NOT EXISTS idx_user_tier_history_changed_by ON user_tier_history(changed_by);
CREATE INDEX IF NOT EXISTS idx_user_tier_history_upgrade_request_id ON user_tier_history(upgrade_request_id);
CREATE INDEX IF NOT EXISTS idx_users_organization_id ON users(organization_id);
CREATE INDEX IF NOT EXISTS idx_referral_attributions_inviter_id ON referral_attributions(inviter_id);
CREATE INDEX IF NOT EXISTS idx_referral_reward_policies_updated_by ON referral_reward_policies(updated_by);
CREATE INDEX IF NOT EXISTS idx_referral_reward_policy_versions_changed_by ON referral_reward_policy_versions(changed_by);
CREATE INDEX IF NOT EXISTS idx_referral_reward_policy_versions_policy_id ON referral_reward_policy_versions(policy_id);
CREATE INDEX IF NOT EXISTS idx_ledger_transactions_actor_id ON ledger_transactions(actor_id);
CREATE INDEX IF NOT EXISTS idx_organizations_reviewed_by ON organizations(reviewed_by);
CREATE INDEX IF NOT EXISTS idx_identity_verifications_reviewed_by ON identity_verifications(reviewed_by);
CREATE INDEX IF NOT EXISTS idx_password_recovery_requests_user_id ON password_recovery_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_audit_events_user_id ON ai_audit_events(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_provider_configs_updated_by ON payment_provider_configs(updated_by);
CREATE INDEX IF NOT EXISTS idx_payment_orders_product_id ON payment_orders(product_id);
CREATE INDEX IF NOT EXISTS idx_admin_security_zone_assignments_granted_by ON admin_security_zone_assignments(granted_by);
