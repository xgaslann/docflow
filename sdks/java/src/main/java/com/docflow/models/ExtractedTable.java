package com.docflow.models;

import java.util.*;

/**
 * Represents an extracted table from a document.
 */
public class ExtractedTable {
    private String markdown;
    private String summary;
    private List<List<String>> data = new ArrayList<>();
    private int pageNumber;

    public ExtractedTable() {
    }

    public ExtractedTable(String markdown) {
        this.markdown = markdown;
    }

    public String getMarkdown() {
        return markdown;
    }

    public void setMarkdown(String markdown) {
        this.markdown = markdown;
    }

    public String getSummary() {
        return summary;
    }

    public void setSummary(String summary) {
        this.summary = summary;
    }

    public List<List<String>> getData() {
        return data;
    }

    public void setData(List<List<String>> data) {
        this.data = data;
    }

    public int getPageNumber() {
        return pageNumber;
    }

    public void setPageNumber(int pageNumber) {
        this.pageNumber = pageNumber;
    }
}
