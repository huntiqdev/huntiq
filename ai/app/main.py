"""HUNTIQ Python AI Service Entry Point (FastAPI)."""

from contextlib import asynccontextmanager
from typing import Any

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from .api.routes import chat, query
from .config import config


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Startup and shutdown lifecycle management."""
    # Startup logging or pool warming can happen here
    yield
    # Cleanup on shutdown


app = FastAPI(
    title="HUNTIQ AI Intelligence Service",
    description="Sovereign AI Natural Language Understanding & Reasoning Platform for Ethiopia Marketplace",
    version="2.0.0",
    lifespan=lifespan,
)

# CORS restricted to internal loopback and authorized origins
app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://127.0.0.1:8080",
        "http://localhost:8080",
        "http://localhost:3000",
    ],
    allow_credentials=True,
    allow_methods=["GET", "POST", "OPTIONS"],
    allow_headers=["*"],
)

# Register routes
app.include_router(chat.router)
app.include_router(query.router)


@app.get("/health", tags=["system"])
async def health_check() -> dict[str, str]:
    """Basic health check probe."""
    return {"status": "ok", "service": "huntiq-ai", "version": "2.0.0"}


@app.get("/ready", tags=["system"])
async def readiness_check() -> dict[str, Any]:
    """Readiness probe verifying operational state."""
    return {
        "ready": True,
        "service": "huntiq-ai",
        "target_go_backend": config.go_backend_url,
        "allowed_tools": config.allowed_tools,
    }


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("app.main:app", host=config.host, port=config.port, reload=False)
