package config

import (
	"errors"
	"time"
)

// SplitBy defines text splitting strategy.
type SplitBy string

const (
	SplitByParagraph SplitBy = "paragraph"
	SplitBySentence  SplitBy = "sentence"
	SplitByToken     SplitBy = "token"
	SplitByCharacter SplitBy = "character"
	SplitByHeading   SplitBy = "heading"
)

// ChunkingConfig configures text chunking.
type ChunkingConfig struct {
	// Size settings
	ChunkSize    int `json:"chunk_size"`
	ChunkOverlap int `json:"chunk_overlap"`
	MinChunkSize int `json:"min_chunk_size"`
	MaxChunkSize int `json:"max_chunk_size"`

	// Splitting strategy
	SplitBy    SplitBy  `json:"split_by"`
	Separators []string `json:"separators"`

	// Token settings
	Tokenizer string `json:"tokenizer"`

	// Heading-aware settings
	RespectHeadings    bool `json:"respect_headings"`
	KeepTablesTogether bool `json:"keep_tables_together"`
	KeepCodeTogether   bool `json:"keep_code_together"`

	// Markers
	AddChunkMarkers bool   `json:"add_chunk_markers"`
	MarkerFormat    string `json:"marker_format"`
}

// DefaultChunkingConfig returns default chunking configuration.
func DefaultChunkingConfig() ChunkingConfig {
	return ChunkingConfig{
		ChunkSize:          1000,
		ChunkOverlap:       200,
		MinChunkSize:       100,
		MaxChunkSize:       2000,
		SplitBy:            SplitByParagraph,
		Separators:         []string{"\n\n", "\n", ". ", " "},
		Tokenizer:          "cl100k_base",
		RespectHeadings:    true,
		KeepTablesTogether: true,
		KeepCodeTogether:   true,
		AddChunkMarkers:    true,
		MarkerFormat:       "[CHUNK %d]",
	}
}

// Validate validates the configuration.
func (c ChunkingConfig) Validate() error {
	if c.ChunkSize <= 0 {
		return errors.New("chunk_size must be positive")
	}
	if c.ChunkOverlap < 0 {
		return errors.New("chunk_overlap cannot be negative")
	}
	if c.ChunkOverlap >= c.ChunkSize {
		return errors.New("chunk_overlap must be less than chunk_size")
	}
	return nil
}

// RetrievalConfig configures retrieval operations.
type RetrievalConfig struct {
	// Basic retrieval
	TopK                int     `json:"top_k"`
	SimilarityThreshold float64 `json:"similarity_threshold"`
	MinScore            float64 `json:"min_score"`

	// Reranking
	Rerank      bool   `json:"rerank"`
	RerankModel string `json:"rerank_model"`
	RerankTopK  int    `json:"rerank_top_k"`

	// Filtering
	FilterDuplicates   bool    `json:"filter_duplicates"`
	DuplicateThreshold float64 `json:"duplicate_threshold"`

	// Context
	IncludeContext bool `json:"include_context"`
	ContextBefore  int  `json:"context_before"`
	ContextAfter   int  `json:"context_after"`

	// Hybrid search
	HybridSearch   bool    `json:"hybrid_search"`
	KeywordWeight  float64 `json:"keyword_weight"`
	SemanticWeight float64 `json:"semantic_weight"`

	// MMR
	UseMMR    bool    `json:"use_mmr"`
	MMRLambda float64 `json:"mmr_lambda"`
}

// DefaultRetrievalConfig returns default retrieval configuration.
func DefaultRetrievalConfig() RetrievalConfig {
	return RetrievalConfig{
		TopK:                5,
		SimilarityThreshold: 0.7,
		MinScore:            0.0,
		Rerank:              false,
		RerankModel:         "cross-encoder/ms-marco-MiniLM-L-6-v2",
		RerankTopK:          3,
		FilterDuplicates:    true,
		DuplicateThreshold:  0.95,
		IncludeContext:      true,
		ContextBefore:       1,
		ContextAfter:        1,
		HybridSearch:        false,
		KeywordWeight:       0.3,
		SemanticWeight:      0.7,
		UseMMR:              false,
		MMRLambda:           0.5,
	}
}

// MetadataConfig configures metadata extraction.
type MetadataConfig struct {
	IncludeFields []string               `json:"include_fields"`
	ExcludeFields []string               `json:"exclude_fields"`
	CustomFields  map[string]interface{} `json:"custom_fields"`

	// Extraction toggles
	ExtractTitle     bool `json:"extract_title"`
	ExtractAuthor    bool `json:"extract_author"`
	ExtractHeadings  bool `json:"extract_headings"`
	ExtractTOC       bool `json:"extract_toc"`
	ExtractWordCount bool `json:"extract_word_count"`
	ExtractPageCount bool `json:"extract_page_count"`
	ExtractEntities  bool `json:"extract_entities"`
	ExtractSummary   bool `json:"extract_summary"`
	ExtractKeyPoints bool `json:"extract_key_points"`
	ExtractLanguage  bool `json:"extract_language"`

	// Heading extraction
	MaxHeadingLevel  int  `json:"max_heading_level"`
	BuildHeadingTree bool `json:"build_heading_tree"`

	// TOC generation
	TOCMaxDepth           int  `json:"toc_max_depth"`
	TOCIncludePageNumbers bool `json:"toc_include_page_numbers"`
}

// DefaultMetadataConfig returns default metadata configuration.
func DefaultMetadataConfig() MetadataConfig {
	return MetadataConfig{
		IncludeFields:         []string{"title", "headings", "table_of_contents", "word_count"},
		ExcludeFields:         []string{},
		CustomFields:          make(map[string]interface{}),
		ExtractTitle:          true,
		ExtractAuthor:         true,
		ExtractHeadings:       true,
		ExtractTOC:            true,
		ExtractWordCount:      true,
		ExtractPageCount:      true,
		ExtractEntities:       false,
		ExtractSummary:        false,
		ExtractKeyPoints:      false,
		ExtractLanguage:       false,
		MaxHeadingLevel:       6,
		BuildHeadingTree:      true,
		TOCMaxDepth:           3,
		TOCIncludePageNumbers: true,
	}
}

// ShouldExtract checks if a field should be extracted.
func (c MetadataConfig) ShouldExtract(fieldName string) bool {
	for _, f := range c.ExcludeFields {
		if f == fieldName {
			return false
		}
	}
	if len(c.IncludeFields) > 0 {
		for _, f := range c.IncludeFields {
			if f == fieldName {
				return true
			}
		}
		return false
	}
	return true
}

// DocIntelProvider defines Document Intelligence provider.
type DocIntelProvider string

const (
	DocIntelProviderAzure DocIntelProvider = "azure"
	DocIntelProviderAWS   DocIntelProvider = "aws"
)

// DocIntelConfig configures Document Intelligence services.
type DocIntelConfig struct {
	Provider DocIntelProvider `json:"provider"`

	// Azure settings
	Endpoint     string   `json:"endpoint"`
	APIKey       string   `json:"api_key"`
	ModelID      string   `json:"model_id"`
	APIVersion   string   `json:"api_version"`
	Locale       string   `json:"locale"`
	Language     string   `json:"language"`
	Features     []string `json:"features"`
	OutputFormat string   `json:"output_format"`
	Pages        string   `json:"pages"`

	// AWS settings
	AWSRegion        string   `json:"aws_region"`
	AWSAccessKey     string   `json:"aws_access_key"`
	AWSSecretKey     string   `json:"aws_secret_key"`
	TextractFeatures []string `json:"textract_features"`

	// Common
	Timeout       time.Duration `json:"timeout"`
	MinConfidence float64       `json:"min_confidence"`
}

// DefaultDocIntelConfig returns default Document Intelligence configuration.
func DefaultDocIntelConfig() DocIntelConfig {
	return DocIntelConfig{
		Provider:         DocIntelProviderAzure,
		ModelID:          "prebuilt-layout",
		APIVersion:       "2024-02-29-preview",
		Locale:           "en-US",
		Features:         []string{"keyValuePairs", "languages"},
		OutputFormat:     "markdown",
		AWSRegion:        "us-east-1",
		TextractFeatures: []string{"TABLES", "FORMS"},
		Timeout:          5 * time.Minute,
		MinConfidence:    0.0,
	}
}
