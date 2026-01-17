package com.docflow.models;

import java.util.List;

/**
 * Information about a document heading.
 */
public class HeadingInfo {
    private String text;
    private int level;
    private int startPos;
    private int endPos;
    private Integer parentIndex;
    private List<Integer> childrenIndices;

    public HeadingInfo() {
    }

    public HeadingInfo(String text, int level, int startPos, int endPos) {
        this.text = text;
        this.level = level;
        this.startPos = startPos;
        this.endPos = endPos;
    }

    // Getters and Setters
    public String getText() {
        return text;
    }

    public void setText(String text) {
        this.text = text;
    }

    public int getLevel() {
        return level;
    }

    public void setLevel(int level) {
        this.level = level;
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

    public Integer getParentIndex() {
        return parentIndex;
    }

    public void setParentIndex(Integer parentIndex) {
        this.parentIndex = parentIndex;
    }

    public List<Integer> getChildrenIndices() {
        return childrenIndices;
    }

    public void setChildrenIndices(List<Integer> childrenIndices) {
        this.childrenIndices = childrenIndices;
    }
}
