package external

import (
	"context"
	"fmt"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
)

type CacheService struct {
	redis *RedisClient
}

func NewCacheService(redis *RedisClient) *CacheService {
	return &CacheService{
		redis: redis,
	}
}

func (c *CacheService) SetUser(ctx context.Context, user *entity.User) error {
	key := fmt.Sprintf("user:%d", user.ID)
	return c.redis.Set(ctx, key, user, 30*time.Minute)
}

func (c *CacheService) GetUser(ctx context.Context, userID uint) (*entity.User, error) {
	key := fmt.Sprintf("user:%d", userID)
	var user entity.User
	err := c.redis.Get(ctx, key, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *CacheService) DeleteUser(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("user:%d", userID)
	return c.redis.Delete(ctx, key)
}

func (c *CacheService) SetAuth(ctx context.Context, auth *entity.Auth) error {
	key := fmt.Sprintf("auth:user:%d", auth.UserID)
	return c.redis.Set(ctx, key, auth, 30*time.Minute)
}

func (c *CacheService) GetAuth(ctx context.Context, userID uint) (*entity.Auth, error) {
	key := fmt.Sprintf("auth:user:%d", userID)
	var auth entity.Auth
	err := c.redis.Get(ctx, key, &auth)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func (c *CacheService) DeleteAuth(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("auth:user:%d", userID)
	return c.redis.Delete(ctx, key)
}

func (c *CacheService) SetSession(ctx context.Context, sessionID string, userID uint, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	sessionData := map[string]interface{}{
		"user_id":    userID,
		"created_at": time.Now(),
	}
	return c.redis.Set(ctx, key, sessionData, expiration)
}

func (c *CacheService) GetSession(ctx context.Context, sessionID string) (uint, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	var sessionData map[string]interface{}
	err := c.redis.Get(ctx, key, &sessionData)
	if err != nil {
		return 0, err
	}

	userID, ok := sessionData["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid session data")
	}

	return uint(userID), nil
}

func (c *CacheService) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.redis.Delete(ctx, key)
}

func (c *CacheService) IncrementRateLimit(ctx context.Context, key string, window time.Duration) (int64, error) {
	fullKey := fmt.Sprintf("rate_limit:%s", key)
	count, err := c.redis.Incr(ctx, fullKey)
	if err != nil {
		return 0, err
	}

	if count == 1 {
		err = c.redis.Expire(ctx, fullKey, window)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (c *CacheService) GetRateLimit(ctx context.Context, key string) (int64, error) {
	fullKey := fmt.Sprintf("rate_limit:%s", key)
	val, err := c.redis.client.Get(ctx, fullKey).Int64()
	if err != nil {
		if err.Error() == "redis: nil" {
			return 0, nil
		}
		return 0, err
	}
	return val, nil
}

func (c *CacheService) AddToBlacklist(ctx context.Context, ip string, expiration time.Duration) error {
	key := fmt.Sprintf("blacklist:ip:%s", ip)
	return c.redis.Set(ctx, key, true, expiration)
}

func (c *CacheService) IsBlacklisted(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("blacklist:ip:%s", ip)
	return c.redis.Exists(ctx, key)
}

func (c *CacheService) RemoveFromBlacklist(ctx context.Context, ip string) error {
	key := fmt.Sprintf("blacklist:ip:%s", ip)
	return c.redis.Delete(ctx, key)
}

func (c *CacheService) BlacklistToken(ctx context.Context, tokenID string, expiration time.Duration) error {
	key := fmt.Sprintf("blacklist:token:%s", tokenID)
	return c.redis.Set(ctx, key, true, expiration)
}

func (c *CacheService) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := fmt.Sprintf("blacklist:token:%s", tokenID)
	return c.redis.Exists(ctx, key)
}
