package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Storage implements Storage interface for AWS S3.
type S3Storage struct {
	client *s3.Client
	bucket string
	prefix string
	region string
}

// S3Config contains configuration for S3 storage.
type S3Config struct {
	Bucket   string
	Region   string
	Prefix   string // Optional prefix for all keys
	Endpoint string // Optional custom endpoint (for MinIO, LocalStack, etc.)
}

// NewS3Storage creates a new S3 storage instance.
func NewS3Storage(cfg S3Config) (*S3Storage, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket is required")
	}

	ctx := context.Background()

	// Load AWS config
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	clientOpts := []func(*s3.Options){}
	if cfg.Endpoint != "" {
		clientOpts = append(clientOpts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true
		})
	}

	client := s3.NewFromConfig(awsCfg, clientOpts...)

	return &S3Storage{
		client: client,
		bucket: cfg.Bucket,
		prefix: strings.TrimSuffix(cfg.Prefix, "/"),
		region: cfg.Region,
	}, nil
}

// Save stores data at the given path.
func (s *S3Storage) Save(path string, data []byte) error {
	return s.SaveReader(path, bytes.NewReader(data))
}

// SaveReader stores data from a reader at the given path.
func (s *S3Storage) SaveReader(filePath string, reader io.Reader) error {
	key := s.fullKey(filePath)

	// Read all data (S3 needs content length for some operations)
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	_, err = s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
}

// Load retrieves data from the given path.
func (s *S3Storage) Load(filePath string) ([]byte, error) {
	reader, err := s.LoadReader(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// LoadReader returns a reader for the data at the given path.
func (s *S3Storage) LoadReader(filePath string) (io.ReadCloser, error) {
	key := s.fullKey(filePath)

	result, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}

	return result.Body, nil
}

// Delete removes the file at the given path.
func (s *S3Storage) Delete(filePath string) error {
	key := s.fullKey(filePath)

	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

// Exists checks if a file exists at the given path.
func (s *S3Storage) Exists(filePath string) (bool, error) {
	key := s.fullKey(filePath)

	_, err := s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		// Check if it's a "not found" error
		return false, nil
	}

	return true, nil
}

// List returns all files in the given directory.
func (s *S3Storage) List(dir string) ([]string, error) {
	prefix := s.fullKey(dir)
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	result, err := s.client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var files []string
	for _, obj := range result.Contents {
		// Remove prefix to get relative path
		key := strings.TrimPrefix(*obj.Key, prefix)
		if key != "" && !strings.Contains(key, "/") {
			files = append(files, key)
		}
	}

	return files, nil
}

// GetURL returns the S3 URL for the file.
func (s *S3Storage) GetURL(filePath string) string {
	key := s.fullKey(filePath)
	return fmt.Sprintf("s3://%s/%s", s.bucket, key)
}

// GetHTTPURL returns the HTTP URL for the file.
func (s *S3Storage) GetHTTPURL(filePath string) string {
	key := s.fullKey(filePath)
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
}

func (s *S3Storage) fullKey(filePath string) string {
	if s.prefix == "" {
		return filePath
	}
	return path.Join(s.prefix, filePath)
}
