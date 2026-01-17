package formats

import (
	"fmt"
	"regexp"
	"strings"
)

// TXTConverter handles plain text to Markdown conversion.
type TXTConverter struct {
	DetectStructure bool
	LineBreakMode   string // "paragraph" or "preserve"
}

// TXTResult represents the result of a TXT conversion.
type TXTResult struct {
	Success  bool
	Content  string
	Error    string
	Metadata map[string]interface{}
}

// NewTXTConverter creates a new TXT converter with default settings.
func NewTXTConverter() *TXTConverter {
	return &TXTConverter{
		DetectStructure: true,
		LineBreakMode:   "paragraph",
	}
}

// ToMarkdown converts plain text to Markdown.
func (c *TXTConverter) ToMarkdown(textData []byte, filename string) TXTResult {
	text := string(textData)
	lines := strings.Split(text, "\n")
	wordCount := len(strings.Fields(text))

	var sb strings.Builder

	// Metadata
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("source: %s\n", filename))
	sb.WriteString("format: txt\n")
	sb.WriteString(fmt.Sprintf("lines: %d\n", len(lines)))
	sb.WriteString(fmt.Sprintf("words: %d\n", wordCount))
	sb.WriteString(fmt.Sprintf("characters: %d\n", len(text)))
	sb.WriteString("---\n\n")

	// Title
	title := strings.TrimSuffix(filename, ".txt")
	title = strings.ReplaceAll(title, "_", " ")
	title = strings.Title(title)
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))

	if c.DetectStructure {
		sb.WriteString(c.detectAndConvert(text))
	} else if c.LineBreakMode == "preserve" {
		sb.WriteString(text)
	} else {
		sb.WriteString(c.paragraphize(text))
	}

	return TXTResult{
		Success: true,
		Content: sb.String(),
		Metadata: map[string]interface{}{
			"lines": len(lines),
			"words": wordCount,
		},
	}
}

// FromMarkdown converts Markdown to plain text.
func (c *TXTConverter) FromMarkdown(markdown string) TXTResult {
	text := markdown

	// Remove frontmatter
	if strings.HasPrefix(text, "---") {
		idx := strings.Index(text[3:], "---")
		if idx > 0 {
			text = strings.TrimSpace(text[idx+6:])
		}
	}

	// Remove markdown formatting
	// Headers
	re := regexp.MustCompile(`(?m)^#{1,6}\s+`)
	text = re.ReplaceAllString(text, "")

	// Bold/italic
	text = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`\*(.+?)\*`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`__(.+?)__`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`_(.+?)_`).ReplaceAllString(text, "$1")

	// Code
	text = regexp.MustCompile("`(.+?)`").ReplaceAllString(text, "$1")
	text = regexp.MustCompile("(?s)```.*?```").ReplaceAllStringFunc(text, func(s string) string {
		return strings.ReplaceAll(s, "```", "")
	})

	// Links
	text = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`).ReplaceAllString(text, "$1")

	// Images
	text = regexp.MustCompile(`!\[([^\]]*)\]\([^)]+\)`).ReplaceAllString(text, "[Image: $1]")

	// Lists
	text = regexp.MustCompile(`(?m)^[-*+]\s+`).ReplaceAllString(text, "• ")
	text = regexp.MustCompile(`(?m)^\d+\.\s+`).ReplaceAllString(text, "")

	// Tables - simple format
	var resultLines []string
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "|") {
			cells := strings.Split(line, "|")
			var filtered []string
			for _, cell := range cells {
				cell = strings.TrimSpace(cell)
				if cell != "" && !isTableSeparator(cell) {
					filtered = append(filtered, cell)
				}
			}
			if len(filtered) > 0 {
				resultLines = append(resultLines, strings.Join(filtered, " | "))
			}
		} else {
			resultLines = append(resultLines, line)
		}
	}
	text = strings.Join(resultLines, "\n")

	// Clean up
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")
	text = strings.TrimSpace(text)

	return TXTResult{
		Success: true,
		Content: text,
	}
}

func (c *TXTConverter) detectAndConvert(text string) string {
	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		stripped := strings.TrimSpace(line)

		if stripped == "" {
			result = append(result, "")
			continue
		}

		// Detect headers (ALL CAPS)
		if stripped == strings.ToUpper(stripped) &&
			len(stripped) > 3 &&
			len(stripped) < 80 &&
			!strings.HasPrefix(stripped, "•") &&
			!strings.HasPrefix(stripped, "-") &&
			len(strings.Fields(stripped)) >= 2 {
			result = append(result, fmt.Sprintf("\n## %s\n", strings.Title(strings.ToLower(stripped))))
			continue
		}

		// Detect bullet points
		for _, bullet := range []string{"•", "●", "○", "▪", "▸"} {
			if strings.HasPrefix(stripped, bullet) {
				result = append(result, fmt.Sprintf("- %s", strings.TrimSpace(stripped[len(bullet):])))
				continue
			}
		}

		// Detect numbered lists
		re := regexp.MustCompile(`^(\d+)[.)]\s+(.+)$`)
		if matches := re.FindStringSubmatch(stripped); matches != nil {
			result = append(result, fmt.Sprintf("%s. %s", matches[1], matches[2]))
			continue
		}

		result = append(result, stripped)
	}

	return strings.Join(result, "\n")
}

func (c *TXTConverter) paragraphize(text string) string {
	paragraphs := regexp.MustCompile(`\n\s*\n`).Split(text, -1)
	var result []string

	for _, para := range paragraphs {
		// Join lines within paragraph
		para = strings.Join(strings.Fields(para), " ")
		if strings.TrimSpace(para) != "" {
			result = append(result, para)
		}
	}

	return strings.Join(result, "\n\n")
}

func isTableSeparator(s string) bool {
	clean := strings.ReplaceAll(s, "-", "")
	clean = strings.ReplaceAll(clean, ":", "")
	return strings.TrimSpace(clean) == ""
}
