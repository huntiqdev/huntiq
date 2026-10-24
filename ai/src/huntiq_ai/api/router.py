from fastapi import APIRouter

from huntiq_ai.api.routes import health

api_router = APIRouter()

# Operational Health Probes
api_router.include_router(health.router)
