package read_test

import (
	"context"
	"testing"
	"time"

	"urlshortener/internal/read"
	"urlshortener/pkg/cache"
	"urlshortener/pkg/config"
)

func TestCache_Init(t *testing.T) {
	cfg := &config.Config{
		RedisAddr: "127.0.0.1:6379",
	}

	redisClient := cache.NewRedisClient(cfg)
	if redisClient != nil {
		defer redisClient.Close()
	}

	c := read.NewCache(redisClient)
	if c == nil {
		t.Error("Expected Cache to be initialized")
	}
}

// TestCache_Interface verifies the expected methods of the cache for the read service.
// This is mostly to drive the TDD design of the interface.
func TestCache_Interface(t *testing.T) {
	var c interface {
		Get(ctx context.Context, shortID string) (string, error)
		Set(ctx context.Context, shortID string, longURL string, expiration time.Duration) error
	} = read.NewCache(nil)

	if c == nil {
		t.Log("Cache is nil, which is fine for this nil-init test")
	}
}
