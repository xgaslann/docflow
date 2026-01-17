package docflow

import "time"

// Options contains configuration options for the DocFlow library.
type Options struct {
	// TempDir is the directory for temporary files.
	TempDir string

	// Timeout is the maximum time for PDF generation.
	Timeout time.Duration

	// ChromePath is the path to Chrome/Chromium executable.
	// If empty, will try to find it automatically.
	ChromePath string

	// UseBrowser determines if browser-based PDF generation should be used.
	// If false, will try alternative methods.
	UseBrowser bool
}

// DefaultOptions returns the default configuration options.
func DefaultOptions() Options {
	return Options{
		TempDir:    "/tmp/docflow",
		Timeout:    60 * time.Second,
		ChromePath: "",
		UseBrowser: true,
	}
}

// ConvertOptions contains options for conversion operations.
type ConvertOptions struct {
	// MergeMode specifies how multiple files should be handled.
	// "merged" - combine all files into one PDF
	// "separate" - create separate PDFs for each file
	MergeMode string

	// OutputName is the name for the output file (used in merged mode).
	OutputName string

	// OutputPath is the path where the output should be saved.
	// If using storage, this is relative to the storage root.
	OutputPath string
}

// Result represents the result of a conversion operation.
type Result struct {
	// Success indicates if the operation was successful.
	Success bool

	// FilePaths contains the paths to generated files.
	FilePaths []string

	// Bytes contains the raw PDF data (if requested).
	Bytes []byte

	// Error contains any error that occurred.
	Error error
}

// ExtractOptions contains options for PDF extraction.
type ExtractOptions struct {
	// PreserveLayout attempts to preserve the original PDF layout.
	PreserveLayout bool

	// FirstPageOnly only extracts the first page for preview.
	FirstPageOnly bool

	// OutputPath is where to save the extracted markdown.
	OutputPath string
}

// ExtractResult represents the result of a PDF extraction.
type ExtractResult struct {
	// Success indicates if the extraction was successful.
	Success bool

	// Markdown contains the extracted markdown content.
	Markdown string

	// FilePath is the path to the saved markdown file.
	FilePath string

	// PageCount is the number of pages in the PDF.
	PageCount int

	// Error contains any error that occurred.
	Error error
}
