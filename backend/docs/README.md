# HUNTIQ Backend Documentation (`backend/docs/`)

This directory contains authoritative backend architecture, security operations, database migrations, and financial ledger specifications for the Go Fiber / PostgreSQL backend.

## 📚 Table of Contents

### 1. Architecture & Data Model
- [Universal Marketplace Taxonomy](architecture/universal-marketplace-taxonomy.md) — Canonical category hierarchy, leaf catalogs, and database seeder design.

### 2. Administrative Role Taxonomy & Security
- [Official Admin Role Taxonomy & Authorization Matrix](audits/admin-authorization-matrix.md) — Enterprise RBAC matrix (`platform_owner`, `platform_admin`, `security_admin`, `operations_admin`, `moderator`, `support_agent`).
- [Admin Authentication Operations](security/admin-authentication-operations.md) — Session management, JWT/cookie lifecycle, and OTP keys.
- [Admin Session Security](audits/admin-session-security.md) — Session stores, sliding expiration, and concurrent login rules.
- [Admin Login Threat Model](audits/admin-login-threat-model.md) — Threat model, brute-force protections, and audit logging.
- [Admin Passkey Implementation](audits/admin-passkey-implementation.md) — WebAuthn / Passkey credential flows.
- [Admin Security Monitoring](audits/admin-security-monitoring.md) — Security alerting, suspicious event monitoring, and audit log schemas.
- [Admin Post-Implementation Security Review](audits/admin-post-implementation-security-review.md) — Verification of backend security hardening.
- [Admin Recovery Procedures](audits/admin-recovery-procedures.md) — Credential reset and emergency recovery protocols.

### 3. Database & Operations
- [Admin Migration and Rollback](security/admin-migration-and-rollback.md) — PostgreSQL migration conventions, zero-downtime rules, and rollback procedures.
- [Admin Role Migration](audits/admin-role-migration.md) — Migration plan for legacy `super`/`normal` roles to enterprise role taxonomy.
- [Platform Owner Bootstrap](security/platform-owner-bootstrap.md) — First-time initialization via `cmd/seed`.
- [Administrative Roles & Responsibilities](admin/admin-roles-and-responsibilities.md) — Operational boundary definitions.

### 4. Financials & Audits
- [FP&A Financial Ledger Reconciliation](audits/fpa-financial-reconciliation.md) — Financial transactions, fees, escrow, and audit reconciliation logic.
- [FP&A Quarterly Growth Reconciliation](audits/fpa-growth-reconciliation.md) — Financial analytics, cash flow, and growth forecasting.
- [Security & Moderation Rules](audits/security-and-moderation.md) — Content moderation policies, dispute mediation, and seller tiers.
- [Membership Access & Product Visibility](audits/membership-access-and-product-visibility.md) — Subscription gates, listing promotion, and search visibility.
- [Dead Code & Cleanup Review](audits/cleanup-review.md) — Backend legacy code removal and handler cleanup audit.
