# DocFlow Go SDK: Complete Developer Manual

[![Go Reference](https://pkg.go.dev/badge/github.com/xgaslan/docflow/sdks/go/docflow.svg)](https://pkg.go.dev/github.com/xgaslan/docflow/sdks/go/docflow)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/xgaslan/docflow/sdks/go)](https://goreportcard.com/report/github.com/xgaslan/docflow/sdks/go)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev/)

**DocFlow Go SDK** is a high-performance, concurrent document processing library optimized for **RAG (Retrieval-Augmented Generation)** pipelines and backend ingestion services.

Unlike Python wrappers, this SDK is designed for **Go-native concurrency**, making it ideal for high-throughput ingestion workers.

This manual covers **every single feature, module, configuration option, and usage pattern**.

---

## 📚 Table of Contents

### Getting Started
1.  [Introduction](#-introduction)
2.  [Why Go for RAG?](#why-go-for-rag)
3.  [Architecture Overview](#architecture-overview)
4.  [Installation](#-installation)
5.  [Quick Start](#-quick-start)

### Core Modules
6.  [Converter Module](#-converter-module)
7.  [Extractor Module](#-extractor-module)
8.  [Template Module](#-template-module)
9.  [Markdown Module](#-markdown-module)
10. [Types & Models](#-types--models)

### Format Converters
11. [CSV Converter](#-csv-converter)
12. [Excel Converter](#-excel-converter)
13. [DOCX Converter](#-docx-converter)
14. [TXT Converter](#-txt-converter)

### RAG System
15. [RAG Processor](#-rag-processor)
16. [Chunker](#-chunker)
17. [LLM Processor](#-llm-processor)
18. [Image Describer](#-image-describer)

### Batch Processing
19. [Batch Processor](#-batch-processor)

### Storage Backends
20. [Local Storage](#-local-storage)
21. [AWS S3 Storage](#-aws-s3-storage)
22. [Azure Blob Storage](#-azure-blob-storage)

### Configuration
23. [RAG Configuration](#-rag-configuration)
24. [LLM Configuration](#-llm-configuration)
25. [Batch Configuration](#-batch-configuration)

### Advanced Topics
26. [Concurrency Patterns](#-concurrency-patterns)
27. [Error Handling](#-error-handling)
28. [Performance Optimization](#-performance-optimization)
29. [Azure Enterprise Pipeline](#-azure-enterprise-pipeline)
30. [Troubleshooting & FAQ](#-troubleshooting--faq)
31. [License](#-license)

---

# 🌟 Introduction

## Why Go for RAG?

| Python Problem | Go Solution |
|----------------|-------------|
| GIL limits true parallelism | Native goroutines, no GIL |
| Complex Celery/Redis for workers | Simple goroutine pools |
| Runtime type errors | Compile-time safety |
| Heavy deployment (venv, pip) | Single static binary |
| Memory-hungry | Efficient memory usage |

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           DocFlow Go SDK                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────┐    ┌───────────┐    ┌──────────┐    ┌─────────────────┐  │
│  │  INPUT   │───▶│ CONVERTER │───▶│   RAG    │───▶│   BATCH         │  │
│  │          │    │           │    │PROCESSOR │    │   PROCESSOR     │  │
│  │ PDF      │    │ Format    │    │ Chunker  │    │ Worker Pool     │  │
│  │ DOCX     │    │ Detection │    │ LLM      │    │ Queues          │  │
│  │ Excel    │    │ to MD     │    │          │    │                 │  │
│  └──────────┘    └───────────┘    └──────────┘    └─────────────────┘  │
│                                                                          │
│  Packages:                                                               │
│  ├── docflow/          Core types, Converter, BatchProcessor            │
│  ├── docflow/config/   All configuration structs                        │
│  ├── docflow/rag/      RAGProcessor, Chunker, LLMProcessor              │
│  ├── docflow/formats/  CSV, Excel, DOCX, TXT converters                 │
│  └── docflow/storage/  Local, S3, Azure storage backends                │
└─────────────────────────────────────────────────────────────────────────┘
```

---

# 📦 Installation

```bash
go get github.com/xgaslan/docflow/sdks/go
```

Update your module:
```bash
go mod tidy
```

## Dependencies

The SDK uses pure Go libraries where possible:
- **PDF**: `ledongthuc/pdf`
- **Excel**: `xuri/excelize`
- **UUID**: `google/uuid`

---

# 🚀 Quick Start

## Example 1: Simple Conversion

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/xgaslan/docflow/sdks/go/docflow"
)

func main() {
    converter := docflow.NewConverter()
    
    result, err := converter.ConvertFileToMarkdown(context.Background(), "document.pdf")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(result.Content)
}
```

## Example 2: RAG Pipeline

```go
package main

import (
    "fmt"
    "log"

    "github.com/xgaslan/docflow/sdks/go/docflow/rag"
    "github.com/xgaslan/docflow/sdks/go/docflow/config"
)

func main() {
    cfg := config.RAGConfig{
        ChunkSize:        800,
        ChunkOverlap:     100,
        ChunkingStrategy: config.ChunkingStrategyHeadingAware,
        ExtractImages:    true,
    }

    processor := rag.NewRAGProcessor(cfg)
    doc, err := processor.ProcessFile("handbook.docx")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created %d chunks\n", len(doc.Chunks))
    for _, chunk := range doc.Chunks {
        fmt.Printf("[%v] %s...\n", chunk.Metadata.HeadingPath, chunk.Content[:50])
    }
}
```

## Example 3: Batch Processing

```go
package main

import (
    "fmt"
    "time"

    "github.com/xgaslan/docflow/sdks/go/docflow"
    "github.com/xgaslan/docflow/sdks/go/docflow/config"
)

func main() {
    ragCfg := config.DefaultRAGConfig()
    batchCfg := config.BatchConfig{
        MaxWorkers:  16,
        QueueSize:   1000,
        RetryFailed: true,
    }

    processor := docflow.NewBatchProcessor(ragCfg, batchCfg)
    
    files := []string{"doc1.pdf", "doc2.docx", "data.xlsx"}
    jobID, _ := processor.Enqueue(files)

    for {
        status, _ := processor.GetStatus(jobID)
        fmt.Printf("Progress: %d/%d\n", status.ProcessedFiles, status.TotalFiles)
        
        if status.Status == config.JobStatusCompleted {
            break
        }
        time.Sleep(time.Second)
    }

    results, _ := processor.GetResult(jobID)
    fmt.Printf("Processed %d documents\n", len(results))
}
```

---

# 📄 Converter Module

**Location**: `docflow/converter.go`

Handles Markdown ↔ PDF conversion.

## Basic Usage

```go
package main

import (
    "context"
    "github.com/xgaslan/docflow/sdks/go/docflow"
)

func main() {
    // With default options
    converter := docflow.NewConverter()
    
    // Convert file to markdown
    result, err := converter.ConvertFileToMarkdown(context.Background(), "doc.pdf")
    if err != nil {
        panic(err)
    }
    
    println(result.Content)
    println(result.Metadata.Title)
}
```

## With Options

```go
converter := docflow.NewConverter(
    docflow.WithStorage(myStorage),
    docflow.WithAPIKey("your-api-key"),
    docflow.WithBaseURL("https://custom.api.com"),
)
```

## Convert Bytes

```go
data, _ := os.ReadFile("document.pdf")
result, err := converter.ConvertToMarkdown(context.Background(), data, "document.pdf")
```

---

# 📤 Extractor Module

**Location**: `docflow/extractor.go`

Handles PDF extraction with layout analysis.

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow"
)

func main() {
    extractor := docflow.NewExtractor()
    
    result, err := extractor.Extract("scanned.pdf")
    if err != nil {
        panic(err)
    }
    
    // Full content
    println(result.Content)
    
    // Extracted images
    for _, img := range result.Images {
        println(img.Filename, img.Description)
    }
    
    // Extracted tables
    for _, table := range result.Tables {
        println(table.Markdown)
    }
}
```

---

# 🎨 Template Module

**Location**: `docflow/template.go`

Go template engine for dynamic document generation.

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow"
)

func main() {
    engine := docflow.NewTemplateEngine()
    
    template := `# Hello {{.Name}}
    
Welcome to {{.Company}}!`

    data := map[string]interface{}{
        "Name":    "Alice",
        "Company": "Acme Inc",
    }
    
    result, _ := engine.Render(template, data)
    println(result)
}
```

## From File

```go
engine := docflow.NewTemplateEngine(docflow.WithTemplateDir("./templates"))
result, _ := engine.RenderFile("invoice.md.tmpl", invoiceData)
```

---

# 📝 Markdown Module

**Location**: `docflow/markdown.go`

Utilities for parsing Markdown.

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow"
)

func main() {
    content := `# Title
## Section 1
Content here...
## Section 2
More content...`

    // Parse headers
    headers := docflow.ParseHeaders(content)
    for _, h := range headers {
        println(h.Level, h.Text, h.Line)
    }
    
    // Extract section
    section := docflow.ExtractSection(content, "Section 1")
    println(section)
}
```

---

# 📊 Types & Models

**Location**: `docflow/types.go`

## Core Types

```go
// MDFile represents a markdown file
type MDFile struct {
    Name    string
    Content string
}

// ConvertResult is returned from conversion
type ConvertResult struct {
    Content  string
    Format   string
    Metadata DocumentMetadata
    Images   []ExtractedImage
    Tables   []ExtractedTable
}

// RAGDocument is the main document container
type RAGDocument struct {
    ID       string
    Filename string
    Content  string
    Chunks   []Chunk
    Images   []ExtractedImage
    Tables   []ExtractedTable
    Metadata DocumentMetadata
}

// Chunk represents a text segment
type Chunk struct {
    Index    int
    Content  string
    Metadata ChunkMetadata
}

// ChunkMetadata contains context information
type ChunkMetadata struct {
    HeadingPath  []string
    SectionTitle string
    HasCode      bool
    HasTable     bool
    StartPos     int
    EndPos       int
}
```

---

# 📊 CSV Converter

**Location**: `docflow/formats/csv.go`

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow/formats"
)

func main() {
    converter := formats.NewCSVConverter()
    
    // From file
    result, _ := converter.ToMarkdown(nil, "data.csv")
    println(result.Content)
    
    // From bytes
    data := []byte("name,age\nAlice,30\nBob,25")
    result, _ = converter.ToMarkdown(data, "users.csv")
}
```

## Options

```go
converter := formats.NewCSVConverter()
converter.Delimiter = ';'      // Semicolon separated
converter.SkipEmptyRows = true
converter.MaxColumnWidth = 50
```

---

# 📗 Excel Converter

**Location**: `docflow/formats/excel.go`

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow/formats"
)

func main() {
    converter := formats.NewExcelConverter()
    converter.IncludeAllSheets = true
    
    result, _ := converter.ToMarkdown(nil, "financials.xlsx")
    
    // Output:
    // # financials.xlsx
    // ## Sheet1
    // | Col1 | Col2 |
    // |------|------|
    // | val1 | val2 |
}
```

---

# 📘 DOCX Converter

**Location**: `docflow/formats/docx.go`

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow/formats"
)

func main() {
    converter := formats.NewDOCXConverter()
    converter.ExtractImages = true
    converter.PreserveFormatting = true
    
    result, _ := converter.ToMarkdown(nil, "document.docx")
    
    // Access extracted images
    for _, img := range result.Images {
        println(img.Filename)
    }
}
```

---

# 📄 TXT Converter

**Location**: `docflow/formats/txt.go`

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow/formats"
)

func main() {
    converter := formats.NewTXTConverter()
    converter.DetectHeaders = true  // Auto-detect ALL CAPS as headers
    converter.DetectLists = true    // Auto-detect numbered/bullet lists
    
    result, _ := converter.ToMarkdown(nil, "notes.txt")
}
```

---

# 🤖 RAG Processor

**Location**: `docflow/rag/processor.go`

The main orchestrator for RAG pipelines.

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow/rag"
    "github.com/xgaslan/docflow/sdks/go/docflow/config"
)

func main() {
    cfg := config.RAGConfig{
        ChunkSize:        1000,
        ChunkOverlap:     200,
        ChunkingStrategy: config.ChunkingStrategyHeadingAware,
        ExtractImages:    true,
        ExtractTables:    true,
    }

    processor := rag.NewRAGProcessor(cfg)
    doc, _ := processor.ProcessFile("report.pdf")

    // Access data
    println("Document ID:", doc.ID)
    println("Chunks:", len(doc.Chunks))
    println("Images:", len(doc.Images))
    println("Tables:", len(doc.Tables))
}
```

## With LLM Enrichment

```go
llmCfg := config.LLMConfig{
    Provider: "openai",
    Model:    "gpt-4o",
    APIKey:   os.Getenv("OPENAI_API_KEY"),
}

cfg := config.RAGConfig{
    ExtractImages: true,
    LLMConfig:     &llmCfg,  // Enable image descriptions
}

processor := rag.NewRAGProcessor(cfg)
```

---

# ✂️ Chunker

**Location**: `docflow/rag/chunker.go`

Intelligent text splitting.

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow/rag"
    "github.com/xgaslan/docflow/sdks/go/docflow/config"
)

func main() {
    cfg := config.RAGConfig{
        ChunkSize:        500,
        ChunkOverlap:     50,
        ChunkingStrategy: config.ChunkingStrategyHeadingAware,
    }

    chunker := rag.NewChunker(cfg)
    
    markdown := `# Title
## Section 1
Content for section 1...
## Section 2
Content for section 2...`

    chunks := chunker.Chunk(markdown)
    
    for _, chunk := range chunks {
        println("Path:", chunk.Metadata.HeadingPath)
        println("Content:", chunk.Content[:50])
    }
}
```

## Strategies

```go
// Simple: Fixed size chunks
config.ChunkingStrategySimple

// Heading-Aware: Respects H1, H2, H3 boundaries (Recommended)
config.ChunkingStrategyHeadingAware

// Semantic: Uses sentence boundaries
config.ChunkingStrategySemantic
```

---

# 🧠 LLM Processor

**Location**: `docflow/rag/llm_processor.go`

Interface to LLM providers.

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow/rag"
    "github.com/xgaslan/docflow/sdks/go/docflow/config"
)

func main() {
    cfg := config.LLMConfig{
        Provider: "openai",
        Model:    "gpt-4o",
        APIKey:   "sk-...",
    }

    processor := rag.NewLLMProcessor(cfg)

    // Describe image
    imageBytes, _ := os.ReadFile("chart.png")
    description, _ := processor.DescribeImage(imageBytes, "Financial chart")
    println(description)

    // Summarize table
    tableMarkdown := "| A | B |\n|---|---|\n| 1 | 2 |"
    summary, _ := processor.SummarizeTable(tableMarkdown)
    println(summary)
}
```

## Providers

```go
// OpenAI
config.LLMConfig{Provider: "openai", Model: "gpt-4o"}

// Anthropic
config.LLMConfig{Provider: "anthropic", Model: "claude-3-opus"}

// Ollama (Local)
config.LLMConfig{Provider: "ollama", Model: "llava", BaseURL: "http://localhost:11434"}

// Azure OpenAI
config.LLMConfig{Provider: "azure", Model: "gpt-4", BaseURL: "https://your.openai.azure.com/"}
```

---

# 🖼️ Image Describer

**Location**: `docflow/rag/image_describer.go`

```go
describer := rag.NewImageDescriber(llmConfig)

description, _ := describer.Describe(imageBytes)
println(description)

// For RAG context
description, _ = describer.DescribeForRAG(
    imageBytes,
    "Q3 Financial Report",
    "As shown in the chart below...",
)
```

---

# 🏭 Batch Processor

**Location**: `docflow/batch_processor.go`

Worker pool pattern for high-throughput processing.

```go
package main

import (
    "fmt"
    "time"

    "github.com/xgaslan/docflow/sdks/go/docflow"
    "github.com/xgaslan/docflow/sdks/go/docflow/config"
)

func main() {
    ragCfg := config.DefaultRAGConfig()
    
    batchCfg := config.BatchConfig{
        MaxWorkers:     32,      // 32 concurrent goroutines
        QueueSize:      5000,    // Buffer size
        RetryFailed:    true,    // Auto-retry
        MaxRetries:     3,
        TimeoutPerFile: 300,     // 5 minutes
    }

    processor := docflow.NewBatchProcessor(ragCfg, batchCfg)

    // Queue files
    files := make([]string, 1000)
    for i := range files {
        files[i] = fmt.Sprintf("doc_%d.pdf", i)
    }
    
    jobID, _ := processor.Enqueue(files)
    fmt.Println("Job started:", jobID)

    // Monitor
    ticker := time.NewTicker(time.Second)
    for range ticker.C {
        status, _ := processor.GetStatus(jobID)
        fmt.Printf("\rProgress: %d/%d", status.ProcessedFiles, status.TotalFiles)
        
        if status.Status == config.JobStatusCompleted {
            break
        }
    }

    // Results
    results, _ := processor.GetResult(jobID)
    for _, doc := range results {
        if doc.Error != nil {
            fmt.Printf("Error: %s - %v\n", doc.Filename, doc.Error)
        } else {
            fmt.Printf("Success: %s - %d chunks\n", doc.Filename, len(doc.Chunks))
        }
    }
}
```

---

# 💾 Local Storage

**Location**: `docflow/storage/local.go`

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow/storage"
)

func main() {
    store, _ := storage.NewLocalStorage("./output")
    
    // Save
    path, _ := store.Save("result.pdf", pdfBytes)
    println(path)  // ./output/result.pdf
    
    // Read
    data, _ := store.Get("result.pdf")
    
    // Check
    exists := store.Exists("result.pdf")
    
    // Delete
    store.Delete("result.pdf")
    
    // List
    files, _ := store.List()
}
```

---

# ☁️ AWS S3 Storage

**Location**: `docflow/storage/s3.go`

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow/storage"
)

func main() {
    store, _ := storage.NewS3Storage(storage.S3Config{
        Bucket:    "my-bucket",
        Region:    "us-east-1",
        AccessKey: "AKIA...",
        SecretKey: "...",
        Prefix:    "documents/",
    })
    
    url, _ := store.Save("report.pdf", pdfBytes)
    println(url)  // s3://my-bucket/documents/report.pdf
}
```

---

# 🔷 Azure Blob Storage

**Location**: `docflow/storage/azure.go`

```go
package main

import (
    "github.com/xgaslan/docflow/sdks/go/docflow/storage"
)

func main() {
    store, _ := storage.NewAzureStorage(storage.AzureConfig{
        AccountName:   "mystorageaccount",
        AccountKey:    "...",
        ContainerName: "documents",
    })
    
    // Or with connection string
    store, _ = storage.NewAzureStorageFromConnectionString(
        "DefaultEndpointsProtocol=https;AccountName=...",
        "documents",
    )
    
    url, _ := store.Save("report.pdf", pdfBytes)
}
```

---

# ⚙️ RAG Configuration

**Location**: `docflow/config/rag.go`

```go
type RAGConfig struct {
    // Chunking
    ChunkSize        int              // Default: 1000
    ChunkOverlap     int              // Default: 200
    ChunkingStrategy ChunkingStrategy // Default: HeadingAware
    
    // Extraction
    ExtractImages    bool  // Default: true
    ExtractTables    bool  // Default: true
    PreserveMetadata bool  // Default: true
    RespectHeadings  bool  // Default: true
    
    // LLM (optional)
    LLMConfig *LLMConfig
}
```

## Default Config

```go
cfg := config.DefaultRAGConfig()
// ChunkSize: 1000, ChunkOverlap: 200, HeadingAware, ExtractImages: true
```

---

# 🤖 LLM Configuration

**Location**: `docflow/config/llm.go`

```go
type LLMConfig struct {
    Provider    string  // openai, anthropic, ollama, azure
    Model       string  // gpt-4o, claude-3-opus, llava
    APIKey      string
    BaseURL     string  // Optional override
    Temperature float64 // Default: 0.0
    MaxTokens   int     // Default: 1000
    Timeout     int     // Seconds
}
```

---

# 🔄 Batch Configuration

**Location**: `docflow/config/batch.go`

```go
type BatchConfig struct {
    MaxWorkers     int  // Concurrent workers
    QueueSize      int  // Buffer size
    RetryFailed    bool // Auto-retry failures
    MaxRetries     int  // Retry count
    TimeoutPerFile int  // Seconds per file
    FailFast       bool // Stop on first error
}
```

---

# ⚡ Concurrency Patterns

## Worker Pool

```go
func ProcessConcurrently(files []string, workers int) []RAGDocument {
    jobs := make(chan string, len(files))
    results := make(chan RAGDocument, len(files))
    
    // Start workers
    for w := 0; w < workers; w++ {
        go func() {
            for file := range jobs {
                doc, _ := processor.ProcessFile(file)
                results <- doc
            }
        }()
    }
    
    // Queue jobs
    for _, file := range files {
        jobs <- file
    }
    close(jobs)
    
    // Collect results
    var docs []RAGDocument
    for i := 0; i < len(files); i++ {
        docs = append(docs, <-results)
    }
    
    return docs
}
```

## Context Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

result, err := converter.ConvertFileToMarkdown(ctx, "huge.pdf")
if errors.Is(err, context.DeadlineExceeded) {
    log.Println("Conversion timed out")
}
```

## Error Groups

```go
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(context.Background())

for _, file := range files {
    file := file
    g.Go(func() error {
        _, err := processor.ProcessFile(file)
        return err
    })
}

if err := g.Wait(); err != nil {
    log.Fatal(err)
}
```

---

# ⚠️ Error Handling

## Custom Errors

```go
import "github.com/xgaslan/docflow/sdks/go/docflow/errors"

result, err := processor.ProcessFile("doc.pdf")
if err != nil {
    switch {
    case errors.Is(err, errors.ErrConversion):
        log.Println("Conversion failed")
    case errors.Is(err, errors.ErrLLM):
        log.Println("LLM API error")
    case errors.Is(err, errors.ErrStorage):
        log.Println("Storage error")
    default:
        log.Println("Unknown error:", err)
    }
}
```

## Retry Pattern

```go
func ProcessWithRetry(file string, maxRetries int) (*RAGDocument, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        doc, err := processor.ProcessFile(file)
        if err == nil {
            return doc, nil
        }
        lastErr = err
        time.Sleep(time.Second * time.Duration(i+1))
    }
    
    return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
```

---

# 🚀 Performance Optimization

## 1. Tune Worker Count

```go
// CPU-bound: Use CPU count
workers := runtime.NumCPU()

// IO-bound (API calls): Use more
workers := runtime.NumCPU() * 4
```

## 2. Disable Unnecessary Features

```go
cfg := config.RAGConfig{
    ExtractImages: false,  // Faster if you don't need images
    ExtractTables: false,  // Faster if you don't need tables
}
```

## 3. Use Connection Pooling

```go
// Reuse HTTP client
client := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 100,
    },
}
```

## 4. Streaming for Large Files

```go
// Process pages one by one
pages := extractor.Stream("huge.pdf")
for page := range pages {
    process(page)
}
```

---

# 🏢 Azure Enterprise Pipeline

Complete end-to-end Azure RAG pipeline in Go.

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/google/uuid"
    "github.com/xgaslan/docflow/sdks/go/docflow"
    "github.com/xgaslan/docflow/sdks/go/docflow/config"
    "github.com/xgaslan/docflow/sdks/go/docflow/rag"
    "github.com/xgaslan/docflow/sdks/go/docflow/storage"
)

type AzureEnterprisePipeline struct {
    blobStorage  *storage.AzureStorage
    ragProcessor *rag.RAGProcessor
    searchClient *AzureSearchClient // Your Azure Search wrapper
}

func NewAzureEnterprisePipeline() *AzureEnterprisePipeline {
    // 1. Azure Blob Storage
    blobStorage, _ := storage.NewAzureStorageFromConnectionString(
        os.Getenv("AZURE_STORAGE_CONNECTION_STRING"),
        "documents",
    )

    // 2. LLM Config
    llmCfg := config.LLMConfig{
        Provider: "openai",
        Model:    "gpt-4o",
        APIKey:   os.Getenv("OPENAI_API_KEY"),
    }

    // 3. RAG Config
    ragCfg := config.RAGConfig{
        ChunkSize:        800,
        ChunkOverlap:     100,
        ChunkingStrategy: config.ChunkingStrategyHeadingAware,
        ExtractImages:    true,
        ExtractTables:    true,
        LLMConfig:        &llmCfg,
    }

    return &AzureEnterprisePipeline{
        blobStorage:  blobStorage,
        ragProcessor: rag.NewRAGProcessor(ragCfg),
        searchClient: NewAzureSearchClient(), // Your implementation
    }
}

func (p *AzureEnterprisePipeline) IngestDocument(filePath string) (map[string]interface{}, error) {
    ctx := context.Background()
    documentID := uuid.New().String()
    filename := filepath.Base(filePath)
    timestamp := time.Now().UTC().Format(time.RFC3339)

    // Step 1: Read file
    fileBytes, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    // Step 2: Store original in Blob
    originalPath := fmt.Sprintf("originals/%s/%s", documentID, filename)
    originalURL, _ := p.blobStorage.Save(originalPath, fileBytes)

    // Step 3: Process with RAG
    doc, err := p.ragProcessor.Process(fileBytes, filename)
    if err != nil {
        return nil, err
    }

    // Step 4: Store Markdown in Blob
    mdPath := fmt.Sprintf("markdown/%s/%s.md", documentID, filename)
    mdURL, _ := p.blobStorage.Save(mdPath, []byte(doc.Content))

    // Step 5: Prepare search documents
    var searchDocs []map[string]interface{}
    
    for _, chunk := range doc.Chunks {
        embedding := p.generateEmbedding(chunk.Content)
        
        searchDoc := map[string]interface{}{
            "id":             fmt.Sprintf("%s_%d", documentID, chunk.Index),
            "document_id":    documentID,
            "chunk_index":    chunk.Index,
            "content":        chunk.Content,
            "content_vector": embedding,
            "filename":       filename,
            "original_url":   originalURL,
            "markdown_url":   mdURL,
            "heading_path":   strings.Join(chunk.Metadata.HeadingPath, " > "),
            "created_at":     timestamp,
        }
        searchDocs = append(searchDocs, searchDoc)
    }

    // Step 6: Index in Azure Search
    p.searchClient.UploadDocuments(ctx, searchDocs)

    return map[string]interface{}{
        "document_id":    documentID,
        "filename":       filename,
        "original_url":   originalURL,
        "markdown_url":   mdURL,
        "chunks_indexed": len(doc.Chunks),
        "images_indexed": len(doc.Images),
    }, nil
}

func (p *AzureEnterprisePipeline) Search(query string, top int) ([]map[string]interface{}, error) {
    queryVector := p.generateEmbedding(query)
    
    return p.searchClient.HybridSearch(context.Background(), query, queryVector, top)
}

func (p *AzureEnterprisePipeline) generateEmbedding(text string) []float32 {
    // Use OpenAI embeddings API
    // ... implementation
    return nil
}

func main() {
    pipeline := NewAzureEnterprisePipeline()
    
    // Ingest
    result, _ := pipeline.IngestDocument("./contracts/agreement.pdf")
    fmt.Printf("Ingested: %+v\n", result)
    
    // Search
    hits, _ := pipeline.Search("What are the payment terms?", 5)
    for _, hit := range hits {
        fmt.Printf("Score: %v\nContent: %s\n---\n", 
            hit["@search.score"], 
            hit["content"].(string)[:200])
    }
}
```

---

# ❓ Troubleshooting & FAQ

**Q: `undefined: config.RAGConfig`**  
A: Import correct path: `github.com/xgaslan/docflow/sdks/go/docflow/config`

**Q: Goroutine leaks?**  
A: Ensure BatchProcessor completes before main exits. Use `sync.WaitGroup` or channels.

**Q: High memory usage?**  
A: Reduce `MaxWorkers`. Processing 50 PDFs concurrently needs significant RAM.

**Q: PDF extraction fails?**  
A: Some PDFs are image-based (scanned). Use Document Intelligence OCR.

**Q: Rate limits with OpenAI?**  
A: Add retry logic with exponential backoff.

---

# � Real-World Examples

## Example 1: Complete RAG Ingestion Service

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"

    "github.com/google/uuid"
    "github.com/xgaslan/docflow/sdks/go/docflow"
    "github.com/xgaslan/docflow/sdks/go/docflow/config"
    "github.com/xgaslan/docflow/sdks/go/docflow/rag"
    "github.com/xgaslan/docflow/sdks/go/docflow/storage"
)

type IngestionService struct {
    processor    *rag.RAGProcessor
    vectorStore  VectorStore
    blobStorage  *storage.AzureStorage
}

func NewIngestionService() *IngestionService {
    // LLM Config
    llmCfg := config.LLMConfig{
        Provider: "openai",
        Model:    "gpt-4o",
        APIKey:   os.Getenv("OPENAI_API_KEY"),
    }

    // RAG Config
    ragCfg := config.RAGConfig{
        ChunkSize:        800,
        ChunkOverlap:     100,
        ChunkingStrategy: config.ChunkingStrategyHeadingAware,
        ExtractImages:    true,
        ExtractTables:    true,
        LLMConfig:        &llmCfg,
    }

    // Storage
    blobStorage, _ := storage.NewAzureStorageFromConnectionString(
        os.Getenv("AZURE_STORAGE_CONNECTION_STRING"),
        "documents",
    )

    return &IngestionService{
        processor:   rag.NewRAGProcessor(ragCfg),
        vectorStore: NewPostgresVectorStore(os.Getenv("DATABASE_URL")),
        blobStorage: blobStorage,
    }
}

func (s *IngestionService) Ingest(ctx context.Context, filePath string) (*IngestResult, error) {
    documentID := uuid.New().String()
    filename := filepath.Base(filePath)

    // Read file
    fileBytes, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    // Store original
    originalURL, _ := s.blobStorage.Save("originals/"+documentID+"/"+filename, fileBytes)

    // Process with RAG
    doc, err := s.processor.Process(fileBytes, filename)
    if err != nil {
        return nil, err
    }

    // Store markdown
    mdURL, _ := s.blobStorage.Save("markdown/"+documentID+"/"+filename+".md", []byte(doc.Content))

    // Generate embeddings and store
    for _, chunk := range doc.Chunks {
        embedding := s.generateEmbedding(chunk.Content)
        s.vectorStore.Upsert(ctx, VectorRecord{
            ID:        documentID + "_" + string(rune(chunk.Index)),
            Content:   chunk.Content,
            Embedding: embedding,
            Metadata: map[string]interface{}{
                "document_id":  documentID,
                "filename":     filename,
                "heading_path": chunk.Metadata.HeadingPath,
            },
        })
    }

    return &IngestResult{
        DocumentID:    documentID,
        Filename:      filename,
        OriginalURL:   originalURL,
        MarkdownURL:   mdURL,
        ChunksCreated: len(doc.Chunks),
        ImagesFound:   len(doc.Images),
    }, nil
}

func (s *IngestionService) Search(ctx context.Context, query string, topK int) ([]SearchResult, error) {
    queryVector := s.generateEmbedding(query)
    return s.vectorStore.Search(ctx, queryVector, topK)
}

func (s *IngestionService) generateEmbedding(text string) []float32 {
    // Call OpenAI embeddings API
    // ... implementation
    return nil
}

// HTTP Handler
func (s *IngestionService) HandleIngest(w http.ResponseWriter, r *http.Request) {
    filePath := r.URL.Query().Get("path")
    
    result, err := s.Ingest(r.Context(), filePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(result)
}

func main() {
    service := NewIngestionService()
    
    http.HandleFunc("/ingest", service.HandleIngest)
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Example 2: Bulk Document Processor

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"

    "github.com/xgaslan/docflow/sdks/go/docflow"
    "github.com/xgaslan/docflow/sdks/go/docflow/config"
)

func ProcessDirectory(dir string, workers int) {
    ragCfg := config.DefaultRAGConfig()
    batchCfg := config.BatchConfig{
        MaxWorkers:     workers,
        QueueSize:      1000,
        RetryFailed:    true,
        MaxRetries:     3,
        TimeoutPerFile: 300,
    }

    processor := docflow.NewBatchProcessor(ragCfg, batchCfg)

    // Collect files
    var files []string
    filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        ext := filepath.Ext(path)
        if ext == ".pdf" || ext == ".docx" || ext == ".xlsx" {
            files = append(files, path)
        }
        return nil
    })

    fmt.Printf("Found %d files to process\n", len(files))

    // Start job
    startTime := time.Now()
    jobID, _ := processor.Enqueue(files)

    // Monitor with progress bar
    for {
        status, _ := processor.GetStatus(jobID)
        percent := float64(status.ProcessedFiles) / float64(status.TotalFiles) * 100
        
        fmt.Printf("\r[%-50s] %.1f%% (%d/%d) | Failed: %d",
            progressBar(percent, 50),
            percent,
            status.ProcessedFiles,
            status.TotalFiles,
            status.FailedFiles,
        )

        if status.Status == config.JobStatusCompleted {
            break
        }
        time.Sleep(500 * time.Millisecond)
    }

    elapsed := time.Since(startTime)
    fmt.Printf("\n\nCompleted in %v\n", elapsed)
    fmt.Printf("Rate: %.2f docs/sec\n", float64(len(files))/elapsed.Seconds())

    // Get results
    results, _ := processor.GetResult(jobID)
    
    var totalChunks, totalImages int
    var errors []string
    
    for _, doc := range results {
        if doc.Error != nil {
            errors = append(errors, fmt.Sprintf("%s: %v", doc.Filename, doc.Error))
        } else {
            totalChunks += len(doc.Chunks)
            totalImages += len(doc.Images)
        }
    }

    fmt.Printf("\nTotal chunks: %d\n", totalChunks)
    fmt.Printf("Total images: %d\n", totalImages)
    fmt.Printf("Errors: %d\n", len(errors))
    
    for _, err := range errors {
        fmt.Println("  -", err)
    }
}

func progressBar(percent float64, width int) string {
    filled := int(percent / 100 * float64(width))
    bar := ""
    for i := 0; i < width; i++ {
        if i < filled {
            bar += "█"
        } else {
            bar += "░"
        }
    }
    return bar
}

func main() {
    ProcessDirectory("./documents", 16)
}
```

## Example 3: Financial Report Analyzer

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/xgaslan/docflow/sdks/go/docflow/config"
    "github.com/xgaslan/docflow/sdks/go/docflow/rag"
)

type FinancialAnalysis struct {
    Revenue         string   `json:"revenue"`
    ProfitMargin    string   `json:"profit_margin"`
    YoYGrowth       string   `json:"yoy_growth"`
    KeyRisks        []string `json:"key_risks"`
    FutureOutlook   string   `json:"future_outlook"`
}

func AnalyzeFinancialReport(filePath string) (*FinancialAnalysis, error) {
    llmCfg := config.LLMConfig{
        Provider: "openai",
        Model:    "gpt-4o",
        APIKey:   os.Getenv("OPENAI_API_KEY"),
    }

    ragCfg := config.RAGConfig{
        ExtractTables: true,
        LLMConfig:     &llmCfg,
    }

    processor := rag.NewRAGProcessor(ragCfg)
    doc, err := processor.ProcessFile(filePath)
    if err != nil {
        return nil, err
    }

    // Build context from tables
    var tableContext string
    for _, table := range doc.Tables {
        tableContext += table.Markdown + "\n\n"
    }

    // Analyze with LLM
    llm := rag.NewLLMProcessor(llmCfg)
    
    prompt := fmt.Sprintf(`
Analyze this financial document and extract key metrics.

Document Content:
%s

Tables Found:
%s

Return JSON with:
- revenue: string
- profit_margin: string  
- yoy_growth: string
- key_risks: []string
- future_outlook: string
`, doc.Content[:min(5000, len(doc.Content))], tableContext)

    response, err := llm.Complete(prompt)
    if err != nil {
        return nil, err
    }

    var analysis FinancialAnalysis
    json.Unmarshal([]byte(response), &analysis)

    return &analysis, nil
}

func main() {
    analysis, err := AnalyzeFinancialReport("annual_report.pdf")
    if err != nil {
        panic(err)
    }
    
    output, _ := json.MarshalIndent(analysis, "", "  ")
    fmt.Println(string(output))
}
```

## Example 4: Document Comparison Tool

```go
package main

import (
    "fmt"
    "strings"

    "github.com/xgaslan/docflow/sdks/go/docflow/config"
    "github.com/xgaslan/docflow/sdks/go/docflow/rag"
)

type SectionDiff struct {
    Section string
    Status  string // "added", "removed", "modified", "unchanged"
    OldText string
    NewText string
}

func CompareDocuments(oldPath, newPath string) ([]SectionDiff, error) {
    cfg := config.RAGConfig{
        ChunkingStrategy: config.ChunkingStrategyHeadingAware,
    }

    processor := rag.NewRAGProcessor(cfg)

    oldDoc, err := processor.ProcessFile(oldPath)
    if err != nil {
        return nil, err
    }

    newDoc, err := processor.ProcessFile(newPath)
    if err != nil {
        return nil, err
    }

    // Map sections by heading
    oldSections := make(map[string]string)
    for _, chunk := range oldDoc.Chunks {
        key := strings.Join(chunk.Metadata.HeadingPath, " > ")
        oldSections[key] = chunk.Content
    }

    newSections := make(map[string]string)
    for _, chunk := range newDoc.Chunks {
        key := strings.Join(chunk.Metadata.HeadingPath, " > ")
        newSections[key] = chunk.Content
    }

    // Compare
    var diffs []SectionDiff
    allSections := make(map[string]bool)
    
    for k := range oldSections {
        allSections[k] = true
    }
    for k := range newSections {
        allSections[k] = true
    }

    for section := range allSections {
        oldText, inOld := oldSections[section]
        newText, inNew := newSections[section]

        var status string
        switch {
        case !inOld:
            status = "added"
        case !inNew:
            status = "removed"
        case oldText == newText:
            status = "unchanged"
        default:
            status = "modified"
        }

        if status != "unchanged" {
            diffs = append(diffs, SectionDiff{
                Section: section,
                Status:  status,
                OldText: oldText,
                NewText: newText,
            })
        }
    }

    return diffs, nil
}

func main() {
    diffs, _ := CompareDocuments("contract_v1.docx", "contract_v2.docx")
    
    for _, diff := range diffs {
        fmt.Printf("[%s] %s\n", diff.Status, diff.Section)
        if diff.Status == "modified" {
            fmt.Printf("  Old: %s...\n", diff.OldText[:min(100, len(diff.OldText))])
            fmt.Printf("  New: %s...\n", diff.NewText[:min(100, len(diff.NewText))])
        }
        fmt.Println()
    }
}
```

---

# 🧪 Testing

## Unit Testing

```go
package docflow_test

import (
    "testing"

    "github.com/xgaslan/docflow/sdks/go/docflow"
    "github.com/xgaslan/docflow/sdks/go/docflow/config"
    "github.com/xgaslan/docflow/sdks/go/docflow/rag"
)

func TestChunker_SimpleMarkdown(t *testing.T) {
    cfg := config.RAGConfig{
        ChunkSize:        100,
        ChunkingStrategy: config.ChunkingStrategyHeadingAware,
    }

    chunker := rag.NewChunker(cfg)

    markdown := `# Title

## Section 1
Content for section 1.

## Section 2
Content for section 2.`

    chunks := chunker.Chunk(markdown)

    if len(chunks) < 2 {
        t.Errorf("Expected at least 2 chunks, got %d", len(chunks))
    }

    for _, chunk := range chunks {
        if len(chunk.Metadata.HeadingPath) == 0 {
            t.Error("Expected heading path in metadata")
        }
    }
}

func TestConverter_DetectsFormat(t *testing.T) {
    converter := docflow.NewConverter()

    // Test that it doesn't panic on unknown format
    _, err := converter.ConvertFileToMarkdown(nil, "test.xyz")
    if err == nil {
        t.Error("Expected error for unknown format")
    }
}
```

## Integration Testing

```go
package integration_test

import (
    "os"
    "testing"

    "github.com/xgaslan/docflow/sdks/go/docflow/config"
    "github.com/xgaslan/docflow/sdks/go/docflow/rag"
)

func TestRAGProcessor_RealPDF(t *testing.T) {
    if os.Getenv("INTEGRATION_TEST") != "true" {
        t.Skip("Skipping integration test")
    }

    cfg := config.RAGConfig{
        ChunkSize:     500,
        ExtractImages: true,
    }

    processor := rag.NewRAGProcessor(cfg)
    
    doc, err := processor.ProcessFile("testdata/sample.pdf")
    if err != nil {
        t.Fatalf("Failed to process: %v", err)
    }

    if len(doc.Chunks) == 0 {
        t.Error("Expected chunks to be created")
    }
    
    if doc.ID == "" {
        t.Error("Expected document ID")
    }
}
```

## Benchmark Testing

```go
func BenchmarkChunker(b *testing.B) {
    cfg := config.RAGConfig{ChunkSize: 500}
    chunker := rag.NewChunker(cfg)
    
    markdown := strings.Repeat("# Section\n\nParagraph content here.\n\n", 100)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        chunker.Chunk(markdown)
    }
}

func BenchmarkBatchProcessor(b *testing.B) {
    // Benchmark with different worker counts
    workerCounts := []int{1, 4, 8, 16, 32}
    
    for _, workers := range workerCounts {
        b.Run(fmt.Sprintf("workers-%d", workers), func(b *testing.B) {
            cfg := config.BatchConfig{MaxWorkers: workers}
            processor := docflow.NewBatchProcessor(config.DefaultRAGConfig(), cfg)
            
            files := make([]string, 100)
            for i := range files {
                files[i] = fmt.Sprintf("testdata/doc_%d.pdf", i%10)
            }

            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                jobID, _ := processor.Enqueue(files)
                for {
                    status, _ := processor.GetStatus(jobID)
                    if status.Status == config.JobStatusCompleted {
                        break
                    }
                }
            }
        })
    }
}
```

---

# 🐳 Docker Deployment

## Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /docflow-worker ./cmd/worker

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /docflow-worker .

ENV TZ=UTC
EXPOSE 8080

ENTRYPOINT ["/app/docflow-worker"]
```

## Docker Compose

```yaml
version: '3.8'

services:
  docflow-worker:
    build: .
    environment:
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - AZURE_STORAGE_CONNECTION_STRING=${AZURE_STORAGE_CONNECTION_STRING}
      - DATABASE_URL=postgres://postgres:postgres@db:5432/docflow
    ports:
      - "8080:8080"
    depends_on:
      - db
    deploy:
      replicas: 3
      resources:
        limits:
          memory: 2G
          cpus: '2'

  db:
    image: pgvector/pgvector:pg16
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: docflow
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"

volumes:
  pgdata:
```

## Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docflow-worker
spec:
  replicas: 5
  selector:
    matchLabels:
      app: docflow-worker
  template:
    metadata:
      labels:
        app: docflow-worker
    spec:
      containers:
      - name: worker
        image: docflow:latest
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "2"
        env:
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: docflow-secrets
              key: openai-key
        - name: AZURE_STORAGE_CONNECTION_STRING
          valueFrom:
            secretKeyRef:
              name: docflow-secrets
              key: azure-storage
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: docflow-worker
spec:
  selector:
    app: docflow-worker
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: docflow-worker-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: docflow-worker
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

---

# 📊 Monitoring & Observability

## Prometheus Metrics

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    DocsProcessed = promauto.NewCounter(prometheus.CounterOpts{
        Name: "docflow_documents_processed_total",
        Help: "Total number of documents processed",
    })

    ProcessingDuration = promauto.NewHistogram(prometheus.HistogramOpts{
        Name:    "docflow_processing_duration_seconds",
        Help:    "Time taken to process documents",
        Buckets: prometheus.ExponentialBuckets(0.1, 2, 10),
    })

    ChunksCreated = promauto.NewCounter(prometheus.CounterOpts{
        Name: "docflow_chunks_created_total",
        Help: "Total number of chunks created",
    })

    ProcessingErrors = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "docflow_processing_errors_total",
        Help: "Total number of processing errors",
    }, []string{"error_type"})
)
```

## Structured Logging

```go
package main

import (
    "log/slog"
    "os"
)

func main() {
    // JSON logging for production
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)

    // Usage
    slog.Info("Processing document",
        slog.String("filename", "report.pdf"),
        slog.Int("chunks", 15),
        slog.Duration("duration", elapsed),
    )

    slog.Error("Processing failed",
        slog.String("filename", "bad.pdf"),
        slog.Any("error", err),
    )
}
```

## OpenTelemetry Tracing

```go
package tracing

import (
    "context"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("docflow")

func ProcessWithTracing(ctx context.Context, filePath string) (*RAGDocument, error) {
    ctx, span := tracer.Start(ctx, "ProcessDocument",
        trace.WithAttributes(
            attribute.String("document.path", filePath),
        ),
    )
    defer span.End()

    doc, err := processor.ProcessFile(filePath)
    if err != nil {
        span.RecordError(err)
        return nil, err
    }

    span.SetAttributes(
        attribute.Int("document.chunks", len(doc.Chunks)),
        attribute.Int("document.images", len(doc.Images)),
    )

    return doc, nil
}
```

---

# 📚 Glossary

| Term | Definition |
|------|------------|
| **RAG** | Retrieval-Augmented Generation - Combining retrieval with LLM generation |
| **Chunk** | A semantic segment of text, optimized for embedding |
| **Embedding** | Vector representation of text for similarity search |
| **Goroutine** | Lightweight thread managed by Go runtime |
| **Context** | Go's context.Context for cancellation and timeouts |
| **Heading-Aware** | Chunking that respects document structure |

---

# �📜 License

MIT License
