# Data Dictionary: Search Refinement and Progressive Filtering (`27_search_refinement_progressive_filtering`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `base_filter` | `object` | Yes | Attribute representing base_filter | - |
| `refinement_input` | `string` | Yes | Attribute representing refinement_input | - |
| `filter_mutation` | `string` | Yes | Attribute representing filter_mutation | ['narrow', 'widen', 'replace', 'remove'] |
| `resulting_filter` | `object` | Yes | Attribute representing resulting_filter | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
