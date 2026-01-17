"""AWS S3 storage backend."""

from io import BytesIO
from typing import BinaryIO, List, Optional

from .base import Storage


class S3Storage(Storage):
    """AWS S3 storage implementation.
    
    Requires boto3: pip install boto3
    
    Example:
        >>> storage = S3Storage(bucket="my-docs", region="eu-west-1")
        >>> storage.save("doc.pdf", pdf_bytes)
    """

    def __init__(
        self,
        bucket: str,
        region: str = "us-east-1",
        prefix: str = "",
        endpoint_url: Optional[str] = None,
    ) -> None:
        """Create S3 storage.
        
        Args:
            bucket: S3 bucket name.
            region: AWS region.
            prefix: Optional prefix for all keys.
            endpoint_url: Optional custom endpoint (for MinIO, LocalStack).
        """
        try:
            import boto3
        except ImportError:
            raise ImportError("boto3 is required for S3 storage: pip install boto3")

        self.bucket = bucket
        self.prefix = prefix.rstrip("/")
        self.region = region

        self.client = boto3.client(
            "s3",
            region_name=region,
            endpoint_url=endpoint_url,
        )

    def save(self, path: str, data: bytes) -> None:
        """Save data to S3."""
        key = self._full_key(path)
        self.client.put_object(
            Bucket=self.bucket,
            Key=key,
            Body=data,
        )

    def save_file(self, path: str, file: BinaryIO) -> None:
        """Save data from a file-like object to S3."""
        key = self._full_key(path)
        self.client.upload_fileobj(file, self.bucket, key)

    def load(self, path: str) -> bytes:
        """Load data from S3."""
        key = self._full_key(path)
        try:
            response = self.client.get_object(Bucket=self.bucket, Key=key)
            return response["Body"].read()
        except self.client.exceptions.NoSuchKey:
            raise FileNotFoundError(f"File not found: {path}")

    def delete(self, path: str) -> None:
        """Delete file from S3."""
        key = self._full_key(path)
        self.client.delete_object(Bucket=self.bucket, Key=key)

    def exists(self, path: str) -> bool:
        """Check if file exists in S3."""
        key = self._full_key(path)
        try:
            self.client.head_object(Bucket=self.bucket, Key=key)
            return True
        except:
            return False

    def list(self, directory: str = "") -> List[str]:
        """List files in S3 prefix."""
        prefix = self._full_key(directory)
        if prefix and not prefix.endswith("/"):
            prefix += "/"

        response = self.client.list_objects_v2(
            Bucket=self.bucket,
            Prefix=prefix,
        )

        files = []
        for obj in response.get("Contents", []):
            key = obj["Key"]
            # Remove prefix to get relative path
            rel_path = key[len(prefix):] if prefix else key
            # Only include files (not nested)
            if rel_path and "/" not in rel_path:
                files.append(rel_path)

        return files

    def get_url(self, path: str) -> Optional[str]:
        """Get S3 URL for the file."""
        key = self._full_key(path)
        return f"s3://{self.bucket}/{key}"

    def get_http_url(self, path: str) -> str:
        """Get HTTP URL for the file."""
        key = self._full_key(path)
        return f"https://{self.bucket}.s3.{self.region}.amazonaws.com/{key}"

    def _full_key(self, path: str) -> str:
        """Get full S3 key for a relative path."""
        if self.prefix:
            return f"{self.prefix}/{path}"
        return path
