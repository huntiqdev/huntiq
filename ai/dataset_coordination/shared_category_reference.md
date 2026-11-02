# Shared Marketplace Category Reference

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Authoritative dictionary linking live PostgreSQL category IDs with AI dataset category enums across all 8 primary verticals.

---

## 2. Detailed Technical & Operational Rules
- 8 Primary Verticals: Products & Electronics, Real Estate & Property, Vehicles, Jobs, Agriculture, Services, Hotels, Event Tickets.
- 45 Subcategories: Full hierarchical tree mapped to canonical slugs.
- Synonyms and Aliases: Local Ethiopian naming variations mapped to canonical slugs.

---

## 3. Enforcement & Verification Protocols
Mirrored directly from backend/internal/models/category.go.

---

## 4. Cross-System References
- Sovereign Go Authority: backend/internal/
- Python AI Inference Runtime: ai/app/
- Authoritative Marketplace Database: PostgreSQL 16
