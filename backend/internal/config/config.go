package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	PDF      PDFConfig
	Storage  StorageConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	BodyLimit    int
}

type PDFConfig struct {
	PageSize     string
	MarginTop    string
	MarginBottom string
	MarginLeft   string
	MarginRight  string
}

type StorageConfig struct {
	TempDir   string
	OutputDir string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", ""),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 60*time.Second),
			BodyLimit:    getIntEnv("SERVER_BODY_LIMIT", 50*1024*1024),
		},
		PDF: PDFConfig{
			PageSize:     getEnv("PDF_PAGE_SIZE", "A4"),
			MarginTop:    getEnv("PDF_MARGIN_TOP", "20mm"),
			MarginBottom: getEnv("PDF_MARGIN_BOTTOM", "20mm"),
			MarginLeft:   getEnv("PDF_MARGIN_LEFT", "20mm"),
			MarginRight:  getEnv("PDF_MARGIN_RIGHT", "20mm"),
		},
		Storage: StorageConfig{
			TempDir:   getEnv("STORAGE_TEMP_DIR", "./temp"),
			OutputDir: getEnv("STORAGE_OUTPUT_DIR", "./output"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
