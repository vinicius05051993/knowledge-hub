# Knowledge Hub

Knowledge Hub is a lightweight document indexing and search platform built with Go, MySQL, and OpenSearch.

The project is designed to provide a simple and scalable way to index documents, apply structured filters, and perform fast searches across multiple namespaces. It follows a clean architecture approach with a clear separation between API, persistence, synchronization, and search layers.

## Project Roadmap

### V1 — Basic Search

**Branch:** `basic`

The first version focuses on a solid foundation:

- REST API for document management
- MySQL as the source of truth
- OpenSearch for indexing and search
- Namespace isolation
- Structured document filters
- API key authentication
- Background synchronization worker
- Automated tests
- Docker-based local development

**Goal:** Provide a reliable and maintainable document search platform.

---

### V2 — Semantic Search

**Branch:** `semantic-search`

Planned improvements:

- Embedding generation pipeline
- Semantic search using vector embeddings
- Hybrid search (keyword + semantic)
- OpenSearch vector indexes
- Improved ranking and relevance
- Search quality benchmarks
- Performance optimizations

**Goal:** Allow users to find relevant documents even when exact keywords are not present.

---

### V3 — AI Chatbot on Kubernetes

**Branch:** `ia-chatbot-kubernets`

Planned improvements:

- Retrieval-Augmented Generation (RAG)
- LLM integration
- AI-powered chatbot interface
- Context-aware document answering
- Kubernetes deployment support
- Horizontal scaling
- Observability and monitoring
- Production-ready cloud architecture

**Goal:** Transform the search platform into an AI knowledge assistant capable of answering questions based on indexed documents.

---

## Long-Term Vision

Knowledge Hub aims to evolve from a traditional search engine into a complete enterprise knowledge platform by combining:

- Structured search
- Semantic retrieval
- AI-assisted question answering

while keeping the architecture simple, scalable, and self-hosted.