# HUNTIQ AI Security Boundary & Trust Architecture

## 1. Trust Boundaries

```
[ Public Internet / Untrusted Visitors ]
                  │
                  ▼ (Public HTTPS)
┌─────────────────────────────────────────────────┐
│ Go Backend API Gateway (Port 8080)             │
│ - Rate Limiter (15 req/min, burst 5)            │
│ - Content-Length limits (<= 1000 chars)         │
│ - Zero guest prompt persistence to DB           │
└───────────────────────┬─────────────────────────┘
                        │ (Private Loopback Network 127.0.0.1)
                        │ Header: X-HUNTIQ-AI-SECRET
                        ▼
┌─────────────────────────────────────────────────┐
│ Python AI Service (Port 8000)                   │
│ - Input Normalization & Prompt Injection Guard  │
│ - Zero direct DB credentials or SQL pool        │
│ - In-memory temporary bounded context only      │
└───────────────────────┬─────────────────────────┘
                        │ (Private Loopback Network 127.0.0.1)
                        │ Header: X-HUNTIQ-AI-SECRET
                        │ POST /api/internal/ai/tools/execute
                        ▼
┌─────────────────────────────────────────────────┐
│ Go Backend Authority (Layer 2 & 3 Validation)  │
│ - Tool name whitelist enforcement               │
│ - Result limit capping (<= 10 items)            │
│ - Price bound sanity checks                     │
│ - Authorization & Listing Visibility filter     │
└───────────────────────┬─────────────────────────┘
                        │
                        ▼ (PostgreSQL Connection Pool)
┌─────────────────────────────────────────────────┐
│ PostgreSQL 18.4 Database                        │
│ - Only approved, active, un-deleted listings    │
└─────────────────────────────────────────────────┘
```

---

## 2. Four Independent Validation Layers

1. **Layer 1: AI Interpretation Validation (Python)**
   - Validates that user input is normalized and maps to recognized intents.
   - Rejects malicious prompt injections trying to override system behavior.
   - Decides whether a clarification is necessary before invoking tools.

2. **Layer 2: API & Tool-Request Validation (Go)**
   - Re-validates every incoming tool call independently at the Go API boundary.
   - Enforces valid tool names (`search_marketplace_listings`, `get_listing_details`, `compare_listings`, `get_marketplace_stats`).
   - Caps result limits (1 to 10) to prevent denial of service.
   - Rejects negative price parameters.

3. **Layer 3: Business Rules & Authorization (Go)**
   - Enforces user membership tiers, seller verification badges, and listing approval states.
   - Filters out hidden, suspended, or unapproved records.
   - Go owns this layer exclusively; Python has zero ability to bypass authorization.

4. **Layer 4: Result Grounding (Python)**
   - Takes Go's verified database results and formats them into natural summaries.
   - Never fabricates listings, prices, contact phone numbers, or inventory counts.
