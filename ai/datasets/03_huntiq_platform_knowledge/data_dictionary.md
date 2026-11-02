# Data Dictionary: HUNTIQ Platform Knowledge (`03_huntiq_platform_knowledge`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `topic` | `string` | Yes | Attribute representing topic | ['products', 'membership', 'listing', 'account', 'payment', 'privacy', 'security', 'support'] |
| `query` | `string` | Yes | Attribute representing query | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `authoritative_answer` | `string` | Yes | Attribute representing authoritative_answer | - |
| `source_doc` | `string` | Yes | Attribute representing source_doc | - |
| `effective_date` | `string` | No | Attribute representing effective_date | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
