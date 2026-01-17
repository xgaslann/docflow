package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear any existing env vars
	envVars := []string{
		"SERVER_HOST", "SERVER_PORT", "SERVER_READ_TIMEOUT",
		"SERVER_WRITE_TIMEOUT", "SERVER_BODY_LIMIT",
		"STORAGE_TEMP_DIR", "STORAGE_OUTPUT_DIR",
	}
	for _, v := range envVars {
		os.Unsetenv(v)
	}

	cfg := Load()

	// Server defaults
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected host '0.0.0.0', got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("expected port '8080', got %q", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 30*time.Second {
		t.Errorf("expected read timeout 30s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.WriteTimeout != 120*time.Second {
		t.Errorf("expected write timeout 120s, got %v", cfg.Server.WriteTimeout)
	}
	if cfg.Server.BodyLimit != 50*1024*1024 {
		t.Errorf("expected body limit 50MB, got %d", cfg.Server.BodyLimit)
	}

	// Storage defaults
	if cfg.Storage.TempDir != "./temp" {
		t.Errorf("expected temp dir './temp', got %q", cfg.Storage.TempDir)
	}
	if cfg.Storage.OutputDir != "./output" {
		t.Errorf("expected output dir './output', got %q", cfg.Storage.OutputDir)
	}
}

func TestLoad_EnvironmentOverrides(t *testing.T) {
	// Set custom env vars
	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("STORAGE_TEMP_DIR", "/custom/temp")
	os.Setenv("STORAGE_OUTPUT_DIR", "/custom/output")

	defer func() {
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("STORAGE_TEMP_DIR")
		os.Unsetenv("STORAGE_OUTPUT_DIR")
	}()

	cfg := Load()

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected host '127.0.0.1', got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != "3000" {
		t.Errorf("expected port '3000', got %q", cfg.Server.Port)
	}
	if cfg.Storage.TempDir != "/custom/temp" {
		t.Errorf("expected temp dir '/custom/temp', got %q", cfg.Storage.TempDir)
	}
	if cfg.Storage.OutputDir != "/custom/output" {
		t.Errorf("expected output dir '/custom/output', got %q", cfg.Storage.OutputDir)
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		envValue   string
		setEnv     bool
		defaultVal string
		expected   string
	}{
		{
			name:       "env set",
			key:        "TEST_VAR_1",
			envValue:   "custom_value",
			setEnv:     true,
			defaultVal: "default",
			expected:   "custom_value",
		},
		{
			name:       "env not set",
			key:        "TEST_VAR_2",
			setEnv:     false,
			defaultVal: "default",
			expected:   "default",
		},
		{
			name:       "env empty",
			key:        "TEST_VAR_3",
			envValue:   "",
			setEnv:     true,
			defaultVal: "default",
			expected:   "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetEnvDuration(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		envValue   string
		setEnv     bool
		defaultVal time.Duration
		expected   time.Duration
	}{
		{
			name:       "valid duration",
			key:        "TEST_DURATION_1",
			envValue:   "5s",
			setEnv:     true,
			defaultVal: 10 * time.Second,
			expected:   5 * time.Second,
		},
		{
			name:       "not set",
			key:        "TEST_DURATION_2",
			setEnv:     false,
			defaultVal: 10 * time.Second,
			expected:   10 * time.Second,
		},
		{
			name:       "invalid duration",
			key:        "TEST_DURATION_3",
			envValue:   "invalid",
			setEnv:     true,
			defaultVal: 10 * time.Second,
			expected:   10 * time.Second,
		},
		{
			name:       "minutes",
			key:        "TEST_DURATION_4",
			envValue:   "2m",
			setEnv:     true,
			defaultVal: 1 * time.Minute,
			expected:   2 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvDuration(tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		envValue   string
		setEnv     bool
		defaultVal int
		expected   int
	}{
		{
			name:       "valid int",
			key:        "TEST_INT_1",
			envValue:   "42",
			setEnv:     true,
			defaultVal: 10,
			expected:   42,
		},
		{
			name:       "not set",
			key:        "TEST_INT_2",
			setEnv:     false,
			defaultVal: 10,
			expected:   10,
		},
		{
			name:       "invalid int",
			key:        "TEST_INT_3",
			envValue:   "not_a_number",
			setEnv:     true,
			defaultVal: 10,
			expected:   10,
		},
		{
			name:       "negative",
			key:        "TEST_INT_4",
			envValue:   "-5",
			setEnv:     true,
			defaultVal: 10,
			expected:   -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvInt(tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
