# Data Dictionary: Service Failure and Error Conversations (`29_service_failure_error_conversations`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `error_condition` | `string` | Yes | Attribute representing error_condition | ['api_timeout', 'db_unavailable', 'rate_limit_exceeded', 'network_split'] |
| `user_query` | `string` | Yes | Attribute representing user_query | - |
| `graceful_response` | `string` | Yes | Attribute representing graceful_response | - |
| `allow_retry` | `boolean` | Yes | Attribute representing allow_retry | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
