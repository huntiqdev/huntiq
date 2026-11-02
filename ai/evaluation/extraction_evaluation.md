# Named Entity & Attribute Extraction Evaluation

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Measures token-level and span-level precision, recall, and F1 score for extracted marketplace entities across all four supported languages.

---

## 2. Detailed Technical & Operational Rules
- Target Entities: Brand, Model, Storage, Color, Condition, Price, Bedrooms, Bathrooms, Subcity, Neighborhood, Transmission, Mileage.
- Multi-Script Consistency: Amharic Ge'ez spans and Latin spans must achieve equal extraction accuracy.
- F1-Score Target: Macro F1 >= 0.90 across all primary entity categories.

---

## 3. Enforcement & Verification Protocols
Validated using  and .

---

## 4. Cross-System References
- Sovereign Go Authority: `backend/internal/`
- Python AI Inference Runtime: `ai/app/`
- Authoritative Marketplace Database: PostgreSQL 16
