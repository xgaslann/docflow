package com.docflow.formats;

import java.util.*;

/**
 * CSV format converter.
 */
public class CSVConverter {

    private char delimiter = ',';
    private boolean hasHeader = true;
    private String tableTitle = null;

    public CSVConverter() {
    }

    public CSVConverter(char delimiter, boolean hasHeader) {
        this.delimiter = delimiter;
        this.hasHeader = hasHeader;
    }

    /**
     * Convert CSV to Markdown.
     */
    public ConvertResult toMarkdown(byte[] csvData, String filename) {
        return toMarkdown(new String(csvData), filename);
    }

    public ConvertResult toMarkdown(String csvData, String filename) {
        try {
            List<String[]> rows = parseCSV(csvData);

            if (rows.isEmpty()) {
                return ConvertResult.error("Empty CSV file");
            }

            StringBuilder sb = new StringBuilder();

            // Frontmatter
            sb.append("---\n");
            sb.append("source: ").append(filename).append("\n");
            sb.append("format: csv\n");
            sb.append("rows: ").append(rows.size()).append("\n");
            sb.append("columns: ").append(rows.get(0).length).append("\n");
            sb.append("---\n\n");

            // Title
            String title = tableTitle != null ? tableTitle : filename.replace(".csv", "").replace("_", " ");
            sb.append("# ").append(title).append("\n\n");

            // Table
            String[] header;
            List<String[]> dataRows;

            if (hasHeader && !rows.isEmpty()) {
                header = rows.get(0);
                dataRows = rows.subList(1, rows.size());
            } else {
                header = new String[rows.get(0).length];
                for (int i = 0; i < header.length; i++) {
                    header[i] = "Column " + (i + 1);
                }
                dataRows = rows;
            }

            // Header row
            sb.append("| ").append(String.join(" | ", escapeCell(header))).append(" |\n");

            // Separator
            sb.append("| ");
            for (int i = 0; i < header.length; i++) {
                sb.append("---");
                if (i < header.length - 1)
                    sb.append(" | ");
            }
            sb.append(" |\n");

            // Data rows
            for (String[] row : dataRows) {
                String[] paddedRow = Arrays.copyOf(row, header.length);
                for (int i = 0; i < paddedRow.length; i++) {
                    if (paddedRow[i] == null)
                        paddedRow[i] = "";
                }
                sb.append("| ").append(String.join(" | ", escapeCell(paddedRow))).append(" |\n");
            }

            return ConvertResult.success(sb.toString(), "csv");

        } catch (Exception e) {
            return ConvertResult.error(e.getMessage());
        }
    }

    /**
     * Convert Markdown to CSV.
     */
    public ConvertResult fromMarkdown(String markdown) {
        try {
            List<String[]> tables = extractTables(markdown);

            if (tables.isEmpty()) {
                return ConvertResult.error("No tables found");
            }

            StringBuilder sb = new StringBuilder();
            for (String[] row : tables) {
                sb.append(String.join(String.valueOf(delimiter), row)).append("\n");
            }

            return ConvertResult.success(sb.toString(), "csv");

        } catch (Exception e) {
            return ConvertResult.error(e.getMessage());
        }
    }

    private List<String[]> parseCSV(String data) {
        List<String[]> rows = new ArrayList<>();
        String[] lines = data.split("\n");

        for (String line : lines) {
            if (line.trim().isEmpty())
                continue;
            String[] cells = line.split(String.valueOf(delimiter), -1);
            for (int i = 0; i < cells.length; i++) {
                cells[i] = cells[i].trim();
            }
            rows.add(cells);
        }

        return rows;
    }

    private String[] escapeCell(String[] cells) {
        String[] escaped = new String[cells.length];
        for (int i = 0; i < cells.length; i++) {
            escaped[i] = cells[i] != null ? cells[i].replace("|", "\\|") : "";
        }
        return escaped;
    }

    private List<String[]> extractTables(String markdown) {
        List<String[]> rows = new ArrayList<>();

        for (String line : markdown.split("\n")) {
            line = line.trim();
            if (line.startsWith("|") && line.endsWith("|")) {
                // Skip separator
                if (line.replace("|", "").replace("-", "").replace(":", "").trim().isEmpty()) {
                    continue;
                }
                String[] parts = line.split("\\|");
                List<String> cells = new ArrayList<>();
                for (int i = 1; i < parts.length - 1; i++) {
                    cells.add(parts[i].trim());
                }
                rows.add(cells.toArray(new String[0]));
            }
        }

        return rows;
    }
}
