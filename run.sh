#!/usr/bin/env bash
# ==============================================================================
# HUNTIQ — Single Command Unified Development Runner
# Runs Go Backend (:8080) + Next.js Frontend (:3000)
# ==============================================================================

set -euo pipefail

# Project Directories
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

# Terminal Colors
C_RESET='\033[0m'
C_CYAN='\033[1;36m'
C_GREEN='\033[1;32m'
C_YELLOW='\033[1;33m'
C_BLUE='\033[1;34m'
C_RED='\033[1;31m'
C_BOLD='\033[1m'

echo -e "${C_CYAN}╔══════════════════════════════════════════════════════════════╗${C_RESET}"
echo -e "${C_CYAN}║             HUNTIQ Universal Marketplace Platform            ║${C_RESET}"
echo -e "${C_CYAN}║         Single Command Runner: Go Backend + Frontend UI      ║${C_RESET}"
echo -e "${C_CYAN}╚══════════════════════════════════════════════════════════════╝${C_RESET}"
echo ""

# ── 1. Process Cleanup Trap (Graceful Shutdown on Ctrl+C) ──────
PIDS=()
cleanup() {
    echo ""
    echo -e "${C_YELLOW}Shutting down HUNTIQ services gracefully...${C_RESET}"
    for pid in "${PIDS[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid" 2>/dev/null || true
        fi
    done
    wait 2>/dev/null || true
    echo -e "${C_GREEN}All services stopped cleanly.${C_RESET}"
    exit 0
}
trap cleanup SIGINT SIGTERM EXIT

# ── 2. Start Go Backend (:8080) ───────────────────────────────
echo -e "${C_BLUE}[1/2] 🔧 Starting Go Backend API (:8080)...${C_RESET}"

# Check if port 8080 is already occupied
if ss -tulpn 2>/dev/null | grep -q ':8080 '; then
    echo -e "      ${C_YELLOW}Port 8080 is already active (using running backend).${C_RESET}"
else
    (
        cd "$BACKEND_DIR"
        go run ./cmd/api
    ) 2>&1 | sed -e "s/^/[BACKEND] /" &
    PIDS+=($!)
    
    # Wait for backend healthcheck
    echo -n "      Waiting for backend to become ready..."
    for i in {1..20}; do
        if curl -s http://127.0.0.1:8080/health >/dev/null 2>&1 || curl -s http://127.0.0.1:8080/api/health >/dev/null 2>&1; then
            echo -e " ${C_GREEN}Ready!${C_RESET}"
            break
        fi
        sleep 0.5
    done
fi

# ── 3. Start Next.js Frontend (:3000) ─────────────────────────
echo ""
echo -e "${C_BLUE}[2/2] 🌐 Starting Next.js Frontend (:3000)...${C_RESET}"

# Check if port 3000 is already occupied
if ss -tulpn 2>/dev/null | grep -q ':3000 '; then
    echo -e "      ${C_YELLOW}Port 3000 is already active (using running frontend).${C_RESET}"
else
    (
        cd "$FRONTEND_DIR"
        npm run dev
    ) 2>&1 | sed -e "s/^/[FRONTEND] /" &
    PIDS+=($!)
    
    # Wait briefly for Next.js initialization
    sleep 2
fi

echo ""
echo -e "${C_GREEN}════════════════════════════════════════════════════════════════${C_RESET}"
echo -e "${C_GREEN}${C_BOLD}✓ HUNTIQ is now running locally!${C_RESET}"
echo ""
echo -e "  🌐 ${C_BOLD}Frontend Application :${C_RESET} ${C_CYAN}http://localhost:3000${C_RESET}"
echo -e "  🔧 ${C_BOLD}Go Backend API       :${C_RESET} ${C_CYAN}http://localhost:8080${C_RESET}"
echo -e "${C_GREEN}════════════════════════════════════════════════════════════════${C_RESET}"
echo -e "${C_YELLOW}Press Ctrl+C to stop all services.${C_RESET}"
echo ""

# Keep running until Ctrl+C
wait
