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

# ── 3. AI Service (Local Qwen3 4B on Ollama) ────────────────
echo ""
echo "🤖 Setting up Local AI service (Ollama + Qwen3 4B)..."
bash ai/scripts/setup_local.sh
echo "✅ Local AI environment and models ready"

# Copy env if missing
if [ ! -f ai/.env ]; then
  cp ai/.env.example ai/.env
  echo "📄 Created ai/.env from .env.example"
fi

echo ""
echo "=============================="
echo "  Setup complete!"
echo "=============================="
echo ""
echo "Next steps:"
echo "  🚀 Start Everything in One Command: ./run.sh"
echo "     (or individually:)"
echo "     • Frontend  : cd frontend && npm run dev"
echo "     • Backend   : cd backend && go run ./cmd/api"
echo "     • Local AI  : ollama run huntiq-assistant"
echo "     • Verify AI : ./ai/scripts/verify_local.sh"
echo ""
