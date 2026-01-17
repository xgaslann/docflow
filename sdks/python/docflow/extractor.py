"""PDF to Markdown extractor for DocFlow."""

import os
import re
import subprocess
import tempfile
import time
from pathlib import Path
from typing import Optional

from .storage.base import Storage
from .types import ExtractResult


class Extractor:
    """Extracts text from PDF files and converts to Markdown.
    
    Uses pdfminer.six for pure Python extraction, with fallback to pdftotext.
    
    Example:
        >>> extractor = Extractor()
        >>> with open("document.pdf", "rb") as f:
        ...     result = extractor.extract_to_markdown(f.read(), "document.pdf")
        >>> print(result.markdown)
    """

    def __init__(
        self,
        storage: Optional[Storage] = None,
        temp_dir: Optional[str] = None,
        use_native: bool = True,
    ) -> None:
        """Initialize the extractor.
        
        Args:
            storage: Storage backend for saving output files.
            temp_dir: Directory for temporary files.
            use_native: Use pdfminer.six (True) or pdftotext CLI (False).
        """
        self.storage = storage
        self.temp_dir = temp_dir or tempfile.gettempdir()
        self.use_native = use_native

    def extract_to_markdown(
        self, pdf_data: bytes, filename: str
    ) -> ExtractResult:
        """Extract text from PDF and convert to Markdown.
        
        Args:
            pdf_data: PDF file bytes.
            filename: Original filename.
            
        Returns:
            ExtractResult with markdown content.
        """
        if not pdf_data:
            return ExtractResult(success=False, error="PDF data is required")

        timestamp = int(time.time())
        base_name = Path(filename).stem
        safe_name = self._sanitize_filename(base_name)

        try:
            # Extract text
            if self.use_native:
                text = self._extract_with_pdfminer(pdf_data)
            else:
                text = self._extract_with_pdftotext(pdf_data, safe_name, timestamp)

            # Get page count
            page_count = self._get_page_count(pdf_data)

            # Convert to markdown
            markdown = self._text_to_markdown(text, base_name)

            # Save if storage configured
            output_path = None
            if self.storage:
                output_path = f"{safe_name}_{timestamp}.md"
                self.storage.save(output_path, markdown.encode("utf-8"))
                output_path = self.storage.get_url(output_path) or output_path

            return ExtractResult(
                success=True,
                markdown=markdown,
                file_path=output_path,
                page_count=page_count,
            )
        except Exception as e:
            return ExtractResult(success=False, error=str(e))

    def extract_from_file(self, path: str) -> ExtractResult:
        """Extract markdown from a PDF file path.
        
        Args:
            path: Path to PDF file.
            
        Returns:
            ExtractResult with markdown content.
        """
        pdf_data = Path(path).read_bytes()
        return self.extract_to_markdown(pdf_data, Path(path).name)

    def preview(self, pdf_data: bytes, filename: str) -> ExtractResult:
        """Get preview of extracted content (first page only).
        
        Args:
            pdf_data: PDF file bytes.
            filename: Original filename.
            
        Returns:
            ExtractResult with truncated markdown preview.
        """
        result = self.extract_to_markdown(pdf_data, filename)
        if result.success and len(result.markdown) > 2000:
            result.markdown = result.markdown[:2000] + "\n\n... (continued)"
        return result

    def get_page_count(self, pdf_data: bytes) -> int:
        """Get number of pages in PDF.
        
        Args:
            pdf_data: PDF file bytes.
            
        Returns:
            Page count.
        """
        return self._get_page_count(pdf_data)

    def _extract_with_pdfminer(self, pdf_data: bytes) -> str:
        """Extract text using pdfminer.six (pure Python)."""
        try:
            from pdfminer.high_level import extract_text
            from io import BytesIO
        except ImportError:
            raise ImportError("pdfminer.six is required: pip install pdfminer.six")

        return extract_text(BytesIO(pdf_data))

    def _extract_with_pdftotext(
        self, pdf_data: bytes, safe_name: str, timestamp: int
    ) -> str:
        """Extract text using pdftotext CLI."""
        temp_pdf = Path(self.temp_dir) / f"{safe_name}_{timestamp}.pdf"
        temp_pdf.write_bytes(pdf_data)

        try:
            result = subprocess.run(
                ["pdftotext", "-layout", "-enc", "UTF-8", str(temp_pdf), "-"],
                capture_output=True,
                text=True,
                timeout=60,
            )
            return result.stdout
        finally:
            temp_pdf.unlink(missing_ok=True)

    def _get_page_count(self, pdf_data: bytes) -> int:
        """Get page count from PDF."""
        try:
            from pdfminer.pdfpage import PDFPage
            from io import BytesIO
            
            return sum(1 for _ in PDFPage.get_pages(BytesIO(pdf_data)))
        except:
            return 0

    def _text_to_markdown(self, text: str, title: str) -> str:
        """Convert extracted text to Markdown format."""
        result = [f"# {title}\n\n"]

        lines = text.split("\n")
        processed_lines = []

        for line in lines:
            line = line.strip()

            if not line:
                if processed_lines and processed_lines[-1] != "":
                    processed_lines.append("")
                continue

            # Detect headers (ALL CAPS)
            if self._is_potential_header(line):
                line = "## " + line.title()

            # Detect bullet points
            for bullet in ["•", "●", "○"]:
                if line.startswith(bullet):
                    line = "- " + line[1:].strip()
                    break

            processed_lines.append(line)

        # Join with proper spacing
        in_paragraph = False
        for i, line in enumerate(processed_lines):
            if not line:
                if in_paragraph:
                    result.append("\n\n")
                    in_paragraph = False
                continue

            if line.startswith("##") or line.startswith("- "):
                if in_paragraph:
                    result.append("\n\n")
                result.append(line + "\n")
                in_paragraph = False
            else:
                if in_paragraph and i > 0 and processed_lines[i - 1]:
                    result.append(" ")
                result.append(line)
                in_paragraph = True

        result.append("\n")
        return "".join(result)

    def _is_potential_header(self, line: str) -> bool:
        """Check if line is likely a header."""
        return (
            3 < len(line) < 60
            and line == line.upper()
            and len(line.split()) >= 2
        )

    def _sanitize_filename(self, name: str) -> str:
        """Sanitize filename."""
        return re.sub(r'[/\\:*?"<>| ]', "_", name)
