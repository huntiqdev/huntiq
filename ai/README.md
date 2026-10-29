# HUNTIQ AI — Local Qwen3 4B System

This directory contains the local model configuration, offline model storage, evaluation suite, and Python tools for the HUNTIQ Local AI Assistant running on Ollama with `qwen3:4b`.

## Architecture Overview

1. **Local Model Storage:** Stored directly inside the project at `ai/models/ollama` (configured via `OLLAMA_MODELS` and container bind mounts).
2. **Local Model Runtime:** Ollama running locally at `http://127.0.0.1:11434` with zero external cloud dependencies.
3. **Base Pretrained Model:** `qwen3:4b` (4.0B parameters, Q4_K_M quantization).
4. **Derived Model:** `huntiq-assistant:latest` built from [`Modelfile`](./Modelfile) with customized multilingual directives (English, Amharic, Afaan Oromoo, Tigrinya) and live marketplace tool definitions.
5. **Trusted Sovereign Authority:** Go Backend (`backend/`) manages all authentication, authorization, session persistence, and parameter-safe database queries.
6. **Database Source of Truth:** PostgreSQL 16. Listings and prices are retrieved dynamically via authorized Go tool executions; never hallucinated.
7. **Frontend:** Next.js 15 + TypeScript conversational interface (`/ai`) and query-understanding search bar (`SearchBar.tsx`).

## Directory Structure

```
ai/
├── Modelfile                  # Ollama Modelfile deriving huntiq-assistant from qwen3:4b
├── README.md                  # System documentation
├── .env.example               # Environment variables template
├── requirements.txt           # Minimal Python dependencies
├── pytest.ini                 # Pytest configuration
├── scripts/
│   ├── setup_local.sh         # Idempotent environment and model setup script
│   └── verify_local.sh        # Diagnostic verification script
├── models/
│   ├── .gitkeep               # Preserves models directory in Git
│   └── ollama/                # Local model weights, blobs, and manifests
├── .venv/                     # Project-local Python virtual environment
├── app/
│   ├── __init__.py
│   ├── config.py              # Ollama and assistant configuration
│   ├── model_client.py        # Python client for Ollama API interaction
│   ├── prompts.py             # Multilingual system prompts and tool schemas
│   └── evaluation.py          # Multilingual and functional evaluation runner
├── data/
│   ├── README.md              # Offline training and evaluation data instructions
│   └── .gitkeep
└── tests/
    ├── __init__.py
    ├── test_model_client.py   # Ollama connection and chat contract tests
    ├── test_prompts.py        # Prompt construction and multilingual integrity tests
    └── test_offline_mode.py   # Verification that zero cloud/external calls occur
```

## Quick Start & Setup

### 1. Automated Idempotent Setup
Run the setup script from the project root:
```bash
./ai/scripts/setup_local.sh
```
This automatically:
- Ensures `ai/models/ollama` and `ai/.venv` exist.
- Installs minimal dependencies from `ai/requirements.txt` into `ai/.venv`.
- Verifies the local Ollama runtime mounted to `ai/models/ollama`.
- Verifies `qwen3:4b` and builds `huntiq-assistant:latest` if missing.

### 2. Verify Local Health & Offline Readiness
Run the verification script:
```bash
./ai/scripts/verify_local.sh
```
This validates:
- Daemon reachability (`127.0.0.1:11434`).
- Storage presence in `ai/models/ollama/blobs`.
- Registration of `qwen3:4b` and `huntiq-assistant:latest`.
- Python unit tests in `ai/.venv`.
- Live test inference from `huntiq-assistant:latest`.
- Strict local loopback enforcement (no cloud LLM API calls).

### 3. Interactive CLI Usage
```bash
ollama run huntiq-assistant
```

## Python Virtual Environment

Activate the local project virtualenv:
```bash
source ai/.venv/bin/activate
pytest ai/tests/ -v
```

## Configuration Reference

| Variable | Default Value | Description |
|---|---|---|
| `OLLAMA_BASE_URL` | `http://127.0.0.1:11434` | Local Ollama daemon endpoint |
| `OLLAMA_MODELS` | `$PROJECT_ROOT/ai/models/ollama` | Directory storing model blobs & manifests |
| `OLLAMA_BASE_MODEL` | `qwen3:4b` | Base pretrained model |
| `OLLAMA_HUNTIQ_MODEL` | `huntiq-assistant:latest` | Derived customized assistant |
| `AI_REQUEST_TIMEOUT` | `90.0` | Timeout in seconds for CPU generation |
| `AI_CONTEXT_LENGTH` | `4096` | Context window size in tokens |
