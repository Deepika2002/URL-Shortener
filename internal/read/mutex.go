package read

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const mutexPrefix = "urlshortener:lock:"

// Mutex defines a distributed lock interface to prevent cache stampedes.
type Mutex interface {
	AcquireLock(ctx context.Context, shortID string, expiration time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, shortID string) error
}

type mutexImpl struct {
	client *redis.Client
}

// NewMutex creates a new Redis-backed distributed lock.
func NewMutex(client *redis.Client) Mutex {
	return &mutexImpl{
		client: client,
	}
}

// AcquireLock attempts to grab an exclusive lock for the given short ID.
func (m *mutexImpl) AcquireLock(ctx context.Context, shortID string, expiration time.Duration) (bool, error) {
	key := mutexPrefix + shortID
	// SetNX (SET if Not eXists) atomically sets the key only if it does not already exist
	acquired, err := m.client.SetNX(ctx, key, "locked", expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock for %s: %w", shortID, err)
	}
	return acquired, nil
}

// ReleaseLock frees the lock so other requests can grab it.
func (m *mutexImpl) ReleaseLock(ctx context.Context, shortID string) error {
	key := mutexPrefix + shortID
	// Note: For strict correctness, a Lua script should be used to verify ownership
	// before deleting. A simple Del is used here for brevity and speed in this architecture.
	err := m.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to release lock for %s: %w", shortID, err)
	}
	return nil
}
