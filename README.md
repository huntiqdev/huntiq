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
| Deploy | Vercel + Render |

---

## Project Structure

```
huntiq/
├── backend/
│   ├── cmd/api/          # API server entry point
│   ├── cmd/seed/         # Admin account seeder
│   ├── internal/         # Handlers, services, repos
│   ├── pkg/              # Shared packages (dto, errors, etc.)
│   ├── .env.example      # Copy to .env
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── app/          # Next.js App Router pages
│   │   ├── components/   # UI components
│   │   └── lib/          # API client, utils, i18n
│   ├── .env.example      # Copy to .env.local
│   └── package.json
├── setup.sh              # One-command installer
└── README.md
```

---

© 2026 HUNTIQ
