# HUNTIQ AI Integration Specification — Go to Python Contract

**Document Version:** 1.0.0  
**Target:** Senior Go Engineers & AI Systems Integrators  
**Status:** Approved Specification  

---

## 1. Overview

This document specifies the exact contract, transport protocol, security parameters, and error handling for communication between the authoritative **HUNTIQ Go Backend** and the internal **HUNTIQ Python AI Service**.

---

## 2. Transport & Network Parameters

- **Protocol:** HTTP/1.1 or HTTP/2 over private network / loopback.
- **Default Base URL:** `http://127.0.0.1:8000` (configurable via `AI_SERVICE_URL`).
- **Connection Pooling:** Go's `http.Client` with `MaxIdleConns: 100`, `IdleConnTimeout: 90s`.
- **Request Timeout:** Hard limit of **500ms** enforced via Go `context.WithTimeout`.
- **Maximum Request Body Size:** **64 KB** (65,536 bytes). Oversized payloads receive `413 Payload Too Large`.

---

## 3. Headers & Inter-Service Security

Every request initiated by the Go backend must transmit:

| Header | Required | Purpose | Example |
| :--- | :--- | :--- | :--- |
| `X-HUNTIQ-AI-SECRET` | Yes | Pre-shared internal authentication secret. | `huntiq_sec_live_9f83a...` |
| `X-Correlation-ID` | Yes | Distributed tracing correlation identifier. | `cid-98a7-4b12-98ef` |
| `Content-Type` | Yes | Must specify JSON. | `application/json` |
| `Accept` | Yes | Expected return format. | `application/json` |

---

## 4. Operational Endpoints

### 4.1 Liveness Probe
- **Method / Path:** `GET /health/live`
- **Auth Required:** No
- **Response (200 OK):**
```json
{
  "status": "live",
  "service": "huntiq-ai",
  "version": "0.1.0"
}
```

### 4.2 Readiness Probe
- **Method / Path:** `GET /health/ready`
- **Auth Required:** No
- **Response (200 OK):**
```json
{
  "status": "ready",
  "service": "huntiq-ai",
  "version": "0.1.0",
  "environment": "production",
  "models_loaded": false,
  "details": {
    "model_count": 0
  }
}
```

---

## 5. Future Predictive Contracts (Planned Schemas)

### 5.1 Query Intent Classification (`POST /api/v1/predict/query-intent`)

#### Request Schema:
```json
{
  "query": "Samsung Galaxy A16 128GB",
  "locale": "am",
  "limit_categories": 5
}
```

#### Response Schema (200 OK):
```json
{
  "intent": "product_search",
  "confidence": 0.94,
  "suggested_category_ids": [142, 105],
  "extracted_entities": {
    "brand": "Samsung",
    "series": "Galaxy A",
    "model": "A16",
    "attributes": {
      "storage": "128GB"
    }
  },
  "execution_time_ms": 12.4
}
```

### 5.2 Category Suggestion (`POST /api/v1/predict/category-suggestion`)

#### Request Schema:
```json
{
  "title": "Toyota Hilux Double Cab 2022 Diesel",
  "description": "Clean white pickup, manual transmission, 4WD, 45000km, imported from Dubai.",
  "locale": "en",
  "top_k": 3
}
```

#### Response Schema (200 OK):
```json
{
  "suggestions": [
    {
      "category_id": 3704,
      "category_slug": "commercial-pickup-trucks",
      "confidence": 0.96
    },
    {
      "category_id": 3700,
      "category_slug": "vehicles-cars",
      "confidence": 0.03
    }
  ],
  "execution_time_ms": 15.1
}
```

---

## 6. Standard Error Format

All error responses from the AI service follow RFC 7807 problem details:

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Missing or invalid internal service authentication secret",
    "correlation_id": "cid-98a7-4b12-98ef",
    "timestamp": "2026-10-24T08:15:30Z"
  }
}
```

---

## 7. Recommended Go Client Implementation Pattern

```go
package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AIClient struct {
	baseURL    string
	secretKey  string
	httpClient *http.Client
}

func NewAIClient(baseURL, secretKey string) *AIClient {
	return &AIClient{
		baseURL:   baseURL,
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 500 * time.Millisecond,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				DisableCompression: false,
			},
		},
	}
}

// PredictWithFallback illustrates strict fallback on timeout or error
func (c *AIClient) CheckHealth(ctx context.Context) (bool, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, c.baseURL+"/health/ready", nil)
	if err != nil {
		return false, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Non-fatal: AI service unreachable, fallback immediately
		return false, fmt.Errorf("ai service unreachable: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
```
