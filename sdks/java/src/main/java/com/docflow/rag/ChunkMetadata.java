package com.docflow.rag;

import java.util.*;

/**
 * Metadata for a chunk.
 */
public class ChunkMetadata {
    private String sectionTitle = "";
    private List<String> headingPath = new ArrayList<>();
    private boolean hasTable = false;
    private boolean hasImage = false;
    private int page = 0;

    public ChunkMetadata() {
    }

    public ChunkMetadata(String sectionTitle, List<String> headingPath, boolean hasTable, boolean hasImage) {
        this.sectionTitle = sectionTitle;
        this.headingPath = headingPath;
        this.hasTable = hasTable;
        this.hasImage = hasImage;
    }

    public String getSectionTitle() {
        return sectionTitle;
    }

    public void setSectionTitle(String sectionTitle) {
        this.sectionTitle = sectionTitle;
    }

    public List<String> getHeadingPath() {
        return headingPath;
    }

    public void setHeadingPath(List<String> headingPath) {
        this.headingPath = headingPath;
    }

    public boolean hasTable() {
        return hasTable;
    }

    public void setHasTable(boolean hasTable) {
        this.hasTable = hasTable;
    }

    public boolean hasImage() {
        return hasImage;
    }

    public void setHasImage(boolean hasImage) {
        this.hasImage = hasImage;
    }

    public int getPage() {
        return page;
    }

    public void setPage(int page) {
        this.page = page;
    }
}
