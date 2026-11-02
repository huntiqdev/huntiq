# HUNTIQ Data Ingestion & Transformation Architecture

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
