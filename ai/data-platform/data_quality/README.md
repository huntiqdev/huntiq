# HUNTIQ Data Quality & Governance Engine

This package provides deterministic validation for datasets, manifests, and data splits before they are admitted into ML pipelines or training runs.

## Key Capabilities

1. **Manifest Conformance**: Validates dataset manifest structures against `data-contracts/schemas/manifest_schema.json`.
2. **PII Detection**: Enforces zero-tolerance policies on unredacted personally identifiable information, including:
   - Ethiopian phone numbers (`+2519...`, `09...`, `07...`)
   - Email addresses
   - Authentication tokens, passwords, and API keys
3. **Duplicate Detection**: Identifies exact text duplicates and duplicated sample identifiers.
4. **Data Leakage Prevention**: Detects samples appearing in both training and test/validation evaluation splits, preventing inflated evaluation metrics.
5. **Dry-Run Validation**: Produces structured JSON audit reports without mutating underlying files.
