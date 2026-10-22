#!/usr/bin/env bash
# ============================================================================
# HUNTIQ Database Migration Runner
# ============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
MIGRATIONS_DIR="${BACKEND_DIR}/internal/database/migrations"

# Load .env if present
if [ -f "${BACKEND_DIR}/.env" ]; then
    # shellcheck disable=SC2046
    export $(grep -v '^#' "${BACKEND_DIR}/.env" | grep -v '^\s*$' | xargs)
fi

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-hulubedej_db}"
DB_PASSWORD="${DB_PASSWORD:-}"

if [ -n "${DB_PASSWORD}" ]; then
    export PGPASSWORD="${DB_PASSWORD}"
fi

PSQL_ARGS=("-h" "${DB_HOST}" "-p" "${DB_PORT}" "-U" "${DB_USER}" "-d" "${DB_NAME}")

echo "==> Connecting to PostgreSQL at ${DB_HOST}:${DB_PORT}/${DB_NAME} as ${DB_USER}..."
psql "${PSQL_ARGS[@]}" -c "SELECT 1;" > /dev/null

echo "==> Ensuring schema_migrations tracking table exists..."
psql "${PSQL_ARGS[@]}" -c "
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
" > /dev/null

APPLIED_COUNT=0
SKIPPED_COUNT=0

while IFS= read -r -d '' file; do
    version="$(basename "${file}")"
    is_applied=$(psql "${PSQL_ARGS[@]}" -t -A -c "SELECT COUNT(*) FROM schema_migrations WHERE version='${version}';")

    if [ "${is_applied}" -eq "0" ]; then
        echo "  [MIGRATING] ${version}..."
        psql "${PSQL_ARGS[@]}" -v ON_ERROR_STOP=1 -f "${file}" > /dev/null
        psql "${PSQL_ARGS[@]}" -c "INSERT INTO schema_migrations (version) VALUES ('${version}');" > /dev/null
        echo "  [DONE]      ${version}"
        ((APPLIED_COUNT++)) || true
    else
        echo "  [SKIPPED]   ${version} (already applied)"
        ((SKIPPED_COUNT++)) || true
    fi
done < <(find "${MIGRATIONS_DIR}" -maxdepth 1 -name "*.sql" -print0 | sort -z)

echo "==> Migration complete! Applied: ${APPLIED_COUNT}, Skipped: ${SKIPPED_COUNT}."
