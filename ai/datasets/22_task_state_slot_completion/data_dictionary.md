# Data Dictionary: Task State and Slot Completion (`22_task_state_slot_completion`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `task_id` | `string` | Yes | Attribute representing task_id | - |
| `goal` | `string` | Yes | Attribute representing goal | - |
| `required_slots` | `array` | Yes | Attribute representing required_slots | - |
| `filled_slots` | `object` | Yes | Attribute representing filled_slots | - |
| `is_complete` | `boolean` | Yes | Attribute representing is_complete | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
