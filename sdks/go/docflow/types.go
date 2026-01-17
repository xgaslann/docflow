package docflow

import (
	"github.com/xgaslan/docflow/sdks/go/docflow/config"
	"github.com/xgaslan/docflow/sdks/go/docflow/rag"
)

// Basic Types

// MDFile represents a Markdown file to be converted.
type MDFile struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Order   int    `json:"order"`
}

// NewMDFile creates a new MDFile with the given name and content.
func NewMDFile(name, content string) MDFile {
	return MDFile{
		ID:      name,
		Name:    name,
		Content: content,
		Order:   0,
	}
}

// NewMDFileWithOrder creates a new MDFile with the given name, content, and order.
func NewMDFileWithOrder(name, content string, order int) MDFile {
	return MDFile{
		ID:      name,
		Name:    name,
		Content: content,
		Order:   order,
	}
}

// ConvertResult represents the result of a format conversion.
type ConvertResult struct {
	Success  bool                   `json:"success"`
	Content  string                 `json:"content,omitempty"`
	Format   string                 `json:"format,omitempty"`
	Error    error                  `json:"error,omitempty"`
	Images   []rag.ExtractedImage   `json:"images,omitempty"`
	Tables   []rag.ExtractedTable   `json:"tables,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// BatchJob represents a batch processing job.
type BatchJob struct {
	JobID          string            `json:"job_id"`
	Status         config.JobStatus  `json:"status"`
	TotalFiles     int               `json:"total_files"`
	ProcessedFiles int               `json:"processed_files"`
	FailedFiles    int               `json:"failed_files"`
	Results        []rag.RAGDocument `json:"results,omitempty"`
	Errors         map[string]string `json:"errors,omitempty"`
	CreatedAt      interface{}       `json:"created_at,omitempty"`   // time.Time
	CompletedAt    interface{}       `json:"completed_at,omitempty"` // time.Time
}
