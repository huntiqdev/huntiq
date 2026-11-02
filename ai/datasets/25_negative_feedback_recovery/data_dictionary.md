# Data Dictionary: Negative Feedback and Recovery (`25_negative_feedback_recovery`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `prior_assistant_response` | `string` | No | Attribute representing prior_assistant_response | - |
| `user_rejection` | `string` | Yes | Attribute representing user_rejection | - |
| `failure_diagnosis` | `string` | Yes | Attribute representing failure_diagnosis | - |
| `recovery_turn` | `string` | Yes | Attribute representing recovery_turn | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
