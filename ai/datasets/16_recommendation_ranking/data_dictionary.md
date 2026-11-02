# Data Dictionary: Recommendation and Ranking (`16_recommendation_ranking`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `user_constraints` | `object` | Yes | Attribute representing user_constraints | - |
| `candidate_item_ids` | `array` | Yes | Attribute representing candidate_item_ids | - |
| `ranking_order` | `array` | Yes | Attribute representing ranking_order | - |
| `tradeoff_explanation` | `string` | No | Attribute representing tradeoff_explanation | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
