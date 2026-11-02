# Data Dictionary: Domain Classification (`04_domain_classification`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `text` | `string` | Yes | Attribute representing text | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `domain_label` | `string` | Yes | Attribute representing domain_label | ['marketplace', 'platform_support', 'relevant_calculation', 'out_of_domain', 'unsafe'] |
| `confidence_threshold` | `number` | No | Attribute representing confidence_threshold | - |
| `redirection_needed` | `boolean` | No | Attribute representing redirection_needed | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
