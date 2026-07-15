package write

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const bloomFilterKey = "urlshortener:bloom:longurls"

// BloomFilter defines the interface for our probabilistic data structure.
type BloomFilter interface {
	Exists(ctx context.Context, item string) (bool, error)
	Add(ctx context.Context, item string) error
}

type bloomFilter struct {
	client *redis.Client
}

// NewBloomFilter creates a new BloomFilter wrapper utilizing RedisBloom via go-redis.
func NewBloomFilter(client *redis.Client) BloomFilter {
	return &bloomFilter{
		client: client,
	}
}

// Exists checks if an item might exist in the Bloom filter.
// It returns true if the item *might* exist, and false if it *definitely* does not exist.
func (b *bloomFilter) Exists(ctx context.Context, item string) (bool, error) {
	// We use Do to send raw RedisBloom commands since they might not be fully mapped in all client versions.
	res, err := b.client.Do(ctx, "BF.EXISTS", bloomFilterKey, item).Bool()
	if err != nil {
		return false, fmt.Errorf("failed to check bloom filter existence: %w", err)
	}
	
	return res, nil
}

// Add safely adds a new item to the Bloom filter.
func (b *bloomFilter) Add(ctx context.Context, item string) error {
	_, err := b.client.Do(ctx, "BF.ADD", bloomFilterKey, item).Result()
	if err != nil {
		return fmt.Errorf("failed to add item to bloom filter: %w", err)
	}
	return nil
}
