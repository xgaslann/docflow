# DocFlow Roadmap

This document outlines planned features and improvements. Things may change based on feedback and priorities.

---

## ✅ Current Version (v1.0) — COMPLETED

### Core Features
- [x] Markdown ↔ PDF conversion
- [x] Multi-file upload and merge
- [x] Live preview with syntax highlighting
- [x] Light/Dark theme support
- [x] Docker deployment

### Multi-Format Support
- [x] DOCX → Markdown conversion
- [x] Excel (XLSX/XLS) → Markdown tables
- [x] CSV → Markdown tables
- [x] TXT → Markdown

### RAG Pipeline
- [x] Semantic chunking with configurable size/overlap
- [x] Heading-aware splitting
- [x] Metadata extraction and preservation
- [x] Image extraction from documents
- [x] Table extraction and formatting

### LLM Integration
- [x] OpenAI (GPT-4, GPT-4o)
- [x] Anthropic (Claude 3)
- [x] Azure OpenAI
- [x] Ollama (local models)
- [x] Vision LLM for image description

### Vector Stores
- [x] PostgreSQL with pgvector
- [x] MongoDB Atlas Vector Search
- [x] HNSW and IVFFlat index support

### Search
- [x] Azure AI Search integration
- [x] Hybrid search (keyword + vector)
- [x] Semantic ranking

### Storage Backends
- [x] Local filesystem
- [x] AWS S3
- [x] Azure Blob Storage

### Batch Processing
- [x] Parallel document processing
- [x] Job queue with status tracking
- [x] Retry logic with exponential backoff
- [x] Progress callbacks

### SDK Documentation
- [x] Go SDK README (2100+ lines)
- [x] Python SDK README (2900+ lines)
- [x] Java SDK README (2200+ lines)
- [x] Real-world examples and patterns

---

## 🚧 Short Term (v1.1)

**Target: Q1 2026**

### Features
- [ ] **PDF parsing with Azure Document Intelligence** — Extract layout, tables, figures
- [ ] **Embedding generation** — Built-in OpenAI/Cohere embeddings
- [ ] **Pinecone integration** — Vector store support
- [ ] **Weaviate integration** — Additional vector store
- [ ] **CLI tool** — Command-line interface for all operations

### Improvements
- [ ] **Streaming responses** — Stream LLM outputs
- [ ] **Async Python SDK** — Full async/await support
- [ ] **Connection pooling** — Database connection management
- [ ] **Caching layer** — Redis/Memcached support

### Developer Experience
- [ ] **SDK versioning** — Semantic versioning across all SDKs
- [ ] **CI/CD pipelines** — Automated testing and releases
- [ ] **API documentation** — OpenAPI/Swagger specs

---

## 📋 Medium Term (v1.2)

**Target: Q2 2026**

### Features
- [ ] **OCR support** — Extract text from scanned documents
- [ ] **PowerPoint support** — PPTX → Markdown
- [ ] **HTML → Markdown** — Web page conversion
- [ ] **Custom templates** — User-defined output formats
- [ ] **Watermarks** — Add watermarks to PDFs
- [ ] **Digital signatures** — Sign generated PDFs

### Enterprise Features
- [ ] **Multi-tenancy** — Isolated workspaces
- [ ] **Rate limiting** — API usage controls
- [ ] **Audit logging** — Track all operations
- [ ] **RBAC** — Role-based access control

### Integrations
- [ ] **Slack** — Document processing bot
- [ ] **Microsoft Teams** — Integration app
- [ ] **Zapier/Make** — Workflow automation

---

## 🔮 Long Term (v2.0)

**Target: Q4 2026**

### Platform
- [ ] **Cloud-hosted version** — Managed DocFlow service
- [ ] **VS Code extension** — Convert documents from editor
- [ ] **Desktop app** — Cross-platform Electron/Tauri app
- [ ] **Browser extension** — Quick web page conversion

### Advanced RAG
- [ ] **Agentic retrieval** — Multi-step reasoning
- [ ] **Knowledge graphs** — Entity extraction and linking
- [ ] **Cross-lingual** — Multi-language document support
- [ ] **Query expansion** — Automatic query refinement

### AI Features
- [ ] **Document summarization** — Auto-generate summaries
- [ ] **Question answering** — Build Q&A over documents
- [ ] **Citation extraction** — Academic reference parsing
- [ ] **Auto-tagging** — ML-based document classification

---

## 💭 Maybe Someday

Ideas that might happen if there's demand:

- Mobile app (iOS/Android)
- LaTeX support
- Presentation mode (MD → slides)
- Version history and diffs
- Real-time collaborative editing
- Custom ML model training
- On-premise deployment package

---

## ❌ Won't Do

Things that are out of scope:

- Full word processor features
- Real-time collaboration (complex, many solutions exist)
- DRM/copy protection
- Paid tiers (keeping core open source)

---

## 📊 Feature Prioritization

Features are prioritized based on:

1. **Community demand** — Most requested features get priority
2. **Complexity vs value** — Quick wins over complex features
3. **Maintainability** — Must be testable and maintainable
4. **Alignment** — Must fit the project's purpose

---

## 💡 Want to Suggest Something?

Open an issue with the `feature-request` label. Include:

- What you want
- Why you need it
- How you'd use it
- Example use case

Good suggestions with clear use cases get prioritized.

---

## 🤝 Contributing to Roadmap Items

Want to work on something from this list? Great!

1. Check if there's an existing issue
2. If not, create one and mention you want to work on it
3. Wait for confirmation (to avoid duplicate work)
4. Start coding

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

<p align="center">
  <i>Last updated: January 2026</i>
</p>
