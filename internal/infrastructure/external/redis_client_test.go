package external_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestRedisClientSet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	key := "test:key"
	value := `{"id":1,"name":"test"}`
	expiration := 30 * time.Minute

	mock.ExpectSet(key, value, expiration).SetVal("OK")

	result := db.Set(ctx, key, value, expiration)
	status, err := result.Result()

	assert.NoError(t, err)
	assert.Equal(t, "OK", status)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRedisClientGet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	key := "test:key"
	expectedValue := `{"id":1,"name":"test"}`

	mock.ExpectGet(key).SetVal(expectedValue)

	result := db.Get(ctx, key)
	val, err := result.Result()

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, val)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRedisClientDelete(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	key := "test:key"

	mock.ExpectDel(key).SetVal(1)

	result := db.Del(ctx, key)
	deletedCount, err := result.Result()

	assert.NoError(t, err)
	assert.Equal(t, int64(1), deletedCount)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRedisClientExists(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	key := "test:key"

	t.Run("key exists", func(t *testing.T) {
		mock.ExpectExists(key).SetVal(1)

		result := db.Exists(ctx, key)
		count, err := result.Result()

		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("key does not exist", func(t *testing.T) {
		mock.ExpectExists(key).SetVal(0)

		result := db.Exists(ctx, key)
		count, err := result.Result()

		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestRedisClientSetNX(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	key := "test:key"
	value := "test_value"
	expiration := 30 * time.Minute

	t.Run("key does not exist - set successful", func(t *testing.T) {
		mock.ExpectSetNX(key, value, expiration).SetVal(true)

		result := db.SetNX(ctx, key, value, expiration)
		success, err := result.Result()

		assert.NoError(t, err)
		assert.True(t, success)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("key already exists - set failed", func(t *testing.T) {
		mock.ExpectSetNX(key, value, expiration).SetVal(false)

		result := db.SetNX(ctx, key, value, expiration)
		success, err := result.Result()

		assert.NoError(t, err)
		assert.False(t, success)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestRedisClientIncr(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	key := "counter:key"

	mock.ExpectIncr(key).SetVal(1)

	result := db.Incr(ctx, key)
	count, err := result.Result()

	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRedisClientExpire(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	key := "test:key"
	expiration := 1 * time.Hour

	mock.ExpectExpire(key, expiration).SetVal(true)

	result := db.Expire(ctx, key, expiration)
	success, err := result.Result()

	assert.NoError(t, err)
	assert.True(t, success)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRedisClientPing(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func() { _ = db.Close() }()

	ctx := context.Background()

	mock.ExpectPing().SetVal("PONG")

	result := db.Ping(ctx)
	pong, err := result.Result()

	assert.NoError(t, err)
	assert.Equal(t, "PONG", pong)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRedisClientErrorHandling(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	key := "test:key"

	t.Run("connection error", func(t *testing.T) {
		mock.ExpectGet(key).SetErr(redis.Nil)

		result := db.Get(ctx, key)
		_, err := result.Result()

		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("timeout error", func(t *testing.T) {
		timeoutErr := context.DeadlineExceeded
		mock.ExpectSet(key, "value", time.Minute).SetErr(timeoutErr)

		result := db.Set(ctx, key, "value", time.Minute)
		err := result.Err()

		assert.Error(t, err)
		assert.Equal(t, timeoutErr, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestComplexDataComparison(t *testing.T) {
	type UserData struct {
		ID       int                    `json:"id"`
		Name     string                 `json:"name"`
		Email    string                 `json:"email"`
		Metadata map[string]interface{} `json:"metadata"`
		Tags     []string               `json:"tags"`
	}

	expected := UserData{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
		Metadata: map[string]interface{}{
			"role":        "admin",
			"permissions": []string{"read", "write"},
			"settings": map[string]interface{}{
				"theme": "dark",
				"lang":  "ja",
			},
		},
		Tags: []string{"premium", "verified"},
	}

	actual := UserData{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
		Metadata: map[string]interface{}{
			"role":        "admin",
			"permissions": []string{"read", "write"},
			"settings": map[string]interface{}{
				"theme": "dark",
				"lang":  "ja",
			},
		},
		Tags: []string{"premium", "verified"},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("UserData mismatch (-want +got):\n%s", diff)
	}

	opts := cmp.Options{
		cmp.FilterPath(func(p cmp.Path) bool {
			return p.String() == "ID"
		}, cmp.Ignore()),
	}

	if diff := cmp.Diff(expected, actual, opts); diff != "" {
		t.Errorf("UserData mismatch (ignoring ID) (-want +got):\n%s", diff)
	}
}

func TestNestedStructComparison(t *testing.T) {
	type Address struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		Country string `json:"country"`
	}

	type User struct {
		ID      int     `json:"id"`
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	expected := User{
		ID:   1,
		Name: "John Doe",
		Address: Address{
			Street:  "123 Main St",
			City:    "Tokyo",
			Country: "Japan",
		},
	}

	actual := User{
		ID:   1,
		Name: "John Doe",
		Address: Address{
			Street:  "123 Main St",
			City:    "Tokyo",
			Country: "Japan",
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("User struct mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(expected.Address, actual.Address); diff != "" {
		t.Errorf("Address mismatch (-want +got):\n%s", diff)
	}
}

func TestSliceComparisonIgnoreOrder(t *testing.T) {
	expected := []string{"apple", "banana", "cherry"}
	actual := []string{"cherry", "apple", "banana"}

	opts := cmp.Options{
		cmp.Transformer("Sort", func(in []string) []string {
			out := make([]string, len(in))
			copy(out, in)
			for i := 0; i < len(out); i++ {
				for j := i + 1; j < len(out); j++ {
					if out[i] > out[j] {
						out[i], out[j] = out[j], out[i]
					}
				}
			}
			return out
		}),
	}

	if diff := cmp.Diff(expected, actual, opts); diff != "" {
		t.Errorf("Slice mismatch (ignoring order) (-want +got):\n%s", diff)
	}
}
