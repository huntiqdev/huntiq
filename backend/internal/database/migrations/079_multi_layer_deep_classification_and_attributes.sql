-- ============================================================================
-- 079_multi_layer_deep_classification_and_attributes.sql
-- Expand hierarchical classification through Layers 3, 4, 5, and 6
-- Inject cross-sector relationship graph mesh and multi-tier attribute rules
-- ============================================================================

BEGIN;

-- 1. Ensure placement role check constraint includes all 7 layers (0 to 6)
ALTER TABLE classification_placements DROP CONSTRAINT IF EXISTS chk_placement_role;
ALTER TABLE classification_placements ADD CONSTRAINT chk_placement_role
    CHECK (role IN ('ROOT', 'SUPERFAMILY', 'CATEGORY', 'SUBCATEGORY', 'SPECIFICATION', 'VARIANT', 'LEAF', 'OTHER'));

-- 2. Expand Power Tools under Tools, Machinery & Equipment to Layer 6
DO $$
DECLARE
    tools_scheme_id INT;
    drill_cat_id    INT;
    tools_root_pid  INT;
    c_cat_id        INT;
    l3_id           INT;
    l4_id           INT;
    l5_id           INT;
    l6_id           INT;
    
    c_l3_id INT; c_l4_id INT; c_l5_id INT; c_l6_id INT;
BEGIN
    -- Ensure SYS_TOOLS scheme exists
    INSERT INTO classification_schemes (code, name_en, name_am, name_om, name_ti, description, is_system, status)
    VALUES ('SYS_TOOLS', 'Tools & Machinery Classification', 'የመሳሪያዎችና ማሽነሪዎች ምደባ', 'Ramaddii Meeshaalee fi Maashinootaa', 'ምደባ መሳርሒታትን ማሽነሪታትን', 'Universal tools, power equipment, and machinery taxonomy', false, 'ACTIVE')
    ON CONFLICT (code) DO NOTHING;

    SELECT id INTO tools_scheme_id FROM classification_schemes WHERE code = 'SYS_TOOLS';

    -- Ensure a root concept & placement exists for SYS_TOOLS
    IF NOT EXISTS (SELECT 1 FROM classification_placements WHERE scheme_id = tools_scheme_id AND role = 'ROOT') THEN
        INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
        VALUES ('CONCEPT_TOOLS_MACHINERY', 'Tools, Machinery & Equipment', 'Tools, Machinery & Equipment', 'የመሳሪያዎችና ማሽነሪዎች', 'Meeshaalee fi Maashinootaa', 'መሳርሒታትን ማሽነሪታትን', 'Universal tools and machinery root', 'EQUIPMENT', 'ACTIVE')
        ON CONFLICT (code) DO NOTHING;

        INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, path, depth, is_leaf, status)
        SELECT tools_scheme_id, id, NULL, 'ROOT', 'tools-machinery', '/tools-machinery', 0, false, 'PUBLISHED'
        FROM classification_concepts WHERE code = 'CONCEPT_TOOLS_MACHINERY'
        ON CONFLICT (scheme_id, slug) DO NOTHING;
    END IF;
    
    -- Find Category (Layer 2) for Power Drills in SYS_TOOLS
    SELECT id INTO drill_cat_id FROM classification_placements 
    WHERE scheme_id = tools_scheme_id AND slug IN ('tools-xref-power-drills-home', 'tools-xref-drills-power') LIMIT 1;
    
    IF drill_cat_id IS NULL THEN
        SELECT id INTO drill_cat_id FROM classification_placements 
        WHERE scheme_id = tools_scheme_id AND role = 'CATEGORY' LIMIT 1;
    END IF;

    IF drill_cat_id IS NULL THEN
        SELECT id INTO tools_root_pid FROM classification_placements WHERE scheme_id = tools_scheme_id AND role = 'ROOT' LIMIT 1;
        INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, description, concept_type, status)
        VALUES ('CONCEPT_POWER_DRILLS', 'Power Drills', 'Power Drills', 'የኤሌክትሪክ ድሪሎች', 'Power rotary drills', 'EQUIPMENT', 'ACTIVE')
        ON CONFLICT (code) DO UPDATE SET name_en = EXCLUDED.name_en RETURNING id INTO c_cat_id;

        INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, path, depth, is_leaf, status)
        VALUES (tools_scheme_id, c_cat_id, tools_root_pid, 'CATEGORY', 'tools-xref-power-drills-home', '/tools-machinery/power-drills', 2, false, 'PUBLISHED')
        ON CONFLICT (scheme_id, slug) DO UPDATE SET role = 'CATEGORY' RETURNING id INTO drill_cat_id;
    END IF;

    -- --- LAYER 3: SUBCATEGORY (Power Drills & Drivers) ---
    INSERT INTO classification_concepts (code, canonical_name, name_en, description, concept_type, status)
    VALUES ('CONCEPT_L3_DRILLS_DRIVERS', 'Power Drills & Drivers', 'Power Drills & Drivers', 'Rotary power drilling and fastening tools', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET name_en = EXCLUDED.name_en RETURNING id INTO c_l3_id;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, path, depth, is_leaf, status)
    VALUES (tools_scheme_id, c_l3_id, drill_cat_id, 'SUBCATEGORY', 'l3-power-drills-drivers', '/tools/drills/drivers', 3, false, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET depth = 3 RETURNING id INTO l3_id;

    -- --- LAYER 4: SPECIFICATION (Rotary Hammer & SDS Drills) ---
    INSERT INTO classification_concepts (code, canonical_name, name_en, description, concept_type, status)
    VALUES ('CONCEPT_L4_ROTARY_HAMMER', 'Rotary Hammer & SDS Drills', 'Rotary Hammer & SDS Drills', 'Heavy duty electropneumatic hammer drills for concrete & masonry', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET name_en = EXCLUDED.name_en RETURNING id INTO c_l4_id;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, path, depth, is_leaf, status)
    VALUES (tools_scheme_id, c_l4_id, l3_id, 'SPECIFICATION', 'l4-rotary-hammer-sds', '/tools/drills/drivers/rotary-hammer', 4, false, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET depth = 4 RETURNING id INTO l4_id;

    -- --- LAYER 5: VARIANT (Cordless Brushless SDS-Plus 18V-20V) ---
    INSERT INTO classification_concepts (code, canonical_name, name_en, description, concept_type, status)
    VALUES ('CONCEPT_L5_CORDLESS_SDS_PLUS', 'Cordless Brushless SDS-Plus 18V-20V', 'Cordless Brushless SDS-Plus 18V-20V', 'High efficiency brushless cordless SDS-plus concrete drills', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET name_en = EXCLUDED.name_en RETURNING id INTO c_l5_id;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, path, depth, is_leaf, status)
    VALUES (tools_scheme_id, c_l5_id, l4_id, 'VARIANT', 'l5-cordless-sds-plus-18v', '/tools/drills/drivers/rotary-hammer/cordless-sds-plus', 5, false, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET depth = 5 RETURNING id INTO l5_id;

    -- --- LAYER 6: LEAF (Heavy Duty 3-Mode SDS-Plus Hammer Drill 3.2J) ---
    INSERT INTO classification_concepts (code, canonical_name, name_en, description, concept_type, status)
    VALUES ('CONCEPT_L6_SDS_PLUS_3MODE_32J', 'Heavy Duty 3-Mode SDS-Plus Hammer Drill 3.2J', 'Heavy Duty 3-Mode SDS-Plus Hammer Drill 3.2J', 'Industrial 3-mode rotary hammer drill with 3.2 Joules impact energy', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET name_en = EXCLUDED.name_en RETURNING id INTO c_l6_id;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, path, depth, is_leaf, status)
    VALUES (tools_scheme_id, c_l6_id, l5_id, 'LEAF', 'l6-sds-plus-3mode-32j', '/tools/drills/drivers/rotary-hammer/cordless-sds-plus/3mode-32j', 6, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET depth = 6, is_leaf = true RETURNING id INTO l6_id;

END $$;

-- 3. Cross-Sector Mapping: Link this specific tool into Construction, Manufacturing, and Electrical sectors
DO $$
DECLARE
    drill_concept_id INT;
    const_root_cid   INT;
    mfg_root_cid     INT;
    elec_root_cid    INT;
    tools_scheme_id  INT;
BEGIN
    SELECT id INTO drill_concept_id FROM classification_concepts WHERE code = 'CONCEPT_L6_SDS_PLUS_3MODE_32J';
    SELECT concept_id INTO const_root_cid FROM classification_placements WHERE slug = 'construction-infrastructure' AND role = 'ROOT';
    SELECT concept_id INTO mfg_root_cid   FROM classification_placements WHERE slug = 'industrial-manufacturing' AND role = 'ROOT';
    SELECT concept_id INTO elec_root_cid  FROM classification_placements WHERE slug = 'electronics' AND role = 'ROOT';
    SELECT id INTO tools_scheme_id        FROM classification_schemes WHERE code = 'SYS_TOOLS';

    IF drill_concept_id IS NOT NULL THEN
        IF const_root_cid IS NOT NULL THEN
            INSERT INTO classification_relationships (source_concept_id, target_concept_id, relationship_type_code, context_scheme_id)
            VALUES (drill_concept_id, const_root_cid, 'USED_IN', tools_scheme_id)
            ON CONFLICT ON CONSTRAINT uq_concept_rel DO NOTHING;
        END IF;

        IF mfg_root_cid IS NOT NULL THEN
            INSERT INTO classification_relationships (source_concept_id, target_concept_id, relationship_type_code, context_scheme_id)
            VALUES (drill_concept_id, mfg_root_cid, 'USED_IN', tools_scheme_id)
            ON CONFLICT ON CONSTRAINT uq_concept_rel DO NOTHING;
        END IF;

        IF elec_root_cid IS NOT NULL THEN
            INSERT INTO classification_relationships (source_concept_id, target_concept_id, relationship_type_code, context_scheme_id)
            VALUES (drill_concept_id, elec_root_cid, 'APPLICABLE_TO', tools_scheme_id)
            ON CONFLICT ON CONSTRAINT uq_concept_rel DO NOTHING;
        END IF;
    END IF;
END $$;

-- 4. Inject Multi-Tier Attribute Rules (Power, Voltage, Chuck Size, Impact Energy)
DO $$
DECLARE
    l4_pid INT;
    l5_pid INT;
    l6_pid INT;
    attr_chuck_id   INT;
    attr_voltage_id INT;
    attr_impact_id  INT;
BEGIN
    SELECT id INTO l4_pid FROM classification_placements WHERE slug = 'l4-rotary-hammer-sds';
    SELECT id INTO l5_pid FROM classification_placements WHERE slug = 'l5-cordless-sds-plus-18v';
    SELECT id INTO l6_pid FROM classification_placements WHERE slug = 'l6-sds-plus-3mode-32j';

    -- Upsert attribute definitions
    INSERT INTO attribute_definitions (name_en, slug, data_type, validation_rules, is_filterable, is_searchable)
    VALUES ('Chuck Type', 'attr-chuck-type', 'select', '{"options": ["SDS-Plus", "SDS-Max", "Spline", "Keyless Chuck"]}', true, true)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en RETURNING id INTO attr_chuck_id;

    INSERT INTO attribute_definitions (name_en, slug, data_type, validation_rules, is_filterable, is_searchable)
    VALUES ('Operating Voltage', 'attr-operating-voltage', 'select', '{"options": ["12V", "18V", "20V MAX", "36V", "40V MAX", "54V/60V MAX"]}', true, true)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en RETURNING id INTO attr_voltage_id;

    INSERT INTO attribute_definitions (name_en, slug, data_type, validation_rules, unit, is_filterable, is_searchable)
    VALUES ('Impact Energy', 'attr-impact-energy', 'number', '{"min": 0.5, "max": 25.0, "step": 0.1}', 'Joules', true, true)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en RETURNING id INTO attr_impact_id;

    -- Bind attributes to Layer 4, Layer 5, and Layer 6
    IF l4_pid IS NOT NULL AND attr_chuck_id IS NOT NULL THEN
        INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
        VALUES (l4_pid, attr_chuck_id, true, true, 1)
        ON CONFLICT (placement_id, attribute_id) DO NOTHING;
    END IF;

    IF l5_pid IS NOT NULL AND attr_voltage_id IS NOT NULL THEN
        INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
        VALUES (l5_pid, attr_voltage_id, true, true, 2)
        ON CONFLICT (placement_id, attribute_id) DO NOTHING;
    END IF;

    IF l6_pid IS NOT NULL AND attr_impact_id IS NOT NULL THEN
        INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
        VALUES (l6_pid, attr_impact_id, true, true, 3)
        ON CONFLICT (placement_id, attribute_id) DO NOTHING;
    END IF;
END $$;

-- 5. Mark non-parent placements as is_leaf efficiently
UPDATE classification_placements p
SET is_leaf = false
WHERE is_leaf = true
  AND EXISTS (SELECT 1 FROM classification_placements child WHERE child.parent_placement_id = p.id);

UPDATE classification_placements p
SET is_leaf = true
WHERE is_leaf = false
  AND NOT EXISTS (SELECT 1 FROM classification_placements child WHERE child.parent_placement_id = p.id);

COMMIT;
