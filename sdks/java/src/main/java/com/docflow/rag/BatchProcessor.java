package com.docflow.rag;

import com.docflow.models.RAGDocument;

import java.util.*;
import java.util.concurrent.*;
import java.util.function.Consumer;

/**
 * BatchProcessor for processing multiple documents in parallel.
 * Provides job queuing, status tracking, and result collection.
 */
public class BatchProcessor {

    private final RAGConfig ragConfig;
    private final BatchConfig batchConfig;
    private final ExecutorService executor;
    private final Map<String, BatchJob> jobs;
    private final RAGProcessor ragProcessor;

    public BatchProcessor(RAGConfig ragConfig, BatchConfig batchConfig) {
        this.ragConfig = ragConfig;
        this.batchConfig = batchConfig;
        this.executor = Executors.newFixedThreadPool(batchConfig.getMaxWorkers());
        this.jobs = new ConcurrentHashMap<>();
        this.ragProcessor = new RAGProcessor(ragConfig);
    }

    /**
     * Enqueue a list of files for processing.
     *
     * @param filePaths List of file paths to process
     * @return Job ID to track progress
     */
    public String enqueue(List<String> filePaths) {
        String jobId = UUID.randomUUID().toString();
        BatchJob job = new BatchJob(jobId, filePaths.size());
        jobs.put(jobId, job);

        for (String filePath : filePaths) {
            executor.submit(() -> processFile(job, filePath));
        }

        return jobId;
    }

    /**
     * Get the status of a batch job.
     *
     * @param jobId The job ID
     * @return Job status
     */
    public BatchStatus getStatus(String jobId) {
        BatchJob job = jobs.get(jobId);
        if (job == null) {
            return null;
        }
        return job.getStatus();
    }

    /**
     * Get the results of a completed batch job.
     *
     * @param jobId The job ID
     * @return List of processed documents
     */
    public List<RAGDocument> getResult(String jobId) {
        BatchJob job = jobs.get(jobId);
        if (job == null) {
            return Collections.emptyList();
        }
        return job.getResults();
    }

    /**
     * Cancel a running batch job.
     *
     * @param jobId The job ID
     */
    public void cancel(String jobId) {
        BatchJob job = jobs.get(jobId);
        if (job != null) {
            job.cancel();
        }
    }

    /**
     * Shutdown the batch processor.
     */
    public void shutdown() {
        executor.shutdown();
        try {
            if (!executor.awaitTermination(60, TimeUnit.SECONDS)) {
                executor.shutdownNow();
            }
        } catch (InterruptedException e) {
            executor.shutdownNow();
            Thread.currentThread().interrupt();
        }
    }

    private void processFile(BatchJob job, String filePath) {
        if (job.isCancelled()) {
            return;
        }

        int retries = 0;
        Exception lastError = null;

        while (retries <= batchConfig.getMaxRetries()) {
            try {
                RAGDocument doc = ragProcessor.processFile(filePath);
                job.addResult(doc);

                if (batchConfig.getCallback() != null) {
                    batchConfig.getCallback().accept(doc);
                }
                return;
            } catch (Exception e) {
                lastError = e;
                retries++;

                if (!batchConfig.isRetryFailed() || retries > batchConfig.getMaxRetries()) {
                    break;
                }

                try {
                    Thread.sleep(1000L * retries); // Exponential backoff
                } catch (InterruptedException ie) {
                    Thread.currentThread().interrupt();
                    break;
                }
            }
        }

        // Record failure
        job.addError(filePath, lastError);

        if (batchConfig.isFailFast()) {
            job.cancel();
        }
    }

    // Use ragConfig if needed for validation
    public RAGConfig getRagConfig() {
        return ragConfig;
    }

    /**
     * Configuration for batch processing.
     */
    public static class BatchConfig {
        private int maxWorkers = 8;
        private int queueSize = 1000;
        private boolean retryFailed = true;
        private int maxRetries = 3;
        private int timeoutPerFile = 300; // seconds
        private boolean failFast = false;
        private Consumer<RAGDocument> callback;

        public BatchConfig() {
        }

        public int getMaxWorkers() {
            return maxWorkers;
        }

        public void setMaxWorkers(int maxWorkers) {
            this.maxWorkers = maxWorkers;
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

        public int getTimeoutPerFile() {
            return timeoutPerFile;
        }

        public void setTimeoutPerFile(int timeoutPerFile) {
            this.timeoutPerFile = timeoutPerFile;
        }

        public boolean isFailFast() {
            return failFast;
        }

        public void setFailFast(boolean failFast) {
            this.failFast = failFast;
        }

        public Consumer<RAGDocument> getCallback() {
            return callback;
        }

        public void setCallback(Consumer<RAGDocument> callback) {
            this.callback = callback;
        }
    }

    /**
     * Status of a batch job.
     */
    public static class BatchStatus {
        private final String jobId;
        private final int totalFiles;
        private final int processedFiles;
        private final int failedFiles;
        private final String status; // PENDING, PROCESSING, COMPLETED, FAILED, CANCELLED

        public BatchStatus(String jobId, int totalFiles, int processedFiles, int failedFiles, String status) {
            this.jobId = jobId;
            this.totalFiles = totalFiles;
            this.processedFiles = processedFiles;
            this.failedFiles = failedFiles;
            this.status = status;
        }

        public String getJobId() {
            return jobId;
        }

        public int getTotalFiles() {
            return totalFiles;
        }

        public int getProcessedFiles() {
            return processedFiles;
        }

        public int getFailedFiles() {
            return failedFiles;
        }

        public String getStatus() {
            return status;
        }
    }

    /**
     * Internal job tracking.
     */
    private static class BatchJob {
        private final String jobId;
        private final int totalFiles;
        private final List<RAGDocument> results;
        private final Map<String, Exception> errors;
        private volatile boolean cancelled;

        public BatchJob(String jobId, int totalFiles) {
            this.jobId = jobId;
            this.totalFiles = totalFiles;
            this.results = Collections.synchronizedList(new ArrayList<>());
            this.errors = new ConcurrentHashMap<>();
            this.cancelled = false;
        }

        public void addResult(RAGDocument doc) {
            results.add(doc);
        }

        public void addError(String filePath, Exception error) {
            errors.put(filePath, error);
        }

        public void cancel() {
            this.cancelled = true;
        }

        public boolean isCancelled() {
            return cancelled;
        }

        public List<RAGDocument> getResults() {
            return new ArrayList<>(results);
        }

        public BatchStatus getStatus() {
            int processed = results.size();
            int failed = errors.size();
            int total = totalFiles;

            String status;
            if (cancelled) {
                status = "CANCELLED";
            } else if (processed + failed >= total) {
                status = failed > 0 ? "COMPLETED_WITH_ERRORS" : "COMPLETED";
            } else if (processed + failed > 0) {
                status = "PROCESSING";
            } else {
                status = "PENDING";
            }

            return new BatchStatus(jobId, total, processed, failed, status);
        }
    }
}
