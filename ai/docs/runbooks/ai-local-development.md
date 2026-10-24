# HUNTIQ AI Local Development Runbook

**Document Version:** 1.0.0  
**Status:** Canonical Local Development Guide  

---

## 1. Prerequisites

Ensure your workstation has:
- **Python:** Version 3.12 or newer (`python3 --version`).
- **uv:** Version 0.10+ installed (`curl -LsSf https://astral.sh/uv/install.sh | bash`).
- **Git:** Standard Git installation.

---

## 2. Quickstart: Setting up the AI Service

### Step 1: Navigate to the AI directory
```bash
cd "GO PROJECT/sellbuy/ai"
```

### Step 2: Initialize configuration
Copy the example environment configuration:
```bash
cp .env.example .env
```
*(The default placeholder secret in `.env.example` is pre-configured for local testing).*

### Step 3: Install dependencies using `uv`
```bash
uv sync
```
This automatically sets up a virtual environment in `ai/.venv` and installs all locked dependencies from `uv.lock`.

### Step 4: Run the automated test suite
```bash
uv run pytest -v
```
All unit, integration, and contract tests should pass in under 2 seconds.

### Step 5: Start the local development server
```bash
uv run uvicorn huntiq_ai.main:app --host 127.0.0.1 --port 8000 --reload
```

---

## 3. Verifying Local Service Endpoints

### 1. Check Liveness Probe (No Auth Required)
```bash
curl -s http://127.0.0.1:8000/health/live | jq .
```
Expected output:
```json
{
  "status": "live",
  "service": "huntiq-ai",
  "version": "0.1.0"
}
```

### 2. Check Readiness Probe (No Auth Required)
```bash
curl -s http://127.0.0.1:8000/health/ready | jq .
```
Expected output:
```json
{
  "status": "ready",
  "service": "huntiq-ai",
  "version": "0.1.0",
  "environment": "development",
  "models_loaded": false,
  "details": {
    "model_count": 0
  }
}
```

### 3. Test Inter-Service Authentication
Verify that requests without the secret header are rejected:
```bash
curl -i -X POST http://127.0.0.1:8000/api/v1/predict/query-intent \
  -H "Content-Type: application/json" \
  -d '{"query": "Samsung A16"}'
```
Expected output: `401 Unauthorized`.

Now test with the secret header:
```bash
curl -i -X POST http://127.0.0.1:8000/api/v1/predict/query-intent \
  -H "Content-Type: application/json" \
  -H "X-HUNTIQ-AI-SECRET: dev_secret_change_in_production" \
  -H "X-Correlation-ID: local-test-1" \
  -d '{"query": "Samsung A16"}'
```

---

## 4. Code Formatting and Linting

Run Ruff to inspect and format Python code:
```bash
cd ai
uv run ruff check .
uv run ruff format .
```

---

## 5. Troubleshooting Common Issues

### Issue 1: `uv: command not found`
**Fix:** Ensure `~/.local/bin` is in your `PATH`:
```bash
export PATH="$HOME/.local/bin:$PATH"
```
Add this line to your `~/.bashrc` or `~/.zshrc`.

### Issue 2: `Port 8000 already in use`
**Fix:** Either terminate the existing process or run Uvicorn on an alternate port:
```bash
uv run uvicorn huntiq_ai.main:app --host 127.0.0.1 --port 8001
```
Remember to update `AI_SERVICE_URL` in `backend/.env` if port is changed.

### Issue 3: `401 Unauthorized` on internal calls
**Fix:** Verify that `HUNTIQ_AI_SECRET_KEY` in `ai/.env` matches `AI_SERVICE_SECRET` in `backend/.env`.
