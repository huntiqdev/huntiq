# Data Dictionary: General Greetings and Social Conversation (`01_general_greetings`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `utterance` | `string` | Yes | Attribute representing utterance | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `greeting_type` | `string` | Yes | Attribute representing greeting_type | ['hello', 'farewell', 'gratitude', 'apology', 'casual'] |
| `has_commercial_intent` | `boolean` | Yes | Attribute representing has_commercial_intent | - |
| `response` | `string` | No | Attribute representing response | - |
| `mixed_languages` | `boolean` | No | Attribute representing mixed_languages | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
