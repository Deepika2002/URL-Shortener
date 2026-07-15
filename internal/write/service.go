package write

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"urlshortener/pkg/kafka"
	"urlshortener/pkg/kgsclient"
)

// Service defines the core business logic for the Write service.
type Service interface {
	ShortenURL(ctx context.Context, longURL string, customAlias string) (string, error)
}

type service struct {
	repo      Repository
	cache     IdempotencyCache
	bloom     BloomFilter
	kgs       kgsclient.Client
	publisher kafka.Publisher
}

// NewService instantiates the core shortening service.
func NewService(repo Repository, cache IdempotencyCache, bloom BloomFilter, kgs kgsclient.Client, publisher kafka.Publisher) Service {
	return &service{
		repo:      repo,
		cache:     cache,
		bloom:     bloom,
		kgs:       kgs,
		publisher: publisher,
	}
}

// ShortenURL orchestrates the generation and persistence of a short URL.
func (s *service) ShortenURL(ctx context.Context, longURL string, customAlias string) (string, error) {
	// 1. Check idempotency cache to avoid duplicate processing for the same longURL
	cachedShortID, err := s.cache.Get(ctx, longURL)
	if err == nil && cachedShortID != "" {
		// Found in cache, return immediately
		return cachedShortID, nil
	}

	var shortID string

	// 2. Determine shortID (Custom Alias vs Auto-Generated)
	if customAlias != "" {
		// Verify custom alias doesn't exist via Bloom Filter
		exists, err := s.bloom.Exists(ctx, customAlias)
		if err != nil {
			return "", fmt.Errorf("failed to check bloom filter for alias: %w", err)
		}
		if exists {
			// A custom alias collision (or false positive).
			// At scale, we reject immediately to protect the DB from query floods.
			return "", fmt.Errorf("custom alias '%s' is already in use", customAlias)
		}
		shortID = customAlias
	} else {
		// Fetch a pre-generated, base62 encoded Snowflake ID from the Key Generation Service
		id, err := s.kgs.GetNextID()
		if err != nil {
			return "", fmt.Errorf("failed to fetch ID from KGS: %w", err)
		}
		shortID = id
	}

	// 3. Save to primary database (ScyllaDB)
	if err := s.repo.SaveMapping(ctx, shortID, longURL); err != nil {
		return "", fmt.Errorf("failed to save mapping to DB: %w", err)
	}

	// 4. Update Idempotency Cache (expire in 24 hours to keep cache lean)
	if err := s.cache.Set(ctx, longURL, shortID, 24*time.Hour); err != nil {
		log.Printf("Warning: Failed to set idempotency cache for %s: %v", longURL, err)
	}

	// 5. Update Bloom Filter (add the alias/shortID to prevent future collisions)
	if err := s.bloom.Add(ctx, shortID); err != nil {
		log.Printf("Warning: Failed to update bloom filter for %s: %v", shortID, err)
	}

	// 6. Publish event to Kafka for asynchronous Analytics/Read-model processing
	event := map[string]string{
		"short_id":   shortID,
		"long_url":   longURL,
		"created_at": time.Now().UTC().Format(time.RFC3339Nano),
	}
	eventBytes, _ := json.Marshal(event)
	
	if s.publisher != nil {
		if err := s.publisher.Publish("url_created", shortID, eventBytes); err != nil {
			log.Printf("Warning: Failed to publish to Kafka: %v", err)
		}
	}

	return shortID, nil
}
