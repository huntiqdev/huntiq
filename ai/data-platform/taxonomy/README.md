# HUNTIQ Taxonomy ML Integration

This directory contains deterministic schemas and tools for consuming authoritative taxonomy snapshots from the Go backend into AI datasets.

## Governance & Sovereignty Principles

1. **Backend Authority**: Canonical taxonomy categories, slugs, icon keys, and parent-child hierarchies are authored and managed strictly within the Go backend (`backend/internal/model/category.go` and database migrations).
2. **Read-Only / Export Only**: AI workflows NEVER mutate, insert, or reassign canonical category IDs.
3. **Deterministic Exports**: When datasets require category labels, they must link against an exported taxonomy manifest matching the schema in `schemas/taxonomy_export.schema.json`.

## Directory Purposes

- `schemas/`: Machine-readable schemas defining taxonomy export contracts.
- `sources/`: Authoritative category metadata dumps exported from PostgreSQL (git-ignored for live data; synthetic fixtures for testing).
- `generated/`: Flattened or indexed taxonomy structures prepared for model training (e.g., label encoders).
- `manifests/`: Cryptographic checksums and version records for exported taxonomy snapshots.
- `reports/`: Audit reports detailing tree depth, leaf count, and validation checks.
