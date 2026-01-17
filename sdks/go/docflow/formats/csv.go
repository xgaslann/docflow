package formats

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
)

// CSVConverter handles CSV to Markdown conversion.
type CSVConverter struct {
	Delimiter  rune
	HasHeader  bool
	TableTitle string
}

// CSVResult represents the result of a CSV conversion.
type CSVResult struct {
	Success  bool
	Content  string
	Error    string
	Metadata map[string]interface{}
}

// NewCSVConverter creates a new CSV converter with default settings.
func NewCSVConverter() *CSVConverter {
	return &CSVConverter{
		Delimiter: ',',
		HasHeader: true,
	}
}

// ToMarkdown converts CSV content to Markdown table format.
func (c *CSVConverter) ToMarkdown(csvData []byte, filename string) CSVResult {
	reader := csv.NewReader(bytes.NewReader(csvData))
	reader.Comma = c.Delimiter

	rows, err := reader.ReadAll()
	if err != nil {
		return CSVResult{Success: false, Error: fmt.Sprintf("failed to parse CSV: %v", err)}
	}

	if len(rows) == 0 {
		return CSVResult{Success: false, Error: "empty CSV file"}
	}

	var sb strings.Builder

	// Metadata frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("source: %s\n", filename))
	sb.WriteString("format: csv\n")
	sb.WriteString(fmt.Sprintf("rows: %d\n", len(rows)))
	if len(rows) > 0 {
		sb.WriteString(fmt.Sprintf("columns: %d\n", len(rows[0])))
	}
	sb.WriteString("---\n\n")

	// Title
	title := c.TableTitle
	if title == "" {
		title = strings.TrimSuffix(filename, ".csv")
		title = strings.ReplaceAll(title, "_", " ")
		title = strings.Title(title)
	}
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))

	// Table
	var header, dataRows [][]string
	if c.HasHeader && len(rows) > 0 {
		header = [][]string{rows[0]}
		dataRows = rows[1:]
	} else {
		// Generate header
		h := make([]string, len(rows[0]))
		for i := range h {
			h[i] = fmt.Sprintf("Column %d", i+1)
		}
		header = [][]string{h}
		dataRows = rows
	}

	// Header row
	sb.WriteString("| ")
	sb.WriteString(strings.Join(escapeCSVCells(header[0]), " | "))
	sb.WriteString(" |\n")

	// Separator
	sb.WriteString("| ")
	seps := make([]string, len(header[0]))
	for i := range seps {
		seps[i] = "---"
	}
	sb.WriteString(strings.Join(seps, " | "))
	sb.WriteString(" |\n")

	// Data rows
	for _, row := range dataRows {
		// Pad row if needed
		paddedRow := make([]string, len(header[0]))
		copy(paddedRow, row)
		sb.WriteString("| ")
		sb.WriteString(strings.Join(escapeCSVCells(paddedRow), " | "))
		sb.WriteString(" |\n")
	}

	return CSVResult{
		Success: true,
		Content: sb.String(),
		Metadata: map[string]interface{}{
			"rows":    len(rows),
			"columns": len(rows[0]),
		},
	}
}

// FromMarkdown extracts tables from Markdown and converts to CSV.
func (c *CSVConverter) FromMarkdown(markdown string) CSVResult {
	tables := extractTablesFromMD(markdown)
	if len(tables) == 0 {
		return CSVResult{Success: false, Error: "no tables found in markdown"}
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = c.Delimiter

	for _, row := range tables[0] {
		if err := writer.Write(row); err != nil {
			return CSVResult{Success: false, Error: fmt.Sprintf("failed to write CSV: %v", err)}
		}
	}
	writer.Flush()

	return CSVResult{
		Success: true,
		Content: buf.String(),
		Metadata: map[string]interface{}{
			"tables_found": len(tables),
		},
	}
}

func escapeCSVCells(cells []string) []string {
	escaped := make([]string, len(cells))
	for i, cell := range cells {
		escaped[i] = strings.ReplaceAll(cell, "|", "\\|")
		escaped[i] = strings.ReplaceAll(escaped[i], "\n", " ")
	}
	return escaped
}

func extractTablesFromMD(markdown string) [][][]string {
	var tables [][][]string
	var currentTable [][]string
	inTable := false

	lines := strings.Split(markdown, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|") {
			// Skip separator rows
			content := strings.Trim(line, "|")
			content = strings.ReplaceAll(content, "-", "")
			content = strings.ReplaceAll(content, ":", "")
			content = strings.TrimSpace(content)
			if content == "" {
				continue
			}

			cells := strings.Split(line, "|")
			if len(cells) > 2 {
				row := make([]string, 0, len(cells)-2)
				for _, cell := range cells[1 : len(cells)-1] {
					row = append(row, strings.TrimSpace(cell))
				}
				currentTable = append(currentTable, row)
			}
			inTable = true
		} else {
			if inTable && len(currentTable) > 0 {
				tables = append(tables, currentTable)
				currentTable = nil
			}
			inTable = false
		}
	}

	if len(currentTable) > 0 {
		tables = append(tables, currentTable)
	}

	return tables
}
