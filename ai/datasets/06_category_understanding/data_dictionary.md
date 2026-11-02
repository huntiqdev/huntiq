# Data Dictionary: Marketplace Category Understanding (`06_category_understanding`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `raw_term` | `string` | Yes | Attribute representing raw_term | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `canonical_slug` | `string` | Yes | Attribute representing canonical_slug | - |
| `category_level` | `integer` | No | Attribute representing category_level | [1, 2, 3] |
| `parent_slug` | `string` | No | Attribute representing parent_slug | - |
| `is_ambiguous` | `boolean` | No | Attribute representing is_ambiguous | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
