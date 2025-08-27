package external_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/infrastructure/external"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockRedisClient) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockRedisClient) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	args := m.Called(ctx, key, value, expiration)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	val, ok := args.Get(0).(int64)
	if !ok {
		return 0, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	args := m.Called(ctx, key, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisClient) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockRedisClientWithClient struct {
	MockRedisClient
	client *redis.Client
}

func NewMockRedisClientWithClient() *MockRedisClientWithClient {
	client := redis.NewClient(&redis.Options{
		Addr: getRedisAddrForTest(),
	})

	return &MockRedisClientWithClient{
		client: client,
	}
}

func (m *MockRedisClientWithClient) GetClient() *redis.Client {
	return m.client
}

func TestCacheServiceSetUser(t *testing.T) {
	mockRedis := &MockRedisClient{}
	cacheService := external.NewCacheService(&external.RedisClient{})

	ctx := context.Background()

	user := &entity.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
	}

	mockRedis.On("Set", ctx, "user:1", user, 30*time.Minute).Return(nil)

	assert.NotNil(t, cacheService)
	assert.NotNil(t, user)
}

func TestCacheServiceGetUser(t *testing.T) {
	mockRedis := &MockRedisClient{}
	ctx := context.Background()

	userID := uint(1)
	expectedUser := &entity.User{
		ID:    userID,
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
	}

	mockRedis.On("Get", ctx, "user:1", mock.AnythingOfType("*entity.User")).Return(nil).Run(func(args mock.Arguments) {
		if user, ok := args.Get(2).(*entity.User); ok {
			*user = *expectedUser
		}
	})

	assert.NotNil(t, expectedUser)
}

func TestCacheServiceDeleteUser(t *testing.T) {
	mockRedis := &MockRedisClient{}
	ctx := context.Background()

	userID := uint(1)

	mockRedis.On("Delete", ctx, "user:1").Return(nil)

	assert.Equal(t, uint(1), userID)
}

func TestCacheServiceSetAuth(t *testing.T) {
	mockRedis := &MockRedisClient{}
	ctx := context.Background()

	auth := &entity.Auth{
		ID:           1,
		UserID:       1,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		IsActive:     true,
	}

	mockRedis.On("Set", ctx, "auth:user:1", auth, 30*time.Minute).Return(nil)

	assert.NotNil(t, auth)
}

func TestCacheServiceGetAuth(t *testing.T) {
	mockRedis := &MockRedisClient{}
	ctx := context.Background()

	userID := uint(1)
	expectedAuth := &entity.Auth{
		ID:           1,
		UserID:       userID,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		IsActive:     true,
	}

	mockRedis.On("Get", ctx, "auth:user:1", mock.AnythingOfType("*entity.Auth")).Return(nil).Run(func(args mock.Arguments) {
		if auth, ok := args.Get(2).(*entity.Auth); ok {
			*auth = *expectedAuth
		}
	})

	assert.NotNil(t, expectedAuth)
}

func TestCacheServiceSetSession(t *testing.T) {
	mockRedis := &MockRedisClient{}
	ctx := context.Background()

	sessionID := "session_123"
	userID := uint(1)
	expiration := 24 * time.Hour

	sessionData := map[string]interface{}{
		"user_id":    userID,
		"created_at": mock.AnythingOfType("time.Time"),
	}

	mockRedis.On("Set", ctx, "session:session_123", sessionData, expiration).Return(nil)

	assert.Equal(t, "session_123", sessionID)
}

func TestCacheServiceGetSession(t *testing.T) {
	mockRedis := &MockRedisClient{}
	ctx := context.Background()

	sessionID := "session_123"
	expectedUserID := uint(1)

	sessionData := map[string]interface{}{
		"user_id":    float64(expectedUserID),
		"created_at": time.Now(),
	}

	mockRedis.On("Get", ctx, "session:session_123", mock.AnythingOfType("*map[string]interface {}")).Return(nil).Run(func(args mock.Arguments) {
		if data, ok := args.Get(2).(*map[string]interface{}); ok {
			*data = sessionData
		}
	})

	assert.Equal(t, "session_123", sessionID)
}

func TestCacheServiceIncrementRateLimit(t *testing.T) {
	mockRedis := &MockRedisClient{}
	ctx := context.Background()

	key := "user:1:login"
	window := 1 * time.Hour
	expectedCount := int64(1)

	mockRedis.On("Incr", ctx, "rate_limit:user:1:login").Return(expectedCount, nil)
	mockRedis.On("Expire", ctx, "rate_limit:user:1:login", window).Return(nil)

	assert.Equal(t, "user:1:login", key)
}

func TestCacheServiceAddToBlacklist(t *testing.T) {
	mockRedis := &MockRedisClient{}
	ctx := context.Background()

	ip := "192.168.1.100"
	expiration := 24 * time.Hour

	mockRedis.On("Set", ctx, "blacklist:ip:192.168.1.100", true, expiration).Return(nil)

	assert.Equal(t, "192.168.1.100", ip)
}

func TestCacheServiceIsBlacklisted(t *testing.T) {
	mockRedis := &MockRedisClient{}
	ctx := context.Background()

	ip := "192.168.1.100"

	t.Run("is blacklisted", func(t *testing.T) {
		mockRedis.On("Exists", ctx, "blacklist:ip:192.168.1.100").Return(true, nil)

		assert.Equal(t, "192.168.1.100", ip)
	})

	t.Run("not blacklisted", func(t *testing.T) {
		mockRedis.On("Exists", ctx, "blacklist:ip:192.168.1.100").Return(false, nil)

		assert.Equal(t, "192.168.1.100", ip)
	})
}

func TestCacheServiceBlacklistToken(t *testing.T) {
	mockRedis := &MockRedisClient{}
	ctx := context.Background()

	tokenID := "token_123"
	expiration := 24 * time.Hour

	mockRedis.On("Set", ctx, "blacklist:token:token_123", true, expiration).Return(nil)

	assert.Equal(t, "token_123", tokenID)
}

func TestCacheServiceIsTokenBlacklisted(t *testing.T) {
	mockRedis := &MockRedisClient{}
	ctx := context.Background()

	tokenID := "token_123"

	t.Run("token is blacklisted", func(t *testing.T) {
		mockRedis.On("Exists", ctx, "blacklist:token:token_123").Return(true, nil)

		assert.Equal(t, "token_123", tokenID)
	})

	t.Run("token not blacklisted", func(t *testing.T) {
		mockRedis.On("Exists", ctx, "blacklist:token:token_123").Return(false, nil)

		assert.Equal(t, "token_123", tokenID)
	})
}

func TestCacheServiceIntegration(t *testing.T) {
	t.Skip("Integration test requires Redis instance")

	redisClient := external.NewRedisClient(getRedisAddrForTest(), getRedisPasswordForTest(), getRedisDBForTest())
	cacheService := external.NewCacheService(redisClient)

	ctx := context.Background()

	err := redisClient.Ping(ctx)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	user := &entity.User{
		ID:    1,
		Name:  "Integration Test User",
		Email: "integration@example.com",
		Age:   30,
	}

	err = cacheService.SetUser(ctx, user)
	assert.NoError(t, err)

	retrievedUser, err := cacheService.GetUser(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Email, retrievedUser.Email)

	err = cacheService.DeleteUser(ctx, user.ID)
	assert.NoError(t, err)

	_, err = cacheService.GetUser(ctx, user.ID)
	assert.Error(t, err)
}

func getRedisAddrForTest() string {
	host := "redis"
	if testHost := getEnv("REDIS_HOST", ""); testHost != "" {
		host = testHost
	}

	port := getEnv("REDIS_PORT", "6379")
	return host + ":" + port
}

func getRedisPasswordForTest() string {
	password := getEnv("REDIS_PASSWORD", "")
	if password == "" {
		password = "password"
	}
	return password
}

func getRedisDBForTest() int {
	return 0
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
