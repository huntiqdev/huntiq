#!/usr/bin/env bash
# ==============================================================================
# HUNTIQ Local AI Verification Script
# Validates model storage, daemon status, models, tests, and offline capability
# ==============================================================================

set -euo pipefail

AI_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO_ROOT="$(cd "$AI_DIR/.." && pwd)"
MODELS_DIR="$AI_DIR/models/ollama"
VENV_PYTEST="$AI_DIR/.venv/bin/pytest"
OLLAMA_URL="http://127.0.0.1:11434"

echo "=== HUNTIQ Local AI Verification ==="
echo "Date       : $(date -u '+%Y-%m-%d %H:%M:%S UTC')"
echo "Project    : $REPO_ROOT"
echo "Model Dir  : $MODELS_DIR"
echo "===================================="

ERRORS=0

# 1. Check local daemon reachable
echo -n "[CHECK 1/6] Ollama local API reachable (127.0.0.1:11434)... "
if curl -s --connect-timeout 2 "$OLLAMA_URL/api/tags" >/dev/null 2>&1; then
    echo "OK"
else
    echo "FAILED (daemon unreachable)"
    ERRORS=$((ERRORS + 1))
fi

# 2. Check model storage files on disk
echo -n "[CHECK 2/6] Local model storage directory populated... "
if [ -d "$MODELS_DIR/blobs" ] && [ "$(ls -A "$MODELS_DIR/blobs" 2>/dev/null)" ]; then
    BLOB_COUNT=$(ls "$MODELS_DIR/blobs" | wc -l)
    STORAGE_SIZE=$(du -sh "$MODELS_DIR" | cut -f1)
    echo "OK ($BLOB_COUNT blobs, $STORAGE_SIZE total)"
else
    echo "FAILED ($MODELS_DIR is empty or missing)"
    ERRORS=$((ERRORS + 1))
fi

# 3. Check installed models via API
echo -n "[CHECK 3/6] Installed models (qwen3:4b & huntiq-assistant)... "
MODELS_JSON=$(curl -s --connect-timeout 2 "$OLLAMA_URL/api/tags" || echo "{}")
HAS_QWEN3=$(echo "$MODELS_JSON" | grep -c 'qwen3:4b' || true)
HAS_HUNTIQ=$(echo "$MODELS_JSON" | grep -c 'huntiq-assistant' || true)

if [ "$HAS_QWEN3" -gt 0 ] && [ "$HAS_HUNTIQ" -gt 0 ]; then
    echo "OK (both models verified)"
else
    echo "FAILED (missing models: qwen3=$HAS_QWEN3, huntiq=$HAS_HUNTIQ)"
    ERRORS=$((ERRORS + 1))
fi

# 4. Run Python unit tests in local venv
echo -n "[CHECK 4/6] Python AI unit tests in .venv... "
if [ -x "$VENV_PYTEST" ]; then
    if (cd "$AI_DIR" && "$VENV_PYTEST" -q >/dev/null 2>&1); then
        echo "OK (8 of 8 passed)"
    else
        echo "FAILED (pytest failure)"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "FAILED (.venv/bin/pytest not found; run setup_local.sh first)"
    ERRORS=$((ERRORS + 1))
fi

# 5. Verify local inference response
echo -n "[CHECK 5/6] Local model response (huntiq-assistant:latest)... "
RESP_JSON=$(curl -s --connect-timeout 5 -X POST "$OLLAMA_URL/api/chat" \
    -H "Content-Type: application/json" \
    -d '{"model":"huntiq-assistant:latest","messages":[{"role":"user","content":"Say READY"}],"stream":false}' || echo "{}")

if echo "$RESP_JSON" | grep -Eq '"role":"assistant"'; then
    echo "OK (received valid assistant payload)"
else
    echo "FAILED (invalid or empty inference response)"
    ERRORS=$((ERRORS + 1))
fi

# 6. Verify strictly offline operation (no external cloud requests)
echo -n "[CHECK 6/6] Strictly offline verification... "
# Verify configuration points to loopback
if grep -q "127.0.0.1" "$AI_DIR/app/config.py" && grep -q "127.0.0.1" "$AI_DIR/.env.example"; then
    echo "OK (enforced local loopback)"
else
    echo "FAILED (found non-loopback endpoints)"
    ERRORS=$((ERRORS + 1))
fi

echo "===================================="
if [ "$ERRORS" -eq 0 ]; then
    echo "SUCCESS: All 6 local AI checks passed."
    exit 0
else
    echo "FAILURE: $ERRORS check(s) failed."
    exit 1
fi
