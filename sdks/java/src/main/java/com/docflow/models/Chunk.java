package com.docflow.models;

import java.util.*;

/**
 * Represents a text chunk for RAG.
 */
public class Chunk {
    private int index;
    private String content;
    private ChunkMetadata metadata;

    public Chunk() {
    }

    public Chunk(int index, String content) {
        this.index = index;
        this.content = content;
    }

    public int getIndex() {
        return index;
    }

    public void setIndex(int index) {
        this.index = index;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public ChunkMetadata getMetadata() {
        return metadata;
    }

    public void setMetadata(ChunkMetadata metadata) {
        this.metadata = metadata;
    }
}
