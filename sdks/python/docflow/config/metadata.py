"""Metadata configuration."""

from dataclasses import dataclass, field
from typing import Any, Callable, Dict, List, Optional


@dataclass
class MetadataConfig:
    """Configuration for document metadata extraction.
    
    Allows fine-grained control over which metadata fields to extract,
    exclude, or add custom fields.
    
    Example:
        >>> config = MetadataConfig(
        ...     include_fields=["title", "headings", "entities"],
        ...     exclude_fields=["author"],
        ...     custom_fields={"project": "MyProject", "version": "1.0"}
        ... )
    """
    
    # Field selection
    include_fields: List[str] = field(default_factory=lambda: [
        "title", "headings", "table_of_contents", "word_count"
    ])
    exclude_fields: List[str] = field(default_factory=list)
    
    # Custom metadata
    custom_fields: Dict[str, Any] = field(default_factory=dict)
    
    # Extraction toggles
    extract_title: bool = True
    extract_author: bool = True
    extract_created_date: bool = True
    extract_modified_date: bool = True
    extract_headings: bool = True
    extract_toc: bool = True
    extract_word_count: bool = True
    extract_page_count: bool = True
    extract_image_count: bool = True
    extract_table_count: bool = True
    extract_language: bool = False  # Requires detection
    extract_keywords: bool = False  # Requires LLM or NLP
    extract_entities: bool = False  # Requires LLM
    extract_summary: bool = False   # Requires LLM
    extract_key_points: bool = False  # Requires LLM
    
    # Language detection
    language_detection_method: str = "auto"  # auto, langdetect, fasttext
    
    # Heading extraction
    max_heading_level: int = 6
    heading_numbering: bool = False
    build_heading_tree: bool = True
    
    # TOC generation
    toc_max_depth: int = 3
    toc_include_page_numbers: bool = True
    
    # Custom extractors
    custom_extractors: Dict[str, Callable] = field(default_factory=dict)
    
    def should_extract(self, field_name: str) -> bool:
        """Check if a field should be extracted."""
        if field_name in self.exclude_fields:
            return False
        if self.include_fields and field_name not in self.include_fields:
            return False
        return True
    
    def add_custom_field(self, name: str, value: Any) -> None:
        """Add a custom metadata field."""
        self.custom_fields[name] = value
    
    def remove_field(self, name: str) -> None:
        """Exclude a field from extraction."""
        if name not in self.exclude_fields:
            self.exclude_fields.append(name)
    
    def add_extractor(self, name: str, extractor: Callable) -> None:
        """Add a custom metadata extractor.
        
        The extractor should be a callable that takes document content
        and returns the extracted value.
        
        Example:
            >>> config.add_extractor("word_count", lambda doc: len(doc.split()))
        """
        self.custom_extractors[name] = extractor
    
    def get_enabled_extractors(self) -> List[str]:
        """Get list of enabled standard extractors."""
        extractors = []
        if self.extract_title:
            extractors.append("title")
        if self.extract_author:
            extractors.append("author")
        if self.extract_headings:
            extractors.append("headings")
        if self.extract_toc:
            extractors.append("toc")
        if self.extract_entities:
            extractors.append("entities")
        if self.extract_summary:
            extractors.append("summary")
        return extractors
