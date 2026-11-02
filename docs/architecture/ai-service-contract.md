# HUNTIQ AI Service Contracts

## 1. Public Visitor AI Chat Contract

**Endpoint:** `POST /api/ai/chat` (on Go Backend :8080) / `POST /api/v1/chat` (on Python Service :8000)

### Request Payload
```json
{
  "message": "I need a used Samsung phone under 20,000 birr in Bole",
  "locale": "en",
  "session_id": "b3c8f921-789a-41d3-a4c1-112233445566",
  "history": [
    {"role": "user", "content": "Hello"},
    {"role": "assistant", "content": "Hello! How can I help you today?"}
  ]
}
```

### Response Payload
```json
{
  "session_id": "b3c8f921-789a-41d3-a4c1-112233445566",
  "greeting_prefix": "Hello! Welcome to HUNTIQ.",
  "answer": "Here are 3 verified Samsung phones under 20,000 ETB in Bole, Addis Ababa:",
  "detected_locale": "en",
  "intents": ["marketplace_search"],
  "category_groups": [
    {
      "category_name": "Mobile Phones",
      "category_slug": "phones-telecommunications",
      "search_term": "samsung",
      "total_count": 3,
      "listings": [
        {
          "id": 104,
          "title": "Samsung Galaxy A54 5G 128GB",
          "price": 18500,
          "currency": "ETB",
          "listing_type": "sale",
          "condition": "used",
          "seller_tier": "verified",
          "city": "Addis Ababa",
          "subcity": "Bole",
          "detail_url": "/products/104",
          "image_url": "https://.../a54.jpg",
          "contact_action": "Buy / Contact Seller",
          "availability_status": "Available"
        }
      ]
    }
  ],
  "listings": [...],
  "clarification": null,
  "latency_ms": 112.4,
  "model": "huntiq-intelligence-v2"
}
```

---

## 2. Internal Marketplace Tool Execution Contract

**Endpoint:** `POST /api/internal/ai/tools/execute` (on Go Backend :8080)  
**Security:** Header `X-HUNTIQ-AI-SECRET: <secret_key>`

### Request Payload
```json
{
  "tool_name": "search_marketplace_listings",
  "arguments": {
    "search_term": "samsung",
    "category_slug": "phones-telecommunications",
    "min_price": 0,
    "max_price": 20000,
    "condition": "used",
    "location": "Bole",
    "listing_type": "sale",
    "limit": 5
  }
}
```

### Response Payload
```json
{
  "success": true,
  "tool_name": "search_marketplace_listings",
  "raw_result": "{\"total\": 3}",
  "total_matching": 3,
  "listings": [
    {
      "id": 104,
      "title": "Samsung Galaxy A54 5G 128GB",
      "price": 18500,
      "currency": "ETB",
      "listing_type": "sale",
      "condition": "used",
      "seller_tier": "verified",
      "city": "Addis Ababa",
      "subcity": "Bole",
      "detail_url": "/products/104",
      "image_url": "https://.../a54.jpg",
      "contact_action": "Buy / Contact Seller",
      "availability_status": "Available"
    }
  ]
}
```
