# Data Dictionary: Constraint Reasoning and Conflict Resolution (`23_constraint_reasoning_conflict_resolution`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `conflicting_constraints` | `object` | Yes | Attribute representing conflicting_constraints | - |
| `conflict_type` | `string` | Yes | Attribute representing conflict_type | - |
| `resolution_strategy` | `string` | No | Attribute representing resolution_strategy | ['relax_price', 'widen_location', 'ask_user_priority'] |
| `suggested_compromise` | `string` | Yes | Attribute representing suggested_compromise | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
