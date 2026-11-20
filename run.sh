#!/usr/bin/env bash
# ==============================================================================
# HUNTIQ — Fast Unified Orchestrator (Native & Docker)
# ==============================================================================
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND="$ROOT/backend"
FRONTEND="$ROOT/frontend"

# Read DB config from .env if present
DB_HOST="/tmp"
DB_PORT="5432"
if [ -f "$BACKEND/.env" ]; then
    DB_HOST=$(grep -E '^[[:space:]]*DB_HOST=' "$BACKEND/.env" | head -n 1 | cut -d '=' -f2- | tr -d ' "\r' || echo "/tmp")
    DB_PORT=$(grep -E '^[[:space:]]*DB_PORT=' "$BACKEND/.env" | head -n 1 | cut -d '=' -f2- | tr -d ' "\r' || echo "5432")
fi
[ -z "$DB_HOST" ] && DB_HOST="/tmp"
[ -z "$DB_PORT" ] && DB_PORT="5432"

PIDS=()

cleanup() {
    echo -e "\n🛑 Stopping all services..."
    for pid in "${PIDS[@]:-}"; do
        kill "$pid" 2>/dev/null || true
    done
    wait 2>/dev/null || true
    exit 0
}
trap cleanup SIGINT SIGTERM

in_use() {
    if command -v lsof >/dev/null 2>&1; then
        lsof -i ":$1" -sTCP:LISTEN >/dev/null 2>&1 && return 0
    fi
    if command -v ss >/dev/null 2>&1; then
        ss -tulpn 2>/dev/null | grep -q ":$1 " && return 0
    fi
    return 1
}

start_postgres() {
    # 1. Check primary configured host & port (5432)
    if psql -h "$DB_HOST" -p "$DB_PORT" -U postgres -d hulubedej_db -c "SELECT 1;" >/dev/null 2>&1; then
        return 0
    fi

    # 2. Check 127.0.0.1 on port 5432 if unix socket is not connected
    if [ "$DB_HOST" != "127.0.0.1" ] && psql -h "127.0.0.1" -p 5432 -U postgres -d hulubedej_db -c "SELECT 1;" >/dev/null 2>&1; then
        DB_HOST="127.0.0.1"
        DB_PORT="5432"
        return 0
    fi

    echo "📦 Starting local PostgreSQL daemon ($DB_PORT)..."
    if [ -f /tmp/pgdata/postmaster.pid ]; then
        local pid
        pid=$(head -n 1 /tmp/pgdata/postmaster.pid 2>/dev/null || echo "")
        if [ -n "$pid" ] && ! kill -0 "$pid" 2>/dev/null; then
            rm -f /tmp/pgdata/postmaster.pid 2>/dev/null || true
        fi
    fi
    if [ -x "/usr/lib/postgresql/18/bin/postgres" ] && [ -d "/tmp/pgdata" ]; then
        /usr/lib/postgresql/18/bin/postgres -D /tmp/pgdata -p "$DB_PORT" -k "$DB_HOST" >/tmp/postgres.log 2>&1 &
        PIDS+=($!)
        for _ in {1..20}; do
            psql -h "$DB_HOST" -p "$DB_PORT" -U postgres -d hulubedej_db -c "SELECT 1;" >/dev/null 2>&1 && return 0
            sleep 0.3
        done
    fi
    echo "❌ PostgreSQL failed to start on $DB_HOST:$DB_PORT" && exit 1
}

cmd_native() {
    echo "🚀 Starting HUNTIQ Native Stack..."
    start_postgres
    echo "✓ Database online (Port: $DB_PORT)"

    # Ensure database schema is migrated
    echo "⚙️ Running Go database migrations..."
    (cd "$BACKEND" && go run ./cmd/tools/migrate)

    if in_use 8080; then
        echo "✓ Backend API already running on :8080"
    else
        echo "🔧 Starting Go Backend (:8080)..."
        (cd "$BACKEND" && go run ./cmd/api) 2>&1 | sed -e 's/^/[BACKEND] /' &
        local backend_pid=$!
        PIDS+=($backend_pid)
        for _ in {1..30}; do
            curl -s http://127.0.0.1:8080/health >/dev/null 2>&1 && break
            sleep 0.3
        done
    fi

    if in_use 3000; then
        echo "✓ Frontend already running on :3000"
    else
        echo "🌐 Starting Next.js Frontend (:3000)..."
        (cd "$FRONTEND" && npm run dev) 2>&1 | sed -e 's/^/[FRONTEND] /' &
        PIDS+=($!)
    fi

    echo -e "\n══════════════════════════════════════════════════"
    echo "✓ HUNTIQ is ready:"
    echo "  • Web App : http://localhost:3000"
    echo "  • API     : http://localhost:8080"
    echo "  • Health  : http://localhost:8080/health"
    echo -e "══════════════════════════════════════════════════\n"
    [ ${#PIDS[@]} -gt 0 ] && wait
}

cmd_docker() {
    echo "🐳 Starting HUNTIQ in Docker Compose..."
    local compose="docker compose"
    docker compose version >/dev/null 2>&1 || compose="docker-compose"
    trap '$compose down; exit 0' SIGINT SIGTERM
    $compose up --build
}

cmd_status() {
    echo "📊 HUNTIQ Status Check:"
    if psql -h "$DB_HOST" -p "$DB_PORT" -U postgres -d hulubedej_db -c "SELECT 1;" >/dev/null 2>&1; then
        local cat_count
        cat_count=$(psql -h "$DB_HOST" -p "$DB_PORT" -U postgres -d hulubedej_db -t -A -c "SELECT count(*) FROM category_hierarchy;" 2>/dev/null || echo "0")
        echo "  • Database : ONLINE (Port $DB_PORT, $cat_count categories)"
    else
        echo "  • Database : OFFLINE"
    fi
    curl -s http://127.0.0.1:8080/health >/dev/null 2>&1 \
        && echo "  • Backend  : ONLINE (:8080)" || echo "  • Backend  : OFFLINE"
    curl -sI http://127.0.0.1:3000 >/dev/null 2>&1 \
        && echo "  • Frontend : ONLINE (:3000)" || echo "  • Frontend : OFFLINE"
}

case "${1:-}" in
    --docker|-d)  cmd_docker ;;
    --status|-s)  cmd_status ;;
    --migrate)    bash "$BACKEND/scripts/migrate.sh" ;;
    --help|-h)    echo "Usage: ./run.sh [--docker|-d | --status|-s | --migrate | --help|-h]" ;;
    "")           cmd_native ;;
    *)            echo "Unknown option: $1" && exit 1 ;;
esac
