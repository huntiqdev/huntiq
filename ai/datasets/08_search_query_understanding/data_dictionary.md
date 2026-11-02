# Data Dictionary: Search Query Understanding (`08_search_query_understanding`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `raw_query` | `string` | Yes | Attribute representing raw_query | - |
| `locale` | `string` | Yes | Attribute representing locale | ['en', 'am', 'or', 'ti'] |
| `search_term` | `string` | Yes | Attribute representing search_term | - |
| `category_slug` | `string` | No | Attribute representing category_slug | - |
| `min_price` | `number` | No | Attribute representing min_price | - |
| `max_price` | `number` | No | Attribute representing max_price | - |
| `location` | `string` | No | Attribute representing location | - |
| `condition` | `string` | No | Attribute representing condition | ['new', 'used', 'any'] |
| `listing_type` | `string` | No | Attribute representing listing_type | ['sale', 'rent', 'any'] |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
