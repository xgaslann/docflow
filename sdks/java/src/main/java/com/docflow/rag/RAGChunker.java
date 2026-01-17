package com.docflow.rag;

import java.util.*;
import java.util.regex.*;

/**
 * Smart chunker for RAG-optimized markdown content.
 */
public class RAGChunker {

    private final RAGConfig config;

    public RAGChunker(RAGConfig config) {
        this.config = config;
    }

    /**
     * Chunk markdown content into RAG-optimized chunks.
     */
    public List<Chunk> chunk(String markdown) {
        String content = extractContent(markdown);

        List<Section> sections;
        if (config.isRespectHeadings()) {
            sections = splitByHeadings(content);
        } else {
            sections = List.of(new Section("", content));
        }

        List<Chunk> chunks = new ArrayList<>();
        int chunkIndex = 0;
        int charOffset = 0;

        for (Section section : sections) {
            List<Chunk> sectionChunks = chunkSection(section.content, section.title, chunkIndex, charOffset);
            chunks.addAll(sectionChunks);
            chunkIndex += sectionChunks.size();
            charOffset += section.content.length();
        }

        if (config.getChunkOverlap() > 0) {
            addOverlap(chunks);
        }

        if (config.isAddChunkMarkers()) {
            for (int i = 0; i < chunks.size(); i++) {
                Chunk chunk = chunks.get(i);
                chunk.setContent(chunk.getContent() + "\n\n<!-- chunk_boundary: " + i + " -->");
            }
        }

        return chunks;
    }

    private String extractContent(String markdown) {
        if (markdown.startsWith("---")) {
            int end = markdown.indexOf("---", 3);
            if (end > 0) {
                return markdown.substring(end + 3).trim();
            }
        }
        return markdown;
    }

    private List<Section> splitByHeadings(String content) {
        Pattern pattern = Pattern.compile("(?m)^(#{1,6})\\s+(.+)$");
        Matcher matcher = pattern.matcher(content);

        List<Section> sections = new ArrayList<>();
        int lastEnd = 0;
        String lastTitle = "";

        List<int[]> matches = new ArrayList<>();
        List<String> titles = new ArrayList<>();

        while (matcher.find()) {
            matches.add(new int[] { matcher.start(), matcher.end() });
            titles.add(matcher.group(2).trim());
        }

        for (int i = 0; i < matches.size(); i++) {
            int start = matches.get(i)[0];
            if (start > lastEnd) {
                String sectionContent = content.substring(lastEnd, start).trim();
                if (!sectionContent.isEmpty()) {
                    sections.add(new Section(lastTitle, sectionContent));
                }
            }
            lastTitle = titles.get(i);
            lastEnd = matches.get(i)[0];
        }

        if (lastEnd < content.length()) {
            String remaining = content.substring(lastEnd).trim();
            if (!remaining.isEmpty()) {
                sections.add(new Section(lastTitle, remaining));
            }
        }

        if (sections.isEmpty()) {
            sections.add(new Section("", content));
        }

        return sections;
    }

    private List<Chunk> chunkSection(String content, String sectionTitle, int startIndex, int charOffset) {
        List<Chunk> chunks = new ArrayList<>();
        String[] lines = content.split("\n");
        StringBuilder currentChunk = new StringBuilder();
        int currentStart = charOffset;
        int chunkIdx = startIndex;

        int i = 0;
        while (i < lines.length) {
            String line = lines[i];

            // Handle code blocks
            if (line.trim().startsWith("```")) {
                StringBuilder block = new StringBuilder(line).append("\n");
                i++;
                while (i < lines.length && !lines[i].trim().startsWith("```")) {
                    block.append(lines[i]).append("\n");
                    i++;
                }
                if (i < lines.length) {
                    block.append(lines[i]);
                }

                if (currentChunk.length() + block.length() > config.getChunkSize() && currentChunk.length() > 0) {
                    chunks.add(createChunk(currentChunk.toString().trim(), chunkIdx++, currentStart,
                            currentStart + currentChunk.length(), sectionTitle));
                    currentChunk = new StringBuilder();
                    currentStart = charOffset + countChars(Arrays.copyOfRange(lines, 0, i));
                }
                currentChunk.append(block).append("\n");
                i++;
                continue;
            }

            // Handle tables
            if (line.trim().startsWith("|") && line.trim().endsWith("|")) {
                StringBuilder table = new StringBuilder(line).append("\n");
                i++;
                while (i < lines.length && lines[i].trim().startsWith("|")) {
                    table.append(lines[i]).append("\n");
                    i++;
                }

                if (currentChunk.length() + table.length() > config.getChunkSize() && currentChunk.length() > 0) {
                    chunks.add(createChunk(currentChunk.toString().trim(), chunkIdx++, currentStart,
                            currentStart + currentChunk.length(), sectionTitle));
                    currentChunk = new StringBuilder();
                }
                currentChunk.append(table);
                continue;
            }

            // Regular line
            if (currentChunk.length() + line.length() > config.getChunkSize() && currentChunk.length() > 0) {
                chunks.add(createChunk(currentChunk.toString().trim(), chunkIdx++, currentStart,
                        currentStart + currentChunk.length(), sectionTitle));
                currentChunk = new StringBuilder();
                currentStart = charOffset + countChars(Arrays.copyOfRange(lines, 0, i));
            }
            currentChunk.append(line).append("\n");
            i++;
        }

        if (currentChunk.toString().trim().length() > 0) {
            chunks.add(createChunk(currentChunk.toString().trim(), chunkIdx, currentStart,
                    currentStart + currentChunk.length(), sectionTitle));
        }

        return chunks;
    }

    private Chunk createChunk(String content, int index, int startChar, int endChar, String sectionTitle) {
        boolean hasTable = content.contains("|") && content.contains("---");
        boolean hasImage = content.contains("![") || content.contains("[Image:");

        ChunkMetadata metadata = new ChunkMetadata(sectionTitle, extractHeadingPath(content), hasTable, hasImage);
        return new Chunk(content, index, startChar, endChar, metadata);
    }

    private List<String> extractHeadingPath(String content) {
        List<String> headings = new ArrayList<>();
        for (String line : content.split("\n")) {
            if (line.startsWith("#")) {
                headings.add(line.replaceFirst("^#+\\s*", "").trim());
            }
        }
        return headings;
    }

    private void addOverlap(List<Chunk> chunks) {
        if (chunks.size() <= 1)
            return;

        int overlapSize = config.getChunkOverlap();

        for (int i = 1; i < chunks.size(); i++) {
            String prevContent = chunks.get(i - 1).getContent();
            String overlap = prevContent.length() > overlapSize
                    ? prevContent.substring(prevContent.length() - overlapSize)
                    : prevContent;

            // Find break point
            for (String breakStr : new String[] { "\n\n", ". ", "\n" }) {
                int idx = overlap.indexOf(breakStr);
                if (idx > 0) {
                    overlap = overlap.substring(idx + breakStr.length());
                    break;
                }
            }

            if (!overlap.trim().isEmpty()) {
                chunks.get(i).setContent("[...] " + overlap.trim() + "\n\n" + chunks.get(i).getContent());
            }
        }
    }

    private int countChars(String[] lines) {
        int total = 0;
        for (String line : lines) {
            total += line.length() + 1;
        }
        return total;
    }

    private static class Section {
        String title;
        String content;

        Section(String title, String content) {
            this.title = title;
            this.content = content;
        }
    }
}
