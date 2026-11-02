# HUNTIQ AI — PHASE 1.5: MISSING CAPABILITIES & GAP ANALYSIS

**Project:** HUNTIQ (SellBuy)  
**Audit Phase:** Phase 1.5 — Independent Architecture & Dataset Readiness Audit  
**Date:** 2026-11-02  
**Auditor Status:** Senior AI Architecture & Data Engineering Independent Review  

---

## 1. Executive Summary

During the Phase 1.5 independent audit of the 30 dataset modules, schemas, manifests, and intelligence contracts, the core marketplace foundation was verified as sound, strictly bounded to HUNTIQ's Ethiopian marketplace domain, and aligned with the sovereign Go backend.

However, an in-depth operational review against Ethiopian marketplace realities identified **5 specific domain gaps and architectural enhancements**. In accordance with the mandate to avoid unnecessary module proliferation, **no new dataset modules are added**. Instead, targeted schema extensions, vocabulary additions, and contract enhancements are recommended for existing modules.

---

## 2. Identified Capabilities & Gap Assessment

### Gap 1: Telebirr / CBE Birr / Chapa Payment Verification Dialogues
- **Marketplace Reality:** In Ethiopia, platform transactions, Gold membership upgrades, and listing promotions rely on Telebirr, CBE Birr, Chapa, and Awash Bank. Users frequently send messages such as:
  - *"I transferred 500 Birr via Telebirr with transaction code 'FT261102ABC', why is my listing not yet promoted?"*
  - *"በቴሌብር 300 ብር ከፍያለሁ ግን ማስታወቂያዬ አልታየም"*
- **Current Limitation in Phase 1:** Module 03 ([03_huntiq_platform_knowledge](file:///home/poven/GO%20PROJECT/sellbuy/ai/datasets/03_huntiq_platform_knowledge)) covers general membership policies, but lacks a structured entity schema for payment transaction codes, payment methods, or payment dispute reconciliation flows.
- **Risk:** The AI might hallucinate payment confirmation or disclose internal transaction details without backend verification.
- **Recommendation:**
  1. **Extend Existing Dataset:** Add payment verification dialogue pairs to Module 03 and fraud prevention cases to Module 18 ([18_safety_policy](file:///home/poven/GO%20PROJECT/sellbuy/ai/datasets/18_safety_policy)).
  2. **Add Fields to Schema:** In `03_huntiq_platform_knowledge/schema.json` and `07_entity_attribute_extraction/schema.json`, add fields:
     - `payment_method`: `enum ["telebirr", "cbe_birr", "chapa", "awash_birr", "bank_transfer"]`
     - `transaction_reference`: `string (regex pattern validated)`
     - `reconciliation_status`: `enum ["pending_verification", "verified", "rejected"]`
  3. **Add Evaluation Benchmark:** Add adversarial payment fraud test cases to `18_safety_policy/evaluation/edge_cases.jsonl` ensuring the assistant never claims payment success without Go backend confirmation.

---

### Gap 2: Local Delivery, Handover Spots & Intercity Logistics
- **Marketplace Reality:** Ethiopian marketplace transactions frequently require negotiating delivery terms:
  - Within Addis Ababa: Local couriers (Tikur Anbessa, Megenagna, Piassa, Bole Edna Mall handover points) or ride-hailing package delivery (Feres, Ride).
  - Regional Intercity: Sending items (electronics, documents, spare parts) via cross-country bus lines (Selam Bus, Sky Bus, Golden Bus, Zemen Bus) from Addis Ababa to regional hubs like Hawassa, Bahir Dar, Adama, or Mekelle.
- **Current Limitation in Phase 1:** Module 07 extracts product specs and Module 10 extracts locations, but neither captures delivery options, handover spots, or logistics negotiation terms.
- **Recommendation:**
  1. **Extend Existing Schema:** In `07_entity_attribute_extraction/schema.json`, add an optional `logistics` object:
     ```json
     "logistics": {
       "type": "object",
       "properties": {
         "delivery_available": { "type": "boolean" },
         "handover_preference": { "type": "string", "enum": ["in_person_pickup", "addis_courier", "regional_bus_cargo", "seller_delivers"] },
         "designated_landmark": { "type": "string" }
       }
     }
     ```
  2. **Extend Slot Completion:** Update `22_task_state_slot_completion` to track `handover_preference` as an optional dialogue slot.

---

### Gap 3: Ethiopian Item Condition & Grading Standards
- **Marketplace Reality:** In Ethiopian electronics and vehicle markets, items are described with culturally specific condition qualifiers:
  - **"Boxed / Sealed" (ቦክስ ያልተከፈተ):** Brand new in original packaging.
  - **"Dubai Used" (ከዱባይ የመጣ):** High-quality imported second-hand item, perceived as superior to locally worn goods.
  - **"Refurbished" (የታደሰ):** Repaired or reconditioned with replacement parts.
  - **"Locally Used" (እዚህ አገር ያገለገለ):** Used within Ethiopia, generally priced lower.
- **Current Limitation in Phase 1:** Module 06 and Module 07 use a generic binary condition field (`"new"` vs `"used"`).
- **Recommendation:**
  1. **Extend Schema:** In `06_category_understanding/schema.json` and `07_entity_attribute_extraction/schema.json`, expand the condition enum to:
     ```json
     "item_condition": {
       "type": "string",
       "enum": ["brand_new_sealed", "dubai_used_import", "refurbished", "locally_used", "parts_only"]
     }
     ```
  2. **Update Multilingual Glossaries:** Add localized condition terminology across `ai/languages/marketplace_glossary/` (`en.json`, `am.json`, `or.json`, `ti.json`).

---

### Gap 4: Ge'ez Numeral Parsing for Traditional Prices & Measurements
- **Marketplace Reality:** While modern Arabic numerals (1, 2, 3...) predominate, traditional Ge'ez numerals (፩, ፪, ፫, ፲, ፻, ፼) and traditional land measurement units ("ጋሻ" / Gasha, "ካሬ" / Kare, "ካሬ ሜትር" / Square Meters) appear in formal property, land, and agricultural contracts.
- **Current Limitation in Phase 1:** Module 09 and Module 17 handle standard Arabic numeric parsing, but lack Ge'ez traditional numeral token translation.
- **Recommendation:**
  1. **Add Normalization Rule:** Add a Ge'ez-to-Arabic numeral converter utility to `ai/app/core/normalizer.py` and document the contract in `ai/intelligence/language_contract.md`.
  2. **Add Test Cases:** Add Ge'ez numeral price expressions (`፳ ሺህ ብር` $\to$ `20,000 ETB`) to `09_price_budget_understanding/evaluation/edge_cases.jsonl`.

---

### Gap 5: Architectural Documentation Completeness (Stub Expansion)
- **Finding:** Phase 1 created the directory hierarchy and all schemas, but left several high-level governance, coordination, and pipeline markdown files with minimal stub text (1–3 lines).
- **Impact:** While schemas and manifests are fully functional, human developers and auditors require explicit narrative specifications.
- **Recommendation:**
  1. **Flesh Out Critical Architecture Files:** Expand [dataset_dependency_map.md](file:///home/poven/GO%20PROJECT/sellbuy/ai/datasets/dataset_dependency_map.md), [dataset_relationships.md](file:///home/poven/GO%20PROJECT/sellbuy/ai/datasets/dataset_relationships.md), and coordination contracts with comprehensive tables and ASCII diagrams.
  2. **Audit Status Classification:** Classify overall Phase 1 readiness as **READY WITH DOCUMENTED LIMITATIONS** pending narrative documentation expansion across future data ingestion phases.

---

## 3. Summary of Gap Resolution Matrix

| Gap Identified | Recommended Action | New Module Needed? | Target Module(s) | Implementation Phase |
| :--- | :--- | :---: | :--- | :---: |
| **Payment Verification Dialogue** | Add fields to schema & safety checks | **NO** | `03_huntiq_platform_knowledge`, `18_safety_policy` | Phase 2 / Phase 3 |
| **Delivery & Logistics Negotiation** | Add `logistics` object to schema | **NO** | `07_entity_attribute_extraction`, `22_task_state_slot_completion` | Phase 3 |
| **Ethiopian Condition Grading** | Expand condition enum to 5 tiers | **NO** | `06_category_understanding`, `07_entity_attribute_extraction` | Phase 2 / Phase 4 |
| **Ge'ez Traditional Numerals** | Add normalizer rule & test cases | **NO** | `09_price_budget_understanding`, `17_comparison_calculations` | Phase 4 |
| **Narrative Documentation Stubs** | Expand architectural markdown | **NO** | `ai/datasets/`, `ai/intelligence/`, `ai/governance/` | Phase 1.5 / Ongoing |
