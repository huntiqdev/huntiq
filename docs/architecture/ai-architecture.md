# HUNTIQ AI Architecture Specification

## 1. System Overview & Four-Tier Separation

HUNTIQ establishes four distinct, strictly bounded layers of responsibility:

```
               HUNTIQ FRONTEND
            Next.js 16 + TypeScript
                       |
                       |  POST /api/ai/chat (JSON)
                       v
               GO BACKEND API
            Fiber v2 (Port 8080)
            [Auth, Rate Limit, Proxy]
                       |
                       |  POST /api/v1/chat (X-HUNTIQ-AI-SECRET)
                       v
               PYTHON AI SERVICE
             FastAPI (Port 8000)
       [Natural Language, Intent, Slots]
                       |
                       |  POST /api/internal/ai/tools/execute
                       v
               GO BACKEND API
         [Layer 2 & 3 Validation, Tool Execution]
                       |
                       v
                 POSTGRESQL 18.4
               [Authoritative DB]
```

### Trust & Authority Matrix

| Layer | Technology | Primary Authority | Invariants |
|---|---|---|---|
| **Presentation** | Next.js 16 (TypeScript) | User interaction & display | Only customer-facing interface. Never talks directly to PostgreSQL or internal Python ports. |
| **Marketplace Authority** | Go (Fiber v2) | Business rules, security, DB access | Sole authority for DB mutations, payments, user auth, listing visibility, rate limiting, and tool validation. |
| **AI Intelligence** | Python (FastAPI) | Language understanding & reasoning | Sole authority for intent classification, slot extraction, context tracking, and response composition. Zero DB write access. |
| **Data Authority** | PostgreSQL 18.4 | Persistent storage | Authoritative truth. Guest chats are never permanently stored. |

---

## 2. Conversation Flow Architecture

```
User Message
    |
    v
Language Detection & Normalization (Python)
    |
    v
In-Memory Session Context Lookup (Python)
    |
    v
Intent Classification & Disambiguation (Python)
    |
    +---> Out-of-Domain Guard ---> Polite Redirection
    |
    +---> Platform FAQ Check ------> Structured Guidance
    |
    +---> Ambiguous Query ---------> Clarification Request
    |
    +---> Conversational Query ----> Direct Greeting / Identity Response
    |
    +---> Marketplace Search ------>
             1. Python builds typed tool arguments
             2. Python sends HTTP tool call to Go (:8080/api/internal/ai/tools/execute)
             3. Go executes Layer 2 & 3 validation
             4. Go queries PostgreSQL for active, approved listings
             5. Go returns structured listings to Python
             6. Python grounds and summarizes results without hallucination
             7. Python returns unified response to Go -> Next.js
```

---

## 3. Graceful Degradation & Fallback Guarantee

If the Python AI service is starting up, temporarily unreachable, or experiencing high load:
1. The Go backend's `AIChatService` catches the error.
2. Go immediately executes its internal deterministic fallback classifier.
3. The marketplace visitor receives a truthful response without 500 errors or page crashes.
