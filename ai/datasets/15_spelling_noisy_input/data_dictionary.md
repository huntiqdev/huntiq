# Data Dictionary: Spelling and Noisy Input (`15_spelling_noisy_input`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `noisy_input` | `string` | Yes | Attribute representing noisy_input | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `normalized_output` | `string` | Yes | Attribute representing normalized_output | - |
| `noise_type` | `string` | No | Attribute representing noise_type | ['typo', 'phonetic', 'script_mix', 'abbreviation', 'keyboard_slip'] |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
