package main

import (
	"os"
	"testing"
	"time"
)

func setupEnv(t *testing.T, key, value string) {
	t.Helper()
	if value != "" {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("Failed to set %s: %v", key, err)
		}
	} else {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("Failed to unset %s: %v", key, err)
		}
	}
}

func cleanupEnv(t *testing.T, key string) {
	t.Helper()
	if err := os.Unsetenv(key); err != nil {
		t.Logf("Warning: Failed to unset %s: %v", key, err)
	}
}

func TestGetDBMaxRetries(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected int
	}{
		{
			name:     "デフォルト値",
			envValue: "",
			expected: 10,
		},
		{
			name:     "環境変数で設定された値",
			envValue: "5",
			expected: 5,
		},
		{
			name:     "無効な値の場合はデフォルト値",
			envValue: "invalid",
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupEnv(t, "DB_MAX_RETRIES", tt.envValue)
			defer cleanupEnv(t, "DB_MAX_RETRIES")

			result := getDBMaxRetries()
			if result != tt.expected {
				t.Errorf("getDBMaxRetries() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetDBRetryInterval(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected time.Duration
	}{
		{
			name:     "デフォルト値",
			envValue: "",
			expected: 5 * time.Second,
		},
		{
			name:     "環境変数で設定された値",
			envValue: "3",
			expected: 3 * time.Second,
		},
		{
			name:     "無効な値の場合はデフォルト値",
			envValue: "invalid",
			expected: 5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupEnv(t, "DB_RETRY_INTERVAL_SECONDS", tt.envValue)
			defer cleanupEnv(t, "DB_RETRY_INTERVAL_SECONDS")

			result := getDBRetryInterval()
			if result != tt.expected {
				t.Errorf("getDBRetryInterval() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetDSN(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "デフォルト値",
			envValue: "",
			expected: "testuser:password@tcp(mysql:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name:     "環境変数で設定された値",
			envValue: "user:pass@tcp(localhost:3306)/mydb",
			expected: "user:pass@tcp(localhost:3306)/mydb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupEnv(t, "DATABASE_URL", tt.envValue)
			defer cleanupEnv(t, "DATABASE_URL")

			result := getDSN()
			if result != tt.expected {
				t.Errorf("getDSN() = %v, want %v", result, tt.expected)
			}
		})
	}
}
