package write_test

import (
	"context"
	"testing"

	"urlshortener/internal/write"
	"urlshortener/pkg/cache"
	"urlshortener/pkg/config"
)

func TestBloomFilter_Init(t *testing.T) {
	cfg := &config.Config{
		RedisAddr: "localhost:6379",
	}
	
	redisClient := cache.NewRedisClient(cfg)
	if redisClient != nil {
		defer redisClient.Close()
	}

	bf := write.NewBloomFilter(redisClient)
	if bf == nil {
		t.Error("Expected Bloom Filter wrapper to be initialized")
	}
}

// TestBloomFilter_Interface verifies that the BloomFilter has the expected methods.
// This is mostly to drive the TDD design of the interface.
func TestBloomFilter_Interface(t *testing.T) {
	var bf interface {
		Exists(ctx context.Context, item string) (bool, error)
		Add(ctx context.Context, item string) error
	} = write.NewBloomFilter(nil)

	if bf == nil {
		t.Log("BloomFilter is nil, which is fine for this nil-init test")
	}
}
