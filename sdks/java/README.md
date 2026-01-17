# DocFlow Java SDK: Complete Developer Manual

[![JitPack](https://jitpack.io/v/xgaslan/docflow.svg)](https://jitpack.io/#xgaslan/docflow)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Java Version](https://img.shields.io/badge/Java-17+-blue.svg)](https://adoptium.net/)
[![Maven Central](https://img.shields.io/maven-central/v/io.docflow/docflow-sdk.svg)](https://search.maven.org/artifact/io.docflow/docflow-sdk)

**DocFlow Java SDK** is the enterprise standard for Document Processing and RAG (Retrieval-Augmented Generation) on the JVM.

Built for **Spring Boot**, **Quarkus**, and **Jakarta EE** environments, it provides type-safe, robust pipelines for turning unstructured documents into AI-ready data.

This manual covers **every single feature, module, configuration option, and usage pattern**.

---

## 📚 Table of Contents

### Getting Started
1.  [Introduction](#-introduction)
2.  [Why Java for RAG?](#why-java-for-rag)
3.  [Architecture Overview](#architecture-overview)
4.  [Installation](#-installation)
5.  [Quick Start](#-quick-start)

### Core Modules
6.  [Converter](#-converter)
7.  [Extractor](#-extractor)
8.  [Template](#-template)
9.  [MarkdownParser](#-markdownparser)
10. [DocFlowClient](#-docflowclient)

### Format Converters
11. [CSV Converter](#-csv-converter)
12. [Excel Converter](#-excel-converter)
13. [DOCX Converter](#-docx-converter)
14. [TXT Converter](#-txt-converter)

### RAG System
15. [RAG Processor](#-rag-processor)
16. [RAG Chunker](#-rag-chunker)
17. [LLM Processor](#-llm-processor)

### Models
18. [Core Models](#-core-models)

### Storage Backends
19. [Local Storage](#-local-storage)
20. [AWS S3 Storage](#-aws-s3-storage)
21. [Azure Blob Storage](#-azure-blob-storage)

### Configuration
22. [RAGConfig](#-ragconfig)
23. [LLMConfig](#-llmconfig)
24. [ChunkingConfig](#-chunkingconfig)

### Advanced Topics
25. [Spring Boot Integration](#-spring-boot-integration)
26. [Error Handling](#-error-handling)
27. [Performance Optimization](#-performance-optimization)
28. [Azure Enterprise Pipeline](#-azure-enterprise-pipeline)
29. [Troubleshooting & FAQ](#-troubleshooting--faq)
30. [License](#-license)

---

# 🌟 Introduction

## Why Java for RAG?

| Python Problem | Java Solution |
|----------------|---------------|
| Experimental, hard to productionize | Enterprise-ready, battle-tested |
| Type errors at runtime | Compile-time safety |
| Threading complexity (GIL) | Native multi-threading |
| Deployment complexity | Container-ready JARs |
| No existing Spring integration | Native Spring Boot beans |

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           DocFlow Java SDK                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────┐    ┌───────────┐    ┌──────────┐    ┌─────────────────┐  │
│  │  INPUT   │───▶│ CONVERTER │───▶│   RAG    │───▶│    STORAGE      │  │
│  │          │    │           │    │          │    │                 │  │
│  │ PDF      │    │ Format    │    │ Chunker  │    │ Local/S3/Azure  │  │
│  │ DOCX     │    │ Detection │    │ LLM      │    │ Vector DBs      │  │
│  │ Excel    │    │ to MD     │    │          │    │                 │  │
│  └──────────┘    └───────────┘    └──────────┘    └─────────────────┘  │
│                                                                          │
│  Packages:                                                               │
│  ├── com.docflow/             Core classes                              │
│  ├── com.docflow.config/      Configuration beans                       │
│  ├── com.docflow.models/      POJOs and DTOs                            │
│  ├── com.docflow.rag/         RAG processing logic                      │
│  ├── com.docflow.formats/     Format converters                         │
│  └── com.docflow.storage/     Storage backends                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

# 📦 Installation

## Maven

Add repository and dependency:

```xml
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

## Gradle

```groovy
repositories {
    mavenCentral()
    maven { url 'https://jitpack.io' }
}

dependencies {
    implementation 'com.github.xgaslan.docflow:sdks-java:main-SNAPSHOT'
}
```

## Required Dependencies

The SDK uses these libraries (included transitively):

```xml
<!-- Apache POI for Office files -->
<dependency>
    <groupId>org.apache.poi</groupId>
    <artifactId>poi-ooxml</artifactId>
    <version>5.2.5</version>
</dependency>

<!-- PDFBox for PDF processing -->
<dependency>
    <groupId>org.apache.pdfbox</groupId>
    <artifactId>pdfbox</artifactId>
    <version>3.0.0</version>
</dependency>
```

---

# 🚀 Quick Start

## Example 1: Simple Conversion

```java
import com.docflow.Converter;
import com.docflow.models.ConvertResult;
import com.docflow.storage.LocalStorage;
import java.nio.file.Path;

public class QuickStart {
    public static void main(String[] args) {
        // Initialize
        LocalStorage storage = new LocalStorage("./output");
        Converter converter = new Converter(storage);
        
        // Convert
        ConvertResult result = converter.convertToMarkdown(Path.of("document.pdf"));
        
        System.out.println(result.getContent());
    }
}
```

## Example 2: RAG Pipeline

```java
import com.docflow.rag.RAGChunker;
import com.docflow.config.RAGConfig;
import com.docflow.models.Chunk;
import java.util.List;

public class RAGExample {
    public static void main(String[] args) {
        // Configure
        RAGConfig config = new RAGConfig();
        config.setChunkSize(800);
        config.setChunkOverlap(100);
        config.setChunkingStrategy(ChunkingStrategy.HEADING_AWARE);
        
        // Initialize
        RAGChunker chunker = new RAGChunker(config);
        
        // Process
        String markdown = "# Title\n\n## Section 1\nContent...";
        List<Chunk> chunks = chunker.chunk(markdown);
        
        for (Chunk chunk : chunks) {
            System.out.printf("[%d] %s%n", chunk.getIndex(), chunk.getContent());
        }
    }
}
```

## Example 3: Full RAG with LLM

```java
import com.docflow.rag.*;
import com.docflow.config.*;
import com.docflow.models.*;

public class FullRAGExample {
    public static void main(String[] args) {
        // LLM Configuration
        LLMConfig llmConfig = new LLMConfig();
        llmConfig.setProvider("openai");
        llmConfig.setModel("gpt-4o");
        llmConfig.setApiKey(System.getenv("OPENAI_API_KEY"));
        
        // RAG Configuration
        RAGConfig ragConfig = new RAGConfig();
        ragConfig.setChunkSize(1000);
        ragConfig.setExtractImages(true);
        ragConfig.setLlmConfig(llmConfig);
        
        // Process
        RAGProcessor processor = new RAGProcessor(ragConfig);
        RAGDocument doc = processor.processFile("handbook.pdf");
        
        System.out.println("Chunks: " + doc.getChunks().size());
        System.out.println("Images: " + doc.getImages().size());
    }
}
```

---

# 📄 Converter

**Location**: `com.docflow.Converter`

Converts Markdown to PDF and vice versa.

## Basic Usage

```java
import com.docflow.Converter;
import com.docflow.models.*;
import com.docflow.storage.LocalStorage;

public class ConverterExample {
    public static void main(String[] args) {
        Converter converter = new Converter(new LocalStorage("./output"));
        
        // Markdown to PDF
        List<MDFile> files = List.of(
            new MDFile("doc.md", "# Hello World\n\nContent here.")
        );
        
        PDFResult result = converter.convertToPdf(files);
        System.out.println("PDF URL: " + result.getDownloadUrl());
        
        // Any file to Markdown
        ConvertResult mdResult = converter.convertToMarkdown(Path.of("report.pdf"));
        System.out.println(mdResult.getContent());
    }
}
```

## With Options

```java
ConvertOptions options = new ConvertOptions();
options.setPageFormat("A4");
options.setMarginTop(20);
options.setMarginBottom(20);

PDFResult result = converter.convertToPdf(files, options);
```

---

# 📤 Extractor

**Location**: `com.docflow.Extractor`

Extracts content from PDFs with layout analysis.

```java
import com.docflow.Extractor;
import com.docflow.models.ExtractResult;

public class ExtractorExample {
    public static void main(String[] args) {
        Extractor extractor = new Extractor();
        
        ExtractResult result = extractor.extract(Path.of("report.pdf"));
        
        // Content
        System.out.println(result.getContent());
        
        // Images
        for (ExtractedImage img : result.getImages()) {
            System.out.println("Image: " + img.getFilename());
        }
        
        // Tables
        for (ExtractedTable table : result.getTables()) {
            System.out.println(table.getMarkdown());
        }
    }
}
```

---

# 🎨 Template

**Location**: `com.docflow.Template`

Template engine for dynamic document generation.

```java
import com.docflow.Template;
import java.util.Map;

public class TemplateExample {
    public static void main(String[] args) {
        Template engine = new Template();
        
        String template = "# Hello ${name}\n\nWelcome to ${company}!";
        
        Map<String, Object> data = Map.of(
            "name", "Alice",
            "company", "Acme Inc"
        );
        
        String result = engine.render(template, data);
        System.out.println(result);
    }
}
```

## From File

```java
Template engine = new Template("./templates");
String result = engine.renderFile("invoice.md.ftl", invoiceData);
```

---

# 📝 MarkdownParser

**Location**: `com.docflow.MarkdownParser`

Utilities for parsing Markdown content.

```java
import com.docflow.MarkdownParser;
import com.docflow.models.HeaderInfo;
import java.util.List;

public class ParserExample {
    public static void main(String[] args) {
        String content = """
            # Title
            ## Section 1
            Content here...
            ## Section 2
            More content...
            """;
        
        // Parse headers
        List<HeaderInfo> headers = MarkdownParser.parseHeaders(content);
        for (HeaderInfo h : headers) {
            System.out.printf("Level %d: %s (line %d)%n", 
                h.getLevel(), h.getText(), h.getLine());
        }
        
        // Extract section
        String section = MarkdownParser.extractSection(content, "Section 1");
        System.out.println(section);
    }
}
```

---

# 🔌 DocFlowClient

**Location**: `com.docflow.DocFlowClient`

HTTP client for DocFlow API.

```java
import com.docflow.DocFlowClient;

public class ClientExample {
    public static void main(String[] args) {
        DocFlowClient client = DocFlowClient.builder()
            .apiKey("your-api-key")
            .baseUrl("https://api.docflow.io")
            .timeout(30000)
            .build();
        
        // Health check
        boolean healthy = client.healthCheck();
        
        // Convert
        byte[] pdfBytes = client.convertToPdf(markdownContent);
    }
}
```

---

# 📊 CSV Converter

**Location**: `com.docflow.formats.CSVConverter`

```java
import com.docflow.formats.CSVConverter;
import com.docflow.models.ConvertResult;

public class CSVExample {
    public static void main(String[] args) {
        CSVConverter converter = new CSVConverter();
        
        // From file
        ConvertResult result = converter.toMarkdown(
            Files.readAllBytes(Path.of("data.csv")),
            "data.csv"
        );
        
        System.out.println(result.getContent());
        // | Name | Age |
        // |------|-----|
        // | Alice | 30 |
    }
}
```

## Options

```java
CSVConverter converter = new CSVConverter();
converter.setDelimiter(';');        // Semicolon separated
converter.setSkipEmptyRows(true);
converter.setMaxColumnWidth(50);
```

---

# 📗 Excel Converter

**Location**: `com.docflow.formats.ExcelConverter`

```java
import com.docflow.formats.ExcelConverter;
import com.docflow.models.ConvertResult;

public class ExcelExample {
    public static void main(String[] args) throws Exception {
        ExcelConverter converter = new ExcelConverter();
        converter.setIncludeAllSheets(true);
        
        byte[] excelBytes = Files.readAllBytes(Path.of("financials.xlsx"));
        ConvertResult result = converter.toMarkdown(excelBytes, "financials.xlsx");
        
        System.out.println(result.getContent());
        // # financials.xlsx
        // ## Sheet1
        // | Col1 | Col2 |
        // |------|------|
    }
}
```

---

# 📘 DOCX Converter

**Location**: `com.docflow.formats.DOCXConverter`

```java
import com.docflow.formats.DOCXConverter;
import com.docflow.models.ConvertResult;

public class DOCXExample {
    public static void main(String[] args) throws Exception {
        DOCXConverter converter = new DOCXConverter();
        converter.setExtractImages(true);
        converter.setPreserveFormatting(true);
        
        byte[] docxBytes = Files.readAllBytes(Path.of("document.docx"));
        ConvertResult result = converter.toMarkdown(docxBytes, "document.docx");
        
        // Extracted images
        for (ExtractedImage img : result.getImages()) {
            System.out.println("Found image: " + img.getFilename());
        }
    }
}
```

---

# 📄 TXT Converter

**Location**: `com.docflow.formats.TXTConverter`

```java
import com.docflow.formats.TXTConverter;
import com.docflow.models.ConvertResult;

public class TXTExample {
    public static void main(String[] args) throws Exception {
        TXTConverter converter = new TXTConverter();
        converter.setDetectHeaders(true);  // ALL CAPS as headers
        converter.setDetectLists(true);    // Numbered/bullet lists
        
        byte[] txtBytes = Files.readAllBytes(Path.of("notes.txt"));
        ConvertResult result = converter.toMarkdown(txtBytes, "notes.txt");
        
        System.out.println(result.getContent());
    }
}
```

---

# 🤖 RAG Processor

**Location**: `com.docflow.rag.RAGProcessor`

Main orchestrator for RAG pipelines (hypothetical, to be implemented).

```java
import com.docflow.rag.RAGProcessor;
import com.docflow.config.RAGConfig;
import com.docflow.config.LLMConfig;
import com.docflow.models.RAGDocument;

public class RAGProcessorExample {
    public static void main(String[] args) {
        // LLM Config
        LLMConfig llmConfig = new LLMConfig();
        llmConfig.setProvider("openai");
        llmConfig.setModel("gpt-4o");
        llmConfig.setApiKey(System.getenv("OPENAI_API_KEY"));
        
        // RAG Config
        RAGConfig config = new RAGConfig();
        config.setChunkSize(1000);
        config.setChunkOverlap(200);
        config.setExtractImages(true);
        config.setExtractTables(true);
        config.setLlmConfig(llmConfig);
        
        // Process
        RAGProcessor processor = new RAGProcessor(config);
        RAGDocument doc = processor.processFile("report.pdf");
        
        // Access data
        System.out.println("Document ID: " + doc.getId());
        System.out.println("Chunks: " + doc.getChunks().size());
        System.out.println("Images: " + doc.getImages().size());
        System.out.println("Tables: " + doc.getTables().size());
    }
}
```

---

# ✂️ RAG Chunker

**Location**: `com.docflow.rag.RAGChunker`

Intelligent text splitting.

```java
import com.docflow.rag.RAGChunker;
import com.docflow.config.RAGConfig;
import com.docflow.models.Chunk;
import java.util.List;

public class ChunkerExample {
    public static void main(String[] args) {
        RAGConfig config = new RAGConfig();
        config.setChunkSize(500);
        config.setChunkOverlap(50);
        config.setChunkingStrategy(ChunkingStrategy.HEADING_AWARE);
        
        RAGChunker chunker = new RAGChunker(config);
        
        String markdown = """
            # Document Title
            
            ## Section 1
            Content for section 1...
            
            ## Section 2
            Content for section 2...
            """;
        
        List<Chunk> chunks = chunker.chunk(markdown);
        
        for (Chunk chunk : chunks) {
            System.out.printf("Chunk %d: %s%n", 
                chunk.getIndex(), 
                chunk.getMetadata().getSectionTitle());
        }
    }
}
```

## Strategies

```java
// Simple: Fixed size
ChunkingStrategy.SIMPLE

// Heading-Aware: Respects H1, H2, H3 (Recommended)
ChunkingStrategy.HEADING_AWARE

// Semantic: Sentence boundaries
ChunkingStrategy.SEMANTIC
```

---

# 🧠 LLM Processor

**Location**: `com.docflow.rag.LLMProcessor` (to be implemented)

```java
import com.docflow.rag.LLMProcessor;
import com.docflow.config.LLMConfig;

public class LLMExample {
    public static void main(String[] args) {
        LLMConfig config = new LLMConfig();
        config.setProvider("openai");
        config.setModel("gpt-4o");
        config.setApiKey("sk-...");
        
        LLMProcessor processor = new LLMProcessor(config);
        
        // Describe image
        byte[] imageBytes = Files.readAllBytes(Path.of("chart.png"));
        String description = processor.describeImage(imageBytes, "Financial chart");
        
        // Summarize table
        String tableMd = "| A | B |\n|---|---|\n| 1 | 2 |";
        String summary = processor.summarizeTable(tableMd);
    }
}
```

## Supported Providers

```java
// OpenAI
config.setProvider("openai");
config.setModel("gpt-4o");

// Anthropic
config.setProvider("anthropic");
config.setModel("claude-3-opus");

// Ollama (Local)
config.setProvider("ollama");
config.setModel("llava");
config.setBaseUrl("http://localhost:11434");

// Azure OpenAI
config.setProvider("azure");
config.setBaseUrl("https://your.openai.azure.com/");
```

---

# 📦 Core Models

**Location**: `com.docflow.models`

## MDFile

```java
public class MDFile {
    private String name;
    private String content;
    
    public MDFile(String name, String content) { ... }
}
```

## ConvertResult

```java
public class ConvertResult {
    private String content;
    private String format;
    private DocumentMetadata metadata;
    private List<ExtractedImage> images;
    private List<ExtractedTable> tables;
}
```

## Chunk

```java
public class Chunk {
    private int index;
    private String content;
    private ChunkMetadata metadata;
}

public class ChunkMetadata {
    private List<String> headingPath;
    private String sectionTitle;
    private boolean hasCode;
    private boolean hasTable;
    private int startPos;
    private int endPos;
}
```

## RAGDocument

```java
public class RAGDocument {
    private String id;
    private String filename;
    private String content;
    private List<Chunk> chunks;
    private List<ExtractedImage> images;
    private List<ExtractedTable> tables;
    private DocumentMetadata metadata;
}
```

---

# 💾 Local Storage

**Location**: `com.docflow.storage.LocalStorage`

```java
import com.docflow.storage.LocalStorage;

public class LocalStorageExample {
    public static void main(String[] args) throws Exception {
        LocalStorage storage = new LocalStorage("./output");
        
        // Save
        String path = storage.save("result.pdf", pdfBytes);
        System.out.println(path);  // ./output/result.pdf
        
        // Read
        byte[] data = storage.get("result.pdf");
        
        // Check
        boolean exists = storage.exists("result.pdf");
        
        // Delete
        storage.delete("result.pdf");
        
        // List
        List<String> files = storage.list();
    }
}
```

---

# ☁️ AWS S3 Storage

**Location**: `com.docflow.storage.S3Storage`

```java
import com.docflow.storage.S3Storage;

public class S3Example {
    public static void main(String[] args) {
        S3Storage storage = new S3Storage(
            "my-bucket",
            "us-east-1",
            "AKIA...",
            "secret..."
        );
        
        String url = storage.save("documents/report.pdf", pdfBytes);
        System.out.println(url);  // s3://my-bucket/documents/report.pdf
    }
}
```

---

# 🔷 Azure Blob Storage

**Location**: `com.docflow.storage.AzureStorage`

```java
import com.docflow.storage.AzureStorage;

public class AzureExample {
    public static void main(String[] args) {
        AzureStorage storage = new AzureStorage(
            "mystorageaccount",
            "documents",
            "accountKey..."
        );
        
        // Or with connection string
        AzureStorage storage2 = AzureStorage.fromConnectionString(
            "DefaultEndpointsProtocol=https;AccountName=...",
            "documents"
        );
        
        String url = storage.save("report.pdf", pdfBytes);
    }
}
```

---

# ⚙️ RAGConfig

**Location**: `com.docflow.config.RAGConfig`

```java
public class RAGConfig {
    // Chunking
    private int chunkSize = 1000;
    private int chunkOverlap = 200;
    private ChunkingStrategy chunkingStrategy = ChunkingStrategy.HEADING_AWARE;
    
    // Extraction
    private boolean extractImages = true;
    private boolean extractTables = true;
    private boolean preserveMetadata = true;
    
    // LLM
    private LLMConfig llmConfig;
    
    // Getters and Setters...
}
```

---

# 🤖 LLMConfig

**Location**: `com.docflow.config.LLMConfig`

```java
public class LLMConfig {
    private String provider;      // openai, anthropic, ollama, azure
    private String model;         // gpt-4o, claude-3-opus
    private String apiKey;
    private String baseUrl;       // Optional override
    private double temperature = 0.0;
    private int maxTokens = 1000;
    private int timeout = 60000;  // milliseconds
    
    // Getters and Setters...
}
```

---

# 📐 ChunkingConfig

**Location**: `com.docflow.config.ChunkingConfig`

```java
public class ChunkingConfig {
    private int chunkSize = 1000;
    private int chunkOverlap = 200;
    private ChunkingStrategy strategy = ChunkingStrategy.HEADING_AWARE;
    private int minChunkSize = 100;
    private int maxChunkSize = 2000;
    private boolean respectCodeBlocks = true;
    private boolean respectTables = true;
}
```

---

# 🌱 Spring Boot Integration

## Configuration Properties

**application.yml**
```yaml
docflow:
  rag:
    chunk-size: 1000
    chunk-overlap: 200
    extract-images: true
  llm:
    provider: openai
    model: gpt-4o
    api-key: ${OPENAI_API_KEY}
```

## Configuration Class

```java
@Configuration
@ConfigurationProperties(prefix = "docflow")
public class DocFlowProperties {
    private RAGProperties rag = new RAGProperties();
    private LLMProperties llm = new LLMProperties();
    
    public static class RAGProperties {
        private int chunkSize = 1000;
        private int chunkOverlap = 200;
        private boolean extractImages = true;
        // Getters and Setters
    }
    
    public static class LLMProperties {
        private String provider;
        private String model;
        private String apiKey;
        // Getters and Setters
    }
}
```

## Bean Configuration

```java
@Configuration
public class DocFlowConfig {

    @Bean
    public RAGConfig ragConfig(DocFlowProperties props) {
        RAGConfig config = new RAGConfig();
        config.setChunkSize(props.getRag().getChunkSize());
        config.setChunkOverlap(props.getRag().getChunkOverlap());
        config.setExtractImages(props.getRag().isExtractImages());
        return config;
    }
    
    @Bean
    public LLMConfig llmConfig(DocFlowProperties props) {
        LLMConfig config = new LLMConfig();
        config.setProvider(props.getLlm().getProvider());
        config.setModel(props.getLlm().getModel());
        config.setApiKey(props.getLlm().getApiKey());
        return config;
    }
    
    @Bean
    public RAGProcessor ragProcessor(RAGConfig ragConfig, LLMConfig llmConfig) {
        ragConfig.setLlmConfig(llmConfig);
        return new RAGProcessor(ragConfig);
    }
}
```

## Service Layer

```java
@Service
public class DocumentIngestionService {
    
    private final RAGProcessor ragProcessor;
    private final VectorStore vectorStore;
    
    public DocumentIngestionService(RAGProcessor ragProcessor, VectorStore vectorStore) {
        this.ragProcessor = ragProcessor;
        this.vectorStore = vectorStore;
    }
    
    public IngestResult ingest(String filePath) {
        RAGDocument doc = ragProcessor.processFile(filePath);
        vectorStore.upsert(doc);
        
        return IngestResult.builder()
            .documentId(doc.getId())
            .chunksCreated(doc.getChunks().size())
            .build();
    }
}
```

## REST Controller

```java
@RestController
@RequestMapping("/api/documents")
public class DocumentController {
    
    private final DocumentIngestionService ingestionService;
    
    @PostMapping("/ingest")
    public ResponseEntity<IngestResult> ingest(@RequestParam String filePath) {
        IngestResult result = ingestionService.ingest(filePath);
        return ResponseEntity.ok(result);
    }
}
```

---

# ⚠️ Error Handling

## Custom Exceptions

```java
package com.docflow.exceptions;

public class DocFlowException extends RuntimeException { }
public class ConversionException extends DocFlowException { }
public class ExtractionException extends DocFlowException { }
public class LLMException extends DocFlowException { }
public class StorageException extends DocFlowException { }
```

## Usage

```java
try {
    RAGDocument doc = processor.processFile("doc.pdf");
} catch (ConversionException e) {
    logger.error("Conversion failed: {}", e.getMessage());
} catch (LLMException e) {
    logger.error("LLM API error: {}", e.getMessage());
} catch (DocFlowException e) {
    logger.error("DocFlow error: {}", e.getMessage());
}
```

## Spring Exception Handler

```java
@ControllerAdvice
public class DocFlowExceptionHandler {
    
    @ExceptionHandler(ConversionException.class)
    public ResponseEntity<ErrorResponse> handleConversionError(ConversionException e) {
        return ResponseEntity.badRequest()
            .body(new ErrorResponse("CONVERSION_ERROR", e.getMessage()));
    }
    
    @ExceptionHandler(LLMException.class)
    public ResponseEntity<ErrorResponse> handleLLMError(LLMException e) {
        return ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE)
            .body(new ErrorResponse("LLM_ERROR", e.getMessage()));
    }
}
```

---

# 🚀 Performance Optimization

## 1. Thread Pool for Batch Processing

```java
ExecutorService executor = Executors.newFixedThreadPool(16);

List<Future<RAGDocument>> futures = files.stream()
    .map(file -> executor.submit(() -> processor.processFile(file)))
    .collect(Collectors.toList());

List<RAGDocument> results = futures.stream()
    .map(f -> {
        try { return f.get(); } 
        catch (Exception e) { return null; }
    })
    .filter(Objects::nonNull)
    .collect(Collectors.toList());

executor.shutdown();
```

## 2. Connection Pooling

```java
// For HTTP clients
OkHttpClient client = new OkHttpClient.Builder()
    .connectionPool(new ConnectionPool(50, 5, TimeUnit.MINUTES))
    .build();
```

## 3. Caching

```java
@Cacheable("embeddings")
public float[] getEmbedding(String text) {
    return embeddingService.generate(text);
}
```

## 4. Async Processing

```java
@Async
public CompletableFuture<RAGDocument> processAsync(String filePath) {
    return CompletableFuture.completedFuture(processor.processFile(filePath));
}
```

---

# 🏢 Azure Enterprise Pipeline

Complete end-to-end Azure RAG pipeline in Java.

```java
package com.example.rag;

import com.docflow.rag.*;
import com.docflow.config.*;
import com.docflow.storage.AzureStorage;
import com.docflow.models.*;

import java.io.*;
import java.time.Instant;
import java.util.*;

public class AzureEnterprisePipeline {
    
    private final AzureStorage blobStorage;
    private final RAGProcessor ragProcessor;
    private final AzureSearchClient searchClient;
    private final EmbeddingService embeddingService;
    
    public AzureEnterprisePipeline() {
        // 1. Azure Blob Storage
        this.blobStorage = AzureStorage.fromConnectionString(
            System.getenv("AZURE_STORAGE_CONNECTION_STRING"),
            "documents"
        );
        
        // 2. LLM Config
        LLMConfig llmConfig = new LLMConfig();
        llmConfig.setProvider("openai");
        llmConfig.setModel("gpt-4o");
        llmConfig.setApiKey(System.getenv("OPENAI_API_KEY"));
        
        // 3. RAG Config
        RAGConfig ragConfig = new RAGConfig();
        ragConfig.setChunkSize(800);
        ragConfig.setChunkOverlap(100);
        ragConfig.setExtractImages(true);
        ragConfig.setExtractTables(true);
        ragConfig.setLlmConfig(llmConfig);
        
        this.ragProcessor = new RAGProcessor(ragConfig);
        this.searchClient = new AzureSearchClient(
            System.getenv("AZURE_SEARCH_ENDPOINT"),
            System.getenv("AZURE_SEARCH_KEY"),
            "enterprise-docs"
        );
        this.embeddingService = new OpenAIEmbeddingService(
            System.getenv("OPENAI_API_KEY")
        );
    }
    
    public Map<String, Object> ingestDocument(String filePath) throws Exception {
        String documentId = UUID.randomUUID().toString();
        String filename = new File(filePath).getName();
        String timestamp = Instant.now().toString();
        
        // Step 1: Read file
        byte[] fileBytes = Files.readAllBytes(Path.of(filePath));
        
        // Step 2: Store original in Blob
        String originalPath = String.format("originals/%s/%s", documentId, filename);
        String originalUrl = blobStorage.save(originalPath, fileBytes);
        
        // Step 3: Process with RAG
        RAGDocument doc = ragProcessor.process(fileBytes, filename);
        
        // Step 4: Store Markdown in Blob
        String mdPath = String.format("markdown/%s/%s.md", documentId, filename);
        String mdUrl = blobStorage.save(mdPath, doc.getContent().getBytes());
        
        // Step 5: Prepare search documents
        List<Map<String, Object>> searchDocs = new ArrayList<>();
        
        for (Chunk chunk : doc.getChunks()) {
            float[] embedding = embeddingService.embed(chunk.getContent());
            
            Map<String, Object> searchDoc = new HashMap<>();
            searchDoc.put("id", documentId + "_" + chunk.getIndex());
            searchDoc.put("document_id", documentId);
            searchDoc.put("chunk_index", chunk.getIndex());
            searchDoc.put("content", chunk.getContent());
            searchDoc.put("content_vector", embedding);
            searchDoc.put("filename", filename);
            searchDoc.put("original_url", originalUrl);
            searchDoc.put("markdown_url", mdUrl);
            searchDoc.put("heading_path", String.join(" > ", 
                chunk.getMetadata().getHeadingPath()));
            searchDoc.put("created_at", timestamp);
            
            searchDocs.add(searchDoc);
        }
        
        // Step 6: Index in Azure Search
        searchClient.uploadDocuments(searchDocs);
        
        return Map.of(
            "document_id", documentId,
            "filename", filename,
            "original_url", originalUrl,
            "markdown_url", mdUrl,
            "chunks_indexed", doc.getChunks().size(),
            "images_indexed", doc.getImages().size()
        );
    }
    
    public List<Map<String, Object>> search(String query, int top) {
        float[] queryVector = embeddingService.embed(query);
        return searchClient.hybridSearch(query, queryVector, top);
    }
    
    public static void main(String[] args) throws Exception {
        AzureEnterprisePipeline pipeline = new AzureEnterprisePipeline();
        
        // Ingest
        Map<String, Object> result = pipeline.ingestDocument("./contract.pdf");
        System.out.println("Ingested: " + result);
        
        // Search
        List<Map<String, Object>> hits = pipeline.search("payment terms", 5);
        for (Map<String, Object> hit : hits) {
            System.out.printf("Score: %s%nContent: %s%n---%n",
                hit.get("@search.score"),
                hit.get("content").toString().substring(0, 200));
        }
    }
}
```

---

# ❓ Troubleshooting & FAQ

**Q: `NoClassDefFoundError: org/apache/poi/...`**  
A: Add POI dependency to your pom.xml/build.gradle

**Q: OutOfMemoryError on large Excel files?**  
A: Increase heap size: `java -Xmx2G -jar app.jar`

**Q: PDF extraction returns empty?**  
A: PDF might be image-based. Enable OCR with Document Intelligence.

**Q: Rate limits with OpenAI?**  
A: Implement exponential backoff retry logic.

**Q: How to use with Spring WebFlux?**  
A: Wrap processor calls in `Mono.fromCallable()`

---

# 📖 Real-World Examples

## Example 1: Complete RAG Ingestion Service

```java
package com.example.rag;

import com.docflow.rag.*;
import com.docflow.config.*;
import com.docflow.storage.*;
import com.docflow.models.*;

import java.io.*;
import java.nio.file.*;
import java.time.Instant;
import java.util.*;

public class RAGIngestionService {
    
    private final RAGProcessor ragProcessor;
    private final AzureStorage blobStorage;
    private final VectorStore vectorStore;
    private final EmbeddingService embeddingService;
    
    public RAGIngestionService() {
        // 1. Azure Blob Storage
        this.blobStorage = AzureStorage.fromConnectionString(
            System.getenv("AZURE_STORAGE_CONNECTION_STRING"),
            "documents"
        );
        
        // 2. LLM Config
        LLMConfig llmConfig = new LLMConfig();
        llmConfig.setProvider("openai");
        llmConfig.setModel("gpt-4o");
        llmConfig.setApiKey(System.getenv("OPENAI_API_KEY"));
        
        // 3. RAG Config
        RAGConfig ragConfig = new RAGConfig();
        ragConfig.setChunkSize(800);
        ragConfig.setChunkOverlap(100);
        ragConfig.setChunkingStrategy(ChunkingStrategy.HEADING_AWARE);
        ragConfig.setExtractImages(true);
        ragConfig.setExtractTables(true);
        ragConfig.setLlmConfig(llmConfig);
        
        this.ragProcessor = new RAGProcessor(ragConfig);
        this.vectorStore = new PostgresVectorStore(System.getenv("DATABASE_URL"));
        this.embeddingService = new OpenAIEmbeddingService(System.getenv("OPENAI_API_KEY"));
    }
    
    public IngestResult ingest(String filePath) throws Exception {
        String documentId = UUID.randomUUID().toString();
        String filename = Paths.get(filePath).getFileName().toString();
        String timestamp = Instant.now().toString();
        
        // Step 1: Read file
        byte[] fileBytes = Files.readAllBytes(Paths.get(filePath));
        
        // Step 2: Store original in Blob
        String originalPath = String.format("originals/%s/%s", documentId, filename);
        String originalUrl = blobStorage.save(originalPath, fileBytes);
        
        // Step 3: Process with RAG
        RAGDocument doc = ragProcessor.process(fileBytes, filename);
        
        // Step 4: Store Markdown in Blob
        String mdPath = String.format("markdown/%s/%s.md", documentId, filename);
        String mdUrl = blobStorage.save(mdPath, doc.getContent().getBytes());
        
        // Step 5: Generate embeddings and store chunks
        List<VectorRecord> records = new ArrayList<>();
        
        for (Chunk chunk : doc.getChunks()) {
            float[] embedding = embeddingService.embed(chunk.getContent());
            
            VectorRecord record = new VectorRecord();
            record.setId(documentId + "_" + chunk.getIndex());
            record.setContent(chunk.getContent());
            record.setEmbedding(embedding);
            record.setMetadata(Map.of(
                "document_id", documentId,
                "filename", filename,
                "heading_path", String.join(" > ", chunk.getMetadata().getHeadingPath()),
                "created_at", timestamp
            ));
            
            records.add(record);
        }
        
        // Step 6: Batch upsert to vector store
        vectorStore.upsertBatch(records);
        
        return new IngestResult(
            documentId,
            filename,
            originalUrl,
            mdUrl,
            doc.getChunks().size(),
            doc.getImages().size()
        );
    }
    
    public List<SearchResult> search(String query, int topK) {
        float[] queryVector = embeddingService.embed(query);
        return vectorStore.search(queryVector, topK);
    }
    
    public static void main(String[] args) throws Exception {
        RAGIngestionService service = new RAGIngestionService();
        
        // Ingest
        IngestResult result = service.ingest("./contracts/agreement.pdf");
        System.out.println("Ingested: " + result);
        
        // Search
        List<SearchResult> hits = service.search("payment terms", 5);
        for (SearchResult hit : hits) {
            System.out.printf("Score: %.3f%n", hit.getScore());
            System.out.printf("Content: %s...%n", hit.getContent().substring(0, 200));
            System.out.println("---");
        }
    }
}
```

## Example 2: Bulk Document Processor

```java
package com.example.batch;

import com.docflow.rag.*;
import com.docflow.config.*;
import com.docflow.models.*;

import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.util.concurrent.*;
import java.util.stream.*;

public class BulkDocumentProcessor {
    
    private final RAGProcessor processor;
    private final ExecutorService executor;
    
    public BulkDocumentProcessor(int workers) {
        RAGConfig config = new RAGConfig();
        config.setChunkSize(800);
        config.setChunkOverlap(100);
        
        this.processor = new RAGProcessor(config);
        this.executor = Executors.newFixedThreadPool(workers);
    }
    
    public List<ProcessResult> processDirectory(String dirPath) throws Exception {
        // Collect files
        List<Path> files = Files.walk(Paths.get(dirPath))
            .filter(Files::isRegularFile)
            .filter(p -> {
                String name = p.getFileName().toString().toLowerCase();
                return name.endsWith(".pdf") || name.endsWith(".docx") || name.endsWith(".xlsx");
            })
            .collect(Collectors.toList());
        
        System.out.printf("Found %d files to process%n", files.size());
        
        // Submit tasks
        List<Future<ProcessResult>> futures = new ArrayList<>();
        long startTime = System.currentTimeMillis();
        
        for (Path file : files) {
            futures.add(executor.submit(() -> processFile(file)));
        }
        
        // Collect results with progress
        List<ProcessResult> results = new ArrayList<>();
        int completed = 0;
        
        for (Future<ProcessResult> future : futures) {
            try {
                ProcessResult result = future.get(5, TimeUnit.MINUTES);
                results.add(result);
            } catch (TimeoutException e) {
                results.add(ProcessResult.error("timeout", e.getMessage()));
            } catch (ExecutionException e) {
                results.add(ProcessResult.error("error", e.getCause().getMessage()));
            }
            
            completed++;
            double percent = (double) completed / files.size() * 100;
            System.out.printf("\rProgress: %.1f%% (%d/%d)", percent, completed, files.size());
        }
        
        long elapsed = System.currentTimeMillis() - startTime;
        System.out.printf("%n%nCompleted in %.2f seconds%n", elapsed / 1000.0);
        System.out.printf("Rate: %.2f docs/sec%n", files.size() / (elapsed / 1000.0));
        
        // Statistics
        long successful = results.stream().filter(r -> r.isSuccess()).count();
        long failed = results.size() - successful;
        int totalChunks = results.stream()
            .filter(ProcessResult::isSuccess)
            .mapToInt(ProcessResult::getChunkCount)
            .sum();
        
        System.out.printf("Successful: %d%n", successful);
        System.out.printf("Failed: %d%n", failed);
        System.out.printf("Total chunks: %d%n", totalChunks);
        
        return results;
    }
    
    private ProcessResult processFile(Path file) {
        try {
            RAGDocument doc = processor.processFile(file.toString());
            return ProcessResult.success(
                file.getFileName().toString(),
                doc.getChunks().size(),
                doc.getImages().size()
            );
        } catch (Exception e) {
            return ProcessResult.error(file.getFileName().toString(), e.getMessage());
        }
    }
    
    public void shutdown() {
        executor.shutdown();
        try {
            executor.awaitTermination(1, TimeUnit.MINUTES);
        } catch (InterruptedException e) {
            executor.shutdownNow();
        }
    }
    
    public static void main(String[] args) throws Exception {
        BulkDocumentProcessor processor = new BulkDocumentProcessor(16);
        
        try {
            List<ProcessResult> results = processor.processDirectory("./documents");
            
            // Print errors
            results.stream()
                .filter(r -> !r.isSuccess())
                .forEach(r -> System.out.printf("Error: %s - %s%n", r.getFilename(), r.getError()));
        } finally {
            processor.shutdown();
        }
    }
}

class ProcessResult {
    private final boolean success;
    private final String filename;
    private final int chunkCount;
    private final int imageCount;
    private final String error;
    
    private ProcessResult(boolean success, String filename, int chunkCount, int imageCount, String error) {
        this.success = success;
        this.filename = filename;
        this.chunkCount = chunkCount;
        this.imageCount = imageCount;
        this.error = error;
    }
    
    public static ProcessResult success(String filename, int chunkCount, int imageCount) {
        return new ProcessResult(true, filename, chunkCount, imageCount, null);
    }
    
    public static ProcessResult error(String filename, String error) {
        return new ProcessResult(false, filename, 0, 0, error);
    }
    
    // Getters...
    public boolean isSuccess() { return success; }
    public String getFilename() { return filename; }
    public int getChunkCount() { return chunkCount; }
    public int getImageCount() { return imageCount; }
    public String getError() { return error; }
}
```

## Example 3: Financial Report Analyzer

```java
package com.example.finance;

import com.docflow.rag.*;
import com.docflow.config.*;
import com.docflow.models.*;
import com.fasterxml.jackson.databind.ObjectMapper;

public class FinancialReportAnalyzer {
    
    private final RAGProcessor processor;
    private final LLMProcessor llm;
    private final ObjectMapper mapper = new ObjectMapper();
    
    public FinancialReportAnalyzer() {
        LLMConfig llmConfig = new LLMConfig();
        llmConfig.setProvider("openai");
        llmConfig.setModel("gpt-4o");
        llmConfig.setApiKey(System.getenv("OPENAI_API_KEY"));
        
        RAGConfig ragConfig = new RAGConfig();
        ragConfig.setExtractTables(true);
        ragConfig.setLlmConfig(llmConfig);
        
        this.processor = new RAGProcessor(ragConfig);
        this.llm = new LLMProcessor(llmConfig);
    }
    
    public FinancialAnalysis analyze(String filePath) throws Exception {
        // Process document
        RAGDocument doc = processor.processFile(filePath);
        
        // Build table context
        StringBuilder tableContext = new StringBuilder();
        for (ExtractedTable table : doc.getTables()) {
            tableContext.append(table.getMarkdown()).append("\n\n");
        }
        
        // Create analysis prompt
        String prompt = String.format("""
            Analyze this financial document and extract key metrics.
            
            Document Content:
            %s
            
            Tables Found:
            %s
            
            Return JSON with:
            - revenue: string
            - profit_margin: string
            - yoy_growth: string
            - key_risks: string[]
            - future_outlook: string
            """,
            doc.getContent().substring(0, Math.min(5000, doc.getContent().length())),
            tableContext.toString()
        );
        
        String response = llm.complete(prompt);
        return mapper.readValue(response, FinancialAnalysis.class);
    }
    
    public static void main(String[] args) throws Exception {
        FinancialReportAnalyzer analyzer = new FinancialReportAnalyzer();
        
        FinancialAnalysis analysis = analyzer.analyze("annual_report_2024.pdf");
        
        System.out.println("Revenue: " + analysis.getRevenue());
        System.out.println("Profit Margin: " + analysis.getProfitMargin());
        System.out.println("YoY Growth: " + analysis.getYoyGrowth());
        System.out.println("Key Risks:");
        for (String risk : analysis.getKeyRisks()) {
            System.out.println("  - " + risk);
        }
        System.out.println("Future Outlook: " + analysis.getFutureOutlook());
    }
}

class FinancialAnalysis {
    private String revenue;
    private String profitMargin;
    private String yoyGrowth;
    private List<String> keyRisks;
    private String futureOutlook;
    
    // Getters and Setters...
    public String getRevenue() { return revenue; }
    public String getProfitMargin() { return profitMargin; }
    public String getYoyGrowth() { return yoyGrowth; }
    public List<String> getKeyRisks() { return keyRisks; }
    public String getFutureOutlook() { return futureOutlook; }
}
```

## Example 4: Document Comparison Tool

```java
package com.example.compare;

import com.docflow.rag.*;
import com.docflow.config.*;
import com.docflow.models.*;

import java.util.*;
import java.util.stream.*;

public class DocumentComparator {
    
    private final RAGProcessor processor;
    
    public DocumentComparator() {
        RAGConfig config = new RAGConfig();
        config.setChunkingStrategy(ChunkingStrategy.HEADING_AWARE);
        this.processor = new RAGProcessor(config);
    }
    
    public List<SectionDiff> compare(String oldPath, String newPath) throws Exception {
        RAGDocument oldDoc = processor.processFile(oldPath);
        RAGDocument newDoc = processor.processFile(newPath);
        
        // Map sections by heading
        Map<String, String> oldSections = new HashMap<>();
        for (Chunk chunk : oldDoc.getChunks()) {
            String key = String.join(" > ", chunk.getMetadata().getHeadingPath());
            oldSections.put(key, chunk.getContent());
        }
        
        Map<String, String> newSections = new HashMap<>();
        for (Chunk chunk : newDoc.getChunks()) {
            String key = String.join(" > ", chunk.getMetadata().getHeadingPath());
            newSections.put(key, chunk.getContent());
        }
        
        // Find all sections
        Set<String> allSections = new HashSet<>();
        allSections.addAll(oldSections.keySet());
        allSections.addAll(newSections.keySet());
        
        // Compare
        List<SectionDiff> diffs = new ArrayList<>();
        
        for (String section : allSections) {
            String oldText = oldSections.get(section);
            String newText = newSections.get(section);
            
            DiffStatus status;
            if (oldText == null) {
                status = DiffStatus.ADDED;
            } else if (newText == null) {
                status = DiffStatus.REMOVED;
            } else if (oldText.equals(newText)) {
                status = DiffStatus.UNCHANGED;
            } else {
                status = DiffStatus.MODIFIED;
            }
            
            if (status != DiffStatus.UNCHANGED) {
                diffs.add(new SectionDiff(section, status, oldText, newText));
            }
        }
        
        return diffs;
    }
    
    public static void main(String[] args) throws Exception {
        DocumentComparator comparator = new DocumentComparator();
        
        List<SectionDiff> diffs = comparator.compare("contract_v1.docx", "contract_v2.docx");
        
        for (SectionDiff diff : diffs) {
            System.out.printf("[%s] %s%n", diff.getStatus(), diff.getSection());
            if (diff.getStatus() == DiffStatus.MODIFIED) {
                System.out.printf("  Old: %s...%n", truncate(diff.getOldText(), 100));
                System.out.printf("  New: %s...%n", truncate(diff.getNewText(), 100));
            }
            System.out.println();
        }
    }
    
    private static String truncate(String text, int maxLength) {
        if (text == null) return "";
        return text.length() > maxLength ? text.substring(0, maxLength) : text;
    }
}

enum DiffStatus { ADDED, REMOVED, MODIFIED, UNCHANGED }

class SectionDiff {
    private final String section;
    private final DiffStatus status;
    private final String oldText;
    private final String newText;
    
    public SectionDiff(String section, DiffStatus status, String oldText, String newText) {
        this.section = section;
        this.status = status;
        this.oldText = oldText;
        this.newText = newText;
    }
    
    public String getSection() { return section; }
    public DiffStatus getStatus() { return status; }
    public String getOldText() { return oldText; }
    public String getNewText() { return newText; }
}
```

---

# 🧪 Testing

## Unit Testing with JUnit 5

```java
package com.docflow.rag;

import com.docflow.config.*;
import com.docflow.models.*;
import org.junit.jupiter.api.*;
import static org.junit.jupiter.api.Assertions.*;

import java.util.List;

class RAGChunkerTest {
    
    private RAGChunker chunker;
    
    @BeforeEach
    void setUp() {
        RAGConfig config = new RAGConfig();
        config.setChunkSize(100);
        config.setChunkingStrategy(ChunkingStrategy.HEADING_AWARE);
        chunker = new RAGChunker(config);
    }
    
    @Test
    void testSimpleChunking() {
        String markdown = """
            # Title
            
            ## Section 1
            Content for section 1.
            
            ## Section 2
            Content for section 2.
            """;
        
        List<Chunk> chunks = chunker.chunk(markdown);
        
        assertTrue(chunks.size() >= 2, "Expected at least 2 chunks");
        
        for (Chunk chunk : chunks) {
            assertNotNull(chunk.getMetadata().getHeadingPath());
            assertFalse(chunk.getContent().isEmpty());
        }
    }
    
    @Test
    void testChunkMetadata() {
        String markdown = "# Document\n\n## Section\nContent here.";
        
        List<Chunk> chunks = chunker.chunk(markdown);
        
        Chunk chunk = chunks.get(0);
        assertEquals(0, chunk.getIndex());
        assertNotNull(chunk.getMetadata());
    }
    
    @Test
    void testEmptyInput() {
        List<Chunk> chunks = chunker.chunk("");
        assertTrue(chunks.isEmpty());
    }
}
```

## Integration Testing

```java
package com.docflow.integration;

import com.docflow.rag.*;
import com.docflow.config.*;
import com.docflow.models.*;
import org.junit.jupiter.api.*;
import org.junit.jupiter.api.condition.*;
import static org.junit.jupiter.api.Assertions.*;

@EnabledIfEnvironmentVariable(named = "INTEGRATION_TESTS", matches = "true")
class RAGProcessorIntegrationTest {
    
    private RAGProcessor processor;
    
    @BeforeEach
    void setUp() {
        RAGConfig config = new RAGConfig();
        config.setChunkSize(500);
        config.setExtractImages(true);
        processor = new RAGProcessor(config);
    }
    
    @Test
    void testProcessRealPDF() {
        RAGDocument doc = processor.processFile("src/test/resources/sample.pdf");
        
        assertNotNull(doc.getId());
        assertFalse(doc.getChunks().isEmpty());
        assertNotNull(doc.getContent());
    }
    
    @Test
    void testProcessRealDOCX() {
        RAGDocument doc = processor.processFile("src/test/resources/sample.docx");
        
        assertNotNull(doc.getId());
        assertTrue(doc.getChunks().size() > 0);
    }
}
```

## Mocking with Mockito

```java
package com.docflow.rag;

import com.docflow.config.*;
import org.junit.jupiter.api.*;
import org.junit.jupiter.api.extension.*;
import org.mockito.*;
import org.mockito.junit.jupiter.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class LLMProcessorTest {
    
    @Mock
    private OpenAIClient openAIClient;
    
    @InjectMocks
    private LLMProcessor processor;
    
    @Test
    void testDescribeImage() {
        byte[] imageBytes = new byte[]{1, 2, 3};
        when(openAIClient.callVisionAPI(any(), any()))
            .thenReturn("A chart showing revenue growth");
        
        String description = processor.describeImage(imageBytes, "Financial chart");
        
        assertEquals("A chart showing revenue growth", description);
        verify(openAIClient).callVisionAPI(any(), any());
    }
    
    @Test
    void testSummarizeTable() {
        String tableMarkdown = "| A | B |\n|---|---|\n| 1 | 2 |";
        when(openAIClient.complete(any()))
            .thenReturn("A simple table with two columns");
        
        String summary = processor.summarizeTable(tableMarkdown);
        
        assertNotNull(summary);
        verify(openAIClient).complete(contains("table"));
    }
}
```

---

# 🐳 Docker Deployment

## Dockerfile

```dockerfile
# Build stage
FROM maven:3.9-eclipse-temurin-17 AS builder

WORKDIR /app
COPY pom.xml .
RUN mvn dependency:go-offline

COPY src ./src
RUN mvn package -DskipTests

# Runtime stage
FROM eclipse-temurin:17-jre-alpine

WORKDIR /app
COPY --from=builder /app/target/*.jar app.jar

ENV JAVA_OPTS="-Xmx2G -Xms512M"
EXPOSE 8080

ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar app.jar"]
```

## Docker Compose

```yaml
version: '3.8'

services:
  docflow-app:
    build: .
    environment:
      - SPRING_PROFILES_ACTIVE=prod
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - AZURE_STORAGE_CONNECTION_STRING=${AZURE_STORAGE_CONNECTION_STRING}
      - SPRING_DATASOURCE_URL=jdbc:postgresql://db:5432/docflow
    ports:
      - "8080:8080"
    depends_on:
      - db
    deploy:
      replicas: 3
      resources:
        limits:
          memory: 2G

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
  name: docflow-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: docflow-service
  template:
    metadata:
      labels:
        app: docflow-service
    spec:
      containers:
      - name: app
        image: docflow:latest
        resources:
          requests:
            memory: "1Gi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "2"
        env:
        - name: JAVA_OPTS
          value: "-Xmx1536M -Xms512M"
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: docflow-secrets
              key: openai-key
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /actuator/health/liveness
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /actuator/health/readiness
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: docflow-service
spec:
  selector:
    app: docflow-service
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: docflow-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: docflow-service
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

## Micrometer Metrics

```java
import io.micrometer.core.instrument.*;
import org.springframework.stereotype.Component;

@Component
public class DocFlowMetrics {
    
    private final Counter documentsProcessed;
    private final Counter processingErrors;
    private final Timer processingDuration;
    private final DistributionSummary chunksPerDocument;
    
    public DocFlowMetrics(MeterRegistry registry) {
        this.documentsProcessed = Counter.builder("docflow.documents.processed")
            .description("Total documents processed")
            .register(registry);
        
        this.processingErrors = Counter.builder("docflow.documents.errors")
            .description("Total processing errors")
            .register(registry);
        
        this.processingDuration = Timer.builder("docflow.processing.duration")
            .description("Document processing duration")
            .register(registry);
        
        this.chunksPerDocument = DistributionSummary.builder("docflow.chunks.per.document")
            .description("Number of chunks created per document")
            .register(registry);
    }
    
    public void recordSuccess(int chunkCount, long durationMs) {
        documentsProcessed.increment();
        chunksPerDocument.record(chunkCount);
        processingDuration.record(java.time.Duration.ofMillis(durationMs));
    }
    
    public void recordError() {
        processingErrors.increment();
    }
}
```

## Structured Logging with Logback

**logback-spring.xml**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<configuration>
    <appender name="JSON" class="ch.qos.logback.core.ConsoleAppender">
        <encoder class="net.logstash.logback.encoder.LogstashEncoder">
            <includeMdcKeyName>document_id</includeMdcKeyName>
            <includeMdcKeyName>correlation_id</includeMdcKeyName>
        </encoder>
    </appender>
    
    <root level="INFO">
        <appender-ref ref="JSON"/>
    </root>
    
    <logger name="com.docflow" level="DEBUG"/>
</configuration>
```

## OpenTelemetry Tracing

```java
import io.opentelemetry.api.trace.*;
import io.opentelemetry.context.*;
import org.springframework.stereotype.Component;

@Component
public class TracedRAGProcessor {
    
    private final Tracer tracer;
    private final RAGProcessor processor;
    
    public TracedRAGProcessor(Tracer tracer, RAGProcessor processor) {
        this.tracer = tracer;
        this.processor = processor;
    }
    
    public RAGDocument process(String filePath) {
        Span span = tracer.spanBuilder("process_document")
            .setAttribute("document.path", filePath)
            .startSpan();
        
        try (Scope scope = span.makeCurrent()) {
            RAGDocument doc = processor.processFile(filePath);
            
            span.setAttribute("document.id", doc.getId());
            span.setAttribute("document.chunks", doc.getChunks().size());
            span.setAttribute("document.images", doc.getImages().size());
            
            return doc;
        } catch (Exception e) {
            span.recordException(e);
            span.setStatus(StatusCode.ERROR);
            throw e;
        } finally {
            span.end();
        }
    }
}
```

---

# 📚 Glossary

| Term | Definition |
|------|------------|
| **RAG** | Retrieval-Augmented Generation - Combining retrieval with LLM generation |
| **Chunk** | A semantic segment of text, optimized for embedding |
| **Embedding** | Vector representation of text for similarity search |
| **JPA** | Java Persistence API for database access |
| **Heading-Aware** | Chunking that respects document structure |
| **Spring Boot** | Java framework for microservices |

---

# 📜 License

MIT License
