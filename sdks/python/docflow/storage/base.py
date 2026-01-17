"""Abstract base class for storage backends."""

from abc import ABC, abstractmethod
from typing import BinaryIO, List, Optional


class Storage(ABC):
    """Abstract base class for file storage backends."""

    @abstractmethod
    def save(self, path: str, data: bytes) -> None:
        """Save data to the given path.
        
        Args:
            path: Relative path where to save the file.
            data: Binary data to save.
        """
        pass

    @abstractmethod
    def save_file(self, path: str, file: BinaryIO) -> None:
        """Save data from a file-like object.
        
        Args:
            path: Relative path where to save the file.
            file: File-like object to read from.
        """
        pass

    @abstractmethod
    def load(self, path: str) -> bytes:
        """Load data from the given path.
        
        Args:
            path: Relative path to the file.
            
        Returns:
            Binary file contents.
            
        Raises:
            FileNotFoundError: If file doesn't exist.
        """
        pass

    @abstractmethod
    def delete(self, path: str) -> None:
        """Delete file at the given path.
        
        Args:
            path: Relative path to the file.
        """
        pass

    @abstractmethod
    def exists(self, path: str) -> bool:
        """Check if file exists at the given path.
        
        Args:
            path: Relative path to check.
            
        Returns:
            True if file exists.
        """
        pass

    @abstractmethod
    def list(self, directory: str) -> List[str]:
        """List files in the given directory.
        
        Args:
            directory: Relative path to directory.
            
        Returns:
            List of filenames in the directory.
        """
        pass

    def get_url(self, path: str) -> Optional[str]:
        """Get URL for accessing the file.
        
        Args:
            path: Relative path to the file.
            
        Returns:
            URL string if supported, None otherwise.
        """
        return None
