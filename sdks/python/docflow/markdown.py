"""Markdown parser for DocFlow."""

from typing import List, Optional

from .types import MDFile


class MarkdownParser:
    """Handles markdown processing and HTML conversion.
    
    Uses markdown-it-py for parsing with extensions support.
    
    Example:
        >>> parser = MarkdownParser()
        >>> html = parser.to_html("# Hello World")
    """

    def __init__(self) -> None:
        """Initialize the markdown parser."""
        try:
            from markdown_it import MarkdownIt
            from mdit_py_plugins.tasklists import tasklists_plugin
            from mdit_py_plugins.footnote import footnote_plugin
        except ImportError:
            raise ImportError(
                "markdown-it-py is required: pip install markdown-it-py mdit-py-plugins"
            )

        self.md = (
            MarkdownIt("gfm-like")
            .enable("table")
            .enable("strikethrough")
            .use(tasklists_plugin)
            .use(footnote_plugin)
        )

    def to_html(self, content: str) -> str:
        """Convert markdown content to HTML.
        
        Args:
            content: Markdown string.
            
        Returns:
            HTML string.
        """
        return self.md.render(content)

    def merge_files(self, files: List[MDFile]) -> str:
        """Merge multiple files into a single content string.
        
        Files are sorted by their order field.
        
        Args:
            files: List of MDFile objects.
            
        Returns:
            Merged markdown content.
        """
        if not files:
            return ""

        # Sort files by order
        sorted_files = sorted(files, key=lambda f: f.order)

        parts = []
        for i, file in enumerate(sorted_files):
            if i > 0:
                parts.append("\n\n---\n\n")
            parts.append(file.content)

        return "".join(parts)

    def merge_files_to_html(self, files: List[MDFile]) -> str:
        """Merge files and convert to HTML with file separators.
        
        Args:
            files: List of MDFile objects.
            
        Returns:
            HTML string with file separators.
        """
        if not files:
            return ""

        sorted_files = sorted(files, key=lambda f: f.order)

        parts = []
        for i, file in enumerate(sorted_files):
            if i > 0:
                parts.append(f'<div class="file-separator"><span>{file.name}</span></div>')
            else:
                parts.append(f'<div class="file-header"><span>{file.name}</span></div>')

            html = self.to_html(file.content)
            parts.append(f'<div class="file-content">{html}</div>')

        return "".join(parts)

    def estimate_page_count(self, content: str) -> int:
        """Estimate the number of PDF pages based on content.
        
        Args:
            content: Markdown or text content.
            
        Returns:
            Estimated page count (minimum 1).
        """
        chars_per_page = 3000
        pages = len(content) // chars_per_page
        return max(1, pages)
