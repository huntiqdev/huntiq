# Data Dictionary: Ethiopian Location Understanding (`10_location_understanding`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `raw_location` | `string` | Yes | Attribute representing raw_location | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `canonical_city` | `string` | Yes | Attribute representing canonical_city | - |
| `subcity` | `string` | No | Attribute representing subcity | - |
| `region` | `string` | No | Attribute representing region | - |
| `transliteration_variants` | `array` | No | Attribute representing transliteration_variants | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
