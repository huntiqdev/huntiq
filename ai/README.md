# HUNTIQ AI — Sovereign Intelligence Platform

HUNTIQ AI is the dedicated Python intelligence authority for Ethiopia's universal marketplace. It provides natural language understanding, multilingual intent classification, slot extraction, context tracking, and response composition.

---

## 1. System Responsibilities

- **Sole Authority for AI Reasoning:** Intent classification, multi-intent extraction, entity slot filling, ambiguity clarification, and zero-fabrication response composition.
- **Languages Supported:** English (`en`), Amharic (`am`), Afaan Oromoo (`om`/`or`), and Tigrinya (`ti`).
- **Zero Direct Database Access:** Python holds no PostgreSQL credentials or connection pools. Marketplace data is retrieved strictly via authorized HTTP tool calls to the Go backend (`POST /api/internal/ai/tools/execute`) authenticated via `X-HUNTIQ-AI-SECRET`.
- **Zero Persistent Guest Storage:** Visitor conversations are kept in bounded, thread-safe in-memory sliding sessions (20-minute TTL, max 6 turns).

---

## 2. Architecture & Directory Layout

```
ai/
├── app/
│   ├── main.py                  # FastAPI server entry point
│   ├── config.py                # Service configuration (ports, timeouts, secrets)
│   ├── api/
│   │   ├── routes/
│   │   │   └── chat.py          # POST /api/v1/chat endpoint
│   │   └── schemas.py           # Pydantic request/response & tool schemas
│   ├── core/
│   │   ├── orchestrator.py      # Conversation pipeline runner
│   │   ├── intent.py            # Multilingual intent classification
│   │   ├── entity_extractor.py  # Slot & entity extraction (prices, locations, condition)
│   │   ├── policy.py            # Platform FAQ check & out-of-domain guard
│   │   ├── clarification.py     # Disambiguation & clarification questions
│   │   ├── context_manager.py   # Bounded sliding-TTL in-memory sessions
│   │   └── response_grounding.py# Zero-fabrication response synthesis
│   ├── languages/
│   │   ├── normalizer.py        # Ge'ez homophone & Oromo orthography normalization
│   │   └── terminology.py       # Multilingual category taxonomy
│   ├── integrations/
│   │   └── go_marketplace_client.py # Restricted client to Go backend (:8080)
│   ├── evaluation/
│   │   ├── benchmark.py         # Evaluation benchmark runner
│   │   └── metrics.py           # Precision, Recall, and F1 metrics
│   ├── model_client.py          # Local model client abstraction
│   └── prompts.py               # Multilingual prompt templates
├── datasets/
│   ├── evaluation/              # Domain, conversational, and entity benchmarks
│   └── schemas/                 # JSON schema for marketplace tool calling
├── tests/                       # Complete Pytest test suite (34 passing tests)
├── requirements.txt             # Pinned production dependencies
└── pytest.ini                   # Pytest test discovery settings
```

---

## 3. Running the Service

### Startup Command
From the `ai/` directory:
```bash
uv run uvicorn app.main:app --host 127.0.0.1 --port 8000
```

### Running Tests
```bash
uv run pytest -v
```

---

## 4. Endpoints

- `GET /health` — Health check probe returning `{ "status": "ok", "service": "huntiq-ai", "version": "2.0.0" }`.
- `GET /ready` — Readiness probe verifying configuration and target Go backend.
- `POST /api/v1/chat` — Conversational chat endpoint accepting `GuestChatRequest` and returning `GuestChatResponse`.

---

## 5. Environment Variables

| Variable | Default | Purpose |
|---|---|---|
| `HOST` | `127.0.0.1` | Loopback bind address |
| `PORT` | `8000` | Port for Python AI service |
| `HUNTIQ_AI_SECRET` | `huntiq-local-ai-secret-key-2026` | Service-to-service authentication secret |
| `GO_BACKEND_URL` | `http://127.0.0.1:8080` | URL of Go backend API |
