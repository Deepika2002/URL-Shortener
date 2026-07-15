package write

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const idempotencyPrefix = "urlshortener:idempotency:"

// IdempotencyCache defines the interface for caching long URL to short URL mappings.
type IdempotencyCache interface {
	Get(ctx context.Context, longURL string) (string, error)
	Set(ctx context.Context, longURL string, shortURL string, expiration time.Duration) error
}

type idempotencyCache struct {
	client *redis.Client
}

// NewIdempotencyCache creates a new Redis-backed idempotency cache.
func NewIdempotencyCache(client *redis.Client) IdempotencyCache {
	return &idempotencyCache{
		client: client,
	}
}

// hashKey creates a safe redis key by hashing the long URL to avoid massive keys.
func hashKey(longURL string) string {
	hash := sha256.Sum256([]byte(longURL))
	return idempotencyPrefix + hex.EncodeToString(hash[:])
}

// Get returns the short URL if it exists, or an empty string and no error if not found.
func (i *idempotencyCache) Get(ctx context.Context, longURL string) (string, error) {
	key := hashKey(longURL)
	val, err := i.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Cache miss, perfectly normal
	}
	if err != nil {
		return "", fmt.Errorf("failed to fetch from idempotency cache: %w", err)
	}
	return val, nil
}

// Set stores the longURL to shortURL mapping with an expiration time.
func (i *idempotencyCache) Set(ctx context.Context, longURL string, shortURL string, expiration time.Duration) error {
	key := hashKey(longURL)
	err := i.client.Set(ctx, key, shortURL, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set idempotency cache: %w", err)
	}
	return nil
}
