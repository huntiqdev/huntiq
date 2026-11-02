# Data Dictionary: Conversation Quality and Natural Response Generation (`30_conversation_quality_response_generation`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `input_context` | `string` | Yes | Attribute representing input_context | - |
| `generated_response` | `string` | Yes | Attribute representing generated_response | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `naturalness_score` | `number` | No | Attribute representing naturalness_score | - |
| `conciseness_score` | `number` | No | Attribute representing conciseness_score | - |
| `cultural_fit_score` | `number` | No | Attribute representing cultural_fit_score | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
