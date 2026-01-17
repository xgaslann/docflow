"""Storage backends for DocFlow."""

from .base import Storage
from .local import LocalStorage

__all__ = ["Storage", "LocalStorage"]

# Optional imports for cloud storage
try:
    from .s3 import S3Storage
    __all__.append("S3Storage")
except ImportError:
    pass

try:
    from .azure import AzureStorage
    __all__.append("AzureStorage")
except ImportError:
    pass

try:
    from .gcs import GCSStorage
    __all__.append("GCSStorage")
except ImportError:
    pass
