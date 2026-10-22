#!/bin/bash
# ============================================================
# API Endpoint Testing Script
# Tests all neutral registration endpoints and redirects
# ============================================================

set -e

# Configuration
API_URL="${API_URL:-http://localhost:8080/api}"
TEST_PHONE="0912345678"
TEST_PASSWORD="Test1234!"
TEST_NAME="Test User"
TOKEN=""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
PASS=0
FAIL=0

# Helper functions
print_test() {
    echo -e "\n${YELLOW}━━━ TEST: $1 ━━━${NC}"
}

print_pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    ((PASS++))
}

print_fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    ((FAIL++))
}

test_endpoint() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    local expected_status="$4"
    local description="$5"
    
    if [ -n "$TOKEN" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$API_URL$endpoint" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $TOKEN" \
            -d "$data")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$API_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" = "$expected_status" ]; then
        print_pass "$description (HTTP $status_code)"
        echo "$body"
        return 0
    else
        print_fail "$description (Expected $expected_status, got $status_code)"
        echo "$body"
        return 1
    fi
}

# ============================================================
# START TESTS
# ============================================================

echo "╔════════════════════════════════════════════════════════╗"
echo "║  Neutral Registration API Test Suite                  ║"
echo "╚════════════════════════════════════════════════════════╝"
echo "API URL: $API_URL"

# Test 1: Health Check
print_test "API Health Check"
if curl -s -o /dev/null -w "%{http_code}" "$API_URL/../health" | grep -q "200"; then
    print_pass "API is reachable"
else
    print_fail "API is not reachable"
    echo "Please start the backend server first"
    exit 1
fi

# Test 2: Public Endpoints (No Auth Required)
print_test "Public Product Listing (SALE)"
test_endpoint "GET" "/products?listing_type=sale&limit=5" "" "200" "Get sale listings"

print_test "Public Product Listing (RENT)"
test_endpoint "GET" "/products?listing_type=rent&limit=5" "" "200" "Get rental listings"

print_test "Get Categories"
test_endpoint "GET" "/categories" "" "200" "Get active categories"

print_test "Get Locations"
test_endpoint "GET" "/locations" "" "200" "Get location list"

# Test 3: Registration Flow
print_test "Seller Registration"
RANDOM_PHONE="0${RANDOM}${RANDOM}"
REG_DATA="{\"name\":\"$TEST_NAME\",\"phone\":\"$RANDOM_PHONE\",\"password\":\"$TEST_PASSWORD\"}"

if test_endpoint "POST" "/seller/register" "$REG_DATA" "200" "Register new seller"; then
    echo "Registration successful"
else
    print_fail "Registration failed - using existing account for tests"
    # Try to login with test account instead
    RANDOM_PHONE="$TEST_PHONE"
fi

# Test 4: Login
print_test "Seller Login"
LOGIN_DATA="{\"phone\":\"$RANDOM_PHONE\",\"password\":\"$TEST_PASSWORD\"}"
response=$(curl -s -X POST "$API_URL/seller/login" \
    -H "Content-Type: application/json" \
    -d "$LOGIN_DATA")

if echo "$response" | grep -q "token"; then
    TOKEN=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    print_pass "Login successful"
    echo "Token: ${TOKEN:0:20}..."
else
    print_fail "Login failed"
    echo "$response"
fi

# Test 5: Demand Selection Endpoint
if [ -n "$TOKEN" ]; then
    print_test "Set Demand to FIND"
    DEMAND_DATA='{"demand":"FIND"}'
    response=$(test_endpoint "POST" "/auth/demand" "$DEMAND_DATA" "200" "Set demand to FIND")
    
    if echo "$response" | grep -q "next_path"; then
        print_pass "Demand endpoint returns next_path"
    else
        print_fail "Demand endpoint missing next_path"
    fi
    
    print_test "Set Demand to LIST"
    DEMAND_DATA='{"demand":"LIST"}'
    test_endpoint "POST" "/auth/demand" "$DEMAND_DATA" "200" "Set demand to LIST"
    
    print_test "Invalid Demand Value"
    DEMAND_DATA='{"demand":"INVALID"}'
    test_endpoint "POST" "/auth/demand" "$DEMAND_DATA" "400" "Reject invalid demand"
else
    print_fail "Skipping demand tests - no auth token"
fi

# Test 6: Authenticated Product Endpoints
if [ -n "$TOKEN" ]; then
    print_test "Get Seller Stats"
    test_endpoint "GET" "/seller/stats" "" "200" "Get seller statistics"
    
    print_test "Get Seller Products"
    test_endpoint "GET" "/seller/products" "" "200" "Get seller product list"
    
    print_test "Fee Preview"
    test_endpoint "GET" "/seller/fee-preview?price=5000&duration=30" "" "200" "Get listing fee preview"
fi

# Test 7: Product Creation (SALE)
if [ -n "$TOKEN" ]; then
    print_test "Create Sale Listing"
    SALE_DATA='{
        "title":"Test Product Sale",
        "price":5000,
        "listing_type":"sale",
        "description":"Test description",
        "condition":"used",
        "image_urls":["https://example.com/image1.jpg","https://example.com/image2.jpg"],
        "region":"Addis Ababa",
        "listing_duration_days":30
    }'
    
    if test_endpoint "POST" "/seller/products" "$SALE_DATA" "200" "Create sale listing"; then
        PRODUCT_ID=$(echo "$response" | grep -o '"id":[0-9]*' | cut -d':' -f2 | head -1)
        print_pass "Sale listing created (ID: $PRODUCT_ID)"
    fi
fi

# Test 8: Product Creation (RENT)
if [ -n "$TOKEN" ]; then
    print_test "Create Rental Listing"
    RENT_DATA='{
        "title":"Test Property Rental",
        "listing_type":"rent",
        "rental_price":3000,
        "rental_duration_value":6,
        "rental_duration_unit":"month",
        "deposit_amount":3000,
        "description":"Test rental property",
        "condition":"used",
        "image_urls":["https://example.com/image1.jpg","https://example.com/image2.jpg"],
        "region":"Addis Ababa",
        "listing_duration_days":30
    }'
    
    test_endpoint "POST" "/seller/products" "$RENT_DATA" "200" "Create rental listing"
fi

# Test 9: Buyer Endpoints
print_test "Buyer Registration"
BUYER_PHONE="09${RANDOM}${RANDOM}"
BUYER_DATA="{\"name\":\"Test Buyer\",\"phone\":\"$BUYER_PHONE\",\"password\":\"$TEST_PASSWORD\"}"
test_endpoint "POST" "/buyer/register" "$BUYER_DATA" "200" "Register new buyer"

print_test "Buyer Login"
BUYER_LOGIN="{\"phone\":\"$BUYER_PHONE\",\"password\":\"$TEST_PASSWORD\"}"
response=$(curl -s -X POST "$API_URL/buyer/login" \
    -H "Content-Type: application/json" \
    -d "$BUYER_LOGIN")

if echo "$response" | grep -q "token"; then
    BUYER_TOKEN=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    print_pass "Buyer login successful"
    
    # Test buyer demand selection
    print_test "Buyer Demand Selection"
    TOKEN="$BUYER_TOKEN"
    DEMAND_DATA='{"demand":"FIND"}'
    test_endpoint "POST" "/auth/demand" "$DEMAND_DATA" "200" "Buyer can set demand"
    TOKEN="" # Reset token
fi

# Test 10: Authorization Tests
print_test "Unauthorized Demand Access"
TOKEN_BACKUP="$TOKEN"
TOKEN=""
test_endpoint "POST" "/auth/demand" '{"demand":"FIND"}' "401" "Reject unauthorized demand request"
TOKEN="$TOKEN_BACKUP"

print_test "Unauthorized Product Creation"
TOKEN=""
test_endpoint "POST" "/seller/products" '{"title":"Test"}' "401" "Reject unauthorized product creation"
TOKEN="$TOKEN_BACKUP"

# Test 11: Edge Cases
print_test "Missing Required Fields"
test_endpoint "POST" "/seller/register" '{"name":"Test"}' "400" "Reject incomplete registration"

print_test "Invalid Phone Format"
test_endpoint "POST" "/seller/register" '{"name":"Test","phone":"invalid","password":"Test1234!"}' "400" "Reject invalid phone"

# ============================================================
# SUMMARY
# ============================================================

echo ""
echo "╔════════════════════════════════════════════════════════╗"
echo "║                  TEST SUMMARY                          ║"
echo "╚════════════════════════════════════════════════════════╝"
echo -e "${GREEN}✓ Passed: $PASS${NC}"
echo -e "${RED}✗ Failed: $FAIL${NC}"

TOTAL=$((PASS + FAIL))
if [ $TOTAL -gt 0 ]; then
    SUCCESS_RATE=$((PASS * 100 / TOTAL))
    echo "Success Rate: $SUCCESS_RATE%"
fi

if [ $FAIL -eq 0 ]; then
    echo -e "\n${GREEN}🎉 ALL TESTS PASSED!${NC}"
    exit 0
else
    echo -e "\n${RED}⚠️  SOME TESTS FAILED${NC}"
    exit 1
fi
