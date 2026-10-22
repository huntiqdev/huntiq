# Enterprise Search Engine Upgrade - Implementation Document

**Project**: Habesha Market Search System Upgrade  
**Start Date**: October 21, 2026  
**Status**: Phase 1 - Audit & Baseline (IN PROGRESS)

---

## Executive Summary

This document tracks the upgrade of the Habesha Market search system from its current working implementation to an enterprise-grade dynamic hierarchical search engine capable of supporting unlimited marketplace categories, Ethiopian localization, and intent-aware search across products, services, property, employment, and agriculture sectors.

**Objectives**:
1. Preserve existing working search functionality
2. Implement data-driven category hierarchy (configurable depth)
3. Add dynamic category-specific attributes and filters
4. Improve query understanding (brands, models, synonyms, typos)
5. Implement intent-aware search (sale vs rent vs jobs)
6. Optimize performance and relevance
7. Integrate with existing admin system
8. Support Ethiopian localization and terminology

---

## Phase 1: Audit and Baseline

### Current Architecture Discovery

#### Technology Stack
- **Backend**: Go 1.x, Fiber framework
- **Database**: PostgreSQL (with pgx/v5 driver)
- **Search**: PostgreSQL Full-Text Search (`search_vector`, `plainto_tsquery`)
- **Frontend**: Next.js, TypeScript, React

#### Existing Search Implementation

**Repository Layer** (`backend/internal/repository/product_repo.go`):
```go
func (r *ProductRepo) ListApproved(ctx context.Context, f dto.ProductFilter) ([]models.Product, error)
```

**Current Filters** (`pkg/dto/product.go`):
```go
type ProductFilter struct {
    CategoryID  *int
    Subcategory *string
    MinPrice    *float64
    MaxPrice    *float64
    LocationID  *int
    Region      *string
    Zone        *string
    Subcity     *string
    Woreda      *string
    Condition   *string
    Search      *string      // Full-text search
    Page        int
    Limit       int
    ListingType *string      // 'sale' or 'rent'
    Item        *string
}
```

**Current Search Query**:
```sql
SELECT p.*, c.*, l.*, u.tier
FROM products p
LEFT JOIN categories c ON c.id = p.category_id
LEFT JOIN locations l ON l.id = p.location_id
LEFT JOIN users u ON u.id = p.seller_id
WHERE p.status = 'approved'
  AND (p.expires_at IS NULL OR p.expires_at > NOW())
  AND u.is_active = true
  AND search_vector @@ plainto_tsquery('simple', $search_term)  -- Full-text search
  AND p.category_id = $category_id  -- Category filter
  AND p.price >= $min_price AND p.price <= $max_price  -- Price range
  AND l.woreda ILIKE $woreda  -- Location hierarchy
  AND p.listing_type = $listing_type  -- Sale vs Rent
ORDER BY p.created_at DESC
LIMIT $limit OFFSET $offset
```

#### Current Category System

**Database Schema** (`backend/internal/database/migrations/004_categories.sql`):
```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);
```

**Product Model**:
- `category_id` (references categories table)
- `subcategory` (VARCHAR, not structured)
- `item` (VARCHAR, not structured)

**Current Limitations**:
1. ❌ Category hierarchy is flat (only 1 level in DB, subcategory/item are text fields)
2. ❌ No parent-child relationships
3. ❌ Cannot add new hierarchy levels without schema changes
4. ❌ No category-specific attributes
5. ❌ No dynamic filter generation based on category

#### Current Location System

**Database Schema** (`backend/internal/database/migrations/003_locations.sql`):
```sql
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    region VARCHAR(100),
    zone VARCHAR(100),
    subcity VARCHAR(100),
    woreda VARCHAR(100)
);
```

**Location Hierarchy**: Region → Zone → Subcity → Woreda (4 levels, hardcoded)

#### Current Search Features ✅

**What Works Well**:
1. ✅ PostgreSQL full-text search with `search_vector`
2. ✅ Basic filtering (category, price, location, condition)
3. ✅ Sale vs Rent listing type distinction
4. ✅ Location hierarchy filtering (Region → Woreda)
5. ✅ Pagination and sorting
6. ✅ Seller tier integration
7. ✅ Expired listing exclusion
8. ✅ Active seller filtering
9. ✅ Popular categories tracking

**Current Ranking**: Created date DESC (simple, predictable, works)

#### Frontend Search Pages

**Existing Pages**:
- `/find` - Main search/browse page
- `/demand` - Demand listings
- `/list` - Seller listings management
- `/rent` - Rental property listings

**Search UI Components** (to be inspected):
- Category filters
- Price range slider
- Location selector
- Condition filter
- Search input

#### API Endpoints

**Current**:
- `GET /api/products` - List approved products with filters
- `GET /api/products/:id` - Get single product
- `GET /api/locations/search` - Location search

---

## Baseline Test Collection

### Representative Query Patterns

#### Product Searches
```
1. "Samsung Galaxy M14" - Exact model
2. "Samsung phone" - Brand + category
3. "used Samsung" - Condition + brand
4. "phone under 5000" - Category + price
5. "electronics" - Broad category
```

#### Property Searches
```
6. "house for rent in Addis Ababa" - Property + intent + location
7. "2 bedroom apartment" - Property attributes
8. "rental property Bahir Dar" - Alternative terminology
```

#### Agricultural Searches
```
9. "tractor for rent" - Equipment + intent
10. "ox ploughing service" - Service type
11. "farm land for sale" - Land + intent
```

#### Employment Searches
```
12. "agriculture job" - Job category
13. "tractor operator" - Specific role
14. "farm labor Oromia" - Job + location
```

#### Vehicle Searches
```
15. "Toyota Corolla 2015" - Make + model + year
16. "used car under 200000" - Condition + category + price
17. "motorcycle spare parts" - Vehicle category + item type
```

### Current Behavior Baseline

**To Be Measured**:
- [ ] Query response time (median, p95, p99)
- [ ] Result relevance (manual evaluation on 50 queries)
- [ ] Zero-result rate
- [ ] Filter usage patterns
- [ ] Click-through rate on top results

---

## Gap Analysis

### Critical Gaps Identified

#### 1. Category Hierarchy ❌
**Current**: Flat categories + text subcategory/item fields  
**Required**: Dynamic N-level parent-child relationships

**Impact**: Cannot support deep category trees without schema changes

#### 2. Category-Specific Attributes ❌
**Current**: All products have same fields (title, price, description, condition)  
**Required**: Phone has "brand, model, storage"; Vehicle has "make, model, year, mileage"

**Impact**: Cannot filter by phone brand or vehicle mileage

#### 3. Query Understanding ❌
**Current**: Simple full-text match on title + description  
**Required**: Recognize brands, models, synonyms, intent, location

**Example**:
- "Samsung M14" → should match "Samsung Galaxy M14"
- "mobile" → should match "phone" listings
- "house for rent" → should filter to rental properties

#### 4. Relevance Ranking ❌
**Current**: Created date DESC only  
**Required**: Model match > brand match > category match > text similarity

**Example**: "Samsung M14" query should rank exact M14 phones above other Samsung products

#### 5. Dynamic Filters ❌
**Current**: Hardcoded filters in frontend  
**Required**: Filters generated from category attributes

**Example**: Electronics category shows "Brand, Model, Storage"; Vehicles shows "Make, Model, Year"

#### 6. Synonym Support ❌
**Current**: No synonym dictionary  
**Required**: "mobile" ↔ "phone", "flat" ↔ "apartment", "plough" ↔ "plow"

#### 7. Ethiopian Localization ❌
**Current**: English-only search, English place names  
**Required**: Amharic, Afaan Oromo support, local terminology

---

## Proposed Architecture

### Phase 2: Category Hierarchy

**New Schema**:
```sql
CREATE TABLE category_hierarchy (
    id SERIAL PRIMARY KEY,
    parent_id INT REFERENCES category_hierarchy(id) ON DELETE RESTRICT,
    name VARCHAR(200) NOT NULL,
    name_am VARCHAR(200),  -- Amharic name
    name_om VARCHAR(200),  -- Afaan Oromo name
    slug VARCHAR(200) UNIQUE NOT NULL,
    level INT NOT NULL DEFAULT 0,  -- 0 = root, 1 = level 1, etc.
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Prevent circular references
    CONSTRAINT no_circular_reference CHECK (id != parent_id)
);

CREATE INDEX idx_category_hierarchy_parent ON category_hierarchy(parent_id);
CREATE INDEX idx_category_hierarchy_level ON category_hierarchy(level);
CREATE INDEX idx_category_hierarchy_slug ON category_hierarchy(slug);
```

**Migration Strategy**:
1. Create new `category_hierarchy` table
2. Migrate existing `categories` to root level
3. Parse existing `subcategory` and `item` text fields into hierarchy
4. Update `products.category_id` to point to leaf categories
5. Keep old `categories` table temporarily for rollback

### Phase 3: Attribute Definitions

**New Schema**:
```sql
CREATE TABLE attribute_definitions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    name_am VARCHAR(100),
    name_om VARCHAR(100),
    slug VARCHAR(100) UNIQUE NOT NULL,
    data_type VARCHAR(50) NOT NULL,  -- 'text', 'number', 'boolean', 'select', 'multiselect', 'range'
    is_required BOOLEAN DEFAULT false,
    is_searchable BOOLEAN DEFAULT true,
    is_filterable BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    validation_rules JSONB,  -- {"min": 0, "max": 1000, "pattern": "regex"}
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE category_attributes (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES category_hierarchy(id) ON DELETE CASCADE,
    attribute_id INT NOT NULL REFERENCES attribute_definitions(id) ON DELETE CASCADE,
    is_inherited BOOLEAN DEFAULT false,  -- Inherited from parent category
    UNIQUE(category_id, attribute_id)
);

CREATE TABLE product_attributes (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    attribute_id INT NOT NULL REFERENCES attribute_definitions(id) ON DELETE CASCADE,
    value_text TEXT,
    value_number NUMERIC,
    value_boolean BOOLEAN,
    value_json JSONB,  -- For multiselect, complex values
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_product_attributes_product ON product_attributes(product_id);
CREATE INDEX idx_product_attributes_attribute ON product_attributes(attribute_id);
CREATE INDEX idx_product_attributes_value_text ON product_attributes(value_text);
CREATE INDEX idx_product_attributes_value_number ON product_attributes(value_number);
```

### Phase 4: Query Understanding

**Components**:
1. Brand/Model Dictionary (database table)
2. Synonym Dictionary (database table)
3. Intent Detection (keyword patterns)
4. Query Normalization (lowercase, trim, common misspellings)

**New Schema**:
```sql
CREATE TABLE brand_dictionary (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL UNIQUE,
    aliases JSONB,  -- ["samsung", "samung", "samsng"]
    category_ids INT[],  -- Applicable categories
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE synonym_dictionary (
    id SERIAL PRIMARY KEY,
    term VARCHAR(200) NOT NULL,
    synonyms JSONB,  -- ["mobile", "phone", "cell phone"]
    context VARCHAR(100),  -- Optional: limit to category/domain
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE search_analytics (
    id BIGSERIAL PRIMARY KEY,
    query TEXT NOT NULL,
    normalized_query TEXT,
    result_count INT,
    category_id INT,
    filters_used JSONB,
    clicked_product_id INT,
    user_id INT,
    session_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_search_analytics_query ON search_analytics(normalized_query);
CREATE INDEX idx_search_analytics_created ON search_analytics(created_at);
```

### Phase 5: Enhanced Ranking

**Ranking Signals** (in order of priority):
1. Exact model match (highest)
2. Brand match
3. Category match (leaf > parent)
4. Attribute match
5. Query term frequency in title
6. Query term frequency in description
7. Listing freshness (newer = better)
8. Seller tier (premium > gold > normal)
9. Product quality signals (if available)

**Implementation**: PostgreSQL `ts_rank()` + custom scoring

### Phase 6: Dynamic Filtering

**Filter Generation Logic**:
```
1. GET category_id from user selection or query
2. FETCH applicable attributes FROM category_attributes WHERE category_id IN (selected + ancestors)
3. FOR EACH attribute:
     IF filterable:
         GENERATE filter UI component based on data_type
         FETCH available values FROM product_attributes WHERE attribute_id = ...
4. RENDER filters dynamically
```

---

## Implementation Plan

### Timeline Estimate

| Phase | Description | Estimated Time |
|-------|-------------|----------------|
| **Phase 1** | Audit & Baseline | 4-6 hours ✅ IN PROGRESS |
| **Phase 2** | Category Hierarchy | 8-12 hours |
| **Phase 3** | Attribute System | 10-14 hours |
| **Phase 4** | Query Understanding | 8-10 hours |
| **Phase 5** | Enhanced Ranking | 6-8 hours |
| **Phase 6** | Dynamic Filters | 10-12 hours |
| **Phase 7** | Frontend Integration | 12-16 hours |
| **Phase 8** | Admin Integration | 8-10 hours |
| **Phase 9** | Performance Optimization | 6-8 hours |
| **Phase 10** | Testing & Rollout | 8-12 hours |
| **Total** | | **80-108 hours** |

### Dependencies
- Database migration tools
- Existing admin permissions
- Frontend component library
- Testing data (real marketplace listings)

---

## Next Steps (Phase 1 Completion)

- [ ] Inspect frontend search pages (`/find`, `/rent`, etc.)
- [ ] Document current filter UI implementation
- [ ] Analyze current database indexes
- [ ] Measure baseline search performance
- [ ] Create representative test query collection
- [ ] Document current categories and subcategories
- [ ] Identify high-value Ethiopian search terms

---

## Risks and Mitigation

### Risk 1: Breaking Existing Search
**Mitigation**: Feature flags, gradual rollout, comprehensive regression tests

### Risk 2: Performance Degradation
**Mitigation**: Benchmark each phase, optimize indexes, use caching

### Risk 3: Data Migration Complexity
**Mitigation**: Safe migration scripts, rollback plan, data validation

### Risk 4: Admin Adoption
**Mitigation**: Training documentation, intuitive UI, backward compatibility

---

**Status**: Phase 1 - Audit in progress  
**Next Phase**: Category Hierarchy Design  
**Last Updated**: October 21, 2026

