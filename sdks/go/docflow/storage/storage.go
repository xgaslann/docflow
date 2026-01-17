// Package storage provides storage backends for DocFlow.
package storage

import (
	"io"
)

// Storage defines the interface for storing and retrieving files.
type Storage interface {
	// Save stores data at the given path.
	Save(path string, data []byte) error

	// SaveReader stores data from a reader at the given path.
	SaveReader(path string, reader io.Reader) error

	// Load retrieves data from the given path.
	Load(path string) ([]byte, error)

	// LoadReader returns a reader for the data at the given path.
	LoadReader(path string) (io.ReadCloser, error)

	// Delete removes the file at the given path.
	Delete(path string) error

	// Exists checks if a file exists at the given path.
	Exists(path string) (bool, error)

	// List returns all files in the given directory.
	List(dir string) ([]string, error)

	// GetURL returns a URL for accessing the file (if supported).
	// Returns empty string if not supported.
	GetURL(path string) string
}

// StorageType represents the type of storage backend.
type StorageType string

const (
	StorageTypeLocal StorageType = "local"
	StorageTypeS3    StorageType = "s3"
	StorageTypeGCS   StorageType = "gcs"
	StorageTypeAzure StorageType = "azure"
)
