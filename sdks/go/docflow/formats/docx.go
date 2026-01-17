package formats

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xgaslan/docflow/sdks/go/docflow"
	"github.com/xgaslan/docflow/sdks/go/docflow/rag"
)

// DOCXConverter converts DOCX files to/from Markdown.
type DOCXConverter struct {
	// ExtractImages extracts images from DOCX
	ExtractImages bool
	// PreserveFormatting attempts to preserve formatting
	PreserveFormatting bool
}

// NewDOCXConverter creates a new DOCX converter.
func NewDOCXConverter() *DOCXConverter {
	return &DOCXConverter{
		ExtractImages:      true,
		PreserveFormatting: true,
	}
}

// ToMarkdown converts DOCX data to Markdown.
func (c *DOCXConverter) ToMarkdown(data []byte, filename string) (*docflow.ConvertResult, error) {
	// Note: Requires github.com/nguyenthenguyen/docx or similar
	// For full implementation, use a proper DOCX library

	doc, err := openDOCXFromBytes(data)
	if err != nil {
		return &docflow.ConvertResult{
			Success: false,
			Error:   fmt.Errorf("failed to open DOCX: %w", err),
		}, nil
	}

	content := doc.GetText()
	images := []rag.ExtractedImage{}

	if c.ExtractImages {
		imgs, err := doc.GetImages()
		if err == nil {
			images = imgs
		}
	}

	// Convert to markdown
	markdown := c.textToMarkdown(content, filename)

	metadata := map[string]interface{}{
		"filename":    filename,
		"image_count": len(images),
	}

	// Get document properties if available
	if props := doc.GetProperties(); props != nil {
		if props.Title != "" {
			metadata["title"] = props.Title
		}
		if props.Author != "" {
			metadata["author"] = props.Author
		}
	}

	return &docflow.ConvertResult{
		Success:  true,
		Content:  markdown,
		Format:   "docx",
		Images:   images,
		Metadata: metadata,
	}, nil
}

// FromMarkdown converts Markdown to DOCX.
func (c *DOCXConverter) FromMarkdown(content string, filename string) ([]byte, error) {
	doc := newDOCXDocument()

	// Parse markdown and add to document
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "# ") {
			doc.AddHeading(strings.TrimPrefix(line, "# "), 1)
		} else if strings.HasPrefix(line, "## ") {
			doc.AddHeading(strings.TrimPrefix(line, "## "), 2)
		} else if strings.HasPrefix(line, "### ") {
			doc.AddHeading(strings.TrimPrefix(line, "### "), 3)
		} else if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			doc.AddListItem(line[2:])
		} else if line != "" {
			doc.AddParagraph(c.processInlineFormatting(line))
		}
	}

	return doc.Save()
}

func (c *DOCXConverter) textToMarkdown(text, title string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))

	lines := strings.Split(text, "\n")
	var processedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			if len(processedLines) > 0 && processedLines[len(processedLines)-1] != "" {
				processedLines = append(processedLines, "")
			}
			continue
		}

		// Detect potential headers (short, capitalized lines)
		if c.isPotentialHeader(line) {
			line = "## " + strings.Title(strings.ToLower(line))
		}

		// Detect bullet points
		for _, bullet := range []string{"•", "●", "○", "◦"} {
			if strings.HasPrefix(line, bullet) {
				line = "- " + strings.TrimSpace(line[len(bullet):])
				break
			}
		}

		processedLines = append(processedLines, line)
	}

	// Join with proper spacing
	inParagraph := false
	for i, line := range processedLines {
		if line == "" {
			if inParagraph {
				sb.WriteString("\n\n")
				inParagraph = false
			}
			continue
		}

		if strings.HasPrefix(line, "##") || strings.HasPrefix(line, "- ") {
			if inParagraph {
				sb.WriteString("\n\n")
			}
			sb.WriteString(line + "\n")
			inParagraph = false
		} else {
			if inParagraph && i > 0 && processedLines[i-1] != "" {
				sb.WriteString(" ")
			}
			sb.WriteString(line)
			inParagraph = true
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

func (c *DOCXConverter) isPotentialHeader(line string) bool {
	return len(line) > 3 && len(line) < 60 &&
		line == strings.ToUpper(line) &&
		len(strings.Fields(line)) >= 2
}

func (c *DOCXConverter) processInlineFormatting(text string) string {
	// Convert markdown bold to plain (for DOCX, actual formatting would be applied differently)
	boldRe := regexp.MustCompile(`\*\*(.+?)\*\*`)
	text = boldRe.ReplaceAllString(text, "$1")

	italicRe := regexp.MustCompile(`\*(.+?)\*`)
	text = italicRe.ReplaceAllString(text, "$1")

	return text
}

// DOCX document interface (requires docx library)
type docxDocument interface {
	GetText() string
	GetImages() ([]rag.ExtractedImage, error)
	GetProperties() *docxProperties
	AddHeading(text string, level int)
	AddParagraph(text string)
	AddListItem(text string)
	Save() ([]byte, error)
}

type docxProperties struct {
	Title  string
	Author string
}

// Placeholder functions - require docx dependency
func openDOCXFromBytes(data []byte) (docxDocument, error) {
	// Requires: go get github.com/nguyenthenguyen/docx
	// Or: go get github.com/unidoc/unioffice
	return nil, fmt.Errorf("docx library not installed")
}

func newDOCXDocument() docxDocument {
	return nil
}
