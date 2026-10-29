# HUNTIQ Backend Documentation (`backend/docs/`)

This directory contains authoritative backend architecture, security operations, database conventions, and administrative specifications for the Go Fiber / PostgreSQL sovereign core.

## 📚 Table of Contents

### 1. Architecture & Taxonomy
- [Universal Marketplace Taxonomy](architecture/universal-marketplace-taxonomy.md) — Canonical 4-level category hierarchy (2,380 categories), materialized path design, and Ethiopian multilingual support.

### 2. Administrative Roles & Governance
- [Administrative Roles & Responsibilities](admin/admin-roles-and-responsibilities.md) — Operational boundary definitions for the 6-tier RBAC system (`platform_owner`, `platform_admin`, `security_admin`, `operations_admin`, `moderator`, `support_agent`).

### 3. Security & Cryptographic Operations
- [Admin Authentication Operations](security/admin-authentication-operations.md) — Session management, Passkey (WebAuthn/FIDO2) lifecycle, and AES-256 encrypted SMS OTP.
- [Platform Owner Bootstrap](security/platform-owner-bootstrap.md) — First-time initialization via `cmd/seed`.
- [Admin Migration and Rollback](security/admin-migration-and-rollback.md) — PostgreSQL migration conventions, zero-downtime rules, and rollback procedures.

### 4. Master Verification & Audits
- [Master Platform Audit & Verification Ledger](../../docs/audits/MASTER-PLATFORM-AUDIT-AND-VERIFICATION.md) — Comprehensive master audit covering all 63 routes, security primitives, financial ledger reconciliation, and test catalogs.
