"""Local filesystem storage backend."""

import os
from pathlib import Path
from typing import BinaryIO, List, Optional

from .base import Storage


class LocalStorage(Storage):
    """Local filesystem storage implementation.
    
    Example:
        >>> storage = LocalStorage("./output")
        >>> storage.save("doc.pdf", pdf_bytes)
        >>> exists = storage.exists("doc.pdf")
    """

    def __init__(self, base_path: str) -> None:
        """Create local storage at the given path.
        
        Args:
            base_path: Root directory for all operations.
        """
        self.base_path = Path(base_path).resolve()
        self.base_path.mkdir(parents=True, exist_ok=True)

    def save(self, path: str, data: bytes) -> None:
        """Save data to the given path."""
        full_path = self._full_path(path)
        full_path.parent.mkdir(parents=True, exist_ok=True)
        full_path.write_bytes(data)

    def save_file(self, path: str, file: BinaryIO) -> None:
        """Save data from a file-like object."""
        full_path = self._full_path(path)
        full_path.parent.mkdir(parents=True, exist_ok=True)
        with open(full_path, "wb") as f:
            while chunk := file.read(8192):
                f.write(chunk)

    def load(self, path: str) -> bytes:
        """Load data from the given path."""
        full_path = self._full_path(path)
        if not full_path.exists():
            raise FileNotFoundError(f"File not found: {path}")
        return full_path.read_bytes()

    def delete(self, path: str) -> None:
        """Delete file at the given path."""
        full_path = self._full_path(path)
        if full_path.exists():
            full_path.unlink()

    def exists(self, path: str) -> bool:
        """Check if file exists at the given path."""
        return self._full_path(path).exists()

    def list(self, directory: str = "") -> List[str]:
        """List files in the given directory."""
        dir_path = self._full_path(directory)
        if not dir_path.exists():
            return []
        return [f.name for f in dir_path.iterdir() if f.is_file()]

    def get_url(self, path: str) -> Optional[str]:
        """Get file:// URL for the file."""
        return f"file://{self._full_path(path)}"

    def get_absolute_path(self, path: str) -> str:
        """Get the absolute filesystem path."""
        return str(self._full_path(path))

    def _full_path(self, path: str) -> Path:
        """Get full path for a relative path."""
        return self.base_path / path
