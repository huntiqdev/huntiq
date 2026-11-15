-- ============================================================================
-- 078_universal_classification_platform.sql
-- Master Migration: Universal Database-Driven Classification, Catalog,
-- Relationship, Attribute, Search and Administration Platform.
-- Single authoritative PostgreSQL engine for multi-scheme taxonomy.
-- ============================================================================

-- ============================================================================
-- PART 1: CLASSIFICATION SCHEMES
-- ============================================================================

CREATE TABLE IF NOT EXISTS classification_schemes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name_en VARCHAR(200) NOT NULL,
    name_am VARCHAR(200),
    name_om VARCHAR(200),
    name_ti VARCHAR(200),
    description TEXT,
    icon VARCHAR(100),
    status VARCHAR(20) DEFAULT 'ACTIVE' NOT NULL,
    version INT DEFAULT 1 NOT NULL,
    revision INT DEFAULT 1 NOT NULL,
    is_system BOOLEAN DEFAULT false NOT NULL,
    rules JSONB DEFAULT '{}'::jsonb NOT NULL,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    updated_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_schemes' AND constraint_name = 'chk_scheme_status'
    ) THEN
        ALTER TABLE classification_schemes ADD CONSTRAINT chk_scheme_status
            CHECK (status IN ('ACTIVE', 'SUSPENDED', 'ARCHIVED', 'DRAFT'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_schemes_status ON classification_schemes USING btree (status);
CREATE INDEX IF NOT EXISTS idx_schemes_code ON classification_schemes USING btree (code);

-- ============================================================================
-- PART 2: UNIVERSAL CANONICAL CONCEPTS
-- ============================================================================

CREATE TABLE IF NOT EXISTS classification_concepts (
    id SERIAL PRIMARY KEY,
    code VARCHAR(100) UNIQUE NOT NULL,
    canonical_name VARCHAR(200) NOT NULL,
    name_en VARCHAR(200) NOT NULL,
    name_am VARCHAR(200),
    name_om VARCHAR(200),
    name_ti VARCHAR(200),
    description TEXT,
    concept_type VARCHAR(50) DEFAULT 'PRODUCT' NOT NULL,
    status VARCHAR(20) DEFAULT 'ACTIVE' NOT NULL,
    brand_id INT REFERENCES brands(id) ON DELETE SET NULL,
    metadata JSONB DEFAULT '{}'::jsonb NOT NULL,
    version INT DEFAULT 1 NOT NULL,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    updated_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_concepts' AND constraint_name = 'chk_concept_type'
    ) THEN
        ALTER TABLE classification_concepts ADD CONSTRAINT chk_concept_type
            CHECK (concept_type IN ('PRODUCT', 'SERVICE', 'MATERIAL', 'EQUIPMENT', 'PROPERTY', 'EMPLOYMENT', 'ORGANIZATION', 'OTHER'));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_concepts' AND constraint_name = 'chk_concept_status'
    ) THEN
        ALTER TABLE classification_concepts ADD CONSTRAINT chk_concept_status
            CHECK (status IN ('ACTIVE', 'DRAFT', 'REVIEW', 'ARCHIVED', 'RETIRED'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_concepts_status ON classification_concepts USING btree (status);
CREATE INDEX IF NOT EXISTS idx_concepts_type ON classification_concepts USING btree (concept_type);
CREATE INDEX IF NOT EXISTS idx_concepts_code ON classification_concepts USING btree (code);
CREATE INDEX IF NOT EXISTS idx_concepts_name_en_trgm ON classification_concepts USING gin (name_en gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_concepts_name_am_trgm ON classification_concepts USING gin (name_am gin_trgm_ops) WHERE name_am IS NOT NULL;

-- ============================================================================
-- PART 3: CLASSIFICATION PLACEMENTS (Hierarchy & Multi-scheme placements)
-- ============================================================================

CREATE TABLE IF NOT EXISTS classification_placements (
    id SERIAL PRIMARY KEY,
    scheme_id INT NOT NULL REFERENCES classification_schemes(id) ON DELETE RESTRICT,
    concept_id INT NOT NULL REFERENCES classification_concepts(id) ON DELETE CASCADE,
    parent_placement_id INT REFERENCES classification_placements(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'CATEGORY' NOT NULL,
    display_name_override VARCHAR(200),
    slug VARCHAR(255) NOT NULL,
    path VARCHAR(1000) DEFAULT '/' NOT NULL,
    depth INT DEFAULT 0 NOT NULL,
    display_order INT DEFAULT 0 NOT NULL,
    is_leaf BOOLEAN DEFAULT false NOT NULL,
    is_active BOOLEAN DEFAULT true NOT NULL,
    status VARCHAR(20) DEFAULT 'PUBLISHED' NOT NULL,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    updated_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_placements' AND constraint_name = 'chk_placement_role'
    ) THEN
        ALTER TABLE classification_placements ADD CONSTRAINT chk_placement_role
            CHECK (role IN ('ROOT', 'SUPERFAMILY', 'CATEGORY', 'SUBCATEGORY', 'LEAF', 'OTHER'));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_placements' AND constraint_name = 'chk_placement_status'
    ) THEN
        ALTER TABLE classification_placements ADD CONSTRAINT chk_placement_status
            CHECK (status IN ('DRAFT', 'REVIEW', 'APPROVED', 'PUBLISHED', 'ARCHIVED'));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_placements' AND constraint_name = 'no_placement_self_parent'
    ) THEN
        ALTER TABLE classification_placements ADD CONSTRAINT no_placement_self_parent
            CHECK (id <> parent_placement_id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_placements' AND constraint_name = 'uq_placement_scheme_slug'
    ) THEN
        ALTER TABLE classification_placements ADD CONSTRAINT uq_placement_scheme_slug
            UNIQUE (scheme_id, slug);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_placements_scheme ON classification_placements USING btree (scheme_id);
CREATE INDEX IF NOT EXISTS idx_placements_concept ON classification_placements USING btree (concept_id);
CREATE INDEX IF NOT EXISTS idx_placements_parent ON classification_placements USING btree (parent_placement_id);
CREATE INDEX IF NOT EXISTS idx_placements_path ON classification_placements USING btree (path varchar_pattern_ops);
CREATE INDEX IF NOT EXISTS idx_placements_depth ON classification_placements USING btree (depth);
CREATE INDEX IF NOT EXISTS idx_placements_role ON classification_placements USING btree (role);
CREATE INDEX IF NOT EXISTS idx_placements_active ON classification_placements USING btree (is_active);

-- ============================================================================
-- PART 4: ALIASES, TRANSLITERATIONS AND SYNONYMS
-- ============================================================================

CREATE TABLE IF NOT EXISTS classification_aliases (
    id SERIAL PRIMARY KEY,
    concept_id INT REFERENCES classification_concepts(id) ON DELETE CASCADE,
    placement_id INT REFERENCES classification_placements(id) ON DELETE CASCADE,
    alias VARCHAR(255) NOT NULL,
    lang VARCHAR(10) DEFAULT 'en' NOT NULL,
    alias_type VARCHAR(50) DEFAULT 'SYNONYM' NOT NULL,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_aliases' AND constraint_name = 'chk_alias_type'
    ) THEN
        ALTER TABLE classification_aliases ADD CONSTRAINT chk_alias_type
            CHECK (alias_type IN ('SYNONYM', 'ALTERNATIVE_NAME', 'TRANSLITERATION', 'MISSPELLING', 'ACRONYM', 'TRADE_NAME'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_aliases_concept ON classification_aliases USING btree (concept_id);
CREATE INDEX IF NOT EXISTS idx_aliases_placement ON classification_aliases USING btree (placement_id);
CREATE INDEX IF NOT EXISTS idx_aliases_lang ON classification_aliases USING btree (lang);
CREATE INDEX IF NOT EXISTS idx_aliases_alias_trgm ON classification_aliases USING gin (alias gin_trgm_ops);

-- ============================================================================
-- PART 5: MULTI-INDUSTRY & MULTI-APPLICATION RELATIONSHIPS
-- ============================================================================

CREATE TABLE IF NOT EXISTS classification_relationship_types (
    code VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_directed BOOLEAN DEFAULT true NOT NULL,
    inverse_type_code VARCHAR(50) REFERENCES classification_relationship_types(code),
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS classification_relationships (
    id SERIAL PRIMARY KEY,
    source_concept_id INT NOT NULL REFERENCES classification_concepts(id) ON DELETE CASCADE,
    target_concept_id INT NOT NULL REFERENCES classification_concepts(id) ON DELETE CASCADE,
    relationship_type_code VARCHAR(50) NOT NULL REFERENCES classification_relationship_types(code) ON DELETE RESTRICT,
    context_scheme_id INT REFERENCES classification_schemes(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'ACTIVE' NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb NOT NULL,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_relationships' AND constraint_name = 'no_self_relationship'
    ) THEN
        ALTER TABLE classification_relationships ADD CONSTRAINT no_self_relationship
            CHECK (source_concept_id <> target_concept_id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_relationships' AND constraint_name = 'uq_concept_rel'
    ) THEN
        ALTER TABLE classification_relationships ADD CONSTRAINT uq_concept_rel
            UNIQUE (source_concept_id, target_concept_id, relationship_type_code, context_scheme_id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_rels_source ON classification_relationships USING btree (source_concept_id);
CREATE INDEX IF NOT EXISTS idx_rels_target ON classification_relationships USING btree (target_concept_id);
CREATE INDEX IF NOT EXISTS idx_rels_type ON classification_relationships USING btree (relationship_type_code);
CREATE INDEX IF NOT EXISTS idx_rels_context ON classification_relationships USING btree (context_scheme_id);

-- ============================================================================
-- PART 6: DYNAMIC ATTRIBUTES INHERITANCE AND SPECIFICATION RULES
-- ============================================================================

CREATE TABLE IF NOT EXISTS classification_attribute_rules (
    id SERIAL PRIMARY KEY,
    placement_id INT REFERENCES classification_placements(id) ON DELETE CASCADE,
    concept_id INT REFERENCES classification_concepts(id) ON DELETE CASCADE,
    attribute_id INT NOT NULL REFERENCES attribute_definitions(id) ON DELETE CASCADE,
    is_required BOOLEAN DEFAULT false NOT NULL,
    is_inherited BOOLEAN DEFAULT true NOT NULL,
    display_order INT DEFAULT 0 NOT NULL,
    override_validation_rules JSONB,
    override_unit VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_attribute_rules' AND constraint_name = 'chk_attr_rule_target'
    ) THEN
        ALTER TABLE classification_attribute_rules ADD CONSTRAINT chk_attr_rule_target
            CHECK (placement_id IS NOT NULL OR concept_id IS NOT NULL);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_attr_rules_placement ON classification_attribute_rules USING btree (placement_id);
CREATE INDEX IF NOT EXISTS idx_attr_rules_concept ON classification_attribute_rules USING btree (concept_id);
CREATE INDEX IF NOT EXISTS idx_attr_rules_attribute ON classification_attribute_rules USING btree (attribute_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_attr_rule_placement ON classification_attribute_rules (placement_id, attribute_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_attr_rule_concept ON classification_attribute_rules (concept_id, attribute_id);

-- ============================================================================
-- PART 7: GOVERNANCE, CHANGE REQUESTS, REVISIONS & AUDIT LOGS
-- ============================================================================

CREATE TABLE IF NOT EXISTS classification_change_requests (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INT,
    action VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'DRAFT' NOT NULL,
    proposed_changes JSONB NOT NULL,
    previous_state JSONB,
    impact_summary JSONB DEFAULT '{}'::jsonb,
    requested_by INT REFERENCES users(id) ON DELETE SET NULL,
    reviewed_by INT REFERENCES users(id) ON DELETE SET NULL,
    rejection_reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_change_requests' AND constraint_name = 'chk_change_req_status'
    ) THEN
        ALTER TABLE classification_change_requests ADD CONSTRAINT chk_change_req_status
            CHECK (status IN ('DRAFT', 'PENDING_REVIEW', 'APPROVED', 'REJECTED', 'PUBLISHED'));
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS classification_audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INT NOT NULL,
    changes JSONB NOT NULL,
    reason TEXT,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_class_audit_entity ON classification_audit_logs USING btree (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_class_audit_user ON classification_audit_logs USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_class_audit_created ON classification_audit_logs USING btree (created_at);

-- ============================================================================
-- PART 8: INTEGRATION WITH PRODUCTS (MARKETPLACE LISTINGS)
-- ============================================================================

ALTER TABLE products ADD COLUMN IF NOT EXISTS concept_id integer REFERENCES classification_concepts(id) ON DELETE RESTRICT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS placement_id integer REFERENCES classification_placements(id) ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_products_concept_id ON products USING btree (concept_id);
CREATE INDEX IF NOT EXISTS idx_products_placement_id ON products USING btree (placement_id);

-- ============================================================================
-- PART 9: HIERARCHY CASCAFE TRIGGERS FOR PLACEMENTS
-- ============================================================================

CREATE OR REPLACE FUNCTION public.update_classification_placement_path()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
DECLARE
    v_parent_path VARCHAR;
    v_parent_depth INT;
    v_parent_scheme_id INT;
BEGIN
    IF NEW.parent_placement_id IS NULL THEN
        NEW.depth := 0;
        NEW.path := '/' || CAST(NEW.id AS VARCHAR) || '/';
    ELSE
        SELECT depth, path, scheme_id
        INTO v_parent_depth, v_parent_path, v_parent_scheme_id
        FROM classification_placements
        WHERE id = NEW.parent_placement_id;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Parent placement ID % does not exist', NEW.parent_placement_id;
        END IF;

        IF v_parent_scheme_id <> NEW.scheme_id THEN
            RAISE EXCEPTION 'Parent placement ID % belongs to scheme %, but child placement is in scheme %',
                NEW.parent_placement_id, v_parent_scheme_id, NEW.scheme_id;
        END IF;

        NEW.depth := v_parent_depth + 1;
        NEW.path := v_parent_path || CAST(NEW.id AS VARCHAR) || '/';
    END IF;

    NEW.updated_at := NOW();
    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION public.cascade_classification_placement_path()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    IF OLD.path IS DISTINCT FROM NEW.path THEN
        UPDATE classification_placements
        SET path = NEW.path || SUBSTRING(path FROM LENGTH(OLD.path) + 1),
            depth = array_length(string_to_array(trim(both '/' from NEW.path || SUBSTRING(path FROM LENGTH(OLD.path) + 1)), '/'), 1) - 1,
            updated_at = NOW()
        WHERE path LIKE OLD.path || '%';
    END IF;
    RETURN NULL;
END;
$function$;

DROP TRIGGER IF EXISTS trg_update_classification_placement_path ON classification_placements;
CREATE TRIGGER trg_update_classification_placement_path
BEFORE INSERT OR UPDATE OF parent_placement_id, id ON classification_placements
FOR EACH ROW EXECUTE FUNCTION update_classification_placement_path();

DROP TRIGGER IF EXISTS trg_cascade_classification_placement_path ON classification_placements;
CREATE TRIGGER trg_cascade_classification_placement_path
AFTER UPDATE OF path ON classification_placements
FOR EACH ROW EXECUTE FUNCTION cascade_classification_placement_path();

-- ============================================================================
-- PART 10: SEED STANDARD SCHEMES, RELATIONSHIP TYPES, AND MIGRATE INITIAL TAXONOMY
-- ============================================================================

-- 10.1 Standard Schemes
INSERT INTO classification_schemes (code, name_en, name_am, name_om, name_ti, description, is_system, status)
VALUES
('SYS_PROD', 'Universal Product Classification', 'ሁለንተናዊ የምርት ምደባ', 'Ramaddii Oomishaa Waliigalaa', 'ኩለመዳላዊ ምደባ ፍርያት', 'Master product marketplace classification scheme', true, 'ACTIVE'),
('SYS_IND', 'Industry Classification', 'የኢንዱስትሪ ምደባ', 'Ramaddii Industirii', 'ምደባ ኢንዱስትሪ', 'Economic sector and business vertical taxonomy', true, 'ACTIVE'),
('SYS_SERV', 'Services & Trades Classification', 'የአገልግሎትና የሙያ ምደባ', 'Ramaddii Tajaajilootaa', 'ምደባ ኣገልግሎታትን ሞያታትን', 'Service providers, contractors and professional trades', true, 'ACTIVE'),
('SYS_PROP', 'Real Estate & Property Classification', 'የሪል እስቴትና የንብረት ምደባ', 'Ramaddii Qabeenya fi Manneen', 'ምደባ ሪል እስቴትን ንብረትን', 'Residential, commercial, and land property classification', true, 'ACTIVE'),
('SYS_VEH', 'Automotive & Mobility Classification', 'የተሽከርካሪና የትራንስፖርት ምደባ', 'Ramaddii Konkolaattotaa', 'ምደባ ተሽከርከርትን ትራንስፖርትን', 'Vehicles, heavy machinery, and spare parts taxonomy', true, 'ACTIVE'),
('SYS_AGRI', 'Agriculture & Farming Classification', 'የግብርናና የእርሻ ምደባ', 'Ramaddii Qonnaa', 'ምደባ ሕርሻን ልምዓትን', 'Crops, livestock, farm machinery and agricultural commodities', true, 'ACTIVE'),
('SYS_JOBS', 'Workforce & Employment Classification', 'የስራና የቅጥር ምደባ', 'Ramaddii Hojii', 'ምደባ ስራሕን ስራሕተኛታትን', 'Occupations, qualifications, and freelance labor taxonomy', true, 'ACTIVE'),
('SYS_TOOLS', 'Tools & Machinery Classification', 'የመሳሪያዎችና ማሽነሪዎች ምደባ', 'Ramaddii Meeshaalee fi Maashinootaa', 'ምደባ መሳርሒታትን ማሽነሪታትን', 'Universal tools, power equipment, and machinery taxonomy', false, 'ACTIVE'),
('SYS_ELEC', 'Electronics Classification', 'የኤሌክትሮኒክስ ምደባ', 'Ramaddii Eleektirooniksii', 'ምደባ ኤሌክትሮኒክስ', 'Universal consumer electronics, devices, and computing hardware', false, 'ACTIVE')
ON CONFLICT (code) DO UPDATE SET
    name_en = EXCLUDED.name_en,
    name_am = EXCLUDED.name_am,
    name_om = EXCLUDED.name_om,
    name_ti = EXCLUDED.name_ti,
    status = 'ACTIVE';

-- 10.2 Standard Relationship Types
INSERT INTO classification_relationship_types (code, name, description, is_directed)
VALUES
('USED_IN', 'Used In', 'Specifies that an entity/product is used within a particular industry or application', true),
('APPLICABLE_TO', 'Applicable To', 'Contextual applicability to another domain or equipment', true),
('PART_OF', 'Part Of', 'Component or structural constituent of a parent assembly', true),
('COMPATIBLE_WITH', 'Compatible With', 'Interoperability or physical/chemical compatibility', false),
('ACCESSORY_FOR', 'Accessory For', 'Complementary accessory designed for primary equipment', true),
('ALTERNATIVE_TO', 'Alternative To', 'Direct functional substitute or equivalent entity', false),
('RELATED_TO', 'Related To', 'General contextual or associative semantic relationship', false),
('MANUFACTURED_FOR', 'Manufactured For', 'Specific target application or user persona requirement', true),
('REQUIRED_BY', 'Required By', 'Prerequisite material or tool required for a process', true),
('SUBTYPE_OF', 'Subtype Of', 'Specialized subtype or variant classification', true)
ON CONFLICT (code) DO NOTHING;

-- 10.3 Initial Migration of Existing category_hierarchy Nodes into Concepts and Placements
DO $$
DECLARE
    v_scheme_id INT;
    rec RECORD;
    v_concept_id INT;
    v_placement_id INT;
    v_parent_placement_id INT;
BEGIN
    SELECT id INTO v_scheme_id FROM classification_schemes WHERE code = 'SYS_PROD';

    IF EXISTS (SELECT 1 FROM classification_placements LIMIT 1) THEN
        RAISE NOTICE 'Classification hierarchy already populated, skipping initial loop.';
        RETURN;
    END IF;

    -- 1. Create concepts for all existing category_hierarchy nodes
    FOR rec IN SELECT * FROM category_hierarchy ORDER BY level ASC, id ASC LOOP
        INSERT INTO classification_concepts (
            code, canonical_name, name_en, name_am, name_om, name_ti,
            description, concept_type, status, metadata
        ) VALUES (
            rec.slug,
            rec.name_en,
            rec.name_en,
            rec.name_am,
            rec.name_om,
            rec.name_ti,
            rec.description,
            'PRODUCT',
            CASE WHEN rec.is_active THEN 'ACTIVE' ELSE 'ARCHIVED' END,
            jsonb_build_object('legacy_category_id', rec.id, 'icon', rec.icon)
        )
        ON CONFLICT (code) DO UPDATE SET
            canonical_name = EXCLUDED.canonical_name,
            name_en = EXCLUDED.name_en,
            name_am = EXCLUDED.name_am,
            name_om = EXCLUDED.name_om,
            name_ti = EXCLUDED.name_ti
        RETURNING id INTO v_concept_id;

        -- Find parent placement if rec has parent_id
        v_parent_placement_id := NULL;
        IF rec.parent_id IS NOT NULL THEN
            SELECT cp.id INTO v_parent_placement_id
            FROM classification_placements cp
            JOIN classification_concepts cc ON cc.id = cp.concept_id
            JOIN category_hierarchy parent_ch ON parent_ch.slug = cc.code
            WHERE cp.scheme_id = v_scheme_id AND parent_ch.id = rec.parent_id
            LIMIT 1;
        END IF;

        -- Create or update placement in SYS_PROD scheme
        INSERT INTO classification_placements (
            scheme_id, concept_id, parent_placement_id, role,
            slug, display_order, is_leaf, is_active, status
        ) VALUES (
            v_scheme_id,
            v_concept_id,
            v_parent_placement_id,
            CASE 
                WHEN rec.level = 0 THEN 'ROOT'
                WHEN rec.level = 1 THEN 'SUPERFAMILY'
                ELSE 'CATEGORY'
            END,
            rec.slug,
            rec.display_order,
            rec.is_leaf,
            rec.is_active,
            'PUBLISHED'
        )
        ON CONFLICT (scheme_id, slug) DO UPDATE SET
            parent_placement_id = EXCLUDED.parent_placement_id,
            role = EXCLUDED.role,
            is_leaf = EXCLUDED.is_leaf,
            is_active = EXCLUDED.is_active;
    END LOOP;

    -- 2. Migrate existing category_aliases into classification_aliases
    INSERT INTO classification_aliases (concept_id, alias, lang, alias_type, is_active)
    SELECT cc.id, ca.alias, 'en', 'SYNONYM', true
    FROM category_aliases ca
    JOIN category_hierarchy ch ON ch.id = ca.category_id
    JOIN classification_concepts cc ON cc.code = ch.slug
    ON CONFLICT DO NOTHING;

    -- 3. Migrate existing category_attributes into classification_attribute_rules
    INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
    SELECT cp.id, ca.attribute_id, ca.is_required, true, ca.display_order
    FROM category_attributes ca
    JOIN category_hierarchy ch ON ch.id = ca.category_id
    JOIN classification_placements cp ON cp.slug = ch.slug AND cp.scheme_id = v_scheme_id
    ON CONFLICT DO NOTHING;

    -- 4. Synchronize products table: link existing products to concept_id and placement_id
    UPDATE products p
    SET concept_id = cp.concept_id,
        placement_id = cp.id
    FROM classification_placements cp
    JOIN category_hierarchy ch ON ch.slug = cp.slug AND cp.scheme_id = v_scheme_id
    WHERE p.category_hierarchy_id = ch.id;

END $$;

