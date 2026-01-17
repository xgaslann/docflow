"""Chunking and retrieval configuration."""

from dataclasses import dataclass, field
from enum import Enum
from typing import List, Optional


class SplitBy(Enum):
    """Text splitting strategy."""
    
    PARAGRAPH = "paragraph"
    SENTENCE = "sentence"
    TOKEN = "token"
    CHARACTER = "character"
    HEADING = "heading"


@dataclass
class ChunkingConfig:
    """Configuration for text chunking.
    
    Example:
        >>> config = ChunkingConfig(
        ...     chunk_size=1000,
        ...     chunk_overlap=200,
        ...     split_by=SplitBy.PARAGRAPH
        ... )
    
    Recommendations:
        - chunk_size: 500-1500 for most use cases
        - chunk_overlap: 10-20% of chunk_size
        - split_by: PARAGRAPH for documents, SENTENCE for precise retrieval
    """
    
    # Size settings
    chunk_size: int = 1000
    chunk_overlap: int = 200
    min_chunk_size: int = 100
    max_chunk_size: int = 2000
    
    # Splitting strategy
    split_by: SplitBy = SplitBy.PARAGRAPH
    
    # Separators (for character/paragraph splitting)
    separators: List[str] = field(default_factory=lambda: ["\n\n", "\n", ". ", " "])
    
    # Token settings (for token splitting)
    tokenizer: str = "cl100k_base"  # tiktoken tokenizer
    
    # Heading-aware settings
    respect_headings: bool = True
    keep_tables_together: bool = True
    keep_code_together: bool = True
    
    # Markers
    add_chunk_markers: bool = True
    marker_format: str = "[CHUNK {index}]"
    
    def validate(self) -> None:
        """Validate configuration."""
        if self.chunk_size <= 0:
            raise ValueError("chunk_size must be positive")
        if self.chunk_overlap < 0:
            raise ValueError("chunk_overlap cannot be negative")
        if self.chunk_overlap >= self.chunk_size:
            raise ValueError("chunk_overlap must be less than chunk_size")
        if self.min_chunk_size > self.max_chunk_size:
            raise ValueError("min_chunk_size cannot exceed max_chunk_size")


@dataclass
class RetrievalConfig:
    """Configuration for retrieval operations.
    
    Example:
        >>> config = RetrievalConfig(
        ...     top_k=5,
        ...     similarity_threshold=0.7,
        ...     rerank=True
        ... )
    
    Recommendations:
        - top_k: 3-10 depending on context window
        - similarity_threshold: 0.6-0.8 for balanced precision/recall
        - Enable rerank for better results (slower)
    """
    
    # Basic retrieval
    top_k: int = 5
    similarity_threshold: float = 0.7
    min_score: float = 0.0
    
    # Reranking
    rerank: bool = False
    rerank_model: str = "cross-encoder/ms-marco-MiniLM-L-6-v2"
    rerank_top_k: int = 3
    
    # Filtering
    filter_duplicates: bool = True
    duplicate_threshold: float = 0.95
    
    # Context
    include_context: bool = True
    context_before: int = 1  # Number of chunks before
    context_after: int = 1   # Number of chunks after
    
    # Hybrid search
    hybrid_search: bool = False
    keyword_weight: float = 0.3
    semantic_weight: float = 0.7
    
    # MMR (Maximal Marginal Relevance)
    use_mmr: bool = False
    mmr_lambda: float = 0.5  # Balance between relevance and diversity
    
    def validate(self) -> None:
        """Validate configuration."""
        if self.top_k <= 0:
            raise ValueError("top_k must be positive")
        if not 0 <= self.similarity_threshold <= 1:
            raise ValueError("similarity_threshold must be between 0 and 1")
        if self.hybrid_search and abs(self.keyword_weight + self.semantic_weight - 1.0) > 0.01:
            raise ValueError("keyword_weight + semantic_weight must equal 1.0")
