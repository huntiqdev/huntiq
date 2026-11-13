-- ============================================================================
-- 068_core_views.sql
-- Core platform views and underlying helper functions
-- ============================================================================

CREATE OR REPLACE FUNCTION public.get_category_ancestors(p_category_id integer)
 RETURNS TABLE(id integer, name_en character varying, slug character varying, level integer)
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    WITH RECURSIVE ancestors AS (
        SELECT 
            ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level
        FROM category_hierarchy ch
        WHERE ch.id = p_category_id
        
        UNION ALL
        
        SELECT 
            ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level
        FROM category_hierarchy ch
        JOIN ancestors a ON ch.id = a.parent_id
    )
    SELECT a.id, a.name_en, a.slug, a.level
    FROM ancestors a
    ORDER BY a.level;
END;
$function$;

CREATE OR REPLACE FUNCTION public.get_category_descendants(p_category_id integer)
 RETURNS TABLE(id integer, name_en character varying, slug character varying, level integer)
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    WITH RECURSIVE descendants AS (
        SELECT 
            ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level
        FROM category_hierarchy ch
        WHERE ch.id = p_category_id
        
        UNION ALL
        
        SELECT 
            ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level
        FROM category_hierarchy ch
        JOIN descendants d ON ch.parent_id = d.id
    )
    SELECT d.id, d.name_en, d.slug, d.level
    FROM descendants d
    ORDER BY d.level;
END;
$function$;

CREATE OR REPLACE FUNCTION public.get_category_full_name(p_category_id integer)
 RETURNS character varying
 LANGUAGE plpgsql
AS $function$
DECLARE
    v_full_name VARCHAR;
BEGIN
    SELECT ARRAY_TO_STRING(breadcrumb_names, ' > ')
    INTO v_full_name
    FROM category_breadcrumbs
    WHERE id = p_category_id;
    
    RETURN COALESCE(v_full_name, '');
END;
$function$;

CREATE OR REPLACE FUNCTION public.get_category_attributes(p_category_id integer)
 RETURNS TABLE(attribute_id integer, name_en character varying, slug character varying, data_type character varying, validation_rules jsonb, unit character varying, is_required boolean, is_filterable boolean, display_order integer)
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    SELECT 
        attr_inh.attribute_id,
        attr_inh.name_en,
        attr_inh.slug,
        attr_inh.data_type,
        attr_inh.validation_rules,
        attr_inh.unit,
        attr_inh.is_required,
        attr_inh.is_filterable,
        attr_inh.display_order
    FROM get_category_attributes_inherited(p_category_id) attr_inh;
END;
$function$;

CREATE OR REPLACE FUNCTION public.update_attribute_timestamp()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION public.update_product_attribute_search_vector()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    IF NEW.value_text IS NOT NULL THEN
        NEW.search_vector := to_tsvector('simple', COALESCE(NEW.value_text, ''));
    END IF;
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION public.user_has_capability(p_user_id integer, p_capability character varying)
 RETURNS boolean
 LANGUAGE plpgsql
 STABLE
AS $function$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM user_capabilities
        WHERE user_id = p_user_id
          AND capability = p_capability
          AND status = 'active'
    );
END;
$function$;

CREATE OR REPLACE FUNCTION public.get_user_capabilities(p_user_id integer)
 RETURNS text[]
 LANGUAGE plpgsql
 STABLE
AS $function$
BEGIN
    RETURN ARRAY(
        SELECT capability FROM user_capabilities
        WHERE user_id = p_user_id AND status = 'active'
        ORDER BY capability ASC
    );
END;
$function$;

DROP TRIGGER IF EXISTS trg_update_attribute_definitions_timestamp ON attribute_definitions;
CREATE TRIGGER trg_update_attribute_definitions_timestamp BEFORE UPDATE ON attribute_definitions FOR EACH ROW EXECUTE FUNCTION update_attribute_timestamp();

DROP TRIGGER IF EXISTS trg_update_product_attribute_search ON product_attributes;
CREATE TRIGGER trg_update_product_attribute_search BEFORE UPDATE ON product_attributes FOR EACH ROW EXECUTE FUNCTION update_product_attribute_search_vector();

CREATE OR REPLACE VIEW referral_balance_view AS
SELECT u.id AS user_id,
    COALESCE(sum(rc.commission_amount), (0)::numeric) AS total_earned,
    COALESCE(( SELECT sum(rw.amount) AS sum
           FROM referral_withdrawals rw
          WHERE ((rw.user_id = u.id) AND ((rw.status)::text = ANY ((ARRAY['approved'::character varying, 'paid'::character varying])::text[])))), (0)::numeric) AS total_withdrawn,
    (COALESCE(sum(rc.commission_amount), (0)::numeric) - COALESCE(( SELECT sum(rw.amount) AS sum
           FROM referral_withdrawals rw
          WHERE ((rw.user_id = u.id) AND ((rw.status)::text = ANY ((ARRAY['approved'::character varying, 'paid'::character varying, 'pending'::character varying])::text[])))), (0)::numeric)) AS available_balance
   FROM (users u
     LEFT JOIN referral_commissions rc ON (((rc.referrer_id = u.id) AND ((rc.status)::text = 'confirmed'::text))))
  GROUP BY u.id;

CREATE OR REPLACE VIEW otp_log_view AS
SELECT id,
    COALESCE(phone, ''::character varying) AS phone,
    COALESCE(email, ''::character varying) AS email,
    purpose,
    user_type,
    used,
    attempts,
    expires_at,
    created_at,
    status,
    COALESCE(code_encrypted, code, ''::text) AS encrypted_code
   FROM otp_codes;

CREATE OR REPLACE VIEW category_breadcrumbs AS
WITH RECURSIVE category_path AS (
         SELECT category_hierarchy.id,
            category_hierarchy.parent_id,
            category_hierarchy.name_en,
            category_hierarchy.slug,
            category_hierarchy.level,
            category_hierarchy.path,
            ARRAY[(category_hierarchy.name_en)::text] AS breadcrumb_names,
            ARRAY[(category_hierarchy.slug)::text] AS breadcrumb_slugs,
            ARRAY[category_hierarchy.id] AS breadcrumb_ids
           FROM category_hierarchy
          WHERE (category_hierarchy.parent_id IS NULL)
        UNION ALL
         SELECT ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level,
            ch.path,
            (cp.breadcrumb_names || (ch.name_en)::text),
            (cp.breadcrumb_slugs || (ch.slug)::text),
            (cp.breadcrumb_ids || ch.id)
           FROM (category_hierarchy ch
             JOIN category_path cp ON ((ch.parent_id = cp.id)))
        )
 SELECT id,
    name_en,
    slug,
    level,
    breadcrumb_names,
    breadcrumb_slugs,
    breadcrumb_ids,
    array_length(breadcrumb_names, 1) AS depth
   FROM category_path;

CREATE OR REPLACE VIEW category_stats AS
SELECT ch.id,
    ch.id AS category_id,
    ch.parent_id,
    ch.name_en,
    ch.slug,
    ch.level,
    ch.is_leaf,
    count(p.id) AS product_count,
    count(p.id) FILTER (WHERE (p.status = 'approved'::product_status)) AS approved_count,
    count(p.id) FILTER (WHERE (p.status = 'pending'::product_status)) AS pending_count
   FROM (category_hierarchy ch
     LEFT JOIN products p ON ((p.category_hierarchy_id = ch.id)))
  GROUP BY ch.id, ch.parent_id, ch.name_en, ch.slug, ch.level, ch.is_leaf;

CREATE OR REPLACE VIEW category_applicable_attributes AS
SELECT ca.category_id,
    ch.name_en AS category_name,
    ch.slug AS category_slug,
    ad.id AS attribute_id,
    ad.name_en AS attribute_name,
    ad.slug AS attribute_slug,
    ad.data_type,
    ad.validation_rules,
    ad.unit,
    ca.is_required,
    ca.display_order,
    ad.is_filterable,
    ad.is_searchable
   FROM ((category_attributes ca
     JOIN category_hierarchy ch ON ((ch.id = ca.category_id)))
     JOIN attribute_definitions ad ON ((ad.id = ca.attribute_id)))
  ORDER BY ca.category_id, ca.display_order;

CREATE OR REPLACE VIEW user_capabilities_view AS
SELECT u.id AS user_id,
    u.name,
    u.phone,
    u.email,
    u.role,
    u.tier AS seller_tier,
    u.buyer_tier,
    u.admin_role,
    u.is_active,
    COALESCE(array_agg(uc.capability) FILTER (WHERE ((uc.status)::text = 'active'::text)), ARRAY['buyer'::character varying]) AS active_capabilities
   FROM (users u
     LEFT JOIN user_capabilities uc ON ((uc.user_id = u.id)))
  GROUP BY u.id;
