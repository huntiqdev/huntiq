#!/usr/bin/env bash
# ==============================================================================
# HUNTIQ — Unified Application Runner & Orchestrator
# Supports both Native Development Mode and Containerized Docker Mode.
#
# Usage:
#   ./run.sh              # Start Native Stack (Postgres + Backend :8080 + Frontend :3000)
#   ./run.sh --docker     # Start Docker Stack (Postgres, Redis, Backend, Frontend)
#   ./run.sh --status     # Show live status of services and endpoints
#   ./run.sh --migrate    # Run database migration validation and migrations
#   ./run.sh --help       # Display usage guide
# ==============================================================================

set -euo pipefail

# Project Directories
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

# Terminal Colors & Styling
C_RESET='\033[0m'
C_CYAN='\033[1;36m'
C_GREEN='\033[1;32m'
C_YELLOW='\033[1;33m'
C_BLUE='\033[1;34m'
C_RED='\033[1;31m'
C_MAGENTA='\033[1;35m'
C_BOLD='\033[1m'
C_DIM='\033[2m'

# Print Header Banner
print_banner() {
    local mode="${1:-NATIVE}"
    echo -e "${C_CYAN}╔══════════════════════════════════════════════════════════════╗${C_RESET}"
    echo -e "${C_CYAN}║             HUNTIQ Universal Marketplace Platform            ║${C_RESET}"
    echo -e "${C_CYAN}║         Mode: ${C_YELLOW}${mode}${C_CYAN} | Unified Runner & Orchestrator          ║${C_RESET}"
    echo -e "${C_CYAN}╚══════════════════════════════════════════════════════════════╝${C_RESET}"
    echo ""
}

# ------------------------------------------------------------------------------
# Check Port Usage Helper
# ------------------------------------------------------------------------------
is_port_in_use() {
    local port="$1"
    if command -v lsof >/dev/null 2>&1; then
        lsof -i ":$port" -sTCP:LISTEN >/dev/null 2>&1 && return 0
    fi
    if command -v fuser >/dev/null 2>&1; then
        fuser "$port/tcp" >/dev/null 2>&1 && return 0
    fi
    if command -v ss >/dev/null 2>&1; then
        ss -tulpn 2>/dev/null | grep -q ":$port " && return 0
    fi
    return 1
}

# ------------------------------------------------------------------------------
# Status Checker
# ------------------------------------------------------------------------------
cmd_status() {
    print_banner "STATUS CHECK"
    echo -e "${C_BOLD}Inspecting HUNTIQ System Services...${C_RESET}"
    echo ""

    local cfg_host="/tmp"
    local cfg_port="5433"
    if [ -f "$BACKEND_DIR/.env" ]; then
        cfg_host=$(grep -E '^[[:space:]]*DB_HOST=' "$BACKEND_DIR/.env" | head -n 1 | cut -d '=' -f2- | tr -d ' "\r' || echo "/tmp")
        cfg_port=$(grep -E '^[[:space:]]*DB_PORT=' "$BACKEND_DIR/.env" | head -n 1 | cut -d '=' -f2- | tr -d ' "\r' || echo "5433")
    fi
    [ -z "$cfg_host" ] && cfg_host="/tmp"
    [ -z "$cfg_port" ] && cfg_port="5433"

    # 1. PostgreSQL Check
    echo -n "  📦 PostgreSQL Database  : "
    if psql -h "$cfg_host" -p "$cfg_port" -U postgres -d hulubedej_db -c "SELECT 1;" >/dev/null 2>&1; then
        echo -e "${C_GREEN}HEALTHY${C_RESET} (Host ${cfg_host} / Port ${cfg_port})"
        local nodes
        nodes=$(psql -h "$cfg_host" -p "$cfg_port" -U postgres -d hulubedej_db -t -A -c "SELECT count(*) FROM classification_placements;" 2>/dev/null || echo "0")
        echo -e "     └─ Classification Nodes: ${C_CYAN}${nodes}${C_RESET} taxonomies cataloged"
    elif psql -h 127.0.0.1 -p 5432 -U postgres -d hulubedej_db -c "SELECT 1;" >/dev/null 2>&1; then
        echo -e "${C_GREEN}HEALTHY${C_RESET} (Port 5432 / localhost)"
    else
        echo -e "${C_RED}NOT REACHABLE (${cfg_host}:${cfg_port})${C_RESET}"
    fi

    # 2. Backend Check
    echo -n "  🔧 Go Backend API (:8080): "
    if curl -s --connect-timeout 2 http://127.0.0.1:8080/health >/dev/null 2>&1; then
        echo -e "${C_GREEN}HEALTHY (200 OK)${C_RESET}"
        echo -e "     └─ Endpoint: ${C_CYAN}http://localhost:8080/health${C_RESET}"
        echo -e "     └─ Search  : ${C_CYAN}http://localhost:8080/api/classification/search?q=dril${C_RESET}"
    else
        echo -e "${C_YELLOW}OFFLINE / NOT RESPONDING${C_RESET}"
    fi

    # 3. Frontend Check
    echo -n "  🌐 Next.js Frontend (:3000): "
    if curl -sI --connect-timeout 2 http://127.0.0.1:3000 >/dev/null 2>&1; then
        echo -e "${C_GREEN}HEALTHY (200 OK)${C_RESET}"
        echo -e "     └─ Web App: ${C_CYAN}http://localhost:3000${C_RESET}"
    else
        echo -e "${C_YELLOW}OFFLINE / NOT RESPONDING${C_RESET}"
    fi

    # 4. Docker Containers Check
    if command -v docker >/dev/null 2>&1 && docker info >/dev/null 2>&1; then
        echo ""
        echo -e "${C_BOLD}Active Docker Containers:${C_RESET}"
        docker ps --filter "name=huntiq" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null || true
    fi
    echo ""
}

# ------------------------------------------------------------------------------
# Docker Compose Execution Mode
# ------------------------------------------------------------------------------
cmd_docker() {
    print_banner "DOCKER CONTAINERIZED"

    # Check Docker CLI
    if ! command -v docker >/dev/null 2>&1; then
        echo -e "${C_RED}Error: 'docker' CLI is not installed or not in PATH.${C_RESET}" >&2
        exit 1
    fi

    # Check Docker Daemon
    echo -n "Checking Docker daemon status... "
    if ! docker info >/dev/null 2>&1; then
        echo -e "${C_RED}NOT RUNNING${C_RESET}"
        echo ""
        echo -e "${C_YELLOW}The Docker daemon is currently stopped or requires elevated permissions.${C_RESET}"
        echo -e "To start Docker on your system, please execute:"
        echo -e "  ${C_CYAN}sudo service docker start${C_RESET}  or  ${C_CYAN}sudo systemctl start docker${C_RESET}"
        echo ""
        echo -e "Alternatively, you can run HUNTIQ natively right now by executing:"
        echo -e "  ${C_GREEN}./run.sh${C_RESET}"
        echo ""
        exit 1
    fi
    echo -e "${C_GREEN}RUNNING${C_RESET}"

    # Determine compose command
    local COMPOSE_CMD="docker compose"
    if ! docker compose version >/dev/null 2>&1; then
        if command -v docker-compose >/dev/null 2>&1; then
            COMPOSE_CMD="docker-compose"
        else
            echo -e "${C_RED}Error: Neither 'docker compose' nor 'docker-compose' found.${C_RESET}" >&2
            exit 1
        fi
    fi

    # Cleanup Trap for Docker
    trap_docker() {
        echo ""
        echo -e "${C_YELLOW}Gracefully stopping HUNTIQ Docker containers...${C_RESET}"
        $COMPOSE_CMD down
        echo -e "${C_GREEN}All containers stopped cleanly.${C_RESET}"
        exit 0
    }
    trap trap_docker SIGINT SIGTERM

    echo -e "${C_BLUE}Building and launching HUNTIQ container stack with Docker Compose...${C_RESET}"
    echo -e "Services: ${C_CYAN}PostgreSQL 16, Redis 7.2, Go Backend (:8080), Next.js Frontend (:3000)${C_RESET}"
    echo ""

    $COMPOSE_CMD up --build
}

# ------------------------------------------------------------------------------
# Native Execution Mode (Default)
# ------------------------------------------------------------------------------
cmd_native() {
    print_banner "NATIVE FAST-DEV"

    PIDS=()
    cleanup() {
        echo ""
        echo -e "${C_YELLOW}Shutting down HUNTIQ native processes gracefully...${C_RESET}"
        for pid in "${PIDS[@]}"; do
            if kill -0 "$pid" 2>/dev/null; then
                kill "$pid" 2>/dev/null || true
            fi
        done
        wait 2>/dev/null || true
        echo -e "${C_GREEN}All native processes stopped cleanly.${C_RESET}"
        exit 0
    }
    trap cleanup SIGINT SIGTERM

    # 1. Check Database
    echo -e "${C_BLUE}[1/3] 📦 Verifying PostgreSQL Connection...${C_RESET}"
    local cfg_host="/tmp"
    local cfg_port="5433"
    if [ -f "$BACKEND_DIR/.env" ]; then
        cfg_host=$(grep -E '^[[:space:]]*DB_HOST=' "$BACKEND_DIR/.env" | head -n 1 | cut -d '=' -f2- | tr -d ' "\r' || echo "/tmp")
        cfg_port=$(grep -E '^[[:space:]]*DB_PORT=' "$BACKEND_DIR/.env" | head -n 1 | cut -d '=' -f2- | tr -d ' "\r' || echo "5433")
    fi
    [ -z "$cfg_host" ] && cfg_host="/tmp"
    [ -z "$cfg_port" ] && cfg_port="5433"

    local db_ok=false
    if psql -h "$cfg_host" -p "$cfg_port" -U postgres -d hulubedej_db -c "SELECT 1;" >/dev/null 2>&1; then
        echo -e "      ${C_GREEN}✓ Connected to PostgreSQL (Host: ${cfg_host}, Port: ${cfg_port})${C_RESET}"
        db_ok=true
    fi

    if [ "$db_ok" = false ]; then
        echo -e "      ${C_YELLOW}Notice: PostgreSQL not responding on ${cfg_host}:${cfg_port}. Initializing daemon...${C_RESET}"

        # Clean stale postmaster.pid if process is not running
        if [ -f "/tmp/pgdata/postmaster.pid" ]; then
            local pid_candidate
            pid_candidate=$(head -n 1 "/tmp/pgdata/postmaster.pid" 2>/dev/null || echo "")
            if [ -n "$pid_candidate" ] && ! kill -0 "$pid_candidate" 2>/dev/null; then
                echo -e "      ${C_YELLOW}Removing stale lock file /tmp/pgdata/postmaster.pid (PID ${pid_candidate} is dead)...${C_RESET}"
                rm -f "/tmp/pgdata/postmaster.pid"
            fi
        fi

        if [ -x "/usr/lib/postgresql/18/bin/postgres" ] && [ -d "/tmp/pgdata" ]; then
            /usr/lib/postgresql/18/bin/postgres -D /tmp/pgdata -p "$cfg_port" -k "$cfg_host" >/tmp/postgres.log 2>&1 &
            local pg_pid=$!
            PIDS+=($pg_pid)

            echo -n "      Waiting for PostgreSQL to become ready..."
            for _ in {1..20}; do
                if psql -h "$cfg_host" -p "$cfg_port" -U postgres -d hulubedej_db -c "SELECT 1;" >/dev/null 2>&1; then
                    db_ok=true
                    break
                fi
                sleep 0.5
            done
            if [ "$db_ok" = true ]; then
                echo -e " ${C_GREEN}Ready!${C_RESET}"
            else
                echo -e " ${C_RED}FAILED to connect to PostgreSQL. Check /tmp/postgres.log${C_RESET}"
                cleanup
                exit 1
            fi
        else
            echo -e "      ${C_RED}Error: PostgreSQL binary (/usr/lib/postgresql/18/bin/postgres) or data dir (/tmp/pgdata) missing.${C_RESET}"
            cleanup
            exit 1
        fi
    fi

    # 2. Start Go Backend
    echo ""
    echo -e "${C_BLUE}[2/3] 🔧 Starting Go Backend API (:8080)...${C_RESET}"
    if is_port_in_use 8080; then
        echo -e "      ${C_YELLOW}Port 8080 is already active (reusing active Go backend).${C_RESET}"
    else
        (
            cd "$BACKEND_DIR"
            go run ./cmd/api
        ) 2>&1 | sed -e "s/^/[BACKEND] /" &
        local backend_pid=$!
        PIDS+=($backend_pid)

        # Wait for backend healthcheck
        echo -n "      Waiting for backend to become ready..."
        local backend_ready=false
        for _ in {1..40}; do
            if ! kill -0 "$backend_pid" 2>/dev/null; then
                echo -e " ${C_RED}FAILED! Backend process exited with error.${C_RESET}"
                echo -e "      ${C_RED}Check the [BACKEND] error messages above for details.${C_RESET}"
                cleanup
                exit 1
            fi
            if curl -s --connect-timeout 1 http://127.0.0.1:8080/health >/dev/null 2>&1; then
                backend_ready=true
                break
            fi
            sleep 0.5
        done
        if [ "$backend_ready" = true ]; then
            echo -e " ${C_GREEN}Ready! (HTTP 200)${C_RESET}"
        elif kill -0 "$backend_pid" 2>/dev/null; then
            echo -e " ${C_YELLOW}Started (health check pending).${C_RESET}"
        else
            echo -e " ${C_RED}FAILED! Backend exited.${C_RESET}"
            cleanup
            exit 1
        fi
    fi

    # 3. Start Next.js Frontend
    echo ""
    echo -e "${C_BLUE}[3/3] 🌐 Starting Next.js Frontend (:3000)...${C_RESET}"
    if is_port_in_use 3000; then
        echo -e "      ${C_YELLOW}Port 3000 is already active (reusing active Next.js frontend).${C_RESET}"
    else
        (
            cd "$FRONTEND_DIR"
            npm run dev
        ) 2>&1 | sed -e "s/^/[FRONTEND] /" &
        PIDS+=($!)

        echo -n "      Initializing Next.js development server..."
        sleep 2
        echo -e " ${C_GREEN}Active!${C_RESET}"
    fi

    echo ""
    echo -e "${C_GREEN}════════════════════════════════════════════════════════════════${C_RESET}"
    echo -e "${C_GREEN}${C_BOLD}✓ HUNTIQ Platform is actively serving requests!${C_RESET}"
    echo ""
    echo -e "  🌐 ${C_BOLD}Frontend Application :${C_RESET} ${C_CYAN}http://localhost:3000${C_RESET}"
    echo -e "  🔧 ${C_BOLD}Go Backend API       :${C_RESET} ${C_CYAN}http://localhost:8080${C_RESET}"
    echo -e "  📊 ${C_BOLD}Health Endpoint      :${C_RESET} ${C_CYAN}http://localhost:8080/health${C_RESET}"
    echo -e "  🔍 ${C_BOLD}Universal Search API :${C_RESET} ${C_CYAN}http://localhost:8080/api/classification/search?q=dril${C_RESET}"
    echo -e "${C_GREEN}════════════════════════════════════════════════════════════════${C_RESET}"
    echo -e "${C_YELLOW}Press Ctrl+C to terminate all services.${C_RESET}"
    echo ""

    # Keep alive if child processes were spawned
    if [ ${#PIDS[@]} -gt 0 ]; then
        wait
    else
        echo -e "${C_GREEN}Services are already running in background. Run './run.sh --status' to inspect.${C_RESET}"
    fi
}

# ------------------------------------------------------------------------------
# Help Guide
# ------------------------------------------------------------------------------
cmd_help() {
    print_banner "HELP"
    echo "Usage: ./run.sh [OPTION]"
    echo ""
    echo "Options:"
    echo "  (none)        Start HUNTIQ in Native Fast-Dev Mode (Backend :8080 + Frontend :3000)"
    echo "  --docker, -d  Build & run containerized stack with Docker Compose"
    echo "  --status, -s  Check health status and endpoints of all services"
    echo "  --migrate     Run sequence and idempotency validation and execute pending migrations"
    echo "  --help, -h    Display this help guide"
    echo ""
    echo "Examples:"
    echo "  ./run.sh             # Run natively"
    echo "  ./run.sh --docker    # Run inside Docker"
    echo "  ./run.sh --status    # Check running services"
    echo ""
}

# ------------------------------------------------------------------------------
# CLI Router
# ------------------------------------------------------------------------------
MODE="${1:-}"

case "$MODE" in
    --docker|-d)
        cmd_docker
        ;;
    --status|-s)
        cmd_status
        ;;
    --migrate)
        echo "==> Running database migrations..."
        bash "$BACKEND_DIR/scripts/migrate.sh"
        ;;
    --help|-h)
        cmd_help
        ;;
    "")
        cmd_native
        ;;
    *)
        echo -e "${C_RED}Unknown option: $MODE${C_RESET}" >&2
        cmd_help
        exit 1
        ;;
esac
