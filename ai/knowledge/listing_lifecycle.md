# Listing Lifecycle & Moderation States

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Complete state machine detailing the lifecycle of marketplace listings on HUNTIQ.

---

## 2. Detailed Technical & Operational Rules
- Draft: Created by user, not yet submitted for review.
- Pending Approval: Automated and manual screening for prohibited items (15-30 min).
- Active: Searchable and visible to all users (30-day default duration).
- Promoted: Highlighted or featured via Telebirr/CBE payment.
- Expired: Automatically hidden after 30 days; renewable for free.
- Sold / Inactive: Marked as sold by seller or removed upon request.

---

## 3. Enforcement & Verification Protocols
Implemented in backend/internal/models/listing.go.

---

## 4. Cross-System References
- Sovereign Go Authority: backend/internal/
- Python AI Inference Runtime: ai/app/
- Authoritative Marketplace Database: PostgreSQL 16
