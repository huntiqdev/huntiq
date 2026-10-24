-- ============================================================================
-- Migration 031: Universal Marketplace Taxonomy, Catalogs & Recruitment System
-- ============================================================================
-- Purpose:
-- 1. Add Tigrinya (name_ti) localization support to category and attribute systems.
-- 2. Create Brands, Product Families, and Product Models catalog tables.
-- 3. Create Occupations, Qualification Types, and Recruitment/CV privacy tables.
-- 4. Create Category Aliases for multi-lingual search understanding.
-- 5. Add Brand and Model foreign keys to products table.
-- 6. Upgrade get_category_attributes() function with full ancestor path inheritance.
-- ============================================================================

-- ============================================================================
-- PART 1: Multi-Lingual Column Upgrades (Tigrinya Support)
-- ============================================================================

ALTER TABLE category_hierarchy 
    ADD COLUMN IF NOT EXISTS name_ti VARCHAR(200);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'uq_category_hierarchy_slug'
    ) THEN
        ALTER TABLE category_hierarchy ADD CONSTRAINT uq_category_hierarchy_slug UNIQUE (slug);
    END IF;
END $$;

ALTER TABLE attribute_definitions 
    ADD COLUMN IF NOT EXISTS name_ti VARCHAR(200);

-- ============================================================================
-- PART 2: Brands and Manufacturers Catalog
-- ============================================================================

CREATE TABLE IF NOT EXISTS brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    category_ids INT[],                               -- Applicable category_hierarchy IDs
    logo_url TEXT,
    country_of_origin VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_brands_slug ON brands(slug);
CREATE INDEX IF NOT EXISTS idx_brands_active ON brands(is_active);
CREATE INDEX IF NOT EXISTS idx_brands_category_ids ON brands USING gin(category_ids);

-- ============================================================================
-- PART 3: Product Families and Series
-- ============================================================================

CREATE TABLE IF NOT EXISTS product_families (
    id SERIAL PRIMARY KEY,
    brand_id INT NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
    category_id INT REFERENCES category_hierarchy(id) ON DELETE SET NULL,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL,
    series_name VARCHAR(100),                         -- e.g. "Galaxy S Series", "ThinkPad T Series"
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_brand_family_slug UNIQUE(brand_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_product_families_brand ON product_families(brand_id);
CREATE INDEX IF NOT EXISTS idx_product_families_category ON product_families(category_id);
CREATE INDEX IF NOT EXISTS idx_product_families_slug ON product_families(slug);

-- ============================================================================
-- PART 4: Exact Product Models Catalog
-- ============================================================================

CREATE TABLE IF NOT EXISTS product_models (
    id SERIAL PRIMARY KEY,
    family_id INT REFERENCES product_families(id) ON DELETE SET NULL,
    brand_id INT NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
    category_id INT REFERENCES category_hierarchy(id) ON DELETE SET NULL,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL,
    model_number VARCHAR(100),                        -- Manufacturer SKU/code, e.g. "SM-S928B"
    release_year INT,
    default_specs JSONB,                              -- Prepopulated hardware specs (RAM, Storage, Engine, etc.)
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_brand_model_slug UNIQUE(brand_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_product_models_family ON product_models(family_id);
CREATE INDEX IF NOT EXISTS idx_product_models_brand ON product_models(brand_id);
CREATE INDEX IF NOT EXISTS idx_product_models_category ON product_models(category_id);
CREATE INDEX IF NOT EXISTS idx_product_models_slug ON product_models(slug);
CREATE INDEX IF NOT EXISTS idx_product_models_year ON product_models(release_year);
CREATE INDEX IF NOT EXISTS idx_product_models_specs ON product_models USING gin(default_specs);

-- ============================================================================
-- PART 5: Professional Occupations & Qualification Taxonomy
-- ============================================================================

CREATE TABLE IF NOT EXISTS occupations (
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES category_hierarchy(id) ON DELETE SET NULL,
    title_en VARCHAR(200) NOT NULL,
    title_am VARCHAR(200),
    title_om VARCHAR(200),
    title_ti VARCHAR(200),
    slug VARCHAR(200) UNIQUE NOT NULL,
    code VARCHAR(50),                                 -- Optional ISCO-08 occupational code
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_occupations_slug ON occupations(slug);
CREATE INDEX IF NOT EXISTS idx_occupations_category ON occupations(category_id);

CREATE TABLE IF NOT EXISTS qualification_types (
    id SERIAL PRIMARY KEY,
    name_en VARCHAR(200) NOT NULL,
    name_am VARCHAR(200),
    name_om VARCHAR(200),
    name_ti VARCHAR(200),
    level_order INT NOT NULL DEFAULT 0,               -- Sorting by education level (1=Certificate, 5=PhD)
    category VARCHAR(50) NOT NULL DEFAULT 'degree',   -- 'degree', 'tvet', 'diploma', 'certificate', 'license'
    is_active BOOLEAN DEFAULT true
);

-- ============================================================================
-- PART 6: Job Postings, Candidate Profiles & Private CV Vault
-- ============================================================================

-- Table: job_postings (Extension of products table for employment domain)
CREATE TABLE IF NOT EXISTS job_postings (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE UNIQUE,
    job_title VARCHAR(250) NOT NULL,
    occupation_id INT REFERENCES occupations(id) ON DELETE SET NULL,
    employer_name VARCHAR(250) NOT NULL,
    workplace_type VARCHAR(50) DEFAULT 'on_site',     -- on_site, hybrid, remote
    employment_type VARCHAR(50) DEFAULT 'full_time',   -- full_time, contract, part_time, internship, freelance
    min_salary NUMERIC(15,2),
    max_salary NUMERIC(15,2),
    currency VARCHAR(10) DEFAULT 'ETB',
    education_level VARCHAR(100),
    experience_years INT DEFAULT 0,
    application_deadline DATE,
    application_instructions TEXT,
    skills_required TEXT[],
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_job_postings_occupation ON job_postings(occupation_id);
CREATE INDEX IF NOT EXISTS idx_job_postings_employment_type ON job_postings(employment_type);
CREATE INDEX IF NOT EXISTS idx_job_postings_workplace_type ON job_postings(workplace_type);
CREATE INDEX IF NOT EXISTS idx_job_postings_deadline ON job_postings(application_deadline);

-- Table: candidate_profiles (Public professional profiles)
CREATE TABLE IF NOT EXISTS candidate_profiles (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    professional_title VARCHAR(200) NOT NULL,
    occupation_id INT REFERENCES occupations(id) ON DELETE SET NULL,
    headline TEXT,
    experience_years INT DEFAULT 0,
    highest_qualification_id INT REFERENCES qualification_types(id) ON DELETE SET NULL,
    skills TEXT[],
    languages TEXT[],
    portfolio_url TEXT,
    city VARCHAR(100),
    is_profile_public BOOLEAN DEFAULT true,
    is_seeking_work BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_candidate_profiles_occupation ON candidate_profiles(occupation_id);
CREATE INDEX IF NOT EXISTS idx_candidate_profiles_seeking ON candidate_profiles(is_seeking_work);

-- Table: candidate_qualifications (Degrees, TVET, Diplomas)
CREATE TABLE IF NOT EXISTS candidate_qualifications (
    id SERIAL PRIMARY KEY,
    candidate_id INT NOT NULL REFERENCES candidate_profiles(id) ON DELETE CASCADE,
    qualification_type_id INT REFERENCES qualification_types(id) ON DELETE SET NULL,
    field_of_study VARCHAR(200) NOT NULL,
    institution VARCHAR(200) NOT NULL,
    graduation_year INT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_candidate_qualifications_cand ON candidate_qualifications(candidate_id);

-- Table: cv_documents (Private CV storage - Never public by default)
CREATE TABLE IF NOT EXISTS cv_documents (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_url TEXT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size INT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    is_private BOOLEAN DEFAULT true,
    download_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cv_documents_user ON cv_documents(user_id);

-- Table: cv_access_requests (Explicit candidate consent workflow)
CREATE TABLE IF NOT EXISTS cv_access_requests (
    id SERIAL PRIMARY KEY,
    candidate_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    employer_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_posting_id INT REFERENCES job_postings(id) ON DELETE SET NULL,
    status VARCHAR(50) DEFAULT 'pending',             -- pending, approved, declined, revoked
    request_note TEXT,
    requested_at TIMESTAMP DEFAULT NOW(),
    responded_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cv_access_candidate ON cv_access_requests(candidate_id);
CREATE INDEX IF NOT EXISTS idx_cv_access_employer ON cv_access_requests(employer_id);

-- ============================================================================
-- PART 7: Multi-Lingual Category Aliases (Synonyms & Transliterations)
-- ============================================================================

CREATE TABLE IF NOT EXISTS category_aliases (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES category_hierarchy(id) ON DELETE CASCADE,
    alias VARCHAR(200) NOT NULL,
    language VARCHAR(10) DEFAULT 'en',                -- 'en', 'am', 'om', 'ti'
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_category_aliases_category ON category_aliases(category_id);
CREATE INDEX IF NOT EXISTS idx_category_aliases_alias ON category_aliases(alias);

-- ============================================================================
-- PART 8: Add Brand and Model Relationships to Products Table
-- ============================================================================

ALTER TABLE products 
    ADD COLUMN IF NOT EXISTS brand_id INT REFERENCES brands(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS product_model_id INT REFERENCES product_models(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS custom_model_name VARCHAR(200);

CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand_id);
CREATE INDEX IF NOT EXISTS idx_products_model ON products(product_model_id);

-- ============================================================================
-- PART 9: Upgraded Dynamic Attribute Inheritance Function
-- ============================================================================

-- Function: get_category_attributes_inherited
-- Retrieves all attributes applicable to a category by walking the ancestor path
CREATE OR REPLACE FUNCTION get_category_attributes_inherited(p_category_id INT)
RETURNS TABLE (
    attribute_id INT,
    name_en VARCHAR,
    name_am VARCHAR,
    name_om VARCHAR,
    name_ti VARCHAR,
    slug VARCHAR,
    data_type VARCHAR,
    validation_rules JSONB,
    unit VARCHAR,
    is_required BOOLEAN,
    is_filterable BOOLEAN,
    is_searchable BOOLEAN,
    is_inherited BOOLEAN,
    display_order INT
) AS $$
DECLARE
    v_path VARCHAR;
    v_ancestor_ids INT[];
BEGIN
    -- Get path of target category
    SELECT path INTO v_path
    FROM category_hierarchy
    WHERE id = p_category_id;
    
    IF v_path IS NOT NULL AND v_path != '' THEN
        -- Convert '/1/5/12/' into integer array [1, 5, 12]
        v_ancestor_ids := string_to_array(trim(both '/' from v_path), '/')::int[];
    ELSE
        v_ancestor_ids := ARRAY[p_category_id];
    END IF;

    RETURN QUERY
    WITH candidate_attributes AS (
        SELECT 
            ad.id AS attribute_id,
            ad.name_en,
            ad.name_am,
            ad.name_om,
            ad.name_ti,
            ad.slug,
            ad.data_type,
            ad.validation_rules,
            ad.unit,
            ca.is_required,
            ad.is_filterable,
            ad.is_searchable,
            (ca.category_id != p_category_id) AS is_inherited,
            ca.display_order,
            ca.category_id AS assigned_category_id,
            ROW_NUMBER() OVER (
                PARTITION BY ad.id 
                ORDER BY 
                    -- Prefer exact category assignment over ancestor assignments
                    CASE WHEN ca.category_id = p_category_id THEN 0 ELSE 1 END,
                    ca.display_order ASC
            ) AS rn
        FROM category_attributes ca
        JOIN attribute_definitions ad ON ad.id = ca.attribute_id
        WHERE ca.category_id = ANY(v_ancestor_ids) 
           OR ca.category_id = p_category_id
    )
    SELECT 
        ca.attribute_id,
        ca.name_en,
        ca.name_am,
        ca.name_om,
        ca.name_ti,
        ca.slug,
        ca.data_type,
        ca.validation_rules,
        ca.unit,
        ca.is_required,
        ca.is_filterable,
        ca.is_searchable,
        ca.is_inherited,
        ca.display_order
    FROM candidate_attributes ca
    WHERE ca.rn = 1
    ORDER BY ca.display_order ASC, ca.attribute_id ASC;
END;
$$ LANGUAGE plpgsql;

-- Also update original get_category_attributes to use inherited logic for full backward compatibility
CREATE OR REPLACE FUNCTION get_category_attributes(p_category_id INT)
RETURNS TABLE (
    attribute_id INT,
    name_en VARCHAR,
    slug VARCHAR,
    data_type VARCHAR,
    validation_rules JSONB,
    unit VARCHAR,
    is_required BOOLEAN,
    is_filterable BOOLEAN,
    display_order INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ai.attribute_id,
        ai.name_en,
        ai.slug,
        ai.data_type,
        ai.validation_rules,
        ai.unit,
        ai.is_required,
        ai.is_filterable,
        ai.display_order
    FROM get_category_attributes_inherited(p_category_id) ai;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- Migration 031 Complete
-- ============================================================================
