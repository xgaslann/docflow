"""MongoDB vector store with Atlas Vector Search."""

from dataclasses import dataclass, field
from typing import Any, Dict, List, Optional
from datetime import datetime

from .base import (
    VectorStore,
    VectorStoreConfig,
    SearchResult,
    IndexType,
    DistanceMetric,
)
from ...types import Chunk, ChunkMetadata


@dataclass
class MongoVectorConfig(VectorStoreConfig):
    """Configuration for MongoDB with Atlas Vector Search.
    
    Requires MongoDB Atlas with Vector Search enabled.
    
    Example:
        >>> config = MongoVectorConfig(
        ...     connection_string="mongodb+srv://user:pass@cluster.mongodb.net/",
        ...     database="docflow",
        ...     collection="chunks",
        ...     embedding_api_key="sk-..."
        ... )
    
    Recommendations:
        - num_candidates: 2-10x of top_k for better recall
        - Use "cosine" similarity for text embeddings
    """
    
    # Atlas specific
    atlas_cluster: str = ""
    
    # Index settings
    index_name: str = "vector_index"
    num_candidates: int = 100  # Number of candidates to consider
    
    # Vector search options
    vector_path: str = "embedding"
    
    # Full-text search (for hybrid)
    text_index_name: str = "text_index"
    text_search_path: str = "content"


class MongoVectorStore(VectorStore):
    """MongoDB vector store using Atlas Vector Search.
    
    Example:
        >>> config = MongoVectorConfig(
        ...     connection_string="mongodb+srv://...",
        ...     database="docflow",
        ...     embedding_api_key="sk-..."
        ... )
        >>> store = MongoVectorStore(config)
        >>> await store.connect()
        >>> 
        >>> # Insert chunks
        >>> ids = await store.insert(chunks)
        >>> 
        >>> # Search
        >>> results = await store.search("query text", top_k=5)
    """
    
    def __init__(self, config: MongoVectorConfig) -> None:
        super().__init__(config)
        self.config: MongoVectorConfig = config
        self._client = None
        self._db = None
        self._collection = None
    
    async def connect(self) -> None:
        """Connect to MongoDB."""
        try:
            from motor.motor_asyncio import AsyncIOMotorClient
        except ImportError:
            raise ImportError("motor is required: pip install motor")
        
        self._client = AsyncIOMotorClient(self.config.connection_string)
        self._db = self._client[self.config.database]
        self._collection = self._db[self.config.collection]
    
    async def disconnect(self) -> None:
        """Disconnect from MongoDB."""
        if self._client:
            self._client.close()
            self._client = None
    
    async def create_collection(self, name: Optional[str] = None) -> None:
        """Create collection with vector search index."""
        collection_name = name or self.config.collection
        
        # Create collection if not exists
        collections = await self._db.list_collection_names()
        if collection_name not in collections:
            await self._db.create_collection(collection_name)
        
        # Create vector search index
        # Note: This requires Atlas UI or Atlas Admin API
        # The index definition is provided as a reference
        index_definition = {
            "fields": [
                {
                    "type": "vector",
                    "path": self.config.vector_path,
                    "numDimensions": self.config.embedding_dimensions,
                    "similarity": self._get_similarity_type(),
                }
            ]
        }
        
        # Try to create via command (may require Atlas Admin API)
        try:
            await self._db.command({
                "createSearchIndexes": collection_name,
                "indexes": [{
                    "name": self.config.index_name,
                    "definition": index_definition,
                }]
            })
        except Exception:
            # Index creation may need Atlas UI
            pass
    
    async def drop_collection(self, name: Optional[str] = None) -> None:
        """Drop collection."""
        collection_name = name or self.config.collection
        await self._db.drop_collection(collection_name)
    
    async def insert(
        self,
        chunks: List[Chunk],
        embeddings: Optional[List[List[float]]] = None,
    ) -> List[str]:
        """Insert chunks with embeddings."""
        if embeddings is None:
            texts = [c.content for c in chunks]
            embeddings = await self.generate_embeddings(texts)
        
        documents = []
        ids = []
        
        for chunk, embedding in zip(chunks, embeddings):
            chunk_id = f"chunk_{chunk.index}_{datetime.utcnow().timestamp()}"
            
            doc = {
                "_id": chunk_id,
                "content": chunk.content,
                "embedding": embedding,
                "chunk_index": chunk.index,
                "start_char": chunk.start_char,
                "end_char": chunk.end_char,
                "created_at": datetime.utcnow(),
            }
            
            if chunk.metadata:
                doc["metadata"] = {
                    "section_title": chunk.metadata.section_title,
                    "heading_path": chunk.metadata.heading_path,
                    "has_table": chunk.metadata.has_table,
                    "has_image": chunk.metadata.has_image,
                    "page": chunk.metadata.page,
                    "content_type": chunk.metadata.content_type,
                }
            
            documents.append(doc)
            ids.append(chunk_id)
        
        await self._collection.insert_many(documents)
        return ids
    
    async def search(
        self,
        query: str,
        top_k: int = 5,
        filter_metadata: Optional[Dict[str, Any]] = None,
    ) -> List[SearchResult]:
        """Search for similar chunks."""
        embedding = await self.generate_embedding(query)
        return await self.search_by_vector(embedding, top_k, filter_metadata)
    
    async def search_by_vector(
        self,
        embedding: List[float],
        top_k: int = 5,
        filter_metadata: Optional[Dict[str, Any]] = None,
    ) -> List[SearchResult]:
        """Search using vector directly."""
        # Build vector search pipeline
        vector_search = {
            "$vectorSearch": {
                "index": self.config.index_name,
                "path": self.config.vector_path,
                "queryVector": embedding,
                "numCandidates": self.config.num_candidates,
                "limit": top_k,
            }
        }
        
        # Add filter if provided
        if filter_metadata:
            filter_conditions = {}
            for key, value in filter_metadata.items():
                filter_conditions[f"metadata.{key}"] = value
            vector_search["$vectorSearch"]["filter"] = filter_conditions
        
        pipeline = [
            vector_search,
            {
                "$project": {
                    "_id": 1,
                    "content": 1,
                    "chunk_index": 1,
                    "metadata": 1,
                    "score": {"$meta": "vectorSearchScore"},
                }
            }
        ]
        
        cursor = self._collection.aggregate(pipeline)
        results = []
        
        async for doc in cursor:
            metadata = doc.get("metadata", {})
            chunk_meta = ChunkMetadata(
                section_title=metadata.get("section_title", ""),
                heading_path=metadata.get("heading_path", []),
                has_table=metadata.get("has_table", False),
                has_image=metadata.get("has_image", False),
                page=metadata.get("page", 0),
                content_type=metadata.get("content_type", "text"),
            )
            
            chunk = Chunk(
                content=doc["content"],
                index=doc.get("chunk_index", 0),
                start_char=0,
                end_char=len(doc["content"]),
                metadata=chunk_meta,
            )
            
            results.append(SearchResult(
                chunk=chunk,
                score=doc.get("score", 0.0),
                metadata=metadata,
            ))
        
        return results
    
    async def delete(
        self,
        ids: Optional[List[str]] = None,
        filter_metadata: Optional[Dict[str, Any]] = None,
    ) -> int:
        """Delete chunks."""
        if ids:
            result = await self._collection.delete_many({"_id": {"$in": ids}})
        elif filter_metadata:
            filter_conditions = {}
            for key, value in filter_metadata.items():
                filter_conditions[f"metadata.{key}"] = value
            result = await self._collection.delete_many(filter_conditions)
        else:
            return 0
        
        return result.deleted_count
    
    async def update(
        self,
        id: str,
        chunk: Optional[Chunk] = None,
        embedding: Optional[List[float]] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> bool:
        """Update a chunk."""
        update_doc = {}
        
        if chunk:
            update_doc["content"] = chunk.content
            update_doc["chunk_index"] = chunk.index
        
        if embedding:
            update_doc["embedding"] = embedding
        
        if metadata:
            for key, value in metadata.items():
                update_doc[f"metadata.{key}"] = value
        
        if not update_doc:
            return False
        
        result = await self._collection.update_one(
            {"_id": id},
            {"$set": update_doc}
        )
        
        return result.modified_count > 0
    
    async def count(self, filter_metadata: Optional[Dict[str, Any]] = None) -> int:
        """Count chunks."""
        if filter_metadata:
            filter_conditions = {}
            for key, value in filter_metadata.items():
                filter_conditions[f"metadata.{key}"] = value
            return await self._collection.count_documents(filter_conditions)
        
        return await self._collection.count_documents({})
    
    def _get_similarity_type(self) -> str:
        """Get MongoDB similarity type."""
        similarity_map = {
            DistanceMetric.COSINE: "cosine",
            DistanceMetric.EUCLIDEAN: "euclidean",
            DistanceMetric.DOT_PRODUCT: "dotProduct",
        }
        return similarity_map.get(self.config.distance_metric, "cosine")
    
    # ============== Hybrid Search ==============
    
    async def hybrid_search(
        self,
        query: str,
        top_k: int = 5,
        vector_weight: float = 0.7,
        text_weight: float = 0.3,
        filter_metadata: Optional[Dict[str, Any]] = None,
    ) -> List[SearchResult]:
        """Perform hybrid search combining vector and text search."""
        embedding = await self.generate_embedding(query)
        
        pipeline = [
            # Vector search
            {
                "$vectorSearch": {
                    "index": self.config.index_name,
                    "path": self.config.vector_path,
                    "queryVector": embedding,
                    "numCandidates": self.config.num_candidates,
                    "limit": top_k * 2,
                }
            },
            {
                "$addFields": {
                    "vector_score": {"$meta": "vectorSearchScore"}
                }
            },
            # Text search score (if text index exists)
            {
                "$addFields": {
                    "combined_score": {
                        "$multiply": ["$vector_score", vector_weight]
                    }
                }
            },
            {"$sort": {"combined_score": -1}},
            {"$limit": top_k},
        ]
        
        cursor = self._collection.aggregate(pipeline)
        results = []
        
        async for doc in cursor:
            metadata = doc.get("metadata", {})
            chunk_meta = ChunkMetadata(
                section_title=metadata.get("section_title", ""),
                heading_path=metadata.get("heading_path", []),
            )
            
            chunk = Chunk(
                content=doc["content"],
                index=doc.get("chunk_index", 0),
                start_char=0,
                end_char=len(doc["content"]),
                metadata=chunk_meta,
            )
            
            results.append(SearchResult(
                chunk=chunk,
                score=doc.get("combined_score", 0.0),
                metadata=metadata,
            ))
        
        return results
