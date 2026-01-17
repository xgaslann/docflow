package model

// FileData represents a single markdown file
type FileData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Order   int    `json:"order"`
}

// PreviewRequest is the request body for preview endpoint
type PreviewRequest struct {
	Content string `json:"content" validate:"required"`
}

// PreviewResponse is the response for preview endpoint
type PreviewResponse struct {
	HTML string `json:"html"`
}

// ConvertRequest is the request body for convert endpoint
type ConvertRequest struct {
	Files      []FileData `json:"files" validate:"required,min=1"`
	MergeMode  MergeMode  `json:"mergeMode" validate:"required,oneof=separate merged"`
	OutputName string     `json:"outputName,omitempty"`
}

// ConvertResponse is the response for convert endpoint
type ConvertResponse struct {
	Success bool     `json:"success"`
	Files   []string `json:"files,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// MergePreviewRequest is the request for merge preview endpoint
type MergePreviewRequest struct {
	Files []FileData `json:"files" validate:"required,min=1"`
}

// MergePreviewResponse is the response for merge preview endpoint
type MergePreviewResponse struct {
	HTML           string `json:"html"`
	TotalFiles     int    `json:"totalFiles"`
	EstimatedPages int    `json:"estimatedPages"`
}

// MergeMode represents how multiple files should be handled
type MergeMode string

const (
	MergeModeSeparate MergeMode = "separate"
	MergeModeMerged   MergeMode = "merged"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

// PDFExtractRequest is the request for PDF extraction
type PDFExtractRequest struct {
	FileName string `json:"fileName"`
	Content  string `json:"content"` // Base64 encoded PDF
}

// PDFExtractResponse is the response for PDF extraction
type PDFExtractResponse struct {
	Success  bool   `json:"success"`
	Markdown string `json:"markdown"`
	FilePath string `json:"filePath,omitempty"`
	FileName string `json:"fileName,omitempty"`
	Error    string `json:"error,omitempty"`
}

// PDFPreviewResponse is the response for PDF preview
type PDFPreviewResponse struct {
	Preview   string `json:"preview"`
	PageCount int    `json:"pageCount"`
	FileName  string `json:"fileName"`
}
