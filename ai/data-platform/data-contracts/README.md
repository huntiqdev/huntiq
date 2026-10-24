# HUNTIQ Data Contracts

This directory defines explicit, machine-readable contracts between the authoritative Go backend data exports, the data-platform training pipelines, and the internal AI serving runtime.

## Directory Structure

```text
data-contracts/
├── README.md                           # This document
└── schemas/
    ├── manifest_schema.json            # JSONSchema for versioned dataset manifests
    ├── query_intent_contract.json      # Specification for search query intent classification
    └── category_suggestion_contract.json # Specification for listing category classification
```

## Contract Invariants

1. **Schema Authority**: The Go backend remains the sovereign authority for marketplace data models (listings, categories, users). Data contracts translate authoritative snapshots into immutable training and evaluation formats.
2. **Backward Compatibility**: Any change to contract schemas must increment the semantic version (`major.minor.patch`). Breaking changes require approval and parallel transition periods.
3. **Data Scrubbing**: Datasets created under these contracts MUST be stripped of personal identifiers (PII), credentials, phone numbers, payment tokens, and authentication cookies.
4. **Deterministic Validation**: Data pipelines must validate incoming data batches against these contracts using `data_quality.validator` before any processing or split generation.
