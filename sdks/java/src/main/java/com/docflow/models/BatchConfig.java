package com.docflow.models;

/**
 * Configuration for batch processing.
 */
public class BatchConfig {
    private int maxWorkers = 4;
    private boolean failFast = false;
    private boolean continueOnError = true;
    private int timeoutPerFile = 300;
    private int queueSize = 100;
    private boolean retryFailed = true;
    private int maxRetries = 3;

    public BatchConfig() {
    }

    /**
     * Create batch config with specified workers.
     */
    public static BatchConfig withWorkers(int maxWorkers) {
        BatchConfig config = new BatchConfig();
        config.maxWorkers = maxWorkers;
        return config;
    }

    // Getters and Setters
    public int getMaxWorkers() {
        return maxWorkers;
    }

    public void setMaxWorkers(int maxWorkers) {
        this.maxWorkers = maxWorkers;
    }

    public boolean isFailFast() {
        return failFast;
    }

    public void setFailFast(boolean failFast) {
        this.failFast = failFast;
    }

    public boolean isContinueOnError() {
        return continueOnError;
    }

    public void setContinueOnError(boolean continueOnError) {
        this.continueOnError = continueOnError;
    }

    public int getTimeoutPerFile() {
        return timeoutPerFile;
    }

    public void setTimeoutPerFile(int timeoutPerFile) {
        this.timeoutPerFile = timeoutPerFile;
    }

    public int getQueueSize() {
        return queueSize;
    }

    public void setQueueSize(int queueSize) {
        this.queueSize = queueSize;
    }

    public boolean isRetryFailed() {
        return retryFailed;
    }

    public void setRetryFailed(boolean retryFailed) {
        this.retryFailed = retryFailed;
    }

    public int getMaxRetries() {
        return maxRetries;
    }

    public void setMaxRetries(int maxRetries) {
        this.maxRetries = maxRetries;
    }
}
