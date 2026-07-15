package write_test

import (
	"context"
	"testing"
	"time"

	"urlshortener/internal/write"
	"urlshortener/pkg/cache"
	"urlshortener/pkg/config"
)

func TestIdempotencyCache_Init(t *testing.T) {
	cfg := &config.Config{
		RedisAddr: "localhost:6379",
	}

	redisClient := cache.NewRedisClient(cfg)
	if redisClient != nil {
		defer redisClient.Close()
	}

	ic := write.NewIdempotencyCache(redisClient)
	if ic == nil {
		t.Error("Expected IdempotencyCache to be initialized")
	}
}

// TestIdempotencyCache_Interface verifies the expected methods of the cache.
// This is to drive the TDD design of the interface.
func TestIdempotencyCache_Interface(t *testing.T) {
	var ic interface {
		Get(ctx context.Context, longURL string) (string, error)
		Set(ctx context.Context, longURL string, shortURL string, expiration time.Duration) error
	} = write.NewIdempotencyCache(nil)

	if ic == nil {
		t.Log("IdempotencyCache is nil, which is fine for this nil-init test")
	}
}
