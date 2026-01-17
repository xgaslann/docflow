package config

// LLMProcessingMode defines what content to process with LLM
type LLMProcessingMode string

const (
	LLMProcessingModeImages LLMProcessingMode = "images"
	LLMProcessingModeTables LLMProcessingMode = "tables"
	LLMProcessingModeText   LLMProcessingMode = "text"
	LLMProcessingModeAll    LLMProcessingMode = "all"
)

// ChunkingStrategy defines how to chunk documents
type ChunkingStrategy string

const (
	ChunkingStrategySimple        ChunkingStrategy = "simple"
	ChunkingStrategyHeadingAware  ChunkingStrategy = "heading_aware"
	ChunkingStrategyDocumentIntel ChunkingStrategy = "doc_intel"
	ChunkingStrategySemantic      ChunkingStrategy = "semantic"
)

// OutputFormat defines the output format
type OutputFormat string

const (
	OutputFormatMarkdown OutputFormat = "markdown"
	OutputFormatPDF      OutputFormat = "pdf"
	OutputFormatHTML     OutputFormat = "html"
)

// RAGConfig configures RAG mode extraction.
type RAGConfig struct {
	Enabled bool `json:"enabled"`

	// Output
	OutputFormat OutputFormat `json:"output_format"`

	// Chunking
	ChunkSize        int              `json:"chunk_size"`
	ChunkOverlap     int              `json:"chunk_overlap"`
	ChunkingStrategy ChunkingStrategy `json:"chunking_strategy"`
	DocIntelConfig   *DocIntelConfig  `json:"doc_intel_config,omitempty"`

	// Extraction options
	ExtractImages    bool `json:"extract_images"`
	ExtractTables    bool `json:"extract_tables"`
	PreserveMetadata bool `json:"preserve_metadata"`
	ExtractHeadings  bool `json:"extract_headings"`
	GenerateTOC      bool `json:"generate_toc"`

	// LLM Processing
	LLMProcessing []LLMProcessingMode `json:"llm_processing,omitempty"`
	LLMConfig     LLMConfig           `json:"llm_config,omitempty"`

	// Chunking behavior
	RespectHeadings    bool `json:"respect_headings"`
	KeepTablesTogether bool `json:"keep_tables_together"`
	AddChunkMarkers    bool `json:"add_chunk_markers"`

	// Parallel processing
	MaxWorkers int `json:"max_workers"`
}

// DefaultRAGConfig returns default RAG configuration.
func DefaultRAGConfig() RAGConfig {
	return RAGConfig{
		Enabled:            true,
		OutputFormat:       OutputFormatMarkdown,
		ChunkSize:          1000,
		ChunkOverlap:       200,
		ChunkingStrategy:   ChunkingStrategyHeadingAware,
		ExtractImages:      true,
		ExtractTables:      true,
		PreserveMetadata:   true,
		ExtractHeadings:    true,
		GenerateTOC:        true,
		RespectHeadings:    true,
		KeepTablesTogether: true,
		AddChunkMarkers:    true,
		MaxWorkers:         4,
	}
}

// BatchConfig configures batch processing.
type BatchConfig struct {
	MaxWorkers      int  `json:"max_workers"`
	FailFast        bool `json:"fail_fast"`
	ContinueOnError bool `json:"continue_on_error"`
	TimeoutPerFile  int  `json:"timeout_per_file"`
	QueueSize       int  `json:"queue_size"`
	RetryFailed     bool `json:"retry_failed"`
	MaxRetries      int  `json:"max_retries"`
}

// DefaultBatchConfig returns default batch configuration.
func DefaultBatchConfig() BatchConfig {
	return BatchConfig{
		MaxWorkers:      4,
		FailFast:        false,
		ContinueOnError: true,
		TimeoutPerFile:  300,
		QueueSize:       100,
		RetryFailed:     true,
		MaxRetries:      3,
	}
}

// JobStatus defines the status of a batch job
type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)
