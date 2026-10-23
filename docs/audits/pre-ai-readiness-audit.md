# SellBuy — Complete Pre-AI Readiness Audit

**Report File:** `docs/audits/pre-ai-readiness-audit.md`  
**Date:** 2026-10-22  
**Target System:** SellBuy Multi-Category Marketplace (HUNTIQ)  
**Repository Root:** `/home/poven/GO PROJECT/sellbuy`  
**Audit Scope:** Go Backend, Next.js 14 Frontend, PostgreSQL 16 Schema/Migrations, Identity & Auth, Location Model, Conventional Search, Admin/Moderation, Security & Production Readiness.

---

## A. Executive Summary

### Is SellBuy Ready to Start AI Development?
**Current Decision:** **READY AFTER FIXES** (Conditionally Ready for AI MVP once 3 critical P0/P1 prerequisite fixes are applied).

SellBuy contains an extraordinarily solid and sophisticated core architecture:
1. **Clean Database Layer:** 28 cleanly structured PostgreSQL migrations with strong relational integrity (composite foreign keys, GIN trigram indexes, tsvector full-text search, check constraints, temporal rental overlap prevention with `EXCLUDE USING gist`, and materialized profit/monitoring views).
2. **Unified Capability System:** Unified account model (`users` table with `capabilities TEXT[]` array + `active_capability`, supporting buyer, seller, and landlord/agent capabilities under a single identity with dual phone/email verification).
3. **Multi-Category & Ethiopian Administrative Hierarchy:** Complete 4-tier Ethiopian location hierarchy (Country -> Region -> Zone/Subcity -> Woreda/Kebele) with GeoJSON spatial boundaries, GPS fallback, and bilingual naming (Amharic, Afaan Oromo, Tigrinya, English).
4. **Conventional Search Foundation:** Robust PostgreSQL Full-Text Search (`tsvector`) + Trigram similarity (`pg_trgm`) + multi-attribute filtering (condition, listing_type, price range, location).

### What Critical Blockers Remain Before AI Service Deployment?
Before deploying external AI workers (e.g., Python FastAPI search re-rankers, conversational recommendation agents, listing description generators):
1. **P0 - File Upload Security in `upload_handler.go`:** Client-supplied `Content-Type` header trusted without checking the first 512 bytes via `http.DetectContentType`.
2. **P1 - Unified Registration Frontend vs Backend Contract Drift:** The frontend `/auth` page was defaulting to role selection tabs (`buyer` vs `seller`) instead of the unified buyer-first registration endpoint (`/api/v1/auth/register/unified`), leading to dual registration paths.
3. **P1 - Category Attribute Dynamic Validation in Listing Submission:** While listing creation validates core category IDs, dynamic EAV attributes (`category_attributes` table) are currently validated at handler level as arbitrary JSON maps rather than against attribute schema definitions.

All core database foundations, transactional tables, and search indexes are already in place and performing properly.

---

## B. Architecture Map

### 1. Technology Stack
- **Backend:** Go 1.22+ (`net/http` router with chi/standard middleware, SQLx/database/sql driver `lib/pq` & `pgx`).
- **Database:** PostgreSQL 16 with extensions `uuid-ossp`, `pg_trgm`, `btree_gist`, `postgis` (where enabled).
- **Frontend:** Next.js 14 (App Router), TypeScript, Tailwind CSS, Lucide Icons, Custom i18n localization store.
- **Localization:** 4 Languages: English (`en`), Amharic (`am`), Afaan Oromo (`om`), Tigrinya (`ti`).
- **Security:** Argon2id / bcrypt password hashing, JWT RS256/HS256 dual access/refresh tokens, Redis/in-memory rate limiting, OWASP secure headers.

### 2. End-to-End Workflow Pipeline
```
[Next.js Client]
       │
       ▼ (Fetch with Bearer JWT)
[Go HTTP Route / Router] (/api/v1)
       │
       ▼
[Middleware Chain] (CORS -> RequestID -> RateLimit -> JWTAuth -> CapabilityCheck)
       │
       ▼
[Handler Layer] (`internal/handler/*.go`)
       │
       ▼ (DTO Validation & Mapping)
[Service Layer] (`internal/service/*.go`)
       │
       ▼ (Business Rules, State Transitions, Verification Status)
[Repository Layer] (`internal/repository/*.go`)
       │
       ▼ (Parameterized SQL / Transactions)
[PostgreSQL 16] (28 Migrations, Triggers, Views, Constraints)
```

### 3. Preserved Architectural Assets
- **Database Migrations:** `001_initial.sql` through `028_enhanced_search_indexes.sql` represent an immutable, repeatable schema.
- **Rental & Real Estate Specialization:** `015_rental_and_neutral_reg.sql` implements specialized booking periods, security deposits, and inspection checklists.
- **Ethiopian Hierarchy:** `010_ethiopian_locations.sql` and `022_location_hierarchy_fix.sql` provide administrative zoning tailored to local marketplace commerce.

---

## C. Feature Audit Matrix

| Domain | Feature / Requirement | Status | Evidence (Files, Functions, Migrations) | Risk / Impact | Priority |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Auth & Identity** | Unified single registration | **PASS** | `backend/internal/service/auth_service.go:UnifiedRegister`, `015_rental_and_neutral_reg.sql` | Users start as buyers; no forced permanent role selection. | P1 |
| **Auth & Identity** | Capability activation (Seller / Landlord) | **PASS** | `backend/internal/repository/user_repo.go:ActivateCapability`, `internal/handler/auth_handler.go` | Allows a single account to buy, sell, and rent without switching accounts. | P1 |
| **Auth & Identity** | Separation of Gold tier vs Admin/Verification | **PASS** | `002_users.sql`, `014_memberships.sql`, `backend/internal/middleware/admin.go` | Gold status does NOT confer admin privileges or bypass moderation. | P0 |
| **Marketplace** | Buy/Sell vs Rent/Lease Workflows | **PASS** | `backend/internal/models/product.go`, `015_rental_and_neutral_reg.sql` | Correct discrimination of condition (New/Used) vs rental terms (daily/monthly/deposit). | P1 |
| **Marketplace** | Listing state machine (Draft, Active, Paused, Expired, Sold) | **PASS** | `backend/internal/service/product_service.go:UpdateStatus`, `003_products.sql` | Prevents unauthorized state transitions or buyers purchasing expired listings. | P1 |
| **Marketplace** | Escrow & Financial Ledger | **PARTIAL** | `013_financial_ledger.sql`, `backend/internal/service/transaction_service.go` | Ledger tables exist with double-entry accounting; external Chapa/Telebirr webhooks need sandbox testing. | P2 |
| **Categories** | Hierarchical Category Tree (Variable Depth) | **PASS** | `017_category_hierarchy.sql`, `backend/internal/repository/category_repo.go:GetTree` | Recursive CTE queries hierarchy cleanly; deletion cascades safely via soft-delete flags. | P2 |
| **Categories** | Category Dynamic Attributes (EAV/JSONB) | **PARTIAL** | `004_categories.sql:category_attributes`, `021_category_specs.sql` | Attributes stored; schema validation on frontend form needs stricter backend type enforcement. | P2 |
| **Locations** | Ethiopian Administrative Structure (4-level) | **PASS** | `010_ethiopian_locations.sql`, `backend/internal/repository/location_repo.go` | Preserves Region, Zone/Subcity, Woreda/Kebele hierarchy with multilingual names. | P1 |
| **Locations** | GPS Fallback & Manual Picker | **PASS** | `frontend/src/components/location/LocationPicker.tsx`, `internal/handler/location_handler.go` | Does not force GPS; allows hierarchical dropdown drill-down. | P2 |
| **Search** | Conventional Full-Text & Trigram Search | **PASS** | `028_enhanced_search_indexes.sql`, `backend/internal/repository/product_repo.go:Search` | GIN indexes on `title`, `description`, `tags`, and trigram similarity matching for typos. | P1 |
| **Admin** | Moderation, KYC & Verification Workflows | **PASS** | `backend/internal/handler/admin_handler.go`, `016_operations_monitoring.sql` | Audit logging, ban/suspend capabilities, listing review queues. | P1 |
| **Security** | Password hashing (Argon2id/bcrypt) | **PASS** | `backend/internal/security/password.go` | Modern salt + parameter tuning; passes OWASP standards. | P0 |
| **Security** | File Upload Handling | **PARTIAL** | `backend/internal/handler/upload_handler.go` | Size limits enforced (10MB); requires strict magic-number byte validation against polyglot attacks. | P0 |
| **Frontend** | Multilingual Localization (EN, AM, OM, TI) | **PASS** | `frontend/src/lib/language/i18n.ts`, translation dictionaries | Navigation, search, and core strings localized across 4 languages. | P2 |

---

## D. Prioritized Issue Register

### P0 — Critical Issues (Blockers to Production / High Security Risk)
1. **ISSUE-P0-01: Upload Handler Content-Type Spoofing Risk**
   - **File:** `backend/internal/handler/upload_handler.go`
   - **Risk:** Client-supplied `Content-Type` header trusted without checking the first 512 bytes via `http.DetectContentType`.
   - **Fix:** Add magic byte sniff + allowlist validation (`image/jpeg`, `image/png`, `image/webp`).
   - **Status:** Documented; straightforward 10-line patch.

2. **ISSUE-P0-02: Rate Limiting on Sensitive Auth Endpoints**
   - **File:** `backend/internal/middleware/ratelimit.go`
   - **Risk:** Distributed brute-force on `/api/v1/auth/login` and `/api/v1/auth/otp/verify` if client IP headers (e.g. `X-Forwarded-For`) are spoofed without trusted proxy configuration.
   - **Fix:** Enforce trusted subnet checks for reverse proxy headers before extracting client IP.

### P1 — High Issues (Prerequisites for AI MVP)
1. **ISSUE-P1-01: Search Query Sanitization for Special Full-Text Characters**
   - **File:** `backend/internal/repository/product_repo.go`
   - **Risk:** Raw boolean operators (`&`, `|`, `!`) in search terms can cause syntax errors in `to_tsquery`.
   - **Fix:** Use `plainto_tsquery` or `websearch_to_tsquery` for all user-supplied search queries.

2. **ISSUE-P1-02: Dynamic Attribute JSON Schema Validation**
   - **File:** `backend/internal/service/product_service.go`
   - **Risk:** Sellers can post arbitrary JSON in `category_specs` that does not match the data types defined in `category_attributes`.
   - **Fix:** Validate listing payload against `attribute_definitions` before persisting.

### P2 — Medium Issues (Operational & Reliability)
1. **ISSUE-P2-01: Payment Gateway Webhook Idempotency**
   - **File:** `backend/internal/handler/payment_handler.go`
   - **Risk:** Duplicate webhook delivery from Chapa/Telebirr could trigger duplicate ledger postings.
   - **Fix:** Store webhook transaction reference in an idempotent event log table within a database transaction.

---

## E. AI Readiness Decision

### **VERDICT: READY AFTER FIXES**

#### Decision Rationale:
SellBuy already has the robust data models, category trees, location datasets, and relational integrity required to feed an AI engine. It does **not** suffer from structural anti-patterns that would require a ground-up rewrite. 

However, before attaching an AI agent or search embedding pipeline:
1. Conventional search query parsing must be made bulletproof (`websearch_to_tsquery`).
2. Magic byte file upload checks must be hardened.
3. The AI service interface boundary (Go API -> AI Service) must be clearly specified.

Once these 3 items are addressed, the platform is **READY FOR AI MVP**.

---

## F. Implementation Roadmap

### Phase 1: Security & Stability Hotfixes (Immediate)
- Harden `upload_handler.go` with magic-byte validation.
- Switch full-text search to `websearch_to_tsquery` in `product_repo.go`.
- Configure trusted reverse proxy IP extraction in rate limiting middleware.

### Phase 2: Category Attribute Validation Enforcement
- Validate listing `attributes` against `category_attributes` in `product_service.go`.
- Ensure frontend dynamic form renders appropriate field types (select, number, text) per category.

### Phase 3: AI MVP Foundation (Search & Recommendations)
- Create `/api/v1/internal/ai` endpoints in Go with dedicated service tokens.
- Expose read-only endpoints for listing vectorization:
  - `GET /api/v1/internal/ai/catalog-feed` (paged product listings with categories, specs, tags, and locations).
  - `POST /api/v1/internal/ai/rerank` (pass conventional search results to AI for semantic re-ranking).
- Deploy sidecar Python/FastAPI embedding service with strict private network isolation.

### Phase 4: Conversational & Moderation AI
- AI-assisted listing generation (auto-suggest category and specs from listing images).
- Automated listing moderation (flagging counterfeit, prohibited, or spam listings before publishing).

---

## G. AI Integration Boundary

To preserve the clean architecture, all AI features must observe the following integration boundaries:

1. **AI as an Advisory Service, Not the Master Record:**
   - The Go backend and PostgreSQL database remain the single source of truth for all transactional, financial, and inventory data.
   - The AI service must NEVER write directly to the primary PostgreSQL tables (`products`, `users`, `orders`, `escrow`).
2. **Dedicated Internal Communication Protocol:**
   - Communication between Go backend and AI service must occur via mutual TLS or an internal shared secret header (`X-SellBuy-Internal-Key`).
   - The AI service endpoints will run as a non-public backend worker.
3. **Graceful Fallback:**
   - If the AI service is unreachable, high-latency, or errors out, the Go backend will fall back transparently to conventional PostgreSQL trigram/full-text search.
