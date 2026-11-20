#!/usr/bin/env bash
# ============================================================================
# HUNTIQ Database Migration Runner (Native Go Migration Tool)
# ============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

exec go run "${BACKEND_DIR}/internal/database/migrate" "$@"
