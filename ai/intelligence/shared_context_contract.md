# Shared Context Contract (Session State Schema)

---

## 1. Context Object Structure
```json
{
  "session_id": "string (UUID v4)",
  "user_role": "guest | authenticated_member",
  "locale": "en | am | or | ti",
  "active_category": "string | null",
  "active_slots": {
    "query": "string | null",
    "listing_type": "rent | sale | null",
    "min_price": "number | null",
    "max_price": "number | null",
    "location": "string | null",
    "attributes": "object"
  },
  "turn_count": "integer",
  "last_updated_epoch": "integer"
}
```
