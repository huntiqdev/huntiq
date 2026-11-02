#!/usr/bin/env python3
"""
Final expansion batch for governance, pipeline, and knowledge architecture docs.
"""

import os

final_docs = {
    "ai/data_pipeline/architecture.md": """# HUNTIQ Data Ingestion & Transformation Architecture

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Purpose:** End-to-end data pipeline architecture moving unnormalized external and internal records into validated, deduplicated, and leakage-proof training splits.

---

## 1. Pipeline Stages
```text
[External Sourced Batches / Internal Logs]
                    │
                    ▼
Stage 1: Ingestion & Provenance Registration
- Verifies SHA-256 checksum against source_schema.json.
- Registers legal custodian and commercial license in sources.json.
                    │
                    ▼
Stage 2: Validation & Quarantine Routing
- Executes jsonschema validation against module schema.json.
- Records failing validation are routed to quarantine/ with rejection codes.
                    │
                    ▼
Stage 3: Linguistic Normalization
- Unicode NFC normalization.
- Ge'ez homophone unification (ሀ/ሐ/ኀ, ሰ/ሠ, ጸ/ፀ, አ/ዐ).
- Afaan Oromoo Qubee punctuation and letter casing standardization.
- Tigrinya dialectal alignment.
                    │
                    ▼
Stage 4: Deduplication & Semantic Filtering
- Exact MD5 string deduplication.
- MinHash LSH (0.90 Jaccard threshold) near-duplicate removal.
- Synthetic outlier detection.
                    │
                    ▼
Stage 5: Deterministic Partitioning & Leakage Audit
- 80% Training (training/train.jsonl)
- 10% Validation (training/validation.jsonl)
- 10% Test (training/test.jsonl)
- Session-grouped hashing prevents customer leakage between splits.
```

---

## 2. Invariants & Security
- Zero silent drops: Every rejected record must have a logged rejection code in rejection_schema.json.
- Zero PII: 100% of telephone numbers and personal identifiers are cryptographically masked.
""",
    "ai/data_pipeline/language_processing.md": """# HUNTIQ Linguistic Processing Specification

**Languages:** English (`en`), Amharic (`am`), Afaan Oromoo (`or`), Tigrinya (`ti`).

---

## 1. Ethiopic (Ge'ez) Script Normalization Rules
1. **Homophone Standardization:**
   - Unified to first-order Fidel: ሀ, ሐ, ኀ $\\to$ ሀ
   - ሰ, ሠ $\\to$ ሰ
   - ጸ, ፀ $\\to$ ጸ
   - አ, ዐ $\\to$ አ
2. **Punctuation Normalization:**
   - Replace Ethiopic full stop (፠, ፡, ፣, ፤, ፥, ፦, ፧, ።) with canonical punctuation tokens.
   - Traditional word separator colon (፡) normalized to whitespace.
3. **Number Systems:**
   - Traditional Ethiopic numerals (፩, ፪, ፫, ፬, ፭, ፮, ፯, ፰, ፱, ፲, ፳, ፴, ፵, ፶, ፷, ፸, ፹, ፺, ፻, ፼) mapped to Arabic integer values.

---

## 2. Afaan Oromoo (Qubee) Normalization Rules
1. **Glottal Stop (Hudhaa) Handling:**
   - Standardize apostrophes (' or ’) to standard ASCII apostrophe `'`.
2. **Vowel and Consonant Lengthening:**
   - Preserve double vowels (aa, ee, ii, oo, uu) and double consonants (bb, dd, gg) essential for semantic distinction.
""",
    "ai/training/model_registry/model_approval_policy.md": """# HUNTIQ Model Registry & Production Approval Policy

---

## 1. Approval Gate Protocol
A candidate fine-tuned model checkpoint may only be marked `STATUS: APPROVED_PRODUCTION` in `model_registry.json` when the following conditions are satisfied:

1. **Frozen Benchmark Evaluation:** Evaluated on all 3 golden test suites (`domain_benchmark.jsonl`, `conversational_benchmark.jsonl`, `entity_extraction.jsonl`).
2. **Macro F1 Score Threshold:** Macro F1 $\\ge 0.92$ on intent classification and entity span extraction.
3. **Zero Security Regressions:** 100% defense against prompt injections in Module 18 evaluation suite.
4. **Latency Verification:** P95 response time $\\le 120\\text{ms}$ under 10 concurrent requests on target CPU inference instance.
5. **Dual Engineering Sign-Off:** Signed off by the AI Architect and Go Backend Lead.
""",
    "ai/knowledge/knowledge_architecture.md": """# HUNTIQ Grounded Knowledge Architecture

---

## 1. Architecture Overview
HUNTIQ's platform knowledge system bridges verified company documentation with real-time conversational retrieval.

```text
[Official Policy Markdown (ai/knowledge/)]
                    │
                    ▼
[Deterministic Chunking & Embedding Pipeline]
                    │
                    ▼
[Local Vector / Keyword Index (Ollama / In-Memory)]
                    │
                    ▼
[Context Injection via Strict Pass-Through Delimiters]
                    │
                    ▼
[Grounded Answer Generation (Module 03 & 30)]
```

## 2. Inviolable Truth Hierarchy
1. Live Database State (PostgreSQL via Go Backend): Sole authority on listing prices, stock, and status.
2. Official Policy Documents (`ai/knowledge/`): Sole authority on fees, terms, and rules.
3. Model Parametric Memory: Never used as an authority on facts.
""",
    "ai/knowledge/marketplace_entity_model.md": """# HUNTIQ Marketplace Entity Model

---

## 1. Relational Entity Schema
- **User:** `user_id`, `phone_number_hash`, `account_tier` (free, gold, enterprise), `verified_badge`, `created_at`.
- **Listing:** `listing_id`, `user_id`, `category_id`, `title`, `description`, `price_etb`, `negotiable`, `listing_type` (rent/sale), `location_id`, `status` (active, expired, sold).
- **Category:** `category_id`, `slug`, `parent_id`, `name_en`, `name_am`, `name_or`, `name_ti`, `attribute_schema`.
- **Location:** `location_id`, `region`, `city`, `subcity`, `woreda`, `neighborhood_landmarks`, `coordinates`.
- **Promotion:** `promotion_id`, `listing_id`, `tier` (featured_7d, top_30d), `paid_via` (telebirr, cbe_birr), `receipt_ref`.
""",
    "ai/knowledge/knowledge_validation.md": """# HUNTIQ Knowledge Validation & Review Procedures

---

## 1. Review Cadence
- **Monthly Operations Review:** Verify that phone support hours, hotline numbers, and office addresses are current.
- **Bi-annual Fee Review:** Review fee schedules, VAT rates, and Telebirr merchant codes against official banking agreements.
- **Immediate Invalidation:** Any policy change in the Go backend immediately triggers a revision of the corresponding document in `ai/knowledge/` and invalidates retrieval chunk caches.
""",
}

for path, content in final_docs.items():
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w") as f:
        f.write(content.strip() + "\n")
    print(f"Final expansion applied to: {path}")

print("Final expansion complete.")
