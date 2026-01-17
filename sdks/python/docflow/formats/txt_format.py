"""TXT format converter for DocFlow."""

import re
from typing import Optional

from ..types import ConvertResult


class TXTConverter:
    """Converts between plain text and Markdown formats.
    
    Example:
        >>> converter = TXTConverter()
        >>> md = converter.to_markdown(text_data, "notes.txt")
    """

    def __init__(
        self,
        detect_structure: bool = True,
        line_break_mode: str = "paragraph",  # "paragraph" or "preserve"
    ) -> None:
        """Initialize TXT converter.
        
        Args:
            detect_structure: Try to detect headings, lists, etc.
            line_break_mode: How to handle line breaks.
        """
        self.detect_structure = detect_structure
        self.line_break_mode = line_break_mode

    def to_markdown(
        self,
        text_data: str,
        filename: str = "document.txt",
        include_metadata: bool = True,
    ) -> ConvertResult:
        """Convert plain text to Markdown.
        
        Args:
            text_data: Text content.
            filename: Original filename.
            include_metadata: Include YAML frontmatter.
            
        Returns:
            ConvertResult with markdown content.
        """
        try:
            if isinstance(text_data, bytes):
                text_data = text_data.decode("utf-8")

            md_parts = []

            # Metadata
            if include_metadata:
                lines_count = len(text_data.split("\n"))
                word_count = len(text_data.split())
                md_parts.append("---")
                md_parts.append(f"source: {filename}")
                md_parts.append("format: txt")
                md_parts.append(f"lines: {lines_count}")
                md_parts.append(f"words: {word_count}")
                md_parts.append(f"characters: {len(text_data)}")
                md_parts.append("---\n")

            # Title
            title = filename.replace(".txt", "").replace("_", " ").title()
            md_parts.append(f"# {title}\n")

            if self.detect_structure:
                md_parts.append(self._detect_and_convert(text_data))
            else:
                if self.line_break_mode == "preserve":
                    md_parts.append(text_data)
                else:
                    md_parts.append(self._paragraphize(text_data))

            return ConvertResult(
                success=True,
                content="\n".join(md_parts),
                format="txt",
                metadata={"lines": len(text_data.split("\n"))},
            )

        except Exception as e:
            return ConvertResult(success=False, error=str(e))

    def from_markdown(self, markdown: str) -> ConvertResult:
        """Convert Markdown to plain text.
        
        Args:
            markdown: Markdown content.
            
        Returns:
            ConvertResult with plain text.
        """
        try:
            text = markdown

            # Remove frontmatter
            if text.startswith("---"):
                end = text.find("---", 3)
                if end > 0:
                    text = text[end + 3:].strip()

            # Remove markdown formatting
            # Headers -> just text
            text = re.sub(r'^#{1,6}\s+', '', text, flags=re.MULTILINE)
            
            # Bold/italic
            text = re.sub(r'\*\*(.+?)\*\*', r'\1', text)
            text = re.sub(r'\*(.+?)\*', r'\1', text)
            text = re.sub(r'__(.+?)__', r'\1', text)
            text = re.sub(r'_(.+?)_', r'\1', text)
            
            # Code
            text = re.sub(r'`(.+?)`', r'\1', text)
            text = re.sub(r'```[\s\S]*?```', lambda m: m.group().replace('```', ''), text)
            
            # Links
            text = re.sub(r'\[([^\]]+)\]\([^)]+\)', r'\1', text)
            
            # Images
            text = re.sub(r'!\[([^\]]*)\]\([^)]+\)', r'[Image: \1]', text)
            
            # Lists
            text = re.sub(r'^[-*+]\s+', '• ', text, flags=re.MULTILINE)
            text = re.sub(r'^\d+\.\s+', '', text, flags=re.MULTILINE)
            
            # Tables -> simple format
            lines = []
            for line in text.split("\n"):
                if line.strip().startswith("|"):
                    cells = [c.strip() for c in line.split("|")[1:-1]]
                    if not all(set(c) <= {"-", ":"} for c in cells):
                        lines.append(" | ".join(cells))
                else:
                    lines.append(line)
            text = "\n".join(lines)
            
            # Clean up
            text = re.sub(r'\n{3,}', '\n\n', text)
            text = text.strip()

            return ConvertResult(
                success=True,
                content=text,
                format="txt",
            )

        except Exception as e:
            return ConvertResult(success=False, error=str(e))

    def _detect_and_convert(self, text: str) -> str:
        """Detect structure in plain text and convert to markdown."""
        lines = text.split("\n")
        result = []
        
        for i, line in enumerate(lines):
            stripped = line.strip()
            
            if not stripped:
                result.append("")
                continue
            
            # Detect potential headers (ALL CAPS with reasonable length)
            if (
                stripped == stripped.upper() 
                and len(stripped) > 3 
                and len(stripped) < 80
                and not stripped.startswith(("•", "-", "*", "●"))
                and len(stripped.split()) >= 2
            ):
                result.append(f"\n## {stripped.title()}\n")
                continue
            
            # Detect bullet points
            for bullet in ["•", "●", "○", "▪", "▸"]:
                if stripped.startswith(bullet):
                    result.append(f"- {stripped[1:].strip()}")
                    break
            else:
                # Detect numbered lists
                match = re.match(r'^(\d+)[.)]\s+(.+)$', stripped)
                if match:
                    result.append(f"{match.group(1)}. {match.group(2)}")
                else:
                    result.append(stripped)
        
        return "\n".join(result)

    def _paragraphize(self, text: str) -> str:
        """Convert single line breaks to spaces, keep double line breaks."""
        paragraphs = re.split(r'\n\s*\n', text)
        result = []
        
        for para in paragraphs:
            # Join lines within paragraph
            para = ' '.join(para.split('\n'))
            para = ' '.join(para.split())  # Normalize whitespace
            if para.strip():
                result.append(para)
        
        return "\n\n".join(result)
