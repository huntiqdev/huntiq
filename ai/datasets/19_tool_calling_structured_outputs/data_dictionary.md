# Data Dictionary: Tool Calling and Structured Outputs (`19_tool_calling_structured_outputs`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `user_intent` | `string` | Yes | Attribute representing user_intent | - |
| `tool_name` | `string` | Yes | Attribute representing tool_name | ['search_marketplace_listings', 'get_listing_details', 'compare_listings', 'get_marketplace_stats'] |
| `arguments` | `object` | Yes | Attribute representing arguments | - |
| `schema_valid` | `boolean` | Yes | Attribute representing schema_valid | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
