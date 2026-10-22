#!/bin/bash

# ============================================================
# Migration Testing Script - Optimized 15 Migrations
# Tests that all 15 consolidated migrations run without errors
# ============================================================

set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}Testing Optimized Migrations (15 files)${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""

# Database configuration
DB_NAME="${POSTGRES_DB:-marketplace_test}"
DB_USER="${POSTGRES_USER:-postgres}"
DB_PASSWORD="${POSTGRES_PASSWORD:-}"
DB_HOST="${POSTGRES_HOST:-localhost}"
DB_PORT="${POSTGRES_PORT:-5432}"

MIGRATIONS_DIR="backend/internal/database/migrations"

# Export password for psql
if [ -n "$DB_PASSWORD" ]; then
    export PGPASSWORD="$DB_PASSWORD"
fi

PSQL_CMD="psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME"

echo -e "${YELLOW}Step 1: Check database connection${NC}"
if $PSQL_CMD -c "SELECT version();" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Connected to database${NC}"
else
    echo -e "${RED}✗ Failed to connect to database${NC}"
    echo -e "${YELLOW}Make sure PostgreSQL is running and credentials are correct${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}Step 2: Create fresh test database${NC}"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;" > /dev/null 2>&1 || true
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;" > /dev/null 2>&1
echo -e "${GREEN}✓ Fresh database created${NC}"

echo ""
echo -e "${YELLOW}Step 3: Run all 15 migrations in order${NC}"

MIGRATION_COUNT=0
FAILED_MIGRATIONS=0

for migration_file in $(ls $MIGRATIONS_DIR/*.sql | sort); do
    migration_name=$(basename $migration_file)
    echo -e "${BLUE}Running: $migration_name${NC}"
    
    if $PSQL_CMD -f "$migration_file" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ $migration_name completed${NC}"
        MIGRATION_COUNT=$((MIGRATION_COUNT + 1))
    else
        echo -e "${RED}✗ $migration_name FAILED${NC}"
        FAILED_MIGRATIONS=$((FAILED_MIGRATIONS + 1))
        # Show error details
        $PSQL_CMD -f "$migration_file"
    fi
done

echo ""
echo -e "${YELLOW}Step 4: Verify database structure${NC}"

# Check critical tables exist
TABLES_TO_CHECK=(
    "users"
    "products"
    "categories"
    "locations"
    "transactions"
    "order_requests"
    "notifications"
    "admin_notifications"
    "admin_logs"
    "referral_commissions"
    "referral_withdrawals"
    "upgrade_requests"
    "user_tier_history"
    "platform_settings"
    "seller_reviews"
    "user_warnings"
)

TABLES_FOUND=0
TABLES_MISSING=0

for table in "${TABLES_TO_CHECK[@]}"; do
    if $PSQL_CMD -c "SELECT 1 FROM $table LIMIT 1;" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Table: $table${NC}"
        TABLES_FOUND=$((TABLES_FOUND + 1))
    else
        echo -e "${RED}✗ Table missing: $table${NC}"
        TABLES_MISSING=$((TABLES_MISSING + 1))
    fi
done

echo ""
echo -e "${YELLOW}Step 5: Verify functions and triggers${NC}"

# Check functions exist
FUNCTIONS_TO_CHECK=(
    "cleanup_unverified_users"
    "cleanup_expired_products"
    "update_seller_rating"
    "refresh_all_materialized_views"
)

FUNCTIONS_FOUND=0
FUNCTIONS_MISSING=0

for func in "${FUNCTIONS_TO_CHECK[@]}"; do
    if $PSQL_CMD -c "SELECT 1 FROM pg_proc WHERE proname='$func';" | grep -q 1; then
        echo -e "${GREEN}✓ Function: $func${NC}"
        FUNCTIONS_FOUND=$((FUNCTIONS_FOUND + 1))
    else
        echo -e "${RED}✗ Function missing: $func${NC}"
        FUNCTIONS_MISSING=$((FUNCTIONS_MISSING + 1))
    fi
done

echo ""
echo -e "${YELLOW}Step 6: Verify materialized views${NC}"

# Check materialized views exist
VIEWS_TO_CHECK=(
    "marketplace_listings"
    "seller_rating_stats"
    "user_warning_counts"
)

VIEWS_FOUND=0
VIEWS_MISSING=0

for view in "${VIEWS_TO_CHECK[@]}"; do
    if $PSQL_CMD -c "SELECT 1 FROM $view LIMIT 1;" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ View: $view${NC}"
        VIEWS_FOUND=$((VIEWS_FOUND + 1))
    else
        echo -e "${RED}✗ View missing: $view${NC}"
        VIEWS_MISSING=$((VIEWS_MISSING + 1))
    fi
done

echo ""
echo -e "${YELLOW}Step 7: Test critical columns exist${NC}"

# Check critical columns
COLUMNS_TO_CHECK=(
    "users.current_demand"
    "users.seller_preference"
    "products.listing_type"
    "products.rental_price"
    "products.rental_price_per_month"
    "categories.is_active"
)

COLUMNS_FOUND=0
COLUMNS_MISSING=0

for col in "${COLUMNS_TO_CHECK[@]}"; do
    IFS='.' read -r table column <<< "$col"
    if $PSQL_CMD -c "SELECT $column FROM $table LIMIT 1;" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Column: $col${NC}"
        COLUMNS_FOUND=$((COLUMNS_FOUND + 1))
    else
        echo -e "${RED}✗ Column missing: $col${NC}"
        COLUMNS_MISSING=$((COLUMNS_MISSING + 1))
    fi
done

echo ""
echo -e "${YELLOW}Step 8: Test indexes exist${NC}"

# Count indexes
INDEX_COUNT=$($PSQL_CMD -t -c "SELECT COUNT(*) FROM pg_indexes WHERE schemaname='public';")
echo -e "${GREEN}✓ Total indexes created: $INDEX_COUNT${NC}"

echo ""
echo -e "${YELLOW}Step 9: Test constraints${NC}"

# Count constraints
CONSTRAINT_COUNT=$($PSQL_CMD -t -c "SELECT COUNT(*) FROM information_schema.table_constraints WHERE constraint_schema='public';")
echo -e "${GREEN}✓ Total constraints: $CONSTRAINT_COUNT${NC}"

echo ""
echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}MIGRATION TEST SUMMARY${NC}"
echo -e "${BLUE}================================================${NC}"
echo -e "Migrations run:        ${GREEN}$MIGRATION_COUNT / 15${NC}"
echo -e "Migrations failed:     ${RED}$FAILED_MIGRATIONS${NC}"
echo -e "Tables found:          ${GREEN}$TABLES_FOUND / ${#TABLES_TO_CHECK[@]}${NC}"
echo -e "Tables missing:        ${RED}$TABLES_MISSING${NC}"
echo -e "Functions found:       ${GREEN}$FUNCTIONS_FOUND / ${#FUNCTIONS_TO_CHECK[@]}${NC}"
echo -e "Functions missing:     ${RED}$FUNCTIONS_MISSING${NC}"
echo -e "Views found:           ${GREEN}$VIEWS_FOUND / ${#VIEWS_TO_CHECK[@]}${NC}"
echo -e "Views missing:         ${RED}$VIEWS_MISSING${NC}"
echo -e "Columns found:         ${GREEN}$COLUMNS_FOUND / ${#COLUMNS_TO_CHECK[@]}${NC}"
echo -e "Columns missing:       ${RED}$COLUMNS_MISSING${NC}"
echo -e "Total indexes:         ${GREEN}$INDEX_COUNT${NC}"
echo -e "Total constraints:     ${GREEN}$CONSTRAINT_COUNT${NC}"
echo ""

# Final verdict
TOTAL_ERRORS=$((FAILED_MIGRATIONS + TABLES_MISSING + FUNCTIONS_MISSING + VIEWS_MISSING + COLUMNS_MISSING))

if [ $TOTAL_ERRORS -eq 0 ]; then
    echo -e "${GREEN}================================================${NC}"
    echo -e "${GREEN}✓✓✓ ALL TESTS PASSED! ✓✓✓${NC}"
    echo -e "${GREEN}15 optimized migrations work perfectly!${NC}"
    echo -e "${GREEN}================================================${NC}"
    exit 0
else
    echo -e "${RED}================================================${NC}"
    echo -e "${RED}✗✗✗ TESTS FAILED ✗✗✗${NC}"
    echo -e "${RED}Total errors: $TOTAL_ERRORS${NC}"
    echo -e "${RED}================================================${NC}"
    exit 1
fi
