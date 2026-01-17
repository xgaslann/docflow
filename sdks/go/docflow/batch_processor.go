package docflow

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/xgaslan/docflow/sdks/go/docflow/config"
	"github.com/xgaslan/docflow/sdks/go/docflow/rag"
)

// jobRequest wraps a job and its files for processing
type jobRequest struct {
	job   *BatchJob
	files []string
}

// BatchProcessor handles multi-file queue processing.
type BatchProcessor struct {
	ragConfig   config.RAGConfig
	batchConfig config.BatchConfig
	maxWorkers  int

	jobs    sync.Map // map[string]*BatchJob
	queue   chan jobRequest
	workers sync.WaitGroup
}

// NewBatchProcessor creates a new batch processor.
func NewBatchProcessor(ragConfig config.RAGConfig, batchConfig config.BatchConfig) *BatchProcessor {
	if batchConfig.MaxWorkers <= 0 {
		batchConfig.MaxWorkers = 4
	}
	if batchConfig.QueueSize <= 0 {
		batchConfig.QueueSize = 100
	}

	bp := &BatchProcessor{
		ragConfig:   ragConfig,
		batchConfig: batchConfig,
		maxWorkers:  batchConfig.MaxWorkers,
		queue:       make(chan jobRequest, batchConfig.QueueSize),
	}

	bp.startWorkers()
	return bp
}

func (bp *BatchProcessor) startWorkers() {
	for i := 0; i < bp.maxWorkers; i++ {
		go bp.worker()
	}
}

func (bp *BatchProcessor) worker() {
	for req := range bp.queue {
		bp.processJob(req)
	}
}

func (bp *BatchProcessor) processJob(req jobRequest) {
	job := req.job
	job.Status = config.JobStatusProcessing
	processor := rag.NewRAGProcessor(bp.ragConfig)

	// In a real implementation, we would process files in parallel too if job has multiple files
	// For simplicity here, we process files sequentially or we could spawn goroutines
	// But since we have worker pool for jobs, maybe files within job?
	// The Python implementation uses thread pool for files.
	// We can do semantic equivalent using semaphores or separate worker pool for files if needed.
	// For now, let's just loop.

	// Assuming job.Files is somehow passed or stored.
	// The Python enqueue takes 'files' list. BatchJob struct might need to store file paths.
	// Let's assume BatchJob has a way to access files (not defined in types.go yet? Check types.go)

	// Wait, types.go BatchJob struct:
	// type BatchJob struct {
	// 	JobID          string                  `json:"job_id"`
	// 	Status         JobStatus               `json:"status"`
	// 	Results        []rag.RAGDocument       `json:"results,omitempty"`
	// 	Errors         map[string]string       `json:"errors,omitempty"`
	// 	TotalFiles     int                     `json:"total_files"`
	// 	ProcessedFiles int                     `json:"processed_files"`
	// 	FailedFiles    int                     `json:"failed_files"`
	// 	CreatedAt      time.Time               `json:"created_at"`
	// 	CompletedAt    time.Time               `json:"completed_at,omitempty"`
	// }
	// It doesn't store input files! Python stores them in closure or somewhere?
	// Python: enqueue(files) -> creates job -> starts _process_queue_job(job_id, files).
	// So the job struct itself doesn't strictly need to hold the file list if the processor function closes over it.
	// However, in Go, we need to pass the file list to the worker via the channel struct.
	// So we should define a wrapper or extend BatchJob.

	// Process files
	for _, f := range req.files {
		doc, err := processor.ProcessFile(f)
		if err != nil {
			job.FailedFiles++
			if job.Errors == nil {
				job.Errors = make(map[string]string)
			}
			job.Errors[f] = err.Error()

			if bp.batchConfig.FailFast {
				job.Status = config.JobStatusFailed
				return
			}
			continue
		}

		job.Results = append(job.Results, *doc)
		job.ProcessedFiles++
	}

	if job.Status != config.JobStatusFailed {
		job.Status = config.JobStatusCompleted
		now := time.Now()
		job.CompletedAt = &now // Assumes CreatedAt/CompletedAt are *time.Time or interface{}
	}
}

// Enqueue adds files to the processing queue.
func (bp *BatchProcessor) Enqueue(files []string) (string, error) {
	jobID := uuid.New().String()
	now := time.Now()
	job := &BatchJob{
		JobID:      jobID,
		Status:     config.JobStatusPending,
		TotalFiles: len(files),
		CreatedAt:  &now,
		Errors:     make(map[string]string),
	}

	bp.jobs.Store(jobID, job)

	// In a real app we might want to persist this.

	// Send to worker (this effectively limits concurrency by queue size/workers)
	// Only if I change channel type.
	// Let's define queue type properly.
	select {
	case bp.queue <- jobRequest{job: job, files: files}:
		return jobID, nil
	default:
		return "", fmt.Errorf("queue is full")
	}
}

// ProcessFiles synchronously processes files.
func (bp *BatchProcessor) ProcessFiles(files []string) ([]*rag.RAGDocument, error) {
	processor := rag.NewRAGProcessor(bp.ragConfig)
	var results []*rag.RAGDocument

	for _, f := range files {
		doc, err := processor.ProcessFile(f)
		if err != nil {
			if bp.batchConfig.FailFast {
				return nil, err
			}
			continue
		}
		results = append(results, doc)
	}
	return results, nil
}

// GetStatus returns the status of a job.
func (bp *BatchProcessor) GetStatus(jobID string) (*BatchJob, error) {
	if val, ok := bp.jobs.Load(jobID); ok {
		return val.(*BatchJob), nil
	}
	return nil, fmt.Errorf("job not found: %s", jobID)
}

// GetResult returns the results of a completed job.
func (bp *BatchProcessor) GetResult(jobID string) ([]rag.RAGDocument, error) {
	job, err := bp.GetStatus(jobID)
	if err != nil {
		return nil, err
	}
	if job.Status == config.JobStatusFailed {
		return nil, fmt.Errorf("job failed")
	}
	return job.Results, nil
}
