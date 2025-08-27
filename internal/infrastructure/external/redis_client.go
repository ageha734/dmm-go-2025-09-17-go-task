package external

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	return NewRedisClientWithRetry(addr, password, db)
}

func NewRedisClientWithRetry(addr, password string, db int) *RedisClient {
	maxRetries := getRedisMaxRetries()
	retryInterval := getRedisRetryInterval()

	var rdb *redis.Client
	var err error

	for i := 0; i < maxRetries; i++ {
		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = rdb.Ping(ctx).Err()
		cancel()

		if err == nil {
			log.Printf("✅ Redisに正常に接続しました！ (試行回数: %d/%d)", i+1, maxRetries)
			return &RedisClient{
				client: rdb,
			}
		}

		log.Printf("❌ Redis接続に失敗しました (試行回数: %d/%d): %v", i+1, maxRetries, err)

		if rdb != nil {
			if closeErr := rdb.Close(); closeErr != nil {
				log.Printf("⚠️ Redis接続のクローズに失敗しました: %v", closeErr)
			}
		}

		if i < maxRetries-1 {
			log.Printf("⏳ %v秒後に再試行します...", retryInterval.Seconds())
			time.Sleep(retryInterval)
		}
	}

	log.Printf("❌ %d回の試行後もRedisに接続できませんでした", maxRetries)
	return &RedisClient{
		client: rdb,
	}
}

func getRedisMaxRetries() int {
	retries := os.Getenv("REDIS_MAX_RETRIES")
	if retries == "" {
		return 10
	}
	val, err := strconv.Atoi(retries)
	if err != nil {
		return 10
	}
	return val
}

func getRedisRetryInterval() time.Duration {
	interval := os.Getenv("REDIS_RETRY_INTERVAL_SECONDS")
	if interval == "" {
		return 5 * time.Second
	}
	val, err := strconv.Atoi(interval)
	if err != nil {
		return 5 * time.Second
	}
	return time.Duration(val) * time.Second
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}

func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return r.client.SetNX(ctx, key, data, expiration).Result()
}

func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
