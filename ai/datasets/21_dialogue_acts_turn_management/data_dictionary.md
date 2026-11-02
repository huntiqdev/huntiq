# Data Dictionary: Dialogue Acts and Turn Management (`21_dialogue_acts_turn_management`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `speaker` | `string` | Yes | Attribute representing speaker | ['user', 'assistant'] |
| `utterance` | `string` | Yes | Attribute representing utterance | - |
| `dialogue_act` | `string` | Yes | Attribute representing dialogue_act | ['inform', 'request', 'confirm', 'reject', 'greet', 'thank', 'clarify', 'farewell'] |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
