"""PostgreSQL vector store with pgvector."""

from dataclasses import dataclass, field
from typing import Any, Dict, List, Optional
import json

from .base import (
    VectorStore,
    VectorStoreConfig,
    SearchResult,
    IndexType,
    DistanceMetric,
)
from ...types import Chunk, ChunkMetadata


@dataclass
class PostgresVectorConfig(VectorStoreConfig):
    """Configuration for PostgreSQL with pgvector.
    
    Requires PostgreSQL with pgvector extension installed.
    
    Example:
        >>> config = PostgresVectorConfig(
        ...     host="localhost",
        ...     port=5432,
        ...     database="docflow",
        ...     user="postgres",
        ...     password="secret",
        ...     embedding_api_key="sk-..."
        ... )
    
    Recommendations:
        - index_type: HNSW for fast queries, IVFFLAT for lower memory
        - ef_construction: 64-128 for good recall
        - m: 16-64 (higher = better recall, more memory)
    """
    
    # Connection
    host: str = "localhost"
    port: int = 5432
    user: str = "postgres"
    password: str = ""
    ssl_mode: str = "prefer"  # disable, allow, prefer, require, verify-ca, verify-full
    
    # Schema
    schema: str = "public"
    table_name: str = "chunks"
    
    # pgvector specific
    vector_extension: str = "vector"
    
    # HNSW index parameters
    m: int = 16  # Max connections per layer
    ef_construction: int = 64  # Size of dynamic candidate list
    ef_search: int = 40  # Size of dynamic candidate list for search
    
    # IVFFlat index parameters
    lists: int = 100  # Number of inverted lists
    probes: int = 10  # Number of lists to probe
    
    def get_dsn(self) -> str:
        """Get PostgreSQL connection string."""
        if self.connection_string:
            return self.connection_string
        return f"postgresql://{self.user}:{self.password}@{self.host}:{self.port}/{self.database}?sslmode={self.ssl_mode}"


class PostgresVectorStore(VectorStore):
    """PostgreSQL vector store using pgvector.
    
    Example:
        >>> config = PostgresVectorConfig(
        ...     host="localhost",
        ...     database="docflow",
        ...     user="postgres",
        ...     password="secret",
        ...     embedding_api_key="sk-..."
        ... )
        >>> store = PostgresVectorStore(config)
        >>> await store.connect()
        >>> await store.create_collection()
        >>> 
        >>> # Insert chunks
        >>> ids = await store.insert(chunks)
        >>> 
        >>> # Search
        >>> results = await store.search("query text", top_k=5)
    """
    
    def __init__(self, config: PostgresVectorConfig) -> None:
        super().__init__(config)
        self.config: PostgresVectorConfig = config
        self._pool = None
    
    async def connect(self) -> None:
        """Connect to PostgreSQL."""
        try:
            import asyncpg
        except ImportError:
            raise ImportError("asyncpg is required: pip install asyncpg")
        
        self._pool = await asyncpg.create_pool(
            self.config.get_dsn(),
            min_size=1,
            max_size=self.config.pool_size,
            timeout=self.config.timeout,
        )
        
        # Enable pgvector extension
        async with self._pool.acquire() as conn:
            await conn.execute(f"CREATE EXTENSION IF NOT EXISTS {self.config.vector_extension}")
    
    async def disconnect(self) -> None:
        """Disconnect from PostgreSQL."""
        if self._pool:
            await self._pool.close()
            self._pool = None
    
    async def create_collection(self, name: Optional[str] = None) -> None:
        """Create table with vector column and index."""
        table = name or self.config.table_name
        dim = self.config.embedding_dimensions
        
        async with self._pool.acquire() as conn:
            # Create table
            await conn.execute(f"""
                CREATE TABLE IF NOT EXISTS {self.config.schema}.{table} (
                    id SERIAL PRIMARY KEY,
                    chunk_id TEXT UNIQUE,
                    content TEXT NOT NULL,
                    embedding vector({dim}),
                    metadata JSONB DEFAULT '{{}}',
                    source_file TEXT,
                    chunk_index INTEGER,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
            """)
            
            # Create vector index
            if self.config.index_type == IndexType.HNSW:
                op_class = self._get_distance_op_class()
                await conn.execute(f"""
                    CREATE INDEX IF NOT EXISTS {table}_embedding_idx
                    ON {self.config.schema}.{table}
                    USING hnsw (embedding {op_class})
                    WITH (m = {self.config.m}, ef_construction = {self.config.ef_construction})
                """)
            elif self.config.index_type == IndexType.IVFFLAT:
                op_class = self._get_distance_op_class()
                await conn.execute(f"""
                    CREATE INDEX IF NOT EXISTS {table}_embedding_idx
                    ON {self.config.schema}.{table}
                    USING ivfflat (embedding {op_class})
                    WITH (lists = {self.config.lists})
                """)
    
    async def drop_collection(self, name: Optional[str] = None) -> None:
        """Drop table."""
        table = name or self.config.table_name
        async with self._pool.acquire() as conn:
            await conn.execute(f"DROP TABLE IF EXISTS {self.config.schema}.{table}")
    
    async def insert(
        self,
        chunks: List[Chunk],
        embeddings: Optional[List[List[float]]] = None,
    ) -> List[str]:
        """Insert chunks with embeddings."""
        if embeddings is None:
            texts = [c.content for c in chunks]
            embeddings = await self.generate_embeddings(texts)
        
        table = self.config.table_name
        ids = []
        
        async with self._pool.acquire() as conn:
            for chunk, embedding in zip(chunks, embeddings):
                chunk_id = f"{chunk.metadata.section_title}_{chunk.index}" if chunk.metadata else f"chunk_{chunk.index}"
                
                metadata = {}
                if chunk.metadata:
                    metadata = {
                        "section_title": chunk.metadata.section_title,
                        "heading_path": chunk.metadata.heading_path,
                        "has_table": chunk.metadata.has_table,
                        "has_image": chunk.metadata.has_image,
                        "page": chunk.metadata.page,
                    }
                
                await conn.execute(f"""
                    INSERT INTO {self.config.schema}.{table}
                    (chunk_id, content, embedding, metadata, chunk_index)
                    VALUES ($1, $2, $3, $4, $5)
                    ON CONFLICT (chunk_id) DO UPDATE SET
                        content = EXCLUDED.content,
                        embedding = EXCLUDED.embedding,
                        metadata = EXCLUDED.metadata
                """, chunk_id, chunk.content, str(embedding), json.dumps(metadata), chunk.index)
                
                ids.append(chunk_id)
        
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
        table = self.config.table_name
        distance_op = self._get_distance_operator()
        
        # Build filter clause
        filter_clause = ""
        if filter_metadata:
            conditions = []
            for key, value in filter_metadata.items():
                conditions.append(f"metadata->>'{key}' = '{value}'")
            filter_clause = "WHERE " + " AND ".join(conditions)
        
        # Set search parameters for HNSW
        async with self._pool.acquire() as conn:
            if self.config.index_type == IndexType.HNSW:
                await conn.execute(f"SET hnsw.ef_search = {self.config.ef_search}")
            elif self.config.index_type == IndexType.IVFFLAT:
                await conn.execute(f"SET ivfflat.probes = {self.config.probes}")
            
            rows = await conn.fetch(f"""
                SELECT chunk_id, content, embedding, metadata, chunk_index,
                       embedding {distance_op} $1 AS distance
                FROM {self.config.schema}.{table}
                {filter_clause}
                ORDER BY embedding {distance_op} $1
                LIMIT $2
            """, str(embedding), top_k)
        
        results = []
        for row in rows:
            metadata = json.loads(row["metadata"]) if row["metadata"] else {}
            chunk_meta = ChunkMetadata(
                section_title=metadata.get("section_title", ""),
                heading_path=metadata.get("heading_path", []),
                has_table=metadata.get("has_table", False),
                has_image=metadata.get("has_image", False),
                page=metadata.get("page", 0),
            )
            
            chunk = Chunk(
                content=row["content"],
                index=row["chunk_index"],
                start_char=0,
                end_char=len(row["content"]),
                metadata=chunk_meta,
            )
            
            # Convert distance to similarity score
            distance = float(row["distance"])
            if self.config.distance_metric == DistanceMetric.COSINE:
                score = 1 - distance
            else:
                score = 1 / (1 + distance)
            
            results.append(SearchResult(chunk=chunk, score=score, metadata=metadata))
        
        return results
    
    async def delete(
        self,
        ids: Optional[List[str]] = None,
        filter_metadata: Optional[Dict[str, Any]] = None,
    ) -> int:
        """Delete chunks."""
        table = self.config.table_name
        
        async with self._pool.acquire() as conn:
            if ids:
                result = await conn.execute(f"""
                    DELETE FROM {self.config.schema}.{table}
                    WHERE chunk_id = ANY($1)
                """, ids)
            elif filter_metadata:
                conditions = []
                for key, value in filter_metadata.items():
                    conditions.append(f"metadata->>'{key}' = '{value}'")
                where_clause = " AND ".join(conditions)
                result = await conn.execute(f"""
                    DELETE FROM {self.config.schema}.{table}
                    WHERE {where_clause}
                """)
            else:
                return 0
            
            return int(result.split()[-1])
    
    async def update(
        self,
        id: str,
        chunk: Optional[Chunk] = None,
        embedding: Optional[List[float]] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> bool:
        """Update a chunk."""
        table = self.config.table_name
        updates = []
        values = [id]
        idx = 2
        
        if chunk:
            updates.append(f"content = ${idx}")
            values.append(chunk.content)
            idx += 1
        
        if embedding:
            updates.append(f"embedding = ${idx}")
            values.append(str(embedding))
            idx += 1
        
        if metadata:
            updates.append(f"metadata = ${idx}")
            values.append(json.dumps(metadata))
            idx += 1
        
        if not updates:
            return False
        
        async with self._pool.acquire() as conn:
            result = await conn.execute(f"""
                UPDATE {self.config.schema}.{table}
                SET {', '.join(updates)}
                WHERE chunk_id = $1
            """, *values)
            
            return "UPDATE 1" in result
    
    async def count(self, filter_metadata: Optional[Dict[str, Any]] = None) -> int:
        """Count chunks."""
        table = self.config.table_name
        
        filter_clause = ""
        if filter_metadata:
            conditions = []
            for key, value in filter_metadata.items():
                conditions.append(f"metadata->>'{key}' = '{value}'")
            filter_clause = "WHERE " + " AND ".join(conditions)
        
        async with self._pool.acquire() as conn:
            result = await conn.fetchval(f"""
                SELECT COUNT(*) FROM {self.config.schema}.{table}
                {filter_clause}
            """)
            return result
    
    def _get_distance_operator(self) -> str:
        """Get PostgreSQL distance operator for the metric."""
        operators = {
            DistanceMetric.COSINE: "<=>",
            DistanceMetric.EUCLIDEAN: "<->",
            DistanceMetric.DOT_PRODUCT: "<#>",
        }
        return operators.get(self.config.distance_metric, "<=>")
    
    def _get_distance_op_class(self) -> str:
        """Get operator class for index creation."""
        op_classes = {
            DistanceMetric.COSINE: "vector_cosine_ops",
            DistanceMetric.EUCLIDEAN: "vector_l2_ops",
            DistanceMetric.DOT_PRODUCT: "vector_ip_ops",
        }
        return op_classes.get(self.config.distance_metric, "vector_cosine_ops")
