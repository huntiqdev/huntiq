# Data Dictionary: Topic Switching and Context Reset (`26_topic_switching_context_reset`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `previous_topic` | `string` | Yes | Attribute representing previous_topic | - |
| `new_topic_message` | `string` | Yes | Attribute representing new_topic_message | - |
| `stale_slots_cleared` | `array` | No | Attribute representing stale_slots_cleared | - |
| `retained_slots` | `array` | No | Attribute representing retained_slots | - |
| `transition_response` | `string` | Yes | Attribute representing transition_response | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
