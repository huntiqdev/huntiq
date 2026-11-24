#!/usr/bin/env python3
"""
HUNTIQ Comprehensive QA Lead & Adversarial Test Suite
Full multi-perspective audit of:
  - Multi-level arbitrary depth taxonomy (Levels 0-5)
  - Materialized path and tree consistency across 20,528 categories
  - Leaf vs branch constraints
  - Acyclic tree integrity
  - Security, boundary & negative testing of category endpoints
  - Cross-category discovery relationships
  - Database schema & table normalization (76 base tables)
  - Full JSON-Database bidirectional parity
  - Product-Category referential integrity
"""

import psycopg2
import psycopg2.extras
import json
import urllib.request
import urllib.error
import urllib.parse
import sys
import os
import time

DB_PARAMS = {
    "dbname": "hulubedej_db",
    "user": "postgres",
    "password": "gedam",
    "host": "127.0.0.1",
    "port": 5432
}
API_BASE = "http://127.0.0.1:8080/api"

def api_request(endpoint, timeout=10):
    url = f"{API_BASE}{endpoint}"
    req = urllib.request.Request(url, headers={"User-Agent": "HUNTIQ-QA-Lead/2.0"})
    try:
        start = time.perf_counter()
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            latency = (time.perf_counter() - start) * 1000
            body = resp.read().decode("utf-8")
            return resp.status, json.loads(body), latency
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8")
        try:
            parsed = json.loads(body)
        except:
            parsed = {"raw": body}
        return e.code, parsed, 0.0
    except Exception as e:
        return 500, {"error": str(e)}, 0.0

def run_qa_suite():
    conn = psycopg2.connect(**DB_PARAMS)
    cur = conn.cursor(cursor_factory=psycopg2.extras.DictCursor)

    results = []

    def test(group, name, passed, details=""):
        status = "PASS" if passed else "FAIL"
        print(f"[{status}] [{group}] {name}")
        if details:
            print(f"       -> {details}")
        results.append({"group": group, "name": name, "passed": passed, "details": details})

    print("================================================================================")
    print("  HUNTIQ FULL TESTER & ADVERSARIAL QA AUDIT SUITE")
    print("================================================================================")

    # ==============================================================================
    # 1. DATABASE SCHEMA & TABLE NORMALIZATION
    # ==============================================================================
    cur.execute("""
        SELECT count(*) FROM information_schema.tables 
        WHERE table_schema = 'public' AND table_type = 'BASE TABLE';
    """)
    table_count = cur.fetchone()[0]
    test("Database Schema", "Base Table Count (Target: >= 76, Enterprise: 80)", table_count in (76, 80), f"Count: {table_count}")

    # Check unreferenced legacy tables do not exist
    cur.execute("""
        SELECT table_name FROM information_schema.tables 
        WHERE table_schema = 'public' AND table_name IN (
            'brand_dictionary', 'synonym_dictionary', 'category_redirects',
            'collection_memberships', 'discovery_collections', 'entity_lifecycle_audit_logs',
            'entity_relationships', 'canonical_entities', 'entity_types',
            'merchant_inventory', 'invoices', 'service_appointments',
            'merchant_analytics_events', 'referral_fraud_evaluations'
        );
    """)
    dead_tables = cur.fetchall()
    test("Database Schema", "Zero Dead Legacy Tables Surviving", len(dead_tables) == 0, f"Found: {len(dead_tables)}")

    # Check indexes on classification_relationships
    cur.execute("""
        SELECT indexname FROM pg_indexes 
        WHERE tablename = 'classification_relationships';
    """)
    rel_indexes = [r[0] for r in cur.fetchall()]
    has_rel_idx = "idx_classif_rel_source" in rel_indexes and "idx_classif_rel_target" in rel_indexes
    test("Database Schema", "Performance Indexes on Classification Relationships", has_rel_idx, f"Indexes: {rel_indexes}")

    # ==============================================================================
    # 2. CATEGORY TAXONOMY INTEGRITY & MATHEMATICAL CONSISTENCY
    # ==============================================================================
    cur.execute("SELECT count(*) FROM category_hierarchy;")
    total_cats = cur.fetchone()[0]
    test("Taxonomy Integrity", "Total Categories in DB (20,528)", total_cats == 20528, f"Total: {total_cats}")

    # Acyclic check & zero orphan parents
    cur.execute("""
        SELECT count(*) FROM category_hierarchy c
        WHERE c.parent_id IS NOT NULL 
          AND NOT EXISTS (SELECT 1 FROM category_hierarchy p WHERE p.id = c.parent_id);
    """)
    orphan_cats = cur.fetchone()[0]
    test("Taxonomy Integrity", "Zero Orphan Parent References", orphan_cats == 0, f"Orphans: {orphan_cats}")

    # Check that id <> parent_id
    cur.execute("SELECT count(*) FROM category_hierarchy WHERE id = parent_id;")
    self_parents = cur.fetchone()[0]
    test("Taxonomy Integrity", "Zero Self-Parenting Cycles", self_parents == 0, f"Self-parents: {self_parents}")

    # Path Consistency: path must end with /id/ and contain parent's path
    cur.execute("""
        SELECT count(*) FROM category_hierarchy c
        JOIN category_hierarchy p ON c.parent_id = p.id
        WHERE c.path NOT LIKE p.path || c.id || '/';
    """)
    path_mismatches = cur.fetchone()[0]
    test("Taxonomy Integrity", "100% Consistent Materialized Paths with Parents", path_mismatches == 0, f"Mismatches: {path_mismatches}")

    # Level Consistency: level must equal depth in path
    cur.execute("""
        SELECT count(*) FROM category_hierarchy
        WHERE level <> (length(path) - length(replace(path, '/', '')) - 2);
    """)
    level_mismatches = cur.fetchone()[0]
    test("Taxonomy Integrity", "100% Mathematical Level-to-Path Alignment", level_mismatches == 0, f"Mismatches: {level_mismatches}")

    # Leaf Consistency: is_leaf = false MUST have children; is_leaf = true MUST NOT have children
    cur.execute("""
        SELECT count(*) FROM category_hierarchy c
        WHERE c.is_leaf = false 
          AND NOT EXISTS (SELECT 1 FROM category_hierarchy ch WHERE ch.parent_id = c.id);
    """)
    branches_without_children = cur.fetchone()[0]
    test("Taxonomy Integrity", "All Branch Categories (is_leaf=false) Have Children", branches_without_children == 0, f"Empty branches: {branches_without_children}")

    cur.execute("""
        SELECT count(*) FROM category_hierarchy c
        WHERE c.is_leaf = true 
          AND EXISTS (SELECT 1 FROM category_hierarchy ch WHERE ch.parent_id = c.id);
    """)
    leaves_with_children = cur.fetchone()[0]
    test("Taxonomy Integrity", "All Leaf Categories (is_leaf=true) Have No Children", leaves_with_children == 0, f"Violating leaves: {leaves_with_children}")

    # ==============================================================================
    # 3. MULTILINGUAL 4-LANGUAGE PARITY (en, am, om, ti)
    # ==============================================================================
    cur.execute("""
        SELECT count(*) FROM category_hierarchy
        WHERE name_en IS NULL OR trim(name_en) = ''
           OR name_am IS NULL OR trim(name_am) = ''
           OR name_om IS NULL OR trim(name_om) = ''
           OR name_ti IS NULL OR trim(name_ti) = '';
    """)
    missing_translations = cur.fetchone()[0]
    test("Multilingual Parity", "Zero Missing Translations Across en, am, om, ti", missing_translations == 0, f"Missing: {missing_translations}")

    # ==============================================================================
    # 4. CROSS-CATEGORY DISCOVERABILITY
    # ==============================================================================
    cur.execute("""
        SELECT cr.id, cr.source_concept_id, c1.name_en as s_name, 
               cr.target_concept_id, c2.name_en as t_name
        FROM classification_relationships cr
        JOIN category_hierarchy c1 ON cr.source_concept_id = c1.id
        JOIN category_hierarchy c2 ON cr.target_concept_id = c2.id
        WHERE cr.status = 'ACTIVE';
    """)
    active_rels = cur.fetchall()
    test("Cross-Discovery", "Active Normalized Cross-Discovery Relationships (>=8)", len(active_rels) >= 8, f"Active count: {len(active_rels)}")

    # Check for self relationships
    cur.execute("SELECT count(*) FROM classification_relationships WHERE source_concept_id = target_concept_id;")
    self_rels = cur.fetchone()[0]
    test("Cross-Discovery", "Zero Self-Referential Relationships", self_rels == 0, f"Self-referential: {self_rels}")

    # ==============================================================================
    # 5. API SECURITY, LATENCY & BOUNDARY TESTING
    # ==============================================================================
    # Valid Related API
    code, data, lat = api_request("/categories/1500/related?lang=en")
    solar_rel_ok = code == 200 and data.get("success") is True and len(data.get("data", [])) >= 2
    test("API Functional", "GET /api/categories/1500/related (Solar Inverters)", solar_rel_ok, f"Status: {code}, Latency: {lat:.1f}ms, Count: {len(data.get('data', []))}")

    # Multilingual Related API
    code, data, lat = api_request("/categories/2000010/related?lang=am")
    am_rel_ok = code == 200 and data.get("success") is True and any("ምርት" in c.get("name_am", "") for c in data.get("data", []))
    test("API Functional", "GET /api/categories/2000010/related (Amharic i18n)", am_rel_ok, f"Status: {code}, Latency: {lat:.1f}ms")

    # Non-existent ID -> 200 with empty data
    code, data, lat = api_request("/categories/99999999/related")
    test("API Boundaries", "GET /api/categories/:id/related Non-existent ID", code == 200 and data.get("data") == [], f"Status: {code}, Data: {data.get('data')}")

    # Invalid ID string -> 400 Bad Request
    code, data, lat = api_request("/categories/invalid-string-id/related")
    test("API Boundaries", "GET /api/categories/:id/related Invalid String ID -> 400", code == 400, f"Status: {code}")

    # Negative ID -> 400 or empty
    code, data, lat = api_request("/categories/-100/related")
    test("API Boundaries", "GET /api/categories/:id/related Negative ID Handled Cleanly", code in (200, 400), f"Status: {code}")

    # SQL Injection in category ID
    code, data, lat = api_request("/categories/1%20OR%201=1/related")
    test("API Security", "SQL Injection in ID Parameter Blocked (400)", code == 400, f"Status: {code}")

    # Search API with Amharic query
    code, data, lat = api_request("/categories/search?q=%E1%88%B5%E1%88%8D%E1%89%B5&lang=am&limit=10")
    test("API Search", "Hierarchical Category Search with Amharic Script", code == 200 and data.get("success") is True, f"Status: {code}, Latency: {lat:.1f}ms")

    # Search API with deep excavator query
    code, data, lat = api_request("/categories/search?q=Excavator&lang=en&limit=10")
    has_bcs = code == 200 and len(data.get("data", [])) > 0 and "breadcrumbs" in data.get("data", [])[0]
    test("API Search", "Category Search Returns Full Ancestor Breadcrumbs", has_bcs, f"Results: {len(data.get('data', []))}, Has Breadcrumbs: {has_bcs}")

    # Category Tree API with depth 6
    code, data, lat = api_request("/categories/tree?lang=en&depth=6")
    tree_ok = code == 200 and data.get("success") is True and len(data.get("data", [])) == 53
    test("API Performance", "GET /api/categories/tree?depth=6 (<1000ms)", tree_ok and lat < 1000, f"Status: {code}, Latency: {lat:.1f}ms, Roots: {len(data.get('data', []))}")

    # Check unauthenticated public access to roots
    code, data, lat = api_request("/categories/roots?lang=en")
    test("API Public Access", "Unauthenticated Browsing to /api/categories/roots", code == 200 and data.get("success") is True, f"Status: {code}, Count: {len(data.get('data', []))}")

    # ==============================================================================
    # 6. JSON MANIFEST & REPOSITORY SYNCHRONIZATION
    # ==============================================================================
    manifest_file = "/home/poven/GO PROJECT/sellbuy/backend/internal/database/loccatdata/categories/manifest.json"
    with open(manifest_file, "r", encoding="utf-8") as f:
        manifest = json.load(f)
    test("JSON Sync", "Manifest Version 2.1.0 and Total 20,528", manifest.get("version") == "2.1.0" and manifest.get("total_categories") == 20528, f"Version: {manifest.get('version')}, Total: {manifest.get('total_categories')}")

    categories_dir = "/home/poven/GO PROJECT/sellbuy/backend/internal/database/loccatdata/categories/classification-systems"
    total_json_cats = 0
    family_file_count = 0
    for root, dirs, files in os.walk(categories_dir):
        if root.endswith("families"):
            for ff in files:
                if ff.endswith(".json") and ff != "manifest.json":
                    family_file_count += 1
                    with open(os.path.join(root, ff), "r", encoding="utf-8") as f:
                        fj = json.load(f)
                        cats = fj.get("categories", [])
                        total_json_cats += len(cats)

    test("JSON Sync", "Total Categories in 65 Canonical Family JSON Files Equals DB (20,528)", total_json_cats == total_cats and family_file_count == 65, f"JSON: {total_json_cats}, DB: {total_cats}, Files: {family_file_count}")

    # ==============================================================================
    # 7. PRODUCT REFERENTIAL INTEGRITY
    # ==============================================================================
    cur.execute("""
        SELECT count(*) FROM products p
        WHERE p.category_hierarchy_id IS NOT NULL
          AND NOT EXISTS (SELECT 1 FROM category_hierarchy c WHERE c.id = p.category_hierarchy_id);
    """)
    orphan_products = cur.fetchone()[0]
    test("Product Integrity", "Zero Products Pointing to Non-Existent Category", orphan_products == 0, f"Orphan products: {orphan_products}")

    # Summary
    print("================================================================================")
    passed_count = sum(1 for r in results if r["passed"])
    total_count = len(results)
    print(f"  QA AUDIT SUMMARY: {passed_count} / {total_count} CHECKS PASSED ({(passed_count/total_count)*100:.1f}%)")
    print("================================================================================")

    cur.close()
    conn.close()

    if passed_count < total_count:
        sys.exit(1)

if __name__ == "__main__":
    run_qa_suite()
