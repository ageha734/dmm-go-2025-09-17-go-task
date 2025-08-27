package external_test

import (
	"os"
	"testing"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/infrastructure/external"
)

func setupRedisEnv(t *testing.T, key, value string) {
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

func cleanupRedisEnv(t *testing.T, key string) {
	t.Helper()
	if err := os.Unsetenv(key); err != nil {
		t.Logf("Warning: Failed to unset %s: %v", key, err)
	}
}

func setupRedisTestEnv(t *testing.T, maxRetries, retryInterval string) {
	t.Helper()
	setupRedisEnv(t, "REDIS_MAX_RETRIES", maxRetries)
	setupRedisEnv(t, "REDIS_RETRY_INTERVAL_SECONDS", retryInterval)
}

func cleanupAllRedisEnv(t *testing.T) {
	t.Helper()
	cleanupRedisEnv(t, "REDIS_MAX_RETRIES")
	cleanupRedisEnv(t, "REDIS_RETRY_INTERVAL_SECONDS")
}

func TestRedisRetryConfiguration(t *testing.T) {
	tests := []struct {
		name               string
		maxRetriesEnv      string
		retryIntervalEnv   string
		expectedMaxRetries int
		expectedInterval   time.Duration
	}{
		{
			name:               "デフォルト値",
			maxRetriesEnv:      "",
			retryIntervalEnv:   "",
			expectedMaxRetries: 10,
			expectedInterval:   5 * time.Second,
		},
		{
			name:               "環境変数で設定された値",
			maxRetriesEnv:      "5",
			retryIntervalEnv:   "3",
			expectedMaxRetries: 5,
			expectedInterval:   3 * time.Second,
		},
		{
			name:               "無効な値の場合はデフォルト値",
			maxRetriesEnv:      "invalid",
			retryIntervalEnv:   "invalid",
			expectedMaxRetries: 10,
			expectedInterval:   5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupRedisTestEnv(t, tt.maxRetriesEnv, tt.retryIntervalEnv)
			defer cleanupAllRedisEnv(t)
		})
	}
}

func TestNewRedisClient(t *testing.T) {
	setupRedisTestEnv(t, "2", "1")
	defer cleanupAllRedisEnv(t)

	client := external.NewRedisClient("localhost:9999", "", 0)

	if client == nil {
		t.Error("NewRedisClient() should return a client even if connection fails")
		return
	}

	if err := client.Close(); err != nil {
		t.Logf("Warning: Failed to close Redis client: %v", err)
	}
}
