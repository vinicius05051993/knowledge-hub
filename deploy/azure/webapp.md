# Azure Web App Deployment

This document describes how to deploy Knowledge Hub to Azure using Azure Web Apps.

## Overview

The application is split into two independent services:

- `knowledge-hub-api`
- `knowledge-hub-worker`

Both services use the same Docker image but run different startup commands.

## Architecture

```text
Azure Web App (API)
    └── knowledge-hub-api

Azure Web App (Worker)
    └── knowledge-hub-worker

Azure Database for MySQL
    └── knowledge_hub

Azure Virtual Machine
    └── OpenSearch
```

## API Web App

The API service handles:

- Document ingestion
- Search requests
- Authentication
- Metrics
- Health checks

### Startup Command

```bash
./api
```

### Required Environment Variables

```env
MYSQL_DSN=
OPENSEARCH_URL=
API_SECRET=
PORT=8080
```

---

## Worker Web App

The worker service is responsible for:

- Synchronizing documents to OpenSearch
- Processing pending upserts
- Processing pending deletes

### Startup Command

```bash
./worker
```

### Required Environment Variables

```env
MYSQL_DSN=
OPENSEARCH_URL=
```

---

## Health Checks

The API should expose:

```http
GET /health
```

Expected response:

```json
{
  "status": "ok"
}
```

---

## Deployment Strategy

Recommended deployment flow:

```text
GitHub Push
    ↓
GitHub Actions
    ↓
Run Tests
    ↓
Build Docker Image
    ↓
Push Container Registry
    ↓
Deploy API Web App
    ↓
Deploy Worker Web App
```

## Future Improvements

- Azure Container Apps
- Kubernetes (AKS)
- Blue/Green deployments
- Automatic scaling
- Distributed tracing