package com.docflow.formats;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Converts DOCX files to/from Markdown.
 * 
 * Requires Apache POI dependency:
 * <dependency>
 * <groupId>org.apache.poi</groupId>
 * <artifactId>poi-ooxml</artifactId>
 * <version>5.2.5</version>
 * <optional>true</optional>
 * </dependency>
 */
public class DOCXConverter {

    private boolean extractImages = true;
    private boolean preserveFormatting = true;

    public DOCXConverter() {
    }

    public DOCXConverter(boolean extractImages) {
        this.extractImages = extractImages;
    }

    /**
     * Convert DOCX data to Markdown.
     */
    public ConvertResult toMarkdown(byte[] data, String filename) {
        try {
            Object document = openDocument(data);
            if (document == null) {
                return new ConvertResult(false, "Apache POI not available. Add poi-ooxml dependency.");
            }

            StringBuilder sb = new StringBuilder();
            String title = filename.replaceAll("\\.[^.]+$", "");
            sb.append("# ").append(title).append("\n\n");

            // Extract text content
            String text = extractText(document);
            List<ExtractedImageInfo> images = new ArrayList<>();

            if (extractImages) {
                images = extractImagesFromDoc(document);
            }

            // Convert to markdown
            String markdown = textToMarkdown(text, title);
            sb.append(markdown.substring(markdown.indexOf("\n\n") + 2)); // Skip title we already added

            Map<String, Object> metadata = new HashMap<>();
            metadata.put("filename", filename);
            metadata.put("image_count", images.size());

            // Get document properties
            Map<String, String> props = getDocumentProperties(document);
            if (props.containsKey("title")) {
                metadata.put("title", props.get("title"));
            }
            if (props.containsKey("author")) {
                metadata.put("author", props.get("author"));
            }

            ConvertResult result = new ConvertResult(true, sb.toString(), "docx", metadata);
            return result;

        } catch (Exception e) {
            return new ConvertResult(false, "DOCX conversion failed: " + e.getMessage());
        }
    }

    /**
     * Convert Markdown to DOCX.
     */
    public byte[] fromMarkdown(String content, String filename) throws Exception {
        Object document = createDocument();
        if (document == null) {
            throw new RuntimeException("Apache POI not available");
        }

        String[] lines = content.split("\n");

        for (String line : lines) {
            line = line.trim();

            if (line.startsWith("# ")) {
                addHeading(document, line.substring(2), 1);
            } else if (line.startsWith("## ")) {
                addHeading(document, line.substring(3), 2);
            } else if (line.startsWith("### ")) {
                addHeading(document, line.substring(4), 3);
            } else if (line.startsWith("- ") || line.startsWith("* ")) {
                addListItem(document, line.substring(2));
            } else if (!line.isEmpty()) {
                addParagraph(document, processInlineFormatting(line));
            }
        }

        return saveDocument(document);
    }

    private String textToMarkdown(String text, String title) {
        StringBuilder sb = new StringBuilder();
        sb.append("# ").append(title).append("\n\n");

        String[] lines = text.split("\n");
        List<String> processed = new ArrayList<>();

        for (String line : lines) {
            line = line.trim();

            if (line.isEmpty()) {
                if (!processed.isEmpty() && !processed.get(processed.size() - 1).isEmpty()) {
                    processed.add("");
                }
                continue;
            }

            // Detect headers
            if (isPotentialHeader(line)) {
                line = "## " + toTitleCase(line.toLowerCase());
            }

            // Detect bullet points
            for (String bullet : new String[] { "•", "●", "○", "◦" }) {
                if (line.startsWith(bullet)) {
                    line = "- " + line.substring(bullet.length()).trim();
                    break;
                }
            }

            processed.add(line);
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

            if (line.startsWith("##") || line.startsWith("- ")) {
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
        return sb.toString();
    }

    private boolean isPotentialHeader(String line) {
        return line.length() > 3 && line.length() < 60 &&
                line.equals(line.toUpperCase()) &&
                line.split("\\s+").length >= 2;
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

    private String processInlineFormatting(String text) {
        // Remove markdown bold
        text = text.replaceAll("\\*\\*(.+?)\\*\\*", "$1");
        // Remove markdown italic
        text = text.replaceAll("\\*(.+?)\\*", "$1");
        return text;
    }

    // Placeholder methods for Apache POI integration
    private Object openDocument(byte[] data) {
        try {
            Class<?> docClass = Class.forName("org.apache.poi.xwpf.usermodel.XWPFDocument");
            return docClass.getConstructor(java.io.InputStream.class)
                    .newInstance(new ByteArrayInputStream(data));
        } catch (Exception e) {
            return null;
        }
    }

    private Object createDocument() {
        try {
            Class<?> docClass = Class.forName("org.apache.poi.xwpf.usermodel.XWPFDocument");
            return docClass.getConstructor().newInstance();
        } catch (Exception e) {
            return null;
        }
    }

    private String extractText(Object document) {
        try {
            Class<?> extractorClass = Class.forName("org.apache.poi.xwpf.extractor.XWPFWordExtractor");
            Object extractor = extractorClass.getConstructor(document.getClass()).newInstance(document);
            return (String) extractorClass.getMethod("getText").invoke(extractor);
        } catch (Exception e) {
            return "";
        }
    }

    private List<ExtractedImageInfo> extractImagesFromDoc(Object document) {
        // Implementation would extract images using POI
        return new ArrayList<>();
    }

    private Map<String, String> getDocumentProperties(Object document) {
        Map<String, String> props = new HashMap<>();
        try {
            Object properties = document.getClass().getMethod("getProperties").invoke(document);
            Object coreProps = properties.getClass().getMethod("getCoreProperties").invoke(properties);

            String title = (String) coreProps.getClass().getMethod("getTitle").invoke(coreProps);
            String creator = (String) coreProps.getClass().getMethod("getCreator").invoke(coreProps);

            if (title != null)
                props.put("title", title);
            if (creator != null)
                props.put("author", creator);
        } catch (Exception e) {
            // Ignore
        }
        return props;
    }

    private void addHeading(Object document, String text, int level) {
        try {
            Object para = document.getClass().getMethod("createParagraph").invoke(document);
            para.getClass().getMethod("setStyle", String.class).invoke(para, "Heading" + level);
            Object run = para.getClass().getMethod("createRun").invoke(para);
            run.getClass().getMethod("setText", String.class).invoke(run, text);
        } catch (Exception e) {
            // Ignore
        }
    }

    private void addParagraph(Object document, String text) {
        try {
            Object para = document.getClass().getMethod("createParagraph").invoke(document);
            Object run = para.getClass().getMethod("createRun").invoke(para);
            run.getClass().getMethod("setText", String.class).invoke(run, text);
        } catch (Exception e) {
            // Ignore
        }
    }

    private void addListItem(Object document, String text) {
        addParagraph(document, "• " + text);
    }

    private byte[] saveDocument(Object document) throws Exception {
        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        document.getClass().getMethod("write", java.io.OutputStream.class).invoke(document, baos);
        return baos.toByteArray();
    }

    // Helper class for extracted images
    private static class ExtractedImageInfo {
        byte[] data;
        String format;
        String filename;

        ExtractedImageInfo(byte[] data, String format, String filename) {
            this.data = data;
            this.format = format;
            this.filename = filename;
        }

        public byte[] getData() {
            return data;
        }

        public String getFormat() {
            return format;
        }

        public String getFilename() {
            return filename;
        }
    }

    // Getters and setters
    public boolean isExtractImages() {
        return extractImages;
    }

    public void setExtractImages(boolean extractImages) {
        this.extractImages = extractImages;
    }

    public boolean isPreserveFormatting() {
        return preserveFormatting;
    }

    public void setPreserveFormatting(boolean preserveFormatting) {
        this.preserveFormatting = preserveFormatting;
    }
}
