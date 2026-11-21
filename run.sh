#!/usr/bin/env bash
# ==============================================================================
# HUNTIQ — Master Orchestrator (Single-Command Auto-Lifecycle)
#
# Usage:
#   ./run.sh            Start PostgreSQL, backend API (:8080), & frontend (:3000)
#   ./run.sh --status   Check status of all stack components
#   ./run.sh --stop     Stop all running HUNTIQ processes & free ports
#   ./run.sh --migrate  Run database migrations only
#   ./run.sh --seed     Seed database master data only
#   ./run.sh --docker   Start via Docker Compose
# ==============================================================================
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND="$ROOT/backend"
FRONTEND="$ROOT/frontend"

# Database defaults
DB_HOST="127.0.0.1"
DB_PORT="5432"
DB_PASSWORD="gedam"
if [ -f "$BACKEND/.env" ]; then
    DB_HOST=$(grep -E '^[[:space:]]*DB_HOST=' "$BACKEND/.env" | head -n 1 | cut -d '=' -f2- | tr -d ' "\r' || echo "127.0.0.1")
    DB_PORT=$(grep -E '^[[:space:]]*DB_PORT=' "$BACKEND/.env" | head -n 1 | cut -d '=' -f2- | tr -d ' "\r' || echo "5432")
    DB_PASSWORD=$(grep -E '^[[:space:]]*DB_PASSWORD=' "$BACKEND/.env" | head -n 1 | cut -d '=' -f2- | tr -d ' "\r' || echo "gedam")
fi
[ -z "$DB_HOST" ] && DB_HOST="127.0.0.1"
[ -z "$DB_PORT" ] && DB_PORT="5432"
export PGPASSWORD="${DB_PASSWORD:-gedam}"

# Process tracking
PG_PID=""
BACKEND_PID=""
FRONTEND_PID=""
TAIL_API_PID=""
TAIL_APP_PID=""
PG_STARTED_BY_US=0
CLEANED_UP=0

in_use() {
    local port="$1"
    if command -v lsof >/dev/null 2>&1; then
        lsof -i ":$port" -sTCP:LISTEN >/dev/null 2>&1 && return 0
    fi
    if command -v ss >/dev/null 2>&1; then
        ss -tulpn 2>/dev/null | grep -q ":$port " && return 0
    fi
    return 1
}

free_port() {
    local port="$1"
    if command -v fuser >/dev/null 2>&1; then
        fuser -k -9 "${port}/tcp" 2>/dev/null || true
    fi
    if command -v lsof >/dev/null 2>&1; then
        local pids
        pids=$(lsof -ti ":$port" 2>/dev/null || echo "")
        for p in $pids; do
            kill -9 "$p" 2>/dev/null || true
        done
    fi
    if [ "$port" = "8080" ]; then
        pkill -9 -f "huntiq_backend" 2>/dev/null || true
    fi
    sleep 0.3
}

cleanup() {
    if [ "$CLEANED_UP" -eq 1 ]; then
        return
    fi
    CLEANED_UP=1
    echo -e "\n🛑 Stopping all HUNTIQ services..."

    # 1. Stop log tailers
    [ -n "$TAIL_API_PID" ] && kill "$TAIL_API_PID" 2>/dev/null || true
    [ -n "$TAIL_APP_PID" ] && kill "$TAIL_APP_PID" 2>/dev/null || true

    # 2. Stop frontend (:3000)
    if [ -n "$FRONTEND_PID" ]; then
        kill "$FRONTEND_PID" 2>/dev/null || true
        pkill -P "$FRONTEND_PID" 2>/dev/null || true
    fi
    free_port 3000

    # 3. Stop backend (:8080)
    if [ -n "$BACKEND_PID" ]; then
        kill "$BACKEND_PID" 2>/dev/null || true
        pkill -P "$BACKEND_PID" 2>/dev/null || true
    fi
    free_port 8080

    # 4. Stop PostgreSQL if started by this script
    if [ "$PG_STARTED_BY_US" -eq 1 ] && [ -n "$PG_PID" ]; then
        echo "🛑 Stopping PostgreSQL daemon (PID $PG_PID)..."
        kill -TERM "$PG_PID" 2>/dev/null || true
        for _ in {1..15}; do
            kill -0 "$PG_PID" 2>/dev/null || break
            sleep 0.2
        done
        kill -9 "$PG_PID" 2>/dev/null || true
    fi

    # Clean stale locks if no postgres is running
    local pg_dir="/tmp/huntiq_pgdata"
    [ ! -d "$pg_dir" ] && [ -d "/tmp/pgdata" ] && pg_dir="/tmp/pgdata"
    if ! pgrep -x postgres >/dev/null 2>&1; then
        rm -f /tmp/.s.PGSQL.* "${pg_dir}/postmaster.pid" 2>/dev/null || true
    fi

    echo "✓ All HUNTIQ services stopped cleanly."
    exit 0
}

start_postgres() {
    # Check if PostgreSQL is already running and accessible
    if psql -h "$DB_HOST" -p "$DB_PORT" -U postgres -d hulubedej_db -c "SELECT 1;" >/dev/null 2>&1; then
        echo "✓ PostgreSQL already running and accepting connections."
        return 0
    fi

    echo "📦 Starting local PostgreSQL daemon ($DB_PORT)..."
    local pg_dir="/tmp/huntiq_pgdata"
    [ ! -d "$pg_dir" ] && [ -d "/tmp/pgdata" ] && pg_dir="/tmp/pgdata"

    if [ ! -d "$pg_dir" ]; then
        echo "❌ PostgreSQL data directory ($pg_dir) not found!"
        exit 1
    fi

    # Fix shared memory type if set to posix (which fails in sandboxed environments)
    if grep -q "dynamic_shared_memory_type = posix" "$pg_dir/postgresql.conf" 2>/dev/null; then
        sed -i "s/^dynamic_shared_memory_type = posix/dynamic_shared_memory_type = sysv/" "$pg_dir/postgresql.conf"
    fi

    # Clean stale lock files and ensure port is completely free
    free_port "$DB_PORT"
    rm -f /tmp/.s.PGSQL.* "$pg_dir/postmaster.pid" 2>/dev/null || true

    # Locate postgres binary
    local pg_bin="/usr/lib/postgresql/18/bin/postgres"
    if [ ! -x "$pg_bin" ]; then
        pg_bin=$(which postgres 2>/dev/null || echo "")
    fi
    if [ -z "$pg_bin" ] || [ ! -x "$pg_bin" ]; then
        echo "❌ postgres binary not found!"
        exit 1
    fi

    # Start postgres (accepts local TCP and unix sockets)
    "$pg_bin" -D "$pg_dir" -p "$DB_PORT" -k "/tmp" -h "127.0.0.1" >/tmp/postgres.log 2>&1 &
    PG_PID=$!
    PG_STARTED_BY_US=1

    # Poll for database ready
    local ready=0
    for _ in {1..30}; do
        if psql -h "$DB_HOST" -p "$DB_PORT" -U postgres -d hulubedej_db -c "SELECT 1;" >/dev/null 2>&1; then
            ready=1
            break
        fi
        sleep 0.2
    done

    if [ "$ready" -ne 1 ]; then
        echo "❌ PostgreSQL failed to start on $DB_HOST:$DB_PORT. Log:"
        cat /tmp/postgres.log 2>/dev/null || true
        exit 1
    fi
}

cmd_native() {
    # Register trap for clean shutdown on Ctrl+C, kill, terminal close, or exit
    trap cleanup SIGINT SIGTERM SIGHUP EXIT

    echo "🚀 Initializing HUNTIQ Stack..."

    # Step 1: Clean stale processes on ports 8080 and 3000
    free_port 8080
    free_port 3000

    # Step 2: Ensure PostgreSQL is up
    start_postgres
    echo "✓ Database online (Host: $DB_HOST, Port: $DB_PORT)"

    # Step 3: Run database migrations and verify master data
    echo "⚙️ Running Go database migrations & verifying master data..."
    (cd "$BACKEND" && go run ./cmd/api -migrate-only)

    # Step 4: Build & start Go Backend (:8080)
    echo "🔧 Building and launching Go Backend (:8080)..."
    local bin_path="/tmp/huntiq_backend"
    (cd "$BACKEND" && go build -o "$bin_path" ./cmd/api)
    free_port 8080
    (cd "$BACKEND" && exec "$bin_path") >/tmp/huntiq_backend.log 2>&1 &
    BACKEND_PID=$!

    # Wait for backend health
    local backend_ready=0
    for _ in {1..40}; do
        if curl -s http://127.0.0.1:8080/health >/dev/null 2>&1; then
            backend_ready=1
            break
        fi
        sleep 0.25
    done
    if [ "$backend_ready" -ne 1 ]; then
        echo "❌ Backend failed to start. Log:"
        tail -n 30 /tmp/huntiq_backend.log 2>/dev/null || true
        exit 1
    fi
    echo "✓ Backend API online (http://localhost:8080)"

    # Step 5: Start Next.js Frontend (:3000)
    echo "🌐 Starting Next.js Frontend (:3000)..."
    free_port 3000
    (cd "$FRONTEND" && npm run dev) >/tmp/huntiq_frontend.log 2>&1 &
    FRONTEND_PID=$!

    # Wait for frontend port
    local frontend_ready=0
    for _ in {1..40}; do
        if in_use 3000; then
            frontend_ready=1
            break
        fi
        sleep 0.3
    done
    if [ "$frontend_ready" -ne 1 ]; then
        echo "❌ Frontend failed to bind port 3000. Log:"
        tail -n 30 /tmp/huntiq_frontend.log 2>/dev/null || true
        exit 1
    fi
    echo "✓ Next.js Frontend online (http://localhost:3000)"

    # Step 6: Print ready banner
    echo -e "\n══════════════════════════════════════════════════════════════════"
    echo "✓ HUNTIQ IS READY AND SERVING:"
    echo "  • Web Application : http://localhost:3000"
    echo "  • REST API        : http://localhost:8080"
    echo "  • Health Check    : http://localhost:8080/health"
    echo "  • Mode            : Pure Manual Marketplace"
    echo "──────────────────────────────────────────────────────────────────"
    echo "  Log streaming active. Press [Ctrl+C] or close terminal to stop."
    echo -e "══════════════════════════════════════════════════════════════════\n"

    # Step 7: Stream logs cleanly
    tail -n 0 -f /tmp/huntiq_backend.log 2>/dev/null | sed -e 's/^/[API] /' &
    TAIL_API_PID=$!
    tail -n 0 -f /tmp/huntiq_frontend.log 2>/dev/null | sed -e 's/^/[APP] /' &
    TAIL_APP_PID=$!

    # Step 8: Monitor health of processes
    while true; do
        if ! kill -0 "$BACKEND_PID" 2>/dev/null; then
            echo -e "\n❌ Backend process (PID $BACKEND_PID) exited unexpectedly."
            tail -n 20 /tmp/huntiq_backend.log 2>/dev/null || true
            break
        fi
        if ! kill -0 "$FRONTEND_PID" 2>/dev/null; then
            echo -e "\n❌ Frontend process (PID $FRONTEND_PID) exited unexpectedly."
            tail -n 20 /tmp/huntiq_frontend.log 2>/dev/null || true
            break
        fi
        sleep 2
    done
}

cmd_stop() {
    echo "🛑 Stopping all HUNTIQ services..."
    free_port 3000
    free_port 8080
    pkill -f "huntiq_backend" 2>/dev/null || true
    pkill -f "next-server" 2>/dev/null || true
    local pg_dir="/tmp/huntiq_pgdata"
    [ ! -d "$pg_dir" ] && [ -d "/tmp/pgdata" ] && pg_dir="/tmp/pgdata"
    if [ -d "$pg_dir" ] && [ -x "/usr/lib/postgresql/18/bin/pg_ctl" ]; then
        /usr/lib/postgresql/18/bin/pg_ctl stop -D "$pg_dir" -m fast 2>/dev/null || true
    fi
    pkill -x postgres 2>/dev/null || true
    rm -f /tmp/.s.PGSQL.* "${pg_dir}/postmaster.pid" 2>/dev/null || true
    echo "✓ All services stopped."
}

cmd_docker() {
    echo "🐳 Starting HUNTIQ in Docker Compose..."
    local compose="docker compose"
    docker compose version >/dev/null 2>&1 || compose="docker-compose"
    trap '$compose down; exit 0' SIGINT SIGTERM SIGHUP EXIT
    $compose up --build
}

cmd_status() {
    echo "📊 HUNTIQ Stack Status Check:"
    if psql -h "$DB_HOST" -p "$DB_PORT" -U postgres -d hulubedej_db -c "SELECT 1;" >/dev/null 2>&1; then
        local cat_count
        cat_count=$(psql -h "$DB_HOST" -p "$DB_PORT" -U postgres -d hulubedej_db -t -A -c "SELECT count(*) FROM category_hierarchy;" 2>/dev/null || echo "0")
        echo "  • Database : ONLINE (Port $DB_PORT, $cat_count categories verified)"
    else
        echo "  • Database : OFFLINE"
    fi

    if curl -s http://127.0.0.1:8080/health >/dev/null 2>&1; then
        echo "  • Backend  : ONLINE (http://localhost:8080)"
    else
        echo "  • Backend  : OFFLINE (:8080)"
    fi

    if in_use 3000; then
        echo "  • Frontend : ONLINE (http://localhost:3000)"
    else
        echo "  • Frontend : OFFLINE (:3000)"
    fi
}

case "${1:-}" in
    --docker|-d)  cmd_docker ;;
    --status|-s)  cmd_status ;;
    --stop)       cmd_stop ;;
    --migrate)    (cd "$BACKEND" && go run ./cmd/api -migrate-only) ;;
    --seed)       (cd "$BACKEND" && go run ./cmd/api -seed-only) ;;
    --help|-h)    echo "Usage: ./run.sh [--docker|-d | --status|-s | --stop | --migrate | --seed | --help|-h]" ;;
    "")           cmd_native ;;
    *)            echo "Unknown option: $1" && exit 1 ;;
esac
