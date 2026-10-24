#!/bin/bash
# HUNTIQ — One-command setup for collaborators
# Run: bash setup.sh

set -e

echo "=============================="
echo "  HUNTIQ — Project Setup"
echo "=============================="

# ── Frontend ──────────────────────────────────────────────
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

# ── Backend ───────────────────────────────────────────────
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

echo ""
echo "=============================="
echo "  Setup complete!"
echo "=============================="
echo ""
echo "Next steps:"
echo "  1. Edit backend/.env  (DB credentials, JWT_SECRET, OTP_ENCRYPT_KEY)"
echo "  2. Edit frontend/.env.local  (NEXT_PUBLIC_API_URL)"
echo "  3. Run backend:   cd backend && go run ./cmd/api"
echo "  4. Run frontend:  cd frontend && npm run dev"
echo "  5. Create admin:  cd backend && go run ./cmd/seed"
echo ""
