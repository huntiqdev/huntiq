# HUNTIQ — Seven Economic Classification Systems & Database-Driven Governance Architecture

**Document:** Authoritative Enterprise Taxonomy Governance Architecture  
**Scope:** Multi-System Economic Classification (Physical Goods, Services, Jobs, Gigs, Opportunities, Media, Promotions)  
**System Version:** 2.1.0  
**Status:** **ACTIVE & INTEGRATED IN POSTGRESQL & GO BACKEND**

---

## 1. The Seven Economic Classification Systems

HUNTIQ classifies all commerce, employment, services, and opportunities across seven orthogonal economic systems:

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────┐
│                           HUNTIQ UNIVERSAL ECONOMIC CLASSIFICATION                              │
├───────────┬───────────────────────────────────┬───────────────────┬─────────────────────────────┤
│ Code      │ System Name                       │ Primary Entities  │ Transactional Modalities    │
├───────────┼───────────────────────────────────┼───────────────────┼─────────────────────────────┤
│ SYS_PROD  │ Products & Physical Assets        │ Goods, Machinery  │ Sale, Rent, Wholesale, Buy  │
│ SYS_SRV  │ Services & Professional Offerings │ Technical, Maint. │ Hourly, Fixed Contract, SLA │
│ SYS_EMP   │ Employment & Occupations          │ Jobs, Careers, CV │ Full-Time, Part-Time, Corp  │
│ SYS_GIG   │ Freelance & Online Gigs           │ Gigs, Creative    │ Milestone, Per-Task, Remote │
│ SYS_OPP   │ Tenders & Procurement             │ RFPs, Grants      │ Sealed Bid, Reverse Auction │
│ SYS_MED   │ News & Media Content              │ Notices, Media    │ Broadcast, Bulletin, Post   │
│ SYS_PRM   │ Promotions & Campaigns            │ Seasonal, Deals   │ Flash Sale, Bundle, Coupon  │
└───────────┴───────────────────────────────────┴───────────────────┴─────────────────────────────┘
```

### Database Implementation in PostgreSQL:
Table `classification_systems` serves as the foundational catalog:
- `code VARCHAR(50) PRIMARY KEY`
- `name_en`, `name_am`, `name_om`, `name_ti` (Localized display titles)
- `description TEXT`
- `icon VARCHAR(100)`
- `display_order INT`
- `is_active BOOLEAN`

`category_hierarchy` contains a foreign key:
```sql
ALTER TABLE category_hierarchy 
    ADD COLUMN classification_system_code VARCHAR(50) REFERENCES classification_systems(code) DEFAULT 'SYS_PROD';
```

---

## 2. Resolving the 139 Database Roots vs. 196 Root Families Architecture

A central governance principle of HUNTIQ is maintaining absolute clarity between **live database row counts** and **enterprise architectural expansion models**:

```
Enterprise Architecture Plan (196 Target Root Families)
├── Actual Database Roots in PostgreSQL (level = 0, parent_id IS NULL):    139
│     ├── Legacy Core Roots (IDs 1-11):                                      11
│     ├── Migrations 010/017 Roots (IDs 12-69):                              58
│     └── Migration 035 Canonical Roots (IDs 7601-7780):                     70
├── Subnodes in Database (instantiated at Level 1 or 2):                     19
├── Cross-System Occupational Talent Domains (System C & D):                 12
├── Logical Classification Groups (System E & G Modalities):                  7
└── Planned Enterprise Expansion Roots (Documented in Schemas, Not in DB):   76
```

### Governance Invariant:
**No synthetic placeholder rows may be inserted into `category_hierarchy`.**
New root rows are instantiated in PostgreSQL only when verified subcategories, localized aliases, attribute definitions, and active listings exist.

---

## 3. Database-Driven Taxonomy Governance Infrastructure

Taxonomy modification is strictly database-driven, removing hardcoded category assumptions from the frontend:

### A. Category Redirects Catalog (`category_redirects`)
Handles merged nodes, legacy URLs, and renamed categories seamlessly:
```sql
CREATE TABLE category_redirects (
    id SERIAL PRIMARY KEY,
    source_category_id INT,
    source_slug VARCHAR(255) NOT NULL UNIQUE,
    target_category_id INT NOT NULL REFERENCES category_hierarchy(id) ON DELETE CASCADE,
    redirect_type VARCHAR(50) NOT NULL DEFAULT 'PERMANENT',
    reason TEXT,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### B. Structured Administrative Audit Logging (`taxonomy_audit_logs`)
Maintains an immutable audit trail of all reparenting, attribute alterations, and alias changes:
```sql
CREATE TABLE taxonomy_audit_logs (
    id SERIAL PRIMARY KEY,
    admin_id INT REFERENCES users(id) ON DELETE SET NULL,
    category_id INT REFERENCES category_hierarchy(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    old_values JSONB,
    new_values JSONB,
    reason TEXT,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 4. Orthogonal Classification Dimensions

To prevent category explosion and conflated titles, the taxonomy enforces orthogonal separation between intrinsic classification and secondary dimensions:

1. **Intrinsic Category:** What the entity fundamentally *is* (e.g. `Laptop`, `Agricultural Tractor`, `Apartment`).
2. **Transaction Mode:** How the entity is traded (`sale`, `rent`, `lease`, `service_order`). Managed by `products.listing_type`.
3. **Trade Tier:** Scale of trade (`retail`, `wholesale`, `b2b_bulk`). Managed by attribute dimension `trade_tier`.
4. **Condition:** State of preservation (`new`, `used_like_new`, `refurbished`). Managed by `products.condition`.
5. **Location:** Geographic availability (`region`, `subcity`, `woreda`). Managed by `location_hierarchy`.
