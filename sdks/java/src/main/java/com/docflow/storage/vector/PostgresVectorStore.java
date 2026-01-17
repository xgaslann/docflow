package com.docflow.storage.vector;

import com.docflow.config.VectorStoreConfig;
import com.docflow.models.Chunk;
import com.docflow.models.RAGDocument;

import java.sql.*;
import java.util.*;

/**
 * PostgreSQL Vector Store using pgvector extension.
 */
public class PostgresVectorStore {

    private final VectorStoreConfig config;
    private Connection connection;

    public PostgresVectorStore(VectorStoreConfig config) {
        this.config = config;
    }

    public PostgresVectorStore(String connectionString) {
        this.config = new VectorStoreConfig();
        this.config.setConnectionString(connectionString);
    }

    /**
     * Initialize the connection and ensure table exists.
     */
    public void initialize() throws SQLException {
        connection = DriverManager.getConnection(config.getDsn());
        ensureTableExists();
    }

    /**
     * Upsert a document and its chunks.
     */
    public void upsert(RAGDocument doc) throws SQLException {
        String sql = String.format("""
                INSERT INTO %s.%s (id, document_id, chunk_index, content, embedding, metadata, created_at)
                VALUES (?, ?, ?, ?, ?::vector, ?::jsonb, NOW())
                ON CONFLICT (id) DO UPDATE SET
                    content = EXCLUDED.content,
                    embedding = EXCLUDED.embedding,
                    metadata = EXCLUDED.metadata,
                    updated_at = NOW()
                """, config.getSchema(), config.getTableName());

        try (PreparedStatement stmt = connection.prepareStatement(sql)) {
            for (Chunk chunk : doc.getChunks()) {
                String id = doc.getId() + "_" + chunk.getIndex();

                stmt.setString(1, id);
                stmt.setString(2, doc.getId());
                stmt.setInt(3, chunk.getIndex());
                stmt.setString(4, chunk.getContent());
                stmt.setString(5, null); // embedding - would need to generate
                stmt.setString(6, toJson(chunk.getMetadata()));

                stmt.addBatch();
            }
            stmt.executeBatch();
        }
    }

    /**
     * Search for similar content.
     */
    public List<SearchResult> search(float[] queryVector, int topK) throws SQLException {
        return search(queryVector, topK, null);
    }

    /**
     * Search with filter.
     */
    public List<SearchResult> search(float[] queryVector, int topK, Map<String, Object> filter) throws SQLException {
        String vectorStr = arrayToVector(queryVector);

        StringBuilder sql = new StringBuilder(String.format("""
                SELECT id, content, metadata, 1 - (embedding <=> '%s'::vector) as score
                FROM %s.%s
                """, vectorStr, config.getSchema(), config.getTableName()));

        if (filter != null && !filter.isEmpty()) {
            sql.append(" WHERE ");
            List<String> conditions = new ArrayList<>();
            for (Map.Entry<String, Object> entry : filter.entrySet()) {
                conditions.add(String.format("metadata->>'%s' = '%s'", entry.getKey(), entry.getValue()));
            }
            sql.append(String.join(" AND ", conditions));
        }

        sql.append(String.format(" ORDER BY embedding <=> '%s'::vector LIMIT %d", vectorStr, topK));

        List<SearchResult> results = new ArrayList<>();

        try (Statement stmt = connection.createStatement();
                ResultSet rs = stmt.executeQuery(sql.toString())) {
            while (rs.next()) {
                SearchResult result = new SearchResult();
                result.setId(rs.getString("id"));
                result.setContent(rs.getString("content"));
                result.setScore(rs.getFloat("score"));
                result.setMetadata(rs.getString("metadata"));
                results.add(result);
            }
        }

        return results;
    }

    /**
     * Delete document chunks.
     */
    public void delete(String documentId) throws SQLException {
        String sql = String.format("DELETE FROM %s.%s WHERE document_id = ?",
                config.getSchema(), config.getTableName());

        try (PreparedStatement stmt = connection.prepareStatement(sql)) {
            stmt.setString(1, documentId);
            stmt.executeUpdate();
        }
    }

    /**
     * Close the connection.
     */
    public void close() throws SQLException {
        if (connection != null && !connection.isClosed()) {
            connection.close();
        }
    }

    private void ensureTableExists() throws SQLException {
        String sql = String.format("""
                CREATE TABLE IF NOT EXISTS %s.%s (
                    id TEXT PRIMARY KEY,
                    document_id TEXT NOT NULL,
                    chunk_index INTEGER NOT NULL,
                    content TEXT NOT NULL,
                    embedding vector(%d),
                    metadata JSONB,
                    created_at TIMESTAMP DEFAULT NOW(),
                    updated_at TIMESTAMP
                )
                """, config.getSchema(), config.getTableName(), config.getEmbeddingDimensions());

        try (Statement stmt = connection.createStatement()) {
            stmt.execute(sql);

            // Create index
            String indexSql = String.format("""
                    CREATE INDEX IF NOT EXISTS %s_embedding_idx ON %s.%s
                    USING %s (embedding %s)
                    """,
                    config.getTableName(),
                    config.getSchema(),
                    config.getTableName(),
                    config.getIndexType().getValue().equals("hnsw") ? "hnsw" : "ivfflat",
                    config.getDistanceMetric().getValue().equals("cosine") ? "vector_cosine_ops" : "vector_l2_ops");

            stmt.execute(indexSql);
        }
    }

    private String arrayToVector(float[] array) {
        StringBuilder sb = new StringBuilder("[");
        for (int i = 0; i < array.length; i++) {
            if (i > 0)
                sb.append(",");
            sb.append(array[i]);
        }
        sb.append("]");
        return sb.toString();
    }

    private String toJson(Object obj) {
        // Simple JSON serialization
        if (obj == null)
            return "{}";
        return "{}"; // Would use a JSON library in production
    }

    /**
     * Search result.
     */
    public static class SearchResult {
        private String id;
        private String content;
        private float score;
        private String metadata;

        public String getId() {
            return id;
        }

        public void setId(String id) {
            this.id = id;
        }

        public String getContent() {
            return content;
        }

        public void setContent(String content) {
            this.content = content;
        }

        public float getScore() {
            return score;
        }

        public void setScore(float score) {
            this.score = score;
        }

        public String getMetadata() {
            return metadata;
        }

        public void setMetadata(String metadata) {
            this.metadata = metadata;
        }
    }
}
