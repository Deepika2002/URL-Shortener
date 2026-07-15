package read_test

import (
	"context"
	"testing"
	"time"

	"urlshortener/internal/read"
	"urlshortener/pkg/cache"
	"urlshortener/pkg/config"
)

func TestMutex_Init(t *testing.T) {
	cfg := &config.Config{
		RedisAddr: "127.0.0.1:6379",
	}

	redisClient := cache.NewRedisClient(cfg)
	if redisClient != nil {
		defer redisClient.Close()
	}

	m := read.NewMutex(redisClient)
	if m == nil {
		t.Error("Expected Mutex to be initialized")
	}
}

// TestMutex_Interface verifies the expected methods of the mutex wrapper for stampede prevention.
func TestMutex_Interface(t *testing.T) {
	var m interface {
		AcquireLock(ctx context.Context, shortID string, expiration time.Duration) (bool, error)
		ReleaseLock(ctx context.Context, shortID string) error
	} = read.NewMutex(nil)

	if m == nil {
		t.Log("Mutex is nil, which is fine for this nil-init test")
	}
}
