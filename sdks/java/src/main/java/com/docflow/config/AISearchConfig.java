package com.docflow.config;

import java.util.*;

/**
 * Configuration for Azure AI Search.
 */
public class AISearchConfig {

    public enum QueryType {
        SIMPLE("simple"),
        FULL("full"),
        SEMANTIC("semantic");

        private final String value;

        QueryType(String value) {
            this.value = value;
        }

        public String getValue() {
            return value;
        }
    }

    public enum SearchMode {
        ANY("any"),
        ALL("all");

        private final String value;

        SearchMode(String value) {
            this.value = value;
        }

        public String getValue() {
            return value;
        }
    }

    // Connection
    private String endpoint = "";
    private String apiKey = "";
    private String apiVersion = "2024-07-01";
    private String indexName = "docflow-index";

    // Vector configuration
    private String vectorSearchProfile = "default-vector-profile";
    private int hnswM = 4;
    private int hnswEfConstruction = 400;
    private int hnswEfSearch = 500;
    private String metric = "cosine";

    // Semantic configuration
    private String semanticConfig = "default-semantic-config";
    private List<String> semanticPrioritizedFields = Arrays.asList("content", "title");

    // Search options
    private QueryType queryType = QueryType.SEMANTIC;
    private SearchMode searchMode = SearchMode.ANY;
    private int top = 10;
    private int skip = 0;

    // Vector search
    private List<String> vectorFields = Arrays.asList("content_vector");
    private int kNearestNeighbors = 50;

    // Hybrid search
    private boolean hybridSearch = true;
    private int maxTextRecallSize = 1000;

    // Semantic reranking
    private boolean semanticReranking = true;
    private int semanticMaxWait = 700;

    // Filters & Fields
    private String filterExpression = "";
    private List<String> searchFields = new ArrayList<>();
    private List<String> selectFields = Arrays.asList("id", "content", "title", "metadata");
    private List<String> orderBy = new ArrayList<>();

    // Facets
    private List<String> facets = new ArrayList<>();

    // Highlighting
    private List<String> highlightFields = new ArrayList<>();
    private String highlightPreTag = "<em>";
    private String highlightPostTag = "</em>";

    // Embedding
    private String embeddingModel = "text-embedding-3-small";
    private int embeddingDimensions = 1536;
    private String embeddingApiKey = "";

    public AISearchConfig() {
    }

    public static AISearchConfig defaultConfig() {
        return new AISearchConfig();
    }

    public void validate() {
        if (endpoint == null || endpoint.isEmpty()) {
            throw new IllegalArgumentException("Azure AI Search endpoint is required");
        }
        if (apiKey == null || apiKey.isEmpty()) {
            throw new IllegalArgumentException("Azure AI Search API key is required");
        }
    }

    // Getters and setters
    public String getEndpoint() {
        return endpoint;
    }

    public void setEndpoint(String endpoint) {
        this.endpoint = endpoint;
    }

    public String getApiKey() {
        return apiKey;
    }

    public void setApiKey(String apiKey) {
        this.apiKey = apiKey;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    public String getIndexName() {
        return indexName;
    }

    public void setIndexName(String indexName) {
        this.indexName = indexName;
    }

    public String getVectorSearchProfile() {
        return vectorSearchProfile;
    }

    public void setVectorSearchProfile(String vectorSearchProfile) {
        this.vectorSearchProfile = vectorSearchProfile;
    }

    public int getHnswM() {
        return hnswM;
    }

    public void setHnswM(int hnswM) {
        this.hnswM = hnswM;
    }

    public int getHnswEfConstruction() {
        return hnswEfConstruction;
    }

    public void setHnswEfConstruction(int hnswEfConstruction) {
        this.hnswEfConstruction = hnswEfConstruction;
    }

    public int getHnswEfSearch() {
        return hnswEfSearch;
    }

    public void setHnswEfSearch(int hnswEfSearch) {
        this.hnswEfSearch = hnswEfSearch;
    }

    public String getMetric() {
        return metric;
    }

    public void setMetric(String metric) {
        this.metric = metric;
    }

    public String getSemanticConfig() {
        return semanticConfig;
    }

    public void setSemanticConfig(String semanticConfig) {
        this.semanticConfig = semanticConfig;
    }

    public List<String> getSemanticPrioritizedFields() {
        return semanticPrioritizedFields;
    }

    public void setSemanticPrioritizedFields(List<String> semanticPrioritizedFields) {
        this.semanticPrioritizedFields = semanticPrioritizedFields;
    }

    public QueryType getQueryType() {
        return queryType;
    }

    public void setQueryType(QueryType queryType) {
        this.queryType = queryType;
    }

    public SearchMode getSearchMode() {
        return searchMode;
    }

    public void setSearchMode(SearchMode searchMode) {
        this.searchMode = searchMode;
    }

    public int getTop() {
        return top;
    }

    public void setTop(int top) {
        this.top = top;
    }

    public int getSkip() {
        return skip;
    }

    public void setSkip(int skip) {
        this.skip = skip;
    }

    public List<String> getVectorFields() {
        return vectorFields;
    }

    public void setVectorFields(List<String> vectorFields) {
        this.vectorFields = vectorFields;
    }

    public int getkNearestNeighbors() {
        return kNearestNeighbors;
    }

    public void setkNearestNeighbors(int kNearestNeighbors) {
        this.kNearestNeighbors = kNearestNeighbors;
    }

    public boolean isHybridSearch() {
        return hybridSearch;
    }

    public void setHybridSearch(boolean hybridSearch) {
        this.hybridSearch = hybridSearch;
    }

    public boolean isSemanticReranking() {
        return semanticReranking;
    }

    public void setSemanticReranking(boolean semanticReranking) {
        this.semanticReranking = semanticReranking;
    }

    public String getFilterExpression() {
        return filterExpression;
    }

    public void setFilterExpression(String filterExpression) {
        this.filterExpression = filterExpression;
    }

    public List<String> getSearchFields() {
        return searchFields;
    }

    public void setSearchFields(List<String> searchFields) {
        this.searchFields = searchFields;
    }

    public List<String> getSelectFields() {
        return selectFields;
    }

    public void setSelectFields(List<String> selectFields) {
        this.selectFields = selectFields;
    }

    public List<String> getHighlightFields() {
        return highlightFields;
    }

    public void setHighlightFields(List<String> highlightFields) {
        this.highlightFields = highlightFields;
    }

    public String getEmbeddingModel() {
        return embeddingModel;
    }

    public void setEmbeddingModel(String embeddingModel) {
        this.embeddingModel = embeddingModel;
    }

    public int getEmbeddingDimensions() {
        return embeddingDimensions;
    }

    public void setEmbeddingDimensions(int embeddingDimensions) {
        this.embeddingDimensions = embeddingDimensions;
    }

    public String getEmbeddingApiKey() {
        return embeddingApiKey;
    }

    public void setEmbeddingApiKey(String embeddingApiKey) {
        this.embeddingApiKey = embeddingApiKey;
    }

    public int getMaxTextRecallSize() {
        return maxTextRecallSize;
    }

    public void setMaxTextRecallSize(int maxTextRecallSize) {
        this.maxTextRecallSize = maxTextRecallSize;
    }

    public int getSemanticMaxWait() {
        return semanticMaxWait;
    }

    public void setSemanticMaxWait(int semanticMaxWait) {
        this.semanticMaxWait = semanticMaxWait;
    }

    public List<String> getOrderBy() {
        return orderBy;
    }

    public void setOrderBy(List<String> orderBy) {
        this.orderBy = orderBy;
    }

    public List<String> getFacets() {
        return facets;
    }

    public void setFacets(List<String> facets) {
        this.facets = facets;
    }

    public String getHighlightPreTag() {
        return highlightPreTag;
    }

    public void setHighlightPreTag(String highlightPreTag) {
        this.highlightPreTag = highlightPreTag;
    }

    public String getHighlightPostTag() {
        return highlightPostTag;
    }

    public void setHighlightPostTag(String highlightPostTag) {
        this.highlightPostTag = highlightPostTag;
    }
}
