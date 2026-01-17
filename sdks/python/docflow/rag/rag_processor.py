"""Main RAG processor for DocFlow."""

import re
from datetime import datetime
from pathlib import Path
from typing import Optional, Union

from ..types import (
    ConvertResult,
    ExtractedImage,
    ExtractedTable,
    LLMConfig,
    RAGConfig,
    RAGDocument,
)
from .chunker import RAGChunker
from .image_describer import LLMImageDescriber


class RAGProcessor:
    """Unified RAG processor for all document formats.
    
    Converts any supported format to RAG-optimized Markdown with:
    - Metadata preservation
    - Image extraction and LLM description
    - Table structure preservation
    - Smart chunking
    
    Example:
        >>> config = RAGConfig(
        ...     extract_images=True,
        ...     describe_images=True,
        ...     llm_config=LLMConfig(provider="openai", api_key="...")
        ... )
        >>> processor = RAGProcessor(config)
        >>> result = processor.process_file("document.pdf")
    """

    def __init__(self, config: Optional[RAGConfig] = None) -> None:
        """Initialize RAG processor.
        
        Args:
            config: RAG configuration options.
        """
        self.config = config or RAGConfig()
        self.chunker = RAGChunker(self.config)
        
        if self.config.describe_images and self.config.llm_config:
            self.image_describer = LLMImageDescriber(self.config.llm_config)
        else:
            self.image_describer = None

    def process(
        self,
        data: Union[str, bytes],
        filename: str,
        format: Optional[str] = None,
    ) -> RAGDocument:
        """Process a document and convert to RAG-optimized format.
        
        Args:
            data: File content as string or bytes.
            filename: Original filename.
            format: File format (auto-detected if not provided).
            
        Returns:
            RAGDocument with markdown content and chunks.
        """
        if format is None:
            format = Path(filename).suffix.lower().lstrip(".")
        
        # Convert to markdown using appropriate converter
        result = self._convert_to_markdown(data, filename, format)
        
        if not result.success:
            raise ValueError(f"Conversion failed: {result.error}")
        
        # Process images with LLM if configured
        images = result.images
        if self.config.describe_images and self.image_describer and images:
            for img in images:
                try:
                    img.description = self.image_describer.describe_for_rag(img)
                except Exception as e:
                    img.description = f"[Image description failed: {e}]"
        
        # Build final markdown with image descriptions
        content = self._build_rag_markdown(result.content, images)
        
        # Extract tables for metadata
        tables = self._extract_tables(content)
        
        # Chunk the content
        chunks = self.chunker.chunk(content)
        
        return RAGDocument(
            content=content,
            chunks=chunks,
            images=images,
            tables=tables,
            metadata=result.metadata,
            source_file=filename,
            source_format=format,
        )

    def process_file(self, filepath: str) -> RAGDocument:
        """Process a file from disk.
        
        Args:
            filepath: Path to the file.
            
        Returns:
            RAGDocument.
        """
        path = Path(filepath)
        
        if path.suffix.lower() in [".pdf", ".docx", ".xlsx", ".xls"]:
            data = path.read_bytes()
        else:
            data = path.read_text(encoding="utf-8")
        
        return self.process(data, path.name)

    def _convert_to_markdown(
        self, data: Union[str, bytes], filename: str, format: str
    ) -> ConvertResult:
        """Convert file to markdown using appropriate converter."""
        
        if format == "pdf":
            from ..extractor import Extractor
            extractor = Extractor()
            result = extractor.extract_to_markdown(data, filename)
            return ConvertResult(
                success=result.success,
                content=result.markdown,
                format="pdf",
                error=result.error,
                metadata={"page_count": result.page_count},
            )
        
        elif format == "docx":
            from ..formats import DOCXConverter
            converter = DOCXConverter(extract_images=self.config.extract_images)
            return converter.to_markdown(data, filename)
        
        elif format in ["xlsx", "xls"]:
            from ..formats import ExcelConverter
            converter = ExcelConverter(include_all_sheets=True)
            return converter.to_markdown(data, filename)
        
        elif format == "csv":
            from ..formats import CSVConverter
            converter = CSVConverter()
            return converter.to_markdown(data, filename)
        
        elif format == "txt":
            from ..formats import TXTConverter
            converter = TXTConverter(detect_structure=True)
            return converter.to_markdown(data, filename)
        
        elif format == "md":
            # Already markdown
            return ConvertResult(
                success=True,
                content=data if isinstance(data, str) else data.decode("utf-8"),
                format="md",
            )
        
        else:
            return ConvertResult(
                success=False,
                error=f"Unsupported format: {format}",
            )

    def _build_rag_markdown(
        self, content: str, images: list
    ) -> str:
        """Build final markdown with image descriptions."""
        if not images:
            return content
        
        # Add image descriptions section if we have described images
        described_images = [img for img in images if img.description]
        
        if described_images:
            content += "\n\n---\n\n## Extracted Images\n\n"
            
            for i, img in enumerate(described_images, 1):
                content += f"### Image {i}: {img.filename}\n\n"
                if img.caption:
                    content += f"**Caption**: {img.caption}\n\n"
                content += f"**AI Analysis**:\n{img.description}\n\n"
        
        return content

    def _extract_tables(self, content: str) -> list:
        """Extract tables from markdown for metadata."""
        tables = []
        lines = content.split("\n")
        
        current_table = []
        in_table = False
        
        for line in lines:
            line = line.strip()
            
            if line.startswith("|") and line.endswith("|"):
                # Skip separator
                if set(line.replace("|", "").replace("-", "").replace(":", "").strip()) == set():
                    continue
                cells = [c.strip() for c in line.split("|")[1:-1]]
                current_table.append(cells)
                in_table = True
            else:
                if in_table and current_table:
                    header = current_table[0] if current_table else []
                    tables.append(ExtractedTable(
                        rows=current_table[1:] if len(current_table) > 1 else [],
                        header=header,
                    ))
                    current_table = []
                in_table = False
        
        if current_table:
            header = current_table[0] if current_table else []
            tables.append(ExtractedTable(
                rows=current_table[1:] if len(current_table) > 1 else [],
                header=header,
            ))
        
        return tables
