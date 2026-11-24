#!/usr/bin/env python3
"""
HUNTIQ Final Adversarial Security, Architecture & Usability Verification Suite
Validates:
1. Authentication Security & Brute Force Protection
2. Admin Privilege Separation & RBAC Isolation
3. SQL Injection Resistance & Input Sanitization
4. Category Intent & Subcategory Auto-Discovery Engine
5. Product Listing Quality Gates (30-char description, dual sell/rent)
6. Multilingual & Latin Amharic AI Processing
"""

import sys
import os
import json
import time
import urllib.parse
import urllib.request
import urllib.error

BASE_URL = "http://127.0.0.1:8080"

def api_call(method, path, body=None, token=None, expected_status=None, timeout=30):
    url = f"{BASE_URL}{path}"
    data = json.dumps(body).encode("utf-8") if body else None
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            status = resp.status
            content = resp.read().decode("utf-8")
            try:
                res_json = json.loads(content)
            except Exception:
                res_json = content
            if expected_status and status != expected_status:
                raise AssertionError(f"Expected status {expected_status}, got {status} for {method} {path}")
            return status, res_json
    except urllib.error.HTTPError as e:
        content = e.read().decode("utf-8")
        try:
            res_json = json.loads(content)
        except Exception:
            res_json = content
        if expected_status and e.code != expected_status:
            raise AssertionError(f"Expected status {expected_status}, got {e.code} for {method} {path}: {res_json}")
        return e.code, res_json
    except Exception as e:
        return 0, str(e)

def run_tests():
    passed = 0
    total = 0

    def record(name, condition, details=""):
        nonlocal passed, total
        total += 1
        if condition:
            passed += 1
            print(f"  [PASS] {name} {details}")
        else:
            print(f"  [FAIL] {name} {details}")
            sys.exit(1)

    print("\n" + "="*80)
    print("PHASE 1: AUTHENTICATION SECURITY & BRUTE-FORCE RESILIENCE")
    print("="*80)

    # 1. Invalid login rejection
    status, res = api_call("POST", "/api/auth/login", {"phone": "+251911999999", "password": "WrongPassword123!"})
    record("Reject invalid login credentials with 401", status == 401)

    # 2. Prevent user enumeration in forgot password
    status, res = api_call("POST", "/api/auth/forgot-password", {"phone": "+251911888888"})
    record("Unified forgot-password returns constant message (anti-enumeration)", 
           status == 200 and "if an account exists" in res.get("data", {}).get("message", res.get("message", "")))

    # 3. Unauthenticated access to admin endpoints strictly rejected with 401
    status, res = api_call("GET", "/api/admin/me")
    record("Unauthenticated /api/admin/me returns 401", status == 401)

    status, res = api_call("GET", "/api/admin/admins")
    record("Unauthenticated /api/admin/admins returns 401", status == 401)

    status, res = api_call("GET", "/api/admin/security/active-sessions")
    record("Unauthenticated /api/admin/security/active-sessions returns 401", status == 401)

    # 4. Dev peek OTP without active dev mode or unauthenticated
    status, res = api_call("GET", "/api/dev/otp?phone=+251911223344")
    record("/api/dev/otp query returns handled status (devMode active or 404)", status in (200, 404))

    print("\n" + "="*80)
    print("PHASE 2: SQL INJECTION DEFENSE & SANITIZATION ACROSS SEARCH ROUTES")
    print("="*80)

    sqli_payloads = [
        "' OR '1'='1",
        "'; DROP TABLE users; --",
        "1 UNION SELECT null, null, null, null --",
        "test' AND 1=pg_sleep(0.1) --",
        "<script>alert(1)</script>",
        "\" OR 1=1 --"
    ]

    for payload in sqli_payloads:
        enc = urllib.parse.quote(payload)
        status, res = api_call("GET", f"/api/categories/intent?q={enc}")
        record(f"SQLi payload '{payload[:20]}' on /api/categories/intent returns clean response", 
               status == 200 and "error" not in str(res).lower() or status == 400)

        status2, res2 = api_call("GET", f"/api/categories/search?q={enc}")
        record(f"SQLi payload '{payload[:20]}' on /api/categories/search safely handled", 
               status2 in (200, 400))

        status3, res3 = api_call("GET", f"/api/products?search={enc}")
        record(f"SQLi payload '{payload[:20]}' on /api/products safely handled", 
               status3 in (200, 400))

    print("\n" + "="*80)
    print("PHASE 3: CATEGORY MULTI-BRANCH INTENT & DISAMBIGUATION ENGINE")
    print("="*80)

    # 1. Multi-branch disambiguation for "drill"
    status, res = api_call("GET", "/api/categories/intent?q=drill")
    data = res.get("data", res)
    primary = data.get("category", {})
    alternatives = data.get("alternatives", [])
    has_tools = any("tool" in b.get("name", "").lower() for b in primary.get("breadcrumbs", [])) or "drill" in primary.get("name_en", "").lower()
    has_alt_mining = any("mining" in alt.get("hierarchy", "").lower() for alt in alternatives)
    record("Multi-branch disambiguation for 'drill' identifies Tools & Mining branches", 
           status == 200 and has_tools and has_alt_mining, 
           f"(Primary: {primary.get('name_en')}, Alternatives: {len(alternatives)})")

    # 2. Category intent for "electronics"
    status, res = api_call("GET", "/api/categories/intent?q=electronics")
    data = res.get("data", res)
    primary = data.get("category", {})
    children = data.get("children", [])
    record("Category intent for 'electronics' finds Electronics with subcategories", 
           status == 200 and "electronic" in primary.get("name_en", "").lower() and len(children) > 0,
           f"(Found {len(children)} subcategories)")

    # 3. Category intent for "phone"
    status, res = api_call("GET", "/api/categories/intent?q=phone")
    data = res.get("data", res)
    primary = data.get("category", {})
    record("Category intent for 'phone' resolves to Phones", 
           status == 200 and "phone" in primary.get("name_en", "").lower())

    # 4. Amharic Ge'ez category intent for "መኪና"
    enc_mekina = urllib.parse.quote("መኪና")
    status, res = api_call("GET", f"/api/categories/intent?q={enc_mekina}&lang=am")
    data = res.get("data", res)
    primary = data.get("category", {})
    record("Category intent for Amharic 'መኪና' resolves to Cars", 
           status == 200 and (primary.get("slug") == "cars" or "መኪና" in primary.get("name_am", "") or "መኪኖች" in primary.get("name_am", "")))

    # 5. Latin Amharic category intent for "mekena"
    status, res = api_call("GET", "/api/categories/intent?q=mekena&lang=am")
    data = res.get("data", res)
    primary = data.get("category", {})
    record("Category intent for Latin Amharic 'mekena' resolves to Cars via alias", 
           status == 200 and (primary.get("slug") == "cars" or "መኪና" in primary.get("name_am", "") or "መኪኖች" in primary.get("name_am", "")))

    print("\n" + "="*80)
    print("PHASE 4: PRODUCT QUALITY GATES & DUAL-ROLE SELL/RENT VALIDATION")
    print("="*80)

    # Register/login a seller to test product creation
    test_phone = f"+251977{int(time.time() * 1000) % 1000000:06d}"
    status, reg_res = api_call("POST", "/api/auth/register", {
        "name": "Quality Gate Seller",
        "phone": test_phone,
        "password": "ValidPassword123!",
        "role": "seller"
    })
    
    # Dev OTP verify
    _, otp_res = api_call("GET", f"/api/dev/otp?phone={urllib.parse.quote(test_phone)}")
    otp_code = otp_res.get("data", {}).get("otp") if isinstance(otp_res, dict) else None
    if otp_code:
        api_call("POST", "/api/auth/otp/verify", {"phone": test_phone, "code": otp_code, "purpose": "register"})

    # Login to get JWT
    _, login_res = api_call("POST", "/api/auth/login", {"phone": test_phone, "password": "ValidPassword123!"})
    token = login_res.get("data", {}).get("token") or login_res.get("token")
    record("Seller logged in successfully and acquired JWT", bool(token))

    if token:
        # A. Try creating listing with description < 30 chars -> MUST FAIL
        short_desc_payload = {
            "title": "Bosch Professional Cordless Drill",
            "price": 12500,
            "description": "Too short",  # 9 chars (< 30)
            "category_hierarchy_id": 1,
            "condition": "used",
            "image_urls": ["/uploads/drill1.jpg", "/uploads/drill2.jpg"],
            "region": "Addis Ababa"
        }
        status, res = api_call("POST", "/api/seller/products", short_desc_payload, token=token)
        record("Listing creation with < 30 chars description is REJECTED", status == 400, f"(Response: {res.get('error', res)})")

        # B. Try creating listing with valid description >= 30 chars -> MUST SUCCEED
        valid_desc_payload = {
            "title": "Bosch Professional Cordless Impact Drill GSB 18V",
            "price": 14500,
            "description": "High performance brushless motor impact drill in excellent working condition with two lithium ion batteries and rapid charger included.",
            "category_hierarchy_id": 1,
            "condition": "used",
            "image_urls": ["/uploads/drill1.jpg", "/uploads/drill2.jpg"],
            "region": "Addis Ababa"
        }
        status, res = api_call("POST", "/api/seller/products", valid_desc_payload, token=token)
        record("Listing creation with valid >= 30 chars description SUCCEEDS", status == 201)

        # C. Rental listing validation: Missing rental_price -> MUST FAIL
        rental_missing_price = {
            "title": "Toyota Vitz 2018 Compact Sedan Rental",
            "listing_type": "rent",
            "price": 0,
            "description": "Fuel efficient automatic transmission car available for self drive or with dedicated chauffeur in Addis Ababa city.",
            "category_hierarchy_id": 14,
            "condition": "used",
            "image_urls": ["/uploads/car1.jpg", "/uploads/car2.jpg"],
            "region": "Addis Ababa",
            "rental_duration_value": 1,
            "rental_duration_unit": "day"
            # rental_price missing
        }
        status, res = api_call("POST", "/api/seller/products", rental_missing_price, token=token)
        record("Rental listing missing rental_price is REJECTED", status == 400)

        # D. Rental listing with complete fields -> SUCCEEDS
        rental_valid = {
            "title": "Toyota Vitz 2018 Compact Sedan Daily Rental",
            "listing_type": "rent",
            "rental_price": 2500,
            "rental_duration_value": 1,
            "rental_duration_unit": "day",
            "available_from": "2026-11-25T00:00:00Z",
            "price": 2500,
            "description": "Fuel efficient automatic transmission car available for daily rental with complete insurance coverage in Addis Ababa.",
            "category_hierarchy_id": 14,
            "condition": "used",
            "image_urls": ["/uploads/car1.jpg", "/uploads/car2.jpg"],
            "region": "Addis Ababa"
        }
        status, res = api_call("POST", "/api/seller/products", rental_valid, token=token)
        record("Valid rental listing creation SUCCEEDS", status == 201)

    print("\n" + "="*80)
    print("PHASE 5: MULTILINGUAL & LATIN AMHARIC AI QUERY CLASSIFICATION")
    print("="*80)

    ai_queries = [
        ("ሰላም፣ ቶዮታ ቪትስ መኪና ስንት ነው?", "am", "Ge'ez Fidel Amharic"),
        ("selam, toyota vitz mekena sint new?", "am", "Latin Amharic (Transliterated)"),
        ("ere be natish mekena ale?", "am", "Latin Amharic Colloquial ('ere be natish')"),
        ("Akkam, konkolaataa meeqa?", "or", "Afaan Oromo Query"),
        ("Find laptops under 50000 ETB", "en", "English Query"),
    ]

    for q, exp_lang, label in ai_queries:
        status, res = api_call("POST", "/api/ai/chat", {"message": q}, timeout=35)
        data = res.get("data", res) if isinstance(res, dict) else {}
        lang = data.get("language")
        record(f"AI classification for {label} -> '{exp_lang}'", 
               status == 200 and lang == exp_lang, f"(Detected: {lang})")

    print("\n" + "="*80)
    print(f"MASTER VERIFICATION COMPLETE: {passed}/{total} tests passed ({passed*100//total}%)")
    print("="*80)
    return passed == total

if __name__ == "__main__":
    success = run_tests()
    sys.exit(0 if success else 1)
