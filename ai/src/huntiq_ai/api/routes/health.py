from fastapi import APIRouter, Depends, status

from huntiq_ai.core.config import Settings, get_settings
from huntiq_ai.schemas.health import LivenessResponse, ReadinessResponse

router = APIRouter(prefix="/health", tags=["Health & Operations"])


@router.get(
    "/live",
    response_model=LivenessResponse,
    status_code=status.HTTP_200_OK,
    summary="Liveness Probe",
    description="Returns HTTP 200 if the process is responsive. Used by container orchestrators.",
)
async def liveness_probe() -> LivenessResponse:
    return LivenessResponse(
        status="live",
        service="huntiq-ai",
        version="0.1.0",
    )


@router.get(
    "/ready",
    response_model=ReadinessResponse,
    status_code=status.HTTP_200_OK,
    summary="Readiness Probe",
    description="Returns HTTP 200 if configuration is valid and models are ready to serve.",
)
async def readiness_probe(settings: Settings = Depends(get_settings)) -> ReadinessResponse:
    # At foundation stage, no models are loaded; service is ready to receive requests
    return ReadinessResponse(
        status="ready",
        service="huntiq-ai",
        version="0.1.0",
        environment=settings.environment,
        models_loaded=False,
        details={
            "model_count": 0,
        },
    )
