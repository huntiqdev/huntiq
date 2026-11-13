# HUNTIQ — Master Marketplace Platform Architecture & Operations Guide

> **Authoritative Engineering Reference & Deployment Manual**  
> **Platform Version:** 2.0.0 (Production Ready)  
> **Core Architecture:** PostgreSQL 18 (Sole Authority) + Go Fiber Backend + Next.js 16 Frontend

---

## 1. System Architecture Overview

HUNTIQ is a multi-vertical, multi-lingual commerce and exchange platform engineered for high reliability, strict transactional integrity, and zero dependence on mock services.

```
                              ┌───────────────────────────────────┐
                              │     Next.js 16 App Router UI      │
                              │ (React 19, TypeScript, Tailwind)  │
                              └─────────────────┬─────────────────┘
                                                │ REST / JSON (HTTP/2)
                                                ▼
                              ┌───────────────────────────────────┐
                              │      Go Fiber API Backend         │
                              │  (Clean Architecture, Fiber v2)   │
                              └─────────┬───────────────────┬─────┘
                                        │                   │
                     pgx/v5 Pool / TLS  │                   │ go-redis / Memory Fallback
                                        ▼                   ▼
                     ┌───────────────────────┐   ┌───────────────────────┐
                     │     PostgreSQL 18     │   │      Redis 7.2        │
                     │ (Authoritative State) │   │ (Session & Fast Cache)│
                     └───────────────────────┘   └───────────────────────┘
```

- **Backend:** Go (Fiber v2, `jackc/pgx/v5`, `go-webauthn`, JWT v5). Provides authoritative business validation, authentication, role enforcement, and taxonomy queries.
- **Frontend:** Next.js 16 (React 19, TypeScript, Tailwind CSS). Dynamic category cascaders, responsive layouts, client shell with toast/locale/auth context, full dark/light theme support.
- **Database:** PostgreSQL 18. Single source of truth for taxonomy, users, listings, capabilities, and audit trails.
- **Languages:** 4 fully localized languages supported across all UI strings and 2,575 database categories:
  - English (`en`)
  - Amharic (`am`)
  - Afaan Oromoo (`or`)
  - Tigrinya (`ti`)

---

## 2. Database & Migration Architecture

PostgreSQL is the single source of truth. All schema changes follow a strict, normalized 75-migration architecture (65 dedicated single-table migrations + 10 supporting routines & views) managed by [`backend/scripts/validate_migrations.sh`](./backend/scripts/validate_migrations.sh).

### Approved 31-Migration Baseline

```
001_initial_schema.sql
002_users.sql
003_locations.sql
004_categories.sql
005_products.sql
006_transactions_and_orders.sql
007_notifications_system.sql
008_admin_system.sql
009_referral_system.sql
010_tier_and_upgrades.sql
011_settings_and_limits.sql
012_security_and_otp.sql
013_reviews_and_warnings.sql
014_functions_and_triggers.sql
015_rental_and_neutral_registration.sql
016_operations_monitoring.sql
017_category_hierarchy.sql
018_attribute_system.sql
019_query_understanding.sql
020_profit_analytics_views.sql
021_seed_category_hierarchy.sql
022_location_hierarchy.sql
023_seed_ethiopian_locations.sql
024_multi_role_capabilities.sql
025_enterprise_admin_roles_and_security.sql
026_financial_planning_and_forecasting.sql
027_universal_taxonomy_and_catalogs.sql
028_taxonomy_governance_and_hierarchy_triggers.sql
029_taxonomy_classification_and_alias_resolution.sql
030_product_search_trigram_indexes.sql
031_user_sessions_and_security.sql
```

### Migration Safety Rules
1. Every migration is completely idempotent (`IF NOT EXISTS`, safe conflict handling).
2. All foreign-key references in seed/reconciliation migrations check `WHERE EXISTS` to allow execution on both clean databases and production snapshots.
3. Migration validation gate (`./backend/scripts/validate_migrations.sh`) is run in CI before any deployment.

---

## 3. Authoritative Taxonomy & Category Architecture

The category hierarchy is 100% database-driven in PostgreSQL:

- **Hierarchy Table:** `category_hierarchy` (contains 2,575 nodes from Level 0 roots down to Level 3/4 terminal leaves).
- **Paths:** Stored as materialized paths (`/root_id/child_id/.../`) maintained by recursive triggers.
- **Aliases:** `category_aliases` (1,048 multilingual search tokens for typo-tolerant trigram query resolution).
- **Redirects:** `category_redirects` (handles legacy flat slugs like `houses-sale` or `houses-rent`, redirecting them to canonical property nodes while separating `listing_type` = sale/rent).
- **Dynamic Attributes:** `category_attributes` & `attribute_definitions` allow arbitrary typed specifications (`number`, `text`, `select`, `boolean`, `range`) per category branch.

To verify taxonomy integrity against the database:
```bash
cd backend
go run ./cmd/taxonomy-tool verify-testcases
```

To export the authoritative canonical JSON dump:
```bash
cd backend
go run ./cmd/taxonomy-tool export-all-db-json
```

---

## 4. Security Architecture

1. **Authentication & Multi-Role Capabilities:**
   - Unified accounts with contextual capability activation (`buyer`, `seller`, `landlord`).
   - Session family tracking with refresh token replay detection and grace window.
   - Enforced phone number verification via encrypted OTP.
2. **Enterprise Admin & WebAuthn:**
   - Multi-role admin permissions (`super`, `platform_owner`, `operations`, `sellers`, `trust_safety`, `finance`, etc.).
   - Hardware-bound WebAuthn / FIDO2 Passkeys supported for administrative consoles.
   - Sudo re-authentication required for sensitive operational mutations.
3. **Transport & Cookie Security:**
   - Double-Submit CSRF protection with cryptographic tokens.
   - `HttpOnly`, `SameSite=Lax`, `Secure` session cookies.
   - Fiber security headers (`X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, `Content-Security-Policy`).
4. **Input Validation:**
   - Strict server-side validation on prices, image counts, text lengths, and Ethiopian phone numbers (`+2519...`, `+2517...`).
   - Prohibited goods filter blocks dangerous/narcotic keywords across English, Amharic, Oromo, and Tigrinya.

---

## 5. Development & Operations Guide

### Prerequisites
- Go 1.24+ (developed on Go 1.26)
- Node.js 20+ (recommended Node 22)
- PostgreSQL 16+ (or Docker)
- Redis 7+ (optional; automatically falls back to high-speed in-memory cache)

### Quick Start (One Command)

```bash
# Setup dependencies and environment configuration
bash setup.sh

# Run both Go API (:8080) and Next.js Frontend (:3000)
bash run.sh
```

### Docker Compose Deployment

```bash
docker compose up -d --build
```
This launches:
- `huntiq-postgres` on port 5432
- `huntiq-redis` on port 6379
- `huntiq-backend` on port 8080
- `huntiq-frontend` on port 3000

### Running Test Suites

**Backend Tests:**
```bash
cd backend
go test -v ./...
```

**Frontend Tests & Linting:**
```bash
cd frontend
npm run typecheck
npm run lint
npm run test
npm run build
```

**Database Migration Validation:**
```bash
./backend/scripts/validate_migrations.sh
```

---

## 6. Project Structure

```
huntiq/
├── backend/
│   ├── cmd/
│   │   ├── admin-recovery/    # CLI for administrative recovery
│   │   ├── api/               # Main Go API entrypoint & graceful lifecycle
│   │   ├── seed/              # Initial super admin bootstrapper
│   │   ├── taxonomy-tool/     # Master taxonomy CLI (verification, audit, i18n export & seeder)
│   │   └── test-session/      # Integration session tester
│   ├── internal/
│   │   ├── cache/             # Redis + in-memory fallback cache
│   │   ├── config/            # Environment & app configuration
│   │   ├── database/          # PostgreSQL connection pool & migrations
│   │   ├── fpa/               # Financial planning & analysis engine
│   │   ├── handler/           # HTTP controllers (Fiber v2)
│   │   ├── middleware/        # Auth, RBAC, CSRF, rate limiters
│   │   ├── models/            # Domain database models
│   │   ├── repository/        # SQL data access layers (pgx/v5)
│   │   ├── security/          # WebAuthn, JWT, OTP crypto, Passkeys
│   │   ├── service/           # Business logic & validations
│   │   └── worker/            # Background scheduled jobs & cleanup
│   ├── pkg/                   # Public DTOs, validations, constants, fees
│   └── scripts/
│       ├── migrate.sh         # Migration runner
│       └── validate_migrations.sh # Gate 4 safety checker
├── frontend/
│   ├── src/
│   │   ├── app/               # Next.js 16 App Router pages & API routes
│   │   ├── components/        # UI components (Admin, Category, Seller, Layout)
│   │   ├── context/           # AuthContext & LocaleContext
│   │   ├── hooks/             # Custom React hooks
│   │   ├── lib/               # API client, i18n, translations, validations
│   │   └── types/             # TypeScript definitions
│   └── tests/                 # Node test runner suite
├── docker-compose.yml         # Production multi-container composition
├── run.sh                     # Unified single-command development runner
└── setup.sh                   # Unified dependency and env setup script
```

