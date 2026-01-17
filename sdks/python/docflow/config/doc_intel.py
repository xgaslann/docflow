"""Document Intelligence configuration."""

from dataclasses import dataclass, field
from enum import Enum
from typing import List, Optional


class DocIntelProvider(Enum):
    """Document Intelligence provider."""
    
    AZURE = "azure"
    AWS = "aws"


class AzureModelId(Enum):
    """Azure Document Intelligence model IDs."""
    
    PREBUILT_READ = "prebuilt-read"
    PREBUILT_LAYOUT = "prebuilt-layout"
    PREBUILT_DOCUMENT = "prebuilt-document"
    PREBUILT_INVOICE = "prebuilt-invoice"
    PREBUILT_RECEIPT = "prebuilt-receipt"
    PREBUILT_ID_DOCUMENT = "prebuilt-idDocument"
    PREBUILT_BUSINESS_CARD = "prebuilt-businessCard"
    PREBUILT_HEALTH_INSURANCE_CARD = "prebuilt-healthInsuranceCard.us"


class AzureFeature(Enum):
    """Azure Document Intelligence features."""
    
    OCR_HIGH_RESOLUTION = "ocrHighResolution"
    FORMULAS = "formulas"
    STYLE_FONT = "styleFont"
    BARCODES = "barcodes"
    LANGUAGES = "languages"
    KEY_VALUE_PAIRS = "keyValuePairs"
    QUERY_FIELDS = "queryFields"


class TextractFeature(Enum):
    """AWS Textract features."""
    
    TABLES = "TABLES"
    FORMS = "FORMS"
    QUERIES = "QUERIES"
    SIGNATURES = "SIGNATURES"
    LAYOUT = "LAYOUT"


@dataclass
class DocIntelConfig:
    """Configuration for Document Intelligence services.
    
    Supports Azure Document Intelligence and AWS Textract with
    all available options exposed.
    
    Example (Azure):
        >>> config = DocIntelConfig(
        ...     provider=DocIntelProvider.AZURE,
        ...     endpoint="https://xxx.cognitiveservices.azure.com/",
        ...     api_key="your-key",
        ...     model_id=AzureModelId.PREBUILT_LAYOUT,
        ...     features=[AzureFeature.TABLES, AzureFeature.KEY_VALUE_PAIRS]
        ... )
    
    Example (AWS):
        >>> config = DocIntelConfig(
        ...     provider=DocIntelProvider.AWS,
        ...     aws_region="us-east-1",
        ...     textract_features=[TextractFeature.TABLES, TextractFeature.FORMS]
        ... )
    """
    
    provider: DocIntelProvider = DocIntelProvider.AZURE
    
    # ============== Azure Document Intelligence ==============
    
    # Connection
    endpoint: str = ""
    api_key: str = ""
    
    # Model
    model_id: str = "prebuilt-layout"
    api_version: str = "2024-02-29-preview"
    
    # Locale & Language
    locale: str = "en-US"
    language: str = ""  # Auto-detect if empty
    
    # Features
    features: List[str] = field(default_factory=lambda: [
        "keyValuePairs", "languages"
    ])
    
    # Output format
    output_content_format: str = "markdown"  # text, markdown
    
    # Pages
    pages: str = ""  # e.g., "1-3,5,7-10" or empty for all
    
    # Query fields (for prebuilt-document)
    query_fields: List[str] = field(default_factory=list)
    
    # Reading order
    reading_order: str = "natural"  # natural, basic
    
    # Polling
    polling_interval: int = 1  # seconds
    max_polling_attempts: int = 120
    
    # ============== AWS Textract ==============
    
    aws_region: str = "us-east-1"
    aws_access_key: str = ""
    aws_secret_key: str = ""
    aws_session_token: str = ""  # For temporary credentials
    
    # Features
    textract_features: List[str] = field(default_factory=lambda: [
        "TABLES", "FORMS"
    ])
    
    # S3 settings (for async processing)
    s3_bucket: str = ""
    s3_prefix: str = "textract/"
    
    # SNS notification
    sns_topic_arn: str = ""
    role_arn: str = ""
    
    # Queries
    textract_queries: List[str] = field(default_factory=list)
    
    # ============== Common ==============
    
    # Timeout
    timeout: int = 300  # seconds
    
    # Confidence threshold
    min_confidence: float = 0.0
    
    # Output
    include_raw_response: bool = False
    include_bounding_boxes: bool = False
    
    def get_azure_features(self) -> List[str]:
        """Get Azure features as string list."""
        return [f.value if isinstance(f, AzureFeature) else f for f in self.features]
    
    def get_textract_features(self) -> List[str]:
        """Get Textract features as string list."""
        return [f.value if isinstance(f, TextractFeature) else f for f in self.textract_features]
    
    def validate(self) -> None:
        """Validate configuration."""
        if self.provider == DocIntelProvider.AZURE:
            if not self.endpoint:
                raise ValueError("Azure endpoint is required")
            if not self.api_key:
                raise ValueError("Azure API key is required")
        elif self.provider == DocIntelProvider.AWS:
            if not self.aws_region:
                raise ValueError("AWS region is required")
