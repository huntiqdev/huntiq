# HUNTIQ

Ethiopia's universal marketplace — buy, sell, and rent anything across Ethiopia.

---

## 🚀 One-Command Setup

```bash
git clone https://github.com/zemenuyes/huntiq.git && cd huntiq && bash setup.sh
```

This single command will:
- Install all Node.js frontend dependencies (`npm install`)
- Download all Go backend modules (`go mod download`)
- Copy `.env.example` files to `.env` / `.env.local`

---

## After Setup

1. **Configure backend** — edit `backend/.env`:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=huntiq_db
   JWT_SECRET=<run: openssl rand -base64 48>
   OTP_ENCRYPT_KEY=<run: openssl rand -base64 32>
   APP_BASE_URL=http://localhost:8080
   FRONTEND_URL=http://localhost:3000
   ALLOW_ORIGIN=http://localhost:3000
   APP_ENV=development
   ```

2. **Configure frontend** — edit `frontend/.env.local`:
   ```
   NEXT_PUBLIC_API_URL=http://localhost:8080/api
   BACKEND_URL=http://localhost:8080
   ```

3. **Create the database**:
   ```bash
   createdb huntiq_db
   ```

4. **Run backend** (terminal 1):
   ```bash
   cd backend && go run ./cmd/api
   ```

5. **Run frontend** (terminal 2):
   ```bash
   cd frontend && npm run dev
   ```

6. **Create admin account**:
   ```bash
   cd backend
   SEED_ADMIN_EMAIL=admin@huntiq.et \
   SEED_ADMIN_PHONE=0911234567 \
   SEED_ADMIN_PASSWORD=Admin@1234 \
   SEED_ADMIN_NAME="Super Admin" \
   go run ./cmd/seed
   ```

7. Open **http://localhost:3000** in your browser

---

## Requirements

| Tool | Version |
|------|---------|
| Node.js | 18+ |
| Go | 1.23+ |
| PostgreSQL | 14+ |

---

## Stack

| Layer | Tech |
|-------|------|
| Frontend | Next.js 15, React 18, Tailwind CSS |
| Backend | Go 1.23, Fiber v2 |
| Database | PostgreSQL + pgx/v5 |
| AI Microservice | Python 3.12+, FastAPI, Uvicorn, uv |
| Data Platform | Python 3.12+, Pydantic, JSONSchema |
| Deploy | Vercel + Render |

---

## Project Structure

```
huntiq/
├── backend/              # PILLAR 1: Authoritative Go API server & business logic
│   ├── cmd/              # API server and administrative CLI seeders
│   ├── internal/         # Handlers, services, repos, DB migrations
│   ├── pkg/              # Shared packages (dto, errors, etc.)
│   ├── docs/             # Architecture, RBAC matrix, auth operations, FP&A
│   └── go.mod
├── frontend/             # PILLAR 2: Next.js 15 web client & localized UI
│   ├── src/app/          # Next.js App Router pages
│   ├── src/components/   # UI components
│   ├── docs/             # Page roadmaps, UI audits, responsive/theme test reports
│   └── package.json
├── ai/                   # PILLAR 3: Self-contained AI Platform
│   ├── src/huntiq_ai/    # FastAPI prediction service, health checks, schemas
│   ├── data-platform/    # Offline dataset governance, contracts, and quality validator
│   ├── docs/             # AI architecture, ADR-001, security, runbooks, and audits
│   ├── tests/            # Contract, unit, and integration tests
│   └── pyproject.toml    # Managed with uv
├── setup.sh              # One-command installer
└── README.md
```

---

## 🧠 AI Platform Foundation

HUNTIQ includes an AI foundation designed with clear separation of concerns:
- **Go Backend Sovereignty**: Identity, transactions, listings, category definitions, and database writes remain strictly in the Go backend.
- **Stateless AI Worker**: The Python AI microservice (`ai/`) acts strictly as an internal, pre-shared-key authenticated worker for classification and search enhancement.
- **Self-Contained AI Platform**: All AI code, offline data governance (`ai/data-platform/`), and architecture documentation (`ai/docs/`) are self-contained in `ai/`.

See documentation for details:
- [AI Platform Documentation Hub](ai/docs/README.md)
- [AI Architecture & Threat Model](ai/docs/architecture/ai-architecture.md)
- [Go/Python Integration Contract](ai/docs/architecture/ai-integration.md)
- [ADR-001: Platform Boundaries](ai/docs/decisions/ADR-001-ai-platform-boundaries.md)
- [Local AI Development Runbook](ai/docs/runbooks/ai-local-development.md)

---

© 2026 HUNTIQ
