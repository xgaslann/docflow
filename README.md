# DocFlow

<p align="center">
  <b>Enterprise Document Processing & RAG Pipeline SDK</b><br>
  <i>Multi-format conversion • Semantic chunking • Vector stores • LLM integration</i>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Python-3.9+-3776AB?style=flat&logo=python" alt="Python">
  <img src="https://img.shields.io/badge/Java-17+-007396?style=flat&logo=openjdk" alt="Java">
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
</p>

---

## 🎯 Overview

DocFlow is a comprehensive document processing system with standalone SDKs for:

- **📄 Multi-format Conversion** — Markdown, PDF, DOCX, Excel, CSV, TXT
- **🧠 RAG Pipeline** — Semantic chunking, heading-aware splitting, metadata extraction
- **🔍 Vector Search** — PostgreSQL (pgvector), MongoDB Atlas, Azure AI Search
- **🤖 LLM Integration** — OpenAI, Anthropic, Azure OpenAI, Ollama
- **☁️ Cloud Storage** — Local, AWS S3, Azure Blob Storage

```
docflow/
├── app/                 # Web Application (Go + React)
├── sdks/                # Standalone SDKs
│   ├── go/              # Go SDK (2100+ lines docs)
│   ├── python/          # Python SDK (2900+ lines docs)
│   └── java/            # Java SDK (2200+ lines docs)
└── examples/            # Usage examples
```

---

## 🚀 Quick Start

### Go SDK

```bash
go get github.com/xgaslan/docflow/sdks/go@latest
```

```go
package main

import (
    "context"
    "github.com/xgaslan/docflow/sdks/go/docflow"
    "github.com/xgaslan/docflow/sdks/go/docflow/rag"
)

func main() {
    // Basic conversion
    converter := docflow.NewConverter()
    files := []docflow.MDFile{docflow.NewMDFile("doc.md", "# Hello World")}
    result, _ := converter.ConvertToPDF(context.Background(), files, docflow.ConvertOptions{})

    // RAG Pipeline
    cfg := rag.DefaultRAGConfig()
    cfg.ChunkSize = 1000
    cfg.ChunkingStrategy = "heading_aware"
    
    processor := rag.NewBatchProcessor(cfg)
    doc, _ := processor.ProcessFile("document.pdf")
    
    for _, chunk := range doc.Chunks {
        fmt.Printf("Chunk %d: %s\n", chunk.Index, chunk.Content[:100])
    }
}
```

### Python SDK

```bash
pip install git+https://github.com/xgaslan/docflow.git#subdirectory=sdks/python
```

```python
from docflow import Converter, MDFile
from docflow.rag import RAGProcessor, RAGConfig
from docflow.storage.vector import PostgresVectorStore

# Basic conversion
converter = Converter()
result = converter.convert_to_pdf([MDFile("doc.md", "# Hello World")])

# RAG Pipeline
config = RAGConfig(
    chunk_size=1000,
    chunking_strategy="heading_aware",
    extract_images=True,
    describe_images=True
)

processor = RAGProcessor(config)
doc = processor.process_file("document.pdf")

# Store in vector database
vector_store = PostgresVectorStore("postgresql://localhost/docflow")
vector_store.upsert(doc)

# Search
results = vector_store.search(query_embedding, top_k=5)
```

### Java SDK

```xml
<!-- Add JitPack repository -->
<repositories>
    <repository>
        <id>jitpack.io</id>
        <url>https://jitpack.io</url>
    </repository>
</repositories>

<dependency>
    <groupId>com.github.xgaslan.docflow</groupId>
    <artifactId>sdks-java</artifactId>
    <version>main-SNAPSHOT</version>
</dependency>
```

```java
import com.docflow.rag.*;
import com.docflow.storage.vector.*;
import com.docflow.search.*;

// RAG Pipeline
RAGConfig config = RAGConfig.defaultConfig();
config.setChunkSize(1000);
config.setDescribeImages(true);

RAGProcessor processor = new RAGProcessor(config);
RAGDocument doc = processor.processFile("document.pdf");

// Vector Store
PostgresVectorStore vectorStore = new PostgresVectorStore(postgresConfig);
vectorStore.initialize();
vectorStore.upsert(doc);

// Azure AI Search
AzureAISearch search = new AzureAISearch(endpoint, apiKey, indexName);
List<SearchResult> results = search.hybridSearch(query, vector, 10);
```

---

## ✨ Features

| Feature | Go | Python | Java |
|---------|:--:|:------:|:----:|
| **Core Conversion** |
| MD → PDF | ✅ | ✅ | ✅ |
| PDF → MD | ✅ | ✅ | ✅ |
| DOCX → MD | ✅ | ✅ | ✅ |
| Excel → MD | ✅ | ✅ | ✅ |
| CSV → MD | ✅ | ✅ | ✅ |
| **RAG Pipeline** |
| Semantic Chunking | ✅ | ✅ | ✅ |
| Heading-aware Split | ✅ | ✅ | ✅ |
| Image Extraction | ✅ | ✅ | ✅ |
| Table Extraction | ✅ | ✅ | ✅ |
| LLM Image Description | ✅ | ✅ | ✅ |
| Batch Processing | ✅ | ✅ | ✅ |
| **Storage** |
| Local | ✅ | ✅ | ✅ |
| AWS S3 | ✅ | ✅ | ✅ |
| Azure Blob | ✅ | ✅ | ✅ |
| **Vector Stores** |
| PostgreSQL (pgvector) | ✅ | ✅ | ✅ |
| MongoDB Atlas | ✅ | ✅ | ✅ |
| **Search** |
| Azure AI Search | ✅ | ✅ | ✅ |
| Hybrid Search | ✅ | ✅ | ✅ |
| **LLM Providers** |
| OpenAI | ✅ | ✅ | ✅ |
| Anthropic | ✅ | ✅ | ✅ |
| Azure OpenAI | ✅ | ✅ | ✅ |
| Ollama | ✅ | ✅ | ✅ |

---

## 📚 Documentation

Each SDK includes comprehensive documentation (2000+ lines):

- **[Go SDK](sdks/go/README.md)** — Concurrency patterns, goroutine pools, context handling
- **[Python SDK](sdks/python/README.md)** — Async support, type hints, dataclasses  
- **[Java SDK](sdks/java/README.md)** — Spring Boot integration, enterprise patterns
- **[Web App](app/README.md)** — Full-featured UI for document conversion

### Quick Links

| Topic | Go | Python | Java |
|-------|:--:|:------:|:----:|
| Installation | [📖](sdks/go/README.md#installation) | [📖](sdks/python/README.md#installation) | [📖](sdks/java/README.md#installation) |
| RAG Pipeline | [📖](sdks/go/README.md#rag-system) | [📖](sdks/python/README.md#rag-system) | [📖](sdks/java/README.md#rag-system) |
| Vector Stores | [📖](sdks/go/README.md#vector-stores) | [📖](sdks/python/README.md#storage-backends) | [📖](sdks/java/README.md#storage-backends) |
| Azure Pipeline | [📖](sdks/go/README.md#azure-enterprise-pipeline) | [📖](sdks/python/README.md#complete-azure-enterprise-pipeline) | [📖](sdks/java/README.md#azure-enterprise-pipeline) |

---

## 🏗️ Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                        DocFlow SDK                                │
├──────────────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐ │
│  │ Converter │  │ Extractor │  │ Template │  │ Format Converters │ │
│  │ (MD↔PDF) │  │ (PDF→MD) │  │ (Custom) │  │ (DOCX,Excel,CSV) │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────────┘ │
├──────────────────────────────────────────────────────────────────┤
│                        RAG System                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐ │
│  │  Chunker │  │ LLM Proc │  │  Batch   │  │  Image Describer │ │
│  │(Semantic)│  │(OpenAI..)│  │ Processor│  │  (Vision LLM)    │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────────┘ │
├──────────────────────────────────────────────────────────────────┤
│                     Storage Layer                                 │
│  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────────────────────┐ │
│  │ Local  │  │ AWS S3 │  │ Azure  │  │      Vector Stores     │ │
│  │Storage │  │Storage │  │  Blob  │  │  (Postgres, MongoDB)   │ │
│  └────────┘  └────────┘  └────────┘  └────────────────────────┘ │
├──────────────────────────────────────────────────────────────────┤
│                     Search & Retrieval                            │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │  Azure AI Search  │  Hybrid Search  │  Semantic Ranking     │ │
│  └─────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🌐 Web Application

For a full UI experience with drag-and-drop, live preview, and more:

```bash
cd app
docker-compose up
```

Open http://localhost:3000

---

## ⚙️ Environment Variables

```bash
# LLM Providers
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com
AZURE_OPENAI_API_KEY=...
OLLAMA_HOST=http://localhost:11434

# Azure Document Intelligence
AZURE_DOCUMENT_INTELLIGENCE_ENDPOINT=https://your-resource.cognitiveservices.azure.com
AZURE_DOCUMENT_INTELLIGENCE_KEY=...

# Storage
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
AZURE_STORAGE_CONNECTION_STRING=...

# Vector Stores
POSTGRES_CONNECTION_STRING=postgresql://user:pass@localhost:5432/docflow
MONGODB_URI=mongodb+srv://...

# Azure AI Search
AZURE_SEARCH_ENDPOINT=https://your-search.search.windows.net
AZURE_SEARCH_API_KEY=...
```

---

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

```bash
# Clone
git clone https://github.com/xgaslan/docflow.git
cd docflow

# Run tests
cd sdks/go && go test ./...
cd ../python && pytest
cd ../java && mvn test
```

---

## 📋 Roadmap

See [ROADMAP.md](ROADMAP.md) for planned features.

**Recent Completions:**
- ✅ Multi-format converters (DOCX, Excel, CSV)
- ✅ RAG pipeline with semantic chunking
- ✅ Vector store integrations (PostgreSQL, MongoDB)
- ✅ Azure AI Search support
- ✅ LLM integration (OpenAI, Anthropic, Azure, Ollama)
- ✅ Batch processing with job tracking
- ✅ Comprehensive documentation (7000+ lines across SDKs)

---

## 📄 License

MIT License - see [LICENSE](LICENSE)

---

<p align="center">
  <b>Built with ❤️ for the developer community</b>
</p>
