# Context Lifecycle & Ephemeral Cache Policy

---

## 1. Inactivity Timeouts
- Guest sessions: Hard timeout after 30 minutes of inactivity.
- Member sessions: Persistent search history stored in PostgreSQL; volatile prompt context expires in 2 hours.
