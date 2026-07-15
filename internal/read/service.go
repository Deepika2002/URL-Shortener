package read

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"urlshortener/pkg/kafka"
)

// Service defines the core business logic for the Read service (URL redirection).
type Service interface {
	GetLongURL(ctx context.Context, shortID string) (string, error)
}

type service struct {
	repo      Repository
	cache     Cache
	mutex     Mutex
	publisher kafka.Publisher
}

// NewService instantiates the core read service orchestration.
func NewService(repo Repository, cache Cache, mutex Mutex, publisher kafka.Publisher) Service {
	return &service{
		repo:      repo,
		cache:     cache,
		mutex:     mutex,
		publisher: publisher,
	}
}

// GetLongURL fetches the original URL using a Cache-Aside pattern with Distributed Locking
// to prevent Cache Stampedes.
func (s *service) GetLongURL(ctx context.Context, shortID string) (string, error) {
	const maxRetries = 3
	const retryDelay = 50 * time.Millisecond
	const lockExpiration = 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		// 1. Check the fast Cache (Redis)
		// Our Cache implementation uses GETEX, so hits automatically refresh their TTL.
		longURL, err := s.cache.Get(ctx, shortID)
		if err == nil && longURL != "" {
			s.publishRedirectEvent(shortID, longURL)
			return longURL, nil
		}

		// 2. Cache Miss -> Attempt to acquire distributed lock
		// This ensures only ONE request hits the database for this specific shortID.
		acquired, err := s.mutex.AcquireLock(ctx, shortID, lockExpiration)
		if err == nil && acquired {
			// We got the lock! Ensure we release it when done.
			defer s.mutex.ReleaseLock(ctx, shortID) // Note: using ctx, might need background ctx if parent cancels early

			// 2.1 Double-check cache (another request might have populated it while we acquired the lock)
			longURL, err = s.cache.Get(ctx, shortID)
			if err == nil && longURL != "" {
				s.publishRedirectEvent(shortID, longURL)
				return longURL, nil
			}

			// 2.2 Fetch from primary Database (ScyllaDB)
			longURL, err = s.repo.GetMapping(ctx, shortID)
			if err != nil {
				return "", err // err could be ErrNotFound, which the HTTP handler will translate to 404
			}

			// 2.3 Populate the Cache
			if err := s.cache.Set(ctx, shortID, longURL, 24*time.Hour); err != nil {
				log.Printf("Warning: failed to populate cache for %s: %v", shortID, err)
			}

			s.publishRedirectEvent(shortID, longURL)
			return longURL, nil
		}

		// 3. Lock not acquired -> Another request is currently querying the database.
		// Wait a tiny bit and loop around to check the cache again.
		time.Sleep(retryDelay)
	}

	return "", fmt.Errorf("timeout waiting for cache population for shortID: %s", shortID)
}

// publishRedirectEvent emits an analytics event to Kafka for asynchronous processing.
func (s *service) publishRedirectEvent(shortID, longURL string) {
	if s.publisher == nil {
		return
	}

	event := map[string]string{
		"short_id":    shortID,
		"long_url":    longURL,
		"accessed_at": time.Now().UTC().Format(time.RFC3339Nano),
	}
	eventBytes, _ := json.Marshal(event)

	// Fire and forget - don't block the critical path of returning the redirect.
	if err := s.publisher.Publish("url_redirected", shortID, eventBytes); err != nil {
		log.Printf("Warning: failed to publish redirect event to Kafka: %v", err)
	}
}
