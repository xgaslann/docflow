"""
DocFlow Python SDK

A standalone Python library for converting between document formats,
optimized for RAG (Retrieval-Augmented Generation) systems.

No server required.

Basic Usage:
    >>> from docflow import Converter, MDFile
    >>> converter = Converter()
    >>> result = converter.convert_to_pdf([
    ...     MDFile(name="doc.md", content="# Hello World")
    ... ])

Multi-Format Conversion:
    >>> from docflow.formats import DOCXConverter, ExcelConverter
    >>> docx_conv = DOCXConverter()
    >>> result = docx_conv.to_markdown(docx_bytes, "document.docx")

RAG Mode:
    >>> from docflow.rag import RAGProcessor
    >>> from docflow import RAGConfig, LLMProcessingMode
    >>> 
    >>> config = RAGConfig(
    ...     llm_processing=[LLMProcessingMode.ALL],
    ...     llm_config=LLMConfig(provider="openai", api_key="...")
    ... )
    >>> processor = RAGProcessor(config)
    >>> doc = processor.process_file("document.pdf")

Batch Processing:
    >>> from docflow import BatchProcessor
    >>> batch = BatchProcessor(rag_config=config, max_workers=4)
    >>> results = batch.process_files(["doc.pdf", "data.xlsx"])
"""

from .types import (
    # Enums
    LLMProcessingMode,
    ChunkingStrategy,
    OutputFormat,
    JobStatus,
    # Basic Types
    MDFile,
    ConvertOptions,
    PDFResult,
    ExtractResult,
    MDResult,
    ConvertResult,
    ExtractedImage,
    ExtractedTable,
    HealthResponse,
    PreviewResult,
    # Metadata
    HeadingInfo,
    TOCItem,
    DocumentMetadata,
    # Config
    LLMConfig,
    DocIntelConfig,
    RAGConfig,
    BatchConfig,
    # RAG Types
    Chunk,
    ChunkMetadata,
    RAGDocument,
    BatchJob,
)
from .converter import Converter
from .extractor import Extractor
from .markdown import MarkdownParser
from .template import Template
from .storage import Storage, LocalStorage
from .batch_processor import BatchProcessor

__version__ = "0.3.0"
__all__ = [
    # Main classes
    "Converter",
    "Extractor",
    "MarkdownParser",
    "Template",
    "BatchProcessor",
    # Enums
    "LLMProcessingMode",
    "ChunkingStrategy",
    "OutputFormat",
    "JobStatus",
    # Types
    "MDFile",
    "ConvertOptions",
    "PDFResult",
    "ExtractResult",
    "MDResult",
    "ConvertResult",
    "ExtractedImage",
    "ExtractedTable",
    "HealthResponse",
    "PreviewResult",
    # Metadata
    "HeadingInfo",
    "TOCItem",
    "DocumentMetadata",
    # Config
    "LLMConfig",
    "DocIntelConfig",
    "RAGConfig",
    "BatchConfig",
    # RAG Types
    "Chunk",
    "ChunkMetadata",
    "RAGDocument",
    "BatchJob",
    # Storage
    "Storage",
    "LocalStorage",
]

# Optional cloud storage exports
try:
    from .storage import S3Storage
    __all__.append("S3Storage")
except ImportError:
    pass

try:
    from .storage import AzureStorage
    __all__.append("AzureStorage")
except ImportError:
    pass

