# Data Dictionary: Confidence and Uncertainty (`28_confidence_uncertainty`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `ambiguous_input` | `string` | Yes | Attribute representing ambiguous_input | - |
| `confidence_score` | `number` | Yes | Attribute representing confidence_score | - |
| `abstention_threshold` | `number` | No | Attribute representing abstention_threshold | - |
| `should_abstain` | `boolean` | Yes | Attribute representing should_abstain | - |
| `uncertainty_response` | `string` | Yes | Attribute representing uncertainty_response | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
