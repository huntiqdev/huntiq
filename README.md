# HUNTIQ

Ethiopia's universal marketplace — buy, sell, and rent anything across Ethiopia.

## Stack

- **Frontend**: Next.js 15, React 18, Tailwind CSS, TypeScript
- **Backend**: Go (Fiber v2), PostgreSQL, pgx/v5
- **Deployment**: Vercel (frontend), Render (backend)

## Project Structure

```
huntiq/
├── frontend/     # Next.js 15 App Router
└── backend/      # Go Fiber API
```

## Getting Started

### Backend

```bash
cd backend
cp .env.example .env   # fill in real values
go run ./cmd/api
```

### Frontend

```bash
cd frontend
cp .env.example .env.local   # fill in NEXT_PUBLIC_API_URL
npm install
npm run dev
```

### Create Admin

```bash
cd backend
go run ./cmd/seed
```

## Features

- Universal listings: sale, rent, property, agriculture, services, jobs
- Seller & buyer tiers (Normal / Gold / Premium)
- OTP-based authentication (phone + email)
- Full admin platform with RBAC
- Multilingual: English, Amharic, Afaan Oromoo, Tigrinya
- Category hierarchy with dynamic attributes
- Referral commission system
- Background jobs: expiry, tier cleanup, OTP cleanup

## License

Proprietary — © 2026 HUNTIQ
