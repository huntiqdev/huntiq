# Data Dictionary: Assistant Identity and Behavior (`02_identity_assistant_behavior`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `query` | `string` | Yes | Attribute representing query | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `inquiry_type` | `string` | Yes | Attribute representing inquiry_type | ['name', 'purpose', 'capabilities', 'limitations', 'guest_access', 'pricing'] |
| `authorized_response` | `string` | Yes | Attribute representing authorized_response | - |
| `is_honest_limitation` | `boolean` | No | Attribute representing is_honest_limitation | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
