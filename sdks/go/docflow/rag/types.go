package rag

import (
	"github.com/xgaslan/docflow/sdks/go/docflow/config"
)

// Chunk represents a chunk of content for RAG.
type Chunk struct {
	Content   string
	Index     int
	StartChar int
	EndChar   int
	Metadata  ChunkMetadata
	Embedding []float64 // Vector embedding
}

// ChunkMetadata contains metadata for a chunk.
type ChunkMetadata struct {
	SectionTitle    string
	HeadingPath     []string
	HeadingLevels   []int
	HasTable        bool
	HasImage        bool
	HasCode         bool
	Page            int
	SectionIndex    int
	SubsectionIndex int
	ContentType     string // text, table, code, image
}

// ExtractedImage represents an extracted image.
type ExtractedImage struct {
	Data            []byte
	Format          string // png, jpg, etc.
	Filename        string
	Caption         string
	Page            int
	Position        [2]float64 // x, y
	SurroundingText string
	Description     string                 // LLM-generated
	LLMAnalysis     map[string]interface{} // Extended LLM analysis
}

// ExtractedTable represents an extracted table.
type ExtractedTable struct {
	Rows        [][]string
	Header      []string
	Caption     string
	Page        int
	Summary     string                 // LLM-generated summary
	LLMAnalysis map[string]interface{} // Extended LLM analysis
}

// RAGDocument represents a fully processed RAG document.
type RAGDocument struct {
	ID           string
	Filename     string
	Status       string
	Content      string
	Chunks       []Chunk
	Images       []ExtractedImage
	Tables       []ExtractedTable
	Metadata     config.MetadataConfig // Or DocumentMetadata?
	RawMetadata  map[string]interface{}
	SourceFile   string
	SourceFormat string
	PDFBytes     []byte      // If output format is PDF
	CreatedAt    interface{} // time.Time
}

// DocumentMetadata contains comprehensive document metadata.
type DocumentMetadata struct {
	Title        string
	Author       string
	CreatedDate  string
	ModifiedDate string
	Headings     []HeadingInfo
	WordCount    int
	CharCount    int
	PageCount    int
	ImageCount   int
	TableCount   int
	Language     string
	Keywords     []string
	Entities     []string
	Summary      string
	KeyPoints    []string
}

// HeadingInfo represents information about a document heading.
type HeadingInfo struct {
	Text     string
	Level    int
	StartPos int
	EndPos   int
}
