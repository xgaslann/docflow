package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage implements Storage interface for local filesystem.
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new local storage instance.
// basePath is the root directory for all operations.
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	return &LocalStorage{basePath: absPath}, nil
}

// Save stores data at the given path.
func (s *LocalStorage) Save(path string, data []byte) error {
	fullPath := s.fullPath(path)

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// SaveReader stores data from a reader at the given path.
func (s *LocalStorage) SaveReader(path string, reader io.Reader) error {
	fullPath := s.fullPath(path)

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Load retrieves data from the given path.
func (s *LocalStorage) Load(path string) ([]byte, error) {
	fullPath := s.fullPath(path)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// LoadReader returns a reader for the data at the given path.
func (s *LocalStorage) LoadReader(path string) (io.ReadCloser, error) {
	fullPath := s.fullPath(path)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete removes the file at the given path.
func (s *LocalStorage) Delete(path string) error {
	fullPath := s.fullPath(path)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already deleted
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Exists checks if a file exists at the given path.
func (s *LocalStorage) Exists(path string) (bool, error) {
	fullPath := s.fullPath(path)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file: %w", err)
	}

	return true, nil
}

// List returns all files in the given directory.
func (s *LocalStorage) List(dir string) ([]string, error) {
	fullPath := s.fullPath(dir)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to list directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// GetURL returns the file path as a file:// URL.
func (s *LocalStorage) GetURL(path string) string {
	return "file://" + s.fullPath(path)
}

// BasePath returns the base path of the storage.
func (s *LocalStorage) BasePath() string {
	return s.basePath
}

func (s *LocalStorage) fullPath(path string) string {
	return filepath.Join(s.basePath, path)
}
