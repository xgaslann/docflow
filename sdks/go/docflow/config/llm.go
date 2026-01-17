package config

// LLMProvider defines LLM provider.
type LLMProvider string

const (
	LLMProviderOpenAI      LLMProvider = "openai"
	LLMProviderAzureOpenAI LLMProvider = "azure_openai"
	LLMProviderAnthropic   LLMProvider = "anthropic"
	LLMProviderOllama      LLMProvider = "ollama"
	LLMProviderGoogle      LLMProvider = "google"
)

// LLMPrompts contains custom prompts for LLM processing.
type LLMPrompts struct {
	ImageDescription       string            `json:"image_description"`
	ImageOCR               string            `json:"image_ocr"`
	TableAnalysis          string            `json:"table_analysis"`
	TableToText            string            `json:"table_to_text"`
	TextSummary            string            `json:"text_summary"`
	TextKeyPoints          string            `json:"text_key_points"`
	EntityExtraction       string            `json:"entity_extraction"`
	KeywordExtraction      string            `json:"keyword_extraction"`
	DocumentClassification string            `json:"document_classification"`
	DocumentQA             string            `json:"document_qa"`
	Custom                 map[string]string `json:"custom"`
}

// DefaultLLMPrompts returns default prompts.
func DefaultLLMPrompts() LLMPrompts {
	return LLMPrompts{
		ImageDescription: `Describe this image in detail for use in a document retrieval system.
Focus on:
1. Key information visible in the image
2. Text, numbers, or data shown
3. Context and relevance to the document`,

		ImageOCR: "Extract all text visible in this image. Preserve formatting and structure.",

		TableAnalysis: `Analyze this table and provide:
1. Brief summary of what the table contains
2. Key insights or patterns
3. Important data points

Table:
%s`,

		TextSummary: `Summarize the following content concisely.
Focus on main points and key information.
Maximum %d characters.

Content:
%s`,

		TextKeyPoints: `Extract the %d most important key points from:
%s

Return as a JSON array of strings.`,

		EntityExtraction: `Extract all named entities from this content.
Include: people, organizations, locations, products, dates, monetary values.
Return as a JSON object.

Content:
%s`,

		Custom: make(map[string]string),
	}
}

// GetPrompt returns a prompt by name.
func (p LLMPrompts) GetPrompt(name string) string {
	if custom, ok := p.Custom[name]; ok {
		return custom
	}
	return ""
}

// SetCustomPrompt sets a custom prompt.
func (p *LLMPrompts) SetCustomPrompt(name, prompt string) {
	if p.Custom == nil {
		p.Custom = make(map[string]string)
	}
	p.Custom[name] = prompt
}

// LLMConfig configures LLM integration.
type LLMConfig struct {
	Provider LLMProvider `json:"provider"`
	Model    string      `json:"model"`
	APIKey   string      `json:"api_key"`

	// Prompts
	Prompts LLMPrompts `json:"prompts"`

	// OpenAI
	Organization string `json:"organization"`
	BaseURL      string `json:"base_url"`

	// Azure OpenAI
	AzureEndpoint   string `json:"azure_endpoint"`
	AzureDeployment string `json:"azure_deployment"`
	APIVersion      string `json:"api_version"`

	// Ollama
	OllamaBaseURL string `json:"ollama_base_url"`

	// Generation Parameters
	Temperature      float64  `json:"temperature"`
	MaxTokens        int      `json:"max_tokens"`
	TopP             float64  `json:"top_p"`
	FrequencyPenalty float64  `json:"frequency_penalty"`
	PresencePenalty  float64  `json:"presence_penalty"`
	StopSequences    []string `json:"stop_sequences"`

	// Vision Parameters
	Detail           string   `json:"detail"`
	MaxImageSize     int64    `json:"max_image_size"`
	SupportedFormats []string `json:"supported_formats"`

	// Retry & Timeout
	Timeout         int     `json:"timeout"`
	RetryCount      int     `json:"retry_count"`
	RetryDelay      float64 `json:"retry_delay"`
	RetryMultiplier float64 `json:"retry_multiplier"`

	// Batch Processing
	BatchSize          int `json:"batch_size"`
	ConcurrentRequests int `json:"concurrent_requests"`

	// Response
	ResponseFormat string `json:"response_format"`
}

// DefaultLLMConfig returns default LLM configuration.
func DefaultLLMConfig() LLMConfig {
	return LLMConfig{
		Provider:           LLMProviderOpenAI,
		Model:              "gpt-4-vision-preview",
		Prompts:            DefaultLLMPrompts(),
		OllamaBaseURL:      "http://localhost:11434",
		Temperature:        0.7,
		MaxTokens:          1000,
		TopP:               1.0,
		FrequencyPenalty:   0.0,
		PresencePenalty:    0.0,
		Detail:             "auto",
		MaxImageSize:       20 * 1024 * 1024,
		SupportedFormats:   []string{"png", "jpg", "jpeg", "gif", "webp"},
		Timeout:            60,
		RetryCount:         3,
		RetryDelay:         1.0,
		RetryMultiplier:    2.0,
		BatchSize:          5,
		ConcurrentRequests: 3,
		ResponseFormat:     "text",
	}
}

// AISearchConfig configures Azure AI Search.
type AISearchConfig struct {
	// Connection
	Endpoint   string `json:"endpoint"`
	APIKey     string `json:"api_key"`
	APIVersion string `json:"api_version"`
	IndexName  string `json:"index_name"`

	// Vector configuration
	VectorSearchProfile string `json:"vector_search_profile"`
	HNSWM               int    `json:"hnsw_m"`
	HNSWEFConstruction  int    `json:"hnsw_ef_construction"`
	HNSWEFSearch        int    `json:"hnsw_ef_search"`
	Metric              string `json:"metric"`

	// Semantic configuration
	SemanticConfig            string   `json:"semantic_config"`
	SemanticPrioritizedFields []string `json:"semantic_prioritized_fields"`

	// Search options
	QueryType  string `json:"query_type"`  // simple, full, semantic
	SearchMode string `json:"search_mode"` // any, all
	Top        int    `json:"top"`
	Skip       int    `json:"skip"`

	// Vector search
	VectorFields      []string `json:"vector_fields"`
	KNearestNeighbors int      `json:"k_nearest_neighbors"`

	// Hybrid search
	HybridSearch      bool `json:"hybrid_search"`
	SemanticReranking bool `json:"semantic_reranking"`

	// Embedding
	EmbeddingModel      string `json:"embedding_model"`
	EmbeddingDimensions int    `json:"embedding_dimensions"`
	EmbeddingAPIKey     string `json:"embedding_api_key"`
}

// DefaultAISearchConfig returns default AI Search configuration.
func DefaultAISearchConfig() AISearchConfig {
	return AISearchConfig{
		APIVersion:                "2024-07-01",
		IndexName:                 "docflow-index",
		VectorSearchProfile:       "default-vector-profile",
		HNSWM:                     4,
		HNSWEFConstruction:        400,
		HNSWEFSearch:              500,
		Metric:                    "cosine",
		SemanticConfig:            "default-semantic-config",
		SemanticPrioritizedFields: []string{"content", "title"},
		QueryType:                 "semantic",
		SearchMode:                "any",
		Top:                       10,
		Skip:                      0,
		VectorFields:              []string{"content_vector"},
		KNearestNeighbors:         50,
		HybridSearch:              true,
		SemanticReranking:         true,
		EmbeddingModel:            "text-embedding-3-small",
		EmbeddingDimensions:       1536,
	}
}

// VectorStoreProvider defines vector store provider.
type VectorStoreProvider string

const (
	VectorStoreProviderPostgres VectorStoreProvider = "postgresql"
	VectorStoreProviderMongoDB  VectorStoreProvider = "mongodb"
)

// VectorStoreConfig configures vector storage.
type VectorStoreConfig struct {
	Provider         VectorStoreProvider `json:"provider"`
	ConnectionString string              `json:"connection_string"`
	Database         string              `json:"database"`
	Collection       string              `json:"collection"`

	// Embedding
	EmbeddingProvider   string `json:"embedding_provider"`
	EmbeddingModel      string `json:"embedding_model"`
	EmbeddingAPIKey     string `json:"embedding_api_key"`
	EmbeddingDimensions int    `json:"embedding_dimensions"`

	// Index
	IndexType      string `json:"index_type"`      // hnsw, ivfflat
	DistanceMetric string `json:"distance_metric"` // cosine, euclidean, dot

	// PostgreSQL specific
	Host           string `json:"host"`
	Port           int    `json:"port"`
	User           string `json:"user"`
	Password       string `json:"password"`
	SSLMode        string `json:"ssl_mode"`
	Schema         string `json:"schema"`
	M              int    `json:"m"`
	EFConstruction int    `json:"ef_construction"`
	EFSearch       int    `json:"ef_search"`

	// MongoDB specific
	AtlasCluster   string `json:"atlas_cluster"`
	NumCandidates  int    `json:"num_candidates"`
	IndexNameMongo string `json:"index_name"`
}

// DefaultVectorStoreConfig returns default vector store configuration.
func DefaultVectorStoreConfig() VectorStoreConfig {
	return VectorStoreConfig{
		Provider:            VectorStoreProviderPostgres,
		Database:            "docflow",
		Collection:          "chunks",
		EmbeddingProvider:   "openai",
		EmbeddingModel:      "text-embedding-3-small",
		EmbeddingDimensions: 1536,
		IndexType:           "hnsw",
		DistanceMetric:      "cosine",
		Host:                "localhost",
		Port:                5432,
		User:                "postgres",
		SSLMode:             "prefer",
		Schema:              "public",
		M:                   16,
		EFConstruction:      64,
		EFSearch:            40,
		NumCandidates:       100,
		IndexNameMongo:      "vector_index",
	}
}
