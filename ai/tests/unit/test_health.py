from fastapi.testclient import TestClient


def test_liveness_probe_returns_200(client: TestClient):
    """Liveness probe must return 200 OK with expected fields."""
    response = client.get("/health/live")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "live"
    assert data["service"] == "huntiq-ai"
    assert data["version"] == "0.1.0"
    assert "X-Correlation-ID" in response.headers


def test_readiness_probe_returns_200(client: TestClient):
    """Readiness probe must return 200 OK with operational details."""
    response = client.get("/health/ready")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ready"
    assert data["service"] == "huntiq-ai"
    assert data["environment"] == "test"
    assert data["models_loaded"] is False
    assert "model_count" in data["details"]


def test_custom_correlation_id_is_preserved(client: TestClient):
    """Preserves client-supplied X-Correlation-ID header."""
    custom_cid = "custom-test-cid-12345"
    response = client.get("/health/live", headers={"X-Correlation-ID": custom_cid})
    assert response.status_code == 200
    assert response.headers["X-Correlation-ID"] == custom_cid
