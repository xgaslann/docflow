package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// AzureStorage implements Storage interface for Azure Blob Storage.
type AzureStorage struct {
	client        *azblob.Client
	containerName string
	prefix        string
	accountName   string
}

// AzureConfig contains configuration for Azure Blob storage.
type AzureConfig struct {
	AccountName   string
	AccountKey    string // Optional, uses DefaultAzureCredential if empty
	ContainerName string
	Prefix        string // Optional prefix for all blob names
}

// NewAzureStorage creates a new Azure Blob storage instance.
func NewAzureStorage(cfg AzureConfig) (*AzureStorage, error) {
	if cfg.AccountName == "" {
		return nil, fmt.Errorf("account name is required")
	}
	if cfg.ContainerName == "" {
		return nil, fmt.Errorf("container name is required")
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", cfg.AccountName)

	var client *azblob.Client
	var err error

	if cfg.AccountKey != "" {
		// Use shared key credential
		cred, err := azblob.NewSharedKeyCredential(cfg.AccountName, cfg.AccountKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create credential: %w", err)
		}
		client, err = azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create client: %w", err)
		}
	} else {
		// Use default credential (requires azure-identity)
		client, err = azblob.NewClientWithNoCredential(serviceURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create client: %w", err)
		}
	}

	return &AzureStorage{
		client:        client,
		containerName: cfg.ContainerName,
		prefix:        strings.TrimSuffix(cfg.Prefix, "/"),
		accountName:   cfg.AccountName,
	}, nil
}

// Save stores data at the given path.
func (s *AzureStorage) Save(filePath string, data []byte) error {
	return s.SaveReader(filePath, bytes.NewReader(data))
}

// SaveReader stores data from a reader at the given path.
func (s *AzureStorage) SaveReader(filePath string, reader io.Reader) error {
	blobName := s.fullKey(filePath)

	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	_, err = s.client.UploadBuffer(context.Background(), s.containerName, blobName, data, nil)
	if err != nil {
		return fmt.Errorf("failed to upload to Azure: %w", err)
	}

	return nil
}

// Load retrieves data from the given path.
func (s *AzureStorage) Load(filePath string) ([]byte, error) {
	blobName := s.fullKey(filePath)

	resp, err := s.client.DownloadStream(context.Background(), s.containerName, blobName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download from Azure: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// LoadReader returns a reader for the data at the given path.
func (s *AzureStorage) LoadReader(filePath string) (io.ReadCloser, error) {
	blobName := s.fullKey(filePath)

	resp, err := s.client.DownloadStream(context.Background(), s.containerName, blobName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download from Azure: %w", err)
	}

	return resp.Body, nil
}

// Delete removes the file at the given path.
func (s *AzureStorage) Delete(filePath string) error {
	blobName := s.fullKey(filePath)

	_, err := s.client.DeleteBlob(context.Background(), s.containerName, blobName, nil)
	if err != nil {
		// Ignore if doesn't exist
		return nil
	}

	return nil
}

// Exists checks if a file exists at the given path.
func (s *AzureStorage) Exists(filePath string) (bool, error) {
	blobName := s.fullKey(filePath)

	_, err := s.client.DownloadStream(context.Background(), s.containerName, blobName, nil)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// List returns all files in the given directory.
func (s *AzureStorage) List(dir string) ([]string, error) {
	prefix := s.fullKey(dir)
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	pager := s.client.NewListBlobsFlatPager(s.containerName, &azblob.ListBlobsFlatOptions{
		Prefix: &prefix,
	})

	var files []string
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to list blobs: %w", err)
		}

		for _, blob := range resp.Segment.BlobItems {
			relPath := strings.TrimPrefix(*blob.Name, prefix)
			if relPath != "" && !strings.Contains(relPath, "/") {
				files = append(files, relPath)
			}
		}
	}

	return files, nil
}

// GetURL returns the Azure Blob URL for the file.
func (s *AzureStorage) GetURL(filePath string) string {
	blobName := s.fullKey(filePath)
	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", s.accountName, s.containerName, blobName)
}

func (s *AzureStorage) fullKey(filePath string) string {
	if s.prefix == "" {
		return filePath
	}
	return path.Join(s.prefix, filePath)
}
