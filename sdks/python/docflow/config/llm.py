"""LLM configuration with custom prompts."""

from dataclasses import dataclass, field
from enum import Enum
from typing import Any, Callable, Dict, List, Optional


class LLMProvider(Enum):
    """LLM provider."""
    
    OPENAI = "openai"
    AZURE_OPENAI = "azure_openai"
    ANTHROPIC = "anthropic"
    OLLAMA = "ollama"
    GOOGLE = "google"
    COHERE = "cohere"


@dataclass
class LLMPrompts:
    """Custom prompts for LLM processing.
    
    All prompts support template variables:
    - {content}: The content being processed
    - {context}: Surrounding context
    - {filename}: Source filename
    - {format}: Source format
    
    Example:
        >>> prompts = LLMPrompts(
        ...     image_description="Describe this image for a RAG system. Focus on: {context}",
        ...     custom={"classify": "Classify this document: {content}"}
        ... )
    """
    
    # Image processing
    image_description: str = """Describe this image in detail for use in a document retrieval system.
Focus on:
1. Key information visible in the image
2. Text, numbers, or data shown
3. Context and relevance to the document
4. Any charts, graphs, or diagrams

Be concise but comprehensive."""

    image_ocr: str = """Extract all text visible in this image.
Preserve formatting and structure where possible."""

    image_data_extraction: str = """Extract structured data from this image.
Return as JSON with appropriate keys for the data found."""

    # Table processing
    table_analysis: str = """Analyze this table and provide:
1. Brief summary of what the table contains
2. Key insights or patterns
3. Important data points
4. Column descriptions

Table:
{content}"""

    table_to_text: str = """Convert this table to a natural language description:
{content}"""

    # Text processing
    text_summary: str = """Summarize the following content concisely.
Focus on main points and key information.
Maximum {max_length} characters.

Content:
{content}"""

    text_key_points: str = """Extract the {max_points} most important key points from:
{content}

Return as a JSON array of strings."""

    entity_extraction: str = """Extract all named entities from this content.
Include: people, organizations, locations, products, dates, monetary values.
Return as a JSON object with entity types as keys and arrays of entities as values.

Content:
{content}"""

    keyword_extraction: str = """Extract the most relevant keywords from:
{content}

Return as a JSON array of strings. Maximum {max_keywords} keywords."""

    # Document processing
    document_classification: str = """Classify this document into one of the following categories:
{categories}

Document:
{content}

Return only the category name."""

    document_qa: str = """Answer the following question based on the document content.
If the answer is not in the document, say "Not found in document."

Question: {question}

Document:
{content}"""

    # Custom prompts
    custom: Dict[str, str] = field(default_factory=dict)
    
    def get_prompt(self, name: str, **kwargs) -> str:
        """Get a prompt by name with variable substitution.
        
        Args:
            name: Prompt name (e.g., "image_description", "text_summary")
            **kwargs: Template variables to substitute
            
        Returns:
            Formatted prompt string.
        """
        # Check custom prompts first
        if name in self.custom:
            template = self.custom[name]
        elif hasattr(self, name):
            template = getattr(self, name)
        else:
            raise ValueError(f"Unknown prompt: {name}")
        
        # Substitute variables
        for key, value in kwargs.items():
            template = template.replace(f"{{{key}}}", str(value))
        
        return template
    
    def set_custom_prompt(self, name: str, prompt: str) -> None:
        """Set a custom prompt."""
        self.custom[name] = prompt
    
    def override_prompt(self, name: str, prompt: str) -> None:
        """Override a built-in prompt."""
        if hasattr(self, name) and name != "custom":
            setattr(self, name, prompt)
        else:
            self.custom[name] = prompt


@dataclass
class LLMConfig:
    """Configuration for LLM integration.
    
    Supports multiple providers with full configuration options.
    
    Example (OpenAI):
        >>> config = LLMConfig(
        ...     provider=LLMProvider.OPENAI,
        ...     model="gpt-4-vision-preview",
        ...     api_key="sk-...",
        ...     prompts=LLMPrompts(
        ...         image_description="Custom prompt..."
        ...     )
        ... )
    
    Example (Azure OpenAI):
        >>> config = LLMConfig(
        ...     provider=LLMProvider.AZURE_OPENAI,
        ...     azure_endpoint="https://xxx.openai.azure.com/",
        ...     azure_deployment="gpt-4-vision",
        ...     api_key="...",
        ...     api_version="2024-02-01"
        ... )
    """
    
    # Provider & Model
    provider: LLMProvider = LLMProvider.OPENAI
    model: str = "gpt-4-vision-preview"
    
    # Authentication
    api_key: str = ""
    
    # Custom prompts
    prompts: LLMPrompts = field(default_factory=LLMPrompts)
    
    # ============== OpenAI ==============
    
    organization: str = ""
    base_url: Optional[str] = None
    
    # ============== Azure OpenAI ==============
    
    azure_endpoint: str = ""
    azure_deployment: str = ""
    api_version: str = "2024-02-01"
    
    # ============== Anthropic ==============
    
    anthropic_version: str = "2023-06-01"
    
    # ============== Ollama ==============
    
    ollama_base_url: str = "http://localhost:11434"
    
    # ============== Google ==============
    
    google_project: str = ""
    google_location: str = "us-central1"
    
    # ============== Generation Parameters ==============
    
    temperature: float = 0.7
    max_tokens: int = 1000
    top_p: float = 1.0
    frequency_penalty: float = 0.0
    presence_penalty: float = 0.0
    
    # Stop sequences
    stop_sequences: List[str] = field(default_factory=list)
    
    # ============== Vision Parameters ==============
    
    detail: str = "auto"  # auto, low, high
    max_image_size: int = 20 * 1024 * 1024  # 20MB
    supported_formats: List[str] = field(default_factory=lambda: [
        "png", "jpg", "jpeg", "gif", "webp"
    ])
    
    # ============== Retry & Timeout ==============
    
    timeout: int = 60
    retry_count: int = 3
    retry_delay: float = 1.0
    retry_multiplier: float = 2.0
    
    # ============== Batch Processing ==============
    
    batch_size: int = 5
    concurrent_requests: int = 3
    
    # ============== Cost Control ==============
    
    max_cost_per_request: float = 0.0  # 0 = no limit
    track_usage: bool = True
    
    # ============== Response Processing ==============
    
    response_format: str = "text"  # text, json_object
    json_schema: Optional[Dict[str, Any]] = None
    
    # Custom response processor
    response_processor: Optional[Callable] = None
    
    def validate(self) -> None:
        """Validate configuration."""
        if not self.api_key and self.provider not in [LLMProvider.OLLAMA]:
            raise ValueError(f"API key is required for {self.provider.value}")
        
        if self.provider == LLMProvider.AZURE_OPENAI:
            if not self.azure_endpoint:
                raise ValueError("Azure endpoint is required")
            if not self.azure_deployment:
                raise ValueError("Azure deployment is required")
        
        if not 0 <= self.temperature <= 2:
            raise ValueError("Temperature must be between 0 and 2")
        
        if self.max_tokens <= 0:
            raise ValueError("max_tokens must be positive")
    
    def get_effective_base_url(self) -> Optional[str]:
        """Get the effective base URL for the provider."""
        if self.provider == LLMProvider.OPENAI:
            return self.base_url
        elif self.provider == LLMProvider.AZURE_OPENAI:
            return self.azure_endpoint
        elif self.provider == LLMProvider.OLLAMA:
            return self.ollama_base_url
        return None
