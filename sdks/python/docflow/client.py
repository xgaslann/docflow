"""DocFlow API client."""

import base64
from typing import List, Optional, Union

import requests

from .types import (
    MDFile,
    ConvertOptions,
    PDFResult,
    MDResult,
    HealthResponse,
    PreviewResult,
)


class DocFlowClient:
    """
    Client for the DocFlow API.
    
    Provides methods for converting between Markdown and PDF formats.
    
    Example:
        >>> client = DocFlowClient("http://localhost:8080")
        >>> result = client.convert_md_to_pdf([
        ...     MDFile(name="doc.md", content="# Hello World")
        ... ])
        >>> print(result.file_paths)
    """
    
    def __init__(self, base_url: str, timeout: int = 60):
        """
        Initialize the DocFlow client.
        
        Args:
            base_url: Base URL of the DocFlow server (e.g., "http://localhost:8080")
            timeout: Request timeout in seconds (default: 60)
        """
        self.base_url = base_url.rstrip("/")
        self.timeout = timeout
        self.session = requests.Session()
        self.session.headers.update({
            "Content-Type": "application/json",
        })
    
    def health(self) -> HealthResponse:
        """
        Check server health.
        
        Returns:
            HealthResponse with status, version, and timestamp.
            
        Raises:
            requests.RequestException: If the request fails.
        """
        resp = self.session.get(
            f"{self.base_url}/api/health",
            timeout=self.timeout
        )
        resp.raise_for_status()
        data = resp.json()
        return HealthResponse(
            status=data["status"],
            version=data["version"],
            timestamp=data["timestamp"],
        )
    
    def preview(self, content: str) -> PreviewResult:
        """
        Generate HTML preview of markdown content.
        
        Args:
            content: Markdown content to preview.
            
        Returns:
            PreviewResult with HTML content.
        """
        resp = self.session.post(
            f"{self.base_url}/api/preview",
            json={"content": content},
            timeout=self.timeout
        )
        resp.raise_for_status()
        data = resp.json()
        return PreviewResult(html=data["html"])
    
    def convert_md_to_pdf(
        self,
        files: List[Union[MDFile, dict]],
        options: Optional[ConvertOptions] = None,
    ) -> PDFResult:
        """
        Convert Markdown files to PDF.
        
        Args:
            files: List of MDFile objects or dicts with name and content.
            options: Conversion options (merge mode, output name).
            
        Returns:
            PDFResult with success status and file paths.
            
        Example:
            >>> result = client.convert_md_to_pdf([
            ...     MDFile(name="doc.md", content="# Hello")
            ... ], ConvertOptions(merge_mode="merged", output_name="output"))
        """
        if options is None:
            options = ConvertOptions()
        
        # Convert files to dict format
        file_list = []
        for i, f in enumerate(files):
            if isinstance(f, MDFile):
                file_dict = f.to_dict()
                if not file_dict["id"]:
                    file_dict["id"] = str(i + 1)
                file_dict["order"] = i
                file_list.append(file_dict)
            else:
                file_list.append({
                    "id": f.get("id", str(i + 1)),
                    "name": f["name"],
                    "content": f["content"],
                    "order": f.get("order", i),
                })
        
        payload = {
            "files": file_list,
            "mergeMode": options.merge_mode,
        }
        if options.output_name:
            payload["outputName"] = options.output_name
        
        resp = self.session.post(
            f"{self.base_url}/api/convert",
            json=payload,
            timeout=self.timeout
        )
        resp.raise_for_status()
        data = resp.json()
        
        if not data.get("success", False):
            return PDFResult(
                success=False,
                error=data.get("error", "Unknown error"),
            )
        
        return PDFResult(
            success=True,
            file_paths=data.get("files", []),
        )
    
    def extract_pdf_to_md(
        self,
        pdf_data: bytes,
        filename: str,
    ) -> MDResult:
        """
        Extract text from PDF and convert to Markdown.
        
        Args:
            pdf_data: Raw PDF file bytes.
            filename: Original filename of the PDF.
            
        Returns:
            MDResult with success status and markdown content.
            
        Example:
            >>> with open("document.pdf", "rb") as f:
            ...     result = client.extract_pdf_to_md(f.read(), "document.pdf")
            >>> print(result.markdown)
        """
        # Encode PDF to base64
        encoded = base64.b64encode(pdf_data).decode("utf-8")
        
        resp = self.session.post(
            f"{self.base_url}/api/pdf/extract",
            json={
                "fileName": filename,
                "content": encoded,
            },
            timeout=self.timeout
        )
        resp.raise_for_status()
        data = resp.json()
        
        if not data.get("success", False):
            return MDResult(
                success=False,
                error=data.get("error", "Unknown error"),
            )
        
        return MDResult(
            success=True,
            markdown=data.get("markdown", ""),
            file_path=data.get("filePath", ""),
            file_name=data.get("fileName", ""),
        )
    
    def download_pdf(self, file_path: str) -> bytes:
        """
        Download a PDF file from the server.
        
        Args:
            file_path: Path returned from convert_md_to_pdf (e.g., "/output/doc.pdf")
            
        Returns:
            Raw PDF bytes.
        """
        resp = self.session.get(
            f"{self.base_url}{file_path}",
            timeout=self.timeout
        )
        resp.raise_for_status()
        return resp.content
