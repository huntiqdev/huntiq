# Data Dictionary: Response Grounding and Quality (`20_grounding_response_quality`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `query` | `string` | Yes | Attribute representing query | - |
| `db_results_count` | `integer` | Yes | Attribute representing db_results_count | - |
| `grounded_response` | `string` | Yes | Attribute representing grounded_response | - |
| `zero_hallucination_verified` | `boolean` | Yes | Attribute representing zero_hallucination_verified | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
