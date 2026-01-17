"""LLM Processor for enhanced RAG features."""

import base64
import json
from typing import Any, Dict, List, Optional, Union

from ..types import (
    ExtractedImage,
    ExtractedTable,
    LLMConfig,
    LLMProcessingMode,
    DocumentMetadata,
)


class LLMProcessor:
    """Unified LLM processor for images, tables, and text.
    
    Supports multiple providers: OpenAI, Anthropic, Ollama.
    Can process images, tables, text sections, or all content.
    
    Example:
        >>> config = LLMConfig(provider="openai", api_key="...")
        >>> processor = LLMProcessor(config)
        >>> 
        >>> # Process image
        >>> description = processor.describe_image(image)
        >>> 
        >>> # Process table
        >>> analysis = processor.analyze_table(table)
        >>> 
        >>> # Generate summary
        >>> summary = processor.generate_summary(content)
    """
    
    def __init__(self, config: LLMConfig) -> None:
        """Initialize LLM processor.
        
        Args:
            config: LLM configuration.
        """
        self.config = config
        self._client = None
    
    # ============== Image Processing ==============
    
    def describe_image(self, image: ExtractedImage, context: str = "") -> str:
        """Generate a description for an image.
        
        Args:
            image: Extracted image object.
            context: Surrounding text for context.
            
        Returns:
            Text description of the image.
        """
        prompt = self._build_image_prompt(image, context)
        return self._call_vision_api(image.data, prompt)
    
    def analyze_image_for_rag(self, image: ExtractedImage) -> Dict[str, Any]:
        """Perform full RAG analysis on an image.
        
        Returns structured analysis including description, entities, and data.
        """
        prompt = """Analyze this image for RAG (Retrieval-Augmented Generation):

1. **Description**: Detailed description of the image content
2. **Key Information**: Important facts, numbers, or data shown
3. **Entities**: People, organizations, locations, products mentioned
4. **Context**: How this image relates to document content
5. **Data Extraction**: Any text, charts, or tables visible

Respond in JSON format."""

        response = self._call_vision_api(image.data, prompt)
        
        try:
            return json.loads(response)
        except json.JSONDecodeError:
            return {
                "description": response,
                "key_information": [],
                "entities": [],
                "context": "",
                "data_extraction": ""
            }
    
    # ============== Table Processing ==============
    
    def analyze_table(self, table: ExtractedTable, context: str = "") -> str:
        """Generate analysis/summary for a table.
        
        Args:
            table: Extracted table object.
            context: Surrounding text for context.
            
        Returns:
            Analysis of the table content.
        """
        table_md = self._table_to_markdown(table)
        
        prompt = f"""Analyze this table and provide:
1. Brief summary of what the table contains
2. Key insights or patterns
3. Important data points

Table:
{table_md}

Context: {context if context else 'No additional context provided.'}"""

        return self._call_text_api(prompt)
    
    def extract_table_data(self, table: ExtractedTable) -> Dict[str, Any]:
        """Extract structured data from a table for RAG.
        
        Returns key-value pairs, trends, and statistics.
        """
        table_md = self._table_to_markdown(table)
        
        prompt = f"""Extract structured information from this table:

{table_md}

Respond with JSON containing:
{{
    "summary": "Brief table summary",
    "columns": ["column descriptions"],
    "key_values": {{"important": "values"}},
    "statistics": {{"if applicable": "stats"}},
    "trends": ["observed patterns"],
    "entities": ["mentioned entities"]
}}"""

        response = self._call_text_api(prompt)
        
        try:
            return json.loads(response)
        except json.JSONDecodeError:
            return {"summary": response}
    
    # ============== Text Processing ==============
    
    def generate_summary(self, content: str, max_length: int = 500) -> str:
        """Generate a summary of the content.
        
        Args:
            content: Document content.
            max_length: Maximum summary length.
            
        Returns:
            Summary text.
        """
        prompt = f"""Summarize the following content in {max_length} characters or less.
Focus on the main points and key information.

Content:
{content[:8000]}"""  # Limit input size

        return self._call_text_api(prompt)
    
    def extract_key_points(self, content: str, max_points: int = 5) -> List[str]:
        """Extract key points from content.
        
        Args:
            content: Document content.
            max_points: Maximum number of points.
            
        Returns:
            List of key points.
        """
        prompt = f"""Extract the {max_points} most important key points from this content.
Return as a JSON array of strings.

Content:
{content[:8000]}"""

        response = self._call_text_api(prompt)
        
        try:
            return json.loads(response)
        except json.JSONDecodeError:
            # Parse as bullet points
            lines = response.strip().split('\n')
            return [line.lstrip('- •0123456789.').strip() for line in lines if line.strip()][:max_points]
    
    def extract_entities(self, content: str) -> List[str]:
        """Extract named entities from content.
        
        Returns list of entities (people, orgs, places, etc.).
        """
        prompt = """Extract all named entities from this content.
Include: people, organizations, locations, products, dates, numbers.
Return as a JSON array of strings.

Content:
""" + content[:6000]

        response = self._call_text_api(prompt)
        
        try:
            return json.loads(response)
        except json.JSONDecodeError:
            return []
    
    def enhance_metadata(self, metadata: DocumentMetadata, content: str) -> DocumentMetadata:
        """Enhance document metadata using LLM.
        
        Adds summary, key points, entities, and keywords.
        """
        # Generate summary
        metadata.summary = self.generate_summary(content)
        
        # Extract key points
        metadata.key_points = self.extract_key_points(content)
        
        # Extract entities
        metadata.entities = self.extract_entities(content)
        
        return metadata
    
    # ============== Batch Processing ==============
    
    def process_images(self, images: List[ExtractedImage]) -> List[ExtractedImage]:
        """Process multiple images."""
        for img in images:
            try:
                analysis = self.analyze_image_for_rag(img)
                img.description = analysis.get("description", "")
                img.llm_analysis = analysis
            except Exception as e:
                img.description = f"[Analysis failed: {e}]"
        return images
    
    def process_tables(self, tables: List[ExtractedTable]) -> List[ExtractedTable]:
        """Process multiple tables."""
        for table in tables:
            try:
                analysis = self.extract_table_data(table)
                table.summary = analysis.get("summary", "")
                table.llm_analysis = analysis
            except Exception as e:
                table.summary = f"[Analysis failed: {e}]"
        return tables
    
    # ============== API Calls ==============
    
    def _call_vision_api(self, image_data: bytes, prompt: str) -> str:
        """Call vision API with image."""
        if self.config.provider == "openai":
            return self._call_openai_vision(image_data, prompt)
        elif self.config.provider == "anthropic":
            return self._call_anthropic_vision(image_data, prompt)
        elif self.config.provider == "ollama":
            return self._call_ollama_vision(image_data, prompt)
        else:
            raise ValueError(f"Unsupported provider: {self.config.provider}")
    
    def _call_text_api(self, prompt: str) -> str:
        """Call text API."""
        if self.config.provider == "openai":
            return self._call_openai_text(prompt)
        elif self.config.provider == "anthropic":
            return self._call_anthropic_text(prompt)
        elif self.config.provider == "ollama":
            return self._call_ollama_text(prompt)
        else:
            raise ValueError(f"Unsupported provider: {self.config.provider}")
    
    def _call_openai_vision(self, image_data: bytes, prompt: str) -> str:
        """Call OpenAI Vision API."""
        try:
            from openai import OpenAI
        except ImportError:
            raise ImportError("openai is required: pip install openai")
        
        client = OpenAI(
            api_key=self.config.api_key,
            base_url=self.config.base_url,
            timeout=self.config.timeout,
        )
        
        base64_image = base64.b64encode(image_data).decode("utf-8")
        
        response = client.chat.completions.create(
            model=self.config.model,
            messages=[{
                "role": "user",
                "content": [
                    {"type": "text", "text": prompt},
                    {
                        "type": "image_url",
                        "image_url": {
                            "url": f"data:image/png;base64,{base64_image}",
                            "detail": self.config.detail,
                        },
                    },
                ],
            }],
            max_tokens=self.config.max_tokens,
            temperature=self.config.temperature,
        )
        
        return response.choices[0].message.content
    
    def _call_openai_text(self, prompt: str) -> str:
        """Call OpenAI Text API."""
        try:
            from openai import OpenAI
        except ImportError:
            raise ImportError("openai is required: pip install openai")
        
        client = OpenAI(
            api_key=self.config.api_key,
            base_url=self.config.base_url,
            timeout=self.config.timeout,
        )
        
        model = self.config.model
        if "vision" in model:
            model = "gpt-4"  # Use text model for text-only requests
        
        response = client.chat.completions.create(
            model=model,
            messages=[{"role": "user", "content": prompt}],
            max_tokens=self.config.max_tokens,
            temperature=self.config.temperature,
        )
        
        return response.choices[0].message.content
    
    def _call_anthropic_vision(self, image_data: bytes, prompt: str) -> str:
        """Call Anthropic Vision API."""
        try:
            import anthropic
        except ImportError:
            raise ImportError("anthropic is required: pip install anthropic")
        
        client = anthropic.Anthropic(api_key=self.config.api_key)
        base64_image = base64.b64encode(image_data).decode("utf-8")
        
        response = client.messages.create(
            model=self.config.model,
            max_tokens=self.config.max_tokens,
            messages=[{
                "role": "user",
                "content": [
                    {
                        "type": "image",
                        "source": {
                            "type": "base64",
                            "media_type": "image/png",
                            "data": base64_image,
                        },
                    },
                    {"type": "text", "text": prompt},
                ],
            }],
        )
        
        return response.content[0].text
    
    def _call_anthropic_text(self, prompt: str) -> str:
        """Call Anthropic Text API."""
        try:
            import anthropic
        except ImportError:
            raise ImportError("anthropic is required: pip install anthropic")
        
        client = anthropic.Anthropic(api_key=self.config.api_key)
        
        response = client.messages.create(
            model=self.config.model.replace("-vision", ""),
            max_tokens=self.config.max_tokens,
            messages=[{"role": "user", "content": prompt}],
        )
        
        return response.content[0].text
    
    def _call_ollama_vision(self, image_data: bytes, prompt: str) -> str:
        """Call Ollama Vision API."""
        import requests
        
        base_url = self.config.base_url or "http://localhost:11434"
        base64_image = base64.b64encode(image_data).decode("utf-8")
        
        response = requests.post(
            f"{base_url}/api/generate",
            json={
                "model": self.config.model,
                "prompt": prompt,
                "images": [base64_image],
                "stream": False,
            },
            timeout=self.config.timeout,
        )
        response.raise_for_status()
        
        return response.json().get("response", "")
    
    def _call_ollama_text(self, prompt: str) -> str:
        """Call Ollama Text API."""
        import requests
        
        base_url = self.config.base_url or "http://localhost:11434"
        
        response = requests.post(
            f"{base_url}/api/generate",
            json={
                "model": self.config.model,
                "prompt": prompt,
                "stream": False,
            },
            timeout=self.config.timeout,
        )
        response.raise_for_status()
        
        return response.json().get("response", "")
    
    # ============== Helpers ==============
    
    def _build_image_prompt(self, image: ExtractedImage, context: str) -> str:
        """Build prompt for image description."""
        prompt = "Describe this image in detail for use in a document retrieval system."
        
        if image.caption:
            prompt += f"\n\nOriginal caption: {image.caption}"
        
        if context:
            prompt += f"\n\nSurrounding context: {context}"
        
        prompt += "\n\nFocus on: key information, text visible, data shown, and relevance to the document."
        
        return prompt
    
    def _table_to_markdown(self, table: ExtractedTable) -> str:
        """Convert table to markdown format."""
        lines = []
        
        if table.header:
            lines.append("| " + " | ".join(table.header) + " |")
            lines.append("| " + " | ".join(["---"] * len(table.header)) + " |")
        
        for row in table.rows:
            lines.append("| " + " | ".join(row) + " |")
        
        return "\n".join(lines)
