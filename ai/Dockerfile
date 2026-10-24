# Multi-stage minimal production Dockerfile for HUNTIQ AI Service
FROM ghcr.io/astral-sh/uv:0.5.21 AS uv-installer

FROM python:3.12-slim-bookworm AS runner

# Create non-root system user
RUN groupadd -g 10001 huntiq && \
    useradd -u 10001 -g huntiq -s /bin/false -m appuser

# Copy uv binary from official builder
COPY --from=uv-installer /uv /uvx /bin/

WORKDIR /app

# Set production environment variables
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PATH="/app/.venv/bin:$PATH" \
    ENVIRONMENT=production

# Install dependencies into virtual environment
COPY pyproject.toml uv.lock ./
RUN uv sync --frozen --no-dev --no-install-project

# Copy application source code
COPY src/ ./src/
RUN uv sync --frozen --no-dev

# Set permissions and switch to non-root user
RUN chown -R appuser:huntiq /app
USER appuser

EXPOSE 8000

# Health check matching readiness probe
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://127.0.0.1:8000/health/live')" || exit 1

ENTRYPOINT ["uvicorn", "huntiq_ai.main:app", "--host", "0.0.0.0", "--port", "8000"]
