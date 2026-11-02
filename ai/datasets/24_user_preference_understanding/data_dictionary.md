# Data Dictionary: User Preference Understanding (`24_user_preference_understanding`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `session_context` | `object` | No | Attribute representing session_context | - |
| `revealed_preference` | `string` | Yes | Attribute representing revealed_preference | - |
| `preference_scope` | `string` | Yes | Attribute representing preference_scope | ['temporary_session', 'ephemeral_turn'] |
| `adapted_response` | `string` | No | Attribute representing adapted_response | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
