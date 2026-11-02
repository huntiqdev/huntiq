# Data Dictionary: Multilingual Understanding (`14_multilingual_understanding`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `semantic_id` | `string` | Yes | Attribute representing semantic_id | - |
| `english` | `string` | Yes | Attribute representing english | - |
| `amharic` | `string` | Yes | Attribute representing amharic | - |
| `oromo` | `string` | Yes | Attribute representing oromo | - |
| `tigrinya` | `string` | Yes | Attribute representing tigrinya | - |
| `canonical_intent` | `string` | Yes | Attribute representing canonical_intent | - |
| `review_status` | `string` | No | Attribute representing review_status | ['human_verified', 'pending_review'] |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
