#!/bin/bash

# Test script for Category Hierarchy API endpoints
# Usage: ./test_category_hierarchy.sh

BASE_URL="http://localhost:8080"
ADMIN_TOKEN=""  # Set after admin login

echo "=================================="
echo "Category Hierarchy API Test Script"
echo "=================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_RUN=0
TESTS_PASSED=0

run_test() {
    local test_name="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    local expected_status="${5:-200}"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    echo -e "${BLUE}Test $TESTS_RUN: $test_name${NC}"
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -d "$data")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    
    if [ "$http_code" == "$expected_status" ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code)"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        echo -e "${RED}✗ FAIL${NC} (Expected HTTP $expected_status, got $http_code)"
        echo "$body"
    fi
    echo ""
}

# ============================================================================
# PUBLIC ENDPOINTS (No auth required)
# ============================================================================

echo "=== PUBLIC ENDPOINTS ==="
echo ""

run_test "Get category tree (English)" \
    "GET" \
    "/api/categories/tree?lang=en"

run_test "Get category tree (Amharic)" \
    "GET" \
    "/api/categories/tree?lang=am"

run_test "Get category by ID" \
    "GET" \
    "/api/categories/1"

run_test "Get category by slug" \
    "GET" \
    "/api/categories/slug/electronics"

run_test "Get category children" \
    "GET" \
    "/api/categories/1/children"

run_test "Get category breadcrumbs" \
    "GET" \
    "/api/categories/1/ancestors"

run_test "List categories (filtered)" \
    "GET" \
    "/api/categories?level=0&is_active=true&page=1&limit=10"

run_test "Search categories by name" \
    "GET" \
    "/api/categories?search=electronics"

# ============================================================================
# ADMIN ENDPOINTS (Requires auth)
# ============================================================================

echo "=== ADMIN ENDPOINTS (Auth Required) ==="
echo ""

if [ -z "$ADMIN_TOKEN" ]; then
    echo -e "${RED}⚠️  ADMIN_TOKEN not set. Admin endpoint tests skipped.${NC}"
    echo "To test admin endpoints:"
    echo "1. Login as admin and copy JWT token"
    echo "2. Set ADMIN_TOKEN variable: export ADMIN_TOKEN='your_jwt_token'"
    echo "3. Re-run this script"
    echo ""
else
    run_test "Get category stats (Admin)" \
        "GET" \
        "/api/admin/category-hierarchy/1/stats"
    
    run_test "Create new category (Admin)" \
        "POST" \
        "/api/admin/category-hierarchy" \
        '{
            "parent_id": 1,
            "name_en": "Test Category",
            "name_am": "ሙከራ ምድብ",
            "name_om": "Gosa Yaalii",
            "slug": "test-category",
            "is_leaf": true,
            "display_order": 99
        }' \
        201
    
    # Note: Update and delete tests would need the created category ID
    echo -e "${BLUE}Note:${NC} Update/delete tests require manual category ID"
fi

# ============================================================================
# PRODUCT FILTERING WITH CATEGORY HIERARCHY
# ============================================================================

echo "=== PRODUCT FILTERING ==="
echo ""

run_test "Get products by category (flat)" \
    "GET" \
    "/api/products?category_id=1&page=1&limit=5"

run_test "Get products by category hierarchy" \
    "GET" \
    "/api/products?category_hierarchy_id=1&page=1&limit=5"

run_test "Get products with descendants included" \
    "GET" \
    "/api/products?category_hierarchy_id=1&include_descendants=true&page=1&limit=5"

# ============================================================================
# ERROR HANDLING TESTS
# ============================================================================

echo "=== ERROR HANDLING ==="
echo ""

run_test "Get non-existent category" \
    "GET" \
    "/api/categories/99999" \
    "" \
    404

run_test "Invalid category ID format" \
    "GET" \
    "/api/categories/invalid" \
    "" \
    400

# ============================================================================
# SUMMARY
# ============================================================================

echo "=================================="
echo "Test Summary"
echo "=================================="
echo "Tests Run: $TESTS_RUN"
echo "Tests Passed: $TESTS_PASSED"
echo "Tests Failed: $((TESTS_RUN - TESTS_PASSED))"
echo ""

if [ $TESTS_PASSED -eq $TESTS_RUN ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed.${NC}"
    exit 1
fi
