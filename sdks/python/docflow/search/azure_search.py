"""Azure AI Search integration."""

from dataclasses import dataclass, field
from enum import Enum
from typing import Any, Dict, List, Optional
import json

from ..types import Chunk, ChunkMetadata


class QueryType(Enum):
    """Search query type."""
    
    SIMPLE = "simple"
    FULL = "full"
    SEMANTIC = "semantic"


class SearchMode(Enum):
    """Search mode."""
    
    ANY = "any"
    ALL = "all"


class VectorSearchProfile(Enum):
    """Vector search algorithm profile."""
    
    HNSW = "hnsw"
    EXHAUSTIVE_KNN = "exhaustive-knn"


@dataclass
class AISearchConfig:
    """Configuration for Azure AI Search.
    
    All available options exposed for full control.
    
    Example:
        >>> config = AISearchConfig(
        ...     endpoint="https://xxx.search.windows.net",
        ...     api_key="your-key",
        ...     index_name="docflow-index",
        ...     query_type=QueryType.SEMANTIC,
        ...     hybrid_search=True
        ... )
    
    Recommendations:
        - query_type: SEMANTIC for best quality
        - hybrid_search: True for balanced keyword + semantic
        - k_nearest_neighbors: 50 for good recall/speed balance
    """
    
    # Connection
    endpoint: str = ""
    api_key: str = ""
    api_version: str = "2024-07-01"
    
    # Index
    index_name: str = "docflow-index"
    
    # ============== Index Schema ==============
    
    # Vector configuration
    vector_search_profile: str = "default-vector-profile"
    vector_algorithm: VectorSearchProfile = VectorSearchProfile.HNSW
    
    # HNSW parameters
    hnsw_m: int = 4
    hnsw_ef_construction: int = 400
    hnsw_ef_search: int = 500
    metric: str = "cosine"  # cosine, euclidean, dotProduct
    
    # Semantic configuration
    semantic_config: str = "default-semantic-config"
    semantic_prioritized_fields: List[str] = field(default_factory=lambda: [
        "content", "title"
    ])
    
    # Scoring profile
    scoring_profile: str = ""
    
    # ============== Search Options ==============
    
    # Query type
    query_type: QueryType = QueryType.SEMANTIC
    search_mode: SearchMode = SearchMode.ANY
    
    # Pagination
    top: int = 10
    skip: int = 0
    
    # Vector search
    vector_fields: List[str] = field(default_factory=lambda: ["content_vector"])
    k_nearest_neighbors: int = 50
    
    # Hybrid search
    hybrid_search: bool = True
    max_text_recall_size: int = 1000
    
    # Semantic reranking
    semantic_reranking: bool = True
    semantic_max_wait: int = 700  # milliseconds
    
    # ============== Filters & Fields ==============
    
    filter_expression: str = ""
    search_fields: List[str] = field(default_factory=list)
    select_fields: List[str] = field(default_factory=lambda: [
        "id", "content", "title", "metadata"
    ])
    order_by: List[str] = field(default_factory=list)
    
    # Facets
    facets: List[str] = field(default_factory=list)
    
    # Highlighting
    highlight_fields: List[str] = field(default_factory=list)
    highlight_pre_tag: str = "<em>"
    highlight_post_tag: str = "</em>"
    
    # ============== Embedding ==============
    
    embedding_model: str = "text-embedding-3-small"
    embedding_dimensions: int = 1536
    embedding_api_key: str = ""
    
    def validate(self) -> None:
        """Validate configuration."""
        if not self.endpoint:
            raise ValueError("Azure AI Search endpoint is required")
        if not self.api_key:
            raise ValueError("Azure AI Search API key is required")


@dataclass
class SearchResult:
    """Result from AI Search."""
    
    id: str
    content: str
    score: float
    semantic_score: Optional[float] = None
    highlights: Dict[str, List[str]] = field(default_factory=dict)
    metadata: Dict[str, Any] = field(default_factory=dict)


class AzureAISearchClient:
    """Azure AI Search client.
    
    Supports vector search, semantic search, and hybrid search.
    
    Example:
        >>> config = AISearchConfig(
        ...     endpoint="https://xxx.search.windows.net",
        ...     api_key="your-key",
        ...     index_name="docflow-index"
        ... )
        >>> client = AzureAISearchClient(config)
        >>> 
        >>> # Create index
        >>> await client.create_index()
        >>> 
        >>> # Upload documents
        >>> await client.upload_chunks(chunks)
        >>> 
        >>> # Search
        >>> results = await client.search("query", top_k=5)
    """
    
    def __init__(self, config: AISearchConfig) -> None:
        self.config = config
        self._client = None
        self._index_client = None
    
    async def connect(self) -> None:
        """Initialize search client."""
        try:
            from azure.search.documents.aio import SearchClient
            from azure.search.documents.indexes.aio import SearchIndexClient
            from azure.core.credentials import AzureKeyCredential
        except ImportError:
            raise ImportError("azure-search-documents is required: pip install azure-search-documents")
        
        credential = AzureKeyCredential(self.config.api_key)
        
        self._client = SearchClient(
            endpoint=self.config.endpoint,
            index_name=self.config.index_name,
            credential=credential,
        )
        
        self._index_client = SearchIndexClient(
            endpoint=self.config.endpoint,
            credential=credential,
        )
    
    async def disconnect(self) -> None:
        """Close client connections."""
        if self._client:
            await self._client.close()
        if self._index_client:
            await self._index_client.close()
    
    async def create_index(self) -> None:
        """Create search index with vector and semantic configuration."""
        from azure.search.documents.indexes.models import (
            SearchIndex,
            SearchField,
            SearchFieldDataType,
            VectorSearch,
            HnswAlgorithmConfiguration,
            VectorSearchProfile,
            SemanticConfiguration,
            SemanticField,
            SemanticPrioritizedFields,
            SemanticSearch,
        )
        
        # Define fields
        fields = [
            SearchField(
                name="id",
                type=SearchFieldDataType.String,
                key=True,
                filterable=True,
            ),
            SearchField(
                name="content",
                type=SearchFieldDataType.String,
                searchable=True,
            ),
            SearchField(
                name="title",
                type=SearchFieldDataType.String,
                searchable=True,
                filterable=True,
            ),
            SearchField(
                name="content_vector",
                type=SearchFieldDataType.Collection(SearchFieldDataType.Single),
                searchable=True,
                vector_search_dimensions=self.config.embedding_dimensions,
                vector_search_profile_name=self.config.vector_search_profile,
            ),
            SearchField(
                name="chunk_index",
                type=SearchFieldDataType.Int32,
                filterable=True,
                sortable=True,
            ),
            SearchField(
                name="source_file",
                type=SearchFieldDataType.String,
                filterable=True,
                facetable=True,
            ),
            SearchField(
                name="metadata",
                type=SearchFieldDataType.String,
            ),
        ]
        
        # Vector search configuration
        vector_search = VectorSearch(
            algorithms=[
                HnswAlgorithmConfiguration(
                    name="hnsw-config",
                    parameters={
                        "m": self.config.hnsw_m,
                        "efConstruction": self.config.hnsw_ef_construction,
                        "efSearch": self.config.hnsw_ef_search,
                        "metric": self.config.metric,
                    }
                ),
            ],
            profiles=[
                VectorSearchProfile(
                    name=self.config.vector_search_profile,
                    algorithm_configuration_name="hnsw-config",
                ),
            ],
        )
        
        # Semantic search configuration
        semantic_search = SemanticSearch(
            configurations=[
                SemanticConfiguration(
                    name=self.config.semantic_config,
                    prioritized_fields=SemanticPrioritizedFields(
                        content_fields=[
                            SemanticField(field_name=f)
                            for f in self.config.semantic_prioritized_fields
                        ],
                    ),
                ),
            ],
        )
        
        # Create index
        index = SearchIndex(
            name=self.config.index_name,
            fields=fields,
            vector_search=vector_search,
            semantic_search=semantic_search,
        )
        
        await self._index_client.create_or_update_index(index)
    
    async def delete_index(self) -> None:
        """Delete the search index."""
        await self._index_client.delete_index(self.config.index_name)
    
    async def upload_chunks(
        self,
        chunks: List[Chunk],
        embeddings: Optional[List[List[float]]] = None,
        source_file: str = "",
    ) -> List[str]:
        """Upload chunks to the index."""
        if embeddings is None:
            embeddings = await self._generate_embeddings([c.content for c in chunks])
        
        documents = []
        ids = []
        
        for i, (chunk, embedding) in enumerate(zip(chunks, embeddings)):
            doc_id = f"{source_file}_{chunk.index}" if source_file else f"chunk_{i}"
            
            doc = {
                "id": doc_id,
                "content": chunk.content,
                "title": chunk.metadata.section_title if chunk.metadata else "",
                "content_vector": embedding,
                "chunk_index": chunk.index,
                "source_file": source_file,
                "metadata": json.dumps({
                    "heading_path": chunk.metadata.heading_path if chunk.metadata else [],
                    "has_table": chunk.metadata.has_table if chunk.metadata else False,
                    "has_image": chunk.metadata.has_image if chunk.metadata else False,
                }),
            }
            
            documents.append(doc)
            ids.append(doc_id)
        
        await self._client.upload_documents(documents)
        return ids
    
    async def search(
        self,
        query: str,
        top_k: int = 5,
        filter_expression: Optional[str] = None,
        **kwargs,
    ) -> List[SearchResult]:
        """Search the index.
        
        Combines vector search, semantic search, and keyword search
        based on configuration.
        """
        from azure.search.documents.models import VectorizedQuery
        
        # Generate query embedding
        query_embedding = await self._generate_embedding(query)
        
        # Build search parameters
        search_params = {
            "search_text": query if not self.config.hybrid_search else query,
            "top": kwargs.get("top", self.config.top),
            "skip": kwargs.get("skip", self.config.skip),
            "select": kwargs.get("select", self.config.select_fields),
        }
        
        # Vector search
        search_params["vector_queries"] = [
            VectorizedQuery(
                vector=query_embedding,
                k_nearest_neighbors=kwargs.get("k", self.config.k_nearest_neighbors),
                fields=",".join(self.config.vector_fields),
            )
        ]
        
        # Semantic search
        if self.config.query_type == QueryType.SEMANTIC:
            search_params["query_type"] = "semantic"
            search_params["semantic_configuration_name"] = self.config.semantic_config
        
        # Filter
        if filter_expression or self.config.filter_expression:
            search_params["filter"] = filter_expression or self.config.filter_expression
        
        # Highlighting
        if self.config.highlight_fields:
            search_params["highlight_fields"] = ",".join(self.config.highlight_fields)
        
        # Execute search
        results = []
        async for result in self._client.search(**search_params):
            metadata = {}
            if result.get("metadata"):
                try:
                    metadata = json.loads(result["metadata"])
                except:
                    pass
            
            search_result = SearchResult(
                id=result["id"],
                content=result.get("content", ""),
                score=result.get("@search.score", 0.0),
                semantic_score=result.get("@search.reranker_score"),
                highlights=dict(result.get("@search.highlights", {})),
                metadata=metadata,
            )
            results.append(search_result)
        
        return results[:top_k]
    
    async def delete_documents(self, ids: List[str]) -> None:
        """Delete documents by ID."""
        documents = [{"id": doc_id} for doc_id in ids]
        await self._client.delete_documents(documents)
    
    async def _generate_embedding(self, text: str) -> List[float]:
        """Generate embedding for text."""
        embeddings = await self._generate_embeddings([text])
        return embeddings[0]
    
    async def _generate_embeddings(self, texts: List[str]) -> List[List[float]]:
        """Generate embeddings for texts."""
        try:
            from openai import AsyncOpenAI
        except ImportError:
            raise ImportError("openai is required: pip install openai")
        
        client = AsyncOpenAI(api_key=self.config.embedding_api_key)
        
        response = await client.embeddings.create(
            model=self.config.embedding_model,
            input=texts,
        )
        
        return [d.embedding for d in response.data]
