package com.docflow.formats;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Converts TXT files to/from Markdown.
 */
public class TXTConverter {

    private boolean detectStructure = true;
    private boolean detectLists = true;
    private int wrapWidth = 80;

    public TXTConverter() {
    }

    public TXTConverter(boolean detectStructure) {
        this.detectStructure = detectStructure;
    }

    /**
     * Convert TXT content to Markdown.
     */
    public ConvertResult toMarkdown(String content, String filename) {
        if (content == null || content.isEmpty()) {
            return new ConvertResult(false, "Empty content");
        }

        StringBuilder sb = new StringBuilder();
        String title = filename.replaceAll("\\.[^.]+$", "");
        sb.append("# ").append(title).append("\n\n");

        String[] lines = content.split("\n");
        List<String> processed = new ArrayList<>();

        for (String line : lines) {
            String trimmed = line.trim();

            if (trimmed.isEmpty()) {
                if (!processed.isEmpty() && !processed.get(processed.size() - 1).isEmpty()) {
                    processed.add("");
                }
                continue;
            }

            // Detect headers
            if (detectStructure && isPotentialHeader(trimmed)) {
                trimmed = "## " + toTitleCase(trimmed.toLowerCase());
            }

            // Detect bullet points
            if (detectLists) {
                for (String bullet : new String[] { "•", "●", "○", "◦", "-", "*" }) {
                    if (trimmed.startsWith(bullet + " ")) {
                        trimmed = "- " + trimmed.substring(bullet.length() + 1).trim();
                        break;
                    }
                }

                // Detect numbered lists
                if (trimmed.matches("^\\d+[.)]\\s+.*")) {
                    trimmed = trimmed.replaceFirst("^\\d+[.)]\\s+", "1. ");
                }
            }

            processed.add(trimmed);
        }

        // Join with proper spacing
        boolean inParagraph = false;
        for (int i = 0; i < processed.size(); i++) {
            String line = processed.get(i);

            if (line.isEmpty()) {
                if (inParagraph) {
                    sb.append("\n\n");
                    inParagraph = false;
                }
                continue;
            }

            if (line.startsWith("##") || line.startsWith("- ") || line.startsWith("1. ")) {
                if (inParagraph) {
                    sb.append("\n\n");
                }
                sb.append(line).append("\n");
                inParagraph = false;
            } else {
                if (inParagraph && i > 0 && !processed.get(i - 1).isEmpty()) {
                    sb.append(" ");
                }
                sb.append(line);
                inParagraph = true;
            }
        }

        sb.append("\n");

        Map<String, Object> metadata = new HashMap<>();
        metadata.put("filename", filename);
        metadata.put("line_count", lines.length);
        metadata.put("word_count", countWords(content));

        return new ConvertResult(true, sb.toString(), "txt", metadata);
    }

    /**
     * Convert Markdown to plain text.
     */
    public String fromMarkdown(String content, String filename) {
        if (content == null) {
            return "";
        }

        StringBuilder sb = new StringBuilder();
        String[] lines = content.split("\n");

        for (String line : lines) {
            // Remove markdown headers
            line = line.replaceFirst("^#{1,6}\\s+", "");

            // Remove bold
            line = line.replaceAll("\\*\\*(.+?)\\*\\*", "$1");

            // Remove italic
            line = line.replaceAll("\\*(.+?)\\*", "$1");

            // Remove inline code
            line = line.replaceAll("`(.+?)`", "$1");

            // Remove links, keep text
            line = line.replaceAll("\\[(.+?)\\]\\(.+?\\)", "$1");

            // Keep list markers but simplify
            line = line.replaceFirst("^-\\s+", "• ");
            line = line.replaceFirst("^\\d+\\.\\s+", "• ");

            sb.append(line).append("\n");
        }

        return sb.toString();
    }

    private boolean isPotentialHeader(String line) {
        return line.length() > 3 && line.length() < 60 &&
                line.equals(line.toUpperCase()) &&
                line.split("\\s+").length >= 2 &&
                !line.contains("•") && !line.contains("-");
    }

    private String toTitleCase(String str) {
        StringBuilder result = new StringBuilder();
        boolean nextUpper = true;

        for (char c : str.toCharArray()) {
            if (Character.isWhitespace(c)) {
                nextUpper = true;
                result.append(c);
            } else if (nextUpper) {
                result.append(Character.toUpperCase(c));
                nextUpper = false;
            } else {
                result.append(c);
            }
        }

        return result.toString();
    }

    private int countWords(String text) {
        if (text == null || text.trim().isEmpty()) {
            return 0;
        }
        return text.trim().split("\\s+").length;
    }

    // Getters and setters
    public boolean isDetectStructure() {
        return detectStructure;
    }

    public void setDetectStructure(boolean detectStructure) {
        this.detectStructure = detectStructure;
    }

    public boolean isDetectLists() {
        return detectLists;
    }

    public void setDetectLists(boolean detectLists) {
        this.detectLists = detectLists;
    }

    public int getWrapWidth() {
        return wrapWidth;
    }

    public void setWrapWidth(int wrapWidth) {
        this.wrapWidth = wrapWidth;
    }
}
