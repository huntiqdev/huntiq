# Data Dictionary: Price and Budget Understanding (`09_price_budget_understanding`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `expression` | `string` | Yes | Attribute representing expression | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `extracted_min` | `number` | No | Attribute representing extracted_min | - |
| `extracted_max` | `number` | No | Attribute representing extracted_max | - |
| `currency` | `string` | Yes | Attribute representing currency | - |
| `is_negotiable` | `boolean` | No | Attribute representing is_negotiable | - |
| `period` | `string` | No | Attribute representing period | ['one_time', 'daily', 'monthly', 'yearly'] |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
