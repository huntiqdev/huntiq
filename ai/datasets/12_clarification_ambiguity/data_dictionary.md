# Data Dictionary: Clarification and Ambiguity (`12_clarification_ambiguity`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `ambiguous_input` | `string` | Yes | Attribute representing ambiguous_input | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `missing_slot` | `string` | Yes | Attribute representing missing_slot | ['listing_type', 'location', 'category', 'budget'] |
| `clarification_question` | `string` | Yes | Attribute representing clarification_question | - |
| `can_proceed_with_default` | `boolean` | No | Attribute representing can_proceed_with_default | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
