package com.docflow.formats;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Converts Excel files to/from Markdown.
 * 
 * Requires Apache POI dependency:
 * <dependency>
 * <groupId>org.apache.poi</groupId>
 * <artifactId>poi-ooxml</artifactId>
 * <version>5.2.5</version>
 * <optional>true</optional>
 * </dependency>
 */
public class ExcelConverter {

    private boolean includeAllSheets = true;
    private String sheetSeparator = "\n\n---\n\n";

    public ExcelConverter() {
    }

    public ExcelConverter(boolean includeAllSheets) {
        this.includeAllSheets = includeAllSheets;
    }

    /**
     * Convert Excel data to Markdown.
     */
    public ConvertResult toMarkdown(byte[] data, String filename) {
        try {
            // Requires Apache POI
            // import org.apache.poi.ss.usermodel.*;
            // import org.apache.poi.xssf.usermodel.XSSFWorkbook;

            Object workbook = openWorkbook(data);
            if (workbook == null) {
                return new ConvertResult(false, "Apache POI not available. Add poi-ooxml dependency.");
            }

            StringBuilder sb = new StringBuilder();
            sb.append("# ").append(filename).append("\n\n");

            List<String> sheetNames = getSheetNames(workbook);
            int sheetCount = 0;

            for (int i = 0; i < sheetNames.size(); i++) {
                if (!includeAllSheets && i > 0) {
                    break;
                }

                if (i > 0) {
                    sb.append(sheetSeparator);
                }

                String sheetName = sheetNames.get(i);
                sb.append("## ").append(sheetName).append("\n\n");

                List<List<String>> rows = getSheetRows(workbook, i);
                if (rows.isEmpty()) {
                    sb.append("*Empty sheet*\n");
                    continue;
                }

                sb.append(rowsToMarkdownTable(rows));
                sheetCount++;
            }

            Map<String, Object> metadata = new HashMap<>();
            metadata.put("sheet_count", sheetCount);
            metadata.put("filename", filename);

            return new ConvertResult(true, sb.toString(), "xlsx", metadata);

        } catch (Exception e) {
            return new ConvertResult(false, "Excel conversion failed: " + e.getMessage());
        }
    }

    /**
     * Convert Markdown tables to Excel.
     */
    public byte[] fromMarkdown(String content, String filename) throws Exception {
        List<List<List<String>>> tables = extractMarkdownTables(content);

        Object workbook = createWorkbook();
        if (workbook == null) {
            throw new RuntimeException("Apache POI not available");
        }

        for (int i = 0; i < tables.size(); i++) {
            String sheetName = "Sheet" + (i + 1);
            addSheet(workbook, sheetName, tables.get(i));
        }

        return saveWorkbook(workbook);
    }

    private String rowsToMarkdownTable(List<List<String>> rows) {
        if (rows.isEmpty()) {
            return "";
        }

        StringBuilder sb = new StringBuilder();
        int maxCols = 0;
        for (List<String> row : rows) {
            maxCols = Math.max(maxCols, row.size());
        }

        // Header
        sb.append("|");
        List<String> header = rows.get(0);
        for (int i = 0; i < maxCols; i++) {
            String cell = i < header.size() ? header.get(i) : "";
            sb.append(" ").append(cell).append(" |");
        }
        sb.append("\n");

        // Separator
        sb.append("|");
        for (int i = 0; i < maxCols; i++) {
            sb.append(" --- |");
        }
        sb.append("\n");

        // Data rows
        for (int r = 1; r < rows.size(); r++) {
            List<String> row = rows.get(r);
            sb.append("|");
            for (int i = 0; i < maxCols; i++) {
                String cell = i < row.size() ? row.get(i) : "";
                sb.append(" ").append(cell).append(" |");
            }
            sb.append("\n");
        }

        return sb.toString();
    }

    private List<List<List<String>>> extractMarkdownTables(String content) {
        List<List<List<String>>> tables = new ArrayList<>();
        List<List<String>> currentTable = new ArrayList<>();
        boolean inTable = false;

        for (String line : content.split("\n")) {
            line = line.trim();

            if (line.startsWith("|") && line.endsWith("|")) {
                // Skip separator line
                if (line.contains("---")) {
                    continue;
                }

                String[] cells = line.split("\\|");
                List<String> row = new ArrayList<>();
                for (int i = 1; i < cells.length - 1; i++) {
                    row.add(cells[i].trim());
                }
                currentTable.add(row);
                inTable = true;
            } else if (inTable && !currentTable.isEmpty()) {
                tables.add(currentTable);
                currentTable = new ArrayList<>();
                inTable = false;
            }
        }

        if (!currentTable.isEmpty()) {
            tables.add(currentTable);
        }

        return tables;
    }

    // Placeholder methods for Apache POI integration
    private Object openWorkbook(byte[] data) {
        try {
            Class<?> workbookClass = Class.forName("org.apache.poi.xssf.usermodel.XSSFWorkbook");
            return workbookClass.getConstructor(java.io.InputStream.class)
                    .newInstance(new ByteArrayInputStream(data));
        } catch (Exception e) {
            return null;
        }
    }

    private Object createWorkbook() {
        try {
            Class<?> workbookClass = Class.forName("org.apache.poi.xssf.usermodel.XSSFWorkbook");
            return workbookClass.getConstructor().newInstance();
        } catch (Exception e) {
            return null;
        }
    }

    private List<String> getSheetNames(Object workbook) {
        List<String> names = new ArrayList<>();
        try {
            int count = (int) workbook.getClass().getMethod("getNumberOfSheets").invoke(workbook);
            for (int i = 0; i < count; i++) {
                names.add((String) workbook.getClass().getMethod("getSheetName", int.class).invoke(workbook, i));
            }
        } catch (Exception e) {
            names.add("Sheet1");
        }
        return names;
    }

    private List<List<String>> getSheetRows(Object workbook, int sheetIndex) {
        List<List<String>> rows = new ArrayList<>();
        try {
            Object sheet = workbook.getClass().getMethod("getSheetAt", int.class).invoke(workbook, sheetIndex);
            if (sheet != null) {
                // Get physical number of rows
                int rowCount = (int) sheet.getClass().getMethod("getPhysicalNumberOfRows").invoke(sheet);
                for (int r = 0; r < rowCount; r++) {
                    Object row = sheet.getClass().getMethod("getRow", int.class).invoke(sheet, r);
                    if (row != null) {
                        List<String> cellValues = new ArrayList<>();
                        int cellCount = (int) row.getClass().getMethod("getPhysicalNumberOfCells").invoke(row);
                        for (int c = 0; c < cellCount; c++) {
                            Object cell = row.getClass().getMethod("getCell", int.class).invoke(row, c);
                            String value = cell != null ? cell.toString() : "";
                            cellValues.add(value);
                        }
                        rows.add(cellValues);
                    }
                }
            }
        } catch (Exception e) {
            // Return empty on any reflection error
        }
        return rows;
    }

    private void addSheet(Object workbook, String name, List<List<String>> data) {
        try {
            Object sheet = workbook.getClass().getMethod("createSheet", String.class).invoke(workbook, name);
            if (sheet != null && data != null) {
                for (int r = 0; r < data.size(); r++) {
                    Object row = sheet.getClass().getMethod("createRow", int.class).invoke(sheet, r);
                    List<String> rowData = data.get(r);
                    for (int c = 0; c < rowData.size(); c++) {
                        Object cell = row.getClass().getMethod("createCell", int.class).invoke(row, c);
                        cell.getClass().getMethod("setCellValue", String.class).invoke(cell, rowData.get(c));
                    }
                }
            }
        } catch (Exception e) {
            // Ignore reflection errors
        }
    }

    private byte[] saveWorkbook(Object workbook) throws Exception {
        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        workbook.getClass().getMethod("write", java.io.OutputStream.class).invoke(workbook, baos);
        return baos.toByteArray();
    }

    // Getters and setters
    public boolean isIncludeAllSheets() {
        return includeAllSheets;
    }

    public void setIncludeAllSheets(boolean includeAllSheets) {
        this.includeAllSheets = includeAllSheets;
    }

    public String getSheetSeparator() {
        return sheetSeparator;
    }

    public void setSheetSeparator(String sheetSeparator) {
        this.sheetSeparator = sheetSeparator;
    }
}
