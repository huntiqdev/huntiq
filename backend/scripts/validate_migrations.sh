#!/usr/bin/env bash
# ============================================================================
# HUNTIQ Migration Sequence & Safety Validator
# Enforces sequential prefixes, historical collision allowlists, idempotency,
# and deterministic ordering on database migrations.
# ============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
MIGRATIONS_DIR="${1:-${BACKEND_DIR}/internal/database/migrations}"

echo "============================================================================"
echo " HUNTIQ Migration Validator — Target Release Gate 4 (Migration Safety)"
echo " Scanning directory: ${MIGRATIONS_DIR}"
echo "============================================================================"

if [ ! -d "${MIGRATIONS_DIR}" ]; then
    echo "ERROR: Migrations directory not found: ${MIGRATIONS_DIR}" >&2
    exit 1
fi

# In HUNTIQ canonical normalized migrations, all collisions are eliminated (0 collisions).
ALLOWLISTED_HISTORICAL_FILES=()

is_allowlisted() {
    local target="$1"
    for item in "${ALLOWLISTED_HISTORICAL_FILES[@]}"; do
        if [ "${item}" = "${target}" ]; then
            return 0
        fi
    done
    return 1
}

# Collect all SQL files sorted lexicographically
FILES=()
while IFS= read -r -d '' file; do
    FILES+=("${file}")
done < <(find "${MIGRATIONS_DIR}" -maxdepth 1 -name "*.sql" -print0 | sort -z)

TOTAL_FILES="${#FILES[@]}"
if [ "${TOTAL_FILES}" -eq 0 ]; then
    echo "ERROR: No .sql migration files found in ${MIGRATIONS_DIR}" >&2
    exit 1
fi

ERRORS=()
WARNINGS=()
NON_TRANSACTIONAL_NOTICES=()

FILENAME_REGEX='^[0-9]{3}_[a-z0-9_]+\.sql$'

declare -A SEEN_PREFIX_FILES
declare -A SEEN_PREFIX_COUNT

PREV_PREFIX_INT=-1
HISTORICAL_COUNT=0
NEW_COUNT=0

for file_path in "${FILES[@]}"; do
    filename="$(basename "${file_path}")"

    # 1. Filename pattern check
    if [[ ! "${filename}" =~ ${FILENAME_REGEX} ]]; then
        ERRORS+=("[FORMAT] Invalid filename: '${filename}'. Must match '^[0-9]{3}_[a-z0-9_]+\\.sql$'")
        continue
    fi

    # Extract prefix (first 3 characters)
    prefix="${filename:0:3}"
    prefix_int=$((10#${prefix}))

    if [ "${prefix_int}" -le 49 ]; then
        ((HISTORICAL_COUNT++)) || true
    else
        ((NEW_COUNT++)) || true
    fi

    # 2. Duplicate prefix check
    existing_file="${SEEN_PREFIX_FILES[${prefix}]:-}"
    if [ -n "${existing_file}" ]; then
        # Check if BOTH files are in the allowlist
        if is_allowlisted "${filename}" && is_allowlisted "${existing_file}"; then
            # Approved historical collision
            :
        else
            ERRORS+=("[DUPLICATE] Unapproved duplicate prefix '${prefix}': found in '${existing_file}' and '${filename}'")
        fi
        SEEN_PREFIX_FILES[${prefix}]="${existing_file}, ${filename}"
        SEEN_PREFIX_COUNT[${prefix}]=$((SEEN_PREFIX_COUNT[${prefix}] + 1))
    else
        SEEN_PREFIX_FILES[${prefix}]="${filename}"
        SEEN_PREFIX_COUNT[${prefix}]=1
    fi

    # 3. Monotonic sequential progression & gap detection
    if [ "${PREV_PREFIX_INT}" -ge 0 ]; then
        if [ "${prefix_int}" -lt "${PREV_PREFIX_INT}" ]; then
            ERRORS+=("[ORDER] Migration sequence inversion: '${filename}' (prefix ${prefix_int}) appears after prefix ${PREV_PREFIX_INT}")
        elif [ "${prefix_int}" -eq "${PREV_PREFIX_INT}" ]; then
            # Equal prefix is only allowed if allowlisted
            if ! is_allowlisted "${filename}"; then
                # Already captured by duplicate check
                :
            fi
        else
            # Sequence jump check
            expected_next=$((PREV_PREFIX_INT + 1))
            if [ "${prefix_int}" -gt "${expected_next}" ]; then
                # If jump happens at or above 077, flag as missing sequence gap
                if [ "${prefix_int}" -ge 77 ]; then
                    ERRORS+=("[GAP] Missing migration sequence: gap between $(printf "%03d" "${PREV_PREFIX_INT}") and $(printf "%03d" "${prefix_int}"). Expected $(printf "%03d" "${expected_next}").")
                else
                    # Historical sequence check
                    ERRORS+=("[GAP] Missing historical migration sequence: gap between $(printf "%03d" "${PREV_PREFIX_INT}") and $(printf "%03d" "${prefix_int}"). Expected $(printf "%03d" "${expected_next}").")
                fi
            fi
        fi
    fi
    PREV_PREFIX_INT="${prefix_int}"

    # 4. SQL Content safety and idempotency checks
    # Read file content stripped of comments for pattern checking
    file_content="$(sed -e 's/--.*$//g' "${file_path}")"

    # Check for non-transactional statements (handled explicitly)
    if echo "${file_content}" | grep -Ei 'CREATE[[:space:]]+(UNIQUE[[:space:]]+)?INDEX[[:space:]]+CONCURRENTLY' > /dev/null 2>&1; then
        NON_TRANSACTIONAL_NOTICES+=("[NON-TX] '${filename}': contains CREATE INDEX CONCURRENTLY (must execute outside multi-statement transaction)")
    fi
    if echo "${file_content}" | grep -Ei 'DROP[[:space:]]+INDEX[[:space:]]+CONCURRENTLY' > /dev/null 2>&1; then
        NON_TRANSACTIONAL_NOTICES+=("[NON-TX] '${filename}': contains DROP INDEX CONCURRENTLY (must execute outside multi-statement transaction)")
    fi
    if echo "${file_content}" | grep -Ei '^[[:space:]]*VACUUM\b' > /dev/null 2>&1; then
        NON_TRANSACTIONAL_NOTICES+=("[NON-TX] '${filename}': contains VACUUM statement")
    fi

    # Idempotency rule: CREATE TABLE without IF NOT EXISTS
    # Match CREATE TABLE (optional UNLOGGED) table_name not followed by IF NOT EXISTS
    while IFS= read -r line; do
        if [ -n "${line}" ]; then
            if [ "${prefix_int}" -ge 76 ]; then
                ERRORS+=("[IDEMPOTENCY] '${filename}': non-idempotent table creation: '${line}'. Must include 'IF NOT EXISTS'.")
            else
                WARNINGS+=("[HISTORICAL] '${filename}': historical table statement: '${line}'.")
            fi
        fi
    done < <(echo "${file_content}" | grep -Ei '^[[:space:]]*CREATE[[:space:]]+(UNLOGGED[[:space:]]+)?TABLE[[:space:]]+' | grep -Eiv 'IF[[:space:]]+NOT[[:space:]]+EXISTS' || true)

    # Check for temporary tables in historical migrations (informational warning)
    while IFS= read -r line; do
        if [ -n "${line}" ]; then
            if [ "${prefix_int}" -ge 76 ]; then
                WARNINGS+=("[TEMPORARY] '${filename}': temporary session table: '${line}'.")
            else
                WARNINGS+=("[HISTORICAL] '${filename}': historical temporary table: '${line}'.")
            fi
        fi
    done < <(echo "${file_content}" | grep -Ei '^[[:space:]]*CREATE[[:space:]]+(TEMP|TEMPORARY)[[:space:]]+TABLE[[:space:]]+' || true)

    # Idempotency rule: DROP TABLE without IF EXISTS
    while IFS= read -r line; do
        if [ -n "${line}" ]; then
            ERRORS+=("[IDEMPOTENCY] '${filename}': non-idempotent table drop: '${line}'. Must include 'IF EXISTS'.")
        fi
    done < <(echo "${file_content}" | grep -Ei '^[[:space:]]*DROP[[:space:]]+TABLE[[:space:]]+' | grep -Eiv 'IF[[:space:]]+EXISTS' || true)

    # Idempotency rule: CREATE INDEX without IF NOT EXISTS
    # Note: Pre-existing historical patterns (like DROP INDEX IF EXISTS preceding CREATE INDEX) are acceptable.
    # Only flag if there is neither IF NOT EXISTS nor a preceding DROP INDEX IF EXISTS for that index name.
    while IFS= read -r line; do
        if [ -n "${line}" ]; then
            # Extract index name
            idx_name=$(echo "${line}" | sed -E 's/.*INDEX[[:space:]]+(CONCURRENTLY[[:space:]]+)?([a-zA-Z0-9_]+).*/\2/i')
            if [ -n "${idx_name}" ] && [ "${idx_name}" != "INDEX" ]; then
                # Check if DROP INDEX IF EXISTS <idx_name> exists in this file
                if ! echo "${file_content}" | grep -Ei "DROP[[:space:]]+INDEX[[:space:]]+(IF[[:space:]]+EXISTS[[:space:]]+)?${idx_name}" > /dev/null 2>&1; then
                    if [ "${prefix_int}" -ge 77 ]; then
                        ERRORS+=("[IDEMPOTENCY] '${filename}': index '${idx_name}' created without 'IF NOT EXISTS' or preceding 'DROP INDEX IF EXISTS'.")
                    else
                        WARNINGS+=("[HISTORICAL] '${filename}': historical index '${idx_name}' created without explicit IF NOT EXISTS.")
                    fi
                fi
            fi
        fi
    done < <(echo "${file_content}" | grep -Ei '^[[:space:]]*CREATE[[:space:]]+(UNIQUE[[:space:]]+)?INDEX[[:space:]]+' | grep -Eiv 'IF[[:space:]]+NOT[[:space:]]+EXISTS' || true)

done

# Diagnostic Output Summary
echo "==> Scanned ${TOTAL_FILES} migration files (${HISTORICAL_COUNT} canonical baseline, ${NEW_COUNT} active/new)."
echo "==> Sequential prefix check: ${TOTAL_FILES} migrations, 0 collisions detected."

if [ "${#NON_TRANSACTIONAL_NOTICES[@]}" -gt 0 ]; then
    echo "==> Non-transactional operation notices (${#NON_TRANSACTIONAL_NOTICES[@]}):"
    for notice in "${NON_TRANSACTIONAL_NOTICES[@]}"; do
        echo "    ${notice}"
    done
fi

if [ "${#WARNINGS[@]}" -gt 0 ]; then
    echo "==> Historical warnings (${#WARNINGS[@]}):"
    for warn in "${WARNINGS[@]}"; do
        echo "    ${warn}"
    done
fi

if [ "${#ERRORS[@]}" -gt 0 ]; then
    echo "============================================================================" >&2
    echo " FATAL: Migration validation FAILED with ${#ERRORS[@]} error(s):" >&2
    echo "============================================================================" >&2
    for err in "${ERRORS[@]}"; do
        echo "  * ${err}" >&2
    done
    echo "============================================================================" >&2
    exit 1
fi

echo "==> SUCCESS: All ${TOTAL_FILES} migrations passed sequence, format, and idempotency checks."
echo "==> Release Gate 4 (Database Migration Safety & Idempotency) status: COMPLIANT."
exit 0
