"""Markdown to PDF converter for DocFlow."""

import os
import re
import tempfile
import time
from pathlib import Path
from typing import List, Optional, Tuple, Union

from .markdown import MarkdownParser
from .template import Template
from .types import MDFile, ConvertOptions, PDFResult
from .storage.base import Storage
from .storage.local import LocalStorage


class Converter:
    """Converts Markdown files to PDF.
    
    Supports multiple PDF engines:
    - weasyprint (default, pure Python)
    - pdfkit/wkhtmltopdf (alternative)
    - playwright/chromium (best quality)
    
    Example:
        >>> converter = Converter(storage=LocalStorage("./output"))
        >>> result = converter.convert_to_pdf([
        ...     MDFile(name="doc.md", content="# Hello World")
        ... ])
    """

    def __init__(
        self,
        storage: Optional[Storage] = None,
        temp_dir: Optional[str] = None,
        engine: str = "weasyprint",
    ) -> None:
        """Initialize the converter.
        
        Args:
            storage: Storage backend for saving output files.
            temp_dir: Directory for temporary files.
            engine: PDF engine to use ('weasyprint', 'pdfkit', or 'playwright').
        """
        self.storage = storage
        self.temp_dir = temp_dir or tempfile.gettempdir()
        self.engine = engine
        self.parser = MarkdownParser()
        self.template = Template()

    def convert_to_pdf(
        self,
        files: List[Union[MDFile, dict]],
        options: Optional[ConvertOptions] = None,
    ) -> PDFResult:
        """Convert markdown files to PDF.
        
        Args:
            files: List of MDFile objects or dicts with name and content.
            options: Conversion options.
            
        Returns:
            PDFResult with success status and file paths.
        """
        if not files:
            return PDFResult(success=False, error="At least one file is required")

        # Normalize files
        md_files = self._normalize_files(files)
        
        if options is None:
            options = ConvertOptions()

        timestamp = int(time.time())

        try:
            if options.merge_mode == "merged":
                path, data = self._convert_merged(md_files, options.output_name, timestamp)
                return PDFResult(success=True, file_paths=[path], bytes_data=data)
            else:
                paths = []
                for file in md_files:
                    path, _ = self._convert_single(file, timestamp)
                    paths.append(path)
                return PDFResult(success=True, file_paths=paths)
        except Exception as e:
            return PDFResult(success=False, error=str(e))

    def convert_to_bytes(self, files: List[Union[MDFile, dict]]) -> bytes:
        """Convert markdown to PDF and return bytes.
        
        Args:
            files: List of MDFile objects or dicts.
            
        Returns:
            PDF bytes.
        """
        md_files = self._normalize_files(files)
        _, data = self._convert_merged(md_files, None, int(time.time()))
        return data

    def preview(self, content: str) -> str:
        """Generate HTML preview of markdown content.
        
        Args:
            content: Markdown content.
            
        Returns:
            HTML string.
        """
        return self.parser.to_html(content)

    def _normalize_files(self, files: List[Union[MDFile, dict]]) -> List[MDFile]:
        """Convert dicts to MDFile objects."""
        result = []
        for i, f in enumerate(files):
            if isinstance(f, MDFile):
                f.order = i
                result.append(f)
            else:
                result.append(MDFile(
                    name=f["name"],
                    content=f["content"],
                    id=f.get("id", f["name"]),
                    order=f.get("order", i),
                ))
        return result

    def _convert_merged(
        self, files: List[MDFile], output_name: Optional[str], timestamp: int
    ) -> Tuple[str, bytes]:
        """Convert merged files to PDF."""
        merged_content = self.parser.merge_files(files)
        
        if not output_name:
            output_name = f"merged_{timestamp}"
        output_name = self._sanitize_filename(output_name)

        return self._generate_pdf(merged_content, output_name)

    def _convert_single(self, file: MDFile, timestamp: int) -> Tuple[str, bytes]:
        """Convert a single file to PDF."""
        base_name = Path(file.name).stem
        output_name = f"{self._sanitize_filename(base_name)}_{timestamp}"
        return self._generate_pdf(file.content, output_name)

    def _generate_pdf(self, md_content: str, output_name: str) -> Tuple[str, bytes]:
        """Generate PDF from markdown content."""
        # Convert to HTML
        html_content = self.parser.to_html(md_content)
        full_html = self.template.generate(html_content)

        # Generate PDF based on engine
        if self.engine == "weasyprint":
            pdf_data = self._generate_with_weasyprint(full_html)
        elif self.engine == "pdfkit":
            pdf_data = self._generate_with_pdfkit(full_html)
        elif self.engine == "playwright":
            pdf_data = self._generate_with_playwright(full_html)
        else:
            # Default to weasyprint
            pdf_data = self._generate_with_weasyprint(full_html)

        # Save to storage
        output_path = f"{output_name}.pdf"
        if self.storage:
            self.storage.save(output_path, pdf_data)
            return self.storage.get_url(output_path) or output_path, pdf_data
        else:
            # Save to temp
            temp_path = Path(self.temp_dir) / output_path
            temp_path.write_bytes(pdf_data)
            return str(temp_path), pdf_data

    def _generate_with_weasyprint(self, html: str) -> bytes:
        """Generate PDF using WeasyPrint."""
        try:
            from weasyprint import HTML
        except ImportError:
            raise ImportError("weasyprint is required: pip install weasyprint")

        return HTML(string=html).write_pdf()

    def _generate_with_pdfkit(self, html: str) -> bytes:
        """Generate PDF using pdfkit (wkhtmltopdf)."""
        try:
            import pdfkit
        except ImportError:
            raise ImportError("pdfkit is required: pip install pdfkit")

        return pdfkit.from_string(html, False, options={
            "page-size": "A4",
            "margin-top": "20mm",
            "margin-bottom": "20mm",
            "margin-left": "20mm",
            "margin-right": "20mm",
            "encoding": "UTF-8",
        })

    def _generate_with_playwright(self, html: str) -> bytes:
        """Generate PDF using Playwright (Chromium)."""
        try:
            from playwright.sync_api import sync_playwright
        except ImportError:
            raise ImportError("playwright is required: pip install playwright")

        # Write HTML to temp file
        temp_html = Path(self.temp_dir) / f"temp_{int(time.time())}.html"
        temp_html.write_text(html, encoding="utf-8")

        try:
            with sync_playwright() as p:
                browser = p.chromium.launch()
                page = browser.new_page()
                page.goto(f"file://{temp_html}")
                pdf_data = page.pdf(
                    format="A4",
                    margin={"top": "20mm", "bottom": "20mm", "left": "20mm", "right": "20mm"},
                    print_background=True,
                )
                browser.close()
                return pdf_data
        finally:
            temp_html.unlink(missing_ok=True)

    def _sanitize_filename(self, name: str) -> str:
        """Sanitize filename for safe file system use."""
        return re.sub(r'[/\\:*?"<>| ]', "_", name)
