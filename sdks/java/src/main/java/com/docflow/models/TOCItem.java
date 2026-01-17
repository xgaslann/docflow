package com.docflow.models;

import java.util.List;

/**
 * Table of contents item.
 */
public class TOCItem {
    private String title;
    private int level;
    private String anchor;
    private int page;
    private List<TOCItem> children;

    public TOCItem() {
    }

    public TOCItem(String title, int level, String anchor) {
        this.title = title;
        this.level = level;
        this.anchor = anchor;
    }

    // Getters and Setters
    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public int getLevel() {
        return level;
    }

    public void setLevel(int level) {
        this.level = level;
    }

    public String getAnchor() {
        return anchor;
    }

    public void setAnchor(String anchor) {
        this.anchor = anchor;
    }

    public int getPage() {
        return page;
    }

    public void setPage(int page) {
        this.page = page;
    }

    public List<TOCItem> getChildren() {
        return children;
    }

    public void setChildren(List<TOCItem> children) {
        this.children = children;
    }
}
