# Marketplace Schema Synchronization Policy

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Defines protocols for updating AI dataset schemas whenever Go backend marketplace categories, pricing policies, or geocoding hierarchies evolve.

---

## 2. Detailed Technical & Operational Rules
- Automated Code Generation: Go struct definitions in internal/models generate matching JSON schemas for Python.
- Synchronous Version Bump: Any addition of a top-level category triggers a minor version bump across all 30 manifests.
- Backward Compatibility: Models must support previous category slugs for 60 days following a taxonomy migration.

---

## 3. Enforcement & Verification Protocols
Managed through Git release hooks and database migration listeners.

---

## 4. Cross-System References
- Sovereign Go Authority: backend/internal/
- Python AI Inference Runtime: ai/app/
- Authoritative Marketplace Database: PostgreSQL 16
