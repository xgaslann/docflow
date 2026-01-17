package com.docflow.models;

import java.util.*;

/**
 * Metadata for a chunk.
 */
public class ChunkMetadata {
    private List<String> headingPath = new ArrayList<>();
    private String sectionTitle;
    private boolean hasCode;
    private boolean hasTable;
    private int startPos;
    private int endPos;

    public ChunkMetadata() {
    }

    public List<String> getHeadingPath() {
        return headingPath;
    }

    public void setHeadingPath(List<String> headingPath) {
        this.headingPath = headingPath;
    }

    public String getSectionTitle() {
        return sectionTitle;
    }

    public void setSectionTitle(String sectionTitle) {
        this.sectionTitle = sectionTitle;
    }

    public boolean isHasCode() {
        return hasCode;
    }

    public void setHasCode(boolean hasCode) {
        this.hasCode = hasCode;
    }

    public boolean isHasTable() {
        return hasTable;
    }

    public void setHasTable(boolean hasTable) {
        this.hasTable = hasTable;
    }

    public int getStartPos() {
        return startPos;
    }

    public void setStartPos(int startPos) {
        this.startPos = startPos;
    }

    public int getEndPos() {
        return endPos;
    }

    public void setEndPos(int endPos) {
        this.endPos = endPos;
    }
}
