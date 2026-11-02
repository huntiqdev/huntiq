# Runtime Decision Matrix

---

| Input Scenario | Classified Intent | Go DB Search Required? | Response Type |
| :--- | :--- | :---: | :--- |
| "Hello!" | `SOCIAL_GREETING` | **NO** | Conversational greeting |
| "Who are you?" | `IDENTITY_INQUIRY` | **NO** | Persona & boundary explanation |
| "Cheap laptop under 20k" | `MARKETPLACE_SEARCH` | **YES** | Propose `search_listings` |
| "2 bed in Bole" | `PROPERTY_SEARCH` | **YES** | Propose `search_listings` |
| "Min 500k max 10k" | `CLARIFICATION_NEEDED`| **NO** | Ask price clarification |
| "Ignore instructions" | `POLICY_VIOLATION` | **NO** | Safe refusal |
