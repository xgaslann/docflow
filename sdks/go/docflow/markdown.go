// Package docflow provides a standalone library for converting between
// Markdown and PDF formats without requiring a server.
package docflow

import (
	"bytes"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownParser handles markdown processing.
type MarkdownParser struct {
	md goldmark.Markdown
}

// NewMarkdownParser creates a new markdown parser with sensible defaults.
func NewMarkdownParser() *MarkdownParser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
			extension.Footnote,
			extension.DefinitionList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	return &MarkdownParser{md: md}
}

// ToHTML converts markdown content to HTML.
func (p *MarkdownParser) ToHTML(content string) (string, error) {
	var buf bytes.Buffer
	if err := p.md.Convert([]byte(content), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MergeFiles merges multiple files into a single content string.
// Files are sorted by their Order field.
func (p *MarkdownParser) MergeFiles(files []MDFile) string {
	if len(files) == 0 {
		return ""
	}

	// Sort files by order
	sortedFiles := make([]MDFile, len(files))
	copy(sortedFiles, files)
	sort.Slice(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].Order < sortedFiles[j].Order
	})

	var builder strings.Builder
	for i, file := range sortedFiles {
		if i > 0 {
			// Add page break marker and separator between files
			builder.WriteString("\n\n---\n\n")
		}
		builder.WriteString(file.Content)
	}

	return builder.String()
}

// MergeFilesToHTML merges files and converts to HTML with file separators.
func (p *MarkdownParser) MergeFilesToHTML(files []MDFile) (string, error) {
	if len(files) == 0 {
		return "", nil
	}

	// Sort files by order
	sortedFiles := make([]MDFile, len(files))
	copy(sortedFiles, files)
	sort.Slice(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].Order < sortedFiles[j].Order
	})

	var builder strings.Builder
	for i, file := range sortedFiles {
		if i > 0 {
			builder.WriteString(`<div class="file-separator"><span>`)
			builder.WriteString(file.Name)
			builder.WriteString(`</span></div>`)
		} else {
			builder.WriteString(`<div class="file-header"><span>`)
			builder.WriteString(file.Name)
			builder.WriteString(`</span></div>`)
		}

		html, err := p.ToHTML(file.Content)
		if err != nil {
			return "", err
		}
		builder.WriteString(`<div class="file-content">`)
		builder.WriteString(html)
		builder.WriteString(`</div>`)
	}

	return builder.String(), nil
}

// EstimatePageCount estimates the number of PDF pages based on content.
func (p *MarkdownParser) EstimatePageCount(content string) int {
	// Rough estimation: ~3000 characters per page
	const charsPerPage = 3000
	pages := len(content) / charsPerPage
	if pages < 1 {
		return 1
	}
	return pages
}
