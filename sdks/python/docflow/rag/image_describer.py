"""LLM image describer for DocFlow RAG mode."""

import base64
from typing import Optional

from ..types import ExtractedImage, LLMConfig


class LLMImageDescriber:
    """Generates descriptions for images using vision LLMs.
    
    Supports:
    - OpenAI GPT-4 Vision
    - Anthropic Claude Vision
    - Ollama (local models like LLaVA)
    
    Example:
        >>> config = LLMConfig(provider="openai", api_key="sk-...")
        >>> describer = LLMImageDescriber(config)
        >>> description = describer.describe(image)
    """

    def __init__(self, config: LLMConfig) -> None:
        """Initialize image describer.
        
        Args:
            config: LLM configuration.
        """
        self.config = config

    def describe(self, image: ExtractedImage) -> str:
        """Generate a description for an image.
        
        Args:
            image: Extracted image with data.
            
        Returns:
            Text description of the image.
        """
        if self.config.provider == "openai":
            return self._describe_openai(image)
        elif self.config.provider == "anthropic":
            return self._describe_anthropic(image)
        elif self.config.provider == "ollama":
            return self._describe_ollama(image)
        else:
            raise ValueError(f"Unknown LLM provider: {self.config.provider}")

    def describe_for_rag(self, image: ExtractedImage) -> str:
        """Generate RAG-optimized description.
        
        Includes:
        - Main content description
        - Data extraction if chart/graph
        - Key entities
        - Semantic context
        
        Args:
            image: Extracted image.
            
        Returns:
            RAG-optimized description.
        """
        prompt = self._build_rag_prompt(image)
        
        if self.config.provider == "openai":
            return self._call_openai(image, prompt)
        elif self.config.provider == "anthropic":
            return self._call_anthropic(image, prompt)
        elif self.config.provider == "ollama":
            return self._call_ollama(image, prompt)
        else:
            raise ValueError(f"Unknown LLM provider: {self.config.provider}")

    def _build_rag_prompt(self, image: ExtractedImage) -> str:
        """Build a prompt optimized for RAG extraction."""
        context = ""
        if image.surrounding_text:
            context = f"\n\nContext from surrounding text: {image.surrounding_text[:500]}"
        
        return f"""Analyze this image and provide a detailed description optimized for RAG (Retrieval-Augmented Generation).

Include:
1. **Main Content**: What is shown in the image
2. **Type**: Is this a chart, graph, diagram, photo, screenshot, etc.
3. **Data Extraction**: If it contains data (charts, tables), extract key values
4. **Key Entities**: People, objects, brands, locations mentioned or shown
5. **Semantic Tags**: 3-5 keywords for retrieval
{context}

Format your response as:
**Description**: [detailed description]
**Type**: [image type]
**Key Data**: [extracted data if applicable]
**Entities**: [list of entities]
**Tags**: [comma-separated tags]"""

    def _describe_openai(self, image: ExtractedImage) -> str:
        """Generate description using OpenAI."""
        prompt = "Describe this image in detail. Include any text, data, or important elements visible."
        return self._call_openai(image, prompt)

    def _call_openai(self, image: ExtractedImage, prompt: str) -> str:
        """Call OpenAI Vision API."""
        try:
            import openai
        except ImportError:
            raise ImportError("openai is required: pip install openai")

        client = openai.OpenAI(
            api_key=self.config.api_key,
            base_url=self.config.base_url,
            timeout=self.config.timeout,
        )

        # Encode image
        b64_image = base64.b64encode(image.data).decode("utf-8")
        media_type = f"image/{image.format}"

        response = client.chat.completions.create(
            model=self.config.model,
            messages=[
                {
                    "role": "user",
                    "content": [
                        {"type": "text", "text": prompt},
                        {
                            "type": "image_url",
                            "image_url": {
                                "url": f"data:{media_type};base64,{b64_image}",
                                "detail": self.config.detail,
                            },
                        },
                    ],
                }
            ],
            max_tokens=self.config.max_tokens,
        )

        return response.choices[0].message.content

    def _describe_anthropic(self, image: ExtractedImage) -> str:
        """Generate description using Anthropic Claude."""
        prompt = "Describe this image in detail. Include any text, data, or important elements visible."
        return self._call_anthropic(image, prompt)

    def _call_anthropic(self, image: ExtractedImage, prompt: str) -> str:
        """Call Anthropic Claude Vision API."""
        try:
            import anthropic
        except ImportError:
            raise ImportError("anthropic is required: pip install anthropic")

        client = anthropic.Anthropic(api_key=self.config.api_key)

        b64_image = base64.b64encode(image.data).decode("utf-8")
        media_type = f"image/{image.format}"

        response = client.messages.create(
            model=self.config.model,
            max_tokens=self.config.max_tokens,
            messages=[
                {
                    "role": "user",
                    "content": [
                        {
                            "type": "image",
                            "source": {
                                "type": "base64",
                                "media_type": media_type,
                                "data": b64_image,
                            },
                        },
                        {"type": "text", "text": prompt},
                    ],
                }
            ],
        )

        return response.content[0].text

    def _describe_ollama(self, image: ExtractedImage) -> str:
        """Generate description using Ollama local model."""
        prompt = "Describe this image in detail."
        return self._call_ollama(image, prompt)

    def _call_ollama(self, image: ExtractedImage, prompt: str) -> str:
        """Call Ollama local model."""
        try:
            import requests
        except ImportError:
            raise ImportError("requests is required: pip install requests")

        b64_image = base64.b64encode(image.data).decode("utf-8")
        
        base_url = self.config.base_url or "http://localhost:11434"
        url = f"{base_url}/api/generate"

        payload = {
            "model": self.config.model or "llava",
            "prompt": prompt,
            "images": [b64_image],
            "stream": False,
        }

        response = requests.post(url, json=payload, timeout=self.config.timeout)
        response.raise_for_status()

        return response.json().get("response", "")
