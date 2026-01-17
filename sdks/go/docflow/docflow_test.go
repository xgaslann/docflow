package docflow

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgaslan/docflow/sdks/go/docflow/storage"
)

func TestMarkdownParser_ToHTML(t *testing.T) {
	parser := NewMarkdownParser()

	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "heading",
			input:    "# Hello World",
			contains: []string{"<h1", "Hello World", "</h1>"},
		},
		{
			name:     "paragraph",
			input:    "This is a paragraph.",
			contains: []string{"<p>", "This is a paragraph.", "</p>"},
		},
		{
			name:     "bold",
			input:    "This is **bold** text.",
			contains: []string{"<strong>", "bold", "</strong>"},
		},
		{
			name:     "italic",
			input:    "This is *italic* text.",
			contains: []string{"<em>", "italic", "</em>"},
		},
		{
			name:     "code",
			input:    "This is `code` text.",
			contains: []string{"<code>", "code", "</code>"},
		},
		{
			name:     "list",
			input:    "- Item 1\n- Item 2",
			contains: []string{"<ul>", "<li>", "Item 1", "Item 2", "</li>", "</ul>"},
		},
		{
			name:     "table",
			input:    "| A | B |\n|---|---|\n| 1 | 2 |",
			contains: []string{"<table>", "<th>", "<td>", "</table>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := parser.ToHTML(tt.input)
			require.NoError(t, err)

			for _, s := range tt.contains {
				assert.Contains(t, html, s)
			}
		})
	}
}

func TestMarkdownParser_MergeFiles(t *testing.T) {
	parser := NewMarkdownParser()

	files := []MDFile{
		{Name: "second.md", Content: "# Second", Order: 1},
		{Name: "first.md", Content: "# First", Order: 0},
		{Name: "third.md", Content: "# Third", Order: 2},
	}

	merged := parser.MergeFiles(files)

	// Check that files are merged in order
	firstIdx := len("# First")
	assert.True(t, len(merged) > firstIdx)
	assert.Contains(t, merged, "# First")
	assert.Contains(t, merged, "# Second")
	assert.Contains(t, merged, "# Third")

	// First should appear before Second
	assert.Less(t,
		indexOf(merged, "# First"),
		indexOf(merged, "# Second"),
	)

	// Second should appear before Third
	assert.Less(t,
		indexOf(merged, "# Second"),
		indexOf(merged, "# Third"),
	)
}

func TestMarkdownParser_EstimatePageCount(t *testing.T) {
	parser := NewMarkdownParser()

	tests := []struct {
		name     string
		length   int
		expected int
	}{
		{"short", 100, 1},
		{"one page", 3000, 1},
		{"two pages", 6000, 2},
		{"five pages", 15000, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := string(make([]byte, tt.length))
			pages := parser.EstimatePageCount(content)
			assert.Equal(t, tt.expected, pages)
		})
	}
}

func TestLocalStorage(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "docflow-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	s, err := storage.NewLocalStorage(tmpDir)
	require.NoError(t, err)

	testData := []byte("Hello, World!")
	testPath := "test/file.txt"

	// Test Save
	err = s.Save(testPath, testData)
	require.NoError(t, err)

	// Test Exists
	exists, err := s.Exists(testPath)
	require.NoError(t, err)
	assert.True(t, exists)

	// Test Load
	loaded, err := s.Load(testPath)
	require.NoError(t, err)
	assert.Equal(t, testData, loaded)

	// Test List
	files, err := s.List("test")
	require.NoError(t, err)
	assert.Contains(t, files, "file.txt")

	// Test Delete
	err = s.Delete(testPath)
	require.NoError(t, err)

	exists, err = s.Exists(testPath)
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestTemplate_Generate(t *testing.T) {
	template := NewTemplate()

	html := template.Generate("<h1>Test</h1>")

	assert.Contains(t, html, "<!DOCTYPE html>")
	assert.Contains(t, html, "<h1>Test</h1>")
	assert.Contains(t, html, "@page")
	assert.Contains(t, html, "size: A4")
}

func TestNewMDFile(t *testing.T) {
	file := NewMDFile("test.md", "# Test Content")

	assert.Equal(t, "test.md", file.Name)
	assert.Equal(t, "test.md", file.ID)
	assert.Equal(t, "# Test Content", file.Content)
	assert.Equal(t, 0, file.Order)
}

func TestNewMDFileWithOrder(t *testing.T) {
	file := NewMDFileWithOrder("test.md", "# Test", 5)

	assert.Equal(t, "test.md", file.Name)
	assert.Equal(t, 5, file.Order)
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"with space", "with_space"},
		{"with/slash", "with_slash"},
		{"with:colon", "with_colon"},
		{"with*star", "with_star"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeFilename(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Integration test (requires Chrome)
func TestConverter_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Skip if Chrome not available
	if _, err := exec.LookPath("google-chrome"); err != nil {
		if _, err := exec.LookPath("chromium"); err != nil {
			t.Skip("Chrome/Chromium not found, skipping integration test")
		}
	}

	tmpDir, err := os.MkdirTemp("", "docflow-integration-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	converter := NewConverter(WithLocalStorage(tmpDir))

	files := []MDFile{
		NewMDFile("test.md", "# Test Document\n\nThis is a test."),
	}

	result, err := converter.ConvertToPDF(context.Background(), files, ConvertOptions{
		MergeMode: "separate",
	})
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.FilePaths, 1)
}
