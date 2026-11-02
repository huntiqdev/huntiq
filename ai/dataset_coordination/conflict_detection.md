# Cross-Dataset Conflict Detection Protocol

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Automated continuous integration checks to detect conflicting label enums, contradictory intent classifications, and overlapping slot definitions across the 30 dataset modules.

---

## 2. Detailed Technical & Operational Rules
- Canonical Label Dictionary: All modules must draw category slugs strictly from the PostgreSQL taxonomy mirror.
- Currency Uniformity: Currency codes must be strictly 'ETB' across all dataset modules.
- Incompatible Constraint Screening: Detection of schemas that define mutually exclusive validation rules.

---

## 3. Enforcement & Verification Protocols
Executed via automated schema comparison scripts running during CI validation.

---

## 4. Cross-System References
- Sovereign Go Authority: backend/internal/
- Python AI Inference Runtime: ai/app/
- Authoritative Marketplace Database: PostgreSQL 16
