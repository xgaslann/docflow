"""Type definitions for DocFlow SDK."""

from dataclasses import dataclass, field
from enum import Enum
from typing import Any, Dict, List, Optional, Union


# ============== Enums ==============

class LLMProcessingMode(Enum):
    """LLM processing modes for RAG."""
    
    IMAGES = "images"       # Process images only
    TABLES = "tables"       # Process tables only
    TEXT = "text"           # Process important text sections
    ALL = "all"             # Process everything


class ChunkingStrategy(Enum):
    """Chunking strategies for RAG."""
    
    SIMPLE = "simple"                       # Character-based chunking
    HEADING_AWARE = "heading_aware"         # Respects document headings
    DOCUMENT_INTELLIGENCE = "doc_intel"     # Azure/AWS Document Intelligence
    SEMANTIC = "semantic"                   # Embedding-based semantic chunking


class OutputFormat(Enum):
    """Output format options."""
    
    MARKDOWN = "markdown"
    PDF = "pdf"
    HTML = "html"


class JobStatus(Enum):
    """Status of a processing job."""
    
    PENDING = "pending"
    PROCESSING = "processing"
    COMPLETED = "completed"
    FAILED = "failed"


# ============== Basic Types ==============

@dataclass
class MDFile:
    """Represents a Markdown file to be converted."""
    
    name: str
    content: str
    id: str = ""
    order: int = 0
    
    def __post_init__(self):
        if not self.id:
            self.id = self.name
    
    def to_dict(self) -> dict:
        """Convert to dictionary for API request."""
        return {
            "id": self.id,
            "name": self.name,
            "content": self.content,
            "order": self.order,
        }


@dataclass
class ConvertOptions:
    """Options for MD to PDF conversion."""
    
    merge_mode: str = "separate"
    output_name: Optional[str] = None
    output_format: OutputFormat = OutputFormat.PDF


@dataclass
class PDFResult:
    """Result of a PDF conversion."""
    
    success: bool
    file_paths: List[str] = field(default_factory=list)
    bytes_data: Optional[bytes] = None
    error: Optional[str] = None


@dataclass
class ExtractResult:
    """Result of a PDF to Markdown extraction."""
    
    success: bool
    markdown: str = ""
    file_path: Optional[str] = None
    page_count: int = 0
    error: Optional[str] = None


@dataclass
class MDResult:
    """Result of a PDF to Markdown extraction (API client compatible)."""
    
    success: bool
    markdown: str = ""
    file_path: Optional[str] = None
    file_name: str = ""
    error: Optional[str] = None


@dataclass
class ConvertResult:
    """Result of a format conversion."""
    
    success: bool
    content: Union[str, bytes] = ""
    format: str = ""
    error: Optional[str] = None
    images: List["ExtractedImage"] = field(default_factory=list)
    tables: List["ExtractedTable"] = field(default_factory=list)
    metadata: Dict[str, Any] = field(default_factory=dict)


@dataclass
class ExtractedImage:
    """Represents an extracted image from a document."""
    
    data: bytes
    format: str  # png, jpg, etc.
    filename: str
    caption: Optional[str] = None
    page: int = 0
    position: Optional[tuple] = None  # (x, y)
    surrounding_text: str = ""
    description: str = ""  # LLM-generated description
    llm_analysis: Optional[Dict[str, Any]] = None  # Full LLM analysis


@dataclass
class ExtractedTable:
    """Represents an extracted table from a document."""
    
    rows: List[List[str]]
    header: List[str] = field(default_factory=list)
    caption: Optional[str] = None
    page: int = 0
    summary: str = ""  # LLM-generated summary
    llm_analysis: Optional[Dict[str, Any]] = None  # Full LLM analysis


@dataclass
class HealthResponse:
    """Health check response from the server."""
    
    status: str
    version: str
    timestamp: int


@dataclass
class PreviewResult:
    """Result of a markdown preview."""
    
    html: str


# ============== Enhanced Metadata Types ==============

@dataclass
class HeadingInfo:
    """Information about a document heading."""
    
    text: str
    level: int  # 1-6
    start_pos: int
    end_pos: int
    parent_index: Optional[int] = None
    children_indices: List[int] = field(default_factory=list)


@dataclass
class TOCItem:
    """Table of contents item."""
    
    title: str
    level: int
    anchor: str
    page: int = 0
    children: List["TOCItem"] = field(default_factory=list)


@dataclass
class DocumentMetadata:
    """Comprehensive document metadata."""
    
    # Basic info
    title: str = ""
    author: str = ""
    created_date: Optional[str] = None
    modified_date: Optional[str] = None
    
    # Structure
    headings: List[HeadingInfo] = field(default_factory=list)
    heading_tree: Dict[str, Any] = field(default_factory=dict)
    table_of_contents: List[TOCItem] = field(default_factory=list)
    
    # Statistics
    word_count: int = 0
    char_count: int = 0
    page_count: int = 0
    image_count: int = 0
    table_count: int = 0
    
    # Language & Content
    language: str = ""
    keywords: List[str] = field(default_factory=list)
    entities: List[str] = field(default_factory=list)  # LLM-extracted entities
    
    # LLM-generated
    summary: str = ""  # LLM-generated document summary
    key_points: List[str] = field(default_factory=list)  # LLM-extracted key points


# ============== LLM & RAG Config ==============

@dataclass
class LLMConfig:
    """Configuration for LLM integration."""
    
    provider: str = "openai"  # openai, anthropic, ollama, azure
    model: str = "gpt-4-vision-preview"
    api_key: str = ""
    base_url: Optional[str] = None
    timeout: int = 60
    max_tokens: int = 1000
    
    # Vision-specific
    detail: str = "auto"  # auto, low, high
    
    # Processing options
    temperature: float = 0.7
    retry_count: int = 3


@dataclass
class DocIntelConfig:
    """Configuration for Document Intelligence services."""
    
    provider: str = "azure"  # azure, aws
    endpoint: str = ""
    api_key: str = ""
    model_id: str = "prebuilt-document"  # Azure model ID
    
    # AWS Textract options
    aws_region: str = "us-east-1"


@dataclass
class RAGConfig:
    """Configuration for RAG mode extraction."""
    
    enabled: bool = True
    
    # Output
    output_format: OutputFormat = OutputFormat.MARKDOWN
    
    # Chunking
    chunk_size: int = 1000
    chunk_overlap: int = 200
    chunking_strategy: ChunkingStrategy = ChunkingStrategy.HEADING_AWARE
    doc_intel_config: Optional[DocIntelConfig] = None
    
    # Extraction options
    extract_images: bool = True
    extract_tables: bool = True
    preserve_metadata: bool = True
    extract_headings: bool = True
    generate_toc: bool = True
    
    # LLM Processing - NEW: Multiple modes supported
    llm_processing: List[LLMProcessingMode] = field(default_factory=list)
    llm_config: Optional[LLMConfig] = None
    
    # Legacy compatibility
    describe_images: bool = False  # Deprecated, use llm_processing
    
    # Chunking behavior
    respect_headings: bool = True
    keep_tables_together: bool = True
    add_chunk_markers: bool = True
    
    # Parallel processing
    max_workers: int = 4
    
    def __post_init__(self):
        # Legacy compatibility: if describe_images is True, add IMAGES to llm_processing
        if self.describe_images and LLMProcessingMode.IMAGES not in self.llm_processing:
            self.llm_processing.append(LLMProcessingMode.IMAGES)


# ============== Chunk Types ==============

@dataclass
class ChunkMetadata:
    """Metadata for a chunk."""
    
    section_title: str = ""
    heading_path: List[str] = field(default_factory=list)
    heading_levels: List[int] = field(default_factory=list)  # Corresponding levels
    has_table: bool = False
    has_image: bool = False
    has_code: bool = False
    page: int = 0
    
    # Position in document hierarchy
    section_index: int = 0
    subsection_index: int = 0
    
    # Content type hints
    content_type: str = "text"  # text, table, code, mixed


@dataclass
class Chunk:
    """Represents a chunk of content for RAG."""
    
    content: str
    index: int
    start_char: int
    end_char: int
    metadata: Optional[ChunkMetadata] = None
    
    # Embeddings (if generated)
    embedding: Optional[List[float]] = None


# ============== RAG Document ==============

@dataclass
class RAGDocument:
    """Complete RAG-processed document."""
    
    content: str  # Full markdown content
    chunks: List[Chunk] = field(default_factory=list)
    images: List[ExtractedImage] = field(default_factory=list)
    tables: List[ExtractedTable] = field(default_factory=list)
    
    # Enhanced metadata
    metadata: DocumentMetadata = field(default_factory=DocumentMetadata)
    raw_metadata: Dict[str, Any] = field(default_factory=dict)
    
    # Source info
    source_file: str = ""
    source_format: str = ""
    
    # Output
    pdf_bytes: Optional[bytes] = None  # If PDF output requested


# ============== Batch Processing Types ==============

@dataclass
class BatchJob:
    """Represents a batch processing job."""
    
    job_id: str
    status: JobStatus = JobStatus.PENDING
    total_files: int = 0
    processed_files: int = 0
    failed_files: int = 0
    results: List[RAGDocument] = field(default_factory=list)
    errors: Dict[str, str] = field(default_factory=dict)  # filename -> error
    created_at: Optional[str] = None
    completed_at: Optional[str] = None


@dataclass
class BatchConfig:
    """Configuration for batch processing."""
    
    max_workers: int = 4
    fail_fast: bool = False  # Stop on first error
    continue_on_error: bool = True
    timeout_per_file: int = 300  # seconds
    
    # Queue settings
    queue_size: int = 100
    retry_failed: bool = True
    max_retries: int = 3

