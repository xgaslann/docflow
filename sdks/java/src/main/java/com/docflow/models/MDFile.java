package com.docflow.models;

/**
 * Represents a Markdown file to be converted.
 */
public class MDFile {
    private String id;
    private String name;
    private String content;
    private int order;

    public MDFile() {}

    public MDFile(String name, String content) {
        this.id = name;
        this.name = name;
        this.content = content;
        this.order = 0;
    }

    public MDFile(String id, String name, String content, int order) {
        this.id = id;
        this.name = name;
        this.content = content;
        this.order = order;
    }

    // Getters and Setters
    public String getId() { return id; }
    public void setId(String id) { this.id = id; }

    public String getName() { return name; }
    public void setName(String name) { this.name = name; }

    public String getContent() { return content; }
    public void setContent(String content) { this.content = content; }

    public int getOrder() { return order; }
    public void setOrder(int order) { this.order = order; }

    public static Builder builder() {
        return new Builder();
    }

    public static class Builder {
        private String id;
        private String name;
        private String content;
        private int order;

        public Builder id(String id) { this.id = id; return this; }
        public Builder name(String name) { this.name = name; return this; }
        public Builder content(String content) { this.content = content; return this; }
        public Builder order(int order) { this.order = order; return this; }

        public MDFile build() {
            MDFile file = new MDFile();
            file.id = this.id != null ? this.id : this.name;
            file.name = this.name;
            file.content = this.content;
            file.order = this.order;
            return file;
        }
    }
}
