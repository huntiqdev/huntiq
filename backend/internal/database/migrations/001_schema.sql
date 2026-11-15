-- ============================================================================
-- 001_schema.sql
-- HUNTIQ Master Schema Architecture
-- Contains global extensions, system enumerated types, and master table definitions.
-- All 65 domain tables are initialized here strictly alphabetically.
-- Column attributes, check constraints, foreign keys, and indexes are added
-- in the individual table migrations (002_abuse_reports.sql to 066_users.sql).
-- ============================================================================

-- Global Extensions
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gist;
CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "ltree";

-- System Enumerated Types
DO $$ BEGIN CREATE TYPE seller_tier          AS ENUM ('normal','gold','premium');                                                                          EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE buyer_tier           AS ENUM ('normal','gold','premium');                                                                          EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE admin_role           AS ENUM ('normal','super','platform_owner','platform_admin','security_admin','operations_admin','moderator','support_agent'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE product_condition    AS ENUM ('new','like_new','good','fair','poor');                                                              EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE product_status       AS ENUM ('pending','approved','rejected','expired','sold');                                                   EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE transaction_status   AS ENUM ('pending','completed','cancelled');                                                                  EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE order_status         AS ENUM ('pending','accepted','declined','completed','cancelled');                                            EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE admin_action         AS ENUM ('approve','reject','flag','edit','delete','note','activate','deactivate','tier_change','permission_change','create','setting_change','invitation'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE listing_type         AS ENUM ('sale','rent');                                                                                      EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE rental_duration_unit AS ENUM ('hour','day','week','month','year');                                                                EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- ============================================================================
-- Master Table Creation (Ordered Alphabetically, No Attributes)
-- ============================================================================
CREATE TABLE IF NOT EXISTS abuse_reports (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS account_appeals (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS admin_auth_attempts (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS admin_auth_challenges (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS admin_logs (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS admin_notifications (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS admin_passkeys (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS admin_permissions (admin_id INT NOT NULL, permission VARCHAR(50) NOT NULL, PRIMARY KEY (admin_id, permission));
CREATE TABLE IF NOT EXISTS admin_reauth_tokens (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS admin_sessions (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS attribute_definitions (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS buyer_tier_payments (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS category_aliases (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS category_attributes (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS category_hierarchy (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS category_redirects (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS financial_actuals (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS financial_assumptions (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS financial_audit_logs (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS financial_scenarios (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS financial_tax_rules (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS hero_banners (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS invitation_requests (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS job_execution_log (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS location_hierarchy (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS locations (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS notifications (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS order_requests (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS otp_codes (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS platform_settings (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS product_attributes (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS product_images (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS products (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS referral_commissions (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS referral_withdrawals (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS refund_requests (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS schema_migrations (version character varying(255) PRIMARY KEY, applied_at timestamp with time zone NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS search_analytics (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS seller_verification_documents (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS settings_change_log (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS taxonomy_audit_logs (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS transactions (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS upgrade_requests (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS user_activity_logs (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS user_capabilities (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS user_refresh_tokens (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS user_tier_history (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS user_warnings (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY);
-- ============================================================================
