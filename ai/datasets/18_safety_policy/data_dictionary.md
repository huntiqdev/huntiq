# Data Dictionary: Safety, Security and Platform Policy (`18_safety_policy`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
| `id` | `string` | Yes | Attribute representing id | - |
| `adversarial_prompt` | `string` | Yes | Attribute representing adversarial_prompt | - |
| `threat_category` | `string` | Yes | Attribute representing threat_category | ['prompt_injection', 'pii_leak', 'prohibited_content', 'auth_bypass', 'data_poisoning'] |
| `safe_refusal_response` | `string` | Yes | Attribute representing safe_refusal_response | - |
| `policy_violation_code` | `string` | No | Attribute representing policy_violation_code | - |

## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
