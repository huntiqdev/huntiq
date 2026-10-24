from fastapi.testclient import TestClient

from huntiq_ai.schemas.health import LivenessResponse, ReadinessResponse


def test_openapi_schema_generation(client: TestClient):
    """OpenAPI schema must be valid and define health endpoints."""
    response = client.get("/openapi.json")
    assert response.status_code == 200
    schema = response.json()
    assert "paths" in schema
    assert "/health/live" in schema["paths"]
    assert "/health/ready" in schema["paths"]


def test_liveness_contract_schema(client: TestClient):
    """Response must parse strictly into LivenessResponse Pydantic model."""
    response = client.get("/health/live")
    assert response.status_code == 200
    model = LivenessResponse.model_validate(response.json())
    assert model.status == "live"
    assert model.service == "huntiq-ai"


def test_readiness_contract_schema(client: TestClient):
    """Response must parse strictly into ReadinessResponse Pydantic model."""
    response = client.get("/health/ready")
    assert response.status_code == 200
    model = ReadinessResponse.model_validate(response.json())
    assert model.status == "ready"
    assert model.service == "huntiq-ai"


def test_error_response_conforms_to_rfc7807(client: TestClient):
    """Error responses must strictly follow the expected contract structure."""
    oversized = "a" * 1500
    response = client.post(
        "/health/live",
        content=oversized,
        headers={"Content-Length": "1500"},
    )
    assert response.status_code == 413
    data = response.json()
    assert "error" in data
    assert "code" in data["error"]
    assert "message" in data["error"]
    assert "correlation_id" in data["error"]
    assert "timestamp" in data["error"]
