package cache_test

import (
	"testing"

	"urlshortener/pkg/cache"
	"urlshortener/pkg/config"
)

func TestNewRedisClient_Init(t *testing.T) {
	cfg := &config.Config{
		RedisAddr: "localhost:6379",
	}

	client := cache.NewRedisClient(cfg)
	if client == nil {
		t.Error("Expected Redis client to be initialized, got nil")
	}
	
	if client != nil {
		client.Close()
	}
}
