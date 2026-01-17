package com.docflow.config;

/**
 * Configuration for retrieval operations.
 */
public class RetrievalConfig {

    private int topK = 5;
    private double similarityThreshold = 0.7;
    private double minScore = 0.0;
    private boolean rerank = false;
    private String rerankModel = "cross-encoder/ms-marco-MiniLM-L-6-v2";
    private int rerankTopK = 3;
    private boolean filterDuplicates = true;
    private double duplicateThreshold = 0.95;
    private boolean includeContext = true;
    private int contextBefore = 1;
    private int contextAfter = 1;
    private boolean hybridSearch = false;
    private double keywordWeight = 0.3;
    private double semanticWeight = 0.7;
    private boolean useMMR = false;
    private double mmrLambda = 0.5;

    public RetrievalConfig() {
    }

    public static RetrievalConfig defaultConfig() {
        return new RetrievalConfig();
    }

    public void validate() {
        if (topK <= 0) {
            throw new IllegalArgumentException("topK must be positive");
        }
        if (similarityThreshold < 0 || similarityThreshold > 1) {
            throw new IllegalArgumentException("similarityThreshold must be between 0 and 1");
        }
    }

    // Getters and setters
    public int getTopK() {
        return topK;
    }

    public void setTopK(int topK) {
        this.topK = topK;
    }

    public double getSimilarityThreshold() {
        return similarityThreshold;
    }

    public void setSimilarityThreshold(double similarityThreshold) {
        this.similarityThreshold = similarityThreshold;
    }

    public double getMinScore() {
        return minScore;
    }

    public void setMinScore(double minScore) {
        this.minScore = minScore;
    }

    public boolean isRerank() {
        return rerank;
    }

    public void setRerank(boolean rerank) {
        this.rerank = rerank;
    }

    public String getRerankModel() {
        return rerankModel;
    }

    public void setRerankModel(String rerankModel) {
        this.rerankModel = rerankModel;
    }

    public int getRerankTopK() {
        return rerankTopK;
    }

    public void setRerankTopK(int rerankTopK) {
        this.rerankTopK = rerankTopK;
    }

    public boolean isFilterDuplicates() {
        return filterDuplicates;
    }

    public void setFilterDuplicates(boolean filterDuplicates) {
        this.filterDuplicates = filterDuplicates;
    }

    public double getDuplicateThreshold() {
        return duplicateThreshold;
    }

    public void setDuplicateThreshold(double duplicateThreshold) {
        this.duplicateThreshold = duplicateThreshold;
    }

    public boolean isIncludeContext() {
        return includeContext;
    }

    public void setIncludeContext(boolean includeContext) {
        this.includeContext = includeContext;
    }

    public int getContextBefore() {
        return contextBefore;
    }

    public void setContextBefore(int contextBefore) {
        this.contextBefore = contextBefore;
    }

    public int getContextAfter() {
        return contextAfter;
    }

    public void setContextAfter(int contextAfter) {
        this.contextAfter = contextAfter;
    }

    public boolean isHybridSearch() {
        return hybridSearch;
    }

    public void setHybridSearch(boolean hybridSearch) {
        this.hybridSearch = hybridSearch;
    }

    public double getKeywordWeight() {
        return keywordWeight;
    }

    public void setKeywordWeight(double keywordWeight) {
        this.keywordWeight = keywordWeight;
    }

    public double getSemanticWeight() {
        return semanticWeight;
    }

    public void setSemanticWeight(double semanticWeight) {
        this.semanticWeight = semanticWeight;
    }

    public boolean isUseMMR() {
        return useMMR;
    }

    public void setUseMMR(boolean useMMR) {
        this.useMMR = useMMR;
    }

    public double getMmrLambda() {
        return mmrLambda;
    }

    public void setMmrLambda(double mmrLambda) {
        this.mmrLambda = mmrLambda;
    }
}
