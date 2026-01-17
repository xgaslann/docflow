package vector

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresVectorStore implements vector storage using PostgreSQL with pgvector.
type PostgresVectorStore struct {
	db        *sql.DB
	config    PostgresConfig
	tableName string
}

// PostgresConfig holds PostgreSQL connection settings.
type PostgresConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	Database   string
	SSLMode    string
	Schema     string
	TableName  string
	Dimensions int
}

// DefaultPostgresConfig returns sensible defaults.
func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:       "localhost",
		Port:       5432,
		User:       "postgres",
		Password:   "",
		Database:   "docflow",
		SSLMode:    "disable",
		Schema:     "public",
		TableName:  "chunks",
		Dimensions: 1536,
	}
}

// NewPostgresVectorStore creates a new PostgreSQL vector store.
func NewPostgresVectorStore(config PostgresConfig) (*PostgresVectorStore, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	store := &PostgresVectorStore{
		db:        db,
		config:    config,
		tableName: fmt.Sprintf("%s.%s", config.Schema, config.TableName),
	}

	if err := store.ensureTable(); err != nil {
		return nil, err
	}

	return store, nil
}

// NewPostgresVectorStoreFromDSN creates a store from a connection string.
func NewPostgresVectorStoreFromDSN(dsn, tableName string, dimensions int) (*PostgresVectorStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	store := &PostgresVectorStore{
		db:        db,
		tableName: tableName,
		config:    PostgresConfig{Dimensions: dimensions},
	}

	if err := store.ensureTable(); err != nil {
		return nil, err
	}

	return store, nil
}

// VectorRecord represents a single vector record.
type VectorRecord struct {
	ID         string
	DocumentID string
	ChunkIndex int
	Content    string
	Embedding  []float32
	Metadata   map[string]interface{}
}

// SearchResult represents a search result.
type SearchResult struct {
	ID       string
	Content  string
	Score    float32
	Metadata map[string]interface{}
}

// Upsert inserts or updates a vector record.
func (s *PostgresVectorStore) Upsert(record VectorRecord) error {
	metadataJSON, _ := json.Marshal(record.Metadata)
	embeddingStr := floatsToVector(record.Embedding)

	query := fmt.Sprintf(`
		INSERT INTO %s (id, document_id, chunk_index, content, embedding, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5::vector, $6::jsonb, NOW())
		ON CONFLICT (id) DO UPDATE SET
			content = EXCLUDED.content,
			embedding = EXCLUDED.embedding,
			metadata = EXCLUDED.metadata,
			updated_at = NOW()
	`, s.tableName)

	_, err := s.db.Exec(query, record.ID, record.DocumentID, record.ChunkIndex, record.Content, embeddingStr, string(metadataJSON))
	return err
}

// UpsertBatch inserts multiple records in a transaction.
func (s *PostgresVectorStore) UpsertBatch(records []VectorRecord) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`
		INSERT INTO %s (id, document_id, chunk_index, content, embedding, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5::vector, $6::jsonb, NOW())
		ON CONFLICT (id) DO UPDATE SET
			content = EXCLUDED.content,
			embedding = EXCLUDED.embedding,
			metadata = EXCLUDED.metadata,
			updated_at = NOW()
	`, s.tableName)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, record := range records {
		metadataJSON, _ := json.Marshal(record.Metadata)
		embeddingStr := floatsToVector(record.Embedding)

		_, err := stmt.Exec(record.ID, record.DocumentID, record.ChunkIndex, record.Content, embeddingStr, string(metadataJSON))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Search performs a similarity search.
func (s *PostgresVectorStore) Search(queryVector []float32, topK int) ([]SearchResult, error) {
	return s.SearchWithFilter(queryVector, topK, nil)
}

// SearchWithFilter performs a similarity search with metadata filter.
func (s *PostgresVectorStore) SearchWithFilter(queryVector []float32, topK int, filter map[string]interface{}) ([]SearchResult, error) {
	vectorStr := floatsToVector(queryVector)

	query := fmt.Sprintf(`
		SELECT id, content, metadata, 1 - (embedding <=> $1::vector) as score
		FROM %s
	`, s.tableName)

	args := []interface{}{vectorStr}

	if len(filter) > 0 {
		conditions := []string{}
		argIndex := 2
		for key, value := range filter {
			conditions = append(conditions, fmt.Sprintf("metadata->>'%s' = $%d", key, argIndex))
			args = append(args, value)
			argIndex++
		}
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY embedding <=> $1::vector LIMIT %d", topK)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var r SearchResult
		var metadataJSON string

		err := rows.Scan(&r.ID, &r.Content, &metadataJSON, &r.Score)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(metadataJSON), &r.Metadata)
		results = append(results, r)
	}

	return results, nil
}

// Delete removes all chunks for a document.
func (s *PostgresVectorStore) Delete(documentID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE document_id = $1", s.tableName)
	_, err := s.db.Exec(query, documentID)
	return err
}

// Close closes the database connection.
func (s *PostgresVectorStore) Close() error {
	return s.db.Close()
}

func (s *PostgresVectorStore) ensureTable() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			document_id TEXT NOT NULL,
			chunk_index INTEGER NOT NULL,
			content TEXT NOT NULL,
			embedding vector(%d),
			metadata JSONB,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP
		)
	`, s.tableName, s.config.Dimensions)

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	// Create index
	indexQuery := fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS %s_embedding_idx ON %s
		USING hnsw (embedding vector_cosine_ops)
	`, s.config.TableName, s.tableName)

	s.db.Exec(indexQuery) // Ignore error if already exists

	return nil
}

func floatsToVector(floats []float32) string {
	if len(floats) == 0 {
		return "[]"
	}
	parts := make([]string, len(floats))
	for i, f := range floats {
		parts[i] = fmt.Sprintf("%f", f)
	}
	return "[" + strings.Join(parts, ",") + "]"
}
