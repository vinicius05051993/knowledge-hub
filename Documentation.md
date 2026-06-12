# Knowledge Hub API Documentation

## Overview

The Knowledge Hub API provides document indexing and search capabilities using namespace-based isolation.

Authentication is performed using API keys sent in the `X-API-Key` header.

---

## Authentication

All protected endpoints require an API key.

### Header

```http
X-API-Key: sk_live_xxxxxxxxxxxxxxxxx
```

### Authentication Errors

#### Missing API Key

```http
401 Unauthorized
```

Response:

```text
missing api key
```

#### Invalid API Key

```http
401 Unauthorized
```

Response:

```text
invalid api key
```

---

# Endpoints

## Health Check

Verifies that the API process is running.

### Request

```http
GET /health
```

### Authentication

Not required.

### Example

```bash
curl http://localhost:8080/health
```

---

## Readiness Check

Verifies application readiness and dependency availability.

### Request

```http
GET /ready
```

### Authentication

Not required.

### Example

```bash
curl http://localhost:8080/ready
```

---

## Validate Authentication

Returns information about the authenticated API key.

### Request

```http
GET /protected
```

### Headers

| Header | Required |
|----------|----------|
| X-API-Key | Yes |

### Example

```bash
curl \
  -H "X-API-Key: sk_live_xxxxxxxxxxxxxxxxx" \
  http://localhost:8080/protected
```

### Response

```json
{
  "api_key_id": 1,
  "namespace": "magento"
}
```

---

# Upsert Documents

Creates or updates one or more documents.

### Request

```http
POST /documents/upsert
```

### Headers

| Header | Required |
|----------|----------|
| X-API-Key | Yes |
| Content-Type: application/json | Yes |

### Body

```json
{
  "documents": [
    {
      "external_id": "product-123",
      "title": "Magento Installation Guide",
      "text": "Install Magento using Composer.",
      "payload": {
        "category": "documentation",
        "version": "2.4"
      }
    }
  ]
}
```

### Fields

| Field | Type | Required |
|---------|---------|----------|
| documents | array | Yes |
| documents[].external_id | string | Yes |
| documents[].title | string | No |
| documents[].text | string | No |
| documents[].payload | object | No |

### Example

```bash
curl -X POST http://localhost:8080/documents/upsert \
  -H "X-API-Key: sk_live_xxxxxxxxxxxxxxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      {
        "external_id": "product-123",
        "title": "Magento Installation Guide",
        "text": "Install Magento using Composer.",
        "payload": {
          "category": "documentation"
        }
      }
    ]
  }'
```

### Response

```json
{
  "success": true,
  "processed": 1
}
```

### Validation Errors

```http
400 Bad Request
```

```text
invalid request
```

or

```text
invalid payload
```

---

# Delete Documents

Deletes one or more documents by external ID.

### Request

```http
DELETE /documents
```

### Headers

| Header | Required |
|----------|----------|
| X-API-Key | Yes |
| Content-Type: application/json | Yes |

### Body

```json
{
  "external_ids": [
    "product-123",
    "product-456"
  ]
}
```

### Fields

| Field | Type | Required |
|---------|---------|----------|
| external_ids | array[string] | Yes |

### Example

```bash
curl -X DELETE http://localhost:8080/documents \
  -H "X-API-Key: sk_live_xxxxxxxxxxxxxxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "external_ids": [
      "product-123",
      "product-456"
    ]
  }'
```

### Response

Successful deletion returns:

```http
204 No Content
```

### Validation Errors

```http
400 Bad Request
```

```text
invalid request
```

---

# Search Documents

Searches documents within the authenticated namespace.

### Request

```http
POST /search
```

### Headers

| Header | Required |
|----------|----------|
| X-API-Key | Yes |
| Content-Type: application/json | Yes |

### Body

```json
{
  "query": "magento",
  "offset": 0,
  "limit": 20,
  "filters": {
    "category": "documentation"
  }
}
```

### Fields

| Field | Type | Required |
|---------|---------|----------|
| query | string | No |
| offset | integer | No |
| limit | integer | No |
| filters | object | No |

### Pagination Rules

| Rule | Value |
|--------|--------|
| Negative offset | Converted to 0 |
| Default limit | 100 |
| Maximum limit | 100 |

### Example

```bash
curl -X POST http://localhost:8080/search \
  -H "X-API-Key: sk_live_xxxxxxxxxxxxxxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "magento",
    "offset": 0,
    "limit": 20,
    "filters": {
      "category": "documentation"
    }
  }'
```

### Response

```json
[
  {
    "namespace": "magento",
    "external_id": "product-123",
    "title": "Magento Installation Guide",
    "text": "Install Magento using Composer.",
    "payload": {
      "category": "documentation"
    },
    "highlights": {
      "text": "<mark>Magento</mark> Installation Guide"
    }
  }
]
```

### Response Fields

| Field | Description |
|---------|-------------|
| namespace | Document namespace |
| external_id | External document identifier |
| title | Document title |
| text | Indexed text |
| payload | Custom JSON payload |
| highlights | Highlighted search fragments (when available) |

### Validation Errors

```http
400 Bad Request
```

```text
invalid request
```

---

# Metrics

Prometheus metrics are available at:

```http
GET /metrics
```

No authentication is required.

---

# Quick Start

### Index Documents

```bash
curl -X POST http://localhost:8080/documents/upsert \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      {
        "external_id": "doc-1",
        "title": "Getting Started",
        "text": "Welcome to Knowledge Hub"
      }
    ]
  }'
```

### Search

```bash
curl -X POST http://localhost:8080/search \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Welcome"
  }'
```

### Delete

```bash
curl -X DELETE http://localhost:8080/documents \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "external_ids": [
      "doc-1"
    ]
  }'
```

# Usage Examples

## Index a Single Document

```bash
curl -X POST http://localhost:8080/documents/upsert \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      {
        "external_id": "product-1001",
        "title": "Running Shoes",
        "text": "Lightweight running shoes for everyday training.",
        "payload": {
          "category": "sports",
          "brand": "Acme"
        }
      }
    ]
  }'
```

---

## Index Multiple Documents (Batch Upsert)

```bash
curl -X POST http://localhost:8080/documents/upsert \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      {
        "external_id": "product-1001",
        "title": "Running Shoes",
        "text": "Lightweight running shoes for everyday training.",
        "payload": {
          "category": "sports",
          "brand": "Acme"
        }
      },
      {
        "external_id": "product-1002",
        "title": "Basketball Shoes",
        "text": "High-top basketball shoes with ankle support.",
        "payload": {
          "category": "sports",
          "brand": "Acme"
        }
      },
      {
        "external_id": "product-1003",
        "title": "Tennis Shoes",
        "text": "Durable shoes designed for tennis courts.",
        "payload": {
          "category": "sports",
          "brand": "Acme"
        }
      }
    ]
  }'
```

---

## Update an Existing Document

Documents with the same `external_id` are updated automatically.

```bash
curl -X POST http://localhost:8080/documents/upsert \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      {
        "external_id": "product-1001",
        "title": "Running Shoes V2",
        "text": "Updated product description.",
        "payload": {
          "category": "sports",
          "brand": "Acme",
          "status": "updated"
        }
      }
    ]
  }'
```

---

## Search by Text

```bash
curl -X POST http://localhost:8080/search \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "running shoes"
  }'
```

---

## Search with Filters

```bash
curl -X POST http://localhost:8080/search \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "shoes",
    "filters": {
      "category": "sports",
      "brand": "Acme"
    }
  }'
```

---

## Search with Pagination

```bash
curl -X POST http://localhost:8080/search \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "shoes",
    "offset": 0,
    "limit": 20
  }'
```

---

## Delete a Single Document

```bash
curl -X DELETE http://localhost:8080/documents \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "external_ids": [
      "product-1001"
    ]
  }'
```

---

## Delete Multiple Documents (Batch Delete)

```bash
curl -X DELETE http://localhost:8080/documents \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "external_ids": [
      "product-1001",
      "product-1002",
      "product-1003",
      "product-1004",
      "product-1005"
    ]
  }'
```

---

## Full Catalog Synchronization Example

### Step 1: Index Products

```bash
curl -X POST http://localhost:8080/documents/upsert \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      {
        "external_id": "sku-001",
        "title": "Product A",
        "text": "Description of Product A"
      },
      {
        "external_id": "sku-002",
        "title": "Product B",
        "text": "Description of Product B"
      },
      {
        "external_id": "sku-003",
        "title": "Product C",
        "text": "Description of Product C"
      }
    ]
  }'
```

### Step 2: Search Products

```bash
curl -X POST http://localhost:8080/search \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Product"
  }'
```

### Step 3: Remove Discontinued Products

```bash
curl -X DELETE http://localhost:8080/documents \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "external_ids": [
      "sku-002",
      "sku-003"
    ]
  }'
```

---

## Recommended Batch Sizes

| Operation | Recommended Batch Size |
|------------|------------------------|
| Upsert | 100 - 1,000 documents |
| Delete | 100 - 1,000 documents |
| Initial Import | 500 - 1,000 documents per request |
| Incremental Updates | 1 - 100 documents |

For large imports, split the dataset into multiple requests instead of sending tens of thousands of documents in a single request.