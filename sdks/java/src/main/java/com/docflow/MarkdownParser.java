package com.docflow;

import com.vladsch.flexmark.html.HtmlRenderer;
import com.vladsch.flexmark.parser.Parser;
import com.vladsch.flexmark.util.ast.Node;
import com.vladsch.flexmark.ext.gfm.strikethrough.StrikethroughExtension;
import com.vladsch.flexmark.ext.tables.TablesExtension;
import com.vladsch.flexmark.ext.gfm.tasklist.TaskListExtension;
import com.vladsch.flexmark.util.data.MutableDataSet;

import com.docflow.models.MDFile;

import java.util.*;
import java.util.stream.Collectors;

/**
 * Markdown parser for DocFlow.
 *
 * <p>
 * Uses Flexmark for parsing with GFM extensions.
 */
public class MarkdownParser {

    private final Parser parser;
    private final HtmlRenderer renderer;

    /**
     * Create a new markdown parser with default extensions.
     */
    public MarkdownParser() {
        MutableDataSet options = new MutableDataSet();
        options.set(Parser.EXTENSIONS, Arrays.asList(
                TablesExtension.create(),
                StrikethroughExtension.create(),
                TaskListExtension.create()));

        this.parser = Parser.builder(options).build();
        this.renderer = HtmlRenderer.builder(options).build();
    }

    /**
     * Convert markdown content to HTML.
     *
     * @param content Markdown string.
     * @return HTML string.
     */
    public String toHtml(String content) {
        Node document = parser.parse(content);
        return renderer.render(document);
    }

    /**
     * Merge multiple files into a single content string.
     * Files are sorted by their order field.
     *
     * @param files List of MDFile objects.
     * @return Merged markdown content.
     */
    public String mergeFiles(List<MDFile> files) {
        if (files == null || files.isEmpty()) {
            return "";
        }

        List<MDFile> sorted = files.stream()
                .sorted(Comparator.comparingInt(MDFile::getOrder))
                .collect(Collectors.toList());

        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < sorted.size(); i++) {
            if (i > 0) {
                sb.append("\n\n---\n\n");
            }
            sb.append(sorted.get(i).getContent());
        }

        return sb.toString();
    }

    /**
     * Merge files and convert to HTML with file separators.
     *
     * @param files List of MDFile objects.
     * @return HTML string with file separators.
     */
    public String mergeFilesToHtml(List<MDFile> files) {
        if (files == null || files.isEmpty()) {
            return "";
        }

        List<MDFile> sorted = files.stream()
                .sorted(Comparator.comparingInt(MDFile::getOrder))
                .collect(Collectors.toList());

        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < sorted.size(); i++) {
            MDFile file = sorted.get(i);
            if (i > 0) {
                sb.append("<div class=\"file-separator\"><span>")
                        .append(file.getName())
                        .append("</span></div>");
            } else {
                sb.append("<div class=\"file-header\"><span>")
                        .append(file.getName())
                        .append("</span></div>");
            }

            sb.append("<div class=\"file-content\">")
                    .append(toHtml(file.getContent()))
                    .append("</div>");
        }

        return sb.toString();
    }

    /**
     * Estimate the number of PDF pages based on content.
     *
     * @param content Markdown or text content.
     * @return Estimated page count (minimum 1).
     */
    public int estimatePageCount(String content) {
        int charsPerPage = 3000;
        int pages = content.length() / charsPerPage;
        return Math.max(1, pages);
    }
}
