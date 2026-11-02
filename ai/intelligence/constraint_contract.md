# Constraint Reasoning & Conflict Resolution Contract

---

## 1. Conflict Invariants
- Mathematical validation: `min_price` must be $\le$ `max_price`.
- Property constraints: `bedrooms` must be $\ge 1$.
- If user supplies contradictory instructions, priority is given to the most recent statement, with a brief confirmation.
