#!/bin/bash
# HUNTIQ — One-command setup for collaborators
# Run: bash setup.sh

set -e

echo "=============================="
echo "  HUNTIQ — Project Setup"
echo "  Architecture: Frontend + Backend + AI"
echo "=============================="

# ── 1. Frontend (Next.js 15) ──────────────────────────────
echo ""
echo "📦 Installing frontend dependencies..."
cd frontend
npm install
echo "✅ Frontend node_modules installed"

# Copy env if missing
if [ ! -f .env.local ]; then
  cp .env.example .env.local
  echo "📄 Created frontend/.env.local from .env.example"
  echo "   → Edit NEXT_PUBLIC_API_URL and NEXT_PUBLIC_BACKEND_HOST"
fi
cd ..

# ── 2. Backend (Go Fiber) ─────────────────────────────────
echo ""
echo "🔧 Installing backend Go modules..."
cd backend
go mod tidy
go mod download
echo "✅ Backend Go modules ready"

# Copy env if missing
if [ ! -f .env ]; then
  cp .env.example .env
  echo "📄 Created backend/.env from .env.example"
  echo "   → Fill in DB_*, JWT_SECRET, OTP_ENCRYPT_KEY before running"
fi
cd ..

# ── 3. AI Service (Python FastAPI) ────────────────────────
echo ""
echo "🤖 Setting up Python AI service dependencies..."
cd ai
if command -v uv >/dev/null 2>&1; then
  uv sync
  echo "✅ AI service virtual environment and dependencies synchronized"
else
  echo "⚠️ 'uv' is not installed. To install: curl -LsSf https://astral.sh/uv/install.sh | sh"
  echo "   Then run: cd ai && uv sync"
fi

# Copy env if missing
if [ ! -f .env ]; then
  cp .env.example .env
  echo "📄 Created ai/.env from .env.example"
  echo "   → Verify HUNTIQ_AI_SECRET_KEY matches backend config"
fi
cd ..

echo ""
echo "=============================="
echo "  Setup complete!"
echo "=============================="
echo ""
echo "Next steps:"
echo "  1. Edit backend/.env         (DB credentials, JWT_SECRET, OTP_ENCRYPT_KEY)"
echo "  2. Edit frontend/.env.local  (NEXT_PUBLIC_API_URL)"
echo "  3. Edit ai/.env              (HUNTIQ_AI_SECRET_KEY)"
echo "  4. Start Backend API:        cd backend && go run ./cmd/api"
echo "  5. Start Next.js Frontend:   cd frontend && npm run dev"
echo "  6. Start AI Worker:          cd ai && uv run uvicorn huntiq_ai.main:app --port 8000"
echo "  7. Seed Admin & Categories:  cd backend && go run ./cmd/seed"
echo ""
