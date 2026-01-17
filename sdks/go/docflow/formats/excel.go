package formats

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xgaslan/docflow/sdks/go/docflow"
)

// ExcelConverter converts Excel files to/from Markdown.
type ExcelConverter struct {
	// IncludeAllSheets includes all sheets in conversion
	IncludeAllSheets bool
	// SheetSeparator is the separator between sheets
	SheetSeparator string
}

// NewExcelConverter creates a new Excel converter.
func NewExcelConverter() *ExcelConverter {
	return &ExcelConverter{
		IncludeAllSheets: true,
		SheetSeparator:   "\n\n---\n\n",
	}
}

// ToMarkdown converts Excel data to Markdown.
func (c *ExcelConverter) ToMarkdown(data []byte, filename string) (*docflow.ConvertResult, error) {
	// Note: Requires github.com/xuri/excelize/v2
	// go get github.com/xuri/excelize/v2

	xlsx, err := openExcelFromBytes(data)
	if err != nil {
		return &docflow.ConvertResult{
			Success: false,
			Error:   fmt.Errorf("failed to open Excel file: %w", err),
		}, nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", filename))

	sheets := xlsx.GetSheetList()
	for i, sheet := range sheets {
		if !c.IncludeAllSheets && i > 0 {
			break
		}

		if i > 0 {
			sb.WriteString(c.SheetSeparator)
		}

		sb.WriteString(fmt.Sprintf("## %s\n\n", sheet))

		rows, err := xlsx.GetRows(sheet)
		if err != nil {
			continue
		}

		if len(rows) == 0 {
			sb.WriteString("*Empty sheet*\n")
			continue
		}

		// Convert to markdown table
		sb.WriteString(c.rowsToMarkdownTable(rows))
	}

	return &docflow.ConvertResult{
		Success: true,
		Content: sb.String(),
		Format:  "xlsx",
		Metadata: map[string]interface{}{
			"sheet_count": len(sheets),
			"filename":    filename,
		},
	}, nil
}

// FromMarkdown converts Markdown table to Excel.
func (c *ExcelConverter) FromMarkdown(content string, filename string) ([]byte, error) {
	xlsx := newExcelFile()

	tables := c.extractMarkdownTables(content)

	for i, table := range tables {
		sheetName := fmt.Sprintf("Sheet%d", i+1)
		if i == 0 {
			xlsx.SetSheetName("Sheet1", sheetName)
		} else {
			xlsx.NewSheet(sheetName)
		}

		for rowIdx, row := range table {
			for colIdx, cell := range row {
				cellRef := fmt.Sprintf("%s%d", columnToLetter(colIdx+1), rowIdx+1)
				xlsx.SetCellValue(sheetName, cellRef, cell)
			}
		}
	}

	buf := new(bytes.Buffer)
	if err := xlsx.Write(buf); err != nil {
		return nil, fmt.Errorf("failed to write Excel: %w", err)
	}

	return buf.Bytes(), nil
}

func (c *ExcelConverter) rowsToMarkdownTable(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	var sb strings.Builder
	maxCols := 0
	for _, row := range rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	// Header
	sb.WriteString("|")
	for i := 0; i < maxCols; i++ {
		if i < len(rows[0]) {
			sb.WriteString(fmt.Sprintf(" %s |", rows[0][i]))
		} else {
			sb.WriteString(" |")
		}
	}
	sb.WriteString("\n")

	// Separator
	sb.WriteString("|")
	for i := 0; i < maxCols; i++ {
		sb.WriteString(" --- |")
	}
	sb.WriteString("\n")

	// Data rows
	for _, row := range rows[1:] {
		sb.WriteString("|")
		for i := 0; i < maxCols; i++ {
			if i < len(row) {
				sb.WriteString(fmt.Sprintf(" %s |", row[i]))
			} else {
				sb.WriteString(" |")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (c *ExcelConverter) extractMarkdownTables(content string) [][][]string {
	var tables [][][]string
	lines := strings.Split(content, "\n")

	var currentTable [][]string
	inTable := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|") {
			// Skip separator line
			if strings.Contains(line, "---") {
				continue
			}

			cells := strings.Split(line, "|")
			var row []string
			for _, cell := range cells[1 : len(cells)-1] {
				row = append(row, strings.TrimSpace(cell))
			}
			currentTable = append(currentTable, row)
			inTable = true
		} else if inTable && len(currentTable) > 0 {
			tables = append(tables, currentTable)
			currentTable = nil
			inTable = false
		}
	}

	if len(currentTable) > 0 {
		tables = append(tables, currentTable)
	}

	return tables
}

func columnToLetter(col int) string {
	result := ""
	for col > 0 {
		col--
		result = string(rune('A'+col%26)) + result
		col /= 26
	}
	return result
}

// Excel file interface (requires excelize)
type excelFile interface {
	GetSheetList() []string
	GetRows(sheet string) ([][]string, error)
	SetSheetName(old, new string) error
	NewSheet(name string) (int, error)
	SetCellValue(sheet, cell string, value interface{}) error
	Write(w *bytes.Buffer) error
}

// Placeholder functions - require excelize dependency
func openExcelFromBytes(data []byte) (excelFile, error) {
	// Requires: go get github.com/xuri/excelize/v2
	// import "github.com/xuri/excelize/v2"
	// return excelize.OpenReader(bytes.NewReader(data))
	return nil, fmt.Errorf("excelize not installed: go get github.com/xuri/excelize/v2")
}

func newExcelFile() excelFile {
	// return excelize.NewFile()
	return nil
}
