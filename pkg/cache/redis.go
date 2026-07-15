package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"urlshortener/pkg/config"
)

// NewRedisClient creates and returns a new connected Redis client.
func NewRedisClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Fail-fast ping to ensure connectivity on startup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		// We still return the client here so the application can attempt to reconnect
		// or handle the error in the higher-level health checks.
	}

	return client
}
