package service

import (
	"strings"
	"testing"

	"github.com/gorkem/md-to-pdf/internal/model"
)

func TestNewMarkdownService(t *testing.T) {
	svc := NewMarkdownService()
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestToHTML(t *testing.T) {
	svc := NewMarkdownService()

	tests := []struct {
		name     string
		input    string
		contains []string
		notContains []string
	}{
		{
			name:     "h1 heading",
			input:    "# Hello World",
			contains: []string{"<h1>", "Hello World", "</h1>"},
		},
		{
			name:     "h2 heading",
			input:    "## Section Title",
			contains: []string{"<h2>", "Section Title", "</h2>"},
		},
		{
			name:     "paragraph",
			input:    "This is a paragraph.",
			contains: []string{"<p>", "This is a paragraph.", "</p>"},
		},
		{
			name:     "bold text",
			input:    "This is **bold** text",
			contains: []string{"<strong>", "bold", "</strong>"},
		},
		{
			name:     "italic text",
			input:    "This is *italic* text",
			contains: []string{"<em>", "italic", "</em>"},
		},
		{
			name:     "unordered list",
			input:    "- Item 1\n- Item 2",
			contains: []string{"<ul>", "<li>", "Item 1", "Item 2", "</li>", "</ul>"},
		},
		{
			name:     "ordered list",
			input:    "1. First\n2. Second",
			contains: []string{"<ol>", "<li>", "First", "Second", "</li>", "</ol>"},
		},
		{
			name:     "code block",
			input:    "```go\nfmt.Println(\"hello\")\n```",
			contains: []string{"<pre>", "<code>", "fmt.Println", "</code>", "</pre>"},
		},
		{
			name:     "inline code",
			input:    "Use `go run` to execute",
			contains: []string{"<code>", "go run", "</code>"},
		},
		{
			name:     "link",
			input:    "[GitHub](https://github.com)",
			contains: []string{"<a", "href=\"https://github.com\"", "GitHub", "</a>"},
		},
		{
			name:     "blockquote",
			input:    "> This is a quote",
			contains: []string{"<blockquote>", "This is a quote", "</blockquote>"},
		},
		{
			name:     "horizontal rule",
			input:    "---",
			contains: []string{"<hr"},
		},
		{
			name:     "table",
			input:    "| A | B |\n|---|---|\n| 1 | 2 |",
			contains: []string{"<table>", "<th>", "A", "B", "<td>", "1", "2", "</table>"},
		},
		{
			name:     "task list",
			input:    "- [x] Done\n- [ ] Todo",
			contains: []string{"<input", "type=\"checkbox\"", "Done", "Todo"},
		},
		{
			name:     "empty input",
			input:    "",
			contains: []string{},
		},
		{
			name:     "multiline",
			input:    "# Title\n\nParagraph 1\n\nParagraph 2",
			contains: []string{"<h1>", "Title", "<p>", "Paragraph 1", "Paragraph 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.ToHTML(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("expected result to contain %q, got: %s", s, result)
				}
			}
			
			for _, s := range tt.notContains {
				if strings.Contains(result, s) {
					t.Errorf("expected result to NOT contain %q, got: %s", s, result)
				}
			}
		})
	}
}

func TestMergeFiles(t *testing.T) {
	svc := NewMarkdownService()

	tests := []struct {
		name     string
		files    []model.FileData
		contains []string
	}{
		{
			name: "single file",
			files: []model.FileData{
				{ID: "1", Name: "doc.md", Content: "# Hello", Order: 0},
			},
			contains: []string{"# Hello"},
		},
		{
			name: "multiple files",
			files: []model.FileData{
				{ID: "1", Name: "first.md", Content: "# First", Order: 0},
				{ID: "2", Name: "second.md", Content: "# Second", Order: 1},
			},
			contains: []string{"# First", "# Second"},
		},
		{
			name: "files sorted by order",
			files: []model.FileData{
				{ID: "1", Name: "third.md", Content: "# Third", Order: 2},
				{ID: "2", Name: "first.md", Content: "# First", Order: 0},
				{ID: "3", Name: "second.md", Content: "# Second", Order: 1},
			},
			contains: []string{"# First", "# Second", "# Third"},
		},
		{
			name:     "empty files",
			files:    []model.FileData{},
			contains: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.MergeFiles(tt.files)

			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("expected result to contain %q, got: %s", s, result)
				}
			}
		})
	}
}

func TestMergeFilesOrder(t *testing.T) {
	svc := NewMarkdownService()

	files := []model.FileData{
		{ID: "1", Name: "c.md", Content: "CCC", Order: 2},
		{ID: "2", Name: "a.md", Content: "AAA", Order: 0},
		{ID: "3", Name: "b.md", Content: "BBB", Order: 1},
	}

	result := svc.MergeFiles(files)

	aIdx := strings.Index(result, "AAA")
	bIdx := strings.Index(result, "BBB")
	cIdx := strings.Index(result, "CCC")

	if aIdx == -1 || bIdx == -1 || cIdx == -1 {
		t.Fatalf("missing content in result: %s", result)
	}

	if !(aIdx < bIdx && bIdx < cIdx) {
		t.Errorf("expected order AAA < BBB < CCC, got indices: A=%d, B=%d, C=%d", aIdx, bIdx, cIdx)
	}
}

func TestMergeFilesToHTML(t *testing.T) {
	svc := NewMarkdownService()

	files := []model.FileData{
		{ID: "1", Name: "first.md", Content: "# First Doc", Order: 0},
		{ID: "2", Name: "second.md", Content: "# Second Doc", Order: 1},
	}

	result, err := svc.MergeFilesToHTML(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"<h1>", "First Doc",
		"Second Doc",
		"file-separator", // separator class
	}

	for _, s := range mustContain {
		if !strings.Contains(result, s) {
			t.Errorf("expected result to contain %q", s)
		}
	}
}

func TestEstimatePageCount(t *testing.T) {
	svc := NewMarkdownService()

	tests := []struct {
		name     string
		content  string
		minPages int
		maxPages int
	}{
		{
			name:     "empty content",
			content:  "",
			minPages: 1,
			maxPages: 1,
		},
		{
			name:     "short content",
			content:  "Hello world",
			minPages: 1,
			maxPages: 1,
		},
		{
			name:     "medium content",
			content:  strings.Repeat("Lorem ipsum dolor sit amet. ", 100),
			minPages: 1,
			maxPages: 2,
		},
		{
			name:     "long content",
			content:  strings.Repeat("Lorem ipsum dolor sit amet consectetur. ", 500),
			minPages: 2,
			maxPages: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.EstimatePageCount(tt.content)

			if result < tt.minPages {
				t.Errorf("expected at least %d pages, got %d", tt.minPages, result)
			}
			if result > tt.maxPages {
				t.Errorf("expected at most %d pages, got %d", tt.maxPages, result)
			}
		})
	}
}
