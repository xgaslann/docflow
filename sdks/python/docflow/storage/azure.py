"""Azure Blob Storage backend."""

from typing import BinaryIO, List, Optional

from .base import Storage


class AzureStorage(Storage):
    """Azure Blob Storage implementation.
    
    Requires azure-storage-blob: pip install azure-storage-blob
    
    Example:
        >>> storage = AzureStorage(
        ...     account_name="myaccount",
        ...     container_name="documents"
        ... )
        >>> storage.save("doc.pdf", pdf_bytes)
    """

    def __init__(
        self,
        account_name: str,
        container_name: str,
        account_key: Optional[str] = None,
        connection_string: Optional[str] = None,
        prefix: str = "",
    ) -> None:
        """Create Azure Blob storage.
        
        Args:
            account_name: Azure storage account name.
            container_name: Container name.
            account_key: Optional account key (uses DefaultAzureCredential if not provided).
            connection_string: Optional connection string.
            prefix: Optional prefix for all blob names.
        """
        try:
            from azure.storage.blob import BlobServiceClient, ContainerClient
        except ImportError:
            raise ImportError(
                "azure-storage-blob is required: pip install azure-storage-blob"
            )

        self.container_name = container_name
        self.prefix = prefix.rstrip("/")
        self.account_name = account_name

        if connection_string:
            self.blob_service = BlobServiceClient.from_connection_string(connection_string)
        elif account_key:
            account_url = f"https://{account_name}.blob.core.windows.net"
            self.blob_service = BlobServiceClient(account_url, credential=account_key)
        else:
            # Use DefaultAzureCredential
            from azure.identity import DefaultAzureCredential
            account_url = f"https://{account_name}.blob.core.windows.net"
            self.blob_service = BlobServiceClient(account_url, credential=DefaultAzureCredential())

        self.container_client = self.blob_service.get_container_client(container_name)

    def save(self, path: str, data: bytes) -> None:
        """Save data to Azure Blob."""
        blob_name = self._full_key(path)
        blob_client = self.container_client.get_blob_client(blob_name)
        blob_client.upload_blob(data, overwrite=True)

    def save_file(self, path: str, file: BinaryIO) -> None:
        """Save data from a file-like object to Azure Blob."""
        blob_name = self._full_key(path)
        blob_client = self.container_client.get_blob_client(blob_name)
        blob_client.upload_blob(file, overwrite=True)

    def load(self, path: str) -> bytes:
        """Load data from Azure Blob."""
        blob_name = self._full_key(path)
        blob_client = self.container_client.get_blob_client(blob_name)
        try:
            downloader = blob_client.download_blob()
            return downloader.readall()
        except Exception as e:
            raise FileNotFoundError(f"File not found: {path}") from e

    def delete(self, path: str) -> None:
        """Delete file from Azure Blob."""
        blob_name = self._full_key(path)
        blob_client = self.container_client.get_blob_client(blob_name)
        try:
            blob_client.delete_blob()
        except:
            pass  # Ignore if doesn't exist

    def exists(self, path: str) -> bool:
        """Check if file exists in Azure Blob."""
        blob_name = self._full_key(path)
        blob_client = self.container_client.get_blob_client(blob_name)
        return blob_client.exists()

    def list(self, directory: str = "") -> List[str]:
        """List files in Azure Blob prefix."""
        prefix = self._full_key(directory)
        if prefix and not prefix.endswith("/"):
            prefix += "/"

        blobs = self.container_client.list_blobs(name_starts_with=prefix)
        
        files = []
        for blob in blobs:
            rel_path = blob.name[len(prefix):] if prefix else blob.name
            if rel_path and "/" not in rel_path:
                files.append(rel_path)

        return files

    def get_url(self, path: str) -> Optional[str]:
        """Get Azure Blob URL for the file."""
        blob_name = self._full_key(path)
        return f"https://{self.account_name}.blob.core.windows.net/{self.container_name}/{blob_name}"

    def _full_key(self, path: str) -> str:
        """Get full blob name for a relative path."""
        if self.prefix:
            return f"{self.prefix}/{path}"
        return path
