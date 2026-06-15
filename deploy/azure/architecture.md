# Architecture

## Current Architecture (V1)

Knowledge Hub follows a simple and scalable architecture.

MySQL acts as the source of truth, while OpenSearch is used exclusively for search operations.

```text
                +------------------+
                |     Clients      |
                +--------+---------+
                         |
                         v
                +------------------+
                |       API        |
                |   Azure Web App  |
                +--------+---------+
                         |
              +----------+----------+
              |                     |
              v                     v
     +----------------+   +----------------+
     |     MySQL      |   |   OpenSearch   |
     | Source of Truth|   | Search Engine  |
     +----------------+   +----------------+
              ^
              |
              |
     +------------------+
     |      Worker      |
     | Azure Web App    |
     +------------------+
```

## Components

### API

Responsible for:

- Document creation
- Document updates
- Search requests
- Authentication
- Validation

### MySQL

Stores:

- Documents
- Filters
- API keys
- Synchronization status

MySQL is considered the source of truth for all data.

### Worker

Responsible for:

- Synchronizing documents to OpenSearch
- Processing pending updates
- Processing pending deletions

This allows search indexing to happen asynchronously.

### OpenSearch

Responsible for:

- Full-text search
- Filtered search
- Fast retrieval of document identifiers

OpenSearch is treated as a read-only search layer and can be rebuilt from MySQL at any time.

---

# Future Roadmap

## V2 — Semantic Search

Branch: `semantic-search`

Planned additions:

- Embedding generation
- Vector indexing
- Semantic search
- Hybrid search (keyword + vector)
- Relevance improvements

Architecture:

```text
Clients
   |
   v
 API
   |
   +------> OpenSearch
   |           |
   |           +--> Keyword Search
   |           +--> Vector Search
   |
   +------> Embedding Service
```

---

## V3 — AI Chatbot on Kubernetes

Branch: `ia-chatbot-kubernets`

Planned additions:

- Retrieval-Augmented Generation (RAG)
- LLM integration
- AI chatbot interface
- Kubernetes deployment
- Horizontal scaling
- Observability stack

Architecture:

```text
Users
   |
   v
Chatbot API
   |
   +--> Retriever
   |       |
   |       +--> OpenSearch
   |
   +--> LLM
   |
   +--> Knowledge Hub API
```

The long-term goal is to evolve Knowledge Hub from a traditional search platform into a complete AI-powered knowledge assistant.