# DocFlow Python SDK: Complete Developer Manual

[![PyPI version](https://badge.fury.io/py/docflow-client.svg)](https://badge.fury.io/py/docflow-client)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Python 3.9+](https://img.shields.io/badge/python-3.9+-blue.svg)](https://www.python.org/downloads/)
[![Code Coverage](https://codecov.io/gh/xgaslan/docflow/branch/main/graph/badge.svg)](https://codecov.io/gh/xgaslan/docflow)
[![Documentation](https://readthedocs.org/projects/docflow/badge/?version=latest)](https://docflow.readthedocs.io)

**DocFlow** is the most comprehensive document processing, RAG (Retrieval-Augmented Generation), and ingestion SDK for Python. It transforms messy, unstructured documents into clean, semantic, AI-ready data.

This manual covers **every single feature, module, configuration option, and usage pattern** available in the SDK.

---

## 📚 Table of Contents

### Getting Started
1.  [Introduction](#-introduction)
2.  [Why DocFlow?](#why-docflow)
3.  [Architecture Overview](#architecture-overview)
4.  [Installation](#-installation)
5.  [Quick Start](#-quick-start)

### Core Modules
6.  [Converter Module](#-converter-module)
7.  [Extractor Module](#-extractor-module)
8.  [Template Module](#-template-module)
9.  [Markdown Module](#-markdown-module)

### Format Converters
10. [CSV Converter](#-csv-converter)
11. [Excel Converter](#-excel-converter)
12. [DOCX Converter](#-docx-converter)
13. [TXT Converter](#-txt-converter)

### RAG System
14. [RAG Processor](#-rag-processor)
15. [Chunker](#-chunker)
16. [LLM Processor](#-llm-processor)
17. [Image Describer](#-image-describer)

### Batch Processing
18. [Batch Processor](#-batch-processor)

### Storage Backends
19. [Local Storage](#-local-storage)
20. [AWS S3 Storage](#-aws-s3-storage)
21. [Azure Blob Storage](#-azure-blob-storage)
22. [Vector Stores](#-vector-stores)
    - [PostgreSQL (pgvector)](#postgresql-pgvector)
    - [MongoDB Atlas](#mongodb-atlas)

### Search & Retrieval
23. [Azure AI Search](#-azure-ai-search)

### Configuration
24. [Chunking Configuration](#-chunking-configuration)
25. [LLM Configuration](#-llm-configuration)
26. [Document Intelligence Configuration](#-document-intelligence-configuration)
27. [Metadata Configuration](#-metadata-configuration)

### Advanced Topics
28. [Type System](#-type-system)
29. [Error Handling](#-error-handling)
30. [Performance Optimization](#-performance-optimization)
31. [Troubleshooting & FAQ](#-troubleshooting--faq)
32. [Contributing](#-contributing)
33. [License](#-license)

---

# 🌟 Introduction

## Why DocFlow?

Building AI applications that work with real-world documents is hard:

| Problem | DocFlow Solution |
|---------|------------------|
| PDFs export as garbage text | Smart layout analysis preserves structure |
| Tables become unreadable | Structured table extraction to Markdown/HTML/DataFrame |
| Images are ignored | Vision LLM integration for image descriptions |
| Headers disconnect from content | Heading-aware chunking preserves context |
| Single file is easy, 100K files is hard | Built-in batch processor with queues |
| Different vector DBs need different code | Unified interface for Postgres, Mongo, etc. |

## Architecture Overview

DocFlow follows a **Pipeline Architecture**:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           DocFlow Pipeline                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────┐    ┌───────────┐    ┌──────────┐    ┌─────────────────┐  │
│  │  INPUT   │───▶│ CONVERTER │───▶│ ENRICHER │───▶│    CHUNKER      │  │
│  │          │    │           │    │          │    │                 │  │
│  │ PDF      │    │ Format    │    │ LLM      │    │ Heading-Aware   │  │
│  │ DOCX     │    │ Detection │    │ Image    │    │ Semantic Split  │  │
│  │ Excel    │    │ to MD     │    │ Table    │    │                 │  │
│  └──────────┘    └───────────┘    └──────────┘    └────────┬────────┘  │
│                                                             │           │
│                                                             ▼           │
│                                                    ┌─────────────────┐  │
│                                                    │   VECTOR STORE  │  │
│                                                    │                 │  │
│                                                    │ Postgres/Mongo  │  │
│                                                    │ Embeddings      │  │
│                                                    └─────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

### Core Components

| Component | Module | Purpose |
|-----------|--------|---------|
| **Converter** | `docflow.converter` | MD → PDF conversion |
| **Extractor** | `docflow.extractor` | PDF → MD extraction |
| **Template** | `docflow.template` | Jinja2 template rendering |
| **Formats** | `docflow.formats` | CSV, Excel, DOCX, TXT parsing |
| **RAG** | `docflow.rag` | Chunking, LLM processing |
| **Storage** | `docflow.storage` | Local, S3, Azure, Vector DBs |
| **Search** | `docflow.search` | Azure AI Search integration |
| **Batch** | `docflow.batch_processor` | Parallel file processing |
| **Config** | `docflow.config` | All configuration classes |

---

# 📦 Installation

## Standard Installation
For basic conversion features:

```bash
pip install docflow-client
```

## RAG Installation (Recommended)
Includes LLM libraries, chunking, and vector store drivers:

```bash
pip install "docflow-client[rag]"
```

## Full Installation
Everything including OCR, all vector DBs, and heavy dependencies:

```bash
pip install "docflow-client[all]"
```

## Specific Extras

```bash
# Only PostgreSQL vector store
pip install "docflow-client[postgres]"

# Only MongoDB vector store
pip install "docflow-client[mongo]"

# Only Azure integrations
pip install "docflow-client[azure]"

# Only AWS integrations
pip install "docflow-client[aws]"
```

## Docker Installation

```dockerfile
FROM python:3.11-slim

WORKDIR /app
RUN pip install "docflow-client[all]"
COPY . .

CMD ["python", "app.py"]
```

## Environment Variables

```bash
# LLM Providers
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export AZURE_OPENAI_API_KEY="..."
export AZURE_OPENAI_ENDPOINT="https://..."

# Document Intelligence (OCR)
export AZURE_FORM_RECOGNIZER_ENDPOINT="https://..."
export AZURE_FORM_RECOGNIZER_KEY="..."

# Storage
export AWS_ACCESS_KEY_ID="..."
export AWS_SECRET_ACCESS_KEY="..."
export AZURE_STORAGE_CONNECTION_STRING="..."

# Vector Stores
export POSTGRES_CONNECTION_STRING="postgresql://..."
export MONGODB_URI="mongodb+srv://..."

# Search
export AZURE_SEARCH_ENDPOINT="https://..."
export AZURE_SEARCH_API_KEY="..."
```

---

# 🚀 Quick Start

## Example 1: Simple PDF to Markdown

```python
from docflow import Converter, LocalStorage

# Initialize with local storage for output
storage = LocalStorage("./output")
converter = Converter(storage=storage)

# Convert any file to markdown
result = converter.convert_to_markdown("document.pdf")

if result.success:
    print(result.content)
else:
    print(f"Error: {result.error}")
```

## Example 2: RAG Pipeline in 10 Lines

```python
from docflow.rag import RAGProcessor, RAGConfig

config = RAGConfig(
    chunk_size=500,
    chunk_overlap=50,
    extract_images=True,
    chunking_strategy="heading_aware"
)

processor = RAGProcessor(config)
doc = processor.process_file("handbook.docx")

print(f"Created {len(doc.chunks)} chunks")
for chunk in doc.chunks:
    print(f"[{chunk.metadata.heading_path}] {chunk.content[:100]}...")
```

## Example 3: Batch Processing 1000 Files

```python
from docflow import BatchProcessor, BatchConfig, RAGConfig

rag_cfg = RAGConfig()
batch_cfg = BatchConfig(max_workers=16, retry_failed=True)

processor = BatchProcessor(rag_cfg, batch_cfg)
job_id = processor.enqueue(["file1.pdf", "file2.docx", ...])

# Monitor
while True:
    status = processor.get_status(job_id)
    print(f"{status.processed_files}/{status.total_files}")
    if status.status == "completed":
        break
```

---

# 📄 Converter Module

**Location**: `docflow.converter`

The Converter module handles **Markdown → PDF** conversion using the DocFlow API.

## Basic Usage

```python
from docflow import Converter, LocalStorage, MDFile

# Initialize
storage = LocalStorage("./output")
converter = Converter(storage=storage)

# Single file
md_file = MDFile(name="doc.md", content="# Hello World\n\nThis is content.")
result = converter.convert_to_pdf([md_file])

print(f"PDF saved to: {result.download_url}")
```

## Multiple Files

```python
files = [
    MDFile(name="chapter1.md", content="# Chapter 1\n\nContent..."),
    MDFile(name="chapter2.md", content="# Chapter 2\n\nMore content..."),
]

# Merge into single PDF
result = converter.convert_to_pdf(files, merge=True)

# Or separate PDFs
result = converter.convert_to_pdf(files, merge=False)
```

## With Options

```python
from docflow import ConvertOptions

options = ConvertOptions(
    page_format="A4",
    margin_top=20,
    margin_bottom=20,
    header_template="<div>Header</div>",
    footer_template="<div>Page {page}</div>"
)

result = converter.convert_to_pdf(files, options=options)
```

## API Reference: Converter

```python
class Converter:
    def __init__(
        self,
        storage: Storage = None,
        api_key: str = None,
        base_url: str = None
    )
    
    def convert_to_pdf(
        self,
        files: List[MDFile],
        options: ConvertOptions = None,
        merge: bool = True
    ) -> PDFResult
    
    def convert_to_markdown(
        self,
        file_path: str
    ) -> ConvertResult
```

---

# 📤 Extractor Module

**Location**: `docflow.extractor`

The Extractor module handles **PDF → Markdown** extraction.

## Basic Usage

```python
from docflow import Extractor

extractor = Extractor()

# From file path
result = extractor.extract("document.pdf")
print(result.content)

# From bytes
with open("document.pdf", "rb") as f:
    result = extractor.extract_bytes(f.read(), filename="document.pdf")
```

## With Options

```python
from docflow import ExtractOptions

options = ExtractOptions(
    extract_images=True,
    extract_tables=True,
    ocr_enabled=True,
    language="en"
)

result = extractor.extract("scanned.pdf", options=options)
```

## Accessing Extracted Data

```python
result = extractor.extract("report.pdf")

# Full markdown content
print(result.content)

# Metadata
print(result.metadata)  # {"title": "...", "author": "...", "pages": 15}

# Extracted images (if enabled)
for img in result.images:
    print(f"Image: {img.filename}, Description: {img.description}")

# Extracted tables (if enabled)
for table in result.tables:
    print(table.markdown)
    df = table.to_pandas()  # Convert to DataFrame
```

---

# 🎨 Template Module

**Location**: `docflow.template`

The Template module provides **Jinja2-based** templating for dynamic document generation.

## Basic Usage

```python
from docflow import TemplateEngine

engine = TemplateEngine()

# Simple template
template = "# Hello {{ name }}\n\nWelcome to {{ company }}!"
result = engine.render(template, name="Alice", company="Acme Inc")
print(result)
# Output:
# # Hello Alice
#
# Welcome to Acme Inc!
```

## Loading Templates from Files

```python
engine = TemplateEngine(template_dir="./templates")

# Load and render template
result = engine.render_file("invoice.md.j2", 
    invoice_number="INV-001",
    items=[
        {"name": "Widget", "price": 10.00},
        {"name": "Gadget", "price": 25.00}
    ]
)
```

## Template Example: Invoice

**templates/invoice.md.j2**
```jinja2
# Invoice {{ invoice_number }}

| Item | Price |
|------|-------|
{% for item in items %}
| {{ item.name }} | ${{ "%.2f"|format(item.price) }} |
{% endfor %}

**Total**: ${{ items|sum(attribute='price')|round(2) }}
```

## Built-in Filters

```python
# Date formatting
"{{ date|dateformat('%Y-%m-%d') }}"

# Currency
"{{ amount|currency('USD') }}"

# Markdown escaping
"{{ user_input|md_escape }}"
```

---

# 📝 Markdown Module

**Location**: `docflow.markdown`

Utilities for parsing and manipulating Markdown content.

## Parsing Headers

```python
from docflow.markdown import parse_headers

content = """
# Title
## Section 1
Some text
## Section 2
More text
### Subsection
"""

headers = parse_headers(content)
# [
#   {"level": 1, "text": "Title", "line": 1},
#   {"level": 2, "text": "Section 1", "line": 2},
#   {"level": 2, "text": "Section 2", "line": 4},
#   {"level": 3, "text": "Subsection", "line": 6}
# ]
```

## Extracting Sections

```python
from docflow.markdown import extract_section

content = """
# Document
## Introduction
This is the intro.
## Methods
This is methodology.
"""

intro = extract_section(content, "Introduction")
# "This is the intro."
```

## Table Parsing

```python
from docflow.markdown import parse_tables

content = """
# Data

| Name | Value |
|------|-------|
| A    | 1     |
| B    | 2     |
"""

tables = parse_tables(content)
# [{"headers": ["Name", "Value"], "rows": [["A", "1"], ["B", "2"]]}]
```

---

# 📊 CSV Converter

**Location**: `docflow.formats.csv_format`

Converts CSV files to/from Markdown tables.

## CSV to Markdown

```python
from docflow.formats import CSVConverter

converter = CSVConverter()

# From file
result = converter.to_markdown("data.csv")
print(result.content)

# From string
csv_data = "name,age,city\nAlice,30,NYC\nBob,25,LA"
result = converter.to_markdown_string(csv_data, filename="users.csv")
```

## Markdown to CSV

```python
markdown = """
| Name | Age |
|------|-----|
| Alice | 30 |
| Bob | 25 |
"""

csv_content = converter.from_markdown(markdown)
# "Name,Age\nAlice,30\nBob,25"
```

## Options

```python
converter = CSVConverter(
    delimiter=";",           # Semicolon separated
    encoding="utf-8",
    skip_empty_rows=True,
    max_column_width=50      # Truncate long values
)
```

## Auto-Detect Delimiter

```python
result = converter.to_markdown("european_data.csv")
# Automatically detects ; , | \t as delimiters
```

---

# 📗 Excel Converter

**Location**: `docflow.formats.excel_format`

Converts Excel (XLSX/XLS) files to/from Markdown.

## Excel to Markdown

```python
from docflow.formats import ExcelConverter

converter = ExcelConverter()

# Convert all sheets
result = converter.to_markdown("financials.xlsx")
print(result.content)

# Output:
# # financials.xlsx
#
# ## Sheet1
#
# | Col1 | Col2 | Col3 |
# |------|------|------|
# | val1 | val2 | val3 |
#
# ## Sheet2
# ...
```

## Specific Sheets

```python
result = converter.to_markdown(
    "data.xlsx",
    sheets=["Revenue", "Expenses"]  # Only these sheets
)
```

## Options

```python
converter = ExcelConverter(
    include_all_sheets=True,
    sheet_separator="\n\n---\n\n",
    header_row=0,            # Row to use as headers
    skip_empty_sheets=True,
    max_rows=1000            # Limit rows per sheet
)
```

## Access as DataFrame

```python
result = converter.to_markdown("data.xlsx")

# Access raw data
for sheet_name, df in result.dataframes.items():
    print(f"Sheet: {sheet_name}")
    print(df.describe())
```

---

# 📘 DOCX Converter

**Location**: `docflow.formats.docx_format`

Converts Word documents (DOCX) to/from Markdown.

## DOCX to Markdown

```python
from docflow.formats import DOCXConverter

converter = DOCXConverter()

result = converter.to_markdown("document.docx")
print(result.content)
```

## Features Preserved

- **Headers** → `#`, `##`, `###`
- **Bold** → `**text**`
- **Italic** → `*text*`
- **Lists** → `-` or `1.`
- **Tables** → Markdown tables
- **Images** → Extracted and referenced

## Extracting Images

```python
converter = DOCXConverter(extract_images=True)
result = converter.to_markdown("presentation.docx")

for img in result.images:
    print(f"Found: {img.filename}")
    # img.data contains the raw bytes
```

## Options

```python
converter = DOCXConverter(
    extract_images=True,
    preserve_formatting=True,
    combine_runs=True       # Merge adjacent text runs
)
```

---

# 📄 TXT Converter

**Location**: `docflow.formats.txt_format`

Converts plain text files to structured Markdown.

## TXT to Markdown

```python
from docflow.formats import TXTConverter

converter = TXTConverter()
result = converter.to_markdown("notes.txt")
```

## Smart Detection

The TXT converter attempts to detect structure:

```python
text = """
COMPANY REPORT 2024

EXECUTIVE SUMMARY

This is the summary section with important information.

KEY FINDINGS

1. Revenue increased by 15%
2. Customer base grew by 20%
3. New markets entered

RECOMMENDATIONS

* Focus on digital
* Expand team
"""

result = converter.to_markdown_string(text, "report.txt")
# Detects all-caps lines as headers
# Detects numbered lists
# Detects bullet points
```

## Options

```python
converter = TXTConverter(
    detect_headers=True,     # Auto-detect headers
    detect_lists=True,       # Auto-detect lists
    paragraph_separator="\n\n",
    min_header_length=3,
    max_header_length=80
)
```

---

# 🤖 RAG Processor

**Location**: `docflow.rag.rag_processor`

The central orchestrator for the RAG pipeline.

## Basic Usage

```python
from docflow.rag import RAGProcessor, RAGConfig

config = RAGConfig(
    chunk_size=1000,
    chunk_overlap=200,
    extract_images=True,
    extract_tables=True
)

processor = RAGProcessor(config)
doc = processor.process_file("document.pdf")
```

## Accessing Results

```python
doc = processor.process_file("report.pdf")

# Document ID (UUID)
print(doc.id)

# Original filename
print(doc.filename)

# Full content
print(doc.content)

# Chunks
for chunk in doc.chunks:
    print(f"Index: {chunk.index}")
    print(f"Content: {chunk.content}")
    print(f"Metadata: {chunk.metadata}")

# Images
for img in doc.images:
    print(f"Image: {img.filename}")
    print(f"Description: {img.description}")  # AI-generated

# Tables
for table in doc.tables:
    print(table.markdown)
    print(table.summary)  # AI-generated
```

## With LLM Enrichment

```python
from docflow.config import LLMConfig

llm_config = LLMConfig(
    provider="openai",
    model="gpt-4o",
    api_key="sk-..."
)

config = RAGConfig(
    describe_images=True,    # Enable image description
    summarize_tables=True,   # Enable table summaries
    llm_config=llm_config
)

processor = RAGProcessor(config)
```

## Processing Bytes

```python
with open("document.pdf", "rb") as f:
    doc = processor.process(f.read(), filename="document.pdf")
```

---

# ✂️ Chunker

**Location**: `docflow.rag.chunker`

Intelligently splits text into semantic chunks.

## Basic Usage

```python
from docflow.rag import Chunker, RAGConfig

config = RAGConfig(
    chunk_size=500,
    chunk_overlap=50
)

chunker = Chunker(config)
chunks = chunker.chunk(markdown_text)

for chunk in chunks:
    print(f"[{chunk.index}] {chunk.content[:100]}...")
```

## Chunking Strategies

### 1. Simple (Fixed Size)

```python
config = RAGConfig(
    chunking_strategy="simple",
    chunk_size=500
)
# Splits every 500 tokens regardless of structure
```

### 2. Heading-Aware (Recommended)

```python
config = RAGConfig(
    chunking_strategy="heading_aware",
    chunk_size=500
)
# Respects H1, H2, H3 boundaries
# Each chunk knows its heading context
```

### 3. Semantic

```python
config = RAGConfig(
    chunking_strategy="semantic",
    chunk_size=500
)
# Uses sentence boundaries and paragraph structure
```

## Chunk Metadata

```python
for chunk in chunks:
    meta = chunk.metadata
    
    print(meta.index)           # Chunk number
    print(meta.heading_path)    # ["Section 1", "Subsection A"]
    print(meta.has_code)        # True if contains code blocks
    print(meta.has_table)       # True if contains tables
    print(meta.content_type)    # "text", "code", "table"
    print(meta.start_pos)       # Character position
    print(meta.end_pos)
```

## Token Counting

```python
from docflow.rag import Chunker

# Uses tiktoken for accurate token counting
chunker = Chunker(config)

token_count = chunker.count_tokens("Some text here")
print(f"Tokens: {token_count}")
```

---

# 🧠 LLM Processor

**Location**: `docflow.rag.llm_processor`

Unified interface for LLM operations: image description, table summarization, metadata extraction.

## Basic Usage

```python
from docflow.rag import LLMProcessor
from docflow.config import LLMConfig

config = LLMConfig(
    provider="openai",
    model="gpt-4o",
    api_key="sk-..."
)

processor = LLMProcessor(config)
```

## Describe Image

```python
with open("chart.png", "rb") as f:
    image_bytes = f.read()

description = processor.describe_image(
    image_bytes, 
    context="This is a financial chart from a Q3 report"
)
print(description)
# "A bar chart showing quarterly revenue..."
```

## Summarize Table

```python
table_markdown = """
| Product | Q1 | Q2 | Q3 |
|---------|----|----|------|
| Widget  | 100 | 150 | 200 |
| Gadget  | 80 | 90 | 120 |
"""

summary = processor.summarize_table(table_markdown)
print(summary)
# "The table shows Widget sales growing from 100 to 200 units..."
```

## Extract Metadata

```python
content = "This document discusses machine learning strategies..."

metadata = processor.extract_metadata(content)
# {"title": "ML Strategies", "topics": ["machine learning"], ...}
```

## Supported Providers

### OpenAI

```python
LLMConfig(
    provider="openai",
    model="gpt-4o",  # or gpt-4-vision-preview, gpt-4-turbo
    api_key="sk-..."
)
```

### Anthropic

```python
LLMConfig(
    provider="anthropic",
    model="claude-3-opus-20240229",
    api_key="sk-ant-..."
)
```

### Azure OpenAI

```python
LLMConfig(
    provider="azure",
    model="gpt-4",
    api_key="...",
    base_url="https://your-resource.openai.azure.com/",
    api_version="2024-02-15-preview"
)
```

### Ollama (Local)

```python
LLMConfig(
    provider="ollama",
    model="llava",  # or llama3, mistral
    base_url="http://localhost:11434"
)
```

---

# 🖼️ Image Describer

**Location**: `docflow.rag.image_describer`

Specialized component for describing images using Vision LLMs.

## Basic Usage

```python
from docflow.rag import ImageDescriber
from docflow.config import LLMConfig

config = LLMConfig(provider="openai", model="gpt-4o")
describer = ImageDescriber(config)

description = describer.describe(image_bytes)
```

## For RAG Context

```python
# Optimized for RAG systems
description = describer.describe_for_rag(
    image_bytes,
    document_context="Financial report Q3 2024",
    surrounding_text="As shown in the chart below..."
)
# Returns detailed description optimized for search
```

## Batch Processing

```python
images = [img1_bytes, img2_bytes, img3_bytes]

descriptions = describer.describe_batch(images)
for img, desc in zip(images, descriptions):
    print(desc)
```

---

# 🏭 Batch Processor

**Location**: `docflow.batch_processor`

Process thousands of files with parallel workers, queuing, and fault tolerance.

## Basic Usage

```python
from docflow import BatchProcessor, BatchConfig, RAGConfig

rag_config = RAGConfig(chunk_size=1000)

batch_config = BatchConfig(
    max_workers=16,
    queue_size=5000,
    fail_fast=False,
    retry_failed=True,
    max_retries=3,
    timeout_per_file=300
)

processor = BatchProcessor(rag_config, batch_config)
```

## Processing Files

```python
# Queue files
files = ["doc1.pdf", "doc2.docx", "data.xlsx", ...]
job_id = processor.enqueue(files)

# Monitor progress
import time
while True:
    status = processor.get_status(job_id)
    
    print(f"Progress: {status.processed_files}/{status.total_files}")
    print(f"Failed: {status.failed_files}")
    print(f"Status: {status.status}")
    
    if status.status in ["completed", "failed"]:
        break
    
    time.sleep(2)
```

## Getting Results

```python
results = processor.get_result(job_id)

for doc in results:
    if doc.error:
        print(f"Error in {doc.filename}: {doc.error}")
    else:
        print(f"Success: {doc.filename} - {len(doc.chunks)} chunks")
```

## Configuration Options

```python
BatchConfig(
    max_workers=16,          # Concurrent workers
    queue_size=5000,         # Max pending jobs
    fail_fast=False,         # Stop on first error?
    retry_failed=True,       # Auto-retry failures
    max_retries=3,           # Retry count
    timeout_per_file=300,    # 5 minutes max
    callback=my_callback     # Called for each result
)
```

## Callback Function

```python
def on_result(doc):
    """Called when each document is processed"""
    print(f"Finished: {doc.filename}")
    # Save to database immediately
    db.save(doc)

processor = BatchProcessor(
    rag_config, 
    BatchConfig(callback=on_result)
)
```

---

# 💾 Local Storage

**Location**: `docflow.storage.local`

Stores files on the local filesystem.

## Basic Usage

```python
from docflow.storage import LocalStorage

storage = LocalStorage("./output")

# Save a file
path = storage.save("result.pdf", pdf_bytes)
print(path)  # ./output/result.pdf

# Read a file
data = storage.read("result.pdf")

# Check existence
exists = storage.exists("result.pdf")

# Delete
storage.delete("result.pdf")

# List files
files = storage.list()
```

## With Converter

```python
from docflow import Converter

storage = LocalStorage("./pdfs")
converter = Converter(storage=storage)

result = converter.convert_to_pdf(files)
# PDF automatically saved to ./pdfs/
```

---

# ☁️ AWS S3 Storage

**Location**: `docflow.storage.s3`

Stores files in Amazon S3.

## Basic Usage

```python
from docflow.storage import S3Storage

storage = S3Storage(
    bucket_name="my-bucket",
    region_name="us-east-1",
    access_key="AKIA...",      # Optional if using AWS credentials
    secret_key="..."           # Optional if using AWS credentials
)

# Save
url = storage.save("documents/report.pdf", pdf_bytes)
print(url)  # s3://my-bucket/documents/report.pdf

# Read
data = storage.read("documents/report.pdf")

# Generate presigned URL
url = storage.get_presigned_url("documents/report.pdf", expires_in=3600)
```

## With Prefix

```python
storage = S3Storage(
    bucket_name="my-bucket",
    prefix="rag-outputs/"  # All files saved under this prefix
)
```

---

# 🔷 Azure Blob Storage

**Location**: `docflow.storage.azure`

Stores files in Azure Blob Storage.

## Basic Usage

```python
from docflow.storage import AzureStorage

storage = AzureStorage(
    account_name="mystorageaccount",
    container_name="documents",
    account_key="..."  # Or use connection_string
)

# Or via connection string
storage = AzureStorage(
    connection_string="DefaultEndpointsProtocol=https;AccountName=..."
)

# Save
url = storage.save("report.pdf", pdf_bytes)

# Read
data = storage.read("report.pdf")
```

---

# 🗄️ Vector Stores

**Location**: `docflow.storage.vector`

Native integrations with vector databases for semantic search.

## PostgreSQL (pgvector)

```python
from docflow.storage.vector import PostgresVectorStore

store = PostgresVectorStore(
    connection_string="postgresql://user:pass@localhost:5432/db",
    table_name="embeddings",
    dimension=1536,  # OpenAI embedding dimension
    distance_metric="cosine"  # or "l2", "inner_product"
)

# Upsert a processed document
store.upsert(doc)

# Search
results = store.search(
    query="revenue growth",
    top_k=5,
    filter={"source": "quarterly_report.pdf"}
)

for result in results:
    print(f"Score: {result.score}")
    print(f"Content: {result.content}")
    print(f"Metadata: {result.metadata}")
```

### Table Schema

```sql
CREATE TABLE embeddings (
    id UUID PRIMARY KEY,
    content TEXT,
    embedding vector(1536),
    metadata JSONB,
    created_at TIMESTAMP
);

CREATE INDEX ON embeddings USING ivfflat (embedding vector_cosine_ops);
```

## MongoDB Atlas

```python
from docflow.storage.vector import MongoDBVectorStore

store = MongoDBVectorStore(
    uri="mongodb+srv://user:pass@cluster.mongodb.net",
    database="rag_db",
    collection="embeddings",
    index_name="vector_index",
    dimension=1536
)

# Upsert
store.upsert(doc)

# Search
results = store.search("machine learning risks", top_k=5)
```

### Index Configuration

```javascript
{
  "mappings": {
    "dynamic": true,
    "fields": {
      "embedding": {
        "type": "knnVector",
        "dimensions": 1536,
        "similarity": "cosine"
      }
    }
  }
}
```

---

# 🔍 Azure AI Search

**Location**: `docflow.search.azure_search`

Enterprise-grade hybrid search with Azure AI Search.

## Basic Usage

```python
from docflow.search import AzureAISearch

client = AzureAISearch(
    endpoint="https://my-search.search.windows.net",
    api_key="...",
    index_name="documents"
)
```

## Vector Search

```python
# Get embedding from OpenAI or your model
query_embedding = get_embedding("machine learning risks")

results = client.vector_search(
    vector=query_embedding,
    top=10,
    select=["content", "title", "source"]
)
```

## Keyword Search

```python
results = client.keyword_search(
    query="quarterly revenue",
    top=10,
    filter="category eq 'finance'"
)
```

## Hybrid Search (Recommended)

```python
results = client.hybrid_search(
    query="What are the main risks?",
    vector=query_embedding,
    top=10,
    hybrid_mode="rrf"  # Reciprocal Rank Fusion
)
```

## Semantic Search (Re-ranking)

```python
results = client.semantic_search(
    query="risk factors",
    semantic_configuration="my-semantic-config",
    top=10,
    query_caption=True,
    query_answer=True
)

for result in results:
    print(result.content)
    print(result.caption)    # Highlighted relevant passage
    print(result.answer)     # Direct answer if found
```

## Index Management

```python
# Create index
client.create_index(
    fields=[
        {"name": "id", "type": "Edm.String", "key": True},
        {"name": "content", "type": "Edm.String", "searchable": True},
        {"name": "embedding", "type": "Collection(Edm.Single)", 
         "dimensions": 1536, "vectorSearchProfile": "my-profile"}
    ]
)

# Upload documents
client.upload_documents([
    {"id": "1", "content": "...", "embedding": [...]}
])
```

---

# ⚙️ Chunking Configuration

**Location**: `docflow.config.chunking`

```python
from docflow.config import ChunkingConfig

config = ChunkingConfig(
    chunk_size=1000,             # Target tokens per chunk
    chunk_overlap=200,           # Token overlap
    strategy="heading_aware",    # "simple", "heading_aware", "semantic"
    min_chunk_size=100,          # Minimum chunk size
    max_chunk_size=2000,         # Maximum chunk size
    respect_code_blocks=True,    # Don't split code
    respect_tables=True,         # Don't split tables
    tokenizer="cl100k_base"      # OpenAI tokenizer
)
```

---

# 🤖 LLM Configuration

**Location**: `docflow.config.llm`

```python
from docflow.config import LLMConfig

config = LLMConfig(
    provider="openai",           # openai, anthropic, azure, ollama
    model="gpt-4o",
    api_key="sk-...",
    base_url=None,               # Custom endpoint
    temperature=0.0,             # 0.0 for deterministic
    max_tokens=1000,
    timeout=60,
    retry_count=3,
    retry_delay=1.0
)
```

---

# 📊 Document Intelligence Configuration

**Location**: `docflow.config.doc_intel`

For OCR and advanced document parsing.

```python
from docflow.config import DocIntelConfig

config = DocIntelConfig(
    provider="azure",            # azure, aws
    endpoint="https://...",
    api_key="...",
    model_id="prebuilt-layout",  # prebuilt-read, prebuilt-document
    locale="en-US",
    features=["tables", "figures", "formulas"]
)
```

---

# 🏷️ Metadata Configuration

**Location**: `docflow.config.metadata`

```python
from docflow.config import MetadataConfig

config = MetadataConfig(
    include_fields=["title", "author", "date"],
    exclude_fields=["internal_id"],
    custom_extractors={
        "department": lambda doc: extract_department(doc)
    }
)
```

---

# 📐 Type System

**Location**: `docflow.types`

All core types used throughout the SDK.

## Enums

```python
from docflow import (
    LLMProcessingMode,    # IMAGE_ONLY, TABLE_ONLY, FULL
    ChunkingStrategy,     # SIMPLE, HEADING_AWARE, SEMANTIC
    OutputFormat,         # MARKDOWN, HTML, JSON
    JobStatus             # PENDING, PROCESSING, COMPLETED, FAILED
)
```

## Dataclasses

```python
from docflow import (
    RAGDocument,          # Main document container
    Chunk,                # Text segment
    ChunkMetadata,        # Chunk context
    ExtractedImage,       # Image with description
    ExtractedTable,       # Table with markdown/summary
    DocumentMetadata,     # File metadata
    ConvertResult,        # Conversion output
    BatchJob              # Job status
)
```

## Type Hints

```python
from docflow.types import ChunkList, ImageList, TableList

def process(chunks: ChunkList) -> None:
    for chunk in chunks:
        ...
```

---

# ⚠️ Error Handling

## Custom Exceptions

```python
from docflow.exceptions import (
    DocFlowError,           # Base exception
    ConversionError,        # File conversion failed
    ExtractionError,        # Content extraction failed
    ChunkingError,          # Chunking failed
    LLMError,               # LLM API error
    StorageError,           # Storage operation failed
    ValidationError         # Invalid input
)
```

## Usage

```python
from docflow.exceptions import ConversionError, LLMError

try:
    doc = processor.process_file("corrupt.pdf")
except ConversionError as e:
    print(f"Could not convert: {e}")
except LLMError as e:
    print(f"LLM failed: {e}")
```

---

# ⚡ Performance Optimization

## 1. Disable Unnecessary Features

```python
# Fastest config (no LLM calls)
config = RAGConfig(
    extract_images=False,
    extract_tables=False,
    describe_images=False,
    summarize_tables=False
)
```

## 2. Batch Processing

```python
# Use BatchProcessor for many files
batch_config = BatchConfig(
    max_workers=os.cpu_count() * 2
)
```

## 3. Caching

```python
from functools import lru_cache

@lru_cache(maxsize=1000)
def get_embedding(text):
    return embedding_model.encode(text)
```

## 4. Streaming for Large Files

```python
# Process in chunks for large files
for page in extractor.stream("huge.pdf"):
    process(page)
```

---

# ❓ Troubleshooting & FAQ

**Q: `ImportError: cannot import name 'RAGProcessor'`**
A: Install RAG extras: `pip install "docflow-client[rag]"`

**Q: Images are not being described**
A: Ensure `describe_images=True` and provide `llm_config` with a vision model

**Q: PDF text is garbled**
A: Enable OCR with `DocIntelConfig`

**Q: Rate limits with OpenAI**
A: Increase `retry_count` and `retry_delay` in `LLMConfig`

**Q: Memory issues with large files**
A: Reduce `BatchConfig.max_workers` or process files individually

**Q: Tables look wrong**
A: Enable `summarize_tables=True` for LLM summarization

---

# 🤝 Contributing

1. Fork the repo
2. `pip install -r requirements-dev.txt`
3. `pytest tests/`
4. Submit PR

---

# 📜 License

MIT License

---

# 📖 Real-World Examples

## Example 1: Complete RAG Chatbot Backend

```python
"""
Complete example: Document ingestion service for a RAG chatbot.
"""
import os
from pathlib import Path
from docflow.rag import RAGProcessor, RAGConfig
from docflow.config import LLMConfig
from docflow.storage.vector import PostgresVectorStore

# Configuration
llm_config = LLMConfig(
    provider="openai",
    model="gpt-4o",
    api_key=os.environ["OPENAI_API_KEY"]
)

rag_config = RAGConfig(
    chunk_size=800,
    chunk_overlap=100,
    chunking_strategy="heading_aware",
    extract_images=True,
    describe_images=True,
    extract_tables=True,
    summarize_tables=True,
    llm_config=llm_config
)

# Vector Store
vector_store = PostgresVectorStore(
    connection_string=os.environ["DATABASE_URL"],
    table_name="knowledge_base",
    dimension=1536
)

# Processor
processor = RAGProcessor(rag_config)


def ingest_document(file_path: str) -> dict:
    """
    Ingest a single document into the knowledge base.
    
    Args:
        file_path: Path to the document
        
    Returns:
        dict with ingestion statistics
    """
    # Process
    doc = processor.process_file(file_path)
    
    # Store
    vector_store.upsert(doc)
    
    return {
        "document_id": doc.id,
        "filename": doc.filename,
        "chunks": len(doc.chunks),
        "images": len(doc.images),
        "tables": len(doc.tables)
    }


def ingest_directory(directory: str) -> list:
    """
    Ingest all documents in a directory.
    """
    results = []
    path = Path(directory)
    
    for file in path.glob("**/*"):
        if file.suffix.lower() in [".pdf", ".docx", ".xlsx", ".csv", ".txt"]:
            try:
                result = ingest_document(str(file))
                results.append({"status": "success", **result})
            except Exception as e:
                results.append({
                    "status": "error",
                    "filename": str(file),
                    "error": str(e)
                })
    
    return results


def search_knowledge_base(query: str, top_k: int = 5) -> list:
    """
    Search the knowledge base.
    """
    return vector_store.search(query, top_k=top_k)


# Usage
if __name__ == "__main__":
    # Ingest
    results = ingest_directory("./documents")
    print(f"Ingested {len(results)} documents")
    
    # Search
    hits = search_knowledge_base("What is our refund policy?")
    for hit in hits:
        print(f"Score: {hit.score:.3f}")
        print(f"Content: {hit.content[:200]}...")
        print("---")
```

## Example 2: Financial Report Analyzer

```python
"""
Extract and analyze financial data from annual reports.
"""
from docflow.rag import RAGProcessor, RAGConfig, LLMProcessor
from docflow.config import LLMConfig, DocIntelConfig
from docflow.formats import ExcelConverter
import json

# Use Azure Document Intelligence for accurate table extraction
doc_intel_config = DocIntelConfig(
    provider="azure",
    endpoint=os.environ["AZURE_FORM_ENDPOINT"],
    api_key=os.environ["AZURE_FORM_KEY"],
    model_id="prebuilt-layout"
)

llm_config = LLMConfig(
    provider="openai",
    model="gpt-4o",
    api_key=os.environ["OPENAI_API_KEY"]
)

rag_config = RAGConfig(
    extract_tables=True,
    summarize_tables=True,
    doc_intel_config=doc_intel_config,
    llm_config=llm_config
)

processor = RAGProcessor(rag_config)
llm = LLMProcessor(llm_config)


def analyze_financial_report(file_path: str) -> dict:
    """
    Analyze a financial report and extract key metrics.
    """
    # Process document
    doc = processor.process_file(file_path)
    
    # Extract tables
    tables_data = []
    for table in doc.tables:
        tables_data.append({
            "page": table.page,
            "markdown": table.markdown,
            "summary": table.summary
        })
    
    # Use LLM to extract key metrics
    prompt = f"""
    Analyze the following financial document and extract key metrics:
    
    Document Content:
    {doc.content[:5000]}
    
    Tables Found:
    {json.dumps(tables_data, indent=2)}
    
    Extract:
    1. Revenue figures
    2. Profit margins
    3. Year-over-year growth
    4. Key risks mentioned
    5. Future outlook
    
    Return as JSON.
    """
    
    analysis = llm.complete(prompt)
    
    return {
        "document": doc.filename,
        "pages": doc.metadata.get("pages", 0),
        "tables_found": len(tables_data),
        "analysis": json.loads(analysis)
    }


# Usage
result = analyze_financial_report("annual_report_2024.pdf")
print(json.dumps(result, indent=2))
```

## Example 3: Multi-Language Document Processor

```python
"""
Process documents in multiple languages with translation.
"""
from docflow.rag import RAGProcessor, RAGConfig
from docflow.config import LLMConfig

llm_config = LLMConfig(
    provider="openai",
    model="gpt-4o"
)

rag_config = RAGConfig(
    chunk_size=1000,
    llm_config=llm_config
)

processor = RAGProcessor(rag_config)


def process_with_translation(file_path: str, target_language: str = "en") -> dict:
    """
    Process a document and translate chunks to target language.
    """
    doc = processor.process_file(file_path)
    
    translated_chunks = []
    for chunk in doc.chunks:
        # Detect language and translate if needed
        translated = llm_config.translate(
            chunk.content, 
            target_language=target_language
        )
        translated_chunks.append({
            "original": chunk.content,
            "translated": translated,
            "metadata": chunk.metadata
        })
    
    return {
        "document_id": doc.id,
        "chunks": translated_chunks
    }
```

## Example 4: Legal Document Comparator

```python
"""
Compare two versions of a legal document.
"""
from docflow.rag import RAGProcessor, Chunker, RAGConfig
from difflib import unified_diff

rag_config = RAGConfig(chunking_strategy="heading_aware")
processor = RAGProcessor(rag_config)


def compare_documents(doc1_path: str, doc2_path: str) -> dict:
    """
    Compare two documents and highlight differences.
    """
    doc1 = processor.process_file(doc1_path)
    doc2 = processor.process_file(doc2_path)
    
    # Compare by sections
    sections1 = {c.metadata.section_title: c.content for c in doc1.chunks}
    sections2 = {c.metadata.section_title: c.content for c in doc2.chunks}
    
    all_sections = set(sections1.keys()) | set(sections2.keys())
    
    comparison = {}
    for section in all_sections:
        text1 = sections1.get(section, "")
        text2 = sections2.get(section, "")
        
        if text1 == text2:
            comparison[section] = {"status": "unchanged"}
        elif section not in sections1:
            comparison[section] = {"status": "added", "content": text2}
        elif section not in sections2:
            comparison[section] = {"status": "removed", "content": text1}
        else:
            diff = list(unified_diff(
                text1.splitlines(),
                text2.splitlines(),
                lineterm=""
            ))
            comparison[section] = {"status": "modified", "diff": diff}
    
    return comparison
```

## Example 5: Automated Report Generator

```python
"""
Generate reports from multiple source documents.
"""
from docflow import TemplateEngine, Converter
from docflow.rag import RAGProcessor, RAGConfig

template_engine = TemplateEngine(template_dir="./templates")
converter = Converter()
processor = RAGProcessor(RAGConfig())


def generate_summary_report(document_paths: list, template: str) -> bytes:
    """
    Generate a PDF summary report from multiple documents.
    """
    # Process all documents
    summaries = []
    for path in document_paths:
        doc = processor.process_file(path)
        summaries.append({
            "filename": doc.filename,
            "content": doc.content[:1000],
            "chunks": len(doc.chunks),
            "tables": len(doc.tables),
            "images": len(doc.images)
        })
    
    # Render template
    markdown = template_engine.render_file(template, documents=summaries)
    
    # Convert to PDF
    result = converter.convert_to_pdf([
        MDFile(name="report.md", content=markdown)
    ])
    
    return result.pdf_bytes
```

---

# 🏗️ Advanced Patterns

## Pattern 1: Pipeline Composition

```python
from docflow.rag import RAGProcessor, Chunker, LLMProcessor
from docflow.config import RAGConfig, LLMConfig

class CustomPipeline:
    """
    Build a custom processing pipeline with granular control.
    """
    
    def __init__(self, rag_config: RAGConfig, llm_config: LLMConfig):
        self.chunker = Chunker(rag_config)
        self.llm = LLMProcessor(llm_config)
        self.processors = []
    
    def add_processor(self, func):
        """Add a processing step."""
        self.processors.append(func)
        return self
    
    def process(self, content: str) -> list:
        """Run the pipeline."""
        # Chunk
        chunks = self.chunker.chunk(content)
        
        # Apply processors
        for processor in self.processors:
            chunks = [processor(chunk) for chunk in chunks]
        
        return chunks


# Usage
pipeline = CustomPipeline(rag_config, llm_config)
pipeline.add_processor(lambda c: sanitize(c))
pipeline.add_processor(lambda c: enrich(c))

chunks = pipeline.process(markdown_content)
```

## Pattern 2: Caching Layer

```python
import hashlib
from functools import lru_cache
import redis

class CachedProcessor:
    """
    Add caching to any processor.
    """
    
    def __init__(self, processor, redis_url: str):
        self.processor = processor
        self.redis = redis.from_url(redis_url)
    
    def _cache_key(self, content: str) -> str:
        return hashlib.sha256(content.encode()).hexdigest()
    
    def process(self, content: str):
        key = self._cache_key(content)
        
        # Check cache
        cached = self.redis.get(key)
        if cached:
            return json.loads(cached)
        
        # Process
        result = self.processor.process(content)
        
        # Cache
        self.redis.setex(key, 3600, json.dumps(result.to_dict()))
        
        return result
```

## Pattern 3: Async Processing

```python
import asyncio
from concurrent.futures import ThreadPoolExecutor
from docflow.rag import RAGProcessor

class AsyncRAGProcessor:
    """
    Async wrapper for RAGProcessor.
    """
    
    def __init__(self, config):
        self.processor = RAGProcessor(config)
        self.executor = ThreadPoolExecutor(max_workers=10)
    
    async def process_file(self, path: str):
        loop = asyncio.get_event_loop()
        return await loop.run_in_executor(
            self.executor,
            self.processor.process_file,
            path
        )
    
    async def process_many(self, paths: list):
        tasks = [self.process_file(p) for p in paths]
        return await asyncio.gather(*tasks)


# Usage
async def main():
    processor = AsyncRAGProcessor(rag_config)
    docs = await processor.process_many(["doc1.pdf", "doc2.pdf"])
```

## Pattern 4: Plugin Architecture

```python
from abc import ABC, abstractmethod
from typing import List

class ProcessorPlugin(ABC):
    """Base class for processor plugins."""
    
    @abstractmethod
    def pre_process(self, content: str) -> str:
        """Called before processing."""
        pass
    
    @abstractmethod
    def post_process(self, chunks: List) -> List:
        """Called after chunking."""
        pass


class PIIRedactionPlugin(ProcessorPlugin):
    """Remove PII from content."""
    
    def pre_process(self, content: str) -> str:
        # Redact emails
        content = re.sub(r'\b[\w.-]+@[\w.-]+\.\w+\b', '[EMAIL]', content)
        # Redact phone numbers
        content = re.sub(r'\b\d{3}[-.]?\d{3}[-.]?\d{4}\b', '[PHONE]', content)
        return content
    
    def post_process(self, chunks: List) -> List:
        return chunks


class PluggableProcessor:
    """Processor with plugin support."""
    
    def __init__(self, base_processor, plugins: List[ProcessorPlugin] = None):
        self.processor = base_processor
        self.plugins = plugins or []
    
    def add_plugin(self, plugin: ProcessorPlugin):
        self.plugins.append(plugin)
    
    def process(self, content: str):
        # Pre-processing
        for plugin in self.plugins:
            content = plugin.pre_process(content)
        
        # Core processing
        result = self.processor.process(content)
        
        # Post-processing
        for plugin in self.plugins:
            result.chunks = plugin.post_process(result.chunks)
        
        return result
```

---

# 📊 Best Practices

## 1. Chunk Size Guidelines

| Use Case | Recommended Size | Overlap |
|----------|------------------|---------|
| Q&A Chatbot | 500-800 | 50-100 |
| Document Search | 800-1200 | 100-200 |
| Summarization | 1500-2000 | 200-300 |
| Code Documentation | 300-500 | 50 |

## 2. LLM Provider Selection

| Provider | Best For | Cost |
|----------|----------|------|
| OpenAI | General purpose, Vision | $$ |
| Anthropic | Long context, Safety | $$$ |
| Ollama | Privacy, Local | Free |
| Azure | Enterprise, Compliance | $$ |

## 3. Storage Strategy

| Volume | Recommended |
|--------|-------------|
| < 10K chunks | PostgreSQL |
| 10K - 1M chunks | PostgreSQL + Indexing |
| > 1M chunks | Dedicated Vector DB (Pinecone, Milvus) |

## 4. Error Handling Best Practices

```python
from docflow.exceptions import DocFlowError
import logging

logger = logging.getLogger(__name__)

def robust_process(file_path: str):
    try:
        return processor.process_file(file_path)
    except DocFlowError as e:
        logger.error(f"Processing error: {e}")
        # Fallback strategy
        return fallback_process(file_path)
    except Exception as e:
        logger.exception(f"Unexpected error: {e}")
        raise
```

## 5. Memory Management

```python
import gc

def process_large_batch(files: list):
    """Process large batches with memory management."""
    results = []
    
    for i, file in enumerate(files):
        result = processor.process_file(file)
        results.append(summarize(result))  # Don't keep full objects
        
        # Periodic cleanup
        if i % 100 == 0:
            gc.collect()
    
    return results
```

---

# 🔧 Development & Testing

## Running Tests

```bash
# Install dev dependencies
pip install -r requirements-dev.txt

# Run all tests
pytest tests/

# Run with coverage
pytest tests/ --cov=docflow --cov-report=html

# Run specific test file
pytest tests/test_rag.py -v

# Run with debug output
pytest tests/ -v -s
```

## Creating Test Fixtures

```python
import pytest
from docflow.rag import RAGProcessor, RAGConfig

@pytest.fixture
def rag_processor():
    config = RAGConfig(chunk_size=500)
    return RAGProcessor(config)

@pytest.fixture
def sample_markdown():
    return """
    # Test Document
    
    ## Section 1
    Content for section 1.
    
    ## Section 2
    Content for section 2.
    """

def test_chunking(rag_processor, sample_markdown):
    chunks = rag_processor.chunker.chunk(sample_markdown)
    assert len(chunks) >= 2
```

## Mocking LLM Calls

```python
from unittest.mock import patch

def test_image_description():
    with patch('docflow.rag.LLMProcessor.describe_image') as mock:
        mock.return_value = "A chart showing growth"
        
        result = processor.process_file("doc_with_image.pdf")
        
        assert result.images[0].description == "A chart showing growth"
```

---

# 📈 Monitoring & Observability

## Logging

```python
import logging

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)

# DocFlow loggers
logging.getLogger('docflow').setLevel(logging.DEBUG)
logging.getLogger('docflow.rag').setLevel(logging.INFO)
```

## Metrics

```python
from prometheus_client import Counter, Histogram

# Define metrics
docs_processed = Counter('docflow_docs_processed_total', 'Total documents processed')
processing_time = Histogram('docflow_processing_seconds', 'Processing time')

@processing_time.time()
def process_with_metrics(file_path: str):
    result = processor.process_file(file_path)
    docs_processed.inc()
    return result
```

## OpenTelemetry Integration

```python
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider

trace.set_tracer_provider(TracerProvider())
tracer = trace.get_tracer(__name__)

def process_with_tracing(file_path: str):
    with tracer.start_as_current_span("process_document") as span:
        span.set_attribute("document.path", file_path)
        
        result = processor.process_file(file_path)
        
        span.set_attribute("document.chunks", len(result.chunks))
        return result
```

---

# 🔐 Security Considerations

## 1. API Key Management

```python
import os
from dotenv import load_dotenv

# Load from .env file
load_dotenv()

# Never hardcode keys
config = LLMConfig(
    api_key=os.environ.get("OPENAI_API_KEY")  # From environment
)
```

## 2. Content Sanitization

```python
def sanitize_content(content: str) -> str:
    """Remove sensitive patterns before processing."""
    # Remove credit card numbers
    content = re.sub(r'\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b', '[REDACTED]', content)
    
    # Remove SSNs
    content = re.sub(r'\b\d{3}-\d{2}-\d{4}\b', '[REDACTED]', content)
    
    return content
```

## 3. Rate Limiting

```python
from ratelimit import limits, sleep_and_retry

@sleep_and_retry
@limits(calls=10, period=60)  # 10 calls per minute
def rate_limited_process(file_path: str):
    return processor.process_file(file_path)
```

---

# 🌐 Deployment

## Docker Compose Example

```yaml
version: '3.8'

services:
  docflow-worker:
    build: .
    environment:
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - DATABASE_URL=postgresql://postgres:postgres@db:5432/docflow
    depends_on:
      - db
      - redis
    
  db:
    image: pgvector/pgvector:pg16
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: docflow
    volumes:
      - pgdata:/var/lib/postgresql/data
  
  redis:
    image: redis:7-alpine
    
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
  replicas: 3
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
```

---

# 📚 Glossary

| Term | Definition |
|------|------------|
| **RAG** | Retrieval-Augmented Generation - Combining retrieval with LLM generation |
| **Chunk** | A semantic segment of text, optimized for embedding |
| **Embedding** | Vector representation of text for similarity search |
| **Document Intelligence** | OCR and layout analysis services (Azure, AWS) |
| **Heading-Aware** | Chunking that respects document structure |
| **Hybrid Search** | Combining keyword (BM25) and vector (semantic) search |
| **pgvector** | PostgreSQL extension for vector similarity search |

---

# 🏢 Complete Azure Enterprise Pipeline

This is the **end-to-end production flow** that many enterprises use:

1. **Files → Markdown** (Conversion)
2. **Markdown → Azure Blob** (Storage)
3. **Document Intelligence** (OCR & Layout Analysis)
4. **Smart Chunking** (Heading-Aware)
5. **Azure AI Search** (Indexing & Retrieval)

## Full Implementation

```python
"""
Complete Azure Enterprise RAG Pipeline

Flow:
1. Convert any file to Markdown
2. Store original + markdown in Azure Blob
3. Use Document Intelligence for OCR (scanned docs)
4. Apply heading-aware chunking
5. Generate embeddings
6. Index in Azure AI Search
7. Provide hybrid search API
"""

import os
import uuid
from datetime import datetime
from typing import List, Optional

from docflow import Converter, Extractor
from docflow.storage import AzureStorage
from docflow.rag import RAGProcessor, Chunker
from docflow.config import RAGConfig, LLMConfig, DocIntelConfig
from docflow.search import AzureAISearch


class AzureEnterprisePipeline:
    """
    Production-ready Azure RAG pipeline.
    """
    
    def __init__(
        self,
        blob_connection_string: str,
        blob_container: str,
        doc_intel_endpoint: str,
        doc_intel_key: str,
        search_endpoint: str,
        search_key: str,
        search_index: str,
        openai_key: str
    ):
        # 1. Azure Blob Storage (for documents & markdown)
        self.blob_storage = AzureStorage(
            connection_string=blob_connection_string,
            container_name=blob_container
        )
        
        # 2. Document Intelligence (OCR for scanned docs)
        self.doc_intel_config = DocIntelConfig(
            provider="azure",
            endpoint=doc_intel_endpoint,
            api_key=doc_intel_key,
            model_id="prebuilt-layout"  # Best for documents with tables
        )
        
        # 3. LLM Config (for embeddings & enrichment)
        self.llm_config = LLMConfig(
            provider="openai",
            model="gpt-4o",
            api_key=openai_key
        )
        
        # 4. RAG Config
        self.rag_config = RAGConfig(
            chunk_size=800,
            chunk_overlap=100,
            chunking_strategy="heading_aware",
            extract_images=True,
            extract_tables=True,
            describe_images=True,
            summarize_tables=True,
            doc_intel_config=self.doc_intel_config,
            llm_config=self.llm_config
        )
        
        # 5. Azure AI Search
        self.search_client = AzureAISearch(
            endpoint=search_endpoint,
            api_key=search_key,
            index_name=search_index
        )
        
        # 6. Processors
        self.converter = Converter(storage=self.blob_storage)
        self.rag_processor = RAGProcessor(self.rag_config)
    
    
    def ingest_document(self, file_path: str) -> dict:
        """
        Complete ingestion pipeline for a single document.
        
        Steps:
        1. Read file
        2. Convert to Markdown
        3. Store original + markdown in Blob
        4. Process with Document Intelligence (if scanned)
        5. Chunk with heading-aware strategy
        6. Generate embeddings
        7. Index in Azure AI Search
        
        Returns:
            dict with ingestion statistics
        """
        document_id = str(uuid.uuid4())
        filename = os.path.basename(file_path)
        timestamp = datetime.utcnow().isoformat()
        
        print(f"[1/7] Reading file: {filename}")
        with open(file_path, "rb") as f:
            file_bytes = f.read()
        
        # Step 1: Store original file in Blob
        print(f"[2/7] Uploading original to Blob Storage")
        original_blob_path = f"originals/{document_id}/{filename}"
        original_url = self.blob_storage.save(original_blob_path, file_bytes)
        
        # Step 2: Convert to Markdown (uses Doc Intel for scanned PDFs)
        print(f"[3/7] Converting to Markdown (with OCR if needed)")
        doc = self.rag_processor.process(file_bytes, filename=filename)
        
        # Step 3: Store Markdown in Blob
        print(f"[4/7] Storing Markdown in Blob Storage")
        markdown_blob_path = f"markdown/{document_id}/{filename}.md"
        markdown_url = self.blob_storage.save(
            markdown_blob_path, 
            doc.content.encode("utf-8")
        )
        
        # Step 4: Prepare chunks for indexing
        print(f"[5/7] Preparing {len(doc.chunks)} chunks for indexing")
        search_documents = []
        
        for chunk in doc.chunks:
            # Generate embedding
            embedding = self._generate_embedding(chunk.content)
            
            search_doc = {
                "id": f"{document_id}_{chunk.index}",
                "document_id": document_id,
                "chunk_index": chunk.index,
                "content": chunk.content,
                "content_vector": embedding,
                
                # Metadata
                "filename": filename,
                "original_url": original_url,
                "markdown_url": markdown_url,
                "section_title": chunk.metadata.section_title or "",
                "heading_path": " > ".join(chunk.metadata.heading_path or []),
                "has_table": chunk.metadata.has_table,
                "has_code": chunk.metadata.has_code,
                
                # Timestamps
                "created_at": timestamp,
                "updated_at": timestamp
            }
            search_documents.append(search_doc)
        
        # Step 5: Index images (if any)
        for img in doc.images:
            embedding = self._generate_embedding(img.description)
            
            search_doc = {
                "id": f"{document_id}_img_{img.index}",
                "document_id": document_id,
                "chunk_index": -1,  # Special marker for images
                "content": img.description,
                "content_vector": embedding,
                "filename": filename,
                "content_type": "image",
                "image_url": img.path,
                "created_at": timestamp
            }
            search_documents.append(search_doc)
        
        # Step 6: Upload to Azure AI Search
        print(f"[6/7] Indexing {len(search_documents)} documents in Azure AI Search")
        self.search_client.upload_documents(search_documents)
        
        print(f"[7/7] Complete!")
        
        return {
            "document_id": document_id,
            "filename": filename,
            "original_url": original_url,
            "markdown_url": markdown_url,
            "chunks_indexed": len(doc.chunks),
            "images_indexed": len(doc.images),
            "tables_found": len(doc.tables)
        }
    
    
    def search(
        self, 
        query: str, 
        top: int = 5,
        filters: Optional[dict] = None
    ) -> List[dict]:
        """
        Hybrid search: Vector + Keyword + Semantic Re-ranking
        """
        # Generate query embedding
        query_vector = self._generate_embedding(query)
        
        # Build filter string
        filter_str = None
        if filters:
            filter_parts = [f"{k} eq '{v}'" for k, v in filters.items()]
            filter_str = " and ".join(filter_parts)
        
        # Execute hybrid search
        results = self.search_client.hybrid_search(
            query=query,
            vector=query_vector,
            vector_fields=["content_vector"],
            top=top,
            filter=filter_str,
            select=["id", "content", "filename", "heading_path", "section_title"]
        )
        
        return results
    
    
    def semantic_search(
        self,
        query: str,
        top: int = 5
    ) -> List[dict]:
        """
        Semantic search with answer extraction.
        """
        query_vector = self._generate_embedding(query)
        
        results = self.search_client.semantic_search(
            query=query,
            vector=query_vector,
            semantic_configuration="my-semantic-config",
            query_caption=True,
            query_answer=True,
            top=top
        )
        
        return results
    
    
    def _generate_embedding(self, text: str) -> List[float]:
        """Generate embedding using OpenAI."""
        import openai
        
        client = openai.OpenAI(api_key=self.llm_config.api_key)
        response = client.embeddings.create(
            model="text-embedding-3-small",
            input=text
        )
        return response.data[0].embedding


# ============================================================
# Usage Example
# ============================================================

if __name__ == "__main__":
    # Initialize Pipeline
    pipeline = AzureEnterprisePipeline(
        blob_connection_string=os.environ["AZURE_STORAGE_CONNECTION_STRING"],
        blob_container="documents",
        doc_intel_endpoint=os.environ["AZURE_FORM_RECOGNIZER_ENDPOINT"],
        doc_intel_key=os.environ["AZURE_FORM_RECOGNIZER_KEY"],
        search_endpoint=os.environ["AZURE_SEARCH_ENDPOINT"],
        search_key=os.environ["AZURE_SEARCH_KEY"],
        search_index="enterprise-docs",
        openai_key=os.environ["OPENAI_API_KEY"]
    )
    
    # Ingest a document
    result = pipeline.ingest_document("./contracts/vendor_agreement.pdf")
    print(f"Ingested: {result}")
    
    # Search
    hits = pipeline.search("What are the payment terms?", top=5)
    for hit in hits:
        print(f"Score: {hit['@search.score']:.3f}")
        print(f"Section: {hit['heading_path']}")
        print(f"Content: {hit['content'][:200]}...")
        print("---")
    
    # Semantic search with answer extraction
    answers = pipeline.semantic_search("What is the contract duration?")
    for answer in answers:
        if answer.get("@search.answers"):
            print(f"Answer: {answer['@search.answers'][0]['text']}")
```

## Azure Search Index Schema

Use this schema when creating your Azure AI Search index:

```json
{
  "name": "enterprise-docs",
  "fields": [
    {"name": "id", "type": "Edm.String", "key": true},
    {"name": "document_id", "type": "Edm.String", "filterable": true},
    {"name": "chunk_index", "type": "Edm.Int32", "sortable": true},
    {"name": "content", "type": "Edm.String", "searchable": true, "analyzer": "en.microsoft"},
    {"name": "content_vector", "type": "Collection(Edm.Single)", 
     "dimensions": 1536, "vectorSearchProfile": "vector-profile"},
    {"name": "filename", "type": "Edm.String", "filterable": true, "facetable": true},
    {"name": "original_url", "type": "Edm.String"},
    {"name": "markdown_url", "type": "Edm.String"},
    {"name": "section_title", "type": "Edm.String", "searchable": true},
    {"name": "heading_path", "type": "Edm.String", "searchable": true},
    {"name": "has_table", "type": "Edm.Boolean", "filterable": true},
    {"name": "has_code", "type": "Edm.Boolean", "filterable": true},
    {"name": "content_type", "type": "Edm.String", "filterable": true},
    {"name": "image_url", "type": "Edm.String"},
    {"name": "created_at", "type": "Edm.DateTimeOffset", "sortable": true}
  ],
  "vectorSearch": {
    "profiles": [
      {"name": "vector-profile", "algorithm": "hnsw-algorithm"}
    ],
    "algorithms": [
      {"name": "hnsw-algorithm", "kind": "hnsw", "hnswParameters": {"m": 4, "efConstruction": 400}}
    ]
  },
  "semantic": {
    "configurations": [
      {
        "name": "my-semantic-config",
        "prioritizedFields": {
          "contentFields": [{"fieldName": "content"}],
          "titleField": {"fieldName": "section_title"}
        }
      }
    ]
  }
}
```

## Environment Variables Required

```bash
# Azure Blob Storage
export AZURE_STORAGE_CONNECTION_STRING="DefaultEndpointsProtocol=https;AccountName=...;AccountKey=...;EndpointSuffix=core.windows.net"

# Azure Document Intelligence (Form Recognizer)
export AZURE_FORM_RECOGNIZER_ENDPOINT="https://your-resource.cognitiveservices.azure.com/"
export AZURE_FORM_RECOGNIZER_KEY="..."

# Azure AI Search
export AZURE_SEARCH_ENDPOINT="https://your-search.search.windows.net"
export AZURE_SEARCH_KEY="..."

# OpenAI (for embeddings)
export OPENAI_API_KEY="sk-..."
```

## Architecture Diagram

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────────┐
│   Input File    │────▶│    Converter     │────▶│   Azure Blob        │
│   (PDF/DOCX)    │     │   + Doc Intel    │     │   (Original + MD)   │
└─────────────────┘     └──────────────────┘     └─────────────────────┘
                                │
                                ▼
                    ┌──────────────────────┐
                    │   RAG Processor      │
                    │   - Extract Images   │
                    │   - Extract Tables   │
                    │   - LLM Enrichment   │
                    └──────────────────────┘
                                │
                                ▼
                    ┌──────────────────────┐
                    │   Heading-Aware      │
                    │   Chunker            │
                    │   - Context Metadata │
                    └──────────────────────┘
                                │
                                ▼
                    ┌──────────────────────┐
                    │   Embedding          │
                    │   (OpenAI/Azure)     │
                    └──────────────────────┘
                                │
                                ▼
                    ┌──────────────────────┐
                    │   Azure AI Search    │
                    │   - Vector Index     │
                    │   - Semantic Config  │
                    │   - Hybrid Search    │
                    └──────────────────────┘
                                │
                                ▼
                    ┌──────────────────────┐
                    │   RAG Application    │
                    │   (Chatbot, Q&A)     │
                    └──────────────────────┘
```

---

# 📞 Support

- **Documentation**: https://docflow.readthedocs.io
- **GitHub Issues**: https://github.com/xgaslan/docflow/issues
- **Discord**: https://discord.gg/docflow
- **Email**: support@docflow.io
