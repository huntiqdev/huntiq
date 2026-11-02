# Data Dictionary: Marketplace Comparison and Calculations (`17_comparison_calculations`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `query` | `string` | Yes | Attribute representing query | - |
| `items_compared` | `array` | Yes | Attribute representing items_compared | - |
| `calculation_formula` | `string` | No | Attribute representing calculation_formula | - |
| `deterministic_result` | `string` | Yes | Attribute representing deterministic_result | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
