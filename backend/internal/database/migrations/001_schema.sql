-- 001_schema.sql
-- Creates all 23 empty tables and enums. No columns. No indexes.
-- Each table's columns are added by its own ALTER file (002–024).

-- ── Extensions ────────────────────────────────────────────────────────────────
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- ── Enums ─────────────────────────────────────────────────────────────────────
DO $$ BEGIN CREATE TYPE seller_tier        AS ENUM ('normal','gold','premium');                                                                          EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE buyer_tier         AS ENUM ('normal','gold','premium');                                                                   EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE admin_role         AS ENUM ('normal','super');                                                                                   EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE product_condition  AS ENUM ('new','like_new','good','fair','poor');                                                              EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE product_status     AS ENUM ('pending','approved','rejected','expired','sold');                                                   EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE transaction_status AS ENUM ('pending','completed','cancelled');                                                                  EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE order_status       AS ENUM ('pending','accepted','declined','completed','cancelled');                                            EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE admin_action       AS ENUM ('approve','reject','flag','edit','delete','note','activate','deactivate','tier_change','permission_change','create','setting_change','invitation'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- ── Tables ────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS users                (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS locations            (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS categories           (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS products             (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS product_images       (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS notifications        (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS transactions         (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS order_requests       (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS admin_logs           (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS otp_codes            (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS platform_settings    (key VARCHAR(100) PRIMARY KEY);
CREATE TABLE IF NOT EXISTS admin_notifications  (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS admin_permissions    (admin_id INT NOT NULL, permission VARCHAR(50) NOT NULL, PRIMARY KEY (admin_id, permission));
CREATE TABLE IF NOT EXISTS hero_banners         (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS user_activity_logs   (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS upgrade_requests     (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS referral_commissions (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS referral_withdrawals (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS seller_reviews       (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS user_tier_history    (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS settings_change_log  (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS invitation_requests  (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS user_warnings        (id SERIAL PRIMARY KEY);
CREATE TABLE IF NOT EXISTS buyer_tier_payments   (id SERIAL PRIMARY KEY);
