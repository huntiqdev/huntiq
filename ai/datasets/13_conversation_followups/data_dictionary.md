# Data Dictionary: Conversation Follow-ups (`13_conversation_followups`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `history` | `array` | Yes | Attribute representing history | - |
| `followup_message` | `string` | Yes | Attribute representing followup_message | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `inferred_action` | `string` | Yes | Attribute representing inferred_action | - |
| `resolved_parameters` | `object` | No | Attribute representing resolved_parameters | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
