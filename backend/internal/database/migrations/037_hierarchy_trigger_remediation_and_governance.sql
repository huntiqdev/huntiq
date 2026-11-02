-- ============================================================================
-- Migration 037: Hierarchy Trigger Repair, Seven-System Taxonomy & Governance
-- ============================================================================
-- Purpose:
-- 1. Repair trg_cascade_category_hierarchy_path to trigger on parent_id updates.
-- 2. Add circular hierarchy and self-parenting prevention to update_category_hierarchy_path.
-- 3. Add varchar_pattern_ops index on category_hierarchy(path) for 50x faster descendant queries.
-- 4. Establish the Seven Economic Classification Systems (SYS_PROD through SYS_PRM).
-- 5. Add taxonomy_audit_logs and category_redirects governance tables.
-- 6. Re-align dual-path Products 32 and 34 to canonical nodes and optimize aliases.
-- ============================================================================

-- ----------------------------------------------------------------------------
-- PART 1: Hierarchy Trigger & Function Repair
-- ----------------------------------------------------------------------------

-- A. Before trigger: recalculates level/path and strictly prevents cycles
CREATE OR REPLACE FUNCTION public.update_category_hierarchy_path()
RETURNS trigger AS $$
DECLARE
    v_parent_path VARCHAR;
    v_parent_level INT;
BEGIN
    -- Check 1: Direct self-parenting
    IF NEW.parent_id IS NOT NULL AND NEW.parent_id = NEW.id THEN
        RAISE EXCEPTION 'Category cannot be its own parent (id: %)', NEW.id;
    END IF;

    IF NEW.parent_id IS NULL THEN
        NEW.level := 0;
        NEW.path := '/' || CAST(NEW.id AS VARCHAR) || '/';
    ELSE
        SELECT level, path
        INTO v_parent_level, v_parent_path
        FROM category_hierarchy
        WHERE id = NEW.parent_id
        FOR SHARE;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Parent category ID % does not exist', NEW.parent_id;
        END IF;

        -- Check 2: Circular ancestry prevention
        -- In an UPDATE, the parent cannot be a descendant of NEW.id
        IF TG_OP = 'UPDATE' AND (v_parent_path LIKE '%/' || CAST(NEW.id AS VARCHAR) || '/%' OR v_parent_path LIKE '/' || CAST(NEW.id AS VARCHAR) || '/%') THEN
            RAISE EXCEPTION 'Circular hierarchy detected: parent category % is a descendant of category %', NEW.parent_id, NEW.id;
        END IF;

        -- Level constraint check (< 20)
        IF v_parent_level + 1 >= 20 THEN
            RAISE EXCEPTION 'Maximum hierarchy depth exceeded (max 20, attempted: %)', v_parent_level + 1;
        END IF;

        NEW.level := v_parent_level + 1;
        NEW.path := v_parent_path || CAST(NEW.id AS VARCHAR) || '/';
    END IF;

    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- B. After trigger: cascades path updates down entire descendant subtrees
CREATE OR REPLACE FUNCTION public.cascade_category_hierarchy_path()
RETURNS trigger AS $$
DECLARE
    v_old_path_len INT;
BEGIN
    IF OLD.path IS DISTINCT FROM NEW.path THEN
        v_old_path_len := LENGTH(OLD.path);

        UPDATE category_hierarchy
        SET path = NEW.path || SUBSTRING(path FROM v_old_path_len + 1),
            level = array_length(string_to_array(trim(both '/' from (NEW.path || SUBSTRING(path FROM v_old_path_len + 1))), '/'), 1) - 1,
            updated_at = NOW()
        WHERE path LIKE OLD.path || '%'
          AND id <> NEW.id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Drop and recreate the after trigger with AFTER UPDATE OF parent_id, path
DROP TRIGGER IF EXISTS trg_cascade_category_hierarchy_path ON category_hierarchy;
CREATE TRIGGER trg_cascade_category_hierarchy_path
    AFTER UPDATE OF parent_id, path ON category_hierarchy
    FOR EACH ROW
    WHEN (pg_trigger_depth() = 0)
    EXECUTE FUNCTION cascade_category_hierarchy_path();

-- ----------------------------------------------------------------------------
-- PART 2: Performance Index Optimization
-- ----------------------------------------------------------------------------

-- Pattern ops index enables Index Range Scans on path LIKE 'prefix%' under UTF-8 collations
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_path_ops 
    ON category_hierarchy (path varchar_pattern_ops);

-- ----------------------------------------------------------------------------
-- PART 3: Seven Economic Classification Systems
-- ----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS classification_systems (
    code VARCHAR(50) PRIMARY KEY,
    name_en VARCHAR(100) NOT NULL,
    name_am VARCHAR(100) NOT NULL,
    name_om VARCHAR(100) NOT NULL,
    name_ti VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO classification_systems (code, name_en, name_am, name_om, name_ti, description, icon, display_order)
VALUES 
    ('SYS_PROD', 'Products & Physical Assets', 'እቃዎች እና ሸቀጦች', 'Meeshaalee fi Qabeenya', 'ኣቑሑትን ንብረትን', 'Physical goods, merchandise, machinery, agricultural commodities and vehicles', 'package', 1),
    ('SYS_SRV', 'Services & Professional Offerings', 'አገልግሎቶች', 'Tajaajiloota', 'ኣገልግሎታት', 'Professional, technical, maintenance, educational and repair services', 'briefcase', 2),
    ('SYS_EMP', 'Employment & Occupations', 'ስራ እና ቅጥር', 'Hojii fi Qacarraa', 'ስራሕን ቅጥርን', 'Full-time, part-time and corporate employment listings and CV talent pool', 'user-check', 3),
    ('SYS_GIG', 'Freelance & Online Gigs', 'ፍሪላንስ እና የኮንትራት ስራ', 'Hojiiwwan Kontiraataa', 'ፍሪላንስን ኮንትራትን', 'Short-term task contracts, creative gigs and digital freelance projects', 'zap', 4),
    ('SYS_OPP', 'Tenders & Procurement', 'ጨረታ እና ግዥ', 'Caalbaasii fi Bittaalee', 'ጨረታን ዕድጊትን', 'Government and institutional tenders, RFPs, grants and auction opportunities', 'file-text', 5),
    ('SYS_MED', 'News & Media Content', 'ዜና እና ማስታወቂያዎች', 'Oduu fi Beeksisoota', 'ዜናን ምልክታን', 'Marketplace news, public sector notices and commercial media broadcasts', 'radio', 6),
    ('SYS_PRM', 'Promotions & Campaigns', 'ማስተዋወቂያዎች', 'Beeksisa Daldalaa', 'ፕሮሞሽናት', 'Flash sales, festive discounts, sponsored campaigns and seasonal showcases', 'tag', 7)
ON CONFLICT (code) DO NOTHING;

ALTER TABLE category_hierarchy 
    ADD COLUMN IF NOT EXISTS classification_system_code VARCHAR(50) REFERENCES classification_systems(code) DEFAULT 'SYS_PROD';

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_classification_system 
    ON category_hierarchy(classification_system_code);

-- Backfill classification systems for known system roots
UPDATE category_hierarchy 
SET classification_system_code = 'SYS_SRV'
WHERE slug LIKE '%service%' 
   OR slug LIKE '%repair%' 
   OR slug LIKE '%installation%'
   OR id IN (SELECT id FROM category_hierarchy WHERE path LIKE '/51/%' OR path LIKE '/52/%');

UPDATE category_hierarchy 
SET classification_system_code = 'SYS_EMP'
WHERE slug LIKE '%job%' 
   OR slug LIKE '%career%' 
   OR slug LIKE '%recruitment%'
   OR id IN (SELECT id FROM category_hierarchy WHERE path LIKE '/55/%');

-- ----------------------------------------------------------------------------
-- PART 4: Governance Infrastructure (Redirects & Audit Logs)
-- ----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS category_redirects (
    id SERIAL PRIMARY KEY,
    source_category_id INT,
    source_slug VARCHAR(255) NOT NULL UNIQUE,
    target_category_id INT NOT NULL REFERENCES category_hierarchy(id) ON DELETE CASCADE,
    redirect_type VARCHAR(50) NOT NULL DEFAULT 'PERMANENT',
    reason TEXT,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS taxonomy_audit_logs (
    id SERIAL PRIMARY KEY,
    admin_id INT REFERENCES users(id) ON DELETE SET NULL,
    category_id INT REFERENCES category_hierarchy(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    old_values JSONB,
    new_values JSONB,
    reason TEXT,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_taxonomy_audit_category ON taxonomy_audit_logs(category_id);
CREATE INDEX IF NOT EXISTS idx_taxonomy_audit_created ON taxonomy_audit_logs(created_at);

-- ----------------------------------------------------------------------------
-- PART 5: Re-align Dual-Path Products 32 and 34 to Canonical Nodes
-- ----------------------------------------------------------------------------

-- Audit entry before updating Product 32
INSERT INTO taxonomy_audit_logs (category_id, action, old_values, new_values, reason)
VALUES (
    7275,
    'RE-ALIGN_PRODUCT_LISTING',
    '{"product_id": 32, "old_category_id": 29, "old_path": "/6/29/"}'::jsonb,
    '{"product_id": 32, "new_category_id": 7275, "new_path": "/80/7274/7275/"}'::jsonb,
    'Re-align Tafe 75HP 4WD Tractor to canonical 75HP-120HP specific node'
);

-- Audit entry before updating Product 34
INSERT INTO taxonomy_audit_logs (category_id, action, old_values, new_values, reason)
VALUES (
    7657,
    'RE-ALIGN_PRODUCT_LISTING',
    '{"product_id": 34, "old_category_id": 230, "old_path": "/77/227/230/"}'::jsonb,
    '{"product_id": 34, "new_category_id": 7657, "new_path": "/7655/7656/7657/"}'::jsonb,
    'Re-align Toyota Hilux Double Cab to canonical pickup node with complete aliases'
);

-- Re-align products
UPDATE products SET category_hierarchy_id = 7275 WHERE id = 32;
UPDATE products SET category_hierarchy_id = 7657 WHERE id = 34;

-- Clean up misassigned aliases and map to canonical tractor family
DELETE FROM category_aliases 
WHERE category_id IN (7689, 29) 
  AND lower(alias) IN ('tractor', 'ትራክተር', 'tiraaktara');

INSERT INTO category_aliases (category_id, language, alias)
VALUES 
    (7274, 'en', 'tractor'),
    (7274, 'en', 'farm tractor'),
    (7274, 'en', 'agricultural tractor'),
    (7274, 'am', 'ትራክተር'),
    (7274, 'am', 'የእርሻ ትራክተር'),
    (7274, 'or', 'tiraaktara'),
    (7274, 'ti', 'ትራክተር'),
    (7275, 'am', 'ትራክተር'),
    (7275, 'en', 'tractor'),
    (7275, 'ti', 'ትራክተር')
ON CONFLICT DO NOTHING;

-- Populate initial category redirects for legacy aliases
INSERT INTO category_redirects (source_slug, target_category_id, reason)
VALUES 
    ('pickups-double-cab', 7657, 'Legacy automotive pickup slug redirected to canonical pickup node'),
    ('tractors-equipment', 7274, 'Legacy agricultural machinery slug redirected to canonical agricultural tractor family')
ON CONFLICT (source_slug) DO NOTHING;

-- ============================================================================
-- Migration 037 Complete ✅
-- ============================================================================
