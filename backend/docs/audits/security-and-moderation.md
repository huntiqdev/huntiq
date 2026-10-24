# SellBuy — Security & Marketplace Moderation Audit Report

**Report File:** `docs/audits/security-and-moderation.md`  
**Date:** 2026-10-22  
**Target Systems:** Go API Backend, PostgreSQL Database, Next.js Security Boundary

---

## 1. Verified Security Controls

| Security Domain | Control Implemented | Verification Evidence | Status |
| :--- | :--- | :--- | :--- |
| **Password Security** | Argon2id + Bcrypt salt + length check (≥ 8 chars) | `backend/internal/security/password.go` | **PASS** |
| **Session & Auth** | Dual token (JWT RS256/HS256) + active user check | `backend/internal/middleware/middleware.go:RequireAuth` | **PASS** |
| **Brute Force Defense** | IP-keyed rate limiters on login (10/min) & OTP (5/min) | `backend/cmd/api/main.go:authLimiter, otpLimiter` | **PASS** |
| **File Upload Defense** | 12-byte Magic byte sniffing + path traversal containment | `backend/internal/handler/upload_handler.go` | **PASS** |
| **SQL Injection** | Strict parameterized SQL queries (`$1, $2`) + `websearch_to_tsquery` | `backend/internal/repository/product_repo.go` | **PASS** |
| **Role Separation** | Gold membership decoupled from admin / permissions | `002_users.sql`, `014_memberships.sql` | **PASS** |
| **Audit Logging** | Granular admin action logging to `admin_logs` table | `backend/internal/service/admin_service.go:InsertAdminLog` | **PASS** |

---

## 2. Marketplace Moderation & Prohibited Items Policy

### A. Current Implementation State
- **Product Review Lifecycle:** All new product submissions enter `status = 'pending'` (`003_products.sql`, `product_service.go:CreateForSeller`).
- **Admin Review Queue:** Listings must be inspected and explicitly approved (`status = 'approved'`) or rejected (`status = 'rejected'`) via `POST /api/admin/products/:id/review`.
- **User Warnings & Bans:** Administrators can issue formal warnings with severity ratings (`low`, `medium`, `high`, `critical`) or suspend accounts via `013_reviews_and_warnings.sql`.

### B. Prohibited Items Categories
The platform prohibits listings in the following critical categories under Ethiopian and international marketplace regulations:
1. **Weapons & Explosives:** Firearms, handguns, rifles, automatic weapons, ammunition, grenades, gunpowder, military ordnance.
2. **Controlled Substances:** Narcotics, prescription drugs without pharmacy license, illegal stimulants, synthetic cannabinoids.
3. **Stolen Goods & Counterfeits:** Pirated media, counterfeit currency, stolen vehicles or electronics without ownership documentation.
4. **Endangered Wildlife:** Ivory, protected flora/fauna, live wild animal sales.
5. **Deceptive Financial Schemes:** Pyramid schemes, unregistered lottery tickets, fraudulent crypto/forex recruitment.

### C. Gaps Identified & Recommended Improvements
1. **Pre-Submission Deterministic Rule Engine:** Currently, prohibited items are not caught at the point of submission; they wait in the pending queue for human review.
   - *Fix:* Implement a fast deterministic prohibited-keyword pre-check in `product_service.go` that rejects or immediately flags high-risk listings (`flagged_for_review`) with a descriptive warning in the four supported languages.
2. **Missing Public Settings Endpoint Registration:**
   - *Issue:* `adminHandler.GetPublicSettings` (which exposes company phone, support email, and platform limits) exists in `admin_handler.go:1270` but was omitted from `backend/cmd/api/main.go` route definitions.
   - *Fix:* Register `api.Get("/public/settings", adminHandler.GetPublicSettings)` in `backend/cmd/api/main.go`.
