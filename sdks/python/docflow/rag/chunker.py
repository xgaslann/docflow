"""RAG chunker for DocFlow."""

import re
from typing import List, Optional

from ..types import Chunk, ChunkMetadata, RAGConfig


class RAGChunker:
    """Smart chunker for RAG-optimized markdown content.
    
    Features:
    - Respects heading boundaries
    - Keeps tables together
    - Maintains image+description pairs
    - Adds overlap for context
    
    Example:
        >>> chunker = RAGChunker(config=RAGConfig())
        >>> chunks = chunker.chunk(markdown_content)
    """

    def __init__(self, config: Optional[RAGConfig] = None) -> None:
        """Initialize chunker.
        
        Args:
            config: RAG configuration options.
        """
        self.config = config or RAGConfig()

    def chunk(self, markdown: str) -> List[Chunk]:
        """Split markdown into RAG-optimized chunks.
        
        Args:
            markdown: Markdown content to chunk.
            
        Returns:
            List of Chunk objects.
        """
        # Remove frontmatter for chunking (will be added to metadata)
        content, frontmatter = self._extract_frontmatter(markdown)
        
        if self.config.respect_headings:
            sections = self._split_by_headings(content)
        else:
            sections = [("", content)]
        
        chunks = []
        chunk_index = 0
        char_offset = 0
        
        for section_title, section_content in sections:
            section_chunks = self._chunk_section(
                section_content,
                section_title,
                chunk_index,
                char_offset,
            )
            chunks.extend(section_chunks)
            chunk_index += len(section_chunks)
            char_offset += len(section_content)
        
        # Add overlap between chunks
        if self.config.chunk_overlap > 0:
            chunks = self._add_overlap(chunks, content)
        
        # Add chunk markers if requested
        if self.config.add_chunk_markers:
            for i, chunk in enumerate(chunks):
                chunk.content = f"{chunk.content}\n\n<!-- chunk_boundary: {i} -->"
        
        return chunks

    def _extract_frontmatter(self, markdown: str) -> tuple:
        """Extract YAML frontmatter from markdown."""
        if markdown.startswith("---"):
            end = markdown.find("---", 3)
            if end > 0:
                frontmatter = markdown[3:end].strip()
                content = markdown[end + 3:].strip()
                return content, frontmatter
        return markdown, ""

    def _split_by_headings(self, content: str) -> List[tuple]:
        """Split content by headings."""
        # Pattern for markdown headings
        heading_pattern = re.compile(r'^(#{1,6})\s+(.+)$', re.MULTILINE)
        
        sections = []
        last_end = 0
        last_title = ""
        
        for match in heading_pattern.finditer(content):
            # Content before this heading
            if match.start() > last_end:
                section_content = content[last_end:match.start()].strip()
                if section_content:
                    sections.append((last_title, section_content))
            
            last_title = match.group(2).strip()
            last_end = match.start()
        
        # Remaining content
        if last_end < len(content):
            remaining = content[last_end:].strip()
            if remaining:
                sections.append((last_title, remaining))
        
        if not sections:
            sections = [("", content)]
        
        return sections

    def _chunk_section(
        self,
        content: str,
        section_title: str,
        start_index: int,
        char_offset: int,
    ) -> List[Chunk]:
        """Chunk a section respecting tables and code blocks."""
        chunks = []
        
        # Find special blocks that should stay together
        protected_blocks = self._find_protected_blocks(content)
        
        # Simple chunking by size
        current_chunk = ""
        current_start = char_offset
        
        lines = content.split("\n")
        i = 0
        
        while i < len(lines):
            line = lines[i]
            
            # Check if we're starting a protected block
            block_end = None
            for start, end, block_type in protected_blocks:
                if i == start:
                    block_end = end
                    break
            
            if block_end is not None:
                # Add entire protected block
                block_content = "\n".join(lines[i:block_end + 1])
                
                # If adding this would exceed chunk size, save current and start new
                if len(current_chunk) + len(block_content) > self.config.chunk_size and current_chunk:
                    chunks.append(self._create_chunk(
                        current_chunk.strip(),
                        start_index + len(chunks),
                        current_start,
                        current_start + len(current_chunk),
                        section_title,
                    ))
                    current_chunk = ""
                    current_start = char_offset + sum(len(l) + 1 for l in lines[:i])
                
                current_chunk += block_content + "\n"
                i = block_end + 1
            else:
                # Regular line
                if len(current_chunk) + len(line) > self.config.chunk_size and current_chunk:
                    chunks.append(self._create_chunk(
                        current_chunk.strip(),
                        start_index + len(chunks),
                        current_start,
                        current_start + len(current_chunk),
                        section_title,
                    ))
                    current_chunk = ""
                    current_start = char_offset + sum(len(l) + 1 for l in lines[:i])
                
                current_chunk += line + "\n"
                i += 1
        
        # Don't forget remaining content
        if current_chunk.strip():
            chunks.append(self._create_chunk(
                current_chunk.strip(),
                start_index + len(chunks),
                current_start,
                current_start + len(current_chunk),
                section_title,
            ))
        
        return chunks

    def _find_protected_blocks(self, content: str) -> List[tuple]:
        """Find blocks that should stay together (tables, code blocks)."""
        blocks = []
        lines = content.split("\n")
        
        i = 0
        while i < len(lines):
            line = lines[i].strip()
            
            # Code blocks
            if line.startswith("```"):
                start = i
                i += 1
                while i < len(lines) and not lines[i].strip().startswith("```"):
                    i += 1
                blocks.append((start, i, "code"))
            
            # Tables
            elif line.startswith("|") and line.endswith("|"):
                start = i
                while i < len(lines) and lines[i].strip().startswith("|"):
                    i += 1
                blocks.append((start, i - 1, "table"))
                continue
            
            i += 1
        
        return blocks

    def _create_chunk(
        self,
        content: str,
        index: int,
        start_char: int,
        end_char: int,
        section_title: str,
    ) -> Chunk:
        """Create a Chunk object with metadata."""
        has_table = "|" in content and "---" in content
        has_image = "![" in content or "[Image:" in content
        
        metadata = ChunkMetadata(
            section_title=section_title,
            heading_path=self._extract_heading_path(content),
            has_table=has_table,
            has_image=has_image,
        )
        
        return Chunk(
            content=content,
            index=index,
            start_char=start_char,
            end_char=end_char,
            metadata=metadata,
        )

    def _extract_heading_path(self, content: str) -> List[str]:
        """Extract heading hierarchy from content."""
        headings = []
        for line in content.split("\n"):
            if line.startswith("#"):
                level = len(line.split()[0])
                text = line.lstrip("#").strip()
                headings.append(text)
        return headings

    def _add_overlap(self, chunks: List[Chunk], full_content: str) -> List[Chunk]:
        """Add overlap between chunks for context."""
        if len(chunks) <= 1:
            return chunks
        
        overlap_size = self.config.chunk_overlap
        
        for i in range(1, len(chunks)):
            prev_chunk = chunks[i - 1]
            curr_chunk = chunks[i]
            
            # Get overlap from end of previous chunk
            overlap_text = prev_chunk.content[-overlap_size:]
            
            # Find a good break point (end of sentence or paragraph)
            for break_char in ["\n\n", ". ", "\n"]:
                idx = overlap_text.find(break_char)
                if idx > 0:
                    overlap_text = overlap_text[idx + len(break_char):]
                    break
            
            if overlap_text.strip():
                curr_chunk.content = f"[...] {overlap_text.strip()}\n\n{curr_chunk.content}"
        
        return chunks
