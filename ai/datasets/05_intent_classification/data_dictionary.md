# Data Dictionary: Intent Classification (`05_intent_classification`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `query` | `string` | Yes | Attribute representing query | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `primary_intent` | `string` | Yes | Attribute representing primary_intent | - |
| `secondary_intents` | `array` | No | Attribute representing secondary_intents | - |
| `is_multi_intent` | `boolean` | No | Attribute representing is_multi_intent | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
