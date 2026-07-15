package read

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const cachePrefix = "urlshortener:mapping:"
const defaultTTL = 24 * time.Hour

// Cache defines the caching interface for the Read service.
type Cache interface {
	Get(ctx context.Context, shortID string) (string, error)
	Set(ctx context.Context, shortID string, longURL string, expiration time.Duration) error
}

type cacheImpl struct {
	client *redis.Client
}

// NewCache creates a new Redis-backed cache for URL mappings.
func NewCache(client *redis.Client) Cache {
	return &cacheImpl{
		client: client,
	}
}

// Get fetches the long URL by short ID and automatically refreshes its TTL to keep it hot.
func (c *cacheImpl) Get(ctx context.Context, shortID string) (string, error) {
	key := cachePrefix + shortID
	
	// GetEx fetches the value and updates the expiration atomically (requires Redis 6.2+)
	val, err := c.client.GetEx(ctx, key, defaultTTL).Result()
	if err == redis.Nil {
		return "", nil // Cache miss, completely normal
	}
	if err != nil {
		return "", fmt.Errorf("failed to fetch from cache: %w", err)
	}

	return val, nil
}

// Set stores a new shortID to longURL mapping with an explicit expiration.
func (c *cacheImpl) Set(ctx context.Context, shortID string, longURL string, expiration time.Duration) error {
	key := cachePrefix + shortID
	err := c.client.Set(ctx, key, longURL, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}
	return nil
}
