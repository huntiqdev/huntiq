#!/usr/bin/env bash
# ==============================================================================
# HUNTIQ Local AI Setup Script
# Configures local Ollama model storage at ai/models/ollama and Python virtualenv
# ==============================================================================

set -euo pipefail

# 1. Resolve absolute paths safely (handles spaces)
AI_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO_ROOT="$(cd "$AI_DIR/.." && pwd)"
MODELS_DIR="$AI_DIR/models/ollama"
VENV_DIR="$AI_DIR/.venv"
MODELFILE="$AI_DIR/Modelfile"
REQ_FILE="$AI_DIR/requirements.txt"

echo "=== HUNTIQ Local AI Setup ==="
echo "Project Root : $REPO_ROOT"
echo "AI Directory : $AI_DIR"
echo "Models Store : $MODELS_DIR"
echo "Virtualenv   : $VENV_DIR"
echo "----------------------------------------------------"

# 2. Ensure model storage directory exists
mkdir -p "$MODELS_DIR"
mkdir -p "$AI_DIR/data"
mkdir -p "$AI_DIR/models"
touch "$AI_DIR/models/.gitkeep"

# 3. Setup Python Virtual Environment (Lean, Local)
if [ ! -d "$VENV_DIR" ] || [ ! -f "$VENV_DIR/bin/pytest" ]; then
    echo "[1/4] Initializing Python virtual environment in $VENV_DIR..."
    python3 -m venv "$VENV_DIR"
    "$VENV_DIR/bin/pip" install --upgrade pip --quiet
    "$VENV_DIR/bin/pip" install -r "$REQ_FILE" --quiet
    echo "       Virtual environment ready."
else
    echo "[1/4] Python virtual environment already present and verified."
fi

# 4. Verify / Start Ollama Container with Local Model Mount
OLLAMA_PORT=11434
OLLAMA_URL="http://127.0.0.1:$OLLAMA_PORT"

echo "[2/4] Verifying Ollama service runtime..."

is_ollama_running() {
    curl -s --connect-timeout 2 "$OLLAMA_URL/api/tags" >/dev/null 2>&1
}

if is_ollama_running; then
    echo "       Ollama service is currently running at $OLLAMA_URL."
else
    # Check if Docker container exists
    if command -v docker >/dev/null 2>&1; then
        if docker ps -a --format '{{.Names}}' | grep -Eq '^ollama$'; then
            echo "       Starting existing 'ollama' container..."
            docker start ollama >/dev/null
        else
            echo "       Creating 'ollama' container with local model mount..."
            docker run -d \
                --name ollama \
                --restart unless-stopped \
                -p 127.0.0.1:$OLLAMA_PORT:11434 \
                -v "$MODELS_DIR:/root/.ollama/models" \
                -e OLLAMA_MODELS=/root/.ollama/models \
                -e OLLAMA_HOST=0.0.0.0:11434 \
                ollama/ollama:latest >/dev/null
        fi
        # Wait for daemon ready
        for i in {1..15}; do
            if is_ollama_running; then
                break
            fi
            sleep 1
        done
    else
        echo "[ERROR] Docker is required to run the local Ollama runtime on this system."
        exit 1
    fi
fi

if ! is_ollama_running; then
    echo "[ERROR] Ollama daemon failed to start at $OLLAMA_URL."
    exit 1
fi

# 5. Check / Pull Base Model
echo "[3/4] Checking base model qwen3:4b..."
INSTALLED_MODELS=$(curl -s "$OLLAMA_URL/api/tags" | grep -o '"name":"[^"]*"' | cut -d'"' -f4 || true)

if echo "$INSTALLED_MODELS" | grep -Eq 'qwen3:4b'; then
    echo "       Base model qwen3:4b is already present in $MODELS_DIR."
else
    echo "       Base model not found. Pulling qwen3:4b into $MODELS_DIR..."
    docker exec -i ollama ollama pull qwen3:4b
fi

# 6. Check / Build Derived HUNTIQ Model
echo "[4/4] Checking derived model huntiq-assistant:latest..."
if echo "$INSTALLED_MODELS" | grep -Eq 'huntiq-assistant:latest'; then
    echo "       Derived model huntiq-assistant:latest is already registered."
else
    echo "       Building huntiq-assistant from $MODELFILE..."
    docker cp "$MODELFILE" ollama:/tmp/Modelfile
    docker exec -i ollama ollama create huntiq-assistant -f /tmp/Modelfile
    docker exec -i ollama rm -f /tmp/Modelfile
fi

echo "----------------------------------------------------"
echo "HUNTIQ Local AI setup completed successfully."
echo "Active models stored at: $MODELS_DIR"
