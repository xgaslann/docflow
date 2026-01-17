"""Base vector store interface."""

from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from enum import Enum
from typing import Any, Dict, List, Optional, Tuple

from ...types import Chunk


class EmbeddingProvider(Enum):
    """Embedding provider."""
    
    OPENAI = "openai"
    AZURE_OPENAI = "azure_openai"
    COHERE = "cohere"
    HUGGINGFACE = "huggingface"
    LOCAL = "local"


class DistanceMetric(Enum):
    """Distance metric for vector similarity."""
    
    COSINE = "cosine"
    EUCLIDEAN = "euclidean"
    DOT_PRODUCT = "dot"
    MANHATTAN = "manhattan"


class IndexType(Enum):
    """Vector index type."""
    
    HNSW = "hnsw"
    IVFFLAT = "ivfflat"
    FLAT = "flat"
    ANNOY = "annoy"


@dataclass
class VectorStoreConfig:
    """Base configuration for vector stores.
    
    Recommendations:
        - embedding_model: "text-embedding-3-small" for cost/performance balance
        - distance_metric: COSINE for text similarity
        - index_type: HNSW for best query performance
    """
    
    # Connection
    connection_string: str = ""
    database: str = "docflow"
    collection: str = "chunks"
    
    # Embedding
    embedding_provider: EmbeddingProvider = EmbeddingProvider.OPENAI
    embedding_model: str = "text-embedding-3-small"
    embedding_api_key: str = ""
    embedding_dimensions: int = 1536
    embedding_batch_size: int = 100
    
    # Index
    index_type: IndexType = IndexType.HNSW
    distance_metric: DistanceMetric = DistanceMetric.COSINE
    
    # Performance
    pool_size: int = 5
    timeout: int = 30
    
    # Metadata
    store_metadata: bool = True
    metadata_fields: List[str] = field(default_factory=lambda: [
        "source_file", "chunk_index", "section_title"
    ])


@dataclass
class SearchResult:
    """Result from vector search."""
    
    chunk: Chunk
    score: float
    metadata: Dict[str, Any] = field(default_factory=dict)


class VectorStore(ABC):
    """Abstract base class for vector stores.
    
    Provides interface for storing and retrieving document chunks
    with vector embeddings for semantic search.
    """
    
    def __init__(self, config: VectorStoreConfig) -> None:
        self.config = config
        self._embedder = None
    
    @abstractmethod
    async def connect(self) -> None:
        """Connect to the vector store."""
        pass
    
    @abstractmethod
    async def disconnect(self) -> None:
        """Disconnect from the vector store."""
        pass
    
    @abstractmethod
    async def create_collection(self, name: Optional[str] = None) -> None:
        """Create a new collection/table for storing vectors."""
        pass
    
    @abstractmethod
    async def drop_collection(self, name: Optional[str] = None) -> None:
        """Drop a collection/table."""
        pass
    
    @abstractmethod
    async def insert(
        self,
        chunks: List[Chunk],
        embeddings: Optional[List[List[float]]] = None,
    ) -> List[str]:
        """Insert chunks with embeddings.
        
        Args:
            chunks: Document chunks to insert.
            embeddings: Pre-computed embeddings (optional).
            
        Returns:
            List of inserted IDs.
        """
        pass
    
    @abstractmethod
    async def search(
        self,
        query: str,
        top_k: int = 5,
        filter_metadata: Optional[Dict[str, Any]] = None,
    ) -> List[SearchResult]:
        """Search for similar chunks.
        
        Args:
            query: Search query text.
            top_k: Number of results to return.
            filter_metadata: Metadata filters.
            
        Returns:
            List of search results with scores.
        """
        pass
    
    @abstractmethod
    async def search_by_vector(
        self,
        embedding: List[float],
        top_k: int = 5,
        filter_metadata: Optional[Dict[str, Any]] = None,
    ) -> List[SearchResult]:
        """Search using a vector directly."""
        pass
    
    @abstractmethod
    async def delete(
        self,
        ids: Optional[List[str]] = None,
        filter_metadata: Optional[Dict[str, Any]] = None,
    ) -> int:
        """Delete chunks by ID or metadata filter.
        
        Returns:
            Number of deleted chunks.
        """
        pass
    
    @abstractmethod
    async def update(
        self,
        id: str,
        chunk: Optional[Chunk] = None,
        embedding: Optional[List[float]] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> bool:
        """Update a chunk."""
        pass
    
    @abstractmethod
    async def count(self, filter_metadata: Optional[Dict[str, Any]] = None) -> int:
        """Count chunks in collection."""
        pass
    
    # ============== Embedding Generation ==============
    
    async def generate_embedding(self, text: str) -> List[float]:
        """Generate embedding for text."""
        embeddings = await self.generate_embeddings([text])
        return embeddings[0]
    
    async def generate_embeddings(self, texts: List[str]) -> List[List[float]]:
        """Generate embeddings for multiple texts."""
        if self.config.embedding_provider == EmbeddingProvider.OPENAI:
            return await self._generate_openai_embeddings(texts)
        elif self.config.embedding_provider == EmbeddingProvider.AZURE_OPENAI:
            return await self._generate_azure_embeddings(texts)
        elif self.config.embedding_provider == EmbeddingProvider.COHERE:
            return await self._generate_cohere_embeddings(texts)
        else:
            raise ValueError(f"Unsupported embedding provider: {self.config.embedding_provider}")
    
    async def _generate_openai_embeddings(self, texts: List[str]) -> List[List[float]]:
        """Generate embeddings using OpenAI."""
        try:
            from openai import AsyncOpenAI
        except ImportError:
            raise ImportError("openai is required: pip install openai")
        
        client = AsyncOpenAI(api_key=self.config.embedding_api_key)
        
        all_embeddings = []
        for i in range(0, len(texts), self.config.embedding_batch_size):
            batch = texts[i:i + self.config.embedding_batch_size]
            response = await client.embeddings.create(
                model=self.config.embedding_model,
                input=batch,
            )
            all_embeddings.extend([d.embedding for d in response.data])
        
        return all_embeddings
    
    async def _generate_azure_embeddings(self, texts: List[str]) -> List[List[float]]:
        """Generate embeddings using Azure OpenAI."""
        try:
            from openai import AsyncAzureOpenAI
        except ImportError:
            raise ImportError("openai is required: pip install openai")
        
        # Azure-specific implementation
        raise NotImplementedError("Azure OpenAI embeddings not yet implemented")
    
    async def _generate_cohere_embeddings(self, texts: List[str]) -> List[List[float]]:
        """Generate embeddings using Cohere."""
        try:
            import cohere
        except ImportError:
            raise ImportError("cohere is required: pip install cohere")
        
        client = cohere.Client(self.config.embedding_api_key)
        response = client.embed(
            texts=texts,
            model=self.config.embedding_model,
            input_type="search_document",
        )
        
        return response.embeddings
