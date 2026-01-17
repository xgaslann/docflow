"""Batch processor for multi-file queue processing."""

import asyncio
import uuid
from concurrent.futures import ThreadPoolExecutor, ProcessPoolExecutor
from datetime import datetime
from pathlib import Path
from typing import Any, Callable, Dict, List, Optional, Union
from queue import Queue
import threading

from .types import (
    BatchConfig,
    BatchJob,
    JobStatus,
    RAGConfig,
    RAGDocument,
    OutputFormat,
)


class BatchProcessor:
    """Batch processor with queue and parallel processing support.
    
    Supports mixed format inputs (PDF, DOCX, Excel, etc.) in the same batch.
    Provides queue management and parallel processing capabilities.
    
    Example:
        >>> config = RAGConfig(
        ...     llm_processing=[LLMProcessingMode.ALL],
        ...     output_format=OutputFormat.PDF
        ... )
        >>> batch = BatchProcessor(config, max_workers=4)
        >>> 
        >>> # Sync processing
        >>> results = batch.process_files([
        ...     "report.pdf", "data.xlsx", "notes.docx"
        ... ])
        >>> 
        >>> # Async processing with queue
        >>> job_id = await batch.enqueue(files)
        >>> status = await batch.get_status(job_id)
        >>> results = await batch.get_result(job_id)
    """
    
    def __init__(
        self,
        rag_config: Optional[RAGConfig] = None,
        batch_config: Optional[BatchConfig] = None,
        max_workers: int = 4,
    ) -> None:
        """Initialize batch processor.
        
        Args:
            rag_config: RAG configuration for processing.
            batch_config: Batch processing configuration.
            max_workers: Number of parallel workers.
        """
        self.rag_config = rag_config or RAGConfig()
        self.batch_config = batch_config or BatchConfig(max_workers=max_workers)
        self.max_workers = self.batch_config.max_workers
        
        # Job storage
        self._jobs: Dict[str, BatchJob] = {}
        self._job_lock = threading.Lock()
        
        # Queue for async processing
        self._queue: Queue = Queue(maxsize=self.batch_config.queue_size)
        self._executor: Optional[ThreadPoolExecutor] = None
    
    # ============== Sync Processing ==============
    
    def process_files(
        self,
        files: List[Union[str, Path, bytes, tuple]],
        parallel: bool = True,
    ) -> List[RAGDocument]:
        """Process multiple files synchronously.
        
        Args:
            files: List of file paths, bytes, or (bytes, filename) tuples.
            parallel: Use parallel processing.
            
        Returns:
            List of RAGDocument results.
        """
        from .rag import RAGProcessor
        
        processor = RAGProcessor(self.rag_config)
        
        if parallel and len(files) > 1:
            return self._process_parallel(files, processor)
        else:
            return self._process_sequential(files, processor)
    
    def process_to_pdf(
        self,
        files: List[Union[str, Path, bytes, tuple]],
        output_path: Optional[str] = None,
        merge: bool = False,
    ) -> Union[bytes, List[bytes]]:
        """Process files and output as PDF.
        
        Args:
            files: Input files.
            output_path: Optional output path for merged PDF.
            merge: Whether to merge all into one PDF.
            
        Returns:
            PDF bytes or list of PDF bytes.
        """
        # Set output format to PDF
        config = RAGConfig(
            **{**self.rag_config.__dict__, "output_format": OutputFormat.PDF}
        )
        
        from .rag import RAGProcessor
        processor = RAGProcessor(config)
        
        results = self._process_parallel(files, processor)
        
        if merge:
            return self._merge_pdfs([r.pdf_bytes for r in results if r.pdf_bytes])
        else:
            return [r.pdf_bytes for r in results if r.pdf_bytes]
    
    def _process_parallel(
        self,
        files: List[Union[str, Path, bytes, tuple]],
        processor: Any,
    ) -> List[RAGDocument]:
        """Process files in parallel using ThreadPoolExecutor."""
        results = []
        errors = []
        
        with ThreadPoolExecutor(max_workers=self.max_workers) as executor:
            futures = {
                executor.submit(self._process_single, f, processor): f
                for f in files
            }
            
            for future in futures:
                try:
                    result = future.result(timeout=self.batch_config.timeout_per_file)
                    results.append(result)
                except Exception as e:
                    if self.batch_config.fail_fast:
                        raise
                    errors.append((futures[future], str(e)))
        
        return results
    
    def _process_sequential(
        self,
        files: List[Union[str, Path, bytes, tuple]],
        processor: Any,
    ) -> List[RAGDocument]:
        """Process files sequentially."""
        results = []
        
        for f in files:
            try:
                result = self._process_single(f, processor)
                results.append(result)
            except Exception as e:
                if self.batch_config.fail_fast:
                    raise
        
        return results
    
    def _process_single(
        self,
        file_input: Union[str, Path, bytes, tuple],
        processor: Any,
    ) -> RAGDocument:
        """Process a single file."""
        if isinstance(file_input, (str, Path)):
            return processor.process_file(str(file_input))
        elif isinstance(file_input, tuple):
            data, filename = file_input
            return processor.process(data, filename)
        elif isinstance(file_input, bytes):
            return processor.process(file_input, "unknown.bin")
        else:
            raise ValueError(f"Unsupported input type: {type(file_input)}")
    
    # ============== Async Queue Processing ==============
    
    async def enqueue(
        self,
        files: List[Union[str, Path, bytes, tuple]],
    ) -> str:
        """Add files to processing queue.
        
        Args:
            files: Files to process.
            
        Returns:
            Job ID for tracking.
        """
        job_id = str(uuid.uuid4())
        
        job = BatchJob(
            job_id=job_id,
            status=JobStatus.PENDING,
            total_files=len(files),
            created_at=datetime.now().isoformat(),
        )
        
        with self._job_lock:
            self._jobs[job_id] = job
        
        # Start processing in background
        asyncio.create_task(self._process_queue_job(job_id, files))
        
        return job_id
    
    async def get_status(self, job_id: str) -> BatchJob:
        """Get status of a job.
        
        Args:
            job_id: Job ID.
            
        Returns:
            BatchJob with current status.
        """
        with self._job_lock:
            job = self._jobs.get(job_id)
            if not job:
                raise ValueError(f"Job not found: {job_id}")
            return job
    
    async def get_result(self, job_id: str, wait: bool = True) -> List[RAGDocument]:
        """Get results of a completed job.
        
        Args:
            job_id: Job ID.
            wait: Wait for job to complete.
            
        Returns:
            List of RAGDocument results.
        """
        while wait:
            job = await self.get_status(job_id)
            
            if job.status == JobStatus.COMPLETED:
                return job.results
            elif job.status == JobStatus.FAILED:
                raise RuntimeError(f"Job failed: {job.errors}")
            
            await asyncio.sleep(0.5)
        
        job = await self.get_status(job_id)
        return job.results
    
    async def cancel(self, job_id: str) -> bool:
        """Cancel a pending/processing job.
        
        Args:
            job_id: Job ID.
            
        Returns:
            True if cancelled.
        """
        with self._job_lock:
            job = self._jobs.get(job_id)
            if not job:
                return False
            
            if job.status in [JobStatus.PENDING, JobStatus.PROCESSING]:
                job.status = JobStatus.FAILED
                job.errors["cancelled"] = "Job cancelled by user"
                return True
        
        return False
    
    async def _process_queue_job(
        self,
        job_id: str,
        files: List[Union[str, Path, bytes, tuple]],
    ) -> None:
        """Process a queued job."""
        from .rag import RAGProcessor
        
        with self._job_lock:
            job = self._jobs[job_id]
            job.status = JobStatus.PROCESSING
        
        processor = RAGProcessor(self.rag_config)
        loop = asyncio.get_event_loop()
        
        try:
            with ThreadPoolExecutor(max_workers=self.max_workers) as executor:
                futures = []
                
                for f in files:
                    future = loop.run_in_executor(
                        executor,
                        self._process_single,
                        f,
                        processor,
                    )
                    futures.append((f, future))
                
                for file_input, future in futures:
                    try:
                        result = await asyncio.wait_for(
                            future,
                            timeout=self.batch_config.timeout_per_file,
                        )
                        
                        with self._job_lock:
                            job.results.append(result)
                            job.processed_files += 1
                    
                    except Exception as e:
                        filename = str(file_input) if isinstance(file_input, (str, Path)) else "unknown"
                        
                        with self._job_lock:
                            job.errors[filename] = str(e)
                            job.failed_files += 1
                        
                        if self.batch_config.fail_fast:
                            raise
            
            with self._job_lock:
                job.status = JobStatus.COMPLETED
                job.completed_at = datetime.now().isoformat()
        
        except Exception as e:
            with self._job_lock:
                job.status = JobStatus.FAILED
                job.errors["fatal"] = str(e)
                job.completed_at = datetime.now().isoformat()
    
    # ============== PDF Merge ==============
    
    def _merge_pdfs(self, pdf_list: List[bytes]) -> bytes:
        """Merge multiple PDFs into one."""
        try:
            from pypdf import PdfMerger
        except ImportError:
            try:
                from PyPDF2 import PdfMerger
            except ImportError:
                raise ImportError("pypdf or PyPDF2 is required for PDF merging")
        
        from io import BytesIO
        
        merger = PdfMerger()
        
        for pdf_bytes in pdf_list:
            if pdf_bytes:
                merger.append(BytesIO(pdf_bytes))
        
        output = BytesIO()
        merger.write(output)
        merger.close()
        
        return output.getvalue()
    
    # ============== Utility Methods ==============
    
    def list_jobs(self) -> List[BatchJob]:
        """List all jobs."""
        with self._job_lock:
            return list(self._jobs.values())
    
    def clear_completed(self) -> int:
        """Clear completed jobs from memory.
        
        Returns:
            Number of jobs cleared.
        """
        with self._job_lock:
            completed = [
                jid for jid, job in self._jobs.items()
                if job.status in [JobStatus.COMPLETED, JobStatus.FAILED]
            ]
            
            for jid in completed:
                del self._jobs[jid]
            
            return len(completed)
