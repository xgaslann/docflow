package com.docflow.rag;

/**
 * Represents a chunk of content for RAG.
 */
public class Chunk {
    private String content;
    private int index;
    private int startChar;
    private int endChar;
    private ChunkMetadata metadata;

    public Chunk(String content, int index, int startChar, int endChar, ChunkMetadata metadata) {
        this.content = content;
        this.index = index;
        this.startChar = startChar;
        this.endChar = endChar;
        this.metadata = metadata;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public int getIndex() {
        return index;
    }

    public int getStartChar() {
        return startChar;
    }

    public int getEndChar() {
        return endChar;
    }

    public ChunkMetadata getMetadata() {
        return metadata;
    }
}
