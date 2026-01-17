package rag

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xgaslan/docflow/sdks/go/docflow/config"
)

// Chunker splits text into chunks based on configuration.
type Chunker struct {
	Config config.RAGConfig
}

// NewChunker creates a new chunker.
func NewChunker(config config.RAGConfig) *Chunker {
	return &Chunker{
		Config: config,
	}
}

// Chunk represents a text chunk.
// Reusing logic from types.go which defines Chunk struct.

// Chunk splits markdown content into RAG-optimized chunks.
func (c *Chunker) Chunk(markdown string) []Chunk {
	// Extract frontmatter
	content, _ := c.extractFrontmatter(markdown)

	// Split by headings if configured
	var sections []sectionData
	if c.Config.RespectHeadings {
		sections = c.splitByHeadings(content)
	} else {
		sections = []sectionData{{"", content}}
	}

	var chunks []Chunk
	chunkIndex := 0
	charOffset := 0

	for _, section := range sections {
		sectionChunks := c.chunkSection(section.content, section.title, chunkIndex, charOffset)
		chunks = append(chunks, sectionChunks...)
		chunkIndex += len(sectionChunks)
		charOffset += len(section.content)
	}

	// Add overlap
	if c.Config.ChunkOverlap > 0 {
		chunks = c.addOverlap(chunks, content)
	}

	// Add chunk markers
	if c.Config.AddChunkMarkers {
		for i := range chunks {
			chunks[i].Content = fmt.Sprintf("%s\n\n<!-- chunk_boundary: %d -->", chunks[i].Content, i)
		}
	}

	return chunks
}

type sectionData struct {
	title   string
	content string
}

func (c *Chunker) extractFrontmatter(markdown string) (string, string) {
	if strings.HasPrefix(markdown, "---") {
		idx := strings.Index(markdown[3:], "---")
		if idx > 0 {
			frontmatter := strings.TrimSpace(markdown[3 : idx+3])
			content := strings.TrimSpace(markdown[idx+6:])
			return content, frontmatter
		}
	}
	return markdown, ""
}

func (c *Chunker) splitByHeadings(content string) []sectionData {
	pattern := regexp.MustCompile(`(?m)^(#{1,6})\s+(.+)$`)
	matches := pattern.FindAllStringSubmatchIndex(content, -1)

	if len(matches) == 0 {
		return []sectionData{{"", content}}
	}

	var sections []sectionData
	lastEnd := 0
	lastTitle := ""

	for _, match := range matches {
		// Content before this heading
		if match[0] > lastEnd {
			sectionContent := strings.TrimSpace(content[lastEnd:match[0]])
			if sectionContent != "" {
				sections = append(sections, sectionData{lastTitle, sectionContent})
			}
		}
		lastTitle = strings.TrimSpace(content[match[4]:match[5]]) // Group 2
		lastEnd = match[0]
	}

	// Remaining content
	if lastEnd < len(content) {
		remaining := strings.TrimSpace(content[lastEnd:])
		if remaining != "" {
			sections = append(sections, sectionData{lastTitle, remaining})
		}
	}

	if len(sections) == 0 {
		sections = []sectionData{{"", content}}
	}

	return sections
}

func (c *Chunker) chunkSection(content, sectionTitle string, startIndex, charOffset int) []Chunk {
	var chunks []Chunk

	lines := strings.Split(content, "\n")
	currentChunk := ""
	currentStart := charOffset
	chunkIdx := startIndex

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Check for protected blocks (tables, code blocks)
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			// Find end of code block
			blockLines := []string{line}
			i++
			for i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), "```") {
				blockLines = append(blockLines, lines[i])
				i++
			}
			if i < len(lines) {
				blockLines = append(blockLines, lines[i])
			}
			block := strings.Join(blockLines, "\n")

			if len(currentChunk)+len(block) > c.Config.ChunkSize && currentChunk != "" {
				chunks = append(chunks, c.createChunk(
					strings.TrimSpace(currentChunk),
					chunkIdx,
					currentStart,
					currentStart+len(currentChunk),
					sectionTitle,
				))
				chunkIdx++
				currentChunk = ""
				currentStart = charOffset + c.countChars(lines[:i])
			}
			currentChunk += block + "\n"
			i++
			continue
		}

		// Check for table
		if strings.HasPrefix(strings.TrimSpace(line), "|") && strings.HasSuffix(strings.TrimSpace(line), "|") {
			tableLines := []string{line}
			i++
			for i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "|") {
				tableLines = append(tableLines, lines[i])
				i++
			}
			block := strings.Join(tableLines, "\n")

			if len(currentChunk)+len(block) > c.Config.ChunkSize && currentChunk != "" {
				chunks = append(chunks, c.createChunk(
					strings.TrimSpace(currentChunk),
					chunkIdx,
					currentStart,
					currentStart+len(currentChunk),
					sectionTitle,
				))
				chunkIdx++
				currentChunk = ""
				currentStart = charOffset + c.countChars(lines[:i-len(tableLines)])
			}
			currentChunk += block + "\n"
			continue
		}

		// Regular line
		if len(currentChunk)+len(line) > c.Config.ChunkSize && currentChunk != "" {
			chunks = append(chunks, c.createChunk(
				strings.TrimSpace(currentChunk),
				chunkIdx,
				currentStart,
				currentStart+len(currentChunk),
				sectionTitle,
			))
			chunkIdx++
			currentChunk = ""
			currentStart = charOffset + c.countChars(lines[:i])
		}
		currentChunk += line + "\n"
		i++
	}

	if strings.TrimSpace(currentChunk) != "" {
		chunks = append(chunks, c.createChunk(
			strings.TrimSpace(currentChunk),
			chunkIdx,
			currentStart,
			currentStart+len(currentChunk),
			sectionTitle,
		))
	}

	return chunks
}

func (c *Chunker) createChunk(content string, index, startChar, endChar int, sectionTitle string) Chunk {
	hasTable := strings.Contains(content, "|") && strings.Contains(content, "---")
	hasImage := strings.Contains(content, "![") || strings.Contains(content, "[Image:")
	hasCode := strings.Contains(content, "```")

	// Determine content type
	contentType := "text"
	if hasTable {
		contentType = "table"
	} else if hasCode {
		contentType = "code"
	} else if hasImage {
		contentType = "image"
	}

	return Chunk{
		Content:   content,
		Index:     index,
		StartChar: startChar,
		EndChar:   endChar,
		Metadata: ChunkMetadata{
			SectionTitle: sectionTitle,
			HeadingPath:  c.extractHeadingPath(content),
			HasTable:     hasTable,
			HasImage:     hasImage,
			HasCode:      hasCode,
			ContentType:  contentType,
		},
	}
}

func (c *Chunker) extractHeadingPath(content string) []string {
	var headings []string
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "#") {
			text := strings.TrimLeft(line, "#")
			headings = append(headings, strings.TrimSpace(text))
		}
	}
	return headings
}

func (c *Chunker) addOverlap(chunks []Chunk, _ string) []Chunk {
	if len(chunks) <= 1 {
		return chunks
	}

	overlapSize := c.Config.ChunkOverlap

	for i := 1; i < len(chunks); i++ {
		prevChunk := chunks[i-1]
		content := prevChunk.Content

		// Get overlap from end of previous chunk
		overlapText := content
		if len(content) > overlapSize {
			overlapText = content[len(content)-overlapSize:]
		}

		// Find a good break point
		for _, breakStr := range []string{"\n\n", ". ", "\n"} {
			idx := strings.Index(overlapText, breakStr)
			if idx > 0 {
				overlapText = overlapText[idx+len(breakStr):]
				break
			}
		}

		if strings.TrimSpace(overlapText) != "" {
			chunks[i].Content = fmt.Sprintf("[...] %s\n\n%s", strings.TrimSpace(overlapText), chunks[i].Content)
		}
	}

	return chunks
}

func (c *Chunker) countChars(lines []string) int {
	total := 0
	for _, line := range lines {
		total += len(line) + 1 // +1 for newline
	}
	return total
}
